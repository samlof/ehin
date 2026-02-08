package resource

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/samlof/ehin/internal/db/model"
	"github.com/samlof/ehin/internal/nordpool"
	"github.com/samlof/ehin/internal/service"
	"github.com/samlof/ehin/internal/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPriceRepository struct {
	mock.Mock
}

func (m *MockPriceRepository) Select1(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPriceRepository) GetPrices(ctx context.Context, from, to time.Time) ([]model.PriceHistoryEntry, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PriceHistoryEntry), args.Error(1)
}

func (m *MockPriceRepository) InsertPrices(ctx context.Context, entries []model.PriceHistoryEntry) (int64, error) {
	args := m.Called(ctx, entries)
	return args.Get(0).(int64), args.Error(1)
}

type MockTimeProvider struct {
	mock.Mock
}

func (m *MockTimeProvider) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

// MockNordPoolClient is a mock implementation of nordpool.NordPoolClient
type MockNordPoolClient struct {
	mock.Mock
}

func (m *MockNordPoolClient) GetDayAheadPrices(date time.Time, market, deliveryArea, currency string) (*nordpool.PriceDataResponse, error) {
	args := m.Called(date, market, deliveryArea, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nordpool.PriceDataResponse), args.Error(1)
}

func TestPriceResource_GetPastPrices(t *testing.T) {
	helsinki, _ := time.LoadLocation("Europe/Helsinki")

	tests := []struct {
		name               string
		dateStr            string
		now                time.Time
		repoReturn         []model.PriceHistoryEntry
		repoError          error
		expectedStatus     int
		expectedCache      string
		expectedExpires    bool
		expectedPriceCount int
	}{
		{
			name:    "Success - Future Prices Available (Long Cache)",
			dateStr: "2023-10-27",
			now:     time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC),
			repoReturn: []model.PriceHistoryEntry{
				{
					Price:         decimal.NewFromFloat(10.0),
					DeliveryStart: time.Date(2023, 10, 28, 0, 0, 0, 0, helsinki), // Next day
				},
			},
			repoError:          nil,
			expectedStatus:     http.StatusOK,
			expectedCache:      utils.CACHE_LONG,
			expectedExpires:    false,
			expectedPriceCount: 1,
		},
		{
			name:    "Success - No Future Prices Yet - Before 11:57 (Expires Header)",
			dateStr: "2023-10-27",
			now:     time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC),
			repoReturn: []model.PriceHistoryEntry{
				{
					Price:         decimal.NewFromFloat(10.0),
					DeliveryStart: time.Date(2023, 10, 27, 23, 0, 0, 0, helsinki), // Same day
				},
			},
			repoError:          nil,
			expectedStatus:     http.StatusOK,
			expectedCache:      utils.CACHE_VAR,
			expectedExpires:    true,
			expectedPriceCount: 1,
		},
		{
			name:    "Success - No Future Prices Yet - After 11:57 (Short Cache)",
			dateStr: "2023-10-27",
			now:     time.Date(2023, 10, 27, 12, 0, 0, 0, time.UTC),
			repoReturn: []model.PriceHistoryEntry{
				{
					Price:         decimal.NewFromFloat(10.0),
					DeliveryStart: time.Date(2023, 10, 27, 23, 0, 0, 0, helsinki),
				},
			},
			repoError:          nil,
			expectedStatus:     http.StatusOK,
			expectedCache:      utils.CACHE_VAR + ", max-age=60",
			expectedExpires:    false,
			expectedPriceCount: 1,
		},
		{
			name:               "Invalid Date Format",
			dateStr:            "27-10-2023",
			now:                time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC),
			repoReturn:         nil,
			repoError:          nil,
			expectedStatus:     http.StatusBadRequest,
			expectedCache:      "",
			expectedExpires:    false,
			expectedPriceCount: 0,
		},
		{
			name:               "Repository Error",
			dateStr:            "2023-10-27",
			now:                time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC),
			repoReturn:         nil,
			repoError:          context.DeadlineExceeded,
			expectedStatus:     http.StatusInternalServerError,
			expectedCache:      "",
			expectedExpires:    false,
			expectedPriceCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPriceRepository)
			mockTime := new(MockTimeProvider)
			res := NewPriceResource(mockRepo, nil, mockTime, "secret")

			req := httptest.NewRequest("GET", "/api/prices/"+tt.dateStr, nil)
			req.SetPathValue("date", tt.dateStr)

			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusInternalServerError {
				if tt.repoError != nil || tt.repoReturn != nil {
					date, _ := time.Parse("2006-01-02", tt.dateStr)
					dateWithTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, helsinki)
					from := dateWithTime.AddDate(0, 0, -1)
					to := dateWithTime.AddDate(0, 0, 3)
					mockRepo.On("GetPrices", mock.Anything, from, to).Return(tt.repoReturn, tt.repoError)
				}
				if tt.repoError == nil && tt.repoReturn != nil {
					mockTime.On("Now").Return(tt.now)
				}
			}

			rr := httptest.NewRecorder()
			res.GetPastPrices(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
				assert.Equal(t, tt.expectedCache, rr.Header().Get(utils.CACHE_CONTROL_HEADER))
				if tt.expectedExpires {
					assert.NotEmpty(t, rr.Header().Get(utils.EXPIRES_HEADER))
				} else {
					assert.Empty(t, rr.Header().Get(utils.EXPIRES_HEADER))
				}

				var prices []model.PriceHistoryEntry
				err := json.NewDecoder(rr.Body).Decode(&prices)
				assert.NoError(t, err)
				assert.Len(t, prices, tt.expectedPriceCount)
			}
		})
	}
}

func TestPriceResource_UpdatePrices(t *testing.T) {
	password := "secret"
	mockRepo := new(MockPriceRepository)
	mockClient := new(MockNordPoolClient)
	mockTime := new(MockTimeProvider)
	pricesService := service.NewPricesService(mockClient, mockTime)
	res := NewPriceResource(mockRepo, pricesService, mockTime, password)

	now := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
	tomorrow := now.AddDate(0, 0, 1)

	t.Run("Wrong Password", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/update-prices?p=wrong", nil)
		rr := httptest.NewRecorder()
		res.UpdatePrices(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockTime.On("Now").Return(now).Once()
		nordPoolResp := &nordpool.PriceDataResponse{
			Market:   "DayAhead",
			Currency: "EUR",
			AreaStates: []nordpool.AreaState{
				{State: "Final", Areas: []string{"FI"}},
			},
			MultiAreaEntries: []nordpool.MultiAreaEntry{
				{
					DeliveryStart: tomorrow,
					DeliveryEnd:   tomorrow.Add(time.Hour),
					EntryPerArea:  map[string]float64{"FI": 10.5},
				},
			},
		}
		mockClient.On("GetDayAheadPrices", tomorrow, "DayAhead", "FI", "EUR").Return(nordPoolResp, nil).Once()
		mockRepo.On("InsertPrices", mock.Anything, mock.Anything).Return(int64(1), nil).Once()

		req := httptest.NewRequest("GET", "/api/update-prices?p="+password, nil)
		rr := httptest.NewRecorder()
		res.UpdatePrices(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp UpdatePricesResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.True(t, resp.Done)
	})

	t.Run("No Prices", func(t *testing.T) {
		mockTime.On("Now").Return(now).Once()
		mockClient.On("GetDayAheadPrices", tomorrow, "DayAhead", "FI", "EUR").Return(nil, nil).Once()

		req := httptest.NewRequest("GET", "/api/update-prices?p="+password, nil)
		rr := httptest.NewRecorder()
		res.UpdatePrices(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp UpdatePricesResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.False(t, resp.Done)
	})
}

func TestPriceResource_UpdatePricesForDate(t *testing.T) {
	password := "secret"
	mockRepo := new(MockPriceRepository)
	mockClient := new(MockNordPoolClient)
	mockTime := new(MockTimeProvider)
	pricesService := service.NewPricesService(mockClient, mockTime)
	res := NewPriceResource(mockRepo, pricesService, mockTime, password)

	dateStr := "2023-10-30"
	date, _ := time.Parse("2006-01-02", dateStr)

	t.Run("Success", func(t *testing.T) {
		nordPoolResp := &nordpool.PriceDataResponse{
			Market:   "DayAhead",
			Currency: "EUR",
			AreaStates: []nordpool.AreaState{
				{State: "Final", Areas: []string{"FI"}},
			},
			MultiAreaEntries: []nordpool.MultiAreaEntry{
				{
					DeliveryStart: date,
					DeliveryEnd:   date.Add(time.Hour),
					EntryPerArea:  map[string]float64{"FI": 12.3},
				},
			},
		}
		mockClient.On("GetDayAheadPrices", date, "DayAhead", "FI", "EUR").Return(nordPoolResp, nil).Once()
		mockRepo.On("InsertPrices", mock.Anything, mock.Anything).Return(int64(1), nil).Once()

		req := httptest.NewRequest("GET", "/api/update-prices/"+dateStr+"?p="+password, nil)
		req.SetPathValue("date", dateStr)
		rr := httptest.NewRecorder()
		res.UpdatePricesForDate(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var resp UpdatePricesResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.True(t, resp.Done)
	})
}

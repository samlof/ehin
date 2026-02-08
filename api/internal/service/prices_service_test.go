package service

import (
	"errors"
	"testing"
	"time"

	"github.com/samlof/ehin/internal/nordpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

// MockTimeProvider is a mock implementation of TimeProvider
type MockTimeProvider struct {
	mock.Mock
}

func (m *MockTimeProvider) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func TestPricesService_GetPrices(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  *nordpool.PriceDataResponse
		mockError     error
		expectedResp  *nordpool.PriceDataResponse
		expectedError error
	}{
		{
			name: "Success",
			mockResponse: &nordpool.PriceDataResponse{
				Market:     "DayAhead",
				Currency:   "EUR",
				AreaStates: []nordpool.AreaState{{State: "Final", Areas: []string{"FI"}}},
			},
			mockError: nil,
			expectedResp: &nordpool.PriceDataResponse{
				Market:     "DayAhead",
				Currency:   "EUR",
				AreaStates: []nordpool.AreaState{{State: "Final", Areas: []string{"FI"}}},
			},
			expectedError: nil,
		},
		{
			name:          "Client Error",
			mockResponse:  nil,
			mockError:     errors.New("client error"),
			expectedResp:  nil,
			expectedError: errors.New("client error"),
		},
		{
			name:          "Validation Error - Nil Response",
			mockResponse:  nil,
			mockError:     nil,
			expectedResp:  nil, // In implementation, invalidPrices returns true -> returns nil, nil
			expectedError: nil,
		},
		{
			name: "Validation Error - Wrong Market",
			mockResponse: &nordpool.PriceDataResponse{
				Market:     "Intraday",
				Currency:   "EUR",
				AreaStates: []nordpool.AreaState{{State: "Final", Areas: []string{"FI"}}},
			},
			mockError:     nil,
			expectedResp:  nil,
			expectedError: nil,
		},
		{
			name: "Validation Error - Wrong Currency",
			mockResponse: &nordpool.PriceDataResponse{
				Market:     "DayAhead",
				Currency:   "SEK",
				AreaStates: []nordpool.AreaState{{State: "Final", Areas: []string{"FI"}}},
			},
			mockError:     nil,
			expectedResp:  nil,
			expectedError: nil,
		},
		{
			name: "Validation Error - Empty AreaStates",
			mockResponse: &nordpool.PriceDataResponse{
				Market:     "DayAhead",
				Currency:   "EUR",
				AreaStates: []nordpool.AreaState{},
			},
			mockError:     nil,
			expectedResp:  nil,
			expectedError: nil,
		},
		{
			name: "Validation Error - Missing FI Area",
			mockResponse: &nordpool.PriceDataResponse{
				Market:     "DayAhead",
				Currency:   "EUR",
				AreaStates: []nordpool.AreaState{{State: "Final", Areas: []string{"SE"}}},
			},
			mockError:     nil,
			expectedResp:  nil,
			expectedError: nil,
		},
		{
			name: "Validation Error - FI State Not Final",
			mockResponse: &nordpool.PriceDataResponse{
				Market:     "DayAhead",
				Currency:   "EUR",
				AreaStates: []nordpool.AreaState{{State: "Preliminary", Areas: []string{"FI"}}},
			},
			mockError:     nil,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockNordPoolClient)
			mockTime := new(MockTimeProvider) // Not used in GetPrices, but required for NewPricesService
			service := NewPricesService(mockClient, mockTime)

			// GetPrices expects a specific date
			date := time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC)

			// Expectation setup
			// If mockResponse is nil and mockError is nil (nil response from client success = invalidPrices check inside service),
			// The implementation calls GetDayAheadPrices first.
			// The test cases cover returning nil from client (checked via mockResponse == nil && mockError == nil special handling in mock?)
			// Wait, my mock implementation:
			// if args.Get(0) == nil { return nil, args.Error(1) }
			// So if I pass nil as mockResponse, it returns nil.

			// Special case for "Validation Error - Nil Response":
			// We want client to return (nil, nil).
			// If mockResponse is nil, mock returns (nil, mockError).
			// So for that case, mockError needs to be nil.

			mockClient.On("GetDayAheadPrices", date, "DayAhead", "FI", "EUR").Return(tt.mockResponse, tt.mockError)

			resp, err := service.GetPrices(date)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResp, resp)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestPricesService_GetTomorrowsPrices(t *testing.T) {
	mockClient := new(MockNordPoolClient)
	mockTime := new(MockTimeProvider)
	service := NewPricesService(mockClient, mockTime)

	now := time.Date(2023, 10, 26, 10, 0, 0, 0, time.UTC)
	tomorrow := now.AddDate(0, 0, 1)

	mockTime.On("Now").Return(now)

	expectedResp := &nordpool.PriceDataResponse{
		Market:     "DayAhead",
		Currency:   "EUR",
		AreaStates: []nordpool.AreaState{{State: "Final", Areas: []string{"FI"}}},
	}

	mockClient.On("GetDayAheadPrices", tomorrow, "DayAhead", "FI", "EUR").Return(expectedResp, nil)

	resp, err := service.GetTomorrowsPrices()

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockTime.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

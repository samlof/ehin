package nordpool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NordPoolClient interface {
	GetDayAheadPrices(date time.Time, market, deliveryArea, currency string) (*PriceDataResponse, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) NordPoolClient {
	return &client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type MultiAreaEntry struct {
	DeliveryStart time.Time          `json:"deliveryStart"`
	DeliveryEnd   time.Time          `json:"deliveryEnd"`
	EntryPerArea  map[string]float64 `json:"entryPerArea"`
}

type AreaState struct {
	State string   `json:"state"`
	Areas []string `json:"areas"`
}

type PriceDataResponse struct {
	DeliveryDateCET  string           `json:"deliveryDateCET"`
	Version          int              `json:"version"`
	DeliveryAreas    []string         `json:"deliveryAreas"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	Market           string           `json:"market"`
	Currency         string           `json:"currency"`
	MultiAreaEntries []MultiAreaEntry `json:"multiAreaEntries"`
	AreaStates       []AreaState      `json:"areaStates"`
}

func (c *client) GetDayAheadPrices(date time.Time, market, deliveryArea, currency string) (*PriceDataResponse, error) {
	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("%s/api/DayAheadPrices?date=%s&market=%s&deliveryArea=%s&currency=%s",
		c.baseURL, dateStr, market, deliveryArea, currency)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices from NordPool: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NordPool API returned unexpected status: %s", resp.Status)
	}

	var priceResp PriceDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return nil, fmt.Errorf("failed to decode NordPool response: %w", err)
	}

	return &priceResp, nil
}

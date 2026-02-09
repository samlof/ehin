package nordpool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetDayAheadPrices(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request parameters
		expectedPath := "/api/DayAheadPrices"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("date") != "2025-09-30" {
			t.Errorf("expected date 2025-09-30, got %s", query.Get("date"))
		}
		if query.Get("market") != "DayAhead" {
			t.Errorf("expected market DayAhead, got %s", query.Get("market"))
		}
		if query.Get("deliveryArea") != "FI" {
			t.Errorf("expected deliveryArea FI, got %s", query.Get("deliveryArea"))
		}
		if query.Get("currency") != "EUR" {
			t.Errorf("expected currency EUR, got %s", query.Get("currency"))
		}

		// Mock response
		resp := PriceDataResponse{
			DeliveryDateCET: "2025-09-30",
			Version:         1,
			DeliveryAreas:   []string{"FI"},
			UpdatedAt:       time.Now(),
			Market:          "DayAhead",
			Currency:        "EUR",
			MultiAreaEntries: []MultiAreaEntry{
				{
					DeliveryStart: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
					DeliveryEnd:   time.Date(2025, 9, 30, 1, 0, 0, 0, time.UTC),
					EntryPerArea: map[string]float64{
						"FI": 10.5,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	date, _ := time.Parse("2006-01-02", "2025-09-30")
	prices, err := client.GetDayAheadPrices(date, "DayAhead", "FI", "EUR")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prices.DeliveryDateCET != "2025-09-30" {
		t.Errorf("expected delivery date 2025-09-30, got %s", prices.DeliveryDateCET)
	}

	if len(prices.MultiAreaEntries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(prices.MultiAreaEntries))
	}

	if prices.MultiAreaEntries[0].EntryPerArea["FI"] != 10.5 {
		t.Errorf("expected price 10.5, got %v", prices.MultiAreaEntries[0].EntryPerArea["FI"])
	}
}

func TestGetDayAheadPrices_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad Request")
	}))
	defer server.Close()

	client := NewClient(server.URL)
	date := time.Now()
	_, err := client.GetDayAheadPrices(date, "DayAhead", "FI", "EUR")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetDayAheadPrices_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{invalid json}")
	}))
	defer server.Close()

	client := NewClient(server.URL)
	date := time.Now()
	_, err := client.GetDayAheadPrices(date, "DayAhead", "FI", "EUR")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

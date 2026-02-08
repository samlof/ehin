package service

import (
	"log"
	"time"

	"github.com/samlof/ehin/internal/db/model"
	"github.com/samlof/ehin/internal/nordpool"
	"github.com/shopspring/decimal"
)

type PricesService struct {
	nordPoolClient nordpool.NordPoolClient
	timeProvider   TimeProvider
}

func NewPricesService(nordPoolClient nordpool.NordPoolClient, timeProvider TimeProvider) *PricesService {
	return &PricesService{
		nordPoolClient: nordPoolClient,
		timeProvider:   timeProvider,
	}
}

func (s *PricesService) GetTomorrowsPrices() (*nordpool.PriceDataResponse, error) {
	tomorrow := s.timeProvider.Now().AddDate(0, 0, 1)
	return s.GetPrices(tomorrow)
}

func (s *PricesService) GetPrices(date time.Time) (*nordpool.PriceDataResponse, error) {
	prices, err := s.nordPoolClient.GetDayAheadPrices(date, "DayAhead", "FI", "EUR")
	if err != nil {
		return nil, err
	}

	if s.invalidPrices(prices) {
		return nil, nil
	}

	return prices, nil
}

func (s *PricesService) ToPriceHistoryEntries(prices *nordpool.PriceDataResponse) []model.PriceHistoryEntry {
	if prices == nil {
		return nil
	}

	entries := make([]model.PriceHistoryEntry, 0, len(prices.MultiAreaEntries))
	for _, entry := range prices.MultiAreaEntries {
		if fiPrice, ok := entry.EntryPerArea["FI"]; ok {
			entries = append(entries, model.PriceHistoryEntry{
				Price:         decimal.NewFromFloat(fiPrice),
				DeliveryStart: entry.DeliveryStart,
				DeliveryEnd:   entry.DeliveryEnd,
			})
		}
	}
	return entries
}

func (s *PricesService) invalidPrices(prices *nordpool.PriceDataResponse) bool {
	if prices == nil {
		log.Println("Expected to find prices but was null")
		return true
	}
	if prices.Market != "DayAhead" {
		log.Printf("Expected market DayAhead but got %s", prices.Market)
		return true
	}
	if prices.Currency != "EUR" {
		log.Printf("Expected currency EUR but got %s", prices.Currency)
		return true
	}
	if len(prices.AreaStates) == 0 {
		log.Println("Expected areaStates to not be empty")
		return true
	}

	var fiState *nordpool.AreaState
	for _, state := range prices.AreaStates {
		for _, area := range state.Areas {
			if area == "FI" {
				fiState = &state
				break
			}
		}
		if fiState != nil {
			break
		}
	}

	if fiState == nil {
		log.Println("Couldn't find FI area from area states")
		return true
	}
	if fiState.State != "Final" {
		log.Printf("Expected state Final but got %s", fiState.State)
		return true
	}

	return false
}

package model

import (
	"time"
)

// PriceHistoryEntry represents a price entry in the database and API.
type PriceHistoryEntry struct {
	Price         float64   `json:"p"`
	DeliveryStart time.Time `json:"s"`
	DeliveryEnd   time.Time `json:"e"`
}

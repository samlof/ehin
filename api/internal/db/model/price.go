package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// PriceHistoryEntry represents a price entry in the database and API.
type PriceHistoryEntry struct {
	Price         decimal.Decimal `json:"p"`
	DeliveryStart time.Time       `json:"s"`
	DeliveryEnd   time.Time       `json:"e"`
}

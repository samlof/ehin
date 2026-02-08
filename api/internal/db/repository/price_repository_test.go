package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/samlof/ehin/internal/db/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPriceRepository_Select1(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	r := NewPriceRepository(mock)

	mock.ExpectQuery("SELECT 1").WillReturnRows(pgxmock.NewRows([]string{"1"}).AddRow(1))

	err = r.Select1(context.Background())
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPriceRepository_GetPrices(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	r := NewPriceRepository(mock)

	from := time.Now()
	to := from.Add(24 * time.Hour)

	rows := pgxmock.NewRows([]string{"price", "delivery_start", "delivery_end"}).
		AddRow(decimal.NewFromFloat(10.5), from, from.Add(time.Hour)).
		AddRow(decimal.NewFromFloat(12.0), from.Add(time.Hour), from.Add(2*time.Hour))

	mock.ExpectQuery("SELECT price, delivery_start, delivery_end FROM price_history").
		WithArgs(from, to).
		WillReturnRows(rows)

	entries, err := r.GetPrices(context.Background(), from, to)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.True(t, entries[0].Price.Equal(decimal.NewFromFloat(10.5)))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPriceRepository_InsertPrices(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	r := NewPriceRepository(mock)

	now := time.Now()
	entries := []model.PriceHistoryEntry{
		{
			Price:         decimal.NewFromFloat(10.5),
			DeliveryStart: now,
			DeliveryEnd:   now.Add(time.Hour),
		},
		{
			Price:         decimal.NewFromFloat(12.0),
			DeliveryStart: now.Add(time.Hour),
			DeliveryEnd:   now.Add(2 * time.Hour),
		},
	}

	mock.ExpectExec("INSERT INTO price_history").
		WithArgs(entries[0].DeliveryStart, entries[0].DeliveryEnd, entries[0].Price, entries[1].DeliveryStart, entries[1].DeliveryEnd, entries[1].Price).
		WillReturnResult(pgxmock.NewResult("INSERT", 2))

	affected, err := r.InsertPrices(context.Background(), entries)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), affected)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

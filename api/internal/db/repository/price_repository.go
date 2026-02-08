package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samlof/ehin/internal/db/model"
)

// PriceRepository defines the interface for price-related database operations.
type PriceRepository interface {
	Select1(ctx context.Context) error
	GetPrices(ctx context.Context, from, to time.Time) ([]model.PriceHistoryEntry, error)
	InsertPrices(ctx context.Context, entries []model.PriceHistoryEntry) (int64, error)
}

// DB defines the interface for database operations, compatible with pgxpool.Pool.
type DB interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type pgPriceRepository struct {
	db DB
}

// NewPriceRepository creates a new PostgreSQL-backed PriceRepository.
func NewPriceRepository(db DB) PriceRepository {
	return &pgPriceRepository{db: db}
}

// Select1 performs a simple health check query.
func (r *pgPriceRepository) Select1(ctx context.Context) error {
	var n int
	err := r.db.QueryRow(ctx, "SELECT 1").Scan(&n)
	return err
}

// GetPrices retrieves prices within the specified time range.
func (r *pgPriceRepository) GetPrices(ctx context.Context, from, to time.Time) ([]model.PriceHistoryEntry, error) {
	query := `
		SELECT price, delivery_start, delivery_end 
		FROM price_history 
		WHERE delivery_start >= $1 AND delivery_start < $2 
		ORDER BY delivery_start
	`
	rows, err := r.db.Query(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query prices: %w", err)
	}
	defer rows.Close()

	var entries []model.PriceHistoryEntry
	for rows.Next() {
		var entry model.PriceHistoryEntry
		if err := rows.Scan(&entry.Price, &entry.DeliveryStart, &entry.DeliveryEnd); err != nil {
			return nil, fmt.Errorf("failed to scan price entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return entries, nil
}

// InsertPrices batch inserts price entries with ON CONFLICT DO NOTHING.
func (r *pgPriceRepository) InsertPrices(ctx context.Context, entries []model.PriceHistoryEntry) (int64, error) {
	if len(entries) == 0 {
		return 0, nil
	}

	// Build dynamic INSERT with multiple VALUES
	// SQL: INSERT INTO price_history (delivery_start, delivery_end, price) VALUES ($1, $2, $3), ($4, $5, $6) ... ON CONFLICT (delivery_start) DO NOTHING

	valueStrings := make([]string, 0, len(entries))
	valueArgs := make([]interface{}, 0, len(entries)*3)

	for i, entry := range entries {
		offset := i * 3
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", offset+1, offset+2, offset+3))
		valueArgs = append(valueArgs, entry.DeliveryStart, entry.DeliveryEnd, entry.Price)
	}

	query := fmt.Sprintf(
		"INSERT INTO price_history (delivery_start, delivery_end, price) VALUES %s ON CONFLICT (delivery_start) DO NOTHING",
		strings.Join(valueStrings, ", "),
	)

	cmdTag, err := r.db.Exec(ctx, query, valueArgs...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert prices: %w", err)
	}

	return cmdTag.RowsAffected(), nil
}

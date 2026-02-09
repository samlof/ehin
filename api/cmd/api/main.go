package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	_ "time/tzdata"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samlof/ehin/internal/api/middleware"
	"github.com/samlof/ehin/internal/api/resource"
	"github.com/samlof/ehin/internal/config"
	"github.com/samlof/ehin/internal/db/repository"
	"github.com/samlof/ehin/internal/nordpool"
	"github.com/samlof/ehin/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	var dbPool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		var err error
		dbPool, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			slog.Error("Unable to connect to database", "error", err)
			os.Exit(1)
		}
		defer dbPool.Close()

		// Test connection
		if err := dbPool.Ping(context.Background()); err != nil {
			slog.Error("Unable to ping database", "error", err)
			os.Exit(1)
		}
		slog.Info("Connected to database")
	} else {
		slog.Info("DATABASE_URL not set, running without database")
	}

	// Repository and Service initialization
	var priceRepo repository.PriceRepository
	if dbPool != nil {
		priceRepo = repository.NewPriceRepository(dbPool)
	}
	
	nordPoolClient := nordpool.NewClient(cfg.NordPoolBaseURL)
	dateService := service.NewDateService()
	pricesService := service.NewPricesService(nordPoolClient, dateService)

	// Resource initialization
	greetingResource := resource.NewGreetingResource()
	priceResource := resource.NewPriceResource(priceRepo, pricesService, dateService, cfg.UpdatePricesPassword)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "EHIN API (Go)")
	})

	mux.HandleFunc("GET /hello", greetingResource.Hello)
	mux.HandleFunc("GET /api/prices/{date}", priceResource.GetPastPrices)
	mux.HandleFunc("GET /api/update-prices", priceResource.UpdatePrices)
	mux.HandleFunc("GET /api/update-prices/{date}", priceResource.UpdatePricesForDate)

	handler := middleware.CORS(cfg)(mux)

	serverAddr := ":" + cfg.Port
	slog.Info("Starting server", "addr", serverAddr)
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
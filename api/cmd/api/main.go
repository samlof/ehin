package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
			log.Fatalf("Unable to connect to database: %v", err)
		}
		defer dbPool.Close()

		// Test connection
		if err := dbPool.Ping(context.Background()); err != nil {
			log.Fatalf("Unable to ping database: %v", err)
		}
		log.Println("Connected to database")
	} else {
		log.Println("DATABASE_URL not set, running without database")
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
		fmt.Fprintf(w, "EHIN API (Go)")
	})

	mux.HandleFunc("GET /hello", greetingResource.Hello)
	mux.HandleFunc("GET /api/prices/{date}", priceResource.GetPastPrices)
	mux.HandleFunc("GET /api/update-prices", priceResource.UpdatePrices)
	mux.HandleFunc("GET /api/update-prices/{date}", priceResource.UpdatePricesForDate)

	handler := middleware.CORS(cfg)(mux)

	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
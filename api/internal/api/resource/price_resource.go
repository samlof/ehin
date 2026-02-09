package resource

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/samlof/ehin/internal/db/repository"
	"github.com/samlof/ehin/internal/service"
	"github.com/samlof/ehin/internal/utils"
)

type PriceResource struct {
	priceRepository      repository.PriceRepository
	pricesService        *service.PricesService
	dateService          service.TimeProvider
	updatePricesPassword string
}

func NewPriceResource(
	priceRepository repository.PriceRepository,
	pricesService *service.PricesService,
	dateService service.TimeProvider,
	updatePricesPassword string,
) *PriceResource {
	return &PriceResource{
		priceRepository:      priceRepository,
		pricesService:        pricesService,
		dateService:          dateService,
		updatePricesPassword: updatePricesPassword,
	}
}

type UpdatePricesResponse struct {
	Done bool `json:"done"`
}

func (res *PriceResource) UpdatePrices(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("p")
	if res.updatePricesPassword == "" || res.updatePricesPassword != password {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	slog.Info("Updating prices", "date", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	prices, err := res.pricesService.GetTomorrowsPrices()
	if err != nil {
		slog.Error("Error fetching tomorrow's prices", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if prices == nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(UpdatePricesResponse{Done: false}); err != nil {
			slog.Error("Error encoding response", "error", err)
		}
		return
	}

	entries := res.pricesService.ToPriceHistoryEntries(prices)
	_, err = res.priceRepository.InsertPrices(r.Context(), entries)
	if err != nil {
		slog.Error("Error inserting prices", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(UpdatePricesResponse{Done: true}); err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

func (res *PriceResource) UpdatePricesForDate(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("p")
	if res.updatePricesPassword == "" || res.updatePricesPassword != password {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	dateStr := r.PathValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	slog.Info("Updating prices", "date", dateStr)
	prices, err := res.pricesService.GetPrices(date)
	if err != nil {
		slog.Error("Error fetching prices", "date", dateStr, "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if prices == nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(UpdatePricesResponse{Done: false}); err != nil {
			slog.Error("Error encoding response", "error", err)
		}
		return
	}

	entries := res.pricesService.ToPriceHistoryEntries(prices)
	_, err = res.priceRepository.InsertPrices(r.Context(), entries)
	if err != nil {
		slog.Error("Error inserting prices", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(UpdatePricesResponse{Done: true}); err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

func (res *PriceResource) GetPastPrices(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	slog.Info("Fetching prices", "date", dateStr)

	if res.priceRepository == nil {
		slog.Warn("Price repository not initialized")
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		slog.Error("Error loading Europe/Helsinki", "error", err)
		helsinki = time.UTC
	}

	// Helsinki 00:00:00 on the requested date
	dateWithTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, helsinki)

	// Range matching Java implementation: date-1 to date+3
	from := dateWithTime.AddDate(0, 0, -1)
	to := dateWithTime.AddDate(0, 0, 3)

	prices, err := res.priceRepository.GetPrices(r.Context(), from, to)
	if err != nil {
		slog.Error("Error fetching prices from repository", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cacheString := utils.CACHE_LONG
	var expiresValue string

	if len(prices) > 0 {
		// Java: newestPriceDate.isAfter(date) ? CACHE_LONG : (isAfter(11:57) ? CACHE_VAR+max-age=60 : EXPIRES+CACHE_VAR)
		lastPrice := prices[len(prices)-1]
		lastPriceDate := lastPrice.DeliveryStart.In(helsinki)
		lastPriceDateOnly := time.Date(lastPriceDate.Year(), lastPriceDate.Month(), lastPriceDate.Day(), 0, 0, 0, 0, helsinki)
		
		if !lastPriceDateOnly.After(dateWithTime) {
			pricesUpdateTime := time.Date(date.Year(), date.Month(), date.Day(), 11, 57, 0, 0, time.UTC)
			if res.dateService.Now().After(pricesUpdateTime) {
				cacheString = utils.CACHE_VAR + ", max-age=60"
			} else {
				expiresValue = utils.GetGmtStringForCache(pricesUpdateTime)
				cacheString = utils.CACHE_VAR
			}
		}
	} else {
		// No prices found, use same logic as if we don't have tomorrow's prices yet
		pricesUpdateTime := time.Date(date.Year(), date.Month(), date.Day(), 11, 57, 0, 0, time.UTC)
		if res.dateService.Now().After(pricesUpdateTime) {
			cacheString = utils.CACHE_VAR + ", max-age=60"
		} else {
			expiresValue = utils.GetGmtStringForCache(pricesUpdateTime)
			cacheString = utils.CACHE_VAR
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(utils.CACHE_CONTROL_HEADER, cacheString)
	if expiresValue != "" {
		w.Header().Set(utils.EXPIRES_HEADER, expiresValue)
	}

	if err := json.NewEncoder(w).Encode(prices); err != nil {
		slog.Error("Error encoding prices", "error", err)
	}
}

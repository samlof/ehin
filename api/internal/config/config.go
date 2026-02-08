package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	UpdatePricesPassword string
	CORSAllowedOrigins   []string
	NordPoolBaseURL      string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var origins []string
	if corsOrigins != "" {
		origins = strings.Split(corsOrigins, ",")
	} else {
		origins = []string{"http://127.0.0.1:5173", "https://ehin.fi", "https://www.ehin.fi"}
	}

	return &Config{
		Port:                 port,
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		UpdatePricesPassword: os.Getenv("UPDATE_PRICES_PASSWORD"),
		CORSAllowedOrigins:   origins,
		NordPoolBaseURL:      "https://dataportal-api.nordpoolgroup.com",
	}
}

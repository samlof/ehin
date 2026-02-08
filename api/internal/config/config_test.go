package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("PORT", "9000")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:4000")

	cfg := LoadConfig()

	if cfg.Port != "9000" {
		t.Errorf("Expected Port 9000, got %s", cfg.Port)
	}

	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Errorf("Expected 2 CORS origins, got %d", len(cfg.CORSAllowedOrigins))
	}

	if cfg.CORSAllowedOrigins[0] != "http://localhost:3000" {
		t.Errorf("Expected first CORS origin http://localhost:3000, got %s", cfg.CORSAllowedOrigins[0])
	}
}

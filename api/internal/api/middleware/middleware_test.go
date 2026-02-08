package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samlof/ehin/internal/config"
)

func TestCORS_AllowedOrigin(t *testing.T) {
	cfg := &config.Config{
		CORSAllowedOrigins: []string{"http://example.com"},
	}
	middleware := CORS(cfg)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin: http://example.com, got %s",
			rr.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	cfg := &config.Config{
		CORSAllowedOrigins: []string{"http://example.com"},
	}
	middleware := CORS(cfg)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://malicious.com")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected Access-Control-Allow-Origin to be empty for disallowed origin, got %s",
			rr.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORS_Preflight(t *testing.T) {
	cfg := &config.Config{
		CORSAllowedOrigins: []string{"http://example.com"},
	}
	middleware := CORS(cfg)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for preflight, got %d", rr.Code)
	}

	if rr.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin: http://example.com, got %s",
			rr.Header().Get("Access-Control-Allow-Origin"))
	}

	if rr.Header().Get("Access-Control-Max-Age") != "86400" {
		t.Errorf("Expected Access-Control-Max-Age: 86400, got %s",
			rr.Header().Get("Access-Control-Max-Age"))
	}
}

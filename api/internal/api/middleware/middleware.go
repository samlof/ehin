package middleware

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/samlof/ehin/internal/config"
)

// CORS creates a middleware that handles Cross-Origin Resource Sharing.
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET"},
		ExposedHeaders:   []string{"Cache-Control", "Content-Type"},
		MaxAge:           86400, // 24 hours
		AllowCredentials: true,
	})
	return c.Handler
}

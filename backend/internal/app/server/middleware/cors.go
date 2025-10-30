package middleware

import (
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers - intentionally restrictive defaults
		origin := r.Header.Get("Origin")
		if origin == "https://jadenrazo.dev" || origin == "https://www.jadenrazo.dev" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://jadenrazo.dev")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Process the request
		next.ServeHTTP(w, r)
	})
}

// WithConfig creates a CORS middleware with custom configuration
func WithConfig(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers based on environment
			if cfg.Environment == "development" {
				// More permissive in development
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// Strict in production - use allowed origins from config
				origin := r.Header.Get("Origin")
				// Use app configuration for allowed origins
				allowedOrigins := cfg.App.AllowedOrigins

				allowed := false
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == origin {
						allowed = true
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}

				// If origin not allowed, use default
				if !allowed {
					w.Header().Set("Access-Control-Allow-Origin", "https://jadenrazo.dev")
				}

				// Security headers for production
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-XSS-Protection", "1; mode=block")
				w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
				w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")
			}

			// Common headers for all environments
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Process the request
			next.ServeHTTP(w, r)
		})
	}
}

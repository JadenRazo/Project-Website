package server

import (
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	appMiddleware "github.com/JadenRazo/Project-Website/backend/internal/app/server/middleware"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/compression"
	"github.com/JadenRazo/Project-Website/backend/internal/common/health"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/ratelimit"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRouter sets up the HTTP routing for the application
func SetupRouter(
	cfg *config.Config,
	auth *auth.Auth,
	cache cache.Cache,
	healthChecker *health.Health,
	rateLimiter ratelimit.RateLimiter,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(logger.HTTPLogger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(appMiddleware.CORS)

	// Add compression middleware if enabled
	if cfg.Compression.Enabled {
		r.Use(compression.CompressionMiddleware(&cfg.Compression))
	}

	// Health check endpoints
	r.Group(func(r chi.Router) {
		r.Get("/health", healthChecker.HTTPHandler().ServeHTTP)
		r.Get("/health/live", healthChecker.LivenessHandler().ServeHTTP)
		r.Get("/health/ready", healthChecker.ReadinessHandler().ServeHTTP)
	})

	// Public routes - no auth required
	r.Group(func(r chi.Router) {
		// Apply rate limiting to public routes
		// Using middleware adapter to make it compatible with chi
		publicRateLimit := func(next http.Handler) http.Handler {
			return rateLimiter.Limit("public", next)
		}
		r.Use(publicRateLimit)

		// Auth routes
		r.Route("/api/auth", func(r chi.Router) {
			r.Post("/login", handleLogin)
			r.Post("/register", handleRegister)
			r.Post("/refresh", handleTokenRefresh)
		})

		// URL shortener public routes
		r.Get("/u/{code}", handleURLRedirect)
	})

	// Protected routes - require authentication
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		// Apply rate limiting to API routes with a different limit
		apiRateLimit := func(next http.Handler) http.Handler {
			return rateLimiter.Limit("api", next)
		}
		r.Use(apiRateLimit)

		// User routes
		r.Route("/api/user", func(r chi.Router) {
			r.Get("/", handleGetCurrentUser)
			r.Put("/", handleUpdateUser)
		})

		// Messaging API routes
		r.Route("/api/messaging", func(r chi.Router) {
			// Channel routes
			r.Route("/channels", func(r chi.Router) {
				r.Get("/", handleGetChannels)
				r.Post("/", handleCreateChannel)
				r.Get("/{id}", handleGetChannel)
				r.Put("/{id}", handleUpdateChannel)
				r.Delete("/{id}", handleDeleteChannel)

				// Message routes
				r.Route("/{channelId}/messages", func(r chi.Router) {
					r.Get("/", handleGetMessages)
					r.Post("/", handleSendMessage)
					r.Get("/{id}", handleGetMessage)
					r.Put("/{id}", handleUpdateMessage)
					r.Delete("/{id}", handleDeleteMessage)
				})
			})
		})

		// URL shortener protected routes
		r.Route("/api/urls", func(r chi.Router) {
			r.Get("/", handleGetURLs)
			r.Post("/", handleCreateURL)
			r.Get("/{id}", handleGetURL)
			r.Delete("/{id}", handleDeleteURL)
			r.Get("/{id}/stats", handleGetURLStats)
		})
	})

	// Admin routes - require admin role
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(appMiddleware.RequireAdmin)

		r.Route("/api/admin", func(r chi.Router) {
			// Admin dashboard routes
			r.Get("/stats", handleAdminStats)
			r.Get("/users", handleGetAllUsers)
			r.Get("/users/{id}", handleGetUserById)
			r.Put("/users/{id}", handleUpdateUserById)
			r.Delete("/users/{id}", handleDeleteUserById)
		})
	})

	// Static file server for development
	if cfg.Environment == "development" {
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		r.Get("/swagger", serveSwaggerUI)
	}

	return r
}

// Handler stubs - these would be implemented elsewhere
func handleLogin(w http.ResponseWriter, r *http.Request)          {}
func handleRegister(w http.ResponseWriter, r *http.Request)       {}
func handleTokenRefresh(w http.ResponseWriter, r *http.Request)   {}
func handleURLRedirect(w http.ResponseWriter, r *http.Request)    {}
func handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {}
func handleUpdateUser(w http.ResponseWriter, r *http.Request)     {}
func handleGetChannels(w http.ResponseWriter, r *http.Request)    {}
func handleCreateChannel(w http.ResponseWriter, r *http.Request)  {}
func handleGetChannel(w http.ResponseWriter, r *http.Request)     {}
func handleUpdateChannel(w http.ResponseWriter, r *http.Request)  {}
func handleDeleteChannel(w http.ResponseWriter, r *http.Request)  {}
func handleGetMessages(w http.ResponseWriter, r *http.Request)    {}
func handleSendMessage(w http.ResponseWriter, r *http.Request)    {}
func handleGetMessage(w http.ResponseWriter, r *http.Request)     {}
func handleUpdateMessage(w http.ResponseWriter, r *http.Request)  {}
func handleDeleteMessage(w http.ResponseWriter, r *http.Request)  {}
func handleGetURLs(w http.ResponseWriter, r *http.Request)        {}
func handleCreateURL(w http.ResponseWriter, r *http.Request)      {}
func handleGetURL(w http.ResponseWriter, r *http.Request)         {}
func handleDeleteURL(w http.ResponseWriter, r *http.Request)      {}
func handleGetURLStats(w http.ResponseWriter, r *http.Request)    {}
func handleAdminStats(w http.ResponseWriter, r *http.Request)     {}
func handleGetAllUsers(w http.ResponseWriter, r *http.Request)    {}
func handleGetUserById(w http.ResponseWriter, r *http.Request)    {}
func handleUpdateUserById(w http.ResponseWriter, r *http.Request) {}
func handleDeleteUserById(w http.ResponseWriter, r *http.Request) {}
func serveSwaggerUI(w http.ResponseWriter, r *http.Request)       {}

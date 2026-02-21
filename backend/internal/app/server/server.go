package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"github.com/JadenRazo/Project-Website/backend/internal/common/health"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/JadenRazo/Project-Website/backend/internal/common/ratelimit"
	mcstatshttp "github.com/JadenRazo/Project-Website/backend/internal/mcstats/delivery/http"
)

// Server represents the HTTP server
type Server struct {
	server          *http.Server
	router          http.Handler
	config          *config.ServerConfig
	logger          logger.Logger
	shutdownTimeout time.Duration
}

// New creates a new server with all dependencies wired up
func New(
	cfg *config.Config,
	db database.Database,
	cacheClient cache.Cache,
	authService *auth.Auth,
	metricsManager *metrics.Manager,
	healthChecker *health.Health,
	mcStatsHandler *mcstatshttp.Handler,
) (*Server, error) {
	// Create rate limiter
	rateLimiter, err := ratelimit.NewRateLimiter(&cfg.RateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limiter: %w", err)
	}

	// Setup router with all handlers
	router := SetupRouter(cfg, authService, cacheClient, healthChecker, rateLimiter, mcStatsHandler)

	return &Server{
		router:          router,
		config:          &cfg.Server,
		logger:          nil, // Use package-level logger functions
		shutdownTimeout: cfg.Server.ShutdownTimeout,
	}, nil
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.ServerConfig, router http.Handler, log logger.Logger) *Server {
	return &Server{
		router:          router,
		config:          cfg,
		logger:          log,
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create the HTTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
		ReadHeaderTimeout: 5 * time.Second, // Security: Prevent Slowloris attacks
		MaxHeaderBytes:    s.config.MaxHeaderBytes,
	}

	// Channel to listen for server errors
	errChan := make(chan error, 1)

	// Start the server in a separate goroutine
	go func() {
		logger.Infof("Starting server on %s", addr)

		if s.config.TLSEnabled {
			logger.Info("TLS enabled, using HTTPS")
			errChan <- s.server.ListenAndServeTLS(s.config.TLSCert, s.config.TLSKey)
		} else {
			logger.Info("TLS disabled, using HTTP")
			errChan <- s.server.ListenAndServe()
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal or server error
	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		logger.Infof("Received signal: %v, shutting down server gracefully", sig)
		return s.Stop()
	}
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	logger.Infof("Shutting down server with %v timeout", s.shutdownTimeout)

	// Gracefully shut down the server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

// Shutdown gracefully stops the HTTP server with the provided context
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	logger.Info("Shutting down server")

	// Gracefully shut down the server
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Info("Server stopped gracefully")
	return nil
}

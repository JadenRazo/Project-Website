package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/app/server"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"github.com/JadenRazo/Project-Website/backend/internal/common/health"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/JadenRazo/Project-Website/backend/internal/common/storage"
	"github.com/JadenRazo/Project-Website/backend/internal/common/tracing"
)

// Application represents the main application structure
type Application struct {
	Config         *config.Config
	Server         *server.Server
	DB             *database.DB
	Cache          cache.Cache
	MetricsManager *metrics.Manager
	Tracer         *tracing.Tracer
	Storage        *storage.Provider
	Auth           *auth.Auth
	HealthChecker  *health.Health
	ShutdownCh     chan os.Signal
}

// New creates a new application instance
func New(cfg *config.Config) (*Application, error) {
	// Initialize logger
	if err := logger.Setup(cfg.Logger); err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	// Create app with default initialized shutdown channel
	app := &Application{
		Config:     cfg,
		ShutdownCh: make(chan os.Signal, 1),
	}

	// Setup signal handling
	signal.Notify(app.ShutdownCh, os.Interrupt, syscall.SIGTERM)

	return app, nil
}

// Boot initializes all application components
func (a *Application) Boot() error {
	logger.Info("Booting application...")
	startTime := time.Now()

	// Initialize database
	db, err := database.New(a.Config.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	a.DB = db
	logger.Info("Database initialized")

	// Initialize cache
	a.Cache, err = cache.NewCache(a.Config.Cache)
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}
	logger.Info("Cache initialized")

	// Initialize metrics
	metricsManager, err := metrics.NewManager(a.Config.Metrics)
	if err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}
	a.MetricsManager = metricsManager
	logger.Info("Metrics initialized")

	// Initialize tracing
	tracer, err := tracing.NewTracer(a.Config.Tracing)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}
	a.Tracer = tracer
	logger.Info("Tracing initialized")

	// Initialize storage
	storageProvider, err := storage.NewProvider(a.Config.Storage)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	a.Storage = storageProvider
	logger.Info("Storage initialized")

	// Initialize authentication
	auth, err := auth.New(a.Config.Auth, a.DB, a.Cache)
	if err != nil {
		return fmt.Errorf("failed to initialize auth: %w", err)
	}
	a.Auth = auth
	logger.Info("Authentication initialized")

	// Initialize health checker
	healthChecker := health.New()
	healthChecker.AddCheck("database", a.DB.HealthCheck)
	healthChecker.AddCheck("cache", a.Cache.HealthCheck)
	a.HealthChecker = healthChecker
	logger.Info("Health checker initialized")

	// Initialize HTTP server
	srv, err := server.New(a.Config, a.DB, a.Cache, a.Auth, a.MetricsManager, a.HealthChecker)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}
	a.Server = srv
	logger.Info("HTTP server initialized")

	// Report startup metrics
	a.MetricsManager.SetAppInfo(a.Config.Version, a.Config.GoVersion, a.Config.BuildDate)

	logger.Infof("Application booted in %v", time.Since(startTime))
	return nil
}

// Start runs the application
func (a *Application) Start() error {
	// Start HTTP server
	go func() {
		logger.Infof("Starting HTTP server on %s", a.Config.Server.Address)
		if err := a.Server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	<-a.ShutdownCh
	logger.Info("Shutdown signal received")

	return a.Shutdown()
}

// Shutdown gracefully stops all components
func (a *Application) Shutdown() error {
	logger.Info("Shutting down application...")
	startTime := time.Now()

	// Create a context with timeout for shutdown operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if a.Server != nil {
		if err := a.Server.Shutdown(ctx); err != nil {
			logger.Errorf("Error shutting down server: %v", err)
		} else {
			logger.Info("HTTP server shut down")
		}
	}

	// Shutdown tracer
	if a.Tracer != nil {
		if err := a.Tracer.Shutdown(ctx); err != nil {
			logger.Errorf("Error shutting down tracer: %v", err)
		} else {
			logger.Info("Tracer shut down")
		}
	}

	// Shutdown metrics
	if a.MetricsManager != nil {
		if err := a.MetricsManager.Shutdown(ctx); err != nil {
			logger.Errorf("Error shutting down metrics: %v", err)
		} else {
			logger.Info("Metrics shut down")
		}
	}

	// Close database connection
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			logger.Errorf("Error closing database: %v", err)
		} else {
			logger.Info("Database connection closed")
		}
	}

	// Close cache connection
	if a.Cache != nil {
		if err := a.Cache.Close(); err != nil {
			logger.Errorf("Error closing cache: %v", err)
		} else {
			logger.Info("Cache connection closed")
		}
	}

	logger.Infof("Application shut down in %v", time.Since(startTime))
	return nil
}

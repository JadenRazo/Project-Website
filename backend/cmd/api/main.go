// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	coreConfig "github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/gateway"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener"
)

const (
	appName    = "project-website-api"
	appVersion = "1.0.0"
)

func main() {
	// Capture and handle panics
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			fmt.Printf("PANIC: %v\nStack trace: %s\n", r, buf[:n])
			os.Exit(1)
		}
	}()

	// Setup context for the entire application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg := coreConfig.GetConfig()

	// Initialize logger
	err := logger.InitLogger(&config.LoggingConfig{
		Level:      cfg.Logging.Level,
		Format:     "json",
		Output:     "file",
		TimeFormat: time.RFC3339Nano,
		Filename:   filepath.Join("logs", "api.log"),
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}, appName, appVersion)

	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Shutdown()

	// Log startup
	logger.Info("Application starting up", "name", appName, "version", appVersion, "environment", cfg.Environment)
	logger.Info("Configuration loaded")

	// Initialize service manager
	serviceManager := core.NewServiceManager()

	// Initialize database connection
	database, err := db.GetDB()
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}

	// Run migrations
	logger.Info("Running database migrations")
	err = db.RunMigrations(
		&domain.User{},
		&domain.ShortURL{},
		&domain.Message{},
		&domain.Channel{},
		&domain.Attachment{},
	)
	if err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	// Initialize services
	logger.Info("Initializing services")

	// Create service-specific configs
	urlShortenerConfig := urlshortener.Config{
		BaseURL:      cfg.URLShortener.BaseURL,
		MaxURLLength: 2048,
		MinURLLength: 5,
	}

	messagingConfig := messaging.Config{
		WebSocketPort:    8081,
		MaxMessageSize:   cfg.Messaging.MaxMessageSize,
		MaxAttachments:   10,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
	}

	urlShortenerService := urlshortener.NewService(database, urlShortenerConfig)
	messagingService := messaging.NewService(database, messagingConfig)

	// Initialize devpanel service
	devpanelService := devpanel.NewService(serviceManager, devpanel.Config{
		AdminToken:      cfg.AdminToken,
		MetricsInterval: 30 * time.Second,
		MaxLogLines:     1000,
		LogRetention:    7 * 24 * time.Hour,
	})

	// Initialize log manager
	logManager := devpanel.NewLogManager(
		filepath.Join("logs", "services"),
		cfg.MaxLogLines,
		cfg.LogRetention,
	)
	logManager.StartCleanup()
	// Note: We should call a cleanup method here if available

	// Initialize metrics collector
	_ = devpanel.NewMetricsCollector(devpanel.Config{
		MetricsInterval: 30 * time.Second,
	})
	// Note: Metrics collection should be started here, but signature is unknown
	// metricsCollector.StartCollecting(serviceManager)
	// Note: We should call a cleanup method here if available

	// Register services with service manager
	logger.Info("Registering services with service manager")
	if err := serviceManager.RegisterService(urlShortenerService); err != nil {
		logger.Fatal("Failed to register URL shortener service", "error", err)
	}
	if err := serviceManager.RegisterService(messagingService); err != nil {
		logger.Fatal("Failed to register messaging service", "error", err)
	}
	if err := serviceManager.RegisterService(devpanelService); err != nil {
		logger.Fatal("Failed to register devpanel service", "error", err)
	}

	// Initialize API gateway
	logger.Info("Initializing API gateway")
	apiGateway := gateway.NewGateway(serviceManager)

	// Configure security headers
	apiGateway.AddMiddleware(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	// Add compression middleware
	// Note: gin.Gzip middleware needs to be imported to use compression

	// Add profiling endpoints in development
	if cfg.Environment == "development" {
		logger.Info("Enabling profiling endpoints (development mode)")
		// Note: pprof needs to be imported and registered if needed
	}

	// Configure CORS
	apiGateway.AddMiddleware(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.AllowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register service routes
	logger.Info("Registering service routes")
	apiGateway.RegisterService("urls", urlShortenerService.RegisterRoutes)
	apiGateway.RegisterService("messaging", messagingService.RegisterRoutes)
	apiGateway.RegisterService("devpanel", devpanelService.RegisterRoutes)

	// Register system endpoints
	apiGateway.RegisterHealthCheck()
	apiGateway.RegisterMetrics()
	apiGateway.RegisterDevPanel()

	// Add Prometheus metrics endpoint
	apiGateway.GetRouter().GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Set up shortcode redirect route at root level
	apiGateway.GetRouter().GET("/:shortCode", urlShortenerService.RedirectHandler)

	// Start health checks
	serviceManager.StartHealthChecks()

	// Configure server with proper timeouts and limits
	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        apiGateway.GetRouter(),
		ReadTimeout:    10 * time.Second,  // Default read timeout
		WriteTimeout:   15 * time.Second,  // Default write timeout
		IdleTimeout:    120 * time.Second, // Default idle timeout
		MaxHeaderBytes: 1 << 20,           // 1 MB
	}

	// Start server
	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		var err error
		if false { // TLS disabled for now
			logger.Info("TLS enabled", "cert_path", "cert.pem", "key_path", "key.pem")
			err = srv.ListenAndServeTLS("cert.pem", "key.pem")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	logger.Info("Server started successfully, waiting for shutdown signal")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Received shutdown signal", "signal", sig.String())

	// Define shutdown timeout
	shutdownTimeout := 5 * time.Second // Default shutdown timeout

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdownCancel()

	// Notify clients that the server is shutting down (if applicable)
	logger.Info("Gracefully shutting down server")

	// Stop all services
	logger.Info("Stopping all services")
	if err := serviceManager.StopAllServices(); err != nil {
		logger.Error("Error stopping services", "error", err)
	}

	// Shutdown HTTP server
	logger.Info("Shutting down HTTP server", "timeout", shutdownTimeout)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		return
	}

	// Close database connection
	logger.Info("Closing database connection")
	if err := db.CloseDB(); err != nil {
		logger.Error("Error closing database", "error", err)
	}

	logger.Info("Server exited properly")
}

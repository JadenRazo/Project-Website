// cmd/api/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	coreConfig "github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel"
	"github.com/JadenRazo/Project-Website/backend/internal/codestats"
	codeStatsHTTP "github.com/JadenRazo/Project-Website/backend/internal/codestats/delivery/http"
	// "github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	projectRepo "github.com/JadenRazo/Project-Website/backend/internal/projects/repository"
	projectHTTP "github.com/JadenRazo/Project-Website/backend/internal/projects/delivery/http"
	projectMemoryService "github.com/JadenRazo/Project-Website/backend/internal/projects/service"
	"github.com/JadenRazo/Project-Website/backend/internal/gateway"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging"
	"github.com/JadenRazo/Project-Website/backend/internal/status"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener"
)

const (
	appName    = "project-website-api"
	appVersion = "1.0.0"
)

// loadEnvFile loads environment variables from a file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return scanner.Err()
}

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

	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := loadEnvFile(".env"); err != nil {
			fmt.Printf("Warning: Failed to load .env file: %v\n", err)
		}
	}

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
	logger.Info("Running automatic migrations")
	
	err = db.RunMigrations(
		&projectRepo.ProjectModel{},
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
	
	// Initialize code stats service
	codeStatsConfig, err := codestats.LoadConfig("config/codestats.yaml")
	if err != nil {
		logger.Warn("Failed to load code stats config, using defaults", "error", err)
		codeStatsConfig = codestats.DefaultConfig()
	}
	codeStatsService := codestats.NewService(database, *codeStatsConfig)

	// Initialize projects service (using memory service for now due to DB issues)
	// projectRepository := projectRepo.NewGormRepository(database)
	// projectService := project.NewService(projectRepository)
	projectService := projectMemoryService.NewMemoryProjectService()

	// Initialize metrics manager
	metricsConfig := metrics.DefaultConfig()
	metricsManager, err := metrics.NewManager(metricsConfig)
	if err != nil {
		logger.Warn("Failed to initialize metrics manager", "error", err)
		metricsManager = nil
	}

	// Initialize status monitoring service with metrics
	statusService := status.NewService(database, 30*time.Second, metricsManager)

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
	metricsCollector := devpanel.NewMetricsCollector(devpanel.Config{
		MetricsInterval: 30 * time.Second,
	})
	_ = metricsCollector // Silencing "declared and not used" error for now.
	// Example: Start collecting metrics. Adapt to the actual method signature.
	// go metricsCollector.StartCollecting(ctx, serviceManager)
	// Note: We should call a cleanup method for metricsCollector during shutdown if it has one.

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
	apiGateway.AddMiddleware(gzip.Gzip(gzip.DefaultCompression))

	// Add Request ID middleware
	apiGateway.AddMiddleware(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("RequestID", requestID)       // Make it available for other handlers/loggers
		c.Header("X-Request-ID", requestID) // Ensure it's in the response
		c.Next()
	})

	// Add latency tracking middleware
	if metricsManager != nil {
		apiGateway.AddMiddleware(metricsManager.GinLatencyMiddleware())
	}

	// Add profiling endpoints in development or if explicitly enabled
	if cfg.Environment == "development" || cfg.EnablePprof { // Assuming cfg.EnablePprof exists
		logger.Info("Enabling profiling endpoints")
		pprofGroup := apiGateway.GetRouter().Group("/debug/pprof")
		{
			pprofGroup.GET("/", gin.WrapF(pprof.Index))
			pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
			pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
			pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
			pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
			pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
			pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		}
	}

	// Configure CORS
	var allowedOrigins []string
	if cfg.AllowedOrigins != "" {
		allowedOrigins = strings.Split(cfg.AllowedOrigins, ",")
	} else {
		// Default to allow all origins in development
		if cfg.Environment == "development" {
			allowedOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
		} else {
			allowedOrigins = []string{"https://jadenrazo.dev"}
		}
	}
	
	apiGateway.AddMiddleware(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
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
	
	// Register code stats routes
	codeStatsHandler := codeStatsHTTP.NewHandler(codeStatsService)
	apiGateway.RegisterService("code", codeStatsHandler.RegisterRoutes)
	
	// Register projects routes - this will create routes at /api/v1/projects/*
	projectHandler := projectHTTP.NewHandler(projectService)
	apiGateway.RegisterService("projects", projectHandler.RegisterRoutes)
	
	// Register status monitoring routes
	apiGateway.RegisterService("status", statusService.RegisterRoutes)

	// Register system endpoints
	apiGateway.RegisterHealthCheck()
	
	// Add Prometheus metrics endpoint (instead of RegisterMetrics to avoid duplicate)
	apiGateway.GetRouter().GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Serve static files
	router := apiGateway.GetRouter()
	
	// Serve frontend build files
	frontendBuildPath := filepath.Join("..", "frontend", "build")
	if _, err := os.Stat(frontendBuildPath); err == nil {
		logger.Info("Serving frontend build files", "path", frontendBuildPath)
		router.Static("/static", filepath.Join(frontendBuildPath, "static"))
		router.StaticFile("/manifest.json", filepath.Join(frontendBuildPath, "manifest.json"))
		router.StaticFile("/favicon.ico", filepath.Join(frontendBuildPath, "favicon.ico"))
		router.StaticFile("/robots.txt", filepath.Join(frontendBuildPath, "robots.txt"))
	}

	// Serve public files (including code_stats.json)
	publicPath := filepath.Join("..", "frontend", "public")
	if _, err := os.Stat(publicPath); err == nil {
		logger.Info("Serving public files", "path", publicPath)
		router.StaticFile("/code_stats.json", filepath.Join(publicPath, "code_stats.json"))
		router.StaticFile("/apple-touch-icon.png", filepath.Join(publicPath, "apple-touch-icon.png"))
		router.StaticFile("/favicon-16x16.png", filepath.Join(publicPath, "favicon-16x16.png"))
		router.StaticFile("/favicon-32x32.png", filepath.Join(publicPath, "favicon-32x32.png"))
	}

	// Set up shortcode redirect route at root level
	router.GET("/:shortCode", urlShortenerService.RedirectHandler)

	// SPA fallback - serve index.html for all unmatched routes
	router.NoRoute(func(c *gin.Context) {
		// Don't serve index.html for API routes
		if strings.HasPrefix(c.Request.URL.Path, "/api") ||
			strings.HasPrefix(c.Request.URL.Path, "/health") ||
			strings.HasPrefix(c.Request.URL.Path, "/metrics") ||
			strings.HasPrefix(c.Request.URL.Path, "/debug") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		
		indexPath := filepath.Join(frontendBuildPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		} else {
			c.JSON(404, gin.H{"error": "Frontend build not found"})
		}
	})

	// Start health checks
	serviceManager.StartHealthChecks()
	
	// Start code stats periodic update
	codeStatsService.StartPeriodicUpdate(1 * time.Hour)

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
		// Use existing fields from cfg.Server for TLS
		if cfg.Server.TLSEnabled {
			logger.Info("TLS enabled", "cert_path", cfg.Server.TLSCert, "key_path", cfg.Server.TLSKey)
			if _, errCert := os.Stat(cfg.Server.TLSCert); os.IsNotExist(errCert) {
				logger.Fatal("TLS cert file not found", "path", cfg.Server.TLSCert)
			}
			if _, errKey := os.Stat(cfg.Server.TLSKey); os.IsNotExist(errKey) {
				logger.Fatal("TLS key file not found", "path", cfg.Server.TLSKey)
			}
			err = srv.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey)
		} else {
			logger.Info("TLS disabled, starting HTTP server")
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

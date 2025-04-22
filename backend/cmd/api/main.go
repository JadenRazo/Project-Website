// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	
	"github.com/gin-contrib/cors"
	
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/gateway"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()
	
	// Initialize service manager
	serviceManager := core.NewServiceManager()
	
	// Initialize database connection
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Run migrations
	err = db.RunMigrations(
		&domain.User{},
		&domain.ShortURL{},
		&domain.Message{},
		&domain.Channel{},
		&domain.Attachment{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Initialize services
	urlShortenerService := urlshortener.NewService(database, cfg.URLShortener)
	messagingService := messaging.NewService(database, cfg.Messaging)
	
	// Initialize devpanel service
	devpanelService := devpanel.NewService(serviceManager, devpanel.Config{
		AdminToken:     cfg.AdminToken,
		MetricsInterval: 30 * time.Second,
		MaxLogLines:    1000,
		LogRetention:   7 * 24 * time.Hour,
	})
	
	// Initialize log manager
	logManager := devpanel.NewLogManager(
		filepath.Join("logs", "services"),
		cfg.MaxLogLines,
		cfg.LogRetention,
	)
	logManager.StartCleanup()
	
	// Initialize metrics collector
	metricsCollector := devpanel.NewMetricsCollector(devpanel.Config{
		MetricsInterval: 30 * time.Second,
	})
	metricsCollector.StartCollecting(serviceManager)
	
	// Register services with service manager
	if err := serviceManager.RegisterService(urlShortenerService); err != nil {
		log.Fatalf("Failed to register URL shortener service: %v", err)
	}
	if err := serviceManager.RegisterService(messagingService); err != nil {
		log.Fatalf("Failed to register messaging service: %v", err)
	}
	if err := serviceManager.RegisterService(devpanelService); err != nil {
		log.Fatalf("Failed to register devpanel service: %v", err)
	}
	
	// Initialize API gateway
	apiGateway := gateway.NewGateway(serviceManager)
	
	// Configure CORS
	apiGateway.AddMiddleware(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Register service routes
	apiGateway.RegisterService("urls", urlShortenerService.RegisterRoutes)
	apiGateway.RegisterService("messaging", messagingService.RegisterRoutes)
	apiGateway.RegisterService("devpanel", devpanelService.RegisterRoutes)
	
	// Register system endpoints
	apiGateway.RegisterHealthCheck()
	apiGateway.RegisterMetrics()
	apiGateway.RegisterDevPanel()
	
	// Set up shortcode redirect route at root level
	apiGateway.GetRouter().GET("/:shortCode", urlShortenerService.RedirectHandler)
	
	// Start health checks
	serviceManager.StartHealthChecks()
	
	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: apiGateway.GetRouter(),
	}
	
	// Graceful shutdown handling
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Stop all services
	if err := serviceManager.StopAllServices(); err != nil {
		log.Printf("Error stopping services: %v", err)
	}
	
	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited properly")
}
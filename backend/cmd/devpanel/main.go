package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	fmt.Println("Starting DevPanel service...")

	// Create service manager for devpanel to manage
	serviceManager := core.NewServiceManager()

	// Initialize devpanel service
	devpanelService := devpanel.NewService(serviceManager, devpanel.Config{
		AdminToken:      cfg.AdminToken,
		MetricsInterval: 30 * time.Second,
		MaxLogLines:     1000,
		LogRetention:    7 * 24 * time.Hour,
	})

	// Register service with manager
	if err := serviceManager.RegisterService(devpanelService); err != nil {
		log.Fatalf("Failed to register devpanel service: %v", err)
	}

	// Set up router
	router := gin.Default()

	// Create a router group for the service
	routerGroup := router.Group("/")

	// Register devpanel routes
	devpanelService.RegisterRoutes(routerGroup)

	// Hardcoded port for DevPanel service (8081 from status report)
	const devPanelPort = 8081

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", devPanelPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("DevPanel service listening on port %d\n", devPanelPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down DevPanel service...")

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("DevPanel service stopped")
}

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

	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging"
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

	fmt.Println("Starting Messaging service...")

	// Initialize database connection
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure messaging service
	messagingConfig := messaging.Config{
		WebSocketPort:    8082, // Match the expected port in the service report
		MaxMessageSize:   cfg.Messaging.MaxMessageSize,
		MaxAttachments:   10,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
	}

	// Initialize service
	messagingService := messaging.NewService(database, messagingConfig)

	// Set up router
	router := gin.Default()

	// Create a router group for the service
	routerGroup := router.Group("/")

	// Register messaging routes
	messagingService.RegisterRoutes(routerGroup)

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         ":8082", // Match the expected port in the service report
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Println("Messaging service listening on port 8082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down Messaging service...")

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Messaging service stopped")
}

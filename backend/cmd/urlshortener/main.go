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
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener"
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

	fmt.Println("Starting URL Shortener service...")

	// Initialize database connection
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure URL shortener service
	urlShortenerConfig := urlshortener.Config{
		BaseURL:      cfg.URLShortener.BaseURL,
		MaxURLLength: 2048,
		MinURLLength: 5,
	}

	// Initialize service
	urlShortenerService := urlshortener.NewService(database, urlShortenerConfig)

	// Set up router
	router := gin.Default()

	// Create a router group for the service
	routerGroup := router.Group("/")

	// Register URL shortener routes
	urlShortenerService.RegisterRoutes(routerGroup)

	// Add redirect handler at root level
	router.GET("/:shortCode", urlShortenerService.RedirectHandler)

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         ":8083", // Match the expected port in the service report
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Println("URL Shortener service listening on port 8083")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down URL Shortener service...")

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("URL Shortener service stopped")
}

// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	
	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener/handlers"
	urlService "github.com/JadenRazo/Project-Website/backend/internal/urlshortener/service"
	urlRepo "github.com/JadenRazo/Project-Website/backend/internal/urlshortener/repository"
	msgHandlers "github.com/JadenRazo/Project-Website/backend/internal/messaging/handlers"
	msgService "github.com/JadenRazo/Project-Website/backend/internal/messaging/service"
	msgRepo "github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()
	
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
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
	
	// Initialize repositories
	urlRepository := urlRepo.NewGormRepository(database)
	msgRepository := msgRepo.NewGormRepository(database)
	
	// Initialize services
	urlShortenerService := urlService.NewURLShortenerService(urlRepository, cfg.URLShortener)
	messagingService := msgService.NewMessagingService(msgRepository, cfg.Messaging)
	
	// Create Gin router
	router := gin.Default()
	
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Set up API routes
	apiV1 := router.Group("/api/v1")
	
	// Set up URL shortener routes
	urlGroup := apiV1.Group("/urls")
	urlHandlers := handlers.NewURLShortenerHandler(urlShortenerService)
	urlHandlers.RegisterRoutes(urlGroup)
	
	// Set up messaging routes
	msgGroup := apiV1.Group("/messaging")
	messageHandlers := msgHandlers.NewMessagingHandler(messagingService)
	messageHandlers.RegisterRoutes(msgGroup)
	
	// Set up shortcode redirect route at root level
	router.GET("/:shortCode", urlHandlers.RedirectHandler)
	
	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
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
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited properly")
}
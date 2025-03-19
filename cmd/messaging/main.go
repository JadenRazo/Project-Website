package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"jadenrazo.dev/internal/messaging/config"
	"jadenrazo.dev/internal/messaging/delivery/websocket"
	"jadenrazo.dev/internal/messaging/repository/sqlite"
	"jadenrazo.dev/internal/messaging/service"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config/messaging.yaml", "Path to configuration file")
	port := flag.String("port", "8082", "HTTP port")
	flag.Parse()

	// Override port from environment variable if set (for use with the developer panel)
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Warning: Failed to load configuration: %v, using defaults", err)
		cfg = config.DefaultConfig()
	}

	// Initialize repository
	repo, err := sqlite.NewRepository(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize connection manager for WebSockets
	connManager := websocket.NewConnectionManager(cfg.Websocket.MaxConnections)

	// Initialize WebSocket hub
	hub := websocket.NewHub(connManager)
	go hub.Run()

	// Initialize service
	messagingService := service.NewMessagingService(repo, hub)

	// Initialize router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Register WebSocket handler
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Register API routes
	r.Route("/api/messaging", func(r chi.Router) {
		// Channels API
		r.Route("/channels", func(r chi.Router) {
			r.Get("/", messagingService.ListChannels)
			r.Post("/", messagingService.CreateChannel)
			r.Get("/{channelID}", messagingService.GetChannel)
			r.Put("/{channelID}", messagingService.UpdateChannel)
			r.Delete("/{channelID}", messagingService.DeleteChannel)

			// Channel messages
			r.Get("/{channelID}/messages", messagingService.GetMessages)
			r.Post("/{channelID}/messages", messagingService.SendMessage)

			// Channel members
			r.Get("/{channelID}/members", messagingService.GetChannelMembers)
			r.Post("/{channelID}/members", messagingService.AddChannelMember)
			r.Delete("/{channelID}/members/{userID}", messagingService.RemoveChannelMember)
		})

		// User API
		r.Route("/users", func(r chi.Router) {
			r.Get("/me/channels", messagingService.GetUserChannels)
			r.Get("/me/status", messagingService.GetUserStatus)
			r.Put("/me/status", messagingService.UpdateUserStatus)
		})
	})

	// Add health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Messaging service is running"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting messaging service on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down messaging service...")
	log.Println("Messaging service stopped")
}

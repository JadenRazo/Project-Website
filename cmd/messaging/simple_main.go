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
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8082", "HTTP port")
	configPath := flag.String("config", "config/messaging.yaml", "Path to configuration file")
	flag.Parse()

	// Log that we're using the simplified implementation
	log.Printf("Starting simplified messaging service (config: %s)", *configPath)

	// Override port from environment variable if set (for use with the developer panel)
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

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

	// Register placeholder WebSocket handler
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("WebSocket endpoint (placeholder)"))
	})

	// Register API routes with placeholder handlers
	r.Route("/api/messaging", func(r chi.Router) {
		// Channels API
		r.Route("/channels", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"channels": []}`))
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"message": "Channel created (placeholder)", "id": 1}`))
			})

			// Other placeholder handlers
			r.Get("/{channelID}", func(w http.ResponseWriter, r *http.Request) {
				channelID := chi.URLParam(r, "channelID")
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(fmt.Sprintf(`{"id": %s, "name": "Placeholder Channel"}`, channelID)))
			})
		})
	})

	// Add health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Messaging service is running (simplified version)"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting simplified messaging service on %s", addr)
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

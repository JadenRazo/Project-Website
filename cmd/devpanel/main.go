package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"jadenrazo.dev/internal/devpanel/config"
	"jadenrazo.dev/internal/devpanel/repository"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config/devpanel.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Warning: Failed to load configuration: %v, using defaults", err)
		// Create a basic default config
		cfg = &config.Config{
			Server: config.ServerConfig{
				Host:     "0.0.0.0",
				Port:     8080,
				BasePath: "/devpanel",
			},
			Auth: config.AuthConfig{
				Enabled:     true,
				AdminUsers:  []string{"admin"},
				TokenExpiry: 24,
			},
			Projects: config.ProjectsConfig{
				ConfigPath: "config/projects.yaml",
			},
		}
	}

	// Create project repository
	projectRepo, err := repository.NewFileProjectRepository(cfg.Projects.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to create project repository: %v", err)
	}

	// Create and initialize the router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Set up API routes
	r.Route(cfg.Server.BasePath, func(r chi.Router) {
		// API routes
		r.Route("/api", func(r chi.Router) {
			// Projects API
			r.Route("/projects", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					projects, err := projectRepo.GetAll()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					respondWithJSON(w, http.StatusOK, projects)
				})

				// Other project endpoints would be implemented here
			})
		})

		// Serve static files for the UI
		fileServer := http.FileServer(http.Dir("./web/devpanel"))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != cfg.Server.BasePath && r.URL.Path != cfg.Server.BasePath+"/" {
				fileServer.ServeHTTP(w, r)
				return
			}
			http.ServeFile(w, r, "./web/devpanel/index.html")
		})
	})

	// Create server instance
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting developer panel server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	log.Println("Server stopped")
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

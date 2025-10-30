package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// Config contains configuration for the devpanel server
type Config struct {
	// Address is the server address to listen on
	Address string `json:"address" mapstructure:"address"`

	// BasePath is the base path for the devpanel API
	BasePath string `json:"base_path" mapstructure:"base_path"`

	// Timeout settings for the HTTP server
	ReadTimeout     time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout" mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout" mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" mapstructure:"shutdown_timeout"`

	// CORS configuration
	CORSAllowedOrigins   []string `json:"cors_allowed_origins" mapstructure:"cors_allowed_origins"`
	CORSAllowedMethods   []string `json:"cors_allowed_methods" mapstructure:"cors_allowed_methods"`
	CORSAllowedHeaders   []string `json:"cors_allowed_headers" mapstructure:"cors_allowed_headers"`
	CORSAllowCredentials bool     `json:"cors_allow_credentials" mapstructure:"cors_allow_credentials"`
	CORSMaxAge           int      `json:"cors_max_age" mapstructure:"cors_max_age"`
}

// DefaultConfig returns the default server configuration
func DefaultConfig() *Config {
	return &Config{
		Address:              ":8081",
		BasePath:             "/api/devpanel",
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		IdleTimeout:          120 * time.Second,
		ShutdownTimeout:      30 * time.Second,
		CORSAllowedOrigins:   []string{"*"},
		CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSAllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		CORSAllowCredentials: false,
		CORSMaxAge:           300,
	}
}

// Server represents the devpanel HTTP server
type Server struct {
	config            *Config
	router            chi.Router
	httpServer        *http.Server
	projectService    *project.Service
	repositoryService *repository.Service
}

// New creates a new devpanel server
func New(
	config *Config,
	projectService *project.Service,
	repositoryService *repository.Service,
) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}

	srv := &Server{
		config:            config,
		projectService:    projectService,
		repositoryService: repositoryService,
	}

	// Initialize router
	router := chi.NewRouter()

	// Basic middleware
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.CORSAllowedOrigins,
		AllowedMethods:   config.CORSAllowedMethods,
		AllowedHeaders:   config.CORSAllowedHeaders,
		AllowCredentials: config.CORSAllowCredentials,
		MaxAge:           config.CORSMaxAge,
	}))

	// Set up routes
	router.Route(config.BasePath, func(r chi.Router) {
		// Health check
		r.Get("/health", srv.handleHealth)

		// Projects
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", srv.handleListProjects)
			r.Post("/", srv.handleCreateProject)
			r.Get("/{projectID}", srv.handleGetProject)
			r.Put("/{projectID}", srv.handleUpdateProject)
			r.Delete("/{projectID}", srv.handleDeleteProject)
			r.Post("/{projectID}/archive", srv.handleArchiveProject)
			r.Post("/{projectID}/activate", srv.handleActivateProject)

			// Project repositories
			r.Route("/{projectID}/repositories", func(r chi.Router) {
				r.Get("/", srv.handleListProjectRepositories)
				r.Post("/", srv.handleCreateRepository)
			})
		})

		// Repositories
		r.Route("/repositories", func(r chi.Router) {
			r.Get("/", srv.handleListRepositories)
			r.Get("/{repoID}", srv.handleGetRepository)
			r.Put("/{repoID}", srv.handleUpdateRepository)
			r.Delete("/{repoID}", srv.handleDeleteRepository)
			r.Post("/{repoID}/sync", srv.handleSyncRepository)
		})
	})

	srv.router = router

	// Create HTTP server
	srv.httpServer = &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return srv, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("server not initialized")
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	return s.httpServer.Shutdown(ctx)
}

// handleHealth is a simple health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Code        int    `json:"code"`
	Description string `json:"description,omitempty"`
}

// renderError renders an error response
func renderError(w http.ResponseWriter, r *http.Request, httpStatus int, err error, description string) {
	errResp := ErrorResponse{
		Error:       http.StatusText(httpStatus),
		Code:        httpStatus,
		Description: description,
	}

	if err != nil {
		errResp.Description = fmt.Sprintf("%s: %v", description, err)
	}

	render.Status(r, httpStatus)
	render.JSON(w, r, errResp)
}

// Project handlers

func (s *Server) handleListProjects(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleGetProject(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleArchiveProject(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleActivateProject(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

// Repository handlers

func (s *Server) handleListRepositories(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleListProjectRepositories(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleCreateRepository(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleGetRepository(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleUpdateRepository(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleDeleteRepository(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

func (s *Server) handleSyncRepository(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, http.StatusNotImplemented, nil, "Not implemented")
}

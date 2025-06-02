package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/JadenRazo/Project-Website/backend/internal/projects/service"
	projectHTTP "github.com/JadenRazo/Project-Website/backend/internal/projects/delivery/http"
)

func main() {
	// Set gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://jadenrazo.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Initialize in-memory project service
	projectService := service.NewMemoryProjectService()

	// Initialize HTTP handler
	projectHandler := projectHTTP.NewHandler(projectService)

	// Register routes
	api := r.Group("/api/v1")
	projectsGroup := api.Group("/projects")
	projectHandler.RegisterRoutes(projectsGroup)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Simple API server running",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting simple API server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
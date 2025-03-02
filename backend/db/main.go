package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	
	"url-shortener/db"
	"url-shortener/handlers"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	
	// Initialize database
	if err := db.InitDB("data/urls.db"); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Set up Gin router
	router := gin.Default()
	
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://jadenrazo.dev"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Serve static files
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	
	// API routes
	api := router.Group("/api")
	{
		api.POST("/shorten", handlers.ShortenURL)
		api.GET("/urls", handlers.GetAllURLs)
	}
	
	// Redirect route
	router.GET("/:shortCode", handlers.RedirectURL)
	
	// Home page for quick testing
	router.GET("/", handlers.HomePage)
	
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(router.Run(":" + port))
}

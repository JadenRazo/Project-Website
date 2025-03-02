package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	
	"./db"
	"./handlers"
	"./middleware"
	"./config"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	
	// Load configuration
	appConfig := config.LoadConfig()
	
	// Create data directory if it doesn't exist
	os.MkdirAll("data", 0755)
	
	// Initialize database
	if err := db.InitDB(appConfig.DatabasePath); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Set Gin mode based on environment
	if appConfig.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// Create rate limiters
	apiLimiter := middleware.NewRateLimiter(appConfig.APIRateLimit, time.Minute)
	redirectLimiter := middleware.NewRateLimiter(appConfig.RedirectRateLimit, time.Minute)
	
	// Set up Gin router
	router := gin.Default()
	
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     appConfig.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Serve static files
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	
	// Public routes
	router.GET("/", handlers.HomePage)
	router.GET("/:shortCode", redirectLimiter.RateLimitMiddleware(), handlers.RedirectURL)
	
	// API routes with rate limiting
	api := router.Group("/api", apiLimiter.RateLimitMiddleware())
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}
		
		// URL routes
		urls := api.Group("/urls")
		{
			// Public endpoint for creating short URLs
			urls.POST("/shorten", handlers.ShortenURL)
			
			// Protected endpoints for URL management
			authURLs := urls.Group("/", middleware.AuthMiddleware())
			{
				authURLs.GET("", handlers.GetUserURLs)
				authURLs.GET("/:shortCode", handlers.GetURLDetails)
				authURLs.PUT("/:shortCode", handlers.UpdateURL)
				authURLs.DELETE("/:shortCode", handlers.DeleteURL)
				authURLs.GET("/:shortCode/analytics", handlers.GetURLAnalytics)
			}
		}
		
		// User routes (all protected)
		user := api.Group("/user", middleware.AuthMiddleware())
		{
			user.GET("/profile", handlers.GetUserProfile)
			user.PUT("/profile", handlers.UpdateUserProfile)
			user.GET("/analytics", handlers.GetUserAnalytics)
			
			// Custom domain routes
			domains := user.Group("/domains")
			{
				domains.GET("", handlers.GetUserDomains)
				domains.POST("", handlers.AddCustomDomain)
			}
		}
		
		// Admin routes (protected and admin-only)
		admin := api.Group("/admin", middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// Add admin endpoints here if needed
		}
	}
	
	// Start server
	port := appConfig.Port
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(router.Run(":" + port))
}

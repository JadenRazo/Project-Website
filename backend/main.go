package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	
	"github.com/JadenRazo/Project-Website/backend/db"
	"github.com/JadenRazo/Project-Website/backend/handlers"
	"github.com/JadenRazo/Project-Website/backend/middleware"
	"github.com/JadenRazo/Project-Website/backend/config"
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
	
	// Initialize database using the new system
	db.Initialize()
	dbManager := db.DefaultManager
	
	// Configure the database connection
	dbManager.SetSQLite(appConfig.DatabasePath)
	if err := dbManager.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Run migrations using the migration manager
	migrationManager := db.NewMigrationManager(dbManager)
	if err := migrationManager.Migrate(); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	
	// Ensure we close the connection when the program exits
	defer func() {
		if err := dbManager.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()
	
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
			// Admin dashboard routes
			admin.GET("/dashboard", handlers.AdminDashboard)
			admin.GET("/users", handlers.GetAllUsers)
			admin.PUT("/users/:id", handlers.UpdateUserStatus)
			admin.GET("/urls", handlers.GetAllURLs)
			admin.DELETE("/urls/:id", handlers.AdminDeleteURL)
		}
	}
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"version": appConfig.Version,
			"environment": appConfig.Environment,
		})
	})
	
	// Start server
	port := appConfig.Port
	
	log.Printf("Starting server in %s mode", appConfig.Environment)
	fmt.Printf("Server running on http://localhost:%s\n", port)
	
	// Run the server with graceful shutdown
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

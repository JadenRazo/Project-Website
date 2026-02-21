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

	"github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath"
	projectpathHttp "github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath/delivery/http"
	projectpathRepo "github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/certification"
	certHttp "github.com/JadenRazo/Project-Website/backend/internal/devpanel/certification/delivery/http"
	certRepo "github.com/JadenRazo/Project-Website/backend/internal/devpanel/certification/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/prompt"
	promptHttp "github.com/JadenRazo/Project-Website/backend/internal/devpanel/prompt/delivery/http"
	promptRepo "github.com/JadenRazo/Project-Website/backend/internal/devpanel/prompt/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/skill"
	skillHttp "github.com/JadenRazo/Project-Website/backend/internal/devpanel/skill/delivery/http"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel/skill/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/visitor"
	"github.com/gin-contrib/cors"
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

	fmt.Println("Starting DevPanel service...")

	// Initialize database
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Database schema should be initialized using backend/schema.sql
	// All tables including skills, certifications, and categories are defined there
	log.Println("Using existing database schema from backend/schema.sql")

	// Initialize cache service
	// Generate a 32-byte encryption key for AES-256
	encryptionKey := make([]byte, 32)
	// In production, load this from environment or config
	copy(encryptionKey, []byte("devpanel-secure-encryption-key!"))

	redisClient, err := cache.NewSecureRedisClient(cache.SecureRedisConfig{
		NetworkAddr:   "localhost:6379", // Or from cfg if available
		MaxRetries:    3,
		EncryptionKey: encryptionKey,
	})

	var cacheService *cache.SecureCache
	if err != nil {
		log.Printf("Warning: Failed to initialize Redis client: %v. Continuing without cache.", err)
		cacheService = nil
	} else {
		cacheService, err = cache.NewSecureCache(redisClient)
		if err != nil {
			log.Printf("Warning: Failed to initialize cache service: %v. Continuing without cache.", err)
			cacheService = nil
		}
	}

	// Initialize metrics manager
	metricsManager, err := metrics.NewManager(metrics.Config{
		Enabled: true,
		Grafana: metrics.GrafanaConfig{
			Enabled:   true,
			Namespace: "devpanel",
			Subsystem: "http",
			Address:   ":9090",
			Endpoint:  "/metrics",
		},
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize metrics: %v", err)
		metricsManager = nil
	}

	// Initialize auth handlers
	authHandlers := auth.NewAdminAuthHandlers(database)

	// Create service manager for devpanel to manage
	serviceManager := core.NewServiceManager()

	// Initialize visitor service with proper dependencies
	visitorService := visitor.NewService(database, cacheService, metricsManager, visitor.Config{
		EnableTracking:         true,
		SessionTimeout:         30 * time.Minute,
		RealtimeTimeout:        5 * time.Minute,
		MaxPageViewsPerSession: 100,
		EnableBotDetection:     true,
		PrivacyMode:            "balanced",
	})

	metricsCollector := devpanel.NewMetricsCollector(devpanel.Config{
		MetricsInterval: 30 * time.Second,
	})

	devpanelService := devpanel.NewService(serviceManager, visitorService, metricsCollector, devpanel.Config{
		AdminToken:      "",
		MetricsInterval: 30 * time.Second,
		MaxLogLines:     1000,
		LogRetention:    7 * 24 * time.Hour,
	})

	// Register service with manager
	if err := serviceManager.RegisterService(devpanelService); err != nil {
		log.Fatalf("Failed to register devpanel service: %v", err)
	}

	// Start metrics collector
	metricsCollector.StartCollecting(serviceManager)

	// Set up router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	if cfg.Environment == "development" {
		config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	} else {
		config.AllowOrigins = []string{"https://jadenrazo.dev", "https://www.jadenrazo.dev"}
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Health check endpoint (for Docker healthcheck)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "devpanel"})
	})

	// Create auth routes (public)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/admin/login", authHandlers.Login)
		authGroup.POST("/admin/validate", authHandlers.ValidateToken)
		authGroup.POST("/admin/setup/request", authHandlers.RequestSetup)
		authGroup.POST("/admin/setup/complete", authHandlers.CompleteSetup)
		authGroup.GET("/admin/setup/status", authHandlers.CheckSetupStatus)
	}

	// Initialize services first
	skillRepo := repository.NewGormRepository(database)
	skillService := skill.NewService(skillRepo)
	skillHandler := skillHttp.NewHandler(skillService)

	certificationRepo := certRepo.NewGormRepository(database)
	certificationService := certification.NewService(certificationRepo)
	certificationHandler := certHttp.NewHandler(certificationService)

	promptRepo := promptRepo.NewGormRepository(database)
	promptService := prompt.NewService(promptRepo)
	promptHandler := promptHttp.NewHandler(promptService)

	// Initialize project path service
	projectPathRepo := projectpathRepo.NewGormRepository(database)
	projectPathService := projectpath.NewService(projectPathRepo)
	projectPathHandler := projectpathHttp.NewHandler(projectPathService)

	// Create public routes group (no auth required) - MUST BE BEFORE PROTECTED ROUTES
	publicGroup := router.Group("/api/v1/devpanel/public")
	{
		// Register public certification routes
		publicGroup.GET("/certifications", certificationHandler.GetVisibleCertifications)
		publicGroup.GET("/certification-categories", certificationHandler.GetVisibleCategories)

		// Register public prompt routes
		publicGroup.GET("/prompts", promptHandler.GetVisiblePrompts)
		publicGroup.GET("/prompt-categories", promptHandler.GetVisibleCategories)

		// Register public skills routes
		publicGroup.GET("/skills/featured", skillHandler.GetFeaturedSkills)
		publicGroup.GET("/skills/categories", skillHandler.GetCategories)
		publicGroup.GET("/skills/proficiency-levels", skillHandler.GetProficiencyLevels)
	}

	// Create a router group for the devpanel service (protected)
	routerGroup := router.Group("/api/v1/devpanel")
	routerGroup.Use(authHandlers.AuthMiddleware())

	// Register devpanel routes
	devpanelService.RegisterRoutes(routerGroup)

	// Register skill routes
	skillGroup := routerGroup.Group("/skills")
	skillHandler.RegisterRoutes(skillGroup)

	// Register certification routes (protected)
	certificationHandler.RegisterRoutes(routerGroup)

	// Register prompt routes (protected)
	promptHandler.RegisterRoutes(routerGroup)

	// Register project path routes (protected)
	projectPathGroup := routerGroup.Group("/project-paths")
	projectPathHandler.RegisterRoutes(projectPathGroup)

	// Register visitor analytics routes (protected)
	visitorGroup := routerGroup.Group("/visitor")
	{
		visitorGroup.GET("/analytics", func(c *gin.Context) {
			analytics, err := visitorService.GetVisitorStats(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
				return
			}
			c.JSON(http.StatusOK, analytics)
		})
		visitorGroup.GET("/realtime", func(c *gin.Context) {
			count := visitorService.GetRealTimeCount(c.Request.Context())
			c.JSON(http.StatusOK, gin.H{"count": count})
		})
		visitorGroup.GET("/metrics", func(c *gin.Context) {
			// Get timeline data instead
			period := c.DefaultQuery("period", "7d")
			interval := c.DefaultQuery("interval", "day")
			metrics, err := visitorService.GetTimelineData(c.Request.Context(), period, interval)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics"})
				return
			}
			c.JSON(http.StatusOK, metrics)
		})
	}

	// Hardcoded port for DevPanel service (8081 from status report)
	const devPanelPort = 8081

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", devPanelPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("DevPanel service listening on port %d\n", devPanelPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down DevPanel service...")

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("DevPanel service stopped")
}
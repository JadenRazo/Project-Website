package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"

	appconfig "github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/blog"
	blogHTTP "github.com/JadenRazo/Project-Website/backend/internal/blog/delivery/http"
	blogRepo "github.com/JadenRazo/Project-Website/backend/internal/blog/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/codestats"
	"github.com/JadenRazo/Project-Website/backend/internal/contact"
	codeStatsHTTP "github.com/JadenRazo/Project-Website/backend/internal/codestats/delivery/http"
	"github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath"
	projectPathHTTP "github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath/delivery/http"
	projectPathRepo "github.com/JadenRazo/Project-Website/backend/internal/codestats/projectpath/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth/oauth"
	"github.com/JadenRazo/Project-Website/backend/internal/common/auth/totp"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/metrics"
	"github.com/JadenRazo/Project-Website/backend/internal/common/middleware"
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	coreConfig "github.com/JadenRazo/Project-Website/backend/internal/core/config"
	"github.com/JadenRazo/Project-Website/backend/internal/core/db"
	"github.com/JadenRazo/Project-Website/backend/internal/devpanel"
	// "github.com/JadenRazo/Project-Website/backend/internal/devpanel/project"
	"github.com/JadenRazo/Project-Website/backend/internal/gateway"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging"
	projectHTTP "github.com/JadenRazo/Project-Website/backend/internal/projects/delivery/http"
	projectMemoryService "github.com/JadenRazo/Project-Website/backend/internal/projects/service"
	"github.com/JadenRazo/Project-Website/backend/internal/status"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener"
	"github.com/JadenRazo/Project-Website/backend/internal/visitor"
	"github.com/JadenRazo/Project-Website/backend/internal/worker"
)

const (
	appName    = "project-website-api"
	appVersion = "1.0.0"
)

type databaseAdapter struct {
	db *gorm.DB
}

func (d *databaseAdapter) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (d *databaseAdapter) GetDB() *gorm.DB {
	return d.db
}

func (d *databaseAdapter) Transaction(fn func(tx *gorm.DB) error) error {
	return d.db.Transaction(fn)
}

func (d *databaseAdapter) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func verifyVisitorTables(db *gorm.DB) error {
	requiredTables := []string{
		"visitor_sessions",
		"page_views",
		"visitor_realtime",
		"visitor_metrics",
		"visitor_daily_summary",
		"privacy_consents",
		"visitor_locations",
	}

	for _, table := range requiredTables {
		var exists bool
		err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", table).Scan(&exists).Error
		if err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("table %s does not exist", table)
		}
	}

	return nil
}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func main() {
	fmt.Println("Application starting...")
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			fmt.Printf("PANIC: %v\nStack trace: %s\n", r, buf[:n])
			os.Exit(1)
		}
	}()

	if _, err := os.Stat(".env"); err == nil {
		if err := loadEnvFile(".env"); err != nil {
			fmt.Printf("Warning: Failed to load .env file: %v\n", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Loading configuration...")
	cfg := coreConfig.GetConfig()
	fmt.Println("Configuration loaded.")

	err := logger.InitLogger(&appconfig.LoggingConfig{
		Level:      cfg.Logging.Level,
		Format:     "json",
		Output:     "file",
		TimeFormat: time.RFC3339Nano,
		Filename:   filepath.Join("logs", "api.log"),
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}, appName, appVersion)

	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Shutdown()

	logger.Info("Application starting up", "name", appName, "version", appVersion, "environment", cfg.Environment)
	logger.Info("Configuration loaded")

	fmt.Println("Initializing service manager...")
	serviceManager := core.NewServiceManager()
	fmt.Println("Service manager initialized.")

	fmt.Println("Connecting to database...")
	gormDB, err := db.GetDB()
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	fmt.Println("Database connection established.")

	databaseWrapper := &databaseAdapter{db: gormDB}

	logger.Info("Database connection established")

	fmt.Println("Initializing secure cache...")
	logger.Info("Initializing secure cache with Redis config", "host", cfg.Redis.Host, "port", cfg.Redis.Port)
	secureRedisConfig := cache.SecureRedisConfig{
		NetworkAddr:   cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password:      cfg.Redis.Password,
	}
	if cfg.Redis.EncryptionKey == "" {
		logger.Info("Generating random encryption key")
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			logger.Fatal("Failed to generate encryption key", "error", err)
		}
		secureRedisConfig.EncryptionKey = key
	} else {
		logger.Info("Using encryption key from config")
		secureRedisConfig.EncryptionKey = []byte(cfg.Redis.EncryptionKey)
	}

	secureRedisClient, err := cache.NewSecureRedisClient(secureRedisConfig)
	if err != nil {
		logger.Fatal("Failed to initialize secure redis client", "error", err)
	}
	secureCacheInstance, err := cache.NewSecureCache(secureRedisClient)
	if err != nil {
		logger.Fatal("Failed to initialize secure cache", "error", err)
	}
	logger.Info("Secure cache initialized")
	fmt.Println("Secure cache initialized.")

	fmt.Println("Initializing auth service...")
	authConfig := &appconfig.AuthConfig{
		JWTSecret:      cfg.Auth.JWTSecret,
		TokenExpiry:    15 * time.Minute,
		RefreshExpiry:  7 * 24 * time.Hour,
		Issuer:         "jadenrazo-api",
		Audience:       "jadenrazo-web",
		CookieSecure:   cfg.Environment == "production",
		CookieHTTPOnly: true,
	}
	authService, err := auth.New(authConfig, databaseWrapper, secureCacheInstance)
	if err != nil {
		logger.Fatal("Failed to initialize auth service", "error", err)
	}
	logger.Info("Auth service initialized")
	fmt.Println("Auth service initialized.")

	adminAuthHandlers := auth.NewAdminAuthHandlers(gormDB)
	logger.Info("Admin auth handlers initialized")
	fmt.Println("Admin auth handlers initialized.")

	fmt.Println("Initializing OAuth providers...")
	oauth2Config := auth.GetOAuth2Config()
	oauthManager := oauth.NewOAuthManager(oauth2Config)
	jwtManager := auth.NewJWTManager(authConfig)
	oauthHandlers := oauth.NewOAuthHandlers(oauthManager, secureCacheInstance, gormDB, jwtManager, oauth2Config)
	logger.Info("OAuth handlers initialized", "providers", len(oauthManager.GetEnabledProviders()))
	fmt.Printf("OAuth initialized with %d providers.\n", len(oauthManager.GetEnabledProviders()))

	// Discord handlers disabled - implementation pending
	// var discordHandlers *oauth.DiscordHandlers
	// if oauth2Config.Discord.Enabled {
	// 	discordProvider := oauth.NewDiscordProvider(...)
	// 	discordHandlers = oauth.NewDiscordHandlers(...)
	// }
	_ = oauth2Config.Discord // Silence unused warning

	logger.Info("Initializing services")
	fmt.Println("Initializing services...")

	urlShortenerConfig := urlshortener.Config{
		BaseURL:      cfg.URLShortener.BaseURL,
		MaxURLLength: 2048,
		MinURLLength: 5,
	}

	messagingConfig := messaging.Config{
		WebSocketPort:    8081,
		MaxMessageSize:   cfg.Messaging.MaxMessageSize,
		MaxAttachments:   10,
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
	}

	urlShortenerService := urlshortener.NewService(gormDB, urlShortenerConfig)
	urlShortenerService.SetAuth(authService)
	messagingService := messaging.NewService(gormDB, messagingConfig)

	projectPathRepository := projectPathRepo.NewGormRepository(gormDB)

	codeStatsService := codestats.NewService(gormDB, projectPathRepository)

	projectPathService := projectpath.NewServiceWithStatsUpdater(projectPathRepository, codeStatsService)

	projectService := projectMemoryService.NewMemoryProjectService()

	metricsConfig := metrics.DefaultConfig()
	metricsManager, err := metrics.NewManagerWithDB(metricsConfig, gormDB)
	if err != nil {
		logger.Warn("Failed to initialize metrics manager", "error", err)
		metricsManager = nil
	}

	statusService := status.NewService(gormDB, 30*time.Second, metricsManager)

	visitorConfig := visitor.Config{
		EnableTracking:     true,
		SessionTimeout:     30 * time.Minute,
		RealtimeTimeout:    5 * time.Minute,
		MaxPageViewsPerSession: 100,
		EnableBotDetection: true,
		PrivacyMode:        "balanced",
	}

	// Check visitor tables BEFORE creating the service
	if err := verifyVisitorTables(gormDB); err != nil {
		logger.Error("Visitor Analytics disabled: Database tables not found", "error", err)
		fmt.Printf("\n=== VISITOR ANALYTICS ERROR ===\n")
		fmt.Printf("Visitor analytics tables are missing from the database.\n")
		fmt.Printf("Please apply the schema: psql -U user -d db < backend/schema.sql\n")
		fmt.Printf("Error details: %v\n", err)
		fmt.Printf("===============================\n\n")
		visitorConfig.EnableTracking = false
	}

	// Now create the service with the properly configured config
	visitorService := visitor.NewService(gormDB, secureCacheInstance, metricsManager, visitorConfig)

	contactEmailConfig := &contact.EmailConfig{
		ResendAPIKey: os.Getenv("RESEND_API_KEY"),
		FromEmail:    cfg.Contact.FromEmail,
		ToEmail:      cfg.Contact.ContactToEmail,
	}
	contactHandler := contact.NewHandler(gormDB, contactEmailConfig)
	if contactEmailConfig.IsConfigured() {
		logger.Info("Contact form email configured via Resend API")
	} else {
		logger.Info("Contact form email not configured, submissions will be logged only")
	}

	blogRepository := blogRepo.NewGormRepository(gormDB)
	blogService := blog.NewService(blogRepository, secureCacheInstance)
	blogHandler := blogHTTP.NewHandler(blogService)

	workerService := worker.NewService(gormDB)

	metricsCollector := devpanel.NewMetricsCollector(devpanel.Config{
		MetricsInterval: 30 * time.Second,
	})

	devpanelService := devpanel.NewService(serviceManager, visitorService, metricsCollector, devpanel.Config{
		AdminToken:      cfg.AdminToken,
		MetricsInterval: 30 * time.Second,
		MaxLogLines:     1000,
		LogRetention:    7 * 24 * time.Hour,
	})

	logManager := devpanel.NewLogManager(
		filepath.Join("logs", "services"),
		cfg.MaxLogLines,
		cfg.LogRetention,
	)
	logManager.StartCleanup()

	fmt.Println("Services initialized.")

	logger.Info("Registering services with service manager")
	fmt.Println("Registering services with service manager...")
	if err := serviceManager.RegisterService(urlShortenerService); err != nil {
		logger.Fatal("Failed to register URL shortener service", "error", err)
	}
	if err := serviceManager.RegisterService(messagingService); err != nil {
		logger.Fatal("Failed to register messaging service", "error", err)
	}
	if err := serviceManager.RegisterService(devpanelService); err != nil {
		logger.Fatal("Failed to register devpanel service", "error", err)
	}
	if err := serviceManager.RegisterService(visitorService); err != nil {
		logger.Fatal("Failed to register visitor service", "error", err)
	}
	if err := serviceManager.RegisterService(workerService); err != nil {
		logger.Fatal("Failed to register worker service", "error", err)
	}
	fmt.Println("Services registered.")

	logger.Info("Starting metrics collector")
	fmt.Println("Starting metrics collector...")
	metricsCollector.StartCollecting(serviceManager)
	fmt.Println("Metrics collector started.")

	logger.Info("Initializing API gateway")
	fmt.Println("Initializing API gateway...")
	apiGateway := gateway.NewGateway(serviceManager)
	fmt.Println("API gateway initialized.")

	apiGateway.AddMiddleware(middleware.ErrorRecovery())

	apiGateway.AddMiddleware(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	apiGateway.AddMiddleware(gzip.Gzip(gzip.DefaultCompression))

	apiGateway.AddMiddleware(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	})

	if metricsManager != nil {
		apiGateway.AddMiddleware(metricsManager.GinLatencyMiddleware())
	}

	apiGateway.AddMiddleware(middleware.APIRateLimiter())

	apiGateway.AddMiddleware(middleware.RequestValidation())
	apiGateway.AddMiddleware(middleware.QueryParamValidation())

	if cfg.Environment == "development" || cfg.EnablePprof {
		logger.Info("Enabling profiling endpoints")
		pprofGroup := apiGateway.GetRouter().Group("/debug/pprof")
		pprofGroup.Use(adminAuthHandlers.AuthMiddleware())
		{
			pprofGroup.GET("/", gin.WrapF(pprof.Index))
			pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
			pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
			pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
			pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
			pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
			pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		}
	}

	var allowedOrigins []string
	if len(cfg.App.AllowedOrigins) > 0 {
		allowedOrigins = cfg.App.AllowedOrigins
		logger.Info("Using configured CORS origins", "origins", allowedOrigins)
	} else if cfg.AllowedOrigins != "" {
		allowedOrigins = strings.Split(cfg.AllowedOrigins, ",")
		logger.Info("Using CORS origins from string config", "origins", allowedOrigins)
	} else {
		if cfg.Environment == "development" {
			allowedOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
		} else {
			allowedOrigins = []string{"https://jadenrazo.dev", "https://www.jadenrazo.dev"}
		}
		logger.Info("Using default CORS origins", "origins", allowedOrigins)
	}

	apiGateway.AddMiddleware(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := apiGateway.GetRouter()

	router.Use(middleware.SecurityHeaders())
	logger.Info("Security headers middleware registered")

	router.GET("/ws/analytics", func(c *gin.Context) {
		visitorService.ServeWs(visitorService.GetHub(), c)
	})

	router.Use(visitor.TrackingMiddleware(visitorService))

	logger.Info("Registering service routes")
	fmt.Println("Registering service routes...")
	totpHandlers := totp.NewTOTPHandlers(gormDB)

	apiGateway.RegisterService("auth", func(rg *gin.RouterGroup) {
		authService.RegisterUserRoutes(rg)

		admin := rg.Group("/admin")
		{
			admin.POST("/login", middleware.StrictRateLimiter(), adminAuthHandlers.Login)
			admin.POST("/validate", adminAuthHandlers.ValidateToken)
			admin.POST("/setup/request", middleware.StrictRateLimiter(), adminAuthHandlers.RequestSetup)
			admin.POST("/setup/complete", middleware.StrictRateLimiter(), adminAuthHandlers.CompleteSetup)
			admin.GET("/setup/status", adminAuthHandlers.CheckSetupStatus)
			admin.POST("/mfa/verify", middleware.StrictRateLimiter(), adminAuthHandlers.VerifyMFA)

			oauth := admin.Group("/oauth")
			{
				oauth.GET("/login/:provider", oauthHandlers.InitiateOAuth)
				oauth.GET("/callback/:provider", oauthHandlers.HandleCallback)
			}

			protected := admin.Group("/mfa")
			protected.Use(adminAuthHandlers.AuthMiddleware())
			{
				protected.POST("/totp/setup", totpHandlers.SetupTOTP)
				protected.POST("/totp/verify", middleware.StrictRateLimiter(), totpHandlers.VerifyTOTP)
				protected.POST("/totp/disable", totpHandlers.DisableTOTP)
				protected.GET("/status", totpHandlers.GetMFAStatus)
				protected.POST("/backup/regenerate", totpHandlers.RegenerateBackupCodes)
			}
		}
	})

	// Discord routes disabled - implementation pending
	// if discordHandlers != nil {
	// 	apiGateway.RegisterService("discord", func(rg *gin.RouterGroup) {
	// 		rg.GET("/connect", discordHandlers.InitiateConnect)
	// 		rg.GET("/callback", discordHandlers.HandleCallback)
	// 		rg.GET("/stats", discordHandlers.GetConnectionStats)
	// 	})
	// 	logger.Info("Discord linked roles routes registered")
	// }

	apiGateway.RegisterService("urls", urlShortenerService.RegisterRoutes)
	apiGateway.RegisterService("messaging", messagingService.RegisterRoutes)
	apiGateway.RegisterService("devpanel", devpanelService.RegisterRoutes)

	codeStatsHandler := codeStatsHTTP.NewHandler(codeStatsService)
	apiGateway.RegisterService("code", codeStatsHandler.RegisterRoutes)

	projectPathHandler := projectPathHTTP.NewHandler(projectPathService)
	apiGateway.RegisterService("code/paths", projectPathHandler.RegisterRoutes)

	projectHandler := projectHTTP.NewHandler(projectService)
	apiGateway.RegisterService("projects", projectHandler.RegisterRoutes)

	apiGateway.RegisterService("status", statusService.RegisterRoutes)

	apiGateway.RegisterService("visitor", visitorService.RegisterRoutes)

	apiGateway.RegisterService("contact", contactHandler.RegisterRoutes)

	apiGateway.RegisterService("blog", blogHandler.RegisterPublicRoutes)

	blogAdminGroup := router.Group("/api/v1/blog/admin")
	blogAdminGroup.Use(adminAuthHandlers.AuthMiddleware())
	blogHandler.RegisterAdminRoutes(blogAdminGroup)

	apiGateway.RegisterHealthCheck()

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	frontendBuildPath := filepath.Join("..", "frontend", "build")
	if _, err := os.Stat(frontendBuildPath); err == nil {
		logger.Info("Serving frontend build files", "path", frontendBuildPath)
		router.Static("/static", filepath.Join(frontendBuildPath, "static"))
		router.StaticFile("/manifest.json", filepath.Join(frontendBuildPath, "manifest.json"))
		router.StaticFile("/favicon.ico", filepath.Join(frontendBuildPath, "favicon.ico"))
		router.StaticFile("/robots.txt", filepath.Join(frontendBuildPath, "robots.txt"))
	}

	publicPath := filepath.Join("..", "frontend", "public")
	if _, err := os.Stat(publicPath); err == nil {
		logger.Info("Serving public files", "path", publicPath)
		router.StaticFile("/code_stats.json", filepath.Join(publicPath, "code_stats.json"))
		router.StaticFile("/apple-touch-icon.png", filepath.Join(publicPath, "apple-touch-icon.png"))
		router.StaticFile("/favicon-16x16.png", filepath.Join(publicPath, "favicon-16x16.png"))
		router.StaticFile("/favicon-32x32.png", filepath.Join(publicPath, "favicon-32x32.png"))
	}

	router.GET("/:shortCode", urlShortenerService.RedirectHandler)

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") ||
			strings.HasPrefix(c.Request.URL.Path, "/health") ||
			strings.HasPrefix(c.Request.URL.Path, "/metrics") ||
			strings.HasPrefix(c.Request.URL.Path, "/debug") {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		indexPath := filepath.Join(frontendBuildPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		} else {
			c.JSON(404, gin.H{"error": "Frontend build not found"})
		}
	})

	fmt.Println("Service routes registered.")

	logger.Info("Starting all registered services")
	fmt.Println("Starting all registered services...")
	var failedServices []string
	for _, service := range serviceManager.GetAllServices() {
		if err := service.Start(); err != nil {
			logger.Error("Failed to start service", "service", service.Name(), "error", err)
			failedServices = append(failedServices, service.Name())
		} else {
			logger.Info("Service started", "service", service.Name())
		}
	}
	if len(failedServices) > 0 {
		logger.Error("Some services failed to start", "failed_services", failedServices)
		fmt.Printf("WARNING: Failed to start services: %v\n", failedServices)
	} else {
		fmt.Println("All services started successfully.")
	}

	serviceManager.StartHealthChecks()


	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        apiGateway.GetRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		fmt.Println("Starting server on port", cfg.Port)
		var err error
		if cfg.Server.TLSEnabled {
			logger.Info("TLS enabled", "cert_path", cfg.Server.TLSCert, "key_path", cfg.Server.TLSKey)
			if _, errCert := os.Stat(cfg.Server.TLSCert); os.IsNotExist(errCert) {
				logger.Fatal("TLS cert file not found", "path", cfg.Server.TLSCert)
			}
			if _, errKey := os.Stat(cfg.Server.TLSKey); os.IsNotExist(errKey) {
				logger.Fatal("TLS key file not found", "path", cfg.Server.TLSKey)
			}
			err = srv.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey)
		} else {
			logger.Info("TLS disabled, starting HTTP server")
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	logger.Info("Server started successfully, waiting for shutdown signal")
	fmt.Println("Server started successfully, waiting for shutdown signal")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Received shutdown signal", "signal", sig.String())

	shutdownTimeout := 5 * time.Second

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdownCancel()

	logger.Info("Gracefully shutting down server")

	logger.Info("Stopping all services")
	if err := serviceManager.StopAllServices(); err != nil {
		logger.Error("Error stopping services", "error", err)
	}

	logger.Info("Shutting down HTTP server", "timeout", shutdownTimeout)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		return
	}

	logger.Info("Closing database connection")
	if err := db.CloseDB(); err != nil {
		logger.Error("Error closing database", "error", err)
	}

	logger.Info("Server exited properly")
	fmt.Println("Server exited properly.")
}

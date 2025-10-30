package urlshortener

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/core"
	"github.com/JadenRazo/Project-Website/backend/internal/urlshortener/entity"
)

// Service handles URL shortener operations
type Service struct {
	*core.BaseService
	db     *gorm.DB
	config Config
	auth   *auth.Auth
}

// Config holds URL shortener service configuration
type Config struct {
	BaseURL      string
	MaxURLLength int
	MinURLLength int
}

// NewService creates a new URL shortener service
func NewService(db *gorm.DB, config Config) *Service {
	return &Service{
		BaseService: core.NewBaseService("urlshortener"),
		db:          db,
		config:      config,
	}
}

// SetAuth sets the auth service for the URL shortener
func (s *Service) SetAuth(auth *auth.Auth) {
	s.auth = auth
}

// RegisterRoutes registers the service's HTTP routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	// Health check endpoint
	router.GET("/health", s.HealthCheckHandler)

	urlshortener := router.Group("/urls")
	{
		// Public endpoints - anyone can shorten URLs
		urlshortener.POST("/shorten", s.ShortenURL) // Allow anonymous shortening
		
		// Protected endpoints - require authentication
		if s.auth != nil {
			protected := urlshortener.Group("", s.auth.GinAuthMiddleware())
			{
				protected.GET("/", s.GetUserURLs)           // List user's URLs
				protected.GET("/:shortCode/stats", s.GetURLStats) // View stats
				protected.DELETE("/:shortCode", s.DeleteURL)      // Delete URL
			}
		} else {
			// Fallback if auth is not set (for backward compatibility)
			urlshortener.GET("/", s.GetUserURLs)
			urlshortener.GET("/:shortCode/stats", s.GetURLStats)
			urlshortener.DELETE("/:shortCode", s.DeleteURL)
		}
	}

	// Redirect endpoint (public, different route pattern)
	router.GET("/:shortCode", s.RedirectHandler)
}

// ShortenURL handles URL shortening requests
func (s *Service) ShortenURL(c *gin.Context) {
	s.IncrementRequests()

	var req entity.ShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.IncrementErrors()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Sanitize and validate URL
	req.URL = entity.SanitizeURL(req.URL)
	if !s.isValidURL(req.URL) {
		s.IncrementErrors()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// Get user ID if authenticated (optional for URL shortening)
	var userID *uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if parsedUID, err := uuid.Parse(uid.(string)); err == nil {
			userID = &parsedUID
		}
	}

	// Generate or validate short code
	shortCode := req.CustomCode
	if shortCode != "" {
		// Validate custom code
		if !entity.IsValidShortCode(shortCode) {
			s.IncrementErrors()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom code format"})
			return
		}

		// Check if custom code is already taken
		var existing entity.ShortenedURL
		err := s.db.Where("short_code = ?", shortCode).First(&existing).Error
		if err == nil {
			s.IncrementErrors()
			c.JSON(http.StatusConflict, gin.H{"error": "Short code already exists"})
			return
		}
		if err != gorm.ErrRecordNotFound {
			s.IncrementErrors()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check short code availability"})
			return
		}
	} else {
		// Generate unique short code
		var err error
		shortCode, err = s.generateShortCode()
		if err != nil {
			s.IncrementErrors()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short code"})
			return
		}
	}

	// Parse expiration date if provided
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		if parsedTime, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			if parsedTime.After(time.Now()) {
				expiresAt = &parsedTime
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Expiration date must be in the future"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration date format. Use RFC3339 format"})
			return
		}
	}

	// Create shortened URL
	shortenedURL := &entity.ShortenedURL{
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		ExpiresAt:   expiresAt,
	}

	err := s.db.Create(shortenedURL).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shortened URL"})
		return
	}

	// Build response
	response := entity.ShortenURLResponse{
		ID:          shortenedURL.ID,
		ShortCode:   shortenedURL.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", s.config.BaseURL, shortenedURL.ShortCode),
		OriginalURL: shortenedURL.OriginalURL,
		Title:       shortenedURL.Title,
		Description: shortenedURL.Description,
		ExpiresAt:   shortenedURL.ExpiresAt,
		CreatedAt:   shortenedURL.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "URL shortened successfully",
		"data":    response,
	})
}

// GetURLStats retrieves statistics for a shortened URL
func (s *Service) GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	// Get user ID for authorization
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find the URL and verify ownership
	var shortenedURL entity.ShortenedURL
	err = s.db.Where("short_code = ? AND user_id = ?", shortCode, uid).First(&shortenedURL).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or access denied"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		}
		return
	}

	// Get click statistics
	var totalClicks int64
	var uniqueClicks int64

	s.db.Model(&entity.URLClick{}).Where("short_url_id = ?", shortenedURL.ID).Count(&totalClicks)
	s.db.Model(&entity.URLClick{}).Where("short_url_id = ?", shortenedURL.ID).Distinct("ip_address").Count(&uniqueClicks)

	// Get clicks by country
	var countryStats []struct {
		Country string
		Count   int64
	}
	s.db.Model(&entity.URLClick{}).Select("country_code as country, COUNT(*) as count").
		Where("short_url_id = ? AND country_code != ''", shortenedURL.ID).
		Group("country_code").Order("count DESC").Limit(10).Find(&countryStats)

	clicksByCountry := make(map[string]int64)
	for _, stat := range countryStats {
		clicksByCountry[stat.Country] = stat.Count
	}

	// Get clicks by date (last 30 days)
	var dateStats []struct {
		Date  string
		Count int64
	}
	s.db.Model(&entity.URLClick{}).Select("DATE(clicked_at) as date, COUNT(*) as count").
		Where("short_url_id = ? AND clicked_at >= ?", shortenedURL.ID, time.Now().AddDate(0, 0, -30)).
		Group("DATE(clicked_at)").Order("date").Find(&dateStats)

	clicksByDate := make(map[string]int64)
	for _, stat := range dateStats {
		clicksByDate[stat.Date] = stat.Count
	}

	// Get top referers
	var refererStats []struct {
		Referer string
		Count   int64
	}
	s.db.Model(&entity.URLClick{}).Select("referer, COUNT(*) as count").
		Where("short_url_id = ? AND referer != ''", shortenedURL.ID).
		Group("referer").Order("count DESC").Limit(10).Find(&refererStats)

	topReferers := make(map[string]int64)
	for _, stat := range refererStats {
		topReferers[stat.Referer] = stat.Count
	}

	// Get recent clicks
	var recentClicks []entity.URLClick
	s.db.Where("short_url_id = ?", shortenedURL.ID).Order("clicked_at DESC").Limit(20).Find(&recentClicks)

	// Build response
	response := entity.URLStatsResponse{
		ShortenedURL:    &shortenedURL,
		TotalClicks:     totalClicks,
		UniqueClicks:    uniqueClicks,
		ClicksByCountry: clicksByCountry,
		ClicksByDate:    clicksByDate,
		RecentClicks:    recentClicks,
		TopReferers:     topReferers,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserURLs retrieves all URLs created by a user
func (s *Service) GetUserURLs(c *gin.Context) {
	// Parse pagination
	var pagination entity.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	// Ensure default values
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 20
	}

	// Check if user is authenticated
	userID, exists := c.Get("user_id")
	var query *gorm.DB

	if !exists {
		// For anonymous users, show all public URLs (those without a user_id)
		query = s.db.Where("user_id IS NULL").Order("created_at DESC")
	} else {
		// For authenticated users, show their URLs
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		query = s.db.Where("user_id = ?", uid).Order("created_at DESC")
	}

	// Get URLs with pagination
	var urls []entity.ShortenedURL
	offset := (pagination.Page - 1) * pagination.PageSize

	var total int64
	query.Model(&entity.ShortenedURL{}).Count(&total)

	err := query.Offset(offset).Limit(pagination.PageSize).Find(&urls).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
		return
	}

	// Get click statistics for each URL
	var urlsWithStats []entity.URLWithStats
	for _, url := range urls {
		var totalClicks int64
		var uniqueClicks int64

		s.db.Model(&entity.URLClick{}).Where("short_url_id = ?", url.ID).Count(&totalClicks)
		s.db.Model(&entity.URLClick{}).Where("short_url_id = ?", url.ID).Distinct("ip_address").Count(&uniqueClicks)

		urlsWithStats = append(urlsWithStats, entity.URLWithStats{
			ShortenedURL: &url,
			TotalClicks:  totalClicks,
			UniqueClicks: uniqueClicks,
		})
	}

	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	response := entity.UserURLsResponse{
		URLs: urlsWithStats,
		Pagination: entity.PaginationResponse{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    pagination.Page < totalPages,
			HasPrev:    pagination.Page > 1,
		},
	}

	c.JSON(http.StatusOK, response)
}

// DeleteURL deletes a shortened URL
func (s *Service) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find and delete the URL (verify ownership)
	result := s.db.Where("short_code = ? AND user_id = ?", shortCode, uid).Delete(&entity.ShortenedURL{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "URL deleted successfully",
	})
}

// RedirectHandler handles redirects for shortened URLs
func (s *Service) RedirectHandler(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	// Find the URL
	var shortenedURL entity.ShortenedURL
	err := s.db.Where("short_code = ?", shortCode).First(&shortenedURL).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		}
		return
	}

	// Check if URL is accessible
	if !shortenedURL.IsAccessible() {
		if shortenedURL.IsURLExpired() {
			c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "URL is not active"})
		}
		return
	}

	// Record click analytics (async to not slow down redirect)
	go s.recordClick(shortenedURL.ID, c)

	// Perform redirect
	c.Redirect(http.StatusMovedPermanently, shortenedURL.OriginalURL)
}

// HealthCheck performs service-specific health checks
func (s *Service) HealthCheck() error {
	if err := s.BaseService.HealthCheck(); err != nil {
		return err
	}

	// Check database connection
	sqlDB, err := s.db.DB()
	if err != nil {
		s.AddError(err)
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		s.AddError(err)
		return err
	}

	return nil
}

// HealthCheckHandler handles HTTP health check requests
func (s *Service) HealthCheckHandler(c *gin.Context) {
	// Test database connection
	var count int64
	err := s.db.Model(&entity.ShortenedURL{}).Count(&count).Error

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"error":   err.Error(),
			"service": "urlshortener",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "healthy",
		"service":    "urlshortener",
		"database":   "connected",
		"total_urls": count,
		"base_url":   s.config.BaseURL,
	})
}

// Helper methods

// generateShortCode generates a unique short code
func (s *Service) generateShortCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	for attempts := 0; attempts < 10; attempts++ {
		code := make([]byte, codeLength)
		for i := range code {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			code[i] = charset[num.Int64()]
		}

		shortCode := string(code)

		// Check if code already exists
		var existing entity.ShortenedURL
		err := s.db.Where("short_code = ?", shortCode).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			return shortCode, nil
		}
		if err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after 10 attempts")
}

// isValidURL validates if a URL is properly formatted
func (s *Service) isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Must have a scheme and host
	if u.Scheme == "" || u.Host == "" {
		return false
	}

	// Only allow http and https
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Validate URL length
	if len(urlStr) < 10 || len(urlStr) > 2048 {
		return false
	}

	return true
}

// recordClick records analytics data for a URL click
func (s *Service) recordClick(urlID uuid.UUID, c *gin.Context) {
	// Get client IP
	clientIP := c.ClientIP()
	var ipAddr net.IP
	if ip := net.ParseIP(clientIP); ip != nil {
		ipAddr = ip
	}

	// Get user agent
	userAgent := c.GetHeader("User-Agent")

	// Get referer
	referer := c.GetHeader("Referer")

	// Clean referer URL
	if referer != "" {
		if u, err := url.Parse(referer); err == nil {
			referer = u.Hostname()
		}
	}

	// TODO: Add GeoIP lookup for country/city
	// For now, we'll leave these empty

	// Create click record
	click := &entity.URLClick{
		ShortURLID:  urlID,
		IPAddress:   ipAddr,
		UserAgent:   userAgent,
		Referer:     referer,
		CountryCode: "",    // TODO: Implement GeoIP lookup
		City:        "",    // TODO: Implement GeoIP lookup
		DeviceType:  "",    // TODO: Parse user agent
		Browser:     "",    // TODO: Parse user agent
		OS:          "",    // TODO: Parse user agent
		IsBot:       false, // TODO: Detect bots
	}

	// Save asynchronously (don't block the redirect)
	err := s.db.Create(click).Error
	if err != nil {
		// Log error but don't fail the redirect
		fmt.Printf("Failed to record click analytics: %v\n", err)
	}
}

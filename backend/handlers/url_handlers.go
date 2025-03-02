package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"github.com/JadenRazo/Project-Website/backend/db"
	"github.com/JadenRazo/Project-Website/backend/utils"
)

// ShortenURLRequest represents the request body for shortening a URL
type ShortenURLRequest struct {
	OriginalURL string `json:"originalUrl" binding:"required,url"`
	CustomCode  string `json:"customCode,omitempty"`
}

// URLResponse represents the response returned for URL operations
type URLResponse struct {
	ShortCode   string    `json:"shortCode"`
	OriginalURL string    `json:"originalUrl"`
	ShortURL    string    `json:"shortUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	Clicks      uint      `json:"clicks"`
}

// ShortenURL handles URL shortening requests
func ShortenURL(c *gin.Context) {
	var req ShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	// Check if URL is valid
	if !utils.IsValidURL(req.OriginalURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}
	
	shortCode := req.CustomCode
	if shortCode == "" {
		// Generate a random short code
		shortCode = generateShortCode(6)
	}
	
	// Check if custom code is already taken
	var existingURL db.URL
	if result := db.DB.Where("short_code = ?", shortCode).First(&existingURL); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "This short code is already in use"})
		return
	}
	
	// Set expiration date (optional)
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
	
	// Create new URL record
	url := db.URL{
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   &expiresAt,
		Clicks:      0,
	}
	
	if result := db.DB.Create(&url); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shortened URL"})
		return
	}
	
	// Construct full short URL (using request host)
	shortURL := c.Request.Host + "/" + shortCode
	if c.Request.TLS != nil {
		shortURL = "https://" + shortURL
	} else {
		shortURL = "http://" + shortURL
	}
	
	c.JSON(http.StatusCreated, URLResponse{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		ShortURL:    shortURL,
		CreatedAt:   url.CreatedAt,
		ExpiresAt:   url.ExpiresAt,
		Clicks:      url.Clicks,
	})
}

// RedirectURL redirects from short URL to original URL
func RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	
	var url db.URL
	if result := db.DB.Where("short_code = ?", shortCode).First(&url); result.Error != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"message": "Short URL not found",
		})
		return
	}
	
	// Check if URL has expired
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		c.HTML(http.StatusGone, "expired.html", gin.H{
			"message": "This link has expired",
		})
		return
	}
	
	// Increment click count
	db.DB.Model(&url).Update("clicks", url.Clicks+1)
	
	// Redirect to original URL
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetAllURLs returns all URLs (could be paginated in a real app)
func GetAllURLs(c *gin.Context) {
	var urls []db.URL
	db.DB.Order("created_at desc").Find(&urls)
	
	var response []URLResponse
	for _, url := range urls {
		shortURL := c.Request.Host + "/" + url.ShortCode
		if c.Request.TLS != nil {
			shortURL = "https://" + shortURL
		} else {
			shortURL = "http://" + shortURL
		}
		
		response = append(response, URLResponse{
			ShortCode:   url.ShortCode,
			OriginalURL: url.OriginalURL,
			ShortURL:    shortURL,
			CreatedAt:   url.CreatedAt,
			ExpiresAt:   url.ExpiresAt,
			Clicks:      url.Clicks,
		})
	}
	
	c.JSON(http.StatusOK, response)
}

// HomePage renders a simple homepage
func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "URL Shortener",
	})
}

// generateShortCode creates a random string for short URLs
func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	// Initialize the random generator with a seed
	rand.Seed(time.Now().UnixNano())
	
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	
	return string(code)
}

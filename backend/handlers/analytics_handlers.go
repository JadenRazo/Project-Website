package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	
	"urlshortener/db"
)

// GetURLAnalytics returns analytics data for a URL
func GetURLAnalytics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You must be logged in"})
		return
	}
	
	shortCode := c.Param("shortCode")
	
	var url db.URL
	if result := db.DB.Where("short_code = ?", shortCode).First(&url); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	
	// Check if user owns the URL or is an admin
	role, _ := c.Get("role")
	if url.UserID != userID.(uint) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this URL's analytics"})
		return
	}
	
	// Get time period from query params
	period := c.DefaultQuery("period", "week")
	
	// Calculate start time based on period
	var startTime time.Time
	now := time.Now()
	
	switch period {
	case "day":
		startTime = now.AddDate(0, 0, -1)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	case "all":
		startTime = time.Time{}
	default:
		startTime = now.AddDate(0, 0, -7) // Default to week
	}
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	
	// Get click data
	var analytics []db.ClickAnalytics
	query := db.DB.Where("url_id = ?", url.ID)
	
	if !startTime.IsZero() {
		query = query.Where("clicked_at >= ?", startTime)
	}
	
	// Count total records for pagination
	var total int64
	query.Model(&db.ClickAnalytics{}).Count(&total)
	
	// Apply pagination and sorting
	offset := (page - 1) * pageSize
	result := query.Order("clicked_at desc").Offset(offset).Limit(pageSize).Find(&analytics)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analytics data"})
		return
	}
	
	// Get aggregate data
	var browserStats []struct {
		Browser string
		Count   int
	}
	db.DB.Raw(`
		SELECT browser, COUNT(*) as count
		FROM click_analytics
		WHERE url_id = ?
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY browser
		ORDER BY count DESC
	`, url.ID, startTime, startTime.IsZero()).Scan(&browserStats)
	
	var deviceStats []struct {
		Device string
		Count  int
	}
	db.DB.Raw(`
		SELECT device, COUNT(*) as count
		FROM click_analytics
		WHERE url_id = ?
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY device
		ORDER BY count DESC
	`, url.ID, startTime, startTime.IsZero()).Scan(&deviceStats)
	
	var countryStats []struct {
		CountryCode string
		Count       int
	}
	db.DB.Raw(`
		SELECT country_code, COUNT(*) as count
		FROM click_analytics
		WHERE url_id = ?
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY country_code
		ORDER BY count DESC
	`, url.ID, startTime, startTime.IsZero()).Scan(&countryStats)
	
	var referrerStats []struct {
		Referrer string
		Count    int
	}
	db.DB.Raw(`
		SELECT referrer, COUNT(*) as count
		FROM click_analytics
		WHERE url_id = ?
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY referrer
		ORDER BY count DESC
	`, url.ID, startTime, startTime.IsZero()).Scan(&referrerStats)
	
	// Format time series data by day
	var dailyClicks []struct {
		Date  string
		Count int
	}
	db.DB.Raw(`
		SELECT strftime('%Y-%m-%d', clicked_at) as date, COUNT(*) as count
		FROM click_analytics
		WHERE url_id = ?
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY date
		ORDER BY date
	`, url.ID, startTime, startTime.IsZero()).Scan(&dailyClicks)
	
	c.JSON(http.StatusOK, gin.H{
		"urlDetails": gin.H{
			"id":          url.ID,
			"shortCode":   url.ShortCode,
			"originalUrl": url.OriginalURL,
			"createdAt":   url.CreatedAt,
			"expiresAt":   url.ExpiresAt,
			"totalClicks": url.Clicks,
		},
		"clicksData": analytics,
		"pagination": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"aggregates": gin.H{
			"browsers":   browserStats,
			"devices":    deviceStats,
			"countries":  countryStats,
			"referrers":  referrerStats,
			"dailyClicks": dailyClicks,
		},
	})
}

// GetUserAnalytics returns aggregate analytics for all user URLs
func GetUserAnalytics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You must be logged in"})
		return
	}
	
	// Get time period from query params
	period := c.DefaultQuery("period", "month")
	
	// Calculate start time based on period
	var startTime time.Time
	now := time.Now()
	
	switch period {
	case "day":
		startTime = now.AddDate(0, 0, -1)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	case "all":
		startTime = time.Time{}
	default:
		startTime = now.AddDate(0, -1, 0) // Default to month
	}
	
	// Get user's URLs
	var urls []db.URL
	if result := db.DB.Where("user_id = ?", userID).Find(&urls); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user URLs"})
		return
	}
	
	// Get URL IDs for analytics query
	var urlIDs []uint
	for _, url := range urls {
		urlIDs = append(urlIDs, url.ID)
	}
	
	if len(urlIDs) == 0 {
		// No URLs, return empty analytics
		c.JSON(http.StatusOK, gin.H{
			"totalClicks": 0,
			"topURLs": []struct{}{},
			"clicksOverTime": []struct{}{},
			"deviceStats": []struct{}{},
			"browserStats": []struct{}{},
		})
		return
	}
	
	// Get total clicks count
	var totalClicks int
	db.DB.Model(&db.ClickAnalytics{}).
		Where("url_id IN (?)", urlIDs).
		Where("clicked_at >= ? OR ? IS NULL", startTime, startTime.IsZero()).
		Count(&totalClicks)
	
	// Get top URLs by clicks
	type TopURL struct {
		ID         uint   `json:"id"`
		ShortCode  string `json:"shortCode"`
		OriginalURL string `json:"originalUrl"`
		Clicks     int    `json:"clicks"`
	}
	var topURLs []TopURL
	db.DB.Raw(`
		SELECT u.id, u.short_code, u.original_url, COUNT(ca.id) as clicks
		FROM urls u
		LEFT JOIN click_analytics ca ON u.id = ca.url_id
		WHERE u.user_id = ?
		AND (ca.clicked_at >= ? OR ? IS NULL OR ca.clicked_at IS NULL)
		GROUP BY u.id
		ORDER BY clicks DESC
		LIMIT 10
	`, userID, startTime, startTime.IsZero()).Scan(&topURLs)
	
	// Get clicks over time
	var clicksOverTime []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	db.DB.Raw(`
		SELECT strftime('%Y-%m-%d', clicked_at) as date, COUNT(*) as count
		FROM click_analytics
		WHERE url_id IN (?)
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY date
		ORDER BY date
	`, urlIDs, startTime, startTime.IsZero()).Scan(&clicksOverTime)
	
	// Get device distribution
	var deviceStats []struct {
		Device string `json:"device"`
		Count  int    `json:"count"`
	}
	db.DB.Raw(`
		SELECT device, COUNT(*) as count
		FROM click_analytics
		WHERE url_id IN (?)
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY device
		ORDER BY count DESC
	`, urlIDs, startTime, startTime.IsZero()).Scan(&deviceStats)
	
	// Get browser distribution
	var browserStats []struct {
		Browser string `json:"browser"`
		Count   int    `json:"count"`
	}
	db.DB.Raw(`
		SELECT browser, COUNT(*) as count
		FROM click_analytics
		WHERE url_id IN (?)
		AND (clicked_at >= ? OR ? IS NULL)
		GROUP BY browser
		ORDER BY count DESC
	`, urlIDs, startTime, startTime.IsZero()).Scan(&browserStats)
	
	c.JSON(http.StatusOK, gin.H{
		"totalClicks":    totalClicks,
		"topURLs":        topURLs,
		"clicksOverTime": clicksOverTime,
		"deviceStats":    deviceStats,
		"browserStats":   browserStats,
	})
}

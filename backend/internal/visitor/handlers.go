package visitor

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers visitor tracking routes
func (s *Service) RegisterRoutes(router *gin.RouterGroup) {
	// Tracking endpoint for frontend page views
	router.POST("/track", s.handleTrackPageView)

	// Privacy consent endpoints
	privacy := router.Group("/privacy")
	{
		privacy.POST("/consent", s.handleRecordConsent)
		privacy.GET("/consent/:sessionId", s.handleGetConsent)
		privacy.DELETE("/data/:sessionId", s.handleDataErasure)
	}

	// Analytics endpoints (aggregated data only)
	analytics := router.Group("/analytics")
	{
		analytics.GET("/overview", s.handleGetOverview)
		analytics.GET("/realtime", s.handleGetRealtime)
		analytics.GET("/timeline", s.handleGetTimeline)
		analytics.GET("/locations", s.handleGetLocations)
		analytics.GET("/devices", s.handleGetDevices)
	}
}

// handleTrackPageView handles explicit page view tracking from frontend
func (s *Service) handleTrackPageView(c *gin.Context) {
	var req struct {
		Path      string `json:"path" binding:"required"`
		Referrer  string `json:"referrer"`
		SessionID string `json:"sessionId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := s.TrackPageView(ctx, c.Request, req.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track page view"})
		return
	}

	realtimeCount := s.GetRealTimeCount(ctx)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"activeVisitors": realtimeCount,
	})
}

// handleRecordConsent records user consent preferences
func (s *Service) handleRecordConsent(c *gin.Context) {
	var req struct {
		SessionHash string            `json:"sessionHash"`
		Consents    map[string]bool   `json:"consents"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Record each consent type
	for consentType, granted := range req.Consents {
		if err := s.RecordConsent(c.Request.Context(), req.SessionHash, consentType, granted); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record consent"})
			return
		}
	}

	// Set consent cookie
	c.SetCookie(
		"privacy_consent",
		"granted",
		86400*365, // 1 year
		"/",
		"",
		true,  // Secure
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// handleGetConsent gets consent status for a session
func (s *Service) handleGetConsent(c *gin.Context) {
	sessionID := c.Param("sessionId")
	
	status, err := s.GetConsentStatus(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consent not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// handleDataErasure handles GDPR right to erasure requests
func (s *Service) handleDataErasure(c *gin.Context) {
	sessionID := c.Param("sessionId")
	
	// Delete all data associated with the session
	if err := s.EraseSessionData(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to erase data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All data associated with this session has been erased",
	})
}

// handleGetOverview returns visitor statistics overview
func (s *Service) handleGetOverview(c *gin.Context) {
	stats, err := s.GetVisitorStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// handleGetRealtime returns real-time visitor data
func (s *Service) handleGetRealtime(c *gin.Context) {
	count := s.GetRealTimeCount(c.Request.Context())
	
	// Get current active pages
	var activePages []struct {
		Page  string `json:"page"`
		Count int    `json:"count"`
	}
	
	s.db.Raw(`
		SELECT current_page as page, COUNT(*) as count
		FROM visitor_realtime
		WHERE last_activity > ?
		GROUP BY current_page
		ORDER BY count DESC
		LIMIT 10
	`, time.Now().Add(-5*time.Minute)).Scan(&activePages)

	c.JSON(http.StatusOK, gin.H{
		"activeVisitors": count,
		"activePages":    activePages,
		"timestamp":      time.Now(),
	})
}

// handleGetTimeline returns visitor timeline data
func (s *Service) handleGetTimeline(c *gin.Context) {
	period := c.DefaultQuery("period", "7d")
	interval := c.DefaultQuery("interval", "day")
	
	data, err := s.GetTimelineData(c.Request.Context(), period, interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get timeline data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period":   period,
		"interval": interval,
		"data":     data,
	})
}

// handleGetLocations returns visitor location distribution
func (s *Service) handleGetLocations(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")
	
	locations, err := s.GetLocationDistribution(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get location data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"period":    period,
		"locations": locations,
	})
}

// handleGetDevices returns device and browser statistics
func (s *Service) handleGetDevices(c *gin.Context) {
	period := c.DefaultQuery("period", "30d")
	
	stats, err := s.GetDeviceBrowserStats(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// EraseSessionData erases all data for a session (GDPR compliance)
func (s *Service) EraseSessionData(ctx context.Context, sessionHash string) error {
	// Start transaction
	tx := s.db.Begin()
	
	// Delete page views
	if err := tx.Exec(`
		DELETE FROM page_views 
		WHERE session_id IN (
			SELECT id FROM visitor_sessions WHERE session_hash = ?
		)
	`, sessionHash).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Delete session
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&VisitorSession{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Delete consents
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&PrivacyConsent{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Delete from realtime
	if err := tx.Where("session_hash = ?", sessionHash).Delete(&VisitorRealtime{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Commit transaction
	return tx.Commit().Error
}
package visitor

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TrackingMiddleware creates a middleware for visitor tracking
func TrackingMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipTracking(c.Request.URL.Path) {
			c.Next()
			return
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			path := c.Request.URL.Path
			if err := service.TrackPageView(ctx, c.Request, path); err != nil {
				log.Printf("[VISITOR TRACKING ERROR] Path: %s | Error: %v", path, err)

				if service.metrics != nil {
					service.metrics.RecordLatency(0, "/visitor/error")
				}
			}
		}()

		c.Next()
	}
}

// ConsentMiddleware handles privacy consent
func ConsentMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for consent cookie or header
		consentCookie, err := c.Cookie("privacy_consent")
		if err == nil && consentCookie != "" {
			// Store consent status in context
			c.Set("privacy_consent", consentCookie)
		}

		c.Next()
	}
}

// shouldSkipTracking determines if a path should be tracked
func shouldSkipTracking(path string) bool {
	// Skip static assets
	staticExtensions := []string{".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".woff", ".woff2", ".ttf", ".eot"}
	for _, ext := range staticExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	// Skip API endpoints (except analytics endpoints)
	if strings.HasPrefix(path, "/api/") && !strings.Contains(path, "/analytics") {
		return true
	}

	// Skip health checks
	if path == "/health" || path == "/ping" {
		return true
	}

	// Skip metrics endpoint
	if path == "/metrics" {
		return true
	}

	return false
}

// RateLimitMiddleware prevents tracking abuse
func RateLimitMiddleware(maxRequestsPerMinute int) gin.HandlerFunc {
	// Simple in-memory rate limiting
	// In production, use Redis-based rate limiting
	requests := make(map[string][]time.Time)
	
	return func(c *gin.Context) {
		clientID := getClientIdentifier(c.Request)
		now := time.Now()
		
		// Clean old entries
		if times, exists := requests[clientID]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					validTimes = append(validTimes, t)
				}
			}
			requests[clientID] = validTimes
		}
		
		// Check rate limit
		if len(requests[clientID]) >= maxRequestsPerMinute {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}
		
		// Add current request
		requests[clientID] = append(requests[clientID], now)
		
		c.Next()
	}
}

// getClientIdentifier creates a client identifier for rate limiting
func getClientIdentifier(r *http.Request) string {
	// Use a combination of IP and User-Agent
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}
	
	ua := r.UserAgent()
	return ip + "|" + ua
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")
		
		// Other security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		c.Next()
	}
}

// PrivacyComplianceMiddleware ensures privacy compliance
func PrivacyComplianceMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user's location from IP (for compliance detection)
		// This is done in memory only, IP is not stored
		location := service.getLocationFromIP(c.Request)
		
		if location != nil {
			// Set compliance requirements based on location
			switch location.CountryCode {
			case "DE", "FR", "IT", "ES", "NL", "BE", "AT", "PL": // EU countries
				c.Set("privacy_regime", "GDPR")
				c.Set("require_consent", true)
			case "US":
				// Check for California
				if location.Region == "CA" {
					c.Set("privacy_regime", "CCPA")
					c.Set("allow_opt_out", true)
				}
			case "BR":
				c.Set("privacy_regime", "LGPD")
				c.Set("require_consent", true)
			case "CA":
				c.Set("privacy_regime", "PIPEDA")
				c.Set("require_consent", false)
			default:
				c.Set("privacy_regime", "DEFAULT")
				c.Set("require_consent", false)
			}
		}
		
		c.Next()
	}
}
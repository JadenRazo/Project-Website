package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
	"github.com/gin-gonic/gin"
)

// RateLimiter represents a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string]*userRequests
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type userRequests struct {
	count       int
	windowStart time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*userRequests),
		rate:     rate,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address)
		clientIP := c.ClientIP()
		requestID := c.GetString("RequestID")

		if rl.isAllowed(clientIP) {
			c.Next()
		} else {
			logger.Warn("Rate limit exceeded",
				"client_ip", clientIP,
				"request_id", requestID,
				"path", c.Request.URL.Path,
			)

			response.SendError(c, 429,
				fmt.Sprintf("Rate limit exceeded. Maximum %d requests per %v allowed.", rl.rate, rl.window),
				nil,
			)
			c.Abort()
		}
	}
}

// isAllowed checks if a request from the given IP is allowed
func (rl *RateLimiter) isAllowed(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get or create user request data
	ur, exists := rl.requests[clientIP]
	if !exists {
		rl.requests[clientIP] = &userRequests{
			count:       1,
			windowStart: now,
		}
		return true
	}

	// Check if window has expired
	if now.Sub(ur.windowStart) > rl.window {
		// Reset window
		ur.count = 1
		ur.windowStart = now
		return true
	}

	// Check if under rate limit
	if ur.count < rl.rate {
		ur.count++
		return true
	}

	return false
}

// cleanup periodically removes old entries from the map
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window * 2)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, ur := range rl.requests {
			if now.Sub(ur.windowStart) > rl.window*2 {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// APIRateLimiter returns a rate limiter suitable for API endpoints
func APIRateLimiter() gin.HandlerFunc {
	// 100 requests per minute
	return NewRateLimiter(100, time.Minute).Middleware()
}

// StrictRateLimiter returns a stricter rate limiter for sensitive endpoints
func StrictRateLimiter() gin.HandlerFunc {
	// 10 requests per minute
	return NewRateLimiter(10, time.Minute).Middleware()
}

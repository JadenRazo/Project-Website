package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// Define rate limiter structure
type RateLimiter struct {
	ipCache  *cache.Cache
	mu       sync.Mutex
	limit    int
	duration time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		ipCache:  cache.New(5*time.Minute, 10*time.Minute),
		limit:    limit,
		duration: duration,
	}
}

// RateLimitMiddleware limits request rate by IP address
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Lock for thread safety
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Check if IP exists in cache
		val, found := rl.ipCache.Get(ip)
		if !found {
			// First request from this IP
			rl.ipCache.Set(ip, 1, rl.duration)
		} else {
			// Increment request count
			count := val.(int)
			if count >= rl.limit {
				// Too many requests
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded. Please try again later.",
				})
				c.Abort()
				return
			}
			rl.ipCache.Set(ip, count+1, rl.duration)
		}

		c.Next()
	}
}

package security

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mutex    sync.RWMutex
	config   RateLimitConfig
	cleanup  *time.Ticker
}

type ClientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}

	if config.Enabled {
		rl.cleanup = time.NewTicker(config.CleanupInterval)
		go rl.cleanupRoutine()
	}

	return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		clientIP := getClientIP(c)
		limiter := rl.getLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"code":        "RATE_LIMIT_EXCEEDED",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		clientIP := getClientIP(c)
		limiter := rl.getAuthLimiter(clientIP)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Authentication rate limit exceeded",
				"code":        "AUTH_RATE_LIMIT_EXCEEDED",
				"retry_after": "300s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) getLimiter(clientIP string) *rate.Limiter {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	limiter, exists := rl.limiters[clientIP]
	if !exists {
		rps := rate.Limit(float64(rl.config.RequestsPerMin) / 60.0)
		limiter = rate.NewLimiter(rps, rl.config.BurstSize)
		rl.limiters[clientIP] = limiter
	}

	return limiter
}

func (rl *RateLimiter) getAuthLimiter(clientIP string) *rate.Limiter {
	key := "auth:" + clientIP
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		rps := rate.Limit(float64(rl.config.AuthEndpoints.RequestsPerMin) / 60.0)
		limiter = rate.NewLimiter(rps, rl.config.AuthEndpoints.BurstSize)
		rl.limiters[key] = limiter
	}

	return limiter
}

func (rl *RateLimiter) cleanupRoutine() {
	for range rl.cleanup.C {
		rl.mutex.Lock()

		// Remove limiters that haven't been used recently
		cutoff := time.Now().Add(-time.Hour)
		for key, limiter := range rl.limiters {
			// Check if the limiter has tokens available (indicating recent use)
			if limiter.TokensAt(cutoff) >= float64(rl.config.BurstSize) {
				delete(rl.limiters, key)
			}
		}

		rl.mutex.Unlock()
	}
}

func (rl *RateLimiter) Stop() {
	if rl.cleanup != nil {
		rl.cleanup.Stop()
	}
}

func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return map[string]interface{}{
		"active_limiters":  len(rl.limiters),
		"enabled":          rl.config.Enabled,
		"requests_per_min": rl.config.RequestsPerMin,
		"burst_size":       rl.config.BurstSize,
	}
}

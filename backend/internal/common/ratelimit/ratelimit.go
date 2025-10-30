package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/response"
	"github.com/go-redis/redis/v8"
)

// RateLimiter interface defines the methods that all rate limiter implementations must provide
type RateLimiter interface {
	// Allow checks if a request should be allowed based on the rate limiting rules
	Allow(ctx context.Context, key string) (bool, error)

	// Limit creates middleware that rate limits HTTP requests for a specific group
	Limit(group string, next http.Handler) http.Handler

	// Close cleans up any resources used by the rate limiter
	Close() error
}

// RedisRateLimiter implements RateLimiter using Redis
type RedisRateLimiter struct {
	client *redis.Client
	window time.Duration
	max    int64
	groups map[string]int64
}

// MemoryRateLimiter implements RateLimiter using in-memory storage
type MemoryRateLimiter struct {
	store  map[string][]time.Time
	window time.Duration
	max    int64
	groups map[string]int64
	mu     sync.RWMutex
}

// NewRateLimiter creates a new rate limiter instance based on configuration
func NewRateLimiter(cfg *config.RateLimitConfig) (RateLimiter, error) {
	window := time.Duration(cfg.WindowSeconds) * time.Second
	max := int64(cfg.MaxRequests)

	// Define rate limits for different groups
	groups := map[string]int64{
		"public": max,
		"api":    max * 2, // Higher limit for authenticated API users
		"admin":  max * 5, // Even higher limit for admin users
	}

	switch cfg.Type {
	case "redis":
		return newRedisRateLimiter(cfg, window, max, groups)
	case "memory":
		return newMemoryRateLimiter(window, max, groups), nil
	default:
		return nil, fmt.Errorf("unsupported rate limiter type: %s", cfg.Type)
	}
}

// newRedisRateLimiter creates a new Redis rate limiter instance
func newRedisRateLimiter(cfg *config.RateLimitConfig, window time.Duration, max int64, groups map[string]int64) (*RedisRateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisRateLimiter{
		client: client,
		window: window,
		max:    max,
		groups: groups,
	}, nil
}

// newMemoryRateLimiter creates a new in-memory rate limiter instance
func newMemoryRateLimiter(window time.Duration, max int64, groups map[string]int64) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		store:  make(map[string][]time.Time),
		window: window,
		max:    max,
		groups: groups,
		mu:     sync.RWMutex{},
	}
}

// Allow checks if a request should be allowed based on Redis rate limiting
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-r.window)

	// Clean up old entries
	pipe := r.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
	pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now.UnixNano()),
		Member: now.UnixNano(),
	})
	pipe.Expire(ctx, key, r.window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := results[1].(*redis.IntCmd).Val()
	return count < r.max, nil
}

// Limit creates HTTP middleware that rate limits requests for a specific group
func (r *RedisRateLimiter) Limit(group string, next http.Handler) http.Handler {
	// Get group limit or use default
	limit := r.max
	if groupLimit, exists := r.groups[group]; exists {
		limit = groupLimit
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Create a key that includes the group
		key := fmt.Sprintf("%s:%s", group, req.RemoteAddr)

		// Check if request should be allowed
		allowed, err := r.Allow(req.Context(), key)
		if err != nil {
			response.InternalError(w, "Internal Server Error")
			return
		}

		if !allowed {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))
			response.TooManyRequests(w, "Too Many Requests")
			return
		}

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		// We can't know exactly how many requests are remaining without another Redis query
		w.Header().Set("X-RateLimit-Remaining", "unknown")
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))

		next.ServeHTTP(w, req)
	})
}

// Close closes the Redis connection
func (r *RedisRateLimiter) Close() error {
	return r.client.Close()
}

// Allow checks if a request should be allowed based on memory rate limiting
func (r *MemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Clean up old entries
	if times, exists := r.store[key]; exists {
		valid := make([]time.Time, 0)
		for _, t := range times {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		r.store[key] = valid
	}

	// Check if we're under the limit
	if len(r.store[key]) >= int(r.max) {
		return false, nil
	}

	// Add new entry
	r.store[key] = append(r.store[key], now)
	return true, nil
}

// Limit creates HTTP middleware that rate limits requests for a specific group
func (r *MemoryRateLimiter) Limit(group string, next http.Handler) http.Handler {
	// Get group limit or use default
	limit := r.max
	if groupLimit, exists := r.groups[group]; exists {
		limit = groupLimit
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Create a key that includes the group
		key := fmt.Sprintf("%s:%s", group, req.RemoteAddr)

		// Check if request should be allowed
		allowed, err := r.Allow(req.Context(), key)
		if err != nil {
			response.InternalError(w, "Internal Server Error")
			return
		}

		if !allowed {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))
			response.TooManyRequests(w, "Too Many Requests")
			return
		}

		// Add rate limit headers and calculate remaining
		r.mu.RLock()
		remaining := int(limit) - len(r.store[key])
		if remaining < 0 {
			remaining = 0
		}
		r.mu.RUnlock()

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))

		next.ServeHTTP(w, req)
	})
}

// Close is a no-op for memory rate limiter
func (r *MemoryRateLimiter) Close() error {
	return nil
}

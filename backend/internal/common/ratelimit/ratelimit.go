package ratelimit

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/JadenRazo/Project-Website/backend/internal/app/config"
)

// RateLimiter interface defines the methods that all rate limiter implementations must provide
type RateLimiter interface {
    Allow(ctx context.Context, key string) (bool, error)
    Close() error
}

// RedisRateLimiter implements RateLimiter using Redis
type RedisRateLimiter struct {
    client *redis.Client
    window time.Duration
    max    int64
}

// MemoryRateLimiter implements RateLimiter using in-memory storage
type MemoryRateLimiter struct {
    store  map[string][]time.Time
    window time.Duration
    max    int64
    mu     sync.RWMutex
}

// NewRateLimiter creates a new rate limiter instance based on configuration
func NewRateLimiter(cfg *config.RateLimitConfig) (RateLimiter, error) {
    window := time.Duration(cfg.WindowSeconds) * time.Second
    max := int64(cfg.MaxRequests)

    switch cfg.Type {
    case "redis":
        return newRedisRateLimiter(cfg, window, max)
    case "memory":
        return newMemoryRateLimiter(window, max), nil
    default:
        return nil, fmt.Errorf("unsupported rate limiter type: %s", cfg.Type)
    }
}

// newRedisRateLimiter creates a new Redis rate limiter instance
func newRedisRateLimiter(cfg *config.RateLimitConfig, window time.Duration, max int64) (*RedisRateLimiter, error) {
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
    }, nil
}

// newMemoryRateLimiter creates a new in-memory rate limiter instance
func newMemoryRateLimiter(window time.Duration, max int64) *MemoryRateLimiter {
    return &MemoryRateLimiter{
        store:  make(map[string][]time.Time),
        window: window,
        max:    max,
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

// Close is a no-op for memory rate limiter
func (r *MemoryRateLimiter) Close() error {
    return nil
}

// RateLimitMiddleware creates a middleware that rate limits HTTP requests
func RateLimitMiddleware(limiter RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Generate rate limit key (can be customized based on IP, user, etc.)
            key := r.RemoteAddr

            // Check if request should be allowed
            allowed, err := limiter.Allow(r.Context(), key)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            if !allowed {
                w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", r.max))
                w.Header().Set("X-RateLimit-Remaining", "0")
                w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }

            // Add rate limit headers
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", r.max))
            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", r.max-1))
            w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(r.window).Unix()))

            next.ServeHTTP(w, r)
        })
    }
} 
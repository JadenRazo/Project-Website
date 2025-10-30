package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/go-redis/redis/v8"
)

// Cache interface defines the methods that all cache implementations must provide
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Close() error
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// MemoryCache implements Cache interface using in-memory storage
type MemoryCache struct {
	store map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewCache creates a new cache instance based on configuration
func NewCache(cfg *config.CacheConfig) (Cache, error) {
	switch cfg.Type {
	case "redis":
		return newRedisCache(cfg)
	case "memory":
		return newMemoryCache(), nil
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cfg.Type)
	}
}

// newRedisCache creates a new Redis cache instance
func newRedisCache(cfg *config.CacheConfig) (*RedisCache, error) {
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

	return &RedisCache{client: client}, nil
}

// newMemoryCache creates a new in-memory cache instance
func newMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]cacheItem),
	}
}

// Get retrieves a value from Redis cache
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set stores a value in Redis cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Delete removes a value from Redis cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Clear removes all values from Redis cache
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// Get retrieves a value from memory cache
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	item, exists := c.store[key]
	if !exists {
		return nil, nil
	}

	if time.Now().After(item.expiration) {
		delete(c.store, key)
		return nil, nil
	}

	return item.value, nil
}

// Set stores a value in memory cache
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.store[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

// Delete removes a value from memory cache
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.store, key)
	return nil
}

// Clear removes all values from memory cache
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.store = make(map[string]cacheItem)
	return nil
}

// Close is a no-op for memory cache
func (c *MemoryCache) Close() error {
	return nil
}

// CacheMiddleware creates a middleware that caches HTTP responses
func CacheMiddleware(cache Cache, ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip non-GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Generate cache key
			key := fmt.Sprintf("%s:%s", r.Method, r.URL.String())

			// Try to get from cache
			ctx := r.Context()
			if cached, err := cache.Get(ctx, key); err == nil && cached != nil {
				// Write cached response
				w.Header().Set("X-Cache", "HIT")
				w.Write(cached.([]byte))
				return
			}

			// Create response recorder
			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)

			// Cache the response
			if recorder.Code == http.StatusOK {
				cache.Set(ctx, key, recorder.Body.Bytes(), ttl)
			}

			// Copy response to actual writer
			w.Header().Set("X-Cache", "MISS")
			for k, v := range recorder.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(recorder.Code)
			w.Write(recorder.Body.Bytes())
		})
	}
}

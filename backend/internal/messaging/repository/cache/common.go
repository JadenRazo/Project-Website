package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Common cache constants
const (
	// Key prefixes
	messageKeyPrefix    = "messaging:message:"
	channelKeyPrefix    = "messaging:channel:"
	channelMessagesKey  = "messaging:channel:%d:messages"
	threadMessagesKey   = "messaging:thread:%d:messages"
	messageReactionsKey = "messaging:message:%d:reactions"
	channelMembersKey   = "messaging:channel:%d:members"
	userChannelsKey     = "messaging:user:%d:channels"
	pinnedMessagesKey   = "messaging:channel:%d:pinned"

	// Cache TTLs
	messageCacheTTL   = 30 * time.Minute
	channelCacheTTL   = 30 * time.Minute
	threadCacheTTL    = 5 * time.Minute
	reactionsCacheTTL = 15 * time.Minute
	membersCacheTTL   = 15 * time.Minute
	pinnedCacheTTL    = 15 * time.Minute
)

// BaseCache provides common caching functionality
type BaseCache struct {
	client redis.UniversalClient
}

// NewBaseCache creates a new base cache instance
func NewBaseCache(client redis.UniversalClient) *BaseCache {
	return &BaseCache{
		client: client,
	}
}

// getFromCache retrieves a value from cache and unmarshals it
func (b *BaseCache) getFromCache(ctx context.Context, key string, value interface{}) error {
	data, err := b.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// setInCache marshals and stores a value in cache
func (b *BaseCache) setInCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return b.client.Set(ctx, key, data, ttl).Err()
}

// deleteFromCache removes a key from cache
func (b *BaseCache) deleteFromCache(ctx context.Context, key string) error {
	return b.client.Del(ctx, key).Err()
}

// deleteFromCache removes multiple keys from cache
func (b *BaseCache) deleteFromCacheMulti(ctx context.Context, keys []string) error {
	return b.client.Del(ctx, keys...).Err()
}

// invalidateCache invalidates all cache entries for a given prefix
func (b *BaseCache) invalidateCache(ctx context.Context, prefix string) error {
	pattern := fmt.Sprintf("%s*", prefix)
	keys, err := b.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return b.deleteFromCacheMulti(ctx, keys)
	}
	return nil
}

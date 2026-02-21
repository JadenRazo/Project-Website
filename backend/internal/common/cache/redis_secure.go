package cache

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
)

// SecureRedisConfig holds configuration for secure Redis connection
type SecureRedisConfig struct {
	SocketPath    string
	NetworkAddr   string // Fallback to network address if socket unavailable
	Password      string
	Username      string
	MaxRetries    int
	TLSConfig     *tls.Config
	EncryptionKey []byte // 32 bytes for AES-256
}

// SecureRedisClient wraps redis client with security features
type SecureRedisClient struct {
	client *redis.Client
	config SecureRedisConfig
}

// NewSecureRedisClient creates a new secure Redis client
func NewSecureRedisClient(cfg SecureRedisConfig) (*SecureRedisClient, error) {
	fmt.Println("Creating new secure Redis client...")
	var client *redis.Client

	// Validate encryption key
	if len(cfg.EncryptionKey) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes for AES-256")
	}
	fmt.Println("Encryption key validated.")

	// Prefer Unix socket over network
	if cfg.SocketPath != "" {
		client = redis.NewClient(&redis.Options{
			Network:      "unix",
			Addr:         cfg.SocketPath,
			Username:     cfg.Username,
			Password:     cfg.Password,
			DB:           0,
			MaxRetries:   cfg.MaxRetries,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,

			// Custom dialer for additional security
			Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if network != "unix" {
					return nil, fmt.Errorf("only Unix socket connections allowed")
				}
				return net.Dial(network, addr)
			},
		})
	} else if cfg.NetworkAddr != "" {
		// Fallback to network connection with strict validation
		client = redis.NewClient(&redis.Options{
			Addr:         cfg.NetworkAddr,
			Username:     cfg.Username,
			Password:     cfg.Password,
			DB:           0,
			MaxRetries:   cfg.MaxRetries,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,
			TLSConfig:    cfg.TLSConfig,

			Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				ip := net.ParseIP(host)
				if ip != nil {
					if ip.IsLoopback() || ip.IsPrivate() {
						return net.Dial(network, addr)
					}
					return nil, fmt.Errorf("only private/loopback connections allowed, got: %s", host)
				}
				resolved, err := net.LookupIP(host)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve host %s: %w", host, err)
				}
				for _, rip := range resolved {
					if rip.IsLoopback() || rip.IsPrivate() {
						return net.Dial(network, addr)
					}
				}
				return nil, fmt.Errorf("only private/loopback connections allowed, got: %s", host)
			},
		})
	} else {
		return nil, fmt.Errorf("either SocketPath or NetworkAddr must be provided")
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &SecureRedisClient{
		client: client,
		config: cfg,
	}, nil
}

// SecureCache provides encrypted caching functionality
type SecureCache struct {
	redis *SecureRedisClient
	gcm   cipher.AEAD
}

// NewSecureCache creates a new secure cache instance
func NewSecureCache(redisClient *SecureRedisClient) (*SecureCache, error) {
	// Create AES cipher
	block, err := aes.NewCipher(redisClient.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &SecureCache{
		redis: redisClient,
		gcm:   gcm,
	}, nil
}

// encrypt encrypts data using AES-GCM
func (c *SecureCache) encrypt(plaintext []byte) (string, error) {
	// Create nonce
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := c.gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts data using AES-GCM
func (c *SecureCache) decrypt(ciphertext string) ([]byte, error) {
	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Extract nonce
	if len(data) < c.gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:c.gcm.NonceSize()], data[c.gcm.NonceSize():]

	// Decrypt
	plaintext, err := c.gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (c *SecureCache) Get(ctx context.Context, key string) (interface{}, error) {
	var result interface{}
	err := c.GetDecrypted(ctx, key, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *SecureCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.SetEncrypted(ctx, key, value, ttl)
}

// SetEncrypted stores encrypted value in cache
func (c *SecureCache) SetEncrypted(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Serialize value
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Encrypt
	encrypted, err := c.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	// Store in Redis
	return c.redis.client.Set(ctx, key, encrypted, ttl).Err()
}

// GetDecrypted retrieves and decrypts value from cache
func (c *SecureCache) GetDecrypted(ctx context.Context, key string, dest interface{}) error {
	// Get from Redis
	encrypted, err := c.redis.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	// Decrypt
	decrypted, err := c.decrypt(encrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}

	// Deserialize
	if err := json.Unmarshal(decrypted, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes key from cache
func (c *SecureCache) Delete(ctx context.Context, key string) error {
	return c.redis.client.Del(ctx, key).Err()
}

// Exists checks if key exists
func (c *SecureCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.redis.client.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX sets value only if key doesn't exist
func (c *SecureCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	// Serialize value
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Encrypt
	encrypted, err := c.encrypt(data)
	if err != nil {
		return false, fmt.Errorf("failed to encrypt value: %w", err)
	}

	// Store in Redis
	return c.redis.client.SetNX(ctx, key, encrypted, ttl).Result()
}

// IncrBy increments a counter (stores unencrypted for performance)
func (c *SecureCache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.redis.client.IncrBy(ctx, key, value).Result()
}

// Expire sets expiration on a key
func (c *SecureCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.redis.client.Expire(ctx, key, ttl).Err()
}

// TTL gets remaining TTL of a key
func (c *SecureCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.redis.client.TTL(ctx, key).Result()
}

// Clear clears the cache
func (c *SecureCache) Clear(ctx context.Context) error {
	return c.redis.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (c *SecureCache) Close() error {
	return c.redis.client.Close()
}

// Ping checks Redis connection
func (c *SecureCache) Ping(ctx context.Context) error {
	return c.redis.client.Ping(ctx).Err()
}

// FlushDB flushes the current database (use with caution)
func (c *SecureCache) FlushDB(ctx context.Context) error {
	// This command might be disabled in production
	return c.redis.client.FlushDB(ctx).Err()
}

// Info returns Redis server information
func (c *SecureCache) Info(ctx context.Context, section ...string) (string, error) {
	return c.redis.client.Info(ctx, section...).Result()
}
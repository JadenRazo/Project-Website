package config

import (
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App            AppConfig          `yaml:"app"`
	Server         ServerConfig       `yaml:"server"`
	Database       DatabaseConfig     `yaml:"database"`
	Cache          CacheConfig        `yaml:"cache"`
	Auth           AuthConfig         `yaml:"auth"`
	Logging        LoggingConfig      `yaml:"logging"`
	Metrics        MetricsConfig      `yaml:"metrics"`
	Tracing        TracingConfig      `yaml:"tracing"`
	RateLimit      RateLimitConfig    `yaml:"rateLimit"`
	Compression    CompressionConfig  `yaml:"compression"`
	Environment    string             `yaml:"environment"`
	AdminToken     string             `yaml:"adminToken"`
	MaxLogLines    int                `yaml:"maxLogLines"`
	LogRetention   time.Duration      `yaml:"logRetention"`
	AllowedOrigins string             `yaml:"allowedOrigins"`
	Port           string             `yaml:"port"`
	URLShortener   URLShortenerConfig `yaml:"urlShortener"`
	Messaging      MessagingConfig    `yaml:"messaging"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	Version        string   `yaml:"version"`
	BaseURL        string   `yaml:"baseUrl"`
	APIPrefix      string   `yaml:"apiPrefix"`
	Timezone       string   `yaml:"timezone"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
	MaxHeaderBytes  int           `yaml:"maxHeaderBytes"`
	TLSEnabled      bool          `yaml:"tlsEnabled"`
	TLSCert         string        `yaml:"tlsCert"`
	TLSKey          string        `yaml:"tlsKey"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver                 string `yaml:"driver"`
	DSN                    string `yaml:"dsn"`
	MaxIdleConns           int    `yaml:"maxIdleConns"`
	MaxOpenConns           int    `yaml:"maxOpenConns"`
	ConnMaxLifetimeMinutes int    `yaml:"connMaxLifetimeMinutes"`
	LogLevel               string `yaml:"logLevel"`
	SlowThresholdMs        int    `yaml:"slowThresholdMs"`
	AutoMigrate            bool   `yaml:"autoMigrate"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Type           string        `yaml:"type"`
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	Password       string        `yaml:"password"`
	DB             int           `yaml:"db"`
	DefaultTimeout time.Duration `yaml:"defaultTimeout"`
	Enabled        bool          `yaml:"enabled"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret      string        `yaml:"jwtSecret"`
	TokenExpiry    time.Duration `yaml:"tokenExpiry"`
	RefreshExpiry  time.Duration `yaml:"refreshExpiry"`
	Issuer         string        `yaml:"issuer"`
	Audience       string        `yaml:"audience"`
	CookieSecure   bool          `yaml:"cookieSecure"`
	CookieHTTPOnly bool          `yaml:"cookieHttpOnly"`
	CookieDomain   string        `yaml:"cookieDomain"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	TimeFormat string `yaml:"timeFormat"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Path        string `yaml:"path"`
	ServiceName string `yaml:"serviceName"`
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled          bool    `yaml:"enabled"`
	Endpoint         string  `yaml:"endpoint"`
	ServiceName      string  `yaml:"serviceName"`
	ServiceVersion   string  `yaml:"serviceVersion"`
	Environment      string  `yaml:"environment"`
	SamplingStrategy string  `yaml:"samplingStrategy"`
	SamplingRatio    float64 `yaml:"samplingRatio"`
}

// RateLimitConfig holds rate limit configuration
type RateLimitConfig struct {
	Type          string `yaml:"type"`
	WindowSeconds int    `yaml:"windowSeconds"`
	MaxRequests   int    `yaml:"maxRequests"`
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Password      string `yaml:"password"`
	DB            int    `yaml:"db"`
	Enabled       bool   `yaml:"enabled"`
}

// CompressionConfig holds compression configuration
type CompressionConfig struct {
	Enabled       bool     `yaml:"enabled"`
	Level         int      `yaml:"level"`
	MinSize       int      `yaml:"minSize"`
	Types         []string `yaml:"types"`
	BrotliQuality int      `yaml:"brotliQuality"`
}

// URLShortenerConfig holds configuration for URL shortener
type URLShortenerConfig struct {
	BaseURL         string `yaml:"baseUrl"`
	ShortCodeLength int    `yaml:"shortCodeLength"`
}

// MessagingConfig holds configuration for messaging service
type MessagingConfig struct {
	MaxMessageSize    int `yaml:"maxMessageSize"`
	MaxAttachmentSize int `yaml:"maxAttachmentSize"`
}

// GetConfig loads and returns the application configuration
func GetConfig() *Config {
	// This is a simplified version just to make imports work
	return &Config{
		Port:           "8080",
		AdminToken:     "admin",
		MaxLogLines:    1000,
		LogRetention:   7 * 24 * time.Hour,
		AllowedOrigins: "*",
		URLShortener: URLShortenerConfig{
			BaseURL:         "http://localhost:8080",
			ShortCodeLength: 6,
		},
		Messaging: MessagingConfig{
			MaxMessageSize:    4096,
			MaxAttachmentSize: 10 * 1024 * 1024,
		},
	}
}

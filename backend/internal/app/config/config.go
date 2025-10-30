package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	App         AppConfig         `yaml:"app"`
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig    `yaml:"database"`
	Cache       CacheConfig       `yaml:"cache"`
	Auth        AuthConfig        `yaml:"auth"`
	Logging     LoggingConfig     `yaml:"logging"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	Tracing     TracingConfig     `yaml:"tracing"`
	RateLimit   RateLimitConfig   `yaml:"rateLimit"`
	Compression CompressionConfig `yaml:"compression"`
	Environment string            `yaml:"environment"`
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

// LoadConfig loads configuration from files
func LoadConfig(path string) (*Config, error) {
	// Load environment variables from .env file if it exists
	envFile := filepath.Join(filepath.Dir(path), ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Load configuration from file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse configuration
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Set default values
	setDefaults(&cfg)

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(cfg *Config) {
	// App defaults
	if len(cfg.App.AllowedOrigins) == 0 {
		cfg.App.AllowedOrigins = []string{"https://jadenrazo.dev", "https://www.jadenrazo.dev"}
	}

	// Server defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10 * time.Second
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 60 * time.Second
	}
	if cfg.Server.ShutdownTimeout == 0 {
		cfg.Server.ShutdownTimeout = 5 * time.Second
	}
	if cfg.Server.MaxHeaderBytes == 0 {
		cfg.Server.MaxHeaderBytes = 1 << 20 // 1 MB
	}

	// Database defaults
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 20
	}
	if cfg.Database.ConnMaxLifetimeMinutes == 0 {
		cfg.Database.ConnMaxLifetimeMinutes = 60
	}
	if cfg.Database.SlowThresholdMs == 0 {
		cfg.Database.SlowThresholdMs = 200
	}

	// Auth defaults
	if cfg.Auth.TokenExpiry == 0 {
		cfg.Auth.TokenExpiry = 15 * time.Minute
	}
	if cfg.Auth.RefreshExpiry == 0 {
		cfg.Auth.RefreshExpiry = 7 * 24 * time.Hour
	}

	// Tracing defaults
	if cfg.Tracing.SamplingStrategy == "" {
		cfg.Tracing.SamplingStrategy = "ratio"
	}
	if cfg.Tracing.SamplingRatio == 0 {
		cfg.Tracing.SamplingRatio = 0.1
	}

	// Set default environment if not specified
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	// Set default app version if not specified
	if cfg.App.Version == "" {
		cfg.App.Version = "1.0.0"
	}
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	// Check required values
	if cfg.App.Name == "" {
		return fmt.Errorf("app.name must be set")
	}

	return nil
}

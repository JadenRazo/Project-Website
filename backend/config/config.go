package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Port              string        `yaml:"port"`
	JWTSecret         string        `yaml:"jwtSecret" sensitive:"true"`
	DatabasePath      string        `yaml:"databasePath"`
	AllowedOrigins    []string      `yaml:"allowedOrigins"`
	Environment       string        `yaml:"environment"`
	BaseURL           string        `yaml:"baseURL"`
	APIRateLimit      int           `yaml:"apiRateLimit"`
	RedirectRateLimit int           `yaml:"redirectRateLimit"`
	LogLevel          string        `yaml:"logLevel"`
	MetricsEnabled    bool          `yaml:"metricsEnabled"`
	TracingEnabled    bool          `yaml:"tracingEnabled"`
	ReadTimeout       time.Duration `yaml:"readTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdownTimeout"`
	TLSEnabled        bool          `yaml:"tlsEnabled"`
	TLSCertPath       string        `yaml:"tlsCertPath"`
	TLSKeyPath        string        `yaml:"tlsKeyPath"`
}

// LogFilter defines a filter for log entries
type LogFilter struct {
	Field    string `yaml:"field"`
	Operator string `yaml:"operator"`
	Value    string `yaml:"value"`
}

// Environmental variable constants
const (
	EnvPort              = "PORT"
	EnvJWTSecret         = "JWT_SECRET"
	EnvDatabasePath      = "DB_PATH"
	EnvAllowedOrigins    = "ALLOWED_ORIGINS"
	EnvEnvironment       = "ENVIRONMENT"
	EnvBaseURL           = "BASE_URL"
	EnvAPIRateLimit      = "API_RATE_LIMIT"
	EnvRedirectRateLimit = "REDIRECT_RATE_LIMIT"
	EnvLogLevel          = "LOG_LEVEL"
	EnvMetricsEnabled    = "METRICS_ENABLED"
	EnvTracingEnabled    = "TRACING_ENABLED"
	EnvReadTimeout       = "READ_TIMEOUT"
	EnvWriteTimeout      = "WRITE_TIMEOUT"
	EnvIdleTimeout       = "IDLE_TIMEOUT"
	EnvShutdownTimeout   = "SHUTDOWN_TIMEOUT"
	EnvTLSEnabled        = "TLS_ENABLED"
	EnvTLSCertPath       = "TLS_CERT_PATH"
	EnvTLSKeyPath        = "TLS_KEY_PATH"
	EnvConfigPath        = "CONFIG_PATH"
)

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	// Try to load .env file first
	_ = godotenv.Load(".env")

	// Determine environment
	env := getEnv(EnvEnvironment, "development")

	// Try to load environment-specific config file
	configPath := getEnv(EnvConfigPath, fmt.Sprintf("config/%s.yaml", env))
	config, err := loadConfigFromFile(configPath)
	if err != nil {
		// If no config file, use environment variables
		config = &Config{}
	}

	// Override with environment variables if they exist
	if os.Getenv(EnvPort) != "" {
		config.Port = os.Getenv(EnvPort)
	}
	if os.Getenv(EnvJWTSecret) != "" {
		config.JWTSecret = os.Getenv(EnvJWTSecret)
	}
	if os.Getenv(EnvDatabasePath) != "" {
		config.DatabasePath = os.Getenv(EnvDatabasePath)
	}
	if os.Getenv(EnvAllowedOrigins) != "" {
		config.AllowedOrigins = strings.Split(os.Getenv(EnvAllowedOrigins), ",")
	}
	if os.Getenv(EnvEnvironment) != "" {
		config.Environment = os.Getenv(EnvEnvironment)
	}
	if os.Getenv(EnvBaseURL) != "" {
		config.BaseURL = os.Getenv(EnvBaseURL)
	}
	if os.Getenv(EnvAPIRateLimit) != "" {
		if limit, err := strconv.Atoi(os.Getenv(EnvAPIRateLimit)); err == nil {
			config.APIRateLimit = limit
		}
	}
	if os.Getenv(EnvRedirectRateLimit) != "" {
		if limit, err := strconv.Atoi(os.Getenv(EnvRedirectRateLimit)); err == nil {
			config.RedirectRateLimit = limit
		}
	}
	if os.Getenv(EnvLogLevel) != "" {
		config.LogLevel = os.Getenv(EnvLogLevel)
	}
	if os.Getenv(EnvMetricsEnabled) != "" {
		if enabled, err := strconv.ParseBool(os.Getenv(EnvMetricsEnabled)); err == nil {
			config.MetricsEnabled = enabled
		}
	}
	if os.Getenv(EnvTracingEnabled) != "" {
		if enabled, err := strconv.ParseBool(os.Getenv(EnvTracingEnabled)); err == nil {
			config.TracingEnabled = enabled
		}
	}
	if os.Getenv(EnvReadTimeout) != "" {
		if duration, err := time.ParseDuration(os.Getenv(EnvReadTimeout)); err == nil {
			config.ReadTimeout = duration
		}
	}
	if os.Getenv(EnvWriteTimeout) != "" {
		if duration, err := time.ParseDuration(os.Getenv(EnvWriteTimeout)); err == nil {
			config.WriteTimeout = duration
		}
	}
	if os.Getenv(EnvIdleTimeout) != "" {
		if duration, err := time.ParseDuration(os.Getenv(EnvIdleTimeout)); err == nil {
			config.IdleTimeout = duration
		}
	}
	if os.Getenv(EnvShutdownTimeout) != "" {
		if duration, err := time.ParseDuration(os.Getenv(EnvShutdownTimeout)); err == nil {
			config.ShutdownTimeout = duration
		}
	}
	if os.Getenv(EnvTLSEnabled) != "" {
		if enabled, err := strconv.ParseBool(os.Getenv(EnvTLSEnabled)); err == nil {
			config.TLSEnabled = enabled
		}
	}
	if os.Getenv(EnvTLSCertPath) != "" {
		config.TLSCertPath = os.Getenv(EnvTLSCertPath)
	}
	if os.Getenv(EnvTLSKeyPath) != "" {
		config.TLSKeyPath = os.Getenv(EnvTLSKeyPath)
	}

	// Set defaults for any missing values
	setDefaults(config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// loadConfigFromFile loads configuration from YAML file
func loadConfigFromFile(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// Helper function to get environment variables
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// setDefaults sets default values for missing configuration
func setDefaults(config *Config) {
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.JWTSecret == "" {
		// Generate a warning but use default in development
		if config.Environment == "production" {
			fmt.Println("WARNING: No JWT secret provided in production environment!")
		}
		config.JWTSecret = "your-secret-key-change-in-production"
	}
	if config.DatabasePath == "" {
		config.DatabasePath = "data/app.db"
	}
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"http://localhost:3000"}
	}
	if config.Environment == "" {
		config.Environment = "development"
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:" + config.Port
	}
	if config.APIRateLimit == 0 {
		config.APIRateLimit = 100
	}
	if config.RedirectRateLimit == 0 {
		config.RedirectRateLimit = 1000
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 10 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 10 * time.Second
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = 60 * time.Second
	}
	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 5 * time.Second
	}

	// Create database directory if it doesn't exist
	dbDir := filepath.Dir(config.DatabasePath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		_ = os.MkdirAll(dbDir, 0755)
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Check required values
	if config.Port == "" {
		return errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return errors.New("JWT secret is required")
	}

	// In production, enforce stronger security requirements
	if config.Environment == "production" {
		// Ensure JWT secret is strong enough
		if len(config.JWTSecret) < 32 {
			return errors.New("JWT secret must be at least 32 characters in production")
		}

		// Ensure TLS is enabled
		if !config.TLSEnabled {
			return errors.New("TLS must be enabled in production")
		}

		// Verify TLS cert and key exist
		if config.TLSEnabled {
			if _, err := os.Stat(config.TLSCertPath); os.IsNotExist(err) {
				return fmt.Errorf("TLS certificate not found at: %s", config.TLSCertPath)
			}
			if _, err := os.Stat(config.TLSKeyPath); os.IsNotExist(err) {
				return fmt.Errorf("TLS key not found at: %s", config.TLSKeyPath)
			}
		}
	}

	return nil
}

// GetAPIRateLimiter returns the rate limit configuration for API endpoints
func (c *Config) GetAPIRateLimiter() (int, time.Duration) {
	return c.APIRateLimit, time.Minute
}

// GetRedirectRateLimiter returns the rate limit configuration for redirect endpoints
func (c *Config) GetRedirectRateLimiter() (int, time.Duration) {
	return c.RedirectRateLimit, time.Minute
}

// IsDevelopment returns true if the environment is set to development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is set to production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTest returns true if the environment is set to test
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

// GetSanitizedConfig returns a copy of the config with sensitive data masked
func (c *Config) GetSanitizedConfig() map[string]interface{} {
	return map[string]interface{}{
		"port":              c.Port,
		"jwtSecret":         "[REDACTED]",
		"databasePath":      c.DatabasePath,
		"allowedOrigins":    c.AllowedOrigins,
		"environment":       c.Environment,
		"baseURL":           c.BaseURL,
		"apiRateLimit":      c.APIRateLimit,
		"redirectRateLimit": c.RedirectRateLimit,
		"logLevel":          c.LogLevel,
		"metricsEnabled":    c.MetricsEnabled,
		"tracingEnabled":    c.TracingEnabled,
		"readTimeout":       c.ReadTimeout.String(),
		"writeTimeout":      c.WriteTimeout.String(),
		"idleTimeout":       c.IdleTimeout.String(),
		"shutdownTimeout":   c.ShutdownTimeout.String(),
		"tlsEnabled":        c.TLSEnabled,
		"tlsCertPath":       c.TLSCertPath,
		"tlsKeyPath":        c.TLSKeyPath,
	}
}

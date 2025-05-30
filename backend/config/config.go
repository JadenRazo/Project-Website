package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Port              string             `yaml:"port"`
	JWTSecret         string             `yaml:"jwtSecret" sensitive:"true"`
	Database          DatabaseConfig     `yaml:"database"`
	AllowedOrigins    []string           `yaml:"allowedOrigins"`
	Environment       string             `yaml:"environment"`
	BaseURL           string             `yaml:"baseURL"`
	APIRateLimit      int                `yaml:"apiRateLimit"`
	RedirectRateLimit int                `yaml:"redirectRateLimit"`
	LogLevel          string             `yaml:"logLevel"`
	MetricsEnabled    bool               `yaml:"metricsEnabled"`
	TracingEnabled    bool               `yaml:"tracingEnabled"`
	ReadTimeout       time.Duration      `yaml:"readTimeout"`
	WriteTimeout      time.Duration      `yaml:"writeTimeout"`
	IdleTimeout       time.Duration      `yaml:"idleTimeout"`
	ShutdownTimeout   time.Duration      `yaml:"shutdownTimeout"`
	TLSEnabled        bool               `yaml:"tlsEnabled"`
	TLSCert           string             `yaml:"tlsCert"`
	TLSKey            string             `yaml:"tlsKey"`
	EnablePprof       bool               `yaml:"enablePprof"`
	URLShortener      URLShortenerConfig `yaml:"urlShortener"`
	Messaging         MessagingConfig    `yaml:"messaging"`
	DevPanel          DevPanelConfig     `yaml:"devPanel"`
}

// DatabaseConfig holds PostgreSQL specific configuration
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user" sensitive:"true"`
	Password string `yaml:"password" sensitive:"true"`
	DBName   string `yaml:"dbName"`
	SSLMode  string `yaml:"sslMode"`
	DSN      string `yaml:"-"`
}

// URLShortenerConfig holds configuration for URL shortener service
type URLShortenerConfig struct {
	BaseURL         string `yaml:"baseURL"`         // e.g., http://localhost:8080 or your public domain
	ShortCodeLength int    `yaml:"shortCodeLength"` // Default length for generated short codes
	// MaxURLLength and MinURLLength are currently hardcoded in main.go
}

// MessagingConfig holds configuration for messaging service
type MessagingConfig struct {
	MaxMessageSize int `yaml:"maxMessageSize"` // Max size of a message in bytes
	// WebSocketPort, MaxAttachments, AllowedFileTypes are currently hardcoded in main.go
}

// DevPanelConfig holds configuration for the developer panel service
type DevPanelConfig struct {
	MetricsInterval string `yaml:"metricsInterval"` // e.g., "30s"
	MaxLogLines     int    `yaml:"maxLogLines"`
	LogRetention    string `yaml:"logRetention"` // e.g., "168h" (for 7 days)
}

// LogFilter defines a filter for log entries
type LogFilter struct {
	Field    string `yaml:"field"`
	Operator string `yaml:"operator"`
	Value    string `yaml:"value"`
}

// Environmental variable constants
const (
	EnvPort                    = "PORT"
	EnvJWTSecret               = "JWT_SECRET"
	EnvAllowedOrigins          = "ALLOWED_ORIGINS"
	EnvEnvironment             = "ENVIRONMENT"
	EnvBaseURL                 = "BASE_URL"
	EnvAPIRateLimit            = "API_RATE_LIMIT"
	EnvRedirectRateLimit       = "REDIRECT_RATE_LIMIT"
	EnvLogLevel                = "LOG_LEVEL"
	EnvMetricsEnabled          = "METRICS_ENABLED"
	EnvTracingEnabled          = "TRACING_ENABLED"
	EnvReadTimeout             = "READ_TIMEOUT"
	EnvWriteTimeout            = "WRITE_TIMEOUT"
	EnvIdleTimeout             = "IDLE_TIMEOUT"
	EnvShutdownTimeout         = "SHUTDOWN_TIMEOUT"
	EnvTLSEnabled              = "TLS_ENABLED"
	EnvTLSCertPath             = "TLS_CERT_PATH"
	EnvTLSKeyPath              = "TLS_KEY_PATH"
	EnvConfigPath              = "CONFIG_PATH"
	EnvEnablePprof             = "ENABLE_PPROF"
	EnvDBDriver                = "DB_DRIVER"
	EnvDBHost                  = "DB_HOST"
	EnvDBPort                  = "DB_PORT"
	EnvDBUser                  = "DB_USER"
	EnvDBPassword              = "DB_PASSWORD"
	EnvDBName                  = "DB_NAME"
	EnvDBSSLMode               = "DB_SSLMODE"
	EnvURLShortenerBaseURL     = "URL_SHORTENER_BASE_URL"
	EnvURLShortenerCodeLength  = "URL_SHORTENER_CODE_LENGTH"
	EnvMessagingMaxMessageSize = "MESSAGING_MAX_MESSAGE_SIZE"
	EnvDevPanelMetricsInterval = "DEV_PANEL_METRICS_INTERVAL"
	EnvDevPanelMaxLogLines     = "DEV_PANEL_MAX_LOG_LINES"
	EnvDevPanelLogRetention    = "DEV_PANEL_LOG_RETENTION"
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
		config.TLSCert = os.Getenv(EnvTLSCertPath)
	}
	if os.Getenv(EnvTLSKeyPath) != "" {
		config.TLSKey = os.Getenv(EnvTLSKeyPath)
	}
	if os.Getenv(EnvEnablePprof) != "" {
		if enabled, err := strconv.ParseBool(os.Getenv(EnvEnablePprof)); err == nil {
			config.EnablePprof = enabled
		}
	}

	// Load DatabaseConfig from environment variables
	if driver := os.Getenv(EnvDBDriver); driver != "" {
		config.Database.Driver = driver
	}
	if host := os.Getenv(EnvDBHost); host != "" {
		config.Database.Host = host
	}
	if portStr := os.Getenv(EnvDBPort); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Database.Port = port
		}
	}
	if user := os.Getenv(EnvDBUser); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv(EnvDBPassword); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv(EnvDBName); dbname != "" {
		config.Database.DBName = dbname
	}
	if sslmode := os.Getenv(EnvDBSSLMode); sslmode != "" {
		config.Database.SSLMode = sslmode
	}

	// Load URLShortenerConfig from environment variables
	if baseURL := os.Getenv(EnvURLShortenerBaseURL); baseURL != "" {
		config.URLShortener.BaseURL = baseURL
	}
	if lengthStr := os.Getenv(EnvURLShortenerCodeLength); lengthStr != "" {
		if length, err := strconv.Atoi(lengthStr); err == nil {
			config.URLShortener.ShortCodeLength = length
		}
	}

	// Load MessagingConfig from environment variables
	if sizeStr := os.Getenv(EnvMessagingMaxMessageSize); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil {
			config.Messaging.MaxMessageSize = size
		}
	}

	// Load DevPanelConfig from environment variables
	if interval := os.Getenv(EnvDevPanelMetricsInterval); interval != "" {
		config.DevPanel.MetricsInterval = interval
	}
	if linesStr := os.Getenv(EnvDevPanelMaxLogLines); linesStr != "" {
		if lines, err := strconv.Atoi(linesStr); err == nil {
			config.DevPanel.MaxLogLines = lines
		}
	}
	if retention := os.Getenv(EnvDevPanelLogRetention); retention != "" {
		config.DevPanel.LogRetention = retention
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

	// Set Database defaults
	if config.Database.Driver == "" {
		config.Database.Driver = "postgres"
	}
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.User == "" {
		// This should ideally cause an error in production if not set
		config.Database.User = "postgres"
		if config.Environment == "production" {
			fmt.Println("WARNING: DB_USER not set, defaulting to 'postgres'. This is not recommended for production.")
		}
	}
	if config.Database.DBName == "" {
		config.Database.DBName = "project_website" // Default DB name
		if config.Environment == "production" {
			fmt.Println("WARNING: DB_NAME not set, defaulting to 'project_website'.")
		}
	}
	if config.Database.Password == "" && config.Environment == "production" {
		// It's critical this is set in prod. We won't default it here for prod.
		fmt.Println("CRITICAL WARNING: DB_PASSWORD is not set for production environment!")
	}
	if config.Database.SSLMode == "" {
		// Default to "disable" for local dev, "require" or "verify-full" for prod is better
		if config.Environment == "production" {
			config.Database.SSLMode = "require"
		} else {
			config.Database.SSLMode = "disable"
		}
	}

	// Construct DSN for PostgreSQL
	config.Database.DSN = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	// Set URLShortener defaults
	if config.URLShortener.BaseURL == "" {
		// It often makes sense for this to be the same as the main BaseURL
		config.URLShortener.BaseURL = config.BaseURL
	}
	if config.URLShortener.ShortCodeLength == 0 {
		config.URLShortener.ShortCodeLength = 7 // Default length
	}

	// Set Messaging defaults
	if config.Messaging.MaxMessageSize == 0 {
		config.Messaging.MaxMessageSize = 4096 // Default 4KB
	}

	// Set DevPanel defaults
	if config.DevPanel.MetricsInterval == "" {
		config.DevPanel.MetricsInterval = "30s"
	}
	if config.DevPanel.MaxLogLines == 0 {
		config.DevPanel.MaxLogLines = 1000
	}
	if config.DevPanel.LogRetention == "" {
		config.DevPanel.LogRetention = "168h" // 7 days
	}

	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"http://localhost:3000"}
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.Port == "" {
		return errors.New("port is required")
	}
	if config.JWTSecret == "" || (config.Environment == "production" && config.JWTSecret == "your-secret-key-change-in-production") {
		return errors.New("JWT secret is required and must be changed for production")
	}
	// Database validation
	if config.Database.Driver == "" {
		return errors.New("database driver is required")
	}
	if config.Database.Driver != "postgres" {
		return fmt.Errorf("unsupported database driver: %s. Only 'postgres' is supported", config.Database.Driver)
	}
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port == 0 {
		return errors.New("database port is required")
	}
	if config.Database.User == "" {
		return errors.New("database user is required")
	}
	if config.Database.DBName == "" {
		return errors.New("database name is required")
	}
	if config.Environment == "production" && config.Database.Password == "" {
		// For production, password is a must
		return errors.New("database password is required for production environment")
	}

	if config.Environment == "production" && (config.Database.SSLMode != "require" && config.Database.SSLMode != "verify-full" && config.Database.SSLMode != "verify-ca") {
		// Forcing secure SSL for prod. "disable" is not allowed.
		// "allow" is also risky.
		fmt.Printf("WARNING: Insecure database SSLMode ('%s') for production. Recommended: 'require', 'verify-ca', or 'verify-full'.\n", config.Database.SSLMode)

	}

	if config.BaseURL == "" {
		return errors.New("base URL is required")
	}
	// Add more validation rules as needed
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
		"database":          c.Database,
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
		"tlsCert":           c.TLSCert,
		"tlsKey":            c.TLSKey,
		"enablePprof":       c.EnablePprof,
		"urlShortener":      c.URLShortener,
		"messaging":         c.Messaging,
		"devPanel":          c.DevPanel,
	}
}

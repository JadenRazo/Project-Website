package config

import (
	"log"
	"strconv"
	"strings"
	"time"

	loaderConfig "github.com/JadenRazo/Project-Website/backend/config" // Import the actual config loader
)

// Config holds all configuration for the application
// This struct should mirror the structure that main.go expects.
type Config struct {
	App            AppConfig
	Server         ServerConfig
	Database       DatabaseConfig
	Cache          CacheConfig
	Redis          RedisConfig
	Auth           AuthConfig
	Logging        LoggingConfig
	Metrics        MetricsConfig
	Tracing        TracingConfig
	RateLimit      RateLimitConfig
	Compression    CompressionConfig
	Environment    string
	AdminToken     string        // For DevPanel, maps to JWTSecret
	MaxLogLines    int           // For main logger, from DevPanelConfig.MaxLogLines (or default)
	LogRetention   time.Duration // For main logger, from DevPanelConfig.LogRetention (or default)
	AllowedOrigins string
	Port           string
	URLShortener   URLShortenerConfig // Populated from loaderConfig.URLShortener
	Messaging      MessagingConfig    // Populated from loaderConfig.Messaging
	DevPanel       DevPanelCoreConfig // Specific core config for DevPanel service params
	Contact        ContactConfig
	EnablePprof    bool
}

// RedisConfig holds redis configuration
type RedisConfig struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	Password      string `yaml:"password"`
	EncryptionKey string `yaml:"encryptionKey"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name           string
	Description    string
	Version        string
	BaseURL        string // Populated from loaderConfig.BaseURL
	APIPrefix      string
	Timezone       string
	AllowedOrigins []string // Populated from loaderConfig.AllowedOrigins
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration // Populated from loaderConfig.ReadTimeout
	WriteTimeout    time.Duration // Populated from loaderConfig.WriteTimeout
	IdleTimeout     time.Duration // Populated from loaderConfig.IdleTimeout
	ShutdownTimeout time.Duration // Populated from loaderConfig.ShutdownTimeout
	MaxHeaderBytes  int
	TLSEnabled      bool   // Populated from loaderConfig.TLSEnabled
	TLSCert         string // Populated from loaderConfig.TLSCert
	TLSKey          string // Populated from loaderConfig.TLSKey
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver                 string // Populated from loaderConfig.Database.Driver
	DSN                    string // Populated from loaderConfig.Database.DSN
	Host                   string
	Port                   int
	User                   string
	Password               string
	DBName                 string
	SSLMode                string
	MaxIdleConns           int
	MaxOpenConns           int
	ConnMaxLifetimeMinutes int
	LogLevel               string // This could be loaderConfig.LogLevel or specific DB log level
	SlowThresholdMs        int
	AutoMigrate            bool
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
	MaxAge     int    `yaml:"maxAge"` // This is int in core, but string (duration) in loader for DevPanel.LogRetention
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

// URLShortenerConfig matches the loader's definition
type URLShortenerConfig struct {
	BaseURL         string `yaml:"baseURL"`
	ShortCodeLength int    `yaml:"shortCodeLength"`
}

// MessagingConfig matches the loader's definition
type MessagingConfig struct {
	MaxMessageSize int `yaml:"maxMessageSize"`
}

// DevPanelCoreConfig holds parameters for DevPanel service initialization in main.go
type DevPanelCoreConfig struct {
	MetricsInterval time.Duration // Parsed from loaderConfig.DevPanel.MetricsInterval
	MaxLogLines     int           // From loaderConfig.DevPanel.MaxLogLines
	LogRetention    time.Duration // Parsed from loaderConfig.DevPanel.LogRetention
}

type ContactConfig struct {
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPass       string
	FromEmail      string
	ContactToEmail string
}

// GetConfig loads and returns the application configuration
// using the actual loader and maps it to the local core.Config struct.
func GetConfig() *Config {
	cfgFromLoader, err := loaderConfig.LoadConfig()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to load configuration: %v", err)
	}

	// --- Parse durations for DevPanel from strings ---
	devPanelMetricsInterval, errPI := time.ParseDuration(cfgFromLoader.DevPanel.MetricsInterval)
	if errPI != nil {
		log.Printf("WARNING: Could not parse DevPanel MetricsInterval '%s', using default 30s. Error: %v", cfgFromLoader.DevPanel.MetricsInterval, errPI)
		devPanelMetricsInterval = 30 * time.Second
	}
	devPanelLogRetention, errLR := time.ParseDuration(cfgFromLoader.DevPanel.LogRetention)
	if errLR != nil {
		log.Printf("WARNING: Could not parse DevPanel LogRetention '%s', using default 168h. Error: %v", cfgFromLoader.DevPanel.LogRetention, errLR)
		devPanelLogRetention = 168 * time.Hour
	}

	// Attempt to parse Port from loaderConfig.Port (string) to int for coreConfig.Server.Port
	serverPortInt := 8080 // Default
	if cfgFromLoader.Port != "" {
		p, err := strconv.Atoi(cfgFromLoader.Port)
		if err == nil {
			serverPortInt = p
		} else {
			log.Printf("WARNING: Could not parse Server Port '%s' to int, using default %d. Error: %v", cfgFromLoader.Port, serverPortInt, err)
		}
	}

	coreCfg := &Config{
		Environment:    cfgFromLoader.Environment,
		Port:           cfgFromLoader.Port,      // For server Addr string in main.go
		AdminToken:     cfgFromLoader.JWTSecret, // For DevPanel service init
		AllowedOrigins: strings.Join(cfgFromLoader.AllowedOrigins, ","),
		EnablePprof:    cfgFromLoader.EnablePprof,
		MaxLogLines:    cfgFromLoader.DevPanel.MaxLogLines, // Main logger uses DevPanel's MaxLogLines
		LogRetention:   devPanelLogRetention,               // Main logger uses DevPanel's LogRetention

		App: AppConfig{
			BaseURL:        cfgFromLoader.BaseURL,
			AllowedOrigins: cfgFromLoader.AllowedOrigins,
		},
		Server: ServerConfig{
			// Host: not directly in top-level loaderConfig.Config, usually "" or "0.0.0.0"
			Port:            serverPortInt, // Parsed int for consistency if needed by other parts
			ReadTimeout:     cfgFromLoader.ReadTimeout,
			WriteTimeout:    cfgFromLoader.WriteTimeout,
			IdleTimeout:     cfgFromLoader.IdleTimeout,
			ShutdownTimeout: cfgFromLoader.ShutdownTimeout,
			TLSEnabled:      cfgFromLoader.TLSEnabled,
			TLSCert:         cfgFromLoader.TLSCert,
			TLSKey:          cfgFromLoader.TLSKey,
		},
		Database: DatabaseConfig{
			Driver:   cfgFromLoader.Database.Driver,
			DSN:      cfgFromLoader.Database.DSN,
			Host:     cfgFromLoader.Database.Host,
			Port:     cfgFromLoader.Database.Port,
			User:     cfgFromLoader.Database.User,
			Password: cfgFromLoader.Database.Password,
			DBName:   cfgFromLoader.Database.DBName,
			SSLMode:  cfgFromLoader.Database.SSLMode,
			LogLevel: cfgFromLoader.LogLevel, // Main LogLevel used for DB too for now
		},
		Redis: RedisConfig{
			Host:          cfgFromLoader.Redis.Host,
			Port:          cfgFromLoader.Redis.Port,
			Password:      cfgFromLoader.Redis.Password,
			EncryptionKey: cfgFromLoader.Redis.EncryptionKey,
		},
		Logging: LoggingConfig{
			Level:      cfgFromLoader.LogLevel,
			Format:     "json",
			Output:     "stdout",
			TimeFormat: time.RFC3339Nano,                       // Default
			Filename:   "logs/api.log",                         // Default
			MaxSize:    100,                                    // Default
			MaxBackups: 10,                                     // Default
			MaxAge:     int(devPanelLogRetention.Hours() / 24), // MaxAge in days for logger from retention
			Compress:   true,                                   // Default
		},
		URLShortener: URLShortenerConfig{
			BaseURL:         cfgFromLoader.URLShortener.BaseURL,
			ShortCodeLength: cfgFromLoader.URLShortener.ShortCodeLength,
		},
		Messaging: MessagingConfig{
			MaxMessageSize: cfgFromLoader.Messaging.MaxMessageSize,
		},
		DevPanel: DevPanelCoreConfig{
			MetricsInterval: devPanelMetricsInterval,
			MaxLogLines:     cfgFromLoader.DevPanel.MaxLogLines,
			LogRetention:    devPanelLogRetention,
		},
		Contact: ContactConfig{
			SMTPHost:       cfgFromLoader.Contact.SMTPHost,
			SMTPPort:       cfgFromLoader.Contact.SMTPPort,
			SMTPUser:       cfgFromLoader.Contact.SMTPUser,
			SMTPPass:       cfgFromLoader.Contact.SMTPPass,
			FromEmail:      cfgFromLoader.Contact.FromEmail,
			ContactToEmail: cfgFromLoader.Contact.ContactToEmail,
		},
	}

	return coreCfg
}

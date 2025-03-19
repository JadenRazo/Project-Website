package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the messaging service
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Websocket WebsocketConfig `yaml:"websocket"`
	Database  DatabaseConfig  `yaml:"database"`
	Auth      AuthConfig      `yaml:"auth"`
	Cache     CacheConfig     `yaml:"cache"`
	Logging   LoggingConfig   `yaml:"logging"`
	Features  FeaturesConfig  `yaml:"features"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	BasePath    string `yaml:"basePath"`
	TLSEnabled  bool   `yaml:"tlsEnabled"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
}

// WebsocketConfig holds WebSocket server configuration
type WebsocketConfig struct {
	Path            string `yaml:"path"`
	MaxConnections  int    `yaml:"maxConnections"`
	PingInterval    int    `yaml:"pingInterval"`
	PongTimeout     int    `yaml:"pongTimeout"`
	WriteBufferSize int    `yaml:"writeBufferSize"`
	ReadBufferSize  int    `yaml:"readBufferSize"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Path     string `yaml:"path"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	SSLMode  string `yaml:"sslMode"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled     bool     `yaml:"enabled"`
	SecretKey   string   `yaml:"secretKey"`
	TokenExpiry int      `yaml:"tokenExpiry"`
	AdminUsers  []string `yaml:"adminUsers"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	Enabled         bool `yaml:"enabled"`
	DefaultExpiry   int  `yaml:"defaultExpiry"`
	CleanupInterval int  `yaml:"cleanupInterval"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// FeaturesConfig holds feature flags
type FeaturesConfig struct {
	ReadReceipts bool `yaml:"readReceipts"`
	Typing       bool `yaml:"typing"`
	FileUploads  bool `yaml:"fileUploads"`
	Embeds       bool `yaml:"embeds"`
	Reactions    bool `yaml:"reactions"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(path string) (*Config, error) {
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	// Read the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse the YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Set defaults for missing values
	setDefaults(config)

	return config, nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Host:     "0.0.0.0",
			Port:     8082,
			BasePath: "/api/messaging",
		},
		Websocket: WebsocketConfig{
			Path:            "/ws",
			MaxConnections:  100,
			PingInterval:    30,
			PongTimeout:     10,
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: "data/messaging.db",
		},
		Auth: AuthConfig{
			Enabled:     true,
			SecretKey:   "default-insecure-key",
			TokenExpiry: 24,
		},
		Cache: CacheConfig{
			Enabled:         true,
			DefaultExpiry:   5,
			CleanupInterval: 10,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Features: FeaturesConfig{
			ReadReceipts: true,
			Typing:       true,
			FileUploads:  false,
			Embeds:       true,
			Reactions:    true,
		},
	}

	return config
}

// setDefaults fills in default values for missing configuration
func setDefaults(config *Config) {
	if config.Server.Port == 0 {
		config.Server.Port = 8082
	}

	if config.Server.BasePath == "" {
		config.Server.BasePath = "/api/messaging"
	}

	if config.Websocket.Path == "" {
		config.Websocket.Path = "/ws"
	}

	if config.Websocket.MaxConnections == 0 {
		config.Websocket.MaxConnections = 100
	}

	if config.Websocket.PingInterval == 0 {
		config.Websocket.PingInterval = 30
	}

	if config.Websocket.PongTimeout == 0 {
		config.Websocket.PongTimeout = 10
	}

	if config.Websocket.WriteBufferSize == 0 {
		config.Websocket.WriteBufferSize = 1024
	}

	if config.Websocket.ReadBufferSize == 0 {
		config.Websocket.ReadBufferSize = 1024
	}

	if config.Database.Type == "" {
		config.Database.Type = "sqlite"
	}

	if config.Database.Type == "sqlite" && config.Database.Path == "" {
		config.Database.Path = "data/messaging.db"
	}

	if config.Auth.TokenExpiry == 0 {
		config.Auth.TokenExpiry = 24
	}

	if config.Cache.DefaultExpiry == 0 {
		config.Cache.DefaultExpiry = 5
	}

	if config.Cache.CleanupInterval == 0 {
		config.Cache.CleanupInterval = 10
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}

	if config.Logging.Format == "" {
		config.Logging.Format = "text"
	}

	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

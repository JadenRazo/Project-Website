package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port              string
	JWTSecret         string
	DatabasePath      string
	AllowedOrigins    []string
	Environment       string
	BaseURL           string
	APIRateLimit      int
	RedirectRateLimit int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	dbPath := getEnv("DB_PATH", "data/urls.db")
	origins := strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ",")
	env := getEnv("ENVIRONMENT", "development")
	baseURL := getEnv("BASE_URL", "http://localhost:"+port)
	
	// Rate limits
	apiLimit, _ := strconv.Atoi(getEnv("API_RATE_LIMIT", "100"))
	redirectLimit, _ := strconv.Atoi(getEnv("REDIRECT_RATE_LIMIT", "1000"))

	return Config{
		Port:              port,
		JWTSecret:         jwtSecret,
		DatabasePath:      dbPath,
		AllowedOrigins:    origins,
		Environment:       env,
		BaseURL:           baseURL,
		APIRateLimit:      apiLimit,
		RedirectRateLimit: redirectLimit,
	}
}

// Helper function to get environment variables
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
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

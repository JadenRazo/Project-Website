package config

import (
	"os"
	"strconv"
	"strings"
)

// DefaultEnvValues contains default configuration values
var DefaultEnvValues = map[string]string{
	// Server configuration
	"PORT":        "8080",
	"BASE_URL":    "http://localhost:8080",
	"ENVIRONMENT": "development",

	// Database configuration
	"DB_PATH": "data/urls.db",

	// Security
	"JWT_SECRET": "replace-this-with-a-long-secure-random-string-in-production",

	// CORS configuration
	"ALLOWED_ORIGINS": "http://localhost:3001,http://localhost:3000,http://localhost:8080",

	// Rate limiting
	"API_RATE_LIMIT":      "100",
	"REDIRECT_RATE_LIMIT": "1000",
}

// LoadEnv loads environment variables with defaults
func LoadEnv() {
	for key, defaultValue := range DefaultEnvValues {
		if os.Getenv(key) == "" {
			os.Setenv(key, defaultValue)
		}
	}
}

// GetEnvString gets an environment variable as string
func GetEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvInt gets an environment variable as integer
func GetEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetEnvBool gets an environment variable as boolean
func GetEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetEnvStringSlice gets an environment variable as string slice
func GetEnvStringSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	return strings.Split(valueStr, ",")
}

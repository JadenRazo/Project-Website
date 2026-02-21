package auth

import (
	"os"
	"strings"
	"time"
)

// GetJWTSecret retrieves the JWT secret from environment
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Use a secure default for development only
		if os.Getenv("APP_ENV") == "development" {
			return "dev-secret-key-change-in-production-minimum-32-chars"
		}
		// In production, this should fail
		panic("JWT_SECRET environment variable is required in production")
	}
	return secret
}

// GetAuthConfig returns the authentication configuration
func GetAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:     GetJWTSecret(),
		TokenExpiry:   15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
		Issuer:        "jadenrazo.dev",
		Audience:      "jadenrazo.dev",
	}
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret     string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	Issuer        string
	Audience      string
}

// OAuth2ProviderConfig holds configuration for a single OAuth2 provider
type OAuth2ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Enabled      bool
}

// OAuth2Config holds all OAuth2 provider configurations
type OAuth2Config struct {
	Google            OAuth2ProviderConfig
	GitHub            OAuth2ProviderConfig
	Microsoft         OAuth2ProviderConfig
	Discord           OAuth2ProviderConfig
	EncryptionKey     []byte
	AllowedEmails     []string
	StateExpiry       time.Duration
	TempTokenExpiry   time.Duration
}

// GetOAuth2Config returns the OAuth2 configuration from environment
func GetOAuth2Config() *OAuth2Config {
	encryptionKey := os.Getenv("OAUTH_ENCRYPTION_KEY")
	if len(encryptionKey) != 64 {
		if os.Getenv("APP_ENV") == "development" {
			encryptionKey = "0000000000000000000000000000000000000000000000000000000000000000"
		} else {
			panic("OAUTH_ENCRYPTION_KEY must be a 32-byte hex string (64 characters)")
		}
	}

	allowedEmails := strings.Split(os.Getenv("ADMIN_ALLOWED_EMAILS"), ",")
	for i, email := range allowedEmails {
		allowedEmails[i] = strings.TrimSpace(email)
	}

	return &OAuth2Config{
		Google: OAuth2ProviderConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  getEnvOrDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/admin/2fa/callback/google"),
			Scopes:       []string{"openid", "email", "profile"},
			Enabled:      os.Getenv("GOOGLE_CLIENT_ID") != "",
		},
		GitHub: OAuth2ProviderConfig{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  getEnvOrDefault("GITHUB_REDIRECT_URL", "http://localhost:8080/api/v1/auth/admin/2fa/callback/github"),
			Scopes:       []string{"read:user", "user:email"},
			Enabled:      os.Getenv("GITHUB_CLIENT_ID") != "",
		},
		Microsoft: OAuth2ProviderConfig{
			ClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
			ClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
			RedirectURL:  getEnvOrDefault("MICROSOFT_REDIRECT_URL", "http://localhost:8080/api/v1/auth/admin/2fa/callback/microsoft"),
			Scopes:       []string{"openid", "email", "profile"},
			Enabled:      os.Getenv("MICROSOFT_CLIENT_ID") != "",
		},
		Discord: OAuth2ProviderConfig{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  getEnvOrDefault("DISCORD_REDIRECT_URL", "http://localhost:8080/api/v1/discord/callback"),
			Scopes:       []string{"identify", "role_connections.write"},
			Enabled:      os.Getenv("DISCORD_CLIENT_ID") != "",
		},
		EncryptionKey:   []byte(encryptionKey),
		AllowedEmails:   allowedEmails,
		StateExpiry:     5 * time.Minute,
		TempTokenExpiry: 5 * time.Minute,
	}
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsEmailAllowed checks if an email is in the allowed list
func (c *OAuth2Config) IsEmailAllowed(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))

	for _, allowed := range c.AllowedEmails {
		allowed = strings.ToLower(strings.TrimSpace(allowed))

		// Exact match
		if email == allowed {
			return true
		}

		// Domain wildcard match (e.g., @jadenrazo.dev)
		if strings.HasPrefix(allowed, "@") {
			if strings.HasSuffix(email, allowed) {
				return true
			}
		}
	}

	return false
}
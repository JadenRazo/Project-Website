package oauth

import (
	"context"
)

// UserInfo represents the standardized user information from OAuth providers
type UserInfo struct {
	ProviderUserID string
	Email          string
	Name           string
	Picture        string
	EmailVerified  bool
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

// Provider defines the interface that all OAuth providers must implement
type Provider interface {
	// GetName returns the provider name (google, github, microsoft)
	GetName() string

	// GetAuthURL generates the OAuth authorization URL with state parameter
	GetAuthURL(state string) string

	// ExchangeCode exchanges the authorization code for access/refresh tokens
	ExchangeCode(ctx context.Context, code string) (*TokenResponse, error)

	// GetUserInfo fetches user information using the access token
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)

	// RefreshToken refreshes an expired access token using the refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// RevokeToken revokes an access token
	RevokeToken(ctx context.Context, token string) error
}

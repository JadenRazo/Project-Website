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

// AuthURLResponse contains the authorization URL and PKCE verifier
type AuthURLResponse struct {
	URL      string
	Verifier string // Empty if provider doesn't support PKCE (e.g., GitHub)
}

// Provider defines the interface that all OAuth providers must implement
type Provider interface {
	// GetName returns the provider name (google, github, microsoft)
	GetName() string

	// GetAuthURL generates the OAuth authorization URL with state and nonce parameters
	// Returns the URL and PKCE verifier (empty string if provider doesn't support PKCE)
	GetAuthURL(state, nonce string) (*AuthURLResponse, error)

	// ExchangeCode exchanges the authorization code for access/refresh tokens
	// verifier should be empty string if provider doesn't support PKCE
	ExchangeCode(ctx context.Context, code, verifier string) (*TokenResponse, error)

	// GetUserInfo fetches user information using the access token
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)

	// RefreshToken refreshes an expired access token using the refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// RevokeToken revokes an access token
	RevokeToken(ctx context.Context, token string) error
}

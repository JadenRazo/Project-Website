package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleProvider implements the Provider interface for Google OAuth2
type GoogleProvider struct {
	config *oauth2.Config
}

// NewGoogleProvider creates a new Google OAuth provider
func NewGoogleProvider(clientID, clientSecret, redirectURL string, scopes []string) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     google.Endpoint,
		},
	}
}

// GetName returns the provider name
func (p *GoogleProvider) GetName() string {
	return "google"
}

// GetAuthURL generates the OAuth authorization URL with PKCE and nonce
func (p *GoogleProvider) GetAuthURL(state, nonce string) (*AuthURLResponse, error) {
	verifier := oauth2.GenerateVerifier()

	url := p.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("nonce", nonce),
	)

	return &AuthURLResponse{
		URL:      url,
		Verifier: verifier,
	}, nil
}

// ExchangeCode exchanges the authorization code for tokens with PKCE verification
func (p *GoogleProvider) ExchangeCode(ctx context.Context, code, verifier string) (*TokenResponse, error) {
	var token *oauth2.Token
	var err error

	if verifier != "" {
		token, err = p.config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	} else {
		token, err = p.config.Exchange(ctx, code)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int64(time.Until(token.Expiry).Seconds()),
		TokenType:    token.TokenType,
	}, nil
}

// googleUserInfo represents the Google userinfo API response
type googleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// GetUserInfo fetches user information from Google
func (p *GoogleProvider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API returned status %d: %s", resp.StatusCode, string(body))
	}

	var googleUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &UserInfo{
		ProviderUserID: googleUser.Sub,
		Email:          googleUser.Email,
		Name:           googleUser.Name,
		Picture:        googleUser.Picture,
		EmailVerified:  googleUser.EmailVerified,
	}, nil
}

// RefreshToken refreshes an expired access token
func (p *GoogleProvider) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := p.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresIn:    int64(time.Until(newToken.Expiry).Seconds()),
		TokenType:    newToken.TokenType,
	}, nil
}

// RevokeToken revokes an access token
func (p *GoogleProvider) RevokeToken(ctx context.Context, token string) error {
	url := fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create revoke request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("google revoke returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

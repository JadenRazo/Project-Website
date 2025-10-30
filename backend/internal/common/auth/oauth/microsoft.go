package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type MicrosoftProvider struct {
	config *oauth2.Config
}

func NewMicrosoftProvider(clientID, clientSecret, redirectURL string, scopes []string) *MicrosoftProvider {
	return &MicrosoftProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     microsoft.AzureADEndpoint("common"),
		},
	}
}

func (p *MicrosoftProvider) GetName() string {
	return "microsoft"
}

func (p *MicrosoftProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *MicrosoftProvider) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int64(token.Expiry.Sub(token.Expiry).Seconds()),
		TokenType:    token.TokenType,
	}, nil
}

type microsoftUserInfo struct {
	ID                string `json:"id"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	DisplayName       string `json:"displayName"`
}

func (p *MicrosoftProvider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://graph.microsoft.com/v1.0/me", nil)
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
		return nil, fmt.Errorf("microsoft API returned status %d: %s", resp.StatusCode, string(body))
	}

	var msUser microsoftUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	email := msUser.Mail
	if email == "" {
		email = msUser.UserPrincipalName
	}

	return &UserInfo{
		ProviderUserID: msUser.ID,
		Email:          email,
		Name:           msUser.DisplayName,
		Picture:        "",
		EmailVerified:  true,
	}, nil
}

func (p *MicrosoftProvider) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
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
		ExpiresIn:    int64(newToken.Expiry.Sub(newToken.Expiry).Seconds()),
		TokenType:    newToken.TokenType,
	}, nil
}

func (p *MicrosoftProvider) RevokeToken(ctx context.Context, token string) error {
	return nil
}

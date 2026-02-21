package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

func (p *MicrosoftProvider) GetAuthURL(state, nonce string) (*AuthURLResponse, error) {
	verifier := oauth2.GenerateVerifier()

	url := p.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("nonce", nonce),
	)

	return &AuthURLResponse{
		URL:      url,
		Verifier: verifier,
	}, nil
}

func (p *MicrosoftProvider) ExchangeCode(ctx context.Context, code, verifier string) (*TokenResponse, error) {
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
		ExpiresIn:    int64(time.Until(newToken.Expiry).Seconds()),
		TokenType:    newToken.TokenType,
	}, nil
}

func (p *MicrosoftProvider) RevokeToken(ctx context.Context, token string) error {
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/logout"

	formData := fmt.Sprintf("token=%s&token_type_hint=access_token&client_id=%s&client_secret=%s",
		token, p.config.ClientID, p.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(formData))
	if err != nil {
		return fmt.Errorf("failed to create revoke request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("microsoft revoke returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

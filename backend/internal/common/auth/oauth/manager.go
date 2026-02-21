package oauth

import (
	"fmt"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
)

type OAuthManager struct {
	providers map[string]Provider
	config    *auth.OAuth2Config
}

func NewOAuthManager(config *auth.OAuth2Config) *OAuthManager {
	manager := &OAuthManager{
		providers: make(map[string]Provider),
		config:    config,
	}

	if config.Google.Enabled {
		manager.providers["google"] = NewGoogleProvider(
			config.Google.ClientID,
			config.Google.ClientSecret,
			config.Google.RedirectURL,
			config.Google.Scopes,
		)
	}

	if config.GitHub.Enabled {
		manager.providers["github"] = NewGitHubProvider(
			config.GitHub.ClientID,
			config.GitHub.ClientSecret,
			config.GitHub.RedirectURL,
			config.GitHub.Scopes,
		)
	}

	if config.Microsoft.Enabled {
		manager.providers["microsoft"] = NewMicrosoftProvider(
			config.Microsoft.ClientID,
			config.Microsoft.ClientSecret,
			config.Microsoft.RedirectURL,
			config.Microsoft.Scopes,
		)
	}

	return manager
}

func (m *OAuthManager) GetProvider(name string) (Provider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found or not enabled", name)
	}
	return provider, nil
}

func (m *OAuthManager) GetEnabledProviders() []string {
	providers := make([]string, 0, len(m.providers))
	for name := range m.providers {
		providers = append(providers, name)
	}
	return providers
}

func (m *OAuthManager) ValidateProvider(name string) bool {
	_, exists := m.providers[name]
	return exists
}

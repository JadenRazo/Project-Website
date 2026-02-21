package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
)

// OAuthState represents the OAuth flow state stored in Redis
type OAuthState struct {
	State      string    `json:"state"`
	Nonce      string    `json:"nonce"`
	Verifier   string    `json:"verifier"`
	Provider   string    `json:"provider"`
	SessionID  string    `json:"session_id"`
	CreatedAt  time.Time `json:"created_at"`
	RedirectURL string   `json:"redirect_url,omitempty"`
}

// GenerateOAuthState creates a new OAuth state with cryptographically random values
func GenerateOAuthState(ctx context.Context, cache cache.Cache, provider, sessionID, redirectURL string) (*OAuthState, error) {
	stateBytes := make([]byte, 32)
	nonceBytes := make([]byte, 32)

	if _, err := rand.Read(stateBytes); err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}
	if _, err := rand.Read(nonceBytes); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	state := &OAuthState{
		State:       base64.URLEncoding.EncodeToString(stateBytes),
		Nonce:       base64.URLEncoding.EncodeToString(nonceBytes),
		Provider:    provider,
		SessionID:   sessionID,
		CreatedAt:   time.Now(),
		RedirectURL: redirectURL,
	}

	key := fmt.Sprintf("oauth:state:%s", state.State)
	data, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := cache.Set(ctx, key, string(data), 5*time.Minute); err != nil {
		return nil, fmt.Errorf("failed to store state: %w", err)
	}

	return state, nil
}

// ValidateOAuthState retrieves and validates the OAuth state from Redis
func ValidateOAuthState(ctx context.Context, cache cache.Cache, state string) (*OAuthState, error) {
	if state == "" {
		return nil, fmt.Errorf("state parameter is required")
	}

	key := fmt.Sprintf("oauth:state:%s", state)

	data, err := cache.Get(ctx, key)
	if err != nil || data == nil {
		return nil, fmt.Errorf("invalid or expired state")
	}

	dataStr, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("invalid state data format")
	}

	var storedState OAuthState
	if err := json.Unmarshal([]byte(dataStr), &storedState); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}

	if err := cache.Delete(ctx, key); err != nil {
		return nil, fmt.Errorf("failed to delete state: %w", err)
	}

	if time.Since(storedState.CreatedAt) > 5*time.Minute {
		return nil, fmt.Errorf("state expired")
	}

	return &storedState, nil
}

// StoreVerifier stores the PKCE verifier with the OAuth state
func (s *OAuthState) StoreVerifier(ctx context.Context, cache cache.Cache, verifier string) error {
	s.Verifier = verifier

	key := fmt.Sprintf("oauth:state:%s", s.State)
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := cache.Set(ctx, key, string(data), 5*time.Minute); err != nil {
		return fmt.Errorf("failed to update state with verifier: %w", err)
	}

	return nil
}

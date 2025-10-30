package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
)

func TestNewJWTManager(t *testing.T) {
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	assert.NotNil(t, manager)
	assert.Equal(t, []byte(cfg.JWTSecret), manager.secret)
	assert.Equal(t, cfg.TokenExpiry, manager.tokenExpiry)
	assert.Equal(t, cfg.RefreshExpiry, manager.refreshExpiry)
	assert.Equal(t, cfg.Issuer, manager.issuer)
	assert.Equal(t, cfg.Audience, manager.audience)
}

func TestJWTManager_GenerateToken(t *testing.T) {
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	tests := []struct {
		name    string
		userID  string
		role    string
		wantErr bool
	}{
		{
			name:    "valid user token",
			userID:  "user123",
			role:    "user",
			wantErr: false,
		},
		{
			name:    "valid admin token",
			userID:  "admin456",
			role:    "admin",
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			role:    "user",
			wantErr: false,
		},
		{
			name:    "empty role",
			userID:  "user123",
			role:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateToken(tt.userID, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Contains(t, token, ".")
			}
		})
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	validToken, err := manager.GenerateToken("user123", "user")
	require.NoError(t, err)

	expiredManager := &JWTManager{
		secret:        []byte("test-secret-key"),
		tokenExpiry:   -time.Hour,
		refreshExpiry: 24 * time.Hour,
		issuer:        "devpanel",
		audience:      "users",
	}

	expiredToken, err := expiredManager.GenerateToken("user123", "user")
	require.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		wantErr    bool
		wantClaims *CustomClaims
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "expired token",
			token:   expiredToken,
			wantErr: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "malformed JWT",
			token:   "not.a.jwt",
			wantErr: true,
		},
		{
			name:    "token signed with different secret",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := manager.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "user123", claims.UserID)
				assert.Equal(t, "user", claims.Role)
			}
		})
	}
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	t.Run("generate refresh token", func(t *testing.T) {
		refreshToken, err := manager.GenerateRefreshToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshToken)
		assert.True(t, len(refreshToken) > 40) // Base64 encoded 32 bytes should be > 40 chars
	})

	t.Run("refresh tokens are unique", func(t *testing.T) {
		token1, err1 := manager.GenerateRefreshToken()
		token2, err2 := manager.GenerateRefreshToken()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestExtractTokenFromBearerString(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		wantErr       bool
	}{
		{
			name:          "valid bearer token",
			authHeader:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr:       false,
		},
		{
			name:       "missing bearer prefix",
			authHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr:    true,
		},
		{
			name:       "wrong prefix",
			authHeader: "Basic dXNlcjpwYXNz",
			wantErr:    true,
		},
		{
			name:       "empty header",
			authHeader: "",
			wantErr:    true,
		},
		{
			name:       "bearer without token",
			authHeader: "Bearer",
			wantErr:    true,
		},
		{
			name:       "bearer with empty token",
			authHeader: "Bearer ",
			wantErr:    true,
		},
		{
			name:          "bearer with extra spaces",
			authHeader:    "Bearer   token123",
			expectedToken: "  token123",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromBearerString(tt.authHeader)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestJWTManager_TokenSecurityValidation(t *testing.T) {
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	t.Run("token_with_different_roles", func(t *testing.T) {
		userToken, err := manager.GenerateToken("user123", "user")
		require.NoError(t, err)

		adminToken, err := manager.GenerateToken("admin456", "admin")
		require.NoError(t, err)

		userClaims, err := manager.ValidateToken(userToken)
		require.NoError(t, err)
		assert.Equal(t, "user", userClaims.Role)

		adminClaims, err := manager.ValidateToken(adminToken)
		require.NoError(t, err)
		assert.Equal(t, "admin", adminClaims.Role)
	})

	t.Run("token_tampering_detection", func(t *testing.T) {
		validToken, err := manager.GenerateToken("user123", "user")
		require.NoError(t, err)

		tamperedToken := validToken[:len(validToken)-5] + "XXXXX"

		_, err = manager.ValidateToken(tamperedToken)
		assert.Error(t, err)
	})

	t.Run("cross_service_token_validation", func(t *testing.T) {
		cfg1 := &config.AuthConfig{
			JWTSecret:     "secret1",
			TokenExpiry:   time.Hour,
			RefreshExpiry: 24 * time.Hour,
			Issuer:        "devpanel",
			Audience:      "users",
		}

		cfg2 := &config.AuthConfig{
			JWTSecret:     "secret2",
			TokenExpiry:   time.Hour,
			RefreshExpiry: 24 * time.Hour,
			Issuer:        "devpanel",
			Audience:      "users",
		}

		manager1 := NewJWTManager(cfg1)
		manager2 := NewJWTManager(cfg2)

		token1, err := manager1.GenerateToken("user123", "user")
		require.NoError(t, err)

		_, err = manager2.ValidateToken(token1)
		assert.Error(t, err, "Token signed with different secret should be invalid")
	})
}

func TestJWTManager_EdgeCases(t *testing.T) {
	cfg := &config.AuthConfig{
		JWTSecret:     "test-secret-key",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	t.Run("very_long_user_id", func(t *testing.T) {
		longUserID := string(make([]byte, 1000))
		for i := range longUserID {
			longUserID = string(append([]byte(longUserID[:i]), 'a'))
		}

		token, err := manager.GenerateToken(longUserID, "user")
		assert.NoError(t, err)

		validatedClaims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, longUserID, validatedClaims.UserID)
	})

	t.Run("special_characters_in_claims", func(t *testing.T) {
		token, err := manager.GenerateToken("user@domain.com", "admin/super-user")
		assert.NoError(t, err)

		validatedClaims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "user@domain.com", validatedClaims.UserID)
		assert.Equal(t, "admin/super-user", validatedClaims.Role)
	})

	t.Run("token_at_exact_expiry", func(t *testing.T) {
		shortCfg := &config.AuthConfig{
			JWTSecret:     "test-secret-key",
			TokenExpiry:   time.Second,
			RefreshExpiry: 24 * time.Hour,
			Issuer:        "devpanel",
			Audience:      "users",
		}

		shortManager := NewJWTManager(shortCfg)

		token, err := shortManager.GenerateToken("user123", "user")
		require.NoError(t, err)

		time.Sleep(2 * time.Second)

		_, err = shortManager.ValidateToken(token)
		assert.Error(t, err, "Token should be expired")
	})
}

func BenchmarkJWTManager_GenerateToken(b *testing.B) {
	cfg := &config.AuthConfig{
		JWTSecret:     "benchmark-secret",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GenerateToken("user123", "user")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJWTManager_ValidateToken(b *testing.B) {
	cfg := &config.AuthConfig{
		JWTSecret:     "benchmark-secret",
		TokenExpiry:   time.Hour,
		RefreshExpiry: 24 * time.Hour,
		Issuer:        "devpanel",
		Audience:      "users",
	}

	manager := NewJWTManager(cfg)

	token, err := manager.GenerateToken("user123", "user")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}

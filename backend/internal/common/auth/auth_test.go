package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth/mock"
	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

func TestService_GenerateTokens(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	tests := []struct {
		name    string
		userID  string
		role    string
		wantErr bool
	}{
		{
			name:    "valid user tokens",
			userID:  "user123",
			role:    "user",
			wantErr: false,
		},
		{
			name:    "valid admin tokens",
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
			tokens, err := authService.GenerateTokens(tt.userID, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.Greater(t, tokens.ExpiresIn, int64(0))
			}
		})
	}
}

func TestService_Register(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	tests := []struct {
		name     string
		email    string
		password string
		username string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid registration",
			email:    "user@example.com",
			password: "ValidPass123!",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "weak password",
			email:    "user2@example.com",
			password: "weak",
			username: "testuser2",
			wantErr:  true,
			errMsg:   "invalid password",
		},
		{
			name:     "duplicate email",
			email:    "user@example.com",
			password: "ValidPass123!",
			username: "testuser3",
			wantErr:  true,
			errMsg:   "user already exists",
		},
		{
			name:     "empty email",
			email:    "",
			password: "ValidPass123!",
			username: "testuser4",
			wantErr:  false,
		},
		{
			name:     "empty username",
			email:    "user5@example.com",
			password: "ValidPass123!",
			username: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := authService.Register(context.Background(), tt.email, tt.password, tt.username)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.username, user.Username)
				assert.True(t, user.IsActive)
				assert.False(t, user.IsVerified)
				assert.NotEqual(t, tt.password, user.HashedPassword)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	testPassword := "ValidPass123!"
	hashedPassword, err := mock.HashPassword(testPassword)
	require.NoError(t, err)

	activeUser := &entity.User{
		Email:          "active@example.com",
		HashedPassword: hashedPassword,
		IsActive:       true,
		IsVerified:     true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	inactiveUser := &entity.User{
		Email:          "inactive@example.com",
		HashedPassword: hashedPassword,
		IsActive:       false,
		IsVerified:     true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = userRepo.Create(context.Background(), activeUser)
	require.NoError(t, err)
	err = userRepo.Create(context.Background(), inactiveUser)
	require.NoError(t, err)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid login",
			email:    "active@example.com",
			password: testPassword,
			wantErr:  false,
		},
		{
			name:     "wrong password",
			email:    "active@example.com",
			password: "WrongPass123!",
			wantErr:  true,
			errMsg:   "invalid credentials",
		},
		{
			name:     "nonexistent user",
			email:    "notfound@example.com",
			password: testPassword,
			wantErr:  true,
			errMsg:   "invalid credentials",
		},
		{
			name:     "inactive user",
			email:    "inactive@example.com",
			password: testPassword,
			wantErr:  true,
			errMsg:   "account is inactive",
		},
		{
			name:     "empty email",
			email:    "",
			password: testPassword,
			wantErr:  true,
			errMsg:   "invalid credentials",
		},
		{
			name:     "empty password",
			email:    "active@example.com",
			password: "",
			wantErr:  true,
			errMsg:   "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := authService.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tokens)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
				assert.Greater(t, tokens.ExpiresIn, int64(0))
			}
		})
	}
}

func TestService_ValidateToken(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	validClaims := &mock.Claims{
		UserID: "user123",
		Role:   "user",
		StandardClaims: mock.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel",
			Audience:  "users",
		},
	}

	validToken, err := jwtService.GenerateToken(validClaims)
	require.NoError(t, err)

	expiredClaims := &mock.Claims{
		UserID: "user123",
		Role:   "user",
		StandardClaims: mock.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
			Issuer:    "devpanel",
			Audience:  "users",
		},
	}

	expiredToken, err := jwtService.GenerateToken(expiredClaims)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
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
			name:    "invalid token",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, validClaims.UserID, claims.UserID)
				assert.Equal(t, validClaims.Role, claims.Role)
			}
		})
	}
}

func TestService_RefreshToken(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	refreshClaims := &mock.Claims{
		UserID: "user123",
		Role:   "refresh",
		StandardClaims: mock.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel",
			Audience:  "users",
		},
	}

	validRefreshToken, err := jwtService.GenerateToken(refreshClaims)
	require.NoError(t, err)

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
	}{
		{
			name:         "valid refresh token",
			refreshToken: validRefreshToken,
			wantErr:      false,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-token",
			wantErr:      true,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTokens, err := authService.RefreshToken(tt.refreshToken)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, newTokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, newTokens)
				assert.NotEmpty(t, newTokens.AccessToken)
				assert.NotEmpty(t, newTokens.RefreshToken)
				assert.NotEqual(t, tt.refreshToken, newTokens.AccessToken)
			}
		})
	}
}

func TestService_Logout(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	validClaims := &mock.Claims{
		UserID: "user123",
		Role:   "user",
		StandardClaims: mock.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "devpanel",
			Audience:  "users",
		},
	}

	token, err := jwtService.GenerateToken(validClaims)
	require.NoError(t, err)

	t.Run("successful logout", func(t *testing.T) {
		err := authService.Logout(token)
		assert.NoError(t, err)

		_, err = authService.ValidateToken(token)
		assert.Error(t, err, "Token should be blacklisted after logout")
	})

	t.Run("logout with invalid token", func(t *testing.T) {
		err := authService.Logout("invalid-token")
		assert.NoError(t, err)
	})
}

func TestService_SessionManagement(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	testPassword := "ValidPass123!"
	hashedPassword, err := mock.HashPassword(testPassword)
	require.NoError(t, err)

	user := &entity.User{
		Email:          "session@example.com",
		HashedPassword: hashedPassword,
		IsActive:       true,
		IsVerified:     true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	t.Run("session creation on login", func(t *testing.T) {
		tokens, err := authService.Login(context.Background(), user.Email, testPassword)
		require.NoError(t, err)

		sessionInfo, err := authService.GetSessionInfo(tokens.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, sessionInfo)
		assert.Equal(t, user.ID.String(), sessionInfo.UserID)
		assert.False(t, sessionInfo.LoginTime.IsZero())
		assert.False(t, sessionInfo.LastSeen.IsZero())
	})

	t.Run("session info for nonexistent token", func(t *testing.T) {
		_, err := authService.GetSessionInfo("nonexistent-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session not found")
	})
}

func TestService_RoleValidation(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	tests := []struct {
		name         string
		role         string
		expectedRole string
	}{
		{
			name:         "user role",
			role:         "user",
			expectedRole: "user",
		},
		{
			name:         "admin role",
			role:         "admin",
			expectedRole: "admin",
		},
		{
			name:         "moderator role",
			role:         "moderator",
			expectedRole: "moderator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &mock.Claims{
				UserID: "user123",
				Role:   tt.role,
				StandardClaims: mock.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
					IssuedAt:  time.Now().Unix(),
					Issuer:    "devpanel",
					Audience:  "users",
				},
			}

			token, err := jwtService.GenerateToken(claims)
			require.NoError(t, err)

			validatedClaims, err := authService.ValidateToken(token)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedRole, validatedClaims.Role)
		})
	}
}

func TestService_ConcurrentAccess(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	jwtService := mock.NewMockJWTService("test-secret")
	authService := mock.NewMockAuthService(userRepo, jwtService)

	testPassword := "ValidPass123!"
	hashedPassword, err := mock.HashPassword(testPassword)
	require.NoError(t, err)

	user := &entity.User{
		Email:          "concurrent@example.com",
		HashedPassword: hashedPassword,
		IsActive:       true,
		IsVerified:     true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	t.Run("concurrent login attempts", func(t *testing.T) {
		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := authService.Login(context.Background(), user.Email, testPassword)
				results <- err
			}()
		}

		for i := 0; i < numGoroutines; i++ {
			err := <-results
			assert.NoError(t, err)
		}
	})

	t.Run("concurrent token validation", func(t *testing.T) {
		tokens, err := authService.Login(context.Background(), user.Email, testPassword)
		require.NoError(t, err)

		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				_, err := authService.ValidateToken(tokens.AccessToken)
				results <- err
			}()
		}

		for i := 0; i < numGoroutines; i++ {
			err := <-results
			assert.NoError(t, err)
		}
	})
}

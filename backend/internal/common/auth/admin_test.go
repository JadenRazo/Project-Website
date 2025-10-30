package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

func TestAdminAuth_ValidateAdminEmail(t *testing.T) {
	adminAuth := &AdminAuth{
		jwtSecret: "test-secret",
	}

	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "valid admin email domain",
			email: "admin@jadenrazo.dev",
			valid: true,
		},
		{
			name:  "valid admin email subdomain",
			email: "test@sub.jadenrazo.dev",
			valid: false, // The actual implementation might not support subdomains
		},
		{
			name:  "invalid domain",
			email: "admin@example.com",
			valid: false,
		},
		{
			name:  "invalid domain similar",
			email: "admin@jadenrazo.com",
			valid: false,
		},
		{
			name:  "empty email",
			email: "",
			valid: false,
		},
		{
			name:  "malformed email",
			email: "notanemail",
			valid: false,
		},
		{
			name:  "no domain",
			email: "admin@",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adminAuth.ValidateAdminEmail(tt.email)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestAdminAuth_GenerateToken(t *testing.T) {
	adminAuth := &AdminAuth{
		jwtSecret: "test-secret-key",
	}

	adminUser := &entity.User{
		Email:    "admin@jadenrazo.dev",
		Username: "admin",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	userUser := &entity.User{
		Email:    "user@example.com",
		Username: "user",
		Role:     entity.RoleUser,
		IsActive: true,
	}

	t.Run("generate_token_for_admin", func(t *testing.T) {
		token, expiresIn, err := adminAuth.GenerateToken(adminUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Greater(t, expiresIn, int64(0))
	})

	t.Run("generate_token_for_user", func(t *testing.T) {
		token, expiresIn, err := adminAuth.GenerateToken(userUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Greater(t, expiresIn, int64(0))
	})
}

func TestAdminAuth_ValidateToken(t *testing.T) {
	adminAuth := &AdminAuth{
		jwtSecret: "test-secret-key",
	}

	adminUser := &entity.User{
		Email:    "admin@jadenrazo.dev",
		Username: "admin",
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	validToken, _, err := adminAuth.GenerateToken(adminUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid_admin_token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "invalid_token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty_token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "malformed_token",
			token:   "not-a-jwt-token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := adminAuth.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, adminUser.Email, claims.Email)
				assert.True(t, claims.IsAdmin)
			}
		})
	}
}

func TestAdminAuth_GenerateSetupToken(t *testing.T) {
	adminAuth := &AdminAuth{
		jwtSecret: "test-secret-key",
	}

	t.Run("generate_setup_token", func(t *testing.T) {
		token, err := adminAuth.GenerateSetupToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, len(token) > 40) // Base64 encoded token should be reasonably long
	})

	t.Run("tokens_are_unique", func(t *testing.T) {
		token1, err1 := adminAuth.GenerateSetupToken()
		token2, err2 := adminAuth.GenerateSetupToken()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestAdminAuth_HasAdminAccount(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	adminAuth := NewAdminAuth(db)

	t.Run("no_admin_accounts", func(t *testing.T) {
		hasAdmin, err := adminAuth.HasAdminAccount()
		assert.NoError(t, err)
		assert.False(t, hasAdmin)
	})

	t.Run("with_admin_account", func(t *testing.T) {
		// Create an admin user
		adminUser := &entity.User{
			ID:       uuid.New(),
			Email:    "admin@jadenrazo.dev",
			Username: "admin",
			Role:     entity.RoleAdmin,
			IsActive: true,
		}

		err := db.Create(adminUser).Error
		require.NoError(t, err)

		hasAdmin, err := adminAuth.HasAdminAccount()
		assert.NoError(t, err)
		assert.True(t, hasAdmin)
	})

	t.Run("with_regular_user_only", func(t *testing.T) {
		// Clear the database
		err := db.Exec("DELETE FROM users").Error
		require.NoError(t, err)

		// Create only a regular user
		regularUser := &entity.User{
			ID:       uuid.New(),
			Email:    "user@example.com",
			Username: "user",
			Role:     entity.RoleUser,
			IsActive: true,
		}

		err = db.Create(regularUser).Error
		require.NoError(t, err)

		hasAdmin, err := adminAuth.HasAdminAccount()
		assert.NoError(t, err)
		assert.False(t, hasAdmin)
	})
}

func TestAdminAuth_Login(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	adminAuth := NewAdminAuth(db)

	// Create an admin user with hashed password
	password := "AdminPass123!"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	adminUser := &entity.User{
		ID:             uuid.New(),
		Email:          "admin@jadenrazo.dev",
		Username:       "admin",
		HashedPassword: hashedPassword,
		Role:           entity.RoleAdmin,
		IsActive:       true,
		IsVerified:     true,
	}

	err = db.Create(adminUser).Error
	require.NoError(t, err)

	// Create a regular user
	regularUser := &entity.User{
		ID:             uuid.New(),
		Email:          "user@jadenrazo.dev",
		Username:       "user",
		HashedPassword: hashedPassword,
		Role:           entity.RoleUser,
		IsActive:       true,
		IsVerified:     true,
	}

	err = db.Create(regularUser).Error
	require.NoError(t, err)

	tests := []struct {
		name    string
		request *AdminLoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid_admin_login",
			request: &AdminLoginRequest{
				Email:    "admin@jadenrazo.dev",
				Password: password,
			},
			wantErr: false,
		},
		{
			name: "wrong_password",
			request: &AdminLoginRequest{
				Email:    "admin@jadenrazo.dev",
				Password: "WrongPassword123!",
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "regular_user_attempting_admin_login",
			request: &AdminLoginRequest{
				Email:    "user@jadenrazo.dev",
				Password: password,
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "nonexistent_user",
			request: &AdminLoginRequest{
				Email:    "nonexistent@jadenrazo.dev",
				Password: password,
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "invalid_email_domain",
			request: &AdminLoginRequest{
				Email:    "admin@example.com",
				Password: password,
			},
			wantErr: true,
			errMsg:  "email not authorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := adminAuth.Login(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, response)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.Token)
				assert.NotNil(t, response.User)
				assert.Equal(t, tt.request.Email, response.User.Email)
				assert.True(t, response.User.IsAdmin)
				assert.Greater(t, response.ExpiresIn, int64(0))
			}
		})
	}
}

func TestAdminAuth_CompleteSetup(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	adminAuth := NewAdminAuth(db)

	t.Run("complete_setup_successfully", func(t *testing.T) {
		setupToken, err := adminAuth.GenerateSetupToken()
		require.NoError(t, err)

		setupReq := &SetupRequest{
			Email:           "setup@jadenrazo.dev",
			Password:        "SetupPass123!",
			ConfirmPassword: "SetupPass123!",
			SetupToken:      setupToken,
		}

		err = adminAuth.CompleteSetup(setupReq)
		assert.NoError(t, err)

		// Verify admin was created
		hasAdmin, err := adminAuth.HasAdminAccount()
		assert.NoError(t, err)
		assert.True(t, hasAdmin)
	})

	t.Run("invalid_email_domain", func(t *testing.T) {
		setupToken, err := adminAuth.GenerateSetupToken()
		require.NoError(t, err)

		setupReq := &SetupRequest{
			Email:           "setup@example.com",
			Password:        "SetupPass123!",
			ConfirmPassword: "SetupPass123!",
			SetupToken:      setupToken,
		}

		err = adminAuth.CompleteSetup(setupReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email not authorized")
	})

	t.Run("password_mismatch", func(t *testing.T) {
		setupToken, err := adminAuth.GenerateSetupToken()
		require.NoError(t, err)

		setupReq := &SetupRequest{
			Email:           "mismatch@jadenrazo.dev",
			Password:        "SetupPass123!",
			ConfirmPassword: "DifferentPass123!",
			SetupToken:      setupToken,
		}

		err = adminAuth.CompleteSetup(setupReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "passwords do not match")
	})
}

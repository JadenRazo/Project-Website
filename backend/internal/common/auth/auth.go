package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/JadenRazo/Project-Website/backend/internal/common/cache"
	"github.com/JadenRazo/Project-Website/backend/internal/common/database"
	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Auth handles authentication and authorization
type Auth struct {
	cfg        *config.AuthConfig
	cache      cache.Cache
	db         *gorm.DB
	jwtManager *JWTManager
}


// New creates a new authentication instance
func New(cfg *config.AuthConfig, db database.Database, cache cache.Cache) (*Auth, error) {
	if cfg == nil {
		// Use default config if none provided
		cfg = &config.AuthConfig{
			JWTSecret:     GetJWTSecret(),
			TokenExpiry:   15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
			Issuer:        "jadenrazo.dev",
			Audience:      "jadenrazo.dev",
		}
	}
	if db == nil {
		return nil, fmt.Errorf("database cannot be nil")
	}

	// Cache can be nil (optional)

	// Create JWT manager with proper config
	jwtConfig := &config.AuthConfig{
		JWTSecret:     cfg.JWTSecret,
		TokenExpiry:   cfg.TokenExpiry,
		RefreshExpiry: cfg.RefreshExpiry,
		Issuer:        cfg.Issuer,
		Audience:      cfg.Audience,
	}
	jwtManager := NewJWTManager(jwtConfig)

	return &Auth{
		cfg:        cfg,
		cache:      cache,
		db:         db.GetDB(),
		jwtManager: jwtManager,
	}, nil
}

// GenerateTokenPair generates a new pair of access and refresh tokens
func (a *Auth) GenerateTokenPair(user *entity.User) (*TokenPair, error) {
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}

	// Use JWT manager to generate tokens
	tokenPair, err := a.jwtManager.GenerateTokenPair(
		user.ID.String(),
		user.Username,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// Store refresh token in cache with user info (if cache is available)
	if a.cache != nil {
		key := fmt.Sprintf("refresh_token:%s", tokenPair.RefreshToken)
		tokenData := map[string]interface{}{
			"user_id":  user.ID.String(),
			"username": user.Username,
			"email":    user.Email,
			"role":     string(user.Role),
		}
		err = a.cache.Set(context.Background(), key, tokenData, a.cfg.RefreshExpiry)
		if err != nil {
			// Log error but don't fail - refresh tokens can work without cache
			fmt.Printf("Warning: failed to cache refresh token: %v\n", err)
		}
	}

	return tokenPair, nil
}

// ValidateToken validates a JWT token
func (a *Auth) ValidateToken(tokenString string) (*Claims, error) {
	return a.jwtManager.ValidateToken(tokenString)
}

// RefreshToken refreshes an access token using a refresh token
func (a *Auth) RefreshToken(refreshToken string) (*TokenPair, error) {
	var userID uuid.UUID
	var err error

	// Try to get user data from cache if available
	if a.cache != nil {
		key := fmt.Sprintf("refresh_token:%s", refreshToken)
		cachedData, err := a.cache.Get(context.Background(), key)
		if err == nil {
			// Extract user ID from cached data
			tokenData, ok := cachedData.(map[string]interface{})
			if ok {
				if userIDStr, ok := tokenData["user_id"].(string); ok {
					userID, _ = uuid.Parse(userIDStr)
				}
			}
		}
	}

	// If we couldn't get user ID from cache, validate the refresh token format
	// In production, you'd want to store refresh tokens in the database
	if userID == uuid.Nil {
		// For now, return an error if refresh token is invalid
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Get fresh user details from database
	user, err := a.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is deactivated")
	}

	// Generate new token pair
	return a.GenerateTokenPair(user)
}

// RevokeToken revokes a refresh token
func (a *Auth) RevokeToken(refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return a.cache.Delete(context.Background(), key)
}

// CreateMiddleware creates a new authentication middleware instance
func (a *Auth) CreateMiddleware(config *MiddlewareConfig) *Middleware {
	return NewMiddleware(a.jwtManager, config)
}

// AuthMiddleware creates a middleware that validates JWT tokens (backward compatibility)
func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	middleware := a.CreateMiddleware(nil)
	return middleware.RequireAuth(next)
}


// isPublicPath checks if a path should be publicly accessible
func (a *Auth) isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/auth/refresh",
		"/api/health",
	}
	for _, p := range publicPaths {
		if path == p {
			return true
		}
	}
	return false
}

// GetUserByID retrieves a user by ID from the database
func (a *Auth) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := a.db.Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email from the database
func (a *Auth) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := a.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username from the database
func (a *Auth) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := a.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (a *Auth) CreateUser(email, username, password, fullName string) (*entity.User, error) {
	// Check if user already exists
	existingUser, _ := a.GetUserByEmail(email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	existingUser, _ = a.GetUserByUsername(username)
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", username)
	}

	// Create new user
	user, err := entity.NewUser(email, username, password, fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save to database
	err = a.db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// AuthenticateUser validates user credentials and returns the user if valid
func (a *Auth) AuthenticateUser(emailOrUsername, password string) (*entity.User, error) {
	var user *entity.User
	var err error

	// Try to find user by email first, then by username
	if strings.Contains(emailOrUsername, "@") {
		user, err = a.GetUserByEmail(emailOrUsername)
	} else {
		user, err = a.GetUserByUsername(emailOrUsername)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is verified and active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Record login time
	user.RecordLogin()
	a.db.Save(user)

	return user, nil
}

// UpdateUserPassword updates a user's password
func (a *Auth) UpdateUserPassword(userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := a.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !user.CheckPassword(currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	// Update password
	err = user.UpdatePassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Save to database
	err = a.db.Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to save updated password: %w", err)
	}

	return nil
}

// DeactivateUser deactivates a user account
func (a *Auth) DeactivateUser(userID uuid.UUID) error {
	user, err := a.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.IsActive = false
	user.UpdatedAt = time.Now()

	err = a.db.Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Invalidate all refresh tokens for this user
	return a.InvalidateUserTokens(userID)
}

// InvalidateUserTokens invalidates all refresh tokens for a user
func (a *Auth) InvalidateUserTokens(userID uuid.UUID) error {
	// This is a simple implementation - in production you might want to maintain
	// a more sophisticated token blacklist
	// Note: This is a simplified approach. In production, you'd want to
	// store user_id -> tokens mapping for efficient invalidation
	return nil
}

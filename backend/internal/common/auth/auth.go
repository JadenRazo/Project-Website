package auth

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/JadenRazo/Project-Website/backend/internal/app/config"
    "github.com/JadenRazo/Project-Website/backend/internal/common/cache"
    "github.com/JadenRazo/Project-Website/backend/internal/common/database"
    "github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// Auth handles authentication and authorization
type Auth struct {
    cfg   *config.AuthConfig
    cache cache.Cache
    db    *gorm.DB
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

// Claims represents the JWT claims
type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// NewAuth creates a new authentication instance
func New(cfg *config.AuthConfig, db database.Database, cache cache.Cache) (*Auth, error) {
    if cfg == nil {
        return nil, fmt.Errorf("auth config cannot be nil")
    }
    if db == nil {
        return nil, fmt.Errorf("database cannot be nil")
    }
    if cache == nil {
        return nil, fmt.Errorf("cache cannot be nil")
    }

    return &Auth{
        cfg:   cfg,
        cache: cache,
        db:    db.GetDB(),
    }, nil
}

// GenerateTokenPair generates a new pair of access and refresh tokens
func (a *Auth) GenerateTokenPair(user *entity.User) (*TokenPair, error) {
    if user == nil {
        return nil, fmt.Errorf("user cannot be nil")
    }

    now := time.Now()
    
    // Generate access token
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
        UserID:   user.ID.String(),
        Username: user.Username,
        Email:    user.Email,
        Role:     string(user.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(a.cfg.TokenExpiry)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
            Issuer:    "jadenrazo-api",
            Subject:   user.ID.String(),
        },
    })

    accessTokenString, err := accessToken.SignedString([]byte(a.cfg.JWTSecret))
    if err != nil {
        return nil, fmt.Errorf("failed to sign access token: %v", err)
    }

    // Generate refresh token
    refreshToken, err := a.generateRefreshToken()
    if err != nil {
        return nil, err
    }

    // Store refresh token in cache with user info
    key := fmt.Sprintf("refresh_token:%s", refreshToken)
    tokenData := map[string]interface{}{
        "user_id": user.ID.String(),
        "username": user.Username,
        "email": user.Email,
        "role": string(user.Role),
    }
    err = a.cache.Set(context.Background(), key, tokenData, a.cfg.RefreshExpiry)
    if err != nil {
        return nil, fmt.Errorf("failed to store refresh token: %v", err)
    }

    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(a.cfg.TokenExpiry.Seconds()),
    }, nil
}

// ValidateToken validates a JWT token
func (a *Auth) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(a.cfg.JWTSecret), nil
    })

    if err != nil {
        return nil, fmt.Errorf("invalid token: %v", err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}

// RefreshToken refreshes an access token using a refresh token
func (a *Auth) RefreshToken(refreshToken string) (*TokenPair, error) {
    // Get user data from cache
    key := fmt.Sprintf("refresh_token:%s", refreshToken)
    cachedData, err := a.cache.Get(context.Background(), key)
    if err != nil {
        return nil, fmt.Errorf("invalid refresh token")
    }

    // Extract user ID from cached data
    tokenData, ok := cachedData.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid token data format")
    }

    userIDStr, ok := tokenData["user_id"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid user ID in token data")
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID format")
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

// AuthMiddleware creates a middleware that validates JWT tokens
func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip auth for certain paths
        if a.isPublicPath(r.URL.Path) {
            next.ServeHTTP(w, r)
            return
        }

        // Get token from header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        // Extract token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "invalid authorization header", http.StatusUnauthorized)
            return
        }

        // Validate token
        claims, err := a.ValidateToken(parts[1])
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        // Add claims to context
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// generateRefreshToken generates a secure random refresh token
func (a *Auth) generateRefreshToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
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
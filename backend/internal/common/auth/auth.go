package auth

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt"
    "github.com/JadenRazo/Project-Website/backend/internal/app/config"
    "github.com/JadenRazo/Project-Website/backend/internal/common/cache"
)

// Auth handles authentication and authorization
type Auth struct {
    cfg   *config.AuthConfig
    cache cache.Cache
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
    Role     string `json:"role"`
    jwt.StandardClaims
}

// NewAuth creates a new authentication instance
func NewAuth(cfg *config.AuthConfig, cache cache.Cache) *Auth {
    return &Auth{
        cfg:   cfg,
        cache: cache,
    }
}

// GenerateTokenPair generates a new pair of access and refresh tokens
func (a *Auth) GenerateTokenPair(userID, username, role string) (*TokenPair, error) {
    // Generate access token
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(a.cfg.TokenExpiry).Unix(),
            IssuedAt:  time.Now().Unix(),
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

    // Store refresh token in cache
    key := fmt.Sprintf("refresh_token:%s", refreshToken)
    err = a.cache.Set(context.Background(), key, userID, a.cfg.RefreshExpiry)
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
    // Get user ID from cache
    key := fmt.Sprintf("refresh_token:%s", refreshToken)
    userID, err := a.cache.Get(context.Background(), key)
    if err != nil {
        return nil, fmt.Errorf("invalid refresh token")
    }

    // Get user details from database (implement this based on your needs)
    user, err := a.getUserByID(userID.(string))
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }

    // Generate new token pair
    return a.GenerateTokenPair(user.ID, user.Username, user.Role)
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

// getUserByID retrieves a user by ID (implement this based on your needs)
type User struct {
    ID       string
    Username string
    Role     string
}

func (a *Auth) getUserByID(id string) (*User, error) {
    // Implement this based on your database
    return nil, fmt.Errorf("not implemented")
} 
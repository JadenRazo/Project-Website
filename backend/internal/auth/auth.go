package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Auth handles authentication and authorization
type Auth struct {
	jwtSecret     string
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	issuer        string
	audience      string
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// NewAuth creates a new authentication instance
func NewAuth(jwtSecret string, tokenExpiry, refreshExpiry time.Duration, issuer, audience string) *Auth {
	// Enforce minimum token expiry for security (15 minutes)
	if tokenExpiry < 15*time.Minute {
		tokenExpiry = 15 * time.Minute
	}

	// Enforce minimum refresh expiry (1 day)
	if refreshExpiry < 24*time.Hour {
		refreshExpiry = 24 * time.Hour
	}

	return &Auth{
		jwtSecret:     jwtSecret,
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
		audience:      audience,
	}
}

// GenerateTokenPair generates a new pair of access and refresh tokens
func (a *Auth) GenerateTokenPair(ctx context.Context, userID, username, role string) (*TokenPair, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Validate inputs
	if userID == "" || username == "" || role == "" {
		return nil, errors.New("userID, username, and role are required")
	}

	// Token ID for JTI claim
	tokenID := uuid.New().String()
	now := time.Now()

	// Generate access token with enhanced security
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(a.tokenExpiry).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    a.issuer,
			Audience:  a.audience,
			Id:        tokenID,
		},
	})

	accessTokenString, err := accessToken.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := a.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(a.tokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates a JWT token
func (a *Auth) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Sanitize input
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	// Parse with tight validation
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Validate issuer and audience
	if claims.Issuer != a.issuer {
		return nil, fmt.Errorf("%w: invalid issuer", ErrInvalidToken)
	}
	if claims.Audience != a.audience {
		return nil, fmt.Errorf("%w: invalid audience", ErrInvalidToken)
	}

	return claims, nil
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

		// Validate token with context
		claims, err := a.ValidateToken(r.Context(), parts[1])
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
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
		return "", fmt.Errorf("failed to generate secure random token: %w", err)
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
		"/swagger",
		"/metrics",
	}
	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

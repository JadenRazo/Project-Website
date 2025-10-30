package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidFormat      = errors.New("invalid authorization format")
	ErrMissingToken       = errors.New("authorization token required")
)

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	secret        []byte
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	issuer        string
	audience      string
}

// Claims represents the JWT claims with additional app-specific fields
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// NewJWTManager creates a new JWT manager with enhanced security defaults
func NewJWTManager(cfg *config.AuthConfig) *JWTManager {
	tokenExpiry := cfg.TokenExpiry
	refreshExpiry := cfg.RefreshExpiry
	
	// Enforce minimum token expiry for security (15 minutes)
	if tokenExpiry < 15*time.Minute {
		tokenExpiry = 15 * time.Minute
	}
	
	// Enforce minimum refresh expiry (1 day)
	if refreshExpiry < 24*time.Hour {
		refreshExpiry = 24 * time.Hour
	}
	
	return &JWTManager{
		secret:        []byte(cfg.JWTSecret),
		tokenExpiry:   tokenExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        cfg.Issuer,
		audience:      cfg.Audience,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (m *JWTManager) GenerateTokenPair(userID, username, email, role string) (*TokenPair, error) {
	// Validate inputs
	if userID == "" || username == "" || role == "" {
		return nil, errors.New("userID, username, and role are required")
	}

	now := time.Now()
	tokenID := uuid.New().String()

	// Create access token claims with enhanced security
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Audience:  jwt.ClaimStrings{m.audience},
			ID:        tokenID,
		},
	}

	// Create and sign access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := m.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.tokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// GenerateToken generates a single JWT token (for backward compatibility)
func (m *JWTManager) GenerateToken(userID, role string) (string, error) {
	return m.GenerateTokenWithDetails(userID, "", "", role)
}

// GenerateTokenWithDetails generates a JWT token with full user details
func (m *JWTManager) GenerateTokenWithDetails(userID, username, email, role string) (string, error) {
	if userID == "" || role == "" {
		return "", errors.New("userID and role are required")
	}

	now := time.Now()
	tokenID := uuid.New().String()

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Audience:  jwt.ClaimStrings{m.audience},
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// Sanitize input
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrMissingToken
	}

	// Parse token with comprehensive validation
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
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

	// Extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Validate issuer
	if claims.Issuer != m.issuer {
		return nil, fmt.Errorf("%w: invalid issuer", ErrInvalidToken)
	}

	// Validate audience - handle both string and []string audiences
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == m.audience {
			validAudience = true
			break
		}
	}
	if !validAudience {
		return nil, fmt.Errorf("%w: invalid audience", ErrInvalidToken)
	}

	return claims, nil
}

// GenerateRefreshToken generates a secure random refresh token
func (m *JWTManager) GenerateRefreshToken() (string, error) {
	// Create a secure random token
	b := make([]byte, 32) // 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ExtractTokenFromBearerString extracts a token from a Bearer string
func ExtractTokenFromBearerString(bearerString string) (string, error) {
	bearerString = strings.TrimSpace(bearerString)
	if !strings.HasPrefix(bearerString, "Bearer ") {
		return "", ErrInvalidFormat
	}
	
	token := strings.TrimSpace(bearerString[7:])
	if token == "" {
		return "", ErrMissingToken
	}
	
	return token, nil
}

// ExtractTokenFromHeader extracts and validates token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingToken
	}
	
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidFormat
	}
	
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrMissingToken
	}
	
	return token, nil
}

// generateTokenID generates a random token ID for the jti claim
func generateTokenID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return time.Now().Format(time.RFC3339Nano)
	}
	return base64.URLEncoding.EncodeToString(b)
}

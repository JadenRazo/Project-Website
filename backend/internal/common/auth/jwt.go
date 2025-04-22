package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
	"github.com/golang-jwt/jwt"
)

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	secret        []byte
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	issuer        string
	audience      string
}

// CustomClaims represents the JWT claims with additional app-specific fields
type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.AuthConfig) *JWTManager {
	return &JWTManager{
		secret:        []byte(cfg.JWTSecret),
		tokenExpiry:   cfg.TokenExpiry,
		refreshExpiry: cfg.RefreshExpiry,
		issuer:        cfg.Issuer,
		audience:      cfg.Audience,
	}
}

// GenerateToken generates a new JWT token for the given user
func (m *JWTManager) GenerateToken(userID, role string) (string, error) {
	now := time.Now()

	// Create token claims with strong security standards
	claims := &CustomClaims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(m.tokenExpiry).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    m.issuer,
			Audience:  m.audience,
			Id:        generateTokenID(),
		},
	}

	// Create token with claims using HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*CustomClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Validate issuer
	if claims.Issuer != m.issuer {
		return nil, errors.New("invalid token issuer")
	}

	// Validate audience
	if claims.Audience != m.audience {
		return nil, errors.New("invalid token audience")
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
	if len(bearerString) > 7 && bearerString[:7] == "Bearer " {
		return bearerString[7:], nil
	}
	return "", errors.New("invalid bearer token format")
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

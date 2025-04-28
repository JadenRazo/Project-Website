package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware handles JWT authentication for the messaging API
type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(cfg *config.Config) *auth.Middleware {
	return &auth.Middleware{
		RequireAuth: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get token from Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					http.Error(w, "Authorization header required", http.StatusUnauthorized)
					return
				}

				// Parse the Authorization header
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
					return
				}
				tokenString := parts[1]

				// Parse and validate the JWT token
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					// Validate signing method
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(cfg.JWT.Secret), nil
				})

				if err != nil {
					http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
					return
				}

				// Check if token is valid
				if !token.Valid {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}

				// Extract claims from token
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					http.Error(w, "Could not parse token claims", http.StatusUnauthorized)
					return
				}

				// Check token expiration
				if exp, ok := claims["exp"].(float64); ok {
					if time.Unix(int64(exp), 0).Before(time.Now()) {
						http.Error(w, "Token expired", http.StatusUnauthorized)
						return
					}
				}

				// Extract user ID from claims
				userID, ok := claims["sub"].(string)
				if !ok {
					http.Error(w, "Invalid token: missing user ID", http.StatusUnauthorized)
					return
				}

				// Add user info to context
				user := &User{
					ID:   userID,
					Name: claims["name"].(string),
				}

				// Add user to context
				ctx := context.WithValue(r.Context(), userKey, user)

				// Call the next handler with the updated context
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}
}

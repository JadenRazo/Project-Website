package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/auth"
	"github.com/JadenRazo/Project-Website/backend/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware handles JWT authentication for API requests
type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: cfg.JWT.Secret,
	}
}

// RequireAuth middleware ensures a valid JWT token is present
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
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
			return []byte(m.jwtSecret), nil
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
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid token: missing user ID", http.StatusUnauthorized)
			return
		}
		userID := uint(userIDFloat)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)

		// Add optional user role if present
		if role, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, auth.UserRoleKey, role)
		}

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware adds user info to context if token is present,
// but allows requests to proceed even without a valid token
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header, proceed without authentication
			next.ServeHTTP(w, r)
			return
		}

		// Parse the Authorization header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, proceed without authentication
			next.ServeHTTP(w, r)
			return
		}
		tokenString := parts[1]

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		// If token parsing fails, proceed without authentication
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// Check token expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Extract user ID from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		userID := uint(userIDFloat)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)

		// Add optional user role if present
		if role, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, auth.UserRoleKey, role)
		}

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware ensures the authenticated user has the required role
func (m *AuthMiddleware) RequireRole(roles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user role from context (set by RequireAuth middleware)
		role, ok := r.Context().Value(auth.UserRoleKey).(string)
		if !ok {
			http.Error(w, "Unauthorized: missing role", http.StatusUnauthorized)
			return
		}

		// Check if user's role is in the allowed roles
		allowed := false
		for _, allowedRole := range roles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware ensures the authenticated user is an admin
func (m *AuthMiddleware) AdminOnly(next http.Handler) http.Handler {
	return m.RequireRole([]string{"admin"}, next)
}

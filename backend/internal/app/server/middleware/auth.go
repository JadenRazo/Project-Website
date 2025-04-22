package middleware

import (
	"context"
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/common/auth"
)

// RequireAdmin middleware ensures the user has admin role
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user claims from context
		claims, ok := r.Context().Value("user").(*auth.Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if user has admin role
		if claims.Role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// Process the request
		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware ensures the user has a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user claims from context
			claims, ok := r.Context().Value("user").(*auth.Claims)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			if claims.Role != role {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			// Process the request
			next.ServeHTTP(w, r)
		})
	}
}

// WithUserID adds the user ID to the request context
func WithUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user claims from context
		claims, ok := r.Context().Value("user").(*auth.Claims)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)

		// Process the request with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

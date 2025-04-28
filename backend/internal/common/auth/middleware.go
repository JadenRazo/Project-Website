package auth

import (
	"net/http"
)

// Middleware represents the authentication middleware
type Middleware struct {
	RequireAuth func(next http.Handler) http.Handler
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(requireAuth func(next http.Handler) http.Handler) *Middleware {
	return &Middleware{
		RequireAuth: requireAuth,
	}
}

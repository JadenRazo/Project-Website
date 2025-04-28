package context

import (
	"context"
	"errors"

	"github.com/JadenRazo/Project-Website/backend/internal/domain/entity"
)

// Custom context keys
type contextKey string

const (
	// UserKey is the context key for the user
	UserKey contextKey = "user"

	// TraceIDKey is the context key for the trace ID
	TraceIDKey contextKey = "trace_id"

	// RequestIDKey is the context key for the request ID
	RequestIDKey contextKey = "request_id"
)

// ErrUserNotFound is returned when the user is not found in the context
var ErrUserNotFound = errors.New("user not found in context")

// ErrInvalidUserType is returned when the user in the context has an invalid type
var ErrInvalidUserType = errors.New("invalid user type in context")

// WithUser adds a user entity to the context
func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

// GetUser retrieves the user entity from the context
func GetUser(ctx context.Context) (*entity.User, error) {
	if ctx == nil {
		return nil, ErrUserNotFound
	}

	userVal := ctx.Value(UserKey)
	if userVal == nil {
		return nil, ErrUserNotFound
	}

	user, ok := userVal.(*entity.User)
	if !ok {
		return nil, ErrInvalidUserType
	}

	return user, nil
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID retrieves the trace ID from the context
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	traceID, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		return ""
	}

	return traceID
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

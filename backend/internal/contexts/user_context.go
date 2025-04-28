package contexts

import (
	"context"
)

// Key is a type for context keys
type Key string

const (
	// UserIDKey is the key for the user ID in the context
	UserIDKey Key = "user_id"
)

// GetUserID gets the user ID from the context
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// SetUserID sets the user ID in the context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

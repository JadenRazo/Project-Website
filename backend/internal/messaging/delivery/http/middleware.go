package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Middleware errors
var (
	ErrUnauthorized     = errors.New("unauthorized access")
	ErrChannelIDMissing = errors.New("channel ID is required")
	ErrInvalidChannelID = errors.New("invalid channel ID format")
	ErrMessageIDMissing = errors.New("message ID is required")
	ErrInvalidMessageID = errors.New("invalid message ID format")
	ErrAccessForbidden  = errors.New("insufficient permissions")
)

// Context keys to avoid collisions
type contextKey string

const (
	userKey          contextKey = "ctx_user"
	channelKey       contextKey = "ctx_channel"
	channelMemberKey contextKey = "ctx_channel_member"
	messageKey       contextKey = "ctx_message"
)

// User represents a user in the system
type User struct {
	ID   uuid.UUID
	Name string
}

// RequireChannelMembership middleware checks if the user is a member of the channel
func RequireChannelMembership(channelRepo domain.ChannelRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil {
				logger.Errorf("Failed to get user from context: %v", err)
				RespondWithError(w, http.StatusUnauthorized, ErrUnauthorized.Error())
				return
			}

			// Get channel ID from URL
			channelIDStr := chi.URLParam(r, "channelID")
			if channelIDStr == "" {
				RespondWithError(w, http.StatusBadRequest, ErrChannelIDMissing.Error())
				return
			}

			channelID, err := uuid.Parse(channelIDStr)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, ErrInvalidChannelID.Error())
				return
			}

			// Get channel
			channel, err := channelRepo.Get(r.Context(), channelID)
			if err != nil {
				if errors.Is(err, domain.ErrChannelNotFound) {
					RespondWithError(w, http.StatusNotFound, "Channel not found")
				} else {
					logger.Errorf("Failed to get channel: %v", err)
					RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve channel information")
				}
				return
			}

			// Public channels are accessible to all users
			if channel.IsPublic() {
				// Store channel in context for later use
				ctx := StoreChannelInContext(r.Context(), channel)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check if user is a member of the channel
			member, err := channelRepo.GetMember(r.Context(), channelID, user.ID)
			if err != nil {
				if errors.Is(err, domain.ErrNotChannelMember) {
					RespondWithError(w, http.StatusForbidden, "You are not a member of this channel")
				} else {
					logger.Errorf("Failed to check channel membership: %v", err)
					RespondWithError(w, http.StatusInternalServerError, "Failed to validate channel membership")
				}
				return
			}

			if member == nil {
				RespondWithError(w, http.StatusForbidden, "You are not a member of this channel")
				return
			}

			// Store channel and member in context for later use
			ctx := StoreChannelInContext(r.Context(), channel)
			ctx = StoreChannelMemberInContext(ctx, member)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireChannelAdmin middleware checks if the user is an admin of the channel
func RequireChannelAdmin(channelRepo domain.ChannelRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil {
				logger.Errorf("Failed to get user from context: %v", err)
				RespondWithError(w, http.StatusUnauthorized, ErrUnauthorized.Error())
				return
			}

			// Get channel ID from URL
			channelIDStr := chi.URLParam(r, "channelID")
			if channelIDStr == "" {
				RespondWithError(w, http.StatusBadRequest, ErrChannelIDMissing.Error())
				return
			}

			channelID, err := uuid.Parse(channelIDStr)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, ErrInvalidChannelID.Error())
				return
			}

			// Get channel
			channel, err := channelRepo.Get(r.Context(), channelID)
			if err != nil {
				if errors.Is(err, domain.ErrChannelNotFound) {
					RespondWithError(w, http.StatusNotFound, "Channel not found")
				} else {
					logger.Errorf("Failed to get channel: %v", err)
					RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve channel information")
				}
				return
			}

			// Allow channel creator to have admin access
			if channel.CreatedBy == user.ID {
				// Store channel in context for later use
				ctx := StoreChannelInContext(r.Context(), channel)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check if user is an admin of the channel
			member, err := channelRepo.GetMember(r.Context(), channelID, user.ID)
			if err != nil {
				if errors.Is(err, domain.ErrNotChannelMember) {
					RespondWithError(w, http.StatusForbidden, "You are not a member of this channel")
				} else {
					logger.Errorf("Failed to check channel membership: %v", err)
					RespondWithError(w, http.StatusInternalServerError, "Failed to validate channel membership")
				}
				return
			}

			if member == nil || !member.IsUserAdmin() {
				RespondWithError(w, http.StatusForbidden, "You don't have admin permissions for this channel")
				return
			}

			// Store channel and member in context for later use
			ctx := StoreChannelInContext(r.Context(), channel)
			ctx = StoreChannelMemberInContext(ctx, member)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateMessageOwnership middleware checks if the user is the owner of the message
func ValidateMessageOwnership(messageRepo domain.MessageRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, err := GetUserFromContext(r.Context())
			if err != nil {
				logger.Errorf("Failed to get user from context: %v", err)
				RespondWithError(w, http.StatusUnauthorized, ErrUnauthorized.Error())
				return
			}

			// Get message ID from URL
			messageIDStr := chi.URLParam(r, "messageID")
			if messageIDStr == "" {
				RespondWithError(w, http.StatusBadRequest, ErrMessageIDMissing.Error())
				return
			}

			messageID, err := uuid.Parse(messageIDStr)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, ErrInvalidMessageID.Error())
				return
			}

			// Get message
			message, err := messageRepo.Get(r.Context(), messageID)
			if err != nil {
				if errors.Is(err, domain.ErrMessageNotFound) {
					RespondWithError(w, http.StatusNotFound, "Message not found")
				} else {
					logger.Errorf("Failed to get message: %v", err)
					RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve message information")
				}
				return
			}

			// Check if message is deleted
			if message.IsDeleted {
				RespondWithError(w, http.StatusGone, "Message has been deleted")
				return
			}

			// Check if user owns the message
			if message.UserID != user.ID {
				// Allow channel admins to access
				channelMember, err := GetChannelMemberFromContext(r.Context())
				if err != nil || !channelMember.IsUserAdmin() {
					RespondWithError(w, http.StatusForbidden, "Only the message owner or channel admin can perform this action")
					return
				}
			}

			// Store message in context for later use
			ctx := StoreMessageInContext(r.Context(), message)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestLogger middleware logs incoming requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Recoverer middleware recovers from panics and logs them
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				RespondWithError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Context helper functions

// StoreChannelInContext stores a channel in the context
func StoreChannelInContext(ctx context.Context, channel *domain.Channel) context.Context {
	return context.WithValue(ctx, channelKey, channel)
}

// StoreChannelMemberInContext stores a channel member in the context
func StoreChannelMemberInContext(ctx context.Context, member *domain.ChannelMember) context.Context {
	return context.WithValue(ctx, channelMemberKey, member)
}

// StoreMessageInContext stores a message in the context
func StoreMessageInContext(ctx context.Context, message *domain.Message) context.Context {
	return context.WithValue(ctx, messageKey, message)
}

// GetUserFromContext gets a user from the context
func GetUserFromContext(ctx context.Context) (*User, error) {
	userValue := ctx.Value(userKey)
	if userValue == nil {
		return nil, ErrUnauthorized
	}

	user, ok := userValue.(*User)
	if !ok {
		return nil, ErrUnauthorized
	}

	return user, nil
}

// GetChannelFromContext gets a channel from the context
func GetChannelFromContext(ctx context.Context) (*domain.Channel, error) {
	channelValue := ctx.Value(channelKey)
	if channelValue == nil {
		return nil, errors.New("channel not found in context")
	}

	channel, ok := channelValue.(*domain.Channel)
	if !ok {
		return nil, errors.New("invalid channel in context")
	}

	return channel, nil
}

// GetChannelMemberFromContext gets a channel member from the context
func GetChannelMemberFromContext(ctx context.Context) (*domain.ChannelMember, error) {
	memberValue := ctx.Value(channelMemberKey)
	if memberValue == nil {
		return nil, errors.New("channel member not found in context")
	}

	member, ok := memberValue.(*domain.ChannelMember)
	if !ok {
		return nil, errors.New("invalid channel member in context")
	}

	return member, nil
}

// GetMessageFromContext gets a message from the context
func GetMessageFromContext(ctx context.Context) (*domain.Message, error) {
	messageValue := ctx.Value(messageKey)
	if messageValue == nil {
		return nil, errors.New("message not found in context")
	}

	message, ok := messageValue.(*domain.Message)
	if !ok {
		return nil, errors.New("invalid message in context")
	}

	return message, nil
}

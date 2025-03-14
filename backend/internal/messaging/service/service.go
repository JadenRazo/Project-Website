// internal/messaging/service/service.go
package service

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// MessagingService defines the interface for messaging operations
type MessagingService interface {
	// CreateChannel creates a new messaging channel
	CreateChannel(ctx context.Context, name, description string, isPrivate bool, creatorID uint) (*domain.Channel, error)

	// GetChannel retrieves a channel by ID
	GetChannel(ctx context.Context, channelID uint) (*domain.Channel, error)

	// ListUserChannels lists all channels a user is a member of
	ListUserChannels(ctx context.Context, userID uint) ([]domain.Channel, error)

	// SendMessage sends a new message to a channel
	SendMessage(ctx context.Context, content string, senderID, channelID uint) (*domain.Message, error)

	// GetChannelMessages retrieves messages from a channel with pagination
	GetChannelMessages(ctx context.Context, channelID uint, lastMessageID uint, limit int) ([]domain.Message, error)

	// AddUserToChannel adds a user to a channel
	AddUserToChannel(ctx context.Context, userID, channelID uint) error

	// RemoveUserFromChannel removes a user from a channel
	RemoveUserFromChannel(ctx context.Context, userID, channelID uint) error
}

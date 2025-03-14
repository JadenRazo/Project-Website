package repository

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// ChannelRepository defines the interface for channel operations
type ChannelRepository interface {
	// CreateChannel creates a new channel
	CreateChannel(ctx context.Context, channel *domain.Channel) error

	// GetChannel retrieves a channel by ID
	GetChannel(ctx context.Context, channelID uint) (*domain.Channel, error)

	// ListUserChannels lists all channels a user is a member of
	ListUserChannels(ctx context.Context, userID uint) ([]domain.Channel, error)

	// AddUserToChannel adds a user to a channel
	AddUserToChannel(ctx context.Context, channelID uint, userID uint) error

	// RemoveUserFromChannel removes a user from a channel
	RemoveUserFromChannel(ctx context.Context, channelID uint, userID uint) error

	// GetChannelMembers retrieves all members of a channel
	GetChannelMembers(ctx context.Context, channelID uint) ([]domain.User, error)

	// UpdateChannel updates channel information
	UpdateChannel(ctx context.Context, channel *domain.Channel) error

	// DeleteChannel deletes a channel
	DeleteChannel(ctx context.Context, channelID uint) error
}

// MessageRepository defines the interface for message operations
type MessageRepository interface {
	// CreateMessage creates a new message
	CreateMessage(ctx context.Context, message *domain.Message) error

	// GetMessage retrieves a message by ID
	GetMessage(ctx context.Context, messageID uint) (*domain.Message, error)

	// GetChannelMessages retrieves messages from a channel with pagination
	GetChannelMessages(ctx context.Context, channelID uint, lastMessageID uint, limit int) ([]domain.Message, error)

	// UpdateMessage updates a message
	UpdateMessage(ctx context.Context, message *domain.Message) error

	// DeleteMessage deletes a message
	DeleteMessage(ctx context.Context, messageID uint) error

	// AddReaction adds a reaction to a message
	AddReaction(ctx context.Context, reaction *domain.Reaction) error

	// RemoveReaction removes a reaction from a message
	RemoveReaction(ctx context.Context, reactionID uint) error

	// GetMessageReactions retrieves all reactions for a message
	GetMessageReactions(ctx context.Context, messageID uint) ([]domain.Reaction, error)

	// SearchMessages searches for messages in a channel by content
	SearchMessages(ctx context.Context, channelID uint, query string, limit int) ([]domain.Message, error)

	// GetUnreadMessages retrieves unread messages for a user in a channel before the specified time
	GetUnreadMessages(ctx context.Context, channelID uint, userID uint, upToTime time.Time) ([]domain.Message, error)
}

package repository

import (
	"context"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// MessagingReadReceiptRepository defines the repository for read receipts
type MessagingReadReceiptRepository interface {
	// AddReadReceipt adds a read receipt
	AddReadReceipt(ctx context.Context, receipt *domain.MessagingReadReceipt) error

	// GetMessageReadReceipts gets all read receipts for a message
	GetMessageReadReceipts(ctx context.Context, messageID uint) ([]domain.MessagingReadReceipt, error)

	// GetUserReadReceipts gets all read receipts for a user
	GetUserReadReceipts(ctx context.Context, userID uint) ([]domain.MessagingReadReceipt, error)

	// MarkChannelAsRead marks all messages in a channel as read for a user
	MarkChannelAsRead(ctx context.Context, channelID, userID uint) error
}

// MessagingMessageRepository defines the repository for messages
type MessagingMessageRepository interface {
	// GetMessage gets a message by ID
	GetMessage(ctx context.Context, messageID uint) (*domain.MessagingMessage, error)

	// GetChannelMessages gets all messages in a channel
	GetChannelMessages(ctx context.Context, channelID uint) ([]domain.MessagingMessage, error)
}

// MessagingChannelRepository defines the repository for channels
type MessagingChannelRepository interface {
	// GetChannel gets a channel by ID
	GetChannel(ctx context.Context, channelID uint) (*domain.MessagingChannel, error)

	// GetUserChannels gets all channels for a user
	GetUserChannels(ctx context.Context, userID uint) ([]domain.MessagingChannel, error)
}

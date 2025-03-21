package usecase

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// MarkAsReadUseCase handles message read status tracking
type MarkAsReadUseCase struct {
	messageRepo     repository.MessagingMessageRepository
	readReceiptRepo repository.MessagingReadReceiptRepository
	dispatcher      events.EventDispatcher
}

// NewMarkAsReadUseCase creates a new mark as read usecase
func NewMarkAsReadUseCase(
	messageRepo repository.MessagingMessageRepository,
	readReceiptRepo repository.MessagingReadReceiptRepository,
	dispatcher events.EventDispatcher,
) *MarkAsReadUseCase {
	return &MarkAsReadUseCase{
		messageRepo:     messageRepo,
		readReceiptRepo: readReceiptRepo,
		dispatcher:      dispatcher,
	}
}

// MarkAsReadInput represents the input for marking a message as read
type MarkAsReadInput struct {
	MessageID uint
	UserID    uint
}

// Execute marks a message as read by a user
func (uc *MarkAsReadUseCase) Execute(ctx context.Context, input MarkAsReadInput) error {
	// Get the message
	message, err := uc.messageRepo.GetMessage(ctx, input.MessageID)
	if err != nil {
		return err
	}

	// Don't create read receipts for your own messages
	if message.SenderID == input.UserID {
		return nil
	}

	// Create the read receipt
	receipt := &domain.MessagingReadReceipt{
		MessageID: input.MessageID,
		UserID:    input.UserID,
		ReadAt:    time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.readReceiptRepo.AddReadReceipt(ctx, receipt); err != nil {
		return err
	}

	// Dispatch read receipt event
	uc.dispatcher.Dispatch(ctx, events.NewReadReceiptEvent(receipt))

	return nil
}

// MarkChannelAsReadInput represents the input for marking all messages in a channel as read
type MarkChannelAsReadInput struct {
	ChannelID uint
	UserID    uint
}

// MarkChannelAsRead marks all messages in a channel as read
func (uc *MarkAsReadUseCase) MarkChannelAsRead(ctx context.Context, input MarkChannelAsReadInput) error {
	if err := uc.readReceiptRepo.MarkChannelAsRead(ctx, input.ChannelID, input.UserID); err != nil {
		return err
	}

	// Dispatch channel read event
	uc.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelRead, nil, input.ChannelID, input.UserID))

	return nil
}

// GetUnreadCountInput represents the input for getting unread message count
type GetUnreadCountInput struct {
	ChannelID uint
	UserID    uint
}

// GetUnreadCount gets the count of unread messages for a user in a channel
func (uc *MarkAsReadUseCase) GetUnreadCount(ctx context.Context, input GetUnreadCountInput) (int, error) {
	return uc.readReceiptRepo.GetUnreadCount(ctx, input.ChannelID, input.UserID)
}

// GetReadReceiptsInput represents the input for getting read receipts
type GetReadReceiptsInput struct {
	MessageID uint
}

// GetReadReceipts gets all read receipts for a message
func (uc *MarkAsReadUseCase) GetReadReceipts(ctx context.Context, input GetReadReceiptsInput) ([]domain.MessagingReadReceipt, error) {
	return uc.readReceiptRepo.GetMessageReadReceipts(ctx, input.MessageID)
}

// GetUserReadReceiptsInput represents the input for getting user's read receipts
type GetUserReadReceiptsInput struct {
	UserID uint
}

// GetUserReadReceipts gets all read receipts by a user
func (uc *MarkAsReadUseCase) GetUserReadReceipts(ctx context.Context, input GetUserReadReceiptsInput) ([]domain.MessagingReadReceipt, error) {
	return uc.readReceiptRepo.GetUserReadReceipts(ctx, input.UserID)
}

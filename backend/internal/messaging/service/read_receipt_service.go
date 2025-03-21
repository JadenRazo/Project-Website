package service

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// ReadReceiptService manages message read status tracking
type ReadReceiptService struct {
	readReceiptRepo repository.MessagingReadReceiptRepository
	messageRepo     repository.MessagingMessageRepository
	dispatcher      events.EventDispatcher
}

// NewReadReceiptService creates a new read receipt service
func NewReadReceiptService(
	readReceiptRepo repository.MessagingReadReceiptRepository,
	messageRepo repository.MessagingMessageRepository,
	dispatcher events.EventDispatcher,
) *ReadReceiptService {
	return &ReadReceiptService{
		readReceiptRepo: readReceiptRepo,
		messageRepo:     messageRepo,
		dispatcher:      dispatcher,
	}
}

// MarkAsRead marks a message as read by a user
func (s *ReadReceiptService) MarkAsRead(
	ctx context.Context,
	messageID uint,
	userID uint,
) error {
	// Get the message to verify it exists and get sender ID
	message, err := s.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}

	// Don't create read receipts for your own messages
	if message.SenderID == userID {
		return nil
	}

	// Create the read receipt
	receipt := &domain.MessagingReadReceipt{
		MessageID: messageID,
		UserID:    userID,
		ReadAt:    time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.readReceiptRepo.AddReadReceipt(ctx, receipt); err != nil {
		return err
	}

	// Dispatch read receipt event
	s.dispatcher.Dispatch(ctx, events.NewReadReceiptEvent(receipt))

	return nil
}

// MarkChannelAsRead marks all messages in a channel as read
func (s *ReadReceiptService) MarkChannelAsRead(
	ctx context.Context,
	channelID uint,
	userID uint,
) error {
	if err := s.readReceiptRepo.MarkChannelAsRead(ctx, channelID, userID); err != nil {
		return err
	}

	// Dispatch channel read event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelRead, nil, channelID, userID))

	return nil
}

// GetReadReceipts gets all read receipts for a message
func (s *ReadReceiptService) GetReadReceipts(
	ctx context.Context,
	messageID uint,
) ([]domain.MessagingReadReceipt, error) {
	return s.readReceiptRepo.GetMessageReadReceipts(ctx, messageID)
}

// GetUserReadReceipts gets all read receipts by a user
func (s *ReadReceiptService) GetUserReadReceipts(
	ctx context.Context,
	userID uint,
) ([]domain.MessagingReadReceipt, error) {
	return s.readReceiptRepo.GetUserReadReceipts(ctx, userID)
}

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

// MarkChannelAsRead marks all messages in a channel as read up to the given timestamp
func (s *ReadReceiptService) MarkChannelAsRead(
	ctx context.Context,
	channelID uint,
	userID uint,
	timestamp time.Time,
) error {
	// If timestamp is zero, use current time
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	// Mark channel as read up to the timestamp
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

// GetUnreadMessageCount gets the count of unread messages in a channel
func (s *ReadReceiptService) GetUnreadMessageCount(
	ctx context.Context,
	channelID uint,
	userID uint,
) (int, error) {
	// Get all messages in the channel
	messages, err := s.messageRepo.GetChannelMessages(ctx, channelID)
	if err != nil {
		return 0, err
	}

	// Get read receipts for the user
	receipts, err := s.readReceiptRepo.GetUserReadReceipts(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Create a map of read message IDs
	readMessageIDs := make(map[uint]bool)
	for _, receipt := range receipts {
		readMessageIDs[receipt.MessageID] = true
	}

	// Count unread messages
	unreadCount := 0
	for _, message := range messages {
		// Skip messages sent by the user
		if message.SenderID == userID {
			continue
		}

		// Check if the message has a read receipt
		if !readMessageIDs[message.ID] {
			unreadCount++
		}
	}

	return unreadCount, nil
}

// GetTotalUnreadMessageCount gets the total count of unread messages across all channels
func (s *ReadReceiptService) GetTotalUnreadMessageCount(
	ctx context.Context,
	userID uint,
) (int, error) {
	// Get all channels the user is a member of
	// This would require a channel repository, which we don't have in this function
	// For now, we'll return a placeholder implementation

	// Get read receipts for the user
	receipts, err := s.readReceiptRepo.GetUserReadReceipts(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Create a map of read message IDs
	readMessageIDs := make(map[uint]bool)
	for _, receipt := range receipts {
		readMessageIDs[receipt.MessageID] = true
	}

	// Get all messages across all channels
	// This is a very inefficient implementation and should be replaced with a proper query
	// For now, we'll return a placeholder implementation
	totalUnread := 0

	// TODO: Implement proper counting of unread messages across all channels

	return totalUnread, nil
}

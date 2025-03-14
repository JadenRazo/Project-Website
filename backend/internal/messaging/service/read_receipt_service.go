package service

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/websocket"
)

// ReadReceiptService manages message read status tracking
type ReadReceiptService struct {
	repository  repository.ReadReceiptRepository
	messageRepo repository.MessageRepository
	hub         *websocket.Hub
}

// NewReadReceiptService creates a new read receipt service
func NewReadReceiptService(
	repository repository.ReadReceiptRepository,
	messageRepo repository.MessageRepository,
	hub *websocket.Hub,
) *ReadReceiptService {
	return &ReadReceiptService{
		repository:  repository,
		messageRepo: messageRepo,
		hub:         hub,
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
		if err == repository.ErrMessageNotFound {
			return errors.ErrMessageNotFound
		}
		return errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to get message")
	}

	// Don't create read receipts for your own messages
	if message.SenderID == userID {
		return nil
	}

	// Check if receipt already exists to avoid duplicates
	exists, err := s.repository.ReadReceiptExists(ctx, messageID, userID)
	if err != nil {
		return errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to check read receipt")
	}

	if !exists {
		// Create the read receipt
		receipt := &repository.ReadReceipt{
			MessageID: messageID,
			UserID:    userID,
			ReadAt:    time.Now(),
		}

		if err := s.repository.CreateReadReceipt(ctx, receipt); err != nil {
			return errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to create read receipt")
		}

		// Notify through WebSocket
		s.hub.SendReadReceipt(messageID, message.ChannelID, userID, message.SenderID)
	}

	return nil
}

// MarkChannelAsRead marks all messages in a channel as read up to a certain time
func (s *ReadReceiptService) MarkChannelAsRead(
	ctx context.Context,
	channelID uint,
	userID uint,
	upToTime time.Time,
) error {
	// Get unread messages in channel before the specified time
	messages, err := s.messageRepo.GetUnreadMessages(ctx, channelID, userID, upToTime)
	if err != nil {
		return errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to get unread messages")
	}

	// Nothing to mark as read
	if len(messages) == 0 {
		return nil
	}

	// Create batch read receipts
	var receipts []repository.ReadReceipt
	for _, message := range messages {
		// Skip messages sent by the user
		if message.SenderID == userID {
			continue
		}

		receipts = append(receipts, repository.ReadReceipt{
			MessageID: message.ID,
			UserID:    userID,
			ReadAt:    time.Now(),
		})
	}

	// Save read receipts in batch
	if len(receipts) > 0 {
		if err := s.repository.CreateBulkReadReceipts(ctx, receipts); err != nil {
			return errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to create read receipts")
		}

		// Notify through WebSocket for the latest message only to reduce traffic
		if len(messages) > 0 {
			lastMessage := messages[len(messages)-1]
			s.hub.SendReadReceipt(lastMessage.ID, channelID, userID, lastMessage.SenderID)
		}
	}

	return nil
}

// GetReadReceipts gets all read receipts for a message
func (s *ReadReceiptService) GetReadReceipts(
	ctx context.Context,
	messageID uint,
) ([]repository.ReadReceipt, error) {
	receipts, err := s.repository.GetMessageReadReceipts(ctx, messageID)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to get read receipts")
	}

	return receipts, nil
}

// GetUnreadMessageCount gets the count of unread messages for a user in a channel
func (s *ReadReceiptService) GetUnreadMessageCount(
	ctx context.Context,
	channelID uint,
	userID uint,
) (int, error) {
	count, err := s.repository.GetUnreadCount(ctx, channelID, userID)
	if err != nil {
		return 0, errors.WrapError(err, errors.ErrorTypeDatabase, "db_error", "Failed to get unread count")
	}

	return count, nil
}

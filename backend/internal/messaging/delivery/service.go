package delivery

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"

	"github.com/google/uuid"
)

// DeliveryService handles message delivery status and tracking
type DeliveryService struct {
	deliveryRepo domain.MessageDeliveryRepository
	messageRepo  domain.MessageRepository
}

// NewDeliveryService creates a new DeliveryService
func NewDeliveryService(deliveryRepo domain.MessageDeliveryRepository, messageRepo domain.MessageRepository) *DeliveryService {
	return &DeliveryService{
		deliveryRepo: deliveryRepo,
		messageRepo:  messageRepo,
	}
}

// UpdateDeliveryStatus updates the delivery status of a message for a specific user
func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, messageID, userID uuid.UUID, status domain.MessageStatus) error {
	delivery := &domain.MessageDelivery{
		ID:        uuid.New(),
		MessageID: messageID,
		UserID:    userID,
		Status:    status,
		UpdatedAt: time.Now(),
	}

	return s.deliveryRepo.Update(ctx, delivery)
}

// GetMessageDeliveryStatus retrieves the delivery status for all recipients of a message
func (s *DeliveryService) GetMessageDeliveryStatus(ctx context.Context, messageID uuid.UUID) ([]*domain.MessageDelivery, error) {
	return s.deliveryRepo.GetMessageDeliveries(ctx, messageID)
}

// GetUserDeliveryStatus retrieves the delivery status of all messages for a specific user
func (s *DeliveryService) GetUserDeliveryStatus(ctx context.Context, userID uuid.UUID) ([]*domain.MessageDelivery, error) {
	return s.deliveryRepo.GetUserDeliveries(ctx, userID)
}

// MarkMessageAsDelivered marks a message as delivered for all recipients
func (s *DeliveryService) MarkMessageAsDelivered(ctx context.Context, messageID uuid.UUID) error {
	message, err := s.messageRepo.Get(ctx, messageID)
	if err != nil {
		return err
	}

	// Get all channel members
	members, err := s.messageRepo.GetChannelMembers(ctx, message.ChannelID)
	if err != nil {
		return err
	}

	// Update delivery status for each member
	for _, member := range members {
		delivery := &domain.MessageDelivery{
			ID:        uuid.New(),
			MessageID: messageID,
			UserID:    member.UserID,
			Status:    domain.MessageStatusDelivered,
			UpdatedAt: time.Now(),
		}

		if err := s.deliveryRepo.Update(ctx, delivery); err != nil {
			return err
		}
	}

	return nil
}

// MarkMessageAsRead marks a message as read for a specific user
func (s *DeliveryService) MarkMessageAsRead(ctx context.Context, messageID, userID uuid.UUID) error {
	delivery := &domain.MessageDelivery{
		ID:        uuid.New(),
		MessageID: messageID,
		UserID:    userID,
		Status:    domain.MessageStatusRead,
		UpdatedAt: time.Now(),
	}

	return s.deliveryRepo.Update(ctx, delivery)
}

// GetUnreadCount gets the count of unread messages for a user in a channel
func (s *DeliveryService) GetUnreadCount(ctx context.Context, channelID, userID uuid.UUID) (int, error) {
	messages, err := s.messageRepo.GetChannelMessages(ctx, channelID, domain.MessageFilter{}, domain.NewDefaultPagination())
	if err != nil {
		return 0, err
	}

	count := 0
	for _, msg := range messages {
		delivery, err := s.deliveryRepo.Get(ctx, msg.ID, userID)
		if err != nil {
			continue
		}
		if delivery.Status != domain.MessageStatusRead {
			count++
		}
	}

	return count, nil
}

// CleanupExpiredDeliveries removes delivery records for messages that have been deleted
func (s *DeliveryService) CleanupExpiredDeliveries(ctx context.Context) error {
	// Get all messages that have been deleted
	messages, err := s.messageRepo.List(ctx, domain.MessageFilter{IsDeleted: true}, domain.NewDefaultPagination())
	if err != nil {
		return err
	}

	// Remove delivery records for deleted messages
	for _, msg := range messages {
		if err := s.deliveryRepo.DeleteByMessageID(ctx, msg.ID); err != nil {
			return err
		}
	}

	return nil
}

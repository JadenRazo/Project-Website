package service

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	messagingRepo "github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// MessagingService implements service layer for messaging operations
type MessagingService struct {
	messageRepo messagingRepo.MessageRepository
	channelRepo messagingRepo.ChannelRepository
	dispatcher  *events.EventDispatcher
}

// NewMessagingService creates a new messaging service
func NewMessagingService(
	messageRepo messagingRepo.MessageRepository,
	channelRepo messagingRepo.ChannelRepository,
	dispatcher *events.EventDispatcher,
) *MessagingService {
	return &MessagingService{
		messageRepo: messageRepo,
		channelRepo: channelRepo,
		dispatcher:  dispatcher,
	}
}

// CreateMessage creates a new message in a channel
func (s *MessagingService) CreateMessage(ctx context.Context, content string, senderID, channelID uint) (*domain.Message, error) {
	// Validate parameters
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	// Verify channel exists and user has access
	channel, err := s.channelRepo.GetChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}

	// Check if user is a member of the channel
	isMember := false
	if channel.OwnerID == senderID {
		isMember = true
	} else {
		members, err := s.channelRepo.GetChannelMembers(ctx, channelID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			if member.UserID == senderID {
				isMember = true
				break
			}
		}
	}

	if !isMember {
		return nil, errors.New("user does not have access to the channel")
	}

	// Create the message
	message := &domain.Message{
		Content:   content,
		SenderID:  senderID,
		ChannelID: channelID,
		CreatedAt: time.Now(),
	}

	// Save the message
	err = s.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

package usecase

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// SendMessageUseCase handles message sending with rich features
type SendMessageUseCase struct {
	messageRepo    repository.MessagingMessageRepository
	channelRepo    repository.MessagingChannelRepository
	attachmentRepo repository.MessagingAttachmentRepository
	dispatcher     events.EventDispatcher
}

// NewSendMessageUseCase creates a new send message usecase
func NewSendMessageUseCase(
	messageRepo repository.MessagingMessageRepository,
	channelRepo repository.MessagingChannelRepository,
	attachmentRepo repository.MessagingAttachmentRepository,
	dispatcher events.EventDispatcher,
) *SendMessageUseCase {
	return &SendMessageUseCase{
		messageRepo:    messageRepo,
		channelRepo:    channelRepo,
		attachmentRepo: attachmentRepo,
		dispatcher:     dispatcher,
	}
}

// SendMessageInput represents the input for sending a message
type SendMessageInput struct {
	ChannelID    uint
	SenderID     uint
	Content      string
	Attachments  []*domain.MessagingAttachment
	Mentions     []uint
	ReplyToID    *uint
	IsPinned     bool
	IsNSFW       bool
	IsSpoiler    bool
	CustomEmoji  *string
	ThreadID     *uint
	MessageFlags uint
}

// Execute sends a new message
func (uc *SendMessageUseCase) Execute(ctx context.Context, input SendMessageInput) (*domain.MessagingMessage, error) {
	// Create message
	message := &domain.MessagingMessage{
		ChannelID:    input.ChannelID,
		SenderID:     input.SenderID,
		Content:      input.Content,
		Mentions:     input.Mentions,
		ReplyToID:    input.ReplyToID,
		IsNSFW:       input.IsNSFW,
		IsSpoiler:    input.IsSpoiler,
		CustomEmoji:  input.CustomEmoji,
		ThreadID:     input.ThreadID,
		MessageFlags: input.MessageFlags,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save message to database
	if err := uc.messageRepo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	// Handle attachments if any
	if len(input.Attachments) > 0 {
		for _, attachment := range input.Attachments {
			attachment.MessageID = message.ID
			if err := uc.attachmentRepo.AddAttachment(ctx, attachment); err != nil {
				return nil, err
			}
		}
	}

	// If message is pinned, update channel's pinned messages
	if input.IsPinned {
		if err := uc.channelRepo.AddPinnedMessage(ctx, message.ID); err != nil {
			return nil, err
		}
	}

	// Dispatch message created event
	uc.dispatcher.Dispatch(ctx, events.NewMessageEvent(events.MessagingEventMessageCreated, message))

	return message, nil
}

package usecase

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// PinMessageUseCase handles message pinning operations
type PinMessageUseCase struct {
	messageRepo repository.MessagingMessageRepository
	channelRepo repository.MessagingChannelRepository
	dispatcher  events.EventDispatcher
}

// NewPinMessageUseCase creates a new pin message usecase
func NewPinMessageUseCase(
	messageRepo repository.MessagingMessageRepository,
	channelRepo repository.MessagingChannelRepository,
	dispatcher events.EventDispatcher,
) *PinMessageUseCase {
	return &PinMessageUseCase{
		messageRepo: messageRepo,
		channelRepo: channelRepo,
		dispatcher:  dispatcher,
	}
}

// PinMessageInput represents the input for pinning a message
type PinMessageInput struct {
	MessageID uint
	UserID    uint
	ChannelID uint
}

// Execute pins a message
func (uc *PinMessageUseCase) Execute(ctx context.Context, input PinMessageInput) error {
	// Get the message
	message, err := uc.messageRepo.GetMessage(ctx, input.MessageID)
	if err != nil {
		return errors.ErrMessageNotFound
	}

	// Check if user has permission to pin messages
	if message.SenderID != input.UserID {
		// TODO: Check if user has admin/mod permissions
		return errors.ErrUnauthorized
	}

	// Check pin limit
	pinnedMessages, err := uc.channelRepo.GetPinnedMessages(ctx, message.ChannelID)
	if err != nil {
		return err
	}

	if len(pinnedMessages) >= 50 {
		return errors.ErrPinLimitExceeded
	}

	// Update message
	message.IsPinned = true
	message.PinnedAt = time.Now()
	message.PinnedBy = input.UserID

	if err := uc.messageRepo.UpdateMessage(ctx, message); err != nil {
		return err
	}

	// Add to channel's pinned messages
	if err := uc.channelRepo.AddPinnedMessage(ctx, message.ID); err != nil {
		return err
	}

	// Dispatch message pinned event
	event := events.NewMessageEvent(events.MessagingEventMessagePinned, message)
	uc.dispatcher.Dispatch(ctx, event)

	return nil
}

// UnpinMessageInput represents the input for unpinning a message
type UnpinMessageInput struct {
	MessageID uint
	UserID    uint
	ChannelID uint
}

// UnpinMessage unpins a message
func (uc *PinMessageUseCase) UnpinMessage(ctx context.Context, input UnpinMessageInput) error {
	// Get the message
	message, err := uc.messageRepo.GetMessage(ctx, input.MessageID)
	if err != nil {
		return errors.ErrMessageNotFound
	}

	// Verify message is in the specified channel
	if message.ChannelID != input.ChannelID {
		return errors.ErrMessageNotFound
	}

	// Update message pin status
	message.IsPinned = false
	message.PinnedAt = time.Time{}
	message.PinnedBy = 0
	message.UpdatedAt = time.Now()

	if err := uc.messageRepo.UpdateMessage(ctx, message); err != nil {
		return err
	}

	// Remove from channel's pinned messages
	if err := uc.channelRepo.RemovePinnedMessage(ctx, message.ID); err != nil {
		return err
	}

	// Dispatch message unpinned event
	uc.dispatcher.Dispatch(ctx, events.NewMessageEvent(events.MessagingEventMessageUnpinned, message))

	return nil
}

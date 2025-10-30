package usecase

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/messaging/errors"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// DeleteMessageUseCase handles message deletion with soft delete
type DeleteMessageUseCase struct {
	messageRepo    repository.MessagingMessageRepository
	channelRepo    repository.MessagingChannelRepository
	attachmentRepo repository.MessagingAttachmentRepository
	dispatcher     events.EventDispatcher
}

// NewDeleteMessageUseCase creates a new delete message usecase
func NewDeleteMessageUseCase(
	messageRepo repository.MessagingMessageRepository,
	channelRepo repository.MessagingChannelRepository,
	attachmentRepo repository.MessagingAttachmentRepository,
	dispatcher events.EventDispatcher,
) *DeleteMessageUseCase {
	return &DeleteMessageUseCase{
		messageRepo:    messageRepo,
		channelRepo:    channelRepo,
		attachmentRepo: attachmentRepo,
		dispatcher:     dispatcher,
	}
}

// DeleteMessageInput represents the input for deleting a message
type DeleteMessageInput struct {
	MessageID  uint
	UserID     uint
	HardDelete bool
}

// Execute deletes a message with soft delete support
func (uc *DeleteMessageUseCase) Execute(ctx context.Context, input DeleteMessageInput) error {
	// Get the message
	message, err := uc.messageRepo.GetMessage(ctx, input.MessageID)
	if err != nil {
		return err
	}

	// Check if user has permission to delete the message
	if message.SenderID != input.UserID {
		return errors.ErrUnauthorized
	}

	if input.HardDelete {
		// Delete attachments
		attachments, err := uc.attachmentRepo.GetMessageAttachments(ctx, input.MessageID)
		if err != nil {
			return err
		}

		for _, attachment := range attachments {
			if err := uc.attachmentRepo.DeleteAttachment(ctx, attachment.ID); err != nil {
				return err
			}
		}

		// Hard delete the message
		if err := uc.messageRepo.DeleteMessage(ctx, input.MessageID); err != nil {
			return err
		}

		// Dispatch message deleted event
		uc.dispatcher.Dispatch(ctx, events.NewMessageEvent(events.MessagingEventMessageDeleted, message))
	} else {
		// Soft delete the message
		message.DeletedAt = time.Now()
		message.UpdatedAt = time.Now()

		if err := uc.messageRepo.UpdateMessage(ctx, message); err != nil {
			return err
		}

		// Dispatch message updated event
		uc.dispatcher.Dispatch(ctx, events.NewMessageEvent(events.MessagingEventMessageUpdated, message))
	}

	return nil
}

package service

import (
	"context"
	"errors"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// ReactionService handles message reactions
type ReactionService struct {
	reactionRepo repository.MessagingReactionRepository
	messageRepo  repository.MessagingMessageRepository
	dispatcher   events.EventDispatcher
}

// NewReactionService creates a new reaction service
func NewReactionService(
	reactionRepo repository.MessagingReactionRepository,
	messageRepo repository.MessagingMessageRepository,
	dispatcher events.EventDispatcher,
) *ReactionService {
	return &ReactionService{
		reactionRepo: reactionRepo,
		messageRepo:  messageRepo,
		dispatcher:   dispatcher,
	}
}

// AddReaction adds a reaction to a message
func (s *ReactionService) AddReaction(ctx context.Context, messageID, userID uint, emojiCode, emojiName string) error {
	// Check if message exists
	_, err := s.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}

	// Create the reaction
	reaction := &domain.MessagingReaction{
		MessageID: messageID,
		UserID:    userID,
		EmojiCode: emojiCode,
		EmojiName: emojiName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.reactionRepo.AddReaction(ctx, reaction); err != nil {
		return err
	}

	// Dispatch reaction added event
	s.dispatcher.Dispatch(ctx, events.NewReactionEvent(events.MessagingEventReactionAdded, reaction))

	return nil
}

// RemoveReaction removes a reaction from a message
func (s *ReactionService) RemoveReaction(ctx context.Context, messageID, userID uint, emojiCode string) error {
	// Check if message exists
	_, err := s.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return err
	}

	// Get the reaction
	reactions, err := s.reactionRepo.GetMessageReactions(ctx, messageID)
	if err != nil {
		return err
	}

	var reaction *domain.MessagingReaction
	for _, r := range reactions {
		if r.UserID == userID && r.EmojiCode == emojiCode {
			reaction = &r
			break
		}
	}

	if reaction == nil {
		return errors.New("reaction not found")
	}

	if err := s.reactionRepo.RemoveReaction(ctx, reaction.ID); err != nil {
		return err
	}

	// Dispatch reaction removed event
	s.dispatcher.Dispatch(ctx, events.NewReactionEvent(events.MessagingEventReactionRemoved, reaction))

	return nil
}

// GetMessageReactions retrieves all reactions on a message
func (s *ReactionService) GetMessageReactions(ctx context.Context, messageID uint) ([]domain.MessagingReaction, error) {
	return s.reactionRepo.GetMessageReactions(ctx, messageID)
}

// GetUserReactions retrieves all reactions by a user
func (s *ReactionService) GetUserReactions(ctx context.Context, userID uint) ([]domain.MessagingReaction, error) {
	return s.reactionRepo.GetUserReactions(ctx, userID)
}

// GetReaction retrieves a specific reaction by ID
func (s *ReactionService) GetReaction(ctx context.Context, reactionID uint) (*domain.MessagingReaction, error) {
	return s.reactionRepo.GetReaction(ctx, reactionID)
}

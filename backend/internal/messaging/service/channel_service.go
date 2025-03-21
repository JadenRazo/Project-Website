package service

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/events"
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/repository"
)

// ChannelService handles channel operations
type ChannelService struct {
	channelRepo repository.MessagingChannelRepository
	dispatcher  events.EventDispatcher
}

// NewChannelService creates a new channel service
func NewChannelService(channelRepo repository.MessagingChannelRepository, dispatcher events.EventDispatcher) *ChannelService {
	return &ChannelService{
		channelRepo: channelRepo,
		dispatcher:  dispatcher,
	}
}

// CreateChannel creates a new channel
func (s *ChannelService) CreateChannel(ctx context.Context, name, description, channelType string, ownerID uint, isPrivate bool, categoryID *uint) (*domain.MessagingChannel, error) {
	// Create the channel
	channel := &domain.MessagingChannel{
		Name:        name,
		Description: description,
		Type:        channelType,
		OwnerID:     ownerID,
		IsPrivate:   isPrivate,
		CategoryID:  categoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.channelRepo.CreateChannel(ctx, channel); err != nil {
		return nil, err
	}

	// Add owner as first member
	if err := s.channelRepo.AddChannelMember(ctx, channel.ID, ownerID); err != nil {
		return nil, err
	}

	// Dispatch channel created event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelCreated, channel, channel.ID, ownerID))

	return channel, nil
}

// GetChannel retrieves a channel by ID
func (s *ChannelService) GetChannel(ctx context.Context, channelID uint) (*domain.MessagingChannel, error) {
	return s.channelRepo.GetChannel(ctx, channelID)
}

// UpdateChannel updates an existing channel
func (s *ChannelService) UpdateChannel(ctx context.Context, channelID uint, name, description string, userID uint) (*domain.MessagingChannel, error) {
	// Get the channel
	channel, err := s.channelRepo.GetChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}

	// Update the channel
	channel.Name = name
	channel.Description = description
	channel.UpdatedAt = time.Now()

	if err := s.channelRepo.UpdateChannel(ctx, channel); err != nil {
		return nil, err
	}

	// Dispatch channel updated event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelUpdated, channel, channelID, userID))

	return channel, nil
}

// DeleteChannel deletes a channel
func (s *ChannelService) DeleteChannel(ctx context.Context, channelID uint, userID uint) error {
	// Get the channel
	channel, err := s.channelRepo.GetChannel(ctx, channelID)
	if err != nil {
		return err
	}

	if err := s.channelRepo.DeleteChannel(ctx, channelID); err != nil {
		return err
	}

	// Dispatch channel deleted event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelDeleted, channel, channelID, userID))

	return nil
}

// ListUserChannels lists all channels a user is a member of
func (s *ChannelService) ListUserChannels(ctx context.Context, userID uint) ([]domain.MessagingChannel, error) {
	return s.channelRepo.ListUserChannels(ctx, userID)
}

// AddChannelMember adds a user to a channel
func (s *ChannelService) AddChannelMember(ctx context.Context, channelID, userID, addedByID uint) error {
	if err := s.channelRepo.AddChannelMember(ctx, channelID, userID); err != nil {
		return err
	}

	// Dispatch member added event
	s.dispatcher.Dispatch(ctx, events.NewMemberEvent(events.MessagingEventMemberAdded, channelID, userID, addedByID))

	return nil
}

// RemoveChannelMember removes a user from a channel
func (s *ChannelService) RemoveChannelMember(ctx context.Context, channelID, userID, removedByID uint) error {
	if err := s.channelRepo.RemoveChannelMember(ctx, channelID, userID); err != nil {
		return err
	}

	// Dispatch member removed event
	s.dispatcher.Dispatch(ctx, events.NewMemberEvent(events.MessagingEventMemberRemoved, channelID, userID, removedByID))

	return nil
}

// GetChannelMembers retrieves all members of a channel
func (s *ChannelService) GetChannelMembers(ctx context.Context, channelID uint) ([]domain.MessagingUser, error) {
	return s.channelRepo.GetChannelMembers(ctx, channelID)
}

// SetChannelSlowMode sets the slow mode for a channel
func (s *ChannelService) SetChannelSlowMode(ctx context.Context, channelID uint, seconds int, userID uint) error {
	// Get the channel
	channel, err := s.channelRepo.GetChannel(ctx, channelID)
	if err != nil {
		return err
	}

	// Update slow mode
	channel.SlowMode = seconds
	channel.UpdatedAt = time.Now()

	if err := s.channelRepo.UpdateChannel(ctx, channel); err != nil {
		return err
	}

	// Dispatch channel updated event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelUpdated, channel, channelID, userID))

	return nil
}

// SetChannelNSFW sets the NSFW flag for a channel
func (s *ChannelService) SetChannelNSFW(ctx context.Context, channelID uint, isNSFW bool, userID uint) error {
	// Get the channel
	channel, err := s.channelRepo.GetChannel(ctx, channelID)
	if err != nil {
		return err
	}

	// Update NSFW flag
	channel.IsNSFW = isNSFW
	channel.UpdatedAt = time.Now()

	if err := s.channelRepo.UpdateChannel(ctx, channel); err != nil {
		return err
	}

	// Dispatch channel updated event
	s.dispatcher.Dispatch(ctx, events.NewChannelEvent(events.MessagingEventChannelUpdated, channel, channelID, userID))

	return nil
}

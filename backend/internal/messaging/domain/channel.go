package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrChannelNotFound is returned when a channel cannot be found
var ErrChannelNotFound = errors.New("channel not found")

// ErrInvalidChannel is returned when channel data is invalid
var ErrInvalidChannel = errors.New("invalid channel data")

// ErrNotChannelMember is returned when the user is not a member of the channel
var ErrNotChannelMember = errors.New("user is not a member of the channel")

// ErrInsufficientPermissions is returned when the user doesn't have sufficient permissions
var ErrInsufficientPermissions = errors.New("insufficient permissions")

// NewChannel creates a new channel
func NewChannel(name, description string, channelType ChannelType, createdBy uuid.UUID) (*Channel, error) {
	if name == "" {
		return nil, errors.New("channel name cannot be empty")
	}

	if createdBy == uuid.Nil {
		return nil, errors.New("creator ID cannot be empty")
	}

	if channelType == "" {
		channelType = ChannelTypePublic
	}

	now := time.Now()

	return &Channel{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Type:        channelType,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsArchived:  false,
	}, nil
}

// NewDirectChannel creates a new direct message channel between two users
func NewDirectChannel(user1ID, user2ID uuid.UUID) (*Channel, error) {
	if user1ID == uuid.Nil || user2ID == uuid.Nil {
		return nil, errors.New("user IDs cannot be empty")
	}

	if user1ID == user2ID {
		return nil, errors.New("cannot create direct channel with the same user")
	}

	now := time.Now()

	return &Channel{
		ID:          uuid.New(),
		Name:        "Direct Channel",
		Description: "Direct messages",
		Type:        ChannelTypeDirect,
		CreatedBy:   user1ID,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsArchived:  false,
	}, nil
}

// NewChannelMember creates a new channel member
func NewChannelMember(channelID, userID uuid.UUID, role ChannelMemberRole) *ChannelMember {
	now := time.Now()

	return &ChannelMember{
		ChannelID:  channelID,
		UserID:     userID,
		JoinedAt:   now,
		Role:       string(role),
		IsAdmin:    role == ChannelMemberRoleOwner || role == ChannelMemberRoleAdmin,
		LastReadAt: now,
		IsMuted:    false,
	}
}

// IsPublic checks if the channel is public
func (c *Channel) IsPublic() bool {
	return c.Type == ChannelTypePublic
}

// IsPrivate checks if the channel is private
func (c *Channel) IsPrivate() bool {
	return c.Type == ChannelTypePrivate
}

// IsDirect checks if the channel is a direct message channel
func (c *Channel) IsDirect() bool {
	return c.Type == ChannelTypeDirect
}

// IsOwner checks if the user is the channel owner
func (cm *ChannelMember) IsOwner() bool {
	return cm.Role == string(ChannelMemberRoleOwner)
}

// IsUserAdmin checks if the user is a channel admin
func (cm *ChannelMember) IsUserAdmin() bool {
	return cm.Role == string(ChannelMemberRoleOwner) || cm.Role == string(ChannelMemberRoleAdmin)
}

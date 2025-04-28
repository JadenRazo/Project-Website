package domain

import (
	"context"
)

// ChannelUseCase defines the use case interface for channel operations
type ChannelUseCase interface {
	// CreateChannel creates a new channel
	CreateChannel(ctx context.Context, channel *MessagingChannel, members []*MessagingMember) (*MessagingChannel, error)

	// GetChannel retrieves a channel by ID
	GetChannel(ctx context.Context, channelID string) (*MessagingChannel, error)

	// UpdateChannel updates an existing channel
	UpdateChannel(ctx context.Context, channel *MessagingChannel) (*MessagingChannel, error)

	// DeleteChannel deletes a channel
	DeleteChannel(ctx context.Context, channelID string) error

	// ListChannels retrieves a list of channels with pagination
	ListChannels(ctx context.Context, userID string, pagination *Pagination) ([]MessagingChannel, int, error)

	// IsChannelAdmin checks if a user is an admin of a channel
	IsChannelAdmin(ctx context.Context, channelID string, userID string) (bool, error)

	// AddChannelMember adds a member to a channel
	AddChannelMember(ctx context.Context, channelID string, userID string, role string) (*MessagingMember, error)

	// ListChannelMembers retrieves members of a channel with pagination
	ListChannelMembers(ctx context.Context, channelID string, pagination *Pagination) ([]MessagingMember, int, error)

	// RemoveChannelMember removes a member from a channel
	RemoveChannelMember(ctx context.Context, channelID string, userID string) error

	// UpdateChannelMember updates a channel member's role
	UpdateChannelMember(ctx context.Context, channelID string, userID string, role string) (*MessagingMember, error)
}

// Pagination defines pagination parameters
type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty"`
	PerPage  int    `json:"-"` // Alias for PageSize for backward compatibility
	Total    int    `json:"total,omitempty"`
}

// MessagingMember represents a channel member
type MessagingMember struct {
	ID        uint   `json:"id"`
	ChannelID uint   `json:"channel_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	JoinedAt  string `json:"joined_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

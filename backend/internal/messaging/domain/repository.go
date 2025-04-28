package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ChannelRepository defines the interface for channel data access
type ChannelRepository interface {
	// Create creates a new channel
	Create(ctx context.Context, channel *Channel) error

	// Get retrieves a channel by ID
	Get(ctx context.Context, id uuid.UUID) (*Channel, error)

	// Update updates an existing channel
	Update(ctx context.Context, channel *Channel) error

	// Delete deletes a channel
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a list of channels with pagination
	List(ctx context.Context, filter ChannelFilter, pagination Pagination) ([]*Channel, int, error)

	// GetUserChannels retrieves channels for a specific user
	GetUserChannels(ctx context.Context, userID uuid.UUID) ([]*Channel, error)

	// GetDirectChannel retrieves a direct channel between two users
	GetDirectChannel(ctx context.Context, user1ID, user2ID uuid.UUID) (*Channel, error)

	// AddMember adds a user to a channel
	AddMember(ctx context.Context, channelID, userID uuid.UUID, role ChannelMemberRole) error

	// RemoveMember removes a user from a channel
	RemoveMember(ctx context.Context, channelID, userID uuid.UUID) error

	// GetMembers retrieves members of a channel
	GetMembers(ctx context.Context, channelID uuid.UUID) ([]*ChannelMember, error)

	// GetMember retrieves a specific member of a channel
	GetMember(ctx context.Context, channelID, userID uuid.UUID) (*ChannelMember, error)

	// UpdateMember updates a channel member
	UpdateMember(ctx context.Context, member *ChannelMember) error
}

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	// Create creates a new message
	Create(ctx context.Context, message *Message) error

	// Get retrieves a message by ID
	Get(ctx context.Context, id uuid.UUID) (*Message, error)

	// Update updates an existing message
	Update(ctx context.Context, message *Message) error

	// Delete marks a message as deleted
	Delete(ctx context.Context, id uuid.UUID) error

	// GetChannelMessages retrieves messages for a channel with pagination
	GetChannelMessages(ctx context.Context, channelID uuid.UUID, filter MessageFilter, pagination Pagination) ([]*Message, int, error)

	// GetThreadMessages retrieves replies to a message
	GetThreadMessages(ctx context.Context, messageID uuid.UUID, pagination Pagination) ([]*Message, int, error)

	// GetMessageAttachments retrieves attachments for a message
	GetMessageAttachments(ctx context.Context, messageID uuid.UUID) ([]*Attachment, error)

	// CreateAttachment creates a new attachment
	CreateAttachment(ctx context.Context, attachment *Attachment) error

	// DeleteAttachment deletes an attachment
	DeleteAttachment(ctx context.Context, id uuid.UUID) error

	// GetUserMentions retrieves mentions for a user
	GetUserMentions(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]*Message, int, error)
}

// ReactionRepository defines the interface for reaction data access
type ReactionRepository interface {
	// Create creates a new reaction
	Create(ctx context.Context, reaction *Reaction) error

	// Delete deletes a reaction
	Delete(ctx context.Context, messageID, userID uuid.UUID, emoji string) error

	// GetMessageReactions retrieves reactions for a message
	GetMessageReactions(ctx context.Context, messageID uuid.UUID) ([]*Reaction, error)

	// GetUserReactions retrieves reactions by a user
	GetUserReactions(ctx context.Context, userID uuid.UUID) ([]*Reaction, error)

	// GetReactionCounts gets reaction counts for a message
	GetReactionCounts(ctx context.Context, messageID uuid.UUID) ([]*ReactionCount, error)
}

// ReadReceiptRepository defines the interface for read receipt data access
type ReadReceiptRepository interface {
	// Create creates a new read receipt
	Create(ctx context.Context, receipt *ReadReceipt) error

	// GetMessageReceipts retrieves read receipts for a message
	GetMessageReceipts(ctx context.Context, messageID uuid.UUID) ([]*ReadReceipt, error)

	// GetChannelReceipts retrieves read receipts for a channel
	GetChannelReceipts(ctx context.Context, channelID uuid.UUID) ([]*ReadReceipt, error)

	// GetUserReceipts retrieves read receipts for a user
	GetUserReceipts(ctx context.Context, userID uuid.UUID) ([]*ReadReceipt, error)

	// GetUserChannelReceipt retrieves the latest read receipt for a user in a channel
	GetUserChannelReceipt(ctx context.Context, channelID, userID uuid.UUID) (*ReadReceipt, error)

	// GetChannelReadStatus gets the read status of a channel for a user
	GetChannelReadStatus(ctx context.Context, channelID, userID uuid.UUID) (*ChannelReadStatus, error)

	// MarkChannelAsRead marks all messages in a channel as read for a user
	MarkChannelAsRead(ctx context.Context, channelID, userID uuid.UUID) error
}

// MessageDeliveryRepository defines the interface for message delivery status data access
type MessageDeliveryRepository interface {
	// Create creates a new message delivery status
	Create(ctx context.Context, delivery *MessageDelivery) error

	// Update updates a message delivery status
	Update(ctx context.Context, delivery *MessageDelivery) error

	// GetMessageDeliveries retrieves delivery statuses for a message
	GetMessageDeliveries(ctx context.Context, messageID uuid.UUID) ([]*MessageDelivery, error)

	// GetUserDeliveries retrieves delivery statuses for a user
	GetUserDeliveries(ctx context.Context, userID uuid.UUID) ([]*MessageDelivery, error)
}

// MessageDraftRepository defines the interface for message draft data access
type MessageDraftRepository interface {
	// Create creates a new message draft
	Create(ctx context.Context, draft *MessageDraft) error

	// Update updates an existing message draft
	Update(ctx context.Context, draft *MessageDraft) error

	// Delete deletes a message draft
	Delete(ctx context.Context, id uuid.UUID) error

	// Get retrieves a message draft by ID
	Get(ctx context.Context, id uuid.UUID) (*MessageDraft, error)

	// GetUserDrafts retrieves drafts for a user
	GetUserDrafts(ctx context.Context, userID uuid.UUID) ([]*MessageDraft, error)

	// GetChannelDrafts retrieves drafts for a channel
	GetChannelDrafts(ctx context.Context, channelID uuid.UUID) ([]*MessageDraft, error)
}

// TypingIndicatorRepository defines the interface for typing indicator data access
type TypingIndicatorRepository interface {
	// Create creates a new typing indicator
	Create(ctx context.Context, indicator *TypingIndicator) error

	// Delete deletes a typing indicator
	Delete(ctx context.Context, id uuid.UUID) error

	// GetChannelTyping retrieves typing indicators for a channel
	GetChannelTyping(ctx context.Context, channelID uuid.UUID) ([]*TypingIndicator, error)

	// CleanupExpired removes expired typing indicators
	CleanupExpired(ctx context.Context, before time.Time) error
}

// ChannelCategoryRepository defines the interface for channel category data access
type ChannelCategoryRepository interface {
	// Create creates a new channel category
	Create(ctx context.Context, category *ChannelCategory) error

	// Update updates an existing channel category
	Update(ctx context.Context, category *ChannelCategory) error

	// Delete deletes a channel category
	Delete(ctx context.Context, id uuid.UUID) error

	// Get retrieves a channel category by ID
	Get(ctx context.Context, id uuid.UUID) (*ChannelCategory, error)

	// List retrieves all channel categories
	List(ctx context.Context) ([]*ChannelCategory, error)
}

// ChannelInviteRepository defines the interface for channel invite data access
type ChannelInviteRepository interface {
	// Create creates a new channel invite
	Create(ctx context.Context, invite *ChannelInvite) error

	// Update updates an existing channel invite
	Update(ctx context.Context, invite *ChannelInvite) error

	// Get retrieves a channel invite by ID
	Get(ctx context.Context, id uuid.UUID) (*ChannelInvite, error)

	// GetUserInvites retrieves invites for a user
	GetUserInvites(ctx context.Context, userID uuid.UUID) ([]*ChannelInvite, error)

	// GetChannelInvites retrieves invites for a channel
	GetChannelInvites(ctx context.Context, channelID uuid.UUID) ([]*ChannelInvite, error)

	// CleanupExpired removes expired invites
	CleanupExpired(ctx context.Context) error
}

// MessageReportRepository defines the interface for message report data access
type MessageReportRepository interface {
	// Create creates a new message report
	Create(ctx context.Context, report *MessageReport) error

	// Update updates an existing message report
	Update(ctx context.Context, report *MessageReport) error

	// Get retrieves a message report by ID
	Get(ctx context.Context, id uuid.UUID) (*MessageReport, error)

	// GetMessageReports retrieves reports for a message
	GetMessageReports(ctx context.Context, messageID uuid.UUID) ([]*MessageReport, error)

	// GetUserReports retrieves reports by a user
	GetUserReports(ctx context.Context, userID uuid.UUID) ([]*MessageReport, error)

	// GetPendingReports retrieves pending reports
	GetPendingReports(ctx context.Context) ([]*MessageReport, error)
}

// PinnedMessageRepository defines the interface for pinned message data access
type PinnedMessageRepository interface {
	// Create creates a new pinned message
	Create(ctx context.Context, pinned *PinnedMessage) error

	// Delete deletes a pinned message
	Delete(ctx context.Context, id uuid.UUID) error

	// GetChannelPinnedMessages retrieves pinned messages for a channel
	GetChannelPinnedMessages(ctx context.Context, channelID uuid.UUID) ([]*PinnedMessage, error)

	// GetMessagePin retrieves a pin for a specific message
	GetMessagePin(ctx context.Context, messageID uuid.UUID) (*PinnedMessage, error)
}

// MessageRetentionPolicyRepository defines the interface for message retention policy data access
type MessageRetentionPolicyRepository interface {
	// Create creates a new message retention policy
	Create(ctx context.Context, policy *MessageRetentionPolicy) error

	// Update updates an existing message retention policy
	Update(ctx context.Context, policy *MessageRetentionPolicy) error

	// Get retrieves a message retention policy by ID
	Get(ctx context.Context, id uuid.UUID) (*MessageRetentionPolicy, error)

	// GetChannelPolicy retrieves the retention policy for a channel
	GetChannelPolicy(ctx context.Context, channelID uuid.UUID) (*MessageRetentionPolicy, error)

	// List retrieves all message retention policies
	List(ctx context.Context) ([]*MessageRetentionPolicy, error)
}

// ChannelFilter defines filtering criteria for channels
type ChannelFilter struct {
	UserID     *uuid.UUID    `json:"user_id,omitempty"`
	Types      []ChannelType `json:"types,omitempty"`
	IsArchived *bool         `json:"is_archived,omitempty"`
	Search     string        `json:"search,omitempty"`
}

// MessageFilter defines filtering criteria for messages
type MessageFilter struct {
	UserID        *uuid.UUID    `json:"user_id,omitempty"`
	Types         []MessageType `json:"types,omitempty"`
	HasAttachment *bool         `json:"has_attachment,omitempty"`
	FromDate      *time.Time    `json:"from_date,omitempty"`
	ToDate        *time.Time    `json:"to_date,omitempty"`
	Search        string        `json:"search,omitempty"`
}

// Pagination defines pagination parameters
type Pagination struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	SortBy    string `json:"sort_by,omitempty"`
	SortOrder string `json:"sort_order,omitempty"`
}

// NewDefaultPagination creates a default pagination
func NewDefaultPagination() Pagination {
	return Pagination{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Adjust ensures pagination values are within acceptable ranges
func (p *Pagination) Adjust() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize < 1 {
		p.PageSize = 20
	}

	if p.PageSize > 100 {
		p.PageSize = 100
	}

	if p.SortBy == "" {
		p.SortBy = "created_at"
	}

	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}
}

// UserSettingsRepository defines the interface for user settings data access
type UserSettingsRepository interface {
	// Get retrieves user settings
	Get(ctx context.Context, userID uuid.UUID) (*UserSettings, error)

	// Update updates user settings
	Update(ctx context.Context, settings *UserSettings) error

	// Create creates new user settings
	Create(ctx context.Context, settings *UserSettings) error

	// GetNotificationSettings gets notification settings
	GetNotificationSettings(ctx context.Context, userID uuid.UUID) (*NotificationSettings, error)

	// UpdateNotificationSettings updates notification settings
	UpdateNotificationSettings(ctx context.Context, userID uuid.UUID, settings *NotificationSettings) error

	// GetPrivacySettings gets privacy settings
	GetPrivacySettings(ctx context.Context, userID uuid.UUID) (*PrivacySettings, error)

	// UpdatePrivacySettings updates privacy settings
	UpdatePrivacySettings(ctx context.Context, userID uuid.UUID, settings *PrivacySettings) error

	// GetMessageSettings gets message settings
	GetMessageSettings(ctx context.Context, userID uuid.UUID) (*MessageSettings, error)

	// UpdateMessageSettings updates message settings
	UpdateMessageSettings(ctx context.Context, userID uuid.UUID, settings *MessageSettings) error
}

// ServerSettingsRepository defines the interface for server settings data access
type ServerSettingsRepository interface {
	// Get retrieves server settings
	Get(ctx context.Context, serverID uuid.UUID) (*ServerSettings, error)

	// Update updates server settings
	Update(ctx context.Context, settings *ServerSettings) error

	// Create creates new server settings
	Create(ctx context.Context, settings *ServerSettings) error

	// GetModerationSettings gets moderation settings
	GetModerationSettings(ctx context.Context, serverID uuid.UUID) (*ModerationSettings, error)

	// UpdateModerationSettings updates moderation settings
	UpdateModerationSettings(ctx context.Context, serverID uuid.UUID, settings *ModerationSettings) error

	// GetInviteSettings gets invite settings
	GetInviteSettings(ctx context.Context, serverID uuid.UUID) (*InviteSettings, error)

	// UpdateInviteSettings updates invite settings
	UpdateInviteSettings(ctx context.Context, serverID uuid.UUID, settings *InviteSettings) error
}

// UserBlockRepository defines the interface for user block data access
type UserBlockRepository interface {
	// Create creates a new user block
	Create(ctx context.Context, block *UserBlock) error

	// Delete deletes a user block
	Delete(ctx context.Context, blockerID, blockedID uuid.UUID) error

	// Get retrieves a user block
	Get(ctx context.Context, blockerID, blockedID uuid.UUID) (*UserBlock, error)

	// ListBlockedUsers lists users blocked by a user
	ListBlockedUsers(ctx context.Context, blockerID uuid.UUID) ([]*User, error)

	// ListBlockedBy lists users who have blocked a user
	ListBlockedBy(ctx context.Context, blockedID uuid.UUID) ([]*User, error)

	// IsBlocked checks if a user is blocked by another user
	IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error)
}

// UserMuteRepository defines the interface for user mute data access
type UserMuteRepository interface {
	// Create creates a new user mute
	Create(ctx context.Context, mute *UserMute) error

	// Delete deletes a user mute
	Delete(ctx context.Context, userID, channelID uuid.UUID) error

	// Get retrieves a user mute
	Get(ctx context.Context, userID, channelID uuid.UUID) (*UserMute, error)

	// ListUserMutes lists channels muted by a user
	ListUserMutes(ctx context.Context, userID uuid.UUID) ([]*Channel, error)

	// IsMuted checks if a user has muted a channel
	IsMuted(ctx context.Context, userID, channelID uuid.UUID) (bool, error)

	// CleanupExpired removes expired mutes
	CleanupExpired(ctx context.Context) error
}

// ModerationRuleRepository defines the interface for moderation rule data access
type ModerationRuleRepository interface {
	// Create creates a new moderation rule
	Create(ctx context.Context, rule *ModerationRule) error

	// Update updates an existing moderation rule
	Update(ctx context.Context, rule *ModerationRule) error

	// Delete deletes a moderation rule
	Delete(ctx context.Context, id uuid.UUID) error

	// Get retrieves a moderation rule by ID
	Get(ctx context.Context, id uuid.UUID) (*ModerationRule, error)

	// List retrieves all moderation rules for a server
	List(ctx context.Context, serverID uuid.UUID) ([]*ModerationRule, error)

	// GetActiveRules retrieves active moderation rules
	GetActiveRules(ctx context.Context, serverID uuid.UUID) ([]*ModerationRule, error)

	// CheckMessage checks a message against all rules
	CheckMessage(ctx context.Context, message *Message) ([]*ModerationRule, error)
}

// ModerationViolationRepository defines the interface for moderation violation data access
type ModerationViolationRepository interface {
	// Create creates a new moderation violation
	Create(ctx context.Context, violation *ModerationViolation) error

	// Update updates an existing moderation violation
	Update(ctx context.Context, violation *ModerationViolation) error

	// Get retrieves a moderation violation by ID
	Get(ctx context.Context, id uuid.UUID) (*ModerationViolation, error)

	// ListUserViolations retrieves violations for a user
	ListUserViolations(ctx context.Context, userID uuid.UUID) ([]*ModerationViolation, error)

	// ListChannelViolations retrieves violations in a channel
	ListChannelViolations(ctx context.Context, channelID uuid.UUID) ([]*ModerationViolation, error)

	// GetUserWarningCount gets the number of warnings for a user
	GetUserWarningCount(ctx context.Context, userID uuid.UUID) (int, error)

	// CleanupExpiredViolations removes expired violations
	CleanupExpiredViolations(ctx context.Context) error
}

// ModerationLogRepository defines the interface for moderation log data access
type ModerationLogRepository interface {
	// Create creates a new moderation log entry
	Create(ctx context.Context, log *ModerationLog) error

	// Get retrieves a moderation log entry by ID
	Get(ctx context.Context, id uuid.UUID) (*ModerationLog, error)

	// ListUserLogs retrieves moderation logs for a user
	ListUserLogs(ctx context.Context, userID uuid.UUID) ([]*ModerationLog, error)

	// ListModeratorLogs retrieves moderation logs by a moderator
	ListModeratorLogs(ctx context.Context, moderatorID uuid.UUID) ([]*ModerationLog, error)

	// CleanupExpiredLogs removes expired log entries
	CleanupExpiredLogs(ctx context.Context) error
}

// WordFilterRepository defines the interface for word filter data access
type WordFilterRepository interface {
	// Create creates a new word filter
	Create(ctx context.Context, filter *WordFilter) error

	// Update updates an existing word filter
	Update(ctx context.Context, filter *WordFilter) error

	// Delete deletes a word filter
	Delete(ctx context.Context, id uuid.UUID) error

	// Get retrieves a word filter by ID
	Get(ctx context.Context, id uuid.UUID) (*WordFilter, error)

	// List retrieves all word filters for a server
	List(ctx context.Context, serverID uuid.UUID) ([]*WordFilter, error)

	// GetByScope retrieves word filters by scope
	GetByScope(ctx context.Context, serverID uuid.UUID, scope WordFilterScope) ([]*WordFilter, error)

	// CheckMessage checks a message against all word filters
	CheckMessage(ctx context.Context, message *Message) ([]*WordFilter, error)
}

package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ChannelUseCase defines the business logic for channel operations
type ChannelUseCase interface {
	// Channel operations
	GetChannel(ctx context.Context, channelID string) (*Channel, error)
	ListChannels(ctx context.Context, userID string, pagination *Pagination) ([]Channel, int, error)
	CreateChannel(ctx context.Context, channel *Channel, members []*ChannelMember) (*Channel, error)
	UpdateChannel(ctx context.Context, channel *Channel) (*Channel, error)
	DeleteChannel(ctx context.Context, channelID string) error
	ArchiveChannel(ctx context.Context, channelID string) error
	UnarchiveChannel(ctx context.Context, channelID string) error

	// Direct channel operations
	CreateDirectChannel(ctx context.Context, userID string, targetUserID string) (*Channel, error)
	ListDirectChannels(ctx context.Context, userID string, offset, limit int) ([]Channel, int, error)
	GetUserChannels(ctx context.Context, userID string, offset, limit int) ([]Channel, int, error)

	// Member operations
	IsChannelMember(ctx context.Context, channelID string, userID string) (bool, error)
	IsChannelAdmin(ctx context.Context, channelID string, userID string) (bool, error)
	ListChannelMembers(ctx context.Context, channelID string, pagination *Pagination) ([]ChannelMember, int, error)
	AddChannelMember(ctx context.Context, channelID string, userID string, role string) (*ChannelMember, error)
	RemoveChannelMember(ctx context.Context, channelID string, userID string) error
	UpdateChannelMember(ctx context.Context, channelID string, userID string, role string) (*ChannelMember, error)

	// Read operations
	MarkChannelAsRead(ctx context.Context, channelID string, userID string) error

	// Additional operations for compatibility
	ListMembers(ctx context.Context, channelID string, offset, limit int) ([]ChannelMember, int, error)
	AddMember(ctx context.Context, channelID interface{}, userID string, role string) (*ChannelMember, error)
	RemoveMember(ctx context.Context, channelID string, userID string) error
}

// MessageUseCase defines the business logic for message operations
type MessageUseCase interface {
	// Message operations
	GetMessage(ctx context.Context, messageID string) (*Message, error)
	CreateMessage(ctx context.Context, channelID string, userID string, content string, parentID string, metadata map[string]string) (*Message, error)
	UpdateMessage(ctx context.Context, messageID string, userID string, content string, metadata map[string]string) (*Message, error)
	DeleteMessage(ctx context.Context, messageID string, userID string) error
	GetChannelMessages(ctx context.Context, channelID string, offset, limit int) ([]Message, int, error)

	// Thread operations
	GetMessageThread(ctx context.Context, messageID string, offset, limit int) ([]Message, int, error)
	ReplyToMessage(ctx context.Context, messageID string, userID string, content string) (*Message, error)
	GetThreadSummary(ctx context.Context, messageID string, userID string) (*ThreadSummary, error)
	FollowThread(ctx context.Context, messageID string, userID string, follow bool) error

	// Reaction operations
	AddReaction(ctx context.Context, messageID string, userID string, emoji string) (*Reaction, error)
	RemoveReaction(ctx context.Context, messageID string, userID string, emoji string) error
	GetMessageReactions(ctx context.Context, messageID string) ([]Reaction, error)

	// Mention operations
	AddMentions(ctx context.Context, messageID string, userIDs []string) error
	GetUserMentions(ctx context.Context, userID string, offset, limit int) ([]Message, int, error)

	// Channel operations for message context
	IsChannelMember(ctx context.Context, channelID string, userID string) (bool, error)
	IsChannelAdmin(ctx context.Context, channelID string, userID string) (bool, error)
}

// ReadReceiptUseCase defines the business logic for read receipt operations
type ReadReceiptUseCase interface {
	MarkMessageAsRead(ctx context.Context, messageID string, userID string) error
	MarkChannelAsRead(ctx context.Context, channelID string, userID string) error
	MarkThreadAsRead(ctx context.Context, threadID string, userID string) error
	GetUnreadCounts(ctx context.Context, userID string) (map[string]int, error)
	GetTotalUnreadCount(ctx context.Context, userID string) (int, error)
	MarkAllAsRead(ctx context.Context, userID string) error
	GetChannelReadStatus(ctx context.Context, channelID string, userID string) (*ChannelReadStatus, error)
	GetReadReceipts(ctx context.Context, messageID string) ([]ReadReceipt, error)
}

// MessagingService defines a comprehensive interface for all messaging operations
type MessagingService interface {
	ChannelUseCase
	MessageUseCase
	ReadReceiptUseCase
}

// MessageDeliveryUseCase defines the business logic for message delivery operations
type MessageDeliveryUseCase interface {
	// Update delivery status
	UpdateDeliveryStatus(ctx context.Context, messageID string, userID string, status MessageStatus) error
	GetMessageDeliveryStatus(ctx context.Context, messageID string) ([]*MessageDelivery, error)
	GetUserDeliveryStatus(ctx context.Context, userID string) ([]*MessageDelivery, error)
}

// MessageDraftUseCase defines the business logic for message draft operations
type MessageDraftUseCase interface {
	// Draft operations
	CreateDraft(ctx context.Context, channelID string, userID string, content string) (*MessageDraft, error)
	UpdateDraft(ctx context.Context, draftID string, content string) (*MessageDraft, error)
	DeleteDraft(ctx context.Context, draftID string) error
	GetDraft(ctx context.Context, draftID string) (*MessageDraft, error)
	GetUserDrafts(ctx context.Context, userID string) ([]*MessageDraft, error)
	GetChannelDrafts(ctx context.Context, channelID string) ([]*MessageDraft, error)
}

// TypingIndicatorUseCase defines the business logic for typing indicator operations
type TypingIndicatorUseCase interface {
	// Typing operations
	StartTyping(ctx context.Context, channelID string, userID string) error
	StopTyping(ctx context.Context, channelID string, userID string) error
	GetChannelTyping(ctx context.Context, channelID string) ([]*TypingIndicator, error)
}

// ChannelCategoryUseCase defines the business logic for channel category operations
type ChannelCategoryUseCase interface {
	// Category operations
	CreateCategory(ctx context.Context, name string, description string) (*ChannelCategory, error)
	UpdateCategory(ctx context.Context, categoryID string, name string, description string) (*ChannelCategory, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategory(ctx context.Context, categoryID string) (*ChannelCategory, error)
	ListCategories(ctx context.Context) ([]*ChannelCategory, error)
	AssignChannelToCategory(ctx context.Context, channelID string, categoryID string) error
	RemoveChannelFromCategory(ctx context.Context, channelID string, categoryID string) error
}

// ChannelInviteUseCase defines the business logic for channel invite operations
type ChannelInviteUseCase interface {
	// Invite operations
	CreateInvite(ctx context.Context, channelID string, inviterID string, inviteeID string) (*ChannelInvite, error)
	AcceptInvite(ctx context.Context, inviteID string) error
	RejectInvite(ctx context.Context, inviteID string) error
	GetInvite(ctx context.Context, inviteID string) (*ChannelInvite, error)
	GetUserInvites(ctx context.Context, userID string) ([]*ChannelInvite, error)
	GetChannelInvites(ctx context.Context, channelID string) ([]*ChannelInvite, error)
}

// MessageReportUseCase defines the business logic for message report operations
type MessageReportUseCase interface {
	// Report operations
	CreateReport(ctx context.Context, messageID string, reporterID string, reason string, description string) (*MessageReport, error)
	UpdateReport(ctx context.Context, reportID string, status string) (*MessageReport, error)
	GetReport(ctx context.Context, reportID string) (*MessageReport, error)
	GetMessageReports(ctx context.Context, messageID string) ([]*MessageReport, error)
	GetUserReports(ctx context.Context, userID string) ([]*MessageReport, error)
	GetPendingReports(ctx context.Context) ([]*MessageReport, error)
}

// PinnedMessageUseCase defines the business logic for pinned message operations
type PinnedMessageUseCase interface {
	// Pin operations
	PinMessage(ctx context.Context, messageID string, userID string) error
	UnpinMessage(ctx context.Context, messageID string) error
	GetChannelPinnedMessages(ctx context.Context, channelID string) ([]*PinnedMessage, error)
	GetMessagePin(ctx context.Context, messageID string) (*PinnedMessage, error)
}

// MessageRetentionUseCase defines the business logic for message retention operations
type MessageRetentionUseCase interface {
	// Retention operations
	CreateRetentionPolicy(ctx context.Context, channelID string, retentionDays int) (*MessageRetentionPolicy, error)
	UpdateRetentionPolicy(ctx context.Context, policyID string, retentionDays int) (*MessageRetentionPolicy, error)
	GetChannelPolicy(ctx context.Context, channelID string) (*MessageRetentionPolicy, error)
	ListPolicies(ctx context.Context) ([]*MessageRetentionPolicy, error)
	EnforceRetentionPolicies(ctx context.Context) error
}

// UserSettingsUseCase defines the business logic for user settings operations
type UserSettingsUseCase interface {
	// Get settings
	GetUserSettings(ctx context.Context, userID string) (*UserSettings, error)
	GetNotificationSettings(ctx context.Context, userID string) (*NotificationSettings, error)
	GetPrivacySettings(ctx context.Context, userID string) (*PrivacySettings, error)
	GetMessageSettings(ctx context.Context, userID string) (*MessageSettings, error)

	// Update settings
	UpdateUserSettings(ctx context.Context, userID string, settings *UserSettings) error
	UpdateNotificationSettings(ctx context.Context, userID string, settings *NotificationSettings) error
	UpdatePrivacySettings(ctx context.Context, userID string, settings *PrivacySettings) error
	UpdateMessageSettings(ctx context.Context, userID string, settings *MessageSettings) error

	// Reset settings
	ResetUserSettings(ctx context.Context, userID string) error
	ResetNotificationSettings(ctx context.Context, userID string) error
	ResetPrivacySettings(ctx context.Context, userID string) error
	ResetMessageSettings(ctx context.Context, userID string) error
}

// ServerSettingsUseCase defines the business logic for server settings operations
type ServerSettingsUseCase interface {
	// Get settings
	GetServerSettings(ctx context.Context, serverID string) (*ServerSettings, error)
	GetModerationSettings(ctx context.Context, serverID string) (*ModerationSettings, error)
	GetInviteSettings(ctx context.Context, serverID string) (*InviteSettings, error)

	// Update settings
	UpdateServerSettings(ctx context.Context, serverID string, settings *ServerSettings) error
	UpdateModerationSettings(ctx context.Context, serverID string, settings *ModerationSettings) error
	UpdateInviteSettings(ctx context.Context, serverID string, settings *InviteSettings) error

	// Reset settings
	ResetServerSettings(ctx context.Context, serverID string) error
	ResetModerationSettings(ctx context.Context, serverID string) error
	ResetInviteSettings(ctx context.Context, serverID string) error
}

// UserBlockUseCase defines the business logic for user block operations
type UserBlockUseCase interface {
	// Block operations
	BlockUser(ctx context.Context, blockerID string, blockedID string, reason string) error
	UnblockUser(ctx context.Context, blockerID string, blockedID string) error
	GetBlock(ctx context.Context, blockerID string, blockedID string) (*UserBlock, error)
	ListBlockedUsers(ctx context.Context, blockerID string) ([]*User, error)
	ListBlockedBy(ctx context.Context, blockedID string) ([]*User, error)
	IsBlocked(ctx context.Context, blockerID string, blockedID string) (bool, error)
}

// UserMuteUseCase defines the business logic for user mute operations
type UserMuteUseCase interface {
	// Mute operations
	MuteChannel(ctx context.Context, userID string, channelID string, duration *time.Duration, reason string) error
	UnmuteChannel(ctx context.Context, userID string, channelID string) error
	GetMute(ctx context.Context, userID string, channelID string) (*UserMute, error)
	ListMutedChannels(ctx context.Context, userID string) ([]*Channel, error)
	IsMuted(ctx context.Context, userID string, channelID string) (bool, error)
}

// EnhancedMessagingService defines a comprehensive interface for all enhanced messaging operations
type EnhancedMessagingService interface {
	MessagingService
	MessageDeliveryUseCase
	MessageDraftUseCase
	TypingIndicatorUseCase
	ChannelCategoryUseCase
	ChannelInviteUseCase
	MessageReportUseCase
	PinnedMessageUseCase
	MessageRetentionUseCase
	UserSettingsUseCase
	ServerSettingsUseCase
	UserBlockUseCase
	UserMuteUseCase
}

// ModerationUseCase defines the interface for moderation operations
type ModerationUseCase interface {
	// CreateRule creates a new moderation rule
	CreateRule(ctx context.Context, rule *ModerationRule) error

	// UpdateRule updates an existing moderation rule
	UpdateRule(ctx context.Context, rule *ModerationRule) error

	// DeleteRule deletes a moderation rule
	DeleteRule(ctx context.Context, id uuid.UUID) error

	// GetRule retrieves a moderation rule by ID
	GetRule(ctx context.Context, id uuid.UUID) (*ModerationRule, error)

	// ListRules retrieves all moderation rules for a server
	ListRules(ctx context.Context, serverID uuid.UUID) ([]*ModerationRule, error)

	// CheckMessage checks a message against all rules and applies appropriate actions
	CheckMessage(ctx context.Context, message *Message) error

	// GetUserViolations retrieves violations for a user
	GetUserViolations(ctx context.Context, userID uuid.UUID) ([]*ModerationViolation, error)

	// GetChannelViolations retrieves violations in a channel
	GetChannelViolations(ctx context.Context, channelID uuid.UUID) ([]*ModerationViolation, error)

	// GetUserWarningCount gets the number of warnings for a user
	GetUserWarningCount(ctx context.Context, userID uuid.UUID) (int, error)

	// GetModerationLogs retrieves moderation logs
	GetModerationLogs(ctx context.Context, filter *ModerationLogFilter) ([]*ModerationLog, error)

	// CleanupExpiredData removes expired violations and logs
	CleanupExpiredData(ctx context.Context) error
}

// ModerationLogFilter defines the filter for retrieving moderation logs
type ModerationLogFilter struct {
	UserID      *uuid.UUID
	ModeratorID *uuid.UUID
	Action      *ModerationActionType
	StartTime   *time.Time
	EndTime     *time.Time
	Limit       int
	Offset      int
}

// WordFilterUseCase defines the interface for word filter operations
type WordFilterUseCase interface {
	// CreateFilter creates a new word filter
	CreateFilter(ctx context.Context, filter *WordFilter) error

	// UpdateFilter updates an existing word filter
	UpdateFilter(ctx context.Context, filter *WordFilter) error

	// DeleteFilter deletes a word filter
	DeleteFilter(ctx context.Context, id uuid.UUID) error

	// GetFilter retrieves a word filter by ID
	GetFilter(ctx context.Context, id uuid.UUID) (*WordFilter, error)

	// ListFilters retrieves all word filters for a server
	ListFilters(ctx context.Context, serverID uuid.UUID) ([]*WordFilter, error)

	// GetFiltersByScope retrieves word filters by scope
	GetFiltersByScope(ctx context.Context, serverID uuid.UUID, scope WordFilterScope) ([]*WordFilter, error)

	// CheckMessage checks a message against all word filters
	CheckMessage(ctx context.Context, message *Message) error
}

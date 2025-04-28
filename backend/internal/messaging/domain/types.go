package domain

import (
	"time"

	"github.com/google/uuid"
)

// ChannelType represents the type of channel
type ChannelType string

const (
	// ChannelTypePublic represents a public channel visible to all users
	ChannelTypePublic ChannelType = "public"

	// ChannelTypePrivate represents a private channel visible only to specific users
	ChannelTypePrivate ChannelType = "private"

	// ChannelTypeDirect represents a direct message channel between two users
	ChannelTypeDirect ChannelType = "direct"
)

// Channel represents a messaging channel
type Channel struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	Type        ChannelType      `json:"type" db:"type"`
	CreatedBy   uuid.UUID        `json:"created_by" db:"created_by"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	IsArchived  bool             `json:"is_archived" db:"is_archived"`
	Members     []*ChannelMember `json:"members,omitempty" db:"-"`
}

// ChannelMemberRole defines the role of a channel member
type ChannelMemberRole string

const (
	// ChannelMemberRoleOwner represents the channel owner
	ChannelMemberRoleOwner ChannelMemberRole = "owner"

	// ChannelMemberRoleAdmin represents a channel admin
	ChannelMemberRoleAdmin ChannelMemberRole = "admin"

	// ChannelMemberRoleMember represents a regular channel member
	ChannelMemberRoleMember ChannelMemberRole = "member"
)

// ChannelMember represents a user's membership in a channel
type ChannelMember struct {
	ChannelID  uuid.UUID  `json:"channel_id" db:"channel_id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	JoinedAt   time.Time  `json:"joined_at" db:"joined_at"`
	Role       string     `json:"role" db:"role"`
	IsAdmin    bool       `json:"is_admin" db:"is_admin"`
	LastReadAt time.Time  `json:"last_read_at" db:"last_read_at"`
	LastReadID uuid.UUID  `json:"last_read_id,omitempty" db:"last_read_id"`
	IsMuted    bool       `json:"is_muted" db:"is_muted"`
	MutedUntil *time.Time `json:"muted_until,omitempty" db:"muted_until"`
}

// MessageType represents the type of message
type MessageType string

const (
	// MessageTypeText represents a regular text message
	MessageTypeText MessageType = "text"

	// MessageTypeSystem represents a system message
	MessageTypeSystem MessageType = "system"

	// MessageTypeAction represents an action message (e.g., "/me is typing")
	MessageTypeAction MessageType = "action"

	// MessageTypeFile represents a file attachment message
	MessageTypeFile MessageType = "file"

	// MessageTypeReaction represents a reaction message
	MessageTypeReaction MessageType = "reaction"
)

// Message represents a message in a channel
type Message struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	ChannelID   uuid.UUID    `json:"channel_id" db:"channel_id"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	Content     string       `json:"content" db:"content"`
	Type        MessageType  `json:"type" db:"type"`
	ReplyToID   *uuid.UUID   `json:"reply_to_id,omitempty" db:"reply_to_id"`
	EditedAt    *time.Time   `json:"edited_at,omitempty" db:"edited_at"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	IsDeleted   bool         `json:"is_deleted" db:"is_deleted"`
	Attachments []Attachment `json:"attachments,omitempty" db:"-"`
	Reactions   []Reaction   `json:"reactions,omitempty" db:"-"`
	ThreadCount int          `json:"thread_count,omitempty" db:"-"`
}

// Attachment represents a file attachment on a message
type Attachment struct {
	ID           uuid.UUID `json:"id" db:"id"`
	MessageID    uuid.UUID `json:"message_id" db:"message_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	FileName     string    `json:"file_name" db:"file_name"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	FileType     string    `json:"file_type" db:"file_type"`
	FileURL      string    `json:"file_url" db:"file_url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Mention represents a user mention in a message
type Mention struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Reaction represents a reaction to a message
type Reaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Emoji     string    `json:"emoji" db:"emoji"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ReactionCount represents the count of a specific reaction
type ReactionCount struct {
	Emoji      string   `json:"emoji" db:"emoji"`
	Count      int      `json:"count" db:"count"`
	Users      []string `json:"users" db:"users"`
	HasReacted bool     `json:"has_reacted" db:"has_reacted"`
}

// ReadReceipt represents a read receipt for a message
type ReadReceipt struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChannelID uuid.UUID `json:"channel_id" db:"channel_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	ReadAt    time.Time `json:"read_at" db:"read_at"`
}

// ChannelReadStatus represents the read status of a channel for a specific user
type ChannelReadStatus struct {
	ChannelID    uuid.UUID `json:"channel_id" db:"channel_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	LastReadID   uuid.UUID `json:"last_read_id" db:"last_read_id"`
	LastReadAt   time.Time `json:"last_read_at" db:"last_read_at"`
	UnreadCount  int       `json:"unread_count" db:"unread_count"`
	MentionCount int       `json:"mention_count" db:"mention_count"`
}

// MessagingChannel is an alias for Channel to maintain compatibility
type MessagingChannel Channel

// MessagingMember is an alias for ChannelMember to maintain compatibility
type MessagingMember ChannelMember

// MessageStatus represents the status of a message
type MessageStatus string

const (
	// MessageStatusSent represents a sent message
	MessageStatusSent MessageStatus = "sent"
	// MessageStatusDelivered represents a delivered message
	MessageStatusDelivered MessageStatus = "delivered"
	// MessageStatusRead represents a read message
	MessageStatusRead MessageStatus = "read"
	// MessageStatusFailed represents a failed message
	MessageStatusFailed MessageStatus = "failed"
)

// MessageDelivery represents the delivery status of a message
type MessageDelivery struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	MessageID uuid.UUID     `json:"message_id" db:"message_id"`
	UserID    uuid.UUID     `json:"user_id" db:"user_id"`
	Status    MessageStatus `json:"status" db:"status"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

// MessageDraft represents a draft message
type MessageDraft struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChannelID uuid.UUID `json:"channel_id" db:"channel_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TypingIndicator represents a user's typing status
type TypingIndicator struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChannelID uuid.UUID `json:"channel_id" db:"channel_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	StartedAt time.Time `json:"started_at" db:"started_at"`
}

// ChannelCategory represents a category for organizing channels
type ChannelCategory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Order       int       `json:"order" db:"order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ChannelInvite represents an invitation to join a channel
type ChannelInvite struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	ChannelID uuid.UUID  `json:"channel_id" db:"channel_id"`
	InviterID uuid.UUID  `json:"inviter_id" db:"inviter_id"`
	InviteeID uuid.UUID  `json:"invitee_id" db:"invitee_id"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// MessageReport represents a report of inappropriate message content
type MessageReport struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	MessageID   uuid.UUID  `json:"message_id" db:"message_id"`
	ReporterID  uuid.UUID  `json:"reporter_id" db:"reporter_id"`
	Reason      string     `json:"reason" db:"reason"`
	Description string     `json:"description" db:"description"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
}

// PinnedMessage represents a pinned message in a channel
type PinnedMessage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	ChannelID uuid.UUID `json:"channel_id" db:"channel_id"`
	PinnedBy  uuid.UUID `json:"pinned_by" db:"pinned_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MessageRetentionPolicy represents a policy for message retention
type MessageRetentionPolicy struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ChannelID     uuid.UUID `json:"channel_id" db:"channel_id"`
	RetentionDays int       `json:"retention_days" db:"retention_days"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SearchScope defines where to search
type SearchScope string

const (
	// SearchScopeAll searches across all channels
	SearchScopeAll SearchScope = "all"
	// SearchScopeChannel searches within a specific channel
	SearchScopeChannel SearchScope = "channel"
	// SearchScopeUser searches within messages from a specific user
	SearchScopeUser SearchScope = "user"
	// SearchScopeThread searches within a specific thread
	SearchScopeThread SearchScope = "thread"
)

// UserSettings represents user-specific messaging settings
type UserSettings struct {
	ID                   uuid.UUID             `json:"id" db:"id"`
	UserID               uuid.UUID             `json:"user_id" db:"user_id"`
	Theme                string                `json:"theme" db:"theme"`
	Language             string                `json:"language" db:"language"`
	NotificationSettings *NotificationSettings `json:"notification_settings" db:"notification_settings"`
	PrivacySettings      *PrivacySettings      `json:"privacy_settings" db:"privacy_settings"`
	MessageSettings      *MessageSettings      `json:"message_settings" db:"message_settings"`
	CreatedAt            time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at" db:"updated_at"`
}

// NotificationSettings represents user notification preferences
type NotificationSettings struct {
	DesktopNotifications bool       `json:"desktop_notifications" db:"desktop_notifications"`
	EmailNotifications   bool       `json:"email_notifications" db:"email_notifications"`
	MobileNotifications  bool       `json:"mobile_notifications" db:"mobile_notifications"`
	MentionNotifications bool       `json:"mention_notifications" db:"mention_notifications"`
	ThreadNotifications  bool       `json:"thread_notifications" db:"thread_notifications"`
	SoundEnabled         bool       `json:"sound_enabled" db:"sound_enabled"`
	DoNotDisturb         bool       `json:"do_not_disturb" db:"do_not_disturb"`
	DoNotDisturbStart    *time.Time `json:"do_not_disturb_start,omitempty" db:"do_not_disturb_start"`
	DoNotDisturbEnd      *time.Time `json:"do_not_disturb_end,omitempty" db:"do_not_disturb_end"`
}

// PrivacySettings represents user privacy preferences
type PrivacySettings struct {
	OnlineStatusVisible bool     `json:"online_status_visible" db:"online_status_visible"`
	TypingStatusVisible bool     `json:"typing_status_visible" db:"typing_status_visible"`
	ReadReceiptsEnabled bool     `json:"read_receipts_enabled" db:"read_receipts_enabled"`
	ProfileVisible      bool     `json:"profile_visible" db:"profile_visible"`
	LastSeenVisible     bool     `json:"last_seen_visible" db:"last_seen_visible"`
	BlockedUsers        []string `json:"blocked_users" db:"blocked_users"`
}

// MessageSettings represents user message preferences
type MessageSettings struct {
	AutoDeleteMessages bool `json:"auto_delete_messages" db:"auto_delete_messages"`
	AutoDeleteAfter    int  `json:"auto_delete_after" db:"auto_delete_after"`
	MessagePreview     bool `json:"message_preview" db:"message_preview"`
	LinkPreview        bool `json:"link_preview" db:"link_preview"`
	EmojiReactions     bool `json:"emoji_reactions" db:"emoji_reactions"`
	MessageFormatting  bool `json:"message_formatting" db:"message_formatting"`
}

// ServerSettings represents server-wide messaging settings
type ServerSettings struct {
	ID                   uuid.UUID           `json:"id" db:"id"`
	ServerID             uuid.UUID           `json:"server_id" db:"server_id"`
	MessageRetentionDays int                 `json:"message_retention_days" db:"message_retention_days"`
	MaxFileSize          int64               `json:"max_file_size" db:"max_file_size"`
	AllowedFileTypes     []string            `json:"allowed_file_types" db:"allowed_file_types"`
	MaxMessageLength     int                 `json:"max_message_length" db:"max_message_length"`
	ModerationSettings   *ModerationSettings `json:"moderation_settings" db:"moderation_settings"`
	InviteSettings       *InviteSettings     `json:"invite_settings" db:"invite_settings"`
	CreatedAt            time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at" db:"updated_at"`
}

// ModerationActionType defines the type of action to take for moderated content
type ModerationActionType string

const (
	// ModerationActionWarn issues a warning to the user
	ModerationActionWarn ModerationActionType = "warn"
	// ModerationActionDelete deletes the message
	ModerationActionDelete ModerationActionType = "delete"
	// ModerationActionMute temporarily mutes the user
	ModerationActionMute ModerationActionType = "mute"
	// ModerationActionBan permanently bans the user
	ModerationActionBan ModerationActionType = "ban"
)

// ModerationLevel defines the severity level of moderation
type ModerationLevel string

const (
	// ModerationLevelLow indicates minor violations
	ModerationLevelLow ModerationLevel = "low"
	// ModerationLevelMedium indicates moderate violations
	ModerationLevelMedium ModerationLevel = "medium"
	// ModerationLevelHigh indicates severe violations
	ModerationLevelHigh ModerationLevel = "high"
)

// ModerationRule represents a rule for content moderation
type ModerationRule struct {
	ID          uuid.UUID            `json:"id" db:"id"`
	ServerID    uuid.UUID            `json:"server_id" db:"server_id"`
	Pattern     string               `json:"pattern" db:"pattern"`
	Description string               `json:"description" db:"description"`
	Action      ModerationActionType `json:"action" db:"action"`
	Level       ModerationLevel      `json:"level" db:"level"`
	Duration    *time.Duration       `json:"duration,omitempty" db:"duration"`
	CreatedBy   uuid.UUID            `json:"created_by" db:"created_by"`
	CreatedAt   time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" db:"updated_at"`
	IsActive    bool                 `json:"is_active" db:"is_active"`
}

// ModerationViolation represents a violation of moderation rules
type ModerationViolation struct {
	ID            uuid.UUID            `json:"id" db:"id"`
	MessageID     uuid.UUID            `json:"message_id" db:"message_id"`
	UserID        uuid.UUID            `json:"user_id" db:"user_id"`
	ChannelID     uuid.UUID            `json:"channel_id" db:"channel_id"`
	RuleID        uuid.UUID            `json:"rule_id" db:"rule_id"`
	ViolationText string               `json:"violation_text" db:"violation_text"`
	Action        ModerationActionType `json:"action" db:"action"`
	ActionTaken   bool                 `json:"action_taken" db:"action_taken"`
	CreatedAt     time.Time            `json:"created_at" db:"created_at"`
	ResolvedAt    *time.Time           `json:"resolved_at,omitempty" db:"resolved_at"`
}

// ModerationLog represents a log entry for moderation actions
type ModerationLog struct {
	ID          uuid.UUID            `json:"id" db:"id"`
	UserID      uuid.UUID            `json:"user_id" db:"user_id"`
	ModeratorID uuid.UUID            `json:"moderator_id" db:"moderator_id"`
	Action      ModerationActionType `json:"action" db:"action"`
	Reason      string               `json:"reason" db:"reason"`
	Duration    *time.Duration       `json:"duration,omitempty" db:"duration"`
	CreatedAt   time.Time            `json:"created_at" db:"created_at"`
	ExpiresAt   *time.Time           `json:"expires_at,omitempty" db:"expires_at"`
}

// ModerationSettings represents server moderation settings
type ModerationSettings struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	ServerID             uuid.UUID       `json:"server_id" db:"server_id"`
	AutoModEnabled       bool            `json:"auto_mod_enabled" db:"auto_mod_enabled"`
	AutoModLevel         ModerationLevel `json:"auto_mod_level" db:"auto_mod_level"`
	MaxWarnings          int             `json:"max_warnings" db:"max_warnings"`
	WarningExpiryDays    int             `json:"warning_expiry_days" db:"warning_expiry_days"`
	MuteDuration         time.Duration   `json:"mute_duration" db:"mute_duration"`
	BanDuration          *time.Duration  `json:"ban_duration,omitempty" db:"ban_duration"`
	LogModerationActions bool            `json:"log_moderation_actions" db:"log_moderation_actions"`
	NotifyOnViolation    bool            `json:"notify_on_violation" db:"notify_on_violation"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// InviteSettings represents server invite settings
type InviteSettings struct {
	InviteOnly          bool `json:"invite_only" db:"invite_only"`
	MaxInvitesPerUser   int  `json:"max_invites_per_user" db:"max_invites_per_user"`
	InviteExpiryDays    int  `json:"invite_expiry_days" db:"invite_expiry_days"`
	RequireVerification bool `json:"require_verification" db:"require_verification"`
}

// UserBlock represents a user blocking another user
type UserBlock struct {
	ID        uuid.UUID `json:"id" db:"id"`
	BlockerID uuid.UUID `json:"blocker_id" db:"blocker_id"`
	BlockedID uuid.UUID `json:"blocked_id" db:"blocked_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Reason    string    `json:"reason,omitempty" db:"reason"`
}

// UserMute represents a user muting a channel
type UserMute struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	ChannelID uuid.UUID  `json:"channel_id" db:"channel_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	Reason    string     `json:"reason,omitempty" db:"reason"`
}

// WordFilterScope defines where the word filter applies
type WordFilterScope string

const (
	// WordFilterScopePrivate applies to private messages
	WordFilterScopePrivate WordFilterScope = "private"
	// WordFilterScopePublic applies to public messages
	WordFilterScopePublic WordFilterScope = "public"
	// WordFilterScopeAll applies to all messages
	WordFilterScopeAll WordFilterScope = "all"
)

// WordFilter represents a word or phrase to be filtered
type WordFilter struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	ServerID    uuid.UUID       `json:"server_id" db:"server_id"`
	Word        string          `json:"word" db:"word"`
	Scope       WordFilterScope `json:"scope" db:"scope"`
	Description string          `json:"description" db:"description"`
	CreatedBy   uuid.UUID       `json:"created_by" db:"created_by"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	IsActive    bool            `json:"is_active" db:"is_active"`
}

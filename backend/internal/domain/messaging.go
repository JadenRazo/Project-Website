package domain

import (
	"time"

	"gorm.io/gorm"
)

// MessagingMessage represents a message in the messaging system
type MessagingMessage struct {
	gorm.Model
	ChannelID    uint   `gorm:"not null"`
	SenderID     uint   `gorm:"not null"`
	Content      string `gorm:"type:text"`
	Mentions     []uint `gorm:"type:integer[]"`
	ReplyToID    *uint
	IsPinned     bool `gorm:"default:false"`
	PinnedAt     *time.Time
	PinnedBy     *uint
	IsNSFW       bool `gorm:"default:false"`
	IsSpoiler    bool `gorm:"default:false"`
	CustomEmoji  string
	ThreadID     *uint
	MessageFlags uint64 `gorm:"default:0"`

	// Relationships
	Channel      MessagingChannel      `gorm:"foreignKey:ChannelID"`
	Sender       User                  `gorm:"foreignKey:SenderID"`
	Attachments  []MessagingAttachment `gorm:"foreignKey:MessageID"`
	Embeds       []MessagingEmbed      `gorm:"foreignKey:MessageID"`
	Reactions    []MessagingReaction   `gorm:"foreignKey:MessageID"`
	Replies      []MessagingMessage    `gorm:"foreignKey:ReplyToID"`
	Thread       *MessagingMessage     `gorm:"foreignKey:ThreadID"`
	PinnedByUser *User                 `gorm:"foreignKey:PinnedBy"`
}

// MessagingChannel represents a channel in the messaging system
type MessagingChannel struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Description  string
	Type         string `gorm:"not null"` // direct, group, public, private
	Icon         string
	Topic        string
	IsArchived   bool   `gorm:"default:false"`
	IsNSFW       bool   `gorm:"default:false"`
	IsPrivate    bool   `gorm:"default:false"`
	IsReadOnly   bool   `gorm:"default:false"`
	ChannelFlags uint64 `gorm:"default:0"`

	// Relationships
	Members        []MessagingChannelMember `gorm:"foreignKey:ChannelID"`
	PinnedMessages []MessagingPinnedMessage `gorm:"foreignKey:ChannelID"`
}

// MessagingChannelMember represents a user's membership in a channel
type MessagingChannelMember struct {
	gorm.Model
	ChannelID   uint   `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	Role        string `gorm:"not null"` // owner, admin, moderator, member
	IsMuted     bool   `gorm:"default:false"`
	IsBanned    bool   `gorm:"default:false"`
	MemberFlags uint64 `gorm:"default:0"`

	// Relationships
	Channel MessagingChannel `gorm:"foreignKey:ChannelID"`
	User    User             `gorm:"foreignKey:UserID"`
}

// MessagingPinnedMessage represents a pinned message in a channel
type MessagingPinnedMessage struct {
	gorm.Model
	ChannelID uint `gorm:"not null"`
	MessageID uint `gorm:"not null"`
	PinnedBy  uint `gorm:"not null"`

	// Relationships
	Channel      MessagingChannel `gorm:"foreignKey:ChannelID"`
	Message      MessagingMessage `gorm:"foreignKey:MessageID"`
	PinnedByUser User             `gorm:"foreignKey:PinnedBy"`
}

// MessagingAttachment represents a file attachment in a message
type MessagingAttachment struct {
	gorm.Model
	MessageID       uint   `gorm:"not null"`
	URL             string `gorm:"not null"`
	Filename        string `gorm:"not null"`
	ContentType     string `gorm:"not null"`
	Size            int64  `gorm:"not null"`
	Hash            string
	Width           *int
	Height          *int
	Duration        *int64
	IsNSFW          bool   `gorm:"default:false"`
	IsSpoiler       bool   `gorm:"default:false"`
	AttachmentFlags uint64 `gorm:"default:0"`

	// Relationships
	Message MessagingMessage `gorm:"foreignKey:MessageID"`
}

// MessagingEmbed represents an embedded content in a message
type MessagingEmbed struct {
	gorm.Model
	MessageID   uint   `gorm:"not null"`
	URL         string `gorm:"not null"`
	Type        string `gorm:"not null"` // link, image, video, audio, etc.
	Title       string
	Description string
	Thumbnail   string
	Width       *int
	Height      *int
	Duration    *int64
	IsNSFW      bool   `gorm:"default:false"`
	IsSpoiler   bool   `gorm:"default:false"`
	EmbedFlags  uint64 `gorm:"default:0"`

	// Relationships
	Message MessagingMessage `gorm:"foreignKey:MessageID"`
}

// MessagingReaction represents a reaction on a message
type MessagingReaction struct {
	gorm.Model
	MessageID     uint   `gorm:"not null"`
	UserID        uint   `gorm:"not null"`
	EmojiCode     string `gorm:"not null"`
	ReactionFlags uint64 `gorm:"default:0"`

	// Relationships
	Message MessagingMessage `gorm:"foreignKey:MessageID"`
	User    User             `gorm:"foreignKey:UserID"`
}

// MessagingReadReceipt represents a read receipt for a message
type MessagingReadReceipt struct {
	gorm.Model
	MessageID uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	ReadAt    time.Time `gorm:"not null"`

	// Relationships
	Message MessagingMessage `gorm:"foreignKey:MessageID"`
	User    User             `gorm:"foreignKey:UserID"`
}

// User represents a user in the system
type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Status    string `gorm:"default:'offline'"`
	UserFlags uint64 `gorm:"default:0"`
}

// MessagingCategory represents a channel category
type MessagingCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ServerID    uint      `json:"serverId"`
	Position    int       `json:"position"` // For ordering
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt,omitempty" gorm:"index"`

	Channels []MessagingChannel `json:"channels,omitempty" gorm:"foreignKey:CategoryID"`
}

// MessagingRole represents a user role with permissions
type MessagingRole struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Permissions uint64    `json:"permissions"` // Bitfield of permissions
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt,omitempty" gorm:"index"`

	Users []MessagingUser `json:"users,omitempty" gorm:"many2many:user_roles;"`
}

// MessagingUser represents a messaging system user
type MessagingUser struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"-"` // Hashed password
	Avatar        string    `json:"avatar"`
	Status        string    `json:"status"` // online, idle, dnd, offline
	StatusMessage string    `json:"statusMessage"`
	LastSeen      time.Time `json:"lastSeen"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt,omitempty" gorm:"index"`

	// New fields
	IsVerified       bool               `json:"isVerified"`
	TwoFactorEnabled bool               `json:"twoFactorEnabled"`
	BlockedUsers     []MessagingUser    `json:"blockedUsers,omitempty" gorm:"many2many:user_blocks;"`
	MutedUsers       []MessagingUser    `json:"mutedUsers,omitempty" gorm:"many2many:user_mutes;"`
	Roles            []MessagingRole    `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	Channels         []MessagingChannel `json:"channels,omitempty" gorm:"many2many:channel_members;"`
}

// Permission constants
const (
	PermissionReadMessages   uint64 = 1 << 0
	PermissionSendMessages   uint64 = 1 << 1
	PermissionManageChannel  uint64 = 1 << 2
	PermissionManageRoles    uint64 = 1 << 3
	PermissionManageMessages uint64 = 1 << 4
	PermissionInviteUsers    uint64 = 1 << 5
	PermissionBanUsers       uint64 = 1 << 6
	PermissionAdmin          uint64 = 1 << 7
)

// Channel types
const (
	ChannelTypeText         = "text"
	ChannelTypeVoice        = "voice"
	ChannelTypeAnnouncement = "announcement"
)

// User statuses
const (
	UserStatusOnline  = "online"
	UserStatusIdle    = "idle"
	UserStatusDND     = "dnd"
	UserStatusOffline = "offline"
)

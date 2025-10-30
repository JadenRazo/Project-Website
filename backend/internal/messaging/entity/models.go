package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Channel represents a messaging channel
type Channel struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description,omitempty" gorm:"type:text"`
	Type        string         `json:"type" gorm:"size:20;not null;default:'public'"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid;not null"`
	IsArchived  bool           `json:"is_archived" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Associations
	Members  []ChannelMember `json:"members,omitempty" gorm:"foreignKey:ChannelID"`
	Messages []Message       `json:"messages,omitempty" gorm:"foreignKey:ChannelID"`
}

// ChannelMember represents a user's membership in a channel
type ChannelMember struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChannelID uuid.UUID `json:"channel_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Role      string    `json:"role" gorm:"size:20;not null;default:'member'"`
	JoinedAt  time.Time `json:"joined_at" gorm:"autoCreateTime"`

	// Associations
	Channel *Channel `json:"channel,omitempty" gorm:"foreignKey:ChannelID"`
}

// Message represents a chat message
type Message struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChannelID   uuid.UUID      `json:"channel_id" gorm:"type:uuid;not null"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	MessageType string         `json:"message_type" gorm:"size:20;not null;default:'text'"`
	ParentID    *uuid.UUID     `json:"parent_id,omitempty" gorm:"type:uuid"`
	IsEdited    bool           `json:"is_edited" gorm:"default:false"`
	IsDeleted   bool           `json:"is_deleted" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Associations
	Channel      *Channel            `json:"channel,omitempty" gorm:"foreignKey:ChannelID"`
	Parent       *Message            `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies      []Message           `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	Attachments  []MessageAttachment `json:"attachments,omitempty" gorm:"foreignKey:MessageID"`
	Reactions    []MessageReaction   `json:"reactions,omitempty" gorm:"foreignKey:MessageID"`
	ReadReceipts []ReadReceipt       `json:"read_receipts,omitempty" gorm:"foreignKey:MessageID"`
}

// MessageAttachment represents a file attachment
type MessageAttachment struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MessageID        uuid.UUID `json:"message_id" gorm:"type:uuid;not null"`
	Filename         string    `json:"filename" gorm:"size:255;not null"`
	OriginalFilename string    `json:"original_filename" gorm:"size:255;not null"`
	FileSize         int64     `json:"file_size" gorm:"not null"`
	MimeType         string    `json:"mime_type" gorm:"size:100;not null"`
	FilePath         string    `json:"file_path" gorm:"type:text;not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Associations
	Message *Message `json:"message,omitempty" gorm:"foreignKey:MessageID"`
}

// MessageReaction represents an emoji reaction to a message
type MessageReaction struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MessageID uuid.UUID `json:"message_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Emoji     string    `json:"emoji" gorm:"size:10;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Associations
	Message *Message `json:"message,omitempty" gorm:"foreignKey:MessageID"`
}

// ReadReceipt represents when a user read a message
type ReadReceipt struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MessageID uuid.UUID `json:"message_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ReadAt    time.Time `json:"read_at" gorm:"autoCreateTime"`

	// Associations
	Message *Message `json:"message,omitempty" gorm:"foreignKey:MessageID"`
}

// Request/Response DTOs

// CreateChannelRequest represents the request to create a channel
type CreateChannelRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty" binding:"omitempty,oneof=public private direct"`
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required,min=1,max=4000"`
	MessageType string `json:"message_type,omitempty" binding:"omitempty,oneof=text image file system"`
	ParentID    string `json:"parent_id,omitempty"`
}

// AddReactionRequest represents the request to add a reaction
type AddReactionRequest struct {
	Emoji string `json:"emoji" binding:"required,min=1,max=10"`
}

// UpdateMessageRequest represents the request to update a message
type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=4000"`
}

// ChannelResponse represents a channel in API responses
type ChannelResponse struct {
	*Channel
	MemberCount  int `json:"member_count"`
	MessageCount int `json:"message_count"`
}

// MessageResponse represents a message in API responses
type MessageResponse struct {
	*Message
	Username      string `json:"username,omitempty"`
	ReactionCount int    `json:"reaction_count"`
	ReplyCount    int    `json:"reply_count"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// MessageListResponse represents a paginated list of messages
type MessageListResponse struct {
	Messages   []MessageResponse  `json:"messages"`
	Pagination PaginationResponse `json:"pagination"`
}

// ChannelListResponse represents a paginated list of channels
type ChannelListResponse struct {
	Channels   []ChannelResponse  `json:"channels"`
	Pagination PaginationResponse `json:"pagination"`
}

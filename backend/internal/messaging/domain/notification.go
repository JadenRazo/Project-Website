package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	// NotificationTypeNewMessage represents a new message notification
	NotificationTypeNewMessage NotificationType = "new_message"

	// NotificationTypeChannelInvite represents a channel invitation notification
	NotificationTypeChannelInvite NotificationType = "channel_invite"

	// NotificationTypeMention represents a mention notification
	NotificationTypeMention NotificationType = "mention"

	// NotificationTypeReaction represents a reaction notification
	NotificationTypeReaction NotificationType = "reaction"
)

// Notification represents a user notification
type Notification struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	UserID    uuid.UUID        `json:"user_id" db:"user_id"`
	Type      NotificationType `json:"type" db:"type"`
	ChannelID uuid.UUID        `json:"channel_id,omitempty" db:"channel_id"`
	MessageID uuid.UUID        `json:"message_id,omitempty" db:"message_id"`
	SenderID  uuid.UUID        `json:"sender_id" db:"sender_id"`
	Content   string           `json:"content" db:"content"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

// NewNotification creates a new notification
func NewNotification(userID, channelID, messageID, senderID uuid.UUID, notificationType NotificationType, content string) *Notification {
	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notificationType,
		ChannelID: channelID,
		MessageID: messageID,
		SenderID:  senderID,
		Content:   content,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

// NewMessageNotification creates a notification for a new message
func NewMessageNotification(userID, channelID, messageID, senderID uuid.UUID, preview string) *Notification {
	return NewNotification(
		userID,
		channelID,
		messageID,
		senderID,
		NotificationTypeNewMessage,
		preview,
	)
}

// NewMentionNotification creates a notification for a mention
func NewMentionNotification(userID, channelID, messageID, senderID uuid.UUID, preview string) *Notification {
	return NewNotification(
		userID,
		channelID,
		messageID,
		senderID,
		NotificationTypeMention,
		preview,
	)
}

// NewReactionNotification creates a notification for a reaction
func NewReactionNotification(userID, channelID, messageID, senderID uuid.UUID, emoji string) *Notification {
	return NewNotification(
		userID,
		channelID,
		messageID,
		senderID,
		NotificationTypeReaction,
		emoji,
	)
}

// NewChannelInviteNotification creates a notification for a channel invitation
func NewChannelInviteNotification(userID, channelID, senderID uuid.UUID, channelName string) *Notification {
	return NewNotification(
		userID,
		channelID,
		uuid.Nil,
		senderID,
		NotificationTypeChannelInvite,
		channelName,
	)
}

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *Notification) error

	// Get retrieves a notification by ID
	Get(ctx context.Context, id uuid.UUID) (*Notification, error)

	// GetUserNotifications retrieves notifications for a user
	GetUserNotifications(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]*Notification, int, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id uuid.UUID) error

	// MarkAllAsRead marks all notifications for a user as read
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error

	// Delete deletes a notification
	Delete(ctx context.Context, id uuid.UUID) error
}

// NotificationUseCase defines the business logic for notification operations
type NotificationUseCase interface {
	// CreateNotification creates a new notification
	CreateNotification(ctx context.Context, notification *Notification) error

	// GetNotification retrieves a notification by ID
	GetNotification(ctx context.Context, id string) (*Notification, error)

	// GetUserNotifications retrieves notifications for a user
	GetUserNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*Notification, int, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id string) error

	// MarkAllAsRead marks all notifications for a user as read
	MarkAllAsRead(ctx context.Context, userID string) error

	// DeleteNotification deletes a notification
	DeleteNotification(ctx context.Context, id string) error
}

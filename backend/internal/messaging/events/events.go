package events

import (
	"context"
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
)

// MessagingEventType represents the type of event
type MessagingEventType string

// Messaging event constants
const (
	MessagingEventMessageCreated    MessagingEventType = "message.created"
	MessagingEventMessageUpdated    MessagingEventType = "message.updated"
	MessagingEventMessageDeleted    MessagingEventType = "message.deleted"
	MessagingEventMessagePinned     MessagingEventType = "message.pinned"
	MessagingEventMessageUnpinned   MessagingEventType = "message.unpinned"
	MessagingEventAttachmentCreated MessagingEventType = "attachment.created"
	MessagingEventAttachmentDeleted MessagingEventType = "attachment.deleted"
	MessagingEventChannelCreated    MessagingEventType = "channel.created"
	MessagingEventChannelUpdated    MessagingEventType = "channel.updated"
	MessagingEventChannelDeleted    MessagingEventType = "channel.deleted"
	MessagingEventChannelRead       MessagingEventType = "channel.read"
	MessagingEventCategoryCreated   MessagingEventType = "category.created"
	MessagingEventCategoryUpdated   MessagingEventType = "category.updated"
	MessagingEventCategoryDeleted   MessagingEventType = "category.deleted"
	MessagingEventStatusUpdated     MessagingEventType = "status.updated"
	MessagingEventUserBlocked       MessagingEventType = "user.blocked"
	MessagingEventUserUnblocked     MessagingEventType = "user.unblocked"
	MessagingEventUserMuted         MessagingEventType = "user.muted"
	MessagingEventUserUnmuted       MessagingEventType = "user.unmuted"
	MessagingEventRoleCreated       MessagingEventType = "role.created"
	MessagingEventRoleUpdated       MessagingEventType = "role.updated"
	MessagingEventRoleDeleted       MessagingEventType = "role.deleted"
	MessagingEventRoleAdded         MessagingEventType = "role.added"
	MessagingEventRoleRemoved       MessagingEventType = "role.removed"
	MessagingEventReactionAdded     MessagingEventType = "reaction.added"
	MessagingEventReactionRemoved   MessagingEventType = "reaction.removed"
	MessagingEventReadReceipt       MessagingEventType = "read.receipt"
	MessagingEventMemberAdded       MessagingEventType = "member.added"
	MessagingEventMemberRemoved     MessagingEventType = "member.removed"
)

// MessagingEvent defines the interface for all messaging events
type MessagingEvent interface {
	Data() interface{}
}

// MessagingEventData structures
type MessagingMessageEventData struct {
	Message *domain.MessagingMessage
}

type MessagingChannelEventData struct {
	Channel   *domain.MessagingChannel
	ChannelID uint
	UserID    uint
}

type MessagingCategoryEventData struct {
	Category *domain.MessagingCategory
}

type MessagingStatusEventData struct {
	UserID        uint
	Status        string
	StatusMessage string
}

type MessagingUserEventData struct {
	UserID    uint
	BlockedID uint
	MutedID   uint
	Duration  time.Duration
}

type MessagingRoleEventData struct {
	Role        *domain.MessagingRole
	UserID      uint
	RoleID      uint
	AddedByID   uint
	RemovedByID uint
}

type MessagingReactionEventData struct {
	Reaction *domain.MessagingReaction
}

type MessagingReadReceiptEventData struct {
	Receipt *domain.MessagingReadReceipt
}

type MessagingMemberEventData struct {
	ChannelID   uint
	UserID      uint
	AddedByID   uint
	RemovedByID uint
}

// BaseMessagingEvent implements the MessagingEvent interface
type BaseMessagingEvent struct {
	eventType MessagingEventType
	data      interface{}
}

// NewBaseMessagingEvent creates a new base event
func NewBaseMessagingEvent(eventType MessagingEventType, data interface{}) *BaseMessagingEvent {
	return &BaseMessagingEvent{
		eventType: eventType,
		data:      data,
	}
}

// Type returns the event type
func (e *BaseMessagingEvent) Type() MessagingEventType {
	return e.eventType
}

// Data returns the event data
func (e *BaseMessagingEvent) Data() interface{} {
	return e.data
}

// MessageEvent represents a message-related event
type MessageEvent struct {
	*BaseMessagingEvent
}

// NewMessageEvent creates a new message event
func NewMessageEvent(eventType MessagingEventType, message *domain.MessagingMessage) *MessageEvent {
	return &MessageEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingMessageEventData{
			Message: message,
		}),
	}
}

// ChannelEvent represents a channel-related event
type ChannelEvent struct {
	*BaseMessagingEvent
}

// NewChannelEvent creates a new channel event
func NewChannelEvent(eventType MessagingEventType, channel *domain.MessagingChannel, channelID, userID uint) *ChannelEvent {
	return &ChannelEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingChannelEventData{
			Channel:   channel,
			ChannelID: channelID,
			UserID:    userID,
		}),
	}
}

// CategoryEvent represents a category-related event
type CategoryEvent struct {
	*BaseMessagingEvent
}

// NewCategoryEvent creates a new category event
func NewCategoryEvent(eventType MessagingEventType, category *domain.MessagingCategory) *CategoryEvent {
	return &CategoryEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingCategoryEventData{
			Category: category,
		}),
	}
}

// StatusEvent represents a status-related event
type StatusEvent struct {
	*BaseMessagingEvent
}

// NewStatusEvent creates a new status event
func NewStatusEvent(userID uint, status, statusMessage string) *StatusEvent {
	return &StatusEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(MessagingEventStatusUpdated, &MessagingStatusEventData{
			UserID:        userID,
			Status:        status,
			StatusMessage: statusMessage,
		}),
	}
}

// UserEvent represents a user-related event
type UserEvent struct {
	*BaseMessagingEvent
}

// NewUserEvent creates a new user event
func NewUserEvent(eventType MessagingEventType, userID, targetID uint, duration time.Duration) *UserEvent {
	return &UserEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingUserEventData{
			UserID:    userID,
			BlockedID: targetID,
			MutedID:   targetID,
			Duration:  duration,
		}),
	}
}

// RoleEvent represents a role-related event
type RoleEvent struct {
	*BaseMessagingEvent
}

// NewRoleEvent creates a new role event
func NewRoleEvent(eventType MessagingEventType, role *domain.MessagingRole, userID, roleID, actorID uint) *RoleEvent {
	return &RoleEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingRoleEventData{
			Role:        role,
			UserID:      userID,
			RoleID:      roleID,
			AddedByID:   actorID,
			RemovedByID: actorID,
		}),
	}
}

// ReactionEvent represents a reaction-related event
type ReactionEvent struct {
	*BaseMessagingEvent
}

// NewReactionEvent creates a new reaction event
func NewReactionEvent(eventType MessagingEventType, reaction *domain.MessagingReaction) *ReactionEvent {
	return &ReactionEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingReactionEventData{
			Reaction: reaction,
		}),
	}
}

// ReadReceiptEvent represents a read receipt event
type ReadReceiptEvent struct {
	*BaseMessagingEvent
}

// NewReadReceiptEvent creates a new read receipt event
func NewReadReceiptEvent(receipt *domain.MessagingReadReceipt) *ReadReceiptEvent {
	return &ReadReceiptEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(MessagingEventReadReceipt, &MessagingReadReceiptEventData{
			Receipt: receipt,
		}),
	}
}

// MemberEvent represents a member-related event
type MemberEvent struct {
	*BaseMessagingEvent
}

// NewMemberEvent creates a new member event
func NewMemberEvent(eventType MessagingEventType, channelID, userID, actorID uint) *MemberEvent {
	return &MemberEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingMemberEventData{
			ChannelID:   channelID,
			UserID:      userID,
			AddedByID:   actorID,
			RemovedByID: actorID,
		}),
	}
}

// AttachmentEvent represents an attachment-related event
type AttachmentEvent struct {
	*BaseMessagingEvent
}

// NewAttachmentEvent creates a new attachment event
func NewAttachmentEvent(eventType MessagingEventType, attachment *domain.MessagingAttachment) *AttachmentEvent {
	return &AttachmentEvent{
		BaseMessagingEvent: NewBaseMessagingEvent(eventType, &MessagingAttachmentEventData{
			Attachment: attachment,
		}),
	}
}

// MessagingAttachmentEventData represents the data for an attachment event
type MessagingAttachmentEventData struct {
	Attachment *domain.MessagingAttachment
}

// EventDispatcher defines the interface for event dispatching
type EventDispatcher interface {
	Dispatch(ctx context.Context, event MessagingEvent) error
}

// MessageCreatedEvent represents a message creation event
type MessageCreatedEvent struct {
	MessageID uint
	ChannelID uint
	UserID    uint
}

// Data returns the event data
func (e *MessageCreatedEvent) Data() interface{} {
	return e
}

// MessageUpdatedEvent represents a message update event
type MessageUpdatedEvent struct {
	MessageID uint
	ChannelID uint
	UserID    uint
}

// Data returns the event data
func (e *MessageUpdatedEvent) Data() interface{} {
	return e
}

// MessageDeletedEvent represents a message deletion event
type MessageDeletedEvent struct {
	MessageID uint
	ChannelID uint
	UserID    uint
}

// Data returns the event data
func (e *MessageDeletedEvent) Data() interface{} {
	return e
}

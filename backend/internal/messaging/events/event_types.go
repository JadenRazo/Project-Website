package events

import (
	"time"

	"Project-Website/backend/internal/messaging/domain"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
	"github.com/google/uuid"
)

// Event types
const (
	EventMessageCreated      = "message.created"
	EventMessageUpdated      = "message.updated"
	EventMessageDeleted      = "message.deleted"
	EventReactionAdded       = "reaction.added"
	EventReactionRemoved     = "reaction.removed"
	EventChannelCreated      = "channel.created"
	EventChannelUpdated      = "channel.updated"
	EventChannelDeleted      = "channel.deleted"
	EventUserJoinedChannel   = "channel.user_joined"
	EventUserLeftChannel     = "channel.user_left"
	EventUserPresenceChanged = "user.presence_changed"
	EventUserTypingStarted   = "user.typing_started"
	EventUserTypingStopped   = "user.typing_stopped"
)

// Event represents a system event
type Event struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
}

// MessageEventData contains data for message-related events
type MessageEventData struct {
	Message   *domain.Message `json:"message"`
	ChannelID uint            `json:"channelId"`
	SenderID  uint            `json:"senderId"`
}

// ChannelEventData contains data for channel-related events
type ChannelEventData struct {
	Channel *domain.Channel `json:"channel"`
	UserID  uint            `json:"userId,omitempty"` // User who triggered the event
}

// UserChannelEventData contains data for user-channel relation events
type UserChannelEventData struct {
	ChannelID uint `json:"channelId"`
	UserID    uint `json:"userId"`
}

// ReactionEventData contains data for reaction-related events
type ReactionEventData struct {
	Reaction  *domain.Reaction `json:"reaction"`
	MessageID uint             `json:"messageId"`
	UserID    uint             `json:"userId"`
}

// TypingEventData contains data for typing indicator events
type TypingEventData struct {
	ChannelID uint `json:"channelId"`
	UserID    uint `json:"userId"`
	IsTyping  bool `json:"isTyping"`
}

// PresenceEventData contains data for user presence events
type PresenceEventData struct {
	UserID    uint   `json:"userId"`
	Status    string `json:"status"`
	StatusMsg string `json:"statusMsg,omitempty"`
}

// EventType represents the type of message event
type EventType string

const (
	// MessageCreatedEvent is emitted when a new message is created
	MessageCreated EventType = "message.created"

	// MessageUpdatedEvent is emitted when a message is updated
	MessageUpdated EventType = "message.updated"

	// MessageDeletedEvent is emitted when a message is deleted
	MessageDeleted EventType = "message.deleted"

	// MessageDeliveredEvent is emitted when a message is delivered
	MessageDelivered EventType = "message.delivered"

	// MessageReadEvent is emitted when a message is read
	MessageRead EventType = "message.read"

	// MessageReactionAddedEvent is emitted when a reaction is added to a message
	MessageReactionAdded EventType = "message.reaction.added"

	// MessageReactionRemovedEvent is emitted when a reaction is removed from a message
	MessageReactionRemoved EventType = "message.reaction.removed"

	// MessagePinnedEvent is emitted when a message is pinned
	MessagePinned EventType = "message.pinned"

	// MessageUnpinnedEvent is emitted when a message is unpinned
	MessageUnpinned EventType = "message.unpinned"

	// MessageReportedEvent is emitted when a message is reported
	MessageReported EventType = "message.reported"

	// MessageModeratedEvent is emitted when a message is moderated
	MessageModerated EventType = "message.moderated"

	// MessageFilteredEvent is emitted when a message is filtered
	MessageFiltered EventType = "message.filtered"
)

// MessageEvent represents a message event
type MessageEvent struct {
	ID        uuid.UUID `json:"id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

// MessageCreatedData represents the data for a message created event
type MessageCreatedData struct {
	Message   *domain.Message `json:"message"`
	ChannelID uuid.UUID       `json:"channel_id"`
	UserID    uuid.UUID       `json:"user_id"`
}

// MessageUpdatedData represents the data for a message updated event
type MessageUpdatedData struct {
	Message   *domain.Message `json:"message"`
	ChannelID uuid.UUID       `json:"channel_id"`
	UserID    uuid.UUID       `json:"user_id"`
}

// MessageDeletedData represents the data for a message deleted event
type MessageDeletedData struct {
	MessageID uuid.UUID `json:"message_id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// MessageDeliveredData represents the data for a message delivered event
type MessageDeliveredData struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// MessageReadData represents the data for a message read event
type MessageReadData struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// MessageReactionData represents the data for a message reaction event
type MessageReactionData struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	Emoji     string    `json:"emoji"`
}

// MessagePinnedData represents the data for a message pinned event
type MessagePinnedData struct {
	MessageID uuid.UUID `json:"message_id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// MessageReportedData represents the data for a message reported event
type MessageReportedData struct {
	MessageID   uuid.UUID `json:"message_id"`
	ReporterID  uuid.UUID `json:"reporter_id"`
	Reason      string    `json:"reason"`
	Description string    `json:"description"`
}

// MessageModeratedData represents the data for a message moderated event
type MessageModeratedData struct {
	MessageID   uuid.UUID                   `json:"message_id"`
	ModeratorID uuid.UUID                   `json:"moderator_id"`
	Action      domain.ModerationActionType `json:"action"`
	Reason      string                      `json:"reason"`
	Duration    *time.Duration              `json:"duration,omitempty"`
}

// MessageFilteredData represents the data for a message filtered event
type MessageFilteredData struct {
	MessageID uuid.UUID              `json:"message_id"`
	FilterID  uuid.UUID              `json:"filter_id"`
	Scope     domain.WordFilterScope `json:"scope"`
	Word      string                 `json:"word"`
}

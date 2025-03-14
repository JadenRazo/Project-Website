package events

import (
	"time"

	"github.com/JadenRazo/Project-Website/backend/internal/domain"
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

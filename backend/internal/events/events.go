package events

import "context"

// EventDispatcher defines the interface for event dispatching
type EventDispatcher interface {
	Dispatch(ctx context.Context, event interface{}) error
}

// MessageCreatedEvent represents a message creation event
type MessageCreatedEvent struct {
	MessageID uint
	ChannelID uint
	UserID    uint
}

// MessageUpdatedEvent represents a message update event
type MessageUpdatedEvent struct {
	MessageID uint
	ChannelID uint
	UserID    uint
}

// MessageDeletedEvent represents a message deletion event
type MessageDeletedEvent struct {
	MessageID uint
	ChannelID uint
	UserID    uint
}

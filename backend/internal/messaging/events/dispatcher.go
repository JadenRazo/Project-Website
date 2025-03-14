package events

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// EventHandler defines a function signature for event handlers
type EventHandler func(ctx context.Context, event *Event) error

// EventDispatcher manages event subscriptions and dispatches events to subscribers
type EventDispatcher struct {
	// Map of event types to handlers
	handlers map[string][]EventHandler

	// Mutex for thread-safe handler management
	mu sync.RWMutex
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type
func (d *EventDispatcher) Subscribe(eventType string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[eventType] = append(d.handlers[eventType], handler)
	log.Printf("Handler subscribed to event type: %s", eventType)
}

// Unsubscribe removes a handler for a specific event type
func (d *EventDispatcher) Unsubscribe(eventType string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Find and remove the handler
	handlers, exists := d.handlers[eventType]
	if !exists {
		return
	}

	// This is somewhat inefficient as we can't directly compare functions
	// Instead, we're skipping the handler instance and rebuilding the slice
	var newHandlers []EventHandler
	for _, h := range handlers {
		// Skip the handler we want to remove
		// Note: In Go, you can't compare function values directly
		// So this is an approximation, in practice you'd need a way to identify handlers
		if fmt.Sprintf("%p", h) != fmt.Sprintf("%p", handler) {
			newHandlers = append(newHandlers, h)
		}
	}

	// Update handlers or remove event type if no handlers left
	if len(newHandlers) > 0 {
		d.handlers[eventType] = newHandlers
	} else {
		delete(d.handlers, eventType)
	}
}

// Dispatch sends an event to all subscribed handlers
func (d *EventDispatcher) Dispatch(ctx context.Context, event *Event) {
	d.mu.RLock()
	handlers, exists := d.handlers[event.Type]
	d.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		// No handlers for this event type
		return
	}

	// Process handlers in a separate goroutine to prevent blocking
	go func() {
		var wg sync.WaitGroup
		for _, handler := range handlers {
			wg.Add(1)
			// Execute each handler in its own goroutine
			go func(h EventHandler) {
				defer wg.Done()
				if err := h(ctx, event); err != nil {
					log.Printf("Error handling event %s: %v", event.Type, err)
				}
			}(handler)
		}
		wg.Wait()
	}()
}

// DispatchSync sends an event to all subscribed handlers and waits for completion
func (d *EventDispatcher) DispatchSync(ctx context.Context, event *Event) error {
	d.mu.RLock()
	handlers, exists := d.handlers[event.Type]
	d.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		// No handlers for this event type
		return nil
	}

	errs := make([]error, 0)
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errs = append(errs, fmt.Errorf("handler error for event %s: %w", event.Type, err))
		}
	}

	// If any handlers returned errors, aggregate them
	if len(errs) > 0 {
		return fmt.Errorf("dispatch errors: %v", errs)
	}

	return nil
}

// HasHandlers checks if there are handlers for a specific event type
func (d *EventDispatcher) HasHandlers(eventType string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers, exists := d.handlers[eventType]
	return exists && len(handlers) > 0
}

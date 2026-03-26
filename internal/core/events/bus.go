// Package events provides event bus functionality for inter-component communication.
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Event represents an event in the event bus.
type Event struct {
	ID        string
	Type      string
	Source    string
	Timestamp time.Time
	Data      interface{}
}

// Handler defines the interface for handling events.
type Handler interface {
	Handle(ctx context.Context, event *Event) error
}

// HandlerFunc is a function type that implements the Handler interface.
type HandlerFunc func(ctx context.Context, event *Event) error

// Handle calls the underlying function.
func (f HandlerFunc) Handle(ctx context.Context, event *Event) error {
	return f(ctx, event)
}

// Bus defines the interface for an event bus.
type Bus interface {
	Publish(ctx context.Context, event *Event) error
	Subscribe(eventType string, handler Handler) (string, error)
	Unsubscribe(subscriptionID string) error
}

// InMemoryBus is an in-memory implementation of the event bus.
type InMemoryBus struct {
	mu          sync.RWMutex
	subscribers map[string]map[string]Handler
	nextID      int
}

// NewInMemoryBus creates a new in-memory event bus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		subscribers: make(map[string]map[string]Handler),
	}
}

// Publish sends an event to all registered subscribers.
func (b *InMemoryBus) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if event.ID == "" {
		event.ID = fmt.Sprintf("%s-%d", event.Type, time.Now().UnixNano())
	}

	b.mu.RLock()
	src := b.subscribers[event.Type]
	handlers := make(map[string]Handler, len(src))
	for id, h := range src {
		handlers[id] = h
	}
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for subID, handler := range handlers {
		wg.Add(1)
		eventCopy := deepCopyEvent(event)
		go func(id string, h Handler, evt *Event) {
			defer wg.Done()
			if err := h.Handle(ctx, evt); err != nil {
				errChan <- fmt.Errorf("handler %s failed: %w", id, err)
			}
		}(subID, handler, eventCopy)
	}

	wg.Wait()
	close(errChan)

	var allErrors []error
	for err := range errChan {
		if err != nil {
			allErrors = append(allErrors, err)
		}
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("%d event handlers failed: %v", len(allErrors), allErrors[0])
	}

	return nil
}

func deepCopyEvent(event *Event) *Event {
	if event == nil {
		return nil
	}
	copy := &Event{
		ID:        event.ID,
		Type:      event.Type,
		Source:    event.Source,
		Timestamp: event.Timestamp,
	}
	if event.Data != nil {
		if data, err := json.Marshal(event.Data); err == nil {
			var newData interface{}
			if err := json.Unmarshal(data, &newData); err == nil {
				copy.Data = newData
			} else {
				copy.Data = event.Data
			}
		} else {
			copy.Data = event.Data
		}
	}
	return copy
}

// Subscribe registers a handler for the given event type.
func (b *InMemoryBus) Subscribe(eventType string, handler Handler) (string, error) {
	if eventType == "" {
		return "", fmt.Errorf("event type cannot be empty")
	}

	if handler == nil {
		return "", fmt.Errorf("handler cannot be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subscribers[eventType] == nil {
		b.subscribers[eventType] = make(map[string]Handler)
	}

	b.nextID++
	subID := fmt.Sprintf("sub-%d", b.nextID)
	b.subscribers[eventType][subID] = handler

	return subID, nil
}

// Unsubscribe removes a subscription by its ID.
func (b *InMemoryBus) Unsubscribe(subscriptionID string) error {
	if subscriptionID == "" {
		return fmt.Errorf("subscription ID cannot be empty")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	for eventType, handlers := range b.subscribers {
		if _, exists := handlers[subscriptionID]; exists {
			delete(handlers, subscriptionID)

			if len(handlers) == 0 {
				delete(b.subscribers, eventType)
			}

			return nil
		}
	}

	return fmt.Errorf("subscription %s not found", subscriptionID)
}

// GetSubscriberCount returns the number of subscribers for an event type.
func (b *InMemoryBus) GetSubscriberCount(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, exists := b.subscribers[eventType]; exists {
		return len(handlers)
	}
	return 0
}

// ClearAllSubscribers removes all subscribers from the bus.
func (b *InMemoryBus) ClearAllSubscribers() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = make(map[string]map[string]Handler)
}

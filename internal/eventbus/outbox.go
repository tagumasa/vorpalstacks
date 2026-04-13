package eventbus

import (
	"context"
	"time"
)

// OutboxStatus represents the lifecycle state of an outbox entry.
type OutboxStatus int

const (
	// OutboxPending indicates the entry has been written but not yet dispatched.
	OutboxPending OutboxStatus = iota
	// OutboxProcessing indicates the entry is currently being dispatched.
	OutboxProcessing
	// OutboxDelivered indicates all handlers have processed the entry successfully.
	OutboxDelivered
	// OutboxFailed indicates the entry exhausted all retry attempts without success.
	OutboxFailed
)

// String returns a human-readable label for the outbox status.
func (s OutboxStatus) String() string {
	switch s {
	case OutboxPending:
		return "PENDING"
	case OutboxProcessing:
		return "PROCESSING"
	case OutboxDelivered:
		return "DELIVERED"
	case OutboxFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

// OutboxEntry represents a single event persisted in the outbox store,
// tracking its serialisation, status, retry state, and per-handler results.
type OutboxEntry struct {
	EventID         string            `json:"event_id"`
	EventType       string            `json:"event_type"`
	Depth           int               `json:"depth"`
	SerializedEvent []byte            `json:"serialized_event"`
	Status          OutboxStatus      `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	DeliveredAt     *time.Time        `json:"delivered_at,omitempty"`
	RetryCount      int32             `json:"retry_count"`
	MaxRetries      int32             `json:"max_retries"`
	LastError       string            `json:"last_error,omitempty"`
	HandlerResults  map[string]string `json:"handler_results,omitempty"`
}

// OutboxStore defines the persistence contract for the event outbox,
// enabling durable async delivery with at-least-once semantics.
type OutboxStore interface {
	Write(ctx context.Context, entry *OutboxEntry) error
	Read(ctx context.Context, eventID string) (*OutboxEntry, error)
	UpdateStatus(ctx context.Context, eventID string, from, to OutboxStatus) (bool, error)
	UpdateEntry(ctx context.Context, entry *OutboxEntry) error
	ListPending(ctx context.Context, limit int) ([]*OutboxEntry, error)
	Delete(ctx context.Context, eventID string) error
	Cleanup(ctx context.Context, deliveredBefore time.Time, failedBefore time.Time) (int, error)
	Close() error
}

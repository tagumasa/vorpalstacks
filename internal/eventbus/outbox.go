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
	// ListPendingFrom returns up to 'limit' pending entries together with
	// the opaque cursor of the last returned entry. Passing that cursor
	// back resumes strictly after it; an empty cursor starts from the
	// head, making the call equivalent to ListPending. Callers paginate
	// through backlogs larger than a single page by feeding the returned
	// cursor into the next call.
	ListPendingFrom(ctx context.Context, limit int, afterCursor string) ([]*OutboxEntry, string, error)
	// ResetStaleProcessing finds entries stuck in OutboxProcessing
	// (left behind by a crash or panic) and resets them to OutboxPending
	// so they can be recovered on the next Start. Returns the number of
	// entries reset.
	ResetStaleProcessing(ctx context.Context) (int, error)
	Delete(ctx context.Context, eventID string) error
	// Cleanup purges delivered entries older than deliveredBefore and
	// failed entries older than failedBefore. Pending entries are never
	// purged: an undelivered event must survive until it delivers or
	// exhausts its retry budget and transitions to Failed.
	Cleanup(ctx context.Context, deliveredBefore time.Time, failedBefore time.Time) (int, error)
	Close() error
}

package eventbus

import (
	"context"
	"time"
)

type OutboxStatus int

const (
	OutboxPending OutboxStatus = iota
	OutboxProcessing
	OutboxDelivered
	OutboxFailed
)

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

// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/logs"
)

// BulkheadConfig defines the configuration for a bulkhead.
type BulkheadConfig struct {
	Name          string
	MaxConcurrent int
	MaxWait       time.Duration
	Logger        logs.Logger
	Metrics       *Metrics
}

// DefaultBulkheadConfig returns the default bulkhead configuration.
func DefaultBulkheadConfig() *BulkheadConfig {
	return &BulkheadConfig{
		MaxConcurrent: 10,
		MaxWait:       5 * time.Second,
	}
}

// Bulkhead implements the bulkhead pattern to limit concurrent executions.
type Bulkhead struct {
	config      *BulkheadConfig
	semaphore   chan struct{}
	activeCount int32
}

// NewBulkhead creates a new bulkhead with the given configuration.
func NewBulkhead(config *BulkheadConfig) *Bulkhead {
	if config == nil {
		config = DefaultBulkheadConfig()
	}

	if config.MaxWait == 0 {
		config.MaxWait = DefaultBulkheadConfig().MaxWait
	}

	return &Bulkhead{
		config:    config,
		semaphore: make(chan struct{}, config.MaxConcurrent),
	}
}

// Execute runs the given function with bulkhead protection, waiting for a slot if necessary.
func (b *Bulkhead) Execute(ctx context.Context, fn func() error) error {
	startTime := time.Now()

	if ctx.Err() != nil {
		if b.config.Logger != nil {
			b.config.Logger.Debug("Bulkhead operation cancelled by context",
				logs.String("name", b.config.Name),
			)
		}
		return fmt.Errorf("bulkhead %s: %w", b.config.Name, ctx.Err())
	}

	waitTimer := time.NewTimer(b.config.MaxWait)
	defer waitTimer.Stop()

	select {
	case b.semaphore <- struct{}{}:
		atomic.AddInt32(&b.activeCount, 1)

		defer func() {
			<-b.semaphore
			atomic.AddInt32(&b.activeCount, -1)

			if b.config.Metrics != nil {
				b.config.Metrics.RecordLatency(b.config.Name, time.Since(startTime))
			}
		}()

		err := fn()
		if err != nil {
			if b.config.Metrics != nil {
				b.config.Metrics.RecordError("bulkhead_operation_failed")
			}
			return err
		}

		return nil

	case <-waitTimer.C:
		if b.config.Logger != nil {
			b.config.Logger.Warn("Bulkhead max wait time exceeded",
				logs.String("name", b.config.Name),
				logs.Int("maxConcurrent", b.config.MaxConcurrent),
				logs.Int("activeCount", b.ActiveCount()),
			)
		}

		if b.config.Metrics != nil {
			b.config.Metrics.RecordError("bulkhead_max_wait_exceeded")
		}

		return NewRateLimit(fmt.Sprintf("bulkhead %s: max wait time exceeded", b.config.Name))

	case <-ctx.Done():
		if b.config.Logger != nil {
			b.config.Logger.Debug("Bulkhead operation cancelled by context",
				logs.String("name", b.config.Name),
			)
		}
		return fmt.Errorf("bulkhead %s: %w", b.config.Name, ctx.Err())
	}
}

// TryExecute runs the given function with bulkhead protection without waiting for a slot.
func (b *Bulkhead) TryExecute(fn func() error) error {
	startTime := time.Now()

	select {
	case b.semaphore <- struct{}{}:
		atomic.AddInt32(&b.activeCount, 1)

		defer func() {
			<-b.semaphore
			atomic.AddInt32(&b.activeCount, -1)

			if b.config.Metrics != nil {
				b.config.Metrics.RecordLatency(b.config.Name, time.Since(startTime))
			}
		}()

		err := fn()
		if err != nil {
			if b.config.Metrics != nil {
				b.config.Metrics.RecordError("bulkhead_operation_failed")
			}
			return err
		}

		return nil

	default:
		if b.config.Logger != nil {
			b.config.Logger.Debug("Bulkhead is full",
				logs.String("name", b.config.Name),
				logs.Int("maxConcurrent", b.config.MaxConcurrent),
				logs.Int("activeCount", b.ActiveCount()),
			)
		}

		if b.config.Metrics != nil {
			b.config.Metrics.RecordError("bulkhead_full")
		}

		return NewRateLimit(fmt.Sprintf("bulkhead %s: no slot available", b.config.Name))
	}
}

// ActiveCount returns the number of currently active executions.
func (b *Bulkhead) ActiveCount() int {
	return int(atomic.LoadInt32(&b.activeCount))
}

// AvailableSlots returns the number of available slots for new executions.
func (b *Bulkhead) AvailableSlots() int {
	return b.config.MaxConcurrent - b.ActiveCount()
}

// MaxConcurrent returns the maximum number of concurrent executions allowed.
func (b *Bulkhead) MaxConcurrent() int {
	return b.config.MaxConcurrent
}

// Stats returns the current statistics of the bulkhead.
func (b *Bulkhead) Stats() BulkheadStats {
	return BulkheadStats{
		Name:           b.config.Name,
		MaxConcurrent:  b.config.MaxConcurrent,
		ActiveCount:    b.ActiveCount(),
		AvailableSlots: b.AvailableSlots(),
	}
}

// BulkheadStats represents statistics about the bulkhead.
type BulkheadStats struct {
	Name           string
	MaxConcurrent  int
	ActiveCount    int
	AvailableSlots int
}

// String returns the string representation of the BulkheadStats.
func (s BulkheadStats) String() string {
	return fmt.Sprintf("Name: %s, MaxConcurrent: %d, ActiveCount: %d, AvailableSlots: %d",
		s.Name, s.MaxConcurrent, s.ActiveCount, s.AvailableSlots)
}

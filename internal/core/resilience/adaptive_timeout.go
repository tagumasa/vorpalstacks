// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
)

// AdaptiveTimeoutConfig holds configuration for an adaptive timeout mechanism.
// It defines the initial, minimum, and maximum timeout values, as well as
// thresholds for adjusting the timeout based on success or failure rates.
type AdaptiveTimeoutConfig struct {
	Name             string
	InitialTimeout   time.Duration
	MinTimeout       time.Duration
	MaxTimeout       time.Duration
	SuccessThreshold int
	FailureThreshold int
	AdjustmentFactor float64
	Logger           logs.Logger
	Metrics          *Metrics
}

// DefaultAdaptiveTimeoutConfig returns the default adaptive timeout configuration.
func DefaultAdaptiveTimeoutConfig() *AdaptiveTimeoutConfig {
	return &AdaptiveTimeoutConfig{
		InitialTimeout:   30 * time.Second,
		MinTimeout:       5 * time.Second,
		MaxTimeout:       60 * time.Second,
		SuccessThreshold: 5,
		FailureThreshold: 2,
		AdjustmentFactor: 1.5,
	}
}

// AdaptiveTimeout dynamically adjusts the operation timeout based on
// success and failure rates. It increases the timeout on repeated failures
// and decreases it when operations succeed consistently.
type AdaptiveTimeout struct {
	config               *AdaptiveTimeoutConfig
	currentTimeout       time.Duration
	consecutiveSuccesses int
	consecutiveFailures  int
	mu                   sync.RWMutex
	done                 chan struct{}
	closeOnce            sync.Once
}

// NewAdaptiveTimeout creates a new adaptive timeout with the given configuration.
func NewAdaptiveTimeout(config *AdaptiveTimeoutConfig) *AdaptiveTimeout {
	if config == nil {
		config = DefaultAdaptiveTimeoutConfig()
	}

	return &AdaptiveTimeout{
		config:         config,
		currentTimeout: config.InitialTimeout,
		done:           make(chan struct{}),
	}
}

// Execute runs the given function with an adaptive timeout.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute
//
// Returns:
//   - error: Any error returned by the function or timeout
func (at *AdaptiveTimeout) Execute(ctx context.Context, fn func(context.Context) error) error {
	timeout := at.GetTimeout()

	if at.config.Logger != nil {
		at.config.Logger.Debug("AdaptiveTimeout Execute starting",
			logs.String("name", at.config.Name),
			logs.Any("timeout", timeout),
		)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultChan <- fmt.Errorf("panic: %v", r)
			}
			close(resultChan)
		}()
		select {
		case resultChan <- fn(timeoutCtx):
		case <-at.done:
		case <-timeoutCtx.Done():
		}
	}()

	select {
	case err, ok := <-resultChan:
		if !ok {
			if at.config.Logger != nil {
				at.config.Logger.Error("adaptive timeout: result channel closed without result", logs.String("name", at.config.Name))
			}
			return NewTimeout("operation timed out")
		}
		at.recordResult(err)
		return err
	case <-timeoutCtx.Done():
		if at.config.Logger != nil {
			at.config.Logger.Warn("adaptive timeout exceeded",
				logs.String("name", at.config.Name),
				logs.Any("timeout", timeout),
			)
		}
		select {
		case err, ok := <-resultChan:
			if ok {
				at.recordResult(err)
			}
		default:
			at.recordResult(timeoutCtx.Err())
		}
		return NewTimeout("operation timed out")
	case <-ctx.Done():
		select {
		case err, ok := <-resultChan:
			if !ok {
				return ctx.Err()
			}
			return err
		default:
			return ctx.Err()
		}
	}
}

// GetTimeout returns the current adaptive timeout duration.
//
// Returns:
//   - time.Duration: The current timeout value
func (at *AdaptiveTimeout) GetTimeout() time.Duration {
	at.mu.RLock()
	defer at.mu.RUnlock()
	return at.currentTimeout
}

// recordResult records the result of an operation for adaptive timeout adjustment.
//
// Parameters:
//   - err: The error result of the operation (nil if successful)
func (at *AdaptiveTimeout) recordResult(err error) {
	at.mu.Lock()
	defer at.mu.Unlock()

	if err == nil {
		at.consecutiveSuccesses++
		at.consecutiveFailures = 0

		if at.consecutiveSuccesses >= at.config.SuccessThreshold {
			at.decreaseTimeoutLocked()
			at.consecutiveSuccesses = 0
		}
	} else {
		at.consecutiveFailures++
		at.consecutiveSuccesses = 0

		if at.consecutiveFailures >= at.config.FailureThreshold {
			at.increaseTimeoutLocked()
			at.consecutiveFailures = 0
		}

		if at.config.Metrics != nil {
			at.config.Metrics.RecordError("adaptive_timeout_operation_failed")
		}
	}
}

func (at *AdaptiveTimeout) increaseTimeoutLocked() {
	oldTimeout := at.currentTimeout
	newTimeout := time.Duration(float64(at.currentTimeout) * at.config.AdjustmentFactor)

	if newTimeout > at.config.MaxTimeout {
		newTimeout = at.config.MaxTimeout
	}

	at.currentTimeout = newTimeout

	if at.config.Logger != nil {
		at.config.Logger.Info("Adaptive timeout increased",
			logs.String("name", at.config.Name),
			logs.Any("oldTimeout", oldTimeout),
			logs.Any("newTimeout", newTimeout),
		)
	}

	if at.config.Metrics != nil {
		at.config.Metrics.RecordLatency(at.config.Name+"_timeout_increase", newTimeout-oldTimeout)
	}
}

func (at *AdaptiveTimeout) decreaseTimeoutLocked() {
	oldTimeout := at.currentTimeout
	newTimeout := time.Duration(float64(at.currentTimeout) / at.config.AdjustmentFactor)

	if newTimeout < at.config.MinTimeout {
		newTimeout = at.config.MinTimeout
	}

	at.currentTimeout = newTimeout

	if at.config.Logger != nil {
		at.config.Logger.Info("Adaptive timeout decreased",
			logs.String("name", at.config.Name),
			logs.Any("oldTimeout", oldTimeout),
			logs.Any("newTimeout", newTimeout),
		)
	}

	if at.config.Metrics != nil {
		at.config.Metrics.RecordLatency(at.config.Name+"_timeout_decrease", oldTimeout-newTimeout)
	}
}

// Reset resets the adaptive timeout to its initial value.
func (at *AdaptiveTimeout) Reset() {
	at.mu.Lock()
	defer at.mu.Unlock()

	at.currentTimeout = at.config.InitialTimeout
	at.consecutiveSuccesses = 0
	at.consecutiveFailures = 0

	if at.config.Logger != nil {
		at.config.Logger.Info("Adaptive timeout reset",
			logs.String("name", at.config.Name),
			logs.Any("timeout", at.currentTimeout),
		)
	}
}

// Close closes the adaptive timeout mechanism.
func (at *AdaptiveTimeout) Close() error {
	at.closeOnce.Do(func() {
		close(at.done)
	})
	return nil
}

// Stats returns the current statistics of the adaptive timeout.
//
// Returns:
//   - AdaptiveTimeoutStats: The current statistics
func (at *AdaptiveTimeout) Stats() AdaptiveTimeoutStats {
	at.mu.RLock()
	defer at.mu.RUnlock()

	return AdaptiveTimeoutStats{
		Name:                 at.config.Name,
		CurrentTimeout:       at.currentTimeout,
		ConsecutiveSuccesses: at.consecutiveSuccesses,
		ConsecutiveFailures:  at.consecutiveFailures,
	}
}

// AdaptiveTimeoutStats holds the current statistics for an adaptive timeout,
// including the current timeout value and counts of consecutive successes and failures.
type AdaptiveTimeoutStats struct {
	Name                 string
	CurrentTimeout       time.Duration
	ConsecutiveSuccesses int
	ConsecutiveFailures  int
}

// String returns a string representation of the adaptive timeout statistics.
func (s AdaptiveTimeoutStats) String() string {
	return fmt.Sprintf("%s: Timeout=%s, Successes=%d, Failures=%d",
		s.Name, s.CurrentTimeout, s.ConsecutiveSuccesses, s.ConsecutiveFailures)
}

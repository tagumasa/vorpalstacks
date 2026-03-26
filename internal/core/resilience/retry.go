// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"vorpalstacks/internal/core/logs"
)

func randomInt63n(n int64) int64 {
	if n <= 0 {
		return 0
	}
	randomNum, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		return 0
	}
	return randomNum.Int64()
}

// RetryPolicy defines the configuration and behaviour for retry operations.
type RetryPolicy struct {
	name         string
	maxAttempts  int
	maxDelay     time.Duration
	initialDelay time.Duration
	backoff      BackoffStrategy
	jitter       JitterStrategy
	logger       logs.Logger
	metrics      *Metrics
}

// BackoffStrategy defines the interface for calculating delay between retries.
type BackoffStrategy interface {
	Delay(attempt int, initialDelay time.Duration) time.Duration
}

// JitterStrategy defines the interface for applying jitter to delay calculations.
type JitterStrategy interface {
	Apply(delay time.Duration) time.Duration
}

// ExponentialBackoff implements an exponential backoff strategy.
type ExponentialBackoff struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	jitter       JitterStrategy
}

// Delay calculates the delay for the given attempt number using exponential backoff.
func (eb *ExponentialBackoff) Delay(attempt int, initialDelay time.Duration) time.Duration {
	delay := initialDelay
	for i := 0; i < attempt; i++ {
		delay = time.Duration(float64(delay) * eb.Multiplier)
	}
	if delay > eb.MaxDelay {
		delay = eb.MaxDelay
	}
	if eb.jitter != nil {
		delay = eb.jitter.Apply(delay)
	}
	return delay
}

// LinearBackoff implements a linear backoff strategy.
type LinearBackoff struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Increment    time.Duration
	jitter       JitterStrategy
}

// Delay calculates the delay for the given attempt number using linear backoff.
func (lb *LinearBackoff) Delay(attempt int, initialDelay time.Duration) time.Duration {
	delay := initialDelay + time.Duration(attempt)*lb.Increment
	if delay > lb.MaxDelay {
		delay = lb.MaxDelay
	}
	if lb.jitter != nil {
		delay = lb.jitter.Apply(delay)
	}
	return delay
}

// FullJitter applies full jitter to the delay.
type FullJitter struct{}

// Apply adds full random jitter to the delay.
func (fj *FullJitter) Apply(delay time.Duration) time.Duration {
	if delay <= 0 {
		return delay
	}
	return time.Duration(randomInt63n(int64(delay)))
}

// EqualJitter applies equal jitter to the delay.
type EqualJitter struct{}

// Apply adds equal jitter to the delay.
func (ej *EqualJitter) Apply(delay time.Duration) time.Duration {
	if delay <= 0 {
		return delay
	}
	half := delay / 2
	jitter := time.Duration(randomInt63n(int64(half)))
	return delay - half + jitter
}

// NoJitter applies no jitter to the delay.
type NoJitter struct{}

// Apply returns the delay unchanged (no jitter).
func (nj *NoJitter) Apply(delay time.Duration) time.Duration {
	return delay
}

// NewRetryPolicy creates a new retry policy with default configuration.
func NewRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		name:         "default",
		maxAttempts:  DefaultRetryMaxAttempts,
		maxDelay:     DefaultRetryMaxDelay,
		initialDelay: DefaultRetryInitialDelay,
		backoff: &ExponentialBackoff{
			InitialDelay: DefaultRetryInitialDelay,
			MaxDelay:     DefaultRetryMaxDelay,
			Multiplier:   DefaultRetryBackoffMultiplier,
			jitter:       &EqualJitter{},
		},
	}
}

// NewRetryPolicyWithBackoff creates a new retry policy with the given backoff strategy.
func NewRetryPolicyWithBackoff(backoff BackoffStrategy) *RetryPolicy {
	return &RetryPolicy{
		name:         "default",
		maxAttempts:  DefaultRetryMaxAttempts,
		maxDelay:     DefaultRetryMaxDelay,
		initialDelay: DefaultRetryInitialDelay,
		backoff:      backoff,
	}
}

// SetMaxAttempts sets the maximum number of retry attempts.
func (rp *RetryPolicy) SetMaxAttempts(maxAttempts int) {
	rp.maxAttempts = maxAttempts
}

// SetMaxDelay sets the maximum delay between retries.
func (rp *RetryPolicy) SetMaxDelay(maxDelay time.Duration) {
	rp.maxDelay = maxDelay
}

// SetInitialDelay sets the initial delay between retries.
func (rp *RetryPolicy) SetInitialDelay(initialDelay time.Duration) {
	rp.initialDelay = initialDelay
}

// SetBackoff sets the backoff strategy for calculating delays.
func (rp *RetryPolicy) SetBackoff(backoff BackoffStrategy) {
	rp.backoff = backoff
}

// SetJitter sets the jitter strategy for adding randomness to delays.
func (rp *RetryPolicy) SetJitter(jitter JitterStrategy) {
	rp.jitter = jitter
	if eb, ok := rp.backoff.(*ExponentialBackoff); ok {
		eb.jitter = jitter
	}
	if lb, ok := rp.backoff.(*LinearBackoff); ok {
		lb.jitter = jitter
	}
}

// SetLogger sets the logger for the retry policy.
func (rp *RetryPolicy) SetLogger(logger logs.Logger) {
	rp.logger = logger
}

// SetName sets the name for the retry policy.
func (rp *RetryPolicy) SetName(name string) {
	rp.name = name
}

// SetMetrics sets the metrics collector for the retry policy.
func (rp *RetryPolicy) SetMetrics(metrics *Metrics) {
	rp.metrics = metrics
}

// Do executes the given function with retry logic.
func (rp *RetryPolicy) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error
	startTime := time.Now()

	for attempt := 0; attempt < rp.maxAttempts; attempt++ {
		if attempt > 0 {
			delay := rp.backoff.Delay(attempt, rp.initialDelay)
			if delay > rp.maxDelay {
				delay = rp.maxDelay
			}

			if rp.metrics != nil {
				rp.metrics.RecordRetry(rp.name)
			}

			if rp.logger != nil {
				rp.logger.Debug("Retrying operation",
					logs.Int("attempt", attempt),
					logs.Any("delay", delay),
					logs.Err(lastErr),
				)
			}

			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			}
		}

		err := fn(ctx)
		if err == nil {
			if rp.logger != nil && attempt > 0 {
				rp.logger.Info("Operation succeeded after retry",
					logs.Int("attempt", attempt+1),
				)
			}
			if rp.metrics != nil {
				rp.metrics.RecordLatency(rp.name, time.Since(startTime))
			}
			return nil
		}

		lastErr = err

		if !rp.isRetryable(err) {
			if rp.logger != nil {
				rp.logger.Warn("Operation failed with non-retryable error",
					logs.Int("attempt", attempt+1),
					logs.Err(err),
				)
			}
			if rp.metrics != nil {
				rp.metrics.RecordError("retry_non_retryable")
			}
			return err
		}
	}

	if rp.logger != nil {
		rp.logger.Error("Max retry attempts reached",
			logs.Int("maxAttempts", rp.maxAttempts),
			logs.Err(lastErr),
		)
	}

	if rp.metrics != nil {
		rp.metrics.RecordError("retry_max_attempts_reached")
	}

	return fmt.Errorf("max retry attempts (%d) reached: %w", rp.maxAttempts, lastErr)
}

func (rp *RetryPolicy) isRetryable(err error) bool {
	if err == nil {
		return false
	}
	return IsRetryable(err)
}

// GetMaxAttempts returns the maximum number of retry attempts.
func (rp *RetryPolicy) GetMaxAttempts() int {
	return rp.maxAttempts
}

// GetMaxDelay returns the maximum delay between retries.
func (rp *RetryPolicy) GetMaxDelay() time.Duration {
	return rp.maxDelay
}

// GetInitialDelay returns the initial delay between retries.
func (rp *RetryPolicy) GetInitialDelay() time.Duration {
	return rp.initialDelay
}

// GetBackoff returns the backoff strategy.
func (rp *RetryPolicy) GetBackoff() BackoffStrategy {
	return rp.backoff
}

// GetJitter returns the jitter strategy.
func (rp *RetryPolicy) GetJitter() JitterStrategy {
	return rp.jitter
}

// GetName returns the name of the retry policy.
func (rp *RetryPolicy) GetName() string {
	return rp.name
}

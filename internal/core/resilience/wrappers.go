// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"context"
	"fmt"
)

// CircuitBreakerWrapper provides a convenient wrapper around the CircuitBreaker
// for executing operations with circuit breaker protection.
type CircuitBreakerWrapper struct {
	circuitBreaker *CircuitBreaker
	name           string
}

// NewCircuitBreakerWrapper creates a new circuit breaker wrapper with the given name and circuit breaker.
func NewCircuitBreakerWrapper(name string, cb *CircuitBreaker) *CircuitBreakerWrapper {
	return &CircuitBreakerWrapper{
		circuitBreaker: cb,
		name:           name,
	}
}

// Execute runs the provided function with circuit breaker protection.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute
//
// Returns:
//   - error: An error if the circuit breaker is open or the function fails
func (w *CircuitBreakerWrapper) Execute(ctx context.Context, fn func() error) error {
	err := w.circuitBreaker.Execute(ctx, fn)
	if err != nil {
		if IsCircuitBreakerOpen(err) {
			return NewCircuitBreakerOpen(w.name)
		}
		return err
	}
	return nil
}

// ExecuteWithResult runs the provided function with circuit breaker protection and returns a result.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute that returns a result and error
//
// Returns:
//   - interface{}: The result from the function
//   - error: An error if the circuit breaker is open or the function fails
func (w *CircuitBreakerWrapper) ExecuteWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	var result interface{}

	err := w.circuitBreaker.Execute(ctx, func() error {
		var err error
		result, err = fn()
		return err
	})

	if err != nil {
		if IsCircuitBreakerOpen(err) {
			return nil, NewCircuitBreakerOpen(w.name)
		}
		return nil, err
	}

	return result, nil
}

// GetState returns the current state of the circuit breaker.
//
// Returns:
//   - State: The current circuit breaker state
func (w *CircuitBreakerWrapper) GetState() State {
	return w.circuitBreaker.State()
}

// GetStats returns the statistics for the circuit breaker.
//
// Returns:
//   - Stats: The circuit breaker statistics
func (w *CircuitBreakerWrapper) GetStats() Stats {
	return w.circuitBreaker.Stats()
}

// Reset resets the circuit breaker to closed state.
func (w *CircuitBreakerWrapper) Reset() {
	w.circuitBreaker.Reset()
}

// BulkheadWrapper provides a convenient wrapper around the Bulkhead
// for limiting concurrent execution of operations.
type BulkheadWrapper struct {
	bulkhead *Bulkhead
	name     string
}

// NewBulkheadWrapper creates a new bulkhead wrapper with the given name and bulkhead.
func NewBulkheadWrapper(name string, bh *Bulkhead) *BulkheadWrapper {
	return &BulkheadWrapper{
		bulkhead: bh,
		name:     name,
	}
}

// Execute runs the provided function with bulkhead (concurrency limit) protection.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute
//
// Returns:
//   - error: An error if the bulkhead is full or the function fails
func (w *BulkheadWrapper) Execute(ctx context.Context, fn func() error) error {
	err := w.bulkhead.Execute(ctx, fn)
	if err != nil {
		if IsRateLimit(err) {
			return NewRateLimit(fmt.Sprintf("bulkhead %s is full", w.name))
		}
		return err
	}
	return nil
}

// ExecuteWithResult runs the provided function with bulkhead protection and returns a result.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute that returns a result and error
//
// Returns:
//   - interface{}: The result from the function
//   - error: An error if the bulkhead is full or the function fails
func (w *BulkheadWrapper) ExecuteWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	var result interface{}

	err := w.bulkhead.Execute(ctx, func() error {
		var err error
		result, err = fn()
		return err
	})

	if err != nil {
		if IsRateLimit(err) {
			return nil, NewRateLimit(fmt.Sprintf("bulkhead %s is full", w.name))
		}
		return nil, err
	}

	return result, nil
}

// GetActiveCount returns the number of currently active executions.
//
// Returns:
//   - int: The number of active executions
func (w *BulkheadWrapper) GetActiveCount() int {
	return w.bulkhead.ActiveCount()
}

// GetAvailableSlots returns the number of available slots for new executions.
//
// Returns:
//   - int: The number of available slots
func (w *BulkheadWrapper) GetAvailableSlots() int {
	return w.bulkhead.AvailableSlots()
}

// GetStats returns the statistics for the bulkhead.
//
// Returns:
//   - BulkheadStats: The bulkhead statistics
func (w *BulkheadWrapper) GetStats() BulkheadStats {
	return w.bulkhead.Stats()
}

// AdaptiveTimeoutWrapper provides a convenient wrapper around the AdaptiveTimeout
// for executing operations with dynamically adjusted timeouts.
type AdaptiveTimeoutWrapper struct {
	adaptiveTimeout *AdaptiveTimeout
	name            string
}

// NewAdaptiveTimeoutWrapper creates a new adaptive timeout wrapper with the given name and adaptive timeout.
func NewAdaptiveTimeoutWrapper(name string, at *AdaptiveTimeout) *AdaptiveTimeoutWrapper {
	return &AdaptiveTimeoutWrapper{
		adaptiveTimeout: at,
		name:            name,
	}
}

// Execute runs the provided function with adaptive timeout protection.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute
//
// Returns:
//   - error: An error if the timeout is exceeded or the function fails
func (w *AdaptiveTimeoutWrapper) Execute(ctx context.Context, fn func(context.Context) error) error {
	err := w.adaptiveTimeout.Execute(ctx, fn)
	if err != nil {
		if IsTimeout(err) {
			return NewTimeout(fmt.Sprintf("adaptive timeout %s exceeded", w.name))
		}
		return err
	}
	return nil
}

// ExecuteWithResult runs the provided function with adaptive timeout protection and returns a result.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute that returns a result and error
//
// Returns:
//   - interface{}: The result from the function
//   - error: An error if the timeout is exceeded or the function fails
func (w *AdaptiveTimeoutWrapper) ExecuteWithResult(ctx context.Context, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	var result interface{}

	err := w.adaptiveTimeout.Execute(ctx, func(ctx context.Context) error {
		var err error
		result, err = fn(ctx)
		return err
	})

	if err != nil {
		if IsTimeout(err) {
			return nil, NewTimeout(fmt.Sprintf("adaptive timeout %s exceeded", w.name))
		}
		return nil, err
	}

	return result, nil
}

// GetTimeout returns the current adaptive timeout value in milliseconds.
//
// Returns:
//   - int64: The timeout in milliseconds
func (w *AdaptiveTimeoutWrapper) GetTimeout() int64 {
	return w.adaptiveTimeout.GetTimeout().Milliseconds()
}

// GetStats returns the statistics for the adaptive timeout.
//
// Returns:
//   - AdaptiveTimeoutStats: The adaptive timeout statistics
func (w *AdaptiveTimeoutWrapper) GetStats() AdaptiveTimeoutStats {
	return w.adaptiveTimeout.Stats()
}

// Reset resets the adaptive timeout to its initial value.
func (w *AdaptiveTimeoutWrapper) Reset() {
	w.adaptiveTimeout.Reset()
}

// ResilienceWrapper combines multiple resilience patterns (circuit breaker, bulkhead,
// adaptive timeout, and retry) into a single wrapper for convenient usage.
type ResilienceWrapper struct {
	circuitBreakerWrapper  *CircuitBreakerWrapper
	bulkheadWrapper        *BulkheadWrapper
	adaptiveTimeoutWrapper *AdaptiveTimeoutWrapper
	retryPolicy            *RetryPolicy
	name                   string
}

// NewResilienceWrapper creates a new resilience wrapper with the given name and components.
func NewResilienceWrapper(name string, cb *CircuitBreaker, bh *Bulkhead, at *AdaptiveTimeout, rp *RetryPolicy) *ResilienceWrapper {
	var cbw *CircuitBreakerWrapper
	if cb != nil {
		cbw = NewCircuitBreakerWrapper(name, cb)
	}

	var bhw *BulkheadWrapper
	if bh != nil {
		bhw = NewBulkheadWrapper(name, bh)
	}

	var atw *AdaptiveTimeoutWrapper
	if at != nil {
		atw = NewAdaptiveTimeoutWrapper(name, at)
	}

	return &ResilienceWrapper{
		circuitBreakerWrapper:  cbw,
		bulkheadWrapper:        bhw,
		adaptiveTimeoutWrapper: atw,
		retryPolicy:            rp,
		name:                   name,
	}
}

// Execute runs the provided function with all configured resilience patterns applied.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute
//
// Returns:
//   - error: An error if any resilience pattern fails
func (w *ResilienceWrapper) Execute(ctx context.Context, fn func() error) error {
	wrappedFn := fn

	if w.circuitBreakerWrapper != nil {
		originalFn := wrappedFn
		wrappedFn = func() error {
			return w.circuitBreakerWrapper.Execute(ctx, originalFn)
		}
	}

	if w.adaptiveTimeoutWrapper != nil {
		originalFn := wrappedFn
		wrappedFn = func() error {
			return w.adaptiveTimeoutWrapper.Execute(ctx, func(ctx context.Context) error {
				return originalFn()
			})
		}
	}

	if w.bulkheadWrapper != nil {
		originalFn := wrappedFn
		wrappedFn = func() error {
			return w.bulkheadWrapper.Execute(ctx, originalFn)
		}
	}

	if w.retryPolicy != nil {
		return w.retryPolicy.Do(ctx, func(ctx context.Context) error {
			return wrappedFn()
		})
	}

	return wrappedFn()
}

// ExecuteWithResult runs the provided function with all configured resilience patterns and returns a result.
//
// Parameters:
//   - ctx: The context for the operation
//   - fn: The function to execute that returns a result and error
//
// Returns:
//   - interface{}: The result from the function
//   - error: An error if any resilience pattern fails
func (w *ResilienceWrapper) ExecuteWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	wrappedFn := fn

	if w.circuitBreakerWrapper != nil {
		originalFn := wrappedFn
		wrappedFn = func() (interface{}, error) {
			return w.circuitBreakerWrapper.ExecuteWithResult(ctx, func() (interface{}, error) {
				return originalFn()
			})
		}
	}

	if w.adaptiveTimeoutWrapper != nil {
		originalFn := wrappedFn
		wrappedFn = func() (interface{}, error) {
			return w.adaptiveTimeoutWrapper.ExecuteWithResult(ctx, func(ctx context.Context) (interface{}, error) {
				return originalFn()
			})
		}
	}

	if w.bulkheadWrapper != nil {
		originalFn := wrappedFn
		wrappedFn = func() (interface{}, error) {
			return w.bulkheadWrapper.ExecuteWithResult(ctx, func() (interface{}, error) {
				return originalFn()
			})
		}
	}

	if w.retryPolicy != nil {
		var result interface{}
		err := w.retryPolicy.Do(ctx, func(ctx context.Context) error {
			var err error
			result, err = wrappedFn()
			return err
		})
		return result, err
	}

	return wrappedFn()
}

// GetStats returns aggregated statistics from all resilience components.
//
// Returns:
//   - ResilienceStats: Combined statistics from all resilience patterns
func (w *ResilienceWrapper) GetStats() ResilienceStats {
	stats := ResilienceStats{
		Name: w.name,
	}

	if w.circuitBreakerWrapper != nil {
		stats.CircuitBreakerStats = w.circuitBreakerWrapper.GetStats()
	}

	if w.bulkheadWrapper != nil {
		stats.BulkheadStats = w.bulkheadWrapper.GetStats()
	}

	if w.adaptiveTimeoutWrapper != nil {
		stats.AdaptiveTimeoutStats = w.adaptiveTimeoutWrapper.GetStats()
	}

	return stats
}

// ResilienceStats aggregates statistics from all resilience components
// (circuit breaker, bulkhead, and adaptive timeout) into a single structure.
type ResilienceStats struct {
	Name                 string
	CircuitBreakerStats  Stats
	BulkheadStats        BulkheadStats
	AdaptiveTimeoutStats AdaptiveTimeoutStats
}

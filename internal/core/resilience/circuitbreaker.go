// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
)

// State represents the state of a circuit breaker.
type State int

// State constants represent the possible states of a circuit breaker.
const (
	// StateClosed allows requests to pass through.
	StateClosed State = iota
	// StateOpen blocks requests from passing through.
	StateOpen
	// StateHalfOpen allows limited requests to test recovery.
	StateHalfOpen
)

// String returns the string representation of the State.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// CircuitBreakerConfig defines the configuration for a circuit breaker.
type CircuitBreakerConfig struct {
	Name             string
	MaxFailures      int
	ResetTimeout     time.Duration
	HalfOpenMaxCalls int
	SuccessThreshold int
	Logger           logs.Logger
	Metrics          *Metrics
}

// DefaultConfig returns the default circuit breaker configuration.
func DefaultConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		MaxFailures:      5,
		ResetTimeout:     60 * time.Second,
		HalfOpenMaxCalls: 3,
		SuccessThreshold: 2,
	}
}

// CircuitBreaker implements the circuit breaker pattern to prevent cascading failures.
type CircuitBreaker struct {
	config        *CircuitBreakerConfig
	state         State
	failures      int
	successes     int
	halfOpenCalls int
	lastFailure   time.Time
	openStartTime time.Time
	mu            sync.RWMutex
	onStateChange func(State, State)
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration.
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	if config == nil {
		config = DefaultConfig()
	}
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// OnStateChange sets a callback function to be called when the circuit breaker state changes.
func (cb *CircuitBreaker) OnStateChange(callback func(State, State)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onStateChange = callback
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Execute runs the given function with circuit breaker protection.
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mu.Lock()
	allowed, newHalfOpenCalls := cb.allowRequestLocked()
	if !allowed {
		currentState := cb.state
		cb.mu.Unlock()

		if cb.config.Logger != nil {
			cb.config.Logger.Warn("Circuit breaker is open",
				logs.String("state", currentState.String()),
			)
		}
		return NewCircuitBreakerOpen("circuit breaker")
	}

	if cb.state == StateHalfOpen {
		cb.halfOpenCalls = newHalfOpenCalls
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	cb.recordResultLocked(err)
	cb.mu.Unlock()

	return err
}

func (cb *CircuitBreaker) allowRequestLocked() (bool, int) {
	switch cb.state {
	case StateClosed:
		return true, 0
	case StateOpen:
		if time.Since(cb.openStartTime) > cb.config.ResetTimeout {
			cb.transitionToLocked(StateHalfOpen)
			return true, 1
		}
		return false, 0
	case StateHalfOpen:
		if cb.halfOpenCalls < cb.config.HalfOpenMaxCalls {
			return true, cb.halfOpenCalls + 1
		}
		return false, cb.halfOpenCalls
	default:
		return false, 0
	}
}

func (cb *CircuitBreaker) recordResultLocked(err error) {
	if err == nil || isClientError(err) || isUnknownErrorType(err) {
		cb.onSuccessLocked()
		return
	}
	cb.onFailureLocked()
}

func isClientError(err error) bool {
	var awsErr interface{ GetHTTPStatusCode() int }
	if errors.As(err, &awsErr) {
		code := awsErr.GetHTTPStatusCode()
		return code >= 400 && code < 500
	}
	return false
}

func isUnknownErrorType(err error) bool {
	var awsErr interface{ GetHTTPStatusCode() int }
	return !errors.As(err, &awsErr)
}

func (cb *CircuitBreaker) onFailureLocked() {
	cb.failures++
	cb.lastFailure = time.Now()

	if cb.config.Logger != nil {
		cb.config.Logger.Debug("Circuit breaker failure",
			logs.Int("failures", cb.failures),
			logs.String("state", cb.state.String()),
		)
	}

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.MaxFailures {
			cb.transitionToLocked(StateOpen)
		}
	case StateHalfOpen:
		cb.transitionToLocked(StateOpen)
	}

	if cb.config.Metrics != nil {
		cb.config.Metrics.RecordError("circuit_breaker_failure")
	}
}

func (cb *CircuitBreaker) onSuccessLocked() {
	cb.successes++

	if cb.config.Logger != nil {
		cb.config.Logger.Debug("Circuit breaker success",
			logs.Int("successes", cb.successes),
			logs.String("state", cb.state.String()),
		)
	}

	switch cb.state {
	case StateClosed:
		cb.failures = 0
		cb.successes = 0
	case StateHalfOpen:
		if cb.successes >= cb.config.SuccessThreshold {
			cb.transitionToLocked(StateClosed)
		}
	}
}

func (cb *CircuitBreaker) transitionToLocked(newState State) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState

	switch newState {
	case StateClosed:
		cb.failures = 0
		cb.successes = 0
		cb.halfOpenCalls = 0
	case StateOpen:
		cb.successes = 0
		cb.halfOpenCalls = 0
		cb.openStartTime = time.Now()
	case StateHalfOpen:
		cb.successes = 0
		cb.halfOpenCalls = 0
	}

	if cb.config.Logger != nil {
		cb.config.Logger.Info("Circuit breaker state changed",
			logs.String("from", oldState.String()),
			logs.String("to", newState.String()),
		)
	}

	if cb.config.Metrics != nil {
		cb.config.Metrics.RecordCircuitStateChange(cb.config.Name, oldState, newState)

		if newState == StateOpen {
			cb.config.Metrics.RecordCircuitOpen(cb.config.Name)
		}

		if newState == StateClosed && !cb.openStartTime.IsZero() {
			duration := time.Since(cb.openStartTime)
			cb.config.Metrics.RecordCircuitOpenDuration(cb.config.Name, duration)
		}
	}

	if cb.onStateChange != nil {
		cb.onStateChange(oldState, newState)
	}
}

// Reset resets the circuit breaker to the closed state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateClosed {
		return
	}

	cb.transitionToLocked(StateClosed)
}

// Stats returns the current statistics of the circuit breaker.
func (cb *CircuitBreaker) Stats() Stats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	var openDuration time.Duration
	if cb.state == StateOpen && !cb.openStartTime.IsZero() {
		openDuration = time.Since(cb.openStartTime)
	}

	return Stats{
		State:         cb.state,
		Failures:      cb.failures,
		Successes:     cb.successes,
		HalfOpenCalls: cb.halfOpenCalls,
		LastFailure:   cb.lastFailure,
		OpenDuration:  openDuration,
	}
}

// Stats represents statistics about the circuit breaker.
type Stats struct {
	State         State
	Failures      int
	Successes     int
	HalfOpenCalls int
	LastFailure   time.Time
	OpenDuration  time.Duration
}

// String returns the string representation of the Stats.
func (s Stats) String() string {
	return fmt.Sprintf("State: %s, Failures: %d, Successes: %d, HalfOpenCalls: %d, LastFailure: %v, OpenDuration: %v",
		s.State, s.Failures, s.Successes, s.HalfOpenCalls, s.LastFailure, s.OpenDuration)
}

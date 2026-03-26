// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"errors"
	"fmt"
)

// Error variables for resilience operations.
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrTimeoutExceeded    = errors.New("timeout exceeded")
	ErrRetryable          = errors.New("retryable error")
)

// CircuitBreakerOpenError is returned when a circuit breaker is open.
type CircuitBreakerOpenError struct {
	Name string
}

// Error returns the error message for CircuitBreakerOpenError.
func (e *CircuitBreakerOpenError) Error() string {
	return fmt.Sprintf("circuit breaker %s is open", e.Name)
}

// Unwrap returns the underlying error for CircuitBreakerOpenError.
func (e *CircuitBreakerOpenError) Unwrap() error {
	return ErrCircuitBreakerOpen
}

// NewCircuitBreakerOpen creates a new circuit breaker open error for the given name.
func NewCircuitBreakerOpen(name string) error {
	return &CircuitBreakerOpenError{Name: name}
}

// IsCircuitBreakerOpen checks whether the given error is a circuit breaker open error.
func IsCircuitBreakerOpen(err error) bool {
	if err == nil {
		return false
	}
	var cbErr *CircuitBreakerOpenError
	return errors.Is(err, ErrCircuitBreakerOpen) || errors.As(err, &cbErr)
}

// RateLimitError is returned when a rate limit is exceeded.
type RateLimitError struct {
	Message string
}

// Error returns the error message for RateLimitError.
func (e *RateLimitError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for RateLimitError.
func (e *RateLimitError) Unwrap() error {
	return ErrRateLimitExceeded
}

// NewRateLimit creates a new rate limit error with the given message.
func NewRateLimit(msg string) error {
	return &RateLimitError{Message: msg}
}

// IsRateLimit checks whether the given error is a rate limit error.
func IsRateLimit(err error) bool {
	if err == nil {
		return false
	}
	var rlErr *RateLimitError
	return errors.Is(err, ErrRateLimitExceeded) || errors.As(err, &rlErr)
}

// TimeoutError is returned when an operation times out.
type TimeoutError struct {
	Message string
}

// Error returns the error message for TimeoutError.
func (e *TimeoutError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error for TimeoutError.
func (e *TimeoutError) Unwrap() error {
	return ErrTimeoutExceeded
}

// NewTimeout creates a new timeout error with the given message.
func NewTimeout(msg string) error {
	return &TimeoutError{Message: msg}
}

// IsTimeout checks whether the given error is a timeout error.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	var toErr *TimeoutError
	return errors.Is(err, ErrTimeoutExceeded) || errors.As(err, &toErr)
}

// HTTPStatusError is an error that contains an HTTP status code.
type HTTPStatusError interface {
	error
	GetHTTPStatus() int
}

// GetHTTPStatus extracts HTTP status code from an error if it implements HTTPStatusError.
// Also supports AWS errors with GetHTTPStatusCode method.
func GetHTTPStatus(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	var httpErr HTTPStatusError
	if errors.As(err, &httpErr) {
		return httpErr.GetHTTPStatus(), true
	}
	// Support AWS errors that use GetHTTPStatusCode instead of GetHTTPStatus
	type awsStatusError interface {
		GetHTTPStatusCode() int
	}
	var awsErr awsStatusError
	if errors.As(err, &awsErr) {
		return awsErr.GetHTTPStatusCode(), true
	}
	return 0, false
}

// IsRetryable checks whether the given error is a retryable error.
// Returns true if the operation should be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrRetryable) {
		return true
	}
	if IsRateLimit(err) {
		return false
	}
	if IsCircuitBreakerOpen(err) {
		return false
	}
	if IsTimeout(err) {
		return false
	}

	if status, ok := GetHTTPStatus(err); ok {
		if status >= 400 && status < 500 {
			return false
		}
		if status >= 500 {
			return true
		}
	}

	return true
}

// IsClientError checks whether the given error is a client-side HTTP error
// (status code 4xx).
func IsClientError(err error) bool {
	if err == nil {
		return false
	}
	var httpErr HTTPStatusError
	if errors.As(err, &httpErr) {
		status := httpErr.GetHTTPStatus()
		return status >= 400 && status < 500
	}
	return false
}

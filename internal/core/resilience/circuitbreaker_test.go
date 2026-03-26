package resilience

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type testInfraError struct {
	code int
	msg  string
}

func (e *testInfraError) Error() string          { return e.msg }
func (e *testInfraError) GetHTTPStatusCode() int { return e.code }

func newInfraError(msg string) *testInfraError {
	return &testInfraError{code: http.StatusInternalServerError, msg: msg}
}

func TestCircuitBreaker_ClosedState(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      3,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", cb.State())
	}

	err := cb.Execute(context.Background(), func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      3,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	for i := 0; i < 3; i++ {
		cb.Execute(context.Background(), func() error {
			return newInfraError("failure")
		})
	}

	if cb.State() != StateOpen {
		t.Fatalf("expected StateOpen after max failures, got %v", cb.State())
	}
}

func TestCircuitBreaker_RejectsRequestsWhenOpen(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      1,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	err := cb.Execute(context.Background(), func() error {
		return nil
	})

	if err == nil {
		t.Fatal("expected circuit breaker open error")
	}

	if !IsCircuitBreakerOpen(err) {
		t.Fatalf("expected circuit breaker open error, got %v", err)
	}
}

func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 2,
	}
	cb := NewCircuitBreaker(config)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	time.Sleep(100 * time.Millisecond)

	stateBefore := cb.State()
	if stateBefore != StateOpen {
		t.Fatalf("expected StateOpen before half-open transition, got %v", stateBefore)
	}

	err := cb.Execute(context.Background(), func() error {
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error during half-open: %v", err)
	}

	if cb.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_ClosesAfterSuccessInHalfOpen(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	time.Sleep(60 * time.Millisecond)

	cb.Execute(context.Background(), func() error {
		return nil
	})

	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed after success in half-open, got %v", cb.State())
	}
}

func TestCircuitBreaker_ReopensOnFailureInHalfOpen(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	time.Sleep(60 * time.Millisecond)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure in half-open")
	})

	if cb.State() != StateOpen {
		t.Fatalf("expected StateOpen after failure in half-open, got %v", cb.State())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      1,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	if cb.State() != StateOpen {
		t.Fatalf("expected StateOpen, got %v", cb.State())
	}

	cb.Reset()

	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed after reset, got %v", cb.State())
	}
}

func TestCircuitBreaker_Stats(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test-stats",
		MaxFailures:      3,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	stats := cb.Stats()

	if stats.State != StateClosed {
		t.Errorf("expected state Closed, got %v", stats.State)
	}
	if stats.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", stats.Failures)
	}
}

func TestCircuitBreaker_OnStateChangeCallback(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      1,
		ResetTimeout:     50 * time.Millisecond,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	var stateChanges []State
	cb.OnStateChange(func(from, to State) {
		stateChanges = append(stateChanges, to)
	})

	cb.Execute(context.Background(), func() error {
		return newInfraError("failure")
	})

	if len(stateChanges) < 1 || stateChanges[0] != StateOpen {
		t.Fatalf("expected state change to Open, got %v", stateChanges)
	}

	time.Sleep(100 * time.Millisecond)

	cb.Execute(context.Background(), func() error {
		return nil
	})

	if len(stateChanges) < 2 || stateChanges[1] != StateHalfOpen {
		t.Fatalf("expected state change to HalfOpen, got %v", stateChanges)
	}
}

func TestCircuitBreaker_ResetsFailuresOnSuccess(t *testing.T) {
	config := &CircuitBreakerConfig{
		Name:             "test",
		MaxFailures:      3,
		ResetTimeout:     1 * time.Second,
		HalfOpenMaxCalls: 2,
		SuccessThreshold: 1,
	}
	cb := NewCircuitBreaker(config)

	for i := 0; i < 2; i++ {
		cb.Execute(context.Background(), func() error {
			return newInfraError(fmt.Sprintf("failure %d", i))
		})
	}

	stats := cb.Stats()
	if stats.Failures != 2 {
		t.Fatalf("expected 2 failures, got %d", stats.Failures)
	}

	cb.Execute(context.Background(), func() error {
		return nil
	})

	stats = cb.Stats()
	if stats.Failures != 0 {
		t.Fatalf("expected failures reset to 0, got %d", stats.Failures)
	}
}

func TestCircuitBreaker_DefaultConfig(t *testing.T) {
	cb := NewCircuitBreaker(nil)

	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", cb.State())
	}

	stats := cb.Stats()
	if stats.Failures != 0 {
		t.Fatalf("expected 0 failures, got %d", stats.Failures)
	}
}

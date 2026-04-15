package resilience

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryPolicy_Success(t *testing.T) {
	policy := NewRetryPolicy()

	callCount := 0
	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 1 {
		t.Fatalf("expected 1 call, got %d", callCount)
	}
}

func TestRetryPolicy_RetriesOnError(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(3)

	callCount := 0
	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

func TestRetryPolicy_ReturnsErrorAfterMaxAttempts(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(3)

	callCount := 0
	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		return errors.New("persistent error")
	})

	if err == nil {
		t.Fatal("expected error after max attempts")
	}

	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

func TestRetryPolicy_DoesNotRetryNonRetryableError(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(3)

	callCount := 0
	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		return NewRateLimit("rate limited")
	})

	if err == nil {
		t.Fatal("expected error")
	}

	if callCount != 1 {
		t.Fatalf("expected 1 call for non-retryable error, got %d", callCount)
	}
}

func TestRetryPolicy_ContextCancel(t *testing.T) {
	policy := NewRetryPolicy()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := policy.Do(ctx, func(ctx context.Context) error {
		return errors.New("error")
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetryPolicy_ExponentialBackoff(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(4)
	policy.SetInitialDelay(10 * time.Millisecond)
	policy.SetMaxDelay(100 * time.Millisecond)

	backoff := &ExponentialBackoff{
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		jitter:       &NoJitter{},
	}
	policy.SetBackoff(backoff)

	callCount := 0
	startTime := time.Now()

	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		if callCount < 4 {
			return errors.New("error")
		}
		return nil
	})

	elapsed := time.Since(startTime)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 4 {
		t.Fatalf("expected 4 calls, got %d", callCount)
	}

	if elapsed < 30*time.Millisecond {
		t.Fatalf("expected at least 30ms delay, got %v", elapsed)
	}
}

func TestRetryPolicy_LinearBackoff(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(3)
	policy.SetInitialDelay(10 * time.Millisecond)

	backoff := &LinearBackoff{
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Increment:    10 * time.Millisecond,
		jitter:       &NoJitter{},
	}
	policy.SetBackoff(backoff)

	callCount := 0
	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			return errors.New("error")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

func TestRetryPolicy_WithMetrics(t *testing.T) {
	policy := NewRetryPolicy()
	metrics := NewMetrics()
	policy.SetMetrics(metrics)

	callCount := 0
	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			return errors.New("error")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	retryCount := metrics.GetRetryCount(policy.name)
	if retryCount != 2 {
		t.Fatalf("expected 2 retries recorded, got %d", retryCount)
	}
}

func TestRetryPolicy_MaxDelay(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(5)
	policy.SetInitialDelay(10 * time.Millisecond)
	policy.SetMaxDelay(50 * time.Millisecond)

	backoff := &ExponentialBackoff{
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     50 * time.Millisecond,
		Multiplier:   10.0,
		jitter:       &NoJitter{},
	}
	policy.SetBackoff(backoff)

	callCount := 0
	startTime := time.Now()

	err := policy.Do(context.Background(), func(ctx context.Context) error {
		callCount++
		if callCount < 5 {
			return errors.New("error")
		}
		return nil
	})

	elapsed := time.Since(startTime)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if elapsed > 300*time.Millisecond {
		t.Fatalf("expected elapsed less than 300ms due to max delay, got %v", elapsed)
	}
}

func TestRetryPolicy_GetMaxAttempts(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(5)

	if policy.GetMaxAttempts() != 5 {
		t.Fatalf("expected 5, got %d", policy.GetMaxAttempts())
	}
}

func TestRetryPolicy_GetMaxDelay(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxDelay(10 * time.Second)

	if policy.GetMaxDelay() != 10*time.Second {
		t.Fatalf("expected 10s, got %v", policy.GetMaxDelay())
	}
}

func TestRetryPolicy_GetInitialDelay(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetInitialDelay(500 * time.Millisecond)

	if policy.GetInitialDelay() != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", policy.GetInitialDelay())
	}
}

func TestRetryPolicy_SetName(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetName("custom-name")

	if policy.GetName() != "custom-name" {
		t.Fatalf("expected 'custom-name', got '%s'", policy.GetName())
	}
}

func TestFullJitter(t *testing.T) {
	jitter := &FullJitter{}

	var min, max time.Duration
	delay := 100 * time.Millisecond
	for i := 0; i < 1000; i++ {
		result := jitter.Apply(delay)
		if result < 0 || result >= delay {
			t.Fatalf("expected result in [0, %v), got %v", delay, result)
		}
		if i == 0 || result < min {
			min = result
		}
		if i == 0 || result > max {
			max = result
		}
	}

	if min == max {
		t.Fatalf("expected variance in jitter output, all values were %v", min)
	}
}

func TestEqualJitter(t *testing.T) {
	jitter := &EqualJitter{}

	for i := 0; i < 100; i++ {
		result := jitter.Apply(100 * time.Millisecond)
		if result < 50*time.Millisecond || result > 100*time.Millisecond {
			t.Fatalf("expected result between 50ms and 100ms, got %v", result)
		}
	}
}

func TestNoJitter(t *testing.T) {
	jitter := &NoJitter{}

	result := jitter.Apply(100 * time.Millisecond)
	if result != 100*time.Millisecond {
		t.Fatalf("expected 100ms, got %v", result)
	}
}

func TestNewRetryPolicyWithBackoff(t *testing.T) {
	backoff := &ExponentialBackoff{
		InitialDelay: 50 * time.Millisecond,
		MaxDelay:     500 * time.Millisecond,
		Multiplier:   2.0,
	}
	policy := NewRetryPolicyWithBackoff(backoff)

	got := policy.GetBackoff()
	if got != backoff {
		t.Fatal("GetBackoff() should return the same backoff instance that was passed to constructor")
	}
}

func TestRetryPolicy_ConcurrentCalls(t *testing.T) {
	policy := NewRetryPolicy()
	policy.SetMaxAttempts(3)

	var callCount int32

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			policy.Do(context.Background(), func(ctx context.Context) error {
				atomic.AddInt32(&callCount, 1)
				return nil
			})
		}()
	}

	wg.Wait()

	if callCount != 10 {
		t.Fatalf("expected 10 total calls, got %d", callCount)
	}
}

func TestRetryPolicy_GetBackoff(t *testing.T) {
	policy := NewRetryPolicyWithBackoff(&ExponentialBackoff{Multiplier: 2.0})

	backoff := policy.GetBackoff()
	if backoff == nil {
		t.Fatal("expected non-nil backoff")
	}
	eb, ok := backoff.(*ExponentialBackoff)
	if !ok {
		t.Fatal("expected ExponentialBackoff")
	}
	if eb.Multiplier != 2.0 {
		t.Fatalf("expected multiplier 2.0, got %v", eb.Multiplier)
	}
}

func TestRetryPolicy_GetJitter(t *testing.T) {
	policy := NewRetryPolicy()
	fj := &FullJitter{}
	policy.SetJitter(fj)

	jitter := policy.GetJitter()
	if jitter != fj {
		t.Fatal("expected FullJitter")
	}
}

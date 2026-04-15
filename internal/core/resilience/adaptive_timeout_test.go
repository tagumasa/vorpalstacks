package resilience

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestAdaptiveTimeout_Execute_Success(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   1 * time.Second,
		MinTimeout:       100 * time.Millisecond,
		MaxTimeout:       5 * time.Second,
		SuccessThreshold: 2,
		FailureThreshold: 2,
		AdjustmentFactor: 1.5,
	}
	at := NewAdaptiveTimeout(config)

	err := at.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdaptiveTimeout_Execute_Timeout(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   50 * time.Millisecond,
		MinTimeout:       10 * time.Millisecond,
		MaxTimeout:       100 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 1,
		AdjustmentFactor: 1.5,
	}
	at := NewAdaptiveTimeout(config)

	err := at.Execute(context.Background(), func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !IsTimeout(err) {
		t.Fatalf("expected timeout error, got %v", err)
	}
}

func TestAdaptiveTimeout_ContextCancel(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   1 * time.Second,
		MinTimeout:       100 * time.Millisecond,
		MaxTimeout:       5 * time.Second,
		SuccessThreshold: 1,
		FailureThreshold: 1,
	}
	at := NewAdaptiveTimeout(config)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := at.Execute(ctx, func(ctx context.Context) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	if err == nil {
		t.Fatal("expected context cancel error")
	}
}

func TestAdaptiveTimeout_IncreasesTimeoutOnFailure(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   100 * time.Millisecond,
		MinTimeout:       50 * time.Millisecond,
		MaxTimeout:       500 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 1,
		AdjustmentFactor: 2.0,
	}
	at := NewAdaptiveTimeout(config)

	initialTimeout := at.GetTimeout()

	at.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("error")
	})

	newTimeout := at.GetTimeout()

	if newTimeout <= initialTimeout {
		t.Fatalf("expected timeout to increase, initial=%v, new=%v", initialTimeout, newTimeout)
	}
}

func TestAdaptiveTimeout_DecreasesTimeoutOnSuccess(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   200 * time.Millisecond,
		MinTimeout:       10 * time.Millisecond,
		MaxTimeout:       500 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 2,
		AdjustmentFactor: 2.0,
	}
	at := NewAdaptiveTimeout(config)

	initialTimeout := at.GetTimeout()

	for i := 0; i < 5; i++ {
		at.Execute(context.Background(), func(ctx context.Context) error {
			return nil
		})
	}

	newTimeout := at.GetTimeout()

	if newTimeout >= initialTimeout {
		t.Fatalf("expected timeout to decrease, initial=%v, new=%v", initialTimeout, newTimeout)
	}

	if newTimeout > config.MinTimeout {
		t.Fatalf("expected timeout <= MinTimeout, got %v", newTimeout)
	}
}

func TestAdaptiveTimeout_CapsAtMaxTimeout(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   100 * time.Millisecond,
		MinTimeout:       50 * time.Millisecond,
		MaxTimeout:       200 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 1,
		AdjustmentFactor: 10.0,
	}
	at := NewAdaptiveTimeout(config)

	for i := 0; i < 5; i++ {
		at.Execute(context.Background(), func(ctx context.Context) error {
			return errors.New("error")
		})
	}

	timeout := at.GetTimeout()
	if timeout > config.MaxTimeout {
		t.Fatalf("timeout should not exceed MaxTimeout, got %v", timeout)
	}
}

func TestAdaptiveTimeout_Reset(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   100 * time.Millisecond,
		MinTimeout:       50 * time.Millisecond,
		MaxTimeout:       500 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 1,
		AdjustmentFactor: 2.0,
	}
	at := NewAdaptiveTimeout(config)

	at.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("error")
	})

	at.Reset()

	timeout := at.GetTimeout()
	if timeout != config.InitialTimeout {
		t.Fatalf("expected timeout to reset to %v, got %v", config.InitialTimeout, timeout)
	}
}

func TestAdaptiveTimeout_Stats(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test-stats",
		InitialTimeout:   100 * time.Millisecond,
		MinTimeout:       50 * time.Millisecond,
		MaxTimeout:       500 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 1,
		AdjustmentFactor: 1.0,
	}
	at := NewAdaptiveTimeout(config)

	at.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})

	stats := at.Stats()

	if stats.Name != "test-stats" {
		t.Errorf("expected name 'test-stats', got '%s'", stats.Name)
	}

	if stats.ConsecutiveSuccesses != 0 {
		t.Errorf("expected 0 consecutive success after reset, got %d", stats.ConsecutiveSuccesses)
	}
}

func TestAdaptiveTimeout_GetTimeout(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		InitialTimeout: 30 * time.Second,
	}
	at := NewAdaptiveTimeout(config)

	timeout := at.GetTimeout()
	if timeout != 30*time.Second {
		t.Fatalf("expected 30s, got %v", timeout)
	}
}

func TestAdaptiveTimeout_MultipleSuccessesDecreaseTimeout(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   100 * time.Millisecond,
		MinTimeout:       25 * time.Millisecond,
		MaxTimeout:       200 * time.Millisecond,
		SuccessThreshold: 2,
		FailureThreshold: 1,
		AdjustmentFactor: 2.0,
	}
	at := NewAdaptiveTimeout(config)

	initialTimeout := at.GetTimeout()

	at.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})

	afterFirst := at.GetTimeout()
	if afterFirst != initialTimeout {
		t.Fatalf("timeout should not change after 1 success (threshold=2), got %v", afterFirst)
	}

	at.Execute(context.Background(), func(ctx context.Context) error {
		return nil
	})

	afterSecond := at.GetTimeout()
	if afterSecond >= initialTimeout {
		t.Fatalf("timeout should decrease after 2 successes, got %v", afterSecond)
	}
}

func TestAdaptiveTimeout_MultipleFailuresIncreaseTimeout(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test",
		InitialTimeout:   100 * time.Millisecond,
		MinTimeout:       25 * time.Millisecond,
		MaxTimeout:       200 * time.Millisecond,
		SuccessThreshold: 1,
		FailureThreshold: 2,
		AdjustmentFactor: 2.0,
	}
	at := NewAdaptiveTimeout(config)

	initialTimeout := at.GetTimeout()

	at.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("error")
	})

	afterFirst := at.GetTimeout()
	if afterFirst != initialTimeout {
		t.Fatalf("timeout should not change after 1 failure (threshold=2), got %v", afterFirst)
	}

	at.Execute(context.Background(), func(ctx context.Context) error {
		return errors.New("error")
	})

	afterSecond := at.GetTimeout()
	if afterSecond <= afterFirst {
		t.Fatalf("timeout should increase after 2 failures, got %v", afterSecond)
	}
}

func TestAdaptiveTimeout_DefaultConfig(t *testing.T) {
	at := NewAdaptiveTimeout(nil)

	timeout := at.GetTimeout()
	if timeout != 30*time.Second {
		t.Fatalf("expected default 30s, got %v", timeout)
	}
}

func TestAdaptiveTimeout_Close(t *testing.T) {
	config := &AdaptiveTimeoutConfig{
		Name:             "test-close",
		InitialTimeout:   1 * time.Second,
		MinTimeout:       100 * time.Millisecond,
		MaxTimeout:       5 * time.Second,
		SuccessThreshold: 2,
		FailureThreshold: 2,
		AdjustmentFactor: 1.5,
	}
	at := NewAdaptiveTimeout(config)

	err := at.Close()
	if err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}

	err = at.Close()
	if err != nil {
		t.Fatalf("double Close() returned error: %v", err)
	}
}

func TestAdaptiveTimeoutStats_String(t *testing.T) {
	s := AdaptiveTimeoutStats{
		Name:                 "test",
		CurrentTimeout:       3 * time.Second,
		ConsecutiveSuccesses: 5,
		ConsecutiveFailures:  1,
	}
	got := s.String()
	if !strings.Contains(got, "test") {
		t.Fatalf("AdaptiveTimeoutStats.String() should contain name, got %q", got)
	}
	if !strings.Contains(got, "Timeout=") {
		t.Fatalf("AdaptiveTimeoutStats.String() should contain Timeout, got %q", got)
	}
	if !strings.Contains(got, "Successes=5") {
		t.Fatalf("AdaptiveTimeoutStats.String() should contain Successes=5, got %q", got)
	}
	if !strings.Contains(got, "Failures=1") {
		t.Fatalf("AdaptiveTimeoutStats.String() should contain Failures=1, got %q", got)
	}
}

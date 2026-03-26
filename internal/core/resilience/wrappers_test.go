package resilience

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreakerWrapper_Execute_Success(t *testing.T) {
	cfg := DefaultConfig()
	cb := NewCircuitBreaker(cfg)
	wrapper := NewCircuitBreakerWrapper("test", cb)

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, StateClosed, wrapper.GetState())
}

func TestCircuitBreakerWrapper_Execute_Failure(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxFailures = 2
	cb := NewCircuitBreaker(cfg)
	wrapper := NewCircuitBreakerWrapper("test", cb)

	wrapper.Execute(context.Background(), func() error {
		return &testInfraError{code: 500, msg: "failure"}
	})
	wrapper.Execute(context.Background(), func() error {
		return &testInfraError{code: 500, msg: "failure"}
	})

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.Error(t, err)
	assert.True(t, IsCircuitBreakerOpen(err))
}

func TestCircuitBreakerWrapper_ExecuteWithResult_Success(t *testing.T) {
	cfg := DefaultConfig()
	cb := NewCircuitBreaker(cfg)
	wrapper := NewCircuitBreakerWrapper("test", cb)

	result, err := wrapper.ExecuteWithResult(context.Background(), func() (interface{}, error) {
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
}

func TestCircuitBreakerWrapper_GetStats(t *testing.T) {
	cfg := DefaultConfig()
	cb := NewCircuitBreaker(cfg)
	wrapper := NewCircuitBreakerWrapper("test", cb)

	stats := wrapper.GetStats()

	assert.NotNil(t, stats)
}

func TestCircuitBreakerWrapper_Reset(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxFailures = 1
	cb := NewCircuitBreaker(cfg)
	wrapper := NewCircuitBreakerWrapper("test", cb)

	wrapper.Execute(context.Background(), func() error {
		return &testInfraError{code: 500, msg: "failure"}
	})

	wrapper.Reset()

	assert.Equal(t, StateClosed, wrapper.GetState())
}

func TestBulkheadWrapper_Execute_Success(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	bh := NewBulkhead(cfg)
	wrapper := NewBulkheadWrapper("test", bh)

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.NoError(t, err)
}

func TestBulkheadWrapper_Execute(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	bh := NewBulkhead(cfg)
	wrapper := NewBulkheadWrapper("test", bh)

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.NoError(t, err)
}

func TestBulkheadWrapper_ExecuteWithResult(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	bh := NewBulkhead(cfg)
	wrapper := NewBulkheadWrapper("test", bh)

	result, err := wrapper.ExecuteWithResult(context.Background(), func() (interface{}, error) {
		return "result", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "result", result)
}

func TestBulkheadWrapper_GetActiveCount(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	bh := NewBulkhead(cfg)
	wrapper := NewBulkheadWrapper("test", bh)

	assert.Equal(t, 0, wrapper.GetActiveCount())
}

func TestBulkheadWrapper_GetAvailableSlots(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	cfg.MaxConcurrent = 10
	bh := NewBulkhead(cfg)
	wrapper := NewBulkheadWrapper("test", bh)

	slots := wrapper.GetAvailableSlots()
	assert.Equal(t, 10, slots)
}

func TestAdaptiveTimeoutWrapper_Execute_Success(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewAdaptiveTimeoutWrapper("test", at)

	err := wrapper.Execute(context.Background(), func(context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestAdaptiveTimeoutWrapper_Execute_Timeout(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	cfg.InitialTimeout = 50 * time.Millisecond
	cfg.MinTimeout = 50 * time.Millisecond
	cfg.MaxTimeout = 50 * time.Millisecond
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewAdaptiveTimeoutWrapper("test", at)

	err := wrapper.Execute(context.Background(), func(context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	assert.True(t, IsTimeout(err))
}

func TestAdaptiveTimeoutWrapper_ExecuteWithResult(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewAdaptiveTimeoutWrapper("test", at)

	result, err := wrapper.ExecuteWithResult(context.Background(), func(context.Context) (interface{}, error) {
		return "value", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "value", result)
}

func TestAdaptiveTimeoutWrapper_GetTimeout(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	cfg.InitialTimeout = 5 * time.Second
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewAdaptiveTimeoutWrapper("test", at)

	timeout := wrapper.GetTimeout()

	assert.Equal(t, int64(5000), timeout)
}

func TestAdaptiveTimeoutWrapper_Reset(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewAdaptiveTimeoutWrapper("test", at)

	wrapper.Reset()

	stats := wrapper.GetStats()
	assert.Equal(t, 0, stats.ConsecutiveFailures)
}

func TestResilienceWrapper_Execute_NoWrappers(t *testing.T) {
	wrapper := NewResilienceWrapper("test", nil, nil, nil, nil)

	executed := false
	err := wrapper.Execute(context.Background(), func() error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestResilienceWrapper_Execute_WithCircuitBreaker(t *testing.T) {
	cfg := DefaultConfig()
	cb := NewCircuitBreaker(cfg)
	wrapper := NewResilienceWrapper("test", cb, nil, nil, nil)

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.NoError(t, err)
}

func TestResilienceWrapper_Execute_WithBulkhead(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	bh := NewBulkhead(cfg)
	wrapper := NewResilienceWrapper("test", nil, bh, nil, nil)

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.NoError(t, err)
}

func TestResilienceWrapper_Execute_WithAdaptiveTimeout(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewResilienceWrapper("test", nil, nil, at, nil)

	err := wrapper.Execute(context.Background(), func() error {
		return nil
	})

	assert.NoError(t, err)
}

func TestResilienceWrapper_Execute_WithRetry(t *testing.T) {
	rp := NewRetryPolicy()
	wrapper := NewResilienceWrapper("test", nil, nil, nil, rp)

	callCount := 0
	err := wrapper.Execute(context.Background(), func() error {
		callCount++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestResilienceWrapper_ExecuteWithResult(t *testing.T) {
	wrapper := NewResilienceWrapper("test", nil, nil, nil, nil)

	result, err := wrapper.ExecuteWithResult(context.Background(), func() (interface{}, error) {
		return "test-result", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "test-result", result)
}

func TestResilienceWrapper_GetStats_Empty(t *testing.T) {
	wrapper := NewResilienceWrapper("test", nil, nil, nil, nil)

	stats := wrapper.GetStats()

	assert.Equal(t, "test", stats.Name)
}

func TestResilienceWrapper_GetStats_WithCircuitBreaker(t *testing.T) {
	cfg := DefaultConfig()
	cb := NewCircuitBreaker(cfg)
	wrapper := NewResilienceWrapper("test", cb, nil, nil, nil)

	stats := wrapper.GetStats()

	assert.NotNil(t, stats.CircuitBreakerStats)
}

func TestResilienceWrapper_GetStats_WithBulkhead(t *testing.T) {
	cfg := DefaultBulkheadConfig()
	bh := NewBulkhead(cfg)
	wrapper := NewResilienceWrapper("test", nil, bh, nil, nil)

	stats := wrapper.GetStats()

	assert.NotNil(t, stats.BulkheadStats)
}

func TestResilienceWrapper_GetStats_WithAdaptiveTimeout(t *testing.T) {
	cfg := DefaultAdaptiveTimeoutConfig()
	at := NewAdaptiveTimeout(cfg)
	wrapper := NewResilienceWrapper("test", nil, nil, at, nil)

	stats := wrapper.GetStats()

	assert.NotNil(t, stats.AdaptiveTimeoutStats)
}

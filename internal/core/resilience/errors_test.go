package resilience

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreakerOpenError(t *testing.T) {
	err := NewCircuitBreakerOpen("test-circuit")
	assert.Equal(t, "circuit breaker test-circuit is open", err.Error())
	assert.True(t, IsCircuitBreakerOpen(err))
	assert.True(t, errors.Is(err, ErrCircuitBreakerOpen))
}

func TestRateLimitError(t *testing.T) {
	err := NewRateLimit("rate limit exceeded")
	assert.Equal(t, "rate limit exceeded", err.Error())
	assert.True(t, IsRateLimit(err))
	assert.True(t, errors.Is(err, ErrRateLimitExceeded))
}

func TestTimeoutError(t *testing.T) {
	err := NewTimeout("operation timed out")
	assert.Equal(t, "operation timed out", err.Error())
	assert.True(t, IsTimeout(err))
	assert.True(t, errors.Is(err, ErrTimeoutExceeded))
}

func TestIsRetryable(t *testing.T) {
	assert.True(t, IsRetryable(assert.AnError))
	assert.False(t, IsRetryable(NewRateLimit("rate limit")))
	assert.False(t, IsRetryable(NewCircuitBreakerOpen("test")))
	assert.False(t, IsRetryable(NewTimeout("timeout")))
	assert.True(t, IsRetryable(ErrRetryable))
}

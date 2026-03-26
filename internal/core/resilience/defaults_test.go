package resilience

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRetryValues(t *testing.T) {
	assert.Equal(t, 3, DefaultRetryMaxAttempts)
	assert.Equal(t, 30*time.Second, DefaultRetryMaxDelay)
	assert.Equal(t, 100*time.Millisecond, DefaultRetryInitialDelay)
	assert.Equal(t, 2.0, DefaultRetryBackoffMultiplier)
}

func TestDefaultCircuitBreakerValues(t *testing.T) {
	assert.Equal(t, 5, DefaultCircuitBreakerMaxFailures)
	assert.Equal(t, 60*time.Second, DefaultCircuitBreakerResetTimeout)
	assert.Equal(t, 3, DefaultCircuitBreakerHalfOpenRequests)
	assert.Equal(t, 2, DefaultCircuitBreakerSuccessThreshold)
}

func TestDefaultBulkheadValues(t *testing.T) {
	assert.Equal(t, 100, DefaultBulkheadMaxConcurrent)
	assert.Equal(t, 1000, DefaultBulkheadMaxQueueSize)
	assert.Equal(t, 5*time.Second, DefaultBulkheadTimeout)
}

func TestDefaultAdaptiveTimeoutValues(t *testing.T) {
	assert.Equal(t, 5*time.Second, DefaultAdaptiveTimeoutMin)
	assert.Equal(t, 60*time.Second, DefaultAdaptiveTimeoutMax)
	assert.Equal(t, 30*time.Second, DefaultAdaptiveTimeoutDefault)
	assert.Equal(t, 5, DefaultAdaptiveTimeoutSuccessThreshold)
	assert.Equal(t, 2, DefaultAdaptiveTimeoutFailureThreshold)
}

func TestDefaultCacheValues(t *testing.T) {
	assert.Equal(t, 5*time.Minute, DefaultCacheTTL)
	assert.Equal(t, 1000, DefaultCacheMaxSize)
	assert.Equal(t, 1*time.Minute, DefaultCacheCleanupInterval)
}

func TestDefaultHealthCheckValues(t *testing.T) {
	assert.Equal(t, 5*time.Second, DefaultHealthCheckTimeout)
	assert.Equal(t, 30*time.Second, DefaultHealthCheckInterval)
	assert.Equal(t, 3, DefaultHealthCheckFailureThreshold)
}

func TestDefaultAdaptiveAdjustmentFactor(t *testing.T) {
	assert.Equal(t, 1.5, DefaultAdaptiveAdjustmentFactor)
}

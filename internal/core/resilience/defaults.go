// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import "time"

// Default constants for resilience patterns.
const (
	// DefaultRetryMaxAttempts is the default maximum number of retry attempts.
	DefaultRetryMaxAttempts = 3
	// DefaultRetryMaxDelay is the default maximum delay between retry attempts.
	DefaultRetryMaxDelay = 30 * time.Second
	// DefaultRetryInitialDelay is the default initial delay for retry attempts.
	DefaultRetryInitialDelay = 100 * time.Millisecond
	// DefaultRetryBackoffMultiplier is the default multiplier for exponential backoff.
	DefaultRetryBackoffMultiplier = 2.0

	// DefaultCircuitBreakerMaxFailures is the default maximum number of failures before opening the circuit.
	DefaultCircuitBreakerMaxFailures = 9
	// DefaultCircuitBreakerResetTimeout is the default timeout before attempting to close the circuit.
	DefaultCircuitBreakerResetTimeout = 60 * time.Second
	// DefaultCircuitBreakerHalfOpenRequests is the default number of requests in half-open state.
	DefaultCircuitBreakerHalfOpenRequests = 5
	// DefaultCircuitBreakerSuccessThreshold is the default number of successes required to close the circuit.
	DefaultCircuitBreakerSuccessThreshold = 2

	// DefaultBulkheadMaxConcurrent is the default maximum number of concurrent executions.
	DefaultBulkheadMaxConcurrent = 100
	// DefaultBulkheadMaxQueueSize is the default maximum size of the bulkhead queue.
	DefaultBulkheadMaxQueueSize = 1000
	// DefaultBulkheadTimeout is the default timeout for bulkhead executions.
	DefaultBulkheadTimeout = 5 * time.Second

	// DefaultAdaptiveTimeoutMin is the default minimum timeout for adaptive timeout.
	DefaultAdaptiveTimeoutMin = 5 * time.Second
	// DefaultAdaptiveTimeoutMax is the default maximum timeout for adaptive timeout.
	DefaultAdaptiveTimeoutMax = 60 * time.Second
	// DefaultAdaptiveTimeoutDefault is the default timeout for adaptive timeout.
	DefaultAdaptiveTimeoutDefault = 30 * time.Second
	// DefaultAdaptiveTimeoutSuccessThreshold is the default number of successes to decrease timeout.
	DefaultAdaptiveTimeoutSuccessThreshold = 5
	// DefaultAdaptiveTimeoutFailureThreshold is the default number of failures to increase timeout.
	DefaultAdaptiveTimeoutFailureThreshold = 2

	// DefaultCacheTTL is the default time-to-live for cached items.
	DefaultCacheTTL = 5 * time.Minute
	// DefaultCacheMaxSize is the default maximum size of the cache.
	DefaultCacheMaxSize = 1000
	// DefaultCacheCleanupInterval is the default interval for cache cleanup.
	DefaultCacheCleanupInterval = 1 * time.Minute

	// DefaultHealthCheckTimeout is the default timeout for health checks.
	DefaultHealthCheckTimeout = 5 * time.Second
	// DefaultHealthCheckInterval is the default interval between health checks.
	DefaultHealthCheckInterval = 30 * time.Second
	// DefaultHealthCheckFailureThreshold is the default number of failures before marking as unhealthy.
	DefaultHealthCheckFailureThreshold = 3

	// DefaultAdaptiveAdjustmentFactor is the default factor for adjusting adaptive timeouts.
	DefaultAdaptiveAdjustmentFactor = 1.5
)

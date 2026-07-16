// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"sort"
	"sync"
	"time"
)

// Metrics collects and tracks various resilience-related metrics including
// error counts, retry counts, latencies, circuit breaker state, and health check results.
type Metrics struct {
	mu sync.RWMutex

	errorCounts          map[string]int64
	errorRates           map[string]float64
	retryCounts          map[string]int64
	circuitOpenCount     int64
	latencies            map[string][]time.Duration
	avgLatencies         map[string]time.Duration
	p95Latencies         map[string]time.Duration
	p99Latencies         map[string]time.Duration
	circuitStateChanges  map[string]int64
	circuitOpenDuration  map[string]time.Duration
	healthCheckFailures  map[string]int64
	healthCheckSuccesses map[string]int64
}

// NewMetrics creates a new metrics collector.
func NewMetrics() *Metrics {
	return &Metrics{
		errorCounts:          make(map[string]int64),
		errorRates:           make(map[string]float64),
		retryCounts:          make(map[string]int64),
		latencies:            make(map[string][]time.Duration),
		avgLatencies:         make(map[string]time.Duration),
		p95Latencies:         make(map[string]time.Duration),
		p99Latencies:         make(map[string]time.Duration),
		circuitStateChanges:  make(map[string]int64),
		circuitOpenDuration:  make(map[string]time.Duration),
		healthCheckFailures:  make(map[string]int64),
		healthCheckSuccesses: make(map[string]int64),
	}
}

// GlobalMetrics is the global metrics instance shared across the application.
var GlobalMetrics = NewMetrics()

// RecordError records an error occurrence.
//
// Parameters:
//   - errorType: The type/category of the error
func (m *Metrics) RecordError(errorType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCounts[errorType]++
}

// GetErrorCount returns the count of errors for a specific error type.
//
// Parameters:
//   - errorType: The type of error
//
// Returns:
//   - int64: The number of errors recorded
func (m *Metrics) GetErrorCount(errorType string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.errorCounts[errorType]
}

// GetAllErrorCounts returns all error counts.
//
// Returns:
//   - map[string]int64: A map of error types to their counts
func (m *Metrics) GetAllErrorCounts() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]int64)
	for k, v := range m.errorCounts {
		result[k] = v
	}
	return result
}

// RecordRetry records a retry occurrence.
//
// Parameters:
//   - operation: The name of the operation that was retried
func (m *Metrics) RecordRetry(operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.retryCounts[operation]++
}

// GetRetryCount returns the count of retries for a specific operation.
//
// Parameters:
//   - operation: The name of the operation
//
// Returns:
//   - int64: The number of retries recorded
func (m *Metrics) GetRetryCount(operation string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.retryCounts[operation]
}

// GetAllRetryCounts returns all retry counts.
//
// Returns:
//   - map[string]int64: A map of operation names to their retry counts
func (m *Metrics) GetAllRetryCounts() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]int64)
	for k, v := range m.retryCounts {
		result[k] = v
	}
	return result
}

// RecordLatency records the latency of an operation for percentile calculations.
//
// Parameters:
//   - operation: The name of the operation
//   - latency: The latency duration
func (m *Metrics) RecordLatency(operation string, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	samples := m.latencies[operation]
	if len(samples) >= 1000 {
		samples = samples[1:]
	}
	samples = append(samples, latency)
	m.latencies[operation] = samples

	m.calculateLatencyMetrics(operation, samples)
}

func (m *Metrics) calculateLatencyMetrics(operation string, samples []time.Duration) {
	if len(samples) == 0 {
		return
	}

	var sum time.Duration
	for _, s := range samples {
		sum += s
	}
	m.avgLatencies[operation] = sum / time.Duration(len(samples))

	sorted := make([]time.Duration, len(samples))
	copy(sorted, samples)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	p95Index := int(float64(len(sorted)) * 0.95)
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	m.p95Latencies[operation] = sorted[p95Index]

	p99Index := int(float64(len(sorted)) * 0.99)
	if p99Index >= len(sorted) {
		p99Index = len(sorted) - 1
	}
	m.p99Latencies[operation] = sorted[p99Index]
}

// GetAvgLatency returns the average latency for a specific operation.
//
// Parameters:
//   - operation: The name of the operation
//
// Returns:
//   - time.Duration: The average latency
func (m *Metrics) GetAvgLatency(operation string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.avgLatencies[operation]
}

// GetP95Latency returns the 95th percentile latency for a specific operation.
//
// Parameters:
//   - operation: The name of the operation
//
// Returns:
//   - time.Duration: The 95th percentile latency
func (m *Metrics) GetP95Latency(operation string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.p95Latencies[operation]
}

// GetP99Latency returns the 99th percentile latency for a specific operation.
//
// Parameters:
//   - operation: The name of the operation
//
// Returns:
//   - time.Duration: The 99th percentile latency
func (m *Metrics) GetP99Latency(operation string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.p99Latencies[operation]
}

// RecordCircuitOpen records a circuit breaker being opened.
//
// Parameters:
//   - circuitName: The name of the circuit breaker
func (m *Metrics) RecordCircuitOpen(circuitName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.circuitOpenCount++
}

// GetCircuitOpenCount returns the total number of times circuits have been opened.
//
// Returns:
//   - int64: The total count of circuit openings
func (m *Metrics) GetCircuitOpenCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.circuitOpenCount
}

// RecordCircuitStateChange records a state change in a circuit breaker.
//
// Parameters:
//   - circuitName: The name of the circuit breaker
//   - from: The previous state
//   - to: The new state
func (m *Metrics) RecordCircuitStateChange(circuitName string, from, to State) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.circuitStateChanges[circuitName]++
}

// GetCircuitStateChanges returns the number of state changes for a circuit breaker.
//
// Parameters:
//   - circuitName: The name of the circuit breaker
//
// Returns:
//   - int64: The number of state changes
func (m *Metrics) GetCircuitStateChanges(circuitName string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.circuitStateChanges[circuitName]
}

// RecordCircuitOpenDuration records how long a circuit breaker remained open.
//
// Parameters:
//   - circuitName: The name of the circuit breaker
//   - duration: The duration the circuit was open
func (m *Metrics) RecordCircuitOpenDuration(circuitName string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.circuitOpenDuration[circuitName] = duration
}

// GetCircuitOpenDuration returns the duration a circuit breaker remained open.
//
// Parameters:
//   - circuitName: The name of the circuit breaker
//
// Returns:
//   - time.Duration: The total open duration
func (m *Metrics) GetCircuitOpenDuration(circuitName string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.circuitOpenDuration[circuitName]
}

// RecordHealthCheckFailure records a health check failure.
//
// Parameters:
//   - checkName: The name of the health check
func (m *Metrics) RecordHealthCheckFailure(checkName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.healthCheckFailures[checkName]++
}

// RecordHealthCheckSuccess records a successful health check.
//
// Parameters:
//   - checkName: The name of the health check
func (m *Metrics) RecordHealthCheckSuccess(checkName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.healthCheckSuccesses[checkName]++
}

// GetHealthCheckFailures returns the number of failures for a health check.
//
// Parameters:
//   - checkName: The name of the health check
//
// Returns:
//   - int64: The number of failures
func (m *Metrics) GetHealthCheckFailures(checkName string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.healthCheckFailures[checkName]
}

// GetHealthCheckSuccesses returns the number of successes for a health check.
//
// Parameters:
//   - checkName: The name of the health check
//
// Returns:
//   - int64: The number of successes
func (m *Metrics) GetHealthCheckSuccesses(checkName string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.healthCheckSuccesses[checkName]
}

// GetHealthCheckFailureRate returns the failure rate for a health check.
//
// Parameters:
//   - checkName: The name of the health check
//
// Returns:
//   - float64: The failure rate (0.0 to 1.0)
func (m *Metrics) GetHealthCheckFailureRate(checkName string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	failures := m.healthCheckFailures[checkName]
	successes := m.healthCheckSuccesses[checkName]
	total := failures + successes

	if total == 0 {
		return 0
	}

	return float64(failures) / float64(total)
}

// Reset clears all collected metrics, resetting all counters and latency samples.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errorCounts = make(map[string]int64)
	m.errorRates = make(map[string]float64)
	m.retryCounts = make(map[string]int64)
	m.circuitOpenCount = 0
	m.latencies = make(map[string][]time.Duration)
	m.avgLatencies = make(map[string]time.Duration)
	m.p95Latencies = make(map[string]time.Duration)
	m.p99Latencies = make(map[string]time.Duration)
	m.circuitStateChanges = make(map[string]int64)
	m.circuitOpenDuration = make(map[string]time.Duration)
	m.healthCheckFailures = make(map[string]int64)
	m.healthCheckSuccesses = make(map[string]int64)
}

// MetricsSnapshot provides a point-in-time copy of all collected metrics.
// It captures error counts, retry counts, latency percentiles, circuit breaker
// state, and health check statistics at the time of capture.
type MetricsSnapshot struct {
	ErrorCounts          map[string]int64
	RetryCounts          map[string]int64
	CircuitOpenCount     int64
	AvgLatencies         map[string]time.Duration
	P95Latencies         map[string]time.Duration
	P99Latencies         map[string]time.Duration
	CircuitStateChanges  map[string]int64
	CircuitOpenDuration  map[string]time.Duration
	HealthCheckFailures  map[string]int64
	HealthCheckSuccesses map[string]int64
}

// Snapshot creates a point-in-time copy of all collected metrics.
// Returns a MetricsSnapshot containing all current error counts, retry counts,
// latency percentiles, circuit breaker statistics, and health check results.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgLatencies := make(map[string]time.Duration)
	for k, v := range m.avgLatencies {
		avgLatencies[k] = v
	}

	p95Latencies := make(map[string]time.Duration)
	for k, v := range m.p95Latencies {
		p95Latencies[k] = v
	}

	p99Latencies := make(map[string]time.Duration)
	for k, v := range m.p99Latencies {
		p99Latencies[k] = v
	}

	circuitStateChanges := make(map[string]int64)
	for k, v := range m.circuitStateChanges {
		circuitStateChanges[k] = v
	}

	circuitOpenDuration := make(map[string]time.Duration)
	for k, v := range m.circuitOpenDuration {
		circuitOpenDuration[k] = v
	}

	healthCheckFailures := make(map[string]int64)
	for k, v := range m.healthCheckFailures {
		healthCheckFailures[k] = v
	}

	healthCheckSuccesses := make(map[string]int64)
	for k, v := range m.healthCheckSuccesses {
		healthCheckSuccesses[k] = v
	}

	return MetricsSnapshot{
		ErrorCounts:          m.GetAllErrorCounts(),
		RetryCounts:          m.GetAllRetryCounts(),
		CircuitOpenCount:     m.circuitOpenCount,
		AvgLatencies:         avgLatencies,
		P95Latencies:         p95Latencies,
		P99Latencies:         p99Latencies,
		CircuitStateChanges:  circuitStateChanges,
		CircuitOpenDuration:  circuitOpenDuration,
		HealthCheckFailures:  healthCheckFailures,
		HealthCheckSuccesses: healthCheckSuccesses,
	}
}

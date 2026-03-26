package resilience

import (
	"testing"
	"time"
)

func TestMetrics_RecordError(t *testing.T) {
	m := NewMetrics()

	m.RecordError("test_error")
	m.RecordError("test_error")
	m.RecordError("another_error")

	count := m.GetErrorCount("test_error")
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}

	allCounts := m.GetAllErrorCounts()
	if allCounts["test_error"] != 2 {
		t.Fatalf("expected 2 for test_error, got %d", allCounts["test_error"])
	}
	if allCounts["another_error"] != 1 {
		t.Fatalf("expected 1 for another_error, got %d", allCounts["another_error"])
	}
}

func TestMetrics_RecordRetry(t *testing.T) {
	m := NewMetrics()

	m.RecordRetry("operation1")
	m.RecordRetry("operation1")
	m.RecordRetry("operation2")

	count := m.GetRetryCount("operation1")
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestMetrics_RecordLatency(t *testing.T) {
	m := NewMetrics()

	m.RecordLatency("test_op", 100*time.Millisecond)
	m.RecordLatency("test_op", 200*time.Millisecond)
	m.RecordLatency("test_op", 300*time.Millisecond)

	avg := m.GetAvgLatency("test_op")
	if avg != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %v", avg)
	}

	p95 := m.GetP95Latency("test_op")
	if p95 < 100*time.Millisecond || p95 > 300*time.Millisecond {
		t.Fatalf("expected p95 between 100ms and 300ms, got %v", p95)
	}

	p99 := m.GetP99Latency("test_op")
	if p99 < 100*time.Millisecond || p99 > 300*time.Millisecond {
		t.Fatalf("expected p99 between 100ms and 300ms, got %v", p99)
	}
}

func TestMetrics_RecordCircuitOpen(t *testing.T) {
	m := NewMetrics()

	m.RecordCircuitOpen("circuit1")
	m.RecordCircuitOpen("circuit1")
	m.RecordCircuitOpen("circuit2")

	count := m.GetCircuitOpenCount()
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestMetrics_RecordCircuitStateChange(t *testing.T) {
	m := NewMetrics()

	m.RecordCircuitStateChange("circuit1", StateClosed, StateOpen)
	m.RecordCircuitStateChange("circuit1", StateOpen, StateHalfOpen)

	changes := m.GetCircuitStateChanges("circuit1")
	if changes != 2 {
		t.Fatalf("expected 2 state changes, got %d", changes)
	}
}

func TestMetrics_RecordCircuitOpenDuration(t *testing.T) {
	m := NewMetrics()

	m.RecordCircuitOpenDuration("circuit1", 100*time.Millisecond)
	m.RecordCircuitOpenDuration("circuit1", 200*time.Millisecond)

	duration := m.GetCircuitOpenDuration("circuit1")
	if duration != 200*time.Millisecond {
		t.Fatalf("expected 200ms, got %v", duration)
	}
}

func TestMetrics_HealthCheck(t *testing.T) {
	m := NewMetrics()

	m.RecordHealthCheckSuccess("check1")
	m.RecordHealthCheckSuccess("check1")
	m.RecordHealthCheckFailure("check1")

	successes := m.GetHealthCheckSuccesses("check1")
	if successes != 2 {
		t.Fatalf("expected 2 successes, got %d", successes)
	}

	failures := m.GetHealthCheckFailures("check1")
	if failures != 1 {
		t.Fatalf("expected 1 failure, got %d", failures)
	}

	rate := m.GetHealthCheckFailureRate("check1")
	if rate < 0.3 || rate > 0.35 {
		t.Fatalf("expected failure rate ~0.333, got %f", rate)
	}
}

func TestMetrics_Reset(t *testing.T) {
	m := NewMetrics()

	m.RecordError("test_error")
	m.RecordRetry("test_op")
	m.RecordLatency("test_op", 100*time.Millisecond)
	m.RecordCircuitOpen("circuit1")

	m.Reset()

	if m.GetErrorCount("test_error") != 0 {
		t.Fatal("errors should be reset")
	}

	if m.GetRetryCount("test_op") != 0 {
		t.Fatal("retries should be reset")
	}

	if m.GetCircuitOpenCount() != 0 {
		t.Fatal("circuit open count should be reset")
	}
}

func TestMetrics_Snapshot(t *testing.T) {
	m := NewMetrics()

	m.RecordError("test_error")
	m.RecordRetry("test_op")
	m.RecordCircuitOpen("circuit1")

	snapshot := m.Snapshot()

	if snapshot.ErrorCounts["test_error"] != 1 {
		t.Fatalf("expected 1 error in snapshot, got %d", snapshot.ErrorCounts["test_error"])
	}

	if snapshot.RetryCounts["test_op"] != 1 {
		t.Fatalf("expected 1 retry in snapshot, got %d", snapshot.RetryCounts["test_op"])
	}

	if snapshot.CircuitOpenCount != 1 {
		t.Fatalf("expected 1 circuit open in snapshot, got %d", snapshot.CircuitOpenCount)
	}
}

func TestMetrics_MaxSamples(t *testing.T) {
	m := NewMetrics()

	for i := 0; i < 1500; i++ {
		m.RecordLatency("test_op", time.Duration(i)*time.Millisecond)
	}

	avg := m.GetAvgLatency("test_op")
	if avg == 0 {
		t.Fatal("expected average to be calculated")
	}
}

func TestMetrics_LatencyP95P99(t *testing.T) {
	m := NewMetrics()

	latencies := []time.Duration{
		10, 20, 30, 40, 50, 60, 70, 80, 90, 100,
		110, 120, 130, 140, 150, 160, 170, 180, 190, 200,
	}

	for _, lat := range latencies {
		m.RecordLatency("test_op", lat*time.Millisecond)
	}

	p95 := m.GetP95Latency("test_op")
	if p95 < 190*time.Millisecond || p95 > 200*time.Millisecond {
		t.Fatalf("expected p95 ~190-200ms, got %v", p95)
	}

	p99 := m.GetP99Latency("test_op")
	if p99 < 190*time.Millisecond || p99 > 200*time.Millisecond {
		t.Fatalf("expected p99 ~190-200ms, got %v", p99)
	}
}

func TestMetrics_EmptyOperation(t *testing.T) {
	m := NewMetrics()

	avg := m.GetAvgLatency("nonexistent")
	if avg != 0 {
		t.Fatalf("expected 0 for nonexistent operation, got %v", avg)
	}

	p95 := m.GetP95Latency("nonexistent")
	if p95 != 0 {
		t.Fatalf("expected 0 for nonexistent operation, got %v", p95)
	}

	p99 := m.GetP99Latency("nonexistent")
	if p99 != 0 {
		t.Fatalf("expected 0 for nonexistent operation, got %v", p99)
	}
}

func TestMetrics_HealthCheckEmpty(t *testing.T) {
	m := NewMetrics()

	rate := m.GetHealthCheckFailureRate("nonexistent")
	if rate != 0 {
		t.Fatalf("expected 0 for nonexistent health check, got %f", rate)
	}
}

func TestGlobalMetrics(t *testing.T) {
	GlobalMetrics.RecordError("global_test")

	count := GlobalMetrics.GetErrorCount("global_test")
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
}

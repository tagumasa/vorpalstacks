package resilience

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthChecker(t *testing.T) {
	hc := NewHealthChecker(nil)

	t.Run("IsHealthy returns true when no checks", func(t *testing.T) {
		assert.True(t, hc.IsHealthy())
	})

	t.Run("AddCheck registers health check", func(t *testing.T) {
		config := &HealthCheckConfig{
			Name:     "test-check",
			Check:    func(ctx context.Context) error { return nil },
			Timeout:  time.Second,
			Interval: time.Second,
		}
		hc.AddCheck(config)

		status := hc.GetStatus("test-check")
		assert.Equal(t, HealthUnknown, status)
	})

	t.Run("GetAllStatuses returns all statuses", func(t *testing.T) {
		config := &HealthCheckConfig{
			Name:     "test-check-2",
			Check:    func(ctx context.Context) error { return nil },
			Timeout:  time.Second,
			Interval: time.Second,
		}
		hc.AddCheck(config)

		statuses := hc.GetAllStatuses()
		assert.Contains(t, statuses, "test-check")
		assert.Contains(t, statuses, "test-check-2")
	})

	t.Run("GetDegradationLevel returns full by default", func(t *testing.T) {
		assert.Equal(t, LevelFull, hc.GetDegradationLevel())
	})

	t.Run("SetDegradationLevel updates level", func(t *testing.T) {
		hc.SetDegradationLevel(LevelMaintenance)
		assert.Equal(t, LevelMaintenance, hc.GetDegradationLevel())
	})
}

func TestDefaultHealthCheckConfig(t *testing.T) {
	config := DefaultHealthCheckConfig()
	assert.NotNil(t, config)

	assert.Equal(t, DefaultHealthCheckTimeout, config.Timeout)
	assert.Equal(t, DefaultHealthCheckInterval, config.Interval)
	assert.Equal(t, DefaultHealthCheckFailureThreshold, config.FailureThreshold)
}

func TestDegradationLevelString(t *testing.T) {
	assert.Equal(t, "Full", LevelFull.String())
	assert.Equal(t, "Partial", LevelPartial.String())
	assert.Equal(t, "Minimal", LevelMinimal.String())
	assert.Equal(t, "Maintenance", LevelMaintenance.String())
	assert.Equal(t, "Unknown", DegradationLevel(99).String())
}

func TestHealthStatusString(t *testing.T) {
	assert.Equal(t, "Healthy", HealthHealthy.String())
	assert.Equal(t, "Degraded", HealthDegraded.String())
	assert.Equal(t, "Unhealthy", HealthUnhealthy.String())
	assert.Equal(t, "Unknown", HealthUnknown.String())
	assert.Equal(t, "Unknown", HealthStatus(99).String())
}

func TestCheckDependency(t *testing.T) {
	hc := NewHealthChecker(nil)

	t.Run("Successful dependency check", func(t *testing.T) {
		err := hc.CheckDependency(context.Background(), "test-dep", func() error { return nil })
		assert.NoError(t, err)
		assert.Equal(t, HealthHealthy, hc.GetStatus("test-dep"))
	})

	t.Run("Failed dependency check", func(t *testing.T) {
		err := hc.CheckDependency(context.Background(), "test-dep-fail", func() error { return assert.AnError })
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test-dep-fail")
	})
}

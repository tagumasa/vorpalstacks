// Package resilience provides resilience patterns and utilities for vorpalstacks.
package resilience

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
)

// DegradationLevel represents the current degradation state of the system.
type DegradationLevel int

// DegradationLevel constants represent possible system degradation levels.
const (
	LevelFull DegradationLevel = iota
	LevelPartial
	LevelMinimal
	LevelMaintenance
)

// String returns the string representation of the degradation level.
func (dl DegradationLevel) String() string {
	switch dl {
	case LevelFull:
		return "Full"
	case LevelPartial:
		return "Partial"
	case LevelMinimal:
		return "Minimal"
	case LevelMaintenance:
		return "Maintenance"
	default:
		return "Unknown"
	}
}

// HealthStatus represents the health status of a component.
type HealthStatus int

// HealthStatus constants represent possible health statuses.
const (
	HealthHealthy HealthStatus = iota
	HealthDegraded
	HealthUnhealthy
	HealthUnknown
)

// String returns the string representation of the health status.
func (hs HealthStatus) String() string {
	switch hs {
	case HealthHealthy:
		return "Healthy"
	case HealthDegraded:
		return "Degraded"
	case HealthUnhealthy:
		return "Unhealthy"
	case HealthUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// HealthCheck is a function that performs a health check.
// It returns nil if the component is healthy, or an error otherwise.
type HealthCheck func(ctx context.Context) error

// HealthCheckConfig holds configuration for a health check.
type HealthCheckConfig struct {
	Name             string
	Check            HealthCheck
	Timeout          time.Duration
	Interval         time.Duration
	FailureThreshold int
}

// DefaultHealthCheckConfig returns the default health check configuration.
func DefaultHealthCheckConfig() *HealthCheckConfig {
	return &HealthCheckConfig{
		Timeout:          DefaultHealthCheckTimeout,
		Interval:         DefaultHealthCheckInterval,
		FailureThreshold: DefaultHealthCheckFailureThreshold,
	}
}

// HealthChecker monitors the health of various components.
type HealthChecker struct {
	configs        []*HealthCheckConfig
	statuses       map[string]HealthStatus
	failures       map[string]int
	degradation    DegradationLevel
	mu             sync.RWMutex
	logger         logs.Logger
	onStatusChange func(string, HealthStatus)
	metrics        *Metrics
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewHealthChecker creates a new health checker with the given logger.
func NewHealthChecker(logger logs.Logger) *HealthChecker {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthChecker{
		statuses:    make(map[string]HealthStatus),
		failures:    make(map[string]int),
		degradation: LevelFull,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// SetMetrics sets the metrics collector for the health checker.
func (hc *HealthChecker) SetMetrics(metrics *Metrics) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.metrics = metrics
}

// AddCheck adds a health check to the checker.
func (hc *HealthChecker) AddCheck(config *HealthCheckConfig) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.configs = append(hc.configs, config)
	hc.statuses[config.Name] = HealthUnknown
	hc.failures[config.Name] = 0

	if hc.logger != nil {
		hc.logger.Info("Health check added",
			logs.String("name", config.Name),
		)
	}
}

// OnStatusChange sets a callback function that is called when the health status changes.
func (hc *HealthChecker) OnStatusChange(callback func(string, HealthStatus)) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.onStatusChange = callback
}

// Start begins the health check loop, running checks at their configured intervals.
func (hc *HealthChecker) Start(ctx context.Context) {
	hc.mu.Lock()
	if hc.cancel != nil {
		hc.cancel()
	}
	if ctx != nil {
		hc.ctx, hc.cancel = context.WithCancel(ctx)
	} else {
		hc.ctx, hc.cancel = context.WithCancel(context.Background())
	}
	configs := make([]*HealthCheckConfig, len(hc.configs))
	copy(configs, hc.configs)
	runCtx := hc.ctx
	hc.mu.Unlock()

	for _, config := range configs {
		go hc.runCheck(runCtx, config)
	}
}

// Stop stops the health check loop.
func (hc *HealthChecker) Stop() {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	if hc.cancel != nil {
		hc.cancel()
	}
}

func (hc *HealthChecker) runCheck(ctx context.Context, config *HealthCheckConfig) {
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.checkHealth(ctx, config)
		}
	}
}

func (hc *HealthChecker) checkHealth(ctx context.Context, config *HealthCheckConfig) {
	checkCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	err := config.Check(checkCtx)

	hc.mu.Lock()
	defer hc.mu.Unlock()

	oldStatus := hc.statuses[config.Name]

	if err != nil {
		hc.failures[config.Name]++
		if hc.failures[config.Name] >= config.FailureThreshold {
			hc.statuses[config.Name] = HealthUnhealthy
		} else {
			hc.statuses[config.Name] = HealthDegraded
		}

		if hc.metrics != nil {
			hc.metrics.RecordHealthCheckFailure(config.Name)
		}

		if hc.logger != nil {
			hc.logger.Warn("Health check failed",
				logs.String("name", config.Name),
				logs.Int("failures", hc.failures[config.Name]),
				logs.Err(err),
			)
		}
	} else {
		hc.failures[config.Name] = 0
		hc.statuses[config.Name] = HealthHealthy

		if hc.metrics != nil {
			hc.metrics.RecordHealthCheckSuccess(config.Name)
		}

		if hc.logger != nil {
			hc.logger.Debug("Health check passed",
				logs.String("name", config.Name),
			)
		}
	}

	newStatus := hc.statuses[config.Name]
	if oldStatus != newStatus && hc.onStatusChange != nil {
		hc.onStatusChange(config.Name, newStatus)
	}

	hc.updateDegradationLocked()
}

func (hc *HealthChecker) updateDegradationLocked() {
	healthyCount := 0
	degradedCount := 0
	unhealthyCount := 0

	for _, status := range hc.statuses {
		switch status {
		case HealthHealthy:
			healthyCount++
		case HealthDegraded:
			degradedCount++
		case HealthUnhealthy:
			unhealthyCount++
		}
	}

	oldDegradation := hc.degradation

	if unhealthyCount > 0 {
		hc.degradation = LevelMaintenance
	} else if degradedCount > 0 {
		hc.degradation = LevelPartial
	} else {
		hc.degradation = LevelFull
	}

	if oldDegradation != hc.degradation && hc.logger != nil {
		hc.logger.Info("Degradation level changed",
			logs.String("from", oldDegradation.String()),
			logs.String("to", hc.degradation.String()),
		)
	}
}

// GetStatus returns the health status for a given check by name.
func (hc *HealthChecker) GetStatus(name string) HealthStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.statuses[name]
}

// GetAllStatuses returns all health statuses.
func (hc *HealthChecker) GetAllStatuses() map[string]HealthStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	result := make(map[string]HealthStatus)
	for k, v := range hc.statuses {
		result[k] = v
	}
	return result
}

// GetDegradationLevel returns the current degradation level.
func (hc *HealthChecker) GetDegradationLevel() DegradationLevel {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.degradation
}

// IsHealthy returns true if all health checks are healthy.
func (hc *HealthChecker) IsHealthy() bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	for _, status := range hc.statuses {
		if status != HealthHealthy {
			return false
		}
	}
	return true
}

// SetDegradationLevel manually sets the degradation level.
func (hc *HealthChecker) SetDegradationLevel(level DegradationLevel) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	oldLevel := hc.degradation
	hc.degradation = level

	if oldLevel != level && hc.logger != nil {
		hc.logger.Info("Degradation level set",
			logs.String("from", oldLevel.String()),
			logs.String("to", level.String()),
		)
	}
}

// CheckDependency checks a dependency and updates its health status based on the result.
func (hc *HealthChecker) CheckDependency(ctx context.Context, name string, check func() error) error {
	hc.mu.RLock()
	degraded := hc.degradation >= LevelMaintenance
	hc.mu.RUnlock()

	if degraded {
		return fmt.Errorf("service %s is in maintenance mode", name)
	}

	err := check()
	if err != nil {
		hc.mu.Lock()
		defer hc.mu.Unlock()
		hc.failures[name]++
		hc.statuses[name] = HealthDegraded
		if hc.metrics != nil {
			hc.metrics.RecordHealthCheckFailure(name)
		}
		if hc.logger != nil {
			hc.logger.Warn("Dependency check failed",
				logs.String("name", name),
				logs.Err(err),
			)
		}
		hc.updateDegradationLocked()
		return fmt.Errorf("dependency %s unavailable: %w", name, err)
	}

	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.failures[name] = 0
	hc.statuses[name] = HealthHealthy
	if hc.metrics != nil {
		hc.metrics.RecordHealthCheckSuccess(name)
	}
	hc.updateDegradationLocked()
	return nil
}

// HealthCheckResult represents the result of a health check.
type HealthCheckResult struct {
	Name      string
	Status    HealthStatus
	Timestamp time.Time
	Error     error
}

// CheckAll runs all health checks immediately and returns their results.
func (hc *HealthChecker) CheckAll(ctx context.Context) []*HealthCheckResult {
	hc.mu.RLock()
	configs := make([]*HealthCheckConfig, len(hc.configs))
	copy(configs, hc.configs)
	hc.mu.RUnlock()

	results := make([]*HealthCheckResult, 0, len(configs))

	for _, config := range configs {
		result := func() *HealthCheckResult {
			checkCtx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			err := config.Check(checkCtx)
			status := HealthHealthy

			if err != nil {
				status = HealthUnhealthy
			}

			return &HealthCheckResult{
				Name:      config.Name,
				Status:    status,
				Timestamp: time.Now(),
				Error:     err,
			}
		}()
		results = append(results, result)
	}

	return results
}

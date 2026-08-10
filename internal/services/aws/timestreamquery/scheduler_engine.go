package timestreamquery

import (
	"context"
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

var tsqTickerInterval = 1 * time.Minute

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		tsqTickerInterval = 1 * time.Second
	}
}

// ScheduledQueryEngine periodically evaluates schedule expressions for all
// ENABLED scheduled queries and triggers execution when the next run time
// has been reached. It reuses the shared cron/rate/at parser from
// internal/common/scheduleexpr.
type ScheduledQueryEngine struct {
	service *TimestreamQueryService

	lastFired sync.Map // key: region/name → time.Time

	running  bool
	runMu    sync.RWMutex
	stopChan chan struct{}
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewScheduledQueryEngine creates a new engine bound to the given service.
func NewScheduledQueryEngine(svc *TimestreamQueryService) *ScheduledQueryEngine {
	return &ScheduledQueryEngine{
		service:  svc,
		stopChan: make(chan struct{}),
	}
}

// Start launches the background goroutine.
func (e *ScheduledQueryEngine) Start() {
	e.runMu.Lock()
	defer e.runMu.Unlock()
	if e.running {
		return
	}
	e.running = true
	e.stopChan = make(chan struct{})
	e.ctx, e.cancel = context.WithCancel(context.Background())

	e.wg.Add(1)
	go e.run()

	logs.Debug("Timestream scheduled-query engine started")
}

// Stop signals the background goroutine to exit and waits for it.
func (e *ScheduledQueryEngine) Stop() {
	e.runMu.Lock()
	defer e.runMu.Unlock()
	if !e.running {
		return
	}
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	close(e.stopChan)
	e.wg.Wait()

	logs.Debug("Timestream scheduled-query engine stopped")
}

func (e *ScheduledQueryEngine) run() {
	defer e.wg.Done()
	defer func() { resilience.RecoverAndRestart("timestream-scheduled-query engine", &e.wg, e.run) }()

	ticker := time.NewTicker(tsqTickerInterval)
	defer ticker.Stop()

	e.checkScheduledQueries()
	e.cleanupOldRuns()

	for {
		select {
		case <-e.stopChan:
			return
		case <-ticker.C:
			e.checkScheduledQueries()
			e.cleanupOldRuns()
		}
	}
}

func (e *ScheduledQueryEngine) checkScheduledQueries() {
	if e.service.storageManager == nil {
		return
	}

	regions := e.service.storageManager.GetActiveRegions()
	now := time.Now().UTC()
	allActiveKeys := make(map[string]bool)

	for _, region := range regions {
		stores, err := e.service.getStoresForRegion(region)
		if err != nil {
			continue
		}

		queries, err := stores.scheduledQueryStore.ListScheduledQueries()
		if err != nil {
			logs.Debug("Failed to list scheduled queries", logs.String("region", region), logs.Err(err))
			continue
		}

		for _, sq := range queries {
			if sq.ScheduledQueryStatus != tsstore.ScheduledQueryStatusEnabled {
				continue
			}

			dedupKey := region + "/" + sq.Name
			allActiveKeys[dedupKey] = true

			if sq.ScheduleConfiguration == nil {
				continue
			}
			expr := sq.ScheduleConfiguration.ScheduleExpression
			if expr == "" {
				continue
			}

			nextTime, err := scheduleexpr.NextExecutionTime(expr, now, sq.CreationTime, nil)
			if err != nil {
				continue
			}

			if last, ok := e.lastFired.Load(dedupKey); ok {
				if lastTime, ok := last.(time.Time); ok && !lastTime.Before(nextTime) {
					continue
				}
			}

			diff := now.Sub(nextTime)
			if diff < 0 || diff >= time.Minute {
				continue
			}

			e.lastFired.Store(dedupKey, now)

			e.wg.Add(1)
			go func(scheduledQuery *tsstore.ScheduledQuery, st *tsQueryStores, sqRegion, key string) {
				defer e.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						logs.Error("panic executing scheduled query", logs.String("name", scheduledQuery.Name), logs.Any("panic", r))
						e.lastFired.Delete(key)
					}
				}()

				logs.Debug("Auto-triggering scheduled query",
					logs.String("name", scheduledQuery.Name),
					logs.String("region", sqRegion),
					logs.String("expression", expr))

				if _, err := e.service.executeScheduledQueryInternal(e.ctx, st, scheduledQuery, nextTime, tsstore.TriggerTypeAuto); err != nil {
					logs.Error("Auto-triggered scheduled query failed",
						logs.String("name", scheduledQuery.Name),
						logs.String("region", sqRegion),
						logs.Err(err))
					e.lastFired.Delete(key)
				} else {
					nextRun, nextErr := scheduleexpr.NextExecutionTime(expr, nextTime.Add(time.Second), scheduledQuery.CreationTime, nil)
					if nextErr == nil && nextRun.After(now) {
						_ = st.scheduledQueryStore.UpdateNextRunTime(scheduledQuery.Name, nextRun)
					}
				}
			}(sq, stores, region, dedupKey)
		}
	}

	e.lastFired.Range(func(key, _ any) bool {
		if k, ok := key.(string); ok {
			if !allActiveKeys[k] {
				e.lastFired.Delete(k)
			}
		}
		return true
	})
}

// cleanupOldRuns removes scheduled-query run records older than 48 hours
// to prevent unbounded accumulation.
func (e *ScheduledQueryEngine) cleanupOldRuns() {
	if e.service.storageManager == nil {
		return
	}

	cutoff := time.Now().UTC().Add(-48 * time.Hour)
	regions := e.service.storageManager.GetActiveRegions()

	for _, region := range regions {
		stores, err := e.service.getStoresForRegion(region)
		if err != nil {
			continue
		}

		queries, err := stores.scheduledQueryStore.ListScheduledQueries()
		if err != nil {
			continue
		}

		for _, sq := range queries {
			runs, err := stores.scheduledQueryRunStore.ListRuns(sq.ARN)
			if err != nil {
				continue
			}
			for _, run := range runs {
				if run.RunStatus == tsstore.ScheduleRunStatusSucceeded ||
					run.RunStatus == tsstore.ScheduleRunStatusFailed ||
					run.RunStatus == tsstore.ScheduleRunStatusCancelled {
					if !run.CompletionTime.IsZero() && run.CompletionTime.Before(cutoff) {
						_ = stores.scheduledQueryRunStore.DeleteRun(run.ARN)
					}
				}
			}
		}
	}
}

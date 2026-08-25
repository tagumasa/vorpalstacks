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

// shouldFireQuery resolves the latest elapsed execution boundary for the
// evaluation time now and reports whether it has not fired yet. The
// returned boundary is used as the execution's scheduled time. The
// boundary contract comes from the shared resolver: rate() never runs on
// the creation boundary (the first run is one full interval after
// creation), cron() recovers a matching minute missed by a slow ticker
// but never a boundary older than the creation, and at() fires once its
// timestamp is reached. previousRun is the record's persisted
// PreviousRunTime: a boundary at or before it has already been consumed
// by an earlier run, in this or a previous process lifetime.
func (e *ScheduledQueryEngine) shouldFireQuery(dedupKey, expr string, now, creation, previousRun time.Time) (time.Time, bool) {
	boundary, ok := scheduleexpr.ElapsedExecutionTime(expr, now, creation, nil, scheduleexpr.RateFiresAfterFirstInterval)
	if !ok {
		return time.Time{}, false
	}
	// The persisted last run is the source of truth across restarts (the
	// in-memory map dies with the process): PreviousRunTime is stamped
	// after every run — successful or failed, auto-triggered or manual —
	// with an invocation time at or after the boundary it consumed, so a
	// boundary at or before it has already run.
	if !previousRun.IsZero() && !previousRun.Before(boundary) {
		return time.Time{}, false
	}
	if last, ok := e.lastFired.Load(dedupKey); ok {
		if lastTime, ok := last.(time.Time); ok && !lastTime.Before(boundary) {
			return time.Time{}, false
		}
	}
	return boundary, true
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

			nextTime, due := e.shouldFireQuery(dedupKey, expr, now, sq.CreationTime, sq.PreviousRunTime)
			if !due {
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
						if err := st.scheduledQueryStore.UpdateNextRunTime(scheduledQuery.Name, nextRun); err != nil {
							logs.Warn("Failed to persist next run time for scheduled query", logs.String("name", scheduledQuery.Name), logs.Err(err))
						}
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
						if err := stores.scheduledQueryRunStore.DeleteRun(run.ARN); err != nil {
							logs.Warn("Failed to clean up old scheduled query run", logs.String("run", run.ARN), logs.Err(err))
						}
					}
				}
			}
		}
	}
}

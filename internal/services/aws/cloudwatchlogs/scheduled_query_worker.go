package cloudwatchlogs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/common/worker"
	logs "vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The scheduled-query worker: the ticker that claims due executions, the
// execution window arithmetic and the query-run lifecycle that carries one
// execution to its recorded delivery outcome. The operation surface the
// worker serves lives in scheduled_query_operations.go.

// scheduledExecutionWindow computes one execution's query window from
// the millisecond execution clock and the scheduled query's offsets.
// The offsets are documented in seconds — "The time offset in seconds
// that defines the lookback period for the query" (CreateScheduledQuery,
// startTimeOffset/endTimeOffset) — so they convert before the
// subtraction; without the conversion a documented one-hour lookback
// was a 3.6-second window.
func scheduledExecutionWindow(now int64, sq *logsstore.ScheduledQuery) (startTime, endTime int64) {
	startTime = now - 60*60*1000
	endTime = now
	if sq.StartTimeOffset > 0 {
		startTime = now - sq.StartTimeOffset*1000
	}
	if sq.EndTimeOffset > 0 {
		endTime = now - sq.EndTimeOffset*1000
	}
	return startTime, endTime
}

// scheduledQueryTickerInterval is the interval between scheduled-query
// evaluations. TEST_MODE shortens it so integration tests observe the
// AWS first-interval contract (the first run happens one full interval
// after creation) without also waiting out the evaluation phase — the
// same convention as the scheduler and timestream-query engines.
var scheduledQueryTickerInterval = worker.Cadence(1*time.Minute, time.Second)

// recordDelivery stamps the execution outcome on the stored record: the
// consumed boundary and the execution clock. A query deleted while its
// execution was in flight has no record left to stamp, which is benign.
func recordDelivery(store *logsstore.Store, id string, boundary, executedAt int64, status string) {
	if err := store.TouchScheduledQueryDelivery(id, boundary, executedAt, status); err != nil && !errors.Is(err, logsstore.ErrResourceNotFound) {
		logs.Error("Failed to record the scheduled query execution outcome",
			logs.String("scheduledQueryId", id),
			logs.Err(err))
	}
}

// startScheduledQueryWorker runs a background goroutine that evaluates
// schedule expressions and triggers enabled scheduled queries.
func (s *LogsService) startScheduledQueryWorker() {
	s.startRespawningWorker("scheduled query worker",
		worker.TickerLoop(s.ctx.Done(), scheduledQueryTickerInterval, s.tickScheduledQueries))
}

// scheduledQueryDue reports whether an ENABLED scheduled query should run
// at the evaluation time now. The schedule expression is evaluated in the
// query's configured timezone (UTC when unset or invalid), inside the
// optional scheduleStartTime/scheduleEndTime execution window. Each
// boundary runs exactly once: the boundary is the latest elapsed
// execution instant of the expression — rate() never runs on the creation
// boundary (the first run is one full interval after creation) and cron()
// recovers a matching minute missed between evaluations — and the query
// runs when that boundary is
// later than the last consumed boundary (LastExecutedBoundary, zero
// meaning never run). The marker holds the boundary value, never the
// execution clock: an execution that runs late must not suppress the
// next unexecuted boundary, and it survives restarts. The evaluated
// boundary is returned so the caller stamps exactly what it consumed.
func scheduledQueryDue(sq *logsstore.ScheduledQuery, now time.Time) (time.Time, bool) {
	nowMillis := now.UnixMilli()
	if sq.ScheduleStartTime != 0 && nowMillis < sq.ScheduleStartTime {
		return time.Time{}, false
	}
	if sq.ScheduleEndTime != 0 && nowMillis > sq.ScheduleEndTime {
		return time.Time{}, false
	}
	creationTime := time.UnixMilli(sq.CreationTime).UTC()
	boundary, elapsed := scheduleexpr.ElapsedExecutionTime(sq.ScheduleExpression, now.In(scheduledQueryLocation(sq)), creationTime, nil, scheduleexpr.RateFiresAfterFirstInterval)
	if !elapsed {
		return time.Time{}, false
	}
	if !boundary.After(time.UnixMilli(sq.LastExecutedBoundary).UTC()) {
		return time.Time{}, false
	}
	return boundary, true
}

// scheduledQueryLocation resolves the timezone the schedule expression is
// evaluated in: the query's configured timezone, or UTC when unset or
// invalid.
func scheduledQueryLocation(sq *logsstore.ScheduledQuery) *time.Location {
	if sq.Timezone != "" {
		if loc, err := time.LoadLocation(sq.Timezone); err == nil {
			return loc
		}
		logs.Debug("Invalid scheduled query timezone, falling back to UTC",
			logs.String("scheduledQuery", sq.Name),
			logs.String("timezone", sq.Timezone))
	}
	return time.UTC
}

func (s *LogsService) tickScheduledQueries() {
	now := time.Now().UTC()

	// The region set comes from the storage manager, not the store
	// instance cache: an ENABLED scheduled query in a region nobody has
	// touched since restart must still fire, and its schedule boundaries
	// must not be skipped silently.
	for _, region := range worker.ActiveRegions(s.storageManager) {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			logs.Error("Failed to resolve logs store for scheduled query evaluation",
				logs.String("region", region), logs.Err(err))
			continue
		}

		queries, err := store.ListScheduledQueries(logsstore.ScheduledQueryStateEnabled)
		if err != nil {
			continue
		}

		for _, sq := range queries {
			// The boundary the evaluation consumes is the exact value
			// the execution stamps; recomputing it later against the
			// execution clock could advance past an unexecuted
			// boundary and suppress it.
			boundary, due := scheduledQueryDue(sq, now)
			if !due {
				continue
			}
			// A still-running execution owns the boundary: it advances
			// the marker only when it finishes, so re-firing now would
			// duplicate the run.
			if scheduledExecutionInFlight(store, sq.Id) {
				continue
			}
			// The boundary is claimed on the ticker goroutine itself —
			// the RUNNING record exists before the task is spawned, so
			// a slow spawn cannot leave a window the next tick's
			// in-flight check reads as idle.
			now, exec := claimScheduledExecution(store, sq)
			if exec == nil {
				continue
			}
			// The execution runs as its own task, off the ticker
			// goroutine: one slow query (or a hanging delivery) must
			// not head-of-line block every scheduled query in every
			// region.
			s.spawnTask(func(taskCtx context.Context) {
				s.runScheduledQueryExecution(taskCtx, region, store, sq, exec, now, boundary)
			})
		}
	}
}

// scheduledExecutionInFlight reports whether the scheduled query has an
// execution record still marked RUNNING — either genuinely running in
// this process or orphaned by a restart until reconciliation fails it.
func scheduledExecutionInFlight(store *logsstore.Store, sqId string) bool {
	execs, err := store.ListScheduledQueryExecutions(sqId, 0, 0)
	if err != nil {
		return false
	}
	for _, exec := range execs {
		if exec.Status == logsstore.ScheduledExecutionStatusRunning {
			return true
		}
	}
	return false
}

// scheduledExecutionCancelled reports whether the execution record was
// marked CANCELLED by StopQuery while the worker was running.
func scheduledExecutionCancelled(store *logsstore.Store, sqId, queryId string) bool {
	execs, err := store.ListScheduledQueryExecutions(sqId, 0, 0)
	if err != nil {
		return false
	}
	for _, exec := range execs {
		if exec.QueryId == queryId {
			return exec.Status == logsstore.ScheduledExecutionStatusCancelled
		}
	}
	return false
}

// claimScheduledExecution writes the execution's RUNNING record on the
// caller's goroutine: the boundary is claimed before the execution task
// is even spawned, so the next ticker pass observes the execution in
// flight instead of re-firing a boundary whose task has not started
// yet. A write failure reports a nil execution and the boundary stays
// unclaimed (the next pass retries it).
func claimScheduledExecution(store *logsstore.Store, sq *logsstore.ScheduledQuery) (int64, *logsstore.ScheduledQueryExecution) {
	now := time.Now().UTC().UnixMilli()

	exec := &logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          fmt.Sprintf("sq-%s-%d", sq.Id, now),
		TriggerTime:      now,
		Status:           logsstore.ScheduledExecutionStatusRunning,
	}
	if err := store.PutScheduledQueryExecution(exec); err != nil {
		logs.Error("Failed to persist scheduled query execution (RUNNING)",
			logs.String("scheduledQueryId", sq.Id),
			logs.Err(err))
		return 0, nil
	}
	return now, exec
}

// triggerScheduledQuery claims the boundary and runs the execution on
// the calling goroutine (the direct-delivery entry point).
func (s *LogsService) triggerScheduledQuery(ctx context.Context, region string, store *logsstore.Store, sq *logsstore.ScheduledQuery, boundary time.Time) {
	now, exec := claimScheduledExecution(store, sq)
	if exec == nil {
		return
	}
	s.runScheduledQueryExecution(ctx, region, store, sq, exec, now, boundary)
}

func (s *LogsService) runScheduledQueryExecution(ctx context.Context, region string, store *logsstore.Store, sq *logsstore.ScheduledQuery, exec *logsstore.ScheduledQueryExecution, now int64, boundary time.Time) {

	defer func() {
		if r := recover(); r != nil {
			exec.Status = logsstore.ScheduledExecutionStatusFailed
			exec.ErrorMessage = fmt.Sprintf("panic: %v", r)
			if err := store.FinaliseScheduledQueryExecution(exec); err != nil {
				logs.Error("Failed to persist scheduled query execution (FAILED after panic)",
					logs.String("scheduledQueryId", sq.Id),
					logs.Err(err))
			}
			// The claimed boundary is consumed on the panic path as
			// well: a terminal failure advances the schedule instead of
			// re-firing the boundary on every worker tick.
			recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusFailed)
		}
	}()

	startTime, endTime := scheduledExecutionWindow(now, sq)

	fail := func(message string) {
		exec.Status = logsstore.ScheduledExecutionStatusFailed
		exec.ErrorMessage = message
		// The gated finaliser: a StopQuery that flipped the record
		// CANCELLED while the query was failing keeps its status — the
		// terminal-wins rule the success and delivery paths already
		// enforce.
		if err := store.FinaliseScheduledQueryExecution(exec); err != nil {
			logs.Error("Failed to persist scheduled query execution (FAILED)",
				logs.String("scheduledQueryId", sq.Id),
				logs.Err(err))
		}
		// The claimed boundary is consumed here too: a persistently
		// failing query retries on the next schedule boundary instead
		// of re-firing the consumed boundary on every worker tick —
		// the same rule the cancellation and delivery-failure paths
		// already enforce.
		recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusFailed)
	}

	// The stored identifiers echo back verbatim; execution resolves the
	// name-or-ARN elements to store name keys.
	groups := resolveLogGroupIdentifiers(sq.LogGroupIdentifiers)
	events := fetchGroupEventsForQuery(store, groups, startTime, endTime)
	execCtx := &execContext{
		startTime:     startTime,
		endTime:       endTime,
		accountID:     s.accountID,
		defaultGroups: groups,
		events:        events,
		fetchEvents: func(groups []string, start, end int64) ([]logEventWithContext, error) {
			return fetchGroupEventsForQuery(store, groups, start, end), nil
		},
		listLogGroups: func() ([]sourceGroupInfo, error) {
			return listSourceGroups(store), nil
		},
		getLookupTable: func(name string) (*parsedLookupTable, error) {
			return s.loadParsedLookupTable(store, region, name)
		},
		subqueryCache: map[string][]interface{}{},
	}
	// The scheduled plane runs under the same documented outer runtime
	// as the interactive one ("Queries time out after 60 minutes of
	// runtime"): without it a hung execution holds a RUNNING record —
	// and its schedule's next boundary — for the process lifetime.
	execCtx.startedMs = time.Now().UnixMilli()
	execCtx.deadline = queryDeadlineNow().Add(queryRuntimeLimit)
	rows, err := executeQueryContext(execCtx, sq.QueryString)
	if err != nil {
		fail(fmt.Sprintf("query failed: %v", err))
		return
	}
	// A breached deadline stops the pipeline without an error; the
	// execution surfaces the enum's Timeout member — the scheduled
	// record's sibling of the query plane's Timeout status — and the
	// claimed boundary is consumed so the schedule advances.
	if execCtx.timedOut {
		exec.Status = logsstore.ScheduledExecutionStatusTimeout
		exec.ErrorMessage = "query execution exceeded the 60-minute runtime limit"
		if err := store.FinaliseScheduledQueryExecution(exec); err != nil {
			logs.Error("Failed to persist scheduled query execution (TIMEOUT)",
				logs.String("scheduledQueryId", sq.Id),
				logs.Err(err))
		}
		recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusTimeout)
		return
	}
	// Cancellation checkpoint before delivery: StopQuery marks a running
	// execution's record CANCELLED, and the terminal status wins — no
	// destination receives the results and the record keeps the
	// canceller's status. The boundary is still recorded so the schedule
	// advances instead of re-firing the consumed boundary. A cancelled
	// context (service shutdown) is the same terminal turn: the record
	// must reach a terminal state here too, because the boundary it
	// claimed is already consumed — a record left RUNNING pairs a
	// consumed boundary with an in-flight guard that then suppresses
	// every later evaluation until a restart reconciles it.
	if scheduledExecutionCancelled(store, sq.Id, exec.QueryId) {
		recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusFailed)
		return
	}
	if ctx != nil && ctx.Err() != nil {
		exec.Status = logsstore.ScheduledExecutionStatusFailed
		exec.ErrorMessage = "interrupted by service shutdown"
		if err := store.FinaliseScheduledQueryExecution(exec); err != nil {
			logs.Error("Failed to persist scheduled query execution (FAILED on shutdown)",
				logs.String("scheduledQueryId", sq.Id),
				logs.Err(err))
		}
		recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusFailed)
		return
	}
	exec.Destinations = s.deliverScheduledQueryResults(ctx, region, store, sq, exec.QueryId, rows)
	for _, dest := range exec.Destinations {
		if dest.Status != destinationStatusComplete {
			exec.Status = logsstore.ScheduledExecutionStatusFailed
			exec.ErrorMessage = fmt.Sprintf("destination delivery failed: %s", dest.ErrorMessage)
			break
		}
	}
	if exec.Status == logsstore.ScheduledExecutionStatusFailed {
		// The finalise gate re-reads the record under the record mutex:
		// a StopQuery that landed after the checkpoint keeps its
		// CANCELLED status against this terminal write.
		if err := store.FinaliseScheduledQueryExecution(exec); err != nil {
			logs.Error("Failed to persist scheduled query execution (delivery FAILED)",
				logs.String("scheduledQueryId", sq.Id),
				logs.Err(err))
		}
		// Record the trigger so a persistently failing destination retries
		// on schedule instead of on every worker tick.
		recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusFailed)
		return
	}

	stats := queryStats{
		recordsScanned: int64(len(execCtx.events)),
	}
	for _, e := range execCtx.events {
		stats.bytesScanned += int64(len(e.message))
	}
	stats.recordsMatched = execCtx.recordsMatched

	exec.Status = logsstore.ScheduledExecutionStatusSuccess
	exec.RecordsScanned = stats.recordsScanned
	exec.RecordsMatched = stats.recordsMatched
	if err := store.FinaliseScheduledQueryExecution(exec); err != nil {
		logs.Error("Failed to persist scheduled query execution (SUCCESS)",
			logs.String("scheduledQueryId", sq.Id),
			logs.Err(err))
	}

	recordDelivery(store, sq.Id, boundary.UnixMilli(), now, logsstore.ScheduledQueryStatusComplete)
}

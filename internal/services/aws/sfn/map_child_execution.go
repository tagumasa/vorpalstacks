package sfn

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// This file implements the Distributed Map child-execution dispatch: in
// Distributed mode every work unit (a single item, or a batch when an
// ItemBatcher is configured) runs as its own child workflow execution
// under the parent state machine's namespace, with its own ARN, history
// and status. ListExecutions surfaces the children through the mapRunArn
// filter (ListExecutions API contract).

// mapChildMeta carries the identity of a dispatched child execution for
// the ResultWriter export records.
type mapChildMeta struct {
	Arn          string
	Name         string
	RedriveCount int64
}

// isDistributedMap reports whether the Map state runs its units as child
// workflow executions.
func isDistributedMap(state *sfnstore.MapState) bool {
	return state.ItemProcessor != nil &&
		state.ItemProcessor.ProcessorConfig != nil &&
		state.ItemProcessor.ProcessorConfig.Mode == "DISTRIBUTED"
}

// childLabel returns the name component a Map state stamps on its child
// executions: the Label when set, the state name otherwise.
func childLabel(state *sfnstore.MapState, currentState string) string {
	if state.Label != "" {
		return state.Label
	}
	return currentState
}

// mapChildName derives the deterministic child-execution name for a work
// unit of a Map Run: the same unit of the same run reclaims the same child
// across redrives, which is what makes the child's identity stable. A run
// of a retried state carries its attempt in the name — "When you retry a
// Map state, it creates a new Map Run" and the retry "applies to all of
// the child workflow executions", so each attempt's children are fresh
// executions disjoint from every prior attempt's, succeeded ones included.
func mapChildName(execCtx *ExecutionContext, state *sfnstore.MapState, ordinal int, attempt int64) string {
	name := fmt.Sprintf("%s:%s-%d", execCtx.Execution.Name, childLabel(state, execCtx.CurrentState), ordinal)
	if attempt > 0 {
		name = fmt.Sprintf("%s-r%d", name, attempt)
	}
	return name
}

// mapChildARN derives the child-execution ARN for a work unit of a Map
// Run. The region and account come from the parent execution's own ARN,
// falling back to the executor's view.
func (e *Executor) mapChildARN(execCtx *ExecutionContext, state *sfnstore.MapState, ordinal int, attempt int64) string {
	childName := mapChildName(execCtx, state, ordinal, attempt)
	region, accountID := e.region, e.accountID
	if parsed, err := svcarn.ParseARN(execCtx.Execution.ExecutionArn); err == nil {
		region, accountID = parsed.Region, parsed.AccountID
	}
	smName := svcarn.ExtractStateMachineNameFromARN(execCtx.Execution.ExecutionArn)
	return svcarn.NewARNBuilder(accountID, region).StepFunctions().Execution(smName, childName)
}

// beginMapChildExecution creates (or reclaims after a redrive) the child
// workflow execution record for one work unit and writes its
// ExecutionStarted event. attempt is the owning Map Run's state-level retry
// attempt — it stamps the child's name, so only a redrive of the same run
// can collide with an existing record. The returned record backs the unit's
// ExecutionContext so the whole iteration history lands on the child. For
// a reclaimed child the returned resume point continues the child from its
// own failed state — its history is a complete execution record, so the
// same analysis as a top-level redrive applies and its succeeded inner
// states do not re-run; it is nil for a fresh child. A non-nil error means
// the child record could not be established: the caller must fail the
// unit, never fall back to inline execution (the unit's topology is
// distributed by definition).
func (e *Executor) beginMapChildExecution(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, mapRunArn string, unit mapWorkUnit, ordinal int, iteratorDef *sfnstore.StateMachineDefinition, attempt int64) (*sfnstore.Execution, int64, *resumePoint, error) {
	childArn := e.mapChildARN(execCtx, state, ordinal, attempt)

	now := time.Now().UTC()
	child := &sfnstore.Execution{
		ExecutionArn:           childArn,
		StateMachineArn:        execCtx.Execution.StateMachineArn,
		StateMachineVersionArn: execCtx.Execution.StateMachineVersionArn,
		StateMachineAliasArn:   execCtx.Execution.StateMachineAliasArn,
		Name:                   mapChildName(execCtx, state, ordinal, attempt),
		Status:                 "RUNNING",
		Input:                  unit.InputJSON,
		StartDate:              now,
		MapRunArn:              mapRunArn,
		ItemCount:              int64(unit.ItemCount),
	}
	if err := e.store.CreateExecution(ctx, child); err != nil {
		if !errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil, 0, nil, err
		}
		// A redrive re-runs the unit: reclaim the prior child record and
		// count the re-run in its redrive counters. The resume point is
		// determined against the pre-transition record first, so a failure
		// here leaves the child untouched in its terminal state and fails
		// the unit instead of leaving a half-redriven child behind.
		var childResume *resumePoint
		rp, rerr := determineResumePoint(ctx, e.store, childArn, iteratorDef)
		if rerr != nil {
			return nil, 0, nil, fmt.Errorf("determine child resume point: %w", rerr)
		}
		childResume = rp

		// The re-run appends past the prior attempt's events: history keys
		// embed the event id, so restarting the child-local sequence at one
		// would overwrite the earlier attempt's events and leave its
		// terminal event stranded mid-sequence. The redrive marker takes
		// the next id and the returned counter continues from it.
		lastId, lerr := e.store.LastExecutionHistoryId(ctx, childArn)
		if lerr != nil {
			return nil, 0, nil, lerr
		}

		// The reclaim is a serialised transition: the child flips to
		// RUNNING — persisted immediately, not at unit completion — only
		// while its fresh status is still one the Map Run may reclaim
		// (an unsuccessful terminal status, or PENDING_REDRIVE for a child
		// parked at the concurrency limit).
		existing, terr := e.store.TransitionExecutionForRedrive(ctx, childArn,
			func(fresh *sfnstore.Execution) bool {
				return isMapRunReclaimStatus(fresh.Status)
			},
			func(fresh *sfnstore.Execution) {
				fresh.Status = "RUNNING"
				fresh.Output = ""
				fresh.Error = ""
				fresh.Cause = ""
				fresh.StopDate = time.Time{}
				fresh.Input = unit.InputJSON
				fresh.ItemCount = int64(unit.ItemCount)
				fresh.MapRunArn = mapRunArn
				fresh.RedriveCount++
				fresh.RedriveDate = now
			})
		if terr != nil {
			return nil, 0, nil, terr
		}
		child = existing

		if err := e.addExecutionHistoryEvent(ctx, child, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:    child.ExecutionArn,
			EventId:         lastId + 1,
			PreviousEventId: lastId,
			Type:            "ExecutionRedriven",
			Timestamp:       now,
			ExecutionRedrivenEventDetails: &sfnstore.ExecutionRedrivenEventDetails{
				RedriveCount: existing.RedriveCount,
			},
		}); err != nil {
			logs.Warn("sfn: failed to record child execution redrive", logs.String("arn", child.ExecutionArn), logs.Err(err))
		}
		return child, lastId + 1, childResume, nil
	}

	roleArn := ""
	if e.currentStateMachine != nil {
		roleArn = e.currentStateMachine.RoleArn
	}
	if err := e.addExecutionHistoryEvent(ctx, child, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: child.ExecutionArn,
		EventId:      1,
		Type:         "ExecutionStarted",
		Timestamp:    now,
		ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{
			Input:                  unit.InputJSON,
			RoleArn:                roleArn,
			StateMachineAliasArn:   child.StateMachineAliasArn,
			StateMachineVersionArn: child.StateMachineVersionArn,
		},
	}); err != nil {
		logs.Warn("sfn: failed to record child execution start", logs.String("arn", child.ExecutionArn), logs.Err(err))
	}

	return child, 1, nil, nil
}

// parkMapChildPendingRedrive persists the prior child execution of a
// redriven unit as PENDING_REDRIVE while the unit waits for a Map Run
// concurrency slot — the child is otherwise indistinguishable from one
// still holding its terminal status. attempt is the owning run's retry
// attempt, stamping the child name the same way the dispatch does. Only
// children in an unsuccessful terminal status are parked: a unit the prior
// attempt never dispatched has no child record to surface, and a parked
// child is never parked twice.
func (e *Executor) parkMapChildPendingRedrive(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, ordinal int, attempt int64) {
	childArn := e.mapChildARN(execCtx, state, ordinal, attempt)
	child, err := e.store.GetExecution(ctx, childArn)
	if err != nil || child == nil {
		return
	}
	if !isRedrivableStatus(child.Status) {
		return
	}
	child.Status = "PENDING_REDRIVE"
	if err := e.store.UpdateExecution(ctx, child); err != nil {
		logs.Warn("sfn: failed to park child execution as pending redrive", logs.String("arn", childArn), logs.Err(err))
	}
}

// finishMapChildExecution writes the terminal event of a dispatched child
// execution and persists its final status. A cancelled child terminates
// ABORTED with the platform's abort vocabulary, mirroring the parent
// execution's StopExecution terminal pair.
func (e *Executor) finishMapChildExecution(ctx context.Context, child *sfnstore.Execution, eventID *int64, output string, execErr *ExecutionError, aborted bool) {
	now := time.Now().UTC()
	child.StopDate = now
	switch {
	case execErr == nil && !aborted:
		child.Status = "SUCCEEDED"
		child.Output = output
		if err := e.addExecutionHistoryEvent(ctx, child, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: child.ExecutionArn,
			EventId:      nextEventID(eventID),
			Type:         "ExecutionSucceeded",
			Timestamp:    now,
			ExecutionSucceededEventDetails: &sfnstore.ExecutionSucceededEventDetails{
				Output: output,
			},
		}); err != nil {
			logs.Warn("sfn: failed to record child execution success", logs.String("arn", child.ExecutionArn), logs.Err(err))
		}
	case aborted:
		child.Status = "ABORTED"
		child.Error = "ExecutionAborted"
		child.Cause = "Execution was aborted by StopExecution"
		if err := e.addExecutionHistoryEvent(ctx, child, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: child.ExecutionArn,
			EventId:      nextEventID(eventID),
			Type:         "ExecutionAborted",
			Timestamp:    now,
			ExecutionAbortedEventDetails: &sfnstore.ExecutionAbortedEventDetails{
				Error: child.Error,
				Cause: child.Cause,
			},
		}); err != nil {
			logs.Warn("sfn: failed to record child execution abort", logs.String("arn", child.ExecutionArn), logs.Err(err))
		}
	default:
		child.Status = "FAILED"
		child.Error = execErr.ErrorCode
		child.Cause = execErr.Cause
		if err := e.addExecutionHistoryEvent(ctx, child, &sfnstore.ExecutionHistoryEvent{
			ExecutionArn: child.ExecutionArn,
			EventId:      nextEventID(eventID),
			Type:         "ExecutionFailed",
			Timestamp:    now,
			ExecutionFailedEventDetails: &sfnstore.ExecutionFailedEventDetails{
				Error: child.Error,
				Cause: child.Cause,
			},
		}); err != nil {
			logs.Warn("sfn: failed to record child execution failure", logs.String("arn", child.ExecutionArn), logs.Err(err))
		}
	}
	if err := e.store.UpdateExecution(ctx, child); err != nil {
		// The Map Run rolls the unit result up regardless; the stale
		// record only costs the child's DescribeExecution freshness.
		logs.Warn("sfn: failed to persist child execution", logs.String("arn", child.ExecutionArn), logs.Err(err))
	}
}

// nextEventID advances a child-local event counter atomically.
func nextEventID(counter *int64) int64 {
	if counter == nil {
		return 1
	}
	return atomic.AddInt64(counter, 1)
}

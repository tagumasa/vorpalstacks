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

// beginMapChildExecution creates (or reclaims after a redrive) the child
// workflow execution record for one work unit and writes its
// ExecutionStarted event. The returned record backs the unit's
// ExecutionContext so the whole iteration history lands on the child.
func (e *Executor) beginMapChildExecution(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, mapRunArn string, unit mapWorkUnit, ordinal int) *sfnstore.Execution {
	label := state.Label
	if label == "" {
		label = execCtx.CurrentState
	}
	childName := fmt.Sprintf("%s:%s-%d", execCtx.Execution.Name, label, ordinal)

	region, accountID := e.region, e.accountID
	smName := svcarn.ExtractStateMachineNameFromARN(execCtx.Execution.ExecutionArn)
	if parsed, err := svcarn.ParseARN(execCtx.Execution.ExecutionArn); err == nil {
		region, accountID = parsed.Region, parsed.AccountID
	}
	childArn := svcarn.NewARNBuilder(accountID, region).StepFunctions().Execution(smName, childName)

	now := time.Now().UTC()
	child := &sfnstore.Execution{
		ExecutionArn:           childArn,
		StateMachineArn:        execCtx.Execution.StateMachineArn,
		StateMachineVersionArn: execCtx.Execution.StateMachineVersionArn,
		StateMachineAliasArn:   execCtx.Execution.StateMachineAliasArn,
		Name:                   childName,
		Status:                 "RUNNING",
		Input:                  unit.InputJSON,
		StartDate:              now,
		MapRunArn:              mapRunArn,
		ItemCount:              int64(unit.ItemCount),
	}
	if err := e.store.CreateExecution(ctx, child); err != nil {
		if !errors.Is(err, sfnstore.ErrExecutionAlreadyExists) {
			return nil
		}
		// A redrive re-runs the unit: reclaim the prior child record and
		// count the re-run in its redrive counters.
		existing, gerr := e.store.GetExecution(ctx, childArn)
		if gerr != nil || existing == nil {
			return nil
		}
		existing.Status = "RUNNING"
		existing.Output = ""
		existing.Error = ""
		existing.Cause = ""
		existing.StopDate = time.Time{}
		existing.Input = unit.InputJSON
		existing.ItemCount = int64(unit.ItemCount)
		existing.MapRunArn = mapRunArn
		existing.RedriveCount++
		existing.RedriveDate = now
		child = existing
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
			Input:           unit.InputJSON,
			RoleArn:         roleArn,
			StateMachineArn: child.StateMachineArn,
			Name:            child.Name,
		},
	}); err != nil {
		logs.Warn("sfn: failed to record child execution start", logs.String("arn", child.ExecutionArn), logs.Err(err))
	}

	return child
}

// finishMapChildExecution writes the terminal event of a dispatched child
// execution and persists its final status.
func (e *Executor) finishMapChildExecution(ctx context.Context, child *sfnstore.Execution, eventID *int64, output string, execErr *ExecutionError) {
	now := time.Now().UTC()
	child.StopDate = now
	if execErr == nil {
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
	} else {
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

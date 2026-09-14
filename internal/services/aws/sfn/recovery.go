package sfn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// RecoverRunningExecutions resumes every execution that was still RUNNING
// when the server stopped. Executions run as in-memory goroutines, so a
// restart would otherwise leave them RUNNING in the store forever with no
// goroutine behind them — permanent zombies that redrive cannot address
// (it is restricted to FAILED/TIMED_OUT/ABORTED executions). Recovery
// reuses the redrive resume machinery: each execution continues from the
// first state it entered but never exited, with already-succeeded states
// not re-executed. Task timeout and heartbeat timers restart from the
// recovery point — a restart-tolerant approximation, since AWS never
// loses an in-flight execution.
//
// Must run once at boot before the server accepts traffic, so that no
// live goroutine exists for the swept executions.
func (s *StepFunctionService) RecoverRunningExecutions() {
	ctx := context.Background()
	recovered := 0
	for _, region := range s.storageManager.GetActiveRegions() {
		store, err := s.getStoreForRegion(region)
		if err != nil {
			logs.Error("sfn: failed to open store for execution recovery", logs.String("region", region), logs.Err(err))
			continue
		}
		recovered += s.recoverRegionExecutions(ctx, region, store)
	}
	if recovered > 0 {
		logs.Info("sfn: recovered running executions after restart", logs.Int("count", recovered))
	}
}

// recoverRegionExecutions sweeps the region's executions and resumes the
// RUNNING ones. The sweep lists executions without a state-machine filter
// so orphans whose state machine record was deleted are still found (and
// failed) rather than staying invisible; the full pagination chain is
// followed.
func (s *StepFunctionService) recoverRegionExecutions(ctx context.Context, region string, store *sfnstore.StepFunctionStore) int {
	recovered := 0
	nextToken := ""
	for {
		result, err := store.ListExecutions(ctx, "", "RUNNING", "", "", sfnstore.MaxPageSize, nextToken)
		if err != nil {
			logs.Error("sfn: failed to list running executions for recovery", logs.String("region", region), logs.Err(err))
			return recovered
		}
		for _, exec := range result.Executions {
			if s.resumeRecoveredExecution(ctx, region, store, exec) {
				recovered++
			}
		}
		if result.NextToken == "" {
			return recovered
		}
		nextToken = result.NextToken
	}
}

// resumeRecoveredExecution resumes a single RUNNING execution found at
// boot. RedriveCount and RedriveDate stay untouched: they count
// user-initiated redrives. An execution whose state machine no longer
// exists or whose definition no longer parses is failed with
// States.Runtime instead of being left as a zombie. Distributed Map child
// executions are never resumed here: their unit of work belongs to the
// parent execution, which the sweep resumes (or fails) itself — a child
// resumed independently would run beside the re-dispatching parent and
// duplicate every item.
func (s *StepFunctionService) resumeRecoveredExecution(ctx context.Context, region string, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) bool {
	if exec.MapRunArn != "" {
		return s.deferMapChildToParent(ctx, store, exec)
	}

	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrStateMachineNotFound) {
			s.failUnrecoverableExecution(ctx, store, exec, fmt.Sprintf("state machine %s no longer exists", exec.StateMachineArn))
			return false
		}
		s.failUnrecoverableExecution(ctx, store, exec, "its state machine record could not be read: "+err.Error())
		return false
	}
	definition, err := parseStateMachineDefinition(sm.Definition)
	if err != nil {
		s.failUnrecoverableExecution(ctx, store, exec, "state machine definition no longer parses")
		return false
	}
	rp, err := determineResumePoint(ctx, store, exec.ExecutionArn, definition)
	if err != nil {
		s.failUnrecoverableExecution(ctx, store, exec, "failed to determine the resume point from the execution history")
		return false
	}

	executionArn := exec.ExecutionArn
	executor := NewExecutorWithStores(store, s.bus, s.accountID, region, s.taskCredentialsAuthz)
	s.launchExecutionGoroutine(store, exec, "recovered execution", func(ctx context.Context) error {
		// Recovery keeps a fresh timeout window: the lost goroutine's
		// elapsed time is unknowable after a restart.
		return executor.ExecuteStateMachineFromState(ctx, exec, rp.StateName, rp.Input, rp.LastEventId, true)
	})

	logs.Info("sfn: resumed execution after restart",
		logs.String("arn", executionArn),
		logs.String("stateMachineArn", exec.StateMachineArn),
		logs.String("state", rp.StateName))
	return true
}

// deferMapChildToParent decides the recovery ownership of a Distributed
// Map child: the parent's own resume re-dispatches unfinished units, so a
// child whose parent is still RUNNING is left entirely to it. A child
// whose Map Run or parent is gone has no dispatcher and is failed instead
// of being left a zombie.
func (s *StepFunctionService) deferMapChildToParent(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) bool {
	run, err := store.GetMapRun(ctx, exec.MapRunArn)
	if err != nil {
		if errors.Is(err, sfnstore.ErrMapRunNotFound) {
			s.failUnrecoverableExecution(ctx, store, exec, "its Map Run no longer exists")
			return false
		}
		s.failUnrecoverableExecution(ctx, store, exec, "its Map Run record could not be read: "+err.Error())
		return false
	}
	parent, err := store.GetExecution(ctx, run.ExecutionArn)
	if err != nil || parent.Status != "RUNNING" {
		s.failUnrecoverableExecution(ctx, store, exec, "its parent execution is no longer running")
		return false
	}
	return false
}

// failUnrecoverableExecution terminates an execution that cannot be
// resumed so it never becomes a permanent zombie. The terminal event
// accompanies the status, as on every other terminal path. A parent failed
// here also fails its RUNNING Distributed Map children: they are never
// resumed directly (the parent dispatches them), so without this cascade
// they would stay RUNNING with no dispatcher behind them.
func (s *StepFunctionService) failUnrecoverableExecution(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution, cause string) {
	exec.Status = "FAILED"
	exec.Error = "States.Runtime"
	exec.Cause = "recovery after restart: " + cause
	exec.StopDate = time.Now().UTC()
	appendTerminalFailureEvent(ctx, store, exec, exec.Error, exec.Cause)
	if err := store.UpdateExecution(ctx, exec); err != nil {
		logs.Error("sfn: failed to mark an unrecoverable execution FAILED", logs.String("arn", exec.ExecutionArn), logs.Err(err))
	}
	logs.Warn("sfn: execution could not be recovered after restart",
		logs.String("arn", exec.ExecutionArn), logs.String("cause", cause))

	runs, err := store.ListMapRunsByExecution(ctx, exec.ExecutionArn)
	if err != nil {
		return
	}
	for _, run := range runs {
		children, cerr := store.ListAllExecutions(ctx, "", "RUNNING", run.MapRunArn, "")
		if cerr != nil {
			continue
		}
		for _, child := range children {
			s.failUnrecoverableExecution(ctx, store, child, "its parent execution could not be recovered")
		}
	}
}

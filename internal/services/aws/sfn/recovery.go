package sfn

import (
	"context"
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
// States.Runtime instead of being left as a zombie.
func (s *StepFunctionService) resumeRecoveredExecution(ctx context.Context, region string, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution) bool {
	sm, err := store.GetStateMachine(ctx, exec.StateMachineArn)
	if err != nil {
		s.failUnrecoverableExecution(ctx, store, exec, fmt.Sprintf("state machine %s no longer exists", exec.StateMachineArn))
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
	executor := NewExecutorWithStores(store, s.bus, s.accountID, region)
	execCtx, cancel := context.WithCancel(context.Background())
	store.RegisterExecution(executionArn, cancel)
	s.asyncWg.Add(1)
	go func() {
		defer s.asyncWg.Done()
		defer store.UnregisterExecution(executionArn)
		defer func() {
			if r := recover(); r != nil {
				logs.Error("sfn: panic in recovered execution", logs.String("arn", executionArn), logs.Any("panic", r))
				exec.Status = "FAILED"
				exec.Error = "States.InternalError"
				exec.Cause = fmt.Sprintf("internal panic: %v", r)
				exec.StopDate = time.Now().UTC()
				_ = store.UpdateExecution(context.Background(), exec)
			}
		}()
		if err := executor.ExecuteStateMachineFromState(execCtx, exec, rp.StateName, rp.Input, rp.LastEventId); err != nil {
			logs.Error("sfn: recovered execution failed", logs.String("arn", executionArn), logs.Err(err))
		}
	}()

	logs.Info("sfn: resumed execution after restart",
		logs.String("arn", executionArn),
		logs.String("stateMachineArn", exec.StateMachineArn),
		logs.String("state", rp.StateName))
	return true
}

// failUnrecoverableExecution terminates an execution that cannot be
// resumed so it never becomes a permanent zombie.
func (s *StepFunctionService) failUnrecoverableExecution(ctx context.Context, store *sfnstore.StepFunctionStore, exec *sfnstore.Execution, cause string) {
	exec.Status = "FAILED"
	exec.Error = "States.Runtime"
	exec.Cause = "recovery after restart: " + cause
	exec.StopDate = time.Now().UTC()
	if err := store.UpdateExecution(ctx, exec); err != nil {
		logs.Error("sfn: failed to mark an unrecoverable execution FAILED", logs.String("arn", exec.ExecutionArn), logs.Err(err))
	}
	logs.Warn("sfn: execution could not be recovered after restart",
		logs.String("arn", exec.ExecutionArn), logs.String("cause", cause))
}

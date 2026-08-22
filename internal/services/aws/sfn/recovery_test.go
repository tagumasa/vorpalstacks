package sfn

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func newRecoveryService(t *testing.T) (*StepFunctionService, *sfnstore.StepFunctionStore) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	t.Cleanup(func() { mgr.Close() })
	svc := NewStepFunctionService(mgr, "000000000000")
	store, err := svc.getStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("store for region: %v", err)
	}
	return svc, store
}

const recoveryDefinition = `{
	"StartAt": "First",
	"States": {
		"First": {"Type": "Pass", "Result": {"stage": "first"}, "Next": "Second"},
		"Second": {"Type": "Pass", "Result": {"stage": "second"}, "End": true}
	}
}`

// waitForExecution polls the execution record until it leaves RUNNING or
// the deadline passes, returning the final record.
func waitForExecution(t *testing.T, store *sfnstore.StepFunctionStore, arn string) *sfnstore.Execution {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		exec, err := store.GetExecution(t.Context(), arn)
		if err != nil {
			t.Fatalf("get execution: %v", err)
		}
		if exec.Status != "RUNNING" {
			return exec
		}
		time.Sleep(20 * time.Millisecond)
	}
	exec, err := store.GetExecution(t.Context(), arn)
	if err != nil {
		t.Fatalf("get execution: %v", err)
	}
	return exec
}

// TestRecoverRunningExecutionsResumesMidFlight pins the boot recovery: a
// RUNNING execution whose history shows First exited but Second never
// entered resumes from Second with First's output and runs to
// completion. RedriveCount stays untouched — recovery is not a
// user-initiated redrive.
func TestRecoverRunningExecutionsResumesMidFlight(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "recover-midflight", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}

	exec := sfnstore.NewExecution(sm.StateMachineArn, "midflight", `{"start":true}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":midflight"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	// History as a restart would have left it: First entered and exited,
	// Second never entered.
	for _, event := range []*sfnstore.ExecutionHistoryEvent{
		{
			ExecutionArn: exec.ExecutionArn, EventId: 1, PreviousEventId: 0,
			Type: "ExecutionStarted", Timestamp: time.Now().UTC(),
			ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{Input: exec.Input},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 2, PreviousEventId: 1,
			Type: "PassStateEntered", Timestamp: time.Now().UTC(),
			PassStateEnteredEventDetails: &sfnstore.PassStateEnteredEventDetails{Name: "First"},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 3, PreviousEventId: 2,
			Type: "PassStateExited", Timestamp: time.Now().UTC(),
			PassStateExitedEventDetails: &sfnstore.PassStateExitedEventDetails{Name: "First", Output: `{"stage":"first"}`},
		},
	} {
		if err := store.AddExecutionHistoryEvent(ctx, event); err != nil {
			t.Fatalf("seed history event: %v", err)
		}
	}

	svc.RecoverRunningExecutions()

	final := waitForExecution(t, store, exec.ExecutionArn)
	if final.Status != "SUCCEEDED" {
		t.Fatalf("recovered execution status = %s (error %q, cause %q), want SUCCEEDED", final.Status, final.Error, final.Cause)
	}
	if final.RedriveCount != 0 {
		t.Errorf("recovery bumped RedriveCount: %d, want 0", final.RedriveCount)
	}
}

// TestRecoverRunningExecutionsFailsWhenStateMachineGone pins that an
// execution whose state machine no longer exists is failed with
// States.Runtime at recovery instead of staying a permanent zombie.
func TestRecoverRunningExecutionsFailsWhenStateMachineGone(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	ghostArn := "arn:aws:states:us-east-1:000000000000:stateMachine:deleted-sm:ghost"
	exec := sfnstore.NewExecution("arn:aws:states:us-east-1:000000000000:stateMachine:deleted-sm", "ghost", `{}`, "")
	exec.ExecutionArn = ghostArn
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	svc.RecoverRunningExecutions()

	final := waitForExecution(t, store, ghostArn)
	if final.Status != "FAILED" {
		t.Fatalf("unrecoverable execution status = %s, want FAILED", final.Status)
	}
	if final.Error != "States.Runtime" {
		t.Errorf("unrecoverable execution error = %q, want States.Runtime", final.Error)
	}
}

// TestRecoverRunningExecutionsNoopWithoutRunning pins that the boot
// sweep leaves non-RUNNING executions untouched: a SUCCEEDED execution
// keeps its terminal state and output.
func TestRecoverRunningExecutionsNoopWithoutRunning(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "recover-noop", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}

	exec := sfnstore.NewExecution(sm.StateMachineArn, "finished", `{"start":true}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":finished"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	// CreateExecution forces the RUNNING status; terminal states are
	// written through updates.
	exec.Status = "SUCCEEDED"
	exec.Output = `{"stage":"second"}`
	if err := store.UpdateExecution(ctx, exec); err != nil {
		t.Fatalf("finish execution: %v", err)
	}

	svc.RecoverRunningExecutions()
	svc.Shutdown()

	got, err := store.GetExecution(ctx, exec.ExecutionArn)
	if err != nil {
		t.Fatalf("get execution: %v", err)
	}
	if got.Status != "SUCCEEDED" || got.Output != `{"stage":"second"}` {
		t.Errorf("finished execution disturbed by recovery: status %s output %q", got.Status, got.Output)
	}
}

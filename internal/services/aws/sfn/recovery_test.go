package sfn

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
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

// TestDetermineResumePointFollowsPagination pins that the resume analysis
// reads the whole history: the store clamps any single page at MaxPageSize,
// so a resume point computed from one page would anchor the resumed event
// ids at the clamp boundary and blindly overwrite persisted events.
func TestDetermineResumePointFollowsPagination(t *testing.T) {
	_, store := newRecoveryService(t)
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "resume-paged", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "paged", `{}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":paged"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	totalEvents := int(sfnstore.MaxPageSize) + 5
	for i := 1; i <= totalEvents; i++ {
		event := &sfnstore.ExecutionHistoryEvent{
			ExecutionArn:             exec.ExecutionArn,
			EventId:                  int64(i),
			Type:                     "PassStateEntered",
			Timestamp:                time.Now().UTC(),
			StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "First"},
		}
		if err := store.AddExecutionHistoryEvent(ctx, event); err != nil {
			t.Fatalf("seed event %d: %v", i, err)
		}
	}

	definition, err := parseStateMachineDefinition(recoveryDefinition)
	if err != nil {
		t.Fatalf("parse definition: %v", err)
	}
	rp, err := determineResumePoint(ctx, store, exec.ExecutionArn, definition)
	if err != nil {
		t.Fatalf("determineResumePoint: %v", err)
	}
	if rp.LastEventId != int64(totalEvents) {
		t.Fatalf("resume last event id = %d, want %d (a clamped single page would report %d)",
			rp.LastEventId, totalEvents, sfnstore.MaxPageSize)
	}
}

// TestDetermineResumePointResumesAtContainer pins that a failure inside a
// Parallel branch resumes at the container: branch inner states write their
// exit events into the parent history, and the resume anchor is the last
// TOP-LEVEL exit — a StartAt fallback here would re-run every succeeded
// top-level state on redrive.
func TestDetermineResumePointResumesAtContainer(t *testing.T) {
	_, store := newRecoveryService(t)
	ctx := t.Context()

	definition := `{
		"StartAt": "A",
		"States": {
			"A": {"Type": "Pass", "Result": {"stage": "a"}, "Next": "P"},
			"P": {"Type": "Parallel", "Branches": [
				{"StartAt": "Inner", "States": {"Inner": {"Type": "Pass", "End": true}}}
			], "End": true}
		}
	}`
	sm := &sfnstore.StateMachine{Name: "resume-container", Definition: definition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "container", `{}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":container"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	// History as a branch failure would have left it: A exited at top
	// level, then the branch's inner Pass exited INSIDE the container.
	for _, event := range []*sfnstore.ExecutionHistoryEvent{
		{
			ExecutionArn: exec.ExecutionArn, EventId: 1, PreviousEventId: 0,
			Type: "ExecutionStarted", Timestamp: time.Now().UTC(),
			ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{Input: exec.Input},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 2, PreviousEventId: 1,
			Type: "PassStateEntered", Timestamp: time.Now().UTC(),
			StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "A"},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 3, PreviousEventId: 2,
			Type: "PassStateExited", Timestamp: time.Now().UTC(),
			StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "A", Output: `{"stage":"a"}`},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 4, PreviousEventId: 3,
			Type: "ParallelStateEntered", Timestamp: time.Now().UTC(),
			StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "P"},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 5, PreviousEventId: 4,
			Type: "PassStateEntered", Timestamp: time.Now().UTC(),
			StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "Inner"},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 6, PreviousEventId: 5,
			Type: "PassStateExited", Timestamp: time.Now().UTC(),
			StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "Inner", Output: `{"inner":true}`},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 7, PreviousEventId: 6,
			Type: "ParallelStateFailed", Timestamp: time.Now().UTC(),
		},
	} {
		if err := store.AddExecutionHistoryEvent(ctx, event); err != nil {
			t.Fatalf("seed history event: %v", err)
		}
	}

	def, err := parseStateMachineDefinition(definition)
	if err != nil {
		t.Fatalf("parse definition: %v", err)
	}
	rp, err := determineResumePoint(ctx, store, exec.ExecutionArn, def)
	if err != nil {
		t.Fatalf("determineResumePoint: %v", err)
	}
	if rp.StateName != "P" {
		t.Fatalf("resume state = %q, want the container P (a StartAt fallback would re-run A)", rp.StateName)
	}
	if rp.Input != `{"stage":"a"}` {
		t.Errorf("resume input = %q, want A's output", rp.Input)
	}
	if rp.LastEventId != 7 {
		t.Errorf("resume last event id = %d, want 7", rp.LastEventId)
	}
}

// TestExecuteStateMachineFromStateTimeoutArming pins the timeout-contract
// split: restart recovery re-arms the state-machine-level timeout, a
// redrive does not ("When you redrive an execution, the state machine
// level timeout, if defined, is reset to 0").
func TestExecuteStateMachineFromStateTimeoutArming(t *testing.T) {
	_, store := newRecoveryService(t)
	ctx := t.Context()

	definition := `{"StartAt": "W", "TimeoutSeconds": 1, "States": {"W": {"Type": "Wait", "Seconds": 2, "End": true}}}`
	run := func(name string, arm bool) string {
		sm := &sfnstore.StateMachine{Name: "timeout-" + name, Definition: definition}
		if err := store.CreateStateMachine(ctx, sm); err != nil {
			t.Fatalf("create state machine: %v", err)
		}
		exec := sfnstore.NewExecution(sm.StateMachineArn, name, `{}`, "")
		exec.ExecutionArn = sm.StateMachineArn + ":" + name
		if err := store.CreateExecution(ctx, exec); err != nil {
			t.Fatalf("create execution: %v", err)
		}
		e := NewExecutorWithStores(store, eventbus.NewEventBus(), "000000000000", "us-east-1", nil)
		_ = e.ExecuteStateMachineFromState(ctx, exec, "W", "", 0, arm)
		final, err := store.GetExecution(ctx, exec.ExecutionArn)
		if err != nil {
			t.Fatalf("get execution: %v", err)
		}
		return final.Status
	}

	if got := run("armed", true); got != "TIMED_OUT" {
		t.Errorf("armed (recovery) resume status = %s, want TIMED_OUT", got)
	}
	if got := run("unarmed", false); got != "SUCCEEDED" {
		t.Errorf("unarmed (redrive) resume status = %s, want SUCCEEDED — a redrive resets the machine-level timeout", got)
	}
}

// TestRecoverRunningExecutionsSkipsMapChildren pins the single-ownership
// rule of the boot sweep for Distributed Map children: the sweep resumes
// the PARENT only, and the parent's re-dispatch reclaims and finishes the
// child — the child is never resumed directly beside the re-dispatching
// parent (duplicated items and colliding history ids).
func TestRecoverRunningExecutionsSkipsMapChildren(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	definition := `{
		"StartAt": "M",
		"States": {
			"M": {"Type": "Map", "ItemsPath": "$.v", "ItemProcessor": {
				"ProcessorConfig": {"Mode": "DISTRIBUTED", "ExecutionType": "STANDARD"},
				"StartAt": "W",
				"States": {"W": {"Type": "Pass", "ResultPath": "$", "End": true}}
			}, "End": true}
		}
	}`
	sm := &sfnstore.StateMachine{Name: "recover-map9", Definition: definition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}

	parentArn := "arn:aws:states:us-east-1:000000000000:execution:recover-map9:map-parent"
	parent := sfnstore.NewExecution(sm.StateMachineArn, "map-parent", `{"v":[1]}`, "")
	parent.ExecutionArn = parentArn
	if err := store.CreateExecution(ctx, parent); err != nil {
		t.Fatalf("create parent: %v", err)
	}
	if err := store.AddExecutionHistoryEvent(ctx, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn: parentArn, EventId: 1, PreviousEventId: 0,
		Type: "ExecutionStarted", Timestamp: time.Now().UTC(),
		ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{Input: parent.Input},
	}); err != nil {
		t.Fatalf("seed parent start: %v", err)
	}

	runArn := "arn:aws:states:us-east-1:000000000000:mapRun:recover-map9/existing"
	if err := store.CreateMapRun(ctx, &sfnstore.MapRun{
		MapRunArn:       runArn,
		ExecutionArn:    parentArn,
		StateMachineArn: sm.StateMachineArn,
		Name:            "M",
		Status:          "RUNNING",
		StartDate:       time.Now().UTC().Unix(),
		ItemCounts:      sfnstore.MapRunItemCounts{Pending: 1, Total: 1},
		ExecutionCounts: sfnstore.MapRunExecutionCounts{Pending: 1, Total: 1},
	}); err != nil {
		t.Fatalf("create map run: %v", err)
	}

	childArn := parentArn + ":M-0"
	child := sfnstore.NewExecution(sm.StateMachineArn, "map-parent:M-0", `1`, "")
	child.ExecutionArn = childArn
	child.MapRunArn = runArn
	child.ItemCount = 1
	if err := store.CreateExecution(ctx, child); err != nil {
		t.Fatalf("create child: %v", err)
	}
	// Mid-child crash shape: started, iterator state entered but never
	// exited.
	for _, event := range []*sfnstore.ExecutionHistoryEvent{
		{
			ExecutionArn: childArn, EventId: 1, PreviousEventId: 0,
			Type: "ExecutionStarted", Timestamp: time.Now().UTC(),
			ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{Input: child.Input},
		},
		{
			ExecutionArn: childArn, EventId: 2, PreviousEventId: 1,
			Type: "PassStateEntered", Timestamp: time.Now().UTC(),
			StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "W"},
		},
	} {
		if err := store.AddExecutionHistoryEvent(ctx, event); err != nil {
			t.Fatalf("seed child history: %v", err)
		}
	}

	svc.RecoverRunningExecutions()

	if final := waitForExecution(t, store, parentArn); final.Status != "SUCCEEDED" {
		t.Fatalf("recovered parent status = %s (error %q), want SUCCEEDED", final.Status, final.Error)
	}
	finalChild := waitForExecution(t, store, childArn)
	if finalChild.Status != "SUCCEEDED" {
		t.Fatalf("child status = %s (error %q), want SUCCEEDED via the parent's re-dispatch", finalChild.Status, finalChild.Error)
	}
	history, _, err := store.GetExecutionHistory(t.Context(), childArn, 100, "", false)
	if err != nil {
		t.Fatalf("child history: %v", err)
	}
	assertSequentialIds(t, history)
	redriven := 0
	for _, ev := range history {
		if ev.Type == "ExecutionRedriven" {
			redriven++
		}
	}
	if redriven != 1 {
		t.Errorf("child history has %d ExecutionRedriven events, want exactly 1 from the parent's single re-dispatch", redriven)
	}
}

// TestRecoverRunningExecutionsFailsOrphanedMapChild pins that a child
// whose parent execution is gone (no dispatcher will ever re-dispatch it)
// is failed at recovery instead of staying a permanent RUNNING zombie.
func TestRecoverRunningExecutionsFailsOrphanedMapChild(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	missingParentArn := "arn:aws:states:us-east-1:000000000000:execution:recover-gone/parent"
	runArn := "arn:aws:states:us-east-1:000000000000:mapRun:recover-gone/run"
	if err := store.CreateMapRun(ctx, &sfnstore.MapRun{
		MapRunArn:       runArn,
		ExecutionArn:    missingParentArn,
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:recover-gone",
		Name:            "M",
		Status:          "RUNNING",
		StartDate:       time.Now().UTC().Unix(),
	}); err != nil {
		t.Fatalf("create map run: %v", err)
	}

	childArn := "arn:aws:states:us-east-1:000000000000:execution:recover-gone/orphan"
	child := sfnstore.NewExecution("arn:aws:states:us-east-1:000000000000:stateMachine:recover-gone", "orphan", `1`, "")
	child.ExecutionArn = childArn
	child.MapRunArn = runArn
	if err := store.CreateExecution(ctx, child); err != nil {
		t.Fatalf("create child: %v", err)
	}

	svc.RecoverRunningExecutions()

	final := waitForExecution(t, store, childArn)
	if final.Status != "FAILED" {
		t.Fatalf("orphaned child status = %s, want FAILED", final.Status)
	}
	if final.Error != "States.Runtime" {
		t.Errorf("orphaned child error = %q, want States.Runtime", final.Error)
	}
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
			StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "First"},
		},
		{
			ExecutionArn: exec.ExecutionArn, EventId: 3, PreviousEventId: 2,
			Type: "PassStateExited", Timestamp: time.Now().UTC(),
			StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "First", Output: `{"stage":"first"}`},
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
	history, _, err := store.GetExecutionHistory(t.Context(), exec.ExecutionArn, 100, "", false)
	if err != nil {
		t.Fatalf("get history: %v", err)
	}
	assertSequentialIds(t, history)
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
	history, _, err := store.GetExecutionHistory(t.Context(), ghostArn, 100, "", false)
	if err != nil {
		t.Fatalf("get history: %v", err)
	}
	if len(history) == 0 || history[len(history)-1].Type != "ExecutionFailed" {
		t.Fatalf("the unrecoverable execution must end with an ExecutionFailed event; history: %+v", history)
	}
	assertSequentialIds(t, history)
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

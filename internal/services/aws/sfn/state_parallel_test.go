package sfn

import (
	"context"
	"errors"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestParallelCheckpointRecordedOnOriginalRun pins that the original (non
// redrive) run records its successful branch results: the FIRST redrive
// must find them and re-run only the failed branch ("Reschedules and
// redrives only those branches that failed or aborted"), and the redrive
// re-entry must not re-enter the checkpointed branch's states.
func TestParallelCheckpointRecordedOnOriginalRun(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "P",
		States: map[string]interface{}{
			"P": map[string]interface{}{
				"Type": "Parallel",
				"Branches": []interface{}{
					map[string]interface{}{
						"StartAt": "OK",
						"States": map[string]interface{}{
							"OK": map[string]interface{}{"Type": "Pass", "Result": "good", "End": true},
						},
					},
					map[string]interface{}{
						"StartAt": "Bad",
						"States": map[string]interface{}{
							"Bad": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
						},
					},
				},
				"End": true,
			},
		},
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/par1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "par1",
		Status:          "RUNNING",
		Input:           `{}`,
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist execution failed: %v", err)
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	eventId := int64(0)
	execCtx := &ExecutionContext{
		Execution:     exec,
		Definition:    def,
		CurrentState:  "P",
		Input:         exec.Input,
		EventId:       &eventId,
		States:        states,
		QueryLanguage: "JSONPath",
		MapItemIndex:  -1,
	}
	if execErr := e.executeStates(context.Background(), execCtx); execErr == nil {
		t.Fatal("the failing branch must fail the Parallel state")
	}

	cp, ok := exec.ParallelCheckpoints["P"]
	if !ok {
		t.Fatal("the original run must checkpoint its successful branches for the first redrive")
	}
	if got := cp.BranchResults[0]; got != `"good"` {
		t.Errorf("branch 0 checkpoint = %q, want the quoted %q", got, "good")
	}
	if _, has := cp.BranchResults[1]; has {
		t.Error("the failed branch must not be checkpointed")
	}

	// Redrive re-entry: the same context with IsRedrive set and the id
	// counter continued, as the resumed executor would run it.
	execCtx.IsRedrive = true
	execCtx.CurrentState = "P"
	execCtx.Input = exec.Input
	if execErr := e.executeStates(context.Background(), execCtx); execErr == nil {
		t.Fatal("the still-failing branch must fail the redriven Parallel state")
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 1000, "", false)
	if herr != nil {
		t.Fatalf("history failed: %v", herr)
	}
	okEntries := 0
	for _, ev := range history {
		if ev.Type == "PassStateEntered" && ev.StateEnteredEventDetails != nil && ev.StateEnteredEventDetails.Name == "OK" {
			okEntries++
		}
	}
	if okEntries != 1 {
		t.Errorf("branch OK entered %d times across the original run and the redrive, want 1 — the checkpointed branch must not re-run", okEntries)
	}
}

// TestParallelMidTaskCancellationClassifiesAborted pins the mid-task
// cancellation contract inside a container state: a StopExecution landing
// while a branch task is in flight must classify the interruption at every
// layer — ParallelStateAborted, never Failed, and no Catcher consumption —
// the abort path owns the history.
func TestParallelMidTaskCancellationClassifiesAborted(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	def := &sfnstore.StateMachineDefinition{
		StartAt: "DoWork",
		States: map[string]interface{}{
			"DoWork": map[string]interface{}{
				"Type": "Parallel",
				"Branches": []interface{}{
					map[string]interface{}{
						"StartAt": "Work",
						"States": map[string]interface{}{
							"Work": map[string]interface{}{"Type": "Task", "Resource": tokenTestActivityARN, "TimeoutSeconds": 10, "End": true},
						},
					},
				},
				"Catch": []interface{}{
					map[string]interface{}{"ErrorEquals": []interface{}{"States.ALL"}, "Next": "Fallback"},
				},
				"End": true,
			},
			"Fallback": map[string]interface{}{"Type": "Pass", "Result": "recovered", "End": true},
		},
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states
	execCtx.Definition = def

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- e.executeStates(ctx, execCtx) }()

	// A claimed activity task means the branch task is in flight and
	// waiting on the worker.
	workerCtx, wcancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer wcancel()
	if _, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1"); err != nil {
		t.Fatalf("no activity task became available: %v", err)
	}

	cancel()

	select {
	case runErr := <-done:
		if !errors.Is(runErr, context.Canceled) {
			t.Fatalf("executeStates returned %v, want a context.Canceled-classified error", runErr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeStates did not return after the cancellation")
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), execCtx.Execution.ExecutionArn, 200, "", false)
	if herr != nil {
		t.Fatalf("history failed: %v", herr)
	}
	sawAborted, sawFailed, sawFallback := false, false, false
	for _, ev := range history {
		switch ev.Type {
		case "ParallelStateAborted":
			sawAborted = true
		case "ParallelStateFailed":
			sawFailed = true
		}
		if ev.StateEnteredEventDetails != nil && ev.StateEnteredEventDetails.Name == "Fallback" {
			sawFallback = true
		}
	}
	if !sawAborted {
		t.Error("history lacks ParallelStateAborted — a mid-task cancellation must record the Aborted variant")
	}
	if sawFailed {
		t.Error("history records ParallelStateFailed for a cancellation — the interruption is not a branch failure")
	}
	if sawFallback {
		t.Error("the States.ALL Catcher consumed the cancellation — no recovery transition may run on a stopping execution")
	}
}

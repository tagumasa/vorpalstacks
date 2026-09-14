package sfn

import (
	"context"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestDeleteStateMachineCascadesActivityTasks pins the lifecycle symmetry:
// deleting a state machine removes its activity-task records, so token
// lookups and heartbeats cannot resolve a task whose execution is gone —
// the same cascade DeleteActivity runs for its own tasks.
func TestDeleteStateMachineCascadesActivityTasks(t *testing.T) {
	svc, store := newRecoveryService(t)
	_ = svc
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "cascade-sm", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "e1", "{}", "")
	exec.ExecutionArn = "arn:aws:states:us-east-1:000000000000:execution:cascade-sm:e1"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	task := &sfnstore.ActivityTask{
		ActivityArn:  "arn:aws:states:us-east-1:000000000000:activity:cascade",
		ExecutionArn: exec.ExecutionArn,
		Input:        "{}",
		TaskToken:    "tok-cascade-1",
	}
	if err := store.CreateActivityTask(task); err != nil {
		t.Fatalf("create activity task: %v", err)
	}

	if err := store.DeleteStateMachine(ctx, sm.StateMachineArn); err != nil {
		t.Fatalf("delete state machine: %v", err)
	}

	if _, err := store.GetActivityTaskByToken("tok-cascade-1"); err == nil {
		t.Error("activity task survived its state machine's deletion")
	}
}

// TestDeleteStateMachineCascadesMapRuns pins the lifecycle symmetry for the
// Map Run family: deleting a state machine removes the map runs of its
// executions, so DescribeMapRun cannot resolve a run whose execution is
// gone.
func TestDeleteStateMachineCascadesMapRuns(t *testing.T) {
	svc, store := newRecoveryService(t)
	_ = svc
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "cascade-maprun", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "e1", "{}", "")
	exec.ExecutionArn = "arn:aws:states:us-east-1:000000000000:execution:cascade-maprun:e1"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	mapRunArn := "arn:aws:states:us-east-1:000000000000:mapRun:cascade-maprun/e1/mapRun-1-cascade"
	if err := store.CreateMapRun(ctx, &sfnstore.MapRun{
		MapRunArn:       mapRunArn,
		ExecutionArn:    exec.ExecutionArn,
		StateMachineArn: sm.StateMachineArn,
		Name:            "M",
		Status:          "SUCCEEDED",
	}); err != nil {
		t.Fatalf("create map run: %v", err)
	}

	if err := store.DeleteStateMachine(ctx, sm.StateMachineArn); err != nil {
		t.Fatalf("delete state machine: %v", err)
	}

	if _, err := store.GetMapRun(ctx, mapRunArn); err == nil {
		t.Error("map run survived its state machine's deletion")
	}
}

// TestStartSyncExecutionStoppable pins that the synchronous start registers
// for StopExecution like the async launch: a stop must cancel the inline
// run, which returns ABORTED instead of running to its own terminal state
// and overwriting the aborted record.
func TestStartSyncExecutionStoppable(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	sm := &sfnstore.StateMachine{
		Name:       "sync-stop",
		Type:       "EXPRESS",
		Definition: `{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":60,"End":true}}}`,
	}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}

	done := make(chan map[string]interface{}, 1)
	go func() {
		resp, err := svc.startSyncExecutionCore(context.Background(), store, StartSyncExecutionInput{StateMachineArn: sm.StateMachineArn})
		if err != nil {
			t.Errorf("sync execution: %v", err)
			return
		}
		done <- resp
	}()

	// Wait until the execution exists and its wait state is entered — by
	// then the inline run has registered for cancellation.
	var execArn string
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		execs, err := store.ListAllExecutions(ctx, sm.StateMachineArn, "", "", "")
		if err == nil && len(execs) > 0 {
			execArn = execs[0].ExecutionArn
			history, _, herr := store.GetExecutionHistory(ctx, execArn, 100, "", false)
			if herr == nil && len(history) > 0 && history[len(history)-1].Type == "WaitStateEntered" {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if execArn == "" {
		t.Fatal("sync execution never appeared")
	}

	if _, err := svc.stopExecutionCore(ctx, store, StopExecutionInput{ExecutionArn: execArn}); err != nil {
		t.Fatalf("stop execution: %v", err)
	}

	select {
	case resp := <-done:
		if resp["status"] != "ABORTED" {
			t.Errorf("sync execution status after stop = %v, want ABORTED", resp["status"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("sync execution did not return after StopExecution — the inline run was not cancelled")
	}
}

package sfn

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestFailActivityTaskKeepsCauseSeparate pins the SendTaskFailure contract:
// the result the executor waits on carries the worker's error name and its
// cause as separate members — the model's ActivityFailedEventDetails pair —
// rather than one folded "error: cause" string.
func TestFailActivityTaskKeepsCauseSeparate(t *testing.T) {
	store := newHistoryTestStore(t)
	ctx := context.Background()

	activity := &Activity{Name: "cause-sep"}
	if err := store.CreateActivity(ctx, activity); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	task := &ActivityTask{
		ActivityArn:  activity.ActivityArn,
		ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:e",
		Input:        `{}`,
		TaskToken:    "token-cause-sep",
	}
	if err := store.CreateActivityTask(task); err != nil {
		t.Fatalf("create activity task: %v", err)
	}

	results := make(chan *ActivityTaskResult, 1)
	waitErr := make(chan error, 1)
	go func() {
		result, err := store.WaitForTaskResult(ctx, task.TaskToken, 5*time.Second, 0)
		if err != nil {
			waitErr <- err
			return
		}
		results <- result
	}()

	// Let the waiter register before reporting the failure.
	time.Sleep(50 * time.Millisecond)
	if err := store.FailActivityTask(task.TaskToken, "CustomActivityError", "the worker cause"); err != nil {
		t.Fatalf("fail activity task: %v", err)
	}

	select {
	case err := <-waitErr:
		t.Fatalf("wait for result: %v", err)
	case result := <-results:
		if result.Error == nil || result.Error.Error() != "CustomActivityError" {
			t.Errorf("result error = %v, want the bare error name CustomActivityError", result.Error)
		}
		if result.Cause != "the worker cause" {
			t.Errorf("result cause = %q, want the separate cause member", result.Cause)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the task result")
	}
}

// TestDeleteStateMachineKeepsForeignWaiters pins the waiter lifecycle
// across an unrelated machine's deletion: only the deleted machine's own
// waiters are closed; a foreign waiter keeps its registration so its
// worker's report still reaches it instead of stranding it for the rest of
// the task timeout.
func TestDeleteStateMachineKeepsForeignWaiters(t *testing.T) {
	store := newHistoryTestStore(t)
	ctx := context.Background()

	smB := &StateMachine{Name: "waiter-b"}
	if err := store.CreateStateMachine(ctx, smB); err != nil {
		t.Fatal(err)
	}

	task := &ActivityTask{
		TaskToken:    "token-foreign-waiter",
		ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:waiter-a:exec1",
		Input:        `{}`,
	}
	if err := store.CreateCallbackTask(task); err != nil {
		t.Fatal(err)
	}

	results := make(chan *ActivityTaskResult, 1)
	errs := make(chan error, 1)
	go func() {
		result, err := store.WaitForTaskResult(ctx, task.TaskToken, 5*time.Second, 0)
		if err != nil {
			errs <- err
			return
		}
		results <- result
	}()

	time.Sleep(50 * time.Millisecond)
	if err := store.DeleteStateMachine(ctx, smB.StateMachineArn); err != nil {
		t.Fatal(err)
	}

	if err := store.CompleteActivityTask(task.TaskToken, `{"ok":true}`); err != nil {
		t.Fatalf("complete after the foreign delete: %v", err)
	}
	select {
	case result := <-results:
		if result.Output != `{"ok":true}` {
			t.Errorf("waiter output = %q, want the worker's report", result.Output)
		}
	case err := <-errs:
		t.Fatalf("waiter failed after the foreign delete: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("the foreign waiter was dropped by the unrelated deletion — the report never arrived")
	}
}

// TestWaitForTaskResultSurvivesOrphanClose pins the orphaned-waiter path:
// the machine's own cascade closes the waiter channel, and that must
// surface as a task error — never as a nil result the consumer would
// dereference.
func TestWaitForTaskResultSurvivesOrphanClose(t *testing.T) {
	store := newHistoryTestStore(t)
	ctx := context.Background()

	sm := &StateMachine{Name: "waiter-own"}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatal(err)
	}

	task := &ActivityTask{
		TaskToken:    "token-own-waiter",
		ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:waiter-own:exec1",
		Input:        `{}`,
	}
	if err := store.CreateCallbackTask(task); err != nil {
		t.Fatal(err)
	}

	type outcome struct {
		result *ActivityTaskResult
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := store.WaitForTaskResult(ctx, task.TaskToken, 5*time.Second, 0)
		done <- outcome{result, err}
	}()

	time.Sleep(50 * time.Millisecond)
	if err := store.DeleteStateMachine(ctx, sm.StateMachineArn); err != nil {
		t.Fatal(err)
	}

	select {
	case o := <-done:
		if o.err == nil {
			t.Fatalf("orphaned waiter returned (%v, nil), want an error", o.result)
		}
		if o.result != nil {
			t.Error("orphaned waiter returned a result alongside the error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the orphaned waiter never returned")
	}
}

// TestDeleteActivityDoesNotBreakPollers pins that deleting an activity
// never crashes an in-flight poll: the queue is removed without closing, so
// the blocked receiver exits through its own context with the documented
// empty answer.
func TestDeleteActivityDoesNotBreakPollers(t *testing.T) {
	store := newHistoryTestStore(t)
	ctx := context.Background()

	arn := "arn:aws:states:us-east-1:000000000000:activity:poll-me"
	if err := store.CreateActivity(ctx, &Activity{Name: "poll-me"}); err != nil {
		t.Fatal(err)
	}

	pollCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() {
		task, err := store.GetActivityTask(pollCtx, arn, "worker-1")
		if err == nil && task != nil {
			err = fmt.Errorf("poller received a task after deletion: %+v", task)
		}
		done <- err
	}()

	time.Sleep(50 * time.Millisecond)
	if err := store.DeleteActivity(ctx, arn); err != nil {
		t.Fatal(err)
	}
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("poller after activity deletion: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the poller never returned after the activity deletion")
	}
}

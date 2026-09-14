package sfn

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// registerWaiter starts a WaitForTaskResult goroutine and blocks until its
// channel is registered in the pending map, so the test controls the
// interleaving from a known state.
type waiterOutcome struct {
	result *ActivityTaskResult
	err    error
}

func registerWaiter(t *testing.T, s *StepFunctionStore, token string, timeout time.Duration) <-chan waiterOutcome {
	t.Helper()
	done := make(chan waiterOutcome, 1)
	go func() {
		res, err := s.WaitForTaskResult(context.Background(), token, timeout, 0)
		done <- waiterOutcome{res, err}
	}()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		s.pendingTasksMu.RLock()
		_, registered := s.pendingTasks[token]
		s.pendingTasksMu.RUnlock()
		if registered {
			return done
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("waiter never registered its channel")
	return done
}

func createPendingTaskMachine(t *testing.T, s *StepFunctionStore, name string) string {
	t.Helper()
	sm := &StateMachine{
		Name:       name,
		Definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}`,
		RoleArn:    "arn:aws:iam::000000000000:role/sm",
	}
	if err := s.CreateStateMachine(context.Background(), sm); err != nil {
		t.Fatalf("create %s: %v", name, err)
	}
	return sm.StateMachineArn
}

// Deleting a state machine closes its own registered waiters: the wait
// surfaces a task error instead of hanging for the rest of the task
// timeout, the task record is cascade-deleted, and a late report on the
// spent token no longer resolves.
func TestDeleteStateMachineClosesOwnWaiters(t *testing.T) {
	s := newHistoryTestStore(t)
	ctx := context.Background()
	smArn := createPendingTaskMachine(t, s, "own-sm")

	execArn := "arn:aws:states:us-east-1:000000000000:execution:own-sm:exec1"
	token := "token-own"
	if err := s.CreateCallbackTask(&ActivityTask{TaskToken: token, ExecutionArn: execArn, Input: "{}"}); err != nil {
		t.Fatalf("register callback task: %v", err)
	}

	done := registerWaiter(t, s, token, 30*time.Second)
	if err := s.DeleteStateMachine(ctx, smArn); err != nil {
		t.Fatalf("delete state machine: %v", err)
	}

	select {
	case out := <-done:
		if !errors.Is(out.err, ErrTaskNotFound) {
			t.Fatalf("own waiter outcome = (%v, %v), want ErrTaskNotFound", out.result, out.err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("own waiter stranded by the deletion")
	}

	if err := s.CompleteActivityTask(token, "{}"); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("late report after cascade = %v, want ErrTaskNotFound", err)
	}
}

// A foreign machine's waiter survives an unrelated DeleteStateMachine and
// still receives its report afterwards: the deletion neither strands the
// waiter nor drops a report that lands once the deletion has finished.
func TestDeleteStateMachineKeepsForeignWaitersReportable(t *testing.T) {
	s := newHistoryTestStore(t)
	ctx := context.Background()
	deletedArn := createPendingTaskMachine(t, s, "deleted-sm")
	createPendingTaskMachine(t, s, "other-sm")

	execArn := "arn:aws:states:us-east-1:000000000000:execution:other-sm:exec1"
	token := "token-foreign"
	if err := s.CreateCallbackTask(&ActivityTask{TaskToken: token, ExecutionArn: execArn, Input: "{}"}); err != nil {
		t.Fatalf("register callback task: %v", err)
	}

	done := registerWaiter(t, s, token, 30*time.Second)
	if err := s.DeleteStateMachine(ctx, deletedArn); err != nil {
		t.Fatalf("delete unrelated machine: %v", err)
	}

	if err := s.CompleteActivityTask(token, `{"ok":true}`); err != nil {
		t.Fatalf("report for a foreign waiter after deletion: %v", err)
	}
	select {
	case out := <-done:
		if out.err != nil || out.result == nil || out.result.Output != `{"ok":true}` {
			t.Fatalf("foreign waiter outcome = (%v, %v), want the reported output", out.result, out.err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("foreign waiter never received its report")
	}
}

// The report must survive even when it races the deletion itself: with the
// foreign/own decision and the map replacement inside one lock hold, a
// concurrent SendTaskSuccess either delivers before the sweep, blocks on
// the mutex and delivers to the kept waiter, or lands after — never into a
// window where the foreign waiter is momentarily unregistered.
func TestDeleteStateMachineForeignReportConcurrent(t *testing.T) {
	s := newHistoryTestStore(t)
	ctx := context.Background()
	createPendingTaskMachine(t, s, "other-sm")

	execArn := "arn:aws:states:us-east-1:000000000000:execution:other-sm:exec1"
	for round := 0; round < 25; round++ {
		deletedArn := createPendingTaskMachine(t, s, "sweep-sm")
		token := "token-race-" + string(rune('a'+round))
		if err := s.CreateCallbackTask(&ActivityTask{TaskToken: token, ExecutionArn: execArn, Input: "{}"}); err != nil {
			t.Fatalf("round %d: register callback task: %v", round, err)
		}

		done := registerWaiter(t, s, token, 30*time.Second)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = s.DeleteStateMachine(ctx, deletedArn)
		}()
		go func() {
			defer wg.Done()
			_ = s.CompleteActivityTask(token, `{"round":true}`)
		}()
		wg.Wait()

		select {
		case out := <-done:
			if out.err != nil || out.result == nil || out.result.Output != `{"round":true}` {
				t.Fatalf("round %d: waiter outcome = (%v, %v), want the reported output", round, out.result, out.err)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("round %d: report lost to the deletion window", round)
		}
	}
}

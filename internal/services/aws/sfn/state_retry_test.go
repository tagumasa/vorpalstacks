package sfn

import (
	"context"
	"errors"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// The container retry contract: "Task, Parallel, and Map states can have a
// field named Retry" and "When a state has both Retry and Catch fields,
// Step Functions uses any appropriate retriers first" (error handling
// documentation). These pins drive the container failure paths end to end:
// a matching retrier re-enters the state (re-emitting its Entered event per
// attempt) before any Catch policy is consulted.

func countStateEvents(t *testing.T, store *sfnstore.StepFunctionStore, arn, eventType, stateName string) int {
	t.Helper()
	history, _, err := store.GetExecutionHistory(context.Background(), arn, 1000, "", false)
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	count := 0
	for _, ev := range history {
		if ev.Type != eventType {
			continue
		}
		if stateName == "" {
			count++
			continue
		}
		if ev.StateEnteredEventDetails != nil && ev.StateEnteredEventDetails.Name == stateName {
			count++
		}
	}
	return count
}

// TestParallelRetryReentersState pins that a Parallel state with a failing
// branch and a matching retrier re-enters the state once per retry attempt
// instead of failing through to Catch immediately.
func TestParallelRetryReentersState(t *testing.T) {
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
						"StartAt": "Bad",
						"States": map[string]interface{}{
							"Bad": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
						},
					},
				},
				"Retry": []interface{}{
					map[string]interface{}{
						"ErrorEquals":     []interface{}{"States.ALL"},
						"IntervalSeconds": 1,
						"BackoffRate":     1,
						"MaxAttempts":     2,
					},
				},
				"End": true,
			},
		},
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:retry1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "retry1",
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
		Execution: exec, Definition: def, CurrentState: "P",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	execErr := e.executeStates(context.Background(), execCtx)
	if execErr == nil {
		t.Fatal("the permanently failing branch must fail the Parallel state after the retries")
	}

	entered := countStateEvents(t, store, exec.ExecutionArn, "ParallelStateEntered", "P")
	if entered != 3 {
		t.Errorf("ParallelStateEntered = %d, want 3 (initial attempt + 2 retries)", entered)
	}
}

// TestMapRetryReentersState pins the same contract on the Map state's
// failure path: one retry attempt re-enters the Map state.
func TestMapRetryReentersState(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.v",
				"ItemProcessor": map[string]interface{}{
					"StartAt": "Bad",
					"States": map[string]interface{}{
						"Bad": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
					},
				},
				"Retry": []interface{}{
					map[string]interface{}{
						"ErrorEquals":     []interface{}{"CustomError"},
						"IntervalSeconds": 1,
						"BackoffRate":     1,
						"MaxAttempts":     1,
					},
				},
				"End": true,
			},
		},
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:retry2",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "retry2",
		Status:          "RUNNING",
		Input:           `{"v":[1]}`,
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
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	if execErr := e.executeStates(context.Background(), execCtx); execErr == nil {
		t.Fatal("the permanently failing iteration must fail the Map state after the retry")
	}

	entered := countStateEvents(t, store, exec.ExecutionArn, "MapStateEntered", "M")
	if entered != 2 {
		t.Errorf("MapStateEntered = %d, want 2 (initial attempt + 1 retry)", entered)
	}
}

// TestParallelRetryWaitCancellationIsInterrupted pins the container retry
// wait's interruption contract: a StopExecution that lands while the
// retrier's backoff interval is running leaves through the interruption
// vehicle (wrapping context.Canceled), never as a plain States.Timeout that
// the surrounding layers would classify as a machine deadline. Task, Map and
// Parallel share the single sleepForRetry wait, so the contract is one.
func TestParallelRetryWaitCancellationIsInterrupted(t *testing.T) {
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
						"StartAt": "Bad",
						"States": map[string]interface{}{
							"Bad": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
						},
					},
				},
				"Retry": []interface{}{
					map[string]interface{}{
						"ErrorEquals":     []interface{}{"States.ALL"},
						"IntervalSeconds": 30,
						"BackoffRate":     1,
						"MaxAttempts":     1,
					},
				},
				"End": true,
			},
		},
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:retrywait1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "retrywait1",
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
		Execution: exec, Definition: def, CurrentState: "P",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- e.executeStates(ctx, execCtx) }()

	// The attempt's ParallelStateFailed event means the backoff wait for the
	// retry is now running.
	waitDeadline := time.Now().Add(10 * time.Second)
	for countStateEvents(t, store, exec.ExecutionArn, "ParallelStateFailed", "") == 0 {
		if time.Now().After(waitDeadline) {
			t.Fatal("no ParallelStateFailed event — the failing branch never ran")
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	select {
	case runErr := <-done:
		if !errors.Is(runErr, context.Canceled) {
			t.Fatalf("executeStates returned %v, want a context.Canceled-classified interruption", runErr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeStates did not return after the cancellation cut the retry wait")
	}
}

// TestTaskRetryWaitCancellationIsInterrupted pins the same contract on the
// Task retry path: the wait is the shared sleepForRetry, so a cancellation
// during the backoff interval propagates with the cancellation identity.
func TestTaskRetryWaitCancellationIsInterrupted(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "T",
		States: map[string]interface{}{
			"T": map[string]interface{}{
				"Type":     "Task",
				"Resource": "arn:aws:states:::nosuch:integration",
				"Retry": []interface{}{
					map[string]interface{}{
						"ErrorEquals":     []interface{}{"States.ALL"},
						"IntervalSeconds": 30,
						"BackoffRate":     1,
						"MaxAttempts":     1,
					},
				},
				"End": true,
			},
		},
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:retrywait2",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "retrywait2",
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
		Execution: exec, Definition: def, CurrentState: "T",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- e.executeStates(ctx, execCtx) }()

	waitDeadline := time.Now().Add(10 * time.Second)
	for countStateEvents(t, store, exec.ExecutionArn, "TaskStartFailed", "") == 0 {
		if time.Now().After(waitDeadline) {
			t.Fatal("no TaskStartFailed event — the failing dispatch never ran")
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	select {
	case runErr := <-done:
		if !errors.Is(runErr, context.Canceled) {
			t.Fatalf("executeStates returned %v, want a context.Canceled-classified interruption", runErr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeStates did not return after the cancellation cut the retry wait")
	}
}

// TestRetryMaxAttemptsDefaultsToThree pins the documented default on the
// real path: "MaxAttempts … a positive integer that represents the maximum
// number of retry attempts (3 by default)". A retrier without MaxAttempts
// retries three times, not zero.
func TestRetryMaxAttemptsDefaultsToThree(t *testing.T) {
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
						"StartAt": "Bad",
						"States": map[string]interface{}{
							"Bad": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
						},
					},
				},
				"Retry": []interface{}{
					map[string]interface{}{
						"ErrorEquals":     []interface{}{"States.ALL"},
						"IntervalSeconds": 1,
						"BackoffRate":     1,
					},
				},
				"End": true,
			},
		},
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:retry3",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "retry3",
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
		Execution: exec, Definition: def, CurrentState: "P",
		Input: exec.Input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	if execErr := e.executeStates(context.Background(), execCtx); execErr == nil {
		t.Fatal("the permanently failing branch must fail the Parallel state after the default retries")
	}

	entered := countStateEvents(t, store, exec.ExecutionArn, "ParallelStateEntered", "P")
	if entered != 4 {
		t.Errorf("ParallelStateEntered = %d, want 4 (initial attempt + 3 default retries)", entered)
	}
}

// TestCalculateBackoffIntervalJitter pins the JitterStrategy arithmetic:
// NONE is deterministic exponential backoff; FULL randomises the interval
// between 0 and the attempt's deterministic interval (MaxDelaySeconds cap
// included) — "the first retry interval is randomized between 0 and 2
// seconds, the second … between 0 and 4 seconds" for IntervalSeconds 2.
func TestCalculateBackoffIntervalJitter(t *testing.T) {
	e := &Executor{}

	none := &sfnstore.RetryPolicy{IntervalSeconds: 2, BackoffRate: 2}
	if got := e.calculateBackoffInterval(none, 1); got != 2*time.Second {
		t.Errorf("attempt 1 without jitter = %v, want 2s", got)
	}
	if got := e.calculateBackoffInterval(none, 3); got != 8*time.Second {
		t.Errorf("attempt 3 without jitter = %v, want 8s", got)
	}

	full := &sfnstore.RetryPolicy{IntervalSeconds: 2, BackoffRate: 2, JitterStrategy: "FULL"}
	for attempt := int32(1); attempt <= 3; attempt++ {
		window := time.Duration(2) * time.Second << (attempt - 1)
		for i := 0; i < 50; i++ {
			got := e.calculateBackoffInterval(full, attempt)
			if got < 0 || got >= window {
				t.Fatalf("FULL jitter attempt %d gave %v, want within [0, %v)", attempt, got, window)
			}
		}
	}

	capped := &sfnstore.RetryPolicy{IntervalSeconds: 2, BackoffRate: 2, MaxDelaySeconds: 5, JitterStrategy: "FULL"}
	for i := 0; i < 50; i++ {
		got := e.calculateBackoffInterval(capped, 6)
		if got < 0 || got > 5*time.Second {
			t.Fatalf("FULL jitter with MaxDelaySeconds gave %v, want within [0, 5s]", got)
		}
	}
}

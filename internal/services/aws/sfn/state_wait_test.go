package sfn

import (
	"context"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// newWaitExecutor builds an executor and execution context bound to a
// real store so history logging works, mirroring the TestState harness.
func newWaitExecutor(t *testing.T) (*Executor, *ExecutionContext) {
	t.Helper()
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	sm := &sfnstore.StateMachine{Name: "wait-unit", Definition: recoveryDefinition}
	if err := store.CreateStateMachine(ctx, sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "wait-unit", `{"expiry":"2024-03-14T01:59:00Z"}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":wait-unit"
	if err := store.CreateExecution(ctx, exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	executor := NewExecutorWithStores(store, svc.bus, "000000000000", "us-east-1")
	eventId := int64(1)
	execCtx := &ExecutionContext{
		Execution:     exec,
		CurrentState:  "W",
		Input:         exec.Input,
		EventId:       &eventId,
		QueryLanguage: "JSONPath",
		VariableScope: NewVariableScope(nil),
	}
	return executor, execCtx
}

// executionErrorCode unwraps the States.* error code of a wait failure.
func executionErrorCode(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if execErr, ok := err.(*ExecutionError); ok {
		return execErr.ErrorCode
	}
	t.Fatalf("error %v is not an ExecutionError", err)
	return ""
}

// TestExecuteWaitInvalidValuesFailExecution pins that the JSONPath Wait
// fails the execution with States.Runtime on contract violations instead
// of silently skipping the wait: an invalid Timestamp literal, a
// TimestampPath that selects nothing / a non-string / an invalid
// timestamp, and a SecondsPath that selects a negative or fractional
// value.
func TestExecuteWaitInvalidValuesFailExecution(t *testing.T) {
	executor, execCtx := newWaitExecutor(t)

	cases := []struct {
		name  string
		state *sfnstore.WaitState
	}{
		{"invalid timestamp literal", &sfnstore.WaitState{Timestamp: "2024-03-14 01:59:00"}},
		{"timestamp path selects nothing", &sfnstore.WaitState{TimestampPath: "$.missing"}},
		{"timestamp path selects a non-string", &sfnstore.WaitState{TimestampPath: "$.count"}},
		{"timestamp path selects an invalid timestamp", &sfnstore.WaitState{TimestampPath: "$.badExpiry"}},
		{"seconds path selects a negative value", &sfnstore.WaitState{SecondsPath: "$.negative"}},
		{"seconds path selects a fractional value", &sfnstore.WaitState{SecondsPath: "$.fraction"}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			execCtx.Input = `{"expiry":"2024-03-14T01:59:00Z","count":5,"badExpiry":"not-a-timestamp","negative":-1,"fraction":1.5}`
			_, _, err := executor.executeWait(context.Background(), execCtx, tt.state)
			if code := executionErrorCode(t, err); code != "States.Runtime" {
				t.Errorf("error code = %q, want States.Runtime", code)
			}
		})
	}
}

// TestExecuteWaitValidTimestampsProceed pins the accepting side: a past
// timestamp and a resolving timestamp path both proceed without waiting,
// and milliseconds are truncated rather than rejected.
func TestExecuteWaitValidTimestampsProceed(t *testing.T) {
	executor, execCtx := newWaitExecutor(t)

	past := &sfnstore.WaitState{Timestamp: "2024-03-14T01:59:00.123Z"}
	if _, _, err := executor.executeWait(context.Background(), execCtx, past); err != nil {
		t.Errorf("past timestamp wait failed: %v", err)
	}

	viaPath := &sfnstore.WaitState{TimestampPath: "$.expiry"}
	if _, _, err := executor.executeWait(context.Background(), execCtx, viaPath); err != nil {
		t.Errorf("timestamp path wait failed: %v", err)
	}
}

// TestExecuteWaitJSONataTimestamp pins the JSONata Wait contract: a
// Timestamp literal or expression is honoured (a past time proceeds
// immediately), while a non-string expression result or a non-timestamp
// value fails with States.Runtime and an evaluation failure surfaces
// States.QueryEvaluationError.
func TestExecuteWaitJSONataTimestamp(t *testing.T) {
	executor, execCtx := newWaitExecutor(t)
	execCtx.QueryLanguage = "JSONata"
	execCtx.Input = `{"expiry":"2024-03-14T01:59:00Z"}`

	past := &sfnstore.WaitState{Timestamp: "2024-03-14T01:59:00Z"}
	if _, _, err := executor.executeWaitJSONata(context.Background(), execCtx, past); err != nil {
		t.Errorf("literal past timestamp wait failed: %v", err)
	}

	viaExpression := &sfnstore.WaitState{Timestamp: "{% $states.input.expiry %}"}
	if _, _, err := executor.executeWaitJSONata(context.Background(), execCtx, viaExpression); err != nil {
		t.Errorf("timestamp expression wait failed: %v", err)
	}

	nonString := &sfnstore.WaitState{Timestamp: "{% 42 %}"}
	_, _, err := executor.executeWaitJSONata(context.Background(), execCtx, nonString)
	if code := executionErrorCode(t, err); code != "States.Runtime" {
		t.Errorf("non-string timestamp expression error code = %q, want States.Runtime", code)
	}

	notATimestamp := &sfnstore.WaitState{Timestamp: "{% 'tomorrow' %}"}
	_, _, err = executor.executeWaitJSONata(context.Background(), execCtx, notATimestamp)
	if code := executionErrorCode(t, err); code != "States.Runtime" {
		t.Errorf("non-timestamp expression error code = %q, want States.Runtime", code)
	}

	brokenExpression := &sfnstore.WaitState{Timestamp: "{% $ }( %}"}
	_, _, err = executor.executeWaitJSONata(context.Background(), execCtx, brokenExpression)
	if code := executionErrorCode(t, err); code != "States.QueryEvaluationError" {
		t.Errorf("broken expression error code = %q, want States.QueryEvaluationError", code)
	}
}

// TestExecuteWaitJSONataSecondsContract pins the JSONata Seconds
// contract: an integer expression waits and a fractional or negative
// result fails with States.Runtime.
func TestExecuteWaitJSONataSecondsContract(t *testing.T) {
	executor, execCtx := newWaitExecutor(t)
	execCtx.QueryLanguage = "JSONata"
	execCtx.Input = `{}`

	zero := &sfnstore.WaitState{Seconds: "{% 0 %}"}
	if _, _, err := executor.executeWaitJSONata(context.Background(), execCtx, zero); err != nil {
		t.Errorf("zero seconds expression wait failed: %v", err)
	}

	fractional := &sfnstore.WaitState{Seconds: "{% 1.5 %}"}
	_, _, err := executor.executeWaitJSONata(context.Background(), execCtx, fractional)
	if code := executionErrorCode(t, err); code != "States.Runtime" {
		t.Errorf("fractional seconds error code = %q, want States.Runtime", code)
	}

	negative := &sfnstore.WaitState{Seconds: "{% -1 %}"}
	_, _, err = executor.executeWaitJSONata(context.Background(), execCtx, negative)
	if code := executionErrorCode(t, err); code != "States.Runtime" {
		t.Errorf("negative seconds error code = %q, want States.Runtime", code)
	}
}

// TestExecuteWaitTruncatesMilliseconds pins the documented truncation:
// a timestamp carrying milliseconds waits only to the whole second.
func TestExecuteWaitTruncatesMilliseconds(t *testing.T) {
	executor, execCtx := newWaitExecutor(t)

	// Half a second past a whole-second boundary roughly one second in
	// the future: the truncated boundary fires by the whole-second mark.
	future := time.Now().UTC().Add(1500 * time.Millisecond).Truncate(time.Second).Add(500 * time.Millisecond)
	state := &sfnstore.WaitState{Timestamp: future.Format("2006-01-02T15:04:05.000Z")}
	started := time.Now()
	if _, _, err := executor.executeWait(context.Background(), execCtx, state); err != nil {
		t.Fatalf("wait failed: %v", err)
	}
	waited := time.Since(started)
	if waited >= 1500*time.Millisecond {
		t.Errorf("waited %v, milliseconds were not truncated", waited)
	}
}

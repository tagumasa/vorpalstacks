package sfn

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// launchHistoryExecution creates a state machine and execution, launches it
// through the real asynchronous path and returns the service, store and
// execution ARN for history assertions.
func launchHistoryExecution(t *testing.T, name, definition string) (*StepFunctionService, *sfnstore.StepFunctionStore, string) {
	t.Helper()
	svc, store := newRecoveryService(t)
	sm := &sfnstore.StateMachine{Name: name, Definition: definition}
	if err := store.CreateStateMachine(t.Context(), sm); err != nil {
		t.Fatalf("create state machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "exec-"+name, `{"in":1}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":exec-" + name
	if err := store.CreateExecution(t.Context(), exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}
	svc.launchExecution(store, exec)
	return svc, store, exec.ExecutionArn
}

// readHistoryUntilTerminal polls the stored history until the last event is
// the wanted terminal type. StopExecution persists the ABORTED status on its
// own record copy before the running goroutine appends the terminal event, so
// the record status alone does not guarantee the event is present yet.
func readHistoryUntilTerminal(t *testing.T, store *sfnstore.StepFunctionStore, arn, terminalType string) []*sfnstore.ExecutionHistoryEvent {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		history, _, err := store.GetExecutionHistory(t.Context(), arn, 100, "", false)
		if err != nil {
			t.Fatalf("get history: %v", err)
		}
		if len(history) > 0 && history[len(history)-1].Type == terminalType {
			return history
		}
		if time.Now().After(deadline) {
			t.Fatalf("terminal event %q never appeared; history: %+v", terminalType, history)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// assertSequentialIds pins the AWS contract that event IDs increment by one
// from 1 with no reuse: the store keys history rows on EventId, so a reused
// ID silently overwrites the earlier event instead of appending.
func assertSequentialIds(t *testing.T, history []*sfnstore.ExecutionHistoryEvent) {
	t.Helper()
	for i, ev := range history {
		if ev.EventId != int64(i+1) {
			t.Errorf("event %d has EventId %d, want %d — IDs must increment by one with no reuse", i, ev.EventId, i+1)
		}
	}
}

// TestTerminalHistoryEventTakesNextEventId runs every terminal path through
// the real launch (and stop) path and pins the full tail sequence: the final
// state event survives and the terminal event carries the previous event's
// ID plus one.
func TestTerminalHistoryEventTakesNextEventId(t *testing.T) {
	t.Run("Succeeded", func(t *testing.T) {
		_, store, arn := launchHistoryExecution(t, "hist-succeeded",
			`{"StartAt":"P","States":{"P":{"Type":"Pass","Result":{"ok":true},"End":true}}}`)
		final := waitForExecution(t, store, arn)
		if final.Status != "SUCCEEDED" {
			t.Fatalf("status = %s, want SUCCEEDED", final.Status)
		}
		history := readHistoryUntilTerminal(t, store, arn, "ExecutionSucceeded")
		assertSequentialIds(t, history)
		last, pre := history[len(history)-1], history[len(history)-2]
		if pre.Type != "PassStateExited" {
			t.Errorf("event before terminal = %s, want PassStateExited — a reused ID overwrites it", pre.Type)
		}
		if last.EventId != pre.EventId+1 {
			t.Errorf("terminal EventId = %d, want %d", last.EventId, pre.EventId+1)
		}
	})

	t.Run("Failed", func(t *testing.T) {
		_, store, arn := launchHistoryExecution(t, "hist-failed",
			`{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"CustomError","Cause":"boom"}}}`)
		final := waitForExecution(t, store, arn)
		if final.Status != "FAILED" {
			t.Fatalf("status = %s, want FAILED", final.Status)
		}
		history := readHistoryUntilTerminal(t, store, arn, "ExecutionFailed")
		assertSequentialIds(t, history)
		last, pre := history[len(history)-1], history[len(history)-2]
		if pre.Type != "FailStateEntered" {
			t.Errorf("event before terminal = %s, want FailStateEntered — a reused ID overwrites it", pre.Type)
		}
		if last.EventId != pre.EventId+1 {
			t.Errorf("terminal EventId = %d, want %d", last.EventId, pre.EventId+1)
		}
	})

	t.Run("TimedOut", func(t *testing.T) {
		_, store, arn := launchHistoryExecution(t, "hist-timedout",
			`{"StartAt":"W","TimeoutSeconds":1,"States":{"W":{"Type":"Wait","Seconds":60,"End":true}}}`)
		final := waitForExecution(t, store, arn)
		if final.Status != "TIMED_OUT" {
			t.Fatalf("status = %s, want TIMED_OUT", final.Status)
		}
		history := readHistoryUntilTerminal(t, store, arn, "ExecutionTimedOut")
		assertSequentialIds(t, history)
		last, pre := history[len(history)-1], history[len(history)-2]
		if pre.Type != "WaitStateEntered" {
			t.Errorf("event before terminal = %s, want WaitStateEntered — a reused ID overwrites it", pre.Type)
		}
		if last.EventId != pre.EventId+1 {
			t.Errorf("terminal EventId = %d, want %d", last.EventId, pre.EventId+1)
		}
	})

	t.Run("Aborted", func(t *testing.T) {
		svc, store, arn := launchHistoryExecution(t, "hist-aborted",
			`{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":60,"End":true}}}`)

		// Stop once the wait state is entered, so the stop lands mid-wait.
		deadline := time.Now().Add(5 * time.Second)
		entered := false
		for !entered && time.Now().Before(deadline) {
			history, _, err := store.GetExecutionHistory(t.Context(), arn, 100, "", false)
			if err != nil {
				t.Fatalf("get history: %v", err)
			}
			for _, ev := range history {
				if ev.Type == "WaitStateEntered" {
					entered = true
					break
				}
			}
			if !entered {
				time.Sleep(10 * time.Millisecond)
			}
		}
		if !entered {
			t.Fatal("wait state was never entered before the stop")
		}

		if _, err := svc.stopExecutionCore(t.Context(), store, StopExecutionInput{ExecutionArn: arn}); err != nil {
			t.Fatalf("stop execution: %v", err)
		}
		final := waitForExecution(t, store, arn)
		if final.Status != "ABORTED" {
			t.Fatalf("status = %s, want ABORTED", final.Status)
		}
		history := readHistoryUntilTerminal(t, store, arn, "ExecutionAborted")
		assertSequentialIds(t, history)
		last, pre := history[len(history)-1], history[len(history)-2]
		// The running state records its modelled Aborted variant before the
		// execution's terminal event.
		if pre.Type != "WaitStateAborted" {
			t.Errorf("event before terminal = %s, want WaitStateAborted — a reused ID overwrites it", pre.Type)
		}
		if last.EventId != pre.EventId+1 {
			t.Errorf("terminal EventId = %d, want %d", last.EventId, pre.EventId+1)
		}
	})
}

// TestRedriveKeepsSequentialEventIds pins the resume counter contract: the
// ExecutionRedrived event and every resumed state event continue the id
// sequence one past the prior history with no gap and a correct
// previousEventId chain, per "Events are numbered sequentially, starting at
// one".
func TestRedriveKeepsSequentialEventIds(t *testing.T) {
	svc, store, arn := launchHistoryExecution(t, "hist-redrive",
		`{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"CustomError","Cause":"boom"}}}`)
	final := waitForExecution(t, store, arn)
	if final.Status != "FAILED" {
		t.Fatalf("status = %s, want FAILED", final.Status)
	}
	first := readHistoryUntilTerminal(t, store, arn, "ExecutionFailed")
	assertSequentialIds(t, first)

	if _, err := svc.redriveExecutionCore(t.Context(), store, RedriveExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("redrive execution: %v", err)
	}
	redriven := waitForExecution(t, store, arn)
	if redriven.Status != "FAILED" {
		t.Fatalf("redriven status = %s, want FAILED", redriven.Status)
	}
	history := readHistoryUntilTerminal(t, store, arn, "ExecutionFailed")
	if len(history) <= len(first) {
		t.Fatalf("history did not grow after the redrive: %d events", len(history))
	}
	assertSequentialIds(t, history)

	var redriveEvent *sfnstore.ExecutionHistoryEvent
	for _, ev := range history {
		if ev.Type == "ExecutionRedriven" {
			redriveEvent = ev
			break
		}
	}
	if redriveEvent == nil {
		t.Fatal("ExecutionRedrived event missing from the redriven history")
	}
	if redriveEvent.EventId != first[len(first)-1].EventId+1 {
		t.Errorf("ExecutionRedrived EventId = %d, want %d — the redrive must continue the sequence",
			redriveEvent.EventId, first[len(first)-1].EventId+1)
	}
	if redriveEvent.PreviousEventId != first[len(first)-1].EventId {
		t.Errorf("ExecutionRedrived previousEventId = %d, want %d", redriveEvent.PreviousEventId, first[len(first)-1].EventId)
	}
	next := history[redriveEvent.EventId]
	if next.EventId != redriveEvent.EventId+1 {
		t.Errorf("first resumed event EventId = %d, want %d — a resume must not skip an id",
			next.EventId, redriveEvent.EventId+1)
	}
	if next.PreviousEventId != redriveEvent.EventId {
		t.Errorf("first resumed event previousEventId = %d, want %d — it must point at the ExecutionRedriven event",
			next.PreviousEventId, redriveEvent.EventId)
	}
}

// TestRedriveAfterChoiceResumesAtTakenTransition pins the resume contract
// for Choice states: a redriven execution whose last exited state is a
// Choice follows the transition the choice actually took (recorded on the
// exit event) instead of restarting the machine from StartAt and
// re-executing states that already succeeded.
func TestRedriveAfterChoiceResumesAtTakenTransition(t *testing.T) {
	svc, store, arn := launchHistoryExecution(t, "hist-choice-redrive",
		`{"StartAt":"C","States":{
			"C":{"Type":"Choice","Choices":[{"Variable":"$.x","StringEquals":"1","Next":"F"}],"Default":"F"},
			"F":{"Type":"Fail","Error":"CustomError","Cause":"boom"}
		}}`)
	final := waitForExecution(t, store, arn)
	if final.Status != "FAILED" {
		t.Fatalf("status = %s, want FAILED", final.Status)
	}

	if _, err := svc.redriveExecutionCore(t.Context(), store, RedriveExecutionInput{ExecutionArn: arn}); err != nil {
		t.Fatalf("redrive execution: %v", err)
	}
	redriven := waitForExecution(t, store, arn)
	if redriven.Status != "FAILED" {
		t.Fatalf("redriven status = %s, want FAILED", redriven.Status)
	}

	history := readHistoryUntilTerminal(t, store, arn, "ExecutionFailed")
	assertSequentialIds(t, history)
	entered, failed := 0, 0
	for _, ev := range history {
		switch ev.Type {
		case "ChoiceStateEntered":
			entered++
		case "FailStateEntered":
			failed++
		}
	}
	if entered != 1 {
		t.Errorf("ChoiceStateEntered appears %d times after the redrive, want 1 — the resume must not re-execute the choice", entered)
	}
	if failed != 2 {
		t.Errorf("FailStateEntered appears %d times, want 2 (once per attempt)", failed)
	}
}

// TestStopExecutionPairSurvivesTheCancellationRace pins the StopExecution
// contract: the caller's error and cause are the execution's reported pair —
// both on DescribeExecution and on the ExecutionAborted history event —
// instead of being overwritten by the executor's generic strings.
func TestStopExecutionPairSurvivesTheCancellationRace(t *testing.T) {
	svc, store, arn := launchHistoryExecution(t, "hist-stop-pair",
		`{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":60,"End":true}}}`)

	deadline := time.Now().Add(5 * time.Second)
	entered := false
	for !entered && time.Now().Before(deadline) {
		history, _, err := store.GetExecutionHistory(t.Context(), arn, 100, "", false)
		if err != nil {
			t.Fatalf("get history: %v", err)
		}
		for _, ev := range history {
			if ev.Type == "WaitStateEntered" {
				entered = true
			}
		}
		if !entered {
			time.Sleep(10 * time.Millisecond)
		}
	}
	if !entered {
		t.Fatal("wait state was never entered before the stop")
	}

	if _, err := svc.stopExecutionCore(t.Context(), store, StopExecutionInput{
		ExecutionArn: arn,
		Error:        "CustomStopError",
		Cause:        "stopped by the test",
	}); err != nil {
		t.Fatalf("stop execution: %v", err)
	}

	final := waitForExecution(t, store, arn)
	if final.Status != "ABORTED" {
		t.Fatalf("status = %s, want ABORTED", final.Status)
	}
	if final.Error != "CustomStopError" || final.Cause != "stopped by the test" {
		t.Errorf("terminal pair = %q / %q, want the caller's CustomStopError / stopped by the test", final.Error, final.Cause)
	}

	history := readHistoryUntilTerminal(t, store, arn, "ExecutionAborted")
	last := history[len(history)-1]
	if d := last.ExecutionAbortedEventDetails; d == nil || d.Error != "CustomStopError" || d.Cause != "stopped by the test" {
		t.Errorf("ExecutionAborted details = %+v, want the caller's pair", last.ExecutionAbortedEventDetails)
	}
}

// TestStatesTaskFailedWildcardSemantics pins the documented wildcard:
// "The name States.TaskFailed also acts a wildcard and matches any error
// except for States.Timeout" — foreign error names match, the timeout does
// not, and exact-name matching keeps working.
func TestStatesTaskFailedWildcardSemantics(t *testing.T) {
	e := &Executor{}
	cases := []struct {
		errorCode, pattern string
		want               bool
	}{
		{"Lambda.TooManyRequestsException", "States.TaskFailed", true},
		{"HandledError", "States.TaskFailed", true},
		{"States.TaskFailed", "States.TaskFailed", true},
		{"States.Timeout", "States.TaskFailed", false},
		{"HandledError", "HandledError", true},
		{"HandledError", "OtherError", false},
	}
	for _, c := range cases {
		if got := e.errorMatchesPattern(c.errorCode, c.pattern); got != c.want {
			t.Errorf("errorMatchesPattern(%q, %q) = %v, want %v", c.errorCode, c.pattern, got, c.want)
		}
	}
}

// TestMapCatchMatchesIterationErrorIdentity pins that a Map-level Catch on
// the iteration's own error name matches: the Map state reports the
// iteration's error identity rather than an undocumented aggregate.
func TestMapCatchMatchesIterationErrorIdentity(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/mapcatch",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "mapcatch", Status: "RUNNING", Input: `{"v":[1]}`,
	}
	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.v",
				"ItemProcessor": map[string]interface{}{
					"StartAt": "F",
					"States": map[string]interface{}{
						"F": map[string]interface{}{"Type": "Fail", "Error": "CustomError", "Cause": "boom"},
					},
				},
				"Catch": []interface{}{
					map[string]interface{}{"ErrorEquals": []interface{}{"CustomError"}, "Next": "Recovery"},
				},
				"Next": "Done",
			},
			"Recovery": map[string]interface{}{"Type": "Pass", "Result": "caught", "End": true},
			"Done":     map[string]interface{}{"Type": "Pass", "Result": "uncought", "End": true},
		},
	}
	if err := store.CreateExecution(t.Context(), exec); err != nil {
		t.Fatalf("persist parent execution failed: %v", err)
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M", Input: exec.Input,
		EventId: ptrEventID(), States: map[string]sfnstore.State{}, QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	output, next, execErr := e.executeMap(t.Context(), execCtx, execCtx.States["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("the Catch on the iteration's error must handle the failure: %v", execErr.Cause)
	}
	if next != "Recovery" {
		t.Errorf("Catch transition Next = %q, want Recovery", next)
	}
	var caught map[string]interface{}
	if err := json.Unmarshal([]byte(output), &caught); err != nil {
		t.Fatalf("catch output not JSON: %v (%s)", err, output)
	}
	if caught["Error"] != "CustomError" {
		t.Errorf("catch output Error = %v, want CustomError — the identity must survive the Map wrap", caught["Error"])
	}
}

// modelHistoryEventTypes is the HistoryEventType enum of the Smithy model,
// transcribed by hand (test files never read the model asset). The wire
// contract this table pins: every event the engine serialises carries a
// type the model defines, and a details member under the exact model member
// name — the SDK drops unknown members, so a misspelt key erases the
// event's payload on the wire.
var modelHistoryEventTypes = map[string]bool{
	"ActivityFailed": true, "ActivityScheduled": true, "ActivityScheduleFailed": true,
	"ActivityStarted": true, "ActivitySucceeded": true, "ActivityTimedOut": true,
	"ChoiceStateEntered": true, "ChoiceStateExited": true,
	"ExecutionAborted": true, "ExecutionFailed": true, "ExecutionStarted": true,
	"ExecutionSucceeded": true, "ExecutionTimedOut": true, "ExecutionRedriven": true,
	"FailStateEntered":     true,
	"LambdaFunctionFailed": true, "LambdaFunctionScheduled": true,
	"LambdaFunctionScheduleFailed": true, "LambdaFunctionStarted": true,
	"LambdaFunctionStartFailed": true, "LambdaFunctionSucceeded": true,
	"LambdaFunctionTimedOut": true,
	"MapIterationAborted":    true, "MapIterationFailed": true, "MapIterationStarted": true,
	"MapIterationSucceeded": true,
	"MapStateAborted":       true, "MapStateEntered": true, "MapStateExited": true,
	"MapStateFailed": true, "MapStateStarted": true, "MapStateSucceeded": true,
	"ParallelStateAborted": true, "ParallelStateEntered": true, "ParallelStateExited": true,
	"ParallelStateFailed": true, "ParallelStateStarted": true, "ParallelStateSucceeded": true,
	"PassStateEntered": true, "PassStateExited": true,
	"SucceedStateEntered": true, "SucceedStateExited": true,
	"TaskFailed": true, "TaskScheduled": true, "TaskStarted": true, "TaskStartFailed": true,
	"TaskStateAborted": true, "TaskStateEntered": true, "TaskStateExited": true,
	"TaskSubmitFailed": true, "TaskSubmitted": true, "TaskSucceeded": true, "TaskTimedOut": true,
	"WaitStateAborted": true, "WaitStateEntered": true, "WaitStateExited": true,
	"MapRunAborted": true, "MapRunFailed": true, "MapRunStarted": true, "MapRunSucceeded": true,
	"MapRunRedriven": true, "EvaluationFailed": true,
}

// TestHistoryEventDetailsMemberKeys walks every detail-bearing event type
// through the serialiser and pins the response member key against the
// HistoryEvent member names of the model.
func TestHistoryEventDetailsMemberKeys(t *testing.T) {
	cases := []struct {
		eventType  string
		wantMember string
		populate   func(*sfnstore.ExecutionHistoryEvent)
	}{
		{"ExecutionStarted", "executionStartedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ExecutionStartedEventDetails = &sfnstore.ExecutionStartedEventDetails{}
		}},
		{"ExecutionSucceeded", "executionSucceededEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ExecutionSucceededEventDetails = &sfnstore.ExecutionSucceededEventDetails{}
		}},
		{"ExecutionFailed", "executionFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ExecutionFailedEventDetails = &sfnstore.ExecutionFailedEventDetails{}
		}},
		{"ExecutionAborted", "executionAbortedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ExecutionAbortedEventDetails = &sfnstore.ExecutionAbortedEventDetails{}
		}},
		{"ExecutionTimedOut", "executionTimedOutEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ExecutionTimedOutEventDetails = &sfnstore.ExecutionTimedOutEventDetails{}
		}},
		{"ExecutionRedriven", "executionRedrivenEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ExecutionRedrivenEventDetails = &sfnstore.ExecutionRedrivenEventDetails{RedriveCount: 2}
		}},
		{"TaskScheduled", "taskScheduledEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskScheduledEventDetails = &sfnstore.TaskScheduledEventDetails{}
		}},
		{"TaskStarted", "taskStartedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskStartedEventDetails = &sfnstore.TaskStartedEventDetails{}
		}},
		{"TaskStartFailed", "taskStartFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskStartFailedEventDetails = &sfnstore.TaskStartFailedEventDetails{}
		}},
		{"TaskSubmitted", "taskSubmittedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskSubmittedEventDetails = &sfnstore.TaskSubmittedEventDetails{}
		}},
		{"TaskSubmitFailed", "taskSubmitFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskSubmitFailedEventDetails = &sfnstore.TaskSubmitFailedEventDetails{}
		}},
		{"TaskSucceeded", "taskSucceededEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskSucceededEventDetails = &sfnstore.TaskSucceededEventDetails{}
		}},
		{"TaskFailed", "taskFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskFailedEventDetails = &sfnstore.TaskFailedEventDetails{}
		}},
		{"TaskTimedOut", "taskTimedOutEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.TaskTimedOutEventDetails = &sfnstore.TaskTimedOutEventDetails{}
		}},
		{"LambdaFunctionScheduled", "lambdaFunctionScheduledEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.LambdaFunctionScheduledEventDetails = &sfnstore.LambdaFunctionScheduledEventDetails{}
		}},
		{"LambdaFunctionScheduleFailed", "lambdaFunctionScheduleFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.LambdaFunctionScheduleFailedEventDetails = &sfnstore.LambdaFunctionScheduleFailedEventDetails{}
		}},
		{"LambdaFunctionStartFailed", "lambdaFunctionStartFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.LambdaFunctionStartFailedEventDetails = &sfnstore.LambdaFunctionStartFailedEventDetails{}
		}},
		{"LambdaFunctionFailed", "lambdaFunctionFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.LambdaFunctionFailedEventDetails = &sfnstore.LambdaFunctionFailedEventDetails{}
		}},
		{"LambdaFunctionTimedOut", "lambdaFunctionTimedOutEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.LambdaFunctionTimedOutEventDetails = &sfnstore.LambdaFunctionTimedOutEventDetails{}
		}},
		{"LambdaFunctionSucceeded", "lambdaFunctionSucceededEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.LambdaFunctionSucceededEventDetails = &sfnstore.LambdaFunctionSucceededEventDetails{}
		}},
		{"ActivityScheduled", "activityScheduledEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ActivityTaskScheduledEventDetails = &sfnstore.ActivityTaskScheduledEventDetails{}
		}},
		{"ActivityScheduleFailed", "activityScheduleFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ActivityScheduleFailedEventDetails = &sfnstore.ActivityScheduleFailedEventDetails{}
		}},
		{"ActivityStarted", "activityStartedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ActivityTaskStartedEventDetails = &sfnstore.ActivityTaskStartedEventDetails{}
		}},
		{"ActivitySucceeded", "activitySucceededEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ActivityTaskSucceededEventDetails = &sfnstore.ActivityTaskSucceededEventDetails{}
		}},
		{"ActivityFailed", "activityFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ActivityTaskFailedEventDetails = &sfnstore.ActivityTaskFailedEventDetails{}
		}},
		{"ActivityTimedOut", "activityTimedOutEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.ActivityTaskTimedOutEventDetails = &sfnstore.ActivityTaskTimedOutEventDetails{}
		}},
		{"PassStateEntered", "stateEnteredEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.StateEnteredEventDetails = &sfnstore.StateEnteredEventDetails{}
		}},
		{"PassStateExited", "stateExitedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.StateExitedEventDetails = &sfnstore.StateExitedEventDetails{}
		}},
		{"MapStateStarted", "mapStateStartedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapStateStartedEventDetails = &sfnstore.MapStateStartedEventDetails{}
		}},
		{"MapRunStarted", "mapRunStartedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapRunStartedEventDetails = &sfnstore.MapRunStartedEventDetails{}
		}},
		{"MapRunFailed", "mapRunFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapRunFailedEventDetails = &sfnstore.MapRunFailedEventDetails{}
		}},
		{"MapRunRedriven", "mapRunRedrivenEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapRunRedrivenEventDetails = &sfnstore.MapRunRedrivenEventDetails{MapRunArn: "arn", RedriveCount: 1}
		}},
		{"MapIterationStarted", "mapIterationStartedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapIterationEventDetails = &sfnstore.MapIterationEventDetails{}
		}},
		{"MapIterationSucceeded", "mapIterationSucceededEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapIterationEventDetails = &sfnstore.MapIterationEventDetails{}
		}},
		{"MapIterationFailed", "mapIterationFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapIterationEventDetails = &sfnstore.MapIterationEventDetails{}
		}},
		{"MapIterationAborted", "mapIterationAbortedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.MapIterationEventDetails = &sfnstore.MapIterationEventDetails{}
		}},
		{"EvaluationFailed", "evaluationFailedEventDetails", func(e *sfnstore.ExecutionHistoryEvent) {
			e.EvaluationFailedEventDetails = &sfnstore.EvaluationFailedEventDetails{}
		}},
	}
	for _, tc := range cases {
		ev := &sfnstore.ExecutionHistoryEvent{Type: tc.eventType}
		tc.populate(ev)
		resp := historyEventToResponse(ev, false)
		if !modelHistoryEventTypes[tc.eventType] {
			t.Errorf("%s: event type is not a HistoryEventType enum value of the model", tc.eventType)
		}
		if _, ok := resp[tc.wantMember]; !ok {
			t.Errorf("%s: response lacks the model member %q (got keys %v) — the SDK drops unknown members, so the details are erased on the wire",
				tc.eventType, tc.wantMember, respKeys(resp))
		}
	}
}

// respKeys renders a response map's keys for a failure message.
func respKeys(resp map[string]interface{}) []string {
	keys := make([]string, 0, len(resp))
	for k := range resp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TestHistoryEventLogsAtLevel pins the level gate against the CloudWatch
// Logs table in the developer guide: a level=ERROR machine emits only
// failure-class events, FATAL only the execution-terminal failures, OFF
// nothing, ALL everything.
func TestHistoryEventLogsAtLevel(t *testing.T) {
	cases := []struct {
		eventType       string
		all, err, fatal bool
	}{
		{"ExecutionStarted", true, false, false},
		{"ExecutionSucceeded", true, false, false},
		{"TaskStarted", true, false, false},
		{"TaskSucceeded", true, false, false},
		{"PassStateEntered", true, false, false},
		{"TaskStateExited", true, false, false},
		{"WaitStateExited", true, false, false},
		{"ActivityScheduled", true, false, false},
		{"ExecutionRedriven", true, false, false},
		{"FailStateEntered", true, true, false},
		{"TaskFailed", true, true, false},
		{"TaskTimedOut", true, true, false},
		{"TaskStateAborted", true, true, false},
		{"ActivityFailed", true, true, false},
		{"ActivityTimedOut", true, true, false},
		{"EvaluationFailed", true, true, false},
		{"MapStateFailed", true, true, false},
		{"ParallelStateAborted", true, true, false},
		{"LambdaFunctionScheduleFailed", true, true, false},
		{"ExecutionFailed", true, true, true},
		{"ExecutionAborted", true, true, true},
		{"ExecutionTimedOut", true, true, true},
	}
	for _, tc := range cases {
		if got := historyEventLogsAtLevel("ALL", tc.eventType); got != tc.all {
			t.Errorf("ALL/%s = %v, want %v", tc.eventType, got, tc.all)
		}
		if got := historyEventLogsAtLevel("ERROR", tc.eventType); got != tc.err {
			t.Errorf("ERROR/%s = %v, want %v", tc.eventType, got, tc.err)
		}
		if got := historyEventLogsAtLevel("FATAL", tc.eventType); got != tc.fatal {
			t.Errorf("FATAL/%s = %v, want %v", tc.eventType, got, tc.fatal)
		}
		if got := historyEventLogsAtLevel("OFF", tc.eventType); got {
			t.Errorf("OFF/%s = true, want false", tc.eventType)
		}
	}
}

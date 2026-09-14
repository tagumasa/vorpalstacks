package sfn

import (
	"context"
	"strings"
	"testing"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// scrubWalkEvents pairs every event type with detail members that
// historyEventToResponse handles, so the scrub test below covers the whole
// vocabulary rather than the types a hand-maintained key list happened to
// name. Event types whose HistoryEvent shape carries no detail member
// (ParallelStateStarted, MapRunSucceeded, the Aborted variants, …) have no
// details map to scrub and stay out of the walk.
func scrubWalkEvents() []*sfnstore.ExecutionHistoryEvent {
	return []*sfnstore.ExecutionHistoryEvent{
		{Type: "ExecutionStarted", ExecutionStartedEventDetails: &sfnstore.ExecutionStartedEventDetails{Input: `{"v":1}`, RoleArn: "arn:r", StateMachineAliasArn: "arn:a", StateMachineVersionArn: "arn:v"}},
		{Type: "ExecutionSucceeded", ExecutionSucceededEventDetails: &sfnstore.ExecutionSucceededEventDetails{Output: `{"v":1}`}},
		{Type: "ExecutionFailed", ExecutionFailedEventDetails: &sfnstore.ExecutionFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "ExecutionAborted", ExecutionAbortedEventDetails: &sfnstore.ExecutionAbortedEventDetails{Error: "E", Cause: "c"}},
		{Type: "ExecutionTimedOut", ExecutionTimedOutEventDetails: &sfnstore.ExecutionTimedOutEventDetails{Error: "E", Cause: "c"}},
		{Type: "ExecutionRedriven", ExecutionRedrivenEventDetails: &sfnstore.ExecutionRedrivenEventDetails{RedriveCount: 2}},
		{Type: "TaskScheduled", TaskScheduledEventDetails: &sfnstore.TaskScheduledEventDetails{Resource: "arn:r", ResourceType: "sqs", Region: "us-east-1", Parameters: `{"v":1}`, TimeoutInSeconds: 30, HeartbeatInSeconds: 5}},
		{Type: "TaskStarted", TaskStartedEventDetails: &sfnstore.TaskStartedEventDetails{Resource: "arn:r", ResourceType: "sqs"}},
		{Type: "TaskStartFailed", TaskStartFailedEventDetails: &sfnstore.TaskStartFailedEventDetails{Resource: "arn:r", ResourceType: "sqs", Error: "E", Cause: "c"}},
		{Type: "TaskSubmitted", TaskSubmittedEventDetails: &sfnstore.TaskSubmittedEventDetails{Resource: "arn:r", ResourceType: "sqs", Output: `{"v":1}`}},
		{Type: "TaskSubmitFailed", TaskSubmitFailedEventDetails: &sfnstore.TaskSubmitFailedEventDetails{Resource: "arn:r", ResourceType: "sqs", Error: "E", Cause: "c"}},
		{Type: "TaskSucceeded", TaskSucceededEventDetails: &sfnstore.TaskSucceededEventDetails{Resource: "arn:r", ResourceType: "sqs", Output: `{"v":1}`}},
		{Type: "TaskFailed", TaskFailedEventDetails: &sfnstore.TaskFailedEventDetails{Resource: "arn:r", ResourceType: "sqs", Error: "E", Cause: "c"}},
		{Type: "TaskTimedOut", TaskTimedOutEventDetails: &sfnstore.TaskTimedOutEventDetails{Resource: "arn:r", ResourceType: "sqs", Error: "E", Cause: "c"}},
		{Type: "LambdaFunctionScheduled", LambdaFunctionScheduledEventDetails: &sfnstore.LambdaFunctionScheduledEventDetails{Resource: "arn:r", Input: `{"v":1}`, TimeoutInSeconds: 30}},
		{Type: "LambdaFunctionScheduleFailed", LambdaFunctionScheduleFailedEventDetails: &sfnstore.LambdaFunctionScheduleFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "LambdaFunctionStartFailed", LambdaFunctionStartFailedEventDetails: &sfnstore.LambdaFunctionStartFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "LambdaFunctionFailed", LambdaFunctionFailedEventDetails: &sfnstore.LambdaFunctionFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "LambdaFunctionTimedOut", LambdaFunctionTimedOutEventDetails: &sfnstore.LambdaFunctionTimedOutEventDetails{Error: "E", Cause: "c"}},
		{Type: "LambdaFunctionSucceeded", LambdaFunctionSucceededEventDetails: &sfnstore.LambdaFunctionSucceededEventDetails{Output: `{"v":1}`}},
		{Type: "PassStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "P", Input: `{"v":1}`}},
		{Type: "PassStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "P", Output: `{"v":1}`}},
		{Type: "ChoiceStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "C", Input: `{"v":1}`}},
		{Type: "ChoiceStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "C", Output: `{"v":1}`, NextState: "N"}},
		{Type: "WaitStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "W", Input: `{"v":1}`}},
		{Type: "WaitStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "W", Output: `{"v":1}`}},
		{Type: "ParallelStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "Par", Input: `{"v":1}`}},
		{Type: "ParallelStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "Par", Output: `{"v":1}`}},
		{Type: "MapStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "M", Input: `{"v":1}`}},
		{Type: "MapStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "M", Output: `{"v":1}`}},
		{Type: "FailStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "F", Input: `{"v":1}`}},
		{Type: "SucceedStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "S", Input: `{"v":1}`}},
		{Type: "SucceedStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "S", Output: `{"v":1}`}},
		{Type: "TaskStateEntered", StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{Name: "T", Input: `{"v":1}`}},
		{Type: "TaskStateExited", StateExitedEventDetails: &sfnstore.StateExitedEventDetails{Name: "T", Output: `{"v":1}`}},
		{Type: "MapStateStarted", MapStateStartedEventDetails: &sfnstore.MapStateStartedEventDetails{Length: 3}},
		{Type: "MapRunStarted", MapRunStartedEventDetails: &sfnstore.MapRunStartedEventDetails{MapRunArn: "arn:mr"}},
		{Type: "MapRunFailed", MapRunFailedEventDetails: &sfnstore.MapRunFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "MapRunRedriven", MapRunRedrivenEventDetails: &sfnstore.MapRunRedrivenEventDetails{MapRunArn: "arn:mr", RedriveCount: 1}},
		{Type: "MapIterationStarted", MapIterationEventDetails: &sfnstore.MapIterationEventDetails{Index: 0, Name: "M"}},
		{Type: "MapIterationSucceeded", MapIterationEventDetails: &sfnstore.MapIterationEventDetails{Index: 0, Name: "M"}},
		{Type: "MapIterationFailed", MapIterationEventDetails: &sfnstore.MapIterationEventDetails{Index: 0, Name: "M"}},
		{Type: "MapIterationAborted", MapIterationEventDetails: &sfnstore.MapIterationEventDetails{Index: 0, Name: "M"}},
		{Type: "ActivityScheduled", ActivityTaskScheduledEventDetails: &sfnstore.ActivityTaskScheduledEventDetails{Resource: "arn:r", Input: `{"v":1}`, TimeoutInSeconds: 30, HeartbeatSeconds: 5}},
		{Type: "ActivityScheduleFailed", ActivityScheduleFailedEventDetails: &sfnstore.ActivityScheduleFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "ActivityStarted", ActivityTaskStartedEventDetails: &sfnstore.ActivityTaskStartedEventDetails{WorkerName: "w"}},
		{Type: "ActivitySucceeded", ActivityTaskSucceededEventDetails: &sfnstore.ActivityTaskSucceededEventDetails{Output: `{"v":1}`}},
		{Type: "ActivityFailed", ActivityTaskFailedEventDetails: &sfnstore.ActivityTaskFailedEventDetails{Error: "E", Cause: "c"}},
		{Type: "ActivityTimedOut", ActivityTaskTimedOutEventDetails: &sfnstore.ActivityTaskTimedOutEventDetails{Error: "E", Cause: "c"}},
		{Type: "EvaluationFailed", EvaluationFailedEventDetails: &sfnstore.EvaluationFailedEventDetails{State: "S", Cause: "c", Error: "E", Location: "loc"}},
	}
}

func assertNoExecutionDataKeys(t *testing.T, prefix string, v map[string]interface{}) {
	t.Helper()
	for k, val := range v {
		if executionDataSensitiveFields[k] {
			t.Errorf("%s: execution data key %q survived includeExecutionData=false", prefix, k)
		}
		if nested, ok := val.(map[string]interface{}); ok {
			assertNoExecutionDataKeys(t, prefix+"."+k, nested)
		}
	}
}

// TestHistoryEventScrubCoversEveryEventType walks the full event
// vocabulary: under includeExecutionData=false no input/output member may
// survive anywhere in the serialised event, and every event type must
// serialise its details map — a bare event means the case is missing from
// the switch entirely.
func TestHistoryEventScrubCoversEveryEventType(t *testing.T) {
	for _, ev := range scrubWalkEvents() {
		t.Run(ev.Type, func(t *testing.T) {
			resp := historyEventToResponse(ev, false)
			assertNoExecutionDataKeys(t, ev.Type, resp)

			wantDetails := strings.ToLower(ev.Type[:1]) + ev.Type[1:] + "EventDetails"
			switch {
			case strings.HasSuffix(ev.Type, "StateEntered"):
				// State-entry events share the generic member.
				wantDetails = "stateEnteredEventDetails"
			case strings.HasSuffix(ev.Type, "StateExited"):
				wantDetails = "stateExitedEventDetails"
			}
			if _, ok := resp[wantDetails].(map[string]interface{}); !ok {
				t.Fatalf("%s must serialise %s as a map, got %v", ev.Type, wantDetails, resp[wantDetails])
			}

			// The Choice exit's taken transition is an internal resume
			// marker, never a wire member.
			if ev.Type == "ChoiceStateExited" {
				if _, ok := resp["stateExitedEventDetails"].(map[string]interface{})["nextState"]; ok {
					t.Error("nextState is an internal resume marker and must not serialise")
				}
			}

			// Positive control on one representative: with execution data
			// included the members are present, so the scrub above cannot
			// pass vacuously.
			if ev.Type == "ExecutionStarted" {
				full := historyEventToResponse(ev, true)
				details, ok := full["executionStartedEventDetails"].(map[string]interface{})
				if !ok || details["input"] != `{"v":1}` {
					t.Errorf("includeExecutionData=true must keep the input, got %v", full["executionStartedEventDetails"])
				}
			}
		})
	}
}

// TestActivityScheduledEventMembersMatchModel pins the member set of the
// serialised activityScheduledEventDetails to the model's
// ActivityScheduledEventDetails: resource, input, inputDetails,
// timeoutInSeconds, heartbeatInSeconds — the task token is not a member;
// it travels in GetActivityTask's response.
func TestActivityScheduledEventMembersMatchModel(t *testing.T) {
	ev := &sfnstore.ExecutionHistoryEvent{
		Type: "ActivityScheduled",
		ActivityTaskScheduledEventDetails: &sfnstore.ActivityTaskScheduledEventDetails{
			Resource: "arn:r", Input: `{"v":1}`, TimeoutInSeconds: 30, HeartbeatSeconds: 5,
		},
	}
	details, ok := historyEventToResponse(ev, true)["activityScheduledEventDetails"].(map[string]interface{})
	if !ok {
		t.Fatal("activityScheduledEventDetails did not serialise as a map")
	}
	want := []string{"resource", "input", "inputDetails", "timeoutInSeconds", "heartbeatInSeconds"}
	if len(details) != len(want) {
		t.Errorf("member count = %d (%v), want %d", len(details), details, len(want))
	}
	for _, k := range want {
		if _, ok := details[k]; !ok {
			t.Errorf("missing member %q in %v", k, details)
		}
	}
	if _, ok := details["taskToken"]; ok {
		t.Error("taskToken is not a model member of activityScheduledEventDetails")
	}
}

// TestHistoryEventDataDetailsShape pins the details member shape:
// HistoryEventExecutionDataDetails and AssignedVariablesDetails carry one
// member, truncated, "Always false for API calls" — the flag does not
// track includeExecutionData, which governs omission of the payload
// members alone.
func TestHistoryEventDataDetailsShape(t *testing.T) {
	entered := &sfnstore.ExecutionHistoryEvent{
		Type: "TaskStateEntered",
		StateEnteredEventDetails: &sfnstore.StateEnteredEventDetails{
			Input: `{"v":1}`, Name: "T",
		},
	}
	resp := historyEventToResponse(entered, false)
	details, ok := resp["stateEnteredEventDetails"].(map[string]interface{})
	if !ok {
		t.Fatal("stateEnteredEventDetails did not serialise")
	}
	if _, present := details["input"]; present {
		t.Error("includeExecutionData=false must omit the input member")
	}
	in, ok := details["inputDetails"].(map[string]interface{})
	if !ok || len(in) != 1 || in["truncated"] != false {
		t.Errorf("inputDetails = %v, want {truncated:false}", details["inputDetails"])
	}

	exited := &sfnstore.ExecutionHistoryEvent{
		Type: "TaskStateExited",
		StateExitedEventDetails: &sfnstore.StateExitedEventDetails{
			Output: `{"v":1}`, Name: "T",
			AssignedVariables: map[string]interface{}{"$x": 1},
		},
	}
	resp = historyEventToResponse(exited, true)
	details, ok = resp["stateExitedEventDetails"].(map[string]interface{})
	if !ok {
		t.Fatal("stateExitedEventDetails did not serialise")
	}
	if out, ok := details["outputDetails"].(map[string]interface{}); !ok || len(out) != 1 || out["truncated"] != false {
		t.Errorf("outputDetails = %v, want {truncated:false}", details["outputDetails"])
	}
	vars, ok := details["assignedVariablesDetails"].(map[string]interface{})
	if !ok {
		t.Fatal("assignedVariablesDetails did not serialise")
	}
	for name, d := range vars {
		detail, ok := d.(map[string]interface{})
		if !ok || len(detail) != 1 || detail["truncated"] != false {
			t.Errorf("assignedVariablesDetails[%s] = %v, want {truncated:false}", name, d)
		}
	}
}

// TestExecutionResponseDetailsShape pins the execution-level details:
// CloudWatchEventsExecutionDataDetails carries one member, included,
// "Always true for API calls" — no type member, and the flag does not
// track whether the execution carried input.
func TestExecutionResponseDetailsShape(t *testing.T) {
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:e1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "e1", Status: "RUNNING",
	}
	store := newCreateTestStore(t)
	resp := executionToResponse(context.Background(), store, exec)
	in, ok := resp["inputDetails"].(map[string]interface{})
	if !ok || len(in) != 1 || in["included"] != true {
		t.Errorf("inputDetails = %v, want {included:true} even without input", resp["inputDetails"])
	}
	if _, present := resp["outputDetails"]; present {
		t.Error("outputDetails must be absent while the execution has no output")
	}

	exec.Output = `{"done":true}`
	resp = executionToResponse(context.Background(), store, exec)
	out, ok := resp["outputDetails"].(map[string]interface{})
	if !ok || len(out) != 1 || out["included"] != true {
		t.Errorf("outputDetails = %v, want {included:true}", resp["outputDetails"])
	}
}

// TestWireTimestampsKeepSubSecondPrecision pins the wire rendering of
// timestamps: Step Functions timestamps carry sub-second digits, so two
// events milliseconds apart must not collapse onto one whole-second
// value.
func TestWireTimestampsKeepSubSecondPrecision(t *testing.T) {
	base := time.Unix(1700000000, 0).UTC()
	first := &sfnstore.ExecutionHistoryEvent{
		EventId:   1,
		Type:      "ExecutionStarted",
		Timestamp: base,
	}
	second := &sfnstore.ExecutionHistoryEvent{
		EventId:   2,
		Type:      "ExecutionStarted",
		Timestamp: base.Add(5 * time.Millisecond),
	}
	a := historyEventToResponse(first, true)["timestamp"].(float64)
	b := historyEventToResponse(second, true)["timestamp"].(float64)
	if a == b {
		t.Errorf("5ms apart collapsed to one timestamp: %v", a)
	}
	if got := b - a; got < 0.004 {
		t.Errorf("timestamp delta = %v, want the 5ms difference within float precision", got)
	}
}

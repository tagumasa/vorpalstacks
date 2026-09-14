package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// stubEventsInvoker records the stored putEvents entries; failOn marks the
// Nth store call (1-based) to fail so per-entry result handling stays
// observable.
type stubEventsInvoker struct {
	stored  []map[string]interface{}
	failOn  int
	calls   int
	publish func(event map[string]interface{})
}

func (s *stubEventsInvoker) PutEvent(_ context.Context, _ string, event any) error {
	s.calls++
	if s.failOn > 0 && s.calls == s.failOn {
		return errors.New("storage unavailable")
	}
	m, _ := event.(map[string]interface{})
	s.stored = append(s.stored, m)
	return nil
}

func newEventsTestExecutor(t *testing.T, stub *stubEventsInvoker) *Executor {
	t.Helper()
	bus := eventbus.NewEventBus()
	bus.SetEventsInvoker(stub)
	e := NewExecutor(nil, bus)
	e.region = "us-east-1"
	e.accountID = "000000000000"
	return e
}

func eventsTestContext() *ExecutionContext {
	return &ExecutionContext{
		Execution: &sfnstore.Execution{
			ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:e1",
			StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		},
	}
}

// TestEventsPutEventsResourcesAndTime pins the documented entry contract:
// "The execution ARN and the state machine ARN are automatically appended
// to the Resources field of each PutEventsRequestEntry" — a user-supplied
// Resources array survives ahead of the appended ARNs, and a supplied
// entry Time is honoured.
func TestEventsPutEventsResourcesAndTime(t *testing.T) {
	stub := &stubEventsInvoker{}
	e := newEventsTestExecutor(t, stub)

	input := `{"Entries":[{"Source":"my.app","DetailType":"x","Detail":{"a":1},` +
		`"Resources":["arn:aws:s3:::my-bucket"],"Time":"2026-01-02T03:04:05Z"}]}`
	out, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", input)
	if err != nil {
		t.Fatalf("putEvents failed: %v", err)
	}

	if len(stub.stored) != 1 {
		t.Fatalf("stored events = %d, want 1", len(stub.stored))
	}
	resources, _ := stub.stored[0]["Resources"].([]string)
	if len(resources) != 3 || resources[0] != "arn:aws:s3:::my-bucket" ||
		resources[1] != "arn:aws:states:us-east-1:000000000000:execution:sm:e1" ||
		resources[2] != "arn:aws:states:us-east-1:000000000000:stateMachine:sm" {
		t.Errorf("resources = %v, want the user entry first and the ARNs appended", resources)
	}
	storedTime, _ := stub.stored[0]["Time"].(interface{ String() string })
	if storedTime == nil || !strings.HasPrefix(storedTime.String(), "2026-01-02 03:04:05") {
		t.Errorf("event Time = %v, want the entry-supplied timestamp", stub.stored[0]["Time"])
	}

	var result map[string]interface{}
	if jerr := json.Unmarshal([]byte(out), &result); jerr != nil {
		t.Fatalf("result not JSON: %v (%s)", jerr, out)
	}
	if result["FailedEntryCount"].(float64) != 0 {
		t.Errorf("FailedEntryCount = %v, want 0", result["FailedEntryCount"])
	}
}

// TestEventsPutEventsFailedEntry pins the error contract: "Step Functions
// checks whether the FailedEntryCount is greater than zero. If it is
// greater than zero, Step Functions fails the state with the error
// EventBridge.FailedEntry" — retry and Catch on that name must be able to
// match.
func TestEventsPutEventsFailedEntry(t *testing.T) {
	stub := &stubEventsInvoker{failOn: 2}
	e := newEventsTestExecutor(t, stub)

	input := `{"Entries":[{"Source":"a","DetailType":"x"},{"Source":"b","DetailType":"x"},{"Source":"c","DetailType":"x"}]}`
	_, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", input)
	if err == nil {
		t.Fatal("a failed entry must fail the task")
	}
	var tf *taskFailure
	if !errors.As(err, &tf) {
		t.Fatalf("error = %v, want the typed task failure", err)
	}
	if tf.code != "EventBridge.FailedEntry" {
		t.Errorf("failure code = %q, want EventBridge.FailedEntry", tf.code)
	}
	if !strings.Contains(tf.cause, "1 of 3") {
		t.Errorf("failure cause = %q, want the per-entry counts", tf.cause)
	}
}

// TestEventsPutEventsRequiresEntries pins that a payload without the
// entries array is an invocation failure: putEvents publishes exactly the
// entries it is given and never invents an event.
func TestEventsPutEventsRequiresEntries(t *testing.T) {
	stub := &stubEventsInvoker{}
	e := newEventsTestExecutor(t, stub)

	_, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", `{"Source":"my.app","DetailType":"t","Detail":{"a":1}}`)
	if err == nil {
		t.Fatal("putEvents without Entries must fail")
	}
	if len(stub.stored) != 0 {
		t.Fatalf("no event may be invented: stored = %v", stub.stored)
	}
}

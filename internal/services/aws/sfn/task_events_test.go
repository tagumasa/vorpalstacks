package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// stubPutEventsHandler records the ingress events the putEvents integration
// publishes; failOn marks the Nth publish (1-based) to fail so per-entry
// result handling stays observable. The handler answers with the payload id
// the real EventBridge handler returns (the assigned event id).
type stubPutEventsHandler struct {
	mu     sync.Mutex
	events []*eventbus.EventBridgePutEventsEvent
	failOn int
	calls  int
}

func (h *stubPutEventsHandler) handle(_ context.Context, evt *eventbus.EventBridgePutEventsEvent) eventbus.HandlerResult {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.calls++
	if h.failOn > 0 && h.calls == h.failOn {
		return eventbus.HandlerResult{Error: errors.New("delivery unavailable")}
	}
	h.events = append(h.events, evt)
	return eventbus.HandlerResult{Payload: []byte(fmt.Sprintf("id-%d", h.calls))}
}

func newEventsTestExecutor(t *testing.T, stub *stubPutEventsHandler) *Executor {
	t.Helper()
	bus := eventbus.NewEventBus()
	// The subscription stays synchronous so PublishSync returns the stub's
	// verdict directly, the way the real handler's does.
	if _, err := eventbus.SubscribeTyped[*eventbus.EventBridgePutEventsEvent](bus, stub.handle); err != nil {
		t.Fatal(err)
	}
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
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

// TestEventsPutEventsRidesTheIngressPlane pins the unified-ingress
// contract: each entry is published as an internal putEvents bus event
// carrying the entry's members — the documented Resources append first,
// the entry-supplied Time — and the result's EventId is the id the handler
// assigned.
func TestEventsPutEventsRidesTheIngressPlane(t *testing.T) {
	stub := &stubPutEventsHandler{}
	e := newEventsTestExecutor(t, stub)

	input := `{"Entries":[{"Source":"my.app","DetailType":"x","Detail":{"a":1},` +
		`"Resources":["arn:aws:s3:::my-bucket"],"Time":"2026-01-02T03:04:05Z"}]}`
	out, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", input)
	if err != nil {
		t.Fatalf("putEvents failed: %v", err)
	}

	if len(stub.events) != 1 {
		t.Fatalf("published ingress events = %d, want 1", len(stub.events))
	}
	published := stub.events[0]
	if published.EventBusName != "default" {
		t.Errorf("event bus = %q, want the default bus", published.EventBusName)
	}
	if published.Region != "us-east-1" {
		t.Errorf("region = %q, want the resource ARN's region", published.Region)
	}
	var entry map[string]interface{}
	if jerr := json.Unmarshal([]byte(published.Input), &entry); jerr != nil {
		t.Fatalf("ingress input not JSON: %v (%s)", jerr, published.Input)
	}
	if entry["Source"] != "my.app" || entry["DetailType"] != "x" {
		t.Errorf("ingress entry = %v, want the entry-supplied Source and DetailType", entry)
	}
	if detail, _ := entry["Detail"].(map[string]interface{}); detail["a"] != float64(1) {
		t.Errorf("ingress detail = %v, want the entry-supplied detail object", entry["Detail"])
	}
	resources, _ := entry["Resources"].([]interface{})
	if len(resources) != 3 || resources[0] != "arn:aws:s3:::my-bucket" ||
		resources[1] != "arn:aws:states:us-east-1:000000000000:execution:sm:e1" ||
		resources[2] != "arn:aws:states:us-east-1:000000000000:stateMachine:sm" {
		t.Errorf("resources = %v, want the user entry first and the ARNs appended", resources)
	}
	if entry["Time"] != "2026-01-02T03:04:05Z" {
		t.Errorf("ingress time = %v, want the entry-supplied timestamp", entry["Time"])
	}

	var result map[string]interface{}
	if jerr := json.Unmarshal([]byte(out), &result); jerr != nil {
		t.Fatalf("result not JSON: %v (%s)", jerr, out)
	}
	if result["FailedEntryCount"].(float64) != 0 {
		t.Errorf("FailedEntryCount = %v, want 0", result["FailedEntryCount"])
	}
	entries, _ := result["Entries"].([]interface{})
	if len(entries) != 1 {
		t.Fatalf("result entries = %v, want one", result["Entries"])
	}
	first, _ := entries[0].(map[string]interface{})
	if first["EventId"] != "id-1" {
		t.Errorf("EventId = %v, want the handler-assigned id", first["EventId"])
	}
}

// TestEventsPutEventsFailedEntry pins the error contract: "Step Functions
// checks whether the FailedEntryCount is greater than zero. If it is
// greater than zero, Step Functions fails the state with the error
// EventBridge.FailedEntry" — retry and Catch on that name must be able to
// match.
func TestEventsPutEventsFailedEntry(t *testing.T) {
	stub := &stubPutEventsHandler{failOn: 2}
	e := newEventsTestExecutor(t, stub)

	input := `{"Entries":[{"Source":"a","DetailType":"x","Detail":{}},{"Source":"b","DetailType":"x","Detail":{}},{"Source":"c","DetailType":"x","Detail":{}}]}`
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

// TestEventsPutEventsForwardsMembersVerbatim pins the verbatim-forwarding
// contract: the entry's Time (string or epoch number, parseable or not) and
// Detail (object, string, or the null literal) ride the ingress event in
// their wire form — the EventBridge handler owns the member semantics, so
// pre-normalising here would silence its rejections.
func TestEventsPutEventsForwardsMembersVerbatim(t *testing.T) {
	stub := &stubPutEventsHandler{}
	e := newEventsTestExecutor(t, stub)

	input := `{"Entries":[` +
		`{"Source":"n","DetailType":"x","Detail":null},` +
		`{"Source":"e","DetailType":"x","Detail":{},"Time":1735689600},` +
		`{"Source":"b","DetailType":"x","Detail":{"k":"v"},"Time":"not-a-timestamp"},` +
		`{"Source":"s","DetailType":"x","Detail":"{\"string\":\"form\"}"}]}`
	out, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", input)
	if err != nil {
		t.Fatalf("putEvents failed: %v", err)
	}

	if len(stub.events) != 4 {
		t.Fatalf("published ingress events = %d, want 4", len(stub.events))
	}
	var n0 map[string]interface{}
	if jerr := json.Unmarshal([]byte(stub.events[0].Input), &n0); jerr != nil {
		t.Fatalf("ingress input not JSON: %v (%s)", jerr, stub.events[0].Input)
	}
	if _, present := n0["Detail"]; !present {
		t.Errorf("the null Detail must ride the ingress entry verbatim, input = %s", stub.events[0].Input)
	} else if n0["Detail"] != nil {
		t.Errorf("Detail = %v, want the JSON null forwarded as-is", n0["Detail"])
	}

	var n1 map[string]interface{}
	if jerr := json.Unmarshal([]byte(stub.events[1].Input), &n1); jerr != nil {
		t.Fatalf("ingress input not JSON: %v", jerr)
	}
	if n1["Time"] != float64(1735689600) {
		t.Errorf("Time = %v, want the epoch number forwarded verbatim", n1["Time"])
	}

	var n2 map[string]interface{}
	if jerr := json.Unmarshal([]byte(stub.events[2].Input), &n2); jerr != nil {
		t.Fatalf("ingress input not JSON: %v", jerr)
	}
	if n2["Time"] != "not-a-timestamp" {
		t.Errorf("Time = %v, want the unparseable timestamp forwarded verbatim for the handler to reject", n2["Time"])
	}

	var result map[string]interface{}
	if jerr := json.Unmarshal([]byte(out), &result); jerr != nil {
		t.Fatalf("result not JSON: %v (%s)", jerr, out)
	}
	if result["FailedEntryCount"].(float64) != 0 {
		t.Errorf("FailedEntryCount = %v, want 0 (the stub applies no member semantics)", result["FailedEntryCount"])
	}
}

// TestEventsPutEventsRequiresDetailMember pins the incomplete-entry rule on
// the SFN plane: PutEvents requires Source, DetailType and Detail on every
// entry, so an entry without a Detail member is a per-entry failure and no
// ingress event is invented for it.
func TestEventsPutEventsRequiresDetailMember(t *testing.T) {
	stub := &stubPutEventsHandler{}
	e := newEventsTestExecutor(t, stub)

	input := `{"Entries":[{"Source":"a","DetailType":"x"},{"Source":"b","DetailType":"x","Detail":{}}]}`
	_, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", input)
	if err == nil {
		t.Fatal("an entry without Detail must fail the task")
	}
	var tf *taskFailure
	if !errors.As(err, &tf) {
		t.Fatalf("error = %v, want the typed task failure", err)
	}
	if tf.code != "EventBridge.FailedEntry" {
		t.Errorf("failure code = %q, want EventBridge.FailedEntry", tf.code)
	}
	if !strings.Contains(tf.cause, "Source, DetailType, and Detail are required") {
		t.Errorf("failure cause = %q, want the incomplete-entry reason", tf.cause)
	}
	if len(stub.events) != 1 {
		t.Fatalf("only the complete entry may publish, published = %d", len(stub.events))
	}
}

// TestEventsPutEventsRequiresEntries pins that a payload without the
// entries array is an invocation failure: putEvents publishes exactly the
// entries it is given and never invents an event.
func TestEventsPutEventsRequiresEntries(t *testing.T) {
	stub := &stubPutEventsHandler{}
	e := newEventsTestExecutor(t, stub)

	_, err := e.executeEventsTask(context.Background(), eventsTestContext(),
		"arn:aws:states:::events:putEvents", `{"Source":"my.app","DetailType":"t","Detail":{"a":1}}`)
	if err == nil {
		t.Fatal("putEvents without Entries must fail")
	}
	if len(stub.events) != 0 {
		t.Fatalf("no event may be invented: published = %v", stub.events)
	}
}

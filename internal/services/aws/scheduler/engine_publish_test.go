package scheduler

import (
	"context"
	"testing"

	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// TestStartStepFunctionExecutionPropagatesHandlerError pins that a
// Step Functions execution-start failure reported by the bus handler
// reaches the caller as an error: the PublishSync contract carries
// handler failures in HandlerResult.Error with a nil error, and the
// delivery must be treated as failed so the retry policy and dead-letter
// routing engage.
func TestStartStepFunctionExecutionPropagatesHandlerError(t *testing.T) {
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	if _, err := bus.Subscribe(func(ctx context.Context, event eventbus.Event) eventbus.HandlerResult {
		return eventbus.HandlerResult{Error: context.DeadlineExceeded}
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	e := &Engine{bus: bus, accountID: "000000000000"}
	schedule := &schedulerstore.Schedule{Name: "sfn-fail", GroupName: "default", Region: "us-east-1"}
	target := &schedulerstore.Target{Arn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm"}
	if err := e.startStepFunctionExecution(context.Background(), schedule, target); err == nil {
		t.Error("startStepFunctionExecution returned nil for a failing handler")
	}
}

// TestSendToEventBridgePropagatesHandlerError pins the same contract for
// the EventBridge PutEvents delivery path.
func TestSendToEventBridgePropagatesHandlerError(t *testing.T) {
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	if _, err := bus.Subscribe(func(ctx context.Context, event eventbus.Event) eventbus.HandlerResult {
		return eventbus.HandlerResult{Error: context.DeadlineExceeded}
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	e := &Engine{bus: bus, accountID: "000000000000"}
	schedule := &schedulerstore.Schedule{Name: "eb-fail", GroupName: "default", Region: "us-east-1"}
	target := &schedulerstore.Target{Arn: "arn:aws:events:us-east-1:000000000000:event-bus/default"}
	if err := e.sendToEventBridge(context.Background(), schedule, target); err == nil {
		t.Error("sendToEventBridge returned nil for a failing handler")
	}
}

// TestSendToEventBridgeUsesNamedBusFromARN pins the bus-name extraction:
// a schedule targeting a custom event bus publishes to THAT bus. The ARN's
// resource segment is "event-bus/<name>" (no leading colon), so the
// published name must come from a prefix extractor over that segment.
// Rule ARNs are tolerated by the shared extractor (a fully qualified
// rule ARN names its own bus) — the pin holds the bus ARN form.
func TestSendToEventBridgeUsesNamedBusFromARN(t *testing.T) {
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	events := make(chan *eventbus.EventBridgePutEventsEvent, 1)
	if _, err := eventbus.SubscribeTyped[*eventbus.EventBridgePutEventsEvent](bus,
		func(ctx context.Context, evt *eventbus.EventBridgePutEventsEvent) eventbus.HandlerResult {
			select {
			case events <- evt:
			default:
			}
			return eventbus.HandlerResult{}
		}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	e := &Engine{bus: bus, accountID: "000000000000"}
	schedule := &schedulerstore.Schedule{Name: "named-bus", GroupName: "default", Region: "us-east-1"}
	target := &schedulerstore.Target{Arn: "arn:aws:events:us-east-1:000000000000:event-bus/vocab-bus"}
	if err := e.sendToEventBridge(context.Background(), schedule, target); err != nil {
		t.Fatalf("sendToEventBridge: %v", err)
	}
	select {
	case evt := <-events:
		if evt.EventBusName != "vocab-bus" {
			t.Errorf("published EventBusName = %q, want vocab-bus", evt.EventBusName)
		}
	default:
		t.Fatal("no EventBridgePutEventsEvent was published")
	}
}

// TestSendToEventBridgeRejectsNonBusARN pins the non-bus posture: an
// events ARN naming something other than an event bus (an archive here)
// is a delivery error feeding the retry/DLQ lifecycle, never a silent
// hop to the default bus.
func TestSendToEventBridgeRejectsNonBusARN(t *testing.T) {
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	e := &Engine{bus: bus, accountID: "000000000000"}
	schedule := &schedulerstore.Schedule{Name: "non-bus", GroupName: "default", Region: "us-east-1"}
	target := &schedulerstore.Target{Arn: "arn:aws:events:us-east-1:000000000000:archive/vocab-archive"}
	if err := e.sendToEventBridge(context.Background(), schedule, target); err == nil {
		t.Error("non-bus events ARN accepted, want a delivery error")
	}
}

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

	bus.Subscribe(func(ctx context.Context, event eventbus.Event) eventbus.HandlerResult {
		return eventbus.HandlerResult{Error: context.DeadlineExceeded}
	})

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

	bus.Subscribe(func(ctx context.Context, event eventbus.Event) eventbus.HandlerResult {
		return eventbus.HandlerResult{Error: context.DeadlineExceeded}
	})

	e := &Engine{bus: bus, accountID: "000000000000"}
	schedule := &schedulerstore.Schedule{Name: "eb-fail", GroupName: "default", Region: "us-east-1"}
	target := &schedulerstore.Target{Arn: "arn:aws:events:us-east-1:000000000000:event-bus/default"}
	if err := e.sendToEventBridge(context.Background(), schedule, target); err == nil {
		t.Error("sendToEventBridge returned nil for a failing handler")
	}
}

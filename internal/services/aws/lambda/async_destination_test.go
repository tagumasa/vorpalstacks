package lambda

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/eventbus"
)

// TestDeliverToEventBridgePublishesTheDocumentedEnvelope pins the
// EventBridge destination contract from the "Capturing records of Lambda
// asynchronous invocations" destination table: the invocation record rides
// as the detail, the source is lambda, the detail-type carries the
// success/failure suffix, and resources holds the function and destination
// ARNs. Delivery reaches the bus's putEvents subscribers — the unified
// ingress the EventBridge service consumes.
func TestDeliverToEventBridgePublishesTheDocumentedEnvelope(t *testing.T) {
	var mu sync.Mutex
	var got []*eventbus.EventBridgePutEventsEvent
	bus := eventbus.NewEventBus()
	// WithAsync matches the EventBridge service's own registration, which
	// is the subscription the production publish reaches.
	if _, err := eventbus.SubscribeTyped[*eventbus.EventBridgePutEventsEvent](bus, func(_ context.Context, evt *eventbus.EventBridgePutEventsEvent) eventbus.HandlerResult {
		mu.Lock()
		defer mu.Unlock()
		got = append(got, evt)
		return eventbus.HandlerResult{}
	}, eventbus.WithAsync()); err != nil {
		t.Fatal(err)
	}
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	s := &LambdaService{bus: bus}

	record := `{"version":"1.0","requestPayload":{"k":"v"}}`
	destARN := "arn:aws:events:us-east-1:000000000000:event-bus/dest-bus"
	functionArn := "arn:aws:lambda:us-east-1:000000000000:function:fn"

	waitFor := func(n int) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		mu.Lock()
		for len(got) < n && time.Now().Before(deadline) {
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
		}
		if len(got) < n {
			mu.Unlock()
			t.Fatalf("published events = %d, want %d", len(got), n)
		}
		mu.Unlock()
	}

	deliverToEventBridge(context.Background(), s, destARN, record, false, functionArn, "us-east-1")
	waitFor(1)

	mu.Lock()
	first := got[0]
	mu.Unlock()
	if first.EventBusName != "dest-bus" || first.Region != "us-east-1" {
		t.Fatalf("bus/region = %q/%q, want dest-bus/us-east-1", first.EventBusName, first.Region)
	}
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(first.Input), &input); err != nil {
		t.Fatalf("input not JSON: %v", err)
	}
	if input["Source"] != "lambda" {
		t.Errorf("source = %v, want lambda", input["Source"])
	}
	if input["DetailType"] != "Lambda Function Invocation Result - Failure" {
		t.Errorf("detail-type = %v, want the failure suffix", input["DetailType"])
	}
	detail, _ := input["Detail"].(map[string]interface{})
	if detail["version"] != "1.0" {
		t.Errorf("detail = %v, want the invocation record", input["Detail"])
	}
	resources, _ := input["Resources"].([]interface{})
	if len(resources) != 2 || resources[0] != functionArn || resources[1] != destARN {
		t.Errorf("resources = %v, want the function and destination ARNs", resources)
	}

	deliverToEventBridge(context.Background(), s, destARN, record, true, functionArn, "us-east-1")
	waitFor(2)
	mu.Lock()
	second := got[1]
	mu.Unlock()
	var okInput map[string]interface{}
	if err := json.Unmarshal([]byte(second.Input), &okInput); err != nil {
		t.Fatalf("input not JSON: %v", err)
	}
	if okInput["DetailType"] != "Lambda Function Invocation Result - Success" {
		t.Errorf("detail-type = %v, want the success suffix", okInput["DetailType"])
	}
}

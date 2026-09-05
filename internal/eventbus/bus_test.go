package eventbus

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPublishSyncNilEvent(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	_, err := bus.PublishSync(context.Background(), nil)
	if err != ErrNilEvent {
		t.Fatalf("expected ErrNilEvent, got %v", err)
	}
}

func TestPublishSyncEmptyType(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	evt := &EventBase{ID: "test-1", Timestamp: time.Now().UTC(), Source: "test", Region: "us-east-1"}
	_, err := bus.PublishSync(context.Background(), evt)
	if err != ErrEmptyType {
		t.Fatalf("expected ErrEmptyType, got %v", err)
	}
}

func TestPublishSyncDepthExceeded(t *testing.T) {
	bus := NewEventBus(WithMaxEventDepth(2))
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	evt := &SNSDeliveryEvent{
		EventBase: EventBase{ID: "depth-test", Timestamp: time.Now().UTC(), Source: "test", Region: "us-east-1"},
		TopicARN:  "arn:aws:sns:us-east-1:123456789012:test",
	}
	evt.EventBase.SetEventDepth(1)

	_, err := bus.PublishSync(context.Background(), evt)
	if err != ErrMaxDepth {
		t.Fatalf("expected ErrMaxDepth, got %v", err)
	}
}

func TestPublishSyncBasic(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	called := atomic.Int32{}
	subID, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		called.Add(1)
		return HandlerResult{StatusCode: 200}
	})
	if err != nil {
		t.Fatal(err)
	}

	evt := &SNSDeliveryEvent{
		EventBase: EventBase{ID: "basic-test", Timestamp: time.Now().UTC(), Source: "test", Region: "us-east-1"},
		TopicARN:  "arn:aws:sns:us-east-1:123456789012:test",
	}

	result, err := bus.PublishSync(context.Background(), evt)
	if err != nil {
		t.Fatal(err)
	}
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if called.Load() != 1 {
		t.Fatalf("expected handler called 1 time, got %d", called.Load())
	}

	if err := bus.Unsubscribe(subID); err != nil {
		t.Fatal(err)
	}
}

func TestPublishSyncHandlerError(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		return HandlerResult{Error: context.DeadlineExceeded}
	})

	evt := &SNSDeliveryEvent{
		EventBase: EventBase{ID: "err-test", Timestamp: time.Now().UTC(), Source: "test", Region: "us-east-1"},
		TopicARN:  "arn:aws:sns:us-east-1:123456789012:test",
	}

	result, err := bus.PublishSync(context.Background(), evt)
	if err != nil {
		t.Fatal("PublishSync should not return error when handler errors, should return HandlerResult with Error")
	}
	if result.Error == nil {
		t.Fatal("expected handler error in result")
	}
}

func TestSubscribeNilHandler(t *testing.T) {
	bus := NewEventBus()
	_, err := bus.Subscribe(nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestUnsubscribeNotFound(t *testing.T) {
	bus := NewEventBus()
	err := bus.Unsubscribe("nonexistent")
	if err != ErrUnknownSub {
		t.Fatalf("expected ErrUnknownSub, got %v", err)
	}
}

func TestShutdownAndRestart(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}

	if err := bus.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}

	_, err := bus.PublishSync(context.Background(), &EventBase{})
	if err != ErrBusShutdown {
		t.Fatalf("expected ErrBusShutdown after shutdown, got %v", err)
	}
}

func TestSubscribeUnsubscribe(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	subID, err := bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		return HandlerResult{}
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := bus.Unsubscribe(subID); err != nil {
		t.Fatal(err)
	}

	if err := bus.Unsubscribe(subID); err != ErrUnknownSub {
		t.Fatalf("expected ErrUnknownSub on double unsubscribe, got %v", err)
	}
}

func TestMultipleSubscribersPriority(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	var order []string
	var mu sync.Mutex

	bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		mu.Lock()
		order = append(order, "low")
		mu.Unlock()
		return HandlerResult{}
	}, WithPriority(1))

	bus.Subscribe(func(ctx context.Context, event Event) HandlerResult {
		mu.Lock()
		order = append(order, "high")
		mu.Unlock()
		return HandlerResult{}
	}, WithPriority(10))

	evt := &SNSDeliveryEvent{
		EventBase: EventBase{ID: "prio-test", Timestamp: time.Now().UTC(), Source: "test", Region: "us-east-1"},
		TopicARN:  "arn:aws:sns:us-east-1:123456789012:test",
	}

	_, err := bus.PublishSync(context.Background(), evt)
	if err != nil {
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(order))
	}
	if order[0] != "high" || order[1] != "low" {
		t.Fatalf("expected [high, low], got %v", order)
	}
}

func TestEventRegistry(t *testing.T) {
	registry := NewEventRegistry()
	registry.Register("sns:deliver", func() Event { return &SNSDeliveryEvent{} })

	evt, err := registry.New("sns:deliver")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := evt.(*SNSDeliveryEvent); !ok {
		t.Fatal("expected *SNSDeliveryEvent")
	}

	_, err = registry.New("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestEventSerialization(t *testing.T) {
	evt := &SNSDeliveryEvent{
		EventBase: EventBase{
			ID:        "ser-test",
			Timestamp: time.Now().UTC(),
			Source:    "test",
			Region:    "us-east-1",
			Caller:    CallerContext{ServicePrincipal: "test.amazonaws.com"},
		},
		TopicARN:  "arn:aws:sns:us-east-1:123456789012:my-topic",
		MessageID: "msg-1",
		Message:   `{"key":"value"}`,
	}

	data, err := SerializeEvent(evt)
	if err != nil {
		t.Fatal(err)
	}

	registry := NewEventRegistry()
	registry.Register("sns:deliver", func() Event { return &SNSDeliveryEvent{} })

	restored, err := registry.Deserialize("sns:deliver", data)
	if err != nil {
		t.Fatal(err)
	}

	typed := restored.(*SNSDeliveryEvent)
	if typed.EventID() != evt.EventID() {
		t.Fatalf("expected ID %q, got %q", evt.EventID(), typed.EventID())
	}
	if typed.Caller.ServicePrincipal != "test.amazonaws.com" {
		t.Fatalf("expected caller principal test.amazonaws.com, got %s", typed.Caller.ServicePrincipal)
	}
}

func TestBusPolicyDocument(t *testing.T) {
	doc := EmptyBusPolicyDocument()
	if doc.Version != "2012-10-17" {
		t.Fatalf("expected version 2012-10-17, got %s", doc.Version)
	}
}

type testEvent struct {
	EventBase
	Name string
}

func (e *testEvent) EventType() string { return "test:event" }

func TestSubscribeTyped(t *testing.T) {
	bus := NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer bus.Shutdown(context.Background())

	called := atomic.Int32{}
	_, err := SubscribeTyped[*testEvent](bus, func(ctx context.Context, event *testEvent) HandlerResult {
		called.Add(1)
		return HandlerResult{StatusCode: 200}
	})
	if err != nil {
		t.Fatal(err)
	}

	evt := &testEvent{
		EventBase: EventBase{ID: "typed-test", Timestamp: time.Now().UTC(), Source: "test", Region: "us-east-1"},
		Name:      "hello",
	}

	_, err = bus.PublishSync(context.Background(), evt)
	if err != nil {
		t.Fatal(err)
	}

	if called.Load() != 1 {
		t.Fatalf("expected handler called 1 time, got %d", called.Load())
	}
}

func TestWithAsyncOption(t *testing.T) {
	cfg := &subscribeConfig{}
	WithAsync()(cfg)
	if !cfg.async {
		t.Fatal("expected async to be true")
	}
}

func TestWithMaxConcurrencyOption(t *testing.T) {
	cfg := &subscribeConfig{}
	WithMaxConcurrency(5)(cfg)
	if cfg.maxConcurrency != 5 {
		t.Fatalf("expected maxConcurrency 5, got %d", cfg.maxConcurrency)
	}
}

func TestOutboxStatusString(t *testing.T) {
	tests := []struct {
		status OutboxStatus
		want   string
	}{
		{OutboxPending, "PENDING"},
		{OutboxProcessing, "PROCESSING"},
		{OutboxDelivered, "DELIVERED"},
		{OutboxFailed, "FAILED"},
	}

	for _, tt := range tests {
		got := tt.status.String()
		if got != tt.want {
			t.Errorf("OutboxStatus(%d).String() = %q, want %q", tt.status, got, tt.want)
		}
	}
}

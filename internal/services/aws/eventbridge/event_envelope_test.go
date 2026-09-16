package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// fullEnvelopeEvent returns an Event with every member populated, including
// a sub-second timestamp and a trace header.
func fullEnvelopeEvent() *eventsstore.Event {
	return &eventsstore.Event{
		ID:           "evt-round-trip",
		Version:      "0",
		DetailType:   "order.created",
		Source:       "com.example.orders",
		Account:      "000000000000",
		Time:         time.Date(2026, 9, 15, 10, 20, 30, 123456789, time.UTC),
		Region:       "us-east-1",
		Resources:    []string{"arn:aws:s3:::demobucket/order-7.json"},
		Detail:       map[string]interface{}{"count": float64(3), "nested": map[string]interface{}{"state": "new"}},
		EventBusName: "original-bus",
		TraceHeader:  "Root=1-5759e988-bp186060bgaddr;Parent=53995c3f42cd8ad8;Sampled=1",
	}
}

// assertEnvelopeRoundTrip checks that the rebuilt Event carries every
// envelope member unchanged (sub-second time included) onto the named bus.
// The trace header is the deliberate exception: AWS propagates it on the
// transport, never on the event, so it must not survive the envelope.
func assertEnvelopeRoundTrip(t *testing.T, rebuilt *eventsstore.Event, original *eventsstore.Event, busName string) {
	t.Helper()
	if rebuilt.ID != original.ID || rebuilt.Version != original.Version ||
		rebuilt.DetailType != original.DetailType || rebuilt.Source != original.Source ||
		rebuilt.Account != original.Account || rebuilt.Region != original.Region {
		t.Fatalf("scalar envelope members changed in the round-trip: got %+v", rebuilt)
	}
	if !rebuilt.Time.Equal(original.Time) {
		t.Fatalf("time changed in the round-trip: got %v, want %v (sub-second precision must survive)",
			rebuilt.Time, original.Time)
	}
	if !reflect.DeepEqual(rebuilt.Resources, original.Resources) {
		t.Fatalf("resources changed in the round-trip: got %v, want %v", rebuilt.Resources, original.Resources)
	}
	if !reflect.DeepEqual(rebuilt.Detail, original.Detail) {
		t.Fatalf("detail changed in the round-trip: got %v, want %v", rebuilt.Detail, original.Detail)
	}
	if rebuilt.EventBusName != busName {
		t.Fatalf("the rebuilt event must land on the caller's bus %q, got %q", busName, rebuilt.EventBusName)
	}
	if rebuilt.TraceHeader != "" {
		t.Fatal("the trace header must not ride the event envelope (AWS propagates it on the transport, never on the event)")
	}
}

// TestEventEnvelopeRoundTrip closes the envelope codec class: the nine
// members plus the sub-second timestamp survive eventEnvelope →
// eventFromEnvelope both in-process and across a JSON hop, the envelope
// key set is exactly the nine AWS members, and the trace header never
// enters the envelope.
func TestEventEnvelopeRoundTrip(t *testing.T) {
	original := fullEnvelopeEvent()
	envelope := eventEnvelope(original)

	wantKeys := []string{
		"version", "id", "detail-type", "source", "account", "time", "region", "resources", "detail",
	}
	if len(envelope) != len(wantKeys) {
		t.Fatalf("the envelope must carry exactly the nine AWS members, got %d: %v", len(envelope), envelope)
	}
	for _, key := range wantKeys {
		if _, present := envelope[key]; !present {
			t.Fatalf("envelope member %q missing", key)
		}
	}
	if got := envelope["time"].(string); got != "2026-09-15T10:20:30.123456789Z" {
		t.Fatalf("sub-second precision lost in the envelope time: %q", got)
	}

	// In-process round-trip (the envelope map keeps Go-typed resources).
	assertEnvelopeRoundTrip(t, eventFromEnvelope(envelope, "rebuilt-bus"), original, "rebuilt-bus")

	// JSON hop (the wire form: bus-to-bus delivery input, archived event).
	wire, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	var hopped map[string]interface{}
	if err := json.Unmarshal(wire, &hopped); err != nil {
		t.Fatal(err)
	}
	assertEnvelopeRoundTrip(t, eventFromEnvelope(hopped, "hopped-bus"), original, "hopped-bus")
}

// TestReplayNameStampRidesTheEnvelope pins the replay-name metadata field's
// envelope behaviour: live events carry exactly the nine AWS members, a
// replayed event gains the replay-name stamp, and the stamp survives the
// JSON hop so a replayed event forwarded through a bus target stays
// recognisable on the downstream bus.
func TestReplayNameStampRidesTheEnvelope(t *testing.T) {
	if _, present := eventEnvelope(fullEnvelopeEvent())["replay-name"]; present {
		t.Fatal("a live event must not carry the replay-name member")
	}

	replayed := fullEnvelopeEvent()
	replayed.ReplayName = "stamp-round-trip"
	envelope := eventEnvelope(replayed)
	if got := envelope["replay-name"]; got != "stamp-round-trip" {
		t.Fatalf("a replayed event must carry the replay-name stamp, got %v", got)
	}
	if len(envelope) != 10 {
		t.Fatalf("a replayed event's envelope must carry the nine members plus the stamp, got %d: %v", len(envelope), envelope)
	}

	// JSON hop (the wire form a bus target receives and re-ingests).
	wire, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	var hopped map[string]interface{}
	if err := json.Unmarshal(wire, &hopped); err != nil {
		t.Fatal(err)
	}
	if rebuilt := eventFromEnvelope(hopped, "downstream-bus"); rebuilt.ReplayName != "stamp-round-trip" {
		t.Fatalf("the replay-name stamp must survive the envelope hop, got %q", rebuilt.ReplayName)
	}
}

// TestBusToBusDeliveryPreservesDetailTypeAndIdentity runs the production
// bus-to-bus path end to end: an envelope delivered onto a bus whose rule
// matches on detail-type must fire that rule, and the delivered payload
// must carry both the original detail-type and the original event ID — the
// Kinesis partition-key default derives from the event's own ID, not from a
// bus-synthesised one.
func TestBusToBusDeliveryPreservesDetailTypeAndIdentity(t *testing.T) {
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "orders"}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateRule(ctx, &eventsstore.Rule{
		Name:         "match-created",
		EventBusName: "orders",
		EventPattern: `{"detail-type":["order.created"]}`,
		State:        eventsstore.RuleStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutTarget(ctx, &eventsstore.Target{
		ID:           "stream",
		RuleName:     "match-created",
		EventBusName: "orders",
		ARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/orders",
	}); err != nil {
		t.Fatal(err)
	}

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	invoker := &recordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	if err := bus.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })

	svc := NewEventsService(mgr, "000000000000")
	svc.SetEventsStore("us-east-1", store)
	svc.SetEventBus(bus)
	t.Cleanup(svc.Close)

	original := fullEnvelopeEvent()
	input, err := json.Marshal(eventEnvelope(original))
	if err != nil {
		t.Fatal(err)
	}
	delivery := &eventbus.EventBridgeDeliveryEvent{
		TargetARN: "arn:aws:events:us-east-1:000000000000:event-bus/orders",
		Input:     input,
	}
	delivery.Region = "us-east-1"
	if res := svc.handleEventBusDelivery(ctx, delivery); res.Error != nil {
		t.Fatalf("bus delivery failed: %v", res.Error)
	}

	// The delivery lands on the Kinesis target asynchronously, through the
	// bus workers.
	deadline := time.Now().Add(5 * time.Second)
	for len(invoker.keys()) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	keys := invoker.keys()
	payloads := invoker.payloads()
	if len(keys) != 1 || len(payloads) != 1 {
		t.Fatalf("expected exactly one Kinesis put for the delivered event, got %d", len(keys))
	}
	if keys[0] != original.ID {
		t.Fatalf("partition key %q must be the event's own ID %q — the documented Kinesis default derives from the event ID",
			keys[0], original.ID)
	}

	var delivered map[string]interface{}
	// The Kinesis delivery path pre-encodes the payload to base64 to match
	// the format the Kinesis SDK's PutRecord would send.
	decoded, err := base64.StdEncoding.DecodeString(string(payloads[0]))
	if err != nil {
		t.Fatalf("the delivered payload must be base64-encoded: %v", err)
	}
	if err := json.Unmarshal(decoded, &delivered); err != nil {
		t.Fatalf("the delivered payload must be a JSON envelope: %v", err)
	}
	if got := delivered["detail-type"]; got != original.DetailType {
		t.Fatalf("detail-type must survive the bus hop: got %v, want %q", got, original.DetailType)
	}
	if got := delivered["id"]; got != original.ID {
		t.Fatalf("the event ID must survive the bus hop: got %v, want %q", got, original.ID)
	}
	if got := delivered["time"]; got != "2026-09-15T10:20:30.123456789Z" {
		t.Fatalf("sub-second precision must survive the bus hop, got %v", got)
	}
}

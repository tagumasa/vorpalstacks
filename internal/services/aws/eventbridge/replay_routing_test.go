package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// newReplayRoutingFixture builds the replay routing stage: one source bus,
// one enabled archive on it, and two ENABLED rules whose patterns match the
// fixture's events and whose Kinesis targets land on distinct recorded
// streams, wired through a real service bus so deliveries run the
// production asynchronous path.
func newReplayRoutingFixture(t *testing.T) (*EventsService, *eventsstore.EventsStore, *streamRecordingInvoker) {
	t.Helper()
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "routing-bus"}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{
		Name:         "routing-archive",
		EventBusName: "routing-bus",
		State:        eventsstore.ArchiveStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}
	for _, ruleName := range []string{"filtered-rule", "other-rule"} {
		if err := store.CreateRule(ctx, &eventsstore.Rule{
			Name:         ruleName,
			EventBusName: "routing-bus",
			// A pattern rule: pattern-less rules are scheduled rules and
			// never match bus events, so the routing pins need a pattern
			// that matches every source the fixture's consumers ingest
			// (com.example.routing, com.example.paging).
			EventPattern: `{"source":[{"prefix":"com.example."}]}`,
			State:        eventsstore.RuleStateEnabled,
		}); err != nil {
			t.Fatal(err)
		}
		if err := store.PutTarget(ctx, &eventsstore.Target{
			ID:           "stream",
			RuleName:     ruleName,
			EventBusName: "routing-bus",
			ARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/replay-" + ruleName,
		}); err != nil {
			t.Fatal(err)
		}
	}

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	invoker := &streamRecordingInvoker{}
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
	return svc, store, invoker
}

// waitForStreamDeliveries blocks until at least n deliveries were recorded
// or the deadline passes, then lets asynchronous stragglers settle.
func waitForStreamDeliveries(t *testing.T, invoker *streamRecordingInvoker, n int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for len(invoker.streams()) < n && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)
}

func streamCounts(streams []string) map[string]int {
	counts := map[string]int{}
	for _, s := range streams {
		counts[s]++
	}
	return counts
}

// TestReplayRoutingContract pins the replay routing controls end to end: a
// replay whose destination names FilterArns fires only the chosen rule's
// targets on the source bus, and the replayed payload carries the
// replay-name stamp that lets consumers distinguish replays from live
// events.
func TestReplayRoutingContract(t *testing.T) {
	ctx := context.Background()
	svc, store, invoker := newReplayRoutingFixture(t)

	// Ingest one event onto the bus: the archive captures it and both
	// rules deliver — the baseline the replay's selectivity is measured
	// against.
	result, err := svc.putEventsCore(ctx, store, PutEventsInput{
		Entries: []interface{}{map[string]interface{}{
			"EventBusName": "routing-bus",
			"Source":       "com.example.routing",
			"DetailType":   "RoutingTest",
			"Detail":       `{"k":"v"}`,
		}},
		Region: "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	originalEventID, _ := result.Entries[0]["EventId"].(string)

	waitForStreamDeliveries(t, invoker, 2)
	if baseline := streamCounts(invoker.streams()); baseline["replay-filtered-rule"] != 1 || baseline["replay-other-rule"] != 1 {
		t.Fatalf("the ingested event must reach both rules' streams exactly once, got %v", baseline)
	}

	replayName := "routing-replay"
	if _, err := svc.startReplayCore(ctx, store, "us-east-1", StartReplayInput{
		ReplayName:     replayName,
		EventSourceArn: "arn:aws:events:us-east-1:000000000000:archive/routing-archive",
		Destination: map[string]interface{}{
			"Arn":        "arn:aws:events:us-east-1:000000000000:event-bus/routing-bus",
			"FilterArns": []interface{}{ruleARNFor("routing-bus", "filtered-rule")},
		},
		EventStartTime: float64(time.Now().Add(-time.Hour).Unix()),
		EventEndTime:   float64(time.Now().Add(time.Hour).Unix()),
	}); err != nil {
		t.Fatalf("start replay: %v", err)
	}

	waitForStreamDeliveries(t, invoker, 3)
	counts := streamCounts(invoker.streams())
	if counts["replay-filtered-rule"] != 2 {
		t.Fatalf("the filtered rule must receive the replayed event, got %v", counts)
	}
	if counts["replay-other-rule"] != 1 {
		t.Fatalf("a rule outside FilterArns must not receive the replayed event, got %v", counts)
	}

	// The second delivery to the filtered stream is the replay: its payload
	// must carry the replay-name stamp and the original event's identity.
	streams := invoker.streams()
	payloads := invoker.payloads()
	replayPayloadIndex := -1
	seen := 0
	for i, s := range streams {
		if s == "replay-filtered-rule" {
			seen++
			if seen == 2 {
				replayPayloadIndex = i
			}
		}
	}
	if replayPayloadIndex < 0 {
		t.Fatalf("the replayed delivery never landed, streams: %v", streams)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(payloads[replayPayloadIndex]))
	if err != nil {
		t.Fatalf("the replayed payload must be base64-encoded: %v", err)
	}
	var replayed map[string]interface{}
	if err := json.Unmarshal(decoded, &replayed); err != nil {
		t.Fatalf("the replayed payload must be a JSON envelope: %v", err)
	}
	if got := replayed["replay-name"]; got != replayName {
		t.Fatalf("the replayed payload must carry the replay-name stamp %q, got %v", replayName, got)
	}
	if got := replayed["id"]; got != originalEventID {
		t.Fatalf("the replayed payload must keep the original event ID %q, got %v", originalEventID, got)
	}
}

// TestStartReplayRejectsCrossBusDestination pins the destination constraint:
// archives replay only onto their source event bus, so a destination naming
// any other existing bus is rejected and no replay record is written.
func TestStartReplayRejectsCrossBusDestination(t *testing.T) {
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	for _, busName := range []string{"source-bus", "other-bus"} {
		if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: busName}); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{
		Name:         "xbus-archive",
		EventBusName: "source-bus",
		State:        eventsstore.ArchiveStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}

	svc := NewEventsService(nil, "000000000000")
	_, err = svc.startReplayCore(ctx, store, "us-east-1", StartReplayInput{
		ReplayName:     "cross-bus-replay",
		EventSourceArn: "arn:aws:events:us-east-1:000000000000:archive/xbus-archive",
		Destination: map[string]interface{}{
			"Arn": "arn:aws:events:us-east-1:000000000000:event-bus/other-bus",
		},
		EventStartTime: float64(time.Now().Add(-time.Hour).Unix()),
		EventEndTime:   float64(time.Now().Add(time.Hour).Unix()),
	})

	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
		t.Fatalf("a cross-bus destination must be rejected as ValidationException, got %v", err)
	}
	if _, err := store.GetReplay(ctx, "cross-bus-replay"); !errors.Is(err, eventsstore.ErrReplayNotFound) {
		t.Fatalf("the rejected replay must not be recorded, got %v", err)
	}
}

package eventbridge

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestBusToBusIngressArchives pins the production bus-target ingress: an
// event delivered onto a bus through handleEventBusDelivery archives onto
// that bus's enabled archives exactly like the PutEvents planes — the
// archives user guide states a pattern filter, never an ingress-path one.
func TestBusToBusIngressArchives(t *testing.T) {
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
	if err := store.CreateArchive(ctx, &eventsstore.Archive{
		Name:         "orders-archive",
		EventBusName: "orders",
		State:        eventsstore.ArchiveStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	svc.SetEventsStore("us-east-1", store)

	envelope := eventEnvelope(fullEnvelopeEvent())
	input, err := json.Marshal(envelope)
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

	if got := archivedCount(t, store, "orders-archive"); got != 1 {
		t.Fatalf("the bus-target ingress must archive onto the target bus, got %d events", got)
	}

	// The bus-missing guard still holds: a delivery onto a deleted bus
	// neither archives nor crashes the handler path's error contract.
	missing := &eventbus.EventBridgeDeliveryEvent{
		TargetARN: "arn:aws:events:us-east-1:000000000000:event-bus/vanished",
		Input:     input,
	}
	missing.Region = "us-east-1"
	if res := svc.handleEventBusDelivery(ctx, missing); res.Error == nil {
		t.Fatal("a bus delivery onto a non-existent bus must fail the handler")
	}
	if got := archivedCount(t, store, "orders-archive"); got != 1 {
		t.Fatalf("a failed bus delivery must not archive, got %d events", got)
	}
}

// TestScheduledFireArchives pins the scheduler ingress: a scheduled fire
// archives its event onto the firing rule's bus — the pattern-only archive
// contract — while staying scoped to the firing rule's own targets.
func TestScheduledFireArchives(t *testing.T) {
	ctx := context.Background()

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "sched-arch-bus"}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{
		Name:         "sched-arch-archive",
		EventBusName: "sched-arch-bus",
		State:        eventsstore.ArchiveStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateRule(ctx, &eventsstore.Rule{
		Name:               "arch-minutely",
		EventBusName:       "sched-arch-bus",
		ScheduleExpression: "rate(1 minute)",
		State:              eventsstore.RuleStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutTarget(ctx, &eventsstore.Target{
		ID:           "stream",
		RuleName:     "arch-minutely",
		EventBusName: "sched-arch-bus",
		ARN:          "arn:aws:kinesis:us-east-1:000000000000:stream/sched-arch",
	}); err != nil {
		t.Fatal(err)
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

	now := time.Now().UTC().Truncate(time.Minute)
	svc.tickScheduledRules(ctx, now.Add(2*time.Minute))

	waitForStreamDeliveries(t, invoker, 1)
	if counts := streamCounts(invoker.streams()); counts["sched-arch"] != 1 {
		t.Fatalf("the scheduled fire must reach its target stream exactly once, got %v", counts)
	}
	if got := archivedCount(t, store, "sched-arch-archive"); got != 1 {
		t.Fatalf("the scheduled fire must archive onto the firing rule's bus, got %d events", got)
	}
}

// TestReplayedEventsAreNotRearchived pins the replay exclusion at the
// unified ingress: replaying an archive onto its own source bus redelivers
// the events to the bus's rules but never grows the archive the events were
// replayed from — the replay-name stamp is the discriminator.
func TestReplayedEventsAreNotRearchived(t *testing.T) {
	ctx := context.Background()
	svc, store, invoker := newReplayRoutingFixture(t)

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
	if result.FailedEntryCount != 0 {
		t.Fatalf("the ingested event must succeed, got %+v", result)
	}

	waitForStreamDeliveries(t, invoker, 2)
	if got := archivedCount(t, store, "routing-archive"); got != 1 {
		t.Fatalf("the ingested event must archive once, got %d events", got)
	}

	// Replay the whole window without FilterArns: both rules receive the
	// replayed event, and the archive must not grow.
	if _, err := svc.startReplayCore(ctx, store, "us-east-1", StartReplayInput{
		ReplayName:     "rearchive-replay",
		EventSourceArn: "arn:aws:events:us-east-1:000000000000:archive/routing-archive",
		Destination: map[string]interface{}{
			"Arn": "arn:aws:events:us-east-1:000000000000:event-bus/routing-bus",
		},
		EventStartTime: float64(time.Now().Add(-time.Hour).Unix()),
		EventEndTime:   float64(time.Now().Add(time.Hour).Unix()),
	}); err != nil {
		t.Fatalf("start replay: %v", err)
	}

	waitForStreamDeliveries(t, invoker, 4)
	counts := streamCounts(invoker.streams())
	if counts["replay-filtered-rule"] != 2 || counts["replay-other-rule"] != 2 {
		t.Fatalf("the unfiltered replay must redeliver to both rules, got %v", counts)
	}
	if got := archivedCount(t, store, "routing-archive"); got != 1 {
		t.Fatalf("a replayed event must not re-enter the archive, got %d events", got)
	}
}

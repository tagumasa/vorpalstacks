package eventbridge

import (
	"context"
	"testing"
)

// TestListEventBusesScopesToBusRecords pins the listing scope: event bus
// records are keyed by their full ARN, and a row the events-eventbuses
// bucket carries under any other key shape (residue from the retired
// raw-write cross-service ingress plane, keyed "events:<bus>:<id>") never
// surfaces as a bus — least of all as an empty-name ghost entry.
func TestListEventBusesScopesToBusRecords(t *testing.T) {
	store, _ := newRecordLockTestStores(t)
	ctx := context.Background()

	if err := store.CreateEventBus(ctx, &EventBus{Name: "real-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	// A row keyed the way the retired raw-write ingress plane wrote them,
	// written straight through the base store like that plane did.
	if err := store.Put("events:real-bus:legacy-row", map[string]interface{}{"id": "legacy-row"}); err != nil {
		t.Fatalf("plant residue row: %v", err)
	}

	result, err := store.ListEventBuses(ctx, "", ListLimitMaximum, "")
	if err != nil {
		t.Fatalf("list buses: %v", err)
	}
	names := map[string]bool{}
	for _, bus := range result.EventBuses {
		if bus.Name == "" {
			t.Fatalf("a non-bus row surfaced as an empty-name bus: %+v", result.EventBuses)
		}
		names[bus.Name] = true
	}
	if !names["real-bus"] {
		t.Fatalf("the created bus must be listed, got %v", names)
	}
	if names["legacy-row"] {
		t.Fatalf("a residue row must never surface as a bus, got %v", names)
	}
}

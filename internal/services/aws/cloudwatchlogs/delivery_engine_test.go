package cloudwatchlogs

import (
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The delivery fetch's late window: an event ingested below a cursor the
// engine has already persisted returns in the late slice, its sort
// position strictly ahead of the normal window, and the next pass with
// the advanced ingestion mark sees nothing — the at-least-once promise
// holds for backdated puts, which no timestamp-window fetch could ever
// recover.
func TestFetchDeliveryEventsReturnsBackdatedLateWindow(t *testing.T) {
	_, store := newTestService(t)
	createTestLogGroup(t, store, "late-engine-group")
	if err := store.CreateLogStream(&logsstore.LogStream{Name: "s", LogGroupName: "late-engine-group"}); err != nil {
		t.Fatalf("stream: %v", err)
	}
	base := time.Now().UnixMilli() - 100_000
	if _, err := store.PutLogEvents("late-engine-group", "s", []logsstore.LogEntry{
		{Timestamp: base, Message: "first", IngestionTime: base},
		{Timestamp: base + 1000, Message: "cursor-event", IngestionTime: base},
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	// The cursor sits on the second event: one delivered event at the
	// triple.
	delivery := &logsstore.Delivery{
		CursorTime:          base + 1000,
		CursorStream:        "s",
		CursorDigest:        eventDigest("cursor-event"),
		CursorCount:         1,
		CursorIngestionMark: base + 10_000,
	}
	// The backdated put lands below the cursor, ingested after the mark.
	if _, err := store.PutLogEvents("late-engine-group", "s", []logsstore.LogEntry{
		{Timestamp: base - 50_000, Message: "backdated", IngestionTime: base + 60_000},
	}); err != nil {
		t.Fatalf("backdated put: %v", err)
	}

	events, late, scanMark, err := fetchDeliveryEvents(store, "late-engine-group", delivery)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("normal window: want empty past the cursor, got %d events", len(events))
	}
	if len(late) != 1 || late[0].Message != "backdated" {
		t.Fatalf("late window: want the single backdated event, got %v", late)
	}
	if scanMark == 0 {
		t.Fatal("scan mark must carry the pass's start clock")
	}

	// The pass that persisted the mark no longer re-finds the event.
	delivery.CursorIngestionMark = scanMark
	_, late, _, err = fetchDeliveryEvents(store, "late-engine-group", delivery)
	if err != nil {
		t.Fatalf("second fetch: %v", err)
	}
	if len(late) != 0 {
		t.Fatalf("advanced mark: want an empty late window, got %v", late)
	}
}

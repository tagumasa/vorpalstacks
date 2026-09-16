package eventbridge

import (
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestDeleteArchiveCoreEventsFirst pins the archive delete ordering: the
// stored events are removed before the archive record, and a failure
// removing them surfaces as an error with the archive record left in
// place — once the record is gone the retention worker no longer
// iterates the archive, so a record-first order with a swallowed events
// error would orphan the events permanently.
func TestDeleteArchiveCoreEventsFirst(t *testing.T) {
	ctx := t.Context()
	svc := newEventBusCoreTestService()

	seed := func(t *testing.T) *eventsstore.EventsStore {
		t.Helper()
		store := newEventBusCoreTestStore(t)
		if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
			t.Fatalf("create bus: %v", err)
		}
		archive := &eventsstore.Archive{Name: "a", EventBusName: "default"}
		if err := store.CreateArchive(ctx, archive); err != nil {
			t.Fatalf("create archive: %v", err)
		}
		if err := store.StoreArchiveEvent(ctx, "a", &eventsstore.ArchivedEvent{
			ID:        "e1",
			Event:     map[string]interface{}{"id": "e1"},
			Timestamp: time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC),
		}); err != nil {
			t.Fatalf("store archive event: %v", err)
		}
		return store
	}

	// Healthy path: both the events and the record are gone.
	healthy := seed(t)
	if err := svc.deleteArchiveCore(ctx, healthy, "a"); err != nil {
		t.Fatalf("deleteArchiveCore: %v", err)
	}
	if _, err := healthy.GetArchive(ctx, "a"); err != eventsstore.ErrArchiveNotFound {
		t.Errorf("archive record survived the delete (err = %v)", err)
	}
	events, err := healthy.ListArchiveEvents(ctx, "a", time.Time{}, time.Time{}, eventsstore.ListLimitMaximum, "")
	if err != nil {
		t.Fatalf("list archive events: %v", err)
	}
	if len(events.Events) != 0 {
		t.Errorf("%d archived events survived the delete", len(events.Events))
	}

	// Faulted events bucket: the failure surfaces and the archive record
	// is left in place for a retry.
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	fault := errors.New("transient pebble fault")
	faulted := &faultedReadStorage{BasicStorage: st, failOn: "events-archived-events-us-east-1", fail: fault}
	faultedStore := eventsstore.NewEventsStore(faulted, "000000000000", "us-east-1")
	if err := faultedStore.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	if err := faultedStore.CreateArchive(ctx, &eventsstore.Archive{Name: "a", EventBusName: "default"}); err != nil {
		t.Fatalf("create archive: %v", err)
	}

	if err := svc.deleteArchiveCore(ctx, faultedStore, "a"); err == nil {
		t.Fatal("events-delete failure was swallowed (deleteArchiveCore returned nil)")
	}
	if _, err := faultedStore.GetArchive(ctx, "a"); err != nil {
		t.Fatalf("archive record deleted despite the events-delete failure (err = %v)", err)
	}
}

// TestStoreArchiveEventAfterArchiveDelete pins the ingress-vs-delete window:
// once the archive is deleted, an ingress write is rejected under the family
// lock instead of landing an event row under the dead archive name — a row
// a same-named archive created later would surface as its own.
func TestStoreArchiveEventAfterArchiveDelete(t *testing.T) {
	ctx := t.Context()
	store := newEventBusCoreTestStore(t)
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{Name: "a", EventBusName: "default"}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	if err := store.DeleteArchive(ctx, "a"); err != nil {
		t.Fatalf("delete archive: %v", err)
	}

	err := store.StoreArchiveEvent(ctx, "a", &eventsstore.ArchivedEvent{
		ID:        "e1",
		Event:     map[string]interface{}{"id": "e1"},
		Timestamp: time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, eventsstore.ErrArchiveNotFound) {
		t.Fatalf("an ingress write after the delete must fail with ErrArchiveNotFound, got %v", err)
	}

	// A same-named archive created later must not surface the rejected row.
	if err := store.CreateArchive(ctx, &eventsstore.Archive{Name: "a", EventBusName: "default"}); err != nil {
		t.Fatalf("recreate archive: %v", err)
	}
	events, err := store.ListArchiveEvents(ctx, "a", time.Time{}, time.Time{}, eventsstore.ListLimitMaximum, "")
	if err != nil {
		t.Fatalf("list archive events: %v", err)
	}
	if len(events.Events) != 0 {
		t.Fatalf("a rejected ingress row surfaced in the recreated archive: %d events", len(events.Events))
	}
}

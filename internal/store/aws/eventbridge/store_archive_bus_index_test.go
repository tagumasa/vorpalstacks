package eventbridge

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// TestArchiveBusIndexScopesListingsToTheBus pins the per-bus archive index:
// ListArchivesForEventBus returns exactly the bus's own archives (a prefix
// scan that cannot bleed into another bus's rows), and deleting an archive
// removes its index row so the archive stops matching event delivery.
func TestArchiveBusIndexScopesListingsToTheBus(t *testing.T) {
	ctx := context.Background()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := NewEventsStore(st, "000000000000", "us-east-1")

	for _, bus := range []string{"alpha-bus", "beta-bus", "alpha-bus-extra"} {
		if err := store.CreateEventBus(ctx, &EventBus{Name: bus}); err != nil {
			t.Fatal(err)
		}
	}
	create := func(bus, name string) {
		t.Helper()
		if err := store.CreateArchive(ctx, &Archive{Name: name, EventBusName: bus, State: ArchiveStateEnabled}); err != nil {
			t.Fatal(err)
		}
	}
	create("alpha-bus", "alpha-one")
	create("alpha-bus", "alpha-two")
	create("beta-bus", "beta-one")

	names := func(bus string) []string {
		t.Helper()
		archives, err := store.ListArchivesForEventBus(ctx, bus)
		if err != nil {
			t.Fatalf("list archives for %s: %v", bus, err)
		}
		got := make([]string, 0, len(archives))
		for _, a := range archives {
			got = append(got, a.Name)
		}
		return got
	}

	if got := names("alpha-bus"); len(got) != 2 {
		t.Fatalf("the alpha bus must list exactly its two archives, got %v", got)
	}
	if got := names("beta-bus"); len(got) != 1 || got[0] != "beta-one" {
		t.Fatalf("the beta bus must list exactly its own archive, got %v", got)
	}

	if err := store.DeleteArchive(ctx, "alpha-one"); err != nil {
		t.Fatal(err)
	}
	if got := names("alpha-bus"); len(got) != 1 || got[0] != "alpha-two" {
		t.Fatalf("a deleted archive must leave the index, got %v", got)
	}
}

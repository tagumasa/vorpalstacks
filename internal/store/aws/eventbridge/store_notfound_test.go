package eventbridge

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// failingGetBucket fails every read with the injected error, standing in
// for a transient storage fault.
type failingGetBucket struct {
	storage.Bucket
	fail error
}

func (b *failingGetBucket) Get(key []byte) ([]byte, error) { return nil, b.fail }

// faultedGetStorage routes every bucket's point reads through the fault.
type faultedGetStorage struct {
	storage.BasicStorage
	fail error
}

func (s *faultedGetStorage) Bucket(name string) storage.Bucket {
	return &failingGetBucket{Bucket: s.BasicStorage.Bucket(name), fail: s.fail}
}

// TestGetMutateDiscriminateFaultFromMissing pins the error contract of
// every family Get*/Mutate*: a storage fault must surface as itself —
// never collapsed into the family's not-found sentinel, which would
// answer "does not exist" (a ResourceNotFoundException on the API plane)
// for a transient I/O failure.
func TestGetMutateDiscriminateFaultFromMissing(t *testing.T) {
	ctx := t.Context()
	fault := errors.New("transient pebble fault")

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := NewEventsStore(&faultedGetStorage{BasicStorage: st, fail: fault}, "000000000000", "us-east-1")

	cases := []struct {
		name string
		call func() error
	}{
		{"GetEventBus", func() error { _, err := store.GetEventBus(ctx, "bus"); return err }},
		{"MutateEventBus", func() error {
			return store.MutateEventBus(ctx, "bus", func(*EventBus) error { return nil })
		}},
		{"GetRule", func() error { _, err := store.GetRule(ctx, "bus", "rule"); return err }},
		{"MutateRule", func() error {
			return store.MutateRule(ctx, "bus", "rule", func(*Rule) error { return nil })
		}},
		{"GetTarget", func() error { _, err := store.GetTarget(ctx, "bus", "rule", "t"); return err }},
		{"GetArchive", func() error { _, err := store.GetArchive(ctx, "arch"); return err }},
		{"MutateArchive", func() error {
			return store.MutateArchive(ctx, "arch", func(*Archive) error { return nil })
		}},
		{"GetConnection", func() error { _, err := store.GetConnection(ctx, "conn"); return err }},
		{"MutateConnection", func() error {
			return store.MutateConnection(ctx, "conn", func(*Connection) error { return nil })
		}},
		{"GetApiDestination", func() error { _, err := store.GetApiDestination(ctx, "dest"); return err }},
		{"MutateApiDestination", func() error {
			return store.MutateApiDestination(ctx, "dest", func(*ApiDestination) error { return nil })
		}},
		{"GetReplay", func() error { _, err := store.GetReplay(ctx, "rep"); return err }},
		{"MutateReplay", func() error {
			return store.MutateReplay(ctx, "rep", func(*Replay) error { return nil })
		}},
	}
	for _, tc := range cases {
		if err := tc.call(); !errors.Is(err, fault) {
			t.Errorf("%s: error = %v, want the storage fault itself (not the not-found sentinel)", tc.name, err)
		}
	}
}

// TestCreateEmptyNameGuard pins the store-level empty-name backstop: the
// Create* methods reject an empty resource name with the dedicated
// sentinel instead of the misleading "invalid ARN".
func TestCreateEmptyNameGuard(t *testing.T) {
	ctx := t.Context()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := NewEventsStore(st, "000000000000", "us-east-1")

	creates := []struct {
		name   string
		create func() error
	}{
		{"eventbus", func() error { return store.CreateEventBus(ctx, &EventBus{}) }},
		{"rule", func() error { return store.CreateRule(ctx, &Rule{}) }},
		{"archive", func() error { return store.CreateArchive(ctx, &Archive{}) }},
		{"connection", func() error { return store.CreateConnection(ctx, &Connection{}) }},
		{"api-destination", func() error { return store.CreateApiDestination(ctx, &ApiDestination{}) }},
		{"replay", func() error { return store.CreateReplay(ctx, &Replay{}) }},
	}
	for _, tc := range creates {
		if err := tc.create(); err != ErrEmptyResourceName {
			t.Errorf("%s: empty-name create = %v, want %v", tc.name, err, ErrEmptyResourceName)
		}
	}
}

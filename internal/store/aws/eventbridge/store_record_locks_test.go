package eventbridge

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// newRecordLockTestStores builds two EventsStore instances over the same
// Pebble keyspace — the production premise of the record locks: the
// scheduler, replay and delivery workers and the API handlers each hold
// their own store instance, so only package-scope locking can serialise a
// record cycle between them.
func newRecordLockTestStores(t *testing.T) (*EventsStore, *EventsStore) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewEventsStore(st, "000000000000", "us-east-1"), NewEventsStore(st, "000000000000", "us-east-1")
}

// runDeleteVsMutate pins the family lock ordering: a mutation paused inside
// its fn holds the family record lock, so a concurrent delete from another
// store instance cannot slip between the mutation's read and its write —
// the record ends up deleted, never resurrected by the mutation's write.
func runDeleteVsMutate(t *testing.T, name string, mutate func(paused func() error) error, delete func() error, missing func() error) {
	t.Helper()
	inFn := make(chan struct{})
	release := make(chan struct{})
	mutateDone := make(chan error, 1)
	go func() {
		mutateDone <- mutate(func() error {
			close(inFn)
			<-release
			return nil
		})
	}()
	<-inFn

	delDone := make(chan error, 1)
	go func() { delDone <- delete() }()
	// The delete must be blocked on the family lock while fn is paused.
	select {
	case err := <-delDone:
		close(release)
		t.Fatalf("%s: delete completed while a mutation held the record lock (err = %v)", name, err)
	case <-time.After(50 * time.Millisecond):
	}

	close(release)
	if err := <-mutateDone; err != nil {
		t.Fatalf("%s: mutate: %v", name, err)
	}
	if err := <-delDone; err != nil {
		t.Fatalf("%s: delete after mutation: %v", name, err)
	}
	if err := missing(); err == nil {
		t.Fatalf("%s: record resurrected after delete", name)
	}
}

// TestDeleteVsMutatePerFamily walks every record family's delete-vs-mutate
// window: the mutation runs on one store instance, the delete on a second
// one over the same keyspace.
func TestDeleteVsMutatePerFamily(t *testing.T) {
	ctx := t.Context()

	t.Run("eventbus", func(t *testing.T) {
		a, b := newRecordLockTestStores(t)
		if err := a.CreateEventBus(ctx, &EventBus{Name: "victim"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		runDeleteVsMutate(t, "eventbus",
			func(paused func() error) error {
				return a.MutateEventBus(ctx, "victim", func(bus *EventBus) error {
					bus.Description = "merged"
					return paused()
				})
			},
			func() error { return b.DeleteEventBus(ctx, "victim") },
			func() error { _, err := b.GetEventBus(ctx, "victim"); return err })
	})

	t.Run("archive", func(t *testing.T) {
		a, b := newRecordLockTestStores(t)
		if err := a.CreateArchive(ctx, &Archive{Name: "victim", EventBusName: "default"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		runDeleteVsMutate(t, "archive",
			func(paused func() error) error {
				return a.MutateArchive(ctx, "victim", func(archive *Archive) error {
					archive.Description = "merged"
					return paused()
				})
			},
			func() error { return b.DeleteArchive(ctx, "victim") },
			func() error { _, err := b.GetArchive(ctx, "victim"); return err })
	})

	t.Run("connection", func(t *testing.T) {
		a, b := newRecordLockTestStores(t)
		if err := a.CreateConnection(ctx, &Connection{Name: "victim", AuthorizationType: "BASIC"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		runDeleteVsMutate(t, "connection",
			func(paused func() error) error {
				return a.MutateConnection(ctx, "victim", func(connection *Connection) error {
					connection.Description = "merged"
					return paused()
				})
			},
			func() error { return b.DeleteConnection(ctx, "victim") },
			func() error { _, err := b.GetConnection(ctx, "victim"); return err })
	})

	t.Run("api-destination", func(t *testing.T) {
		a, b := newRecordLockTestStores(t)
		if err := a.CreateApiDestination(ctx, &ApiDestination{Name: "victim", ConnectionARN: "arn:aws:events:us-east-1:000000000000:connection/c"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		runDeleteVsMutate(t, "api-destination",
			func(paused func() error) error {
				return a.MutateApiDestination(ctx, "victim", func(apiDest *ApiDestination) error {
					apiDest.Description = "merged"
					return paused()
				})
			},
			func() error { return b.DeleteApiDestination(ctx, "victim") },
			func() error { _, err := b.GetApiDestination(ctx, "victim"); return err })
	})

	t.Run("replay", func(t *testing.T) {
		a, b := newRecordLockTestStores(t)
		if err := a.CreateReplay(ctx, &Replay{Name: "victim"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		runDeleteVsMutate(t, "replay",
			func(paused func() error) error {
				return a.MutateReplay(ctx, "victim", func(replay *Replay) error {
					replay.Description = "merged"
					return paused()
				})
			},
			func() error { return b.DeleteReplay(ctx, "victim") },
			func() error { _, err := b.GetReplay(ctx, "victim"); return err })
	})
}

// TestCreateReplayCappedAtomic pins the capped create's atomicity: with the
// cap at ten and twelve concurrent creators across separate store instances
// (the API handlers' premise), exactly ten records land — the count-then-
// create pair is one locked section, so no two creators can observe the same
// sub-cap count and both create.
func TestCreateReplayCappedAtomic(t *testing.T) {
	a, b := newRecordLockTestStores(t)
	ctx := t.Context()

	const creators = 12
	var wg sync.WaitGroup
	errs := make([]error, creators)
	for i := 0; i < creators; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			store := a
			if i%2 == 1 {
				store = b
			}
			errs[i] = store.CreateReplayCapped(ctx, &Replay{Name: fmt.Sprintf("racer-%d", i)}, MaxConcurrentReplays)
		}(i)
	}
	wg.Wait()

	created, capped := 0, 0
	for i, err := range errs {
		switch {
		case err == nil:
			created++
		case errors.Is(err, ErrReplayCapReached):
			capped++
		default:
			t.Fatalf("creator %d: unexpected error %v", i, err)
		}
	}
	if created != MaxConcurrentReplays || capped != creators-MaxConcurrentReplays {
		t.Fatalf("created = %d, capped = %d; want exactly %d created and %d capped",
			created, capped, MaxConcurrentReplays, creators-MaxConcurrentReplays)
	}
}

// TestMutateArchiveCountersAcrossInstances pins the archive lost-update
// window: counter increments from one store instance (the delivery path's
// pattern) racing configuration merges from a second instance (the API
// plane's pattern) must both survive — the merge's read-modify-write can
// never regress EventCount/SizeBytes.
func TestMutateArchiveCountersAcrossInstances(t *testing.T) {
	a, b := newRecordLockTestStores(t)
	ctx := t.Context()
	if err := a.CreateArchive(ctx, &Archive{Name: "counted", EventBusName: "default"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	const increments = 100
	const merges = 25
	const eventSize = int64(128)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < increments; i++ {
			if err := a.IncrementArchiveCounters(ctx, "counted", eventSize); err != nil {
				t.Errorf("increment: %v", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < merges; i++ {
			if err := b.MutateArchive(ctx, "counted", func(archive *Archive) error {
				archive.Description = "merged"
				return nil
			}); err != nil {
				t.Errorf("merge: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	got, err := a.GetArchive(ctx, "counted")
	if err != nil {
		t.Fatalf("final get: %v", err)
	}
	if got.EventCount != increments {
		t.Errorf("eventCount = %d, want %d (a merge regressed the counters)", got.EventCount, increments)
	}
	if got.SizeBytes != increments*eventSize {
		t.Errorf("sizeBytes = %d, want %d", got.SizeBytes, increments*eventSize)
	}
	if got.Description != "merged" {
		t.Errorf("description = %q, want %q (an increment overwrote the merge)", got.Description, "merged")
	}
}

// TestPutTargetCreatedAtPreservedAcrossInstances pins the target record
// lock: concurrent puts of the same target ID from two store instances
// must preserve the original creation time — the read-preserve-write cycle
// is atomic, so no put can observe a half-written or reset CreatedAt.
func TestPutTargetCreatedAtPreservedAcrossInstances(t *testing.T) {
	a, b := newRecordLockTestStores(t)
	ctx := t.Context()

	first := &Target{ID: "t1", RuleName: "r", EventBusName: "default", ARN: "arn:aws:events:us-east-1:000000000000:event-bus/default"}
	if err := a.PutTarget(ctx, first); err != nil {
		t.Fatalf("initial put: %v", err)
	}
	original := first.CreatedAt
	if original.IsZero() {
		t.Fatal("initial put did not stamp CreatedAt")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	put := func(s *EventsStore, input string) {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := s.PutTarget(ctx, &Target{
				ID:           "t1",
				RuleName:     "r",
				EventBusName: "default",
				ARN:          "arn:aws:events:us-east-1:000000000000:event-bus/default",
				Input:        input,
			}); err != nil {
				t.Errorf("put: %v", err)
				return
			}
		}
	}
	go put(a, `{"from":"a"}`)
	go put(b, `{"from":"b"}`)
	wg.Wait()

	got, err := b.GetTarget(ctx, "default", "r", "t1")
	if err != nil {
		t.Fatalf("final get: %v", err)
	}
	if !got.CreatedAt.Equal(original) {
		t.Errorf("createdAt = %v, want the original %v (a concurrent put reset it)", got.CreatedAt, original)
	}
}

// TestDeleteTargetNotResurrectedByPut pins the target family's
// delete-vs-mutate window: a put racing the delete from another store
// instance can never leave the record present after the delete returned.
func TestDeleteTargetNotResurrectedByPut(t *testing.T) {
	a, b := newRecordLockTestStores(t)
	ctx := t.Context()

	for i := 0; i < 50; i++ {
		if err := a.PutTarget(ctx, &Target{
			ID:           "victim",
			RuleName:     "r",
			EventBusName: "default",
			ARN:          "arn:aws:events:us-east-1:000000000000:event-bus/default",
		}); err != nil {
			t.Fatalf("put (iteration %d): %v", i, err)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = a.PutTarget(ctx, &Target{
				ID:           "victim",
				RuleName:     "r",
				EventBusName: "default",
				ARN:          "arn:aws:events:us-east-1:000000000000:event-bus/default",
			})
		}()
		go func() {
			defer wg.Done()
			_ = b.DeleteTarget(ctx, "default", "r", "victim")
		}()
		wg.Wait()

		// A deterministic final delete settles the record either way; the
		// assertion is that the racing put never resurrected it before
		// this point and cannot after it.
		if err := b.DeleteTarget(ctx, "default", "r", "victim"); err != nil && err != ErrTargetNotFound {
			t.Fatalf("iteration %d: settle delete: %v", i, err)
		}
		if _, err := a.GetTarget(ctx, "default", "r", "victim"); err != ErrTargetNotFound {
			t.Fatalf("iteration %d: target resurrected after delete (err = %v)", i, err)
		}
	}
}

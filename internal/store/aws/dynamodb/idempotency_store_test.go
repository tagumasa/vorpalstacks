package dynamodb

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func TestIdempotencyStoreSweepExpired(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewIdempotencyStore(st, "us-east-1")
	now := time.Now()

	if err := store.Record("live", "hash-live", IdempotencyStateCompleted, now.Add(5*time.Minute), nil); err != nil {
		t.Fatalf("record live token: %v", err)
	}
	if err := store.Record("expired-1", "hash-one", IdempotencyStateCompleted, now.Add(-time.Minute), nil); err != nil {
		t.Fatalf("record expired token: %v", err)
	}
	if err := store.Record("expired-2", "hash-two", IdempotencyStateCompleted, now.Add(-2*time.Minute), nil); err != nil {
		t.Fatalf("record expired token: %v", err)
	}

	removed, err := store.SweepExpired(now)
	if err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if removed != 2 {
		t.Fatalf("expected 2 swept records, got %d", removed)
	}

	if _, _, _, found, lookupErr := store.Lookup("expired-1"); lookupErr != nil || found {
		t.Fatalf("expired token must be gone after the sweep (found=%v err=%v)", found, lookupErr)
	}
	hash, state, _, found, lookupErr := store.Lookup("live")
	if lookupErr != nil || !found || hash != "hash-live" || state != IdempotencyStateCompleted {
		t.Fatalf("live token must survive the sweep (found=%v hash=%q state=%q err=%v)", found, hash, state, lookupErr)
	}

	removed, err = store.SweepExpired(now)
	if err != nil {
		t.Fatalf("second sweep: %v", err)
	}
	if removed != 0 {
		t.Fatalf("a second sweep must be a no-op, removed %d", removed)
	}
}

func TestIdempotencyStoreClaimStates(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewIdempotencyStore(st, "us-east-1")
	expires := time.Now().Add(5 * time.Minute)

	if err := store.Record("token", "hash-a", IdempotencyStateInProgress, expires, nil); err != nil {
		t.Fatalf("record in-progress claim: %v", err)
	}
	hash, state, _, found, lookupErr := store.Lookup("token")
	if lookupErr != nil || !found || hash != "hash-a" || state != IdempotencyStateInProgress {
		t.Fatalf("in-progress claim must be observable (found=%v hash=%q state=%q err=%v)", found, hash, state, lookupErr)
	}

	if err := store.Record("token", "hash-a", IdempotencyStateCompleted, expires, nil); err != nil {
		t.Fatalf("record completed claim: %v", err)
	}
	_, state, _, found, lookupErr = store.Lookup("token")
	if lookupErr != nil || !found || state != IdempotencyStateCompleted {
		t.Fatalf("completed claim must be observable (found=%v state=%q err=%v)", found, state, lookupErr)
	}

	if err := store.Delete("token"); err != nil {
		t.Fatalf("delete claim: %v", err)
	}
	if _, _, _, found, _ := store.Lookup("token"); found {
		t.Fatal("deleted claim must be absent")
	}
}

// A missing token is reported as absent with no error; undecodable record
// bytes are a storage failure, not absence — silently re-executing a
// request the caller believes deduplicated is the failure mode the error
// exists to prevent.
func TestIdempotencyLookupErrorSemantics(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	s := NewIdempotencyStore(st, "us-east-1")

	if _, _, _, found, err := s.Lookup("absent"); found || err != nil {
		t.Fatalf("absent token: found=%v err=%v, want false/nil", found, err)
	}

	if err := s.BaseStore.PutRaw("corrupt", []byte("not-proto")); err != nil {
		t.Fatalf("seed corrupt record: %v", err)
	}
	if _, _, _, _, err := s.Lookup("corrupt"); err == nil {
		t.Fatal("corrupt record decoded as absent")
	}

	if err := s.Record("live", "hash", IdempotencyStateCompleted, time.Now().Add(time.Hour), nil); err != nil {
		t.Fatalf("record: %v", err)
	}
	hash, state, _, found, err := s.Lookup("live")
	if err != nil || !found || hash != "hash" || state != IdempotencyStateCompleted {
		t.Fatalf("live record: hash=%q state=%q found=%v err=%v", hash, state, found, err)
	}
}

// A completed record carries the per-table replay read units through the
// round trip; a claim recorded without units reports none.
func TestIdempotencyRecordReplayReadUnits(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	s := NewIdempotencyStore(st, "us-east-1")

	units := map[string]float64{"T1": 6, "T2": 4}
	if err := s.Record("tok", "hash", IdempotencyStateCompleted, time.Now().Add(time.Hour), units); err != nil {
		t.Fatalf("record: %v", err)
	}
	_, state, got, found, err := s.Lookup("tok")
	if err != nil || !found || state != IdempotencyStateCompleted {
		t.Fatalf("lookup: found=%v state=%q err=%v", found, state, err)
	}
	if len(got) != 2 || got["T1"] != 6 || got["T2"] != 4 {
		t.Fatalf("replay read units = %v, want T1:6 T2:4", got)
	}

	if err := s.Record("bare", "hash", IdempotencyStateCompleted, time.Now().Add(time.Hour), nil); err != nil {
		t.Fatalf("record bare: %v", err)
	}
	if _, _, got, _, err := s.Lookup("bare"); err != nil || got != nil {
		t.Fatalf("bare record units = %v err=%v, want nil/nil", got, err)
	}
}

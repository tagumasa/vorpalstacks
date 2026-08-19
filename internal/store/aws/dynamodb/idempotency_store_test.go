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

	if err := store.Record("live", "hash-live", IdempotencyStateCompleted, now.Add(5*time.Minute)); err != nil {
		t.Fatalf("record live token: %v", err)
	}
	if err := store.Record("expired-1", "hash-one", IdempotencyStateCompleted, now.Add(-time.Minute)); err != nil {
		t.Fatalf("record expired token: %v", err)
	}
	if err := store.Record("expired-2", "hash-two", IdempotencyStateCompleted, now.Add(-2*time.Minute)); err != nil {
		t.Fatalf("record expired token: %v", err)
	}

	removed, err := store.SweepExpired(now)
	if err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if removed != 2 {
		t.Fatalf("expected 2 swept records, got %d", removed)
	}

	if _, _, found, lookupErr := store.Lookup("expired-1"); lookupErr != nil || found {
		t.Fatalf("expired token must be gone after the sweep (found=%v err=%v)", found, lookupErr)
	}
	hash, state, found, lookupErr := store.Lookup("live")
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

	if err := store.Record("token", "hash-a", IdempotencyStateInProgress, expires); err != nil {
		t.Fatalf("record in-progress claim: %v", err)
	}
	hash, state, found, lookupErr := store.Lookup("token")
	if lookupErr != nil || !found || hash != "hash-a" || state != IdempotencyStateInProgress {
		t.Fatalf("in-progress claim must be observable (found=%v hash=%q state=%q err=%v)", found, hash, state, lookupErr)
	}

	if err := store.Record("token", "hash-a", IdempotencyStateCompleted, expires); err != nil {
		t.Fatalf("record completed claim: %v", err)
	}
	_, state, found, lookupErr = store.Lookup("token")
	if lookupErr != nil || !found || state != IdempotencyStateCompleted {
		t.Fatalf("completed claim must be observable (found=%v state=%q err=%v)", found, state, lookupErr)
	}

	if err := store.Delete("token"); err != nil {
		t.Fatalf("delete claim: %v", err)
	}
	if _, _, found, _ := store.Lookup("token"); found {
		t.Fatal("deleted claim must be absent")
	}
}

package eventbus

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cockroachdb/pebble/v2"
)

func newTestDB(t *testing.T) *pebble.DB {
	t.Helper()
	dir := t.TempDir()
	db, err := pebble.Open(filepath.Join(dir, "test.db"), &pebble.Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestPebbleOutboxWriteRead(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	entry := &OutboxEntry{
		EventID:         "evt-1",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{"test": true}`),
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}

	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	read, err := store.Read(context.Background(), "evt-1")
	if err != nil {
		t.Fatal(err)
	}
	if read == nil {
		t.Fatal("expected entry, got nil")
	}
	if read.EventID != "evt-1" {
		t.Fatalf("expected event ID evt-1, got %s", read.EventID)
	}
	if read.Status != OutboxPending {
		t.Fatalf("expected status PENDING, got %s", read.Status)
	}
}

func TestPebbleOutboxReadNotFound(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	read, err := store.Read(context.Background(), "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if read != nil {
		t.Fatal("expected nil for nonexistent entry")
	}
}

func TestPebbleOutboxUpdateStatus(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	entry := &OutboxEntry{
		EventID:         "evt-2",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}

	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	updated, err := store.UpdateStatus(context.Background(), "evt-2", OutboxPending, OutboxProcessing)
	if err != nil {
		t.Fatal(err)
	}
	if !updated {
		t.Fatal("expected CAS to succeed")
	}

	updated, err = store.UpdateStatus(context.Background(), "evt-2", OutboxPending, OutboxProcessing)
	if err != nil {
		t.Fatal(err)
	}
	if updated {
		t.Fatal("expected CAS to fail (status already PROCESSING)")
	}

	read, _ := store.Read(context.Background(), "evt-2")
	if read.Status != OutboxProcessing {
		t.Fatalf("expected status PROCESSING, got %s", read.Status)
	}
}

func TestPebbleOutboxListPending(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	for i := 0; i < 5; i++ {
		entry := &OutboxEntry{
			EventID:         string(rune('a' + i)),
			EventType:       "service:invoke",
			Depth:           0,
			SerializedEvent: []byte(`{}`),
			Status:          OutboxPending,
			CreatedAt:       time.Now().UTC(),
			MaxRetries:      3,
			HandlerResults:  map[string]string{},
		}
		if err := store.Write(context.Background(), entry); err != nil {
			t.Fatal(err)
		}
	}

	pending, err := store.ListPending(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 5 {
		t.Fatalf("expected 5 pending entries, got %d", len(pending))
	}
}

// The pagination cursor excludes the entry it names, so consecutive pages
// tile the pending set without overlap or gaps.
func TestPebbleOutboxListPendingFromPagination(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	for i := 0; i < 5; i++ {
		entry := &OutboxEntry{
			EventID:         fmt.Sprintf("evt-%02d", i),
			EventType:       "service:invoke",
			Depth:           0,
			SerializedEvent: []byte(`{}`),
			Status:          OutboxPending,
			CreatedAt:       time.Now().UTC(),
			MaxRetries:      3,
			HandlerResults:  map[string]string{},
		}
		if err := store.Write(context.Background(), entry); err != nil {
			t.Fatal(err)
		}
	}

	page1, cursor1, err := store.ListPendingFrom(context.Background(), 2, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 2 || page1[0].EventID != "evt-00" || page1[1].EventID != "evt-01" {
		t.Fatalf("unexpected page 1: %s, %s", page1[0].EventID, page1[1].EventID)
	}
	if cursor1 == "" {
		t.Fatal("expected a continuation cursor for a full page")
	}

	page2, cursor2, err := store.ListPendingFrom(context.Background(), 2, cursor1)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 || page2[0].EventID != "evt-02" || page2[1].EventID != "evt-03" {
		t.Fatalf("unexpected page 2: %s, %s", page2[0].EventID, page2[1].EventID)
	}

	page3, cursor3, err := store.ListPendingFrom(context.Background(), 2, cursor2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page3) != 1 || page3[0].EventID != "evt-04" {
		t.Fatalf("unexpected page 3: %d entries", len(page3))
	}

	page4, _, err := store.ListPendingFrom(context.Background(), 2, cursor3)
	if err != nil {
		t.Fatal(err)
	}
	if len(page4) != 0 {
		t.Fatalf("expected empty page past the tail, got %d entries", len(page4))
	}
}

func TestPebbleOutboxDelete(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	entry := &OutboxEntry{
		EventID:         "evt-del",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}

	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	if err := store.Delete(context.Background(), "evt-del"); err != nil {
		t.Fatal(err)
	}

	read, _ := store.Read(context.Background(), "evt-del")
	if read != nil {
		t.Fatal("expected nil after delete")
	}
}

func TestPebbleOutboxDeleteNonexistent(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	if err := store.Delete(context.Background(), "nonexistent"); err != nil {
		t.Fatal(err)
	}
}

func TestPebbleOutboxUpdateEntry(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	entry := &OutboxEntry{
		EventID:         "evt-upd",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}

	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	entry.RetryCount = 1
	entry.LastError = "test error"
	if err := store.UpdateEntry(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	read, _ := store.Read(context.Background(), "evt-upd")
	if read.RetryCount != 1 {
		t.Fatalf("expected RetryCount 1, got %d", read.RetryCount)
	}
	if read.LastError != "test error" {
		t.Fatalf("expected LastError 'test error', got %s", read.LastError)
	}
}

func TestPebbleOutboxCleanup(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	now := time.Now().UTC()

	deliveredEntry := &OutboxEntry{
		EventID:         "delivered-old",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxDelivered,
		CreatedAt:       now.Add(-2 * time.Hour),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	deliveredTime := now.Add(-2 * time.Hour)
	deliveredEntry.DeliveredAt = &deliveredTime

	if err := store.Write(context.Background(), deliveredEntry); err != nil {
		t.Fatal(err)
	}

	failedEntry := &OutboxEntry{
		EventID:         "failed-old",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxFailed,
		CreatedAt:       now.Add(-25 * time.Hour),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}

	if err := store.Write(context.Background(), failedEntry); err != nil {
		t.Fatal(err)
	}

	recentEntry := &OutboxEntry{
		EventID:         "delivered-recent",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxDelivered,
		CreatedAt:       now.Add(-30 * time.Minute),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	recentTime := now.Add(-30 * time.Minute)
	recentEntry.DeliveredAt = &recentTime

	if err := store.Write(context.Background(), recentEntry); err != nil {
		t.Fatal(err)
	}

	pendingOld := &OutboxEntry{
		EventID:         "pending-old",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxPending,
		CreatedAt:       now.Add(-48 * time.Hour),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	if err := store.Write(context.Background(), pendingOld); err != nil {
		t.Fatal(err)
	}

	purged, err := store.Cleanup(context.Background(), now.Add(-1*time.Hour), now.Add(-24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	if purged != 2 {
		t.Fatalf("expected 2 purged entries, got %d", purged)
	}

	if _, err := store.Read(context.Background(), "delivered-old"); err != nil {
		t.Fatalf("read after cleanup: %v", err)
	}
	read, _ := store.Read(context.Background(), "delivered-old")
	if read != nil {
		t.Fatal("expected delivered-old to be cleaned up")
	}

	read2, _ := store.Read(context.Background(), "delivered-recent")
	if read2 == nil {
		t.Fatal("expected delivered-recent to still exist")
	}
}

// Pending entries are never purged by age: an undelivered event must survive
// cleanup until it delivers or exhausts its retry budget, otherwise a
// delivery backlog turns into silent event loss.
func TestCleanupNeverPurgesPending(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	now := time.Now().UTC()
	entry := &OutboxEntry{
		EventID:         "pending-ancient",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxPending,
		CreatedAt:       now.Add(-72 * time.Hour),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	purged, err := store.Cleanup(context.Background(), now.Add(-1*time.Hour), now.Add(-24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if purged != 0 {
		t.Fatalf("expected nothing purged, got %d", purged)
	}

	read, _ := store.Read(context.Background(), "pending-ancient")
	if read == nil {
		t.Fatal("pending entry was purged by age")
	}
	pending, err := store.ListPending(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 1 || pending[0].EventID != "pending-ancient" {
		t.Fatalf("pending scan lost the entry: %d entries", len(pending))
	}
}

// Every status transition must keep the record and its status index in
// step: after an entry walks Pending→Processing→Delivered the pending scan
// no longer sees it, retention purges it once the delivery timestamp ages
// out, and no index residue survives to haunt later scans.
func TestStatusIndexConsistencyAcrossTransitions(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	entry := &OutboxEntry{
		EventID:         "evt-transition",
		EventType:       "service:invoke",
		Depth:           0,
		SerializedEvent: []byte(`{}`),
		Status:          OutboxPending,
		CreatedAt:       time.Now().UTC(),
		MaxRetries:      3,
		HandlerResults:  map[string]string{},
	}
	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	pending, err := store.ListPending(context.Background(), 10)
	if err != nil || len(pending) != 1 {
		t.Fatalf("expected the fresh entry in the pending scan: %d, err=%v", len(pending), err)
	}

	if ok, err := store.UpdateStatus(context.Background(), "evt-transition", OutboxPending, OutboxProcessing); err != nil || !ok {
		t.Fatalf("Pending→Processing failed: ok=%v err=%v", ok, err)
	}
	pending, err = store.ListPending(context.Background(), 10)
	if err != nil || len(pending) != 0 {
		t.Fatalf("processing entry still visible in the pending scan: %d, err=%v", len(pending), err)
	}

	// Crash recovery resets Processing back to Pending via ResetStaleProcessing.
	reset, err := store.ResetStaleProcessing(context.Background())
	if err != nil || reset != 1 {
		t.Fatalf("ResetStaleProcessing failed: reset=%d err=%v", reset, err)
	}
	pending, err = store.ListPending(context.Background(), 10)
	if err != nil || len(pending) != 1 {
		t.Fatalf("recovered entry missing from the pending scan: %d, err=%v", len(pending), err)
	}

	if ok, err := store.UpdateStatus(context.Background(), "evt-transition", OutboxPending, OutboxProcessing); err != nil || !ok {
		t.Fatalf("Pending→Processing (2) failed: ok=%v err=%v", ok, err)
	}
	if ok, err := store.UpdateStatus(context.Background(), "evt-transition", OutboxProcessing, OutboxDelivered); err != nil || !ok {
		t.Fatalf("Processing→Delivered failed: ok=%v err=%v", ok, err)
	}
	pending, err = store.ListPending(context.Background(), 10)
	if err != nil || len(pending) != 0 {
		t.Fatalf("delivered entry still visible in the pending scan: %d, err=%v", len(pending), err)
	}

	// A recent delivery is retained; a delivery older than the retention
	// window is purged together with its index entry.
	now := time.Now().UTC()
	if purged, err := store.Cleanup(context.Background(), now.Add(-1*time.Hour), now.Add(-24*time.Hour)); err != nil || purged != 0 {
		t.Fatalf("recent delivered entry must be retained: purged=%d err=%v", purged, err)
	}
	if purged, err := store.Cleanup(context.Background(), now, now); err != nil || purged != 1 {
		t.Fatalf("aged delivered entry must be purged: purged=%d err=%v", purged, err)
	}

	// Re-writing the same event ID must be visible again (index recreated).
	if err := store.Write(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	pending, err = store.ListPending(context.Background(), 10)
	if err != nil || len(pending) != 1 {
		t.Fatalf("re-written entry missing from the pending scan: %d, err=%v", len(pending), err)
	}
}

func TestPebbleOutboxClose(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestPebbleOutboxConcurrentWrites(t *testing.T) {
	db := newTestDB(t)
	store := NewPebbleOutboxStore(db)

	const n = 100
	type result struct {
		id  int
		err error
	}
	results := make(chan result, n)
	for i := 0; i < n; i++ {
		go func(id int) {
			entry := &OutboxEntry{
				EventID:         fmt.Sprintf("concurrent-%d", id),
				EventType:       "service:invoke",
				Depth:           0,
				SerializedEvent: []byte(`{}`),
				Status:          OutboxPending,
				CreatedAt:       time.Now().UTC(),
				MaxRetries:      3,
				HandlerResults:  map[string]string{},
			}
			results <- result{id: id, err: store.Write(context.Background(), entry)}
		}(i)
	}

	writeErrors := 0
	for i := 0; i < n; i++ {
		r := <-results
		if r.err != nil {
			t.Errorf("concurrent write %d failed: %v", r.id, r.err)
			writeErrors++
		}
	}
	if writeErrors > 0 {
		t.Fatalf("%d/%d concurrent writes failed", writeErrors, n)
	}

	pending, err := store.ListPending(context.Background(), n+10)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != n {
		t.Fatalf("expected %d pending entries after concurrent writes, got %d", n, len(pending))
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

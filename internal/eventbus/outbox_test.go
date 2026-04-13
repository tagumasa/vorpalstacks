package eventbus

import (
	"context"
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

	done := make(chan struct{})
	const n = 100
	for i := 0; i < n; i++ {
		go func(id int) {
			defer func() { done <- struct{}{} }()
			entry := &OutboxEntry{
				EventID:         string(rune('0' + id)),
				EventType:       "service:invoke",
				Depth:           0,
				SerializedEvent: []byte(`{}`),
				Status:          OutboxPending,
				CreatedAt:       time.Now().UTC(),
				MaxRetries:      3,
				HandlerResults:  map[string]string{},
			}
			_ = store.Write(context.Background(), entry)
		}(i)
	}

	for i := 0; i < n; i++ {
		<-done
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

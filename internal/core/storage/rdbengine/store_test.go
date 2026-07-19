package rdbengine

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newTestStore creates a Pebble-backed rdbengine.Store for testing.
func newTestStore(t *testing.T) *Store {
	t.Helper()

	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })

	bucket := ps.Bucket("rdb").(storage.BatchBucket)
	store, err := New(bucket, DefaultOptions())
	if err != nil {
		t.Fatalf("rdbengine.New: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	ctx := context.Background()
	if err := store.CreateDatabase(ctx, "testdb"); err != nil {
		t.Fatalf("CreateDatabase: %v", err)
	}

	schema := &TableSchema{
		Name: "users",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "email", Type: ColumnTypeString},
			{Name: "name", Type: ColumnTypeString},
		},
	}
	if err := store.CreateTable(ctx, "testdb", schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	if err := store.CreateIndex(ctx, "testdb", "users", "idx_email", []string{"email"}, false); err != nil {
		t.Fatalf("CreateIndex: %v", err)
	}

	return store
}

func makeRow(id int32, email, name string) Row {
	return Row{
		"id":    {Type: ColumnTypeInt32, Value: id},
		"email": {Type: ColumnTypeString, Value: email},
		"name":  {Type: ColumnTypeString, Value: name},
	}
}

func pkOf(id int32) []byte {
	enc, err := encodeValue(ColumnValue{Type: ColumnTypeInt32, Value: id})
	if err != nil {
		panic(err)
	}
	return enc
}

// TestTxnDeleteRowRemovesIndexEntries verifies that TxnDeleteRow cleans up
// secondary index entries, preventing phantom index references after delete.
func TestTxnDeleteRowRemovesIndexEntries(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// Insert a row non-transactionally (creates index entries).
	row := makeRow(1, "alice@example.com", "Alice")
	pk := pkOf(1)
	if err := store.InsertRow(ctx, "testdb", "users", pk, row); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// Verify the index scan finds the row.
	it, err := store.ScanIndex(ctx, "testdb", "users", "idx_email", IndexScanOptions{})
	if err != nil {
		t.Fatalf("ScanIndex before delete: %v", err)
	}
	count := 0
	for it.Next() {
		count++
	}
	it.Close()
	if count != 1 {
		t.Fatalf("Expected 1 index entry before delete, got %d", count)
	}

	// Delete the row transactionally.
	tb := store.NewTxnBatch()
	if err := store.TxnDeleteRow(tb, "testdb", "users", pk); err != nil {
		tb.Rollback()
		t.Fatalf("TxnDeleteRow: %v", err)
	}
	if err := tb.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	// Verify the index scan finds NO entries after transactional delete.
	it2, err := store.ScanIndex(ctx, "testdb", "users", "idx_email", IndexScanOptions{})
	if err != nil {
		t.Fatalf("ScanIndex after delete: %v", err)
	}
	count = 0
	for it2.Next() {
		count++
	}
	it2.Close()
	if count != 0 {
		t.Errorf("Expected 0 index entries after txn delete, got %d (stale index entries remain)", count)
	}
}

// TestTxnUpdateRowRemovesOldIndexEntries verifies that TxnUpdateRow removes old
// index entries when a property value changes, preventing stale index lookups.
func TestTxnUpdateRowRemovesOldIndexEntries(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// Insert a row.
	row := makeRow(1, "alice@example.com", "Alice")
	pk := pkOf(1)
	if err := store.InsertRow(ctx, "testdb", "users", pk, row); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	// Update the email transactionally.
	updatedRow := makeRow(1, "bob@example.com", "Alice")
	tb := store.NewTxnBatch()
	if err := store.TxnUpdateRow(tb, "testdb", "users", pk, updatedRow); err != nil {
		tb.Rollback()
		t.Fatalf("TxnUpdateRow: %v", err)
	}
	if err := tb.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	// Scan index for the OLD email value — should find nothing.
	oldVal, err := encodeValue(ColumnValue{Type: ColumnTypeString, Value: "alice@example.com"})
	if err != nil {
		t.Fatalf("encodeValue old: %v", err)
	}
	it, err := store.ScanIndex(ctx, "testdb", "users", "idx_email", IndexScanOptions{
		Start: oldVal,
		End:   append(oldVal, 0xFF),
	})
	if err != nil {
		t.Fatalf("ScanIndex old: %v", err)
	}
	count := 0
	for it.Next() {
		count++
	}
	it.Close()
	if count != 0 {
		t.Errorf("Expected 0 index entries for old email after update, got %d (stale entries remain)", count)
	}

	// Scan index for the NEW email value — should find 1 entry.
	newVal, err := encodeValue(ColumnValue{Type: ColumnTypeString, Value: "bob@example.com"})
	if err != nil {
		t.Fatalf("encodeValue new: %v", err)
	}
	it2, err := store.ScanIndex(ctx, "testdb", "users", "idx_email", IndexScanOptions{
		Start: newVal,
		End:   append(newVal, 0xFF),
	})
	if err != nil {
		t.Fatalf("ScanIndex new: %v", err)
	}
	count = 0
	for it2.Next() {
		count++
	}
	it2.Close()
	if count != 1 {
		t.Errorf("Expected 1 index entry for new email after update, got %d", count)
	}
}

// TestNonTxnDeleteRowRemovesIndex is a baseline test confirming the
// non-transactional path correctly cleans up indexes.
func TestNonTxnDeleteRowRemovesIndex(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	row := makeRow(1, "alice@example.com", "Alice")
	pk := pkOf(1)
	if err := store.InsertRow(ctx, "testdb", "users", pk, row); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	if err := store.DeleteRow(ctx, "testdb", "users", pk); err != nil {
		t.Fatalf("DeleteRow: %v", err)
	}

	it, err := store.ScanIndex(ctx, "testdb", "users", "idx_email", IndexScanOptions{})
	if err != nil {
		t.Fatalf("ScanIndex: %v", err)
	}
	count := 0
	for it.Next() {
		count++
	}
	it.Close()
	if count != 0 {
		t.Errorf("Expected 0 index entries after delete, got %d", count)
	}
}

// TestTxnInsertRowCreatesIndex verifies that TxnInsertRow creates index entries.
func TestTxnInsertRowCreatesIndex(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	row := makeRow(1, "alice@example.com", "Alice")
	pk := pkOf(1)

	tb := store.NewTxnBatch()
	if err := store.TxnInsertRow(tb, "testdb", "users", pk, row); err != nil {
		tb.Rollback()
		t.Fatalf("TxnInsertRow: %v", err)
	}
	if err := tb.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	// Index scan should find the row.
	it, err := store.ScanIndex(ctx, "testdb", "users", "idx_email", IndexScanOptions{})
	if err != nil {
		t.Fatalf("ScanIndex: %v", err)
	}
	count := 0
	for it.Next() {
		count++
	}
	it.Close()
	if count != 1 {
		t.Errorf("Expected 1 index entry after txn insert, got %d", count)
	}
}

func TestGetRowNotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	_, err := store.GetRow(ctx, "testdb", "users", pkOf(999))
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestInsertRowDuplicate(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	row := makeRow(1, "alice@example.com", "Alice")
	pk := pkOf(1)
	if err := store.InsertRow(ctx, "testdb", "users", pk, row); err != nil {
		t.Fatalf("InsertRow first: %v", err)
	}
	if err := store.InsertRow(ctx, "testdb", "users", pk, row); err != ErrAlreadyExists {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}
}

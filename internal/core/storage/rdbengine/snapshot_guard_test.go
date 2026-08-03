package rdbengine

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// TestSnapshotReadOnlyGuard verifies that every write method on a
// snapshot Store returns errReadOnly. Without this guard, reads would
// observe the snapshot while writes silently went to the live backend,
// producing an inconsistent view across operations on the same handle.
func TestSnapshotReadOnlyGuard(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "guard_tbl",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "val", Type: ColumnTypeString},
		},
	}
	if err := store.CreateTable(ctx, "testdb", schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	snap, err := store.NewSnapshotReader()
	if err != nil {
		t.Fatalf("NewSnapshotReader: %v", err)
	}
	defer snap.Close()

	row := Row{"id": ColumnValue{Type: ColumnTypeInt32, Value: int32(1)}, "val": ColumnValue{Type: ColumnTypeString, Value: "x"}}
	pk := pkOfInt32(1)

	// --- DDL ---
	if err := snap.CreateDatabase(ctx, "readonly_db"); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.CreateDatabase: want errReadOnly, got %v", err)
	}
	if err := snap.CreateTable(ctx, "testdb", schema); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.CreateTable: want errReadOnly, got %v", err)
	}
	if err := snap.DropTable(ctx, "testdb", "guard_tbl"); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.DropTable: want errReadOnly, got %v", err)
	}
	if err := snap.CreateIndex(ctx, "testdb", "guard_tbl", "idx", []string{"val"}, false); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.CreateIndex: want errReadOnly, got %v", err)
	}
	if err := snap.DropIndex(ctx, "testdb", "guard_tbl", "idx"); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.DropIndex: want errReadOnly, got %v", err)
	}

	// --- DML ---
	if err := snap.InsertRow(ctx, "testdb", "guard_tbl", pk, row); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.InsertRow: want errReadOnly, got %v", err)
	}
	if err := snap.UpdateRow(ctx, "testdb", "guard_tbl", pk, row); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.UpdateRow: want errReadOnly, got %v", err)
	}
	if err := snap.DeleteRow(ctx, "testdb", "guard_tbl", pk); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.DeleteRow: want errReadOnly, got %v", err)
	}
	if _, err := snap.TruncateTable(ctx, "testdb", "guard_tbl"); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.TruncateTable: want errReadOnly, got %v", err)
	}

	// --- Auto-increment ---
	if err := snap.SetAutoIncrement("testdb", "guard_tbl", 100); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.SetAutoIncrement: want errReadOnly, got %v", err)
	}
	if _, err := snap.NextAutoIncrement("testdb", "guard_tbl"); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.NextAutoIncrement: want errReadOnly, got %v", err)
	}

	// --- Transaction ---
	tb := snap.NewTxnBatch()
	if tb != nil {
		t.Errorf("snap.NewTxnBatch: want nil on read-only store, got non-nil")
	}
	if err := snap.TxnInsertRow(tb, "testdb", "guard_tbl", pk, row); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.TxnInsertRow: want errReadOnly, got %v", err)
	}
	if err := snap.TxnUpdateRow(tb, "testdb", "guard_tbl", pk, row); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.TxnUpdateRow: want errReadOnly, got %v", err)
	}
	if err := snap.TxnDeleteRow(tb, "testdb", "guard_tbl", pk); !errors.Is(err, errReadOnly) {
		t.Errorf("snap.TxnDeleteRow: want errReadOnly, got %v", err)
	}
}

// TestSnapshotCloseReleases verifies that Close() on a snapshot Store
// releases the underlying Pebble snapshot handle. After Close, the
// source store must remain fully functional with no leaked locks or
// iterator caps.
func TestSnapshotCloseReleases(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "close_tbl",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
		},
	}
	if err := store.CreateTable(ctx, "testdb", schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// Create and close many snapshot readers in sequence. Pebble has a
	// default max open snapshot limit (typically very high, but a leak
	// would eventually exhaust file descriptors or hit internal caps
	// on long-running servers).
	for i := 0; i < 200; i++ {
		snap, err := store.NewSnapshotReader()
		if err != nil {
			t.Fatalf("iter %d: NewSnapshotReader: %v", i, err)
		}
		if err := snap.Close(); err != nil {
			t.Fatalf("iter %d: Close: %v", i, err)
		}
	}

	// Source store should still work.
	row := Row{"id": ColumnValue{Type: ColumnTypeInt32, Value: int32(42)}}
	pk := pkOfInt32(42)
	if err := store.InsertRow(ctx, "testdb", "close_tbl", pk, row); err != nil {
		t.Fatalf("InsertRow after snapshot close: %v", err)
	}
	got, err := store.GetRow(ctx, "testdb", "close_tbl", pk)
	if err != nil {
		t.Fatalf("GetRow: %v", err)
	}
	if got == nil {
		t.Fatal("GetRow returned nil after snapshot close cycle")
	}
}

// TestSnapshotCloseSnapshotIdempotent verifies that CloseSnapshot and
// Close can both be called without panicking, even in mixed order.
func TestSnapshotCloseSnapshotIdempotent(t *testing.T) {
	store := newTestStoreWithDB(t)

	snap, err := store.NewSnapshotReader()
	if err != nil {
		t.Fatalf("NewSnapshotReader: %v", err)
	}

	snap.CloseSnapshot()
	snap.CloseSnapshot() // double CloseSnapshot: should not panic
	snap.Close()         // Close after CloseSnapshot: should not panic
	snap.Close()         // double Close: should not panic
}

// ===== Snapshot read consistency =====

// TestSnapshotReadConsistency verifies that a snapshot reader observes
// a point-in-time view even after the source store is modified. This
// is the core guarantee that SnapshotData relies on for cross-table
// consistency.
func TestSnapshotReadConsistency(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "snap_tbl",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "val", Type: ColumnTypeString},
		},
	}
	if err := store.CreateTable(ctx, "testdb", schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// Insert row 1 before snapshot.
	row1 := Row{
		"id":  ColumnValue{Type: ColumnTypeInt32, Value: int32(1)},
		"val": ColumnValue{Type: ColumnTypeString, Value: "before"},
	}
	if err := store.InsertRow(ctx, "testdb", "snap_tbl", pkOfInt32(1), row1); err != nil {
		t.Fatalf("InsertRow 1: %v", err)
	}

	// Take snapshot.
	snap, err := store.NewSnapshotReader()
	if err != nil {
		t.Fatalf("NewSnapshotReader: %v", err)
	}
	defer snap.Close()

	// After snapshot, insert row 2 and update row 1.
	row2 := Row{
		"id":  ColumnValue{Type: ColumnTypeInt32, Value: int32(2)},
		"val": ColumnValue{Type: ColumnTypeString, Value: "after"},
	}
	if err := store.InsertRow(ctx, "testdb", "snap_tbl", pkOfInt32(2), row2); err != nil {
		t.Fatalf("InsertRow 2: %v", err)
	}
	row1Updated := Row{
		"id":  ColumnValue{Type: ColumnTypeInt32, Value: int32(1)},
		"val": ColumnValue{Type: ColumnTypeString, Value: "modified"},
	}
	if err := store.UpdateRow(ctx, "testdb", "snap_tbl", pkOfInt32(1), row1Updated); err != nil {
		t.Fatalf("UpdateRow 1: %v", err)
	}

	// Snapshot should see row 1 with "before" value and NOT see row 2.
	snapRow1, err := snap.GetRow(ctx, "testdb", "snap_tbl", pkOfInt32(1))
	if err != nil {
		t.Fatalf("snap.GetRow 1: %v", err)
	}
	if snapRow1 == nil {
		t.Fatal("snap.GetRow 1: expected row, got nil")
	}
	valCol, ok := snapRow1["val"]
	if !ok {
		t.Fatal("snap.GetRow 1: missing 'val' column")
	}
	if s, ok := valCol.Value.(string); !ok || s != "before" {
		t.Errorf("snap.GetRow 1 val: want 'before', got %v", valCol.Value)
	}

	snapRow2, err := snap.GetRow(ctx, "testdb", "snap_tbl", pkOfInt32(2))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("snap.GetRow 2: want ErrNotFound (row inserted after snapshot), got row=%v err=%v", snapRow2, err)
	}

	// Source store should see the updated state.
	liveRow1, err := store.GetRow(ctx, "testdb", "snap_tbl", pkOfInt32(1))
	if err != nil {
		t.Fatalf("live.GetRow 1: %v", err)
	}
	liveVal, ok := liveRow1["val"]
	if !ok {
		t.Fatal("live.GetRow 1: missing 'val' column")
	}
	if s, ok := liveVal.Value.(string); !ok || s != "modified" {
		t.Errorf("live.GetRow 1 val: want 'modified', got %v", liveVal.Value)
	}

	// Snapshot scan should return only row 1.
	iter, err := snap.ScanRows(ctx, "testdb", "snap_tbl", ScanOptions{})
	if err != nil {
		t.Fatalf("snap.ScanRows: %v", err)
	}
	count := 0
	for iter.Next() {
		count++
	}
	if err := iter.Error(); err != nil {
		t.Fatalf("snap.ScanRows iter error: %v", err)
	}
	iter.Close()
	if count != 1 {
		t.Errorf("snap.ScanRows: want 1 row, got %d", count)
	}
}

// TestSnapshotEmptyInstance verifies that snapshotting a store with no
// databases (a freshly created instance) does not error.
func TestSnapshotEmptyInstance(t *testing.T) {
	ctx := context.Background()

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

	snap, err := store.NewSnapshotReader()
	if err != nil {
		t.Fatalf("NewSnapshotReader: %v", err)
	}
	defer snap.Close()

	dbs, err := snap.ListDatabases(ctx)
	if err != nil {
		t.Fatalf("ListDatabases on empty snapshot: %v", err)
	}
	if len(dbs) != 0 {
		t.Errorf("empty snapshot ListDatabases: want 0, got %d", len(dbs))
	}
}

package rdbengine

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// newTestStoreWithDB creates a fresh store with a database but no tables,
// giving individual tests full control over schema setup.
func newTestStoreWithDB(t *testing.T) *Store {
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
	return store
}

// createTable is a helper that creates a table and returns its schema.
func createTable(t *testing.T, store *Store, schema *TableSchema) *TableSchema {
	t.Helper()
	ctx := context.Background()
	if err := store.CreateTable(ctx, "testdb", schema); err != nil {
		t.Fatalf("CreateTable %s: %v", schema.Name, err)
	}
	return schema
}

// pkOfInt32 encodes an int32 PK column value for simple single-column PKs.
func pkOfInt32(id int32) []byte {
	enc, err := encodeValue(ColumnValue{Type: ColumnTypeInt32, Value: id})
	if err != nil {
		panic(err)
	}
	return enc
}

// ===== Primary Key Type Coverage =====

func TestPKAllTypes(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	tests := []struct {
		name   string
		colDef ColumnDef
		pkVal  interface{}
	}{
		{"int32", ColumnDef{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true}, int32(42)},
		{"int64", ColumnDef{Name: "id", Type: ColumnTypeInt64, PrimaryKey: true}, int64(99)},
		{"string", ColumnDef{Name: "id", Type: ColumnTypeString, PrimaryKey: true}, "key1"},
		{"bool", ColumnDef{Name: "id", Type: ColumnTypeBool, PrimaryKey: true}, true},
		{"float64", ColumnDef{Name: "id", Type: ColumnTypeFloat64, PrimaryKey: true}, 3.14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tblName := "pk_" + tt.name
			schema := &TableSchema{
				Name:    tblName,
				Columns: []ColumnDef{tt.colDef, {Name: "v", Type: ColumnTypeString}},
			}
			createTable(t, store, schema)

			row := Row{
				"id": {Type: tt.colDef.Type, Value: tt.pkVal},
				"v":  {Type: ColumnTypeString, Value: "x"},
			}
			pk, err := EncodePK(schema, row)
			if err != nil {
				t.Fatalf("EncodePK: %v", err)
			}
			if err := store.InsertRow(ctx, "testdb", tblName, pk, row); err != nil {
				t.Fatalf("InsertRow: %v", err)
			}

			got, err := store.GetRow(ctx, "testdb", tblName, pk)
			if err != nil {
				t.Fatalf("GetRow: %v", err)
			}
			if got["v"].Value != "x" {
				t.Fatalf("expected v=x, got %v", got["v"].Value)
			}
		})
	}
}

// ===== NULL Primary Key Rejection =====

func TestNullPKRejected(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "npk",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "v", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	// Row with NULL primary key (Value=nil)
	row := Row{
		"id": {Type: ColumnTypeInt32, Value: nil},
		"v":  {Type: ColumnTypeString, Value: "x"},
	}
	_, err := EncodePK(schema, row)
	if err == nil {
		t.Fatal("expected EncodePK to reject NULL PK, but got nil error")
	}

	// Direct insert with nil PK via encodeValue should also error
	_, err = encodeValue(ColumnValue{Type: ColumnTypeInt32, Value: nil})
	if err == nil {
		t.Fatal("expected encodeValue to reject nil value")
	}

	// Verify the table is still usable with a valid PK
	row2 := Row{
		"id": {Type: ColumnTypeInt32, Value: int32(1)},
		"v":  {Type: ColumnTypeString, Value: "valid"},
	}
	pk2, _ := EncodePK(schema, row2)
	if err := store.InsertRow(ctx, "testdb", "npk", pk2, row2); err != nil {
		t.Fatalf("InsertRow valid: %v", err)
	}
}

// ===== NULL in UNIQUE Index (Multiple NULLs Allowed) =====

func TestUniqueIndexMultipleNullsAllowed(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "uniq_null",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "email", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)
	if err := store.CreateIndex(ctx, "testdb", "uniq_null", "uidx", []string{"email"}, true); err != nil {
		t.Fatalf("CreateIndex: %v", err)
	}

	// Insert two rows with NULL email — both should succeed (SQL standard: NULL != NULL)
	for i := int32(1); i <= 2; i++ {
		row := Row{
			"id":    {Type: ColumnTypeInt32, Value: i},
			"email": {Type: ColumnTypeString, Value: nil},
		}
		pk := pkOfInt32(i)
		if err := store.InsertRow(ctx, "testdb", "uniq_null", pk, row); err != nil {
			t.Fatalf("InsertRow %d with NULL email should succeed: %v", i, err)
		}
	}

	// Insert a third row with a non-NULL email
	row3 := Row{
		"id":    {Type: ColumnTypeInt32, Value: 3},
		"email": {Type: ColumnTypeString, Value: "a@b.com"},
	}
	pk3 := pkOfInt32(3)
	if err := store.InsertRow(ctx, "testdb", "uniq_null", pk3, row3); err != nil {
		t.Fatalf("InsertRow with non-NULL email should succeed: %v", err)
	}

	// Insert a fourth row with the SAME non-NULL email — should fail (unique violation)
	row4 := Row{
		"id":    {Type: ColumnTypeInt32, Value: 4},
		"email": {Type: ColumnTypeString, Value: "a@b.com"},
	}
	pk4 := pkOfInt32(4)
	if err := store.InsertRow(ctx, "testdb", "uniq_null", pk4, row4); err == nil {
		t.Fatal("InsertRow with duplicate non-NULL email should fail")
	}
}

// ===== Composite Primary Key =====

func TestCompositePK(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "composite",
		Columns: []ColumnDef{
			{Name: "tenant", Type: ColumnTypeString, PrimaryKey: true},
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "data", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	rows := []Row{
		{"tenant": {Type: ColumnTypeString, Value: "a"}, "id": {Type: ColumnTypeInt32, Value: 1}, "data": {Type: ColumnTypeString, Value: "r1"}},
		{"tenant": {Type: ColumnTypeString, Value: "a"}, "id": {Type: ColumnTypeInt32, Value: 2}, "data": {Type: ColumnTypeString, Value: "r2"}},
		{"tenant": {Type: ColumnTypeString, Value: "b"}, "id": {Type: ColumnTypeInt32, Value: 1}, "data": {Type: ColumnTypeString, Value: "r3"}},
	}

	for _, row := range rows {
		pk, err := EncodePK(schema, row)
		if err != nil {
			t.Fatalf("EncodePK: %v", err)
		}
		if err := store.InsertRow(ctx, "testdb", "composite", pk, row); err != nil {
			t.Fatalf("InsertRow: %v", err)
		}
	}

	// Lookup (a,2)
	target := Row{
		"tenant": {Type: ColumnTypeString, Value: "a"},
		"id":     {Type: ColumnTypeInt32, Value: 2},
	}
	pk, _ := EncodePK(schema, target)
	got, err := store.GetRow(ctx, "testdb", "composite", pk)
	if err != nil {
		t.Fatalf("GetRow: %v", err)
	}
	if got["data"].Value != "r2" {
		t.Fatalf("expected data=r2, got %v", got["data"].Value)
	}

	// Duplicate composite PK should fail
	dup := Row{
		"tenant": {Type: ColumnTypeString, Value: "a"},
		"id":     {Type: ColumnTypeInt32, Value: 1},
		"data":   {Type: ColumnTypeString, Value: "dup"},
	}
	pkDup, _ := EncodePK(schema, dup)
	if err := store.InsertRow(ctx, "testdb", "composite", pkDup, dup); err == nil {
		t.Fatal("duplicate composite PK should fail")
	}
}

// ===== Transaction Commit and Rollback =====

func TestTransactionCommit(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "txn_commit",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "v", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	batch := store.NewTxnBatch()
	row := Row{"id": {Type: ColumnTypeInt32, Value: 1}, "v": {Type: ColumnTypeString, Value: "committed"}}
	pk := pkOfInt32(1)

	if err := store.TxnInsertRow(batch, "testdb", "txn_commit", pk, row); err != nil {
		t.Fatalf("TxnInsertRow: %v", err)
	}
	if err := batch.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	got, err := store.GetRow(ctx, "testdb", "txn_commit", pk)
	if err != nil {
		t.Fatalf("GetRow after commit: %v", err)
	}
	if got["v"].Value != "committed" {
		t.Fatalf("expected v=committed, got %v", got["v"].Value)
	}
}

func TestTransactionRollback(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "txn_rollback",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "v", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	batch := store.NewTxnBatch()
	row := Row{"id": {Type: ColumnTypeInt32, Value: 1}, "v": {Type: ColumnTypeString, Value: "rolled-back"}}
	pk := pkOfInt32(1)

	if err := store.TxnInsertRow(batch, "testdb", "txn_rollback", pk, row); err != nil {
		t.Fatalf("TxnInsertRow: %v", err)
	}
	batch.Rollback()

	_, err := store.GetRow(ctx, "testdb", "txn_rollback", pk)
	if err == nil {
		t.Fatal("expected error after rollback (row should not exist)")
	}
}

// ===== Index Backfill on CreateIndex =====

func TestCreateIndexBackfill(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "backfill",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "cat", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	// Insert rows BEFORE creating the index
	for i := int32(1); i <= 5; i++ {
		row := Row{
			"id":  {Type: ColumnTypeInt32, Value: i},
			"cat": {Type: ColumnTypeString, Value: fmt.Sprintf("c%d", i%3)},
		}
		pk := pkOfInt32(i)
		if err := store.InsertRow(ctx, "testdb", "backfill", pk, row); err != nil {
			t.Fatalf("InsertRow %d: %v", i, err)
		}
	}

	// Now create a non-unique index on "cat" — should backfill all 5 rows
	if err := store.CreateIndex(ctx, "testdb", "backfill", "idx_cat", []string{"cat"}, false); err != nil {
		t.Fatalf("CreateIndex: %v", err)
	}

	// Scan the index — all 5 rows should be present
	iter, err := store.ScanIndex(ctx, "testdb", "backfill", "idx_cat", IndexScanOptions{})
	if err != nil {
		t.Fatalf("ScanIndex: %v", err)
	}
	defer iter.Close()

	count := 0
	for iter.Next() {
		count++
	}
	if err := iter.Error(); err != nil {
		t.Fatalf("iter.Error: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected 5 rows in index after backfill, got %d", count)
	}
}

// ===== TruncateTable =====

func TestTruncateTable(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "trunc",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "v", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	for i := int32(1); i <= 10; i++ {
		row := Row{"id": {Type: ColumnTypeInt32, Value: i}, "v": {Type: ColumnTypeString, Value: "x"}}
		pk := pkOfInt32(i)
		if err := store.InsertRow(ctx, "testdb", "trunc", pk, row); err != nil {
			t.Fatalf("InsertRow %d: %v", i, err)
		}
	}

	affected, err := store.TruncateTable(ctx, "testdb", "trunc")
	if err != nil {
		t.Fatalf("TruncateTable: %v", err)
	}
	// TruncateTable returns 0 affected (MySQL TRUNCATE semantics)
	if affected != 0 {
		t.Fatalf("expected 0 affected, got %d", affected)
	}

	// Verify table is empty
	iter, err := store.ScanRows(ctx, "testdb", "trunc", ScanOptions{})
	if err != nil {
		t.Fatalf("ScanRows: %v", err)
	}
	defer iter.Close()
	count := 0
	for iter.Next() {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 rows after truncate, got %d", count)
	}

	// Table should still exist and accept new inserts
	row := Row{"id": {Type: ColumnTypeInt32, Value: 100}, "v": {Type: ColumnTypeString, Value: "new"}}
	pk := pkOfInt32(100)
	if err := store.InsertRow(ctx, "testdb", "trunc", pk, row); err != nil {
		t.Fatalf("InsertRow after truncate: %v", err)
	}
}

// ===== DropDatabase Atomicity =====

func TestDropDatabaseAtomic(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	for _, tbl := range []string{"t1", "t2"} {
		schema := &TableSchema{
			Name: tbl,
			Columns: []ColumnDef{
				{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			},
		}
		createTable(t, store, schema)
		row := Row{"id": {Type: ColumnTypeInt32, Value: 1}}
		pk := pkOfInt32(1)
		if err := store.InsertRow(ctx, "testdb", tbl, pk, row); err != nil {
			t.Fatalf("InsertRow %s: %v", tbl, err)
		}
	}

	if err := store.DropDatabase(ctx, "testdb"); err != nil {
		t.Fatalf("DropDatabase: %v", err)
	}

	dbs, err := store.ListDatabases(ctx)
	if err != nil {
		t.Fatalf("ListDatabases: %v", err)
	}
	for _, db := range dbs {
		if db == "testdb" {
			t.Fatal("testdb should have been dropped")
		}
	}

	_, err = store.ListTables(ctx, "testdb")
	if err == nil {
		t.Fatal("expected error listing tables of dropped database")
	}
}

// ===== Identifier Validation =====

func TestValidateIdentifierRejectsInvalid(t *testing.T) {
	ctx := context.Background()

	tests := []string{
		"",
		strings.Repeat("a", 65), // exceeds 64-byte limit
		"has/slash",
		"has\x00null",
		"has\xffbyte",
	}

	for _, badName := range tests {
		t.Run("reject", func(t *testing.T) {
			store := newTestStoreWithDB(t)
			err := store.CreateDatabase(ctx, badName)
			if err == nil {
				t.Fatalf("expected error for invalid database name %q", badName)
			}
		})
	}
}

func TestValidateIdentifierAcceptsValid(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	validNames := []string{
		"mytbl",
		"MyTable123",
		"_private",
		"table-with-dash",
		strings.Repeat("a", 64), // exactly 64 bytes
	}

	for _, name := range validNames {
		schema := &TableSchema{
			Name: name,
			Columns: []ColumnDef{
				{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			},
		}
		if err := store.CreateTable(ctx, "testdb", schema); err != nil {
			t.Fatalf("expected CreateTable with valid name %q to succeed: %v", name, err)
		}
	}
}

// ===== Timestamp Values =====

func TestTimestampValue(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "ts_test",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "ts", Type: ColumnTypeTimestamp},
		},
	}
	createTable(t, store, schema)

	now := time.Now().UTC()
	row := Row{
		"id": {Type: ColumnTypeInt32, Value: int32(1)},
		"ts": {Type: ColumnTypeTimestamp, Value: now},
	}
	pk := pkOfInt32(1)
	if err := store.InsertRow(ctx, "testdb", "ts_test", pk, row); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	got, err := store.GetRow(ctx, "testdb", "ts_test", pk)
	if err != nil {
		t.Fatalf("GetRow: %v", err)
	}

	gotTime, ok := got["ts"].Value.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", got["ts"].Value)
	}
	if !gotTime.Equal(now) {
		t.Fatalf("timestamp mismatch: expected %v, got %v", now, gotTime)
	}
}

// ===== ScanRows Full Table =====

func TestScanRowsFullTable(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "scan_tbl",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "v", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)

	for i := int32(1); i <= 20; i++ {
		row := Row{
			"id": {Type: ColumnTypeInt32, Value: i},
			"v":  {Type: ColumnTypeString, Value: fmt.Sprintf("val%d", i)},
		}
		pk := pkOfInt32(i)
		if err := store.InsertRow(ctx, "testdb", "scan_tbl", pk, row); err != nil {
			t.Fatalf("InsertRow %d: %v", i, err)
		}
	}

	iter, err := store.ScanRows(ctx, "testdb", "scan_tbl", ScanOptions{})
	if err != nil {
		t.Fatalf("ScanRows: %v", err)
	}
	defer iter.Close()

	count := 0
	for iter.Next() {
		count++
	}
	if err := iter.Error(); err != nil {
		t.Fatalf("iter.Error: %v", err)
	}
	if count != 20 {
		t.Fatalf("expected 20 rows, got %d", count)
	}
}

// ===== Non-Txn Update Preserves Index Consistency =====

func TestUpdateRowNonTxnPreservesIndex(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "upd_idx",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
			{Name: "email", Type: ColumnTypeString},
		},
	}
	createTable(t, store, schema)
	if err := store.CreateIndex(ctx, "testdb", "upd_idx", "idx_email", []string{"email"}, false); err != nil {
		t.Fatalf("CreateIndex: %v", err)
	}

	row := Row{
		"id":    {Type: ColumnTypeInt32, Value: 1},
		"email": {Type: ColumnTypeString, Value: "old@test.com"},
	}
	pk := pkOfInt32(1)
	if err := store.InsertRow(ctx, "testdb", "upd_idx", pk, row); err != nil {
		t.Fatalf("InsertRow: %v", err)
	}

	updatedRow := Row{
		"id":    {Type: ColumnTypeInt32, Value: 1},
		"email": {Type: ColumnTypeString, Value: "new@test.com"},
	}
	if err := store.UpdateRow(ctx, "testdb", "upd_idx", pk, updatedRow); err != nil {
		t.Fatalf("UpdateRow: %v", err)
	}

	iter, err := store.ScanIndex(ctx, "testdb", "upd_idx", "idx_email", IndexScanOptions{})
	if err != nil {
		t.Fatalf("ScanIndex: %v", err)
	}
	defer iter.Close()

	count := 0
	for iter.Next() {
		count++
		email := iter.Row()["email"]
		if email.Value != "new@test.com" {
			t.Fatalf("expected new@test.com, got %v", email.Value)
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 index entry after update, got %d", count)
	}
}

// ===== AutoIncrement =====

func TestAutoIncrement(t *testing.T) {
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "autoinc",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeSerial, PrimaryKey: true},
		},
	}
	createTable(t, store, schema)

	v1, err := store.NextAutoIncrement("testdb", "autoinc")
	if err != nil {
		t.Fatalf("NextAutoIncrement: %v", err)
	}
	if v1 != 1 {
		t.Fatalf("expected first auto-increment=1, got %d", v1)
	}

	v2, err := store.NextAutoIncrement("testdb", "autoinc")
	if err != nil {
		t.Fatalf("NextAutoIncrement: %v", err)
	}
	if v2 != 2 {
		t.Fatalf("expected second auto-increment=2, got %d", v2)
	}

	cur, err := store.GetAutoIncrement("testdb", "autoinc")
	if err != nil {
		t.Fatalf("GetAutoIncrement: %v", err)
	}
	// GetAutoIncrement returns the next-to-be-assigned value, which is
	// one past the last returned value (2+1=3).
	if cur != 3 {
		t.Fatalf("expected current auto-increment=3, got %d", cur)
	}

	if err := store.SetAutoIncrement("testdb", "autoinc", 100); err != nil {
		t.Fatalf("SetAutoIncrement: %v", err)
	}
	// SetAutoIncrement(100) stores 100 as the next value to return.
	v3, _ := store.NextAutoIncrement("testdb", "autoinc")
	if v3 != 100 {
		t.Fatalf("expected auto-increment=100 after set 100, got %d", v3)
	}
}

// ===== DropTable =====

func TestDropTable(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	schema := &TableSchema{
		Name: "drop_me",
		Columns: []ColumnDef{
			{Name: "id", Type: ColumnTypeInt32, PrimaryKey: true},
		},
	}
	createTable(t, store, schema)

	tables, err := store.ListTables(ctx, "testdb")
	if err != nil {
		t.Fatalf("ListTables: %v", err)
	}
	found := false
	for _, tbl := range tables {
		if tbl == "drop_me" {
			found = true
		}
	}
	if !found {
		t.Fatal("table should exist before drop")
	}

	if err := store.DropTable(ctx, "testdb", "drop_me"); err != nil {
		t.Fatalf("DropTable: %v", err)
	}

	tables, _ = store.ListTables(ctx, "testdb")
	for _, tbl := range tables {
		if tbl == "drop_me" {
			t.Fatal("table should not exist after drop")
		}
	}

	_, err = store.GetTableSchema(ctx, "testdb", "drop_me")
	if err == nil {
		t.Fatal("expected error getting schema of dropped table")
	}
}

// ===== GetRow on Non-Existent Table =====

func TestGetRowNonExistentTable(t *testing.T) {
	ctx := context.Background()
	store := newTestStoreWithDB(t)

	pk := []byte{1, 2, 3}
	_, err := store.GetRow(ctx, "testdb", "nonexistent", pk)
	if err == nil {
		t.Fatal("expected error getting row from non-existent table")
	}
}

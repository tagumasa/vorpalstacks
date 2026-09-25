package dynamodb

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// createMetricsTestTable creates a single-hash-key table for the metric
// flush tests.
func createMetricsTestTable(t *testing.T, store *DynamoDBStore, name string) {
	t.Helper()
	if _, err := store.Tables().Create(CreateTableParams{
		Name:                 name,
		KeySchema:            []*KeySchemaElement{{AttributeName: "pk", KeyType: KeyTypeHash}},
		AttributeDefinitions: []*AttributeDefinition{{AttributeName: "pk", AttributeType: ScalarAttributeTypeS}},
		BillingMode:          BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table %s: %v", name, err)
	}
}

// TestRestoreChunkAggregatesCounters mirrors a chunked restore: many items
// written with per-item counter calls inside ONE transaction. The queued
// deltas must aggregate — the chunk lands the full item count and size, not
// the last write's increment.
func TestRestoreChunkAggregatesCounters(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	createMetricsTestTable(t, store, "RestoreTbl")

	const items = 500
	const itemSize = 7
	err = store.Update(context.Background(), func(txn *DynamoDBTxn) error {
		for i := 0; i < items; i++ {
			key := map[string]*AttributeValue{"pk": strAttr(fmt.Sprintf("k%04d", i))}
			if err := txn.PutItem("RestoreTbl", key, key); err != nil {
				return err
			}
			if err := txn.UpdateItemCount("RestoreTbl", 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize("RestoreTbl", itemSize); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("write restore chunk: %v", err)
	}

	table, err := store.Tables().Get("RestoreTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	if table.ItemCount != items {
		t.Fatalf("ItemCount = %d, want %d", table.ItemCount, items)
	}
	if table.TableSizeBytes != items*itemSize {
		t.Fatalf("TableSizeBytes = %d, want %d", table.TableSizeBytes, items*itemSize)
	}
}

// TestTwoPhaseTransactionMetricDeltas drives two item-writing executors in
// one two-phase transaction — the TransactWriteItems shape — and applies
// the collected deltas after the commit: two puts to one table must land
// exactly +2 items, not the +1 a per-operation read-modify-write produced.
func TestTwoPhaseTransactionMetricDeltas(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	createMetricsTestTable(t, store, "TransactTbl")

	twoPhase := store.Storage().TwoPhaseTransaction()
	var metricDeltas map[string]TableMetricDelta
	for i := 0; i < 2; i++ {
		i := i
		twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
			dbTxn := store.NewTxn(txn)
			key := map[string]*AttributeValue{"pk": strAttr(fmt.Sprintf("op%d", i))}
			if err := dbTxn.PutItem("TransactTbl", key, key); err != nil {
				return err
			}
			if err := dbTxn.UpdateItemCount("TransactTbl", 1); err != nil {
				return err
			}
			if err := dbTxn.UpdateTableSize("TransactTbl", 5); err != nil {
				return err
			}
			metricDeltas = MergeTableMetricDeltas(metricDeltas, dbTxn.TakeTableMetricDeltas())
			return nil
		}))
	}
	if err := twoPhase.Commit(context.Background()); err != nil {
		t.Fatalf("two-phase commit: %v", err)
	}
	store.FlushTableMetrics(metricDeltas)

	table, err := store.Tables().Get("TransactTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	if table.ItemCount != 2 {
		t.Fatalf("ItemCount = %d, want 2", table.ItemCount)
	}
	if table.TableSizeBytes != 10 {
		t.Fatalf("TableSizeBytes = %d, want 10", table.TableSizeBytes)
	}
}

// TestConcurrentWritesExactCounters hammers one table with concurrent
// single-item transactions, the concurrent-PutItem shape: the per-table
// locked flush must land every increment — the final counters match the
// number of committed writes exactly.
func TestConcurrentWritesExactCounters(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	createMetricsTestTable(t, store, "ConcurrentTbl")

	const writers = 16
	const perWriter = 25
	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perWriter; i++ {
				key := map[string]*AttributeValue{"pk": strAttr(fmt.Sprintf("w%d-i%d", w, i))}
				err := store.Update(context.Background(), func(txn *DynamoDBTxn) error {
					if err := txn.PutItem("ConcurrentTbl", key, key); err != nil {
						return err
					}
					if err := txn.UpdateItemCount("ConcurrentTbl", 1); err != nil {
						return err
					}
					return txn.UpdateTableSize("ConcurrentTbl", 3)
				})
				if err != nil {
					t.Errorf("writer %d iteration %d: %v", w, i, err)
					return
				}
			}
		}()
	}
	wg.Wait()

	table, err := store.Tables().Get("ConcurrentTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	want := int64(writers * perWriter)
	if table.ItemCount != want {
		t.Fatalf("ItemCount = %d, want %d", table.ItemCount, want)
	}
	if table.TableSizeBytes != want*3 {
		t.Fatalf("TableSizeBytes = %d, want %d", table.TableSizeBytes, want*3)
	}
}

// TestMetricFlushSkipsASameNameSuccessor pins the generation guard of the
// metric flush: deltas queued against one generation of a table name apply
// only to that generation's record — a same-name successor created between
// the carrying write's commit and the flush inherits nothing, the same
// not-an-error the deleted-table half of the window applies.
func TestMetricFlushSkipsASameNameSuccessor(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	createMetricsTestTable(t, store, "GenTbl")
	first, err := store.Tables().Get("GenTbl")
	if err != nil {
		t.Fatalf("first generation: %v", err)
	}

	// The successor: delete and recreate under the same name, through the
	// production deletion path (the cascade inside one transaction).
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteTableCascade("GenTbl")
	}); err != nil {
		t.Fatalf("delete first generation: %v", err)
	}
	createMetricsTestTable(t, store, "GenTbl")
	second, err := store.Tables().Get("GenTbl")
	if err != nil {
		t.Fatalf("second generation: %v", err)
	}
	if second.TableId == first.TableId {
		t.Fatal("the recreated table kept the first generation's id")
	}

	// The superseded generation's deltas are dropped, not inherited.
	if err := store.tables.applyMetricDeltas("GenTbl", first.TableId, 7, 700); err != nil {
		t.Fatalf("superseded flush errored: %v", err)
	}
	if second, err := store.Tables().Get("GenTbl"); err != nil || second.ItemCount != 0 || second.TableSizeBytes != 0 {
		t.Fatalf("successor inherited the superseded generation's counters: %+v (%v)", second, err)
	}

	// The live generation's deltas apply.
	if err := store.tables.applyMetricDeltas("GenTbl", second.TableId, 3, 300); err != nil {
		t.Fatalf("live flush errored: %v", err)
	}
	live, err := store.Tables().Get("GenTbl")
	if err != nil {
		t.Fatalf("reload successor: %v", err)
	}
	if live.ItemCount != 3 || live.TableSizeBytes != 300 {
		t.Fatalf("live generation counters = %d/%d, want 3/300", live.ItemCount, live.TableSizeBytes)
	}
}

// TestMergedTransactionDeltasKeepGeneration pins the merge half of the
// generation guard: the two-phase transaction plane folds per-operation
// queued deltas through MergeTableMetricDeltas before flushing, and the
// fold must carry the queued TableId — a same-name successor created
// between the carrying transaction's commit and the flush inherits
// nothing, exactly as the single-transaction flush path does.
func TestMergedTransactionDeltasKeepGeneration(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	createMetricsTestTable(t, store, "MergeGenTbl")
	first, err := store.Tables().Get("MergeGenTbl")
	if err != nil {
		t.Fatalf("first generation: %v", err)
	}

	// One two-phase operation queues its deltas against the first
	// generation; the wrapper merge is the TransactWriteItems shape.
	twoPhase := store.Storage().TwoPhaseTransaction()
	var metricDeltas map[string]TableMetricDelta
	twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
		dbTxn := store.NewTxn(txn)
		key := map[string]*AttributeValue{"pk": strAttr("orphaned")}
		if err := dbTxn.PutItem("MergeGenTbl", key, key); err != nil {
			return err
		}
		if err := dbTxn.UpdateItemCount("MergeGenTbl", 1); err != nil {
			return err
		}
		if err := dbTxn.UpdateTableSize("MergeGenTbl", 9); err != nil {
			return err
		}
		metricDeltas = MergeTableMetricDeltas(metricDeltas, dbTxn.TakeTableMetricDeltas())
		return nil
	}))
	if err := twoPhase.Commit(context.Background()); err != nil {
		t.Fatalf("two-phase commit: %v", err)
	}

	// The successor replaces the table generation before the flush runs,
	// through the production deletion path (the cascade inside one
	// transaction).
	if err := store.Update(t.Context(), func(txn *DynamoDBTxn) error {
		return txn.DeleteTableCascade("MergeGenTbl")
	}); err != nil {
		t.Fatalf("delete first generation: %v", err)
	}
	createMetricsTestTable(t, store, "MergeGenTbl")
	second, err := store.Tables().Get("MergeGenTbl")
	if err != nil {
		t.Fatalf("second generation: %v", err)
	}
	if second.TableId == first.TableId {
		t.Fatal("the recreated table kept the first generation's id")
	}

	// The merged flush belongs to the superseded generation and is dropped.
	store.FlushTableMetrics(metricDeltas)
	successor, err := store.Tables().Get("MergeGenTbl")
	if err != nil {
		t.Fatalf("reload successor: %v", err)
	}
	if successor.ItemCount != 0 || successor.TableSizeBytes != 0 {
		t.Fatalf("successor inherited the superseded generation's counters through the merge: %d/%d",
			successor.ItemCount, successor.TableSizeBytes)
	}
}

// TestStoreOpenSweepsRestartOrphanedJobs pins the restart reconciliation:
// an export or import record left IN_PROGRESS by a dead process cannot
// complete — its carrier goroutine died with that process — so a store
// opening on the persisted state fails the record before any request can
// observe it, while terminal records pass through untouched.
func TestStoreOpenSweepsRestartOrphanedJobs(t *testing.T) {
	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	first := NewDynamoDBStore(st, st, "123456789012", "us-east-1")

	orphanExport, err := first.Exports().Create("arn:aws:dynamodb:us-east-1:123456789012:table/SweepTbl", "sweep-tid", "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("create orphaned export: %v", err)
	}
	orphanImport, err := first.Imports().Create("arn:aws:dynamodb:us-east-1:123456789012:table/SweepTbl", "SweepTbl")
	if err != nil {
		t.Fatalf("create orphaned import: %v", err)
	}
	doneExport, err := first.Exports().Create("arn:aws:dynamodb:us-east-1:123456789012:table/SweepTbl", "sweep-tid2", "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("create completed export: %v", err)
	}
	doneExport.ExportStatus = "COMPLETED"
	if err := first.Exports().Put(doneExport); err != nil {
		t.Fatalf("complete export: %v", err)
	}
	first.Close()
	st.Close()

	// The restart: a fresh store opens on the same persisted state.
	st2, err := storage.Open(dir)
	if err != nil {
		t.Fatalf("reopen storage: %v", err)
	}
	defer st2.Close()
	second := NewDynamoDBStore(st2, st2, "123456789012", "us-east-1")
	defer second.Close()

	finalExport, err := second.Exports().Get(orphanExport.ExportArn)
	if err != nil {
		t.Fatalf("reload orphaned export: %v", err)
	}
	if finalExport.ExportStatus != "FAILED" || finalExport.FailureCode != restartOrphanedJobCode {
		t.Fatalf("orphaned export = %s/%s, want FAILED/%s", finalExport.ExportStatus, finalExport.FailureCode, restartOrphanedJobCode)
	}
	finalImport, err := second.Imports().Get(orphanImport.ImportArn)
	if err != nil {
		t.Fatalf("reload orphaned import: %v", err)
	}
	if finalImport.ImportStatus != "FAILED" || finalImport.FailureCode != restartOrphanedJobCode {
		t.Fatalf("orphaned import = %s/%s, want FAILED/%s", finalImport.ImportStatus, finalImport.FailureCode, restartOrphanedJobCode)
	}
	survivor, err := second.Exports().Get(doneExport.ExportArn)
	if err != nil {
		t.Fatalf("reload completed export: %v", err)
	}
	if survivor.ExportStatus != "COMPLETED" {
		t.Fatalf("terminal record touched by the sweep: %s", survivor.ExportStatus)
	}
}

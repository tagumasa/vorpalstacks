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

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
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

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
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

	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
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

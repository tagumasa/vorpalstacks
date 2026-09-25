package dynamodb

import (
	"fmt"

	"vorpalstacks/internal/core/storage"
)

// ---------------------------------------------------------------------------
// Transaction index plane — the DynamoDBTxn methods that maintain and query
// index entries: the GSI/LSI/vector entry put and delete paths every item
// write walks, the single-index variants backfill and index deletion use,
// and the GSI/LSI query delegation with its shared option record. The
// GSI/LSI key construction and bucket operations live in index_store.go,
// the vector backend in vector_store.go.
// ---------------------------------------------------------------------------

// updateIndexEntries runs the shared shape of the item index-entry
// maintenance: load the table once, apply the vector step when the table
// carries vector indexes, then the GSI/LSI step — the put and delete paths
// differ only in the direction of the two steps.
func (t *DynamoDBTxn) updateIndexEntries(tableName, op string, item *Item, vectorStep func(storage.Transaction, string, *Table, *Item) error, indexStep func(*IndexStore, storage.Transaction, *Table, *Item) error) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for %s: %w", tableName, op, err)
	}
	if len(table.VectorIndexes) > 0 {
		if err := vectorStep(t.txn, t.region(), table, item); err != nil {
			return err
		}
	}
	return indexStep(t.indexStore, t.txn, table, item)
}

// PutIndexEntries stores index entries for an item in the transaction.
// Delegates to IndexStore for the GSI/LSI key construction and bucket
// operations, and maintains the vector index entries the same way so every
// item write path keeps all index families current.
func (t *DynamoDBTxn) PutIndexEntries(tableName string, item *Item) error {
	return t.updateIndexEntries(tableName, "PutIndexEntries", item, putVectorEntries, (*IndexStore).PutIndexEntries)
}

// DeleteIndexEntries removes index entries for an item from the
// transaction. Delegates to IndexStore for GSI/LSI entries and removes the
// vector index entries alongside them.
func (t *DynamoDBTxn) DeleteIndexEntries(tableName string, item *Item) error {
	return t.updateIndexEntries(tableName, "DeleteIndexEntries", item, deleteVectorEntries, (*IndexStore).DeleteIndexEntries)
}

// PutVectorEntriesForIndex stores only the named vector index's entry for an
// item in the transaction, mirroring PutGSIEntriesForIndex for vector index
// backfill.
func (t *DynamoDBTxn) PutVectorEntriesForIndex(tableName, indexName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for PutVectorEntriesForIndex: %w", tableName, err)
	}
	return putVectorEntriesForIndex(t.txn, t.region(), table, indexName, item)
}

// DeleteVectorEntriesForIndex removes every entry of the named vector index
// so a deleted index's data does not outlive the index, mirroring
// DeleteIndexEntriesForIndex.
func (t *DynamoDBTxn) DeleteVectorEntriesForIndex(tableName, indexName string) error {
	return deleteVectorEntriesForIndex(t.txn, t.region(), tableName, indexName)
}

// VectorTopK performs a brute-force similarity search over the named vector
// index within the transaction.
func (t *DynamoDBTxn) VectorTopK(tableName, indexName string, query []float64, k int, filter func(*Item) bool) ([]VectorSearchHit, error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s for VectorTopK: %w", tableName, err)
	}
	return vectorTopK(t.txn, t.region(), table, indexName, query, k, filter)
}

// PutGSIEntriesForIndex stores only the named GSI's index entry for an item
// in the transaction. Used by GSI backfill so a newly added index is
// populated without rewriting every other index's entries.
func (t *DynamoDBTxn) PutGSIEntriesForIndex(tableName, indexName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for PutGSIEntriesForIndex: %w", tableName, err)
	}
	return t.indexStore.PutIndexEntriesForIndex(t.txn, table, indexName, item)
}

// DeleteIndexEntriesForIndex removes every index entry of the named GSI in
// the transaction. Called when UpdateTable deletes the index so its entries
// cannot outlive it. Delegates to IndexStore.
func (t *DynamoDBTxn) DeleteIndexEntriesForIndex(tableName, indexName string) error {
	return t.indexStore.DeleteIndexEntriesForIndex(t.txn, tableName, indexName)
}

// QueryByGSI queries a global secondary index for items matching the
// hash key. Delegates to IndexStore.
func (t *DynamoDBTxn) QueryByGSI(tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error) {
	_, err := t.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s for GSI query: %w", tableName, err)
	}
	return t.indexStore.QueryGSI(t.txn, tableName, indexName, hashKeyValue, opts)
}

// QueryByLSI queries a local secondary index for items matching the
// hash key. Delegates to IndexStore.
func (t *DynamoDBTxn) QueryByLSI(tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error) {
	_, err := t.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s for LSI query: %w", tableName, err)
	}
	return t.indexStore.QueryLSI(t.txn, tableName, indexName, hashKeyValue, opts)
}

// IndexQueryOptions defines options for querying indexes. Filter, when
// non-nil, is applied to each resolved item during iteration — before the
// Limit is counted — so a filtered query reads only as far as its page
// requires. Marker resumes the walk strictly after (forward) or strictly
// before (Reverse) the stored index key it names.
type IndexQueryOptions struct {
	Limit   int
	Reverse bool
	Marker  string
	Filter  func(*Item) bool
}

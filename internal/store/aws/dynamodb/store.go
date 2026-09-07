// Package dynamodb provides DynamoDB data store implementations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

const KeySep = "\x00"

// Default throughput quotas per the DynamoDB quotas documentation, per
// Region: a table may provision at most 40,000 read and 40,000 write
// capacity units, and the account-wide total across all provisioned tables
// and GSIs is 80,000 read and 80,000 write capacity units.
const (
	AccountMaxReadCapacityUnits  = 80000
	AccountMaxWriteCapacityUnits = 80000
	TableMaxReadCapacityUnits    = 40000
	TableMaxWriteCapacityUnits   = 40000
)

// DynamoDBStore provides a unified interface to all DynamoDB store components.
type DynamoDBStore struct {
	tables       *TableStore
	items        *ItemStore
	indexes      *IndexStore
	backups      *BackupStore
	globalTables *GlobalTableStore
	exports      *ExportStore
	imports      *ImportStore
	streams      *StreamStore
	idempotency  *IdempotencyStore
	journal      *JournalStore
	contributors *ContributorStore
	storage      storage.TransactionalStorageWith2PC
	ttlWorker    *ttlWorker
	// contributorMu serialises contributor counter updates across the
	// store: the storage layer offers no cross-transaction isolation, so a
	// plain read-modify-write can lose increments to a concurrent update.
	contributorMu sync.Mutex
}

// NewDynamoDBStore creates a new DynamoDB store with the specified storage, account ID, and region.
func NewDynamoDBStore(store storage.TransactionalStorageWith2PC, accountID, region string) *DynamoDBStore {
	tableStore := NewTableStore(store, accountID, region)
	itemStore := NewItemStore(store, tableStore)
	indexStore := NewIndexStore(region)
	backupStore := NewBackupStore(store, accountID, region)
	globalTableStore := NewGlobalTableStore(store, accountID, region)
	exportStore := NewExportStore(store, accountID, region)
	importStore := NewImportStore(store, accountID, region)
	streamStore := NewStreamStore(store, accountID, region)
	idempotencyStore := NewIdempotencyStore(store, region)
	journalStore := NewJournalStore(store, region)
	contributorStore := NewContributorStore(store, region)

	s := &DynamoDBStore{
		tables:       tableStore,
		items:        itemStore,
		indexes:      indexStore,
		backups:      backupStore,
		globalTables: globalTableStore,
		exports:      exportStore,
		imports:      importStore,
		streams:      streamStore,
		idempotency:  idempotencyStore,
		journal:      journalStore,
		contributors: contributorStore,
		storage:      store,
	}
	s.ttlWorker = newTTLWorker(s)
	return s
}

// Close stops the TTL cleanup worker and releases resources.
func (s *DynamoDBStore) Close() {
	if s.ttlWorker != nil {
		s.ttlWorker.Close()
	}
}

// Tables returns the table store for managing DynamoDB table metadata.
func (s *DynamoDBStore) Tables() TableStoreInterface {
	return s.tables
}

// Items returns the item store for managing DynamoDB items.
func (s *DynamoDBStore) Items() ItemStoreInterface {
	return s.items
}

// Indexes returns the index store for managing GSI and LSI index entries.
func (s *DynamoDBStore) Indexes() *IndexStore {
	return s.indexes
}

// Idempotency returns the store for client request token idempotency.
func (s *DynamoDBStore) Idempotency() *IdempotencyStore {
	return s.idempotency
}

// Journal returns the item-mutation journal store backing point-in-time
// recovery.
func (s *DynamoDBStore) Journal() *JournalStore {
	return s.journal
}

// Backups returns the backup store for managing DynamoDB backups.
func (s *DynamoDBStore) Backups() BackupStoreInterface {
	return s.backups
}

// GlobalTables returns the global table store for managing DynamoDB global tables.
func (s *DynamoDBStore) GlobalTables() GlobalTableStoreInterface {
	return s.globalTables
}

// Exports returns the export store for managing DynamoDB exports to S3.
func (s *DynamoDBStore) Exports() ExportStoreInterface {
	return s.exports
}

// Imports returns the import store for managing DynamoDB imports from S3.
func (s *DynamoDBStore) Imports() ImportStoreInterface {
	return s.imports
}

// Streams returns the stream record store for managing DynamoDB Streams.
func (s *DynamoDBStore) Streams() *StreamStore {
	return s.streams
}

// Contributors returns the contributor access aggregation store.
func (s *DynamoDBStore) Contributors() *ContributorStore {
	return s.contributors
}

// Storage returns the underlying storage for this DynamoDB store.
func (s *DynamoDBStore) Storage() storage.TransactionalStorageWith2PC {
	return s.storage
}

// View executes a read-only transaction on the DynamoDB store.
func (s *DynamoDBStore) View(ctx context.Context, fn func(txn *DynamoDBTxn) error) error {
	return s.storage.View(ctx, func(txn storage.Transaction) error {
		return fn(&DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes})
	})
}

// Update executes a read-write transaction on the DynamoDB store. When the
// transaction commits, the contributor events queued by its item writes are
// applied to the access counters, and the table metric deltas queued by its
// counter calls are applied to the table records.
func (s *DynamoDBStore) Update(ctx context.Context, fn func(txn *DynamoDBTxn) error) error {
	var dtxn *DynamoDBTxn
	err := s.storage.Update(ctx, func(txn storage.Transaction) error {
		dtxn = &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return fn(dtxn)
	})
	if err != nil {
		return err
	}
	s.FlushContributorWrites(ctx, dtxn.TakeContributorWrites())
	s.FlushTableMetrics(dtxn.TakeTableMetricDeltas())
	return nil
}

// FlushTableMetrics applies queued per-table metric deltas to the table
// records. Tables are visited in name order and each record is updated
// under its key lock with an immediate read-modify-write, so concurrent
// transactions and the locked TableStore methods cannot lose increments.
// A flush failure is logged and never fails the observed operation — the
// counters are documented approximate values that AWS itself refreshes
// only periodically.
//
// Callers must not hold a table's record lock while queueing metric deltas
// for that table: the flush takes the same lock. No path does — deltas are
// queued exclusively from item-write transactions, which hold no table
// record locks.
func (s *DynamoDBStore) FlushTableMetrics(deltas map[string]TableMetricDelta) {
	if len(deltas) == 0 {
		return
	}
	names := make([]string, 0, len(deltas))
	for name := range deltas {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := s.tables.applyMetricDeltas(name, deltas[name].ItemCount, deltas[name].SizeBytes); err != nil {
			logs.Warn("failed to apply table metric deltas",
				logs.String("table", name), logs.Err(err))
		}
	}
}

// FlushContributorWrites applies queued contributor events to the access
// counters in one transaction serialised against every other counter
// update, so concurrent reads and writes cannot lose increments. The
// carrying item transaction has already committed: a flush failure is
// logged and never fails the observed operation.
func (s *DynamoDBStore) FlushContributorWrites(ctx context.Context, events []ContributorWriteEvent) {
	if len(events) == 0 {
		return
	}
	s.contributorMu.Lock()
	defer s.contributorMu.Unlock()
	err := s.storage.Update(ctx, func(txn storage.Transaction) error {
		dtxn := &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return dtxn.applyContributorWrites(events)
	})
	if err != nil {
		logs.Warn("failed to record contributor writes", logs.Err(err))
	}
}

// TwoPhaseTransaction returns a two-phase transaction interface for the DynamoDB store.
func (s *DynamoDBStore) TwoPhaseTransaction() storage.TwoPhaseTransaction {
	return s.storage.TwoPhaseTransaction()
}

// NewTxn creates a new DynamoDB transaction wrapper over the given storage transaction.
func (s *DynamoDBStore) NewTxn(txn storage.Transaction) *DynamoDBTxn {
	return &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
}

// DynamoDBTxn represents a DynamoDB transaction for atomic operations.
type DynamoDBTxn struct {
	txn        storage.Transaction
	tableStore *TableStore
	indexStore *IndexStore
	// contributorWrites collects the item writes observed in this
	// transaction; the store applies them to the access counters after the
	// transaction commits.
	contributorWrites []ContributorWriteEvent
	// tableMetricDeltas accumulates the item-count and size changes queued
	// by this transaction's counter calls, per table; the store applies
	// them to the table records after the transaction commits.
	tableMetricDeltas map[string]TableMetricDelta
}

// TableMetricDelta accumulates the item-count and size-bytes changes one
// transaction observes for a table. The counters are fields of the whole
// table record, so the deltas are aggregated during the transaction and
// applied once, under the table's record lock, after it commits.
type TableMetricDelta struct {
	ItemCount int64
	SizeBytes int64
}

func (t *DynamoDBTxn) region() string {
	if t.tableStore != nil {
		return t.tableStore.region
	}
	return ""
}

// RawTxn returns the underlying storage transaction so callers can write
// to additional buckets (e.g. stream records) within the same atomic
// transaction as item mutations.
func (t *DynamoDBTxn) RawTxn() storage.Transaction {
	return t.txn
}

// GetTable retrieves a table by name within the transaction.
func (t *DynamoDBTxn) GetTable(name string) (*Table, error) {
	bucket := t.txn.Bucket(tableBucketName(t.region()))
	data, err := bucket.Get([]byte(name))
	if err != nil {
		return nil, fmt.Errorf("get table %s: %w", name, err)
	}
	if data == nil {
		return nil, ErrTableNotFound
	}
	var pbTable pb.Table
	if err := proto.Unmarshal(data, &pbTable); err != nil {
		return nil, fmt.Errorf("unmarshal table %s: %w", name, err)
	}
	return ProtoToTable(&pbTable), nil
}

// GetItem retrieves an item from a table by its key.
func (t *DynamoDBTxn) GetItem(tableName string, key map[string]*AttributeValue) (*Item, error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s for GetItem: %w", tableName, err)
	}
	itemKey := EncodeItemKey(tableName, key, table)
	if itemKey == "" {
		return nil, ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	data, err := bucket.Get([]byte(itemKey))
	if err != nil {
		return nil, fmt.Errorf("get item %s#%v: %w", tableName, key, err)
	}
	if data == nil {
		return nil, ErrItemNotFound
	}
	var pbItem pb.Item
	if err := proto.Unmarshal(data, &pbItem); err != nil {
		return nil, fmt.Errorf("unmarshal item %s#%v: %w", tableName, key, err)
	}
	return itemFromProto(&pbItem), nil
}

// PutItem stores an item in a table.
func (t *DynamoDBTxn) PutItem(tableName string, key map[string]*AttributeValue, attributes map[string]*AttributeValue) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for PutItem: %w", tableName, err)
	}

	mergedAttrs := make(map[string]*AttributeValue, len(attributes)+len(key))
	for k, v := range attributes {
		mergedAttrs[k] = v
	}
	for k, v := range key {
		mergedAttrs[k] = v
	}

	pbItem := &pb.Item{
		TableName:  tableName,
		Key:        attributeValueMapToProtoDirect(key),
		Attributes: attributeValueMapToProtoDirect(mergedAttrs),
	}
	itemKey := EncodeItemKey(tableName, key, table)
	if itemKey == "" {
		return ErrInvalidKey
	}
	data, err := proto.Marshal(pbItem)
	if err != nil {
		return fmt.Errorf("marshal item %s: %w", tableName, err)
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	if err := bucket.Put([]byte(itemKey), data); err != nil {
		return err
	}
	if err := t.journalPut(table, key); err != nil {
		return err
	}
	t.queueContributorWrite(table, key)
	return nil
}

// queueContributorWrite defers the contributor access accounting of an
// item write until the carrying transaction commits. The counters update
// in their own serialised transaction afterwards, which is the only way to
// keep the read-modify-write safe against concurrent updates on a storage
// layer without cross-transaction isolation. Tables without contributor
// insights queue nothing, so their writes open no extra transaction.
func (t *DynamoDBTxn) queueContributorWrite(table *Table, key map[string]*AttributeValue) {
	if table == nil || !table.ContributorInsightsEnabled {
		return
	}
	event := ContributorWriteEvent{TableName: table.Name, Key: make(map[string]*AttributeValue, len(key))}
	for name, value := range key {
		event.Key[name] = value
	}
	t.contributorWrites = append(t.contributorWrites, event)
}

// TakeContributorWrites drains the contributor events queued by the item
// writes of this transaction. The store applies them after the carrying
// transaction commits.
func (t *DynamoDBTxn) TakeContributorWrites() []ContributorWriteEvent {
	events := t.contributorWrites
	t.contributorWrites = nil
	return events
}

// applyContributorWrites credits each queued write event under every
// contributor layout of its table. A write counts as three units of
// ConsumedThroughputUnits. Events on the same counter key are aggregated
// first: the storage layer's read-modify-write reads the committed state,
// so unaggregated repeats inside one transaction would overwrite each
// other.
func (t *DynamoDBTxn) applyContributorWrites(events []ContributorWriteEvent) error {
	tables := make(map[string]*Table)
	type counterTarget struct {
		table  *Table
		layout string
		keyStr string
	}
	aggregated := make(map[counterTarget]int64)
	var order []counterTarget
	for _, event := range events {
		table := tables[event.TableName]
		if _, resolved := tables[event.TableName]; !resolved {
			fetched, err := t.GetTable(event.TableName)
			if err != nil || fetched == nil || !fetched.ContributorInsightsEnabled {
				fetched = nil
			}
			tables[event.TableName] = fetched
			table = fetched
		}
		if table == nil {
			continue
		}
		for _, layout := range ContributorLayouts(table) {
			target := counterTarget{table: table, layout: layout, keyStr: ContributorKeyString(table, event.Key, layout)}
			if _, seen := aggregated[target]; !seen {
				order = append(order, target)
			}
			aggregated[target]++
		}
	}
	at := time.Now()
	for _, target := range order {
		if err := RecordAccessTxn(t.txn, t.region(), target.table.Name, target.layout, target.keyStr, at, aggregated[target], ContributorWriteUnits); err != nil {
			return err
		}
	}
	return nil
}

// contributorRecordReads credits one read event for every key under each
// contributor layout of the table, inside the caller's transaction. Reads
// on the same counter key (items sharing a partition key) are aggregated
// first: the storage layer's read-modify-write reads the committed state,
// so unaggregated repeats inside one transaction would overwrite each
// other.
func (t *DynamoDBTxn) contributorRecordReads(tableName string, keys []map[string]*AttributeValue) error {
	if len(keys) == 0 {
		return nil
	}
	table, err := t.GetTable(tableName)
	if err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	type counterTarget struct {
		layout string
		keyStr string
	}
	aggregated := make(map[counterTarget]int64)
	var order []counterTarget
	for _, key := range keys {
		for _, layout := range ContributorLayouts(table) {
			target := counterTarget{layout: layout, keyStr: ContributorKeyString(table, key, layout)}
			if _, seen := aggregated[target]; !seen {
				order = append(order, target)
			}
			aggregated[target]++
		}
	}
	at := time.Now()
	for _, target := range order {
		if err := RecordAccessTxn(t.txn, t.region(), table.Name, target.layout, target.keyStr, at, aggregated[target], ContributorReadUnits); err != nil {
			return err
		}
	}
	return nil
}

// contributorRecordQueryEvent credits the single event a Query contributes,
// on the partition-key series only: a Query is one read event regardless of
// how many items it returns, and a result set spans many sort keys.
func (t *DynamoDBTxn) contributorRecordQueryEvent(tableName string, key map[string]*AttributeValue) error {
	table, err := t.GetTable(tableName)
	if err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	keyStr := ContributorKeyString(table, key, ContributorLayoutPartitionKey)
	if keyStr == "" {
		return nil
	}
	return RecordAccessTxn(t.txn, t.region(), table.Name, ContributorLayoutPartitionKey, keyStr, time.Now(), 1, ContributorReadUnits)
}

// RecordContributorReads credits one read event per key under every
// contributor layout of the table. The update runs under the contributor
// lock, so concurrent counter updates cannot lose increments.
func (s *DynamoDBStore) RecordContributorReads(ctx context.Context, tableName string, keys []map[string]*AttributeValue) error {
	if len(keys) == 0 {
		return nil
	}
	if table, err := s.Tables().Get(tableName); err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	s.contributorMu.Lock()
	defer s.contributorMu.Unlock()
	return s.storage.Update(ctx, func(txn storage.Transaction) error {
		dtxn := &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return dtxn.contributorRecordReads(tableName, keys)
	})
}

// RecordContributorQuery credits the single read event a Query contributes
// to the partition-key series. The update runs under the contributor lock,
// so concurrent counter updates cannot lose increments.
func (s *DynamoDBStore) RecordContributorQuery(ctx context.Context, tableName string, key map[string]*AttributeValue) error {
	if table, err := s.Tables().Get(tableName); err != nil || table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	s.contributorMu.Lock()
	defer s.contributorMu.Unlock()
	return s.storage.Update(ctx, func(txn storage.Transaction) error {
		dtxn := &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
		return dtxn.contributorRecordQueryEvent(tableName, key)
	})
}

// journalPut appends a journal record for a put on a table with
// point-in-time recovery enabled. The record stores the pre-write item so a
// restore can undo the change; the append shares the caller's transaction,
// so the journal never diverges from the item state.
func (t *DynamoDBTxn) journalPut(table *Table, key map[string]*AttributeValue) error {
	if !pitrEnabled(table) {
		return nil
	}
	beforeImage := t.itemBeforeImage(table.Name, key)
	return appendJournalTxn(t.txn, t.region(), table.Name, JournalOperationPut, key, beforeImage)
}

// journalDelete appends a journal record for a delete. Deletes of keys that
// do not exist are not journaled because they change nothing.
func (t *DynamoDBTxn) journalDelete(table *Table, key map[string]*AttributeValue) error {
	if !pitrEnabled(table) {
		return nil
	}
	beforeImage := t.itemBeforeImage(table.Name, key)
	if beforeImage == nil {
		return nil
	}
	return appendJournalTxn(t.txn, t.region(), table.Name, JournalOperationDelete, key, beforeImage)
}

// pitrEnabled reports whether the table currently has point-in-time
// recovery enabled.
func pitrEnabled(table *Table) bool {
	return table != nil && table.PointInTimeRecovery != nil && table.PointInTimeRecovery.Status == PITRStatusEnabled
}

// itemBeforeImage reads the item's current attribute map through the
// transaction, or nil when the key does not exist.
func (t *DynamoDBTxn) itemBeforeImage(tableName string, key map[string]*AttributeValue) map[string]*AttributeValue {
	existing, err := t.GetItem(tableName, key)
	if err != nil || existing == nil {
		return nil
	}
	return existing.Attributes
}

// DeleteItem removes an item from a table by its key.
func (t *DynamoDBTxn) DeleteItem(tableName string, key map[string]*AttributeValue) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for DeleteItem: %w", tableName, err)
	}
	itemKey := EncodeItemKey(tableName, key, table)
	if itemKey == "" {
		return ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	if err := bucket.Delete([]byte(itemKey)); err != nil {
		return err
	}
	if err := t.journalDelete(table, key); err != nil {
		return err
	}
	t.queueContributorWrite(table, key)
	return nil
}

// ItemExists checks whether an item with the given key exists in the table.
func (t *DynamoDBTxn) ItemExists(tableName string, key map[string]*AttributeValue) (bool, error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return false, fmt.Errorf("get table %s for ItemExists: %w", tableName, err)
	}
	itemKey := EncodeItemKey(tableName, key, table)
	if itemKey == "" {
		return false, ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	return bucket.Has([]byte(itemKey)), nil
}

// UpdateItemCount queues an item-count change for a table. The counter is a
// field of the whole table record and the storage layer is read-committed:
// reading and rewriting the record inside the caller's transaction would
// overwrite same-transaction and concurrent increments. The queued delta is
// applied once, under the table's record lock, after the carrying
// transaction commits. Table existence is guaranteed by the item write that
// precedes every counter call on each path.
func (t *DynamoDBTxn) UpdateItemCount(tableName string, delta int64) error {
	if t.tableMetricDeltas == nil {
		t.tableMetricDeltas = make(map[string]TableMetricDelta)
	}
	d := t.tableMetricDeltas[tableName]
	d.ItemCount += delta
	t.tableMetricDeltas[tableName] = d
	return nil
}

// UpdateTableSize queues a table-size change for a table, aggregated and
// applied like UpdateItemCount after the carrying transaction commits.
func (t *DynamoDBTxn) UpdateTableSize(tableName string, delta int64) error {
	if t.tableMetricDeltas == nil {
		t.tableMetricDeltas = make(map[string]TableMetricDelta)
	}
	d := t.tableMetricDeltas[tableName]
	d.SizeBytes += delta
	t.tableMetricDeltas[tableName] = d
	return nil
}

// TakeTableMetricDeltas drains the metric deltas queued by this
// transaction. The store applies them after the carrying transaction
// commits.
func (t *DynamoDBTxn) TakeTableMetricDeltas() map[string]TableMetricDelta {
	deltas := t.tableMetricDeltas
	t.tableMetricDeltas = nil
	return deltas
}

// MergeTableMetricDeltas folds every delta of src into dst and returns dst.
// Two-phase transactions run one wrapper per operation; merging their
// queued deltas keeps the post-commit flush one application per table.
func MergeTableMetricDeltas(dst, src map[string]TableMetricDelta) map[string]TableMetricDelta {
	if dst == nil {
		dst = make(map[string]TableMetricDelta, len(src))
	}
	for name, d := range src {
		acc := dst[name]
		acc.ItemCount += d.ItemCount
		acc.SizeBytes += d.SizeBytes
		dst[name] = acc
	}
	return dst
}

// PutIndexEntries stores index entries for an item in the transaction.
// Delegates to IndexStore for the GSI/LSI key construction and bucket
// operations, and maintains the vector index entries the same way so every
// item write path keeps all index families current.
func (t *DynamoDBTxn) PutIndexEntries(tableName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for PutIndexEntries: %w", tableName, err)
	}
	if len(table.VectorIndexes) > 0 {
		if err := putVectorEntries(t.txn, t.region(), table, item); err != nil {
			return err
		}
	}
	return t.indexStore.PutIndexEntries(t.txn, table, item)
}

// DeleteIndexEntries removes index entries for an item from the
// transaction. Delegates to IndexStore for GSI/LSI entries and removes the
// vector index entries alongside them.
func (t *DynamoDBTxn) DeleteIndexEntries(tableName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for DeleteIndexEntries: %w", tableName, err)
	}
	if len(table.VectorIndexes) > 0 {
		if err := deleteVectorEntries(t.txn, t.region(), table, item); err != nil {
			return err
		}
	}
	return t.indexStore.DeleteIndexEntries(t.txn, table, item)
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

// Scan scans all items in a table within the transaction.
func (t *DynamoDBTxn) Scan(tableName string, fn func(item *Item) error) error {
	prefix := tableName + KeySep
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		var pbItem pb.Item
		if err := proto.Unmarshal(iter.Value(), &pbItem); err != nil {
			return fmt.Errorf("unmarshal item during scan: %w", err)
		}
		item := itemFromProto(&pbItem)
		if err := fn(item); err != nil {
			return err
		}
	}
	return iter.Error()
}

// ScanByPartitionKey scans items with a specific partition key within the
// transaction. partitionKeyValue must already be an EncodeKeyValue rendering;
// the encoded component is prefix-free, so the scan cannot cross into
// another partition and no trailing separator is needed.
func (t *DynamoDBTxn) ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for ScanByPartitionKey: %w", tableName, err)
	}

	prefix := tableName + KeySep + partitionKeyValue
	pkName, _ := schemaKeyNames(table.KeySchema)

	bucket := t.txn.Bucket(itemBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		var pbItem pb.Item
		if err := proto.Unmarshal(iter.Value(), &pbItem); err != nil {
			return fmt.Errorf("unmarshal item during ScanByPartitionKey: %w", err)
		}
		item := itemFromProto(&pbItem)
		itemPkValue := EncodeKeyValue(item.Key[pkName])
		if itemPkValue != partitionKeyValue {
			continue
		}
		if err := fn(item); err != nil {
			return err
		}
	}
	return iter.Error()
}

// DeleteTableCascade removes all data associated with a table within a single
// transaction: items, GSI/LSI index entries, the table record itself, backups,
// exports, imports, tags, and global table entries.
// It does NOT check DeletionProtectionEnabled — that is the caller's responsibility.
func (t *DynamoDBTxn) DeleteTableCascade(name string) error {
	if err := t.deleteAllByPrefix(itemBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete items for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(gsiIndexBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete GSI index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(lsiIndexBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete LSI index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(vectorIndexBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete vector index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(streamBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete stream records for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(journalBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete journal records for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(contributorBucketName(t.region()), name+KeySep); err != nil {
		return fmt.Errorf("delete contributor counters for table %s: %w", name, err)
	}

	if err := t.deleteBackupsForTable(name); err != nil {
		return fmt.Errorf("delete backups for table %s: %w", name, err)
	}

	if err := t.deleteExportsForTable(name); err != nil {
		return fmt.Errorf("delete exports for table %s: %w", name, err)
	}

	if err := t.deleteImportsForTable(name); err != nil {
		return fmt.Errorf("delete imports for table %s: %w", name, err)
	}

	tableBucket := t.txn.Bucket(tableBucketName(t.region()))
	if err := tableBucket.Delete([]byte(name)); err != nil {
		return fmt.Errorf("delete table record %s: %w", name, err)
	}

	tagMainBucket := t.txn.Bucket(tagMainBucketName(t.region()))
	if tagMainBucket != nil {
		if err := tagMainBucket.Delete([]byte(name)); err != nil {
			return fmt.Errorf("delete tags for table %s: %w", name, err)
		}
	}

	globalTableBucket := t.txn.Bucket(globalTableBucketName(t.region()))
	if globalTableBucket != nil {
		if err := globalTableBucket.Delete([]byte(name)); err != nil {
			return fmt.Errorf("delete global table entry for %s: %w", name, err)
		}
	}

	tagIdxBucket := t.txn.Bucket(tagIndexBucketName(t.region()))
	if tagIdxBucket != nil {
		suffix := "\x00" + name
		iter := tagIdxBucket.ScanPrefix(nil)
		defer iter.Close()
		var idxKeysToDelete []string
		for iter.Next() {
			k := string(iter.Key())
			if strings.HasSuffix(k, suffix) {
				idxKeysToDelete = append(idxKeysToDelete, k)
			}
		}
		if err := iter.Error(); err != nil {
			return fmt.Errorf("scan tag index for table %s: %w", name, err)
		}
		for _, k := range idxKeysToDelete {
			if err := tagIdxBucket.Delete([]byte(k)); err != nil {
				return fmt.Errorf("delete tag index entry %s: %w", k, err)
			}
		}
	}

	return nil
}

func (t *DynamoDBTxn) deleteBackupsForTable(tableName string) error {
	backupBucket := t.txn.Bucket(backupBucketName(t.region()))
	if backupBucket == nil {
		return nil
	}
	iter := backupBucket.ScanPrefix(nil)
	defer iter.Close()
	var keysToDelete []string
	for iter.Next() {
		var backup Backup
		if err := json.Unmarshal(iter.Value(), &backup); err != nil {
			logs.Warn("cascade delete: corrupted backup record, will delete", logs.String("key", string(iter.Key())), logs.Err(err))
			keysToDelete = append(keysToDelete, string(iter.Key()))
			continue
		}
		if backup.SourceTableName == tableName {
			keysToDelete = append(keysToDelete, string(iter.Key()))
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	for _, k := range keysToDelete {
		if err := backupBucket.Delete([]byte(k)); err != nil {
			return err
		}
	}
	return nil
}

func (t *DynamoDBTxn) deleteExportsForTable(tableName string) error {
	exportBucket := t.txn.Bucket(exportBucketName(t.region()))
	if exportBucket == nil {
		return nil
	}
	iter := exportBucket.ScanPrefix(nil)
	defer iter.Close()
	var keysToDelete []string
	for iter.Next() {
		var export pb.ExportDescription
		if err := proto.Unmarshal(iter.Value(), &export); err != nil {
			logs.Warn("cascade delete: corrupted export record, will delete", logs.String("key", string(iter.Key())), logs.Err(err))
			keysToDelete = append(keysToDelete, string(iter.Key()))
			continue
		}
		if strings.HasSuffix(export.TableArn, "/table/"+tableName) {
			keysToDelete = append(keysToDelete, string(iter.Key()))
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	for _, k := range keysToDelete {
		if err := exportBucket.Delete([]byte(k)); err != nil {
			return err
		}
	}
	return nil
}

func (t *DynamoDBTxn) deleteImportsForTable(tableName string) error {
	importBucket := t.txn.Bucket(importBucketName(t.region()))
	if importBucket == nil {
		return nil
	}
	iter := importBucket.ScanPrefix(nil)
	defer iter.Close()
	var keysToDelete []string
	for iter.Next() {
		var imp pb.ImportTableDescription
		if err := proto.Unmarshal(iter.Value(), &imp); err != nil {
			logs.Warn("cascade delete: corrupted import record, will delete", logs.String("key", string(iter.Key())), logs.Err(err))
			keysToDelete = append(keysToDelete, string(iter.Key()))
			continue
		}
		if strings.HasSuffix(imp.TableArn, "/table/"+tableName) {
			keysToDelete = append(keysToDelete, string(iter.Key()))
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	for _, k := range keysToDelete {
		if err := importBucket.Delete([]byte(k)); err != nil {
			return err
		}
	}
	return nil
}

// StoreItemWrite stores attrs as the item under key inside the transaction,
// retiring the replaced item's index entries and adjusting the table's item
// count and size: a new item increments the count and adds its size, a
// replacement keeps the count and adjusts the size by the delta. oldItem is
// the item whose index entries the write retires (nil when the key was
// vacant or the caller retired them itself); existed reports whether this
// write replaces an item and oldSize that item's pre-write size — the metric
// inputs may survive only as a size when the caller did not materialise an
// old-image copy. Every item-mutating path — single items, transactions,
// batches, PartiQL statements, imports — routes its write through here so
// the index swap and the metric deltas cannot drift apart.
func (t *DynamoDBTxn) StoreItemWrite(tableName string, key, attrs map[string]*AttributeValue, oldItem *Item, existed bool, oldSize int64) error {
	if oldItem != nil {
		if err := t.DeleteIndexEntries(tableName, oldItem); err != nil {
			return err
		}
	}
	if err := t.PutItem(tableName, key, attrs); err != nil {
		return err
	}
	newItem := &Item{TableName: tableName, Key: key, Attributes: attrs}
	if err := t.PutIndexEntries(tableName, newItem); err != nil {
		return err
	}
	newSize := CalculateItemSize(attrs)
	if !existed {
		if err := t.UpdateItemCount(tableName, 1); err != nil {
			return err
		}
		return t.UpdateTableSize(tableName, newSize)
	}
	if delta := newSize - oldSize; delta != 0 {
		return t.UpdateTableSize(tableName, delta)
	}
	return nil
}

// DeleteItemWrite removes the item under key inside the transaction,
// deleting its index entries and adjusting the table's item count and size
// by the removed item's. existed reports whether an item was removed and
// oldSize its pre-delete size; a delete of a vacant key still issues the
// item delete (it succeeds silently) but adjusts no metrics. Callers whose
// contract forbids touching a vacant key guard the call themselves.
func (t *DynamoDBTxn) DeleteItemWrite(tableName string, key map[string]*AttributeValue, existing *Item, existed bool, oldSize int64) error {
	if existing != nil {
		if err := t.DeleteIndexEntries(tableName, existing); err != nil {
			return err
		}
	}
	if err := t.DeleteItem(tableName, key); err != nil {
		return err
	}
	if !existed {
		return nil
	}
	if err := t.UpdateItemCount(tableName, -1); err != nil {
		return err
	}
	if oldSize > 0 {
		return t.UpdateTableSize(tableName, -oldSize)
	}
	return nil
}

// prefixDeleteBatchSize bounds how many keys one batch of a prefix delete
// holds before the deletes are issued, keeping a large sweep's pending set
// bounded inside the transaction.
const prefixDeleteBatchSize = 500

// deletePrefixBatched deletes every key under prefix from the bucket,
// collecting keys in bounded batches — the single implementation behind
// table drops and per-index entry sweeps.
func deletePrefixBatched(bucket storage.Bucket, prefix string) error {
	if bucket == nil {
		return nil
	}

	var keysBatch []string

	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		keysBatch = append(keysBatch, string(iter.Key()))
		if len(keysBatch) < prefixDeleteBatchSize {
			continue
		}
		if err := flushPrefixDeletes(bucket, &keysBatch); err != nil {
			return err
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return flushPrefixDeletes(bucket, &keysBatch)
}

// flushPrefixDeletes issues one batch's deletes and resets the batch.
func flushPrefixDeletes(bucket storage.Bucket, keysBatch *[]string) error {
	for _, k := range *keysBatch {
		if err := bucket.Delete([]byte(k)); err != nil {
			return err
		}
	}
	*keysBatch = (*keysBatch)[:0]
	return nil
}

// deleteAllByPrefix deletes all keys with the given prefix from the specified bucket.
func (t *DynamoDBTxn) deleteAllByPrefix(bucketName, prefix string) error {
	return deletePrefixBatched(t.txn.Bucket(bucketName), prefix)
}

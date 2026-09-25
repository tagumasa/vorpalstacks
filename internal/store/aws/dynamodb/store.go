// Package dynamodb provides DynamoDB data store implementations for vorpalstacks.
package dynamodb

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"

	"google.golang.org/protobuf/proto"
)

// This file holds the store root: the DynamoDBStore facade (construction,
// component accessors, View/Update transaction entry points) and the
// DynamoDBTxn wrapper with the read primitive every transactional family
// starts from (GetTable). The transaction planes live in txn_item.go (item
// CRUD, unified write paths, scans), txn_index.go (index and vector entry
// maintenance and queries), txn_accounting.go (contributor and table-metric
// counter queues), and txn_cascade.go (table-deletion sweep and delete-time
// system backup); the journal write hooks are in journal_store.go. The
// package-shared List shape (listProtoConverted) and the bounded
// prefix-delete machinery also stay here.

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

// NewDynamoDBStore creates a new DynamoDB store with the specified storage,
// account ID, and region. The globalStorage argument is the account-global
// storage handle backing the global-table record: global-table membership is
// account metadata every region resolves identically, so the regional store
// and the global-table bucket must not share one physical location.
func NewDynamoDBStore(store storage.TransactionalStorageWith2PC, globalStorage storage.BasicStorage, accountID, region string) *DynamoDBStore {
	tableStore := NewTableStore(store, accountID, region)
	itemStore := NewItemStore(store, tableStore)
	indexStore := NewIndexStore(region)
	backupStore := NewBackupStore(store, accountID, region)
	globalTableStore := NewGlobalTableStore(globalStorage, accountID, region)
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
	// A store opening on persisted state inherits any job a previous
	// process left IN_PROGRESS — its carrier goroutine died with that
	// process, so the sweep fails the record before a request can observe
	// a job that can never complete.
	s.SweepOrphanedJobs()
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

// listProtoConverted runs the List shape every resource family of this
// store shares: one proto page walk over the family's records, each
// accepted record converted to its store type, and the continuation
// marker passed through when the walk is truncated.
func listProtoConverted[P proto.Message, T any](s *common.BaseStore, marker string, limit int, newProto func() P, convert func(P) T, filter func(P) bool) ([]T, string, error) {
	opts := common.ListOptions{
		Marker:   marker,
		MaxItems: limit,
	}
	result, err := common.ListProto[P](s, opts, newProto, filter)
	if err != nil {
		return nil, "", err
	}
	items := make([]T, len(result.Items))
	for i, pbItem := range result.Items {
		items[i] = convert(pbItem)
	}
	if !result.IsTruncated {
		return items, "", nil
	}
	return items, result.NextMarker, nil
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

// Package dynamodb provides DynamoDB data store implementations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

const KeySep = "\x00"

const keySep = KeySep

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
	storage      storage.TransactionalStorageWith2PC
	ttlWorker    *ttlWorker
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

// Update executes a read-write transaction on the DynamoDB store.
func (s *DynamoDBStore) Update(ctx context.Context, fn func(txn *DynamoDBTxn) error) error {
	return s.storage.Update(ctx, func(txn storage.Transaction) error {
		return fn(&DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes})
	})
}

// TwoPhaseTransaction returns a two-phase transaction interface for the DynamoDB store.
func (s *DynamoDBStore) TwoPhaseTransaction() storage.TwoPhaseTransaction {
	return s.storage.TwoPhaseTransaction()
}

// NewTxn creates a new DynamoDB transaction wrapper over the given storage transaction.
func (s *DynamoDBStore) NewTxn(txn storage.Transaction) *DynamoDBTxn {
	return &DynamoDBTxn{txn: txn, tableStore: s.tables, indexStore: s.indexes}
}

// NewDynamoDBTxn creates a new DynamoDB transaction with the given storage transaction and table store.
func NewDynamoDBTxn(txn storage.Transaction, tableStore *TableStore) *DynamoDBTxn {
	return &DynamoDBTxn{txn: txn, tableStore: tableStore, indexStore: NewIndexStore(tableStore.region)}
}

// DynamoDBTxn represents a DynamoDB transaction for atomic operations.
type DynamoDBTxn struct {
	txn        storage.Transaction
	tableStore *TableStore
	indexStore *IndexStore
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

// PutTable stores a table in the transaction.
func (t *DynamoDBTxn) PutTable(table *Table) error {
	bucket := t.txn.Bucket(tableBucketName(t.region()))
	data, err := proto.Marshal(TableToProto(table))
	if err != nil {
		return fmt.Errorf("marshal table %s: %w", table.Name, err)
	}
	return bucket.Put([]byte(table.Name), data)
}

// GetItem retrieves an item from a table by its key.
func (t *DynamoDBTxn) GetItem(tableName string, key map[string]*AttributeValue) (*Item, error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return nil, fmt.Errorf("get table %s for GetItem: %w", tableName, err)
	}
	itemKey := buildItemKeyFromTable(tableName, key, table)
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
	return &Item{
		TableName:  pbItem.TableName,
		Key:        protoToAttributeValueMapDirect(pbItem.Key),
		Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
	}, nil
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
	itemKey := buildItemKeyFromTable(tableName, key, table)
	if itemKey == "" {
		return ErrInvalidKey
	}
	data, err := proto.Marshal(pbItem)
	if err != nil {
		return fmt.Errorf("marshal item %s: %w", tableName, err)
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	return bucket.Put([]byte(itemKey), data)
}

// DeleteItem removes an item from a table by its key.
func (t *DynamoDBTxn) DeleteItem(tableName string, key map[string]*AttributeValue) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for DeleteItem: %w", tableName, err)
	}
	itemKey := buildItemKeyFromTable(tableName, key, table)
	if itemKey == "" {
		return ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	return bucket.Delete([]byte(itemKey))
}

// ItemExists checks whether an item with the given key exists in the table.
func (t *DynamoDBTxn) ItemExists(tableName string, key map[string]*AttributeValue) (bool, error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return false, fmt.Errorf("get table %s for ItemExists: %w", tableName, err)
	}
	itemKey := buildItemKeyFromTable(tableName, key, table)
	if itemKey == "" {
		return false, ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	return bucket.Has([]byte(itemKey)), nil
}

// UpdateItemCount increments or decrements the item count for a table within the current transaction.
func (t *DynamoDBTxn) UpdateItemCount(tableName string, delta int64) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for UpdateItemCount: %w", tableName, err)
	}
	table.ItemCount += delta
	if table.ItemCount < 0 {
		table.ItemCount = 0
	}
	return t.PutTable(table)
}

// UpdateTableSize increments or decrements the table size in bytes within the current transaction.
func (t *DynamoDBTxn) UpdateTableSize(tableName string, delta int64) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for UpdateTableSize: %w", tableName, err)
	}
	table.TableSizeBytes += delta
	if table.TableSizeBytes < 0 {
		table.TableSizeBytes = 0
	}
	return t.PutTable(table)
}

func buildItemKeyFromTable(tableName string, key map[string]*AttributeValue, table *Table) string {
	pkName := ""
	skName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeHash {
			pkName = ks.AttributeName
		} else if ks.KeyType == KeyTypeRange {
			skName = ks.AttributeName
		}
	}

	pkValue := attributeValueToString(key[pkName])
	if pkValue == "" {
		return ""
	}

	if skName != "" {
		if key[skName] == nil {
			return ""
		}
		skValue := attributeValueToString(key[skName])
		if skValue == "" {
			return ""
		}
		return tableName + keySep + pkValue + keySep + skValue
	}

	return tableName + keySep + pkValue
}

func attributeValueToString(av *AttributeValue) string {
	if av == nil {
		return ""
	}
	if av.S != nil {
		return *av.S
	}
	if av.N != nil {
		return formatNumberForSort(*av.N)
	}
	if av.B != nil {
		return string(av.B)
	}
	if av.BOOL != nil {
		return ""
	}
	if av.NULL != nil && *av.NULL {
		return ""
	}
	return ""
}

func formatNumberForSort(numStr string) string {
	if numStr == "" {
		return "1" + strings.Repeat("0", 80)
	}
	rat := new(big.Rat)
	if _, ok := rat.SetString(numStr); !ok {
		return "1" + numStr
	}
	sign := rat.Sign()
	if sign == 0 {
		return "1" + strings.Repeat("0", 80)
	}

	absRat := new(big.Rat).Abs(rat)
	numerator := absRat.Num()
	denominator := absRat.Denom()
	floatVal := new(big.Float).SetRat(new(big.Rat).SetFrac(numerator, denominator))
	floatStr := floatVal.Text('f', 38)
	floatStr = strings.TrimRight(floatStr, "0")
	if strings.HasSuffix(floatStr, ".") {
		floatStr = floatStr[:len(floatStr)-1]
	}

	intPart, fracPart := floatStr, ""
	if dotIdx := strings.Index(floatStr, "."); dotIdx >= 0 {
		intPart = floatStr[:dotIdx]
		fracPart = floatStr[dotIdx+1:]
	}

	intPadded := intPart
	if len(intPadded) < 40 {
		intPadded = strings.Repeat("0", 40-len(intPadded)) + intPadded
	}
	fracPadded := fracPart
	if len(fracPadded) < 40 {
		fracPadded = fracPadded + strings.Repeat("0", 40-len(fracPadded))
	} else if len(fracPadded) > 40 {
		fracPadded = fracPadded[:40]
	}
	digits := intPadded + fracPadded

	if sign > 0 {
		return "1" + digits
	}

	complemented := make([]byte, len(digits))
	for i := 0; i < len(digits); i++ {
		complemented[i] = '9' - digits[i] + '0'
	}
	return "0" + string(complemented)
}

// PutIndexEntries stores index entries for an item in the transaction.
// Delegates to IndexStore for the actual key construction and bucket
// operations.
func (t *DynamoDBTxn) PutIndexEntries(tableName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for PutIndexEntries: %w", tableName, err)
	}
	return t.indexStore.PutIndexEntries(t.txn, table, item)
}

// DeleteIndexEntries removes index entries for an item from the
// transaction. Delegates to IndexStore.
func (t *DynamoDBTxn) DeleteIndexEntries(tableName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for DeleteIndexEntries: %w", tableName, err)
	}
	return t.indexStore.DeleteIndexEntries(t.txn, table, item)
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

// IndexQueryOptions defines options for querying indexes.
type IndexQueryOptions struct {
	Limit   int
	Reverse bool
}

// Scan scans all items in a table within the transaction.
func (t *DynamoDBTxn) Scan(tableName string, fn func(item *Item) error) error {
	prefix := tableName + keySep
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		var pbItem pb.Item
		if err := proto.Unmarshal(iter.Value(), &pbItem); err != nil {
			return fmt.Errorf("unmarshal item during scan: %w", err)
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		if err := fn(item); err != nil {
			return err
		}
	}
	return iter.Error()
}

// ScanByPartitionKey scans items with a specific partition key within the transaction.
func (t *DynamoDBTxn) ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return fmt.Errorf("get table %s for ScanByPartitionKey: %w", tableName, err)
	}

	prefix := tableName + keySep + partitionKeyValue
	pkName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeRange {
			prefix += keySep
			break
		}
		if ks.KeyType == KeyTypeHash {
			pkName = ks.AttributeName
		}
	}

	bucket := t.txn.Bucket(itemBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		var pbItem pb.Item
		if err := proto.Unmarshal(iter.Value(), &pbItem); err != nil {
			return fmt.Errorf("unmarshal item during ScanByPartitionKey: %w", err)
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		itemPkValue := attributeValueToString(item.Key[pkName])
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
	if err := t.deleteAllByPrefix(itemBucketName(t.region()), name+keySep); err != nil {
		return fmt.Errorf("delete items for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(gsiIndexBucketName(t.region()), name+keySep); err != nil {
		return fmt.Errorf("delete GSI index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(lsiIndexBucketName(t.region()), name+keySep); err != nil {
		return fmt.Errorf("delete LSI index entries for table %s: %w", name, err)
	}

	if err := t.deleteAllByPrefix(streamBucketName(t.region()), name+keySep); err != nil {
		return fmt.Errorf("delete stream records for table %s: %w", name, err)
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

// deleteAllByPrefix deletes all keys with the given prefix from the specified bucket.
func (t *DynamoDBTxn) deleteAllByPrefix(bucketName, prefix string) error {
	bucket := t.txn.Bucket(bucketName)
	if bucket == nil {
		return nil
	}

	const batchSize = 500
	var keysBatch []string

	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		keysBatch = append(keysBatch, string(iter.Key()))
		if len(keysBatch) >= batchSize {
			for _, k := range keysBatch {
				if err := bucket.Delete([]byte(k)); err != nil {
					return err
				}
			}
			keysBatch = keysBatch[:0]
		}
	}

	if err := iter.Error(); err != nil {
		return err
	}

	for _, k := range keysBatch {
		if err := bucket.Delete([]byte(k)); err != nil {
			return err
		}
	}

	return nil
}

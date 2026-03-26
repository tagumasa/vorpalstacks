// Package dynamodb provides DynamoDB data store implementations for vorpalstacks.
package dynamodb

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// DynamoDBStore provides a unified interface to all DynamoDB store components.
type DynamoDBStore struct {
	tables       *TableStore
	items        *ItemStore
	backups      *BackupStore
	globalTables *GlobalTableStore
	exports      *ExportStore
	imports      *ImportStore
	storage      storage.TransactionalStorageWith2PC
}

// NewDynamoDBStore creates a new DynamoDB store with the specified storage, account ID, and region.
func NewDynamoDBStore(store storage.TransactionalStorageWith2PC, accountID, region string) *DynamoDBStore {
	tableStore := NewTableStore(store, accountID, region)
	itemStore := NewItemStore(store, tableStore)
	backupStore := NewBackupStore(store, accountID, region)
	globalTableStore := NewGlobalTableStore(store, accountID, region)
	exportStore := NewExportStore(store, accountID, region)
	importStore := NewImportStore(store, accountID, region)

	return &DynamoDBStore{
		tables:       tableStore,
		items:        itemStore,
		backups:      backupStore,
		globalTables: globalTableStore,
		exports:      exportStore,
		imports:      importStore,
		storage:      store,
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

// Storage returns the underlying storage for this DynamoDB store.
func (s *DynamoDBStore) Storage() storage.TransactionalStorageWith2PC {
	return s.storage
}

// View executes a read-only transaction on the DynamoDB store.
func (s *DynamoDBStore) View(ctx context.Context, fn func(txn *DynamoDBTxn) error) error {
	return s.storage.View(ctx, func(txn storage.Transaction) error {
		return fn(&DynamoDBTxn{txn: txn, tableStore: s.tables})
	})
}

// Update executes a read-write transaction on the DynamoDB store.
func (s *DynamoDBStore) Update(ctx context.Context, fn func(txn *DynamoDBTxn) error) error {
	return s.storage.Update(ctx, func(txn storage.Transaction) error {
		return fn(&DynamoDBTxn{txn: txn, tableStore: s.tables})
	})
}

// TwoPhaseTransaction returns a two-phase transaction interface for the DynamoDB store.
func (s *DynamoDBStore) TwoPhaseTransaction() storage.TwoPhaseTransaction {
	return s.storage.TwoPhaseTransaction()
}

// NewDynamoDBTxn creates a new DynamoDB transaction with the given storage transaction and table store.
func NewDynamoDBTxn(txn storage.Transaction, tableStore *TableStore) *DynamoDBTxn {
	return &DynamoDBTxn{txn: txn, tableStore: tableStore}
}

// DynamoDBTxn represents a DynamoDB transaction for atomic operations.
type DynamoDBTxn struct {
	txn        storage.Transaction
	tableStore *TableStore
}

func (t *DynamoDBTxn) region() string {
	if t.tableStore != nil {
		return t.tableStore.region
	}
	return ""
}

// GetTable retrieves a table by name within the transaction.
func (t *DynamoDBTxn) GetTable(name string) (*Table, error) {
	bucket := t.txn.Bucket(tableBucketName(t.region()))
	data, err := bucket.Get([]byte(name))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrTableNotFound
	}
	var pbTable pb.Table
	if err := proto.Unmarshal(data, &pbTable); err != nil {
		return nil, err
	}
	return ProtoToTable(&pbTable), nil
}

// PutTable stores a table in the transaction.
func (t *DynamoDBTxn) PutTable(table *Table) error {
	bucket := t.txn.Bucket(tableBucketName(t.region()))
	data, err := proto.Marshal(TableToProto(table))
	if err != nil {
		return err
	}
	return bucket.Put([]byte(table.Name), data)
}

// GetItem retrieves an item from a table by its key.
func (t *DynamoDBTxn) GetItem(tableName string, key map[string]*AttributeValue) (*Item, error) {
	table, err := t.GetTable(tableName)
	if err != nil {
		return nil, err
	}
	itemKey := buildItemKeyFromTable(tableName, key, table)
	if itemKey == "" {
		return nil, ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	data, err := bucket.Get([]byte(itemKey))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrItemNotFound
	}
	var pbItem pb.Item
	if err := proto.Unmarshal(data, &pbItem); err != nil {
		return nil, err
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
		return err
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
		return err
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	return bucket.Put([]byte(itemKey), data)
}

// DeleteItem removes an item from a table by its key.
func (t *DynamoDBTxn) DeleteItem(tableName string, key map[string]*AttributeValue) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return err
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
		return false, err
	}
	itemKey := buildItemKeyFromTable(tableName, key, table)
	if itemKey == "" {
		return false, ErrInvalidKey
	}
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	return bucket.Has([]byte(itemKey)), nil
}

// UpdateItemCount increments or decrements the item count for a table.
func (t *DynamoDBTxn) UpdateItemCount(tableName string, delta int64) error {
	if t.tableStore != nil {
		t.tableStore.mu.Lock()
		defer t.tableStore.mu.Unlock()
	}
	table, err := t.GetTable(tableName)
	if err != nil {
		return err
	}
	table.ItemCount += delta
	if table.ItemCount < 0 {
		table.ItemCount = 0
	}
	return t.PutTable(table)
}

// UpdateTableSize increments or decrements the table size in bytes.
func (t *DynamoDBTxn) UpdateTableSize(tableName string, delta int64) error {
	if t.tableStore != nil {
		t.tableStore.mu.Lock()
		defer t.tableStore.mu.Unlock()
	}
	table, err := t.GetTable(tableName)
	if err != nil {
		return err
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
		return tableName + "#" + pkValue + "#" + skValue
	}

	return tableName + "#" + pkValue
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
		return "+"
	}
	rat := new(big.Rat)
	if _, ok := rat.SetString(numStr); !ok {
		return "+" + numStr
	}
	sign := rat.Sign()
	if sign == 0 {
		return "+00000000000000000000"
	}
	absRat := new(big.Rat).Abs(rat)
	numerator := absRat.Num()
	denominator := absRat.Denom()
	floatVal := new(big.Float).SetRat(new(big.Rat).SetFrac(numerator, denominator))
	floatStr := floatVal.Text('f', 20)
	floatStr = strings.TrimRight(strings.TrimRight(floatStr, "0"), ".")
	intPart, fracPart := floatStr, ""
	if dotIdx := strings.Index(floatStr, "."); dotIdx >= 0 {
		intPart = floatStr[:dotIdx]
		fracPart = floatStr[dotIdx+1:]
	}
	if sign > 0 {
		return fmt.Sprintf("+%020d.%s", len(intPart), intPart+fracPart)
	}
	inverted := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(10), big.NewInt(100), nil), numerator)
	return fmt.Sprintf("-%020d.%s", len(intPart), inverted.String()+fracPart)
}

// PutIndexEntries stores index entries for an item in the transaction.
func (t *DynamoDBTxn) PutIndexEntries(tableName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return err
	}
	return t.putIndexEntriesForTable(table, item)
}

// DeleteIndexEntries removes index entries for an item from the transaction.
func (t *DynamoDBTxn) DeleteIndexEntries(tableName string, item *Item) error {
	table, err := t.GetTable(tableName)
	if err != nil {
		return err
	}
	return t.deleteIndexEntriesForTable(table, item)
}

func (t *DynamoDBTxn) putIndexEntriesForTable(table *Table, item *Item) error {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if err := t.putGSIIndexEntry(table, gsi, item); err != nil {
			return err
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if err := t.putLSIIndexEntry(table, lsi, item); err != nil {
			return err
		}
	}
	return nil
}

func (t *DynamoDBTxn) deleteIndexEntriesForTable(table *Table, item *Item) error {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if err := t.deleteGSIIndexEntry(table, gsi, item); err != nil {
			return err
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if err := t.deleteLSIIndexEntry(table, lsi, item); err != nil {
			return err
		}
	}
	return nil
}

func (t *DynamoDBTxn) putGSIIndexEntry(table *Table, gsi *GlobalSecondaryIndex, item *Item) error {
	indexKey := t.buildGSIIndexKey(table, gsi, item)
	if indexKey == "" {
		return nil
	}
	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	bucket := t.txn.Bucket(gsiIndexBucketName(t.region()))
	return bucket.Put([]byte(indexKey), []byte(primaryKey))
}

func (t *DynamoDBTxn) deleteGSIIndexEntry(table *Table, gsi *GlobalSecondaryIndex, item *Item) error {
	indexKey := t.buildGSIIndexKey(table, gsi, item)
	if indexKey == "" {
		return nil
	}
	bucket := t.txn.Bucket(gsiIndexBucketName(t.region()))
	return bucket.Delete([]byte(indexKey))
}

func (t *DynamoDBTxn) putLSIIndexEntry(table *Table, lsi *LocalSecondaryIndex, item *Item) error {
	indexKey := t.buildLSIIndexKey(table, lsi, item)
	if indexKey == "" {
		return nil
	}
	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	bucket := t.txn.Bucket(lsiIndexBucketName(t.region()))
	return bucket.Put([]byte(indexKey), []byte(primaryKey))
}

func (t *DynamoDBTxn) deleteLSIIndexEntry(table *Table, lsi *LocalSecondaryIndex, item *Item) error {
	indexKey := t.buildLSIIndexKey(table, lsi, item)
	if indexKey == "" {
		return nil
	}
	bucket := t.txn.Bucket(lsiIndexBucketName(t.region()))
	return bucket.Delete([]byte(indexKey))
}

func (t *DynamoDBTxn) buildGSIIndexKey(table *Table, gsi *GlobalSecondaryIndex, item *Item) string {
	var hashKeyName, rangeKeyName string
	for _, ks := range gsi.KeySchema {
		if ks.KeyType == KeyTypeHash {
			hashKeyName = ks.AttributeName
		} else if ks.KeyType == KeyTypeRange {
			rangeKeyName = ks.AttributeName
		}
	}

	hashValue := t.getAttributeValueForIndex(item, hashKeyName)
	if hashValue == "" {
		return ""
	}

	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	if primaryKey == "" {
		return ""
	}

	if rangeKeyName != "" {
		rangeValue := t.getAttributeValueForIndex(item, rangeKeyName)
		if rangeValue == "" {
			return ""
		}
		return table.Name + "#" + gsi.IndexName + "#" + hashValue + "#" + rangeValue + "#" + primaryKey
	}
	return table.Name + "#" + gsi.IndexName + "#" + hashValue + "#" + primaryKey
}

func (t *DynamoDBTxn) buildLSIIndexKey(table *Table, lsi *LocalSecondaryIndex, item *Item) string {
	var rangeKeyName string
	for _, ks := range lsi.KeySchema {
		if ks.KeyType == KeyTypeRange {
			rangeKeyName = ks.AttributeName
		}
	}

	var tableHashKeyName string
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeHash {
			tableHashKeyName = ks.AttributeName
			break
		}
	}

	hashValue := attributeValueToString(item.Key[tableHashKeyName])
	if hashValue == "" {
		return ""
	}

	rangeValue := t.getAttributeValueForIndex(item, rangeKeyName)
	if rangeValue == "" {
		return ""
	}

	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	if primaryKey == "" {
		return ""
	}

	return table.Name + "#" + lsi.IndexName + "#" + hashValue + "#" + rangeValue + "#" + primaryKey
}

func (t *DynamoDBTxn) getAttributeValueForIndex(item *Item, attrName string) string {
	if item.Key != nil && item.Key[attrName] != nil {
		return attributeValueToString(item.Key[attrName])
	}
	if item.Attributes != nil && item.Attributes[attrName] != nil {
		return attributeValueToString(item.Attributes[attrName])
	}
	return ""
}

// QueryByGSI queries items from a global secondary index.
func (t *DynamoDBTxn) QueryByGSI(tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error) {
	_, err := t.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	prefix := tableName + "#" + indexName + "#" + hashKeyValue + "#"
	bucket := t.txn.Bucket(gsiIndexBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	var items []*Item

	for iter.Next() {
		if !opts.Reverse && opts.Limit > 0 && len(items) >= opts.Limit {
			break
		}

		primaryKey := string(iter.Value())
		itemBucket := t.txn.Bucket(itemBucketName(t.region()))
		data, err := itemBucket.Get([]byte(primaryKey))
		if err != nil {
			return nil, fmt.Errorf("failed to get item from index %s key %s: %w", indexName, primaryKey, err)
		}
		if data == nil {
			continue
		}

		var pbItem pb.Item
		if err := proto.Unmarshal(data, &pbItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item from index %s key %s: %w", indexName, primaryKey, err)
		}

		items = append(items, &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		})
	}

	if opts.Reverse {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
		if opts.Limit > 0 && len(items) > opts.Limit {
			items = items[:opts.Limit]
		}
	}

	return items, iter.Error()
}

// QueryByLSI queries items from a local secondary index.
func (t *DynamoDBTxn) QueryByLSI(tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error) {
	_, err := t.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	prefix := tableName + "#" + indexName + "#" + hashKeyValue + "#"
	bucket := t.txn.Bucket(lsiIndexBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	var items []*Item

	for iter.Next() {
		if !opts.Reverse && opts.Limit > 0 && len(items) >= opts.Limit {
			break
		}

		primaryKey := string(iter.Value())
		itemBucket := t.txn.Bucket(itemBucketName(t.region()))
		data, err := itemBucket.Get([]byte(primaryKey))
		if err != nil {
			return nil, fmt.Errorf("failed to get item from index %s key %s: %w", indexName, primaryKey, err)
		}
		if data == nil {
			continue
		}

		var pbItem pb.Item
		if err := proto.Unmarshal(data, &pbItem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item from index %s key %s: %w", indexName, primaryKey, err)
		}

		items = append(items, &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		})
	}

	if opts.Reverse {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
		if opts.Limit > 0 && len(items) > opts.Limit {
			items = items[:opts.Limit]
		}
	}

	return items, iter.Error()
}

// IndexQueryOptions defines options for querying indexes.
type IndexQueryOptions struct {
	Limit   int
	Reverse bool
}

// Scan scans all items in a table within the transaction.
func (t *DynamoDBTxn) Scan(tableName string, fn func(item *Item) error) error {
	prefix := tableName + "#"
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		var pbItem pb.Item
		if err := proto.Unmarshal(iter.Value(), &pbItem); err != nil {
			continue
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
		return err
	}

	prefix := tableName + "#" + partitionKeyValue
	hasSortKey := false
	pkName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeRange {
			prefix += "#"
			hasSortKey = true
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
			continue
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		if !hasSortKey {
			itemPkValue := attributeValueToString(item.Key[pkName])
			if itemPkValue != partitionKeyValue {
				continue
			}
		}
		if err := fn(item); err != nil {
			return err
		}
	}
	return iter.Error()
}

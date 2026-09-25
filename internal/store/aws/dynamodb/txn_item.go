package dynamodb

import (
	"fmt"

	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"

	"google.golang.org/protobuf/proto"
)

// ---------------------------------------------------------------------------
// Transaction item plane — the DynamoDBTxn methods that read, write, and
// walk items: single-key CRUD, the unified write paths every item-mutating
// operation routes through (index-entry swap plus metric deltas), and the
// full-table and partition-key scans. The transaction wrapper and the
// shared GetTable read primitive live in store.go; the journal hooks the
// write paths append live in journal_store.go.
// ---------------------------------------------------------------------------

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

// Scan scans all items in a table within the transaction, through the
// shared walk core every item walk of this store runs on.
func (t *DynamoDBTxn) Scan(tableName string, fn func(item *Item) error) error {
	bucket := t.txn.Bucket(itemBucketName(t.region()))
	_, err := bucketScanWalk(bucket, tableName+KeySep, ScanOptions{}, nil, fn)
	return err
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
	// The partition recheck runs as the walk's accept filter: the encoded
	// partition component is prefix-free, so the scan cannot cross into
	// another partition, and the recheck keeps the delivered set exact.
	_, err = bucketScanWalk(bucket, prefix, ScanOptions{}, func(item *Item) bool {
		return EncodeKeyValue(item.Key[pkName]) == partitionKeyValue
	}, fn)
	return err
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

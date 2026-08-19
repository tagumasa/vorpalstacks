// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
)

var errScanLimitReached = errors.New("scan limit reached")

func itemBucketName(region string) string {
	return "dynamodb_items-" + region
}

func gsiIndexBucketName(region string) string {
	return "dynamodb_gsi_index-" + region
}

func lsiIndexBucketName(region string) string {
	return "dynamodb_lsi_index-" + region
}

func tagMainBucketName(region string) string {
	return "dynamodb-tags-" + region
}

func tagIndexBucketName(region string) string {
	return "dynamodb-tag-idx-" + region
}

// ItemStore manages DynamoDB table items in persistent storage.
type ItemStore struct {
	*common.BaseStore
	tableStore *TableStore
	storage    storage.BasicStorage
	region     string
}

// NewItemStore creates a new store for DynamoDB items.
func NewItemStore(store storage.BasicStorage, tableStore *TableStore) *ItemStore {
	region := ""
	if tableStore != nil {
		region = tableStore.region
	}
	s := &ItemStore{
		BaseStore:  common.NewBaseStore(store.Bucket(itemBucketName(region)), "dynamodb_items"),
		tableStore: tableStore,
		storage:    store,
		region:     region,
	}
	return s
}

func (s *ItemStore) buildItemKey(tableName string, key map[string]*AttributeValue) string {
	table, err := s.tableStore.Get(tableName)
	if err != nil {
		return ""
	}

	pkName := s.tableStore.GetPartitionKey(table)
	skName := s.tableStore.GetSortKey(table)

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
		return fmt.Sprintf("%s"+keySep+"%s"+keySep+"%s", tableName, pkValue, skValue)
	}

	return fmt.Sprintf("%s"+keySep+"%s", tableName, pkValue)
}

// Get retrieves a DynamoDB item by table name and key.
func (s *ItemStore) Get(tableName string, key map[string]*AttributeValue) (*Item, error) {
	itemKey := s.buildItemKey(tableName, key)
	if itemKey == "" {
		return nil, ErrInvalidKey
	}

	var pbItem pb.Item
	if err := s.BaseStore.GetProto(itemKey, &pbItem); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &Item{
		TableName:  pbItem.TableName,
		Key:        protoToAttributeValueMapDirect(pbItem.Key),
		Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
	}, nil
}

// Exists checks if a DynamoDB item exists.
func (s *ItemStore) Exists(tableName string, key map[string]*AttributeValue) bool {
	itemKey := s.buildItemKey(tableName, key)
	return itemKey != "" && s.BaseStore.Exists(itemKey)
}

// List returns a list of DynamoDB items with pagination.
func (s *ItemStore) List(tableName string, marker string, limit int) ([]*Item, string, error) {
	prefix := tableName + keySep
	var items []*Item
	var lastKey string

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if marker != "" && key <= marker {
			return nil
		}
		if limit > 0 && len(items) >= limit {
			lastKey = key
			return errScanLimitReached
		}

		var pbItem pb.Item
		if err := proto.Unmarshal(value, &pbItem); err != nil {
			return err
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		items = append(items, item)
		return nil
	})

	if err != nil && !errors.Is(err, errScanLimitReached) {
		return nil, "", err
	}
	return items, lastKey, nil
}

// ScanOptions controls the behaviour of a storage-level scan.
type ScanOptions struct {
	Limit  int
	Marker string
}

// Scan scans all items in a DynamoDB table.
func (s *ItemStore) Scan(tableName string, fn func(item *Item) error) error {
	_, err := s.ScanWithOptions(tableName, ScanOptions{}, func(item *Item) error {
		return fn(item)
	})
	return err
}

// ScanWithOptions scans items with limit and marker support for pagination.
func (s *ItemStore) ScanWithOptions(tableName string, opts ScanOptions, fn func(item *Item) error) (string, error) {
	prefix := tableName + keySep
	var lastKey string
	count := 0

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if opts.Marker != "" && key <= opts.Marker {
			return nil
		}
		if opts.Limit > 0 && count >= opts.Limit {
			lastKey = key
			return errScanLimitReached
		}

		var pbItem pb.Item
		if err := proto.Unmarshal(value, &pbItem); err != nil {
			return err
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		count++
		return fn(item)
	})

	if err != nil && !errors.Is(err, errScanLimitReached) {
		return "", err
	}
	return lastKey, nil
}

// ScanByPartitionKey scans items with a specific partition key value.
func (s *ItemStore) ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error {
	table, err := s.tableStore.Get(tableName)
	if err != nil {
		return err
	}
	_, err = s.scanByPartitionKeyWithTable(tableName, table, partitionKeyValue, ScanOptions{}, fn)
	return err
}

// ScanByPartitionKeyWithTable scans items with a specific partition key value using a pre-fetched table,
// avoiding a redundant table store lookup.
func (s *ItemStore) ScanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, opts ScanOptions, fn func(item *Item) error) (string, error) {
	return s.scanByPartitionKeyWithTable(tableName, table, partitionKeyValue, opts, fn)
}

func (s *ItemStore) scanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, opts ScanOptions, fn func(item *Item) error) (string, error) {
	prefix := tableName + keySep + partitionKeyValue
	hasSortKey := s.tableStore.GetSortKey(table) != ""
	if hasSortKey {
		prefix += keySep
	}

	pkName := s.tableStore.GetPartitionKey(table)
	var lastKey string
	count := 0

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if opts.Marker != "" && key <= opts.Marker {
			return nil
		}
		var pbItem pb.Item
		if err := proto.Unmarshal(value, &pbItem); err != nil {
			return err
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		itemPkValue := attributeValueToString(item.Key[pkName])
		if itemPkValue != partitionKeyValue {
			return nil
		}
		if opts.Limit > 0 && count >= opts.Limit {
			lastKey = key
			return errScanLimitReached
		}
		count++
		return fn(item)
	})

	if err != nil && !errors.Is(err, errScanLimitReached) {
		return "", err
	}
	return lastKey, nil
}

// Count returns the number of items in a DynamoDB table.
func (s *ItemStore) Count(tableName string) (int64, error) {
	var count int64
	prefix := tableName + keySep
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		count++
		return nil
	})
	return count, err
}

// DeleteAllForTable removes all items from a DynamoDB table.
// Items are deleted first; GSI/LSI indexes are cleaned afterwards.
// If index cleanup fails, orphan index entries may remain, but this is
// preferable to losing GSI/LSI data while items still exist.
func (s *ItemStore) DeleteAllForTable(tableName string) error {
	prefix := tableName + keySep

	const batchSize = 1000
	var keysBatch []string

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		keysBatch = append(keysBatch, key)
		if len(keysBatch) >= batchSize {
			for _, k := range keysBatch {
				if delErr := s.BaseStore.Delete(k); delErr != nil {
					return delErr
				}
			}
			keysBatch = keysBatch[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, key := range keysBatch {
		if delErr := s.BaseStore.Delete(key); delErr != nil {
			return delErr
		}
	}

	gsiBucket := s.storage.Bucket(gsiIndexBucketName(s.region))
	if gsiBucket != nil {
		if err := s.deleteByPrefix(gsiBucket, prefix); err != nil {
			return err
		}
	}

	lsiBucket := s.storage.Bucket(lsiIndexBucketName(s.region))
	if lsiBucket != nil {
		if err := s.deleteByPrefix(lsiBucket, prefix); err != nil {
			return err
		}
	}

	table, err := s.tableStore.Get(tableName)
	if err != nil {
		if !IsTableNotFound(err) {
			return fmt.Errorf("get table for size reset: %w", err)
		}
		return nil
	}
	if err := s.tableStore.UpdateTableSize(tableName, -table.TableSizeBytes); err != nil {
		return fmt.Errorf("reset table size after delete all: %w", err)
	}
	if err := s.tableStore.UpdateItemCount(tableName, -table.ItemCount); err != nil {
		return fmt.Errorf("reset item count after delete all: %w", err)
	}
	return nil
}

func (s *ItemStore) deleteByPrefix(bucket storage.Bucket, prefix string) error {
	var keysBatch []string
	const batchSize = 1000

	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	for iter.Next() {
		keysBatch = append(keysBatch, string(iter.Key()))
		if len(keysBatch) >= batchSize {
			for _, k := range keysBatch {
				if delErr := bucket.Delete([]byte(k)); delErr != nil {
					return delErr
				}
			}
			keysBatch = keysBatch[:0]
		}
	}

	if err := iter.Error(); err != nil {
		return err
	}

	for _, k := range keysBatch {
		if delErr := bucket.Delete([]byte(k)); delErr != nil {
			return delErr
		}
	}
	return nil
}

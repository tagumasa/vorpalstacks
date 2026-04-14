// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
)

func itemBucketName(region string) string {
	return "dynamodb_items-" + region
}

func gsiIndexBucketName(region string) string {
	return "dynamodb_gsi_index-" + region
}

func lsiIndexBucketName(region string) string {
	return "dynamodb_lsi_index-" + region
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
	return &ItemStore{
		BaseStore:  common.NewBaseStore(store.Bucket(itemBucketName(region)), "dynamodb_items"),
		tableStore: tableStore,
		storage:    store,
		region:     region,
	}
}

func (s *ItemStore) buildItemKey(tableName string, key map[string]*AttributeValue) string {
	table, err := s.tableStore.Get(tableName)
	if err != nil {
		return ""
	}

	pkName := s.tableStore.GetPartitionKey(table)
	skName := s.tableStore.GetSortKey(table)

	pkValue := s.attributeValueToString(key[pkName])
	if pkValue == "" {
		return ""
	}

	if skName != "" {
		if key[skName] == nil {
			return ""
		}
		skValue := s.attributeValueToString(key[skName])
		if skValue == "" {
			return ""
		}
		return fmt.Sprintf("%s#%s#%s", tableName, pkValue, skValue)
	}

	return fmt.Sprintf("%s#%s", tableName, pkValue)
}

func (s *ItemStore) attributeValueToString(av *AttributeValue) string {
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

// Put stores a DynamoDB item in a table.
func (s *ItemStore) Put(tableName string, key map[string]*AttributeValue, attributes map[string]*AttributeValue) (*Item, error) {
	itemKey := s.buildItemKey(tableName, key)
	if itemKey == "" {
		return nil, ErrInvalidKey
	}

	if attributes == nil {
		attributes = make(map[string]*AttributeValue)
	}
	for k, v := range key {
		attributes[k] = v
	}

	pbItem := &pb.Item{
		TableName:  tableName,
		Key:        attributeValueMapToProtoDirect(key),
		Attributes: attributeValueMapToProtoDirect(attributes),
	}

	if err := s.BaseStore.PutProto(itemKey, pbItem); err != nil {
		return nil, err
	}

	return &Item{
		TableName:  tableName,
		Key:        key,
		Attributes: attributes,
	}, nil
}

// Delete removes a DynamoDB item by table name and key.
func (s *ItemStore) Delete(tableName string, key map[string]*AttributeValue) error {
	itemKey := s.buildItemKey(tableName, key)
	if itemKey == "" {
		return ErrInvalidKey
	}
	return s.BaseStore.Delete(itemKey)
}

// Exists checks if a DynamoDB item exists.
func (s *ItemStore) Exists(tableName string, key map[string]*AttributeValue) bool {
	itemKey := s.buildItemKey(tableName, key)
	return itemKey != "" && s.BaseStore.Exists(itemKey)
}

// List returns a list of DynamoDB items with pagination.
func (s *ItemStore) List(tableName string, marker string, limit int) ([]*Item, string, error) {
	prefix := tableName + "#"
	var items []*Item
	var lastKey string

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if marker != "" && key <= marker {
			return nil
		}
		if limit > 0 && len(items) >= limit {
			lastKey = key
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
		items = append(items, item)
		return nil
	})

	if err != nil {
		return nil, "", err
	}
	return items, lastKey, nil
}

// Scan scans all items in a DynamoDB table.
func (s *ItemStore) Scan(tableName string, fn func(item *Item) error) error {
	prefix := tableName + "#"
	return s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbItem pb.Item
		if err := proto.Unmarshal(value, &pbItem); err != nil {
			return err
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		return fn(item)
	})
}

// ScanByPartitionKey scans items with a specific partition key value.
func (s *ItemStore) ScanByPartitionKey(tableName, partitionKeyValue string, fn func(item *Item) error) error {
	table, err := s.tableStore.Get(tableName)
	if err != nil {
		return err
	}
	return s.scanByPartitionKeyWithTable(tableName, table, partitionKeyValue, fn)
}

// ScanByPartitionKeyWithTable scans items with a specific partition key value using a pre-fetched table,
// avoiding a redundant table store lookup.
func (s *ItemStore) ScanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, fn func(item *Item) error) error {
	return s.scanByPartitionKeyWithTable(tableName, table, partitionKeyValue, fn)
}

func (s *ItemStore) scanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, fn func(item *Item) error) error {
	prefix := tableName + "#" + partitionKeyValue
	hasSortKey := s.tableStore.GetSortKey(table) != ""
	if hasSortKey {
		prefix += "#"
	}

	pkName := s.tableStore.GetPartitionKey(table)

	return s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbItem pb.Item
		if err := proto.Unmarshal(value, &pbItem); err != nil {
			return err
		}
		item := &Item{
			TableName:  pbItem.TableName,
			Key:        protoToAttributeValueMapDirect(pbItem.Key),
			Attributes: protoToAttributeValueMapDirect(pbItem.Attributes),
		}
		if !hasSortKey {
			itemPkValue := s.attributeValueToString(item.Key[pkName])
			if itemPkValue != partitionKeyValue {
				return nil
			}
		}
		return fn(item)
	})
}

// Count returns the number of items in a DynamoDB table.
func (s *ItemStore) Count(tableName string) (int64, error) {
	var count int64
	prefix := tableName + "#"
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
	prefix := tableName + "#"

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

	if table, err := s.tableStore.Get(tableName); err == nil {
		if err := s.tableStore.UpdateTableSize(tableName, -table.TableSizeBytes); err != nil {
			return err
		}
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

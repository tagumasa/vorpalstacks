package dynamodb

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// IndexStore manages GSI and LSI secondary index entries for DynamoDB
// tables. It encapsulates key construction, CRUD operations, and
// querying for secondary indexes, following the same separation pattern
// used by CloudTrail (eventIDIndexStore / arnIndexStore) and AppSync
// (mergedApiAssocIndexStore).
//
// All mutation methods accept a storage.Transaction so that index updates
// are atomic with the corresponding item writes within the caller's
// transaction.
type IndexStore struct {
	region string
}

// NewIndexStore creates a new IndexStore for the given region.
func NewIndexStore(region string) *IndexStore {
	return &IndexStore{region: region}
}

// Region returns the region this IndexStore operates on.
func (s *IndexStore) Region() string {
	return s.region
}

// ---------------------------------------------------------------------------
// Key construction
// ---------------------------------------------------------------------------

// BuildGSIKey constructs the Pebble key for a GSI index entry.
// Returns "" when the item lacks the GSI hash or range key attributes.
func (s *IndexStore) BuildGSIKey(table *Table, gsi *GlobalSecondaryIndex, item *Item) string {
	var hashKeyName, rangeKeyName string
	for _, ks := range gsi.KeySchema {
		if ks.KeyType == KeyTypeHash {
			hashKeyName = ks.AttributeName
		} else if ks.KeyType == KeyTypeRange {
			rangeKeyName = ks.AttributeName
		}
	}

	hashValue := s.getAttributeValueForIndex(item, hashKeyName)
	if hashValue == "" {
		return ""
	}

	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	if primaryKey == "" {
		return ""
	}

	if rangeKeyName != "" {
		rangeValue := s.getAttributeValueForIndex(item, rangeKeyName)
		if rangeValue == "" {
			return ""
		}
		return table.Name + keySep + gsi.IndexName + keySep + hashValue + keySep + rangeValue + keySep + primaryKey
	}
	return table.Name + keySep + gsi.IndexName + keySep + hashValue + keySep + primaryKey
}

// BuildLSIKey constructs the Pebble key for an LSI index entry.
// Returns "" when the item lacks the required key attributes.
func (s *IndexStore) BuildLSIKey(table *Table, lsi *LocalSecondaryIndex, item *Item) string {
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

	rangeValue := s.getAttributeValueForIndex(item, rangeKeyName)
	if rangeValue == "" {
		return ""
	}

	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	if primaryKey == "" {
		return ""
	}

	return table.Name + keySep + lsi.IndexName + keySep + hashValue + keySep + rangeValue + keySep + primaryKey
}

func (s *IndexStore) getAttributeValueForIndex(item *Item, attrName string) string {
	if item.Key != nil && item.Key[attrName] != nil {
		return attributeValueToString(item.Key[attrName])
	}
	if item.Attributes != nil && item.Attributes[attrName] != nil {
		return attributeValueToString(item.Attributes[attrName])
	}
	return ""
}

// ---------------------------------------------------------------------------
// Index entry CRUD (transactional)
// ---------------------------------------------------------------------------

// PutIndexEntries writes all GSI and LSI index entries for an item
// within the given transaction.
func (s *IndexStore) PutIndexEntries(txn storage.Transaction, table *Table, item *Item) error {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if err := s.putGSIEntry(txn, table, gsi, item); err != nil {
			return err
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if err := s.putLSIEntry(txn, table, lsi, item); err != nil {
			return err
		}
	}
	return nil
}

// DeleteIndexEntries removes all GSI and LSI index entries for an item
// within the given transaction.
func (s *IndexStore) DeleteIndexEntries(txn storage.Transaction, table *Table, item *Item) error {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if err := s.deleteGSIEntry(txn, table, gsi, item); err != nil {
			return err
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if err := s.deleteLSIEntry(txn, table, lsi, item); err != nil {
			return err
		}
	}
	return nil
}

func (s *IndexStore) putGSIEntry(txn storage.Transaction, table *Table, gsi *GlobalSecondaryIndex, item *Item) error {
	indexKey := s.BuildGSIKey(table, gsi, item)
	if indexKey == "" {
		return nil
	}
	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	bucket := txn.Bucket(gsiIndexBucketName(s.region))
	return bucket.Put([]byte(indexKey), []byte(primaryKey))
}

func (s *IndexStore) deleteGSIEntry(txn storage.Transaction, table *Table, gsi *GlobalSecondaryIndex, item *Item) error {
	indexKey := s.BuildGSIKey(table, gsi, item)
	if indexKey == "" {
		return nil
	}
	bucket := txn.Bucket(gsiIndexBucketName(s.region))
	return bucket.Delete([]byte(indexKey))
}

func (s *IndexStore) putLSIEntry(txn storage.Transaction, table *Table, lsi *LocalSecondaryIndex, item *Item) error {
	indexKey := s.BuildLSIKey(table, lsi, item)
	if indexKey == "" {
		return nil
	}
	primaryKey := buildItemKeyFromTable(table.Name, item.Key, table)
	bucket := txn.Bucket(lsiIndexBucketName(s.region))
	return bucket.Put([]byte(indexKey), []byte(primaryKey))
}

func (s *IndexStore) deleteLSIEntry(txn storage.Transaction, table *Table, lsi *LocalSecondaryIndex, item *Item) error {
	indexKey := s.BuildLSIKey(table, lsi, item)
	if indexKey == "" {
		return nil
	}
	bucket := txn.Bucket(lsiIndexBucketName(s.region))
	return bucket.Delete([]byte(indexKey))
}

// ---------------------------------------------------------------------------
// Index query (transactional read)
// ---------------------------------------------------------------------------

// QueryGSI queries a global secondary index for items matching the hash key.
func (s *IndexStore) QueryGSI(txn storage.Transaction, tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error) {
	return s.queryByIndex(txn, tableName, indexName, hashKeyValue, gsiIndexBucketName(s.region), opts)
}

// QueryLSI queries a local secondary index for items matching the hash key.
func (s *IndexStore) QueryLSI(txn storage.Transaction, tableName, indexName, hashKeyValue string, opts IndexQueryOptions) ([]*Item, error) {
	return s.queryByIndex(txn, tableName, indexName, hashKeyValue, lsiIndexBucketName(s.region), opts)
}

func (s *IndexStore) queryByIndex(txn storage.Transaction, tableName, indexName, hashKeyValue, bucketName string, opts IndexQueryOptions) ([]*Item, error) {
	prefix := tableName + keySep + indexName + keySep + hashKeyValue + keySep
	bucket := txn.Bucket(bucketName)
	iter := bucket.ScanPrefix([]byte(prefix))
	defer iter.Close()

	var items []*Item

	for iter.Next() {
		if !opts.Reverse && opts.Limit > 0 && len(items) >= opts.Limit {
			break
		}

		primaryKey := string(iter.Value())
		itemBucket := txn.Bucket(itemBucketName(s.region))
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

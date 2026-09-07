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
	hashKeyName, rangeKeyName := schemaKeyNames(gsi.KeySchema)

	hashValue := s.getAttributeValueForIndex(item, hashKeyName)
	if hashValue == "" {
		return ""
	}

	primaryKey := EncodeItemKey(table.Name, item.Key, table)
	if primaryKey == "" {
		return ""
	}

	if rangeKeyName != "" {
		rangeValue := s.getAttributeValueForIndex(item, rangeKeyName)
		if rangeValue == "" {
			return ""
		}
		return table.Name + KeySep + gsi.IndexName + KeySep + hashValue + KeySep + rangeValue + KeySep + primaryKey
	}
	return table.Name + KeySep + gsi.IndexName + KeySep + hashValue + KeySep + primaryKey
}

// BuildLSIKey constructs the Pebble key for an LSI index entry.
// Returns "" when the item lacks the required key attributes.
func (s *IndexStore) BuildLSIKey(table *Table, lsi *LocalSecondaryIndex, item *Item) string {
	_, rangeKeyName := schemaKeyNames(lsi.KeySchema)
	tableHashKeyName, _ := schemaKeyNames(table.KeySchema)

	hashValue := EncodeKeyValue(item.Key[tableHashKeyName])
	if hashValue == "" {
		return ""
	}

	rangeValue := s.getAttributeValueForIndex(item, rangeKeyName)
	if rangeValue == "" {
		return ""
	}

	primaryKey := EncodeItemKey(table.Name, item.Key, table)
	if primaryKey == "" {
		return ""
	}

	return table.Name + KeySep + lsi.IndexName + KeySep + hashValue + KeySep + rangeValue + KeySep + primaryKey
}

func (s *IndexStore) getAttributeValueForIndex(item *Item, attrName string) string {
	if item.Key != nil && item.Key[attrName] != nil {
		return EncodeKeyValue(item.Key[attrName])
	}
	if item.Attributes != nil && item.Attributes[attrName] != nil {
		return EncodeKeyValue(item.Attributes[attrName])
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
	primaryKey := EncodeItemKey(table.Name, item.Key, table)
	bucket := txn.Bucket(gsiIndexBucketName(s.region))
	return bucket.Put([]byte(indexKey), []byte(primaryKey))
}

// PutIndexEntriesForIndex writes only the named GSI's index entry for an
// item within the given transaction, leaving every other index untouched.
// A name absent from the table's schema writes nothing: the index may have
// been deleted between the schema update and this backfill write.
func (s *IndexStore) PutIndexEntriesForIndex(txn storage.Transaction, table *Table, indexName string, item *Item) error {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName != indexName {
			continue
		}
		return s.putGSIEntry(txn, table, gsi, item)
	}
	return nil
}

// DeleteIndexEntriesForIndex removes every index entry of the named GSI:
// deleting an index removes its data with it, so the entries must not
// outlive the index (a re-created same-name index would otherwise inherit
// stale entries pointing at items that no longer hold the index key
// attributes). Keys are deleted in batches to bound the delete set,
// mirroring DynamoDBTxn.deleteAllByPrefix.
func (s *IndexStore) DeleteIndexEntriesForIndex(txn storage.Transaction, tableName, indexName string) error {
	return deletePrefixBatched(txn.Bucket(gsiIndexBucketName(s.region)), tableName+KeySep+indexName+KeySep)
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
	primaryKey := EncodeItemKey(table.Name, item.Key, table)
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

// queryByIndex walks one hash range of a secondary index in the direction
// opts selects. Index keys encode the sort key with the sort-correct value
// encoder, so storage order is semantic order: a forward walk yields
// ascending sort keys, a reverse walk descending. The limit therefore
// bounds both directions, and the filter runs before the limit is counted
// so a filtered query stops as soon as its page is full.
func (s *IndexStore) queryByIndex(txn storage.Transaction, tableName, indexName, hashKeyValue, bucketName string, opts IndexQueryOptions) ([]*Item, error) {
	prefix := tableName + KeySep + indexName + KeySep + hashKeyValue + KeySep
	bucket := txn.Bucket(bucketName)

	var items []*Item

	collect := func(indexKey string, primaryKey []byte) (bool, error) {
		itemBucket := txn.Bucket(itemBucketName(s.region))
		data, err := itemBucket.Get(primaryKey)
		if err != nil {
			return false, fmt.Errorf("failed to get item from index %s key %s: %w", indexName, indexKey, err)
		}
		if data == nil {
			return false, nil
		}

		var pbItem pb.Item
		if err := proto.Unmarshal(data, &pbItem); err != nil {
			return false, fmt.Errorf("failed to unmarshal item from index %s key %s: %w", indexName, indexKey, err)
		}

		item := itemFromProto(&pbItem)
		if opts.Filter != nil && !opts.Filter(item) {
			return false, nil
		}
		items = append(items, item)
		if opts.Limit > 0 && len(items) >= opts.Limit {
			return true, nil
		}
		return false, nil
	}

	var iter storage.Iterator
	if opts.Reverse {
		var before []byte
		if opts.Marker != "" {
			before = []byte(opts.Marker)
		}
		iter = bucket.ScanPrefixReverse([]byte(prefix), before)
	} else {
		iter = bucket.ScanPrefix([]byte(prefix))
	}
	defer iter.Close()

	for iter.Next() {
		indexKey := string(iter.Key())
		if !opts.Reverse && opts.Marker != "" && indexKey <= opts.Marker {
			continue
		}
		stop, err := collect(indexKey, iter.Value())
		if err != nil {
			return nil, err
		}
		if stop {
			break
		}
	}

	return items, iter.Error()
}

// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"errors"

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

func vectorIndexBucketName(region string) string {
	return "dynamodb_vector_index-" + region
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
	return EncodeItemKey(tableName, key, table)
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
	return itemFromProto(&pbItem), nil
}

// Exists checks if a DynamoDB item exists.
func (s *ItemStore) Exists(tableName string, key map[string]*AttributeValue) bool {
	itemKey := s.buildItemKey(tableName, key)
	return itemKey != "" && s.BaseStore.Exists(itemKey)
}

// List returns a list of DynamoDB items with pagination. The returned
// marker names the last item the page delivered; an empty marker means the
// walk is exhausted. Resuming with the marker delivers the remaining items
// exactly once.
func (s *ItemStore) List(tableName string, marker string, limit int) ([]*Item, string, error) {
	var items []*Item
	lastKey, err := s.scanWalk(tableName+KeySep, ScanOptions{Limit: limit, Marker: marker}, nil, func(item *Item) error {
		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return items, lastKey, nil
}

// ScanOptions controls the behaviour of a storage-level scan. Filter, when
// non-nil, is applied to each item during iteration — before the Limit is
// counted — so a filtered scan reads only as far as its page requires.
// Reverse walks descending; its Marker names the key to start strictly
// before (the forward walk resumes strictly after it).
type ScanOptions struct {
	Limit   int
	Marker  string
	Reverse bool
	Filter  func(*Item) bool
}

// Scan scans all items in a DynamoDB table.
func (s *ItemStore) Scan(tableName string, fn func(item *Item) error) error {
	_, err := s.ScanWithOptions(tableName, ScanOptions{}, func(item *Item) error {
		return fn(item)
	})
	return err
}

// scanWalk drives one bounded item walk over a key prefix: the marker
// resume, the reverse direction, the per-item decode, the limit, and the
// errScanLimitReached epilogue are owned here. accept, when non-nil,
// decides whether a decoded item belongs to the page (partition membership,
// filter) — accepted items are delivered to fn and count towards the limit,
// rejected items are skipped without consuming it. The returned marker
// names the last delivered key and is set only when the limit stopped the
// walk; the forward walk resumes strictly after it, the reverse walk
// strictly before it.
func (s *ItemStore) scanWalk(prefix string, opts ScanOptions, accept func(*Item) bool, fn func(*Item) error) (string, error) {
	var lastKey string
	count := 0

	visit := func(key string, value []byte) error {
		if opts.Limit > 0 && count >= opts.Limit {
			return errScanLimitReached
		}

		var pbItem pb.Item
		if err := proto.Unmarshal(value, &pbItem); err != nil {
			return err
		}
		item := itemFromProto(&pbItem)
		if accept != nil && !accept(item) {
			return nil
		}
		count++
		if err := fn(item); err != nil {
			return err
		}
		lastKey = key
		return nil
	}

	var err error
	if opts.Reverse {
		err = s.BaseStore.ScanPrefixReverse(prefix, opts.Marker, visit)
	} else {
		err = s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
			if opts.Marker != "" && key <= opts.Marker {
				return nil
			}
			return visit(key, value)
		})
	}

	if err != nil && !errors.Is(err, errScanLimitReached) {
		return "", err
	}
	if errors.Is(err, errScanLimitReached) {
		return lastKey, nil
	}
	return "", nil
}

// ScanWithOptions scans items with limit and marker support for pagination.
// The returned marker names the last item the page delivered and is set
// only when the limit stopped the walk — an empty marker means the walk is
// exhausted. Reverse walks descending and resume strictly before the
// marker, the same direction contract the partition scan honours.
func (s *ItemStore) ScanWithOptions(tableName string, opts ScanOptions, fn func(*Item) error) (string, error) {
	return s.scanWalk(tableName+KeySep, opts, opts.Filter, fn)
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

// scanByPartitionKeyWithTable scans one partition in the direction opts
// selects. partitionKeyValue must already be an EncodeKeyValue rendering;
// the encoded component is prefix-free, so no trailing separator is needed
// to keep the scan inside the partition. Storage keys encode the sort key
// with the sort-correct value encoder, so a forward walk yields ascending
// sort keys and a reverse walk descending ones — the limit bounds both
// directions, and the filter runs before the limit is counted. The returned
// marker names the last item the page delivered and is set only when the
// limit stopped the walk; the forward walk resumes strictly after it, the
// reverse walk strictly before it.
func (s *ItemStore) scanByPartitionKeyWithTable(tableName string, table *Table, partitionKeyValue string, opts ScanOptions, fn func(item *Item) error) (string, error) {
	pkName, _ := schemaKeyNames(table.KeySchema)
	return s.scanWalk(tableName+KeySep+partitionKeyValue, opts, func(item *Item) bool {
		if EncodeKeyValue(item.Key[pkName]) != partitionKeyValue {
			return false
		}
		return opts.Filter == nil || opts.Filter(item)
	}, fn)
}

// Count returns the number of items in a DynamoDB table.
func (s *ItemStore) Count(tableName string) (int64, error) {
	var count int64
	prefix := tableName + KeySep
	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		count++
		return nil
	})
	return count, err
}

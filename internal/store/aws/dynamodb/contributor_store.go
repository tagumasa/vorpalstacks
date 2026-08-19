package dynamodb

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// Contributor access aggregation backs CloudWatch Contributor Insights for
// DynamoDB: every item read or write on a table with contributor insights
// enabled is counted per tracked key, so the most accessed partition (and
// partition+sort) keys can be reported.

const (
	// ContributorLayoutPartitionKey tracks partition keys only.
	ContributorLayoutPartitionKey = "PKC"
	// ContributorLayoutFullKey tracks partition and sort key pairs.
	ContributorLayoutFullKey = "SKC"
	// ContributorReadUnits is the ConsumedThroughputUnits weight of one read
	// event: (3 x write capacity units) + consumed read capacity units.
	ContributorReadUnits = 1.0
	// ContributorWriteUnits is the ConsumedThroughputUnits weight of one
	// write event: (3 x write capacity units) + consumed read capacity units.
	ContributorWriteUnits = 3.0
)

// ContributorWriteEvent is one item write observed on a table with
// contributor insights enabled, applied to the access counters after the
// carrying item transaction commits.
type ContributorWriteEvent struct {
	TableName string
	Key       map[string]*AttributeValue
}

// ContributorKeyStat is one aggregated key of one table.
type ContributorKeyStat struct {
	Key   string  `json:"key"`
	Count int64   `json:"count"`
	Units float64 `json:"units"`
}

type contributorAccess struct {
	Count int64   `json:"count"`
	Units float64 `json:"units"`
}

// ContributorStore manages the per-key access counters.
type ContributorStore struct {
	*common.BaseStore
}

func contributorBucketName(region string) string {
	return "dynamodb_contributors-" + region
}

// NewContributorStore creates a contributor access counter store for the
// given region.
func NewContributorStore(store storage.BasicStorage, region string) *ContributorStore {
	return &ContributorStore{
		BaseStore: common.NewBaseStore(store.Bucket(contributorBucketName(region)), "dynamodb-contributors"),
	}
}

// ContributorLayouts lists the key layouts tracked for a table: partition
// keys, plus partition and sort key pairs when the key schema has a range
// element.
func ContributorLayouts(table *Table) []string {
	layouts := []string{ContributorLayoutPartitionKey}
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeRange {
			return []string{ContributorLayoutPartitionKey, ContributorLayoutFullKey}
		}
	}
	return layouts
}

// ContributorKeyString renders the tracked portion of a primary key for the
// given layout. The values are JSON-encoded, so they cannot contain the
// key separator used by the counter keys.
func ContributorKeyString(table *Table, key map[string]*AttributeValue, layout string) string {
	var names []string
	for _, ks := range table.KeySchema {
		if ks.KeyType == KeyTypeHash || layout == ContributorLayoutFullKey {
			names = append(names, ks.AttributeName)
		}
	}
	values := make([]string, 0, len(names))
	for _, n := range names {
		if v, ok := key[n]; ok && v != nil {
			values = append(values, encodeContributorValue(v))
		}
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// encodeContributorValue renders one attribute value as a stable string.
func encodeContributorValue(v *AttributeValue) string {
	switch {
	case v.S != nil:
		return "s:" + *v.S
	case v.N != nil:
		return "n:" + *v.N
	case v.B != nil:
		return "b:" + string(v.B)
	case v.BOOL != nil:
		return "bool:" + strconv.FormatBool(*v.BOOL)
	default:
		return "other"
	}
}

func contributorAccessKey(tableName, layout, keyStr string, minute int64) string {
	return tableName + keySep + layout + keySep + keyStr + keySep + fmt.Sprintf("%020d", minute)
}

// RecordAccessTxn credits the given number of access events to one
// counter key inside a storage transaction so the aggregation commits
// atomically with the item write that produced it. Counters are bucketed
// per minute so a report window can select them. Events targeting the
// same key must be aggregated before this call: the storage layer's
// read-modify-write reads the committed state, so it cannot see earlier
// writes of the same transaction and a second unaggregated call would
// overwrite the first.
func RecordAccessTxn(txn storage.Transaction, region, tableName, layout, keyStr string, at time.Time, events int64, unitsPerEvent float64) error {
	if keyStr == "" || events <= 0 {
		return nil
	}
	bucket := txn.Bucket(contributorBucketName(region))
	key := []byte(contributorAccessKey(tableName, layout, keyStr, at.Unix()/60))
	var acc contributorAccess
	if data, err := bucket.Get(key); err != nil {
		return fmt.Errorf("read contributor counter: %w", err)
	} else if len(data) > 0 {
		if err := json.Unmarshal(data, &acc); err != nil {
			return fmt.Errorf("unmarshal contributor counter: %w", err)
		}
	}
	acc.Count += events
	acc.Units += float64(events) * unitsPerEvent
	data, err := json.Marshal(acc)
	if err != nil {
		return fmt.Errorf("marshal contributor counter: %w", err)
	}
	return bucket.Put(key, data)
}

// TopKeys returns the most accessed tracked keys of a table within the
// half-open time window, ordered by consumed throughput units.
func (s *ContributorStore) TopKeys(tableName, layout string, start, end time.Time, limit int) ([]ContributorKeyStat, error) {
	if limit <= 0 {
		limit = 10
	}
	prefix := tableName + keySep + layout + keySep
	startMinute := start.Unix() / 60
	endMinute := end.Unix() / 60
	agg := make(map[string]*ContributorKeyStat)

	err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		parts := strings.Split(key, keySep)
		if len(parts) != 4 {
			return nil
		}
		minute, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return nil
		}
		if minute < startMinute || minute >= endMinute {
			return nil
		}
		var acc contributorAccess
		if err := json.Unmarshal(value, &acc); err != nil {
			return nil
		}
		stat, ok := agg[parts[2]]
		if !ok {
			stat = &ContributorKeyStat{Key: parts[2]}
			agg[parts[2]] = stat
		}
		stat.Count += acc.Count
		stat.Units += acc.Units
		return nil
	})
	if err != nil {
		return nil, err
	}

	stats := make([]ContributorKeyStat, 0, len(agg))
	for _, stat := range agg {
		stats = append(stats, *stat)
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Units != stats[j].Units {
			return stats[i].Units > stats[j].Units
		}
		return stats[i].Key < stats[j].Key
	})
	if len(stats) > limit {
		stats = stats[:limit]
	}
	return stats, nil
}

// SweepTableOlderThan removes the access counters of one table that fell
// out of the retention window. Keys are collected before any delete so the
// prefix scan never mutates while iterating.
func (s *ContributorStore) SweepTableOlderThan(tableName string, cutoff time.Time) error {
	prefix := tableName + keySep
	cutoffMinute := cutoff.Unix() / 60
	var doomed []string
	err := s.BaseStore.ScanPrefix(prefix, func(key string, _ []byte) error {
		parts := strings.Split(key, keySep)
		if len(parts) != 4 {
			return nil
		}
		minute, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil || minute >= cutoffMinute {
			return nil
		}
		doomed = append(doomed, key)
		return nil
	})
	if err != nil {
		return err
	}
	for _, key := range doomed {
		if err := s.BaseStore.Delete(key); err != nil {
			return fmt.Errorf("sweep contributor counter: %w", err)
		}
	}
	return nil
}

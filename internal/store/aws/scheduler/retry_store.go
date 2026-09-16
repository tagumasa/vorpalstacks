package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// RetryStore manages the persistence of RetryRecords in Pebble so that
// pending retries survive server restarts.
//
// Key format: "<zero-padded 20-digit UnixNano>:<record.ID>"
// Pebble stores keys in lexicographic byte order, so zero-padded numeric
// timestamps sort chronologically. GetDueRetryRecords uses ScanRange to
// read only records whose NextAttemptAt is at or before the cutoff,
// avoiding a full-table scan on every tick cycle.
type RetryStore struct {
	store *common.BaseStore
}

// NewRetryStore creates a RetryStore backed by a bucket in the given storage.
func NewRetryStore(baseStorage storage.BasicStorage, region string) *RetryStore {
	bucket := baseStorage.Bucket("scheduler-retries-" + region)
	return &RetryStore{
		store: common.NewBaseStore(bucket, "scheduler-retries"),
	}
}

// retryStorageKey builds the Pebble key for a RetryRecord using its
// NextAttemptAt timestamp (zero-padded to 20 digits for correct sorting)
// followed by the record ID for uniqueness.
func retryStorageKey(nextAttemptAt time.Time, recordID string) string {
	return fmt.Sprintf("%020d:%s", nextAttemptAt.UnixNano(), recordID)
}

// SaveRetryRecord persists a RetryRecord. The storage key is derived from
// NextAttemptAt so that records sort chronologically in Pebble.
func (rs *RetryStore) SaveRetryRecord(record *RetryRecord) error {
	key := retryStorageKey(record.NextAttemptAt, record.ID)
	return rs.store.Put(key, record)
}

// DeleteRetryRecord removes a RetryRecord. The caller must supply the
// NextAttemptAt that was used when the record was last saved, because the
// storage key is derived from it.
func (rs *RetryStore) DeleteRetryRecord(recordID string, nextAttemptAt time.Time) error {
	key := retryStorageKey(nextAttemptAt, recordID)
	return rs.store.Delete(key)
}

// GetDueRetryRecords returns all RetryRecords whose NextAttemptAt is at or
// before the given cutoff time. Uses ScanRange to read only the relevant
// prefix of the keyspace instead of iterating every record. A record whose
// stored value fails deserialisation is reaped: it can never be processed,
// and its key sorts at or before every future cutoff, so retaining it would
// rescan the corruption on every engine tick.
func (rs *RetryStore) GetDueRetryRecords(cutoff time.Time) ([]*RetryRecord, error) {
	// Include records whose timestamp equals the cutoff by adding 1ns
	// to the end key (ScanRange is exclusive on the end boundary).
	endKey := fmt.Sprintf("%020d", cutoff.Add(time.Nanosecond).UnixNano())
	iter := rs.store.Bucket().ScanRange(nil, []byte(endKey))
	defer iter.Close()

	var due []*RetryRecord
	var corruptKeys []string
	for iter.Next() {
		var record RetryRecord
		if err := json.Unmarshal(iter.Value(), &record); err != nil {
			// string(iter.Key()) copies the bytes, so the collected key
			// survives the iterator reusing its buffer. Deletion happens
			// after the scan, never while it is in progress.
			corruptKeys = append(corruptKeys, string(iter.Key()))
			continue // skip corrupt records
		}
		due = append(due, &record)
	}
	iterErr := iter.Error()
	for _, key := range corruptKeys {
		if err := rs.store.Delete(key); err != nil {
			// A failed reap stays visible to the next tick's scan, which
			// retries it; suppressing the healthy due records gathered
			// above over one undeletable key would be the worse failure.
			logs.Error("Failed to reap corrupt retry record",
				logs.String("key", key),
				logs.Err(err))
			continue
		}
		logs.Warn("Reaped corrupt retry record", logs.String("key", key))
	}
	return due, iterErr
}

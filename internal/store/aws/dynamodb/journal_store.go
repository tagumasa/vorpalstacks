package dynamodb

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

func journalBucketName(region string) string {
	return "dynamodb_journal-" + region
}

// Journal operations: every journaled item mutation is either a put (the key
// may or may not have existed before) or a delete.
const (
	JournalOperationPut    = "PUT"
	JournalOperationDelete = "DELETE"
)

// journalSequence breaks ordering ties between records appended within the
// same nanosecond so journal keys remain unique and totally ordered.
var journalSequence atomic.Uint64

// journalTimeFormat renders the journal key time component with fixed width
// so lexicographic bucket order matches chronological order.
const (
	journalTimeWidth = 20
	journalSeqWidth  = 10
)

// journalRecord is the persisted form of one item mutation on a table with
// point-in-time recovery enabled. BeforeImage holds the complete attribute
// map of the item as it was before the mutation (nil when the key did not
// exist), which is exactly the state needed to undo the change.
type journalRecord struct {
	Timestamp   int64                      `json:"timestamp"`
	Operation   string                     `json:"operation"`
	Key         map[string]*AttributeValue `json:"key"`
	BeforeImage map[string]*AttributeValue `json:"before_image,omitempty"`
}

// journalRecordKey builds the bucket key for a record: table, the append
// time, and the tie-breaking sequence, all ordered so a prefix scan of the
// table yields records oldest-first.
func journalRecordKey(tableName string, at time.Time) string {
	return tableName + keySep + fmt.Sprintf("%0*d%0*d",
		journalTimeWidth, at.UnixNano(), journalSeqWidth, journalSequence.Add(1))
}

// appendJournalTxnAt appends one journal record inside the given transaction
// so the journal entry commits atomically with the item mutation it
// describes. The append time is injected for testability.
func appendJournalTxnAt(txn storage.Transaction, region, tableName, operation string, key, beforeImage map[string]*AttributeValue, at time.Time) error {
	record := journalRecord{
		Timestamp:   at.UnixNano(),
		Operation:   operation,
		Key:         key,
		BeforeImage: beforeImage,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal journal record for table %s: %w", tableName, err)
	}
	bucket := txn.Bucket(journalBucketName(region))
	return bucket.Put([]byte(journalRecordKey(tableName, at)), data)
}

// appendJournalTxn appends one journal record at the current time.
func appendJournalTxn(txn storage.Transaction, region, tableName, operation string, key, beforeImage map[string]*AttributeValue) error {
	return appendJournalTxnAt(txn, region, tableName, operation, key, beforeImage, time.Now())
}

// JournalStore reads and prunes the item-mutation journal that backs
// point-in-time recovery.
type JournalStore struct {
	*common.BaseStore
	region string
}

// NewJournalStore creates a JournalStore for the given region.
func NewJournalStore(store storage.BasicStorage, region string) *JournalStore {
	return &JournalStore{
		BaseStore: common.NewBaseStore(store.Bucket(journalBucketName(region)), "dynamodb_journal"),
		region:    region,
	}
}

// journalRecordFromJSON decodes a persisted journal record.
func journalRecordFromJSON(data []byte) (*journalRecord, error) {
	var record journalRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("unmarshal journal record: %w", err)
	}
	return &record, nil
}

// ReverseReplay hands the caller every journaled mutation of the table that
// is newer than the given time, newest first. Replaying a record means
// restoring its BeforeImage (removing the key when the image is nil), which
// reconstructs the table state at the given time when applied over the
// current state in that order.
func (s *JournalStore) ReverseReplay(tableName string, from time.Time, fn func(record *JournalChange) error) error {
	prefix := tableName + keySep
	fromNanos := from.UnixNano()

	var records []*journalRecord
	if err := s.BaseStore.ScanPrefix(prefix, func(_ string, value []byte) error {
		record, err := journalRecordFromJSON(value)
		if err != nil {
			return err
		}
		if record.Timestamp > fromNanos {
			records = append(records, record)
		}
		return nil
	}); err != nil {
		return err
	}

	for i := len(records) - 1; i >= 0; i-- {
		record := records[i]
		change := &JournalChange{
			Timestamp:   time.Unix(0, record.Timestamp),
			Operation:   record.Operation,
			Key:         record.Key,
			BeforeImage: record.BeforeImage,
		}
		if err := fn(change); err != nil {
			return err
		}
	}
	return nil
}

// DeleteOlderThan removes every journal record of the table appended at or
// before the cutoff and returns how many were removed. Records at or before
// the table's EarliestRestorableDateTime can never be replayed by a
// restore, so pruning them keeps the journal bounded.
func (s *JournalStore) DeleteOlderThan(tableName string, cutoff time.Time) (int, error) {
	prefix := tableName + keySep
	cutoffNanos := cutoff.UnixNano()

	var stale []string
	if err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		record, err := journalRecordFromJSON(value)
		if err != nil {
			return err
		}
		if record.Timestamp <= cutoffNanos {
			stale = append(stale, key)
		}
		return nil
	}); err != nil {
		return 0, err
	}

	for _, key := range stale {
		if err := s.BaseStore.Delete(key); err != nil {
			return len(stale), err
		}
	}
	return len(stale), nil
}

// DeleteAllForTable removes the whole journal of the table. Disabling
// point-in-time recovery invalidates the journal because re-enabling starts
// a new restorable window at the re-enable time.
func (s *JournalStore) DeleteAllForTable(tableName string) error {
	return s.BaseStore.DeleteByPrefix(tableName + keySep)
}

// JournalChange is the caller-facing form of one journaled mutation.
type JournalChange struct {
	Timestamp   time.Time
	Operation   string
	Key         map[string]*AttributeValue
	BeforeImage map[string]*AttributeValue
}

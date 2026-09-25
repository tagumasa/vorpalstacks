package dynamodb

import (
	"fmt"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
	"vorpalstacks/internal/store/aws/common"
)

func journalBucketName(region string) string {
	return "dynamodb_journal-" + region
}

// JournalOperation distinguishes the two journaled item mutations: every
// record is either a put (the key may or may not have existed before) or a
// delete.
type JournalOperation string

const (
	JournalOperationPut    JournalOperation = "PUT"
	JournalOperationDelete JournalOperation = "DELETE"
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

// journalSeqModulus bounds the tie-breaking sequence at its key width: the
// rendered component must stay exactly journalSeqWidth digits, because a
// wider component sorts lexicographically BEFORE every narrower one and
// would invert the journal's total order. The tie-breaker needs uniqueness
// only among records appended within the same nanosecond, and 10^10
// same-nanosecond appends cannot occur, so wrapping the counter loses
// nothing.
const journalSeqModulus = 10_000_000_000

// journalRecord is the in-memory form of one item mutation on a table with
// point-in-time recovery enabled. BeforeImage holds the complete attribute
// map of the item as it was before the mutation (nil when the key did not
// exist), which is exactly the state needed to undo the change.
type journalRecord struct {
	Timestamp   int64
	Operation   JournalOperation
	Key         map[string]*AttributeValue
	BeforeImage map[string]*AttributeValue
}

// journalRecordKey builds the bucket key for a record: table, the append
// time, and the tie-breaking sequence, all ordered so a prefix scan of the
// table yields records oldest-first.
func journalRecordKey(tableName string, at time.Time) string {
	return tableName + KeySep + fmt.Sprintf("%0*d%0*d",
		journalTimeWidth, at.UnixNano(), journalSeqWidth, journalSequence.Add(1)%journalSeqModulus)
}

// appendJournalTxnAt appends one journal record inside the given transaction
// so the journal entry commits atomically with the item mutation it
// describes. The append time is injected for testability.
func appendJournalTxnAt(txn storage.Transaction, region, tableName string, operation JournalOperation, key, beforeImage map[string]*AttributeValue, at time.Time) error {
	data, err := proto.Marshal(journalRecordToProto(&journalRecord{
		Timestamp:   at.UnixNano(),
		Operation:   operation,
		Key:         key,
		BeforeImage: beforeImage,
	}))
	if err != nil {
		return fmt.Errorf("marshal journal record for table %s: %w", tableName, err)
	}
	bucket := txn.Bucket(journalBucketName(region))
	return bucket.Put([]byte(journalRecordKey(tableName, at)), data)
}

// appendJournalTxn appends one journal record at the current time.
func appendJournalTxn(txn storage.Transaction, region, tableName string, operation JournalOperation, key, beforeImage map[string]*AttributeValue) error {
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

// journalRecordFromBytes decodes a persisted journal record.
func journalRecordFromBytes(data []byte) (*journalRecord, error) {
	var pbRecord pb.JournalRecord
	if err := proto.Unmarshal(data, &pbRecord); err != nil {
		return nil, fmt.Errorf("unmarshal journal record: %w", err)
	}
	return protoToJournalRecord(&pbRecord), nil
}

// ReverseReplay hands the caller every journaled mutation of the table that
// is newer than the given time, newest first. Replaying a record means
// restoring its BeforeImage (removing the key when the image is nil), which
// reconstructs the table state at the given time when applied over the
// current state in that order.
func (s *JournalStore) ReverseReplay(tableName string, from time.Time, fn func(record *JournalChange) error) error {
	prefix := tableName + KeySep
	fromNanos := from.UnixNano()

	var records []*journalRecord
	if err := s.BaseStore.ScanPrefix(prefix, func(_ string, value []byte) error {
		record, err := journalRecordFromBytes(value)
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

// WindowReplay hands the caller every journaled mutation of the table whose
// append time lies in the half-open window [from, to), oldest first. The
// incremental export reads exactly this window: the specification documents
// ExportFromTime as the inclusive start of the exported change range and
// ExportToTime as its exclusive end.
func (s *JournalStore) WindowReplay(tableName string, from, to time.Time, fn func(record *JournalChange) error) error {
	prefix := tableName + KeySep
	fromNanos, toNanos := from.UnixNano(), to.UnixNano()
	return s.BaseStore.ScanPrefix(prefix, func(_ string, value []byte) error {
		record, err := journalRecordFromBytes(value)
		if err != nil {
			return err
		}
		if record.Timestamp < fromNanos || record.Timestamp >= toNanos {
			return nil
		}
		change := &JournalChange{
			Timestamp:   time.Unix(0, record.Timestamp),
			Operation:   record.Operation,
			Key:         record.Key,
			BeforeImage: record.BeforeImage,
		}
		return fn(change)
	})
}

// DeleteOlderThan removes every journal record of the table appended at or
// before the cutoff and returns how many were removed. Records at or before
// the table's EarliestRestorableDateTime can never be replayed by a
// restore, so pruning them keeps the journal bounded.
func (s *JournalStore) DeleteOlderThan(tableName string, cutoff time.Time) (int, error) {
	prefix := tableName + KeySep
	cutoffNanos := cutoff.UnixNano()

	var stale []string
	if err := s.BaseStore.ScanPrefix(prefix, func(key string, value []byte) error {
		record, err := journalRecordFromBytes(value)
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
	return s.BaseStore.DeleteByPrefix(tableName + KeySep)
}

// JournalChange is the caller-facing form of one journaled mutation.
type JournalChange struct {
	Timestamp   time.Time
	Operation   JournalOperation
	Key         map[string]*AttributeValue
	BeforeImage map[string]*AttributeValue
}

// journalPut appends a journal record for a put on a table with
// point-in-time recovery enabled. The record stores the pre-write item so a
// restore can undo the change; the append shares the caller's transaction,
// so the journal never diverges from the item state.
func (t *DynamoDBTxn) journalPut(table *Table, key map[string]*AttributeValue) error {
	if !pitrEnabled(table) {
		return nil
	}
	beforeImage := t.itemBeforeImage(table.Name, key)
	return appendJournalTxn(t.txn, t.region(), table.Name, JournalOperationPut, key, beforeImage)
}

// journalDelete appends a journal record for a delete. Deletes of keys that
// do not exist are not journaled because they change nothing.
func (t *DynamoDBTxn) journalDelete(table *Table, key map[string]*AttributeValue) error {
	if !pitrEnabled(table) {
		return nil
	}
	beforeImage := t.itemBeforeImage(table.Name, key)
	if beforeImage == nil {
		return nil
	}
	return appendJournalTxn(t.txn, t.region(), table.Name, JournalOperationDelete, key, beforeImage)
}

// pitrEnabled reports whether the table currently has point-in-time
// recovery enabled.
func pitrEnabled(table *Table) bool {
	return table != nil && table.PointInTimeRecovery != nil && table.PointInTimeRecovery.Status == PITRStatusEnabled
}

// itemBeforeImage reads the item's current attribute map through the
// transaction, or nil when the key does not exist.
func (t *DynamoDBTxn) itemBeforeImage(tableName string, key map[string]*AttributeValue) map[string]*AttributeValue {
	existing, err := t.GetItem(tableName, key)
	if err != nil || existing == nil {
		return nil
	}
	return existing.Attributes
}

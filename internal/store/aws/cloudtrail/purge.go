package cloudtrail

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_cloudtrail"

	"google.golang.org/protobuf/proto"
)

// EventHistoryRetention is how long event-history entries are kept before
// the retention worker purges them. AWS CloudTrail event history covers the
// past 90 days of management events per region; no AWS-plane operation
// deletes event history, so expiry is enforced platform-side.
const EventHistoryRetention = 90 * 24 * time.Hour

// defaultPurgeBatchSize is the number of events removed per storage
// transaction when PurgeEventHistoryBefore is called without an explicit
// batch size. Small transactions keep the purge off the recorder's hot
// path; the deleted (old) and appended (new) key ranges never overlap.
const defaultPurgeBatchSize = 256

// PurgeResult reports the per-bucket deletion counts of one purge run. A
// second run over the same cutoff removes nothing (every counter zero).
type PurgeResult struct {
	Events           int
	EventIDIndex     int
	TimeIndex        int
	EventNameIndex   int
	UsernameIndex    int
	EventSourceIndex int
	Batches          int
}

// PurgeEventHistoryBefore deletes every event-history entry whose event
// time is before cutoff, together with the five companion entries each
// write creates (the eventID index and the four ct_idx_* lookup indexes).
// Events are removed in batches of batchSize storage transactions; a
// batchSize <= 0 selects the default. The operation is idempotent.
func (s *CloudTrailStore) PurgeEventHistoryBefore(cutoff time.Time, batchSize int) (*PurgeResult, error) {
	if batchSize <= 0 {
		batchSize = defaultPurgeBatchSize
	}
	result := &PurgeResult{}

	// Event record keys are "<unixnano>#<eventID>". Every nanosecond
	// timestamp in the int64 range has 19 digits until the year 2262, so
	// lexicographic order equals numeric order; an exclusive upper bound
	// at "<cutoff>#" selects exactly the events strictly older than the
	// cutoff.
	bound := []byte(fmt.Sprintf("%d#", cutoff.UnixNano()))

	var resumeAfter string
	for {
		batch, lastKey, err := s.collectPurgeBatch(bound, resumeAfter, batchSize)
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			return result, nil
		}
		if err := s.deletePurgeBatch(batch, result); err != nil {
			return nil, err
		}
		result.Batches++
		resumeAfter = lastKey
	}
}

// purgeEntry holds one scanned event: the exact storage key of the record
// plus the event with the index attributes the lookup-index keys are built
// from.
type purgeEntry struct {
	eventKey string
	event    *Event
}

// collectPurgeBatch scans the events bucket from just after resumeAfter up
// to the exclusive bound and decodes up to batchSize records. It returns
// the raw key of the last scanned record so the next batch can resume
// strictly after it.
func (s *CloudTrailStore) collectPurgeBatch(bound []byte, resumeAfter string, batchSize int) ([]*purgeEntry, string, error) {
	var start []byte
	if resumeAfter != "" {
		// Appending 0x00 (the smallest byte value) yields the smallest
		// key strictly greater than resumeAfter in lexicographic order.
		start = append([]byte(resumeAfter), 0x00)
	}
	iter := s.eventsStore.Bucket().ScanRange(start, bound)
	defer iter.Close()

	batch := make([]*purgeEntry, 0, batchSize)
	lastKey := ""
	for iter.Next() {
		key := string(iter.Key())
		// The iterator value slice is only valid until the next Next
		// call, so the record bytes are copied before advancing.
		value := make([]byte, len(iter.Value()))
		copy(value, iter.Value())

		event, err := purgeEventFromRecord(key, value)
		if err != nil {
			return nil, "", err
		}
		batch = append(batch, &purgeEntry{eventKey: key, event: event})
		lastKey = key
		if len(batch) >= batchSize {
			break
		}
	}
	if err := iter.Error(); err != nil {
		return nil, "", err
	}
	return batch, lastKey, nil
}

// purgeEventFromRecord rebuilds the fields the index keys are derived
// from. The exact event time and ID come from the storage key: the record
// proto stores only millisecond precision while the index keys encode the
// full nanosecond timestamp, so the record value alone cannot reproduce
// them. The indexed attributes (event name, username, event source) come
// from the record.
func purgeEventFromRecord(key string, value []byte) (*Event, error) {
	nanos, eventID, ok := strings.Cut(key, "#")
	if !ok {
		return nil, fmt.Errorf("cloudtrail purge: malformed event key %q", key)
	}
	ts, err := strconv.ParseInt(nanos, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cloudtrail purge: malformed event key %q: %w", key, err)
	}
	var p pb.Event
	if err := proto.Unmarshal(value, &p); err != nil {
		return nil, fmt.Errorf("cloudtrail purge: corrupt event record %q: %w", key, err)
	}
	event := ProtoToEvent(&p)
	event.EventTime = time.Unix(0, ts).UTC()
	event.EventID = eventID
	return event, nil
}

// deletePurgeBatch removes one batch of events and their companion entries
// in a single transaction per batch. The count accumulation runs on a
// transaction-local result merged only after a successful commit, so a
// failed batch never inflates the reported counts.
func (s *CloudTrailStore) deletePurgeBatch(batch []*purgeEntry, result *PurgeResult) error {
	if s.storage != nil {
		return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			batchResult := &PurgeResult{}
			for _, entry := range batch {
				if err := s.deleteEventEntriesInTxn(txn, entry, batchResult); err != nil {
					return err
				}
			}
			mergePurgeResult(result, batchResult)
			return nil
		})
	}
	for _, entry := range batch {
		if err := s.deleteEventEntries(entry, result); err != nil {
			return err
		}
	}
	return nil
}

// deleteEventEntriesInTxn removes all six entries of one event inside a
// transaction, counting each deletion per bucket. Index keys are built by
// the shared key builders so the purge deletes exactly what the write path
// created.
func (s *CloudTrailStore) deleteEventEntriesInTxn(txn storage.Transaction, entry *purgeEntry, res *PurgeResult) error {
	if err := txn.Bucket(eventBucketName(s.region)).Delete([]byte(entry.eventKey)); err != nil {
		return err
	}
	res.Events++
	if err := txn.Bucket(eventIDIndexBucketName(s.region)).Delete([]byte(entry.event.EventID)); err != nil {
		return err
	}
	res.EventIDIndex++
	return s.deleteIndexEntries(entry.event, res, func(bucketName string, key []byte) error {
		return txn.Bucket(bucketName).Delete(key)
	})
}

// deleteEventEntries is the non-transactional fallback used when the
// backing storage provides no transactions; it mirrors the write path's
// own fallback structure.
func (s *CloudTrailStore) deleteEventEntries(entry *purgeEntry, res *PurgeResult) error {
	if err := s.eventsStore.Bucket().Delete([]byte(entry.eventKey)); err != nil {
		return err
	}
	res.Events++
	if err := s.eventIDIndexStore.Bucket().Delete([]byte(entry.event.EventID)); err != nil {
		return err
	}
	res.EventIDIndex++
	return s.deleteIndexEntries(entry.event, res, func(bucketName string, key []byte) error {
		return s.indexer.storage.Bucket(bucketName).Delete(key)
	})
}

// deleteIndexEntries removes the ct_idx_* lookup-index entries of one
// event through the supplied delete function and counts them per index.
// The index keys are built with this store's account and region — the same
// values the write path's indexer used.
func (s *CloudTrailStore) deleteIndexEntries(event *Event, res *PurgeResult, deleteFn func(bucketName string, key []byte) error) error {
	for _, key := range buildIndexKeys(s.accountID, s.region, event) {
		if err := deleteFn(key.EncodePrefix(), indexStorageKey(key)); err != nil {
			return err
		}
		switch key.IndexType {
		case IndexByTime:
			res.TimeIndex++
		case IndexByEventName:
			res.EventNameIndex++
		case IndexByUsername:
			res.UsernameIndex++
		case IndexByEventSource:
			res.EventSourceIndex++
		}
	}
	return nil
}

func mergePurgeResult(dst, src *PurgeResult) {
	dst.Events += src.Events
	dst.EventIDIndex += src.EventIDIndex
	dst.TimeIndex += src.TimeIndex
	dst.EventNameIndex += src.EventNameIndex
	dst.UsernameIndex += src.UsernameIndex
	dst.EventSourceIndex += src.EventSourceIndex
}

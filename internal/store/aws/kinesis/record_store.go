package kinesis

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	"google.golang.org/protobuf/proto"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

const (
	iteratorCleanupInterval = 5 * time.Minute
	retentionTrimInterval   = 5 * time.Minute
)

// appendRecord is the single record-write primitive every ingestion path
// (single put, shard-selected put, batch) goes through: it allocates the
// sequence number from ONE clock read — the same stamp keys the sequence
// number and rides the record as its arrival time, so the retention trim
// (which compares the key's embedded time) and the read path's visibility
// filter (which compares the arrival time) can never disagree at the
// boundary — advances the shard's LatestSequenceNumber cursor, and stages
// the record write on the caller's transaction; persisting the shard
// cursor is the caller's side of the same transaction, so a record is
// never readable while its cursor write is still pending (or failed).
// Record ingestion is not a stream-configuration change — it never
// rewrites the stream record.
func (s *KinesisStore) appendRecord(txn storage.Transaction, streamName string, shard *Shard, partitionKey, data string, now time.Time) (*Record, error) {
	seqNum := s.generateSequenceNumberAt(shard.ShardID, now)
	shard.LatestSequenceNumber = seqNum

	record := &Record{
		SequenceNumber:              seqNum,
		ApproximateArrivalTimestamp: now,
		Data:                        data,
		PartitionKey:                partitionKey,
	}

	key := fmt.Sprintf("%s#%s#%s", streamName, shard.ShardID, seqNum)
	recordBytes, err := proto.Marshal(RecordToProto(record))
	if err != nil {
		return nil, err
	}
	if err := txn.Bucket(s.recordsBucketName()).Put([]byte(key), recordBytes); err != nil {
		return nil, err
	}
	return record, nil
}

// activeShardsOf filters a stream's shards down to the open ones — the
// routing space; a stream whose shards are all closed yields none and the
// callers answer ErrNoActiveShards.
func activeShardsOf(shards []*Shard) []*Shard {
	var activeShards []*Shard
	for _, shard := range shards {
		if shard.SequenceNumberRange == nil || shard.SequenceNumberRange.EndingSequenceNumber == "" {
			activeShards = append(activeShards, shard)
		}
	}
	return activeShards
}

// routeRecord selects the shard one record writes to: the explicit hash key
// when present, else the partition key's hash — the same placement rule on
// every write path. Routing failures are server states; client-side key
// validity is settled before the store, and the inner error names the hash
// origin.
func (s *KinesisStore) routeRecord(activeShards []*Shard, explicitHashKey, partitionKey string) (*Shard, error) {
	if explicitHashKey != "" {
		return s.selectShardByHashKey(activeShards, explicitHashKey)
	}
	return s.selectShardByPartitionKey(activeShards, partitionKey)
}

// PutRecordWithShardSelection selects the appropriate shard by partition key
// or explicit hash key and writes a single record, all under a single lock:
// the record and its shard-cursor advance commit as one transaction, so a
// failure leaves neither a readable record without its cursor nor a
// reported failure whose record actually landed (a client retry would
// duplicate it).
func (s *KinesisStore) PutRecordWithShardSelection(streamName, partitionKey, data, explicitHashKey string) (*Record, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	shards, err := s.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, "", err
	}

	activeShards := activeShardsOf(shards)
	if len(activeShards) == 0 {
		return nil, "", ErrNoActiveShards
	}

	// Routing failures are server states — client-side key validity is
	// settled before the store — and surface as errors (mapping to
	// InternalFailure) rather than a silent fallback shard the client would
	// be told a wrong ShardId for.
	targetShard, err := s.routeRecord(activeShards, explicitHashKey, partitionKey)
	if err != nil {
		return nil, "", fmt.Errorf("route record: %w", err)
	}

	now := time.Now().UTC()
	var record *Record
	err = s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		var aerr error
		record, aerr = s.appendRecord(txn, streamName, targetShard, partitionKey, data, now)
		if aerr != nil {
			return aerr
		}
		return s.putShardInTxn(txn, targetShard)
	})
	if err != nil {
		return nil, "", err
	}

	s.maybeTrimRetentionLocked(streamName)

	return record, targetShard.ShardID, nil
}

// PutRecords writes multiple data records to a Kinesis stream: each entry
// is routed, and every record write plus every touched shard's cursor
// advance commits as ONE transaction — a storage failure fails the whole
// call instead of leaving readable records behind a reported failure (a
// client retry would duplicate them). A routing failure is a server state
// that touches no storage: that entry alone reports InternalFailure, the
// model's per-entry write-time channel.
func (s *KinesisStore) PutRecords(streamName string, records []PutRecordRequest) ([]PutRecordResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	shards, err := s.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, err
	}

	activeShards := activeShardsOf(shards)
	if len(activeShards) == 0 {
		return nil, ErrNoActiveShards
	}

	results := make([]PutRecordResult, len(records))
	entryShard := make([]*Shard, len(records))
	touched := make(map[*Shard]bool)
	for i, req := range records {
		// A routing failure is a server state (client-side key validity is
		// settled before the store), so the entry reports InternalFailure —
		// never a throughput identity, and never a fallback shard. The
		// model fixes the entry's message for this code.
		shard, err := s.routeRecord(activeShards, req.ExplicitHashKey, req.PartitionKey)
		if err != nil {
			results[i] = PutRecordResult{
				ErrorCode:    "InternalFailure",
				ErrorMessage: InternalFailureEntryMessage,
			}
			continue
		}
		entryShard[i] = shard
		touched[shard] = true
	}

	now := time.Now().UTC()
	err = s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		for i, req := range records {
			if entryShard[i] == nil {
				continue
			}
			record, aerr := s.appendRecord(txn, streamName, entryShard[i], req.PartitionKey, req.Data, now)
			if aerr != nil {
				return aerr
			}
			results[i] = PutRecordResult{
				SequenceNumber: record.SequenceNumber,
				ShardID:        entryShard[i].ShardID,
			}
		}
		for shard := range touched {
			if perr := s.putShardInTxn(txn, shard); perr != nil {
				return perr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.maybeTrimRetentionLocked(streamName)

	return results, nil
}

func (s *KinesisStore) selectShardByPartitionKey(shards []*Shard, partitionKey string) (*Shard, error) {
	if len(shards) == 0 {
		return nil, nil
	}

	return s.shardForHash(shards, s.hashPartitionKey(partitionKey), "partition key hash")
}

func (s *KinesisStore) hashPartitionKey(partitionKey string) *big.Int {
	h := md5.Sum([]byte(partitionKey))
	// MD5's 128-bit digest already spans the full shard hash-key space
	// [0, 2^128-1] — no reduction is needed.
	return new(big.Int).SetBytes(h[:])
}

// selectShardByHashKey selects a shard using an explicit hash key value
// (bypassing MD5 partition key hashing).
func (s *KinesisStore) selectShardByHashKey(shards []*Shard, hashKey string) (*Shard, error) {
	if len(shards) == 0 {
		return nil, nil
	}

	hash, ok := new(big.Int).SetString(hashKey, 10)
	if !ok {
		return nil, fmt.Errorf("invalid explicit hash key: %s", hashKey)
	}
	return s.shardForHash(shards, hash, "explicit hash key")
}

// shardForHash locates the shard whose hash-key range contains hash. The
// ranges tile [0, MaxShardHashKey] contiguously, so exactly one shard holds
// any in-space value; a miss means the tiling itself is defective. The kind
// string names the hash's origin in the no-shard error.
func (s *KinesisStore) shardForHash(shards []*Shard, hash *big.Int, kind string) (*Shard, error) {
	for _, shard := range shards {
		// A shard record without its hash-key submessage reads as
		// empty-string bounds — the absent-submessage reading the
		// family's guards and formatters take — so the parses below
		// report an unparseable bound instead of panicking on the
		// absent record.
		startBound, endBound := "", ""
		if shard.HashKeyRange != nil {
			startBound, endBound = shard.HashKeyRange.StartingHashKey, shard.HashKeyRange.EndingHashKey
		}
		start, ok := new(big.Int).SetString(startBound, 10)
		if !ok {
			return nil, fmt.Errorf("invalid starting hash key: %s", startBound)
		}
		end, ok := new(big.Int).SetString(endBound, 10)
		if !ok {
			return nil, fmt.Errorf("invalid ending hash key: %s", endBound)
		}
		if hash.Cmp(start) >= 0 && hash.Cmp(end) <= 0 {
			return shard, nil
		}
	}
	return nil, fmt.Errorf("no shard found for %s %s", kind, hash.String())
}

// GetRecords returns up to limit records at or after the given position in
// the shard, skipping records the retention window has aged out, plus the
// advance position: the sequence number of the last record the scan
// consumed (returned or aged out), or the starting position when the scan
// consumed none. Aged-out records never consume page slots — the page
// fills only with in-retention records — and an aged-out record the scan
// passes while the page still has room advances the position, so a reader
// whose whole page has expired still moves forward to the tip instead of
// re-reading the same expired page forever; past a full page the scan
// stops advancing, because the records there belong to the next read. A
// zero cutoff disables the filter — the AT_TIMESTAMP
// position scan reads raw shard history.
func (s *KinesisStore) GetRecords(streamName, shardID, startingSeqNum string, limit int32, includeStart bool, retentionCutoff time.Time) ([]*Record, string, error) {
	prefix := fmt.Sprintf("%s#%s#", streamName, shardID)
	var page []*Record
	advance := startingSeqNum

	err := s.recordsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbRecord pb.Record
		if err := proto.Unmarshal(value, &pbRecord); err != nil {
			return err
		}
		record := ProtoToRecord(&pbRecord)

		atOrAfter := startingSeqNum == "" ||
			(includeStart && record.SequenceNumber >= startingSeqNum) ||
			(!includeStart && record.SequenceNumber > startingSeqNum)
		if !atOrAfter || int32(len(page)) >= limit {
			return nil
		}
		// Every record the reader passes advances it, whether or not the
		// retention window still admits the record into the page.
		advance = record.SequenceNumber
		if !retentionCutoff.IsZero() && record.ApproximateArrivalTimestamp.Before(retentionCutoff) {
			return nil
		}
		page = append(page, record)
		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return page, advance, nil
}

// TrimExpiredRecords physically deletes the stream's records whose age
// exceeds the stream's retention window. Visibility is governed by the read
// path's age filter; this pass reclaims the storage. The caller-facing form
// takes the mutex; the write paths reach the locked variant through
// maybeTrimRetentionLocked.
func (s *KinesisStore) TrimExpiredRecords(streamName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.trimExpiredRecordsLocked(streamName)
}

// trimExpiredRecordsLocked reclaims the aged-out records of one stream; the
// caller must hold s.mu.
func (s *KinesisStore) trimExpiredRecordsLocked(streamName string) error {
	stream, err := s.GetStream(streamName)
	if err != nil {
		return err
	}
	cutoffSeq := formatSequenceNumber(RetentionCutoff(stream.RetentionPeriodHours).UnixNano(), 0, 0)

	shards, err := s.ListShards(streamName, nil, "", 0)
	if err != nil {
		return err
	}
	for _, shard := range shards {
		prefix := fmt.Sprintf("%s#%s#", streamName, shard.ShardID)
		if err := s.recordsStore.ScanPrefix(prefix, func(key string, _ []byte) error {
			// Record keys carry the arrival-time nanoseconds as the leading
			// sequence-number digits, so the age comparison runs on the key
			// alone. While every live timestamp shares one digit width, key
			// order is chronological; a width change (year 2262) can only
			// leave an old record untrimmed for a further interval, never
			// trim a newer record — the read-side age filter stays the
			// authority on visibility regardless.
			if key[len(prefix):] >= cutoffSeq {
				return nil
			}
			return s.recordsStore.Delete(key)
		}); err != nil {
			return fmt.Errorf("failed to trim expired records for shard %s: %w", shard.ShardID, err)
		}
	}
	return nil
}

// maybeTrimRetentionLocked runs the physical retention trim at most once
// per interval per store, from the record-write paths — the plane whose
// writes create the reclaim obligation; the caller must hold s.mu. The
// cadence is a platform mechanism AWS leaves undocumented; visibility of
// aged records never depends on it.
func (s *KinesisStore) maybeTrimRetentionLocked(streamName string) {
	now := time.Now().UTC()
	if now.Before(s.nextRetentionTrim) {
		return
	}
	s.nextRetentionTrim = now.Add(retentionTrimInterval)
	_ = s.trimExpiredRecordsLocked(streamName)
}

// PutAgedRecord writes one record whose arrival timestamp — and the
// sequence-number key encoding it — sits the given age in the past.
// Production ingestion stamps arrival at write time through the
// appendRecord paths; this entry point exists so retention behaviour can
// be exercised against shard history at a chosen depth without waiting
// out a real retention window. It takes the store mutex like every other
// record write — an unlocked write could interleave with a retention trim.
func (s *KinesisStore) PutAgedRecord(streamName, shardID string, age time.Duration, counter int64) (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	arrival := time.Now().UTC().Add(-age)
	record := &Record{
		SequenceNumber:              formatSequenceNumber(arrival.UnixNano(), counter, 0),
		ApproximateArrivalTimestamp: arrival,
		Data:                        "aged-record",
		PartitionKey:                "pk",
	}
	key := fmt.Sprintf("%s#%s#%s", streamName, shardID, record.SequenceNumber)
	if err := s.recordsStore.PutProto(key, RecordToProto(record)); err != nil {
		return nil, err
	}
	return record, nil
}

// cleanExpiredIterators removes shard iterators that have passed their expiry time.
// Runs at most once per iteratorCleanupInterval to avoid O(n) scans on every CreateShardIterator call.
func (s *KinesisStore) cleanExpiredIterators() {
	now := time.Now().UTC()
	if now.Before(s.nextIteratorCleanup) {
		return
	}
	s.nextIteratorCleanup = now.Add(iteratorCleanupInterval)
	_ = s.iteratorsStore.ForEach(func(key string, value []byte) error {
		var it ShardIterator
		if err := json.Unmarshal(value, &it); err != nil {
			return nil
		}
		if now.After(it.ExpiresAt) {
			_ = s.iteratorsStore.Delete(key)
		}
		return nil
	})
}

// CreateShardIterator creates a shard iterator for reading from a Kinesis stream.
// The whole creation — including the expiry-sweep bookkeeping on the
// store's cleanup field — runs under the store mutex.
func (s *KinesisStore) CreateShardIterator(streamName, shardID, iteratorType, seqNum string, timestamp *time.Time) (*ShardIterator, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanExpiredIterators()

	shard, err := s.GetShard(streamName, shardID)
	if err != nil {
		return nil, err
	}
	if shard.SequenceNumberRange == nil {
		return nil, fmt.Errorf("shard %s has no sequence number range", shardID)
	}

	var startingSeqNum string
	switch iteratorType {
	case "TRIM_HORIZON":
		startingSeqNum = shard.SequenceNumberRange.StartingSequenceNumber
	case "LATEST":
		if shard.LatestSequenceNumber != "" {
			startingSeqNum = shard.LatestSequenceNumber
		} else if shard.SequenceNumberRange.EndingSequenceNumber != "" {
			startingSeqNum = shard.SequenceNumberRange.EndingSequenceNumber
		} else {
			startingSeqNum = shard.SequenceNumberRange.StartingSequenceNumber
		}
	case "AT_SEQUENCE_NUMBER":
		startingSeqNum = seqNum
	case "AFTER_SEQUENCE_NUMBER":
		startingSeqNum = seqNum
	case "AT_TIMESTAMP":
		if timestamp != nil {
			records, _, err := s.GetRecords(streamName, shardID, shard.SequenceNumberRange.StartingSequenceNumber, ATTimestampScanBound, true, time.Time{})
			if err != nil {
				return nil, fmt.Errorf("failed to get records for AT_TIMESTAMP iterator: %w", err)
			}
			startingSeqNum = shard.SequenceNumberRange.StartingSequenceNumber
			for _, rec := range records {
				if !rec.ApproximateArrivalTimestamp.Before(*timestamp) {
					break
				}
				startingSeqNum = rec.SequenceNumber
			}
		}
	default:
		startingSeqNum = shard.SequenceNumberRange.StartingSequenceNumber
	}

	now := time.Now().UTC()
	iterator := &ShardIterator{
		IteratorID:     uuid.New().String(),
		StreamName:     streamName,
		ShardID:        shardID,
		IteratorType:   iteratorType,
		SequenceNumber: startingSeqNum,
		ExpiresAt:      now.Add(5 * time.Minute),
	}

	if err := s.iteratorsStore.Put(iterator.IteratorID, iterator); err != nil {
		return nil, err
	}

	return iterator, nil
}

// GetShardIterator retrieves a shard iterator by its ID.
func (s *KinesisStore) GetShardIterator(iteratorID string) (*ShardIterator, error) {
	var iterator ShardIterator
	if err := s.iteratorsStore.Get(iteratorID, &iterator); err != nil {
		return nil, ErrInvalidIterator
	}

	if time.Now().UTC().After(iterator.ExpiresAt) {
		_ = s.iteratorsStore.Delete(iteratorID)
		return nil, ErrExpiredIterator
	}

	return &iterator, nil
}

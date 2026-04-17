package kinesis

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sort"
	"time"

	"github.com/google/uuid"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

// PutRecord writes a single data record to a Kinesis stream shard.
func (s *KinesisStore) PutRecord(streamName, shardID, partitionKey, data string) (*Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	shard, err := s.GetShard(streamName, shardID)
	if err != nil {
		return nil, err
	}

	if shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber != "" {
		return nil, ErrShardClosed
	}

	seqNum := s.generateSequenceNumber(shard.ShardID)
	shard.LatestSequenceNumber = seqNum
	if err := s.UpdateShard(shard); err != nil {
		return nil, err
	}

	record := &Record{
		SequenceNumber:              seqNum,
		ApproximateArrivalTimestamp: time.Now().UTC(),
		Data:                        data,
		PartitionKey:                partitionKey,
	}

	key := fmt.Sprintf("%s#%s#%s", streamName, shardID, seqNum)
	if err := s.recordsStore.PutProto(key, RecordToProto(record)); err != nil {
		return nil, err
	}

	return record, nil
}

// PutRecords writes multiple data records to a Kinesis stream.
func (s *KinesisStore) PutRecords(streamName string, records []PutRecordRequest) ([]PutRecordResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, err := s.GetStream(streamName)
	if err != nil {
		return nil, err
	}

	shards, err := s.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, err
	}

	var activeShards []*Shard
	for _, shard := range shards {
		if shard.SequenceNumberRange == nil || shard.SequenceNumberRange.EndingSequenceNumber == "" {
			activeShards = append(activeShards, shard)
		}
	}

	if len(activeShards) == 0 {
		return nil, ErrNoActiveShards
	}

	results := make([]PutRecordResult, len(records))
	for i, req := range records {
		shard, err := s.selectShardByPartitionKey(activeShards, req.PartitionKey)
		if err != nil || shard == nil {
			results[i] = PutRecordResult{ErrorCode: "ProvisionedThroughputExceededException"}
			continue
		}

		seqNum := s.generateSequenceNumber(shard.ShardID)
		shard.LatestSequenceNumber = seqNum
		if err := s.UpdateShard(shard); err != nil {
			results[i] = PutRecordResult{ErrorCode: "InternalFailure"}
			continue
		}

		record := &Record{
			SequenceNumber:              seqNum,
			ApproximateArrivalTimestamp: time.Now().UTC(),
			Data:                        req.Data,
			PartitionKey:                req.PartitionKey,
		}

		key := fmt.Sprintf("%s#%s#%s", streamName, shard.ShardID, seqNum)
		if err := s.recordsStore.PutProto(key, RecordToProto(record)); err != nil {
			results[i] = PutRecordResult{ErrorCode: "InternalFailure"}
			continue
		}

		results[i] = PutRecordResult{
			SequenceNumber: seqNum,
			ShardID:        shard.ShardID,
		}
	}

	stream.LastModifiedAt = time.Now().UTC()
	if err := s.UpdateStream(stream); err != nil {
		return nil, err
	}

	return results, nil
}

// SelectShardByPartitionKey selects a shard for the given partition key.
func (s *KinesisStore) SelectShardByPartitionKey(shards []*Shard, partitionKey string) *Shard {
	shard, _ := s.selectShardByPartitionKey(shards, partitionKey)
	return shard
}

func (s *KinesisStore) selectShardByPartitionKey(shards []*Shard, partitionKey string) (*Shard, error) {
	if len(shards) == 0 {
		return nil, nil
	}

	hash := s.hashPartitionKey(partitionKey)
	for _, shard := range shards {
		start, ok := new(big.Int).SetString(shard.HashKeyRange.StartingHashKey, 10)
		if !ok {
			return nil, fmt.Errorf("invalid starting hash key: %s", shard.HashKeyRange.StartingHashKey)
		}
		end, ok := new(big.Int).SetString(shard.HashKeyRange.EndingHashKey, 10)
		if !ok {
			return nil, fmt.Errorf("invalid ending hash key: %s", shard.HashKeyRange.EndingHashKey)
		}
		if hash.Cmp(start) >= 0 && hash.Cmp(end) <= 0 {
			return shard, nil
		}
	}
	return nil, fmt.Errorf("no shard found for partition key hash %s", hash.String())
}

func (s *KinesisStore) hashPartitionKey(partitionKey string) *big.Int {
	h := md5.Sum([]byte(partitionKey))
	hashInt := new(big.Int).SetBytes(h[:])
	maxInt := new(big.Int).SetUint64(math.MaxUint64)
	return hashInt.Mod(hashInt, maxInt)
}

// GetRecords retrieves records from a Kinesis stream shard.
func (s *KinesisStore) GetRecords(streamName, shardID, startingSeqNum string, limit int32, includeStart bool) ([]*Record, string, error) {
	prefix := fmt.Sprintf("%s#%s#", streamName, shardID)
	var allRecords []*Record

	err := s.recordsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbRecord pb.Record
		if err := proto.Unmarshal(value, &pbRecord); err != nil {
			return err
		}
		record := ProtoToRecord(&pbRecord)

		include := startingSeqNum == "" ||
			(includeStart && record.SequenceNumber >= startingSeqNum) ||
			(!includeStart && record.SequenceNumber > startingSeqNum)
		if include {
			allRecords = append(allRecords, record)
		}
		return nil
	})

	if err != nil {
		return nil, "", err
	}

	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].SequenceNumber < allRecords[j].SequenceNumber
	})

	if int32(len(allRecords)) > limit {
		allRecords = allRecords[:limit]
	}

	var lastSeqNum string
	if len(allRecords) > 0 {
		lastSeqNum = allRecords[len(allRecords)-1].SequenceNumber
	}

	return allRecords, lastSeqNum, nil
}

// cleanExpiredIterators removes shard iterators that have passed their expiry time.
func (s *KinesisStore) cleanExpiredIterators() {
	now := time.Now().UTC()
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
func (s *KinesisStore) CreateShardIterator(streamName, shardID, iteratorType, seqNum string, timestamp *time.Time) (*ShardIterator, error) {
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
			records, _, err := s.GetRecords(streamName, shardID, shard.SequenceNumberRange.StartingSequenceNumber, 10000, true)
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
		Timestamp:      timestamp,
		CreatedAt:      now,
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

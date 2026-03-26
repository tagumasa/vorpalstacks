// Package kinesis provides AWS Kinesis data stream storage functionality for vorpalstacks.
package kinesis

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// KinesisStore provides data storage operations for AWS Kinesis streams.
type KinesisStore struct {
	*common.BaseStore
	shardsStore    *common.BaseStore
	recordsStore   *common.BaseStore
	consumersStore *common.BaseStore
	iteratorsStore *common.BaseStore
	*common.TagStore
	arnBuilder      *svcarn.ARNBuilder
	accountID       string
	region          string
	mu              sync.Mutex
	sequenceCounter int64
	shardIDCounter  int64
	storage         storage.TransactionalStorageWith2PC
}

// NewKinesisStore creates a new KinesisStore instance with the specified storage, account ID, and region.
func NewKinesisStore(store storage.TransactionalStorageWith2PC, accountID, region string) *KinesisStore {
	return &KinesisStore{
		BaseStore:      common.NewBaseStore(store.Bucket("kinesis-streams-"+region), "kinesis-streams"),
		shardsStore:    common.NewBaseStore(store.Bucket("kinesis-shards-"+region), "kinesis-shards"),
		recordsStore:   common.NewBaseStore(store.Bucket("kinesis-records-"+region), "kinesis-records"),
		consumersStore: common.NewBaseStore(store.Bucket("kinesis-consumers-"+region), "kinesis-consumers"),
		iteratorsStore: common.NewBaseStore(store.Bucket("kinesis-iterators-"+region), "kinesis-iterators"),
		TagStore:       common.NewTagStoreWithRegion(store, "kinesis", region),
		arnBuilder:     svcarn.NewARNBuilder(accountID, region),
		accountID:      accountID,
		region:         region,
		storage:        store,
	}
}

// GetAccountID returns the account ID associated with this store.
func (s *KinesisStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the region associated with this store.
func (s *KinesisStore) GetRegion() string {
	return s.region
}

func (s *KinesisStore) buildStreamARN(streamName string) string {
	return s.arnBuilder.Kinesis().Stream(streamName)
}

// BuildStreamARN constructs the ARN for a Kinesis stream.
func (s *KinesisStore) BuildStreamARN(streamName string) string {
	return s.buildStreamARN(streamName)
}

func (s *KinesisStore) buildConsumerARN(streamName, consumerName string) string {
	return s.arnBuilder.Kinesis().Consumer(streamName, consumerName)
}

// CreateStream creates a new Kinesis stream with the specified name, shard count, and stream mode.
func (s *KinesisStore) CreateStream(streamName string, shardCount int32, streamMode StreamMode) (*Stream, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := streamName
	if s.Exists(key) {
		return nil, ErrStreamAlreadyExists
	}

	stream := NewStream(streamName, shardCount, streamMode)
	stream.StreamARN = s.buildStreamARN(streamName)

	if stream.StreamModeDetails.StreamMode == StreamModeOnDemand {
		shardCount = 4
	}
	stream.ShardCount = shardCount

	if err := s.PutProto(key, StreamToProto(stream)); err != nil {
		return nil, err
	}

	for i := 0; i < int(shardCount); i++ {
		shardID := fmt.Sprintf("shardId-%012d", i)
		startingHashKey := uint64(i) * (math.MaxUint64 / uint64(shardCount))
		endingHashKey := uint64(i+1)*(math.MaxUint64/uint64(shardCount)) - 1
		if i == int(shardCount)-1 {
			endingHashKey = math.MaxUint64
		}

		shard := NewShard(shardID, streamName, fmt.Sprintf("%d", startingHashKey), fmt.Sprintf("%d", endingHashKey))
		shard.SequenceNumberRange = &SequenceNumberRange{
			StartingSequenceNumber: s.generateSequenceNumber(shardID),
		}

		if err := s.PutShard(shard); err != nil {
			return nil, err
		}
	}

	return stream, nil
}

// GetStream retrieves a Kinesis stream by its name.
func (s *KinesisStore) GetStream(streamName string) (*Stream, error) {
	var pbStream pb.Stream
	if err := s.GetProto(streamName, &pbStream); err != nil {
		return nil, ErrStreamNotFound
	}
	return ProtoToStream(&pbStream), nil
}

func (s *KinesisStore) getStreamInTxn(txn storage.Transaction, streamName string) (*Stream, error) {
	bucket := txn.Bucket("kinesis-streams-" + s.region)
	data, err := bucket.Get([]byte(streamName))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrStreamNotFound
	}
	var pbStream pb.Stream
	if err := proto.Unmarshal(data, &pbStream); err != nil {
		return nil, err
	}
	return ProtoToStream(&pbStream), nil
}

func (s *KinesisStore) updateStreamInTxn(txn storage.Transaction, stream *Stream) error {
	stream.LastModifiedAt = time.Now().UTC()
	bucket := txn.Bucket("kinesis-streams-" + s.region)
	data, err := proto.Marshal(StreamToProto(stream))
	if err != nil {
		return err
	}
	return bucket.Put([]byte(stream.StreamName), data)
}

// GetStreamByARN retrieves a Kinesis stream by its ARN.
func (s *KinesisStore) GetStreamByARN(streamARN string) (*Stream, error) {
	normalizedARN := normalizeARN(streamARN, s.accountID)

	streams, err := s.ListStreams()
	if err != nil {
		return nil, err
	}
	for _, stream := range streams {
		if normalizeARN(stream.StreamARN, s.accountID) == normalizedARN {
			return stream, nil
		}
	}
	return nil, ErrStreamNotFound
}

// UpdateStream updates an existing Kinesis stream.
func (s *KinesisStore) UpdateStream(stream *Stream) error {
	stream.LastModifiedAt = time.Now().UTC()
	return s.PutProto(stream.StreamName, StreamToProto(stream))
}

// DeleteStream deletes a Kinesis stream by its name.
func (s *KinesisStore) DeleteStream(streamName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(streamName) {
		return ErrStreamNotFound
	}

	shards, err := s.ListShards(streamName, nil, "", 0)
	if err != nil {
		return err
	}
	for _, shard := range shards {
		prefix := fmt.Sprintf("%s#%s#", streamName, shard.ShardID)
		if err := s.recordsStore.ScanPrefix(prefix, func(key string, value []byte) error {
			return s.recordsStore.Delete(key)
		}); err != nil {
			return fmt.Errorf("failed to delete records for shard %s: %w", shard.ShardID, err)
		}
		if err := s.shardsStore.Delete(fmt.Sprintf("%s#%s", streamName, shard.ShardID)); err != nil {
			return fmt.Errorf("failed to delete shard %s: %w", shard.ShardID, err)
		}
	}

	consumers, err := s.ListStreamConsumers(streamName)
	if err != nil {
		return err
	}
	for _, consumer := range consumers {
		if err := s.consumersStore.Delete(consumer.ConsumerARN); err != nil {
			return fmt.Errorf("failed to delete consumer %s: %w", consumer.ConsumerARN, err)
		}
	}

	if err := s.TagStore.Delete(streamName); err != nil {
		return fmt.Errorf("failed to delete tags: %w", err)
	}

	return s.BaseStore.Delete(streamName)
}

// ListStreams returns all Kinesis streams in the store.
func (s *KinesisStore) ListStreams() ([]*Stream, error) {
	opts := common.ListOptions{MaxItems: 10000}
	result, err := common.ListProto(s.BaseStore, opts, func() *pb.Stream { return &pb.Stream{} }, nil)
	if err != nil {
		return nil, err
	}
	streams := make([]*Stream, len(result.Items))
	for i, pbStream := range result.Items {
		streams[i] = ProtoToStream(pbStream)
	}
	return streams, nil
}

// PutShard stores a shard in the data store.
func (s *KinesisStore) PutShard(shard *Shard) error {
	key := fmt.Sprintf("%s#%s", shard.StreamName, shard.ShardID)
	return s.shardsStore.PutProto(key, ShardToProto(shard))
}

// GetShard retrieves a shard from a stream by its ID.
func (s *KinesisStore) GetShard(streamName, shardID string) (*Shard, error) {
	key := fmt.Sprintf("%s#%s", streamName, shardID)
	var pbShard pb.Shard
	if err := s.shardsStore.GetProto(key, &pbShard); err != nil {
		return nil, ErrShardNotFound
	}
	return ProtoToShard(&pbShard), nil
}

// UpdateShard updates an existing shard in the data store.
func (s *KinesisStore) UpdateShard(shard *Shard) error {
	return s.PutShard(shard)
}

// ListShards returns all shards in a stream, optionally filtered by a ShardFilter.
var errStopIteration = errors.New("stop iteration")

// ListShards returns all shards in a stream, optionally filtered by a ShardFilter and paginated.
func (s *KinesisStore) ListShards(streamName string, filter *ShardFilter, exclusiveStartShardID string, limit int) ([]*Shard, error) {
	prefix := streamName + "#"
	var shards []*Shard
	skip := exclusiveStartShardID != ""

	err := s.shardsStore.ScanPrefix(prefix, func(key string, value []byte) error {
		if skip {
			var pbShard pb.Shard
			if err := proto.Unmarshal(value, &pbShard); err != nil {
				return err
			}
			if pbShard.ShardId == exclusiveStartShardID {
				skip = false
			}
			return nil
		}

		var pbShard pb.Shard
		if err := proto.Unmarshal(value, &pbShard); err != nil {
			return err
		}
		shard := ProtoToShard(&pbShard)
		if filter == nil || s.shardMatchesFilter(shard, filter) {
			shards = append(shards, shard)
			if limit > 0 && len(shards) >= limit {
				return errStopIteration
			}
		}
		return nil
	})
	if errors.Is(err, errStopIteration) {
		err = nil
	}
	return shards, err
}

func (s *KinesisStore) listShardsInTxn(txn storage.Transaction, streamName string) ([]*Shard, error) {
	bucket := txn.Bucket("kinesis-shards-" + s.region)
	prefix := []byte(streamName + "#")
	var shards []*Shard

	iter := bucket.ScanPrefix(prefix)
	for iter.Next() {
		var pbShard pb.Shard
		if err := proto.Unmarshal(iter.Value(), &pbShard); err != nil {
			return nil, err
		}
		shard := ProtoToShard(&pbShard)
		if shard.SequenceNumberRange.EndingSequenceNumber == "" {
			shards = append(shards, shard)
		}
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return shards, nil
}

func (s *KinesisStore) shardMatchesFilter(shard *Shard, filter *ShardFilter) bool {
	switch filter.Type {
	case "AT_SHARD_ID":
		return shard.ShardID == filter.ShardID
	case "FROM_TIMESTAMP":
		if filter.Timestamp == nil {
			return true
		}
		return shard.CreatedAt.After(*filter.Timestamp) || shard.CreatedAt.Equal(*filter.Timestamp)
	case "AFTER_SHARD_ID":
		return shard.ShardID > filter.ShardID
	case "TRIM_HORIZON", "LATEST":
		return true
	}
	return true
}

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

// CreateShardIterator creates a shard iterator for reading from a Kinesis stream.
func (s *KinesisStore) CreateShardIterator(streamName, shardID, iteratorType, seqNum string, timestamp *time.Time) (*ShardIterator, error) {
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

// SplitShard splits a shard into two new shards.
func (s *KinesisStore) SplitShard(streamName, shardID, newStartingHashKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.splitShardInternal(streamName, shardID, newStartingHashKey)
}

func (s *KinesisStore) splitShardInternal(streamName, shardID, newStartingHashKey string) error {
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return s.splitShardInTxn(txn, streamName, shardID, newStartingHashKey)
	})
}

func (s *KinesisStore) splitShardInTxn(txn storage.Transaction, streamName, shardID, newStartingHashKey string) error {
	shard, err := s.getShardInTxn(txn, streamName, shardID)
	if err != nil {
		return err
	}

	if shard.SequenceNumberRange.EndingSequenceNumber != "" {
		return ErrShardClosed
	}

	shard.SequenceNumberRange.EndingSequenceNumber = shard.LatestSequenceNumber
	if shard.LatestSequenceNumber == "" {
		shard.SequenceNumberRange.EndingSequenceNumber = shard.SequenceNumberRange.StartingSequenceNumber
	}
	if err := s.putShardInTxn(txn, shard); err != nil {
		return err
	}

	startHash, ok := new(big.Int).SetString(shard.HashKeyRange.StartingHashKey, 10)
	if !ok {
		return fmt.Errorf("invalid starting hash key: %s", shard.HashKeyRange.StartingHashKey)
	}
	endHash, ok := new(big.Int).SetString(shard.HashKeyRange.EndingHashKey, 10)
	if !ok {
		return fmt.Errorf("invalid ending hash key: %s", shard.HashKeyRange.EndingHashKey)
	}
	var midHash *big.Int
	if newStartingHashKey != "" {
		var ok bool
		midHash, ok = new(big.Int).SetString(newStartingHashKey, 10)
		if !ok {
			return fmt.Errorf("invalid new starting hash key: %s", newStartingHashKey)
		}
	}
	if midHash == nil {
		midHash = new(big.Int).Div(new(big.Int).Add(startHash, endHash), big.NewInt(2))
	}

	shardCounter1 := atomic.AddInt64(&s.shardIDCounter, 1)
	newShardID1 := fmt.Sprintf("shardId-%012d", shardCounter1)
	newShard1 := NewShard(newShardID1, streamName, shard.HashKeyRange.StartingHashKey, midHash.String())
	newShard1.ParentShardID = shardID
	newShard1.SequenceNumberRange = &SequenceNumberRange{
		StartingSequenceNumber: s.generateSequenceNumber(newShardID1),
	}

	shardCounter2 := atomic.AddInt64(&s.shardIDCounter, 1)
	newShardID2 := fmt.Sprintf("shardId-%012d", shardCounter2)
	newShard2 := NewShard(newShardID2, streamName, midHash.String(), shard.HashKeyRange.EndingHashKey)
	newShard2.ParentShardID = shardID
	newShard2.SequenceNumberRange = &SequenceNumberRange{
		StartingSequenceNumber: s.generateSequenceNumber(newShardID2),
	}

	if err := s.putShardInTxn(txn, newShard1); err != nil {
		return err
	}
	if err := s.putShardInTxn(txn, newShard2); err != nil {
		return err
	}

	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return err
	}
	stream.ShardCount++
	return s.updateStreamInTxn(txn, stream)
}

func (s *KinesisStore) getShardInTxn(txn storage.Transaction, streamName, shardID string) (*Shard, error) {
	bucket := txn.Bucket("kinesis-shards-" + s.region)
	key := []byte(fmt.Sprintf("%s#%s", streamName, shardID))
	data, err := bucket.Get(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrShardNotFound
	}
	var pbShard pb.Shard
	if err := proto.Unmarshal(data, &pbShard); err != nil {
		return nil, err
	}
	return ProtoToShard(&pbShard), nil
}

func (s *KinesisStore) putShardInTxn(txn storage.Transaction, shard *Shard) error {
	bucket := txn.Bucket("kinesis-shards-" + s.region)
	key := []byte(fmt.Sprintf("%s#%s", shard.StreamName, shard.ShardID))
	data, err := proto.Marshal(ShardToProto(shard))
	if err != nil {
		return err
	}
	return bucket.Put(key, data)
}

// MergeShards merges two adjacent shards into one.
func (s *KinesisStore) MergeShards(streamName, shardID1, shardID2 string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mergeShardsInternal(streamName, shardID1, shardID2)
}

func (s *KinesisStore) mergeShardsInternal(streamName, shardID1, shardID2 string) error {
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return s.mergeShardsInTxn(txn, streamName, shardID1, shardID2)
	})
}

func (s *KinesisStore) mergeShardsInTxn(txn storage.Transaction, streamName, shardID1, shardID2 string) error {
	shard1, err := s.getShardInTxn(txn, streamName, shardID1)
	if err != nil {
		return err
	}
	shard2, err := s.getShardInTxn(txn, streamName, shardID2)
	if err != nil {
		return err
	}

	if !s.areAdjacentShards(shard1, shard2) {
		return ErrShardsNotAdjacent
	}

	shard1.SequenceNumberRange.EndingSequenceNumber = shard1.LatestSequenceNumber
	if shard1.LatestSequenceNumber == "" {
		shard1.SequenceNumberRange.EndingSequenceNumber = shard1.SequenceNumberRange.StartingSequenceNumber
	}
	shard2.SequenceNumberRange.EndingSequenceNumber = shard2.LatestSequenceNumber
	if shard2.LatestSequenceNumber == "" {
		shard2.SequenceNumberRange.EndingSequenceNumber = shard2.SequenceNumberRange.StartingSequenceNumber
	}
	if err := s.putShardInTxn(txn, shard1); err != nil {
		return err
	}
	if err := s.putShardInTxn(txn, shard2); err != nil {
		return err
	}

	start1, ok1 := new(big.Int).SetString(shard1.HashKeyRange.StartingHashKey, 10)
	start2, ok2 := new(big.Int).SetString(shard2.HashKeyRange.StartingHashKey, 10)
	if !ok1 || !ok2 {
		return fmt.Errorf("invalid hash key range")
	}
	startHash := shard1.HashKeyRange.StartingHashKey
	endHash := shard2.HashKeyRange.EndingHashKey
	if start1.Cmp(start2) > 0 {
		startHash = shard2.HashKeyRange.StartingHashKey
		endHash = shard1.HashKeyRange.EndingHashKey
	}

	shardCounter := atomic.AddInt64(&s.shardIDCounter, 1)
	newShardID := fmt.Sprintf("shardId-%012d", shardCounter)
	newShard := NewShard(newShardID, streamName, startHash, endHash)
	newShard.ParentShardID = shardID1
	newShard.AdjacentParentShardID = shardID2
	newShard.SequenceNumberRange = &SequenceNumberRange{
		StartingSequenceNumber: s.generateSequenceNumber(newShardID),
	}

	if err := s.putShardInTxn(txn, newShard); err != nil {
		return err
	}

	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return err
	}
	stream.ShardCount--
	return s.updateStreamInTxn(txn, stream)
}

func (s *KinesisStore) areAdjacentShards(shard1, shard2 *Shard) bool {
	end1, ok1 := new(big.Int).SetString(shard1.HashKeyRange.EndingHashKey, 10)
	start2, ok2 := new(big.Int).SetString(shard2.HashKeyRange.StartingHashKey, 10)
	if ok1 && ok2 {
		diff := new(big.Int).Sub(start2, end1)
		if diff.Cmp(big.NewInt(1)) == 0 {
			return true
		}
	}

	end2, ok2 := new(big.Int).SetString(shard2.HashKeyRange.EndingHashKey, 10)
	start1, ok1 := new(big.Int).SetString(shard1.HashKeyRange.StartingHashKey, 10)
	if ok1 && ok2 {
		diff := new(big.Int).Sub(start1, end2)
		if diff.Cmp(big.NewInt(1)) == 0 {
			return true
		}
	}

	return false
}

// UpdateShardCount updates the shard count for a stream.
func (s *KinesisStore) UpdateShardCount(streamName string, targetCount int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		return s.updateShardCountInTxn(txn, streamName, targetCount)
	})
}

func (s *KinesisStore) updateShardCountInTxn(txn storage.Transaction, streamName string, targetCount int32) error {
	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return err
	}

	if stream.StreamModeDetails != nil && stream.StreamModeDetails.StreamMode == StreamModeOnDemand {
		return ErrOperationNotSupportedOnOnDemandStream
	}

	currentShards, err := s.listShardsInTxn(txn, streamName)
	if err != nil {
		return err
	}

	var activeShards []*Shard
	for _, shard := range currentShards {
		if shard.SequenceNumberRange.EndingSequenceNumber == "" {
			activeShards = append(activeShards, shard)
		}
	}

	currentCount := int32(len(activeShards))
	if currentCount == targetCount {
		return nil
	}
	if currentCount == 0 {
		return ErrInvalidShardCount
	}

	ratio := float64(targetCount) / float64(currentCount)
	if ratio > 2.0 || ratio < 0.5 {
		return ErrInvalidShardCount
	}

	if targetCount > currentCount {
		splitsNeeded := targetCount - currentCount
		for i := 0; i < int(splitsNeeded); i++ {
			refreshedShards, err := s.listShardsInTxn(txn, streamName)
			if err != nil {
				return err
			}
			var splitCandidates []*Shard
			for _, shard := range refreshedShards {
				if shard.SequenceNumberRange.EndingSequenceNumber == "" {
					splitCandidates = append(splitCandidates, shard)
				}
			}
			if len(splitCandidates) == 0 {
				break
			}
			shardToSplit := splitCandidates[i%len(splitCandidates)]
			if err := s.splitShardInTxn(txn, streamName, shardToSplit.ShardID, ""); err != nil {
				return err
			}
		}
	} else {
		mergesNeeded := currentCount - targetCount
		for i := 0; i < int(mergesNeeded); i++ {
			refreshedShards, err := s.listShardsInTxn(txn, streamName)
			if err != nil {
				return err
			}
			var mergeCandidates []*Shard
			for _, shard := range refreshedShards {
				if shard.SequenceNumberRange.EndingSequenceNumber == "" {
					mergeCandidates = append(mergeCandidates, shard)
				}
			}
			if len(mergeCandidates) < 2 {
				break
			}
			sort.Slice(mergeCandidates, func(i, j int) bool {
				return mergeCandidates[i].HashKeyRange.StartingHashKey < mergeCandidates[j].HashKeyRange.StartingHashKey
			})
			merged := false
			for j := 0; j < len(mergeCandidates)-1; j++ {
				if s.areAdjacentShards(mergeCandidates[j], mergeCandidates[j+1]) {
					if err := s.mergeShardsInTxn(txn, streamName, mergeCandidates[j].ShardID, mergeCandidates[j+1].ShardID); err != nil {
						return err
					}
					merged = true
					break
				}
			}
			if !merged {
				break
			}
		}
	}

	finalShards, err := s.listShardsInTxn(txn, streamName)
	if err != nil {
		return err
	}
	stream.ShardCount = int32(len(finalShards))
	return s.updateStreamInTxn(txn, stream)
}

// RegisterStreamConsumer registers a consumer for a Kinesis stream.
func (s *KinesisStore) RegisterStreamConsumer(streamARN, consumerName string) (*Consumer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, err := s.GetStreamByARN(streamARN)
	if err != nil {
		return nil, err
	}

	consumerARN := s.buildConsumerARN(stream.StreamName, consumerName)
	if s.consumersStore.Exists(consumerARN) {
		return nil, ErrConsumerAlreadyExists
	}

	consumer := NewConsumer(consumerName, streamARN)
	consumer.ConsumerARN = consumerARN

	if err := s.consumersStore.PutProto(consumerARN, ConsumerToProto(consumer)); err != nil {
		return nil, err
	}

	stream.ConsumerCount++
	if err := s.UpdateStream(stream); err != nil {
		return nil, err
	}

	return consumer, nil
}

// DeregisterStreamConsumer deregisters a consumer from a Kinesis stream.
func (s *KinesisStore) DeregisterStreamConsumer(consumerARN string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var pbConsumer pb.Consumer
	if err := s.consumersStore.GetProto(consumerARN, &pbConsumer); err != nil {
		return ErrConsumerNotFound
	}
	consumer := ProtoToConsumer(&pbConsumer)

	stream, err := s.GetStreamByARN(consumer.StreamARN)
	if err != nil {
		return err
	}

	if err := s.consumersStore.Delete(consumerARN); err != nil {
		return err
	}

	stream.ConsumerCount--
	if err := s.UpdateStream(stream); err != nil {
		return err
	}

	return nil
}

// GetStreamConsumer retrieves a consumer by its ARN.
func (s *KinesisStore) GetStreamConsumer(consumerARN string) (*Consumer, error) {
	var pbConsumer pb.Consumer
	if err := s.consumersStore.GetProto(consumerARN, &pbConsumer); err != nil {
		return nil, ErrConsumerNotFound
	}
	return ProtoToConsumer(&pbConsumer), nil
}

// GetStreamConsumerByName retrieves a consumer by stream ARN and consumer name.
func (s *KinesisStore) GetStreamConsumerByName(streamARN, consumerName string) (*Consumer, error) {
	stream, err := s.GetStreamByARN(streamARN)
	if err != nil {
		return nil, err
	}

	consumerARN := s.buildConsumerARN(stream.StreamName, consumerName)
	return s.GetStreamConsumer(consumerARN)
}

// ListStreamConsumers lists all consumers for a stream.
func (s *KinesisStore) ListStreamConsumers(streamName string) ([]*Consumer, error) {
	stream, err := s.GetStream(streamName)
	if err != nil {
		return nil, err
	}

	prefix := stream.StreamARN + "/"
	var consumers []*Consumer

	err = s.consumersStore.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbConsumer pb.Consumer
		if err := proto.Unmarshal(value, &pbConsumer); err != nil {
			return err
		}
		consumer := ProtoToConsumer(&pbConsumer)
		if consumer.StreamARN == stream.StreamARN {
			consumers = append(consumers, consumer)
		}
		return nil
	})

	return consumers, err
}

func (s *KinesisStore) generateSequenceNumber(shardID string) string {
	ts := time.Now().UTC().UnixNano()
	counter := atomic.AddInt64(&s.sequenceCounter, 1)
	return fmt.Sprintf("%d%012d%012d", ts, counter, s.hashToInt(shardID))
}

func (s *KinesisStore) hashToInt(str string) int64 {
	h := md5.Sum([]byte(str))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

// PutRecordRequest represents a request to put a record into a Kinesis stream.
type PutRecordRequest struct {
	Data         string `json:"data"`
	PartitionKey string `json:"partitionKey"`
}

// PutRecordResult represents the result of putting a record into a Kinesis stream.
type PutRecordResult struct {
	SequenceNumber string `json:"sequenceNumber,omitempty"`
	ShardID        string `json:"shardId,omitempty"`
	ErrorCode      string `json:"errorCode,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
}

func normalizeARN(arn string, accountID string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) >= 5 && parts[4] == "" {
		parts[4] = accountID
		return strings.Join(parts, ":")
	}
	return arn
}

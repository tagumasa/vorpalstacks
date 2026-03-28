package kinesis

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync/atomic"

	"vorpalstacks/internal/core/storage"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

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
				return fmt.Errorf("unmarshal shard for skip check: %w", err)
			}
			if pbShard.ShardId == exclusiveStartShardID {
				skip = false
			}
			return nil
		}

		var pbShard pb.Shard
		if err := proto.Unmarshal(value, &pbShard); err != nil {
			return fmt.Errorf("unmarshal shard: %w", err)
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
	if err != nil {
		return shards, fmt.Errorf("scan shards prefix: %w", err)
	}
	return shards, nil
}

func (s *KinesisStore) listShardsInTxn(txn storage.Transaction, streamName string) ([]*Shard, error) {
	bucket := txn.Bucket("kinesis-shards-" + s.region)
	prefix := []byte(streamName + "#")
	var shards []*Shard

	iter := bucket.ScanPrefix(prefix)
	for iter.Next() {
		var pbShard pb.Shard
		if err := proto.Unmarshal(iter.Value(), &pbShard); err != nil {
			return nil, fmt.Errorf("unmarshal shard: %w", err)
		}
		shard := ProtoToShard(&pbShard)
		if shard.SequenceNumberRange.EndingSequenceNumber == "" {
			shards = append(shards, shard)
		}
	}
	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("scan shards iterator: %w", err)
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

func (s *KinesisStore) getShardInTxn(txn storage.Transaction, streamName, shardID string) (*Shard, error) {
	bucket := txn.Bucket("kinesis-shards-" + s.region)
	key := []byte(fmt.Sprintf("%s#%s", streamName, shardID))
	data, err := bucket.Get(key)
	if err != nil {
		return nil, fmt.Errorf("get shard from bucket: %w", err)
	}
	if data == nil {
		return nil, ErrShardNotFound
	}
	var pbShard pb.Shard
	if err := proto.Unmarshal(data, &pbShard); err != nil {
		return nil, fmt.Errorf("unmarshal shard data: %w", err)
	}
	return ProtoToShard(&pbShard), nil
}

func (s *KinesisStore) putShardInTxn(txn storage.Transaction, shard *Shard) error {
	bucket := txn.Bucket("kinesis-shards-" + s.region)
	key := []byte(fmt.Sprintf("%s#%s", shard.StreamName, shard.ShardID))
	data, err := proto.Marshal(ShardToProto(shard))
	if err != nil {
		return fmt.Errorf("marshal shard: %w", err)
	}
	return bucket.Put(key, data)
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
		return fmt.Errorf("get existing shard: %w", err)
	}

	if shard.SequenceNumberRange.EndingSequenceNumber != "" {
		return ErrShardClosed
	}

	shard.SequenceNumberRange.EndingSequenceNumber = shard.LatestSequenceNumber
	if shard.LatestSequenceNumber == "" {
		shard.SequenceNumberRange.EndingSequenceNumber = shard.SequenceNumberRange.StartingSequenceNumber
	}
	if err := s.putShardInTxn(txn, shard); err != nil {
		return fmt.Errorf("put closed shard: %w", err)
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
		return fmt.Errorf("put new shard 1: %w", err)
	}
	if err := s.putShardInTxn(txn, newShard2); err != nil {
		return fmt.Errorf("put new shard 2: %w", err)
	}

	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("get stream for split: %w", err)
	}
	stream.ShardCount++
	return s.updateStreamInTxn(txn, stream)
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
		return fmt.Errorf("get shard 1: %w", err)
	}
	shard2, err := s.getShardInTxn(txn, streamName, shardID2)
	if err != nil {
		return fmt.Errorf("get shard 2: %w", err)
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
		return fmt.Errorf("put closed shard 1: %w", err)
	}
	if err := s.putShardInTxn(txn, shard2); err != nil {
		return fmt.Errorf("put closed shard 2: %w", err)
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
		return fmt.Errorf("put merged shard: %w", err)
	}

	stream, err := s.getStreamInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("get stream for merge: %w", err)
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
		return fmt.Errorf("get stream: %w", err)
	}

	if stream.StreamModeDetails != nil && stream.StreamModeDetails.StreamMode == StreamModeOnDemand {
		return ErrOperationNotSupportedOnOnDemandStream
	}

	currentShards, err := s.listShardsInTxn(txn, streamName)
	if err != nil {
		return fmt.Errorf("list current shards: %w", err)
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
				return fmt.Errorf("list shards for split: %w", err)
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
				return fmt.Errorf("split shard: %w", err)
			}
		}
	} else {
		mergesNeeded := currentCount - targetCount
		for i := 0; i < int(mergesNeeded); i++ {
			refreshedShards, err := s.listShardsInTxn(txn, streamName)
			if err != nil {
				return fmt.Errorf("list shards for merge: %w", err)
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
						return fmt.Errorf("merge adjacent shards: %w", err)
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
		return fmt.Errorf("list final shards: %w", err)
	}
	stream.ShardCount = int32(len(finalShards))
	return s.updateStreamInTxn(txn, stream)
}

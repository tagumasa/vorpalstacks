package kinesis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

// CreateStream creates a new Kinesis stream with the specified parameters.
// The stream record, its ARN-index entry, every initial shard and the
// create-time tag set are one transaction: a mid-creation failure leaves
// no half-created stream behind for a retry to collide with, and no
// tagless stream either — the tags commit with the resource or not at
// all.
func (s *KinesisStore) CreateStream(streamName string, shardCount int32, streamMode StreamMode, maxRecordSizeInKiB int32, warmThroughputMiBps int32, tags map[string]string) (*Stream, error) {
	// The initial layout divides the hash-key space by the shard count
	// below; a count below one is a caller programming error that must
	// surface as an error, never reach the division.
	if shardCount < 1 {
		return nil, ErrInvalidShardCount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := streamName
	if s.Exists(key) {
		return nil, ErrStreamAlreadyExists
	}

	stream := NewStream(streamName, shardCount, streamMode, maxRecordSizeInKiB, warmThroughputMiBps)
	stream.StreamARN = s.buildStreamARN(streamName)

	if stream.StreamModeDetails.StreamMode == StreamModeOnDemand {
		shardCount = DefaultOnDemandShardCount
	}
	stream.ShardCount = shardCount

	streamData, err := proto.Marshal(StreamToProto(stream))
	if err != nil {
		return nil, err
	}

	shardData := make([][]byte, shardCount)
	shardIDs := make([]string, shardCount)
	// The hash-key space is the model's 128-bit domain: shards tile
	// [0, 2^128-1] with inclusive ranges, evenly for power-of-two counts
	// (AWS's own single-shard streams report [0, 2^128-1]), and the last
	// shard's end is pinned to the space top so no value is unroutable.
	space := MaxShardHashKeyInt()
	width := new(big.Int).Add(space, big.NewInt(1))
	width.Div(width, big.NewInt(int64(shardCount)))
	for i := 0; i < int(shardCount); i++ {
		shardID := fmt.Sprintf("shardId-%012d", i)
		startingHashKey := new(big.Int).Mul(width, big.NewInt(int64(i)))
		endingHashKey := new(big.Int).Mul(width, big.NewInt(int64(i+1)))
		endingHashKey.Sub(endingHashKey, big.NewInt(1))
		if i == int(shardCount)-1 {
			endingHashKey = space
		}

		shard := NewShard(shardID, streamName, startingHashKey.String(), endingHashKey.String())
		shard.SequenceNumberRange = &SequenceNumberRange{
			StartingSequenceNumber: s.generateSequenceNumber(shardID),
		}
		data, err := proto.Marshal(ShardToProto(shard))
		if err != nil {
			return nil, err
		}
		shardData[i] = data
		shardIDs[i] = shardID
	}

	err = s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		streamsBucket := txn.Bucket(s.streamsBucketName())
		if err := streamsBucket.Put([]byte(stream.StreamName), streamData); err != nil {
			return err
		}
		arnBucket := txn.Bucket(s.arnIndexBucketName())
		if err := arnBucket.Put([]byte("arn#"+stream.StreamARN), []byte(stream.StreamName)); err != nil {
			return err
		}
		shardsBucket := txn.Bucket(s.shardsBucketName())
		for i, data := range shardData {
			key := fmt.Sprintf("%s#%s", stream.StreamName, shardIDs[i])
			if err := shardsBucket.Put([]byte(key), data); err != nil {
				return err
			}
		}
		// The initial tag set commits with the stream itself: a tag
		// failure leaves no half-created stream for a retry to collide
		// with, and the create-time cap is the fresh set's own size.
		return s.TagStore.TagInTxn(txn, streamName, tags, MaxTagsPerResource)
	})
	if err != nil {
		return nil, err
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
	bucket := txn.Bucket(s.streamsBucketName())
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
	bucket := txn.Bucket(s.streamsBucketName())
	data, err := proto.Marshal(StreamToProto(stream))
	if err != nil {
		return err
	}
	return bucket.Put([]byte(stream.StreamName), data)
}

// GetStreamByARN retrieves a Kinesis stream by its ARN through the ARN
// index — the single lookup path. The index is written by CreateStream
// and removed by DeleteStream, so a miss is a not-found.
func (s *KinesisStore) GetStreamByARN(streamARN string) (*Stream, error) {
	normalizedARN := s.normalizeARN(streamARN)

	data, err := s.arnIndexBucket.Get([]byte("arn#" + normalizedARN))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrStreamNotFound
	}
	return s.GetStream(string(data))
}

// UpdateStreamFields applies fn to the stream record under the store
// mutex: the read, the mutation and the write are one critical section, so
// a configuration change can neither write a concurrently deleted stream
// back into existence nor revert a counter another mutation advanced. fn
// returning an error aborts without writing.
func (s *KinesisStore) UpdateStreamFields(streamName string, fn func(*Stream) error) (*Stream, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream, err := s.GetStream(streamName)
	if err != nil {
		return nil, err
	}
	if err := fn(stream); err != nil {
		return nil, err
	}
	if err := s.updateStreamLocked(stream); err != nil {
		return nil, err
	}
	return stream, nil
}

// updateStreamLocked persists the stream record; the caller must hold s.mu.
func (s *KinesisStore) updateStreamLocked(stream *Stream) error {
	return s.PutProto(stream.StreamName, StreamToProto(stream))
}

// DeleteStream deletes a Kinesis stream by its name. Deleting a stream
// dissociates everything attached to it: the model states shards and tags
// go with the stream, and the platform extends the same completeness to the
// registered consumers and their tags, the attached resource policy, and
// the stream's shard iterators — a recreated stream inherits nothing.
func (s *KinesisStore) DeleteStream(streamName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(streamName) {
		return ErrStreamNotFound
	}

	stream, err := s.GetStream(streamName)
	if err != nil {
		return err
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
		if err := s.TagStore.Delete(consumer.ConsumerARN); err != nil {
			return fmt.Errorf("failed to delete tags for consumer %s: %w", consumer.ConsumerARN, err)
		}
		if err := s.consumersStore.Delete(consumer.ConsumerARN); err != nil {
			return fmt.Errorf("failed to delete consumer %s: %w", consumer.ConsumerARN, err)
		}
	}

	if err := s.TagStore.Delete(streamName); err != nil {
		return fmt.Errorf("failed to delete tags: %w", err)
	}

	if err := s.deleteResourcePolicyLocked(stream.StreamARN); err != nil {
		return fmt.Errorf("failed to delete the resource policy: %w", err)
	}

	// Shard iterators are server-side read cursors the stream invalidates;
	// they are keyed by ID, so the sweep walks the bucket and matches on the
	// owning stream.
	if err := s.iteratorsStore.ForEach(func(key string, value []byte) error {
		var it ShardIterator
		if err := json.Unmarshal(value, &it); err != nil {
			return nil
		}
		if it.StreamName == streamName {
			return s.iteratorsStore.Delete(key)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to sweep shard iterators: %w", err)
	}

	if err := s.arnIndexBucket.Delete([]byte("arn#" + s.normalizeARN(stream.StreamARN))); err != nil {
		return fmt.Errorf("failed to delete the stream ARN index: %w", err)
	}

	return s.BaseStore.Delete(streamName)
}

// ListStreams returns Kinesis streams with pagination support.
func (s *KinesisStore) ListStreams(opts common.ListOptions) (*common.ListResult[Stream], error) {
	result, err := common.ListProto(s.BaseStore, opts, func() *pb.Stream { return &pb.Stream{} }, nil)
	if err != nil {
		return nil, err
	}
	streams := make([]*Stream, len(result.Items))
	for i, pbStream := range result.Items {
		streams[i] = ProtoToStream(pbStream)
	}
	return &common.ListResult[Stream]{
		Items:       streams,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

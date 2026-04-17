package kinesis

import (
	"fmt"
	"math"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"

	"google.golang.org/protobuf/proto"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

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

	arnBucket := s.storage.Bucket("kinesis-streams-arn-" + s.region)
	_ = arnBucket.Put([]byte("arn#"+stream.StreamARN), []byte(stream.StreamName))

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

	arnBucket := s.storage.Bucket("kinesis-streams-arn-" + s.region)
	arnKey := []byte("arn#" + normalizedARN)
	if data, err := arnBucket.Get(arnKey); err == nil && data != nil {
		streamName := string(data)
		stream, err := s.GetStream(streamName)
		if err == nil {
			return stream, nil
		}
	}

	streams, err := s.ListStreams()
	if err != nil {
		return nil, err
	}
	for _, stream := range streams {
		if normalizeARN(stream.StreamARN, s.accountID) == normalizedARN {
			_ = arnBucket.Put(arnKey, []byte(stream.StreamName))
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
		if err := s.consumersStore.Delete(consumer.ConsumerARN); err != nil {
			return fmt.Errorf("failed to delete consumer %s: %w", consumer.ConsumerARN, err)
		}
	}

	if err := s.TagStore.Delete(streamName); err != nil {
		return fmt.Errorf("failed to delete tags: %w", err)
	}

	arnBucket := s.storage.Bucket("kinesis-streams-arn-" + s.region)
	_ = arnBucket.Delete([]byte("arn#" + normalizeARN(stream.StreamARN, s.accountID)))

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

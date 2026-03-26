package kinesis

import (
	"time"
)

// KinesisStoreInterface defines operations for managing Kinesis streams.
type KinesisStoreInterface interface {
	GetAccountID() string
	GetRegion() string
	BuildStreamARN(streamName string) string
	CreateStream(streamName string, shardCount int32, streamMode StreamMode) (*Stream, error)
	GetStream(streamName string) (*Stream, error)
	GetStreamByARN(streamARN string) (*Stream, error)
	UpdateStream(stream *Stream) error
	DeleteStream(streamName string) error
	ListStreams() ([]*Stream, error)
	PutShard(shard *Shard) error
	GetShard(streamName, shardID string) (*Shard, error)
	UpdateShard(shard *Shard) error
	ListShards(streamName string, filter *ShardFilter, exclusiveStartShardID string, limit int) ([]*Shard, error)
	PutRecord(streamName, shardID, partitionKey, data string) (*Record, error)
	PutRecords(streamName string, records []PutRecordRequest) ([]PutRecordResult, error)
	SelectShardByPartitionKey(shards []*Shard, partitionKey string) *Shard
	GetRecords(streamName, shardID, startingSeqNum string, limit int32, includeStart bool) ([]*Record, string, error)
	CreateShardIterator(streamName, shardID, iteratorType, seqNum string, timestamp *time.Time) (*ShardIterator, error)
	GetShardIterator(iteratorID string) (*ShardIterator, error)
	SplitShard(streamName, shardID, newStartingHashKey string) error
	MergeShards(streamName, shardID1, shardID2 string) error
	UpdateShardCount(streamName string, targetCount int32) error
	RegisterStreamConsumer(streamARN, consumerName string) (*Consumer, error)
	DeregisterStreamConsumer(consumerARN string) error
	GetStreamConsumer(consumerARN string) (*Consumer, error)
	GetStreamConsumerByName(streamARN, consumerName string) (*Consumer, error)
	ListStreamConsumers(streamName string) ([]*Consumer, error)
}

var _ KinesisStoreInterface = (*KinesisStore)(nil)

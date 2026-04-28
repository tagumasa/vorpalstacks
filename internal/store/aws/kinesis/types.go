// Package kinesis provides Kinesis data store functionality for vorpalstacks.
package kinesis

import (
	"time"
)

// StreamStatus represents the status of a Kinesis stream.
type StreamStatus string

// StreamStatus constants define the possible statuses of a Kinesis stream.
const (
	StreamStatusCreating StreamStatus = "CREATING"
	StreamStatusDeleting StreamStatus = "DELETING"
	StreamStatusActive   StreamStatus = "ACTIVE"
	StreamStatusUpdating StreamStatus = "UPDATING"
)

// StreamMode represents the mode of a Kinesis stream.
type StreamMode string

// StreamMode constants define the possible modes of a Kinesis stream.
const (
	StreamModeProvisioned StreamMode = "PROVISIONED"
	StreamModeOnDemand    StreamMode = "ON_DEMAND"
)

// Stream represents a Kinesis stream.
type Stream struct {
	StreamName           string                `json:"streamName"`
	StreamARN            string                `json:"streamArn"`
	StreamStatus         StreamStatus          `json:"streamStatus"`
	StreamModeDetails    *StreamModeDetails    `json:"streamModeDetails,omitempty"`
	ShardCount           int32                 `json:"shardCount"`
	RetentionPeriodHours int32                 `json:"retentionPeriodHours"`
	EnhancedMonitoring   []EnhancedMonitoring  `json:"enhancedMonitoring,omitempty"`
	EncryptionType       string                `json:"encryptionType,omitempty"`
	KeyID                string                `json:"keyId,omitempty"`
	ConsumerCount        int32                 `json:"consumerCount"`
	CreatedAt            time.Time             `json:"createdAt"`
	LastModifiedAt       time.Time             `json:"lastModifiedAt"`
	MaxRecordSizeInKiB   int32                 `json:"maxRecordSizeInKiB,omitempty"`
	OnDemandStreamConfig *OnDemandStreamConfig `json:"onDemandStreamConfig,omitempty"`
}

// OnDemandStreamConfig represents on-demand stream configuration.
type OnDemandStreamConfig struct {
	ShardCount    int32 `json:"shardCount,omitempty"`
	MaxShardCount int32 `json:"maxShardCount,omitempty"`
	OnDemandMode  bool  `json:"onDemandMode,omitempty"`
}

// StreamModeDetails represents the mode details of a Kinesis stream.
type StreamModeDetails struct {
	StreamMode StreamMode `json:"streamMode"`
}

// StreamSummary represents a summary of a Kinesis stream.
type StreamSummary struct {
	StreamName        string             `json:"streamName"`
	StreamARN         string             `json:"streamArn"`
	StreamStatus      StreamStatus       `json:"streamStatus"`
	StreamModeDetails *StreamModeDetails `json:"streamModeDetails,omitempty"`
	ConsumerCount     int32              `json:"consumerCount"`
}

// Shard represents a Kinesis shard.
type Shard struct {
	ShardID               string               `json:"shardId"`
	StreamName            string               `json:"streamName"`
	ParentShardID         string               `json:"parentShardId,omitempty"`
	AdjacentParentShardID string               `json:"adjacentParentShardId,omitempty"`
	HashKeyRange          *HashKeyRange        `json:"hashKeyRange"`
	SequenceNumberRange   *SequenceNumberRange `json:"sequenceNumberRange"`
	LatestSequenceNumber  string               `json:"latestSequenceNumber,omitempty"`
	CreatedAt             time.Time            `json:"createdAt"`
}

// HashKeyRange represents the hash key range for a Kinesis shard.
type HashKeyRange struct {
	StartingHashKey string `json:"startingHashKey"`
	EndingHashKey   string `json:"endingHashKey"`
}

// SequenceNumberRange represents the sequence number range for a Kinesis shard.
type SequenceNumberRange struct {
	StartingSequenceNumber string `json:"startingSequenceNumber"`
	EndingSequenceNumber   string `json:"endingSequenceNumber,omitempty"`
}

// Record represents a Kinesis record.
type Record struct {
	SequenceNumber              string    `json:"sequenceNumber"`
	ApproximateArrivalTimestamp time.Time `json:"approximateArrivalTimestamp"`
	Data                        string    `json:"data"`
	PartitionKey                string    `json:"partitionKey"`
	EncryptionType              string    `json:"encryptionType,omitempty"`
}

// Consumer represents a Kinesis consumer.
type Consumer struct {
	ConsumerName              string    `json:"consumerName"`
	ConsumerARN               string    `json:"consumerArn"`
	StreamARN                 string    `json:"streamArn"`
	ConsumerStatus            string    `json:"consumerStatus"`
	ConsumerCreationTimestamp time.Time `json:"consumerCreationTimestamp"`
}

// ConsumerSummary represents a summary of a Kinesis consumer.
type ConsumerSummary struct {
	ConsumerName              string    `json:"consumerName"`
	ConsumerARN               string    `json:"consumerArn"`
	StreamARN                 string    `json:"streamArn"`
	ConsumerStatus            string    `json:"consumerStatus"`
	ConsumerCreationTimestamp time.Time `json:"consumerCreationTimestamp"`
}

// ShardIterator represents a Kinesis shard iterator.
type ShardIterator struct {
	IteratorID     string     `json:"iteratorId"`
	StreamName     string     `json:"streamName"`
	ShardID        string     `json:"shardId"`
	IteratorType   string     `json:"iteratorType"`
	SequenceNumber string     `json:"sequenceNumber,omitempty"`
	Timestamp      *time.Time `json:"timestamp,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	ExpiresAt      time.Time  `json:"expiresAt"`
}

// EnhancedMonitoring represents enhanced monitoring settings for Kinesis.
type EnhancedMonitoring struct {
	ShardLevelMetrics []string `json:"shardLevelMetrics"`
}

// ShardFilter represents filter criteria for Kinesis shards.
type ShardFilter struct {
	Type      string     `json:"type"`
	ShardID   string     `json:"shardId,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// NewStream creates a new Kinesis stream with the specified name and shard count.
func NewStream(name string, shardCount int32, streamMode StreamMode) *Stream {
	now := time.Now().UTC()
	return &Stream{
		StreamName:           name,
		StreamStatus:         StreamStatusActive,
		StreamModeDetails:    &StreamModeDetails{StreamMode: streamMode},
		ShardCount:           shardCount,
		RetentionPeriodHours: 24,
		EnhancedMonitoring:   []EnhancedMonitoring{{ShardLevelMetrics: []string{}}},
		MaxRecordSizeInKiB:   1024,
		CreatedAt:            now,
		LastModifiedAt:       now,
	}
}

// NewShard creates a new Kinesis shard with the specified parameters.
func NewShard(shardID, streamName string, startingHashKey, endingHashKey string) *Shard {
	return &Shard{
		ShardID:    shardID,
		StreamName: streamName,
		HashKeyRange: &HashKeyRange{
			StartingHashKey: startingHashKey,
			EndingHashKey:   endingHashKey,
		},
		CreatedAt: time.Now().UTC(),
	}
}

// NewConsumer creates a new Kinesis consumer with the specified name and stream ARN.
func NewConsumer(consumerName, streamARN string) *Consumer {
	now := time.Now().UTC()
	return &Consumer{
		ConsumerName:              consumerName,
		StreamARN:                 streamARN,
		ConsumerStatus:            "ACTIVE",
		ConsumerCreationTimestamp: now,
	}
}

package kinesis

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrStreamNotFound is returned when the specified Kinesis stream
	// does not exist.
	ErrStreamNotFound = errors.New("stream not found")

	// ErrStreamAlreadyExists is returned when attempting to create a stream
	// that already exists.
	ErrStreamAlreadyExists = errors.New("stream already exists")

	// ErrShardNotFound is returned when the specified shard does not exist.
	ErrShardNotFound = errors.New("shard not found")

	// ErrShardClosed is returned when the shard has been closed and can no
	// longer accept records.
	ErrShardClosed = errors.New("shard is closed")

	// ErrNoActiveShards is returned when there are no active shards available
	// to read from.
	ErrNoActiveShards = errors.New("no active shards available")

	// ErrInvalidIterator is returned when the shard iterator is not valid.
	ErrInvalidIterator = errors.New("invalid shard iterator")

	// ErrExpiredIterator is returned when the shard iterator has expired
	// (typically after 5 minutes).
	ErrExpiredIterator = errors.New("shard iterator has expired")

	// ErrShardsNotAdjacent is returned when the shards are not adjacent and
	// cannot be merged or split.
	ErrShardsNotAdjacent = errors.New("shards are not adjacent")

	// ErrOperationNotSupportedOnOnDemandStream is returned when the operation
	// is not supported on an on-demand stream.
	ErrOperationNotSupportedOnOnDemandStream = errors.New("operation not supported on on-demand stream")

	// ErrInvalidShardCount is returned when the shard count is not valid for
	// the scaling operation.
	ErrInvalidShardCount = errors.New("invalid shard count for scaling operation")

	// ErrConsumerNotFound is returned when the specified consumer
	// does not exist.
	ErrConsumerNotFound = errors.New("consumer not found")

	// ErrConsumerAlreadyExists is returned when attempting to create a consumer
	// that already exists.
	ErrConsumerAlreadyExists = errors.New("consumer already exists")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = common.ErrInvalidParameter

	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = errors.New("resource not found")
)

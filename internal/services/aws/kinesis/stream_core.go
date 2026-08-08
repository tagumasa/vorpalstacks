package kinesis

import (
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// AdminCreateStreamInput is the transport-agnostic input for the admin
// console CreateStream operation.
type AdminCreateStreamInput struct {
	StreamName string
	ShardCount int32
}

// AdminListStreamsInput is the transport-agnostic input for the admin
// console ListStreams operation.
type AdminListStreamsInput struct {
	ExclusiveStartStreamName string
	Limit                    int
}

// AdminDescribeStreamInput is the transport-agnostic input for the admin
// console DescribeStream operation.
type AdminDescribeStreamInput struct {
	StreamName string
	StreamARN  string
}

// AdminDeleteStreamInput is the transport-agnostic input for the admin
// console DeleteStream operation.
type AdminDeleteStreamInput struct {
	StreamName string
}

// ListStreamsResult contains the result of a listStreamsCore call.
// Store types are used directly so that the admin handler convert file
// can translate them to proto types.
type ListStreamsResult struct {
	Streams     []*kinesisstore.Stream
	NextMarker  string
	IsTruncated bool
}

// DescribeStreamResult contains the result of a describeStreamCore call.
type DescribeStreamResult struct {
	Stream *kinesisstore.Stream
	Shards []*kinesisstore.Shard
}

// createStreamCore creates a new Kinesis stream. Used by the admin console
// gRPC-Web handler.
func (s *KinesisService) createStreamCore(store *kinesisstore.KinesisStore, input AdminCreateStreamInput) error {
	if !validateStreamName(input.StreamName) {
		return ErrInvalidArgument
	}

	shardCount := input.ShardCount
	if shardCount == 0 {
		shardCount = 1
	}
	if !validateShardCount(shardCount) {
		return ErrInvalidArgument
	}

	_, err := store.CreateStream(input.StreamName, shardCount, kinesisstore.StreamModeProvisioned, 0, 0)
	if err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// deleteStreamCore deletes a Kinesis stream. Used by the admin console
// gRPC-Web handler.
func (s *KinesisService) deleteStreamCore(store *kinesisstore.KinesisStore, input AdminDeleteStreamInput) error {
	if !validateStreamName(input.StreamName) {
		return ErrInvalidArgument
	}

	if err := store.DeleteStream(input.StreamName); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// listStreamsCore lists Kinesis streams with pagination. Used by the admin
// console gRPC-Web handler.
func (s *KinesisService) listStreamsCore(store *kinesisstore.KinesisStore, input AdminListStreamsInput) (ListStreamsResult, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 100
	}
	if !validateListStreamsLimit(limit) {
		return ListStreamsResult{}, ErrInvalidArgument
	}

	result, err := store.ListStreams(storecommon.ListOptions{
		Marker:   input.ExclusiveStartStreamName,
		MaxItems: limit,
	})
	if err != nil {
		return ListStreamsResult{}, s.mapStoreError(err)
	}

	return ListStreamsResult{
		Streams:     result.Items,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

// describeStreamCore describes a Kinesis stream. Used by the admin console
// gRPC-Web handler.
func (s *KinesisService) describeStreamCore(store *kinesisstore.KinesisStore, input AdminDescribeStreamInput) (DescribeStreamResult, error) {
	streamName := input.StreamName

	if streamName == "" && input.StreamARN != "" {
		stream, err := store.GetStreamByARN(input.StreamARN)
		if err != nil {
			return DescribeStreamResult{}, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if !validateStreamName(streamName) {
		return DescribeStreamResult{}, ErrInvalidArgument
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return DescribeStreamResult{}, s.mapStoreError(err)
	}

	shards, err := store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return DescribeStreamResult{}, s.mapStoreError(err)
	}

	return DescribeStreamResult{
		Stream: stream,
		Shards: shards,
	}, nil
}

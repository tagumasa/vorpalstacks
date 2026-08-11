package kinesis

import (
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateStreamInput is the transport-agnostic input for CreateStream.
type CreateStreamInput struct {
	StreamName             string
	ShardCount             int32
	StreamMode             kinesisstore.StreamMode
	MaxRecordSizeInKiB     int32
	HasMaxRecordSizeInKiB  bool
	WarmThroughputMiBps    int32
	HasWarmThroughputMiBps bool
	Tags                   []types.Tag
}

// DeleteStreamInput is the transport-agnostic input for DeleteStream.
type DeleteStreamInput struct {
	StreamName string
	StreamARN  string
}

// ListStreamsInput is the transport-agnostic input for ListStreams.
type ListStreamsInput struct {
	ExclusiveStartStreamName string
	Limit                    int
	NextToken                string
}

// DescribeStreamInput is the transport-agnostic input for DescribeStream.
type DescribeStreamInput struct {
	StreamName string
	StreamARN  string
}

// ListStreamsResult contains the result of a listStreamsCore call.
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

// createStreamCore is the single entry point for creating a Kinesis stream,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) createStreamCore(store *kinesisstore.KinesisStore, input CreateStreamInput) (*kinesisstore.Stream, error) {
	if !validateStreamName(input.StreamName) {
		return nil, ErrInvalidArgument
	}

	shardCount := input.ShardCount
	if shardCount == 0 {
		shardCount = 1
	}
	if !validateShardCount(shardCount) {
		return nil, ErrInvalidArgument
	}

	streamMode := input.StreamMode
	if streamMode == "" {
		streamMode = kinesisstore.StreamModeProvisioned
	}

	if input.HasMaxRecordSizeInKiB && !validateMaxRecordSizeInKiB(input.MaxRecordSizeInKiB) {
		return nil, ErrInvalidArgument
	}

	if input.HasWarmThroughputMiBps && !validateWarmThroughputMiBps(input.WarmThroughputMiBps) {
		return nil, ErrInvalidArgument
	}

	stream, err := store.CreateStream(input.StreamName, shardCount, streamMode, input.MaxRecordSizeInKiB, input.WarmThroughputMiBps)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if len(input.Tags) > 0 {
		tagMap := make(map[string]string, len(input.Tags))
		for _, t := range input.Tags {
			tagMap[t.Key] = t.Value
		}
		if err := store.Tag(input.StreamName, tagMap); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	return stream, nil
}

// deleteStreamCore is the single entry point for deleting a Kinesis stream,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) deleteStreamCore(store *kinesisstore.KinesisStore, input DeleteStreamInput) error {
	streamName := input.StreamName

	if input.StreamARN != "" {
		stream, err := store.GetStreamByARN(input.StreamARN)
		if err != nil {
			return s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if !validateStreamName(streamName) {
		return ErrInvalidArgument
	}

	if err := store.DeleteStream(streamName); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// listStreamsCore is the single entry point for listing Kinesis streams,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) listStreamsCore(store *kinesisstore.KinesisStore, input ListStreamsInput) (ListStreamsResult, error) {
	exclusiveStartName := input.ExclusiveStartStreamName
	if exclusiveStartName == "" && input.NextToken != "" {
		exclusiveStartName = input.NextToken
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 100
	}
	if !validateListStreamsLimit(limit) {
		return ListStreamsResult{}, ErrInvalidArgument
	}

	result, err := store.ListStreams(storecommon.ListOptions{
		Marker:   exclusiveStartName,
		MaxItems: limit + 1,
	})
	if err != nil {
		return ListStreamsResult{}, s.mapStoreError(err)
	}

	hasMore := len(result.Items) > limit
	if hasMore {
		result.Items = result.Items[:limit]
	}

	nextMarker := ""
	if hasMore && len(result.Items) > 0 {
		nextMarker = result.Items[len(result.Items)-1].StreamName
	}

	return ListStreamsResult{
		Streams:     result.Items,
		NextMarker:  nextMarker,
		IsTruncated: hasMore,
	}, nil
}

// describeStreamCore is the single entry point for describing a Kinesis
// stream, shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) describeStreamCore(store *kinesisstore.KinesisStore, input DescribeStreamInput) (DescribeStreamResult, error) {
	streamName := input.StreamName

	if input.StreamARN != "" {
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

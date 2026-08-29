package kinesis

import (
	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
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
	HasLimit                 bool
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

// DescribeStreamSummaryInput is the transport-agnostic input for
// DescribeStreamSummary.
type DescribeStreamSummaryInput struct {
	StreamName string
	StreamARN  string
}

// DescribeStreamSummaryResult contains the result of a
// describeStreamSummaryCore call.
type DescribeStreamSummaryResult struct {
	Stream *kinesisstore.Stream
}

// UpdateStreamModeInput is the transport-agnostic input for UpdateStreamMode.
type UpdateStreamModeInput struct {
	StreamARN           string
	StreamMode          kinesisstore.StreamMode
	WarmThroughputMiBps int32
	HasWarmThroughput   bool
}

// UpdateStreamModeResult contains the result of an updateStreamModeCore call.
type UpdateStreamModeResult struct {
	StreamARN string
}

// resolveStreamNameCore resolves a stream name from the StreamName/StreamARN
// pair (the ARN wins), requiring a valid name.
func (s *KinesisService) resolveStreamNameCore(store *kinesisstore.KinesisStore, streamName, streamARN string) (string, error) {
	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return "", s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return "", ErrInvalidArgument
	}

	if !validateStreamName(streamName) {
		return "", ErrInvalidArgument
	}

	return streamName, nil
}

// resolveStreamNameOptionalCore resolves a stream name from the
// StreamName/StreamARN pair but does not require one (used by ListShards
// which accepts optional stream identification).
func (s *KinesisService) resolveStreamNameOptionalCore(store *kinesisstore.KinesisStore, streamName, streamARN string) (string, error) {
	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return "", s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName != "" && !validateStreamName(streamName) {
		return "", ErrInvalidArgument
	}

	return streamName, nil
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

// The Limit input window and the effective page follow two different
// documented rules: the ListStreamsInputLimit range trait accepts 1-10000
// on the wire, while the Limit member documentation fixes the default page
// at 100 and returns at most 100 results per call.
const (
	DefaultListStreamsResults = 100
	MaxListStreamsResults     = 100
)

// listStreamsCore is the single entry point for listing Kinesis streams,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) listStreamsCore(store *kinesisstore.KinesisStore, input ListStreamsInput) (ListStreamsResult, error) {
	exclusiveStartName := input.ExclusiveStartStreamName
	if exclusiveStartName == "" && input.NextToken != "" {
		exclusiveStartName = input.NextToken
	}

	limit := input.Limit
	if input.HasLimit {
		if !validateListStreamsLimit(limit) {
			return ListStreamsResult{}, ErrInvalidArgument
		}
		if limit > MaxListStreamsResults {
			limit = MaxListStreamsResults
		}
	} else {
		limit = DefaultListStreamsResults
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

// describeStreamSummaryCore is the single entry point for describing a
// Kinesis stream summary. It serves the HTTP DescribeStreamSummary
// operation.
func (s *KinesisService) describeStreamSummaryCore(store *kinesisstore.KinesisStore, input DescribeStreamSummaryInput) (DescribeStreamSummaryResult, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return DescribeStreamSummaryResult{}, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return DescribeStreamSummaryResult{}, s.mapStoreError(err)
	}

	return DescribeStreamSummaryResult{Stream: stream}, nil
}

// updateStreamModeCore switches a stream between provisioned and on-demand
// modes, optionally carrying a warm-throughput target.
func (s *KinesisService) updateStreamModeCore(reqCtx *request.RequestContext, input UpdateStreamModeInput) (UpdateStreamModeResult, error) {
	if input.StreamARN == "" || !validateStreamMode(string(input.StreamMode)) {
		return UpdateStreamModeResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return UpdateStreamModeResult{}, err
	}

	stream, err := store.GetStreamByARN(input.StreamARN)
	if err != nil {
		return UpdateStreamModeResult{}, s.mapStoreError(err)
	}

	stream.StreamModeDetails = &kinesisstore.StreamModeDetails{StreamMode: input.StreamMode}

	if input.HasWarmThroughput {
		if !validateWarmThroughputMiBps(input.WarmThroughputMiBps) {
			return UpdateStreamModeResult{}, ErrInvalidArgument
		}
		stream.WarmThroughputMiBps = input.WarmThroughputMiBps
	}

	if err := store.UpdateStream(stream); err != nil {
		return UpdateStreamModeResult{}, s.mapStoreError(err)
	}

	return UpdateStreamModeResult{StreamARN: stream.StreamARN}, nil
}

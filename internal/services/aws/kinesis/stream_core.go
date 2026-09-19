package kinesis

import (
	"encoding/base64"

	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// CreateStreamInput is the transport-agnostic input for CreateStream.
type CreateStreamInput struct {
	StreamName             string
	ShardCount             int32
	HasShardCount          bool
	StreamMode             kinesisstore.StreamMode
	HasStreamModeDetails   bool
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

	// Limit is the requested shard page (range trait 1-10000; the
	// member's documentation caps the effective page at one hundred)
	// and ExclusiveStartShardId resumes the page after that shard.
	Limit                 int
	HasLimit              bool
	ExclusiveStartShardId string
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

	// HasMoreShards reports whether further shard pages follow — the
	// StreamDescription member the model marks required.
	HasMoreShards bool
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

// createStreamCore is the single entry point for creating a Kinesis stream,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) createStreamCore(reqCtx *request.RequestContext, input CreateStreamInput) (*kinesisstore.Stream, error) {
	if !validateStreamName(input.StreamName) {
		return nil, ErrInvalidArgument
	}

	// ShardCount targets the PositiveIntegerObject shape (range min 1):
	// only an absent member takes the single-shard default — an explicit
	// zero is an out-of-range value the range trait rejects, not a value
	// to substitute.
	shardCount := input.ShardCount
	if !input.HasShardCount {
		shardCount = 1
	}
	if !validateShardCount(shardCount) {
		return nil, ErrInvalidArgument
	}

	// The StreamModeDetails member is optional: an absent one defaults to
	// provisioned. A present one must carry its own StreamMode member and a
	// value from the StreamMode enum.
	streamMode := input.StreamMode
	if input.HasStreamModeDetails && streamMode == "" {
		return nil, ErrInvalidArgument
	}
	if streamMode == "" {
		streamMode = kinesisstore.StreamModeProvisioned
	}
	if !validateStreamMode(string(streamMode)) {
		return nil, ErrInvalidArgument
	}

	if input.HasMaxRecordSizeInKiB && !validateMaxRecordSizeInKiB(input.MaxRecordSizeInKiB) {
		return nil, ErrInvalidArgument
	}

	if input.HasWarmThroughputMiBps && !validateWarmThroughputMiBps(input.WarmThroughputMiBps) {
		return nil, ErrInvalidArgument
	}

	// Create-time tags pass the same shared tag-set validation every tag
	// write path applies, before the stream is created.
	if len(input.Tags) > 0 {
		if err := types.ValidateTags(input.Tags); err != nil {
			return nil, ErrInvalidArgument
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var tagMap map[string]string
	if len(input.Tags) > 0 {
		tagMap = make(map[string]string, len(input.Tags))
		for _, t := range input.Tags {
			tagMap[t.Key] = t.Value
		}
	}

	stream, err := store.CreateStream(input.StreamName, shardCount, streamMode, input.MaxRecordSizeInKiB, input.WarmThroughputMiBps, tagMap)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return stream, nil
}

// deleteStreamCore is the single entry point for deleting a Kinesis stream,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) deleteStreamCore(reqCtx *request.RequestContext, input DeleteStreamInput) error {
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return err
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
func (s *KinesisService) listStreamsCore(reqCtx *request.RequestContext, input ListStreamsInput) (ListStreamsResult, error) {
	exclusiveStartName := input.ExclusiveStartStreamName
	if exclusiveStartName == "" && input.NextToken != "" {
		// The resumption token is an opaque envelope over the page's last
		// stream name; one that cannot be decoded is expired — the error
		// ListStreams declares for it — never a silent restart at page one.
		decoded, err := base64.StdEncoding.DecodeString(input.NextToken)
		if err != nil {
			return ListStreamsResult{}, ErrExpiredNextToken
		}
		exclusiveStartName = string(decoded)
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

	store, err := s.store(reqCtx)
	if err != nil {
		return ListStreamsResult{}, err
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
		nextMarker = base64.StdEncoding.EncodeToString([]byte(result.Items[len(result.Items)-1].StreamName))
	}

	return ListStreamsResult{
		Streams:     result.Items,
		NextMarker:  nextMarker,
		IsTruncated: hasMore,
	}, nil
}

// describeStreamCore is the single entry point for describing a Kinesis
// stream, shared by the HTTP API and the admin gRPC-Web handler.
func (s *KinesisService) describeStreamCore(reqCtx *request.RequestContext, input DescribeStreamInput) (DescribeStreamResult, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return DescribeStreamResult{}, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return DescribeStreamResult{}, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return DescribeStreamResult{}, s.mapStoreError(err)
	}

	// The Limit member's range trait accepts 1-10000, while its
	// documentation fixes the effective page: the default is one hundred
	// and "if you specify a value greater than 100, at most 100 results
	// are returned". One over-fetch past the effective page decides
	// HasMoreShards, so the exact final page reports false.
	if input.HasLimit && (input.Limit < 1 || input.Limit > kinesisstore.MaxListResultsLimit) {
		return DescribeStreamResult{}, ErrInvalidArgument
	}
	limit := kinesisstore.DescribeStreamShardPageCeiling
	if input.HasLimit && input.Limit < limit {
		limit = input.Limit
	}

	shards, err := store.ListShards(streamName, nil, input.ExclusiveStartShardId, limit+1)
	if err != nil {
		return DescribeStreamResult{}, s.mapStoreError(err)
	}

	hasMoreShards := len(shards) > limit
	if hasMoreShards {
		shards = shards[:limit]
	}

	return DescribeStreamResult{
		Stream:        stream,
		Shards:        shards,
		HasMoreShards: hasMoreShards,
	}, nil
}

// describeStreamSummaryCore is the single entry point for describing a
// Kinesis stream summary. It serves the HTTP DescribeStreamSummary
// operation.
func (s *KinesisService) describeStreamSummaryCore(reqCtx *request.RequestContext, input DescribeStreamSummaryInput) (DescribeStreamSummaryResult, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return DescribeStreamSummaryResult{}, err
	}

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
// modes, optionally carrying a warm-throughput target. StreamModeDetails is
// a required member whose inner StreamMode member is required in turn: an
// absent member or an inner-less one both read as an empty mode, which the
// enum check rejects. The operation's output is Unit — no result to carry.
func (s *KinesisService) updateStreamModeCore(reqCtx *request.RequestContext, input UpdateStreamModeInput) error {
	if input.StreamARN == "" || !validateStreamMode(string(input.StreamMode)) {
		return ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	stream, err := store.GetStreamByARN(input.StreamARN)
	if err != nil {
		return s.mapStoreError(err)
	}

	if input.HasWarmThroughput && !validateWarmThroughputMiBps(input.WarmThroughputMiBps) {
		return ErrInvalidArgument
	}

	if _, err := store.UpdateStreamFields(stream.StreamName, func(stream *kinesisstore.Stream) error {
		stream.StreamModeDetails = &kinesisstore.StreamModeDetails{StreamMode: input.StreamMode}
		if input.HasWarmThroughput {
			stream.WarmThroughputMiBps = input.WarmThroughputMiBps
		}
		return nil
	}); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

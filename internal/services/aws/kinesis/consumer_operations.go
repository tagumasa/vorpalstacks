package kinesis

import (
	"context"
	"io"
	"net/http"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
)

// SubscribeToShardEventStream represents an event stream for SubscribeToShard operations.
type SubscribeToShardEventStream struct {
	reader io.Reader
}

// GetStream returns the reader for the stream.
func (e *SubscribeToShardEventStream) GetStream() io.Reader {
	return e.reader
}

// GetStreamHeaders returns the HTTP headers for the stream.
func (e *SubscribeToShardEventStream) GetStreamHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/vnd.amazon.eventstream")
	return headers
}

// RegisterStreamConsumer registers a consumer for a Kinesis stream.
func (s *KinesisService) RegisterStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.registerStreamConsumerCore(reqCtx, RegisterStreamConsumerInput{
		StreamARN:    request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ConsumerName: request.GetParamLowerFirst(req.Parameters, "ConsumerName"),
		Tags:         tags.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Consumer": formatConsumer(result.Consumer),
	}, nil
}

// DeregisterStreamConsumer deregisters a consumer from a Kinesis stream.
func (s *KinesisService) DeregisterStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deregisterStreamConsumerCore(reqCtx, DeregisterStreamConsumerInput{
		StreamARN:    request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ConsumerARN:  request.GetParamLowerFirst(req.Parameters, "ConsumerARN"),
		ConsumerName: request.GetParamLowerFirst(req.Parameters, "ConsumerName"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeStreamConsumer returns details about a Kinesis stream consumer.
func (s *KinesisService) DescribeStreamConsumer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.describeStreamConsumerCore(reqCtx, DescribeStreamConsumerInput{
		StreamARN:    request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ConsumerARN:  request.GetParamLowerFirst(req.Parameters, "ConsumerARN"),
		ConsumerName: request.GetParamLowerFirst(req.Parameters, "ConsumerName"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ConsumerDescription": formatConsumerDescription(result.Consumer),
	}, nil
}

// ListStreamConsumers lists consumers of a Kinesis stream.
func (s *KinesisService) ListStreamConsumers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults, err := strictIntParam(req.Parameters, "MaxResults")
	if err != nil {
		return nil, err
	}

	// The SDKs serialise StreamCreationTimestamp as an epoch-second JSON
	// number; the strict reader carries both that form and the documented
	// string notations, so the Core's generation check actually operates.
	streamCreationTimestamp, _, err := strictTimestampParam(req.Parameters, "StreamCreationTimestamp")
	if err != nil {
		return nil, err
	}

	result, err := s.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:               request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		StreamCreationTimestamp: streamCreationTimestamp,
		MaxResults:              maxResults,
		HasMaxResults:           hasMaxResults,
		NextToken:               request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	formattedConsumers := make([]map[string]interface{}, len(result.Consumers))
	for i, c := range result.Consumers {
		formattedConsumers[i] = formatConsumer(c)
	}

	resp := map[string]interface{}{
		"Consumers": formattedConsumers,
	}
	// The output shape carries Consumers and NextToken only; the
	// continuation token is present when a page follows.
	if result.NextToken != nil {
		resp["NextToken"] = *result.NextToken
	}
	return resp, nil
}

// SubscribeToShard subscribes a consumer to receive records from a Kinesis shard.
func (s *KinesisService) SubscribeToShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	startingPosition := ""
	seqNum := ""
	tsStr := ""
	if sp := request.GetMapParam(req.Parameters, "StartingPosition"); sp != nil {
		startingPosition = request.GetStringParam(sp, "Type")
		seqNum = request.GetStringParam(sp, "SequenceNumber")
		// The SDKs serialise the position's Timestamp as an epoch-second
		// JSON number; the normaliser keeps it from being dropped the way
		// a string-only read would.
		if raw, ok := sp["Timestamp"]; ok && raw != nil {
			s, err := request.NormalizeTimestampValue(raw)
			if err != nil {
				return nil, ErrInvalidArgument
			}
			tsStr = s
		}
	}

	return s.subscribeToShardCore(ctx, reqCtx, SubscribeToShardInput{
		ConsumerARN:            request.GetParamLowerFirst(req.Parameters, "ConsumerARN"),
		ShardId:                request.GetParamLowerFirst(req.Parameters, "ShardId"),
		StartingPositionType:   startingPosition,
		StartingSequenceNumber: seqNum,
		Timestamp:              tsStr,
		DryRun:                 request.GetBoolParam(req.Parameters, "DryRun"),
	})
}

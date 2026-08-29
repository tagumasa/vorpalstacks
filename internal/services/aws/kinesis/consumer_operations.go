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
		"ConsumerDescription": formatConsumer(result.Consumer),
	}, nil
}

// ListStreamConsumers lists consumers of a Kinesis stream.
func (s *KinesisService) ListStreamConsumers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults := 0
	hasMaxResults := false
	if _, ok := req.Parameters["MaxResults"]; ok {
		maxResults = request.GetIntParam(req.Parameters, "MaxResults")
		hasMaxResults = true
	}

	result, err := s.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:               request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		StreamCreationTimestamp: request.GetParamLowerFirst(req.Parameters, "StreamCreationTimestamp"),
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

	return map[string]interface{}{
		"Consumers":        formattedConsumers,
		"HasMoreConsumers": result.HasMore,
		"NextToken":        result.NextToken,
	}, nil
}

// SubscribeToShard subscribes a consumer to receive records from a Kinesis shard.
func (s *KinesisService) SubscribeToShard(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	startingPosition := ""
	seqNum := ""
	tsStr := ""
	if sp := request.GetMapParam(req.Parameters, "StartingPosition"); sp != nil {
		startingPosition = request.GetStringParam(sp, "Type")
		seqNum = request.GetStringParam(sp, "SequenceNumber")
		tsStr = request.GetStringParam(sp, "Timestamp")
	}

	return s.subscribeToShardCore(ctx, reqCtx, SubscribeToShardInput{
		ConsumerARN:            request.GetParamLowerFirst(req.Parameters, "ConsumerARN"),
		ShardId:                request.GetParamLowerFirst(req.Parameters, "ShardId"),
		StartingPositionType:   startingPosition,
		StartingSequenceNumber: seqNum,
		Timestamp:              tsStr,
	})
}

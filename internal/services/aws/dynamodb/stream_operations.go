package dynamodb

import (
	"context"
	"math"

	"vorpalstacks/internal/common/request"
)

// DescribeStream returns information about a DynamoDB stream.
//
// AWS API: DynamoDB Streams — DescribeStream
// Protocol: JSON (X-Amz-Target: DynamoDBStreams_20120810.DescribeStream)
func (s *DynamoDBService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := DescribeStreamInput{
		StreamArn:             request.GetStringParam(req.Parameters, "StreamArn"),
		ExclusiveStartShardId: request.GetStringParam(req.Parameters, "ExclusiveStartShardId"),
	}
	limit, limitPresent, err := intParamWithPresence(req.Parameters, "Limit")
	if err != nil {
		return nil, err
	}
	in.Limit, in.LimitSet = limit, limitPresent
	if filterMap, ok := req.Parameters["ShardFilter"].(map[string]interface{}); ok {
		in.ShardFilter = &ShardFilterSpec{
			Type:    request.GetStringParam(filterMap, "Type"),
			ShardId: request.GetStringParam(filterMap, "ShardId"),
		}
	}

	result, err := s.describeStreamCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamDescription": buildStreamDescriptionResponse(result),
	}, nil
}

// GetShardIterator returns a shard iterator positioned according to the
// requested iterator type.
//
// AWS API: DynamoDB Streams — GetShardIterator
func (s *DynamoDBService) GetShardIterator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getShardIteratorCore(store,
		request.GetStringParam(req.Parameters, "StreamArn"),
		request.GetStringParam(req.Parameters, "ShardId"),
		request.GetStringParam(req.Parameters, "ShardIteratorType"),
		request.GetStringParam(req.Parameters, "SequenceNumber"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ShardIterator": result.ShardIterator,
	}, nil
}

// GetRecords retrieves stream records from the given shard iterator
// position.
//
// AWS API: DynamoDB Streams — GetRecords
func (s *DynamoDBService) GetRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, limitSet, err := intParamWithPresence(req.Parameters, "Limit")
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getRecordsCore(store, request.GetStringParam(req.Parameters, "ShardIterator"), limit, limitSet)
	if err != nil {
		return nil, err
	}

	records := result.Records
	if records == nil {
		records = []interface{}{}
	}

	return map[string]interface{}{
		"Records":           records,
		"NextShardIterator": result.NextShardIterator,
	}, nil
}

// ListStreams returns stream ARNs associated with the current account and
// endpoint.
//
// AWS API: DynamoDB Streams — ListStreams
func (s *DynamoDBService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, limitSet, err := intParamWithPresence(req.Parameters, "Limit")
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listStreamsCore(store,
		request.GetStringParam(req.Parameters, "TableName"),
		request.GetStringParam(req.Parameters, "ExclusiveStartStreamArn"),
		limit, limitSet)
	if err != nil {
		return nil, err
	}

	streamList := make([]map[string]interface{}, 0, len(result.Streams))
	for _, st := range result.Streams {
		streamList = append(streamList, map[string]interface{}{
			"StreamArn":   st.StreamArn,
			"TableName":   st.TableName,
			"StreamLabel": st.StreamLabel,
		})
	}

	resp := map[string]interface{}{
		"Streams": streamList,
	}
	if result.LastEvaluatedStreamArn != "" {
		resp["LastEvaluatedStreamArn"] = result.LastEvaluatedStreamArn
	}
	return resp, nil
}

// intParamWithPresence extracts an integer-valued member and reports whether
// the member was present, so an explicit zero stays distinguishable from an
// omitted value — the model's PositiveIntegerObject range rejects the
// former and defaults the latter. A present member of a non-numeric type is
// a validation error, and so is a JSON number that is not integral or that
// falls outside the platform's integer range: a wire number is never
// silently truncated or overflowed into an integer.
func intParamWithPresence(params map[string]interface{}, name string) (value int, present bool, err error) {
	raw, ok := params[name]
	if !ok {
		return 0, false, nil
	}
	switch v := raw.(type) {
	case float64:
		if v != math.Trunc(v) || v < math.MinInt64 || v > math.MaxInt64 {
			return 0, true, ErrInvalidParameter
		}
		return int(v), true, nil
	case int:
		return v, true, nil
	default:
		return 0, true, ErrInvalidParameter
	}
}

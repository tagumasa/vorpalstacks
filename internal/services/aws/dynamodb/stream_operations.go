package dynamodb

import (
	"context"

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

	result, err := s.describeStreamCore(store, request.GetStringParam(req.Parameters, "StreamArn"))
	if err != nil {
		return nil, err
	}

	shards := make([]interface{}, 0, len(result.Shards))
	for _, sh := range result.Shards {
		seqRange := map[string]interface{}{
			"StartingSequenceNumber": sh.StartingSequenceNumber,
		}
		if sh.EndingSequenceNumber != "" {
			seqRange["EndingSequenceNumber"] = sh.EndingSequenceNumber
		}
		shards = append(shards, map[string]interface{}{
			"ShardId":             sh.ShardID,
			"SequenceNumberRange": seqRange,
		})
	}

	return map[string]interface{}{
		"StreamDescription": map[string]interface{}{
			"StreamArn":               result.StreamArn,
			"StreamLabel":             result.StreamLabel,
			"StreamStatus":            result.StreamStatus,
			"StreamViewType":          result.StreamViewType,
			"TableName":               result.TableName,
			"KeySchema":               buildKeySchemaResponse(result.KeySchema),
			"Shards":                  shards,
			"CreationRequestDateTime": result.CreationRequestDateTime,
		},
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
	limit := 0
	limitSet := false
	if limitVal, ok := req.Parameters["Limit"]; ok {
		switch v := limitVal.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		default:
			return nil, ErrInvalidParameter
		}
		limitSet = true
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listStreamsCore(store,
		request.GetStringParam(req.Parameters, "TableName"),
		request.GetStringParam(req.Parameters, "ExclusiveStartStreamArn"),
		request.GetIntParam(req.Parameters, "Limit"))
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

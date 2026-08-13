package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
)

// DescribeStream returns information about a DynamoDB stream.
//
// AWS API: DynamoDB Streams — DescribeStream
// Protocol: JSON (X-Amz-Target: DynamoDBStreams_20120810.DescribeStream)
func (s *DynamoDBService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamArn := request.GetStringParam(req.Parameters, "StreamArn")
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeStreamCore(store, streamArn)
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
	streamArn := request.GetStringParam(req.Parameters, "StreamArn")
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}
	shardId := request.GetStringParam(req.Parameters, "ShardId")
	if shardId == "" {
		return nil, ErrInvalidParameter
	}
	iteratorType := request.GetStringParam(req.Parameters, "ShardIteratorType")
	if iteratorType == "" {
		return nil, ErrInvalidParameter
	}
	sequenceNumber := request.GetStringParam(req.Parameters, "SequenceNumber")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getShardIteratorCore(store, streamArn, iteratorType, sequenceNumber)
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
	iterator := request.GetStringParam(req.Parameters, "ShardIterator")
	if iterator == "" {
		return nil, ErrInvalidParameter
	}

	limit := getRecordsDefaultLimit
	if limitVal, ok := req.Parameters["Limit"]; ok {
		switch v := limitVal.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		default:
			return nil, ErrInvalidParameter
		}
	}
	if limit < 0 || limit > getRecordsMaxLimit {
		return nil, ErrInvalidParameter
	}

	tableName, fromSeq, err := decodeShardIterator(iterator)
	if err != nil {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getRecordsCore(store, tableName, fromSeq, limit)
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

// extractTableNameFromStreamArn parses a DynamoDB stream ARN to extract
// the table name.
func extractTableNameFromStreamArn(streamArn string) string {
	idx := strings.Index(streamArn, "table/")
	if idx < 0 {
		return ""
	}
	rest := streamArn[idx+6:]
	slashIdx := strings.Index(rest, "/")
	if slashIdx < 0 {
		return rest
	}
	return rest[:slashIdx]
}

// encodeShardIterator creates an opaque iterator string from the table
// name and sequence number. Format: "tableName|seqNum".
func encodeShardIterator(tableName string, seq int64) string {
	return fmt.Sprintf("%s|%d", tableName, seq)
}

// decodeShardIterator parses an iterator string back into table name and
// sequence number.
func decodeShardIterator(iterator string) (string, int64, error) {
	idx := strings.LastIndex(iterator, "|")
	if idx < 0 {
		return "", 0, fmt.Errorf("invalid iterator format")
	}
	tableName := iterator[:idx]
	seqStr := iterator[idx+1:]
	seq, err := strconv.ParseInt(seqStr, 10, 64)
	if err != nil {
		return "", 0, err
	}
	return tableName, seq, nil
}

// streamTimeNow returns the current time. Extracted for potential testing.
var streamTimeNow = func() time.Time { return time.Now().UTC() }

// ListStreams returns stream ARNs associated with the current account and
// endpoint.
//
// AWS API: DynamoDB Streams — ListStreams
func (s *DynamoDBService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableNameFilter := request.GetStringParam(req.Parameters, "TableName")
	exclusiveStartStreamArn := request.GetStringParam(req.Parameters, "ExclusiveStartStreamArn")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit == 0 {
		limit = listStreamsDefaultLimit
	}
	if limit > listStreamsMaxLimit {
		limit = listStreamsMaxLimit
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listStreamsCore(store, tableNameFilter, exclusiveStartStreamArn, limit)
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

package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DescribeStream returns information about a DynamoDB stream. The stream
// ARN is extracted from the StreamArn parameter (which may be the table's
// StreamArn or a full stream ARN).
//
// AWS API: DynamoDB Streams — DescribeStream
// Protocol: JSON (X-Amz-Target: DynamoDBStreams_20120810.DescribeStream)
func (s *DynamoDBService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamArn := request.GetStringParam(req.Parameters, "StreamArn")
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}

	tableName := extractTableNameFromStreamArn(streamArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.GetCachedStoreForRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	table, err := store.Tables().Get(tableName)
	if err != nil || table == nil || table.StreamArn != streamArn {
		return nil, ErrResourceNotFound
	}

	latestSeq, _ := store.Streams().GetLatestSequence(tableName)

	seqRange := map[string]interface{}{
		"StartingSequenceNumber": fmt.Sprintf("%d", int64(1)),
	}
	if latestSeq >= 1 {
		seqRange["EndingSequenceNumber"] = fmt.Sprintf("%d", latestSeq)
	}

	shard := map[string]interface{}{
		"ShardId":             dbstore.ShardID,
		"SequenceNumberRange": seqRange,
	}

	return map[string]interface{}{
		"StreamDescription": map[string]interface{}{
			"StreamArn":               streamArn,
			"StreamLabel":             table.LatestStreamLabel,
			"StreamStatus":            "ENABLED",
			"StreamViewType":          string(table.StreamSpecification.StreamViewType),
			"TableName":               tableName,
			"KeySchema":               buildKeySchemaResponse(table.KeySchema),
			"Shards":                  []interface{}{shard},
			"CreationRequestDateTime": table.CreationDateTime.Unix(),
		},
	}, nil
}

// GetShardIterator returns a shard iterator positioned according to the
// requested iterator type:
//
//	TRIM_HORIZON     — from the beginning (sequence 0)
//	LATEST           — from the latest position
//	AT_SEQUENCE_NUMBER   — at the given sequence number
//	AFTER_SEQUENCE_NUMBER — after the given sequence number
//
// The iterator is a string encoding of the starting sequence number.
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

	tableName := extractTableNameFromStreamArn(streamArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.GetCachedStoreForRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	// Validate the table exists and the stream ARN matches; AWS returns
	// ResourceNotFoundException for either.
	table, err := store.Tables().Get(tableName)
	if err != nil || table == nil || table.StreamArn != streamArn {
		return nil, ErrResourceNotFound
	}

	latestSeq, err := store.Streams().GetLatestSequence(tableName)
	if err != nil {
		return nil, err
	}

	var startSeq int64
	switch iteratorType {
	case "TRIM_HORIZON":
		startSeq = 0
	case "LATEST":
		startSeq = latestSeq
	case "AT_SEQUENCE_NUMBER":
		seqStr := request.GetStringParam(req.Parameters, "SequenceNumber")
		startSeq, err = strconv.ParseInt(seqStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		startSeq-- // AT means inclusive, so iterator starts at seq-1 (GetRecords reads from seq+1)
	case "AFTER_SEQUENCE_NUMBER":
		seqStr := request.GetStringParam(req.Parameters, "SequenceNumber")
		startSeq, err = strconv.ParseInt(seqStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidParameter
		}
	default:
		return nil, ErrInvalidParameter
	}

	iterator := encodeShardIterator(tableName, startSeq)

	return map[string]interface{}{
		"ShardIterator": iterator,
	}, nil
}

// GetRecords retrieves stream records from the given shard iterator
// position. Returns up to Limit records (default 100, max 1000) and a
// new iterator for the next page. If no records are available, returns
// an empty list and the same iterator position.
//
// AWS API: DynamoDB Streams — GetRecords
func (s *DynamoDBService) GetRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	iterator := request.GetStringParam(req.Parameters, "ShardIterator")
	if iterator == "" {
		return nil, ErrInvalidParameter
	}

	limit := 100
	if limitVal, ok := req.Parameters["Limit"]; ok {
		switch v := limitVal.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		}
	}
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	tableName, fromSeq, err := decodeShardIterator(iterator)
	if err != nil {
		return nil, ErrInvalidParameter
	}

	store, err := s.GetCachedStoreForRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	// Validate the table exists and has streaming enabled.
	table, err := store.Tables().Get(tableName)
	if err != nil || table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return nil, ErrResourceNotFound
	}

	records, nextSeq, err := store.Streams().GetRecords(tableName, fromSeq, limit)
	if err != nil {
		logs.Warn("failed to get stream records",
			logs.String("table", tableName), logs.Err(err))
		return nil, ErrInternal
	}

	nextIterator := encodeShardIterator(tableName, nextSeq)

	var recordsResp []interface{}
	for _, rec := range records {
		recordsResp = append(recordsResp, rec)
	}
	if recordsResp == nil {
		recordsResp = []interface{}{}
	}

	return map[string]interface{}{
		"Records":           recordsResp,
		"NextShardIterator": nextIterator,
	}, nil
}

// extractTableNameFromStreamArn parses a DynamoDB stream ARN to extract
// the table name. Stream ARN format:
//
//	arn:aws:dynamodb:region:accountId:table/tableName/stream/label
func extractTableNameFromStreamArn(streamArn string) string {
	idx := indexOf(streamArn, "table/")
	if idx < 0 {
		return ""
	}
	rest := streamArn[idx+6:]
	slashIdx := indexOf(rest, "/")
	if slashIdx < 0 {
		return rest
	}
	return rest[:slashIdx]
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// encodeShardIterator creates an opaque iterator string from the table
// name and sequence number. Format: "tableName|seqNum".
func encodeShardIterator(tableName string, seq int64) string {
	return fmt.Sprintf("%s|%d", tableName, seq)
}

// decodeShardIterator parses an iterator string back into table name and
// sequence number.
func decodeShardIterator(iterator string) (string, int64, error) {
	for i := len(iterator) - 1; i >= 0; i-- {
		if iterator[i] == '|' {
			tableName := iterator[:i]
			seqStr := iterator[i+1:]
			seq, err := strconv.ParseInt(seqStr, 10, 64)
			if err != nil {
				return "", 0, err
			}
			return tableName, seq, nil
		}
	}
	return "", 0, fmt.Errorf("invalid iterator format")
}

// streamTimeNow returns the current time. Extracted for potential testing.
var streamTimeNow = func() time.Time { return time.Now().UTC() }

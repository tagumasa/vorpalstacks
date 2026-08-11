package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

	latestSeq, seqErr := store.Streams().GetLatestSequence(tableName)
	if seqErr != nil {
		return nil, seqErr
	}

	seqRange := map[string]interface{}{
		"StartingSequenceNumber": fmt.Sprintf("%d", int64(1)),
	}
	if latestSeq >= 1 {
		seqRange["EndingSequenceNumber"] = fmt.Sprintf("%d", latestSeq)
	}

	// StreamStatus reflects the table's streaming state.
	streamStatus := "ENABLED"
	if table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		streamStatus = "DISABLED"
	}

	shard := map[string]interface{}{
		"ShardId":             dbstore.ShardIDForStream(streamArn),
		"SequenceNumberRange": seqRange,
	}

	return map[string]interface{}{
		"StreamDescription": map[string]interface{}{
			"StreamArn":               streamArn,
			"StreamLabel":             table.LatestStreamLabel,
			"StreamStatus":            streamStatus,
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
		default:
			return nil, ErrInvalidParameter
		}
	}
	if limit <= 0 || limit > 1000 {
		return nil, ErrInvalidParameter
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

// ListStreams returns stream ARNs associated with the current account and
// endpoint. If TableName is provided, only streams for that table are returned.
// AWS API: DynamoDB Streams — ListStreams
func (s *DynamoDBService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableNameFilter := request.GetStringParam(req.Parameters, "TableName")
	exclusiveStartStreamArn := request.GetStringParam(req.Parameters, "ExclusiveStartStreamArn")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit == 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}

	store, err := s.GetCachedStoreForRegion(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	// Collect ALL tables by walking every store page. A single List call is
	// capped at DefaultMaxItems (100), so tables beyond that would be
	// invisible if we only fetched one page.
	type streamEntry struct {
		StreamArn   string
		TableName   string
		StreamLabel string
	}

	var allStreams []streamEntry
	{
		var marker string
		for {
			page, nextMarker, err := store.Tables().List(marker, 0)
			if err != nil {
				return nil, err
			}
			for _, t := range page {
				if tableNameFilter != "" && t.Name != tableNameFilter {
					continue
				}
				if t.StreamSpecification == nil || !t.StreamSpecification.StreamEnabled {
					continue
				}
				allStreams = append(allStreams, streamEntry{
					StreamArn:   t.StreamArn,
					TableName:   t.Name,
					StreamLabel: t.LatestStreamLabel,
				})
			}
			if nextMarker == "" {
				break
			}
			marker = nextMarker
		}
	}

	// Apply cursor-based pagination over the complete stream list.
	var streams []streamEntry
	if exclusiveStartStreamArn == "" {
		streams = allStreams
	} else {
		// Find the cursor position. AWS returns an empty page when the
		// cursor is invalid (stale or nonexistent), NOT a restart.
		cursorIdx := -1
		for i, st := range allStreams {
			if st.StreamArn == exclusiveStartStreamArn {
				cursorIdx = i
				break
			}
		}
		if cursorIdx == -1 {
			streams = nil // invalid cursor → empty page
		} else {
			streams = allStreams[cursorIdx+1:]
		}
	}

	hasMore := limit > 0 && len(streams) > limit
	if hasMore {
		streams = streams[:limit]
	}

	streamList := make([]map[string]interface{}, 0, len(streams))
	for _, st := range streams {
		streamList = append(streamList, map[string]interface{}{
			"StreamArn":   st.StreamArn,
			"TableName":   st.TableName,
			"StreamLabel": st.StreamLabel,
		})
	}

	result := map[string]interface{}{
		"Streams": streamList,
	}

	if hasMore && len(streamList) > 0 {
		result["LastEvaluatedStreamArn"] = streams[len(streams)-1].StreamArn
	}

	return result, nil
}

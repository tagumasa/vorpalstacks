package dynamodb

import (
	"fmt"
	"strconv"

	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Stream Core — single validation + persistence path for DynamoDB Streams
// operations.
//
// These methods encapsulate stream lifecycle logic. Both the HTTP API
// handlers (stream_operations.go) and any future admin handler delegate to
// these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// DescribeStreamResult is the service-layer result of DescribeStream.
type DescribeStreamResult struct {
	StreamArn               string
	StreamLabel             string
	StreamStatus            string
	StreamViewType          string
	TableName               string
	KeySchema               []*dbstore.KeySchemaElement
	Shards                  []ShardInfo
	CreationRequestDateTime int64
}

// ShardInfo describes a single shard within a stream.
type ShardInfo struct {
	ShardID                string
	StartingSequenceNumber string
	EndingSequenceNumber   string
}

// describeStreamCore builds the description of a single DynamoDB stream.
// Returns ErrResourceNotFound when the table or stream does not exist.
func (s *DynamoDBService) describeStreamCore(store dbstore.DynamoDBStoreInterface, streamArn string) (*DescribeStreamResult, error) {
	tableName := extractTableNameFromStreamArn(streamArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil || table == nil || table.StreamArn != streamArn {
		return nil, ErrResourceNotFound
	}

	latestSeq, seqErr := store.Streams().GetLatestSequence(tableName)
	if seqErr != nil {
		return nil, seqErr
	}

	streamStatus := "ENABLED"
	streamViewType := ""
	if table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		streamStatus = "DISABLED"
	} else {
		streamViewType = string(table.StreamSpecification.StreamViewType)
	}

	result := &DescribeStreamResult{
		StreamArn:               streamArn,
		StreamLabel:             table.LatestStreamLabel,
		StreamStatus:            streamStatus,
		StreamViewType:          streamViewType,
		TableName:               tableName,
		KeySchema:               table.KeySchema,
		CreationRequestDateTime: table.CreationDateTime.Unix(),
	}
	shard := ShardInfo{
		ShardID:                dbstore.ShardIDForStream(streamArn),
		StartingSequenceNumber: "1",
	}
	if latestSeq >= 1 {
		shard.EndingSequenceNumber = fmt.Sprintf("%d", latestSeq)
	}
	result.Shards = []ShardInfo{shard}
	return result, nil
}

// GetShardIteratorResult is the service-layer result of GetShardIterator.
type GetShardIteratorResult struct {
	ShardIterator string
}

// getShardIteratorCore computes and encodes a shard iterator for the
// requested position. Returns ErrResourceNotFound when the table or stream
// does not exist, ErrInvalidParameter for unknown iterator types or bad
// sequence numbers.
func (s *DynamoDBService) getShardIteratorCore(store dbstore.DynamoDBStoreInterface, streamArn, iteratorType, sequenceNumber string) (*GetShardIteratorResult, error) {
	tableName := extractTableNameFromStreamArn(streamArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

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
		seq, parseErr := parseSequenceNumber(sequenceNumber)
		if parseErr != nil {
			return nil, ErrInvalidParameter
		}
		startSeq = seq - 1
	case "AFTER_SEQUENCE_NUMBER":
		seq, parseErr := parseSequenceNumber(sequenceNumber)
		if parseErr != nil {
			return nil, ErrInvalidParameter
		}
		startSeq = seq
	default:
		return nil, ErrInvalidParameter
	}

	return &GetShardIteratorResult{
		ShardIterator: encodeShardIterator(tableName, startSeq),
	}, nil
}

// GetRecordsResult is the service-layer result of GetRecords.
type GetRecordsResult struct {
	Records           []interface{}
	NextShardIterator string
}

// getRecordsCore retrieves up to limit stream records starting from the
// decoded iterator position. Returns ErrResourceNotFound when the table
// does not exist or streaming is disabled.
func (s *DynamoDBService) getRecordsCore(store dbstore.DynamoDBStoreInterface, tableName string, fromSeq int64, limit int) (*GetRecordsResult, error) {
	table, err := store.Tables().Get(tableName)
	if err != nil || table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return nil, ErrResourceNotFound
	}

	// Reads starting at or below the trimmed floor would miss records the
	// retention window has already removed.
	floor, err := store.Streams().OldestSequence(tableName)
	if err != nil {
		logs.Warn("failed to read stream trim floor",
			logs.String("table", tableName), logs.Err(err))
		return nil, ErrInternal
	}
	if fromSeq < floor {
		return nil, ErrTrimmedDataAccess
	}

	records, nextSeq, err := store.Streams().GetRecords(tableName, fromSeq, limit)
	if err != nil {
		logs.Warn("failed to get stream records",
			logs.String("table", tableName),
			logs.Err(err))
		return nil, ErrInternal
	}

	recordsResp := make([]interface{}, 0, len(records))
	for _, rec := range records {
		recordsResp = append(recordsResp, rec)
	}

	return &GetRecordsResult{
		Records:           recordsResp,
		NextShardIterator: encodeShardIterator(tableName, nextSeq),
	}, nil
}

// StreamEntry is a single entry in the ListStreams result.
type StreamEntry struct {
	StreamArn   string
	TableName   string
	StreamLabel string
}

// ListStreamsResult is the service-layer result of ListStreams.
type ListStreamsResult struct {
	Streams                []StreamEntry
	LastEvaluatedStreamArn string
}

// listStreamsCore walks every table in the store, collects those with
// streaming enabled, applies the optional TableName filter, and returns a
// single page of stream entries.
func (s *DynamoDBService) listStreamsCore(store dbstore.DynamoDBStoreInterface, tableNameFilter, exclusiveStartStreamArn string, limit int) (*ListStreamsResult, error) {
	// Collect ALL tables by walking every store page. A single List call
	// is capped at DefaultMaxItems (100), so tables beyond that would be
	// invisible if we only fetched one page.
	var allStreams []StreamEntry
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
				allStreams = append(allStreams, StreamEntry{
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
	var streams []StreamEntry
	if exclusiveStartStreamArn == "" {
		streams = allStreams
	} else {
		cursorIdx := -1
		for i, st := range allStreams {
			if st.StreamArn == exclusiveStartStreamArn {
				cursorIdx = i
				break
			}
		}
		if cursorIdx == -1 {
			streams = nil
		} else {
			streams = allStreams[cursorIdx+1:]
		}
	}

	hasMore := limit > 0 && len(streams) > limit
	if hasMore {
		streams = streams[:limit]
	}

	result := &ListStreamsResult{
		Streams: streams,
	}
	if hasMore && len(streams) > 0 {
		result.LastEvaluatedStreamArn = streams[len(streams)-1].StreamArn
	}
	return result, nil
}

// parseSequenceNumber parses a sequence-number string used by
// AT_SEQUENCE_NUMBER / AFTER_SEQUENCE_NUMBER iterator types.
func parseSequenceNumber(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

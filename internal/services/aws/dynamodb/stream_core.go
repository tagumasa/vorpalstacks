package dynamodb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/internal/utils/aws/arn"
	crypto "vorpalstacks/internal/utils/crypto"
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
// Returns ErrInvalidParameter when StreamArn is missing and
// ErrResourceNotFound when the table or stream does not exist.
func (s *DynamoDBService) describeStreamCore(store dbstore.DynamoDBStoreInterface, streamArn string) (*DescribeStreamResult, error) {
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}
	tableName := arn.ParseStreamARN(streamArn)
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

// encodeShardIterator creates an opaque iterator string from the table
// name and sequence number. Format: "tableName|seqNum".
// shardIteratorTTL is the documented shard iterator lifetime: a shard
// iterator expires fifteen minutes after it was issued.
const shardIteratorTTL = 15 * time.Minute

// encodeShardIterator renders an opaque, signed iterator for a read
// position. The payload carries the table, the sequence number to read
// from, and the issue time; the HMAC-SHA256 signature under the
// store-persisted signing key makes the token unforgeable — a client can
// neither read a crafted position into the stream nor tamper with an
// issued one, because only the server holds the key.
func encodeShardIterator(signingKey []byte, tableName string, seq int64) string {
	payload := fmt.Sprintf("%s|%d|%d", tableName, seq, streamTimeNow().Unix())
	mac := crypto.HMACSHA256(signingKey, []byte(payload))
	return base64.RawURLEncoding.EncodeToString(append([]byte(payload), mac...))
}

// decodeShardIterator verifies an iterator's signature and returns the
// table name, sequence number, and issue time it carries. A token whose
// signature does not verify — anything a client constructed or altered —
// is rejected, as is any malformed encoding.
func decodeShardIterator(signingKey []byte, iterator string) (string, int64, int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(iterator)
	if err != nil || len(raw) <= sha256.Size {
		return "", 0, 0, fmt.Errorf("invalid iterator format")
	}
	split := len(raw) - sha256.Size
	payload, mac := raw[:split], raw[split:]
	if !hmac.Equal(crypto.HMACSHA256(signingKey, payload), mac) {
		return "", 0, 0, fmt.Errorf("invalid iterator signature")
	}
	parts := strings.Split(string(payload), "|")
	if len(parts) != 3 {
		return "", 0, 0, fmt.Errorf("invalid iterator format")
	}
	tableName := parts[0]
	seq, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, 0, err
	}
	issuedAt, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", 0, 0, err
	}
	return tableName, seq, issuedAt, nil
}

// shardIteratorExpired reports whether an iterator issued at the given
// unix time has passed the documented fifteen-minute lifetime.
func shardIteratorExpired(issuedAtUnix int64, now time.Time) bool {
	return now.Unix()-issuedAtUnix >= int64(shardIteratorTTL/time.Second)
}

// streamTimeNow returns the current time. Extracted for potential testing.
var streamTimeNow = func() time.Time { return time.Now().UTC() }

// getShardIteratorCore computes and encodes a shard iterator for the
// requested position. Returns ErrInvalidParameter when a required parameter
// is missing or the iterator type or sequence number is unknown, and
// ErrResourceNotFound when the table or stream does not exist.
func (s *DynamoDBService) getShardIteratorCore(store dbstore.DynamoDBStoreInterface, streamArn, shardId, iteratorType, sequenceNumber string) (*GetShardIteratorResult, error) {
	if streamArn == "" || shardId == "" || iteratorType == "" {
		return nil, ErrInvalidParameter
	}
	tableName := arn.ParseStreamARN(streamArn)
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

	signingKey, err := store.Streams().IteratorSigningKey()
	if err != nil {
		return nil, err
	}

	return &GetShardIteratorResult{
		ShardIterator: encodeShardIterator(signingKey, tableName, startSeq),
	}, nil
}

// GetRecordsResult is the service-layer result of GetRecords.
type GetRecordsResult struct {
	Records           []interface{}
	NextShardIterator string
}

// getRecordsCore retrieves up to limit stream records starting from the
// decoded iterator position. The iterator's format and fifteen-minute
// lifetime and the Limit range are validated here — the model binds Limit
// to PositiveLongObject (minimum 1) and documents values above 1000 as a
// LimitExceededException. Returns ErrResourceNotFound when the table does
// not exist or streaming is disabled.
func (s *DynamoDBService) getRecordsCore(store dbstore.DynamoDBStoreInterface, iterator string, limit int, limitSet bool) (*GetRecordsResult, error) {
	// An explicit Limit below 1 is invalid; a plain value check cannot tell
	// "unset" apart from an explicit zero, so presence is signalled
	// separately by the caller.
	if limitSet && limit < 1 {
		return nil, ErrInvalidParameter
	}
	if limit > getRecordsMaxLimit {
		return nil, ErrStreamsLimitExceeded
	}
	if limit == 0 {
		limit = getRecordsDefaultLimit
	}
	if iterator == "" {
		return nil, ErrInvalidParameter
	}
	signingKey, err := store.Streams().IteratorSigningKey()
	if err != nil {
		return nil, err
	}
	tableName, fromSeq, issuedAt, err := decodeShardIterator(signingKey, iterator)
	if err != nil {
		return nil, ErrInvalidParameter
	}
	if shardIteratorExpired(issuedAt, streamTimeNow()) {
		return nil, ErrExpiredIterator
	}

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
		NextShardIterator: encodeShardIterator(signingKey, tableName, nextSeq),
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
// single page of stream entries. The Limit default (and its cap) is applied
// here so every caller shares one range policy.
func (s *DynamoDBService) listStreamsCore(store dbstore.DynamoDBStoreInterface, tableNameFilter, exclusiveStartStreamArn string, limit int) (*ListStreamsResult, error) {
	if limit == 0 {
		limit = listStreamsDefaultLimit
	}
	if limit > listStreamsMaxLimit {
		limit = listStreamsMaxLimit
	}

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

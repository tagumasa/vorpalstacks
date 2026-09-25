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
	// LastEvaluatedShardId names the last shard of the returned page when
	// more shards exist; empty means the whole shard list was served.
	LastEvaluatedShardId string
}

// ShardInfo describes a single shard within a stream.
type ShardInfo struct {
	ShardID                string
	StartingSequenceNumber string
	EndingSequenceNumber   string
}

// DescribeStreamInput is the service-layer input of DescribeStream. LimitSet
// distinguishes an omitted Limit from an explicit zero because the model's
// PositiveIntegerObject range (minimum 1) rejects the latter while the former
// takes the default.
type DescribeStreamInput struct {
	StreamArn             string
	Limit                 int
	LimitSet              bool
	ExclusiveStartShardId string
	ShardFilter           *ShardFilterSpec
}

// ShardFilterSpec is the DescribeStream ShardFilter. The model's
// ShardFilterType enum defines CHILD_SHARDS alone, and the member's
// documentation names the shard whose children are requested.
type ShardFilterSpec struct {
	Type    string
	ShardId string
}

// describeStreamCore builds the description of a single DynamoDB stream.
// Returns ErrInvalidParameter when StreamArn is missing, a member carries an
// out-of-range value, or the ShardFilter names a type outside the model's
// enum; ErrResourceNotFound when the table or stream does not exist.
func (s *DynamoDBService) describeStreamCore(store dbstore.DynamoDBStoreInterface, in DescribeStreamInput) (*DescribeStreamResult, error) {
	if in.StreamArn == "" {
		return nil, ErrInvalidParameter
	}
	if in.LimitSet && in.Limit < 1 {
		return nil, ErrInvalidParameter
	}
	if in.ExclusiveStartShardId != "" && (len(in.ExclusiveStartShardId) < 28 || len(in.ExclusiveStartShardId) > 65) {
		return nil, ErrInvalidParameter
	}
	if in.ShardFilter != nil && in.ShardFilter.Type != "CHILD_SHARDS" {
		return nil, ErrInvalidParameter
	}
	limit := in.Limit
	if limit == 0 {
		limit = describeStreamDefaultLimit
	}
	if limit > describeStreamMaxLimit {
		limit = describeStreamMaxLimit
	}

	tableName := arn.ParseStreamARN(in.StreamArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		// Absence maps to the stream's not-found answer; a storage fault
		// (bucket I/O, an unreadable record) is reported as the storage
		// error it is — never masquerading as the stream's absence.
		if storeRecordMissing(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	if table == nil || table.StreamArn != in.StreamArn {
		return nil, ErrResourceNotFound
	}

	latestSeq, seqErr := store.Streams().GetLatestSequenceForStream(tableName, in.StreamArn)
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
		StreamArn:               in.StreamArn,
		StreamLabel:             table.LatestStreamLabel,
		StreamStatus:            streamStatus,
		StreamViewType:          streamViewType,
		TableName:               tableName,
		KeySchema:               table.KeySchema,
		CreationRequestDateTime: table.CreationDateTime.Unix(),
	}
	shard := ShardInfo{
		ShardID:                dbstore.ShardIDForStream(in.StreamArn),
		StartingSequenceNumber: dbstore.FormatStreamSequenceNumber(1),
	}
	if latestSeq >= 1 {
		shard.EndingSequenceNumber = dbstore.FormatStreamSequenceNumber(latestSeq)
	}
	shards := []ShardInfo{shard}

	// The CHILD_SHARDS filter answers with the children of the named
	// shard; the platform's streams expose a single shard that never
	// splits, so no shard has children and the filtered list is empty.
	if in.ShardFilter != nil {
		shards = nil
	}

	// The exclusive-start cursor positions the page after the matching
	// shard; a cursor naming no shard of this stream leaves nothing to
	// serve, the same policy listStreamsCore applies to its cursor.
	if in.ExclusiveStartShardId != "" {
		cursorIdx := -1
		for i, sh := range shards {
			if sh.ShardID == in.ExclusiveStartShardId {
				cursorIdx = i
				break
			}
		}
		if cursorIdx == -1 {
			shards = nil
		} else {
			shards = shards[cursorIdx+1:]
		}
	}

	if len(shards) > limit {
		shards = shards[:limit]
		result.LastEvaluatedShardId = shards[len(shards)-1].ShardID
	}
	result.Shards = shards
	return result, nil
}

// GetShardIteratorResult is the service-layer result of GetShardIterator.
type GetShardIteratorResult struct {
	ShardIterator string
}

// encodeShardIterator creates an opaque iterator string from the stream
// ARN, table name, and sequence number. Format:
// "streamArn|tableName|seqNum|issuedAt|iteratorType".
// shardIteratorTTL is the documented shard iterator lifetime: a shard
// iterator expires fifteen minutes after it was issued.
const shardIteratorTTL = 15 * time.Minute

// encodeShardIterator renders an opaque, signed iterator for a read
// position. The payload carries the issuing stream's ARN, the table, the
// sequence number to read from, the issue time, and the iterator type the
// position was issued under; the HMAC-SHA256 signature under the
// store-persisted signing key makes the token unforgeable — a client can
// neither read a crafted position into the stream nor tamper with an
// issued one, because only the server holds the key. The ARN binds the
// iterator to its stream generation: a table whose stream is disabled and
// re-enabled carries a fresh ARN, and an iterator of the superseded
// generation must never read the successor stream. The type travels with
// the position so the read side can apply the type's own semantics — a
// TRIM_HORIZON position follows the retention floor at read time, not the
// floor's historical value at issue time.
func encodeShardIterator(signingKey []byte, streamArn, tableName string, seq int64, iteratorType string) string {
	payload := fmt.Sprintf("%s|%s|%d|%d|%s", streamArn, tableName, seq, streamTimeNow().Unix(), iteratorType)
	mac := crypto.HMACSHA256(signingKey, []byte(payload))
	return base64.RawURLEncoding.EncodeToString(append([]byte(payload), mac...))
}

// decodeShardIterator verifies an iterator's signature and returns the
// issuing stream ARN, table name, sequence number, issue time, and
// iterator type it carries. A token whose signature does not verify —
// anything a client constructed or altered — is rejected, as is any
// malformed encoding.
func decodeShardIterator(signingKey []byte, iterator string) (string, string, int64, int64, string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(iterator)
	if err != nil || len(raw) <= sha256.Size {
		return "", "", 0, 0, "", fmt.Errorf("invalid iterator format")
	}
	split := len(raw) - sha256.Size
	payload, mac := raw[:split], raw[split:]
	if !hmac.Equal(crypto.HMACSHA256(signingKey, payload), mac) {
		return "", "", 0, 0, "", fmt.Errorf("invalid iterator signature")
	}
	parts := strings.Split(string(payload), "|")
	if len(parts) != 5 {
		return "", "", 0, 0, "", fmt.Errorf("invalid iterator format")
	}
	streamArn := parts[0]
	tableName := parts[1]
	seq, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", "", 0, 0, "", err
	}
	issuedAt, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return "", "", 0, 0, "", err
	}
	return streamArn, tableName, seq, issuedAt, parts[4], nil
}

// shardIteratorExpired reports whether an iterator issued at the given
// unix time has passed the documented fifteen-minute lifetime.
func shardIteratorExpired(issuedAtUnix int64, now time.Time) bool {
	return now.Unix()-issuedAtUnix >= int64(shardIteratorTTL/time.Second)
}

// streamTimeNow is the stream plane's clock seam: the core-validation and
// sweeper tests shift it past the shard-iterator expiry and retention
// horizons to pin those behaviours without waiting real time; every
// production read treats it as the UTC now.
var streamTimeNow = func() time.Time { return time.Now().UTC() }

// getShardIteratorCore computes and encodes a shard iterator for the
// requested position. Returns ErrInvalidParameter when a required parameter
// is missing or the iterator type, shard id, or sequence number is invalid,
// ErrResourceNotFound when the table, stream, or shard does not exist, and
// ErrTrimmedDataAccess when the requested sequence number lies at or below
// the retention trim floor.
func (s *DynamoDBService) getShardIteratorCore(store dbstore.DynamoDBStoreInterface, streamArn, shardId, iteratorType, sequenceNumber string) (*GetShardIteratorResult, error) {
	if streamArn == "" || shardId == "" || iteratorType == "" {
		return nil, ErrInvalidParameter
	}
	// The model bounds ShardId at 28-65 characters.
	if len(shardId) < 28 || len(shardId) > 65 {
		return nil, ErrInvalidParameter
	}
	tableName := arn.ParseStreamARN(streamArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		// The same fault-versus-absence split the describe path applies.
		if storeRecordMissing(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	if table == nil || table.StreamArn != streamArn {
		return nil, ErrResourceNotFound
	}

	// The stream exposes exactly one shard, so any shard id other than the
	// stream's own names a shard that does not exist.
	if shardId != dbstore.ShardIDForStream(streamArn) {
		return nil, ErrResourceNotFound
	}

	latestSeq, err := store.Streams().GetLatestSequenceForStream(tableName, streamArn)
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

	// A sequence-number position whose first record the retention window
	// has already removed is the documented TrimmedDataAccessException at
	// iterator-issuing time.
	if iteratorType == "AT_SEQUENCE_NUMBER" || iteratorType == "AFTER_SEQUENCE_NUMBER" {
		floor, floorErr := store.Streams().OldestSequence(tableName)
		if floorErr != nil {
			return nil, floorErr
		}
		firstServed := startSeq + 1
		if firstServed <= floor {
			return nil, ErrTrimmedDataAccess
		}
	}

	signingKey, err := store.Streams().IteratorSigningKey()
	if err != nil {
		return nil, err
	}

	return &GetShardIteratorResult{
		ShardIterator: encodeShardIterator(signingKey, streamArn, tableName, startSeq, iteratorType),
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
	iterStreamArn, tableName, fromSeq, issuedAt, iteratorType, err := decodeShardIterator(signingKey, iterator)
	if err != nil {
		return nil, ErrInvalidParameter
	}
	if shardIteratorExpired(issuedAt, streamTimeNow()) {
		return nil, ErrExpiredIterator
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		// The same fault-versus-absence split the describe path applies.
		if storeRecordMissing(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	if table == nil || table.StreamSpecification == nil || !table.StreamSpecification.StreamEnabled {
		return nil, ErrResourceNotFound
	}
	// The iterator belongs to the stream generation that issued it. A
	// table whose stream is disabled and re-enabled carries a fresh ARN,
	// so only this comparison keeps a superseded generation's iterator
	// from reading the successor stream.
	if table.StreamArn != iterStreamArn {
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
	// A TRIM_HORIZON position reads the oldest retained record wherever
	// the retention floor has moved to: the read clamps the position up to
	// the floor, so a trim that happened after the iterator was issued —
	// or between the pages of a long read — still serves the oldest
	// survivor. A LATEST position is the same class: the iterator froze a
	// position the window has since moved past, and the read serves the
	// oldest survivor (or empty) rather than erroring. TrimmedDataAccess
	// is the answer to a position the caller named (AT/AFTER a sequence
	// number the window has removed), not to the horizon itself.
	if (iteratorType == "TRIM_HORIZON" || iteratorType == "LATEST") && fromSeq < floor {
		fromSeq = floor
	}
	if fromSeq < floor {
		return nil, ErrTrimmedDataAccess
	}

	records, nextSeq, err := store.Streams().GetRecords(tableName, iterStreamArn, fromSeq, limit)
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
		NextShardIterator: encodeShardIterator(signingKey, iterStreamArn, tableName, nextSeq, iteratorType),
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
// single page of stream entries. The Limit range (explicit values must be at
// least 1), its default, and its cap are applied here so every caller shares
// one policy.
func (s *DynamoDBService) listStreamsCore(store dbstore.DynamoDBStoreInterface, tableNameFilter, exclusiveStartStreamArn string, limit int, limitSet bool) (*ListStreamsResult, error) {
	if limitSet && limit < 1 {
		return nil, ErrInvalidParameter
	}
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

	hasMore := len(streams) > limit
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
// AT_SEQUENCE_NUMBER / AFTER_SEQUENCE_NUMBER iterator types. The model
// bounds the SequenceNumber type at 21-40 characters; the platform's own
// sequence numbers are zero-padded to that minimum on the wire, so a shorter
// or longer string never names one of the stream's records.
func parseSequenceNumber(s string) (int64, error) {
	if len(s) < 21 || len(s) > 40 {
		return 0, fmt.Errorf("sequence number length outside the model's 21-40 range")
	}
	return strconv.ParseInt(s, 10, 64)
}

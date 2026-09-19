package kinesis

import (
	"encoding/base64"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// PutRecordInput is the transport-agnostic input for PutRecord. The record
// members are validated inside the Core in the original precedence order
// (partition key before the stream fetch, data size after it). HasData
// carries the wire presence of the required Data member — a string alone
// cannot distinguish an absent member from an explicitly empty payload,
// which stays a size-legal blob.
type PutRecordInput struct {
	StreamName                string
	StreamARN                 string
	Data                      string
	HasData                   bool
	PartitionKey              string
	ExplicitHashKey           string
	SequenceNumberForOrdering string
	DryRun                    bool
}

// PutRecordResult carries the write receipt of a single record.
type PutRecordResult struct {
	ShardID        string
	SequenceNumber string
	EncryptionType string
}

// PutRecordsInput is the transport-agnostic input for PutRecords. The raw
// wire list travels untyped (nil-able) because the per-entry validation runs
// in the Core after the stream fetch, in the original precedence order.
type PutRecordsInput struct {
	StreamName string
	StreamARN  string
	Records    interface{}
	DryRun     bool
}

// PutRecordsResult carries the store-level per-entry outcomes.
type PutRecordsResult struct {
	Results        []kinesisstore.PutRecordResult
	EncryptionType string
}

// GetRecordsInput is the transport-agnostic input for GetRecords.
type GetRecordsInput struct {
	ShardIterator string
	Limit         int32
	DryRun        bool
}

// GetRecordsResult carries the records page plus the continuation state.
type GetRecordsResult struct {
	Records            []*kinesisstore.Record
	NextShardIterator  interface{}
	MillisBehindLatest int64
	EncryptionType     string
	ChildShards        []*kinesisstore.Shard
}

// GetShardIteratorInput is the transport-agnostic input for GetShardIterator.
type GetShardIteratorInput struct {
	StreamName             string
	StreamARN              string
	ShardId                string
	ShardIteratorType      string
	StartingSequenceNumber string
	Timestamp              *time.Time
	DryRun                 bool
}

// putRecordCore validates and writes a single record into the stream.
func (s *KinesisService) putRecordCore(reqCtx *request.RequestContext, input PutRecordInput) (PutRecordResult, error) {
	// Data is a required member of the input shape: an absent one is a
	// shape violation rejecting the request — the same reading the batch
	// twin applies per entry — never an empty record to ingest. The wire
	// type is the strict string read's contract at the parse seam: a
	// non-string value never reaches the Core as a coerced empty payload.
	if !input.HasData {
		return PutRecordResult{}, ErrInvalidArgument
	}
	if !validatePartitionKey(input.PartitionKey) {
		return PutRecordResult{}, ErrInvalidArgument
	}
	if !validateExplicitHashKey(input.ExplicitHashKey) {
		return PutRecordResult{}, ErrInvalidArgument
	}
	// SequenceNumberForOrdering carries the SequenceNumber shape's pattern;
	// the value needs no further effect here — the platform assigns
	// sequence numbers from the arrival clock and a per-store counter, so
	// consecutive puts to the same shard (a partition key always routes to
	// the same shard) already receive strictly increasing numbers while
	// the clock is monotonic, which is the member's documented guarantee
	// delivered by assignment rather than by honouring the supplied value.
	if !validateSequenceNumber(input.SequenceNumberForOrdering) {
		return PutRecordResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return PutRecordResult{}, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return PutRecordResult{}, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return PutRecordResult{}, s.mapStoreError(err)
	}

	if !validateRecordSize(input.Data, input.PartitionKey, stream.MaxRecordSizeInKiB) {
		return PutRecordResult{}, ErrInvalidArgument
	}

	// A DryRun request that passes every check answers the model's
	// dedicated error — the request "was rejected because the DryRun
	// parameter was specified" — without performing the write.
	if input.DryRun {
		return PutRecordResult{}, ErrDryRunOperation
	}

	record, targetShardID, err := store.PutRecordWithShardSelection(streamName, input.PartitionKey, input.Data, input.ExplicitHashKey)
	if err != nil {
		return PutRecordResult{}, s.mapStoreError(err)
	}

	return PutRecordResult{
		ShardID:        targetShardID,
		SequenceNumber: record.SequenceNumber,
		EncryptionType: resolveEncryptionType(stream),
	}, nil
}

// putRecordsCore validates and writes a batch of records into the stream.
// Batch-shape problems — a missing, mistyped, empty or oversized Records
// member — reject the request. Entry validation is request-level too: the
// model's per-entry failure channel is bounded to the two write-time
// outcomes the ErrorCode member documents (ProvisionedThroughputExceededException
// or InternalFailure), so a member that violates a single entry's shape — a
// missing Data, a malformed partition or hash key, an oversized record —
// has no per-entry identity to report and rejects the whole request.
func (s *KinesisService) putRecordsCore(reqCtx *request.RequestContext, input PutRecordsInput) (PutRecordsResult, error) {
	if input.Records == nil {
		return PutRecordsResult{}, ErrInvalidArgument
	}
	recordsList, ok := input.Records.([]interface{})
	if !ok {
		return PutRecordsResult{}, ErrInvalidArgument
	}
	// The PutRecordsRequestEntryList length trait bounds the batch to 1-500
	// entries; an empty or oversized member is a batch-shape violation
	// answering the identity the operation declares for invalid input.
	if len(recordsList) < 1 || len(recordsList) > kinesisstore.MaxPutRecordsEntries {
		return PutRecordsResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return PutRecordsResult{}, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return PutRecordsResult{}, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return PutRecordsResult{}, s.mapStoreError(err)
	}

	// Entry validation covers the whole batch before any write happens:
	// every entry must satisfy its shape, and a violation answers the
	// request's identity. Data is a required member of the entry shape —
	// an entry without it (or with a non-string value) is a shape
	// violation, not an empty record to ingest. The same pass sums the
	// request's size the way the operation's documented whole-request cap
	// counts it: every entry's decoded payload plus its partition key.
	// Every accepted entry appends exactly one request, so the store's
	// result list aligns with the entry list by position.
	entries := make([]kinesisstore.PutRecordResult, len(recordsList))
	requests := make([]kinesisstore.PutRecordRequest, 0, len(recordsList))
	var requestBytes int64
	for _, r := range recordsList {
		rm, ok := r.(map[string]interface{})
		if !ok {
			return PutRecordsResult{}, ErrInvalidArgument
		}
		rawData, hasData := rm["Data"]
		if !hasData {
			return PutRecordsResult{}, ErrInvalidArgument
		}
		data, ok := rawData.(string)
		if !ok {
			return PutRecordsResult{}, ErrInvalidArgument
		}
		pk, _ := rm["PartitionKey"].(string)
		ehk, _ := rm["ExplicitHashKey"].(string)
		if !validatePartitionKey(pk) {
			return PutRecordsResult{}, ErrInvalidArgument
		}
		if !validateExplicitHashKey(ehk) {
			return PutRecordsResult{}, ErrInvalidArgument
		}
		if !validateRecordSize(data, pk, stream.MaxRecordSizeInKiB) {
			return PutRecordsResult{}, ErrInvalidArgument
		}
		// The per-record rule already required the payload to decode.
		payload, _ := base64.StdEncoding.DecodeString(data)
		requestBytes += int64(len(payload)) + int64(len(pk))
		requests = append(requests, kinesisstore.PutRecordRequest{
			Data:            data,
			PartitionKey:    pk,
			ExplicitHashKey: ehk,
		})
	}
	if requestBytes > kinesisstore.MaxPutRecordsRequestBytes {
		return PutRecordsResult{}, ErrInvalidArgument
	}

	// A DryRun request that passes every check answers the model's
	// dedicated error without performing the writes — the response's
	// per-entry receipts describe writes that must not happen.
	if input.DryRun {
		return PutRecordsResult{}, ErrDryRunOperation
	}

	results, err := store.PutRecords(streamName, requests)
	if err != nil {
		return PutRecordsResult{}, s.mapStoreError(err)
	}
	copy(entries, results)

	return PutRecordsResult{
		Results:        entries,
		EncryptionType: resolveEncryptionType(stream),
	}, nil
}

// shardPage is one read page from a shard: the records at or after the
// given sequence number, the follow-up position, the distance from the
// stream tip, and the shard's closure state with its children.
type shardPage struct {
	Records            []*kinesisstore.Record
	LastSeqNum         string
	MillisBehindLatest int64
	ShardClosed        bool
	ChildShards        []*kinesisstore.Shard
}

// readShardPage is the single record-read primitive shared by the
// GetRecords path and the SubscribeToShard pump: it reads the page after
// the given sequence number under the stream's retention window, computes
// the advance position (aged-out records move it forward, so neither
// transport can strand on a fully expired page) and the distance from the
// tip, and collects the shard's closure state and children. The stream
// lookup here is load-bearing, not a caller-duplicated waste: it is how a
// stream deleted mid-subscription surfaces as an error on the pump, whose
// own stream pointer went stale at delete time (the GetRecords path keeps
// its own earlier lookup for the encryption type and the same liveness
// question). The shard record's own read stays advisory — a shard record
// that cannot be read reads as an open shard — but a closed shard whose
// children cannot be read is a failure: reporting closure without its
// continuation would tell the client the stream ended there.
func (s *KinesisService) readShardPage(store *kinesisstore.KinesisStore, streamName, shardID, afterSeqNum string, limit int32, includeStart bool) (shardPage, error) {
	stream, err := store.GetStream(streamName)
	if err != nil {
		return shardPage{}, err
	}
	retentionCutoff := kinesisstore.RetentionCutoff(stream.RetentionPeriodHours)

	records, advanceSeqNum, err := store.GetRecords(streamName, shardID, afterSeqNum, limit, includeStart, retentionCutoff)
	if err != nil {
		return shardPage{}, err
	}

	var millisBehindLatest int64
	if len(records) > 0 {
		millis := time.Since(records[len(records)-1].ApproximateArrivalTimestamp).Milliseconds()
		if millis > 0 {
			millisBehindLatest = millis
		}
	}

	shard, _ := store.GetShard(streamName, shardID)
	page := shardPage{
		Records:            records,
		LastSeqNum:         advanceSeqNum,
		MillisBehindLatest: millisBehindLatest,
		ShardClosed:        shard != nil && shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber != "",
	}
	if page.ShardClosed {
		children, cerr := store.GetChildShards(streamName, shardID)
		if cerr != nil {
			return shardPage{}, cerr
		}
		page.ChildShards = children
	}
	return page, nil
}

// getRecordsCore reads one page of records through a shard iterator and
// builds the follow-up iterator at the page's advance position.
func (s *KinesisService) getRecordsCore(reqCtx *request.RequestContext, input GetRecordsInput) (GetRecordsResult, error) {
	if !validateGetRecordsLimit(input.Limit) {
		return GetRecordsResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return GetRecordsResult{}, err
	}

	iterator, err := store.GetShardIterator(input.ShardIterator)
	if err != nil {
		return GetRecordsResult{}, s.mapStoreError(err)
	}

	// The stream lookup is not advisory: an iterator can outlive its stream
	// (created just before a delete), and the encryption type below reads
	// the stream record.
	stream, err := store.GetStream(iterator.StreamName)
	if err != nil {
		return GetRecordsResult{}, s.mapStoreError(err)
	}

	// A DryRun request that passes every check answers the model's
	// dedicated error — a success payload here would carry a null
	// NextShardIterator, which the member's documentation reads as the
	// shard having closed.
	if input.DryRun {
		return GetRecordsResult{}, ErrDryRunOperation
	}

	includeStart := iterator.IteratorType == "AT_SEQUENCE_NUMBER"
	page, err := s.readShardPage(store, iterator.StreamName, iterator.ShardID, iterator.SequenceNumber, input.Limit, includeStart)
	if err != nil {
		return GetRecordsResult{}, s.mapStoreError(err)
	}
	records := page.Records

	// The model's contract: an open shard always returns a continuation,
	// and a closed one returns null only once the reader has consumed its
	// records. The follow-up sits at the advance position — the last
	// consumed record, aged-out records included — so a fully expired page
	// moves the reader to the tip instead of looping on the same position.
	// A creation failure propagates: a page without its continuation would
	// read as a closed shard to SDK clients.
	var nextIterator interface{}
	if len(records) > 0 || !page.ShardClosed {
		newIterator, err := store.CreateShardIterator(
			iterator.StreamName,
			iterator.ShardID,
			"AFTER_SEQUENCE_NUMBER",
			page.LastSeqNum,
			nil,
		)
		if err != nil {
			return GetRecordsResult{}, s.mapStoreError(err)
		}
		nextIterator = newIterator.IteratorID
	}

	return GetRecordsResult{
		Records:            records,
		NextShardIterator:  nextIterator,
		MillisBehindLatest: page.MillisBehindLatest,
		EncryptionType:     resolveEncryptionType(stream),
		ChildShards:        page.ChildShards,
	}, nil
}

// getShardIteratorCore resolves the stream (the ARN wins), validates the
// iterator request and creates the shard iterator.
func (s *KinesisService) getShardIteratorCore(reqCtx *request.RequestContext, input GetShardIteratorInput) (string, error) {
	if input.ShardId == "" || !validateShardId(input.ShardId) || !validateIteratorType(input.ShardIteratorType) {
		return "", ErrInvalidArgument
	}
	if !validateStartingPosition(input.ShardIteratorType, input.StartingSequenceNumber, input.Timestamp != nil) {
		return "", ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return "", err
	}

	// A DryRun request runs every check — the shard must resolve — then
	// answers the model's dedicated error instead of creating the cursor.
	if input.DryRun {
		if _, err := store.GetShard(streamName, input.ShardId); err != nil {
			return "", s.mapStoreError(err)
		}
		return "", ErrDryRunOperation
	}

	iterator, err := store.CreateShardIterator(streamName, input.ShardId, input.ShardIteratorType, input.StartingSequenceNumber, input.Timestamp)
	if err != nil {
		return "", s.mapStoreError(err)
	}

	return iterator.IteratorID, nil
}

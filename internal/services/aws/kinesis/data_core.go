package kinesis

import (
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// PutRecordInput is the transport-agnostic input for PutRecord. The record
// members are validated inside the Core in the original precedence order
// (partition key before the stream fetch, data size after it).
type PutRecordInput struct {
	StreamName      string
	StreamARN       string
	Data            string
	PartitionKey    string
	ExplicitHashKey string
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
}

// putRecordCore validates and writes a single record into the stream.
func (s *KinesisService) putRecordCore(store *kinesisstore.KinesisStore, input PutRecordInput) (PutRecordResult, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return PutRecordResult{}, err
	}

	if !validatePartitionKey(input.PartitionKey) {
		return PutRecordResult{}, ErrInvalidArgument
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return PutRecordResult{}, s.mapStoreError(err)
	}

	if !validateRecordDataSize(input.Data, stream.MaxRecordSizeInKiB) {
		return PutRecordResult{}, ErrValidation
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
func (s *KinesisService) putRecordsCore(store *kinesisstore.KinesisStore, input PutRecordsInput) (PutRecordsResult, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return PutRecordsResult{}, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return PutRecordsResult{}, s.mapStoreError(err)
	}

	if input.Records == nil {
		return PutRecordsResult{}, ErrInvalidArgument
	}

	var requests []kinesisstore.PutRecordRequest
	if recordsList, ok := input.Records.([]interface{}); ok {
		for _, r := range recordsList {
			rm, ok := r.(map[string]interface{})
			if !ok {
				continue
			}
			data, _ := rm["Data"].(string)
			pk, _ := rm["PartitionKey"].(string)
			ehk, _ := rm["ExplicitHashKey"].(string)
			if !validatePartitionKey(pk) {
				return PutRecordsResult{}, ErrValidation
			}
			if !validateRecordDataSize(data, stream.MaxRecordSizeInKiB) {
				return PutRecordsResult{}, ErrValidation
			}
			requests = append(requests, kinesisstore.PutRecordRequest{
				Data:            data,
				PartitionKey:    pk,
				ExplicitHashKey: ehk,
			})
		}
	}

	if len(requests) == 0 {
		return PutRecordsResult{}, ErrInvalidArgument
	}

	if len(requests) > 500 {
		return PutRecordsResult{}, ErrValidation
	}

	results, err := store.PutRecords(streamName, requests)
	if err != nil {
		return PutRecordsResult{}, s.mapStoreError(err)
	}

	return PutRecordsResult{
		Results:        results,
		EncryptionType: resolveEncryptionType(stream),
	}, nil
}

// getRecordsCore reads one page of records through a shard iterator,
// applying the retention window, building the follow-up iterator and
// collecting child shards for closed shards.
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

	includeStart := iterator.IteratorType == "AT_SEQUENCE_NUMBER"
	records, lastSeqNum, err := store.GetRecords(
		iterator.StreamName,
		iterator.ShardID,
		iterator.SequenceNumber,
		input.Limit,
		includeStart,
	)
	if err != nil {
		return GetRecordsResult{}, s.mapStoreError(err)
	}

	stream, _ := store.GetStream(iterator.StreamName)
	encType := resolveEncryptionType(stream)

	// Filter out records older than the retention period.
	retentionCutoff := time.Now().UTC().Add(-time.Duration(stream.RetentionPeriodHours) * time.Hour)
	if len(records) > 0 {
		filtered := records[:0]
		for _, r := range records {
			if !r.ApproximateArrivalTimestamp.Before(retentionCutoff) {
				filtered = append(filtered, r)
			}
		}
		records = filtered
		if len(records) == 0 {
			lastSeqNum = iterator.SequenceNumber
		} else {
			lastSeqNum = records[len(records)-1].SequenceNumber
		}
	}

	// Calculate MillisBehindLatest from the last record's arrival time.
	var millisBehindLatest int64
	if len(records) > 0 {
		last := records[len(records)-1]
		millisBehindLatest = time.Since(last.ApproximateArrivalTimestamp).Milliseconds()
		if millisBehindLatest < 0 {
			millisBehindLatest = 0
		}
	}

	var nextIterator interface{}
	shard, _ := store.GetShard(iterator.StreamName, iterator.ShardID)
	shardClosed := shard != nil && shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber != ""

	if len(records) > 0 {
		newIterator, err := store.CreateShardIterator(
			iterator.StreamName,
			iterator.ShardID,
			"AFTER_SEQUENCE_NUMBER",
			lastSeqNum,
			nil,
		)
		if err == nil {
			nextIterator = newIterator.IteratorID
		}
	} else if !shardClosed {
		newIterator, err := store.CreateShardIterator(
			iterator.StreamName,
			iterator.ShardID,
			"AFTER_SEQUENCE_NUMBER",
			iterator.SequenceNumber,
			nil,
		)
		if err == nil {
			nextIterator = newIterator.IteratorID
		}
	}

	var childShards []*kinesisstore.Shard
	// When the shard is closed (split or merged), collect the child shards so
	// consumers know which shards to read from next.
	if shardClosed {
		childShards, _ = store.GetChildShards(iterator.StreamName, iterator.ShardID)
	}

	return GetRecordsResult{
		Records:            records,
		NextShardIterator:  nextIterator,
		MillisBehindLatest: millisBehindLatest,
		EncryptionType:     encType,
		ChildShards:        childShards,
	}, nil
}

// getShardIteratorCore resolves the stream (the ARN wins), validates the
// iterator request and creates the shard iterator.
func (s *KinesisService) getShardIteratorCore(store *kinesisstore.KinesisStore, input GetShardIteratorInput) (string, error) {
	streamName := input.StreamName

	if input.StreamARN != "" {
		stream, err := store.GetStreamByARN(input.StreamARN)
		if err != nil {
			return "", s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" || input.ShardId == "" || !validateIteratorType(input.ShardIteratorType) {
		return "", ErrInvalidArgument
	}

	iterator, err := store.CreateShardIterator(streamName, input.ShardId, input.ShardIteratorType, input.StartingSequenceNumber, input.Timestamp)
	if err != nil {
		return "", s.mapStoreError(err)
	}

	return iterator.IteratorID, nil
}

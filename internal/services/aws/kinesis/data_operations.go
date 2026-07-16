package kinesis

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// validateRecordDataSize decodes the base64-encoded Data and checks the
// decoded byte length against the stream's max record size. AWS measures
// record size on the raw payload, not the base64-encoded representation.
func validateRecordDataSize(b64Data string, maxKiB int32) bool {
	decoded, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return false
	}
	maxBytes := int(maxKiB) * 1024
	if maxBytes <= 0 {
		maxBytes = 1048576 // 1 MiB default
	}
	return len(decoded) <= maxBytes
}

// PutRecord writes a single data record into a Kinesis stream.
func (s *KinesisService) PutRecord(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	data := request.GetParamLowerFirst(req.Parameters, "Data")
	partitionKey := request.GetParamLowerFirst(req.Parameters, "PartitionKey")
	explicitHashKey := request.GetParamLowerFirst(req.Parameters, "ExplicitHashKey")
	if partitionKey == "" {
		return nil, ErrInvalidArgument
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if !validateRecordDataSize(data, stream.MaxRecordSizeInKiB) {
		return nil, ErrValidation
	}

	record, targetShardID, err := store.PutRecordWithShardSelection(streamName, partitionKey, data, explicitHashKey)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ShardId":        targetShardID,
		"SequenceNumber": record.SequenceNumber,
		"EncryptionType": resolveEncryptionType(stream),
	}, nil
}

// PutRecords writes multiple data records into a Kinesis stream.
func (s *KinesisService) PutRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	recordsRaw := req.Parameters["Records"]
	if recordsRaw == nil {
		return nil, ErrInvalidArgument
	}

	var requests []kinesisstore.PutRecordRequest
	switch v := recordsRaw.(type) {
	case []interface{}:
		for _, r := range v {
			rm, ok := r.(map[string]interface{})
			if !ok {
				continue
			}
			data, _ := rm["Data"].(string)
			pk, _ := rm["PartitionKey"].(string)
			ehk, _ := rm["ExplicitHashKey"].(string)
			if pk == "" {
				return nil, ErrValidation
			}
			if len(pk) > 256 {
				return nil, ErrValidation
			}
			if !validateRecordDataSize(data, stream.MaxRecordSizeInKiB) {
				return nil, ErrValidation
			}
			requests = append(requests, kinesisstore.PutRecordRequest{
				Data:            data,
				PartitionKey:    pk,
				ExplicitHashKey: ehk,
			})
		}
	}

	if len(requests) == 0 {
		return nil, ErrInvalidArgument
	}

	if len(requests) > 500 {
		return nil, ErrValidation
	}

	results, err := store.PutRecords(streamName, requests)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	encType := resolveEncryptionType(stream)

	var failedCount int32
	formattedResults := make([]map[string]interface{}, len(results))
	for i, r := range results {
		entry := map[string]interface{}{
			"SequenceNumber": r.SequenceNumber,
			"ShardId":        r.ShardID,
			"EncryptionType": encType,
		}
		if r.ErrorCode != "" {
			failedCount++
			entry["ErrorCode"] = r.ErrorCode
			entry["ErrorMessage"] = r.ErrorMessage
		}
		formattedResults[i] = entry
	}

	return map[string]interface{}{
		"FailedRecordCount": failedCount,
		"Records":           formattedResults,
		"EncryptionType":    encType,
	}, nil
}

// GetRecords retrieves records from a Kinesis stream shard.
func (s *KinesisService) GetRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	iteratorID := request.GetParamLowerFirst(req.Parameters, "ShardIterator")
	limit := int32(10000)
	if val := request.GetParamLowerFirst(req.Parameters, "Limit"); val != "" {
		limit = int32(request.GetIntParam(req.Parameters, "Limit"))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	iterator, err := store.GetShardIterator(iteratorID)
	if err != nil {
		return nil, ErrExpiredIterator
	}

	includeStart := iterator.IteratorType == "AT_SEQUENCE_NUMBER"
	records, lastSeqNum, err := store.GetRecords(
		iterator.StreamName,
		iterator.ShardID,
		iterator.SequenceNumber,
		limit,
		includeStart,
	)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream, _ := store.GetStream(iterator.StreamName)
	encType := resolveEncryptionType(stream)

	// L4: Filter out records older than the retention period.
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

	// M1: Calculate MillisBehindLatest from the last record's arrival time.
	var millisBehindLatest int64
	if len(records) > 0 {
		last := records[len(records)-1]
		millisBehindLatest = time.Since(last.ApproximateArrivalTimestamp).Milliseconds()
		if millisBehindLatest < 0 {
			millisBehindLatest = 0
		}
	}

	formattedRecords := make([]map[string]interface{}, len(records))
	for i, r := range records {
		formattedRecords[i] = map[string]interface{}{
			"SequenceNumber":              r.SequenceNumber,
			"ApproximateArrivalTimestamp": r.ApproximateArrivalTimestamp.Unix(),
			"Data":                        r.Data,
			"PartitionKey":                r.PartitionKey,
			"EncryptionType":              encType,
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

	return map[string]interface{}{
		"Records":            formattedRecords,
		"NextShardIterator":  nextIterator,
		"MillisBehindLatest": millisBehindLatest,
	}, nil
}

// GetShardIterator gets a shard iterator for reading from a Kinesis stream shard.
func (s *KinesisService) GetShardIterator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	shardID := request.GetParamLowerFirst(req.Parameters, "ShardId")
	iteratorType := request.GetParamLowerFirst(req.Parameters, "ShardIteratorType")
	startingSeqNum := request.GetParamLowerFirst(req.Parameters, "StartingSequenceNumber")

	var timestamp *time.Time
	if ts := request.GetParamLowerFirst(req.Parameters, "Timestamp"); ts != "" {
		if unixTs, err := strconv.ParseInt(ts, 10, 64); err == nil {
			t := time.Unix(unixTs, 0).UTC()
			timestamp = &t
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if streamARN != "" {
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" || shardID == "" || iteratorType == "" {
		return nil, ErrInvalidArgument
	}

	iterator, err := store.CreateShardIterator(streamName, shardID, iteratorType, startingSeqNum, timestamp)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ShardIterator": iterator.IteratorID,
	}, nil
}

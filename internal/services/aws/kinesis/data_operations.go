package kinesis

import (
	"context"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// PutRecord writes a single data record into a Kinesis stream.
func (s *KinesisService) PutRecord(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	data := request.GetParamLowerFirst(req.Parameters, "Data")
	partitionKey := request.GetParamLowerFirst(req.Parameters, "PartitionKey")

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

	if streamName == "" || data == "" || partitionKey == "" {
		return nil, ErrInvalidArgument
	}

	shards, err := store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var activeShards []*kinesisstore.Shard
	for _, shard := range shards {
		if shard.SequenceNumberRange != nil && shard.SequenceNumberRange.EndingSequenceNumber == "" {
			activeShards = append(activeShards, shard)
		}
	}

	if len(activeShards) == 0 {
		return nil, ErrResourceNotFound
	}

	targetShard := store.SelectShardByPartitionKey(activeShards, partitionKey)
	if targetShard == nil {
		targetShard = activeShards[0]
	}

	record, err := store.PutRecord(streamName, targetShard.ShardID, partitionKey, data)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"ShardId":        targetShard.ShardID,
		"SequenceNumber": record.SequenceNumber,
		"EncryptionType": "NONE",
	}, nil
}

// PutRecords writes multiple data records into a Kinesis stream.
func (s *KinesisService) PutRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

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

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	recordsRaw := req.Parameters["Records"]
	if recordsRaw == nil {
		return nil, ErrInvalidArgument
	}

	var requests []kinesisstore.PutRecordRequest
	switch v := recordsRaw.(type) {
	case []interface{}:
		for _, r := range v {
			if rm, ok := r.(map[string]interface{}); ok {
				data := ""
				if d, ok := rm["Data"].(string); ok {
					data = d
				}
				pk := ""
				if p, ok := rm["PartitionKey"].(string); ok {
					pk = p
				}
				if pk == "" {
					return nil, ErrValidation
				}
				if len(pk) > 256 {
					return nil, ErrValidation
				}
				if len(data) > 1048576 {
					return nil, ErrValidation
				}
				requests = append(requests, kinesisstore.PutRecordRequest{
					Data:         data,
					PartitionKey: pk,
				})
			}
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

	var failedCount int32
	var formattedResults []map[string]interface{}
	for _, r := range results {
		entry := map[string]interface{}{
			"SequenceNumber": r.SequenceNumber,
			"ShardId":        r.ShardID,
		}
		if r.ErrorCode != "" {
			failedCount++
			entry["ErrorCode"] = r.ErrorCode
			entry["ErrorMessage"] = r.ErrorMessage
		}
		formattedResults = append(formattedResults, entry)
	}

	return map[string]interface{}{
		"FailedRecordCount": failedCount,
		"Records":           formattedResults,
		"EncryptionType":    "NONE",
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

	formattedRecords := make([]map[string]interface{}, 0)
	for _, r := range records {
		formattedRecords = append(formattedRecords, map[string]interface{}{
			"SequenceNumber":              r.SequenceNumber,
			"ApproximateArrivalTimestamp": r.ApproximateArrivalTimestamp.Unix(),
			"Data":                        r.Data,
			"PartitionKey":                r.PartitionKey,
			"EncryptionType":              "NONE",
		})
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
		"MillisBehindLatest": 0,
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

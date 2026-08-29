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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putRecordCore(store, PutRecordInput{
		StreamName:      request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:       request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		Data:            request.GetParamLowerFirst(req.Parameters, "Data"),
		PartitionKey:    request.GetParamLowerFirst(req.Parameters, "PartitionKey"),
		ExplicitHashKey: request.GetParamLowerFirst(req.Parameters, "ExplicitHashKey"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ShardId":        result.ShardID,
		"SequenceNumber": result.SequenceNumber,
		"EncryptionType": result.EncryptionType,
	}, nil
}

// PutRecords writes multiple data records into a Kinesis stream.
func (s *KinesisService) PutRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putRecordsCore(store, PutRecordsInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		Records:    req.Parameters["Records"],
	})
	if err != nil {
		return nil, err
	}

	var failedCount int32
	formattedResults := make([]map[string]interface{}, len(result.Results))
	for i, r := range result.Results {
		entry := map[string]interface{}{
			"SequenceNumber": r.SequenceNumber,
			"ShardId":        r.ShardID,
			"EncryptionType": result.EncryptionType,
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
		"EncryptionType":    result.EncryptionType,
	}, nil
}

// GetRecords retrieves records from a Kinesis stream shard.
func (s *KinesisService) GetRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit := int32(10000)
	if _, ok := req.Parameters["Limit"]; ok {
		limit = int32(request.GetIntParam(req.Parameters, "Limit"))
	}

	result, err := s.getRecordsCore(reqCtx, GetRecordsInput{
		ShardIterator: request.GetParamLowerFirst(req.Parameters, "ShardIterator"),
		Limit:         limit,
	})
	if err != nil {
		return nil, err
	}

	formattedRecords := make([]map[string]interface{}, len(result.Records))
	for i, r := range result.Records {
		formattedRecords[i] = map[string]interface{}{
			"SequenceNumber":              r.SequenceNumber,
			"ApproximateArrivalTimestamp": r.ApproximateArrivalTimestamp.Unix(),
			"Data":                        r.Data,
			"PartitionKey":                r.PartitionKey,
			"EncryptionType":              result.EncryptionType,
		}
	}

	resp := map[string]interface{}{
		"Records":            formattedRecords,
		"NextShardIterator":  result.NextShardIterator,
		"MillisBehindLatest": result.MillisBehindLatest,
	}

	// When the shard is closed (split or merged), include ChildShards so
	// consumers know which shards to read from next.
	if len(result.ChildShards) > 0 {
		resp["ChildShards"] = formatChildShards(result.ChildShards)
	}

	return resp, nil
}

// GetShardIterator gets a shard iterator for reading from a Kinesis stream shard.
func (s *KinesisService) GetShardIterator(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var timestamp *time.Time
	if ts := request.GetParamLowerFirst(req.Parameters, "Timestamp"); ts != "" {
		if unixTs, err := strconv.ParseInt(ts, 10, 64); err == nil {
			t := time.Unix(unixTs, 0).UTC()
			timestamp = &t
		}
	}

	iteratorID, err := s.getShardIteratorCore(store, GetShardIteratorInput{
		StreamName:             request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:              request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardId:                request.GetParamLowerFirst(req.Parameters, "ShardId"),
		ShardIteratorType:      request.GetParamLowerFirst(req.Parameters, "ShardIteratorType"),
		StartingSequenceNumber: request.GetParamLowerFirst(req.Parameters, "StartingSequenceNumber"),
		Timestamp:              timestamp,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ShardIterator": iteratorID,
	}, nil
}

// formatChildShards formats child shards for GetRecords and SubscribeToShard
// responses. Each ChildShard contains ShardId, ParentShards (list), and
// HashKeyRange.
func formatChildShards(shards []*kinesisstore.Shard) []interface{} {
	result := make([]interface{}, 0, len(shards))
	for _, shard := range shards {
		parentShards := []string{}
		if shard.ParentShardID != "" {
			parentShards = append(parentShards, shard.ParentShardID)
		}
		if shard.AdjacentParentShardID != "" {
			parentShards = append(parentShards, shard.AdjacentParentShardID)
		}
		m := map[string]interface{}{
			"ShardId":      shard.ShardID,
			"ParentShards": parentShards,
			"HashKeyRange": map[string]interface{}{
				"StartingHashKey": shard.HashKeyRange.StartingHashKey,
				"EndingHashKey":   shard.HashKeyRange.EndingHashKey,
			},
		}
		result = append(result, m)
	}
	return result
}

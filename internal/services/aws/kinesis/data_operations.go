package kinesis

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// PutRecord writes a single data record into a Kinesis stream.
func (s *KinesisService) PutRecord(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Data is a required blob member — the wire form is a base64 string, so
	// the strict string read rejects a non-string value here instead of
	// letting the lenient coercion degrade it to an empty payload; its
	// present flag is the required-member presence (a JSON null reads as
	// absent, the protocol's dropped-null rule).
	data, hasData, err := strictStringParam(req.Parameters, "Data")
	if err != nil {
		return nil, err
	}
	result, err := s.putRecordCore(reqCtx, PutRecordInput{
		StreamName:                request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:                 request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		Data:                      data,
		HasData:                   hasData,
		PartitionKey:              request.GetParamLowerFirst(req.Parameters, "PartitionKey"),
		ExplicitHashKey:           request.GetParamLowerFirst(req.Parameters, "ExplicitHashKey"),
		SequenceNumberForOrdering: request.GetParamLowerFirst(req.Parameters, "SequenceNumberForOrdering"),
		DryRun:                    request.GetBoolParam(req.Parameters, "DryRun"),
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
	result, err := s.putRecordsCore(reqCtx, PutRecordsInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		Records:    req.Parameters["Records"],
		DryRun:     request.GetBoolParam(req.Parameters, "DryRun"),
	})
	if err != nil {
		return nil, err
	}

	var failedCount int32
	formattedResults := make([]map[string]interface{}, len(result.Results))
	for i, r := range result.Results {
		// The result-entry shape carries the write receipt and the failure
		// identity only — encryption reports at the batch level. The model's
		// own form: a successful record includes SequenceNumber and ShardId,
		// a failed one includes ErrorCode and ErrorMessage.
		var entry map[string]interface{}
		if r.ErrorCode != "" {
			failedCount++
			entry = map[string]interface{}{
				"ErrorCode":    r.ErrorCode,
				"ErrorMessage": r.ErrorMessage,
			}
		} else {
			entry = map[string]interface{}{
				"SequenceNumber": r.SequenceNumber,
				"ShardId":        r.ShardID,
			}
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
	limit := int32(kinesisstore.DefaultGetRecordsLimit)
	limitValue, hasLimit, err := strictIntParam(req.Parameters, "Limit")
	if err != nil {
		return nil, err
	}
	if hasLimit {
		limit = int32(limitValue)
	}

	result, err := s.getRecordsCore(reqCtx, GetRecordsInput{
		ShardIterator: request.GetParamLowerFirst(req.Parameters, "ShardIterator"),
		Limit:         limit,
		DryRun:        request.GetBoolParam(req.Parameters, "DryRun"),
	})
	if err != nil {
		return nil, err
	}

	formattedRecords := make([]map[string]interface{}, len(result.Records))
	for i, r := range result.Records {
		formattedRecords[i] = map[string]interface{}{
			"SequenceNumber":              r.SequenceNumber,
			"ApproximateArrivalTimestamp": formatEpochSeconds(r.ApproximateArrivalTimestamp),
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
	var timestamp *time.Time
	if ts, present, err := strictTimestampParam(req.Parameters, "Timestamp"); err != nil {
		return nil, err
	} else if present {
		t, err := parseTimestampMember(ts)
		if err != nil {
			return nil, err
		}
		timestamp = &t
	}

	iteratorID, err := s.getShardIteratorCore(reqCtx, GetShardIteratorInput{
		StreamName:             request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:              request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardId:                request.GetParamLowerFirst(req.Parameters, "ShardId"),
		ShardIteratorType:      request.GetParamLowerFirst(req.Parameters, "ShardIteratorType"),
		StartingSequenceNumber: request.GetParamLowerFirst(req.Parameters, "StartingSequenceNumber"),
		Timestamp:              timestamp,
		DryRun:                 request.GetBoolParam(req.Parameters, "DryRun"),
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
// HashKeyRange — a required output member, so a shard record without its
// hash-key submessage formats with empty-string bounds instead of panicking,
// the same absent-submessage reading the twin formatter and the store's own
// guards take.
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
		startHash, endHash := "", ""
		if shard.HashKeyRange != nil {
			startHash, endHash = shard.HashKeyRange.StartingHashKey, shard.HashKeyRange.EndingHashKey
		}
		m := map[string]interface{}{
			"ShardId":      shard.ShardID,
			"ParentShards": parentShards,
			"HashKeyRange": map[string]interface{}{
				"StartingHashKey": startHash,
				"EndingHashKey":   endHash,
			},
		}
		result = append(result, m)
	}
	return result
}

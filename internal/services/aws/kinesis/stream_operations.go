package kinesis

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// CreateStream creates a new Kinesis stream.
func (s *KinesisService) CreateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, hasMaxRecordSize := req.Parameters["MaxRecordSizeInKiB"]
	_, hasWarmThroughput := req.Parameters["WarmThroughputMiBps"]
	stream, err := s.createStreamCore(store, CreateStreamInput{
		StreamName:             request.GetParamLowerFirst(req.Parameters, "StreamName"),
		ShardCount:             int32(request.GetIntParam(req.Parameters, "ShardCount")),
		StreamMode:             parseStreamModeDetails(req.Parameters),
		MaxRecordSizeInKiB:     int32(request.GetIntParam(req.Parameters, "MaxRecordSizeInKiB")),
		HasMaxRecordSizeInKiB:  hasMaxRecordSize,
		WarmThroughputMiBps:    int32(request.GetIntParam(req.Parameters, "WarmThroughputMiBps")),
		HasWarmThroughputMiBps: hasWarmThroughput,
		Tags:                   tags.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

// DeleteStream deletes a Kinesis stream.
func (s *KinesisService) DeleteStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteStreamCore(store, DeleteStreamInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeStream returns detailed information about a Kinesis stream.
func (s *KinesisService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeStreamCore(store, DescribeStreamInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
	})
	if err != nil {
		return nil, err
	}

	stream := result.Stream
	shards := result.Shards

	return map[string]interface{}{
		"StreamDescription": map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"Shards":                  formatShards(shards),
			"HasMoreShards":           false,
			"RetentionPeriodHours":    stream.RetentionPeriodHours,
			"StreamCreationTimestamp": float64(stream.CreatedAt.Unix()),
			"EnhancedMonitoring":      formatEnhancedMonitoring(stream.EnhancedMonitoring),
			"EncryptionType":          resolveEncryptionType(stream),
			"KeyId":                   stream.KeyID,
		},
	}, nil
}

// DescribeStreamSummary returns summary information about a Kinesis stream.
func (s *KinesisService) DescribeStreamSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeStreamSummaryCore(store, DescribeStreamSummaryInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
	})
	if err != nil {
		return nil, err
	}

	stream := result.Stream

	return map[string]interface{}{
		"StreamDescriptionSummary": map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"ConsumerCount":           stream.ConsumerCount,
			"OpenShardCount":          stream.ShardCount,
			"RetentionPeriodHours":    stream.RetentionPeriodHours,
			"StreamCreationTimestamp": float64(stream.CreatedAt.Unix()),
			"EnhancedMonitoring":      formatEnhancedMonitoring(stream.EnhancedMonitoring),
			"EncryptionType":          resolveEncryptionType(stream),
			"KeyId":                   stream.KeyID,
			"MaxRecordSizeInKiB":      stream.MaxRecordSizeInKiB,
		},
	}, nil
}

// ListStreams lists the Kinesis streams.
func (s *KinesisService) ListStreams(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, hasLimit := req.Parameters["Limit"]
	result, err := s.listStreamsCore(store, ListStreamsInput{
		ExclusiveStartStreamName: request.GetStringParam(req.Parameters, "ExclusiveStartStreamName"),
		Limit:                    request.GetIntParam(req.Parameters, "Limit"),
		HasLimit:                 hasLimit,
		NextToken:                request.GetStringParam(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	streamNames := make([]string, 0, len(result.Streams))
	streamSummaries := make([]map[string]interface{}, 0, len(result.Streams))
	for _, stream := range result.Streams {
		streamNames = append(streamNames, stream.StreamName)
		streamSummaries = append(streamSummaries, map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"StreamCreationTimestamp": float64(stream.CreatedAt.Unix()),
		})
	}

	resp := map[string]interface{}{
		"StreamNames":     streamNames,
		"StreamSummaries": streamSummaries,
		"HasMoreStreams":  result.IsTruncated,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// UpdateStreamMode updates the stream mode of a Kinesis stream.
func (s *KinesisService) UpdateStreamMode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, hasWarmThroughput := req.Parameters["WarmThroughputMiBps"]
	result, err := s.updateStreamModeCore(reqCtx, UpdateStreamModeInput{
		StreamARN:           request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		StreamMode:          parseStreamModeDetails(req.Parameters),
		WarmThroughputMiBps: int32(request.GetIntParam(req.Parameters, "WarmThroughputMiBps")),
		HasWarmThroughput:   hasWarmThroughput,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN": result.StreamARN,
	}, nil
}

func formatShards(shards []*kinesisstore.Shard) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(shards))
	for _, shard := range shards {
		m := map[string]interface{}{
			"ShardId": shard.ShardID,
			"HashKeyRange": map[string]interface{}{
				"StartingHashKey": shard.HashKeyRange.StartingHashKey,
				"EndingHashKey":   shard.HashKeyRange.EndingHashKey,
			},
			"SequenceNumberRange": map[string]interface{}{
				"StartingSequenceNumber": shard.SequenceNumberRange.StartingSequenceNumber,
			},
		}
		if shard.ParentShardID != "" {
			m["ParentShardId"] = shard.ParentShardID
		}
		if shard.AdjacentParentShardID != "" {
			m["AdjacentParentShardId"] = shard.AdjacentParentShardID
		}
		if shard.SequenceNumberRange.EndingSequenceNumber != "" {
			m["SequenceNumberRange"].(map[string]interface{})["EndingSequenceNumber"] = shard.SequenceNumberRange.EndingSequenceNumber
		}
		result = append(result, m)
	}
	return result
}

func formatEnhancedMonitoring(em []kinesisstore.EnhancedMonitoring) []map[string]interface{} {
	result := make([]map[string]interface{}, len(em))
	for i, m := range em {
		metrics := m.ShardLevelMetrics
		if metrics == nil {
			metrics = []string{}
		}
		result[i] = map[string]interface{}{
			"ShardLevelMetrics": metrics,
		}
	}
	return result
}

func formatStreamModeDetails(smd *kinesisstore.StreamModeDetails) map[string]interface{} {
	if smd != nil {
		return map[string]interface{}{
			"StreamMode": string(smd.StreamMode),
		}
	}
	return map[string]interface{}{
		"StreamMode": "PROVISIONED",
	}
}

func mergeMetrics(current, added []string) []string {
	seen := make(map[string]bool, len(current)+len(added))
	for _, m := range current {
		seen[m] = true
	}
	result := make([]string, 0, len(current)+len(added))
	result = append(result, current...)
	for _, m := range added {
		if !seen[m] {
			result = append(result, m)
			seen[m] = true
		}
	}
	return result
}

func subtractMetrics(current, removed []string) []string {
	removeSet := make(map[string]bool, len(removed))
	for _, m := range removed {
		removeSet[m] = true
	}
	result := make([]string, 0, len(current))
	for _, m := range current {
		if !removeSet[m] {
			result = append(result, m)
		}
	}
	return result
}

func (s *KinesisService) mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	mapped := awserrors.MapStoreError(err, storeErrorMappings)
	if mapped != err {
		return mapped
	}
	return ErrInvalidArgument
}

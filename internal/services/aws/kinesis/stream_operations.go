package kinesis

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// CreateStream creates a new Kinesis stream.
func (s *KinesisService) CreateStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	shardCount := int32(request.GetIntParam(req.Parameters, "ShardCount"))
	if shardCount == 0 {
		shardCount = 1
	}
	if shardCount < 1 {
		return nil, ErrInvalidArgument
	}

	streamMode := parseStreamModeDetails(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.CreateStream(streamName, shardCount, streamMode)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if tagList := tags.ParseTags(req.Parameters, "Tags"); len(tagList) > 0 {
		tagMap := make(map[string]string, len(tagList))
		for _, t := range tagList {
			tagMap[t.Key] = t.Value
		}
		if err := store.Tag(streamName, tagMap); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

// DeleteStream deletes a Kinesis stream.
func (s *KinesisService) DeleteStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteStream(streamName); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeStream returns detailed information about a Kinesis stream.
func (s *KinesisService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	shards, err := store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

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
	store, streamName, err := s.resolveStreamName(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

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

	exclusiveStartName := request.GetStringParam(req.Parameters, "ExclusiveStartStreamName")
	if exclusiveStartName == "" {
		exclusiveStartName = request.GetStringParam(req.Parameters, "NextToken")
	}
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit <= 0 || limit > 10000 {
		limit = 100
	}

	result, err := store.ListStreams(storecommon.ListOptions{
		Marker:   exclusiveStartName,
		MaxItems: limit + 1,
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	hasMore := len(result.Items) > limit
	if hasMore {
		result.Items = result.Items[:limit]
	}

	streamNames := make([]string, 0, len(result.Items))
	streamSummaries := make([]map[string]interface{}, 0, len(result.Items))
	for _, stream := range result.Items {
		streamNames = append(streamNames, stream.StreamName)
		streamSummaries = append(streamSummaries, map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"StreamCreationTimestamp": float64(stream.CreatedAt.Unix()),
		})
	}

	nextToken := ""
	if hasMore && len(streamSummaries) > 0 {
		nextToken = streamSummaries[len(streamSummaries)-1]["StreamName"].(string)
	}

	return map[string]interface{}{
		"StreamNames":     streamNames,
		"StreamSummaries": streamSummaries,
		"HasMoreStreams":  hasMore,
		"NextToken":       nextToken,
	}, nil
}

// UpdateStreamMode updates the stream mode of a Kinesis stream.
func (s *KinesisService) UpdateStreamMode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	streamMode := parseStreamModeDetails(req.Parameters)

	if streamARN == "" || streamMode == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.GetStreamByARN(streamARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream.StreamModeDetails = &kinesisstore.StreamModeDetails{StreamMode: streamMode}
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
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

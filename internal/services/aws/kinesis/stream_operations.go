package kinesis

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
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

	streamMode := kinesisstore.StreamModeProvisioned
	streamModeDetails := request.GetMapParam(req.Parameters, "StreamModeDetails")
	if streamModeDetails == nil {
		streamModeDetails = request.GetMapParam(req.Parameters, "streamModeDetails")
	}
	if streamModeDetails != nil {
		if v, ok := streamModeDetails["StreamMode"].(string); ok {
			streamMode = kinesisstore.StreamMode(v)
		} else if v, ok := streamModeDetails["streamMode"].(string); ok {
			streamMode = kinesisstore.StreamMode(v)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := store.CreateStream(streamName, shardCount, streamMode)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if tagList := tags.ParseTags(req.Parameters, "Tags"); len(tagList) > 0 {
		tagMap := make(map[string]string)
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

	if err := store.DeleteStream(streamName); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeStream returns detailed information about a Kinesis stream.
func (s *KinesisService) DescribeStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	shards, err := store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	encryptionType := stream.EncryptionType
	if encryptionType == "" {
		encryptionType = "NONE"
	}

	return map[string]interface{}{
		"StreamDescription": map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"Shards":                  s.formatShards(shards),
			"HasMoreShards":           false,
			"RetentionPeriodHours":    stream.RetentionPeriodHours,
			"StreamCreationTimestamp": float64(stream.CreatedAt.Unix()),
			"EnhancedMonitoring":      formatEnhancedMonitoring(stream.EnhancedMonitoring),
			"EncryptionType":          encryptionType,
			"KeyId":                   stream.KeyID,
		},
	}, nil
}

// DescribeStreamSummary returns summary information about a Kinesis stream.
func (s *KinesisService) DescribeStreamSummary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	encryptionType := stream.EncryptionType
	if encryptionType == "" {
		encryptionType = "NONE"
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
			"EncryptionType":          encryptionType,
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

	streams, err := store.ListStreams()
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	exclusiveStartName := request.GetStringParam(req.Parameters, "ExclusiveStartStreamName")
	// The SDK paginator sends NextToken instead of ExclusiveStartStreamName
	// on subsequent pages; accept both for compatibility.
	if exclusiveStartName == "" {
		exclusiveStartName = request.GetStringParam(req.Parameters, "NextToken")
	}
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	streamNames := []string{}
	streamSummaries := []map[string]interface{}{}
	started := exclusiveStartName == ""
	for _, stream := range streams {
		if !started {
			if stream.StreamName == exclusiveStartName {
				started = true
			}
			continue
		}
		streamNames = append(streamNames, stream.StreamName)
		streamSummaries = append(streamSummaries, map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       formatStreamModeDetails(stream.StreamModeDetails),
			"StreamCreationTimestamp": float64(stream.CreatedAt.Unix()),
		})
		if len(streamSummaries) >= limit {
			break
		}
	}

	hasMore := len(streamSummaries) >= limit
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

	streamModeDetails := request.GetMapParam(req.Parameters, "StreamModeDetails")
	if streamModeDetails == nil {
		streamModeDetails = request.GetMapParam(req.Parameters, "streamModeDetails")
	}
	var streamMode string
	if streamModeDetails != nil {
		if v, ok := streamModeDetails["StreamMode"].(string); ok {
			streamMode = v
		} else if v, ok := streamModeDetails["streamMode"].(string); ok {
			streamMode = v
		}
	}

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

	stream.StreamModeDetails = &kinesisstore.StreamModeDetails{StreamMode: kinesisstore.StreamMode(streamMode)}
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

func (s *KinesisService) formatShards(shards []*kinesisstore.Shard) []map[string]interface{} {
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
	seen := make(map[string]bool)
	for _, m := range current {
		seen[m] = true
	}
	result := make([]string, len(current))
	copy(result, current)
	for _, m := range added {
		if !seen[m] {
			result = append(result, m)
			seen[m] = true
		}
	}
	return result
}

func subtractMetrics(current, removed []string) []string {
	removeSet := make(map[string]bool)
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
	switch {
	case errors.Is(err, kinesisstore.ErrStreamNotFound):
		return ErrResourceNotFound
	case errors.Is(err, kinesisstore.ErrStreamAlreadyExists):
		return ErrResourceInUse
	case errors.Is(err, kinesisstore.ErrShardNotFound):
		return ErrResourceNotFound
	case errors.Is(err, kinesisstore.ErrConsumerNotFound):
		return ErrResourceNotFound
	case errors.Is(err, kinesisstore.ErrConsumerAlreadyExists):
		return ErrResourceAlreadyExists
	case errors.Is(err, kinesisstore.ErrShardClosed):
		return ErrShardClosed
	case errors.Is(err, kinesisstore.ErrNoActiveShards):
		return ErrProvisionedThroughputExceeded
	case errors.Is(err, kinesisstore.ErrInvalidIterator):
		return ErrInvalidIterator
	case errors.Is(err, kinesisstore.ErrExpiredIterator):
		return ErrExpiredIterator
	case errors.Is(err, kinesisstore.ErrShardsNotAdjacent):
		return ErrInvalidArgument
	case errors.Is(err, kinesisstore.ErrOperationNotSupportedOnOnDemandStream):
		return ErrInvalidArgument
	case errors.Is(err, kinesisstore.ErrInvalidShardCount):
		return ErrInvalidArgument
	case errors.Is(err, kinesisstore.ErrInvalidParameter):
		return ErrInvalidArgument
	case errors.Is(err, kinesisstore.ErrResourceNotFound):
		return ErrResourceNotFound
	default:
		return ErrInvalidArgument
	}
}

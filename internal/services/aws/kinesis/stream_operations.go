package kinesis

import (
	"context"
	"errors"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	"vorpalstacks/internal/services/aws/common/tags"
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
		if err := store.TagResource(streamName, tagMap); err != nil {
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

	var streamModeDetails map[string]interface{}
	if stream.StreamModeDetails != nil {
		streamModeDetails = map[string]interface{}{
			"StreamMode": string(stream.StreamModeDetails.StreamMode),
		}
	}

	return map[string]interface{}{
		"StreamDescription": map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"ShardCount":              stream.ShardCount,
			"StreamModeDetails":       streamModeDetails,
			"Shards":                  s.formatShards(shards),
			"HasMoreShards":           false,
			"RetentionPeriodHours":    stream.RetentionPeriodHours,
			"StreamCreationTimestamp": stream.CreatedAt.Unix(),
			"EnhancedMonitoring":      stream.EnhancedMonitoring,
			"EncryptionType":          encryptionType,
			"KeyId":                   stream.KeyID,
			"ConsumerCount":           stream.ConsumerCount,
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

	var streamModeDetails map[string]interface{}
	if stream.StreamModeDetails != nil {
		streamModeDetails = map[string]interface{}{
			"StreamMode": string(stream.StreamModeDetails.StreamMode),
		}
	}

	return map[string]interface{}{
		"StreamDescriptionSummary": map[string]interface{}{
			"StreamName":              stream.StreamName,
			"StreamARN":               stream.StreamARN,
			"StreamStatus":            stream.StreamStatus,
			"StreamModeDetails":       streamModeDetails,
			"ConsumerCount":           stream.ConsumerCount,
			"OpenShardCount":          stream.ShardCount,
			"RetentionPeriodHours":    stream.RetentionPeriodHours,
			"StreamCreationTimestamp": stream.CreatedAt.Unix(),
			"EnhancedMonitoring":      stream.EnhancedMonitoring,
			"EncryptionType":          encryptionType,
			"KeyId":                   stream.KeyID,
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
			"StreamName":        stream.StreamName,
			"StreamARN":         stream.StreamARN,
			"StreamStatus":      stream.StreamStatus,
			"StreamModeDetails": stream.StreamModeDetails,
			"ConsumerCount":     stream.ConsumerCount,
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
	var result []map[string]interface{}
	for _, shard := range shards {
		result = append(result, map[string]interface{}{
			"ShardId":               shard.ShardID,
			"ParentShardId":         shard.ParentShardID,
			"AdjacentParentShardId": shard.AdjacentParentShardID,
			"HashKeyRange": map[string]interface{}{
				"StartingHashKey": shard.HashKeyRange.StartingHashKey,
				"EndingHashKey":   shard.HashKeyRange.EndingHashKey,
			},
			"SequenceNumberRange": map[string]interface{}{
				"StartingSequenceNumber": shard.SequenceNumberRange.StartingSequenceNumber,
				"EndingSequenceNumber":   shard.SequenceNumberRange.EndingSequenceNumber,
			},
		})
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
		return ErrResourceAlreadyExists
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

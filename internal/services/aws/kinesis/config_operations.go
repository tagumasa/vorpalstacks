package kinesis

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// IncreaseStreamRetentionPeriod increases the retention period of a Kinesis stream.
func (s *KinesisService) IncreaseStreamRetentionPeriod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	retentionHours := int32(request.GetIntParam(req.Parameters, "RetentionPeriodHours"))

	if streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if retentionHours < stream.RetentionPeriodHours || retentionHours < 24 || retentionHours > 8760 {
		return nil, ErrInvalidArgument
	}

	stream.RetentionPeriodHours = retentionHours
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamName":                  streamName,
		"CurrentRetentionPeriodHours": stream.RetentionPeriodHours,
		"StreamARN":                   stream.StreamARN,
	}, nil
}

// DecreaseStreamRetentionPeriod decreases the retention period of a Kinesis stream.
func (s *KinesisService) DecreaseStreamRetentionPeriod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	retentionHours := int32(request.GetIntParam(req.Parameters, "RetentionPeriodHours"))

	if streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if retentionHours > stream.RetentionPeriodHours || retentionHours < 24 || retentionHours > 8760 {
		return nil, ErrInvalidArgument
	}

	stream.RetentionPeriodHours = retentionHours
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamName":                  streamName,
		"CurrentRetentionPeriodHours": stream.RetentionPeriodHours,
		"StreamARN":                   stream.StreamARN,
	}, nil
}

// DescribeLimits returns the Kinesis service limits.
func (s *KinesisService) DescribeLimits(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"ShardLimit":               500,
		"OpenShardCount":           0,
		"OnDemandStreamCount":      0,
		"OnDemandStreamCountLimit": 50,
	}, nil
}

// DescribeAccountSettings returns the Kinesis account settings.
func (s *KinesisService) DescribeAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{}, nil
}

// UpdateAccountSettings updates the Kinesis account settings.
func (s *KinesisService) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{}, nil
}

// EnableEnhancedMonitoring enables enhanced monitoring for a Kinesis stream.
func (s *KinesisService) EnableEnhancedMonitoring(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	if streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	requestedMetrics := request.GetStringList(req.Parameters, "ShardLevelMetrics")

	var currentMetrics []string
	if len(stream.EnhancedMonitoring) > 0 {
		currentMetrics = stream.EnhancedMonitoring[0].ShardLevelMetrics
	}
	if currentMetrics == nil {
		currentMetrics = []string{}
	}

	desiredMetrics := mergeMetrics(currentMetrics, requestedMetrics)

	stream.EnhancedMonitoring = []kinesisstore.EnhancedMonitoring{
		{ShardLevelMetrics: desiredMetrics},
	}
	if err := store.UpdateStream(stream); err != nil {
		return nil, fmt.Errorf("failed to update stream monitoring: %w", err)
	}

	return map[string]interface{}{
		"StreamName":               streamName,
		"CurrentShardLevelMetrics": currentMetrics,
		"DesiredShardLevelMetrics": desiredMetrics,
		"StreamARN":                stream.StreamARN,
	}, nil
}

// DisableEnhancedMonitoring disables enhanced monitoring for a Kinesis stream.
func (s *KinesisService) DisableEnhancedMonitoring(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	if streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	requestedMetrics := request.GetStringList(req.Parameters, "ShardLevelMetrics")

	var currentMetrics []string
	if len(stream.EnhancedMonitoring) > 0 {
		currentMetrics = stream.EnhancedMonitoring[0].ShardLevelMetrics
	}
	if currentMetrics == nil {
		currentMetrics = []string{}
	}

	desiredMetrics := subtractMetrics(currentMetrics, requestedMetrics)

	stream.EnhancedMonitoring = []kinesisstore.EnhancedMonitoring{
		{ShardLevelMetrics: desiredMetrics},
	}
	if err := store.UpdateStream(stream); err != nil {
		return nil, fmt.Errorf("failed to update stream monitoring: %w", err)
	}

	return map[string]interface{}{
		"StreamName":               streamName,
		"CurrentShardLevelMetrics": currentMetrics,
		"DesiredShardLevelMetrics": desiredMetrics,
		"StreamARN":                stream.StreamARN,
	}, nil
}

// StartStreamEncryption starts server-side encryption for a Kinesis stream.
func (s *KinesisService) StartStreamEncryption(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	encryptionType := request.GetParamLowerFirst(req.Parameters, "EncryptionType")
	keyID := request.GetParamLowerFirst(req.Parameters, "KeyId")

	if streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream.EncryptionType = encryptionType
	stream.KeyID = keyID
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

// StopStreamEncryption stops server-side encryption for a Kinesis stream.
func (s *KinesisService) StopStreamEncryption(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

	if streamARN != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		stream, err := store.GetStreamByARN(streamARN)
		if err != nil {
			return nil, s.mapStoreError(err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream.EncryptionType = ""
	stream.KeyID = ""
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

// GetResourcePolicy retrieves the resource policy for a Kinesis stream.
func (s *KinesisService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.GetResourcePolicy(resourceARN)
	if err != nil {
		policy = ""
	}

	return map[string]interface{}{
		"Policy": policy,
	}, nil
}

// PutResourcePolicy attaches a resource policy to a Kinesis stream.
func (s *KinesisService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	policy := request.GetParamLowerFirst(req.Parameters, "Policy")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.PutResourcePolicy(resourceARN, policy); err != nil {
		return nil, fmt.Errorf("failed to put resource policy: %w", err)
	}

	return response.EmptyResponse(), nil
}

// DeleteResourcePolicy removes the resource policy from a Kinesis stream.
func (s *KinesisService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteResourcePolicy(resourceARN); err != nil {
		logs.Warn("Failed to delete resource policy", logs.Err(err))
	}

	return response.EmptyResponse(), nil
}

// UpdateMaxRecordSize updates the maximum record size for a Kinesis stream.
func (s *KinesisService) UpdateMaxRecordSize(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	maxRecordSizeInKiB := int32(request.GetIntParam(req.Parameters, "MaxRecordSizeInKiB"))

	if streamARN == "" {
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

	if maxRecordSizeInKiB < 1 || maxRecordSizeInKiB > 1024 {
		return nil, ErrInvalidArgument
	}

	stream.MaxRecordSizeInKiB = maxRecordSizeInKiB
	stream.LastModifiedAt = time.Now()
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UpdateStreamWarmThroughput updates the warm throughput capacity for an on-demand Kinesis stream.
func (s *KinesisService) UpdateStreamWarmThroughput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")
	streamName := request.GetParamLowerFirst(req.Parameters, "StreamName")
	warmThroughputMiBps := int32(request.GetIntParam(req.Parameters, "WarmThroughputMiBps"))

	if streamARN == "" {
		return nil, ErrInvalidArgument
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

	if streamName == "" {
		return nil, ErrInvalidArgument
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	if stream.OnDemandStreamConfig == nil {
		stream.OnDemandStreamConfig = &kinesisstore.OnDemandStreamConfig{}
	}
	stream.OnDemandStreamConfig.OnDemandMode = true
	stream.LastModifiedAt = time.Now()
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN":  streamARN,
		"StreamName": streamName,
		"WarmThroughput": map[string]interface{}{
			"CurrentMiBps": warmThroughputMiBps,
			"TargetMiBps":  warmThroughputMiBps,
		},
	}, nil
}

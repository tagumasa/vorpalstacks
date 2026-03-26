package kinesis

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
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
	return map[string]interface{}{
		"AccountSettingsList": []map[string]interface{}{
			{
				"AccountId":                reqCtx.GetAccountID(),
				"ShardLimit":               500,
				"OnDemandStreamCount":      0,
				"OnDemandStreamCountLimit": 50,
			},
		},
	}, nil
}

// UpdateAccountSettings updates the Kinesis account settings.
func (s *KinesisService) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"AccountId":                reqCtx.GetAccountID(),
		"ShardLimit":               500,
		"OnDemandStreamCount":      0,
		"OnDemandStreamCountLimit": 50,
	}, nil
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

	stream.EnhancedMonitoring = []kinesisstore.EnhancedMonitoring{
		{ShardLevelMetrics: []string{"IncomingBytes", "IncomingRecords", "OutgoingBytes", "OutgoingRecords"}},
	}
	if err := store.UpdateStream(stream); err != nil {
		return nil, fmt.Errorf("failed to update stream monitoring: %w", err)
	}

	return map[string]interface{}{
		"StreamName":               streamName,
		"CurrentShardLevelMetrics": []string{"IncomingBytes", "IncomingRecords", "OutgoingBytes", "OutgoingRecords"},
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

	stream.EnhancedMonitoring = []kinesisstore.EnhancedMonitoring{
		{ShardLevelMetrics: []string{}},
	}
	if err := store.UpdateStream(stream); err != nil {
		return nil, fmt.Errorf("failed to update stream monitoring: %w", err)
	}

	return map[string]interface{}{
		"StreamName":               streamName,
		"CurrentShardLevelMetrics": []string{},
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

	return map[string]interface{}{
		"ResourceARN": resourceARN,
		"Policy":      "",
	}, nil
}

// PutResourcePolicy attaches a resource policy to a Kinesis stream.
func (s *KinesisService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
	}

	return map[string]interface{}{
		"ResourceARN": resourceARN,
	}, nil
}

// DeleteResourcePolicy removes the resource policy from a Kinesis stream.
func (s *KinesisService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceARN")

	if resourceARN == "" {
		return nil, ErrInvalidArgument
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

	return map[string]interface{}{
		"StreamARN":          streamARN,
		"MaxRecordSizeInKiB": maxRecordSizeInKiB,
	}, nil
}

// UpdateStreamWarmThroughput updates the warm throughput capacity for an on-demand Kinesis stream.
func (s *KinesisService) UpdateStreamWarmThroughput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	streamARN := request.GetParamLowerFirst(req.Parameters, "StreamARN")

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

	if stream.OnDemandStreamConfig == nil {
		stream.OnDemandStreamConfig = &kinesisstore.OnDemandStreamConfig{}
	}
	stream.OnDemandStreamConfig.OnDemandMode = true
	stream.LastModifiedAt = time.Now()
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{
		"StreamARN":            streamARN,
		"OnDemandStreamConfig": stream.OnDemandStreamConfig,
	}, nil
}

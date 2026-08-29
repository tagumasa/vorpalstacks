package kinesis

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// IncreaseStreamRetentionPeriod increases the retention period of a Kinesis stream.
func (s *KinesisService) IncreaseStreamRetentionPeriod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, streamName, err := s.updateRetentionPeriodCore(store, UpdateRetentionPeriodInput{
		StreamName:           request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:            request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		RetentionPeriodHours: int32(request.GetIntParam(req.Parameters, "RetentionPeriodHours")),
	}, true)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamName":                  streamName,
		"CurrentRetentionPeriodHours": stream.RetentionPeriodHours,
		"StreamARN":                   stream.StreamARN,
	}, nil
}

// DecreaseStreamRetentionPeriod decreases the retention period of a Kinesis stream.
func (s *KinesisService) DecreaseStreamRetentionPeriod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, streamName, err := s.updateRetentionPeriodCore(store, UpdateRetentionPeriodInput{
		StreamName:           request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:            request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		RetentionPeriodHours: int32(request.GetIntParam(req.Parameters, "RetentionPeriodHours")),
	}, false)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamName":                  streamName,
		"CurrentRetentionPeriodHours": stream.RetentionPeriodHours,
		"StreamARN":                   stream.StreamARN,
	}, nil
}

// DescribeLimits returns the Kinesis service limits.
func (s *KinesisService) DescribeLimits(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeLimitsCore(store)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ShardLimit":               result.ShardLimit,
		"OpenShardCount":           result.OpenShardCount,
		"OnDemandStreamCount":      result.OnDemandStreamCount,
		"OnDemandStreamCountLimit": result.OnDemandStreamCountLimit,
	}, nil
}

// DescribeAccountSettings returns the Kinesis account settings.
func (s *KinesisService) DescribeAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"MinimumThroughputBillingCommitment": map[string]interface{}{
			"Status": "DISABLED",
		},
	}, nil
}

// UpdateAccountSettings updates the Kinesis account settings.
func (s *KinesisService) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	status := "DISABLED"
	if mtbc, ok := req.Parameters["MinimumThroughputBillingCommitment"]; ok {
		if mtbcMap, ok := mtbc.(map[string]interface{}); ok {
			if s, ok := mtbcMap["Status"].(string); ok {
				status = s
			}
		}
	}
	return map[string]interface{}{
		"MinimumThroughputBillingCommitment": map[string]interface{}{
			"Status": status,
		},
	}, nil
}

// EnableEnhancedMonitoring enables enhanced monitoring for a Kinesis stream.
func (s *KinesisService) EnableEnhancedMonitoring(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.enableEnhancedMonitoringCore(store, EnhancedMonitoringInput{
		StreamName:        request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:         request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardLevelMetrics: request.GetStringList(req.Parameters, "ShardLevelMetrics"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamName":               result.StreamName,
		"CurrentShardLevelMetrics": result.CurrentMetrics,
		"DesiredShardLevelMetrics": result.DesiredMetrics,
		"StreamARN":                result.StreamARN,
	}, nil
}

// DisableEnhancedMonitoring disables enhanced monitoring for a Kinesis stream.
func (s *KinesisService) DisableEnhancedMonitoring(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.disableEnhancedMonitoringCore(store, EnhancedMonitoringInput{
		StreamName:        request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:         request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		ShardLevelMetrics: request.GetStringList(req.Parameters, "ShardLevelMetrics"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamName":               result.StreamName,
		"CurrentShardLevelMetrics": result.CurrentMetrics,
		"DesiredShardLevelMetrics": result.DesiredMetrics,
		"StreamARN":                result.StreamARN,
	}, nil
}

// StartStreamEncryption starts server-side encryption for a Kinesis stream.
func (s *KinesisService) StartStreamEncryption(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := s.startStreamEncryptionCore(store, StartStreamEncryptionInput{
		StreamName:     request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:      request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		EncryptionType: request.GetParamLowerFirst(req.Parameters, "EncryptionType"),
		KeyId:          request.GetParamLowerFirst(req.Parameters, "KeyId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

// StopStreamEncryption stops server-side encryption for a Kinesis stream.
func (s *KinesisService) StopStreamEncryption(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stream, err := s.stopStreamEncryptionCore(store, StopStreamEncryptionInput{
		StreamName: request.GetParamLowerFirst(req.Parameters, "StreamName"),
		StreamARN:  request.GetParamLowerFirst(req.Parameters, "StreamARN"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN": stream.StreamARN,
	}, nil
}

// GetResourcePolicy retrieves the resource policy for a Kinesis stream.
func (s *KinesisService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policy, err := s.getResourcePolicyCore(reqCtx, ResourcePolicyInput{
		ResourceARN: request.GetParamLowerFirst(req.Parameters, "ResourceARN"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy": policy,
	}, nil
}

// PutResourcePolicy attaches a resource policy to a Kinesis stream.
func (s *KinesisService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.putResourcePolicyCore(reqCtx, ResourcePolicyInput{
		ResourceARN: request.GetParamLowerFirst(req.Parameters, "ResourceARN"),
		Policy:      request.GetParamLowerFirst(req.Parameters, "Policy"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteResourcePolicy removes the resource policy from a Kinesis stream.
func (s *KinesisService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteResourcePolicyCore(reqCtx, ResourcePolicyInput{
		ResourceARN: request.GetParamLowerFirst(req.Parameters, "ResourceARN"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateMaxRecordSize updates the maximum record size for a Kinesis stream.
func (s *KinesisService) UpdateMaxRecordSize(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.updateMaxRecordSizeCore(reqCtx, UpdateMaxRecordSizeInput{
		StreamARN:          request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		MaxRecordSizeInKiB: int32(request.GetIntParam(req.Parameters, "MaxRecordSizeInKiB")),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateStreamWarmThroughput updates the warm throughput capacity for an on-demand Kinesis stream.
func (s *KinesisService) UpdateStreamWarmThroughput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
		StreamARN:           request.GetParamLowerFirst(req.Parameters, "StreamARN"),
		WarmThroughputMiBps: int32(request.GetIntParam(req.Parameters, "WarmThroughputMiBps")),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"StreamARN":  result.StreamARN,
		"StreamName": result.StreamName,
		"WarmThroughput": map[string]interface{}{
			"CurrentMiBps": result.WarmThroughputMiBps,
			"TargetMiBps":  result.WarmThroughputMiBps,
		},
	}, nil
}

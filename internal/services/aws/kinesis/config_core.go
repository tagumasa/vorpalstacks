package kinesis

import (
	"encoding/json"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// UpdateRetentionPeriodInput is the transport-agnostic input for the
// IncreaseStreamRetentionPeriod/DecreaseStreamRetentionPeriod pair.
type UpdateRetentionPeriodInput struct {
	StreamName           string
	StreamARN            string
	RetentionPeriodHours int32
}

// DescribeLimitsResult carries the aggregated account limits.
type DescribeLimitsResult struct {
	ShardLimit               int32
	OpenShardCount           int32
	OnDemandStreamCount      int32
	OnDemandStreamCountLimit int32
}

// EnhancedMonitoringInput is the transport-agnostic input for the
// EnableEnhancedMonitoring/DisableEnhancedMonitoring pair.
type EnhancedMonitoringInput struct {
	StreamName        string
	StreamARN         string
	ShardLevelMetrics []string
}

// EnhancedMonitoringResult carries the metrics transition reported by the
// enhanced-monitoring pair.
type EnhancedMonitoringResult struct {
	StreamName     string
	CurrentMetrics []string
	DesiredMetrics []string
	StreamARN      string
}

// StartStreamEncryptionInput is the transport-agnostic input for
// StartStreamEncryption.
type StartStreamEncryptionInput struct {
	StreamName     string
	StreamARN      string
	EncryptionType string
	KeyId          string
}

// StopStreamEncryptionInput is the transport-agnostic input for
// StopStreamEncryption.
type StopStreamEncryptionInput struct {
	StreamName string
	StreamARN  string
}

// ResourcePolicyInput is the transport-agnostic input for the resource-policy
// family (GetResourcePolicy/PutResourcePolicy/DeleteResourcePolicy).
type ResourcePolicyInput struct {
	ResourceARN string
	Policy      string
}

// UpdateMaxRecordSizeInput is the transport-agnostic input for
// UpdateMaxRecordSize.
type UpdateMaxRecordSizeInput struct {
	StreamARN          string
	MaxRecordSizeInKiB int32
}

// UpdateStreamWarmThroughputInput is the transport-agnostic input for
// UpdateStreamWarmThroughput.
type UpdateStreamWarmThroughputInput struct {
	StreamARN           string
	WarmThroughputMiBps int32
}

// UpdateStreamWarmThroughputResult carries the warm-throughput transition.
type UpdateStreamWarmThroughputResult struct {
	StreamARN           string
	StreamName          string
	WarmThroughputMiBps int32
}

// updateRetentionPeriodCore applies a retention-period change in the given
// direction: an increase rejects values below the current period, a decrease
// rejects values above it.
func (s *KinesisService) updateRetentionPeriodCore(store *kinesisstore.KinesisStore, input UpdateRetentionPeriodInput, increase bool) (*kinesisstore.Stream, string, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return nil, "", err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, "", s.mapStoreError(err)
	}

	if increase {
		if !validateRetentionPeriod(input.RetentionPeriodHours) || input.RetentionPeriodHours < stream.RetentionPeriodHours {
			return nil, "", ErrInvalidArgument
		}
	} else {
		if !validateRetentionPeriod(input.RetentionPeriodHours) || input.RetentionPeriodHours > stream.RetentionPeriodHours {
			return nil, "", ErrInvalidArgument
		}
	}

	stream.RetentionPeriodHours = input.RetentionPeriodHours
	if err := store.UpdateStream(stream); err != nil {
		return nil, "", s.mapStoreError(err)
	}

	return stream, streamName, nil
}

// describeLimitsCore aggregates the actual shard and on-demand stream counts
// across all streams of the account.
func (s *KinesisService) describeLimitsCore(store *kinesisstore.KinesisStore) (DescribeLimitsResult, error) {
	// Aggregate actual shard and on-demand stream counts across all streams.
	openShardCount := int32(0)
	onDemandStreamCount := int32(0)

	result, err := store.ListStreams(storecommon.ListOptions{MaxItems: 10000})
	if err != nil {
		return DescribeLimitsResult{}, s.mapStoreError(err)
	}
	for _, stream := range result.Items {
		if stream.StreamModeDetails != nil && stream.StreamModeDetails.StreamMode == kinesisstore.StreamModeOnDemand {
			onDemandStreamCount++
		}
		// Use stream.ShardCount (maintained by split/merge/updateShardCount)
		// to avoid an N+1 ListShards query per stream.
		openShardCount += stream.ShardCount
	}

	return DescribeLimitsResult{
		ShardLimit:               int32(500),
		OpenShardCount:           openShardCount,
		OnDemandStreamCount:      onDemandStreamCount,
		OnDemandStreamCountLimit: int32(50),
	}, nil
}

// enableEnhancedMonitoringCore merges the requested metrics into the stream's
// enhanced-monitoring configuration.
func (s *KinesisService) enableEnhancedMonitoringCore(store *kinesisstore.KinesisStore, input EnhancedMonitoringInput) (EnhancedMonitoringResult, error) {
	streamName, stream, currentMetrics, err := s.loadStreamForMonitoring(store, input)
	if err != nil {
		return EnhancedMonitoringResult{}, err
	}

	desiredMetrics := mergeMetrics(currentMetrics, input.ShardLevelMetrics)
	if err := s.storeMonitoringMetrics(store, stream, desiredMetrics); err != nil {
		return EnhancedMonitoringResult{}, err
	}

	return EnhancedMonitoringResult{
		StreamName:     streamName,
		CurrentMetrics: currentMetrics,
		DesiredMetrics: desiredMetrics,
		StreamARN:      stream.StreamARN,
	}, nil
}

// disableEnhancedMonitoringCore subtracts the requested metrics from the
// stream's enhanced-monitoring configuration.
func (s *KinesisService) disableEnhancedMonitoringCore(store *kinesisstore.KinesisStore, input EnhancedMonitoringInput) (EnhancedMonitoringResult, error) {
	streamName, stream, currentMetrics, err := s.loadStreamForMonitoring(store, input)
	if err != nil {
		return EnhancedMonitoringResult{}, err
	}

	desiredMetrics := subtractMetrics(currentMetrics, input.ShardLevelMetrics)
	if err := s.storeMonitoringMetrics(store, stream, desiredMetrics); err != nil {
		return EnhancedMonitoringResult{}, err
	}

	return EnhancedMonitoringResult{
		StreamName:     streamName,
		CurrentMetrics: currentMetrics,
		DesiredMetrics: desiredMetrics,
		StreamARN:      stream.StreamARN,
	}, nil
}

// loadStreamForMonitoring resolves the stream, validates the requested
// metrics against the Smithy enum and returns the current metrics list.
func (s *KinesisService) loadStreamForMonitoring(store *kinesisstore.KinesisStore, input EnhancedMonitoringInput) (string, *kinesisstore.Stream, []string, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return "", nil, nil, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return "", nil, nil, s.mapStoreError(err)
	}

	for _, m := range input.ShardLevelMetrics {
		if !validateShardLevelMetric(m) {
			return "", nil, nil, ErrInvalidArgument
		}
	}

	var currentMetrics []string
	if len(stream.EnhancedMonitoring) > 0 {
		currentMetrics = stream.EnhancedMonitoring[0].ShardLevelMetrics
	}
	if currentMetrics == nil {
		currentMetrics = []string{}
	}

	return streamName, stream, currentMetrics, nil
}

// storeMonitoringMetrics persists the desired metrics list on the stream.
func (s *KinesisService) storeMonitoringMetrics(store *kinesisstore.KinesisStore, stream *kinesisstore.Stream, desiredMetrics []string) error {
	stream.EnhancedMonitoring = []kinesisstore.EnhancedMonitoring{
		{ShardLevelMetrics: desiredMetrics},
	}
	if err := store.UpdateStream(stream); err != nil {
		return s.mapStoreError(err)
	}
	return nil
}

// startStreamEncryptionCore enables server-side encryption on a stream.
func (s *KinesisService) startStreamEncryptionCore(store *kinesisstore.KinesisStore, input StartStreamEncryptionInput) (*kinesisstore.Stream, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return nil, err
	}

	if input.EncryptionType != "KMS" {
		return nil, ErrInvalidArgument
	}
	if !validateKeyId(input.KeyId) {
		return nil, ErrInvalidArgument
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	stream.EncryptionType = input.EncryptionType
	stream.KeyID = input.KeyId
	if err := store.UpdateStream(stream); err != nil {
		return nil, s.mapStoreError(err)
	}

	return stream, nil
}

// stopStreamEncryptionCore disables server-side encryption on a stream.
func (s *KinesisService) stopStreamEncryptionCore(store *kinesisstore.KinesisStore, input StopStreamEncryptionInput) (*kinesisstore.Stream, error) {
	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
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

	return stream, nil
}

// getResourcePolicyCore reads the resource policy attached to a stream ARN,
// returning an empty policy when none is attached.
func (s *KinesisService) getResourcePolicyCore(reqCtx *request.RequestContext, input ResourcePolicyInput) (string, error) {
	if !validateResourceARN(input.ResourceARN) {
		return "", ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	policy, err := store.GetResourcePolicy(input.ResourceARN)
	if err != nil {
		policy = ""
	}

	return policy, nil
}

// putResourcePolicyCore attaches a JSON resource policy to a stream ARN.
func (s *KinesisService) putResourcePolicyCore(reqCtx *request.RequestContext, input ResourcePolicyInput) error {
	if input.Policy == "" {
		return ErrInvalidArgument
	}
	if !validateResourceARN(input.ResourceARN) {
		return ErrInvalidArgument
	}
	if !json.Valid([]byte(input.Policy)) {
		return ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := store.PutResourcePolicy(input.ResourceARN, input.Policy); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// deleteResourcePolicyCore removes the resource policy from a stream ARN.
func (s *KinesisService) deleteResourcePolicyCore(reqCtx *request.RequestContext, input ResourcePolicyInput) error {
	if !validateResourceARN(input.ResourceARN) {
		return ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := store.DeleteResourcePolicy(input.ResourceARN); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// updateMaxRecordSizeCore sets the maximum record size of the stream
// identified by ARN.
func (s *KinesisService) updateMaxRecordSizeCore(reqCtx *request.RequestContext, input UpdateMaxRecordSizeInput) error {
	if input.StreamARN == "" {
		return ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	stream, err := store.GetStreamByARN(input.StreamARN)
	if err != nil {
		return s.mapStoreError(err)
	}

	if !validateMaxRecordSizeInKiB(input.MaxRecordSizeInKiB) {
		return ErrInvalidArgument
	}

	stream.MaxRecordSizeInKiB = input.MaxRecordSizeInKiB
	if err := store.UpdateStream(stream); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// updateStreamWarmThroughputCore sets the warm throughput capacity of the
// on-demand stream identified by ARN.
func (s *KinesisService) updateStreamWarmThroughputCore(reqCtx *request.RequestContext, input UpdateStreamWarmThroughputInput) (UpdateStreamWarmThroughputResult, error) {
	if input.StreamARN == "" {
		return UpdateStreamWarmThroughputResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return UpdateStreamWarmThroughputResult{}, err
	}

	stream, err := store.GetStreamByARN(input.StreamARN)
	if err != nil {
		return UpdateStreamWarmThroughputResult{}, s.mapStoreError(err)
	}

	if stream.OnDemandStreamConfig == nil {
		stream.OnDemandStreamConfig = &kinesisstore.OnDemandStreamConfig{}
	}
	stream.OnDemandStreamConfig.OnDemandMode = true
	stream.WarmThroughputMiBps = input.WarmThroughputMiBps
	if err := store.UpdateStream(stream); err != nil {
		return UpdateStreamWarmThroughputResult{}, s.mapStoreError(err)
	}

	return UpdateStreamWarmThroughputResult{
		StreamARN:           input.StreamARN,
		StreamName:          stream.StreamName,
		WarmThroughputMiBps: input.WarmThroughputMiBps,
	}, nil
}

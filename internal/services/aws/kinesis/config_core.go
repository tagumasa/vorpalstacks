package kinesis

import (
	"encoding/json"
	"errors"

	"vorpalstacks/internal/common/kmsutil"
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
	StreamName     string
	StreamARN      string
	EncryptionType string
	KeyId          string
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
	StreamName             string
	StreamARN              string
	WarmThroughputMiBps    int32
	HasWarmThroughputMiBps bool
}

// UpdateStreamWarmThroughputResult carries the warm-throughput transition.
type UpdateStreamWarmThroughputResult struct {
	StreamARN           string
	StreamName          string
	WarmThroughputMiBps int32
}

// UpdateAccountSettingsInput is the transport-agnostic input for
// UpdateAccountSettings: the nested commitment status is the operation's
// only member, required at both levels.
type UpdateAccountSettingsInput struct {
	Status string
}

// UpdateAccountSettingsResult carries the effective commitment status.
type UpdateAccountSettingsResult struct {
	Status string
}

// updateRetentionPeriodCore applies a retention-period change in the given
// direction: an increase rejects values below the current period, a decrease
// rejects values above it.
func (s *KinesisService) updateRetentionPeriodCore(reqCtx *request.RequestContext, input UpdateRetentionPeriodInput, increase bool) (*kinesisstore.Stream, string, error) {
	if !validateRetentionPeriod(input.RetentionPeriodHours) {
		return nil, "", ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return nil, "", err
	}

	stream, err := store.UpdateStreamFields(streamName, func(stream *kinesisstore.Stream) error {
		// The member documentation is directional: an increase "must be more
		// than the current retention period", a decrease "must be less" —
		// equal values satisfy neither direction and reject.
		if increase {
			if input.RetentionPeriodHours <= stream.RetentionPeriodHours {
				return ErrInvalidArgument
			}
		} else {
			if input.RetentionPeriodHours >= stream.RetentionPeriodHours {
				return ErrInvalidArgument
			}
		}
		stream.RetentionPeriodHours = input.RetentionPeriodHours
		return nil
	})
	if err != nil {
		return nil, "", s.mapStoreError(err)
	}

	// A decrease makes the newly aged-out records immediately inaccessible
	// (the read path's cutoff follows the stored period), and their storage
	// is reclaimed at once rather than waiting for the write-plane trim
	// interval. Best-effort: a reclaim failure leaves invisible records on
	// disk for the interval sweep, never a failed configuration change.
	if !increase {
		_ = store.TrimExpiredRecords(streamName)
	}

	return stream, streamName, nil
}

// describeLimitsCore aggregates the actual shard and on-demand stream counts
// across all streams of the account.
func (s *KinesisService) describeLimitsCore(reqCtx *request.RequestContext) (DescribeLimitsResult, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return DescribeLimitsResult{}, err
	}

	// Aggregate actual shard and on-demand stream counts across all
	// streams. The walk pages at the shared paginator's hard item cap
	// until every stream is counted — a single unbounded request would be
	// silently clamped and undercount beyond the cap, with no truncation
	// signal the shape could carry.
	openShardCount := int32(0)
	onDemandStreamCount := int32(0)

	marker := ""
	for {
		result, err := store.ListStreams(storecommon.ListOptions{Marker: marker, MaxItems: kinesisstore.DescribeLimitsStreamsWalkBound})
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
		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}

	return DescribeLimitsResult{
		ShardLimit:               kinesisstore.ShardLimit,
		OpenShardCount:           openShardCount,
		OnDemandStreamCount:      onDemandStreamCount,
		OnDemandStreamCountLimit: kinesisstore.OnDemandStreamCountLimit,
	}, nil
}

// enableEnhancedMonitoringCore merges the requested metrics into the stream's
// enhanced-monitoring configuration.
func (s *KinesisService) enableEnhancedMonitoringCore(reqCtx *request.RequestContext, input EnhancedMonitoringInput) (EnhancedMonitoringResult, error) {
	result, err := s.applyMonitoringCore(reqCtx, input, mergeMetrics)
	if err != nil {
		return EnhancedMonitoringResult{}, err
	}
	return result, nil
}

// disableEnhancedMonitoringCore subtracts the requested metrics from the
// stream's enhanced-monitoring configuration.
func (s *KinesisService) disableEnhancedMonitoringCore(reqCtx *request.RequestContext, input EnhancedMonitoringInput) (EnhancedMonitoringResult, error) {
	result, err := s.applyMonitoringCore(reqCtx, input, subtractMetrics)
	if err != nil {
		return EnhancedMonitoringResult{}, err
	}
	return result, nil
}

// applyMonitoringCore applies a metrics set operation to the stream's
// enhanced-monitoring configuration as one locked read-modify-write.
func (s *KinesisService) applyMonitoringCore(reqCtx *request.RequestContext, input EnhancedMonitoringInput, op func(current, requested []string) []string) (EnhancedMonitoringResult, error) {
	// The member is required and the MetricsNameList length trait bounds it
	// to 1-7 names on both operations of the pair.
	if !validateShardLevelMetricsList(input.ShardLevelMetrics) {
		return EnhancedMonitoringResult{}, ErrInvalidArgument
	}
	for _, m := range input.ShardLevelMetrics {
		if !validateShardLevelMetric(m) {
			return EnhancedMonitoringResult{}, ErrInvalidArgument
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return EnhancedMonitoringResult{}, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return EnhancedMonitoringResult{}, err
	}

	result := EnhancedMonitoringResult{StreamName: streamName}
	_, err = store.UpdateStreamFields(streamName, func(stream *kinesisstore.Stream) error {
		var currentMetrics []string
		if len(stream.EnhancedMonitoring) > 0 {
			currentMetrics = stream.EnhancedMonitoring[0].ShardLevelMetrics
		}
		if currentMetrics == nil {
			currentMetrics = []string{}
		}
		desiredMetrics := op(currentMetrics, input.ShardLevelMetrics)
		stream.EnhancedMonitoring = []kinesisstore.EnhancedMonitoring{
			{ShardLevelMetrics: desiredMetrics},
		}
		result.CurrentMetrics = currentMetrics
		result.DesiredMetrics = desiredMetrics
		result.StreamARN = stream.StreamARN
		return nil
	})
	if err != nil {
		return EnhancedMonitoringResult{}, s.mapStoreError(err)
	}

	return result, nil
}

// startStreamEncryptionCore enables server-side encryption on a stream.
func (s *KinesisService) startStreamEncryptionCore(reqCtx *request.RequestContext, input StartStreamEncryptionInput) (*kinesisstore.Stream, error) {
	if input.EncryptionType != "KMS" {
		return nil, ErrInvalidArgument
	}
	if !validateKeyId(input.KeyId) {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return nil, err
	}

	// The declared KMS error family is reachable: before the stream adopts
	// a key, the key must resolve, be enabled and carry the right usage —
	// the same cross-service contract SQS's KmsMasterKeyId validates
	// through. A nil checker (KMS not constructed) leaves the member
	// shape-checked alone.
	if s.kmsChecker != nil {
		if err := s.kmsChecker.CheckKey(reqCtx, reqCtx.Region, input.KeyId); err != nil {
			return nil, mapKMSError(err)
		}
	}

	stream, err := store.UpdateStreamFields(streamName, func(stream *kinesisstore.Stream) error {
		stream.EncryptionType = input.EncryptionType
		stream.KeyID = input.KeyId
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return stream, nil
}

// mapKMSError maps the KMS checker's sentinel failures to the identities
// the encryption operations declare; an unmapped failure is invalid input.
func mapKMSError(err error) error {
	switch {
	case errors.Is(err, kmsutil.ErrKeyNotFound):
		return ErrKMSNotFound
	case errors.Is(err, kmsutil.ErrKeyDisabled):
		return ErrKMSDisabled
	case errors.Is(err, kmsutil.ErrKeyInvalidState):
		return ErrKMSInvalidState
	case errors.Is(err, kmsutil.ErrKeyInvalidUsage):
		return ErrKMSAccessDenied
	default:
		return ErrInvalidArgument
	}
}

// stopStreamEncryptionCore disables server-side encryption on a stream.
// Both the EncryptionType and KeyId members are required; the model's
// member documentation fixes the encryption type's only valid value at KMS.
func (s *KinesisService) stopStreamEncryptionCore(reqCtx *request.RequestContext, input StopStreamEncryptionInput) (*kinesisstore.Stream, error) {
	if input.EncryptionType != "KMS" {
		return nil, ErrInvalidArgument
	}
	if !validateKeyId(input.KeyId) {
		return nil, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return nil, err
	}

	stream, err := store.UpdateStreamFields(streamName, func(stream *kinesisstore.Stream) error {
		stream.EncryptionType = ""
		stream.KeyID = ""
		return nil
	})
	if err != nil {
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
		// A stream without a policy is the operation's declared
		// ResourceNotFoundException; any other failure is infrastructure
		// and propagates — never an empty-string success for either.
		return "", s.mapStoreError(err)
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

	if _, err := store.UpdateStreamFields(stream.StreamName, func(stream *kinesisstore.Stream) error {
		stream.MaxRecordSizeInKiB = input.MaxRecordSizeInKiB
		return nil
	}); err != nil {
		return s.mapStoreError(err)
	}

	return nil
}

// updateStreamWarmThroughputCore sets the warm throughput capacity of the
// stream identified by name or ARN. The operation's documentation scopes it
// to on-demand data streams: a provisioned stream rejects instead of being
// silently switched, and the write never touches the stream's mode — the
// warm throughput is a capacity figure the summary reports.
func (s *KinesisService) updateStreamWarmThroughputCore(reqCtx *request.RequestContext, input UpdateStreamWarmThroughputInput) (UpdateStreamWarmThroughputResult, error) {
	// WarmThroughputMiBps is a required member: an absent one is rejected,
	// not read as the range-legal zero that would overwrite the stored
	// target with nothing.
	if !input.HasWarmThroughputMiBps {
		return UpdateStreamWarmThroughputResult{}, ErrInvalidArgument
	}
	if !validateWarmThroughputMiBps(input.WarmThroughputMiBps) {
		return UpdateStreamWarmThroughputResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return UpdateStreamWarmThroughputResult{}, err
	}

	streamName, err := s.resolveStreamNameCore(store, input.StreamName, input.StreamARN)
	if err != nil {
		return UpdateStreamWarmThroughputResult{}, err
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		return UpdateStreamWarmThroughputResult{}, s.mapStoreError(err)
	}
	if stream.StreamModeDetails == nil || stream.StreamModeDetails.StreamMode != kinesisstore.StreamModeOnDemand {
		return UpdateStreamWarmThroughputResult{}, ErrInvalidArgument
	}

	updated, err := store.UpdateStreamFields(streamName, func(stream *kinesisstore.Stream) error {
		stream.WarmThroughputMiBps = input.WarmThroughputMiBps
		return nil
	})
	if err != nil {
		return UpdateStreamWarmThroughputResult{}, s.mapStoreError(err)
	}

	return UpdateStreamWarmThroughputResult{
		StreamARN:           updated.StreamARN,
		StreamName:          updated.StreamName,
		WarmThroughputMiBps: input.WarmThroughputMiBps,
	}, nil
}

// updateAccountSettingsCore validates the commitment status request — the
// status member is required at both nesting levels and takes the InputStatus
// enum's two values — and persists it, so a later Describe reports what the
// Update answered.
func (s *KinesisService) updateAccountSettingsCore(reqCtx *request.RequestContext, input UpdateAccountSettingsInput) (UpdateAccountSettingsResult, error) {
	if !validateMTBCStatus(input.Status) {
		return UpdateAccountSettingsResult{}, ErrInvalidArgument
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return UpdateAccountSettingsResult{}, err
	}

	if err := store.SetMinimumThroughputBillingCommitmentStatus(input.Status); err != nil {
		return UpdateAccountSettingsResult{}, s.mapStoreError(err)
	}

	return UpdateAccountSettingsResult{Status: input.Status}, nil
}

// describeAccountSettingsCore reports the persisted account settings: the
// commitment status Update last answered, DISABLED until one is set.
func (s *KinesisService) describeAccountSettingsCore(reqCtx *request.RequestContext) (UpdateAccountSettingsResult, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return UpdateAccountSettingsResult{}, err
	}

	status, err := store.MinimumThroughputBillingCommitmentStatus()
	if err != nil {
		return UpdateAccountSettingsResult{}, s.mapStoreError(err)
	}

	return UpdateAccountSettingsResult{Status: status}, nil
}

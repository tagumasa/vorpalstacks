package kinesis

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/kinesis"
	kinesisconnect "vorpalstacks/internal/pb/aws/kinesis/kinesisconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// AdminHandler provides Kinesis service administration functionality.
// It implements the KinesisServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	kinesisconnect.UnimplementedKinesisServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ kinesisconnect.KinesisServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Kinesis AdminHandler.
// It initialises the handler with the provided storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getKinesisStoreByRegion(region string) (*kinesisstore.KinesisStore, error) {
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	tstore, ok := regionStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("storage does not support transactions")
	}
	return kinesisstore.NewKinesisStore(tstore, h.accountId, region), nil
}

// ListStreams lists Kinesis streams.
// It returns all stream names in the region.
func (h *AdminHandler) ListStreams(ctx context.Context, req *connect.Request[pb.ListStreamsInput]) (*connect.Response[pb.ListStreamsOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getKinesisStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	streams, err := store.ListStreams()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	streamNames := make([]string, len(streams))
	for i, s := range streams {
		streamNames[i] = s.StreamName
	}

	return connect.NewResponse(&pb.ListStreamsOutput{
		Streamnames:    streamNames,
		Hasmorestreams: false,
	}), nil
}

// DescribeStream retrieves the detailed description of a Kinesis stream.
// It returns the stream's current status, shards, encryption settings, and other metadata.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeStream(ctx context.Context, req *connect.Request[pb.DescribeStreamInput]) (*connect.Response[pb.DescribeStreamOutput], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.getKinesisStoreByRegion(region)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	streamName := req.Msg.Streamname
	if streamName == "" && req.Msg.Streamarn != "" {
		stream, err := store.GetStreamByARN(req.Msg.Streamarn)
		if err != nil {
			if err == kinesisstore.ErrStreamNotFound {
				return nil, connect.NewError(connect.CodeNotFound, err)
			}
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		streamName = stream.StreamName
	}

	if streamName == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("StreamName or StreamARN is required"))
	}

	stream, err := store.GetStream(streamName)
	if err != nil {
		if err == kinesisstore.ErrStreamNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	shards, err := store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DescribeStreamOutput{
		Streamdescription: toPbStreamDescription(stream, shards),
	}), nil
}

// toPbStreamDescription converts an internal Stream object to a protobuf StreamDescription.
// It populates all relevant fields including stream status, shards, and encryption details.
func toPbStreamDescription(stream *kinesisstore.Stream, shards []*kinesisstore.Shard) *pb.StreamDescription {
	sd := &pb.StreamDescription{
		Streamname:              stream.StreamName,
		Streamarn:               stream.StreamARN,
		Streamstatus:            toPbStreamStatus(stream.StreamStatus),
		Retentionperiodhours:    stream.RetentionPeriodHours,
		Streamcreationtimestamp: stream.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Hasmoreshards:           false,
	}

	if stream.StreamModeDetails != nil {
		sd.Streammodedetails = &pb.StreamModeDetails{
			Streammode: toPbStreamMode(stream.StreamModeDetails.StreamMode),
		}
	}

	sd.Shards = make([]*pb.Shard, len(shards))
	for i, shard := range shards {
		sd.Shards[i] = toPbShard(shard)
	}

	if len(stream.EnhancedMonitoring) > 0 {
		sd.Enhancedmonitoring = make([]*pb.EnhancedMetrics, len(stream.EnhancedMonitoring))
		for i, em := range stream.EnhancedMonitoring {
			sd.Enhancedmonitoring[i] = &pb.EnhancedMetrics{
				Shardlevelmetrics: toPbMetricsNames(em.ShardLevelMetrics),
			}
		}
	}

	if stream.EncryptionType != "" {
		sd.Encryptiontype = toPbEncryptionType(stream.EncryptionType)
		sd.Keyid = stream.KeyID
	}

	return sd
}

// toPbShard converts an internal Shard object to a protobuf Shard.
// It includes the shard ID, parent shard ID, hash key range, and sequence number range.
func toPbShard(shard *kinesisstore.Shard) *pb.Shard {
	s := &pb.Shard{
		Shardid:       shard.ShardID,
		Parentshardid: shard.ParentShardID,
	}

	if shard.HashKeyRange != nil {
		s.Hashkeyrange = &pb.HashKeyRange{
			Startinghashkey: shard.HashKeyRange.StartingHashKey,
			Endinghashkey:   shard.HashKeyRange.EndingHashKey,
		}
	}

	if shard.SequenceNumberRange != nil {
		s.Sequencenumberrange = &pb.SequenceNumberRange{
			Startingsequencenumber: shard.SequenceNumberRange.StartingSequenceNumber,
			Endingsequencenumber:   shard.SequenceNumberRange.EndingSequenceNumber,
		}
	}

	if shard.AdjacentParentShardID != "" {
		s.Adjacentparentshardid = shard.AdjacentParentShardID
	}

	return s
}

// toPbStreamStatus converts an internal StreamStatus to a protobuf StreamStatus.
func toPbStreamStatus(status kinesisstore.StreamStatus) pb.StreamStatus {
	switch status {
	case kinesisstore.StreamStatusCreating:
		return pb.StreamStatus_STREAM_STATUS_CREATING
	case kinesisstore.StreamStatusActive:
		return pb.StreamStatus_STREAM_STATUS_ACTIVE
	case kinesisstore.StreamStatusDeleting:
		return pb.StreamStatus_STREAM_STATUS_DELETING
	case kinesisstore.StreamStatusUpdating:
		return pb.StreamStatus_STREAM_STATUS_UPDATING
	default:
		return pb.StreamStatus_STREAM_STATUS_ACTIVE
	}
}

// toPbStreamMode converts an internal StreamMode to a protobuf StreamMode.
func toPbStreamMode(mode kinesisstore.StreamMode) pb.StreamMode {
	switch mode {
	case kinesisstore.StreamModeProvisioned:
		return pb.StreamMode_STREAM_MODE_PROVISIONED
	case kinesisstore.StreamModeOnDemand:
		return pb.StreamMode_STREAM_MODE_ON_DEMAND
	default:
		return pb.StreamMode_STREAM_MODE_PROVISIONED
	}
}

// toPbEncryptionType converts an encryption type string to a protobuf EncryptionType.
func toPbEncryptionType(encryptionType string) pb.EncryptionType {
	if encryptionType == "KMS" {
		return pb.EncryptionType_ENCRYPTION_TYPE_KMS
	}
	return pb.EncryptionType_ENCRYPTION_TYPE_NONE
}

// toPbMetricsNames converts a slice of metric strings to a slice of protobuf MetricsNames.
func toPbMetricsNames(metrics []string) []pb.MetricsName {
	result := make([]pb.MetricsName, 0, len(metrics))
	for _, m := range metrics {
		result = append(result, toPbMetricsName(m))
	}
	return result
}

// toPbMetricsName converts a single metric string to a protobuf MetricsName.
func toPbMetricsName(metric string) pb.MetricsName {
	switch metric {
	case "IncomingBytes":
		return pb.MetricsName_METRICS_NAME_INCOMING_BYTES
	case "IncomingRecords":
		return pb.MetricsName_METRICS_NAME_INCOMING_RECORDS
	case "OutgoingBytes":
		return pb.MetricsName_METRICS_NAME_OUTGOING_BYTES
	case "OutgoingRecords":
		return pb.MetricsName_METRICS_NAME_OUTGOING_RECORDS
	case "WriteProvisionedThroughputExceeded":
		return pb.MetricsName_METRICS_NAME_WRITE_PROVISIONED_THROUGHPUT_EXCEEDED
	case "ReadProvisionedThroughputExceeded":
		return pb.MetricsName_METRICS_NAME_READ_PROVISIONED_THROUGHPUT_EXCEEDED
	case "IteratorAgeMilliseconds":
		return pb.MetricsName_METRICS_NAME_ITERATOR_AGE_MILLISECONDS
	default:
		return pb.MetricsName_METRICS_NAME_INCOMING_BYTES
	}
}

// AddTagsToStream adds tags to a Kinesis stream.
// Tags are key-value pairs that help you to group resources for billing purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AddTagsToStream(ctx context.Context, req *connect.Request[pb.AddTagsToStreamInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateStream creates a new Kinesis stream with the specified name and shard count.
// The stream must be explicitly created before any data can be sent to it.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateStream(ctx context.Context, req *connect.Request[pb.CreateStreamInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DecreaseStreamRetentionPeriod decreases the retention period of a Kinesis stream.
// This determines how long stream data is kept before it is automatically deleted.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DecreaseStreamRetentionPeriod(ctx context.Context, req *connect.Request[pb.DecreaseStreamRetentionPeriodInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteResourcePolicy deletes the resource-based policy attached to a Kinesis stream.
// This removes the cross-account access permissions for the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteResourcePolicy(ctx context.Context, req *connect.Request[pb.DeleteResourcePolicyInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteStream deletes a Kinesis stream and all its data.
// The stream must be in the ACTIVE state before it can be deleted.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteStream(ctx context.Context, req *connect.Request[pb.DeleteStreamInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeregisterStreamConsumer removes a stream consumer from a Kinesis stream.
// This disconnects the consumer from receiving data from the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeregisterStreamConsumer(ctx context.Context, req *connect.Request[pb.DeregisterStreamConsumerInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeAccountSettings returns the current Kinesis account settings.
// This includes information about the account's stream limits and usage.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeAccountSettings(ctx context.Context, req *connect.Request[pb.DescribeAccountSettingsInput]) (*connect.Response[pb.DescribeAccountSettingsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeLimits returns the current account limits for Kinesis streams and shards.
// This includes the maximum number of streams and shards allowed per region.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeLimits(ctx context.Context, req *connect.Request[pb.DescribeLimitsInput]) (*connect.Response[pb.DescribeLimitsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeStreamConsumer returns the detailed description of a stream consumer.
// Consumers are applications that read data from a Kinesis stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeStreamConsumer(ctx context.Context, req *connect.Request[pb.DescribeStreamConsumerInput]) (*connect.Response[pb.DescribeStreamConsumerOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeStreamSummary returns a summary description of a Kinesis stream.
// This includes the stream name, ARN, status, and other high-level information.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeStreamSummary(ctx context.Context, req *connect.Request[pb.DescribeStreamSummaryInput]) (*connect.Response[pb.DescribeStreamSummaryOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableEnhancedMonitoring disables enhanced monitoring for a Kinesis stream.
// This stops the collection of shard-level metrics for the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableEnhancedMonitoring(ctx context.Context, req *connect.Request[pb.DisableEnhancedMonitoringInput]) (*connect.Response[pb.EnhancedMonitoringOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableEnhancedMonitoring enables enhanced monitoring for a Kinesis stream.
// This enables the collection of shard-level metrics such as incoming bytes and records.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableEnhancedMonitoring(ctx context.Context, req *connect.Request[pb.EnableEnhancedMonitoringInput]) (*connect.Response[pb.EnhancedMonitoringOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetRecords retrieves records from a Kinesis stream shard.
// This is the primary operation for reading data from a Kinesis stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRecords(ctx context.Context, req *connect.Request[pb.GetRecordsInput]) (*connect.Response[pb.GetRecordsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetResourcePolicy retrieves the resource-based policy attached to a Kinesis stream.
// This policy controls cross-account access to the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResourcePolicy(ctx context.Context, req *connect.Request[pb.GetResourcePolicyInput]) (*connect.Response[pb.GetResourcePolicyOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetShardIterator gets a shard iterator for reading records from a Kinesis stream shard.
// The iterator is used in subsequent GetRecords calls to read data.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetShardIterator(ctx context.Context, req *connect.Request[pb.GetShardIteratorInput]) (*connect.Response[pb.GetShardIteratorOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// IncreaseStreamRetentionPeriod increases the retention period of a Kinesis stream.
// This extends how long stream data is kept before it is automatically deleted.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) IncreaseStreamRetentionPeriod(ctx context.Context, req *connect.Request[pb.IncreaseStreamRetentionPeriodInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListShards lists the shards of a Kinesis stream.
// This returns information about each shard including its ID, hash range, and parent shard.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListShards(ctx context.Context, req *connect.Request[pb.ListShardsInput]) (*connect.Response[pb.ListShardsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListStreamConsumers lists the registered consumers of a Kinesis stream.
// Consumers are applications that read data from the stream using enhanced fan-out.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListStreamConsumers(ctx context.Context, req *connect.Request[pb.ListStreamConsumersInput]) (*connect.Response[pb.ListStreamConsumersOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForResource lists tags for a Kinesis stream or consumer.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceInput]) (*connect.Response[pb.ListTagsForResourceOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListTagsForStream lists the tags attached to a Kinesis stream.
// Tags help you to group and categorize resources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForStream(ctx context.Context, req *connect.Request[pb.ListTagsForStreamInput]) (*connect.Response[pb.ListTagsForStreamOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// MergeShards merges two adjacent shards of a Kinesis stream into a single shard.
// This reduces the number of shards while maintaining data ordering.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) MergeShards(ctx context.Context, req *connect.Request[pb.MergeShardsInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutRecord writes a single data record into a Kinesis stream.
// Each record has a partition key that determines which shard the data goes to.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutRecord(ctx context.Context, req *connect.Request[pb.PutRecordInput]) (*connect.Response[pb.PutRecordOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutRecords writes multiple data records into a Kinesis stream in a single call.
// This is more efficient than calling PutRecord multiple times.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutRecords(ctx context.Context, req *connect.Request[pb.PutRecordsInput]) (*connect.Response[pb.PutRecordsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutResourcePolicy attaches a resource-based policy to a Kinesis stream.
// This policy controls which accounts can access the stream and what actions they can perform.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutResourcePolicy(ctx context.Context, req *connect.Request[pb.PutResourcePolicyInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RegisterStreamConsumer registers a consumer with a Kinesis stream.
// Consumers can use enhanced fan-out to receive data from the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RegisterStreamConsumer(ctx context.Context, req *connect.Request[pb.RegisterStreamConsumerInput]) (*connect.Response[pb.RegisterStreamConsumerOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RemoveTagsFromStream removes tags from a Kinesis stream.
// This removes the specified tag keys from the stream's tag set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveTagsFromStream(ctx context.Context, req *connect.Request[pb.RemoveTagsFromStreamInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SplitShard splits a shard into two new shards in a Kinesis stream.
// This is used to increase the capacity of a stream by dividing a shard's hash range.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SplitShard(ctx context.Context, req *connect.Request[pb.SplitShardInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// StartStreamEncryption enables server-side encryption for a Kinesis stream.
// This encrypts all data stored in the stream using the specified KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) StartStreamEncryption(ctx context.Context, req *connect.Request[pb.StartStreamEncryptionInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// StopStreamEncryption disables server-side encryption for a Kinesis stream.
// This removes encryption from all data stored in the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) StopStreamEncryption(ctx context.Context, req *connect.Request[pb.StopStreamEncryptionInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// SubscribeToShard subscribes a consumer to a Kinesis stream for receiving data.
// This provides enhanced fan-out with lower latency than regular consumers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) SubscribeToShard(ctx context.Context, req *connect.Request[pb.SubscribeToShardInput]) (*connect.Response[pb.SubscribeToShardOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagResource adds tags to a Kinesis stream or consumer.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagResource removes tags from a Kinesis stream or consumer.
// This removes the specified tag keys from the resource's tag set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateAccountSettings updates the Kinesis account settings.
// This allows you to configure on-demand stream limits and other account preferences.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAccountSettings(ctx context.Context, req *connect.Request[pb.UpdateAccountSettingsInput]) (*connect.Response[pb.UpdateAccountSettingsOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateMaxRecordSize updates the maximum record size for a Kinesis stream.
// This determines the maximum size of each record that can be written to the stream.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateMaxRecordSize(ctx context.Context, req *connect.Request[pb.UpdateMaxRecordSizeInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateShardCount updates the number of shards in a Kinesis stream.
// This scales the stream capacity up or down by splitting or merging shards.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateShardCount(ctx context.Context, req *connect.Request[pb.UpdateShardCountInput]) (*connect.Response[pb.UpdateShardCountOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateStreamMode updates the stream mode of a Kinesis stream.
// This switches between provisioned throughput and on-demand capacity modes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateStreamMode(ctx context.Context, req *connect.Request[pb.UpdateStreamModeInput]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateStreamWarmThroughput updates the warm throughput for a Kinesis stream.
// This configures the read and write throughput for on-demand streams.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateStreamWarmThroughput(ctx context.Context, req *connect.Request[pb.UpdateStreamWarmThroughputInput]) (*connect.Response[pb.UpdateStreamWarmThroughputOutput], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

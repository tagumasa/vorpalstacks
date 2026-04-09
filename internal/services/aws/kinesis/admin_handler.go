package kinesis

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/kinesis"
	kinesisconnect "vorpalstacks/internal/pb/aws/kinesis/kinesisconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// AdminHandler implements the Kinesis admin console gRPC-Web handler.
type AdminHandler struct {
	kinesisconnect.UnimplementedKinesisServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ kinesisconnect.KinesisServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Kinesis admin handler with the given storage manager and account ID.
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

// ListStreams returns a list of Kinesis stream names via the admin console gRPC-Web interface.
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

// DescribeStream returns detailed information about a Kinesis stream via the admin console gRPC-Web interface.
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

func toPbEncryptionType(encryptionType string) pb.EncryptionType {
	if encryptionType == "KMS" {
		return pb.EncryptionType_ENCRYPTION_TYPE_KMS
	}
	return pb.EncryptionType_ENCRYPTION_TYPE_NONE
}

func toPbMetricsNames(metrics []string) []pb.MetricsName {
	result := make([]pb.MetricsName, 0, len(metrics))
	for _, m := range metrics {
		result = append(result, toPbMetricsName(m))
	}
	return result
}

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

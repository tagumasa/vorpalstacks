package kinesis

import (
	"net/http"

	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/kinesis"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// This file is the sole exception to the store-import prohibition for admin
// handler files: it is the only admin handler file that imports the store
// package. It contains the store-access helper and pure proto conversion
// helpers (toPb* functions) that translate store types to proto types for
// response marshalling.

// getStores returns the KinesisStore for the region specified in the
// request headers.
func (h *AdminHandler) getStores(headers http.Header) (*kinesisstore.KinesisStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.getStoreForRegion(region)
}

// toPbStreamDescription converts store Stream and Shard objects to a proto
// StreamDescription for the admin console DescribeStream response.
func toPbStreamDescription(stream *kinesisstore.Stream, shards []*kinesisstore.Shard) *pb.StreamDescription {
	sd := &pb.StreamDescription{
		Streamname:              stream.StreamName,
		Streamarn:               stream.StreamARN,
		Streamstatus:            toPbStreamStatus(stream.StreamStatus),
		Retentionperiodhours:    stream.RetentionPeriodHours,
		Streamcreationtimestamp: stream.CreatedAt.Format(timeutils.ISO8601UTCFormat),
		Hasmoreshards:           proto.Bool(false),
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

// toPbStreamSummary converts a store Stream to a proto StreamSummary for
// the admin console ListStreams response.
func toPbStreamSummary(s *kinesisstore.Stream) *pb.StreamSummary {
	summary := &pb.StreamSummary{
		Streamname:              s.StreamName,
		Streamarn:               s.StreamARN,
		Streamstatus:            toPbStreamStatus(s.StreamStatus),
		Streamcreationtimestamp: s.CreatedAt.Format(timeutils.ISO8601UTCFormat),
	}
	if s.StreamModeDetails != nil {
		summary.Streammodedetails = &pb.StreamModeDetails{
			Streammode: toPbStreamMode(s.StreamModeDetails.StreamMode),
		}
	}
	return summary
}

// toPbShard converts a store Shard to a proto Shard.
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

// toPbStreamStatus converts a store StreamStatus to a proto StreamStatus.
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

// toPbStreamMode converts a store StreamMode to a proto StreamMode.
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

// toPbEncryptionType converts an encryption type string to a proto EncryptionType.
func toPbEncryptionType(encryptionType string) pb.EncryptionType {
	if encryptionType == "KMS" {
		return pb.EncryptionType_ENCRYPTION_TYPE_KMS
	}
	return pb.EncryptionType_ENCRYPTION_TYPE_NONE
}

// toPbMetricsNames converts a slice of metric name strings to proto MetricsName values.
func toPbMetricsNames(metrics []string) []pb.MetricsName {
	result := make([]pb.MetricsName, 0, len(metrics))
	for _, m := range metrics {
		result = append(result, toPbMetricsName(m))
	}
	return result
}

// toPbMetricsName converts a single metric name string to a proto MetricsName.
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

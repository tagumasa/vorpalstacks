package kinesis

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_kinesis"
)

// StreamMode conversion

// StreamModeToProto converts a StreamMode to its protobuf representation.
func StreamModeToProto(m StreamMode) pb.StreamMode {
	switch m {
	case StreamModeProvisioned:
		return pb.StreamMode_STREAM_MODE_PROVISIONED
	case StreamModeOnDemand:
		return pb.StreamMode_STREAM_MODE_ON_DEMAND
	default:
		return pb.StreamMode_STREAM_MODE_UNSPECIFIED
	}
}

// ProtoToStreamMode converts a protobuf StreamMode to its internal representation.
func ProtoToStreamMode(p pb.StreamMode) StreamMode {
	switch p {
	case pb.StreamMode_STREAM_MODE_PROVISIONED:
		return StreamModeProvisioned
	case pb.StreamMode_STREAM_MODE_ON_DEMAND:
		return StreamModeOnDemand
	default:
		return StreamModeProvisioned
	}
}

// StreamStatus conversion

// StreamStatusToProto converts a StreamStatus to its protobuf representation.
func StreamStatusToProto(s StreamStatus) string {
	return string(s)
}

// ProtoToStreamStatus converts a protobuf StreamStatus to its internal representation.
func ProtoToStreamStatus(p string) StreamStatus {
	return StreamStatus(p)
}

// HashKeyRange conversion

// HashKeyRangeToProto converts a HashKeyRange to its protobuf representation.
func HashKeyRangeToProto(h *HashKeyRange) *pb.HashKeyRange {
	if h == nil {
		return nil
	}
	return &pb.HashKeyRange{
		StartingHashKey: h.StartingHashKey,
		EndingHashKey:   h.EndingHashKey,
	}
}

// ProtoToHashKeyRange converts a protobuf HashKeyRange to its internal representation.
func ProtoToHashKeyRange(p *pb.HashKeyRange) *HashKeyRange {
	if p == nil {
		return nil
	}
	return &HashKeyRange{
		StartingHashKey: p.StartingHashKey,
		EndingHashKey:   p.EndingHashKey,
	}
}

// SequenceNumberRange conversion

// SequenceNumberRangeToProto converts a SequenceNumberRange to its protobuf representation.
func SequenceNumberRangeToProto(s *SequenceNumberRange) *pb.SequenceNumberRange {
	if s == nil {
		return nil
	}
	return &pb.SequenceNumberRange{
		StartingSequenceNumber: s.StartingSequenceNumber,
		EndingSequenceNumber:   s.EndingSequenceNumber,
	}
}

// ProtoToSequenceNumberRange converts a protobuf SequenceNumberRange to its internal representation.
func ProtoToSequenceNumberRange(p *pb.SequenceNumberRange) *SequenceNumberRange {
	if p == nil {
		return nil
	}
	return &SequenceNumberRange{
		StartingSequenceNumber: p.StartingSequenceNumber,
		EndingSequenceNumber:   p.EndingSequenceNumber,
	}
}

// StreamModeDetails conversion

// StreamModeDetailsToProto converts a StreamModeDetails to its protobuf representation.
func StreamModeDetailsToProto(s *StreamModeDetails) *pb.StreamModeDetails {
	if s == nil {
		return nil
	}
	return &pb.StreamModeDetails{
		StreamMode: StreamModeToProto(s.StreamMode),
	}
}

// ProtoToStreamModeDetails converts a protobuf StreamModeDetails to its internal representation.
func ProtoToStreamModeDetails(p *pb.StreamModeDetails) *StreamModeDetails {
	if p == nil {
		return nil
	}
	return &StreamModeDetails{
		StreamMode: ProtoToStreamMode(p.StreamMode),
	}
}

// EnhancedMonitoringToProto converts EnhancedMonitoring to its protobuf representation.
func EnhancedMonitoringToProto(e []EnhancedMonitoring) []*pb.EnhancedMonitoring {
	if e == nil {
		return nil
	}
	result := make([]*pb.EnhancedMonitoring, len(e))
	for i, em := range e {
		result[i] = &pb.EnhancedMonitoring{
			ShardLevelMetrics: em.ShardLevelMetrics,
		}
	}
	return result
}

// ProtoToEnhancedMonitoring converts protobuf EnhancedMonitoring to its internal representation.
func ProtoToEnhancedMonitoring(p []*pb.EnhancedMonitoring) []EnhancedMonitoring {
	if p == nil {
		return nil
	}
	result := make([]EnhancedMonitoring, len(p))
	for i, em := range p {
		result[i] = EnhancedMonitoring{
			ShardLevelMetrics: em.ShardLevelMetrics,
		}
	}
	return result
}

// Stream conversion

// StreamToProto converts a Stream to its protobuf representation.
func StreamToProto(s *Stream) *pb.Stream {
	if s == nil {
		return nil
	}
	return &pb.Stream{
		StreamName:           s.StreamName,
		StreamArn:            s.StreamARN,
		StreamStatus:         StreamStatusToProto(s.StreamStatus),
		StreamModeDetails:    StreamModeDetailsToProto(s.StreamModeDetails),
		ShardCount:           s.ShardCount,
		RetentionPeriodHours: s.RetentionPeriodHours,
		EnhancedMonitoring:   EnhancedMonitoringToProto(s.EnhancedMonitoring),
		EncryptionType:       s.EncryptionType,
		KeyId:                s.KeyID,
		ConsumerCount:        s.ConsumerCount,
		MaxRecordSizeInKib:   s.MaxRecordSizeInKiB,
		CreatedAt:            timestamppb.New(s.CreatedAt),
		LastModifiedAt:       timestamppb.New(s.LastModifiedAt),
	}
}

// ProtoToStream converts a protobuf Stream to its internal representation.
func ProtoToStream(p *pb.Stream) *Stream {
	if p == nil {
		return nil
	}
	return &Stream{
		StreamName:           p.StreamName,
		StreamARN:            p.StreamArn,
		StreamStatus:         ProtoToStreamStatus(p.StreamStatus),
		StreamModeDetails:    ProtoToStreamModeDetails(p.StreamModeDetails),
		ShardCount:           p.ShardCount,
		RetentionPeriodHours: p.RetentionPeriodHours,
		EnhancedMonitoring:   ProtoToEnhancedMonitoring(p.EnhancedMonitoring),
		EncryptionType:       p.EncryptionType,
		KeyID:                p.KeyId,
		ConsumerCount:        p.ConsumerCount,
		MaxRecordSizeInKiB:   p.MaxRecordSizeInKib,
		CreatedAt:            p.CreatedAt.AsTime(),
		LastModifiedAt:       p.LastModifiedAt.AsTime(),
	}
}

// StreamSummary conversion

// StreamSummaryToProto converts a StreamSummary to its protobuf representation.
func StreamSummaryToProto(s *StreamSummary) *pb.StreamSummary {
	if s == nil {
		return nil
	}
	return &pb.StreamSummary{
		StreamName:        s.StreamName,
		StreamArn:         s.StreamARN,
		StreamStatus:      StreamStatusToProto(s.StreamStatus),
		StreamModeDetails: StreamModeDetailsToProto(s.StreamModeDetails),
		ConsumerCount:     s.ConsumerCount,
	}
}

// ProtoToStreamSummary converts a protobuf StreamSummary to its internal representation.
func ProtoToStreamSummary(p *pb.StreamSummary) *StreamSummary {
	if p == nil {
		return nil
	}
	return &StreamSummary{
		StreamName:        p.StreamName,
		StreamARN:         p.StreamArn,
		StreamStatus:      ProtoToStreamStatus(p.StreamStatus),
		StreamModeDetails: ProtoToStreamModeDetails(p.StreamModeDetails),
		ConsumerCount:     p.ConsumerCount,
	}
}

// Shard conversion

// ShardToProto converts a Shard to its protobuf representation.
func ShardToProto(s *Shard) *pb.Shard {
	if s == nil {
		return nil
	}
	return &pb.Shard{
		ShardId:               s.ShardID,
		StreamName:            s.StreamName,
		ParentShardId:         s.ParentShardID,
		AdjacentParentShardId: s.AdjacentParentShardID,
		HashKeyRange:          HashKeyRangeToProto(s.HashKeyRange),
		SequenceNumberRange:   SequenceNumberRangeToProto(s.SequenceNumberRange),
		LatestSequenceNumber:  s.LatestSequenceNumber,
		CreatedAt:             timestamppb.New(s.CreatedAt),
	}
}

// ProtoToShard converts a protobuf Shard to its internal representation.
func ProtoToShard(p *pb.Shard) *Shard {
	if p == nil {
		return nil
	}
	return &Shard{
		ShardID:               p.ShardId,
		StreamName:            p.StreamName,
		ParentShardID:         p.ParentShardId,
		AdjacentParentShardID: p.AdjacentParentShardId,
		HashKeyRange:          ProtoToHashKeyRange(p.HashKeyRange),
		SequenceNumberRange:   ProtoToSequenceNumberRange(p.SequenceNumberRange),
		LatestSequenceNumber:  p.LatestSequenceNumber,
		CreatedAt:             p.CreatedAt.AsTime(),
	}
}

// Record conversion

// RecordToProto converts a Record to its protobuf representation.
func RecordToProto(r *Record) *pb.Record {
	if r == nil {
		return nil
	}
	return &pb.Record{
		SequenceNumber:              r.SequenceNumber,
		ApproximateArrivalTimestamp: timestamppb.New(r.ApproximateArrivalTimestamp),
		Data:                        r.Data,
		PartitionKey:                r.PartitionKey,
		EncryptionType:              r.EncryptionType,
	}
}

// ProtoToRecord converts a protobuf Record to its internal representation.
func ProtoToRecord(p *pb.Record) *Record {
	if p == nil {
		return nil
	}
	var arrivalTime time.Time
	if p.ApproximateArrivalTimestamp != nil {
		arrivalTime = p.ApproximateArrivalTimestamp.AsTime()
	}
	return &Record{
		SequenceNumber:              p.SequenceNumber,
		ApproximateArrivalTimestamp: arrivalTime,
		Data:                        p.Data,
		PartitionKey:                p.PartitionKey,
		EncryptionType:              p.EncryptionType,
	}
}

// Consumer conversion

// ConsumerToProto converts a Consumer to its protobuf representation.
func ConsumerToProto(c *Consumer) *pb.Consumer {
	if c == nil {
		return nil
	}
	return &pb.Consumer{
		ConsumerName:              c.ConsumerName,
		ConsumerArn:               c.ConsumerARN,
		StreamArn:                 c.StreamARN,
		ConsumerStatus:            c.ConsumerStatus,
		ConsumerCreationTimestamp: timestamppb.New(c.ConsumerCreationTimestamp),
	}
}

// ProtoToConsumer converts a protobuf Consumer to its internal representation.
func ProtoToConsumer(p *pb.Consumer) *Consumer {
	if p == nil {
		return nil
	}
	return &Consumer{
		ConsumerName:              p.ConsumerName,
		ConsumerARN:               p.ConsumerArn,
		StreamARN:                 p.StreamArn,
		ConsumerStatus:            p.ConsumerStatus,
		ConsumerCreationTimestamp: p.ConsumerCreationTimestamp.AsTime(),
	}
}

// ConsumerSummary conversion

// ConsumerSummaryToProto converts a ConsumerSummary to its protobuf representation.
func ConsumerSummaryToProto(c *ConsumerSummary) *pb.ConsumerSummary {
	if c == nil {
		return nil
	}
	return &pb.ConsumerSummary{
		ConsumerName:              c.ConsumerName,
		ConsumerArn:               c.ConsumerARN,
		StreamArn:                 c.StreamARN,
		ConsumerStatus:            c.ConsumerStatus,
		ConsumerCreationTimestamp: timestamppb.New(c.ConsumerCreationTimestamp),
	}
}

// ProtoToConsumerSummary converts a protobuf ConsumerSummary to its internal representation.
func ProtoToConsumerSummary(p *pb.ConsumerSummary) *ConsumerSummary {
	if p == nil {
		return nil
	}
	return &ConsumerSummary{
		ConsumerName:              p.ConsumerName,
		ConsumerARN:               p.ConsumerArn,
		StreamARN:                 p.StreamArn,
		ConsumerStatus:            p.ConsumerStatus,
		ConsumerCreationTimestamp: p.ConsumerCreationTimestamp.AsTime(),
	}
}

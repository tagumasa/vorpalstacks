package cloudwatchlogs

import (
	"time"

	pb "vorpalstacks/internal/pb/storage/storage_cloudwatchlogs"
)

// LogGroupToProto converts a LogGroup to its protobuf representation.
func LogGroupToProto(lg *LogGroup) *pb.LogGroup {
	if lg == nil {
		return nil
	}
	return &pb.LogGroup{
		Name:              lg.Name,
		Arn:               lg.ARN,
		Region:            lg.Region,
		AccountId:         lg.AccountID,
		CreatedAt:         lg.CreatedAt.UnixMilli(),
		RetentionInDays:   lg.RetentionInDays,
		MetricFilterCount: lg.MetricFilterCount,
		StoredBytes:       lg.StoredBytes,
		Tags:              lg.Tags,
	}
}

// ProtoToLogGroup converts a protobuf LogGroup to its internal representation.
func ProtoToLogGroup(p *pb.LogGroup) *LogGroup {
	if p == nil {
		return nil
	}
	return &LogGroup{
		Name:              p.Name,
		ARN:               p.Arn,
		Region:            p.Region,
		AccountID:         p.AccountId,
		CreatedAt:         time.UnixMilli(p.CreatedAt),
		RetentionInDays:   p.RetentionInDays,
		MetricFilterCount: p.MetricFilterCount,
		StoredBytes:       p.StoredBytes,
		Tags:              p.Tags,
	}
}

// LogStreamToProto converts a LogStream to its protobuf representation.
func LogStreamToProto(ls *LogStream) *pb.LogStream {
	if ls == nil {
		return nil
	}
	return &pb.LogStream{
		Name:                ls.Name,
		LogGroupName:        ls.LogGroupName,
		Arn:                 ls.ARN,
		CreatedAt:           ls.CreatedAt.UnixMilli(),
		FirstEventTs:        ls.FirstEventTs,
		LastEventTs:         ls.LastEventTs,
		LastIngestionTs:     ls.LastIngestionTs,
		UploadSequenceToken: ls.UploadSequenceToken,
	}
}

// ProtoToLogStream converts a protobuf LogStream to its internal representation.
func ProtoToLogStream(p *pb.LogStream) *LogStream {
	if p == nil {
		return nil
	}
	return &LogStream{
		Name:                p.Name,
		LogGroupName:        p.LogGroupName,
		ARN:                 p.Arn,
		CreatedAt:           time.UnixMilli(p.CreatedAt),
		FirstEventTs:        p.FirstEventTs,
		LastEventTs:         p.LastEventTs,
		LastIngestionTs:     p.LastIngestionTs,
		UploadSequenceToken: p.UploadSequenceToken,
	}
}

// ChunkMetaToProto converts a ChunkMeta to its protobuf representation.
func ChunkMetaToProto(cm *ChunkMeta) *pb.ChunkMeta {
	if cm == nil {
		return nil
	}
	return &pb.ChunkMeta{
		ChunkId:    cm.ChunkID,
		LogGroup:   cm.LogGroup,
		LogStream:  cm.LogStream,
		MinTs:      cm.MinTs,
		MaxTs:      cm.MaxTs,
		EntryCount: int32(cm.EntryCount),
		ChunkPath:  cm.ChunkPath,
	}
}

// ProtoToChunkMeta converts a protobuf ChunkMeta to its internal representation.
func ProtoToChunkMeta(p *pb.ChunkMeta) *ChunkMeta {
	if p == nil {
		return nil
	}
	return &ChunkMeta{
		ChunkID:    p.ChunkId,
		LogGroup:   p.LogGroup,
		LogStream:  p.LogStream,
		MinTs:      p.MinTs,
		MaxTs:      p.MaxTs,
		EntryCount: int(p.EntryCount),
		ChunkPath:  p.ChunkPath,
	}
}

// MetricFilterToProto converts a MetricFilter to its protobuf representation.
func MetricFilterToProto(mf *MetricFilter) *pb.MetricFilter {
	if mf == nil {
		return nil
	}
	transformations := make([]*pb.MetricTransformation, len(mf.MetricTransformations))
	for i, t := range mf.MetricTransformations {
		transformations[i] = &pb.MetricTransformation{
			MetricName:      t.MetricName,
			MetricNamespace: t.MetricNamespace,
			MetricValue:     t.MetricValue,
			DefaultValue:    t.DefaultValue,
			DefaultValueSet: t.DefaultValueSet,
		}
	}
	return &pb.MetricFilter{
		Name:                  mf.Name,
		LogGroupName:          mf.LogGroupName,
		FilterPattern:         mf.FilterPattern,
		MetricTransformations: transformations,
		CreatedAt:             mf.CreatedAt.UnixMilli(),
	}
}

// ProtoToMetricFilter converts a protobuf MetricFilter to its internal representation.
func ProtoToMetricFilter(p *pb.MetricFilter) *MetricFilter {
	if p == nil {
		return nil
	}
	transformations := make([]MetricTransformation, len(p.MetricTransformations))
	for i, t := range p.MetricTransformations {
		transformations[i] = MetricTransformation{
			MetricName:      t.MetricName,
			MetricNamespace: t.MetricNamespace,
			MetricValue:     t.MetricValue,
			DefaultValue:    t.DefaultValue,
			DefaultValueSet: t.DefaultValueSet,
		}
	}
	return &MetricFilter{
		Name:                  p.Name,
		LogGroupName:          p.LogGroupName,
		FilterPattern:         p.FilterPattern,
		MetricTransformations: transformations,
		CreatedAt:             time.UnixMilli(p.CreatedAt),
	}
}

// SubscriptionFilterToProto converts a SubscriptionFilter to its protobuf representation.
func SubscriptionFilterToProto(sf *SubscriptionFilter) *pb.SubscriptionFilter {
	if sf == nil {
		return nil
	}
	return &pb.SubscriptionFilter{
		LogGroupName:   sf.LogGroupName,
		FilterName:     sf.FilterName,
		FilterPattern:  sf.FilterPattern,
		DestinationArn: sf.DestinationArn,
		RoleArn:        sf.RoleArn,
		Distribution:   sf.Distribution,
		CreatedAt:      sf.CreationTime.UnixMilli(),
	}
}

// ProtoToSubscriptionFilter converts a protobuf SubscriptionFilter to its internal representation.
func ProtoToSubscriptionFilter(p *pb.SubscriptionFilter) *SubscriptionFilter {
	if p == nil {
		return nil
	}
	return &SubscriptionFilter{
		LogGroupName:   p.LogGroupName,
		FilterName:     p.FilterName,
		FilterPattern:  p.FilterPattern,
		DestinationArn: p.DestinationArn,
		RoleArn:        p.RoleArn,
		Distribution:   p.Distribution,
		CreationTime:   time.UnixMilli(p.CreatedAt),
	}
}

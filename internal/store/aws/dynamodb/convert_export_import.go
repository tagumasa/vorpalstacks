package dynamodb

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "vorpalstacks/internal/pb/storage/storage_dynamodb"
)

// ExportDescriptionToProto converts an ExportDescription to its protobuf representation.
func ExportDescriptionToProto(e *ExportDescription) *pb.ExportDescription {
	if e == nil {
		return nil
	}
	return &pb.ExportDescription{
		ExportArn:         e.ExportArn,
		ExportStatus:      exportStatusToProto(e.ExportStatus),
		StartTime:         timestamppb.New(e.StartTime),
		EndTime:           timestamppb.New(e.EndTime),
		ManifestFilesSize: e.ManifestFilesSize,
		ItemCount:         e.ItemCount,
		BilledSizeBytes:   e.BilledSizeBytes,
		ExportTime:        timestamppb.New(e.ExportTime),
		TableArn:          e.TableArn,
		TableId:           e.TableId,
		ExportFormat:      exportFormatToProto(e.ExportFormat),
		S3BucketSource: &pb.S3BucketSource{
			S3Bucket:      e.S3Bucket,
			S3Prefix:      e.S3Prefix,
			S3BucketOwner: e.S3BucketOwner,
		},
		FailureCode:    e.FailureCode,
		FailureMessage: e.FailureMessage,
		ClientToken:    e.ClientToken,
		S3SseKmsKeyId:  e.S3SseKmsKeyId,
		ExportManifest: e.ExportManifest,
		ExportType:     exportTypeToProto(e.ExportType),
		S3SseAlgorithm: s3SseAlgorithmToProto(e.S3SseAlgorithm),
		ExportViewType: exportViewTypeToProto(e.ExportViewType),
		ExportFromTime: timestamppb.New(e.ExportFromTime),
		ExportToTime:   timestamppb.New(e.ExportToTime),
	}
}

// ProtoToExportDescription converts a protobuf ExportDescription to its internal representation.
func ProtoToExportDescription(p *pb.ExportDescription) *ExportDescription {
	if p == nil {
		return nil
	}
	var s3Bucket, s3Prefix, s3BucketOwner string
	if src := p.S3BucketSource; src != nil {
		s3Bucket, s3Prefix, s3BucketOwner = src.S3Bucket, src.S3Prefix, src.S3BucketOwner
	}
	return &ExportDescription{
		ExportArn:         p.ExportArn,
		ExportStatus:      protoToExportStatus(p.ExportStatus),
		StartTime:         p.StartTime.AsTime(),
		EndTime:           p.EndTime.AsTime(),
		ManifestFilesSize: p.ManifestFilesSize,
		ItemCount:         p.ItemCount,
		BilledSizeBytes:   p.BilledSizeBytes,
		ExportTime:        p.ExportTime.AsTime(),
		TableArn:          p.TableArn,
		TableId:           p.TableId,
		ExportFormat:      protoToExportFormat(p.ExportFormat),
		S3Bucket:          s3Bucket,
		S3Prefix:          s3Prefix,
		FailureCode:       p.FailureCode,
		FailureMessage:    p.FailureMessage,
		ClientToken:       p.ClientToken,
		S3BucketOwner:     s3BucketOwner,
		S3SseKmsKeyId:     p.S3SseKmsKeyId,
		ExportManifest:    p.ExportManifest,
		ExportType:        protoToExportType(p.ExportType),
		S3SseAlgorithm:    protoToS3SseAlgorithm(p.S3SseAlgorithm),
		ExportViewType:    protoToExportViewType(p.ExportViewType),
		ExportFromTime:    p.ExportFromTime.AsTime(),
		ExportToTime:      p.ExportToTime.AsTime(),
	}
}

// ImportTableDescriptionToProto converts an ImportTableDescription to its protobuf representation.
func ImportTableDescriptionToProto(i *ImportTableDescription) *pb.ImportTableDescription {
	if i == nil {
		return nil
	}
	return &pb.ImportTableDescription{
		ImportArn:            i.ImportArn,
		ImportStatus:         importStatusToProto(i.ImportStatus),
		TableArn:             i.TableArn,
		TableId:              i.TableId,
		StartTime:            timestamppb.New(i.StartTime),
		EndTime:              timestamppb.New(i.EndTime),
		ProcessedItemCount:   i.ProcessedItemCount,
		ProcessedSizeBytes:   i.ProcessedSizeBytes,
		ImportedItemCount:    i.ImportedItemCount,
		ErrorCount:           i.ErrorCount,
		InputFormat:          inputFormatToProto(i.InputFormat),
		S3BucketSource:       s3BucketSourceToProto(i.S3BucketSource),
		FailureCode:          i.FailureCode,
		FailureMessage:       i.FailureMessage,
		ClientToken:          i.ClientToken,
		InputCompressionType: inputCompressionTypeToProto(i.InputCompressionType),
	}
}

// ProtoToImportTableDescription converts a protobuf ImportTableDescription to its internal representation.
func ProtoToImportTableDescription(p *pb.ImportTableDescription) *ImportTableDescription {
	if p == nil {
		return nil
	}
	return &ImportTableDescription{
		ImportArn:            p.ImportArn,
		ImportStatus:         protoToImportStatus(p.ImportStatus),
		TableArn:             p.TableArn,
		TableId:              p.TableId,
		StartTime:            p.StartTime.AsTime(),
		EndTime:              p.EndTime.AsTime(),
		ProcessedItemCount:   p.ProcessedItemCount,
		ProcessedSizeBytes:   p.ProcessedSizeBytes,
		ImportedItemCount:    p.ImportedItemCount,
		ErrorCount:           p.ErrorCount,
		InputFormat:          protoToInputFormat(p.InputFormat),
		S3BucketSource:       protoToS3BucketSource(p.S3BucketSource),
		FailureCode:          p.FailureCode,
		FailureMessage:       p.FailureMessage,
		ClientToken:          p.ClientToken,
		InputCompressionType: protoToInputCompressionType(p.InputCompressionType),
	}
}

func s3BucketSourceToProto(s *S3BucketSource) *pb.S3BucketSource {
	if s == nil {
		return nil
	}
	return &pb.S3BucketSource{
		S3Bucket:      s.S3Bucket,
		S3Prefix:      s.S3Prefix,
		S3BucketOwner: s.S3BucketOwner,
	}
}

func protoToS3BucketSource(p *pb.S3BucketSource) *S3BucketSource {
	if p == nil {
		return nil
	}
	return &S3BucketSource{
		S3Bucket:      p.S3Bucket,
		S3Prefix:      p.S3Prefix,
		S3BucketOwner: p.S3BucketOwner,
	}
}

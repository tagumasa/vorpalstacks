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
		ExportStatus:      e.ExportStatus,
		StartTime:         timestamppb.New(e.StartTime),
		EndTime:           timestamppb.New(e.EndTime),
		ManifestFilesSize: e.ManifestFilesSize,
		ItemCount:         e.ItemCount,
		TableArn:          e.TableArn,
		TableId:           e.TableId,
		ExportFormat:      e.ExportFormat,
		S3Bucket:          e.S3Bucket,
		S3Prefix:          e.S3Prefix,
		FailureCode:       e.FailureCode,
		FailureMessage:    e.FailureMessage,
	}
}

// ProtoToExportDescription converts a protobuf ExportDescription to its internal representation.
func ProtoToExportDescription(p *pb.ExportDescription) *ExportDescription {
	if p == nil {
		return nil
	}
	return &ExportDescription{
		ExportArn:         p.ExportArn,
		ExportStatus:      p.ExportStatus,
		StartTime:         p.StartTime.AsTime(),
		EndTime:           p.EndTime.AsTime(),
		ManifestFilesSize: p.ManifestFilesSize,
		ItemCount:         p.ItemCount,
		TableArn:          p.TableArn,
		TableId:           p.TableId,
		ExportFormat:      p.ExportFormat,
		S3Bucket:          p.S3Bucket,
		S3Prefix:          p.S3Prefix,
		FailureCode:       p.FailureCode,
		FailureMessage:    p.FailureMessage,
	}
}

// ImportTableDescriptionToProto converts an ImportTableDescription to its protobuf representation.
func ImportTableDescriptionToProto(i *ImportTableDescription) *pb.ImportTableDescription {
	if i == nil {
		return nil
	}
	return &pb.ImportTableDescription{
		ImportArn:          i.ImportArn,
		ImportStatus:       i.ImportStatus,
		TableArn:           i.TableArn,
		TableId:            i.TableId,
		StartTime:          timestamppb.New(i.StartTime),
		EndTime:            timestamppb.New(i.EndTime),
		ProcessedItemCount: i.ProcessedItemCount,
		ProcessedSizeBytes: i.ProcessedSizeBytes,
		InputFormat:        i.InputFormat,
		S3BucketSource:     s3BucketSourceToProto(i.S3BucketSource),
		FailureCode:        i.FailureCode,
		FailureMessage:     i.FailureMessage,
	}
}

// ProtoToImportTableDescription converts a protobuf ImportTableDescription to its internal representation.
func ProtoToImportTableDescription(p *pb.ImportTableDescription) *ImportTableDescription {
	if p == nil {
		return nil
	}
	return &ImportTableDescription{
		ImportArn:          p.ImportArn,
		ImportStatus:       p.ImportStatus,
		TableArn:           p.TableArn,
		TableId:            p.TableId,
		StartTime:          p.StartTime.AsTime(),
		EndTime:            p.EndTime.AsTime(),
		ProcessedItemCount: p.ProcessedItemCount,
		ProcessedSizeBytes: p.ProcessedSizeBytes,
		InputFormat:        p.InputFormat,
		S3BucketSource:     protoToS3BucketSource(p.S3BucketSource),
		FailureCode:        p.FailureCode,
		FailureMessage:     p.FailureMessage,
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

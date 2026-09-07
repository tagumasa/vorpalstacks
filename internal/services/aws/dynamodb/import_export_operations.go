// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
)

// s3invoker returns the S3 invoker from the EventBus, or nil if unavailable.
func (s *DynamoDBService) s3invoker() invokers.S3Invoker {
	if s.bus == nil {
		return nil
	}
	return s.bus.S3Invoker()
}

// ExportTableToPointInTime exports a DynamoDB table to S3.
func (s *DynamoDBService) ExportTableToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	outcome, err := s.exportTableCore(ctx, reqCtx, exportTableInput{
		Parameters: req.Parameters,
	})
	if err != nil {
		return nil, err
	}
	export := outcome.Export
	exportTime := outcome.ExportTime

	result := map[string]interface{}{
		"ExportArn":      export.ExportArn,
		"ExportStatus":   export.ExportStatus,
		"StartTime":      export.StartTime.Unix(),
		"ExportTime":     exportTime.Unix(),
		"TableArn":       export.TableArn,
		"ExportFormat":   export.ExportFormat,
		"S3Bucket":       export.S3Bucket,
		"S3Prefix":       export.S3Prefix,
		"S3BucketOwner":  export.S3BucketOwner,
		"S3SseKmsKeyId":  export.S3SseKmsKeyId,
		"ExportManifest": export.ExportManifest,
		"ExportType":     export.ExportType,
	}
	if export.ItemCount > 0 {
		result["ItemCount"] = export.ItemCount
		result["BilledSizeBytes"] = export.BilledSizeBytes
	}
	if !export.EndTime.IsZero() {
		result["EndTime"] = export.EndTime.Unix()
	}
	if export.FailureCode != "" {
		result["FailureCode"] = export.FailureCode
	}
	if export.FailureMessage != "" {
		result["FailureMessage"] = export.FailureMessage
	}

	return map[string]interface{}{
		"ExportDescription": result,
	}, nil
}

// DescribeExport returns information about a table export.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_DescribeExport.html
func (s *DynamoDBService) DescribeExport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	export, err := s.describeExportCore(ctx, reqCtx, request.GetStringParam(req.Parameters, "ExportArn"))
	if err != nil {
		return nil, err
	}

	description := map[string]interface{}{
		"ExportArn":      export.ExportArn,
		"ExportStatus":   export.ExportStatus,
		"StartTime":      export.StartTime.Unix(),
		"TableArn":       export.TableArn,
		"ExportFormat":   export.ExportFormat,
		"S3Bucket":       export.S3Bucket,
		"S3Prefix":       export.S3Prefix,
		"S3BucketOwner":  export.S3BucketOwner,
		"S3SseKmsKeyId":  export.S3SseKmsKeyId,
		"ExportManifest": export.ExportManifest,
		"ExportType":     export.ExportType,
	}
	if !export.ExportTime.IsZero() {
		description["ExportTime"] = export.ExportTime.Unix()
	}
	if !export.EndTime.IsZero() {
		description["EndTime"] = export.EndTime.Unix()
	}
	if export.ItemCount > 0 {
		description["ItemCount"] = export.ItemCount
		description["BilledSizeBytes"] = export.BilledSizeBytes
	}
	if export.FailureCode != "" {
		description["FailureCode"] = export.FailureCode
		description["FailureMessage"] = export.FailureMessage
	}

	return map[string]interface{}{
		"ExportDescription": description,
	}, nil
}

// ListExports lists the exports for a table.
func (s *DynamoDBService) ListExports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	exports, nextToken, err := s.listExportsCore(ctx, reqCtx, listExportsInput{
		Parameters: req.Parameters,
	})
	if err != nil {
		return nil, err
	}

	exportSummaries := make([]map[string]interface{}, 0)
	for _, e := range exports {
		summary := map[string]interface{}{
			"ExportArn":    e.ExportArn,
			"ExportStatus": e.ExportStatus,
			"ExportType":   e.ExportType,
		}
		exportSummaries = append(exportSummaries, summary)
	}

	return pagination.BuildListResponse("ExportSummaries", exportSummaries, nextToken), nil
}

// ImportTable imports table data from S3.
func (s *DynamoDBService) ImportTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	outcome, err := s.importTableCore(ctx, reqCtx, importTableInput{
		Parameters: req.Parameters,
	})
	if err != nil {
		return nil, err
	}
	imp := outcome.Import
	s3Bucket := outcome.S3Bucket
	s3Prefix := outcome.S3Prefix
	s3BucketOwner := outcome.S3BucketOwner

	description := map[string]interface{}{
		"ImportArn":            imp.ImportArn,
		"ImportStatus":         imp.ImportStatus,
		"StartTime":            imp.StartTime.Unix(),
		"TableArn":             imp.TableArn,
		"InputFormat":          imp.InputFormat,
		"InputCompressionType": imp.InputCompressionType,
		"S3BucketSource": map[string]interface{}{
			"S3Bucket":      s3Bucket,
			"S3KeyPrefix":   s3Prefix,
			"S3BucketOwner": s3BucketOwner,
		},
	}
	if !imp.EndTime.IsZero() {
		description["EndTime"] = imp.EndTime.Unix()
	}
	if imp.ProcessedItemCount > 0 {
		description["ProcessedItemCount"] = imp.ProcessedItemCount
		description["ProcessedSizeBytes"] = imp.ProcessedSizeBytes
	}
	if imp.ImportedItemCount > 0 {
		description["ImportedItemCount"] = imp.ImportedItemCount
	}
	if imp.ErrorCount > 0 {
		description["ErrorCount"] = imp.ErrorCount
	}
	if imp.FailureCode != "" {
		description["FailureCode"] = imp.FailureCode
		description["FailureMessage"] = imp.FailureMessage
	}

	return map[string]interface{}{
		"ImportTableDescription": description,
	}, nil
}

// DescribeImport returns information about a table import.
func (s *DynamoDBService) DescribeImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	imp, err := s.describeImportCore(ctx, reqCtx, request.GetStringParam(req.Parameters, "ImportArn"))
	if err != nil {
		return nil, err
	}

	description := map[string]interface{}{
		"ImportArn":            imp.ImportArn,
		"ImportStatus":         imp.ImportStatus,
		"StartTime":            imp.StartTime.Unix(),
		"TableArn":             imp.TableArn,
		"ClientToken":          imp.ClientToken,
		"InputFormat":          imp.InputFormat,
		"InputCompressionType": imp.InputCompressionType,
		"ProcessedItemCount":   imp.ProcessedItemCount,
		"ImportedItemCount":    imp.ImportedItemCount,
		"ErrorCount":           imp.ErrorCount,
	}
	if imp.S3BucketSource != nil {
		description["S3BucketSource"] = map[string]interface{}{
			"S3Bucket":      imp.S3BucketSource.S3Bucket,
			"S3KeyPrefix":   imp.S3BucketSource.S3Prefix,
			"S3BucketOwner": imp.S3BucketSource.S3BucketOwner,
		}
	}
	if !imp.EndTime.IsZero() {
		description["EndTime"] = imp.EndTime.Unix()
	}
	if imp.FailureCode != "" {
		description["FailureCode"] = imp.FailureCode
		description["FailureMessage"] = imp.FailureMessage
	}

	return map[string]interface{}{
		"ImportTableDescription": description,
	}, nil
}

// ListImports lists the imports for a table.
func (s *DynamoDBService) ListImports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	imports, nextToken, err := s.listImportsCore(ctx, reqCtx, listImportsInput{
		Parameters: req.Parameters,
	})
	if err != nil {
		return nil, err
	}

	importSummaries := make([]map[string]interface{}, 0)
	for _, i := range imports {
		summary := map[string]interface{}{
			"ImportArn":    i.ImportArn,
			"ImportStatus": i.ImportStatus,
			"TableArn":     i.TableArn,
			"InputFormat":  i.InputFormat,
			"StartTime":    i.StartTime.Unix(),
		}
		if i.S3BucketSource != nil {
			summary["S3BucketSource"] = map[string]interface{}{
				"S3Bucket":      i.S3BucketSource.S3Bucket,
				"S3KeyPrefix":   i.S3BucketSource.S3Prefix,
				"S3BucketOwner": i.S3BucketSource.S3BucketOwner,
			}
		}
		if !i.EndTime.IsZero() {
			summary["EndTime"] = i.EndTime.Unix()
		}
		importSummaries = append(importSummaries, summary)
	}

	return pagination.BuildListResponse("ImportSummaryList", importSummaries, nextToken), nil
}

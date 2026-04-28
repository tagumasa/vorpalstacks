// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ExportTableToPointInTime exports a DynamoDB table to S3.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ExportTableToPointInTime.html
func (s *DynamoDBService) ExportTableToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableArn := request.GetStringParam(req.Parameters, "TableArn")
	if tableArn == "" {
		return nil, ErrInvalidParameter
	}

	s3Bucket := request.GetStringParam(req.Parameters, "S3Bucket")
	if s3Bucket == "" {
		return nil, ErrInvalidParameter
	}

	s3Prefix := request.GetStringParam(req.Parameters, "S3Prefix")

	if len(s3Bucket) < 3 || len(s3Bucket) > 63 {
		return nil, ErrInvalidParameter
	}

	tableName := svcarn.ParseTableARN(tableArn)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrTableNotFound
	}

	pitr, err := store.Tables().GetPointInTimeRecovery(tableName)
	if err != nil || pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
		return nil, ErrPITRNotEnabled
	}

	exportFormat := request.GetStringParam(req.Parameters, "ExportFormat")
	if exportFormat == "" {
		exportFormat = "DYNAMODB_JSON"
	}

	export, err := store.Exports().Create(tableArn, tableName, exportFormat)
	if err != nil {
		return nil, err
	}

	export.ExportStatus = "COMPLETED"
	export.S3Bucket = s3Bucket
	export.S3Prefix = s3Prefix
	export.EndTime = time.Now()
	export.ItemCount = table.ItemCount
	if err := store.Exports().Put(export); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ExportDescription": map[string]interface{}{
			"ExportArn":    export.ExportArn,
			"ExportStatus": export.ExportStatus,
			"StartTime":    export.StartTime.Unix(),
			"EndTime":      export.EndTime.Unix(),
			"TableArn":     export.TableArn,
			"ExportFormat": export.ExportFormat,
			"S3Bucket":     export.S3Bucket,
			"S3Prefix":     export.S3Prefix,
			"ItemCount":    export.ItemCount,
		},
	}, nil
}

// DescribeExport returns information about a table export.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_DescribeExport.html
func (s *DynamoDBService) DescribeExport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	exportArn := request.GetStringParam(req.Parameters, "ExportArn")
	if exportArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	export, err := store.Exports().Get(exportArn)
	if err != nil {
		return nil, ErrExportNotFound
	}

	return map[string]interface{}{
		"ExportDescription": map[string]interface{}{
			"ExportArn":    export.ExportArn,
			"ExportStatus": export.ExportStatus,
			"StartTime":    export.StartTime.Unix(),
			"EndTime":      export.EndTime.Unix(),
			"TableArn":     export.TableArn,
			"ExportFormat": export.ExportFormat,
			"ItemCount":    export.ItemCount,
		},
	}, nil
}

// ListExports lists the exports for a table.
func (s *DynamoDBService) ListExports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableArn := request.GetStringParam(req.Parameters, "TableArn")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exports, nextToken, err := store.Exports().List(tableArn, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	exportSummaries := make([]map[string]interface{}, 0)
	for _, e := range exports {
		exportSummaries = append(exportSummaries, map[string]interface{}{
			"ExportArn":    e.ExportArn,
			"ExportStatus": e.ExportStatus,
			"StartTime":    e.StartTime.Unix(),
			"TableArn":     e.TableArn,
		})
	}

	return pagination.BuildListResponse("ExportSummaries", exportSummaries, nextToken), nil
}

// ImportTable imports table data from S3.
func (s *DynamoDBService) ImportTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	inputFormat := request.GetStringParam(req.Parameters, "InputFormat")
	if inputFormat == "" {
		inputFormat = "DYNAMODB_JSON"
	}

	validFormats := map[string]bool{
		"DYNAMODB_JSON": true,
		"ION":           true,
		"CSV":           true,
	}
	if !validFormats[inputFormat] {
		return nil, ErrInvalidParameter
	}

	s3BucketSourceParam, ok := req.Parameters["S3BucketSource"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	s3Bucket, _ := s3BucketSourceParam["S3Bucket"].(string)
	if s3Bucket == "" {
		return nil, ErrInvalidParameter
	}

	if len(s3Bucket) < 3 || len(s3Bucket) > 63 {
		return nil, ErrInvalidParameter
	}

	s3Prefix, _ := s3BucketSourceParam["S3Prefix"].(string)
	s3BucketOwner, _ := s3BucketSourceParam["S3BucketOwner"].(string)

	tableCreationParams, ok := req.Parameters["TableCreationParameters"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	tableName, _ := tableCreationParams["TableName"].(string)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	if len(tableName) < 3 || len(tableName) > 255 {
		return nil, ErrInvalidParameter
	}

	if !tableNameRegex.MatchString(tableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if store.Tables().Exists(tableName) {
		return nil, ErrTableAlreadyExists
	}

	keySchema := parseKeySchema(tableCreationParams)
	if len(keySchema) == 0 {
		keySchema = []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}}
	}

	attrDefs := parseAttributeDefinitions(tableCreationParams)
	if len(attrDefs) == 0 {
		attrDefs = []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}}
	}

	billingMode := dbstore.BillingMode(request.GetStringParam(tableCreationParams, "BillingMode"))
	if billingMode == "" {
		billingMode = dbstore.BillingModePayPerRequest
	}

	provThroughput := parseProvisionedThroughput(tableCreationParams)
	if billingMode == dbstore.BillingModeProvisioned && provThroughput == nil {
		provThroughput = &dbstore.ProvisionedThroughput{
			ReadCapacityUnits:  5,
			WriteCapacityUnits: 5,
		}
	}

	table, err := store.Tables().Create(
		tableName,
		keySchema,
		attrDefs,
		billingMode,
		provThroughput,
		parseGlobalSecondaryIndexes(tableCreationParams),
		parseLocalSecondaryIndexes(tableCreationParams),
		nil,
		nil,
		false,
	)
	if err != nil {
		return nil, err
	}

	imp, err := store.Imports().Create(table.ARN, "")
	if err != nil {
		if delErr := store.Tables().Delete(tableName); delErr != nil {
			logs.Error("Failed to rollback table creation during import", logs.Err(delErr), logs.String("tableName", tableName))
		}
		return nil, err
	}

	imp.ImportStatus = "COMPLETED"
	imp.TableArn = table.ARN
	imp.InputFormat = inputFormat
	imp.S3BucketSource = &dbstore.S3BucketSource{
		S3Bucket:      s3Bucket,
		S3Prefix:      s3Prefix,
		S3BucketOwner: s3BucketOwner,
	}
	imp.EndTime = time.Now()
	if err := store.Imports().Put(imp); err != nil {
		if delErr := store.Tables().Delete(tableName); delErr != nil {
			logs.Error("Failed to rollback table creation during import Put", logs.Err(delErr), logs.String("tableName", tableName))
		}
		return nil, err
	}

	return map[string]interface{}{
		"ImportTableDescription": map[string]interface{}{
			"ImportArn":    imp.ImportArn,
			"ImportStatus": imp.ImportStatus,
			"StartTime":    imp.StartTime.Unix(),
			"EndTime":      imp.EndTime.Unix(),
			"TableArn":     imp.TableArn,
			"InputFormat":  imp.InputFormat,
			"S3BucketSource": map[string]interface{}{
				"S3Bucket":      s3Bucket,
				"S3Prefix":      s3Prefix,
				"S3BucketOwner": s3BucketOwner,
			},
		},
	}, nil
}

// DescribeImport returns information about a table import.
func (s *DynamoDBService) DescribeImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	importArn := request.GetStringParam(req.Parameters, "ImportArn")
	if importArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	imp, err := store.Imports().Get(importArn)
	if err != nil {
		return nil, ErrImportNotFound
	}

	return map[string]interface{}{
		"ImportTableDescription": map[string]interface{}{
			"ImportArn":    imp.ImportArn,
			"ImportStatus": imp.ImportStatus,
			"StartTime":    imp.StartTime.Unix(),
			"EndTime":      imp.EndTime.Unix(),
			"TableArn":     imp.TableArn,
		},
	}, nil
}

// ListImports lists the imports for a table.
func (s *DynamoDBService) ListImports(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableArn := request.GetStringParam(req.Parameters, "TableArn")
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	imports, nextToken, err := store.Imports().List(tableArn, nextToken, maxResults)
	if err != nil {
		return nil, err
	}

	importSummaries := make([]map[string]interface{}, 0)
	for _, i := range imports {
		importSummaries = append(importSummaries, map[string]interface{}{
			"ImportArn":    i.ImportArn,
			"ImportStatus": i.ImportStatus,
			"TableArn":     i.TableArn,
		})
	}

	return pagination.BuildListResponse("ImportSummaryList", importSummaries, nextToken), nil
}

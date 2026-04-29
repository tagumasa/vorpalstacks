// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// s3invoker returns the S3 invoker from the EventBus, or nil if unavailable.
func (s *DynamoDBService) s3invoker() eventbus.S3Invoker {
	if s.bus == nil {
		return nil
	}
	return s.bus.S3Invoker()
}

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

	region := reqCtx.GetRegion()
	itemCount := int64(0)
	billedSizeBytes := int64(0)

	s3 := s.s3invoker()
	if s3 != nil {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)

		err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
			itemCount++
			itemJSON := buildDynamoDBJSONItem(item)
			line, err := json.Marshal(itemJSON)
			if err != nil {
				return err
			}
			if _, err := writer.Write(line); err != nil {
				return err
			}
			if _, err := writer.Write([]byte("\n")); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			logs.Error("Failed to scan items for export",
				logs.Err(err),
				logs.String("tableName", tableName),
			)
		}

		if err := writer.Flush(); err != nil {
			logs.Error("Failed to flush export buffer",
				logs.Err(err),
				logs.String("tableName", tableName),
			)
		}

		billedSizeBytes = int64(buf.Len())
		objectKey := fmt.Sprintf("%s/AWSDynamoDB/%s/data/%s.json",
			s3Prefix, export.ExportArn, export.ExportArn,
		)
		if err := s3.PutObject(ctx, region, s3Bucket, objectKey, buf.Bytes()); err != nil {
			logs.Error("Failed to write export to S3",
				logs.Err(err),
				logs.String("bucket", s3Bucket),
				logs.String("key", objectKey),
			)
		}
	} else {
		itemCount = table.ItemCount
	}

	export.ExportStatus = "COMPLETED"
	export.S3Bucket = s3Bucket
	export.S3Prefix = s3Prefix
	export.EndTime = time.Now()
	export.ItemCount = itemCount
	export.BilledSizeBytes = billedSizeBytes
	if err := store.Exports().Put(export); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ExportDescription": map[string]interface{}{
			"ExportArn":       export.ExportArn,
			"ExportStatus":    export.ExportStatus,
			"StartTime":       export.StartTime.Unix(),
			"EndTime":         export.EndTime.Unix(),
			"TableArn":        export.TableArn,
			"ExportFormat":    export.ExportFormat,
			"S3Bucket":        export.S3Bucket,
			"S3Prefix":        export.S3Prefix,
			"ItemCount":       export.ItemCount,
			"BilledSizeBytes": export.BilledSizeBytes,
		},
	}, nil
}

// buildDynamoDBJSONItem converts a store Item to the DynamoDB JSON format
// used by import/export: {"Item": {"attr": {"S": "val"}, ...}}.
func buildDynamoDBJSONItem(item *dbstore.Item) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range item.Key {
		merged[k] = buildAttributeValueResponse(v)
	}
	for k, v := range item.Attributes {
		merged[k] = buildAttributeValueResponse(v)
	}
	return map[string]interface{}{
		"Item": merged,
	}
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
			"ExportArn":       export.ExportArn,
			"ExportStatus":    export.ExportStatus,
			"StartTime":       export.StartTime.Unix(),
			"EndTime":         export.EndTime.Unix(),
			"TableArn":        export.TableArn,
			"ExportFormat":    export.ExportFormat,
			"ItemCount":       export.ItemCount,
			"BilledSizeBytes": export.BilledSizeBytes,
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
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ImportTable.html
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

	importedCount := int64(0)
	billedSizeBytes := int64(0)

	s3 := s.s3invoker()
	if s3 != nil && inputFormat == "DYNAMODB_JSON" {
		region := reqCtx.GetRegion()
		keys, listErr := s3.ListObjects(ctx, region, s3Bucket, s3Prefix)
		if listErr != nil {
			logs.Error("Failed to list S3 objects for import",
				logs.Err(listErr),
				logs.String("bucket", s3Bucket),
				logs.String("prefix", s3Prefix),
			)
		} else {
			for _, key := range keys {
				data, getErr := s3.GetObject(ctx, region, s3Bucket, key)
				if getErr != nil {
					logs.Error("Failed to read S3 object for import",
						logs.Err(getErr),
						logs.String("bucket", s3Bucket),
						logs.String("key", key),
					)
					continue
				}
				billedSizeBytes += int64(len(data))
				count, parseErr := importDynamoDBJSONData(data, tableName, store)
				if parseErr != nil {
					logs.Error("Failed to parse DynamoDB JSON for import",
						logs.Err(parseErr),
						logs.String("key", key),
					)
					continue
				}
				importedCount += count
			}
		}
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
	imp.ProcessedItemCount = importedCount
	imp.ProcessedSizeBytes = billedSizeBytes
	if err := store.Imports().Put(imp); err != nil {
		if delErr := store.Tables().Delete(tableName); delErr != nil {
			logs.Error("Failed to rollback table creation during import Put", logs.Err(delErr), logs.String("tableName", tableName))
		}
		return nil, err
	}

	return map[string]interface{}{
		"ImportTableDescription": map[string]interface{}{
			"ImportArn":          imp.ImportArn,
			"ImportStatus":       imp.ImportStatus,
			"StartTime":          imp.StartTime.Unix(),
			"EndTime":            imp.EndTime.Unix(),
			"TableArn":           imp.TableArn,
			"InputFormat":        imp.InputFormat,
			"ProcessedItemCount": imp.ProcessedItemCount,
			"S3BucketSource": map[string]interface{}{
				"S3Bucket":      s3Bucket,
				"S3Prefix":      s3Prefix,
				"S3BucketOwner": s3BucketOwner,
			},
		},
	}, nil
}

// importDynamoDBJSONData parses DYNAMODB_JSON data (newline-delimited
// {"Item": {...}} objects) and writes each item to the store.
func importDynamoDBJSONData(data []byte, tableName string, store dbstore.DynamoDBStoreInterface) (int64, error) {
	var count int64
	scanner := bufio.NewScanner(bytes.NewReader(data))
	// Allow lines up to 10MB per item.
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var itemWrapper map[string]json.RawMessage
		if err := json.Unmarshal(line, &itemWrapper); err != nil {
			continue
		}

		itemRaw, ok := itemWrapper["Item"]
		if !ok {
			continue
		}

		var itemMap map[string]interface{}
		if err := json.Unmarshal(itemRaw, &itemMap); err != nil {
			continue
		}

		attrs := parseAttributeValueMap(itemMap)
		if len(attrs) == 0 {
			continue
		}

		key := make(map[string]*dbstore.AttributeValue)
		table, err := store.Tables().Get(tableName)
		if err != nil {
			continue
		}
		for _, ks := range table.KeySchema {
			if v, ok := attrs[ks.AttributeName]; ok {
				key[ks.AttributeName] = v
			}
		}
		if len(key) == 0 {
			continue
		}

		if _, err := store.Items().Put(tableName, key, attrs); err != nil {
			continue
		}
		count++
	}
	return count, scanner.Err()
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
			"ImportArn":          imp.ImportArn,
			"ImportStatus":       imp.ImportStatus,
			"StartTime":          imp.StartTime.Unix(),
			"EndTime":            imp.EndTime.Unix(),
			"TableArn":           imp.TableArn,
			"ProcessedItemCount": imp.ProcessedItemCount,
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

// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
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
		"ExportArn":    export.ExportArn,
		"ExportStatus": export.ExportStatus,
		"StartTime":    export.StartTime.Unix(),
		"ExportTime":   exportTime.Unix(),
		"TableArn":     export.TableArn,
		"ExportFormat": export.ExportFormat,
		"S3Bucket":     export.S3Bucket,
		"S3Prefix":     export.S3Prefix,
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
	export, err := s.describeExportCore(ctx, reqCtx, request.GetStringParam(req.Parameters, "ExportArn"))
	if err != nil {
		return nil, err
	}

	description := map[string]interface{}{
		"ExportArn":    export.ExportArn,
		"ExportStatus": export.ExportStatus,
		"StartTime":    export.StartTime.Unix(),
		"TableArn":     export.TableArn,
		"ExportFormat": export.ExportFormat,
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
		"ImportArn":    imp.ImportArn,
		"ImportStatus": imp.ImportStatus,
		"StartTime":    imp.StartTime.Unix(),
		"TableArn":     imp.TableArn,
		"InputFormat":  imp.InputFormat,
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
	if imp.FailureCode != "" {
		description["FailureCode"] = imp.FailureCode
		description["FailureMessage"] = imp.FailureMessage
	}

	return map[string]interface{}{
		"ImportTableDescription": description,
	}, nil
}

// importDynamoDBJSONData parses DYNAMODB_JSON data (newline-delimited
// {"Item": {...}} objects) and writes each item to the store with proper
// index entries and table metric updates.
func importDynamoDBJSONData(ctx context.Context, data []byte, tableName string, store dbstore.DynamoDBStoreInterface) (int64, error) {
	var count int64

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return 0, fmt.Errorf("get table %s for import: %w", tableName, err)
	}

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

		attrs, parseErr := parseAttributeValueMap(itemMap)
		if parseErr != nil {
			continue
		}
		if len(attrs) == 0 {
			continue
		}

		key := make(map[string]*dbstore.AttributeValue)
		for _, ks := range table.KeySchema {
			if v, ok := attrs[ks.AttributeName]; ok {
				key[ks.AttributeName] = v
			}
		}
		if len(key) == 0 {
			continue
		}

		itemSize := calculateItemSize(attrs)
		putErr := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.PutItem(tableName, key, attrs); err != nil {
				return err
			}
			newItem := &dbstore.Item{
				TableName:  tableName,
				Key:        key,
				Attributes: attrs,
			}
			if err := txn.PutIndexEntries(tableName, newItem); err != nil {
				return err
			}
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, itemSize); err != nil {
				return err
			}
			return nil
		})
		if putErr != nil {
			continue
		}
		count++
	}
	return count, scanner.Err()
}

// importCSVData parses CSV-formatted data and imports rows as DynamoDB items.
// Each column becomes a String attribute. If headerList is empty, the first
// row is treated as the header. delimiter defaults to comma.
func importCSVData(ctx context.Context, data []byte, tableName string, store dbstore.DynamoDBStoreInterface, delimiter string, headerList []string) (int64, error) {
	var count int64

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return 0, fmt.Errorf("get table %s for CSV import: %w", tableName, err)
	}

	if delimiter == "" {
		delimiter = ","
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = rune(delimiter[0])

	headers := headerList
	if len(headers) == 0 {
		headerRecord, err := reader.Read()
		if err != nil {
			return 0, fmt.Errorf("failed to read CSV header: %w", err)
		}
		headers = headerRecord
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		attrs := make(map[string]*dbstore.AttributeValue)
		for i, val := range record {
			if i >= len(headers) {
				break
			}
			attrs[headers[i]] = &dbstore.AttributeValue{S: &val}
		}

		key := make(map[string]*dbstore.AttributeValue)
		for _, ks := range table.KeySchema {
			if v, ok := attrs[ks.AttributeName]; ok {
				key[ks.AttributeName] = v
			}
		}
		if len(key) == 0 {
			continue
		}

		itemSize := calculateItemSize(attrs)
		putErr := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.PutItem(tableName, key, attrs); err != nil {
				return err
			}
			newItem := &dbstore.Item{
				TableName:  tableName,
				Key:        key,
				Attributes: attrs,
			}
			if err := txn.PutIndexEntries(tableName, newItem); err != nil {
				return err
			}
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, itemSize); err != nil {
				return err
			}
			return nil
		})
		if putErr != nil {
			continue
		}
		count++
	}
	return count, nil
}

// DescribeImport returns information about a table import.
func (s *DynamoDBService) DescribeImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	imp, err := s.describeImportCore(ctx, reqCtx, request.GetStringParam(req.Parameters, "ImportArn"))
	if err != nil {
		return nil, err
	}

	description := map[string]interface{}{
		"ImportArn":          imp.ImportArn,
		"ImportStatus":       imp.ImportStatus,
		"StartTime":          imp.StartTime.Unix(),
		"TableArn":           imp.TableArn,
		"ProcessedItemCount": imp.ProcessedItemCount,
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
		importSummaries = append(importSummaries, map[string]interface{}{
			"ImportArn":    i.ImportArn,
			"ImportStatus": i.ImportStatus,
			"TableArn":     i.TableArn,
		})
	}

	return pagination.BuildListResponse("ImportSummaryList", importSummaries, nextToken), nil
}

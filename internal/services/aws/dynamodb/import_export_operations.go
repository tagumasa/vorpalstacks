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
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
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
	if !validateTableArn(tableArn) {
		return nil, ErrInvalidParameter
	}

	s3Bucket := request.GetStringParam(req.Parameters, "S3Bucket")
	if !validateS3Bucket(s3Bucket) {
		return nil, ErrInvalidParameter
	}

	s3Prefix := request.GetStringParam(req.Parameters, "S3Prefix")
	if !validateS3Prefix(s3Prefix) {
		return nil, ErrInvalidParameter
	}

	s3BucketOwner := request.GetStringParam(req.Parameters, "S3BucketOwner")
	if !validateS3BucketOwner(s3BucketOwner) {
		return nil, ErrInvalidParameter
	}

	s3SseKmsKeyId := request.GetStringParam(req.Parameters, "S3SseKmsKeyId")
	if !validateS3SseKmsKeyId(s3SseKmsKeyId) {
		return nil, ErrInvalidParameter
	}

	clientToken := request.GetStringParam(req.Parameters, "ClientToken")
	if !validateClientToken(clientToken) {
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

	if _, err := store.Tables().Get(tableName); err != nil {
		// Smithy declares TableNotFoundException (not the general
		// ResourceNotFoundException) for ExportTableToPointInTime.
		return nil, ErrTableNotFoundException
	}

	pitr, err := store.Tables().GetPointInTimeRecovery(tableName)
	if err != nil || pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
		return nil, ErrPITRNotEnabled
	}

	// The export snapshots the table at the requested time, which must lie
	// inside the restorable window; an omitted time exports the present.
	now := time.Now()
	exportTime, hasExportTime := parseTimestampParam(req.Parameters, "ExportTime")
	if !hasExportTime {
		exportTime = now
	}
	if exportTime.Before(pitr.EarliestRestorableDateTime) || exportTime.After(now) {
		return nil, ErrInvalidExportTime
	}

	exportFormat := request.GetStringParam(req.Parameters, "ExportFormat")
	if exportFormat == "" {
		exportFormat = "DYNAMODB_JSON"
	}
	validExportFormats := map[string]bool{
		"DYNAMODB_JSON": true,
		"ION":           true,
	}
	if !validExportFormats[exportFormat] {
		return nil, ErrInvalidParameter
	}

	export, err := s.exportTableCore(ctx, store, ExportTableCoreInput{
		TableArn:     tableArn,
		TableName:    tableName,
		ExportFormat: exportFormat,
		S3Bucket:     s3Bucket,
		S3Prefix:     s3Prefix,
		ClientToken:  clientToken,
		Region:       reqCtx.GetRegion(),
		ExportTime:   exportTime,
	})
	if err != nil {
		return nil, err
	}

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
	exportArn := request.GetStringParam(req.Parameters, "ExportArn")
	if !validateExportArn(exportArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	export, err := s.describeExportCore(store, exportArn)
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
	tableArn := request.GetStringParam(req.Parameters, "TableArn")
	if tableArn != "" {
		if !validateTableArn(tableArn) {
			return nil, ErrInvalidParameter
		}
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := listExportsMaxLimit
	if _, ok := req.Parameters["MaxResults"]; ok {
		v := request.GetIntParam(req.Parameters, "MaxResults")
		if !validateListExportsLimit(v) {
			return nil, ErrInvalidParameter
		}
		maxResults = v
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	exports, nextToken, err := s.listExportsCore(store, tableArn, nextToken, maxResults)
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
	if !validateS3Bucket(s3Bucket) {
		return nil, ErrInvalidParameter
	}

	// The S3BucketSource member is named S3KeyPrefix on the wire (the
	// Smithy member name; ImportTable has no jsonName override).
	s3Prefix, _ := s3BucketSourceParam["S3KeyPrefix"].(string)
	if !validateS3Prefix(s3Prefix) {
		return nil, ErrInvalidParameter
	}
	s3BucketOwner, _ := s3BucketSourceParam["S3BucketOwner"].(string)
	if !validateS3BucketOwner(s3BucketOwner) {
		return nil, ErrInvalidParameter
	}

	clientToken := request.GetStringParam(req.Parameters, "ClientToken")
	if !validateClientToken(clientToken) {
		return nil, ErrInvalidParameter
	}

	// Validate CSV-specific options when InputFormat=CSV.
	if inputFormat == "CSV" {
		if csvOpts, ok := req.Parameters["CsvOptions"].(map[string]interface{}); ok {
			if delim, ok := csvOpts["Delimiter"].(string); ok {
				if !validateCsvDelimiter(delim) {
					return nil, ErrInvalidParameter
				}
			}
			if headerList, ok := csvOpts["HeaderList"].([]interface{}); ok {
				if !validateCsvHeaderList(len(headerList)) {
					return nil, ErrInvalidParameter
				}
				for _, h := range headerList {
					hs, ok := h.(string)
					if !ok {
						return nil, ErrInvalidParameter
					}
					if !validateCsvHeader(hs) {
						return nil, ErrInvalidParameter
					}
				}
			}
		}
	}

	tableCreationParams, ok := req.Parameters["TableCreationParameters"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	tableName, _ := tableCreationParams["TableName"].(string)
	if !validateResourceName(tableName) {
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
		return nil, ErrInvalidParameter
	}

	attrDefs := parseAttributeDefinitions(tableCreationParams)
	if len(attrDefs) == 0 {
		return nil, ErrInvalidParameter
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

	importGSI, err := parseGlobalSecondaryIndexes(tableCreationParams)
	if err != nil {
		return nil, err
	}
	importLSI, err := parseLocalSecondaryIndexes(tableCreationParams)
	if err != nil {
		return nil, err
	}

	// Parse CSV options for the import loop (already validated above).
	var csvDelimiter string
	var csvHeaderList []string
	if inputFormat == "CSV" {
		if csvOpts, ok := req.Parameters["CsvOptions"].(map[string]interface{}); ok {
			csvDelimiter, _ = csvOpts["Delimiter"].(string)
			if hl, ok := csvOpts["HeaderList"].([]interface{}); ok {
				for _, h := range hl {
					if hs, ok := h.(string); ok {
						csvHeaderList = append(csvHeaderList, hs)
					}
				}
			}
		}
	}

	imp, err := s.importTableCore(store, reqCtx.GetRegion(), ImportTableCoreInput{
		TableName:      tableName,
		KeySchema:      keySchema,
		AttributeDefs:  attrDefs,
		BillingMode:    billingMode,
		ProvThroughput: provThroughput,
		GSI:            importGSI,
		LSI:            importLSI,
		InputFormat:    inputFormat,
		S3Bucket:       s3Bucket,
		S3Prefix:       s3Prefix,
		S3BucketOwner:  s3BucketOwner,
		ClientToken:    clientToken,
		CSVDelimiter:   csvDelimiter,
		CSVHeaderList:  csvHeaderList,
	})
	if err != nil {
		return nil, err
	}

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
	importArn := request.GetStringParam(req.Parameters, "ImportArn")
	if !validateImportArn(importArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	imp, err := s.describeImportCore(store, importArn)
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
	tableArn := request.GetStringParam(req.Parameters, "TableArn")
	if tableArn != "" {
		if !validateTableArn(tableArn) {
			return nil, ErrInvalidParameter
		}
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	if token := request.GetStringParam(req.Parameters, "NextToken"); token != "" {
		if !validateImportNextToken(token) {
			return nil, ErrInvalidParameter
		}
	}
	pageSize := listImportsMaxLimit
	if _, ok := req.Parameters["PageSize"]; ok {
		v := request.GetIntParam(req.Parameters, "PageSize")
		if !validateListImportsLimit(v) {
			return nil, ErrInvalidParameter
		}
		pageSize = v
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	imports, nextToken, err := s.listImportsCore(store, tableArn, nextToken, pageSize)
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

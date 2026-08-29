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
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Import / Export Core — single validation + persistence path for table
// import and export operations.
//
// These methods encapsulate import/export lifecycle logic, including S3
// invocations for data transfer. Both the HTTP API handlers
// (import_export_operations.go) and any future admin handler delegate to
// these methods.
// ---------------------------------------------------------------------------

// describeExportCore validates the request, then returns an export by ARN.
func (s *DynamoDBService) describeExportCore(ctx context.Context, reqCtx *request.RequestContext, exportArn string) (*dbstore.ExportDescription, error) {
	if !validateExportArn(exportArn) {
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
	return export, nil
}

// listExportsInput carries the raw wire parameters for ListExports.
type listExportsInput struct {
	Parameters map[string]interface{}
}

// listExportsCore validates the request, then returns a paginated list of
// exports filtered by table ARN.
func (s *DynamoDBService) listExportsCore(ctx context.Context, reqCtx *request.RequestContext, in listExportsInput) ([]*dbstore.ExportDescription, string, error) {
	tableArn := request.GetStringParam(in.Parameters, "TableArn")
	if tableArn != "" {
		if !validateTableArn(tableArn) {
			return nil, "", ErrInvalidParameter
		}
	}
	nextToken := pagination.GetMarker(in.Parameters, "NextToken")
	maxResults := listExportsMaxLimit
	if _, ok := in.Parameters["MaxResults"]; ok {
		v := request.GetIntParam(in.Parameters, "MaxResults")
		if !validateListExportsLimit(v) {
			return nil, "", ErrInvalidParameter
		}
		maxResults = v
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}
	return store.Exports().List(tableArn, nextToken, maxResults)
}

// describeImportCore validates the request, then returns an import by ARN.
func (s *DynamoDBService) describeImportCore(ctx context.Context, reqCtx *request.RequestContext, importArn string) (*dbstore.ImportTableDescription, error) {
	if !validateImportArn(importArn) {
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
	return imp, nil
}

// listImportsInput carries the raw wire parameters for ListImports.
type listImportsInput struct {
	Parameters map[string]interface{}
}

// listImportsCore validates the request, then returns a paginated list of
// imports filtered by table ARN.
func (s *DynamoDBService) listImportsCore(ctx context.Context, reqCtx *request.RequestContext, in listImportsInput) ([]*dbstore.ImportTableDescription, string, error) {
	tableArn := request.GetStringParam(in.Parameters, "TableArn")
	if tableArn != "" {
		if !validateTableArn(tableArn) {
			return nil, "", ErrInvalidParameter
		}
	}
	nextToken := pagination.GetMarker(in.Parameters, "NextToken")
	if token := request.GetStringParam(in.Parameters, "NextToken"); token != "" {
		if !validateImportNextToken(token) {
			return nil, "", ErrInvalidParameter
		}
	}
	pageSize := listImportsMaxLimit
	if _, ok := in.Parameters["PageSize"]; ok {
		v := request.GetIntParam(in.Parameters, "PageSize")
		if !validateListImportsLimit(v) {
			return nil, "", ErrInvalidParameter
		}
		pageSize = v
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}
	return store.Imports().List(tableArn, nextToken, pageSize)
}

// exportTableInput carries the raw wire parameters for
// ExportTableToPointInTime.
type exportTableInput struct {
	Parameters map[string]interface{}
}

// exportTableResult carries the created export record and the resolved
// export time for wire serialisation.
type exportTableResult struct {
	Export     *dbstore.ExportDescription
	ExportTime time.Time
}

// ExportTableCoreInput is the service-layer DTO for the export job.
type ExportTableCoreInput struct {
	TableArn     string
	TableName    string
	ExportFormat string
	S3Bucket     string
	S3Prefix     string
	ClientToken  string
	Region       string
	ExportTime   time.Time
}

// exportTableCore validates the request, then creates the export record in
// the IN_PROGRESS state and starts the export in the background; clients
// poll DescribeExport for the final state. The returned description is the
// initial one.
func (s *DynamoDBService) exportTableCore(ctx context.Context, reqCtx *request.RequestContext, in exportTableInput) (*exportTableResult, error) {
	tableArn := request.GetStringParam(in.Parameters, "TableArn")
	if !validateTableArn(tableArn) {
		return nil, ErrInvalidParameter
	}

	s3Bucket := request.GetStringParam(in.Parameters, "S3Bucket")
	if !validateS3Bucket(s3Bucket) {
		return nil, ErrInvalidParameter
	}

	s3Prefix := request.GetStringParam(in.Parameters, "S3Prefix")
	if !validateS3Prefix(s3Prefix) {
		return nil, ErrInvalidParameter
	}

	s3BucketOwner := request.GetStringParam(in.Parameters, "S3BucketOwner")
	if !validateS3BucketOwner(s3BucketOwner) {
		return nil, ErrInvalidParameter
	}

	s3SseKmsKeyId := request.GetStringParam(in.Parameters, "S3SseKmsKeyId")
	if !validateS3SseKmsKeyId(s3SseKmsKeyId) {
		return nil, ErrInvalidParameter
	}

	clientToken := request.GetStringParam(in.Parameters, "ClientToken")
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
		// ResourceNotFoundException) for exportTableToPointInTime.
		return nil, ErrTableNotFoundException
	}

	pitr, err := store.Tables().GetPointInTimeRecovery(tableName)
	if err != nil || pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
		return nil, ErrPITRNotEnabled
	}

	// The export snapshots the table at the requested time, which must lie
	// inside the restorable window; an omitted time exports the present.
	now := time.Now()
	exportTime, hasExportTime := parseTimestampParam(in.Parameters, "ExportTime")
	if !hasExportTime {
		exportTime = now
	}
	if exportTime.Before(pitr.EarliestRestorableDateTime) || exportTime.After(now) {
		return nil, ErrInvalidExportTime
	}

	exportFormat := request.GetStringParam(in.Parameters, "ExportFormat")
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

	job := ExportTableCoreInput{
		TableArn:     tableArn,
		TableName:    tableName,
		ExportFormat: exportFormat,
		S3Bucket:     s3Bucket,
		S3Prefix:     s3Prefix,
		ClientToken:  clientToken,
		Region:       reqCtx.GetRegion(),
		ExportTime:   exportTime,
	}

	export, err := store.Exports().Create(job.TableArn, job.TableName, job.ExportFormat)
	if err != nil {
		return nil, err
	}
	export.S3Bucket = job.S3Bucket
	export.S3Prefix = job.S3Prefix
	if err := store.Exports().Put(export); err != nil {
		return nil, err
	}

	s.bgWg.Add(1)
	go func() {
		defer s.bgWg.Done()
		defer resilience.RecoverPanic("dynamodb export job")
		s.runExportJob(store, job, export.ExportArn)
	}()
	return &exportTableResult{Export: export, ExportTime: exportTime}, nil
}

// runExportJob snapshots the table as of the export time, writes the data
// to S3, and records the export's final state.
func (s *DynamoDBService) runExportJob(store dbstore.DynamoDBStoreInterface, in ExportTableCoreInput, exportArn string) {
	export, err := store.Exports().Get(exportArn)
	if err != nil {
		logs.Error("Failed to load export record for job",
			logs.String("exportArn", exportArn), logs.Err(err))
		return
	}

	failExport := func(code, message string) {
		export.ExportStatus = "FAILED"
		export.FailureCode = code
		export.FailureMessage = message
		export.EndTime = time.Now()
		if err := store.Exports().Put(export); err != nil {
			logs.Error("Failed to persist failed export",
				logs.String("exportArn", exportArn), logs.Err(err))
		}
	}

	items, err := snapshotItemsAsOf(store, in.TableName, in.ExportTime)
	if err != nil {
		logs.Error("Failed to snapshot items for export",
			logs.String("tableName", in.TableName), logs.Err(err))
		failExport("InternalFailure", fmt.Sprintf("failed to snapshot items: %v", err))
		return
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	for _, item := range items {
		line, mErr := json.Marshal(buildDynamoDBJSONItem(item))
		if mErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to marshal item: %v", mErr))
			return
		}
		if _, wErr := writer.Write(line); wErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to write item: %v", wErr))
			return
		}
		if _, wErr := writer.Write([]byte("\n")); wErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to write item: %v", wErr))
			return
		}
	}
	if flushErr := writer.Flush(); flushErr != nil {
		failExport("InternalFailure", fmt.Sprintf("failed to flush buffer: %v", flushErr))
		return
	}

	if s3 := s.s3invoker(); s3 != nil {
		dataFileID := fmt.Sprintf("%d", time.Now().UnixNano())
		objectKey := fmt.Sprintf("AWSDynamoDB/%s/data/%s.json", export.ExportArn, dataFileID)
		if in.S3Prefix != "" {
			objectKey = in.S3Prefix + "/" + objectKey
		}
		if putErr := s3.PutObject(s.bgCtx, in.Region, in.S3Bucket, objectKey, buf.Bytes(), "application/octet-stream"); putErr != nil {
			logs.Error("Failed to write export to S3",
				logs.Err(putErr),
				logs.String("bucket", in.S3Bucket),
				logs.String("key", objectKey))
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
	}

	export.ExportStatus = "COMPLETED"
	export.EndTime = time.Now()
	export.ExportTime = in.ExportTime
	export.ItemCount = int64(len(items))
	export.BilledSizeBytes = int64(buf.Len())
	if err := store.Exports().Put(export); err != nil {
		logs.Error("Failed to persist completed export",
			logs.String("exportArn", exportArn), logs.Err(err))
	}
}

// snapshotItemsAsOf returns the table's items at the given time: the
// current state with every journaled mutation newer than the point undone.
// The undo replays newest first, so each before-image overwrites the state
// the newer mutations left behind.
func snapshotItemsAsOf(store dbstore.DynamoDBStoreInterface, tableName string, at time.Time) ([]*dbstore.Item, error) {
	items := make(map[string]*dbstore.Item)
	order := make([]string, 0)
	if err := store.Items().Scan(tableName, func(item *dbstore.Item) error {
		key := itemKeyString(item.Key)
		if _, seen := items[key]; !seen {
			order = append(order, key)
		}
		items[key] = item
		return nil
	}); err != nil {
		return nil, err
	}

	if err := store.Journal().ReverseReplay(tableName, at, func(change *dbstore.JournalChange) error {
		key := itemKeyString(change.Key)
		if change.BeforeImage == nil {
			delete(items, key)
			return nil
		}
		if _, seen := items[key]; !seen {
			order = append(order, key)
		}
		items[key] = &dbstore.Item{
			TableName:  tableName,
			Key:        change.Key,
			Attributes: change.BeforeImage,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	snapshot := make([]*dbstore.Item, 0, len(order))
	for _, key := range order {
		if item, present := items[key]; present {
			snapshot = append(snapshot, item)
		}
	}
	return snapshot, nil
}

// itemKeyString renders a primary key map as a canonical string. Go's JSON
// encoder writes map keys in sorted order, so equal keys always render
// equally.
func itemKeyString(key map[string]*dbstore.AttributeValue) string {
	data, err := json.Marshal(key)
	if err != nil {
		return ""
	}
	return string(data)
}

// importTableInput carries the raw wire parameters for ImportTable.
type importTableInput struct {
	Parameters map[string]interface{}
}

// importTableResult carries the created import record and the echoed S3
// source members for wire serialisation.
type importTableResult struct {
	Import        *dbstore.ImportTableDescription
	S3Bucket      string
	S3Prefix      string
	S3BucketOwner string
}

// ImportTableCoreInput is the service-layer DTO for the import job. It
// contains all parameters needed to create the target table and import
// data from S3.
type ImportTableCoreInput struct {
	TableName      string
	KeySchema      []*dbstore.KeySchemaElement
	AttributeDefs  []*dbstore.AttributeDefinition
	BillingMode    dbstore.BillingMode
	ProvThroughput *dbstore.ProvisionedThroughput
	GSI            []*dbstore.GlobalSecondaryIndex
	LSI            []*dbstore.LocalSecondaryIndex
	InputFormat    string
	S3Bucket       string
	S3Prefix       string
	S3BucketOwner  string
	ClientToken    string
	CSVDelimiter   string
	CSVHeaderList  []string
}

// importTableCore validates the request, then creates the import record in
// the IN_PROGRESS state and runs the import in the background; clients poll
// DescribeImport for the final state. The returned description is the
// initial one.
func (s *DynamoDBService) importTableCore(ctx context.Context, reqCtx *request.RequestContext, in importTableInput) (*importTableResult, error) {
	inputFormat := request.GetStringParam(in.Parameters, "InputFormat")
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

	s3BucketSourceParam, ok := in.Parameters["S3BucketSource"].(map[string]interface{})
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

	clientToken := request.GetStringParam(in.Parameters, "ClientToken")
	if !validateClientToken(clientToken) {
		return nil, ErrInvalidParameter
	}

	// Validate CSV-specific options when InputFormat=CSV.
	if inputFormat == "CSV" {
		if csvOpts, ok := in.Parameters["CsvOptions"].(map[string]interface{}); ok {
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

	tableCreationParams, ok := in.Parameters["TableCreationParameters"].(map[string]interface{})
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
		if csvOpts, ok := in.Parameters["CsvOptions"].(map[string]interface{}); ok {
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

	job := ImportTableCoreInput{
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
	}

	tableArn := store.Tables().ARNBuilder().Table(job.TableName)
	imp, err := store.Imports().Create(tableArn, job.TableName)
	if err != nil {
		return nil, err
	}

	s.bgWg.Add(1)
	go func() {
		defer s.bgWg.Done()
		defer resilience.RecoverPanic("dynamodb import job")
		s.runImportJob(store, reqCtx.GetRegion(), job, imp.ImportArn)
	}()
	return &importTableResult{
		Import:        imp,
		S3Bucket:      s3Bucket,
		S3Prefix:      s3Prefix,
		S3BucketOwner: s3BucketOwner,
	}, nil
}

// runImportJob creates the target table, reads the data from S3, writes the
// items, and records the import's final state.
func (s *DynamoDBService) runImportJob(store dbstore.DynamoDBStoreInterface, region string, in ImportTableCoreInput, importArn string) {
	imp, err := store.Imports().Get(importArn)
	if err != nil {
		logs.Error("Failed to load import record for job",
			logs.String("importArn", importArn), logs.Err(err))
		return
	}
	imp.InputFormat = in.InputFormat
	imp.S3BucketSource = &dbstore.S3BucketSource{
		S3Bucket:      in.S3Bucket,
		S3Prefix:      in.S3Prefix,
		S3BucketOwner: in.S3BucketOwner,
	}

	failImport := func(code, message string) {
		imp.ImportStatus = "FAILED"
		imp.FailureCode = code
		imp.FailureMessage = message
		imp.EndTime = time.Now()
		if err := store.Imports().Put(imp); err != nil {
			logs.Error("Failed to persist failed import",
				logs.String("importArn", importArn), logs.Err(err))
		}
	}

	table, err := store.Tables().Create(
		in.TableName, in.KeySchema, in.AttributeDefs, in.BillingMode,
		in.ProvThroughput, in.GSI, in.LSI, nil, nil, false,
	)
	if err != nil {
		failImport("TableAlreadyExists", fmt.Sprintf("failed to create target table: %v", err))
		return
	}
	imp.TableArn = table.ARN

	importedCount := int64(0)
	processedSizeBytes := int64(0)

	s3 := s.s3invoker()
	if s3 != nil {
		// A missing source bucket is reported as a failed import rather
		// than an empty one, because listing cannot distinguish the two.
		exists, existsErr := s3.BucketExists(s.bgCtx, region, in.S3Bucket)
		if existsErr != nil {
			failImport("S3AccessDenied", fmt.Sprintf("failed to check source bucket: %v", existsErr))
			return
		}
		if !exists {
			failImport("S3NoSuchBucket", fmt.Sprintf("the source bucket %s does not exist", in.S3Bucket))
			return
		}
		keys, listErr := s3.ListObjects(s.bgCtx, region, in.S3Bucket, in.S3Prefix, 0)
		if listErr != nil {
			logs.Error("Failed to list S3 objects for import",
				logs.Err(listErr),
				logs.String("bucket", in.S3Bucket),
				logs.String("prefix", in.S3Prefix))
			failImport("S3AccessDenied", fmt.Sprintf("failed to list S3 objects: %v", listErr))
			return
		}
		for _, key := range keys {
			data, getErr := s3.GetObject(s.bgCtx, region, in.S3Bucket, key, 0)
			if getErr != nil {
				logs.Error("Failed to read S3 object for import",
					logs.Err(getErr),
					logs.String("bucket", in.S3Bucket),
					logs.String("key", key))
				failImport("S3AccessDenied", fmt.Sprintf("failed to read S3 object %s: %v", key, getErr))
				return
			}
			processedSizeBytes += int64(len(data))
			var count int64
			var parseErr error
			switch in.InputFormat {
			case "CSV":
				count, parseErr = importCSVData(s.bgCtx, data, in.TableName, store, in.CSVDelimiter, in.CSVHeaderList)
			case "ION":
				failImport("UnsupportedFormat", "ION format parsing is not yet implemented")
				return
			default:
				count, parseErr = importDynamoDBJSONData(s.bgCtx, data, in.TableName, store)
			}
			if parseErr != nil {
				logs.Error("Failed to parse import data",
					logs.Err(parseErr),
					logs.String("key", key))
				failImport("DocumentAccessException", fmt.Sprintf("failed to parse data in %s: %v", key, parseErr))
				return
			}
			importedCount += count
		}
	}

	imp.ImportStatus = "COMPLETED"
	imp.ProcessedItemCount = importedCount
	imp.ProcessedSizeBytes = processedSizeBytes
	imp.EndTime = time.Now()
	if err := store.Imports().Put(imp); err != nil {
		logs.Error("Failed to persist completed import",
			logs.String("importArn", importArn), logs.Err(err))
	}
}

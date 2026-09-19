package dynamodb

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
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

// clientTokenWindow is the documented idempotency window of an export or
// import ClientToken: a token is valid for eight hours after the first
// request that used it completed.
const clientTokenWindow = 8 * time.Hour

// clientTokenHash derives the idempotency payload hash of an export or
// import request: the canonical JSON of the request parameters with the
// ClientToken itself removed (map marshalling is key-sorted, so the hash is
// stable across identical retries).
func clientTokenHash(params map[string]interface{}) string {
	filtered := make(map[string]interface{}, len(params))
	for k, v := range params {
		if k != "ClientToken" {
			filtered[k] = v
		}
	}
	encoded, err := json.Marshal(filtered)
	if err != nil {
		encoded = []byte(fmt.Sprintf("%v", filtered))
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

// findExportByClientToken returns the export a ClientToken created, or nil
// when no export of the table carries the token.
func findExportByClientToken(store dbstore.DynamoDBStoreInterface, tableArn, clientToken string) *dbstore.ExportDescription {
	marker := ""
	for {
		exports, next, err := store.Exports().List(tableArn, marker, 100)
		if err != nil {
			return nil
		}
		for _, e := range exports {
			if e.ClientToken == clientToken {
				return e
			}
		}
		if next == "" {
			return nil
		}
		marker = next
	}
}

// findImportByClientToken returns the import a ClientToken created, or nil
// when no import of the table carries the token.
func findImportByClientToken(store dbstore.DynamoDBStoreInterface, tableArn, clientToken string) *dbstore.ImportTableDescription {
	marker := ""
	for {
		imports, next, err := store.Imports().List(tableArn, marker, 100)
		if err != nil {
			return nil
		}
		for _, i := range imports {
			if i.ClientToken == clientToken {
				return i
			}
		}
		if next == "" {
			return nil
		}
		marker = next
	}
}

// exportManifestFileEntry is one JSON line of manifest-files.json: the
// checksums and key of one exported data file.
type exportManifestFileEntry struct {
	ItemCount     int64  `json:"itemCount"`
	Md5Checksum   string `json:"md5Checksum"`
	Etag          string `json:"etag"`
	DataFileS3Key string `json:"dataFileS3Key"`
}

// exportManifestSummary is the manifest-summary.json object of a completed
// export (Developer Guide, DynamoDB table export output format).
type exportManifestSummary struct {
	Version            string  `json:"version"`
	ExportArn          string  `json:"exportArn"`
	StartTime          string  `json:"startTime"`
	EndTime            string  `json:"endTime"`
	TableArn           string  `json:"tableArn"`
	TableId            string  `json:"tableId"`
	ExportTime         string  `json:"exportTime"`
	S3Bucket           string  `json:"s3Bucket"`
	S3Prefix           string  `json:"s3Prefix"`
	S3SseAlgorithm     string  `json:"s3SseAlgorithm"`
	S3SseKmsKeyId      *string `json:"s3SseKmsKeyId"`
	ManifestFilesS3Key string  `json:"manifestFilesS3Key"`
	BilledSizeBytes    int64   `json:"billedSizeBytes"`
	ItemCount          int64   `json:"itemCount"`
	OutputFormat       string  `json:"outputFormat"`
}

// s3SseAlgorithmForManifest reports the SSE algorithm the export requested:
// KMS when a customer key was supplied, the AES256 default otherwise. The
// destination bucket's default-encryption policy governs what S3 applies at
// rest.
func s3SseAlgorithmForManifest(kmsKeyId string) string {
	if kmsKeyId != "" {
		return "KMS"
	}
	return "AES256"
}

// nilIfEmpty maps the empty string to a JSON null.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// decompressImportObject gunzips one GZIP-compressed import source object.
func decompressImportObject(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
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
	if exportTime.Before(pitrEarliestRestorable(pitr, now)) || exportTime.After(now) {
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

	// The ClientToken makes identical retries within the documented window
	// return the original export; a retried token with a changed payload is
	// the documented ExportConflictException.
	idemKey := "export:" + clientToken
	requestHash := clientTokenHash(in.Parameters)
	if clientToken != "" {
		recordedHash, _, _, found, lookupErr := store.Idempotency().Lookup(idemKey)
		if lookupErr != nil {
			return nil, lookupErr
		}
		if found {
			if recordedHash != requestHash {
				return nil, ErrExportConflict
			}
			if replay := findExportByClientToken(store, tableArn, clientToken); replay != nil {
				return &exportTableResult{Export: replay, ExportTime: replay.ExportTime}, nil
			}
			// The token record outlived the export itself (the export was
			// deleted); the call creates a replacement export.
		}
	}

	export, err := store.Exports().Create(job.TableArn, job.TableName, job.ExportFormat)
	if err != nil {
		return nil, err
	}
	export.S3Bucket = job.S3Bucket
	export.S3Prefix = job.S3Prefix
	export.S3BucketOwner = s3BucketOwner
	export.S3SseKmsKeyId = s3SseKmsKeyId
	export.ClientToken = clientToken
	export.ExportType = "FULL_EXPORT"
	if err := store.Exports().Put(export); err != nil {
		return nil, err
	}
	if clientToken != "" {
		if err := store.Idempotency().Record(idemKey, requestHash, dbstore.IdempotencyStateCompleted, time.Now().Add(clientTokenWindow), nil); err != nil {
			logs.Warn("failed to record export client token",
				logs.String("exportArn", export.ExportArn), logs.Err(err))
		}
	}

	s.bgWg.Add(1)
	go func() {
		defer s.bgWg.Done()
		defer func() {
			if r := recover(); r != nil {
				resilience.LogPanic("dynamodb export job", r)
			}
		}()
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
		// The S3 layout follows the documented export output format
		// (Developer Guide): every object lives under
		// <prefix>/AWSDynamoDB/<ExportId>/, whose ExportId is the export
		// ARN's trailing segment. The data is gzip-compressed JSON lines;
		// the manifests and their checksums describe it. Storage-side
		// encryption follows the destination bucket's default-encryption
		// policy; the manifest records the export's own SSE request.
		exportID := export.ExportArn[strings.LastIndex(export.ExportArn, "/")+1:]
		baseKey := "AWSDynamoDB/" + exportID
		if in.S3Prefix != "" {
			baseKey = in.S3Prefix + "/" + baseKey
		}
		putObject := func(key string, data []byte) error {
			return s3.PutObject(s.bgCtx, in.Region, in.S3Bucket, key, data, "application/octet-stream")
		}
		if putErr := putObject(baseKey+"/_started", nil); putErr != nil {
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}

		var compressed bytes.Buffer
		gzWriter := gzip.NewWriter(&compressed)
		if _, gzErr := gzWriter.Write(buf.Bytes()); gzErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to compress export data: %v", gzErr))
			return
		}
		if gzErr := gzWriter.Close(); gzErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to compress export data: %v", gzErr))
			return
		}
		dataKey := baseKey + "/data/" + fmt.Sprintf("%d.json.gz", time.Now().UnixNano())
		if putErr := putObject(dataKey, compressed.Bytes()); putErr != nil {
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}

		dataSum := md5.Sum(compressed.Bytes())
		manifestEntry := exportManifestFileEntry{
			ItemCount:     int64(len(items)),
			Md5Checksum:   base64.StdEncoding.EncodeToString(dataSum[:]),
			Etag:          hex.EncodeToString(dataSum[:]),
			DataFileS3Key: dataKey,
		}
		manifestFiles, mErr := json.Marshal(manifestEntry)
		if mErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to marshal files manifest: %v", mErr))
			return
		}
		manifestFilesKey := baseKey + "/manifest-files.json"
		if putErr := putObject(manifestFilesKey, append(manifestFiles, '\n')); putErr != nil {
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
		filesSum := md5.Sum(append(manifestFiles, '\n'))
		if putErr := putObject(baseKey+"/manifest-files.checksum", []byte(base64.StdEncoding.EncodeToString(filesSum[:]))); putErr != nil {
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}

		summary := exportManifestSummary{
			Version:            "2020-06-30",
			ExportArn:          export.ExportArn,
			StartTime:          export.StartTime.UTC().Format("2006-01-02T15:04:05.000Z"),
			EndTime:            time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			TableArn:           export.TableArn,
			TableId:            export.TableId,
			ExportTime:         in.ExportTime.UTC().Format("2006-01-02T15:04:05.000Z"),
			S3Bucket:           in.S3Bucket,
			S3Prefix:           in.S3Prefix,
			S3SseAlgorithm:     s3SseAlgorithmForManifest(export.S3SseKmsKeyId),
			S3SseKmsKeyId:      nilIfEmpty(export.S3SseKmsKeyId),
			ManifestFilesS3Key: manifestFilesKey,
			BilledSizeBytes:    int64(buf.Len()),
			ItemCount:          int64(len(items)),
			OutputFormat:       export.ExportFormat,
		}
		summaryBytes, sumErr := json.Marshal(summary)
		if sumErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to marshal summary manifest: %v", sumErr))
			return
		}
		if putErr := putObject(baseKey+"/manifest-summary.json", summaryBytes); putErr != nil {
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
		summarySum := md5.Sum(summaryBytes)
		if putErr := putObject(baseKey+"/manifest-summary.checksum", []byte(base64.StdEncoding.EncodeToString(summarySum[:]))); putErr != nil {
			failExport("S3AccessDenied", fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
		export.ExportManifest = "s3://" + in.S3Bucket + "/" + manifestFilesKey
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
	TableName            string
	KeySchema            []*dbstore.KeySchemaElement
	AttributeDefs        []*dbstore.AttributeDefinition
	BillingMode          dbstore.BillingMode
	ProvThroughput       *dbstore.ProvisionedThroughput
	GSI                  []*dbstore.GlobalSecondaryIndex
	LSI                  []*dbstore.LocalSecondaryIndex
	InputFormat          string
	InputCompressionType string
	S3Bucket             string
	S3Prefix             string
	S3BucketOwner        string
	ClientToken          string
	CSVDelimiter         string
	CSVHeaderList        []string
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

	// InputCompressionType defaults to NONE; GZIP is the only other modelled
	// value and decompresses each source object before parsing.
	inputCompressionType := request.GetStringParam(in.Parameters, "InputCompressionType")
	if inputCompressionType == "" {
		inputCompressionType = "NONE"
	}
	if inputCompressionType != "NONE" && inputCompressionType != "GZIP" {
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
		TableName:            tableName,
		KeySchema:            keySchema,
		AttributeDefs:        attrDefs,
		BillingMode:          billingMode,
		ProvThroughput:       provThroughput,
		GSI:                  importGSI,
		LSI:                  importLSI,
		InputFormat:          inputFormat,
		InputCompressionType: inputCompressionType,
		S3Bucket:             s3Bucket,
		S3Prefix:             s3Prefix,
		S3BucketOwner:        s3BucketOwner,
		ClientToken:          clientToken,
		CSVDelimiter:         csvDelimiter,
		CSVHeaderList:        csvHeaderList,
	}

	tableArn := store.Tables().ARNBuilder().Table(job.TableName)

	// The ClientToken makes identical retries within the documented window
	// return the original import; a retried token with a changed payload is
	// the documented IdempotentParameterMismatchException.
	idemKey := "import:" + clientToken
	requestHash := clientTokenHash(in.Parameters)
	if clientToken != "" {
		recordedHash, _, _, found, lookupErr := store.Idempotency().Lookup(idemKey)
		if lookupErr != nil {
			return nil, lookupErr
		}
		if found {
			if recordedHash != requestHash {
				return nil, ErrIdempotentParameterMismatch
			}
			if replay := findImportByClientToken(store, tableArn, clientToken); replay != nil {
				result := &importTableResult{Import: replay}
				if replay.S3BucketSource != nil {
					result.S3Bucket = replay.S3BucketSource.S3Bucket
					result.S3Prefix = replay.S3BucketSource.S3Prefix
					result.S3BucketOwner = replay.S3BucketSource.S3BucketOwner
				}
				return result, nil
			}
			// The token record outlived the import itself (the import table
			// was deleted); the call creates a replacement import.
		}
	}

	imp, err := store.Imports().Create(tableArn, job.TableName)
	if err != nil {
		return nil, err
	}
	imp.ClientToken = clientToken
	imp.InputCompressionType = inputCompressionType
	if err := store.Imports().Put(imp); err != nil {
		return nil, err
	}
	if clientToken != "" {
		if err := store.Idempotency().Record(idemKey, requestHash, dbstore.IdempotencyStateCompleted, time.Now().Add(clientTokenWindow), nil); err != nil {
			logs.Warn("failed to record import client token",
				logs.String("importArn", imp.ImportArn), logs.Err(err))
		}
	}

	s.bgWg.Add(1)
	go func() {
		defer s.bgWg.Done()
		defer func() {
			if r := recover(); r != nil {
				resilience.LogPanic("dynamodb import job", r)
			}
		}()
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

	table, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                   in.TableName,
		KeySchema:              in.KeySchema,
		AttributeDefinitions:   in.AttributeDefs,
		BillingMode:            in.BillingMode,
		ProvisionedThroughput:  in.ProvThroughput,
		GlobalSecondaryIndexes: in.GSI,
		LocalSecondaryIndexes:  in.LSI,
	})
	if err != nil {
		failImport("TableAlreadyExists", fmt.Sprintf("failed to create target table: %v", err))
		return
	}
	imp.TableArn = table.ARN

	importedCount := int64(0)
	processedCount := int64(0)
	errorCount := int64(0)
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
			if in.InputCompressionType == "GZIP" {
				uncompressed, gzErr := decompressImportObject(data)
				if gzErr != nil {
					failImport("GZIPError", fmt.Sprintf("failed to decompress S3 object %s: %v", key, gzErr))
					return
				}
				data = uncompressed
			}
			processedSizeBytes += int64(len(data))
			var counts importCounts
			var parseErr error
			switch in.InputFormat {
			case "CSV":
				counts, parseErr = importCSVData(s.bgCtx, data, in.TableName, store, in.CSVDelimiter, in.CSVHeaderList)
			default:
				counts, parseErr = importDynamoDBJSONData(s.bgCtx, data, in.TableName, store)
			}
			if parseErr != nil {
				logs.Error("Failed to parse import data",
					logs.Err(parseErr),
					logs.String("key", key))
				failImport("DocumentAccessException", fmt.Sprintf("failed to parse data in %s: %v", key, parseErr))
				return
			}
			processedCount += counts.processed
			importedCount += counts.imported
			errorCount += counts.errors
		}
	}

	imp.ProcessedItemCount = processedCount
	imp.ImportedItemCount = importedCount
	imp.ErrorCount = errorCount
	imp.ProcessedSizeBytes = processedSizeBytes
	imp.EndTime = time.Now()
	// Items that failed validation were skipped and counted as errors; the
	// job reports the failure so the counts surface to the caller.
	if errorCount > 0 {
		imp.ImportStatus = "FAILED"
		imp.FailureCode = "ItemValidationError"
		imp.FailureMessage = "Some of the items failed validation checks and were not imported."
	} else {
		imp.ImportStatus = "COMPLETED"
	}
	if err := store.Imports().Put(imp); err != nil {
		logs.Error("Failed to persist completed import",
			logs.String("importArn", importArn), logs.Err(err))
	}
}

// importWriteItem writes one imported item in a single transaction with
// index entries and honest table metrics: an overwrite retires the replaced
// item's index entries, adjusts the table size by the delta, and does not
// increment the item count, so the table's reported counts match the
// distinct items it holds.
func importWriteItem(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, key, attrs map[string]*dbstore.AttributeValue) error {
	return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existing, err := txn.GetItem(table.Name, key)
		if err != nil && !errors.Is(err, dbstore.ErrItemNotFound) {
			return err
		}
		var oldSize int64
		if existing != nil {
			oldSize = dbstore.CalculateItemSize(existing.Attributes)
		}
		return txn.StoreItemWrite(table.Name, key, attrs, existing, existing != nil, oldSize)
	})
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

// importLineItem validates one parsed source line against the table schema
// and writes it. A line whose key attributes are incomplete or whose key
// attribute types disagree with the attribute definitions, and a line whose
// write fails, is counted as an item error and skipped — the import job
// continues with the next item.
func importLineItem(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, attrs map[string]*dbstore.AttributeValue, counts *importCounts) {
	counts.processed++

	key := make(map[string]*dbstore.AttributeValue, len(table.KeySchema))
	for _, ks := range table.KeySchema {
		if v, ok := attrs[ks.AttributeName]; ok {
			key[ks.AttributeName] = v
		}
	}
	if len(key) != len(table.KeySchema) || validateKeyTypes(table, key) != nil {
		counts.errors++
		return
	}
	if err := importWriteItem(ctx, store, table, key, attrs); err != nil {
		counts.errors++
		return
	}
	counts.imported++
}

// importCounts carries the outcome tallies of an import pass. The three
// counts are separate members of the import description: processed counts
// every item attempted from the source file (validation failures included),
// imported counts items successfully loaded (duplicate keys overwrite and
// count on every occurrence — an overwrite is not an error), and errors
// count items that failed validation or loading. Lines that cannot be
// parsed into an item at all never become processed items and count only
// as errors.
type importCounts struct {
	processed int64
	imported  int64
	errors    int64
}

// importDynamoDBJSONData parses DYNAMODB_JSON data (newline-delimited
// {"Item": {...}} objects) and writes each item to the store with proper
// index entries and table metric updates.
func importDynamoDBJSONData(ctx context.Context, data []byte, tableName string, store dbstore.DynamoDBStoreInterface) (importCounts, error) {
	counts := importCounts{}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return counts, fmt.Errorf("get table %s for import: %w", tableName, err)
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
			counts.errors++
			continue
		}

		itemRaw, ok := itemWrapper["Item"]
		if !ok {
			counts.errors++
			continue
		}

		var itemMap map[string]interface{}
		if err := json.Unmarshal(itemRaw, &itemMap); err != nil {
			counts.errors++
			continue
		}

		attrs, parseErr := parseAttributeValueMap(itemMap)
		if parseErr != nil || len(attrs) == 0 {
			counts.errors++
			continue
		}

		importLineItem(ctx, store, table, attrs, &counts)
	}
	return counts, scanner.Err()
}

// importCSVData parses CSV-formatted data and imports rows as DynamoDB items.
// Each column becomes a String attribute. If headerList is empty, the first
// row is treated as the header. delimiter defaults to comma.
func importCSVData(ctx context.Context, data []byte, tableName string, store dbstore.DynamoDBStoreInterface, delimiter string, headerList []string) (importCounts, error) {
	counts := importCounts{}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return counts, fmt.Errorf("get table %s for CSV import: %w", tableName, err)
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
			return counts, fmt.Errorf("failed to read CSV header: %w", err)
		}
		headers = headerRecord
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			counts.errors++
			continue
		}

		attrs := make(map[string]*dbstore.AttributeValue)
		for i, val := range record {
			if i >= len(headers) {
				break
			}
			attrs[headers[i]] = &dbstore.AttributeValue{S: &val}
		}

		importLineItem(ctx, store, table, attrs, &counts)
	}
	return counts, nil
}

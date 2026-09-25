package dynamodb

import (
	"bufio"
	"bytes"
	"compress/gzip"

	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
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
// Export Core — single validation + persistence path for the export family
// (ExportTableToPointInTime, DescribeExport, ListExports) and the S3 job
// that writes the export artefacts. The HTTP API handlers
// (import_export_operations.go) delegate to these methods; the machinery
// shared with the import family lives in import_export_core.go.
// ---------------------------------------------------------------------------

// describeExportCore validates the request, then returns an export by ARN.
func (s *DynamoDBService) describeExportCore(ctx context.Context, reqCtx *request.RequestContext, exportArn string) (*dbstore.ExportDescription, error) {
	return describeByArn(s, reqCtx, exportArn, validateExportArn(exportArn), ErrExportNotFound, func(store dbstore.DynamoDBStoreInterface) (*dbstore.ExportDescription, error) {
		return store.Exports().Get(exportArn)
	})
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

// findExportByClientToken returns the export a ClientToken created, or nil
// when no export of the table carries the token.
func findExportByClientToken(store dbstore.DynamoDBStoreInterface, tableArn, clientToken string) *dbstore.ExportDescription {
	return findByClientToken(
		func(marker string) ([]*dbstore.ExportDescription, string, error) {
			return store.Exports().List(tableArn, marker, 100)
		},
		func(e *dbstore.ExportDescription) string { return e.ClientToken },
		clientToken,
	)
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
// export (Developer Guide, DynamoDB table export output format). The
// incremental-only outputView and exportType members follow the
// incremental manifest-summary example, which carries both alongside the
// shared fields; the full-export example shows neither.
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
	OutputView         *string `json:"outputView,omitempty"`
	ExportType         string  `json:"exportType,omitempty"`
}

// s3SseAlgorithmForManifest reports the SSE algorithm the export requested:
// the algorithm the request named, KMS when only a customer key was
// supplied, the AES256 default otherwise. The destination bucket's
// default-encryption policy governs what S3 applies at rest.
func s3SseAlgorithmForManifest(algorithm, kmsKeyId string) string {
	if algorithm != "" {
		return algorithm
	}
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

	s3SseAlgorithm := request.GetStringParam(in.Parameters, "S3SseAlgorithm")
	if !validateS3SseAlgorithm(s3SseAlgorithm) {
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

	table, err := store.Tables().Get(tableName)
	if err != nil {
		// Smithy declares TableNotFoundException (not the general
		// ResourceNotFoundException) for exportTableToPointInTime.
		return nil, ErrTableNotFoundException
	}

	pitr, err := store.Tables().GetPointInTimeRecovery(tableName)
	if err != nil {
		// A failed read is the storage error it is: the not-enabled
		// sentinel belongs to the enabled-state judgement below, and a
		// storage fault masquerading as "PITR not enabled" would send the
		// client to enable recovery against a healthy table.
		return nil, err
	}
	if pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
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

	// ExportType selects the export's data path (API reference): the
	// default FULL_EXPORT snapshots the whole table, while
	// INCREMENTAL_EXPORT requires the IncrementalExportSpecification that
	// bounds the change window.
	exportType := request.GetStringParam(in.Parameters, "ExportType")
	if !validateExportType(exportType) {
		return nil, ErrInvalidParameter
	}
	if exportType == "" {
		exportType = "FULL_EXPORT"
	}
	var incSpec *incrementalExportSpec
	if raw, ok := in.Parameters["IncrementalExportSpecification"]; ok {
		incSpec, err = parseIncrementalExportSpecification(raw)
		if err != nil {
			return nil, err
		}
	}
	if exportType == "INCREMENTAL_EXPORT" && incSpec == nil {
		return nil, ErrInvalidParameter
	}
	if exportType == "FULL_EXPORT" && incSpec != nil {
		// The specification only carries meaning for the incremental data
		// path; accepting it on a full export would silently ignore a
		// member the client believes it sent.
		return nil, ErrInvalidParameter
	}
	if incSpec != nil {
		// The window must lie inside the restorable range the journal can
		// reconstruct, and the exclusive end cannot precede the inclusive
		// start. An omitted end resolves to the present, fixed here so the
		// recorded window and the echoed description agree.
		if incSpec.From.IsZero() {
			// The specification documents no default for an omitted start;
			// the earliest restorable instant is the complete window the
			// journal can reconstruct.
			incSpec.From = pitrEarliestRestorable(pitr, now)
		}
		if incSpec.From.Before(pitrEarliestRestorable(pitr, now)) || incSpec.From.After(now) {
			return nil, ErrInvalidExportTime
		}
		if incSpec.To.IsZero() {
			incSpec.To = now
		}
		if incSpec.To.Before(incSpec.From) || incSpec.To.After(now) {
			return nil, ErrInvalidExportTime
		}
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

	// The export record carries the table's minted UUID identity, not its
	// name: descriptions and the manifest tableId field state the pattern
	// the model documents for TableId.
	export, err := store.Exports().Create(job.TableArn, table.TableId, dbstore.ExportFormat(job.ExportFormat))
	if err != nil {
		return nil, err
	}
	export.S3Bucket = job.S3Bucket
	export.S3Prefix = job.S3Prefix
	export.S3BucketOwner = s3BucketOwner
	export.S3SseKmsKeyId = s3SseKmsKeyId
	export.S3SseAlgorithm = dbstore.S3SseAlgorithm(s3SseAlgorithm)
	export.ClientToken = clientToken
	export.ExportType = dbstore.ExportType(exportType)
	if incSpec != nil {
		export.ExportFromTime = incSpec.From
		export.ExportToTime = incSpec.To
		export.ExportViewType = dbstore.ExportViewType(incSpec.ViewType)
	}
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

	// The data path follows the export type: a full export serialises the
	// snapshot at the export time, an incremental export the change records
	// the journal window [from, to) yields.
	var items []*dbstore.Item
	var records []*incrementalExportRecord
	if export.ExportType == "INCREMENTAL_EXPORT" {
		records, err = buildIncrementalExportRecords(store, in.TableName, export.ExportFromTime, export.ExportToTime, string(export.ExportViewType))
		if err != nil {
			logs.Error("Failed to build incremental export records",
				logs.String("tableName", in.TableName), logs.Err(err))
			failExport("InternalFailure", fmt.Sprintf("failed to build incremental export records: %v", err))
			return
		}
	} else {
		items, err = snapshotItemsAsOf(store, in.TableName, in.ExportTime)
		if err != nil {
			logs.Error("Failed to snapshot items for export",
				logs.String("tableName", in.TableName), logs.Err(err))
			failExport("InternalFailure", fmt.Sprintf("failed to snapshot items: %v", err))
			return
		}
	}
	outputCount := len(items)
	if records != nil {
		outputCount = len(records)
	}

	// The serialiser follows the requested format: DynamoDB JSON lines, or
	// the documented Ion text form (one $ion_1_0-prefixed line per item)
	// when ExportFormat=ION.
	serializeLine := func(i int) ([]byte, error) {
		if export.ExportType == "INCREMENTAL_EXPORT" {
			if export.ExportFormat == "ION" {
				return buildIonIncrementalRecordLine(records[i]), nil
			}
			return json.Marshal(buildIncrementalWireRecord(records[i]))
		}
		if export.ExportFormat == "ION" {
			return buildIonItemLine(items[i]), nil
		}
		return json.Marshal(buildDynamoDBJSONItem(items[i]))
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	for i := 0; i < outputCount; i++ {
		line, mErr := serializeLine(i)
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
		// ARN's trailing segment. The data is gzip-compressed lines in the
		// requested format (.json.gz for DynamoDB JSON, .ion.gz for Ion);
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
			failExport(s3FailureCode(putErr), fmt.Sprintf("failed to write to S3: %v", putErr))
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
		dataExt := ".json.gz"
		if export.ExportFormat == "ION" {
			dataExt = ".ion.gz"
		}
		dataKey := baseKey + "/data/" + fmt.Sprintf("%d%s", time.Now().UnixNano(), dataExt)
		if putErr := putObject(dataKey, compressed.Bytes()); putErr != nil {
			failExport(s3FailureCode(putErr), fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}

		dataSum := md5.Sum(compressed.Bytes())
		manifestEntry := exportManifestFileEntry{
			ItemCount:   int64(outputCount),
			Md5Checksum: base64.StdEncoding.EncodeToString(dataSum[:]),
			// The documented manifest example carries the S3 multipart etag
			// form — the md5 hex followed by the part count — so the export
			// reports the same "-1" suffix.
			Etag:          hex.EncodeToString(dataSum[:]) + "-1",
			DataFileS3Key: dataKey,
		}
		manifestFiles, mErr := json.Marshal(manifestEntry)
		if mErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to marshal files manifest: %v", mErr))
			return
		}
		manifestFilesKey := baseKey + "/manifest-files.json"
		if putErr := putObject(manifestFilesKey, append(manifestFiles, '\n')); putErr != nil {
			failExport(s3FailureCode(putErr), fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
		filesSum := md5.Sum(append(manifestFiles, '\n'))
		if putErr := putObject(baseKey+"/manifest-files.checksum", []byte(base64.StdEncoding.EncodeToString(filesSum[:]))); putErr != nil {
			failExport(s3FailureCode(putErr), fmt.Sprintf("failed to write to S3: %v", putErr))
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
			S3SseAlgorithm:     s3SseAlgorithmForManifest(string(export.S3SseAlgorithm), export.S3SseKmsKeyId),
			S3SseKmsKeyId:      nilIfEmpty(export.S3SseKmsKeyId),
			ManifestFilesS3Key: manifestFilesKey,
			BilledSizeBytes:    int64(buf.Len()),
			ItemCount:          int64(outputCount),
			OutputFormat:       string(export.ExportFormat),
		}
		// The incremental summary carries the view type and the export
		// type alongside the shared fields (developer guide, incremental
		// manifest-summary example); the full-export example shows neither.
		if export.ExportType == "INCREMENTAL_EXPORT" {
			view := string(export.ExportViewType)
			summary.OutputView = &view
			summary.ExportType = string(export.ExportType)
		}
		summaryBytes, sumErr := json.Marshal(summary)
		if sumErr != nil {
			failExport("InternalFailure", fmt.Sprintf("failed to marshal summary manifest: %v", sumErr))
			return
		}
		if putErr := putObject(baseKey+"/manifest-summary.json", summaryBytes); putErr != nil {
			failExport(s3FailureCode(putErr), fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
		summarySum := md5.Sum(summaryBytes)
		if putErr := putObject(baseKey+"/manifest-summary.checksum", []byte(base64.StdEncoding.EncodeToString(summarySum[:]))); putErr != nil {
			failExport(s3FailureCode(putErr), fmt.Sprintf("failed to write to S3: %v", putErr))
			return
		}
		export.ExportManifest = "s3://" + in.S3Bucket + "/" + manifestFilesKey
	} else {
		// An unreachable S3 plane is the export destination's missing
		// bucket: a job that fell through to COMPLETED would fake success
		// over data written nowhere — the worst failure mode for a client
		// that deletes source data after a "successful" export.
		failExport("S3AccessDenied", "the S3 destination plane is unavailable")
		return
	}

	export.ExportStatus = "COMPLETED"
	export.EndTime = time.Now()
	export.ExportTime = in.ExportTime
	export.ItemCount = int64(outputCount)
	export.BilledSizeBytes = int64(buf.Len())
	if err := store.Exports().Put(export); err != nil {
		logs.Error("Failed to persist completed export",
			logs.String("exportArn", exportArn), logs.Err(err))
	}
}

// buildDynamoDBJSONItem converts a store Item to the DynamoDB JSON format
// used by import/export: {"Item": {"attr": {"S": "val"}, ...}}.
func buildDynamoDBJSONItem(item *dbstore.Item) map[string]interface{} {
	return map[string]interface{}{
		"Item": buildAttributeValueWireMap(fullAttributes(item)),
	}
}

// fullAttributes returns an item's complete attribute map in store form,
// key attributes included.
func fullAttributes(item *dbstore.Item) map[string]*dbstore.AttributeValue {
	merged := make(map[string]*dbstore.AttributeValue, len(item.Attributes)+len(item.Key))
	for k, v := range item.Key {
		merged[k] = v
	}
	for k, v := range item.Attributes {
		merged[k] = v
	}
	return merged
}

// buildAttributeValueWireMap converts a store attribute map to the
// DynamoDB JSON wire form: name → type descriptor → value.
func buildAttributeValueWireMap(attrs map[string]*dbstore.AttributeValue) map[string]interface{} {
	wire := make(map[string]interface{}, len(attrs))
	for k, v := range attrs {
		wire[k] = buildAttributeValueResponse(v)
	}
	return wire
}

// incrementalExportSpec is the resolved IncrementalExportSpecification of
// an incremental export (API reference): the change window and the image
// set the output carries. From is the window's inclusive start, To its
// exclusive end; a zero To denotes the documented "latest time with data
// available", which the caller resolves to the present.
type incrementalExportSpec struct {
	From     time.Time
	To       time.Time
	ViewType string
}

// parseIncrementalExportSpecification validates and resolves the wire form
// of the incremental export specification. The view type defaults to
// NEW_AND_OLD_IMAGES, the documented default.
func parseIncrementalExportSpecification(raw interface{}) (*incrementalExportSpec, error) {
	obj, ok := raw.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	spec := &incrementalExportSpec{}
	if from, ok := parseTimestampParam(obj, "ExportFromTime"); ok {
		spec.From = from
	}
	if to, ok := parseTimestampParam(obj, "ExportToTime"); ok {
		spec.To = to
	}
	viewType := request.GetStringParam(obj, "ExportViewType")
	if !validateExportViewType(viewType) {
		return nil, ErrInvalidParameter
	}
	if viewType == "" {
		viewType = "NEW_AND_OLD_IMAGES"
	}
	spec.ViewType = viewType
	return spec, nil
}

// incrementalExportRecord is one output record of an incremental export:
// the item's last write time inside the window, its keys, and the new and
// old images the view type selected. The images are store attribute maps;
// the DynamoDB JSON and Ion serialisers render them.
type incrementalExportRecord struct {
	WriteTimestampMicros int64
	Keys                 map[string]*dbstore.AttributeValue
	NewImage             map[string]*dbstore.AttributeValue
	OldImage             map[string]*dbstore.AttributeValue
}

// buildIncrementalExportRecords constructs the change records of the
// window [from, to) (Developer Guide, DynamoDB table export output
// format). The keys the journal reports as changed in the window are
// classified against the table state just prior to each boundary: present
// at the end only is an insert, present at both an update, present at the
// start only a delete, and present at neither the insert+delete pair the
// guide documents as producing no output. The boundary snapshots and the
// window share one clock, so a key's classification and its images always
// come from the same instant.
func buildIncrementalExportRecords(store dbstore.DynamoDBStoreInterface, tableName string, from, to time.Time, viewType string) ([]*incrementalExportRecord, error) {
	oldItems, err := snapshotItemsAsOf(store, tableName, from.Add(-time.Nanosecond))
	if err != nil {
		return nil, err
	}
	newItems, err := snapshotItemsAsOf(store, tableName, to.Add(-time.Nanosecond))
	if err != nil {
		return nil, err
	}
	oldState := make(map[string]*dbstore.Item, len(oldItems))
	for _, item := range oldItems {
		oldState[itemKeyString(item.Key)] = item
	}
	newState := make(map[string]*dbstore.Item, len(newItems))
	for _, item := range newItems {
		newState[itemKeyString(item.Key)] = item
	}

	// One pass over the window, oldest first: the touched keys keep their
	// journal order, and a key's last record supplies the write timestamp
	// the record's Metadata carries.
	var order []string
	keys := make(map[string]map[string]*dbstore.AttributeValue)
	stamps := make(map[string]time.Time)
	if err := store.Journal().WindowReplay(tableName, from, to, func(change *dbstore.JournalChange) error {
		key := itemKeyString(change.Key)
		if _, seen := keys[key]; !seen {
			order = append(order, key)
			keys[key] = change.Key
		}
		stamps[key] = change.Timestamp
		return nil
	}); err != nil {
		return nil, err
	}

	includeOld := viewType == "NEW_AND_OLD_IMAGES"
	records := make([]*incrementalExportRecord, 0, len(order))
	for _, key := range order {
		newItem, newPresent := newState[key]
		oldItem, oldPresent := oldState[key]
		switch {
		case newPresent:
			// Insert or update: both view types carry the new image; the
			// old image rides along only for an update under
			// NEW_AND_OLD_IMAGES — an insert has none.
			record := &incrementalExportRecord{
				WriteTimestampMicros: stamps[key].UnixMicro(),
				Keys:                 keys[key],
				NewImage:             fullAttributes(newItem),
			}
			if includeOld && oldPresent {
				record.OldImage = fullAttributes(oldItem)
			}
			records = append(records, record)
		case oldPresent:
			// Delete: the new-images-only form is the bare keys; the
			// paired form adds the old image.
			record := &incrementalExportRecord{
				WriteTimestampMicros: stamps[key].UnixMicro(),
				Keys:                 keys[key],
			}
			if includeOld {
				record.OldImage = fullAttributes(oldItem)
			}
			records = append(records, record)
		default:
			// Insert and delete inside the window: no output.
		}
	}
	return records, nil
}

// incrementalWireMetadata renders the record's Metadata member: the write
// timestamp as a DynamoDB number attribute, the form the guide's examples
// show.
type incrementalWireMetadata struct {
	WriteTimestampMicros map[string]interface{} `json:"WriteTimestampMicros"`
}

// incrementalWireRecord is the DynamoDB JSON form of one incremental
// export record: Metadata, Keys, and the images the view type selected —
// absent members are omitted, matching the per-operation structures the
// guide's view-type table documents.
type incrementalWireRecord struct {
	Metadata *incrementalWireMetadata `json:"Metadata"`
	Keys     map[string]interface{}   `json:"Keys"`
	NewImage map[string]interface{}   `json:"NewImage,omitempty"`
	OldImage map[string]interface{}   `json:"OldImage,omitempty"`
}

func buildIncrementalWireRecord(record *incrementalExportRecord) *incrementalWireRecord {
	wire := &incrementalWireRecord{
		Metadata: &incrementalWireMetadata{
			WriteTimestampMicros: map[string]interface{}{
				"N": strconv.FormatInt(record.WriteTimestampMicros, 10),
			},
		},
		Keys: buildAttributeValueWireMap(record.Keys),
	}
	if record.NewImage != nil {
		wire.NewImage = buildAttributeValueWireMap(record.NewImage)
	}
	if record.OldImage != nil {
		wire.OldImage = buildAttributeValueWireMap(record.OldImage)
	}
	return wire
}

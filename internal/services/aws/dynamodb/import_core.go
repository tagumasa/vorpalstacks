package dynamodb

import (
	"bufio"
	"bytes"
	"compress/gzip"

	"context"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/klauspost/compress/zstd"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Import Core — single validation + persistence path for the import family
// (ImportTable, DescribeImport, ListImports) and the S3 job that reads the
// source data into the target table. The HTTP API handlers
// (import_export_operations.go) delegate to these methods; the machinery
// shared with the export family lives in import_export_core.go.
// ---------------------------------------------------------------------------

// describeImportCore validates the request, then returns an import by ARN.
func (s *DynamoDBService) describeImportCore(ctx context.Context, reqCtx *request.RequestContext, importArn string) (*dbstore.ImportTableDescription, error) {
	return describeByArn(s, reqCtx, importArn, validateImportArn(importArn), ErrImportNotFound, func(store dbstore.DynamoDBStoreInterface) (*dbstore.ImportTableDescription, error) {
		return store.Imports().Get(importArn)
	})
}

// listImportsInput carries the raw wire parameters for ListImports.
type listImportsInput struct {
	Parameters map[string]interface{}
}

// encodeImportListToken renders an import record key as the pagination
// token ListImports exchanges. The model types the token as 112-1024
// hexadecimal characters in sixteen-character blocks, so the key travels
// as its hex encoding, left-padded with zeros to the block alignment and
// the documented minimum; the padding bytes never occur at the start of
// a real ARN, so decoding strips them without ambiguity.
func encodeImportListToken(key string) string {
	hexKey := hex.EncodeToString([]byte(key))
	for len(hexKey) < importNextTokenMinLen || len(hexKey)%16 != 0 {
		hexKey = "0" + hexKey
	}
	return hexKey
}

// decodeImportListToken reverses encodeImportListToken for a token that
// already passed the model's pattern validation; a malformed hex string
// answers the parameter error.
func decodeImportListToken(token string) (string, error) {
	raw, err := hex.DecodeString(token)
	if err != nil {
		return "", ErrInvalidParameter
	}
	return string(bytes.TrimLeft(raw, "\x00")), nil
}

// listImportsCore validates the request, then returns a paginated list of
// imports filtered by table ARN. The store keys imports by their ARN, so
// the pagination marker crosses the wire as the token-encoded key and
// returns the same way — a raw ARN marker would fail the model's
// next-token pattern on the very next page.
func (s *DynamoDBService) listImportsCore(ctx context.Context, reqCtx *request.RequestContext, in listImportsInput) ([]*dbstore.ImportTableDescription, string, error) {
	tableArn := request.GetStringParam(in.Parameters, "TableArn")
	if tableArn != "" {
		if !validateTableArn(tableArn) {
			return nil, "", ErrInvalidParameter
		}
	}
	marker := pagination.GetMarker(in.Parameters, "NextToken")
	if token := request.GetStringParam(in.Parameters, "NextToken"); token != "" {
		if !validateImportNextToken(token) {
			return nil, "", ErrInvalidParameter
		}
		key, decodeErr := decodeImportListToken(token)
		if decodeErr != nil {
			return nil, "", decodeErr
		}
		marker = key
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
	imports, nextMarker, err := store.Imports().List(tableArn, marker, pageSize)
	if err != nil {
		return nil, "", err
	}
	if nextMarker != "" {
		nextMarker = encodeImportListToken(nextMarker)
	}
	return imports, nextMarker, nil
}

// findImportByClientToken returns the import a ClientToken created, or nil
// when no import of the table carries the token.
func findImportByClientToken(store dbstore.DynamoDBStoreInterface, tableArn, clientToken string) *dbstore.ImportTableDescription {
	return findByClientToken(
		func(marker string) ([]*dbstore.ImportTableDescription, string, error) {
			return store.Imports().List(tableArn, marker, 100)
		},
		func(i *dbstore.ImportTableDescription) string { return i.ClientToken },
		clientToken,
	)
}

// decompressImportObject decompresses one compressed import source object.
// GZIP and ZSTD are the compressed formats the model's
// InputCompressionType enum defines and the developer guide names ("Data
// can be compressed in ZSTD or GZIP format"); the failure code names the
// format that failed to decompress, following the GZIPError precedent.
func decompressImportObject(data []byte, compressionType string) ([]byte, error) {
	switch compressionType {
	case "GZIP":
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)
	case "ZSTD":
		reader, err := zstd.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)
	default:
		return data, nil
	}
}

// importCreateFailureCode names the failure code for a target-table
// creation failure: the duplicate-name race reports TableAlreadyExists,
// any other store failure is internal.
func importCreateFailureCode(err error) string {
	if dbstore.IsTableAlreadyExists(err) {
		return "TableAlreadyExists"
	}
	return "InternalFailure"
}

// importTableInput carries the raw wire parameters for ImportTable.
type importTableInput struct {
	Parameters map[string]interface{}
}

// importTableResult carries the created import record; the record itself
// holds every member the response faces render.
type importTableResult struct {
	Import *dbstore.ImportTableDescription
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
	SSEDescription       *dbstore.SSEDescription
	OnDemandThroughput   *dbstore.OnDemandThroughput
	InputFormat          string
	InputCompressionType string
	S3Bucket             string
	S3Prefix             string
	S3BucketOwner        string
	ClientToken          string
	CSVDelimiter         string
	CSVHeaderList        []string
}

// extractCSVOptions reads the CSV processing options from the request's
// InputFormatOptions member: the model nests CsvOptions under
// InputFormatOptions.Csv, and the wire names carry no jsonName override.
// An absent member leaves the defaults (comma delimiter, header taken
// from the first data row).
func extractCSVOptions(params map[string]interface{}) map[string]interface{} {
	opts, ok := params["InputFormatOptions"].(map[string]interface{})
	if !ok {
		return nil
	}
	csv, ok := opts["Csv"].(map[string]interface{})
	if !ok {
		return nil
	}
	return csv
}

// importTableCore validates the request, then creates the import record in
// the IN_PROGRESS state and runs the import in the background; clients poll
// DescribeImport for the final state. The returned description is the
// initial one.
func (s *DynamoDBService) importTableCore(ctx context.Context, reqCtx *request.RequestContext, in importTableInput) (*importTableResult, error) {
	// InputFormat is a required member (the model marks it so): an
	// omitted member is a validation failure, never a silently applied
	// default; a present non-string value is the same class of failure.
	inputFormat, formatPresent, formatErr := wireStringMember(in.Parameters["InputFormat"])
	if formatErr != nil {
		return nil, formatErr
	}
	if !formatPresent || inputFormat == "" {
		return nil, ErrInvalidParameter
	}

	validFormats := map[string]bool{
		"DYNAMODB_JSON": true,
		"CSV":           true,
	}
	if !validFormats[inputFormat] {
		return nil, ErrInvalidParameter
	}

	// InputCompressionType defaults to NONE; GZIP and ZSTD are the other
	// modelled values and each decompresses the source objects before
	// parsing (developer guide: "Data can be compressed in ZSTD or GZIP
	// format, or can be directly imported in uncompressed form").
	inputCompressionType := request.GetStringParam(in.Parameters, "InputCompressionType")
	if inputCompressionType == "" {
		inputCompressionType = "NONE"
	}
	switch inputCompressionType {
	case "NONE", "GZIP", "ZSTD":
	default:
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

	// The CSV options are extracted exactly once: validated where the
	// request's other members are checked, and the same parsed values feed
	// the import job below. The model nests them under
	// InputFormatOptions.Csv (no jsonName override), so the wire member is
	// InputFormatOptions and its Csv child carries the processing options.
	var csvDelimiter string
	var csvHeaderList []string
	if inputFormat == "CSV" {
		if csvOpts := extractCSVOptions(in.Parameters); csvOpts != nil {
			if delim, ok := csvOpts["Delimiter"].(string); ok {
				if !validateCsvDelimiter(delim) {
					return nil, ErrInvalidParameter
				}
				csvDelimiter = delim
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
					csvHeaderList = append(csvHeaderList, hs)
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

	// The created table carries the same billing-mode contract every table
	// creation path enforces: the documented default is provisioned, and a
	// provisioned request must carry a complete, quota-bounded throughput
	// pair — no invented inline pair, and no on-demand default of the
	// import's own. The indexes follow the mode the same way CreateTable's
	// indexes do.
	billingMode := dbstore.BillingMode(request.GetStringParam(tableCreationParams, "BillingMode"))
	provThroughput := parseProvisionedThroughput(tableCreationParams)
	billingMode, err = resolveBillingModeForCreation(billingMode, provThroughput)
	if err != nil {
		return nil, err
	}

	importGSI, err := parseGlobalSecondaryIndexes(tableCreationParams)
	if err != nil {
		return nil, err
	}
	if !validateIndexThroughputModeConsistency(billingMode, importGSI) {
		return nil, ErrInvalidParameter
	}
	// The model's TableCreationParameters carries no LocalSecondaryIndexes
	// member, and the developer guide states the limitation outright: "You
	// can create this table with global secondary indexes (local secondary
	// indexes are not supported)". A request naming local secondary indexes
	// is rejected rather than silently building an index the import
	// contract does not define.
	if _, present := tableCreationParams["LocalSecondaryIndexes"]; present {
		return nil, ErrInvalidParameter
	}

	// The SSE and on-demand maximum settings are TableCreationParameters
	// members the created table carries; the SSE member runs through the
	// table-creation SSE finaliser (disable policy and key resolution).
	sseDesc, sseDisable, sseErr := parseSSESpecification(tableCreationParams["SSESpecification"])
	if sseErr != nil {
		return nil, sseErr
	}
	sseDesc, err = s.finaliseSSESpecification(ctx, reqCtx, sseDesc, sseDisable, sseOpTableCreate)
	if err != nil {
		return nil, err
	}
	onDemandThroughput, odtErr := parseOnDemandThroughput(tableCreationParams)
	if odtErr != nil {
		return nil, odtErr
	}
	if billingMode == dbstore.BillingModeProvisioned && onDemandThroughput != nil {
		return nil, ErrInvalidParameter
	}

	job := ImportTableCoreInput{
		TableName:            tableName,
		KeySchema:            keySchema,
		AttributeDefs:        attrDefs,
		BillingMode:          billingMode,
		ProvThroughput:       provThroughput,
		GSI:                  importGSI,
		SSEDescription:       sseDesc,
		OnDemandThroughput:   onDemandThroughput,
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
				return &importTableResult{Import: replay}, nil
			}
			// The token record outlived the import itself (the import table
			// was deleted); the call creates a replacement import.
		}
	}

	// The import's target table does not exist yet; its UUID is minted at
	// creation inside the job and recorded on the import there.
	imp, err := store.Imports().Create(tableArn, "")
	if err != nil {
		return nil, err
	}
	imp.ClientToken = clientToken
	imp.InputCompressionType = dbstore.InputCompressionType(inputCompressionType)
	// The request has already fixed the input format and the bucket
	// source: they belong on the record at creation, so DescribeImport
	// answers both throughout IN_PROGRESS instead of from the job's first
	// write onwards.
	imp.InputFormat = dbstore.InputFormat(inputFormat)
	imp.S3BucketSource = &dbstore.S3BucketSource{
		S3Bucket:      s3Bucket,
		S3Prefix:      s3Prefix,
		S3BucketOwner: s3BucketOwner,
	}
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
	return &importTableResult{Import: imp}, nil
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

	importedCount := int64(0)
	processedCount := int64(0)
	errorCount := int64(0)
	processedSizeBytes := int64(0)

	// failImport persists the FAILED description with the figures the run
	// had already counted: the counters are observable state a client
	// polling DescribeImport reads off either terminal state, so a
	// mid-run failure must not discard them.
	failImport := func(code, message string) {
		imp.ImportStatus = "FAILED"
		imp.FailureCode = code
		imp.FailureMessage = message
		imp.ProcessedItemCount = processedCount
		imp.ImportedItemCount = importedCount
		imp.ErrorCount = errorCount
		imp.ProcessedSizeBytes = processedSizeBytes
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
	})
	if err != nil {
		// The code follows the store's error class: a duplicate-name race
		// (the synchronous preflight passed, another client created the
		// table inside the window) reports TableAlreadyExists; any other
		// Create failure is not the client's naming and reports the
		// internal code instead.
		code := importCreateFailureCode(err)
		failImport(code, fmt.Sprintf("failed to create target table: %v", err))
		return
	}
	imp.TableArn = table.ARN
	imp.TableId = table.TableId
	// The SSE specification and on-demand maximum settings the request's
	// TableCreationParameters carried are part of the created table; the
	// update runs under the table's record lock so any metric deltas the
	// import has already flushed survive.
	if in.SSEDescription != nil || in.OnDemandThroughput != nil {
		if _, err := store.Tables().Update(table.Name, func(t *dbstore.Table) error {
			if in.SSEDescription != nil {
				t.SSEDescription = in.SSEDescription
			}
			if in.OnDemandThroughput != nil {
				t.OnDemandThroughput = in.OnDemandThroughput
			}
			return nil
		}); err != nil {
			failImport("InternalFailure", fmt.Sprintf("failed to apply table settings: %v", err))
			return
		}
	}

	s3 := s.s3invoker()
	if s3 == nil {
		// The same ambiguity the missing-bucket branch below refuses: an
		// unreachable S3 plane cannot be distinguished from an empty
		// source, and a job that reported COMPLETED would fake success
		// over data read from nowhere.
		failImport("S3AccessDenied", "the S3 source plane is unavailable")
		return
	}
	// A missing source bucket is reported as a failed import rather
	// than an empty one, because listing cannot distinguish the two.
	exists, existsErr := s3.BucketExists(s.bgCtx, region, in.S3Bucket)
	if existsErr != nil {
		failImport(s3FailureCode(existsErr), fmt.Sprintf("failed to check source bucket: %v", existsErr))
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
		failImport(s3FailureCode(listErr), fmt.Sprintf("failed to list S3 objects: %v", listErr))
		return
	}
	for _, key := range keys {
		data, getErr := s3.GetObject(s.bgCtx, region, in.S3Bucket, key, 0)
		if getErr != nil {
			logs.Error("Failed to read S3 object for import",
				logs.Err(getErr),
				logs.String("bucket", in.S3Bucket),
				logs.String("key", key))
			failImport(s3FailureCode(getErr), fmt.Sprintf("failed to read S3 object %s: %v", key, getErr))
			return
		}
		if in.InputCompressionType != "NONE" {
			uncompressed, dcErr := decompressImportObject(data, in.InputCompressionType)
			if dcErr != nil {
				code := "GZIPError"
				if in.InputCompressionType == "ZSTD" {
					code = "ZSTDError"
				}
				failImport(code, fmt.Sprintf("failed to decompress S3 object %s: %v", key, dcErr))
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
		if len(record) > len(headers) {
			// A row carrying more values than the header list names cannot
			// be parsed into an item: the extra columns have no attribute
			// to bind to, and silently dropping them would import an item
			// the source file does not contain. The row counts as an error
			// item like any other unparseable line.
			counts.errors++
			continue
		}
		for i, val := range record {
			attrs[headers[i]] = &dbstore.AttributeValue{S: &val}
		}

		importLineItem(ctx, store, table, attrs, &counts)
	}
	return counts, nil
}

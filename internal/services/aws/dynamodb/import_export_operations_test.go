package dynamodb

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// newImportTestStore opens a real store with one table keyed by a string
// "id" attribute, matching the import fixtures below.
func newImportTestStore(t *testing.T, keyType dbstore.ScalarAttributeType) *dbstore.DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	store := dbstore.NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "ImpTbl",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: keyType}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return store
}

func TestImportDynamoDBJSONDataCountsDistinctItems(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)

	// The second i1 line overwrites the first: every line is processed and
	// imported, but the table ends holding two distinct items.
	data := []byte("{\"Item\":{\"id\":{\"S\":\"i1\"},\"v\":{\"S\":\"one\"}}}\n" +
		"{\"Item\":{\"id\":{\"S\":\"i2\"},\"v\":{\"S\":\"two\"}}}\n" +
		"{\"Item\":{\"id\":{\"S\":\"i1\"},\"v\":{\"S\":\"again\"}}}\n")
	counts, err := importDynamoDBJSONData(t.Context(), data, "ImpTbl", store)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if counts.processed != 3 || counts.imported != 3 || counts.errors != 0 {
		t.Fatalf("counts = %+v, want processed=3 imported=3 errors=0", counts)
	}

	table, err := store.Tables().Get("ImpTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	if table.ItemCount != 2 {
		t.Fatalf("table ItemCount = %d, want 2 distinct items (overwrites must not double-count)", table.ItemCount)
	}

	// Last write wins for the duplicated key.
	item, err := store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {S: ptrStr("i1")}})
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if item.Attributes["v"].S == nil || *item.Attributes["v"].S != "again" {
		t.Fatalf("overwritten item value = %v, want again", item.Attributes["v"])
	}
}

func TestImportDynamoDBJSONDataValidatesKeyTypes(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeN)

	data := []byte("{\"Item\":{\"id\":{\"N\":\"1\"},\"v\":{\"S\":\"ok\"}}}\n" +
		// A string value in the number-typed key attribute is a data
		// validation error: the item is skipped, not stored.
		"{\"Item\":{\"id\":{\"S\":\"x\"},\"v\":{\"S\":\"bad\"}}}\n" +
		// A line missing the key attribute parses into an item and then
		// fails key-completeness validation.
		"{\"Item\":{\"v\":{\"S\":\"keyless\"}}}\n" +
		"not-json-at-all\n")
	counts, err := importDynamoDBJSONData(t.Context(), data, "ImpTbl", store)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	// The keyless line parsed into an item, so it counts as processed and
	// then fails validation; only the unparsable line skips the processed
	// tally. Errors: wrong key type, missing key, unparsable line.
	if counts.processed != 3 || counts.imported != 1 || counts.errors != 3 {
		t.Fatalf("counts = %+v, want processed=3 imported=1 errors=3", counts)
	}

	table, err := store.Tables().Get("ImpTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	if table.ItemCount != 1 {
		t.Fatalf("table ItemCount = %d, want 1 (only the well-typed item is stored)", table.ItemCount)
	}

	item, err := store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {N: ptrStr("1")}})
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if item.Attributes["v"].S == nil || *item.Attributes["v"].S != "ok" {
		t.Fatalf("imported item value = %v, want ok", item.Attributes["v"])
	}
}

// The export and import ClientToken persists with its job record so a
// retry inside the idempotency window can find the original job, and the
// payload hash ignores the token itself while staying stable across map
// iteration orders.
func TestClientTokenIdempotencyHelpers(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)

	h1 := clientTokenHash(map[string]interface{}{"S3Bucket": "b", "ClientToken": "tok"})
	h2 := clientTokenHash(map[string]interface{}{"ClientToken": "tok", "S3Bucket": "b"})
	if h1 != h2 {
		t.Fatalf("hash not stable across map order: %s vs %s", h1, h2)
	}
	if clientTokenHash(map[string]interface{}{"S3Bucket": "other", "ClientToken": "tok"}) == h1 {
		t.Fatal("changed payload produced the same hash")
	}
	if clientTokenHash(map[string]interface{}{"S3Bucket": "b", "ClientToken": "different"}) != h1 {
		t.Fatal("token itself must not enter the payload hash")
	}

	tableArn := store.Tables().ARNBuilder().Table("ImpTbl")
	exp, err := store.Exports().Create(tableArn, "ImpTbl", "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("create export: %v", err)
	}
	exp.ClientToken = "tok-export"
	if err := store.Exports().Put(exp); err != nil {
		t.Fatalf("put export: %v", err)
	}
	if findExportByClientToken(store, tableArn, "tok-export") == nil {
		t.Fatal("export not found by its client token")
	}
	if findExportByClientToken(store, tableArn, "tok-other") != nil {
		t.Fatal("export matched a foreign client token")
	}

	imp, err := store.Imports().Create(tableArn, "ImpTbl")
	if err != nil {
		t.Fatalf("create import: %v", err)
	}
	imp.ClientToken = "tok-import"
	if err := store.Imports().Put(imp); err != nil {
		t.Fatalf("put import: %v", err)
	}
	if findImportByClientToken(store, tableArn, "tok-import") == nil {
		t.Fatal("import not found by its client token")
	}
}

// GZIP-imported source objects decompress before parsing, so a gzipped
// DYNAMODB_JSON source imports exactly the items it contains.
func TestImportGZIPCompressedData(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)

	raw := []byte("{\"Item\":{\"id\":{\"S\":\"g1\"},\"v\":{\"S\":\"one\"}}}\n" +
		"{\"Item\":{\"id\":{\"S\":\"g2\"},\"v\":{\"S\":\"two\"}}}\n")
	var compressed bytes.Buffer
	gzWriter := gzip.NewWriter(&compressed)
	if _, err := gzWriter.Write(raw); err != nil {
		t.Fatalf("gzip: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}

	uncompressed, err := decompressImportObject(compressed.Bytes(), "GZIP")
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	counts, err := importDynamoDBJSONData(t.Context(), uncompressed, "ImpTbl", store)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if counts.imported != 2 {
		t.Fatalf("imported %d, want 2", counts.imported)
	}
	if got, gErr := store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {S: ptrStr("g1")}}); gErr != nil || got == nil {
		t.Fatalf("g1 not imported: %v", gErr)
	}
}

// The export manifest summary renders the documented SSE request shape and
// a JSON null for an absent KMS key. The incremental-only members follow
// the documented incremental summary example: present on incremental
// exports and absent otherwise.
func TestExportManifestSummarySSEFields(t *testing.T) {
	if got := s3SseAlgorithmForManifest("", ""); got != "AES256" {
		t.Fatalf("default algorithm = %s, want AES256", got)
	}
	if got := s3SseAlgorithmForManifest("", "alias/foo"); got != "KMS" {
		t.Fatalf("kms algorithm = %s, want KMS", got)
	}
	if got := s3SseAlgorithmForManifest("KMS", ""); got != "KMS" {
		t.Fatalf("explicit kms algorithm = %s, want KMS", got)
	}
	if got := s3SseAlgorithmForManifest("AES256", "alias/foo"); got != "AES256" {
		t.Fatalf("named algorithm governs = %s, want AES256", got)
	}
	summary := exportManifestSummary{S3SseKmsKeyId: nilIfEmpty("")}
	encoded, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(encoded), `"s3SseKmsKeyId":null`) {
		t.Fatalf("absent key must render as null: %s", encoded)
	}
	if strings.Contains(string(encoded), "outputView") || strings.Contains(string(encoded), "exportType") {
		t.Fatalf("full-export summary carries incremental members: %s", encoded)
	}
	view := "NEW_AND_OLD_IMAGES"
	inc := exportManifestSummary{OutputView: &view, ExportType: "INCREMENTAL_EXPORT"}
	encoded, err = json.Marshal(inc)
	if err != nil {
		t.Fatalf("marshal incremental summary: %v", err)
	}
	if !strings.Contains(string(encoded), `"outputView":"NEW_AND_OLD_IMAGES"`) || !strings.Contains(string(encoded), `"exportType":"INCREMENTAL_EXPORT"`) {
		t.Fatalf("incremental summary members: %s", encoded)
	}
}

// TestJobsFailWhenTheS3PlaneIsUnavailable pins the nil-invoker contract:
// with no S3 invoker registered, both jobs record FAILED with the S3
// access failure code instead of falling through to COMPLETED — the
// fail-loud principle the import path already applies to a missing source
// bucket (an unreachable plane cannot be distinguished from an empty one).
func TestJobsFailWhenTheS3PlaneIsUnavailable(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)
	table, err := store.Tables().Get("ImpTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}

	// A bare service carries no event bus, so s3invoker() returns nil.
	svc := &DynamoDBService{}

	export, err := store.Exports().Create(table.ARN, "ImpTbl", "DYNAMODB_JSON")
	if err != nil {
		t.Fatalf("create export record: %v", err)
	}
	svc.runExportJob(store, ExportTableCoreInput{
		TableArn:     table.ARN,
		TableName:    "ImpTbl",
		ExportFormat: "DYNAMODB_JSON",
		S3Bucket:     "dest-bucket",
		Region:       "us-east-1",
		ExportTime:   time.Now(),
	}, export.ExportArn)
	finalExport, err := store.Exports().Get(export.ExportArn)
	if err != nil {
		t.Fatalf("reload export: %v", err)
	}
	if finalExport.ExportStatus != "FAILED" || finalExport.FailureCode != "S3AccessDenied" {
		t.Fatalf("export with nil invoker = %s/%s, want FAILED/S3AccessDenied",
			finalExport.ExportStatus, finalExport.FailureCode)
	}

	imp, err := store.Imports().Create(table.ARN, "ImpTbl")
	if err != nil {
		t.Fatalf("create import record: %v", err)
	}
	svc.runImportJob(store, "us-east-1", ImportTableCoreInput{
		TableName:     "ImportedNoS3",
		KeySchema:     []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefs: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:   dbstore.BillingModePayPerRequest,
		S3Bucket:      "source-bucket",
		InputFormat:   "DYNAMODB_JSON",
	}, imp.ImportArn)
	finalImport, err := store.Imports().Get(imp.ImportArn)
	if err != nil {
		t.Fatalf("reload import: %v", err)
	}
	if finalImport.ImportStatus != "FAILED" || finalImport.FailureCode != "S3AccessDenied" {
		t.Fatalf("import with nil invoker = %s/%s, want FAILED/S3AccessDenied",
			finalImport.ImportStatus, finalImport.FailureCode)
	}
}

// TestImportTableCreationParametersMemberSet pins the import creation
// member contract: a request naming LocalSecondaryIndexes is rejected (the
// model's TableCreationParameters has no such member and the developer
// guide states "local secondary indexes are not supported"), while the
// SSESpecification and OnDemandThroughput members the model does define are
// applied to the created table, and the on-demand pair is rejected on a
// provisioned table.
func TestImportTableCreationParametersMemberSet(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	base := func(table string, extra map[string]interface{}) map[string]interface{} {
		params := map[string]interface{}{
			"S3BucketSource": map[string]interface{}{"S3Bucket": "import-bucket"},
			"InputFormat":    "DYNAMODB_JSON",
			"TableCreationParameters": map[string]interface{}{
				"TableName":            table,
				"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
				"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
				"BillingMode":          "PAY_PER_REQUEST",
			},
		}
		for k, v := range extra {
			params[k] = v
		}
		return params
	}
	withCreation := func(table string, creationExtra map[string]interface{}) map[string]interface{} {
		params := base(table, nil)
		tcp := params["TableCreationParameters"].(map[string]interface{})
		for k, v := range creationExtra {
			tcp[k] = v
		}
		return params
	}

	// LocalSecondaryIndexes is not a member the import contract defines.
	lsiParams := withCreation("ImpLsiTbl", map[string]interface{}{
		"LocalSecondaryIndexes": []interface{}{map[string]interface{}{
			"IndexName":  "lsi-1",
			"KeySchema":  []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"Projection": map[string]interface{}{"ProjectionType": "ALL"},
		}},
	})
	if _, err := svc.ImportTable(ctx, reqCtx, &request.ParsedRequest{Parameters: lsiParams}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("import with LSI: err = %v, want ErrInvalidParameter", err)
	}

	// The on-demand pair is rejected on a provisioned table.
	provParams := withCreation("ImpProvOdt", map[string]interface{}{
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": 5.0, "WriteCapacityUnits": 5.0},
		"OnDemandThroughput":    map[string]interface{}{"MaxReadRequestUnits": 9.0},
	})
	if _, err := svc.ImportTable(ctx, reqCtx, &request.ParsedRequest{Parameters: provParams}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("provisioned import with on-demand pair: err = %v, want ErrInvalidParameter", err)
	}

	// SSE and the on-demand pair apply to the created table. The job then
	// fails at the S3 plane (no invoker on a bare service), after the table
	// exists — the failure leaves the settings observable.
	sseParams := withCreation("ImpSseTbl", map[string]interface{}{
		"SSESpecification":   map[string]interface{}{"Enabled": true},
		"OnDemandThroughput": map[string]interface{}{"MaxReadRequestUnits": 9.0, "MaxWriteRequestUnits": 8.0},
	})
	if _, err := svc.ImportTable(ctx, reqCtx, &request.ParsedRequest{Parameters: sseParams}); err != nil {
		t.Fatalf("import with sse: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		table, getErr := store.Tables().Get("ImpSseTbl")
		if getErr == nil {
			if table.SSEDescription == nil || table.SSEDescription.SSEType != "KMS" {
				t.Fatalf("imported table SSE = %+v, want KMS", table.SSEDescription)
			}
			if table.OnDemandThroughput == nil || table.OnDemandThroughput.MaxReadRequestUnits != 9 ||
				table.OnDemandThroughput.MaxWriteRequestUnits != 8 {
				t.Fatalf("imported table on-demand = %+v, want 9/8", table.OnDemandThroughput)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("imported table never appeared: %v", getErr)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestImportZSTDCompressedData pins the ZSTD half of the
// InputCompressionType contract: a modelled, documented enum value
// ("Data can be compressed in ZSTD or GZIP format") that decompresses
// through the already-vendored klauspost/compress zstd reader, with a
// corrupt stream failing under the format-naming failure code.
func TestImportZSTDCompressedData(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)

	var compressed bytes.Buffer
	zw, err := zstd.NewWriter(&compressed)
	if err != nil {
		t.Fatalf("zstd writer: %v", err)
	}
	if _, err := zw.Write([]byte(
		"{\"Item\":{\"id\":{\"S\":\"z1\"},\"v\":{\"S\":\"one\"}}}\n" +
			"{\"Item\":{\"id\":{\"S\":\"z2\"},\"v\":{\"S\":\"two\"}}}\n")); err != nil {
		t.Fatalf("zstd write: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zstd close: %v", err)
	}

	uncompressed, err := decompressImportObject(compressed.Bytes(), "ZSTD")
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	counts, err := importDynamoDBJSONData(t.Context(), uncompressed, "ImpTbl", store)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if counts.imported != 2 {
		t.Fatalf("imported %d, want 2", counts.imported)
	}
	if got, gErr := store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {S: ptrStr("z1")}}); gErr != nil || got == nil {
		t.Fatalf("z1 not imported: %v", gErr)
	}

	// Garbage does not decode as a zstd stream: the decompression failure
	// is surfaced, not swallowed as an empty import.
	if _, err := decompressImportObject([]byte("not a zstd stream"), "ZSTD"); err == nil {
		t.Fatal("garbage decoded as zstd without error")
	}

	// The request plane accepts the value and the job carries it: a bare
	// service has no S3 invoker, so a valid ZSTD request passes validation,
	// records the compression type from creation, and terminates FAILED at
	// the S3 plane — not a validation rejection.
	svc, reqCtx := billingModePlaneFixture(t)
	importResp, err := svc.ImportTable(t.Context(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"S3BucketSource":       map[string]interface{}{"S3Bucket": "import-bucket"},
		"InputFormat":          "DYNAMODB_JSON",
		"InputCompressionType": "ZSTD",
		"TableCreationParameters": map[string]interface{}{
			"TableName":            "ImpZstdTbl",
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		},
	}})
	if err != nil {
		t.Fatalf("zstd import request: %v", err)
	}
	desc := importResp.(map[string]interface{})["ImportTableDescription"].(map[string]interface{})
	if desc["InputCompressionType"] != "ZSTD" {
		t.Fatalf("request-plane compression type = %v, want ZSTD", desc["InputCompressionType"])
	}
	importStore, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		final, getErr := importStore.Imports().Get(desc["ImportArn"].(string))
		if getErr != nil {
			t.Fatalf("reload import: %v", getErr)
		}
		if final.ImportStatus != "IN_PROGRESS" {
			if final.ImportStatus != "FAILED" || final.FailureCode != "S3AccessDenied" {
				t.Fatalf("zstd import job = %s/%s, want FAILED/S3AccessDenied", final.ImportStatus, final.FailureCode)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("zstd import job never left IN_PROGRESS")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestImportCSVOptionsValidation pins ImportTable's CSV options handling
// at the request boundary — the options travel under InputFormatOptions.Csv
// on the wire: a valid delimiter and header list pass, while a
// multi-character delimiter, a non-string header element and an empty
// header list are validation errors.
func TestImportCSVOptionsValidation(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	base := func(csv map[string]interface{}) importTableInput {
		return importTableInput{Parameters: map[string]interface{}{
			"InputFormat":        "CSV",
			"S3BucketSource":     map[string]interface{}{"S3Bucket": "csv-pin-bucket"},
			"InputFormatOptions": map[string]interface{}{"Csv": csv},
			"TableCreationParameters": map[string]interface{}{
				"TableName":            "CsvPinTable",
				"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "id", "KeyType": "HASH"}},
				"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "id", "AttributeType": "S"}},
				"BillingMode":          "PAY_PER_REQUEST",
			},
		}}
	}

	if _, err := svc.importTableCore(ctx, reqCtx, base(map[string]interface{}{
		"Delimiter":  ";",
		"HeaderList": []interface{}{"id", "v"},
	})); err != nil {
		t.Fatalf("valid csv options: %v", err)
	}
	if _, err := svc.importTableCore(ctx, reqCtx, base(map[string]interface{}{
		"Delimiter": ";;",
	})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("multi-character delimiter: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := svc.importTableCore(ctx, reqCtx, base(map[string]interface{}{
		"HeaderList": []interface{}{"id", 1.0},
	})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("non-string header: err = %v, want ErrInvalidParameter", err)
	}
	if _, err := svc.importTableCore(ctx, reqCtx, base(map[string]interface{}{
		"HeaderList": []interface{}{},
	})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty header list: err = %v, want ErrInvalidParameter", err)
	}
}

// TestImportCSVDataDeliveredOptions pins the data-plane half of the CSV
// options: the delimiter the request carried splits the rows, the header
// list names the columns without consuming a data row, and the default
// path takes the header from the first row.
func TestImportCSVDataDeliveredOptions(t *testing.T) {
	store := newImportTestStore(t, dbstore.ScalarAttributeTypeS)

	counts, err := importCSVData(t.Context(), []byte("1;one\n2;two\n"), "ImpTbl", store, ";", []string{"id", "v"})
	if err != nil {
		t.Fatalf("csv import with delivered options: %v", err)
	}
	if counts.processed != 2 || counts.imported != 2 || counts.errors != 0 {
		t.Fatalf("counts = %+v, want processed=2 imported=2 errors=0", counts)
	}
	item, err := store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {S: ptrStr("1")}})
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if item.Attributes["v"].S == nil || *item.Attributes["v"].S != "one" {
		t.Fatalf("delivered header list: item = %v, want v=one", item.Attributes)
	}

	counts, err = importCSVData(t.Context(), []byte("id,v\n3,three\n"), "ImpTbl", store, "", nil)
	if err != nil {
		t.Fatalf("csv import with default options: %v", err)
	}
	if counts.processed != 1 || counts.imported != 1 || counts.errors != 0 {
		t.Fatalf("counts = %+v, want processed=1 imported=1 errors=0", counts)
	}
	item, err = store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {S: ptrStr("3")}})
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if item.Attributes["v"].S == nil || *item.Attributes["v"].S != "three" {
		t.Fatalf("first-row header: item = %v, want v=three", item.Attributes)
	}

	// A row longer than the header list cannot bind its extra columns and
	// counts as an error item — never a silently truncated import.
	counts, err = importCSVData(t.Context(), []byte("id,v\n4,four,orphan\n5,five\n"), "ImpTbl", store, "", nil)
	if err != nil {
		t.Fatalf("csv import with a long row: %v", err)
	}
	if counts.processed != 1 || counts.imported != 1 || counts.errors != 1 {
		t.Fatalf("counts = %+v, want processed=1 imported=1 errors=1 (the long row)", counts)
	}
	if _, err := store.Items().Get("ImpTbl", map[string]*dbstore.AttributeValue{"id": {S: ptrStr("4")}}); err == nil {
		t.Fatal("the over-long row must not import a truncated item")
	}
}

// TestListImportsPaginationFollowsIssuedTokens pins the ListImports
// pagination contract end to end: every issued NextToken satisfies the
// model's hexadecimal next-token pattern, the page after a token carries
// the remaining imports, and a marker that is not a token the issuer
// could have produced — a raw import ARN, the record key the store marks
// by — is a validation error rather than a silently wrong page.
func TestListImportsPaginationFollowsIssuedTokens(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	tableArn := store.Tables().ARNBuilder().Table("TokTable")
	for i := 0; i < 3; i++ {
		if _, err := store.Imports().Create(tableArn, ""); err != nil {
			t.Fatalf("seed import %d: %v", i, err)
		}
	}

	listPage := func(token string) ([]*dbstore.ImportTableDescription, string, error) {
		params := map[string]interface{}{"PageSize": 1.0, "TableArn": tableArn}
		if token != "" {
			params["NextToken"] = token
		}
		return svc.listImportsCore(ctx, reqCtx, listImportsInput{Parameters: params})
	}

	seen := make(map[string]bool)
	token := ""
	for range 4 {
		page, next, err := listPage(token)
		if err != nil {
			t.Fatalf("list page after %q: %v", token, err)
		}
		for _, imp := range page {
			seen[imp.ImportArn] = true
		}
		if next == "" {
			break
		}
		if !validateImportNextToken(next) {
			t.Fatalf("issued token %q fails the model's next-token pattern", next)
		}
		token = next
	}
	if len(seen) != 3 {
		t.Fatalf("pagination saw %d distinct imports, want 3", len(seen))
	}

	if _, _, err := listPage(tableArn); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("raw ARN as next token: expected ErrInvalidParameter, got %v", err)
	}
}

// enableLegacyPitr turns point-in-time recovery on for the legacy fixture's
// table: the journal only records mutations while recovery is enabled, so
// every export-window test calls this before its first write.
func enableLegacyPitr(t *testing.T, svc *DynamoDBService, reqCtx *request.RequestContext) {
	t.Helper()
	if _, err := svc.UpdateContinuousBackups(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "LegacyTable",
		"PointInTimeRecoverySpecification": map[string]interface{}{
			"PointInTimeRecoveryEnabled": true,
		},
	}}); err != nil {
		t.Fatalf("enable pitr: %v", err)
	}
}

// TestExportTypeAndIncrementalSpecificationValidation pins the ExportType
// contract and the incremental specification's acceptance rules: the enum
// validation, the pairing rules (an incremental export requires the
// specification, a full export rejects it), the window checks against the
// restorable range, and the documented defaults — FULL_EXPORT for an
// omitted type, NEW_AND_OLD_IMAGES for an omitted view. The accepted
// incremental export then runs its job on a bare service and terminates
// FAILED at the S3 plane, which also pins the incremental branch of the
// job.
func TestExportTypeAndIncrementalSpecificationValidation(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("LegacyTable")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	enableLegacyPitr(t, svc, reqCtx)
	// The window's start must postdate the restorable range's beginning,
	// which the enablement stamps on the record: derive the wait and the
	// accepted window's start from that stamp rather than a fixed sleep —
	// the poll exits as soon as the stamp's age reaches the margin, which
	// the enablement round-trip has normally already provided.
	pitrTable, err := store.Tables().Get("LegacyTable")
	if err != nil {
		t.Fatalf("reload table for pitr stamp: %v", err)
	}
	enabled := pitrTable.PointInTimeRecovery.EarliestRestorableDateTime
	const stampMargin = 10 * time.Millisecond
	for time.Since(enabled) < stampMargin {
		time.Sleep(time.Millisecond)
	}
	now := time.Now()

	epoch := func(ts time.Time) float64 { return float64(ts.UnixMilli()) / 1000.0 }
	base := func(extra map[string]interface{}) map[string]interface{} {
		params := map[string]interface{}{
			"TableArn": table.ARN,
			"S3Bucket": "dest-bucket",
		}
		for k, v := range extra {
			params[k] = v
		}
		return params
	}
	incremental := func(spec map[string]interface{}) map[string]interface{} {
		return base(map[string]interface{}{"ExportType": "INCREMENTAL_EXPORT", "IncrementalExportSpecification": spec})
	}
	reject := func(params map[string]interface{}, want error) {
		t.Helper()
		if _, err := svc.exportTableCore(context.Background(), reqCtx, exportTableInput{Parameters: params}); !errors.Is(err, want) {
			t.Fatalf("exportTableCore: err = %v, want %v", err, want)
		}
	}

	// Enum validation.
	reject(base(map[string]interface{}{"S3SseAlgorithm": "DES3"}), ErrInvalidParameter)
	reject(base(map[string]interface{}{"ExportType": "PARTIAL"}), ErrInvalidParameter)
	// The S3Bucket pattern applies unconditionally: the empty string — a
	// required member's omission — is rejected at request time, never left
	// to fail inside the background job.
	reject(base(map[string]interface{}{"S3Bucket": ""}), ErrInvalidParameter)
	// Pairing rules: an incremental export without a specification, and a
	// specification on a full export.
	reject(base(map[string]interface{}{"ExportType": "INCREMENTAL_EXPORT"}), ErrInvalidParameter)
	reject(base(map[string]interface{}{"ExportType": "FULL_EXPORT", "IncrementalExportSpecification": map[string]interface{}{}}), ErrInvalidParameter)
	// Specification shape.
	reject(incremental(map[string]interface{}{"ExportViewType": "OLD_IMAGE", "ExportFromTime": epoch(now)}), ErrInvalidParameter)
	reject(base(map[string]interface{}{"ExportType": "INCREMENTAL_EXPORT", "IncrementalExportSpecification": "bogus"}), ErrInvalidParameter)
	// Window checks: a start before the restorable range, an end before
	// the start, an end in the future.
	reject(incremental(map[string]interface{}{"ExportFromTime": epoch(now.Add(-400 * 24 * time.Hour))}), ErrInvalidExportTime)
	reject(incremental(map[string]interface{}{"ExportFromTime": epoch(now), "ExportToTime": epoch(now.Add(-time.Hour))}), ErrInvalidExportTime)
	reject(incremental(map[string]interface{}{"ExportFromTime": epoch(now), "ExportToTime": epoch(now.Add(time.Hour))}), ErrInvalidExportTime)

	// The accepted incremental export echoes the resolved specification,
	// with the documented view default applied. The window's start sits
	// inside the restorable range by construction: after the enablement
	// stamp, before the polled now.
	from := enabled.Add(stampMargin / 2)
	validParams := incremental(map[string]interface{}{"ExportFromTime": epoch(from)})
	validParams["S3SseAlgorithm"] = "KMS"
	outcome, err := svc.exportTableCore(context.Background(), reqCtx, exportTableInput{Parameters: validParams})
	if err != nil {
		t.Fatalf("incremental export: %v", err)
	}
	if outcome.Export.ExportType != "INCREMENTAL_EXPORT" {
		t.Fatalf("export type = %s, want INCREMENTAL_EXPORT", outcome.Export.ExportType)
	}
	if outcome.Export.ExportViewType != "NEW_AND_OLD_IMAGES" {
		t.Fatalf("view type = %s, want NEW_AND_OLD_IMAGES", outcome.Export.ExportViewType)
	}
	if outcome.Export.S3SseAlgorithm != "KMS" {
		t.Fatalf("sse algorithm = %s, want KMS", outcome.Export.S3SseAlgorithm)
	}
	if diff := outcome.Export.ExportFromTime.Sub(from); diff < -2*time.Millisecond || diff > 2*time.Millisecond {
		t.Fatalf("from-time echo = %v, want ~%v", outcome.Export.ExportFromTime, from)
	}
	// The omitted end resolves to the present at request time.
	if end := outcome.Export.ExportToTime; end.Before(from) || end.After(time.Now()) {
		t.Fatalf("to-time echo = %v, want within [%v, now]", end, from)
	}

	// The background job runs against a bare service (no S3 plane) and
	// must terminate FAILED rather than fall through to COMPLETED.
	deadline := time.Now().Add(5 * time.Second)
	for {
		final, getErr := store.Exports().Get(outcome.Export.ExportArn)
		if getErr != nil {
			t.Fatalf("reload export: %v", getErr)
		}
		if final.ExportStatus != "IN_PROGRESS" {
			if final.ExportStatus != "FAILED" || final.FailureCode != "S3AccessDenied" {
				t.Fatalf("incremental job with nil invoker = %s/%s, want FAILED/S3AccessDenied",
					final.ExportStatus, final.FailureCode)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("incremental export job never left IN_PROGRESS")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestIncrementalExportRecordShapes pins the incremental export's record
// construction against the documented output-view table: an insert record
// carries the new image, an update record the new image plus the old one
// under NEW_AND_OLD_IMAGES, a delete record the keys plus the old image,
// and an insert-then-delete pair inside the window produces no output.
// Records keep journal order and carry the last write's microsecond
// timestamp; the DynamoDB JSON and Ion serialisers render the documented
// member names.
func TestIncrementalExportRecordShapes(t *testing.T) {
	svc, reqCtx := newLegacyTestService(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	enableLegacyPitr(t, svc, reqCtx)

	put := func(id, v string) {
		t.Helper()
		params := map[string]interface{}{
			"TableName": "LegacyTable",
			"Item": map[string]interface{}{
				"id": map[string]interface{}{"S": id},
				"sk": map[string]interface{}{"S": "a"},
				"v":  map[string]interface{}{"S": v},
			},
		}
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params}); err != nil {
			t.Fatalf("put %s: %v", id, err)
		}
	}
	remove := func(id string) {
		t.Helper()
		params := map[string]interface{}{
			"TableName": "LegacyTable",
			"Key": map[string]interface{}{
				"id": map[string]interface{}{"S": id},
				"sk": map[string]interface{}{"S": "a"},
			},
		}
		if _, err := svc.DeleteItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params}); err != nil {
			t.Fatalf("delete %s: %v", id, err)
		}
	}

	// Pre-window state supplies the update's old image and the delete's
	// only image; keep stays untouched throughout and tmp is written and
	// removed inside the window, so neither may appear in the output.
	put("upd", "v1")
	put("del", "v1")
	put("keep", "v1")
	from := time.Now()
	put("ins", "ins")
	put("upd", "v2")
	remove("del")
	put("tmp", "tmp")
	remove("tmp")
	to := time.Now()

	records, err := buildIncrementalExportRecords(store, "LegacyTable", from, to, "NEW_AND_OLD_IMAGES")
	if err != nil {
		t.Fatalf("build records: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("got %d records, want 3 (ins, upd, del; the insert+delete pair emits nothing)", len(records))
	}
	ins, upd, del := records[0], records[1], records[2]
	for _, record := range records {
		if record.Keys["id"] == nil || record.Keys["sk"] == nil {
			t.Fatalf("record keys = %v, want id and sk", record.Keys)
		}
		if record.WriteTimestampMicros < from.UnixMicro() || record.WriteTimestampMicros > to.UnixMicro() {
			t.Fatalf("write timestamp %d outside [%d, %d]",
				record.WriteTimestampMicros, from.UnixMicro(), to.UnixMicro())
		}
	}
	if got := ins.NewImage["v"].S; got == nil || *got != "ins" {
		t.Fatalf("insert new image v = %v, want ins", ins.NewImage["v"])
	}
	if ins.OldImage != nil {
		t.Fatalf("insert carries an old image: %v", ins.OldImage)
	}
	if got := upd.NewImage["v"].S; got == nil || *got != "v2" {
		t.Fatalf("update new image v = %v, want v2", upd.NewImage["v"])
	}
	if got := upd.OldImage["v"].S; got == nil || *got != "v1" {
		t.Fatalf("update old image v = %v, want v1", upd.OldImage["v"])
	}
	if del.NewImage != nil {
		t.Fatalf("delete carries a new image: %v", del.NewImage)
	}
	if got := del.OldImage["v"].S; got == nil || *got != "v1" {
		t.Fatalf("delete old image v = %v, want v1", del.OldImage["v"])
	}

	// The new-images-only view keeps every record but drops the old
	// images: an update keeps its new image, a delete reduces to the bare
	// keys.
	newOnly, err := buildIncrementalExportRecords(store, "LegacyTable", from, to, "NEW_IMAGE")
	if err != nil {
		t.Fatalf("build new-image records: %v", err)
	}
	if len(newOnly) != 3 {
		t.Fatalf("new-image view: got %d records, want 3", len(newOnly))
	}
	if newOnly[1].OldImage != nil || newOnly[2].OldImage != nil {
		t.Fatalf("new-image view carries old images: %+v", newOnly[1:])
	}
	if newOnly[2].Keys["id"] == nil {
		t.Fatalf("new-image delete lost its keys: %v", newOnly[2].Keys)
	}

	// The serialisers render the documented member names.
	encoded, err := json.Marshal(buildIncrementalWireRecord(upd))
	if err != nil {
		t.Fatalf("marshal wire record: %v", err)
	}
	stamp := strconv.FormatInt(upd.WriteTimestampMicros, 10)
	for _, want := range []string{
		`"Metadata":{"WriteTimestampMicros":{"N":"` + stamp + `"}}`,
		`"Keys":`,
		`"NewImage":`,
		`"OldImage":`,
	} {
		if !strings.Contains(string(encoded), want) {
			t.Fatalf("wire record missing %s: %s", want, encoded)
		}
	}
	ionLine := string(buildIonIncrementalRecordLine(upd))
	for _, want := range []string{"Metadata", "WriteTimestampMicros", `"N"`, "Keys", "NewImage", "OldImage"} {
		if !strings.Contains(ionLine, want) {
			t.Fatalf("ion line missing %s: %s", want, ionLine)
		}
	}
}

// The creation-time record carries the request's input format and bucket
// source immediately: both members are fixed by the request, so the
// persisted record answers them from creation — DescribeImport reports
// them throughout IN_PROGRESS, not from the background job's first write
// onwards.
func TestImportRecordCarriesRequestMembersAtCreation(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}

	outcome, err := svc.importTableCore(context.Background(), reqCtx, importTableInput{Parameters: map[string]interface{}{
		"S3BucketSource": map[string]interface{}{"S3Bucket": "src-bucket", "S3KeyPrefix": "pfx"},
		"InputFormat":    "DYNAMODB_JSON",
		"TableCreationParameters": map[string]interface{}{
			"TableName":            "ImpEchoTbl",
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		},
	}})
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	persisted, err := store.Imports().Get(outcome.Import.ImportArn)
	if err != nil {
		t.Fatalf("reload import: %v", err)
	}
	if persisted.InputFormat != "DYNAMODB_JSON" {
		t.Fatalf("persisted input format = %q, want DYNAMODB_JSON from creation", persisted.InputFormat)
	}
	if persisted.S3BucketSource == nil || persisted.S3BucketSource.S3Bucket != "src-bucket" || persisted.S3BucketSource.S3Prefix != "pfx" {
		t.Fatalf("persisted bucket source = %+v, want the request's src-bucket/pfx from creation", persisted.S3BucketSource)
	}
}

// TestImportRequiresInputFormat pins the required InputFormat member at
// the core: an omitted member is a ValidationException, never a silently
// applied DYNAMODB_JSON default — the SDK's client-side validation makes
// the omission unreachable through the wire, so the contract is pinned
// here.
func TestImportRequiresInputFormat(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)

	base := func(extra map[string]interface{}) importTableInput {
		params := map[string]interface{}{
			"S3BucketSource": map[string]interface{}{"S3Bucket": "src-bucket", "S3KeyPrefix": "pfx"},
			"TableCreationParameters": map[string]interface{}{
				"TableName":            "ImpFmtTbl",
				"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
				"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
				"BillingMode":          "PAY_PER_REQUEST",
			},
		}
		for k, v := range extra {
			params[k] = v
		}
		return importTableInput{Parameters: params}
	}

	if _, err := svc.importTableCore(context.Background(), reqCtx, base(nil)); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("omitted InputFormat: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.importTableCore(context.Background(), reqCtx, base(map[string]interface{}{"InputFormat": ""})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty InputFormat: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.importTableCore(context.Background(), reqCtx, base(map[string]interface{}{"InputFormat": float64(1)})); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("non-string InputFormat: expected ErrInvalidParameter, got %v", err)
	}
}

// TestImportDescriptionFacesAgree pins the one-builder contract of the
// import description: DescribeImport renders exactly the shared builder's
// members for the same record — TableId once the target table exists,
// ClientToken only when the record carries one — so the ImportTable and
// DescribeImport faces cannot disagree.
func TestImportDescriptionFacesAgree(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "FacesSrcTbl",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create source table: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("FacesSrcTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}

	// An empty source completes the job synchronously with its target
	// table created, so the record carries the TableId the model defines.
	svc.SetEventBus(ionExportBus{s3: &importSourceInvoker{}})
	imp, err := store.Imports().Create(table.ARN, table.TableId)
	if err != nil {
		t.Fatalf("create import record: %v", err)
	}
	svc.runImportJob(store, "us-east-1", ImportTableCoreInput{
		TableName:     "FacesTarget",
		KeySchema:     table.KeySchema,
		AttributeDefs: table.AttributeDefinitions,
		BillingMode:   dbstore.BillingModePayPerRequest,
		S3Bucket:      "src-bucket",
		InputFormat:   "DYNAMODB_JSON",
	}, imp.ImportArn)

	describe := func() map[string]interface{} {
		t.Helper()
		resp, err := svc.DescribeImport(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"ImportArn": imp.ImportArn,
		}})
		if err != nil {
			t.Fatalf("describe import: %v", err)
		}
		return resp.(map[string]interface{})["ImportTableDescription"].(map[string]interface{})
	}

	final, err := store.Imports().Get(imp.ImportArn)
	if err != nil {
		t.Fatalf("reload import: %v", err)
	}
	desc := describe()
	if !reflect.DeepEqual(desc, buildImportTableDescription(final)) {
		t.Fatalf("DescribeImport face disagrees with the shared builder:\n got %#v\nwant %#v",
			desc, buildImportTableDescription(final))
	}
	if _, has := desc["TableId"]; !has {
		t.Fatal("TableId must render once the target table exists — the model defines the member")
	}
	if _, has := desc["ClientToken"]; has {
		t.Fatal("an unset ClientToken must not render — never an empty-string member")
	}

	// A record carrying a client token renders it through the same
	// builder on the same face.
	final.ClientToken = "faces-token"
	if err := store.Imports().Put(final); err != nil {
		t.Fatalf("put tokened import: %v", err)
	}
	desc = describe()
	if desc["ClientToken"] != "faces-token" {
		t.Fatalf("set ClientToken = %v, want faces-token", desc["ClientToken"])
	}
}

package dynamodb

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"strings"
	"testing"

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

	store := dbstore.NewDynamoDBStore(st, "123456789012", "us-east-1")
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

	uncompressed, err := decompressImportObject(compressed.Bytes())
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
// a JSON null for an absent KMS key.
func TestExportManifestSummarySSEFields(t *testing.T) {
	if got := s3SseAlgorithmForManifest(""); got != "AES256" {
		t.Fatalf("default algorithm = %s, want AES256", got)
	}
	if got := s3SseAlgorithmForManifest("alias/foo"); got != "KMS" {
		t.Fatalf("kms algorithm = %s, want KMS", got)
	}
	summary := exportManifestSummary{S3SseKmsKeyId: nilIfEmpty("")}
	encoded, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(encoded), `"s3SseKmsKeyId":null`) {
		t.Fatalf("absent key must render as null: %s", encoded)
	}
}

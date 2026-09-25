// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// recordingS3Invoker captures every PutObject so an export job's S3 output
// can be inspected. The read-side methods are unreachable on the export
// path and return errors if called.
type recordingS3Invoker struct {
	puts map[string][]byte
}

func newRecordingS3Invoker() *recordingS3Invoker {
	return &recordingS3Invoker{puts: map[string][]byte{}}
}

func (r *recordingS3Invoker) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	return nil, errors.New("unexpected GetObject on export path")
}

func (r *recordingS3Invoker) GetObjectVersion(ctx context.Context, region, bucket, key, versionID string, maxBytes int64) ([]byte, error) {
	return nil, errors.New("unexpected GetObjectVersion on export path")
}

func (r *recordingS3Invoker) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	r.puts[key] = data
	return nil
}

func (r *recordingS3Invoker) PutObjectWithMetadata(ctx context.Context, region, bucket, key string, data []byte, contentType string, metadata map[string]string) error {
	r.puts[key] = data
	return nil
}

func (r *recordingS3Invoker) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	return nil, errors.New("unexpected ListObjects on export path")
}

func (r *recordingS3Invoker) ListObjectEntries(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]invokers.S3ObjectEntry, error) {
	return nil, errors.New("unexpected ListObjectEntries on export path")
}

func (r *recordingS3Invoker) BucketExists(ctx context.Context, region, bucket string) (bool, error) {
	return true, nil
}

func (r *recordingS3Invoker) GetBucketPolicy(ctx context.Context, region, bucket string) (string, error) {
	return "", nil
}

func (r *recordingS3Invoker) EnsureBucket(ctx context.Context, region, bucket string) error {
	return nil
}

func (r *recordingS3Invoker) DeleteObject(ctx context.Context, region, bucket, key string) error {
	return nil
}

// ionExportBus wires the recording invoker into the service's invoker
// lookup; the embedded nil interfaces cover the bus and registry members
// the export path never touches.
type ionExportBus struct {
	eventbus.Bus
	invokers.Registry
	s3 invokers.S3Invoker
}

func (b ionExportBus) S3Invoker() invokers.S3Invoker { return b.s3 }

// ionTestItem builds one item exercising every documented Ion mapping.
func ionTestItem() *dbstore.Item {
	str := func(s string) *string { return &s }
	tr := true
	return &dbstore.Item{
		TableName: "IonTbl",
		Key:       map[string]*dbstore.AttributeValue{"pk": {S: str("i1")}},
		Attributes: map[string]*dbstore.AttributeValue{
			"pk":     {S: str("i1")},
			"num":    {N: str("103")},
			"price":  {N: str("1.5")},
			"big":    {N: str("2E3")},
			"flag":   {BOOL: &tr},
			"off":    {BOOL: new(bool)},
			"nullv":  {NULL: &tr},
			"blob":   {B: []byte("hello")},
			"tags":   {SS: []string{"a", "b"}},
			"nums":   {NS: []string{"103", "6"}},
			"blobs":  {BS: [][]byte{[]byte("up")}},
			"nested": {L: []*dbstore.AttributeValue{{S: str("x")}, {M: map[string]*dbstore.AttributeValue{"inner": {S: str("y")}}}}},
			"m":      {M: map[string]*dbstore.AttributeValue{"k": {S: str("v")}}},
		},
	}
}

// TestBuildIonItemLine pins the documented Ion export line: the Ion
// version marker prefix, the Item wrapper, and the datatype mapping
// (S→string, N→decimal with the int→decimal dot and the d-exponent, BOOL,
// blob, sets with their $dynamodb_* annotations, list, struct, null), with
// sorted field names for determinism.
func TestBuildIonItemLine(t *testing.T) {
	line := buildIonItemLine(ionTestItem())
	want := `$ion_1_0 {"Item":{"big":2d3,"blob":{{aGVsbG8=}},"blobs":$dynamodb_BS::[{{dXA=}}],"flag":true,"m":{"k":"v"},"nested":["x",{"inner":"y"}],"nullv":null,"num":103.,"nums":$dynamodb_NS::[103.,6.],"off":false,"pk":"i1","price":1.5,"tags":$dynamodb_SS::["a","b"]}}`
	if string(line) != want {
		t.Fatalf("ion line:\n got %s\nwant %s", line, want)
	}
}

// TestExportJobWritesIonFormat pins the job-level contract: an ION export
// writes gzip-compressed Ion text lines under a .ion.gz data key, and the
// manifests describe the same format the data carries — the label and the
// bytes finally agree.
func TestExportJobWritesIonFormat(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "IonTbl",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := svc.PutItem(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "IonTbl",
		"Item": map[string]interface{}{
			"pk":   map[string]interface{}{"S": "i1"},
			"num":  map[string]interface{}{"N": "103"},
			"tags": map[string]interface{}{"SS": []interface{}{"a", "b"}},
		},
	}}); err != nil {
		t.Fatalf("put item: %v", err)
	}

	recorder := newRecordingS3Invoker()
	svc.SetEventBus(ionExportBus{s3: recorder})

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	table, err := store.Tables().Get("IonTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	export, err := store.Exports().Create(table.ARN, "IonTbl", "ION")
	if err != nil {
		t.Fatalf("create export record: %v", err)
	}

	svc.runExportJob(store, ExportTableCoreInput{
		TableArn:     table.ARN,
		TableName:    "IonTbl",
		ExportFormat: "ION",
		S3Bucket:     "dest-bucket",
		Region:       "us-east-1",
		ExportTime:   time.Now(),
	}, export.ExportArn)

	final, err := store.Exports().Get(export.ExportArn)
	if err != nil {
		t.Fatalf("reload export: %v", err)
	}
	if final.ExportStatus != "COMPLETED" || final.ExportFormat != "ION" {
		t.Fatalf("export = %s/%s, want COMPLETED/ION", final.ExportStatus, final.ExportFormat)
	}

	var dataKey string
	for key := range recorder.puts {
		// The documented export layout fixes the folder: the data object's
		// key is AWSDynamoDB/<export id>/data/<unique name>.
		if strings.HasPrefix(key, "AWSDynamoDB/") && strings.Contains(key, "/data/") {
			dataKey = key
		}
	}
	if dataKey == "" {
		t.Fatalf("no data object written: %v", recorder.puts)
	}
	if !strings.HasSuffix(dataKey, ".ion.gz") {
		t.Fatalf("data key %s lacks the .ion.gz suffix", dataKey)
	}
	zr, err := gzip.NewReader(bytes.NewReader(recorder.puts[dataKey]))
	if err != nil {
		t.Fatalf("data object is not gzip: %v", err)
	}
	data, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("gunzip: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) != 1 || !strings.HasPrefix(lines[0], "$ion_1_0 {\"Item\":") {
		t.Fatalf("data is not Ion lines: %q", data)
	}
	if !strings.Contains(lines[0], `"num":103.`) || !strings.Contains(lines[0], `"tags":$dynamodb_SS::["a","b"]`) {
		t.Fatalf("Ion mapping missing from line: %q", lines[0])
	}

	var manifestFilesKey, manifestSummaryKey string
	for key := range recorder.puts {
		if strings.HasSuffix(key, "/manifest-files.json") {
			manifestFilesKey = key
		}
		if strings.HasSuffix(key, "/manifest-summary.json") {
			manifestSummaryKey = key
		}
	}
	if manifestFilesKey == "" || manifestSummaryKey == "" {
		t.Fatalf("manifests not written: %v", recorder.puts)
	}
	if !strings.Contains(string(recorder.puts[manifestFilesKey]), dataKey) {
		t.Fatalf("files manifest does not name the data key %s: %s", dataKey, recorder.puts[manifestFilesKey])
	}
	// The etag follows the documented multipart form: the md5 hex with the
	// part-count suffix, not the bare digest.
	var manifestEntry exportManifestFileEntry
	if err := json.Unmarshal(recorder.puts[manifestFilesKey], &manifestEntry); err != nil {
		t.Fatalf("unmarshal files manifest: %v", err)
	}
	if len(manifestEntry.Etag) != 34 || !strings.HasSuffix(manifestEntry.Etag, "-1") {
		t.Fatalf("etag = %q, want the md5 hex with the -1 part-count suffix", manifestEntry.Etag)
	}
	if _, err := hex.DecodeString(manifestEntry.Etag[:32]); err != nil {
		t.Fatalf("etag prefix is not hex: %q", manifestEntry.Etag)
	}
	if !strings.Contains(string(recorder.puts[manifestSummaryKey]), `"outputFormat":"ION"`) {
		t.Fatalf("summary manifest labels the wrong format: %s", recorder.puts[manifestSummaryKey])
	}
}

// failingS3Invoker fails one chosen operation with a caller-supplied
// error, so a job's failure code can be pinned to the error class.
type failingS3Invoker struct {
	recordingS3Invoker
	failPut    error
	failGet    error
	failList   error
	failExists error
}

func (f *failingS3Invoker) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	if f.failPut != nil {
		return f.failPut
	}
	return f.recordingS3Invoker.PutObject(ctx, region, bucket, key, data, contentType)
}

func (f *failingS3Invoker) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	if f.failGet != nil {
		return nil, f.failGet
	}
	return f.recordingS3Invoker.GetObject(ctx, region, bucket, key, maxBytes)
}

func (f *failingS3Invoker) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	if f.failList != nil {
		return nil, f.failList
	}
	return f.recordingS3Invoker.ListObjects(ctx, region, bucket, prefix, maxKeys)
}

func (f *failingS3Invoker) BucketExists(ctx context.Context, region, bucket string) (bool, error) {
	if f.failExists != nil {
		return false, f.failExists
	}
	return f.recordingS3Invoker.BucketExists(ctx, region, bucket)
}

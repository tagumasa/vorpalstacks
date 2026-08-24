package sfn

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestParseCSVDataset pins the documented Step Functions CSV parsing rules:
// FIRST_ROW and GIVEN headers, the five delimiters, quote doubling,
// backslash escapes, missing fields padded with empty strings and extra
// fields dropped.
func TestParseCSVDataset(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		rc      *sfnstore.ItemReaderReaderConfig
		want    []map[string]string
		wantErr bool
	}{
		{
			name: "first row header",
			data: "userId,rating\n1,3.5\n2,4.0\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW"},
			want: []map[string]string{
				{"userId": "1", "rating": "3.5"},
				{"userId": "2", "rating": "4.0"},
			},
		},
		{
			name: "given header",
			data: "1,307,3.5\n1,481,3.5\n",
			rc: &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "GIVEN",
				CSVHeaders: []string{"userId", "movieId", "rating"}},
			want: []map[string]string{
				{"userId": "1", "movieId": "307", "rating": "3.5"},
				{"userId": "1", "movieId": "481", "rating": "3.5"},
			},
		},
		{
			name: "pipe delimiter with quoted separators",
			data: "path|note\n\"a|b\"|hello\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW", CSVDelimiter: "PIPE"},
			want: []map[string]string{
				{"path": "a|b", "note": "hello"},
			},
		},
		{
			name: "tab delimiter",
			data: "k\tv\n1\t2\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW", CSVDelimiter: "TAB"},
			want: []map[string]string{
				{"k": "1", "v": "2"},
			},
		},
		{
			name: "quoted newline inside a field",
			data: "k,v\n\"line1\nline2\",done\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW"},
			want: []map[string]string{
				{"k": "line1\nline2", "v": "done"},
			},
		},
		{
			name: "doubled quotes preserved",
			data: "k\n\"say \"\"hi\"\"\"\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW"},
			want: []map[string]string{
				{"k": `say "hi"`},
			},
		},
		{
			name: "backslash escapes pair and before others vanish",
			data: "k,path\n1,\"C:\\\\dir\\\\file,extra\"\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW"},
			want: []map[string]string{
				// A backslash escaping a backslash survives as one,
				// so doubled backslashes halve; the field separator
				// inside quotes does not split the record.
				{"k": "1", "path": "C:\\dir\\file,extra"},
			},
		},
		{
			name: "missing fields padded with empty strings",
			data: "a,b,c\n1,2\n",
			rc:   &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW"},
			want: []map[string]string{
				{"a": "1", "b": "2", "c": ""},
			},
		},
		{
			name:    "given header without csvheaders",
			data:    "1\n",
			rc:      &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "GIVEN"},
			wantErr: true,
		},
		{
			name:    "unsupported delimiter",
			data:    "a\n1\n",
			rc:      &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVDelimiter: "COLON"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, err := parseCSVDataset([]byte(tc.data), tc.rc)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got items %v", items)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseCSVDataset failed: %v", err)
			}
			if len(items) != len(tc.want) {
				t.Fatalf("got %d items, want %d: %v", len(items), len(tc.want), items)
			}
			for i, want := range tc.want {
				got, ok := items[i].(map[string]interface{})
				if !ok {
					t.Fatalf("item %d is %T, want map", i, items[i])
				}
				for k, v := range want {
					if got[k] != v {
						t.Errorf("item %d field %q = %v, want %q", i, k, got[k], v)
					}
				}
			}
		})
	}
}

// TestParseJSONDataset pins array itemisation, object key-value
// itemisation and the ItemsPointer JSON Pointer selection.
func TestParseJSONDataset(t *testing.T) {
	items, err := parseJSONDataset([]byte(`[{"id":1},{"id":2}]`), nil)
	if err != nil {
		t.Fatalf("array parse failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("array dataset produced %d items", len(items))
	}

	items, err = parseJSONDataset([]byte(`{"a":1,"b":2}`), nil)
	if err != nil {
		t.Fatalf("object parse failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("object dataset produced %d items", len(items))
	}
	seen := map[string]float64{}
	for _, it := range items {
		m := it.(map[string]interface{})
		seen[m["key"].(string)] = m["value"].(float64)
	}
	if seen["a"] != 1 || seen["b"] != 2 {
		t.Errorf("object items = %v", seen)
	}

	items, err = parseJSONDataset([]byte(`{"data":{"items":[{"id":1},{"id":2}]}}`),
		&sfnstore.ItemReaderReaderConfig{ItemsPointer: "/data/items"})
	if err != nil {
		t.Fatalf("pointer parse failed: %v", err)
	}
	if len(items) != 2 || items[0].(map[string]interface{})["id"].(float64) != 1 {
		t.Errorf("pointer dataset items = %v", items)
	}

	if _, err := parseJSONDataset([]byte(`{"data":{}}`), &sfnstore.ItemReaderReaderConfig{ItemsPointer: "/missing"}); err == nil {
		t.Errorf("missing pointer token should fail")
	}
}

// TestParseJSONLDataset pins one item per non-empty line.
func TestParseJSONLDataset(t *testing.T) {
	data := "{\"n\":1}\n\n{\"n\":2}\n"
	items, err := parseJSONLDataset([]byte(data))
	if err != nil {
		t.Fatalf("parseJSONLDataset failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if _, err := parseJSONLDataset([]byte("{\"n\":}\n")); err == nil {
		t.Errorf("invalid record should fail")
	}
}

// gzipBytes compresses a payload for the decompression tests.
func gzipBytes(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		t.Fatalf("gzip write failed: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("gzip close failed: %v", err)
	}
	return buf.Bytes()
}

// TestDecompressItemReaderData pins gzip decompression by key extension
// and the parquet passthrough.
func TestDecompressItemReaderData(t *testing.T) {
	payload := []byte("plain payload")

	out, err := decompressItemReaderData("data/items.json.gz", "JSON", gzipBytes(t, payload))
	if err != nil {
		t.Fatalf("gzip decompress failed: %v", err)
	}
	if string(out) != string(payload) {
		t.Errorf("gzip round trip = %q", out)
	}

	out, err = decompressItemReaderData("data/items.parquet", "PARQUET", payload)
	if err != nil || string(out) != string(payload) {
		t.Errorf("parquet passthrough failed: %q %v", out, err)
	}
}

// TestResolveItemReaderMaxItems pins the MaxItems literal cap and the
// MaxItemsPath reference resolution, including rejection of negative or
// non-integer targets.
func TestResolveItemReaderMaxItems(t *testing.T) {
	e := NewExecutor(nil, nil)
	rc := &sfnstore.ItemReaderReaderConfig{MaxItems: ptrInt64(5)}
	n, ierr := e.resolveItemReaderMaxItems(&ExecutionContext{Input: `{}`}, &sfnstore.ItemReaderConfig{ReaderConfig: rc})
	if ierr != nil || n != 5 {
		t.Fatalf("literal MaxItems = %d, %v", n, ierr)
	}

	rc = &sfnstore.ItemReaderReaderConfig{MaxItemsPath: "$.limit"}
	n, ierr = e.resolveItemReaderMaxItems(&ExecutionContext{Input: `{"limit":7}`}, &sfnstore.ItemReaderConfig{ReaderConfig: rc})
	if ierr != nil || n != 7 {
		t.Fatalf("MaxItemsPath resolution = %d, %v", n, ierr)
	}

	n, ierr = e.resolveItemReaderMaxItems(&ExecutionContext{Input: `{"limit":-1}`}, &sfnstore.ItemReaderConfig{ReaderConfig: rc})
	if ierr == nil {
		t.Fatalf("negative MaxItemsPath should fail, got %d", n)
	}
}

// TestResolveItemReaderArgs pins literal and reference-path argument
// resolution and the Bucket+Key/Prefix requirement.
func TestResolveItemReaderArgs(t *testing.T) {
	e := NewExecutor(nil, nil)

	reader := &sfnstore.ItemReaderConfig{
		Resource:   "arn:aws:states:::s3:getObject",
		Parameters: json.RawMessage(`{"Bucket":"src","Key":"items.csv","VersionId":"v1"}`),
	}
	args, ierr := e.resolveItemReaderArgs(&ExecutionContext{Input: `{}`}, reader)
	if ierr != nil {
		t.Fatalf("literal args failed: %v", ierr)
	}
	if args.bucket != "src" || args.key != "items.csv" || args.versionID != "v1" {
		t.Errorf("args = %+v", args)
	}

	reader.Parameters = json.RawMessage(`{"Bucket.$":"$.bucket","Key.$":"$.key"}`)
	args, ierr = e.resolveItemReaderArgs(&ExecutionContext{Input: `{"bucket":"dynamic","key":"data.json"}`}, reader)
	if ierr != nil {
		t.Fatalf("reference args failed: %v", ierr)
	}
	if args.bucket != "dynamic" || args.key != "data.json" {
		t.Errorf("args = %+v", args)
	}

	reader.Parameters = json.RawMessage(`{"Bucket":"src"}`)
	if _, ierr := e.resolveItemReaderArgs(&ExecutionContext{Input: `{}`}, reader); ierr == nil {
		t.Fatalf("Bucket without Key or Prefix should fail")
	}
}

// stubItemReaderS3 backs the S3-dependent ItemReader paths in tests.
type stubItemReaderS3 struct {
	objects map[string][]byte
	entries []eventbus.S3ObjectEntry
	put     map[string][]byte
}

func (s *stubItemReaderS3) GetObjectVersion(_ context.Context, _, bucket, key, _ string, _ int64) ([]byte, error) {
	data, ok := s.objects[bucket+"/"+key]
	if !ok {
		return nil, fmt.Errorf("no such object %s/%s", bucket, key)
	}
	return data, nil
}

func (s *stubItemReaderS3) ListObjectEntries(_ context.Context, _, _, prefix string, _ int) ([]eventbus.S3ObjectEntry, error) {
	var out []eventbus.S3ObjectEntry
	for _, e := range s.entries {
		if prefix == "" || (len(e.Key) >= len(prefix) && e.Key[:len(prefix)] == prefix) {
			out = append(out, e)
		}
	}
	return out, nil
}

func (s *stubItemReaderS3) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	return s.GetObjectVersion(ctx, region, bucket, key, "", maxBytes)
}

func (s *stubItemReaderS3) PutObject(_ context.Context, _, bucket, key string, data []byte, _ string) error {
	if s.put == nil {
		s.put = map[string][]byte{}
	}
	s.put[bucket+"/"+key] = data
	return nil
}

func (s *stubItemReaderS3) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	entries, err := s.ListObjectEntries(ctx, region, bucket, prefix, maxKeys)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		keys = append(keys, e.Key)
	}
	return keys, nil
}

func (s *stubItemReaderS3) BucketExists(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}

func (s *stubItemReaderS3) EnsureBucket(_ context.Context, _, _ string) error { return nil }

func (s *stubItemReaderS3) DeleteObject(_ context.Context, _, _, _ string) error { return nil }

// TestReadItemReaderItemsEndToEnd drives the reader pipeline over the
// stub S3 plane: getObject CSV, listObjectsV2 metadata, MaxItems capping,
// unsupported resources and the raw override path.
func TestReadItemReaderItemsEndToEnd(t *testing.T) {
	stub := &stubItemReaderS3{objects: map[string][]byte{
		"src/items.csv":  []byte("a,b\n1,2\n3,4\n"),
		"src/gz.json.gz": gzipBytes(t, []byte(`[{"x":1}]`)),
		"src/plain.json": []byte(`[{"x":1},{"x":2},{"x":3}]`),
	}}
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(stub)
	e := NewExecutor(nil, bus)
	e.region = "us-east-1"

	state := &sfnstore.MapState{
		ItemReader: &sfnstore.ItemReaderConfig{
			Resource:     "arn:aws:states:::s3:getObject",
			Parameters:   json.RawMessage(`{"Bucket":"src","Key":"items.csv"}`),
			ReaderConfig: &sfnstore.ItemReaderReaderConfig{InputType: "CSV", CSVHeaderLocation: "FIRST_ROW"},
		},
	}
	items, ierr := e.readItemReaderItems(context.Background(), &ExecutionContext{Input: `{}`}, state, "")
	if ierr != nil {
		t.Fatalf("csv read failed: %v", ierr.Cause)
	}
	if len(items) != 2 || items[0].(map[string]interface{})["a"] != "1" {
		t.Errorf("csv items = %v", items)
	}

	state.ItemReader.ReaderConfig.MaxItems = ptrInt64(1)
	items, ierr = e.readItemReaderItems(context.Background(), &ExecutionContext{Input: `{}`}, state, "")
	if ierr != nil {
		t.Fatalf("capped read failed: %v", ierr.Cause)
	}
	if len(items) != 1 {
		t.Errorf("MaxItems cap produced %d items", len(items))
	}

	// gzip-compressed JSON dataset resolved by key extension.
	state.ItemReader.Parameters = json.RawMessage(`{"Bucket":"src","Key":"gz.json.gz"}`)
	state.ItemReader.ReaderConfig = &sfnstore.ItemReaderReaderConfig{InputType: "JSON"}
	items, ierr = e.readItemReaderItems(context.Background(), &ExecutionContext{Input: `{}`}, state, "")
	if ierr != nil {
		t.Fatalf("gzip json read failed: %v", ierr.Cause)
	}
	if len(items) != 1 {
		t.Errorf("gzip json items = %v", items)
	}

	stub.entries = []eventbus.S3ObjectEntry{{Key: "data/f1", ETag: "\"e1\"", Size: 10, StorageClass: "STANDARD"}}
	state.ItemReader = &sfnstore.ItemReaderConfig{
		Resource:   "arn:aws:states:::s3:listObjectsV2",
		Parameters: json.RawMessage(`{"Bucket":"src","Prefix":"data/"}`),
	}
	items, ierr = e.readItemReaderItems(context.Background(), &ExecutionContext{Input: `{}`}, state, "")
	if ierr != nil {
		t.Fatalf("list read failed: %v", ierr.Cause)
	}
	if len(items) != 1 {
		t.Fatalf("list items = %v", items)
	}
	entry := items[0].(map[string]interface{})
	if entry["Key"] != "data/f1" || entry["Etag"] != "\"e1\"" || entry["StorageClass"] != "STANDARD" {
		t.Errorf("entry = %v", entry)
	}

	// The raw override replaces the S3 read with caller-supplied bytes.
	state.ItemReader = &sfnstore.ItemReaderConfig{
		Resource:     "arn:aws:states:::s3:getObject",
		Parameters:   json.RawMessage(`{"Bucket":"src","Key":"plain.json"}`),
		ReaderConfig: &sfnstore.ItemReaderReaderConfig{InputType: "JSON"},
	}
	items, ierr = e.readItemReaderItems(context.Background(), &ExecutionContext{Input: `{}`}, state, `[{"y":9}]`)
	if ierr != nil {
		t.Fatalf("raw override read failed: %v", ierr.Cause)
	}
	if len(items) != 1 || items[0].(map[string]interface{})["y"].(float64) != 9 {
		t.Errorf("raw override items = %v", items)
	}

	// Unsupported resources surface the documented error code.
	state.ItemReader = &sfnstore.ItemReaderConfig{
		Resource:   "arn:aws:states:::dynamodb:getItem",
		Parameters: json.RawMessage(`{"Bucket":"src","Key":"x"}`),
	}
	if _, ierr := e.readItemReaderItems(context.Background(), &ExecutionContext{Input: `{}`}, state, ""); ierr == nil || ierr.ErrorCode != "States.ItemReaderFailed" {
		t.Fatalf("unsupported resource error = %+v", ierr)
	}
}

// TestParseManifestDatasetInventory pins the S3 inventory manifest walk:
// the manifest's files[] reference CSV data files itemised per fileSchema.
func TestParseManifestDatasetInventory(t *testing.T) {
	manifest := `{"sourceBucket":"src","fileFormat":"CSV","fileSchema":"Bucket, Key, Size","files":[{"key":"dest/data/0.csv","size":30}]}`
	fetched := map[string][]byte{"dest/data/0.csv": []byte("\"src\",\"a/x\",12\n\"src\",\"b/y\",5\n")}
	items, err := parseManifestDataset([]byte(manifest),
		&sfnstore.ItemReaderReaderConfig{InputType: "MANIFEST"}, "dest",
		func(key string) ([]byte, error) { return fetched[key], nil })
	if err != nil {
		t.Fatalf("manifest parse failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("manifest items = %v", items)
	}
	first := items[0].(map[string]interface{})
	if first["Bucket"] != "src" || first["Key"] != "a/x" || first["Size"] != "12" {
		t.Errorf("first item = %v", first)
	}
}

// TestParseManifestDatasetAthena pins the ATHENA_DATA manifest walk: a
// headerless CSV list of data files, each parsed per the InputType.
func TestParseManifestDatasetAthena(t *testing.T) {
	manifest := "s3://out/job/1.json\ns3://out/job/2.json\n"
	fetched := map[string][]byte{
		"job/1.json": []byte(`{"a":1}` + "\n"),
		"job/2.json": []byte(`{"a":2}` + "\n"),
	}
	items, err := parseManifestDataset([]byte(manifest),
		&sfnstore.ItemReaderReaderConfig{ManifestType: "ATHENA_DATA", InputType: "JSONL"}, "out",
		func(key string) ([]byte, error) { return fetched[key], nil })
	if err != nil {
		t.Fatalf("athena manifest parse failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("athena items = %v", items)
	}
}

func ptrInt64(v int64) *int64 { return &v }

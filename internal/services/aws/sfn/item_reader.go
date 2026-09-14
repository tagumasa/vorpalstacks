package sfn

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"

	"github.com/klauspost/compress/zstd"
	"github.com/parquet-go/parquet-go"
)

// itemReaderArgs carries the resolved ItemReader argument object.
type itemReaderArgs struct {
	bucket    string
	key       string
	prefix    string
	versionID string
}

// readItemReaderItems resolves the dataset of a Map state's ItemReader and
// returns it with the reader arguments the dataset was read through, so the
// caller derives the context object's Map.Item.Source from the same
// resolved values (a single interpretation of the reader parameters). A
// nil ItemReader returns a nil slice and zero arguments. Every failure
// surfaces as a States.ItemReaderFailed execution error, matching the
// documented error contract ("A Map state failed because it couldn't read
// from the item source specified in the ItemReader field"). rawOverride,
// when non-empty, replaces the S3 read with the caller-supplied raw source
// bytes (the TestState mapItemReaderData contract: the data as found at its
// original source).
func (e *Executor) readItemReaderItems(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, rawOverride string) ([]interface{}, itemReaderArgs, *ExecutionError) {
	reader := state.ItemReader
	if reader == nil {
		return nil, itemReaderArgs{}, nil
	}

	rc := reader.ReaderConfig
	inputType := ""
	if rc != nil {
		inputType = rc.InputType
	}

	fetchObject := func(bucket, key, versionID string) ([]byte, *ExecutionError) {
		if rawOverride != "" {
			return []byte(rawOverride), nil
		}
		if e.bus == nil || e.bus.S3Invoker() == nil {
			return nil, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "S3 integration is not available"}
		}
		data, err := e.bus.S3Invoker().GetObjectVersion(ctx, e.region, bucket, key, versionID, sfnstore.MaxItemReaderFileBytes)
		if err != nil {
			return nil, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: fmt.Sprintf("failed to read s3://%s/%s: %s", bucket, key, err.Error())}
		}
		return data, nil
	}

	args, ierr := e.resolveItemReaderArgs(ctx, execCtx, reader)
	if ierr != nil {
		return nil, itemReaderArgs{}, ierr
	}

	resource := reader.Resource
	var items []interface{}
	switch {
	case strings.HasSuffix(resource, ":s3:getObject"), strings.HasSuffix(resource, ":aws-sdk:s3:getObject"):
		if inputType == "" {
			inputType = "JSON"
		}
		var data []byte
		if rawOverride != "" {
			data = []byte(rawOverride)
		} else {
			// The reader's own object honours a resolved VersionID; the
			// manifest/listing entries below address live keys, so their
			// fetches carry no version pin.
			fetched, ferr := fetchObject(args.bucket, args.key, args.versionID)
			if ferr != nil {
				return nil, args, ferr
			}
			data = fetched
		}
		// The manifest reader has three documented spellings: InputType
		// MANIFEST, ManifestType ATHENA_DATA (InputType CSV, JSONL or
		// Parquet names the data files' format), and ManifestType
		// S3_INVENTORY (no InputType; the type is assumed CSV).
		manifestReader := inputType == "MANIFEST" || (rc != nil && rc.ManifestType != "")
		if manifestReader {
			parsed, merr := parseManifestDataset(data, rc, args.bucket, func(key string) ([]byte, error) {
				fileData, ferr := fetchObject(args.bucket, key, "")
				if ferr != nil {
					return nil, fmt.Errorf("%s", ferr.Cause)
				}
				return fileData, nil
			})
			if merr != nil {
				return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: merr.Error()}
			}
			items = parsed
		} else {
			data, err := decompressItemReaderData(args.key, inputType, data)
			if err != nil {
				return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: err.Error()}
			}
			parsed, err := parseItemReaderData(data, inputType, rc)
			if err != nil {
				return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: err.Error()}
			}
			items = parsed
		}
	case strings.HasSuffix(resource, ":s3:listObjectsV2"):
		if rawOverride != "" {
			return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "raw reader data cannot substitute an s3:listObjectsV2 ItemReader"}
		}
		if e.bus == nil || e.bus.S3Invoker() == nil {
			return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "S3 integration is not available"}
		}
		entries, err := e.bus.S3Invoker().ListObjectEntries(ctx, e.region, args.bucket, args.prefix, 0)
		if err != nil {
			return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: fmt.Sprintf("failed to list s3://%s/%s: %s", args.bucket, args.prefix, err.Error())}
		}
		if rc != nil && rc.Transformation == "LOAD_AND_FLATTEN" {
			if inputType == "" {
				inputType = "JSON"
			}
			for _, entry := range entries {
				data, gerr := fetchObject(args.bucket, entry.Key, "")
				if gerr != nil {
					return nil, args, gerr
				}
				data, err = decompressItemReaderData(entry.Key, inputType, data)
				if err != nil {
					return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: err.Error()}
				}
				parsed, err := parseItemReaderData(data, inputType, rc)
				if err != nil {
					return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: err.Error()}
				}
				items = append(items, parsed...)
			}
		} else {
			for _, entry := range entries {
				items = append(items, map[string]interface{}{
					"Etag":         entry.ETag,
					"Key":          entry.Key,
					"LastModified": entry.LastModified,
					"Size":         entry.Size,
					"StorageClass": entry.StorageClass,
				})
			}
		}
	default:
		return nil, args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: fmt.Sprintf("unsupported ItemReader Resource: %s", resource)}
	}

	maxItems, ierr := e.resolveItemReaderMaxItems(execCtx, reader)
	if ierr != nil {
		return nil, args, ierr
	}
	if maxItems >= 0 && int64(len(items)) > maxItems {
		items = items[:maxItems]
	}
	return items, args, nil
}

// itemSourceFromReader derives the provenance the context object reports
// as Map.Item.Source, from the same resolved reader arguments the dataset
// was read through: "For state input, the value will be: STATE_DATA"; a
// keyed object read reports "the Amazon S3 URI. For example:
// S3://bucket-name/object-key"; a listing read (no key) reports "the S3
// URI for the bucket. For example: S3://bucket-name".
func itemSourceFromReader(reader *sfnstore.ItemReaderConfig, args itemReaderArgs) string {
	if reader == nil || args.bucket == "" {
		return "STATE_DATA"
	}
	if args.key != "" {
		return "S3://" + args.bucket + "/" + args.key
	}
	return "S3://" + args.bucket
}

// resolveReaderArguments renders the ItemReader/ResultWriter argument
// object for both dialects. JSONPath Parameters pass through for the
// "Field.$" reference-path lookup; JSONata Arguments ("For JSONata
// states, Parameters will be replaced with Arguments") template-resolve
// their {% %} expressions against the state input. Both empty returns a
// nil map: the WriterConfig-only preview form carries no arguments.
func (e *Executor) resolveReaderArguments(ctx context.Context, execCtx *ExecutionContext, parameters, arguments json.RawMessage) (map[string]interface{}, error) {
	if len(parameters) == 0 && len(arguments) == 0 {
		return nil, nil
	}
	if len(parameters) == 0 {
		// The reader/writer Arguments are raw JSON on the store shape;
		// parse them so the template walk sees the object, not the bytes.
		var argumentsValue interface{}
		if err := json.Unmarshal(arguments, &argumentsValue); err != nil {
			return nil, err
		}
		var inputData interface{}
		if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
			inputData = nil
		}
		statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
		resolved, err := e.applyJSONataArguments(ctx, argumentsValue, statesVar, execCtx.VariableScope)
		if err != nil {
			return nil, err
		}
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(resolved), &params); err != nil {
			return nil, fmt.Errorf("Arguments did not resolve to a JSON object: %s", err.Error())
		}
		return params, nil
	}
	var params map[string]interface{}
	if err := json.Unmarshal(parameters, &params); err != nil {
		return nil, err
	}
	return params, nil
}

// readerParamLookup resolves a reader/writer parameter by name: a literal
// string value, or the JSONPath reference-path form ("Field.$") selecting
// from the state input.
func readerParamLookup(params, input map[string]interface{}) func(string) (string, bool) {
	return func(name string) (string, bool) {
		if v, ok := params[name]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
		// JSONPath reference-path form: "Bucket.$": "$.bucket".
		if v, ok := params[name+".$"]; ok {
			if s, ok := v.(string); ok && input != nil && strings.HasPrefix(s, "$") {
				if resolved, err := getJSONPathValue(input, s); err == nil {
					switch rv := resolved.(type) {
					case string:
						return rv, true
					case float64:
						return strconv.FormatFloat(rv, 'f', -1, 64), true
					}
				}
			}
		}
		return "", false
	}
}

// resolveItemReaderArgs resolves the Bucket, Key, Prefix and VersionId
// arguments. JSONPath definitions carry reference paths on "Field.$" keys;
// JSONata definitions may template the values. Literal values pass through
// unchanged.
func (e *Executor) resolveItemReaderArgs(ctx context.Context, execCtx *ExecutionContext, reader *sfnstore.ItemReaderConfig) (itemReaderArgs, *ExecutionError) {
	args := itemReaderArgs{}
	// The failing member is named by the dialect the definition carries:
	// JSONPath Parameters or the JSONata Arguments that replaced them.
	member := "Parameters"
	if len(reader.Parameters) == 0 {
		member = "Arguments"
	}
	params, perr := e.resolveReaderArguments(ctx, execCtx, reader.Parameters, reader.Arguments)
	if perr != nil {
		return args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "ItemReader " + member + " is not a JSON object: " + perr.Error()}
	}

	var input map[string]interface{}
	if execCtx.Input != "" {
		if err := json.Unmarshal([]byte(execCtx.Input), &input); err != nil {
			input = nil
		}
	}
	lookup := readerParamLookup(params, input)

	if v, ok := lookup("Bucket"); ok {
		args.bucket = v
	}
	if v, ok := lookup("Key"); ok {
		args.key = v
	}
	if v, ok := lookup("Prefix"); ok {
		args.prefix = v
	}
	if v, ok := lookup("VersionId"); ok {
		args.versionID = v
	}

	if args.bucket == "" || (args.key == "" && args.prefix == "") {
		return args, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "ItemReader requires Bucket with Key or Prefix"}
	}
	return args, nil
}

// resolveItemReaderMaxItems resolves the effective MaxItems limit. A
// MaxItemsPath reference resolves against the state input; an unset limit
// returns -1 (no cap).
func (e *Executor) resolveItemReaderMaxItems(execCtx *ExecutionContext, reader *sfnstore.ItemReaderConfig) (int64, *ExecutionError) {
	rc := reader.ReaderConfig
	if rc == nil {
		return -1, nil
	}
	if rc.MaxItems != nil {
		return *rc.MaxItems, nil
	}
	if rc.MaxItemsPath == "" {
		return -1, nil
	}
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(execCtx.Input), &input); err != nil {
		return -1, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "MaxItemsPath requires an object state input: " + err.Error()}
	}
	resolved, err := getJSONPathValue(input, rc.MaxItemsPath)
	if err != nil {
		return -1, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: "MaxItemsPath resolution failed: " + err.Error()}
	}
	n, ok := resolved.(float64)
	if !ok || n < 0 {
		return -1, &ExecutionError{ErrorCode: "States.ItemReaderFailed", Cause: fmt.Sprintf("MaxItemsPath must resolve to a non-negative integer, got %v", resolved)}
	}
	return int64(n), nil
}

// decompressItemReaderData undoes the documented external compression
// types (GZIP, ZSTD) detected from the object key extension. Parquet
// datasets carry their compression internally, so they pass through.
func decompressItemReaderData(key, inputType string, data []byte) ([]byte, error) {
	if inputType == "PARQUET" {
		return data, nil
	}
	lower := strings.ToLower(key)
	switch {
	case strings.HasSuffix(lower, ".gz") || strings.HasSuffix(lower, ".gzip"):
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decompress %s: %s", key, err.Error())
		}
		defer r.Close()
		out, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress %s: %s", key, err.Error())
		}
		return out, nil
	case strings.HasSuffix(lower, ".zst") || strings.HasSuffix(lower, ".zstd"):
		r, err := zstd.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decompress %s: %s", key, err.Error())
		}
		defer r.Close()
		out, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress %s: %s", key, err.Error())
		}
		return out, nil
	default:
		return data, nil
	}
}

// parseItemReaderData parses raw dataset bytes into items per the
// InputType. Manifest datasets are handled by parseManifestDataset, which
// needs the S3 fetcher to follow the manifest's file references.
func parseItemReaderData(data []byte, inputType string, rc *sfnstore.ItemReaderReaderConfig) ([]interface{}, error) {
	switch inputType {
	case "JSON":
		return parseJSONDataset(data, rc)
	case "JSONL":
		return parseJSONLDataset(data)
	case "CSV":
		return parseCSVDataset(data, rc)
	case "PARQUET":
		return parseParquetDataset(data)
	default:
		return nil, fmt.Errorf("unsupported ItemReader InputType: %s", inputType)
	}
}

// parseJSONDataset turns a JSON array into items, or a JSON object into
// one {"key": …, "value": …} item per key-value pair. ItemsPointer selects
// a nested array or object first when present.
func parseJSONDataset(data []byte, rc *sfnstore.ItemReaderReaderConfig) ([]interface{}, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("JSON dataset is empty")
	}
	var doc interface{}
	if err := json.Unmarshal(trimmed, &doc); err != nil {
		return nil, fmt.Errorf("JSON dataset is not valid JSON: %s", err.Error())
	}
	if rc != nil && rc.ItemsPointer != "" {
		pointed, err := resolveJSONPointer(doc, rc.ItemsPointer)
		if err != nil {
			return nil, err
		}
		doc = pointed
	}
	return jsonDocumentItems(doc)
}

// jsonDocumentItems converts a decoded JSON document to dataset items.
func jsonDocumentItems(doc interface{}) ([]interface{}, error) {
	switch v := doc.(type) {
	case []interface{}:
		return v, nil
	case map[string]interface{}:
		items := make([]interface{}, 0, len(v))
		for k, val := range v {
			items = append(items, map[string]interface{}{"key": k, "value": val})
		}
		return items, nil
	default:
		return nil, fmt.Errorf("JSON dataset must be an array or object")
	}
}

// resolveJSONPointer walks an RFC 6901 JSON Pointer (~0 escapes ~, ~1
// escapes /) over a decoded JSON document.
func resolveJSONPointer(doc interface{}, pointer string) (interface{}, error) {
	if pointer == "" {
		return doc, nil
	}
	if !strings.HasPrefix(pointer, "/") {
		return nil, fmt.Errorf("ItemsPointer must be a JSON Pointer starting with '/': %s", pointer)
	}
	current := doc
	for _, rawToken := range strings.Split(pointer[1:], "/") {
		token := strings.ReplaceAll(strings.ReplaceAll(rawToken, "~1", "/"), "~0", "~")
		switch node := current.(type) {
		case map[string]interface{}:
			next, ok := node[token]
			if !ok {
				return nil, fmt.Errorf("ItemsPointer token '%s' not found", token)
			}
			current = next
		case []interface{}:
			idx, err := strconv.Atoi(token)
			if err != nil || idx < 0 || idx >= len(node) {
				return nil, fmt.Errorf("ItemsPointer array index '%s' out of range", token)
			}
			current = node[idx]
		default:
			return nil, fmt.Errorf("ItemsPointer token '%s' cannot descend a scalar", token)
		}
	}
	return current, nil
}

// parseJSONLDataset turns newline-delimited JSON into one item per line.
func parseJSONLDataset(data []byte) ([]interface{}, error) {
	var items []interface{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var item interface{}
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			return nil, fmt.Errorf("JSONL dataset has an invalid record: %s", err.Error())
		}
		items = append(items, item)
	}
	return items, nil
}

// csvDelimiterRune maps the documented CSVDelimiter values to their
// separator runes. COMMA is the default.
func csvDelimiterRune(name string) (rune, error) {
	switch name {
	case "", "COMMA":
		return ',', nil
	case "PIPE":
		return '|', nil
	case "SEMICOLON":
		return ';', nil
	case "SPACE":
		return ' ', nil
	case "TAB":
		return '\t', nil
	default:
		return 0, fmt.Errorf("unsupported CSVDelimiter: %s", name)
	}
}

// parseCSVDataset parses a text delimited dataset under the documented
// Step Functions rules: the configured delimiter separates fields,
// newlines separate records, quoted fields may contain delimiters and
// newlines (quotes escape by doubling), backslashes escape backslashes,
// quotes and the field separator (a backslash before any other character
// is silently removed), rows with fewer fields than the header gain empty
// strings and rows with more fields drop the extras. All field values are
// strings.
func parseCSVDataset(data []byte, rc *sfnstore.ItemReaderReaderConfig) ([]interface{}, error) {
	delimiter, err := csvDelimiterRune("")
	if rc != nil {
		delimiter, err = csvDelimiterRune(rc.CSVDelimiter)
	}
	if err != nil {
		return nil, err
	}

	records, err := scanCSVRecords(string(data), delimiter)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	var headers []string
	headerLocation := "FIRST_ROW"
	if rc != nil && rc.CSVHeaderLocation != "" {
		headerLocation = rc.CSVHeaderLocation
	}
	dataRows := records
	switch headerLocation {
	case "FIRST_ROW":
		headers = records[0]
		dataRows = records[1:]
	case "GIVEN":
		headers = rc.CSVHeaders
		if len(headers) == 0 {
			return nil, fmt.Errorf("CSVHeaderLocation GIVEN requires CSVHeaders")
		}
	default:
		return nil, fmt.Errorf("unsupported CSVHeaderLocation: %s", headerLocation)
	}

	items := make([]interface{}, 0, len(dataRows))
	for _, row := range dataRows {
		item := make(map[string]interface{}, len(headers))
		for i, h := range headers {
			if i < len(row) {
				item[h] = row[i]
			} else {
				item[h] = ""
			}
		}
		items = append(items, item)
	}
	return items, nil
}

// scanCSVRecords splits raw CSV text into string records honouring quoted
// fields and the documented backslash escape rules.
func scanCSVRecords(text string, delimiter rune) ([][]string, error) {
	var records [][]string
	var record []string
	var field strings.Builder
	inQuotes := false
	flushField := func() {
		record = append(record, field.String())
		field.Reset()
	}
	flushRecord := func() {
		flushField()
		records = append(records, record)
		record = nil
	}

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case inQuotes:
			switch c {
			case '"':
				if i+1 < len(runes) && runes[i+1] == '"' {
					field.WriteRune('"')
					i++
				} else {
					inQuotes = false
				}
			case '\\':
				if i+1 < len(runes) {
					next := runes[i+1]
					if next == '\\' || next == '"' || next == delimiter {
						field.WriteRune(next)
						i++
					}
					// A backslash before any other character is
					// silently removed.
				}
			default:
				field.WriteRune(c)
			}
		case c == '"':
			inQuotes = true
		case c == delimiter:
			flushField()
		case c == '\r':
			// A CRLF pair is one record break, not two — two breaks
			// would emit an empty record after every row.
			if i+1 < len(runes) && runes[i+1] == '\n' {
				i++
			}
			flushRecord()
		case c == '\n':
			flushRecord()
		default:
			field.WriteRune(c)
		}
	}
	if field.Len() > 0 || len(record) > 0 {
		flushRecord()
	}
	return records, nil
}

// parseParquetDataset reads parquet rows into string-keyed item objects.
func parseParquetDataset(data []byte) ([]interface{}, error) {
	reader := parquet.NewReader(bytes.NewReader(data))
	defer reader.Close()
	var items []interface{}
	for {
		row := make(map[string]interface{})
		err := reader.Read(&row)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read parquet row: %s", err.Error())
		}
		items = append(items, row)
	}
	return items, nil
}

// s3InventoryManifest models the manifest.json of an Amazon S3 inventory
// report: the data files it lists hold the CSV item rows.
type s3InventoryManifest struct {
	SourceBucket      string `json:"sourceBucket"`
	DestinationBucket string `json:"destinationBucket"`
	FileFormat        string `json:"fileFormat"`
	FileSchema        string `json:"fileSchema"`
	Files             []struct {
		Key  string `json:"key"`
		Size int64  `json:"size"`
	} `json:"files"`
}

// parseManifestDataset reads a manifest and the data files it references.
// The default manifest is the Amazon S3 inventory manifest.json; the
// ATHENA_DATA variant is a headerless CSV list of data files parsed per
// the manifest's InputType.
func parseManifestDataset(data []byte, rc *sfnstore.ItemReaderReaderConfig, bucket string, fetch func(key string) ([]byte, error)) ([]interface{}, error) {
	manifestType := ""
	inputType := ""
	if rc != nil {
		manifestType = rc.ManifestType
		inputType = rc.InputType
	}

	switch manifestType {
	case "ATHENA_DATA":
		if inputType == "" {
			return nil, fmt.Errorf("ManifestType ATHENA_DATA requires InputType")
		}
		// Athena UNLOAD manifests are headerless single-column CSV
		// lists of the data files.
		records, err := scanCSVRecords(string(data), ',')
		if err != nil {
			return nil, err
		}
		var items []interface{}
		for _, rec := range records {
			if len(rec) == 0 || strings.TrimSpace(rec[0]) == "" {
				continue
			}
			ref := strings.TrimPrefix(rec[0], "s3://")
			key := ref
			if idx := strings.Index(ref, "/"); idx >= 0 {
				key = ref[idx+1:]
			}
			fileData, err := fetch(key)
			if err != nil {
				return nil, fmt.Errorf("failed to read manifest data file %s: %s", key, err.Error())
			}
			fileData, err = decompressItemReaderData(key, inputType, fileData)
			if err != nil {
				return nil, err
			}
			parsed, err := parseItemReaderData(fileData, inputType, rc)
			if err != nil {
				return nil, err
			}
			items = append(items, parsed...)
		}
		return items, nil
	default:
		// Default: the Amazon S3 inventory manifest.
		var manifest s3InventoryManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return nil, fmt.Errorf("inventory manifest is not valid JSON: %s", err.Error())
		}
		columns := strings.Split(manifest.FileSchema, ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}
		if len(columns) == 0 || (len(columns) == 1 && columns[0] == "") {
			return nil, fmt.Errorf("inventory manifest has no fileSchema")
		}
		var items []interface{}
		for _, f := range manifest.Files {
			fileData, err := fetch(f.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to read inventory data file %s: %s", f.Key, err.Error())
			}
			fileData, err = decompressItemReaderData(f.Key, "CSV", fileData)
			if err != nil {
				return nil, err
			}
			records, err := scanCSVRecords(string(fileData), ',')
			if err != nil {
				return nil, err
			}
			for _, row := range records {
				item := make(map[string]interface{}, len(columns))
				for i, col := range columns {
					if i < len(row) {
						item[col] = row[i]
					} else {
						item[col] = ""
					}
				}
				items = append(items, item)
			}
		}
		return items, nil
	}
}

// toleratedFailureThresholds carries the Distributed Map tolerated-failure
// thresholds resolved once against the processed state input — the same
// stage MaxConcurrency resolves at.
type toleratedFailureThresholds struct {
	count    float64
	hasCount bool
	pct      float64
	hasPct   bool
}

// resolveToleratedFailureThresholds resolves the literal, *Path and
// JSONata-expression threshold members. The *Path forms select from the
// processed input (after InputPath); a path or expression that fails to
// resolve surfaces its error instead of silently behaving as an
// unconfigured threshold — the unconfigured default is fail-fast, so a
// swallowed failure would invert the operator's configured tolerance.
func (e *Executor) resolveToleratedFailureThresholds(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, processedInput string) (toleratedFailureThresholds, *ExecutionError) {
	var thresholds toleratedFailureThresholds

	var input map[string]interface{}
	if processedInput != "" {
		if err := json.Unmarshal([]byte(processedInput), &input); err != nil {
			input = nil
		}
	}
	resolvePath := func(path, field string) (float64, bool, *ExecutionError) {
		if path == "" {
			return 0, false, nil
		}
		if input == nil {
			return 0, false, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("%s requires object input", field)}
		}
		resolved, err := getJSONPathValue(input, path)
		if err != nil {
			return 0, false, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("%s failed to resolve: %s", field, err.Error())}
		}
		if n, ok := resolved.(float64); ok {
			return n, true, nil
		}
		return 0, false, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("%s must resolve to a number", field)}
	}
	resolveMember := func(field string, raw interface{}, min, max float64) (float64, bool, *ExecutionError) {
		if raw == nil {
			return 0, false, nil
		}
		f, isNumber := toFloat64(raw)
		if isNumber && (f < min || (max > 0 && f > max)) {
			// An out-of-range literal is rejected, not discarded: the
			// JSONata-expression resolver errors on the same condition,
			// and discarding would silently restore the fail-fast
			// default.
			return 0, false, &ExecutionError{ErrorCode: "States.InvalidInput", Cause: fmt.Sprintf("Map state %s is outside the acceptable range", field)}
		}
		if !isNumber {
			// "In JSONata states, you can specify a JSONata expression
			// that evaluates to an integer" (Distributed-mode thresholds).
			resolved, present, rerr := e.resolveMapNumericMember(ctx, execCtx, field, raw, processedInput, min, max)
			if rerr != nil {
				return 0, false, rerr
			}
			if !present {
				return 0, false, nil
			}
			return resolved, true, nil
		}
		return f, true, nil
	}

	if n, ok, err := resolvePath(state.ToleratedFailureCountPath, "ToleratedFailureCountPath"); err != nil {
		return thresholds, err
	} else if ok {
		thresholds.count, thresholds.hasCount = n, true
	}
	if c, ok, err := resolveMember("ToleratedFailureCount", state.ToleratedFailureCount, 0, 0); err != nil {
		return thresholds, err
	} else if ok {
		thresholds.count, thresholds.hasCount = c, true
	}
	if n, ok, err := resolvePath(state.ToleratedFailurePercentagePath, "ToleratedFailurePercentagePath"); err != nil {
		return thresholds, err
	} else if ok {
		thresholds.pct, thresholds.hasPct = n, true
	}
	if p, ok, err := resolveMember("ToleratedFailurePercentage", state.ToleratedFailurePercentage, 0, 100); err != nil {
		return thresholds, err
	} else if ok {
		thresholds.pct, thresholds.hasPct = p, true
	}
	return thresholds, nil
}

// evaluateToleratedFailure applies the Distributed Map tolerated-failure
// thresholds. Inline maps keep the classic fail-fast behaviour, so both
// results are false. For Distributed maps the documented defaults apply
// (a zero count and a zero percentage: any failed item fails the Map Run),
// and a failure within the configured thresholds is tolerated instead.
// The first return value reports whether the failures stay within the
// thresholds; the second reports whether an exceeded threshold (rather
// than the generic iterator failure) caused the failure.
func (e *Executor) evaluateToleratedFailure(state *sfnstore.MapState, thresholds toleratedFailureThresholds, itemsFailed, totalItems int64) (tolerated bool, exceededThreshold bool) {
	if state.ItemProcessor == nil || state.ItemProcessor.ProcessorConfig == nil ||
		state.ItemProcessor.ProcessorConfig.Mode != "DISTRIBUTED" {
		return false, false
	}

	// The documented defaults are a zero count and a zero percentage, so
	// without any configured threshold a single failed item fails the
	// Map Run.
	exceeded := false
	if thresholds.hasCount && float64(itemsFailed) > thresholds.count {
		exceeded = true
	}
	if thresholds.hasPct && totalItems > 0 && float64(itemsFailed)*100.0/float64(totalItems) > thresholds.pct {
		exceeded = true
	}
	if !thresholds.hasCount && !thresholds.hasPct && itemsFailed > 0 {
		exceeded = true
	}
	if exceeded {
		return false, true
	}
	return itemsFailed > 0, false
}

// writeMapResultWriter renders the formatted child results of a Map Run
// and either returns them as the state output (the WriterConfig-only
// preview form) or exports them to S3 and returns the Map Run ARN plus
// the export location. Distributed units carry their dispatched child
// execution identity; inline units fall back to the synthesised per-unit
// identifiers. Failures surface as States.ResultWriterFailed.
func (e *Executor) writeMapResultWriter(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.MapState, mapRunRecord *sfnstore.MapRun, renderedResults []string, itemErrors []error, itemInputs []string, childMetas []mapChildMeta) (string, *ExecutionError) {
	// The documented field combinations are WriterConfig alone (preview),
	// Resource with Parameters (or Arguments), or all three. The preview
	// form carries no export destination: "Non-existent Resource and
	// Parameters fields ... imply the state output resource. The results
	// are passed on to the next state."
	exporting := len(state.ResultWriter.Parameters) > 0 || len(state.ResultWriter.Arguments) > 0

	transformation := "NONE"
	outputType := "JSON"
	if wc := state.ResultWriter.WriterConfig; wc != nil {
		if wc.Transformation == "COMPACT" || wc.Transformation == "FLATTEN" {
			transformation = wc.Transformation
		}
		if wc.OutputType == "JSONL" {
			outputType = "JSONL"
		}
	}

	// The export directory hangs off the Map Run identifier.
	mapRunID := mapRunRecord.MapRunArn
	if idx := strings.LastIndex(mapRunID, "/"); idx >= 0 {
		mapRunID = mapRunID[idx+1:]
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	type itemRecord struct {
		ExecutionArn string `json:"ExecutionArn"`
		Input        string `json:"Input"`
		InputDetails struct {
			Included bool `json:"Included"`
		} `json:"InputDetails"`
		Name          string `json:"Name"`
		Output        string `json:"Output"`
		OutputDetails struct {
			Included bool `json:"Included"`
		} `json:"OutputDetails"`
		RedriveCount        int    `json:"RedriveCount"`
		RedriveStatus       string `json:"RedriveStatus"`
		RedriveStatusReason string `json:"RedriveStatusReason"`
		StartDate           string `json:"StartDate"`
		StateMachineArn     string `json:"StateMachineArn"`
		Status              string `json:"Status"`
		StopDate            string `json:"StopDate"`
	}

	// collectRecords gathers the per-child execution records with their
	// workflow metadata; the empty status collects every child, which the
	// preview's NONE presentation uses.
	collectRecords := func(status string) []interface{} {
		var records []interface{}
		for i, res := range renderedResults {
			failed := i < len(itemErrors) && itemErrors[i] != nil
			if status != "" && (status == "SUCCEEDED") == failed {
				continue
			}
			recStatus := status
			if recStatus == "" {
				recStatus = "SUCCEEDED"
				if failed {
					recStatus = "FAILED"
				}
			}
			itemInput := "null"
			if i < len(itemInputs) {
				itemInput = itemInputs[i]
			}
			rec := itemRecord{
				ExecutionArn:        fmt.Sprintf("%s/%s/%s-%d", execCtx.Execution.ExecutionArn, mapRunRecord.Name, mapRunID, i),
				Input:               itemInput,
				Name:                fmt.Sprintf("%s-%d", mapRunID, i),
				Output:              res,
				RedriveStatus:       "NOT_REDRIVABLE",
				RedriveStatusReason: "Execution results are exported by the Map Run ResultWriter",
				StartDate:           now,
				StateMachineArn:     execCtx.Execution.StateMachineArn,
				Status:              recStatus,
				StopDate:            now,
			}
			if i < len(childMetas) && childMetas[i].Arn != "" {
				rec.ExecutionArn = childMetas[i].Arn
				rec.Name = childMetas[i].Name
				rec.RedriveCount = int(childMetas[i].RedriveCount)
			}
			rec.InputDetails.Included = true
			rec.OutputDetails.Included = true
			records = append(records, rec)
		}
		return records
	}

	// The transformation shapes the succeeded results: COMPACT keeps each
	// child output in its original structure, FLATTEN splices array
	// results into one array, and NONE (the default) returns the per-child
	// execution records. "If a child workflow execution fails, Step
	// Functions returns its execution result unchanged. The results would
	// be equivalent to having set Transformation to NONE", so a failed
	// child anywhere keeps the record presentation.
	renderTransformed := func() ([]interface{}, bool) {
		if transformation == "NONE" {
			return nil, false
		}
		payload := []interface{}{}
		for i, res := range renderedResults {
			if i < len(itemErrors) && itemErrors[i] != nil {
				return nil, false
			}
			var parsed interface{}
			if err := json.Unmarshal([]byte(res), &parsed); err != nil {
				parsed = res
			}
			if transformation == "FLATTEN" {
				if arr, ok := parsed.([]interface{}); ok {
					payload = append(payload, arr...)
					continue
				}
			}
			payload = append(payload, parsed)
		}
		return payload, true
	}

	// OutputType JSON renders the elements as a JSON array; JSONL renders
	// them as JSON Lines — for the preview form that formatted text is
	// itself the state output.
	encodeElements := func(elements []interface{}) interface{} {
		if outputType != "JSONL" {
			return elements
		}
		lines := make([]string, 0, len(elements))
		for _, el := range elements {
			if b, err := json.Marshal(el); err == nil {
				lines = append(lines, string(b))
			}
		}
		return strings.Join(lines, "\n") + "\n"
	}
	encode := func(elements []interface{}) []byte {
		b, err := json.Marshal(encodeElements(elements))
		if err != nil {
			return []byte("[]")
		}
		return b
	}

	if !exporting {
		elements, ok := renderTransformed()
		if !ok {
			elements = collectRecords("")
		}
		return string(encode(elements)), nil
	}

	if e.bus == nil || e.bus.S3Invoker() == nil {
		return "", &ExecutionError{ErrorCode: "States.ResultWriterFailed", Cause: "S3 integration is not available"}
	}
	// The failing member is named by the dialect the definition carries:
	// JSONPath Parameters or the JSONata Arguments that replaced them.
	writerMember := "Parameters"
	if len(state.ResultWriter.Parameters) == 0 {
		writerMember = "Arguments"
	}
	params, perr := e.resolveReaderArguments(ctx, execCtx, state.ResultWriter.Parameters, state.ResultWriter.Arguments)
	if perr != nil {
		return "", &ExecutionError{ErrorCode: "States.ResultWriterFailed", Cause: "ResultWriter " + writerMember + " is not a JSON object: " + perr.Error()}
	}
	var input map[string]interface{}
	if execCtx.Input != "" {
		if err := json.Unmarshal([]byte(execCtx.Input), &input); err != nil {
			input = nil
		}
	}
	lookup := readerParamLookup(params, input)
	bucket, _ := lookup("Bucket")
	prefix, _ := lookup("Prefix")
	if bucket == "" {
		return "", &ExecutionError{ErrorCode: "States.ResultWriterFailed", Cause: "ResultWriter requires a Bucket parameter"}
	}
	prefix = strings.Trim(prefix, "/")

	dir := prefix
	if dir != "" {
		dir += "/"
	}
	dir += mapRunID + "/"

	put := func(key string, data []byte) *ExecutionError {
		if err := e.bus.S3Invoker().PutObject(ctx, e.region, bucket, key, data, "application/json"); err != nil {
			return &ExecutionError{ErrorCode: "States.ResultWriterFailed", Cause: fmt.Sprintf("failed to write s3://%s/%s: %s", bucket, key, err.Error())}
		}
		return nil
	}

	resultFiles := map[string][]string{}
	failedCount := 0
	for _, err := range itemErrors {
		if err != nil {
			failedCount++
		}
	}
	if len(renderedResults)-failedCount > 0 {
		succeeded := collectRecords("SUCCEEDED")
		if payload, ok := renderTransformed(); ok {
			succeeded = payload
		}
		if ierr := put(dir+"SUCCEEDED_0.json", encode(succeeded)); ierr != nil {
			return "", ierr
		}
		resultFiles["SUCCEEDED"] = []string{"SUCCEEDED_0.json"}
	} else {
		resultFiles["SUCCEEDED"] = []string{}
	}
	if failedCount > 0 {
		if ierr := put(dir+"FAILED_0.json", encode(collectRecords("FAILED"))); ierr != nil {
			return "", ierr
		}
		resultFiles["FAILED"] = []string{"FAILED_0.json"}
	} else {
		resultFiles["FAILED"] = []string{}
	}
	resultFiles["PENDING"] = []string{}
	resultFiles["RUNNING"] = []string{}

	manifest := map[string]interface{}{
		"MapRunArn":           mapRunRecord.MapRunArn,
		"ResultLocation":      fmt.Sprintf("s3://%s/%s", bucket, prefix),
		"NumberOfFailedItems": failedCount,
		"ResultFiles":         resultFiles,
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return "", &ExecutionError{ErrorCode: "States.ResultWriterFailed", Cause: "failed to encode manifest: " + err.Error()}
	}
	manifestKey := dir + "manifest.json"
	if ierr := put(manifestKey, manifestJSON); ierr != nil {
		return "", ierr
	}

	exported := map[string]interface{}{
		"MapRunArn": mapRunRecord.MapRunArn,
		"ResultWriterDetails": map[string]interface{}{
			"Bucket": bucket,
			"Key":    manifestKey,
		},
	}
	out, err := json.Marshal(exported)
	if err != nil {
		return "", &ExecutionError{ErrorCode: "States.ResultWriterFailed", Cause: "failed to encode export summary: " + err.Error()}
	}
	return string(out), nil
}

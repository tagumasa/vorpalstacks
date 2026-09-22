package cloudwatchlogs

import (
	"encoding/json"
	"strings"
)

// The log transformer recipe engine. Every processor's semantics follow
// the documented processor pages (CloudWatch-Logs-Transformation:
// Configurable / StringMutate / JSONMutate / Datatype); the engine works
// on one record per event — a JSON object whose "@message" member holds
// the original message until a parser reads it — and emits the JSON
// serialisation of the processed record without "@message" as the
// transformed message. An event the recipe leaves without any extracted
// key has no transformed copy.

// TransformContext carries the event's platform identity — the
// metadata sources the copyValue processor documents (@logGroupName,
// @logGroupStream, @accountId, @regionName).
type TransformContext struct {
	LogGroup  string
	LogStream string
	AccountID string
	Region    string
}

// maxTransformMessageBytes is the documented per-event transformation
// bound: "logs transformation can only handle log event of size less
// than 512kb. Any log events greater than 512kb will fail in
// transformation and emit an error" (Transform logs during ingestion).
const maxTransformMessageBytes = 512 * 1024

// maxTransformedFieldCount is the documented extraction bound: "Each
// transformer can extract up to 200 fields from a log event."
const maxTransformedFieldCount = 200

// transformationErrorField is the documented failure marker: "Whenever
// CloudWatch Logs tries and fails to transform a log event, it adds a
// @transformationError system field to that log event. ... You can
// query for this field with a query such as filter
// ispresent(@transformationError)".
const transformationErrorField = "@transformationError"

// The failure reasons the marker carries. The documented contract is
// the field's presence; the value is the platform's own diagnostic.
const (
	transformationErrTooLarge     = "Log event exceeds the 512KB transformation size limit"
	transformationErrFieldCeiling = "Log event transformation exceeded the 200-field extraction limit"
	transformationErrUnserialised = "Transformed log event is not JSON-serialisable"
)

// transformFailureRecord renders the marker-only record a failed
// transformation persists as its transformed copy.
func transformFailureRecord(reason string) string {
	out, err := json.Marshal(map[string]interface{}{transformationErrorField: reason})
	if err != nil {
		return ""
	}
	return string(out)
}

// transformEvent applies the processor pipeline to one event message
// and returns the transformed message ("" when the recipe added
// nothing) and whether the transformation failed — a failure carries
// the @transformationError marker record in its transformed copy. A
// processor failure is non-fatal: the documented pipeline is
// best-effort per processor, and the event keeps its accumulated state.
func transformEvent(config []map[string]interface{}, message string, tctx TransformContext) (string, bool) {
	if len(message) >= maxTransformMessageBytes {
		return transformFailureRecord(transformationErrTooLarge), true
	}
	record := map[string]interface{}{"@message": message}
	for _, processor := range config {
		applyProcessor(record, processor, tctx)
	}
	delete(record, "@message")
	if len(record) == 0 {
		return "", false
	}
	if len(record) > maxTransformedFieldCount {
		return transformFailureRecord(transformationErrFieldCeiling), true
	}
	out, err := json.Marshal(record)
	if err != nil {
		return transformFailureRecord(transformationErrUnserialised), true
	}
	return string(out), false
}

// processorKind names the single processor-type member one Processor
// element carries — the Processor structure is a union, so an element
// with more than one member (or a non-object member) is not a processor
// configuration at all.
func processorKind(processor map[string]interface{}) (string, map[string]interface{}, bool) {
	if len(processor) != 1 {
		return "", nil, false
	}
	for kind, cfg := range processor {
		if m, ok := cfg.(map[string]interface{}); ok {
			return kind, m, true
		}
	}
	return "", nil, false
}

func applyProcessor(record map[string]interface{}, processor map[string]interface{}, tctx TransformContext) {
	kind, cfg, ok := processorKind(processor)
	if !ok {
		return
	}
	switch kind {
	case "parseJSON":
		applyParseJSON(record, cfg)
	case "grok":
		applyGrok(record, cfg)
	case "csv":
		applyCSV(record, cfg)
	case "parseCloudfront":
		applyParseCloudfront(record, cfg)
	case "parseRoute53":
		applyParseRoute53(record, cfg)
	case "parseVPC":
		applyParseVPC(record, cfg)
	case "parseWAF":
		applyParseWAF(record, cfg)
	case "parsePostgres":
		applyParsePostgres(record, cfg)
	case "parseKeyValue":
		applyParseKeyValue(record, cfg)
	case "addKeys":
		applyAddKeys(record, cfg)
	case "deleteKeys":
		applyDeleteKeys(record, cfg)
	case "moveKeys":
		applyMoveKeys(record, cfg)
	case "renameKeys":
		applyRenameKeys(record, cfg)
	case "copyValue":
		applyCopyValue(record, cfg, tctx)
	case "listToMap":
		applyListToMap(record, cfg)
	case "lowerCaseString":
		applyCaseString(record, cfg, strings.ToLower)
	case "upperCaseString":
		applyCaseString(record, cfg, strings.ToUpper)
	case "trimString":
		applyTrimString(record, cfg)
	case "splitString":
		applySplitString(record, cfg)
	case "substituteString":
		applySubstituteString(record, cfg)
	case "typeConverter":
		applyTypeConverter(record, cfg)
	case "dateTimeConverter":
		applyDateTimeConverter(record, cfg)
	}
}

// --- path helpers ---

// getPath reads a dotted path out of the record; the source "@message"
// addresses the original message.
func getPath(record map[string]interface{}, path string) (interface{}, bool) {
	if path == "@message" {
		v, ok := record["@message"]
		return v, ok
	}
	parts := strings.Split(path, ".")
	var current interface{} = record
	for _, part := range parts {
		asMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = asMap[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// setPath writes a dotted path into the record, creating intermediate
// objects; overwrite governs an existing leaf. The return reports
// whether the write landed: a blocked write (an existing leaf or a
// non-object intermediate without overwrite) leaves the record
// untouched.
func setPath(record map[string]interface{}, path string, value interface{}, overwrite bool) bool {
	parts := strings.Split(path, ".")
	current := record
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]interface{})
		if !ok {
			if _, exists := current[part]; exists && !overwrite {
				return false
			}
			next = map[string]interface{}{}
			current[part] = next
		}
		current = next
	}
	leaf := parts[len(parts)-1]
	if _, exists := current[leaf]; exists && !overwrite {
		return false
	}
	current[leaf] = value
	return true
}

// deletePath removes a dotted path; deleting the last member of an
// object leaves the empty parent in place (the documented deleteKeys
// example keeps "outer_key": {}).
func deletePath(record map[string]interface{}, path string) {
	parts := strings.Split(path, ".")
	current := record
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]interface{})
		if !ok {
			return
		}
		current = next
	}
	delete(current, parts[len(parts)-1])
}

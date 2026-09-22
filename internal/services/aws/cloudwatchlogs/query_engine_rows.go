package cloudwatchlogs

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// The query result row: field storage with insertion-ordered columns,
// and the builders that turn gathered events into result rows — the row
// shape, the event pointer vocabulary, JSON field discovery and the
// timestamp rendering every result row carries.

// queryResultRow represents a single row in the query results. columns
// records the order in which field names were first written so that output
// field order and "first field" selection stay deterministic. base carries
// the pre-projection source row once a projection command has narrowed the
// row, so later projection commands resolve against the unprojected record
// — the documented fields-union and last-display rules. structs carries
// the decoded trees of the structurally-nested fields — the values that
// were structure at ingest or creation time, the only values dot notation
// traverses ("Dot notation traverses only fields that are stored as
// structurally nested JSON objects at ingest time"); a JSON-encoded
// string value has no structs entry and stays an opaque leaf. The trees
// are read-only once stored; cloned rows share them.
type queryResultRow struct {
	fields  map[string]string
	columns []string
	base    *queryResultRow
	structs map[string]interface{}
}

// set stores a field value, appending new field names to the column order.
func (r *queryResultRow) set(k, v string) {
	if r.fields == nil {
		r.fields = make(map[string]string)
	}
	if _, ok := r.fields[k]; !ok {
		r.columns = append(r.columns, k)
	}
	r.fields[k] = v
}

// setStruct stores a structural field value: the canonical JSON string
// serves display and string comparison, while the decoded tree serves dot
// notation — the two views of one ingest-time structure.
func (r *queryResultRow) setStruct(k string, v interface{}) {
	r.set(k, storeValue(v))
	if r.structs == nil {
		r.structs = make(map[string]interface{})
	}
	r.structs[k] = v
}

// isStructure reports whether a value is a structural tree — a map or a
// list — that a field carries in the row's structural view.
func isStructure(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return true
	}
	return false
}

// ordered returns the row's field names in insertion order. Fields written
// directly to the map without set are appended in sorted order as a
// fallback so the result is always deterministic.
func (r *queryResultRow) ordered() []string {
	out := make([]string, 0, len(r.fields))
	seen := make(map[string]bool, len(r.fields))
	for _, k := range r.columns {
		if _, ok := r.fields[k]; !ok || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, k)
	}
	var rest []string
	for k := range r.fields {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	return append(out, rest...)
}

// cloneRow deep-copies a row including its column order and pre-projection
// source. The struct trees are read-only once stored, so the structs map is
// copied shallowly and the trees shared.
func cloneRow(src queryResultRow) queryResultRow {
	fields := make(map[string]string, len(src.fields))
	for k, v := range src.fields {
		fields[k] = v
	}
	columns := make([]string, len(src.columns))
	copy(columns, src.columns)
	var structs map[string]interface{}
	if src.structs != nil {
		structs = make(map[string]interface{}, len(src.structs))
		for k, v := range src.structs {
			structs[k] = v
		}
	}
	return queryResultRow{fields: fields, columns: columns, base: src.base, structs: structs}
}

// queryStats holds statistics about the query execution.

func buildRows(events []logEventWithContext, accountID string) []queryResultRow {
	rows := make([]queryResultRow, 0, len(events))
	for _, evt := range events {
		row := queryResultRow{}
		row.set("@timestamp", formatTimestampField(evt.timestamp))
		row.set("@message", evt.message)
		row.set("@logStream", evt.logStream)
		// @log identifies the event's log group as account:group-name.
		row.set("@log", accountID+":"+evt.logGroup)
		row.set("@ingestionTime", formatTimestampField(evt.ingestionTime))
		// @ptr addresses the event for GetLogRecord; aggregate rows built
		// by later commands never carry it.
		row.set("@ptr", eventPointer(evt.logGroup, evt.logStream, evt.timestamp, evt.message))
		discoverJSONFields(evt.message, &row, evt.transformed)
		rows = append(rows, row)
	}
	return rows
}

// eventPointer encodes the logGroup|logStream|timestamp|message pointer
// that GetLogRecord accepts, in base64. The group and stream fields carry
// the delimiter escaping splitPointer reverses (a log stream name may
// legally contain '|'); the timestamp never does and the message is the
// final raw field.
func eventPointer(group, stream string, ts int64, msg string) string {
	joined := escapePointerField(group) + "|" + escapePointerField(stream) + "|" + strconv.FormatInt(ts, 10) + "|" + msg
	return base64.StdEncoding.EncodeToString([]byte(joined))
}

// escapePointerField escapes the pointer delimiter and the escape
// character itself inside the group and stream fields.
func escapePointerField(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "|", "\\|")
}

// discoverJSONFields merges the top-level keys of a JSON object message
// into the row fields. Discoverable @-fields are never overwritten; field
// names starting with @ are displayed with an additional @ prefix, per the
// documented discovery behaviour. The transformer's @transformationError
// system field is the one exception: in a transformed copy it keeps its
// single @ — "You can query for this field with a query such as filter
// ispresent(@transformationError)". An original message carrying a
// @-prefixed JSON key of the same spelling still renders doubled.
func discoverJSONFields(message string, row *queryResultRow, transformedCopy bool) {
	trimmed := strings.TrimSpace(message)
	if !strings.HasPrefix(trimmed, "{") {
		return
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &decoded); err != nil {
		return
	}
	keys := make([]string, 0, len(decoded))
	for k := range decoded {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	discovered := 0
	for _, k := range keys {
		name := k
		if strings.HasPrefix(k, "@") {
			name = "@" + k
			if transformedCopy && k == transformationErrorField {
				name = k
			}
		}
		if _, exists := row.fields[name]; exists {
			continue
		}
		if discovered >= maxDiscoveredJSONFields {
			return
		}
		// Nested objects and arrays are structure at ingest time — the
		// fields dot notation traverses; scalars stay plain values.
		if isStructure(decoded[k]) {
			row.setStruct(name, decoded[k])
		} else {
			row.set(name, storeValue(decoded[k]))
		}
		discovered++
	}
}

// formatTimestampField renders epoch milliseconds; numeric formatting keeps
// the value numeric for later comparisons.
func formatTimestampField(ms int64) string {
	return formatNumber(float64(ms))
}

// executeQueryContext compiles and runs the query within an execution
// context. It returns the result rows and any compile error.

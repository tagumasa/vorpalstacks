package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// The parser processors: parseJSON, grok, CSV and parseKeyValue —
// the recipes that read the raw message into record keys.

// --- parser processors ---

// applyParseJSON parses a JSON field into the record. With the default
// source (@message) and no destination the parsed keys merge into the
// root ("the new keys are added to the message; the original @message
// content is not changed").
func applyParseJSON(record map[string]interface{}, cfg map[string]interface{}) {
	source := stringMember(cfg, "source", "@message")
	raw, ok := getPath(record, source)
	if !ok {
		return
	}
	text, ok := raw.(string)
	if !ok {
		return
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return
	}
	destination, _ := cfg["destination"].(string)
	if destination == "" {
		for k, v := range parsed {
			if _, exists := record[k]; !exists {
				record[k] = v
			}
		}
		return
	}
	setPath(record, destination, parsed, true)
}

// applyGrok matches the configured grok pattern against the source
// field and adds every named capture. Dotted field names nest.
func applyGrok(record map[string]interface{}, cfg map[string]interface{}) {
	source := stringMember(cfg, "source", "@message")
	match, _ := cfg["match"].(string)
	if match == "" {
		return
	}
	re, ok := compiledGrokFor(match)
	if !ok {
		return
	}
	raw, exists := getPath(record, source)
	if !exists {
		return
	}
	text, ok := raw.(string)
	if !ok {
		return
	}
	compiled := re
	found := compiled.re.FindStringSubmatchIndex(text)
	if found == nil {
		return
	}
	for i, name := range compiled.re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		// An optional capture that did not participate is omitted, not
		// stored as an empty string ("capturing the port as \"port\" if
		// present"): its match bounds are -1, and skipping it also
		// stops a repeated alias's absent occurrence overwriting the
		// value its earlier occurrence captured.
		start, end := found[2*i], found[2*i+1]
		if start < 0 || end < 0 {
			continue
		}
		// Dotted field names are JSON paths (the documentation's grok
		// example 3); the placeholder restores them after the group-name
		// round trip. The common-log composites' numeric groups carry
		// JSON numbers, their httpdate capture converts to the
		// timestamp field, and QUOTEDSTRING captures drop their quotes.
		typed := strings.TrimPrefix(name, grokNumericPrefix)
		field := strings.ReplaceAll(typed, grokDotPlaceholder, ".")
		value := text[start:end]
		if typed != name {
			setPath(record, field, grokNumericCapture(value), true)
			continue
		}
		if converted := strings.TrimPrefix(name, grokTimestampPrefix); converted != name {
			field := strings.ReplaceAll(converted, grokDotPlaceholder, ".")
			setPath(record, field, httpdateToISO(value), true)
			continue
		}
		// A protocol's nil token is the field's absence, not its value:
		// the capture is omitted when it holds the token and written
		// under the stripped name otherwise.
		if nilStripped := strings.TrimPrefix(name, grokNilPrefix); nilStripped != name {
			if value == "-" {
				continue
			}
			field := strings.ReplaceAll(nilStripped, grokDotPlaceholder, ".")
			setPath(record, field, grokCaptureValue(field, value), true)
			continue
		}
		if compiled.quotedFields[name] {
			value = strings.TrimSuffix(strings.TrimPrefix(value, `"`), `"`)
			value = strings.TrimSuffix(strings.TrimPrefix(value, "'"), "'")
		}
		setPath(record, field, grokCaptureValue(field, value), true)
	}
}

// applyCSV parses delimited columns into named (or column_N) keys.
func applyCSV(record map[string]interface{}, cfg map[string]interface{}) {
	source := stringMember(cfg, "source", "@message")
	raw, exists := getPath(record, source)
	if !exists {
		return
	}
	text, ok := raw.(string)
	if !ok {
		return
	}
	delimiter := stringMember(cfg, "delimiter", ",")
	switch delimiter {
	case `\t`:
		delimiter = "\t"
	case `\s`:
		delimiter = " "
	}
	quote := stringMember(cfg, "quoteCharacter", `"`)
	var columns []string
	if rawCols, ok := cfg["columns"].([]interface{}); ok {
		for _, c := range rawCols {
			if s, ok := c.(string); ok {
				columns = append(columns, s)
			}
		}
	}
	parsed := parseCSVLine(text, delimiter, quote)
	if len(parsed) == 0 {
		return
	}
	destination, _ := cfg["destination"].(string)
	for i, value := range parsed {
		name := fmt.Sprintf("column_%d", i+1)
		if i < len(columns) && columns[i] != "" {
			name = columns[i]
		}
		if destination == "" {
			setPath(record, name, value, true)
		} else {
			setPath(record, destination+"."+name, value, true)
		}
	}
}

// parseCSVLine splits one delimited line honouring the quote character
// (quoted segments keep embedded delimiters; quotes are stripped).
func parseCSVLine(text, delimiter, quote string) []string {
	var fields []string
	var current strings.Builder
	inQuotes := false
	runes := []rune(text)
	d := []rune(delimiter)
	q := rune(quote[0])
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == q:
			inQuotes = !inQuotes
		case !inQuotes && matchesDelimiter(runes, i, d):
			fields = append(fields, current.String())
			current.Reset()
			i += len(d) - 1
		default:
			current.WriteRune(r)
		}
	}
	fields = append(fields, current.String())
	return fields
}

func matchesDelimiter(runes []rune, at int, d []rune) bool {
	if at+len(d) > len(runes) {
		return false
	}
	for i, dr := range d {
		if runes[at+i] != dr {
			return false
		}
	}
	return true
}

// applyParseKeyValue splits the source into key-value pairs.
func applyParseKeyValue(record map[string]interface{}, cfg map[string]interface{}) {
	source := stringMember(cfg, "source", "@message")
	raw, exists := getPath(record, source)
	if !exists {
		return
	}
	text, ok := raw.(string)
	if !ok {
		return
	}
	fieldDelimiter := stringMember(cfg, "fieldDelimiter", "&")
	keyValueDelimiter := stringMember(cfg, "keyValueDelimiter", "=")
	nonMatchValue, _ := cfg["nonMatchValue"].(string)
	keyPrefix, _ := cfg["keyPrefix"].(string)
	destination, _ := cfg["destination"].(string)
	overwrite, _ := cfg["overwriteIfExists"].(bool)

	pairs := map[string]interface{}{}
	for _, pair := range strings.Split(text, fieldDelimiter) {
		if pair == "" {
			continue
		}
		key, value, found := strings.Cut(pair, keyValueDelimiter)
		if !found {
			if nonMatchValue == "" {
				continue
			}
			pairs[keyPrefix+key] = nonMatchValue
			continue
		}
		pairs[keyPrefix+key] = value
	}
	if len(pairs) == 0 {
		return
	}
	if destination == "" {
		for k, v := range pairs {
			setPath(record, k, v, overwrite)
		}
		return
	}
	setPath(record, destination, pairs, true)
}

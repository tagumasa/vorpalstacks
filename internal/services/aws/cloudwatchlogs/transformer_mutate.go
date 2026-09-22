package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strings"
)

// The mutate processors: the JSON-mutate family (add/delete/move/
// rename keys, copyValue, listToMap) and the string-mutate family
// (case, trim, splitString, substituteString).

// --- JSON mutate processors ---

func applyAddKeys(record map[string]interface{}, cfg map[string]interface{}) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := entry["key"].(string)
		value, hasValue := entry["value"]
		if key == "" || !hasValue {
			continue
		}
		overwrite, _ := entry["overwriteIfExists"].(bool)
		setPath(record, key, value, overwrite)
	}
}

func applyDeleteKeys(record map[string]interface{}, cfg map[string]interface{}) {
	keys, _ := cfg["withKeys"].([]interface{})
	for _, raw := range keys {
		if key, ok := raw.(string); ok {
			deletePath(record, key)
		}
	}
}

func applyMoveKeys(record map[string]interface{}, cfg map[string]interface{}) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := entry["source"].(string)
		target, _ := entry["target"].(string)
		if source == "" || target == "" {
			continue
		}
		overwrite, _ := entry["overwriteIfExists"].(bool)
		value, exists := getPath(record, source)
		if !exists {
			continue
		}
		// The moved value keeps its leaf name inside the target (the
		// documented example moves outer_key1.inner_key1 to target
		// outer_key2 and the result carries inner_key1 inside it).
		leaf := source[strings.LastIndex(source, ".")+1:]
		// The destination is written before the source is deleted, and
		// the delete happens only when the write landed: a move that
		// cannot overwrite (the destination exists, overwriteIfExists
		// false) leaves the value where it was instead of losing it. A
		// degenerate move onto the source's own path keeps the value.
		dest := target + "." + leaf
		if setPath(record, dest, value, overwrite) && dest != source {
			deletePath(record, source)
		}
	}
}

func applyRenameKeys(record map[string]interface{}, cfg map[string]interface{}) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		key := stringMember(entry, "key", stringMember(entry, "source", ""))
		// The wire member is renameTo (RenameKeyEntry in the service
		// model); the user guide's tables spell it "target", but every
		// SDK and CLI client emits renameTo.
		target := stringMember(entry, "renameTo", "")
		if key == "" || target == "" {
			continue
		}
		overwrite, _ := entry["overwriteIfExists"].(bool)
		value, exists := getPath(record, key)
		if !exists {
			continue
		}
		// Write before delete, and only on a landed write: a rename that
		// cannot overwrite leaves the value under its original key. A
		// degenerate rename onto the same key keeps the value.
		if setPath(record, target, value, overwrite) && target != key {
			deletePath(record, key)
		}
	}
}

// metadataSourceValue resolves copyValue's documented metadata sources —
// "You can also use this processor to add metadata to log events, by
// copying the values of the following metadata keys into the log events:
// @logGroupName, @logGroupStream, @accountId, @regionName" — from the
// event's platform identity. The metadata reading takes precedence over a
// record path of the same spelling: the four keys are designated metadata
// keys wherever they appear as a source.
func metadataSourceValue(tctx TransformContext, source string) (interface{}, bool) {
	switch source {
	case "@logGroupName":
		return tctx.LogGroup, tctx.LogGroup != ""
	case "@logGroupStream":
		return tctx.LogStream, tctx.LogStream != ""
	case "@accountId":
		return tctx.AccountID, tctx.AccountID != ""
	case "@regionName":
		return tctx.Region, tctx.Region != ""
	}
	return nil, false
}

func applyCopyValue(record map[string]interface{}, cfg map[string]interface{}, tctx TransformContext) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source := stringMember(entry, "source", stringMember(entry, "key", ""))
		target, _ := entry["target"].(string)
		if source == "" || target == "" {
			continue
		}
		overwrite, _ := entry["overwriteIfExists"].(bool)
		value, exists := metadataSourceValue(tctx, source)
		if !exists {
			value, exists = getPath(record, source)
		}
		if !exists {
			continue
		}
		setPath(record, target, value, overwrite)
	}
}

func applyListToMap(record map[string]interface{}, cfg map[string]interface{}) {
	source, _ := cfg["source"].(string)
	keyField, _ := cfg["key"].(string)
	if source == "" || keyField == "" {
		return
	}
	valueKey, _ := cfg["valueKey"].(string)
	target, _ := cfg["target"].(string)
	flatten, _ := cfg["flatten"].(bool)
	flattenedElement, _ := cfg["flattenedElement"].(string)

	raw, exists := getPath(record, source)
	if !exists {
		return
	}
	list, ok := raw.([]interface{})
	if !ok {
		return
	}
	result := map[string]interface{}{}
	for _, item := range list {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		k, ok := obj[keyField]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", k)
		var value interface{} = obj
		if valueKey != "" {
			v, ok := obj[valueKey]
			if !ok {
				continue
			}
			value = v
		}
		if flatten && (flattenedElement == "first" || flattenedElement == "last") {
			_, seen := result[key]
			switch flattenedElement {
			case "first":
				if !seen {
					result[key] = value
				}
			case "last":
				result[key] = value
			}
			continue
		}
		existing, isList := result[key].([]interface{})
		if !isList {
			existing = []interface{}{}
		}
		result[key] = append(existing, value)
	}
	if len(result) == 0 {
		return
	}
	if target == "" {
		for k, v := range result {
			setPath(record, k, v, true)
		}
		return
	}
	setPath(record, target, result, true)
}

// --- string mutate processors ---

func applyCaseString(record map[string]interface{}, cfg map[string]interface{}, cast func(string) string) {
	keys, _ := cfg["withKeys"].([]interface{})
	for _, raw := range keys {
		key, ok := raw.(string)
		if !ok {
			continue
		}
		value, exists := getPath(record, key)
		if !exists {
			continue
		}
		text, ok := value.(string)
		if !ok {
			continue
		}
		setPath(record, key, cast(text), true)
	}
}

func applyTrimString(record map[string]interface{}, cfg map[string]interface{}) {
	keys, _ := cfg["withKeys"].([]interface{})
	for _, raw := range keys {
		key, ok := raw.(string)
		if !ok {
			continue
		}
		value, exists := getPath(record, key)
		if !exists {
			continue
		}
		text, ok := value.(string)
		if !ok {
			continue
		}
		setPath(record, key, strings.TrimSpace(text), true)
	}
}

func applySplitString(record map[string]interface{}, cfg map[string]interface{}) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := entry["source"].(string)
		delimiter, _ := entry["delimiter"].(string)
		if source == "" || delimiter == "" {
			continue
		}
		value, exists := getPath(record, source)
		if !exists {
			continue
		}
		text, ok := value.(string)
		if !ok {
			continue
		}
		parts := strings.Split(text, delimiter)
		values := make([]interface{}, len(parts))
		for i, p := range parts {
			values[i] = p
		}
		setPath(record, source, values, true)
	}
}

func applySubstituteString(record map[string]interface{}, cfg map[string]interface{}) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := entry["source"].(string)
		from, _ := entry["from"].(string)
		to, _ := entry["to"].(string)
		if source == "" || from == "" {
			continue
		}
		value, exists := getPath(record, source)
		if !exists {
			continue
		}
		text, ok := value.(string)
		if !ok {
			continue
		}
		re, err := regexp.Compile(from)
		if err != nil {
			continue
		}
		// The documented backreference forms are $n and ${name} — the
		// Go replacement vocabulary, applied verbatim.
		setPath(record, source, re.ReplaceAllString(text, to), true)
	}
}

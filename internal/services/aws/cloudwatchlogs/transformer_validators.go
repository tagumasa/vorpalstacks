package cloudwatchlogs

import (
	"fmt"
)

// The transformer's per-processor member validation: one validator per
// processor kind, enforcing the documented configuration tables' rows
// (required members, entry ceilings, member lengths and depths). The
// dispatch lives in validateTransformerProcessorMembers; the shared
// member helpers here carry the processor kind for their error
// identities.

// transformerEntryCeiling enforces an entries member's 1..max length
// trait; a required member (requiredMembers non-empty) also rejects its
// absence.
func transformerEntryCeiling(kind string, cfg map[string]interface{}, member string, max int, required bool) error {
	raw, has := cfg[member].([]interface{})
	if !has {
		if required {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The %s processor requires %s", kind, member), 400)
		}
		return nil
	}
	if len(raw) == 0 || len(raw) > max {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The %s processor's %s accepts 1 to %d entries", kind, member, max), 400)
	}
	return nil
}

// transformerPathMember enforces an optional path member's documented
// length and nested-key-depth bounds.
func transformerPathMember(kind string, cfg map[string]interface{}, member string) error {
	value, has := cfg[member].(string)
	if !has {
		return nil
	}
	if !validTransformerPath(value) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The %s processor's %s accepts 1 to %d characters and at most %d nested key levels",
				kind, member, transformerMaxPathLength, transformerMaxPathDepth), 400)
	}
	return nil
}

// transformerStringMemberBound enforces an optional string member's
// 1..max length bound.
func transformerStringMemberBound(kind string, cfg map[string]interface{}, member string, max int) error {
	value, has := cfg[member].(string)
	if !has {
		return nil
	}
	if len(value) == 0 || len(value) > max {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The %s processor's %s accepts 1 to %d characters", kind, member, max), 400)
	}
	return nil
}

// validateWithKeysMembers enforces the withKeys member's entry ceiling:
// deleteKeys carries the 1-5 length trait ("Maximum entries: 5"); the
// case and trim siblings keep their ten.
func validateWithKeysMembers(kind string, cfg map[string]interface{}, max int) error {
	return transformerEntryCeiling(kind, cfg, "withKeys", max, true)
}

// validateGrokMembers enforces the grok processor's rows: the required
// match expression with its length ceiling and structural rules, and
// the optional source path.
func validateGrokMembers(cfg map[string]interface{}) error {
	match, _ := cfg["match"].(string)
	if match == "" {
		return NewLogsError("InvalidParameterException",
			"The grok processor requires match", 400)
	}
	if len(match) > transformerMaxGrokMatchLen {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The grok processor's match accepts at most %d characters", transformerMaxGrokMatchLen), 400)
	}
	if err := validateGrokMatch(match); err != nil {
		return err
	}
	return transformerPathMember("grok", cfg, "source")
}

// validateParseJSONMembers enforces the parseJSON processor's optional
// source and destination paths.
func validateParseJSONMembers(kind string, cfg map[string]interface{}) error {
	if err := transformerPathMember(kind, cfg, "source"); err != nil {
		return err
	}
	return transformerPathMember(kind, cfg, "destination")
}

// validateCSVMembers enforces the csv processor's rows: the optional
// source and destination paths, the one-character delimiter ("Maximum
// length: 1 unless the value is \t or \s" — the two-character escapes
// are the documented exception), the one-character quoteCharacter, and
// the bounded column-name list over the path alphabet.
func validateCSVMembers(kind string, cfg map[string]interface{}) error {
	if err := validateParseJSONMembers(kind, cfg); err != nil {
		return err
	}
	if delimiter, has := cfg["delimiter"].(string); has {
		if len(delimiter) > 1 && delimiter != `\t` && delimiter != `\s` {
			return NewLogsError("InvalidParameterException",
				"The csv processor's delimiter accepts one character, \\t or \\s", 400)
		}
	}
	if quote, has := cfg["quoteCharacter"].(string); has && len(quote) != 1 {
		return NewLogsError("InvalidParameterException",
			"The csv processor's quoteCharacter accepts one character", 400)
	}
	if columns, has := cfg["columns"].([]interface{}); has {
		if len(columns) > transformerMaxColumns {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The csv processor's columns accept at most %d entries", transformerMaxColumns), 400)
		}
		for _, raw := range columns {
			column, _ := raw.(string)
			if !validTransformerPath(column) {
				return NewLogsError("InvalidParameterException",
					fmt.Sprintf("The csv processor's column names accept 1 to %d characters and at most %d nested key levels",
						transformerMaxPathLength, transformerMaxPathDepth), 400)
			}
		}
	}
	return nil
}

// validateParseKeyValueMembers enforces the parseKeyValue processor's
// rows: the optional source and destination paths and the four bounded
// string members.
func validateParseKeyValueMembers(kind string, cfg map[string]interface{}) error {
	if err := validateParseJSONMembers(kind, cfg); err != nil {
		return err
	}
	for _, member := range []string{"fieldDelimiter", "keyValueDelimiter", "nonMatchValue", "keyPrefix"} {
		if err := transformerStringMemberBound(kind, cfg, member, transformerMaxPathLength); err != nil {
			return err
		}
	}
	return nil
}

// validateAddKeysMembers enforces the addKeys processor's rows: the
// entry ceiling and each entry's key path and required value.
func validateAddKeysMembers(cfg map[string]interface{}) error {
	if err := transformerEntryCeiling("addKeys", cfg, "entries", 5, true); err != nil {
		return err
	}
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := entry["key"].(string)
		if !validTransformerPath(key) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The addKeys processor's key accepts 1 to %d characters and at most %d nested key levels",
					transformerMaxPathLength, transformerMaxPathDepth), 400)
		}
		// value is a required member (AddKeyValue, length 1-256): an
		// entry without it once passed Put and silently added nothing.
		value, has := entry["value"].(string)
		if !has || len(value) == 0 || len(value) > transformerMaxValueLen {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The addKeys processor's value accepts 1 to %d characters", transformerMaxValueLen), 400)
		}
	}
	return nil
}

// validateMoveKeysMembers enforces the moveKeys and copyValue rows: the
// entry ceiling and each entry's required source and target paths.
func validateMoveKeysMembers(kind string, cfg map[string]interface{}) error {
	if err := transformerEntryCeiling(kind, cfg, "entries", 5, true); err != nil {
		return err
	}
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		for _, member := range []string{"source", "target"} {
			value, _ := entry[member].(string)
			if value == "" || !validTransformerPath(value) {
				return NewLogsError("InvalidParameterException",
					fmt.Sprintf("The %s processor's entries require %s as a path of 1 to %d characters and at most %d nested key levels",
						kind, member, transformerMaxPathLength, transformerMaxPathDepth), 400)
			}
		}
	}
	return nil
}

// validateRenameKeysMembers enforces the renameKeys rows: the entry
// ceiling and each entry's required key and renameTo members.
func validateRenameKeysMembers(cfg map[string]interface{}) error {
	if err := transformerEntryCeiling("renameKeys", cfg, "entries", 5, true); err != nil {
		return err
	}
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := entry["key"].(string)
		if key == "" || len(key) > transformerMaxPathLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The renameKeys processor's key accepts 1 to %d characters", transformerMaxPathLength), 400)
		}
		renameTo, _ := entry["renameTo"].(string)
		if renameTo == "" || len(renameTo) > transformerMaxPathLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The renameKeys processor's renameTo accepts 1 to %d characters", transformerMaxPathLength), 400)
		}
	}
	return nil
}

// validateSplitStringMembers enforces the splitString rows: the entry
// ceiling, each entry's required source path and its required
// delimiter.
func validateSplitStringMembers(kind string, cfg map[string]interface{}) error {
	if err := transformerEntryCeiling(kind, cfg, "entries", 10, true); err != nil {
		return err
	}
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := entry["source"].(string)
		if source == "" || !validTransformerPath(source) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The splitString processor's source accepts 1 to %d characters and at most %d nested key levels",
					transformerMaxPathLength, transformerMaxPathDepth), 400)
		}
		if err := requiredEntryString(kind, entry, "delimiter", transformerMaxPathLength); err != nil {
			return err
		}
	}
	return nil
}

// validateSubstituteStringMembers enforces the substituteString rows:
// the entry ceiling, each entry's required source path and its required
// from and to members.
func validateSubstituteStringMembers(kind string, cfg map[string]interface{}) error {
	if err := transformerEntryCeiling(kind, cfg, "entries", 10, true); err != nil {
		return err
	}
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		source, _ := entry["source"].(string)
		if source == "" || !validTransformerPath(source) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The substituteString processor's source accepts 1 to %d characters and at most %d nested key levels",
					transformerMaxPathLength, transformerMaxPathDepth), 400)
		}
		for _, member := range []string{"from", "to"} {
			if err := requiredEntryString(kind, entry, member, transformerMaxPathLength); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateTypeConverterMembers enforces the typeConverter rows: the
// entry ceiling ("Maximum number of 5 items") and each entry's required
// key path and its cast target from the documented four-value set.
func validateTypeConverterMembers(kind string, cfg map[string]interface{}) error {
	if err := transformerEntryCeiling(kind, cfg, "entries", 5, true); err != nil {
		return err
	}
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := entry["key"].(string)
		if key == "" || !validTransformerPath(key) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The typeConverter processor's key accepts 1 to %d characters and at most %d nested key levels",
					transformerMaxPathLength, transformerMaxPathDepth), 400)
		}
		targetType, _ := entry["type"].(string)
		switch targetType {
		case "integer", "double", "string", "boolean":
		default:
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The typeConverter processor's type accepts integer, double, string or boolean, not %q", targetType), 400)
		}
	}
	return nil
}

// validateDateTimeConverterMembers enforces the dateTimeConverter rows:
// the required source, target and locale members with the paths' bounds,
// the 1..5 matchPatterns list, the bounded targetFormat and the
// timezone members' minimum length of one.
func validateDateTimeConverterMembers(kind string, cfg map[string]interface{}) error {
	source, _ := cfg["source"].(string)
	target, _ := cfg["target"].(string)
	locale, _ := cfg["locale"].(string)
	if source == "" || target == "" || locale == "" {
		return NewLogsError("InvalidParameterException",
			"The dateTimeConverter processor requires source, target and locale", 400)
	}
	if !validTransformerPath(source) || !validTransformerPath(target) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The dateTimeConverter processor's source and target accept 1 to %d characters and at most %d nested key levels",
				transformerMaxPathLength, transformerMaxPathDepth), 400)
	}
	patterns, ok := cfg["matchPatterns"].([]interface{})
	if !ok || len(patterns) == 0 {
		return NewLogsError("InvalidParameterException",
			"The dateTimeConverter processor requires matchPatterns", 400)
	}
	if len(patterns) > transformerMaxPatterns {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The dateTimeConverter processor's matchPatterns accept at most %d entries", transformerMaxPatterns), 400)
	}
	if err := transformerStringMemberBound(kind, cfg, "targetFormat", transformerMaxTargetFmt); err != nil {
		return err
	}
	// The timezone members and locale carry "Minimum length:1" alone.
	for _, member := range []string{"sourceTimezone", "targetTimezone"} {
		if v, has := cfg[member].(string); has && v == "" {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The dateTimeConverter processor's %s accepts at least one character", member), 400)
		}
	}
	return nil
}

// validateListToMapMembers enforces the listToMap rows: the required
// source and key members with the source path's bounds, the bounded
// string members, the target path's depth bound, and the
// flatten/flattenedElement pair ("Required when flatten is set to true"
// and its "Value can only be first or last").
func validateListToMapMembers(kind string, cfg map[string]interface{}) error {
	source, _ := cfg["source"].(string)
	key, _ := cfg["key"].(string)
	if source == "" || key == "" {
		return NewLogsError("InvalidParameterException",
			"The listToMap processor requires source and key", 400)
	}
	if !validTransformerPath(source) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The listToMap processor's source accepts 1 to %d characters and at most %d nested key levels",
				transformerMaxPathLength, transformerMaxPathDepth), 400)
	}
	for _, member := range []string{"key", "valueKey", "target"} {
		if err := transformerStringMemberBound(kind, cfg, member, transformerMaxPathLength); err != nil {
			return err
		}
	}
	if target, has := cfg["target"].(string); has && target != "" && !validTransformerPath(target) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The listToMap processor's target accepts at most %d nested key levels",
				transformerMaxPathDepth), 400)
	}
	if flatten, _ := cfg["flatten"].(bool); flatten {
		element, _ := cfg["flattenedElement"].(string)
		if element != "first" && element != "last" {
			return NewLogsError("InvalidParameterException",
				"The listToMap processor's flattenedElement is required when flatten is true and accepts only first or last", 400)
		}
	}
	return nil
}

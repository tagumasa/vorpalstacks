package sqs

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
)

// getQueryListByStem reads a flattened query-wire string list whose Smithy
// member carries a singular xmlName trait: the indexed entries arrive as
// <stem>.N (AttributeName.1, TagKey.1, MessageAttributeName.1) — a key form
// the generic list helpers never try, since they probe only
// <member>.member.N and <member>.N under the member's own name. The JSON
// arm of each member stays the caller's responsibility.
//
// The list is contiguous: a missing or empty slot with entries beyond it is
// a malformed list that rejects the request, never a silent truncation that
// drops the tail.
func getQueryListByStem(params map[string]interface{}, stem string) ([]string, error) {
	var result []string
	for i := 1; ; i++ {
		key := stem + "." + strconv.Itoa(i)
		val := request.GetParamCaseInsensitive(params, key)
		if val == "" {
			if hasQueryListEntriesBeyond(params, stem, i) {
				return nil, ErrInvalidParameterValue
			}
			break
		}
		result = append(result, val)
	}
	return result, nil
}

// hasQueryListEntriesBeyond reports whether any parameter key carries an
// indexed entry of the stem with an index greater than the given one — the
// signature of a gap (or of a present-but-empty value) inside what is
// otherwise a populated list.
func hasQueryListEntriesBeyond(params map[string]interface{}, stem string, after int) bool {
	prefix := strings.ToLower(stem + ".")
	for k := range params {
		lower := strings.ToLower(k)
		if !strings.HasPrefix(lower, prefix) {
			continue
		}
		rest := lower[len(prefix):]
		idxStr := rest
		if dot := strings.IndexByte(rest, '.'); dot >= 0 {
			idxStr = rest[:dot]
		}
		if idx, err := strconv.Atoi(idxStr); err == nil && idx > after {
			return true
		}
	}
	return false
}

// hasQueryEntryMembers reports whether any parameter key carries the given
// query-wire entry prefix. An indexed batch entry whose other members are
// present while its @required Id is not is a malformed entry that rejects
// the request — not a list terminator silently truncating the batch.
func hasQueryEntryMembers(params map[string]interface{}, entryPrefix string) bool {
	lowerPrefix := strings.ToLower(entryPrefix)
	for k := range params {
		if strings.HasPrefix(strings.ToLower(k), lowerPrefix) {
			return true
		}
	}
	return false
}

// getQueryEntriesByStem reads a flattened query-wire map-shaped entry list
// whose Smithy member carries xmlName + xmlFlattened — the entries arrive as
// <stem>.N.<nameField> / <stem>.N.<valueField> (Attribute.N.Name/Value,
// Tag.N.Key/Value) — a key form the generic map helpers never try. The JSON
// arm of each member stays the caller's responsibility.
//
// The contiguity rule of the string lists holds here too: a missing or empty
// slot with entries beyond it, or a value half arriving without its name
// half, is a malformed list that rejects the request, never a silent
// truncation that drops the tail.
func getQueryEntriesByStem(params map[string]interface{}, stem, nameField, valueField string) (map[string]string, error) {
	result := make(map[string]string)
	for i := 1; ; i++ {
		slot := stem + "." + strconv.Itoa(i) + "."
		name := request.GetParamCaseInsensitive(params, slot+nameField)
		if name == "" {
			// The value half is a leaf on these stems (Attribute.N.Value,
			// Tag.N.Value), so both its exact key and any child members
			// (none today) count as a half-present slot.
			if request.GetParamCaseInsensitive(params, slot+valueField) != "" ||
				hasQueryEntryMembers(params, slot+valueField+".") ||
				hasQueryListEntriesBeyond(params, stem, i) {
				return nil, ErrInvalidParameterValue
			}
			break
		}
		result[name] = request.GetParamCaseInsensitive(params, slot+valueField)
	}
	return result, nil
}

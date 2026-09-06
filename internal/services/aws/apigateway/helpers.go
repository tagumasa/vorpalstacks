package apigateway

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// MaxPaginationLimit is the maximum number of items returned per page,
// matching the AWS API Gateway limit of 500.
const MaxPaginationLimit = 500

// DefaultPaginationLimit is the default page size when the caller does not
// supply a limit parameter.
const DefaultPaginationLimit = 25

// ResolvePaginationLimit extracts the "limit" parameter and validates it.
// When the parameter is absent the default page size is returned. An
// explicit value of zero, a negative number, or a value exceeding
// MaxPaginationLimit produces a validation error — matching the AWS API
// Gateway behaviour of returning 400 BadRequest instead of silently
// clamping.
func ResolvePaginationLimit(params map[string]interface{}) (int, error) {
	raw, ok := params["limit"]
	if !ok || raw == nil {
		return DefaultPaginationLimit, nil
	}

	var limit int
	switch v := raw.(type) {
	case int:
		limit = v
	case int32:
		limit = int(v)
	case int64:
		limit = int(v)
	case float64:
		limit = int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, NewBadRequestException("limit must be an integer")
		}
		limit = n
	default:
		return 0, NewBadRequestException("limit must be an integer")
	}

	if limit <= 0 {
		return 0, NewBadRequestException("limit must be greater than 0")
	}
	if limit > MaxPaginationLimit {
		return 0, NewBadRequestException(fmt.Sprintf("limit must not exceed %d", MaxPaginationLimit))
	}
	return limit, nil
}

// PatchOperation represents a JSON Patch operation for updating API Gateway resources.
type PatchOperation struct {
	Op    string
	Path  string
	Value string
	From  string
}

// patchOpSet is the set of operations a patch path admits, transcribed from
// the row of the official Patch Operations table for the target resource.
type patchOpSet uint8

const (
	opAdd patchOpSet = 1 << iota
	opReplace
	opRemove
	opCopy
)

// requirePatchOp enforces the (path, op) matrix of the official Patch
// Operations page (docs.aws.amazon.com/apigateway/latest/api/patch-operations.html),
// which "lists information about the supported patch operations to update
// resources": the op member documentation states that "Attempts to apply an
// unsupported operation on a resource will return an error message", so a
// cell marked Not supported rejects here with the same wording an
// unrecognised path produces.
func requirePatchOp(po PatchOperation, allowed patchOpSet) *ApiGatewayError {
	var op patchOpSet
	switch po.Op {
	case "add":
		op = opAdd
	case "replace":
		op = opReplace
	case "remove":
		op = opRemove
	case "copy":
		op = opCopy
	}
	if allowed&op == 0 {
		return unknownPatchPathError(po)
	}
	return nil
}

// unknownPatchPathError builds the BadRequestException AWS returns for a
// patch operation whose path the target resource does not support — real AWS
// wording: "Invalid patch path '/name' specified for op 'replace'". Every
// patch applier routes unrecognised paths here instead of silently ignoring
// them.
func unknownPatchPathError(po PatchOperation) *ApiGatewayError {
	return NewBadRequestException(fmt.Sprintf("Invalid patch path '%s' specified for op '%s'", po.Path, po.Op))
}

// parsePatchOperations extracts the patchOperations request member. Only
// add, remove, replace, and copy are implemented in the patch appliers; any
// other operation (including the RFC 6902 move and test values that the Op
// enum admits) is rejected with BadRequestException instead of being silently
// skipped — AWS documents that applying an unsupported operation on a
// resource returns an error message. copy is admitted only with a non-empty
// from path, its source member per the PatchOperation definition.
func parsePatchOperations(params map[string]interface{}) ([]PatchOperation, error) {
	var ops []PatchOperation
	patchOps, ok := params["patchOperations"].([]interface{})
	if !ok {
		return ops, nil
	}

	validOps := map[string]bool{
		"add": true, "remove": true, "replace": true, "copy": true,
	}

	for _, op := range patchOps {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}
		po := PatchOperation{}
		if o, ok := opMap["op"].(string); ok {
			po.Op = o
		}
		if !validOps[po.Op] {
			return nil, NewBadRequestException(fmt.Sprintf(
				"unsupported patch operation %q; supported operations are add, remove, replace and copy", po.Op))
		}
		if p, ok := opMap["path"].(string); ok {
			po.Path = p
		}
		if v, ok := opMap["value"].(string); ok {
			po.Value = v
		}
		if f, ok := opMap["from"].(string); ok {
			po.From = f
		}
		if po.Op == "copy" && po.From == "" {
			return nil, NewBadRequestException("copy operation requires a from path")
		}
		ops = append(ops, po)
	}
	return ops, nil
}

func parseInt64(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err
}

func parseInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	return int32(v), err
}

func parseFloat64(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	return v, err
}

// unescapePointerToken reverses the RFC 6901 escape sequences of a single
// JSON Pointer reference token (~1 becomes /, ~0 becomes ~). The PatchOperation
// path member is a JSON Pointer: any slash in a name must travel escaped, so
// every key or value extracted from a path token must be unescaped here.
func unescapePointerToken(token string) string {
	token = strings.ReplaceAll(token, "~1", "/")
	return strings.ReplaceAll(token, "~0", "~")
}

// splitPatchTokens splits a patch path into its raw JSON Pointer tokens:
// exactly one leading "/" is stripped and the remainder splits on "/".
// Tokens keep their wire form — "~1" and "~0" stay escaped and "*" stays a
// literal token — because escaping is consumed per member, never here.
// An empty path yields no tokens; a lone "/" yields the single root token.
func splitPatchTokens(path string) []string {
	if path == "" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}

// methodMapKey derives the map key the {resource_path}/{http_method}-keyed
// maps share (stage methodSettings and usage-plan per-method throttle).
// The official CLI update-stage example output shows the key as the
// addressed pointer token itself — "~1resourceName/GET" — and the all-methods
// wildcard as the literal "*/*"; the Stage methodSettings documentation
// defines the keys as "{resource_path}/{http_method} ... or /*/* for
// overriding all methods in the stage". Resource tokens therefore stay as
// addressed (escaped forms keep their escapes), and the root's single empty
// token renders as "/" so the root method GET keys as "//GET".
// For a correctly escaped address this equals
// apigatewaystore.MethodSettingsKey(resource_path, http_method) — the
// derivation the execution plane looks the settings up with; the invariant
// is pinned by TestMethodMapKeyMatchesStoreDerivation.
func methodMapKey(resourceTokens []string, httpMethod string) string {
	resource := strings.Join(resourceTokens, "/")
	if len(resourceTokens) == 1 && resourceTokens[0] == "" {
		resource = "/"
	}
	return resource + "/" + httpMethod
}

// isIndexToken reports whether a path element token is numeric-index or
// JSON-Pointer append-marker addressing. Numeric indexes appear in no
// official patch table, so list-member patches reject them.
func isIndexToken(token string) bool {
	if token == "" || token == "-" {
		return true
	}
	_, err := strconv.Atoi(token)
	return err == nil
}

// applyMapPatch applies an add/remove/replace patch operation to a
// string-valued map member: the member name is taken from the path after
// prefix and unescaped per RFC 6902 (~1 becomes /, ~0 becomes ~); "remove"
// deletes the entry, any other operation sets it. The documented patch
// forms address a named entry, so a path whose key token is empty (a
// trailing slash) is rejected with the unknown-patch-path error instead
// of writing an empty-string key. validateName and validateValue gate the
// entries where the target member documents constraints and may be nil.
func applyMapPatch(target map[string]string, po PatchOperation, prefix string, validateName, validateValue func(string) bool) error {
	name := unescapePointerToken(strings.TrimPrefix(po.Path, prefix))
	if name == "" {
		return unknownPatchPathError(po)
	}
	if po.Op != "remove" {
		if validateName != nil && !validateName(name) {
			return NewBadRequestException(fmt.Sprintf(
				"Invalid patch path '%s': invalid entry name '%s'", po.Path, name))
		}
		if validateValue != nil && !validateValue(po.Value) {
			return NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': invalid entry value for '%s'", po.Path, name))
		}
	}
	if po.Op == "remove" {
		delete(target, name)
	} else {
		target[name] = po.Value
	}
	return nil
}

// applyBoolMapPatch is applyMapPatch for bool-valued members, where the
// patch value is the string form of the boolean.
func applyBoolMapPatch(target map[string]bool, po PatchOperation, prefix string, validateName, validateValue func(string) bool) error {
	name := unescapePointerToken(strings.TrimPrefix(po.Path, prefix))
	if name == "" {
		return unknownPatchPathError(po)
	}
	if po.Op != "remove" {
		if validateName != nil && !validateName(name) {
			return NewBadRequestException(fmt.Sprintf(
				"Invalid patch path '%s': invalid entry name '%s'", po.Path, name))
		}
		if validateValue != nil && !validateValue(po.Value) {
			return NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': invalid entry value for '%s'", po.Path, name))
		}
	}
	if po.Op == "remove" {
		delete(target, name)
	} else {
		target[name] = po.Value == "true"
	}
	return nil
}

// parseWholeStringMapValue decodes the value carried by a whole-member map
// patch (a path naming the map itself, no key token). The official
// PatchOperation value documentation states that updating a property of a
// JSON value passes the JSON object in the string value, so a whole-map
// replace carries the map as a JSON object string. Entry names must be
// non-empty; validateName and validateValue gate the entries where the
// target member documents constraints and may be nil.
func parseWholeStringMapValue(po PatchOperation, validateName, validateValue func(string) bool) (map[string]string, error) {
	var raw map[string]string
	if err := json.Unmarshal([]byte(po.Value), &raw); err != nil {
		return nil, NewBadRequestException(fmt.Sprintf(
			"Invalid patch value for '%s': expected a JSON object of string to string", po.Path))
	}
	parsed := make(map[string]string, len(raw))
	for name, value := range raw {
		if name == "" {
			return nil, NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': entry names must not be empty", po.Path))
		}
		if validateName != nil && !validateName(name) {
			return nil, NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': invalid entry name '%s'", po.Path, name))
		}
		if validateValue != nil && !validateValue(value) {
			return nil, NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': invalid entry value for '%s'", po.Path, name))
		}
		parsed[name] = value
	}
	return parsed, nil
}

// parseWholeBoolMapValue is parseWholeStringMapValue for bool-valued maps.
func parseWholeBoolMapValue(po PatchOperation) (map[string]bool, error) {
	var raw map[string]bool
	if err := json.Unmarshal([]byte(po.Value), &raw); err != nil {
		return nil, NewBadRequestException(fmt.Sprintf(
			"Invalid patch value for '%s': expected a JSON object of string to boolean", po.Path))
	}
	parsed := make(map[string]bool, len(raw))
	for name := range raw {
		if name == "" {
			return nil, NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': entry names must not be empty", po.Path))
		}
		parsed[name] = raw[name]
	}
	return parsed, nil
}

// parseWholeStringListValue decodes the value carried by a whole-member list
// patch whose row documents replace: following the PatchOperation value
// documentation's JSON single-quote rule and the whole-map replace form (a
// JSON object), the whole-list replace carries a JSON array of strings.
func parseWholeStringListValue(po PatchOperation) ([]string, error) {
	var raw []string
	if err := json.Unmarshal([]byte(po.Value), &raw); err != nil {
		return nil, NewBadRequestException(fmt.Sprintf(
			"Invalid patch value for '%s': expected a JSON array of strings", po.Path))
	}
	for _, item := range raw {
		if item == "" {
			return nil, NewBadRequestException(fmt.Sprintf(
				"Invalid patch value for '%s': entries must not be empty", po.Path))
		}
	}
	return raw, nil
}

// applyWholeStringMapPatch applies a whole-member map patch for the
// members whose official patch table rows allow add, replace and remove:
// add and replace set the map from the JSON object value, remove clears
// it. Any other operation is unsupported for these paths and returns the
// unknown-patch-path error.
func applyWholeStringMapPatch(target *map[string]string, po PatchOperation, validateName, validateValue func(string) bool) error {
	switch po.Op {
	case "remove":
		*target = nil
	case "add", "replace":
		parsed, err := parseWholeStringMapValue(po, validateName, validateValue)
		if err != nil {
			return err
		}
		*target = parsed
	default:
		return unknownPatchPathError(po)
	}
	return nil
}

// applyWholeBoolMapPatch is applyWholeStringMapPatch for bool-valued maps.
func applyWholeBoolMapPatch(target *map[string]bool, po PatchOperation) error {
	switch po.Op {
	case "remove":
		*target = nil
	case "add", "replace":
		parsed, err := parseWholeBoolMapValue(po)
		if err != nil {
			return err
		}
		*target = parsed
	default:
		return unknownPatchPathError(po)
	}
	return nil
}

// paginateItems applies position-based pagination using the "id" key.
func paginateItems(items []interface{}, position string, limit int) ([]interface{}, string, bool) {
	return paginateItemsWithKey(items, position, limit, "id")
}

// paginateItemsWithKey applies position-based pagination using a custom key
// to identify items (e.g. "stageName" for stages which have no "id" field).
// Returns (page, nextPosition, positionFound). positionFound is false when a
// non-empty position token does not match any item.
func paginateItemsWithKey(items []interface{}, position string, limit int, key string) ([]interface{}, string, bool) {
	start := 0
	positionFound := true
	if position != "" {
		positionFound = false
		for i, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				if id, _ := m[key].(string); id == position {
					start = i + 1
					positionFound = true
					break
				}
			}
		}
	}

	end := len(items)
	if limit > 0 && start+limit < end {
		end = start + limit
	}
	page := items[start:end]

	nextPosition := ""
	if end < len(items) && len(page) > 0 {
		if m, ok := page[len(page)-1].(map[string]interface{}); ok {
			nextPosition, _ = m[key].(string)
		}
	}

	return page, nextPosition, positionFound
}

// sortedKeys returns the keys of a map in sorted order for deterministic pagination.
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// paginateAdminList applies offset-based pagination for admin handler list
// operations. The position string is an integer offset (0-based). Returns
// the start index, end index, the next position string (empty if no
// more items), and ok (false when a non-empty position token is invalid
// or out of range).
func paginateAdminList(totalLen int, position string, limit int) (int, int, string, bool) {
	if limit <= 0 {
		limit = DefaultPaginationLimit
	}
	if limit > MaxPaginationLimit {
		limit = MaxPaginationLimit
	}
	start := 0
	if position != "" {
		idx, err := strconv.Atoi(position)
		if err != nil || idx < 0 || idx >= totalLen {
			return 0, 0, "", false
		}
		start = idx
	}
	end := start + limit
	if end > totalLen {
		end = totalLen
	}
	nextPos := ""
	if end < totalLen {
		nextPos = strconv.Itoa(end)
	}
	return start, end, nextPos, true
}

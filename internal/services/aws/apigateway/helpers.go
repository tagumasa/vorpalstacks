package apigateway

import (
	"fmt"
	"sort"
	"strconv"
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
}

// parsePatchOperations extracts the patchOperations request member. Only
// add, remove, and replace are implemented in the patch appliers; any other
// operation (including the RFC 6902 move, copy, and test values that the Op
// enum admits) is rejected with BadRequestException instead of being silently
// skipped — AWS documents that applying an unsupported operation on a
// resource returns an error message.
func parsePatchOperations(params map[string]interface{}) ([]PatchOperation, error) {
	var ops []PatchOperation
	patchOps, ok := params["patchOperations"].([]interface{})
	if !ok {
		return ops, nil
	}

	validOps := map[string]bool{
		"add": true, "remove": true, "replace": true,
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
				"unsupported patch operation %q; supported operations are add, remove and replace", po.Op))
		}
		if p, ok := opMap["path"].(string); ok {
			po.Path = p
		}
		if v, ok := opMap["value"].(string); ok {
			po.Value = v
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

// sliceContains checks whether slice contains the target string.
func sliceContains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
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

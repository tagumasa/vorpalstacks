package apigateway

import (
	"sort"
	"strconv"
)

// PatchOperation represents a JSON Patch operation for updating API Gateway resources.
type PatchOperation struct {
	Op    string
	Path  string
	Value string
}

func parsePatchOperations(params map[string]interface{}) []PatchOperation {
	var ops []PatchOperation
	patchOps, ok := params["patchOperations"].([]interface{})
	if !ok {
		return ops
	}

	for _, op := range patchOps {
		if opMap, ok := op.(map[string]interface{}); ok {
			po := PatchOperation{}
			if o, ok := opMap["op"].(string); ok {
				po.Op = o
			}
			if po.Op == "" {
				po.Op = "replace"
			}
			if p, ok := opMap["path"].(string); ok {
				po.Path = p
			}
			if v, ok := opMap["value"].(string); ok {
				po.Value = v
			}
			ops = append(ops, po)
		}
	}
	return ops
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

// containsAny checks whether slice contains the target string.
func containsAny(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// paginateItems applies position-based pagination to a pre-sorted slice of items.
// Returns the paginated slice and the position value for the next page (empty if no more items).
func paginateItems(items []interface{}, position string, limit int) ([]interface{}, string) {
	start := 0
	if position != "" {
		for i, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				if id, _ := m["id"].(string); id == position {
					start = i + 1
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
			nextPosition, _ = m["id"].(string)
		}
	}

	return page, nextPosition
}

// sortedKeys returns the keys of a map in sorted order for deterministic pagination.
func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

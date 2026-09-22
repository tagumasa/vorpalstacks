package cloudwatchlogs

import (
	"encoding/json"
	"strings"
)

// The JSON function family of the documented Logs Insights QL function
// library. These take structures by their own contracts — the
// structure-null coercion the string, number, and datetime families
// apply does not run for them.

// evalJSONFunctions holds the JSON family's dispatch: false means the
// name is not this family's.
func evalJSONFunctions(name string, f funcArgs) (interface{}, bool) {
	switch name {
	case "jsonparse":
		if decoded, ok := parseJSONValue(f.str(0)); ok {
			return decoded, true
		}
		return nil, true
	case "jsonstringify":
		return asString(f.arg(0)), true
	case "jsonarraysize":
		var decoded interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(f.str(0))), &decoded); err != nil {
			return float64(0), true
		}
		if list, ok := decoded.([]interface{}); ok {
			return float64(len(list)), true
		}
		return float64(0), true
	case "jsonarraycontains":
		var decoded interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(f.str(0))), &decoded); err != nil {
			return false, true
		}
		list, ok := decoded.([]interface{})
		if !ok {
			return false, true
		}
		for _, item := range list {
			if valuesEqual(item, f.arg(1)) {
				return true, true
			}
		}
		return false, true
	}
	return nil, false
}

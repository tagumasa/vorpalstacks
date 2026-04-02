// Package request provides AWS request handling utilities for vorpalstacks.
package request

import (
	"strconv"
	"strings"
)

// HasParam checks whether a parameter key exists in the params map.
func HasParam(params map[string]interface{}, key string) bool {
	if _, ok := params[key]; ok {
		return true
	}
	if _, ok := params[strings.ToLower(key)]; ok {
		return true
	}
	return false
}

// GetStringParam extracts a string parameter from the params map.
func GetStringParam(params map[string]interface{}, key string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	if v, ok := params[strings.ToLower(key)]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetIntParam extracts an integer parameter from the params map.
func GetIntParam(params map[string]interface{}, key string) int {
	if v, ok := params[key]; ok {
		if n, ok := asInt(v); ok {
			return n
		}
	}
	if v, ok := params[strings.ToLower(key)]; ok {
		if n, ok := asInt(v); ok {
			return n
		}
	}
	return 0
}

func asInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	case string:
		parsed, err := strconv.Atoi(n)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

// GetInt64Param extracts a 64-bit integer parameter from the params map.
func GetInt64Param(params map[string]interface{}, key string) int64 {
	if v, ok := params[key]; ok {
		if n, ok := asInt64(v); ok {
			return n
		}
	}
	if v, ok := params[strings.ToLower(key)]; ok {
		if n, ok := asInt64(v); ok {
			return n
		}
	}
	return 0
}

func asInt64(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case float64:
		return int64(n), true
	case string:
		parsed, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

// GetFloatParam extracts a float parameter from the params map.
func GetFloatParam(params map[string]interface{}, key string) float64 {
	val := GetParamLowerFirst(params, key)
	if val == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

// GetBoolParam extracts a boolean parameter from the params map.
func GetBoolParam(params map[string]interface{}, key string) bool {
	if v, ok := params[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
		if s, ok := v.(string); ok {
			return strings.ToLower(s) == "true"
		}
	}
	if v, ok := params[LowerFirst(key)]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
		if s, ok := v.(string); ok {
			return strings.ToLower(s) == "true"
		}
	}
	return false
}

// GetMapParam extracts a map parameter from the params map.
func GetMapParam(params map[string]interface{}, key string) map[string]interface{} {
	if v, ok := params[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// GetParamCaseInsensitive extracts a parameter from the params map using case-insensitive key matching.
func GetParamCaseInsensitive(params map[string]interface{}, key string) string {
	if v := GetStringParam(params, key); v != "" {
		return v
	}
	return GetStringParam(params, strings.ToLower(key))
}

// GetParamLowerFirst extracts a parameter from the params map using lowercased first letter key.
func GetParamLowerFirst(params map[string]interface{}, key string) string {
	if v := GetStringParam(params, key); v != "" {
		return v
	}
	return GetStringParam(params, LowerFirst(key))
}

// LowerFirst converts the first character of a string to lowercase.
func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]|32) + s[1:]
}

// ParseAttributes parses attributes from the params map into a string map.
func ParseAttributes(params map[string]interface{}, key string) map[string]string {
	result := make(map[string]string)
	attrs := GetMapParam(params, key)
	if attrs == nil {
		attrs = GetMapParam(params, LowerFirst(key))
	}
	if attrs == nil {
		attrs = GetMapParam(params, strings.ToLower(key))
	}
	if attrs == nil {
		return result
	}

	for k, v := range attrs {
		if vs, ok := v.(string); ok {
			result[k] = vs
		}
	}
	return result
}

// GetStringList extracts a string list parameter from the params map.
func GetStringList(params map[string]interface{}, key string) []string {
	var result []string
	listIf, ok := params[key]
	if !ok {
		listIf, ok = params[strings.ToLower(key)]
	}
	if !ok {
		lowerFirst := strings.ToLower(string(key[0])) + key[1:]
		listIf, ok = params[lowerFirst]
	}
	if ok {
		if list, ok := listIf.([]interface{}); ok {
			for _, item := range list {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			if len(result) > 0 {
				return result
			}
		} else if list, ok := listIf.([]string); ok {
			result = append(result, list...)
			if len(result) > 0 {
				return result
			}
		}
	}

	for i := 1; ; i++ {
		memberKey := key + ".member." + strconv.Itoa(i)
		val := GetStringParam(params, memberKey)
		if val == "" {
			val = GetStringParam(params, strings.ToLower(memberKey))
		}
		if val == "" {
			break
		}
		result = append(result, val)
	}
	return result
}

// GetMapParamCaseInsensitive extracts a map parameter from the params map using case-insensitive key matching.
func GetMapParamCaseInsensitive(params map[string]interface{}, key string) map[string]interface{} {
	if m := GetMapParam(params, key); m != nil {
		return m
	}
	return GetMapParam(params, strings.ToLower(key))
}

// GetArrayParam extracts an array parameter from the params map.
func GetArrayParam(params map[string]interface{}, key string) []interface{} {
	if v, ok := params[key]; ok {
		if a, ok := v.([]interface{}); ok {
			return a
		}
	}
	if v, ok := params[strings.ToLower(key)]; ok {
		if a, ok := v.([]interface{}); ok {
			return a
		}
	}
	return nil
}

// GetSliceParam extracts a slice parameter from the params map
func GetSliceParam(params map[string]interface{}, key string) []interface{} {
	if v := GetArrayParam(params, key); v != nil {
		return v
	}
	return GetArrayParam(params, LowerFirst(key))
}

// GetArrayParamLowerFirst extracts an array parameter trying original key, LowerFirst, and lowercase variants
func GetArrayParamLowerFirst(params map[string]interface{}, key string) []interface{} {
	if v := GetArrayParam(params, key); v != nil {
		return v
	}
	if v := GetArrayParam(params, LowerFirst(key)); v != nil {
		return v
	}
	return GetArrayParam(params, strings.ToLower(key))
}

// ParseQueryAttributes parses query-style attributes from the params map.
func ParseQueryAttributes(params map[string]interface{}, prefix string) map[string]string {
	result := make(map[string]string)
	for i := 1; ; i++ {
		nameKey := prefix + "." + strconv.Itoa(i) + ".Name"
		valueKey := prefix + "." + strconv.Itoa(i) + ".Value"

		name := GetStringParam(params, nameKey)
		if name == "" {
			name = GetStringParam(params, strings.ToLower(nameKey))
		}
		if name == "" {
			break
		}

		value := GetStringParam(params, valueKey)
		if value == "" {
			value = GetStringParam(params, strings.ToLower(valueKey))
		}

		result[name] = value
	}
	return result
}

// ParseQueryTags parses query-style tags from the params map.
func ParseQueryTags(params map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for i := 1; ; i++ {
		keyKey := "Tag." + strconv.Itoa(i) + ".Key"
		valueKey := "Tag." + strconv.Itoa(i) + ".Value"

		key := GetStringParam(params, keyKey)
		if key == "" {
			key = GetStringParam(params, strings.ToLower(keyKey))
		}
		if key == "" {
			break
		}

		value := GetStringParam(params, valueKey)
		if value == "" {
			value = GetStringParam(params, strings.ToLower(valueKey))
		}

		result[key] = value
	}
	return result
}

// ParseQueryList parses AWS Query protocol list parameters like "Prefix.member.1.Field" into []map[string]interface{}
func ParseQueryList(params map[string]interface{}, prefix string) []map[string]interface{} {
	var result []map[string]interface{}
	entries := make(map[int]map[string]interface{})

	prefixLower := strings.ToLower(prefix)
	memberPrefix := prefix + ".member."
	memberPrefixLower := prefixLower + ".member."

	for key, value := range params {
		var idx int
		var field string

		if strings.HasPrefix(key, memberPrefix) {
			rest := strings.TrimPrefix(key, memberPrefix)
			parts := strings.SplitN(rest, ".", 2)
			if len(parts) == 2 {
				if i, err := strconv.Atoi(parts[0]); err == nil {
					idx = i
					field = parts[1]
				}
			} else if len(parts) == 1 {
				if i, err := strconv.Atoi(parts[0]); err == nil {
					idx = i
					field = "value"
				}
			}
		} else if strings.HasPrefix(key, memberPrefixLower) {
			rest := strings.TrimPrefix(key, memberPrefixLower)
			parts := strings.SplitN(rest, ".", 2)
			if len(parts) == 2 {
				if i, err := strconv.Atoi(parts[0]); err == nil {
					idx = i
					field = parts[1]
				}
			} else if len(parts) == 1 {
				if i, err := strconv.Atoi(parts[0]); err == nil {
					idx = i
					field = "value"
				}
			}
		}

		if idx > 0 && field != "" {
			if entries[idx] == nil {
				entries[idx] = make(map[string]interface{})
			}
			entries[idx][field] = value
		}
	}

	for i := 1; i <= len(entries); i++ {
		if entries[i] != nil {
			result = append(result, entries[i])
		}
	}

	return result
}

// GetListParam gets a list parameter, trying both direct array and AWS Query protocol format
func GetListParam(params map[string]interface{}, key string) []map[string]interface{} {
	if arr := GetArrayParam(params, key); arr != nil {
		result := make([]map[string]interface{}, 0, len(arr))
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
		return result
	}

	if list := ParseQueryList(params, key); len(list) > 0 {
		return list
	}

	return nil
}

// GetListParamLowerFirst gets a list parameter, trying original key, LowerFirst, and lowercase variants
func GetListParamLowerFirst(params map[string]interface{}, key string) []map[string]interface{} {
	if list := GetListParam(params, key); len(list) > 0 {
		return list
	}
	if list := GetListParam(params, LowerFirst(key)); len(list) > 0 {
		return list
	}
	if list := GetListParam(params, strings.ToLower(key)); len(list) > 0 {
		return list
	}
	return nil
}

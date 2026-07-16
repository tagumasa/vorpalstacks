/*
 * Copyright 2026 Vorpalstacks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package filterpattern

import (
	"reflect"
	"strconv"
	"strings"
)

// IsTruthy determines whether a value should be treated as truthy in
// conditional expressions. It handles various Go types including nil,
// booleans, numbers, strings, slices, and maps.
func IsTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case int64:
		return val != 0
	case int32:
		return val != 0
	case float64:
		return val != 0
	case float32:
		return val != 0
	case string:
		return val != ""
	case []any:
		return len(val) > 0
	case map[string]any:
		return len(val) > 0
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map:
			return rv.Len() > 0
		case reflect.Ptr:
			return !rv.IsNil()
		}
		return true
	}
}

// ToFloat64 attempts to convert a value to a float64.
// Returns 0 and false if the conversion is not possible.
func ToFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

// ToString converts a value to its string representation.
func ToString(v any) string {
	return toStringRaw(v)
}

// CompareValues compares two values using the specified operator.
// It handles numeric comparison when both values are numbers,
// otherwise performs string comparison.
func CompareValues(left, right any, op string) bool {
	leftFloat, leftIsNum := ToFloat64(left)
	rightFloat, rightIsNum := ToFloat64(right)

	if leftIsNum && rightIsNum {
		switch op {
		case "==", "=":
			return leftFloat == rightFloat
		case "!=":
			return leftFloat != rightFloat
		case ">":
			return leftFloat > rightFloat
		case "<":
			return leftFloat < rightFloat
		case ">=":
			return leftFloat >= rightFloat
		case "<=":
			return leftFloat <= rightFloat
		}
	}

	switch op {
	case ">=", "<=", ">", "<":
		return false
	}

	leftStr := toStringRaw(left)
	rightStr := toStringRaw(right)

	switch op {
	case "==", "=":
		return leftStr == rightStr
	case "!=":
		return leftStr != rightStr
	}
	return false
}

func toStringRaw(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	default:
		return ""
	}
}

// ContainsValue checks if a haystack contains the needle value.
// Supports slices, strings, and maps as haystacks.
func ContainsValue(haystack, needle any) bool {
	if haystack == nil {
		return false
	}

	switch h := haystack.(type) {
	case []any:
		for _, item := range h {
			if CompareValues(item, needle, "==") {
				return true
			}
		}
	case string:
		n, ok := needle.(string)
		if ok {
			return strings.Contains(h, n)
		}
	case map[string]any:
		n, ok := needle.(string)
		if ok {
			_, exists := h[n]
			return exists
		}
	}
	return false
}

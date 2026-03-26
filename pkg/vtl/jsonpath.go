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

package vtl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regex pattern for matching array index access in JSON paths.
// This pattern matches expressions like "items[0]" to access array elements.
var arrayIndexRegex = regexp.MustCompile(`^(\w+)\[(\d+)\]$`)

// ExtractJSONPath extracts a value from a JSON object using a dot-separated
// path. The path can include array indices using bracket notation.
// Returns the value as a string, or an empty string if the path is not found.
func ExtractJSONPath(obj interface{}, path string) string {
	path = strings.TrimPrefix(path, "$.")
	if path == "" {
		result, _ := json.Marshal(obj)
		return string(result)
	}

	parts := strings.Split(path, ".")
	current := obj
	for _, part := range parts {
		current = navigatePath(current, part)
		if current == nil {
			return ""
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "null"
	default:
		result, _ := json.Marshal(v)
		return string(result)
	}
}

// ExtractJSONPathRaw is similar to ExtractJSONPath but returns the raw
// interface{} value instead of converting it to a string.
func ExtractJSONPathRaw(obj interface{}, path string) interface{} {
	path = strings.TrimPrefix(path, "$.")
	if path == "" {
		return obj
	}

	parts := strings.Split(path, ".")
	current := obj
	for _, part := range parts {
		current = navigatePath(current, part)
		if current == nil {
			return nil
		}
	}
	return current
}

// navigatePath traverses a single part of a JSON path and returns the
// resulting value. It handles both object property access and array index access.
func navigatePath(current interface{}, part string) interface{} {
	if matches := arrayIndexRegex.FindStringSubmatch(part); matches != nil {
		key := matches[1]
		index, _ := strconv.Atoi(matches[2])

		obj, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		arr, ok := obj[key].([]interface{})
		if !ok || index < 0 || index >= len(arr) {
			return nil
		}
		return arr[index]
	}

	switch v := current.(type) {
	case map[string]interface{}:
		return v[part]
	default:
		return nil
	}
}

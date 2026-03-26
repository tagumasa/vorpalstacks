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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Regular expressions for matching VTL utility functions in templates.
// These patterns identify and extract function calls for processing.
var (
	utilUrlEncodeRegex    = regexp.MustCompile(`\$util\.urlEncode\(([^)]+)\)`)
	utilUrlDecodeRegex    = regexp.MustCompile(`\$util\.urlDecode\(([^)]+)\)`)
	utilBase64EncodeRegex = regexp.MustCompile(`\$util\.base64Encode\(([^)]+)\)`)
	utilBase64DecodeRegex = regexp.MustCompile(`\$util\.base64Decode\(([^)]+)\)`)
	utilParseJsonRegex    = regexp.MustCompile(`\$util\.parseJson\(([^)]+)\)`)
	utilEscapeJsRegex     = regexp.MustCompile(`\$util\.escapeJavaScript\(([^)]+)\)`)
)

// processUtil is the main entry point for processing all utility functions
// in a template. It applies each utility processor in sequence.
func (e *Engine) processUtil(templateStr string) string {
	result := templateStr
	result = e.processUtilUrlEncode(result)
	result = e.processUtilUrlDecode(result)
	result = e.processUtilBase64Encode(result)
	result = e.processUtilBase64Decode(result)
	result = e.processUtilParseJson(result)
	result = e.processUtilEscapeJavaScript(result)
	return result
}

// processUtilUrlEncode handles the $util.urlEncode() function which URL-encodes
// a string. The encoded string is safe for use in URLs.
func (e *Engine) processUtilUrlEncode(templateStr string) string {
	return utilUrlEncodeRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		value := utilUrlEncodeRegex.FindStringSubmatch(match)[1]
		value = strings.Trim(value, `"'`)
		return url.QueryEscape(value)
	})
}

// processUtilUrlDecode handles the $util.urlDecode() function which URL-decodes
// a previously encoded string.
func (e *Engine) processUtilUrlDecode(templateStr string) string {
	return utilUrlDecodeRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		value := utilUrlDecodeRegex.FindStringSubmatch(match)[1]
		value = strings.Trim(value, `"'`)
		decoded, err := url.QueryUnescape(value)
		if err != nil {
			return value
		}
		return decoded
	})
}

// processUtilBase64Encode handles the $util.base64Encode() function which
// encodes a string to base64 format.
func (e *Engine) processUtilBase64Encode(templateStr string) string {
	return utilBase64EncodeRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		value := utilBase64EncodeRegex.FindStringSubmatch(match)[1]
		value = strings.Trim(value, `"'`)
		return base64.StdEncoding.EncodeToString([]byte(value))
	})
}

// processUtilBase64Decode handles the $util.base64Decode() function which
// decodes a base64-encoded string.
func (e *Engine) processUtilBase64Decode(templateStr string) string {
	return utilBase64DecodeRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		value := utilBase64DecodeRegex.FindStringSubmatch(match)[1]
		value = strings.Trim(value, `"'`)
		decoded, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return value
		}
		return string(decoded)
	})
}

// processUtilParseJson handles the $util.parseJson() function which parses
// a JSON string and returns the parsed result as a JSON string.
func (e *Engine) processUtilParseJson(templateStr string) string {
	return utilParseJsonRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		value := utilParseJsonRegex.FindStringSubmatch(match)[1]
		value = strings.Trim(value, `"'`)
		var parsed interface{}
		if err := json.Unmarshal([]byte(value), &parsed); err != nil {
			return value
		}
		jsonBytes, _ := json.Marshal(parsed)
		return string(jsonBytes)
	})
}

// processUtilEscapeJavaScript handles the $util.escapeJavaScript() function
// which escapes special characters for safe use in JavaScript strings.
func (e *Engine) processUtilEscapeJavaScript(templateStr string) string {
	return utilEscapeJsRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		value := utilEscapeJsRegex.FindStringSubmatch(match)[1]
		value = strings.Trim(value, `"'`)
		return escapeJavaScript(value)
	})
}

// escapeJavaScript escapes special characters in a string for safe inclusion
// in JavaScript code. This includes quotes, backslashes, newlines, and other
// control characters.
func escapeJavaScript(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			result.WriteString("\\\\")
		case '"':
			result.WriteString("\\\"")
		case '\'':
			result.WriteString("\\'")
		case '\n':
			result.WriteString("\\n")
		case '\r':
			result.WriteString("\\r")
		case '\t':
			result.WriteString("\\t")
		case '\b':
			result.WriteString("\\b")
		case '\f':
			result.WriteString("\\f")
		default:
			if r < 32 || r > 126 {
				result.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}

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
	"regexp"
	"strconv"
	"strings"
)

// SplitLogical splits a logical expression by the given operator (AND/OR),
// while respecting nested parentheses, brackets, braces, and quoted strings.
// This ensures that operators inside strings or nested structures are not
// treated as top-level logical operators.
func SplitLogical(expr, op string) []string {
	var parts []string
	var cur strings.Builder
	depth := 0
	inSingle := false
	inDouble := false
	escaped := false

	opLen := len(op)
	i := 0
	for i < len(expr) {
		c := expr[i]

		if escaped {
			escaped = false
			cur.WriteByte(c)
			i++
			continue
		}

		if c == '\\' {
			escaped = true
			cur.WriteByte(c)
			i++
			continue
		}

		if c == '\'' && !inDouble {
			inSingle = !inSingle
			cur.WriteByte(c)
			i++
			continue
		}

		if c == '"' && !inSingle {
			inDouble = !inDouble
			cur.WriteByte(c)
			i++
			continue
		}

		if inSingle || inDouble {
			cur.WriteByte(c)
			i++
			continue
		}

		if c == '(' || c == '[' || c == '{' {
			depth++
			cur.WriteByte(c)
			i++
			continue
		}

		if c == ')' || c == ']' || c == '}' {
			depth--
			cur.WriteByte(c)
			i++
			continue
		}

		if depth == 0 && i+opLen <= len(expr) && strings.EqualFold(expr[i:i+opLen], op) {
			parts = append(parts, strings.TrimSpace(cur.String()))
			cur.Reset()
			i += opLen
			continue
		}

		cur.WriteByte(c)
		i++
	}

	if cur.Len() > 0 {
		parts = append(parts, strings.TrimSpace(cur.String()))
	}

	return parts
}

// FindComparisonOp finds the position of a comparison operator in an expression.
// It searches for operators like ==, !=, >=, <=, >, <, and = while respecting
// quotes and escape sequences. Returns -1 if no operator is found.
func FindComparisonOp(expr string) int {
	ops := []string{"==", "!=", ">=", "<="}
	for _, op := range ops {
		if idx := findOp(expr, op); idx >= 0 {
			return idx
		}
	}
	for _, op := range []string{">", "<", "="} {
		if idx := findSingleOp(expr, op[0]); idx >= 0 {
			return idx
		}
	}
	return -1
}

func findOp(expr, op string) int {
	inSingle := false
	inDouble := false
	escaped := false

	for i := 0; i < len(expr)-1; i++ {
		c := expr[i]

		if escaped {
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if c == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if inSingle || inDouble {
			continue
		}
		if expr[i:i+2] == op {
			return i
		}
	}
	return -1
}

func findSingleOp(expr string, op byte) int {
	inSingle := false
	inDouble := false
	escaped := false

	for i := 0; i < len(expr); i++ {
		c := expr[i]

		if escaped {
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if c == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if inSingle || inDouble {
			continue
		}
		if c == op {
			if i > 0 && (expr[i-1] == '!' || expr[i-1] == '>' || expr[i-1] == '<' || expr[i-1] == '=') {
				continue
			}
			if i+1 < len(expr) && expr[i+1] == '=' {
				continue
			}
			return i
		}
	}
	return -1
}

// DetectPatternType determines the type of filter pattern based on its syntax.
// Returns PatternTypeJSON for JSON object patterns, PatternTypeDelimited for
// array patterns, or PatternTypeUnstructured for plain text patterns.
func DetectPatternType(pattern string) PatternType {
	pattern = strings.TrimSpace(pattern)

	if strings.HasPrefix(pattern, "{") && strings.HasSuffix(pattern, "}") {
		return PatternTypeJSON
	}
	if strings.HasPrefix(pattern, "[") && strings.Contains(pattern, "]") {
		depth := 0
		balanced := false
		for _, r := range pattern {
			if r == '[' {
				depth++
			} else if r == ']' {
				depth--
			}
			if depth == 0 {
				balanced = true
			}
		}
		if balanced && depth == 0 {
			return PatternTypeDelimited
		}
	}
	return PatternTypeUnstructured
}

// ParseStringLiteral extracts the string content from a quoted literal.
// Handles both single and double quotes.
func ParseStringLiteral(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// ParseNumber attempts to parse a string as a numeric value.
// Returns the parsed float64 and a boolean indicating success.
func ParseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if len(s) > 0 && (s[0] == '"' || s[0] == '\'') {
		return 0, false
	}
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil
}

var regexPattern = regexp.MustCompile(`^%(.+)%$`)

// ExtractRegex extracts a regex pattern from a string wrapped in percent signs.
// For example, "%error.*%" returns "error.*" and true.
func ExtractRegex(s string) (string, bool) {
	matches := regexPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}

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
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Matcher is the main engine for evaluating CloudWatch Logs filter patterns.
// It supports unstructured text patterns, JSON patterns, and delimited patterns.
type Matcher struct{}

// NewMatcher creates a new filter pattern matcher instance.
func NewMatcher() *Matcher {
	return &Matcher{}
}

// Matches determines whether a log message matches the given filter pattern.
// This is the main entry point for pattern matching. It automatically detects
// the pattern type and applies the appropriate matching logic.
func (m *Matcher) Matches(pattern, message string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" || pattern == `" "` || pattern == `"\"" "` {
		return true
	}

	pt := DetectPatternType(pattern)

	switch pt {
	case PatternTypeJSON:
		return m.matchJSON(pattern, message)
	case PatternTypeDelimited:
		return m.matchDelimited(pattern, message)
	default:
		return m.matchUnstructured(pattern, message)
	}
}

// matchUnstructured handles plain text or keyword-based filter patterns.
// It supports logical OR (space-separated), quoted phrases, optional terms (?prefix),
// and exclusion terms (-prefix).
func (m *Matcher) matchUnstructured(pattern, message string) bool {
	pattern = strings.TrimSpace(pattern)

	if parts := SplitLogical(pattern, " or "); len(parts) > 1 {
		for _, p := range parts {
			if m.matchUnstructured(p, message) {
				return true
			}
		}
		return false
	}

	if len(pattern) >= 2 && ((pattern[0] == '"' && pattern[len(pattern)-1] == '"') ||
		(pattern[0] == '\'' && pattern[len(pattern)-1] == '\'')) {
		return m.matchSingleTerm(pattern, message)
	}

	parts := splitUnstructuredTerms(pattern)
	if len(parts) <= 1 {
		return m.matchSingleTerm(pattern, message)
	}

	hasOptional := false
	hasNonOptional := false
	for _, p := range parts {
		if strings.HasPrefix(p, "?") {
			hasOptional = true
		} else {
			hasNonOptional = true
		}
	}

	if hasOptional && !hasNonOptional {
		for _, p := range parts {
			if strings.HasPrefix(p, "?") {
				if m.matchSingleTerm(p[1:], message) {
					return true
				}
			}
		}
		return false
	}

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "?") {
			continue
		}
		if strings.HasPrefix(p, "-") {
			if m.matchSingleTerm(p[1:], message) {
				return false
			}
			continue
		}
		if !m.matchSingleTerm(p, message) {
			return false
		}
	}
	return true
}

// splitUnstructuredTerms parses an unstructured pattern into individual terms,
// respecting quoted strings as single terms.
func splitUnstructuredTerms(pattern string) []string {
	var terms []string
	var cur strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(pattern); i++ {
		c := pattern[i]

		if c == '\'' && !inDouble {
			inSingle = !inSingle
			cur.WriteByte(c)
			continue
		}
		if c == '"' && !inSingle {
			inDouble = !inDouble
			cur.WriteByte(c)
			continue
		}

		if inSingle || inDouble {
			cur.WriteByte(c)
			continue
		}

		if c == ' ' {
			if cur.Len() > 0 {
				terms = append(terms, cur.String())
				cur.Reset()
			}
			continue
		}

		cur.WriteByte(c)
	}

	if cur.Len() > 0 {
		terms = append(terms, cur.String())
	}

	return terms
}

// matchSingleTerm matches a single pattern term against a message.
// Handles regex patterns (%pattern%), quoted phrases, and wildcards.
func (m *Matcher) matchSingleTerm(term, message string) bool {
	term = strings.TrimSpace(term)
	if term == "" {
		return true
	}

	if regex, ok := ExtractRegex(term); ok {
		re, err := regexp.Compile(regex)
		if err != nil {
			return false
		}
		return re.MatchString(message)
	}

	if len(term) >= 2 && ((term[0] == '"' && term[len(term)-1] == '"') ||
		(term[0] == '\'' && term[len(term)-1] == '\'')) {
		phrase := term[1 : len(term)-1]
		return strings.Contains(message, phrase)
	}

	if strings.Contains(term, "*") {
		words := strings.Fields(message)
		for _, word := range words {
			if wildcardMatchSimple(term, word) {
				return true
			}
		}
		return false
	}

	return strings.Contains(message, term)
}

func wildcardMatchSimple(pattern, text string) bool {
	if pattern == "*" {
		return true
	}

	starIdx := strings.Index(pattern, "*")
	if starIdx < 0 {
		return pattern == text
	}

	prefix := pattern[:starIdx]
	if !strings.HasPrefix(text, prefix) {
		return false
	}

	text = text[len(prefix):]
	suffix := pattern[starIdx+1:]

	if suffix == "" {
		return true
	}

	if !strings.Contains(suffix, "*") {
		return strings.HasSuffix(text, suffix)
	}

	return wildcardMatchSimple(suffix, text)
}

// matchJSON evaluates a JSON filter pattern against a log message.
// The pattern should be a JSON object with filter conditions.
func (m *Matcher) matchJSON(pattern, message string) bool {
	pattern = strings.TrimSpace(pattern)
	if !strings.HasPrefix(pattern, "{") || !strings.HasSuffix(pattern, "}") {
		return false
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		return false
	}

	inner := strings.TrimSpace(pattern[1 : len(pattern)-1])
	return m.evalJSONExpr(inner, data)
}

func (m *Matcher) evalJSONExpr(expr string, data map[string]any) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true
	}

	if parts := SplitLogical(expr, "||"); len(parts) > 1 {
		for _, p := range parts {
			if m.evalJSONExpr(p, data) {
				return true
			}
		}
		return false
	}

	if parts := SplitLogical(expr, "&&"); len(parts) > 1 {
		for _, p := range parts {
			if !m.evalJSONExpr(p, data) {
				return false
			}
		}
		return true
	}

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return m.evalJSONExpr(expr[1:len(expr)-1], data)
	}

	if idx := strings.Index(strings.ToUpper(expr), " NOT EXISTS"); idx > 0 {
		path := strings.TrimSpace(expr[:idx])
		_, exists := getJSONValue(data, path)
		return !exists
	}

	if idx := strings.Index(strings.ToUpper(expr), " IS NULL"); idx > 0 {
		path := strings.TrimSpace(expr[:idx])
		val, exists := getJSONValue(data, path)
		return exists && val == nil
	}

	if idx := strings.Index(strings.ToUpper(expr), " IS NOT NULL"); idx > 0 {
		path := strings.TrimSpace(expr[:idx])
		val, exists := getJSONValue(data, path)
		return exists && val != nil
	}

	if idx := strings.Index(strings.ToUpper(expr), " IS "); idx > 0 {
		path := strings.TrimSpace(expr[:idx])
		test := strings.TrimSpace(expr[idx+4:])
		val, _ := getJSONValue(data, path)
		return evalIsTest(val, test)
	}

	if idx := FindComparisonOp(expr); idx >= 0 {
		opLen := 1
		if idx+1 < len(expr) && (expr[idx+1] == '=' || expr[idx+1] == '>' || expr[idx+1] == '<') {
			if expr[idx] == '!' || expr[idx] == '>' || expr[idx] == '<' || expr[idx] == '=' {
				opLen = 2
			}
		}
		op := expr[idx : idx+opLen]
		left := strings.TrimSpace(expr[:idx])
		right := strings.TrimSpace(expr[idx+opLen:])

		leftVal, _ := getJSONValue(data, left)
		rightVal := parseJSONValue(right)

		return CompareValues(leftVal, rightVal, op)
	}

	return false
}

func getJSONValue(data map[string]any, path string) (any, bool) {
	path = strings.TrimSpace(path)
	if !strings.HasPrefix(path, "$.") {
		return nil, false
	}
	path = path[2:]

	if path == "" {
		return data, true
	}

	return getNestedValue(data, path)
}

func getNestedValue(data any, path string) (any, bool) {
	if path == "" {
		return data, true
	}

	if strings.HasPrefix(path, "['") || strings.HasPrefix(path, `["`) {
		endQuote := strings.Index(path[2:], path[1:2])
		if endQuote < 0 {
			return nil, false
		}
		key := path[2 : 2+endQuote]
		rest := path[2+endQuote+2:]

		if m, ok := data.(map[string]any); ok {
			if v, exists := m[key]; exists {
				return getNestedValue(v, rest)
			}
		}
		return nil, false
	}

	if strings.HasPrefix(path, "[") {
		endBracket := strings.Index(path, "]")
		if endBracket < 0 {
			return nil, false
		}
		idxStr := path[1:endBracket]
		rest := path[endBracket+1:]

		if idxStr == "*" {
			return nil, false
		}

		if arr, ok := data.([]any); ok {
			var idx int
			for _, c := range idxStr {
				if c >= '0' && c <= '9' {
					idx = idx*10 + int(c-'0')
				} else {
					return nil, false
				}
			}
			if idx < len(arr) {
				return getNestedValue(arr[idx], rest)
			}
		}
		return nil, false
	}

	dotIdx := strings.Index(path, ".")
	bracketIdx := strings.Index(path, "[")

	var key string
	var rest string

	if dotIdx < 0 && bracketIdx < 0 {
		key = path
		rest = ""
	} else if dotIdx < 0 {
		key = path[:bracketIdx]
		rest = path[bracketIdx:]
	} else if bracketIdx < 0 {
		key = path[:dotIdx]
		rest = path[dotIdx+1:]
	} else {
		if dotIdx < bracketIdx {
			key = path[:dotIdx]
			rest = path[dotIdx+1:]
		} else {
			key = path[:bracketIdx]
			rest = path[bracketIdx:]
		}
	}

	if m, ok := data.(map[string]any); ok {
		if v, exists := m[key]; exists {
			return getNestedValue(v, rest)
		}
	}
	return nil, false
}

func parseJSONValue(s string) any {
	s = strings.TrimSpace(s)

	if s == "null" || s == "NULL" {
		return nil
	}
	if s == "true" || s == "TRUE" {
		return true
	}
	if s == "false" || s == "FALSE" {
		return false
	}

	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}

	if regex, ok := ExtractRegex(s); ok {
		return &regexValue{pattern: regex}
	}

	if strings.Contains(s, "*") {
		return &wildcardValue{pattern: s}
	}

	if n, ok := ParseNumber(s); ok {
		return n
	}

	return s
}

type regexValue struct {
	pattern string
}

type wildcardValue struct {
	pattern string
}

func evalIsTest(val any, test string) bool {
	test = strings.TrimSpace(test)
	switch strings.ToUpper(test) {
	case "NULL":
		return val == nil
	case "TRUE":
		return val == true
	case "FALSE":
		return val == false
	}
	return false
}

// matchDelimited evaluates a delimited (space/comma-separated) filter pattern.
// The pattern uses positional fields (w1, w2, etc.) to match specific columns.
// Supports bracket notation with conditions after the closing bracket,
// e.g. [level, ...] = "ERROR", ...
func (m *Matcher) matchDelimited(pattern, message string) bool {
	pattern = strings.TrimSpace(pattern)
	if !strings.HasPrefix(pattern, "[") {
		return false
	}

	bracketEnd := strings.Index(pattern, "]")
	if bracketEnd < 0 {
		return false
	}

	inner := pattern[1:bracketEnd]
	afterBracket := strings.TrimSpace(pattern[bracketEnd+1:])

	fields := parseDelimitedFields(message)
	conditions := splitDelimitedConditions(inner)

	fieldMap := make(map[string]string)
	for i, f := range fields {
		fieldMap[fmt.Sprintf("w%d", i+1)] = f
	}

	namedFieldIdx := 0
	for i, cond := range conditions {
		cond = strings.TrimSpace(cond)
		if cond == "..." {
			continue
		}

		if strings.HasPrefix(cond, "...") {
			cond = strings.TrimSpace(cond[3:])
			if cond == "" {
				continue
			}
		}

		if !evalDelimitedConditionAt(cond, fields, fieldMap, i) {
			return false
		}
		namedFieldIdx = i + 1
	}

	if afterBracket != "" {
		afterConds := splitDelimitedConditions(afterBracket)
		for _, ac := range afterConds {
			ac = strings.TrimSpace(ac)
			if ac == "..." {
				continue
			}
			if strings.HasPrefix(ac, "...") {
				ac = strings.TrimSpace(ac[3:])
				if ac == "" {
					continue
				}
			}
			if !evalDelimitedConditionAt(ac, fields, fieldMap, namedFieldIdx) {
				return false
			}
		}
	}
	return true
}

func evalDelimitedConditionAt(cond string, fields []string, fieldMap map[string]string, pos int) bool {
	cond = strings.TrimSpace(cond)

	if parts := SplitLogical(cond, "||"); len(parts) > 1 {
		for _, p := range parts {
			if evalDelimitedConditionAt(p, fields, fieldMap, pos) {
				return true
			}
		}
		return false
	}

	if parts := SplitLogical(cond, "&&"); len(parts) > 1 {
		for _, p := range parts {
			if !evalDelimitedConditionAt(p, fields, fieldMap, pos) {
				return false
			}
		}
		return true
	}

	var name string
	var op string
	var value string

	if idx := FindComparisonOp(cond); idx >= 0 {
		opLen := 1
		if idx+1 < len(cond) {
			nextChar := cond[idx+1]
			if nextChar == '=' && (cond[idx] == '!' || cond[idx] == '>' || cond[idx] == '<' || cond[idx] == '=') {
				opLen = 2
			}
		}
		name = strings.TrimSpace(cond[:idx])
		op = cond[idx : idx+opLen]
		value = stripQuotes(strings.TrimSpace(cond[idx+opLen:]))
	} else {
		name = cond
	}

	if strings.HasPrefix(name, "w") && len(name) > 1 {
		if fieldVal, ok := fieldMap[name]; ok {
			if op == "" {
				return true
			}
			return compareField(fieldVal, value, op)
		}
	}

	if op != "" {
		for _, f := range fields {
			if compareField(f, value, op) {
				return true
			}
		}
		return false
	}

	return true
}

func parseDelimitedFields(message string) []string {
	var fields []string
	var cur strings.Builder
	inBrackets := 0
	inQuotes := false

	for i := 0; i < len(message); i++ {
		c := message[i]

		if inQuotes {
			if c == '"' {
				inQuotes = false
			} else {
				cur.WriteByte(c)
			}
			continue
		}

		if c == '"' {
			inQuotes = true
			continue
		}

		if c == '[' {
			inBrackets++
			cur.WriteByte(c)
			continue
		}
		if c == ']' {
			inBrackets--
			cur.WriteByte(c)
			continue
		}

		if inBrackets > 0 {
			cur.WriteByte(c)
			continue
		}

		if c == ' ' {
			if cur.Len() > 0 {
				fields = append(fields, cur.String())
				cur.Reset()
			}
			continue
		}

		cur.WriteByte(c)
	}

	if cur.Len() > 0 {
		fields = append(fields, cur.String())
	}

	return fields
}

func splitDelimitedConditions(inner string) []string {
	var conditions []string
	var cur strings.Builder
	depth := 0

	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if c == '(' || c == '[' || c == '{' {
			depth++
		} else if c == ')' || c == ']' || c == '}' {
			depth--
		}

		if c == ',' && depth == 0 {
			conditions = append(conditions, strings.TrimSpace(cur.String()))
			cur.Reset()
		} else {
			cur.WriteByte(c)
		}
	}

	if cur.Len() > 0 {
		conditions = append(conditions, strings.TrimSpace(cur.String()))
	}

	return conditions
}

func compareField(field, value, op string) bool {
	if regex, ok := ExtractRegex(value); ok {
		re, err := regexp.Compile(regex)
		if err != nil {
			return false
		}
		matched := re.MatchString(field)
		if op == "!=" {
			return !matched
		}
		return matched
	}

	if strings.Contains(value, "*") {
		matched := wildcardMatchSimple(value, field)
		if op == "!=" {
			return !matched
		}
		return matched
	}

	return CompareValues(field, value, op)
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func Matches(pattern, message string) bool {
	return NewMatcher().Matches(pattern, message)
}

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
	"testing"
)

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
	}{
		{nil, false},
		{false, false},
		{true, true},
		{0, false},
		{1, true},
		{-1, true},
		{"", false},
		{"hello", true},
		{[]any{}, false},
		{[]any{1}, true},
		{map[string]any{}, false},
		{map[string]any{"a": 1}, true},
	}

	for _, tt := range tests {
		result := IsTruthy(tt.input)
		if result != tt.expected {
			t.Errorf("IsTruthy(%v) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestDetectPatternType(t *testing.T) {
	tests := []struct {
		pattern  string
		expected PatternType
	}{
		{`ERROR`, PatternTypeUnstructured},
		{`{ $.level = "ERROR" }`, PatternTypeJSON},
		{`[ip, user, status]`, PatternTypeDelimited},
		{`"phrase match"`, PatternTypeUnstructured},
	}

	for _, tt := range tests {
		result := DetectPatternType(tt.pattern)
		if result != tt.expected {
			t.Errorf("DetectPatternType(%q) = %v, expected %v", tt.pattern, result, tt.expected)
		}
	}
}

func TestMatchUnstructured(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		pattern  string
		message  string
		expected bool
	}{
		{`ERROR`, `[ERROR 400] BAD REQUEST`, true},
		{`ERROR`, `[INFO 200] OK`, false},
		{`ERROR ARGUMENTS`, `[ERROR 419] MISSING ARGUMENTS`, true},
		{`ERROR ARGUMENTS`, `[ERROR 400] BAD REQUEST`, false},
		{`?ERROR ?INFO`, `[INFO 200] OK`, true},
		{`?ERROR ?INFO`, `[WARN 300] WARNING`, false},
		{`"INTERNAL SERVER ERROR"`, `[ERROR 500] INTERNAL SERVER ERROR`, true},
		{`"INTERNAL SERVER ERROR"`, `[ERROR 500] BAD GATEWAY`, false},
		{`ERROR -ARGUMENTS`, `[ERROR 400] BAD REQUEST`, true},
		{`ERROR -ARGUMENTS`, `[ERROR 419] MISSING ARGUMENTS`, false},
		{`*`, `any message`, true},
		{`4*`, `status 404`, true},
		{`4*`, `status 500`, false},
		{`%[0-9]+%`, `error code 12345`, true},
		{`%[0-9]+%`, `error code abcde`, false},
	}

	for _, tt := range tests {
		result := matcher.Matches(tt.pattern, tt.message)
		if result != tt.expected {
			t.Errorf("Matches(%q, %q) = %v, expected %v", tt.pattern, tt.message, result, tt.expected)
		}
	}
}

func TestMatchJSON(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		pattern  string
		message  string
		expected bool
	}{
		{
			`{ $.eventType = "UpdateTrail" }`,
			`{"eventType": "UpdateTrail", "sourceIPAddress": "1.2.3.4"}`,
			true,
		},
		{
			`{ $.eventType = "DeleteTrail" }`,
			`{"eventType": "UpdateTrail"}`,
			false,
		},
		{
			`{ $.latency > 50 }`,
			`{"latency": 100}`,
			true,
		},
		{
			`{ $.latency > 50 }`,
			`{"latency": 30}`,
			false,
		},
		{
			`{ $.count >= 10 && $.count <= 100 }`,
			`{"count": 50}`,
			true,
		},
		{
			`{ $.level = "ERROR" || $.level = "WARN" }`,
			`{"level": "WARN"}`,
			true,
		},
		{
			`{ $.missing NOT EXISTS }`,
			`{"level": "ERROR"}`,
			true,
		},
		{
			`{ $.level NOT EXISTS }`,
			`{"level": "ERROR"}`,
			false,
		},
		{
			`{ $.value IS NULL }`,
			`{"value": null}`,
			true,
		},
		{
			`{ $.array[0] = "first" }`,
			`{"array": ["first", "second"]}`,
			true,
		},
		{
			`{ $.user.name = "John" }`,
			`{"user": {"name": "John"}}`,
			true,
		},
	}

	for _, tt := range tests {
		result := matcher.Matches(tt.pattern, tt.message)
		if result != tt.expected {
			t.Errorf("Matches(%q, %q) = %v, expected %v", tt.pattern, tt.message, result, tt.expected)
		}
	}
}

func TestMatchJSONPatternComparisons(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		name     string
		pattern  string
		message  string
		expected bool
	}{
		// "You can use any conditional regular expression when creating
		// filter patterns to match terms in JSON log events" — the guide's
		// worked example: { $.eventType = %Trail% } matches "UpdateTrail".
		{"regex matches term", `{ $.eventType = %Trail% }`, `{"eventType": "UpdateTrail"}`, true},
		{"regex anchors exclude", `{ $.eventType = %Trail% }`, `{"eventType": "Deployment"}`, false},
		{"regex missing field", `{ $.eventType = %Trail% }`, `{"other": 1}`, false},
		{"regex numeric field", `{ $.code = %4\d\d% }`, `{"code": 404}`, true},
		{"negated regex", `{ $.eventType != %Trail% }`, `{"eventType": "Deployment"}`, true},
		// "Use the asterisk (*) as a wild card to match text."
		{"wildcard prefix", `{ $.host = web-* }`, `{"host": "web-prod-1"}`, true},
		{"wildcard miss", `{ $.host = web-* }`, `{"host": "api-prod-1"}`, false},
		{"wildcard missing field", `{ $.host = web-* }`, `{"other": 1}`, false},
		// The metric filter guide's any-value form: "{ $.latency = * }
		// metricValue: $.latency".
		{"any-value numeric", `{ $.latency = * }`, `{"latency": 42}`, true},
		{"any-value string", `{ $.latency = * }`, `{"latency": "high"}`, true},
		{"any-value missing field", `{ $.latency = * }`, `{"other": 1}`, false},
		{"any-value null field", `{ $.latency = * }`, `{"latency": null}`, false},
		{"negated any-value", `{ $.latency != * }`, `{"other": 1}`, true},
	}

	for _, tt := range tests {
		if result := matcher.Matches(tt.pattern, tt.message); result != tt.expected {
			t.Errorf("%s: Matches(%q, %q) = %v, expected %v", tt.name, tt.pattern, tt.message, result, tt.expected)
		}
	}
}

// The documented grammar rejects invalid patterns instead of storing
// them as match-nothing filters: unbalanced quotes, unterminated
// delimited/JSON shapes, regex outside the dialect ("Filter patterns
// with regex can only include the following:" alphanumerics, the symbol
// characters ": _ # = @ / ; , -" and the operators "^ $ ? [ ] { } | \ *
// + ." — "%something!% would be rejected since '!' is not supported",
// "The ( and ) operators are not supported", "Multi-byte characters are
// not supported"), non-compiling regex, and more than two regex in a
// delimited or JSON pattern.
func TestValidatePatternRejectsInvalidForms(t *testing.T) {
	invalid := []string{
		`%something!%`,
		`%(alpha|beta)%`,
		`%マルチバイト%`,
		`%[unclosed%`,
		`[w1=%time%`,
		`{ $.eventType = %Trail% `,
		`{ $.a = "unbalanced }`,
		`ERROR 'unterminated`,
		`{ $.a = %r1% && $.b = %r2% && $.c = %r3% }`,
		`[w1=%a%, w2=%b%, w3=%c%]`,
	}
	for _, pattern := range invalid {
		if err := ValidatePattern(pattern); err == nil {
			t.Errorf("ValidatePattern(%q) accepted an invalid pattern", pattern)
		}
	}

	valid := []string{
		``,
		`ERROR`,
		`"quoted phrase"`,
		`%Trail%`,
		`%[a-z]+ \d{2,4}%`,
		`[w1=%time%, w2=level, w3=%code%]`,
		`{ $.eventType = %Trail% }`,
		`{ $.latency = * }`,
		`{ $.a = %r1% && $.b = %r2% && $.c = "50%" }`,
		`50% off`,
	}
	for _, pattern := range valid {
		if err := ValidatePattern(pattern); err != nil {
			t.Errorf("ValidatePattern(%q) rejected a valid pattern: %v", pattern, err)
		}
	}

	if PatternContainsRegex(`{ $.a = "50%" }`) {
		t.Error("a quoted percent sign is a literal, not a regex span")
	}
	if !PatternContainsRegex(`{ $.a = %r1% }`) {
		t.Error("a regex span must register in the census")
	}
}

func TestMatchDelimited(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		pattern  string
		message  string
		expected bool
	}{
		// The worked example's seven-field alignment: each entry addresses
		// its column ("The following filter pattern parses seven fields").
		{
			`[ip, user, username, timestamp, request =*.html*, status_code = 4*, bytes]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			true,
		},
		// A positional condition never matches another column that happens
		// to satisfy it: the status_code column alone decides.
		{
			`[ip, user, username, timestamp, request =*.html*, status_code = 500, bytes]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			false,
		},
		{
			`[w1=ERROR, w2]`,
			`ERROR 09/25/2014 Something went wrong`,
			true,
		},
		{
			`[w1=INFO, w2]`,
			`ERROR 09/25/2014 Something went wrong`,
			false,
		},
		// The ellipsis references the unnamed leading fields; the
		// conditions after it align from the last field backward.
		{
			`[..., request =*.html*, status_code = 4*, bytes]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			true,
		},
		{
			`[..., status_code = 4*]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			false,
		},
		// The compound expression inside the seven-field form: "the value
		// of status_code must be 404 or 410".
		{
			`[ip, user, username, timestamp, request =*.html*, status_code = 404 || status_code = 410, bytes]`,
			`127.0.0.1 Prod frank [timestamp] "request" 410 1534`,
			false,
		},
		{
			`[..., status_code = 404 || status_code = 410, bytes]`,
			`127.0.0.1 Prod frank [timestamp] "request" 410 1534`,
			true,
		},
		// A pattern with more entries than the event's fields cannot
		// match.
		{
			`[ip, user, username, timestamp, request, status_code = 4*, bytes, extra]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			false,
		},
	}

	for _, tt := range tests {
		result := matcher.Matches(tt.pattern, tt.message)
		if result != tt.expected {
			t.Errorf("Matches(%q, %q) = %v, expected %v", tt.pattern, tt.message, result, tt.expected)
		}
	}
}

func TestSplitLogical(t *testing.T) {
	tests := []struct {
		expr     string
		op       string
		expected []string
	}{
		{`a and b`, `and`, []string{`a`, `b`}},
		{`a or b or c`, `or`, []string{`a`, `b`, `c`}},
		{`(a or b) and c`, `and`, []string{`(a or b)`, `c`}},
		{`"hello world" and foo`, `and`, []string{`"hello world"`, `foo`}},
	}

	for _, tt := range tests {
		result := SplitLogical(tt.expr, tt.op)
		if len(result) != len(tt.expected) {
			t.Errorf("SplitLogical(%q, %q) = %v, expected %v", tt.expr, tt.op, result, tt.expected)
			continue
		}
		for i, v := range result {
			if v != tt.expected[i] {
				t.Errorf("SplitLogical(%q, %q)[%d] = %q, expected %q", tt.expr, tt.op, i, v, tt.expected[i])
			}
		}
	}
}

// TestWildcardSelectors pins the documented wildcard forms: "Wildcard
// selector — You can use the JSON wildcard to select any array element
// or any JSON object field", with the page's worked examples
// { $.arrayKey[*] = %val.{2}% } and
// { $.* = %111\.111\.111\.1[0-9]{1,2}% }, the bracket-notation key
// carrying a period, and the documented quotas ("You can only use up to
// one wildcard selector in a property selector"; "You can use up to
// three wildcard selectors in a filter pattern with compound
// expressions").
func TestWildcardSelectors(t *testing.T) {
	matcher := NewMatcher()

	event := `{"eventTime":"2026-09-21","ip":"111.111.111.15","arrayKey":["one","two","three"],"cluster.name":"c"}`

	cases := []struct {
		pattern  string
		expected bool
	}{
		// Any element of the selected array satisfies the comparison.
		{`{ $.arrayKey[*] = "three" }`, true},
		{`{ $.arrayKey[*] = "four" }`, false},
		// The page's worked regex forms.
		{`{ $.arrayKey[*] = %t.{2}% }`, true},
		{`{ $.* = %111\.111\.111\.1[0-9]{1,2}% }`, true},
		{`{ $.* = %222\.222\.222% }`, false},
		// Bracket-notation keys carry their periods.
		{`{ $.['cluster.name'] = "c" }`, true},
		{`{ $.['cluster.name'] = "d" }`, false},
	}
	for _, tc := range cases {
		if got := matcher.Matches(tc.pattern, event); got != tc.expected {
			t.Errorf("Matches(%q) = %v, expected %v", tc.pattern, got, tc.expected)
		}
	}

	// The quotas: two wildcards in one selector reject; four across a
	// compound pattern reject; three across a compound pattern stay
	// legal.
	if err := ValidatePattern(`{ $.*.inner.* = "x" }`); err == nil {
		t.Error("two wildcard selectors in one property selector must reject")
	}
	if err := ValidatePattern(`{ $.a[*] = "x" && $.* = "y" && $.b[*] = "z" && $.c[*] = "w" }`); err == nil {
		t.Error("four wildcard selectors across a compound pattern must reject")
	}
	if err := ValidatePattern(`{ $.a[*] = "x" && $.* = "y" && $.b[*] = "z" }`); err != nil {
		t.Errorf("three wildcard selectors across a compound pattern stay legal: %v", err)
	}
}

// TestUndocumentedGrammarFormsNeverMatch pins the grammar's documented
// extent: unstructured patterns carry no keyword operator (the page's
// only composition rule is the ?-prefixed optional term — "You can't
// combine the question mark ("?") with other filter patterns"), so "or"
// is a literal term like any other; and "The variables IS NOT and EXISTS
// currently aren't supported".
func TestUndocumentedGrammarFormsNeverMatch(t *testing.T) {
	matcher := NewMatcher()

	// "ERROR or Exception" is three AND-combined literal terms: the
	// message must contain the word "or" itself.
	if matcher.Matches(`ERROR or Exception`, "ERROR Exception") {
		t.Error("an unstructured pattern combines its terms with AND; 'or' is a literal term")
	}
	if !matcher.Matches(`ERROR or Exception`, "ERROR or Exception") {
		t.Error("the literal terms all present must match")
	}

	jsonEvent := `{"a":"x"}`
	if matcher.Matches(`{ $.a IS NOT NULL }`, jsonEvent) {
		t.Error("the IS NOT variable is not supported and never matches")
	}
	if matcher.Matches(`{ $.a EXISTS }`, jsonEvent) {
		t.Error("the EXISTS variable is not supported and never matches")
	}
	// The documented forms keep matching.
	if !matcher.Matches(`{ $.a IS NULL }`, `{"a":null}`) {
		t.Error("IS NULL is documented and must match a null field")
	}
	if !matcher.Matches(`{ $.a NOT EXISTS }`, `{"b":1}`) {
		t.Error("NOT EXISTS is documented and must match an absent field")
	}
}

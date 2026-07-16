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

func TestMatchDelimited(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		pattern  string
		message  string
		expected bool
	}{
		{
			`[ip, user, timestamp, request, status_code = 404, bytes]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			true,
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
		{
			`[..., status_code = 4*]`,
			`127.0.0.1 Prod frank [10/Oct/2000:13:25:15 -0700] "GET /index.html HTTP/1.0" 404 1534`,
			true,
		},
		{
			`[status_code = 404 || status_code = 410]`,
			`127.0.0.1 Prod frank [timestamp] "request" 410 1534`,
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

package eventbridge

import "testing"

// The retry budget counts retries after the initial attempt: a
// MaximumRetryAttempts of N permits N+1 total attempts, and zero permits a
// single attempt with no retries.
func TestRetriesExhaustedSemantics(t *testing.T) {
	cases := []struct {
		attempt    int32 // failed attempts observed so far
		maxRetries int32
		exhausted  bool
	}{
		{attempt: 1, maxRetries: 0, exhausted: true},
		{attempt: 1, maxRetries: 1, exhausted: false},
		{attempt: 2, maxRetries: 1, exhausted: true},
		{attempt: 2, maxRetries: 2, exhausted: false},
		{attempt: 3, maxRetries: 2, exhausted: true},
		{attempt: 186, maxRetries: 185, exhausted: true},
	}
	for _, tc := range cases {
		if got := retriesExhausted(tc.attempt, tc.maxRetries); got != tc.exhausted {
			t.Errorf("retriesExhausted(attempt=%d, maxRetries=%d) = %v, want %v",
				tc.attempt, tc.maxRetries, got, tc.exhausted)
		}
	}
}

// Wildcard matching is anchored at both ends: a fragment adjacent to the
// start or end of the pattern must sit at the corresponding end of the
// value, and only interior fragments may float.
func TestMatchWildcardPatternAnchoring(t *testing.T) {
	cases := []struct {
		value   string
		pattern string
		want    bool
	}{
		// Leading fragment anchors at the start of the value.
		{"abcdef", "abc*", true},
		{"xabcx", "abc*", false},
		{"barfoo", "foo*", false},
		// Trailing fragment anchors at the end of the value.
		{"xabc", "*abc", true},
		{"xabcy", "*abc", false},
		{"abc.log", "*.log", true},
		{"abc.log.txt", "*.log", false},
		// Interior fragments float.
		{"xabcx", "*abc*", true},
		{"aXbXc", "a*b*c", true},
		{"foobar", "foo*bar", true},
		{"aXb", "a*b*c", false},
		{"Xabc", "a*b*c", false},
		{"fooba", "foo*bar", false},
		// No wildcard means exact equality.
		{"abc", "abc", true},
		{"xabc", "abc", false},
		{"", "", true},
		{"abc", "", false},
		// A bare wildcard matches anything, including the empty value.
		{"", "*", true},
		{"anything", "*", true},
		// Consecutive wildcards collapse.
		{"aXb", "a**b", true},
		{"ab", "a**b", true},
		{"aXbXcXb", "a*b*b", true},
	}
	for _, tc := range cases {
		if got := matchWildcardPattern(tc.value, tc.pattern); got != tc.want {
			t.Errorf("matchWildcardPattern(%q, %q) = %v, want %v",
				tc.value, tc.pattern, got, tc.want)
		}
	}
}

package policy

import "testing"

func TestWildcardMatch(t *testing.T) {
	tests := []struct {
		s       string
		pattern string
		want    bool
	}{
		{"axb", "*a*b", true},
		{"axb", "a*b", true},
		{"abcde", "a*c*e", true},
		{"iot:GetThing", "iot:*", true},
		{"aws-iot:GetThing", "*Get*", true},
		{"anything", "*", true},
		{"", "*", true},
		{"abc", "abc", true},
		{"abc", "abcd", false},
		{"abc", "ab", false},
		{"abc", "a*c", true},
		{"", "a", false},
		{"abc", "", false},
		{"", "", true},
		{"hello.world", "hello.world", true},
		{"helloworld", "hello.world", false},
		{"iot:GetThing", "iot:Connect", false},
		{"iot:DeleteThing", "iot:*Thing", true},
		{"foo/bar", "foo/bar", true},
		{"foobar", "foo/bar", false},
		{"", "a*", false},
		{"ab", "a*", true},
		{"ba", "*a", true},
	}
	for _, tt := range tests {
		t.Run(tt.s+"~"+tt.pattern, func(t *testing.T) {
			got := wildcardMatch(tt.s, tt.pattern)
			if got != tt.want {
				t.Errorf("wildcardMatch(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
			}
		})
	}
}

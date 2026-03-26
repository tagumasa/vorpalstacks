package parsers

import (
	"testing"
)

func TestDecodeURLValues(t *testing.T) {
	tests := []struct {
		name     string
		values   map[string]string
		expected map[string]string
	}{
		{
			name:     "nil values",
			values:   nil,
			expected: map[string]string{},
		},
		{
			name:     "empty values",
			values:   map[string]string{},
			expected: map[string]string{},
		},
		{
			name: "with URL encoded values",
			values: map[string]string{
				"key1": "value%20with%20spaces",
				"key2": "value%2Fwith%2Fslashes",
			},
			expected: map[string]string{
				"key1": "value with spaces",
				"key2": "value/with/slashes",
			},
		},
		{
			name: "with non-encoded values",
			values: map[string]string{
				"key1": "simplevalue",
				"key2": "another+value",
			},
			expected: map[string]string{
				"key1": "simplevalue",
				"key2": "another value",
			},
		},
		{
			name: "with mixed values",
			values: map[string]string{
				"encoded":    "hello%20world",
				"notencoded": "test",
			},
			expected: map[string]string{
				"encoded":    "hello world",
				"notencoded": "test",
			},
		},
		{
			name: "with special characters",
			values: map[string]string{
				"special": "%21%40%23%24%25%5E%26%2A%28%29",
			},
			expected: map[string]string{
				"special": "!@#$%^&*()",
			},
		},
		{
			name: "with invalid encoding",
			values: map[string]string{
				"key": "%ZZ",
			},
			expected: map[string]string{
				"key": "%ZZ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeURLValues(tt.values)

			if len(result) != len(tt.expected) {
				t.Errorf("DecodeURLValues() returned %d items, want %d", len(result), len(tt.expected))
			}

			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("DecodeURLValues()[%q] = %q, want %q", k, result[k], v)
				}
			}
		})
	}
}

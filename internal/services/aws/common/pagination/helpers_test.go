package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMaxItems(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]interface{}
		defaultVal int
		expected   int
	}{
		{
			name:       "valid MaxItems",
			params:     map[string]interface{}{"MaxItems": 50},
			defaultVal: 100,
			expected:   50,
		},
		{
			name:       "zero MaxItems uses default",
			params:     map[string]interface{}{"MaxItems": 0},
			defaultVal: 100,
			expected:   100,
		},
		{
			name:       "negative MaxItems uses default",
			params:     map[string]interface{}{"MaxItems": -5},
			defaultVal: 100,
			expected:   100,
		},
		{
			name:       "no MaxItems uses default",
			params:     map[string]interface{}{},
			defaultVal: 100,
			expected:   100,
		},
		{
			name:       "MaxItems as string",
			params:     map[string]interface{}{"MaxItems": "25"},
			defaultVal: 100,
			expected:   25,
		},
		{
			name:       "default zero uses DefaultMaxItems constant",
			params:     map[string]interface{}{},
			defaultVal: 0,
			expected:   DefaultMaxItems,
		},
		{
			name:       "negative default uses DefaultMaxItems",
			params:     map[string]interface{}{},
			defaultVal: -1,
			expected:   DefaultMaxItems,
		},
		{
			name:       "large MaxItems",
			params:     map[string]interface{}{"MaxItems": 10000},
			defaultVal: 100,
			expected:   1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMaxItems(tt.params, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMarker(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		expected string
	}{
		{
			name:     "valid marker",
			params:   map[string]interface{}{"Marker": "abc123"},
			expected: "abc123",
		},
		{
			name:     "empty marker",
			params:   map[string]interface{}{"Marker": ""},
			expected: "",
		},
		{
			name:     "no marker parameter",
			params:   map[string]interface{}{},
			expected: "",
		},
		{
			name:     "marker with special characters",
			params:   map[string]interface{}{"Marker": "abc+def/ghi=="},
			expected: "abc+def/ghi==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMarker(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildNextMarker(t *testing.T) {
	tests := []struct {
		name              string
		count             int
		maxItems          int
		lastKey           string
		expectedKey       string
		expectedTruncated bool
	}{
		{
			name:              "exactly at maxItems with key",
			count:             10,
			maxItems:          10,
			lastKey:           "key10",
			expectedKey:       "key10",
			expectedTruncated: true,
		},
		{
			name:              "below maxItems",
			count:             5,
			maxItems:          10,
			lastKey:           "key5",
			expectedKey:       "",
			expectedTruncated: false,
		},
		{
			name:              "at maxItems but no key",
			count:             10,
			maxItems:          10,
			lastKey:           "",
			expectedKey:       "",
			expectedTruncated: false,
		},
		{
			name:              "above maxItems",
			count:             15,
			maxItems:          10,
			lastKey:           "key15",
			expectedKey:       "key15",
			expectedTruncated: true,
		},
		{
			name:              "zero count",
			count:             0,
			maxItems:          10,
			lastKey:           "",
			expectedKey:       "",
			expectedTruncated: false,
		},
		{
			name:              "zero maxItems",
			count:             5,
			maxItems:          0,
			lastKey:           "key5",
			expectedKey:       "key5",
			expectedTruncated: true,
		},
		{
			name:              "count equals maxItems minus one",
			count:             9,
			maxItems:          10,
			lastKey:           "key9",
			expectedKey:       "",
			expectedTruncated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, truncated := BuildNextMarker(tt.count, tt.maxItems, tt.lastKey)
			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, tt.expectedTruncated, truncated)
		})
	}
}

func TestDefaultMaxItemsConstant(t *testing.T) {
	assert.Equal(t, 100, DefaultMaxItems)
	assert.True(t, DefaultMaxItems > 0)
}

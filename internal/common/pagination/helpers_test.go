package pagination

import (
	"fmt"
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

func TestDefaultMaxItemsConstant(t *testing.T) {
	assert.Equal(t, 100, DefaultMaxItems)
	assert.True(t, DefaultMaxItems > 0)
}

func TestResolveMaxItems(t *testing.T) {
	errFor := func(got int) error { return fmt.Errorf("out of range: %d", got) }
	tests := []struct {
		name       string
		value      int
		defaultVal int
		min        int
		max        int
		want       int
		wantErr    bool
	}{
		{"zero maps to default", 0, 25, 1, 50, 25, false},
		{"in range passes through", 10, 25, 1, 50, 10, false},
		{"upper bound passes through", 50, 25, 1, 50, 50, false},
		{"lower bound passes through", 1, 25, 1, 50, 1, false},
		{"above max rejected", 51, 25, 1, 50, 0, true},
		{"negative rejected", -1, 25, 1, 50, 0, true},
		{"below min rejected when min above one", 1, 25, 2, 50, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveMaxItems(tt.value, tt.defaultVal, tt.min, tt.max, errFor)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ResolveMaxItems() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ResolveMaxItems() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestClampMaxItems(t *testing.T) {
	tests := []struct {
		name       string
		value      int
		defaultVal int
		max        int
		want       int
	}{
		{"zero maps to default", 0, 50, 50, 50},
		{"negative maps to default", -3, 50, 50, 50},
		{"in range passes through", 7, 50, 50, 7},
		{"above max clamps", 500, 50, 50, 50},
		{"distinct default and max", 0, 100, 25, 100},
		{"distinct max clamps", 30, 100, 25, 25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClampMaxItems(tt.value, tt.defaultVal, tt.max); got != tt.want {
				t.Errorf("ClampMaxItems() = %d, want %d", got, tt.want)
			}
		})
	}
}

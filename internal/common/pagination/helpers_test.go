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

// PaginateSliceByPosition walks positions with an opaque marker: repeated
// keys never rewind the walk, and an unusable token yields an empty page.
func TestPaginateSliceByPosition(t *testing.T) {
	items := []string{"a", "a", "a", "b"}

	page := PaginateSliceByPosition(items, "", 2)
	if len(page.Items) != 2 || page.Items[0] != "a" || page.Items[1] != "a" || !page.IsTruncated || page.NextMarker != "1" {
		t.Fatalf("first page: got items=%v truncated=%v marker=%q, want [a a] truncated marker 1", page.Items, page.IsTruncated, page.NextMarker)
	}

	// The second page resumes after the served duplicate, not after the
	// first item sharing its key.
	page = PaginateSliceByPosition(items, page.NextMarker, 2)
	if len(page.Items) != 2 || page.Items[0] != "a" || page.Items[1] != "b" || page.IsTruncated || page.NextMarker != "" {
		t.Fatalf("second page: got items=%v truncated=%v marker=%q, want [a b] untruncated", page.Items, page.IsTruncated, page.NextMarker)
	}

	// A page ending exactly at the list boundary reports no continuation.
	page = PaginateSliceByPosition(items, "", 4)
	if len(page.Items) != 4 || page.IsTruncated || page.NextMarker != "" {
		t.Fatalf("whole-list page: got items=%v truncated=%v marker=%q, want all four untruncated", page.Items, page.IsTruncated, page.NextMarker)
	}

	// An unusable token yields an empty page rather than restarting the
	// walk, which would silently re-deliver items.
	page = PaginateSliceByPosition(items, "not-a-number", 2)
	if len(page.Items) != 0 || page.IsTruncated {
		t.Fatalf("invalid marker: got items=%v truncated=%v, want empty untruncated", page.Items, page.IsTruncated)
	}
	page = PaginateSliceByPosition(items, "9", 2)
	if len(page.Items) != 0 || page.IsTruncated {
		t.Fatalf("out-of-range marker: got items=%v truncated=%v, want empty untruncated", page.Items, page.IsTruncated)
	}

	// A marker at the last position leaves nothing to serve.
	page = PaginateSliceByPosition(items, "3", 2)
	if len(page.Items) != 0 || page.IsTruncated || page.NextMarker != "" {
		t.Fatalf("exhausted marker: got items=%v truncated=%v, want empty untruncated", page.Items, page.IsTruncated)
	}

	empty := PaginateSliceByPosition([]string{}, "", 2)
	if len(empty.Items) != 0 || empty.IsTruncated || empty.NextMarker != "" {
		t.Fatalf("empty list: got items=%v truncated=%v, want empty untruncated", empty.Items, empty.IsTruncated)
	}
}

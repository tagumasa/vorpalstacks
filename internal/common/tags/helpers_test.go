package tags

import (
	"testing"

	"vorpalstacks/internal/store/aws/common"

	"github.com/stretchr/testify/assert"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		key      string
		expected []common.Tag
	}{
		{
			name: "single tag",
			params: map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "Environment",
						"Value": "Production",
					},
				},
			},
			key: "Tags",
			expected: []common.Tag{
				{Key: "Environment", Value: "Production"},
			},
		},
		{
			name: "multiple tags",
			params: map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "Environment",
						"Value": "Production",
					},
					map[string]interface{}{
						"Key":   "Team",
						"Value": "Backend",
					},
				},
			},
			key: "Tags",
			expected: []common.Tag{
				{Key: "Environment", Value: "Production"},
				{Key: "Team", Value: "Backend"},
			},
		},
		{
			name:     "no tags parameter",
			params:   map[string]interface{}{},
			key:      "Tags",
			expected: nil,
		},
		{
			name: "empty tags array",
			params: map[string]interface{}{
				"Tags": []interface{}{},
			},
			key:      "Tags",
			expected: nil,
		},
		{
			name: "tag with empty value",
			params: map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "EmptyKey",
						"Value": "",
					},
				},
			},
			key: "Tags",
			expected: []common.Tag{
				{Key: "EmptyKey", Value: ""},
			},
		},
		{
			name: "different key name",
			params: map[string]interface{}{
				"ResourceTags": []interface{}{
					map[string]interface{}{
						"Key":   "Key1",
						"Value": "Value1",
					},
				},
			},
			key: "ResourceTags",
			expected: []common.Tag{
				{Key: "Key1", Value: "Value1"},
			},
		},
		{
			name: "invalid tag format - not a map",
			params: map[string]interface{}{
				"Tags": []interface{}{
					"not a map",
				},
			},
			key:      "Tags",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTags(tt.params, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTagKeys(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		key      string
		expected map[string]bool
	}{
		{
			name: "single tag key",
			params: map[string]interface{}{
				"TagKeys": []interface{}{"Environment"},
			},
			key: "TagKeys",
			expected: map[string]bool{
				"Environment": true,
			},
		},
		{
			name: "multiple tag keys",
			params: map[string]interface{}{
				"TagKeys": []interface{}{"Environment", "Team", "Project"},
			},
			key: "TagKeys",
			expected: map[string]bool{
				"Environment": true,
				"Team":        true,
				"Project":     true,
			},
		},
		{
			name:     "no tag keys parameter",
			params:   map[string]interface{}{},
			key:      "TagKeys",
			expected: map[string]bool{},
		},
		{
			name: "empty tag keys array",
			params: map[string]interface{}{
				"TagKeys": []interface{}{},
			},
			key:      "TagKeys",
			expected: map[string]bool{},
		},
		{
			name: "invalid format - not a string",
			params: map[string]interface{}{
				"TagKeys": []interface{}{123, "valid"},
			},
			key: "TagKeys",
			expected: map[string]bool{
				"valid": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTagKeys(tt.params, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToResponse(t *testing.T) {
	tests := []struct {
		name     string
		tags     []common.Tag
		expected []map[string]interface{}
	}{
		{
			name: "single tag",
			tags: []common.Tag{
				{Key: "Environment", Value: "Production"},
			},
			expected: []map[string]interface{}{
				{"Key": "Environment", "Value": "Production"},
			},
		},
		{
			name:     "empty tags",
			tags:     []common.Tag{},
			expected: nil,
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: nil,
		},
		{
			name: "multiple tags",
			tags: []common.Tag{
				{Key: "Environment", Value: "Production"},
				{Key: "Team", Value: "Backend"},
			},
			expected: []map[string]interface{}{
				{"Key": "Environment", "Value": "Production"},
				{"Key": "Team", "Value": "Backend"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToResponse(tt.tags)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name     string
		existing []common.Tag
		newTags  []common.Tag
		expected []common.Tag
	}{
		{
			name:     "add new tags",
			existing: []common.Tag{{Key: "Env", Value: "Prod"}},
			newTags:  []common.Tag{{Key: "Team", Value: "Backend"}},
			expected: []common.Tag{{Key: "Env", Value: "Prod"}, {Key: "Team", Value: "Backend"}},
		},
		{
			name:     "override existing tag",
			existing: []common.Tag{{Key: "Env", Value: "Prod"}},
			newTags:  []common.Tag{{Key: "Env", Value: "Dev"}},
			expected: []common.Tag{{Key: "Env", Value: "Dev"}},
		},
		{
			name:     "empty existing tags",
			existing: []common.Tag{},
			newTags:  []common.Tag{{Key: "Env", Value: "Prod"}},
			expected: []common.Tag{{Key: "Env", Value: "Prod"}},
		},
		{
			name:     "empty new tags",
			existing: []common.Tag{{Key: "Env", Value: "Prod"}},
			newTags:  []common.Tag{},
			expected: []common.Tag{{Key: "Env", Value: "Prod"}},
		},
		{
			name:     "both empty",
			existing: []common.Tag{},
			newTags:  []common.Tag{},
			expected: []common.Tag{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Apply(tt.existing, tt.newTags)
			assert.Equal(t, len(tt.expected), len(result))
			for _, exp := range tt.expected {
				found := false
				for _, r := range result {
					if r.Key == exp.Key && r.Value == exp.Value {
						found = true
						break
					}
				}
				assert.True(t, found, "expected tag %+v not found in result", exp)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name         string
		tags         []common.Tag
		keysToRemove map[string]bool
		expected     []common.Tag
	}{
		{
			name: "remove single tag",
			tags: []common.Tag{
				{Key: "Env", Value: "Prod"},
				{Key: "Team", Value: "Backend"},
			},
			keysToRemove: map[string]bool{"Env": true},
			expected: []common.Tag{
				{Key: "Team", Value: "Backend"},
			},
		},
		{
			name: "remove multiple tags",
			tags: []common.Tag{
				{Key: "Env", Value: "Prod"},
				{Key: "Team", Value: "Backend"},
				{Key: "Project", Value: "Alpha"},
			},
			keysToRemove: map[string]bool{"Env": true, "Project": true},
			expected: []common.Tag{
				{Key: "Team", Value: "Backend"},
			},
		},
		{
			name: "remove non-existent tag",
			tags: []common.Tag{
				{Key: "Env", Value: "Prod"},
			},
			keysToRemove: map[string]bool{"NonExistent": true},
			expected: []common.Tag{
				{Key: "Env", Value: "Prod"},
			},
		},
		{
			name:         "empty tags",
			tags:         []common.Tag{},
			keysToRemove: map[string]bool{"Env": true},
			expected:     []common.Tag{},
		},
		{
			name: "empty keys to remove",
			tags: []common.Tag{
				{Key: "Env", Value: "Prod"},
			},
			keysToRemove: map[string]bool{},
			expected: []common.Tag{
				{Key: "Env", Value: "Prod"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Remove(tt.tags, tt.keysToRemove)
			assert.Equal(t, tt.expected, result)
		})
	}
}

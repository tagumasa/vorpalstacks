package tags

import (
	"testing"

	"vorpalstacks/internal/utils/aws/types"

	"github.com/stretchr/testify/assert"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		params   map[string]interface{}
		key      string
		expected []types.Tag
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
			expected: []types.Tag{
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
			expected: []types.Tag{
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
			expected: []types.Tag{
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
			expected: []types.Tag{
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
		tags     []types.Tag
		expected []map[string]interface{}
	}{
		{
			name: "single tag",
			tags: []types.Tag{
				{Key: "Environment", Value: "Production"},
			},
			expected: []map[string]interface{}{
				{"Key": "Environment", "Value": "Production"},
			},
		},
		{
			name:     "empty tags",
			tags:     []types.Tag{},
			expected: []map[string]interface{}{},
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: []map[string]interface{}{},
		},
		{
			name: "multiple tags",
			tags: []types.Tag{
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
		existing []types.Tag
		newTags  []types.Tag
		expected []types.Tag
	}{
		{
			name:     "add new tags",
			existing: []types.Tag{{Key: "Env", Value: "Prod"}},
			newTags:  []types.Tag{{Key: "Team", Value: "Backend"}},
			expected: []types.Tag{{Key: "Env", Value: "Prod"}, {Key: "Team", Value: "Backend"}},
		},
		{
			name:     "override existing tag",
			existing: []types.Tag{{Key: "Env", Value: "Prod"}},
			newTags:  []types.Tag{{Key: "Env", Value: "Dev"}},
			expected: []types.Tag{{Key: "Env", Value: "Dev"}},
		},
		{
			name:     "empty existing tags",
			existing: []types.Tag{},
			newTags:  []types.Tag{{Key: "Env", Value: "Prod"}},
			expected: []types.Tag{{Key: "Env", Value: "Prod"}},
		},
		{
			name:     "empty new tags",
			existing: []types.Tag{{Key: "Env", Value: "Prod"}},
			newTags:  []types.Tag{},
			expected: []types.Tag{{Key: "Env", Value: "Prod"}},
		},
		{
			name:     "both empty",
			existing: []types.Tag{},
			newTags:  []types.Tag{},
			expected: []types.Tag{},
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
		tags         []types.Tag
		keysToRemove map[string]bool
		expected     []types.Tag
	}{
		{
			name: "remove single tag",
			tags: []types.Tag{
				{Key: "Env", Value: "Prod"},
				{Key: "Team", Value: "Backend"},
			},
			keysToRemove: map[string]bool{"Env": true},
			expected: []types.Tag{
				{Key: "Team", Value: "Backend"},
			},
		},
		{
			name: "remove multiple tags",
			tags: []types.Tag{
				{Key: "Env", Value: "Prod"},
				{Key: "Team", Value: "Backend"},
				{Key: "Project", Value: "Alpha"},
			},
			keysToRemove: map[string]bool{"Env": true, "Project": true},
			expected: []types.Tag{
				{Key: "Team", Value: "Backend"},
			},
		},
		{
			name: "remove non-existent tag",
			tags: []types.Tag{
				{Key: "Env", Value: "Prod"},
			},
			keysToRemove: map[string]bool{"NonExistent": true},
			expected: []types.Tag{
				{Key: "Env", Value: "Prod"},
			},
		},
		{
			name:         "empty tags",
			tags:         []types.Tag{},
			keysToRemove: map[string]bool{"Env": true},
			expected:     []types.Tag{},
		},
		{
			name: "empty keys to remove",
			tags: []types.Tag{
				{Key: "Env", Value: "Prod"},
			},
			keysToRemove: map[string]bool{},
			expected: []types.Tag{
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

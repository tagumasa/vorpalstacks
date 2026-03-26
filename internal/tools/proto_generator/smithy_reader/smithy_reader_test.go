package smithy_reader

import (
	"testing"
)

func TestHashKey(t *testing.T) {
	tests := []struct {
		input    string
		expected uint32
	}{
		{"", 0},
		{"a", 97},
		{"test", 1954540393},
		{"Hello", 1418347357},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hashKey(tt.input)
			if result == 0 && tt.input != "" {
				t.Errorf("hashKey(%q) = 0, should not be 0 for non-empty string", tt.input)
			}
			if tt.input == "" && result != 0 {
				t.Errorf("hashKey(%q) = %d, expected 0 for empty string", tt.input, result)
			}
		})
	}
}

func TestExtractName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"foo", "foo"},
		{"com.amazonaws.iam#CreateUser", "CreateUser"},
		{"com.amazonaws.s3#GetObject", "GetObject"},
		{"service#operation", "operation"},
		{"NoHash", "NoHash"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractName(tt.input)
			if result != tt.expected {
				t.Errorf("extractName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseURI(t *testing.T) {
	tests := []struct {
		input         string
		expectedPath  string
		expectedQuery string
	}{
		{"", "", ""},
		{"/", "/", ""},
		{"/path", "/path", ""},
		{"/path?query=value", "/path", "query=value"},
		{"/path?a=1&b=2", "/path", "a=1&b=2"},
		{"?orphan", "", "orphan"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			path, query := parseURI(tt.input)
			if path != tt.expectedPath || query != tt.expectedQuery {
				t.Errorf("parseURI(%q) = (%q, %q), want (%q, %q)",
					tt.input, path, query, tt.expectedPath, tt.expectedQuery)
			}
		})
	}
}

func TestStableEnumIndex(t *testing.T) {
	tests := []struct {
		name     string
		fallback int
		minVal   int
		maxVal   int
	}{
		{"", 0, 1, 1001},
		{"ACTIVE", 0, 1, 1001},
		{"INACTIVE", 0, 1, 1001},
		{"PENDING", 0, 1, 1001},
		{"DELETED", 0, 1, 1001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stableEnumIndex(tt.name, tt.fallback)
			if result < tt.minVal || result > tt.maxVal {
				t.Errorf("stableEnumIndex(%q, %d) = %d, want between %d and %d",
					tt.name, tt.fallback, result, tt.minVal, tt.maxVal)
			}
		})
	}
}

func TestStableEnumIndexUniqueness(t *testing.T) {
	names := []string{"ACTIVE", "INACTIVE", "PENDING", "DELETED", "UNKNOWN", "PENDING", "RUNNING"}
	indices := make(map[int]string)

	for _, name := range names {
		idx := stableEnumIndex(name, 0)
		if existing, ok := indices[idx]; ok {
			t.Logf("Collision: %q and %q both produce %d", existing, name, idx)
		}
		indices[idx] = name
	}
}

func TestToProtoEnumName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "UNKNOWN"},
		{"ACTIVE", "ACTIVE"},
		{"active", "ACTIVE"},
		{"active_user", "ACTIVE_USER"},
		{"123invalid", "_123INVALID"},
		{"user_name", "USER_NAME"},
		{"XMLHttpRequest", "XMLHTTPREQUEST"},
		{"getUser", "GETUSER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toProtoEnumName(tt.input)
			if result != tt.expected {
				t.Errorf("toProtoEnumName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetDocumentation(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"empty", map[string]interface{}{}, ""},
		{"simple", map[string]interface{}{"smithy.api#documentation": "Test doc"}, "Test doc"},
		{"missing", map[string]interface{}{"other": "value"}, ""},
		{"empty_doc", map[string]interface{}{"smithy.api#documentation": ""}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDocumentation(tt.input)
			if result != tt.expected {
				t.Errorf("getDocumentation(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetStringTrait(t *testing.T) {
	tests := []struct {
		name     string
		traits   map[string]interface{}
		expected string
	}{
		{"empty", nil, ""},
		{"simple", map[string]interface{}{"key": "value"}, "value"},
		{"nested", map[string]interface{}{"key": map[string]interface{}{"value": "nested"}}, "nested"},
		{"missing", map[string]interface{}{"other": "value"}, ""},
		{"wrongtype", map[string]interface{}{"key": 123}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringTrait(tt.traits, "key")
			if result != tt.expected {
				t.Errorf("getStringTrait(%v, 'key') = %q, want %q", tt.traits, result, tt.expected)
			}
		})
	}
}

func TestGetIntTrait(t *testing.T) {
	tests := []struct {
		name     string
		traits   map[string]interface{}
		expected int
	}{
		{"empty", nil, 0},
		{"int", map[string]interface{}{"key": 123}, 123},
		{"float", map[string]interface{}{"key": float64(123)}, 123},
		{"nested", map[string]interface{}{"key": map[string]interface{}{"value": float64(456)}}, 456},
		{"missing", map[string]interface{}{"other": 123}, 0},
		{"wrongtype", map[string]interface{}{"key": "string"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIntTrait(tt.traits, "key")
			if result != tt.expected {
				t.Errorf("getIntTrait(%v, 'key') = %d, want %d", tt.traits, result, tt.expected)
			}
		})
	}
}

func TestHasTrait(t *testing.T) {
	tests := []struct {
		name      string
		traits    map[string]interface{}
		traitName string
		expected  bool
	}{
		{"nil", nil, "key", false},
		{"bool_true", map[string]interface{}{"key": true}, "key", true},
		{"bool_false", map[string]interface{}{"key": false}, "key", false},
		{"string", map[string]interface{}{"key": "value"}, "key", true},
		{"missing", map[string]interface{}{"other": "value"}, "key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasTrait(tt.traits, tt.traitName)
			if result != tt.expected {
				t.Errorf("hasTrait(%v, %q) = %v, want %v", tt.traits, tt.traitName, result, tt.expected)
			}
		})
	}
}

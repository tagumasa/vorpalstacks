package proto_generator

import (
	"testing"
)

func TestToProtoFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"ID", "id"},
		{"iD", "id"},
		{"eTag", "etag"},
		{"ID123", "id123"},
		{"XMLSummaries", "xmlsummaries"},
		{"sSEKMSKeyid", "ssekmskeyid"},
		{"userName", "username"},
		{"UserName", "username"},
		{"getUser", "getuser"},
		{"HTTPHeader", "httpheader"},
		{"ACL", "acl"},
		{"aCL", "acl"},
		{"getHTTPResponse", "gethttpresponse"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toProtoFieldName(tt.input)
			if result != tt.expected {
				t.Errorf("toProtoFieldName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStableFieldNumber(t *testing.T) {
	tests := []struct {
		input   string
		minVal  int
		maxVal  int
		notZero bool
	}{
		{"", 1, 536870912, true},
		{"id", 1, 536870912, true},
		{"testField", 1, 536870912, true},
		{"Field1", 1, 536870912, true},
		{"Field19000", 1, 536870912, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stableFieldNumber(tt.input)
			if tt.notZero && result < 1 {
				t.Errorf("stableFieldNumber(%q) = %d, should be >= 1", tt.input, result)
			}
			if result < tt.minVal || result > tt.maxVal {
				t.Errorf("stableFieldNumber(%q) = %d, want between %d and %d", tt.input, result, tt.minVal, tt.maxVal)
			}
		})
	}
}

func TestStableFieldNumberUniqueness(t *testing.T) {
	numbers := make(map[int]string)
	fields := []string{"field1", "field2", "field3", "id", "name", "value", "type", "timestamp", "data", "result"}

	for _, f := range fields {
		num := stableFieldNumber(f)
		if existing, ok := numbers[num]; ok {
			t.Logf("Collision: %q and %q both produce %d", existing, f, num)
		}
		numbers[num] = f
	}

	if len(numbers) != len(fields) {
		t.Errorf("Expected unique numbers for each field, got collisions")
	}
}

func TestFnv32a(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"a"},
		{"test"},
		{"hello"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := fnv32a(tt.input)
			if result == 0 && tt.input != "" {
				t.Errorf("fnv32a(%q) = %d, should not be 0", tt.input, result)
			}
			if tt.input == "" && result != 2166136261 {
				t.Errorf("fnv32a(%q) = %d, expected FNV offset basis", tt.input, result)
			}
		})
	}
}

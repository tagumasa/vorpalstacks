package naming

import (
	"strings"
	"testing"
)

func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with namespace", "com.amazonaws.lambda#DeleteFunction", "DeleteFunction"},
		{"with namespace simple", "com.amazonaws.s3#GetObject", "GetObject"},
		{"no namespace", "DeleteFunction", "DeleteFunction"},
		{"empty string", "", ""},
		{"only namespace", "com.amazonaws.lambda#", ""},
		{"edge case single char", "a#b", "b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractServiceName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractServiceName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractServiceNameFromNamespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"standard namespace", "com.amazonaws.lambda", "lambda"},
		{"two parts", "com.amazonaws", "amazonaws"},
		{"single part", "lambda", "lambda"},
		{"empty string", "", ""},
		{"edge case single char", "a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractServiceNameFromNamespace(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractServiceNameFromNamespace(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractNamespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with namespace", "com.amazonaws.lambda#DeleteFunction", "com.amazonaws.lambda"},
		{"with namespace s3", "com.amazonaws.s3#GetObject", "com.amazonaws.s3"},
		{"no namespace", "DeleteFunction", ""},
		{"empty string", "", ""},
		{"only hash", "#DeleteFunction", ""},
		{"edge case single char", "a#b", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractNamespace(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractNamespace(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateRandomHex(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"length 0", 0, true},
		{"length 1", 1, false},
		{"length 16", 16, false},
		{"length 32", 32, false},
		{"length 64", 64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateRandomHex(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if len(result) != tt.length {
					t.Errorf("GenerateRandomHex() length = %d, want %d", len(result), tt.length)
				}
				if tt.length > 0 {
					for _, c := range result {
						if !strings.Contains("0123456789abcdef", string(c)) {
							t.Errorf("GenerateRandomHex() contains invalid char: %c", c)
						}
					}
				}
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"length 0", 0, true},
		{"length 1", 1, false},
		{"length 16", 16, false},
		{"length 32", 32, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateRandomString(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if len(result) != tt.length {
					t.Errorf("GenerateRandomString() length = %d, want %d", len(result), tt.length)
				}
				if tt.length > 0 {
					for _, c := range result {
						if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", string(c)) {
							t.Errorf("GenerateRandomString() contains invalid char: %c", c)
						}
					}
				}
			}
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	id1 := GenerateRequestID()
	id2 := GenerateRequestID()

	if id1 == "" {
		t.Error("GenerateRequestID() should not return empty string")
	}
	if len(id1) != 36 {
		t.Errorf("GenerateRequestID() length = %d, want 36", len(id1))
	}
	if id1 == id2 {
		t.Error("GenerateRequestID() should generate unique IDs")
	}
}

func TestGenerateResourceID(t *testing.T) {
	id := GenerateResourceID("bucket")

	if !strings.HasPrefix(id, "bucket-") {
		t.Errorf("GenerateResourceID() = %v, want prefix bucket-", id)
	}
	if len(id) <= 8 {
		t.Error("GenerateResourceID() should have more than prefix")
	}
}

func TestGenerateShortID(t *testing.T) {
	id := GenerateShortID()

	if len(id) != 8 {
		t.Errorf("GenerateShortID() length = %d, want 8", len(id))
	}
}

func TestGenerateTimestampID(t *testing.T) {
	id := GenerateTimestampID("api")

	if !strings.HasPrefix(id, "api-") {
		t.Errorf("GenerateTimestampID() = %v, want prefix api-", id)
	}
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		t.Errorf("GenerateTimestampID() should have exactly 2 parts, got %d", len(parts))
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()

	if uuid1 == "" {
		t.Error("GenerateUUID() should not return empty string")
	}
	if len(uuid1) != 36 {
		t.Errorf("GenerateUUID() length = %d, want 36", len(uuid1))
	}
	if uuid1 == uuid2 {
		t.Error("GenerateUUID() should generate unique UUIDs")
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid with uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"valid without hyphens", "550e8400e29b41d4a716446655440000", true},
		{"invalid - not UUID", "invalid", false},
		{"invalid - empty", "", false},
		{"invalid - partial", "550e8400-e29b-41d4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUUID(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidUUID(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

package paramvalidation

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringLength(t *testing.T) {
	lengthErr := func(field string, length, min, max int) error {
		return fmt.Errorf("%s length must be between %d and %d (got %d)", field, min, max, length)
	}
	tests := []struct {
		name    string
		value   string
		min     int
		max     int
		wantErr bool
	}{
		{"within range", "abc", 1, 10, false},
		{"at min", "a", 1, 10, false},
		{"at max", "abcdefghij", 1, 10, false},
		{"below min", "", 1, 10, true},
		{"above max", "abcdefghijk", 1, 10, true},
		{"zero min allows empty", "", 0, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StringLength("Field", tt.value, tt.min, tt.max, lengthErr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("StringLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnumValue(t *testing.T) {
	enumErr := func(field, value string) error {
		return fmt.Errorf("Invalid %s: %s", field, value)
	}
	valid := map[string]bool{"A": true, "B": true}
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"member passes", "A", false},
		{"empty is unset and passes", "", false},
		{"non-member fails", "C", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnumValue("field", tt.value, valid, enumErr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("EnumValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnumList(t *testing.T) {
	enumErr := func(field, value string) error {
		return fmt.Errorf("Invalid %s: %s", field, value)
	}
	valid := map[string]bool{"A": true, "B": true}
	tests := []struct {
		name    string
		values  []string
		wantErr bool
	}{
		{"empty list passes", nil, false},
		{"all members pass", []string{"A", "B", "A"}, false},
		{"first offender reported", []string{"A", "C", "D"}, true},
		{"empty member fails", []string{"A", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnumList("field", tt.values, valid, enumErr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("EnumList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestStringLengthUnicodeCharacters pins that length bounds are counted
// in Unicode characters, matching Smithy @length semantics.
func TestStringLengthUnicodeCharacters(t *testing.T) {
	cjk := strings.Repeat("\u65e5", 128) // 128 characters, 384 bytes
	if err := StringLength("field", cjk, 1, 128, func(f string, n, min, max int) error {
		return nil
	}); err != nil {
		t.Errorf("128-character CJK value rejected: %v", err)
	}
	cjk129 := strings.Repeat("\u65e5", 129)
	err := StringLength("field", cjk129, 1, 128, func(f string, n, min, max int) error {
		return fmt.Errorf("length %d out of range", n)
	})
	if err == nil {
		t.Errorf("129-character CJK value accepted, want rejection")
	}
}

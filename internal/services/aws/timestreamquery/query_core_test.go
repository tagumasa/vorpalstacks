package timestreamquery

import (
	"strings"
	"testing"
)

// TestParseMaxRows pins the MaxRows wire contract: the member targets the
// integer MaxQueryResults shape, so a non-integer value is rejected with
// SerializationException instead of silently falling back to the default
// page size; an absent member keeps the default and an out-of-range
// integer is a ValidationException.
func TestParseMaxRows(t *testing.T) {
	rows, err := parseMaxRows("")
	if err != nil || rows != maxQueryRows {
		t.Fatalf("absent MaxRows must keep the default %d, got (%d, %v)", maxQueryRows, rows, err)
	}
	rows, err = parseMaxRows("5")
	if err != nil || rows != 5 {
		t.Fatalf("integer MaxRows must parse, got (%d, %v)", rows, err)
	}

	_, err = parseMaxRows("not-a-number")
	if err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Fatalf("non-integer MaxRows must be rejected with SerializationException, got %v", err)
	}

	_, err = parseMaxRows("0")
	if err == nil || !strings.Contains(err.Error(), "ValidationException") {
		t.Fatalf("out-of-range integer MaxRows must be rejected with ValidationException, got %v", err)
	}
}

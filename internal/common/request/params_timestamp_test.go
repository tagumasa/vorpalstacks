package request

import (
	"errors"
	"math"
	"testing"
)

// TestNormalizeTimestampValueAcceptsWireForms pins the canonical-string
// conversion across every numeric wire surface: the epoch-second double
// the JSON protocols carry, the unsigned integers the CBOR path produces
// (asInt accepts them; the timestamp read must not falsely reject a
// well-formed CBOR member), the signed integers, and both documented
// string notations.
func TestNormalizeTimestampValueAcceptsWireForms(t *testing.T) {
	cases := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"float64 epoch seconds", 1758000000.5, "1758000000.5"},
		{"uint64 (CBOR positive integer)", uint64(1758000000), "1758000000"},
		{"uint32", uint32(1758), "1758"},
		{"int", 42, "42"},
		{"int32", int32(-7), "-7"},
		{"int64", int64(-1758000000), "-1758000000"},
		{"string epoch seconds", "1758000000.5", "1758000000.5"},
		// ParseFloat's wider grammar canonicalises through the formatted
		// decimal — hex-float and exponent spellings never round-trip as
		// their input form.
		{"hex-float spelling", "0x1p4", "16"},
		{"exponent spelling", "1e3", "1000"},
	}
	for _, tc := range cases {
		got, err := NormalizeTimestampValue(tc.value)
		if err != nil || got != tc.want {
			t.Fatalf("%s: want (%q, nil), got (%q, %v)", tc.name, tc.want, got, err)
		}
	}

	rfc, err := NormalizeTimestampValue("2016-04-04T19:58:46.480Z")
	if err != nil || rfc != "1459799926.48" {
		t.Fatalf("RFC 3339 notation: want 1459799926.48, got (%q, %v)", rfc, err)
	}
}

// TestNormalizeTimestampValueRejectsUnusableValues pins the violation set:
// non-finite numbers in float and string form (ParseFloat accepts "NaN" and
// "Inf" spellings), RFC 3339 instants beyond the nanosecond-representable
// range the conversion reads through UnixNano, nil, and non-timestamp
// types.
func TestNormalizeTimestampValueRejectsUnusableValues(t *testing.T) {
	invalid := []struct {
		name  string
		value interface{}
	}{
		{"float NaN", math.NaN()},
		{"float +Inf", math.Inf(1)},
		{"float -Inf", math.Inf(-1)},
		{"string NaN", "NaN"},
		{"string nan", "nan"},
		{"string Inf", "Inf"},
		{"string -Infinity", "-Infinity"},
		{"RFC 3339 past the nanosecond range", "9999-12-31T00:00:00Z"},
		{"RFC 3339 before the nanosecond range", "1600-01-01T00:00:00Z"},
		{"nil", nil},
		{"boolean", true},
	}
	for _, tc := range invalid {
		if _, err := NormalizeTimestampValue(tc.value); !errors.Is(err, ErrNonTimestampParameter) {
			t.Fatalf("%s: want ErrNonTimestampParameter, got %v", tc.name, err)
		}
	}
}

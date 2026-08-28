package rds

import "testing"

// TestResolveDescribeMaxRecords pins the documented MaxRecords window:
// "Default: 100. Constraints: Minimum 20, maximum 100." — unset resolves to
// the default and out-of-range values clamp to the nearest bound.
func TestResolveDescribeMaxRecords(t *testing.T) {
	cases := []struct {
		name       string
		maxRecords int
		want       int
	}{
		{"unset resolves to default", 0, 100},
		{"negative clamps up to minimum", -5, 20},
		{"below minimum clamps up", 1, 20},
		{"below minimum clamps up from nineteen", 19, 20},
		{"minimum is kept", 20, 20},
		{"in-range value is kept", 50, 50},
		{"maximum is kept", 100, 100},
		{"above maximum clamps down", 101, 100},
		{"far above maximum clamps down", 500, 100},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ResolveDescribeMaxRecords(tc.maxRecords); got != tc.want {
				t.Fatalf("ResolveDescribeMaxRecords(%d) = %d, want %d", tc.maxRecords, got, tc.want)
			}
		})
	}
}

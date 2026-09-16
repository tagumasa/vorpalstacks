package eventbridge

import (
	"strings"
	"testing"
)

// TestPatternValidationRejections is the acceptance-side class-closure test:
// every documented rejection form for an event pattern, each of which must
// surface as InvalidEventPatternException instead of being persisted to
// fail silently at every match.
func TestPatternValidationRejections(t *testing.T) {
	rejected := []struct {
		name    string
		pattern string
	}{
		{"not JSON", `{"source": ["aws.ec2"`},
		{"top-level array", `["aws.ec2"]`},
		{"top-level scalar", `"aws.ec2"`},
		{"unknown operator in a list", `{"detail": {"state": [{"regexp": "x"}]}]`},
		{"prefix with a list operand", `{"detail": {"state": [{"prefix": ["init", "stop"]}]}}`},
		{"prefix conjunction with an unknown inner operator", `{"detail": {"state": [{"prefix": {"equals": "x"}}]}}`},
		{"equals-ignore-case with a list operand", `{"detail-type": [{"equals-ignore-case": ["a", "b"]}]}`},
		{"wildcard with a non-string operand", `{"detail": {"FileName": [{"wildcard": 5}]}}`},
		{"wildcard consecutive stars", `{"detail": {"FileName": [{"wildcard": "a**b"}]}}`},
		{"wildcard escape of another character", `{"detail": {"FileName": [{"wildcard": "a\\xb"}]}}`},
		{"wildcard trailing lone backslash", `{"detail": {"FileName": [{"wildcard": "ab\\"}]}}`},
		{"numeric without pairs", `{"detail": {"c-count": [{"numeric": [">"]}]}}`},
		{"numeric with an odd operand count", `{"detail": {"c-count": [{"numeric": [">", 1, "<"]}]}}`},
		{"numeric with an unknown comparator", `{"detail": {"c-count": [{"numeric": ["!=", 1]}]}}`},
		{"numeric with a non-number value", `{"detail": {"c-count": [{"numeric": [">", "5"]}]}}`},
		{"exists with a non-boolean", `{"detail": {"state": [{"exists": "yes"}]}}`},
		{"cidr with a non-string", `{"detail": {"ip": [{"cidr": 10}]}}`},
		{"anything-but mixed list", `{"detail": {"x": [{"anything-but": ["a", 1]}]}}`},
		{"anything-but empty list", `{"detail": {"x": [{"anything-but": []}]}}`},
		{"anything-but boolean scalar", `{"detail": {"x": [{"anything-but": true}]}}`},
		{"anything-but null scalar", `{"detail": {"x": [{"anything-but": null}]}}`},
		{"anything-but with an unknown conjunction", `{"detail": {"x": [{"anything-but": {"regexp": "a"}}]}}`},
		{"anything-but wildcard with consecutive stars", `{"detail": {"x": [{"anything-but": {"wildcard": ["a**b"]}}]}}`},
		{"$or not an array", `{"$or": {"source": ["aws.ec2"]}}`},
		{"$or with a scalar entry", `{"$or": ["aws.ec2"]}`},
		{"$or empty", `{"$or": []}`},
	}
	for _, tc := range rejected {
		t.Run(tc.name, func(t *testing.T) {
			err := validateEventPatternStructure(tc.pattern)
			if err == nil {
				t.Fatalf("pattern must be rejected: %s", tc.pattern)
			}
			if !strings.Contains(err.Error(), "InvalidEventPatternException") {
				t.Fatalf("rejection must be an InvalidEventPatternException, got: %v", err)
			}
		})
	}

	// $or combination products multiply: one $or with 40 entries and one
	// with 26 give 1040 combinations, over the documented ceiling of 1000,
	// while each array alone stays under it.
	entries := func(n int) string {
		parts := make([]string, n)
		for i := range parts {
			parts[i] = `{"detail": {"c-count": [{"numeric": ["=", ` + string(rune('0'+i%10)) + `]}]}}`
		}
		return strings.Join(parts, ",")
	}
	overCeiling := `{"$or": [` + entries(40) + `], "detail": {"$or": [` + entries(26) + `]}}`
	if err := validateEventPatternStructure(overCeiling); err == nil {
		t.Fatal("a $or combination product above 1000 must be rejected")
	} else if !strings.Contains(err.Error(), "1040") {
		t.Fatalf("the rejection must report the combination count, got: %v", err)
	}

	// Within the ceiling: a single 1000-entry $or is acceptable, 1001 is not.
	atCeiling := `{"$or": [` + entries(1000) + `]}`
	if err := validateEventPatternStructure(atCeiling); err != nil {
		t.Fatalf("exactly 1000 combinations must be accepted: %v", err)
	}
	overCeilingSingle := `{"$or": [` + entries(1001) + `]}`
	if err := validateEventPatternStructure(overCeilingSingle); err == nil {
		t.Fatal("1001 combinations must be rejected")
	}

	// The empty pattern is the "unset" value and passes.
	if err := validateEventPatternStructure(""); err != nil {
		t.Fatalf("an empty pattern must pass: %v", err)
	}

	// The documentation's complex example validates as a whole.
	valid := `{
		"time": [{"prefix": "2017-10-02"}],
		"detail": {
			"state": [{"anything-but": "initializing"}],
			"c-count": [{"numeric": [">", 0, "<=", 5]}],
			"d-count": [{"numeric": ["<", 10]}],
			"x-limit": [{"anything-but": [100, 200, 300]}]
		}
	}`
	if err := validateEventPatternStructure(valid); err != nil {
		t.Fatalf("the documentation's complex example must validate: %v", err)
	}
}

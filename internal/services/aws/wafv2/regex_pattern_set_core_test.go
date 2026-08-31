package wafv2

import (
	"strings"
	"testing"
)

// TestParseRegularExpressionListRequiredMember pins the
// RegularExpressionList contract of CreateRegexPatternSet and
// UpdateRegexPatternSet: the member is required, so an omitted (nil)
// value, a non-list value and non-object entries are rejected, while the
// empty array stays valid. An object without a RegexString member is
// accepted and contributes no pattern (the member is not required on the
// Regex shape).
func TestParseRegularExpressionListRequiredMember(t *testing.T) {
	if _, err := parseRegularExpressionList(nil); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("nil RegularExpressionList must be rejected with WAFInvalidParameterException, got %v", err)
	}
	if _, err := parseRegularExpressionList(".*"); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("non-list RegularExpressionList must be rejected, got %v", err)
	}
	if _, err := parseRegularExpressionList([]interface{}{"not-an-object"}); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("non-object RegularExpressionList entry must be rejected, got %v", err)
	}
	patterns, err := parseRegularExpressionList([]interface{}{})
	if err != nil {
		t.Fatalf("empty RegularExpressionList array must be valid: %v", err)
	}
	if len(patterns) != 0 {
		t.Fatalf("empty RegularExpressionList must parse to zero patterns, got %d", len(patterns))
	}
	patterns, err = parseRegularExpressionList([]interface{}{map[string]interface{}{}})
	if err != nil {
		t.Fatalf("entry without RegexString must be accepted: %v", err)
	}
	if len(patterns) != 0 {
		t.Fatalf("entry without RegexString must contribute no pattern, got %d", len(patterns))
	}
}

// TestParseRegularExpressionListEntries pins entry-level parsing: a
// compilable RegexString is kept, a broken one is rejected.
func TestParseRegularExpressionListEntries(t *testing.T) {
	patterns, err := parseRegularExpressionList([]interface{}{
		map[string]interface{}{"RegexString": "[0-9]+"},
	})
	if err != nil {
		t.Fatalf("valid pattern must parse: %v", err)
	}
	if len(patterns) != 1 || patterns[0] != "[0-9]+" {
		t.Fatalf("expected [0-9]+ to be kept, got %v", patterns)
	}
	if _, err := parseRegularExpressionList([]interface{}{
		map[string]interface{}{"RegexString": "[unclosed"},
	}); err == nil || !strings.Contains(err.Error(), "WAFInvalidParameterException") {
		t.Fatalf("broken regex must be rejected, got %v", err)
	}
}

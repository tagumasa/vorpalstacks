package athena

import (
	"strings"
	"testing"
)

// TestValidateTagsReportsOffenderAccurately pins that tag violation
// messages quote the offender found by the shared checker: the reported
// length is the Unicode character count of the actual offending key or
// value, never an in-range count of a rune-legal multibyte key, and the
// choice does not depend on map iteration order.
func TestValidateTagsReportsOffenderAccurately(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	// The 129-character ASCII key is the offender; the 100-character CJK
	// key (300 bytes) is rune-legal and must never be blamed for it.
	tags := map[string]string{
		strings.Repeat("a", 129): "v",
		strings.Repeat(cjk, 100): "v",
	}
	want := "TagKey length must be between 1 and 128 (got 129)"
	for i := 0; i < 100; i++ {
		err := validateTags(tags)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("validateTags = %v, want message containing %q", err, want)
		}
	}

	// Value lengths are counted in Unicode characters, not bytes.
	err := validateTags(map[string]string{"env": strings.Repeat(cjk, 257)})
	want = "TagValue length must be between 0 and 256 (got 257)"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("validateTags 257-character CJK value = %v, want message containing %q", err, want)
	}
	if err := validateTags(map[string]string{"env": strings.Repeat(cjk, 256)}); err != nil {
		t.Errorf("validateTags 256-character CJK value = %v, want nil", err)
	}
}

// TestValidateQueryStringSizeUnicodeLengths pins that the QueryString
// @length(1, 262144) trait is counted in Unicode characters: the shape
// carries no pattern, so a query whose text is rune-legal must not be
// rejected on byte length.
func TestValidateQueryStringSizeUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateQueryStringSize(strings.Repeat(cjk, 262144)); err != nil {
		t.Errorf("262144-character CJK query rejected: %v", err)
	}
	if err := validateQueryStringSize(strings.Repeat(cjk, 262145)); err == nil {
		t.Error("262145-character CJK query accepted")
	}
	if err := validateQueryStringSize(""); err == nil {
		t.Error("empty query accepted")
	}
}

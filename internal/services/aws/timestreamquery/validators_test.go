package timestreamquery

import (
	"strings"
	"testing"
)

// TestValidateQueryStringUnicodeLengths pins that QueryString @length(1,
// 262144) is counted in Unicode characters; the shape carries no pattern,
// so multibyte SQL literals are valid input and must not be rejected on
// byte length.
func TestValidateQueryStringUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateQueryString(strings.Repeat(cjk, 262144)); err != nil {
		t.Errorf("262144-character CJK query rejected: %v", err)
	}
	if err := validateQueryString(strings.Repeat(cjk, 262145)); err == nil {
		t.Error("262145-character CJK query accepted")
	}
	if err := validateQueryString(""); err == nil {
		t.Error("empty query accepted")
	}
}

// TestValidateTagKeyUnicodeLengths pins that TagKey @length(1,128) is
// counted in Unicode characters (the Timestream Query TagKey shape carries
// no pattern, unlike the general AWS tag guidance).
func TestValidateTagKeyUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateTagKey(strings.Repeat(cjk, 128)); err != nil {
		t.Errorf("128-character CJK tag key rejected: %v", err)
	}
	if err := validateTagKey(strings.Repeat(cjk, 129)); err == nil {
		t.Error("129-character CJK tag key accepted")
	}
	if err := validateTagValue(strings.Repeat(cjk, 256)); err != nil {
		t.Errorf("256-character CJK tag value rejected: %v", err)
	}
	if err := validateTagValue(strings.Repeat(cjk, 257)); err == nil {
		t.Error("257-character CJK tag value accepted")
	}
}

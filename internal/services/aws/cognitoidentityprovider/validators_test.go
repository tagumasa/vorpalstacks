package cognitoidentityprovider

import (
	"strings"
	"testing"
)

// TestValidateCustomAttributeNameUnicodeLengths pins that
// CustomAttributeNameType follows the Smithy @length(1, 20) trait counted in
// Unicode characters; the shape's pattern uses Unicode categories, so
// multibyte names are valid input and must not be rejected on byte length.
func TestValidateCustomAttributeNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateCustomAttributeName(strings.Repeat(cjk, 20)); err != nil {
		t.Errorf("20-character CJK custom attribute name rejected: %v", err)
	}
	if err := validateCustomAttributeName(strings.Repeat(cjk, 21)); err == nil {
		t.Error("21-character CJK custom attribute name accepted")
	}
}

// TestValidateUsernamePatternUnicodeLengths pins that UsernameType follows
// the Smithy @length(1, 128) trait counted in Unicode characters; the
// pattern uses Unicode categories (\p{L} and friends).
func TestValidateUsernamePatternUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if !validateUsernamePattern(strings.Repeat(cjk, 128)) {
		t.Error("128-character CJK username rejected")
	}
	if validateUsernamePattern(strings.Repeat(cjk, 129)) {
		t.Error("129-character CJK username accepted")
	}
	if validateUsernamePattern("") {
		t.Error("empty username accepted")
	}
}

// TestValidateRegionNameUnicodeLengths pins that RegionNameType follows the
// Smithy @length(5, 32) trait counted in Unicode characters (no pattern).
func TestValidateRegionNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateRegionName(strings.Repeat(cjk, 32)); err != nil {
		t.Errorf("32-character CJK region name rejected: %v", err)
	}
	if err := validateRegionName(strings.Repeat(cjk, 33)); err == nil {
		t.Error("33-character CJK region name accepted")
	}
}

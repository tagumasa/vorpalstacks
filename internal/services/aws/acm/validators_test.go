package acm

import (
	"strings"
	"testing"
)

// TestIsValidFilterStringUnicodeLengths pins that the Smithy FilterString
// @length(1, 256) trait is counted in Unicode characters; the shape carries
// no pattern, so rune-legal multibyte filter values must not be rejected on
// byte length.
func TestIsValidFilterStringUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if !isValidFilterString(strings.Repeat(cjk, 256)) {
		t.Error("256-character CJK filter string rejected")
	}
	if isValidFilterString(strings.Repeat(cjk, 257)) {
		t.Error("257-character CJK filter string accepted")
	}
	if isValidFilterString("") {
		t.Error("empty filter string accepted")
	}
}

// TestValidateNextTokenUnicodeLengths pins that the Smithy NextToken
// @length(1, 10000) trait is counted in Unicode characters; the shape's
// pattern permits Latin-1 characters, which encode as two bytes.
func TestValidateNextTokenUnicodeLengths(t *testing.T) {
	latin1 := "\u00e9" // é, two bytes in UTF-8

	if err := validateNextToken(strings.Repeat(latin1, 10000)); err != nil {
		t.Errorf("10000-character Latin-1 next token rejected: %v", err)
	}
	if err := validateNextToken(strings.Repeat(latin1, 10001)); err == nil {
		t.Error("10001-character Latin-1 next token accepted")
	}
	if err := validateNextToken(""); err != nil {
		t.Errorf("empty next token rejected: %v", err)
	}
}

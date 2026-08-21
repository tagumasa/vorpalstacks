package appsync

import (
	"strings"
	"testing"
)

// TestValidateDescriptionUnicodeLengths pins that AppSync descriptions
// follow the Smithy Description @length(0, 255) trait counted in Unicode
// characters; the shape's "^.*$" pattern admits multibyte text.
func TestValidateDescriptionUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateDescription(strings.Repeat(cjk, 255)); err != nil {
		t.Errorf("255-character CJK description rejected: %v", err)
	}
	if err := validateDescription(strings.Repeat(cjk, 256)); err == nil {
		t.Error("256-character CJK description accepted")
	}
	if err := validateDescription(""); err != nil {
		t.Errorf("empty description rejected: %v", err)
	}
}

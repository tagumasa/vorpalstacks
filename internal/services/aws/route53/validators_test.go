package route53

import (
	"strings"
	"testing"
)

// TestValidateCommentUnicodeLengths pins that hosted-zone comments follow
// the Smithy ResourceDescription @length(0, 256) trait counted in Unicode
// characters; the shape carries no pattern, so multibyte comments are valid
// input and must not be rejected on byte length.
func TestValidateCommentUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateComment(strings.Repeat(cjk, 128)); err != nil {
		t.Errorf("128-character CJK comment rejected: %v", err)
	}
	if err := validateComment(strings.Repeat(cjk, 256)); err != nil {
		t.Errorf("256-character CJK comment rejected: %v", err)
	}
	if err := validateComment(strings.Repeat(cjk, 257)); err == nil {
		t.Error("257-character CJK comment accepted")
	}
	if err := validateComment(""); err != nil {
		t.Errorf("empty comment rejected: %v", err)
	}
}

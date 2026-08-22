package dynamodb

import (
	"strings"
	"testing"
)

// TestValidateLengthUnicodeCharacters pins that the shared length helper
// counts Unicode characters, not bytes: every @length-bearing string shape
// it guards (PartiQL statements, attribute names, S3 prefixes, tag keys and
// values, …) is validated per the Smithy string semantics where multibyte
// input is legal unless an ASCII pattern coexists.
func TestValidateLengthUnicodeCharacters(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	// 100 CJK runes are 300 bytes: rune-legal for max 255 despite exceeding
	// 255 bytes.
	if !validateLength(strings.Repeat(cjk, 100), 1, 255) {
		t.Error("100-character CJK string rejected for max 255")
	}
	if !validateLength(strings.Repeat(cjk, 255), 1, 255) {
		t.Error("255-character CJK string rejected for max 255")
	}
	if validateLength(strings.Repeat(cjk, 256), 1, 255) {
		t.Error("256-character CJK string accepted for max 255")
	}
	if validateLength("", 1, 255) {
		t.Error("empty string accepted for min 1")
	}
	if !validateLength(strings.Repeat(cjk, 1), 1, 36) {
		t.Error("single CJK character rejected for max 36")
	}
}

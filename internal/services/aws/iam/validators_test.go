package iam

import (
	"strings"
	"testing"

	tagutil "vorpalstacks/internal/common/tags"
)

// TestValidateTagEntriesUnicodeLengths pins that tag key and value limits
// are counted in Unicode characters, matching the Smithy tagKeyType and
// tagValueType @length traits.
func TestValidateTagEntriesUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateTagEntries([]tagutil.Tag{{Key: strings.Repeat(cjk, 100), Value: strings.Repeat(cjk, 200)}}); err != nil {
		t.Errorf("100-character CJK key with 200-character value rejected: %v", err)
	}
	if err := validateTagEntries([]tagutil.Tag{{Key: strings.Repeat(cjk, 129), Value: "v"}}); err == nil {
		t.Error("129-character CJK key accepted, want rejection")
	}
	if err := validateTagEntries([]tagutil.Tag{{Key: "k", Value: strings.Repeat(cjk, 257)}}); err == nil {
		t.Error("257-character CJK value accepted, want rejection")
	}
}

// TestValidateRoleDescriptionLatin1 pins that roleDescriptionType lengths
// are counted in Unicode characters: the shape's pattern admits Latin-1
// supplement characters (2 bytes each in UTF-8), so 600 such characters
// (1200 bytes) fall within the 1000-character bound.
func TestValidateRoleDescriptionLatin1(t *testing.T) {
	latin1 := "é" // U+00E9, inside the pattern's \u00A1-\u00FF range

	if !validateRoleDescription(strings.Repeat(latin1, 600)) {
		t.Error("600-character Latin-1 role description rejected")
	}
	if validateRoleDescription(strings.Repeat(latin1, 1001)) {
		t.Error("1001-character Latin-1 role description accepted")
	}
	if !validateRoleDescription(strings.Repeat("a", 1000)) {
		t.Error("1000-character ASCII role description rejected")
	}
	if validateRoleDescription("✓") {
		t.Error("character outside the Latin-1 pattern accepted")
	}
}

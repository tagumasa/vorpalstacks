package secretsmanager

import (
	"strings"
	"testing"
)

// TestValidateSecretNameUnicodeLengths pins that NameType @length(1,512) is
// counted in Unicode characters; the shape carries no pattern.
func TestValidateSecretNameUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateSecretName(strings.Repeat(cjk, 512)); err != nil {
		t.Errorf("512-character CJK secret name rejected: %v", err)
	}
	if err := validateSecretName(strings.Repeat(cjk, 513)); err == nil {
		t.Error("513-character CJK secret name accepted")
	}
}

// TestValidateDescriptionUnicodeLengths pins that DescriptionType
// @length(0,2048) is counted in Unicode characters; the shape carries no
// pattern, so multibyte descriptions are valid input.
func TestValidateDescriptionUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateDescription(strings.Repeat(cjk, 2048)); err != nil {
		t.Errorf("2048-character CJK description rejected: %v", err)
	}
	if err := validateDescription(strings.Repeat(cjk, 2049)); err == nil {
		t.Error("2049-character CJK description accepted")
	}
	if err := validateDescription(""); err != nil {
		t.Errorf("empty description rejected: %v", err)
	}
}

// TestValidateExcludeCharactersUnicodeLengths pins that
// ExcludeCharactersType @length(0,4096) is counted in Unicode characters;
// multibyte exclusion sets (for example CJK ranges) are valid input.
func TestValidateExcludeCharactersUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateExcludeCharacters(strings.Repeat(cjk, 4096)); err != nil {
		t.Errorf("4096-character CJK exclude set rejected: %v", err)
	}
	if err := validateExcludeCharacters(strings.Repeat(cjk, 4097)); err == nil {
		t.Error("4097-character CJK exclude set accepted")
	}
}

// TestValidateSecretStringLengthByteCeiling pins the byte-based ceiling on
// SecretString: the AWS quota page caps the secret value at 65,536 bytes,
// and because every Unicode character occupies at least one byte, a byte
// ceiling of 65536 is equivalent to enforcing both the Smithy character
// constraint and the storage quota. A rune-legal multibyte value whose
// encoded size exceeds the quota must therefore stay rejected.
func TestValidateSecretStringLengthByteCeiling(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateSecretStringLength(strings.Repeat("a", 65536)); err != nil {
		t.Errorf("65536-byte ASCII secret rejected: %v", err)
	}
	if err := validateSecretStringLength(strings.Repeat("a", 65537)); err == nil {
		t.Error("65537-byte ASCII secret accepted")
	}
	if err := validateSecretStringLength(strings.Repeat(cjk, 30000)); err == nil {
		t.Error("30000-character CJK secret (90000 bytes) accepted despite the byte quota")
	}
}

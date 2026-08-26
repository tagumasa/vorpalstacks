package secretsmanager

import (
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
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

// TestSecretsManagerErrorCodesAreModelled pins that validators only emit
// error codes defined by the secrets-manager-2017-10-17 model. The model's
// error inventory has no ValidationException shape (unlike Scheduler), so
// both the batch-size limiter and the SecretBinary decoder must report
// InvalidParameterException; an alien code can never match a typed error
// in the AWS SDKs.
func TestSecretsManagerErrorCodesAreModelled(t *testing.T) {
	err := validateSecretIdList(make([]string, maxSecretIdListItems+1))
	if err == nil {
		t.Fatal("oversized SecretIdList accepted")
	}
	if code := errorCode(err); code != "InvalidParameterException" {
		t.Errorf("SecretIdList overflow error code = %q, want InvalidParameterException", code)
	}

	_, err = decodeAndValidateSecretBinary("not-base64-$$$")
	if err == nil {
		t.Fatal("invalid base64 SecretBinary accepted")
	}
	if code := errorCode(err); code != "InvalidParameterException" {
		t.Errorf("SecretBinary decode error code = %q, want InvalidParameterException", code)
	}
}

// errorCode extracts the wire code from an *awserrors.AWSError.
func errorCode(err error) string {
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		return ""
	}
	return awsErr.GetCode()
}

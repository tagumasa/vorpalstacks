package sts

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// TestValidateARNUnicodeLengths pins that the STS arnType @length(20,2048)
// is counted in Unicode characters; the arnType pattern admits code points
// well beyond ASCII (U+00A0-U+D7FF and beyond), so rune-legal multibyte
// ARNs must not be rejected on byte length.
func TestValidateARNUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes, within U+00A0-U+D7FF

	prefix := "arn:aws:iam::123456789012:role/"
	// prefix plus CJK padding: 33 + 700 runes = 733 runes (2199 bytes).
	if err := validateARN(prefix + strings.Repeat(cjk, 700)); err != nil {
		t.Errorf("733-character ARN with CJK tail rejected: %v", err)
	}
	if err := validateARN(prefix + strings.Repeat(cjk, 2032)); err == nil {
		t.Error("2065-character ARN with CJK tail accepted")
	}
	if err := validateRoleArn(prefix + strings.Repeat(cjk, 700)); err != nil {
		t.Errorf("733-character role ARN with CJK tail rejected: %v", err)
	}
	if err := validateRoleArn(strings.Repeat("a", 19)); err == nil {
		t.Error("19-character role ARN accepted")
	}
}

// TestValidateSessionPolicyUnicodeLengths pins that
// sessionPolicyDocumentType @length(1,2048) is counted in Unicode
// characters; the shape's pattern is Latin-1, permitting two-byte
// characters.
func TestValidateSessionPolicyUnicodeLengths(t *testing.T) {
	latin1 := "\u00e9" // é, two bytes in UTF-8, within \x20-\xff

	policy := `{"Version":"2012-10-17","Statement":[]}`
	if err := validateSessionPolicy(policy); err != nil {
		t.Errorf("ASCII session policy rejected: %v", err)
	}
	// A rune-legal Latin-1 policy of 2048 characters is 4096 bytes on the
	// wire but within the Smithy character bound.
	long := `{"é":"` + strings.Repeat(latin1, 2025) + `"}`
	if n := len([]rune(long)); n > maxSessionPolicyLen {
		t.Fatalf("test bug: policy is %d runes", n)
	}
	if err := validateSessionPolicy(long); err != nil {
		t.Errorf("2048-character Latin-1 session policy rejected: %v", err)
	}
}

// TestFitWebIdentitySubjectRuneSafeTruncation pins that an over-long
// webIdentitySubjectType value is truncated to 255 Unicode characters on
// rune boundaries — never mid-sequence — and short fallbacks are padded to
// the minimum of 6 characters.
func TestFitWebIdentitySubjectRuneSafeTruncation(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	long := strings.Repeat(cjk, 260)
	got := fitWebIdentitySubject(long)
	if n := utf8.RuneCountInString(got); n != 255 {
		t.Fatalf("truncated subject length = %d runes, want 255", n)
	}
	if !utf8.ValidString(got) {
		t.Error("truncated subject is not valid UTF-8")
	}
	if got != strings.Repeat(cjk, 255) {
		t.Error("truncated subject content mismatch")
	}

	if got := fitWebIdentitySubject(strings.Repeat(cjk, 255)); len(got) != len(strings.Repeat(cjk, 255)) {
		t.Error("255-character subject should pass through unchanged")
	}
	if got := fitWebIdentitySubject("ab"); got != "____ab" {
		t.Errorf("short subject padding = %q, want ____ab", got)
	}
	if got := fitWebIdentitySubject(strings.Repeat(cjk, 2)); utf8.RuneCountInString(got) != 6 {
		t.Errorf("2-character CJK subject (6 bytes) padded to %d runes, want 6", utf8.RuneCountInString(got))
	}
}

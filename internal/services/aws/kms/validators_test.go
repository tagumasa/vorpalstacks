package kms

import (
	"encoding/base64"
	"strings"
	"testing"
)

// TestValidateDescriptionLengthUnicodeLengths pins that DescriptionType
// @length(0,8192) is counted in Unicode characters; the shape carries no
// pattern, so rune-legal multibyte descriptions must not be rejected on
// byte length.
func TestValidateDescriptionLengthUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateDescriptionLength(strings.Repeat(cjk, 8192)); err != nil {
		t.Errorf("8192-character CJK description rejected: %v", err)
	}
	if err := validateDescriptionLength(strings.Repeat(cjk, 8193)); err == nil {
		t.Error("8193-character CJK description accepted")
	}
	if err := validateDescriptionLength(""); err != nil {
		t.Errorf("empty description rejected: %v", err)
	}
}

// TestValidateMarkerLengthUnicodeLengths pins that MarkerType @length(1,1024)
// is counted in Unicode characters; the shape's pattern is Latin-1, whose
// characters encode as two UTF-8 bytes.
func TestValidateMarkerLengthUnicodeLengths(t *testing.T) {
	latin1 := "\u00e9" // é, two bytes in UTF-8

	if err := validateMarkerLength(strings.Repeat(latin1, 1024)); err != nil {
		t.Errorf("1024-character Latin-1 marker rejected: %v", err)
	}
	if err := validateMarkerLength(strings.Repeat(latin1, 1025)); err == nil {
		t.Error("1025-character Latin-1 marker accepted")
	}
	if err := validateMarkerLength(""); err != nil {
		t.Errorf("empty marker rejected: %v", err)
	}
}

// TestValidateImportBlobDecodedSize pins that the ImportKeyMaterial blob
// guards measure the DECODED byte size against the Smithy CiphertextType
// @length(1,6144) bound: a legitimate 6144-byte blob arrives as 8192 base64
// characters and must be accepted, while a 6145-byte blob and a non-base64
// value are rejected.
func TestValidateImportBlobDecodedSize(t *testing.T) {
	exact := base64.StdEncoding.EncodeToString(make([]byte, 6144))
	over := base64.StdEncoding.EncodeToString(make([]byte, 6145))

	if _, err := decodeEncryptedKeyMaterial(exact); err != nil {
		t.Errorf("6144-byte blob (8192 base64 chars) rejected: %v", err)
	}
	if err := validateImportTokenSize(exact); err != nil {
		t.Errorf("6144-byte token blob rejected: %v", err)
	}
	if _, err := decodeEncryptedKeyMaterial(over); err == nil {
		t.Error("6145-byte blob accepted")
	}
	if err := validateImportTokenSize("not-base64!!"); err == nil {
		t.Error("non-base64 token accepted")
	}
	// Encoded input longer than the derived bound (8192 characters for a
	// 6144-byte blob) must fail the size check before base64 validity is
	// even assessed, so oversized payloads never reach the decoder.
	long := strings.Repeat("A", 8193)
	if err := validateImportTokenSize(long); err == nil || !strings.Contains(err.Error(), "exceeds maximum length") {
		t.Errorf("8193-character encoded token not size-rejected: %v", err)
	}
	// CiphertextType @length(1,6144) also carries a minimum: an empty
	// ImportToken or EncryptedKeyMaterial must be rejected, not deferred
	// to the downstream store comparison.
	if err := validateImportTokenSize(""); err == nil {
		t.Error("empty token accepted")
	}
	if _, err := decodeEncryptedKeyMaterial(""); err == nil {
		t.Error("empty key material accepted")
	}
}

package kms

import (
	"encoding/base64"
	"strings"
	"testing"

	kmsstore "vorpalstacks/internal/store/aws/kms"
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

// TestAlgorithmContractErrors pins the error codes of the algorithm member
// contract: a value outside the modelled enum is a shape violation rejected
// with SerializationException (the aws-json-1.1 protocol contract, which no
// service model enumerates), while an enum value the key does not support is
// rejected with InvalidKeyUsageException ("the encryption algorithm or
// signing algorithm specified for the operation is incompatible with the
// type of key material in the KMS key"). The InvalidAlgorithmException code
// appears nowhere in the KMS model.
func TestAlgorithmContractErrors(t *testing.T) {
	symmetric := &kmsstore.Key{
		KeySpec:              kmsstore.KeySpecSymmetricDefault,
		EncryptionAlgorithms: []string{"SYMMETRIC_DEFAULT"},
	}
	rsaKey := &kmsstore.Key{
		KeySpec:              kmsstore.KeySpecRSA2048,
		EncryptionAlgorithms: []string{"RSAES_OAEP_SHA_1", "RSAES_OAEP_SHA_256"},
	}

	// Enum violations are shape violations.
	if _, err := resolveEncryptionAlgorithm(symmetric, map[string]interface{}{
		"EncryptionAlgorithm": "InvalidAlgorithmValue",
	}, "EncryptionAlgorithm"); err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Errorf("enum-invalid encryption algorithm: got %v, want SerializationException", err)
	}
	if err := resolveSigningAlgorithm("NOT_AN_ALGORITHM", rsaKey); err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Errorf("enum-invalid signing algorithm: got %v, want SerializationException", err)
	}
	macKey := &kmsstore.Key{MacAlgorithms: []string{"HMAC_SHA_256"}}
	if err := resolveMacAlgorithm("HMAC_NOT_REAL", macKey); err == nil || !strings.Contains(err.Error(), "SerializationException") {
		t.Errorf("enum-invalid MAC algorithm: got %v, want SerializationException", err)
	}

	// Enum values unsupported by the key are InvalidKeyUsageException.
	if _, err := resolveEncryptionAlgorithm(symmetric, map[string]interface{}{
		"EncryptionAlgorithm": "RSAES_OAEP_SHA_256",
	}, "EncryptionAlgorithm"); err == nil || !strings.Contains(err.Error(), "InvalidKeyUsageException") {
		t.Errorf("unsupported encryption algorithm: got %v, want InvalidKeyUsageException", err)
	}
	rsaSignKey := &kmsstore.Key{SigningAlgorithms: []string{"RSASSA_PKCS1_V1_5_SHA_256"}}
	if err := resolveSigningAlgorithm("ECDSA_SHA_256", rsaSignKey); err == nil || !strings.Contains(err.Error(), "InvalidKeyUsageException") {
		t.Errorf("unsupported signing algorithm: got %v, want InvalidKeyUsageException", err)
	}
	if err := resolveMacAlgorithm("HMAC_SHA_512", macKey); err == nil || !strings.Contains(err.Error(), "InvalidKeyUsageException") {
		t.Errorf("unsupported MAC algorithm: got %v, want InvalidKeyUsageException", err)
	}

	// Explicit supported values and omitted members resolve.
	if alg, err := resolveEncryptionAlgorithm(rsaKey, map[string]interface{}{
		"EncryptionAlgorithm": "RSAES_OAEP_SHA_1",
	}, "EncryptionAlgorithm"); err != nil || alg != "RSAES_OAEP_SHA_1" {
		t.Errorf("explicit supported encryption algorithm: %q %v", alg, err)
	}
	if alg, err := resolveEncryptionAlgorithm(rsaKey, map[string]interface{}{}, "EncryptionAlgorithm"); err != nil || alg != "RSAES_OAEP_SHA_256" {
		t.Errorf("omitted encryption algorithm default: %q %v", alg, err)
	}

	// The ReEncrypt members resolve through their own member names.
	if alg, err := resolveEncryptionAlgorithm(rsaKey, map[string]interface{}{
		"SourceEncryptionAlgorithm": "RSAES_OAEP_SHA_1",
	}, "SourceEncryptionAlgorithm"); err != nil || alg != "RSAES_OAEP_SHA_1" {
		t.Errorf("SourceEncryptionAlgorithm: %q %v", alg, err)
	}
	// An operation reading the Source member must not pick up a value sent
	// under the single-operation member name: the lookup is by member name,
	// so the keyspec default applies.
	if alg, err := resolveEncryptionAlgorithm(rsaKey, map[string]interface{}{
		"EncryptionAlgorithm": "RSAES_OAEP_SHA_1",
	}, "SourceEncryptionAlgorithm"); err != nil || alg != "RSAES_OAEP_SHA_256" {
		t.Errorf("SourceEncryptionAlgorithm fell back to the EncryptionAlgorithm member: %q %v", alg, err)
	}
}

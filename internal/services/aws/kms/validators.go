package kms

import (
	"fmt"
	"regexp"

	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// maxImportBlobSize is the AWS-documented maximum length (base64-encoded)
// for EncryptedKeyMaterial and ImportToken parameters.
const maxImportBlobSize = 6144

// validateOriginParams rejects CustomKeyStoreId and XksKeyId values when
// the platform does not support Custom Key Stores (CloudHSM / External
// Key Store). AWS_CLOUDHSM origin is already rejected upstream by an
// explicit ErrUnsupportedOperation guard; any CustomKeyStoreId or
// XksKeyId reaching this validator is therefore a caller error because
// no supported origin accepts it. Without this check the parameters are
// silently dropped, causing the caller to believe a key was created in
// a custom key store when it was not.
func validateOriginParams(customKeyStoreID string, xksKeyID string) error {
	if customKeyStoreID != "" {
		return NewValidationError("CustomKeyStoreId is not supported; the platform does not implement Custom Key Stores")
	}
	if xksKeyID != "" {
		return NewValidationError("XksKeyId is not supported; the platform does not implement External Key Stores")
	}
	return nil
}

// validatePrimaryRegion rejects an empty PrimaryRegion value for
// UpdatePrimaryRegion. AWS returns ValidationException when the
// parameter is omitted; without this check the store would persist a
// broken PrimaryKeyInfo with an empty Region field.
func validatePrimaryRegion(region string) error {
	if region == "" {
		return NewValidationError("PrimaryRegion is required")
	}
	return nil
}

// validateEncryptedKeyMaterialSize enforces the AWS-documented length
// constraint (max 6144 base64-encoded characters) on the
// EncryptedKeyMaterial parameter of ImportKeyMaterial. Without this
// guard an attacker could send an arbitrarily large blob to exhaust
// server memory during base64 decode.
func validateEncryptedKeyMaterialSize(b64 string) error {
	if len(b64) > maxImportBlobSize {
		return NewValidationError(fmt.Sprintf("EncryptedKeyMaterial exceeds maximum length of %d", maxImportBlobSize))
	}
	return nil
}

// validateImportTokenSize enforces the same AWS-documented length
// constraint on the ImportToken parameter.
func validateImportTokenSize(b64 string) error {
	if len(b64) > maxImportBlobSize {
		return NewValidationError(fmt.Sprintf("ImportToken exceeds maximum length of %d", maxImportBlobSize))
	}
	return nil
}

// validateRotationKeyEligibility enforces the AWS constraint that
// automatic key rotation is available only for SYMMETRIC_DEFAULT keys
// with AWS_KMS origin. Both EnableKeyRotation and DisableKeyRotation
// must apply this check to maintain API contract consistency.
func validateRotationKeyEligibility(spec kmsstore.KeySpec, origin kmsstore.OriginType) error {
	if spec != kmsstore.KeySpecSymmetricDefault || origin != kmsstore.OriginTypeAWSKMS {
		return ErrUnsupportedOperation
	}
	return nil
}

// maxDescriptionLen is the AWS-documented maximum length for the
// Description parameter (Smithy DescriptionType: length 0-8192).
const maxDescriptionLen = 8192

// maxPolicyLen is the AWS-documented maximum length for a key policy
// document (Smithy PolicyType: length 1-131072).
const maxPolicyLen = 131072

// validateKeyUsageSpecCombo enforces the AWS constraint matrix that
// defines which KeySpec values are valid for each KeyUsage. The
// constraint table mirrors keyUsageKeySpecConstraints (defined in
// key_operations.go). Both CreateKey and ReplicateKey must apply this
// check to prevent creating keys with incompatible usage/spec pairs.
func validateKeyUsageSpecCombo(usage kmsstore.KeyUsage, spec kmsstore.KeySpec) error {
	allowed, ok := keyUsageKeySpecConstraints[usage]
	if !ok {
		return NewValidationError(fmt.Sprintf("Invalid KeyUsage: %q", usage))
	}
	for _, valid := range allowed {
		if spec == valid {
			return nil
		}
	}
	return NewValidationError(fmt.Sprintf("Invalid KeySpec %q for KeyUsage %q", spec, usage))
}

// validateDescriptionLength enforces the Smithy DescriptionType length
// constraint (0-8192). Both CreateKey and UpdateKeyDescription must
// apply this check.
func validateDescriptionLength(desc string) error {
	if len(desc) > maxDescriptionLen {
		return NewValidationError(fmt.Sprintf("Description length %d exceeds maximum of %d", len(desc), maxDescriptionLen))
	}
	return nil
}

// validatePolicySize enforces the Smithy PolicyType length constraint
// (1-131072). An empty policy string is accepted here because callers
// fall back to the platform default policy when the parameter is unset;
// the length ceiling protects server memory during JSON parsing.
func validatePolicySize(policy string) error {
	if len(policy) > maxPolicyLen {
		return NewValidationError(fmt.Sprintf("Policy length %d exceeds maximum of %d", len(policy), maxPolicyLen))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Smithy constraint constants and patterns
// ---------------------------------------------------------------------------

const (
	// KeyIdType: length 1-2048.
	maxKeyIdLen = 2048
	// MarkerType: length 1-1024.
	maxMarkerLen = 1024
	// GrantIdType: length 1-128.
	maxGrantIdLen = 128
	// GrantTokenType: length 1-8192.
	maxGrantTokenLen = 8192
	// PrincipalIdType: length 1-256.
	maxPrincipalIdLen = 256
	// RegionType: length 1-32.
	maxRegionLen = 32
	// CiphertextType: blob length 1-6144.
	maxCiphertextLen = 6144
	// GrantConstraintSourceArnType: length 20-512.
	minGrantSourceArnLen = 20
	maxGrantSourceArnLen = 512
	// GrantTokenList: length 0-10.
	maxGrantTokenListSize = 10
)

var (
	// PrincipalIdType pattern: word chars, plus = + , . @ : / -.
	principalIdPattern = regexp.MustCompile(`^[\w+=,.@:/-]+$`)
	// RegionType pattern: e.g. us-east-1, ap-southeast-2, cn-north-1.
	regionPattern = regexp.MustCompile(`^([a-z]+-){2,3}\d+$`)
	// MarkerType pattern: Latin-1 range (U+0020-U+00FF).
	markerPattern = regexp.MustCompile(`^[\x20-\xFF]*$`)
	// GrantConstraintSourceArnType pattern: ARN with 12-digit account.
	grantSourceArnPattern = regexp.MustCompile(`^arn:aws[a-z0-9-]*:[a-z0-9-]+:[a-z0-9-]*:[0-9]{12}:.+$`)
)

// validateKeyIdLength enforces the Smithy KeyIdType length constraint
// (1-2048). Applied centrally in resolveKey so that every operation
// accepting a KeyId parameter is covered.
func validateKeyIdLength(keyID string) error {
	if len(keyID) == 0 || len(keyID) > maxKeyIdLen {
		return NewValidationError(fmt.Sprintf("KeyId length %d does not fit range 1-%d", len(keyID), maxKeyIdLen))
	}
	return nil
}

// validateMarkerLength enforces the Smithy MarkerType length constraint
// (1-1024) and pattern (Latin-1 characters only). Applied to all
// paginated list operations.
func validateMarkerLength(marker string) error {
	if marker == "" {
		return nil
	}
	if len(marker) > maxMarkerLen || !markerPattern.MatchString(marker) {
		return NewValidationError(fmt.Sprintf("Marker exceeds maximum length of %d or contains invalid characters", maxMarkerLen))
	}
	return nil
}

// validateGrantIdLength enforces the Smithy GrantIdType length
// constraint (1-128). Applied to RevokeGrant, RetireGrant, and
// ListGrants.
func validateGrantIdLength(grantID string) error {
	if len(grantID) == 0 || len(grantID) > maxGrantIdLen {
		return NewValidationError(fmt.Sprintf("GrantId length %d does not fit range 1-%d", len(grantID), maxGrantIdLen))
	}
	return nil
}

// validateGrantTokenLength enforces the Smithy GrantTokenType length
// constraint (1-8192). Applied to RetireGrant.
func validateGrantTokenLength(token string) error {
	if len(token) == 0 || len(token) > maxGrantTokenLen {
		return NewValidationError(fmt.Sprintf("GrantToken length %d does not fit range 1-%d", len(token), maxGrantTokenLen))
	}
	return nil
}

// validateGrantTokenListSize enforces the Smithy GrantTokenList length
// constraint (0-10). Applied to CreateGrant.
func validateGrantTokenListSize(tokens []interface{}) error {
	if len(tokens) > maxGrantTokenListSize {
		return NewValidationError(fmt.Sprintf("GrantTokens list size %d exceeds maximum of %d", len(tokens), maxGrantTokenListSize))
	}
	return nil
}

// validatePrincipalId enforces the Smithy PrincipalIdType length
// constraint (1-256) and pattern. Applied to GranteePrincipal and
// RetiringPrincipal in CreateGrant.
func validatePrincipalId(principal string) error {
	if len(principal) == 0 || len(principal) > maxPrincipalIdLen {
		return NewValidationError(fmt.Sprintf("Principal length %d does not fit range 1-%d", len(principal), maxPrincipalIdLen))
	}
	if !principalIdPattern.MatchString(principal) {
		return NewValidationError("Principal contains invalid characters")
	}
	return nil
}

// validateReplicaRegion enforces the Smithy RegionType length
// constraint (1-32) and pattern. Applied to ReplicaRegion in
// ReplicateKey.
func validateReplicaRegion(region string) error {
	if len(region) == 0 || len(region) > maxRegionLen {
		return NewValidationError(fmt.Sprintf("ReplicaRegion length %d does not fit range 1-%d", len(region), maxRegionLen))
	}
	if !regionPattern.MatchString(region) {
		return NewValidationError(fmt.Sprintf("ReplicaRegion %q is not a valid AWS region", region))
	}
	return nil
}

// validateCiphertextLength enforces the Smithy CiphertextType blob
// length constraint (1-6144 bytes). Applied to Decrypt and ReEncrypt
// after base64 decoding.
func validateCiphertextLength(n int) error {
	if n == 0 || n > maxCiphertextLen {
		return NewValidationError(fmt.Sprintf("CiphertextBlob length %d does not fit range 1-%d", n, maxCiphertextLen))
	}
	return nil
}

// validateGrantSourceArn enforces the Smithy
// GrantConstraintSourceArnType length constraint (20-512) and ARN
// pattern. Applied to the SourceArn constraint in CreateGrant.
func validateGrantSourceArn(arn string) error {
	if len(arn) < minGrantSourceArnLen || len(arn) > maxGrantSourceArnLen {
		return NewValidationError(fmt.Sprintf("SourceArn length %d does not fit range %d-%d", len(arn), minGrantSourceArnLen, maxGrantSourceArnLen))
	}
	if !grantSourceArnPattern.MatchString(arn) {
		return NewValidationError("SourceArn is not a valid ARN")
	}
	return nil
}

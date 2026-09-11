package iam

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	// UserIDPrefix is the prefix for IAM user IDs.
	UserIDPrefix = "AIDA"
	// AccessKeyIDPrefix is the prefix for IAM access key IDs.
	AccessKeyIDPrefix = "AKIA"
	// GroupIDPrefix is the prefix for IAM group IDs.
	GroupIDPrefix = "AGPA"
	// RoleIDPrefix is the prefix for IAM role IDs.
	RoleIDPrefix = "AROA"
	// PolicyIDPrefix is the prefix for IAM policy IDs.
	PolicyIDPrefix = "ANPA"
	// InstanceProfileIDPrefix is the prefix for IAM instance profile IDs.
	InstanceProfileIDPrefix = "AIPA"
	// ServerCertificateIDPrefix is the prefix for IAM server certificate IDs.
	ServerCertificateIDPrefix = "ASCA"
	// SSHPublicKeyIDPrefix is the prefix for IAM SSH public key IDs (the
	// IAM identifiers table: APKA = public key).
	SSHPublicKeyIDPrefix = "APKA"
	// ServiceSpecificCredentialIDPrefix is the prefix for IAM
	// service-specific credential IDs (the IAM identifiers table:
	// ACCA = context-specific credential).
	ServiceSpecificCredentialIDPrefix = "ACCA"
)

var base32Encoder = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)

// GenerateUserID generates a unique IAM user ID.
func GenerateUserID() (string, error) {
	return generateID(UserIDPrefix)
}

// GenerateAccessKeyID generates a unique IAM access key ID.
func GenerateAccessKeyID() (string, error) {
	return generateID(AccessKeyIDPrefix)
}

// GenerateGroupID generates a unique IAM group ID.
func GenerateGroupID() (string, error) {
	return generateID(GroupIDPrefix)
}

// GenerateRoleID generates a unique IAM role ID.
func GenerateRoleID() (string, error) {
	return generateID(RoleIDPrefix)
}

// GeneratePolicyID generates a unique IAM policy ID.
func GeneratePolicyID() (string, error) {
	return generateID(PolicyIDPrefix)
}

// GenerateInstanceProfileID generates a unique IAM instance profile ID.
func GenerateInstanceProfileID() (string, error) {
	return generateID(InstanceProfileIDPrefix)
}

// GenerateServerCertificateID generates a unique IAM server certificate ID.
func GenerateServerCertificateID() (string, error) {
	return generateID(ServerCertificateIDPrefix)
}

// GenerateSecretAccessKey generates a secure random secret access key.
func GenerateSecretAccessKey() (string, error) {
	bytes := make([]byte, 30)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	trimmed := strings.ReplaceAll(encoded, "=", "")
	if len(trimmed) < 40 {
		return "", fmt.Errorf("generated secret key too short: %d < 40", len(trimmed))
	}
	return trimmed[:40], nil
}

func generateID(prefix string) (string, error) {
	bytes := make([]byte, 10)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + base32Encoder.EncodeToString(bytes), nil
}

// generatePrivateKeyID generates a privateKeyIdType identifier, used for
// SAMLPrivateKey.KeyId and the SAML provider UUID.  Smithy constraint:
// length [22,64], pattern ^[A-Z0-9]+$.  14 crypto-random bytes encoded as
// base32 → 23 characters, satisfying both constraints with strong
// collision resistance for edge/on-prem use.
func generatePrivateKeyID() (string, error) {
	bytes := make([]byte, 14)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32Encoder.EncodeToString(bytes), nil
}

// GenerateSSHPublicKeyID generates a unique IAM SSH public key ID.
func GenerateSSHPublicKeyID() (string, error) {
	return generateID(SSHPublicKeyIDPrefix)
}

// GenerateSigningCertificateID generates a unique IAM signing certificate
// ID. AWS documents no unique-ID prefix for signing certificates — the IAM
// identifiers table's ASCA row covers server certificates, and the
// documented UploadSigningCertificate response carries an unprefixed
// base32-style certificate ID — so the ID is minted without a prefix.
func GenerateSigningCertificateID() (string, error) {
	return generateID("")
}

// GenerateServiceSpecificCredentialID generates a unique IAM
// service-specific credential ID.
func GenerateServiceSpecificCredentialID() (string, error) {
	return generateID(ServiceSpecificCredentialIDPrefix)
}

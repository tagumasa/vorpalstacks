package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

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
	// SAMLProviderIDPrefix is the prefix for IAM SAML provider IDs.
	SAMLProviderIDPrefix = "ARPA"
	// OpenIDConnectProviderIDPrefix is the prefix for IAM OIDC provider IDs.
	OpenIDConnectProviderIDPrefix = "AROA"
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

// GenerateSAMLProviderID generates a unique IAM SAML provider ID.
func GenerateSAMLProviderID() (string, error) {
	return generateID(SAMLProviderIDPrefix)
}

// GenerateOpenIDConnectProviderID generates a unique IAM OIDC provider ID.
func GenerateOpenIDConnectProviderID() (string, error) {
	return generateID(OpenIDConnectProviderIDPrefix)
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

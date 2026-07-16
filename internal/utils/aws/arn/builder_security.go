// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import "strings"

// STSBuilder provides methods for constructing STS (Security Token Service) ARNs.
type STSBuilder struct{ *ARNBuilder }

// STS returns an STSBuilder for constructing STS ARNs.
func (b *ARNBuilder) STS() *STSBuilder { return &STSBuilder{b} }

// AssumedRole constructs an ARN for an assumed role session.
func (b *STSBuilder) AssumedRole(roleName, sessionName string) string {
	return b.Build("sts", "assumed-role/"+roleName+"/"+sessionName)
}

// FederatedUser constructs an ARN for a federated user.
func (b *STSBuilder) FederatedUser(name string) string { return b.Build("sts", "federated-user/"+name) }

// User constructs an ARN for an STS user.
func (b *STSBuilder) User(name string) string { return b.Build("sts", "user/"+name) }

// KMSBuilder provides methods for constructing KMS (Key Management Service) ARNs.
type KMSBuilder struct{ *ARNBuilder }

// KMS returns a KMSBuilder for constructing KMS ARNs.
func (b *ARNBuilder) KMS() *KMSBuilder { return &KMSBuilder{b} }

// Key constructs an ARN for a KMS key.
func (b *KMSBuilder) Key(keyId string) string { return b.Build("kms", "key/"+keyId) }

// Alias constructs an ARN for a KMS alias.
func (b *KMSBuilder) Alias(aliasName string) string {
	if !strings.HasPrefix(aliasName, "alias/") {
		aliasName = "alias/" + aliasName
	}
	return b.Build("kms", aliasName)
}

// ParseKeyID extracts the key ID from a KMS key ARN or returns the input if it's already a key ID.
func (b *KMSBuilder) ParseKeyID(keyOrArn string) string {
	if !strings.HasPrefix(keyOrArn, "arn:") {
		return keyOrArn
	}
	parts := strings.Split(keyOrArn, ":")
	if len(parts) >= 6 && strings.HasPrefix(parts[5], "key/") {
		return strings.TrimPrefix(parts[5], "key/")
	}
	return keyOrArn
}

// ParseAliasName extracts the alias name from a KMS alias ARN or returns the properly formatted alias name.
func (b *KMSBuilder) ParseAliasName(aliasOrArn string) string {
	if strings.HasPrefix(aliasOrArn, "arn:") {
		parts := strings.Split(aliasOrArn, ":")
		if len(parts) >= 6 {
			return parts[5]
		}
	}
	if !strings.HasPrefix(aliasOrArn, "alias/") {
		return "alias/" + aliasOrArn
	}
	return aliasOrArn
}

// IsAlias returns true if the given name is a KMS alias.
func (b *KMSBuilder) IsAlias(name string) bool { return strings.HasPrefix(name, "alias/") }

// ACMBuilder provides methods for constructing ACM (AWS Certificate Manager) ARNs.
type ACMBuilder struct{ *ARNBuilder }

// ACM returns an ACMBuilder for constructing ACM ARNs.
func (b *ARNBuilder) ACM() *ACMBuilder { return &ACMBuilder{b} }

// Certificate constructs an ARN for an ACM certificate.
func (b *ACMBuilder) Certificate(id string) string { return b.Build("acm", "certificate/"+id) }

// ParseCertificateID extracts the certificate ID from a certificate ARN.
func (b *ACMBuilder) ParseCertificateID(arn string) string {
	if !strings.HasPrefix(arn, "arn:") {
		return arn
	}
	parts := strings.Split(arn, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return arn
}

// SecretsManagerBuilder provides methods for constructing Secrets Manager ARNs.
type SecretsManagerBuilder struct{ *ARNBuilder }

// SecretsManager returns a SecretsManagerBuilder for constructing Secrets Manager ARNs.
func (b *ARNBuilder) SecretsManager() *SecretsManagerBuilder { return &SecretsManagerBuilder{b} }

// Secret constructs an ARN for a Secrets Manager secret.
func (b *SecretsManagerBuilder) Secret(name string) string {
	return b.Build("secretsmanager", "secret:"+name)
}

// SecretVersion constructs an ARN for a specific version of a Secrets Manager secret.
func (b *SecretsManagerBuilder) SecretVersion(name, versionId string) string {
	return b.Build("secretsmanager", "secret:"+name+":"+versionId)
}

// ParseSecretName extracts the secret name from a Secrets Manager secret ARN.
func (b *SecretsManagerBuilder) ParseSecretName(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "secret:") {
		parts := strings.Split(strings.TrimPrefix(resource, "secret:"), ":")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

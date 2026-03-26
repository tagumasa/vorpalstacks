package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.KMSBuilder for backward compatibility.
type ARNBuilder struct {
	builder *svcarn.KMSBuilder
}

// NewARNBuilder creates a new KMS ARN builder with the specified account ID and region.
func NewARNBuilder(accountId, region string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).KMS(),
	}
}

// KeyArn returns the ARN for a KMS key.
func (b *ARNBuilder) KeyArn(keyID string) string {
	return b.builder.Key(keyID)
}

// AliasArn returns the ARN for a KMS key alias.
func (b *ARNBuilder) AliasArn(aliasName string) string {
	return b.builder.Alias(aliasName)
}

// ParseKeyID extracts the key ID from a key ID or ARN.
func (b *ARNBuilder) ParseKeyID(keyIDOrArn string) string {
	return b.builder.ParseKeyID(keyIDOrArn)
}

// ParseAliasName extracts the alias name from an alias name or ARN.
func (b *ARNBuilder) ParseAliasName(aliasNameOrArn string) string {
	return b.builder.ParseAliasName(aliasNameOrArn)
}

// IsAlias checks if the given name is an alias (starts with "alias/").
func (b *ARNBuilder) IsAlias(name string) bool {
	return b.builder.IsAlias(name)
}

// AccountId returns the account ID used by the ARN builder.
func (b *ARNBuilder) AccountId() string {
	return b.builder.AccountId()
}

// Region returns the region used by the ARN builder.
func (b *ARNBuilder) Region() string {
	return b.builder.Region()
}

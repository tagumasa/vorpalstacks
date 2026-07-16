// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

// IAMBuilder provides methods for constructing IAM resource ARNs.
type IAMBuilder struct{ *ARNBuilder }

// IAM returns an IAMBuilder for constructing IAM ARNs.
func (b *ARNBuilder) IAM() *IAMBuilder { return &IAMBuilder{b} }

// Root constructs an ARN for the root IAM user of the account.
func (b *IAMBuilder) Root() string { return b.BuildNoRegion("iam", "root") }

// OIDCProvider constructs an ARN for an IAM OpenID Connect provider.
func (b *IAMBuilder) OIDCProvider(url string) string {
	return b.BuildNoRegion("iam", "oidc-provider/"+url)
}

// User constructs an ARN for an IAM user with the given name.
func (b *IAMBuilder) User(name string) string { return b.BuildNoRegion("iam", "user/"+name) }

// Role constructs an ARN for an IAM role with the given name.
func (b *IAMBuilder) Role(name string) string { return b.BuildNoRegion("iam", "role/"+name) }

// Group constructs an ARN for an IAM group with the given name.
func (b *IAMBuilder) Group(name string) string { return b.BuildNoRegion("iam", "group/"+name) }

// Policy constructs an ARN for an IAM policy with the given name.
func (b *IAMBuilder) Policy(name string) string { return b.BuildNoRegion("iam", "policy/"+name) }

// InstanceProfile constructs an ARN for an IAM instance profile with the given name.
func (b *IAMBuilder) InstanceProfile(name string) string {
	return b.BuildNoRegion("iam", "instance-profile/"+name)
}

// AccessKey constructs an ARN for an IAM access key.
func (b *IAMBuilder) AccessKey(keyId string) string {
	return b.BuildNoRegion("iam", "user/unknown/accessKey/"+keyId)
}

// MFADevice constructs an ARN for an IAM MFA device.
func (b *IAMBuilder) MFADevice(serial string) string { return b.BuildNoRegion("iam", "mfa/"+serial) }

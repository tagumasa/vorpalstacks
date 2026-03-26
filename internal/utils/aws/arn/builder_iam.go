// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

// IAMBuilder provides methods for constructing IAM resource ARNs.
type IAMBuilder struct{ *ARNBuilder }

// IAM returns an IAMBuilder for constructing IAM ARNs.
func (b *ARNBuilder) IAM() *IAMBuilder { return &IAMBuilder{b} }

// User constructs an ARN for an IAM user with the given name.
func (b *IAMBuilder) User(name string) string { return b.Build("iam", "user/"+name) }

// Role constructs an ARN for an IAM role with the given name.
func (b *IAMBuilder) Role(name string) string { return b.Build("iam", "role/"+name) }

// Group constructs an ARN for an IAM group with the given name.
func (b *IAMBuilder) Group(name string) string { return b.Build("iam", "group/"+name) }

// Policy constructs an ARN for an IAM policy with the given name.
func (b *IAMBuilder) Policy(name string) string { return b.Build("iam", "policy/"+name) }

// InstanceProfile constructs an ARN for an IAM instance profile with the given name.
func (b *IAMBuilder) InstanceProfile(name string) string {
	return b.Build("iam", "instance-profile/"+name)
}

// AccessKey constructs an ARN for an IAM access key.
func (b *IAMBuilder) AccessKey(keyId string) string {
	return b.Build("iam", "user/unknown/accessKey/"+keyId)
}

// MFADevice constructs an ARN for an IAM MFA device.
func (b *IAMBuilder) MFADevice(serial string) string { return b.Build("iam", "mfa/"+serial) }

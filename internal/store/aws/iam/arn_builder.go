package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"fmt"
	"strings"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder provides IAM-specific ARN construction with path support
type ARNBuilder struct {
	accountId string
	region    string
	partition string
}

// NewARNBuilder creates a new IAM ARN builder for the given account ID.
func NewARNBuilder(accountId string) *ARNBuilder {
	return &ARNBuilder{
		accountId: accountId,
		region:    "",
		partition: "aws",
	}
}

// WithRegion sets the region for the ARN builder.
func (b *ARNBuilder) WithRegion(region string) *ARNBuilder {
	b.region = region
	return b
}

// WithPartition sets the partition for the ARN builder.
func (b *ARNBuilder) WithPartition(partition string) *ARNBuilder {
	b.partition = partition
	return b
}

// UserARN constructs an ARN for an IAM user.
func (b *ARNBuilder) UserARN(path, userName string) string {
	fullPath := b.normalizePath(path)
	return fmt.Sprintf("arn:%s:iam::%s:user%s%s", b.partition, b.accountId, fullPath, userName)
}

// GroupARN constructs an ARN for an IAM group.
func (b *ARNBuilder) GroupARN(path, groupName string) string {
	fullPath := b.normalizePath(path)
	return fmt.Sprintf("arn:%s:iam::%s:group%s%s", b.partition, b.accountId, fullPath, groupName)
}

// RoleARN constructs an ARN for an IAM role.
func (b *ARNBuilder) RoleARN(path, roleName string) string {
	fullPath := b.normalizePath(path)
	return fmt.Sprintf("arn:%s:iam::%s:role%s%s", b.partition, b.accountId, fullPath, roleName)
}

// PolicyARN constructs an ARN for an IAM policy.
func (b *ARNBuilder) PolicyARN(path, policyName string) string {
	fullPath := b.normalizePath(path)
	return fmt.Sprintf("arn:%s:iam::%s:policy%s%s", b.partition, b.accountId, fullPath, policyName)
}

// InstanceProfileARN constructs an ARN for an IAM instance profile.
func (b *ARNBuilder) InstanceProfileARN(path, instanceProfileName string) string {
	fullPath := b.normalizePath(path)
	return fmt.Sprintf("arn:%s:iam::%s:instance-profile%s%s", b.partition, b.accountId, fullPath, instanceProfileName)
}

// AccessKeyARN constructs an ARN for an IAM access key.
func (b *ARNBuilder) AccessKeyARN(userName, accessKeyId string) string {
	return fmt.Sprintf("arn:%s:iam::%s:user/%s/access-key/%s", b.partition, b.accountId, userName, accessKeyId)
}

// MFADeviceARN constructs an ARN for an MFA device.
func (b *ARNBuilder) MFADeviceARN(serialNumber string) string {
	return fmt.Sprintf("arn:%s:iam::%s:mfa/%s", b.partition, b.accountId, serialNumber)
}

// VirtualMFADeviceARN constructs an ARN for a virtual MFA device.
func (b *ARNBuilder) VirtualMFADeviceARN(serialNumber string) string {
	return fmt.Sprintf("arn:%s:iam::%s:mfa/%s", b.partition, b.accountId, serialNumber)
}

// ServerCertificateARN constructs an ARN for a server certificate.
func (b *ARNBuilder) ServerCertificateARN(certName string) string {
	return fmt.Sprintf("arn:%s:iam::%s:server-certificate/%s", b.partition, b.accountId, certName)
}

// SAMLProviderARN constructs an ARN for a SAML provider.
func (b *ARNBuilder) SAMLProviderARN(providerName string) string {
	return fmt.Sprintf("arn:%s:iam::%s:saml-provider/%s", b.partition, b.accountId, providerName)
}

// OpenIDConnectProviderARN constructs an ARN for an OIDC provider.
func (b *ARNBuilder) OpenIDConnectProviderARN(providerURL string) string {
	return fmt.Sprintf("arn:%s:iam::%s:oidc-provider/%s", b.partition, b.accountId, providerURL)
}

// AccountARN constructs an ARN for the AWS account.
func (b *ARNBuilder) AccountARN() string {
	return fmt.Sprintf("arn:%s:iam::%s:root", b.partition, b.accountId)
}

func (b *ARNBuilder) normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}

// ParseARN delegates to svcarn package
func ParseARN(arn string) (partition, service, region, accountId, resource string, err error) {
	p, s, r, a, res := svcarn.SplitARN(arn)
	if p == "" {
		err = fmt.Errorf("invalid ARN format: %s", arn)
		return
	}
	return p, s, r, a, res, nil
}

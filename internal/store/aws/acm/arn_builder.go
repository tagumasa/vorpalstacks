// Package acm provides ACM (AWS Certificate Manager) storage functionality for vorpalstacks.
package acm

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.ACMBuilder for backward compatibility
type ARNBuilder struct {
	builder *svcarn.ACMBuilder
}

// NewARNBuilder creates a new ARNBuilder for ACM resources.
func NewARNBuilder(accountId, region string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).ACM(),
	}
}

// BuildCertificateARN builds an ACM certificate ARN from a certificate ID.
func (b *ARNBuilder) BuildCertificateARN(certificateId string) string {
	return b.builder.Certificate(certificateId)
}

// ParseCertificateArn extracts the certificate ID from an ACM certificate ARN.
func (b *ARNBuilder) ParseCertificateArn(arn string) (string, error) {
	certId := b.builder.ParseCertificateID(arn)
	if certId == "" {
		return "", ErrInvalidArn
	}
	return certId, nil
}

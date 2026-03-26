// Package cloudfront provides CloudFront storage functionality for vorpalstacks.
package cloudfront

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.CloudFrontBuilder for backward compatibility
type ARNBuilder struct {
	builder *svcarn.CloudFrontBuilder
}

// NewARNBuilder creates a new ARNBuilder for CloudFront resources.
func NewARNBuilder(accountId string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, "").CloudFront(),
	}
}

// BuildDistributionARN builds a CloudFront distribution ARN from an ID.
func (b *ARNBuilder) BuildDistributionARN(id string) string {
	return b.builder.Distribution(id)
}

// BuildCachePolicyARN builds a CloudFront cache policy ARN from an ID.
func (b *ARNBuilder) BuildCachePolicyARN(id string) string {
	return b.builder.CachePolicy(id)
}

// BuildOriginRequestPolicyARN builds a CloudFront origin request policy ARN from an ID.
func (b *ARNBuilder) BuildOriginRequestPolicyARN(id string) string {
	return b.builder.OriginRequestPolicy(id)
}

// BuildOriginAccessControlARN builds a CloudFront origin access control ARN from an ID.
func (b *ARNBuilder) BuildOriginAccessControlARN(id string) string {
	return b.builder.OriginAccessControl(id)
}

// BuildResponseHeadersPolicyARN builds a CloudFront response headers policy ARN from an ID.
func (b *ARNBuilder) BuildResponseHeadersPolicyARN(id string) string {
	return b.builder.ResponseHeadersPolicy(id)
}

// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.APIGatewayBuilder for ARN construction
type ARNBuilder struct {
	builder *svcarn.APIGatewayBuilder
}

// NewARNBuilder creates a new ARN builder for API Gateway.
func NewARNBuilder(accountId, region string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).APIGateway(),
	}
}

// DomainNameArn returns the ARN for a domain name.
func (b *ARNBuilder) DomainNameArn(domainName string) string {
	return b.builder.DomainName(domainName)
}

// GenerateApiId generates a unique API ID.
func (b *ARNBuilder) GenerateApiId() string {
	return generateId("api", 10)
}

// GenerateResourceId generates a unique resource ID.
func (b *ARNBuilder) GenerateResourceId() string {
	return generateId("res", 6)
}

// GenerateDeploymentId generates a unique deployment ID.
func (b *ARNBuilder) GenerateDeploymentId() string {
	return generateId("dep", 6)
}

// GenerateValidatorId generates a unique validator ID.
func (b *ARNBuilder) GenerateValidatorId() string {
	return generateId("val", 6)
}

// GenerateModelId generates a unique model ID.
func (b *ARNBuilder) GenerateModelId() string {
	return generateId("mod", 6)
}

// GenerateAuthorizerId generates a unique authorizer ID.
func (b *ARNBuilder) GenerateAuthorizerId() string {
	return generateId("auth", 6)
}

// GenerateApiKeyId generates a unique API key ID.
func (b *ARNBuilder) GenerateApiKeyId() string {
	return generateId("apikey", 10)
}

// GenerateUsagePlanId generates a unique usage plan ID.
func (b *ARNBuilder) GenerateUsagePlanId() string {
	return generateId("usageplan", 10)
}

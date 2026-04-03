// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.APIGatewayBuilder for backward compatibility
type ARNBuilder struct {
	builder *svcarn.APIGatewayBuilder
}

// NewARNBuilder creates a new ARN builder for API Gateway.
func NewARNBuilder(accountId, region string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).APIGateway(),
	}
}

// RestApiArn returns the ARN for a REST API.
func (b *ARNBuilder) RestApiArn(apiId string) string {
	return b.builder.RestApi(apiId)
}

// ResourceArn returns the ARN for an API resource.
func (b *ARNBuilder) ResourceArn(apiId, resourceId string) string {
	return b.builder.Resource(apiId, resourceId)
}

// MethodArn returns the ARN for an API method.
func (b *ARNBuilder) MethodArn(apiId, resourceId, httpMethod string) string {
	return b.builder.Method(apiId, resourceId, httpMethod)
}

// DeploymentArn returns the ARN for an API deployment.
func (b *ARNBuilder) DeploymentArn(apiId, deploymentId string) string {
	return b.builder.Deployment(apiId, deploymentId)
}

// StageArn returns the ARN for an API stage.
func (b *ARNBuilder) StageArn(apiId, stageName string) string {
	return b.builder.Stage(apiId, stageName)
}

// ApiKeyArn returns the ARN for an API key.
func (b *ARNBuilder) ApiKeyArn(apiKeyId string) string {
	return b.builder.ApiKey(apiKeyId)
}

// UsagePlanArn returns the ARN for a usage plan.
func (b *ARNBuilder) UsagePlanArn(usagePlanId string) string {
	return b.builder.UsagePlan(usagePlanId)
}

// ModelArn returns the ARN for an API model.
func (b *ARNBuilder) ModelArn(apiId, modelName string) string {
	return b.builder.Model(apiId, modelName)
}

// RequestValidatorArn returns the ARN for a request validator.
func (b *ARNBuilder) RequestValidatorArn(apiId, validatorId string) string {
	return b.builder.RequestValidator(apiId, validatorId)
}

// AuthorizerArn returns the ARN for an authorizer.
func (b *ARNBuilder) AuthorizerArn(apiId, authorizerId string) string {
	return b.builder.Authorizer(apiId, authorizerId)
}

// DomainNameArn returns the ARN for a domain name.
func (b *ARNBuilder) DomainNameArn(domainName string) string {
	return b.builder.DomainName(domainName)
}

// BasePathMappingArn returns the ARN for a base path mapping.
func (b *ARNBuilder) BasePathMappingArn(domainName, basePath string) string {
	return b.builder.BasePathMapping(domainName, basePath)
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

// Package apigateway provides API Gateway storage functionality for vorpalstacks.
package apigateway

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ArnBuilder wraps svcarn.APIGatewayBuilder for backward compatibility
type ArnBuilder struct {
	builder *svcarn.APIGatewayBuilder
}

// NewArnBuilder creates a new ARN builder for API Gateway.
func NewArnBuilder(accountId, region string) *ArnBuilder {
	return &ArnBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).APIGateway(),
	}
}

// RestApiArn returns the ARN for a REST API.
func (b *ArnBuilder) RestApiArn(apiId string) string {
	return b.builder.RestApi(apiId)
}

// ResourceArn returns the ARN for an API resource.
func (b *ArnBuilder) ResourceArn(apiId, resourceId string) string {
	return b.builder.Resource(apiId, resourceId)
}

// MethodArn returns the ARN for an API method.
func (b *ArnBuilder) MethodArn(apiId, resourceId, httpMethod string) string {
	return b.builder.Method(apiId, resourceId, httpMethod)
}

// DeploymentArn returns the ARN for an API deployment.
func (b *ArnBuilder) DeploymentArn(apiId, deploymentId string) string {
	return b.builder.Deployment(apiId, deploymentId)
}

// StageArn returns the ARN for an API stage.
func (b *ArnBuilder) StageArn(apiId, stageName string) string {
	return b.builder.Stage(apiId, stageName)
}

// ApiKeyArn returns the ARN for an API key.
func (b *ArnBuilder) ApiKeyArn(apiKeyId string) string {
	return b.builder.ApiKey(apiKeyId)
}

// UsagePlanArn returns the ARN for a usage plan.
func (b *ArnBuilder) UsagePlanArn(usagePlanId string) string {
	return b.builder.UsagePlan(usagePlanId)
}

// ModelArn returns the ARN for an API model.
func (b *ArnBuilder) ModelArn(apiId, modelName string) string {
	return b.builder.Model(apiId, modelName)
}

// RequestValidatorArn returns the ARN for a request validator.
func (b *ArnBuilder) RequestValidatorArn(apiId, validatorId string) string {
	return b.builder.RequestValidator(apiId, validatorId)
}

// AuthorizerArn returns the ARN for an authorizer.
func (b *ArnBuilder) AuthorizerArn(apiId, authorizerId string) string {
	return b.builder.Authorizer(apiId, authorizerId)
}

// DomainNameArn returns the ARN for a domain name.
func (b *ArnBuilder) DomainNameArn(domainName string) string {
	return b.builder.DomainName(domainName)
}

// BasePathMappingArn returns the ARN for a base path mapping.
func (b *ArnBuilder) BasePathMappingArn(domainName, basePath string) string {
	return b.builder.BasePathMapping(domainName, basePath)
}

// GenerateApiId generates a unique API ID.
func (b *ArnBuilder) GenerateApiId() string {
	return generateId("api", 10)
}

// GenerateResourceId generates a unique resource ID.
func (b *ArnBuilder) GenerateResourceId() string {
	return generateId("res", 6)
}

// GenerateDeploymentId generates a unique deployment ID.
func (b *ArnBuilder) GenerateDeploymentId() string {
	return generateId("dep", 6)
}

// GenerateValidatorId generates a unique validator ID.
func (b *ArnBuilder) GenerateValidatorId() string {
	return generateId("val", 6)
}

// GenerateModelId generates a unique model ID.
func (b *ArnBuilder) GenerateModelId() string {
	return generateId("mod", 6)
}

// GenerateAuthorizerId generates a unique authorizer ID.
func (b *ArnBuilder) GenerateAuthorizerId() string {
	return generateId("auth", 6)
}

// GenerateApiKeyId generates a unique API key ID.
func (b *ArnBuilder) GenerateApiKeyId() string {
	return generateId("apikey", 10)
}

// GenerateUsagePlanId generates a unique usage plan ID.
func (b *ArnBuilder) GenerateUsagePlanId() string {
	return generateId("usageplan", 10)
}

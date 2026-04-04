// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import (
	"fmt"
	"strings"
)

// APIGatewayBuilder provides methods for constructing API Gateway ARNs.
type APIGatewayBuilder struct{ *ARNBuilder }

// APIGateway returns an APIGatewayBuilder for constructing API Gateway ARNs.
func (b *ARNBuilder) APIGateway() *APIGatewayBuilder { return &APIGatewayBuilder{b} }

// RestApi constructs an ARN for an API Gateway REST API.
func (b *APIGatewayBuilder) RestApi(id string) string { return b.Build("apigateway", "/restapis/"+id) }

// Resource constructs an ARN for an API Gateway resource.
func (b *APIGatewayBuilder) Resource(apiId, resId string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/resources/%s", apiId, resId))
}

// Method constructs an ARN for an API Gateway method.
func (b *APIGatewayBuilder) Method(apiId, resId, method string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/resources/%s/methods/%s", apiId, resId, method))
}

// Deployment constructs an ARN for an API Gateway deployment.
func (b *APIGatewayBuilder) Deployment(apiId, depId string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/deployments/%s", apiId, depId))
}

// Stage constructs an ARN for an API Gateway stage.
func (b *APIGatewayBuilder) Stage(apiId, name string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/stages/%s", apiId, name))
}

// ApiKey constructs an ARN for an API Gateway API key.
func (b *APIGatewayBuilder) ApiKey(id string) string { return b.Build("apigateway", "/apikeys/"+id) }

// UsagePlan constructs an ARN for an API Gateway usage plan.
func (b *APIGatewayBuilder) UsagePlan(id string) string {
	return b.Build("apigateway", "/usageplans/"+id)
}

// Model constructs an ARN for an API Gateway model.
func (b *APIGatewayBuilder) Model(apiId, name string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/models/%s", apiId, name))
}

// Authorizer constructs an ARN for an API Gateway authorizer.
func (b *APIGatewayBuilder) Authorizer(apiId, id string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/authorizers/%s", apiId, id))
}

// DomainName constructs an ARN for an API Gateway domain name.
func (b *APIGatewayBuilder) DomainName(name string) string {
	return b.Build("apigateway", "/domainnames/"+name)
}

// RequestValidator constructs an ARN for an API Gateway request validator.
func (b *APIGatewayBuilder) RequestValidator(apiId, id string) string {
	return b.Build("apigateway", fmt.Sprintf("/restapis/%s/requestvalidators/%s", apiId, id))
}

// BasePathMapping constructs an ARN for an API Gateway base path mapping.
func (b *APIGatewayBuilder) BasePathMapping(domainName, basePath string) string {
	return b.Build("apigateway", fmt.Sprintf("/domainnames/%s/basepathmappings/%s", domainName, basePath))
}

// CognitoBuilder provides methods for constructing Cognito ARNs.
type CognitoBuilder struct{ *ARNBuilder }

// Cognito returns a CognitoBuilder for constructing Cognito ARNs.
func (b *ARNBuilder) Cognito() *CognitoBuilder { return &CognitoBuilder{b} }

// UserPool constructs an ARN for a Cognito User Pool.
func (b *CognitoBuilder) UserPool(id string) string { return b.Build("cognito-idp", "userpool/"+id) }

// User constructs an ARN for a Cognito User Pool user.
func (b *CognitoBuilder) User(poolId, username string) string {
	return b.Build("cognito-idp", fmt.Sprintf("userpool/%s/user/%s", poolId, username))
}

// Group constructs an ARN for a Cognito User Pool group.
func (b *CognitoBuilder) Group(poolId, name string) string {
	return b.Build("cognito-idp", fmt.Sprintf("userpool/%s/group/%s", poolId, name))
}

// IdentityPool constructs an ARN for a Cognito Identity Pool.
func (b *CognitoBuilder) IdentityPool(id string) string {
	return b.Build("cognito-identity", "identitypool/"+id)
}

// IdentityProvider constructs an ARN for a Cognito Identity Provider.
func (b *CognitoBuilder) IdentityProvider(poolId, name string) string {
	return b.Build("cognito-idp", fmt.Sprintf("userpool/%s/identityprovider/%s", poolId, name))
}

// UserPoolClient constructs an ARN for a Cognito User Pool client.
func (b *CognitoBuilder) UserPoolClient(poolId, clientId string) string {
	return b.Build("cognito-idp", fmt.Sprintf("userpool/%s/client/%s", poolId, clientId))
}

// ParseUserPoolID extracts the user pool ID from a Cognito User Pool ARN.
func (b *CognitoBuilder) ParseUserPoolID(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "userpool/") {
		return strings.TrimPrefix(resource, "userpool/")
	}
	return ""
}

// SESv2Builder provides methods for constructing SESv2 (Simple Email Service v2) ARNs.
type SESv2Builder struct{ *ARNBuilder }

// SESv2 returns an SESv2Builder for constructing SESv2 ARNs.
func (b *ARNBuilder) SESv2() *SESv2Builder { return &SESv2Builder{b} }

// Identity constructs an ARN for an SESv2 identity.
func (b *SESv2Builder) Identity(name string) string { return b.Build("ses", "identity/"+name) }

// ConfigSet constructs an ARN for an SESv2 configuration set.
func (b *SESv2Builder) ConfigSet(name string) string {
	return b.Build("ses", "configuration-set/"+name)
}

// Template constructs an ARN for an SESv2 template.
func (b *SESv2Builder) Template(name string) string { return b.Build("ses", "template/"+name) }

// ContactList constructs an ARN for an SESv2 contact list.
func (b *SESv2Builder) ContactList(name string) string { return b.Build("ses", "contact-list/"+name) }

// ParseIdentity extracts the identity name from an SESv2 identity ARN.
func (b *SESv2Builder) ParseIdentity(arn string) string { return extractSuffix(arn, "identity/") }

func extractSuffix(arn, prefix string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, prefix) {
		return strings.TrimPrefix(resource, prefix)
	}
	return ""
}

// AppSyncBuilder provides methods for constructing AppSync ARNs.
type AppSyncBuilder struct{ *ARNBuilder }

// AppSync returns an AppSyncBuilder for constructing AppSync ARNs.
func (b *ARNBuilder) AppSync() *AppSyncBuilder { return &AppSyncBuilder{b} }

// Api constructs an ARN for an AppSync API.
func (b *AppSyncBuilder) Api(id string) string { return b.Build("appsync", "apis/"+id) }

// ChannelNamespace constructs an ARN for an AppSync channel namespace.
func (b *AppSyncBuilder) ChannelNamespace(apiId, name string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/channelNamespaces/%s", apiId, name))
}

// GraphQLApi constructs an ARN for an AppSync GraphQL API.
func (b *AppSyncBuilder) GraphQLApi(id string) string { return b.Build("appsync", "apis/"+id) }

// DataSource constructs an ARN for an AppSync data source.
func (b *AppSyncBuilder) DataSource(apiId, name string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/datasources/%s", apiId, name))
}

// Resolver constructs an ARN for an AppSync resolver.
func (b *AppSyncBuilder) Resolver(apiId, typeName, fieldName string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/types/%s/resolvers/%s", apiId, typeName, fieldName))
}

// Function constructs an ARN for an AppSync function.
func (b *AppSyncBuilder) Function(apiId, functionId string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/functions/%s", apiId, functionId))
}

// DomainName constructs an ARN for an AppSync domain name.
func (b *AppSyncBuilder) DomainName(name string) string {
	return b.Build("appsync", "domainnames/"+name)
}

// ApiKey constructs an ARN for an AppSync API key.
func (b *AppSyncBuilder) ApiKey(apiId, id string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/apikeys/%s", apiId, id))
}

// Type constructs an ARN for an AppSync type.
func (b *AppSyncBuilder) Type(apiId, typeName string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/types/%s", apiId, typeName))
}

// ApiCache constructs an ARN for an AppSync API cache.
func (b *AppSyncBuilder) ApiCache(apiId string) string {
	return b.Build("appsync", fmt.Sprintf("apis/%s/ApiCaches", apiId))
}

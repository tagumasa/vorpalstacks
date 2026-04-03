// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import (
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ARNBuilder wraps svcarn.LambdaBuilder for backward compatibility.
type ARNBuilder struct {
	builder *svcarn.LambdaBuilder
}

// NewARNBuilder creates a new Lambda ARN builder with the specified account ID and region.
func NewARNBuilder(accountId, region string) *ARNBuilder {
	return &ARNBuilder{
		builder: svcarn.NewARNBuilder(accountId, region).Lambda(),
	}
}

// FunctionArn returns the ARN for a Lambda function.
func (b *ARNBuilder) FunctionArn(functionName string) string {
	return b.builder.Function(functionName)
}

// FunctionVersionArn returns the ARN for a specific version of a Lambda function.
func (b *ARNBuilder) FunctionVersionArn(functionName, version string) string {
	return b.builder.FunctionVersion(functionName, version)
}

// FunctionAliasArn returns the ARN for a Lambda function alias.
func (b *ARNBuilder) FunctionAliasArn(functionName, aliasName string) string {
	return b.builder.FunctionAlias(functionName, aliasName)
}

// LayerArn returns the ARN for a Lambda layer.
func (b *ARNBuilder) LayerArn(layerName string) string {
	return b.builder.Layer(layerName)
}

// LayerVersionArn returns the ARN for a specific version of a Lambda layer.
func (b *ARNBuilder) LayerVersionArn(layerName string, version int64) string {
	return b.builder.LayerVersion(layerName, version)
}

// EventSourceMappingArn returns the ARN for an event source mapping.
func (b *ARNBuilder) EventSourceMappingArn(uuid string) string {
	return b.builder.EventSourceMapping(uuid)
}

// CodeSigningConfigArn returns the ARN for a code signing configuration.
func (b *ARNBuilder) CodeSigningConfigArn(configId string) string {
	return b.builder.CodeSigningConfig(configId)
}

// ParseFunctionNameFromArn extracts the function name from a Lambda function ARN.
func (b *ARNBuilder) ParseFunctionNameFromArn(arn string) string {
	return b.builder.ParseFunctionName(arn)
}

// ParseLayerNameFromArn extracts the layer name from a Lambda layer ARN.
func (b *ARNBuilder) ParseLayerNameFromArn(arn string) string {
	return b.builder.ParseLayerName(arn)
}

// ParseLayerVersionFromArn extracts the layer version from a Lambda layer ARN.
func (b *ARNBuilder) ParseLayerVersionFromArn(arn string) int64 {
	return b.builder.ParseLayerVersion(arn)
}

// AccountId returns the account ID used by the ARN builder.
func (b *ARNBuilder) AccountId() string {
	return b.builder.AccountId()
}

// Region returns the region used by the ARN builder.
func (b *ARNBuilder) Region() string {
	return b.builder.Region()
}

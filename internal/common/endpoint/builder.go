// Package endpoint provides URL builders for AWS service endpoints.
package endpoint

import (
	"fmt"
	"strings"

	appconfig "vorpalstacks/internal/config"
)

// Builder constructs URLs for AWS service endpoints.
type Builder struct {
	baseURL string
	region  string
}

// NewBuilder creates a new endpoint builder with default configuration.
func NewBuilder() *Builder {
	return &Builder{
		baseURL: appconfig.BaseURL(),
		region:  appconfig.AWSRegion(),
	}
}

// NewBuilderWithConfig creates a new endpoint builder with custom configuration.
func NewBuilderWithConfig(baseURL, region string) *Builder {
	return &Builder{
		baseURL: baseURL,
		region:  region,
	}
}

// BaseURL returns the base URL for the endpoint builder.
func (b *Builder) BaseURL() string {
	return b.baseURL
}

// SQSQueueURL constructs the URL for an SQS queue.
func (b *Builder) SQSQueueURL(accountID, queueName string) string {
	return fmt.Sprintf("%s/%s/%s", b.baseURL, accountID, queueName)
}

// S3WebsiteURL constructs the URL for an S3 bucket configured as a website.
func (b *Builder) S3WebsiteURL(bucket, region string) string {
	port, _ := appconfig.GetResourcePort("ports.s3_website", bucket)
	suffix := appconfig.GetString("endpoints.s3_website_suffix")
	domain := strings.Replace(suffix, "{region}", region, 1)
	if port != 0 && port != 80 {
		return fmt.Sprintf("http://%s:%d", domain, port)
	}
	return fmt.Sprintf("http://%s", domain)
}

// APIGatewayInvokeURL constructs the invoke URL for an API Gateway endpoint.
func (b *Builder) APIGatewayInvokeURL(apiID, stage, region string) string {
	port, _ := appconfig.GetResourcePort("ports.apigateway", apiID+"_"+stage)
	suffix := appconfig.GetString("endpoints.apigateway_suffix")
	domain := strings.Replace(suffix, "{region}", region, 1)
	domain = strings.Replace(domain, "{api_id}", apiID, 1)
	if port != 0 && port != 80 {
		return fmt.Sprintf("https://%s:%d/%s", domain, port, stage)
	}
	return fmt.Sprintf("https://%s/%s", domain, stage)
}

// CognitoHostedUIURL constructs the URL for a Cognito hosted UI.
func (b *Builder) CognitoHostedUIURL(userPoolDomain, region string) string {
	port, _ := appconfig.GetResourcePort("ports.cognito_hosted", userPoolDomain)
	suffix := appconfig.GetString("endpoints.cognito_suffix")
	domain := strings.Replace(suffix, "{region}", region, 1)
	if port != 0 && port != 80 {
		return fmt.Sprintf("https://%s.%s:%d", userPoolDomain, domain, port)
	}
	return fmt.Sprintf("https://%s.%s", userPoolDomain, domain)
}

// CloudFrontURL constructs the URL for a CloudFront distribution.
func (b *Builder) CloudFrontURL(distributionID string) string {
	port, _ := appconfig.GetResourcePort("ports.cloudfront", distributionID)
	if port != 0 && port != 80 {
		return fmt.Sprintf("https://%s.cloudfront.net:%d", distributionID, port)
	}
	return fmt.Sprintf("https://%s.cloudfront.net", distributionID)
}

// LambdaFunctionURL constructs the URL for a Lambda function URL.
func (b *Builder) LambdaFunctionURL(functionName, region string) string {
	port, _ := appconfig.GetResourcePort("ports.lambda_url", functionName)
	if port != 0 && port != 80 {
		return fmt.Sprintf("https://%s.lambda-url.%s.on.aws:%d", functionName, region, port)
	}
	return fmt.Sprintf("https://%s.lambda-url.%s.on.aws", functionName, region)
}

var defaultBuilder *Builder

// Initialise initializes the default endpoint builder.
func Initialise() {
	defaultBuilder = NewBuilder()
}

// GetBuilder returns the default endpoint builder.
func GetBuilder() *Builder {
	if defaultBuilder == nil {
		Initialise()
	}
	return defaultBuilder
}

// SQSQueueURL constructs the URL for an SQS queue using the default builder.
func SQSQueueURL(accountID, queueName string) string {
	return GetBuilder().SQSQueueURL(accountID, queueName)
}

// CloudFrontURL constructs the URL for a CloudFront distribution using the default builder.
func CloudFrontURL(distributionID string) string {
	return GetBuilder().CloudFrontURL(distributionID)
}

// LambdaFunctionURL constructs the URL for a Lambda function URL using the default builder.
func LambdaFunctionURL(functionName, region string) string {
	return GetBuilder().LambdaFunctionURL(functionName, region)
}

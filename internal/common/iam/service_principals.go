// Package iam provides IAM service operations for vorpalstacks.
package iam

// ServicePrincipal represents an AWS service principal.
type ServicePrincipal string

// Service principal constants for AWS services.
const (
	// ServicePrincipalLambda is the service principal for AWS Lambda.
	ServicePrincipalLambda ServicePrincipal = "lambda.amazonaws.com"
	// ServicePrincipalStates is the service principal for AWS Step Functions.
	ServicePrincipalStates ServicePrincipal = "states.amazonaws.com"
	// ServicePrincipalEvents is the service principal for Amazon EventBridge.
	ServicePrincipalEvents ServicePrincipal = "events.amazonaws.com"
	// ServicePrincipalScheduler is the service principal for Amazon EventBridge Scheduler.
	ServicePrincipalScheduler ServicePrincipal = "scheduler.amazonaws.com"
	// ServicePrincipalCloudTrail is the service principal for AWS CloudTrail.
	ServicePrincipalCloudTrail ServicePrincipal = "cloudtrail.amazonaws.com"
	// ServicePrincipalTimestream is the service principal for Amazon Timestream.
	ServicePrincipalTimestream ServicePrincipal = "timestream.amazonaws.com"
	// ServicePrincipalCognito is the service principal for Amazon Cognito.
	ServicePrincipalCognito ServicePrincipal = "cognito-identity.amazonaws.com"
)

// GetServicePrincipal returns the service principal for a given AWS service name.
func GetServicePrincipal(serviceName string) ServicePrincipal {
	switch serviceName {
	case "lambda":
		return ServicePrincipalLambda
	case "states":
		return ServicePrincipalStates
	case "events":
		return ServicePrincipalEvents
	case "scheduler":
		return ServicePrincipalScheduler
	case "cloudtrail":
		return ServicePrincipalCloudTrail
	case "timestream":
		return ServicePrincipalTimestream
	case "cognito", "cognito-identity":
		return ServicePrincipalCognito
	default:
		return ServicePrincipal(serviceName + ".amazonaws.com")
	}
}

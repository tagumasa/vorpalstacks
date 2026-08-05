package lambda

import (
	"fmt"
	"regexp"
	"strings"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Smithy-derived regex patterns
// ---------------------------------------------------------------------------

// functionNamePattern validates function names per AWS Lambda docs:
// alphanumeric, hyphens, and underscores, 1-64 characters.
var functionNamePattern = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// ---------------------------------------------------------------------------
// Runtime validation (Smithy Runtime enum)
// ---------------------------------------------------------------------------

// validRuntimes contains every Lambda runtime recognised by this platform.
// Values are lowercase canonical forms; callers compare case-insensitively.
var validRuntimes = []string{
	"nodejs24.x", "nodejs22.x",
	"python3.14", "python3.13", "python3.12", "python3.11", "python3.10",
	"java25", "java21", "java17", "java11", "java8.al2",
	"dotnet10", "dotnet9", "dotnet8",
	"ruby4.0", "ruby3.4", "ruby3.3",
	"provided.al2023", "provided.al2",
}

// ValidateRuntime checks if the provided runtime is a valid Lambda runtime.
func ValidateRuntime(runtime string) bool {
	runtimeLower := strings.ToLower(runtime)
	for _, r := range validRuntimes {
		if r == runtimeLower {
			return true
		}
	}
	return false
}

// ValidateHandler validates the handler string for a Lambda function.
// Checks that the handler is not empty and conforms to runtime-specific
// format requirements.
func ValidateHandler(runtime, handler string) error {
	if handler == "" {
		return NewInvalidParameter("Handler", "Handler cannot be empty")
	}

	if strings.HasPrefix(runtime, "python") {
		if !strings.Contains(handler, ".") {
			return NewInvalidParameter("Handler", "Python handler must be in the format module.function")
		}
	}

	if strings.HasPrefix(runtime, "nodejs") {
		if !strings.Contains(handler, ".") {
			return NewInvalidParameter("Handler", "Node.js handler must be in the format file.function")
		}
	}

	if strings.HasPrefix(runtime, "java") {
		if !strings.Contains(handler, "::") && !strings.Contains(handler, ".") {
			return NewInvalidParameter("Handler", "Java handler must be in the format package.Class::method")
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Function name validation
// ---------------------------------------------------------------------------

func validateFunctionName(name string) error {
	if len(name) == 0 || len(name) > 64 {
		return NewInvalidParameter("FunctionName", "Function name must be between 1 and 64 characters")
	}
	if !functionNamePattern.MatchString(name) {
		return NewInvalidParameter("FunctionName", "Function name can only contain alphanumeric characters, hyphens, and underscores")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Timeout and MemorySize validation (Smithy range traits)
// ---------------------------------------------------------------------------

func validateTimeout(timeout int32) error {
	if timeout < 1 || timeout > 900 {
		return NewInvalidParameter("Timeout", "Timeout must be between 1 and 900 seconds")
	}
	return nil
}

func validateMemorySize(memorySize int32) error {
	if memorySize < 128 || memorySize > 10240 {
		return NewInvalidParameter("MemorySize", "MemorySize must be between 128 and 10240 MB")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Event source mapping validation
// ---------------------------------------------------------------------------

// validateEventSourceArn checks that the ARN refers to a supported event
// source service (SQS, Kinesis, or DynamoDB streams). Other services
// (Kafka, MSK, DocumentDB) are accepted by AWS but not polled by this
// implementation; rejecting them prevents silent no-op mappings.
func validateEventSourceArn(arn string) error {
	service := arnutil.GetServiceFromARN(arn)
	switch service {
	case "sqs", "kinesis", "dynamodb":
		return nil
	default:
		return NewInvalidParameter("EventSourceArn",
			fmt.Sprintf("Unsupported event source service %q: only sqs, kinesis, and dynamodb are supported", service))
	}
}

// validateStartingPosition validates the EventSourcePosition enum per
// the Smithy model (TRIM_HORIZON, LATEST, AT_TIMESTAMP).
func validateStartingPosition(pos string) error {
	switch pos {
	case "TRIM_HORIZON", "LATEST", "AT_TIMESTAMP":
		return nil
	default:
		return NewInvalidParameter("StartingPosition",
			fmt.Sprintf("StartingPosition must be one of TRIM_HORIZON, LATEST, AT_TIMESTAMP; got %q", pos))
	}
}

// validateStartingPositionForStream checks that StartingPosition is
// provided when the event source is a Kinesis or DynamoDB stream.
// AWS Lambda docs: "Required for Amazon Kinesis, Amazon MSK, and DynamoDB
// Streams sources."
func validateStartingPositionForStream(startingPosition, eventSourceArn string) error {
	service := arnutil.GetServiceFromARN(eventSourceArn)
	if (service == "kinesis" || service == "dynamodb") && startingPosition == "" {
		return NewInvalidParameter("StartingPosition",
			"StartingPosition is required for Kinesis and DynamoDB stream sources")
	}
	return nil
}

// validateStartingPositionTimestamp checks that StartingPositionTimestamp
// is provided when StartingPosition is AT_TIMESTAMP.
func validateStartingPositionTimestamp(startingPosition string, hasTimestamp bool) error {
	if startingPosition == "AT_TIMESTAMP" && !hasTimestamp {
		return NewInvalidParameter("StartingPositionTimestamp",
			"StartingPositionTimestamp is required when StartingPosition is AT_TIMESTAMP")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Function URL config validation (Smithy AuthType required, InvokeMode enum)
// ---------------------------------------------------------------------------

// validateAuthType checks that AuthType is provided (Smithy REQUIRED) and
// is one of NONE or AWS_IAM. An empty value is rejected rather than
// silently defaulting to the insecure NONE option.
func validateAuthType(authType string) error {
	if authType == "" {
		return NewInvalidParameter("AuthType", "AuthType is required (NONE or AWS_IAM)")
	}
	if authType != "NONE" && authType != "AWS_IAM" {
		return NewInvalidParameter("AuthType",
			fmt.Sprintf("AuthType must be NONE or AWS_IAM; got %q", authType))
	}
	return nil
}

// validateInvokeMode checks that InvokeMode is one of the Smithy enum
// values (BUFFERED, RESPONSE_STREAM). An empty value is accepted and
// defaults to BUFFERED at the store layer.
func validateInvokeMode(mode string) error {
	if mode == "" {
		return nil
	}
	if mode != "BUFFERED" && mode != "RESPONSE_STREAM" {
		return NewInvalidParameter("InvokeMode",
			fmt.Sprintf("InvokeMode must be BUFFERED or RESPONSE_STREAM; got %q", mode))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Event invoke config validation (Smithy range traits)
// ---------------------------------------------------------------------------

// validateMaximumEventAgeInSeconds enforces the Smithy range
// MaximumEventAgeInSeconds: min 60, max 21600.
func validateMaximumEventAgeInSeconds(v int32) error {
	if v < 60 || v > 21600 {
		return NewInvalidParameter("MaximumEventAgeInSeconds",
			"MaximumEventAgeInSeconds must be between 60 and 21600 seconds")
	}
	return nil
}

// validateMaximumRetryAttempts enforces the Smithy range
// MaximumRetryAttempts: min 0, max 2.
func validateMaximumRetryAttempts(v int32) error {
	if v < 0 || v > 2 {
		return NewInvalidParameter("MaximumRetryAttempts",
			"MaximumRetryAttempts must be between 0 and 2")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Code Signing Config validation
// ---------------------------------------------------------------------------

// validateCodeSigningConfigArn rejects CodeSigningConfigArn because the
// Code Signing Config operations are not implemented on this platform.
// Accepting a non-empty value would create functions that silently bypass
// code signature verification, failing later at invocation time with a
// misleading CodeVerificationFailedException.
func validateCodeSigningConfigArn(arn string) error {
	if arn != "" {
		return NewInvalidParameter("CodeSigningConfigArn",
			"Code Signing Config is not supported")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Layer ARN validation
// ---------------------------------------------------------------------------

func isValidLayerARN(arnStr string) bool {
	if arnStr == "" {
		return false
	}
	if _, service, _, _, _ := arnutil.SplitARN(arnStr); service != "lambda" {
		return false
	}
	resource := arnutil.ExtractResourceFromARN(arnStr)
	if resource == "" {
		return false
	}
	return strings.HasPrefix(resource, "layer:")
}

// ---------------------------------------------------------------------------
// Resource-based policy permission validation
// ---------------------------------------------------------------------------

// validServicePrincipals contains AWS service principal hostnames commonly
// used in Lambda resource-based policies. An unknown ".amazonaws.com"
// suffix is rejected to prevent typos (e.g. "lamda.amazonaws.com") and
// spoofing (e.g. "evil.amazonaws.com").
var validServicePrincipals = map[string]bool{
	"lambda.amazonaws.com":       true,
	"events.amazonaws.com":       true,
	"sns.amazonaws.com":          true,
	"sqs.amazonaws.com":          true,
	"s3.amazonaws.com":           true,
	"kinesis.amazonaws.com":      true,
	"dynamodb.amazonaws.com":     true,
	"logs.amazonaws.com":         true,
	"apigateway.amazonaws.com":   true,
	"cloudwatch.amazonaws.com":   true,
	"config.amazonaws.com":       true,
	"iot.amazonaws.com":          true,
	"ses.amazonaws.com":          true,
	"states.amazonaws.com":       true,
	"firehose.amazonaws.com":     true,
	"codecommit.amazonaws.com":   true,
	"codepipeline.amazonaws.com": true,
	"codebuild.amazonaws.com":    true,
	"ecr.amazonaws.com":          true,
	"ecs.amazonaws.com":          true,
	"eks.amazonaws.com":          true,
	"glue.amazonaws.com":         true,
	"alexa-appkit.amazon.com":    true,
	"scheduler.amazonaws.com":    true,
}

// principalType determines the IAM policy Principal JSON key for a given
// principal string. Returns "AWS" for IAM ARNs, "Service" for recognised
// service principals, and "" for the wildcard "*".
func principalType(principal string) string {
	if principal == "*" {
		return ""
	}
	if strings.HasPrefix(principal, "arn:") {
		_, service, _, _, _ := arnutil.SplitARN(principal)
		if service == "iam" {
			return "AWS"
		}
	}
	return "Service"
}

// validatePermission validates the Principal and Action fields of a
// resource-based policy statement. Principal must be an IAM ARN, a
// recognised service principal, or "*". Action must start with "lambda:".
func validatePermission(p *lambdastore.FunctionPolicy) error {
	if p.Principal == "" {
		return NewInvalidParameter("Principal", "Principal is required")
	}
	if !isValidPrincipal(p.Principal) {
		return NewInvalidParameter("Principal",
			fmt.Sprintf("Principal %q is not a valid IAM ARN, recognised service principal, or wildcard", p.Principal))
	}
	if p.Action == "" {
		return NewInvalidParameter("Action", "Action is required")
	}
	if !strings.HasPrefix(p.Action, "lambda:") {
		return NewInvalidParameter("Action",
			fmt.Sprintf("Action %q must start with 'lambda:'", p.Action))
	}
	return nil
}

// isValidPrincipal checks whether the principal string is a recognised
// format: wildcard "*", an IAM ARN, or a known AWS service principal.
func isValidPrincipal(principal string) bool {
	if principal == "*" {
		return true
	}
	if strings.HasPrefix(principal, "arn:") {
		_, service, _, _, _ := arnutil.SplitARN(principal)
		return service == "iam"
	}
	return validServicePrincipals[principal]
}

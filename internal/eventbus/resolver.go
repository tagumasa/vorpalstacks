package eventbus

import (
	"fmt"
	"strings"
)

var (
	// ErrInvalidARN is returned when an ARN cannot be parsed.
	ErrInvalidARN = fmt.Errorf("eventbus: invalid ARN format")
	// ErrUnsupportedTarget is returned when an ARN refers to a service not
	// supported by the resolver.
	ErrUnsupportedTarget = fmt.Errorf("eventbus: unsupported target service")
)

// TargetAction describes a resolved target for event dispatch, including
// the service type, identifier, region, and account ID.
type TargetAction struct {
	Type       string
	Identifier string
	Region     string
	AccountID  string
}

// TargetResolver defines the contract for resolving an ARN into a
// dispatchable TargetAction.
type TargetResolver interface {
	Resolve(arn string) (*TargetAction, error)
}

// ARNResolver resolves AWS ARNs into TargetActions by parsing the ARN
// components and mapping the service to the appropriate identifier format.
type ARNResolver struct{}

// NewARNResolver creates a new ARNResolver instance.
func NewARNResolver() *ARNResolver {
	return &ARNResolver{}
}

// Resolve parses an ARN string and returns a TargetAction for Lambda, SQS,
// SNS, or Kinesis targets.
func (r *ARNResolver) Resolve(arn string) (*TargetAction, error) {
	if arn == "" {
		return nil, ErrInvalidARN
	}
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return nil, ErrInvalidARN
	}
	arnPrefix := parts[0]
	if arnPrefix != "arn" {
		return nil, ErrInvalidARN
	}
	service := parts[2]
	region := parts[3]
	accountID := parts[4]
	resource := parts[5]

	switch service {
	case "lambda":
		fnName := extractLambdaFunctionName(resource)
		return &TargetAction{
			Type:       "lambda",
			Identifier: fnName,
			Region:     region,
			AccountID:  accountID,
		}, nil
	case "sqs":
		queueURL := buildQueueURL(accountID, region, resource)
		return &TargetAction{
			Type:       "sqs",
			Identifier: queueURL,
			Region:     region,
			AccountID:  accountID,
		}, nil
	case "sns":
		return &TargetAction{
			Type:       "sns",
			Identifier: arn,
			Region:     region,
			AccountID:  accountID,
		}, nil
	case "kinesis":
		streamName := extractStreamName(resource)
		return &TargetAction{
			Type:       "kinesis",
			Identifier: streamName,
			Region:     region,
			AccountID:  accountID,
		}, nil
	default:
		return nil, ErrUnsupportedTarget
	}
}

func extractLambdaFunctionName(resource string) string {
	parts := strings.Split(resource, ":")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return resource
}

func buildQueueURL(accountID, region, resource string) string {
	parts := strings.Split(resource, ":")
	queueName := parts[len(parts)-1]
	return fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, queueName)
}

func extractStreamName(resource string) string {
	parts := strings.Split(resource, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return resource
}

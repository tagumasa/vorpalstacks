package eventbus

import (
	"fmt"
	"strings"
)

var (
	ErrInvalidARN        = fmt.Errorf("eventbus: invalid ARN format")
	ErrUnsupportedTarget = fmt.Errorf("eventbus: unsupported target service")
)

type TargetAction struct {
	Type       string
	Identifier string
	Region     string
	AccountID  string
}

type TargetResolver interface {
	Resolve(arn string) (*TargetAction, error)
}

type ARNResolver struct{}

func NewARNResolver() *ARNResolver {
	return &ARNResolver{}
}

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

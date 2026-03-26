// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import "strings"

// ExtractResourceFromARN extracts the resource portion from an ARN.
func ExtractResourceFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	return resource
}

// ExtractFunctionNameFromARN extracts the Lambda function name from a Lambda function ARN.
func ExtractFunctionNameFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	for i, p := range parts {
		if p == "function" && i+1 < len(parts) {
			return strings.Split(parts[i+1], ":")[0]
		}
	}
	return ""
}

// ExtractRoleNameFromARN extracts the IAM role name from a role ARN.
func ExtractRoleNameFromARN(arn string) string {
	for i := len(arn) - 1; i >= 0; i-- {
		if arn[i] == '/' {
			return arn[i+1:]
		}
	}
	_, _, _, _, resource := SplitARN(arn)
	return strings.TrimPrefix(resource, "role/")
}

// ExtractQueueNameFromARN extracts the SQS queue name from a queue ARN.
func ExtractQueueNameFromARN(arn string) string {
	resource := ExtractResourceFromARN(arn)
	if resource == "" {
		return ""
	}
	parts := strings.Split(resource, ":")
	if len(parts) >= 2 && parts[0] == "" {
		return parts[1]
	}
	return resource
}

// ExtractQueueNameFromURL extracts the queue name from an SQS queue URL.
func ExtractQueueNameFromURL(url string) string {
	if url == "" {
		return ""
	}
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

// ExtractTopicNameFromARN extracts the SNS topic name from a topic ARN.
func ExtractTopicNameFromARN(arn string) string {
	_, service, _, _, resource := SplitARN(arn)
	if service != "sns" || resource == "" {
		return ""
	}
	if idx := strings.LastIndex(resource, ":"); idx >= 0 {
		return resource[idx+1:]
	}
	return resource
}

// ExtractLogGroupNameFromARN extracts the CloudWatch Logs log group name from an ARN.
func ExtractLogGroupNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "log-group:") {
		name := strings.TrimPrefix(resource, "log-group:")
		if idx := strings.Index(name, ":log-stream:"); idx != -1 {
			return name[:idx]
		}
		return name
	}
	return ""
}

// ExtractLogStreamNameFromARN extracts the CloudWatch Logs log stream name from an ARN.
func ExtractLogStreamNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "log-group:") {
		name := strings.TrimPrefix(resource, "log-group:")
		if idx := strings.Index(name, ":log-stream:"); idx != -1 {
			return name[idx+12:]
		}
	}
	return ""
}

// ExtractStreamNameFromARN extracts the Kinesis stream name from a stream ARN.
func ExtractStreamNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "stream/") {
		parts := strings.Split(strings.TrimPrefix(resource, "stream/"), "/")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// ExtractEventBusNameFromARN extracts the EventBridge event bus name from an ARN.
func ExtractEventBusNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "event-bus/") {
		return strings.TrimPrefix(resource, "event-bus/")
	}
	if strings.HasPrefix(resource, "rule/") {
		parts := strings.Split(strings.TrimPrefix(resource, "rule/"), "/")
		if len(parts) >= 2 {
			return parts[0]
		}
	}
	return ""
}

// ExtractRuleNameFromARN extracts the EventBridge rule name from a rule ARN.
func ExtractRuleNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "rule/") {
		parts := strings.Split(strings.TrimPrefix(resource, "rule/"), "/")
		switch len(parts) {
		case 1:
			return parts[0]
		case 2:
			return parts[1]
		}
	}
	return ""
}

// ExtractArchiveNameFromARN extracts the EventBridge archive name from an archive ARN.
func ExtractArchiveNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "archive/") {
		return strings.TrimPrefix(resource, "archive/")
	}
	return ""
}

// ExtractStateMachineNameFromARN extracts the Step Functions state machine name from an ARN.
func ExtractStateMachineNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.HasPrefix(resource, "stateMachine:") {
		return strings.TrimPrefix(resource, "stateMachine:")
	}
	if strings.HasPrefix(resource, "execution:") {
		parts := strings.Split(strings.TrimPrefix(resource, "execution:"), ":")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// IsLambdaARN returns true if the ARN is for a Lambda function.
func IsLambdaARN(arn string) bool { return strings.Contains(arn, ":lambda:") }

// IsKinesisARN returns true if the ARN is for a Kinesis stream.
func IsKinesisARN(arn string) bool { return strings.Contains(arn, ":kinesis:") }

// IsSQSARN returns true if the ARN is for an SQS queue.
func IsSQSARN(arn string) bool { return strings.Contains(arn, ":sqs:") }

// IsSNSARN returns true if the ARN is for an SNS topic.
func IsSNSARN(arn string) bool { return strings.Contains(arn, ":sns:") }

// IsEventBridgeARN returns true if the ARN is for an EventBridge resource.
func IsEventBridgeARN(arn string) bool { return strings.Contains(arn, ":events:") }

// IsStateMachineARN returns true if the ARN is for a Step Functions state machine.
func IsStateMachineARN(arn string) bool { return strings.Contains(arn, ":states:") }

// ParseLogGroupARN extracts the log group name from a CloudWatch Logs log group ARN.
// This is an alias for ExtractLogGroupNameFromARN for backward compatibility.
func ParseLogGroupARN(arn string) string {
	return ExtractLogGroupNameFromARN(arn)
}

// ParseTableARN extracts the table name from a DynamoDB table ARN.
func ParseTableARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if len(resource) > 6 && resource[:6] == "table/" {
		rest := resource[6:]
		if idx := strings.Index(rest, "/stream/"); idx != -1 {
			return rest[:idx]
		}
		return rest
	}
	return ""
}

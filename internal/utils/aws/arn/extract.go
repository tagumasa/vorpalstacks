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

// ExtractPlatformApplicationFromARN extracts the platform and application
// name from an SNS platform application ARN (resource form
// app/<platform>/<name>). Either return is empty when the ARN is not an
// SNS platform application ARN of exactly that shape.
func ExtractPlatformApplicationFromARN(arn string) (platform, name string) {
	_, service, _, _, resource := SplitARN(arn)
	if service != "sns" || !strings.HasPrefix(resource, "app/") {
		return "", ""
	}
	parts := strings.Split(strings.TrimPrefix(resource, "app/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", ""
	}
	return parts[0], parts[1]
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

// ExtractScheduledQueryIdFromARN extracts the scheduled query id from a
// CloudWatch Logs scheduled-query ARN.
func ExtractScheduledQueryIdFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if id, ok := strings.CutPrefix(resource, "scheduled-query:"); ok {
		return id
	}
	return ""
}

// ExtractLookupTableNameFromARN extracts the lookup table name from a
// CloudWatch Logs lookup-table ARN.
func ExtractLookupTableNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if name, ok := strings.CutPrefix(resource, "lookup-table:"); ok {
		return name
	}
	return ""
}

// ExtractDeliverySourceNameFromARN extracts the delivery source name from
// a CloudWatch Logs delivery-source ARN.
func ExtractDeliverySourceNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if name, ok := strings.CutPrefix(resource, "delivery-source:"); ok {
		return name
	}
	return ""
}

// ExtractDeliveryDestinationNameFromARN extracts the delivery destination
// name from a CloudWatch Logs delivery-destination ARN.
func ExtractDeliveryDestinationNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if name, ok := strings.CutPrefix(resource, "delivery-destination:"); ok {
		return name
	}
	return ""
}

// ExtractDeliveryIdFromARN extracts the delivery id from a CloudWatch
// Logs delivery ARN.
func ExtractDeliveryIdFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if id, ok := strings.CutPrefix(resource, "delivery:"); ok {
		return id
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

// ExtractConnectionNameFromARN extracts the EventBridge connection name from
// a connection ARN (resource form connection/<name>/<id>; the trailing id
// segment is the ARN's unique suffix, not part of the name).
func ExtractConnectionNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if !strings.HasPrefix(resource, "connection/") {
		return ""
	}
	parts := strings.Split(strings.TrimPrefix(resource, "connection/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	return parts[0]
}

// ExtractApiDestinationNameFromARN extracts the EventBridge API destination
// name from an API destination ARN (resource form
// api-destination/<name>/<id>).
func ExtractApiDestinationNameFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if !strings.HasPrefix(resource, "api-destination/") {
		return ""
	}
	parts := strings.Split(strings.TrimPrefix(resource, "api-destination/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	return parts[0]
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
func IsLambdaARN(arn string) bool { return GetServiceFromARN(arn) == "lambda" }

// IsKinesisARN returns true if the ARN is for a Kinesis stream.
func IsKinesisARN(arn string) bool { return GetServiceFromARN(arn) == "kinesis" }

// IsSQSARN returns true if the ARN is for an SQS queue.
func IsSQSARN(arn string) bool { return GetServiceFromARN(arn) == "sqs" }

// IsSNSARN returns true if the ARN is for an SNS topic.
func IsSNSARN(arn string) bool { return GetServiceFromARN(arn) == "sns" }

// IsEventBridgeARN returns true if the ARN is for an EventBridge resource.
func IsEventBridgeARN(arn string) bool { return GetServiceFromARN(arn) == "events" }

// IsStateMachineARN returns true if the ARN is for a Step Functions state machine.
func IsStateMachineARN(arn string) bool { return GetServiceFromARN(arn) == "states" }

// ParseTableARN extracts the table name from a DynamoDB table ARN.
// Also accepts stream ARNs (table/<name>/stream/<label>) and returns the
// table name portion.
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

// ParseStreamARN extracts the table name from a DynamoDB stream ARN.
// Returns "" for non-stream ARNs.
func ParseStreamARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if strings.Contains(resource, "/stream/") {
		return ParseTableARN(arn)
	}
	return ""
}

// ExtractAPIGatewayFunctionRef parses an API Gateway Lambda invocation URI
// (arn:partition:apigateway:region:lambda:path/2015-03-31/functions/{ref}/invocations)
// and returns the function reference between the functions path segment and
// the invocations suffix — a function ARN, partial ARN, or function name.
// ok is false when the URI is not an invocation URI of that grammar.
func ExtractAPIGatewayFunctionRef(uri string) (ref string, ok bool) {
	parsed, err := ParseARN(uri)
	if err != nil || parsed.Service != "apigateway" {
		return "", false
	}
	if parsed.AccountID != "lambda" || !strings.HasPrefix(parsed.Resource, "path/2015-03-31/functions/") {
		return "", false
	}
	rest := strings.TrimPrefix(parsed.Resource, "path/2015-03-31/functions/")
	if !strings.HasSuffix(rest, "/invocations") {
		return "", false
	}
	ref = strings.TrimSuffix(rest, "/invocations")
	if ref == "" {
		return "", false
	}
	return ref, true
}

// ExtractBackupIdFromARN extracts the identity segment from a DynamoDB
// backup ARN (table/<table>/backup/<id>) — the generated id current
// records are keyed by.
func ExtractBackupIdFromARN(arn string) string {
	_, _, _, _, resource := SplitARN(arn)
	if idx := strings.Index(resource, "/backup/"); idx != -1 {
		return resource[idx+len("/backup/"):]
	}
	return ""
}

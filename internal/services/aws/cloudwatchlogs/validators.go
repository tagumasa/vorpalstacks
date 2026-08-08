package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strings"

	"vorpalstacks/internal/utils/aws/arn"
)

// --- Pattern constants (Smithy-derived) ---

var (
	// logGroupNamePattern matches the Smithy trait on LogGroupName:
	// ^[\.\-_/#A-Za-z0-9]+$  length 1-512
	logGroupNamePattern = regexp.MustCompile(`^[\.\-_/#A-Za-z0-9]+$`)

	// noColonAsteriskPattern matches FilterName, DestinationName, LogStreamName
	// whose Smithy trait is ^[^:*]*$ with a length constraint of 1-512.
	noColonAsteriskPattern = regexp.MustCompile(`^[^:*]*$`)

	// metricNamePattern matches MetricName and MetricNamespace:
	// ^[^:*$]*$  length 0-255
	metricNamePattern = regexp.MustCompile(`^[^:*$]*$`)

	// kmsKeyArnPattern validates a KMS key ARN.
	kmsKeyArnPattern   = regexp.MustCompile(`^arn:aws[a-z\-]*:kms:[a-z0-9-]+:\d{12}:key/[a-f0-9\-]+$`)
	kmsAliasArnPattern = regexp.MustCompile(`^arn:aws[a-z\-]*:kms:[a-z0-9-]+:\d{12}:alias/.+$`)
)

// --- Enum maps ---

var (
	validLogGroupClasses = map[string]bool{
		"STANDARD":          true,
		"INFREQUENT_ACCESS": true,
		"DELIVERY":          true,
	}

	validPolicyTypes = map[string]bool{
		"DATA_PROTECTION_POLICY":     true,
		"SUBSCRIPTION_FILTER_POLICY": true,
		"FIELD_INDEX_POLICY":         true,
		"TRANSFORMER_POLICY":         true,
		"METRIC_EXTRACTION_POLICY":   true,
	}

	validDistributions = map[string]bool{
		"Random":      true,
		"ByLogStream": true,
	}

	validScheduledQueryStates = map[string]bool{
		"ENABLED":  true,
		"DISABLED": true,
	}
)

// --- Name / pattern validators ---

// validateLogGroupName validates a log group name against the Smithy pattern
// ^[\.\-_/#A-Za-z0-9]+$ and length constraint 1-512.
func validateLogGroupName(name string) error {
	if name == "" {
		return ErrMissingParameter
	}
	if len(name) > 512 || !logGroupNamePattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid log group name: %s. Must match ^[.\\-_/#A-Za-z0-9]+$ and be 1-512 characters", name), 400)
	}
	return nil
}

// validateLogStreamName validates a log stream name against the Smithy pattern
// ^[^:*]*$ and length constraint 1-512.
func validateLogStreamName(name string) error {
	if name == "" {
		return ErrMissingParameter
	}
	if len(name) > 512 || !noColonAsteriskPattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid log stream name: %s. Must not contain ':' or '*' and be 1-512 characters", name), 400)
	}
	return nil
}

// validateFilterName validates a metric or subscription filter name.
// Smithy: ^[^:*]*$  length 1-512.
func validateFilterName(name string) error {
	if name == "" {
		return ErrMissingParameter
	}
	if len(name) > 512 || !noColonAsteriskPattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid filter name: %s. Must not contain ':' or '*' and be 1-512 characters", name), 400)
	}
	return nil
}

// validateDestinationName validates a destination name.
// Smithy: ^[^:*]*$  length 1-512.
func validateDestinationName(name string) error {
	if name == "" {
		return ErrMissingParameter
	}
	if len(name) > 512 || !noColonAsteriskPattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid destination name: %s. Must not contain ':' or '*' and be 1-512 characters", name), 400)
	}
	return nil
}

// validateFilterPattern validates a filter pattern string.
// Smithy: length 0-1024.
func validateFilterPattern(pattern string) error {
	if len(pattern) > 1024 {
		return NewLogsError("InvalidParameterException",
			"Filter pattern must be between 0 and 1024 characters", 400)
	}
	return nil
}

// validatePolicyDocument validates a policy document string.
// Smithy: length 1-51200.
func validatePolicyDocument(doc string) error {
	if doc == "" {
		return NewLogsError("InvalidParameterException",
			"Policy document is required", 400)
	}
	if len(doc) > 51200 {
		return NewLogsError("InvalidParameterException",
			"Policy document must not exceed 51200 characters", 400)
	}
	return nil
}

// validateAccessPolicy validates a destination access policy string.
// Smithy: AccessPolicy shape — @required, @length min:1.
func validateAccessPolicy(policy string) error {
	if policy == "" {
		return NewLogsError("InvalidParameterException",
			"Access policy is required", 400)
	}
	return nil
}

// validateQueryDefinitionName validates a query definition name.
// Smithy: length 1-255.
func validateQueryDefinitionName(name string) error {
	if name == "" {
		return ErrMissingParameter
	}
	if len(name) > 255 {
		return NewLogsError("InvalidParameterException",
			"Query definition name must be between 1 and 255 characters", 400)
	}
	return nil
}

// validateQueryString validates a query string.
// Smithy: length 0-10000.
func validateQueryString(qs string) error {
	if len(qs) > 10000 {
		return NewLogsError("InvalidParameterException",
			"Query string must not exceed 10000 characters", 400)
	}
	return nil
}

// validateScheduledQueryName validates a scheduled query name.
// Smithy: length 1-300.
func validateScheduledQueryName(name string) error {
	if name == "" {
		return ErrMissingParameter
	}
	if len(name) > 300 {
		return NewLogsError("InvalidParameterException",
			"Scheduled query name must be between 1 and 300 characters", 400)
	}
	return nil
}

// validateExportDestinationBucket validates an S3 bucket name for export tasks.
// Smithy: length 1-512.
func validateExportDestinationBucket(bucket string) error {
	if bucket == "" {
		return ErrMissingParameter
	}
	if len(bucket) > 512 {
		return NewLogsError("InvalidParameterException",
			"Export destination bucket must be between 1 and 512 characters", 400)
	}
	return nil
}

// validateMetricName validates a metric transformation metric name.
// Smithy: ^[^:*$]*$  length 0-255.
func validateMetricName(name string) error {
	if len(name) > 255 || !metricNamePattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid metric name: %s. Must not contain ':', '*', or '$' and be 0-255 characters", name), 400)
	}
	return nil
}

// validateMetricNamespace validates a metric transformation namespace.
// Smithy: ^[^:*$]*$  length 0-255.
func validateMetricNamespace(ns string) error {
	if len(ns) > 255 || !metricNamePattern.MatchString(ns) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid metric namespace: %s. Must not contain ':', '*', or '$' and be 0-255 characters", ns), 400)
	}
	return nil
}

// validateMetricValue validates a metric transformation value expression.
// Smithy: length 0-100.
func validateMetricValue(val string) error {
	if len(val) > 100 {
		return NewLogsError("InvalidParameterException",
			"Metric value must not exceed 100 characters", 400)
	}
	return nil
}

// --- KMS key validators ---

// validateKmsKeyId validates a KMS key ID or ARN.
func validateKmsKeyId(kmsKeyId string) error {
	if kmsKeyId == "" {
		return nil
	}
	if kmsKeyArnPattern.MatchString(kmsKeyId) || kmsAliasArnPattern.MatchString(kmsKeyId) {
		return nil
	}
	parsed, err := arn.ParseARN(kmsKeyId)
	if err != nil || parsed.Service != "kms" {
		return NewLogsError("InvalidParameterException",
			"kmsKeyId must be a valid KMS key ARN", 400)
	}
	if !strings.HasPrefix(parsed.Resource, "key/") && !strings.HasPrefix(parsed.Resource, "alias/") {
		return NewLogsError("InvalidParameterException",
			"kmsKeyId resource must be a key UUID (key/...) or alias (alias/...)", 400)
	}
	return nil
}

// --- Enum validators ---

func validateLogGroupClass(class string) bool {
	return validLogGroupClasses[class]
}

func validatePolicyType(t string) bool {
	return validPolicyTypes[t]
}

func validateDistribution(d string) bool {
	return validDistributions[d]
}

func validateScheduledQueryState(state string) bool {
	return validScheduledQueryStates[state]
}

// --- Range / limit validators ---

// validateListLimit applies a default when limit is unset (<= 0) and enforces
// an upper bound per the Smithy range trait. Returns the resolved limit or an
// InvalidParameterException when the caller-supplied value exceeds max.
func validateListLimit(limit, defaultVal, maxVal int32) (int32, error) {
	if limit <= 0 {
		return defaultVal, nil
	}
	if limit > maxVal {
		return 0, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Limit must be between 1 and %d", maxVal), 400)
	}
	return limit, nil
}

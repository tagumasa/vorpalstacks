package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
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

// --- Policy document JSON validators ---

// validatePolicyDocumentJSON validates that a policy document is non-empty,
// within the Smithy length limit, and parseable as valid JSON.
func validatePolicyDocumentJSON(doc string) error {
	if err := validatePolicyDocument(doc); err != nil {
		return err
	}
	if !json.Valid([]byte(doc)) {
		return NewLogsError("InvalidParameterException",
			"Policy document must be valid JSON", 400)
	}
	return nil
}

// validateAccessPolicyJSON validates that an access policy is non-empty and
// parseable as valid JSON.
func validateAccessPolicyJSON(policy string) error {
	if err := validateAccessPolicy(policy); err != nil {
		return err
	}
	if !json.Valid([]byte(policy)) {
		return NewLogsError("InvalidParameterException",
			"Access policy must be valid JSON", 400)
	}
	return nil
}

// --- Reserved name prefix validators ---

// validatePolicyNamePrefix rejects policy names that start with the reserved
// AWS prefix. AWS documentation states that policy names must not begin with
// "aws/" or "AWS:".
func validatePolicyNamePrefix(name string) error {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "aws/") || strings.HasPrefix(lower, "aws:") {
		return NewLogsError("InvalidParameterException",
			"Policy names starting with 'aws/' are reserved and not allowed", 400)
	}
	return nil
}

// --- Selection criteria length validator (25 KB max) ---

const maxSelectionCriteriaBytes = 25 * 1024

// validateSelectionCriteria checks that the selection criteria string does
// not exceed the AWS-documented maximum of 25 KB.
func validateSelectionCriteria(sc string) error {
	if len(sc) > maxSelectionCriteriaBytes {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("SelectionCriteria must not exceed %d bytes", maxSelectionCriteriaBytes), 400)
	}
	return nil
}

// --- Query language enum validator ---

var validQueryLanguages = map[string]bool{
	"CWLI": true,
	"SQL":  true,
	"PPL":  true,
}

func validateQueryLanguage(ql string) bool {
	if ql == "" {
		return true
	}
	return validQueryLanguages[ql]
}

// --- IAM role ARN validator ---

// validateIAMRoleArn parses an ARN and confirms it belongs to the IAM service
// with a role/ resource prefix.
func validateIAMRoleArn(roleArn string) error {
	if roleArn == "" {
		return nil
	}
	parsed, err := arn.ParseARN(roleArn)
	if err != nil {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid role ARN: %s", roleArn), 400)
	}
	if parsed.Service != "iam" {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Role ARN must be an IAM ARN, got service: %s", parsed.Service), 400)
	}
	if !strings.HasPrefix(parsed.Resource, "role/") {
		return NewLogsError("InvalidParameterException",
			"Role ARN resource must start with 'role/'", 400)
	}
	return nil
}

// --- Kinesis / Firehose target ARN validator ---

// validateKinesisOrFirehoseArn confirms that the targetArn refers to a
// Kinesis Data Stream or Kinesis Data Firehose delivery stream.
func validateKinesisOrFirehoseArn(targetArn string) error {
	if targetArn == "" {
		return nil
	}
	parsed, err := arn.ParseARN(targetArn)
	if err != nil {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid target ARN: %s", targetArn), 400)
	}
	if parsed.Service != "kinesis" && parsed.Service != "firehose" {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Target ARN must be a Kinesis or Firehose ARN, got service: %s", parsed.Service), 400)
	}
	return nil
}

// --- Field selection criteria length validator (2000 char max) ---

const maxFieldSelectionCriteriaLen = 2000

// validateFieldSelectionCriteria enforces the Smithy @length constraint of
// 0-2000 characters on FieldSelectionCriteria.
func validateFieldSelectionCriteria(s string) error {
	if len(s) > maxFieldSelectionCriteriaLen {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("FieldSelectionCriteria must not exceed %d characters", maxFieldSelectionCriteriaLen), 400)
	}
	return nil
}

// --- Destination prefix length validator (1024 char max) ---

const maxDestinationPrefixLen = 1024

// validateDestinationPrefix enforces the documented maximum of 1024
// characters on the export destination prefix.
func validateDestinationPrefix(p string) error {
	if len(p) > maxDestinationPrefixLen {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("DestinationPrefix must not exceed %d characters", maxDestinationPrefixLen), 400)
	}
	return nil
}

// --- Log group identifier count validator (50 max) ---

const maxLogGroupIdentifiers = 50

// validateLogGroupIdentifierCount enforces the AWS-documented limit of 1-50
// log group identifiers per scheduled query.
func validateLogGroupIdentifierCount(ids []string) error {
	if len(ids) > maxLogGroupIdentifiers {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Number of log group identifiers must not exceed %d", maxLogGroupIdentifiers), 400)
	}
	return nil
}

// --- StartFromHead date constraint ---

// jan012024Millis is Jan 1, 2024 00:00:00 UTC in epoch milliseconds.
// AWS requires startTime on or after this date when startFromHead=false.
const jan012024Millis = 1704067200000

// validateStartFromHeadDate enforces the AWS requirement that when
// startFromHead is false, startTime must be on or after Jan 1, 2024.
func validateStartFromHeadDate(startFromHead bool, startTime int64) error {
	if !startFromHead && startTime > 0 && startTime < jan012024Millis {
		return NewLogsError("InvalidParameterException",
			"Setting startFromHead to false is supported only when startTime is on or after Jan 1, 2024 00:00:00 UTC", 400)
	}
	return nil
}

// --- Log event validation (PutLogEvents) ---

const (
	// maxEventsTimeSpan is the maximum allowed time span (in milliseconds)
	// for a single PutLogEvents batch. AWS rejects the entire batch if the
	// span between the earliest and latest event exceeds 24 hours.
	maxEventsTimeSpan int64 = 24 * 60 * 60 * 1000

	// tooNewThreshold is the maximum future offset (in milliseconds) for
	// an event timestamp. Events more than 2 hours in the future are
	// rejected individually.
	tooNewThreshold int64 = 2 * 60 * 60 * 1000

	// tooOldThreshold is the maximum age (in milliseconds) for an event
	// timestamp. Events older than 14 days are rejected individually.
	tooOldThreshold int64 = 14 * 24 * 60 * 60 * 1000
)

// validateLogEvents checks that log events satisfy the PutLogEvents
// constraints required by AWS CloudWatch Logs:
//   - Events must be in chronological order (by timestamp).
//   - No event may be more than 2 hours in the future.
//   - No event may be older than 14 days.
//   - The timespan between the earliest and latest valid event must not
//     exceed 24 hours.
//
// Events that fall outside the age thresholds are silently excluded from
// the returned valid slice; the caller receives information about the
// rejected indices via the returned map, which is suitable for inclusion
// in the response as rejectedLogEventsInfo.
func validateLogEvents(events []logsstore.LogEntry) ([]logsstore.LogEntry, map[string]interface{}, error) {
	now := time.Now().UnixMilli()

	// Chronological order check must be performed on ALL events in the
	// batch, not just the age-valid subset. AWS rejects the entire batch
	// if any event is out of order, regardless of whether some events are
	// later individually rejected for being too old or too new.
	for i := 1; i < len(events); i++ {
		if events[i].Timestamp < events[i-1].Timestamp {
			return nil, nil, awserrors.NewAWSError("InvalidParameterException",
				"log events in the batch must be in chronological order", 400)
		}
	}

	var valid []logsstore.LogEntry
	var tooOldEndIndex int
	tooNewStartIndex := -1

	for i, e := range events {
		if e.Timestamp > now+tooNewThreshold {
			if tooNewStartIndex == -1 || i < tooNewStartIndex {
				tooNewStartIndex = i
			}
			continue
		}
		if e.Timestamp < now-tooOldThreshold {
			tooOldEndIndex = i + 1
			continue
		}
		valid = append(valid, e)
	}

	if len(valid) == 0 {
		rejected := buildRejectedInfo(tooOldEndIndex, tooNewStartIndex, len(events))
		return nil, rejected, nil
	}

	span := valid[len(valid)-1].Timestamp - valid[0].Timestamp
	if span > maxEventsTimeSpan {
		return nil, nil, awserrors.NewAWSError("InvalidParameterException",
			"Events span must not exceed 24 hours", 400)
	}

	rejected := buildRejectedInfo(tooOldEndIndex, tooNewStartIndex, len(events))
	return valid, rejected, nil
}

// buildRejectedInfo constructs the rejectedLogEventsInfo map from the
// computed too-old and too-new indices. If no events were rejected an
// empty map is returned.
func buildRejectedInfo(tooOldEndIndex, tooNewStartIndex, totalEvents int) map[string]interface{} {
	if tooOldEndIndex == 0 && tooNewStartIndex == -1 {
		return nil
	}
	info := make(map[string]interface{})
	if tooOldEndIndex > 0 {
		info["tooOldLogEventEndIndex"] = tooOldEndIndex
	}
	if tooNewStartIndex >= 0 {
		info["tooNewLogEventStartIndex"] = tooNewStartIndex
	}
	return info
}

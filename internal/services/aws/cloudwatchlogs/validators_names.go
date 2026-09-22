package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/filterpattern"
)

// The name, pattern and member-shape validators: one member's
// length, alphabet or ARN form per function, the Smithy traits cited in
// place.

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

// --- Name / pattern validators ---

// validateLogGroupName validates a log group name against the Smithy pattern
// ^[\.\-_/#A-Za-z0-9]+$ and length constraint 1-512.
func validateLogGroupName(name string) error {
	if name == "" {
		return errRequiredMember("logGroupName")
	}
	if len(name) > logsstore.MaxLogGroupNameLength || !logGroupNamePattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid log group name: %s. Must match ^[.\\-_/#A-Za-z0-9]+$ and be 1-%d characters", name, logsstore.MaxLogGroupNameLength), 400)
	}
	return nil
}

// validateLogStreamName validates a log stream name against the Smithy pattern
// ^[^:*]*$ and length constraint 1-512 counted in Unicode characters.
func validateLogStreamName(name string) error {
	if name == "" {
		return errRequiredMember("logStreamName")
	}
	if utf8.RuneCountInString(name) > logsstore.MaxLogStreamNameLength || !noColonAsteriskPattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid log stream name: %s. Must not contain ':' or '*' and be 1-%d characters", name, logsstore.MaxLogStreamNameLength), 400)
	}
	return nil
}

// validateLogStreamNamePrefix validates a logStreamNamePrefix member, which
// targets LogStreamName and so carries its pattern (^[^:*]*$) and length
// ceiling; the member is optional, so the empty prefix is absent.
func validateLogStreamNamePrefix(prefix string) error {
	if prefix == "" {
		return nil
	}
	if utf8.RuneCountInString(prefix) > logsstore.MaxLogStreamNameLength || !noColonAsteriskPattern.MatchString(prefix) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid logStreamNamePrefix: %s. Must not contain ':' or '*' and be at most %d characters", prefix, logsstore.MaxLogStreamNameLength), 400)
	}
	return nil
}

// validateFilterName validates a metric or subscription filter name.
// Smithy: ^[^:*]*$  length 1-512 counted in Unicode characters.
func validateFilterName(name string) error {
	if name == "" {
		return errRequiredMember("filterName")
	}
	if utf8.RuneCountInString(name) > 512 || !noColonAsteriskPattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid filter name: %s. Must not contain ':' or '*' and be 1-%d characters", name, logsstore.MaxFilterNameLength), 400)
	}
	return nil
}

// validateDestinationName validates a destination name.
// Smithy: ^[^:*]*$  length 1-512 counted in Unicode characters.
func validateDestinationName(name string) error {
	if name == "" {
		return errRequiredMember("destinationName")
	}
	if utf8.RuneCountInString(name) > 512 || !noColonAsteriskPattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid destination name: %s. Must not contain ':' or '*' and be 1-%d characters", name, logsstore.MaxDestinationNameLength), 400)
	}
	return nil
}

// validateFilterPattern validates a filter pattern string.
// Smithy: length 0-1024 counted in Unicode characters. The syntax pass
// rejects the documented invalid forms — unbalanced quotes,
// unterminated delimited/JSON shapes, regex outside the documented
// dialect, more than two regex in a delimited or JSON pattern — with
// InvalidParameterException instead of storing them as match-nothing
// filters.
func validateFilterPattern(pattern string) error {
	if utf8.RuneCountInString(pattern) > logsstore.MaxFilterPatternLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Filter pattern must be between 0 and %d characters", logsstore.MaxFilterPatternLength), 400)
	}
	if err := filterpattern.ValidatePattern(pattern); err != nil {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid filter pattern: %s", pattern), 400)
	}
	return nil
}

// validatePolicyDocument validates a policy document string.
// Smithy: length 1-51200 counted in Unicode characters.
func validatePolicyDocument(doc string) error {
	if doc == "" {
		return NewLogsError("InvalidParameterException",
			"Policy document is required", 400)
	}
	if utf8.RuneCountInString(doc) > logsstore.MaxPolicyDocumentLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Policy document must not exceed %d characters", logsstore.MaxPolicyDocumentLength), 400)
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
// Smithy: length 1-255 counted in Unicode characters.
func validateQueryDefinitionName(name string) error {
	if name == "" {
		return errRequiredMember("name")
	}
	if utf8.RuneCountInString(name) > logsstore.MaxQueryDefinitionNameLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Query definition name must be between 1 and %d characters", logsstore.MaxQueryDefinitionNameLength), 400)
	}
	return nil
}

// validateQueryString validates a query string's length domain against
// the caller's documented floor: StartQuery's queryString carries
// "Minimum length of 0" (an empty present member runs the zero-command
// return-all) while PutQueryDefinition's queryDefinitionString and the
// scheduled-query members carry a minimum of 1; every caller caps at
// 10000.
func validateQueryString(qs string, minRunes int) error {
	if n := utf8.RuneCountInString(qs); n < minRunes || n > logsstore.MaxQueryStringLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Query string must be between %d and %d characters", minRunes, logsstore.MaxQueryStringLength), 400)
	}
	return nil
}

// queryParameterNamePattern is the documented parameter-name constraint:
// "must start with a letter or underscore, and contain only letters,
// digits, and underscores".
var queryParameterNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// validateQueryParameters validates a saved query's parameter list: at
// most twenty parameters, each with a constrained name (required, letters,
// digits and underscores, starting with a letter or underscore) and
// bounded default value and description.
func validateQueryParameters(params []logsstore.QueryParameter) error {
	if len(params) > logsstore.MaxQueryParameters {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("A query definition accepts at most %d parameters", logsstore.MaxQueryParameters), 400)
	}
	for _, p := range params {
		if p.Name == "" {
			return errRequiredMember("parameters.name")
		}
		if len(p.Name) > logsstore.MaxQueryParameterNameLength || !queryParameterNamePattern.MatchString(p.Name) {
			return NewLogsError("InvalidParameterException",
				"Parameter names must start with a letter or underscore and contain only letters, digits, and underscores", 400)
		}
		if len(p.DefaultValue) > logsstore.MaxQueryParameterDefaultValueLen {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Parameter %s default value must not exceed %d characters", p.Name, logsstore.MaxQueryParameterDefaultValueLen), 400)
		}
		if len(p.Description) > logsstore.MaxQueryParameterDescriptionLen {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Parameter %s description must not exceed %d characters", p.Name, logsstore.MaxQueryParameterDescriptionLen), 400)
		}
	}
	return nil
}

// validateScheduledQueryName validates a scheduled query name.
// Smithy: length 1-300 counted in Unicode characters.
func validateScheduledQueryName(name string) error {
	if name == "" {
		return errRequiredMember("name")
	}
	if utf8.RuneCountInString(name) > logsstore.MaxScheduledQueryNameLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Scheduled query name must be between 1 and %d characters", logsstore.MaxScheduledQueryNameLength), 400)
	}
	return nil
}

// validateExportDestinationBucket validates an S3 bucket name for export tasks.
// Smithy: length 1-512 counted in Unicode characters.
func validateExportDestinationBucket(bucket string) error {
	if bucket == "" {
		return errRequiredMember("destination")
	}
	if utf8.RuneCountInString(bucket) > logsstore.MaxExportDestinationLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Export destination bucket must be between 1 and %d characters", logsstore.MaxExportDestinationLength), 400)
	}
	return nil
}

// validateMetricName validates a metric transformation metric name.
// Smithy: ^[^:*$]*$  length 0-255 counted in Unicode characters.
func validateMetricName(name string) error {
	if utf8.RuneCountInString(name) > logsstore.MaxMetricNameLength || !metricNamePattern.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid metric name: %s. Must not contain ':', '*', or '$' and be 0-%d characters", name, logsstore.MaxMetricNameLength), 400)
	}
	return nil
}

// validateMetricNamespace validates a metric transformation namespace.
// Smithy: ^[^:*$]*$  length 0-255 counted in Unicode characters.
func validateMetricNamespace(ns string) error {
	if utf8.RuneCountInString(ns) > logsstore.MaxMetricNameLength || !metricNamePattern.MatchString(ns) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid metric namespace: %s. Must not contain ':', '*', or '$' and be 0-%d characters", ns, logsstore.MaxMetricNameLength), 400)
	}
	return nil
}

// validateMetricValue validates a metric transformation value expression.
// Smithy: length 0-100 counted in Unicode characters.
func validateMetricValue(val string) error {
	if utf8.RuneCountInString(val) > logsstore.MaxMetricValueLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Metric value must not exceed %d characters", logsstore.MaxMetricValueLength), 400)
	}
	return nil
}

// --- KMS key validators ---

// validateKmsKeyId validates a KMS key ID or ARN. Every kmsKeyId member
// targets the one KmsKeyId shape, so the shape's length ceiling is
// enforced here once — an unbounded alias pattern makes a form-valid
// over-long value reachable on any surface that checks the form alone.
func validateKmsKeyId(kmsKeyId string) error {
	if kmsKeyId == "" {
		return nil
	}
	if utf8.RuneCountInString(kmsKeyId) > logsstore.MaxKmsKeyIdLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("kmsKeyId exceeds %d characters", logsstore.MaxKmsKeyIdLength), 400)
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

// --- Import identifier validator ---

// importIdPattern is the Smithy pattern trait on ImportId.
var importIdPattern = regexp.MustCompile(`^[\-a-zA-Z0-9]+$`)

// validateImportId enforces the ImportId traits — 1-256 characters
// matching [a-zA-Z0-9-] — on the operations that address an import: every
// import operation's declared error list carries InvalidParameterException,
// so a malformed id is a parameter error, not a missing resource.
func validateImportId(importId string) error {
	if importId == "" {
		return errRequiredMember("importId")
	}
	if len(importId) > logsstore.MaxImportIdLength || !importIdPattern.MatchString(importId) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid importId: %s. Must be 1-%d characters matching [a-zA-Z0-9-]", importId, logsstore.MaxImportIdLength), 400)
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

// validateIAMRoleArn confirms the ARN is a valid IAM role ARN through
// the shared platform validator — which also requires the twelve-digit
// account every IAM role ARN carries and a non-empty role name — and
// maps a rejection onto the member's InvalidParameterException.
func validateIAMRoleArn(roleArn string) error {
	if roleArn == "" {
		return nil
	}
	if !arn.IsValidRoleARN(roleArn) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid role ARN: %s", roleArn), 400)
	}
	return nil
}

// --- Kinesis target ARN validator ---

// validateKinesisStreamArn confirms that the targetArn refers to a Kinesis
// Data Stream. The PutDestinationRequest targetArn member documentation
// names the Kinesis stream alone as the destination target; Firehose
// delivery streams would be rejected at delivery anyway, since the platform
// Firehose service does not exist yet.
func validateKinesisStreamArn(targetArn string) error {
	if targetArn == "" {
		return nil
	}
	parsed, err := arn.ParseARN(targetArn)
	if err != nil {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid target ARN: %s", targetArn), 400)
	}
	if parsed.Service != "kinesis" {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Target ARN must be a Kinesis stream ARN, got service: %s", parsed.Service), 400)
	}
	return nil
}

// --- Field selection criteria length validator (2000 char max) ---

const maxFieldSelectionCriteriaLen = 2000

// validateFieldSelectionCriteria enforces the Smithy @length constraint of
// 0-2000 characters (counted in Unicode characters) on FieldSelectionCriteria.
func validateFieldSelectionCriteria(s string) error {
	if utf8.RuneCountInString(s) > maxFieldSelectionCriteriaLen {
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

// validStandardUnits is the StandardUnit enum the transformation's unit
// member targets ("Valid Values: Seconds | Microseconds | ... | None").
var validStandardUnits = map[string]bool{
	"Seconds": true, "Microseconds": true, "Milliseconds": true,
	"Bytes": true, "Kilobytes": true, "Megabytes": true, "Gigabytes": true, "Terabytes": true,
	"Bits": true, "Kilobits": true, "Megabits": true, "Gigabits": true, "Terabits": true,
	"Percent": true, "Count": true,
	"Bytes/Second": true, "Kilobytes/Second": true, "Megabytes/Second": true,
	"Gigabytes/Second": true, "Terabytes/Second": true,
	"Bits/Second": true, "Kilobits/Second": true, "Megabits/Second": true,
	"Gigabits/Second": true, "Terabits/Second": true,
	"Count/Second": true, "None": true,
}

// validateTransformDimensions enforces the dimensions member's traits:
// "One metric filter can include as many as three dimensions" (the
// emitSystemFieldDimensions members count toward the same total) with
// keys and values of at most 255 characters, and the documented conflict
// rule "If you assign dimensions to a metric created by a metric filter,
// you can't assign a default value for that metric".
func validateTransformDimensions(t logsstore.MetricTransformation, emitSystemFieldDimensions []string) error {
	total := len(t.Dimensions) + len(emitSystemFieldDimensions)
	if total > logsstore.MaxMetricFilterDimensions {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("One metric filter can include as many as %d dimensions (emitSystemFieldDimensions count toward the total); the request carries %d",
				logsstore.MaxMetricFilterDimensions, total), 400)
	}
	for k, v := range t.Dimensions {
		if utf8.RuneCountInString(k) > logsstore.MaxMetricDimensionLength || utf8.RuneCountInString(v) > logsstore.MaxMetricDimensionLength {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Dimension key %q and value must be at most %d characters", k, logsstore.MaxMetricDimensionLength), 400)
		}
	}
	if total > 0 && t.DefaultValueSet {
		return NewLogsError("InvalidParameterException",
			"If you assign dimensions to a metric created by a metric filter, you can't assign a default value for that metric", 400)
	}
	return nil
}

// validateEmitSystemFieldDimensions enforces the member's documented
// value set: "Valid values are @aws.account and @aws.region."
func validateEmitSystemFieldDimensions(fields []string) error {
	for _, f := range fields {
		if f != "@aws.account" && f != "@aws.region" {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid emitSystemFieldDimensions entry %s. Valid values: @aws.account, @aws.region", f), 400)
		}
	}
	return nil
}

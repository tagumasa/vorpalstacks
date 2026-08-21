package lambda

import (
	"fmt"
	"regexp"
	"strconv"
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
	"java25", "java21", "java17", "java11",
	"java17.al2023", "java11.al2023", "java8.al2023",
	"java8.al2",
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

// validateEnvironmentVariables enforces the documented environment rules:
// "Environment variables beginning with AWS_LAMBDA_ are reserved" for the
// service's own use, "Keys start with a letter and are at least two
// characters. Keys only contain letters, numbers, and the underscore
// character (_)", and "The total size of all environment variables
// doesn't exceed 4 KB" (keys plus values).
func validateEnvironmentVariables(env *lambdastore.Environment) error {
	if env == nil {
		return nil
	}
	envKeyPattern := regexp.MustCompile(lambdastore.EnvironmentVariableKeyPattern)
	total := 0
	for k, v := range env.Variables {
		if strings.HasPrefix(k, lambdastore.ReservedEnvironmentPrefix) {
			return NewInvalidParameter("Environment.Variables",
				fmt.Sprintf("Environment variable %q uses the reserved AWS_LAMBDA_ prefix", k))
		}
		if !envKeyPattern.MatchString(k) {
			return NewInvalidParameter("Environment.Variables",
				fmt.Sprintf("Environment variable keys must start with a letter and contain only letters, numbers, and underscores; got %q", k))
		}
		total += len(k) + len(v)
	}
	if total > lambdastore.MaxEnvironmentVariablesSizeBytes {
		return NewInvalidParameter("Environment.Variables",
			fmt.Sprintf("The total size of environment variables must not exceed %d bytes, got %d",
				lambdastore.MaxEnvironmentVariablesSizeBytes, total))
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
	// "AT_TIMESTAMP is supported only for Amazon Kinesis streams, Amazon
	// DocumentDB, Amazon MSK, and self-managed Apache Kafka." — the
	// CreateEventSourceMapping model excludes DynamoDB streams from the
	// timestamp start positions.
	if service == "dynamodb" && startingPosition == "AT_TIMESTAMP" {
		return NewInvalidParameter("StartingPosition",
			"AT_TIMESTAMP is supported only for Amazon Kinesis streams, Amazon DocumentDB, Amazon MSK, and self-managed Apache Kafka")
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
// Invocation validation (Smithy InvocationType, LogType enums)
// ---------------------------------------------------------------------------

// validateInvocationType validates the InvocationType parameter per the
// Smithy enum: "Event" (async), "RequestResponse" (sync), "DryRun".
// An empty value defaults to RequestResponse (sync) at the handler level.
func validateInvocationType(v string) error {
	if v == "" {
		return nil
	}
	switch v {
	case "Event", "RequestResponse", "DryRun":
		return nil
	default:
		return NewInvalidParameter("InvocationType",
			fmt.Sprintf("InvocationType must be one of Event, RequestResponse, DryRun; got %q", v))
	}
}

// validateLogType validates the LogType parameter per the Smithy enum:
// "None" (no logs returned) or "Tail" (last 4 KB of logs returned).
// An empty value defaults to None at the handler level.
func validateLogType(v string) error {
	if v == "" {
		return nil
	}
	switch v {
	case "None", "Tail":
		return nil
	default:
		return NewInvalidParameter("LogType",
			fmt.Sprintf("LogType must be None or Tail; got %q", v))
	}
}

// ---------------------------------------------------------------------------
// Alias and StatementId validation (Smithy pattern + length traits)
// ---------------------------------------------------------------------------

// aliasNamePattern enforces the Smithy Alias name pattern:
// alphanumeric, hyphens, and underscores, NOT all-numeric.
var aliasNamePattern = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// validateAliasName validates an alias name per the Smithy model:
// length 1-128, pattern ^(?!^[0-9]+$)[a-zA-Z0-9-_]+$
// (alphanumeric/hyphen/underscore, not purely numeric).
func validateAliasName(name string) error {
	if len(name) < 1 || len(name) > 128 {
		return NewInvalidParameter("Name", "Alias name must be between 1 and 128 characters")
	}
	if !aliasNamePattern.MatchString(name) {
		return NewInvalidParameter("Name", "Alias name can only contain alphanumeric characters, hyphens, and underscores")
	}
	// Reject purely numeric names per the negative lookahead in the Smithy pattern.
	isNumeric := true
	for _, c := range name {
		if c < '0' || c > '9' {
			isNumeric = false
			break
		}
	}
	if isNumeric {
		return NewInvalidParameter("Name", "Alias name cannot be entirely numeric")
	}
	return nil
}

// statementIdPattern enforces the Smithy StatementId pattern:
// ^([a-zA-Z0-9-_]+)$
var statementIdPattern = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// validateStatementId validates a policy statement ID per the Smithy model:
// length 1-100, pattern ^([a-zA-Z0-9-_]+)$.
func validateStatementId(id string) error {
	if len(id) < 1 || len(id) > 100 {
		return NewInvalidParameter("StatementId", "StatementId must be between 1 and 100 characters")
	}
	if !statementIdPattern.MatchString(id) {
		return NewInvalidParameter("StatementId", "StatementId can only contain alphanumeric characters, hyphens, and underscores")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Package type, Architecture, EphemeralStorage, SnapStart validation
// ---------------------------------------------------------------------------

// validatePackageType validates the PackageType per the Smithy enum
// (Zip, Image). An empty value defaults to Zip at the store layer.
func validatePackageType(v string) error {
	if v == "" {
		return nil
	}
	switch v {
	case "Zip", "Image":
		return nil
	default:
		return NewInvalidParameter("PackageType",
			fmt.Sprintf("PackageType must be Zip or Image; got %q", v))
	}
}

// validateArchitecture validates a single architecture value per the
// Smithy enum (x86_64, arm64).
func validateArchitecture(arch string) error {
	switch arch {
	case "x86_64", "arm64":
		return nil
	default:
		return NewInvalidParameter("Architectures",
			fmt.Sprintf("Architecture must be x86_64 or arm64; got %q", arch))
	}
}

// validateEphemeralStorageSize validates the EphemeralStorage.Size per
// the Smithy range: min 512, max 10240 MB.
func validateEphemeralStorageSize(size int32) error {
	if size < 512 || size > 10240 {
		return NewInvalidParameter("EphemeralStorage.Size",
			"EphemeralStorage.Size must be between 512 and 10240 MB")
	}
	return nil
}

// validateSnapStartApplyOn validates the SnapStart.ApplyOn per the
// Smithy enum (PublishedVersions, None).
func validateSnapStartApplyOn(v string) error {
	switch v {
	case "PublishedVersions", "None":
		return nil
	default:
		return NewInvalidParameter("SnapStart.ApplyOn",
			fmt.Sprintf("SnapStart.ApplyOn must be PublishedVersions or None; got %q", v))
	}
}

// snapStartSupportedRuntime reports whether the runtime family supports
// SnapStart, per the AWS documentation: "SnapStart is available for the
// following Lambda managed runtimes: Java 11 and later, Python 3.12 and
// later, .NET 8 and later."
func snapStartSupportedRuntime(runtime string) bool {
	r := strings.ToLower(runtime)
	switch {
	case strings.HasPrefix(r, "java"):
		// Java 8 runtimes predate the Java 11 baseline.
		return r != "java8" && r != "java8.al2"
	case strings.HasPrefix(r, "python3."):
		minor, err := strconv.Atoi(strings.TrimPrefix(r, "python3."))
		return err == nil && minor >= 12
	case r == "dotnet8", r == "dotnet9", r == "dotnet10":
		return true
	default:
		return false
	}
}

// validateSnapStartForRuntime rejects enabling SnapStart on runtimes
// outside the supported families. ApplyOn None disables SnapStart and is
// accepted on every runtime.
func validateSnapStartForRuntime(runtime string, snapStart *lambdastore.SnapStart) error {
	if snapStart == nil || snapStart.ApplyOn != "PublishedVersions" {
		return nil
	}
	if !snapStartSupportedRuntime(runtime) {
		return NewInvalidParameter("SnapStart",
			"SnapStart is supported on Java 11 and later, Python 3.12 and later, and .NET 8 and later runtimes")
	}
	return nil
}

// ---------------------------------------------------------------------------
// KMS Key ARN validation
// ---------------------------------------------------------------------------

// validateKMSKeyArn validates that the KMSKeyArn is a well-formed KMS ARN
// or empty (meaning the default service key is used).
func validateKMSKeyArn(arn string) error {
	if arn == "" {
		return nil
	}
	if !strings.HasPrefix(arn, "arn:") {
		return NewInvalidParameter("KMSKeyArn",
			"KMSKeyArn must be a valid ARN")
	}
	_, service, _, _, resource := arnutil.SplitARN(arn)
	if service != "kms" {
		return NewInvalidParameter("KMSKeyArn",
			fmt.Sprintf("KMSKeyArn must be a KMS ARN; got service %q", service))
	}
	if resource == "" || !strings.HasPrefix(resource, "key/") {
		return NewInvalidParameter("KMSKeyArn",
			"KMSKeyArn must reference a KMS key (arn:...:kms:...:key/<uuid>)")
	}
	return nil
}

// ---------------------------------------------------------------------------
// MaxItems pagination validation (Smithy ranges)
// ---------------------------------------------------------------------------

// maxListItemsCap is the Smithy-specified maximum for Lambda MaxItems
// (ListFunctions: "returns a maximum of 50 items … even if you set the
// number higher").
const maxListItemsCap = 50

// maxEventSourceMappingListItemsCap is the documented ListEventSourceMappings
// maximum: "ListEventSourceMappings returns a maximum of 100 items in each
// response, even if you set the number higher."
const maxEventSourceMappingListItemsCap = 100

// validateMaxItems validates and applies the standard MaxItems range
// (min 1, max 50). A value <= 0 returns the cap.
func validateMaxItems(v int) int {
	return validateMaxItemsCapped(v, maxListItemsCap)
}

// validateMaxItemsCapped validates a MaxItems value against a per-operation
// cap. A value <= 0 or above the cap returns the cap.
func validateMaxItemsCapped(v, cap int) int {
	if v <= 0 || v > cap {
		return cap
	}
	return v
}

// ---------------------------------------------------------------------------
// Event source mapping range validation (Smithy range traits)
// ---------------------------------------------------------------------------

// validateESMBatchSize enforces the Smithy range for BatchSize: min 1, max 10000.
func validateESMBatchSize(v int32) error {
	if v < 1 || v > 10000 {
		return NewInvalidParameter("BatchSize", "BatchSize must be between 1 and 10000")
	}
	return nil
}

// ESM batch size defaults per event source, from the API model docs:
// "Amazon Kinesis – Default 100. Max 10,000. … Amazon DynamoDB Streams –
// Default 100. Max 10,000. … Amazon Simple Queue Service – Default 10. For
// standard queues the max is 10,000. For FIFO queues the max is 10."
const (
	defaultStreamESMBatchSize = int32(100)
	defaultQueueESMBatchSize  = int32(10)
	maxFIFOQueueESMBatchSize  = int32(10)
)

// defaultESMBatchSize returns the batch size AWS applies when the request
// omits BatchSize: 10 for SQS queues and 100 for stream sources.
func defaultESMBatchSize(eventSourceArn string) int32 {
	if arnutil.GetServiceFromARN(eventSourceArn) == "sqs" {
		return defaultQueueESMBatchSize
	}
	return defaultStreamESMBatchSize
}

// validateESMBatchSizeForSource enforces the per-source batch size rules on
// top of the generic range: FIFO queue names end in ".fifo" and accept at
// most 10 records per batch.
func validateESMBatchSizeForSource(v int32, eventSourceArn string) error {
	if err := validateESMBatchSize(v); err != nil {
		return err
	}
	if arnutil.GetServiceFromARN(eventSourceArn) == "sqs" {
		queueName := arnutil.ExtractResourceFromARN(eventSourceArn)
		if strings.HasSuffix(queueName, ".fifo") && v > maxFIFOQueueESMBatchSize {
			return NewInvalidParameter("BatchSize", "BatchSize for FIFO queues must be at most 10")
		}
	}
	return nil
}

// validateESMBatchingWindow enforces the Smithy range for
// MaximumBatchingWindowInSeconds: min 0, max 300.
func validateESMBatchingWindow(v int32) error {
	if v < 0 || v > 300 {
		return NewInvalidParameter("MaximumBatchingWindowInSeconds",
			"MaximumBatchingWindowInSeconds must be between 0 and 300")
	}
	return nil
}

// validateESMParallelFactor enforces the Smithy range for
// ParallelizationFactor: min 1, max 10.
func validateESMParallelFactor(v int32) error {
	if v < lambdastore.MinParallelizationFactor || v > lambdastore.MaxParallelizationFactor {
		return NewInvalidParameter("ParallelizationFactor",
			fmt.Sprintf("ParallelizationFactor must be between %d and %d",
				lambdastore.MinParallelizationFactor, lambdastore.MaxParallelizationFactor))
	}
	return nil
}

// validateESMBatchWindowPair enforces the documented pairing between the
// batch size and the batching window: "For Kinesis, DynamoDB, and Amazon
// SQS event sources, when you set BatchSize to a value greater than 10,
// you must set MaximumBatchingWindowInSeconds to at least 1."
// (CreateEventSourceMapping model, MaximumBatchingWindowInSeconds member
// documentation.) SQS, Kinesis, and DynamoDB streams are exactly the event
// sources this platform polls, so the rule applies unconditionally.
func validateESMBatchWindowPair(batchSize, batchingWindowSeconds int32) error {
	if batchSize > 10 && batchingWindowSeconds < 1 {
		return NewInvalidParameter("MaximumBatchingWindowInSeconds",
			"MaximumBatchingWindowInSeconds must be at least 1 when BatchSize is greater than 10")
	}
	return nil
}

// validateESMMaxRecordAge enforces the Smithy range for
// MaximumRecordAgeInSeconds: min -1, max 604800.
func validateESMMaxRecordAge(v int32) error {
	if v < -1 || v > 604800 {
		return NewInvalidParameter("MaximumRecordAgeInSeconds",
			"MaximumRecordAgeInSeconds must be between -1 and 604800")
	}
	return nil
}

// validateESMMaxRetry enforces the Smithy range for
// MaximumRetryAttemptsEventSourceMapping: min -1, max 10000.
func validateESMMaxRetry(v int32) error {
	if v < -1 || v > 10000 {
		return NewInvalidParameter("MaximumRetryAttempts",
			"MaximumRetryAttempts must be between -1 and 10000")
	}
	return nil
}

// validateESMTumblingWindow enforces the Smithy range for
// TumblingWindowInSeconds: min 0, max 900.
func validateESMTumblingWindow(v int32) error {
	if v < 0 || v > 900 {
		return NewInvalidParameter("TumblingWindowInSeconds",
			"TumblingWindowInSeconds must be between 0 and 900")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Layer permission validation
// ---------------------------------------------------------------------------

// validateLayerPermission validates the Principal and Action fields of a
// layer version resource-based policy statement, applying the same rules
// as validatePermission but for layer-version-scoped policies.
func validateLayerPermission(p *lambdastore.LayerPolicy) error {
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

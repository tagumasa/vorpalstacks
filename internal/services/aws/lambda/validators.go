package lambda

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/pagination"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Smithy-derived regex patterns
// ---------------------------------------------------------------------------

// functionNamePattern validates function names per AWS Lambda docs:
// alphanumeric, hyphens, and underscores, 1-64 characters.
var functionNamePattern = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// namespacedFunctionNamePattern is the NamespacedFunctionName pattern from
// the Smithy model: an optional arn:partition:lambda: prefix followed by
// optional region, account and "function:" segments, a name that may
// contain dots, and an optional qualifier ($LATEST, $LATEST.PUBLISHED or a
// plain alias/version). Members typed by this shape accept every form the
// pattern admits; FunctionName-typed members keep the stricter
// validateFunctionName bound.
var namespacedFunctionNamePattern = regexp.MustCompile(`^(arn:(aws[a-zA-Z-]*)?:lambda:)?((eusc-)?[a-z]{2}((-gov)|(-iso([a-z]?)))?-[a-z]+-\d{1}:)?(\d{12}:)?(function:)?([a-zA-Z0-9-_.]+)(:(\$LATEST(\.PUBLISHED)?|[a-zA-Z0-9-_]+))?$`)

// resourcePolicyArnPattern is the PolicyResourceArn pattern from the
// Smithy model: a complete function ARN, optionally qualified, without
// wildcard characters. The resource-policy operations bind this as their
// ResourceArn path label.
var resourcePolicyArnPattern = regexp.MustCompile(`^arn:(aws[a-zA-Z-]*)?:lambda:(eusc-)?[a-z]{2}((-gov)|(-iso([a-z]?)))?-[a-z]+-\d{1}:\d{12}:function:[a-zA-Z0-9-_]+(:(\$LATEST(\.PUBLISHED)?|[a-zA-Z0-9-_])+)?$`)

// validateResourcePolicyArn validates the ResourceArn of a resource-policy
// operation against the PolicyResourceArn shape (@length 0-256 plus the
// ARN pattern above) and splits it into the function name and an embedded
// qualifier.
func validateResourcePolicyArn(resourceArn string) (functionName, qualifier string, err error) {
	if resourceArn == "" {
		return "", "", NewInvalidParameter("ResourceArn", "ResourceArn is required")
	}
	if len(resourceArn) > lambdastore.MaxPolicyResourceArnLength || !resourcePolicyArnPattern.MatchString(resourceArn) {
		return "", "", NewInvalidParameter("ResourceArn", "ResourceArn must be a complete function ARN without wildcard characters")
	}
	resource := arnutil.ExtractResourceFromARN(resourceArn)
	functionName, qualifier = splitNameQualifier(strings.TrimPrefix(resource, "function:"))
	return functionName, qualifier, nil
}

// ---------------------------------------------------------------------------
// Runtime validation (Smithy Runtime enum)
// ---------------------------------------------------------------------------

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
	if len(name) == 0 || len(name) > maxFunctionNameLength {
		return NewInvalidParameter("FunctionName", fmt.Sprintf("Function name must be between 1 and %d characters", maxFunctionNameLength))
	}
	if !functionNamePattern.MatchString(name) {
		return NewInvalidParameter("FunctionName", "Function name can only contain alphanumeric characters, hyphens, and underscores")
	}
	return nil
}

// validateNamespacedFunctionName validates the whole wire form of a
// FunctionName reference on operations whose member the Smithy model types
// as NamespacedFunctionName. Unlike validateFunctionName, which checks the
// extracted bare name, this checks the raw reference: a malformed ARN
// segment is rejected here instead of passing through to the not-found
// path, while dotted names and the wider reference length the pattern
// admits are accepted and left to resolve (or not) against the store.
func validateNamespacedFunctionName(ref string) error {
	if len(ref) == 0 || len(ref) > maxNamespacedFunctionRefLength {
		return NewInvalidParameter("FunctionName", fmt.Sprintf("Function name must be between 1 and %d characters", maxNamespacedFunctionRefLength))
	}
	if !namespacedFunctionNamePattern.MatchString(ref) {
		return NewInvalidParameter("FunctionName", "Function name reference form is not valid")
	}
	return nil
}

// validateFileSystemConfigs enforces the S3FilesConfig contract on the
// FileSystemConfigs member: DirectS3Read, when provided, must be one of
// the modelled enum values, and an S3FilesConfig is valid only on an
// Amazon S3 Files access point — "If you specify a different access point
// type (for example, Amazon Elastic File System), the operation returns an
// InvalidParameterException".
func validateFileSystemConfigs(configs []lambdastore.FileSystemConfig) error {
	for _, config := range configs {
		if config.S3FilesConfig == nil {
			continue
		}
		if direct := config.S3FilesConfig.DirectS3Read; direct != "" {
			switch direct {
			case lambdastore.DirectS3ReadAuto, lambdastore.DirectS3ReadEnabled, lambdastore.DirectS3ReadDisabled:
			default:
				return NewInvalidParameter("FileSystemConfigs.S3FilesConfig.DirectS3Read",
					fmt.Sprintf("DirectS3Read must be one of '%s', '%s' or '%s'",
						lambdastore.DirectS3ReadAuto, lambdastore.DirectS3ReadEnabled, lambdastore.DirectS3ReadDisabled))
			}
		}
		if arnutil.GetServiceFromARN(config.Arn) != "s3files" {
			return NewInvalidParameter("FileSystemConfigs.S3FilesConfig",
				fmt.Sprintf("S3FilesConfig is valid only on an Amazon S3 Files access point ARN, got %q", config.Arn))
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Timeout and MemorySize validation (Smithy range traits)
// ---------------------------------------------------------------------------

func validateTimeout(timeout int32) error {
	if timeout < lambdastore.MinTimeoutSeconds || timeout > lambdastore.MaxTimeoutSeconds {
		return NewInvalidParameter("Timeout",
			fmt.Sprintf("Timeout must be between %d and %d seconds",
				lambdastore.MinTimeoutSeconds, lambdastore.MaxTimeoutSeconds))
	}
	return nil
}

func validateMemorySize(memorySize int32) error {
	if memorySize < lambdastore.MinMemorySizeMB || memorySize > lambdastore.MaxMemorySizeMB {
		return NewInvalidParameter("MemorySize",
			fmt.Sprintf("MemorySize must be between %d and %d MB",
				lambdastore.MinMemorySizeMB, lambdastore.MaxMemorySizeMB))
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

// validateCorsConfig enforces the modelled CORS bounds on a function URL
// configuration: Cors.MaxAge targets the MaxAge shape whose @range is 0
// to 86400 seconds.
func validateCorsConfig(cors *lambdastore.CorsConfig) error {
	if cors == nil {
		return nil
	}
	if cors.MaxAge < 0 || cors.MaxAge > lambdastore.MaxCorsMaxAgeSeconds {
		return NewInvalidParameter("Cors.MaxAge",
			fmt.Sprintf("Cors.MaxAge must be between 0 and %d seconds; got %d",
				lambdastore.MaxCorsMaxAgeSeconds, cors.MaxAge))
	}
	return nil
}

// ---------------------------------------------------------------------------
// Event invoke config validation (Smithy range traits)
// ---------------------------------------------------------------------------

// validateMaximumEventAgeInSeconds enforces the Smithy range
// MaximumEventAgeInSeconds: min 60, max 21600.
func validateMaximumEventAgeInSeconds(v int32) error {
	if v < lambdastore.MinEventAgeSeconds || v > lambdastore.MaxEventAgeSeconds {
		return NewInvalidParameter("MaximumEventAgeInSeconds",
			fmt.Sprintf("MaximumEventAgeInSeconds must be between %d and %d seconds",
				lambdastore.MinEventAgeSeconds, lambdastore.MaxEventAgeSeconds))
	}
	return nil
}

// validateMaximumRetryAttempts enforces the Smithy range
// MaximumRetryAttempts: min 0, max 2.
func validateMaximumRetryAttempts(v int32) error {
	if v < lambdastore.MinEventInvokeRetryAttempts || v > lambdastore.MaxEventInvokeRetryAttempts {
		return NewInvalidParameter("MaximumRetryAttempts",
			fmt.Sprintf("MaximumRetryAttempts must be between %d and %d",
				lambdastore.MinEventInvokeRetryAttempts, lambdastore.MaxEventInvokeRetryAttempts))
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
	if len(name) < 1 || len(name) > maxAliasNameLength {
		return NewInvalidParameter("Name", fmt.Sprintf("Alias name must be between 1 and %d characters", maxAliasNameLength))
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
	if len(id) < 1 || len(id) > lambdastore.MaxStatementIdLength {
		return NewInvalidParameter("StatementId", fmt.Sprintf("StatementId must be between 1 and %d characters", lambdastore.MaxStatementIdLength))
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
// the Smithy range: min 512, max 32768 MB.
func validateEphemeralStorageSize(size int32) error {
	if size < lambdastore.MinEphemeralStorageSizeMB || size > lambdastore.MaxEphemeralStorageSizeMB {
		return NewInvalidParameter("EphemeralStorage.Size",
			fmt.Sprintf("EphemeralStorage.Size must be between %d and %d MB",
				lambdastore.MinEphemeralStorageSizeMB, lambdastore.MaxEphemeralStorageSizeMB))
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
	case r == "dotnet8", r == "dotnet10":
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
// LoggingConfig / ImageConfig validation (model enums, lengths, pattern)
// ---------------------------------------------------------------------------

// logGroupPattern enforces the modelled LogGroup shape: letters, digits,
// and the characters . - _ / #, length 1-512.
var logGroupPattern = regexp.MustCompile(lambdastore.LogGroupPattern)

// validateLoggingConfig enforces the modelled LoggingConfig member
// constraints: the LogFormat enum (JSON, Text), the ApplicationLogLevel
// enum (TRACE, DEBUG, INFO, WARN, ERROR, FATAL), the SystemLogLevel enum
// (DEBUG, INFO, WARN), and the LogGroup length and pattern. Empty members
// keep their platform defaults and are accepted.
func validateLoggingConfig(lc *lambdastore.LoggingConfig) error {
	if lc == nil {
		return nil
	}
	switch lc.LogFormat {
	case "", "JSON", "Text":
	default:
		return NewInvalidParameter("LoggingConfig.LogFormat",
			fmt.Sprintf("LogFormat must be JSON or Text; got %q", lc.LogFormat))
	}
	switch lc.ApplicationLogLevel {
	case "", "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL":
	default:
		return NewInvalidParameter("LoggingConfig.ApplicationLogLevel",
			fmt.Sprintf("ApplicationLogLevel must be one of TRACE, DEBUG, INFO, WARN, ERROR, FATAL; got %q", lc.ApplicationLogLevel))
	}
	switch lc.SystemLogLevel {
	case "", "DEBUG", "INFO", "WARN":
	default:
		return NewInvalidParameter("LoggingConfig.SystemLogLevel",
			fmt.Sprintf("SystemLogLevel must be one of DEBUG, INFO, WARN; got %q", lc.SystemLogLevel))
	}
	if lc.LogGroup != "" {
		if len(lc.LogGroup) > lambdastore.MaxLogGroupLength {
			return NewInvalidParameter("LoggingConfig.LogGroup",
				fmt.Sprintf("LogGroup must be at most %d characters", lambdastore.MaxLogGroupLength))
		}
		if !logGroupPattern.MatchString(lc.LogGroup) {
			return NewInvalidParameter("LoggingConfig.LogGroup",
				fmt.Sprintf("LogGroup may only contain letters, numbers, and the characters . - _ / #; got %q", lc.LogGroup))
		}
	}
	return nil
}

// validateImageConfig enforces the modelled ImageConfig member constraints:
// EntryPoint and Command target the StringList shape (at most 1500 entries
// each) and WorkingDirectory is at most 1000 characters.
func validateImageConfig(ic *lambdastore.ImageConfig) error {
	if ic == nil {
		return nil
	}
	if len(ic.EntryPoint) > lambdastore.MaxImageConfigListLength {
		return NewInvalidParameter("ImageConfig.EntryPoint",
			fmt.Sprintf("ImageConfig.EntryPoint must hold at most %d entries", lambdastore.MaxImageConfigListLength))
	}
	if len(ic.Command) > lambdastore.MaxImageConfigListLength {
		return NewInvalidParameter("ImageConfig.Command",
			fmt.Sprintf("ImageConfig.Command must hold at most %d entries", lambdastore.MaxImageConfigListLength))
	}
	if len(ic.WorkingDirectory) > lambdastore.MaxImageConfigWorkingDirectory {
		return NewInvalidParameter("ImageConfig.WorkingDirectory",
			fmt.Sprintf("ImageConfig.WorkingDirectory must be at most %d characters", lambdastore.MaxImageConfigWorkingDirectory))
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
	return pagination.ClampMaxItems(v, cap, cap)
}

// ---------------------------------------------------------------------------
// Event source mapping range validation (Smithy range traits)
// ---------------------------------------------------------------------------

// validateESMBatchSize enforces the Smithy range for BatchSize: min 1, max 10000.
func validateESMBatchSize(v int32) error {
	if v < lambdastore.MinESMBatchSize || v > lambdastore.MaxESMBatchSize {
		return NewInvalidParameter("BatchSize",
			fmt.Sprintf("BatchSize must be between %d and %d",
				lambdastore.MinESMBatchSize, lambdastore.MaxESMBatchSize))
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

// clampESMBatchSize applies the defensive default-and-cap every poller
// path runs on the stored BatchSize: a non-positive value falls back to
// the per-source documented default and the result is capped at the
// modelled maximum. The SQS fallback is 10 — the documented queue
// default — never the stream default 100.
func clampESMBatchSize(batchSize int32, eventSourceArn string) int32 {
	if batchSize <= 0 {
		batchSize = defaultESMBatchSize(eventSourceArn)
	}
	if batchSize > lambdastore.MaxESMBatchSize {
		batchSize = lambdastore.MaxESMBatchSize
	}
	return batchSize
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
	if v < lambdastore.MinESMBatchingWindowSeconds || v > lambdastore.MaxESMBatchingWindowSeconds {
		return NewInvalidParameter("MaximumBatchingWindowInSeconds",
			fmt.Sprintf("MaximumBatchingWindowInSeconds must be between %d and %d",
				lambdastore.MinESMBatchingWindowSeconds, lambdastore.MaxESMBatchingWindowSeconds))
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

// validateESMDestinationConfig enforces the modelled DestinationConfig
// contract for event source mappings: the model marks OnSuccess as "not
// supported in CreateEventSourceMapping or UpdateEventSourceMapping
// requests", and an on-failure destination must be the ARN of an SQS queue,
// SNS topic, or S3 bucket — the destinations the delivery path implements
// for stream sources.
func validateESMDestinationConfig(dc *lambdastore.DestinationConfig) error {
	if dc == nil {
		return nil
	}
	if dc.OnSuccess != nil && dc.OnSuccess.Destination != "" {
		return NewInvalidParameter("DestinationConfig.OnSuccess",
			"OnSuccess destinations are not supported for event source mappings")
	}
	if dc.OnFailure != nil && dc.OnFailure.Destination != "" {
		_, service, _, _, _ := arnutil.SplitARN(dc.OnFailure.Destination)
		switch service {
		case "sqs", "sns", "s3":
		default:
			return NewInvalidParameter("DestinationConfig.OnFailure",
				fmt.Sprintf("OnFailure destination must be an SQS queue, SNS topic, or S3 bucket ARN; got service %q", service))
		}
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
const esmBatchWindowPairingThreshold = int32(10)

func validateESMBatchWindowPair(batchSize, batchingWindowSeconds int32) error {
	if batchSize > esmBatchWindowPairingThreshold && batchingWindowSeconds < 1 {
		return NewInvalidParameter("MaximumBatchingWindowInSeconds",
			"MaximumBatchingWindowInSeconds must be at least 1 when BatchSize is greater than 10")
	}
	return nil
}

// validateESMMaxRecordAge enforces the Smithy range for
// MaximumRecordAgeInSeconds: min -1, max 604800.
func validateESMMaxRecordAge(v int32) error {
	if v < lambdastore.MinESMRecordAgeSeconds || v > lambdastore.MaxESMRecordAgeSeconds {
		return NewInvalidParameter("MaximumRecordAgeInSeconds",
			fmt.Sprintf("MaximumRecordAgeInSeconds must be between %d and %d",
				lambdastore.MinESMRecordAgeSeconds, lambdastore.MaxESMRecordAgeSeconds))
	}
	return nil
}

// validateESMMaxRetry enforces the Smithy range for
// MaximumRetryAttemptsEventSourceMapping: min -1, max 10000.
func validateESMMaxRetry(v int32) error {
	if v < lambdastore.MinESMRetryAttempts || v > lambdastore.MaxESMRetryAttempts {
		return NewInvalidParameter("MaximumRetryAttempts",
			fmt.Sprintf("MaximumRetryAttempts must be between %d and %d",
				lambdastore.MinESMRetryAttempts, lambdastore.MaxESMRetryAttempts))
	}
	return nil
}

// validateESMTumblingWindow enforces the Smithy range for
// TumblingWindowInSeconds: min 0, max 900.
func validateESMTumblingWindow(v int32) error {
	if v < lambdastore.MinESMTumblingWindowSeconds || v > lambdastore.MaxESMTumblingWindowSeconds {
		return NewInvalidParameter("TumblingWindowInSeconds",
			fmt.Sprintf("TumblingWindowInSeconds must be between %d and %d",
				lambdastore.MinESMTumblingWindowSeconds, lambdastore.MaxESMTumblingWindowSeconds))
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
	return validatePolicyStatement(p.Principal, p.Action)
}

// validatePolicyStatement holds the shared member contract of the function
// and layer-version resource policies: a required Principal (an IAM ARN, a
// recognised service principal, or "*") and a required Action scoped to the
// "lambda:" namespace.
func validatePolicyStatement(principal, action string) error {
	if principal == "" {
		return NewInvalidParameter("Principal", "Principal is required")
	}
	if !isValidPrincipal(principal) {
		return NewInvalidParameter("Principal",
			fmt.Sprintf("Principal %q is not a valid IAM ARN, recognised service principal, or wildcard", principal))
	}
	if action == "" {
		return NewInvalidParameter("Action", "Action is required")
	}
	if !strings.HasPrefix(action, "lambda:") {
		return NewInvalidParameter("Action",
			fmt.Sprintf("Action %q must start with 'lambda:'", action))
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
	return validatePolicyStatement(p.Principal, p.Action)
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

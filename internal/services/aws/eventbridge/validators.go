package eventbridge

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Pattern constants extracted from the Smithy model
// (third_party/api-models-aws/models/eventbridge/service/2015-10-07/).
//
// resourceNamePattern applies to RuleName, ArchiveName, ConnectionName,
// ApiDestinationName and NonPartnerEventBusName (Smithy @pattern
// ^[\.\-_A-Za-z0-9]+$).
var resourceNamePattern = regexp.MustCompile(`^[\.\-_A-Za-z0-9]+$`)

// eventBusNamePattern applies to EventBusName (Smithy @pattern includes '/').
var eventBusNamePattern = regexp.MustCompile(`^[/\.\-_A-Za-z0-9]+$`)

// eventBusArnPattern applies to EventBusArn (Smithy @pattern
// ^arn:aws([a-z]|\-)*:events:([a-z]|\d|\-)*:([0-9]{12})?:.+\/.+$; the
// account-id segment is optional).
var eventBusArnPattern = regexp.MustCompile(`^arn:aws([a-z]|\-)*:events:([a-z]|\d|\-)*:([0-9]{12})?:.+\/.+$`)

// ---------------------------------------------------------------------------
// Name validators
// ---------------------------------------------------------------------------

// validateResourceName validates RuleName (1-64), ArchiveName (1-48),
// ConnectionName (1-64) and ApiDestinationName (1-64) against the Smithy
// @length and @pattern traits.
func validateResourceName(name, kind string) bool {
	if name == "" {
		return false
	}
	maxLen := 64
	if kind == "archive" {
		maxLen = 48
	}
	if len(name) > maxLen {
		return false
	}
	return resourceNamePattern.MatchString(name)
}

// validateEventBusName validates EventBusName per Smithy:
// pattern ^[/\.\-_A-Za-z0-9]+$, length 1-256.
func validateEventBusName(name string) bool {
	if name == "" || len(name) > 256 {
		return false
	}
	return eventBusNamePattern.MatchString(name)
}

// validateNonPartnerEventBusName validates NonPartnerEventBusName per
// Smithy: pattern ^[\.\-_A-Za-z0-9]+$, length 1-256. This is the permission
// family's bus reference — the character class excludes both the ARN form
// (colons) and the slash-bearing partner event-source names.
func validateNonPartnerEventBusName(name string) bool {
	if name == "" || len(name) > 256 {
		return false
	}
	return resourceNamePattern.MatchString(name)
}

// validateEventBusArn validates EventBusArn per Smithy: pattern
// ^arn:aws([a-z]|\-)*:events:([a-z]|\d|\-)*:([0-9]{12})?:.+\/.+$, length
// 1-1600 — the archive family's EventSourceArn member form.
func validateEventBusArn(arn string) bool {
	if arn == "" || len(arn) > 1600 {
		return false
	}
	return eventBusArnPattern.MatchString(arn)
}

// validateReplayName validates ReplayName per Smithy:
// pattern ^[\.\-_A-Za-z0-9]+$, length 1-64.
func validateReplayName(name string) bool {
	if name == "" || len(name) > 64 {
		return false
	}
	return resourceNamePattern.MatchString(name)
}

// ---------------------------------------------------------------------------
// List query validation
// ---------------------------------------------------------------------------

// normaliseListLimit applies the model's list Limit window: every list
// operation's Limit member targets the LimitMax100 shape (@range 1-100,
// vendored model). An omitted member arrives as the int32 zero value and
// defaults to the maximum, so an unbounded request pages the full window
// instead of a smaller silent default.
func normaliseListLimit(limit int32) (int32, error) {
	if limit == 0 {
		return eventsstore.ListLimitMaximum, nil
	}
	if limit < 1 || limit > eventsstore.ListLimitMaximum {
		return 0, awserrors.NewValidationException(
			fmt.Sprintf("Limit must be between 1 and %d", eventsstore.ListLimitMaximum))
	}
	return limit, nil
}

// validateListNamePrefix validates a NamePrefix against the charset and
// length window of the family's name shape: every NamePrefix member
// targets the family's name shape (jq over the vendored model), and those
// patterns are pure charsets, so a prefix satisfying the shape is exactly
// a prefix some valid name can have. An invalid prefix is a
// ValidationException rather than a silent empty match.
func validateListNamePrefix(prefix, kind string) error {
	var ok bool
	switch kind {
	case "event-bus":
		ok = validateEventBusName(prefix)
	case "replay":
		ok = validateReplayName(prefix)
	default:
		ok = validateResourceName(prefix, kind)
	}
	if !ok {
		return awserrors.NewValidationException(
			"NamePrefix must match the " + kind + " name pattern and length window")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Enum validators
// ---------------------------------------------------------------------------

// validRuleStates mirrors the Smithy RuleState enum which has three members:
// ENABLED, DISABLED and ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS.
var validRuleStates = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
	"ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS": true,
}

// validateRuleState returns true when *state* is a valid Smithy RuleState enum
// value (ENABLED, DISABLED or ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS).
func validateRuleState(state string) bool {
	return validRuleStates[state]
}

// ---------------------------------------------------------------------------
// Length validators (AWS docs — Smithy targets unadorned String)
// ---------------------------------------------------------------------------

// validateDetailType enforces the AWS PutEventsRequestEntry.DetailType maximum
// length of 128 characters, counted in Unicode characters (member
// documentation: "maximum of 128 characters").
func validateDetailType(s string) bool {
	n := utf8.RuneCountInString(s)
	return n > 0 && n <= 128
}

// validateSource enforces the AWS PutEventsRequestEntry.Source constraint
// "Length constraints: minimum length of 1, maximum length of 256" (API
// reference; the model member itself carries no length documentation),
// counted in Unicode characters like the sibling DetailType bound.
func validateSource(s string) bool {
	n := utf8.RuneCountInString(s)
	return n > 0 && n <= 256
}

// validateTraceHeader enforces the Smithy TraceHeader @length(min=1,max=500),
// counted in Unicode characters (the shape carries no pattern). An empty
// string is accepted because the parameter is optional; a non-empty value
// must satisfy max=500.
func validateTraceHeader(s string) bool {
	return utf8.RuneCountInString(s) <= 500
}

// ---------------------------------------------------------------------------
// Range validators
// ---------------------------------------------------------------------------

// validateRetryPolicy validates MaximumRetryAttempts (0-185) and
// MaximumEventAgeInSeconds (60-86400) per the modelled @range traits and
// the API reference (RetryPolicy). A nil policy is valid. The age member's
// zero value means "not set": a caller that parsed the member off the wire
// passes maxAgeProvided, so an explicitly supplied zero — below the
// modelled minimum of 60 — is rejected instead of read as the default,
// while an omitted member keeps the deployment default (the delivery path
// applies in-range stored values only).
func validateRetryPolicy(rp *eventsstore.RetryPolicy, maxAgeProvided bool) bool {
	if rp == nil {
		return true
	}
	if rp.MaximumRetryAttempts < 0 || rp.MaximumRetryAttempts > eventsstore.RetryPolicyMaxRetryAttempts {
		return false
	}
	if rp.MaximumEventAgeInSeconds != 0 {
		if rp.MaximumEventAgeInSeconds < eventsstore.RetryPolicyMinEventAgeSeconds || rp.MaximumEventAgeInSeconds > eventsstore.RetryPolicyMaxEventAgeSeconds {
			return false
		}
	} else if maxAgeProvided {
		return false
	}
	return true
}

// inputPathsMapKeyPattern is the Smithy @pattern on
// InputTransformerPathKey. Because the pattern admits no dots, a key can
// never collide with the reserved aws.events.* variable names, and the
// separate "keys cannot start with \"AWS.\"" note is subsumed as well.
var inputPathsMapKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_\-]+$`)

// placeholderObjectKeyPattern matches a quoted placeholder in object-key
// position. The restriction applies when the template is a JSON object
// ("If InputTemplate is a JSON object (surrounded by curly braces) ...
// The placeholder cannot be used as an object key", Smithy model doc).
var placeholderObjectKeyPattern = regexp.MustCompile(`"<[^">]*>"\s*:`)

// validateTargetInputConfiguration enforces the input-configuration
// contract on one PutTargets target: "Input, InputPath, and
// InputTransformer are mutually exclusive and optional parameters of a
// target" (PutTargets API reference); Input is valid JSON text within the
// TargetInput length trait; InputPath stays within the TargetInputPath
// bound; and the transformer obeys the TransformerInput,
// TransformerPaths and InputTransformerPathKey traits (template length,
// map entry count, key pattern and path lengths). Every length is counted
// in Unicode scalar values — the unit the Smithy @length trait defines
// for strings — so multibyte input below a bound is not over-rejected.
func validateTargetInputConfiguration(target *eventsstore.Target) error {
	configured := 0
	if target.Input != "" {
		configured++
	}
	if target.InputPath != "" {
		configured++
	}
	if target.InputTransformer != nil {
		configured++
	}
	if configured > 1 {
		return fmt.Errorf("Input, InputPath, and InputTransformer are mutually exclusive and cannot be specified together")
	}

	if target.Input != "" {
		if utf8.RuneCountInString(target.Input) > eventsstore.TargetInputMaxLength {
			return fmt.Errorf("Input must be at most %d characters", eventsstore.TargetInputMaxLength)
		}
		if !json.Valid([]byte(target.Input)) {
			return fmt.Errorf("Input must be valid JSON text")
		}
	}

	if utf8.RuneCountInString(target.InputPath) > eventsstore.TargetInputPathMaxLength {
		return fmt.Errorf("InputPath must be at most %d characters", eventsstore.TargetInputPathMaxLength)
	}

	if target.InputTransformer != nil {
		template := target.InputTransformer.InputTemplate
		if n := utf8.RuneCountInString(template); n < eventsstore.InputTemplateMinLength || n > eventsstore.InputTemplateMaxLength {
			return fmt.Errorf("InputTemplate must be between %d and %d characters",
				eventsstore.InputTemplateMinLength, eventsstore.InputTemplateMaxLength)
		}
		if len(target.InputTransformer.InputPathsMap) > eventsstore.InputPathsMapMaxEntries {
			return fmt.Errorf("InputPathsMap must contain at most %d entries", eventsstore.InputPathsMapMaxEntries)
		}
		for key, path := range target.InputTransformer.InputPathsMap {
			if utf8.RuneCountInString(key) > eventsstore.InputPathsMapKeyMaxLength || !inputPathsMapKeyPattern.MatchString(key) {
				return fmt.Errorf("InputPathsMap key %q must match [A-Za-z0-9_-] and be at most %d characters",
					key, eventsstore.InputPathsMapKeyMaxLength)
			}
			if utf8.RuneCountInString(path) > eventsstore.TargetInputPathMaxLength {
				return fmt.Errorf("InputPathsMap value for key %q must be at most %d characters",
					key, eventsstore.TargetInputPathMaxLength)
			}
		}
		if strings.HasPrefix(strings.TrimSpace(template), "{") &&
			placeholderObjectKeyPattern.MatchString(template) {
			return fmt.Errorf("the placeholder cannot be used as an object key in a JSON object InputTemplate")
		}
	}

	return nil
}

// validateDeadLetterQueueARN enforces the documented DLQ contract at
// PutTargets: the ARN names a standard SQS queue ("The ARN of the SQS
// queue specified as the target for the dead-letter queue", Smithy model
// doc; "Only standard queues are supported. You can't use a FIFO queue
// for a DLQ", user guide) in the rule's own region ("The Amazon SQS queue
// you use must be in the same Region in which you create the rule"). An
// empty ARN means no DLQ and passes; an empty rule region skips the
// region match for callers without a region context.
func validateDeadLetterQueueARN(arn, ruleRegion string) error {
	if arn == "" {
		return nil
	}
	_, service, arnRegion, _, resource := svcarn.SplitARN(arn)
	if service != "sqs" {
		return fmt.Errorf("DeadLetterConfig.Arn %s: EventBridge dead-letter queues are SQS queues, got service %q", arn, service)
	}
	if strings.HasSuffix(resource, ".fifo") {
		return fmt.Errorf("DeadLetterConfig.Arn %s: only standard queues are supported as EventBridge dead-letter queues", arn)
	}
	if ruleRegion != "" && arnRegion != ruleRegion {
		return fmt.Errorf("DeadLetterConfig.Arn %s: the dead-letter queue must be in the same region as the rule (%s)", arn, ruleRegion)
	}
	return nil
}

// validateKmsKeyIdentifier verifies that the value is a usable KMS key
// identifier. The model documentation for the member (CreateConnection and
// CreateArchive requests) states: "The identifier can be the key Amazon
// Resource Name (ARN), KeyId, key alias, or key alias ARN." — so a bare key
// ID or alias is as valid as a full ARN. An ARN form must be a KMS ARN
// naming a key or alias resource; a non-ARN value must be either an alias
// path (alias/...) or a bare key ID. An empty string is valid — the
// parameter is optional. Every form shares the shape's @length(0,2048)
// ceiling, counted in Unicode scalar values.
func validateKmsKeyIdentifier(identifier string) bool {
	if identifier == "" {
		return true
	}
	if utf8.RuneCountInString(identifier) > eventsstore.KmsKeyIdentifierMaxLength {
		return false
	}
	if strings.HasPrefix(identifier, "arn:") {
		_, service, _, _, resource := svcarn.SplitARN(identifier)
		if service != "kms" {
			return false
		}
		return strings.HasPrefix(resource, "key/") || strings.HasPrefix(resource, "alias/")
	}
	if strings.HasPrefix(identifier, "alias/") {
		return kmsAliasPathPattern.MatchString(identifier)
	}
	return kmsKeyIDPattern.MatchString(identifier)
}

// kmsAliasPathPattern constrains a bare key alias path (the characters AWS
// permits in alias names).
var kmsAliasPathPattern = regexp.MustCompile(`^alias/[a-zA-Z0-9/_-]+$`)

// kmsKeyIDPattern constrains a bare KMS key ID: the 36-character key
// identifiers AWS assigns (UUID-shaped, including the mrk- multi-Region
// prefix).
var kmsKeyIDPattern = regexp.MustCompile(`^[a-zA-Z0-9-]{36}$`)

// connectionArnPattern is the Smithy @pattern on ConnectionArn:
// arn:<partition>:events:<region>:<account>:connection/<name>/<id>.
var connectionArnPattern = regexp.MustCompile(
	`^arn:aws([a-z]|\-)*:events:([a-z]|\d|\-)*:([0-9]{12})?:connection\/[\.\-_A-Za-z0-9]+\/[\-A-Za-z0-9]+$`)

// apiDestinationArnPattern is the Smithy @pattern on
// ApiDestinationArn: arn:<partition>:events:<region>:<account>:
// api-destination/<name>/<id>.
var apiDestinationArnPattern = regexp.MustCompile(
	`^arn:aws([a-z]|\-)*:events:([a-z]|\d|\-)*:([0-9]{12})?:api-destination\/[\.\-_A-Za-z0-9]+\/[\-A-Za-z0-9]+$`)

// The four connection HTTP-parameter patterns below are the Smithy
// @pattern traits copied verbatim (HeaderKey, HeaderValueSensitive,
// QueryStringKey, QueryStringValueSensitive). The header-key class keeps
// the model's "+-." range exactly as written; the header-value form is
// printable ASCII runs separated by spaces or tabs; the query-string
// forms exclude control characters, the value form still admitting CR
// and LF.
var (
	connectionHeaderKeyPattern = regexp.MustCompile("^[!#$%&'*+-.^_`|~0-9a-zA-Z]+$")
	// The interpreted-string escapes below (\\t, \\x20-\\x7E, ...) become
	// the literal control characters in the compiled class, which is
	// semantically identical to the model's escape sequences.
	connectionHeaderValuePattern      = regexp.MustCompile("^[ \t]*[\x20-\x7E]+([ \t]+[\x20-\x7E]+)*[ \t]*$")
	connectionQueryStringKeyPattern   = regexp.MustCompile("^[^\x00-\x1F\x7F]+$")
	connectionQueryStringValuePattern = regexp.MustCompile("^[^\x00-\x09\x0B\x0C\x0E-\x1F\x7F]+$")
)

// invocationEndpointPattern is the Smithy HttpsEndpoint @pattern (a URI
// component charset — every character of an https URL, including the
// scheme separators, is inside it). The https scheme requirement comes from
// the console ("The endpoint URL must start with https") and the user
// guide's description of API destinations as HTTPS endpoints.
var invocationEndpointPattern = regexp.MustCompile(
	`^((%[0-9A-Fa-f]{2}|[-()_.!~*';/?:@&=+$,A-Za-z0-9])+)([).!';/?:,])?$`)

// validateInvocationEndpoint enforces the HttpsEndpoint contract on an API
// destination invocation endpoint: https scheme, @length(1,2048) and the
// shape's URI charset pattern.
func validateInvocationEndpoint(endpoint string) bool {
	if !strings.HasPrefix(endpoint, "https://") {
		return false
	}
	if len(endpoint) < 1 || len(endpoint) > 2048 {
		return false
	}
	return invocationEndpointPattern.MatchString(endpoint)
}

// ---------------------------------------------------------------------------
// Additional length validators (Smithy @length traits)
// ---------------------------------------------------------------------------

// maxDescriptionLength is the Smithy @length(0,512) bound shared by
// RuleDescription, EventBusDescription, ArchiveDescription,
// ConnectionDescription, ApiDestinationDescription and ReplayDescription.
const maxDescriptionLength = 512

// maxEventBusPolicyBytes is the documented event bus resource policy size
// ceiling: "The permission policy on the event bus cannot exceed 10 KB in
// size" (PutPermission API Reference and the operation documentation in the
// Smithy model). The bound is applied to the byte size: AWS reports
// violations as "Maximum policy size of 10240 bytes exceeded", even though
// the quotas page phrases the same limit as 10,240 characters. JSON policy
// documents are ASCII-dominated, so the two readings only diverge for
// multibyte policy text.
const maxEventBusPolicyBytes = 10240

// validateEventBusPolicySize enforces the documented 10 KB ceiling on the
// event bus resource policy supplied via PutPermission's Policy parameter.
func validateEventBusPolicySize(policy string) error {
	if len(policy) > maxEventBusPolicyBytes {
		return awserrors.NewPolicyLengthExceededException(
			fmt.Sprintf("Event bus policy length %d exceeds the maximum allowed length of %d bytes", len(policy), maxEventBusPolicyBytes))
	}
	return nil
}

// validateDescription enforces the Smithy @length(0,512) trait shared by
// RuleDescription, EventBusDescription, ArchiveDescription,
// ConnectionDescription, ApiDestinationDescription and ReplayDescription,
// counted in Unicode characters (none of the shapes carries a pattern).
func validateDescription(s string) bool {
	return utf8.RuneCountInString(s) <= maxDescriptionLength
}

// errDescriptionTooLong is the validation error for a description member
// exceeding the shared @length(0,512) bound.
func errDescriptionTooLong() error {
	return awserrors.NewValidationException(
		fmt.Sprintf("Description must be at most %d characters", maxDescriptionLength))
}

// validateEventPatternLength enforces the Smithy EventPattern @length(0,4096),
// counted in Unicode characters (the shape carries no pattern, so JSON
// patterns with multibyte values are rune-legal). The caller should
// additionally check the pattern structure via validateEventPatternStructure.
func validateEventPatternLength(pattern string) bool {
	return utf8.RuneCountInString(pattern) <= 4096
}

// validateInvocationRateLimit enforces the model's @range(min:1) on
// InvocationRateLimitPerSecond (the API reference documents the same
// "Minimum value of 1" and no maximum). Zero means the member was omitted
// and the destination is not throttled below the service limits.
func validateInvocationRateLimit(rate int32) bool {
	return rate >= 1
}

// validateRetentionDays enforces the model's @range(min:0) on the archive
// RetentionDays member. Zero is the documented default and retains events
// indefinitely ("If set to 0, events are retained indefinitely", API
// reference); only negative values are out of range. The model declares
// no maximum.
func validateRetentionDays(days int32) bool {
	return days >= 0
}

// isValidLogIncludeDetail validates the bus LogConfig IncludeDetail member
// against the Smithy enum values: NONE or FULL.
func isValidLogIncludeDetail(v string) bool {
	return v == "NONE" || v == "FULL"
}

// isValidLogLevel validates the bus LogConfig Level member against the
// Smithy enum values: OFF, ERROR, INFO, or TRACE.
func isValidLogLevel(v string) bool {
	switch v {
	case "OFF", "ERROR", "INFO", "TRACE":
		return true
	}
	return false
}

// validOAuthHttpMethods is the Smithy ConnectionOAuthHttpMethod enum: an
// OAuth token request supports GET, POST and PUT only (a strict subset of
// the API destination method enum).
var validOAuthHttpMethods = map[string]bool{
	"GET":  true,
	"POST": true,
	"PUT":  true,
}

// validAuthTypes is the Smithy ConnectionAuthorizationType enum.
var validAuthTypes = map[string]bool{
	"API_KEY":                  true,
	"BASIC":                    true,
	"OAUTH_CLIENT_CREDENTIALS": true,
}

// validHttpMethods is the Smithy ApiDestinationHttpMethod enum used by the
// API destination family.
var validHttpMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"HEAD":    true,
	"OPTIONS": true,
	"PATCH":   true,
}

// ---------------------------------------------------------------------------
// PutPermission statement-mode member validators (API Reference)
// ---------------------------------------------------------------------------

// The statement-mode members carry documented patterns: Action
// `events:[a-zA-Z]+` with length 1-64, Principal `(\d{12}|\*)` (the
// 12-digit account ID or the any-account wildcard) and StatementId
// `[a-zA-Z0-9-_]+` with length 1-64. The Condition members support exactly
// one documented pair: Type "StringEquals" with Key "aws:PrincipalOrgID"
// ("Currently the only supported value is StringEquals" / "Currently the
// only supported key is aws:PrincipalOrgID" — Smithy Condition member
// documentation).
var putPermissionActionPattern = regexp.MustCompile(`^events:[a-zA-Z]+$`)
var putPermissionPrincipalPattern = regexp.MustCompile(`^(\d{12}|\*)$`)
var putPermissionStatementIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

const maxPutPermissionStatementIDLength = 64
const maxPutPermissionActionLength = 64

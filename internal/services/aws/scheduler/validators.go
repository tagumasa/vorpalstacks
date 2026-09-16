package scheduler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/scheduleexpr"
	tagutil "vorpalstacks/internal/common/tags"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// Smithy @length maxima. Every @length bound counts Unicode characters
// (the trait's counting basis), so each maximum is checked against the
// rune count; charset constraints (@pattern) stay regex matches and are
// unaffected by the basis.
const (
	maxTimezoneLength         = 50  // ScheduleExpressionTimezone
	maxScheduleDescriptionLen = 512 // Description
	maxPlatformVersionLength  = 64  // PlatformVersion
	maxEcsGroupLength         = 255 // Group
	maxReferenceIdLength      = 1024
	// PlacementConstraintExpression @length(max 2000) and
	// PlacementStrategyField @length(max 255).
	maxPlacementExpressionLength    = 2000
	maxPlacementStrategyFieldLength = 255
)

// AWS specification limit values shared by the validators, the engine, and
// the operation handlers. Each value is declared exactly once here and every
// other site references these names; the Smithy shape carrying the trait is
// noted per constant.
const (
	// RetryPolicy.MaximumEventAgeInSeconds @range(60, 86400).
	MinRetryPolicyEventAgeSeconds = 60
	MaxRetryPolicyEventAgeSeconds = 86400
	// RetryPolicy.MaximumRetryAttempts @range(0, 185).
	MaxRetryPolicyAttempts = 185
	// FlexibleTimeWindow.MaximumWindowInMinutes @range(1, 1440).
	MaxFlexibleWindowMinutes = 1440
	// TagList @length(0, 200): tags per schedule group.
	MaxTagsPerResource = 200
	// MaxResults @range(1, 100) on both list operations.
	DefaultListMaxResults = 100
	MaxListMaxResults     = 100
	// NextToken @length(1, 2048) on both list operations.
	maxNextTokenLength = 2048
	// KmsKeyArn @length(1, 2048).
	maxKmsKeyArnLength = 2048
	// TagResourceArn @length(1, 1011).
	maxTagResourceArnLength = 1011
	// Target.Input documented maximum: "The maximum size of the Input
	// field is 256 KB." (Target, AWS API reference).
	MaxTargetInputBytes = 256 * 1024
	// Target.Arn / Target.RoleArn / ResourceArn (DeadLetterConfig.Arn) /
	// TaskDefinitionArn @length(1, 1600).
	MaxTargetArnLength = 1600
	// TaskCount @range(1, 10).
	MaxEcsTaskCount = 10
	// CapacityProviderStrategy @length(max 6).
	MaxCapacityProviderStrategyItems = 6
	// PlacementConstraints @length(max 10).
	MaxPlacementConstraintItems = 10
	// PlacementStrategies @length(max 5).
	MaxPlacementStrategyItems = 5
	// EcsParameters Tags @length(0, 50): TagMap entries per ECS target.
	MaxEcsTagsItems = 50
	// CapacityProvider @length(1, 255).
	MinCapacityProviderLength   = 1
	MaxCapacityProviderLength   = 255
	MaxCapacityProviderWeight   = 1000   // CapacityProviderStrategyItemWeight @range(0, 1000)
	MaxCapacityProviderBase     = 100000 // CapacityProviderStrategyItemBase @range(0, 100000)
	MinSubnetsPerTask           = 1      // Subnets @length(1, 16)
	MaxSubnetsPerTask           = 16
	MinSecurityGroupsPerTask    = 1 // SecurityGroups @length(1, 5)
	MaxSecurityGroupsPerTask    = 5
	MaxSubnetIdLength           = 1000 // Subnet @length(1, 1000)
	MaxSecurityGroupIdLength    = 1000 // SecurityGroup @length(1, 1000)
	MaxMessageGroupIdLength     = 128  // MessageGroupId @length(1, 128)
	MaxDetailTypeLength         = 128  // DetailType @length(1, 128)
	MaxSourceLength             = 256  // Source @length(1, 256)
	MaxTargetPartitionKeyLength = 256  // TargetPartitionKey @length(1, 256)
	maxTagKeyLength             = 128  // TagKey @length(1, 128)
	maxTagValueLength           = 256  // TagValue @length(1, 256)
)

// namePattern matches the AWS Scheduler Name/GroupName constraint:
// 1-64 chars of alphanumeric, hyphen, underscore, and period.
var namePattern = regexp.MustCompile(`^[0-9a-zA-Z-_.]{1,64}$`)

// dateLayouts lists the timestamp formats accepted by the Scheduler API.
var dateLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	timeutils.ISO8601UTCFormat,
	timeutils.ISO8601NoZFormat,
	"2006-01-02",
}

// supportedTargetServices maps an ARN service segment to the delivery
// function that handles it in engine_targets.go's deliverToTarget switch.
// Templated targets pointing to services outside this set are rejected at
// validation time; ARNs of the universal-target form
// arn:aws:scheduler:::aws-sdk:{service}:{action} are validated separately
// (parseUniversalTargetARN / validateUniversalTarget) and dispatched by
// deliverUniversalTarget in engine_universal.go.
//
// The set of templated targets is defined by the AWS EventBridge Scheduler
// User Guide ("Using templated targets in EventBridge Scheduler"). The full
// AWS list is: CodeBuild, CodePipeline, ECS, EventBridge, Inspector, Kinesis,
// Firehose, Lambda, SageMaker AI, SNS, SQS, Step Functions. Note that SSM and
// AppSync are EventBridge Rules targets only; the Scheduler has no templated
// target for them (they are reachable only via universal targets).
//
// Currently supported (delivery implemented):
//
//	lambda, sqs, sns, kinesis, states (Step Functions), events (EventBridge)
//
// CloudWatch Logs is NOT a templated target in AWS and must not be
// accepted here.
//
// Accepted with stub delivery (accepted at validation, delivery fails until
// the backing service is implemented on this platform, mirroring the
// EventBridge accept-and-fail-at-delivery behaviour):
//
//	ecs — EcsParameters types/parsing/validation fully implemented. Delivery
//	      returns "not available" until the ECS service exists on this platform.
//	firehose — No sub-parameters (Smithy model has no FirehoseParameters).
//	           Delivery returns "not available" until the Firehose service
//	           exists on this platform.
//
// Out of scope (permanently unsupported on this edge/on-prem platform).
// Templated ARNs for these services are rejected at validation; universal
// aws-sdk ARNs naming them are accepted and fail at delivery with a cause
// naming the recorded exclusion (deliverUniversalTarget):
//
//	sagemaker — ML pipeline service (types stripped; the model's
//	            SageMakerPipelineParameters member is not carried anywhere)
//	codebuild — CI/CD build service
//	codepipeline — CI/CD pipeline orchestration
//	inspector — Security assessment service
var supportedTargetServices = map[string]bool{
	"lambda":   true,
	"sqs":      true,
	"sns":      true,
	"kinesis":  true,
	"states":   true,
	"events":   true,
	"ecs":      true,
	"firehose": true,
}

// universalTargetARNService is the ARN service segment of the universal-
// target form: arn:aws:scheduler:::aws-sdk:{service}:{action} ("Using
// universal targets in EventBridge Scheduler", AWS User Guide — the
// {service} value is the AWS SDK service identifier for the target
// service, which can differ from the endpoint prefix, e.g. sfn or
// eventbridge).
const universalTargetARNService = "scheduler"

// universalResourcePrefix starts the resource segment of a universal-target
// ARN: aws-sdk:{service}:{action}.
const universalResourcePrefix = "aws-sdk:"

// readOnlyActionPrefixes lists the API action prefixes AWS does not support
// as universal targets ("EventBridge Scheduler does not support read-only
// API actions, such as common GET operations, that begin with the following
// list of prefixes"). The documentation's own examples match the prefix
// case-insensitively — getQueueURL and ListBrokers are both rejected.
var readOnlyActionPrefixes = []string{
	"get", "describe", "list", "poll", "receive", "search", "scan", "query",
	"select", "read", "lookup", "discover", "validate", "batchget",
	"batchdescribe", "batchread", "transactget", "adminget", "adminlist",
	"testmigration", "retrieve", "testconnection", "translatedocument",
	"isauthorized", "invokemodel",
}

// parseUniversalTargetARN parses the universal-target ARN form
// arn:partition:scheduler:::aws-sdk:{service}:{action} into the target SDK
// service identifier and API action. ok is false for every other ARN,
// including a scheduler-service ARN whose resource is not of the aws-sdk
// form (validateTarget rejects that shape as a malformed universal ARN).
func parseUniversalTargetARN(arn string) (sdkService, action string, ok bool) {
	_, service, _, _, resource := svcarn.SplitARN(arn)
	if service != universalTargetARNService {
		return "", "", false
	}
	rest := strings.TrimPrefix(resource, universalResourcePrefix)
	if rest == resource {
		return "", "", false
	}
	parts := strings.SplitN(rest, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// validateUniversalTarget applies the creation-time contract of a universal
// target (arn:aws:scheduler:::aws-sdk:{service}:{action}):
//   - read-only API actions are not supported targets (the AWS universal-
//     target documentation's prefix list);
//   - Input, when present, must be well-formed JSON ("a well-formed JSON
//     you specify with the request parameters that EventBridge Scheduler
//     sends to the target API");
//   - the templated sub-parameter members belong to templated targets — a
//     universal target carries its request in Input, so any of them present
//     is rejected rather than carried inert.
//
// Which services and operations DELIVER, and the missing-substrate and
// recorded-exclusion failure causes for the rest, are delivery-time
// behaviour (deliverUniversalTarget): AWS accepts the ARN form for any SDK
// service and fails per target at invocation time, so validation does not
// gate on the named service.
func validateUniversalTarget(target *schedulerstore.Target, sdkService, action string) error {
	lower := strings.ToLower(action)
	for _, prefix := range readOnlyActionPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return awserrors.NewValidationException(fmt.Sprintf(
				"universal target action %q begins with the read-only prefix %q; read-only API actions are not supported as targets",
				action, prefix))
		}
	}
	if target.Input != "" {
		var probe interface{}
		if err := json.Unmarshal([]byte(target.Input), &probe); err != nil {
			return awserrors.NewValidationException(
				"Target.Input must be well-formed JSON carrying the request parameters of the universal target")
		}
	}
	if target.EcsParameters != nil || target.EventBridgeParameters != nil ||
		target.KinesisParameters != nil || target.SqsParameters != nil {
		return awserrors.NewValidationException(
			"templated target parameters cannot be specified on a universal target; pass the request parameters in Target.Input")
	}
	return nil
}

// validEcsLaunchTypes lists the Smithy enum values for EcsParameters.LaunchType.
var validEcsLaunchTypes = map[string]bool{
	"EC2":      true,
	"FARGATE":  true,
	"EXTERNAL": true,
}

// validPropagateTags lists the Smithy enum values for EcsParameters.PropagateTags.
var validPropagateTags = map[string]bool{
	"TASK_DEFINITION": true,
}

// validPlacementConstraintTypes lists the Smithy enum values for
// PlacementConstraint.type.
var validPlacementConstraintTypes = map[string]bool{
	"distinctInstance": true,
	"memberOf":         true,
}

// validPlacementStrategyTypes lists the Smithy enum values for
// PlacementStrategy.type.
var validPlacementStrategyTypes = map[string]bool{
	"random":  true,
	"spread":  true,
	"binpack": true,
}

// validAssignPublicIp lists the Smithy enum values for
// AwsVpcConfiguration.AssignPublicIp.
var validAssignPublicIp = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

// EventBridge Source field Smithy pattern decomposition.
// The original Smithy regex uses lookahead assertions that Go's RE2
// cannot express, so it is split into two checks:
//   - sourceFirstCharRe: the first character must be in the allowed set
//     [/ . - _ A-Za-z0-9] (replaces the positive lookahead)
//   - sourceJSONPathRe: the JSONPath alternative ($ followed by dot-
//     separated segments with optional array indices)
//
// The negative lookahead (?!aws\.) is implemented as a separate
// strings.HasPrefix check in validateEventBridgeParameters.
var (
	sourceFirstCharRe = regexp.MustCompile(`^[/.\-_A-Za-z0-9]`)
	sourceJSONPathRe  = regexp.MustCompile(`^\$(\.[\w_-]+(\[(\d+|\*)\])*)*$`)
)

// Target.RoleArn @pattern from the service model, applied verbatim: an IAM
// role ARN (service iam, role/ resource) in the aws partition family, with
// an empty region, a 12-digit account, and a non-empty role path/name
// limited to [\w+=,.@/-].
var targetRoleArnPattern = regexp.MustCompile(`^arn:aws(-[a-z]+)?:iam::\d{12}:role/[\w+=,.@/-]+$`)

// KmsKeyArn @pattern from the service model, applied verbatim: a KMS key or
// alias ARN in the aws partition family, with a region, a 12-digit account,
// and a key-id or alias-name portion over [0-9a-zA-Z-_].
var kmsKeyArnPattern = regexp.MustCompile(`^arn:aws(-[a-z]+)?:kms:[a-z0-9\-]+:\d{12}:(key|alias)\/[0-9a-zA-Z-_]*$`)

// DeadLetterConfig.Arn member @pattern from the service model, applied
// verbatim: an SQS queue ARN whose queue-name portion is [a-zA-Z0-9-_]+.
var deadLetterQueueArnPattern = regexp.MustCompile(`^arn:aws(-[a-z]+)?:sqs:[a-z0-9\-]+:\d{12}:[a-zA-Z0-9\-_]+$`)

// ScheduleSpec is the common input structure for schedule creation and
// update, used by both the HTTP API and the admin console to guarantee
// identical validation through the shared validation layer.
type ScheduleSpec struct {
	Name                       string
	GroupName                  string
	ScheduleExpression         string
	ScheduleExpressionTimezone string
	Description                string
	State                      string
	KmsKeyArn                  string
	StartDate                  string
	EndDate                    string
	ActionAfterCompletion      string
	Target                     *schedulerstore.Target
	FlexibleTimeWindow         *schedulerstore.FlexibleTimeWindow
}

// ValidatedSchedule holds the parsed and validated schedule fields ready
// for store persistence. Produced by validateScheduleFields so that both
// the HTTP API and admin paths build the store model identically.
type ValidatedSchedule struct {
	State                 schedulerstore.ScheduleState
	ActionAfterCompletion schedulerstore.ActionAfterCompletion
	StartDate             *time.Time
	EndDate               *time.Time
}

// validateClientToken validates the ClientToken format per Smithy spec:
// length [1, 64] in characters, pattern ^[a-zA-Z0-9-_]+$ (the pattern's
// ASCII charset subsumes the basis distinction; conforming tokens count
// identically in bytes and runes).
func validateClientToken(token string) error {
	if n := utf8.RuneCountInString(token); n < 1 || n > 64 {
		return awserrors.NewValidationException("ClientToken must be 1-64 characters")
	}
	for _, c := range token {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_') {
			return awserrors.NewValidationException("ClientToken contains invalid characters; allowed: alphanumeric, hyphen, underscore")
		}
	}
	return nil
}

// validateScheduleGroupName checks the ScheduleGroupName shape constraints
// (@length(1,64), @pattern ^[0-9a-zA-Z-_.]+$). The charset and bounds equal
// the Name shape's, so the single namePattern definition serves both.
func validateScheduleGroupName(name string) error {
	if !namePattern.MatchString(name) {
		return ErrValidation
	}
	return nil
}

// resolveListMaxResults applies MaxResults @range(1,100); a nil value means
// the member was absent and takes the default page size. The pointer form
// keeps an explicitly-passed zero distinguishable from an omitted member.
func resolveListMaxResults(maxResults *int32) (int32, error) {
	if maxResults == nil {
		return DefaultListMaxResults, nil
	}
	if *maxResults < 1 || *maxResults > MaxListMaxResults {
		return 0, ErrValidation
	}
	return *maxResults, nil
}

// validateListStateFilter checks the ListSchedules State filter enum
// (ENABLED | DISABLED); an empty value means no filtering.
func validateListStateFilter(state string) error {
	filter := schedulerstore.ScheduleState(state)
	if state != "" && filter != schedulerstore.ScheduleStateEnabled && filter != schedulerstore.ScheduleStateDisabled {
		return ErrValidation
	}
	return nil
}

// validateListNamePrefix checks the NamePrefix / ScheduleGroupNamePrefix
// shape constraints (@length(1,64), @pattern ^[0-9a-zA-Z-_.]+$); an empty
// value means no filtering.
func validateListNamePrefix(prefix string) error {
	if prefix != "" && !namePattern.MatchString(prefix) {
		return ErrValidation
	}
	return nil
}

// validateNextToken checks the NextToken @length(1,2048) bound; an empty
// value simply starts a new traversal.
func validateNextToken(token string) error {
	if utf8.RuneCountInString(token) > maxNextTokenLength {
		return ErrValidation
	}
	return nil
}

// validateDateFlexible parses a date string using multiple AWS-accepted
// formats, returning the parsed time and an error.
func validateDateFlexible(s string) (time.Time, error) {
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}

// validateScheduleFields validates all schedule fields per AWS Smithy spec.
// Called by both HTTP API and admin console paths (shared validation).
// Returns the normalised schedule fields or an error.
func validateScheduleFields(spec *ScheduleSpec) (*ValidatedSchedule, error) {
	if spec.Name == "" || !namePattern.MatchString(spec.Name) {
		return nil, ErrValidation
	}

	if spec.ScheduleExpression == "" {
		return nil, ErrValidation
	}
	if !isValidScheduleExpression(spec.ScheduleExpression) {
		return nil, ErrInvalidScheduleExpression
	}

	if spec.Target == nil {
		return nil, ErrInvalidTarget
	}
	if err := validateTarget(spec.Target); err != nil {
		return nil, err
	}

	// FlexibleTimeWindow is a required member of CreateScheduleInput and
	// UpdateScheduleInput, and Mode is a required member of the
	// FlexibleTimeWindow shape; both absences are rejected here.
	if spec.FlexibleTimeWindow == nil {
		return nil, awserrors.NewValidationException("FlexibleTimeWindow is required")
	}
	if err := validateFlexibleTimeWindow(spec.FlexibleTimeWindow); err != nil {
		return nil, err
	}

	result := &ValidatedSchedule{
		State:                 schedulerstore.ScheduleStateEnabled,
		ActionAfterCompletion: schedulerstore.ActionAfterCompletionNone,
	}
	if spec.State != "" {
		state := schedulerstore.ScheduleState(spec.State)
		if state != schedulerstore.ScheduleStateEnabled && state != schedulerstore.ScheduleStateDisabled {
			return nil, ErrInvalidState
		}
		result.State = state
	}

	if spec.ActionAfterCompletion != "" {
		action := schedulerstore.ActionAfterCompletion(spec.ActionAfterCompletion)
		if action != schedulerstore.ActionAfterCompletionNone && action != schedulerstore.ActionAfterCompletionDelete {
			return nil, ErrInvalidActionAfterCompletion
		}
		result.ActionAfterCompletion = action
	}

	if spec.KmsKeyArn != "" {
		// The KmsKeyArn @pattern, applied verbatim (see kmsKeyArnPattern),
		// governs the ARN grammar: the ARN must name a KMS key or alias.
		// Whether the key exists and is a symmetric encryption key is a
		// separate creation-time check (validateKmsKey).
		if !kmsKeyArnPattern.MatchString(spec.KmsKeyArn) {
			return nil, awserrors.NewValidationException(fmt.Sprintf(
				"KmsKeyArn %q must be the ARN of a KMS key or alias (arn:aws:kms:region:account:key/... or alias/...)",
				spec.KmsKeyArn))
		}
		if utf8.RuneCountInString(spec.KmsKeyArn) > maxKmsKeyArnLength {
			return nil, awserrors.NewValidationException(fmt.Sprintf(
				"KmsKeyArn must be at most %d characters", maxKmsKeyArnLength))
		}
	}

	if spec.ScheduleExpressionTimezone != "" {
		// ScheduleExpressionTimezone @length(1,50) counts Unicode characters
		// (no pattern); the non-empty branch guarantees the minimum.
		if n := utf8.RuneCountInString(spec.ScheduleExpressionTimezone); n > maxTimezoneLength {
			return nil, awserrors.NewValidationException("ScheduleExpressionTimezone must be 1-50 characters")
		}
		// Go's time.LoadLocation resolves "Local" to the host timezone; the
		// Scheduler contract takes IANA time zone database names, and
		// "Local" is a Go runtime alias rather than a database entry.
		if spec.ScheduleExpressionTimezone == "Local" {
			return nil, awserrors.NewValidationException(
				"ScheduleExpressionTimezone must name an IANA time zone; \"Local\" is not an IANA zone name")
		}
		if _, err := time.LoadLocation(spec.ScheduleExpressionTimezone); err != nil {
			return nil, awserrors.NewValidationException("invalid ScheduleExpressionTimezone: not a valid IANA timezone")
		}
	}

	// Description @length(0,512) counts Unicode characters (no pattern).
	if utf8.RuneCountInString(spec.Description) > maxScheduleDescriptionLen {
		return nil, awserrors.NewValidationException("Description must be 0-512 characters")
	}

	if spec.StartDate != "" {
		t, err := validateDateFlexible(spec.StartDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
		result.StartDate = &t
	}
	if spec.EndDate != "" {
		t, err := validateDateFlexible(spec.EndDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
		result.EndDate = &t
	}
	if result.StartDate != nil && result.EndDate != nil && result.StartDate.After(*result.EndDate) {
		return nil, ErrInvalidDate
	}

	return result, nil
}

// validateTarget validates the Target structure comprehensively:
//   - ARN format for Target, RoleArn, and DeadLetterConfig
//   - Target ARN names a supported templated delivery type, or carries the
//     universal-target aws-sdk form (validateUniversalTarget)
//   - DeadLetterConfig ARN must be SQS only
//   - RetryPolicy ranges (Smithy)
//   - Sub-parameter / service cross-check (sub-parameters must match
//     the target ARN service — e.g. EcsParameters only on ECS targets)
//   - Sub-parameter detailed validation against Smithy ranges and patterns
func validateTarget(target *schedulerstore.Target) error {
	if target.Arn == "" {
		return ErrInvalidTarget
	}
	parsedArn, err := svcarn.ParseARN(target.Arn)
	if err != nil {
		return ErrInvalidTarget
	}
	// Target.Arn @length(1, 1600) in characters; the minimum is covered by
	// the empty check above.
	if utf8.RuneCountInString(target.Arn) > MaxTargetArnLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"Target.Arn must be at most %d characters", MaxTargetArnLength))
	}

	// A universal target carries the aws-sdk ARN form and its request in
	// Input; every other ARN must name a templated target service.
	sdkService, action, universal := parseUniversalTargetARN(target.Arn)
	if universal {
		if err := validateUniversalTarget(target, sdkService, action); err != nil {
			return err
		}
	} else if parsedArn.Service == universalTargetARNService {
		// A scheduler-service ARN that is not of the aws-sdk form names no
		// templated service and no universal operation.
		return awserrors.NewValidationException(fmt.Sprintf(
			"a scheduler-service target ARN must use the universal-target form arn:aws:scheduler:::aws-sdk:{service}:{action}, got %s",
			target.Arn))
	} else if err := validateTargetService(parsedArn.Service); err != nil {
		// Reject templated targets pointing to services we cannot deliver to.
		return err
	}

	// A templated target ARN names the region its resource lives in: the
	// modelled ARN grammar carries a non-empty region segment, and the
	// deliverers resolve the region from the ARN alone. The universal-target
	// form above is the one modelled exception — its grammar fixes an empty
	// region by design — and an empty-region ARN on any other service can
	// never resolve its resource.
	if !universal && parsedArn.Region == "" {
		return awserrors.NewValidationException(fmt.Sprintf(
			"Target.Arn %s must carry the target resource's region", target.Arn))
	}

	if target.RoleArn == "" {
		return ErrInvalidTarget
	}
	// RoleArn must satisfy the RoleArn @pattern (see targetRoleArnPattern):
	// a full IAM role ARN, not merely an iam-service ARN with a role/
	// resource prefix.
	if !targetRoleArnPattern.MatchString(target.RoleArn) {
		return ErrInvalidTarget
	}
	// RoleArn @length(1, 1600) in characters, shared with Target.Arn.
	if utf8.RuneCountInString(target.RoleArn) > MaxTargetArnLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"Target.RoleArn must be at most %d characters", MaxTargetArnLength))
	}

	// Target.Input carries a documented 256 KB maximum; the shape's
	// @length(min 1) is not enforceable here because the DTO string cannot
	// distinguish an omitted member from an explicitly empty one.
	if len(target.Input) > MaxTargetInputBytes {
		return awserrors.NewValidationException(fmt.Sprintf(
			"Target.Input must be at most %d bytes", MaxTargetInputBytes))
	}

	// "If you are configuring a templated Lambda, AWS Step Functions, or
	// Amazon EventBridge target, the input must be a well-formed JSON"
	// (Target.Input, model documentation) — any JSON value qualifies; the
	// object form is a delivery-time concern of the events family, whose
	// PutEvents pipeline enforces it, not a creation-time restriction. All
	// other target types carry free text ("For all other target types, a
	// JSON is not required").
	if isJSONBoundTargetService(parsedArn.Service) && target.Input != "" {
		if !json.Valid([]byte(target.Input)) {
			return awserrors.NewValidationException(
				"Target.Input must be a well-formed JSON for templated Lambda, Step Functions and EventBridge targets")
		}
	}

	// DeadLetterConfig ARN must be SQS only (AWS spec).
	if target.DeadLetterConfig != nil && target.DeadLetterConfig.Arn != "" {
		if err := validateDLQService(target.DeadLetterConfig.Arn); err != nil {
			return err
		}
	}

	// RetryPolicy ranges.
	if target.RetryPolicy != nil {
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil {
			v := *target.RetryPolicy.MaximumEventAgeInSeconds
			if v < MinRetryPolicyEventAgeSeconds || v > MaxRetryPolicyEventAgeSeconds {
				return awserrors.NewValidationException(fmt.Sprintf(
					"RetryPolicy.MaximumEventAgeInSeconds must be between %d and %d",
					MinRetryPolicyEventAgeSeconds, MaxRetryPolicyEventAgeSeconds))
			}
		}
		if target.RetryPolicy.MaximumRetryAttempts != nil {
			v := *target.RetryPolicy.MaximumRetryAttempts
			if v < 0 || v > MaxRetryPolicyAttempts {
				return awserrors.NewValidationException(fmt.Sprintf(
					"RetryPolicy.MaximumRetryAttempts must be between 0 and %d",
					MaxRetryPolicyAttempts))
			}
		}
	}

	// Cross-check: service-specific sub-parameters must match the target
	// ARN service. AWS rejects e.g. EcsParameters on a Lambda target.
	if err := validateSubParametersForService(parsedArn.Service, target); err != nil {
		return err
	}

	// Sub-parameter detailed validation (per Smithy traits).
	if target.EcsParameters != nil {
		if err := validateEcsParameters(target.EcsParameters); err != nil {
			return err
		}
	}
	if target.EventBridgeParameters != nil {
		if err := validateEventBridgeParameters(target.EventBridgeParameters); err != nil {
			return err
		}
	}
	if target.KinesisParameters != nil {
		if n := utf8.RuneCountInString(target.KinesisParameters.PartitionKey); n < 1 || n > MaxTargetPartitionKeyLength {
			return awserrors.NewValidationException(fmt.Sprintf(
				"KinesisParameters.PartitionKey must be 1-%d characters", MaxTargetPartitionKeyLength))
		}
	}
	if target.SqsParameters != nil {
		if err := validateSqsParameters(target.SqsParameters); err != nil {
			return err
		}
	}

	return nil
}

// validateTargetService rejects target ARNs whose service segment is not
// in the supported set. This prevents schedules being created for
// delivery types that have no implementation.
func validateTargetService(service string) error {
	if !supportedTargetServices[service] {
		return awserrors.NewValidationException(
			fmt.Sprintf("unsupported target service %q; supported services: lambda, sqs, sns, kinesis, states, events, ecs, firehose", service),
		)
	}
	return nil
}

// validateDLQService enforces the AWS specification that the
// DeadLetterConfig ARN reference an SQS queue: the member @pattern, applied
// verbatim (see deadLetterQueueArnPattern), and the ResourceArn
// @length(1, 1600) bound.
func validateDLQService(arn string) error {
	// DeadLetterConfig.Arn @length(1, 1600) in characters.
	if utf8.RuneCountInString(arn) > MaxTargetArnLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"DeadLetterConfig.Arn must be at most %d characters", MaxTargetArnLength))
	}
	if !deadLetterQueueArnPattern.MatchString(arn) {
		return awserrors.NewValidationException(fmt.Sprintf(
			"DeadLetterConfig.Arn %q must be the ARN of an SQS queue (arn:aws:sqs:region:account:queue-name)",
			arn))
	}
	return nil
}

// isJSONBoundTargetService reports whether the templated target family's
// Input is bound to well-formed JSON by the model's Target.Input
// documentation: templated Lambda, Step Functions (states) and EventBridge.
// An omitted Input stays legal for every family — the delivery path
// substitutes the default notification payload.
func isJSONBoundTargetService(service string) bool {
	return service == "lambda" || service == "states" || service == "events"
}

// validateSubParametersForService enforces the AWS constraint that
// service-specific sub-parameters on a Target may only be specified when
// the target ARN's service matches. For example, EcsParameters is only
// valid on ECS targets and KinesisParameters only on Kinesis targets.
// EventBridgeParameters stays optional on events targets (the model marks
// no sub-parameter required): a schedule without them is created and its
// delivery fails with the PutEvents cause the eventbridge handler reports
// for missing Source and DetailType — the accept-and-fail posture the
// platform's recorded substrate rules already apply.
// DeadLetterConfig and RetryPolicy are universal and exempt.
func validateSubParametersForService(service string, target *schedulerstore.Target) error {
	if target.EcsParameters != nil && service != "ecs" {
		return awserrors.NewValidationException(
			"EcsParameters can only be specified for ECS targets")
	}
	if target.EventBridgeParameters != nil && service != "events" {
		return awserrors.NewValidationException(
			"EventBridgeParameters can only be specified for EventBridge targets")
	}
	if target.KinesisParameters != nil && service != "kinesis" {
		return awserrors.NewValidationException(
			"KinesisParameters can only be specified for Kinesis targets")
	}
	if target.SqsParameters != nil && service != "sqs" {
		return awserrors.NewValidationException(
			"SqsParameters can only be specified for SQS targets")
	}
	return nil
}

// validateSqsParameters validates SqsParameters per Smithy traits.
// MessageGroupId carries no @required in the model — {"sqsParameters": {}}
// is valid input, and the delivery path falls back to the schedule name
// for FIFO queues when it is unset — and the DTO string cannot distinguish
// an omitted member from an explicitly empty one, so the bound applies
// only when a value is present. When present: length [1, 128] in
// characters.
func validateSqsParameters(sqs *schedulerstore.SqsParameters) error {
	if sqs.MessageGroupId == "" {
		return nil
	}
	if n := utf8.RuneCountInString(sqs.MessageGroupId); n > MaxMessageGroupIdLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"SqsParameters.MessageGroupId must be 1-%d characters", MaxMessageGroupIdLength))
	}
	return nil
}

// validateEcsParameters validates all EcsParameters fields per Smithy
// traits and AWS documentation.
func validateEcsParameters(ecs *schedulerstore.EcsParameters) error {
	if ecs.TaskDefinitionArn == "" {
		return awserrors.NewValidationException("EcsParameters.TaskDefinitionArn is required")
	}
	if _, err := svcarn.ParseARN(ecs.TaskDefinitionArn); err != nil {
		return awserrors.NewValidationException("EcsParameters.TaskDefinitionArn must be a valid ARN")
	}
	// TaskDefinitionArn @length(1, 1600) in characters, shared with
	// Target.Arn and RoleArn; the minimum is covered by the empty check
	// above.
	if utf8.RuneCountInString(ecs.TaskDefinitionArn) > MaxTargetArnLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EcsParameters.TaskDefinitionArn must be at most %d characters", MaxTargetArnLength))
	}
	if ecs.TaskCount != nil {
		v := *ecs.TaskCount
		if v < 1 || v > MaxEcsTaskCount {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.TaskCount must be between 1 and %d", MaxEcsTaskCount))
		}
	}
	if ecs.LaunchType != "" && !validEcsLaunchTypes[ecs.LaunchType] {
		return awserrors.NewValidationException(
			fmt.Sprintf("EcsParameters.LaunchType must be one of EC2, FARGATE, EXTERNAL; got %q", ecs.LaunchType),
		)
	}
	if len(ecs.CapacityProviderStrategy) > MaxCapacityProviderStrategyItems {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EcsParameters.CapacityProviderStrategy must have at most %d items", MaxCapacityProviderStrategyItems))
	}
	if len(ecs.PlacementConstraints) > MaxPlacementConstraintItems {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EcsParameters.PlacementConstraints must have at most %d items", MaxPlacementConstraintItems))
	}
	if len(ecs.PlacementStrategy) > MaxPlacementStrategyItems {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EcsParameters.PlacementStrategy must have at most %d items", MaxPlacementStrategyItems))
	}
	// PlatformVersion / Group / ReferenceId carry pattern-less @length
	// traits, so lengths count Unicode characters. The checks are maxima:
	// an empty string means the member was omitted.
	if utf8.RuneCountInString(ecs.PlatformVersion) > maxPlatformVersionLength {
		return awserrors.NewValidationException("EcsParameters.PlatformVersion must be at most 64 characters")
	}
	if utf8.RuneCountInString(ecs.Group) > maxEcsGroupLength {
		return awserrors.NewValidationException("EcsParameters.Group must be at most 255 characters")
	}
	if utf8.RuneCountInString(ecs.ReferenceId) > maxReferenceIdLength {
		return awserrors.NewValidationException("EcsParameters.ReferenceId must be at most 1024 characters")
	}
	if ecs.PropagateTags != "" && !validPropagateTags[ecs.PropagateTags] {
		return awserrors.NewValidationException(
			fmt.Sprintf("EcsParameters.PropagateTags must be TASK_DEFINITION; got %q", ecs.PropagateTags),
		)
	}
	// EcsParameters.Tags is a list of TagMap (map<TagKey, TagValue>): at
	// most fifty entries, and every pair carries a key of 1-128 and a
	// value of 1-256 characters. The model places no single-pair bound
	// on a TagMap, so one entry may carry several pairs and each pair
	// is validated; the minimum of 1 makes an empty key or value invalid.
	if len(ecs.Tags) > MaxEcsTagsItems {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EcsParameters.Tags must have at most %d items", MaxEcsTagsItems))
	}
	for _, tagMap := range ecs.Tags {
		for key, value := range tagMap {
			if n := utf8.RuneCountInString(key); n < 1 || n > maxTagKeyLength {
				return awserrors.NewValidationException(fmt.Sprintf(
					"EcsParameters.Tags keys must be 1-%d characters", maxTagKeyLength))
			}
			if n := utf8.RuneCountInString(value); n < 1 || n > maxTagValueLength {
				return awserrors.NewValidationException(fmt.Sprintf(
					"EcsParameters.Tags values must be 1-%d characters", maxTagValueLength))
			}
		}
	}
	// CapacityProviderStrategyItem: capacityProvider is required with
	// length 1-255; weight and base are bounded ranges.
	for i, item := range ecs.CapacityProviderStrategy {
		if n := utf8.RuneCountInString(item.CapacityProvider); n < MinCapacityProviderLength || n > MaxCapacityProviderLength {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.CapacityProviderStrategy[%d].capacityProvider must be %d-%d characters",
				i, MinCapacityProviderLength, MaxCapacityProviderLength))
		}
		if item.Weight != nil && (*item.Weight < 0 || *item.Weight > MaxCapacityProviderWeight) {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.CapacityProviderStrategy[%d].weight must be between 0 and %d",
				i, MaxCapacityProviderWeight))
		}
		if item.Base != nil && (*item.Base < 0 || *item.Base > MaxCapacityProviderBase) {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.CapacityProviderStrategy[%d].base must be between 0 and %d",
				i, MaxCapacityProviderBase))
		}
	}
	for i, pc := range ecs.PlacementConstraints {
		if pc.Type != "" && !validPlacementConstraintTypes[pc.Type] {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.PlacementConstraints[%d].type must be one of distinctInstance, memberOf; got %q",
				i, pc.Type))
		}
		// PlacementConstraintExpression @length(max 2000), counted in
		// Unicode characters (pattern-less shape).
		if utf8.RuneCountInString(pc.Expression) > maxPlacementExpressionLength {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.PlacementConstraints[%d].expression must be at most %d characters",
				i, maxPlacementExpressionLength))
		}
	}
	for i, ps := range ecs.PlacementStrategy {
		if ps.Type != "" && !validPlacementStrategyTypes[ps.Type] {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.PlacementStrategy[%d].type must be one of random, spread, binpack; got %q",
				i, ps.Type))
		}
		// PlacementStrategyField @length(max 255), counted in Unicode
		// characters (pattern-less shape).
		if utf8.RuneCountInString(ps.Field) > maxPlacementStrategyFieldLength {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.PlacementStrategy[%d].field must be at most %d characters",
				i, maxPlacementStrategyFieldLength))
		}
	}
	if ecs.NetworkConfiguration != nil && ecs.NetworkConfiguration.AwsVpcConfiguration != nil {
		vpc := ecs.NetworkConfiguration.AwsVpcConfiguration
		if l := len(vpc.Subnets); l < MinSubnetsPerTask || l > MaxSubnetsPerTask {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.NetworkConfiguration.awsvpcConfiguration.subnets must contain %d-%d subnets, got %d",
				MinSubnetsPerTask, MaxSubnetsPerTask, l))
		}
		for i, subnet := range vpc.Subnets {
			if n := utf8.RuneCountInString(subnet); n < 1 || n > MaxSubnetIdLength {
				return awserrors.NewValidationException(fmt.Sprintf(
					"EcsParameters.NetworkConfiguration.awsvpcConfiguration.subnets[%d] must be 1-%d characters",
					i, MaxSubnetIdLength))
			}
		}
		// A provided securityGroups list carries at least one entry; the
		// parse layer keeps an explicitly empty list distinguishable from
		// an omitted member (the latter means the VPC default group).
		if vpc.SecurityGroups != nil && len(vpc.SecurityGroups) < MinSecurityGroupsPerTask {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.NetworkConfiguration.awsvpcConfiguration.securityGroups must contain at least %d item when provided",
				MinSecurityGroupsPerTask))
		}
		if len(vpc.SecurityGroups) > MaxSecurityGroupsPerTask {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.NetworkConfiguration.awsvpcConfiguration.securityGroups must contain at most %d items",
				MaxSecurityGroupsPerTask))
		}
		for i, sg := range vpc.SecurityGroups {
			if n := utf8.RuneCountInString(sg); n < 1 || n > MaxSecurityGroupIdLength {
				return awserrors.NewValidationException(fmt.Sprintf(
					"EcsParameters.NetworkConfiguration.awsvpcConfiguration.securityGroups[%d] must be 1-%d characters",
					i, MaxSecurityGroupIdLength))
			}
		}
		if vpc.AssignPublicIp != "" && !validAssignPublicIp[vpc.AssignPublicIp] {
			return awserrors.NewValidationException(fmt.Sprintf(
				"EcsParameters.NetworkConfiguration.awsvpcConfiguration.assignPublicIp must be ENABLED or DISABLED; got %q",
				vpc.AssignPublicIp))
		}
	}
	return nil
}

// validateEventBridgeParameters validates EventBridgeParameters per
// Smithy traits and AWS documentation. The Source field is checked
// against the full Smithy pattern (decomposed for RE2 compatibility).
func validateEventBridgeParameters(eb *schedulerstore.EventBridgeParameters) error {
	if n := utf8.RuneCountInString(eb.DetailType); n < 1 || n > MaxDetailTypeLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EventBridgeParameters.DetailType must be 1-%d characters", MaxDetailTypeLength))
	}
	// Source @length(1, 256) in characters: the pattern's free-form
	// alternative restricts only the first character, leaving the tail
	// free to carry multibyte text, so the byte count cannot stand in for
	// the character count here.
	if n := utf8.RuneCountInString(eb.Source); n < 1 || n > MaxSourceLength {
		return awserrors.NewValidationException(fmt.Sprintf(
			"EventBridgeParameters.Source must be 1-%d characters", MaxSourceLength))
	}
	// JSONPath alternative (e.g. "$.detail.event").
	if sourceJSONPathRe.MatchString(eb.Source) {
		return nil
	}
	// Free-form alternative: first char must be in the allowed charset.
	if !sourceFirstCharRe.MatchString(eb.Source) {
		return awserrors.NewValidationException(
			fmt.Sprintf("EventBridgeParameters.Source %q does not match the required pattern", eb.Source))
	}
	// Must not use the reserved "aws." prefix (negative lookahead equivalent).
	if strings.HasPrefix(eb.Source, "aws.") {
		return awserrors.NewValidationException(
			"EventBridgeParameters.Source must not use the reserved 'aws.' prefix")
	}
	return nil
}

// ValidateScheduleGroupTags enforces the tag-set constraints shared by the
// TagResource path and CreateScheduleGroup: at most MaxTagsPerResource tags
// (TagList @length(0, 200)), each key 1-128 characters (TagKey @length) and
// each value 1-256 characters (TagValue @length).
func ValidateScheduleGroupTags(tags []tagutil.Tag) error {
	if len(tags) > MaxTagsPerResource {
		return awserrors.NewValidationException(fmt.Sprintf(
			"Tags must contain at most %d items", MaxTagsPerResource))
	}
	for _, t := range tags {
		if n := utf8.RuneCountInString(t.Key); n < 1 || n > maxTagKeyLength {
			return awserrors.NewValidationException(fmt.Sprintf(
				"tag keys must be 1-%d characters", maxTagKeyLength))
		}
		if n := utf8.RuneCountInString(t.Value); n < 1 || n > maxTagValueLength {
			return awserrors.NewValidationException(fmt.Sprintf(
				"tag values must be 1-%d characters", maxTagValueLength))
		}
	}
	return nil
}

// validateFlexibleTimeWindow validates the FlexibleTimeWindow Mode enum and
// the MaximumWindowInMinutes range — the @range(1, 1440) trait sits on the
// shape itself, so a provided window bound is range-checked in either mode,
// while FLEXIBLE mode additionally requires the member to be present.
func validateFlexibleTimeWindow(ftw *schedulerstore.FlexibleTimeWindow) error {
	if ftw.Mode == "" {
		return awserrors.NewValidationException("FlexibleTimeWindow.Mode is required")
	}
	if ftw.Mode != schedulerstore.FlexibleTimeWindowModeOff && ftw.Mode != schedulerstore.FlexibleTimeWindowModeFlexible {
		return ErrInvalidFlexibleTimeWindow
	}
	if ftw.MaximumWindowInMinutes != nil {
		if *ftw.MaximumWindowInMinutes < 1 || *ftw.MaximumWindowInMinutes > MaxFlexibleWindowMinutes {
			return ErrInvalidFlexibleTimeWindow
		}
	}
	if ftw.Mode == schedulerstore.FlexibleTimeWindowModeFlexible && ftw.MaximumWindowInMinutes == nil {
		return ErrInvalidFlexibleTimeWindow
	}
	return nil
}

// isValidScheduleExpression delegates to the shared schedule expression
// validator in internal/common/scheduleexpr.
func isValidScheduleExpression(expr string) bool {
	return scheduleexpr.ValidateExpression(expr)
}

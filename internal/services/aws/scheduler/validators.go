package scheduler

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/scheduleexpr"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// Smithy @length maxima for pattern-less string shapes, counted in Unicode
// characters.
const (
	maxTimezoneLength         = 50  // ScheduleExpressionTimezone
	maxScheduleDescriptionLen = 512 // Description
	maxPlatformVersionLength  = 64  // PlatformVersion
	maxEcsGroupLength         = 255 // Group
	maxReferenceIdLength      = 1024
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
// function that handles it in engine.go's deliverToTarget switch. Targets
// pointing to services outside this set are rejected at validation time.
//
// The set of templated targets is defined by the AWS EventBridge Scheduler
// User Guide ("Using templated targets in EventBridge Scheduler"). The full
// AWS list is: CodeBuild, CodePipeline, ECS, EventBridge, Inspector, Kinesis,
// Firehose, Lambda, SageMaker AI, SNS, SQS, Step Functions. Note that SSM and
// AppSync are EventBridge Rules targets only; the Scheduler has no templated
// target for them (they are reachable only via universal targets, which this
// platform does not implement).
//
// Currently supported (delivery implemented):
//
//	lambda, sqs, sns, kinesis, states (Step Functions), events (EventBridge), logs (CloudWatch Logs)
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
// Out of scope (permanently unsupported on this edge/on-prem platform):
//
//	sagemaker — ML pipeline service (types stripped)
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
	"logs":     true,
	"ecs":      true,
	"firehose": true,
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
// length [1, 64], pattern ^[a-zA-Z0-9-_]+$.
func validateClientToken(token string) error {
	if len(token) < 1 || len(token) > 64 {
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

	if spec.FlexibleTimeWindow != nil {
		if err := validateFlexibleTimeWindow(spec.FlexibleTimeWindow); err != nil {
			return nil, err
		}
	}

	result := &ValidatedSchedule{
		State:                 schedulerstore.ScheduleStateEnabled,
		ActionAfterCompletion: schedulerstore.ActionAfterCompletionNone,
	}
	if spec.State != "" {
		if spec.State != "ENABLED" && spec.State != "DISABLED" {
			return nil, ErrInvalidState
		}
		result.State = schedulerstore.ScheduleState(spec.State)
	}

	if spec.ActionAfterCompletion != "" {
		if spec.ActionAfterCompletion != "NONE" && spec.ActionAfterCompletion != "DELETE" {
			return nil, ErrInvalidActionAfterCompletion
		}
		result.ActionAfterCompletion = schedulerstore.ActionAfterCompletion(spec.ActionAfterCompletion)
	}

	if spec.KmsKeyArn != "" {
		parsed, err := svcarn.ParseARN(spec.KmsKeyArn)
		if err != nil {
			return nil, awserrors.NewValidationException("invalid KmsKeyArn ARN format")
		}
		if parsed.Service != "kms" {
			return nil, awserrors.NewValidationException("KmsKeyArn must reference a KMS key")
		}
	}

	if spec.ScheduleExpressionTimezone != "" {
		// ScheduleExpressionTimezone @length(1,50) counts Unicode characters
		// (no pattern); the non-empty branch guarantees the minimum.
		if n := utf8.RuneCountInString(spec.ScheduleExpressionTimezone); n > maxTimezoneLength {
			return nil, awserrors.NewValidationException("ScheduleExpressionTimezone must be 1-50 characters")
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
//   - Target ARN service must be a supported delivery type
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

	// Reject targets pointing to services we cannot deliver to.
	if err := validateTargetService(parsedArn.Service); err != nil {
		return err
	}

	if target.RoleArn == "" {
		return ErrInvalidTarget
	}
	if _, err := svcarn.ParseARN(target.RoleArn); err != nil {
		return ErrInvalidTarget
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
			if v < 60 || v > 86400 {
				return awserrors.NewValidationException("RetryPolicy.MaximumEventAgeInSeconds must be between 60 and 86400")
			}
		}
		if target.RetryPolicy.MaximumRetryAttempts != nil {
			v := *target.RetryPolicy.MaximumRetryAttempts
			if v < 0 || v > 185 {
				return awserrors.NewValidationException("RetryPolicy.MaximumRetryAttempts must be between 0 and 185")
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
		if l := len(target.KinesisParameters.PartitionKey); l < 1 || l > 256 {
			return awserrors.NewValidationException("KinesisParameters.PartitionKey must be 1-256 characters")
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
			fmt.Sprintf("unsupported target service %q; supported services: lambda, sqs, sns, kinesis, states, events, logs, ecs, firehose", service),
		)
	}
	return nil
}

// validateDLQService enforces the AWS specification that DeadLetterConfig
// ARN must reference an SQS queue.
func validateDLQService(arn string) error {
	_, service, _, _, _ := svcarn.SplitARN(arn)
	if service == "" {
		return ErrInvalidTarget
	}
	if service != "sqs" {
		return awserrors.NewValidationException(
			fmt.Sprintf("DeadLetterConfig ARN must reference an SQS queue, got service %q", service),
		)
	}
	return nil
}

// validateSubParametersForService enforces the AWS constraint that
// service-specific sub-parameters on a Target may only be specified when
// the target ARN's service matches. For example, EcsParameters is only
// valid on ECS targets and KinesisParameters only on Kinesis targets.
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
// MessageGroupId: length [1, 128].
func validateSqsParameters(sqs *schedulerstore.SqsParameters) error {
	if l := len(sqs.MessageGroupId); l < 1 || l > 128 {
		return awserrors.NewValidationException(
			"SqsParameters.MessageGroupId must be 1-128 characters")
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
	if ecs.TaskCount != nil {
		v := *ecs.TaskCount
		if v < 1 || v > 10 {
			return awserrors.NewValidationException("EcsParameters.TaskCount must be between 1 and 10")
		}
	}
	if ecs.LaunchType != "" && !validEcsLaunchTypes[ecs.LaunchType] {
		return awserrors.NewValidationException(
			fmt.Sprintf("EcsParameters.LaunchType must be one of EC2, FARGATE, EXTERNAL; got %q", ecs.LaunchType),
		)
	}
	if len(ecs.CapacityProviderStrategy) > 6 {
		return awserrors.NewValidationException("EcsParameters.CapacityProviderStrategy must have at most 6 items")
	}
	if len(ecs.PlacementConstraints) > 10 {
		return awserrors.NewValidationException("EcsParameters.PlacementConstraints must have at most 10 items")
	}
	if len(ecs.PlacementStrategy) > 5 {
		return awserrors.NewValidationException("EcsParameters.PlacementStrategy must have at most 5 items")
	}
	// PlatformVersion / Group / ReferenceId carry pattern-less @length
	// traits, so lengths count Unicode characters.
	if utf8.RuneCountInString(ecs.PlatformVersion) > maxPlatformVersionLength {
		return awserrors.NewValidationException("EcsParameters.PlatformVersion must be 1-64 characters")
	}
	if utf8.RuneCountInString(ecs.Group) > maxEcsGroupLength {
		return awserrors.NewValidationException("EcsParameters.Group must be 1-255 characters")
	}
	if utf8.RuneCountInString(ecs.ReferenceId) > maxReferenceIdLength {
		return awserrors.NewValidationException("EcsParameters.ReferenceId must be at most 1024 characters")
	}
	if ecs.PropagateTags != "" && !validPropagateTags[ecs.PropagateTags] {
		return awserrors.NewValidationException(
			fmt.Sprintf("EcsParameters.PropagateTags must be TASK_DEFINITION; got %q", ecs.PropagateTags),
		)
	}
	return nil
}

// validateEventBridgeParameters validates EventBridgeParameters per
// Smithy traits and AWS documentation. The Source field is checked
// against the full Smithy pattern (decomposed for RE2 compatibility).
func validateEventBridgeParameters(eb *schedulerstore.EventBridgeParameters) error {
	if l := len(eb.DetailType); l < 1 || l > 128 {
		return awserrors.NewValidationException("EventBridgeParameters.DetailType must be 1-128 characters")
	}
	if l := len(eb.Source); l < 1 || l > 256 {
		return awserrors.NewValidationException("EventBridgeParameters.Source must be 1-256 characters")
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

// validateFlexibleTimeWindow validates the FlexibleTimeWindow Mode enum
// and the MaximumWindowInMinutes range when Mode is FLEXIBLE.
func validateFlexibleTimeWindow(ftw *schedulerstore.FlexibleTimeWindow) error {
	if ftw.Mode != schedulerstore.FlexibleTimeWindowModeOff && ftw.Mode != schedulerstore.FlexibleTimeWindowModeFlexible {
		return ErrInvalidFlexibleTimeWindow
	}
	if ftw.Mode == schedulerstore.FlexibleTimeWindowModeFlexible {
		if ftw.MaximumWindowInMinutes == nil || *ftw.MaximumWindowInMinutes < 1 || *ftw.MaximumWindowInMinutes > 1440 {
			return ErrInvalidFlexibleTimeWindow
		}
	}
	return nil
}

// isValidScheduleExpression delegates to the shared schedule expression
// validator in internal/common/scheduleexpr.
func isValidScheduleExpression(expr string) bool {
	return scheduleexpr.ValidateExpression(expr)
}

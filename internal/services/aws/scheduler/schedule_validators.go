package scheduler

import (
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// ScheduleSpec is the common input structure for schedule creation and
// update, used by both the HTTP API and the admin console to guarantee
// identical validation (H1 shared validation layer).
type ScheduleSpec struct {
	Name                       string
	GroupName                  string
	ScheduleExpression         string
	ScheduleExpressionTimezone string
	Description                string
	State                      string
	KmsKeyArn                  string
	StartDate                  string // raw RFC3339 / ISO8601 string
	EndDate                    string // raw RFC3339 / ISO8601 string
	ActionAfterCompletion      string
	Target                     *schedulerstore.Target
	FlexibleTimeWindow         *schedulerstore.FlexibleTimeWindow
}

// ValidatedSchedule holds the parsed and validated schedule fields ready
// for store persistence. Produced by ValidateScheduleFields so that both
// the HTTP API and admin paths build the store model identically.
type ValidatedSchedule struct {
	State                 schedulerstore.ScheduleState
	ActionAfterCompletion schedulerstore.ActionAfterCompletion
	StartDate             *time.Time
	EndDate               *time.Time
}

// validateClientToken validates the ClientToken format per Smithy spec:
// length [1, 64], pattern ^[a-zA-Z0-9-_]+$ (M5).
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

// dateLayouts lists the timestamp formats accepted by the Scheduler API.
// AWS accepts RFC3339, ISO8601 (with and without milliseconds), and the
// simple date-only format from the admin console.
var dateLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	timeutils.ISO8601UTCFormat,
	timeutils.ISO8601NoZFormat,
	"2006-01-02",
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

// ValidateScheduleFields validates all schedule fields per AWS Smithy spec.
// Called by both HTTP API and admin console paths (H1 — shared validation).
// Returns the normalised schedule fields or an error.
func ValidateScheduleFields(spec *ScheduleSpec) (*ValidatedSchedule, error) {
	// 1. Name pattern (H1).
	if spec.Name == "" || !namePattern.MatchString(spec.Name) {
		return nil, ErrValidation
	}

	// 2. Schedule expression validity including cron field count (H1, H5)
	//    and at() semantic date check (H6).
	if spec.ScheduleExpression == "" {
		return nil, ErrValidation
	}
	if !isValidScheduleExpression(spec.ScheduleExpression) {
		return nil, ErrInvalidScheduleExpression
	}

	// 3. Target required + ARN format validation (H1, M1, H2, L5).
	if spec.Target == nil {
		return nil, ErrInvalidTarget
	}
	if err := validateTarget(spec.Target); err != nil {
		return nil, err
	}

	// 4. FlexibleTimeWindow Mode enum (H1, H3).
	if spec.FlexibleTimeWindow != nil {
		if err := validateFlexibleTimeWindow(spec.FlexibleTimeWindow); err != nil {
			return nil, err
		}
	}

	// State defaults to ENABLED and ActionAfterCompletion defaults to NONE
	// per the AWS Scheduler API contract — both for CreateSchedule and the
	// PUT semantics of UpdateSchedule, where omitted fields are reset to
	// the documented default.
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

	// 6. ActionAfterCompletion enum (H1).
	if spec.ActionAfterCompletion != "" {
		if spec.ActionAfterCompletion != "NONE" && spec.ActionAfterCompletion != "DELETE" {
			return nil, ErrInvalidActionAfterCompletion
		}
		result.ActionAfterCompletion = schedulerstore.ActionAfterCompletion(spec.ActionAfterCompletion)
	}

	// 7. KmsKeyArn ARN validation (M2).
	if spec.KmsKeyArn != "" {
		parsed, err := svcarn.ParseARN(spec.KmsKeyArn)
		if err != nil {
			return nil, awserrors.NewValidationException("invalid KmsKeyArn ARN format")
		}
		if parsed.Service != "kms" {
			return nil, awserrors.NewValidationException("KmsKeyArn must reference a KMS key")
		}
	}

	// 8. ScheduleExpressionTimezone length + IANA TZ database (M3).
	if spec.ScheduleExpressionTimezone != "" {
		if len(spec.ScheduleExpressionTimezone) < 1 || len(spec.ScheduleExpressionTimezone) > 50 {
			return nil, awserrors.NewValidationException("ScheduleExpressionTimezone must be 1-50 characters")
		}
		if _, err := time.LoadLocation(spec.ScheduleExpressionTimezone); err != nil {
			return nil, awserrors.NewValidationException("invalid ScheduleExpressionTimezone: not a valid IANA timezone")
		}
	}

	// 9. Description length 0-512 (M6).
	if len(spec.Description) > 512 {
		return nil, awserrors.NewValidationException("Description must be 0-512 characters")
	}

	// 10. StartDate / EndDate parse + ordering (M4).
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

// validateTarget validates the Target structure: ARN format for Target,
// RoleArn, and DeadLetterConfig (M1); RetryPolicy ranges (H2); and
// EventBridge / Kinesis sub-parameter lengths (L5).
func validateTarget(target *schedulerstore.Target) error {
	if target.Arn == "" {
		return ErrInvalidTarget
	}
	if _, err := svcarn.ParseARN(target.Arn); err != nil {
		return ErrInvalidTarget
	}
	if target.RoleArn == "" {
		return ErrInvalidTarget
	}
	if _, err := svcarn.ParseARN(target.RoleArn); err != nil {
		return ErrInvalidTarget
	}

	// DeadLetterConfig ARN format (M1).
	if target.DeadLetterConfig != nil && target.DeadLetterConfig.Arn != "" {
		if _, err := svcarn.ParseARN(target.DeadLetterConfig.Arn); err != nil {
			return ErrInvalidTarget
		}
	}

	// RetryPolicy ranges — reject out-of-range instead of silently
	// dropping (H2). Smithy: MaximumEventAgeInSeconds [60, 86400],
	// MaximumRetryAttempts [0, 185].
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

	// EventBridge / Kinesis sub-parameter lengths (L5).
	// Smithy: DetailType [1, 128], Source [1, 256], PartitionKey [1, 256].
	if target.EventBridgeParameters != nil {
		if l := len(target.EventBridgeParameters.DetailType); l < 1 || l > 128 {
			return awserrors.NewValidationException("EventBridgeParameters.DetailType must be 1-128 characters")
		}
		if l := len(target.EventBridgeParameters.Source); l < 1 || l > 256 {
			return awserrors.NewValidationException("EventBridgeParameters.Source must be 1-256 characters")
		}
	}
	if target.KinesisParameters != nil {
		if l := len(target.KinesisParameters.PartitionKey); l < 1 || l > 256 {
			return awserrors.NewValidationException("KinesisParameters.PartitionKey must be 1-256 characters")
		}
	}

	return nil
}

// validateFlexibleTimeWindow validates the FlexibleTimeWindow Mode enum
// (H3 — Smithy enum [OFF, FLEXIBLE]) and the MaximumWindowInMinutes range
// when Mode is FLEXIBLE (Smithy range [1, 1440]).
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

// validateScheduleExpressionFull performs the same checks as
// isValidScheduleExpression but also includes cron field-count (H5) and
// at() semantic date validation (H6) so that invalid expressions are
// rejected at the API layer rather than silently stored.
func isValidScheduleExpressionFull(expr string) bool {
	// Smithy model: ScheduleExpression max length is 256 characters.
	if len(expr) > 256 {
		return false
	}

	// at() — validate semantic date correctness (H6). The regex captures
	// the date component; we additionally parse it to reject impossible
	// values such as month 13 or hour 99.
	if matches := atExpressionRegex.FindStringSubmatch(expr); len(matches) == 2 {
		if _, err := time.Parse(timeutils.ISO8601NoZFormat, matches[1]); err != nil {
			return false
		}
		return true
	}

	// rate() expression.
	if validateRateExpression(expr) {
		return true
	}

	// cron() — validate exactly 6 fields (H5). AWS cron requires:
	// Minutes Hours Day-of-month Month Day-of-week Year.
	if matches := cronExpressionRegex.FindStringSubmatch(expr); len(matches) == 2 {
		fields := strings.Fields(matches[1])
		if len(fields) != 6 {
			return false
		}
		return true
	}

	return false
}

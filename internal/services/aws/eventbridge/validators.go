package eventbridge

import (
	"fmt"
	"regexp"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Pattern constants extracted from the Smithy model
// (third_party/api-models-aws/models/eventbridge/service/2015-10-07/).
//
// resourceNamePattern applies to RuleName, ArchiveName, ConnectionName and
// ApiDestinationName (Smithy @pattern ^[\.\-_A-Za-z0-9]+$).
var resourceNamePattern = regexp.MustCompile(`^[\.\-_A-Za-z0-9]+$`)

// eventBusNamePattern applies to EventBusName (Smithy @pattern includes '/').
var eventBusNamePattern = regexp.MustCompile(`^[/\.\-_A-Za-z0-9]+$`)

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

// validateReplayName validates ReplayName per Smithy:
// pattern ^[\.\-_A-Za-z0-9]+$, length 1-64.
func validateReplayName(name string) bool {
	if name == "" || len(name) > 64 {
		return false
	}
	return resourceNamePattern.MatchString(name)
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

// validLaunchTypes mirrors the Smithy LaunchType enum (EC2, FARGATE, EXTERNAL).
var validLaunchTypes = map[string]bool{
	"EC2":      true,
	"FARGATE":  true,
	"EXTERNAL": true,
}

// validateLaunchType returns true when *lt* is a valid Smithy LaunchType.
func validateLaunchType(lt string) bool {
	return validLaunchTypes[lt]
}

// ---------------------------------------------------------------------------
// Length validators (AWS docs — Smithy targets unadorned String)
// ---------------------------------------------------------------------------

// validateDetailType enforces the AWS PutEventsRequestEntry.DetailType maximum
// length of 128 characters (documented in the AWS API Reference).
func validateDetailType(s string) bool {
	return len(s) > 0 && len(s) <= 128
}

// validateSource enforces the AWS PutEventsRequestEntry.Source maximum length
// of 256 characters (documented in the AWS API Reference).
func validateSource(s string) bool {
	return len(s) > 0 && len(s) <= 256
}

// validateTraceHeader enforces the Smithy TraceHeader @length(min=1,max=500).
// An empty string is accepted because the parameter is optional; a non-empty
// value must satisfy max=500.
func validateTraceHeader(s string) bool {
	return len(s) <= 500
}

// ---------------------------------------------------------------------------
// Range validators
// ---------------------------------------------------------------------------

// validateRetryPolicy validates MaximumRetryAttempts (Smithy @range 0-185) and
// MaximumEventAgeInSeconds (Smithy @range 60-86400).  A nil policy is valid.
// A zero MaximumEventAgeInSeconds means "use default" and is accepted.
func validateRetryPolicy(rp *eventsstore.RetryPolicy) bool {
	if rp == nil {
		return true
	}
	if rp.MaximumRetryAttempts < 0 || rp.MaximumRetryAttempts > 185 {
		return false
	}
	if rp.MaximumEventAgeInSeconds != 0 {
		if rp.MaximumEventAgeInSeconds < 60 || rp.MaximumEventAgeInSeconds > 86400 {
			return false
		}
	}
	return true
}

// validateKmsKeyIdentifier verifies that the value is a well-formed KMS key
// ARN across any partition (aws, aws-cn, aws-us-gov).  An empty string is
// valid — the parameter is optional.
func validateKmsKeyIdentifier(arn string) bool {
	if arn == "" {
		return true
	}
	partition, service, _, _, _ := svcarn.SplitARN(arn)
	return strings.HasPrefix(partition, "aws") && service == "kms"
}

// ---------------------------------------------------------------------------
// Additional length validators (Smithy @length traits)
// ---------------------------------------------------------------------------

// maxDescriptionLength is the Smithy @length(0,512) bound shared by
// RuleDescription, EventBusDescription, ArchiveDescription,
// ConnectionDescription, ApiDestinationDescription and ReplayDescription.
const maxDescriptionLength = 512

// validateDescription enforces the Smithy @length(0,512) trait shared by
// RuleDescription, EventBusDescription, ArchiveDescription,
// ConnectionDescription, ApiDestinationDescription and ReplayDescription.
func validateDescription(s string) bool {
	return len(s) <= maxDescriptionLength
}

// errDescriptionTooLong is the validation error for a description member
// exceeding the shared @length(0,512) bound.
func errDescriptionTooLong() error {
	return awserrors.NewValidationException(
		fmt.Sprintf("Description must be at most %d characters", maxDescriptionLength))
}

// validateEventPatternLength enforces the Smithy EventPattern @length(0,4096).
// The caller should additionally check JSON validity via isValidEventPattern.
func validateEventPatternLength(pattern string) bool {
	return len(pattern) <= 4096
}

// validateInvocationRateLimit enforces the AWS-documented maximum of 300
// invocations per second for an API destination.
func validateInvocationRateLimit(rate int32) bool {
	return rate >= 1 && rate <= 300
}

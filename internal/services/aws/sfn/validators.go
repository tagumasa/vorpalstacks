package sfn

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// ---------------------------------------------------------------------------
// Enum validation maps (Smithy-derived)
// ---------------------------------------------------------------------------

var validStateMachineTypes = map[string]bool{
	"STANDARD": true,
	"EXPRESS":  true,
}

var validLogLevels = map[string]bool{
	"ALL":   true,
	"ERROR": true,
	"FATAL": true,
	"OFF":   true,
}

var validEncryptionTypes = map[string]bool{
	"AWS_OWNED_KEY":            true,
	"CUSTOMER_MANAGED_KMS_KEY": true,
}

var validRedrivableStatuses = map[string]bool{
	"FAILED":    true,
	"TIMED_OUT": true,
	"ABORTED":   true,
}

var validTerminalStatuses = map[string]bool{
	"SUCCEEDED": true,
	"FAILED":    true,
	"TIMED_OUT": true,
	"ABORTED":   true,
}

// ---------------------------------------------------------------------------
// Validators
// ---------------------------------------------------------------------------

// forbiddenNameRunes mirrors the documented name contract shared by state
// machines, executions and activities: a name must not contain white
// space, brackets, wildcard characters, or the listed special characters.
var forbiddenNameRunes = "<>{}[]?*\"#%\\^|~`$&,;:/"

// isForbiddenNameRune reports whether a rune is excluded from resource
// names: the documented special characters, white space, the control
// ranges U+0000-001F and U+007F-009F, the noncharacters U+FFFE-FFFF, the
// surrogate range U+D800-DFFF and the invalid character U+10FFFF (listed
// separately from the noncharacters).
func isForbiddenNameRune(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}
	if strings.ContainsRune(forbiddenNameRunes, r) {
		return true
	}
	if r <= 0x001F || (r >= 0x007F && r <= 0x009F) {
		return true
	}
	if r >= 0xFFFE && r <= 0xFFFF {
		return true
	}
	if r >= 0xD800 && r <= 0xDFFF {
		return true
	}
	if r == unicode.MaxRune {
		return true
	}
	return false
}

// validateResourceName validates a state machine or activity name against
// the Smithy Name shape: @length(min=1, max=80) plus the documented
// forbidden-character contract.
func validateResourceName(name string) error {
	if name == "" {
		return NewInvalidName("name is required")
	}
	if len(name) > sfnstore.MaxResourceNameLength {
		return NewInvalidName(fmt.Sprintf("name must be 1-80 characters, got %d", len(name)))
	}
	if err := validateNameCharacters(name); err != nil {
		return err
	}
	return nil
}

// validateExecutionName validates the optional StartExecution /
// StartSyncExecution name against the same Name shape and
// forbidden-character contract.
func validateExecutionName(name string) error {
	if len(name) > sfnstore.MaxResourceNameLength {
		return NewInvalidName(fmt.Sprintf("name must be 1-80 characters, got %d", len(name)))
	}
	return validateNameCharacters(name)
}

// validateNameCharacters rejects names carrying any forbidden rune.
func validateNameCharacters(name string) error {
	for _, r := range name {
		if isForbiddenNameRune(r) {
			return NewInvalidName(fmt.Sprintf("name must not contain white space, brackets, wildcard or special characters (found %q)", r))
		}
	}
	return nil
}

// validateStateMachineType validates the type parameter against the Smithy
// StateMachineType enum (STANDARD | EXPRESS). Returns the default "STANDARD"
// when empty.
func validateStateMachineType(smType string) (string, error) {
	if smType == "" {
		return "STANDARD", nil
	}
	if !validStateMachineTypes[smType] {
		return "", NewStateMachineTypeNotSupported(fmt.Sprintf("State Machine type must be STANDARD or EXPRESS, got: %s", smType))
	}
	return smType, nil
}

// validateLogLevel validates the level field of LoggingConfiguration against
// the Smithy LogLevel enum (ALL | ERROR | FATAL | OFF).
func validateLogLevel(level string) error {
	if level == "" {
		return nil
	}
	if !validLogLevels[level] {
		return NewInvalidLoggingConfiguration(fmt.Sprintf("loggingConfiguration.level must be one of ALL, ERROR, FATAL, OFF, got: %s", level))
	}
	return nil
}

// validateLoggingConfiguration validates a LoggingConfiguration struct against
// AWS specifications:
//   - level must be one of ALL/ERROR/FATAL/OFF (Smithy LogLevel enum)
//   - destinations is limited to size 1 (AWS docs: "Limited to size 1")
//   - destinations is required when level is not OFF (AWS docs)
func validateLoggingConfiguration(lc *sfnstore.LoggingConfiguration) error {
	if lc == nil {
		return nil
	}

	if err := validateLogLevel(lc.Level); err != nil {
		return err
	}

	if len(lc.Destinations) > 1 {
		return NewInvalidLoggingConfiguration("loggingConfiguration.destinations is limited to size 1")
	}

	if lc.Level != "OFF" && len(lc.Destinations) == 0 {
		return NewInvalidLoggingConfiguration("loggingConfiguration.destinations is required when level is not OFF")
	}

	return nil
}

// validateEncryptionConfiguration validates an EncryptionConfiguration struct:
//   - type must be AWS_OWNED_KEY or CUSTOMER_MANAGED_KMS_KEY (Smithy enum)
//   - kmsKeyId is required when type is CUSTOMER_MANAGED_KMS_KEY
//   - kmsDataKeyReusePeriodSeconds must be in [60, 900] (Smithy @range)
func validateEncryptionConfiguration(ec *sfnstore.EncryptionConfiguration) error {
	if ec.Type == "" {
		ec.Type = "AWS_OWNED_KEY"
	}
	if !validEncryptionTypes[ec.Type] {
		return NewInvalidEncryptionConfiguration(fmt.Sprintf("encryptionConfiguration.type must be AWS_OWNED_KEY or CUSTOMER_MANAGED_KMS_KEY, got %s", ec.Type))
	}
	if ec.Type == "CUSTOMER_MANAGED_KMS_KEY" && ec.KmsKeyId == "" {
		return NewInvalidEncryptionConfiguration("encryptionConfiguration.kmsKeyId is required when type is CUSTOMER_MANAGED_KMS_KEY")
	}
	if ec.KmsDataKeyReusePeriod != 0 {
		if ec.KmsDataKeyReusePeriod < 60 || ec.KmsDataKeyReusePeriod > 900 {
			return NewInvalidEncryptionConfiguration(fmt.Sprintf("encryptionConfiguration.kmsDataKeyReusePeriodSeconds must be in [60, 900], got %d", ec.KmsDataKeyReusePeriod))
		}
	}
	return nil
}

// validateMapRunUpdateParams validates the three numeric fields of UpdateMapRun
// against their Smithy @range traits:
//   - maxConcurrency: @range(min=0)
//   - toleratedFailureCount: @range(min=0)
//   - toleratedFailurePercentage: @range(min=0, max=100)
func validateMapRunUpdateParams(maxConcurrency int64, toleratedFailureCount int64, toleratedFailurePercentage float32) error {
	if maxConcurrency < 0 {
		return NewValidationException(fmt.Sprintf("maxConcurrency must be >= 0, got %d", maxConcurrency))
	}
	if toleratedFailureCount < 0 {
		return NewValidationException(fmt.Sprintf("toleratedFailureCount must be >= 0, got %d", toleratedFailureCount))
	}
	if toleratedFailurePercentage < 0 || toleratedFailurePercentage > 100 {
		return NewValidationException(fmt.Sprintf("toleratedFailurePercentage must be in [0, 100], got %f", toleratedFailurePercentage))
	}
	return nil
}

// validateSeverity validates the severity parameter of
// ValidateStateMachineDefinition and TestState against the Smithy enum
// (ERROR | WARNING).
func validateSeverity(severity string) (string, error) {
	if severity == "" {
		return "ERROR", nil
	}
	if severity != "ERROR" && severity != "WARNING" {
		return "", NewValidationException(fmt.Sprintf("severity must be ERROR or WARNING, got %s", severity))
	}
	return severity, nil
}

// validateDefinitionJSON validates that a state machine definition string
// is non-empty, valid JSON, and within the Definition shape's one-mebibyte
// bound (@length(1, 1048576)).
func validateDefinitionJSON(definition string) error {
	if definition == "" {
		return NewInvalidDefinitionException("State Machine definition is required")
	}
	if len(definition) > sfnstore.MaxDefinitionLength {
		return NewInvalidDefinitionException(fmt.Sprintf("State Machine definition must be at most %d bytes, got %d", sfnstore.MaxDefinitionLength, len(definition)))
	}
	if !json.Valid([]byte(definition)) {
		return NewInvalidDefinitionException("State Machine definition is not valid JSON")
	}
	return nil
}

// waitTimestampPattern pins the AWS Wait-state timestamp profile: the
// RFC3339 profile of ISO 8601 with an uppercase T separating date and
// time, an uppercase Z when no numeric offset is present, and fractional
// seconds of zero, three, six, or nine digits (per the ISO 8601 profile
// Step Functions follows).
var waitTimestampPattern = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3}|\.\d{6}|\.\d{9})?(Z|[+-]\d{2}:\d{2})$`)

// parseWaitTimestamp parses a Wait-state timestamp literal under the
// strict AWS profile, reporting false when the value does not conform.
func parseWaitTimestamp(value string) (time.Time, bool) {
	if !waitTimestampPattern.MatchString(value) {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}
	// AWS truncates wait timestamps to whole seconds.
	return t, true
}

// validateRoutingConfiguration validates the routing configuration of a state
// machine alias. Each weight must be in [0, 100] and the sum of all weights
// must equal 100.
func validateRoutingConfiguration(rc []sfnstore.RoutingConfiguration) error {
	if len(rc) == 0 {
		return nil
	}
	var totalWeight int32
	for _, entry := range rc {
		if entry.Weight < 0 || entry.Weight > 100 {
			return NewValidationException(fmt.Sprintf("routingConfiguration weight must be in [0, 100], got %d", entry.Weight))
		}
		if entry.StateMachineVersionArn == "" {
			return NewValidationException("routingConfiguration entry must include stateMachineVersionArn")
		}
		totalWeight += entry.Weight
	}
	if totalWeight != 100 {
		return NewValidationException(fmt.Sprintf("routingConfiguration weights must sum to 100, got %d", totalWeight))
	}
	return nil
}

// validateMaxResults validates a maxResults parameter against the given bounds.
// A value of 0 is treated as "use default" and is allowed.
func validateMaxResults(maxResults int32, minVal, maxVal int32, paramName string) error {
	if maxResults == 0 {
		return nil
	}
	if maxResults < minVal || maxResults > maxVal {
		return NewValidationException(fmt.Sprintf("%s must be in [%d, %d], got %d", paramName, minVal, maxVal, maxResults))
	}
	return nil
}

// validateArnRequired checks that a required ARN parameter is non-empty.
// The Smithy Arn shape carries @length(1,256) with no pattern, so lengths
// count Unicode characters.
func validateArnRequired(arn, paramName string) error {
	if arn == "" {
		return NewInvalidArnException(fmt.Sprintf("%s is required", paramName))
	}
	if n := utf8.RuneCountInString(arn); n > sfnstore.MaxArnLength {
		return NewInvalidArnException(fmt.Sprintf("%s must be 1-%d characters, got %d", paramName, sfnstore.MaxArnLength, n))
	}
	return nil
}

// validateRoleArnRequired checks that the roleArn parameter is provided.
// The Smithy model marks roleArn as @required on CreateStateMachineInput;
// the Arn shape carries no pattern, so lengths count Unicode characters.
func validateRoleArnRequired(roleArn string) error {
	if roleArn == "" {
		return NewMissingRequiredParameter("roleArn is a required parameter")
	}
	if n := utf8.RuneCountInString(roleArn); n > sfnstore.MaxArnLength {
		return NewInvalidArnException(fmt.Sprintf("roleArn must be 1-%d characters, got %d", sfnstore.MaxArnLength, n))
	}
	return nil
}

// validateRoleArnOptional validates roleArn when it is provided.
// On UpdateStateMachineInput, roleArn is optional. Lengths count Unicode
// characters (the Arn shape carries no pattern).
func validateRoleArnOptional(roleArn string) error {
	if roleArn == "" {
		return nil
	}
	if n := utf8.RuneCountInString(roleArn); n > sfnstore.MaxArnLength {
		return NewInvalidArnException(fmt.Sprintf("roleArn must be 1-%d characters, got %d", sfnstore.MaxArnLength, n))
	}
	return nil
}

// isValidExecutionStatus checks whether a status string is a member of the
// Smithy ExecutionStatus enum (RUNNING, SUCCEEDED, FAILED, TIMED_OUT,
// ABORTED, PENDING_REDRIVE). The empty string means "no filter" and is
// valid for the ListExecutions statusFilter parameter.
func isValidExecutionStatus(status string) bool {
	if status == "" {
		return true
	}
	switch status {
	case "RUNNING", "SUCCEEDED", "FAILED", "TIMED_OUT", "ABORTED", "PENDING_REDRIVE":
		return true
	}
	return false
}

// redriveStatusFor checks whether an execution status is redrivable per the
// AWS SFN specification.
func isRedrivableStatus(status string) bool {
	return validRedrivableStatuses[status]
}

// isTerminalStatus checks whether a status is terminal (no further
// transitions possible).
func isTerminalStatus(status string) bool {
	return validTerminalStatuses[status]
}

// validateSNSTopicInput validates that an SNS publish task input contains
// either a TopicArn or a TopicName. AWS requires TopicArn as a mandatory
// parameter; vorpalstacks also accepts TopicName for convenience but must
// not silently default to "default-topic" when neither is provided.
func validateSNSTopicInput(topicArn, topicName string) error {
	if topicArn == "" && topicName == "" {
		return fmt.Errorf("SNS publish requires TopicArn or TopicName")
	}
	if strings.Contains(topicName, ":") {
		return fmt.Errorf("TopicName must not contain ':' — use TopicArn for fully-qualified ARNs")
	}
	return nil
}

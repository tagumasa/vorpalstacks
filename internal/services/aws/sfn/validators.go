package sfn

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
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

// validateResourceName validates a state machine or activity name against the
// Smithy Name shape: @length(min=1, max=80).
func validateResourceName(name string) error {
	if name == "" {
		return NewInvalidName("name is required")
	}
	if len(name) > sfnstore.MaxResourceNameLength {
		return NewInvalidName(fmt.Sprintf("name must be 1-80 characters, got %d", len(name)))
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

// validateDefinitionJSON validates that a state machine definition string is
// non-empty and valid JSON.
func validateDefinitionJSON(definition string) error {
	if definition == "" {
		return NewInvalidDefinitionException("State Machine definition is required")
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

// waitStateDiagnostics inspects every Wait state of a definition and
// returns an ERROR diagnostic for each field-contract violation: a JSON
// Path Wait must specify exactly one of Seconds, Timestamp, SecondsPath,
// or TimestampPath; a JSONata Wait exactly one of Seconds or Timestamp.
// Timestamp literals must conform to the strict RFC3339 profile and
// Seconds literals must be integers in [0, MaxWaitSeconds]. Expressions
// in JSONata states are accepted by presence — their values are only
// checkable at run time.
func waitStateDiagnostics(definition string) []map[string]string {
	var skeleton struct {
		QueryLanguage string                     `json:"QueryLanguage"`
		States        map[string]json.RawMessage `json:"States"`
	}
	if err := json.Unmarshal([]byte(definition), &skeleton); err != nil || skeleton.States == nil {
		return nil
	}

	diagnostics := []map[string]string{}
	add := func(state, message string) {
		diagnostics = append(diagnostics, map[string]string{
			"severity": "ERROR",
			"code":     "InvalidDefinition",
			"message":  "SCHEMA_VALIDATION_FAILED: " + message + " at /States/" + state,
		})
	}

	for name, raw := range skeleton.States {
		var state struct {
			QueryLanguage string       `json:"QueryLanguage"`
			Type          string       `json:"Type"`
			Seconds       *interface{} `json:"Seconds"`
			Timestamp     *interface{} `json:"Timestamp"`
			SecondsPath   *interface{} `json:"SecondsPath"`
			TimestampPath *interface{} `json:"TimestampPath"`
		}
		if err := json.Unmarshal(raw, &state); err != nil {
			continue // malformed members surface through the JSONata and runtime checks
		}
		if state.Type != "Wait" {
			continue
		}

		queryLanguage := state.QueryLanguage
		if queryLanguage == "" {
			queryLanguage = skeleton.QueryLanguage
		}
		jsonata := queryLanguage == "JSONata"

		if jsonata && state.SecondsPath != nil {
			add(name, "Wait state SecondsPath is only supported in JSONPath states")
			continue
		}
		if jsonata && state.TimestampPath != nil {
			add(name, "Wait state TimestampPath is only supported in JSONPath states")
			continue
		}

		present := 0
		if state.Seconds != nil {
			present++
		}
		if state.Timestamp != nil {
			present++
		}
		if !jsonata {
			if state.SecondsPath != nil {
				present++
			}
			if state.TimestampPath != nil {
				present++
			}
		}
		if present != 1 {
			if jsonata {
				add(name, "Wait state must specify exactly one of Seconds or Timestamp")
			} else {
				add(name, "Wait state must specify exactly one of Seconds, Timestamp, SecondsPath, or TimestampPath")
			}
			continue
		}

		if state.Timestamp != nil {
			value, ok := (*state.Timestamp).(string)
			if !ok {
				add(name, "Wait state Timestamp must be a string")
				continue
			}
			if jsonata && IsExpression(value) {
				continue // the expression's result is only checkable at run time
			}
			if _, ok := parseWaitTimestamp(value); !ok {
				add(name, fmt.Sprintf("Wait state Timestamp %q must conform to the RFC3339 profile of ISO 8601 with an uppercase T and an uppercase Z or numeric offset, for example \"2024-03-14T01:59:00Z\"", value))
			}
			continue
		}

		if state.Seconds != nil {
			if value, ok := (*state.Seconds).(string); ok && jsonata && IsExpression(value) {
				continue // the expression's result is only checkable at run time
			}
			value, ok := (*state.Seconds).(float64)
			if !ok || value != math.Trunc(value) || value < 0 || value > sfnstore.MaxWaitSeconds {
				add(name, fmt.Sprintf("Wait state Seconds must be an integer value from 0 to %d", sfnstore.MaxWaitSeconds))
			}
			continue
		}

		// Remaining JSONPath cases: the path fields must be non-empty strings.
		for _, field := range []struct {
			name  string
			value *interface{}
		}{{"SecondsPath", state.SecondsPath}, {"TimestampPath", state.TimestampPath}} {
			if field.value == nil {
				continue
			}
			if value, ok := (*field.value).(string); !ok || value == "" {
				add(name, "Wait state "+field.name+" must be a non-empty path string")
			}
		}
	}
	return diagnostics
}

// validateWaitStates rejects a definition whose Wait states violate the
// field contract with the creation-time InvalidDefinitionException shape.
func validateWaitStates(definition string) error {
	diagnostics := waitStateDiagnostics(definition)
	if len(diagnostics) == 0 {
		return nil
	}
	return NewInvalidDefinitionException(diagnostics[0]["message"])
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
	if n := utf8.RuneCountInString(arn); n > 256 {
		return NewInvalidArnException(fmt.Sprintf("%s must be 1-256 characters, got %d", paramName, n))
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
	if n := utf8.RuneCountInString(roleArn); n > 256 {
		return NewInvalidArnException(fmt.Sprintf("roleArn must be 1-256 characters, got %d", n))
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
	if n := utf8.RuneCountInString(roleArn); n > 256 {
		return NewInvalidArnException(fmt.Sprintf("roleArn must be 1-256 characters, got %d", n))
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

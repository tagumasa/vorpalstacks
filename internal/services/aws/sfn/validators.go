package sfn

import (
	"encoding/json"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
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
	if len(name) > 80 {
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
		return "", NewInvalidExecutionType(fmt.Sprintf("State Machine type must be STANDARD or EXPRESS, got: %s", smType))
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
		return NewInvalidParameterValue(fmt.Sprintf("loggingConfiguration.level must be one of ALL, ERROR, FATAL, OFF, got: %s", level))
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
		return NewInvalidParameterValue("loggingConfiguration.destinations is limited to size 1")
	}

	if lc.Level != "OFF" && len(lc.Destinations) == 0 {
		return NewInvalidParameterValue("loggingConfiguration.destinations is required when level is not OFF")
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
		return NewInvalidParameterValue(fmt.Sprintf("encryptionConfiguration.type must be AWS_OWNED_KEY or CUSTOMER_MANAGED_KMS_KEY, got %s", ec.Type))
	}
	if ec.Type == "CUSTOMER_MANAGED_KMS_KEY" && ec.KmsKeyId == "" {
		return NewInvalidParameterValue("encryptionConfiguration.kmsKeyId is required when type is CUSTOMER_MANAGED_KMS_KEY")
	}
	if ec.KmsDataKeyReusePeriod != 0 {
		if ec.KmsDataKeyReusePeriod < 60 || ec.KmsDataKeyReusePeriod > 900 {
			return NewInvalidParameterValue(fmt.Sprintf("encryptionConfiguration.kmsDataKeyReusePeriodSeconds must be in [60, 900], got %d", ec.KmsDataKeyReusePeriod))
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
		return NewInvalidParameterValue(fmt.Sprintf("maxConcurrency must be >= 0, got %d", maxConcurrency))
	}
	if toleratedFailureCount < 0 {
		return NewInvalidParameterValue(fmt.Sprintf("toleratedFailureCount must be >= 0, got %d", toleratedFailureCount))
	}
	if toleratedFailurePercentage < 0 || toleratedFailurePercentage > 100 {
		return NewInvalidParameterValue(fmt.Sprintf("toleratedFailurePercentage must be in [0, 100], got %f", toleratedFailurePercentage))
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
		return "", NewInvalidParameterValue(fmt.Sprintf("severity must be ERROR or WARNING, got %s", severity))
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
			return NewInvalidParameterValue(fmt.Sprintf("routingConfiguration weight must be in [0, 100], got %d", entry.Weight))
		}
		if entry.StateMachineVersionArn == "" {
			return NewInvalidParameterValue("routingConfiguration entry must include stateMachineVersionArn")
		}
		totalWeight += entry.Weight
	}
	if totalWeight != 100 {
		return NewInvalidParameterValue(fmt.Sprintf("routingConfiguration weights must sum to 100, got %d", totalWeight))
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
		return NewInvalidParameterValue(fmt.Sprintf("%s must be in [%d, %d], got %d", paramName, minVal, maxVal, maxResults))
	}
	return nil
}

// validateArnRequired checks that a required ARN parameter is non-empty.
func validateArnRequired(arn, paramName string) error {
	if arn == "" {
		return NewInvalidArnException(fmt.Sprintf("%s is required", paramName))
	}
	if len(arn) > 256 {
		return NewInvalidArnException(fmt.Sprintf("%s must be 1-256 characters, got %d", paramName, len(arn)))
	}
	return nil
}

// validateRoleArnRequired checks that the roleArn parameter is provided.
// The Smithy model marks roleArn as @required on CreateStateMachineInput.
func validateRoleArnRequired(roleArn string) error {
	if roleArn == "" {
		return awserrors.NewAWSError("MissingRequiredParameter",
			"roleArn is a required parameter", 400)
	}
	if len(roleArn) > 256 {
		return NewInvalidArnException(fmt.Sprintf("roleArn must be 1-256 characters, got %d", len(roleArn)))
	}
	return nil
}

// validateRoleArnOptional validates roleArn when it is provided.
// On UpdateStateMachineInput, roleArn is optional.
func validateRoleArnOptional(roleArn string) error {
	if roleArn == "" {
		return nil
	}
	if len(roleArn) > 256 {
		return NewInvalidArnException(fmt.Sprintf("roleArn must be 1-256 characters, got %d", len(roleArn)))
	}
	return nil
}

// validateExecutionStatus checks whether a status string is a recognised
// execution status.
func isValidExecutionStatus(status string) bool {
	switch status {
	case "RUNNING", "SUCCEEDED", "FAILED", "TIMED_OUT", "ABORTED":
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

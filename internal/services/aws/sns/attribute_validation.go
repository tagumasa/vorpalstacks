package sns

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

// kmsKeyIDRegex validates a bare KMS key ID in UUID hex format
// (8-4-4-4-12 lowercase hex digits, case-insensitive).
var kmsKeyIDRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// validateSubscriptionAttribute validates well-known subscription attributes
// that have structured values. Unknown attributes pass through without
// validation (forward-compatible with future AWS additions).
func validateSubscriptionAttribute(name, value string) error {
	switch name {
	case "FilterPolicy":
		return validateFilterPolicy(value)
	case "FilterPolicyScope":
		return validateFilterPolicyScope(value)
	case "RedrivePolicy":
		return validateRedrivePolicy(value)
	}
	return nil
}

// validateTopicAttribute validates well-known topic attributes that have
// structured values. Unknown attributes pass through without validation
// (forward-compatible with future AWS additions).
func validateTopicAttribute(name, value string) error {
	switch name {
	case "DeliveryPolicy":
		return validateJSONAttribute(name, value)
	case "Policy":
		return validateJSONAttribute(name, value)
	case "DisplayName":
		if len(value) > 100 {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("DisplayName too long: %d characters (maximum 100)", len(value)))
		}
	case "KmsMasterKeyId":
		if value != "" && !strings.HasPrefix(value, "arn:") && !isValidKmsKeyId(value) {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid KmsMasterKeyId: %s", value))
		}
	}
	return nil
}

// validateJSONAttribute validates that the value is valid JSON.
func validateJSONAttribute(name, value string) error {
	if value == "" {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid %s: not valid JSON: %s", name, err.Error()))
	}
	return nil
}

// isValidKmsKeyId checks whether a value is a valid bare KMS key identifier:
// a UUID-format key ID or an alias name (prefixed with "alias/").
// Key ARNs and alias ARNs are handled by the caller via the "arn:" prefix
// check.
func isValidKmsKeyId(value string) bool {
	if strings.HasPrefix(value, "alias/") {
		return len(value) > len("alias/")
	}
	return kmsKeyIDRegex.MatchString(value)
}

// validateFilterPolicy validates the JSON structure of a subscription filter
// policy. AWS allows at most 100 attribute names in a single policy.
func validateFilterPolicy(value string) error {
	if value == "" {
		return nil
	}

	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(value), &policy); err != nil {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid filter policy: %s", err.Error()))
	}

	if len(policy) > 100 {
		return ErrFilterLimitExceeded
	}

	for attrName, raw := range policy {
		if attrName == "" {
			return awserrors.NewInvalidParameterException("Invalid filter policy: attribute name must not be empty")
		}
		if err := validateFilterPolicyValue(raw); err != nil {
			return err
		}
	}

	return nil
}

// validateFilterPolicyValue validates the value array for one attribute in
// the filter policy. Values must be an array of strings, numbers, booleans,
// or condition objects.
func validateFilterPolicyValue(raw interface{}) error {
	values, ok := raw.([]interface{})
	if !ok {
		return awserrors.NewInvalidParameterException("Invalid filter policy: value must be an array")
	}

	for _, v := range values {
		switch cond := v.(type) {
		case string, float64, bool:
			continue
		case map[string]interface{}:
			if len(cond) != 1 {
				return awserrors.NewInvalidParameterException(
					"Invalid filter policy: a condition object must have exactly one operator")
			}
			for operator := range cond {
				switch operator {
				case "prefix", "anything-but", "numeric", "exists":
				default:
					return awserrors.NewInvalidParameterException(
						fmt.Sprintf("Invalid filter policy: unknown operator %q", operator))
				}
			}
		default:
			return awserrors.NewInvalidParameterException("Invalid filter policy: unsupported value type")
		}
	}

	return nil
}

// validateFilterPolicyScope validates the FilterPolicyScope attribute.
// AWS accepts only "MessageAttributes" (default) or "MessageBody".
func validateFilterPolicyScope(value string) error {
	switch value {
	case "MessageAttributes", "MessageBody":
		return nil
	default:
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Invalid FilterPolicyScope: %s. Valid values: MessageAttributes, MessageBody", value))
	}
}

// validateRedrivePolicy validates the JSON structure of a subscription
// redrive policy. AWS requires deadLetterTargetArn.
func validateRedrivePolicy(value string) error {
	if value == "" {
		return nil
	}

	var rp struct {
		DeadLetterTargetArn string `json:"deadLetterTargetArn"`
	}
	if err := json.Unmarshal([]byte(value), &rp); err != nil {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid redrive policy: %s", err.Error()))
	}

	if strings.TrimSpace(rp.DeadLetterTargetArn) == "" {
		return awserrors.NewInvalidParameterException("Invalid redrive policy: deadLetterTargetArn is required")
	}

	return nil
}

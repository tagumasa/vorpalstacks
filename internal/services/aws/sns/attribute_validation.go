package sns

import (
	"encoding/json"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

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
// AWS accepts only "MessageAttributes" or "MessageBodyAttributes".
func validateFilterPolicyScope(value string) error {
	switch value {
	case "MessageAttributes", "MessageBodyAttributes":
		return nil
	default:
		return awserrors.NewInvalidParameterException(
			fmt.Sprintf("Invalid FilterPolicyScope: %s. Valid values: MessageAttributes, MessageBodyAttributes", value))
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

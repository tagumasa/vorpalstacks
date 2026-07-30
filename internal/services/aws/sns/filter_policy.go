package sns

import (
	"encoding/json"
	"strconv"
	"strings"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// matchFilterPolicy evaluates whether a message satisfies the subscription's
// filter policy, respecting the FilterPolicyScope setting.
//
// AWS SNS filter policy semantics:
//   - Each key in the policy must match (AND semantics across keys).
//   - Multiple values for a key use OR semantics.
//   - A missing message attribute fails the key unless the policy value
//     contains {"exists": false}.
//
// When scope is "MessageBody", the policy keys refer to top-level properties
// of the JSON message body instead of message attributes.
//
// Supported operators: exact match, {"prefix": "..."}, {"anything-but": [...]},
// {"numeric": [op, val, ...]}, {"exists": bool}.
func matchFilterPolicy(policyJSON string, scope string, attrs map[string]*snsstore.MessageAttribute, messageBody string) bool {
	if strings.TrimSpace(policyJSON) == "" || policyJSON == "{}" {
		return true
	}

	var policy map[string]json.RawMessage
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		return false
	}

	var matchAttrs map[string]*snsstore.MessageAttribute
	if scope == "MessageBody" {
		matchAttrs = parseMessageBodyAttributes(messageBody)
	} else {
		matchAttrs = attrs
	}

	for attrName, rawValues := range policy {
		var values []interface{}
		if err := json.Unmarshal(rawValues, &values); err != nil {
			return false
		}

		msgAttr, attrExists := matchAttrs[attrName]
		if !matchPolicyValues(values, msgAttr, attrExists) {
			return false
		}
	}

	return true
}

// parseMessageBodyAttributes parses a JSON message body into a pseudo
// MessageAttribute map for MessageBody-scope filter policy matching.
// Only top-level properties are extracted; nested objects/arrays are ignored.
func parseMessageBodyAttributes(body string) map[string]*snsstore.MessageAttribute {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		return nil
	}

	result := make(map[string]*snsstore.MessageAttribute, len(bodyMap))
	for k, v := range bodyMap {
		switch val := v.(type) {
		case string:
			result[k] = &snsstore.MessageAttribute{Type: "String", StringValue: val}
		case float64:
			result[k] = &snsstore.MessageAttribute{Type: "Number", StringValue: strconv.FormatFloat(val, 'f', -1, 64)}
		case bool:
			result[k] = &snsstore.MessageAttribute{Type: "String", StringValue: strconv.FormatBool(val)}
		default:
			raw, err := json.Marshal(val)
			if err != nil {
				continue
			}
			result[k] = &snsstore.MessageAttribute{Type: "String.Array", StringValue: string(raw)}
		}
	}
	return result
}

// matchPolicyValues checks whether any of the policy values match the
// given message attribute (OR semantics within a single key).
func matchPolicyValues(values []interface{}, msgAttr *snsstore.MessageAttribute, attrExists bool) bool {
	for _, v := range values {
		if matchSinglePolicyValue(v, msgAttr, attrExists) {
			return true
		}
	}
	return false
}

// matchSinglePolicyValue evaluates one policy entry against the message attribute.
func matchSinglePolicyValue(entry interface{}, msgAttr *snsstore.MessageAttribute, attrExists bool) bool {
	switch val := entry.(type) {
	case string:
		return attrExists && msgAttr.StringValue == val
	case bool:
		return attrExists && msgAttr.StringValue == strconv.FormatBool(val)
	case float64:
		return attrExists && msgAttr.StringValue == strconv.FormatFloat(val, 'f', -1, 64)
	case map[string]interface{}:
		for operator, operand := range val {
			matched := applyFilterOperator(operator, operand, msgAttr, attrExists)
			if !matched {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// applyFilterOperator dispatches a single filter operator.
func applyFilterOperator(operator string, operand interface{}, msgAttr *snsstore.MessageAttribute, attrExists bool) bool {
	switch operator {
	case "exists":
		want, _ := operand.(bool)
		if want {
			return attrExists
		}
		return !attrExists

	case "prefix":
		if !attrExists {
			return false
		}
		prefix, _ := operand.(string)
		return strings.HasPrefix(msgAttr.StringValue, prefix)

	case "anything-but":
		if !attrExists {
			return false
		}
		return !matchAnythingBut(operand, msgAttr.StringValue)

	case "numeric":
		if !attrExists {
			return false
		}
		return matchNumericRange(operand, msgAttr.StringValue)

	default:
		return false
	}
}

// matchAnythingBut returns true when the attribute value equals one of the
// excluded values.
func matchAnythingBut(operand interface{}, attrValue string) bool {
	switch v := operand.(type) {
	case string:
		return attrValue == v
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok && attrValue == s {
				return true
			}
			if f, ok := item.(float64); ok {
				if attrValue == strconv.FormatFloat(f, 'f', -1, 64) {
					return true
				}
			}
		}
	}
	return false
}

// matchNumericRange evaluates {"numeric": [">=", 0, "<", 100]} style ranges.
func matchNumericRange(operand interface{}, attrValue string) bool {
	conditions, ok := operand.([]interface{})
	if !ok || len(conditions) < 2 || len(conditions)%2 != 0 {
		return false
	}

	target, err := strconv.ParseFloat(attrValue, 64)
	if err != nil {
		return false
	}

	for i := 0; i < len(conditions); i += 2 {
		op, _ := conditions[i].(string)
		bound, _ := conditions[i+1].(float64)

		switch op {
		case "=":
			if target != bound {
				return false
			}
		case ">":
			if !(target > bound) {
				return false
			}
		case ">=":
			if !(target >= bound) {
				return false
			}
		case "<":
			if !(target < bound) {
				return false
			}
		case "<=":
			if !(target <= bound) {
				return false
			}
		default:
			return false
		}
	}

	return true
}

package eventbridge

import (
	"encoding/json"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

// maxPatternCombinations is the ceiling on $or rule combinations: "APIs
// that accept an event pattern (such as PutRule, CreateArchive,
// UpdateArchive, and TestEventPattern) will throw an
// InvalidEventPatternException if the use of $or results in over 1000 rule
// combinations."
const maxPatternCombinations = 1000

// validateEventPatternStructure checks an event pattern against the
// documented pattern vocabulary so that an unmatchable or AWS-illegal
// pattern is rejected at acceptance with InvalidEventPatternException
// instead of being persisted to fail silently at every match. A pattern
// value must be a scalar, a list of those, a nested object level, or a
// single-key operator object with a documented operand form; "$or" must be
// an array of pattern objects, and the product of every $or array's entry
// count must stay within the combination ceiling.
func validateEventPatternStructure(pattern string) error {
	if pattern == "" {
		return nil
	}
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(pattern), &root); err != nil {
		return awserrors.NewInvalidEventPatternException(fmt.Sprintf("EventPattern is not valid JSON: %s", err))
	}

	combinations := 1
	if err := validatePatternLevel(root, &combinations); err != nil {
		return err
	}
	if combinations > maxPatternCombinations {
		return awserrors.NewInvalidEventPatternException(fmt.Sprintf(
			"EventPattern use of $or results in %d rule combinations; the maximum is %d",
			combinations, maxPatternCombinations))
	}
	return nil
}

func validatePatternLevel(level map[string]interface{}, combinations *int) error {
	for key, value := range level {
		if key == "$or" {
			entries, ok := value.([]interface{})
			if !ok || len(entries) == 0 {
				return awserrors.NewInvalidEventPatternException("$or must be a non-empty array of event pattern objects")
			}
			for _, entry := range entries {
				entryMap, ok := entry.(map[string]interface{})
				if !ok {
					return awserrors.NewInvalidEventPatternException("$or entries must be event pattern objects")
				}
				if err := validatePatternLevel(entryMap, combinations); err != nil {
					return err
				}
			}
			*combinations *= len(entries)
			continue
		}
		if err := validatePatternValue(key, value, combinations); err != nil {
			return err
		}
	}
	return nil
}

// validatePatternValue validates one member value of a pattern level:
// scalars pass, lists validate per element (list objects must be operator
// forms), and a direct object value is either a single-key operator (whose
// operand form is checked) or a nested level.
func validatePatternValue(member string, value interface{}, combinations *int) error {
	switch v := value.(type) {
	case map[string]interface{}:
		if op, operand, isOp := soleOperator(v); isOp {
			if err := validateOperand(op, operand); err != nil {
				return awserrors.NewInvalidEventPatternException(fmt.Sprintf("EventPattern member %q: %s", member, err.Error()))
			}
			return nil
		}
		return validatePatternLevel(v, combinations)
	case []interface{}:
		for _, item := range v {
			if err := validatePatternListElement(member, item); err != nil {
				return err
			}
		}
		return nil
	case string, float64, bool, nil:
		return nil
	default:
		return awserrors.NewInvalidEventPatternException(fmt.Sprintf(
			"EventPattern member %q: unsupported value type %T", member, value))
	}
}

// validatePatternListElement validates one element of a list-valued member:
// scalars and operator objects are documented; a nested pattern level
// inside a list is not, and neither is an unknown operator — both could
// never match and are rejected at acceptance.
func validatePatternListElement(member string, item interface{}) error {
	switch v := item.(type) {
	case string, float64, bool, nil:
		return nil
	case map[string]interface{}:
		op, operand, isOp := soleOperator(v)
		if !isOp {
			return awserrors.NewInvalidEventPatternException(fmt.Sprintf(
				"EventPattern member %q: unknown operator %q (operators are prefix, suffix, equals-ignore-case, wildcard, numeric, anything-but, exists, cidr)",
				member, soleKey(v)))
		}
		if err := validateOperand(op, operand); err != nil {
			return awserrors.NewInvalidEventPatternException(fmt.Sprintf("EventPattern member %q: %s", member, err.Error()))
		}
		return nil
	default:
		return awserrors.NewInvalidEventPatternException(fmt.Sprintf(
			"EventPattern member %q: unsupported value type %T in list", member, item))
	}
}

// soleOperator reports the operator name and operand when the map is a
// single-entry known-operator form.
func soleOperator(m map[string]interface{}) (string, interface{}, bool) {
	if len(m) != 1 {
		return "", nil, false
	}
	for key, value := range m {
		if isKnownOperator(key) {
			return key, value, true
		}
	}
	return "", nil, false
}

// soleKey returns the only key of a single-entry map, or "" otherwise.
func soleKey(m map[string]interface{}) string {
	if len(m) != 1 {
		return ""
	}
	for key := range m {
		return key
	}
	return ""
}

func validateOperand(op string, operand interface{}) error {
	switch op {
	case "exists":
		if _, ok := operand.(bool); !ok {
			return fmt.Errorf("exists takes true or false")
		}
		return nil
	case "prefix", "suffix":
		if _, _, ok := prefixSuffixOperand(operand); ok {
			return nil
		}
		return fmt.Errorf("%s takes a string or {\"equals-ignore-case\": \"...\"}", op)
	case "equals-ignore-case":
		if _, ok := operand.(string); !ok {
			return fmt.Errorf("equals-ignore-case takes a string")
		}
		return nil
	case "wildcard":
		return validateWildcardPattern(operand)
	case "numeric":
		return validateNumericOperand(operand)
	case "cidr":
		if _, ok := operand.(string); !ok {
			return fmt.Errorf("cidr takes a string")
		}
		return nil
	case "anything-but":
		return validateAnythingButOperand(operand)
	default:
		return fmt.Errorf("unknown operator %q", op)
	}
}

// validateWildcardPattern enforces the documented wildcard constraints:
// "consecutive wildcard characters are not supported" and the backslash
// escapes only '*' and '\'.
func validateWildcardPattern(operand interface{}) error {
	pattern, ok := operand.(string)
	if !ok {
		return fmt.Errorf("wildcard takes a string")
	}
	if _, splittable := splitWildcardPattern(pattern); !splittable {
		return fmt.Errorf("wildcard escapes may only escape '*' and '\\'")
	}
	if strings.Contains(pattern, "**") {
		return fmt.Errorf("consecutive wildcard characters are not supported")
	}
	return nil
}

// validateNumericOperand enforces the numeric operator/value pair form
// documented for numeric matching: [comparator, number] pairs, one or more.
func validateNumericOperand(operand interface{}) error {
	operands, ok := operand.([]interface{})
	if !ok || len(operands) < 2 || len(operands)%2 != 0 {
		return fmt.Errorf("numeric takes an array of comparator and number pairs, e.g. [\"=\", 100] or [\">\", 10, \"<=\", 20]")
	}
	for i := 0; i < len(operands); i += 2 {
		comp, ok := operands[i].(string)
		if !ok || !isNumericComparator(comp) {
			return fmt.Errorf("numeric comparators are =, <, <=, > and >=")
		}
		if _, ok := operands[i+1].(float64); !ok {
			return fmt.Errorf("numeric comparison values must be numbers")
		}
	}
	return nil
}

func isNumericComparator(op string) bool {
	switch op {
	case "=", "<", "<=", ">", ">=":
		return true
	default:
		return false
	}
}

// validateAnythingButOperand enforces the documented anything-but operand
// forms: string and number scalars ("You can use anything-but matching
// with strings and numeric values", the comparison-operators page), lists
// of only strings or only numbers, and the nested
// prefix/suffix/equals-ignore-case/wildcard conjunctions, each taking a
// string or a list of strings. Booleans and null are exact-match values
// but not anything-but values.
func validateAnythingButOperand(operand interface{}) error {
	switch v := operand.(type) {
	case string, float64:
		return nil
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("anything-but lists must not be empty")
		}
		sawString, sawNumber := false, false
		for _, item := range v {
			switch item.(type) {
			case string:
				sawString = true
			case float64:
				sawNumber = true
			default:
				return fmt.Errorf("anything-but lists contain only strings or only numbers")
			}
		}
		if sawString && sawNumber {
			return fmt.Errorf("anything-but lists contain only strings or only numbers")
		}
		return nil
	case map[string]interface{}:
		if len(v) != 1 {
			return fmt.Errorf("anything-but conjunctions take exactly one operator")
		}
		for op, inner := range v {
			switch op {
			case "prefix", "suffix", "equals-ignore-case":
				if _, ok := stringListOperand(inner); !ok {
					return fmt.Errorf("anything-but %s takes a string or a list of strings", op)
				}
				return nil
			case "wildcard":
				list, ok := stringListOperand(inner)
				if !ok {
					return fmt.Errorf("anything-but wildcard takes a string or a list of strings")
				}
				for _, item := range list {
					if err := validateWildcardPattern(item); err != nil {
						return err
					}
				}
				return nil
			default:
				return fmt.Errorf("anything-but conjunctions support prefix, suffix, equals-ignore-case and wildcard")
			}
		}
	}
	return fmt.Errorf("anything-but takes a scalar, a list of scalars, or a single conjunction object")
}

// Package policy provides IAM policy evaluation functionality for vorpalstacks.
package policy

import (
	"net"
	"path/filepath"
	"strings"
	"time"

	"vorpalstacks/internal/utils/timeutils"
)

// ConditionOperator represents the type of condition operator used in IAM policy evaluation.
type ConditionOperator string

// Condition operators for string comparisons.
const (
	ConditionStringEquals              ConditionOperator = "StringEquals"
	ConditionStringNotEquals           ConditionOperator = "StringNotEquals"
	ConditionStringEqualsIgnoreCase    ConditionOperator = "StringEqualsIgnoreCase"
	ConditionStringNotEqualsIgnoreCase ConditionOperator = "StringNotEqualsIgnoreCase"
	ConditionStringLike                ConditionOperator = "StringLike"
	ConditionStringNotLike             ConditionOperator = "StringNotLike"
)

// Condition operators for numeric comparisons.
const (
	ConditionNumericEquals            ConditionOperator = "NumericEquals"
	ConditionNumericNotEquals         ConditionOperator = "NumericNotEquals"
	ConditionNumericLessThan          ConditionOperator = "NumericLessThan"
	ConditionNumericLessThanEquals    ConditionOperator = "NumericLessThanEquals"
	ConditionNumericGreaterThan       ConditionOperator = "NumericGreaterThan"
	ConditionNumericGreaterThanEquals ConditionOperator = "NumericGreaterThanEquals"
)

// Condition operators for date comparisons.
const (
	ConditionDateEquals            ConditionOperator = "DateEquals"
	ConditionDateNotEquals         ConditionOperator = "DateNotEquals"
	ConditionDateLessThan          ConditionOperator = "DateLessThan"
	ConditionDateLessThanEquals    ConditionOperator = "DateLessThanEquals"
	ConditionDateGreaterThan       ConditionOperator = "DateGreaterThan"
	ConditionDateGreaterThanEquals ConditionOperator = "DateGreaterThanEquals"
)

// Condition operators for other comparisons.
const (
	ConditionBool                        ConditionOperator = "Bool"
	ConditionIpAddress                   ConditionOperator = "IpAddress"
	ConditionNotIpAddress                ConditionOperator = "NotIpAddress"
	ConditionArnEquals                   ConditionOperator = "ArnEquals"
	ConditionArnLike                     ConditionOperator = "ArnLike"
	ConditionArnNotEquals                ConditionOperator = "ArnNotEquals"
	ConditionArnNotLike                  ConditionOperator = "ArnNotLike"
	ConditionNull                        ConditionOperator = "Null"
	ConditionBinaryEquals                ConditionOperator = "BinaryEquals"
	ConditionForAnyValueStringEquals     ConditionOperator = "ForAnyValue:StringEquals"
	ConditionForAllValuesStringEquals    ConditionOperator = "ForAllValues:StringEquals"
	ConditionForAnyValueStringNotEquals  ConditionOperator = "ForAnyValue:StringNotEquals"
	ConditionForAllValuesStringNotEquals ConditionOperator = "ForAllValues:StringNotEquals"
)

// ConditionEvaluator evaluates IAM policy conditions against an evaluation context.
type ConditionEvaluator struct{}

// NewConditionEvaluator creates a new ConditionEvaluator instance.
func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{}
}

// Evaluate evaluates IAM policy conditions against an evaluation context.
func (ce *ConditionEvaluator) Evaluate(conditions ConditionMap, context *EvaluationContext) bool {
	if len(conditions) == 0 {
		return true
	}

	for operator, keyValues := range conditions {
		op := ConditionOperator(operator)

		if op == ConditionNull {
			for key, values := range keyValues {
				resolved := context.ResolveVariable(key)
				keyExists := resolved != ""
				if len(values) > 0 {
					expectNull := strings.EqualFold(values[0], "true")
					if expectNull && keyExists {
						return false
					}
					if !expectNull && !keyExists {
						return false
					}
				}
			}
			continue
		}

		if op == ConditionForAllValuesStringEquals {
			for key, values := range keyValues {
				resolvedKey := context.ResolveVariable(key)
				contextValues := strings.Split(resolvedKey, ",")
				if len(contextValues) == 0 {
					continue
				}
				for _, cv := range contextValues {
					cv = strings.TrimSpace(cv)
					matched := false
					for _, pv := range values {
						if cv == pv {
							matched = true
							break
						}
					}
					if !matched {
						return false
					}
				}
			}
			continue
		}

		if op == ConditionForAllValuesStringNotEquals {
			for key, values := range keyValues {
				resolvedKey := context.ResolveVariable(key)
				contextValues := strings.Split(resolvedKey, ",")
				if len(contextValues) == 0 {
					continue
				}
				for _, cv := range contextValues {
					cv = strings.TrimSpace(cv)
					matched := false
					for _, pv := range values {
						if cv != pv {
							matched = true
							break
						}
					}
					if !matched {
						return false
					}
				}
			}
			continue
		}

		if op == ConditionForAnyValueStringEquals {
			for key, values := range keyValues {
				resolvedKey := context.ResolveVariable(key)
				contextValues := strings.Split(resolvedKey, ",")
				if len(contextValues) == 0 {
					return false
				}
				anyMatch := false
				for _, cv := range contextValues {
					cv = strings.TrimSpace(cv)
					for _, pv := range values {
						if cv == pv {
							anyMatch = true
							break
						}
					}
					if anyMatch {
						break
					}
				}
				if !anyMatch {
					return false
				}
			}
			continue
		}

		if op == ConditionForAnyValueStringNotEquals {
			for key, values := range keyValues {
				resolvedKey := context.ResolveVariable(key)
				contextValues := strings.Split(resolvedKey, ",")
				if len(contextValues) == 0 {
					return false
				}
				anyMatch := false
				for _, cv := range contextValues {
					cv = strings.TrimSpace(cv)
					for _, pv := range values {
						if cv != pv {
							anyMatch = true
							break
						}
					}
					if anyMatch {
						break
					}
				}
				if !anyMatch {
					return false
				}
			}
			continue
		}

		for key, values := range keyValues {
			resolvedKey := context.ResolveVariable(key)
			if !ce.evaluateOperator(op, resolvedKey, values, context) {
				return false
			}
		}
	}
	return true
}

func (ce *ConditionEvaluator) evaluateOperator(op ConditionOperator, key string, values []string, ctx *EvaluationContext) bool {
	contextValue := ctx.GetContextValue(key)
	if contextValue == "" {
		return false
	}

	for _, value := range values {
		if ce.matchValue(op, contextValue, value) {
			return true
		}
	}
	return false
}

func (ce *ConditionEvaluator) matchValue(op ConditionOperator, contextValue, policyValue string) bool {
	switch op {
	case ConditionStringEquals:
		return contextValue == policyValue
	case ConditionStringNotEquals:
		return contextValue != policyValue
	case ConditionStringEqualsIgnoreCase:
		return strings.EqualFold(contextValue, policyValue)
	case ConditionStringNotEqualsIgnoreCase:
		return !strings.EqualFold(contextValue, policyValue)
	case ConditionStringLike:
		return matchLikePattern(policyValue, contextValue)
	case ConditionStringNotLike:
		return !matchLikePattern(policyValue, contextValue)
	case ConditionNumericEquals:
		return compareNumeric(contextValue, policyValue) == 0
	case ConditionNumericNotEquals:
		return compareNumeric(contextValue, policyValue) != 0
	case ConditionNumericLessThan:
		return compareNumeric(contextValue, policyValue) < 0
	case ConditionNumericLessThanEquals:
		return compareNumeric(contextValue, policyValue) <= 0
	case ConditionNumericGreaterThan:
		return compareNumeric(contextValue, policyValue) > 0
	case ConditionNumericGreaterThanEquals:
		return compareNumeric(contextValue, policyValue) >= 0
	case ConditionDateEquals:
		return compareDates(contextValue, policyValue) == 0
	case ConditionDateNotEquals:
		return compareDates(contextValue, policyValue) != 0
	case ConditionDateLessThan:
		return compareDates(contextValue, policyValue) < 0
	case ConditionDateLessThanEquals:
		return compareDates(contextValue, policyValue) <= 0
	case ConditionDateGreaterThan:
		return compareDates(contextValue, policyValue) > 0
	case ConditionDateGreaterThanEquals:
		return compareDates(contextValue, policyValue) >= 0
	case ConditionBool:
		return strings.EqualFold(contextValue, policyValue)
	case ConditionIpAddress:
		return matchIPAddress(contextValue, policyValue)
	case ConditionNotIpAddress:
		return !matchIPAddress(contextValue, policyValue)
	case ConditionArnEquals:
		return contextValue == policyValue
	case ConditionArnLike:
		return matchArnLike(policyValue, contextValue)
	case ConditionArnNotEquals:
		return contextValue != policyValue
	case ConditionArnNotLike:
		return !matchArnLike(policyValue, contextValue)
	case ConditionBinaryEquals:
		return contextValue == policyValue
	case ConditionNull:
		return false
	default:
		return false
	}
}

func matchLikePattern(pattern, value string) bool {
	return wildcardMatch(pattern, value)
}

func wildcardMatch(pattern, value string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == value
	}

	if !strings.HasPrefix(value, parts[0]) {
		return false
	}
	value = value[len(parts[0]):]

	for i := 1; i < len(parts)-1; i++ {
		idx := strings.Index(value, parts[i])
		if idx < 0 {
			return false
		}
		value = value[idx+len(parts[i]):]
	}

	if len(parts) > 0 && parts[len(parts)-1] != "" {
		return strings.HasSuffix(value, parts[len(parts)-1])
	}
	return true
}

func matchArnLike(pattern, arn string) bool {
	segments := strings.Split(pattern, ":")
	arnSegments := strings.Split(arn, ":")

	if len(segments) != len(arnSegments) {
		return false
	}

	for i, seg := range segments {
		if seg == "*" {
			continue
		}
		if strings.Contains(seg, "*") {
			parts := strings.Split(seg, "*")
			if !strings.HasPrefix(arnSegments[i], parts[0]) {
				return false
			}
			if len(parts) > 1 && parts[1] != "" && !strings.HasSuffix(arnSegments[i], parts[1]) {
				return false
			}
		} else if seg != arnSegments[i] {
			return false
		}
	}
	return true
}

func matchIPAddress(value, cidr string) bool {
	ip := net.ParseIP(value)
	if ip == nil {
		return false
	}

	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return value == cidr
	}

	return network.Contains(ip)
}

func compareNumeric(a, b string) int {
	aFloat := parseFloat(a)
	bFloat := parseFloat(b)
	if aFloat < bFloat {
		return -1
	} else if aFloat > bFloat {
		return 1
	}
	return 0
}

func parseFloat(s string) float64 {
	var result float64
	var sign float64 = 1
	var decimal float64 = 1
	var inDecimal bool

	for _, c := range s {
		if c == '-' {
			sign = -1
		} else if c == '.' {
			inDecimal = true
		} else if c >= '0' && c <= '9' {
			digit := float64(c - '0')
			if inDecimal {
				decimal *= 10
				result += digit / decimal
			} else {
				result = result*10 + digit
			}
		}
	}
	return result * sign
}

func compareDates(a, b string) int {
	timeA, errA := parseDate(a)
	timeB, errB := parseDate(b)

	if errA != nil || errB != nil {
		return strings.Compare(a, b)
	}

	if timeA.Before(timeB) {
		return -1
	} else if timeA.After(timeB) {
		return 1
	}
	return 0
}

func parseDate(s string) (time.Time, error) {
	formats := []string{
		timeutils.ISO8601SimpleFormat,
		timeutils.ISO8601Format,
		timeutils.ISO8601WithTimezoneFormat,
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Parse(time.RFC3339, s)
}

func filepathMatch(pattern, name string) bool {
	matched, _ := filepath.Match(pattern, name)
	return matched
}

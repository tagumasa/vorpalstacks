package eventbridge

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func (s *EventsService) matchEventPattern(event *eventsstore.Event, pattern string) bool {
	var patternMap map[string]interface{}
	if err := json.Unmarshal([]byte(pattern), &patternMap); err != nil {
		return false
	}

	eventMap := map[string]interface{}{
		"version":     event.Version,
		"id":          event.ID,
		"source":      event.Source,
		"detail-type": event.DetailType,
		"time":        event.Time.Format(time.RFC3339),
		"region":      event.Region,
		"resources":   event.Resources,
		"detail":      event.Detail,
		"account":     event.Account,
	}

	for key, patternValue := range patternMap {
		eventValue, exists := eventMap[key]
		if !exists {
			if isExistsFalsePattern(patternValue) {
				continue
			}
			return false
		}

		if !s.matchValue(eventValue, patternValue) {
			return false
		}
	}

	return true
}

// isExistsFalsePattern checks whether a pattern value is equivalent to
// {"exists": false} (either directly or as a single-element array),
// which should match when the field is absent from the event.
func isExistsFalsePattern(patternValue interface{}) bool {
	check := func(obj map[string]interface{}) bool {
		if len(obj) != 1 {
			return false
		}
		if op, ok := obj["exists"]; ok {
			if b, ok := op.(bool); ok && !b {
				return true
			}
		}
		return false
	}
	switch p := patternValue.(type) {
	case map[string]interface{}:
		return check(p)
	case []interface{}:
		if len(p) == 1 {
			if item, ok := p[0].(map[string]interface{}); ok {
				return check(item)
			}
		}
	}
	return false
}

func (s *EventsService) matchValue(eventValue, patternValue interface{}) bool {
	switch p := patternValue.(type) {
	case []interface{}:
		for _, item := range p {
			if s.matchValue(eventValue, item) {
				return true
			}
		}
		return false
	case string:
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		return evStr == p
	case map[string]interface{}:
		if len(p) == 1 {
			for key, operand := range p {
				if isKnownOperator(key) {
					return s.matchOperator(eventValue, key, operand)
				}
				break
			}
		}
		evMap, ok := eventValue.(map[string]interface{})
		if !ok {
			return false
		}
		for k, v := range p {
			if !s.matchValue(evMap[k], v) {
				return false
			}
		}
		return true
	default:
		return fmt.Sprintf("%v", eventValue) == fmt.Sprintf("%v", patternValue)
	}
}

func (s *EventsService) matchOperator(eventValue interface{}, op string, operand interface{}) bool {
	switch op {
	case "prefix":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		prefix, ok := operand.(string)
		if !ok {
			return false
		}
		return strings.HasPrefix(evStr, prefix)
	case "numeric":
		evNum, ok := toFloat64(eventValue)
		if !ok {
			return false
		}
		operands, ok := operand.([]interface{})
		if !ok || len(operands) < 2 {
			return false
		}
		for i := 0; i < len(operands)-1; i++ {
			compOp, ok := operands[i].(string)
			if !ok {
				return false
			}
			compVal, ok := toFloat64(operands[i+1])
			if !ok {
				return false
			}
			if !compareNumeric(evNum, compOp, compVal) {
				return false
			}
			i++
		}
		return true
	case "anything-but":
		return !s.matchValue(eventValue, operand)
	case "exists":
		existsVal, ok := operand.(bool)
		if !ok {
			return false
		}
		if existsVal {
			return eventValue != nil
		}
		return eventValue == nil
	case "suffix":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		suffix, ok := operand.(string)
		if !ok {
			return false
		}
		return strings.HasSuffix(evStr, suffix)
	case "equals-ignore-case":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		operandStr, ok := operand.(string)
		if !ok {
			return false
		}
		return strings.EqualFold(evStr, operandStr)
	case "wildcard":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		pattern, ok := operand.(string)
		if !ok {
			return false
		}
		return matchWildcardPattern(evStr, pattern)
	case "cidr":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		cidr, ok := operand.(string)
		if !ok {
			return false
		}
		return matchCIDRBlock(evStr, cidr)
	default:
		return false
	}
}

// matchWildcardPattern reports whether s matches the glob pattern, where
// '*' matches any sequence of characters (including none). The match is
// anchored at both ends: a pattern without a leading '*' must match from
// the start of s, and a pattern without a trailing '*' must match to the
// end of s. Interior segments may occur anywhere in order.
func matchWildcardPattern(s, pattern string) bool {
	fragments := strings.Split(pattern, "*")

	// A pattern without any wildcard must match the whole string.
	if len(fragments) == 1 {
		return s == fragments[0]
	}

	// The leading fragment is anchored at the start of s.
	head := fragments[0]
	if !strings.HasPrefix(s, head) {
		return false
	}
	s = s[len(head):]

	// The trailing fragment is anchored at the end of s.
	tail := fragments[len(fragments)-1]
	if !strings.HasSuffix(s, tail) {
		return false
	}
	s = s[:len(s)-len(tail)]

	// Interior fragments must appear in order in what remains.
	for _, frag := range fragments[1 : len(fragments)-1] {
		idx := strings.Index(s, frag)
		if idx == -1 {
			return false
		}
		s = s[idx+len(frag):]
	}
	return true
}

func matchCIDRBlock(ipStr, cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ipNet.Contains(ip)
}

func isKnownOperator(key string) bool {
	switch key {
	case "prefix", "numeric", "anything-but", "exists", "suffix", "equals-ignore-case", "wildcard", "cidr":
		return true
	default:
		return false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

func compareNumeric(val float64, op string, comp float64) bool {
	switch op {
	case "<":
		return val < comp
	case "<=":
		return val <= comp
	case ">":
		return val > comp
	case ">=":
		return val >= comp
	case "=":
		return val == comp
	default:
		return false
	}
}

package eventbridge

import (
	"encoding/json"
	"net"
	"strings"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func (s *EventsService) matchEventPattern(event *eventsstore.Event, pattern string) bool {
	var patternMap map[string]interface{}
	if err := json.Unmarshal([]byte(pattern), &patternMap); err != nil {
		return false
	}

	return s.matchPatternLevel(eventEnvelope(event), patternMap)
}

// matchPatternLevel evaluates one object level of an event pattern against
// the corresponding event map: every pattern key must be satisfied (AND),
// and "$or" is satisfied when any of its entry objects matches. The loop is
// shared by rule matching, archive ingest matching and TestEventPattern so
// the top-level semantics cannot drift between them. A key absent from the
// event only satisfies an explicit exists-false pattern.
func (s *EventsService) matchPatternLevel(eventMap map[string]interface{}, patternMap map[string]interface{}) bool {
	for key, patternValue := range patternMap {
		if key == "$or" {
			entries, ok := patternValue.([]interface{})
			if !ok {
				return false
			}
			matched := false
			for _, entry := range entries {
				entryMap, ok := entry.(map[string]interface{})
				if !ok {
					return false
				}
				if s.matchPatternLevel(eventMap, entryMap) {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
			continue
		}

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
	// An array-valued event member matches when any of its elements
	// matches — the AWS intersection rule: "If the value in the event is an
	// array, then the event pattern matches if the intersection of the
	// event pattern array and the event array is non-empty."
	switch elements := eventValue.(type) {
	case []interface{}:
		for _, element := range elements {
			if s.matchValue(element, patternValue) {
				return true
			}
		}
		return false
	case []string:
		for _, element := range elements {
			if s.matchValue(element, patternValue) {
				return true
			}
		}
		return false
	}

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
		return s.matchPatternLevel(evMap, p)
	default:
		return scalarEqual(eventValue, patternValue)
	}
}

// scalarEqual compares the non-string, non-object pattern scalars with their
// JSON type: "The values that event patterns match follow JSON rules. You
// can include strings enclosed in quotation marks ("), numbers, and the
// keywords true, false, and null" (event pattern syntax, Considerations) —
// numbers compare numerically, booleans and null by identity, and a value of
// one JSON type never matches a value of another.
func scalarEqual(eventValue, patternValue interface{}) bool {
	switch p := patternValue.(type) {
	case float64:
		ev, ok := toFloat64(eventValue)
		if !ok {
			return false
		}
		return ev == p
	case bool:
		ev, ok := eventValue.(bool)
		return ok && ev == p
	case nil:
		return eventValue == nil
	default:
		return false
	}
}

func (s *EventsService) matchOperator(eventValue interface{}, op string, operand interface{}) bool {
	switch op {
	case "prefix":
		evStr, ok := eventValue.(string)
		if !ok {
			return false
		}
		target, ignoreCase, ok := prefixSuffixOperand(operand)
		if !ok {
			return false
		}
		if ignoreCase {
			return strings.HasPrefix(strings.ToLower(evStr), strings.ToLower(target))
		}
		return strings.HasPrefix(evStr, target)
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
		return !s.matchAnythingBut(eventValue, operand)
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
		target, ignoreCase, ok := prefixSuffixOperand(operand)
		if !ok {
			return false
		}
		if ignoreCase {
			return strings.HasSuffix(strings.ToLower(evStr), strings.ToLower(target))
		}
		return strings.HasSuffix(evStr, target)
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

// prefixSuffixOperand expands a prefix or suffix operand: AWS documents a
// plain string and the nested {"equals-ignore-case": "<string>"}
// conjunction ("using equals-ignore-case in conjunction with prefix").
func prefixSuffixOperand(operand interface{}) (target string, ignoreCase bool, ok bool) {
	switch v := operand.(type) {
	case string:
		return v, false, true
	case map[string]interface{}:
		if len(v) == 1 {
			if inner, isStr := v["equals-ignore-case"].(string); isStr {
				return inner, true, true
			}
		}
	}
	return "", false, false
}

// matchAnythingBut reports whether the event value is excluded by an
// anything-but operand: AWS documents scalars, lists of only strings or
// only numbers, and the nested prefix/suffix/equals-ignore-case/wildcard
// conjunctions, each taking a string or a list of strings.
func (s *EventsService) matchAnythingBut(eventValue interface{}, operand interface{}) bool {
	if conjunction, ok := operand.(map[string]interface{}); ok {
		if len(conjunction) != 1 {
			return false
		}
		for op, inner := range conjunction {
			switch op {
			case "prefix", "suffix", "equals-ignore-case", "wildcard":
				list, ok := stringListOperand(inner)
				if !ok {
					return false
				}
				for _, item := range list {
					if s.matchOperator(eventValue, op, item) {
						return true
					}
				}
				return false
			default:
				return false
			}
		}
	}
	return s.matchValue(eventValue, operand)
}

// stringListOperand expands a string-or-list-of-strings operand.
func stringListOperand(operand interface{}) ([]string, bool) {
	switch v := operand.(type) {
	case string:
		return []string{v}, true
	case []interface{}:
		list := make([]string, 0, len(v))
		for _, item := range v {
			str, ok := item.(string)
			if !ok {
				return nil, false
			}
			list = append(list, str)
		}
		return list, true
	}
	return nil, false
}

// matchWildcardPattern reports whether s matches the glob pattern, where
// '*' matches any sequence of characters (including none) and a backslash
// escapes '*' and '\' into literals. The match is anchored at both ends: a
// pattern without a leading '*' must match from the start of s, and a
// pattern without a trailing '*' must match to the end of s. Interior
// segments may occur anywhere in order.
func matchWildcardPattern(s, pattern string) bool {
	fragments, ok := splitWildcardPattern(pattern)
	if !ok {
		return false
	}

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

// splitWildcardPattern splits a wildcard pattern into the literal fragments
// between unescaped '*' characters. "\*" and "\\" are the literals '*' and
// '\'; a backslash preceding any other character, or ending the pattern, is
// not a documented escape and reports ok=false ("Using the backslash to
// escape other characters is not supported").
func splitWildcardPattern(pattern string) ([]string, bool) {
	fragments := []string{}
	current := strings.Builder{}
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			fragments = append(fragments, current.String())
			current.Reset()
		case '\\':
			if i+1 >= len(pattern) || (pattern[i+1] != '*' && pattern[i+1] != '\\') {
				return nil, false
			}
			i++
			current.WriteByte(pattern[i])
		default:
			current.WriteByte(pattern[i])
		}
	}
	fragments = append(fragments, current.String())
	return fragments, true
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

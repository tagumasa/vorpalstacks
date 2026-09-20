package sns

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// Subscription filter policies: the documented value grammar, its
// set-time validation, and the delivery-time matcher. Validation and
// matching live in one file on purpose — the incident this file answers
// found a validator that accepted operand shapes the matcher silently
// mis-evaluated (a non-string prefix operand became "" and matched
// everything; an unknown anything-but operand matched everything). Every
// operator's operand rules are validated here to exactly the shapes
// applyFilterOperator evaluates, and the matcher itself stays fail-closed:
// an operand shape it cannot evaluate never matches.
//
// Operator set and semantics follow the SNS developer guide's message
// filtering pages (string value matching, numeric value matching, filter
// policy constraints), whose sentences the rule comments quote. Nested
// policies are payload-based filtering alone: "Amazon SNS accepts a nested
// filter policy for payload-based filtering" while attribute-based
// filtering "doesn't accept a nested filter policy" — payload-scope
// policies may nest JSON-object values whose leaf keys the constraints
// page counts and whose nested level its combination formula multiplies
// in, and the matcher walks the body's nested structure along the same
// paths.

// matchFilterPolicy evaluates whether a message satisfies the subscription's
// filter policy, respecting the FilterPolicyScope setting.
//
// AWS SNS filter policy semantics:
//   - Each key in the policy must match (AND semantics across keys).
//   - Multiple values for a key use OR semantics.
//   - A missing message attribute fails the key unless a value entry
//     matches the absence ({"exists": false}).
//
// When scope is "MessageBody", the policy keys refer to properties of the
// JSON message body instead of message attributes, and the policy may nest
// object levels — "Amazon SNS accepts a nested filter policy for
// payload-based filtering" — which matchPolicyNode walks against the
// body's nested structure.
func matchFilterPolicy(policyJSON string, scope string, attrs map[string]*snsstore.MessageAttribute, messageBody string) bool {
	if strings.TrimSpace(policyJSON) == "" || policyJSON == "{}" {
		return true
	}

	var policy map[string]json.RawMessage
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		return false
	}

	if scope != snsstore.FilterPolicyScopeBody {
		// Attribute-based matching stays flat: "Amazon SNS doesn't accept
		// a nested filter policy for attribute-based filtering" — a nested
		// value shape unmarshals as no value array and fails its key.
		for attrName, rawValues := range policy {
			var values []interface{}
			if err := json.Unmarshal(rawValues, &values); err != nil {
				return false
			}

			msgAttr, attrExists := attrs[attrName]
			if !matchPolicyValues(values, msgAttr, attrExists) {
				return false
			}
		}
		return true
	}

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(messageBody), &body); err != nil {
		// An unparseable body exposes no properties to match against;
		// only exists:false leaves can still hold.
		body = nil
	}
	return matchPolicyNode(policy, body)
}

// matchPolicyNode evaluates one object level of the policy against the
// corresponding level of the message body. A policy key's value is either
// a leaf value array — matched against the body property through
// matchPolicyValues — or a nested object walked against the body
// property's object ("The policy property values are nested JSON
// objects", payload-based filtering). A body path that is absent, or that
// holds a scalar or array where the policy nests, leaves the nested level
// without properties: every leaf below it is absent, so the level holds
// only through an exists:false leaf.
func matchPolicyNode(policy map[string]json.RawMessage, body map[string]interface{}) bool {
	for key, raw := range policy {
		if strings.HasPrefix(strings.TrimSpace(string(raw)), "{") {
			var nestedPolicy map[string]json.RawMessage
			if err := json.Unmarshal(raw, &nestedPolicy); err != nil || len(nestedPolicy) == 0 {
				return false
			}
			var nestedBody map[string]interface{}
			if body != nil {
				nestedBody, _ = body[key].(map[string]interface{})
			}
			if !matchPolicyNode(nestedPolicy, nestedBody) {
				return false
			}
			continue
		}

		var values []interface{}
		if err := json.Unmarshal(raw, &values); err != nil {
			return false
		}
		var msgAttr *snsstore.MessageAttribute
		attrExists := false
		if body != nil {
			if v, ok := body[key]; ok {
				attrExists = true
				msgAttr = bodyValueToAttribute(v)
			}
		}
		if !matchPolicyValues(values, msgAttr, attrExists) {
			return false
		}
	}
	return true
}

// nullAttributeType marks a message-body property whose JSON value is null;
// message attributes themselves can never be null, so only the body scope
// produces this marker. A policy's null exact value matches it.
const nullAttributeType = "Null"

// bodyValueToAttribute converts one JSON message-body property into the
// pseudo MessageAttribute form the value matcher consumes: scalars map to
// String/Number (booleans render as their JSON literal, matching the
// boolean policy values attribute matching accepts), a null property
// carries the null marker a policy's null value matches, and arrays and
// objects ride as their JSON text under the String.Array type, whose
// element matching then treats as separate values.
func bodyValueToAttribute(v interface{}) *snsstore.MessageAttribute {
	switch val := v.(type) {
	case string:
		return &snsstore.MessageAttribute{Type: "String", StringValue: val}
	case float64:
		return &snsstore.MessageAttribute{Type: "Number", StringValue: strconv.FormatFloat(val, 'f', -1, 64)}
	case bool:
		return &snsstore.MessageAttribute{Type: "String", StringValue: strconv.FormatBool(val)}
	case nil:
		return &snsstore.MessageAttribute{Type: nullAttributeType}
	default:
		raw, err := json.Marshal(val)
		if err != nil {
			return nil
		}
		return &snsstore.MessageAttribute{Type: "String.Array", StringValue: string(raw)}
	}
}

// filterCandidates expands one message attribute into the scalar string
// values matching considers. A String or Number attribute is its single
// value (numeric operators parse it, and the documented exponent notation
// 3.015e2 survives as text); for String.Array and Number.Array attributes
// each element "is treated as a separate string for matching purposes"
// (string value matching, numeric value matching). Binary attributes are
// ignored — "Amazon SNS ignores message attributes with the Binary data
// type" (filter policy constraints).
func filterCandidates(attr *snsstore.MessageAttribute) []string {
	if attr == nil || attr.Type == "Binary" {
		return nil
	}
	if !strings.HasSuffix(attr.Type, ".Array") {
		return []string{attr.StringValue}
	}
	var elements []interface{}
	if err := json.Unmarshal([]byte(attr.StringValue), &elements); err != nil {
		// An unparseable array value has no elements to match.
		return nil
	}
	candidates := make([]string, 0, len(elements))
	for _, el := range elements {
		switch v := el.(type) {
		case string:
			candidates = append(candidates, v)
		case float64:
			candidates = append(candidates, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			candidates = append(candidates, strconv.FormatBool(v))
		}
	}
	return candidates
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
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return c == val })
	case bool:
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return c == strconv.FormatBool(val) })
	case float64:
		// Numeric exact matching "matches any message attribute or message
		// body property values that have the same name and an equal numeric
		// value" — 301.5 matches an attribute carrying 3.015e2 — so the
		// comparison parses both sides instead of comparing text.
		return attrExists && anyCandidate(msgAttr, func(c string) bool {
			parsed, err := strconv.ParseFloat(c, 64)
			return err == nil && parsed == val
		})
	case nil:
		// The null keyword matches a body property whose value is null.
		return attrExists && msgAttr != nil && msgAttr.Type == nullAttributeType
	case map[string]interface{}:
		if len(val) != 1 {
			return false
		}
		for operator, operand := range val {
			if !applyFilterOperator(operator, operand, msgAttr, attrExists) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// anyCandidate reports whether any candidate value of the attribute
// satisfies pred.
func anyCandidate(msgAttr *snsstore.MessageAttribute, pred func(string) bool) bool {
	for _, c := range filterCandidates(msgAttr) {
		if pred(c) {
			return true
		}
	}
	return false
}

// applyFilterOperator dispatches a single filter operator. Unknown
// operators and operand shapes the operator cannot evaluate return false —
// a filter policy never matches everything by accident of a type fallback.
func applyFilterOperator(operator string, operand interface{}, msgAttr *snsstore.MessageAttribute, attrExists bool) bool {
	switch operator {
	case "exists":
		want, ok := operand.(bool)
		if !ok {
			return false
		}
		if want {
			return attrExists
		}
		return !attrExists

	case "prefix":
		prefix, ok := operand.(string)
		if !ok {
			return false
		}
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return strings.HasPrefix(c, prefix) })

	case "suffix":
		suffix, ok := operand.(string)
		if !ok {
			return false
		}
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return strings.HasSuffix(c, suffix) })

	case "wildcard":
		pattern, ok := operand.(string)
		if !ok {
			return false
		}
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return wildcardMatch(pattern, c) })

	case "equals-ignore-case":
		want, ok := operand.(string)
		if !ok {
			return false
		}
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return strings.EqualFold(c, want) })

	case "cidr":
		subnet, ok := operand.(string)
		if !ok {
			return false
		}
		return attrExists && anyCandidate(msgAttr, func(c string) bool { return cidrContains(subnet, c) })

	case "anything-but":
		if !attrExists {
			return false
		}
		// An array attribute is excluded only when EVERY element is
		// excluded: ["rugby", "baseball"] matches anything-but
		// ["rugby", "tennis"] "because it contains a value that isn't
		// rugby or tennis", while ["rugby"] does not match (string value
		// matching, numeric value matching).
		candidates := filterCandidates(msgAttr)
		if len(candidates) == 0 {
			return false
		}
		for _, c := range candidates {
			excluded, evaluable := anythingButOperandExcludes(operand, c)
			if !evaluable {
				// Fail-closed: an operand shape the matcher cannot
				// evaluate withholds the message instead of matching.
				return false
			}
			if !excluded {
				// The attribute contains a value that isn't excluded.
				return true
			}
		}
		return false

	case "numeric":
		if !attrExists {
			return false
		}
		return numericRangeHolds(operand, msgAttr)

	default:
		return false
	}
}

// anythingButOperandExcludes reports whether one candidate value is
// excluded by the anything-but operand, and whether the operand shape is
// evaluable at all. Documented operand forms: a string, a number, a list
// of those, and one nested prefix/suffix/wildcard pattern (string value
// matching, numeric value matching). Set-time validation rejects every
// other shape, so the unevaluable branch is unreachable except for
// records that predate validation.
func anythingButOperandExcludes(operand interface{}, candidate string) (excluded, evaluable bool) {
	switch v := operand.(type) {
	case string:
		return candidate == v, true
	case float64:
		parsed, err := strconv.ParseFloat(candidate, 64)
		if err != nil {
			return false, true
		}
		return parsed == v, true
	case []interface{}:
		for _, item := range v {
			itemExcludes, ok := anythingButOperandExcludes(item, candidate)
			if !ok {
				return false, false
			}
			if itemExcludes {
				return true, true
			}
		}
		return false, true
	case map[string]interface{}:
		if len(v) != 1 {
			return false, false
		}
		for operator, nested := range v {
			switch operator {
			case "prefix":
				prefix, ok := nested.(string)
				if !ok {
					return false, false
				}
				return strings.HasPrefix(candidate, prefix), true
			case "suffix":
				suffix, ok := nested.(string)
				if !ok {
					return false, false
				}
				return strings.HasSuffix(candidate, suffix), true
			case "wildcard":
				pattern, ok := nested.(string)
				if !ok {
					return false, false
				}
				return wildcardMatch(pattern, candidate), true
			default:
				return false, false
			}
		}
		return false, false
	default:
		return false, false
	}
}

// numericRangeHolds evaluates {"numeric": [">=", 0, "<", 100]} style
// ranges: an even-length array of alternating operator and bound, every
// condition must hold for the attribute's numeric value.
func numericRangeHolds(operand interface{}, msgAttr *snsstore.MessageAttribute) bool {
	conditions, ok := operand.([]interface{})
	if !ok || len(conditions) < 2 || len(conditions)%2 != 0 {
		return false
	}
	return anyCandidate(msgAttr, func(c string) bool {
		target, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return false
		}
		for i := 0; i < len(conditions); i += 2 {
			op, opOk := conditions[i].(string)
			bound, boundOk := conditions[i+1].(float64)
			if !opOk || !boundOk {
				return false
			}
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
	})
}

// wildcardMatch reports whether s matches the pattern, where '*' matches
// any sequence of characters (including empty). All other characters match
// themselves.
func wildcardMatch(pattern, s string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == s
	}
	parts := strings.Split(pattern, "*")
	if !strings.HasPrefix(s, parts[0]) {
		return false
	}
	s = s[len(parts[0]):]
	last := parts[len(parts)-1]
	if len(s) < len(last) || !strings.HasSuffix(s, last) {
		return false
	}
	rest := s[:len(s)-len(last)]
	for _, mid := range parts[1 : len(parts)-1] {
		idx := strings.Index(rest, mid)
		if idx < 0 {
			return false
		}
		rest = rest[idx+len(mid):]
	}
	return true
}

// cidrContains reports whether the candidate IP address falls inside the
// policy's subnet ("the cidr operator ... check[s] whether an incoming
// message originates from a specific IP address or subnet"). A policy
// operand naming a single address (no mask) matches that address alone.
func cidrContains(subnet, candidate string) bool {
	ip := net.ParseIP(candidate)
	if ip == nil {
		return false
	}
	if strings.Contains(subnet, "/") {
		_, network, err := net.ParseCIDR(subnet)
		if err != nil {
			return false
		}
		return network.Contains(ip)
	}
	single := net.ParseIP(subnet)
	return single != nil && single.Equal(ip)
}

// ---------------------------------------------------------------------------
// Filter-policy validation — the set-time half of the shared shape
// contract: everything the matcher above can evaluate, and nothing it
// cannot.
// ---------------------------------------------------------------------------

// validNumericOperators lists the numeric range operators: "a numeric
// policy property can include the following operators: <, <=, >, and >="
// beside "=" (numeric value matching).
var validNumericOperators = map[string]bool{
	"=": true, "<": true, "<=": true, ">": true, ">=": true,
}

// validAnythingButNestedOperators lists the operators an anything-but
// operand may nest: "you can also use a prefix with the anything-but
// operator" (string value matching), with suffix and wildcard carried by
// the same anything-but section.
var validAnythingButNestedOperators = map[string]bool{
	"prefix": true, "suffix": true, "wildcard": true,
}

// filterPolicyStats accumulates the constraints-page counters computed
// across a policy's leaf keys.
type filterPolicyStats struct {
	leaves             int
	combination        int
	wildcardComplexity int
}

// validateFilterPolicyDocument validates the JSON structure and every
// operand of a filter policy against the documented constraints: at most
// five keys ("For a nested policy, only leaf keys are counted towards the
// five key limit"), a total value combination of at most 150, at most
// 256 KB, and the operator grammar of the string/numeric value matching
// pages. Payload-scope policies (FilterPolicyScope MessageBody) may carry
// nested JSON-object values; attribute-scope policies may not. The policy
// is accepted only when the matcher can evaluate every entry, so an
// invalid policy can never behave as match-all at delivery time.
func validateFilterPolicyDocument(value, scope string) error {
	if value == "" {
		return nil
	}
	if len(value) > snsstore.MaxFilterPolicySizeBytes {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid filter policy: size %d bytes exceeds the maximum of %d bytes", len(value), snsstore.MaxFilterPolicySizeBytes))
	}

	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(value), &policy); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid filter policy: %s", err.Error()))
	}

	stats := filterPolicyStats{combination: 1}
	return validateFilterPolicyLevel(policy, 1, scope == snsstore.FilterPolicyScopeBody, &stats)
}

// validateFilterPolicyLevel validates one object level of a filter policy.
// depth is the nested level of this level's keys, counting from one at the
// top of the policy: the constraints page's payload-based combination
// example multiplies each leaf's match-operator count "and the nested
// level for each key" (key_c at a nested level of three contributes 4 x 3,
// key_e at level two 3 x 2, totalling 72). Payload-scope policies may
// nest object levels ("Amazon SNS accepts a nested filter policy for
// payload-based filtering"); attribute-scope policies may not ("Amazon SNS
// doesn't accept a nested filter policy for attribute-based filtering").
func validateFilterPolicyLevel(level map[string]interface{}, depth int, allowNesting bool, stats *filterPolicyStats) error {
	for attrName, raw := range level {
		if attrName == "" {
			return NewInvalidParameter("Invalid filter policy: attribute name must not be empty")
		}
		switch values := raw.(type) {
		case []interface{}:
			stats.leaves++
			if stats.leaves > snsstore.MaxFilterPolicyKeys {
				return NewInvalidParameter(fmt.Sprintf(
					"Invalid filter policy: %d leaf keys exceeds the maximum of %d", stats.leaves, snsstore.MaxFilterPolicyKeys))
			}
			points, patterns := 0, 0
			for _, entry := range values {
				entryPoints, err := validateFilterPolicyEntry(entry)
				if err != nil {
					return err
				}
				if entryPoints > 0 {
					points += entryPoints
					patterns++
				}
			}
			// "Field complexity = (Sum of pattern points) × (Number of
			// patterns)" with "single wildcard: 1 point", "multiple wildcards:
			// 3 points each" and "anything-but: 1 point" (wildcard pattern
			// usage guidelines).
			if patterns > 0 {
				stats.wildcardComplexity += points * patterns
				if stats.wildcardComplexity > snsstore.MaxWildcardPatternComplexity {
					return NewInvalidParameter(fmt.Sprintf(
						"Invalid filter policy: wildcard complexity %d exceeds the maximum of %d points",
						stats.wildcardComplexity, snsstore.MaxWildcardPatternComplexity))
				}
			}
			if len(values) > 0 {
				stats.combination *= len(values) * depth
				if stats.combination > snsstore.MaxFilterPolicyCombination {
					return NewInvalidParameter(fmt.Sprintf(
						"Invalid filter policy: total combination of values exceeds the maximum of %d",
						snsstore.MaxFilterPolicyCombination))
				}
			}
		case map[string]interface{}:
			if !allowNesting {
				return NewInvalidParameter(
					"Invalid filter policy: a nested filter policy requires FilterPolicyScope MessageBody")
			}
			if len(values) == 0 {
				return NewInvalidParameter(
					"Invalid filter policy: a nested object must name at least one property")
			}
			if err := validateFilterPolicyLevel(values, depth+1, allowNesting, stats); err != nil {
				return err
			}
		default:
			return NewInvalidParameter("Invalid filter policy: value must be an array")
		}
	}
	return nil
}

// validateFilterPolicyEntry validates one array entry of a policy value
// and reports its wildcard-complexity points (0 for entries the guidelines
// do not count). Scalars follow the constraints page: "The JSON of the
// filter policy can contain the following: Strings enclosed in quotation
// marks; Numbers; The keywords true, false, and null". Condition objects
// carry exactly one operator whose operand shape the matcher evaluates.
func validateFilterPolicyEntry(entry interface{}) (int, error) {
	switch cond := entry.(type) {
	case string:
		return 0, nil
	case bool:
		return 0, nil
	case nil:
		return 0, nil
	case float64:
		return 0, validateFilterNumericValue(cond)
	case map[string]interface{}:
		if len(cond) != 1 {
			return 0, NewInvalidParameter(
				"Invalid filter policy: a condition object must have exactly one operator")
		}
		for operator, operand := range cond {
			return validateFilterOperator(operator, operand)
		}
		return 0, nil
	default:
		return 0, NewInvalidParameter("Invalid filter policy: unsupported value type")
	}
}

// validateFilterOperator validates one operator's operand and reports its
// wildcard-complexity points.
func validateFilterOperator(operator string, operand interface{}) (int, error) {
	switch operator {
	case "prefix", "suffix", "equals-ignore-case":
		// "Prefixes are supported for string matching only" (numeric value
		// matching); the sibling string operators take string operands.
		if _, ok := operand.(string); !ok {
			return 0, NewInvalidParameter(fmt.Sprintf(
				"Invalid filter policy: %s requires a string value", operator))
		}
		return 0, nil

	case "wildcard":
		pattern, ok := operand.(string)
		if !ok {
			return 0, NewInvalidParameter("Invalid filter policy: wildcard requires a string value")
		}
		wildcards := strings.Count(pattern, "*")
		if wildcards > snsstore.MaxWildcardFiltersPerPattern {
			return 0, NewInvalidParameter(fmt.Sprintf(
				"Invalid filter policy: wildcard pattern %q exceeds the maximum of %d wildcards",
				pattern, snsstore.MaxWildcardFiltersPerPattern))
		}
		switch {
		case wildcards == 0:
			// No wildcard character: the pattern is an exact match and the
			// guidelines' points do not apply.
			return 0, nil
		case wildcards == 1:
			// "Single wildcard: 1 point".
			return 1, nil
		default:
			// "Multiple wildcards: 3 points each".
			return 3 * wildcards, nil
		}

	case "cidr":
		subnet, ok := operand.(string)
		if !ok {
			return 0, NewInvalidParameter("Invalid filter policy: cidr requires a string value")
		}
		if strings.Contains(subnet, "/") {
			if _, _, err := net.ParseCIDR(subnet); err != nil {
				return 0, NewInvalidParameter(fmt.Sprintf(
					"Invalid filter policy: invalid CIDR notation %q", subnet))
			}
			return 0, nil
		}
		if net.ParseIP(subnet) == nil {
			return 0, NewInvalidParameter(fmt.Sprintf(
				"Invalid filter policy: invalid IP address %q", subnet))
		}
		return 0, nil

	case "exists":
		if _, ok := operand.(bool); !ok {
			return 0, NewInvalidParameter("Invalid filter policy: exists requires true or false")
		}
		return 0, nil

	case "anything-but":
		points, err := validateAnythingButOperand(operand)
		if err != nil {
			return 0, err
		}
		// "Anything-but: 1 point" — the flat cost the guidelines assign.
		return points + 1, nil

	case "numeric":
		conditions, ok := operand.([]interface{})
		if !ok || len(conditions) < 2 || len(conditions)%2 != 0 {
			return 0, NewInvalidParameter(
				"Invalid filter policy: numeric requires an even-length array of operator and value pairs")
		}
		for i := 0; i < len(conditions); i += 2 {
			op, opOk := conditions[i].(string)
			if !opOk || !validNumericOperators[op] {
				return 0, NewInvalidParameter(fmt.Sprintf(
					"Invalid filter policy: unknown numeric operator %v", conditions[i]))
			}
			bound, boundOk := conditions[i+1].(float64)
			if !boundOk {
				return 0, NewInvalidParameter(
					"Invalid filter policy: numeric bounds must be numbers")
			}
			if err := validateFilterNumericValue(bound); err != nil {
				return 0, err
			}
		}
		return 0, nil

	default:
		return 0, NewInvalidParameter(fmt.Sprintf(
			"Invalid filter policy: unknown operator %q", operator))
	}
}

// validateAnythingButOperand validates the documented anything-but operand
// forms — a string, a number, a list of those, or one nested
// prefix/suffix/wildcard pattern — and reports the nested wildcard points.
func validateAnythingButOperand(operand interface{}) (int, error) {
	switch v := operand.(type) {
	case string:
		return 0, nil
	case float64:
		return 0, validateFilterNumericValue(v)
	case []interface{}:
		if len(v) == 0 {
			return 0, NewInvalidParameter(
				"Invalid filter policy: anything-but list must not be empty")
		}
		for _, item := range v {
			switch item.(type) {
			case string:
			case float64:
				if err := validateFilterNumericValue(item.(float64)); err != nil {
					return 0, err
				}
			default:
				return 0, NewInvalidParameter(
					"Invalid filter policy: anything-but list items must be strings or numbers")
			}
		}
		return 0, nil
	case map[string]interface{}:
		if len(v) != 1 {
			return 0, NewInvalidParameter(
				"Invalid filter policy: an anything-but condition object must have exactly one operator")
		}
		for operator, nested := range v {
			if !validAnythingButNestedOperators[operator] {
				return 0, NewInvalidParameter(fmt.Sprintf(
					"Invalid filter policy: anything-but cannot nest operator %q", operator))
			}
			return validateFilterOperator(operator, nested)
		}
		return 0, nil
	default:
		return 0, NewInvalidParameter(
			"Invalid filter policy: anything-but requires a string, number, list, or a nested prefix/suffix/wildcard pattern")
	}
}

// validateFilterNumericValue enforces the documented numeric bounds: "the
// value can range from -10^9 to 10^9 (-1 billion to 1 billion), with five
// digits of accuracy after the decimal point" (filter policy constraints).
func validateFilterNumericValue(v float64) error {
	if v > snsstore.MaxFilterPolicyNumericValue || v < -snsstore.MaxFilterPolicyNumericValue {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid filter policy: numeric value %v is outside the supported range of -%d to %d",
			v, snsstore.MaxFilterPolicyNumericValue, snsstore.MaxFilterPolicyNumericValue))
	}
	scale := math.Pow10(snsstore.MaxFilterPolicyNumericAccuracyDigits)
	if v != math.Round(v*scale)/scale {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid filter policy: numeric value %v carries more than five digits of accuracy after the decimal point", v))
	}
	return nil
}

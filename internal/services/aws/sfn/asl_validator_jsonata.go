package sfn

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// jsonataResultRef matches a $states.result node reference inside an
// expression body (a trailing identifier character means a different
// name, e.g. $states.resultX).
var jsonataResultRef = regexp.MustCompile(`\$states\.result([^A-Za-z0-9_]|$)`)

// jsonataErrorOutputRef matches a $states.errorOutput node reference.
var jsonataErrorOutputRef = regexp.MustCompile(`\$states\.errorOutput([^A-Za-z0-9_]|$)`)

// validateJSONataExpressions enforces the two creation-time JSONata
// contracts on a state's expression-bearing values. First, expression
// strings must be properly delimited: "the string must start with `{%`
// with no leading spaces, and must end with `%}` with no trailing spaces.
// Improperly opening or closing the expression will result in a
// validation error." Second, "Attempting to access `$states.result` or
// `$states.errorOutput` in fields and states where they are not
// accessible will be caught at creation, update, or validation of the
// state machine": $states.result belongs to Task, Parallel and Map
// states, and $states.errorOutput is "used in the `Catch` field's
// `Assign` or `Output`".
func (v *aslValidatorContext) validateJSONataExpressions(name, stateType string, stateMap map[string]interface{}, stateLoc string) {
	hasResult := stateType == "Task" || stateType == "Parallel" || stateType == "Map"

	// Container members whose strings the engine evaluates as JSONata.
	// ItemSelector is the documented spelling the engine parses — an
	// underscore variant would be a dead walk entry.
	for _, field := range []string{
		"Output", "Arguments", "Assign", "Items",
		"ItemSelector", "ItemBatcher", "ItemReader", "ResultWriter",
		"Error", "Cause",
	} {
		if raw, ok := stateMap[field]; ok {
			v.walkJSONataValue(name, stateLoc+"/"+field, raw, hasResult, false)
		}
	}

	// Scalar expression members: a well-formed expression only needs the
	// accessibility check — an improperly delimited string already fails
	// the member's type validation.
	for _, field := range []string{
		"TimeoutSeconds", "HeartbeatSeconds", "Seconds", "Timestamp",
		"MaxConcurrency", "ToleratedFailureCount", "ToleratedFailurePercentage",
	} {
		if s, ok := stateMap[field].(string); ok && IsExpression(s) {
			v.checkJSONataAccessibility(name, stateLoc+"/"+field, s, hasResult, false)
		}
	}

	if choices, ok := stateMap["Choices"].([]interface{}); ok {
		for i, raw := range choices {
			rule, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			v.walkJSONataChoiceRule(name, stateLoc+"/Choices/"+strconv.Itoa(i), rule, false)
		}
	}

	if catches, ok := stateMap["Catch"].([]interface{}); ok {
		for i, raw := range catches {
			policy, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			loc := stateLoc + "/Catch/" + strconv.Itoa(i)
			// "The $states.errorOutput can be used in the Catch field's
			// Assign or Output" — the only two errorOutput-accessible
			// fields.
			for _, field := range []string{"Assign", "Output"} {
				if value, ok := policy[field]; ok {
					v.walkJSONataValue(name, loc+"/"+field, value, hasResult, true)
				}
			}
		}
	}
}

// walkJSONataChoiceRule validates one JSONata Choice rule and recurses the
// And/Or/Not combinators: "the value of an And or Or operator must be a
// non-empty array of Choice Rules" and "the value of a Not operator must
// be a single Choice Rule", and the nested rules carry Conditions and
// Assigns of their own — a walk that stopped at the top level left them
// outside both the delimiter and the accessibility checks.
func (v *aslValidatorContext) walkJSONataChoiceRule(name, loc string, rule map[string]interface{}, errorOutputAllowed bool) {
	// The Condition routes through the same walk as every other
	// engine-evaluated string: a well-formed expression gets the
	// accessibility check and an improperly delimited one the documented
	// validation error.
	if cond, ok := rule["Condition"].(string); ok {
		v.walkJSONataValue(name, loc+"/Condition", cond, false, errorOutputAllowed)
	}
	if assign, ok := rule["Assign"]; ok {
		v.walkJSONataValue(name, loc+"/Assign", assign, false, errorOutputAllowed)
	}
	for _, combinator := range []string{"And", "Or", "Not"} {
		switch nested := rule[combinator].(type) {
		case map[string]interface{}:
			v.walkJSONataChoiceRule(name, loc+"/"+combinator, nested, errorOutputAllowed)
		case []interface{}:
			for i, entry := range nested {
				if entryMap, ok := entry.(map[string]interface{}); ok {
					v.walkJSONataChoiceRule(name, loc+"/"+combinator+"/"+strconv.Itoa(i), entryMap, errorOutputAllowed)
				}
			}
		}
	}
}

// walkJSONataValue recurses through a JSONata-evaluated value tree,
// rejecting improperly delimited expression strings and checking node
// accessibility of the well-formed ones. Map keys are walked in sorted
// order so diagnostics are deterministic.
func (v *aslValidatorContext) walkJSONataValue(name, loc string, value interface{}, resultAllowed, errorOutputAllowed bool) {
	switch val := value.(type) {
	case string:
		if IsExpression(val) {
			v.checkJSONataAccessibility(name, loc, val, resultAllowed, errorOutputAllowed)
			return
		}
		if strings.Contains(val, "{%") || strings.Contains(val, "%}") {
			v.schemaf(loc,
				"State '%s': the string %q is not a valid JSONata expression: the string must start with '%s' with no leading spaces, and must end with '%s' with no trailing spaces",
				name, truncateForDiagnostic(val), "{%", "%}")
		}
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v.walkJSONataValue(name, loc+"/"+k, val[k], resultAllowed, errorOutputAllowed)
		}
	case []interface{}:
		for i, ev := range val {
			v.walkJSONataValue(name, loc+"/"+strconv.Itoa(i), ev, resultAllowed, errorOutputAllowed)
		}
	}
}

// checkJSONataAccessibility rejects references to $states.result or
// $states.errorOutput where the nodes are not accessible.
func (v *aslValidatorContext) checkJSONataAccessibility(name, loc, expression string, resultAllowed, errorOutputAllowed bool) {
	body := UnwrapExpression(expression)
	if !resultAllowed && jsonataResultRef.MatchString(body) {
		v.schemaf(loc,
			"State '%s': $states.result is only accessible in Task, Parallel, and Map states",
			name)
	}
	if !errorOutputAllowed && jsonataErrorOutputRef.MatchString(body) {
		v.schemaf(loc,
			"State '%s': $states.errorOutput is only accessible in the Catch field's Assign or Output",
			name)
	}
}

// truncateForDiagnostic bounds an offending string in a message. The cut
// is at a rune boundary — a byte-indexed cut splits a multi-byte
// character and prints mojibake.
func truncateForDiagnostic(s string) string {
	runes := []rune(s)
	if len(runes) > 40 {
		return string(runes[:40]) + "…"
	}
	return s
}

// JSONata variable-reference extraction: the state machine store keeps a
// per-state view of the variables each JSONata state references, built
// from Assign blocks across every scope and the {% ... %} expressions a
// state carries.
package sfn

import (
	"encoding/json"
	"sort"
	"strings"
	"unicode/utf8"
)

// scanVariableRefs returns every $name reference an expression carries,
// using the same UAX #31 grammar ValidateVariableName enforces — "The
// first character of a variable name must be a Unicode ID_Start
// character, and the second and subsequent characters must be Unicode
// ID_Continue characters". A character-class regex cannot express that
// grammar (the Other_ID_* additions minus the Pattern_* exclusions), so
// the scan walks runes through the shared predicates and a reference to
// a legal CJK-named or mark-carrying variable stays visible to the scan,
// not just to the validator.
func scanVariableRefs(expr string) []string {
	var names []string
	for i := 0; i < len(expr); {
		if expr[i] != '$' {
			i++
			continue
		}
		name, width := scanVariableName(expr[i+1:])
		if width > 0 {
			names = append(names, name)
			i += 1 + width
		} else {
			i++
		}
	}
	return names
}

// scanVariableName consumes one variable name from the front of s,
// reporting its width in bytes; an empty result means s does not begin a
// name.
func scanVariableName(s string) (string, int) {
	consumed := 0
	started := false
	for _, r := range s {
		if !started {
			if !isIDStart(r) {
				return "", 0
			}
			started = true
		} else if !isIDContinue(r) {
			break
		}
		consumed += utf8.RuneLen(r)
	}
	if !started {
		return "", 0
	}
	return s[:consumed], consumed
}

func extractVariableReferences(definition string) map[string][]string {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		return nil
	}

	// Variables are scoped per States object, but the variableReferences
	// map is keyed by state name alone, so every scope's assignments and
	// references are collected into one view.
	assignedVars := make(map[string]bool)
	collectAssigned := func(states map[string]interface{}) {
		for _, stateData := range states {
			stateMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}
			if assign, ok := stateMap["Assign"].(map[string]interface{}); ok {
				for name := range assign {
					clean := strings.TrimPrefix(name, "$")
					assignedVars[clean] = true
				}
			}
			if choices, ok := stateMap["Choices"].([]interface{}); ok {
				for _, choice := range choices {
					if choiceMap, ok := choice.(map[string]interface{}); ok {
						if assign, ok := choiceMap["Assign"].(map[string]interface{}); ok {
							for name := range assign {
								clean := strings.TrimPrefix(name, "$")
								assignedVars[clean] = true
							}
						}
					}
				}
			}
		}
	}

	result := make(map[string][]string)

	collectRefs := func(states map[string]interface{}) {
		for stateName, stateData := range states {
			stateMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}

			refs := collectVariableRefsFromState(stateMap, assignedVars)
			if len(refs) > 0 {
				sort.Strings(refs)
				result[stateName] = refs
			}
		}
	}

	// Nested scopes (Parallel branches, Map processors) contribute their
	// own states to the map.
	nestedStates := func(states map[string]interface{}) []map[string]interface{} {
		var nested []map[string]interface{}
		for _, stateData := range states {
			stateMap, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}
			if branches, ok := stateMap["Branches"].([]interface{}); ok {
				for _, branch := range branches {
					if branchMap, ok := branch.(map[string]interface{}); ok {
						if branchStates, ok := branchMap["States"].(map[string]interface{}); ok {
							nested = append(nested, branchStates)
						}
					}
				}
			}
			for _, field := range []string{"ItemProcessor", "Iterator"} {
				if sub, ok := stateMap[field].(map[string]interface{}); ok {
					if subStates, ok := sub["States"].(map[string]interface{}); ok {
						nested = append(nested, subStates)
					}
				}
			}
		}
		return nested
	}

	// Worklist over the scope tree: assignments from every scope are
	// collected before references so the merged view is complete.
	scopes := []map[string]interface{}{states}
	for i := 0; i < len(scopes); i++ {
		scopes = append(scopes, nestedStates(scopes[i])...)
	}
	for _, scope := range scopes {
		collectAssigned(scope)
	}
	for _, scope := range scopes {
		collectRefs(scope)
	}

	return result
}

func collectVariableRefsFromState(stateMap map[string]interface{}, assignedVars map[string]bool) []string {
	seen := make(map[string]bool)
	var refs []string

	for _, field := range []string{"Assign", "Output", "Arguments", "Items", "Condition"} {
		if field == "Condition" {
			if choices, ok := stateMap["Choices"].([]interface{}); ok {
				for _, choice := range choices {
					if choiceMap, ok := choice.(map[string]interface{}); ok {
						scanValueForVariableRefs(choiceMap["Condition"], seen, &refs, assignedVars)
					}
				}
			}
			continue
		}
		if val, ok := stateMap[field]; ok {
			scanValueForVariableRefs(val, seen, &refs, assignedVars)
		}
	}

	// Scalar expression members may reference variables too — the same
	// member list the definition validator treats as expression-bearing.
	// Non-string values (a plain numeric TimeoutSeconds) scan as no-ops:
	// only {% ... %} strings carry references.
	for _, field := range []string{
		"TimeoutSeconds", "HeartbeatSeconds", "Seconds", "Timestamp",
		"MaxConcurrency", "ToleratedFailureCount", "ToleratedFailurePercentage",
	} {
		if val, ok := stateMap[field]; ok {
			scanValueForVariableRefs(val, seen, &refs, assignedVars)
		}
	}

	return refs
}

var jsonataBuiltins = map[string]bool{
	"states": true, "context": true,
	"abs": true, "count": true, "sum": true, "max": true, "min": true, "average": true,
	"string": true, "substring": true, "length": true, "uppercase": true, "lowercase": true,
	"trim": true, "pad": true, "contains": true, "split": true, "join": true,
	"match": true, "replace": true, "base64encode": true, "base64decode": true,
	"number": true, "round": true, "floor": true, "ceil": true, "sqrt": true, "power": true,
	"random": true, "boolean": true, "not": true, "exists": true, "type": true,
	"each": true, "filter": true, "flatten": true, "keys": true, "lookup": true,
	"map": true, "merge": true, "reverse": true, "sort": true, "spread": true,
	"sift": true, "distinct": true, "single": true, "tail": true, "append": true,
	"errors": true, "fromMillis": true, "toMillis": true, "millis": true,
	"now": true, "uuid": true, "parse": true, "hash": true, "partition": true, "range": true,
	"description": true, "url": true, "encodeUrlComponent": true,
	"assert": true, "error": true, "order": true,
}

func scanValueForVariableRefs(v interface{}, seen map[string]bool, refs *[]string, assignedVars map[string]bool) {
	switch val := v.(type) {
	case string:
		if IsExpression(val) {
			expr := UnwrapExpression(val)
			for _, name := range scanVariableRefs(expr) {
				if !seen[name] && (!jsonataBuiltins[name] || assignedVars[name]) {
					seen[name] = true
					*refs = append(*refs, name)
				}
			}
		}
	case map[string]interface{}:
		for _, child := range val {
			scanValueForVariableRefs(child, seen, refs, assignedVars)
		}
	case []interface{}:
		for _, child := range val {
			scanValueForVariableRefs(child, seen, refs, assignedVars)
		}
	}
}

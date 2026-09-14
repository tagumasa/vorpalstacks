package sfn

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// scanPathShapeWarnings emits the documented path-shape warnings over a
// state's plain string fields: NO_PATH for values that look like a path
// under a field name that does not end with 'Path' and carries no ".$"
// suffix, and NO_DOLLAR for intrinsic-function-looking values on fields
// without the ".$" suffix.
func (v *aslValidatorContext) scanPathShapeWarnings(name string, stateMap map[string]interface{}, stateLoc string) {
	var walk func(value interface{}, fieldPath string)
	walk = func(value interface{}, fieldPath string) {
		switch node := value.(type) {
		case map[string]interface{}:
			keys := make([]string, 0, len(node))
			for k := range node {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				child := node[k]
				if strings.HasSuffix(k, ".$") || k == "Comment" || k == "Type" || k == "QueryLanguage" {
					continue
				}
				walk(child, fieldPath+"/"+k)
			}
		case []interface{}:
			for _, child := range node {
				walk(child, fieldPath)
			}
		case string:
			if strings.HasSuffix(fieldPath, "Path") || fieldPath == "" {
				return
			}
			if looksLikeJSONPath(node) {
				v.add("WARNING", "NO_PATH",
					fmt.Sprintf("The value of the field '%s' in state '%s' looks like a path, but the field name does not end with 'Path'", strings.TrimPrefix(fieldPath, "/"), name),
					stateLoc+fieldPath)
			} else if looksLikeIntrinsic(node) {
				v.add("WARNING", "NO_DOLLAR",
					fmt.Sprintf("No .$ on the field '%s' of state '%s' that appears to be a JSONPath or Intrinsic Function", strings.TrimPrefix(fieldPath, "/"), name),
					stateLoc+fieldPath)
			}
		}
	}
	stateKeys := make([]string, 0, len(stateMap))
	for k := range stateMap {
		stateKeys = append(stateKeys, k)
	}
	sort.Strings(stateKeys)
	for _, k := range stateKeys {
		child := stateMap[k]
		if strings.HasSuffix(k, ".$") || k == "Comment" || k == "Type" || k == "QueryLanguage" ||
			k == "Next" || k == "Default" || k == "Resource" || k == "Variable" || k == "StartAt" {
			continue
		}
		// A Pass Result that looks like a path carries its own dedicated
		// diagnostic (PASS_RESULT_IS_STATIC, emitted by the states walk);
		// the scan must not stack NO_PATH on top of it at the same
		// location.
		if k == "Result" && stateMap["Type"] == "Pass" {
			continue
		}
		walk(child, "/"+k)
	}
}

// checkTerminalReachability reports MISSING_END_STATE when no terminal
// state (End: true, Succeed or Fail) is reachable from the scope's start
// state.
func (v *aslValidatorContext) checkTerminalReachability(startAt string, states map[string]interface{}) {
	visited := map[string]bool{}
	queue := []string{startAt}
	reachesTerminal := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true

		stateMap, ok := states[current].(map[string]interface{})
		if !ok {
			continue
		}
		stateType, _ := stateMap["Type"].(string)
		if end, ok := stateMap["End"].(bool); ok && end {
			reachesTerminal = true
			break
		}
		if stateType == "Succeed" || stateType == "Fail" {
			reachesTerminal = true
			break
		}
		if next, ok := stateMap["Next"].(string); ok && next != "" {
			queue = append(queue, next)
		}
		// A Catcher's Next is an ordinary transition edge — on a matching
		// error the interpreter "transitions the machine to the state
		// named in the value of the 'Next' field" — so a terminal state
		// reachable only through catch targets still terminates the
		// workflow.
		if catchers, ok := stateMap["Catch"].([]interface{}); ok {
			for _, raw := range catchers {
				if catcher, ok := raw.(map[string]interface{}); ok {
					if next, ok := catcher["Next"].(string); ok && next != "" {
						queue = append(queue, next)
					}
				}
			}
		}
		if stateType == "Choice" {
			if rules, ok := stateMap["Choices"].([]interface{}); ok {
				for _, raw := range rules {
					if rule, ok := raw.(map[string]interface{}); ok {
						if next, ok := rule["Next"].(string); ok && next != "" {
							queue = append(queue, next)
						}
					}
				}
			}
			if def, ok := stateMap["Default"].(string); ok && def != "" {
				queue = append(queue, def)
			}
		}
	}

	if !reachesTerminal {
		v.add("ERROR", "MISSING_END_STATE", "The workflow does not have a terminal state reachable from StartAt", "")
	}
}

// looksLikeJSONPath reports whether a string value appears to be a
// JSONPath reference (the documented NO_DOLLAR / NO_PATH heuristic).
func looksLikeJSONPath(value string) bool {
	return strings.HasPrefix(value, "$.") || strings.HasPrefix(value, "$[")
}

// looksLikeIntrinsic reports whether a string value appears to be an
// intrinsic function invocation (the documented NO_DOLLAR heuristic).
func looksLikeIntrinsic(value string) bool {
	return strings.HasPrefix(value, "States.") && strings.HasSuffix(value, ")")
}

// tokenFrame tracks one decoded JSON object or array for duplicate-key
// detection.
type tokenFrame struct {
	isArray       bool
	isStatesValue bool
	expectKey     bool
	nextIsStates  bool
	keys          map[string]bool
}

// scanDuplicateStateKeys walks the raw JSON token stream and reports
// duplicate keys inside any States object, which map decoding collapses
// silently. Object members alternate key and value tokens while array
// elements are always values; nested containers push their own frames
// until they close.
func scanDuplicateStateKeys(definition string) []aslDiagnostic {
	dec := json.NewDecoder(strings.NewReader(definition))
	dec.UseNumber()

	var diagnostics []aslDiagnostic
	var stack []*tokenFrame

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case json.Delim:
			switch t {
			case '{':
				isStates := false
				if len(stack) > 0 {
					top := stack[len(stack)-1]
					isStates = top.nextIsStates
					top.nextIsStates = false
					top.expectKey = false
				}
				stack = append(stack, &tokenFrame{isStatesValue: isStates, expectKey: true, keys: map[string]bool{}})
			case '[':
				if len(stack) > 0 {
					top := stack[len(stack)-1]
					top.nextIsStates = false
					top.expectKey = false
				}
				stack = append(stack, &tokenFrame{isArray: true, keys: map[string]bool{}})
			case '}', ']':
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
				// A closed container completes its parent's member
				// value, so the parent expects a key again.
				if len(stack) > 0 && !stack[len(stack)-1].isArray {
					stack[len(stack)-1].expectKey = true
				}
			}
		case string:
			if len(stack) == 0 {
				continue
			}
			top := stack[len(stack)-1]
			if top.expectKey {
				if top.keys[t] && top.isStatesValue {
					diagnostics = append(diagnostics, aslDiagnostic{
						Severity: "ERROR",
						Code:     "DUPLICATE_STATE_NAME",
						Message:  fmt.Sprintf("The state name '%s' appears more than once", t),
						Location: "/States/" + t,
					})
				}
				top.keys[t] = true
				if t == "States" {
					top.nextIsStates = true
				}
				top.expectKey = false
			} else if !top.isArray {
				top.expectKey = true
				top.nextIsStates = false
			}
		default:
			// Scalar values (numbers, booleans, null) end the current
			// object member, so the next string is a key again.
			if len(stack) > 0 && !stack[len(stack)-1].isArray {
				stack[len(stack)-1].expectKey = true
				stack[len(stack)-1].nextIsStates = false
			}
		}
	}
	return diagnostics
}

// isForbiddenLabelRune reports whether a rune is forbidden in Map labels:
// whitespace, the wildcard, bracket and special characters and the
// control ranges from the Distributed Map documentation.
func isForbiddenLabelRune(r rune) bool {
	if unicode.IsSpace(r) || r <= 0x1f || (r >= 0x7f && r <= 0x9f) {
		return true
	}
	switch r {
	case '?', '*', '<', '>', '{', '}', '[', ']',
		':', ';', ',', '\\', '|', '^', '~', '$', '#', '%', '&', '`', '"':
		return true
	}
	return false
}

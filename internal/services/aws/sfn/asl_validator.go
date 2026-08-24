package sfn

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// aslDiagnostic is one ValidateStateMachineDefinition diagnostic: the
// documented severity/code/message triple plus the optional location
// member in the documented "/States/<StateName>/<FieldName>" form.
type aslDiagnostic struct {
	Severity string
	Code     string
	Message  string
	Location string
}

// aslKnownTopLevelFields is the verified top-level definition census from
// the Amazon States Language specification: Comment, StartAt, States,
// Version, TimeoutSeconds and the QueryLanguage extension.
var aslKnownTopLevelFields = map[string]bool{
	"Comment": true, "StartAt": true, "States": true,
	"Version": true, "TimeoutSeconds": true, "QueryLanguage": true,
}

// aslValidStateTypes is the state-type enumeration from the Amazon States
// Language specification.
var aslValidStateTypes = map[string]bool{
	"Pass": true, "Task": true, "Choice": true, "Wait": true,
	"Succeed": true, "Fail": true, "Parallel": true, "Map": true,
}

// aslStateFields is the closed-world census of legal state members per
// state type: the awslabs J2119 StateMachine schema unioned with the
// documented AWS extensions the engine parses (QueryLanguage, the JSONata
// Input and Output members, Assign, Arguments, the JSONata Choice
// Condition and the Distributed Map surface). The J2119 schema rejects
// any member a node does not declare, and the documented
// SCHEMA_VALIDATION_FAILED semantics carry that rejection onto create and
// update. Fail only allows Type and Comment of the common fields.
var aslStateFields = map[string]map[string]bool{
	"Pass": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
		"Result": true, "ResultPath": true, "ResultSelector": true,
		"Parameters": true, "Assign": true, "Next": true, "End": true,
	},
	"Task": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
		"Resource": true, "TimeoutSeconds": true, "TimeoutSecondsPath": true,
		"HeartbeatSeconds": true, "HeartbeatSecondsPath": true,
		"Parameters": true, "Arguments": true, "ResultPath": true,
		"ResultSelector": true, "Retry": true, "Catch": true,
		"Assign": true, "Credentials": true, "Next": true, "End": true,
	},
	"Choice": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
		"Choices": true, "Default": true, "Condition": true,
	},
	"Wait": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
		"Seconds": true, "SecondsPath": true, "Timestamp": true,
		"TimestampPath": true, "Assign": true, "Next": true, "End": true,
	},
	"Succeed": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
	},
	"Fail": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Error": true, "Cause": true, "ErrorPath": true, "CausePath": true,
	},
	"Parallel": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
		"Branches": true, "Parameters": true, "Arguments": true,
		"ResultPath": true, "ResultSelector": true, "Retry": true,
		"Catch": true, "Assign": true, "Next": true, "End": true,
	},
	"Map": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Input": true, "Output": true, "InputPath": true, "OutputPath": true,
		"Iterator": true, "ItemProcessor": true, "ItemsPath": true,
		"ItemSelector": true, "Parameters": true,
		"Items": true, "ItemBatcher": true, "ItemReader": true,
		"ResultWriter": true, "MaxConcurrency": true, "MaxConcurrencyPath": true,
		"Label": true, "ToleratedFailureCount": true,
		"ToleratedFailureCountPath": true, "ToleratedFailurePercentage": true,
		"ToleratedFailurePercentagePath": true, "ResultPath": true,
		"ResultSelector": true, "Retry": true, "Catch": true,
		"Assign": true, "Next": true, "End": true,
	},
}

// aslChoiceComparators is the JSONPath Choice-rule comparator census from
// the Amazon States Language specification (the Is* tests take no Path
// variant).
var aslChoiceComparators = map[string]bool{
	"StringEquals": true, "StringNotEquals": true,
	"StringEqualsPath": true, "StringNotEqualsPath": true,
	"StringLessThan": true, "StringGreaterThan": true,
	"StringLessThanEquals": true, "StringGreaterThanEquals": true,
	"StringLessThanPath": true, "StringGreaterThanPath": true,
	"StringLessThanEqualsPath": true, "StringGreaterThanEqualsPath": true,
	"NumericEquals": true, "NumericNotEquals": true,
	"NumericEqualsPath": true, "NumericNotEqualsPath": true,
	"NumericLessThan": true, "NumericGreaterThan": true,
	"NumericLessThanEquals": true, "NumericGreaterThanEquals": true,
	"NumericLessThanPath": true, "NumericGreaterThanPath": true,
	"NumericLessThanEqualsPath": true, "NumericGreaterThanEqualsPath": true,
	"BooleanEquals": true, "BooleanEqualsPath": true,
	"TimestampEquals": true, "TimestampLessThan": true,
	"TimestampGreaterThan": true, "TimestampLessThanEquals": true,
	"TimestampGreaterThanEquals": true,
	"TimestampEqualsPath":        true, "TimestampLessThanPath": true,
	"TimestampGreaterThanPath": true, "TimestampLessThanEqualsPath": true,
	"TimestampGreaterThanEqualsPath": true,
	"IsPresent":                      true, "IsNull": true, "IsString": true,
	"IsNumeric": true, "IsBoolean": true, "IsTimestamp": true,
	"StringMatches": true,
}

// aslValidatorContext carries the cross-scope validation state: state
// names are unique across the entire state machine and Map labels are
// unique within the definition.
type aslValidatorContext struct {
	smType      string
	stateNames  map[string]string
	labels      map[string]string
	diagnostics []aslDiagnostic
}

func (v *aslValidatorContext) add(severity, code, message, location string) {
	v.diagnostics = append(v.diagnostics, aslDiagnostic{
		Severity: severity, Code: code, Message: message, Location: location,
	})
}

func (v *aslValidatorContext) schemaf(location, format string, args ...interface{}) {
	v.add("ERROR", "SCHEMA_VALIDATION_FAILED", fmt.Sprintf(format, args...), location)
}

// validateASLStructure performs the structural validation of a state
// machine definition against the documented diagnostic code set: the
// severity levels and codes follow the ValidateStateMachineDefinition
// diagnostic reference, and only result stability (OK/FAIL) is a
// compatibility contract. smType gates the EXPRESS-only constraint on
// Distributed processing mode.
func validateASLStructure(definition, smType string) []aslDiagnostic {
	v := &aslValidatorContext{
		smType:     smType,
		stateNames: map[string]string{},
		labels:     map[string]string{},
	}

	if strings.TrimSpace(definition) == "" {
		v.add("ERROR", "MISSING_DESCRIPTION", "Received a null or empty workflow definition", "")
		return v.diagnostics
	}

	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &doc); err != nil {
		v.add("ERROR", "INVALID_JSON_DESCRIPTION", "JSON syntax problem found: "+err.Error(), "")
		return v.diagnostics
	}

	// Duplicate state-name keys collapse silently under map decoding, so
	// the raw token stream supplies the duplicate detection.
	for _, d := range scanDuplicateStateKeys(definition) {
		v.diagnostics = append(v.diagnostics, d)
	}

	for field := range doc {
		if !aslKnownTopLevelFields[field] {
			v.schemaf("/"+field, "Unknown top-level field %q in the state machine definition", field)
		}
	}
	if raw, ok := doc["QueryLanguage"]; ok {
		ql, _ := raw.(string)
		if ql != "JSONPath" && ql != "JSONata" {
			v.schemaf("/QueryLanguage", "QueryLanguage must be JSONPath or JSONata, got %q", ql)
		}
	}
	startAt, hasStartAt := doc["StartAt"].(string)
	if !hasStartAt || startAt == "" {
		v.schemaf("/StartAt", "State machine definition must include a string 'StartAt'")
	}
	states, hasStates := doc["States"].(map[string]interface{})
	if !hasStates {
		v.schemaf("/States", "State machine definition must include a 'States' object")
		return v.diagnostics
	}
	if hasStartAt {
		if _, exists := states[startAt]; !exists {
			v.schemaf("/StartAt", "StartAt '%s' does not reference a state in this scope", startAt)
		} else {
			v.checkTerminalReachability(startAt, states)
		}
	}

	v.validateStatesScope(doc, states, "")

	return v.diagnostics
}

// validateStatesScope validates the states of one States object (the
// top-level scope or a Parallel branch / Map processor scope). The
// effective query language governs which fields are legal per state.
func (v *aslValidatorContext) validateStatesScope(parent map[string]interface{}, states map[string]interface{}, scope string) {
	topQL, _ := parent["QueryLanguage"].(string)

	names := make([]string, 0, len(states))
	for name := range states {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if len([]rune(name)) > sfnstore.MaxStateNameLength {
			v.add("ERROR", "INVALID_STATE_NAME",
				fmt.Sprintf("The state name '%s' exceeds the allowed length of %d characters", name, sfnstore.MaxStateNameLength),
				"/States/"+name)
		}
		if prior, dup := v.stateNames[name]; dup {
			v.add("ERROR", "DUPLICATE_STATE_NAME",
				fmt.Sprintf("The state name '%s' appears more than once (also declared under %s)", name, prior),
				"/States/"+name)
		} else {
			v.stateNames[name] = strings.TrimPrefix(scope, "/")
		}

		stateMap, ok := states[name].(map[string]interface{})
		if !ok {
			v.schemaf("/States/"+name, "State '%s' must be a JSON object", name)
			continue
		}

		stateQL, _ := stateMap["QueryLanguage"].(string)
		effectiveQL := stateQL
		if effectiveQL == "" {
			effectiveQL = topQL
		}
		if effectiveQL == "" {
			effectiveQL = "JSONPath"
		}
		if effectiveQL != "JSONPath" && effectiveQL != "JSONata" {
			v.schemaf("/States/"+name+"/QueryLanguage", "QueryLanguage must be JSONPath or JSONata, got %q", effectiveQL)
			continue
		}

		stateType, _ := stateMap["Type"].(string)
		if stateType == "" {
			v.schemaf("/States/"+name, "State '%s' is missing the 'Type' field", name)
			continue
		}
		if !aslValidStateTypes[stateType] {
			v.schemaf("/States/"+name+"/Type", "State '%s' has an invalid Type '%s'", name, stateType)
			continue
		}

		v.validateOneState(name, stateType, stateMap, states, effectiveQL, scope)
	}
}

// validateOneState validates a single state object against its type
// contract.
func (v *aslValidatorContext) validateOneState(name, stateType string, stateMap, states map[string]interface{}, queryLanguage, scope string) {
	stateLoc := "/States/" + name
	jsonata := queryLanguage == "JSONata"

	// Query-language field separation: JSONata states reject the JSONPath
	// filters and vice versa.
	if jsonata {
		for _, f := range getJSONPathOnlyFields(stateType, stateMap) {
			v.schemaf(stateLoc+"/"+f, "State '%s' uses JSONata QueryLanguage but contains the JSONPath-only field '%s'", name, f)
		}
	} else {
		for _, f := range getJSONataOnlyFields(stateType, stateMap) {
			v.schemaf(stateLoc+"/"+f, "State '%s' uses JSONPath QueryLanguage but contains the JSONata-only field '%s'", name, f)
		}
	}

	// Transition shape: Task, Pass, Wait, Parallel and Map require exactly
	// one of Next or End; Choice, Succeed and Fail are terminal forms.
	next, hasNext := stateMap["Next"].(string)
	_, hasEnd := stateMap["End"]
	needsTransition := stateType == "Task" || stateType == "Pass" || stateType == "Wait" ||
		stateType == "Parallel" || stateType == "Map"
	if needsTransition {
		if hasNext == hasEnd {
			v.schemaf(stateLoc, "State '%s' must specify exactly one of 'Next' or 'End'", name)
		}
	} else {
		if hasNext {
			v.schemaf(stateLoc+"/Next", "A %s state must not specify 'Next'", stateType)
		}
		if hasEnd {
			v.schemaf(stateLoc+"/End", "A %s state must not specify 'End'", stateType)
		}
	}
	if hasNext {
		if _, exists := states[next]; !exists {
			v.add("ERROR", "MISSING_TRANSITION_TARGET",
				fmt.Sprintf("The value of the 'Next' field '%s' in state '%s' does not match a state name in this scope", next, name),
				stateLoc+"/Next")
		}
	}

	switch stateType {
	case "Task":
		resource, hasResource := stateMap["Resource"].(string)
		if !hasResource || resource == "" {
			v.schemaf(stateLoc, "Task state '%s' is missing the required 'Resource' field", name)
		} else {
			trimmed := resource
			for _, suffix := range []string{".sync", ".waitForTaskToken"} {
				if strings.HasSuffix(trimmed, suffix) {
					trimmed = strings.TrimSuffix(trimmed, suffix)
					break
				}
			}
			if _, err := arnutil.ParseARN(trimmed); err != nil {
				v.add("ERROR", "INVALID_RESOURCE",
					fmt.Sprintf("The value of the Task state '%s' resource field is invalid: %s", name, resource),
					stateLoc+"/Resource")
			}
		}
		v.validateTimeoutHeartbeat(name, stateMap, stateLoc)
	case "Pass":
		if result, ok := stateMap["Result"].(string); ok && looksLikeJSONPath(result) {
			v.add("WARNING", "PASS_RESULT_IS_STATIC",
				fmt.Sprintf("The Result of the Pass state '%s' looks like a path but is treated as a static value", name),
				stateLoc+"/Result")
		}
	case "Wait":
		v.validateWaitState(name, stateMap, stateLoc, jsonata)
	case "Choice":
		v.validateChoiceState(name, stateMap, states, stateLoc, jsonata)
	case "Parallel":
		v.validateParallelState(name, stateMap, stateLoc)
	case "Map":
		v.validateMapState(name, stateMap, stateLoc)
	case "Fail":
		v.validateFailState(name, stateMap, stateLoc)
	case "Succeed":
		// No members beyond the common contract.
	}

	// Unknown state members: the ASL schema is closed-world (the awslabs
	// J2119 StateMachine schema rejects any member a node does not
	// declare), so a member outside the census for the state's type is a
	// schema violation.
	if known := aslStateFields[stateType]; known != nil {
		for field := range stateMap {
			if !known[field] {
				v.schemaf(stateLoc+"/"+field, "State '%s' of type %s declares the unknown field '%s'", name, stateType, field)
			}
		}
	}

	v.validateRetryCatch(name, stateMap, states, stateLoc)
	v.scanPathShapeWarnings(name, stateMap, stateLoc)
}

// validateFailState ports the Fail field contract: only Type and Comment
// of the common fields, and each of Error/Cause excludes its reference
// path form (Fail state documentation).
func (v *aslValidatorContext) validateFailState(name string, stateMap map[string]interface{}, stateLoc string) {
	for _, field := range []string{"InputPath", "OutputPath", "Input", "Output"} {
		if _, ok := stateMap[field]; ok {
			v.schemaf(stateLoc+"/"+field, "Fail state '%s' must not specify '%s'", name, field)
		}
	}
	_, hasError := stateMap["Error"]
	_, hasErrorPath := stateMap["ErrorPath"]
	if hasError && hasErrorPath {
		v.schemaf(stateLoc+"/ErrorPath", "Fail state '%s' cannot specify both Error and ErrorPath", name)
	}
	_, hasCause := stateMap["Cause"]
	_, hasCausePath := stateMap["CausePath"]
	if hasCause && hasCausePath {
		v.schemaf(stateLoc+"/CausePath", "Fail state '%s' cannot specify both Cause and CausePath", name)
	}
}

// validateWaitState ports the Wait-state field contract: exactly one of
// Seconds, SecondsPath, Timestamp or TimestampPath (JSONata states only
// Seconds or Timestamp), integer Seconds within the documented ceiling,
// and the strict RFC3339 timestamp profile.
func (v *aslValidatorContext) validateWaitState(name string, stateMap map[string]interface{}, stateLoc string, jsonata bool) {
	if jsonata {
		if _, ok := stateMap["SecondsPath"]; ok {
			v.schemaf(stateLoc+"/SecondsPath", "Wait state SecondsPath is only supported in JSONPath states")
		}
		if _, ok := stateMap["TimestampPath"]; ok {
			v.schemaf(stateLoc+"/TimestampPath", "Wait state TimestampPath is only supported in JSONPath states")
		}
	}

	present := 0
	for _, field := range []string{"Seconds", "Timestamp"} {
		if _, ok := stateMap[field]; ok {
			present++
		}
	}
	if !jsonata {
		for _, field := range []string{"SecondsPath", "TimestampPath"} {
			if _, ok := stateMap[field]; ok {
				present++
			}
		}
	}
	if present != 1 {
		if jsonata {
			v.schemaf(stateLoc, "Wait state '%s' must specify exactly one of Seconds or Timestamp", name)
		} else {
			v.schemaf(stateLoc, "Wait state '%s' must specify exactly one of Seconds, Timestamp, SecondsPath, or TimestampPath", name)
		}
		return
	}

	if raw, ok := stateMap["Timestamp"]; ok {
		value, isString := raw.(string)
		if !isString {
			v.schemaf(stateLoc+"/Timestamp", "Wait state Timestamp must be a string")
			return
		}
		if jsonata && IsExpression(value) {
			return // the expression's result is only checkable at run time
		}
		if _, valid := parseWaitTimestamp(value); !valid {
			v.schemaf(stateLoc+"/Timestamp",
				"Wait state Timestamp %q must conform to the RFC3339 profile of ISO 8601 with an uppercase T and an uppercase Z or numeric offset, for example \"2024-03-14T01:59:00Z\"", value)
		}
		return
	}

	if raw, ok := stateMap["Seconds"]; ok {
		if value, isString := raw.(string); isString && jsonata && IsExpression(value) {
			return // the expression's result is only checkable at run time
		}
		value, isNumber := raw.(float64)
		if !isNumber || value != math.Trunc(value) || value < 0 || value > float64(sfnstore.MaxWaitSeconds) {
			v.schemaf(stateLoc+"/Seconds", "Wait state Seconds must be an integer value from 0 to %d", sfnstore.MaxWaitSeconds)
		}
		return
	}

	for _, field := range []string{"SecondsPath", "TimestampPath"} {
		if raw, ok := stateMap[field]; ok {
			value, isString := raw.(string)
			if !isString || value == "" {
				v.schemaf(stateLoc+"/"+field, "Wait state %s must be a non-empty path string", field)
			}
		}
	}
}

// validateChoiceState checks the Choice contract: a non-empty Choices
// array whose rules each carry exactly one comparator (or a nested
// And/Or/Not) plus an in-scope Next, with Default also in scope when
// present.
func (v *aslValidatorContext) validateChoiceState(name string, stateMap, states map[string]interface{}, stateLoc string, jsonata bool) {
	choices, ok := stateMap["Choices"].([]interface{})
	if !ok || len(choices) == 0 {
		v.schemaf(stateLoc+"/Choices", "Choice state '%s' must specify a non-empty 'Choices' array", name)
	} else {
		for i, raw := range choices {
			rule, ok := raw.(map[string]interface{})
			if !ok {
				v.schemaf(fmt.Sprintf("%s/Choices/%d", stateLoc, i), "Choice rules must be JSON objects")
				continue
			}
			v.validateChoiceRule(name, rule, states, stateLoc, jsonata)
		}
	}

	if raw, ok := stateMap["Default"]; ok {
		target, isString := raw.(string)
		if !isString {
			v.schemaf(stateLoc+"/Default", "Choice state Default must be a state name string")
		} else if _, exists := states[target]; !exists {
			v.add("ERROR", "MISSING_TRANSITION_TARGET",
				fmt.Sprintf("The value of the 'Default' field '%s' in state '%s' does not match a state name in this scope", target, name),
				stateLoc+"/Default")
		}
	}
}

// validateChoiceRule validates one Choice rule (recursively for the
// And/Or/Not combinators).
func (v *aslValidatorContext) validateChoiceRule(name string, rule, states map[string]interface{}, stateLoc string, jsonata bool) {
	combinators := 0
	for _, combinator := range []string{"And", "Or", "Not"} {
		if _, ok := rule[combinator]; ok {
			combinators++
		}
	}
	if combinators > 1 {
		v.schemaf(stateLoc+"/Choices", "Choice rule for state '%s' must not combine And, Or and Not", name)
	}

	comparators := 0
	for field := range rule {
		if aslChoiceComparators[field] || field == "Condition" {
			comparators++
		}
	}
	if combinators == 0 && comparators != 1 {
		v.schemaf(stateLoc+"/Choices", "Choice rule for state '%s' must specify exactly one comparator", name)
	}
	if !jsonata && comparators == 1 && combinators == 0 {
		if _, hasVariable := rule["Variable"]; !hasVariable {
			hasIs := false
			for field := range rule {
				if strings.HasPrefix(field, "Is") {
					hasIs = true
				}
			}
			if !hasIs {
				v.schemaf(stateLoc+"/Choices", "Choice rule for state '%s' must specify a 'Variable'", name)
			}
		}
	}

	if raw, ok := rule["Not"].(map[string]interface{}); ok {
		v.validateChoiceRule(name, raw, states, stateLoc, jsonata)
	}
	for _, combinator := range []string{"And", "Or"} {
		if raws, ok := rule[combinator].([]interface{}); ok {
			for _, raw := range raws {
				if nested, ok := raw.(map[string]interface{}); ok {
					v.validateChoiceRule(name, nested, states, stateLoc, jsonata)
				}
			}
		}
	}

	if raw, ok := rule["Next"]; ok {
		target, isString := raw.(string)
		if !isString {
			v.schemaf(stateLoc+"/Choices", "Choice rule 'Next' must be a state name string")
		} else if _, exists := states[target]; !exists {
			v.add("ERROR", "MISSING_TRANSITION_TARGET",
				fmt.Sprintf("The value of the 'Next' field '%s' in a rule of state '%s' does not match a state name in this scope", target, name),
				stateLoc+"/Choices")
		}
	} else {
		v.schemaf(stateLoc+"/Choices", "Choice rules for state '%s' must specify 'Next'", name)
	}
}

// validateParallelState validates the Branches array: each branch is a
// complete sub-machine with its own StartAt/States scope, its own
// terminal reachability, and no transitions crossing the branch boundary.
func (v *aslValidatorContext) validateParallelState(name string, stateMap map[string]interface{}, stateLoc string) {
	branches, ok := stateMap["Branches"].([]interface{})
	if !ok || len(branches) == 0 {
		v.schemaf(stateLoc+"/Branches", "Parallel state '%s' must specify a non-empty 'Branches' array", name)
		return
	}
	for i, raw := range branches {
		branch, ok := raw.(map[string]interface{})
		if !ok {
			v.schemaf(fmt.Sprintf("%s/Branches/%d", stateLoc, i), "Parallel branches must be state machine objects")
			continue
		}
		branchLoc := fmt.Sprintf("%s/Branches/%d", stateLoc, i)
		branchStartAt, hasStartAt := branch["StartAt"].(string)
		branchStates, hasStates := branch["States"].(map[string]interface{})
		if !hasStartAt || branchStartAt == "" {
			v.schemaf(branchLoc+"/StartAt", "Parallel branch of state '%s' must specify a string 'StartAt'", name)
		}
		if !hasStates {
			v.schemaf(branchLoc+"/States", "Parallel branch of state '%s' must specify a 'States' object", name)
			continue
		}
		if hasStartAt {
			if _, exists := branchStates[branchStartAt]; !exists {
				v.schemaf(branchLoc+"/StartAt", "StartAt '%s' of the branch does not reference a state in the branch", branchStartAt)
			} else {
				v.checkTerminalReachability(branchStartAt, branchStates)
			}
		}
		v.validateStatesScope(branch, branchStates, branchLoc)
	}
}

// validateMapState validates the Map contract: exactly one of Iterator or
// ItemProcessor, the ItemProcessor processing mode (Distributed mode is
// Standard-only), the ItemReader and ResultWriter structures, the
// tolerated-failure thresholds and the Label naming rules.
func (v *aslValidatorContext) validateMapState(name string, stateMap map[string]interface{}, stateLoc string) {
	_, hasIterator := stateMap["Iterator"]
	_, hasItemProcessor := stateMap["ItemProcessor"]
	if hasIterator == hasItemProcessor {
		v.schemaf(stateLoc, "Map state '%s' must specify exactly one of Iterator or ItemProcessor", name)
	}
	if raw, ok := stateMap["Iterator"].(map[string]interface{}); ok {
		v.validateMapSubMachine(name, "Iterator", raw, stateLoc)
	}
	if raw, ok := stateMap["ItemProcessor"].(map[string]interface{}); ok {
		v.validateMapSubMachine(name, "ItemProcessor", raw, stateLoc)
	}

	if raw, ok := stateMap["MaxConcurrency"]; ok {
		value, isNumber := raw.(float64)
		if !isNumber || value != math.Trunc(value) || value < 0 {
			v.schemaf(stateLoc+"/MaxConcurrency", "Map state MaxConcurrency must be a non-negative integer")
		}
		if _, also := stateMap["MaxConcurrencyPath"]; also {
			v.schemaf(stateLoc+"/MaxConcurrency", "Map state cannot include both MaxConcurrency and MaxConcurrencyPath")
		}
	}

	// Label naming: at most 40 characters, unique within the definition,
	// and free of the documented forbidden characters.
	if raw, ok := stateMap["Label"]; ok {
		label, isString := raw.(string)
		if !isString {
			v.schemaf(stateLoc+"/Label", "Map state Label must be a string")
		} else {
			if len([]rune(label)) > sfnstore.MaxMapLabelLength {
				v.add("ERROR", "INVALID_LABEL_NAME",
					fmt.Sprintf("The label '%s' of state '%s' exceeds the allowed length of %d characters", label, name, sfnstore.MaxMapLabelLength),
					stateLoc+"/Label")
			}
			if strings.ContainsFunc(label, isForbiddenLabelRune) {
				v.add("ERROR", "INVALID_LABEL_NAME",
					fmt.Sprintf("The label '%s' of state '%s' contains forbidden characters", label, name),
					stateLoc+"/Label")
			}
			if prior, dup := v.labels[label]; dup {
				v.add("ERROR", "DUPLICATE_LABEL_NAME",
					fmt.Sprintf("The label name '%s' appears more than once (also used by state %s)", label, prior),
					stateLoc+"/Label")
			} else {
				v.labels[label] = name
			}
		}
	}

	v.validateItemReader(name, stateMap, stateLoc)
	v.validateItemBatcher(name, stateMap, stateLoc)
	v.validateResultWriter(name, stateMap, stateLoc)
	v.validateToleratedFailure(name, stateMap, stateLoc)
}

// validateMapSubMachine validates the Iterator or ItemProcessor
// sub-machine and its Distributed-mode constraints.
func (v *aslValidatorContext) validateMapSubMachine(name, field string, sub map[string]interface{}, stateLoc string) {
	subLoc := stateLoc + "/" + field
	startAt, hasStartAt := sub["StartAt"].(string)
	subStates, hasStates := sub["States"].(map[string]interface{})
	if !hasStartAt || startAt == "" {
		v.schemaf(subLoc+"/StartAt", "The %s of Map state '%s' must specify a string 'StartAt'", field, name)
	}
	if !hasStates {
		v.schemaf(subLoc+"/States", "The %s of Map state '%s' must specify a 'States' object", field, name)
		return
	}
	if hasStartAt {
		if _, exists := subStates[startAt]; !exists {
			v.schemaf(subLoc+"/StartAt", "StartAt '%s' of the %s does not reference a state in it", startAt, field)
		} else {
			v.checkTerminalReachability(startAt, subStates)
		}
	}

	if pc, ok := sub["ProcessorConfig"].(map[string]interface{}); ok {
		mode, _ := pc["Mode"].(string)
		if mode != "" && mode != "INLINE" && mode != "DISTRIBUTED" {
			v.schemaf(subLoc+"/ProcessorConfig/Mode", "ProcessorConfig Mode must be INLINE or DISTRIBUTED, got %q", mode)
		}
		if mode == "DISTRIBUTED" {
			executionType, _ := pc["ExecutionType"].(string)
			if executionType != "STANDARD" && executionType != "EXPRESS" {
				v.schemaf(subLoc+"/ProcessorConfig/ExecutionType", "DISTRIBUTED processing requires ExecutionType STANDARD or EXPRESS")
			}
			if v.smType == "EXPRESS" {
				v.schemaf(subLoc+"/ProcessorConfig/Mode", "Distributed mode is supported in Standard workflows but not in Express workflows")
			}
		}
	}
	v.validateStatesScope(sub, subStates, subLoc)
}

// validateItemReader validates the ItemReader structure: resource, reader
// configuration enums, the MaxItems exclusivity and range, CSV header
// requirements and the manifest constraints.
func (v *aslValidatorContext) validateItemReader(name string, stateMap map[string]interface{}, stateLoc string) {
	raw, ok := stateMap["ItemReader"].(map[string]interface{})
	if !ok {
		return
	}
	readerLoc := stateLoc + "/ItemReader"

	resource, _ := raw["Resource"].(string)
	if resource == "" {
		v.schemaf(readerLoc, "The ItemReader of Map state '%s' must specify a Resource", name)
	} else {
		known := strings.HasSuffix(resource, ":s3:getObject") ||
			strings.HasSuffix(resource, ":aws-sdk:s3:getObject") ||
			strings.HasSuffix(resource, ":s3:listObjectsV2")
		if !known {
			v.schemaf(readerLoc+"/Resource", "The ItemReader Resource '%s' of Map state '%s' is not a supported reader", resource, name)
		}
	}
	_, hasParams := raw["Parameters"]
	_, hasArgs := raw["Arguments"]
	if hasParams && hasArgs {
		v.schemaf(readerLoc, "The ItemReader of Map state '%s' cannot specify both Parameters and Arguments", name)
	}

	rc, hasRC := raw["ReaderConfig"].(map[string]interface{})
	if !hasRC {
		return
	}
	rcLoc := readerLoc + "/ReaderConfig"

	if inputType, ok := rc["InputType"].(string); ok {
		switch inputType {
		case "CSV", "JSON", "JSONL", "PARQUET", "MANIFEST":
		default:
			v.schemaf(rcLoc+"/InputType", "ReaderConfig InputType '%s' is not one of CSV, JSON, JSONL, PARQUET or MANIFEST", inputType)
		}
	}
	if csvHeaderLocation, ok := rc["CSVHeaderLocation"].(string); ok {
		if csvHeaderLocation != "FIRST_ROW" && csvHeaderLocation != "GIVEN" {
			v.schemaf(rcLoc+"/CSVHeaderLocation", "ReaderConfig CSVHeaderLocation must be FIRST_ROW or GIVEN")
		}
		if csvHeaderLocation == "GIVEN" {
			headers, hasHeaders := rc["CSVHeaders"].([]interface{})
			if !hasHeaders || len(headers) == 0 {
				v.schemaf(rcLoc+"/CSVHeaders", "ReaderConfig CSVHeaders is required when CSVHeaderLocation is GIVEN")
			} else {
				total := 0
				for _, h := range headers {
					if s, ok := h.(string); ok {
						total += len(s)
					}
				}
				if total > sfnstore.MaxCSVHeaderBytes {
					v.schemaf(rcLoc+"/CSVHeaders", "ReaderConfig CSVHeaders exceed the %d byte header ceiling", sfnstore.MaxCSVHeaderBytes)
				}
			}
		}
	}
	if delimiter, ok := rc["CSVDelimiter"].(string); ok {
		switch delimiter {
		case "COMMA", "PIPE", "SEMICOLON", "SPACE", "TAB":
		default:
			v.schemaf(rcLoc+"/CSVDelimiter", "ReaderConfig CSVDelimiter '%s' is not one of COMMA, PIPE, SEMICOLON, SPACE or TAB", delimiter)
		}
	}
	_, hasMaxItems := rc["MaxItems"]
	maxItemsPath, hasMaxItemsPath := rc["MaxItemsPath"].(string)
	if hasMaxItems && hasMaxItemsPath {
		v.schemaf(rcLoc+"/MaxItems", "ReaderConfig cannot specify both MaxItems and MaxItemsPath")
	}
	if hasMaxItems {
		if value, ok := rc["MaxItems"].(float64); ok {
			if value != math.Trunc(value) || value < 0 || value > float64(sfnstore.MaxItemReaderItems) {
				v.schemaf(rcLoc+"/MaxItems", "ReaderConfig MaxItems must be an integer from 0 to %d", sfnstore.MaxItemReaderItems)
			}
		} else {
			v.schemaf(rcLoc+"/MaxItems", "ReaderConfig MaxItems must be an integer")
		}
	}
	if hasMaxItemsPath && maxItemsPath == "" {
		v.schemaf(rcLoc+"/MaxItemsPath", "ReaderConfig MaxItemsPath must be a non-empty path string")
	}
	if pointer, ok := rc["ItemsPointer"].(string); ok && pointer != "" {
		if inputType, _ := rc["InputType"].(string); inputType != "JSON" {
			v.schemaf(rcLoc+"/ItemsPointer", "ReaderConfig ItemsPointer can only be specified when InputType is JSON")
		}
	}
	if transformation, ok := rc["Transformation"].(string); ok {
		if transformation != "NONE" && transformation != "LOAD_AND_FLATTEN" {
			v.schemaf(rcLoc+"/Transformation", "ReaderConfig Transformation must be NONE or LOAD_AND_FLATTEN")
		}
		if transformation == "LOAD_AND_FLATTEN" {
			if inputType, _ := rc["InputType"].(string); inputType == "" {
				v.schemaf(rcLoc+"/InputType", "InputType is required when Transformation is LOAD_AND_FLATTEN")
			}
		}
	}
	if manifestType, ok := rc["ManifestType"].(string); ok {
		if manifestType != "ATHENA_DATA" && manifestType != "S3_INVENTORY" {
			v.schemaf(rcLoc+"/ManifestType", "ReaderConfig ManifestType must be ATHENA_DATA or S3_INVENTORY")
		}
		if manifestType == "S3_INVENTORY" {
			if _, hasInputType := rc["InputType"]; hasInputType {
				v.schemaf(rcLoc+"/InputType", "InputType cannot be specified when ManifestType is S3_INVENTORY")
			}
		}
		if manifestType == "ATHENA_DATA" {
			if inputType, _ := rc["InputType"].(string); inputType == "" {
				v.schemaf(rcLoc+"/InputType", "InputType is required when ManifestType is ATHENA_DATA")
			}
		}
	}
}

// validateItemBatcher validates the ItemBatcher structure: the literal and
// reference-path forms are mutually exclusive within each pair, at least
// one sizing value is required, item counts are positive integers, the
// batch byte cap stays within the 256 KiB child-execution input bound and
// the fixed BatchInput is an object (ItemBatcher documentation).
func (v *aslValidatorContext) validateItemBatcher(name string, stateMap map[string]interface{}, stateLoc string) {
	raw, ok := stateMap["ItemBatcher"].(map[string]interface{})
	if !ok {
		return
	}
	batcherLoc := stateLoc + "/ItemBatcher"

	_, hasMaxItems := raw["MaxItemsPerBatch"]
	maxItemsPath, hasMaxItemsPath := raw["MaxItemsPerBatchPath"].(string)
	if hasMaxItems && hasMaxItemsPath {
		v.schemaf(batcherLoc+"/MaxItemsPerBatch", "ItemBatcher cannot specify both MaxItemsPerBatch and MaxItemsPerBatchPath")
	}
	if hasMaxItems {
		if value, ok := raw["MaxItemsPerBatch"].(float64); ok {
			if value != math.Trunc(value) || value < 1 {
				v.schemaf(batcherLoc+"/MaxItemsPerBatch", "ItemBatcher MaxItemsPerBatch must be a positive integer")
			}
		} else if _, isString := raw["MaxItemsPerBatch"].(string); !isString {
			// A string value is the JSONata expression form; the numeric
			// bound applies to literal counts.
			v.schemaf(batcherLoc+"/MaxItemsPerBatch", "ItemBatcher MaxItemsPerBatch must be a positive integer")
		}
	}
	if hasMaxItemsPath && maxItemsPath == "" {
		v.schemaf(batcherLoc+"/MaxItemsPerBatchPath", "ItemBatcher MaxItemsPerBatchPath must be a non-empty reference path")
	}

	_, hasMaxBytes := raw["MaxInputBytesPerBatch"]
	maxBytesPath, hasMaxBytesPath := raw["MaxInputBytesPerBatchPath"].(string)
	if hasMaxBytes && hasMaxBytesPath {
		v.schemaf(batcherLoc+"/MaxInputBytesPerBatch", "ItemBatcher cannot specify both MaxInputBytesPerBatch and MaxInputBytesPerBatchPath")
	}
	if hasMaxBytes {
		if value, ok := raw["MaxInputBytesPerBatch"].(float64); ok {
			if value != math.Trunc(value) || value < 1 || value > float64(sfnstore.MaxExecutionDataBytes) {
				v.schemaf(batcherLoc+"/MaxInputBytesPerBatch", "ItemBatcher MaxInputBytesPerBatch must be an integer from 1 to %d", sfnstore.MaxExecutionDataBytes)
			}
		} else {
			v.schemaf(batcherLoc+"/MaxInputBytesPerBatch", "ItemBatcher MaxInputBytesPerBatch must be a positive integer")
		}
	}
	if hasMaxBytesPath && maxBytesPath == "" {
		v.schemaf(batcherLoc+"/MaxInputBytesPerBatchPath", "ItemBatcher MaxInputBytesPerBatchPath must be a non-empty reference path")
	}

	if !hasMaxItems && !hasMaxItemsPath && !hasMaxBytes && !hasMaxBytesPath {
		v.schemaf(batcherLoc, "ItemBatcher must specify MaxItemsPerBatch, MaxInputBytesPerBatch or both to batch items")
	}

	_, hasBatchInput := raw["BatchInput"]
	_, hasBatchInputPath := raw["BatchInputPath"]
	if hasBatchInput && hasBatchInputPath {
		v.schemaf(batcherLoc+"/BatchInput", "ItemBatcher cannot specify both BatchInput and BatchInputPath")
	}
	if hasBatchInput {
		if _, isObject := raw["BatchInput"].(map[string]interface{}); !isObject {
			v.schemaf(batcherLoc+"/BatchInput", "ItemBatcher BatchInput must be a JSON object")
		}
	}
}

// validateResultWriter validates the documented required field
// combinations: WriterConfig alone, Resource with Parameters, or all
// three.
func (v *aslValidatorContext) validateResultWriter(name string, stateMap map[string]interface{}, stateLoc string) {
	raw, ok := stateMap["ResultWriter"].(map[string]interface{})
	if !ok {
		return
	}
	writerLoc := stateLoc + "/ResultWriter"

	_, hasResource := raw["Resource"]
	_, hasParams := raw["Parameters"]
	_, hasArgs := raw["Arguments"]
	_, hasWriterConfig := raw["WriterConfig"]

	if hasResource && !hasParams && !hasArgs && !hasWriterConfig {
		v.schemaf(writerLoc, "The ResultWriter of Map state '%s' must specify Parameters with Resource", name)
	}
	if (hasParams || hasArgs) && !hasResource && !hasWriterConfig {
		v.schemaf(writerLoc, "The ResultWriter of Map state '%s' must specify a Resource with Parameters", name)
	}
	if hasParams && hasArgs {
		v.schemaf(writerLoc, "The ResultWriter of Map state '%s' cannot specify both Parameters and Arguments", name)
	}
	if resource, ok := raw["Resource"].(string); ok && resource != "" &&
		!strings.HasSuffix(resource, ":s3:putObject") {
		v.schemaf(writerLoc+"/Resource", "The ResultWriter Resource '%s' is not the s3:putObject writer", resource)
	}
	if wc, ok := raw["WriterConfig"].(map[string]interface{}); ok {
		if transformation, ok := wc["Transformation"].(string); ok {
			if transformation != "NONE" && transformation != "COMPACT" && transformation != "FLATTEN" {
				v.schemaf(writerLoc+"/WriterConfig/Transformation", "WriterConfig Transformation must be NONE, COMPACT or FLATTEN")
			}
		}
		if outputType, ok := wc["OutputType"].(string); ok {
			if outputType != "JSON" && outputType != "JSONL" {
				v.schemaf(writerLoc+"/WriterConfig/OutputType", "WriterConfig OutputType must be JSON or JSONL")
			}
		}
	}
}

// validateToleratedFailure validates the threshold fields: non-negative
// counts, percentages within zero to one hundred, and the value/path
// exclusivity.
func (v *aslValidatorContext) validateToleratedFailure(name string, stateMap map[string]interface{}, stateLoc string) {
	if raw, ok := stateMap["ToleratedFailureCount"]; ok {
		value, isNumber := raw.(float64)
		if !isNumber || value != math.Trunc(value) || value < 0 {
			v.schemaf(stateLoc+"/ToleratedFailureCount", "Map state ToleratedFailureCount must be a non-negative integer")
		}
		if _, also := stateMap["ToleratedFailureCountPath"]; also {
			v.schemaf(stateLoc+"/ToleratedFailureCount", "Map state cannot specify both ToleratedFailureCount and ToleratedFailureCountPath")
		}
	}
	if raw, ok := stateMap["ToleratedFailurePercentage"]; ok {
		value, isNumber := raw.(float64)
		if !isNumber || value < 0 || value > 100 {
			v.schemaf(stateLoc+"/ToleratedFailurePercentage", "Map state ToleratedFailurePercentage must be between zero and 100")
		}
		if _, also := stateMap["ToleratedFailurePercentagePath"]; also {
			v.schemaf(stateLoc+"/ToleratedFailurePercentage", "Map state cannot specify both ToleratedFailurePercentage and ToleratedFailurePercentagePath")
		}
	}
}

// validateTimeoutHeartbeat checks the Task timeout contract: positive
// integers with HeartbeatSeconds smaller than TimeoutSeconds.
func (v *aslValidatorContext) validateTimeoutHeartbeat(name string, stateMap map[string]interface{}, stateLoc string) {
	timeout, hasTimeout := stateMap["TimeoutSeconds"].(float64)
	heartbeat, hasHeartbeat := stateMap["HeartbeatSeconds"].(float64)
	if raw, ok := stateMap["TimeoutSeconds"]; ok {
		if value, isString := raw.(string); isString && !IsExpression(value) {
			v.schemaf(stateLoc+"/TimeoutSeconds", "Task state TimeoutSeconds must be a positive integer")
		} else if !isString && (!hasTimeout || timeout != math.Trunc(timeout) || timeout <= 0) {
			v.schemaf(stateLoc+"/TimeoutSeconds", "Task state TimeoutSeconds must be a positive integer")
		}
	}
	if raw, ok := stateMap["HeartbeatSeconds"]; ok {
		if value, isString := raw.(string); isString && !IsExpression(value) {
			v.schemaf(stateLoc+"/HeartbeatSeconds", "Task state HeartbeatSeconds must be a positive integer")
		} else if !isString && (!hasHeartbeat || heartbeat != math.Trunc(heartbeat) || heartbeat <= 0) {
			v.schemaf(stateLoc+"/HeartbeatSeconds", "Task state HeartbeatSeconds must be a positive integer")
		}
	}
	if hasTimeout && hasHeartbeat && heartbeat >= timeout {
		v.schemaf(stateLoc+"/HeartbeatSeconds", "Task state HeartbeatSeconds must be smaller than TimeoutSeconds")
	}
	// Each reference-path form excludes its literal form, and the paths
	// are JSONPath-only members.
	_, hasTimeoutPath := stateMap["TimeoutSecondsPath"]
	if _, hasLiteral := stateMap["TimeoutSeconds"]; hasTimeoutPath && hasLiteral {
		v.schemaf(stateLoc+"/TimeoutSecondsPath", "Task state cannot specify both TimeoutSeconds and TimeoutSecondsPath")
	}
	_, hasHeartbeatPath := stateMap["HeartbeatSecondsPath"]
	if _, hasLiteral := stateMap["HeartbeatSeconds"]; hasHeartbeatPath && hasLiteral {
		v.schemaf(stateLoc+"/HeartbeatSecondsPath", "Task state cannot specify both HeartbeatSeconds and HeartbeatSecondsPath")
	}
}

// validateRetryCatch validates the Retrier and Catcher contracts:
// non-empty ErrorEquals, documented ranges for the retry timing fields,
// and in-scope Next targets for catchers.
func (v *aslValidatorContext) validateRetryCatch(name string, stateMap, states map[string]interface{}, stateLoc string) {
	if raws, ok := stateMap["Retry"].([]interface{}); ok {
		for i, raw := range raws {
			retrier, ok := raw.(map[string]interface{})
			if !ok {
				v.schemaf(fmt.Sprintf("%s/Retry/%d", stateLoc, i), "Retriers must be JSON objects")
				continue
			}
			retrierLoc := fmt.Sprintf("%s/Retry/%d", stateLoc, i)
			if equals, ok := retrier["ErrorEquals"].([]interface{}); !ok || len(equals) == 0 {
				v.schemaf(retrierLoc+"/ErrorEquals", "Retriers of state '%s' must specify a non-empty ErrorEquals array", name)
			}
			if raw, ok := retrier["IntervalSeconds"]; ok {
				if value, isNumber := raw.(float64); !isNumber || value != math.Trunc(value) || value <= 0 || value > float64(sfnstore.MaxRetryIntervalSeconds) {
					v.schemaf(retrierLoc+"/IntervalSeconds", "Retry IntervalSeconds must be an integer from 1 to %d", sfnstore.MaxRetryIntervalSeconds)
				}
			}
			if raw, ok := retrier["MaxAttempts"]; ok {
				if value, isNumber := raw.(float64); !isNumber || value != math.Trunc(value) || value < 0 || value > float64(sfnstore.MaxRetryAttempts) {
					v.schemaf(retrierLoc+"/MaxAttempts", "Retry MaxAttempts must be an integer from 0 to %d", sfnstore.MaxRetryAttempts)
				}
			}
			if raw, ok := retrier["BackoffRate"]; ok {
				if value, isNumber := raw.(float64); !isNumber || value < 1.0 {
					v.schemaf(retrierLoc+"/BackoffRate", "Retry BackoffRate must be greater than or equal to 1.0")
				}
			}
			if raw, ok := retrier["MaxDelaySeconds"]; ok {
				if value, isNumber := raw.(float64); !isNumber || value != math.Trunc(value) || value <= 0 || value >= float64(sfnstore.MaxRetryDelaySeconds) {
					v.schemaf(retrierLoc+"/MaxDelaySeconds", "Retry MaxDelaySeconds must be greater than 0 and less than %d", sfnstore.MaxRetryDelaySeconds)
				}
			}
			if strategy, ok := retrier["JitterStrategy"].(string); ok {
				if strategy != "FULL" && strategy != "NONE" {
					v.schemaf(retrierLoc+"/JitterStrategy", "Retry JitterStrategy must be FULL or NONE")
				}
			}
		}
	}

	if raws, ok := stateMap["Catch"].([]interface{}); ok {
		for i, raw := range raws {
			catcher, ok := raw.(map[string]interface{})
			if !ok {
				v.schemaf(fmt.Sprintf("%s/Catch/%d", stateLoc, i), "Catchers must be JSON objects")
				continue
			}
			catcherLoc := fmt.Sprintf("%s/Catch/%d", stateLoc, i)
			if equals, ok := catcher["ErrorEquals"].([]interface{}); !ok || len(equals) == 0 {
				v.schemaf(catcherLoc+"/ErrorEquals", "Catchers of state '%s' must specify a non-empty ErrorEquals array", name)
			}
			if target, ok := catcher["Next"].(string); !ok || target == "" {
				v.schemaf(catcherLoc+"/Next", "Catchers of state '%s' must specify a Next state", name)
			} else if _, exists := states[target]; !exists {
				v.add("ERROR", "MISSING_TRANSITION_TARGET",
					fmt.Sprintf("The value of the 'Next' field '%s' in a Catcher of state '%s' does not match a state name in this scope", target, name),
					catcherLoc+"/Next")
			}
		}
	}
}

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
			for k, child := range node {
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
	for k, child := range stateMap {
		if strings.HasSuffix(k, ".$") || k == "Comment" || k == "Type" || k == "QueryLanguage" ||
			k == "Next" || k == "Default" || k == "Resource" || k == "Variable" || k == "StartAt" {
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

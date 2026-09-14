package sfn

import (
	sfnstore "vorpalstacks/internal/store/aws/sfn"

	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// aslStateFields is the closed-world census of legal state members per
// state type: the awslabs J2119 StateMachine schema unioned with the
// documented AWS extensions the engine parses (QueryLanguage, the JSONata
// Output member, Assign, Arguments, the JSONata Choice Condition and the
// Distributed Map surface). The J2119 schema rejects any member a node
// does not declare, and the documented SCHEMA_VALIDATION_FAILED semantics
// carry that rejection onto create and update. Fail only allows Type and
// Comment of the common fields.
var aslStateFields = map[string]map[string]bool{
	"Pass": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
		"Result": true, "ResultPath": true, "ResultSelector": true,
		"Parameters": true, "Assign": true, "Next": true, "End": true,
	},
	"Task": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
		"Resource": true, "TimeoutSeconds": true, "TimeoutSecondsPath": true,
		"HeartbeatSeconds": true, "HeartbeatSecondsPath": true,
		"Parameters": true, "Arguments": true, "ResultPath": true,
		"ResultSelector": true, "Retry": true, "Catch": true,
		"Assign": true, "Credentials": true, "Next": true, "End": true,
	},
	"Choice": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
		"Choices": true, "Default": true, "Assign": true,
	},
	"Wait": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
		"Seconds": true, "SecondsPath": true, "Timestamp": true,
		"TimestampPath": true, "Assign": true, "Next": true, "End": true,
	},
	"Succeed": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
	},
	"Fail": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Error": true, "Cause": true, "ErrorPath": true, "CausePath": true,
	},
	"Parallel": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
		"Branches": true, "Parameters": true, "Arguments": true,
		"ResultPath": true, "ResultSelector": true, "Retry": true,
		"Catch": true, "Assign": true, "Next": true, "End": true,
	},
	"Map": {
		"Type": true, "Comment": true, "QueryLanguage": true,
		"Output": true, "InputPath": true, "OutputPath": true,
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

// getJSONPathOnlyFields and getJSONataOnlyFields enforce the
// query-language field separation: the members each dialect owns (the
// JSONPath *Path filters, the JSONata Arguments/Items/Condition surface)
// are rejected on states written in the other dialect.

func getJSONPathOnlyFields(stateType string, stateMap map[string]interface{}) []string {
	var forbidden []string

	switch stateType {
	case "Pass", "Task", "Parallel", "Map", "Succeed":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
		if _, ok := stateMap["Parameters"]; ok {
			forbidden = append(forbidden, "Parameters")
		}
	case "Choice":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
	case "Wait":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
		if _, ok := stateMap["SecondsPath"]; ok {
			forbidden = append(forbidden, "SecondsPath")
		}
		if _, ok := stateMap["TimestampPath"]; ok {
			forbidden = append(forbidden, "TimestampPath")
		}
	}

	switch stateType {
	case "Task":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
		if _, ok := stateMap["TimeoutSecondsPath"]; ok {
			forbidden = append(forbidden, "TimeoutSecondsPath")
		}
		if _, ok := stateMap["HeartbeatSecondsPath"]; ok {
			forbidden = append(forbidden, "HeartbeatSecondsPath")
		}
	case "Pass":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Map":
		if _, ok := stateMap["ItemsPath"]; ok {
			forbidden = append(forbidden, "ItemsPath")
		}
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
		if _, ok := stateMap["MaxConcurrencyPath"]; ok {
			forbidden = append(forbidden, "MaxConcurrencyPath")
		}
		if _, ok := stateMap["ToleratedFailureCountPath"]; ok {
			forbidden = append(forbidden, "ToleratedFailureCountPath")
		}
		if _, ok := stateMap["ToleratedFailurePercentagePath"]; ok {
			forbidden = append(forbidden, "ToleratedFailurePercentagePath")
		}
		if _, ok := stateMap["ItemSelector"]; ok {
			forbidden = append(forbidden, "ItemSelector")
		}
	case "Parallel":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Fail":
		if _, ok := stateMap["CausePath"]; ok {
			forbidden = append(forbidden, "CausePath")
		}
		if _, ok := stateMap["ErrorPath"]; ok {
			forbidden = append(forbidden, "ErrorPath")
		}
	}

	return forbidden
}

func getJSONataOnlyFields(stateType string, stateMap map[string]interface{}) []string {
	var forbidden []string

	// "The Arguments and Output fields only support JSONata, so it is
	// invalid to use them with workflows that use JSONPath." — Output is a
	// JSONata-only member of every state type that carries it; a state
	// whose census does not know Output (Fail) gets its single diagnosis
	// from the unknown-field walk. A state-level Input member is not part
	// of the documented common field set at all and the census rejects it
	// on both dialects.
	for _, field := range []string{"Output"} {
		if _, ok := stateMap[field]; ok {
			if aslStateFields[stateType][field] {
				forbidden = append(forbidden, field)
			}
		}
	}

	if _, ok := stateMap["Arguments"]; ok {
		if stateType == "Task" || stateType == "Parallel" {
			forbidden = append(forbidden, "Arguments")
		}
	}

	if _, ok := stateMap["Items"]; ok {
		if stateType == "Map" {
			forbidden = append(forbidden, "Items")
		}
	}

	// Condition is documented only inside Choice rules; a state-level
	// Condition is unknown on both dialects and the census rejects it.

	if stateType == "Choice" {
		if choices, ok := stateMap["Choices"].([]interface{}); ok {
			for _, choice := range choices {
				if choiceMap, ok := choice.(map[string]interface{}); ok {
					if _, ok := choiceMap["Condition"]; ok {
						forbidden = append(forbidden, "Condition")
						break
					}
				}
			}
		}
	}

	return forbidden
}

// aslChoiceComparators is the JSONPath Choice-rule comparator census from
// the documented supported-operator list (the states-language
// specification and the Choice page agree on it — there are no NotEquals
// variants); the Is* tests take no Path variant.
var aslChoiceComparators = map[string]bool{
	"StringEquals":     true,
	"StringEqualsPath": true,
	"StringLessThan":   true, "StringGreaterThan": true,
	"StringLessThanEquals": true, "StringGreaterThanEquals": true,
	"StringLessThanPath": true, "StringGreaterThanPath": true,
	"StringLessThanEqualsPath": true, "StringGreaterThanEqualsPath": true,
	"NumericEquals":     true,
	"NumericEqualsPath": true,
	"NumericLessThan":   true, "NumericGreaterThan": true,
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

// choiceComparatorValueKind classifies each documented JSONPath comparison
// operator by the JSON type its value must carry: "For each of these
// operators, the corresponding value must be of the appropriate type:
// string, number, Boolean, or timestamp." A timestamp literal is a string
// that must additionally conform to the ASL timestamp profile, and the
// Path-suffixed operators carry a path string ("For those operators that
// end with 'Path', the value MUST be a Path").
var choiceComparatorValueKind = map[string]string{
	"StringEquals": "string", "StringLessThan": "string",
	"StringGreaterThan": "string", "StringLessThanEquals": "string",
	"StringGreaterThanEquals": "string", "StringMatches": "string",
	"NumericEquals": "number", "NumericLessThan": "number",
	"NumericGreaterThan": "number", "NumericLessThanEquals": "number",
	"NumericGreaterThanEquals": "number",
	"BooleanEquals":            "boolean",
	"IsPresent":                "boolean", "IsNull": "boolean", "IsString": "boolean",
	"IsNumeric": "boolean", "IsBoolean": "boolean", "IsTimestamp": "boolean",
	"TimestampEquals": "timestamp", "TimestampLessThan": "timestamp",
	"TimestampGreaterThan": "timestamp", "TimestampLessThanEquals": "timestamp",
	"TimestampGreaterThanEquals": "timestamp",
	"StringEqualsPath":           "path", "StringLessThanPath": "path",
	"StringGreaterThanPath": "path", "StringLessThanEqualsPath": "path",
	"StringGreaterThanEqualsPath": "path",
	"NumericEqualsPath":           "path", "NumericLessThanPath": "path",
	"NumericGreaterThanPath": "path", "NumericLessThanEqualsPath": "path",
	"NumericGreaterThanEqualsPath": "path",
	"BooleanEqualsPath":            "path",
	"TimestampEqualsPath":          "path", "TimestampLessThanPath": "path",
	"TimestampGreaterThanPath": "path", "TimestampLessThanEqualsPath": "path",
	"TimestampGreaterThanEqualsPath": "path",
}

// validateStatesScope validates the states of one States object (the
// top-level scope or a Parallel branch / Map processor scope). The
// effective query language governs which fields are legal per state.
func (v *aslValidatorContext) validateStatesScope(parent map[string]interface{}, states map[string]interface{}, scope string) {
	topQL, _ := parent["QueryLanguage"].(string)
	if topQL == "" {
		// "the default query language for each state inside the Map's
		// ItemProcessor or the Parallel's Branches fields is the state
		// machine's query language, and is independent of the Map or
		// Parallel State's query language."
		topQL = v.machineQL
	}

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
		if stateQL == "JSONPath" && v.machineQL == "JSONata" {
			// "You cannot revert a top-level JSONata-based state machine
			// to a mix of JSONata and JSONPath states." — the incremental
			// upgrade direction (JSONPath machine, JSONata state) stays
			// legal.
			v.schemaf("/States/"+name+"/QueryLanguage",
				"A state in a JSONata state machine cannot override QueryLanguage to JSONPath")
			continue
		}
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

		v.validateOneState(name, stateType, stateMap, states, effectiveQL)
	}
}

// validateOneState validates a single state object against its type
// contract.
func (v *aslValidatorContext) validateOneState(name, stateType string, stateMap, states map[string]interface{}, queryLanguage string) {
	stateLoc := "/States/" + name
	jsonata := queryLanguage == "JSONata"

	// Query-language field separation: JSONata states reject the JSONPath
	// filters and vice versa.
	if jsonata {
		for _, f := range getJSONPathOnlyFields(stateType, stateMap) {
			v.schemaf(stateLoc+"/"+f, "State '%s' uses JSONata QueryLanguage but contains the JSONPath-only field '%s'", name, f)
		}
		v.validateJSONataExpressions(name, stateType, stateMap, stateLoc)
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
			pattern := ""
			trimmed := resource
			switch {
			case strings.HasSuffix(trimmed, ".waitForTaskToken"):
				pattern = "waitForTaskToken"
				trimmed = strings.TrimSuffix(trimmed, ".waitForTaskToken")
			case strings.HasSuffix(trimmed, ".sync:2"):
				pattern = "sync2"
				trimmed = strings.TrimSuffix(trimmed, ".sync:2")
			case strings.HasSuffix(trimmed, ".sync"):
				pattern = "sync"
				trimmed = strings.TrimSuffix(trimmed, ".sync")
			}
			if _, err := arnutil.ParseARN(trimmed); err != nil {
				v.add("ERROR", "INVALID_RESOURCE",
					fmt.Sprintf("The value of the Task state '%s' resource field is invalid: %s", name, resource),
					stateLoc+"/Resource")
			} else if msg := validateTaskResourceSupport(name, trimmed, pattern, resource); msg != "" {
				v.add("ERROR", "UNSUPPORTED_INTEGRATION", msg, stateLoc+"/Resource")
			}
			if pattern != "" && v.smType == "EXPRESS" {
				// "Express Workflows only support Request Response
				// integrations" — the .sync, .sync:2 and .waitForTaskToken
				// suffixes name patterns beyond Request Response and are
				// rejected at definition validation.
				v.schemaf(stateLoc+"/Resource",
					"Task state '%s' resource '%s' uses the '%s' integration pattern, but Express Workflows only support Request Response integrations",
					name, resource, pattern)
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
		v.validateMapState(name, stateMap, stateLoc, jsonata)
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
		fields := make([]string, 0, len(stateMap))
		for field := range stateMap {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		for _, field := range fields {
			if !known[field] {
				v.schemaf(stateLoc+"/"+field, "State '%s' of type %s declares the unknown field '%s'", name, stateType, field)
			}
		}
	}

	// Retry/Catch are legal members only on Task, Parallel and Map (the
	// census above); on every other type the unknown-field census already
	// rejects them, so validating their contents would double-diagnose one
	// defect.
	if aslStateFields[stateType]["Retry"] {
		v.validateRetryCatch(name, stateMap, states, stateLoc, queryLanguage)
	}
	v.validatePayloadTemplateKeys(name, stateMap, stateLoc, jsonata)
	v.scanPathShapeWarnings(name, stateMap, stateLoc)
}

// validatePayloadTemplateKeys rejects payload-template entries whose
// ".$"-suffixed key carries a non-string value: the suffix marks a path
// reference, so the value must be a path string. Creation time rejects
// what the runtime would otherwise fail with States.ParameterPathFailure.
// The JSONata surface (Arguments, Items) carries expressions instead of
// the ".$" convention and is excluded. A Distributed Map's nested
// ItemReader.Parameters and ResultWriter.Parameters are payload templates
// of the same shape and join the census through the nested walk.
func (v *aslValidatorContext) validatePayloadTemplateKeys(name string, stateMap map[string]interface{}, stateLoc string, jsonata bool) {
	if jsonata {
		return
	}
	for _, field := range []string{"Parameters", "ItemSelector", "ItemBatcher", "ResultSelector"} {
		raw, ok := stateMap[field]
		if !ok {
			continue
		}
		v.walkPayloadTemplateKeys(name, stateLoc+"/"+field, raw)
	}
	for _, field := range []string{"ItemReader", "ResultWriter"} {
		nested, ok := stateMap[field].(map[string]interface{})
		if !ok {
			continue
		}
		raw, ok := nested["Parameters"]
		if !ok {
			continue
		}
		v.walkPayloadTemplateKeys(name, stateLoc+"/"+field+"/Parameters", raw)
	}
}

// walkPayloadTemplateKeys recurses a payload-template value tree; map keys
// are visited in sorted order so diagnostics are deterministic.
func (v *aslValidatorContext) walkPayloadTemplateKeys(name, loc string, value interface{}) {
	switch val := value.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if strings.HasSuffix(k, ".$") {
				if _, isString := val[k].(string); !isString {
					v.schemaf(loc+"/"+k, "State '%s': the payload template key '%s' must carry a path string", name, k)
				}
			}
			v.walkPayloadTemplateKeys(name, loc+"/"+k, val[k])
		}
	case []interface{}:
		for i, entry := range val {
			v.walkPayloadTemplateKeys(name, loc+"/"+strconv.Itoa(i), entry)
		}
	}
}

// validateFailState ports the Fail field contract: each of Error/Cause
// excludes its reference path form (Fail state documentation). InputPath,
// OutputPath and Output are not in the Fail census, so the unknown-field
// walk already reports them — listing them here would double-diagnose one
// defect with two wordings.
func (v *aslValidatorContext) validateFailState(name string, stateMap map[string]interface{}, stateLoc string) {
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
	// A JSONata Wait carrying SecondsPath/TimestampPath is reported once,
	// by the JSONPath-only field separation in validateOneState.

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
		if _, valid := parseASLTimestamp(value); !valid {
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
			v.validateChoiceRule(name, rule, states, stateLoc, jsonata, false)
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

// validateChoiceRuleValue enforces the documented value typing on one
// comparison operator: the literal families take their JSON type (string,
// number or Boolean), a timestamp literal is a string conforming to the
// ASL timestamp profile, and the Path-suffixed operators take a path
// string.
func (v *aslValidatorContext) validateChoiceRuleValue(name, field, kind string, raw interface{}, stateLoc string) {
	switch kind {
	case "string":
		if _, ok := raw.(string); !ok {
			v.schemaf(stateLoc+"/Choices", "Choice rule operator '%s' in state '%s' must have a string value", field, name)
		}
	case "number":
		if _, ok := raw.(float64); !ok {
			v.schemaf(stateLoc+"/Choices", "Choice rule operator '%s' in state '%s' must have a numeric value", field, name)
		}
	case "boolean":
		if _, ok := raw.(bool); !ok {
			v.schemaf(stateLoc+"/Choices", "Choice rule operator '%s' in state '%s' must have a Boolean value", field, name)
		}
	case "timestamp":
		value, isString := raw.(string)
		if !isString {
			v.schemaf(stateLoc+"/Choices", "Choice rule operator '%s' in state '%s' must have a timestamp string value", field, name)
			return
		}
		if _, valid := parseASLTimestamp(value); !valid {
			v.schemaf(stateLoc+"/Choices",
				"Choice rule operator '%s' in state '%s' has timestamp %q outside the RFC3339 profile of ISO 8601 with an uppercase T and an uppercase Z or numeric offset, for example \"2024-03-14T01:59:00Z\"", field, name, value)
		}
	case "path":
		if _, ok := raw.(string); !ok {
			v.schemaf(stateLoc+"/Choices", "Choice rule operator '%s' in state '%s' must have a path string value", field, name)
		}
	}
}

// validateChoiceRule validates one Choice rule (recursively for the
// And/Or/Not combinators). nested marks rules inside a combinator: "The
// values of the And and Or operators must be non-empty arrays of Choice
// Rules that must not themselves contain Next fields. Likewise, the value
// of a Not operator must be a single Choice Rule that must not contain
// Next fields." — "the Next field can appear only in a top-level Choice
// Rule."
func (v *aslValidatorContext) validateChoiceRule(name string, rule, states map[string]interface{}, stateLoc string, jsonata, nested bool) {
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
	if combinators > 0 && comparators > 0 {
		// "A Choice Rule MUST be either a Boolean Expression or a
		// Data-test Expression" — a combinator and a comparator on the
		// same rule is neither.
		v.schemaf(stateLoc+"/Choices", "Choice rule for state '%s' must not combine And, Or or Not with a comparison operator", name)
	}
	if combinators == 0 && comparators != 1 {
		v.schemaf(stateLoc+"/Choices", "Choice rule for state '%s' must specify exactly one comparator", name)
	}
	if !jsonata {
		// A value whose JSON type contradicts its operator — or a
		// timestamp literal outside the RFC3339 profile — is a schema
		// failure at creation, not a rule that silently never matches at
		// run time. Sorted for a deterministic diagnostic order.
		fields := make([]string, 0, len(rule))
		for field := range rule {
			if _, isComparator := choiceComparatorValueKind[field]; isComparator {
				fields = append(fields, field)
			}
		}
		sort.Strings(fields)
		for _, field := range fields {
			v.validateChoiceRuleValue(name, field, choiceComparatorValueKind[field], rule[field], stateLoc)
		}
		// "The values of Variable and Path fields in a comparison must be
		// valid Reference Paths" — the type dimension of that contract is
		// a string value.
		if raw, hasVariable := rule["Variable"]; hasVariable {
			if _, isString := raw.(string); !isString {
				v.schemaf(stateLoc+"/Choices", "Choice rule 'Variable' in state '%s' must be a path string", name)
			}
		}
	}
	if jsonata {
		// "Variables and comparison fields are only available for
		// JSONPath. Condition is only available for JSONata." The
		// JSONPath-only members are rejected at every nesting level; the
		// Condition requirement applies to the leaf rules a combinator
		// wraps, not to the combinator itself.
		if _, hasVariable := rule["Variable"]; hasVariable {
			v.schemaf(stateLoc+"/Choices", "Choice rule in the JSONata state '%s' uses the JSONPath-only field 'Variable'", name)
		}
		for field := range rule {
			if aslChoiceComparators[field] && field != "Condition" {
				v.schemaf(stateLoc+"/Choices", "Choice rule in the JSONata state '%s' uses the JSONPath-only comparison field '%s'", name, field)
			}
		}
		if combinators == 0 {
			if _, hasCondition := rule["Condition"]; !hasCondition {
				v.schemaf(stateLoc+"/Choices", "Choice rules in the JSONata state '%s' must specify a Condition", name)
			}
		}
	}
	if !jsonata && comparators == 1 && combinators == 0 {
		// A comparison is two fields — the Variable and the operator —
		// including for the Is* tests.
		if _, hasVariable := rule["Variable"]; !hasVariable {
			v.schemaf(stateLoc+"/Choices", "Choice rule for state '%s' must specify a 'Variable'", name)
		}
	}

	if raw, ok := rule["Not"]; ok {
		nestedRule, isRule := raw.(map[string]interface{})
		if !isRule {
			v.schemaf(stateLoc+"/Choices", "The value of a Not operator in state '%s' must be a single Choice Rule", name)
		} else {
			v.validateChoiceRule(name, nestedRule, states, stateLoc, jsonata, true)
		}
	}
	for _, combinator := range []string{"And", "Or"} {
		if raw, ok := rule[combinator]; ok {
			rules, isArr := raw.([]interface{})
			if !isArr || len(rules) == 0 {
				v.schemaf(stateLoc+"/Choices", "The value of an %s operator in state '%s' must be a non-empty array of Choice Rules", combinator, name)
			}
			// "MUST be a non-empty object array of Choice Rules": an entry
			// of any other JSON type is a schema failure at creation, not
			// a skipped entry that dies at run-time decode.
			sawNonObject := false
			for _, raw := range rules {
				if nestedRule, ok := raw.(map[string]interface{}); ok {
					v.validateChoiceRule(name, nestedRule, states, stateLoc, jsonata, true)
				} else {
					sawNonObject = true
				}
			}
			if sawNonObject {
				v.schemaf(stateLoc+"/Choices", "The value of an %s operator in state '%s' must be a non-empty array of Choice Rules", combinator, name)
			}
		}
	}

	if raw, ok := rule["Next"]; ok {
		if nested {
			v.schemaf(stateLoc+"/Choices", "The 'Next' field can appear only in a top-level Choice Rule of state '%s'", name)
		} else {
			target, isString := raw.(string)
			if !isString {
				v.schemaf(stateLoc+"/Choices", "Choice rule 'Next' must be a state name string")
			} else if _, exists := states[target]; !exists {
				v.add("ERROR", "MISSING_TRANSITION_TARGET",
					fmt.Sprintf("The value of the 'Next' field '%s' in a rule of state '%s' does not match a state name in this scope", target, name),
					stateLoc+"/Choices")
			}
		}
	} else if !nested {
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

// aslRetrierFields is the retrier member census (the error-handling
// page's field list: ErrorEquals, IntervalSeconds, MaxAttempts,
// BackoffRate, MaxDelaySeconds, JitterStrategy).
var aslRetrierFields = map[string]bool{
	"ErrorEquals": true, "IntervalSeconds": true, "MaxAttempts": true,
	"BackoffRate": true, "MaxDelaySeconds": true, "JitterStrategy": true,
}

// aslCatcherFields is the catcher member census (ErrorEquals, Next,
// ResultPath for JSONPath states, Assign and Output for JSONata states).
var aslCatcherFields = map[string]bool{
	"ErrorEquals": true, "Next": true, "ResultPath": true,
	"Assign": true, "Output": true,
}

// checkStatesALLPlacement enforces the wildcard placement contract:
// "It must appear alone in the ErrorEquals array and must appear in the
// last retrier in the Retry array" — and the same sentence for the last
// catcher.
func (v *aslValidatorContext) checkStatesALLPlacement(equals []interface{}, loc, kind string, isLast bool) {
	for _, e := range equals {
		s, _ := e.(string)
		if s != "States.ALL" {
			continue
		}
		if len(equals) != 1 {
			v.schemaf(loc+"/ErrorEquals", "%s ErrorEquals containing States.ALL must list it alone", kind)
		}
		if !isLast {
			v.schemaf(loc+"/ErrorEquals", "States.ALL must appear in the last %s of the array", strings.ToLower(kind))
		}
	}
}

// validateRetryCatch validates the Retrier and Catcher contracts:
// non-empty ErrorEquals, documented ranges for the retry timing fields,
// in-scope Next targets for catchers, and the query-language member
// separation on catchers — a JSONata Catcher carries its transformation in
// Assign/Output ("The $states.errorOutput can be used in the Catch field's
// Assign or Output"), while a JSONPath Catcher folds the error through
// ResultPath, so each dialect rejects the other's members.
func (v *aslValidatorContext) validateRetryCatch(name string, stateMap, states map[string]interface{}, stateLoc, queryLanguage string) {
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
			} else {
				v.checkStatesALLPlacement(equals, retrierLoc, "Retrier", i == len(raws)-1)
			}
			for field := range retrier {
				if !aslRetrierFields[field] {
					v.schemaf(retrierLoc+"/"+field, "Retrier of state '%s' declares the unknown field '%s'", name, field)
				}
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
			} else {
				v.checkStatesALLPlacement(equals, catcherLoc, "Catcher", i == len(raws)-1)
			}
			for field := range catcher {
				if !aslCatcherFields[field] {
					v.schemaf(catcherLoc+"/"+field, "Catcher of state '%s' declares the unknown field '%s'", name, field)
				}
			}
			if queryLanguage == "JSONata" {
				if _, ok := catcher["ResultPath"]; ok {
					v.schemaf(catcherLoc+"/ResultPath", "Catchers of state '%s' use JSONata QueryLanguage but contain the JSONPath-only field 'ResultPath'", name)
				}
			} else {
				for _, field := range []string{"Assign", "Output"} {
					if _, ok := catcher[field]; ok {
						v.schemaf(catcherLoc+"/"+field, "Catchers of state '%s' use JSONPath QueryLanguage but contain the JSONata-only field '%s'", name, field)
					}
				}
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

// The Task-resource support matrix. AWS validates, at definition time,
// that the integration-pattern suffix and the API action are supported
// for the named service (the Step Functions developer guide's
// integration-patterns support table); the platform enforces the same
// combinations for the integration families it carries. Resources naming
// services the platform does not carry stay valid — they are valid AWS
// definitions — and fail at run time with a cause naming the unavailable
// integration.
//
// Patterns: "" is the plain Request Response form, "waitForTaskToken" the
// callback form, "sync"/"sync2" the Run-a-Job form and its JSON-output
// variant.

// supportedTaskPatterns lists the integration patterns each carried
// resource family accepts. The empty pattern (plain Request Response) is
// implied for every entry.
var supportedTaskPatterns = map[string]map[string]bool{
	"arn:aws:states:::lambda:invoke":         {"waitForTaskToken": true},
	"arn:aws:states:::sqs:sendMessage":       {"waitForTaskToken": true},
	"arn:aws:states:::sns:publish":           {"waitForTaskToken": true},
	"arn:aws:states:::events:putEvents":      {"waitForTaskToken": true},
	"arn:aws:states:::dynamodb:getItem":      {},
	"arn:aws:states:::dynamodb:putItem":      {},
	"arn:aws:states:::dynamodb:updateItem":   {},
	"arn:aws:states:::dynamodb:deleteItem":   {},
	"arn:aws:states:::states:startExecution": {"sync": true, "sync2": true, "waitForTaskToken": true},
}

// validateTaskResourceSupport reports an UNSUPPORTED_INTEGRATION message
// when a carried family's resource names an unsupported action or pattern,
// or when a pattern suffix decorates a resource form that takes none (a
// Lambda function ARN or an activity ARN). Resources outside the carried
// families return no message: the AWS SDK namespace accepts any
// service/action with the plain and callback patterns (Run a Job is not
// supported there), and unrecognised services keep their AWS-side
// validity, with run time reporting the unavailable integration.
func validateTaskResourceSupport(stateName, base, pattern, fullResource string) string {
	if pattern != "" {
		switch {
		case strings.HasPrefix(base, "arn:aws:lambda:"):
			return fmt.Sprintf("The Resource of Task state '%s' must not carry an integration-pattern suffix: %s", stateName, fullResource)
		case isActivityResource(base):
			return fmt.Sprintf("The Resource of Task state '%s' must not carry an integration-pattern suffix: %s", stateName, fullResource)
		}
	}

	if patterns, carried := supportedTaskPatterns[base]; carried {
		if pattern != "" && !patterns[pattern] {
			return fmt.Sprintf("The integration pattern of Task state '%s' is not supported for the resource: %s", stateName, fullResource)
		}
		return ""
	}

	if strings.HasPrefix(base, "arn:aws:states:::") {
		family := strings.TrimPrefix(base, "arn:aws:states:::")
		switch {
		case strings.HasPrefix(family, "sqs:"), strings.HasPrefix(family, "sns:"),
			strings.HasPrefix(family, "events:"), strings.HasPrefix(family, "dynamodb:"),
			strings.HasPrefix(family, "lambda:"), strings.HasPrefix(family, "states:"):
			return fmt.Sprintf("The resource of Task state '%s' is not a carried integration API: %s", stateName, fullResource)
		case strings.HasPrefix(family, "aws-sdk:"):
			// The AWS SDK namespace: any service and action with the plain
			// or callback pattern; Run a Job is not supported there.
			if pattern == "sync" || pattern == "sync2" {
				return fmt.Sprintf("The integration pattern of Task state '%s' is not supported for the resource: %s", stateName, fullResource)
			}
		}
	}
	return ""
}

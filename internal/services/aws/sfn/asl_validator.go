package sfn

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
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

// aslValidatorContext carries the cross-scope validation state: state
// names are unique across the entire state machine and Map labels are
// unique within the definition.
type aslValidatorContext struct {
	smType      string
	machineQL   string
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

	topFields := make([]string, 0, len(doc))
	for field := range doc {
		topFields = append(topFields, field)
	}
	sort.Strings(topFields)
	for _, field := range topFields {
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

	v.machineQL, _ = doc["QueryLanguage"].(string)
	v.validateStatesScope(doc, states, "")

	return v.diagnostics
}

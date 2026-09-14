package sfn

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// TestChoiceEmptyStringLiteralMatches pins that a literal empty string is a
// usable comparison operand: {"StringEquals": ""} matches an empty input
// value instead of being treated as an absent operator, and a rule list
// without Default no longer fails the execution over it.
func TestChoiceEmptyStringLiteralMatches(t *testing.T) {
	definition := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"Variable": "$.name", "StringEquals": "", "Next": "Empty"},
					{"Variable": "$.name", "StringEquals": "x", "Next": "Named"}
				]
			},
			"Empty":  {"Type": "Pass", "Result": "empty", "End": true},
			"Named":  {"Type": "Pass", "Result": "named", "End": true}
		}
	}`
	var def sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		t.Fatalf("unmarshal definition: %v", err)
	}
	_, states, err := parseDefinitionJSON(definition)
	if err != nil {
		t.Fatalf("parse definition: %v", err)
	}

	e := NewExecutor(newMapTestStore(t), nil)
	_, next, err := e.executeChoice(t.Context(), &ExecutionContext{
		Execution:     &sfnstore.Execution{ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:e"},
		Definition:    &def,
		CurrentState:  "C",
		Input:         `{"name":""}`,
		EventId:       ptrEventID(),
		States:        states,
		QueryLanguage: "JSONPath",
	}, states["C"].(*sfnstore.ChoiceState))
	if err != nil {
		t.Fatalf("the empty-string literal must evaluate as an operator: %v", err)
	}
	if next != "Empty" {
		t.Errorf("choice transition = %q, want Empty — the empty literal must match the empty value", next)
	}
}

// runChoice executes one Choice machine against the given input and
// returns the transition (or the execution error).
func runChoice(t *testing.T, definition, input string) (string, error) {
	t.Helper()
	var def sfnstore.StateMachineDefinition
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		t.Fatalf("unmarshal definition: %v", err)
	}
	_, states, err := parseDefinitionJSON(definition)
	if err != nil {
		t.Fatalf("parse definition: %v", err)
	}
	e := NewExecutor(newMapTestStore(t), nil)
	_, next, execErr := e.executeChoice(t.Context(), &ExecutionContext{
		Execution:     &sfnstore.Execution{ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:sm:e"},
		Definition:    &def,
		CurrentState:  "C",
		Input:         input,
		EventId:       ptrEventID(),
		States:        states,
		QueryLanguage: "JSONPath",
		VariableScope: NewVariableScope(nil),
		MapItemIndex:  -1,
	}, states["C"].(*sfnstore.ChoiceState))
	return next, execErr
}

// TestChoiceStringOperatorsRequireStringOperand pins the documented type
// discipline: "Step Functions doesn't attempt to match a numeric field to
// a string value" — a numeric (or Boolean) variable never stringifies into
// a match — while "because timestamp fields are logically strings, it's
// possible that a field considered to be a timestamp can be matched by a
// StringEquals comparator".
func TestChoiceStringOperatorsRequireStringOperand(t *testing.T) {
	definition := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"Variable": "$.foo", "StringEquals": "5", "Next": "Str"},
					{"Variable": "$.foo", "StringMatches": "5*", "Next": "Str"},
					{"Variable": "$.flag", "StringEquals": "true", "Next": "Str"}
				],
				"Default": "Other"
			},
			"Str":   {"Type": "Pass", "Result": "string", "End": true},
			"Other": {"Type": "Pass", "Result": "other", "End": true}
		}
	}`
	next, err := runChoice(t, definition, `{"foo":5,"flag":true}`)
	if err != nil {
		t.Fatalf("numeric/Boolean operands must evaluate as non-matching, not fail: %v", err)
	}
	if next != "Other" {
		t.Errorf("transition = %q, want Other — numeric 5 must not match StringEquals \"5\" and true must not match \"true\"", next)
	}

	// A timestamp-shaped string is a string: StringEquals matches it.
	tsDef := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"Variable": "$.at", "StringEquals": "2024-01-01T00:00:00Z", "Next": "Str"}
				],
				"Default": "Other"
			},
			"Str":   {"Type": "Pass", "Result": "string", "End": true},
			"Other": {"Type": "Pass", "Result": "other", "End": true}
		}
	}`
	next, err = runChoice(t, tsDef, `{"at":"2024-01-01T00:00:00Z"}`)
	if err != nil {
		t.Fatalf("timestamp string vs StringEquals: %v", err)
	}
	if next != "Str" {
		t.Errorf("transition = %q, want Str — a timestamp field is logically a string and matches StringEquals", next)
	}

	// The *Path comparator form carries the same discipline: a numeric
	// variable never matches a string comparator operand.
	pathDef := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"Variable": "$.foo", "StringEqualsPath": "$.bar", "Next": "Str"}
				],
				"Default": "Other"
			},
			"Str":   {"Type": "Pass", "Result": "string", "End": true},
			"Other": {"Type": "Pass", "Result": "other", "End": true}
		}
	}`
	next, err = runChoice(t, pathDef, `{"foo":5,"bar":"5"}`)
	if err != nil {
		t.Fatalf("numeric variable vs StringEqualsPath: %v", err)
	}
	if next != "Other" {
		t.Errorf("transition = %q, want Other — the *Path string operators must require a string variable", next)
	}
}

// TestChoiceUndefinedVariableFailsRuntime pins that a Choice Variable
// referencing an undefined workflow variable fails the execution — "If $x
// had not been previously assigned, the example would fail because $x
// would be undefined" — while an absent plain input path stays a
// non-match (the documented IsPresent pattern) and IsPresent on an
// unresolvable reference reports absent.
func TestChoiceUndefinedVariableFailsRuntime(t *testing.T) {
	definition := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"Variable": "$missing", "StringEquals": "x", "Next": "Str"}
				],
				"Default": "Other"
			},
			"Str":   {"Type": "Pass", "Result": "string", "End": true},
			"Other": {"Type": "Pass", "Result": "other", "End": true}
		}
	}`
	_, err := runChoice(t, definition, `{}`)
	if err == nil {
		t.Fatal("referencing the undefined variable $missing must fail the execution, not fall to Default")
	}
	if !strings.Contains(err.Error(), "States.Runtime") {
		t.Errorf("error = %v, want States.Runtime", err)
	}

	// An absent plain input path remains a plain non-match.
	absentDef := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"Variable": "$.missing.path", "StringEquals": "x", "Next": "Str"}
				],
				"Default": "Other"
			},
			"Str":   {"Type": "Pass", "Result": "string", "End": true},
			"Other": {"Type": "Pass", "Result": "other", "End": true}
		}
	}`
	next, err := runChoice(t, absentDef, `{"other":1}`)
	if err != nil {
		t.Fatalf("an absent input path must stay a non-match: %v", err)
	}
	if next != "Other" {
		t.Errorf("transition = %q, want Other", next)
	}

	// IsPresent observes absence without failing.
	presentDef := `{
		"StartAt": "C",
		"States": {
			"C": {
				"Type": "Choice",
				"Choices": [
					{"And": [
						{"Variable": "$.missing.path", "IsPresent": false},
						{"Variable": "$missing2", "IsPresent": false}
					], "Next": "Str"}
				],
				"Default": "Other"
			},
			"Str":   {"Type": "Pass", "Result": "string", "End": true},
			"Other": {"Type": "Pass", "Result": "other", "End": true}
		}
	}`
	next, err = runChoice(t, presentDef, `{"other":1}`)
	if err != nil {
		t.Fatalf("IsPresent=false must observe absence without failing: %v", err)
	}
	if next != "Str" {
		t.Errorf("transition = %q, want Str — both absent path and unresolvable variable report IsPresent false", next)
	}
}

// TestChoiceMistypedComparatorLiteralNeverMatches pins the decode
// uniformity the spec states: "For each operator which compares values, if
// the values are not both of the appropriate type (String, number,
// boolean, or Timestamp) the comparison will return false" — a mistyped
// literal decodes (as the operator-absent representation) and the rule
// never matches, on every operator family alike. The operators the spec
// types with an explicit MUST (StringMatches' value is a String, *Path
// values are Paths) keep rejecting mistyped values at decode.
func TestChoiceMistypedComparatorLiteralNeverMatches(t *testing.T) {
	cases := []struct {
		operator string
		value    string
		input    string
	}{
		{"StringEquals", `5`, `{"x":"s"}`},
		{"StringLessThan", `5`, `{"x":"s"}`},
		{"StringGreaterThan", `5`, `{"x":"s"}`},
		{"StringLessThanEquals", `5`, `{"x":"s"}`},
		{"StringGreaterThanEquals", `true`, `{"x":"s"}`},
		{"NumericEquals", `"7"`, `{"x":7}`},
		{"NumericLessThan", `"7"`, `{"x":7}`},
		{"NumericGreaterThan", `"7"`, `{"x":7}`},
		{"NumericLessThanEquals", `"7"`, `{"x":7}`},
		{"NumericGreaterThanEquals", `false`, `{"x":7}`},
		{"BooleanEquals", `"yes"`, `{"x":true}`},
		{"TimestampEquals", `5`, `{"x":"2024-01-01T00:00:00Z"}`},
		{"TimestampLessThan", `5`, `{"x":"2024-01-01T00:00:00Z"}`},
		{"TimestampGreaterThan", `5`, `{"x":"2024-01-01T00:00:00Z"}`},
		{"TimestampLessThanEquals", `5`, `{"x":"2024-01-01T00:00:00Z"}`},
		{"TimestampGreaterThanEquals", `true`, `{"x":"2024-01-01T00:00:00Z"}`},
	}
	for _, tc := range cases {
		t.Run(tc.operator, func(t *testing.T) {
			def := fmt.Sprintf(`{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[`+
				`{"Variable":"$.x","%s":%s,"Next":"Hit"}],"Default":"Other"},`+
				`"Hit":{"Type":"Pass","Result":"hit","End":true},`+
				`"Other":{"Type":"Pass","Result":"other","End":true}}}`, tc.operator, tc.value)
			next, err := runChoice(t, def, tc.input)
			if err != nil {
				t.Fatalf("%s with a mistyped literal must evaluate as a non-match, not fail: %v", tc.operator, err)
			}
			if next != "Other" {
				t.Fatalf("%s with a mistyped literal transitioned to %q, want Other (never match)", tc.operator, next)
			}
		})
	}

	// The explicit-MUST values stay strict: a non-string StringMatches
	// value and a non-string *Path value reject the definition at decode.
	strict := []struct {
		operator string
		value    string
	}{
		{"StringMatches", `5`},
		{"StringEqualsPath", `5`},
	}
	for _, tc := range strict {
		t.Run(tc.operator+"_stays_strict", func(t *testing.T) {
			def := fmt.Sprintf(`{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[`+
				`{"Variable":"$.x","%s":%s,"Next":"Hit"}],"Default":"Other"},`+
				`"Hit":{"Type":"Pass","Result":"hit","End":true},`+
				`"Other":{"Type":"Pass","Result":"other","End":true}}}`, tc.operator, tc.value)
			if _, _, perr := parseDefinitionJSON(def); perr == nil {
				t.Fatalf("a mistyped %s value must reject the definition — its value carries an explicit MUST", tc.operator)
			}
		})
	}
}

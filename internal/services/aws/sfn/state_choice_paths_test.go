package sfn

import (
	"context"
	"errors"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// runStatesForChoice drives executeStates over a JSONPath definition and
// returns the terminal error (nil when the machine completes).
func runStatesForChoice(t *testing.T, def *sfnstore.StateMachineDefinition, input string) (*ExecutionContext, error) {
	t.Helper()
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:ch-" + def.StartAt,
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "ch-" + def.StartAt,
		Status:          "RUNNING",
		Input:           input,
	}
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist execution failed: %v", err)
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	eventId := int64(0)
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: def.StartAt,
		Input: input, EventId: &eventId, States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
		VariableScope: NewVariableScope(nil),
	}
	return execCtx, e.executeStates(context.Background(), execCtx)
}

// choicePathDefinition builds a Choice machine with the given InputPath and
// OutputPath so the pins vary only the filtering fields.
func choicePathDefinition(inputPath, outputPath string) *sfnstore.StateMachineDefinition {
	pick := map[string]interface{}{
		"Type": "Choice",
		"Choices": []interface{}{
			map[string]interface{}{
				"Variable":      "$.val",
				"NumericEquals": 5,
				"Next":          "Done",
			},
		},
		"Default": "Other",
	}
	if inputPath != "" {
		pick["InputPath"] = inputPath
	}
	if outputPath != "" {
		pick["OutputPath"] = outputPath
	}
	return &sfnstore.StateMachineDefinition{
		StartAt: "Pick",
		States: map[string]interface{}{
			"Pick": pick,
			"Done": map[string]interface{}{"Type": "Pass", "End": true},
			"Other": map[string]interface{}{
				"Type":       "Pass",
				"Parameters": map[string]interface{}{"path": "default"},
				"End":        true,
			},
		},
	}
}

// TestChoiceInputPathFiltersRuleInput pins that Choice rules evaluate
// against the InputPath-filtered input: with InputPath $.data the rule
// Variable $.val resolves inside the filtered object, a match that is
// invisible against the raw state input.
func TestChoiceInputPathFiltersRuleInput(t *testing.T) {
	execCtx, err := runStatesForChoice(t, choicePathDefinition("$.data", ""), `{"data":{"val":5},"other":1}`)
	if err != nil {
		t.Fatalf("machine failed: %v", err)
	}
	if strings.Contains(execCtx.Output, "default") {
		t.Fatalf("output = %s, want the matched branch (the rule must see the InputPath-filtered input)", execCtx.Output)
	}
}

// TestChoiceOutputPathFiltersTransitionOutput pins that the OutputPath is
// applied to the effective input before the value reaches the next state.
func TestChoiceOutputPathFiltersTransitionOutput(t *testing.T) {
	execCtx, err := runStatesForChoice(t, choicePathDefinition("$.data", "$.val"), `{"data":{"val":5},"other":1}`)
	if err != nil {
		t.Fatalf("machine failed: %v", err)
	}
	if execCtx.Output != "5" {
		t.Fatalf("output = %s, want 5 (OutputPath $.val applied to the filtered input)", execCtx.Output)
	}
}

// TestChoiceDefaultRouteKeepsFiltering pins that a Default transition flows
// through the same InputPath/OutputPath filtering as a matched rule.
func TestChoiceDefaultRouteKeepsFiltering(t *testing.T) {
	execCtx, err := runStatesForChoice(t, choicePathDefinition("$.data", "$.val"), `{"data":{"val":9},"other":1}`)
	if err != nil {
		t.Fatalf("machine failed: %v", err)
	}
	if !strings.Contains(execCtx.Output, "default") {
		t.Fatalf("output = %s, want the Default branch", execCtx.Output)
	}
}

// TestChoiceNoMatchWithoutDefaultFailsNoChoiceMatched pins the spec error
// name: "States.NoChoiceMatched — A Choice State failed to find a match for
// the condition field extracted from its input."
func TestChoiceNoMatchWithoutDefaultFailsNoChoiceMatched(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Pick",
		States: map[string]interface{}{
			"Pick": map[string]interface{}{
				"Type": "Choice",
				"Choices": []interface{}{
					map[string]interface{}{
						"Variable":      "$.val",
						"NumericEquals": 5,
						"Next":          "Done",
					},
				},
			},
			"Done": map[string]interface{}{"Type": "Pass", "End": true},
		},
	}
	_, err := runStatesForChoice(t, def, `{"val":9}`)
	if err == nil {
		t.Fatal("a Choice state without a matching rule and without Default must fail")
	}
	var execErr *ExecutionError
	if !errors.As(err, &execErr) {
		t.Fatalf("expected *ExecutionError, got %T: %v", err, err)
	}
	if execErr.ErrorCode != "States.NoChoiceMatched" {
		t.Fatalf("error code = %s, want States.NoChoiceMatched", execErr.ErrorCode)
	}

	jsonataDef := &sfnstore.StateMachineDefinition{
		StartAt: "Pick",
		States: map[string]interface{}{
			"Pick": map[string]interface{}{
				"Type": "Choice",
				"Choices": []interface{}{
					map[string]interface{}{
						"Condition": "{% false %}",
						"Next":      "Done",
					},
				},
			},
			"Done": map[string]interface{}{"Type": "Pass", "End": true},
		},
	}
	_, err = runJSONataStates(t, jsonataDef, `{}`)
	if err == nil {
		t.Fatal("a JSONata Choice state without a matching condition and without Default must fail")
	}
	if !errors.As(err, &execErr) || execErr.ErrorCode != "States.NoChoiceMatched" {
		t.Fatalf("JSONata no-match error = %v, want States.NoChoiceMatched", err)
	}
}

// TestChoiceStateLevelAssign pins the documented two-level Assign
// contract: "You can use the Assign field at two levels in a Choice
// state: Top level … Inside a Choice Rule … These two placements are
// mutually exclusive at runtime. If a Choice Rule matches, Step Functions
// evaluates only that rule's Assign field and does not evaluate the
// top-level Assign field. Step Functions evaluates the top-level Assign
// only when no Choice Rule matches and the workflow transitions to the
// Default state."
func TestChoiceStateLevelAssign(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Pick",
		States: map[string]interface{}{
			"Pick": map[string]interface{}{
				"Type": "Choice",
				"Assign": map[string]interface{}{
					"path.$": "$.kind",
				},
				"Choices": []interface{}{
					map[string]interface{}{
						"Variable":     "$.kind",
						"StringEquals": "fast",
						"Next":         "Use",
						"Assign":       map[string]interface{}{"path": "rule"},
					},
				},
				"Default": "Use",
			},
			"Use": map[string]interface{}{
				"Type":       "Pass",
				"Parameters": map[string]interface{}{"taken.$": "$path"},
				"End":        true,
			},
		},
	}
	// Default route: the state-level Assign applies.
	execCtx, err := runStatesForChoice(t, def, `{"kind":"slow"}`)
	if err != nil {
		t.Fatalf("default-route machine failed: %v", err)
	}
	if execCtx.Output != `{"taken":"slow"}` {
		t.Fatalf("default-route output = %s, want the state-level Assign value", execCtx.Output)
	}

	// Matched rule: only the rule's Assign applies.
	execCtx, err = runStatesForChoice(t, def, `{"kind":"fast"}`)
	if err != nil {
		t.Fatalf("matched-rule machine failed: %v", err)
	}
	if execCtx.Output != `{"taken":"rule"}` {
		t.Fatalf("matched-rule output = %s, want the rule-level Assign value", execCtx.Output)
	}
}

// TestChoiceStateLevelAssignAcceptedByValidator pins the creation-time
// side: the page's own example (state-level Assign on a Choice) passes
// structure validation.
func TestChoiceStateLevelAssignAcceptedByValidator(t *testing.T) {
	def := `{
		"StartAt": "Check Score",
		"States": {
			"Check Score": {
				"Type": "Choice",
				"Assign": {"outputValue": "default path taken"},
				"Choices": [
					{"Variable": "$.score", "NumericGreaterThan": 90, "Assign": {"outputValue.$": "$.score"}, "Next": "High"},
					{"Variable": "$.score", "NumericGreaterThan": 50, "Next": "Medium"}
				],
				"Default": "Low"
			},
			"High": {"Type": "Succeed"},
			"Medium": {"Type": "Succeed"},
			"Low": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(def, "STANDARD"); err != nil {
		t.Fatalf("the documented state-level Choice Assign was rejected: %v", err)
	}
}

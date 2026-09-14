package sfn

import (
	"context"
	"errors"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// runJSONataStates drives executeStates over a JSONata definition and
// returns the terminal execution error (nil when the machine completes).
func runJSONataStates(t *testing.T, def *sfnstore.StateMachineDefinition, input string) (*ExecutionContext, error) {
	t.Helper()
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:qe-" + def.StartAt,
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "qe-" + def.StartAt,
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
		QueryLanguage: "JSONata", MapItemIndex: -1,
		VariableScope: NewVariableScope(nil),
	}
	return execCtx, e.executeStates(context.Background(), execCtx)
}

// wantQueryEvalError asserts the run terminated with
// States.QueryEvaluationError (unwrap-friendly: Pass, Choice, and Wait
// handlers return plain errors that executeStates wraps) and a cause
// mentioning causePart when it is non-empty.
func wantQueryEvalError(t *testing.T, execCtx *ExecutionContext, err error, causePart string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected a query evaluation error, machine completed with output %s", execCtx.Output)
	}
	var execErr *ExecutionError
	if !errors.As(err, &execErr) {
		t.Fatalf("expected *ExecutionError in the chain, got %T: %v", err, err)
	}
	if execErr.ErrorCode != "States.QueryEvaluationError" {
		t.Fatalf("error code = %s, want States.QueryEvaluationError", execErr.ErrorCode)
	}
	if causePart != "" && !strings.Contains(execErr.Cause, causePart) {
		t.Fatalf("cause %q does not mention %s", execErr.Cause, causePart)
	}
}

// TestJSONataUndefinedOutputFails pins the "failure to return a result"
// error class: "JSON cannot represent an undefined value expression, so the
// expression {% $data.thisFieldDoesNotExist %} would result in an error" —
// the state fails instead of degrading to JSON null.
func TestJSONataUndefinedOutputFails(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Emit",
		States: map[string]interface{}{
			"Emit": map[string]interface{}{
				"Type":   "Pass",
				"Output": map[string]interface{}{"city": "{% $states.input.shipTo.city %}"},
				"End":    true,
			},
		},
	}
	execCtx, err := runJSONataStates(t, def, `{"shipTo":{"zip":"123"}}`)
	wantQueryEvalError(t, execCtx, err, "undefined")
}

// TestJSONataTaskSecondsFieldValidation pins the "type incompatibility"
// and "value out of range" error classes for the seconds fields: a
// TimeoutSeconds or HeartbeatSeconds expression that evaluates to a
// non-numeric, fractional, or out-of-range value fails the state.
func TestJSONataTaskSecondsFieldValidation(t *testing.T) {
	cases := []struct {
		name     string
		field    string
		expr     string
		input    string
		location string
	}{
		{"TimeoutSeconds string result", "TimeoutSeconds", "{% $states.input.name %}", `{"name":"medium"}`, "TimeoutSeconds"},
		{"TimeoutSeconds fractional", "TimeoutSeconds", "{% 30.5 %}", `{}`, "TimeoutSeconds"},
		{"TimeoutSeconds zero", "TimeoutSeconds", "{% 0 %}", `{}`, "TimeoutSeconds"},
		{"TimeoutSeconds negative", "TimeoutSeconds", "{% -1 %}", `{}`, "TimeoutSeconds"},
		{"HeartbeatSeconds string result", "HeartbeatSeconds", "{% $states.input.name %}", `{"name":"medium"}`, "HeartbeatSeconds"},
		{"HeartbeatSeconds fractional", "HeartbeatSeconds", "{% 0.5 %}", `{}`, "HeartbeatSeconds"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := &sfnstore.StateMachineDefinition{
				StartAt: "Work",
				States: map[string]interface{}{
					"Work": map[string]interface{}{
						"Type":     "Task",
						"Resource": "arn:aws:states:::lambda:invoke",
						tc.field:   tc.expr,
						"End":      true,
					},
				},
			}
			execCtx, err := runJSONataStates(t, def, tc.input)
			wantQueryEvalError(t, execCtx, err, tc.location)
		})
	}
}

// TestQueryEvalSecondsValue unit-pins the seconds-field validation bounds:
// an integer in 1..MaxWaitSeconds is accepted, everything else fails.
func TestQueryEvalSecondsValue(t *testing.T) {
	if v, err := queryEvalSecondsValue("TimeoutSeconds", float64(45)); err != nil || v != 45 {
		t.Fatalf("45: got (%d, %v), want (45, nil)", v, err)
	}
	if v, err := queryEvalSecondsValue("TimeoutSeconds", float64(sfnstore.MaxWaitSeconds)); err != nil || v != sfnstore.MaxWaitSeconds {
		t.Fatalf("max: got (%d, %v), want (%d, nil)", v, err, sfnstore.MaxWaitSeconds)
	}
	for _, bad := range []interface{}{
		float64(0), float64(-1), float64(30.5), float64(sfnstore.MaxWaitSeconds + 1),
		"45", true, nil, map[string]interface{}{"s": 45},
	} {
		if _, err := queryEvalSecondsValue("TimeoutSeconds", bad); err == nil {
			t.Fatalf("value %v (%T): expected an error", bad, bad)
		}
	}
}

// TestJSONataFailExpressionFailureFailsExecution pins that a Fail state
// whose Error or Cause expression fails evaluation terminates the execution
// with States.QueryEvaluationError, not with a fabricated error name, and
// that a resolving expression still yields the Fail state's own error.
func TestJSONataFailExpressionFailureFailsExecution(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Die",
		States: map[string]interface{}{
			"Die": map[string]interface{}{
				"Type":  "Fail",
				"Error": "{% $states.input.missing.error %}",
			},
		},
	}
	execCtx, err := runJSONataStates(t, def, `{"other":1}`)
	wantQueryEvalError(t, execCtx, err, "undefined")

	healthy := &sfnstore.StateMachineDefinition{
		StartAt: "Die",
		States: map[string]interface{}{
			"Die": map[string]interface{}{
				"Type":  "Fail",
				"Error": "{% $states.input.code %}",
			},
		},
	}
	_, err = runJSONataStates(t, healthy, `{"code":"Custom"}`)
	execErr, ok := err.(*ExecutionError)
	if !ok {
		t.Fatalf("expected *ExecutionError, got %T: %v", err, err)
	}
	if execErr.ErrorCode != "Custom" {
		t.Fatalf("error code = %s, want Custom", execErr.ErrorCode)
	}
}

// TestJSONataChoiceNonBooleanConditionFails pins the type requirement on
// Choice conditions: a non-boolean evaluation fails the state instead of
// routing to the next rule or Default, while false conditions keep routing
// to Default.
func TestJSONataChoiceNonBooleanConditionFails(t *testing.T) {
	def := &sfnstore.StateMachineDefinition{
		StartAt: "Pick",
		States: map[string]interface{}{
			"Pick": map[string]interface{}{
				"Type": "Choice",
				"Choices": []interface{}{
					map[string]interface{}{
						"Condition": "{% $states.input.name %}",
						"Next":      "Went",
					},
				},
				"Default": "Fell",
			},
			"Went": map[string]interface{}{"Type": "Succeed", "Output": map[string]interface{}{"path": "choice"}},
			"Fell": map[string]interface{}{"Type": "Succeed", "Output": map[string]interface{}{"path": "default"}},
		},
	}
	execCtx, err := runJSONataStates(t, def, `{"name":"nope"}`)
	wantQueryEvalError(t, execCtx, err, "boolean")

	falsy := &sfnstore.StateMachineDefinition{
		StartAt: "Pick",
		States: map[string]interface{}{
			"Pick": map[string]interface{}{
				"Type": "Choice",
				"Choices": []interface{}{
					map[string]interface{}{
						"Condition": "{% false %}",
						"Next":      "Went",
					},
				},
				"Default": "Fell",
			},
			"Went": map[string]interface{}{"Type": "Succeed", "Output": map[string]interface{}{"path": "choice"}},
			"Fell": map[string]interface{}{"Type": "Succeed", "Output": map[string]interface{}{"path": "default"}},
		},
	}
	execCtx, err = runJSONataStates(t, falsy, `{}`)
	if err != nil {
		t.Fatalf("machine failed: %v", err)
	}
	if !strings.Contains(execCtx.Output, "default") {
		t.Fatalf("output = %s, want the Default branch", execCtx.Output)
	}
}

// TestJSONataWaitExpressionClassifications pins that Wait expression
// results of the wrong type or value fail as query-evaluation errors, the
// same error class the seconds and timestamp fields share.
func TestJSONataWaitExpressionClassifications(t *testing.T) {
	seconds := &sfnstore.StateMachineDefinition{
		StartAt: "Hold",
		States: map[string]interface{}{
			"Hold": map[string]interface{}{
				"Type":    "Wait",
				"Seconds": "{% $states.input.name %}",
				"End":     true,
			},
		},
	}
	execCtx, err := runJSONataStates(t, seconds, `{"name":"soon"}`)
	wantQueryEvalError(t, execCtx, err, "not an integer value")

	timestamp := &sfnstore.StateMachineDefinition{
		StartAt: "Hold",
		States: map[string]interface{}{
			"Hold": map[string]interface{}{
				"Type":      "Wait",
				"Timestamp": "{% $states.input.when %}",
				"End":       true,
			},
		},
	}
	execCtx, err = runJSONataStates(t, timestamp, `{"when":"not-a-timestamp"}`)
	wantQueryEvalError(t, execCtx, err, "RFC3339")
}

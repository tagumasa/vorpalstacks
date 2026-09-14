package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// mapDef wraps a single Map state in a definition rooted at it.
func mapDef(state map[string]interface{}) *sfnstore.StateMachineDefinition {
	return &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States:  map[string]interface{}{"M": state},
	}
}

// mapEchoProcessor builds an ItemProcessor whose Succeed state echoes the
// iteration input.
func mapEchoProcessor() map[string]interface{} {
	return map[string]interface{}{
		"StartAt": "W",
		"States": map[string]interface{}{
			"W": map[string]interface{}{"Type": "Succeed"},
		},
	}
}

// runMapItems runs the machine and decodes the Map result array.
func runMapItems(t *testing.T, def *sfnstore.StateMachineDefinition, input string, jsonata bool) []interface{} {
	t.Helper()
	var output string
	var err error
	if jsonata {
		var execCtx *ExecutionContext
		execCtx, err = runJSONataStates(t, def, input)
		output = execCtx.Output
	} else {
		var execCtx *ExecutionContext
		execCtx, err = runStatesForChoice(t, def, input)
		output = execCtx.Output
	}
	if err != nil {
		t.Fatalf("machine failed: %v", err)
	}
	var items []interface{}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		t.Fatalf("output %s is not a JSON array: %v", output, err)
	}
	return items
}

func pairItems(t *testing.T, items []interface{}, want map[string]float64) {
	t.Helper()
	if len(items) != len(want) {
		t.Fatalf("produced %d items, want %d", len(items), len(want))
	}
	seen := map[string]float64{}
	for _, it := range items {
		m, ok := it.(map[string]interface{})
		if !ok {
			t.Fatalf("item %v is not an object", it)
		}
		key, _ := m["key"].(string)
		value, _ := m["value"].(float64)
		seen[key] = value
	}
	for k, v := range want {
		if seen[k] != v {
			t.Fatalf("items = %v, want pair %s=%v", seen, k, v)
		}
	}
}

// TestJSONataMapDefaultItemsIteratesInput pins the omitted-ItemReader
// dataset contract: "if your dataset is a JSON array passed from a previous
// step in the workflow, the ItemReader field is omitted" and "Step
// Functions will iterate directly over the elements of an array, or the
// key-value pairs of a JSON object."
func TestJSONataMapDefaultItemsIteratesInput(t *testing.T) {
	state := map[string]interface{}{
		"Type":          "Map",
		"ItemProcessor": mapEchoProcessor(),
		"End":           true,
	}
	items := runMapItems(t, mapDef(state), `[1,2,3]`, true)
	if len(items) != 3 {
		t.Fatalf("array input produced %d items, want 3", len(items))
	}

	items = runMapItems(t, mapDef(state), `{"a":1,"b":2}`, true)
	pairItems(t, items, map[string]float64{"a": 1, "b": 2})
}

// TestJSONataMapItemsObjectForm pins the Items field contract: "The Map
// state Items field will accept a JSON array, a JSON object, or a JSONata
// expression that must evaluate to an array or object."
func TestJSONataMapItemsObjectForm(t *testing.T) {
	state := map[string]interface{}{
		"Type":          "Map",
		"Items":         "{% $states.input.byName %}",
		"ItemProcessor": mapEchoProcessor(),
		"End":           true,
	}
	items := runMapItems(t, mapDef(state), `{"byName":{"a":1,"b":2}}`, true)
	pairItems(t, items, map[string]float64{"a": 1, "b": 2})
}

// TestJSONataMapItemsScalarFails pins the boundary: an Items expression
// that evaluates to neither an array nor an object fails the state.
func TestJSONataMapItemsScalarFails(t *testing.T) {
	state := map[string]interface{}{
		"Type":          "Map",
		"Items":         "{% 42 %}",
		"ItemProcessor": mapEchoProcessor(),
		"End":           true,
	}
	_, err := runJSONataStates(t, mapDef(state), `{}`)
	if err == nil || !strings.Contains(err.Error(), "States.InvalidItems") {
		t.Fatalf("scalar Items error = %v, want States.InvalidItems", err)
	}
}

// TestJSONPathMapObjectInputIteratesPairs pins the JSONPath branch of the
// same sentence: "If the input to the Map state is a JSON object, it runs
// an iteration for each key-value pair in the object, passing the pair to
// the iteration as input" — with the default ItemsPath and with ItemsPath
// selecting a nested object.
func TestJSONPathMapObjectInputIteratesPairs(t *testing.T) {
	state := map[string]interface{}{
		"Type":          "Map",
		"ItemProcessor": mapEchoProcessor(),
		"End":           true,
	}
	items := runMapItems(t, mapDef(state), `{"a":1,"b":2}`, false)
	pairItems(t, items, map[string]float64{"a": 1, "b": 2})

	nested := map[string]interface{}{
		"Type":          "Map",
		"ItemsPath":     "$.servers",
		"ItemProcessor": mapEchoProcessor(),
		"End":           true,
	}
	items = runMapItems(t, mapDef(nested), `{"servers":{"web":1},"unused":0}`, false)
	pairItems(t, items, map[string]float64{"web": 1})
}

// TestMapIterationEventsNameParentStateWithLabel pins the
// MapIterationEventDetails name contract: "The name of the iteration's
// parent Map state" — a configured Label travels in the Map Run ARN alone
// and must not replace the state name in the iteration events.
func TestMapIterationEventsNameParentStateWithLabel(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:labelled",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "labelled",
		Status:          "RUNNING",
		Input:           `{"v":[1]}`,
	}
	def := mapDef(map[string]interface{}{
		"Type":      "Map",
		"Label":     "CustomLabel",
		"ItemsPath": "$.v",
		"ItemProcessor": map[string]interface{}{
			"StartAt": "W",
			"States": map[string]interface{}{
				"W": map[string]interface{}{"Type": "Succeed"},
			},
		},
		"End": true,
	})
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("persist execution failed: %v", err)
	}
	eventId := int64(0)
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M", Input: exec.Input,
		EventId: &eventId, States: map[string]sfnstore.State{}, QueryLanguage: "JSONPath",
		MapItemIndex: -1, VariableScope: NewVariableScope(nil),
	}
	states, err := extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	execCtx.States = states

	if _, _, execErr := e.executeMap(context.Background(), execCtx, execCtx.States["M"].(*sfnstore.MapState)); execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 100, "", false)
	if herr != nil {
		t.Fatalf("history failed: %v", herr)
	}
	iterationEvents := 0
	for _, ev := range history {
		if ev.MapIterationEventDetails == nil {
			continue
		}
		iterationEvents++
		if ev.MapIterationEventDetails.Name != "M" {
			t.Errorf("%s name = %q, want the parent Map state name %q — the Label belongs to the Map Run ARN only",
				ev.Type, ev.MapIterationEventDetails.Name, "M")
		}
	}
	if iterationEvents == 0 {
		t.Fatal("no MapIteration events recorded — the inline iteration framing is missing")
	}
}

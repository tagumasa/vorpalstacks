package sfn

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// evalIntrinsic is the test harness for the evaluator: it parses the state
// input as the data root and evaluates one invocation at nesting depth one.
func evalIntrinsic(t *testing.T, e *Executor, input, invocation string) (interface{}, error) {
	t.Helper()
	var data interface{}
	if input != "" {
		if err := json.Unmarshal([]byte(input), &data); err != nil {
			t.Fatalf("invalid test input JSON: %v", err)
		}
	}
	return e.evaluateIntrinsic("", invocation, data, 1)
}

func mustIntrinsic(t *testing.T, e *Executor, input, invocation string) interface{} {
	t.Helper()
	v, err := evalIntrinsic(t, e, input, invocation)
	if err != nil {
		t.Fatalf("evaluateIntrinsic(%s) failed: %v", invocation, err)
	}
	return v
}

// wantIntrinsicValue asserts a whole-result match after a JSON round-trip so
// number formatting (float64 1 marshals as 1) does not create noise.
func wantIntrinsicValue(t *testing.T, e *Executor, input, invocation, wantJSON string) {
	t.Helper()
	got := mustIntrinsic(t, e, input, invocation)
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal result failed: %v", err)
	}
	if string(gotJSON) != wantJSON {
		t.Errorf("evaluateIntrinsic(%s) = %s, want %s", invocation, gotJSON, wantJSON)
	}
}

func wantIntrinsicError(t *testing.T, e *Executor, input, invocation, wantPart string) {
	t.Helper()
	_, err := evalIntrinsic(t, e, input, invocation)
	if err == nil {
		t.Fatalf("evaluateIntrinsic(%s) unexpectedly succeeded", invocation)
	}
	if !strings.Contains(err.Error(), wantPart) {
		t.Errorf("evaluateIntrinsic(%s) error = %q, want it to contain %q", invocation, err, wantPart)
	}
}

// TestIntrinsicDocumentedExamples pins every function against its example in
// the intrinsic-functions documentation.
func TestIntrinsicDocumentedExamples(t *testing.T) {
	e := &Executor{}

	t.Run("Array", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"Id":123456}`, `States.Array($.Id)`, `[123456]`)
	})
	t.Run("ArrayPartition", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputArray":[1,2,3,4,5,6,7,8,9]}`,
			`States.ArrayPartition($.inputArray,4)`, `[[1,2,3,4],[5,6,7,8],[9]]`)
	})
	t.Run("ArrayContains", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputArray":[1,2,3,4,5,6,7,8,9],"lookingFor":5}`,
			`States.ArrayContains($.inputArray, $.lookingFor)`, `true`)
	})
	t.Run("ArrayRange", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{}`, `States.ArrayRange(1, 9, 2)`, `[1,3,5,7,9]`)
	})
	t.Run("ArrayGetItem", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputArray":[1,2,3,4,5,6,7,8,9],"index":5}`,
			`States.ArrayGetItem($.inputArray, $.index)`, `6`)
	})
	t.Run("ArrayLength", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputArray":[1,2,3,4,5,6,7,8,9]}`,
			`States.ArrayLength($.inputArray)`, `9`)
	})
	t.Run("ArrayUnique", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputArray":[1,2,3,3,3,3,3,3,4]}`,
			`States.ArrayUnique($.inputArray)`, `[1,2,3,4]`)
	})
	t.Run("Base64Encode", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"input":"Data to encode"}`,
			`States.Base64Encode($.input)`, `"RGF0YSB0byBlbmNvZGU="`)
	})
	t.Run("Base64Decode", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"base64":"RGF0YSB0byBlbmNvZGU="}`,
			`States.Base64Decode($.base64)`, `"Data to encode"`)
	})
	t.Run("Hash", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"Data":"input data","Algorithm":"SHA-1"}`,
			`States.Hash($.Data, $.Algorithm)`, `"aaff4a450a104cd177d28d18d74485e8cae074b7"`)
	})
	t.Run("JsonMerge", func(t *testing.T) {
		wantIntrinsicValue(t, e,
			`{"json1":{"a":{"a1":1,"a2":2},"b":2},"json2":{"a":{"a3":1,"a4":2},"c":3}}`,
			`States.JsonMerge($.json1, $.json2, false)`,
			`{"a":{"a3":1,"a4":2},"b":2,"c":3}`)
	})
	t.Run("StringToJson", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"escapedJsonString":"{\"foo\": \"bar\"}"}`,
			`States.StringToJson($.escapedJsonString)`, `{"foo":"bar"}`)
	})
	t.Run("JsonToString", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"unescapedJson":{"foo":"bar"}}`,
			`States.JsonToString($.unescapedJson)`, `"{\"foo\":\"bar\"}"`)
	})
	t.Run("MathRandomInRange", func(t *testing.T) {
		v := mustIntrinsic(t, e, `{"start":1,"end":999}`, `States.MathRandom($.start, $.end)`)
		n, ok := v.(float64)
		if !ok || n < 1 || n >= 999 {
			t.Errorf("MathRandom(1,999) = %v, want 1 <= n < 999", v)
		}
	})
	t.Run("MathAdd", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"value1":111,"step":-1}`,
			`States.MathAdd($.value1, $.step)`, `110`)
	})
	t.Run("StringSplit single delimiter", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputString":"1,2,3,4,5","splitter":","}`,
			`States.StringSplit($.inputString, $.splitter)`, `["1","2","3","4","5"]`)
	})
	t.Run("StringSplit delimiter set", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"inputString":"This.is+a,test=string","splitter":".+,="}`,
			`States.StringSplit($.inputString, $.splitter)`,
			`["This","is","a","test","string"]`)
	})
	t.Run("UUID", func(t *testing.T) {
		v := mustIntrinsic(t, e, "", `States.UUID()`)
		s, ok := v.(string)
		if !ok {
			t.Fatalf("UUID() = %T, want string", v)
		}
		// Version 4 UUID: 4 marks the version, [89ab] the variant bits.
		pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
		if !pattern.MatchString(s) {
			t.Errorf("UUID() = %q, not a v4 UUID", s)
		}
	})
	t.Run("Format literal template", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"name":"Arnav"}`,
			`States.Format('Hello, my name is {}.', $.name)`,
			`"Hello, my name is Arnav."`)
	})
	t.Run("Format path template", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"name":"Arnav","template":"Hello, my name is {}."}`,
			`States.Format($.template, $.name)`,
			`"Hello, my name is Arnav."`)
	})
}

// TestIntrinsicValidationErrors pins the documented argument validation rules.
func TestIntrinsicValidationErrors(t *testing.T) {
	e := &Executor{}

	t.Run("ArrayPartition non-array first argument", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"v":1}`, `States.ArrayPartition($.v,4)`, "must be an array")
	})
	t.Run("ArrayPartition zero chunk", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"v":[1,2]}`, `States.ArrayPartition($.v,0)`, "non-zero positive integer")
	})
	t.Run("ArrayPartition non-integer chunk rounds", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"v":[1,2,3,4,5,6,7,8,9]}`,
			`States.ArrayPartition($.v,3.7)`, `[[1,2,3,4],[5,6,7,8],[9]]`)
	})
	t.Run("ArrayContains non-array", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"v":1}`, `States.ArrayContains($.v, 1)`, "must be an array")
	})
	t.Run("ArrayRange zero step", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.ArrayRange(1, 9, 0)`, "non-zero")
	})
	t.Run("ArrayRange descending", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{}`, `States.ArrayRange(9, 1, -2)`, `[9,7,5,3,1]`)
	})
	t.Run("ArrayRange over 1000 elements", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.ArrayRange(1, 1001, 1)`, "1000 elements")
	})
	t.Run("ArrayGetItem out of range", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"v":[1,2]}`, `States.ArrayGetItem($.v, 5)`, "out of range")
	})
	t.Run("Hash unknown algorithm", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"Data":"x","Algorithm":"SHA-333"}`,
			`States.Hash($.Data, $.Algorithm)`, "MD5, SHA-1, SHA-256, SHA-384 or SHA-512")
	})
	t.Run("JsonMerge deep true rejected", func(t *testing.T) {
		wantIntrinsicError(t, e,
			`{"json1":{"a":1},"json2":{"b":2}}`,
			`States.JsonMerge($.json1, $.json2, true)`, "boolean false")
	})
	t.Run("JsonMerge non-object", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"json1":1,"json2":{"b":2}}`,
			`States.JsonMerge($.json1, $.json2, false)`, "must be an object")
	})
	t.Run("StringToJson invalid JSON", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"s":"not json {"}`, `States.StringToJson($.s)`, "not valid JSON")
	})
	t.Run("MathAdd int32 overflow", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.MathAdd(2147483647, 1)`, "2147483648")
	})
	t.Run("MathRandom empty range", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.MathRandom(5, 5)`, "greater than start")
	})
	t.Run("MathAdd below int32 floor", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.MathAdd(-2147483648, -1)`, "2147483648")
	})
	t.Run("Base64Decode invalid", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"s":"!!!"}`, `States.Base64Decode($.s)`, "not a valid Base64")
	})
	t.Run("Format placeholder count mismatch", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"n":1}`, `States.Format('{} {}', $.n)`, "placeholders")
	})
	t.Run("Format non-string template", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"n":1}`, `States.Format($.n, $.n)`, "must be a string")
	})
	t.Run("unknown function", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.Nope(1)`, "not a supported intrinsic function")
	})
	t.Run("missing path argument", func(t *testing.T) {
		wantIntrinsicError(t, e, `{"a":1}`, `States.ArrayLength($.absent)`, "selected no value")
	})
}

// TestIntrinsicLimits pins the documented size bounds: the 10,000 character
// data-string ceiling and the 256 KiB array payload ceiling.
func TestIntrinsicLimits(t *testing.T) {
	e := &Executor{}

	t.Run("Base64Encode over 10000 characters", func(t *testing.T) {
		big := strings.Repeat("x", 10001)
		input := `{"s":"` + big + `"}`
		wantIntrinsicError(t, e, input, `States.Base64Encode($.s)`, "10000 characters")
	})
	t.Run("Base64Encode at 10000 characters", func(t *testing.T) {
		big := strings.Repeat("x", 10000)
		input := `{"s":"` + big + `"}`
		_ = mustIntrinsic(t, e, input, `States.Base64Encode($.s)`)
	})
	t.Run("ArrayPartition over payload limit", func(t *testing.T) {
		big := make([]interface{}, 40000)
		for i := range big {
			big[i] = "0123456789"
		}
		data := map[string]interface{}{"v": big}
		if _, err := e.evaluateIntrinsic("", `States.ArrayPartition($.v,4)`, data, 1); err == nil {
			t.Fatal("ArrayPartition over the payload limit unexpectedly succeeded")
		} else if !strings.Contains(err.Error(), "payload limit") {
			t.Errorf("error = %q, want payload limit", err)
		}
	})
}

// TestIntrinsicEscapes pins the reserved-character escape grammar: \' \{ \}
// and \\ decode to their literals, any other backslash is an open escape
// backslash runtime error, and escaped braces in a States.Format template
// never act as placeholders.
func TestIntrinsicEscapes(t *testing.T) {
	e := &Executor{}

	t.Run("quote escape", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{}`, `States.Array('don\'t')`, `["don't"]`)
	})
	t.Run("escaped braces as literal values", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{}`, `States.Array('\{', '\}')`, `["{","}"]`)
	})
	t.Run("escaped backslash", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{}`, `States.Array('a\\b')`, `["a\\b"]`)
	})
	t.Run("escaped braces in Format template stay literal", func(t *testing.T) {
		wantIntrinsicValue(t, e, `{"x":"V"}`,
			`States.Format('literal \{\} and {}', $.x)`,
			`"literal {} and V"`)
	})
	t.Run("open escape backslash in string argument", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.Array('a\q')`, "open escape backslash")
	})
	t.Run("trailing backslash escapes the closing quote", func(t *testing.T) {
		wantIntrinsicError(t, e, `{}`, `States.Array('a\')`, "unterminated string literal")
	})
}

// TestIntrinsicNesting pins the documented ten-level nesting bound.
func TestIntrinsicNesting(t *testing.T) {
	e := &Executor{}
	data := map[string]interface{}{}

	build := func(depth int) string {
		inv := `States.Array(1,2,3)`
		for i := 1; i < depth; i++ {
			inv = `States.ArrayUnique(` + inv + `)`
		}
		return inv
	}

	gotJSON, _ := json.Marshal(mustIntrinsic(t, e, "", build(10)))
	if string(gotJSON) != `[1,2,3]` {
		t.Errorf("depth 10 chain = %s, want [1,2,3]", gotJSON)
	}
	if _, err := e.evaluateIntrinsic("", build(11), data, 1); err == nil {
		t.Fatal("depth 11 chain unexpectedly succeeded")
	} else if !strings.Contains(err.Error(), "nesting") {
		t.Errorf("depth 11 error = %q, want nesting", err)
	}
}

// TestIntrinsicMathRandomSeed pins that the optional seed makes draws
// reproducible.
func TestIntrinsicMathRandomSeed(t *testing.T) {
	e := &Executor{}
	first := mustIntrinsic(t, e, `{"start":1,"end":999}`, `States.MathRandom($.start, $.end, 42)`)
	second := mustIntrinsic(t, e, `{"start":1,"end":999}`, `States.MathRandom($.start, $.end, 42)`)
	if first != second {
		t.Errorf("seeded draws differ: %v vs %v", first, second)
	}
	third := mustIntrinsic(t, e, `{"start":1,"end":999}`, `States.MathRandom($.start, $.end, 43)`)
	if first == third {
		t.Errorf("different seeds produced the same draw %v", first)
	}
}

// TestApplyParametersIntrinsics pins the payload-template integration: the
// documented examples of Pass/Task Parameters with intrinsics.
func TestApplyParametersIntrinsics(t *testing.T) {
	e := &Executor{}
	params := &sfnstore.Parameters{Values: map[string]interface{}{
		"greeting.$": "States.Format('Hello, my name is {}.', $.name)",
		"buildId.$":  "States.Array($.Id)",
		"nested": map[string]interface{}{
			"merged.$":    "States.JsonMerge($.base, $.overlay, false)",
			"split.$":     "States.StringSplit('a.b,c', '.,')",
			"arrayPart.$": "States.ArrayPartition($.nine, 4)",
		},
	}}
	out, evalErr := e.applyParameters("", `{"name":"Arnav","Id":123456,"base":{"x":1},"overlay":{"y":2},"nine":[1,2,3,4,5,6,7,8,9]}`, params)
	if evalErr != nil {
		t.Fatalf("applyParameters failed: %v", evalErr.Cause)
	}
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output not JSON: %v", err)
	}
	if got["greeting"] != "Hello, my name is Arnav." {
		t.Errorf("greeting = %v", got["greeting"])
	}
	if buildID, ok := got["buildId"].([]interface{}); !ok || len(buildID) != 1 || buildID[0] != float64(123456) {
		t.Errorf("buildId = %v", got["buildId"])
	}
	nested, ok := got["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("nested = %v", got["nested"])
	}
	merged, ok := nested["merged"].(map[string]interface{})
	if !ok || merged["x"] != float64(1) || merged["y"] != float64(2) {
		t.Errorf("merged = %v", nested["merged"])
	}
	if split, ok := nested["split"].([]interface{}); !ok || len(split) != 3 || split[0] != "a" || split[1] != "b" || split[2] != "c" {
		t.Errorf("split = %v", nested["split"])
	}
	if parts, ok := nested["arrayPart"].([]interface{}); !ok || len(parts) != 3 {
		t.Errorf("arrayPart = %v", nested["arrayPart"])
	}
}

// TestApplyParametersIntrinsicFailure pins the failure classification: an
// invalid intrinsic in Parameters surfaces as States.Runtime.
func TestApplyParametersIntrinsicFailure(t *testing.T) {
	e := &Executor{}
	params := &sfnstore.Parameters{Values: map[string]interface{}{
		"bad.$": "States.ArrayRange(1, 9, 0)",
	}}
	_, evalErr := e.applyParameters("", `{}`, params)
	if evalErr == nil {
		t.Fatal("invalid intrinsic unexpectedly succeeded")
	}
	if evalErr.ErrorCode != "States.Runtime" {
		t.Errorf("error code = %s, want States.Runtime", evalErr.ErrorCode)
	}
	if !strings.Contains(evalErr.Cause, "Parameters") {
		t.Errorf("cause = %q, want the Parameters field name", evalErr.Cause)
	}
}

// TestApplyResultSelectorArrayRoot pins that ResultSelector data may be an
// array (the Parallel/Map combined result, an array Task result): paths use
// numeric segments and intrinsics see the array as the root.
func TestApplyResultSelectorArrayRoot(t *testing.T) {
	e := &Executor{}
	selector := &sfnstore.ResultSelector{Fields: map[string]interface{}{
		"first.$":  "$.[0].total",
		"second.$": "$.1.total",
		"len.$":    "States.ArrayLength($)",
	}}
	out, evalErr := e.applyResultSelector(`[{"total":1},{"total":2}]`, selector, "")
	if evalErr != nil {
		t.Fatalf("applyResultSelector failed: %v", evalErr.Cause)
	}
	want := `{"first":1,"len":2,"second":2}`
	if out != want {
		t.Errorf("applyResultSelector = %s, want %s", out, want)
	}
}

// TestExecutePassIntrinsicsThroughEngine runs a Pass state whose Parameters
// use intrinsics end to end through executeStates.
func TestExecutePassIntrinsicsThroughEngine(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "P",
		States: map[string]interface{}{
			"P": map[string]interface{}{
				"Type": "Pass",
				"Parameters": map[string]interface{}{
					"message.$":  "States.Format('{} world', $.word)",
					"uuidLike.$": "States.Base64Encode($.word)",
				},
				"ResultPath": "$",
				"End":        true,
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/intp1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "intp1",
		Status:          "RUNNING",
		Input:           `{"word":"hello"}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "P",
		Input: exec.Input, EventId: ptrEventID(), States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	if err := e.executeStates(context.Background(), execCtx); err != nil {
		t.Fatalf("executeStates failed: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(execCtx.Output), &got); err != nil {
		t.Fatalf("output not JSON: %v", err)
	}
	if got["message"] != "hello world" {
		t.Errorf("message = %v", got["message"])
	}
	if got["uuidLike"] != "aGVsbG8=" {
		t.Errorf("uuidLike = %v", got["uuidLike"])
	}
}

// TestExecuteFailIntrinsicPaths pins the Fail state contract that ErrorPath
// and CausePath may carry intrinsic invocations returning a string.
func TestExecuteFailIntrinsicPaths(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "F",
		States: map[string]interface{}{
			"F": map[string]interface{}{
				"Type":      "Fail",
				"ErrorPath": "States.Format('Code-{}', $.n)",
				"CausePath": "States.Format('failed at {}', $.n)",
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/intf1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "intf1",
		Status:          "RUNNING",
		Input:           `{"n":7}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "F",
		Input: exec.Input, EventId: ptrEventID(), States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	err = e.executeStates(context.Background(), execCtx)
	if err == nil {
		t.Fatal("Fail state must fail the execution")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Code-7") || !strings.Contains(msg, "failed at 7") {
		t.Errorf("Fail error = %q, want intrinsic-resolved Code-7 / failed at 7", msg)
	}
}

// TestParallelParametersAndResultSelector pins that a Parallel state applies
// its Parameters payload template to the branch input and its ResultSelector
// to the combined branch output array.
func TestParallelParametersAndResultSelector(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "Par",
		States: map[string]interface{}{
			"Par": map[string]interface{}{
				"Type": "Parallel",
				"Parameters": map[string]interface{}{
					"word.$": "States.Format('{}-par', $.w)",
				},
				"Branches": []*sfnstore.StateMachineDefinition{
					{
						StartAt: "B1",
						States: map[string]interface{}{
							"B1": map[string]interface{}{
								"Type": "Pass", "ResultPath": "$.b1", "End": true,
							},
						},
					},
					{
						StartAt: "B2",
						States: map[string]interface{}{
							"B2": map[string]interface{}{
								"Type": "Pass", "ResultPath": "$.b2", "End": true,
							},
						},
					},
				},
				"ResultSelector": map[string]interface{}{
					"firstWord.$": "$.[0].word",
					"count.$":     "States.ArrayLength($)",
				},
				"ResultPath": "$.combined",
				"End":        true,
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/intpar1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "intpar1",
		Status:          "RUNNING",
		Input:           `{"w":"go"}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "Par",
		Input: exec.Input, EventId: ptrEventID(), States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	parOutput, _, execErr := e.executeParallel(context.Background(), execCtx, states["Par"].(*sfnstore.ParallelState))
	if execErr != nil {
		t.Fatalf("executeParallel failed: %v", execErr.Cause)
	}
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(parOutput), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, parOutput)
	}
	combined, ok := got["combined"].(map[string]interface{})
	if !ok {
		t.Fatalf("combined = %v", got["combined"])
	}
	if combined["firstWord"] != "go-par" {
		t.Errorf("firstWord = %v", combined["firstWord"])
	}
	if combined["count"] != float64(2) {
		t.Errorf("count = %v", combined["count"])
	}
	// ResultPath folds into the post-Parameters input, so the transformed
	// word survives next to the combined selector result.
	if got["word"] != "go-par" {
		t.Errorf("post-Parameters input lost under ResultPath: %v", got["word"])
	}
}

// TestMapLegacyParametersAndResultSelector pins that a Map state written
// against the legacy Parameters name (the pre-ItemSelector name of the
// per-item payload template) still shapes item input, and that the Map
// ResultSelector filters the result array.
func TestMapLegacyParametersAndResultSelector(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	def := &sfnstore.StateMachineDefinition{
		StartAt: "M",
		States: map[string]interface{}{
			"M": map[string]interface{}{
				"Type":      "Map",
				"ItemsPath": "$.items",
				"Parameters": map[string]interface{}{
					"item.$": "$.value",
					"size.$": "States.ArrayLength($.value.tags)",
				},
				"Iterator": &sfnstore.StateMachineDefinition{
					StartAt: "W",
					States: map[string]interface{}{
						"W": map[string]interface{}{
							"Type": "Pass", "ResultPath": "$.doubled", "End": true,
						},
					},
				},
				"ResultSelector": map[string]interface{}{
					"firstSize.$": "$.[0].size",
					"count.$":     "States.ArrayLength($)",
				},
				"ResultPath": "$.combined",
				"End":        true,
			},
		},
	}
	states, err := e.extractStatesFromDefinition(def)
	if err != nil {
		t.Fatalf("extract states failed: %v", err)
	}
	exec := &sfnstore.Execution{
		ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm/intm1",
		StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		Name:            "intm1",
		Status:          "RUNNING",
		Input:           `{"items":[{"value":{"n":1,"tags":["a","b"]}},{"value":{"n":2,"tags":["c"]}}]}`,
	}
	execCtx := &ExecutionContext{
		Execution: exec, Definition: def, CurrentState: "M",
		Input: exec.Input, EventId: ptrEventID(), States: states,
		QueryLanguage: "JSONPath", MapItemIndex: -1,
	}
	mapOutput, _, execErr := e.executeMap(context.Background(), execCtx, states["M"].(*sfnstore.MapState))
	if execErr != nil {
		t.Fatalf("executeMap failed: %v", execErr.Cause)
	}
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(mapOutput), &got); err != nil {
		t.Fatalf("output not JSON: %v (%s)", err, mapOutput)
	}
	combined, ok := got["combined"].(map[string]interface{})
	if !ok {
		t.Fatalf("combined = %v", got["combined"])
	}
	if combined["firstSize"] != float64(2) {
		t.Errorf("firstSize = %v", combined["firstSize"])
	}
	if combined["count"] != float64(2) {
		t.Errorf("count = %v", combined["count"])
	}
}

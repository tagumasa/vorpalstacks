package sfn

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func newTestStateStore(t *testing.T) *sfnstore.StepFunctionStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")
}

const taskDefForMock = `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","Next":"Done"},"Done":{"Type":"Succeed"}}}`

func TestTestStateMockResultBypassesInvocation(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Mock: &TestStateMock{
			Result:         `{"value":42}`,
			ResultProvided: true,
		},
	})
	if err != nil {
		t.Fatalf("mocked run: %v", err)
	}
	if resp["status"] != "SUCCEEDED" {
		t.Fatalf("status = %v, want SUCCEEDED", resp["status"])
	}
	if resp["output"] != `{"value":42}` {
		t.Fatalf("output = %v, want the mocked result", resp["output"])
	}
	if resp["nextState"] != "Done" {
		t.Fatalf("nextState = %v, want Done", resp["nextState"])
	}
}

// TestTestStateMockJSONataTaskOutputTemplateApplies pins that the mocked
// output tail resolves a JSONata Task state's Output template through the
// same lazy resolver the real executor uses — a mocked run must not skip
// the Output member that a real run applies.
func TestTestStateMockJSONataTaskOutputTemplateApplies(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"QueryLanguage":"JSONata","StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","Output":{"doubled":"{% $states.input.value * 2 %}"},"End":true}}}`
	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "T",
		Input:      `{"value":21}`,
		Mock: &TestStateMock{
			Result:         `{"value":21}`,
			ResultProvided: true,
		},
	})
	if err != nil {
		t.Fatalf("mocked run: %v", err)
	}
	if resp["output"] != `{"doubled":42}` {
		t.Fatalf("output = %v, want the Output template resolved against the mocked result", resp["output"])
	}
}

func TestTestStateMockResultPathApplies(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","ResultPath":"$.wrapped","End":true}}}`
	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "T",
		Input:      `{"keep":1}`,
		Mock: &TestStateMock{
			Result:         `{"value":42}`,
			ResultProvided: true,
		},
	})
	if err != nil {
		t.Fatalf("mocked run: %v", err)
	}
	if resp["output"] != `{"keep":1,"wrapped":{"value":42}}` {
		t.Fatalf("output = %v, want the mocked result placed at $.wrapped", resp["output"])
	}
}

func TestTestStateMockErrorOutput(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Mock: &TestStateMock{
			Error:         "CustomError",
			Cause:         "mocked failure",
			ErrorProvided: true,
		},
	})
	if err != nil {
		t.Fatalf("mocked run: %v", err)
	}
	if resp["status"] != "FAILED" {
		t.Fatalf("status = %v, want FAILED", resp["status"])
	}
	if resp["error"] != "CustomError" || resp["cause"] != "mocked failure" {
		t.Fatalf("error/cause = %v/%v, want the mocked pair", resp["error"], resp["cause"])
	}
}

func TestTestStateMockRejectedOnNonMockableState(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":{"a":1},"End":true}}}`
	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "P",
		Mock: &TestStateMock{
			Result:         `{"value":42}`,
			ResultProvided: true,
		},
	})
	requireAWSCode(t, err, "ValidationException")
}

func TestTestStateContextRequiresMock(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Context:    `{"Execution":{"Name":"x"}}`,
	})
	requireAWSCode(t, err, "ValidationException")
}

func TestTestStateMockStrictValidationMode(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Mock: &TestStateMock{
			Result:              `{not json`,
			ResultProvided:      true,
			FieldValidationMode: "STRICT",
		},
	})
	requireAWSCode(t, err, "ValidationException")

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Mock: &TestStateMock{
			Result:              `{not json`,
			ResultProvided:      true,
			FieldValidationMode: "NONE",
		},
	})
	if err != nil {
		t.Fatalf("NONE mode must skip the mocked result validation: %v", err)
	}
	if resp["status"] != "SUCCEEDED" {
		t.Fatalf("status = %v, want SUCCEEDED", resp["status"])
	}

	// LENIENT is not a member of the MockResponseValidationMode enum
	// (STRICT, PRESENT, NONE) and must be rejected.
	_, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Mock: &TestStateMock{
			Result:              `{}`,
			ResultProvided:      true,
			FieldValidationMode: "LENIENT",
		},
	})
	requireAWSCode(t, err, "ValidationException")
}

func TestTestStateStateConfigurationValidation(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		StateConfig: &TestStateConfiguration{
			ErrorCausedByState: "NoSuchState",
		},
	})
	requireAWSCode(t, err, "ValidationException")

	_, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		StateConfig: &TestStateConfiguration{
			MapIterationFailureCount: -1,
		},
	})
	requireAWSCode(t, err, "ValidationException")
}

// TestTestStateRetrierRetryCount pins the documented retrierRetryCount
// semantics: with one retry already spent on a matching Retry block the
// status is RETRIABLE, errorDetails.retryIndex is the matching block
// index and retryBackoffIntervalSeconds is the next attempt's interval
// (IntervalSeconds x BackoffRate^count).
func TestTestStateRetrierRetryCount(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f",`
	def += `"Retry":[{"ErrorEquals":["Lambda.ServiceException"],"IntervalSeconds":2,"MaxAttempts":3,"BackoffRate":2.0}],`
	def += `"Catch":[{"ErrorEquals":["States.ALL"],"ResultPath":"$.error","Next":"Handle"}],"Next":"Done"},`
	def += `"Handle":{"Type":"Pass","End":true},"Done":{"Type":"Succeed"}}}`

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "T",
		Input:      `{"data":"value"}`,
		Mock: &TestStateMock{
			Error:         "Lambda.ServiceException",
			Cause:         "Service error",
			ErrorProvided: true,
		},
		StateConfig:     &TestStateConfiguration{RetrierRetryCount: 1},
		InspectionLevel: "DEBUG",
	})
	if err != nil {
		t.Fatalf("retriable mock failed: %v", err)
	}
	if resp["status"] != "RETRIABLE" {
		t.Fatalf("status = %v, want RETRIABLE", resp["status"])
	}
	inspection := resp["inspectionData"].(map[string]interface{})
	details := inspection["errorDetails"].(map[string]interface{})
	if details["retryIndex"].(int) != 0 {
		t.Errorf("retryIndex = %v, want 0", details["retryIndex"])
	}
	if details["retryBackoffIntervalSeconds"].(int32) != 4 {
		t.Errorf("retryBackoffIntervalSeconds = %v, want 4 (2 x 2.0^1)", details["retryBackoffIntervalSeconds"])
	}

	// Exhausting the attempts falls through to the Catcher: the worked
	// example shape (status CAUGHT_ERROR, catchIndex, ResultPath output).
	resp, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "T",
		Input:      `{"data":"value"}`,
		Mock: &TestStateMock{
			Error:         "Lambda.ServiceException",
			Cause:         "Service error",
			ErrorProvided: true,
		},
		StateConfig:     &TestStateConfiguration{RetrierRetryCount: 3},
		InspectionLevel: "DEBUG",
	})
	if err != nil {
		t.Fatalf("exhausted mock failed: %v", err)
	}
	if resp["status"] != "CAUGHT_ERROR" {
		t.Fatalf("status = %v, want CAUGHT_ERROR", resp["status"])
	}
	if resp["nextState"] != "Handle" {
		t.Fatalf("nextState = %v, want Handle", resp["nextState"])
	}
	inspection = resp["inspectionData"].(map[string]interface{})
	details = inspection["errorDetails"].(map[string]interface{})
	if details["catchIndex"].(int) != 0 {
		t.Errorf("catchIndex = %v, want 0", details["catchIndex"])
	}
	var caughtOutput map[string]interface{}
	if err := json.Unmarshal([]byte(resp["output"].(string)), &caughtOutput); err != nil {
		t.Fatalf("caught output not JSON: %v", err)
	}
	if caughtOutput["error"] == nil || caughtOutput["data"] == nil {
		t.Errorf("ResultPath must add the error to the input: %v", caughtOutput)
	}
}

// TestTestStateMapIterationFailureCount pins the mocked Map failure
// semantics: within-threshold failures keep the state SUCCEEDED while
// exceeding the threshold fails with
// States.ExceedToleratedFailureThreshold (or the Catch handler).
func TestTestStateMapIterationFailureCount(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemsPath":"$.items",`
	def += `"ItemProcessor":{"ProcessorConfig":{"Mode":"DISTRIBUTED","ExecutionType":"STANDARD"},"StartAt":"W","States":{"W":{"Type":"Pass","End":true}}},`
	def += `"ToleratedFailureCount":1,"End":true}}}`

	// One failed iteration of three stays within the tolerated count.
	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "M",
		Input:      `{"items":[1,2,3]}`,
		Mock: &TestStateMock{
			Result:         "[1,2,3]",
			ResultProvided: true,
		},
		StateConfig:     &TestStateConfiguration{MapIterationFailureCount: 1},
		InspectionLevel: "DEBUG",
	})
	if err != nil {
		t.Fatalf("tolerated failure mock failed: %v", err)
	}
	if resp["status"] != "SUCCEEDED" {
		t.Fatalf("status = %v, want SUCCEEDED", resp["status"])
	}
	inspection := resp["inspectionData"].(map[string]interface{})
	if inspection["afterItemsPath"] == nil {
		t.Errorf("afterItemsPath missing from inspection data: %v", inspection)
	}
	if inspection["toleratedFailureCount"] == nil {
		t.Errorf("toleratedFailureCount missing from inspection data: %v", inspection)
	}

	// Two failed iterations exceed the tolerated count of one.
	resp, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "M",
		Input:      `{"items":[1,2,3]}`,
		Mock: &TestStateMock{
			Result:         "[1,2,3]",
			ResultProvided: true,
		},
		StateConfig:     &TestStateConfiguration{MapIterationFailureCount: 2},
		InspectionLevel: "DEBUG",
	})
	if err != nil {
		t.Fatalf("exceeded failure mock failed: %v", err)
	}
	if resp["status"] != "FAILED" {
		t.Fatalf("status = %v, want FAILED", resp["status"])
	}
	if resp["error"] != "States.ExceedToleratedFailureThreshold" {
		t.Fatalf("error = %v, want States.ExceedToleratedFailureThreshold", resp["error"])
	}

	// The failure count cannot exceed the item count.
	resp, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "M",
		Input:      `{"items":[1,2,3]}`,
		Mock: &TestStateMock{
			Result:         "[1,2,3]",
			ResultProvided: true,
		},
		StateConfig:     &TestStateConfiguration{MapIterationFailureCount: 4},
		InspectionLevel: "DEBUG",
	})
	if err == nil || !strings.Contains(err.Error(), "exceeds the number of items") {
		t.Fatalf("oversized count: err=%v resp=%v, want the ValidationException request error", err, resp)
	}
}

// TestTestStateMockExclusivity pins the documented mutual exclusions: a
// mock carries either a result or an errorOutput, and revealSecrets
// cannot accompany a mock.
func TestTestStateMockExclusivity(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: taskDefForMock,
		StateName:  "T",
		Mock: &TestStateMock{
			Result:         "{}",
			ResultProvided: true,
			Error:          "E",
			Cause:          "C",
			ErrorProvided:  true,
		},
	})
	requireAWSCode(t, err, "ValidationException")

	_, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition:    taskDefForMock,
		StateName:     "T",
		RevealSecrets: true,
		Mock: &TestStateMock{
			Result:         "{}",
			ResultProvided: true,
		},
	})
	requireAWSCode(t, err, "ValidationException")
}

// TestTestStateParallelMockBranchCount pins that a Parallel mocked result
// must be an array with one element per branch.
func TestTestStateParallelMockBranchCount(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"P","States":{"P":{"Type":"Parallel",`
	def += `"Branches":[{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}},{"StartAt":"B","States":{"B":{"Type":"Pass","End":true}}}],`
	def += `"End":true}}}`

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "P",
		Input:      `{}`,
		Mock: &TestStateMock{
			Result:         `["a","b"]`,
			ResultProvided: true,
		},
	})
	if err != nil {
		t.Fatalf("two-branch result failed: %v", err)
	}
	if resp["status"] != "SUCCEEDED" {
		t.Fatalf("status = %v", resp["status"])
	}

	resp, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "P",
		Input:      `{}`,
		Mock: &TestStateMock{
			Result:         `["a"]`,
			ResultProvided: true,
		},
	})
	if err == nil || !strings.Contains(err.Error(), "the Parallel mocked result must have") {
		t.Fatalf("branch mismatch: err=%v resp=%v, want the ValidationException request error", err, resp)
	}
}

// TestTestStateNestedStateName pins that stateName can address a state
// inside a Map ItemProcessor.
func TestTestStateNestedStateName(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"Inner","States":{"Inner":{"Type":"Pass","Result":"deep","End":true}}},"End":true}}}`

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "Inner",
		Input:      `{}`,
	})
	if err != nil {
		t.Fatalf("nested state test failed: %v", err)
	}
	if resp["status"] != "SUCCEEDED" {
		t.Fatalf("status = %v", resp["status"])
	}
	if resp["output"] != `"deep"` {
		t.Fatalf("output = %v, want the quoted JSON string \"deep\"", resp["output"])
	}
}

// TestTestStateLeavesNoHistoryRecords pins the no-resource contract:
// "You can test a state without creating a state machine or updating an
// existing state machine" — a TestState run persists no execution history
// under the synthetic test ARN.
func TestTestStateLeavesNoHistoryRecords(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":{"v":1},"End":true}}}`,
		StateName:  "P",
		Input:      `{}`,
	})
	if err != nil {
		t.Fatalf("testStateCore failed: %v", err)
	}
	if resp["status"] != "SUCCEEDED" {
		t.Fatalf("status = %v, want SUCCEEDED", resp["status"])
	}

	execARN := "arn:aws:states:us-east-1:000000000000:execution:test-state:P"
	events, _, gerr := store.GetExecutionHistory(context.Background(), execARN, 0, "", false)
	if gerr != nil {
		t.Fatalf("history read failed: %v", gerr)
	}
	if len(events) != 0 {
		t.Fatalf("TestState persisted %d orphan history events under %s", len(events), execARN)
	}

	// A Wait state (which logs entered/exited events) leaves nothing too.
	_, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":0,"End":true}}}`,
		StateName:  "W",
		Input:      `{}`,
	})
	if err != nil {
		t.Fatalf("wait testStateCore failed: %v", err)
	}
	events, _, gerr = store.GetExecutionHistory(context.Background(), "arn:aws:states:us-east-1:000000000000:execution:test-state:W", 0, "", false)
	if gerr != nil {
		t.Fatalf("history read failed: %v", gerr)
	}
	if len(events) != 0 {
		t.Fatalf("Wait TestState persisted %d orphan history events", len(events))
	}
}

// TestTestStateMockOnlyCategories pins the isolation contract: "TestState
// only supports the following when a mock is specified: Activity tasks,
// .sync or .waitForTaskToken service integration patterns, Parallel, or
// Map states."
func TestTestStateMockOnlyCategories(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)
	cases := []struct {
		name       string
		definition string
		stateName  string
	}{
		{"Parallel without mock", `{"StartAt":"P","States":{"P":{"Type":"Parallel","Branches":[{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}}],"End":true}}}`, "P"},
		{"Map without mock", `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"W","States":{"W":{"Type":"Pass","End":true}}},"End":true}}}`, "M"},
		{"Activity without mock", `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:us-east-1:000000000000:activity:act","End":true}}}`, "T"},
		{"sync without mock", `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::states:startExecution.sync","End":true}}}`, "T"},
		{"waitForTaskToken without mock", `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::sqs:sendMessage.waitForTaskToken","End":true}}}`, "T"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.testStateCore(context.Background(), store, TestStateInput{
				Definition: tc.definition,
				StateName:  tc.stateName,
				Input:      `{}`,
			})
			if err == nil || !strings.Contains(err.Error(), "when a mock is specified") {
				t.Fatalf("error = %v, want the mock-only rejection", err)
			}
		})
	}

	// The same categories run with a mock.
	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"W","States":{"W":{"Type":"Pass","End":true}}},"End":true}}}`,
		StateName:  "M",
		Input:      `[1]`,
		Mock:       &TestStateMock{Result: `[2]`, ResultProvided: true},
	})
	if err != nil {
		t.Fatalf("mocked Map run failed: %v", err)
	}
}

// TestTestStateRoleArnValidated pins the InvalidArn member: the role ARN
// is validated when supplied.
func TestTestStateRoleArnValidated(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)
	_, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"P","States":{"P":{"Type":"Pass","End":true}}}`,
		StateName:  "P",
		Input:      `{}`,
		RoleArn:    "not-an-arn",
	})
	if err == nil || !strings.Contains(err.Error(), "InvalidArn") {
		t.Fatalf("error = %v, want InvalidArn", err)
	}
	_, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"P","States":{"P":{"Type":"Pass","End":true}}}`,
		StateName:  "P",
		Input:      `{}`,
		RoleArn:    "arn:aws:iam::000000000000:role/test",
	})
	if err != nil {
		t.Fatalf("valid role ARN rejected: %v", err)
	}
}

// TestTestStateFailErrorCauseSplit pins the separate output members: a
// Fail state's error name and cause are distinct members, resolved
// through the real Fail semantics.
func TestTestStateFailErrorCauseSplit(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"CustomError","Cause":"boom"}}}`,
		StateName:  "F",
		Input:      `{}`,
	})
	if err != nil {
		t.Fatalf("fail-state test failed: %v", err)
	}
	if resp["status"] != "FAILED" {
		t.Fatalf("status = %v", resp["status"])
	}
	if resp["error"] != "CustomError" {
		t.Fatalf("error = %v, want CustomError", resp["error"])
	}
	if resp["cause"] != "boom" {
		t.Fatalf("cause = %v, want boom", resp["cause"])
	}

	// A JSONata Fail resolves its expressions like the real executor.
	resp, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"QueryLanguage":"JSONata","StartAt":"F","States":{"F":{"Type":"Fail","Error":"{% $states.input.code %}"}}}`,
		StateName:  "F",
		Input:      `{"code":"E1"}`,
	})
	if err != nil {
		t.Fatalf("jsonata fail-state test failed: %v", err)
	}
	if resp["error"] != "E1" {
		t.Fatalf("error = %v, want the resolved E1", resp["error"])
	}
}

// TestTestStateMockedTailAppliesSelectors pins the mocked output tail:
// ResultSelector and OutputPath run for Map and Parallel like the real
// executor.
func TestTestStateMockedTailAppliesSelectors(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"P","States":{"P":{"Type":"Parallel",`
	def += `"ResultSelector":{"picked":{"x.$":"$.0.x"}},`
	def += `"ResultPath":"$.payload","OutputPath":"$.payload.picked",`
	def += `"Branches":[{"StartAt":"A","States":{"A":{"Type":"Pass","End":true}}},{"StartAt":"B","States":{"B":{"Type":"Pass","End":true}}}],`
	def += `"End":true}}}`

	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		StateName:  "P",
		Input:      `{}`,
		Mock: &TestStateMock{
			Result:         `[{"x":1,"y":2},{"x":3,"y":4}]`,
			ResultProvided: true,
		},
	})
	if err != nil {
		t.Fatalf("mocked parallel tail failed: %v", err)
	}
	if resp["output"] != `{"x":1}` {
		t.Fatalf("output = %v, want the selected and projected {\"x\":1}", resp["output"])
	}
}

// TestTestStateFiveMinuteCap pins the documented run ceiling: "The
// TestState API can run for up to five minutes. If the execution of a
// state exceeds this duration, it fails with the States.Timeout error."
// The cap is shortened here so an hour-long Wait fails fast instead of
// hanging the request.
func TestTestStateFiveMinuteCap(t *testing.T) {
	origTimeout := testStateRunTimeout
	testStateRunTimeout = 150 * time.Millisecond
	t.Cleanup(func() { testStateRunTimeout = origTimeout })

	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":3600,"End":true}}}`

	done := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	go func() {
		resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
			Definition: def,
			StateName:  "W",
			Input:      `{}`,
		})
		if err != nil {
			errCh <- err
			return
		}
		done <- resp
	}()

	select {
	case err := <-errCh:
		t.Fatalf("the cap must fail the run as a result, not an API error: %v", err)
	case resp := <-done:
		if resp["status"] != "FAILED" {
			t.Fatalf("status = %v, want FAILED", resp["status"])
		}
		if resp["error"] != "States.Timeout" {
			t.Fatalf("error = %v, want States.Timeout", resp["error"])
		}
	case <-time.After(5 * time.Second):
		t.Fatal("TestState hung: the five-minute cap never interrupted the Wait state")
	}
}

// TestTestStateMockProcessingFailureSurfaces pins the mock-path processing
// contract: the mock replaces the state's invocation, not its data
// processing, so a processing failure is the execution of the state failing
// and is reported through the error output (status FAILED with error and
// cause) — never as a silently-successful run, a silently-zero tolerance,
// or a caught-error result whose Catch output could not be built.
func TestTestStateMockProcessingFailureSurfaces(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)
	ctx := context.Background()

	// InputPath selecting no value fails the run instead of falling back
	// to the raw input.
	ipDef := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","InputPath":"$.missing","End":true}}}`
	resp, err := svc.testStateCore(ctx, store, TestStateInput{
		Definition: ipDef,
		StateName:  "T",
		Input:      `{"present":1}`,
		Mock:       &TestStateMock{Result: `{"value":42}`, ResultProvided: true},
	})
	if err != nil {
		t.Fatalf("mocked run: %v", err)
	}
	if resp["status"] != "FAILED" {
		t.Fatalf("InputPath failure: status = %v, want FAILED", resp["status"])
	}
	if resp["error"] != "States.Runtime" {
		t.Fatalf("InputPath failure: error = %v, want States.Runtime", resp["error"])
	}
	if cause, _ := resp["cause"].(string); !strings.Contains(cause, "InputPath") {
		t.Fatalf("InputPath failure: cause = %v, want the InputPath failure detail", resp["cause"])
	}

	// A Map threshold path that fails to resolve reports its own failure,
	// not a silently-zero tolerance verdict.
	mapDef := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemsPath":"$.items","ItemProcessor":{"StartAt":"W","States":{"W":{"Type":"Wait","Seconds":0,"End":true}}},"ToleratedFailureCountPath":"$.missing","End":true}}}`
	resp, err = svc.testStateCore(ctx, store, TestStateInput{
		Definition: mapDef,
		StateName:  "M",
		Input:      `{"items":[1,2,3]}`,
		Mock:       &TestStateMock{Result: `[1,2,3]`, ResultProvided: true},
		StateConfig: &TestStateConfiguration{
			MapIterationFailureCount: 1,
		},
	})
	if err != nil {
		t.Fatalf("mocked map run: %v", err)
	}
	if resp["status"] != "FAILED" {
		t.Fatalf("threshold failure: status = %v, want FAILED", resp["status"])
	}
	if resp["error"] != "States.InvalidInput" {
		t.Fatalf("threshold failure: error = %v, want States.InvalidInput", resp["error"])
	}
	if cause, _ := resp["cause"].(string); !strings.Contains(cause, "ToleratedFailureCountPath") {
		t.Fatalf("threshold failure: cause = %v, want the ToleratedFailureCountPath resolution failure", resp["cause"])
	}

	// A Catch whose ResultPath cannot apply fails the tested state; the
	// caught-error result is only reported when the Catch output itself
	// could be built.
	catchDef := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","Catch":[{"ErrorEquals":["States.ALL"],"ResultPath":"$.err","Next":"H"}],"End":true},"H":{"Type":"Succeed"}}}`
	resp, err = svc.testStateCore(ctx, store, TestStateInput{
		Definition: catchDef,
		StateName:  "T",
		Input:      `[1,2]`,
		Mock:       &TestStateMock{Error: "TestFail", Cause: "boom", ErrorProvided: true},
	})
	if err != nil {
		t.Fatalf("mocked error run: %v", err)
	}
	if resp["status"] != "FAILED" {
		t.Fatalf("catch output failure: status = %v, want FAILED", resp["status"])
	}
	if resp["error"] != "States.ResultPathMatchFailure" {
		t.Fatalf("catch output failure: error = %v, want States.ResultPathMatchFailure", resp["error"])
	}
}

// TestTestStateMachineFormValidatesDefinition pins the stateName form's
// definition contract: "If this field is specified, the definition must
// contain a fully-formed state machine definition" — the same structural
// validation the creation paths apply rejects a malformed machine with
// InvalidDefinition before any state executes.
func TestTestStateMachineFormValidatesDefinition(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	rows := []struct {
		name string
		def  string
	}{
		{"dangling Next", `{"StartAt":"T","States":{"T":{"Type":"Pass","Next":"Missing"}}}`},
		{"missing StartAt", `{"States":{"T":{"Type":"Pass","End":true}}}`},
		{"no states", `{"StartAt":"T","States":{}}`},
	}
	for _, row := range rows {
		_, err := svc.testStateCore(context.Background(), store, TestStateInput{
			Definition: row.def,
			StateName:  "T",
		})
		if err == nil {
			t.Errorf("%s: TestState must reject a definition that is not a fully-formed state machine", row.name)
			continue
		}
		requireAWSCode(t, err, "InvalidDefinition")
	}
}

// TestTestStateSingleStateForm pins the documented primary form: "Accepts
// the definition of a single state and executes it" — with no stateName the
// definition itself is the state under test, and the documented CLI example
// runs a Choice whose Next targets exist nowhere, so no machine-level
// structural validation applies to this form.
func TestTestStateSingleStateForm(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	// The documented Choice example: dangling Next targets succeed.
	def := `{"Type":"Choice","Choices":[{"Variable":"$.number","NumericEquals":1,"Next":"Equals 1"},{"Variable":"$.number","NumericEquals":2,"Next":"Equals 2"}],"Default":"No Match"}`
	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		Input:      `{"number":2}`,
	})
	if err != nil {
		t.Fatalf("the single-state form must run without stateName: %v", err)
	}
	if resp["status"] != "SUCCEEDED" || resp["nextState"] != "Equals 2" {
		t.Fatalf("status/nextState = %v/%v, want SUCCEEDED/Equals 2: %v", resp["status"], resp["nextState"], resp)
	}
	if out, _ := resp["output"].(string); !strings.Contains(out, `"number":2`) {
		t.Fatalf("output = %v, want the unfiltered input", resp["output"])
	}

	// A machine definition is not a single state: the form requires the
	// definition to be the state itself.
	_, err = svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: `{"StartAt":"T","States":{"T":{"Type":"Pass","End":true}}}`,
	})
	requireAWSCode(t, err, "InvalidDefinition")
}

// TestTestStateErrorCausedByStateAcceptsNestedState pins the documented
// errorCausedByState shape: the guide's example names "ProcessItem", a
// state inside the Map under test — not a top-level machine state.
func TestTestStateErrorCausedByStateAcceptsNestedState(t *testing.T) {
	svc := &StepFunctionService{}
	store := newTestStateStore(t)

	def := `{"Type":"Map","ItemProcessor":{"StartAt":"ProcessItem","States":{"ProcessItem":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","End":true}}},"End":true}`
	resp, err := svc.testStateCore(context.Background(), store, TestStateInput{
		Definition: def,
		Input:      `[1]`,
		Mock: &TestStateMock{
			Error:         "States.TaskFailed",
			Cause:         "task failed",
			ErrorProvided: true,
		},
		StateConfig: &TestStateConfiguration{ErrorCausedByState: "ProcessItem"},
	})
	if err != nil {
		t.Fatalf("an errorCausedByState nested inside the state under test must be accepted: %v", err)
	}
	if resp["status"] != "FAILED" || resp["error"] != "States.TaskFailed" {
		t.Fatalf("status/error = %v/%v, want FAILED/States.TaskFailed: %v", resp["status"], resp["error"], resp)
	}
}

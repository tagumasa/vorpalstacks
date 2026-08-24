package sfn

import (
	"context"
	"encoding/json"
	"testing"

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
	if err != nil {
		t.Fatalf("oversized failure count call failed: %v", err)
	}
	if resp["status"] != "FAILED" || resp["error"] != "ValidationException" {
		t.Fatalf("oversized count: status=%v error=%v", resp["status"], resp["error"])
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
	if err != nil {
		t.Fatalf("one-branch result call failed: %v", err)
	}
	if resp["status"] != "FAILED" || resp["error"] != "ValidationException" {
		t.Fatalf("branch mismatch: status=%v error=%v", resp["status"], resp["error"])
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

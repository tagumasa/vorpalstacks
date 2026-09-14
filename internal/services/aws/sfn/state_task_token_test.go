package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

const tokenTestActivityARN = "arn:aws:states:us-east-1:000000000000:activity:work"

func newTaskTokenTestHarness(t *testing.T) (*Executor, *sfnstore.StepFunctionStore) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := sfnstore.NewStepFunctionStore(st, "000000000000", "us-east-1")
	if err := store.CreateActivity(context.Background(), &sfnstore.Activity{Name: "work"}); err != nil {
		t.Fatal(err)
	}
	return NewExecutorWithStores(store, nil, "000000000000", "us-east-1", nil), store
}

func newTaskTokenExecCtx(input string) *ExecutionContext {
	eventId := int64(1)
	return &ExecutionContext{
		Execution: &sfnstore.Execution{
			ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:exec-1",
			StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
			Name:            "exec-1",
			Input:           "{}",
			StartDate:       time.Now().UTC(),
		},
		Definition:    &sfnstore.StateMachineDefinition{},
		CurrentState:  "DoWork",
		Input:         input,
		EventId:       &eventId,
		States:        map[string]sfnstore.State{},
		VariableScope: NewVariableScope(nil),
		MapItemIndex:  -1,
	}
}

// The token the worker receives via GetActivityTask must be the exact value
// $$.Task.Token resolved to while Parameters were built, and it must be
// non-empty so the worker can actually report a result.
func TestActivityTaskTokenMatchesParametersValue(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{"orderId":"o-1"}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       tokenTestActivityARN,
		End:            true,
		TimeoutSeconds: float64(10),
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"token.$": "$$.Task.Token",
			"order.$": "$.orderId",
		}},
	}

	type taskResult struct {
		output string
		err    *ExecutionError
	}
	done := make(chan taskResult, 1)
	go func() {
		output, _, execErr := e.executeTask(context.Background(), execCtx, state)
		done <- taskResult{output, execErr}
	}()

	workerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1")
	if err != nil || task == nil {
		t.Fatalf("no activity task became available: task=%v err=%v", task, err)
	}

	if task.TaskToken == "" {
		t.Fatal("activity task persisted with an empty token")
	}
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(task.Input), &input); err != nil {
		t.Fatalf("task input is not JSON: %v", err)
	}
	if got := input["token"]; got != task.TaskToken {
		t.Fatalf("$$.Task.Token mismatch: input carries %v, task token is %v", got, task.TaskToken)
	}
	if got := input["order"]; got != "o-1" {
		t.Fatalf("parameter path resolution broken: order=%v", got)
	}

	if err := store.CompleteActivityTask(task.TaskToken, `{"done":true}`); err != nil {
		t.Fatalf("complete failed: %v", err)
	}

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("executeTask failed: %v", r.err)
		}
		if !strings.Contains(r.output, "done") {
			t.Fatalf("unexpected output: %s", r.output)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeTask did not return after completion")
	}
}

// Every retry attempt must mint its own token: the activity record is keyed
// by token, so attempt 2 overwriting attempt 1's record would let a worker
// holding the stale token complete the retry. After a timeout the stale
// token must be rejected outright.
func TestActivityRetryMintsDistinctTokens(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       tokenTestActivityARN,
		End:            true,
		TimeoutSeconds: float64(1),
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"token.$": "$$.Task.Token",
		}},
		Retry: []*sfnstore.RetryPolicy{{
			ErrorEquals:     []string{"States.Timeout"},
			IntervalSeconds: 0,
			MaxAttempts:     2,
		}},
	}

	done := make(chan *ExecutionError, 1)
	go func() {
		_, _, execErr := e.executeTask(context.Background(), execCtx, state)
		done <- execErr
	}()

	workerCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	first, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1")
	if err != nil || first == nil {
		t.Fatalf("attempt 1 task never arrived: task=%v err=%v", first, err)
	}

	// Wait for the attempt to time out and its record to flip to TIMED_OUT.
	deadline := time.Now().Add(15 * time.Second)
	for {
		rec, recErr := store.GetActivityTaskByToken(first.TaskToken)
		if recErr == nil && rec != nil && rec.Status == "TIMED_OUT" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("attempt 1 task never reached TIMED_OUT")
		}
		time.Sleep(50 * time.Millisecond)
	}

	if err := store.CompleteActivityTask(first.TaskToken, `{"stale":true}`); err == nil {
		t.Fatal("stale token from attempt 1 must be rejected after its timeout")
	}

	second, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-2")
	if err != nil || second == nil {
		t.Fatalf("attempt 2 task never arrived: task=%v err=%v", second, err)
	}
	if second.TaskToken == "" {
		t.Fatal("attempt 2 persisted with an empty token")
	}
	if second.TaskToken == first.TaskToken {
		t.Fatal("retry reused the attempt-1 token")
	}
	var secondInput map[string]interface{}
	if err := json.Unmarshal([]byte(second.Input), &secondInput); err != nil {
		t.Fatalf("attempt 2 input is not JSON: %v", err)
	}
	if got := secondInput["token"]; got != second.TaskToken {
		t.Fatalf("attempt 2 input carries %v but its token is %v", got, second.TaskToken)
	}

	if err := store.CompleteActivityTask(second.TaskToken, `{"ok":true}`); err != nil {
		t.Fatalf("completing attempt 2 failed: %v", err)
	}

	select {
	case execErr := <-done:
		if execErr != nil {
			t.Fatalf("executeTask failed after successful retry: %v", execErr)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("executeTask did not return after retry completion")
	}
}

// Concurrent Map branches share one Executor but must never share or clobber
// each other's task tokens: each branch's record carries its own token and
// the input handed to the worker embeds that same token.
func TestMapBranchesGetDistinctTaskTokens(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{"orders":[{"id":"a"},{"id":"b"}]}`)
	iterator := &sfnstore.StateMachineDefinition{
		StartAt: "DoWork",
		States: map[string]interface{}{
			"DoWork": map[string]interface{}{
				"Type":     "Task",
				"Resource": tokenTestActivityARN,
				"End":      true,
				// A generous task timeout: this deadline only needs to
				// outlive the polling window, and a tight value races
				// with executor goroutine wake-ups when the package runs
				// under load.
				"TimeoutSeconds": 30,
				"Parameters": map[string]interface{}{
					"token.$": "$$.Task.Token",
				},
			},
		},
	}
	state := &sfnstore.MapState{
		Type:           "Map",
		End:            true,
		Iterator:       iterator,
		ItemsPath:      "$.orders",
		MaxConcurrency: 2,
	}

	done := make(chan *ExecutionError, 1)
	go func() {
		_, _, execErr := e.executeMap(context.Background(), execCtx, state)
		done <- execErr
	}()

	// The deadlines below are hang detectors, not latency assertions:
	// generous bounds keep them from racing with executor goroutine
	// wake-ups when the package runs under garbage-collection pressure.
	workerCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	first, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1")
	if err != nil || first == nil {
		t.Fatalf("branch A task never arrived: task=%v err=%v", first, err)
	}
	second, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-2")
	if err != nil || second == nil {
		t.Fatalf("branch B task never arrived: task=%v err=%v", second, err)
	}

	if first.TaskToken == "" || second.TaskToken == "" {
		t.Fatalf("empty task token in Map branch: %q / %q", first.TaskToken, second.TaskToken)
	}
	if first.TaskToken == second.TaskToken {
		t.Fatal("Map branches shared one task token")
	}

	for _, task := range []*sfnstore.ActivityTask{first, second} {
		var input map[string]interface{}
		if err := json.Unmarshal([]byte(task.Input), &input); err != nil {
			t.Fatalf("branch task input is not JSON: %v", err)
		}
		if got := input["token"]; got != task.TaskToken {
			t.Fatalf("branch input carries token %v but its record token is %v", got, task.TaskToken)
		}
	}

	for _, task := range []*sfnstore.ActivityTask{first, second} {
		if err := store.CompleteActivityTask(task.TaskToken, `{"ok":true}`); err != nil {
			t.Fatalf("completing branch task failed: %v", err)
		}
	}

	select {
	case execErr := <-done:
		if execErr != nil {
			t.Fatalf("executeMap failed: %v", execErr)
		}
	case <-time.After(60 * time.Second):
		t.Fatal("executeMap did not return after both branches completed")
	}
}

// A JSONata task's Arguments reference the attempt token through
// $states.context.Task.Token; the worker must receive that exact token just
// like the JSONPath Parameters dialect.
func TestJSONataActivityTaskTokenInArguments(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{"orderId":"o-1"}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       tokenTestActivityARN,
		End:            true,
		TimeoutSeconds: float64(10),
		QueryLanguage:  "JSONata",
		Arguments: map[string]interface{}{
			"token": "{% $states.context.Task.Token %}",
			"order": "{% $states.input.orderId %}",
		},
	}

	type taskResult struct {
		output string
		err    *ExecutionError
	}
	done := make(chan taskResult, 1)
	go func() {
		output, _, execErr := e.executeTask(context.Background(), execCtx, state)
		done <- taskResult{output, execErr}
	}()

	workerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1")
	if err != nil || task == nil {
		t.Fatalf("no activity task became available: task=%v err=%v", task, err)
	}

	var input map[string]interface{}
	if err := json.Unmarshal([]byte(task.Input), &input); err != nil {
		t.Fatalf("task input is not JSON: %v", err)
	}
	if got := input["token"]; got != task.TaskToken {
		t.Fatalf("$states.context.Task.Token mismatch: input carries %v, task token is %v", got, task.TaskToken)
	}
	if got := input["order"]; got != "o-1" {
		t.Fatalf("arguments resolution broken: order=%v", got)
	}

	if err := store.CompleteActivityTask(task.TaskToken, `{"done":true}`); err != nil {
		t.Fatalf("complete failed: %v", err)
	}

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("executeTask failed: %v", r.err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeTask did not return after completion")
	}
}

// The context object exposes $$.Task.Token in ResultSelector; for an
// activity task it must resolve to the token the successful attempt
// actually ran under, not an unrelated value.
func TestActivityResultSelectorResolvesAttemptToken(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       tokenTestActivityARN,
		End:            true,
		TimeoutSeconds: float64(10),
		ResultSelector: &sfnstore.ResultSelector{Fields: map[string]interface{}{
			"tok.$": "$$.Task.Token",
		}},
	}

	type taskResult struct {
		output string
		err    *ExecutionError
	}
	done := make(chan taskResult, 1)
	go func() {
		output, _, execErr := e.executeTask(context.Background(), execCtx, state)
		done <- taskResult{output, execErr}
	}()

	workerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1")
	if err != nil || task == nil {
		t.Fatalf("no activity task became available: task=%v err=%v", task, err)
	}

	var workerResult map[string]interface{}
	if err := json.Unmarshal([]byte(task.Input), &workerResult); err != nil {
		t.Fatalf("task input is not JSON: %v", err)
	}

	if err := store.CompleteActivityTask(task.TaskToken, `{"done":true}`); err != nil {
		t.Fatalf("complete failed: %v", err)
	}

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("executeTask failed: %v", r.err)
		}
		var output map[string]interface{}
		if err := json.Unmarshal([]byte(r.output), &output); err != nil {
			t.Fatalf("output is not JSON: %v (%s)", err, r.output)
		}
		if got := output["tok"]; got != task.TaskToken {
			t.Fatalf("ResultSelector token %v does not match the completed attempt's token %v", got, task.TaskToken)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeTask did not return after completion")
	}
}

// Referencing $$.Task.Token outside a task attempt that carries a token
// must fail the evaluation instead of fabricating a random token no worker
// could ever answer. The failure must classify identically no matter which
// state type evaluated the reference: JSONPath processing failures are
// States.Runtime everywhere, and States.QueryEvaluationError stays reserved
// for JSONata expression failures.
func TestUnbackedTaskTokenClassifiedAsRuntimeEverywhere(t *testing.T) {
	e, _ := newTaskTokenTestHarness(t)

	runtimeTokenErr := func(err *ExecutionError, site string) {
		t.Helper()
		if err == nil {
			t.Fatalf("%s accepted an unbacked $$.Task.Token", site)
		}
		if err.ErrorCode != "States.Runtime" {
			t.Fatalf("%s classified the reference as %s, want States.Runtime (cause: %s)", site, err.ErrorCode, err.Cause)
		}
		if !strings.Contains(err.Cause, "Task.Token") {
			t.Fatalf("%s error cause does not mention Task.Token: %s", site, err.Cause)
		}
	}

	// Parameters of a non-activity task: the evaluation fails before any
	// dispatch, so no integration side effects occur.
	_, _, execErr := e.executeTask(context.Background(), newTaskTokenExecCtx(`{}`), &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:lambda:us-east-1:000000000000:function:not-reached",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"tok.$": "$$.Task.Token",
		}},
	})
	runtimeTokenErr(execErr, "task Parameters")

	// Parameters of a Pass state have no task token.
	_, _, err := e.executePass(context.Background(), newTaskTokenExecCtx(`{}`), &sfnstore.PassState{
		Type: "Pass",
		End:  true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"tok.$": "$$.Task.Token",
		}},
	})
	var passExecErr *ExecutionError
	if !errors.As(err, &passExecErr) {
		t.Fatalf("Pass Parameters returned an unclassified error: %v", err)
	}
	runtimeTokenErr(passExecErr, "Pass Parameters")

	// ResultSelector of a Pass state likewise.
	_, _, err = e.executePass(context.Background(), newTaskTokenExecCtx(`{}`), &sfnstore.PassState{
		Type: "Pass",
		End:  true,
		Result: map[string]interface{}{
			"payload": "done",
		},
		ResultSelector: &sfnstore.ResultSelector{Fields: map[string]interface{}{
			"tok.$": "$$.Task.Token",
		}},
	})
	if !errors.As(err, &passExecErr) {
		t.Fatalf("Pass ResultSelector returned an unclassified error: %v", err)
	}
	runtimeTokenErr(passExecErr, "Pass ResultSelector")

	// ItemSelector of a JSONPath Map state evaluates outside any task.
	_, _, mapErr := e.executeMap(context.Background(), newTaskTokenExecCtx(`{"orders":[{"id":"a"}]}`), &sfnstore.MapState{
		Type:           "Map",
		End:            true,
		ItemsPath:      "$.orders",
		MaxConcurrency: 2,
		ItemSelector: map[string]interface{}{
			"tok.$": "$$.Task.Token",
		},
		Iterator: &sfnstore.StateMachineDefinition{
			StartAt: "Noop",
			States: map[string]interface{}{
				"Noop": map[string]interface{}{"Type": "Pass", "End": true},
			},
		},
	})
	runtimeTokenErr(mapErr, "Map ItemSelector")

	// Direct context resolution mirrors the paths above.
	if _, ctxErr := e.getContextValue(nil, "", "$$.Task.Token"); ctxErr == nil {
		t.Fatal("getContextValue fabricated a token for an empty attempt token")
	}
}

// Both query dialects fail a tokenless token reference, by the same
// contract: JSONPath fails with States.Runtime (see
// TestUnbackedTaskTokenClassifiedAsRuntimeEverywhere) and JSONata fails
// with States.QueryEvaluationError, because the context object carries no
// Task section outside a token-backed attempt and JSON cannot represent
// an undefined value expression.
func TestJSONataUnbackedTaskTokenFails(t *testing.T) {
	e, _ := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{"orderId":"o-1"}`)

	var inputData interface{}
	if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
		t.Fatal(err)
	}

	// Without a token the context object carries no Task section at all.
	if _, hasTask := e.buildContextObject(execCtx)["Task"]; hasTask {
		t.Fatal("context object exposed a Task section without a token")
	}

	arguments := map[string]interface{}{
		"token": "{% $states.context.Task.Token %}",
		"order": "{% $states.input.orderId %}",
	}

	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
	_, err := e.applyJSONataArguments(context.Background(), arguments, statesVar, execCtx.VariableScope)
	if err == nil {
		t.Fatal("expected the tokenless token reference to fail the expression")
	}
	if !strings.Contains(err.Error(), "undefined") {
		t.Fatalf("tokenless token failure does not report undefined: %v", err)
	}

	// Control: with a token minted the same expression resolves to it.
	execCtx.TaskToken = "tok-1"
	if _, hasTask := e.buildContextObject(execCtx)["Task"]; !hasTask {
		t.Fatal("context object omitted the Task section for a token-backed attempt")
	}
	statesVar = e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
	var resolved map[string]interface{}
	out, err := e.applyJSONataArguments(context.Background(), arguments, statesVar, execCtx.VariableScope)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(out), &resolved); err != nil {
		t.Fatal(err)
	}
	if got := resolved["token"]; got != "tok-1" {
		t.Fatalf("token argument resolved to %v, want tok-1", got)
	}
}

// TestActivityWaitTimeoutClassifiesTyped pins the typed timeout signal: an
// activity task whose report never arrives fails with States.Timeout
// through the waiter layer's timedOut flag, not through error-text
// matching.
func TestActivityWaitTimeoutClassifiesTyped(t *testing.T) {
	e, _ := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       tokenTestActivityARN,
		End:            true,
		TimeoutSeconds: float64(1),
	}

	_, _, execErr := e.executeTask(context.Background(), execCtx, state)
	if execErr == nil {
		t.Fatal("the unreported activity task must fail the state")
	}
	if execErr.ErrorCode != "States.Timeout" {
		t.Fatalf("timeout classification = %q, want States.Timeout", execErr.ErrorCode)
	}
}

// TestSendTaskFailureIdentityNotReclassifiedByCauseText pins that a worker
// error whose NAME or CAUSE merely quotes "States.Timeout" keeps its own
// identity: the timeout classification consumes typed signals only, so
// Catch on the real error name fires.
func TestSendTaskFailureIdentityNotReclassifiedByCauseText(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       tokenTestActivityARN,
		End:            true,
		TimeoutSeconds: float64(10),
	}

	type taskResult struct {
		err *ExecutionError
	}
	done := make(chan taskResult, 1)
	go func() {
		_, _, execErr := e.executeTask(context.Background(), execCtx, state)
		done <- taskResult{execErr}
	}()

	workerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, err := store.GetActivityTask(workerCtx, tokenTestActivityARN, "worker-1")
	if err != nil || task == nil {
		t.Fatalf("no activity task became available: task=%v err=%v", task, err)
	}
	if err := store.FailActivityTask(task.TaskToken, "MyStates.TimeoutError", `upstream reported "States.Timeout"`); err != nil {
		t.Fatalf("fail failed: %v", err)
	}

	select {
	case r := <-done:
		if r.err == nil {
			t.Fatal("the failed report must fail the state")
		}
		if r.err.ErrorCode != "MyStates.TimeoutError" {
			t.Fatalf("error identity = %q, want MyStates.TimeoutError (quoted timeout text must not reclassify)", r.err.ErrorCode)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("executeTask did not return after the failure report")
	}
}

// TestTaskCredentialsTemplateResolves pins the Credentials payload
// template: "you can also specify a JSONPath value or an intrinsic
// function that resolves to an IAM role ARN at runtime based on the
// execution input" — both dialects resolve, and the resolved ARN rides
// the scheduled event's taskCredentials member. The per-task assume-role
// step itself is the recorded platform exclusion.
func TestTaskCredentialsTemplateResolves(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutor(store, nil)
	e.region = "us-east-1"

	execCtx := &ExecutionContext{
		Execution: &sfnstore.Execution{
			ExecutionArn:    "arn:aws:states:us-east-1:000000000000:execution:sm:cred",
			StateMachineArn: "arn:aws:states:us-east-1:000000000000:stateMachine:sm",
		},
		Input:         `{"acct":"123","role":"arn:aws:iam::9:role/X"}`,
		QueryLanguage: "JSONPath",
	}

	state := &sfnstore.TaskState{
		Resource: "arn:aws:states:::sns:publish",
		Credentials: map[string]interface{}{
			"RoleArn.$": "States.Format('arn:aws:iam::{}:role/R', $.acct)",
		},
	}
	creds, cerr := e.evaluateTaskCredentials(context.Background(), execCtx, state, execCtx.Input)
	if cerr != nil {
		t.Fatalf("JSONPath credentials failed: %v", cerr.Cause)
	}
	if creds.RoleArn != "arn:aws:iam::123:role/R" {
		t.Errorf("JSONPath roleArn = %q, want the formatted ARN", creds.RoleArn)
	}

	state.QueryLanguage = "JSONata"
	state.Credentials = map[string]interface{}{"RoleArn": "{% $states.input.role %}"}
	creds, cerr = e.evaluateTaskCredentials(context.Background(), execCtx, state, execCtx.Input)
	if cerr != nil {
		t.Fatalf("JSONata credentials failed: %v", cerr.Cause)
	}
	if creds.RoleArn != "arn:aws:iam::9:role/X" {
		t.Errorf("JSONata roleArn = %q, want the input role", creds.RoleArn)
	}

	state.Credentials = map[string]interface{}{"Other.$": "$.acct"}
	if _, cerr = e.evaluateTaskCredentials(context.Background(), execCtx, state, execCtx.Input); cerr == nil {
		t.Error("a template without RoleArn must fail the evaluation")
	}

	// The scheduled event serialises the modelled taskCredentials member.
	evt := &sfnstore.ExecutionHistoryEvent{
		Type: "TaskScheduled",
		TaskScheduledEventDetails: &sfnstore.TaskScheduledEventDetails{
			Resource:        state.Resource,
			ResourceType:    "sns",
			TaskCredentials: &sfnstore.TaskScheduledCredentials{RoleArn: "arn:aws:iam::123:role/R"},
		},
	}
	details, ok := historyEventToResponse(evt, true)["taskScheduledEventDetails"].(map[string]interface{})
	if !ok {
		t.Fatal("taskScheduledEventDetails did not serialise")
	}
	tc, _ := details["taskCredentials"].(map[string]interface{})
	if tc["roleArn"] != "arn:aws:iam::123:role/R" {
		t.Errorf("taskCredentials = %v, want the roleArn member", details["taskCredentials"])
	}
}

// TestMachineTimeoutMidTaskRecordsNoTaskTimeout pins the deadline-source
// discrimination: a machine-level TimeoutSeconds expiring while a task
// without its own timeout is in flight terminates the execution
// TIMED_OUT without minting the task's timeout event and without a
// Catcher consuming the interruption.
func TestMachineTimeoutMidTaskRecordsNoTaskTimeout(t *testing.T) {
	e, store := newTaskTokenTestHarness(t)
	definition := `{"StartAt":"T","TimeoutSeconds":1,"States":{
		"T":{"Type":"Task","Resource":"` + tokenTestActivityARN + `","Catch":[{"ErrorEquals":["States.ALL"],"Next":"Fallback"}]},
		"Fallback":{"Type":"Pass","Result":"recovered","End":true}}}`
	sm := &sfnstore.StateMachine{Name: "mtimeout", Definition: definition}
	if err := store.CreateStateMachine(context.Background(), sm); err != nil {
		t.Fatalf("create machine: %v", err)
	}
	exec := sfnstore.NewExecution(sm.StateMachineArn, "mt-1", `{}`, "")
	exec.ExecutionArn = sm.StateMachineArn + ":mt-1"
	if err := store.CreateExecution(context.Background(), exec); err != nil {
		t.Fatalf("create execution: %v", err)
	}

	// Direct-drive discrimination: with the machine context dead,
	// executeTask must return the interruption vehicle, not consume the
	// deadline through its Catch handler.
	{
		directCtx, directCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer directCancel()
		directExec := newTaskTokenExecCtx(`{}`)
		directState := &sfnstore.TaskState{
			Type:     "Task",
			Resource: tokenTestActivityARN,
			End:      true,
			Catch:    []*sfnstore.CatchPolicy{{ErrorEquals: []string{"States.ALL"}, Next: "Fallback"}},
		}
		time.Sleep(150 * time.Millisecond)
		_, next, directErr := e.executeTask(directCtx, directExec, directState)
		if directErr == nil {
			t.Error("executeTask consumed the machine deadline through its Catch handler — the terminal path owns a dead context")
		} else if directErr.ErrorCode != "States.Timeout" || next != "" {
			t.Errorf("executeTask returned (%s, %q), want the States.Timeout interruption vehicle with no transition", directErr.ErrorCode, next)
		}
	}

	// The machine deadline propagates as the interruption vehicle; the
	// caller (the async launcher) logs it, the record is what matters.
	_ = e.ExecuteStateMachine(context.Background(), exec)
	if exec.Status != "TIMED_OUT" {
		t.Fatalf("status = %s, want TIMED_OUT", exec.Status)
	}

	history, _, herr := store.GetExecutionHistory(context.Background(), exec.ExecutionArn, 200, "", false)
	if herr != nil {
		t.Fatalf("history: %v", herr)
	}
	sawExecutionTimedOut, sawTaskTimeout, sawFallback := false, false, false
	for _, ev := range history {
		switch ev.Type {
		case "ExecutionTimedOut":
			sawExecutionTimedOut = true
		case "TaskTimedOut", "ActivityTimedOut", "LambdaFunctionTimedOut":
			sawTaskTimeout = true
		}
		if ev.StateEnteredEventDetails != nil && ev.StateEnteredEventDetails.Name == "Fallback" {
			sawFallback = true
		}
	}
	if !sawExecutionTimedOut {
		t.Error("history lacks ExecutionTimedOut — the machine deadline must terminate the execution")
	}
	if sawTaskTimeout {
		t.Error("history records a task timeout event for a task that configured no timeout — the machine deadline is not the task's own")
	}
	if sawFallback {
		t.Error("the States.ALL Catcher consumed the machine deadline — no recovery transition may run")
	}
}

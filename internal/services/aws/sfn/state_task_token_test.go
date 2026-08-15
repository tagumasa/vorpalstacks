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
	return NewExecutorWithStores(store, nil, "000000000000", "us-east-1"), store
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
				"Type":           "Task",
				"Resource":       tokenTestActivityARN,
				"End":            true,
				"TimeoutSeconds": 10,
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

	workerCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
	case <-time.After(20 * time.Second):
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
// could ever answer.
func TestTaskTokenReferenceWithoutTokenFails(t *testing.T) {
	e, _ := newTaskTokenTestHarness(t)

	// Parameters of a Pass state have no task token.
	_, _, err := e.executePass(context.Background(), newTaskTokenExecCtx(`{}`), &sfnstore.PassState{
		Type: "Pass",
		End:  true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"tok.$": "$$.Task.Token",
		}},
	})
	if err == nil || !strings.Contains(err.Error(), "Task.Token") {
		t.Fatalf("Pass Parameters accepted an unbacked $$.Task.Token: %v", err)
	}

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
	if err == nil || !strings.Contains(err.Error(), "Task.Token") {
		t.Fatalf("Pass ResultSelector accepted an unbacked $$.Task.Token: %v", err)
	}

	// Direct context resolution mirrors the two paths above.
	if _, ctxErr := e.getContextValue("", "$$.Task.Token"); ctxErr == nil {
		t.Fatal("getContextValue fabricated a token for an empty attempt token")
	}
}

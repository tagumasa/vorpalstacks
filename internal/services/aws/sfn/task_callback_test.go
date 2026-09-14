package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// callbackSQSInvoker records the messages the callback pattern delivers so
// tests can extract the task token exactly the way a queue consumer would.
type callbackSQSInvoker struct {
	mu      sync.Mutex
	bodies  []string
	failMsg bool
}

func (f *callbackSQSInvoker) GetQueueByName(ctx context.Context, region, queueName string) (string, error) {
	return "https://sqs." + region + ".amazonaws.com/000000000000/" + queueName, nil
}

func (f *callbackSQSInvoker) GetQueueARN(ctx context.Context, region, queueURL string) (string, error) {
	return "", nil
}

func (f *callbackSQSInvoker) SendMessage(ctx context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failMsg {
		return "", "", fmt.Errorf("queue unavailable")
	}
	f.bodies = append(f.bodies, body)
	return "msg-id-1", "0123456789abcdef0123456789abcdef", nil
}

func (f *callbackSQSInvoker) ReceiveMessage(ctx context.Context, region, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, nil
}

func (f *callbackSQSInvoker) DeleteMessage(ctx context.Context, region, queueURL, receiptHandle string) error {
	return nil
}

func (f *callbackSQSInvoker) sentBodies() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.bodies...)
}

// callbackLambdaInvoker records the optimised integration's invoke and
// answers with a canned status and payload.
type callbackLambdaInvoker struct {
	mu         sync.Mutex
	names      []string
	payloads   [][]byte
	statusCode int64
	response   []byte
	invokeErr  error
}

func (f *callbackLambdaInvoker) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.names = append(f.names, functionName)
	f.payloads = append(f.payloads, payload)
	return f.statusCode, f.response, f.invokeErr
}

func (f *callbackLambdaInvoker) InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	return invokers.LambdaInvocation{}, nil
}

func (f *callbackLambdaInvoker) GetFunctionARN(ctx context.Context, functionName string) (string, error) {
	return functionName, nil
}

// panickingLambdaInvoker records the submit payload and then dies, standing
// in for a worker crash between registration and the wait.
type panickingLambdaInvoker struct {
	mu      sync.Mutex
	payload []byte
}

func (p *panickingLambdaInvoker) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	p.mu.Lock()
	p.payload = payload
	p.mu.Unlock()
	panic("simulated submit crash")
}

func (p *panickingLambdaInvoker) InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	return invokers.LambdaInvocation{}, nil
}

func (p *panickingLambdaInvoker) GetFunctionARN(ctx context.Context, functionName string) (string, error) {
	return functionName, nil
}

func (p *panickingLambdaInvoker) recordedPayload() []byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.payload
}

// awaitCondition polls cond until it holds or the deadline passes.
func awaitCondition(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// taskScheduledToken extracts the task token the payload carried from the
// TaskScheduled event's parameters — the value $$.Task.Token resolved to.
func taskScheduledToken(t *testing.T, store *sfnstore.StepFunctionStore, execArn string) string {
	t.Helper()
	events, _, err := store.GetExecutionHistory(context.Background(), execArn, 100, "", false)
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range events {
		if ev.Type != "TaskScheduled" || ev.TaskScheduledEventDetails == nil {
			continue
		}
		raw, ok := ev.TaskScheduledEventDetails.Parameters.(string)
		if !ok {
			continue
		}
		var params struct {
			MessageBody struct {
				TaskToken string `json:"TaskToken"`
			} `json:"MessageBody"`
		}
		if err := json.Unmarshal([]byte(raw), &params); err == nil && params.MessageBody.TaskToken != "" {
			return params.MessageBody.TaskToken
		}
	}
	t.Fatal("no TaskScheduled event carried a task token in MessageBody")
	return ""
}

func historyTypes(store *sfnstore.StepFunctionStore, execArn string) []string {
	events, _, err := store.GetExecutionHistory(context.Background(), execArn, 200, "", false)
	if err != nil {
		return nil
	}
	var out []string
	for _, ev := range events {
		out = append(out, ev.Type)
	}
	return out
}

// TestCallbackPatternDeliversTokenAndResumes pins the whole callback
// round-trip: the delivered message carries the exact token the attempt
// registered, SendTaskSuccess resumes the suspended task with its output
// as the task result, and the history records the modelled TaskSubmitted
// event with the suffix-form resource and the submit response.
func TestCallbackPatternDeliversTokenAndResumes(t *testing.T) {
	store := newMapTestStore(t)
	sqs := &callbackSQSInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	e := NewExecutorWithStores(store, bus, "000000000000", "us-east-1", nil)
	execCtx := newTaskTokenExecCtx(`{"orderId":"o-1"}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::sqs:sendMessage.waitForTaskToken",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"QueueUrl": "https://sqs.us-east-1.amazonaws.com/000000000000/callback-queue",
			"MessageBody": map[string]interface{}{
				"orderId.$":   "$.orderId",
				"TaskToken.$": "$$.Task.Token",
			},
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

	awaitCondition(t, "the SQS message to be delivered", func() bool {
		return len(sqs.sentBodies()) > 0
	})

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(sqs.sentBodies()[0]), &body); err != nil {
		t.Fatalf("delivered body is not JSON: %v (%s)", err, sqs.sentBodies()[0])
	}
	if got := body["orderId"]; got != "o-1" {
		t.Fatalf("delivered body lost the state input: orderId=%v", got)
	}
	token, _ := body["TaskToken"].(string)
	if token == "" {
		t.Fatal("delivered body carries no task token")
	}

	if err := store.CompleteActivityTask(token, `{"approved":true}`); err != nil {
		t.Fatalf("SendTaskSuccess failed: %v", err)
	}

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("callback task failed: %v", r.err)
		}
		if strings.TrimSpace(r.output) != `{"approved":true}` {
			t.Fatalf("task output must be the SendTaskSuccess output, got %s", r.output)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("callback task did not resume after SendTaskSuccess")
	}

	types := historyTypes(store, execCtx.Execution.ExecutionArn)
	hasSubmitted := false
	for _, typ := range types {
		if typ == "TaskSubmitted" {
			hasSubmitted = true
		}
	}
	if !hasSubmitted {
		t.Fatalf("history lacks TaskSubmitted: %v", types)
	}

	events, _, err := store.GetExecutionHistory(context.Background(), execCtx.Execution.ExecutionArn, 200, "", false)
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range events {
		if ev.Type != "TaskSubmitted" || ev.TaskSubmittedEventDetails == nil {
			continue
		}
		d := ev.TaskSubmittedEventDetails
		if d.Resource != "arn:aws:states:::sqs:sendMessage.waitForTaskToken" {
			t.Fatalf("TaskSubmitted resource must keep the pattern suffix, got %s", d.Resource)
		}
		if !strings.Contains(d.Output, "MessageId") {
			t.Fatalf("TaskSubmitted output must carry the submit response, got %s", d.Output)
		}
	}

	// The settled record stays: a duplicate report keeps the closed-task
	// answer class (TaskTimedOut — "the task token has either expired or
	// the task associated with the token has already been closed"), not
	// TaskDoesNotExist.
	if err := store.CompleteActivityTask(token, "{}"); !errors.Is(err, sfnstore.ErrTaskNotRunning) {
		t.Fatalf("duplicate report = %v, want ErrTaskNotRunning", err)
	}
}

// TestCallbackPanicCleansRegisteredTask pins the registration's panic
// safety: a crash between registration and the wait must not strand a
// RUNNING record — the deferred cleanup removes it while the frame
// unwinds (the execution runner recovers above this point in production).
func TestCallbackPanicCleansRegisteredTask(t *testing.T) {
	store := newMapTestStore(t)
	lam := &panickingLambdaInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetLambdaInvoker(lam)
	e := NewExecutorWithStores(store, bus, "000000000000", "us-east-1", nil)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::lambda:invoke.waitForTaskToken",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"FunctionName": "arn:aws:lambda:us-east-1:000000000000:function:approver",
			"Payload":      map[string]interface{}{"TaskToken.$": "$$.Task.Token"},
		}},
	}

	recovered := func() (r interface{}) {
		defer func() { r = recover() }()
		_, _, _ = e.executeTask(context.Background(), execCtx, state)
		return nil
	}()
	if recovered == nil {
		t.Fatal("the submit panic never propagated — the cleanup path is unexercised")
	}

	// invokeLambdaIntegration marshals the resolved Payload member, so the
	// recorded bytes are the token-carrying object itself.
	var params struct {
		TaskToken string `json:"TaskToken"`
	}
	if err := json.Unmarshal(lam.recordedPayload(), &params); err != nil {
		t.Fatalf("submit payload is not JSON: %v", err)
	}
	if params.TaskToken == "" {
		t.Fatal("submit payload carried no task token")
	}
	if _, err := store.GetActivityTaskByToken(params.TaskToken); !errors.Is(err, sfnstore.ErrTaskNotFound) {
		t.Fatalf("registered record survived the panic: %v", err)
	}
}

// TestCallbackPatternFailureKeepsErrorIdentity: the worker's
// SendTaskFailure error name is the task's error identity, reachable by
// Catch, not the generic wildcard.
func TestCallbackPatternFailureKeepsErrorIdentity(t *testing.T) {
	store := newMapTestStore(t)
	sqs := &callbackSQSInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	e := NewExecutorWithStores(store, bus, "000000000000", "us-east-1", nil)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::sqs:sendMessage.waitForTaskToken",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"QueueUrl":    "https://sqs.us-east-1.amazonaws.com/000000000000/q",
			"MessageBody": map[string]interface{}{"TaskToken.$": "$$.Task.Token"},
		}},
	}

	done := make(chan *ExecutionError, 1)
	go func() {
		_, _, execErr := e.executeTask(context.Background(), execCtx, state)
		done <- execErr
	}()

	awaitCondition(t, "the SQS message to be delivered", func() bool {
		return len(sqs.sentBodies()) > 0
	})
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(sqs.sentBodies()[0]), &body); err != nil {
		t.Fatal(err)
	}
	if err := store.FailActivityTask(body["TaskToken"].(string), "Approver.Rejected", "limit exceeded"); err != nil {
		t.Fatal(err)
	}

	select {
	case execErr := <-done:
		if execErr == nil || execErr.ErrorCode != "Approver.Rejected" {
			t.Fatalf("error identity lost: %v", execErr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("callback task did not return after SendTaskFailure")
	}
}

// TestCallbackPatternSubmitFailureRecordsTaskSubmitFailed: when the
// submission itself fails, the history carries the modelled TaskSubmitFailed
// event and the registered token is spent — a later report on it cannot
// resolve.
func TestCallbackPatternSubmitFailureRecordsTaskSubmitFailed(t *testing.T) {
	store := newMapTestStore(t)
	sqs := &callbackSQSInvoker{failMsg: true}
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	e := NewExecutorWithStores(store, bus, "000000000000", "us-east-1", nil)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::sns:publish.waitForTaskToken",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"TopicArn":    "arn:aws:sns:us-east-1:000000000000:topic",
			"MessageBody": map[string]interface{}{"TaskToken.$": "$$.Task.Token"},
		}},
	}

	// The SNS invoker is absent, so the publish fails before any delivery:
	// the submit phase fails and TaskSubmitFailed must be recorded.
	_, _, execErr := e.executeTask(context.Background(), execCtx, state)
	if execErr == nil {
		t.Fatal("submit failure must fail the task")
	}

	types := historyTypes(store, execCtx.Execution.ExecutionArn)
	submitted, submitFailed := false, false
	for _, typ := range types {
		switch typ {
		case "TaskSubmitted":
			submitted = true
		case "TaskSubmitFailed":
			submitFailed = true
		}
	}
	if submitted {
		t.Fatalf("a failed submit must not record TaskSubmitted: %v", types)
	}
	if !submitFailed {
		t.Fatalf("history lacks TaskSubmitFailed: %v", types)
	}
}

// TestCallbackPatternTimeoutSpendsToken: with no callback inside
// TimeoutSeconds the task fails with States.Timeout and the token is spent
// — a late SendTaskSuccess is rejected.
func TestCallbackPatternTimeoutSpendsToken(t *testing.T) {
	store := newMapTestStore(t)
	sqs := &callbackSQSInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	e := NewExecutorWithStores(store, bus, "000000000000", "us-east-1", nil)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:           "Task",
		Resource:       "arn:aws:states:::sqs:sendMessage.waitForTaskToken",
		End:            true,
		TimeoutSeconds: float64(1),
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"QueueUrl":    "https://sqs.us-east-1.amazonaws.com/000000000000/q",
			"MessageBody": map[string]interface{}{"TaskToken.$": "$$.Task.Token"},
		}},
	}

	_, _, execErr := e.executeTask(context.Background(), execCtx, state)
	if execErr == nil || execErr.ErrorCode != "States.Timeout" {
		t.Fatalf("a callback task with no report must time out, got %v", execErr)
	}

	token := taskScheduledToken(t, store, execCtx.Execution.ExecutionArn)
	if err := store.CompleteActivityTask(token, `{}`); err == nil {
		t.Fatal("a timed-out token must not accept a late SendTaskSuccess")
	}

	types := historyTypes(store, execCtx.Execution.ExecutionArn)
	found := false
	for _, typ := range types {
		if typ == "TaskTimedOut" {
			found = true
		}
	}
	if !found {
		t.Fatalf("history lacks TaskTimedOut: %v", types)
	}
}

// TestOptimisedLambdaInvokeWrapperAndError pins the optimised
// lambda:invoke contract: FunctionName and Payload are API parameters, the
// task result is the documented invoke response wrapper, and a function
// error keeps its own error name.
func TestOptimisedLambdaInvokeWrapperAndError(t *testing.T) {
	lam := &callbackLambdaInvoker{statusCode: 200, response: []byte(`{"ok":1}`)}
	bus := eventbus.NewEventBus()
	bus.SetLambdaInvoker(lam)
	e := NewExecutorWithStores(newMapTestStore(t), bus, "000000000000", "us-east-1", nil)

	input := `{"FunctionName":"arn:aws:lambda:us-east-1:000000000000:function:approver","Payload":{"a":1}}`
	out, err := e.invokeLambdaIntegration(context.Background(), input)
	if err != nil {
		t.Fatalf("optimised invoke failed: %v", err)
	}
	var wrapper map[string]interface{}
	if jerr := json.Unmarshal([]byte(out), &wrapper); jerr != nil {
		t.Fatalf("invoke response is not JSON: %v (%s)", jerr, out)
	}
	if got, _ := wrapper["Payload"].(map[string]interface{}); got["ok"] != float64(1) {
		t.Fatalf("wrapper Payload must carry the function result, got %v", wrapper["Payload"])
	}
	if wrapper["StatusCode"] != float64(200) || wrapper["ExecutedVersion"] != "$LATEST" {
		t.Fatalf("wrapper members wrong: %v", wrapper)
	}
	if len(lam.names) != 1 || lam.names[0] != "arn:aws:lambda:us-east-1:000000000000:function:approver" {
		t.Fatalf("FunctionName not passed through: %v", lam.names)
	}
	if string(lam.payloads[0]) != `{"a":1}` {
		t.Fatalf("Payload not passed through: %s", lam.payloads[0])
	}

	lam.statusCode, lam.response = 500, []byte(`{"errorType":"Custom.Error","errorMessage":"boom"}`)
	_, err = e.invokeLambdaIntegration(context.Background(), input)
	var tf *taskFailure
	if !errors.As(err, &tf) || tf.code != "Custom.Error" {
		t.Fatalf("function error must keep its error name, got %v", err)
	}
}

const startExecutionTestTargetARN = "arn:aws:states:us-east-1:000000000000:stateMachine:target-sm"

func newStartExecutionHarness(t *testing.T) (*Executor, *sfnstore.StepFunctionStore) {
	t.Helper()
	store := newMapTestStore(t)
	if err := store.CreateStateMachine(context.Background(), &sfnstore.StateMachine{
		StateMachineArn: startExecutionTestTargetARN,
		Name:            "target-sm",
		Type:            "STANDARD",
		Definition:      `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":{"answer":42},"End":true}}}`,
	}); err != nil {
		t.Fatal(err)
	}
	return NewExecutorWithStores(store, eventbus.NewEventBus(), "000000000000", "us-east-1", nil), store
}

func awaitChildTerminal(t *testing.T, store *sfnstore.StepFunctionStore, childArn string) *sfnstore.Execution {
	t.Helper()
	var last *sfnstore.Execution
	awaitCondition(t, "the child execution to finish", func() bool {
		exec, err := store.GetExecution(context.Background(), childArn)
		if err != nil || exec == nil {
			return false
		}
		last = exec
		switch exec.Status {
		case "SUCCEEDED", "FAILED", "TIMED_OUT", "ABORTED":
			return true
		}
		return false
	})
	return last
}

// TestStartExecutionIntegrationPlain: the plain form starts the child and
// returns the StartExecution response pair without waiting.
func TestStartExecutionIntegrationPlain(t *testing.T) {
	e, store := newStartExecutionHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::states:startExecution",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"StateMachineArn": startExecutionTestTargetARN,
			"Name":            "child-plain",
			"Input":           map[string]interface{}{"x": 1},
		}},
	}

	output, _, execErr := e.executeTask(context.Background(), execCtx, state)
	if execErr != nil {
		t.Fatalf("plain startExecution failed: %v", execErr)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		t.Fatalf("response is not JSON: %v (%s)", err, output)
	}
	childArn, _ := resp["ExecutionArn"].(string)
	if !strings.HasSuffix(childArn, ":target-sm:child-plain") {
		t.Fatalf("child ARN wrong: %s", childArn)
	}

	child := awaitChildTerminal(t, store, childArn)
	if child.Status != "SUCCEEDED" {
		t.Fatalf("child status %s", child.Status)
	}
	if child.Input != `{"x":1}` {
		t.Fatalf("child input not the parameter Input: %s", child.Input)
	}
}

// TestStartExecutionIntegrationSyncOutputs: .sync returns the child's
// output as a string and .sync:2 as parsed JSON, both after the child
// finished.
func TestStartExecutionIntegrationSyncOutputs(t *testing.T) {
	e, _ := newStartExecutionHarness(t)

	syncState := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::states:startExecution.sync",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"StateMachineArn": startExecutionTestTargetARN,
			"Name":            "child-sync",
		}},
	}
	output, _, execErr := e.executeTask(context.Background(), newTaskTokenExecCtx(`{}`), syncState)
	if execErr != nil {
		t.Fatalf(".sync failed: %v", execErr)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		t.Fatalf(".sync response is not JSON: %v (%s)", err, output)
	}
	if resp["Output"] != `{"answer":42}` {
		t.Fatalf(".sync Output must be the child output string, got %v", resp["Output"])
	}
	if resp["Status"] != "SUCCEEDED" {
		t.Fatalf(".sync Status wrong: %v", resp["Status"])
	}

	sync2State := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::states:startExecution.sync:2",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"StateMachineArn": startExecutionTestTargetARN,
			"Name":            "child-sync2",
		}},
	}
	output, _, execErr = e.executeTask(context.Background(), newTaskTokenExecCtx(`{}`), sync2State)
	if execErr != nil {
		t.Fatalf(".sync:2 failed: %v", execErr)
	}
	resp = nil
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		t.Fatalf(".sync:2 response is not JSON: %v (%s)", err, output)
	}
	parsed, ok := resp["Output"].(map[string]interface{})
	if !ok || parsed["answer"] != float64(42) {
		t.Fatalf(".sync:2 Output must be parsed JSON, got %v", resp["Output"])
	}
}

// TestChildSyncTimeoutSurfacesAsTaskFailed pins the .sync error mapping:
// "If a nested state machine throws a States.Timeout, the parent will
// receive a States.TaskFailed error" — the timeout identity alone is
// remapped, its cause passes through untouched, and every other child
// error keeps its own identity.
func TestChildSyncTimeoutSurfacesAsTaskFailed(t *testing.T) {
	e, _ := newStartExecutionHarness(t)

	cases := []struct {
		name      string
		child     *sfnstore.Execution
		wantCode  string
		wantCause string
	}{
		{
			"timed out child",
			&sfnstore.Execution{Status: "TIMED_OUT", Error: "States.Timeout", Cause: "exceeded 60 seconds"},
			"States.TaskFailed",
			"exceeded 60 seconds",
		},
		{
			"custom child error",
			&sfnstore.Execution{Status: "FAILED", Error: "CustomError", Cause: "boom"},
			"CustomError",
			"boom",
		},
		{
			"errorless child",
			&sfnstore.Execution{Status: "FAILED", Cause: "no error name"},
			"States.TaskFailed",
			"no error name",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := e.childSyncResult(tc.child, false)
			tf, ok := err.(*taskFailure)
			if !ok {
				t.Fatalf("childSyncResult error = %v, want a taskFailure", err)
			}
			if tf.code != tc.wantCode {
				t.Errorf("code = %s, want %s", tf.code, tc.wantCode)
			}
			if tf.cause != tc.wantCause {
				t.Errorf("cause = %q, want %q", tf.cause, tc.wantCause)
			}
		})
	}
}

// TestStartExecutionIntegrationCallback: the callback form passes the task
// token into the child input and resumes on SendTaskSuccess.
func TestStartExecutionIntegrationCallback(t *testing.T) {
	e, store := newStartExecutionHarness(t)
	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::states:startExecution.waitForTaskToken",
		End:      true,
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"StateMachineArn": startExecutionTestTargetARN,
			"Name":            "child-cb",
			"Input": map[string]interface{}{
				"token.$": "$$.Task.Token",
			},
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

	childArn := "arn:aws:states:us-east-1:000000000000:execution:target-sm:child-cb"
	var token string
	awaitCondition(t, "the child execution to carry the token", func() bool {
		exec, err := store.GetExecution(context.Background(), childArn)
		if err != nil || exec == nil {
			return false
		}
		var input map[string]interface{}
		if jerr := json.Unmarshal([]byte(exec.Input), &input); jerr == nil {
			token, _ = input["token"].(string)
		}
		return token != ""
	})

	if err := store.CompleteActivityTask(token, `{"done":true}`); err != nil {
		t.Fatalf("SendTaskSuccess failed: %v", err)
	}

	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("callback startExecution failed: %v", r.err)
		}
		if strings.TrimSpace(r.output) != `{"done":true}` {
			t.Fatalf("task output must be the SendTaskSuccess output, got %s", r.output)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("callback startExecution did not resume")
	}

	types := historyTypes(store, execCtx.Execution.ExecutionArn)
	found := false
	for _, typ := range types {
		if typ == "TaskSubmitted" {
			found = true
		}
	}
	if !found {
		t.Fatalf("history lacks TaskSubmitted: %v", types)
	}
}

// TestTaskResourceSupportMatrix pins the validator's per-family grammar:
// every carried family accepts exactly the patterns AWS documents for it,
// the AWS SDK namespace accepts the plain and callback patterns for any
// service, and resources naming uncarried services keep their AWS-side
// validity (they fail at run time instead).
func TestTaskResourceSupportMatrix(t *testing.T) {
	cases := []struct {
		resource string
		wantErr  bool
	}{
		{"arn:aws:states:::sqs:sendMessage", false},
		{"arn:aws:states:::sqs:sendMessage.waitForTaskToken", false},
		{"arn:aws:states:::sqs:sendMessage.sync", true},
		{"arn:aws:states:::sns:publish", false},
		{"arn:aws:states:::sns:publish.waitForTaskToken", false},
		{"arn:aws:states:::events:putEvents.waitForTaskToken", false},
		{"arn:aws:states:::dynamodb:putItem", false},
		{"arn:aws:states:::dynamodb:updateItem.waitForTaskToken", true},
		{"arn:aws:states:::lambda:invoke", false},
		{"arn:aws:states:::lambda:invoke.waitForTaskToken", false},
		{"arn:aws:states:::lambda:invoke.sync", true},
		{"arn:aws:states:::states:startExecution", false},
		{"arn:aws:states:::states:startExecution.sync", false},
		{"arn:aws:states:::states:startExecution.sync:2", false},
		{"arn:aws:states:::states:startExecution.waitForTaskToken", false},
		{"arn:aws:lambda:us-east-1:000000000000:function:fn", false},
		{"arn:aws:lambda:us-east-1:000000000000:function:fn.waitForTaskToken", true},
		{"arn:aws:states:us-east-1:000000000000:activity:work", false},
		{"arn:aws:states:us-east-1:000000000000:activity:work.sync", true},
		{"arn:aws:states:::aws-sdk:sqs:sendMessage", false},
		{"arn:aws:states:::aws-sdk:lambda:invoke.waitForTaskToken", false},
		{"arn:aws:states:::aws-sdk:sns:publish.sync", true},
		{"arn:aws:states:::batch:submitJob.sync", false},
		{"arn:aws:states:::sqs:receiveMessage", true},
		{"arn:aws:states:::states:describeExecution", true},
	}
	for _, tc := range cases {
		def := fmt.Sprintf(`{"StartAt":"T","States":{"T":{"Type":"Task","Resource":%q,"End":true}}}`, tc.resource)
		diags := validateASLStructure(def, "STANDARD")
		hasErr := false
		for _, d := range diags {
			if d.Severity == "ERROR" && d.Location == "/States/T/Resource" {
				hasErr = true
			}
		}
		if hasErr != tc.wantErr {
			t.Errorf("resource %s: wantErr=%v, diagnostics=%+v", tc.resource, tc.wantErr, diags)
		}
	}
}

// TestSendTaskHeartbeatAcceptsCallbackTask pins the callback heartbeat
// path: a .waitForTaskToken task record is RUNNING from creation, so a
// consumer's SendTaskHeartbeat lands (and the heartbeat window extends)
// instead of failing with TaskDoesNotExist.
func TestSendTaskHeartbeatAcceptsCallbackTask(t *testing.T) {
	svc, store := newRecoveryService(t)
	ctx := t.Context()

	task := &sfnstore.ActivityTask{
		TaskToken:    "token-callback-hb",
		ExecutionArn: "arn:aws:states:us-east-1:000000000000:execution:hb-sm:exec1",
		Input:        `{}`,
	}
	if err := store.CreateCallbackTask(task); err != nil {
		t.Fatalf("create callback task: %v", err)
	}

	if err := svc.sendTaskHeartbeatCore(ctx, store, task.TaskToken); err != nil {
		t.Fatalf("SendTaskHeartbeat on a live callback task = %v, want nil", err)
	}
	updated, err := store.GetActivityTaskByToken(task.TaskToken)
	if err != nil {
		t.Fatal(err)
	}
	if updated.LastHeartbeatAt.IsZero() {
		t.Error("the heartbeat never landed on the task record")
	}
}

// TestTaskWaitTimeoutDefault pins the unset-TimeoutSeconds wait for both
// token-waiting paths (activity workers and .waitForTaskToken callbacks):
// the Task-state reference fixes the default at 99,999,999 seconds, not a
// transport-style short cap.
func TestTaskWaitTimeoutDefault(t *testing.T) {
	if got := taskWaitTimeout(0); got != taskTimeoutDefaultSeconds*time.Second {
		t.Errorf("unset TimeoutSeconds waits %v, want the documented %d-second default", got, taskTimeoutDefaultSeconds)
	}
	if got := taskWaitTimeout(30); got != 30*time.Second {
		t.Errorf("TimeoutSeconds=30 waits %v, want 30s", got)
	}
}

package sfn

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// recordingTaskAuthz records every authorisation request the dispatch path
// makes and answers with the configured verdict.
type recordingTaskAuthz struct {
	mu       sync.Mutex
	calls    [][3]string
	verdicts []error
}

func (r *recordingTaskAuthz) AuthoriseTaskCredentials(ctx context.Context, taskRoleArn, action, resourceArn string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, [3]string{taskRoleArn, action, resourceArn})
	if len(r.verdicts) > 0 {
		err := r.verdicts[0]
		r.verdicts = r.verdicts[1:]
		return err
	}
	return nil
}

func (r *recordingTaskAuthz) recorded() [][3]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([][3]string(nil), r.calls...)
}

// authzSQSInvoker resolves queue names and URLs to ARNs so the SQS target
// derivation has the same store-backed resolution the send path uses.
type authzSQSInvoker struct {
	invokers.SQSInvoker
}

func (f *authzSQSInvoker) GetQueueByName(ctx context.Context, region, queueName string) (string, error) {
	return "https://sqs." + region + ".amazonaws.com/123456789012/" + queueName, nil
}

func (f *authzSQSInvoker) GetQueueARN(ctx context.Context, region, queueURL string) (string, error) {
	name := queueURL[strings.LastIndex(queueURL, "/")+1:]
	return "arn:aws:sqs:" + region + ":123456789012:" + name, nil
}

func newAuthzExecutor(t *testing.T, authz TaskCredentialsAuthoriser) *Executor {
	t.Helper()
	store := newMapTestStore(t)
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(&authzSQSInvoker{})
	return NewExecutorWithStores(store, bus, "123456789012", "us-east-1", authz)
}

// The target derivation mirrors the dispatch families: each integration
// names its IAM action and derives the resource ARN from the parameters.
func TestIntegrationAuthorisationTargets(t *testing.T) {
	e := newAuthzExecutor(t, nil)

	cases := []struct {
		name       string
		resource   string
		input      string
		wantAction string
		wantARN    string
	}{
		{
			name:       "lambda ARN resource authorises the ARN itself",
			resource:   "arn:aws:lambda:us-east-1:123456789012:function:Echo",
			wantAction: "lambda:InvokeFunction",
			wantARN:    "arn:aws:lambda:us-east-1:123456789012:function:Echo",
		},
		{
			name:       "optimised lambda invoke with an ARN FunctionName",
			resource:   "arn:aws:states:::lambda:invoke",
			input:      `{"FunctionName":"arn:aws:lambda:us-east-1:123456789012:function:Echo"}`,
			wantAction: "lambda:InvokeFunction",
			wantARN:    "arn:aws:lambda:us-east-1:123456789012:function:Echo",
		},
		{
			name:       "lambda invoke with a bare name builds the function ARN",
			resource:   "arn:aws:states:::aws-sdk:lambda:invoke",
			input:      `{"FunctionName":"Echo"}`,
			wantAction: "lambda:InvokeFunction",
			wantARN:    "arn:aws:lambda:us-east-1:123456789012:function:Echo",
		},
		{
			name:       "sqs sendMessage resolves the queue ARN from the URL",
			resource:   "arn:aws:states:::sqs:sendMessage",
			input:      `{"QueueUrl":"https://sqs.us-east-1.amazonaws.com/123456789012/queue-a","MessageBody":"m"}`,
			wantAction: "sqs:SendMessage",
			wantARN:    "arn:aws:sqs:us-east-1:123456789012:queue-a",
		},
		{
			name:       "sqs sendMessage resolves the queue ARN from the name",
			resource:   "arn:aws:states:::aws-sdk:sqs:sendMessage",
			input:      `{"QueueName":"queue-a","MessageBody":"m"}`,
			wantAction: "sqs:SendMessage",
			wantARN:    "arn:aws:sqs:us-east-1:123456789012:queue-a",
		},
		{
			name:       "sns publish authorises the topic ARN",
			resource:   "arn:aws:states:::sns:publish",
			input:      `{"TopicArn":"arn:aws:sns:us-east-1:123456789012:topic-a","Message":"m"}`,
			wantAction: "sns:Publish",
			wantARN:    "arn:aws:sns:us-east-1:123456789012:topic-a",
		},
		{
			name:       "sns publish with a bare topic name builds the topic ARN",
			resource:   "arn:aws:states:::aws-sdk:sns:publish",
			input:      `{"TopicName":"topic-a","Message":"m"}`,
			wantAction: "sns:Publish",
			wantARN:    "arn:aws:sns:us-east-1:123456789012:topic-a",
		},
		{
			name:       "events putEvents authorises the named bus",
			resource:   "arn:aws:states:::events:putEvents",
			input:      `{"Entries":[{"EventBusName":"bus-a","Source":"s","DetailType":"d","Detail":"{}"}]}`,
			wantAction: "events:PutEvents",
			wantARN:    "arn:aws:events:us-east-1:123456789012:event-bus/bus-a",
		},
		{
			name:       "events putEvents without a bus name authorises the default bus",
			resource:   "arn:aws:states:::aws-sdk:eventbridge:putEvents",
			input:      `{"Entries":[{"Source":"s","DetailType":"d","Detail":"{}"}]}`,
			wantAction: "events:PutEvents",
			wantARN:    "arn:aws:events:us-east-1:123456789012:event-bus/default",
		},
		{
			name:       "dynamodb putItem authorises the table ARN",
			resource:   "arn:aws:states:::dynamodb:putItem",
			input:      `{"TableName":"table-a","Item":{}}`,
			wantAction: "dynamodb:PutItem",
			wantARN:    "arn:aws:dynamodb:us-east-1:123456789012:table/table-a",
		},
		{
			name:       "dynamodb getItem via the SDK namespace",
			resource:   "arn:aws:states:::aws-sdk:dynamodb:getItem",
			input:      `{"TableName":"table-a","Key":{}}`,
			wantAction: "dynamodb:GetItem",
			wantARN:    "arn:aws:dynamodb:us-east-1:123456789012:table/table-a",
		},
		{
			name:       "startExecution authorises the child state machine",
			resource:   "arn:aws:states:::states:startExecution",
			input:      `{"StateMachineArn":"arn:aws:states:us-east-1:123456789012:stateMachine:child"}`,
			wantAction: "states:StartExecution",
			wantARN:    "arn:aws:states:us-east-1:123456789012:stateMachine:child",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			targets := e.integrationAuthorisationTargets(context.Background(), tc.resource, tc.input, nil)
			if len(targets) != 1 {
				t.Fatalf("expected exactly one target, got %d: %+v", len(targets), targets)
			}
			if targets[0].action != tc.wantAction {
				t.Errorf("action = %s, want %s", targets[0].action, tc.wantAction)
			}
			if targets[0].resourceArn != tc.wantARN {
				t.Errorf("resourceArn = %s, want %s", targets[0].resourceArn, tc.wantARN)
			}
		})
	}
}

// Activity resources carry no Credentials contract and no SFN-side
// invocation, so they yield no authorisation target.
func TestIntegrationAuthorisationTargetsSkipsActivities(t *testing.T) {
	e := newAuthzExecutor(t, nil)
	targets := e.integrationAuthorisationTargets(context.Background(), tokenTestActivityARN, `{}`, nil)
	if len(targets) != 0 {
		t.Fatalf("expected no targets for an activity resource, got %+v", targets)
	}
}

// countingSQSInvoker counts the queue-name resolutions an authorised
// QueueName-form SendMessage dispatch performs.
type countingSQSInvoker struct {
	invokers.SQSInvoker
	mu        sync.Mutex
	nameCalls int
}

func (f *countingSQSInvoker) GetQueueByName(ctx context.Context, region, queueName string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nameCalls++
	return "https://sqs." + region + ".amazonaws.com/123456789012/" + queueName, nil
}

func (f *countingSQSInvoker) GetQueueARN(ctx context.Context, region, queueURL string) (string, error) {
	name := queueURL[strings.LastIndex(queueURL, "/")+1:]
	return "arn:aws:sqs:" + region + ":123456789012:" + name, nil
}

func (f *countingSQSInvoker) SendMessage(ctx context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	return "msg-1", "md5-body", nil
}

func (f *countingSQSInvoker) queueNameResolutions() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.nameCalls
}

// TestAuthorisedQueueNameSendResolvesTheQueueOnce pins the shared target
// resolution: the Credentials authorisation and the SendMessage dispatch
// derive the queue from one resolver, so an authorised QueueName-form send
// resolves the queue name exactly once per attempt and authorises the ARN
// the send actually targets.
func TestAuthorisedQueueNameSendResolvesTheQueueOnce(t *testing.T) {
	sqs := &countingSQSInvoker{}
	store := newMapTestStore(t)
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	authz := &recordingTaskAuthz{}
	e := NewExecutorWithStores(store, bus, "123456789012", "us-east-1", authz)
	e.region = "us-east-1"

	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::sqs:sendMessage",
		End:      true,
		Credentials: map[string]interface{}{
			"RoleArn": "arn:aws:iam::123456789012:role/task",
		},
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"QueueName":   "queue-a",
			"MessageBody": "m",
		}},
	}
	if _, _, execErr := e.executeTask(context.Background(), execCtx, state); execErr != nil {
		t.Fatalf("authorised send failed: %v", execErr)
	}

	if calls := sqs.queueNameResolutions(); calls != 1 {
		t.Fatalf("queue-name resolutions = %d, want 1 (the authorisation and the send share one resolution)", calls)
	}
	recorded := authz.recorded()
	if len(recorded) != 1 || recorded[0][1] != "sqs:SendMessage" || recorded[0][2] != "arn:aws:sqs:us-east-1:123456789012:queue-a" {
		t.Fatalf("authorisation calls = %v, want one sqs:SendMessage on the resolved queue ARN", recorded)
	}
}

// A denial fails the attempt before any integration call: the state fails
// with States.Permissions, the invocation never reaches the service, and
// the history records the start-failure event for the resource class.
func TestTaskCredentialsDenialFailsInvocationBeforeDispatch(t *testing.T) {
	lambda := &callbackLambdaInvoker{statusCode: 200, response: []byte(`{"ok":true}`)}
	store := newMapTestStore(t)
	bus := eventbus.NewEventBus()
	bus.SetLambdaInvoker(lambda)
	authz := &recordingTaskAuthz{verdicts: []error{errors.New("role is not authorised for lambda:InvokeFunction")}}
	e := NewExecutorWithStores(store, bus, "123456789012", "us-east-1", authz)

	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:lambda:us-east-1:123456789012:function:Echo",
		End:      true,
		Credentials: map[string]interface{}{
			"RoleArn": "arn:aws:iam::123456789012:role/TaskRole",
		},
	}

	_, _, execErr := e.executeTask(context.Background(), execCtx, state)
	if execErr == nil {
		t.Fatal("expected the denied credentials to fail the task")
	}
	if execErr.ErrorCode != "States.Permissions" {
		t.Fatalf("error code = %s, want States.Permissions", execErr.ErrorCode)
	}
	if len(lambda.names) != 0 {
		t.Fatalf("the invocation must not reach the service under a denial, got %d calls", len(lambda.names))
	}

	calls := authz.recorded()
	if len(calls) != 1 || calls[0][0] != "arn:aws:iam::123456789012:role/TaskRole" ||
		calls[0][1] != "lambda:InvokeFunction" || calls[0][2] != "arn:aws:lambda:us-east-1:123456789012:function:Echo" {
		t.Fatalf("unexpected authorisation calls: %+v", calls)
	}

	history := historyTypes(store, execCtx.Execution.ExecutionArn)
	want := "LambdaFunctionStartFailed"
	found := false
	for _, typ := range history {
		if typ == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a %s event in the history, got %v", want, history)
	}
}

// An allowed authorisation passes the invocation through unchanged.
func TestTaskCredentialsAllowPassesInvocationThrough(t *testing.T) {
	lambda := &callbackLambdaInvoker{statusCode: 200, response: []byte(`{"ok":true}`)}
	store := newMapTestStore(t)
	bus := eventbus.NewEventBus()
	bus.SetLambdaInvoker(lambda)
	authz := &recordingTaskAuthz{}
	e := NewExecutorWithStores(store, bus, "123456789012", "us-east-1", authz)

	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:lambda:us-east-1:123456789012:function:Echo",
		End:      true,
		Credentials: map[string]interface{}{
			"RoleArn": "arn:aws:iam::123456789012:role/TaskRole",
		},
	}

	output, _, execErr := e.executeTask(context.Background(), execCtx, state)
	if execErr != nil {
		t.Fatalf("expected the allowed credentials to invoke, got %v", execErr)
	}
	if len(lambda.names) != 1 || lambda.names[0] != "arn:aws:lambda:us-east-1:123456789012:function:Echo" {
		t.Fatalf("unexpected invocations: %+v", lambda.names)
	}
	if output != `{"ok":true}` {
		t.Fatalf("output = %s, want the invocation payload", output)
	}
	if calls := authz.recorded(); len(calls) != 1 {
		t.Fatalf("expected one authorisation call, got %+v", calls)
	}
}

// The scheduled event carries the taskCredentials member for the
// non-Lambda resource classes too — the model defines it on
// TaskScheduledEventDetails.
func TestTaskScheduledEventCarriesTaskCredentials(t *testing.T) {
	store := newMapTestStore(t)
	e := NewExecutorWithStores(store, eventbus.NewEventBus(), "123456789012", "us-east-1", &recordingTaskAuthz{})

	execCtx := newTaskTokenExecCtx(`{}`)
	state := &sfnstore.TaskState{
		Type:     "Task",
		Resource: "arn:aws:states:::sns:publish",
		End:      true,
		Credentials: map[string]interface{}{
			"RoleArn": "arn:aws:iam::123456789012:role/TaskRole",
		},
		Parameters: &sfnstore.Parameters{Values: map[string]interface{}{
			"TopicArn": "arn:aws:sns:us-east-1:123456789012:topic-a",
			"Message":  "m",
		}},
	}

	// The invocation itself fails (no SNS invoker on the bus); the
	// assertion targets the TaskScheduled event that precedes it.
	_, _, _ = e.executeTask(context.Background(), execCtx, state)

	events, _, err := store.GetExecutionHistory(context.Background(), execCtx.Execution.ExecutionArn, 200, "", false)
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range events {
		if ev.Type == "TaskScheduled" {
			if ev.TaskScheduledEventDetails == nil || ev.TaskScheduledEventDetails.TaskCredentials == nil ||
				ev.TaskScheduledEventDetails.TaskCredentials.RoleArn != "arn:aws:iam::123456789012:role/TaskRole" {
				t.Fatalf("TaskScheduled event lost the taskCredentials member: %+v", ev.TaskScheduledEventDetails)
			}
			return
		}
	}
	t.Fatal("no TaskScheduled event in the history")
}

package scheduler

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// TestUniversalTargetValidation pins the creation-time contract of the
// universal-target ARN form: every aws-sdk ARN with a well-formed
// {service}:{action} pair is accepted (the AWS contract — delivery decides
// per target), read-only action prefixes are rejected, Input must be
// well-formed JSON, templated sub-parameters are rejected, and a
// scheduler-service ARN outside the aws-sdk form names nothing deliverable.
func TestUniversalTargetValidation(t *testing.T) {
	roleArn := "arn:aws:iam::000000000000:role/universal-pin"

	cases := []struct {
		name                    string
		target                  *schedulerstore.Target
		wantErr                 bool
		wantInMsg               string
		wantValidationException bool
	}{
		{
			name: "delivered operation accepted",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:lambda:invoke",
				RoleArn: roleArn,
				Input:   `{"FunctionName":"arn:aws:lambda:us-east-1:000000000000:function:f"}`,
			},
		},
		{
			name: "unimplemented service accepted",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:batch:submitJob",
				RoleArn: roleArn,
				Input:   `{"JobName":"pin"}`,
			},
		},
		{
			name: "excluded service accepted",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:sagemaker:createPipeline",
				RoleArn: roleArn,
				Input:   `{"PipelineName":"pin"}`,
			},
		},
		{
			name: "read-only prefix rejected",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:sqs:getQueueUrl",
				RoleArn: roleArn,
			},
			wantErr:                 true,
			wantInMsg:               "read-only",
			wantValidationException: true,
		},
		{
			name: "read-only prefix matches case-insensitively",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:mq:ListBrokers",
				RoleArn: roleArn,
			},
			wantErr:   true,
			wantInMsg: "read-only",
		},
		// The multi-word prefixes themselves carry inner capitals in the
		// documentation; the action's case must not decide the match.
		{
			name: "multi-word read-only prefix rejected (BatchGetItem)",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:dynamodb:BatchGetItem",
				RoleArn: roleArn,
			},
			wantErr:   true,
			wantInMsg: "read-only",
		},
		{
			name: "multi-word read-only prefix rejected (TransactGetItems)",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:dynamodb:TransactGetItems",
				RoleArn: roleArn,
			},
			wantErr:   true,
			wantInMsg: "read-only",
		},
		{
			name: "multi-word read-only prefix rejected (InvokeModel)",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:bedrock:InvokeModel",
				RoleArn: roleArn,
			},
			wantErr:   true,
			wantInMsg: "read-only",
		},
		{
			name: "non-JSON Input rejected",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:lambda:invoke",
				RoleArn: roleArn,
				Input:   `not-json`,
			},
			wantErr:   true,
			wantInMsg: "well-formed JSON",
		},
		{
			name: "templated sub-parameters rejected",
			target: &schedulerstore.Target{
				Arn:           "arn:aws:scheduler:::aws-sdk:sqs:sendMessage",
				RoleArn:       roleArn,
				Input:         `{"QueueUrl":"http://q","MessageBody":"m"}`,
				SqsParameters: &schedulerstore.SqsParameters{MessageGroupId: "g"},
			},
			wantErr:   true,
			wantInMsg: "templated target parameters",
		},
		{
			name: "malformed scheduler-service ARN rejected",
			target: &schedulerstore.Target{
				Arn:     "arn:aws:scheduler:::aws-sdk:sqs",
				RoleArn: roleArn,
			},
			wantErr:   true,
			wantInMsg: "universal-target form",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTarget(tt.target)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("validateTarget accepted %s", tt.target.Arn)
				}
				if !strings.Contains(err.Error(), tt.wantInMsg) {
					t.Fatalf("error %q does not name %q", err.Error(), tt.wantInMsg)
				}
				if tt.wantValidationException {
					var awsErr *awserrors.AWSError
					if !errors.As(err, &awsErr) || awsErr.Code != "ValidationException" {
						t.Fatalf("error must be the modelled ValidationException, got %#v", err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("validateTarget rejected %s: %v", tt.target.Arn, err)
			}
		})
	}
}

// recordingUniversalLambdaInvoker captures the Invoke contract the
// universal lambda deliverer addresses.
type recordingUniversalLambdaInvoker struct {
	mu           sync.Mutex
	triggerCalls []string // functionName per InvokeForTrigger
	gatewayCalls []string // functionName per InvokeForGateway
	lastPayload  []byte
}

func (r *recordingUniversalLambdaInvoker) InvokeForGateway(_ context.Context, functionName string, payload []byte) (int64, []byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gatewayCalls = append(r.gatewayCalls, functionName)
	r.lastPayload = payload
	return 202, nil, nil
}

func (r *recordingUniversalLambdaInvoker) InvokeForTrigger(_ context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.triggerCalls = append(r.triggerCalls, functionName)
	r.lastPayload = payload
	return invokers.LambdaInvocation{StatusCode: 200}, nil
}

func (r *recordingUniversalLambdaInvoker) GetFunctionARN(_ context.Context, functionName string) (string, error) {
	return functionName, nil
}

// recordingUniversalSQSInvoker captures the SendMessage contract.
type recordingUniversalSQSInvoker struct {
	mu       sync.Mutex
	regions  []string
	bodies   []string
	lastOpts invokers.SQSSendOptions
	lastURL  string
}

func (r *recordingUniversalSQSInvoker) GetQueueByName(context.Context, string, string) (string, error) {
	return "", fmt.Errorf("not expected in the universal pin")
}

func (r *recordingUniversalSQSInvoker) GetQueueARN(context.Context, string, string) (string, error) {
	return "", nil
}

func (r *recordingUniversalSQSInvoker) SendMessage(_ context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.regions = append(r.regions, region)
	r.bodies = append(r.bodies, body)
	r.lastOpts = opts
	r.lastURL = queueURL
	return "id", "md5", nil
}

func (r *recordingUniversalSQSInvoker) ReceiveMessage(context.Context, string, string, int32, *int32, int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, nil
}

func (r *recordingUniversalSQSInvoker) DeleteMessage(context.Context, string, string, string) error {
	return nil
}

// recordingUniversalSNSInvoker captures the PublishToTopic contract.
type recordingUniversalSNSInvoker struct {
	mu         sync.Mutex
	calls      int
	topic      string
	message    string
	subject    string
	attributes map[string]string
	publishErr error
}

func (r *recordingUniversalSNSInvoker) GetTopic(context.Context, string) (string, error) {
	return "", nil
}

func (r *recordingUniversalSNSInvoker) GetTopicPolicy(context.Context, string) (string, error) {
	return "", nil
}

func (r *recordingUniversalSNSInvoker) ListSubscriptionsByTopic(context.Context, string) ([]invokers.SubscriptionInfo, error) {
	return nil, nil
}

func (r *recordingUniversalSNSInvoker) PublishToTopic(_ context.Context, topicArn, message, subject string, messageAttributes map[string]string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.publishErr != nil {
		return "", r.publishErr
	}
	r.calls++
	r.topic = topicArn
	r.message = message
	r.subject = subject
	r.attributes = messageAttributes
	return "msg-id", nil
}

func (r *recordingUniversalSNSInvoker) StoreMessage(context.Context, string, any) error { return nil }

func (r *recordingUniversalSNSInvoker) DeleteStoredMessage(context.Context, string) error { return nil }

// universalDeliveryEngine wires a fresh engine over a real bus with the
// given invokers and the synchronous captures for the bus-event families.
func universalDeliveryEngine(t *testing.T, l invokers.LambdaInvoker, s invokers.SQSInvoker, n invokers.SNSInvoker, k invokers.KinesisInvoker) (*Engine, *eventbus.EventBus) {
	t.Helper()
	svc := newLifecycleService(t)
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	if l != nil {
		bus.SetLambdaInvoker(l)
	}
	if s != nil {
		bus.SetSQSInvoker(s)
	}
	if n != nil {
		bus.SetSNSInvoker(n)
	}
	if k != nil {
		bus.SetKinesisInvoker(k)
	}
	svc.engine.SetEventBus(bus)
	return svc.engine, bus
}

func universalSchedule(target *schedulerstore.Target) *schedulerstore.Schedule {
	return &schedulerstore.Schedule{
		Name:      "universal-pin",
		GroupName: "default",
		Region:    "us-east-1",
		Target:    target,
	}
}

// TestUniversalLambdaDeliveryPins translates the Invoke request: the
// payload's string form unwraps to its bytes, the default and Event
// invocation types route to the trigger and gateway contracts, and an
// invalid InvocationType fails with a cause.
func TestUniversalLambdaDeliveryPins(t *testing.T) {
	fnARN := "arn:aws:lambda:us-east-1:000000000000:function:universal-pin"
	invoker := &recordingUniversalLambdaInvoker{}
	engine, _ := universalDeliveryEngine(t, invoker, nil, nil, nil)

	target := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:lambda:invoke",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"FunctionName":%q,"Payload":"{\"k\":1}"}`, fnARN),
	}
	if err := engine.deliverUniversalLambda(t.Context(), universalSchedule(target), target); err != nil {
		t.Fatalf("trigger invocation: %v", err)
	}
	if len(invoker.triggerCalls) != 1 || invoker.triggerCalls[0] != fnARN {
		t.Fatalf("InvokeForTrigger calls = %v, want [%s]", invoker.triggerCalls, fnARN)
	}
	if string(invoker.lastPayload) != `{"k":1}` {
		t.Fatalf("payload = %q, want the unwrapped JSON string bytes", invoker.lastPayload)
	}

	invoker2 := &recordingUniversalLambdaInvoker{}
	engine2, _ := universalDeliveryEngine(t, invoker2, nil, nil, nil)
	target2 := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:lambda:invoke",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"FunctionName":%q,"InvocationType":"Event","Payload":"{}"}`, fnARN),
	}
	if err := engine2.deliverUniversalLambda(t.Context(), universalSchedule(target2), target2); err != nil {
		t.Fatalf("event invocation: %v", err)
	}
	if len(invoker2.gatewayCalls) != 1 || len(invoker2.triggerCalls) != 0 {
		t.Fatalf("Event invocation must route to the gateway contract, got gateway=%v trigger=%v", invoker2.gatewayCalls, invoker2.triggerCalls)
	}

	target3 := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:lambda:invoke",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"FunctionName":%q,"InvocationType":"Bogus"}`, fnARN),
	}
	if err := engine2.deliverUniversalLambda(t.Context(), universalSchedule(target3), target3); err == nil ||
		!strings.Contains(err.Error(), "InvocationType") {
		t.Fatalf("invalid InvocationType must fail naming the member, got %v", err)
	}
}

// TestUniversalSQSDeliveryPins translates the SendMessage request members
// and the region contract: an AWS-form QueueUrl addresses its own region,
// a host-local platform URL falls back to the schedule's region.
func TestUniversalSQSDeliveryPins(t *testing.T) {
	invoker := &recordingUniversalSQSInvoker{}
	engine, _ := universalDeliveryEngine(t, nil, invoker, nil, nil)

	target := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:sqs:sendMessage",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   `{"QueueUrl":"http://localhost:50080/000000000000/universal-q","MessageBody":"m","DelaySeconds":5,"MessageGroupId":"g","MessageDeduplicationId":"dedup-1","MessageAttributes":{"a":{"DataType":"String","StringValue":"v"},"b":{"DataType":"Binary","BinaryValue":"dmFs"}}}`,
	}
	if err := engine.deliverUniversalSQS(t.Context(), universalSchedule(target), target); err != nil {
		t.Fatalf("host-local queue URL: %v", err)
	}
	if len(invoker.bodies) != 1 || invoker.bodies[0] != "m" {
		t.Fatalf("SendMessage bodies = %v, want [m]", invoker.bodies)
	}
	if invoker.regions[0] != "us-east-1" {
		t.Fatalf("host-local URL region = %q, want the schedule's region us-east-1", invoker.regions[0])
	}
	if invoker.lastOpts.DelaySeconds != 5 || invoker.lastOpts.MessageGroupID != "g" {
		t.Fatalf("SendMessage opts = %+v, want DelaySeconds 5 and MessageGroupId g", invoker.lastOpts)
	}
	// The dedup id travels with the send: without it a FIFO queue without
	// content-based deduplication refuses the delivery AWS would accept.
	if invoker.lastOpts.MessageDeduplicationID != "dedup-1" {
		t.Fatalf("SendMessage opts = %+v, want MessageDeduplicationId dedup-1", invoker.lastOpts)
	}
	typed := invoker.lastOpts.TypedMessageAttributes
	if typed["a"].DataType != "String" || typed["a"].StringValue != "v" {
		t.Fatalf("String attribute mistranslated: %+v", typed["a"])
	}
	if string(typed["b"].BinaryValue) != "val" {
		t.Fatalf("Binary attribute must base64-decode to val, got %q", typed["b"].BinaryValue)
	}

	target2 := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:sqs:sendMessage",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   `{"QueueUrl":"https://sqs.eu-west-1.amazonaws.com/000000000000/universal-q","MessageBody":"m"}`,
	}
	if err := engine.deliverUniversalSQS(t.Context(), universalSchedule(target2), target2); err != nil {
		t.Fatalf("AWS-form queue URL: %v", err)
	}
	if invoker.regions[1] != "eu-west-1" {
		t.Fatalf("AWS-form URL region = %q, want eu-west-1", invoker.regions[1])
	}
}

// TestUniversalSNSDeliveryPins translates the Publish request members; a
// Binary message attribute fails the delivery with a cause instead of
// being dropped silently.
func TestUniversalSNSDeliveryPins(t *testing.T) {
	invoker := &recordingUniversalSNSInvoker{}
	engine, _ := universalDeliveryEngine(t, nil, nil, invoker, nil)

	topicArn := "arn:aws:sns:us-east-1:000000000000:universal-pin"
	target := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:sns:publish",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"TopicArn":%q,"Message":"msg","Subject":"subj","MessageAttributes":{"k":{"DataType":"String","StringValue":"v"}}}`, topicArn),
	}
	if err := engine.deliverUniversalSNS(t.Context(), universalSchedule(target), target); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if invoker.topic != topicArn || invoker.message != "msg" || invoker.subject != "subj" {
		t.Fatalf("PublishToTopic = %s/%s/%s, want topic/message/subj", invoker.topic, invoker.message, invoker.subject)
	}
	if invoker.attributes["k"] != "v" {
		t.Fatalf("MessageAttributes = %v, want k=v", invoker.attributes)
	}

	binaryTarget := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:sns:publish",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"TopicArn":%q,"Message":"msg","MessageAttributes":{"k":{"DataType":"Binary","BinaryValue":"dg=="}}}`, topicArn),
	}
	if err := engine.deliverUniversalSNS(t.Context(), universalSchedule(binaryTarget), binaryTarget); err == nil ||
		!strings.Contains(err.Error(), "Binary") {
		t.Fatalf("a Binary message attribute must fail naming the member, got %v", err)
	}
}

// recordingUniversalKinesisInvoker captures the full PutRecord contract:
// region, stream, partition key and decoded data.
type recordingUniversalKinesisInvoker struct {
	mu           sync.Mutex
	putRegions   map[string]string
	partitionKey string
	data         []byte
}

func (r *recordingUniversalKinesisInvoker) ListShards(context.Context, string, string) ([]invokers.ShardInfo, error) {
	return nil, nil
}

func (r *recordingUniversalKinesisInvoker) PutRecord(_ context.Context, region, streamName, partitionKey string, data []byte) (string, string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.putRegions == nil {
		r.putRegions = make(map[string]string)
	}
	r.putRegions[streamName] = region
	r.partitionKey = partitionKey
	r.data = data
	return "seq-1", "", nil
}

func (r *recordingUniversalKinesisInvoker) CreateShardIterator(context.Context, string, string, string, string, string, *time.Time) (string, error) {
	return "", nil
}

func (r *recordingUniversalKinesisInvoker) GetRecords(context.Context, string, string, string, string, int32, bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", nil
}

func (r *recordingUniversalKinesisInvoker) StreamExists(context.Context, string, string) (bool, error) {
	return false, nil
}

// TestUniversalKinesisDeliveryPins translates the PutRecord request: Data
// rides in its base64 wire form (the local store keeps the string as-is and
// GetRecords serves it to SDK clients verbatim), and the StreamARN form
// addresses the ARN's region while the StreamName form falls back to the
// schedule's region.
func TestUniversalKinesisDeliveryPins(t *testing.T) {
	invoker := &recordingUniversalKinesisInvoker{}
	engine, _ := universalDeliveryEngine(t, nil, nil, nil, invoker)

	streamARN := "arn:aws:kinesis:eu-west-1:000000000000:stream/universal-pin"
	dataB64 := base64.StdEncoding.EncodeToString([]byte("rec"))
	target := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:kinesis:putRecord",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"StreamARN":%q,"Data":%q,"PartitionKey":"pk"}`, streamARN, dataB64),
	}
	if err := engine.deliverUniversalKinesis(t.Context(), universalSchedule(target), target); err != nil {
		t.Fatalf("StreamARN form: %v", err)
	}
	if got := invoker.putRegions["universal-pin"]; got != "eu-west-1" {
		t.Fatalf("StreamARN region = %q, want eu-west-1", got)
	}
	if string(invoker.data) != dataB64 || invoker.partitionKey != "pk" {
		t.Fatalf("PutRecord data/partitionKey = %q/%q, want the base64 wire form %q/pk", invoker.data, invoker.partitionKey, dataB64)
	}

	nameTarget := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:kinesis:putRecord",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"StreamName":"universal-name-pin","Data":%q,"PartitionKey":"pk"}`, base64.StdEncoding.EncodeToString([]byte("rec"))),
	}
	if err := engine.deliverUniversalKinesis(t.Context(), universalSchedule(nameTarget), nameTarget); err != nil {
		t.Fatalf("StreamName form: %v", err)
	}
	if got := invoker.putRegions["universal-name-pin"]; got != "us-east-1" {
		t.Fatalf("StreamName region = %q, want the schedule's region us-east-1", got)
	}
}

// TestUniversalStepFunctionsDeliveryPins translates the StartExecution
// request onto the bus event, defaulting an omitted Input to the empty
// JSON object.
func TestUniversalStepFunctionsDeliveryPins(t *testing.T) {
	engine, bus := universalDeliveryEngine(t, nil, nil, nil, nil)

	var mu sync.Mutex
	var got []*eventbus.StepFunctionsStartExecutionEvent
	if _, err := eventbus.SubscribeTyped[*eventbus.StepFunctionsStartExecutionEvent](bus,
		func(_ context.Context, evt *eventbus.StepFunctionsStartExecutionEvent) eventbus.HandlerResult {
			mu.Lock()
			defer mu.Unlock()
			got = append(got, evt)
			return eventbus.HandlerResult{}
		}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	smARN := "arn:aws:states:us-east-1:000000000000:stateMachine:universal-pin"
	target := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:sfn:startExecution",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"StateMachineArn":%q,"Input":"{\"k\":1}"}`, smARN),
	}
	if err := engine.deliverUniversalStepFunctions(t.Context(), universalSchedule(target), target); err != nil {
		t.Fatalf("startExecution: %v", err)
	}

	noInputTarget := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:sfn:startExecution",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   fmt.Sprintf(`{"StateMachineArn":%q}`, smARN),
	}
	if err := engine.deliverUniversalStepFunctions(t.Context(), universalSchedule(noInputTarget), noInputTarget); err != nil {
		t.Fatalf("startExecution without Input: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("bus events = %d, want 2", len(got))
	}
	if got[0].StateMachineArn != smARN || got[0].Input != `{"k":1}` || got[0].Region != "us-east-1" {
		t.Fatalf("first event = %+v", got[0])
	}
	if got[1].Input != "{}" {
		t.Fatalf("omitted Input must default to {}, got %q", got[1].Input)
	}
}

// TestUniversalEventBridgeDeliveryPins translates each PutEvents entry
// onto its own bus event: the struct members are authoritative and the
// entry rides in Input so the handler applies Detail/Resources/Time with
// the PutEvents semantics; an ARN-form EventBusName addresses its region.
func TestUniversalEventBridgeDeliveryPins(t *testing.T) {
	engine, bus := universalDeliveryEngine(t, nil, nil, nil, nil)

	var mu sync.Mutex
	var got []*eventbus.EventBridgePutEventsEvent
	if _, err := eventbus.SubscribeTyped[*eventbus.EventBridgePutEventsEvent](bus,
		func(_ context.Context, evt *eventbus.EventBridgePutEventsEvent) eventbus.HandlerResult {
			mu.Lock()
			defer mu.Unlock()
			got = append(got, evt)
			return eventbus.HandlerResult{}
		}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	target := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:eventbridge:putEvents",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input: `{"Entries":[` +
			`{"Source":"vorpal.pin","DetailType":"universal-pin","Detail":"{\"k\":1}","EventBusName":"arn:aws:events:eu-west-1:000000000000:event-bus/named-pin"},` +
			`{"Source":"vorpal.pin2","DetailType":"universal-pin-2","Detail":"{\"k\":2}"}` +
			`]}`,
	}
	if err := engine.deliverUniversalEventBridge(t.Context(), universalSchedule(target), target); err != nil {
		t.Fatalf("putEvents: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("bus events = %d, want one per entry", len(got))
	}
	if got[0].Source != "vorpal.pin" || got[0].DetailType != "universal-pin" {
		t.Fatalf("first entry struct members = %+v", got[0])
	}
	if got[0].Region != "eu-west-1" {
		t.Fatalf("ARN-form EventBusName region = %q, want eu-west-1", got[0].Region)
	}
	if !strings.Contains(got[0].Input, `"Detail"`) {
		t.Fatalf("the entry must ride in Input for the handler's Detail semantics, got %s", got[0].Input)
	}
	if got[1].Region != "us-east-1" {
		t.Fatalf("name-form EventBusName region = %q, want the schedule's region us-east-1", got[1].Region)
	}

	emptyTarget := &schedulerstore.Target{
		Arn:     "arn:aws:scheduler:::aws-sdk:eventbridge:putEvents",
		RoleArn: "arn:aws:iam::000000000000:role/universal-pin",
		Input:   `{"Entries":[]}`,
	}
	if err := engine.deliverUniversalEventBridge(t.Context(), universalSchedule(emptyTarget), emptyTarget); err == nil ||
		!strings.Contains(err.Error(), "Entries") {
		t.Fatalf("an empty Entries array must fail naming the member, got %v", err)
	}
}

// TestUniversalTargetAcceptThenFailCauses pins the delivery-time causes of
// the accepted-but-undeliverable universal ARNs: the missing substrate for
// unimplemented services, the recorded exclusion for the permanently
// excluded families, and the family-operation gap for delivered services'
// other actions.
func TestUniversalTargetAcceptThenFailCauses(t *testing.T) {
	cases := []struct {
		arn       string
		input     string
		wantInMsg string
	}{
		{
			arn:       "arn:aws:scheduler:::aws-sdk:batch:submitJob",
			input:     `{"JobName":"pin"}`,
			wantInMsg: "batch service is not implemented",
		},
		{
			arn:       "arn:aws:scheduler:::aws-sdk:sagemaker:createPipeline",
			input:     `{"PipelineName":"pin"}`,
			wantInMsg: "permanently excluded",
		},
		{
			arn:       "arn:aws:scheduler:::aws-sdk:sqs:createQueue",
			input:     `{"QueueName":"pin"}`,
			wantInMsg: "not available in this deployment",
		},
		{
			arn:       "arn:aws:scheduler:::aws-sdk:lambda:invoke",
			input:     "",
			wantInMsg: "requires Target.Input",
		},
	}
	for _, tt := range cases {
		t.Run(tt.arn, func(t *testing.T) {
			// The empty-Input case reaches a family deliverer, which needs
			// its invoker before it parses the request.
			engine, _ := universalDeliveryEngine(t, &recordingUniversalLambdaInvoker{}, nil, nil, nil)
			target := &schedulerstore.Target{Arn: tt.arn, Input: tt.input}
			err := engine.deliverUniversalTarget(t.Context(), universalSchedule(target), target)
			if err == nil {
				t.Fatalf("delivery of %s must fail", tt.arn)
			}
			if !strings.Contains(err.Error(), tt.wantInMsg) {
				t.Fatalf("error %q does not name %q", err.Error(), tt.wantInMsg)
			}
		})
	}

	// A malformed stored ARN fails with a cause instead of dispatching.
	bare := &Engine{}
	target := &schedulerstore.Target{Arn: "arn:aws:scheduler:::not-the-form"}
	if err := bare.deliverUniversalTarget(t.Context(), universalSchedule(target), target); err == nil ||
		!strings.Contains(err.Error(), "malformed universal target ARN") {
		t.Fatalf("malformed ARN must fail with a cause, got %v", err)
	}
}

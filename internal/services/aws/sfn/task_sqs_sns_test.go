package sfn

import (
	"context"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
)

// stubSQSInvoker records the SendMessage calls the integrations make.
type stubSQSInvoker struct {
	sentBodies  []string
	sentOpts    []invokers.SQSSendOptions
	sentRegions []string
}

func (s *stubSQSInvoker) GetQueueByName(_ context.Context, _, _ string) (string, error) {
	return "https://sqs.us-east-1.amazonaws.com/000000000000/q", nil
}

func (s *stubSQSInvoker) GetQueueARN(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (s *stubSQSInvoker) SendMessage(_ context.Context, region, _ string, body string, opts invokers.SQSSendOptions) (string, string, error) {
	s.sentBodies = append(s.sentBodies, body)
	s.sentOpts = append(s.sentOpts, opts)
	s.sentRegions = append(s.sentRegions, region)
	return "msg-1", "md5-body", nil
}

func (s *stubSQSInvoker) ReceiveMessage(_ context.Context, _, _ string, _ int32, _ *int32, _ int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, nil
}

func (s *stubSQSInvoker) DeleteMessage(_ context.Context, _, _, _ string) error { return nil }

// stubSNSInvoker records the stored messages the SNS integration writes.
type stubSNSInvoker struct {
	stored []map[string]interface{}
}

func (s *stubSNSInvoker) GetTopic(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (s *stubSNSInvoker) ListSubscriptionsByTopic(_ context.Context, _ string) ([]invokers.SubscriptionInfo, error) {
	return nil, nil
}

func (s *stubSNSInvoker) PublishToTopic(_ context.Context, _, _, _ string, _ map[string]string) (string, error) {
	return "msg-1", nil
}

func (s *stubSNSInvoker) StoreMessage(_ context.Context, _ string, data any) error {
	m, _ := data.(map[string]interface{})
	s.stored = append(s.stored, m)
	return nil
}

func (s *stubSNSInvoker) DeleteStoredMessage(_ context.Context, _ string) error { return nil }

func newSQSTestExecutor(t *testing.T, sqs *stubSQSInvoker) *Executor {
	t.Helper()
	bus := eventbus.NewEventBus()
	bus.SetSQSInvoker(sqs)
	e := NewExecutor(nil, bus)
	e.region = "us-east-1"
	return e
}

// TestSQSSendMessageRequiresBody pins that a parameter set without
// MessageBody is an invocation failure: SendMessage's MessageBody is a
// required member, so the queue never receives the serialised parameter
// object as the message body.
func TestSQSSendMessageRequiresBody(t *testing.T) {
	sqs := &stubSQSInvoker{}
	e := newSQSTestExecutor(t, sqs)

	_, err := e.executeSQSTask(context.Background(),
		"arn:aws:states:::sqs:sendMessage",
		`{"QueueUrl":"https://sqs.us-east-1.amazonaws.com/000000000000/q","DelaySeconds":5}`, nil)
	if err == nil {
		t.Fatal("SendMessage without MessageBody must fail")
	}
	if len(sqs.sentBodies) != 0 {
		t.Fatalf("nothing may be sent: bodies = %v", sqs.sentBodies)
	}
}

// TestSQSSendMessageBodyAndTypedAttributes pins the body and attribute
// contract: a structured MessageBody serialises that object alone, and
// message attributes keep their DataType with BinaryValue decoded from
// the Base64 transport.
func TestSQSSendMessageBodyAndTypedAttributes(t *testing.T) {
	sqs := &stubSQSInvoker{}
	e := newSQSTestExecutor(t, sqs)

	input := `{"QueueUrl":"https://sqs.us-east-1.amazonaws.com/000000000000/q",` +
		`"MessageBody":{"Input":"x","MyTaskToken":"tok"},` +
		`"MessageAttributes":{"retries":{"DataType":"Number","StringValue":"3"},` +
		`"blob":{"DataType":"Binary","BinaryValue":"aGVsbG8="}}}`
	if _, err := e.executeSQSTask(context.Background(), "arn:aws:states:::sqs:sendMessage", input, nil); err != nil {
		t.Fatalf("sendMessage failed: %v", err)
	}

	if len(sqs.sentBodies) != 1 || sqs.sentBodies[0] != `{"Input":"x","MyTaskToken":"tok"}` {
		t.Fatalf("message body = %v, want the structured MessageBody alone", sqs.sentBodies)
	}
	attrs := sqs.sentOpts[0].TypedMessageAttributes
	if attrs["retries"].DataType != "Number" || attrs["retries"].StringValue != "3" {
		t.Errorf("retries attribute = %+v, want DataType Number with the value", attrs["retries"])
	}
	if attrs["blob"].DataType != "Binary" || string(attrs["blob"].BinaryValue) != "hello" {
		t.Errorf("blob attribute = %+v, want DataType Binary with the decoded bytes", attrs["blob"])
	}
}

// TestSNSPublishRequiresMessage pins that a parameter set without Message
// is an invocation failure: Publish's Message is a required member, so
// the topic never receives the serialised parameter object as the
// message.
func TestSNSPublishRequiresMessage(t *testing.T) {
	bus := eventbus.NewEventBus()
	sns := &stubSNSInvoker{}
	bus.SetSNSInvoker(sns)
	e := NewExecutor(nil, bus)
	e.region = "us-east-1"

	_, err := e.executeSNSTask(context.Background(), "arn:aws:states:::sns:publish",
		`{"TopicArn":"arn:aws:sns:us-east-1:000000000000:topic","Subject":"s"}`)
	if err == nil {
		t.Fatal("Publish without Message must fail")
	}
	if len(sns.stored) != 0 {
		t.Fatalf("nothing may be stored: %v", sns.stored)
	}
}

// TestSNSPublishMessageAndTypedAttributes pins the message and attribute
// contract: a structured Message serialises that value alone, and the
// stored message carries the typed attributes (DataType with the value,
// not a flattened string map).
func TestSNSPublishMessageAndTypedAttributes(t *testing.T) {
	bus := eventbus.NewEventBus()
	sns := &stubSNSInvoker{}
	bus.SetSNSInvoker(sns)
	e := NewExecutor(nil, bus)
	e.region = "us-east-1"

	input := `{"TopicArn":"arn:aws:sns:us-east-1:000000000000:topic",` +
		`"Message":{"alert":"high"},` +
		`"MessageAttributes":{"priority":{"DataType":"Number","StringValue":"1"}}}`
	if _, err := e.executeSNSTask(context.Background(), "arn:aws:states:::sns:publish", input); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	if len(sns.stored) != 1 || sns.stored[0]["Message"] != `{"alert":"high"}` {
		t.Fatalf("stored message = %v, want the structured Message alone", sns.stored)
	}
	attrs, _ := sns.stored[0]["MessageAttributes"].(map[string]interface{})
	priority, _ := attrs["priority"].(snsAttributeTransport)
	if priority.Type != "Number" || priority.StringValue != "1" {
		t.Errorf("priority attribute = %+v, want the typed Number attribute", attrs["priority"])
	}
}

// stubDynamoDBInvoker backs the DynamoDB integration pins.
type stubDynamoDBInvoker struct {
	item map[string]interface{}
}

func (s *stubDynamoDBInvoker) ContributorRules(_ context.Context, _ string) ([]invokers.ContributorInsightRule, error) {
	return nil, nil
}

func (s *stubDynamoDBInvoker) ContributorStats(_ context.Context, _, _ string, _ string, _, _ time.Time, _ int) ([]invokers.ContributorKeyStat, error) {
	return nil, nil
}

func (s *stubDynamoDBInvoker) GetItem(_ context.Context, _, _ string, _ map[string]interface{}) (map[string]interface{}, error) {
	return s.item, nil
}

func (s *stubDynamoDBInvoker) PutItem(_ context.Context, _, _ string, key, _ map[string]interface{}) (map[string]interface{}, error) {
	return key, nil
}

func (s *stubDynamoDBInvoker) DeleteItem(_ context.Context, _, _ string, _ map[string]interface{}) error {
	return nil
}

func (s *stubDynamoDBInvoker) Scan(_ context.Context, _, _ string, _ int) ([]map[string]interface{}, error) {
	return nil, nil
}

func (s *stubDynamoDBInvoker) Query(_ context.Context, _, _, _ string, _ int) ([]map[string]interface{}, error) {
	return nil, nil
}

func (s *stubDynamoDBInvoker) UpdateItem(_ context.Context, _, _ string, _ map[string]interface{}, _ map[string]interface{}) error {
	return nil
}

func (s *stubDynamoDBInvoker) ScanWithPagination(_ context.Context, _, _ string, _ int, _ string) ([]map[string]interface{}, string, error) {
	return nil, "", nil
}

func (s *stubDynamoDBInvoker) QueryWithPagination(_ context.Context, _, _, _ string, _ int, _ string) ([]map[string]interface{}, string, error) {
	return nil, "", nil
}

// TestDynamoDBGetItemOmitsItemWhenAbsent pins the GetItem result for a
// key that matches nothing: the Item element is omitted entirely, never
// serialised as an explicit null.
func TestDynamoDBGetItemOmitsItemWhenAbsent(t *testing.T) {
	bus := eventbus.NewEventBus()
	db := &stubDynamoDBInvoker{}
	bus.SetDynamoDBInvoker(db)
	e := NewExecutor(nil, bus)
	e.region = "us-east-1"

	out, err := e.executeDynamoDBTask(context.Background(),
		"arn:aws:states:::dynamodb:getItem",
		`{"TableName":"t","Key":{"id":{"S":"missing"}}}`)
	if err != nil {
		t.Fatalf("getItem on a missing key failed: %v", err)
	}
	if out != "{}" {
		t.Fatalf("getItem result = %s, want the empty object with no Item member", out)
	}

	db.item = map[string]interface{}{"id": "present"}
	out, err = e.executeDynamoDBTask(context.Background(),
		"arn:aws:states:::dynamodb:getItem",
		`{"TableName":"t","Key":{"id":{"S":"present"}}}`)
	if err != nil {
		t.Fatalf("getItem failed: %v", err)
	}
	if !strings.Contains(out, `"Item"`) || strings.Contains(out, "null") {
		t.Fatalf("matched getItem result = %s, want the Item element", out)
	}
}

// TestActivityResourceAnyPartition pins the activity dispatch: it parses
// the ARN instead of matching the aws partition's literal prefix, so an
// activity in any partition dispatches (and reports its resource type)
// while non-activity resources do not.
func TestActivityResourceAnyPartition(t *testing.T) {
	for _, arn := range []string{
		"arn:aws:states:us-east-1:000000000000:activity:Approve",
		"arn:aws-cn:states:cn-north-1:000000000000:activity:Approve",
		"arn:aws-us-gov:states:us-gov-west-1:000000000000:activity:Approve",
	} {
		if !isActivityResource(arn) {
			t.Errorf("isActivityResource(%s) = false, want true", arn)
		}
		if got := taskResourceType(arn); got != "activity" {
			t.Errorf("taskResourceType(%s) = %s, want activity", arn, got)
		}
	}
	if isActivityResource("arn:aws:states:::sqs:sendMessage") {
		t.Error("an sqs integration resource is not an activity")
	}
	if isActivityResource("arn:aws:lambda:us-east-1:000000000000:function:f") {
		t.Error("a lambda ARN is not an activity")
	}
}

// TestSQSRegionFollowsQueueURL pins the region derivation: the queue URL
// is the only member that names the queue's region — the integration
// resource is a region-less states pseudo-ARN — so a queue in another
// region is addressed, not silently run against the default region.
func TestSQSRegionFollowsQueueURL(t *testing.T) {
	if got := regionFromQueueURL("https://sqs.eu-west-2.amazonaws.com/000000000000/q"); got != "eu-west-2" {
		t.Errorf("regionFromQueueURL = %q, want eu-west-2", got)
	}
	if got := regionFromQueueURL("not a url"); got != "" {
		t.Errorf("regionFromQueueURL(garbage) = %q, want empty", got)
	}

	sqs := &stubSQSInvoker{}
	e := newSQSTestExecutor(t, sqs)
	e.region = "us-east-1"

	input := `{"QueueUrl":"https://sqs.eu-west-2.amazonaws.com/000000000000/q","MessageBody":"m"}`
	if _, err := e.executeSQSTask(context.Background(), "arn:aws:states:::sqs:sendMessage", input, nil); err != nil {
		t.Fatalf("sendMessage failed: %v", err)
	}
	if got := sqs.sentRegions[len(sqs.sentRegions)-1]; got != "eu-west-2" {
		t.Errorf("sendMessage region = %q, want the queue URL's eu-west-2", got)
	}
}

// TestBaseStateMachineArnRebuild pins the version-ARN strip: the
// unqualified state machine ARN is rebuilt by the ARN builder with the
// parsed fields.
func TestBaseStateMachineArnRebuild(t *testing.T) {
	base, err := baseStateMachineArn("arn:aws:states:us-east-1:000000000000:stateMachine:sm:3")
	if err != nil {
		t.Fatalf("strip failed: %v", err)
	}
	if base != "arn:aws:states:us-east-1:000000000000:stateMachine:sm" {
		t.Errorf("base ARN = %s", base)
	}
	if _, err := baseStateMachineArn("arn:aws:states:us-east-1:000000000000:stateMachine:sm"); err == nil {
		t.Error("an unqualified ARN must not parse as a version ARN")
	}
}

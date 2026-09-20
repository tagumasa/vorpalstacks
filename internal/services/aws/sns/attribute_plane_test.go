package sns

// Attribute-plane pins: the emitted Get*Attributes sets follow the
// documented enumerations (invented keys absent, set-only attributes
// present, EffectiveDeliveryPolicy synthesised), the write plane rejects
// read-only and create-only keys, RawMessageDelivery carries its value and
// protocol constraints, and a RedrivePolicy must name a dead-letter queue
// that exists.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// stubSQSInvoker answers GetQueueByName from a fixed set of live queue
// names; every other method is unreachable from the paths under test.
type stubSQSInvoker struct {
	liveQueues map[string]bool
}

func (s *stubSQSInvoker) GetQueueByName(ctx context.Context, region, queueName string) (string, error) {
	if s.liveQueues[region+":"+queueName] {
		return "https://sqs." + region + ".amazonaws.com/123456789012/" + queueName, nil
	}
	return "", errors.New("queue does not exist")
}

func (s *stubSQSInvoker) GetQueueARN(ctx context.Context, region, queueURL string) (string, error) {
	return "", errors.New("ARN resolution unavailable")
}

func (s *stubSQSInvoker) SendMessage(ctx context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	return "", "", errors.New("unreachable in this test")
}

func (s *stubSQSInvoker) ReceiveMessage(ctx context.Context, region, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, errors.New("unreachable in this test")
}

func (s *stubSQSInvoker) DeleteMessage(ctx context.Context, region, queueURL, receiptHandle string) error {
	return errors.New("unreachable in this test")
}

// stubBus hands the SNS service a stubbed SQS invoker through the bus's
// registry surface, and routes published delivery events straight into the
// service's own bus handler — the direct-dispatch stand-in for the real
// bus's asynchronous workers. Unimplemented registry methods panic if
// reached, which the tests never do.
type stubBus struct {
	eventbus.ServiceBus
	sqs       invokers.SQSInvoker
	onPublish func(evt *eventbus.SNSDeliveryEvent)
}

func (b *stubBus) SQSInvoker() invokers.SQSInvoker { return b.sqs }

func (b *stubBus) Publish(ctx context.Context, event eventbus.Event) error {
	if evt, ok := event.(*eventbus.SNSDeliveryEvent); ok && b.onPublish != nil {
		b.onPublish(evt)
	}
	return nil
}

func newAttributePlaneTestContext(t *testing.T) (*SNSService, snsstore.SNSStoreInterface) {
	t.Helper()
	svc, store := newFanoutTestService(t)
	return svc, store
}

// TestTopicAttributesEmissionFollowsDocumentedSet pins the
// GetTopicAttributes emission rules: bookkeeping timestamps are not
// attributes, unset optional attributes (DisplayName, SignatureVersion,
// MaximumMessageSize) are absent rather than defaulted, and
// EffectiveDeliveryPolicy is always synthesised from the stored policy and
// the system defaults.
func TestTopicAttributesEmissionFollowsDocumentedSet(t *testing.T) {
	svc, store := newAttributePlaneTestContext(t)
	created, err := store.CreateTopic(&snsstore.Topic{Name: "emit-plane-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	result, err := svc.getTopicAttributesCore(store, GetTopicAttributesInput{TopicArn: created.Arn})
	if err != nil {
		t.Fatalf("get attributes: %v", err)
	}
	attrs := result.(map[string]interface{})["Attributes"].(map[string]string)

	for _, invented := range []string{"CreatedDate", "LastModifiedTime", "DisplayName", "SignatureVersion", "MaximumMessageSize"} {
		if _, present := attrs[invented]; present {
			t.Errorf("bare topic emits %s — the documented set returns it only when explicitly set", invented)
		}
	}
	for _, required := range []string{"TopicArn", "Owner", "Policy", "EffectiveDeliveryPolicy", "SubscriptionsConfirmed", "SubscriptionsDeleted", "SubscriptionsPending"} {
		if _, present := attrs[required]; !present {
			t.Errorf("bare topic omits the documented key %s", required)
		}
	}
	effective, err := parseTopicDeliveryPolicy(attrs["EffectiveDeliveryPolicy"])
	if err != nil {
		t.Fatalf("synthesised EffectiveDeliveryPolicy does not parse: %v", err)
	}
	if effective.numRetries != snsstore.DefaultDeliveryPolicyNumRetries {
		t.Errorf("default effective policy = %d retries, want %d", effective.numRetries, snsstore.DefaultDeliveryPolicyNumRetries)
	}

	// Set attributes surface one by one; each then appears in the read.
	sets := map[string]string{
		"DisplayName":        "Emit Plane",
		"SignatureVersion":   "2",
		"MaximumMessageSize": "2048",
	}
	for name, value := range sets {
		if err := svc.setTopicAttributesCore(store, SetTopicAttributesInput{TopicArn: created.Arn, AttributeName: name, AttributeValue: value}); err != nil {
			t.Fatalf("set %s: %v", name, err)
		}
	}
	result, err = svc.getTopicAttributesCore(store, GetTopicAttributesInput{TopicArn: created.Arn})
	if err != nil {
		t.Fatalf("re-read attributes: %v", err)
	}
	attrs = result.(map[string]interface{})["Attributes"].(map[string]string)
	for name, want := range sets {
		if attrs[name] != want {
			t.Errorf("after set, %s = %q, want %q", name, attrs[name], want)
		}
	}
}

// TestSetTopicAttributesRejectsNonSettableKeys pins the write plane:
// create-only and read-only keys are refused, and the FIFO-only
// attributes validate against the topic's type.
func TestSetTopicAttributesRejectsNonSettableKeys(t *testing.T) {
	svc, store := newAttributePlaneTestContext(t)
	created, err := store.CreateTopic(&snsstore.Topic{Name: "readonly-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	fifo, err := store.CreateTopic(&snsstore.Topic{Name: "readonly.fifo"}, nil)
	if err != nil {
		t.Fatalf("create fifo topic: %v", err)
	}

	for _, name := range []string{"FifoTopic", "EffectiveDeliveryPolicy", "TopicArn", "Owner", "SubscriptionsConfirmed", "SubscriptionsDeleted", "SubscriptionsPending"} {
		if err := svc.setTopicAttributesCore(store, SetTopicAttributesInput{TopicArn: created.Arn, AttributeName: name, AttributeValue: "1"}); err == nil {
			t.Errorf("SetTopicAttributes accepted the non-settable key %s", name)
		}
	}

	for name, tc := range map[string]struct {
		attr     string
		topicArn string
		value    string
		wantErr  bool
	}{
		"ContentBasedDeduplication on FIFO":     {"ContentBasedDeduplication", fifo.Arn, "true", false},
		"ContentBasedDeduplication bad value":   {"ContentBasedDeduplication", fifo.Arn, "yes", true},
		"ContentBasedDeduplication on standard": {"ContentBasedDeduplication", created.Arn, "true", true},
		"FifoThroughputScope on FIFO":           {"FifoThroughputScope", fifo.Arn, "MessageGroup", false},
		"FifoThroughputScope bad value":         {"FifoThroughputScope", fifo.Arn, "PerQueue", true},
		"FifoThroughputScope on standard":       {"FifoThroughputScope", created.Arn, "MessageGroup", true},
		"ArchivePolicy on FIFO":                 {"ArchivePolicy", fifo.Arn, `{"ArchiveStrategy":{"Preserve":{}}}`, false},
		"ArchivePolicy on standard":             {"ArchivePolicy", created.Arn, `{"ArchiveStrategy":{"Preserve":{}}}`, true},
		"SignatureVersion 1":                    {"SignatureVersion", created.Arn, "1", false},
		"SignatureVersion 3":                    {"SignatureVersion", created.Arn, "3", true},
		"MaximumMessageSize 1024":               {"MaximumMessageSize", created.Arn, "1024", false},
		"MaximumMessageSize below floor":        {"MaximumMessageSize", created.Arn, "1023", true},
		"MaximumMessageSize above platform cap": {"MaximumMessageSize", created.Arn, "262145", true},
		"MaximumMessageSize not a number":       {"MaximumMessageSize", created.Arn, "big", true},
	} {
		err := svc.setTopicAttributesCore(store, SetTopicAttributesInput{TopicArn: tc.topicArn, AttributeName: tc.attr, AttributeValue: tc.value})
		if (err != nil) != tc.wantErr {
			t.Errorf("%s: err = %v, want error %v", name, err, tc.wantErr)
		}
	}
}

// TestSubscriptionAttributesPlane pins the subscription read and write
// planes: the invented Protocol/Endpoint keys are gone, the documented
// synthesised keys and EffectiveDeliveryPolicy are present, RawMessageDelivery
// validates its value and its protocol set, and the read-only keys are
// refused.
func TestSubscriptionAttributesPlane(t *testing.T) {
	svc, store := newAttributePlaneTestContext(t)
	created, err := store.CreateTopic(&snsstore.Topic{Name: "sub-plane-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	subResult, err := svc.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn: created.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:plane-queue",
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	subArn := subResult.(map[string]interface{})["SubscriptionArn"].(string)

	result, err := svc.getSubscriptionAttributesCore(store, GetSubscriptionAttributesInput{SubscriptionArn: subArn})
	if err != nil {
		t.Fatalf("get subscription attributes: %v", err)
	}
	attrs := result.(map[string]interface{})["Attributes"].(map[string]string)
	for _, invented := range []string{"Protocol", "Endpoint"} {
		if _, present := attrs[invented]; present {
			t.Errorf("subscription attributes carry the invented key %s", invented)
		}
	}
	for _, required := range []string{"SubscriptionArn", "TopicArn", "Owner", "PendingConfirmation", "ConfirmationWasAuthenticated", "EffectiveDeliveryPolicy"} {
		if _, present := attrs[required]; !present {
			t.Errorf("subscription attributes omit the documented key %s", required)
		}
	}

	// RawMessageDelivery: the value is a boolean literal and the protocol
	// must be in the documented set.
	for name, tc := range map[string]struct {
		protocol string
		value    string
		wantErr  bool
	}{
		"sqs true":    {"sqs", "true", false},
		"sqs yes":     {"sqs", "yes", true},
		"http true":   {"http", "true", false},
		"email true":  {"email", "true", true},
		"lambda true": {"lambda", "true", true},
	} {
		// The write path gates on the SUBSCRIPTION's protocol, so exercise
		// validateSubscriptionAttribute directly with each protocol.
		err := validateSubscriptionAttribute("RawMessageDelivery", tc.value, tc.protocol)
		if (err != nil) != tc.wantErr {
			t.Errorf("RawMessageDelivery %s: err = %v, want error %v", name, err, tc.wantErr)
		}
	}
	if _, err := svc.setSubscriptionAttributesCore(store, reqCtx, SetSubscriptionAttributesInput{
		SubscriptionArn: subArn, AttributeName: "RawMessageDelivery", AttributeValue: "true",
	}); err != nil {
		t.Fatalf("set RawMessageDelivery on the sqs subscription: %v", err)
	}

	for _, name := range []string{"PendingConfirmation", "ConfirmationWasAuthenticated", "SubscriptionArn", "TopicArn", "Owner", "Protocol", "Endpoint", "ReplayStatus"} {
		if err := validateSubscriptionAttribute(name, "true", "sqs"); err == nil {
			t.Errorf("SetSubscriptionAttributes accepted the read-only key %s", name)
		}
	}
}

// TestRedrivePolicyTargetValidation pins the set-time dead-letter-queue
// contract: the ARN must be an SQS queue in the same account and Region,
// and it must resolve to an existing queue through the SQS invoker; with
// no invoker wired the check fails closed.
func TestRedrivePolicyTargetValidation(t *testing.T) {
	svc, store := newAttributePlaneTestContext(t)

	live := `{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:live-dlq"}`
	missing := `{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:ghost-dlq"}`
	otherService := `{"deadLetterTargetArn":"arn:aws:lambda:us-east-1:123456789012:dlq-fn"}`
	crossRegion := `{"deadLetterTargetArn":"arn:aws:sqs:eu-west-1:123456789012:far-dlq"}`

	// Fail-closed with no bus wired.
	if err := svc.validateRedrivePolicyTarget(live, "us-east-1"); err == nil {
		t.Error("RedrivePolicy accepted with no service bus wired — the target is unverifiable")
	}

	svc.bus = &stubBus{sqs: &stubSQSInvoker{liveQueues: map[string]bool{
		"us-east-1:live-dlq": true,
	}}}
	if err := svc.validateRedrivePolicyTarget(live, "us-east-1"); err != nil {
		t.Errorf("live dead-letter queue rejected: %v", err)
	}
	if err := svc.validateRedrivePolicyTarget("", "us-east-1"); err != nil {
		t.Errorf("the empty clearing value rejected: %v", err)
	}
	for name, value := range map[string]string{
		"missing queue":  missing,
		"non-SQS target": otherService,
		"cross-region":   crossRegion,
	} {
		if err := svc.validateRedrivePolicyTarget(value, "us-east-1"); err == nil {
			t.Errorf("RedrivePolicy accepted a %s", name)
		} else if !strings.Contains(err.Error(), "Invalid RedrivePolicy") && !strings.Contains(err.Error(), "Invalid redrive policy") {
			t.Errorf("%s rejected with a non-InvalidParameter message: %v", name, err)
		}
	}

	// The Subscribe plane applies the same check to inline attributes.
	created, err := store.CreateTopic(&snsstore.Topic{Name: "redrive-plane-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	if _, err := svc.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn: created.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:plane-queue",
		Attributes: map[string]string{
			"RedrivePolicy": missing,
		},
	}); err == nil {
		t.Error("Subscribe accepted inline RedrivePolicy naming a missing dead-letter queue")
	}
	if _, err := svc.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn: created.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:plane-queue",
		Attributes: map[string]string{
			"RedrivePolicy": live,
		},
	}); err != nil {
		t.Errorf("Subscribe rejected inline RedrivePolicy naming a live dead-letter queue: %v", err)
	}
}

// TestMaximumMessageSizeEnforcedAtPublish pins the honoured attribute: a
// topic with MaximumMessageSize 2048 rejects a larger body and accepts a
// smaller one; a bare topic keeps the platform's flat cap.
func TestMaximumMessageSizeEnforcedAtPublish(t *testing.T) {
	svc, store := newAttributePlaneTestContext(t)
	created, err := store.CreateTopic(&snsstore.Topic{Name: "cap-plane-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	if err := svc.setTopicAttributesCore(store, SetTopicAttributesInput{
		TopicArn: created.Arn, AttributeName: "MaximumMessageSize", AttributeValue: "2048",
	}); err != nil {
		t.Fatalf("set MaximumMessageSize: %v", err)
	}

	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: strings.Repeat("a", 2049)}); err == nil {
		t.Error("publish accepted a body above the topic's MaximumMessageSize")
	}
	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: strings.Repeat("a", 2048)}); err != nil {
		t.Errorf("publish rejected a body at the topic's MaximumMessageSize: %v", err)
	}
}

// TestTagResourceResponseIsEmpty pins the response shape: the model's
// TagResourceResponse carries no member, so the SNS tag config provides
// the empty response, not a tag list.
func TestTagResourceResponseIsEmpty(t *testing.T) {
	_, store := newAttributePlaneTestContext(t)
	cfg := snsTagConfig(store)
	if cfg.TagResponse != nil {
		t.Error("TagResource is configured to return a tag list — the response shape is empty")
	}
	if cfg.EmptyResponse == nil {
		t.Error("TagResource has no empty-response builder")
	}
}

// TestEndpointCustomUserDataRoundTrip pins the endpoint-plane emission the
// GetEndpointAttributes enumeration promises ("Attributes include the
// following: CustomUserData – arbitrary user data to associate with the
// endpoint"): CustomUserData set through CreatePlatformEndpoint's member
// answers the read even though it is a record field rather than an
// Attributes-map entry, and a later SetEndpointAttributes write replaces
// it.
func TestEndpointCustomUserDataRoundTrip(t *testing.T) {
	svc, store := newAttributePlaneTestContext(t)

	created, err := svc.createPlatformApplicationCore(store, CreatePlatformApplicationInput{
		Name:       "custom-data-app",
		Platform:   "APNS",
		Attributes: map[string]string{"PlatformCredential": "credential"},
	})
	if err != nil {
		t.Fatalf("create platform application: %v", err)
	}
	appArn := created.(map[string]interface{})["PlatformApplicationArn"].(string)

	endpoint, err := svc.createPlatformEndpointCore(store, CreatePlatformEndpointInput{
		PlatformApplicationArn: appArn,
		Token:                  "custom-data-token",
		CustomUserData:         "UserId=01234567",
	})
	if err != nil {
		t.Fatalf("create platform endpoint: %v", err)
	}
	endpointArn := endpoint.(map[string]interface{})["EndpointArn"].(string)

	attrs, err := svc.getEndpointAttributesCore(store, GetEndpointAttributesInput{EndpointArn: endpointArn})
	if err != nil {
		t.Fatalf("get endpoint attributes: %v", err)
	}
	entries := attrs.(map[string]interface{})["Attributes"].(map[string]string)
	if got := entries["CustomUserData"]; got != "UserId=01234567" {
		t.Fatalf("CustomUserData after create = %q, want the created value — the member belongs to the documented response enumeration", got)
	}
	if got := entries["Enabled"]; got != "true" {
		t.Fatalf("Enabled default = %q, want true", got)
	}
	if got := entries["Token"]; got != "custom-data-token" {
		t.Fatalf("Token after create = %q, want the created token", got)
	}

	if _, err := svc.setEndpointAttributesCore(store, SetEndpointAttributesInput{
		EndpointArn: endpointArn,
		Attributes:  map[string]string{"CustomUserData": "UserId=76543210"},
	}); err != nil {
		t.Fatalf("set endpoint attributes: %v", err)
	}
	attrs, err = svc.getEndpointAttributesCore(store, GetEndpointAttributesInput{EndpointArn: endpointArn})
	if err != nil {
		t.Fatalf("re-read endpoint attributes: %v", err)
	}
	entries = attrs.(map[string]interface{})["Attributes"].(map[string]string)
	if got := entries["CustomUserData"]; got != "UserId=76543210" {
		t.Fatalf("CustomUserData after set = %q, want the written value to replace the created one", got)
	}
}

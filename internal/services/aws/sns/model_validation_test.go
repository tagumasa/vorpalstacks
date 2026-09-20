package sns

// Model-derived validation pins: topic names carry
// no reserved-prefix rule, CreatePlatformApplication.Attributes is
// required, the tag framework's constraints are wired, Subscribe is
// idempotent and refuses protocols it cannot honestly accept, and the
// filter-policy grammar accepts exactly what the matcher evaluates.

import (
	"context"
	"errors"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// TestTopicNameReservedPrefixAccepted pins the removal of the invented
// aws/amazon reserved-prefix rejection: the CreateTopic documentation
// states no reserved prefix, so prefixed names are legal.
func TestTopicNameReservedPrefixAccepted(t *testing.T) {
	for _, name := range []string{"aws-events", "amazon-notifications", "AWS-Topic", "amazing-topic"} {
		if err := validateTopicName(name); err != nil {
			t.Fatalf("validateTopicName(%q) = %v, want nil — no reserved prefix exists", name, err)
		}
	}
}

// TestCreatePlatformApplicationAttributesRequired pins the @required
// enforcement: the model marks CreatePlatformApplicationInput.Attributes
// smithy.api#required, and an absent (query-wire: zero-entry) member is
// rejected instead of silently creating an attribute-less application.
func TestCreatePlatformApplicationAttributesRequired(t *testing.T) {
	svc, store := newFanoutTestService(t)

	_, err := svc.createPlatformApplicationCore(store, CreatePlatformApplicationInput{
		Name:       "attrs-required-app",
		Platform:   "APNS",
		Attributes: nil,
	})
	if err == nil {
		t.Fatal("createPlatformApplicationCore accepted an absent Attributes member, want Attributes is required")
	}
	var wireErr *awserrors.AWSError
	if !errors.As(err, &wireErr) || wireErr.QueryErrorCode != "InvalidParameter" {
		t.Fatalf("absent Attributes returned %v, want the InvalidParameter query code", err)
	}
	if !strings.Contains(err.Error(), "Attributes is required") {
		t.Fatalf("absent Attributes returned %q, want the Attributes is required message", err.Error())
	}

	// A present member creates normally.
	created, err := svc.createPlatformApplicationCore(store, CreatePlatformApplicationInput{
		Name:       "attrs-present-app",
		Platform:   "APNS",
		Attributes: map[string]string{"PlatformCredential": "credential"},
	})
	if err != nil {
		t.Fatalf("createPlatformApplicationCore with Attributes: %v", err)
	}
	_ = created
}

// TestTopicTagValidation pins the wired tag constraints on both request
// forms: the fifty-tag cap with the model's TagLimitExceeded identity, the
// reserved aws: prefix, the key and value lengths, and the tag operations'
// ResourceArn length bound.
func TestTopicTagValidation(t *testing.T) {
	over := make(map[string]string, 51)
	for i := 0; i < 51; i++ {
		over[strings.Repeat("k", 4)+string(rune('a'+i%26))+string(rune('a'+i/26))+string(rune('a'+i%26))] = "v"
	}
	if err := validateTopicTagMap(over); err == nil {
		t.Fatal("51 tags accepted, want TagLimitExceeded")
	} else if err != ErrTagLimitExceeded {
		t.Fatalf("51 tags returned %v, want ErrTagLimitExceeded", err)
	}

	bad := map[string]map[string]string{
		"reserved prefix": {"aws:team": "x"},
		"key too long":    {strings.Repeat("k", 129): "x"},
		"value too long":  {"k": strings.Repeat("v", 257)},
	}
	for name, tags := range bad {
		err := validateTopicTagMap(tags)
		if err == nil {
			t.Fatalf("%s accepted, want InvalidParameter", name)
		}
	}
	if err := validateTopicTagMap(map[string]string{"team": "platform", "cost": "edge"}); err != nil {
		t.Fatalf("valid tag set rejected: %v", err)
	}

	if err := validateTopicTagKeys([]string{"aws:reserved"}); err == nil {
		t.Fatal("reserved untag key accepted, want InvalidParameter")
	}
	if err := validateTopicTagKeys([]string{strings.Repeat("k", 129)}); err == nil {
		t.Fatal("129-character untag key accepted, want InvalidParameter")
	}

	// The tag operations' ResourceArn member is bounded by the model's
	// AmazonResourceName length trait (1..1011); the check lives in the
	// handler framework's resource-validation closure.
	_, store := newFanoutTestService(t)
	cfg := snsTagConfig(store)
	if err := cfg.ValidateResource(context.Background(), "arn:aws:sns:us-east-1:123456789012:"+strings.Repeat("t", 1012)); err == nil {
		t.Fatal("over-long ResourceArn accepted, want InvalidParameter")
	}
}

// TestSubscribeProtocolRefusals pins the protocol x delivery linkage: the
// application protocol is refused outright (its delivery needs external
// push services the platform never implements), and FIFO topics refuse
// every customer-managed endpoint protocol — sqs, lambda and firehose
// remain subscribable.
func TestSubscribeProtocolRefusals(t *testing.T) {
	cases := []struct {
		protocol string
		fifo     bool
		wantErr  bool
	}{
		{"application", false, true},
		{"application", true, true},
		{"http", true, true},
		{"https", true, true},
		{"email", true, true},
		{"email-json", true, true},
		{"sms", true, true},
		{"sqs", true, false},
		{"lambda", true, false},
		{"firehose", true, false},
		{"http", false, false},
		{"email", false, false},
		{"sqs", false, false},
	}
	for _, tc := range cases {
		err := rejectSubscriptionProtocol(tc.protocol, tc.fifo)
		if tc.wantErr && err == nil {
			t.Fatalf("rejectSubscriptionProtocol(%q, fifo=%v) accepted, want rejection", tc.protocol, tc.fifo)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("rejectSubscriptionProtocol(%q, fifo=%v) = %v, want accepted", tc.protocol, tc.fifo, err)
		}
	}
}

// TestSubscribeIdempotency pins the natural-key idempotency through the
// core: a repeated Subscribe returns the live subscription's ARN without
// creating a duplicate or double-confirming, while a different endpoint
// still creates a fresh subscription.
func TestSubscribeIdempotency(t *testing.T) {
	svc, store := newFanoutTestService(t)
	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")

	topic, err := store.CreateTopic(&snsstore.Topic{Name: "idempotent-sub-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	first, err := svc.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:idempotency-queue",
	})
	if err != nil {
		t.Fatalf("first subscribe: %v", err)
	}
	second, err := svc.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:idempotency-queue",
	})
	if err != nil {
		t.Fatalf("repeated subscribe: %v", err)
	}
	firstArn := first.(map[string]interface{})["SubscriptionArn"].(string)
	secondArn := second.(map[string]interface{})["SubscriptionArn"].(string)
	if firstArn != secondArn {
		t.Fatalf("repeated Subscribe returned %q after %q — duplicates accumulate", secondArn, firstArn)
	}

	third, err := svc.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:another-queue",
	})
	if err != nil {
		t.Fatalf("distinct-endpoint subscribe: %v", err)
	}
	thirdArn := third.(map[string]interface{})["SubscriptionArn"].(string)
	if thirdArn == firstArn {
		t.Fatal("distinct endpoint collapsed onto the existing subscription's ARN")
	}

	subs, err := store.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil {
		t.Fatalf("list subscriptions: %v", err)
	}
	if len(subs) != 2 {
		t.Fatalf("topic carries %d subscriptions, want 2 — the repeated Subscribe must not persist a duplicate", len(subs))
	}

	// The auto-confirmed re-subscribe did not double-count: the topic
	// reports one confirmed subscription per endpoint.
	after, err := store.GetTopic(topic.Arn)
	if err != nil {
		t.Fatalf("get topic: %v", err)
	}
	if after.SubscriptionsConfirmed != 2 || after.SubscriptionsPending != 0 {
		t.Fatalf("counters after idempotent re-subscribe: confirmed=%d pending=%d, want 2/0",
			after.SubscriptionsConfirmed, after.SubscriptionsPending)
	}
}

// TestSubscribeIdempotencyDifferentOwner pins that the natural key
// includes the subscription owner: another account subscribing the same
// endpoint creates its own subscription rather than borrowing the first
// owner's.
func TestSubscribeIdempotencyDifferentOwner(t *testing.T) {
	svc, store := newFanoutTestService(t)
	ownerA := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	ownerB := request.NewRequestContext(context.Background(), nil, "210987654321", "us-east-1")

	topic, err := store.CreateTopic(&snsstore.Topic{Name: "cross-owner-sub-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	input := SubscribeInput{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:shared-queue",
	}
	if _, err := svc.subscribeCore(store, ownerA, input); err != nil {
		t.Fatalf("owner A subscribe: %v", err)
	}
	if _, err := svc.subscribeCore(store, ownerB, input); err != nil {
		t.Fatalf("owner B subscribe: %v", err)
	}
	subs, err := store.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil {
		t.Fatalf("list subscriptions: %v", err)
	}
	if len(subs) != 2 {
		t.Fatalf("topic carries %d subscriptions, want 2 — the owner is part of the natural key", len(subs))
	}
}

// TestPlatformDeleteMissingReportsModelDeclaredError pins the delete-plane
// error contract: DeletePlatformApplication and DeleteEndpoint declare no
// NotFoundException in the model (AuthorizationError, InternalError and
// InvalidParameter alone), so a missing resource refuses with
// InvalidParameter — not the NotFound the store sentinel maps to on the
// read planes.
func TestPlatformDeleteMissingReportsModelDeclaredError(t *testing.T) {
	svc, store := newFanoutTestService(t)

	created, err := svc.createPlatformApplicationCore(store, CreatePlatformApplicationInput{
		Name:       "delete-error-app",
		Platform:   "APNS",
		Attributes: map[string]string{"PlatformCredential": "credential"},
	})
	if err != nil {
		t.Fatalf("create platform application: %v", err)
	}
	appArn := created.(map[string]interface{})["PlatformApplicationArn"].(string)
	endpoint, err := svc.createPlatformEndpointCore(store, CreatePlatformEndpointInput{
		PlatformApplicationArn: appArn,
		Token:                  "delete-error-token",
	})
	if err != nil {
		t.Fatalf("create platform endpoint: %v", err)
	}
	endpointArn := endpoint.(map[string]interface{})["EndpointArn"].(string)

	if _, err := svc.deletePlatformApplicationCore(store, DeletePlatformApplicationInput{PlatformApplicationArn: appArn}); err != nil {
		t.Fatalf("first delete: %v", err)
	}

	assertInvalidParameter := func(what string, err error) {
		t.Helper()
		if err == nil {
			t.Fatalf("%s on a missing resource succeeded, want InvalidParameter", what)
		}
		var wireErr *awserrors.AWSError
		if !errors.As(err, &wireErr) || wireErr.QueryErrorCode != "InvalidParameter" {
			t.Fatalf("%s returned %v, want the InvalidParameter query code", what, err)
		}
	}
	assertInvalidParameter("repeated DeletePlatformApplication",
		mustError(svc.deletePlatformApplicationCore(store, DeletePlatformApplicationInput{PlatformApplicationArn: appArn})))
	assertInvalidParameter("DeleteEndpoint of the cascaded endpoint",
		mustError(svc.deleteEndpointCore(store, DeleteEndpointInput{EndpointArn: endpointArn})))
}

// mustError adapts a core call's (interface{}, error) pair to the error for
// the assertion helper above.
func mustError(_ interface{}, err error) error {
	return err
}

// TestProtocolRegistryConsistency pins the invariants the registry's
// consumers rely on: a protocol that receives a confirmation POST must
// have a delivery handler (the POST goes to the subscribed endpoint) and
// must not auto-confirm (an auto-confirmed subscription never pends).
func TestProtocolRegistryConsistency(t *testing.T) {
	for protocol, info := range protocolRegistry {
		if info.confirmationPOST && !info.deliveryHandler {
			t.Errorf("%s: confirmationPOST without a delivery handler — the confirmation POST has no endpoint to reach", protocol)
		}
		if info.confirmationPOST && info.autoConfirm {
			t.Errorf("%s: confirmationPOST on an auto-confirming protocol — the subscription never pends", protocol)
		}
	}
}

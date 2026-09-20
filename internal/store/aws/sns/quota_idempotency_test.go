package sns

// Store pins for the model-derived bounds:
// subscription creation is idempotent on its natural key, the documented
// per-topic quotas (FIFO 100 subscriptions, 200 filter policies) hold
// under the create/attach paths, the per-account FIFO topic quota caps
// CreateTopic, and AddPermission refuses a duplicate label instead of
// overwriting the existing statement.

import (
	"errors"
	"fmt"
	"testing"
)

// TestCreateSubscriptionIdempotentNaturalKey pins the idempotency at the
// store seam: the same (topic, protocol, endpoint, owner) returns the live
// record with created=false and no duplicate persists; a differing key
// creates a fresh record.
func TestCreateSubscriptionIdempotentNaturalKey(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "idempotency-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	mk := func(endpoint string) *Subscription {
		return &Subscription{
			TopicArn: topic.Arn,
			Protocol: "sqs",
			Endpoint: endpoint,
			Owner:    "123456789012",
		}
	}

	first, created, err := store.CreateSubscription(mk("arn:aws:sqs:us-east-1:123456789012:q"))
	if err != nil || !created {
		t.Fatalf("first create: created=%v err=%v", created, err)
	}
	again, created, err := store.CreateSubscription(mk("arn:aws:sqs:us-east-1:123456789012:q"))
	if err != nil {
		t.Fatalf("repeated create: %v", err)
	}
	if created {
		t.Fatal("repeated create reported created=true — duplicates accumulate")
	}
	if again.SubscriptionArn != first.SubscriptionArn {
		t.Fatalf("repeated create returned %s, want the live record %s", again.SubscriptionArn, first.SubscriptionArn)
	}

	other, created, err := store.CreateSubscription(mk("arn:aws:sqs:us-east-1:123456789012:other"))
	if err != nil || !created {
		t.Fatalf("distinct-endpoint create: created=%v err=%v", created, err)
	}
	if other.SubscriptionArn == first.SubscriptionArn {
		t.Fatal("distinct endpoint collapsed onto the existing subscription")
	}

	// A different owner is a different natural key even for the same
	// endpoint.
	third, created, err := store.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:q",
		Owner:    "210987654321",
	})
	if err != nil || !created {
		t.Fatalf("cross-owner create: created=%v err=%v", created, err)
	}
	if third.SubscriptionArn == first.SubscriptionArn {
		t.Fatal("cross-owner create borrowed the other owner's subscription")
	}

	subs, err := store.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(subs) != 3 {
		t.Fatalf("topic carries %d subscriptions, want 3", len(subs))
	}

	// The pending counter moved once per genuinely new subscription.
	after, err := store.GetTopic(topic.Arn)
	if err != nil {
		t.Fatalf("get topic: %v", err)
	}
	if after.SubscriptionsPending != 3 {
		t.Fatalf("pending counter = %d, want 3 (the idempotent repeat must not count)", after.SubscriptionsPending)
	}
}

// TestFifoSubscriptionQuota pins the documented FIFO bound: "FIFO: 100 per
// topic" — the 101st subscription is refused SubscriptionLimitExceeded,
// and a standard topic keeps admitting.
func TestFifoSubscriptionQuota(t *testing.T) {
	store := newTestSNSStore(t)

	fifoTopic, err := store.CreateTopic(&Topic{Name: "fifo-quota-topic.fifo"}, nil)
	if err != nil {
		t.Fatalf("create fifo topic: %v", err)
	}
	for i := 0; i < MaxFifoSubscriptionsPerTopic; i++ {
		_, _, err := store.CreateSubscription(&Subscription{
			TopicArn: fifoTopic.Arn,
			Protocol: "sqs",
			Endpoint: fmt.Sprintf("arn:aws:sqs:us-east-1:123456789012:q-%d.fifo", i),
			Owner:    "123456789012",
		})
		if err != nil {
			t.Fatalf("subscription %d under the quota: %v", i, err)
		}
	}
	_, _, err = store.CreateSubscription(&Subscription{
		TopicArn: fifoTopic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:over-quota.fifo",
		Owner:    "123456789012",
	})
	if !errors.Is(err, ErrSubscriptionLimitExceeded) {
		t.Fatalf("subscription past the FIFO quota returned %v, want ErrSubscriptionLimitExceeded", err)
	}

	standardTopic, err := store.CreateTopic(&Topic{Name: "standard-quota-topic"}, nil)
	if err != nil {
		t.Fatalf("create standard topic: %v", err)
	}
	if _, _, err := store.CreateSubscription(&Subscription{
		TopicArn: standardTopic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:q",
		Owner:    "123456789012",
	}); err != nil {
		t.Fatalf("standard-topic subscription: %v", err)
	}
}

// TestFilterPolicyQuota pins the documented "200 filter policies per
// topic" on both attach paths: the 201st subscription carrying a filter
// policy is refused at Subscribe, a filterless subscription still joins,
// and SetSubscriptionAttributes refuses to attach the 201st policy while
// updating an existing one stays free.
func TestFilterPolicyQuota(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "filter-quota-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	const withPolicy = `{"k": ["v"]}`
	var policiedArn string
	for i := 0; i < MaxFilterPoliciesPerTopic; i++ {
		policied, _, err := store.CreateSubscription(&Subscription{
			TopicArn: topic.Arn,
			Protocol: "sqs",
			Endpoint: fmt.Sprintf("arn:aws:sqs:us-east-1:123456789012:q-%d", i),
			Owner:    "123456789012",
			Attributes: map[string]string{
				"FilterPolicy": withPolicy,
			},
		})
		if err != nil {
			t.Fatalf("filtered subscription %d under the quota: %v", i, err)
		}
		policiedArn = policied.SubscriptionArn
	}

	_, _, err = store.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:q-over",
		Owner:    "123456789012",
		Attributes: map[string]string{
			"FilterPolicy": withPolicy,
		},
	})
	if !errors.Is(err, ErrFilterPolicyLimitExceeded) {
		t.Fatalf("filtered subscription past the quota returned %v, want ErrFilterPolicyLimitExceeded", err)
	}

	// A filterless subscription joins at the full filter budget.
	filterless, _, err := store.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:q-plain",
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("filterless subscription at full filter budget: %v", err)
	}

	// Attaching the 201st policy through SetSubscriptionAttributes is
	// refused; updating one of the 200 existing policies is not a new
	// attachment. The updated ARN comes from the creation loop — the
	// list order is random (UUID ARN suffixes in key order), so a
	// positional pick could land on the filterless subscription this
	// test just created.
	if err := store.SetSubscriptionAttributes(filterless.SubscriptionArn, map[string]string{
		"FilterPolicy": withPolicy,
	}); !errors.Is(err, ErrFilterPolicyLimitExceeded) {
		t.Fatalf("attaching the 201st filter policy returned %v, want ErrFilterPolicyLimitExceeded", err)
	}
	if err := store.SetSubscriptionAttributes(policiedArn, map[string]string{
		"FilterPolicy": `{"k": ["w"]}`,
	}); err != nil {
		t.Fatalf("updating an existing filter policy at the budget: %v", err)
	}
}

// TestFifoTopicQuota pins the documented "FIFO: 1,000 per account" topic
// quota at the create path (the standard-topic bound of 100,000 walks the
// same code path).
func TestFifoTopicQuota(t *testing.T) {
	store := newTestSNSStore(t)

	for i := 0; i < MaxFifoTopicsPerAccount; i++ {
		if _, err := store.CreateTopic(&Topic{Name: fmt.Sprintf("fifo-quota-%04d.fifo", i)}, nil); err != nil {
			t.Fatalf("FIFO topic %d under the quota: %v", i, err)
		}
	}
	if _, err := store.CreateTopic(&Topic{Name: "fifo-quota-over.fifo"}, nil); !errors.Is(err, ErrTopicLimitExceeded) {
		t.Fatalf("FIFO topic past the quota returned %v, want ErrTopicLimitExceeded", err)
	}
	// The budgets are per type: a standard topic is still admitted.
	if _, err := store.CreateTopic(&Topic{Name: "standard-after-fifo-budget"}, nil); err != nil {
		t.Fatalf("standard topic after the FIFO budget is full: %v", err)
	}
}

// TestAddPermissionDuplicateLabelRefused pins the model-derived duplicate
// rejection: Label is "A unique identifier for the new policy statement",
// so a second statement under a live label is refused instead of
// overwriting the first.
func TestAddPermissionDuplicateLabelRefused(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "permission-label-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	if err := store.AddPermission(topic.Arn, &Permission{
		Label:      "publish-label",
		Principals: []string{"210987654321"},
		Actions:    []string{"Publish"},
	}); err != nil {
		t.Fatalf("first AddPermission: %v", err)
	}
	if err := store.AddPermission(topic.Arn, &Permission{
		Label:      "publish-label",
		Principals: []string{"210987654321", "111111111111"},
		Actions:    []string{"Publish", "Receive"},
	}); !errors.Is(err, ErrPermissionLabelExists) {
		t.Fatalf("duplicate label returned %v, want ErrPermissionLabelExists", err)
	}

	after, err := store.GetTopic(topic.Arn)
	if err != nil {
		t.Fatalf("get topic: %v", err)
	}
	if len(after.Permissions) != 1 {
		t.Fatalf("topic carries %d permission statements, want 1 — the duplicate must not persist", len(after.Permissions))
	}
	if len(after.Permissions[0].Principals) != 1 || len(after.Permissions[0].Actions) != 1 {
		t.Fatalf("the refused duplicate overwrote the original statement: %+v", after.Permissions[0])
	}

	// A distinct label still appends.
	if err := store.AddPermission(topic.Arn, &Permission{
		Label:      "second-label",
		Principals: []string{"210987654321"},
		Actions:    []string{"Receive"},
	}); err != nil {
		t.Fatalf("distinct-label AddPermission: %v", err)
	}
}

// TestFirehoseSubscriptionQuota pins the protocol-specific quota: "For
// Firehose delivery streams, 5 per topic, per subscription owner" (SNS
// endpoints and quotas) — the sixth firehose subscription is refused while
// the general per-topic quota still admits other protocols.
func TestFirehoseSubscriptionQuota(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "firehose-quota-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	for i := 0; i < MaxFirehoseSubscriptionsPerTopic; i++ {
		_, _, err := store.CreateSubscription(&Subscription{
			TopicArn: topic.Arn,
			Protocol: "firehose",
			Endpoint: fmt.Sprintf("arn:aws:firehose:us-east-1:123456789012:deliverystream/stream-%d", i),
			Owner:    "123456789012",
		})
		if err != nil {
			t.Fatalf("firehose subscription %d: %v", i, err)
		}
	}
	_, _, err = store.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "firehose",
		Endpoint: "arn:aws:firehose:us-east-1:123456789012:deliverystream/stream-overflow",
		Owner:    "123456789012",
	})
	if err != ErrSubscriptionLimitExceeded {
		t.Fatalf("sixth firehose subscription error = %v, want ErrSubscriptionLimitExceeded", err)
	}
	_, _, err = store.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:other-protocol-still-admitted",
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("non-firehose subscription past the firehose quota: %v", err)
	}
}

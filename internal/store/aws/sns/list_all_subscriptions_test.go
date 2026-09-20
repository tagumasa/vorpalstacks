package sns

import (
	"fmt"
	"testing"

	"vorpalstacks/internal/store/aws/common"
)

// TestListAllSubscriptionsByTopicWalksAllPages pins the unbounded walk: a
// topic holding more subscriptions than one list page (DefaultMaxItems is
// 100) returns every subscription through ListAllSubscriptionsByTopic —
// delivery fan-out built on a single capped page would silently drop
// subscription 101+.
func TestListAllSubscriptionsByTopicWalksAllPages(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "fanout-walk-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	const want = 101
	for i := 0; i < want; i++ {
		if _, _, err := store.CreateSubscription(&Subscription{
			TopicArn: topic.Arn,
			Protocol: "sqs",
			Endpoint: fmt.Sprintf("arn:aws:sqs:us-east-1:123456789012:q-%03d", i),
			Owner:    "123456789012",
		}); err != nil {
			t.Fatalf("create subscription %d: %v", i, err)
		}
	}

	// A single page cannot hold the set — the precondition of the pin.
	paged, err := store.ListSubscriptionsByTopic(topic.Arn, common.ListOptions{})
	if err != nil {
		t.Fatalf("single-page list: %v", err)
	}
	if len(paged.Items) != 100 || !paged.IsTruncated {
		t.Fatalf("precondition: one page held %d items, truncated=%v, want 100 truncated", len(paged.Items), paged.IsTruncated)
	}

	all, err := store.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil {
		t.Fatalf("walk-all list: %v", err)
	}
	if len(all) != want {
		t.Fatalf("ListAllSubscriptionsByTopic returned %d subscriptions, want %d — a page of the walk went unread", len(all), want)
	}

	endpoints := make(map[string]bool, len(all))
	for _, sub := range all {
		endpoints[sub.Endpoint] = true
	}
	if len(endpoints) != want {
		t.Fatalf("walk returned %d distinct endpoints, want %d — a page repeated or dropped subscriptions", len(endpoints), want)
	}
}

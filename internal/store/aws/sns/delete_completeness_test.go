package sns

// Deletion completeness pins: a reported-successful delete leaves none of
// the record's companions behind — subscriptions and tags leave with the
// topic, endpoints, indexes and tags leave with the platform application,
// all inside the same transaction that removes the record itself.

import (
	"testing"

	types "vorpalstacks/internal/common/tags"
)

// TestDeleteTopicRemovesSubscriptionsAndTags pins the topic deletion's
// completeness: every subscription record and the topic's tags leave
// storage with the topic record, so a re-created topic name starts clean
// and no orphan survives the call.
func TestDeleteTopicRemovesSubscriptionsAndTags(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "delete-complete-topic"}, map[string]string{"env": "pin"})
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	if tags, err := store.List(topic.Arn); err != nil || len(tags) == 0 {
		t.Fatalf("seed tags: %v (%v)", tags, err)
	}
	for i := 0; i < 3; i++ {
		endpoint := "arn:aws:sqs:us-east-1:123456789012:q" + string(rune('a'+i))
		if _, _, err := store.CreateSubscription(&Subscription{
			TopicArn: topic.Arn,
			Protocol: "sqs",
			Endpoint: endpoint,
			Owner:    "123456789012",
		}); err != nil {
			t.Fatalf("create subscription %d: %v", i, err)
		}
	}

	if err := store.DeleteTopic(topic.Arn); err != nil {
		t.Fatalf("delete topic: %v", err)
	}

	subs, err := store.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil {
		t.Fatalf("list subscriptions after delete: %v", err)
	}
	if len(subs) != 0 {
		t.Fatalf("%d subscription records survived the topic deletion", len(subs))
	}
	if tags, err := store.List(topic.Arn); err != nil || len(tags) != 0 {
		t.Fatalf("tags survived the topic deletion: %v (%v)", tags, err)
	}
	if _, err := store.GetTopic(topic.Arn); err != ErrTopicNotFound {
		t.Fatalf("topic read after delete = %v, want the not-found sentinel", err)
	}
}

// TestDeletePlatformApplicationRemovesEndpointsIndexesAndTags pins the
// platform application deletion's completeness: endpoint records, both
// index families and the application's tags leave with the application
// record.
func TestDeletePlatformApplicationRemovesEndpointsIndexesAndTags(t *testing.T) {
	store := newTestSNSStore(t)

	app, err := store.CreatePlatformApplication(&PlatformApplication{Name: "delete-complete-app"})
	if err != nil {
		t.Fatalf("create platform application: %v", err)
	}
	if err := store.Tag(app.PlatformApplicationArn, []types.Tag{{Key: "env", Value: "pin"}}); err != nil {
		t.Fatalf("seed tags: %v", err)
	}
	token := "delete-complete-token"
	endpoint, err := store.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  token,
	})
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	// Positive control: both index families hold the endpoint's entries
	// before the delete, so their absence afterwards is the deletion's
	// work and not a seed gap.
	appIdxKey := app.PlatformApplicationArn + "\x00" + endpoint.EndpointArn
	tokenIdxKey := app.PlatformApplicationArn + "\x00" + token
	if !store.platformAppEndpointsIndex.Has([]byte(appIdxKey)) {
		t.Fatal("the app-endpoints index never held the endpoint")
	}
	if mapped, err := store.endpointTokenIndex.Get([]byte(tokenIdxKey)); err != nil || string(mapped) != endpoint.EndpointArn {
		t.Fatalf("the token index never mapped the token: %q (%v)", mapped, err)
	}

	if err := store.DeletePlatformApplication(app.PlatformApplicationArn); err != nil {
		t.Fatalf("delete platform application: %v", err)
	}

	if _, err := store.GetEndpoint(endpoint.EndpointArn); err != ErrEndpointNotFound {
		t.Fatalf("endpoint read after delete = %v, want the not-found sentinel", err)
	}
	// Both index families leave with the application. The probe reads the
	// index buckets directly: a leftover token-index entry is invisible to
	// behaviour — with its endpoint record gone, the next endpoint create
	// overwrites it — so only the bucket read can pin the invariant.
	if store.platformAppEndpointsIndex.Has([]byte(appIdxKey)) {
		t.Fatal("the app-endpoints index still holds the deleted endpoint")
	}
	if mapped, err := store.endpointTokenIndex.Get([]byte(tokenIdxKey)); err != nil || len(mapped) != 0 {
		t.Fatalf("the token index still maps the deleted endpoint: %q (%v)", mapped, err)
	}
	if _, err := store.GetPlatformApplication(app.PlatformApplicationArn); err != ErrPlatformApplicationNotFound {
		t.Fatalf("application read after delete = %v, want the not-found sentinel", err)
	}
	if tags, err := store.List(app.PlatformApplicationArn); err != nil || len(tags) != 0 {
		t.Fatalf("tags survived the application deletion: %v (%v)", tags, err)
	}
}

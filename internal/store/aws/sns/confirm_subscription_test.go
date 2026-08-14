package sns

import (
	"fmt"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newTestSNSStore creates an SNSStore backed by a temporary Pebble storage.
func newTestSNSStore(t *testing.T) *SNSStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewSNSStore(st, "123456789012", "us-east-1")
}

// TestConfirmSubscription_AuthenticateOnUnsubscribeSemantics verifies the
// two distinct concepts: ConfirmationWasAuthenticated must always be true
// for the signed ConfirmSubscription API call, while the
// AuthenticateOnUnsubscribe input flag is persisted as a separate attribute
// only when the parameter was sent.
func TestConfirmSubscription_AuthenticateOnUnsubscribeSemantics(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "semantics-topic"})
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	sub, err := store.CreateSubscription(&Subscription{
		TopicArn:            topic.Arn,
		Protocol:            "email",
		Endpoint:            "user@example.com",
		Owner:               "123456789012",
		PendingConfirmation: true,
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	token := sub.ConfirmationToken
	if token == "" {
		t.Fatal("subscription has no confirmation token")
	}

	cases := []struct {
		name string
		flag *bool
		want string // expected Attributes["AuthenticateOnUnsubscribe"], "" = absent
	}{
		{"nil flag leaves attribute absent", nil, ""},
		{"true flag persists true", boolPtr(true), "true"},
		{"false flag persists false", boolPtr(false), "false"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the subscription to pending with a fresh token.
			sub.PendingConfirmation = true
			sub.ConfirmationToken = fmt.Sprintf("token-%s", tc.name)
			sub.ConfirmationWasAuthenticated = false
			sub.Attributes = nil
			if err := store.UpdateSubscription(sub); err != nil {
				t.Fatalf("reset subscription: %v", err)
			}

			confirmed, err := store.ConfirmSubscription(sub.SubscriptionArn, sub.ConfirmationToken, tc.flag)
			if err != nil {
				t.Fatalf("confirm: %v", err)
			}
			if !confirmed.ConfirmationWasAuthenticated {
				t.Error("ConfirmationWasAuthenticated = false, want true (signed API call)")
			}
			if confirmed.PendingConfirmation {
				t.Error("PendingConfirmation = true after confirm")
			}
			got := confirmed.Attributes["AuthenticateOnUnsubscribe"]
			if got != tc.want {
				t.Errorf("AuthenticateOnUnsubscribe = %q, want %q", got, tc.want)
			}

			// GetSubscriptionAttributes must not expose the internal flag.
			attrs, err := store.GetSubscriptionAttributes(sub.SubscriptionArn)
			if err != nil {
				t.Fatalf("get attrs: %v", err)
			}
			if _, ok := attrs["AuthenticateOnUnsubscribe"]; ok {
				t.Error("GetSubscriptionAttributes exposed AuthenticateOnUnsubscribe")
			}
		})
	}
}

func boolPtr(b bool) *bool { return &b }

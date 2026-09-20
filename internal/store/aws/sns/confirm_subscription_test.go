package sns

import (
	"errors"
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

	topic, err := store.CreateTopic(&Topic{Name: "semantics-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
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
			// Each case confirms a freshly created pending subscription —
			// production has no re-pend path, so the fixture is a new
			// natural key (a distinct endpoint) carrying its own token.
			sub, _, err := store.CreateSubscription(&Subscription{
				TopicArn: topic.Arn,
				Protocol: "email",
				Endpoint: fmt.Sprintf("user+%s@example.com", tc.name),
				Owner:    "123456789012",
			})
			if err != nil {
				t.Fatalf("create subscription: %v", err)
			}
			token := sub.ConfirmationToken
			if token == "" {
				t.Fatal("subscription has no confirmation token")
			}

			confirmed, err := store.ConfirmSubscription(sub.SubscriptionArn, token, tc.flag)
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

// TestFindSubscriptionByTokenDistinguishesStorageFailure pins the token
// scan's error identity: a miss is the not-found sentinel, while a record
// the scanner cannot decode is a storage failure that must surface as
// itself — ConfirmSubscription's core answers only the miss with
// InvalidParameter and passes storage failures to the internal plane.
func TestFindSubscriptionByTokenDistinguishesStorageFailure(t *testing.T) {
	st := newTestSNSStore(t)
	topicArn := "arn:aws:sns:us-east-1:123456789012:token-scan"

	if _, err := st.FindSubscriptionByToken(topicArn, "token-x"); !errors.Is(err, ErrSubscriptionNotFound) {
		t.Fatalf("empty scan: err = %v, want ErrSubscriptionNotFound", err)
	}

	if err := st.topicSubscriptionsStore.PutRaw("corrupt-subscription-record", []byte("{not json")); err != nil {
		t.Fatalf("seed corrupt record: %v", err)
	}
	if _, err := st.FindSubscriptionByToken(topicArn, "token-x"); err == nil || errors.Is(err, ErrSubscriptionNotFound) {
		t.Fatalf("scan over a corrupt record: err = %v, want a storage failure distinct from the not-found sentinel", err)
	}
}

// TestConfirmSubscriptionTokenSurvivesConfirmation pins the token's
// documented lifetime: "Confirmation tokens are valid for two days"
// (Subscribe) — a repeated ConfirmSubscription with the same token within
// the window succeeds instead of failing InvalidParameter, so confirmation
// is idempotent.
func TestConfirmSubscriptionTokenSurvivesConfirmation(t *testing.T) {
	store := newTestSNSStore(t)

	topic, err := store.CreateTopic(&Topic{Name: "token-lifetime-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	created, _, err := store.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "http",
		Endpoint: "https://example.test/subscribe",
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	token := created.ConfirmationToken

	if _, err := store.ConfirmSubscription(created.SubscriptionArn, token, nil); err != nil {
		t.Fatalf("first confirmation: %v", err)
	}
	again, err := store.ConfirmSubscription(created.SubscriptionArn, token, nil)
	if err != nil {
		t.Fatalf("re-confirmation with the live token: %v", err)
	}
	if again.PendingConfirmation {
		t.Fatal("re-confirmation left the subscription pending")
	}
}

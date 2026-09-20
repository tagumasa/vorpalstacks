package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// newFanoutTestService builds an SNS service over a temporary regional
// storage (the signing key the HTTP delivery envelope carries is generated
// and persisted there) plus its default-region store, pre-seeded in the
// service's store cache the way the server wires it.
func newFanoutTestService(t *testing.T) (*SNSService, snsstore.SNSStoreInterface) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	t.Cleanup(func() { mgr.Close() })
	st, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("regional storage: %v", err)
	}
	store := snsstore.NewSNSStore(st, "123456789012", "us-east-1")
	svc := NewSNSService(mgr, "123456789012", "us-east-1")
	svc.SetSNSStore("us-east-1", store)
	return svc, store
}

// TestPublishFanOutWalksAllSubscriptionPages pins the fan-out seam end to
// end on the direct (bus-less) dispatch path: a topic holding more
// subscriptions than one list page (DefaultMaxItems is 100) delivers to
// EVERY endpoint — the 101st subscriber is not silently dropped — with the
// completion barrier Close provides before the assertions read the count.
func TestPublishFanOutWalksAllSubscriptionPages(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var mu sync.Mutex
	var envelopes []map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var envelope map[string]interface{}
		_ = json.Unmarshal(body, &envelope)
		mu.Lock()
		envelopes = append(envelopes, envelope)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	topic, err := store.CreateTopic(&snsstore.Topic{Name: "fanout-deliver-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}

	const want = 101
	for i := 0; i < want; i++ {
		created, _, err := store.CreateSubscription(&snsstore.Subscription{
			TopicArn: topic.Arn,
			Protocol: "http",
			// Distinct endpoints: the store is idempotent on the
			// subscription natural key (protocol+endpoint+owner), so one
			// shared endpoint would collapse the 101 subscriptions into a
			// single record before the fan-out ever runs.
			Endpoint: fmt.Sprintf("%s/%d", server.URL, i),
			Owner:    "123456789012",
		})
		if err != nil {
			t.Fatalf("create subscription %d: %v", i, err)
		}
		// The store creates every subscription pending its confirmation;
		// fan-out skips pending subscriptions, so confirm before publishing.
		if err := store.AutoConfirmSubscription(created); err != nil {
			t.Fatalf("confirm subscription %d: %v", i, err)
		}
	}

	result, err := svc.publishCore(store, "us-east-1", PublishInput{
		TopicArn: topic.Arn,
		Message:  "fan-out walks every page",
	})
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	messageID, _ := result.(map[string]interface{})["MessageId"].(string)
	if messageID == "" {
		t.Fatalf("publish result carries no MessageId: %v", result)
	}

	svc.Close()

	mu.Lock()
	defer mu.Unlock()
	if len(envelopes) != want {
		t.Fatalf("%d of %d endpoints received the notification — the fan-out stopped inside the first page", len(envelopes), want)
	}
	for i, envelope := range envelopes {
		if envelope["MessageId"] != messageID || envelope["Message"] != "fan-out walks every page" {
			t.Fatalf("envelope %d = %v, want the published MessageId and message", i, envelope)
		}
		if _, hasSignature := envelope["Signature"]; !hasSignature {
			t.Fatalf("envelope %d carries no Signature — HTTP delivery signs the notification", i)
		}
	}
}

// TestHTTPNotificationCarriesDocumentedHeaders pins the HTTP/S header set
// SNS POSTs carry: message-type, message-id, topic-arn and the subscription
// ARN — "x-amz-sns-subscription-arn – The ARN for the subscription to this
// endpoint" (HTTP/HTTPS headers).
func TestHTTPNotificationCarriesDocumentedHeaders(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var mu sync.Mutex
	var headers http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		headers = r.Header.Clone()
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	topic, err := store.CreateTopic(&snsstore.Topic{Name: "header-pin-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	created, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: topic.Arn,
		Protocol: "http",
		Endpoint: server.URL,
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(created); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}

	result, err := svc.publishCore(store, "us-east-1", PublishInput{
		TopicArn: topic.Arn,
		Message:  "header pin",
	})
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	messageID, _ := result.(map[string]interface{})["MessageId"].(string)

	svc.Close()

	mu.Lock()
	defer mu.Unlock()
	if headers == nil {
		t.Fatal("no notification reached the endpoint")
	}
	for header, want := range map[string]string{
		"x-amz-sns-message-type":     "Notification",
		"x-amz-sns-message-id":       messageID,
		"x-amz-sns-topic-arn":        topic.Arn,
		"x-amz-sns-subscription-arn": created.SubscriptionArn,
	} {
		if got := headers.Get(header); got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}
}

// TestSubscribeConfirmationResendPolicy pins the confirmation POST's
// resend policy: a pending subscription re-subscribed gets its
// SubscriptionConfirmation re-sent (the endpoint never confirmed — the
// retry gives it another chance), while a CONFIRMED subscription does not
// — the repeated Subscribe returns the live subscription and sends the
// endpoint no request it has no reason to act on again. The confirmation
// rides a tracked goroutine, so each phase polls for its arrival instead
// of closing the service mid-test (Close is terminal for the spawn path).
func TestSubscribeConfirmationResendPolicy(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var mu sync.Mutex
	var confirmations int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-amz-sns-message-type") == "SubscriptionConfirmation" {
			mu.Lock()
			confirmations++
			mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	waitForConfirmations := func(want int) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		for {
			mu.Lock()
			got := confirmations
			mu.Unlock()
			if got >= want {
				return
			}
			if time.Now().After(deadline) {
				t.Fatalf("confirmations = %d, want at least %d", got, want)
			}
			time.Sleep(2 * time.Millisecond)
		}
	}

	topic, err := store.CreateTopic(&snsstore.Topic{Name: "resend-policy-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	reqCtx := request.NewRequestContext(context.Background(), nil, "123456789012", "us-east-1")
	in := SubscribeInput{
		TopicArn: topic.Arn,
		Protocol: "http",
		Endpoint: server.URL,
	}

	if _, err := svc.subscribeCore(store, reqCtx, in); err != nil {
		t.Fatalf("first subscribe: %v", err)
	}
	waitForConfirmations(1)

	// The pending subscription re-subscribed: the confirmation re-sends.
	if _, err := svc.subscribeCore(store, reqCtx, in); err != nil {
		t.Fatalf("second subscribe (pending): %v", err)
	}
	waitForConfirmations(2)

	// The confirmed subscription re-subscribed: no further confirmation.
	subs, err := store.ListAllSubscriptionsByTopic(topic.Arn)
	if err != nil || len(subs) != 1 {
		t.Fatalf("list subscriptions: %d (%v)", len(subs), err)
	}
	if _, err := store.ConfirmSubscription(subs[0].SubscriptionArn, subs[0].ConfirmationToken, nil); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if _, err := svc.subscribeCore(store, reqCtx, in); err != nil {
		t.Fatalf("third subscribe (confirmed): %v", err)
	}
	svc.Close()
	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if confirmations != 2 {
		t.Fatalf("confirmed re-subscribe sent %d confirmations in total, want 2 — no re-send after confirmation", confirmations)
	}
}

// TestPublishCoreSurfacesSubscriptionListError pins the error policy of the
// fan-out seam: a subscription listing that fails is returned to the
// caller, never swallowed into a success response with zero delivery.
func TestPublishCoreSurfacesSubscriptionListError(t *testing.T) {
	svc, _ := newFanoutTestService(t)

	_, err := svc.publishCore(failingListStore{}, "us-east-1", PublishInput{
		TopicArn: "arn:aws:sns:us-east-1:123456789012:missing-topic",
		Message:  "m",
	})
	if err == nil {
		t.Fatal("publishCore swallowed the subscription listing failure")
	}
	if want := "list subscriptions: store unavailable"; err.Error() != want {
		t.Fatalf("publishCore error = %q, want the surfaced listing error %q", err.Error(), want)
	}
}

// failingListStore is a store whose every operation fails: GetTopic
// succeeds so the failure under test is the subscription listing alone.
type failingListStore struct {
	snsstore.SNSStoreInterface
}

func (failingListStore) GetTopic(topicArn string) (*snsstore.Topic, error) {
	return &snsstore.Topic{Arn: topicArn, Name: "missing-topic"}, nil
}

func (failingListStore) ListAllSubscriptionsByTopic(topicArn string) ([]*snsstore.Subscription, error) {
	return nil, fmt.Errorf("list subscriptions: store unavailable")
}

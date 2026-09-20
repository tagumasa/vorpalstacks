package sns

// Delivery-engine pins: the HTTP/S delivery policy's retry ladder runs
// until success or exhaustion, permanent failures do not retry, the SQS
// authorisation check fails closed, and the effective policy resolution
// honours the subscription override and the topic default.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// TestHTTPRetryLadderRetriesUntilSuccess pins the honoured ladder
// end to end: a subscription policy of three one-second retries drives
// three attempts at a server that fails twice with 500 (a documented
// retryable status) before accepting — the delivery succeeds on the
// ladder's third attempt, after the two scheduled one-second delays.
func TestHTTPRetryLadderRetriesUntilSuccess(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	created, err := store.CreateTopic(&snsstore.Topic{Name: "ladder-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	sub, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: created.Arn,
		Protocol: "http",
		Endpoint: server.URL,
		Owner:    "123456789012",
		Attributes: map[string]string{
			// minDelayTarget 1 keeps the ladder fast while still crossing
			// the documented constraint floor.
			"DeliveryPolicy": `{"healthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":1,"numRetries":3}}`,
		},
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(sub); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}

	started := time.Now()
	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: "ladder body"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	waitForAttempts(t, &attempts, 3)
	elapsed := time.Since(started)
	svc.Close()

	if got := attempts.Load(); got != 3 {
		t.Fatalf("server saw %d attempts, want 3 (two 500 retries then success)", got)
	}
	if elapsed < 2*time.Second {
		t.Errorf("three attempts completed in %s — the two scheduled one-second delays did not run", elapsed)
	}
	if elapsed := time.Since(started); elapsed < 2*time.Second {
		t.Errorf("three attempts completed in %s — the two scheduled one-second delays did not run", elapsed)
	}
}

// TestHTTPPermanentFailureDoesNotRetry pins the documented error
// classification: "All other errors are considered as permanent failures
// and retries will not be attempted" — a 404 endpoint receives exactly one
// attempt however many retries the policy allows.
func TestHTTPPermanentFailureDoesNotRetry(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	created, err := store.CreateTopic(&snsstore.Topic{Name: "permanent-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	sub, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: created.Arn,
		Protocol: "http",
		Endpoint: server.URL,
		Owner:    "123456789012",
		Attributes: map[string]string{
			"DeliveryPolicy": `{"healthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":1,"numRetries":3}}`,
		},
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(sub); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}

	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: "permanent body"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	waitForAttempts(t, &attempts, 1)
	// Give a would-be retry a moment to prove its absence: a permanent
	// failure must add no further attempts.
	time.Sleep(300 * time.Millisecond)
	svc.Close()

	if got := attempts.Load(); got != 1 {
		t.Fatalf("server saw %d attempts, want 1 — a 404 is a permanent failure and must not retry", got)
	}
}

// TestHTTPRetryExhaustionRoutesToDLQ pins the ladder's terminal
// behaviour: every attempt failing with a retryable status exhausts the
// policy and routes the message to the subscription's dead-letter queue.
func TestHTTPRetryExhaustionRoutesToDLQ(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var dlqBodies atomic.Int32
	// Route bus deliveries straight into the service's own handler (the
	// direct-dispatch stand-in for the real bus's workers) while the DLQ
	// redrive reaches the recording SQS invoker.
	svc.bus = &stubBus{
		sqs:       &recordingSQSInvoker{sends: func() { dlqBodies.Add(1) }},
		onPublish: func(evt *eventbus.SNSDeliveryEvent) { svc.handleBusDelivery(context.Background(), evt) },
	}

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	created, err := store.CreateTopic(&snsstore.Topic{Name: "exhaust-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	sub, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: created.Arn,
		Protocol: "http",
		Endpoint: server.URL,
		Owner:    "123456789012",
		Attributes: map[string]string{
			"DeliveryPolicy": `{"healthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":1,"numRetries":2}}`,
			"RedrivePolicy":  `{"deadLetterTargetArn":"arn:aws:sqs:us-east-1:123456789012:exhaust-dlq"}`,
		},
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(sub); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}

	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: "exhaust body"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	svc.Close()

	if got := attempts.Load(); got != 3 {
		t.Fatalf("server saw %d attempts, want 3 (initial attempt plus 2 retries)", got)
	}
	if got := dlqBodies.Load(); got != 1 {
		t.Fatalf("DLQ received %d redrives, want 1 after the exhausted ladder", got)
	}
}

// waitForAttempts blocks until the server has seen want attempts or a
// deadline passes. The assertions must not close the service first: Close
// cancels pending retries by design, so the counter is read live.
func waitForAttempts(t *testing.T, attempts *atomic.Int32, want int32) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for attempts.Load() < want && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
}

// recordingSQSInvoker records SendMessage calls and fails ARN resolution —
// the shape the fail-closed and redrive pins need.
type recordingSQSInvoker struct {
	sends func()
}

func (s *recordingSQSInvoker) GetQueueByName(ctx context.Context, region, queueName string) (string, error) {
	return "https://sqs." + region + ".amazonaws.com/123456789012/" + queueName, nil
}

func (s *recordingSQSInvoker) GetQueueARN(ctx context.Context, region, queueURL string) (string, error) {
	return "", errors.New("ARN resolution unavailable")
}

func (s *recordingSQSInvoker) SendMessage(ctx context.Context, region, queueURL, body string, opts invokers.SQSSendOptions) (string, string, error) {
	if s.sends != nil {
		s.sends()
	}
	return "recorded", "digest", nil
}

func (s *recordingSQSInvoker) ReceiveMessage(ctx context.Context, region, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]invokers.ReceivedSQSMessage, error) {
	return nil, errors.New("unreachable in this test")
}

func (s *recordingSQSInvoker) DeleteMessage(ctx context.Context, region, queueURL, receiptHandle string) error {
	return errors.New("unreachable in this test")
}

// TestSQSDeliveryFailsClosedOnARNResolution pins the authorisation seam:
// a queue whose ARN resolution fails must NOT be delivered to — the SQS
// twin now fails closed exactly as the Lambda twin always has, instead of
// skipping the resource-policy evaluation and sending unverified.
func TestSQSDeliveryFailsClosedOnARNResolution(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var sends atomic.Int32
	svc.bus = &stubBus{
		sqs:       &recordingSQSInvoker{sends: func() { sends.Add(1) }},
		onPublish: func(evt *eventbus.SNSDeliveryEvent) { svc.handleBusDelivery(context.Background(), evt) },
	}

	created, err := store.CreateTopic(&snsstore.Topic{Name: "failclosed-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	sub, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: created.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:unresolvable-queue",
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(sub); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}

	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: "fail-closed body"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	svc.Close()

	if got := sends.Load(); got != 0 {
		t.Fatalf("SQS delivery sent %d messages past a failed ARN resolution, want 0", got)
	}
}

// TestTopicDeliveryPolicyAppliesToSubscriptions pins the topic-level
// policy as the default its HTTP/S subscriptions resolve: a topic policy
// of one fast retry drives a subscription that sets no policy of its own.
func TestTopicDeliveryPolicyAppliesToSubscriptions(t *testing.T) {
	svc, store := newFanoutTestService(t)

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	created, err := store.CreateTopic(&snsstore.Topic{
		Name: "topic-policy-topic",
		Attributes: map[string]string{
			"DeliveryPolicy": `{"http":{"defaultHealthyRetryPolicy":{"minDelayTarget":1,"maxDelayTarget":1,"numRetries":1}}}`,
		},
	}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	sub, _, err := store.CreateSubscription(&snsstore.Subscription{
		TopicArn: created.Arn,
		Protocol: "http",
		Endpoint: server.URL,
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if err := store.AutoConfirmSubscription(sub); err != nil {
		t.Fatalf("confirm subscription: %v", err)
	}

	if _, err := svc.publishCore(store, "us-east-1", PublishInput{TopicArn: created.Arn, Message: "topic policy body"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	waitForAttempts(t, &attempts, 2)
	svc.Close()

	if got := attempts.Load(); got != 2 {
		t.Fatalf("server saw %d attempts, want 2 — the topic's default retry policy must reach subscriptions without their own", got)
	}
}

// TestResolveQueueURLRegionFallbacks pins the SQS endpoint region
// resolution: an AWS-form queue URL names its own region, a region-less
// platform URL falls back to the topic's region (never the global
// default), and a queue ARN names its region with the URL resolved
// through the invoker.
func TestResolveQueueURLRegionFallbacks(t *testing.T) {
	// The invoker is only consulted for ARN endpoints; a nil-safe stub is
	// unnecessary for the URL arms.
	var invoker invokers.SQSInvoker

	if url, region := resolveQueueURL(invoker, "https://sqs.eu-west-1.amazonaws.com/123456789012/queue", "us-east-1"); region != "eu-west-1" || url != "https://sqs.eu-west-1.amazonaws.com/123456789012/queue" {
		t.Fatalf("AWS-form URL resolved %q/%q, want the URL's own region eu-west-1", url, region)
	}
	if _, region := resolveQueueURL(invoker, "http://127.0.0.1:50080/queue/123456789012/queue", "ap-southeast-2"); region != "ap-southeast-2" {
		t.Fatalf("platform URL resolved region %q, want the topic's region ap-southeast-2", region)
	}
}

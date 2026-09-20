package sns

// Error-plane pins for the read sites: a storage read failure surfaces as
// itself at every one of them, never as the not-found sentinel — the miss
// and the fault are distinct answers (a 404 the client can react to vs an
// infrastructure error the service layer reports as InternalError).

import (
	"errors"

	"testing"

	"vorpalstacks/internal/core/storage"
)

// failingReadStorage delegates every operation to the wrapped storage but
// fails bucket reads with the injected error.
type failingReadStorage struct {
	storage.TransactionalStorageWith2PC
	failErr error
}

func (f failingReadStorage) Bucket(name string) storage.Bucket {
	return failingReadBucket{Bucket: f.TransactionalStorageWith2PC.Bucket(name), failErr: f.failErr}
}

type failingReadBucket struct {
	storage.Bucket
	failErr error
}

func (b failingReadBucket) Get(key []byte) ([]byte, error) { return nil, b.failErr }

// tokenIndexFailingStorage delegates to the wrapped storage but fails
// reads on the endpoint-token index bucket alone.
type tokenIndexFailingStorage struct {
	storage.TransactionalStorageWith2PC
	failErr error
}

func (f tokenIndexFailingStorage) Bucket(name string) storage.Bucket {
	b := f.TransactionalStorageWith2PC.Bucket(name)
	if name == "sns-endpoint-token-index-us-east-1" {
		return failingReadBucket{Bucket: b, failErr: f.failErr}
	}
	return b
}

// TestTokenIndexReadFailureIsNotAMiss pins the token-index leg of endpoint
// creation: a read failure on the index must fail the call — proceeding
// would mint a second endpoint for a token the index still maps — while
// the healthy twin resolves the same token to the existing endpoint.
func TestTokenIndexReadFailureIsNotAMiss(t *testing.T) {
	injected := errors.New("injected index fault")

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	seed := NewSNSStore(st, "123456789012", "us-east-1")
	t.Cleanup(seed.Close)
	app, err := seed.CreatePlatformApplication(&PlatformApplication{Name: "token-fault-app"})
	if err != nil {
		t.Fatalf("create platform application: %v", err)
	}
	first, err := seed.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  "shared-token",
	})
	if err != nil {
		t.Fatalf("seed endpoint: %v", err)
	}

	// Healthy twin: the same token resolves to the existing endpoint.
	merged, err := seed.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  "shared-token",
	})
	if err != nil {
		t.Fatalf("healthy re-create: %v", err)
	}
	if merged.EndpointArn != first.EndpointArn {
		t.Fatalf("healthy re-create minted %s, want the existing %s", merged.EndpointArn, first.EndpointArn)
	}

	failing := NewSNSStore(tokenIndexFailingStorage{TransactionalStorageWith2PC: st, failErr: injected}, "123456789012", "us-east-1")
	t.Cleanup(failing.Close)
	_, err = failing.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  "shared-token",
	})
	if err == nil {
		t.Fatal("token-index read failure was treated as an absent token — a duplicate endpoint was minted")
	}
	if !errors.Is(err, injected) {
		t.Fatalf("token-index read failure surfaced as %v, want the injected fault", err)
	}
}

// TestReadFailureIsNotNotFound drives every resource-read site over a
// storage whose reads fail: each must return the fault (the injected error
// in the chain, not the not-found sentinel), while the same keys over the
// healthy storage resolve — the miss, by contrast, is the sentinel.
func TestReadFailureIsNotNotFound(t *testing.T) {
	injected := errors.New("injected storage fault")

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	seed := NewSNSStore(st, "123456789012", "us-east-1")
	t.Cleanup(seed.Close)

	topic, err := seed.CreateTopic(&Topic{Name: "fault-pin-topic"}, nil)
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	sub, _, err := seed.CreateSubscription(&Subscription{
		TopicArn: topic.Arn,
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:fault-pin-queue",
		Owner:    "123456789012",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	app, err := seed.CreatePlatformApplication(&PlatformApplication{Name: "fault-pin-app"})
	if err != nil {
		t.Fatalf("create platform application: %v", err)
	}
	endpoint, err := seed.CreatePlatformEndpoint(&PlatformEndpoint{
		PlatformApplicationArn: app.PlatformApplicationArn,
		Token:                  "fault-pin-token",
	})
	if err != nil {
		t.Fatalf("create platform endpoint: %v", err)
	}

	failing := NewSNSStore(failingReadStorage{TransactionalStorageWith2PC: st, failErr: injected}, "123456789012", "us-east-1")
	t.Cleanup(failing.Close)

	reads := map[string]struct {
		sentinel error
		call     func() error
	}{
		"GetTopic":                {ErrTopicNotFound, func() error { _, err := failing.GetTopic(topic.Arn); return err }},
		"SetTopicAttributes":      {ErrTopicNotFound, func() error { return failing.SetTopicAttributes(topic.Arn, map[string]string{AttrDisplayName: "x"}) }},
		"GetDataProtectionPolicy": {ErrTopicNotFound, func() error { _, err := failing.GetDataProtectionPolicy(topic.Arn); return err }},
		"AddPermission":           {ErrTopicNotFound, func() error { return failing.AddPermission(topic.Arn, &Permission{Label: "l"}) }},
		"RemovePermission":        {ErrTopicNotFound, func() error { return failing.RemovePermission(topic.Arn, "l") }},
		"DeleteTopic":             {ErrTopicNotFound, func() error { return failing.DeleteTopic(topic.Arn) }},
		"CreateSubscription topic read": {ErrTopicNotFound, func() error {
			_, _, err := failing.CreateSubscription(&Subscription{TopicArn: topic.Arn, Protocol: "sqs", Endpoint: "arn:aws:sqs:us-east-1:123456789012:other", Owner: "123456789012"})
			return err
		}},
		"GetSubscription": {ErrSubscriptionNotFound, func() error { _, err := failing.GetSubscription(sub.SubscriptionArn); return err }},
		"SetSubscriptionAttributes": {ErrSubscriptionNotFound, func() error {
			return failing.SetSubscriptionAttributes(sub.SubscriptionArn, map[string]string{AttrDeliveryPolicy: "{}"})
		}},
		"GetPlatformApplication": {ErrPlatformApplicationNotFound, func() error { _, err := failing.GetPlatformApplication(app.PlatformApplicationArn); return err }},
		"SetPlatformApplicationAttributes": {ErrPlatformApplicationNotFound, func() error {
			return failing.SetPlatformApplicationAttributes(app.PlatformApplicationArn, map[string]string{"key": "v"})
		}},
		"CreatePlatformEndpoint app read": {ErrPlatformApplicationNotFound, func() error {
			_, err := failing.CreatePlatformEndpoint(&PlatformEndpoint{PlatformApplicationArn: app.PlatformApplicationArn, Token: "another-token"})
			return err
		}},
		"DeletePlatformApplication": {ErrPlatformApplicationNotFound, func() error { return failing.DeletePlatformApplication(app.PlatformApplicationArn) }},
		"GetEndpoint":               {ErrEndpointNotFound, func() error { _, err := failing.GetEndpoint(endpoint.EndpointArn); return err }},
		"SetEndpointAttributes": {ErrEndpointNotFound, func() error {
			return failing.SetEndpointAttributes(endpoint.EndpointArn, map[string]string{"key": "v"})
		}},
		"DeleteEndpoint": {ErrEndpointNotFound, func() error { return failing.DeleteEndpoint(endpoint.EndpointArn) }},
	}

	for name, tc := range reads {
		err := tc.call()
		if err == nil {
			t.Errorf("%s: read failure swallowed", name)
			continue
		}
		if errors.Is(err, tc.sentinel) {
			t.Errorf("%s: read failure surfaced as the not-found sentinel %v", name, tc.sentinel)
		}
		if !errors.Is(err, injected) {
			t.Errorf("%s: read failure surfaced as %v, want the injected fault in the chain", name, err)
		}
	}

	// The miss arm: over healthy storage the same lookups answer with
	// their sentinels — the distinction the sweep preserves.
	for name, tc := range map[string]struct {
		sentinel error
		call     func() error
	}{
		"GetTopic miss":               {ErrTopicNotFound, func() error { _, err := seed.GetTopic(topic.Arn + ":missing"); return err }},
		"GetSubscription miss":        {ErrSubscriptionNotFound, func() error { _, err := seed.GetSubscription(sub.SubscriptionArn + "x"); return err }},
		"GetPlatformApplication miss": {ErrPlatformApplicationNotFound, func() error { _, err := seed.GetPlatformApplication(app.PlatformApplicationArn + "x"); return err }},
		"GetEndpoint miss":            {ErrEndpointNotFound, func() error { _, err := seed.GetEndpoint(endpoint.EndpointArn + "x"); return err }},
	} {
		if err := tc.call(); !errors.Is(err, tc.sentinel) {
			t.Errorf("%s: got %v, want %v", name, err, tc.sentinel)
		}
	}
}

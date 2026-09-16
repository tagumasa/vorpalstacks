package eventbridge

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// failingReadBucket fails every read with the injected error, standing in
// for a transient storage fault during bus delivery.
type failingReadBucket struct {
	storage.Bucket
	fail error
}

func (b *failingReadBucket) Get(key []byte) ([]byte, error)           { return nil, b.fail }
func (b *failingReadBucket) ForEach(fn func(k, v []byte) error) error { return b.fail }

// faultedReadStorage routes one bucket's reads through the fault.
type faultedReadStorage struct {
	storage.BasicStorage
	failOn string
	fail   error
}

func (s *faultedReadStorage) Bucket(name string) storage.Bucket {
	b := s.BasicStorage.Bucket(name)
	if name == s.failOn {
		return &failingReadBucket{Bucket: b, fail: s.fail}
	}
	return b
}

// newBlockedStorageManager returns a manager whose region directories can
// never be created — the blocker occupies the parent path — so store
// acquisition surfaces a real storage fault instead of a nil store.
func newBlockedStorageManager(t *testing.T) *storage.RegionStorageManager {
	t.Helper()
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("occupies the region-directory name"), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: blocker})
	if err != nil {
		t.Fatal(err)
	}
	return mgr
}

// TestPutEventsHandlerAnswersDropsAsFailures pins the bus contract on the
// PutEvents bus handler: every condition that drops or fails the delivery —
// an event without a region, a storage fault, a malformed Input, a missing
// Source/DetailType pair — is reported in HandlerResult.Error. Synchronous
// publishers (the Scheduler engine publishes via PublishSync and folds
// result.Error into their retry and dead-letter policy) would otherwise
// record a dropped event as a successful delivery and lose it silently.
func TestPutEventsHandlerAnswersDropsAsFailures(t *testing.T) {
	ctx := context.Background()

	healthyMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}

	healthy := func() *EventsService { return NewEventsService(healthyMgr, "000000000000") }

	cases := []struct {
		name   string
		svc    *EventsService
		region string
		input  string
	}{
		{"empty region", healthy(), "", `{"Source":"aws.test","DetailType":"test"}`},
		{"storage fault", NewEventsService(newBlockedStorageManager(t), "000000000000"),
			"us-east-1", `{"Source":"aws.test","DetailType":"test"}`},
		{"malformed input", healthy(), "us-east-1", `"a scalar is not a PutEvents payload"`},
		{"missing source and detail type", healthy(), "us-east-1", `{"note":"no source members"}`},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			evt := &eventbus.EventBridgePutEventsEvent{Input: tt.input}
			evt.Region = tt.region
			result := tt.svc.handlePutEventsEvent(ctx, evt)
			if result.Error == nil {
				t.Fatal("a dropped or failed delivery must be reported in HandlerResult.Error, not recorded as success")
			}
		})
	}

	// A deliverable event stays successful: no rules match on the fresh
	// store, so the delivery completes without a failure.
	ok := &eventbus.EventBridgePutEventsEvent{
		Input: `{"Source":"aws.test","DetailType":"test","Detail":{"k":"v"}}`,
	}
	ok.Region = "us-east-1"
	if result := healthy().handlePutEventsEvent(ctx, ok); result.Error != nil {
		t.Fatalf("a deliverable event must stay successful: %v", result.Error)
	}
}

// TestBusDeliveryHandlerAnswersDropsAsFailures pins the same contract on the
// event-bus target delivery handler: a malformed payload, a failing rule
// listing and an unavailable store are reported in HandlerResult.Error, and
// a region without a cached store acquires one on demand instead of
// dropping the delivery.
func TestBusDeliveryHandlerAnswersDropsAsFailures(t *testing.T) {
	ctx := context.Background()
	busARN := "arn:aws:events:us-east-1:000000000000:event-bus/default"

	malformed := &eventbus.EventBridgeDeliveryEvent{
		TargetARN: busARN,
		Input:     []byte("not-json"),
	}
	malformed.Region = "us-east-1"
	if res := NewEventsService(nil, "000000000000").handleEventBusDelivery(ctx, malformed); res.Error == nil {
		t.Fatal("a malformed delivery payload must be reported as a failed delivery")
	}

	fault := errors.New("pebble: transient io fault")
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	// The default bus must exist so the handler's bus-existence check
	// passes and the injected rule-listing fault is what surfaces.
	seed := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	if err := seed.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatal(err)
	}
	faulted := eventsstore.NewEventsStore(&faultedReadStorage{
		BasicStorage: st,
		failOn:       "events-rules-us-east-1",
		fail:         fault,
	}, "000000000000", "us-east-1")
	faultSvc := NewEventsService(nil, "000000000000")
	faultSvc.SetEventsStore("us-east-1", faulted)
	delivery := &eventbus.EventBridgeDeliveryEvent{TargetARN: busARN, Input: []byte(`{}`)}
	delivery.Region = "us-east-1"
	res := faultSvc.handleEventBusDelivery(ctx, delivery)
	if res.Error == nil {
		t.Fatal("a failing rule listing must surface as a failed delivery")
	}
	if !strings.Contains(res.Error.Error(), fault.Error()) {
		t.Fatalf("the fault must surface as itself, got: %v", res.Error)
	}

	blocked := NewEventsService(newBlockedStorageManager(t), "000000000000")
	if res := blocked.handleEventBusDelivery(ctx, delivery); res.Error == nil {
		t.Fatal("an unavailable store must surface as a failed delivery")
	}

	// A region without a cached store must acquire one and deliver the
	// event rather than silently skipping it.
	healthyMgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	fresh := NewEventsService(healthyMgr, "000000000000")
	if res := fresh.handleEventBusDelivery(ctx, delivery); res.Error != nil {
		t.Fatalf("an uncached region must acquire its store and deliver: %v", res.Error)
	}
}

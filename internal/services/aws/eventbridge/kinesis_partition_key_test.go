package eventbridge

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// recordingKinesisInvoker captures the partition key and payload of every
// PutRecord so tests can assert what a retried delivery actually sent.
type recordingKinesisInvoker struct {
	mu       sync.Mutex
	putKeys  []string
	putData  [][]byte
	failNext int
}

func (r *recordingKinesisInvoker) ListShards(ctx context.Context, streamName string) ([]invokers.ShardInfo, error) {
	return nil, nil
}

func (r *recordingKinesisInvoker) PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.putKeys = append(r.putKeys, partitionKey)
	r.putData = append(r.putData, append([]byte(nil), data...))
	if r.failNext > 0 {
		r.failNext--
		return "", errors.New("kinesis unavailable")
	}
	return "seq-1", nil
}

func (r *recordingKinesisInvoker) CreateShardIterator(ctx context.Context, streamName string, shardID string, iteratorType string, startingSequenceNumber string, timestamp *time.Time) (string, error) {
	return "", nil
}

func (r *recordingKinesisInvoker) GetRecords(ctx context.Context, streamName string, shardID string, startingSequenceNumber string, limit int32, includeStart bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", nil
}

func (r *recordingKinesisInvoker) StreamExists(ctx context.Context, region, streamARN string) (bool, error) {
	return true, nil
}

func (r *recordingKinesisInvoker) keys() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.putKeys...)
}

func (r *recordingKinesisInvoker) payloads() [][]byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([][]byte(nil), r.putData...)
}

// The Kinesis partition key must default to the event ID (the behaviour
// documented for Kinesis targets in the EventBridge API reference) and
// stay stable across retries: a per-attempt random key would scatter one
// logical event across different shards.
func TestKinesisPartitionKeyIsStableAcrossRetries(t *testing.T) {
	// TEST_MODE selects the short retry backoff so the retry fires quickly.
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	invoker := &recordingKinesisInvoker{failNext: 1}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc.SetEventBus(bus)

	event := &eventsstore.Event{ID: "evt-partition-stability"}
	target := eventsstore.Target{
		ARN: "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
	}

	svc.dispatchToTarget(context.Background(), "us-east-1", "", event, target, []byte(`{"detail":{}}`))

	keys := invoker.keys()
	if len(keys) != 2 {
		t.Fatalf("expected a failed attempt and a retry, got %d puts", len(keys))
	}
	for _, key := range keys {
		if key != event.ID {
			t.Fatalf("partition key %q does not match the event ID %q", key, event.ID)
		}
	}
}

// An explicit PartitionKeyPath that resolves keeps precedence over the
// event-ID default, and the path resolves against the original event —
// not the transformed payload the target receives.
func TestKinesisPartitionKeyPathTakesPrecedence(t *testing.T) {
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	invoker := &recordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc.SetEventBus(bus)

	event := &eventsstore.Event{
		ID:     "evt-with-path",
		Detail: map[string]interface{}{"key": "shard-me"},
	}
	target := eventsstore.Target{
		ARN:               "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		KinesisParameters: &eventsstore.KinesisParameters{PartitionKeyPath: "$.detail.key"},
	}

	// The payload carries a different key: only the original event may
	// feed the partition-key path.
	svc.dispatchToTarget(context.Background(), "us-east-1", "", event, target, []byte(`{"detail":{"key":"transformed-key"}}`))

	keys := invoker.keys()
	if len(keys) != 1 {
		t.Fatalf("expected exactly one put, got %d puts", len(keys))
	}
	if keys[0] != "shard-me" {
		t.Fatalf("expected the key extracted from the original event, got %q", keys[0])
	}
}

// On the bus delivery path the publisher resolves the partition key from
// the original event and the resolved key rides the bus event; the
// handler's identity stub must not override it.
func TestKinesisPartitionKeyRidesTheBus(t *testing.T) {
	t.Setenv("TEST_MODE", "true")

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	invoker := &recordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	svc.SetEventBus(bus)

	delivery := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:           "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		Input:               []byte(`{"detail":{"key":"transformed-key"}}`),
		EventBridgeEventID:  "evt-bus-path",
		KinesisPartitionKey: "publisher-key",
	}
	delivery.Region = "us-east-1"

	if res := svc.handleBusDelivery(context.Background(), delivery); res.Error != nil {
		t.Fatalf("delivery must succeed: %v", res.Error)
	}

	keys := invoker.keys()
	if len(keys) != 1 {
		t.Fatalf("expected exactly one put, got %d puts", len(keys))
	}
	if keys[0] != "publisher-key" {
		t.Fatalf("expected the publisher-resolved key, got %q", keys[0])
	}
}

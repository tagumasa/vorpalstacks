package eventbridge

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// recordingKinesisInvoker captures the partition key of every PutRecord so
// tests can assert what a retried delivery actually sent.
type recordingKinesisInvoker struct {
	mu       sync.Mutex
	putKeys  []string
	failNext int
}

func (r *recordingKinesisInvoker) ListShards(ctx context.Context, streamName string) ([]eventbus.ShardInfo, error) {
	return nil, nil
}

func (r *recordingKinesisInvoker) PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.putKeys = append(r.putKeys, partitionKey)
	if r.failNext > 0 {
		r.failNext--
		return "", errors.New("kinesis unavailable")
	}
	return "seq-1", nil
}

func (r *recordingKinesisInvoker) CreateShardIterator(ctx context.Context, streamName string, shardID string, iteratorType string, startingSequenceNumber string, timestamp *time.Time) (string, error) {
	return "", nil
}

func (r *recordingKinesisInvoker) GetRecords(ctx context.Context, streamName string, shardID string, startingSequenceNumber string, limit int32, includeStart bool) ([]eventbus.KinesisRecord, string, error) {
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

	svc.dispatchToTarget(context.Background(), "us-east-1", event, target, []byte(`{"detail":{}}`))

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
// event-ID default.
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

	event := &eventsstore.Event{ID: "evt-with-path"}
	target := eventsstore.Target{
		ARN:               "arn:aws:kinesis:us-east-1:000000000000:stream/demo",
		KinesisParameters: &eventsstore.KinesisParameters{PartitionKeyPath: "$.detail.key"},
	}

	svc.dispatchToTarget(context.Background(), "us-east-1", event, target, []byte(`{"detail":{"key":"shard-me"}}`))

	keys := invoker.keys()
	if len(keys) != 1 {
		t.Fatalf("expected exactly one put, got %d", len(keys))
	}
	if keys[0] != "shard-me" {
		t.Fatalf("expected the extracted partition key, got %q", keys[0])
	}
}

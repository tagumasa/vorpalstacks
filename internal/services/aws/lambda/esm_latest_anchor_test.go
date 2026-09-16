package lambda

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// memBucket is an in-memory storage.Bucket for poller tests that need a real
// EventSourceStore (the poller reports processing results into it).
type memBucket struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newMemBucket() *memBucket { return &memBucket{data: map[string][]byte{}} }

func (b *memBucket) Get(key []byte) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	v, ok := b.data[string(key)]
	if !ok {
		// Missing keys answer nil, nil — the production bucket contract
		// that BaseStore.Get's not-found detection relies on.
		return nil, nil
	}
	return append([]byte(nil), v...), nil
}

func (b *memBucket) Put(key, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[string(key)] = append([]byte(nil), value...)
	return nil
}

func (b *memBucket) Delete(key []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.data, string(key))
	return nil
}

func (b *memBucket) Has(key []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.data[string(key)]
	return ok
}

func (b *memBucket) ForEach(fn func(k, v []byte) error) error {
	b.mu.Lock()
	keys := make([]string, 0, len(b.data))
	for k := range b.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	snapshot := make(map[string][]byte, len(keys))
	for _, k := range keys {
		snapshot[k] = append([]byte(nil), b.data[k]...)
	}
	b.mu.Unlock()
	for _, k := range keys {
		if err := fn([]byte(k), snapshot[k]); err != nil {
			return err
		}
	}
	return nil
}

func (b *memBucket) ScanPrefix(prefix []byte) storage.Iterator {
	var keys []string
	_ = b.ForEach(func(k, v []byte) error {
		if strings.HasPrefix(string(k), string(prefix)) {
			keys = append(keys, string(k))
		}
		return nil
	})
	return &memIterator{bucket: b, keys: keys}
}

// ScanPrefixReverse walks the matching keys in descending order, starting at
// the largest key strictly less than before when before is non-nil.
func (b *memBucket) ScanPrefixReverse(prefix, before []byte) storage.Iterator {
	var keys []string
	_ = b.ForEach(func(k, v []byte) error {
		ks := string(k)
		if strings.HasPrefix(ks, string(prefix)) && (before == nil || ks < string(before)) {
			keys = append(keys, ks)
		}
		return nil
	})
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return &memIterator{bucket: b, keys: keys}
}

func (b *memBucket) ScanRange(start, end []byte) storage.Iterator {
	var keys []string
	_ = b.ForEach(func(k, v []byte) error {
		if string(k) >= string(start) && string(k) < string(end) {
			keys = append(keys, string(k))
		}
		return nil
	})
	return &memIterator{bucket: b, keys: keys}
}

func (b *memBucket) Count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.data)
}

type memIterator struct {
	bucket *memBucket
	keys   []string
	pos    int
}

func (it *memIterator) Next() bool { it.pos++; return it.pos <= len(it.keys) }
func (it *memIterator) Key() []byte {
	if it.pos == 0 || it.pos > len(it.keys) {
		return nil
	}
	return []byte(it.keys[it.pos-1])
}
func (it *memIterator) Value() []byte {
	if it.pos == 0 || it.pos > len(it.keys) {
		return nil
	}
	v, _ := it.bucket.Get([]byte(it.keys[it.pos-1]))
	return v
}
func (it *memIterator) Error() error { return nil }
func (it *memIterator) Close()       {}

// memStorage hands out the same in-memory bucket for every name; poller tests
// never mix buckets.
type memStorage struct{ bucket *memBucket }

func (s *memStorage) Close() error                 { return nil }
func (s *memStorage) Bucket(string) storage.Bucket { return s.bucket }
func (s *memStorage) CreateBucket(string) error    { return nil }
func (s *memStorage) DeleteBucket(string) error    { return nil }
func (s *memStorage) ListBuckets() []string        { return nil }

// scriptedKinesisStream mirrors the store semantics the poller relies on:
// a LATEST iterator anchors at the shard's latest sequence (the shard floor
// for an empty shard), and GetRecords reads strictly after the position —
// inclusively only when includeStart is set.
type scriptedKinesisStream struct {
	mu      sync.Mutex
	shards  []invokers.ShardInfo
	floor   string
	records []invokers.KinesisRecord
	invokes int
	reads   int
	anchors []string
	regions []string
}

func (s *scriptedKinesisStream) publish(seq string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, invokers.KinesisRecord{
		SequenceNumber:              seq,
		PartitionKey:                "pk-" + seq,
		Data:                        []byte("data-" + seq),
		ApproximateArrivalTimestamp: time.Now().UTC(),
	})
}

func (s *scriptedKinesisStream) ListShards(_ context.Context, region, _ string) ([]invokers.ShardInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.regions = append(s.regions, region)
	return s.shards, nil
}

func (s *scriptedKinesisStream) PutRecord(context.Context, string, string, string, []byte) (string, error) {
	return "", fmt.Errorf("not implemented in test")
}

func (s *scriptedKinesisStream) latest() string {
	if len(s.records) == 0 {
		return s.floor
	}
	return s.records[len(s.records)-1].SequenceNumber
}

func (s *scriptedKinesisStream) seenRegions() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.regions...)
}

func (s *scriptedKinesisStream) CreateShardIterator(_ context.Context, region, _ string, _ string, iteratorType string, _ string, _ *time.Time) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invokes++
	s.regions = append(s.regions, region)
	switch iteratorType {
	case "LATEST":
		anchor := s.latest()
		s.anchors = append(s.anchors, anchor)
		return anchor, nil
	default:
		return s.floor, nil
	}
}

func (s *scriptedKinesisStream) GetRecords(_ context.Context, region, _ string, _ string, from string, limit int32, includeStart bool) ([]invokers.KinesisRecord, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reads++
	s.regions = append(s.regions, region)
	var out []invokers.KinesisRecord
	for _, r := range s.records {
		if from != "" && ((includeStart && r.SequenceNumber < from) || (!includeStart && r.SequenceNumber <= from)) {
			continue
		}
		out = append(out, r)
		if int32(len(out)) >= limit {
			break
		}
	}
	next := from
	if len(out) > 0 {
		next = out[len(out)-1].SequenceNumber
	}
	return out, next, nil
}

func (s *scriptedKinesisStream) StreamExists(context.Context, string, string) (bool, error) {
	return true, nil
}

// fakeKinesisBus exposes only the Kinesis invoker; the poller calls nothing
// else on the bus for a mapping without failure destinations.
type fakeKinesisBus struct {
	eventbus.ServiceBus
	invoker *scriptedKinesisStream
}

func (b *fakeKinesisBus) KinesisInvoker() invokers.KinesisInvoker { return b.invoker }

// TestKinesisLatestAnchorPersistedForBufferedMapping pins the anchoring of a
// LATEST mapping: the first empty poll must durably record the anchor so a
// burst arriving between polls is read in full. Without the durable anchor
// every cycle re-anchors at the newest record and only that record survives.
func TestKinesisLatestAnchorPersistedForBufferedMapping(t *testing.T) {
	const (
		streamARN = "arn:aws:kinesis:us-east-1:123456789012:stream/anchor-stream"
		funcARN   = "arn:aws:lambda:us-east-1:123456789012:function:anchor-fn"
	)
	esmStore := lambdastore.NewEventSourceStore(&memStorage{bucket: newMemBucket()}, "123456789012", "us-east-1")
	created, err := esmStore.Create(&lambdastore.EventSourceMapping{
		EventSourceArn:                 streamARN,
		FunctionArn:                    funcARN,
		StartingPosition:               "LATEST",
		BatchSize:                      2,
		MaximumBatchingWindowInSeconds: 1,
		ParallelizationFactor:          2,
		State:                          "Enabled",
	})
	if err != nil {
		t.Fatalf("create mapping: %v", err)
	}

	stream := &scriptedKinesisStream{
		shards: []invokers.ShardInfo{{ShardID: "shard-0"}},
		floor:  "0",
	}
	invoke := &capturePayloads{}
	p := &esmPoller{
		bus:       &fakeKinesisBus{invoker: stream},
		invoke:    invoke.invoke,
		kinesisCP: make(map[string]string),
		buffers:   make(map[string]*streamBuffer),
		esmStore:  esmStore,
		// The rendered records carry the invoke identity ARN, which is
		// built from the service account ID.
		lambdaSvc: &LambdaService{accountID: "123456789012"},
	}
	cpKey := fmt.Sprintf("%s:%s:%s", created.UUID, "anchor-stream", "shard-0")

	// Cycle 1: the stream is empty. The anchor must be recorded durably.
	p.processKinesisMapping(context.Background(), created)
	if got := p.kinesisCP[cpKey]; got != "0" {
		t.Fatalf("after the first empty poll the anchor must be persisted at the shard floor, got checkpoint %q", got)
	}

	// A burst arrives between polls.
	stream.publish("1")
	stream.publish("2")
	stream.publish("3")

	// Cycle 2: the buffer gathers the burst — the batching window flushes it
	// in batch-size chunks, so every record of the burst must have been
	// invoked in this single cycle.
	p.processKinesisMapping(context.Background(), created)
	allPayloads := strings.Join(invoke.snapshot(), "\n")
	for _, seq := range []string{"1", "2", "3"} {
		if !strings.Contains(allPayloads, fmt.Sprintf(`"sequenceNumber":%q`, seq)) {
			t.Fatalf("invocation payloads must carry record %q of the burst, got %s", seq, allPayloads)
		}
	}
	if got := p.kinesisCP[cpKey]; got != "3" {
		t.Fatalf("checkpoint must advance past the flushed burst, got %q", got)
	}
}

// TestKinesisESMPollerAddressesStreamARNRegion pins that a Kinesis mapping's
// poll addresses the stream in the ARN's region for every read it makes:
// the shard listing, the iterator creation and the record fetch must all
// resolve the ARN's regional store — a default-region read never sees a
// cross-region stream and the mapping silently delivers nothing.
func TestKinesisESMPollerAddressesStreamARNRegion(t *testing.T) {
	const (
		streamARN = "arn:aws:kinesis:eu-west-1:123456789012:stream/cross-region-stream"
		funcARN   = "arn:aws:lambda:eu-west-1:123456789012:function:cross-region-fn"
	)
	esmStore := lambdastore.NewEventSourceStore(&memStorage{bucket: newMemBucket()}, "123456789012", "eu-west-1")
	created, err := esmStore.Create(&lambdastore.EventSourceMapping{
		EventSourceArn:   streamARN,
		FunctionArn:      funcARN,
		StartingPosition: "TRIM_HORIZON",
		BatchSize:        10,
		State:            "Enabled",
	})
	if err != nil {
		t.Fatalf("create mapping: %v", err)
	}

	stream := &scriptedKinesisStream{
		shards: []invokers.ShardInfo{{ShardID: "shard-0"}},
		floor:  "0",
	}
	stream.publish("1")
	p := &esmPoller{
		bus:       &fakeKinesisBus{invoker: stream},
		invoke:    (&capturePayloads{}).invoke,
		kinesisCP: make(map[string]string),
		buffers:   make(map[string]*streamBuffer),
		esmStore:  esmStore,
		lambdaSvc: &LambdaService{accountID: "123456789012"},
	}

	p.processKinesisMapping(context.Background(), created)

	regions := stream.seenRegions()
	if len(regions) < 3 {
		t.Fatalf("a poll cycle must list shards, create an iterator and read records, saw %d region-carrying calls", len(regions))
	}
	for _, region := range regions {
		if region != "eu-west-1" {
			t.Errorf("Kinesis read addressed region %q, want the stream ARN's eu-west-1 (all calls: %v)", region, regions)
		}
	}
}

// capturePayloads records every invoke payload and answers success.
type capturePayloads struct {
	mu       sync.Mutex
	payloads []string
}

func (c *capturePayloads) invoke(_ context.Context, _ string, payload []byte) (*lambdastore.InvocationResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payloads = append(c.payloads, string(payload))
	return &lambdastore.InvocationResult{StatusCode: 200, Payload: []byte(`{"ok":true}`)}, nil
}

func (c *capturePayloads) callCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.payloads)
}

func (c *capturePayloads) snapshot() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.payloads...)
}

package kinesis

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// ended reports whether the pair's registered subscription has stopped —
// the pump is gone, and only the window tombstone remains. A test-file
// diagnostic over the registry; production code never needs the state.
func (r *subscriptionRegistry) ended(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.entries[key] == nil || r.entries[key].ended
}

// TestSubscribeStartingPositionValidation pins the StartingPosition member
// pairing rules at the subscribe seam: the type must be one of the five
// modelled values, sequence-number positions must satisfy the
// SequenceNumber pattern and travel with their types, and a timestamp is
// required by AT_TIMESTAMP — no silent degradation to the horizon.
func TestSubscribeStartingPositionValidation(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subpos", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "pos-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	ctx := context.Background()

	cases := []struct {
		name string
		in   SubscribeToShardInput
	}{
		{"unknown type", SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "AT_MOON"}},
		{"empty type", SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000"}},
		{"AT_SEQUENCE_NUMBER without sequence", SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "AT_SEQUENCE_NUMBER"}},
		{"AFTER_SEQUENCE_NUMBER without sequence", SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "AFTER_SEQUENCE_NUMBER"}},
		{"AT_TIMESTAMP without timestamp", SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "AT_TIMESTAMP"}},
		{"malformed sequence number", SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "AT_SEQUENCE_NUMBER", StartingSequenceNumber: "12x"}},
	}
	for _, c := range cases {
		_, err := svc.subscribeToShardCore(ctx, reqCtx, c.in)
		requireAWSCode(t, c.name, err, "InvalidArgumentException")
	}
}

// TestSubscribeTakeoverWindow pins the documented re-subscription rules:
// a repeat subscription for the same consumer and shard within the
// takeover window answers ResourceInUseException, and one at or beyond
// the window takes the subscription over — the previous connection is
// cancelled and its stream ends.
func TestSubscribeTakeoverWindow(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subtake", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "takeover-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "TRIM_HORIZON"}
	first, err := svc.subscribeToShardCore(ctx, reqCtx, in)
	if err != nil {
		t.Fatalf("first subscription: %v", err)
	}

	// The client never reads: the pump parks on its initial write, which
	// is exactly the state the takeover must release.
	if _, err := svc.subscribeToShardCore(ctx, reqCtx, in); err == nil {
		t.Fatalf("second subscription within the window succeeded")
	} else {
		requireAWSCode(t, "second within window", err, "ResourceInUseException")
	}

	// Rewind the first subscription past the documented window.
	key := subscriptionKey(consumer.ConsumerARN, "shardId-000000000000")
	svc.subscriptions.mu.Lock()
	if sub := svc.subscriptions.entries[key]; sub != nil {
		sub.startedAt = time.Now().Add(-2 * kinesisstore.SubscribeTakeoverWindow)
	}
	svc.subscriptions.mu.Unlock()

	if _, err := svc.subscribeToShardCore(ctx, reqCtx, in); err != nil {
		t.Fatalf("takeover subscription: %v", err)
	}

	// The taken-over pump was cancelled, so the first stream ends.
	if _, err := io.ReadAll(first.GetStream()); err != nil {
		t.Fatalf("first stream after takeover: %v", err)
	}
}

// TestSubscribeLifetimeEndsStream pins the documented subscription
// lifetime: the stream serves events up to its duration and then ends —
// the client renews by calling SubscribeToShard again.
func TestSubscribeLifetimeEndsStream(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subttl", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "ttl-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	svc.subscriptions.ttlOverride = 150 * time.Millisecond

	es, err := svc.subscribeToShardCore(context.Background(), reqCtx, SubscribeToShardInput{
		ConsumerARN:          consumer.ConsumerARN,
		ShardId:              "shardId-000000000000",
		StartingPositionType: "TRIM_HORIZON",
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	data, err := io.ReadAll(es.GetStream())
	if err != nil {
		t.Fatalf("read stream: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("lifetime ended the stream before any event was served")
	}

	// The pump has exited; the registry holds only the window tombstone.
	key := subscriptionKey(consumer.ConsumerARN, "shardId-000000000000")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if svc.subscriptions.ended(key) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("subscription still live after its lifetime")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestSubscribeStalledClientReleases pins the stall handling: a client
// that stops reading cannot hold the pump — cancelling the request
// context unblocks the parked write and the pump exits through its
// deferred cleanup, its registry entry left as the window tombstone.
func TestSubscribeStalledClientReleases(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("substall", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "stall-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	es, err := svc.subscribeToShardCore(ctx, reqCtx, SubscribeToShardInput{
		ConsumerARN:          consumer.ConsumerARN,
		ShardId:              "shardId-000000000000",
		StartingPositionType: "TRIM_HORIZON",
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	// Let the pump park on a write the stalled client never reads.
	time.Sleep(200 * time.Millisecond)
	cancel()

	key := subscriptionKey(consumer.ConsumerARN, "shardId-000000000000")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if svc.subscriptions.ended(key) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("pump held by a stalled client after context cancellation")
		}
		time.Sleep(20 * time.Millisecond)
	}
	_, _ = io.ReadAll(es.GetStream())
}

// TestSubscribeEncryptionTypeThreaded pins the encryption parity between
// the two read transports: a KMS-encrypted stream's subscribe events
// report KMS, never the hardcoded NONE.
func TestSubscribeEncryptionTypeThreaded(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subenc", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "enc-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	if _, err := svc.startStreamEncryptionCore(reqCtx, StartStreamEncryptionInput{
		StreamName:     "subenc",
		EncryptionType: "KMS",
		KeyId:          "key-1",
	}); err != nil {
		t.Fatalf("start encryption: %v", err)
	}
	if _, _, err := store.PutRecordWithShardSelection("subenc", "pk", "ZGF0YQ==", ""); err != nil {
		t.Fatalf("seed record: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	es, err := svc.subscribeToShardCore(ctx, reqCtx, SubscribeToShardInput{
		ConsumerARN:          consumer.ConsumerARN,
		ShardId:              "shardId-000000000000",
		StartingPositionType: "TRIM_HORIZON",
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	drained := make(chan []byte, 1)
	go func() {
		b, _ := io.ReadAll(es.GetStream())
		drained <- b
	}()
	time.Sleep(1500 * time.Millisecond)
	cancel()

	select {
	case b := <-drained:
		if !bytes.Contains(b, []byte(`"EncryptionType":"KMS"`)) {
			t.Fatalf("subscribe records do not report the stream's KMS encryption")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("stream did not end after cancellation")
	}
}

// TestSubscribeReadErrorIdentity pins the error taxonomy: a read failure
// surfaces under its mapped identity — a deleted stream reports
// ResourceNotFoundException — instead of a blanket not-found for every
// failure class.
func TestSubscribeReadErrorIdentity(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("suberr", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "err-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	es, err := svc.subscribeToShardCore(ctx, reqCtx, SubscribeToShardInput{
		ConsumerARN:          consumer.ConsumerARN,
		ShardId:              "shardId-000000000000",
		StartingPositionType: "TRIM_HORIZON",
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	drained := make(chan []byte, 1)
	go func() {
		b, _ := io.ReadAll(es.GetStream())
		drained <- b
	}()

	time.Sleep(300 * time.Millisecond)
	if err := store.DeleteStream("suberr"); err != nil {
		t.Fatalf("delete stream: %v", err)
	}

	select {
	case b := <-drained:
		if !bytes.Contains(b, []byte("ResourceNotFoundException")) {
			t.Fatalf("error event does not carry the mapped identity")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("stream did not end after the read failure")
	}
	cancel()
}

// TestSubscribeTombstoneWindow pins the window's anchor: it runs from the
// successful call, not from the connection's lifetime — a re-subscription
// for the same pair within the window answers ResourceInUseException even
// after the previous pump has ended, and one at or beyond the window
// succeeds.
func TestSubscribeTombstoneWindow(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subtomb", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "tomb-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := SubscribeToShardInput{ConsumerARN: consumer.ConsumerARN, ShardId: "shardId-000000000000", StartingPositionType: "TRIM_HORIZON"}
	if _, err := svc.subscribeToShardCore(ctx, reqCtx, in); err != nil {
		t.Fatalf("first subscription: %v", err)
	}
	cancel()

	key := subscriptionKey(consumer.ConsumerARN, "shardId-000000000000")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if svc.subscriptions.ended(key) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("first pump did not end after cancellation")
		}
		time.Sleep(20 * time.Millisecond)
	}

	_, err = svc.subscribeToShardCore(context.Background(), reqCtx, in)
	requireAWSCode(t, "re-subscription within window after end", err, "ResourceInUseException")

	svc.subscriptions.mu.Lock()
	if sub := svc.subscriptions.entries[key]; sub != nil {
		sub.startedAt = time.Now().Add(-2 * kinesisstore.SubscribeTakeoverWindow)
	}
	svc.subscriptions.mu.Unlock()

	if _, err := svc.subscribeToShardCore(context.Background(), reqCtx, in); err != nil {
		t.Fatalf("re-subscription past window: %v", err)
	}
}

// TestSubscribeFailedCallLeavesNoTombstone pins the failed-call rule: the
// takeover window the documentation states runs from a successful call, so
// a subscription whose iterator creation failed after registration leaves
// no tombstone — an immediate retry for the same pair proceeds instead of
// answering ResourceInUseException.
func TestSubscribeFailedCallLeavesNoTombstone(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subfail", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "fail-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}

	// A well-formed shard that does not exist fails at iterator creation —
	// after the pair is registered.
	in := SubscribeToShardInput{
		ConsumerARN:          consumer.ConsumerARN,
		ShardId:              "shardId-000000000042",
		StartingPositionType: "TRIM_HORIZON",
	}
	_, err = svc.subscribeToShardCore(context.Background(), reqCtx, in)
	requireAWSCode(t, "failed subscribe", err, "ResourceNotFoundException")

	key := subscriptionKey(consumer.ConsumerARN, in.ShardId)
	svc.subscriptions.mu.Lock()
	_, present := svc.subscriptions.entries[key]
	svc.subscriptions.mu.Unlock()
	if present {
		t.Fatal("failed subscribe left a registry entry — the retry window would block a corrected call")
	}

	// The immediate retry fails on the shard again, never on the window.
	_, err = svc.subscribeToShardCore(context.Background(), reqCtx, in)
	requireAWSCode(t, "immediate retry after a failed call", err, "ResourceNotFoundException")
}

// TestSubscribeLifetimeReleasesParkedWrite pins the expiry in the parked
// state: a client that never reads parks the pump on its initial write,
// and the subscription's lifetime still ends the pump with no context
// cancellation — the lifetime arms the event writes themselves, not only
// the loop select a parked write cannot reach.
func TestSubscribeLifetimeReleasesParkedWrite(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("subpark", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "park-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	svc.subscriptions.ttlOverride = 200 * time.Millisecond

	// The context is never cancelled: only the lifetime can end the pump.
	if _, err := svc.subscribeToShardCore(context.Background(), reqCtx, SubscribeToShardInput{
		ConsumerARN:          consumer.ConsumerARN,
		ShardId:              "shardId-000000000000",
		StartingPositionType: "TRIM_HORIZON",
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	key := subscriptionKey(consumer.ConsumerARN, "shardId-000000000000")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if svc.subscriptions.ended(key) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("pump held by a stalled client past its lifetime — the parked write never observed the expiry")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestWriteCtxSurvivesWritePanic pins the writer goroutine's panic
// containment: a write that panics delivers an error to the waiting side
// through the buffered channel — the process survives and the caller ends
// the stream through its own cleanup instead of hanging on a delivery
// that would never come.
func TestWriteCtxSurvivesWritePanic(t *testing.T) {
	pr, pw := io.Pipe()
	defer pr.Close()

	err := writeCtx(context.Background(), pw, nil, func() error {
		panic("encoder bug")
	})
	if err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Fatalf("panicked write: got %v, want the panic-turned error", err)
	}
}

// TestSubscriptionTombstoneSelfRemoves pins the registry's bounded
// residency: an ended subscription's tombstone removes itself once the
// takeover window has definitely elapsed, so a pair never subscribed to
// again does not keep its entry forever.
func TestSubscriptionTombstoneSelfRemoves(t *testing.T) {
	r := &subscriptionRegistry{windowOverride: 20 * time.Millisecond}
	ctx := context.Background()
	if _, sub, ok := r.begin("c#s", ctx); !ok {
		t.Fatal("begin rejected the first subscription")
	} else {
		r.end("c#s", sub)
	}
	time.Sleep(60 * time.Millisecond)
	r.mu.Lock()
	remaining := len(r.entries)
	r.mu.Unlock()
	if remaining != 0 {
		t.Fatalf("registry entries after the window: got %d, want the tombstone gone", remaining)
	}
}

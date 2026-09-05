package lambda

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// TestParallelizationFactorOf_ClampsToModelRange pins the clamping rules:
// unset (0) behaves as the default single batch and out-of-range values
// snap to the documented bounds.
func TestParallelizationFactorOf_ClampsToModelRange(t *testing.T) {
	cases := []struct {
		setting, want int32
	}{
		{0, lambdastore.MinParallelizationFactor},
		{1, 1},
		{5, 5},
		{10, 10},
		{15, lambdastore.MaxParallelizationFactor},
		{-3, lambdastore.MinParallelizationFactor},
	}
	for _, tc := range cases {
		mapping := &lambdastore.EventSourceMapping{ParallelizationFactor: tc.setting}
		if got := parallelizationFactorOf(mapping); got != int(tc.want) {
			t.Fatalf("ParallelizationFactor %d resolved to %d, want %d", tc.setting, got, tc.want)
		}
	}
}

// TestRunOrderedBatches_DisjointKeysRunConcurrently pins the concurrency
// contract: batches without shared record keys overlap in flight.
func TestRunOrderedBatches_DisjointKeysRunConcurrently(t *testing.T) {
	var inFlight, maxInFlight int32
	keys := []map[string]struct{}{
		{"a": {}}, {"b": {}}, {"c": {}},
	}
	outcomes := runOrderedBatches(context.Background(), 3, keys, func(ctx context.Context, idx int) batchOutcome {
		current := atomic.AddInt32(&inFlight, 1)
		for {
			observed := atomic.LoadInt32(&maxInFlight)
			if current <= observed || atomic.CompareAndSwapInt32(&maxInFlight, observed, current) {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&inFlight, -1)
		return batchOutcome{lastConsumed: "seq", delivered: true}
	})
	if len(outcomes) != 3 {
		t.Fatalf("expected 3 outcomes, got %d", len(outcomes))
	}
	if atomic.LoadInt32(&maxInFlight) < 2 {
		t.Fatalf("disjoint batches must overlap in flight, max observed %d", maxInFlight)
	}
}

// TestRunOrderedBatches_SharedKeysRunInOrder pins the concurrency contract
// alongside the ordering one: unrelated batches still overlap while the
// shared-key pair is serialised (asserted strictly in the NeverOverlap
// test below).
func TestRunOrderedBatches_SharedKeysRunInOrder(t *testing.T) {
	var inFlight, maxInFlight int32
	keys := []map[string]struct{}{
		{"shared": {}}, {"shared": {}, "other": {}}, {"unrelated": {}},
	}
	runOrderedBatches(context.Background(), 3, keys, func(ctx context.Context, idx int) batchOutcome {
		current := atomic.AddInt32(&inFlight, 1)
		for {
			observed := atomic.LoadInt32(&maxInFlight)
			if current <= observed || atomic.CompareAndSwapInt32(&maxInFlight, observed, current) {
				break
			}
		}
		time.Sleep(30 * time.Millisecond)
		atomic.AddInt32(&inFlight, -1)
		return batchOutcome{lastConsumed: "seq", delivered: true}
	})
	if atomic.LoadInt32(&maxInFlight) < 2 {
		t.Fatalf("unrelated batch must still run concurrently, max observed %d", maxInFlight)
	}
}

// TestRunOrderedBatches_SharedKeysNeverOverlap asserts the stricter
// property directly: two batches sharing a key are never in flight at the
// same time.
func TestRunOrderedBatches_SharedKeysNeverOverlap(t *testing.T) {
	var inFlight int32
	var overlap int32
	keys := []map[string]struct{}{
		{"shared": {}}, {"shared": {}},
	}
	runOrderedBatches(context.Background(), 2, keys, func(ctx context.Context, idx int) batchOutcome {
		if atomic.AddInt32(&inFlight, 1) > 1 {
			atomic.StoreInt32(&overlap, 1)
		}
		time.Sleep(30 * time.Millisecond)
		atomic.AddInt32(&inFlight, -1)
		return batchOutcome{}
	})
	if atomic.LoadInt32(&overlap) != 0 {
		t.Fatal("batches sharing a record key must never overlap in flight")
	}
}

// TestPrefixOutcome_StopsAtFirstFailure pins the checkpoint rule: only the
// contiguous prefix of consumed batches advances the position; a failed
// batch suspends everything after it even if those batches succeeded.
func TestPrefixOutcome_StopsAtFirstFailure(t *testing.T) {
	outcomes := []batchOutcome{
		{lastConsumed: "10", delivered: true},
		{lastConsumed: "20"},
		{err: errors.New("blocked")},
		{lastConsumed: "40", delivered: true},
	}
	lastConsumed, delivered, discarded, reported, failure := prefixOutcome(outcomes)
	if failure == nil {
		t.Fatal("failure must surface")
	}
	if reported != 0 {
		t.Fatalf("no partial responses in the cycle, got %d reported", reported)
	}
	if lastConsumed != "20" {
		t.Fatalf("checkpoint must stop at the failure, got %q", lastConsumed)
	}
	if !delivered || discarded {
		t.Fatalf("flags wrong: delivered=%v discarded=%v", delivered, discarded)
	}
}

// TestPrefixOutcome_AllConsumed pins the healthy path: every consumed batch
// advances the position and discards are reported.
func TestPrefixOutcome_AllConsumed(t *testing.T) {
	outcomes := []batchOutcome{
		{lastConsumed: "10", delivered: true},
		{lastConsumed: "20", discarded: true},
		{lastConsumed: "30", delivered: true},
	}
	lastConsumed, delivered, discarded, _, failure := prefixOutcome(outcomes)
	if failure != nil {
		t.Fatalf("unexpected failure: %v", failure)
	}
	if lastConsumed != "30" || !delivered || !discarded {
		t.Fatalf("outcome wrong: %q delivered=%v discarded=%v", lastConsumed, delivered, discarded)
	}
}

// TestRunOrderedBatches_CancelledWaitFailsTheBatch pins the shutdown rule:
// a batch that never started because its key dependency was still running
// when the context closed must surface a failure, not a zero outcome that
// prefixOutcome would fold into the consumed prefix.
func TestRunOrderedBatches_CancelledWaitFailsTheBatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entered := make(chan struct{})
	keys := []map[string]struct{}{
		{"shared": {}}, {"shared": {}},
	}
	var outcomes []batchOutcome
	done := make(chan struct{})
	go func() {
		outcomes = runOrderedBatches(ctx, 2, keys, func(ctx context.Context, idx int) batchOutcome {
			if idx == 0 {
				close(entered)
				// The in-flight batch also observes the cancellation and
				// reports its own consumed position.
				<-ctx.Done()
				return batchOutcome{lastConsumed: "10", delivered: true}
			}
			return batchOutcome{lastConsumed: "20", delivered: true}
		})
		close(done)
	}()

	<-entered
	cancel()
	<-done

	if outcomes[1].err == nil {
		t.Fatal("a batch cancelled before start must carry the context error, not a zero outcome")
	}
	lastConsumed, _, _, _, failure := prefixOutcome(outcomes)
	if failure == nil {
		t.Fatal("prefixOutcome must surface the cancelled batch as a failure")
	}
	if lastConsumed != "10" {
		t.Fatalf("checkpoint must keep the prefix before the cancelled batch, got %q", lastConsumed)
	}
}

// TestKinesisRecordKeys_CollectsPartitionKeys verifies the key extraction
// the ordering scheduler depends on.
func TestKinesisRecordKeys_CollectsPartitionKeys(t *testing.T) {
	keys := kinesisRecordKeys([]invokers.KinesisRecord{
		{PartitionKey: "pk-1"}, {PartitionKey: "pk-2"}, {PartitionKey: "pk-1"},
	})
	if len(keys) != 2 {
		t.Fatalf("expected 2 distinct keys, got %d", len(keys))
	}
	if _, ok := keys["pk-1"]; !ok {
		t.Fatal("pk-1 missing")
	}
}

// TestDynamoDBRecordKeys_CanonicalisesItemKeys verifies that item keys
// with identical attribute maps canonicalise to the same scheduler key.
func TestDynamoDBRecordKeys_CanonicalisesItemKeys(t *testing.T) {
	record := func(id string) invokers.DynamoDBStreamRecord {
		return invokers.DynamoDBStreamRecord{
			Dynamodb: map[string]interface{}{
				"Keys": map[string]interface{}{"Id": map[string]interface{}{"N": id}},
			},
		}
	}
	keys := dynamoDBRecordKeys([]invokers.DynamoDBStreamRecord{record("1"), record("1"), record("2")})
	if len(keys) != 2 {
		t.Fatalf("expected 2 distinct item keys, got %d", len(keys))
	}
}

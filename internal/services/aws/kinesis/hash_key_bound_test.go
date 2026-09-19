package kinesis

import (
	"errors"
	"testing"
	"time"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// TestPutRecordExplicitHashKeyWindow pins the explicit-hash-key validation
// window: the model's HashKey pattern (optional "0" or a non-zero-leading
// decimal) governs the shape, and the shard hash-key space bounds the
// value — the space is the model's 128-bit domain, so 2^64 and the space
// top 2^128-1 are both routable, while a value past the top can address no
// shard and is rejected up front instead of failing to route at write
// time.
func TestPutRecordExplicitHashKeyWindow(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("ehk_window", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	valid := []string{"", "0", "1", "18446744073709551616", kinesisstore.MaxShardHashKey}
	for _, key := range valid {
		if _, err := svc.putRecordCore(reqCtx, PutRecordInput{
			StreamName:      "ehk_window",
			Data:            "data",
			HasData:         true,
			PartitionKey:    "pk",
			ExplicitHashKey: key,
		}); err != nil {
			t.Fatalf("put with explicit hash key %q: %v, want accepted", key, err)
		}
	}

	invalid := []string{
		"340282366920938463463374607431768211456", // 2^128 — first value past the space
		"0123", // leading zero
		"-1",   // negative
		"abc",  // non-decimal
	}
	for _, key := range invalid {
		_, err := svc.putRecordCore(reqCtx, PutRecordInput{
			StreamName:      "ehk_window",
			Data:            "data",
			HasData:         true,
			PartitionKey:    "pk",
			ExplicitHashKey: key,
		})
		if !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("put with explicit hash key %q: %v, want ErrInvalidArgument", key, err)
		}
	}
}

// A hash key past the space is an entry-shape violation: the per-entry
// failure channel carries only the write-time outcomes the ErrorCode member
// documents for it, so the batch rejects as a whole — the identity the
// single PutRecord path answers — and nothing from the batch is written.
func TestPutRecordsRejectsOutOfSpaceExplicitHashKey(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("ehk_batch", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	_, err := svc.putRecordsCore(reqCtx, PutRecordsInput{
		StreamName: "ehk_batch",
		Records: []interface{}{
			map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
			map[string]interface{}{"Data": "data", "PartitionKey": "pk", "ExplicitHashKey": "340282366920938463463374607431768211456"},
			map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
		},
	})
	requireAWSCode(t, "batch with one out-of-space explicit hash key", err, "InvalidArgumentException")

	records, _, rerr := store.GetRecords("ehk_batch", "shardId-000000000000", "", 10, false, time.Time{})
	if rerr != nil {
		t.Fatalf("read back: %v", rerr)
	}
	if len(records) != 0 {
		t.Fatalf("records after rejected batch: got %d, want none — the valid entries must not land", len(records))
	}
}

// TestSplitShardRejectsOutOfSpaceHashKey pins the split form of the bound:
// a NewStartingHashKey past the hash-key space is rejected up front.
func TestSplitShardRejectsOutOfSpaceHashKey(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("ehk_split", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	_, err := svc.splitShardCore(reqCtx, SplitShardInput{
		StreamName:         "ehk_split",
		ShardToSplit:       "shardId-000000000000",
		NewStartingHashKey: "340282366920938463463374607431768211456",
	})
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("split with an out-of-space new starting hash key: %v, want ErrInvalidArgument", err)
	}
}

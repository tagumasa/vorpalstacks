package kinesis

import (
	"math/big"
	"testing"
)

// narrowShardRange rewrites a stream's only shard range so most of the
// routing space is uncovered — a fixture for routing-failure pins.
func narrowShardRange(t *testing.T, store *KinesisStore, streamName, shardID string) {
	t.Helper()
	shard, err := store.GetShard(streamName, shardID)
	if err != nil {
		t.Fatalf("get shard: %v", err)
	}
	shard.HashKeyRange.StartingHashKey = "0"
	shard.HashKeyRange.EndingHashKey = "100"
	if err := store.PutShard(shard); err != nil {
		t.Fatalf("rewrite shard range: %v", err)
	}
}

// TestPutRecordSurfacesRoutingFailure pins the removal of the silent
// shard[0] fallback: when no shard covers the computed hash, the write
// fails with an error instead of telling the client a wrong ShardId.
func TestPutRecordSurfacesRoutingFailure(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("route_gap", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	narrowShardRange(t, store, "route_gap", "shardId-000000000000")

	if _, _, err := store.PutRecordWithShardSelection("route_gap", "pk", "data", ""); err == nil {
		t.Fatal("put with an uncovered routing space succeeded — the silent shard[0] fallback is back")
	}

	results, err := store.PutRecords("route_gap", []PutRecordRequest{{Data: "data", PartitionKey: "pk"}})
	if err != nil {
		t.Fatalf("put records: %v", err)
	}
	if len(results) != 1 || results[0].ErrorCode != "InternalFailure" {
		t.Fatalf("batch entry under an uncovered routing space: got %+v, want ErrorCode InternalFailure", results)
	}
	if results[0].ErrorMessage != InternalFailureEntryMessage {
		t.Fatalf("InternalFailure entry message: got %q, want the model's fixed %q", results[0].ErrorMessage, InternalFailureEntryMessage)
	}
}

// TestHashKeyEndOfSpaceRoutesToLastShard pins the inclusive top of the
// routing space: the explicit hash key at the space's maximum value
// selects the last shard.
func TestHashKeyEndOfSpaceRoutesToLastShard(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("route_end", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	record, shardID, err := store.PutRecordWithShardSelection("route_end", "pk", "data", MaxShardHashKey)
	if err != nil {
		t.Fatalf("put at the top of the hash-key space: %v", err)
	}
	if shardID != "shardId-000000000000" || record.SequenceNumber == "" {
		t.Fatalf("put at the top of the hash-key space: shard %q, want the last (only) shard", shardID)
	}
}

// TestHashPartitionKeySpansRoutingSpace pins the partition-key reduction:
// hashes land in the shard hash-key space, [0, 2^128-1] inclusive.
func TestHashPartitionKeySpansRoutingSpace(t *testing.T) {
	store := newLockTestStore(t)
	space := new(big.Int).Lsh(big.NewInt(1), 128)
	for _, pk := range []string{"", "a", "partition-key-3", "日本語キー"} {
		hash := store.hashPartitionKey(pk)
		if hash.Sign() < 0 || hash.Cmp(space) >= 0 {
			t.Fatalf("hash of %q = %s, want within [0, 2^128-1]", pk, hash.String())
		}
	}
}

// TestCreateTilesHashKeySpace pins creation-time range construction over
// the model's 128-bit domain: shards tile [0, 2^128-1] with inclusive
// ranges — evenly for power-of-two counts, and AWS's own single-shard
// describe-stream output reports exactly [0, 2^128-1].
func TestCreateTilesHashKeySpace(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("route_tile", 2, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	shards, err := store.ListShards("route_tile", nil, "", 0)
	if err != nil {
		t.Fatalf("list shards: %v", err)
	}
	if len(shards) != 2 {
		t.Fatalf("shards: got %d, want 2", len(shards))
	}
	midpoint := new(big.Int).Lsh(big.NewInt(1), 127)
	if shards[0].HashKeyRange.StartingHashKey != "0" ||
		shards[0].HashKeyRange.EndingHashKey != new(big.Int).Sub(midpoint, big.NewInt(1)).String() {
		t.Fatalf("first shard range: got [%s, %s], want [0, 2^127-1]",
			shards[0].HashKeyRange.StartingHashKey, shards[0].HashKeyRange.EndingHashKey)
	}
	if shards[1].HashKeyRange.StartingHashKey != midpoint.String() ||
		shards[1].HashKeyRange.EndingHashKey != MaxShardHashKey {
		t.Fatalf("second shard range: got [%s, %s], want [2^127, 2^128-1]",
			shards[1].HashKeyRange.StartingHashKey, shards[1].HashKeyRange.EndingHashKey)
	}
}

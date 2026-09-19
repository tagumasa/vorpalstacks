package kinesis

import (
	"testing"
)

// TestNilRangeShardReadsAsOpen pins the absent-submessage reading at the
// three sites that dereferenced a shard's range without a guard: a split
// of a range-less shard record proceeds (the shard reads as open, and
// sealing it materialises the range), closeShard on a range-less shard
// initialises the range it will carry, and the split's children carry
// well-formed ranges.
func TestNilRangeShardReadsAsOpen(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("nil_range", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	// Rewrite the only shard's record without its range submessage — the
	// pathological record the guards must survive.
	shard, err := store.GetShard("nil_range", "shardId-000000000000")
	if err != nil {
		t.Fatalf("get shard: %v", err)
	}
	shard.SequenceNumberRange = nil
	if err := store.PutShard(shard); err != nil {
		t.Fatalf("rewrite shard without range: %v", err)
	}
	reread, err := store.GetShard("nil_range", "shardId-000000000000")
	if err != nil {
		t.Fatalf("reread shard: %v", err)
	}
	if reread.SequenceNumberRange != nil {
		t.Fatal("fixture check: the rewritten shard still carries a range")
	}

	// closeShard on a range-less shard materialises the range instead of
	// panicking.
	bare := &Shard{ShardID: "shardId-000000000001", StreamName: "nil_range"}
	closeShard(bare)
	if bare.SequenceNumberRange == nil {
		t.Fatal("closeShard left the range absent")
	}

	// The split reads the range-less shard as open and completes; every
	// shard record carries a range afterwards — the sealed parent through
	// closeShard's materialisation, the children by construction.
	if err := store.SplitShard("nil_range", "shardId-000000000000", ""); err != nil {
		t.Fatalf("split a range-less shard: %v", err)
	}
	shards, err := store.ListShards("nil_range", nil, "", 0)
	if err != nil {
		t.Fatalf("list shards after split: %v", err)
	}
	if len(shards) < 3 {
		t.Fatalf("shards after split: got %d, want the sealed parent plus two children", len(shards))
	}
	for _, got := range shards {
		if got.SequenceNumberRange == nil {
			t.Fatalf("shard %s carries no range after the split", got.ShardID)
		}
	}
}

// TestHashKeyRangelessShardRoutesToError pins the absent-submessage reading
// on the write path: a shard record without its HashKeyRange submessage
// stays routable (activeShardsOf filters on the sequence-number range
// alone), and the router answers the unparseable-bound error instead of
// panicking — the same empty-string-bounds reading the guards and
// formatters take.
func TestHashKeyRangelessShardRoutesToError(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("nil_hash", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	shard, err := store.GetShard("nil_hash", "shardId-000000000000")
	if err != nil {
		t.Fatalf("get shard: %v", err)
	}
	shard.HashKeyRange = nil
	if err := store.PutShard(shard); err != nil {
		t.Fatalf("rewrite shard without hash-key range: %v", err)
	}

	// The partition-key route hashes into the shard set; the hash-range-less
	// record must surface as the routing error, never as a panic.
	_, _, err = store.PutRecordWithShardSelection("nil_hash", "pk", "data", "")
	if err == nil {
		t.Fatal("put through a hash-range-less shard: succeeded, want the routing error")
	}
}

package kinesis

import (
	"math/big"
	"testing"
)

// TestHashKeyLessIsNumeric pins the ordering primitive: hash key strings
// span 1 to 20 digits, so lexicographic order diverges from numeric order
// exactly on the initial shard layout of a four-shard stream.
func TestHashKeyLessIsNumeric(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"0", "4611686018427387903", true},
		// 9223372036854775806 < 13835058055282163709 numerically, though
		// the strings sort the other way ("9..." > "1..." lexicographically).
		{"9223372036854775806", "13835058055282163709", true},
		{"13835058055282163709", "18446744073709551615", true},
		{"42", "42", false},
	}
	for _, tc := range cases {
		if got := hashKeyLess(tc.a, tc.b); got != tc.want {
			t.Fatalf("hashKeyLess(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

// TestUpdateShardCountScaleDownConverges pins the scale-down convergence:
// merging four initial shards to two must reach the target — the merge loop
// orders candidates numerically, so adjacency in the list is adjacency in
// the key space — and the surviving open shards must tile the whole hash
// key space contiguously.
func TestUpdateShardCountScaleDownConverges(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("scale_down", 4, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	if err := store.UpdateShardCount("scale_down", 2); err != nil {
		t.Fatalf("update shard count 4 to 2: %v", err)
	}

	stream, err := store.GetStream("scale_down")
	if err != nil {
		t.Fatalf("get stream: %v", err)
	}
	if stream.ShardCount != 2 {
		t.Fatalf("stream shard count after scale-down: got %d, want 2 — the update must converge, not no-op", stream.ShardCount)
	}

	shards, err := store.ListShards("scale_down", nil, "", 0)
	if err != nil {
		t.Fatalf("list shards: %v", err)
	}
	var open []*Shard
	for _, shard := range shards {
		if shard.SequenceNumberRange.EndingSequenceNumber == "" {
			open = append(open, shard)
		}
	}
	if len(open) != 2 {
		t.Fatalf("open shards after scale-down: got %d, want 2", len(open))
	}

	// The open shards tile the hash-key space [0, 2^128-1] with contiguous
	// inclusive ranges.
	sorted := make([]*Shard, len(open))
	copy(sorted, open)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && hashKeyLess(sorted[j].HashKeyRange.StartingHashKey, sorted[j-1].HashKeyRange.StartingHashKey); j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	spaceTop := MaxShardHashKeyInt()
	if start, _ := new(big.Int).SetString(sorted[0].HashKeyRange.StartingHashKey, 10); start.Sign() != 0 {
		t.Fatalf("key space tiling: lowest open shard starts at %s, want 0", sorted[0].HashKeyRange.StartingHashKey)
	}
	if end, _ := new(big.Int).SetString(sorted[len(sorted)-1].HashKeyRange.EndingHashKey, 10); end.Cmp(spaceTop) != 0 {
		t.Fatalf("key space tiling: highest open shard ends at %s, want %s", sorted[len(sorted)-1].HashKeyRange.EndingHashKey, spaceTop.String())
	}
	for i := 1; i < len(sorted); i++ {
		prevEnd, _ := new(big.Int).SetString(sorted[i-1].HashKeyRange.EndingHashKey, 10)
		start, _ := new(big.Int).SetString(sorted[i].HashKeyRange.StartingHashKey, 10)
		if new(big.Int).Sub(start, prevEnd).Cmp(big.NewInt(1)) != 0 {
			t.Fatalf("key space tiling: gap between shard %s (ends %s) and shard %s (starts %s)",
				sorted[i-1].ShardID, sorted[i-1].HashKeyRange.EndingHashKey, sorted[i].ShardID, sorted[i].HashKeyRange.StartingHashKey)
		}
	}
}

// TestUpdateShardCountScaleUpSplitsEachOriginalOnce pins the scale-up
// generation rule: every original open shard splits before any child does,
// so a two-to-four scaling closes exactly the two originals and opens four
// children of the one generation — the picker must not wander onto a child
// it just created and close it again inside the same operation.
func TestUpdateShardCountScaleUpSplitsEachOriginalOnce(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("scale_up", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	if err := store.UpdateShardCount("scale_up", 2); err != nil {
		t.Fatalf("update shard count 1 to 2: %v", err)
	}
	if err := store.UpdateShardCount("scale_up", 4); err != nil {
		t.Fatalf("update shard count 2 to 4: %v", err)
	}

	shards, err := store.ListShards("scale_up", nil, "", 0)
	if err != nil {
		t.Fatalf("list shards: %v", err)
	}

	// The first split's children are the second scaling's originals: both
	// close, and their four children alone are open.
	openByParent := map[string]int{}
	for _, shard := range shards {
		if shard.SequenceNumberRange.EndingSequenceNumber == "" {
			openByParent[shard.ParentShardID]++
		}
	}
	if len(openByParent) != 2 || openByParent["shardId-000000000001"] != 2 || openByParent["shardId-000000000002"] != 2 {
		t.Fatalf("open shards grouped by parent: got %v, want two children each under shards 001 and 002", openByParent)
	}
	for _, shard := range shards {
		switch shard.ShardID {
		case "shardId-000000000001", "shardId-000000000002":
			if shard.SequenceNumberRange.EndingSequenceNumber == "" {
				t.Fatalf("shard %s: still open after the second scaling, want closed", shard.ShardID)
			}
		case "shardId-000000000005", "shardId-000000000006":
			if shard.ParentShardID != "shardId-000000000002" {
				t.Fatalf("shard %s: parent %s, want shardId-000000000002", shard.ShardID, shard.ParentShardID)
			}
		case "shardId-000000000003", "shardId-000000000004":
			if shard.ParentShardID != "shardId-000000000001" {
				t.Fatalf("shard %s: parent %s, want shardId-000000000001", shard.ShardID, shard.ParentShardID)
			}
		}
	}

	// The widest-range-first order also yields uniform quarters: every open
	// shard spans a quarter of the hash-key space.
	quarter := new(big.Int).Add(MaxShardHashKeyInt(), big.NewInt(1))
	quarter.Div(quarter, big.NewInt(4))
	for _, shard := range shards {
		if shard.SequenceNumberRange.EndingSequenceNumber != "" {
			continue
		}
		if width := hashRangeWidth(shard); width.Cmp(quarter) != 0 {
			t.Fatalf("open shard %s width %s, want the uniform quarter %s", shard.ShardID, width.String(), quarter.String())
		}
	}
}

// TestCreateStreamRejectsSubOneShardCount pins the store's own guard: a
// caller passing a count below one gets the invalid-count error, never the
// hash-space division a zero count would panic on.
func TestCreateStreamRejectsSubOneShardCount(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("zero_count", 0, StreamModeProvisioned, 0, 0, nil); err != ErrInvalidShardCount {
		t.Fatalf("zero shard count: got %v, want ErrInvalidShardCount", err)
	}
}

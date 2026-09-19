package apps

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
	svckinesis "vorpalstacks/internal/services/aws/kinesis"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

func newKinesisInvokerTestService(t *testing.T) *svckinesis.KinesisService {
	t.Helper()
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("region storage manager: %v", err)
	}
	svc := svckinesis.NewKinesisService("000000000000")
	svc.SetStorageManager(sm)
	return svc
}

// TestKinesisInvokerRoutesByPartitionKey pins the invoker seam's placement
// contract: shard selection hashes the partition key exactly as the public
// PutRecord does, so the same key keeps its shard across writes (the
// retry-stability guarantee EventBridge's dispatcher documents) and the
// placement spreads across an evenly split stream instead of piling onto
// the first open shard.
func TestKinesisInvokerRoutesByPartitionKey(t *testing.T) {
	svc := newKinesisInvokerTestService(t)
	adapter := &kinesisInvokerAdapter{provider: svc}
	store, err := svc.GetStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("store for region: %v", err)
	}
	if _, err := store.CreateStream("route", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	if err := store.UpdateShardCount("route", 2); err != nil {
		t.Fatalf("split stream: %v", err)
	}

	ctx := context.Background()
	shardOf := func(partitionKey string) string {
		t.Helper()
		_, reportedShard, err := adapter.PutRecord(ctx, "us-east-1", "route", partitionKey, []byte("data"))
		if err != nil {
			t.Fatalf("put %q: %v", partitionKey, err)
		}
		shards, err := store.ListShards("route", &kinesisstore.ShardFilter{Type: "AT_LATEST"}, "", 0)
		if err != nil {
			t.Fatalf("list open shards: %v", err)
		}
		if len(shards) != 2 {
			t.Fatalf("open shards: %d, want 2", len(shards))
		}
		for _, shard := range shards {
			it, err := store.CreateShardIterator("route", shard.ShardID, "TRIM_HORIZON", "", nil)
			if err != nil {
				t.Fatalf("iterator %s: %v", shard.ShardID, err)
			}
			records, _, err := store.GetRecords("route", shard.ShardID, it.SequenceNumber, 10000, false, shard.CreatedAt)
			if err != nil {
				t.Fatalf("get records %s: %v", shard.ShardID, err)
			}
			for _, record := range records {
				if record.PartitionKey == partitionKey {
					if reportedShard != shard.ShardID {
						t.Fatalf("key %q: PutRecord reported shard %s, the record landed on %s — the contract's ShardId return is not the receiving shard",
							partitionKey, reportedShard, shard.ShardID)
					}
					return shard.ShardID
				}
			}
		}
		t.Fatalf("record for key %q not found on any open shard", partitionKey)
		return ""
	}

	keys := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel"}
	placement := make(map[string]string, len(keys))
	shardsSeen := make(map[string]bool)
	for _, key := range keys {
		shard := shardOf(key)
		placement[key] = shard
		shardsSeen[shard] = true
	}
	if len(shardsSeen) < 2 {
		t.Fatalf("all keys routed to one shard (%v) — placement ignores the partition key", placement)
	}

	// Retry stability: the same key keeps its shard.
	for _, key := range keys {
		if again := shardOf(key); again != placement[key] {
			t.Fatalf("key %q moved from %s to %s across retries", key, placement[key], again)
		}
	}
}

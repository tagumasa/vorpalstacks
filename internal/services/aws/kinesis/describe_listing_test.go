package kinesis

import (
	"testing"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// TestAccountSettingsRoundTrip pins the persistence contract: Describe
// reports what Update last answered — DISABLED before any update, the
// updated status after — instead of an echo stub that discards the write.
func TestAccountSettingsRoundTrip(t *testing.T) {
	svc, _, reqCtx := newListPageTestEnv(t)

	result, err := svc.describeAccountSettingsCore(reqCtx)
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	if result.Status != kinesisstore.MTBCStatusDisabled {
		t.Fatalf("fresh account reports %q, want DISABLED", result.Status)
	}

	if _, err := svc.updateAccountSettingsCore(reqCtx, UpdateAccountSettingsInput{Status: kinesisstore.MTBCStatusEnabled}); err != nil {
		t.Fatalf("update ENABLED: %v", err)
	}
	result, err = svc.describeAccountSettingsCore(reqCtx)
	if err != nil {
		t.Fatalf("describe after update: %v", err)
	}
	if result.Status != kinesisstore.MTBCStatusEnabled {
		t.Fatalf("describe after update reports %q, want ENABLED", result.Status)
	}

	_, err = svc.updateAccountSettingsCore(reqCtx, UpdateAccountSettingsInput{Status: "PAUSED"})
	requireAWSCode(t, "invalid commitment status", err, "InvalidArgumentException")

	if _, err := svc.updateAccountSettingsCore(reqCtx, UpdateAccountSettingsInput{Status: kinesisstore.MTBCStatusDisabled}); err != nil {
		t.Fatalf("update DISABLED: %v", err)
	}
	result, err = svc.describeAccountSettingsCore(reqCtx)
	if err != nil {
		t.Fatalf("describe after reset: %v", err)
	}
	if result.Status != kinesisstore.MTBCStatusDisabled {
		t.Fatalf("describe after reset reports %q, want DISABLED", result.Status)
	}
}

// TestGetResourcePolicyErrorPaths pins the two failure identities: a
// resource with no policy answers the operation's declared
// ResourceNotFoundException, and a malformed ARN is rejected before
// storage — never an empty-string success.
func TestGetResourcePolicyErrorPaths(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("pol", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}

	_, err = svc.getResourcePolicyCore(reqCtx, ResourcePolicyInput{ResourceARN: stream.StreamARN})
	requireAWSCode(t, "stream without policy", err, "ResourceNotFoundException")

	if err := store.PutResourcePolicy(stream.StreamARN, `{"Version":"2012-10-17"}`); err != nil {
		t.Fatalf("put policy: %v", err)
	}
	policy, err := svc.getResourcePolicyCore(reqCtx, ResourcePolicyInput{ResourceARN: stream.StreamARN})
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	if policy == "" {
		t.Fatalf("stored policy reads back empty")
	}

	_, err = svc.getResourcePolicyCore(reqCtx, ResourcePolicyInput{ResourceARN: stream.StreamARN + "-extra"})
	requireAWSCode(t, "unknown resource", err, "ResourceNotFoundException")

	_, err = svc.getResourcePolicyCore(reqCtx, ResourcePolicyInput{ResourceARN: "not-an-arn"})
	requireAWSCode(t, "malformed ARN", err, "InvalidArgumentException")
}

// TestDescribeStreamPagination pins the shard-page members: Limit serves
// at most one hundred shards per the member's documentation (its range
// trait still rejects out-of-window values), ExclusiveStartShardId
// resumes the page after that shard, and HasMoreShards reports whether a
// further page follows — the exact final page says false.
func TestDescribeStreamPagination(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("dsp", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	// Scale 1→2→4 (the store's ratio window), leaving four open shards.
	for _, target := range []int32{2, 4} {
		if err := store.UpdateShardCount("dsp", target); err != nil {
			t.Fatalf("update shard count to %d: %v", target, err)
		}
	}

	// DescribeStream serves the full shard history — closed parents
	// included — so 1→2→4 leaves seven shards.
	full, err := svc.describeStreamCore(reqCtx, DescribeStreamInput{StreamName: "dsp"})
	if err != nil {
		t.Fatalf("full describe: %v", err)
	}
	total := len(full.Shards)
	if total != 7 || full.HasMoreShards {
		t.Fatalf("full describe: %d shards, HasMoreShards %v", total, full.HasMoreShards)
	}

	// Walk pages of two: every shard served exactly once, the final page
	// alone reporting HasMoreShards false.
	var walked []string
	in := DescribeStreamInput{StreamName: "dsp", Limit: 2, HasLimit: true}
	pages := 0
	for {
		page, err := svc.describeStreamCore(reqCtx, in)
		if err != nil {
			t.Fatalf("page %d: %v", pages+1, err)
		}
		for _, shard := range page.Shards {
			walked = append(walked, shard.ShardID)
		}
		pages++
		if !page.HasMoreShards {
			if len(page.Shards) == 0 {
				t.Fatalf("walk ended on an empty page before serving every shard")
			}
			break
		}
		if len(page.Shards) != 2 {
			t.Fatalf("page %d served %d shards, want 2", pages, len(page.Shards))
		}
		in.ExclusiveStartShardId = page.Shards[len(page.Shards)-1].ShardID
	}
	if len(walked) != total {
		t.Fatalf("paged walk served %d shards, want %d", len(walked), total)
	}
	seen := make(map[string]bool, len(walked))
	for _, id := range walked {
		if seen[id] {
			t.Fatalf("shard %s served twice", id)
		}
		seen[id] = true
	}

	// A Limit above the documented hundred still serves every shard (the
	// range trait alone accepts it); a value outside the trait's window
	// rejects.
	capped, err := svc.describeStreamCore(reqCtx, DescribeStreamInput{StreamName: "dsp", Limit: 150, HasLimit: true})
	if err != nil {
		t.Fatalf("limit 150: %v", err)
	}
	if len(capped.Shards) != total || capped.HasMoreShards {
		t.Fatalf("limit 150: %d shards, HasMoreShards %v", len(capped.Shards), capped.HasMoreShards)
	}

	for _, bad := range []int{0, kinesisstore.MaxListResultsLimit + 1} {
		_, err := svc.describeStreamCore(reqCtx, DescribeStreamInput{StreamName: "dsp", Limit: bad, HasLimit: true})
		requireAWSCode(t, "limit out of window", err, "InvalidArgumentException")
	}
}

// TestFormatShardsSurvivesAbsentSubmessages pins the emitted shard shape
// against a pathological record: a shard without its range or hash-key
// submessage formats with empty-string bounds instead of panicking — the
// same absent-submessage reading the store's own guards take.
func TestFormatShardsSurvivesAbsentSubmessages(t *testing.T) {
	out := formatShards([]*kinesisstore.Shard{{ShardID: "shardId-000000000000"}})
	if len(out) != 1 {
		t.Fatalf("formatted %d shards, want 1", len(out))
	}
	hashRange, ok := out[0]["HashKeyRange"].(map[string]interface{})
	if !ok || hashRange["StartingHashKey"] != "" || hashRange["EndingHashKey"] != "" {
		t.Fatalf("hash key range of a submessage-less shard: got %+v, want empty-string bounds", out[0]["HashKeyRange"])
	}
	seqRange, ok := out[0]["SequenceNumberRange"].(map[string]interface{})
	if !ok || seqRange["StartingSequenceNumber"] != "" {
		t.Fatalf("sequence number range of a submessage-less shard: got %+v, want an empty start", out[0]["SequenceNumberRange"])
	}
	if _, hasEnd := seqRange["EndingSequenceNumber"]; hasEnd {
		t.Fatal("sequence number range of a submessage-less shard carries an ending number")
	}
}

// TestFormatChildShardsSurvivesAbsentSubmessages pins the ChildShards
// emission against the same pathological record: a child shard without its
// hash-key submessage formats with empty-string bounds instead of panicking
// — HashKeyRange is a required output member, so the emission must stay
// well-formed.
func TestFormatChildShardsSurvivesAbsentSubmessages(t *testing.T) {
	out := formatChildShards([]*kinesisstore.Shard{{ShardID: "shardId-000000000001"}})
	if len(out) != 1 {
		t.Fatalf("formatted %d child shards, want 1", len(out))
	}
	entry, ok := out[0].(map[string]interface{})
	if !ok {
		t.Fatalf("child shard entry: got %T, want a map", out[0])
	}
	hashRange, ok := entry["HashKeyRange"].(map[string]interface{})
	if !ok || hashRange["StartingHashKey"] != "" || hashRange["EndingHashKey"] != "" {
		t.Fatalf("hash key range of a submessage-less child shard: got %+v, want empty-string bounds", entry["HashKeyRange"])
	}
	if parents, ok := entry["ParentShards"].([]string); !ok || len(parents) != 0 {
		t.Fatalf("parent shards of a submessage-less child shard: got %#v, want an empty list", entry["ParentShards"])
	}
}

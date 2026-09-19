package kinesis

import (
	"strconv"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// newListShardsFixture opens a throwaway store with one provisioned stream
// whose initial shard has been split: a closed parent and two open
// children, the lifecycle pair every filter-semantics row needs. The
// parent's close timestamp is the arrival time encoded in its ending
// sequence number, which the horizon-relative semantics read; overlapTime
// is captured while the parent is still open, so the timestamp filters have
// a moment the parent's lifecycle spans.
func newListShardsFixture(t *testing.T) (*KinesisService, *kinesisstore.KinesisStore, *request.RequestContext, string, time.Time) {
	t.Helper()
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("lsh", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create fixture stream: %v", err)
	}
	// One record after the overlap moment gives the parent a last-written
	// position, so its closing sequence number (and the end timestamp the
	// filters read from it) lands after the overlap while its creation
	// precedes it.
	overlapTime := time.Now().UTC()
	if _, _, err := store.PutRecordWithShardSelection("lsh", "pk", "ZGF0YQ==", ""); err != nil {
		t.Fatalf("seed fixture record: %v", err)
	}
	if _, err := svc.splitShardCore(reqCtx, SplitShardInput{
		StreamName: "lsh", ShardToSplit: "shardId-000000000000",
		NewStartingHashKey: "9223372036854775808",
	}); err != nil {
		t.Fatalf("split fixture shard: %v", err)
	}
	return svc, store, reqCtx, "lsh", overlapTime
}

// shardIDsOf collects a result page's shard IDs in listed order.
func shardIDsOf(t *testing.T, result ListShardsResult) []string {
	t.Helper()
	ids := make([]string, 0, len(result.Shards))
	for _, shard := range result.Shards {
		ids = append(ids, shard.ShardID)
	}
	return ids
}

// TestListShardsFilterSemantics pins every ShardFilterType's documented
// membership rule against the closed-parent/open-children fixture, the
// absent-member default, and the enum and member-pairing rejections.
func TestListShardsFilterSemantics(t *testing.T) {
	svc, _, reqCtx, streamName, overlap := newListShardsFixture(t)
	future := time.Now().UTC().Add(time.Hour)

	cases := []struct {
		name    string
		filter  *kinesisstore.ShardFilter
		want    []string
		wantErr string
	}{
		{name: "absent member is the FROM_TRIM_HORIZON default", want: []string{
			"shardId-000000000000", "shardId-000000000001", "shardId-000000000002",
		}},
		{name: "FROM_TRIM_HORIZON carries trim to tip", filter: &kinesisstore.ShardFilter{Type: "FROM_TRIM_HORIZON"}, want: []string{
			"shardId-000000000000", "shardId-000000000001", "shardId-000000000002",
		}},
		{name: "AT_LATEST is the open set alone", filter: &kinesisstore.ShardFilter{Type: "AT_LATEST"}, want: []string{
			"shardId-000000000001", "shardId-000000000002",
		}},
		// Nothing existed at the trim horizon (a day back), so the set
		// open at it is empty — the contrast with FROM_TRIM_HORIZON's
		// trim-to-tip walk.
		{name: "AT_TRIM_HORIZON is the set open at the horizon", filter: &kinesisstore.ShardFilter{Type: "AT_TRIM_HORIZON"}, want: []string{}},
		{name: "AFTER_SHARD_ID starts past the given ID", filter: &kinesisstore.ShardFilter{
			Type: "AFTER_SHARD_ID", ShardID: "shardId-000000000000",
		}, want: []string{"shardId-000000000001", "shardId-000000000002"}},
		// The overlap moment sits inside the parent's lifecycle alone:
		// the parent had started and had not yet closed, while the
		// children did not exist.
		{name: "AT_TIMESTAMP inside the parent's lifecycle is the parent alone", filter: &kinesisstore.ShardFilter{
			Type: "AT_TIMESTAMP", Timestamp: &overlap,
		}, want: []string{"shardId-000000000000"}},
		// A future timestamp matches the open set: every open shard
		// started before it and remains open, while the closed parent's
		// end timestamp precedes it.
		{name: "AT_TIMESTAMP in the future keeps the open set alone", filter: &kinesisstore.ShardFilter{
			Type: "AT_TIMESTAMP", Timestamp: &future,
		}, want: []string{"shardId-000000000001", "shardId-000000000002"}},
		{name: "FROM_TIMESTAMP keeps the closed-past-the-moment shards and all open", filter: &kinesisstore.ShardFilter{
			Type: "FROM_TIMESTAMP", Timestamp: &overlap,
		}, want: []string{
			"shardId-000000000000", "shardId-000000000001", "shardId-000000000002",
		}},
		{name: "FROM_TIMESTAMP before the horizon is corrected to it", filter: &kinesisstore.ShardFilter{
			Type: "FROM_TIMESTAMP", Timestamp: &time.Time{},
		}, want: []string{
			"shardId-000000000000", "shardId-000000000001", "shardId-000000000002",
		}},

		{name: "a Type outside the enum rejects", filter: &kinesisstore.ShardFilter{Type: "FROM_LATEST"}, wantErr: "InvalidArgumentException"},
		{name: "an empty Type rejects", filter: &kinesisstore.ShardFilter{}, wantErr: "InvalidArgumentException"},
		{name: "AFTER_SHARD_ID without its ShardId rejects", filter: &kinesisstore.ShardFilter{Type: "AFTER_SHARD_ID"}, wantErr: "InvalidArgumentException"},
		{name: "a ShardId outside AFTER_SHARD_ID rejects", filter: &kinesisstore.ShardFilter{
			Type: "AT_LATEST", ShardID: "shardId-000000000000",
		}, wantErr: "InvalidArgumentException"},
		{name: "AT_TIMESTAMP without its Timestamp rejects", filter: &kinesisstore.ShardFilter{Type: "AT_TIMESTAMP"}, wantErr: "InvalidArgumentException"},
		{name: "a Timestamp outside the timestamp types rejects", filter: &kinesisstore.ShardFilter{
			Type: "AT_LATEST", Timestamp: &future,
		}, wantErr: "InvalidArgumentException"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.listShardsCore(reqCtx, ListShardsInput{
				StreamName: streamName, ShardFilter: tc.filter,
			})
			if tc.wantErr != "" {
				requireAWSCode(t, tc.name, err, tc.wantErr)
				return
			}
			if err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			got := shardIDsOf(t, result)
			if len(got) != len(tc.want) {
				t.Fatalf("%s: got %v, want %v", tc.name, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("%s: got %v, want %v", tc.name, got, tc.want)
				}
			}
		})
	}
}

// TestListShardsPaginationBoundary pins the over-fetch page contract: the
// page that exactly fills MaxResults as the final page carries no
// continuation token, and a token-only follow-up lands on the remaining
// shards.
func TestListShardsPaginationBoundary(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("lspage", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create fixture stream: %v", err)
	}
	// The scaling operation's ratio window is [0.5, 2.0]: reach four
	// open shards in two doublings.
	if err := store.UpdateShardCount("lspage", 2); err != nil {
		t.Fatalf("reshape fixture to two open shards: %v", err)
	}
	if err := store.UpdateShardCount("lspage", 4); err != nil {
		t.Fatalf("reshape fixture to four open shards: %v", err)
	}
	open := &kinesisstore.ShardFilter{Type: "AT_LATEST"}

	exact, err := svc.listShardsCore(reqCtx, ListShardsInput{
		StreamName: "lspage", ShardFilter: open, MaxResults: 4, HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("list exactly the open set: %v", err)
	}
	if len(exact.Shards) != 4 || exact.NextToken != "" {
		t.Fatalf("exact final page: got %d shards token=%q, want 4 without a token", len(exact.Shards), exact.NextToken)
	}

	first, err := svc.listShardsCore(reqCtx, ListShardsInput{
		StreamName: "lspage", ShardFilter: open, MaxResults: 2, HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	if len(first.Shards) != 2 || first.NextToken == "" {
		t.Fatalf("first page: got %d shards, want 2 with a token", len(first.Shards))
	}

	// The follow-up carries the identification in the envelope alone.
	second, err := svc.listShardsCore(reqCtx, ListShardsInput{
		NextToken: first.NextToken, ShardFilter: open, MaxResults: 2, HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("token follow-up: %v", err)
	}
	if len(second.Shards) != 2 || second.NextToken != "" {
		t.Fatalf("exact-boundary follow-up: got %d shards token=%q, want 2 without a token", len(second.Shards), second.NextToken)
	}
	if got := shardIDsOf(t, second); got[0] != "shardId-000000000005" || got[1] != "shardId-000000000006" {
		t.Fatalf("token follow-up IDs: got %v, want the second half", got)
	}
}

// TestListShardsTokenContract pins the continuation envelope: it binds the
// stream (name and generation), excludes the explicit resumption member,
// and answers the declared errors for a malformed envelope, a foreign
// stream, and a recreated generation.
func TestListShardsTokenContract(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	streamFixture, err := store.CreateStream("lstok", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create fixture stream: %v", err)
	}
	if err := store.UpdateShardCount("lstok", 2); err != nil {
		t.Fatalf("reshape fixture to two open shards: %v", err)
	}
	if err := store.UpdateShardCount("lstok", 4); err != nil {
		t.Fatalf("reshape fixture to four open shards: %v", err)
	}
	open := &kinesisstore.ShardFilter{Type: "AT_LATEST"}

	first, err := svc.listShardsCore(reqCtx, ListShardsInput{
		StreamName: "lstok", ShardFilter: open, MaxResults: 2, HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	token := first.NextToken

	if _, err := store.CreateStream("lsother", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create second stream: %v", err)
	}

	requireAWSCode(t, "token with a foreign StreamName", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{NextToken: token, StreamName: "lsother"})
		return err
	}(), "InvalidArgumentException")

	// The identification members cannot travel with the token even when
	// they agree with it — "the latter unambiguously identifies the stream".
	requireAWSCode(t, "token with its agreeing StreamName", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{NextToken: token, StreamName: "lstok"})
		return err
	}(), "InvalidArgumentException")

	requireAWSCode(t, "token with an agreeing StreamCreationTimestamp", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{
			NextToken:               token,
			StreamCreationTimestamp: strconv.FormatInt(streamFixture.CreatedAt.Unix(), 10),
		})
		return err
	}(), "InvalidArgumentException")

	requireAWSCode(t, "token with the explicit resumption member", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{
			NextToken: token, ExclusiveStartShardId: "shardId-000000000001",
		})
		return err
	}(), "InvalidArgumentException")

	// An envelope this service never issued — base64-decodable but not the
	// three-field form — is an expired token, never a silent restart.
	requireAWSCode(t, "malformed envelope", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{
			NextToken: "bm90LWEtdG9rZW4=",
		})
		return err
	}(), "ExpiredNextTokenException")
	// A structurally valid envelope naming a stream that never existed
	// answers the stream's own error.
	requireAWSCode(t, "envelope naming a nonexistent stream", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{
			NextToken: encodeListShardsToken("never-created", 0, time.Now().UnixNano(), "shardId-000000000000"),
		})
		return err
	}(), "ResourceNotFoundException")

	// The documented validity window: an envelope issued more than the
	// documented lifetime ago is expired whatever it names.
	requireAWSCode(t, "envelope past its validity window", func() error {
		old := time.Now().Add(-kinesisstore.NextTokenValidity - time.Minute).UnixNano()
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{
			NextToken: encodeListShardsToken("lstok", streamFixture.CreatedAt.Unix(), old, "shardId-000000000000"),
		})
		return err
	}(), "ExpiredNextTokenException")

	// The exclusive member alone resumes past the given shard.
	past, err := svc.listShardsCore(reqCtx, ListShardsInput{
		StreamName: "lstok", ShardFilter: open,
		ExclusiveStartShardId: "shardId-000000000002",
	})
	if err != nil {
		t.Fatalf("exclusive start alone: %v", err)
	}
	if got := shardIDsOf(t, past); len(got) != 4 || got[0] != "shardId-000000000003" || got[3] != "shardId-000000000006" {
		t.Fatalf("exclusive start IDs: got %v, want the four shards past the given ID", got)
	}

	// A token against a deleted stream answers the stream's own error; a
	// token against a recreated generation is expired — the creation
	// timestamp the envelope carries identifies its own generation.
	if err := store.DeleteStream("lstok"); err != nil {
		t.Fatalf("delete fixture stream: %v", err)
	}
	requireAWSCode(t, "token against a deleted stream", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{NextToken: token})
		return err
	}(), "ResourceNotFoundException")

	time.Sleep(1100 * time.Millisecond)
	if _, err := store.CreateStream("lstok", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("recreate fixture stream: %v", err)
	}
	requireAWSCode(t, "token against a recreated generation", func() error {
		_, err := svc.listShardsCore(reqCtx, ListShardsInput{NextToken: token})
		return err
	}(), "ExpiredNextTokenException")
}

// TestListShardsTokenAgreesWithARN pins the identity agreement on the ARN
// side: a request carrying a NextToken and another stream's StreamARN
// answers invalid input — the token's stream and the addressed stream
// must agree, exactly as the StreamName form does.
func TestListShardsTokenAgreesWithARN(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	streamA, err := store.CreateStream("shardtok_a", 2, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream A: %v", err)
	}
	streamB, err := store.CreateStream("shardtok_b", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream B: %v", err)
	}

	first, err := svc.listShardsCore(reqCtx, ListShardsInput{
		StreamARN: streamA.StreamARN, MaxResults: 1, HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("first page: %v", err)
	}

	_, err = svc.listShardsCore(reqCtx, ListShardsInput{
		StreamARN: streamB.StreamARN, NextToken: first.NextToken,
	})
	requireAWSCode(t, "token addressed to another stream", err, "InvalidArgumentException")

	if _, err := svc.listShardsCore(reqCtx, ListShardsInput{
		StreamARN: streamA.StreamARN, NextToken: first.NextToken,
	}); err != nil {
		t.Fatalf("token with its own stream's ARN: %v", err)
	}
}

package kinesis

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/kinesis"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// newAdminHandlerEnv opens a throwaway region storage and returns the
// admin handler over it plus the store, so the console plane's thin
// adapters are pinned against the same Core paths the HTTP plane uses.
// Empty request headers resolve the default region, the same form the
// console client's requests carry.
func newAdminHandlerEnv(t *testing.T) (*AdminHandler, *kinesisstore.KinesisStore) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	bs, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get region storage: %v", err)
	}
	store := kinesisstore.NewKinesisStore(bs, "000000000000", "us-east-1")
	svc := NewKinesisService("000000000000")
	svc.SetStorageManager(mgr)
	return NewAdminHandler(svc), store
}

// TestAdminDescribeStreamPaginatesShards pins the console plane's shard
// pagination: the Limit and ExclusiveStartShardId members page the
// documented shard window, and the response's HasMoreShards is the core's
// truth — never a constant false that would tell the console a long
// stream ends at its first page.
func TestAdminDescribeStreamPaginatesShards(t *testing.T) {
	h, store := newAdminHandlerEnv(t)
	if _, err := store.CreateStream("adminpage", 2, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	ctx := context.Background()

	full, err := h.DescribeStream(ctx, connect.NewRequest(&pb.DescribeStreamInput{Streamname: proto.String("adminpage")}))
	if err != nil {
		t.Fatalf("full describe: %v", err)
	}
	desc := full.Msg.GetStreamdescription()
	if len(desc.GetShards()) != 2 || desc.GetHasmoreshards() {
		t.Fatalf("full describe: %d shards hasMore=%v, want the whole set without more", len(desc.GetShards()), desc.GetHasmoreshards())
	}

	first, err := h.DescribeStream(ctx, connect.NewRequest(&pb.DescribeStreamInput{
		Streamname: proto.String("adminpage"), Limit: proto.Int32(1),
	}))
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	desc = first.Msg.GetStreamdescription()
	if len(desc.GetShards()) != 1 || !desc.GetHasmoreshards() {
		t.Fatalf("first page: %d shards hasMore=%v, want one shard with more following", len(desc.GetShards()), desc.GetHasmoreshards())
	}
	lastShard := desc.GetShards()[0].GetShardid()

	second, err := h.DescribeStream(ctx, connect.NewRequest(&pb.DescribeStreamInput{
		Streamname: proto.String("adminpage"), Limit: proto.Int32(1), Exclusivestartshardid: proto.String(lastShard),
	}))
	if err != nil {
		t.Fatalf("second page: %v", err)
	}
	desc = second.Msg.GetStreamdescription()
	if len(desc.GetShards()) != 1 || desc.GetHasmoreshards() {
		t.Fatalf("second page: %d shards hasMore=%v, want the last shard without more", len(desc.GetShards()), desc.GetHasmoreshards())
	}
	if got := desc.GetShards()[0].GetShardid(); got == lastShard {
		t.Fatalf("second page resumed at %s, the first page's last shard", got)
	}
}

// TestAdminCreateStreamThreadsModeAndTags pins the console plane's create
// surface: the mode selection a console client sends is the mode the
// stream gets (an on-demand creation must not silently become a
// provisioned one), and create-time tags ride the locked creation.
func TestAdminCreateStreamThreadsModeAndTags(t *testing.T) {
	h, store := newAdminHandlerEnv(t)
	ctx := context.Background()

	if _, err := h.CreateStream(ctx, connect.NewRequest(&pb.CreateStreamInput{
		Streamname:        "adminod",
		Shardcount:        proto.Int32(1),
		Streammodedetails: &pb.StreamModeDetails{Streammode: pb.StreamMode_STREAM_MODE_ON_DEMAND},
		Tags:              map[string]string{"console": "yes"},
	})); err != nil {
		t.Fatalf("console create: %v", err)
	}

	stream, err := store.GetStream("adminod")
	if err != nil {
		t.Fatalf("get the created stream: %v", err)
	}
	if stream.StreamModeDetails == nil || stream.StreamModeDetails.StreamMode != kinesisstore.StreamModeOnDemand {
		t.Fatalf("created mode: %+v, want on-demand", stream.StreamModeDetails)
	}
	if stream.ShardCount != kinesisstore.DefaultOnDemandShardCount {
		t.Fatalf("on-demand shard count: %d, want the platform's four-shard default", stream.ShardCount)
	}
	tagList, err := store.ListAsSlice("adminod")
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}
	if len(tagList) != 1 || tagList[0].Key != "console" || tagList[0].Value != "yes" {
		t.Fatalf("create-time tags: %+v, want the console pair", tagList)
	}

	// A create without mode details is the provisioned default.
	if _, err := h.CreateStream(ctx, connect.NewRequest(&pb.CreateStreamInput{
		Streamname: "adminprov", Shardcount: proto.Int32(2),
	})); err != nil {
		t.Fatalf("provisioned console create: %v", err)
	}
	prov, err := store.GetStream("adminprov")
	if err != nil {
		t.Fatalf("get the provisioned stream: %v", err)
	}
	if prov.StreamModeDetails == nil || prov.StreamModeDetails.StreamMode != kinesisstore.StreamModeProvisioned {
		t.Fatalf("default mode: %+v, want provisioned", prov.StreamModeDetails)
	}
}

// TestAdminTagTrioRoundTrip pins the console plane's tag surface: the
// three RPCs the console page's tag section calls (list, tag, untag) ride
// the same ARN-addressed tag cores the HTTP framework serves, and a
// resource that does not exist answers not-found rather than tagging
// junk keys.
func TestAdminTagTrioRoundTrip(t *testing.T) {
	h, store := newAdminHandlerEnv(t)
	stream, err := store.CreateStream("admintag", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	ctx := context.Background()

	if _, err := h.TagResource(ctx, connect.NewRequest(&pb.TagResourceInput{
		Resourcearn: stream.StreamARN,
		Tags:        map[string]string{"owner": "console", "env": "test"},
	})); err != nil {
		t.Fatalf("console tag: %v", err)
	}

	listed, err := h.ListTagsForResource(ctx, connect.NewRequest(&pb.ListTagsForResourceInput{
		Resourcearn: stream.StreamARN,
	}))
	if err != nil {
		t.Fatalf("console list tags: %v", err)
	}
	got := map[string]string{}
	for _, tag := range listed.Msg.GetTags() {
		got[tag.GetKey()] = tag.GetValue()
	}
	if len(got) != 2 || got["owner"] != "console" || got["env"] != "test" {
		t.Fatalf("listed tags: %v, want both console pairs", got)
	}

	if _, err := h.UntagResource(ctx, connect.NewRequest(&pb.UntagResourceInput{
		Resourcearn: stream.StreamARN,
		Tagkeys:     []string{"env"},
	})); err != nil {
		t.Fatalf("console untag: %v", err)
	}

	listed, err = h.ListTagsForResource(ctx, connect.NewRequest(&pb.ListTagsForResourceInput{
		Resourcearn: stream.StreamARN,
	}))
	if err != nil {
		t.Fatalf("console list tags after untag: %v", err)
	}
	if tags := listed.Msg.GetTags(); len(tags) != 1 || tags[0].GetKey() != "owner" {
		t.Fatalf("tags after untag: %v, want the owner pair alone", tags)
	}

	if _, err := h.ListTagsForResource(ctx, connect.NewRequest(&pb.ListTagsForResourceInput{
		Resourcearn: "arn:aws:kinesis:us-east-1:000000000000:stream/no-such-stream",
	})); err == nil {
		t.Fatal("list tags on a stream that does not exist succeeded")
	}
}

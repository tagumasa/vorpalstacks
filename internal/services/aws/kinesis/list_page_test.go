package kinesis

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// newListPageTestEnv opens a throwaway region storage and returns the
// service, a store over that storage and a request context wired to it.
func newListPageTestEnv(t *testing.T) (*KinesisService, *kinesisstore.KinesisStore, *request.RequestContext) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("new region storage manager: %v", err)
	}
	bs, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("get region storage: %v", err)
	}
	tstore, ok := bs.(storage.TransactionalStorageWith2PC)
	if !ok {
		t.Fatal("region storage does not support TransactionalStorageWith2PC")
	}
	store := kinesisstore.NewKinesisStore(tstore, "000000000000", "us-east-1")
	svc := NewKinesisService("000000000000", "us-east-1")
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	return svc, store, reqCtx
}

// TestListStreamConsumersPageSize pins the documented MaxResults semantics:
// the accepted input window is 1-10000 while the effective page carries a
// default of 100 and never exceeds 100, with the remainder resumable
// through NextToken.
func TestListStreamConsumersPageSize(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	stream, err := store.CreateStream("page_consumers", 1, kinesisstore.StreamModeProvisioned, 0, 0)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	for i := 0; i < 101; i++ {
		if _, err := store.RegisterStreamConsumer(stream.StreamARN, fmt.Sprintf("reader-%03d", i)); err != nil {
			t.Fatalf("register consumer %d: %v", i, err)
		}
	}

	first, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{StreamARN: stream.StreamARN})
	if err != nil {
		t.Fatalf("list with omitted MaxResults: %v", err)
	}
	if len(first.Consumers) != 100 {
		t.Fatalf("default page: got %d consumers, want 100", len(first.Consumers))
	}
	if !first.HasMore || first.NextToken == nil {
		t.Fatalf("default page must report the remainder, hasMore=%v nextToken=%v", first.HasMore, first.NextToken)
	}

	above, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:     stream.StreamARN,
		MaxResults:    150,
		HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("list with MaxResults above the page cap: %v", err)
	}
	if len(above.Consumers) != 100 {
		t.Fatalf("MaxResults 150: got %d consumers, want the capped 100", len(above.Consumers))
	}

	half, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:     stream.StreamARN,
		MaxResults:    50,
		HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("list with MaxResults 50: %v", err)
	}
	if len(half.Consumers) != 50 || !half.HasMore || half.NextToken == nil {
		t.Fatalf("MaxResults 50: got %d consumers hasMore=%v, want 50 with a token", len(half.Consumers), half.HasMore)
	}

	tail, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN: stream.StreamARN,
		NextToken: *first.NextToken,
	})
	if err != nil {
		t.Fatalf("list following the default-page token: %v", err)
	}
	if len(tail.Consumers) != 1 || tail.HasMore {
		t.Fatalf("token follow-up: got %d consumers hasMore=%v, want 1 without more", len(tail.Consumers), tail.HasMore)
	}

	for _, maxResults := range []int{0, 10001} {
		if _, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
			StreamARN:     stream.StreamARN,
			MaxResults:    maxResults,
			HasMaxResults: true,
		}); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("MaxResults %d: expected InvalidArgumentException, got: %v", maxResults, err)
		}
	}
}

// TestListStreamsLimitClamp pins the documented Limit semantics: the input
// window is 1-10000, the effective page defaults to 100 and never exceeds
// 100, and an explicitly provided out-of-window value is rejected instead
// of being folded into the default.
func TestListStreamsLimitClamp(t *testing.T) {
	svc, store, _ := newListPageTestEnv(t)

	for i := 0; i < 101; i++ {
		if _, err := store.CreateStream(fmt.Sprintf("stream-%03d", i), 1, kinesisstore.StreamModeProvisioned, 0, 0); err != nil {
			t.Fatalf("create stream %d: %v", i, err)
		}
	}

	first, err := svc.listStreamsCore(store, ListStreamsInput{})
	if err != nil {
		t.Fatalf("list with omitted Limit: %v", err)
	}
	if len(first.Streams) != 100 || !first.IsTruncated || first.NextMarker == "" {
		t.Fatalf("default page: got %d streams truncated=%v, want 100 with a marker", len(first.Streams), first.IsTruncated)
	}

	above, err := svc.listStreamsCore(store, ListStreamsInput{Limit: 150, HasLimit: true})
	if err != nil {
		t.Fatalf("list with Limit above the page cap: %v", err)
	}
	if len(above.Streams) != 100 {
		t.Fatalf("Limit 150: got %d streams, want the capped 100", len(above.Streams))
	}

	if _, err := svc.listStreamsCore(store, ListStreamsInput{Limit: 0, HasLimit: true}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("explicit Limit 0: expected InvalidArgumentException, got: %v", err)
	}

	if _, err := svc.listStreamsCore(store, ListStreamsInput{Limit: 10001, HasLimit: true}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("Limit 10001: expected InvalidArgumentException, got: %v", err)
	}
}

// TestListShardsMaxResults pins the documented MaxResults semantics: the
// input window is 1-10000, the effective page defaults to 1000 and never
// exceeds 1000, and an explicitly provided out-of-window value is rejected
// instead of being folded into the default.
func TestListShardsMaxResults(t *testing.T) {
	svc, store, _ := newListPageTestEnv(t)

	if _, err := store.CreateStream("page_shards", 2, kinesisstore.StreamModeProvisioned, 0, 0); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	one, err := svc.listShardsCore(store, ListShardsInput{StreamName: "page_shards", MaxResults: 1, HasMaxResults: true})
	if err != nil {
		t.Fatalf("list with MaxResults 1: %v", err)
	}
	if len(one.Shards) != 1 {
		t.Fatalf("MaxResults 1: got %d shards, want 1", len(one.Shards))
	}

	all, err := svc.listShardsCore(store, ListShardsInput{StreamName: "page_shards"})
	if err != nil {
		t.Fatalf("list with omitted MaxResults: %v", err)
	}
	if len(all.Shards) != 2 {
		t.Fatalf("omitted MaxResults: got %d shards, want 2", len(all.Shards))
	}

	for _, maxResults := range []int{0, 10001} {
		if _, err := svc.listShardsCore(store, ListShardsInput{StreamName: "page_shards", MaxResults: maxResults, HasMaxResults: true}); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("MaxResults %d: expected InvalidArgumentException, got: %v", maxResults, err)
		}
	}
}

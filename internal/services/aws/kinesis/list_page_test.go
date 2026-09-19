package kinesis

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

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
	store := kinesisstore.NewKinesisStore(bs, "000000000000", "us-east-1")
	svc := NewKinesisService("000000000000")
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	return svc, store, reqCtx
}

// TestListStreamConsumersPageSize pins the documented MaxResults semantics
// against the per-stream consumer quota: the accepted input window is
// 1-10000, an in-window value pages with a resumable token, and an
// out-of-window value is rejected. (The store's locked registration now
// enforces the twenty-consumer quota, so the fixture registers twenty —
// the hundred-entry default page itself is pinned by the ListStreams
// counterpart, which has no per-stream quota.)
func TestListStreamConsumersPageSize(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	stream, err := store.CreateStream("page_consumers", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	for i := 0; i < kinesisstore.MaxConsumersPerStream; i++ {
		if _, err := store.RegisterStreamConsumer(stream.StreamARN, fmt.Sprintf("reader-%03d", i), nil); err != nil {
			t.Fatalf("register consumer %d: %v", i, err)
		}
	}

	first, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{StreamARN: stream.StreamARN})
	if err != nil {
		t.Fatalf("list with omitted MaxResults: %v", err)
	}
	if len(first.Consumers) != kinesisstore.MaxConsumersPerStream || first.NextToken != nil {
		t.Fatalf("default page: got %d consumers, want the full quota without a token", len(first.Consumers))
	}

	page, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:     stream.StreamARN,
		MaxResults:    15,
		HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("list with MaxResults 15: %v", err)
	}
	if len(page.Consumers) != 15 || page.NextToken == nil {
		t.Fatalf("MaxResults 15: got %d consumers, want 15 with a token", len(page.Consumers))
	}

	tail, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN: stream.StreamARN,
		NextToken: *page.NextToken,
	})
	if err != nil {
		t.Fatalf("list following the page token: %v", err)
	}
	if len(tail.Consumers) != 5 || tail.NextToken != nil {
		t.Fatalf("token follow-up: got %d consumers, want the remaining 5 without a token", len(tail.Consumers))
	}

	above, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:     stream.StreamARN,
		MaxResults:    25,
		HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("list with MaxResults above the consumer count: %v", err)
	}
	if len(above.Consumers) != kinesisstore.MaxConsumersPerStream || above.NextToken != nil {
		t.Fatalf("MaxResults 25: got %d consumers, want all without a token", len(above.Consumers))
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
	svc, store, reqCtx := newListPageTestEnv(t)

	for i := 0; i < 101; i++ {
		if _, err := store.CreateStream(fmt.Sprintf("stream-%03d", i), 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
			t.Fatalf("create stream %d: %v", i, err)
		}
	}

	first, err := svc.listStreamsCore(reqCtx, ListStreamsInput{})
	if err != nil {
		t.Fatalf("list with omitted Limit: %v", err)
	}
	if len(first.Streams) != 100 || !first.IsTruncated || first.NextMarker == "" {
		t.Fatalf("default page: got %d streams truncated=%v, want 100 with a marker", len(first.Streams), first.IsTruncated)
	}

	above, err := svc.listStreamsCore(reqCtx, ListStreamsInput{Limit: 150, HasLimit: true})
	if err != nil {
		t.Fatalf("list with Limit above the page cap: %v", err)
	}
	if len(above.Streams) != 100 {
		t.Fatalf("Limit 150: got %d streams, want the capped 100", len(above.Streams))
	}

	if _, err := svc.listStreamsCore(reqCtx, ListStreamsInput{Limit: 0, HasLimit: true}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("explicit Limit 0: expected InvalidArgumentException, got: %v", err)
	}

	if _, err := svc.listStreamsCore(reqCtx, ListStreamsInput{Limit: 10001, HasLimit: true}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("Limit 10001: expected InvalidArgumentException, got: %v", err)
	}
}

// TestRegisterStreamConsumerQuota pins the documented per-stream
// registered-consumer quota: twenty for Provisioned and On-demand Standard
// streams, with a deregistration freeing its slot again.
func TestRegisterStreamConsumerQuota(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	stream, err := store.CreateStream("quota_consumers", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	for i := 0; i < 20; i++ {
		if _, err := store.RegisterStreamConsumer(stream.StreamARN, fmt.Sprintf("reader-%02d", i), nil); err != nil {
			t.Fatalf("register consumer %d: %v", i, err)
		}
	}

	if _, err := svc.registerStreamConsumerCore(reqCtx, RegisterStreamConsumerInput{
		StreamARN:    stream.StreamARN,
		ConsumerName: "reader-20",
	}); !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("21st consumer: expected LimitExceededException, got: %v", err)
	}

	// Deregistering one consumer frees its slot for the next registration.
	if err := svc.deregisterStreamConsumerCore(reqCtx, DeregisterStreamConsumerInput{
		StreamARN:    stream.StreamARN,
		ConsumerName: "reader-00",
	}); err != nil {
		t.Fatalf("deregister a consumer: %v", err)
	}
	if _, err := svc.registerStreamConsumerCore(reqCtx, RegisterStreamConsumerInput{
		StreamARN:    stream.StreamARN,
		ConsumerName: "reader-20",
	}); err != nil {
		t.Fatalf("register after freeing a slot: %v", err)
	}
}

// TestListOpsExpiredNextToken pins the declared pagination-token error:
// every list operation that declares ExpiredNextTokenException answers it
// for a token that cannot be decoded — never a silent restart at page one —
// and an emitted token round-trips to the next page.
func TestListOpsExpiredNextToken(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	consumerStream, err := store.CreateStream("token_consumers", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create consumer stream: %v", err)
	}
	for i := 0; i < 3; i++ {
		if _, err := store.CreateStream(fmt.Sprintf("token-stream-%02d", i), 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
			t.Fatalf("create stream %d: %v", i, err)
		}
	}

	first, err := svc.listStreamsCore(reqCtx, ListStreamsInput{Limit: 2, HasLimit: true})
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	if len(first.Streams) != 2 || !first.IsTruncated || first.NextMarker == "" {
		t.Fatalf("first page: got %d streams truncated=%v, want 2 with a marker", len(first.Streams), first.IsTruncated)
	}

	second, err := svc.listStreamsCore(reqCtx, ListStreamsInput{Limit: 2, HasLimit: true, NextToken: first.NextMarker})
	if err != nil {
		t.Fatalf("token follow-up: %v", err)
	}
	if len(second.Streams) != 2 || second.IsTruncated {
		t.Fatalf("token follow-up: got %d streams truncated=%v, want the remaining 2 without more", len(second.Streams), second.IsTruncated)
	}

	for _, tc := range []struct {
		name   string
		invoke func() error
	}{
		{"ListStreams", func() error {
			_, err := svc.listStreamsCore(reqCtx, ListStreamsInput{NextToken: "!!!not-base64!!!"})
			return err
		}},
		{"ListStreamConsumers", func() error {
			_, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
				StreamARN: consumerStream.StreamARN,
				NextToken: "!!!not-base64!!!",
			})
			return err
		}},
		{"ListStreamConsumers token past its validity window", func() error {
			old := fmt.Sprintf("%s#%d", consumerStream.StreamARN, time.Now().Add(-kinesisstore.NextTokenValidity-time.Minute).UnixNano())
			_, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
				StreamARN: consumerStream.StreamARN,
				NextToken: base64.StdEncoding.EncodeToString([]byte(old)),
			})
			return err
		}},
		{"ListShards handler", func() error {
			_, err := svc.ListShards(context.Background(), failingStorageReqCtx(t), parsedRequest(map[string]interface{}{
				"NextToken": "!!!not-base64!!!",
			}))
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.invoke(); !errors.Is(err, ErrExpiredNextToken) {
				t.Fatalf("want ExpiredNextTokenException, got: %v", err)
			}
		})
	}

	// The token "unambiguously identifies the stream": the generation
	// member cannot travel with it, agreeing or not.
	requireAWSCode(t, "ListStreamConsumers token with StreamCreationTimestamp", func() error {
		_, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
			StreamARN:               consumerStream.StreamARN,
			NextToken:               "!!!not-base64!!!",
			StreamCreationTimestamp: "1",
		})
		return err
	}(), "InvalidArgumentException")
}

// TestConsumerVisibleUnderUnnormalisedStreamARN pins the write-boundary ARN
// normalisation: a consumer registered against a stream ARN spelling with an
// omitted account segment still carries the stream's stored ARN, so it
// remains visible in its own stream's consumer list.
func TestConsumerVisibleUnderUnnormalisedStreamARN(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	stream, err := store.CreateStream("arn_norm_consumers", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	unnormalised := "arn:aws:kinesis:us-east-1::stream/arn_norm_consumers"

	if _, err := svc.registerStreamConsumerCore(reqCtx, RegisterStreamConsumerInput{
		StreamARN:    unnormalised,
		ConsumerName: "reader-arn",
	}); err != nil {
		t.Fatalf("register through the un-normalised ARN: %v", err)
	}

	listed, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN: stream.StreamARN,
	})
	if err != nil {
		t.Fatalf("list consumers: %v", err)
	}
	if len(listed.Consumers) != 1 || listed.Consumers[0].ConsumerName != "reader-arn" {
		t.Fatalf("consumer registered via an un-normalised ARN must stay listed, got %d consumers", len(listed.Consumers))
	}
	if listed.Consumers[0].StreamARN != stream.StreamARN {
		t.Fatalf("consumer StreamARN: got %q, want the stream's stored ARN %q", listed.Consumers[0].StreamARN, stream.StreamARN)
	}
}

// TestListShardsMaxResults pins the documented MaxResults semantics: the
// input window is 1-10000, the effective page defaults to 1000 and never
// exceeds 1000, and an explicitly provided out-of-window value is rejected
// instead of being folded into the default.
func TestListShardsMaxResults(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	if _, err := store.CreateStream("page_shards", 2, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	one, err := svc.listShardsCore(reqCtx, ListShardsInput{StreamName: "page_shards", MaxResults: 1, HasMaxResults: true})
	if err != nil {
		t.Fatalf("list with MaxResults 1: %v", err)
	}
	if len(one.Shards) != 1 {
		t.Fatalf("MaxResults 1: got %d shards, want 1", len(one.Shards))
	}

	all, err := svc.listShardsCore(reqCtx, ListShardsInput{StreamName: "page_shards"})
	if err != nil {
		t.Fatalf("list with omitted MaxResults: %v", err)
	}
	if len(all.Shards) != 2 {
		t.Fatalf("omitted MaxResults: got %d shards, want 2", len(all.Shards))
	}

	for _, maxResults := range []int{0, 10001} {
		if _, err := svc.listShardsCore(reqCtx, ListShardsInput{StreamName: "page_shards", MaxResults: maxResults, HasMaxResults: true}); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("MaxResults %d: expected InvalidArgumentException, got: %v", maxResults, err)
		}
	}
}

// TestListStreamConsumersTokenAnchorExpired pins the resumption rule on
// the anchor side: a decodable NextToken whose anchor consumer is no
// longer registered answers ExpiredNextTokenException — the list changed
// under the client — never a silent restart from page one.
func TestListStreamConsumersTokenAnchorExpired(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("page_expired", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	for i := 0; i < kinesisstore.MaxConsumersPerStream; i++ {
		if _, err := store.RegisterStreamConsumer(stream.StreamARN, fmt.Sprintf("reader-%03d", i), nil); err != nil {
			t.Fatalf("register consumer %d: %v", i, err)
		}
	}

	page, err := svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN:     stream.StreamARN,
		MaxResults:    15,
		HasMaxResults: true,
	})
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	if page.NextToken == nil {
		t.Fatal("first page carried no token")
	}

	// The token is the encoded anchor ARN with its issuance time; strip
	// the stamp and deregister the anchor itself.
	anchorBytes, err := base64.StdEncoding.DecodeString(*page.NextToken)
	if err != nil {
		t.Fatalf("decode the page token: %v", err)
	}
	anchorARN := string(anchorBytes)
	if sep := strings.LastIndexByte(anchorARN, '#'); sep >= 0 {
		anchorARN = anchorARN[:sep]
	}
	if err := store.DeregisterStreamConsumer(anchorARN); err != nil {
		t.Fatalf("deregister the anchor: %v", err)
	}

	_, err = svc.listStreamConsumersCore(reqCtx, ListStreamConsumersInput{
		StreamARN: stream.StreamARN,
		NextToken: *page.NextToken,
	})
	requireAWSCode(t, "token whose anchor was deregistered", err, "ExpiredNextTokenException")
}

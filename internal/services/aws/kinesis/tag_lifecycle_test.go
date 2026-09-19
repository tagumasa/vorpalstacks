package kinesis

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"fmt"

	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// tagMapOf converts a ListTags* response's Tags slice to a key-value map.
func tagMapOf(t *testing.T, resp interface{}) map[string]string {
	t.Helper()
	m, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("response is %T, want map", resp)
	}
	raw, ok := m["Tags"].([]map[string]interface{})
	if !ok {
		t.Fatalf("Tags is %T, want []map[string]interface{}", m["Tags"])
	}
	out := make(map[string]string, len(raw))
	for _, entry := range raw {
		key, _ := entry["Key"].(string)
		value, _ := entry["Value"].(string)
		out[key] = value
	}
	return out
}

// registerTaggedConsumer registers a consumer through the public handler
// with the given create-time tags and returns its ARN.
func registerTaggedConsumer(t *testing.T, svc *KinesisService, reqCtx *request.RequestContext, streamARN string, tagsParam map[string]interface{}) string {
	t.Helper()
	resp, err := svc.RegisterStreamConsumer(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamARN":    streamARN,
		"ConsumerName": "tagged-reader",
		"Tags":         tagsParam,
	}))
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	consumer, ok := resp.(map[string]interface{})["Consumer"].(map[string]interface{})
	if !ok {
		t.Fatalf("register response carries no Consumer object: %T", resp)
	}
	arn, ok := consumer["ConsumerARN"].(string)
	if !ok || arn == "" {
		t.Fatalf("register response carries no ConsumerARN: %v", consumer["ConsumerARN"])
	}
	return arn
}

// TestConsumerTagLifecycleThroughTagAPI pins the consumer half of the tag
// resource space: create-time tags are reachable through ListTagsForResource,
// TagResource and UntagResource address the consumer by ARN, and the
// consumer's entries never leak into the owning stream's tag set.
func TestConsumerTagLifecycleThroughTagAPI(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("ctag", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	if err := store.Tag("ctag", map[string]string{"scope": "stream"}); err != nil {
		t.Fatalf("tag stream: %v", err)
	}

	consumerARN := registerTaggedConsumer(t, svc, reqCtx, stream.StreamARN, map[string]interface{}{"origin": "register"})

	resp, err := svc.ListTagsForResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": consumerARN,
	}))
	if err != nil {
		t.Fatalf("list consumer tags: %v", err)
	}
	tags := tagMapOf(t, resp)
	if tags["origin"] != "register" {
		t.Fatalf("register-time tag unreachable through ListTagsForResource: %v", tags)
	}

	if _, err := svc.TagResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": consumerARN,
		"Tags":        map[string]interface{}{"stage": "api"},
	})); err != nil {
		t.Fatalf("tag consumer: %v", err)
	}

	resp, err = svc.ListTagsForResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": consumerARN,
	}))
	if err != nil {
		t.Fatalf("list consumer tags after tag: %v", err)
	}
	tags = tagMapOf(t, resp)
	if len(tags) != 2 || tags["origin"] != "register" || tags["stage"] != "api" {
		t.Fatalf("consumer tags after TagResource: %v", tags)
	}

	// The stream's own tag set is unaffected by the consumer's entries.
	resp, err = svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctag",
	}))
	if err != nil {
		t.Fatalf("list stream tags: %v", err)
	}
	if streamTags := tagMapOf(t, resp); len(streamTags) != 1 || streamTags["scope"] != "stream" {
		t.Fatalf("consumer tags leaked into the stream tag set: %v", streamTags)
	}

	if _, err := svc.UntagResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": consumerARN,
		"TagKeys":     []interface{}{"origin"},
	})); err != nil {
		t.Fatalf("untag consumer: %v", err)
	}

	resp, err = svc.ListTagsForResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": consumerARN,
	}))
	if err != nil {
		t.Fatalf("list consumer tags after untag: %v", err)
	}
	if tags = tagMapOf(t, resp); len(tags) != 1 || tags["stage"] != "api" {
		t.Fatalf("consumer tags after UntagResource: %v", tags)
	}

	// An ARN that is neither a stream nor a registered consumer is
	// not-found, never a silent write keyed by the raw ARN.
	_, err = svc.TagResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": stream.StreamARN + "/consumer/ghost:1",
		"Tags":        map[string]interface{}{"k": "v"},
	}))
	requireAWSCode(t, "TagResource unknown consumer", err, "ResourceNotFoundException")
}

// TestTagValidationBounds pins the request-level tag bounds the model
// declares: the tag set passes the shared validation on every write path
// (add, create-time stream, create-time consumer), TagKeys obeys the
// TagKeyList length ceiling, and ListTagsForStream's Limit window rejects
// out-of-range values instead of clamping them.
func TestTagValidationBounds(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("ctagb", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}

	tooMany := make(map[string]interface{}, 51)
	for i := 0; i < 51; i++ {
		tooMany["k"+strconv.Itoa(i)] = "v"
	}
	longKey := map[string]interface{}{strings.Repeat("k", 129): "v"}
	longValue := map[string]interface{}{"k": strings.Repeat("v", 257)}

	_, err = svc.AddTagsToStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagb",
		"Tags":       tooMany,
	}))
	requireAWSCode(t, "AddTagsToStream 51 tags", err, "InvalidArgumentException")
	_, err = svc.AddTagsToStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagb",
		"Tags":       longKey,
	}))
	requireAWSCode(t, "AddTagsToStream 129-char key", err, "InvalidArgumentException")
	_, err = svc.TagResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"ResourceARN": stream.StreamARN,
		"Tags":        longValue,
	}))
	requireAWSCode(t, "TagResource 257-char value", err, "InvalidArgumentException")

	tagKeys51 := make([]interface{}, 51)
	for i := range tagKeys51 {
		tagKeys51[i] = "k" + strconv.Itoa(i)
	}
	_, err = svc.RemoveTagsFromStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagb",
		"TagKeys":    tagKeys51,
	}))
	requireAWSCode(t, "RemoveTagsFromStream 51 tag keys", err, "InvalidArgumentException")

	_, err = svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagb",
		"Limit":      51,
	}))
	requireAWSCode(t, "ListTagsForStream Limit 51", err, "InvalidArgumentException")
	_, err = svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagb",
		"Limit":      0,
	}))
	requireAWSCode(t, "ListTagsForStream Limit 0", err, "InvalidArgumentException")

	_, err = svc.RegisterStreamConsumer(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamARN":    stream.StreamARN,
		"ConsumerName": "bounds-reader",
		"Tags":         tooMany,
	}))
	requireAWSCode(t, "RegisterStreamConsumer 51 tags", err, "InvalidArgumentException")

	_, err = svc.CreateStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagb-invalid-tags",
		"ShardCount": 1,
		"Tags":       longKey,
	}))
	requireAWSCode(t, "CreateStream 129-char key", err, "InvalidArgumentException")
	if _, err := store.GetStream("ctagb-invalid-tags"); err == nil {
		t.Fatalf("CreateStream with invalid tags created the stream")
	}
}

// TestTagListPaginationWindow pins the in-window behaviour the bound change
// must preserve: Limit pages lexically from ExclusiveStartTagKey and reports
// HasMoreTags, while an absent Limit serves one full page.
func TestTagListPaginationWindow(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("ctagp", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	page := make(map[string]interface{}, 3)
	for _, k := range []string{"alpha", "beta", "gamma"} {
		page[k] = "v"
	}
	if _, err := svc.AddTagsToStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagp",
		"Tags":       page,
	})); err != nil {
		t.Fatalf("add tags: %v", err)
	}

	resp, err := svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagp",
		"Limit":      2,
	}))
	if err != nil {
		t.Fatalf("list page 1: %v", err)
	}
	first := resp.(map[string]interface{})
	if got := tagMapOf(t, first); len(got) != 2 || got["alpha"] != "v" || got["beta"] != "v" {
		t.Fatalf("page 1 keys: %v", got)
	}
	if first["HasMoreTags"] != true {
		t.Fatalf("page 1 HasMoreTags: %v", first["HasMoreTags"])
	}

	resp, err = svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName":           "ctagp",
		"Limit":                2,
		"ExclusiveStartTagKey": "beta",
	}))
	if err != nil {
		t.Fatalf("list page 2: %v", err)
	}
	second := resp.(map[string]interface{})
	if got := tagMapOf(t, second); len(got) != 1 || got["gamma"] != "v" {
		t.Fatalf("page 2 keys: %v", got)
	}
	if second["HasMoreTags"] != false {
		t.Fatalf("page 2 HasMoreTags: %v", second["HasMoreTags"])
	}

	// An absent Limit serves the whole set — the model's Tags output is
	// required even when the resource carries no tags.
	resp, err = svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "ctagp",
	}))
	if err != nil {
		t.Fatalf("list full page: %v", err)
	}
	if got := tagMapOf(t, resp); len(got) != 3 {
		t.Fatalf("full page keys: %v", got)
	}
}

// TestTagTotalCapFifty pins the per-resource tag total the tag members'
// documentation states: the merged key count across writes stays within
// fifty — a later write whose new keys would cross the bound answers
// InvalidArgumentException, while overwriting existing keys at the cap
// stays legal because overwritten keys add nothing to the merged set.
func TestTagTotalCapFifty(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("tag_cap", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}

	build := func(prefix string, n int) []types.Tag {
		set := make([]types.Tag, n)
		for i := range set {
			set[i] = types.Tag{Key: fmt.Sprintf("%s-%02d", prefix, i), Value: "v"}
		}
		return set
	}

	if err := svc.tagResourceCore(reqCtx, stream.StreamARN, build("k", kinesisstore.MaxTagsPerResource)); err != nil {
		t.Fatalf("first write at the cap: %v", err)
	}
	err = svc.tagResourceCore(reqCtx, stream.StreamARN, []types.Tag{{Key: "over", Value: "v"}})
	requireAWSCode(t, "write past the cap", err, "InvalidArgumentException")
	if err := svc.tagResourceCore(reqCtx, stream.StreamARN, []types.Tag{{Key: "k-00", Value: "rewritten"}}); err != nil {
		t.Fatalf("overwrite at the cap: %v", err)
	}

	tags, err := svc.listResourceTagsCore(reqCtx, stream.StreamARN)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(tags) != kinesisstore.MaxTagsPerResource {
		t.Fatalf("tag count: got %d, want the cap", len(tags))
	}
	for _, tag := range tags {
		if tag.Key == "k-00" && tag.Value != "rewritten" {
			t.Fatalf("overwrite did not land: %+v", tag)
		}
	}
}

// TestUntagAndListCoresRoundTrip pins the ARN-addressed untag and list
// cores through the store calls whose failures now map with the same wrap
// the tag core applies: tags land, a removed key leaves the listing, and
// removing the last key empties it.
func TestUntagAndListCoresRoundTrip(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("untag_rt", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}

	if err := svc.tagResourceCore(reqCtx, stream.StreamARN, []types.Tag{
		{Key: "env", Value: "prod"}, {Key: "team", Value: "core"},
	}); err != nil {
		t.Fatalf("tag: %v", err)
	}
	if err := svc.untagResourceCore(reqCtx, stream.StreamARN, []string{"env"}); err != nil {
		t.Fatalf("untag: %v", err)
	}
	tags, err := svc.listResourceTagsCore(reqCtx, stream.StreamARN)
	if err != nil {
		t.Fatalf("list after removal: %v", err)
	}
	if len(tags) != 1 || tags[0].Key != "team" || tags[0].Value != "core" {
		t.Fatalf("listing after removal: got %+v, want the surviving team tag", tags)
	}
	if err := svc.untagResourceCore(reqCtx, stream.StreamARN, []string{"team"}); err != nil {
		t.Fatalf("untag the last key: %v", err)
	}
	tags, err = svc.listResourceTagsCore(reqCtx, stream.StreamARN)
	if err != nil {
		t.Fatalf("list after emptying: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("listing after emptying: got %+v, want no tags", tags)
	}
}

// TestCreateWithTagsIsAtomic pins the creation-transaction rule on both
// creation paths: when the initial tag set cannot commit — here by
// crossing the per-resource cap — the resource itself does not exist
// either, and the stream's consumer count never moved.
func TestCreateWithTagsIsAtomic(t *testing.T) {
	_, store, _ := newListPageTestEnv(t)

	over := make(map[string]string, kinesisstore.MaxTagsPerResource+1)
	for i := 0; i <= kinesisstore.MaxTagsPerResource; i++ {
		over[fmt.Sprintf("k-%02d", i)] = "v"
	}

	if _, err := store.CreateStream("atomic_stream", 1, kinesisstore.StreamModeProvisioned, 0, 0, over); err == nil {
		t.Fatal("create with an over-cap tag set: accepted, want rejection")
	}
	if _, err := store.GetStream("atomic_stream"); err == nil {
		t.Fatal("stream exists after its tag set failed to commit")
	}

	base, err := store.CreateStream("atomic_base", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create base: %v", err)
	}
	if _, err := store.RegisterStreamConsumer(base.StreamARN, "atomic", over); err == nil {
		t.Fatal("register with an over-cap tag set: accepted, want rejection")
	}
	consumers, err := store.ListStreamConsumers("atomic_base")
	if err != nil {
		t.Fatalf("list consumers: %v", err)
	}
	if len(consumers) != 0 {
		t.Fatalf("consumer exists after its tag set failed to commit: %d", len(consumers))
	}
	after, err := store.GetStream("atomic_base")
	if err != nil {
		t.Fatalf("get stream: %v", err)
	}
	if after.ConsumerCount != base.ConsumerCount {
		t.Fatalf("consumer count moved: got %d, want %d", after.ConsumerCount, base.ConsumerCount)
	}
}

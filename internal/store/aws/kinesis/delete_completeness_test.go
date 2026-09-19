package kinesis

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"
)

// iteratorIDs censuses the iterator bucket as a map of iterator ID to the
// owning stream name.
func iteratorIDs(t *testing.T, store *KinesisStore) map[string]string {
	t.Helper()
	ids := make(map[string]string)
	if err := store.iteratorsStore.ForEach(func(key string, value []byte) error {
		var it ShardIterator
		if err := json.Unmarshal(value, &it); err != nil {
			return err
		}
		ids[it.IteratorID] = it.StreamName
		return nil
	}); err != nil {
		t.Fatalf("census iterators: %v", err)
	}
	return ids
}

// TestDeleteStreamRemovesResourcePolicy pins the policy entry into the
// delete-time sweep: every ARN spelling normalises onto one entry, the entry
// dies with the stream, and a stream recreated under the same name inherits
// nothing.
func TestDeleteStreamRemovesResourcePolicy(t *testing.T) {
	store := newLockTestStore(t)
	stream, err := store.CreateStream("policy_sweep", 1, StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}

	canonical := stream.StreamARN
	elidedAccount := "arn:aws:kinesis:us-east-1::stream/policy_sweep"
	if err := store.PutResourcePolicy(canonical, `{"doc":1}`); err != nil {
		t.Fatalf("put policy under canonical spelling: %v", err)
	}
	if err := store.PutResourcePolicy(elidedAccount, `{"doc":2}`); err != nil {
		t.Fatalf("put policy under elided-account spelling: %v", err)
	}
	if got, err := store.GetResourcePolicy(canonical); err != nil || got != `{"doc":2}` {
		t.Fatalf("policy under canonical spelling: got %q, %v — want the normalised single entry", got, err)
	}

	if err := store.DeleteStream("policy_sweep"); err != nil {
		t.Fatalf("delete stream: %v", err)
	}
	if _, err := store.GetResourcePolicy(canonical); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("policy after delete, canonical spelling: %v, want ErrResourceNotFound", err)
	}
	if _, err := store.GetResourcePolicy(elidedAccount); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("policy after delete, elided spelling: %v, want ErrResourceNotFound", err)
	}

	if _, err := store.CreateStream("policy_sweep", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("recreate stream: %v", err)
	}
	if _, err := store.GetResourcePolicy(canonical); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("policy after recreation: %v, want ErrResourceNotFound — a recreated stream inherits nothing", err)
	}
}

// TestDeleteStreamSweepsShardIterators pins the iterator sweep: deleting a
// stream invalidates every iterator that stream owns and leaves other
// streams' iterators untouched.
func TestDeleteStreamSweepsShardIterators(t *testing.T) {
	store := newLockTestStore(t)
	if _, err := store.CreateStream("iter_sweep", 2, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create swept stream: %v", err)
	}
	if _, err := store.CreateStream("iter_keep", 1, StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create kept stream: %v", err)
	}

	sweptA, err := store.CreateShardIterator("iter_sweep", "shardId-000000000000", "LATEST", "", nil)
	if err != nil {
		t.Fatalf("create iterator A: %v", err)
	}
	sweptB, err := store.CreateShardIterator("iter_sweep", "shardId-000000000001", "TRIM_HORIZON", "", nil)
	if err != nil {
		t.Fatalf("create iterator B: %v", err)
	}
	kept, err := store.CreateShardIterator("iter_keep", "shardId-000000000000", "LATEST", "", nil)
	if err != nil {
		t.Fatalf("create kept iterator: %v", err)
	}

	if err := store.DeleteStream("iter_sweep"); err != nil {
		t.Fatalf("delete stream: %v", err)
	}

	remaining := iteratorIDs(t, store)
	if len(remaining) != 1 {
		t.Fatalf("iterator census after delete: %d remain, want 1 (%v)", len(remaining), remaining)
	}
	if _, ok := remaining[kept.IteratorID]; !ok {
		t.Fatalf("the untouched stream's iterator was swept: %v", remaining)
	}
	if _, ok := remaining[sweptA.IteratorID]; ok {
		t.Fatalf("swept iterator A survived the delete: %v", remaining)
	}
	if _, ok := remaining[sweptB.IteratorID]; ok {
		t.Fatalf("swept iterator B survived the delete: %v", remaining)
	}
	if _, err := store.GetShardIterator(sweptA.IteratorID); !errors.Is(err, ErrInvalidIterator) {
		t.Fatalf("swept iterator resolves: %v, want ErrInvalidIterator", err)
	}
}

// TestDeleteStreamSweepsConsumerTags pins the consumer-tag sweep at stream
// deletion: consumer tag entries key by consumer ARN and die with the stream.
func TestDeleteStreamSweepsConsumerTags(t *testing.T) {
	store := newLockTestStore(t)
	stream, err := store.CreateStream("ctag_sweep", 1, StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	if err := store.Tag(consumer.ConsumerARN, map[string]string{"env": "test"}); err != nil {
		t.Fatalf("tag consumer: %v", err)
	}
	if err := store.Tag("ctag_sweep", map[string]string{"team": "kinesis"}); err != nil {
		t.Fatalf("tag stream: %v", err)
	}

	if err := store.DeleteStream("ctag_sweep"); err != nil {
		t.Fatalf("delete stream: %v", err)
	}

	if tags, err := store.List(consumer.ConsumerARN); err != nil || len(tags) != 0 {
		t.Fatalf("consumer tags after stream delete: %v, %d entries — want none", err, len(tags))
	}
	if tags, err := store.List("ctag_sweep"); err != nil || len(tags) != 0 {
		t.Fatalf("stream tags after stream delete: %v, %d entries — want none", err, len(tags))
	}
}

// TestDeregisterStreamConsumerRemovesTags pins the consumer-tag sweep at
// deregistration: the entry dies with the consumer.
func TestDeregisterStreamConsumerRemovesTags(t *testing.T) {
	store := newLockTestStore(t)
	stream, err := store.CreateStream("ctag_dereg", 1, StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	if err := store.Tag(consumer.ConsumerARN, map[string]string{"env": "test"}); err != nil {
		t.Fatalf("tag consumer: %v", err)
	}

	if err := store.DeregisterStreamConsumer(consumer.ConsumerARN); err != nil {
		t.Fatalf("deregister consumer: %v", err)
	}

	if tags, err := store.List(consumer.ConsumerARN); err != nil || len(tags) != 0 {
		t.Fatalf("consumer tags after deregister: %v, %d entries — want none", err, len(tags))
	}
}

// TestConsumerARNCarriesItsCreationTimestamp pins the consumer ARN
// family's pattern: the ARN ends in the consumer's creation stamp, so a
// deregistered name re-registered later is a different consumer with a
// different ARN, and name addressing still finds the live registration
// (the stamp makes a name underivable — lookup scans the stream's
// registrations).
func TestConsumerARNCarriesItsCreationTimestamp(t *testing.T) {
	store := newLockTestStore(t)
	stream, err := store.CreateStream("consumer_arn", 1, StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}

	first, err := store.RegisterStreamConsumer(stream.StreamARN, "dup", nil)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	pattern := regexp.MustCompile(`^arn:aws.*:kinesis:.*:\d{12}:.*stream/[a-zA-Z0-9_.-]+/consumer/[a-zA-Z0-9_.-]+:[0-9]+$`)
	if !pattern.MatchString(first.ConsumerARN) {
		t.Fatalf("consumer ARN %q does not match the family pattern", first.ConsumerARN)
	}

	if err := store.DeregisterStreamConsumer(first.ConsumerARN); err != nil {
		t.Fatalf("deregister: %v", err)
	}
	// The stamp's granularity is AWS's own — seconds — so the
	// cross-generation ARN distinctness the family documents holds across
	// a second boundary; the gap crosses one before the re-registration.
	// (Inside one second the re-registration reuses the ARN, as AWS's
	// format does — the previous record is gone.)
	time.Sleep(1100 * time.Millisecond)
	second, err := store.RegisterStreamConsumer(stream.StreamARN, "dup", nil)
	if err != nil {
		t.Fatalf("re-register: %v", err)
	}
	if second.ConsumerARN == first.ConsumerARN {
		t.Fatalf("re-registration reused the ARN %q", first.ConsumerARN)
	}

	byName, err := store.GetStreamConsumerByName(stream.StreamARN, "dup")
	if err != nil {
		t.Fatalf("name lookup: %v", err)
	}
	if byName.ConsumerARN != second.ConsumerARN {
		t.Fatalf("name lookup: got %q, want the live registration %q", byName.ConsumerARN, second.ConsumerARN)
	}
}

// TestRegisterStreamConsumerNameUniquePerStream pins the model's
// ConsumerName contract: names are unique within one stream (a duplicate
// registration answers ErrConsumerAlreadyExists — the ARN probe cannot
// catch it, since every registration builds a fresh timestamped ARN) and
// independent across streams.
func TestRegisterStreamConsumerNameUniquePerStream(t *testing.T) {
	store := newLockTestStore(t)
	stream, err := store.CreateStream("consumer_names", 1, StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	other, err := store.CreateStream("consumer_names_other", 1, StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create other stream: %v", err)
	}

	if _, err := store.RegisterStreamConsumer(stream.StreamARN, "shared", nil); err != nil {
		t.Fatalf("first registration: %v", err)
	}
	if _, err := store.RegisterStreamConsumer(stream.StreamARN, "shared", nil); err != ErrConsumerAlreadyExists {
		t.Fatalf("duplicate name on one stream: got %v, want ErrConsumerAlreadyExists", err)
	}
	// The same name on a different stream is a different consumer: the
	// uniqueness scope is the stream, not the account.
	if _, err := store.RegisterStreamConsumer(other.StreamARN, "shared", nil); err != nil {
		t.Fatalf("same name on another stream: %v", err)
	}

	// A deregistered name is free again: the registration that held it is
	// gone from the stream's set.
	consumers, err := store.ListStreamConsumers("consumer_names")
	if err != nil {
		t.Fatalf("list consumers: %v", err)
	}
	if len(consumers) != 1 {
		t.Fatalf("unexpected consumer count: %d", len(consumers))
	}
	if err := store.DeregisterStreamConsumer(consumers[0].ConsumerARN); err != nil {
		t.Fatalf("deregister: %v", err)
	}
	if _, err := store.RegisterStreamConsumer(stream.StreamARN, "shared", nil); err != nil {
		t.Fatalf("re-registration after deregister: %v", err)
	}
}

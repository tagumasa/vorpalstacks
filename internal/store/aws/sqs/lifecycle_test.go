package sqs

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

// The pins below cover the store's lifecycle seams: a panicking cleanup sweep
// is isolated per pass, the receive write pair is transactional (no orphan
// in-flight state when the transaction fails), DeleteQueue reclaims the
// queue's in-memory state, the store-level tag operations require an
// existing queue, and UpdateQueue tolerates a hand-built Queue without the
// raw attribute map.

// failingUpdateStorage wraps a real storage and can be switched to fail
// every transaction, so the receive path's transactional boundary is
// observable from a test.
type failingUpdateStorage struct {
	storage.TransactionalStorage
	fail bool
}

func (f *failingUpdateStorage) Update(ctx context.Context, fn func(txn storage.Transaction) error) error {
	if f.fail {
		return errors.New("injected transaction failure")
	}
	return f.TransactionalStorage.Update(ctx, fn)
}

// failingTagPutStorage fails writes into the tag store's main bucket while
// armed, standing in for a tag-persistence failure at queue creation. The
// flag is consulted per Put because the tag store binds its buckets once at
// construction.
type failingTagPutStorage struct {
	storage.TransactionalStorage
	tagsBucket string
	fail       atomic.Bool
}

func (f *failingTagPutStorage) Bucket(name string) storage.Bucket {
	if name == f.tagsBucket {
		return &failingTagPutBucket{Bucket: f.TransactionalStorage.Bucket(name), fail: &f.fail}
	}
	return f.TransactionalStorage.Bucket(name)
}

type failingTagPutBucket struct {
	storage.Bucket
	fail *atomic.Bool
}

func (b *failingTagPutBucket) Put(k, v []byte) error {
	if b.fail.Load() {
		return errors.New("injected tag write failure")
	}
	return b.Bucket.Put(k, v)
}

// TestCreateQueueTagFailureLeavesNoQueue pins the create's all-or-nothing
// contract: a tag-persistence failure after the queue record committed is
// compensated by removing the record, so no tagless queue survives — the
// documented idempotent retry would otherwise return the existing URL
// without ever re-applying the requested tags.
func TestCreateQueueTagFailureLeavesNoQueue(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	wrapped := &failingTagPutStorage{TransactionalStorage: st, tagsBucket: "sqs-tags-us-east-1"}
	store := NewSQSStore(wrapped, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})

	q := NewQueue("tag-fail", "us-east-1", "123456789012")
	q.Tags = map[string]string{"team": "core"}
	wrapped.fail.Store(true)
	if _, err := store.CreateQueue(q); err == nil {
		t.Fatal("CreateQueue with a failing tag write: expected an error")
	}
	if store.Exists(store.buildQueueURL("tag-fail")) {
		t.Fatal("queue record survived its failed tag write; the queue exists tagless")
	}

	// The immediate retry succeeds in full — no recreate-prohibition marker
	// was written for a queue no client ever observed. A fresh Queue value,
	// because the failed create already stamped the caller's copy with the
	// record's derived attributes.
	wrapped.fail.Store(false)
	retry := NewQueue("tag-fail", "us-east-1", "123456789012")
	retry.Tags = map[string]string{"team": "core"}
	created, err := store.CreateQueue(retry)
	if err != nil {
		t.Fatalf("retry after compensation: %v", err)
	}
	tags, err := store.TagStore.List(created.URL)
	if err != nil {
		t.Fatalf("list retried tags: %v", err)
	}
	if tags["team"] != "core" {
		t.Fatalf("retry's tags = %v, want team=core", tags)
	}
}

// TestCleanupSweepPanicIsolated pins the per-sweep recover: a panicking
// sweep neither propagates nor stops the following sweeps in the same tick.
func TestCleanupSweepPanicIsolated(t *testing.T) {
	store, _ := newDeleteLockTestStore(t)

	panicked := func() {
		store.runSweep("test", func() { panic("injected sweep panic") })
	}
	panicked() // must not propagate

	succeeded := false
	store.runSweep("test", func() { succeeded = true })
	if !succeeded {
		t.Error("the sweep after a panicking sweep did not run; the ticker must stay alive")
	}
}

// TestReceiveWritePairTransactional pins the atomicity of the in-flight
// state and its receipt entry: when the transaction fails, the message
// record keeps its pre-receive state (no orphan in-flight message whose
// handle resolves to nothing) and no receipt entry exists.
func TestReceiveWritePairTransactional(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	wrapped := &failingUpdateStorage{TransactionalStorage: st}
	store := NewSQSStore(wrapped, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() { store.Close() })

	created, err := store.CreateQueue(NewQueue("txn-receive-q", "us-east-1", "123456789012"))
	if err != nil {
		t.Fatalf("create queue: %v", err)
	}
	seed := NewMessage("transactional receive probe")
	sent, err := store.SendMessage(created.URL, seed)
	if err != nil {
		t.Fatalf("seed send: %v", err)
	}

	wrapped.fail = true
	received, err := store.ReceiveMessage(created.URL, 1, nil, 0, "")
	if err != nil {
		t.Fatalf("receive under injected failure: %v", err)
	}
	if len(received) != 0 {
		t.Fatalf("receive returned %d messages under a failed transaction; want 0", len(received))
	}

	// The message record must be untouched: still addressable, not in
	// flight, with no handle that resolves to it.
	var after pb.Message
	if err := store.messagesStore.GetProto(messageKey(created.URL, sent.ID), &after); err != nil {
		t.Fatalf("message record lost after failed transaction: %v", err)
	}
	if after.ReceiptHandle != "" || after.ReceivedAt != nil {
		t.Errorf("failed transaction left orphan in-flight state (handle %q, receivedAt %v)", after.ReceiptHandle, after.ReceivedAt)
	}
	receipts := st.Bucket(store.receiptsBucket)
	found := false
	_ = receipts.ForEach(func(_, v []byte) error {
		if string(v) == messageKey(created.URL, sent.ID) {
			found = true
		}
		return nil
	})
	if found {
		t.Error("receipt entry persisted although the transaction failed")
	}
}

// TestDeleteQueueReclaimsInMemoryState pins the reclamation of the queue's
// in-memory state: sequence counters, receive-attempt cache entries and the
// purge timestamp do not outlive the deleted queue.
func TestDeleteQueueReclaimsInMemoryState(t *testing.T) {
	store, _ := newDeleteLockTestStore(t)

	// The seed message carries FIFO identifiers so the send registers a
	// deduplication entry — the third in-memory map under assertion.
	// Deduplication semantics exist only on FIFO queues (a deduplication
	// id on a standard queue is rejected), so the fixture queue is FIFO.
	fifoQueue := NewQueue("reclaim-pin-q.fifo", "us-east-1", "123456789012")
	fifoQueue.FifoQueue = true
	fifoQueue.ContentBasedDeduplication = true
	queue, err := store.CreateQueue(fifoQueue)
	if err != nil {
		t.Fatalf("create FIFO queue: %v", err)
	}

	// Populate all three maps for this queue.
	msg := NewMessage("reclamation probe")
	msg.MessageGroupID = "g"
	msg.MessageDeduplicationID = "reclaim-key"
	if _, err := store.SendMessage(queue.URL, msg); err != nil {
		t.Fatalf("seed send: %v", err)
	}
	store.receiveAttemptMu.Lock()
	store.receiveAttemptCache[queue.URL+"#attempt-1"] = &receiveAttemptEntry{
		messageIDs: []string{"some-id"},
		createdAt:  time.Now(),
	}
	store.receiveAttemptMu.Unlock()
	store.purgeMutex.Lock()
	store.purgeInProgress[queue.URL] = time.Now()
	store.purgeMutex.Unlock()

	if err := store.DeleteQueue(queue.URL); err != nil {
		t.Fatalf("delete queue: %v", err)
	}

	store.sequenceMu.Lock()
	_, seqLeft := store.sequenceCounters[queue.URL]
	store.sequenceMu.Unlock()
	if seqLeft {
		t.Error("sequence counter survived DeleteQueue")
	}

	store.receiveAttemptMu.Lock()
	attemptsLeft := 0
	for k := range store.receiveAttemptCache {
		if k == queue.URL+"#attempt-1" {
			attemptsLeft++
		}
	}
	store.receiveAttemptMu.Unlock()
	if attemptsLeft != 0 {
		t.Error("receive-attempt cache entry survived DeleteQueue")
	}

	store.purgeMutex.Lock()
	_, purgeLeft := store.purgeInProgress[queue.URL]
	store.purgeMutex.Unlock()
	if purgeLeft {
		t.Error("purge timestamp survived DeleteQueue")
	}
}

// TestStoreTagOpsRequireExistingQueue pins the store-level existence check:
// tag records for unknown queue URLs are rejected, not written.
func TestStoreTagOpsRequireExistingQueue(t *testing.T) {
	store, queue := newDeleteLockTestStore(t)
	ghost := queue.URL + "-ghost"

	if err := store.TagQueue(ghost, map[string]string{"k": "v"}); !errors.Is(err, ErrQueueNotFound) {
		t.Errorf("TagQueue on unknown queue: got %v, want ErrQueueNotFound", err)
	}
	if err := store.UntagQueue(ghost, []string{"k"}); !errors.Is(err, ErrQueueNotFound) {
		t.Errorf("UntagQueue on unknown queue: got %v, want ErrQueueNotFound", err)
	}
	if _, err := store.ListQueueTags(ghost); !errors.Is(err, ErrQueueNotFound) {
		t.Errorf("ListQueueTags on unknown queue: got %v, want ErrQueueNotFound", err)
	}

	if err := store.TagQueue(queue.URL, map[string]string{"k": "v"}); err != nil {
		t.Errorf("TagQueue on existing queue: %v", err)
	}
	tags, err := store.ListQueueTags(queue.URL)
	if err != nil || tags["k"] != "v" {
		t.Errorf("ListQueueTags after TagQueue: tags=%v err=%v", tags, err)
	}
}

// TestUpdateQueueToleratesNilAttributes pins the guard on the raw attribute
// map: a hand-built Queue without it updates without panicking.
func TestUpdateQueueToleratesNilAttributes(t *testing.T) {
	store, queue := newDeleteLockTestStore(t)

	handBuilt := &Queue{
		URL:         queue.URL,
		Name:        queue.Name,
		Region:      queue.Region,
		ARN:         queue.ARN,
		Tags:        map[string]string{},
		Permissions: map[string]*Permission{},
	}
	if err := store.UpdateQueue(handBuilt); err != nil {
		t.Fatalf("UpdateQueue with nil Attributes: %v", err)
	}
	if handBuilt.Attributes["LastModifiedTimestamp"] == "" {
		t.Error("UpdateQueue did not stamp LastModifiedTimestamp")
	}
}

// TestFifoSendDeduplicationRegistrationIsAtomic pins the transactional
// pairing of the FIFO message persist and its deduplication entry: when the
// transaction fails, the send surfaces the error and neither record exists
// — a message can never be queued with its deduplication window silently
// missing (a same-key resend inside the window would be delivered as new).
func TestFifoSendDeduplicationRegistrationIsAtomic(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	wrapped := &failingUpdateStorage{TransactionalStorage: st}
	store := NewSQSStore(wrapped, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() { store.Close() })

	fifo := NewQueue("txn-send-q.fifo", "us-east-1", "123456789012")
	fifo.FifoQueue = true
	fifo.ContentBasedDeduplication = true
	created, err := store.CreateQueue(fifo)
	if err != nil {
		t.Fatalf("create queue: %v", err)
	}
	seed := NewMessage("seed")
	seed.MessageGroupID = "g1"
	if _, err := store.SendMessage(created.URL, seed); err != nil {
		t.Fatalf("seed send: %v", err)
	}

	wrapped.fail = true
	failing := NewMessage("racing send")
	failing.MessageGroupID = "g2"
	if _, err := store.SendMessage(created.URL, failing); err == nil {
		t.Fatal("send under injected transaction failure: expected an error")
	}

	var count int
	_ = common.ForEachAllProto[*pb.Message](store.messagesStore, messagePrefix(created.URL),
		func() *pb.Message { return &pb.Message{} }, nil,
		func(m *pb.Message) error {
			if m.Body == "racing send" {
				count++
			}
			return nil
		})
	if count != 0 {
		t.Fatalf("failed transaction still persisted the message (%d copies)", count)
	}
}

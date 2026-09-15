package sqs

import (
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_sqs"
	"vorpalstacks/internal/store/aws/common"
)

// The two pins below interleave a DeleteQueue with each lock plane's
// in-flight mutator deterministically, by holding the store's own mutex from
// the test (the package's white-box seam): they encode the exact races the
// lock matrix exists to prevent — a queue mutator resurrecting the record
// after the delete, and a receive persisting a message after the wipe.

func newDeleteLockTestStore(t *testing.T) (*SQSStore, *Queue) {
	t.Helper()
	store := newSQSTestStore(t)
	queue := createRedriveQueue(t, store, "lock-pin-q", false)
	return store, queue
}

// blockedFor asserts that a goroutine's completion signal does not fire
// within the grace window.
func blockedFor(t *testing.T, done <-chan error, context string) {
	t.Helper()
	select {
	case <-done:
		t.Fatalf("%s: completed while the lock plane was held", context)
	case <-time.After(200 * time.Millisecond):
	}
}

// TestDeleteQueueWaitsForQueueMutator pins the resurrect-after-delete race:
// a queue mutator that already resolved the queue (here: holds queueMutex
// across its UpdateQueue whole-record write) must finish before DeleteQueue
// runs — the delete cannot interleave between the mutator's read and write
// and leave a resurrected record behind.
func TestDeleteQueueWaitsForQueueMutator(t *testing.T) {
	store, queue := newDeleteLockTestStore(t)

	store.queueMutex.Lock()
	done := make(chan error, 1)
	go func() { done <- store.DeleteQueue(queue.URL) }()

	// The mutator re-Puts the whole queue record — the write a racing
	// SetQueueAttributes or AddPermission performs via UpdateQueue.
	blockedFor(t, done, "DeleteQueue under queueMutex")
	queue.Attributes["VisibilityTimeout"] = "45"
	if err := store.UpdateQueue(queue); err != nil {
		t.Fatalf("mutator UpdateQueue: %v", err)
	}
	store.queueMutex.Unlock()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("DeleteQueue: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("DeleteQueue did not finish after queueMutex was released")
	}

	// The delete ran after the mutator's write: the record is gone, and no
	// later write can resurrect it.
	if store.Exists(queue.URL) {
		t.Error("queue record survived DeleteQueue (resurrected after delete)")
	}
}

// TestDeleteQueueWipeExcludesInFlightReceive pins the receive-during-delete
// race: a receive that holds msgMutex (already scanned the queue, about to
// persist) must finish before the message wipe runs — its persist lands
// before the wipe, so no orphan message or receipt survives in the deleted
// queue.
func TestDeleteQueueWipeExcludesInFlightReceive(t *testing.T) {
	store, queue := newDeleteLockTestStore(t)

	// Seed one message so the queue's message bucket has content the
	// in-flight "receive" can also address.
	seed := NewMessage("seed payload")
	if _, err := store.SendMessage(queue.URL, seed); err != nil {
		t.Fatalf("seed send: %v", err)
	}

	store.msgMutex.Lock()
	done := make(chan error, 1)
	go func() { done <- store.DeleteQueue(queue.URL) }()

	blockedFor(t, done, "DeleteQueue under msgMutex")
	// The receive's persist: a message written after its scan, which the
	// prefix wipe must still cover.
	orphan := NewMessage("persisted after scan")
	orphan.QueueURL = queue.URL
	if err := store.messagesStore.PutProto(messageKey(queue.URL, orphan.ID), MessageToProto(orphan)); err != nil {
		t.Fatalf("in-flight receive persist: %v", err)
	}
	store.msgMutex.Unlock()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("DeleteQueue: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("DeleteQueue did not finish after msgMutex was released")
	}

	var leftover pb.Message
	if err := store.messagesStore.GetProto(messageKey(queue.URL, orphan.ID), &leftover); err == nil {
		t.Error("message persisted during the delete survived the wipe")
	}
	if store.Exists(queue.URL) {
		t.Error("queue record survived DeleteQueue")
	}
}

// TestDeleteQueueWipeExcludesInFlightSend pins the send-during-delete race:
// a send that has not yet resolved the queue (here: parks on the queue
// lock's read side while the delete's critical section runs) must fail with
// QueueNotFound after the wipe — never persist a message (or re-arm
// deduplication state) under the deleted queue's prefix, where retention
// cleanup would never sweep it and a recreated same-name queue would
// deliver it.
func TestDeleteQueueWipeExcludesInFlightSend(t *testing.T) {
	store, queue := newDeleteLockTestStore(t)

	store.queueMutex.Lock()
	done := make(chan error, 1)
	go func() {
		_, err := store.SendMessage(queue.URL, NewMessage("racing send"))
		done <- err
	}()

	blockedFor(t, done, "SendMessage under queueMutex")
	// The delete's critical section: the whole record+messages+receipts+
	// dedup transaction (deleteQueueRecordLocked is deleteQueueRecord with
	// the queue lock the test already holds).
	if _, err := store.deleteQueueRecordLocked(queue.URL); err != nil {
		t.Fatalf("delete critical section: %v", err)
	}
	store.queueMutex.Unlock()

	select {
	case err := <-done:
		if err == nil || !errors.Is(err, ErrQueueNotFound) {
			t.Fatalf("racing send: got %v, want ErrQueueNotFound", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("SendMessage did not finish after queueMutex was released")
	}

	var leftover int
	if err := common.ForEachAllProto[*pb.Message](store.messagesStore, messagePrefix(queue.URL),
		func() *pb.Message { return &pb.Message{} }, nil,
		func(_ *pb.Message) error {
			leftover++
			return nil
		}); err != nil {
		t.Fatalf("leftover scan under the deleted queue's prefix: %v", err)
	}
	if leftover != 0 {
		t.Fatalf("%d messages survived under the deleted queue's prefix", leftover)
	}
}

// failingScanStorage wraps a real storage and can be switched to fail every
// bucket scan, so the delete path's key-collection abort is observable from
// a test.
type failingScanStorage struct {
	storage.TransactionalStorage
	fail bool
}

// failingScanBucket is failingScanStorage's per-bucket view: only the scan
// fails; every other operation passes through to the wrapped bucket.
type failingScanBucket struct {
	storage.Bucket
	fail bool
}

func (f *failingScanStorage) Bucket(name string) storage.Bucket {
	return &failingScanBucket{Bucket: f.TransactionalStorage.Bucket(name), fail: f.fail}
}

func (b *failingScanBucket) ForEach(fn func(k, v []byte) error) error {
	if b.fail {
		return errors.New("injected scan failure")
	}
	return b.Bucket.ForEach(fn)
}

// TestDeleteQueueAbortsWhenKeyScanFails pins the pre-transaction abort: with
// the key scan failing, the deletion returns an error and the queue stays
// fully intact — committing the transaction on an incomplete key set would
// delete the queue record while the messages the scan missed survive under
// the deleted queue's prefix, the half-deleted state the single-transaction
// design exists to prevent.
func TestDeleteQueueAbortsWhenKeyScanFails(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	healthy := NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	queue := createRedriveQueue(t, healthy, "scan-abort", false)
	first, err := healthy.SendMessage(queue.URL, NewMessage("survives the failed scan"))
	if err != nil {
		t.Fatalf("seed message: %v", err)
	}
	zero := int32(0)
	received, err := healthy.ReceiveMessage(queue.URL, 1, &zero, 0, "")
	if err != nil || len(received) != 1 {
		t.Fatalf("receive: %v (%d messages)", err, len(received))
	}
	handle := received[0].ReceiptHandle
	healthy.Close()

	failing := NewSQSStore(&failingScanStorage{TransactionalStorage: st, fail: true},
		"123456789012", "us-east-1", "http://localhost:50080")

	if err := failing.DeleteQueue(queue.URL); err == nil {
		t.Fatal("DeleteQueue succeeded against a failing key scan")
	}
	failing.Close()

	// The assertions re-read through a healthy view: the failing wrapper
	// breaks every scan, including the count path the assertions use.
	verify := NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		verify.Close()
		st.Close()
	})
	if !verify.Exists(queue.URL) {
		t.Fatal("queue record vanished under a failed key scan")
	}
	if visible, notVisible, delayed := verify.GetMessageCounts(queue.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("queue holds %d/%d/%d messages after the aborted delete, want the 1 intact", visible, notVisible, delayed)
	}
	if _, err := verify.resolveQueueReceiptHandle(queue.URL, handle); err != nil {
		t.Fatalf("receipt entry lost under a failed key scan: %v", err)
	}
	var leftover pb.Message
	if err := verify.messagesStore.GetProto(messageKey(queue.URL, first.ID), &leftover); err != nil {
		t.Fatalf("message record lost under a failed key scan: %v", err)
	}
}

// TestDeleteQueueAtomicUnderTransactionFailure pins the deletion's
// transactional atomicity: with the commit itself failing, the queue stays
// fully intact — record, messages, receipt entries — instead of the
// half-deleted state the previous unsequenced series produced (messages
// wiped, record and tags alive).
func TestDeleteQueueAtomicUnderTransactionFailure(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	healthy := NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	queue := createRedriveQueue(t, healthy, "atomic-delete", false)
	first, err := healthy.SendMessage(queue.URL, NewMessage("survives the failed delete"))
	if err != nil {
		t.Fatalf("seed message: %v", err)
	}
	zero := int32(0)
	received, err := healthy.ReceiveMessage(queue.URL, 1, &zero, 0, "")
	if err != nil || len(received) != 1 {
		t.Fatalf("receive: %v (%d messages)", err, len(received))
	}
	handle := received[0].ReceiptHandle
	healthy.Close()

	failing := NewSQSStore(&failingUpdateStorage{TransactionalStorage: st, fail: true},
		"123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		failing.Close()
		st.Close()
	})

	if err := failing.DeleteQueue(queue.URL); err == nil {
		t.Fatal("DeleteQueue succeeded against failing transactions")
	}
	if !failing.Exists(queue.URL) {
		t.Fatal("queue record vanished under a failed deletion transaction")
	}
	if visible, notVisible, delayed := failing.GetMessageCounts(queue.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("queue holds %d/%d/%d messages after the failed delete, want the 1 intact", visible, notVisible, delayed)
	}
	if _, err := failing.resolveQueueReceiptHandle(queue.URL, handle); err != nil {
		t.Fatalf("receipt entry lost under a failed deletion transaction: %v", err)
	}
	var leftover pb.Message
	if err := failing.messagesStore.GetProto(messageKey(queue.URL, first.ID), &leftover); err != nil {
		t.Fatalf("message record lost under a failed deletion transaction: %v", err)
	}
}

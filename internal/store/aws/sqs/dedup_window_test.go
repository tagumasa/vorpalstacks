package sqs

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// The pins below cover the FIFO deduplication interval contract: the
// check-put sequence is atomic per key, the interval starts at the first
// send (a re-read never extends it), expired persisted entries are swept,
// and the interval is independent of message existence ("Amazon SQS
// continues to keep track of the message deduplication ID even after the
// message is received and deleted." — SendMessage API reference).

func newDedupTestStore(t *testing.T, queueName string) (*SQSStore, *Queue) {
	t.Helper()
	store := newSQSTestStore(t)
	queue := NewQueue(queueName, "us-east-1", "123456789012")
	queue.FifoQueue = true
	queue.ContentBasedDeduplication = false
	created, err := store.CreateQueue(queue)
	if err != nil {
		t.Fatalf("create queue: %v", err)
	}
	return store, created
}

func dedupSend(t *testing.T, store *SQSStore, queueURL, body, dedupID string) *Message {
	t.Helper()
	msg := NewMessage(body)
	msg.MessageGroupID = "group-1"
	msg.MessageDeduplicationID = dedupID
	sent, err := store.SendMessage(queueURL, msg)
	if err != nil {
		t.Fatalf("send %q: %v", dedupID, err)
	}
	return sent
}

// TestConcurrentSameDedupKeySendDeliversSingleMessage pins the atomicity of
// the dedup check-put: concurrent sends carrying the same
// MessageDeduplicationId are all acknowledged, with one MessageId, and
// exactly one message exists in the queue.
func TestConcurrentSameDedupKeySendDeliversSingleMessage(t *testing.T) {
	store, queue := newDedupTestStore(t, "dedup-atomic.fifo")

	const senders = 16
	ids := make(chan string, senders)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < senders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			msg := NewMessage("same body")
			msg.MessageGroupID = "group-1"
			msg.MessageDeduplicationID = "race-key"
			sent, err := store.SendMessage(queue.URL, msg)
			if err != nil {
				t.Error("concurrent send failed:", err)
				return
			}
			ids <- sent.ID
		}()
	}
	close(start)
	wg.Wait()
	close(ids)

	seen := make(map[string]bool)
	for id := range ids {
		seen[id] = true
	}
	if len(seen) != 1 {
		t.Errorf("%d concurrent same-key sends acknowledged %d distinct MessageIds; the check-put must be atomic", senders, len(seen))
	}
	visible, notVisible, delayed := store.GetMessageCounts(queue.URL)
	if visible+notVisible+delayed != 1 {
		t.Errorf("queue holds %d messages after %d concurrent same-key sends; want 1", visible+notVisible+delayed, senders)
	}
}

// TestDedupWindowNotExtendedByReread pins the interval semantics: a re-read
// of a persisted entry arms the in-memory TTL from the persisted expiry, so
// a key that keeps being hit still expires at first-send + window. The test
// shortens the persisted expiry to bound the wall time.
func TestDedupWindowNotExtendedByReread(t *testing.T) {
	store, queue := newDedupTestStore(t, "dedup-window.fifo")

	first := dedupSend(t, store, queue.URL, "windowed body", "window-key")
	dedupKey := queue.URL + "#window-key"

	// Force the persisted-read path and shorten the persisted expiry to two
	// seconds from now.
	store.deduplicationMu.Lock()
	delete(store.deduplicationCache, dedupKey)
	store.deduplicationMu.Unlock()
	bucket := store.storage.Bucket(store.dedupBucket)
	val, err := bucket.Get([]byte(dedupKey))
	if err != nil || len(val) == 0 {
		t.Fatalf("persisted dedup entry missing: err=%v len=%d", err, len(val))
	}
	// The persisted form is "<messageKey>\x01<sequenceNumber>\x01<expiryMs>":
	// the rewrite shortens the tail expiry segment only.
	last := bytes.LastIndexByte(val, '\x01')
	if last <= 0 {
		t.Fatalf("persisted dedup entry malformed: %q", val)
	}
	short := time.Now().Add(2 * time.Second).UnixMilli()
	rewritten := append(val[:last+1], []byte(strconv.FormatInt(short, 10))...)
	if err := bucket.Put([]byte(dedupKey), rewritten); err != nil {
		t.Fatalf("rewrite dedup expiry: %v", err)
	}

	// A re-read inside the shortened window hits and must not extend the
	// expiry beyond the persisted value.
	second := dedupSend(t, store, queue.URL, "windowed body", "window-key")
	if second.ID != first.ID {
		t.Fatalf("re-read within the window missed: %s vs %s", second.ID, first.ID)
	}

	time.Sleep(3 * time.Second)

	// Past the persisted expiry the window is over, re-read or not: the
	// next same-key send delivers a new message.
	third := dedupSend(t, store, queue.URL, "windowed body", "window-key")
	if third.ID == first.ID {
		t.Error("the dedup window was extended by the re-read; the interval must end at first-send + window")
	}
}

// TestDeduplicationCleanupSweepsExpiredEntries pins the persisted-bucket
// sweep: expired entries are deleted, fresh and unparseable entries stay
// (the latter remain owned by the expired-re-read path).
func TestDeduplicationCleanupSweepsExpiredEntries(t *testing.T) {
	store, _ := newDedupTestStore(t, "dedup-sweep.fifo")
	bucket := store.storage.Bucket(store.dedupBucket)
	entries := map[string][]byte{
		"k-expired":   []byte(fmt.Sprintf("q\x00m\x01%d", time.Now().Add(-time.Hour).UnixMilli())),
		"k-fresh":     []byte(fmt.Sprintf("q\x00m\x01%d", time.Now().Add(time.Hour).UnixMilli())),
		"k-malformed": []byte("no separator"),
	}
	for k, v := range entries {
		if err := bucket.Put([]byte(k), v); err != nil {
			t.Fatalf("seed dedup entry %s: %v", k, err)
		}
	}

	store.doDeduplicationCleanup()

	for _, k := range []string{"k-fresh", "k-malformed"} {
		if v, _ := bucket.Get([]byte(k)); len(v) == 0 {
			t.Errorf("entry %s was deleted by the sweep; only expired entries go", k)
		}
	}
	if v, _ := bucket.Get([]byte("k-expired")); len(v) > 0 {
		t.Error("expired dedup entry survived the sweep")
	}
}

// TestDedupSuppressesResendAfterDeletion pins the existence independence of
// the interval: after the original message is consumed and its record
// deleted, a same-key resend within the window is still suppressed —
// acknowledged with the original's MessageId, delivering nothing.
func TestDedupSuppressesResendAfterDeletion(t *testing.T) {
	store, queue := newDedupTestStore(t, "dedup-deleted.fifo")

	first := dedupSend(t, store, queue.URL, "deleted body", "deleted-key")
	if err := store.messagesStore.Delete(messageKey(queue.URL, first.ID)); err != nil {
		t.Fatalf("delete original message record: %v", err)
	}

	resent := dedupSend(t, store, queue.URL, "deleted body", "deleted-key")
	if resent.ID != first.ID {
		t.Fatalf("resend after deletion was delivered as a new message (%s); the interval must outlive the original", resent.ID)
	}
	if resent.MD5OfBody != calculateMD5("deleted body") {
		t.Errorf("resend acknowledgement carries MD5 %s; want the request body's digest", resent.MD5OfBody)
	}
	// The synthetic acknowledgment echoes the original's sequence number: the
	// acknowledgment surface of a resend must not depend on whether the
	// original still exists (the recorded branch returns it with the message).
	if first.SequenceNumber == "" {
		t.Fatal("FIFO send carried no SequenceNumber to echo")
	}
	if resent.SequenceNumber != first.SequenceNumber {
		t.Errorf("synthetic acknowledgment SequenceNumber = %q; want the original's %q", resent.SequenceNumber, first.SequenceNumber)
	}

	visible, notVisible, delayed := store.GetMessageCounts(queue.URL)
	if visible+notVisible+delayed != 0 {
		t.Errorf("suppressed resend persisted a message (%d/%d/%d); the queue must stay empty", visible, notVisible, delayed)
	}
}

// TestConcurrentFIFOStampsAgreeOnOrder pins the atomic stamp invariant:
// FIFO delivery order is the SentTimestamp sort, so across concurrent
// same-queue sends the SequenceNumber must be monotonic with the timestamp
// — no pair of messages may show the two orderings disagreeing.
func TestConcurrentFIFOStampsAgreeOnOrder(t *testing.T) {
	store, created := newDedupTestStore(t, "stamp-order.fifo")

	const sends = 32
	type stamp struct {
		sentNano int64
		sequence int64
	}
	stamps := make([]stamp, sends)
	var wg sync.WaitGroup
	for i := 0; i < sends; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msg := NewMessage(fmt.Sprintf("stamp %d", i))
			msg.MessageGroupID = "g" + strconv.Itoa(i%4)
			msg.MessageDeduplicationID = "stamp-" + strconv.Itoa(i)
			sent, err := store.SendMessage(created.URL, msg)
			if err != nil {
				t.Errorf("send %d: %v", i, err)
				return
			}
			seq, err := strconv.ParseInt(sent.SequenceNumber, 10, 64)
			if err != nil {
				t.Errorf("send %d sequence %q: %v", i, sent.SequenceNumber, err)
				return
			}
			stamps[i] = stamp{sentNano: sent.SentTimestamp.UnixNano(), sequence: seq}
		}(i)
	}
	wg.Wait()

	for i := range stamps {
		for j := range stamps {
			if i == j {
				continue
			}
			if stamps[i].sentNano < stamps[j].sentNano && stamps[i].sequence >= stamps[j].sequence {
				t.Fatalf("order inversion: message with timestamp %d has sequence %d, later timestamp %d has lower-or-equal sequence %d",
					stamps[i].sentNano, stamps[i].sequence, stamps[j].sentNano, stamps[j].sequence)
			}
		}
	}
}

// TestMoveToDLQRegistersDedupWindow pins the redrive-path registration: a
// FIFO message redriven into a CBD FIFO dead-letter queue opens the
// destination's deduplication window, so a same-body send into the DLQ
// inside the interval is suppressed as the duplicate of a message the queue
// already holds. The converse is pinned too: a window opened by a direct
// send does not SUPPRESS a later redrive — dropping the redriven copy would
// lose the message outright, which no source licenses.
func TestMoveToDLQRegistersDedupWindow(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "register-dlq.fifo", true)
	source := createRedriveQueue(t, store, "register-src.fifo", true)
	armRedrive(t, store, source, dlq)

	// Registration leg: redrive first, then a direct same-body send is
	// suppressed against the redriven copy.
	redrivePastCount(t, store, source, "registered by the redrive")
	zero := int32(0)
	redriven, recvErr := store.ReceiveMessage(dlq.URL, 1, &zero, 0, "")
	if recvErr != nil || len(redriven) != 1 {
		t.Fatalf("DLQ receive: %v (%d messages)", recvErr, len(redriven))
	}
	directMsg := NewMessage("registered by the redrive")
	directMsg.MessageGroupID = "direct-group"
	direct, err := store.SendMessage(dlq.URL, directMsg)
	if err != nil {
		t.Fatalf("direct send to DLQ: %v", err)
	}
	if direct.ID != redriven[0].ID {
		t.Errorf("same-body send after the redrive delivered as new (%s); the redrive must open the window", direct.ID)
	}
	if visible, notVisible, delayed := store.GetMessageCounts(dlq.URL); visible+notVisible+delayed != 1 {
		t.Errorf("DLQ holds %d/%d/%d messages; the suppressed duplicate must not persist", visible, notVisible, delayed)
	}

	// No-consult leg: a window opened by a direct send does not suppress a
	// later redrive of the same body — both copies exist.
	seedMsg := NewMessage("opened by a direct send")
	seedMsg.MessageGroupID = "direct-group"
	if _, err := store.SendMessage(dlq.URL, seedMsg); err != nil {
		t.Fatalf("seed direct send: %v", err)
	}
	redrivePastCount(t, store, source, "opened by a direct send")
	if visible, notVisible, delayed := store.GetMessageCounts(dlq.URL); visible+notVisible+delayed != 3 {
		t.Errorf("DLQ holds %d/%d/%d messages, want 3 total (two distinct bodies plus the unsuppressed redriven duplicate)", visible, notVisible, delayed)
	}
}

// TestMoveTaskRegistersWithoutConsulting pins the move-task arm of the
// transfer-dedup rule (the policy arm is pinned above): a FIFO message
// redriven into its CBD FIFO DLQ and moved back home inside the original
// send's own deduplication window is DELIVERED — the move registers at the
// destination but never consults, because suppressing the relocation would
// delete the DLQ copy and deliver nothing. The registration leg moves a
// second body to an untouched CBD FIFO queue and pins that the move opened
// that destination's window: a same-body send inside the interval is
// suppressed against the moved copy. The registration leg runs FIRST: the
// move-back's copy returns home with a receive count already past the
// source's armed maxReceiveCount, and a later receive-based redrive on the
// source would sweep it back into the DLQ instead of the leg's own body.
func TestMoveTaskRegistersWithoutConsulting(t *testing.T) {
	store := newRedriveTestStore(t)
	dlq := createRedriveQueue(t, store, "move-task-dedup-dlq.fifo", true)
	source := createRedriveQueue(t, store, "move-task-dedup-src.fifo", true)
	armRedrive(t, store, source, dlq)

	// Registration leg: the destination queue never saw this body before
	// the move, so a live window afterwards can only be the move's own.
	fresh := createRedriveQueue(t, store, "move-task-dedup-dest.fifo", true)
	redrivePastCount(t, store, source, "opens the destination window")
	moved, err := store.StartMessageMoveTask(dlq.ARN, fresh.ARN, 0)
	if err != nil {
		t.Fatalf("start explicit-destination move task: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, moved.TaskId, MoveTaskStatusCompleted)
	if visible, notVisible, delayed := store.GetMessageCounts(fresh.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("destination holds %d/%d/%d messages after the move, want 1", visible, notVisible, delayed)
	}
	directMsg := NewMessage("opens the destination window")
	directMsg.MessageGroupID = "direct-group"
	if _, err := store.SendMessage(fresh.URL, directMsg); err != nil {
		t.Fatalf("direct send to the moved destination: %v", err)
	}
	if visible, notVisible, delayed := store.GetMessageCounts(fresh.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("destination holds %d/%d/%d after the same-body send; the window the move opened must suppress it", visible, notVisible, delayed)
	}

	// No-consult leg: the source's window for this body is live (the send
	// happened seconds ago); a consulting move-back would suppress the
	// relocation and lose the message.
	redrivePastCount(t, store, source, "moves home inside its own window")
	home, err := store.StartMessageMoveTask(dlq.ARN, "", 0)
	if err != nil {
		t.Fatalf("start move-back task: %v", err)
	}
	waitMoveTaskTerminalStatus(t, store, home.TaskId, MoveTaskStatusCompleted)
	if visible, notVisible, delayed := store.GetMessageCounts(source.URL); visible+notVisible+delayed != 1 {
		t.Fatalf("source holds %d/%d/%d messages after the move-back; the relocation must be delivered, not suppressed", visible, notVisible, delayed)
	}
	if visible, notVisible, delayed := store.GetMessageCounts(dlq.URL); visible+notVisible+delayed != 0 {
		t.Fatalf("DLQ holds %d/%d/%d messages after the move-back, want 0", visible, notVisible, delayed)
	}
}

// TestDeduplicationCacheBoundIsReal pins the in-memory cache bound: distinct
// live keys beyond the cap are evicted (earliest-expiring first) instead of
// growing the map without bound — the persisted bucket answers an evicted
// key's next lookup.
func TestDeduplicationCacheBoundIsReal(t *testing.T) {
	store, queue := newDedupTestStore(t, "dedup-bound.fifo")

	for i := 0; i < deduplicationCacheMaxSize+40; i++ {
		dedupSend(t, store, queue.URL, fmt.Sprintf("bound body %d", i), fmt.Sprintf("bound-key-%d", i))
	}

	store.deduplicationMu.RLock()
	size := len(store.deduplicationCache)
	store.deduplicationMu.RUnlock()
	if size > deduplicationCacheMaxSize {
		t.Fatalf("in-memory dedup cache holds %d live entries, want at most the cap %d", size, deduplicationCacheMaxSize)
	}

	// An evicted key still deduplicates through the persisted bucket.
	first := dedupSend(t, store, queue.URL, "bound body 0", "bound-key-0")
	store.deduplicationMu.Lock()
	delete(store.deduplicationCache, queue.URL+"#bound-key-0")
	store.deduplicationMu.Unlock()
	again := dedupSend(t, store, queue.URL, "bound body 0", "bound-key-0")
	if again.ID != first.ID {
		t.Fatalf("evicted key re-delivered as new (%s vs %s); the persisted bucket must stay authoritative", again.ID, first.ID)
	}
}

// TestDeduplicationCacheBoundHoldsOnRearm pins the read-path bound: entries
// re-armed from the persisted bucket (a lookup of a live key) count against
// the in-memory cap exactly as registered entries do — a run of re-read keys
// cannot grow the cache past deduplicationCacheMaxSize.
func TestDeduplicationCacheBoundHoldsOnRearm(t *testing.T) {
	store, queue := newDedupTestStore(t, "rearm-bound.fifo")
	bucket := store.storage.Bucket(store.dedupBucket)
	for i := 0; i < deduplicationCacheMaxSize+1; i++ {
		key := fmt.Sprintf("%s#rearm-key-%03d", queue.URL, i)
		value := deduplicationEntryValue(fmt.Sprintf("%s\x00rearm-msg-%03d", queue.URL, i), "")
		if err := bucket.Put([]byte(key), value); err != nil {
			t.Fatalf("seed persisted entry %d: %v", i, err)
		}
	}

	for i := 0; i < deduplicationCacheMaxSize+1; i++ {
		if _, _, ok := store.getDeduplicationMessageID(fmt.Sprintf("%s#rearm-key-%03d", queue.URL, i)); !ok {
			t.Fatalf("re-arm read %d missed its persisted entry", i)
		}
	}

	store.deduplicationMu.RLock()
	size := len(store.deduplicationCache)
	store.deduplicationMu.RUnlock()
	if size > deduplicationCacheMaxSize {
		t.Fatalf("in-memory dedup cache holds %d entries after re-arm reads, want at most %d", size, deduplicationCacheMaxSize)
	}
}

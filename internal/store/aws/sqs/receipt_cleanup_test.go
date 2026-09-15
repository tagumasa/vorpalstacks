package sqs

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// The sweep must delete only entries older than the retention window;
// fresh entries and keys without an embedded timestamp stay untouched.
func TestReceiptHandleCleanupDeletesOnlyStaleEntries(t *testing.T) {
	store := newSQSTestStore(t)

	receipts := store.storage.Bucket(store.receiptsBucket)
	fresh := fmt.Sprintf("%s#%d", uuid.New(), time.Now().UnixNano())
	stale := fmt.Sprintf("%s#%d", uuid.New(),
		time.Now().Add(-receiptHandleRetention-time.Hour).UnixNano())
	unparseable := "not-a-receipt-handle"
	require.NoError(t, receipts.Put([]byte(fresh), []byte("msg-a")))
	require.NoError(t, receipts.Put([]byte(stale), []byte("msg-b")))
	require.NoError(t, receipts.Put([]byte(unparseable), []byte("msg-c")))

	store.doReceiptHandleCleanup()

	// The storage layer returns (nil, nil) for absent keys, so presence is
	// asserted on the value length.
	if val, _ := receipts.Get([]byte(stale)); len(val) > 0 {
		t.Error("stale receipt entry survived the cleanup sweep")
	}
	if val, _ := receipts.Get([]byte(fresh)); len(val) == 0 {
		t.Error("fresh receipt entry was deleted")
	}
	if val, _ := receipts.Get([]byte(unparseable)); len(val) == 0 {
		t.Error("entry without a timestamp suffix was deleted")
	}
}

// Every receive issues a new handle without removing the previous one, so a
// message received twice has two resolvable entries; deleting via the newest
// handle leaves the old handle usable, and an old-handle delete of an already
// deleted message succeeds idempotently (the documented old-handle
// behaviour: "the request will succeed, but the message might not be
// deleted."). The TTL sweep is what eventually bounds those entries.
func TestOldReceiptHandlesAccumulateAndStayResolvable(t *testing.T) {
	store := newSQSTestStore(t)

	q, err := store.CreateQueue(&Queue{
		Name:                   "receipts-test",
		VisibilityTimeout:      30,
		MaximumMessageSize:     MaxMaximumMessageSize,
		MessageRetentionPeriod: MinMessageRetentionPeriod,
	})
	require.NoError(t, err)

	_, err = store.SendMessage(q.URL, &Message{Body: "hello"})
	require.NoError(t, err)

	// Zero visibility per receive so the second receive redelivers the same
	// message under a fresh handle.
	zeroVisibility := int32(0)
	recv1, err := store.ReceiveMessage(q.URL, 1, &zeroVisibility, 0, "")
	require.NoError(t, err)
	require.Len(t, recv1, 1)
	recv2, err := store.ReceiveMessage(q.URL, 1, &zeroVisibility, 0, "")
	require.NoError(t, err)
	require.Len(t, recv2, 1)

	receipts := store.storage.Bucket(store.receiptsBucket)
	for _, handle := range []string{recv1[0].ReceiptHandle, recv2[0].ReceiptHandle} {
		if val, _ := receipts.Get([]byte(handle)); len(val) == 0 {
			t.Fatalf("receipt entry missing after receive")
		}
	}

	require.NoError(t, store.DeleteMessage(q.URL, recv2[0].ReceiptHandle))
	// The newest handle's entry is consumed, the old one still resolves.
	if val, _ := receipts.Get([]byte(recv2[0].ReceiptHandle)); len(val) > 0 {
		t.Error("used receipt entry was not removed")
	}
	if val, _ := receipts.Get([]byte(recv1[0].ReceiptHandle)); len(val) == 0 {
		t.Fatal("old receipt entry must stay resolvable")
	}
	// Deleting the already-deleted message through the old handle succeeds.
	require.NoError(t, store.DeleteMessage(q.URL, recv1[0].ReceiptHandle))
}

// The deletion-marker sweep must remove markers older than the
// recreate-prohibition window and markers with unparseable timestamps (they
// can no longer answer the window question), while fresh markers stay.
func TestDeletionMarkerCleanupSweepsStaleMarkers(t *testing.T) {
	store := newSQSTestStore(t)

	bucket := store.storage.Bucket(store.deletionsBucket)
	freshURL := "http://localhost:50080/123456789012/fresh-queue"
	staleURL := "http://localhost:50080/123456789012/stale-queue"
	corruptURL := "http://localhost:50080/123456789012/corrupt-queue"
	require.NoError(t, bucket.Put([]byte(freshURL), []byte(strconv.FormatInt(time.Now().Unix(), 10))))
	require.NoError(t, bucket.Put([]byte(staleURL), []byte(strconv.FormatInt(time.Now().Add(-2*queueDeletionWindow).Unix(), 10))))
	require.NoError(t, bucket.Put([]byte(corruptURL), []byte("not-a-timestamp")))

	store.doDeletionMarkerCleanup()

	freshVal, err := bucket.Get([]byte(freshURL))
	require.NoError(t, err, "fresh deletion marker must survive the sweep")
	require.NotEmpty(t, freshVal)
	// The storage layer reports a missing key as an empty read, matching
	// deletedRecently's own absence convention.
	staleVal, _ := bucket.Get([]byte(staleURL))
	require.Empty(t, staleVal, "stale deletion marker must be swept")
	corruptVal, _ := bucket.Get([]byte(corruptURL))
	require.Empty(t, corruptVal, "unparseable deletion marker must be swept")
}

// TestPurgeQueueRemovesReceiptEntries pins the purge's receipts leg: every
// receipt-handle entry resolving into the purged queue's message prefix is
// gone after PurgeQueue — a doubled prefix at the call site would match
// nothing and leave the entries resolvable until the age sweep.
func TestPurgeQueueRemovesReceiptEntries(t *testing.T) {
	store := newSQSTestStore(t)

	created, err := store.CreateQueue(NewQueue("purge-receipts", "us-east-1", "123456789012"))
	require.NoError(t, err)
	_, err = store.SendMessage(created.URL, NewMessage("purged body"))
	require.NoError(t, err)
	zero := int32(0)
	received, err := store.ReceiveMessage(created.URL, 1, &zero, 0, "")
	require.NoError(t, err)
	require.Len(t, received, 1)
	handle := received[0].ReceiptHandle

	require.NoError(t, store.PurgeQueue(created.URL))

	receipts := store.storage.Bucket(store.receiptsBucket)
	got, err := receipts.Get([]byte(handle))
	require.NoError(t, err)
	require.Empty(t, got, "receipt entry survived the purge")
}

package sqs

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vorpalstacks/internal/core/storage"
)

func newReceiptTestStore(t *testing.T) (*SQSStore, func()) {
	t.Helper()
	tmpDir := "./tmp/sqs-receipt-test-" + t.Name()
	require.NoError(t, os.MkdirAll(tmpDir, 0o755))
	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	store := NewSQSStore(s, "123456789012", "us-east-1", "http://localhost:50080")
	cleanup := func() {
		store.Close()
		s.Close()
		os.RemoveAll(tmpDir)
	}
	return store, cleanup
}

// The sweep must delete only entries older than the retention window;
// fresh entries and keys without an embedded timestamp stay untouched.
func TestReceiptHandleCleanupDeletesOnlyStaleEntries(t *testing.T) {
	store, cleanup := newReceiptTestStore(t)
	defer cleanup()

	receipts := store.storage.Bucket("sqs-receipts-us-east-1")
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
	store, cleanup := newReceiptTestStore(t)
	defer cleanup()

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

	receipts := store.storage.Bucket("sqs-receipts-us-east-1")
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

package kinesis

import (
	"errors"
	"testing"
	"time"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// TestGetRecordsReturnsInRetentionRecordsPastAgedWindow is the client-visible
// strand-loop regression: a reader positioned behind a wall of aged-out
// records still receives the in-retention records within its page — the age
// filter runs before the limit cut — and the follow-up iterator advances
// past the aged window instead of re-reading it forever.
func TestGetRecordsReturnsInRetentionRecordsPastAgedWindow(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("strand", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	shardID := "shardId-000000000000"

	for i := 1; i <= 3; i++ {
		if _, err := store.PutAgedRecord("strand", shardID, 48*time.Hour, int64(i)); err != nil {
			t.Fatalf("seed aged record: %v", err)
		}
	}
	fresh, freshShard, err := store.PutRecordWithShardSelection("strand", "pk", "fresh-record", "")
	if err != nil {
		t.Fatalf("put fresh record: %v", err)
	}

	iterator, err := store.CreateShardIterator("strand", freshShard, "TRIM_HORIZON", "", nil)
	if err != nil {
		t.Fatalf("create iterator: %v", err)
	}

	result, err := svc.getRecordsCore(reqCtx, GetRecordsInput{ShardIterator: iterator.IteratorID, Limit: 2})
	if err != nil {
		t.Fatalf("get records behind an aged window: %v", err)
	}
	if len(result.Records) != 1 || result.Records[0].SequenceNumber != fresh.SequenceNumber {
		t.Fatalf("page: got %d records, want the one in-retention record %s", len(result.Records), fresh.SequenceNumber)
	}
	if result.NextShardIterator == nil {
		t.Fatal("open shard returned no NextShardIterator — the reader cannot continue")
	}

	// The follow-up sits after the aged window: reading through it returns
	// only records that arrive later, never the aged ones again.
	next, ok := result.NextShardIterator.(string)
	if !ok {
		t.Fatalf("NextShardIterator type: %T, want string", result.NextShardIterator)
	}
	later, _, err := store.PutRecordWithShardSelection("strand", "pk", "later-record", "")
	if err != nil {
		t.Fatalf("put later record: %v", err)
	}

	follow, err := svc.getRecordsCore(reqCtx, GetRecordsInput{ShardIterator: next, Limit: 10})
	if err != nil {
		t.Fatalf("get records through the follow-up iterator: %v", err)
	}
	if len(follow.Records) != 1 || follow.Records[0].SequenceNumber != later.SequenceNumber {
		t.Fatalf("follow-up page: got %d records, want only the later record %s", len(follow.Records), later.SequenceNumber)
	}
}

// TestGetRecordsDeletedStreamAnswersDocumentedError pins the deleted-stream
// read route end-to-end: an iterator that outlived its stream used to
// nil-panic GetRecords. The delete-time sweep now removes the iterator with
// the stream, so the sequential route answers the iterator's documented
// invalid-iterator error; the guarded stream lookup in getRecordsCore
// covers the remaining interleaving (iterator resolved, stream deleted
// before the stream fetch) and maps it to ResourceNotFoundException.
func TestGetRecordsDeletedStreamAnswersDocumentedError(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("gone", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	iterator, err := store.CreateShardIterator("gone", "shardId-000000000000", "LATEST", "", nil)
	if err != nil {
		t.Fatalf("create iterator: %v", err)
	}

	if err := store.DeleteStream("gone"); err != nil {
		t.Fatalf("delete stream: %v", err)
	}

	_, err = svc.getRecordsCore(reqCtx, GetRecordsInput{ShardIterator: iterator.IteratorID, Limit: 10})
	if !errors.Is(err, ErrInvalidIterator) {
		t.Fatalf("get records on a deleted stream: %v, want the swept iterator's ErrInvalidIterator", err)
	}
}

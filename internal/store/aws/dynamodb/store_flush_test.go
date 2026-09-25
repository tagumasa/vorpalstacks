package dynamodb

import (
	"context"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// recordingUpdateStorage delegates to the wrapped storage and records the
// cancellation state of every Update call's context, so a test can observe
// which context the store handed to each transaction it opens.
type recordingUpdateStorage struct {
	storage.TransactionalStorageWith2PC
	updateCtxErrs []error
}

func (r *recordingUpdateStorage) Update(ctx context.Context, fn func(storage.Transaction) error) error {
	r.updateCtxErrs = append(r.updateCtxErrs, ctx.Err())
	return r.TransactionalStorageWith2PC.Update(ctx, fn)
}

// TestPostCommitContributorFlushSurvivesRequestCancellation pins the
// post-commit contract of FlushContributorWrites: the carrying transaction
// has committed by the time the flush drains, so the flush must not inherit
// the request's cancellation — a request that ends between commit and drain
// must not lose the accounting the commit owes. The cancellation is armed
// inside the transaction function, after the contributor event is queued,
// so the primary Update still enters on a live context while the drain's
// Update (the second the recorder sees) must observe a live one too. The
// persisted table carries no contributor insights, so the drain applies no
// counters — the pin targets the drain's context, not the counting.
func TestPostCommitContributorFlushSurvivesRequestCancellation(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()

	store := NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	createMetricsTestTable(t, store, "FlushTbl")

	rec := &recordingUpdateStorage{TransactionalStorageWith2PC: store.storage}
	store.storage = rec

	ctx, cancel := context.WithCancel(context.Background())
	err = store.Update(ctx, func(txn *DynamoDBTxn) error {
		txn.queueContributorWrite(&Table{Name: "FlushTbl", ContributorInsightsEnabled: true},
			map[string]*AttributeValue{"pk": strAttr("k1")})
		cancel()
		return nil
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	if len(rec.updateCtxErrs) != 2 {
		t.Fatalf("Update calls recorded = %d, want 2 (primary transaction + drain)", len(rec.updateCtxErrs))
	}
	if got := rec.updateCtxErrs[0]; got != nil {
		t.Fatalf("primary transaction entered Update with an already-cancelled context: %v", got)
	}
	if got := rec.updateCtxErrs[1]; got != nil {
		t.Fatalf("post-commit drain observed a cancelled context (the flush describes committed state and must survive the request ending): %v", got)
	}
}

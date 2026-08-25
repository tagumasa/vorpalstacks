package athena

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// testQueryExecutionStore creates a Pebble-backed query execution store
// whose temp directory is cleaned up by the testing.T.
func testQueryExecutionStore(t *testing.T) *QueryExecutionStore {
	t.Helper()
	s, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("storage.Open failed: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return NewQueryExecutionStore(s, "us-east-1")
}

func queuedExecution(id string) *QueryExecution {
	return &QueryExecution{
		QueryExecutionId: id,
		Query:            "SELECT 1",
		Status: &QueryExecutionStatus{
			State:              QueryExecutionStateQueued,
			SubmissionDateTime: time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}
}

// TestCompleteQueryExecution pins the from-state check of the worker's
// terminal write: a completion lands only while the stored state is one
// of the given from states, so a cancellation recorded by a stop request
// in between can never be overwritten.
func TestCompleteQueryExecution(t *testing.T) {
	store := testQueryExecutionStore(t)

	// A completion offered while the execution is still QUEUED, with
	// from=RUNNING, does not land and leaves the state untouched.
	qe := queuedExecution("qe-complete-queued")
	if err := store.CreateQueryExecution(qe); err != nil {
		t.Fatalf("CreateQueryExecution failed: %v", err)
	}
	qe.Status.State = QueryExecutionStateFailed
	persisted, err := store.CompleteQueryExecution(qe, QueryExecutionStateRunning)
	if err != nil {
		t.Fatalf("CompleteQueryExecution failed: %v", err)
	}
	if persisted {
		t.Error("completion persisted from QUEUED with from=RUNNING")
	}
	got, err := store.GetQueryExecution(qe.QueryExecutionId)
	if err != nil {
		t.Fatalf("GetQueryExecution failed: %v", err)
	}
	if got.Status.State != QueryExecutionStateQueued {
		t.Errorf("state = %s, want QUEUED", got.Status.State)
	}

	// A stop request that cancels a RUNNING execution before the worker
	// finishes must survive the worker's terminal write.
	qe = queuedExecution("qe-complete-raced")
	if err := store.CreateQueryExecution(qe); err != nil {
		t.Fatalf("CreateQueryExecution failed: %v", err)
	}
	if _, transitioned, err := store.TransitionQueryExecutionState(qe.QueryExecutionId, QueryExecutionStateRunning, QueryExecutionStateQueued); err != nil || !transitioned {
		t.Fatalf("QUEUED->RUNNING transition = %v, %v", transitioned, err)
	}
	if _, transitioned, err := store.TransitionQueryExecutionState(qe.QueryExecutionId, QueryExecutionStateCancelled, QueryExecutionStateRunning); err != nil || !transitioned {
		t.Fatalf("RUNNING->CANCELLED transition = %v, %v", transitioned, err)
	}
	qe.Status.State = QueryExecutionStateSucceeded
	qe.Statistics = &QueryExecutionStatistics{TotalExecutionTimeInMillis: 42}
	persisted, err = store.CompleteQueryExecution(qe, QueryExecutionStateRunning)
	if err != nil {
		t.Fatalf("CompleteQueryExecution failed: %v", err)
	}
	if persisted {
		t.Error("completion persisted over a recorded CANCELLED state")
	}
	got, err = store.GetQueryExecution(qe.QueryExecutionId)
	if err != nil {
		t.Fatalf("GetQueryExecution failed: %v", err)
	}
	if got.Status.State != QueryExecutionStateCancelled {
		t.Errorf("state = %s, want CANCELLED", got.Status.State)
	}

	// A terminal write offered while the execution is RUNNING lands with
	// the fully populated execution, statistics included.
	qe = queuedExecution("qe-complete-ok")
	if err := store.CreateQueryExecution(qe); err != nil {
		t.Fatalf("CreateQueryExecution failed: %v", err)
	}
	if _, transitioned, err := store.TransitionQueryExecutionState(qe.QueryExecutionId, QueryExecutionStateRunning, QueryExecutionStateQueued); err != nil || !transitioned {
		t.Fatalf("QUEUED->RUNNING transition = %v, %v", transitioned, err)
	}
	qe.Status.State = QueryExecutionStateSucceeded
	qe.Status.CompletionDateTime = time.Date(2027, 1, 1, 12, 0, 5, 0, time.UTC)
	qe.Statistics = &QueryExecutionStatistics{TotalExecutionTimeInMillis: 42}
	persisted, err = store.CompleteQueryExecution(qe, QueryExecutionStateRunning)
	if err != nil {
		t.Fatalf("CompleteQueryExecution failed: %v", err)
	}
	if !persisted {
		t.Error("completion from RUNNING was not persisted")
	}
	got, err = store.GetQueryExecution(qe.QueryExecutionId)
	if err != nil {
		t.Fatalf("GetQueryExecution failed: %v", err)
	}
	if got.Status.State != QueryExecutionStateSucceeded {
		t.Errorf("state = %s, want SUCCEEDED", got.Status.State)
	}
	if got.Statistics == nil || got.Statistics.TotalExecutionTimeInMillis != 42 {
		t.Errorf("statistics = %+v, want TotalExecutionTimeInMillis 42", got.Statistics)
	}
}

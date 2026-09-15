package sqs

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newSQSTestStore opens a throwaway storage and returns a store over it,
// closing both when the test ends — the single scaffold every store-side
// test builds on; domain-specific fixtures (queues, policies, seeds) are the
// callers' own.
func newSQSTestStore(t *testing.T) *SQSStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	store := NewSQSStore(st, "123456789012", "us-east-1", "http://localhost:50080")
	t.Cleanup(func() {
		store.Close()
		st.Close()
	})
	return store
}

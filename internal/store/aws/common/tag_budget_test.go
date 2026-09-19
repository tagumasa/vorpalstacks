package common

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/storage"
)

func newBudgetTagStore(t *testing.T, budget TagBudget) *TagStore {
	t.Helper()
	tmpDir := "./tmp/tag-budget-test-" + t.Name()
	require.NoError(t, os.MkdirAll(tmpDir, 0o755))
	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		s.Close()
		os.RemoveAll(tmpDir)
	})
	return NewTagStore(s, "budget-test", budget)
}

func tagSet(prefix string, n int) map[string]string {
	m := make(map[string]string, n)
	for i := 0; i < n; i++ {
		m[fmt.Sprintf("%s-%02d", prefix, i)] = "v"
	}
	return m
}

// The budget bounds the MERGED set: two sequential writes that each sit
// under the cap but together pass it must leave the second rejected, and
// the rejection carries the service-injected wire identity, not the
// generic sentinel. Overwriting existing keys at the cap costs nothing.
func TestTagBudgetBoundsMergedSet(t *testing.T) {
	ts := newBudgetTagStore(t, StandardTagBudget("ValidationException"))

	require.NoError(t, ts.Tag("res", tagSet("a", 30)))
	require.NoError(t, ts.Tag("res", tagSet("b", 20)))

	err := ts.Tag("res", map[string]string{"extra": "v"})
	require.Error(t, err)
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok, "budget rejection must carry the injected identity, got %T", err)
	require.Equal(t, "ValidationException", apiErr.Code)

	require.NoError(t, ts.Tag("res", map[string]string{"a-00": "rewritten"}))
	tags, err := ts.List("res")
	require.NoError(t, err)
	require.Len(t, tags, 50)
	require.Equal(t, "rewritten", tags["a-00"])
}

// TagFromSlice inherits the budget: the slice form of a write past the cap
// answers with the same identity.
func TestTagBudgetBoundsTagFromSlice(t *testing.T) {
	ts := newBudgetTagStore(t, StandardTagBudget("TagLimitExceeded"))

	require.NoError(t, ts.TagFromSlice("res", sliceOf(tagSet("a", 30))))
	err := ts.TagFromSlice("res", sliceOf(tagSet("b", 21)))
	require.Error(t, err)
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok)
	require.Equal(t, "TagLimitExceeded", apiErr.Code)
}

// Replace swaps the whole set, so its bound is the input itself.
func TestTagBudgetBoundsReplace(t *testing.T) {
	ts := newBudgetTagStore(t, StandardTagBudget("InvalidInput"))

	require.NoError(t, ts.Replace("res", tagSet("a", 10)))
	err := ts.Replace("res", tagSet("b", 51))
	require.Error(t, err)
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok)
	require.Equal(t, "InvalidInput", apiErr.Code)

	require.NoError(t, ts.Replace("res", tagSet("c", 50)))
}

// A zero budget is the no-documented-total service: unbounded merge.
func TestTagBudgetZeroStaysUnbounded(t *testing.T) {
	ts := newBudgetTagStore(t, TagBudget{})

	require.NoError(t, ts.Tag("res", tagSet("a", 80)))
	tags, err := ts.List("res")
	require.NoError(t, err)
	require.Len(t, tags, 80)
}

// A budget without an identity falls back to the generic sentinel, and a
// fresh resource's initial set (the transactional form) enforces the
// caller's explicit bound with that same sentinel.
func TestTagBudgetFallbackAndTxnForm(t *testing.T) {
	ts := newBudgetTagStore(t, TagBudget{MaxKeys: 5})

	err := ts.Tag("res", tagSet("a", 6))
	require.ErrorIs(t, err, ErrTagQuotaExceeded)

	// The transactional initial-set form keeps its explicit bound.
	s, err := storage.Open("./tmp/tag-budget-txn-" + t.Name())
	require.NoError(t, err)
	t.Cleanup(func() {
		s.Close()
		os.RemoveAll("./tmp/tag-budget-txn-" + t.Name())
	})
	bounded := NewTagStore(s, "txn-test", TagBudget{})
	err = s.Update(context.Background(), func(txn storage.Transaction) error {
		return bounded.TagInTxn(txn, "res", tagSet("a", 6), 5)
	})
	require.ErrorIs(t, err, ErrTagQuotaExceeded)
}

// The bound holds against concurrent writers: one hundred goroutines each
// add one distinct key to a fresh resource under a fifty-key budget, and
// exactly fifty land — two writers that each pass a pre-read count cannot
// end above the cap together because the check runs under the store lock.
func TestTagBudgetHoldsUnderConcurrency(t *testing.T) {
	ts := newBudgetTagStore(t, StandardTagBudget("ValidationException"))

	var accepted, rejected atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := ts.Tag("res", map[string]string{fmt.Sprintf("k-%02d", i): "v"}); err != nil {
				rejected.Add(1)
				if _, ok := err.(*awserrors.AWSError); !ok {
					t.Errorf("rejection %d carries no wire identity: %T", i, err)
				}
			} else {
				accepted.Add(1)
			}
		}(i)
	}
	wg.Wait()

	require.Equal(t, int64(50), accepted.Load())
	require.Equal(t, int64(50), rejected.Load())
	tags, err := ts.List("res")
	require.NoError(t, err)
	require.Len(t, tags, 50)
}

func sliceOf(m map[string]string) []types.Tag {
	out := make([]types.Tag, 0, len(m))
	for k, v := range m {
		out = append(out, types.Tag{Key: k, Value: v})
	}
	return out
}

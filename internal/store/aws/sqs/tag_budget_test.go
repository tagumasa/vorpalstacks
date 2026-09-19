package sqs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// The queue tag store enforces the documented fifty on the MERGED set: a
// second write that passes the cap answers the store's TooManyTags
// sentinel (the service table maps it to InvalidParameterValue), and an
// overwrite of an existing key at the cap still lands.
func TestQueueTagBudgetBoundsMergedSet(t *testing.T) {
	store := newSQSTestStore(t)

	q, err := store.CreateQueue(&Queue{
		Name:                   "budget-q",
		VisibilityTimeout:      30,
		MaximumMessageSize:     MaxMaximumMessageSize,
		MessageRetentionPeriod: MinMessageRetentionPeriod,
	})
	require.NoError(t, err)

	first := make(map[string]string, 30)
	for i := 0; i < 30; i++ {
		first[fmt.Sprintf("a-%02d", i)] = "v"
	}
	second := make(map[string]string, 20)
	for i := 0; i < 20; i++ {
		second[fmt.Sprintf("b-%02d", i)] = "v"
	}

	require.NoError(t, store.TagQueue(q.URL, first))
	require.NoError(t, store.TagQueue(q.URL, second))

	err = store.TagQueue(q.URL, map[string]string{"extra": "v"})
	require.ErrorIs(t, err, ErrTooManyTags)

	require.NoError(t, store.TagQueue(q.URL, map[string]string{"a-00": "rewritten"}))
	tags, err := store.ListQueueTags(q.URL)
	require.NoError(t, err)
	require.Len(t, tags, 50)
	require.Equal(t, "rewritten", tags["a-00"])
}

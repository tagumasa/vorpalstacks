package cloudtrail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMutateTrailLoadErrors pins the mutate surface's honest load-error
// contract: a missing record answers ErrTrailNotFound, and any other load
// failure — here a record that does not unmarshal — propagates, so a
// caller can tell an absent trail from a broken read instead of treating
// every failure as the sentinel not-found.
func TestMutateTrailLoadErrors(t *testing.T) {
	store, _ := newPurgeTestStore(t)

	// A missing record is the sentinel.
	_, err := store.MutateTrail("absent", func(*Trail) error { return nil })
	require.ErrorIs(t, err, ErrTrailNotFound)

	// A record that does not unmarshal propagates its load error.
	require.NoError(t, store.BaseStore.PutRaw("corrupt", []byte{0xFF, 0xFE}))
	_, err = store.MutateTrail("corrupt", func(*Trail) error { return nil })
	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrTrailNotFound)
}

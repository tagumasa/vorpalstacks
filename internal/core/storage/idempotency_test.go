package storage

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage/pebbledb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPebbleIdempotencyStore(t *testing.T) {
	db, err := pebbledb.Open(pebbledb.WithPath(t.TempDir()))
	require.NoError(t, err)
	defer db.Close()

	store := NewPebbleIdempotencyStore(db)

	t.Run("CheckAndStore creates new token", func(t *testing.T) {
		isNew, err := store.CheckAndStore("test-token-1", time.Minute)
		assert.NoError(t, err)
		assert.True(t, isNew)
	})

	t.Run("CheckAndStore returns false for duplicate", func(t *testing.T) {
		isNew, err := store.CheckAndStore("test-token-1", time.Minute)
		assert.NoError(t, err)
		assert.False(t, isNew)
	})

	t.Run("StoreResult and GetResult", func(t *testing.T) {
		err := store.StoreResult("test-token-1", []byte("test-result"), time.Minute)
		assert.NoError(t, err)

		result, found, err := store.GetResult("test-token-1")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, []byte("test-result"), result)
	})

	t.Run("GetResult returns false for non-existent token", func(t *testing.T) {
		_, found, err := store.GetResult("non-existent-token")
		assert.NoError(t, err)
		assert.False(t, found)
	})

	t.Run("Delete removes token", func(t *testing.T) {
		err := store.Delete("test-token-1")
		assert.NoError(t, err)

		_, found, err := store.GetResult("test-token-1")
		assert.NoError(t, err)
		assert.False(t, found)
	})
}

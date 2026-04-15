package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPebbleStorage(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	t.Run("Bucket Operations", func(t *testing.T) {
		bucket := s.Bucket("test")

		err := bucket.Put([]byte("key1"), []byte("value1"))
		require.NoError(t, err)

		val, err := bucket.Get([]byte("key1"))
		require.NoError(t, err)
		assert.Equal(t, []byte("value1"), val)

		assert.True(t, bucket.Has([]byte("key1")))
		assert.False(t, bucket.Has([]byte("nonexistent")))

		err = bucket.Delete([]byte("key1"))
		require.NoError(t, err)

		val, err = bucket.Get([]byte("key1"))
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("ForEach", func(t *testing.T) {
		bucket := s.Bucket("foreach-test")
		for i := 0; i < 5; i++ {
			err := bucket.Put([]byte{byte(i)}, []byte{byte(i + 10)})
			require.NoError(t, err)
		}

		count := 0
		err := bucket.ForEach(func(k, v []byte) error {
			count++
			assert.Equal(t, k[0]+10, v[0])
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("ScanPrefix", func(t *testing.T) {
		bucket := s.Bucket("prefix-test")
		bucket.Put([]byte("user:1"), []byte("alice"))
		bucket.Put([]byte("user:2"), []byte("bob"))
		bucket.Put([]byte("item:1"), []byte("apple"))

		iter := bucket.ScanPrefix([]byte("user:"))
		defer iter.Close()

		count := 0
		for iter.Next() {
			count++
			assert.True(t, len(iter.Key()) > 0)
		}
		assert.NoError(t, iter.Error())
		assert.Equal(t, 2, count)
	})

	t.Run("ScanRange", func(t *testing.T) {
		bucket := s.Bucket("range-test")
		bucket.Put([]byte("a"), []byte("1"))
		bucket.Put([]byte("b"), []byte("2"))
		bucket.Put([]byte("c"), []byte("3"))
		bucket.Put([]byte("d"), []byte("4"))

		iter := bucket.ScanRange([]byte("b"), []byte("d"))
		defer iter.Close()

		keys := [][]byte{}
		for iter.Next() {
			keys = append(keys, iter.Key())
		}
		assert.NoError(t, iter.Error())
		assert.Len(t, keys, 2)
	})

	t.Run("Transaction", func(t *testing.T) {
		err := s.Update(context.Background(), func(txn Transaction) error {
			bucket := txn.Bucket("txn-test")
			return bucket.Put([]byte("txn-key"), []byte("txn-value"))
		})
		require.NoError(t, err)

		bucket := s.Bucket("txn-test")
		val, err := bucket.Get([]byte("txn-key"))
		require.NoError(t, err)
		assert.Equal(t, []byte("txn-value"), val)
	})

	t.Run("Rollback", func(t *testing.T) {
		err := s.Update(context.Background(), func(txn Transaction) error {
			bucket := txn.Bucket("rollback-test")
			bucket.Put([]byte("key"), []byte("value"))
			return assert.AnError
		})
		assert.Error(t, err)

		bucket := s.Bucket("rollback-test")
		val, err := bucket.Get([]byte("key"))
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("ListBuckets", func(t *testing.T) {
		b1 := s.Bucket("list-1")
		b1.Put([]byte("key"), []byte("value"))
		b2 := s.Bucket("list-2")
		b2.Put([]byte("key"), []byte("value"))

		buckets := s.ListBuckets()
		assert.Contains(t, buckets, "list-1")
		assert.Contains(t, buckets, "list-2")
	})

	t.Run("DeleteBucket", func(t *testing.T) {
		bucket := s.Bucket("delete-me")
		bucket.Put([]byte("key"), []byte("value"))

		err := s.DeleteBucket("delete-me")
		require.NoError(t, err)

		count := bucket.Count()
		assert.Equal(t, 0, count)
	})

	t.Run("Stats", func(t *testing.T) {
		stats := s.Stats()
		assert.GreaterOrEqual(t, stats.BucketCount, 0)
	})
}

func TestOpenBackend(t *testing.T) {
	t.Run("Pebble Backend", func(t *testing.T) {
		tmpDir := filepath.Join(os.TempDir(), "pebble-backend-test")
		defer os.RemoveAll(tmpDir)

		s, err := OpenBackend(BackendPebble, tmpDir)
		require.NoError(t, err)
		defer s.Close()

		bucket := s.Bucket("test")
		bucket.Put([]byte("key"), []byte("value"))

		val, err := bucket.Get([]byte("key"))
		require.NoError(t, err)
		assert.Equal(t, []byte("value"), val)
	})
}

func TestConditionalBucket(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-conditional-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	t.Run("PutIf NotExists", func(t *testing.T) {
		bucket := s.ConditionalBucket("cond-test-1")

		ok, current, err := bucket.PutIf([]byte("key1"), []byte("value1"), ConditionNotExists())
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Nil(t, current)

		ok, current, err = bucket.PutIf([]byte("key1"), []byte("value2"), ConditionNotExists())
		require.NoError(t, err)
		assert.False(t, ok)
		assert.Equal(t, []byte("value1"), current)
	})

	t.Run("PutIf Equals", func(t *testing.T) {
		bucket := s.ConditionalBucket("cond-test-2")
		bucket.PutIf([]byte("key2"), []byte("initial"), ConditionNotExists())

		ok, current, err := bucket.PutIf([]byte("key2"), []byte("updated"), ConditionEquals([]byte("initial")))
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, []byte("initial"), current)

		ok, current, err = bucket.PutIf([]byte("key2"), []byte("wrong"), ConditionEquals([]byte("initial")))
		require.NoError(t, err)
		assert.False(t, ok)
		assert.Equal(t, []byte("updated"), current)
	})

	t.Run("DeleteIf", func(t *testing.T) {
		bucket := s.ConditionalBucket("cond-test-3")
		bucket.PutIf([]byte("key3"), []byte("to-delete"), ConditionNotExists())

		ok, old, err := bucket.DeleteIf([]byte("key3"), ConditionExists())
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, []byte("to-delete"), old)

		ok, _, err = bucket.DeleteIf([]byte("key3"), ConditionExists())
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("PutReturnOld", func(t *testing.T) {
		bucket := s.ConditionalBucket("cond-test-4")
		bucket.PutIf([]byte("key4"), []byte("old"), ConditionNotExists())

		old, err := bucket.PutReturnOld([]byte("key4"), []byte("new"))
		require.NoError(t, err)
		assert.Equal(t, []byte("old"), old)
	})
}

func TestIdempotencyStore(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-idempotency-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir, WithTTL(true, 500*time.Millisecond))
	require.NoError(t, err)
	defer s.Close()

	store := s.IdempotencyStore()

	t.Run("CheckAndStore", func(t *testing.T) {
		token := "token-1"

		isNew, err := store.CheckAndStore(token, 10*time.Minute)
		require.NoError(t, err)
		assert.True(t, isNew)

		isNew, err = store.CheckAndStore(token, 10*time.Minute)
		require.NoError(t, err)
		assert.False(t, isNew)
	})

	t.Run("StoreAndGetResult", func(t *testing.T) {
		token := "token-2"
		result := []byte(`{"status":"success"}`)

		isNew, err := store.CheckAndStore(token, 10*time.Minute)
		require.NoError(t, err)
		assert.True(t, isNew)

		err = store.StoreResult(token, result, 10*time.Minute)
		require.NoError(t, err)

		got, exists, err := store.GetResult(token)
		require.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, result, got)
	})

	t.Run("GetResultNotCompleted", func(t *testing.T) {
		token := "token-3"

		_, _ = store.CheckAndStore(token, 10*time.Minute)

		_, exists, err := store.GetResult(token)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		token := "token-4"

		_, _ = store.CheckAndStore(token, 10*time.Minute)
		_ = store.StoreResult(token, []byte("result"), 10*time.Minute)

		err := store.Delete(token)
		require.NoError(t, err)

		_, exists, err := store.GetResult(token)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("TTL Expiry", func(t *testing.T) {
		token := "token-5"

		isNew, err := store.CheckAndStore(token, 2*time.Second)
		require.NoError(t, err)
		assert.True(t, isNew)

		time.Sleep(3 * time.Second)

		isNew, err = store.CheckAndStore(token, 10*time.Minute)
		require.NoError(t, err)
		assert.True(t, isNew)
	})
}

func TestTwoPhaseTransaction(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-twophase-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	t.Run("Successful Commit", func(t *testing.T) {
		txn := s.TwoPhaseTransaction()

		validationCalled := false
		executionCalled := false

		txn.AddValidator(ValidatorFunc(func(ctx context.Context, t Transaction) error {
			validationCalled = true
			return nil
		}))

		txn.AddExecutor(ExecutorFunc(func(ctx context.Context, t Transaction) error {
			executionCalled = true
			bucket := t.Bucket("twophase-test-1")
			return bucket.Put([]byte("key1"), []byte("value1"))
		}))

		err := txn.Commit(context.Background())
		require.NoError(t, err)
		assert.True(t, validationCalled)
		assert.True(t, executionCalled)

		bucket := s.Bucket("twophase-test-1")
		val, err := bucket.Get([]byte("key1"))
		require.NoError(t, err)
		assert.Equal(t, []byte("value1"), val)
	})

	t.Run("Validation Failure", func(t *testing.T) {
		txn := s.TwoPhaseTransaction()

		executionCalled := false

		txn.AddValidator(ValidatorFunc(func(ctx context.Context, t Transaction) error {
			return errors.New("validation failed")
		}))

		txn.AddExecutor(ExecutorFunc(func(ctx context.Context, t Transaction) error {
			executionCalled = true
			return nil
		}))

		err := txn.Commit(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.False(t, executionCalled)
	})

	t.Run("Execution Failure", func(t *testing.T) {
		txn := s.TwoPhaseTransaction()

		txn.AddValidator(ValidatorFunc(func(ctx context.Context, t Transaction) error {
			return nil
		}))

		txn.AddExecutor(ExecutorFunc(func(ctx context.Context, t Transaction) error {
			return errors.New("execution failed")
		}))

		err := txn.Commit(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "execution failed")
	})

	t.Run("Multiple Validators and Executors", func(t *testing.T) {
		txn := s.TwoPhaseTransaction()

		order := []string{}

		txn.AddValidator(ValidatorFunc(func(ctx context.Context, t Transaction) error {
			order = append(order, "v1")
			return nil
		}))

		txn.AddValidator(ValidatorFunc(func(ctx context.Context, t Transaction) error {
			order = append(order, "v2")
			return nil
		}))

		txn.AddExecutor(ExecutorFunc(func(ctx context.Context, t Transaction) error {
			order = append(order, "e1")
			return nil
		}))

		txn.AddExecutor(ExecutorFunc(func(ctx context.Context, t Transaction) error {
			order = append(order, "e2")
			return nil
		}))

		err := txn.Commit(context.Background())
		require.NoError(t, err)
		assert.Equal(t, []string{"v1", "v2", "e1", "e2"}, order)
	})

	t.Run("Clear and Reuse", func(t *testing.T) {
		txn := s.TwoPhaseTransaction()

		txn.AddValidator(ValidatorFunc(func(ctx context.Context, t Transaction) error {
			return nil
		}))
		txn.AddExecutor(ExecutorFunc(func(ctx context.Context, t Transaction) error {
			return nil
		}))

		assert.Equal(t, 1, txn.ValidatorCount())
		assert.Equal(t, 1, txn.ExecutorCount())

		txn.Clear()

		assert.Equal(t, 0, txn.ValidatorCount())
		assert.Equal(t, 0, txn.ExecutorCount())
	})
}

func TestMultiItemTransaction(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-multiitem-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	t.Run("Multiple Puts", func(t *testing.T) {
		txn := s.MultiItemTransaction()

		txn.Put("multi-test-1", []byte("key1"), []byte("value1"))
		txn.Put("multi-test-1", []byte("key2"), []byte("value2"))
		txn.Put("multi-test-1", []byte("key3"), []byte("value3"))

		err := txn.Commit(context.Background())
		require.NoError(t, err)

		bucket := s.Bucket("multi-test-1")
		val1, _ := bucket.Get([]byte("key1"))
		val2, _ := bucket.Get([]byte("key2"))
		val3, _ := bucket.Get([]byte("key3"))

		assert.Equal(t, []byte("value1"), val1)
		assert.Equal(t, []byte("value2"), val2)
		assert.Equal(t, []byte("value3"), val3)
	})

	t.Run("ConditionCheck Fails", func(t *testing.T) {
		bucket := s.Bucket("multi-test-2")
		bucket.Put([]byte("existing"), []byte("original"))

		txn := s.MultiItemTransaction()
		txn.ConditionCheck("multi-test-2", []byte("existing"), ConditionEquals([]byte("wrong")))
		txn.Put("multi-test-2", []byte("key"), []byte("value"))

		err := txn.Commit(context.Background())
		require.Error(t, err)

		var condErr *ConditionalCheckFailedError
		assert.True(t, errors.As(err, &condErr))
	})

	t.Run("ConditionCheck Succeeds", func(t *testing.T) {
		bucket := s.Bucket("multi-test-3")
		bucket.Put([]byte("existing"), []byte("original"))

		txn := s.MultiItemTransaction()
		txn.ConditionCheck("multi-test-3", []byte("existing"), ConditionEquals([]byte("original")))
		txn.Put("multi-test-3", []byte("new-key"), []byte("new-value"))

		err := txn.Commit(context.Background())
		require.NoError(t, err)

		val, _ := bucket.Get([]byte("new-key"))
		assert.Equal(t, []byte("new-value"), val)
	})

	t.Run("Delete Operations", func(t *testing.T) {
		bucket := s.Bucket("multi-test-4")
		bucket.Put([]byte("to-delete"), []byte("value"))

		txn := s.MultiItemTransaction()
		txn.Delete("multi-test-4", []byte("to-delete"))

		err := txn.Commit(context.Background())
		require.NoError(t, err)

		val, _ := bucket.Get([]byte("to-delete"))
		assert.Nil(t, val)
	})
}

func TestVersionedBucket(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-versioning-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	t.Run("PutWithVersion and GetLatest", func(t *testing.T) {
		bucket := s.VersionedBucket("ver-test-1")

		vv1, err := bucket.PutWithVersion([]byte("key1"), []byte("value1"))
		require.NoError(t, err)
		assert.NotZero(t, vv1.Version)
		assert.False(t, vv1.Deleted)
		assert.Equal(t, []byte("value1"), vv1.Value)

		vv2, err := bucket.PutWithVersion([]byte("key1"), []byte("value2"))
		require.NoError(t, err)
		assert.Greater(t, vv2.Version, vv1.Version)

		latest, err := bucket.GetLatest([]byte("key1"))
		require.NoError(t, err)
		assert.Equal(t, vv2.Version, latest.Version)
		assert.Equal(t, []byte("value2"), latest.Value)
	})

	t.Run("GetVersion", func(t *testing.T) {
		bucket := s.VersionedBucket("ver-test-2")

		vv1, _ := bucket.PutWithVersion([]byte("key2"), []byte("v1"))
		bucket.PutWithVersion([]byte("key2"), []byte("v2"))

		old, err := bucket.GetVersion([]byte("key2"), vv1.Version)
		require.NoError(t, err)
		assert.Equal(t, []byte("v1"), old.Value)
	})

	t.Run("DeleteWithVersion", func(t *testing.T) {
		bucket := s.VersionedBucket("ver-test-3")
		bucket.PutWithVersion([]byte("key3"), []byte("value"))

		vv, err := bucket.DeleteWithVersion([]byte("key3"))
		require.NoError(t, err)
		assert.True(t, vv.Deleted)

		latest, err := bucket.GetLatest([]byte("key3"))
		require.NoError(t, err)
		assert.True(t, latest.Deleted)
	})

	t.Run("ListVersions", func(t *testing.T) {
		bucket := s.VersionedBucket("ver-test-4")

		for i := 0; i < 5; i++ {
			bucket.PutWithVersion([]byte("key4"), []byte{byte(i)})
		}

		versions, err := bucket.ListVersions([]byte("key4"), VersionListOptions{})
		require.NoError(t, err)
		assert.Len(t, versions, 5)

		versionsLimited, err := bucket.ListVersions([]byte("key4"), VersionListOptions{Limit: 2})
		require.NoError(t, err)
		assert.Len(t, versionsLimited, 2)
	})

	t.Run("PurgeVersions", func(t *testing.T) {
		bucket := s.VersionedBucket("ver-test-5")

		for i := 0; i < 5; i++ {
			bucket.PutWithVersion([]byte("key5"), []byte{byte(i)})
		}

		err := bucket.PurgeVersions([]byte("key5"), 2)
		require.NoError(t, err)

		versions, err := bucket.ListVersions([]byte("key5"), VersionListOptions{})
		require.NoError(t, err)
		assert.Len(t, versions, 2)
	})
}

func TestLockManager(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "pebble-lock-test")
	defer os.RemoveAll(tmpDir)

	s, err := Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	t.Run("TryLock and Unlock", func(t *testing.T) {
		lm := s.LockManager()

		handle, err := lm.TryLock([]byte("lock1"), LockModeExclusive, 10*time.Minute)
		require.NoError(t, err)
		assert.NotEmpty(t, handle.Token)

		_, err = lm.TryLock([]byte("lock1"), LockModeExclusive, 10*time.Minute)
		require.Error(t, err)
		assert.IsType(t, &LockConflictError{}, err)

		err = lm.Unlock(handle)
		require.NoError(t, err)

		handle2, err := lm.TryLock([]byte("lock1"), LockModeExclusive, 10*time.Minute)
		require.NoError(t, err)
		assert.NotEmpty(t, handle2.Token)
	})

	t.Run("Lock with context", func(t *testing.T) {
		lm := s.LockManager()

		_, _ = lm.TryLock([]byte("lock2"), LockModeExclusive, 10*time.Minute)

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		_, err := lm.Lock(ctx, []byte("lock2"), LockModeExclusive, 10*time.Minute)
		require.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("IsLocked", func(t *testing.T) {
		lm := s.LockManager()

		locked, _, err := lm.IsLocked([]byte("lock3"))
		require.NoError(t, err)
		assert.False(t, locked)

		lm.TryLock([]byte("lock3"), LockModeExclusive, 10*time.Minute)

		locked, handle, err := lm.IsLocked([]byte("lock3"))
		require.NoError(t, err)
		assert.True(t, locked)
		assert.NotNil(t, handle)
	})

	t.Run("Extend lock", func(t *testing.T) {
		lm := s.LockManager()

		handle, _ := lm.TryLock([]byte("lock4"), LockModeExclusive, 1*time.Second)
		originalExpiry := handle.ExpiresAt

		err := lm.Extend(handle, 10*time.Minute)
		require.NoError(t, err)
		assert.Greater(t, handle.ExpiresAt, originalExpiry)
	})

	t.Run("Unlock with wrong token", func(t *testing.T) {
		lm := s.LockManager()

		_, _ = lm.TryLock([]byte("lock5"), LockModeExclusive, 10*time.Minute)

		wrongHandle := &LockHandle{
			Key:   []byte("lock5"),
			Token: "wrong-token",
		}

		err := lm.Unlock(wrongHandle)
		require.Error(t, err)
		assert.IsType(t, &LockTokenMismatchError{}, err)
	})
}

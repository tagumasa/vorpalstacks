package pebbledb

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

func TestOpen(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.Close()
	assert.NoError(t, err)
}

func TestSetGet(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	key := []byte("testkey")
	value := []byte("testvalue")

	err = db.Set(key, value)
	require.NoError(t, err)

	got, err := db.Get(key)
	require.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestGetNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Get([]byte("nonexistent"))
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	key := []byte("testkey")
	value := []byte("testvalue")

	err = db.Set(key, value)
	require.NoError(t, err)

	err = db.Delete(key)
	require.NoError(t, err)

	_, err = db.Get(key)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestSetGetJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	type TestStruct struct {
		Name  string
		Value int
	}

	key := []byte("jsonkey")
	value := TestStruct{Name: "test", Value: 42}

	err = db.SetJSON(key, value)
	require.NoError(t, err)

	var got TestStruct
	err = db.GetJSON(key, &got)
	require.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestEncryption(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	db, err := Open(WithPath(dbPath), WithEncryption(key))
	require.NoError(t, err)
	defer db.Close()

	testKey := []byte("secretkey")
	testValue := []byte("secretvalue")

	err = db.Set(testKey, testValue)
	require.NoError(t, err)

	got, err := db.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, got)

	rawVal, err := db.GetRaw(testKey)
	require.NoError(t, err)
	assert.NotEqual(t, testValue, rawVal)
}

func TestAESGCMEncryptor(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	enc, err := NewAESGCMEncryptor(key)
	require.NoError(t, err)

	plaintext := []byte("hello world")

	ciphertext, err := enc.Encrypt(plaintext)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, ciphertext)

	decrypted, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestAESGCMEncryptorInvalidKey(t *testing.T) {
	_, err := NewAESGCMEncryptor([]byte("short"))
	assert.ErrorIs(t, err, ErrInvalidKeySize)
}

func TestTTL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TTL test in short mode")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath), WithTTL(true, 500*time.Millisecond, 1*time.Second))
	require.NoError(t, err)
	defer db.Close()

	key := []byte("ttlkey")
	value := []byte("ttlvalue")

	err = db.SetWithTTL(key, value, 1*time.Second)
	require.NoError(t, err)

	got, err := db.Get(key)
	require.NoError(t, err)
	assert.Equal(t, value, got)

	time.Sleep(2 * time.Second)

	db.cleanExpiredKeys()

	_, err = db.Get(key)
	assert.ErrorIs(t, err, ErrKeyNotFound)
}

func TestListKeys(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	prefix := []byte("prefix:")
	keys := [][]byte{
		[]byte("prefix:a"),
		[]byte("prefix:b"),
		[]byte("prefix:c"),
		[]byte("other:d"),
	}

	for _, k := range keys {
		err := db.Set(k, []byte("value"))
		require.NoError(t, err)
	}

	gotKeys, err := db.ListKeys(prefix)
	require.NoError(t, err)
	assert.Len(t, gotKeys, 3)

	for _, k := range gotKeys {
		assert.True(t, len(k) >= 7 && string(k[:7]) == "prefix:")
	}
}

func TestDeleteRange(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	keys := [][]byte{
		[]byte("key:a"),
		[]byte("key:b"),
		[]byte("key:c"),
		[]byte("key:d"),
	}

	for _, k := range keys {
		err := db.Set(k, []byte("value"))
		require.NoError(t, err)
	}

	err = db.DeleteRange([]byte("key:b"), []byte("key:d"))
	require.NoError(t, err)

	for _, k := range keys {
		_, err := db.Get(k)
		keyStr := string(k)
		if keyStr == "key:a" || keyStr == "key:d" {
			assert.NoError(t, err, "key %s should exist", k)
		} else {
			assert.ErrorIs(t, err, ErrKeyNotFound, "key %s should be deleted", k)
		}
	}
}

func TestClosedDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)

	_, err = db.Get([]byte("key"))
	assert.ErrorIs(t, err, ErrClosed)

	err = db.Set([]byte("key"), []byte("value"))
	assert.ErrorIs(t, err, ErrClosed)
}

func TestNoopEncryptor(t *testing.T) {
	enc := NewNoopEncryptor()

	plaintext := []byte("hello")

	ciphertext, err := enc.Encrypt(plaintext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, ciphertext)

	decrypted, err := enc.Decrypt(ciphertext)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestTTLManager(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	tm := db.TTLManager()

	key := []byte("ttlkey")
	value := []byte("ttlvalue")

	err = tm.SetWithExpiry(key, value, 1*time.Hour)
	require.NoError(t, err)

	got, exists, err := tm.Get(key)
	require.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, value, got)

	err = tm.Delete(key)
	require.NoError(t, err)

	_, exists, err = tm.Get(key)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestEncryptionWithEnvKey(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	envKey := "VORPALHAMSTER_TEST_KEY"
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i + 1)
	}

	os.Setenv(envKey, string(testKey))
	defer os.Unsetenv(envKey)

	db, err := Open(WithPath(dbPath), WithEncryption(testKey))
	require.NoError(t, err)
	defer db.Close()

	key := []byte("testkey")
	value := []byte("testvalue")

	err = db.Set(key, value)
	require.NoError(t, err)

	got, err := db.Get(key)
	require.NoError(t, err)
	assert.Equal(t, value, got)
}

func TestCompact(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	for i := 0; i < 100; i++ {
		key := []byte("key" + string(rune(i)))
		err := db.Set(key, []byte("value"))
		require.NoError(t, err)
	}

	for i := 0; i < 50; i++ {
		key := []byte("key" + string(rune(i)))
		err := db.Delete(key)
		require.NoError(t, err)
	}

	err = db.Compact(context.Background())
	assert.NoError(t, err)
}

func TestMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	db, err := Open(WithPath(dbPath))
	require.NoError(t, err)
	defer db.Close()

	err = db.Set([]byte("key"), []byte("value"))
	require.NoError(t, err)

	metrics := db.Metrics()
	assert.NotNil(t, metrics)
}

func TestClosedOperations(t *testing.T) {
	tests := []struct {
		name string
		fn   func(db *DB) error
	}{
		{
			name: "delete",
			fn:   func(db *DB) error { return db.Delete([]byte("key")) },
		},
		{
			name: "delete_range",
			fn:   func(db *DB) error { return db.DeleteRange([]byte("a"), []byte("b")) },
		},
		{
			name: "list_keys",
			fn: func(db *DB) error {
				_, err := db.ListKeys([]byte("prefix"))
				return err
			},
		},
		{
			name: "set_raw",
			fn:   func(db *DB) error { return db.SetRaw([]byte("key"), []byte("value")) },
		},
		{
			name: "get_raw",
			fn: func(db *DB) error {
				_, err := db.GetRaw([]byte("key"))
				return err
			},
		},
		{
			name: "compact",
			fn:   func(db *DB) error { return db.Compact(context.Background()) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			dbPath := filepath.Join(tmpDir, "testdb")

			db, err := Open(WithPath(dbPath))
			require.NoError(t, err)

			err = db.Close()
			require.NoError(t, err)

			err = tt.fn(db)
			assert.ErrorIs(t, err, ErrClosed)
		})
	}
}

func TestTTLValue(t *testing.T) {
	data := []byte("test data")
	expiresAt := int64(1234567890)

	ttlVal := &TTLValue{
		ExpiresAt: expiresAt,
		Data:      data,
	}

	marshaled, err := ttlVal.Marshal()
	require.NoError(t, err)

	unmarshaled, err := UnmarshalTTLValue(marshaled)
	require.NoError(t, err)
	assert.Equal(t, expiresAt, unmarshaled.ExpiresAt)
	assert.Equal(t, data, unmarshaled.Data)
}

func TestUnmarshalTTLValueInvalid(t *testing.T) {
	_, err := UnmarshalTTLValue([]byte("short"))
	assert.Error(t, err)
}

func TestErrors(t *testing.T) {
	assert.True(t, errors.Is(ErrClosed, ErrClosed))
	assert.True(t, errors.Is(ErrKeyNotFound, ErrKeyNotFound))
	assert.True(t, errors.Is(ErrInvalidTTLKey, ErrInvalidTTLKey))
	assert.True(t, errors.Is(ErrEncryptionDisabled, ErrEncryptionDisabled))
	assert.True(t, errors.Is(ErrInvalidKeySize, ErrInvalidKeySize))
}

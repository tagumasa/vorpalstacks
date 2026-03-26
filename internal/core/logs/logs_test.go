package logs

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
)

func TestLogger(t *testing.T) {
	t.Run("Basic Logging", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&Config{
			Level:  LevelDebug,
			Output: &buf,
		})

		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")

		output := buf.String()
		assert.Contains(t, output, "DEBUG")
		assert.Contains(t, output, "INFO")
		assert.Contains(t, output, "WARN")
		assert.Contains(t, output, "ERROR")
	})

	t.Run("With Fields", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&Config{
			Level:  LevelInfo,
			Output: &buf,
		})

		logger.Info("test message", String("key1", "value1"), Int("key2", 123))

		output := buf.String()
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "key1=value1")
		assert.Contains(t, output, "key2=123")
	})

	t.Run("With Chained Logger", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&Config{
			Level:  LevelInfo,
			Output: &buf,
		})

		chained := logger.With(String("service", "test")).WithSource("handler.go")
		chained.Info("chained message")

		output := buf.String()
		assert.Contains(t, output, "service=test")
	})

	t.Run("Level Filtering", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&Config{
			Level:  LevelWarn,
			Output: &buf,
		})

		logger.Debug("should not appear")
		logger.Info("should not appear")
		logger.Warn("should appear")
		logger.Error("should appear")

		output := buf.String()
		assert.NotContains(t, output, "should not appear")
		assert.Contains(t, output, "should appear")
	})
}

func TestLogStore(t *testing.T) {
	tmpDir := os.TempDir() + "/logs-test"
	defer os.RemoveAll(tmpDir)

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	defer s.Close()

	store := NewLogStore(s)

	t.Run("Write and Get", func(t *testing.T) {
		entry := &LogEntry{
			Level:   LevelInfo,
			Message: "test message",
			Fields:  []Field{{Key: "key1", Value: "value1"}},
		}

		err := store.Write(entry)
		require.NoError(t, err)
		assert.NotEmpty(t, entry.ID)

		got, err := store.Get("default", "local", entry.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, entry.Message, got.Message)
		assert.Equal(t, entry.Level, got.Level)
	})

	t.Run("Query by Time", func(t *testing.T) {
		now := time.Now()
		entry := &LogEntry{
			Level:     LevelInfo,
			Message:   "time query test",
			Timestamp: now,
		}

		err := store.Write(entry)
		require.NoError(t, err)

		start := now.Add(-time.Hour)
		end := now.Add(time.Hour)
		entries, err := store.Query(*NewQueryFilter().WithTimeRange(start, end))
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(entries), 1)
	})

	t.Run("Query by Level", func(t *testing.T) {
		entry := &LogEntry{
			Level:   LevelError,
			Message: "error entry",
		}

		err := store.Write(entry)
		require.NoError(t, err)

		entries, err := store.Query(*NewQueryFilter().WithLevel(LevelError))
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(entries), 1)
		for _, e := range entries {
			assert.Equal(t, LevelError, e.Level)
		}
	})

	t.Run("Delete Before", func(t *testing.T) {
		oldEntry := &LogEntry{
			Level:     LevelInfo,
			Message:   "old entry",
			Timestamp: time.Now().Add(-48 * time.Hour),
		}
		newEntry := &LogEntry{
			Level:     LevelInfo,
			Message:   "new entry",
			Timestamp: time.Now(),
		}

		err := store.Write(oldEntry)
		require.NoError(t, err)
		err = store.Write(newEntry)
		require.NoError(t, err)

		cutoff := time.Now().Add(-24 * time.Hour)
		err = store.DeleteBefore(cutoff)
		require.NoError(t, err)

		got, err := store.Get("default", "local", oldEntry.ID)
		require.NoError(t, err)
		assert.Nil(t, got)

		got, err = store.Get("default", "local", newEntry.ID)
		require.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestKey(t *testing.T) {
	t.Run("Key Encoding", func(t *testing.T) {
		key := NewKey("logs", "tenant1", "us-east-1", "entry123")
		encoded := key.Encode()
		assert.Equal(t, "logs:tenant1:us-east-1:entry123", encoded)

		decoded := DecodeKey(encoded)
		require.NotNil(t, decoded)
		assert.Equal(t, "logs", decoded.Namespace)
		assert.Equal(t, "tenant1", decoded.TenantID)
		assert.Equal(t, "us-east-1", decoded.Region)
		assert.Equal(t, "entry123", decoded.ID)
	})

	t.Run("Log Key", func(t *testing.T) {
		key := NewLogKey("tenant1", "region1", "id123")
		assert.Equal(t, "logs", key.Namespace)
		assert.Equal(t, "tenant1", key.TenantID)
	})
}

func TestIndexKey(t *testing.T) {
	t.Run("Time Index Key", func(t *testing.T) {
		key := NewTimeIndexKey("tenant1", "region1", "2024-02-15:10", "entry123")
		encoded := key.Encode()
		assert.Contains(t, encoded, "idx_time")
		assert.Contains(t, encoded, "tenant1")
		assert.Contains(t, encoded, "2024-02-15:10")
	})

	t.Run("Level Index Key", func(t *testing.T) {
		key := NewLevelIndexKey("tenant1", "region1", "ERROR", 123456, "entry123")
		encoded := key.Encode()
		assert.Contains(t, encoded, "idx_level")
		assert.Contains(t, encoded, "ERROR")
	})
}

func TestRotation(t *testing.T) {
	t.Run("Rotation Config", func(t *testing.T) {
		cfg := DefaultRotationConfig()
		assert.Equal(t, 7*24*time.Hour, cfg.MaxAge)
		assert.Equal(t, int64(100*1024*1024), cfg.MaxSize)
	})

	t.Run("Rotator", func(t *testing.T) {
		tmpDir := os.TempDir() + "/rotation-test"
		defer os.RemoveAll(tmpDir)

		s, err := storage.Open(tmpDir)
		require.NoError(t, err)
		defer s.Close()

		store := NewLogStore(s)
		rotator := NewRotator(store, &RotationConfig{
			MaxAge: time.Hour,
		})

		err = rotator.RunOnce()
		assert.NoError(t, err)
	})
}

func TestLevel(t *testing.T) {
	assert.Equal(t, "DEBUG", LevelDebug.String())
	assert.Equal(t, "INFO", LevelInfo.String())
	assert.Equal(t, "WARN", LevelWarn.String())
	assert.Equal(t, "ERROR", LevelError.String())

	assert.Equal(t, LevelDebug, ParseLevel("debug"))
	assert.Equal(t, LevelInfo, ParseLevel("INFO"))
	assert.Equal(t, LevelWarn, ParseLevel("warning"))
}

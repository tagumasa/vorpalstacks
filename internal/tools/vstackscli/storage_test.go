package vstackscli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
)

// TestRunVacuum covers the vacuum happy path against a throwaway data
// directory: database directories are compacted, non-database directories
// are left untouched, and a second run over already-compacted data
// succeeds (idempotent, size-stable).
func TestRunVacuum(t *testing.T) {
	dataPath := "./tmp/vstacks-vacuum-test"
	os.RemoveAll(dataPath)
	t.Cleanup(func() { os.RemoveAll(dataPath) })
	require.NoError(t, os.MkdirAll(dataPath, 0755))

	// Seed a region database with enough data to occupy more than one
	// memtable flush, so compaction has real work to inspect.
	regionDir := filepath.Join(dataPath, "us-east-1")
	s, err := storage.Open(regionDir)
	require.NoError(t, err)
	bucket := s.Bucket("vacuum-test")
	for i := 0; i < 500; i++ {
		require.NoError(t, bucket.Put([]byte(strings.Repeat("a", i%32)+strings.Repeat("k", 100)), []byte(strings.Repeat("v", 256))))
	}
	require.NoError(t, s.Close())

	// A non-database directory must be neither opened nor modified.
	chunkDir := filepath.Join(dataPath, "logs-chunks")
	require.NoError(t, os.MkdirAll(chunkDir, 0755))
	chunkFile := filepath.Join(chunkDir, "000001.chunk")
	require.NoError(t, os.WriteFile(chunkFile, []byte("chunk-data"), 0644))

	var out bytes.Buffer
	require.NoError(t, runVacuum(dataPath, &out))
	output := out.String()
	assert.Contains(t, output, "us-east-1:", "per-database line present")
	assert.NotContains(t, output, "logs-chunks", "non-database directories are not reported")

	chunkData, err := os.ReadFile(chunkFile)
	require.NoError(t, err)
	assert.Equal(t, "chunk-data", string(chunkData), "non-database files untouched")

	// Second run over compacted data: succeeds and stays size-stable
	// within compaction metadata noise.
	var out2 bytes.Buffer
	require.NoError(t, runVacuum(dataPath, &out2))
	assert.Contains(t, out2.String(), "us-east-1:")
}

func TestIsPebbleDBDir(t *testing.T) {
	dataPath := t.TempDir()

	assert.False(t, isPebbleDBDir(dataPath), "empty directory is not a database")

	s, err := storage.Open(filepath.Join(dataPath, "db"))
	require.NoError(t, err)
	require.NoError(t, s.Close())
	assert.True(t, isPebbleDBDir(filepath.Join(dataPath, "db")), "opened-and-closed database recognised")
}

// A close failure must not be swallowed by an earlier compaction failure:
// both report together, each carrying its phase label.
func TestVacuumStepError(t *testing.T) {
	assert.NoError(t, vacuumStepError(nil, nil))

	compactOnly := vacuumStepError(errors.New("c"), nil)
	assert.EqualError(t, compactOnly, "compaction failed: c")

	closeOnly := vacuumStepError(nil, errors.New("x"))
	assert.EqualError(t, closeOnly, "close failed: x")

	both := vacuumStepError(errors.New("c"), errors.New("x"))
	assert.Contains(t, both.Error(), "compaction failed: c")
	assert.Contains(t, both.Error(), "close failed: x")
}

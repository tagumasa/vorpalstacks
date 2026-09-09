package storage

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Object writes create directory nodes (MkdirAll per prefix segment);
// deletes must reclaim them, or the blob root accumulates one empty
// skeleton per deleted prefix forever. Bucket deletion removes the
// bucket's whole tree.
func TestHybridBlobDeleteReclaimsDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	bs, err := Open(tmpDir)
	require.NoError(t, err)
	defer bs.Close()
	s, err := NewHybridBlobStore(bs, tmpDir)
	require.NoError(t, err)
	ctx := context.Background()

	// Above the small-object threshold so the blob lands in the file
	// tier the pruning operates on.
	large := bytes.Repeat([]byte("x"), int(SmallObjectThreshold)+1)

	// Two blobs share a prefix directory: deleting one leaf must leave
	// the shared chain in place for the survivor.
	_, err = s.Put(ctx, "prune-a", "deep/nested/leaf.txt", bytes.NewReader(large), nil)
	require.NoError(t, err)
	_, err = s.Put(ctx, "prune-a", "deep/sibling.txt", bytes.NewReader(large), nil)
	require.NoError(t, err)
	require.NoError(t, s.Delete(ctx, "prune-a", "deep/nested/leaf.txt"))
	require.DirExists(t, filepath.Join(tmpDir, "blobs", "prune-a", "d:deep"))
	require.NoDirExists(t, filepath.Join(tmpDir, "blobs", "prune-a", "d:deep", "d:nested"))

	// Deleting the survivor prunes the chain down to the blob root.
	require.NoError(t, s.Delete(ctx, "prune-a", "deep/sibling.txt"))
	require.NoDirExists(t, filepath.Join(tmpDir, "blobs", "prune-a"))

	// The versioned delete path prunes the same way.
	_, err = s.PutWithVersion(ctx, "prune-b", "v/k", "v1", bytes.NewReader(large), nil)
	require.NoError(t, err)
	require.NoError(t, s.DeleteWithVersion(ctx, "prune-b", "v/k", "v1"))
	require.NoDirExists(t, filepath.Join(tmpDir, "blobs", "prune-b"))

	// Bucket deletion removes the whole tree, empty or not.
	_, err = s.Put(ctx, "prune-c", "x/y", bytes.NewReader(large), nil)
	require.NoError(t, err)
	require.NoError(t, s.DeleteBucket(ctx, "prune-c"))
	require.NoDirExists(t, filepath.Join(tmpDir, "blobs", "prune-c"))
}

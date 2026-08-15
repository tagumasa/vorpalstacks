package cloudwatchlogs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage/chunk"
)

// When the header read of a freshly flushed chunk file fails, the file
// must be removed: it will never be indexed, and no sweep knows about
// unindexed files, so keeping it leaks storage for ever.
func TestWriteChunkFileRemovesOrphanOnHeaderReadFailure(t *testing.T) {
	s := newLogsTestStore(t)

	orig := readChunkHeader
	readChunkHeader = func(string) (*chunk.Header, error) {
		return nil, errors.New("header read failed")
	}
	t.Cleanup(func() { readChunkHeader = orig })

	entries := []LogEntry{{
		Timestamp:     time.Now().UnixMilli(),
		Message:       "orphan candidate",
		IngestionTime: time.Now().UnixMilli(),
	}}

	if _, _, err := s.writeChunkFile(entries); err == nil {
		t.Fatal("writeChunkFile reported success despite the header failure")
	}

	remaining, readErr := os.ReadDir(s.chunksDir)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return // the directory itself never appeared: nothing leaked
		}
		t.Fatal(readErr)
	}
	for _, f := range remaining {
		if !f.IsDir() {
			t.Fatalf("orphaned chunk file leaked: %s", filepath.Join(s.chunksDir, f.Name()))
		}
	}
}

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

// Chunk IDs must sort lexicographically in chronological order even when
// the entry counts and sequence numbers have different digit widths:
// variable-width components let a later chunk (seq 8, 10 entries) sort
// before an earlier one (seq 2, 7 entries) at the same timestamp.
func TestChunkIDLexicographicOrderMatchesChronological(t *testing.T) {
	s := newLogsTestStore(t)

	group := "order-group"
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: "stream"}); err != nil {
		t.Fatal(err)
	}

	ts := int64(1721237155000)
	// Same timestamp, different entry counts, forced sequence order.
	if _, err := s.PutLogEvents(group, "stream", []LogEntry{
		{Timestamp: ts, Message: "first-batch-1", IngestionTime: ts},
		{Timestamp: ts, Message: "first-batch-2", IngestionTime: ts},
		{Timestamp: ts, Message: "first-batch-3", IngestionTime: ts},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PutLogEvents(group, "stream", []LogEntry{
		{Timestamp: ts, Message: "second-batch-1", IngestionTime: ts},
		{Timestamp: ts, Message: "second-batch-2", IngestionTime: ts},
	}); err != nil {
		t.Fatal(err)
	}

	events, _, _, err := s.GetLogEvents(group, "stream", 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 5 {
		t.Fatalf("got %d events, want 5", len(events))
	}
	wantOrder := []string{
		"first-batch-1", "first-batch-2", "first-batch-3",
		"second-batch-1", "second-batch-2",
	}
	for i, e := range events {
		if e.Message != wantOrder[i] {
			t.Fatalf("event %d = %q, want %q (full order: %v)", i, e.Message, wantOrder, messages(events))
		}
	}
}

func messages(events []*OutputLogEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.Message
	}
	return out
}

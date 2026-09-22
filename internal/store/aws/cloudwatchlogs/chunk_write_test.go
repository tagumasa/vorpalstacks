package cloudwatchlogs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
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
// before an earlier one (seq 2, 7 entries) at the same timestamp. The
// read engine does not derive event order from the chunk scan any more
// (it sorts by (timestamp, ingestion time, stream, message digest,
// chunk-and-position ordinal) — the documented key with the ordinal
// standing in for the per-put request ID), so the chronological
// invariant is pinned on the chunk listing itself; the event read is
// pinned for completeness of the served set, not the tie order.
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

	chunks, err := s.ListChunksForStream(group, "stream")
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 2 {
		t.Fatalf("got %d chunks, want 2", len(chunks))
	}
	if chunks[0].EntryCount != 3 || chunks[1].EntryCount != 2 {
		t.Fatalf("chunk order not chronological: entry counts = [%d, %d], want [3, 2]",
			chunks[0].EntryCount, chunks[1].EntryCount)
	}
	if chunks[0].ChunkID >= chunks[1].ChunkID {
		t.Fatalf("chunk IDs not in chronological order: %q >= %q", chunks[0].ChunkID, chunks[1].ChunkID)
	}

	events, _, _, err := s.GetLogEvents(group, "stream", 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 5 {
		t.Fatalf("got %d events, want 5", len(events))
	}
	seen := make(map[string]bool, len(events))
	for _, e := range events {
		seen[e.Message] = true
	}
	for _, want := range []string{"first-batch-1", "first-batch-2", "first-batch-3", "second-batch-1", "second-batch-2"} {
		if !seen[want] {
			t.Fatalf("event %q missing from the read (served: %v)", want, messages(events))
		}
	}
}

// A maximal PutLogEvents batch — exactly MaxChunkSize events, the
// documented per-batch ceiling — must commit and read back: the writer's
// auto-flush boundary sits above the batch size, so the caller's explicit
// Flush still sees the buffered entries and returns a real path instead
// of the empty path an intervening auto-flush would leave behind.
func TestWriteChunkFileMaximalBatch(t *testing.T) {
	s := newLogsTestStore(t)

	group := "maximal-batch-group"
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: "stream"}); err != nil {
		t.Fatal(err)
	}

	ts := time.Now().UnixMilli()
	entries := make([]LogEntry, MaxChunkSize)
	for i := range entries {
		entries[i] = LogEntry{Timestamp: ts, Message: "m", IngestionTime: ts}
	}
	if _, err := s.PutLogEvents(group, "stream", entries); err != nil {
		t.Fatalf("maximal batch of exactly %d events failed: %v", MaxChunkSize, err)
	}

	events, _, _, err := s.GetLogEvents(group, "stream", 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != MaxChunkSize {
		t.Fatalf("got %d events back, want %d", len(events), MaxChunkSize)
	}
}

func messages(events []*OutputLogEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.Message
	}
	return out
}

// A restart replays its in-memory identity counters from zero, so the
// store must derive both counters' floors from the durable state it
// inherits: the index-key sequence from the chunk index, the file-name
// sequence from the chunk files themselves. Replaying the same backfill
// (same first-event timestamp, same batch size, same ordinal position)
// after a restart must issue fresh identities — the earlier chunk's index
// record survives (its events stay readable) and the earlier file is not
// truncated by a recreated name.
func TestChunkIdentitySurvivesRestartReplay(t *testing.T) {
	storageDir := t.TempDir()
	dataDir := t.TempDir()

	batch := func() []LogEntry {
		ts := int64(1721237155000)
		return []LogEntry{
			{Timestamp: ts, Message: "replay-1", IngestionTime: ts},
			{Timestamp: ts, Message: "replay-2", IngestionTime: ts},
		}
	}

	open := func() (*Store, *storage.PebbleStorage) {
		st, err := storage.Open(storageDir)
		if err != nil {
			t.Fatal(err)
		}
		s, err := NewStore(st, st.Bucket("logs-us-east-1"), "000000000000", "us-east-1", dataDir)
		if err != nil {
			t.Fatal(err)
		}
		return s, st
	}

	s1, st1 := open()
	if err := s1.CreateLogGroup(&LogGroup{Name: "replay-group"}); err != nil {
		t.Fatal(err)
	}
	if err := s1.CreateLogStream(&LogStream{LogGroupName: "replay-group", Name: "s"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s1.PutLogEvents("replay-group", "s", batch()); err != nil {
		t.Fatal(err)
	}
	if err := st1.Close(); err != nil {
		t.Fatal(err)
	}

	// The restart: a fresh store over the same durable state — the same
	// index storage and the same chunks directory.
	s2, st2 := open()
	t.Cleanup(func() { st2.Close() })
	if _, err := s2.PutLogEvents("replay-group", "s", batch()); err != nil {
		t.Fatal(err)
	}

	chunks, err := s2.ListChunksForStream("replay-group", "s")
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 2 {
		t.Fatalf("restart replay listed %d chunks, want 2 (the replay must not overwrite the first chunk's index record)", len(chunks))
	}
	if chunks[0].ChunkID == chunks[1].ChunkID {
		t.Fatalf("replayed batch reused chunk ID %q", chunks[0].ChunkID)
	}

	events, _, _, err := s2.GetLogEvents("replay-group", "s", 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 4 {
		t.Fatalf("read %d events after the replay, want 4", len(events))
	}
	counts := map[string]int{}
	for _, e := range events {
		counts[e.Message]++
	}
	for _, m := range []string{"replay-1", "replay-2"} {
		if counts[m] != 2 {
			t.Fatalf("message %q appeared %d times, want 2 (the first chunk's file must not be truncated)", m, counts[m])
		}
	}

	for _, c := range chunks {
		if _, err := os.Stat(filepath.Join(dataDir, "logs-chunks", filepath.Base(c.ChunkPath))); err != nil {
			t.Fatalf("chunk file %s missing: %v", c.ChunkPath, err)
		}
	}
}

// A fallback (non-transactional) PutLogEvents whose commit fails midway
// must leave no drift: no chunk index record, no orphaned chunk file, and
// the LogGroup/LogStream records back at their pre-call state — the
// sequential writes that landed are compensated, so a failed call is
// invisible to the accounting.
func TestFallbackCommitLeavesNoDrift(t *testing.T) {
	for _, tc := range []struct {
		name       string
		failPrefix string
	}{
		{"chunk index write fails", keyPrefixChunk},
		{"log stream write fails", keyPrefixLogStream},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, arm := newFallbackStore(t, tc.failPrefix, true, false)
			const group, stream = "fb", "s"
			if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
				t.Fatal(err)
			}
			if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: stream}); err != nil {
				t.Fatal(err)
			}
			arm()

			ts := time.Now().UnixMilli()
			if _, err := s.PutLogEvents(group, stream, []LogEntry{
				{Timestamp: ts, Message: "drifting", IngestionTime: ts},
			}); err == nil {
				t.Fatal("PutLogEvents reported success despite the injected commit failure")
			}

			lg, err := s.GetLogGroup(group)
			if err != nil {
				t.Fatal(err)
			}
			if lg.StoredBytes != 0 {
				t.Fatalf("StoredBytes drifted to %d after the failed commit, want 0", lg.StoredBytes)
			}
			ls, err := s.GetLogStream(group, stream)
			if err != nil {
				t.Fatal(err)
			}
			if ls.UploadSequenceToken != "0" {
				t.Fatalf("UploadSequenceToken drifted to %q after the failed commit, want %q", ls.UploadSequenceToken, "0")
			}
			if ls.FirstEventTs != 0 || ls.LastEventTs != 0 {
				t.Fatalf("event timestamps drifted after the failed commit: first=%d last=%d", ls.FirstEventTs, ls.LastEventTs)
			}
			chunks, err := s.ListChunksForStream(group, stream)
			if err != nil {
				t.Fatal(err)
			}
			if len(chunks) != 0 {
				t.Fatalf("a chunk index record survived the failed commit: %d listed", len(chunks))
			}
			remaining, err := os.ReadDir(s.chunksDir)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range remaining {
				if !f.IsDir() {
					t.Fatalf("orphaned chunk file survived the failed commit: %s", f.Name())
				}
			}
		})
	}
}

// Every chunk read resolves through the safe-path check, absolute forms
// included: a stored ChunkPath is the writer's relative form, so an
// absolute or escaping value is a corrupt record and reads nothing.
// TestReadChunkFileValidatesEveryPath pins both recorded chunk-path
// forms and the rejection rows. The chunk writer returns an absolute
// path, and the index records the relative form when the chunks
// directory allows the conversion — a server started with a relative
// data path cannot convert, so its records carry the absolute form and
// the read must open it as named. Every form resolves through the
// safe-path check: a traversal or a path outside the chunks directory
// is refused on either form.
func TestReadChunkFileValidatesEveryPath(t *testing.T) {
	s := newLogsTestStore(t)
	ts := time.Now().UnixMilli()
	actualPath, _, err := s.writeChunkFile([]LogEntry{
		{Timestamp: ts, Message: "absolute form", IngestionTime: ts},
		{Timestamp: ts, Message: "relative form", IngestionTime: ts + 1},
	})
	if err != nil {
		t.Fatalf("write chunk: %v", err)
	}

	// The absolute form the writer hands its caller (the server-recorded
	// shape) reads back in full, and the seam resolves it to the very
	// file — the identity the removal paths' os.Remove relies on.
	entries, err := s.readChunkFile(actualPath)
	if err != nil {
		t.Fatalf("read back the absolute chunk path: %v", err)
	}
	if len(entries) != 2 || entries[0].Message != "absolute form" {
		t.Fatalf("absolute-form entries = %+v", entries)
	}
	resolved, err := s.safeChunkPath(actualPath)
	if err != nil {
		t.Fatalf("safe-path the absolute chunk path: %v", err)
	}
	if resolved != filepath.Clean(actualPath) {
		t.Fatalf("safeChunkPath(absolute) = %q, want the file as named", resolved)
	}

	// The relative form the index records when the conversion succeeds
	// (the unit-recorded shape) reads back through the same seam.
	rel, err := filepath.Rel(s.chunksDir, actualPath)
	if err != nil {
		t.Fatalf("relative form: %v", err)
	}
	if _, err := s.readChunkFile(rel); err != nil {
		t.Fatalf("read back the relative chunk path: %v", err)
	}

	if _, err := s.readChunkFile("/etc/passwd"); err == nil {
		t.Fatal("absolute path outside the chunks directory must not read")
	}
	if _, err := s.readChunkFile("../../../../etc/passwd"); err == nil {
		t.Fatal("escaping path read must not succeed")
	}
}

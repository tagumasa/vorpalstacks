package cloudwatchlogs

import (
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newLogsTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	s, err := NewStore(st, st.Bucket("logs-us-east-1"), "000000000000", "us-east-1", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// TestPurgeExpiredChunksAdvancesStreamBoundaries pins that the retention
// purge advances the affected streams' boundary members to the surviving
// data: firstEventTimestamp and lastEventTimestamp describe stored
// events, so a purged era must not keep its stamps on the record — a
// partially purged stream keeps its surviving era's bounds and a fully
// purged one returns to the never-written shape.
func TestPurgeExpiredChunksAdvancesStreamBoundaries(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "g"}); err != nil {
		t.Fatal(err)
	}
	for _, stream := range []string{"partial", "fully"} {
		if err := s.CreateLogStream(&LogStream{LogGroupName: "g", Name: stream}); err != nil {
			t.Fatal(err)
		}
	}
	const old, newer = 1_000_000_000_000, 1_700_000_000_000
	if _, err := s.PutLogEvents("g", "partial", []LogEntry{{Timestamp: old, Message: "ancient"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PutLogEvents("g", "partial", []LogEntry{{Timestamp: newer, Message: "recent"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PutLogEvents("g", "fully", []LogEntry{{Timestamp: old, Message: "ancient"}}); err != nil {
		t.Fatal(err)
	}

	const cutoff = 1_600_000_000_000
	if _, err := s.PurgeExpiredChunks("g", cutoff); err != nil {
		t.Fatal(err)
	}

	partial, err := s.GetLogStream("g", "partial")
	if err != nil {
		t.Fatal(err)
	}
	if partial.FirstEventTs != newer || partial.LastEventTs != newer {
		t.Fatalf("partial: first=%d last=%d, want the surviving era %d", partial.FirstEventTs, partial.LastEventTs, newer)
	}
	fully, err := s.GetLogStream("g", "fully")
	if err != nil {
		t.Fatal(err)
	}
	if fully.FirstEventTs != 0 || fully.LastEventTs != 0 {
		t.Fatalf("fully purged: first=%d last=%d, want the never-written shape", fully.FirstEventTs, fully.LastEventTs)
	}

	events, _, _, err := s.GetLogEvents("g", "partial", 0, 0, 10, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Message != "recent" {
		t.Fatalf("surviving events = %v", events)
	}
}

// TestListLogStreamsMarkerRoundTrip pins the marker contract the listing
// documents: the NextMarker a page returns is the boundary stream NAME,
// so feeding it back resumes strictly after that stream. A raw store key
// in that position double-prefixes on the next call and silently
// truncates every marker-driven walk (stream teardown, fetch-all
// consumers).
func TestListLogStreamsMarkerRoundTrip(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"s-1", "s-2", "s-3", "s-4", "s-5"}
	for _, name := range want {
		if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: name}); err != nil {
			t.Fatal(err)
		}
	}

	var got []string
	marker := ""
	for {
		streams, nextMarker, err := s.ListLogStreams("group", "", marker, 2)
		if err != nil {
			t.Fatal(err)
		}
		for _, ls := range streams {
			got = append(got, ls.Name)
		}
		if nextMarker == "" {
			break
		}
		if strings.HasPrefix(nextMarker, "log-stream:") {
			t.Fatalf("NextMarker leaked a raw store key: %q", nextMarker)
		}
		marker = nextMarker
	}

	if len(got) != len(want) {
		t.Fatalf("marker walk must traverse every stream, got %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("marker walk must preserve order: want %v, got %v", want, got)
		}
	}
}

// A partial delete that removed a chunk file but left its index record
// behind must not wedge the log group: DeleteLogGroup tolerates the
// missing file and clears the orphaned index record so a retry always
// makes progress.
func TestDeleteLogGroupToleratesPartiallyDeletedChunks(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: "stream"}); err != nil {
		t.Fatal(err)
	}

	idxKey := s.chunkIndexKey("group", "stream", "chunk-1")
	if err := s.PutProto(idxKey, ChunkMetaToProto(&ChunkMeta{
		ChunkID:      "chunk-1",
		LogGroupName: "group",
		LogStream:    "stream",
		ChunkPath:    "missing-file",
	})); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatalf("delete failed on partially-deleted group: %v", err)
	}

	if _, err := s.GetLogGroup("group"); err == nil {
		t.Fatal("log group still exists after delete")
	}
	if s.Exists(idxKey) {
		t.Fatal("orphaned chunk index record survived the delete")
	}
}

// Deleting a group that owns streams and filters removes every trace:
// stream records, chunk index records, filters, and the group itself.
func TestDeleteLogGroupRemovesStreamsAndFilters(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: "stream-a"}); err != nil {
		t.Fatal(err)
	}
	idxKey := s.chunkIndexKey("group", "stream-a", "chunk-1")
	if err := s.PutProto(idxKey, ChunkMetaToProto(&ChunkMeta{
		ChunkID:      "chunk-1",
		LogGroupName: "group",
		LogStream:    "stream-a",
		ChunkPath:    "missing-file",
	})); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatal(err)
	}

	if _, err := s.GetLogStream("group", "stream-a"); err == nil {
		t.Fatal("log stream still exists after group delete")
	}
	if s.Exists(idxKey) {
		t.Fatal("chunk index record survived the group delete")
	}
	if _, err := s.GetLogGroup("group"); err == nil {
		t.Fatal("log group still exists after delete")
	}
}

// A teardown that fails part-way must keep the group record: deleting it
// anyway would orphan the surviving sub-resources for ever, because a
// retry stops at the missing group before it can reach them. Once the
// blocker is gone the retry must succeed.
func TestDeleteLogGroupKeepsRecordWhenSubresourceFails(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: "group", Name: "stream-a"}); err != nil {
		t.Fatal(err)
	}
	// A chunk index record whose path escapes the chunks directory makes
	// the chunk teardown fail deterministically.
	badKey := s.chunkIndexKey("group", "stream-a", "chunk-bad")
	if err := s.PutProto(badKey, ChunkMetaToProto(&ChunkMeta{
		ChunkID:      "chunk-bad",
		LogGroupName: "group",
		LogStream:    "stream-a",
		ChunkPath:    "../escape",
	})); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroup("group"); err == nil {
		t.Fatal("delete reported success despite a failing sub-resource")
	}
	if _, err := s.GetLogGroup("group"); err != nil {
		t.Fatalf("group record removed while sub-resources survived: %v", err)
	}

	// Remove the blocker and retry: the retry must reach the survivors.
	if err := s.Delete(badKey); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatalf("retry after clearing the blocker failed: %v", err)
	}
	if _, err := s.GetLogGroup("group"); err == nil {
		t.Fatal("log group still exists after successful retry")
	}
	if _, err := s.GetLogStream("group", "stream-a"); err == nil {
		t.Fatal("log stream orphaned by the partial delete")
	}
}

// A metric-filter listing failure must keep the group record: proceeding
// to delete the group would strand the filters with no way to reach them.
func TestDeleteLogGroupKeepsRecordWhenFilterListFails(t *testing.T) {
	s := newLogsTestStore(t)

	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	// Garbage bytes under a metric-filter key make ListMetricFilters fail
	// to decode the record.
	if err := s.PutRaw("metric-filter:group:broken", []byte{0xFF, 0xFF, 0xFF, 0xFF}); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroup("group"); err == nil {
		t.Fatal("delete reported success despite a failing filter list")
	}
	if _, err := s.GetLogGroup("group"); err != nil {
		t.Fatalf("group record removed despite undiscoverable filters: %v", err)
	}

	// Clearing the corrupt record lets the retry complete.
	if err := s.Delete("metric-filter:group:broken"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatalf("retry after clearing the corrupt record failed: %v", err)
	}
}

// A purge pass whose chunk removal fails keeps the stream's boundary
// members describing the surviving read surface: the chunk whose index
// record outlives the failed removal still serves reads, so zeroing the
// boundaries would claim an empty stream over readable events. The
// failure lever is the safe-path check — an index record whose chunk
// path escapes the chunks directory refuses to open, exactly as a
// tampered record would at read time.
func TestPurgeKeepsStreamBoundsWhenChunkRemovalFails(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "purge-fail-group", "purge-fail-stream"
	mustCreateGroupStream(t, s, group, stream)
	old := time.Now().UnixMilli() - 40*24*60*60*1000
	if _, err := s.PutLogEvents(group, stream, []LogEntry{{
		Timestamp: old, Message: "aged", IngestionTime: old,
	}}); err != nil {
		t.Fatal(err)
	}

	// Tamper the chunk record's path so the removal ladder fails at the
	// safe-path step while the index record stays in place.
	chunks, err := s.ListChunksForStream(group, stream)
	if err != nil || len(chunks) != 1 {
		t.Fatalf("chunk listing: %v (%d)", err, len(chunks))
	}
	tampered := *chunks[0]
	tampered.ChunkPath = "../../escape"
	if err := s.PutProto(s.chunkIndexKey(group, stream, tampered.ChunkID), ChunkMetaToProto(&tampered)); err != nil {
		t.Fatalf("tamper chunk record: %v", err)
	}

	if _, err := s.PurgeExpiredChunks(group, old+1000); err != nil {
		t.Fatalf("purge: %v", err)
	}
	ls, err := s.GetLogStream(group, stream)
	if err != nil {
		t.Fatal(err)
	}
	if ls.FirstEventTs == 0 && ls.LastEventTs == 0 {
		t.Fatal("stream boundaries must keep describing the chunk whose removal failed")
	}
	remaining, err := s.ListChunksForStream(group, stream)
	if err != nil || len(remaining) != 1 {
		t.Fatalf("the failed chunk's index record must survive for the retry: %v (%d)", err, len(remaining))
	}
}

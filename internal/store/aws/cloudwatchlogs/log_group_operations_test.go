package cloudwatchlogs

import (
	"testing"

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
		ChunkID:   "chunk-1",
		LogGroup:  "group",
		LogStream: "stream",
		ChunkPath: "missing-file",
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
		ChunkID:   "chunk-1",
		LogGroup:  "group",
		LogStream: "stream-a",
		ChunkPath: "missing-file",
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
		ChunkID:   "chunk-bad",
		LogGroup:  "group",
		LogStream: "stream-a",
		ChunkPath: "../escape",
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
	if err := s.PutRaw("metric-filter:group#broken", []byte{0xFF, 0xFF, 0xFF, 0xFF}); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroup("group"); err == nil {
		t.Fatal("delete reported success despite a failing filter list")
	}
	if _, err := s.GetLogGroup("group"); err != nil {
		t.Fatalf("group record removed despite undiscoverable filters: %v", err)
	}

	// Clearing the corrupt record lets the retry complete.
	if err := s.Delete("metric-filter:group#broken"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatalf("retry after clearing the corrupt record failed: %v", err)
	}
}

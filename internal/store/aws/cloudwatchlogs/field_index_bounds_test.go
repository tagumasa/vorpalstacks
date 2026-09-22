package cloudwatchlogs

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// newFieldIndexBoundsTestStore opens one store with a single log group.
func newFieldIndexBoundsTestStore(t *testing.T) *Store {
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
	if err := s.CreateLogGroup(&LogGroup{Name: "group"}); err != nil {
		t.Fatal(err)
	}
	return s
}

// The ingestion scan's merge accumulates per-field matching-event bounds
// across batches: a backdated event extends First, a non-matching event
// carries nothing, a dotted JSON path resolves, the scan stamp advances
// on every scan, and a zero event timestamp is a real bound (presence
// keyed, not a zero sentinel).
func TestMergeFieldIndexBoundsAccumulatesAcrossBatches(t *testing.T) {
	s := newFieldIndexBoundsTestStore(t)

	s.MergeFieldIndexBounds("group", []string{"RequestId", "userIdentity.accessKeyId", "traceId"}, []FieldScanMessage{
		{Timestamp: 2000, Message: `{"RequestId":"r2","userIdentity":{"accessKeyId":"k"}}`},
		{Timestamp: 3000, Message: "plain text"},
	})
	s.MergeFieldIndexBounds("group", []string{"RequestId", "Zero"}, []FieldScanMessage{
		{Timestamp: 1000, Message: `{"RequestId":"r1"}`},
		{Timestamp: 4000, Message: `{"Other":"x"}`},
		{Timestamp: 0, Message: `{"Zero":1}`},
	})
	s.MergeFieldIndexBounds("group", []string{"Zero"}, []FieldScanMessage{
		{Timestamp: 5000, Message: `{"Zero":1}`},
	})

	rec := s.ReadFieldIndexBounds("group")
	if r := rec.Fields["RequestId"]; r.First != 1000 || r.Last != 2000 {
		t.Fatalf("RequestId bounds: %+v, want 1000..2000", r)
	}
	if r := rec.Fields["userIdentity.accessKeyId"]; r.First != 2000 || r.Last != 2000 {
		t.Fatalf("dotted path bounds: %+v, want 2000..2000", r)
	}
	if r := rec.Fields["Zero"]; r.First != 0 || r.Last != 5000 {
		t.Fatalf("zero-timestamp bounds: %+v, want 0..5000", r)
	}
	if rec.LastScanTs == 0 {
		t.Fatal("the scan stamp never advanced")
	}
	if _, has := rec.Fields["traceId"]; has {
		t.Fatal("traceId matched without a message carrying it")
	}
}

// A batch that lands after the group's deletion must not resurrect the
// bounds record: the merge re-checks existence under the group lock.
func TestMergeFieldIndexBoundsSkipsDeletedGroup(t *testing.T) {
	s := newFieldIndexBoundsTestStore(t)
	s.MergeFieldIndexBounds("group", []string{"RequestId"}, []FieldScanMessage{
		{Timestamp: 1000, Message: `{"RequestId":"r"}`},
	})
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatal(err)
	}
	s.MergeFieldIndexBounds("group", []string{"RequestId"}, []FieldScanMessage{
		{Timestamp: 2000, Message: `{"RequestId":"r"}`},
	})
	if rec := s.ReadFieldIndexBounds("group"); len(rec.Fields) != 0 {
		t.Fatalf("a deleted group's merge resurrected the bounds record: %+v", rec.Fields)
	}
}

// The post-commit derived-data writers (the transformed-message copies
// and the field-index bounds) serialise with the group teardown through
// the group's write lock: neither can complete while the teardown holds
// the lock, and a transformed batch whose group vanished mid-flight
// writes nothing.
func TestDerivedWritesSerialiseWithTeardown(t *testing.T) {
	s := newFieldIndexBoundsTestStore(t)
	entries := map[string]TransformedMessage{
		TransformedMessageDigest(1000, "s1", "m"): {Timestamp: 1000, Message: "t"},
	}

	gLock := s.groupLock("group")
	gLock.Lock()
	boundsDone := make(chan struct{})
	transformedDone := make(chan struct{})
	go func() {
		defer close(boundsDone)
		s.MergeFieldIndexBounds("group", []string{"RequestId"}, []FieldScanMessage{
			{Timestamp: 1000, Message: `{"RequestId":"r"}`},
		})
	}()
	go func() {
		defer close(transformedDone)
		_ = s.PutTransformedMessages("group", entries)
	}()
	select {
	case <-boundsDone:
		t.Fatal("the bounds merge completed while the teardown held the group lock")
	case <-transformedDone:
		t.Fatal("the transformed write completed while the teardown held the group lock")
	case <-time.After(50 * time.Millisecond):
	}
	gLock.Unlock()
	<-boundsDone
	<-transformedDone
	if r := s.ReadFieldIndexBounds("group").Fields["RequestId"]; r.First != 1000 {
		t.Fatalf("bounds after the lock released: %+v", r)
	}
	m, err := s.TransformedMessageMap("group")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 {
		t.Fatalf("transformed records after the lock released: %d", len(m))
	}

	// A batch whose group vanished mid-flight writes nothing: the guard
	// under the group lock observes the removal and skips silently.
	if err := s.DeleteLogGroup("group"); err != nil {
		t.Fatal(err)
	}
	if err := s.PutTransformedMessages("group", entries); err != nil {
		t.Fatalf("a removed group must be a silent skip, got %v", err)
	}
	if m, err := s.TransformedMessageMap("group"); err != nil || len(m) != 0 {
		t.Fatalf("a removed group's transformed write resurrected records: %v %d", err, len(m))
	}
}

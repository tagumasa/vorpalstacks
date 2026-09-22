// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"errors"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// putEvents is the ingestion helper: a single PutLogEvents call whose
// messages total wantBytes, so StoredBytes accounting is exactly
// predictable.
func putEvents(t *testing.T, s *Store, group, stream string, msg string) {
	t.Helper()
	if _, err := s.PutLogEvents(group, stream, []LogEntry{{
		Timestamp: time.Now().UnixMilli(),
		Message:   msg,
	}}); err != nil {
		t.Fatal(err)
	}
}

// nonDirEntries names the regular files a directory holds, for
// assertions that a teardown unlinked the chunk files it owed.
func nonDirEntries(t *testing.T, dir string) []string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, de := range des {
		if !de.IsDir() {
			names = append(names, de.Name())
		}
	}
	return names
}

// TestPurgeExpiredChunksDecrementsExactlyOncePerChunk pins the purge-time
// half of the StoredBytes accounting: each chunk's decrement lands with
// its index record, so a purge retried after the files are gone is a
// no-op — never a second decrement and never a lost one.
func TestPurgeExpiredChunksDecrementsExactlyOncePerChunk(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "acct"}); err != nil {
		t.Fatal(err)
	}
	for _, stream := range []string{"a", "b"} {
		if err := s.CreateLogStream(&LogStream{LogGroupName: "acct", Name: stream}); err != nil {
			t.Fatal(err)
		}
	}
	msgA := strings.Repeat("x", 300)
	msgB := strings.Repeat("y", 170)
	putEvents(t, s, "acct", "a", msgA)
	putEvents(t, s, "acct", "b", msgB)

	// The chunk files disappear before the purge runs (the tolerated
	// os.IsNotExist path); the decrements must still land exactly once.
	for _, stream := range []string{"a", "b"} {
		chunks, err := s.ListChunksForStream("acct", stream)
		if err != nil {
			t.Fatal(err)
		}
		for _, chunk := range chunks {
			path, err := s.safeChunkPath(chunk.ChunkPath)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}
	}

	cutoff := time.Now().Add(time.Hour).UnixMilli()
	removed, err := s.PurgeExpiredChunks("acct", cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len(msgA) + len(msgB)); removed != want {
		t.Fatalf("first purge must remove the ingested byte total %d, got %d", want, removed)
	}
	lg, err := s.GetLogGroup("acct")
	if err != nil {
		t.Fatal(err)
	}
	if lg.StoredBytes != 0 {
		t.Fatalf("after the purge StoredBytes must be 0, got %d", lg.StoredBytes)
	}

	removed, err = s.PurgeExpiredChunks("acct", cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if removed != 0 {
		t.Fatalf("second purge must remove nothing, got %d", removed)
	}
	lg, err = s.GetLogGroup("acct")
	if err != nil {
		t.Fatal(err)
	}
	if lg.StoredBytes != 0 {
		t.Fatalf("after the second purge StoredBytes must still be 0, got %d", lg.StoredBytes)
	}
}

// TestDeleteLogStreamDecrementsStoredBytes pins the delete-time half of
// the StoredBytes accounting: the stream's chunks leave the group's
// counter on the same basis they entered it (ingested message bytes, not
// the compressed file size), and deleting every stream returns the
// counter to zero.
func TestDeleteLogStreamDecrementsStoredBytes(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "acct"}); err != nil {
		t.Fatal(err)
	}
	for _, stream := range []string{"a", "b"} {
		if err := s.CreateLogStream(&LogStream{LogGroupName: "acct", Name: stream}); err != nil {
			t.Fatal(err)
		}
	}

	msgA := strings.Repeat("x", 300)
	msgB := strings.Repeat("y", 170)
	putEvents(t, s, "acct", "a", msgA)
	putEvents(t, s, "acct", "b", msgB)

	lg, err := s.GetLogGroup("acct")
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len(msgA) + len(msgB)); lg.StoredBytes != want {
		t.Fatalf("after ingestion StoredBytes must be the message-byte total %d, got %d", want, lg.StoredBytes)
	}

	if err := s.DeleteLogStream("acct", "a"); err != nil {
		t.Fatal(err)
	}
	// The deletion removes the stream's chunk files with its records: a
	// removal that missed the file (a path resolution that re-rooted the
	// recorded form) would leak storage no sweep knows about.
	if left := nonDirEntries(t, s.chunksDir); len(left) != 1 {
		t.Fatalf("after deleting stream a %d chunk files remain, want stream b's alone: %v", len(left), left)
	}
	lg, err = s.GetLogGroup("acct")
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len(msgB)); lg.StoredBytes != want {
		t.Fatalf("after deleting stream a StoredBytes must be %d, got %d", want, lg.StoredBytes)
	}

	if err := s.DeleteLogStream("acct", "b"); err != nil {
		t.Fatal(err)
	}
	lg, err = s.GetLogGroup("acct")
	if err != nil {
		t.Fatal(err)
	}
	if lg.StoredBytes != 0 {
		t.Fatalf("after deleting every stream StoredBytes must be 0, got %d", lg.StoredBytes)
	}
}

// TestPurgeExpiredChunksDecrementsByIngestedBytes pins the purge path to
// the ingestion accounting basis: the decrement equals the purged
// chunks' recorded message bytes even though the compressed chunk files
// on disk are smaller.
func TestPurgeExpiredChunksDecrementsByIngestedBytes(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "purge", RetentionInDays: 1}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: "purge", Name: "s"}); err != nil {
		t.Fatal(err)
	}

	msg := strings.Repeat("compressible ", 100)
	putEvents(t, s, "purge", "s", msg)

	lg, err := s.GetLogGroup("purge")
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len(msg)); lg.StoredBytes != want {
		t.Fatalf("after ingestion StoredBytes must be the message-byte total %d, got %d", want, lg.StoredBytes)
	}

	// A cutoff beyond the chunk's timestamps expires it.
	if _, err := s.PurgeExpiredChunks("purge", time.Now().Add(time.Hour).UnixMilli()); err != nil {
		t.Fatal(err)
	}
	lg, err = s.GetLogGroup("purge")
	if err != nil {
		t.Fatal(err)
	}
	if lg.StoredBytes != 0 {
		t.Fatalf("the purge must decrement StoredBytes on the ingested-byte basis to 0, got %d", lg.StoredBytes)
	}
	chunks, err := s.ListChunksForStream("purge", "s")
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 0 {
		t.Fatalf("the purge must remove every chunk index, %d remain", len(chunks))
	}
}

// TestChunkScanFailurePropagates pins the error discipline of the chunk
// metadata listing: a backend scan failure must reach the caller — the
// event reads answer an error instead of a shortened page, and a stream
// teardown leaves the stream record in place so a retry can reach the
// chunks it could not list.
func TestChunkScanFailurePropagates(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "scan-fail-group", "scan-fail-stream"
	mustCreateGroupStream(t, s, group, stream)
	putEvents(t, s, group, stream, "payload")

	orig := scanPrefix
	scanPrefix = func(*Store, string, func(string, []byte) error) error {
		return errors.New("backend unavailable")
	}
	defer func() { scanPrefix = orig }()

	if _, _, _, err := s.GetLogEvents(group, stream, 0, 0, 0, true, ""); err == nil {
		t.Fatal("GetLogEvents must propagate the scan failure, not serve an empty page")
	}
	if _, _, err := s.FilterLogEvents(group, nil, 0, 0, "", 0, true, ""); err == nil {
		t.Fatal("FilterLogEvents must propagate the scan failure, not serve an empty page")
	}
	if err := s.DeleteLogStream(group, stream); err == nil {
		t.Fatal("DeleteLogStream must fail while the chunk list is unreadable")
	}
	if _, err := s.GetLogStream(group, stream); err != nil {
		t.Fatalf("the stream record must survive the failed teardown: %v", err)
	}
	if err := s.DeleteLogGroup(group); err == nil {
		t.Fatal("DeleteLogGroup must report the unreadable chunk backstop")
	}
	if _, err := s.GetLogGroup(group); err != nil {
		t.Fatalf("the group record must survive the failed teardown: %v", err)
	}

	scanPrefix = orig
	if err := s.DeleteLogGroup(group); err != nil {
		t.Fatalf("teardown must succeed once the scan answers again: %v", err)
	}
}

// TestDeleteLogGroupTearsDownDataProtectionPolicy pins the census's
// orphaned-family cell: the policy record keyed on the group name dies
// with the group, so a same-named recreation starts policy-free.
func TestDeleteLogGroupTearsDownDataProtectionPolicy(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "dpp"}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutDataProtectionPolicy(&DataProtectionPolicy{
		LogGroupIdentifier: "dpp",
		PolicyDocument:     `{"Name":"p"}`,
	}); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteLogGroup("dpp"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetDataProtectionPolicy("dpp"); err != ErrResourceNotFound {
		t.Fatalf("the policy must die with the group, got %v", err)
	}

	// A same-named recreation inherits nothing.
	if err := s.CreateLogGroup(&LogGroup{Name: "dpp"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetDataProtectionPolicy("dpp"); err != ErrResourceNotFound {
		t.Fatalf("a recreated group must start policy-free, got %v", err)
	}
}

// TestDeleteLogGroupTearsDownEveryFilter pins the filter families'
// teardown: no metric or subscription filter record survives the group
// delete (the metric listing walks pages, so the loop is exercised on
// multi-page-safe ground even below the page size).
func TestDeleteLogGroupTearsDownEveryFilter(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(&LogGroup{Name: "filters"}); err != nil {
		t.Fatal(err)
	}
	transform := []MetricTransformation{{MetricName: "Count", MetricNamespace: "N", MetricValue: "1"}}
	for i := 0; i < 3; i++ {
		name := string(rune('a' + i))
		if err := s.PutMetricFilter(NewMetricFilter("filters", name, "ERROR", transform)); err != nil {
			t.Fatal(err)
		}
		if err := s.PutSubscriptionFilter(&SubscriptionFilter{
			LogGroupName:   "filters",
			FilterName:     name,
			DestinationArn: s.arnBuilder.CloudWatch().LogGroup("filters"),
		}); err != nil {
			t.Fatal(err)
		}
	}

	if err := s.DeleteLogGroup("filters"); err != nil {
		t.Fatal(err)
	}
	if mf, _, err := s.ListMetricFilters("filters", "", "", ListingPageSize); err != nil {
		t.Fatal(err)
	} else if len(mf) != 0 {
		t.Fatalf("no metric filter may survive the group delete, %d remain", len(mf))
	}
	if sf, err := s.ListSubscriptionFilters("filters", ""); err != nil {
		t.Fatal(err)
	} else if len(sf) != 0 {
		t.Fatalf("no subscription filter may survive the group delete, %d remain", len(sf))
	}
}

// failingBucket forwards to the wrapped bucket except for keys carrying
// the fail prefix, whose puts and/or deletes fail per the flags — the
// injection point for the sequential-write paths' partial-failure
// discipline. An armed gate lets a fixture build its records over the
// same bucket before the injection starts.
type failingBucket struct {
	storage.Bucket
	failPrefix  string
	failPuts    bool
	failDeletes bool
	armed       *atomic.Bool
}

func (b failingBucket) Put(key, value []byte) error {
	if b.failPuts && strings.HasPrefix(string(key), b.failPrefix) &&
		(b.armed == nil || b.armed.Load()) {
		return errors.New("injected write failure")
	}
	return b.Bucket.Put(key, value)
}

func (b failingBucket) Delete(key []byte) error {
	if b.failDeletes && strings.HasPrefix(string(key), b.failPrefix) &&
		(b.armed == nil || b.armed.Load()) {
		return errors.New("injected delete failure")
	}
	return b.Bucket.Delete(key)
}

// nonTransactional hides the storage's transaction methods so a store
// built over it exercises the sequential-write fallback paths.
type nonTransactional struct {
	storage.BasicStorage
}

// newFallbackStore builds a fallback-mode store (no transactional
// commits) whose bucket fails the flagged operations on keys carrying
// failPrefix. The returned arm function starts the injections, so the
// fixture can be built first; callers that inject from the start may
// ignore it.
func newFallbackStore(t *testing.T, failPrefix string, failPuts, failDeletes bool) (*Store, func()) {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	var armed atomic.Bool
	s, err := NewStore(nonTransactional{st}, failingBucket{
		Bucket:      st.Bucket("logs-us-east-1"),
		failPrefix:  failPrefix,
		failPuts:    failPuts,
		failDeletes: failDeletes,
		armed:       &armed,
	}, "000000000000", "us-east-1", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s, func() { armed.Store(true) }
}

// TestDeleteStreamKeepsBytesWhenIndexDeleteFails pins the teardown's
// accounting discipline against the chunk-index delete failing: the
// record survives for the retry, so its bytes must stay counted —
// decrementing here would double-decrement when the retry reaches the
// record again.
func TestDeleteStreamKeepsBytesWhenIndexDeleteFails(t *testing.T) {
	s, arm := newFallbackStore(t, keyPrefixChunk, false, true)
	// The fixture performs no chunk-index deletes, so the injection can
	// start immediately.
	arm()
	const group, stream = "acct", "s"
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: stream}); err != nil {
		t.Fatal(err)
	}
	const msg = "payload"
	putEvents(t, s, group, stream, msg)

	lg, err := s.GetLogGroup(group)
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len(msg)); lg.StoredBytes != want {
		t.Fatalf("after ingestion StoredBytes must be %d, got %d", want, lg.StoredBytes)
	}

	if err := s.DeleteLogStream(group, stream); err != nil {
		t.Fatal(err)
	}

	lg, err = s.GetLogGroup(group)
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len(msg)); lg.StoredBytes != want {
		t.Fatalf("with the index delete failing, StoredBytes must stay %d for the retry, got %d",
			want, lg.StoredBytes)
	}
	chunks, err := s.ListChunksForStream(group, stream)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 1 {
		t.Fatalf("the chunk index record must survive the failed delete for the retry, %d listed", len(chunks))
	}
}

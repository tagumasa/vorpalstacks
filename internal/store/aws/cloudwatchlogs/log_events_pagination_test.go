package cloudwatchlogs

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// putTestEvents writes one batch of events with sequential timestamps
// into the given stream.
func putTestEvents(t *testing.T, s *Store, group, stream string, base int64, messages []string) {
	t.Helper()
	entries := make([]LogEntry, len(messages))
	for i, m := range messages {
		entries[i] = LogEntry{Timestamp: base + int64(i), Message: m, IngestionTime: base + int64(i)}
	}
	if _, err := s.PutLogEvents(group, stream, entries); err != nil {
		t.Fatalf("PutLogEvents: %v", err)
	}
}

// GetLogEvents pages forward through scoped tokens and terminates by
// re-offering the presented token (the model's documented end-of-stream
// rule); both tokens are always present.
func TestGetLogEventsScopedTokensPageAndTerminate(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "page-group", "page-stream"
	mustCreateGroupStream(t, s, group, stream)

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, stream, base, []string{"e1", "e2", "e3", "e4", "e5"})

	var served []string
	token := ""
	pages := 0
	for {
		pages++
		if pages > 10 {
			t.Fatal("pagination did not terminate")
		}
		events, fwd, bwd, err := s.GetLogEvents(group, stream, 0, 0, 2, true, token)
		if err != nil {
			t.Fatalf("GetLogEvents page %d: %v", pages, err)
		}
		if fwd == "" || bwd == "" {
			t.Fatalf("page %d: tokens must never be null (fwd=%q bwd=%q)", pages, fwd, bwd)
		}
		for _, e := range events {
			served = append(served, e.Message)
		}
		if token != "" && fwd == token {
			// End of stream: the same token came back.
			break
		}
		if len(events) == 0 {
			t.Fatal("empty page before token repetition")
		}
		token = fwd
	}

	want := []string{"e1", "e2", "e3", "e4", "e5"}
	if len(served) != len(want) {
		t.Fatalf("served %v, want %v", served, want)
	}
	for i := range want {
		if served[i] != want[i] {
			t.Fatalf("served[%d] = %q, want %q (full: %v)", i, served[i], want[i], served)
		}
	}
}

// Byte-identical duplicate events are distinct at a page boundary: the
// documented ordering key ends with the ID of the PutLogEvents request,
// and the cursor's chunk-and-position ordinal stands in for it, so a
// page that stops on the first copy serves the second copy next instead
// of dropping it.
func TestGetLogEventsServesDuplicateEventsAcrossPages(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "dup-group", "dup-stream"
	mustCreateGroupStream(t, s, group, stream)

	ts := time.Now().UnixMilli()
	if _, err := s.PutLogEvents(group, stream, []LogEntry{
		{Timestamp: ts, Message: "same", IngestionTime: ts},
		{Timestamp: ts, Message: "same", IngestionTime: ts},
		{Timestamp: ts, Message: "tail", IngestionTime: ts},
	}); err != nil {
		t.Fatalf("PutLogEvents: %v", err)
	}

	all, _, _, err := s.GetLogEvents(group, stream, 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("unpaged read must serve every ingested copy, got %d of 3", len(all))
	}

	var served []string
	token := ""
	for pages := 0; pages < 6; pages++ {
		events, fwd, _, err := s.GetLogEvents(group, stream, 0, 0, 1, true, token)
		if err != nil {
			t.Fatalf("page %d: %v", pages+1, err)
		}
		for _, e := range events {
			served = append(served, e.Message)
		}
		if token != "" && fwd == token {
			break
		}
		if len(events) == 0 {
			t.Fatal("empty page before token repetition")
		}
		token = fwd
	}
	want := []string{"same", "same", "tail"}
	if len(served) != len(want) {
		t.Fatalf("served %v, want %v", served, want)
	}
	for i := range want {
		if served[i] != want[i] {
			t.Fatalf("served[%d] = %q, want %q (full: %v)", i, served[i], want[i], served)
		}
	}
}

// A GetLogEvents token is scoped to its request: replaying it against a
// different stream or window is rejected, not silently repositioned.
func TestGetLogEventsTokenScopeIsRequestIdentity(t *testing.T) {
	s := newLogsTestStore(t)
	const group = "scope-group"
	mustCreateGroupStream(t, s, group, "stream-a")
	mustCreateGroupStream(t, s, group, "stream-b")

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, "stream-a", base, []string{"a1", "a2"})
	putTestEvents(t, s, group, "stream-b", base, []string{"b1", "b2"})

	_, fwd, _, err := s.GetLogEvents(group, "stream-a", 0, 0, 1, true, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, _, _, err := s.GetLogEvents(group, "stream-b", 0, 0, 1, true, fwd); !errors.Is(err, ErrInvalidPaginationToken) {
		t.Fatalf("token from stream-a replayed against stream-b: got %v, want ErrInvalidPaginationToken", err)
	}
	if _, _, _, err := s.GetLogEvents(group, "stream-a", base+100, base+200, 1, true, fwd); !errors.Is(err, ErrInvalidPaginationToken) {
		t.Fatalf("token replayed against a different window: got %v, want ErrInvalidPaginationToken", err)
	}
	if _, _, _, err := s.GetLogEvents(group, "stream-a", 0, 0, 1, true, "not-a-token"); !errors.Is(err, ErrInvalidPaginationToken) {
		t.Fatalf("garbage token: got %v, want ErrInvalidPaginationToken", err)
	}
}

// A later write delivering older timestamps must not disturb the pages:
// the cursor resumes strictly after the last served event, so page two
// serves exactly the events that followed page one.
func TestGetLogEventsBackfillStableAcrossPages(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "backfill-group", "backfill-stream"
	mustCreateGroupStream(t, s, group, stream)

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, stream, base, []string{"m1", "m2", "m3", "m4", "m5"})

	page1, fwd, _, err := s.GetLogEvents(group, stream, 0, 0, 2, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 2 || page1[0].Message != "m1" || page1[1].Message != "m2" {
		t.Fatalf("page 1 = %v", page1)
	}

	// Backfill: a legal second batch whose timestamps precede page one's.
	putTestEvents(t, s, group, stream, base-1000, []string{"old1", "old2"})

	page2, fwd2, _, err := s.GetLogEvents(group, stream, 0, 0, 2, true, fwd)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 || page2[0].Message != "m3" || page2[1].Message != "m4" {
		t.Fatalf("page 2 after backfill = %v, want [m3 m4]", page2)
	}
	page3, fwd3, _, err := s.GetLogEvents(group, stream, 0, 0, 2, true, fwd2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page3) != 1 || page3[0].Message != "m5" {
		t.Fatalf("page 3 = %v, want [m5]", page3)
	}
	// The backfilled events sort before the page-one cursor: a fresh walk
	// from the head serves them first, but the paged walk never rewinds.
	fresh, _, _, err := s.GetLogEvents(group, stream, 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(fresh) != 7 || fresh[0].Message != "old1" || fresh[1].Message != "old2" || fresh[2].Message != "m1" {
		t.Fatalf("fresh read = %v", fresh)
	}
	_ = fwd3
}

// Backward reads (the GetLogEvents default direction) serve newest first
// and page toward the head through the same vocabulary.
func TestGetLogEventsBackwardPagesTowardHead(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "backward-group", "backward-stream"
	mustCreateGroupStream(t, s, group, stream)

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, stream, base, []string{"n1", "n2", "n3", "n4", "n5"})

	page1, _, bwd, err := s.GetLogEvents(group, stream, 0, 0, 2, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 2 || page1[0].Message != "n5" || page1[1].Message != "n4" {
		t.Fatalf("backward page 1 = %v, want [n5 n4]", page1)
	}

	page2, fwd2, bwd2, err := s.GetLogEvents(group, stream, 0, 0, 2, false, bwd)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 || page2[0].Message != "n3" || page2[1].Message != "n2" {
		t.Fatalf("backward page 2 = %v, want [n3 n2]", page2)
	}

	// The forward token of a backward page continues after the newest
	// event that page covered, letting a client flip direction mid-walk:
	// page 2 covered [n2, n3], so its forward token resumes at n4.
	resumed, _, _, err := s.GetLogEvents(group, stream, 0, 0, 2, true, fwd2)
	if err != nil {
		t.Fatal(err)
	}
	if len(resumed) != 2 || resumed[0].Message != "n4" || resumed[1].Message != "n5" {
		t.Fatalf("forward resume from backward page 2 = %v, want [n4 n5]", resumed)
	}
	_ = bwd2
}

// A resumed walk gathers only the unserved side of the cursor — chunks
// whose bounds sit entirely on the served side are skipped instead of
// being re-read and re-sorted on every page — while chunks whose bounds
// merely touch the cursor still serve their same-timestamp events
// (those order by tie-break members the chunk bounds cannot see). Both
// directions page the full scope exactly once across the chunk eras.
func TestGetLogEventsResumedWalkServesFullScopeAcrossChunkEras(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "narrow-group", "narrow-stream"
	mustCreateGroupStream(t, s, group, stream)

	// Three chunks (each put persists its own): an ascending era, a
	// same-millisecond trio whose internal order is the digest
	// tie-break, and a closing era.
	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, stream, base, []string{"a1", "a2", "a3", "a4"})
	sameMs := base + 4
	sameEntries := make([]LogEntry, 3)
	for i, m := range []string{"b1", "b2", "b3"} {
		sameEntries[i] = LogEntry{Timestamp: sameMs, Message: m, IngestionTime: sameMs}
	}
	if _, err := s.PutLogEvents(group, stream, sameEntries); err != nil {
		t.Fatalf("PutLogEvents same-millisecond trio: %v", err)
	}
	putTestEvents(t, s, group, stream, base+5, []string{"c1", "c2", "c3"})

	want := []string{"a1", "a2", "a3", "a4", "b1", "b2", "b3", "c1", "c2", "c3"}

	// Forward: page size two keeps a page boundary inside every era,
	// including the same-millisecond trio's.
	served := walkAllPages(t, s, group, stream, 2, true)
	assertScopeComplete(t, served, want)
	if aFirst, aLast := eraBounds(served, "a"); aFirst != 0 || aLast != 3 {
		t.Fatalf("forward walk's first era misplaced: [%d,%d]", aFirst, aLast)
	}
	if cFirst, _ := eraBounds(served, "c"); cFirst != 7 {
		t.Fatalf("forward walk's last era starts at %d, want 7", cFirst)
	}

	// Backward: the same scope newest first.
	served = walkAllPages(t, s, group, stream, 2, false)
	assertScopeComplete(t, served, want)
	if cFirst, cLast := eraBounds(served, "c"); cFirst != 0 || cLast != 2 {
		t.Fatalf("backward walk's first era misplaced: [%d,%d]", cFirst, cLast)
	}
	if aFirst, _ := eraBounds(served, "a"); aFirst != 6 {
		t.Fatalf("backward walk's last era starts at %d, want 6", aFirst)
	}

	// The unpaged read agrees on the complete scope.
	fresh, _, _, err := s.GetLogEvents(group, stream, 0, 0, 0, true, "")
	if err != nil {
		t.Fatal(err)
	}
	assertScopeComplete(t, fresh2messages(fresh), want)
}

// walkAllPages pages one direction to the documented termination (the
// re-offered token) and returns the served messages in page order.
func walkAllPages(t *testing.T, s *Store, group, stream string, limit int, startFromHead bool) []string {
	t.Helper()
	var served []string
	token := ""
	for pages := 0; pages < 12; pages++ {
		events, fwd, bwd, err := s.GetLogEvents(group, stream, 0, 0, limit, startFromHead, token)
		if err != nil {
			t.Fatalf("page %d: %v", pages+1, err)
		}
		for _, e := range events {
			served = append(served, e.Message)
		}
		next := fwd
		if !startFromHead {
			next = bwd
		}
		if token != "" && next == token {
			return served
		}
		if len(events) == 0 {
			t.Fatal("empty page before token repetition")
		}
		token = next
	}
	t.Fatal("pagination did not terminate")
	return nil
}

// assertScopeComplete checks the served messages are exactly the wanted
// multiset: no duplication, no loss (the intra-era order of a
// same-millisecond group is the digest tie-break, which the assertion
// deliberately does not pin).
func assertScopeComplete(t *testing.T, served, want []string) {
	t.Helper()
	if len(served) != len(want) {
		t.Fatalf("served %d events, want %d: %v", len(served), len(want), served)
	}
	counts := make(map[string]int, len(want))
	for _, m := range want {
		counts[m]++
	}
	for _, m := range served {
		counts[m]--
		if counts[m] < 0 {
			t.Fatalf("served %q more times than it exists: %v", m, served)
		}
	}
}

// eraBounds returns the first and last index of the messages carrying
// the era prefix.
func eraBounds(served []string, era string) (first, last int) {
	first, last = -1, -1
	for i, m := range served {
		if strings.HasPrefix(m, era) {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	return first, last
}

// fresh2messages projects output events onto their messages.
func fresh2messages(events []*OutputLogEvent) []string {
	out := make([]string, len(events))
	for i, e := range events {
		out[i] = e.Message
	}
	return out
}

// FilterLogEvents omits the next token when its direction is exhausted
// (the model documents the absent token as the end of pagination).
func TestFilterLogEventsPaginationOmitsTokenAtEnd(t *testing.T) {
	s := newLogsTestStore(t)
	const group = "filter-page-group"
	mustCreateGroupStream(t, s, group, "fs1")
	mustCreateGroupStream(t, s, group, "fs2")

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, "fs1", base, []string{"f1", "f3"})
	putTestEvents(t, s, group, "fs2", base+1, []string{"f2", "f4"})

	var served []string
	token := ""
	pages := 0
	for {
		pages++
		if pages > 10 {
			t.Fatal("pagination did not terminate")
		}
		events, next, err := s.FilterLogEvents(group, nil, 0, 0, "", 2, true, token)
		if err != nil {
			t.Fatalf("FilterLogEvents page %d: %v", pages, err)
		}
		for _, e := range events {
			served = append(served, e.Message)
		}
		if next == "" {
			break
		}
		token = next
	}
	// fs1's events sit at base/base+1 and fs2's at base+1/base+2: the
	// interleaved order is by timestamp, so f3 (fs1, base+1) precedes f2
	// (fs2, base+1) on the stream tiebreak.
	wantOrder := []string{"f1", "f3", "f2", "f4"}
	if len(served) != len(wantOrder) {
		t.Fatalf("served %v, want %v", served, wantOrder)
	}
	for i := range wantOrder {
		if served[i] != wantOrder[i] {
			t.Fatalf("served[%d] = %q, want %q (full: %v)", i, served[i], wantOrder[i], served)
		}
	}
}

// FilterLogEvents tokens carry the pattern and stream set in their scope.
func TestFilterLogEventsTokenScopeIncludesPatternAndStreams(t *testing.T) {
	s := newLogsTestStore(t)
	const group = "filter-scope-group"
	mustCreateGroupStream(t, s, group, "s1")
	mustCreateGroupStream(t, s, group, "s2")

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, "s1", base, []string{"ERROR one", "ERROR one-b", "ok two"})
	putTestEvents(t, s, group, "s2", base, []string{"ERROR three", "ok four"})

	_, tok, err := s.FilterLogEvents(group, []string{"s1"}, 0, 0, "ERROR", 1, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if tok == "" {
		t.Fatal("first page exhausted the match set; no token to replay")
	}
	if _, _, err := s.FilterLogEvents(group, []string{"s1"}, 0, 0, "", 1, true, tok); !errors.Is(err, ErrInvalidPaginationToken) {
		t.Fatalf("token replayed with a different pattern: got %v, want ErrInvalidPaginationToken", err)
	}
	if _, _, err := s.FilterLogEvents(group, []string{"s2"}, 0, 0, "ERROR", 1, true, tok); !errors.Is(err, ErrInvalidPaginationToken) {
		t.Fatalf("token replayed against a different stream set: got %v, want ErrInvalidPaginationToken", err)
	}
	if _, _, err := s.FilterLogEvents(group, nil, 0, 0, "ERROR", 1, true, tok); !errors.Is(err, ErrInvalidPaginationToken) {
		t.Fatalf("token replayed against the whole group: got %v, want ErrInvalidPaginationToken", err)
	}
}

// ListLogStreams takes a stream NAME marker and resumes strictly after
// it; callers never handle raw store keys.
func TestListLogStreamsNameMarkerResumesAfterBoundary(t *testing.T) {
	s := newLogsTestStore(t)
	const group = "marker-group"
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a", "b", "c"} {
		if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: name}); err != nil {
			t.Fatal(err)
		}
	}

	names := func(streams []*LogStream) []string {
		out := make([]string, len(streams))
		for i, ls := range streams {
			out[i] = ls.Name
		}
		return out
	}

	if got := names(mustListStreams(t, s, group, "", 10)); len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("no marker: %v", got)
	}
	if got := names(mustListStreams(t, s, group, "a", 10)); len(got) != 2 || got[0] != "b" || got[1] != "c" {
		t.Fatalf("marker a: %v", got)
	}
	if got := names(mustListStreams(t, s, group, "b", 10)); len(got) != 1 || got[0] != "c" {
		t.Fatalf("marker b: %v", got)
	}
	if got := mustListStreams(t, s, group, "c", 10); len(got) != 0 {
		t.Fatalf("marker c: %v", got)
	}
}

func mustListStreams(t *testing.T, s *Store, group, marker string, limit int) []*LogStream {
	t.Helper()
	streams, _, err := s.ListLogStreams(group, "", marker, limit)
	if err != nil {
		t.Fatalf("ListLogStreams: %v", err)
	}
	return streams
}

func mustCreateGroupStream(t *testing.T, s *Store, group, stream string) {
	t.Helper()
	if err := s.CreateLogGroup(&LogGroup{Name: group}); err != nil && !errors.Is(err, ErrLogGroupAlreadyExists) {
		t.Fatalf("CreateLogGroup: %v", err)
	}
	if err := s.CreateLogStream(&LogStream{LogGroupName: group, Name: stream}); err != nil {
		t.Fatalf("CreateLogStream: %v", err)
	}
}

// The two event reads carry opposite endTime boundary semantics: GetLogEvents
// excludes the event stamped exactly at endTime, FilterLogEvents keeps it.
func TestGetLogEventsEndTimeExcludesBoundaryEvent(t *testing.T) {
	s := newLogsTestStore(t)
	const group, stream = "end-boundary-group", "end-boundary-stream"
	mustCreateGroupStream(t, s, group, stream)

	base := time.Now().UnixMilli()
	putTestEvents(t, s, group, stream, base, []string{"before", "boundary", "after"})

	events, _, _, err := s.GetLogEvents(group, stream, 0, base+1, 10, true, "")
	if err != nil {
		t.Fatalf("GetLogEvents: %v", err)
	}
	if len(events) != 1 || events[0].Message != "before" {
		t.Fatalf("GetLogEvents endTime must exclude the boundary event, got %d events", len(events))
	}

	filtered, _, err := s.FilterLogEvents(group, nil, 0, base+1, "", 10, true, "")
	if err != nil {
		t.Fatalf("FilterLogEvents: %v", err)
	}
	if len(filtered) != 2 {
		t.Fatalf("FilterLogEvents endTime must keep the boundary event, got %d events", len(filtered))
	}
}

// The logStreamNames member is a set, not an ordered sequence: a repeated
// name must not list its stream's chunks twice — every event of the
// stream returns once, exactly as the un-repeated request serves it.
func TestFilterLogEventsRepeatedStreamNamesServeEachEventOnce(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(NewLogGroup("dup-names-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := s.CreateLogStream(&LogStream{Name: "s", LogGroupName: "dup-names-group"}); err != nil {
		t.Fatalf("stream: %v", err)
	}
	now := time.Now().UnixMilli()
	entries := make([]LogEntry, 3)
	for i := range entries {
		entries[i] = LogEntry{Timestamp: now + int64(i), Message: "m", IngestionTime: now}
	}
	if _, err := s.PutLogEvents("dup-names-group", "s", entries); err != nil {
		t.Fatalf("seed: %v", err)
	}
	once, _, err := s.FilterLogEvents("dup-names-group", []string{"s"}, 0, 0, "", 0, true, "")
	if err != nil {
		t.Fatalf("single name: %v", err)
	}
	repeated, _, err := s.FilterLogEvents("dup-names-group", []string{"s", "s"}, 0, 0, "", 0, true, "")
	if err != nil {
		t.Fatalf("repeated name: %v", err)
	}
	if len(repeated) != len(once) {
		t.Fatalf("repeated stream name: want %d events (the single-name count), got %d", len(once), len(repeated))
	}
}

// The late window the delivery engine reads below a cursor: an event put
// after the pass that moved the cursor past its timestamp is findable by
// (timestamp below the cursor, ingestion time above the mark), and the
// chunk-index ingestion bound keeps streams without late arrivals at one
// meta scan.
func TestLateIngestionEventsBelowCursor(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.CreateLogGroup(NewLogGroup("late-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("group: %v", err)
	}
	if err := s.CreateLogStream(&LogStream{Name: "s", LogGroupName: "late-group"}); err != nil {
		t.Fatalf("stream: %v", err)
	}
	mark := time.Now().UnixMilli()
	// The original window the cursor will sit past.
	if _, err := s.PutLogEvents("late-group", "s", []LogEntry{
		{Timestamp: mark - 10_000, Message: "old", IngestionTime: mark - 10_000},
		{Timestamp: mark - 9_000, Message: "newest", IngestionTime: mark - 9_000},
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	// A backdated put landing below the cursor after the mark.
	if _, err := s.PutLogEvents("late-group", "s", []LogEntry{
		{Timestamp: mark - 50_000, Message: "backdated", IngestionTime: mark + 5},
	}); err != nil {
		t.Fatalf("backdated put: %v", err)
	}
	late, err := s.LateIngestionEvents("late-group", "s", mark, mark-9_000)
	if err != nil {
		t.Fatalf("late window: %v", err)
	}
	if len(late) != 1 || late[0].Message != "backdated" {
		t.Fatalf("late window: want the single backdated event, got %v", late)
	}
	// Everything at or past the cursor, or ingested at or before the
	// mark, stays outside the window; a mark at the newest ingestion
	// time empties it.
	late, err = s.LateIngestionEvents("late-group", "s", mark+5, mark-9_000)
	if err != nil {
		t.Fatalf("advanced mark: %v", err)
	}
	if len(late) != 0 {
		t.Fatalf("advanced mark: want an empty window, got %v", late)
	}
}

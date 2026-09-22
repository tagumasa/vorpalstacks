package cloudwatchlogs

import (
	"context"
	"errors"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// newReadTestService builds a service backed by temporary storage with
// its group created, for the read-layer tests.
func newReadTestService(t *testing.T, group string) (*LogsService, *logsstore.Store) {
	t.Helper()
	svc, store := newTestService(t)
	createTestLogGroup(t, store, group)
	return svc, store
}

func putReadEvents(t *testing.T, store *logsstore.Store, group, stream string, base int64, n int) {
	t.Helper()
	remaining := n
	ts := base
	for remaining > 0 {
		batch := remaining
		if batch > 5000 {
			batch = 5000
		}
		entries := make([]logsstore.LogEntry, batch)
		for i := range entries {
			entries[i] = logsstore.LogEntry{Timestamp: ts, Message: "m", IngestionTime: ts}
			ts++
		}
		if _, err := store.PutLogEvents(group, stream, entries); err != nil {
			t.Fatalf("PutLogEvents: %v", err)
		}
		remaining -= batch
	}
}

// The fetch-all layer pages past the wire limit: a stream holding more
// events than one GetLogEvents page serves is read to exhaustion.
func TestFetchAllLogEventsPagesPastWireLimit(t *testing.T) {
	const group, stream = "fetchall-group", "fetchall-stream"
	_, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}

	base := time.Now().UnixMilli()
	putReadEvents(t, store, group, stream, base, logsstore.MaxEventsLimit+5)

	events, err := fetchAllLogEvents(store, group, stream, 0, 0)
	if err != nil {
		t.Fatalf("fetchAllLogEvents: %v", err)
	}
	if len(events) != logsstore.MaxEventsLimit+5 {
		t.Fatalf("fetched %d events, want %d", len(events), logsstore.MaxEventsLimit+5)
	}
	for i := 1; i < len(events); i++ {
		if events[i-1].Timestamp > events[i].Timestamp {
			t.Fatalf("events not in deterministic order at %d", i)
		}
	}
}

// The bounded form of the fetch layer stops at its explicit sample cap.
func TestFetchEventsUpToBoundsSample(t *testing.T) {
	const group, stream = "sample-group", "sample-stream"
	_, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}

	putReadEvents(t, store, group, stream, time.Now().UnixMilli(), 100)

	capped, err := fetchEventsUpTo(store, group, stream, 0, 0, 20)
	if err != nil {
		t.Fatalf("capped fetch: %v", err)
	}
	if len(capped) != 20 {
		t.Fatalf("capped fetch returned %d events, want 20", len(capped))
	}
	all, err := fetchEventsUpTo(store, group, stream, 0, 0, 0)
	if err != nil {
		t.Fatalf("unbounded fetch: %v", err)
	}
	if len(all) != 100 {
		t.Fatalf("unbounded fetch returned %d events, want 100", len(all))
	}
}

// The fetch layer's two end-boundary forms: the inclusive read keeps the
// event whose timestamp equals the window end (the documented StartQuery
// and CreateExportTask windows), while the exclusive default keeps the
// GetLogEvents boundary.
func TestFetchAllLogEventsEndBoundaryForms(t *testing.T) {
	const group, stream = "endbound-group", "endbound-stream"
	_, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}

	base := time.Now().UnixMilli()
	putReadEvents(t, store, group, stream, base, 3)

	exclusive, err := fetchAllLogEvents(store, group, stream, 0, base+1)
	if err != nil {
		t.Fatalf("exclusive fetch: %v", err)
	}
	if len(exclusive) != 1 || exclusive[0].Timestamp != base {
		t.Fatalf("exclusive window returned %d events, want the pre-boundary one alone", len(exclusive))
	}
	inclusive, err := fetchAllLogEventsEndInclusive(store, group, stream, 0, base+1)
	if err != nil {
		t.Fatalf("inclusive fetch: %v", err)
	}
	if len(inclusive) != 2 || inclusive[1].Timestamp != base+1 {
		t.Fatalf("inclusive window returned %d events, want the boundary event included", len(inclusive))
	}
}

// DescribeLogStreams speaks one token vocabulary across every ordering:
// ascending pages, event-time ordering and descending pages all mint
// scoped tokens that are neither raw stream names nor positional
// offsets, and a token replayed against a different request shape is
// rejected.
func TestDescribeLogStreamsOneTokenVocabularyAcrossOrderings(t *testing.T) {
	const group = "vocab-group"
	svc, store := newReadTestService(t, group)

	base := time.Now().UnixMilli()
	for i, name := range []string{"vs-a", "vs-b", "vs-c", "vs-d", "vs-e"} {
		if err := store.CreateLogStream(logsstore.NewLogStream(name, group)); err != nil {
			t.Fatal(err)
		}
		putReadEvents(t, store, group, name, base+int64(i)*1000, 1)
	}

	// Ascending name order: pages resume strictly after the cursor.
	var ascending []string
	token := ""
	for {
		res, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{
			LogGroupName: group, NextToken: token, Limit: 2, Region: "us-east-1",
		})
		if err != nil {
			t.Fatalf("ascending page: %v", err)
		}
		for _, ls := range res.LogStreams {
			ascending = append(ascending, ls.Name)
		}
		if res.NextToken == "" {
			break
		}
		if res.NextToken == "vs-b" || res.NextToken == "vs-d" {
			t.Fatalf("next token leaks the raw stream name: %q", res.NextToken)
		}
		token = res.NextToken
	}
	wantAsc := []string{"vs-a", "vs-b", "vs-c", "vs-d", "vs-e"}
	for i := range wantAsc {
		if ascending[i] != wantAsc[i] {
			t.Fatalf("ascending walk = %v, want %v", ascending, wantAsc)
		}
	}

	// LastEventTime ordering pages by its own (event time, name) cursor,
	// newest first.
	var byEventTime []string
	token = ""
	for {
		res, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{
			LogGroupName: group, OrderBy: "LastEventTime", Descending: true,
			NextToken: token, Limit: 2, Region: "us-east-1",
		})
		if err != nil {
			t.Fatalf("event-time page: %v", err)
		}
		for _, ls := range res.LogStreams {
			byEventTime = append(byEventTime, ls.Name)
		}
		if res.NextToken == "" {
			break
		}
		token = res.NextToken
	}
	wantEvent := []string{"vs-e", "vs-d", "vs-c", "vs-b", "vs-a"}
	for i := range wantEvent {
		if byEventTime[i] != wantEvent[i] {
			t.Fatalf("event-time walk = %v, want %v", byEventTime, wantEvent)
		}
	}

	// Descending name order through the same vocabulary.
	res, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{
		LogGroupName: group, Descending: true, Limit: 2, Region: "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.LogStreams) != 2 || res.LogStreams[0].Name != "vs-e" || res.LogStreams[1].Name != "vs-d" {
		t.Fatalf("descending first page = %v", res.LogStreams)
	}

	// A token minted for one request shape does not serve another.
	res2, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{
		LogGroupName: group, Descending: true, Limit: 2, Region: "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{
		LogGroupName: group, OrderBy: "LastEventTime", NextToken: res2.NextToken,
		Limit: 2, Region: "us-east-1",
	}); err == nil {
		t.Fatal("descending-name token accepted by the event-time ordering")
	}
	if _, err := svc.describeLogStreamsCore(DescribeLogStreamsInput{
		LogGroupName: group, NextToken: "garbage-token", Limit: 2, Region: "us-east-1",
	}); err == nil {
		t.Fatal("garbage token accepted")
	}
}

// DescribeQueries lists deterministically (newest first) and pages
// through a scoped cursor; a token replayed against a different filter
// is rejected.
func TestDescribeQueriesDeterministicOrderAndScopedToken(t *testing.T) {
	svc, _ := newReadTestService(t, "queries-group")

	now := time.Now()
	for i := 4; i >= 0; i-- {
		qs := &queryState{
			queryId:       fmtQueryId(i),
			status:        "Complete",
			createdAt:     now.Add(-time.Duration(i) * time.Minute),
			logGroupNames: []string{"queries-group"},
		}
		svc.queries.Store(qs.queryId, qs)
	}

	page1, next, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 2 || page1[0].queryId != fmtQueryId(0) || page1[1].queryId != fmtQueryId(1) {
		t.Fatalf("page 1 = %s,%s want the two newest", page1[0].queryId, page1[1].queryId)
	}
	if next == "" {
		t.Fatal("page 1 did not mint a continuation token")
	}

	page2, next2, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2, NextToken: next})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 || page2[0].queryId != fmtQueryId(2) || page2[1].queryId != fmtQueryId(3) {
		t.Fatalf("page 2 = %s,%s want the middle pair", page2[0].queryId, page2[1].queryId)
	}

	page3, next3, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2, NextToken: next2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page3) != 1 || page3[0].queryId != fmtQueryId(4) || next3 != "" {
		t.Fatalf("page 3 = %d items, token %q", len(page3), next3)
	}

	// The scope binds the filter: an unfiltered walk rejects a token
	// minted under a status filter (and vice versa).
	filtered, fnext, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2, StatusFilter: "Complete"})
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 2 || fnext == "" {
		t.Fatalf("filtered page = %d items", len(filtered))
	}
	if _, _, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2, NextToken: fnext}); err == nil {
		t.Fatal("filtered token accepted by an unfiltered walk")
	}
	if _, _, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2, NextToken: "garbage"}); err == nil {
		t.Fatal("garbage token accepted")
	}
}

// A token minted by another operation's vocabulary — here
// DescribeLogStreams' — rejects at DescribeQueries with the operation's
// invalid-parameter error: the version namespace keeps the
// vocabularies mutually undecodable, so the foreign token cannot
// silently restart the walk at the top of the listing.
func TestDescribeQueriesRejectsForeignVocabularyToken(t *testing.T) {
	svc, _ := newReadTestService(t, "queries-group")

	streamTok, err := encodeStreamPageToken(streamPageToken{Group: "queries-group"})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 2, NextToken: streamTok})
	if err == nil {
		t.Fatal("DescribeQueries accepted a stream-vocabulary token")
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.Code != "InvalidParameterException" {
		t.Fatalf("foreign-vocabulary token error = %v, want InvalidParameterException", err)
	}
}

// An equal-timestamp group spanning page boundaries pages without
// duplication or loss: the resume predicate must mirror the sort's
// ascending-id tie-break, so a same-millisecond group (scheduled
// executions carry TriggerTime-millisecond stamps) resumes at the
// larger unserved ids rather than re-serving the served smaller ones.
func TestDescribeQueriesSameTimestampGroupAcrossPages(t *testing.T) {
	svc, _ := newReadTestService(t, "queries-group")

	sameMs := time.Now()
	for i := 9; i >= 0; i-- {
		qs := &queryState{
			queryId:       fmtQueryId(i),
			status:        "Complete",
			createdAt:     sameMs,
			logGroupNames: []string{"queries-group"},
		}
		svc.queries.Store(qs.queryId, qs)
	}

	// The listing's order for one timestamp is the ascending id.
	want := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		want = append(want, fmtQueryId(i))
	}

	var served []string
	next := ""
	for page := 0; ; page++ {
		if page > 10 {
			t.Fatal("same-timestamp walk did not terminate")
		}
		queries, token, err := svc.describeQueriesCore(&DescribeQueriesInput{MaxResults: 3, NextToken: next})
		if err != nil {
			t.Fatal(err)
		}
		for _, qs := range queries {
			served = append(served, qs.queryId)
		}
		if token == "" {
			break
		}
		next = token
	}
	if len(served) != len(want) {
		t.Fatalf("served %d queries, want %d: %v", len(served), len(want), served)
	}
	for i, id := range want {
		if served[i] != id {
			t.Fatalf("served[%d] = %s, want %s (full walk: %v)", i, served[i], id, served)
		}
	}
}

func fmtQueryId(i int) string {
	return "query-" + string(rune('0'+i))
}

// Two byte-identical events are distinct events: the FilteredLogEvent
// eventId identifies the event, not the content, so duplicates must
// carry distinct ids — and a re-read of the same stored events must
// mint the same ids again.
func TestFilterLogEventsDuplicateEventIDs(t *testing.T) {
	const group, stream = "eventid-group", "eventid-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UnixMilli() - 2000
	if _, err := store.PutLogEvents(group, stream, []logsstore.LogEntry{
		{Timestamp: ts, Message: "dup"},
		{Timestamp: ts, Message: "dup"},
	}); err != nil {
		t.Fatal(err)
	}

	read := func(t *testing.T) []string {
		t.Helper()
		reqCtx := request.NewRequestContext(context.Background(), nil, "000000000000", "us-east-1")
		resp, err := svc.FilterLogEvents(context.Background(), reqCtx, &request.ParsedRequest{
			Operation:  "FilterLogEvents",
			Parameters: map[string]interface{}{"logGroupName": group},
		})
		if err != nil {
			t.Fatal(err)
		}
		events := resp.(map[string]interface{})["events"].([]map[string]interface{})
		if len(events) != 2 {
			t.Fatalf("filtered %d events, want the two duplicates", len(events))
		}
		ids := make([]string, len(events))
		for i, ev := range events {
			ids[i] = ev["eventId"].(string)
		}
		return ids
	}

	first := read(t)
	if first[0] == "" || first[0] == first[1] {
		t.Fatalf("byte-identical duplicates share one eventId: %q", first[0])
	}
	second := read(t)
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("re-read minted a different id at %d: %q then %q", i, first[i], second[i])
		}
	}
}

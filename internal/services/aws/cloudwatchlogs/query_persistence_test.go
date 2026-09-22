package cloudwatchlogs

import (
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// waitPersistedRecord polls until the query record carries the wanted
// status, closing the window between the in-memory status flip and the
// write-through that follows it.
func waitPersistedRecord(t *testing.T, store *logsstore.Store, queryId, status string) *logsstore.QueryRecord {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		records, err := store.ListQueryRecords()
		if err != nil {
			t.Fatalf("list query records: %v", err)
		}
		for _, rec := range records {
			if rec.QueryId == queryId && rec.Status == status {
				return rec
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("query record %s never reached status %s", queryId, status)
	return nil
}

// A restart is not the availability boundary for query state: a query the
// previous process completed stays queryable — GetQueryResults serves its
// rows and DescribeQueries lists it — because the terminal state wrote
// through to the store and the fresh construction loaded it back.
func TestQueryStateSurvivesRestart(t *testing.T) {
	sm := newTestStorageManager(t)
	svcA := newTestServiceOnManager(t, sm)

	storeA, err := svcA.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	if err := storeA.CreateLogGroup(logsstore.NewLogGroup("persist-group", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if err := storeA.CreateLogStream(logsstore.NewLogStream("s1", "persist-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, storeA, "persist-group", "s1", time.Now().UnixMilli()-60000, 3)

	now := time.Now().UnixMilli()
	queryId, err := svcA.startQueryCore(&StartQueryInput{
		StartTime:      now/1000 - 3600,
		EndTime:        now/1000 + 60,
		QueryString:    "fields @message",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"persist-group"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	waitTerminalQuery(t, svcA, queryId)
	waitPersistedRecord(t, storeA, queryId, queryStatusComplete)

	// The restart: a fresh service over the same storage.
	svcB := NewLogsService(sm, "000000000000", t.TempDir())
	t.Cleanup(svcB.Stop)

	res, err := svcB.getQueryResultsCore(&GetQueryResultsInput{QueryId: queryId})
	if err != nil {
		t.Fatalf("post-restart GetQueryResults: %v", err)
	}
	if res.Status != queryStatusComplete {
		t.Fatalf("post-restart status = %q, want Complete", res.Status)
	}
	if len(res.Results) != 3 {
		t.Fatalf("post-restart results = %d rows, want 3", len(res.Results))
	}
	for i := range res.Results {
		row := &res.Results[i]
		if got := row.fields["@message"]; got != "m" {
			t.Fatalf("row %d @message = %q, want the stored event text", i, got)
		}
		if len(row.ordered()) == 0 {
			t.Fatalf("row %d lost its column order through the round trip", i)
		}
	}
	// The statistics envelope survives the restart with the record's
	// scanned set: one group scanned, and the result count keeps the
	// whole final set rather than the served page.
	if res.LogGroupsScanned != 1 {
		t.Fatalf("post-restart logGroupsScanned = %d, want 1", res.LogGroupsScanned)
	}
	if res.ResultCount != 3 {
		t.Fatalf("post-restart resultCount = %d, want 3", res.ResultCount)
	}

	// A page smaller than the final set reports the set's size: the count
	// spans every page, never the slice the caller holds.
	paged, err := svcB.getQueryResultsCore(&GetQueryResultsInput{QueryId: queryId, MaxItems: 2})
	if err != nil {
		t.Fatalf("paged GetQueryResults: %v", err)
	}
	if len(paged.Results) != 2 || paged.ResultCount != 3 || paged.NextToken == "" {
		t.Fatalf("paged results = %d rows with resultCount %d, want 2 rows over a 3-row set", len(paged.Results), paged.ResultCount)
	}

	items, _, err := svcB.describeQueriesCore(&DescribeQueriesInput{})
	if err != nil {
		t.Fatalf("post-restart DescribeQueries: %v", err)
	}
	found := false
	for _, qs := range items {
		qs.mu.RLock()
		id, status, group := qs.queryId, qs.status, qs.logGroupNames[0]
		qs.mu.RUnlock()
		if id == queryId {
			found = true
			if status != queryStatusComplete || group != "persist-group" {
				t.Fatalf("post-restart listing = (%q, %q), want (Complete, persist-group)", status, group)
			}
		}
	}
	if !found {
		t.Fatal("post-restart DescribeQueries does not list the pre-restart query")
	}
}

// A RUNNING record a restart orphaned reconciles to Failed — the same rule
// export and import tasks follow — both in the loaded state and in the
// persisted record, and StopQuery on it rejects as an ended query.
func TestRestartMarksOrphanedRunningQueryFailed(t *testing.T) {
	sm := newTestStorageManager(t)
	svcA := newTestServiceOnManager(t, sm)

	storeA, err := svcA.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	orphan := &logsstore.QueryRecord{
		QueryId:           "query-orphan",
		Region:            "us-east-1",
		LogGroupNames:     []string{"persist-group"},
		QueryString:       "fields @message",
		QueryLanguage:     "CWLI",
		Status:            queryStatusRunning,
		CreatedAtUnixNano: time.Now().UnixNano(),
	}
	if err := storeA.PutQueryRecord(orphan); err != nil {
		t.Fatal(err)
	}

	svcB := NewLogsService(sm, "000000000000", t.TempDir())
	t.Cleanup(svcB.Stop)

	val, ok := svcB.queries.Load("query-orphan")
	if !ok {
		t.Fatal("orphaned query record was not loaded into the registry")
	}
	qs := val.(*queryState)
	qs.mu.RLock()
	status := qs.status
	qs.mu.RUnlock()
	if status != queryStatusFailed {
		t.Fatalf("orphaned RUNNING query loaded as %q, want Failed", status)
	}
	rec := waitPersistedRecord(t, storeA, "query-orphan", queryStatusFailed)
	if len(rec.Results) != 0 {
		t.Fatalf("reconciled orphan carries %d result rows", len(rec.Results))
	}
	if err := svcB.stopQueryCore("query-orphan"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("stop of a restart-failed query = %v, want InvalidParameterException", err)
	}
}

// The retention sweep bounds durability as well as memory: a persisted
// record past the retention window is deleted from the store while a
// fresh one survives.
func TestEvictExpiredQueryRecords(t *testing.T) {
	svc, store := newReadTestService(t, "lifecycle-group")

	old := &logsstore.QueryRecord{
		QueryId:           "query-old",
		Region:            "us-east-1",
		Status:            queryStatusComplete,
		CreatedAtUnixNano: time.Now().Add(-logsstore.QueryResultRetentionPeriod - time.Minute).UnixNano(),
	}
	fresh := &logsstore.QueryRecord{
		QueryId:           "query-fresh",
		Region:            "us-east-1",
		Status:            queryStatusComplete,
		CreatedAtUnixNano: time.Now().UnixNano(),
	}
	for _, rec := range []*logsstore.QueryRecord{old, fresh} {
		if err := store.PutQueryRecord(rec); err != nil {
			t.Fatal(err)
		}
	}

	svc.evictExpiredQueries()

	records, err := store.ListQueryRecords()
	if err != nil {
		t.Fatal(err)
	}
	for _, rec := range records {
		if rec.QueryId == "query-old" {
			t.Fatal("expired query record survived the sweep")
		}
	}
	found := false
	for _, rec := range records {
		if rec.QueryId == "query-fresh" {
			found = true
		}
	}
	if !found {
		t.Fatal("fresh query record was evicted")
	}
}

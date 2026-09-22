package cloudwatchlogs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// waitTerminalQuery polls the query state until it leaves Running.
func waitTerminalQuery(t *testing.T, svc *LogsService, queryId string) *queryState {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		val, ok := svc.queries.Load(queryId)
		if !ok {
			t.Fatalf("query %s vanished from the state map", queryId)
		}
		qs := val.(*queryState)
		qs.mu.RLock()
		status := qs.status
		qs.mu.RUnlock()
		if status != queryStatusRunning {
			return qs
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("query %s did not leave Running in time", queryId)
	return nil
}

func logsErrorCode(err error) string {
	var ae *awserrors.AWSError
	if errors.As(err, &ae) {
		return ae.GetCode()
	}
	return ""
}

// StartQuery rejects a log group that does not exist instead of issuing a
// queryId that completes with zero rows.
func TestStartQueryRejectsMissingGroup(t *testing.T) {
	svc, _ := newReadTestService(t, "lifecycle-group")

	_, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      time.Now().Unix() - 3600,
		EndTime:        time.Now().Unix() + 60,
		QueryString:    "fields @message | limit 1",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"lifecycle-group", "no-such-group"},
		Region:         "us-east-1",
	})
	if logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("error = %v, want ResourceNotFoundException", err)
	}
}

// The group-member rules: exactly one of logGroupName, logGroupNames or
// logGroupIdentifiers ("A StartQuery operation must include one of the
// following: Either exactly one of the following parameters... Or the
// queryString must include a SOURCE command to select log groups for the
// query"), at most fifty named groups ("You can include up to 50 log
// groups"), and a SOURCE-first query with no group members starts and
// completes — its groups resolve at execution.
func TestStartQueryGroupMemberRules(t *testing.T) {
	svc, _ := newReadTestService(t, "lifecycle-group")
	now := time.Now().Unix()

	base := StartQueryInput{
		StartTime:      now - 3600,
		EndTime:        now + 60,
		QueryString:    "fields @message | limit 1",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		Region:         "us-east-1",
	}

	rows := []struct {
		name  string
		input StartQueryInput
	}{
		{"name plus names", StartQueryInput{LogGroupName: "lifecycle-group", LogGroupNames: []string{"lifecycle-group"}}},
		{"name plus identifiers", StartQueryInput{LogGroupName: "lifecycle-group", LogGroupIdentifiers: []string{"arn:aws:logs:us-east-1:000000000000:log-group:lifecycle-group"}}},
		{"names plus identifiers", StartQueryInput{LogGroupNames: []string{"lifecycle-group"}, LogGroupIdentifiers: []string{"arn:aws:logs:us-east-1:000000000000:log-group:lifecycle-group"}}},
	}
	for _, row := range rows {
		in := base
		in.LogGroupName = row.input.LogGroupName
		in.LogGroupNames = row.input.LogGroupNames
		in.LogGroupIdentifiers = row.input.LogGroupIdentifiers
		if _, err := svc.startQueryCore(&in); logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("%s: error = %v, want InvalidParameterException", row.name, err)
		}
	}

	noMembers := base
	if _, err := svc.startQueryCore(&noMembers); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("no group members and no SOURCE: error = %v, want InvalidParameterException", err)
	}

	oversized := base
	oversized.LogGroupNames = make([]string, logsstore.QueryLogGroupsMax+1)
	if _, err := svc.startQueryCore(&oversized); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("fifty-one groups: error = %v, want InvalidParameterException", err)
	}

	sourceOnly := base
	sourceOnly.QueryString = "SOURCE logGroups(namePrefix: ['lifecycle-']) | fields @message | limit 1"
	queryId, err := svc.startQueryCore(&sourceOnly)
	if err != nil {
		t.Fatalf("SOURCE-only query rejected: %v", err)
	}
	qs := waitTerminalQuery(t, svc, queryId)
	qs.mu.RLock()
	status, groups := qs.status, qs.scannedGroups
	qs.mu.RUnlock()
	if status != queryStatusComplete {
		t.Fatalf("SOURCE-only query status = %q", status)
	}
	if len(groups) != 1 || groups[0] != "lifecycle-group" {
		t.Fatalf("SOURCE-only scanned groups = %v, want the prefix-matched set", groups)
	}
}

// An inverted time window is malformed input, not an empty query.
func TestStartQueryRejectsInvertedWindow(t *testing.T) {
	svc, _ := newReadTestService(t, "lifecycle-group")
	now := time.Now().Unix()

	for _, tc := range []struct {
		name  string
		start int64
		end   int64
	}{
		{"start after end", now + 600, now},
		{"start equal to end", now, now},
	} {
		_, err := svc.startQueryCore(&StartQueryInput{
			StartTime:      tc.start,
			EndTime:        tc.end,
			QueryString:    "fields @message | limit 1",
			StartTimeSet:   true,
			EndTimeSet:     true,
			QueryStringSet: true,
			LogGroupNames:  []string{"lifecycle-group"},
			Region:         "us-east-1",
		})
		if logsErrorCode(err) != "InvalidParameterException" {
			t.Fatalf("%s: error = %v, want InvalidParameterException", tc.name, err)
		}
	}
}

// The status written by the worker never overwrites a terminal status: a
// query cancelled before its result commit stays Cancelled with empty
// results, and a completed query cannot be re-failed or re-stopped.
func TestQueryTerminalStatusWins(t *testing.T) {
	svc, store := newReadTestService(t, "lifecycle-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "lifecycle-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "lifecycle-group", "s1", time.Now().UnixMilli(), 5)

	// A query state pre-cancelled the way StopQuery leaves it: the worker
	// runs to its commit point and must keep the terminal status.
	queryId := fmt.Sprintf("query-%d", time.Now().UnixNano())
	qs := &queryState{
		queryId:       queryId,
		logGroupNames: []string{"lifecycle-group"},
		queryString:   "fields @message",
		queryLanguage: "CWLI",
		status:        queryStatusCancelled,
		cancelled:     true,
		createdAt:     time.Now(),
	}
	svc.queries.Store(queryId, qs)

	svc.executeQuery(context.Background(), "us-east-1", queryId,
		"fields @message", []string{"lifecycle-group"},
		time.Now().UnixMilli()-3600000, time.Now().UnixMilli()+60000, 10000)

	qs.mu.RLock()
	status, results := qs.status, qs.results
	qs.mu.RUnlock()
	if status != queryStatusCancelled {
		t.Fatalf("status = %q, want Cancelled to survive the worker commit", status)
	}
	if len(results) != 0 {
		t.Fatalf("cancelled query committed %d result rows", len(results))
	}

	// failQuery respects the same rule.
	svc.failQuery(queryId, "late failure")
	qs.mu.RLock()
	status = qs.status
	qs.mu.RUnlock()
	if status != queryStatusCancelled {
		t.Fatalf("status = %q after failQuery, want Cancelled", status)
	}
}

// StopQuery on a Running query cancels it; on an ended query it rejects
// with the not-running error the operation documents.
func TestStopQuerySemantics(t *testing.T) {
	svc, _ := newReadTestService(t, "lifecycle-group")

	seed := func(id, status string) *queryState {
		qs := &queryState{
			queryId:       id,
			logGroupNames: []string{"lifecycle-group"},
			queryString:   "fields @message",
			queryLanguage: "CWLI",
			status:        status,
			createdAt:     time.Now(),
		}
		svc.queries.Store(id, qs)
		return qs
	}

	running := seed("query-running", queryStatusRunning)
	if err := svc.stopQueryCore("query-running"); err != nil {
		t.Fatalf("stop of a running query failed: %v", err)
	}
	running.mu.RLock()
	status, cancelled := running.status, running.cancelled
	running.mu.RUnlock()
	if status != queryStatusCancelled || !cancelled {
		t.Fatalf("status = %q cancelled = %v, want Cancelled with the flag raised", status, cancelled)
	}

	seed("query-complete", queryStatusComplete)
	if err := svc.stopQueryCore("query-complete"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("stop of a completed query = %v, want InvalidParameterException", err)
	}
	seed("query-failed", queryStatusFailed)
	if err := svc.stopQueryCore("query-failed"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("stop of a failed query = %v, want InvalidParameterException", err)
	}

	// Unknown query ids stay ResourceNotFound.
	if err := svc.stopQueryCore("query-nope"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("unknown stop = %v, want ResourceNotFoundException", err)
	}
}

// A normal query completes and keeps its language default.
func TestStartQueryCompletesWithDefaults(t *testing.T) {
	svc, store := newReadTestService(t, "lifecycle-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "lifecycle-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "lifecycle-group", "s1", time.Now().UnixMilli(), 5)

	// The StartQuery members are documented as epoch seconds ("Specified
	// as epoch time, the number of seconds since January 1, 1970") while
	// the stored events carry millisecond timestamps — a seconds window
	// must still select the ingested rows.
	queryId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      time.Now().Unix() - 3600,
		EndTime:        time.Now().Unix() + 60,
		QueryString:    "fields @message | limit 2",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"lifecycle-group"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	qs := waitTerminalQuery(t, svc, queryId)
	qs.mu.RLock()
	status, language, results := qs.status, qs.queryLanguage, qs.results
	qs.mu.RUnlock()
	if status != queryStatusComplete {
		t.Fatalf("status = %q, want Complete", status)
	}
	if language != "CWLI" {
		t.Fatalf("queryLanguage = %q, want the CWLI default", language)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d rows, want the query limit 2", len(results))
	}
}

// DescribeQueries lists scheduled-query executions alongside interactive
// queries: an in-progress execution carries Scheduled, a succeeded one
// Complete; the queryLanguage filter matches the CWLI default; the status
// filter addresses the scheduled entry.
func TestDescribeQueriesSurfacesScheduledExecutions(t *testing.T) {
	svc, store := newReadTestService(t, "lifecycle-group")
	now := time.Now().UnixMilli()

	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-lifecycle",
		Name:                "lifecycle",
		QueryString:         "fields @message",
		LogGroupIdentifiers: []string{"lifecycle-group"},
		State:               "ENABLED",
		QueryLanguage:       "",
		CreationTime:        now,
	}
	if err := store.PutScheduledQuery(sq); err != nil {
		t.Fatal(err)
	}
	if err := store.PutScheduledQueryExecution(&logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          "sq-lifecycle-1",
		TriggerTime:      now,
		Status:           "RUNNING",
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.PutScheduledQueryExecution(&logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          "sq-lifecycle-2",
		TriggerTime:      now - 1000,
		Status:           "SUCCESS",
	}); err != nil {
		t.Fatal(err)
	}

	scheduled, _, err := svc.describeQueriesCore(&DescribeQueriesInput{StatusFilter: "Scheduled"})
	if err != nil {
		t.Fatal(err)
	}
	if len(scheduled) != 1 || scheduled[0].queryId != "sq-lifecycle-1" {
		t.Fatalf("Scheduled filter returned %d entries, want the running execution alone", len(scheduled))
	}
	scheduled[0].mu.RLock()
	language, group := scheduled[0].queryLanguage, scheduled[0].logGroupNames
	scheduled[0].mu.RUnlock()
	if language != "CWLI" {
		t.Fatalf("execution queryLanguage = %q, want the CWLI default", language)
	}
	if len(group) != 1 || group[0] != "lifecycle-group" {
		t.Fatalf("execution groups = %v, want the resolved group", group)
	}

	complete, _, err := svc.describeQueriesCore(&DescribeQueriesInput{StatusFilter: "Complete"})
	if err != nil {
		t.Fatal(err)
	}
	if len(complete) != 1 || complete[0].queryId != "sq-lifecycle-2" {
		t.Fatalf("Complete filter returned %d entries, want the succeeded execution", len(complete))
	}

	// The queryLanguage member filters: a CWLI filter keeps both
	// executions, a foreign language drops them.
	cwli, _, err := svc.describeQueriesCore(&DescribeQueriesInput{QueryLanguage: "CWLI"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cwli) != 2 {
		t.Fatalf("CWLI filter returned %d entries, want both executions", len(cwli))
	}
	ppl, _, err := svc.describeQueriesCore(&DescribeQueriesInput{QueryLanguage: "PPL"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ppl) != 0 {
		t.Fatalf("PPL filter returned %d entries, want none", len(ppl))
	}
}

// StopQuery cancels a running scheduled execution by marking its
// persisted record; an ended execution rejects with not-running.
func TestStopQueryCancelsScheduledExecution(t *testing.T) {
	svc, store := newReadTestService(t, "lifecycle-group")
	now := time.Now().UnixMilli()
	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-cancel",
		Name:                "cancel",
		QueryString:         "fields @message",
		LogGroupIdentifiers: []string{"lifecycle-group"},
		State:               "ENABLED",
		CreationTime:        now,
	}
	if err := store.PutScheduledQuery(sq); err != nil {
		t.Fatal(err)
	}
	if err := store.PutScheduledQueryExecution(&logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          "sq-cancel-1",
		TriggerTime:      now,
		Status:           "RUNNING",
	}); err != nil {
		t.Fatal(err)
	}

	if err := svc.stopQueryCore("sq-cancel-1"); err != nil {
		t.Fatalf("stop of a running execution failed: %v", err)
	}
	execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(execs) != 1 || execs[0].Status != "CANCELLED" {
		t.Fatalf("execution status = %v, want CANCELLED", execs)
	}

	// Stopping the now-cancelled execution rejects: it is not running.
	if err := svc.stopQueryCore("sq-cancel-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("second stop = %v, want InvalidParameterException", err)
	}
}

// The retention sweep drops query state past the documented window and
// keeps fresh state.
func TestEvictExpiredQueries(t *testing.T) {
	svc, _ := newReadTestService(t, "lifecycle-group")

	fresh := &queryState{queryId: "fresh", status: queryStatusComplete, createdAt: time.Now()}
	old := &queryState{queryId: "old", status: queryStatusComplete,
		createdAt: time.Now().Add(-logsstore.QueryResultRetentionPeriod - time.Minute)}
	svc.queries.Store("fresh", fresh)
	svc.queries.Store("old", old)

	svc.evictExpiredQueries()

	if _, ok := svc.queries.Load("fresh"); !ok {
		t.Fatal("fresh query state was evicted")
	}
	if _, ok := svc.queries.Load("old"); ok {
		t.Fatal("expired query state survived the sweep")
	}
}

// The ticker enumerates configured regions rather than the store cache: a
// scheduled query in a region whose store the service never resolved
// still fires, and a still-running execution suppresses a duplicate fire.
func TestTickScheduledQueriesFiresForUntouchedRegion(t *testing.T) {
	svc, store := newReadTestService(t, "lifecycle-group")

	// The schedule is due immediately: rate(1 minute) with a creation
	// time in the past has elapsed its first boundary.
	now := time.Now()
	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-tick",
		Name:                "tick",
		QueryString:         "fields @message",
		LogGroupIdentifiers: []string{"lifecycle-group"},
		State:               "ENABLED",
		ScheduleExpression:  "rate(1 minute)",
		CreationTime:        now.Add(-2 * time.Minute).UnixMilli(),
	}
	if err := store.PutScheduledQuery(sq); err != nil {
		t.Fatal(err)
	}

	// Drop the region's store from the instance cache to prove the tick
	// is not cache-bound.
	svc.logsStores.Delete("us-east-1")

	svc.tickScheduledQueries()

	deadline := time.Now().Add(5 * time.Second)
	var execs []*logsstore.ScheduledQueryExecution
	for time.Now().Before(deadline) {
		var err error
		execs, err = store.ListScheduledQueryExecutions(sq.Id, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(execs) == 1 && logsstore.IsTerminalScheduledExecutionStatus(execs[0].Status) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if len(execs) != 1 {
		t.Fatalf("executions = %d, want exactly one fire", len(execs))
	}
	if execs[0].Status != "SUCCESS" {
		t.Fatalf("execution status = %q, want SUCCESS", execs[0].Status)
	}

	// Reset the consumed boundary marker while an execution is RUNNING
	// again (as if a slow one were in flight): the boundary looks due, but
	// the in-flight guard must keep the tick from duplicating the fire.
	execs[0].Status = "RUNNING"
	if err := store.PutScheduledQueryExecution(execs[0]); err != nil {
		t.Fatal(err)
	}
	if err := store.MutateScheduledQuery(sq.Id, func(q *logsstore.ScheduledQuery) error {
		q.LastExecutedBoundary = 0
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	svc.tickScheduledQueries()
	time.Sleep(100 * time.Millisecond)
	execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(execs) != 1 {
		t.Fatalf("executions after guarded tick = %d, want the in-flight one alone", len(execs))
	}
}

// CreateScheduledQuery carries the model's required members (queryString,
// queryLanguage, scheduleExpression, executionRoleArn) and rejects spec
// violations with ValidationException — the identity its declared error
// list carries, not InvalidParameterException.
func TestCreateScheduledQueryRequiredMembersAndErrorIdentity(t *testing.T) {
	svc, store := newReadTestService(t, "sq-create-contract-group")
	base := func() *CreateScheduledQueryInput {
		return &CreateScheduledQueryInput{
			Name:               "sq-contract",
			QueryString:        "fields @message",
			QueryLanguage:      "CWLI",
			ScheduleExpression: "rate(1 minute)",
			ExecutionRoleArn:   "arn:aws:iam::000000000000:role/deliver",
		}
	}

	missing := func(missing string, in *CreateScheduledQueryInput) {
		t.Helper()
		_, err := svc.createScheduledQueryCore(store, in)
		if err == nil {
			t.Fatalf("create with missing %s succeeded", missing)
		}
		if code := logsErrorCode(err); code != "ValidationException" {
			t.Fatalf("missing %s: code = %s, want ValidationException", missing, code)
		}
	}
	in := base()
	in.QueryLanguage = ""
	missing("queryLanguage", in)
	in = base()
	in.ExecutionRoleArn = ""
	missing("executionRoleArn", in)
	in = base()
	in.State = "GARBAGE"
	missing("state", in)
	in = base()
	in.Name = fmt.Sprintf("%300dc", 1)
	missing("name length", in)
	in = base()
	in.LogGroupIdentifiers = make([]string, 51)
	missing("logGroupIdentifiers count", in)
	// ScheduledQueryLogGroupIdentifiers carries @length {1,50} ("Array
	// Members: Minimum number of 1 item."): an explicitly empty list (the
	// non-nil empty slice the parse marks present) is a member violation,
	// not an absent member.
	in = base()
	in.LogGroupIdentifiers = []string{}
	missing("empty logGroupIdentifiers", in)
	// "The timezone for evaluating the schedule expression" names an IANA
	// zone; a value that resolves to no zone — or the Go runtime alias
	// for the host clock — rejects at the door instead of silently
	// evaluating the schedule in UTC.
	in = base()
	in.Timezone = "Bogus/Zone"
	missing("timezone", in)
	in = base()
	in.Timezone = "Local"
	missing("timezone Local alias", in)
	// The member documents the cron grammar; the at() one-shot form is
	// the EventBridge Scheduler profile's alone.
	in = base()
	in.ScheduleExpression = "at(2030-01-01T00:00:00)"
	missing("at() schedule expression", in)

	if _, err := svc.createScheduledQueryCore(store, base()); err != nil {
		t.Fatalf("valid create rejected: %v", err)
	}
	// A named IANA zone is the documented value form and stays servable.
	zoned := base()
	zoned.Name = "sq-contract-tz"
	zoned.Timezone = "Asia/Tokyo"
	if created, err := svc.createScheduledQueryCore(store, zoned); err != nil || created.Timezone != "Asia/Tokyo" {
		t.Fatalf("IANA timezone create: created=%v err=%v", created, err)
	}
	cronForm := base()
	cronForm.Name = "sq-contract-cron"
	cronForm.ScheduleExpression = "cron(0 6 * * ? *)"
	if _, err := svc.createScheduledQueryCore(store, cronForm); err != nil {
		t.Fatalf("cron() schedule expression create rejected: %v", err)
	}
	// "The name must be unique within your account and region" — the
	// duplicate create rejects with the declared ConflictException while
	// a differently named one succeeds.
	if _, err := svc.createScheduledQueryCore(store, base()); logsErrorCode(err) != "ConflictException" {
		t.Fatalf("duplicate name: error = %v, want ConflictException", err)
	}
	renamed := base()
	renamed.Name = "sq-contract-two"
	if _, err := svc.createScheduledQueryCore(store, renamed); err != nil {
		t.Fatalf("differently named create rejected: %v", err)
	}
	// The update path rides the same spec validator, so a timezone that
	// would degrade silently cannot enter through an update either.
	_, err := svc.updateScheduledQueryCore(store, &UpdateScheduledQueryInput{
		Identifier:         "sq-contract",
		QueryString:        "fields @message",
		QueryLanguage:      "CWLI",
		ScheduleExpression: "rate(1 minute)",
		ExecutionRoleArn:   "arn:aws:iam::000000000000:role/deliver",
		Timezone:           "Bogus/Zone",
	})
	if logsErrorCode(err) != "ValidationException" {
		t.Fatalf("update with bogus timezone: error = %v, want ValidationException", err)
	}
	// The empty update list is the same @length(min 1) violation: it
	// rejects instead of clearing the query's group set.
	_, err = svc.updateScheduledQueryCore(store, &UpdateScheduledQueryInput{
		Identifier:          "sq-contract",
		QueryString:         "fields @message",
		QueryLanguage:       "CWLI",
		ScheduleExpression:  "rate(1 minute)",
		ExecutionRoleArn:    "arn:aws:iam::000000000000:role/deliver",
		LogGroupIdentifiers: []string{},
	})
	if logsErrorCode(err) != "ValidationException" {
		t.Fatalf("update with empty logGroupIdentifiers: error = %v, want ValidationException", err)
	}
}

// The required members' zero values sit inside their documented domains
// ("Valid Range: Minimum value of 0" for the epoch pair, "Minimum length
// of 0" for queryString), so a present zero or empty value runs while an
// omitted member rejects as missing — and a negative or oversized value
// rejects on its own rule.
func TestStartQueryPresentZeroMembers(t *testing.T) {
	svc, _ := newReadTestService(t, "lifecycle-group")
	now := time.Now().Unix()

	base := StartQueryInput{
		StartTime:      now - 3600,
		EndTime:        now + 60,
		QueryString:    "fields @message | limit 1",
		LogGroupNames:  []string{"lifecycle-group"},
		Region:         "us-east-1",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
	}

	epochZero := base
	epochZero.StartTime = 0
	if _, err := svc.startQueryCore(&epochZero); err != nil {
		t.Fatalf("present-zero startTime rejected: %v", err)
	}

	emptyQuery := base
	emptyQuery.QueryString = ""
	if _, err := svc.startQueryCore(&emptyQuery); err != nil {
		t.Fatalf("present-empty queryString rejected: %v", err)
	}

	omitted := base
	omitted.StartTimeSet = false
	if _, err := svc.startQueryCore(&omitted); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("omitted startTime: error = %v, want InvalidParameterException", err)
	}

	negative := base
	negative.EndTime = -1
	if _, err := svc.startQueryCore(&negative); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("negative endTime: error = %v, want InvalidParameterException", err)
	}

	oversized := base
	oversized.QueryString = strings.Repeat("a", 10001)
	if _, err := svc.startQueryCore(&oversized); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("oversized queryString: error = %v, want InvalidParameterException", err)
	}
}

// The sixty-minute runtime bound ("Queries time out after 60 minutes of
// runtime") stamps the Timeout status the status vocabulary carries for
// it — a terminal sibling of Failed, never overwritten by a late
// completion.
func TestTimeoutQueryStampsTimeoutStatus(t *testing.T) {
	svc, _ := newReadTestService(t, "timeout-group")
	qs := &queryState{queryId: "timeout-run", status: queryStatusRunning, createdAt: time.Now()}
	svc.queries.Store("timeout-run", qs)
	svc.timeoutQuery("timeout-run")
	loaded, ok := svc.queries.Load("timeout-run")
	if !ok {
		t.Fatal("query left the registry")
	}
	got := loaded.(*queryState)
	got.mu.RLock()
	status := got.status
	got.mu.RUnlock()
	if status != queryStatusTimeout {
		t.Fatalf("status = %s, want Timeout", status)
	}
	if !isTerminalQueryStatus(queryStatusTimeout) {
		t.Fatal("Timeout must be a terminal status")
	}
}

// The worker reaches the Timeout branch of its execution ladder: a query
// whose pipeline breaches the sixty-minute bound comes back Timeout, not
// Failed — the breach is not routed through the error path. The deadline
// clock seam advances the pipeline past the worker's mint so the breach
// fires deterministically without the real sixty-minute wait.
func TestExecuteQueryStampsTimeoutOnDeadlineBreach(t *testing.T) {
	svc, store := newReadTestService(t, "timeout-worker-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "timeout-worker-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "timeout-worker-group", "s1", time.Now().UnixMilli(), 3)

	restore := queryDeadlineNow
	mints := 0
	queryDeadlineNow = func() time.Time {
		mints++
		if mints == 1 {
			return restore() // the worker's deadline mint
		}
		return restore().Add(queryRuntimeLimit + time.Minute) // past the outer bound
	}
	defer func() { queryDeadlineNow = restore }()

	queryId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      time.Now().Unix() - 3600,
		EndTime:        time.Now().Unix() + 60,
		QueryString:    "fields @message | limit 2",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"timeout-worker-group"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	qs := waitTerminalQuery(t, svc, queryId)
	qs.mu.RLock()
	status := qs.status
	qs.mu.RUnlock()
	if status != queryStatusTimeout {
		t.Fatalf("status = %q, want Timeout", status)
	}
}

// The execution window honours the documented second units of the
// offsets ("The time offset in seconds that defines the lookback period
// for the query", CreateScheduledQuery startTimeOffset/endTimeOffset):
// a 3600-second lookback is one hour of the millisecond execution
// clock, not 3.6 seconds.
func TestScheduledExecutionWindowConvertsSecondOffsets(t *testing.T) {
	now := int64(1_700_000_000_000)
	start, end := scheduledExecutionWindow(now, &logsstore.ScheduledQuery{})
	if start != now-60*60*1000 || end != now {
		t.Fatalf("default window: %d..%d", start, end)
	}
	start, end = scheduledExecutionWindow(now, &logsstore.ScheduledQuery{StartTimeOffset: 3600, EndTimeOffset: 60})
	if start != now-3_600_000 {
		t.Fatalf("startTimeOffset 3600 must look back one hour, got %d ms", now-start)
	}
	if end != now-60_000 {
		t.Fatalf("endTimeOffset 60 must offset one minute, got %d ms", now-end)
	}
}

// The scheduled plane runs under the same documented sixty-minute outer
// runtime as the interactive one: an execution whose pipeline breaches
// the bound terminally reports the enum's Timeout member and consumes
// its boundary, instead of holding a RUNNING record indefinitely. The
// deadline clock seam advances the pipeline past the execution's mint so
// the breach fires deterministically without the real sixty-minute wait.
func TestScheduledQueryExecutionTimesOutAtRuntimeLimit(t *testing.T) {
	svc, store := newReadTestService(t, "scheduled-timeout-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "scheduled-timeout-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "scheduled-timeout-group", "s1", time.Now().UnixMilli(), 3)

	restore := queryDeadlineNow
	mints := 0
	queryDeadlineNow = func() time.Time {
		mints++
		if mints == 1 {
			return restore() // the execution's deadline mint
		}
		return restore().Add(queryRuntimeLimit + time.Minute) // past the outer bound
	}
	defer func() { queryDeadlineNow = restore }()

	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-timeout",
		Name:                "sq-timeout",
		QueryString:         "fields @message | limit 2",
		LogGroupIdentifiers: []string{"scheduled-timeout-group"},
		ScheduleExpression:  "rate(1 minute)",
		State:               logsstore.ScheduledQueryStateEnabled,
		CreationTime:        time.Now().UnixMilli(),
	}
	if err := store.PutScheduledQuery(sq); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().UnixMilli()
	exec := &logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          fmt.Sprintf("sq-%s-%d", sq.Id, now),
		TriggerTime:      now,
		Status:           logsstore.ScheduledExecutionStatusRunning,
	}
	if err := store.PutScheduledQueryExecution(exec); err != nil {
		t.Fatal(err)
	}
	boundary := time.UnixMilli(now).UTC().Add(time.Minute)
	svc.runScheduledQueryExecution(context.Background(), "us-east-1", store, sq, exec, now, boundary)

	execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
	if err != nil || len(execs) != 1 {
		t.Fatalf("executions: %v (%d)", err, len(execs))
	}
	if execs[0].Status != logsstore.ScheduledExecutionStatusTimeout {
		t.Fatalf("execution status = %q, want TIMEOUT", execs[0].Status)
	}
	if !logsstore.IsTerminalScheduledExecutionStatus(execs[0].Status) {
		t.Fatal("TIMEOUT must be a terminal execution status")
	}
	if got := mapExecutionStatus(execs[0].Status); got != logsstore.ScheduledQueryStatusTimeout {
		t.Fatalf("mapExecutionStatus(TIMEOUT) = %q, want Timeout", got)
	}
	got, err := store.GetScheduledQuery(sq.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got.LastExecutedBoundary != boundary.UnixMilli() {
		t.Fatalf("LastExecutedBoundary = %d, want the consumed boundary %d", got.LastExecutedBoundary, boundary.UnixMilli())
	}
	if got.LastExecutionStatus != logsstore.ScheduledQueryStatusTimeout {
		t.Fatalf("lastExecutionStatus = %q, want Timeout", got.LastExecutionStatus)
	}
}

// A failed execution consumes its boundary: the query-failure path
// records the trigger so a persistently failing query retries on the
// next schedule boundary instead of re-firing the claimed boundary on
// every worker tick — the same rule the cancellation and
// delivery-failure paths already enforce.
func TestScheduledQueryFailureConsumesBoundary(t *testing.T) {
	svc, store := newReadTestService(t, "scheduled-fail-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "scheduled-fail-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "scheduled-fail-group", "s1", time.Now().UnixMilli(), 3)

	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-fail",
		Name:                "sq-fail",
		QueryString:         "fields @message | frobnicate",
		LogGroupIdentifiers: []string{"scheduled-fail-group"},
		ScheduleExpression:  "rate(1 minute)",
		State:               logsstore.ScheduledQueryStateEnabled,
		CreationTime:        time.Now().UnixMilli(),
	}
	if err := store.PutScheduledQuery(sq); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().UnixMilli()
	exec := &logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          fmt.Sprintf("sq-%s-%d", sq.Id, now),
		TriggerTime:      now,
		Status:           logsstore.ScheduledExecutionStatusRunning,
	}
	if err := store.PutScheduledQueryExecution(exec); err != nil {
		t.Fatal(err)
	}
	boundary := time.UnixMilli(now).UTC().Add(time.Minute)
	svc.runScheduledQueryExecution(context.Background(), "us-east-1", store, sq, exec, now, boundary)

	execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
	if err != nil || len(execs) != 1 {
		t.Fatalf("executions: %v (%d)", err, len(execs))
	}
	if execs[0].Status != logsstore.ScheduledExecutionStatusFailed {
		t.Fatalf("execution status = %q, want FAILED", execs[0].Status)
	}
	got, err := store.GetScheduledQuery(sq.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got.LastExecutedBoundary != boundary.UnixMilli() {
		t.Fatalf("LastExecutedBoundary = %d, want the consumed boundary %d — a failed execution must not re-fire every tick", got.LastExecutedBoundary, boundary.UnixMilli())
	}
	if got.LastExecutionStatus != logsstore.ScheduledQueryStatusFailed {
		t.Fatalf("lastExecutionStatus = %q, want Failed", got.LastExecutionStatus)
	}
	// With the boundary consumed, the query is not due again until the
	// next boundary: the marker, not the failure, decides.
	if _, due := scheduledQueryDue(got, boundary.Add(30*time.Second)); due {
		t.Fatal("query re-fired within the consumed boundary's period")
	}
}

// The scheduled-query operations' Identifier is documented ARN-or-name:
// the bare name form resolves to the id-keyed record, so the console and
// CLI name form reads, updates and deletes a query that exists.
func TestScheduledQueryNameFormIdentifierResolves(t *testing.T) {
	svc, store := newReadTestService(t, "sq-name-form-group")
	in := &CreateScheduledQueryInput{
		Name:               "name-form-query",
		QueryString:        "fields @message",
		QueryLanguage:      "CWLI",
		ScheduleExpression: "rate(1 minute)",
		ExecutionRoleArn:   "arn:aws:iam::000000000000:role/deliver",
	}
	if _, err := svc.createScheduledQueryCore(store, in); err != nil {
		t.Fatalf("create: %v", err)
	}

	id, err := resolveScheduledQueryRef(store, "name-form-query")
	if err != nil {
		t.Fatal(err)
	}
	if id == "" || id == "name-form-query" {
		t.Fatalf("name form did not resolve to the record id: %q", id)
	}
	sq, err := svc.getScheduledQueryCore(store, id)
	if err != nil {
		t.Fatalf("get by resolved id: %v", err)
	}
	if sq.Name != "name-form-query" {
		t.Fatalf("resolved record name: %q", sq.Name)
	}

	// The update core accepts the name form the same way (the update
	// validation requires the full spec members, so the input re-carries
	// them).
	updated, err := svc.updateScheduledQueryCore(store, &UpdateScheduledQueryInput{
		Identifier:         "name-form-query",
		QueryString:        "fields @message",
		QueryLanguage:      "CWLI",
		ScheduleExpression: "rate(1 minute)",
		ExecutionRoleArn:   "arn:aws:iam::000000000000:role/deliver",
		State:              "DISABLED",
	})
	if err != nil {
		t.Fatalf("update through the name form: %v", err)
	}
	if updated.State != "DISABLED" || updated.Id != id {
		t.Fatalf("updated record: %+v", updated)
	}

	// An unknown name keeps its input so the lookup reports not-found.
	unknown, err := resolveScheduledQueryRef(store, "never-was")
	if err != nil {
		t.Fatal(err)
	}
	if unknown != "never-was" {
		t.Fatalf("unknown name must pass through, got %q", unknown)
	}
	if _, err := svc.getScheduledQueryCore(store, unknown); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("unknown name must report not-found, got %v", err)
	}

	if err := svc.deleteScheduledQueryCore(store, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

// DescribeQueries reports the full statistics envelope: createTime in
// epoch seconds (the documented example carries ten-digit values), the
// query's duration, the bytes it scanned, and the ARN of the principal
// that started it.
func TestDescribeQueriesStatisticsMembers(t *testing.T) {
	svc, store := newReadTestService(t, "dq-stats-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "dq-stats-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "dq-stats-group", "s1", time.Now().UnixMilli()-60000, 3)

	reqCtx := request.NewRequestContext(context.Background(), nil, "000000000000", "us-east-1")
	reqCtx.PrincipalType = request.PrincipalTypeUser
	reqCtx.Principal = "alice"

	now := time.Now().Unix()
	startResp, err := svc.StartQuery(context.Background(), reqCtx, vocabRequest(map[string]interface{}{
		"startTime":     now - 3600,
		"endTime":       now + 60,
		"queryString":   "fields @message",
		"logGroupNames": []interface{}{"dq-stats-group"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	queryId := startResp.(map[string]interface{})["queryId"].(string)
	waitTerminalQuery(t, svc, queryId)

	listResp, err := svc.DescribeQueries(context.Background(), reqCtx, vocabRequest(nil))
	if err != nil {
		t.Fatal(err)
	}
	var entry map[string]interface{}
	for _, q := range listResp.(map[string]interface{})["queries"].([]map[string]interface{}) {
		if q["queryId"] == queryId {
			entry = q
			break
		}
	}
	if entry == nil {
		t.Fatal("the started query is missing from the DescribeQueries listing")
	}
	createTime, ok := entry["createTime"].(int64)
	if !ok || createTime < now-120 || createTime > now+120 {
		t.Fatalf("createTime = %#v, want epoch seconds within two minutes of now", entry["createTime"])
	}
	duration, ok := entry["queryDuration"].(int64)
	if !ok || duration < 0 {
		t.Fatalf("queryDuration = %#v, want a non-negative millisecond count", entry["queryDuration"])
	}
	scanned, ok := entry["bytesScanned"].(int64)
	if !ok || scanned <= 0 {
		t.Fatalf("bytesScanned = %#v, want the scanned byte volume", entry["bytesScanned"])
	}
	if identity, _ := entry["userIdentity"].(string); identity != "arn:aws:iam::000000000000:user/alice" {
		t.Fatalf("userIdentity = %#v, want the started principal's ARN", entry["userIdentity"])
	}
}

// The account's concurrent-query quota: the hundred-and-first running
// query rejects with the operation's declared LimitExceededException, and
// a slot freed by a terminal transition admits the next start.
func TestStartQueryConcurrentQuota(t *testing.T) {
	svc, _ := newReadTestService(t, "quota-group")

	for i := 0; i < logsstore.MaxConcurrentQueries; i++ {
		svc.queries.Store(fmt.Sprintf("quota-running-%d", i), &queryState{
			queryId:   fmt.Sprintf("quota-running-%d", i),
			status:    queryStatusRunning,
			createdAt: time.Now(),
		})
	}

	now := time.Now().Unix()
	_, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      now - 3600,
		EndTime:        now + 60,
		QueryString:    "fields @message | limit 1",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"quota-group"},
		Region:         "us-east-1",
	})
	if logsErrorCode(err) != "LimitExceededException" {
		t.Fatalf("error = %v, want LimitExceededException", err)
	}

	freed := &queryState{queryId: "quota-running-0", status: queryStatusRunning, createdAt: time.Now()}
	freed.mu.Lock()
	freed.markTerminal(queryStatusComplete)
	freed.mu.Unlock()
	svc.queries.Store("quota-running-0", freed)

	queryId, err := svc.startQueryCore(&StartQueryInput{
		StartTime:      now - 3600,
		EndTime:        now + 60,
		QueryString:    "fields @message | limit 1",
		StartTimeSet:   true,
		EndTimeSet:     true,
		QueryStringSet: true,
		LogGroupNames:  []string{"quota-group"},
		Region:         "us-east-1",
	})
	if err != nil {
		t.Fatalf("start after freeing a slot: %v", err)
	}
	waitTerminalQuery(t, svc, queryId)
}

// The record resolution behind GetLogRecord and GetLogObject validates
// the pointer against the store: the group, the stream and the event
// itself (an identical message at the pointer's timestamp) must exist —
// a well-formed pointer to deleted or fabricated content answers the
// operations' declared ResourceNotFoundException instead of a
// reconstructed 200 record.
func TestResolveEventPointerRecordValidatesAgainstStore(t *testing.T) {
	const group, stream = "ptr-group", "ptr-stream"
	svc, store := newReadTestService(t, group)
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, group)); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	const msg = `{"lvl":"INFO"}`
	if _, err := store.PutLogEvents(group, stream, []logsstore.LogEntry{{Timestamp: now, Message: msg}}); err != nil {
		t.Fatal(err)
	}
	ts := fmt.Sprintf("%d", now)

	record, err := svc.resolveEventPointerCore(store, "000000000000", group, stream, ts, msg)
	if err != nil {
		t.Fatalf("live pointer rejected: %v", err)
	}
	if record["@message"] != msg || record["@logStream"] != stream || record["lvl"] != "INFO" {
		t.Fatalf("record fields: %v", record)
	}

	if _, err := svc.resolveEventPointerCore(store, "000000000000", "no-such-group", stream, ts, msg); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("missing group: %v", err)
	}
	if _, err := svc.resolveEventPointerCore(store, "000000000000", group, "no-such-stream", ts, msg); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("missing stream: %v", err)
	}
	if _, err := svc.resolveEventPointerCore(store, "000000000000", group, stream, ts, "never written"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("fabricated event: %v", err)
	}
	if _, err := svc.resolveEventPointerCore(store, "000000000000", group, stream, "not-a-ts", msg); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("malformed timestamp: %v", err)
	}
}

// The shutdown checkpoint terminalises the execution record before it
// consumes the boundary: a cancelled context after the query succeeded
// must leave a FAILED record (the restart reconciliation's state written
// in place), not a RUNNING orphan pairing a consumed boundary with an
// in-flight guard that suppresses every later evaluation.
func TestScheduledQueryShutdownFinalisesExecution(t *testing.T) {
	svc, store := newReadTestService(t, "scheduled-shutdown-group")
	if err := store.CreateLogStream(logsstore.NewLogStream("s1", "scheduled-shutdown-group")); err != nil {
		t.Fatal(err)
	}
	putReadEvents(t, store, "scheduled-shutdown-group", "s1", time.Now().UnixMilli(), 3)

	sq := &logsstore.ScheduledQuery{
		Id:                  "sq-shutdown",
		Name:                "sq-shutdown",
		QueryString:         "fields @message | limit 2",
		LogGroupIdentifiers: []string{"scheduled-shutdown-group"},
		ScheduleExpression:  "rate(1 minute)",
		State:               logsstore.ScheduledQueryStateEnabled,
		CreationTime:        time.Now().UnixMilli(),
	}
	if err := store.PutScheduledQuery(sq); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().UnixMilli()
	exec := &logsstore.ScheduledQueryExecution{
		ScheduledQueryId: sq.Id,
		QueryId:          fmt.Sprintf("sq-%s-%d", sq.Id, now),
		TriggerTime:      now,
		Status:           logsstore.ScheduledExecutionStatusRunning,
	}
	if err := store.PutScheduledQueryExecution(exec); err != nil {
		t.Fatal(err)
	}
	// The service is shutting down: the task's context is already
	// cancelled when the checkpoint after query execution reads it.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	boundary := time.UnixMilli(now).UTC().Add(time.Minute)
	svc.runScheduledQueryExecution(ctx, "us-east-1", store, sq, exec, now, boundary)

	execs, err := store.ListScheduledQueryExecutions(sq.Id, 0, 0)
	if err != nil || len(execs) != 1 {
		t.Fatalf("executions: %v (%d)", err, len(execs))
	}
	if execs[0].Status != logsstore.ScheduledExecutionStatusFailed {
		t.Fatalf("execution status = %q, want FAILED", execs[0].Status)
	}
	if !strings.Contains(execs[0].ErrorMessage, "shutdown") {
		t.Fatalf("error message = %q, want the shutdown interruption reason", execs[0].ErrorMessage)
	}
	if scheduledExecutionInFlight(store, sq.Id) {
		t.Fatal("no execution record may remain RUNNING after the shutdown checkpoint")
	}
	got, err := store.GetScheduledQuery(sq.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got.LastExecutedBoundary != boundary.UnixMilli() {
		t.Fatalf("LastExecutedBoundary = %d, want the consumed boundary %d", got.LastExecutedBoundary, boundary.UnixMilli())
	}
}

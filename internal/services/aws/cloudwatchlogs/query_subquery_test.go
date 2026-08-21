package cloudwatchlogs

import (
	"strconv"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

// Tests for subqueries and the join / source command family.

func newTestCtx(events []logEventWithContext) *execContext {
	return &execContext{
		startTime:     0,
		endTime:       1,
		events:        events,
		defaultGroups: []string{"test-group"},
		fetchEvents: func(groups []string, start, end int64) ([]logEventWithContext, error) {
			return events, nil
		},
		listLogGroups: func() ([]sourceGroupInfo, error) {
			return []sourceGroupInfo{{Name: "test-group"}}, nil
		},
		subqueryCache: map[string][]interface{}{},
	}
}

func TestFilterInSubquery(t *testing.T) {
	events := testEvents(
		`{"req": "r1", "kind": "a"}`,
		`{"req": "r2", "kind": "a"}`,
		`{"req": "r1", "kind": "b"}`,
	)
	// The subquery selects req values of kind=b rows; the outer filter
	// keeps kind=a rows whose req appears there.
	rows, err := executeQueryContext(newTestCtx(events),
		`filter kind = "a" and req in (filter kind = "b" | fields req)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rowField(t, rows, 0, "req") != "r1" {
		t.Fatalf("subquery filter rows = %+v", rows)
	}
}

func TestJoinWithSource(t *testing.T) {
	left := testEvents(`{"requestId": "q1", "status": 500}`)
	right := []logEventWithContext{{timestamp: 1, message: `{"requestId": "q1", "duration": 42}`}}
	ctx := newTestCtx(left)
	ctx.fetchEvents = func(groups []string, start, end int64) ([]logEventWithContext, error) {
		if len(groups) > 0 && groups[0] == "right-group" {
			return right, nil
		}
		return left, nil
	}
	ctx.listLogGroups = func() ([]sourceGroupInfo, error) {
		return []sourceGroupInfo{{Name: "right-group"}, {Name: "test-group"}}, nil
	}
	rows, err := executeQueryContext(ctx,
		`join type=inner left=api right=lambda where api.requestId=lambda.requestId (SOURCE logGroups(namePrefix: ['right']))`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("join rows = %d, want 1", len(rows))
	}
	if rowField(t, rows, 0, "lambda.duration") != "42" {
		t.Fatalf("joined fields: %+v", rows[0].fields)
	}
}

func TestAppendCols(t *testing.T) {
	events := testEvents(`{"svc": "a"}`, `{"svc": "b"}`)
	ctx := newTestCtx(events)
	rows, err := executeQueryContext(ctx,
		`fields svc | appendcols override=true ( fields svc | stats count(*) as total )`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d", len(rows))
	}
	if _, ok := rows[0].fields["total"]; !ok {
		t.Fatalf("appendcols should add total: %+v", rows[0].fields)
	}
}

// --- SOURCE ---

func TestSourceNamePrefix(t *testing.T) {
	fetched := map[string]bool{}
	ctx := &execContext{
		startTime:     0,
		endTime:       1,
		defaultGroups: []string{"default"},
		fetchEvents: func(groups []string, start, end int64) ([]logEventWithContext, error) {
			for _, g := range groups {
				fetched[g] = true
			}
			return nil, nil
		},
		listLogGroups: func() ([]sourceGroupInfo, error) {
			return []sourceGroupInfo{
				{Name: "/aws/lambda/app1"},
				{Name: "/aws/lambda/app2"},
				{Name: "other-group"},
			}, nil
		},
		subqueryCache: map[string][]interface{}{},
	}
	if _, err := executeQueryContext(ctx, `SOURCE logGroups(namePrefix: ['/aws/lambda']) | fields @message`); err != nil {
		t.Fatal(err)
	}
	if !fetched["/aws/lambda/app1"] || !fetched["/aws/lambda/app2"] || fetched["other-group"] {
		t.Fatalf("SOURCE selected wrong groups: %v", fetched)
	}
}

func TestSourceTagFilters(t *testing.T) {
	cmdAny, err := parseCommand(mustLex(t, "SOURCE logGroupTags([{key: \"env\", values: [\"prod*\", \"!dev\"]}])"))
	if err != nil {
		t.Fatal(err)
	}
	src := cmdAny.(*sourceCommand)
	if !tagFiltersMatch(src.tagFilters, map[string]string{"env": "production"}) {
		t.Fatal("prod* should match production")
	}
	if tagFiltersMatch(src.tagFilters, map[string]string{"env": "dev"}) {
		t.Fatal("!dev should reject dev")
	}
	if tagFiltersMatch(src.tagFilters, map[string]string{}) {
		t.Fatal("missing tag should not match")
	}
}

func TestSourceTagKeyOnlyFilter(t *testing.T) {
	// The values array is optional: a key-only filter matches every group
	// that carries the tag, regardless of the value.
	cmdAny, err := parseCommand(mustLex(t, "SOURCE logGroupTags([{key: \"env\"}])"))
	if err != nil {
		t.Fatal(err)
	}
	src := cmdAny.(*sourceCommand)
	if !tagFiltersMatch(src.tagFilters, map[string]string{"env": "anything"}) {
		t.Fatal("key-only filter should match a group carrying the tag")
	}
	if tagFiltersMatch(src.tagFilters, map[string]string{"other": "x"}) {
		t.Fatal("key-only filter should not match a group without the tag")
	}
}

// TestSourceReplacesRowSet verifies that the SOURCE selection determines
// the row set, not just the fetch arguments: only events from the selected
// groups appear in the results.
func TestSourceReplacesRowSet(t *testing.T) {
	lambdaEvents := testEvents(`{"svc": "lambda-a"}`)
	otherEvents := testEvents(`{"svc": "other-b"}`)
	ctx := newTestCtx(otherEvents)
	ctx.fetchEvents = func(groups []string, start, end int64) ([]logEventWithContext, error) {
		for _, g := range groups {
			if g == "/aws/lambda/app1" {
				return lambdaEvents, nil
			}
		}
		return otherEvents, nil
	}
	ctx.listLogGroups = func() ([]sourceGroupInfo, error) {
		return []sourceGroupInfo{
			{Name: "/aws/lambda/app1"},
			{Name: "other-group"},
		}, nil
	}
	rows, err := executeQueryContext(ctx, `SOURCE logGroups(namePrefix: ['/aws/lambda']) | fields svc`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].fields["svc"] != "lambda-a" {
		t.Fatalf("SOURCE should replace the row set with selected groups' events, got %+v", rows)
	}
}

func TestSourceQuotas(t *testing.T) {
	cases := []struct {
		name  string
		query string
	}{
		{"six name prefixes", "SOURCE logGroups(namePrefix: ['a', 'b', 'c', 'd', 'e', 'f'])"},
		{"twenty-one accounts", "SOURCE logGroups(accountIdentifier: ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12', '13', '14', '15', '16', '17', '18', '19', '20', '21'])"},
		{"six tag filters", "SOURCE logGroupTags([{key: \"a\"}, {key: \"b\"}, {key: \"c\"}, {key: \"d\"}, {key: \"e\"}, {key: \"f\"}])"},
		{"six tag values", "SOURCE logGroupTags([{key: \"a\", values: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\"]}])"},
		{"eleven data sources", "SOURCE dataSource(['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11'])"},
	}
	for _, tc := range cases {
		if err := validateQueryPipeline(tc.query); err == nil {
			t.Fatalf("%s: expected MalformedQueryException", tc.name)
		}
	}
}

func TestSourceTooManyGroupsFails(t *testing.T) {
	groups := make([]sourceGroupInfo, maxSourceLogGroups+1)
	for i := range groups {
		groups[i] = sourceGroupInfo{Name: "g" + strconv.Itoa(i)}
	}
	ctx := newTestCtx(nil)
	ctx.listLogGroups = func() ([]sourceGroupInfo, error) {
		return groups, nil
	}
	_, err := executeQueryContext(ctx, `SOURCE logGroups() | fields @message`)
	if err == nil {
		t.Fatal("expected LimitExceededException when selection exceeds the group ceiling")
	}
	ae, ok := err.(*awserrors.AWSError)
	if !ok || ae.Code != "LimitExceededException" {
		t.Fatalf("expected LimitExceededException, got %v", err)
	}
}

// --- subquery determinism / join validation ---

func TestSubqueryFirstFieldIsFirstColumn(t *testing.T) {
	events := testEvents(
		`{"req": "r1", "kind": "b-value"}`,
		`{"req": "r2", "kind": "a-value"}`,
	)
	// Each query runs on its own context: the subquery cache keys on token
	// offsets, which are only unique within one query string.
	rows, err := executeQueryContext(newTestCtx(events), `filter req in (fields kind, req)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("first column (kind) must supply the in-set, got %d rows", len(rows))
	}
	rows, err = executeQueryContext(newTestCtx(events), `filter req in (fields req, kind)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("req-first subquery should match both rows, got %d", len(rows))
	}
}

func TestJoinWithoutSourceRejected(t *testing.T) {
	err := validateQueryPipeline(`join type=inner where a.id=b.id`)
	if err == nil {
		t.Fatal("join without a parenthesised SOURCE must be rejected at compile time")
	}
	ae, ok := err.(*awserrors.AWSError)
	if !ok || ae.Code != "MalformedQueryException" {
		t.Fatalf("expected MalformedQueryException, got %v", err)
	}
}

// --- function arity validation ---

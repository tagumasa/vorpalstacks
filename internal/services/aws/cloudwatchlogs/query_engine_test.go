package cloudwatchlogs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/utils/aws/errors"
)

func testEvents(msgs ...string) []logEventWithContext {
	var evts []logEventWithContext
	for i, m := range msgs {
		evts = append(evts, logEventWithContext{
			timestamp:     int64(1700000000000 + i*60000),
			message:       m,
			ingestionTime: int64(1700000000000 + i*60000),
			logGroup:      "test-group",
			logStream:     "stream-1",
		})
	}
	return evts
}

func rowField(t *testing.T, rows []queryResultRow, idx int, field string) string {
	t.Helper()
	if idx >= len(rows) {
		t.Fatalf("row %d out of range (len=%d)", idx, len(rows))
	}
	return rows[idx].fields[field]
}

// --- compile validation ---

func TestValidateQueryPipelineUnknownCommandOffset(t *testing.T) {
	// The command name also occurs inside an earlier string literal; the
	// reported offset must point at the actual command token, not the
	// literal occurrence.
	q := `parse @message "zz * zzz" as p | bogus p`
	err := validateQueryPipeline(q)
	if err == nil {
		t.Fatal("expected MalformedQueryException for unknown command")
	}
	ae, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	if ae.Code != "MalformedQueryException" {
		t.Fatalf("unexpected code %s", ae.Code)
	}
	rf := ae.RawFields
	qce := rf["queryCompileError"].(map[string]interface{})
	loc := qce["location"].(map[string]interface{})
	want := strings.Index(q, "| bogus") + 2
	if got := loc["startCharOffset"].(int); got != want {
		t.Fatalf("startCharOffset = %d, want %d (actual command position)", got, want)
	}
}

func TestValidateQueryPipelineRules(t *testing.T) {
	cases := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"fields ok", "fields @message | limit 10", false},
		{"where alias", "fields @message | where strcontains(@message, \"x\")", false},
		{"unknown command", "fields @message | frobnicate", true},
		{"dedup then limit ok", "fields @message | dedup @message | limit 5", false},
		{"dedup then filter rejected", "fields @message | dedup @message | filter @message like /x/", true},
		{"sort before stats rejected", "sort @timestamp desc | stats count(*)", true},
		{"sort after stats ok", "stats count(*) by bin(1m) | sort @timestamp desc", false},
		{"pattern after sort rejected", "sort @timestamp desc | pattern @message", true},
		{"pattern before sort ok", "pattern @message | sort @sampleCount asc", false},
		{"source not first rejected", "fields @message | SOURCE logGroups(namePrefix: ['a'])", true},
		{"source first ok", "SOURCE logGroups(namePrefix: ['a']) | fields @message", false},
		{"comments ignored", "fields @message # trailing comment | limit 3", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateQueryPipeline(tc.query)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q", tc.query)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.query, err)
			}
		})
	}
}

func TestValidateQueryPipelineStatsLimit(t *testing.T) {
	var b strings.Builder
	b.WriteString("stats count(*)")
	for i := 0; i < 10; i++ {
		b.WriteString(" | stats count(*)")
	}
	if err := validateQueryPipeline(b.String()); err == nil {
		t.Fatal("expected error for more than 10 stats commands")
	}
}

// --- expression evaluation ---

func TestExpressionFunctions(t *testing.T) {
	cases := []struct {
		expr  string
		want  string
		input string
	}{
		{"toupper(name)", "HELLO", "hello"},
		{"strcontains(name, \"ell\")", "1", "hello"},
		{"strlen(name)", "5", "hello"},
		{"isNumeric(name)", "false", "hello"},
		{"isNumeric(num)", "true", "42"},
		{"coalesce(missing, name)", "hello", "hello"},
		{"substr(name, 1, 2)", "el", "hello"},
		{"abs(-5)", "5", ""},
		{"round(2.567, 2)", "2.57", ""},
		{"isValidIp(ip)", "true", "192.168.1.1"},
		{"isValidIp(name)", "false", "hello"},
		{"jsonArraySize(arr)", "2", `["a","b"]`},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			row := queryResultRow{fields: map[string]string{
				"name": tc.input,
				"num":  "42",
				"ip":   "192.168.1.1",
				"arr":  `["a","b"]`,
			}}
			toks, err := lexQuery(tc.expr)
			if err != nil {
				t.Fatal(err)
			}
			node, err := parseExprTokens(toks)
			if err != nil {
				t.Fatal(err)
			}
			got := asString(node.eval(&row, nil))
			if got != tc.want {
				t.Fatalf("eval(%s) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

func TestExpressionStructureAccess(t *testing.T) {
	row := queryResultRow{fields: map[string]string{
		"json_message": `{"users": [{"action": "PutData"}, {"action": "GetData"}], "status": 200}`,
	}}
	node, err := parseExprTokens(mustLex(t, "json_message.users[1].action"))
	if err != nil {
		t.Fatal(err)
	}
	if got := asString(node.eval(&row, nil)); got != "GetData" {
		t.Fatalf("structure access = %q, want GetData", got)
	}
	node2, err := parseExprTokens(mustLex(t, "json_message.status"))
	if err != nil {
		t.Fatal(err)
	}
	if got := asString(node2.eval(&row, nil)); got != "200" {
		t.Fatalf("map access = %q, want 200", got)
	}
}

func TestExpressionLikeRegexAndIn(t *testing.T) {
	row := queryResultRow{fields: map[string]string{"msg": "Request failed with error 503"}}
	node, err := parseExprTokens(mustLex(t, `msg like /error \d+/`))
	if err != nil {
		t.Fatal(err)
	}
	if !truthy(node.eval(&row, nil)) {
		t.Fatal("regex like should match")
	}
	node2, err := parseExprTokens(mustLex(t, `503 in [500, 502, 503]`))
	if err != nil {
		t.Fatal(err)
	}
	if !truthy(node2.eval(&row, nil)) {
		t.Fatal("in list should match")
	}
	node3, err := parseExprTokens(mustLex(t, `msg not like "Request*"`))
	if err != nil {
		t.Fatal(err)
	}
	if truthy(node3.eval(&row, nil)) {
		t.Fatal("negated glob like should not match")
	}
}

func mustLex(t *testing.T, s string) []token {
	t.Helper()
	toks, err := lexQuery(s)
	if err != nil {
		t.Fatal(err)
	}
	return toks
}

// --- commands ---

func TestFilterWhereEquivalence(t *testing.T) {
	events := testEvents("alpha", "beta", "alphabet")
	filterRows, _ := executeQuery(`filter @message like "alp*"`, events)
	whereRows, _ := executeQuery(`where @message like "alp*"`, events)
	if len(filterRows) != 2 || len(whereRows) != 2 {
		t.Fatalf("filter=%d where=%d, want 2/2", len(filterRows), len(whereRows))
	}
}

func TestFieldsProjectionWithAliasAndArithmetic(t *testing.T) {
	events := testEvents(`{"bytes": 100}`, `{"bytes": 300}`)
	rows, err := executeQueryContext(newTestCtx(events), `fields jsonParse(@message).bytes * 2 as doubled`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows, 0, "doubled") != "200" || rowField(t, rows, 1, "doubled") != "600" {
		t.Fatalf("doubled = %q, %q", rowField(t, rows, 0, "doubled"), rowField(t, rows, 1, "doubled"))
	}
}

func TestStatsAggregationsAndBin(t *testing.T) {
	events := testEvents(
		`{"v": 1}`, `{"v": 2}`, `{"v": 3}`,
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`stats sum(v) as total, avg(v) as mean, count(*) as n, max(v) as hi by bin(1h)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rowField(t, rows, 0, "total") != "6" || rowField(t, rows, 0, "mean") != "2" || rowField(t, rows, 0, "n") != "3" {
		t.Fatalf("aggregates: %+v", rows[0].fields)
	}
	if rowField(t, rows, 0, "hi") != "3" {
		t.Fatalf("max = %q", rowField(t, rows, 0, "hi"))
	}
}

func TestStatsPercentileAndStddev(t *testing.T) {
	events := testEvents(`{"v": 10}`, `{"v": 20}`, `{"v": 30}`, `{"v": 40}`)
	rows, err := executeQueryContext(newTestCtx(events), `stats pct(v, 50) as p50, stddev(v) as sd`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows, 0, "p50") != "20" {
		t.Fatalf("p50 = %q", rowField(t, rows, 0, "p50"))
	}
	if rowField(t, rows, 0, "sd") != "11.180339887498949" {
		t.Fatalf("sd = %q", rowField(t, rows, 0, "sd"))
	}
}

func TestDedupDefaultOrderAndNullRetention(t *testing.T) {
	// Three events sharing server=a; the latest (desc default) must win.
	events := []logEventWithContext{}
	base := int64(1700000000000)
	for i, m := range []string{
		`{"server": "a", "seq": 1}`,
		`{"server": "a", "seq": 2}`,
		`{"server": "b", "seq": 3}`,
		`{"other": 4}`,
	} {
		events = append(events, logEventWithContext{timestamp: base + int64(i)*1000, message: m})
	}
	rows, err := executeQueryContext(newTestCtx(events), `dedup server`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3 (a, b, null-retained)", len(rows))
	}
	seqs := map[string]bool{}
	for _, r := range rows {
		if v, ok := r.fields["seq"]; ok {
			seqs[v] = true
		}
	}
	if !seqs["2"] {
		t.Fatalf("latest a row should be retained, got %+v", rows)
	}
}

func TestFilldownAccumAutoregress(t *testing.T) {
	events := testEvents(`{"h": "", "v": 1}`, `{"h": "x", "v": 2}`, `{"h": "", "v": 3}`)
	rows, err := executeQueryContext(newTestCtx(events), `filldown h | accum v as total | autoregress v p=1-2`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows, 2, "h") != "x" {
		t.Fatalf("filldown: %+v", rows[2].fields)
	}
	if rowField(t, rows, 2, "total") != "6" {
		t.Fatalf("accum total = %q", rowField(t, rows, 2, "total"))
	}
	if rowField(t, rows, 2, "v_p1") != "2" || rowField(t, rows, 2, "v_p2") != "1" {
		t.Fatalf("autoregress lags: %+v", rows[2].fields)
	}
}

func TestUnnestAndExpand(t *testing.T) {
	events := testEvents(`{"events": [{"name": "a"}, {"name": "b"}], "host": "w"}`)
	rows, err := executeQueryContext(newTestCtx(events), `fields jsonParse(@message) as jm | unnest jm.events into ev | display ev.name`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("unnest rows = %d, want 2", len(rows))
	}

	events2 := testEvents(`{"items": ["apple", "banana", "cherry"], "host": "web"}`)
	rows2, err := executeQueryContext(newTestCtx(events2), `fields jsonParse(@message) as jm | expand jm.items`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 3 {
		t.Fatalf("expand rows = %d, want 3", len(rows2))
	}
}

func TestPatternMaskingAndSeverity(t *testing.T) {
	events := testEvents(
		"2023-01-01 19:00:01 [INFO] Calling DynamoDB id 12342342k124-12345",
		"2023-01-01 19:00:02 [INFO] Calling DynamoDB id 324892398123-12345",
		"2023-01-01 19:00:03 [ERROR] disk failure on device sda1",
	)
	rows, err := executeQueryContext(newTestCtx(events), `pattern @message`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("patterns = %d, want 2", len(rows))
	}
	var severityByID = rowField(t, rows, 0, "@severityLabel")
	if severityByID != "INFO" && severityByID != "ERROR" {
		t.Fatalf("severity = %q", severityByID)
	}
	pat := rowField(t, rows, 0, "@pattern")
	if !strings.Contains(pat, "[INFO] Calling DynamoDB") && !strings.Contains(pat, "[ERROR] disk failure") {
		t.Fatalf("pattern = %q", pat)
	}
	if !strings.Contains(pat, "<ID-1>") && !strings.Contains(pat, "<Time-1>") {
		t.Fatalf("dynamic token missing in %q", pat)
	}
	if rowField(t, rows, 0, "@sampleCount") != "2" && rowField(t, rows, 0, "@sampleCount") != "1" {
		t.Fatalf("sampleCount = %q", rowField(t, rows, 0, "@sampleCount"))
	}
}

func TestCountFrequentAndOutlierAndSessionize(t *testing.T) {
	events := testEvents(`{"m": "a", "s": 1}`, `{"m": "a", "s": 2}`, `{"m": "b", "s": 3}`)
	rows, err := executeQueryContext(newTestCtx(events), `countFrequent m`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rowField(t, rows, 0, "_approxcount") != "2" {
		t.Fatalf("countFrequent: %+v", rows)
	}

	eventsOut := testEvents(`{"cpu": 10}`, `{"cpu": 11}`, `{"cpu": 12}`, `{"cpu": 100}`)
	rows2, err := executeQueryContext(newTestCtx(eventsOut), `outlier action=transform param=1.0 cpu`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows2, 3, "cpu") == "100" {
		t.Fatalf("outlier should clamp: %+v", rows2[3].fields)
	}

	base := int64(1700000000000)
	var sessionEvents []logEventWithContext
	for i, m := range []string{
		`{"user": "u1", "n": 1}`,
		`{"user": "u1", "n": 2}`,
		`{"user": "u1", "n": 3}`,
	} {
		sessionEvents = append(sessionEvents, logEventWithContext{
			timestamp: base + int64(i)*10*60*1000, // 10-minute gaps
			message:   m,
		})
	}
	rows3, err := executeQueryContext(newTestCtx(sessionEvents), `sessionize user maxspan 5m as sid`)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, r := range rows3 {
		ids[r.fields["sid"]] = true
	}
	if len(ids) != 3 {
		t.Fatalf("5-minute maxspan should split into 3 sessions, got %d", len(ids))
	}
}

func TestAddTotalsAndFillmissing(t *testing.T) {
	events := testEvents(`{"a": 1, "b": 2}`, `{"a": 3, "b": 4}`)
	rows, err := executeQueryContext(newTestCtx(events), `addtotals a, b`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows, 0, "Total") != "3" || rowField(t, rows, 1, "Total") != "7" {
		t.Fatalf("addtotals: %+v", rows)
	}

	// Two events 10 minutes apart with a 5-minute bin: the middle bin is
	// missing and fillmissing must synthesise it.
	base := int64(1700000000000)
	var binned []logEventWithContext
	for i := 0; i < 2; i++ {
		binned = append(binned, logEventWithContext{timestamp: base + int64(i)*10*60*1000, message: `{"v": 1}`})
	}
	rows2, err := executeQueryContext(newTestCtx(binned),
		`stats sum(v) as s by bin(5m) | fillmissing with 0 for s`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 3 {
		t.Fatalf("fillmissing rows = %d, want 3 (missing middle bin synthesised): %+v", len(rows2), rows2)
	}
}

func TestDiffComparesPreviousWindow(t *testing.T) {
	base := int64(1700000000000)
	var cur []logEventWithContext
	for i, m := range []string{"[INFO] started id 1", "[INFO] started id 2", "[ERROR] crashed hard"} {
		cur = append(cur, logEventWithContext{timestamp: base + int64(i)*1000, message: m})
	}
	var prev []logEventWithContext
	for i, m := range []string{"[INFO] started id 9"} {
		prev = append(prev, logEventWithContext{timestamp: base - 3600000 + int64(i)*1000, message: m})
	}
	ctx := newTestCtx(cur)
	window := base + 3600000
	ctx.startTime = window - 3600000
	ctx.endTime = window
	ctx.fetchEvents = func(groups []string, start, end int64) ([]logEventWithContext, error) {
		if end <= window-3600000+1000 {
			return prev, nil
		}
		return cur, nil
	}
	rows, err := executeQueryContext(ctx, `diff`)
	if err != nil {
		t.Fatal(err)
	}
	// started-pattern diff +1, crashed-pattern is new.
	var startedDiff, crashedCount string
	for _, r := range rows {
		if strings.Contains(r.fields["@pattern"], "started") {
			startedDiff = r.fields["@diffEventCount"]
		}
		if strings.Contains(r.fields["@pattern"], "crashed") {
			crashedCount = r.fields["@differenceDescription"]
		}
	}
	if startedDiff != "1" {
		t.Fatalf("started pattern diff = %q, want 1; rows=%+v", startedDiff, rows)
	}
	if crashedCount != "new" {
		t.Fatalf("crashed pattern should be new, got %q", crashedCount)
	}
}

// --- subquery / join / appendcols ---

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

func TestFunctionArityValidation(t *testing.T) {
	cases := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"if no args", "fields if()", true},
		{"if two args", "fields if(@message, 1)", true},
		{"if three args", "fields if(ispresent(x), 1, 2)", false},
		{"topk no args", "stats topk()", true},
		{"substr one arg", "fields substr(@message)", true},
		{"substr two args", "fields substr(@message, 2)", false},
		{"haversine three args", "fields haversine(1, 2, 3)", true},
		{"coalesce one arg", "fields coalesce(x)", false},
		{"coalesce no args", "fields coalesce()", true},
		{"count no args", "stats count(*)", false},
		{"messageSize no args", "fields messageSize()", true},
		{"messageSize one arg", "fields messageSize(@message) as sz", false},
	}
	for _, tc := range cases {
		err := validateQueryPipeline(tc.query)
		if (err != nil) != tc.wantErr {
			t.Fatalf("%s: err=%v wantErr=%v", tc.name, err, tc.wantErr)
		}
	}
}

func TestCaseBranchLimit(t *testing.T) {
	// case supports up to ten condition/value branches plus an optional
	// trailing default; eleven branches are rejected at compile time.
	build := func(pairs int, withDefault bool) string {
		q := "fields case("
		for i := 0; i < pairs; i++ {
			if i > 0 {
				q += ", "
			}
			q += "false, 0"
		}
		if withDefault {
			q += ", -1"
		}
		return q + ") as c"
	}
	if err := validateQueryPipeline(build(10, true)); err != nil {
		t.Fatalf("ten branches with default should compile: %v", err)
	}
	if err := validateQueryPipeline(build(10, false)); err != nil {
		t.Fatalf("ten branches without default should compile: %v", err)
	}
	if err := validateQueryPipeline(build(11, false)); err == nil {
		t.Fatal("eleven branches should be rejected")
	}
	if err := validateQueryPipeline(build(11, true)); err == nil {
		t.Fatal("eleven branches with default should be rejected")
	}
}

func TestStatsTopkOnEmptyInputNoPanic(t *testing.T) {
	ctx := newTestCtx(nil)
	rows, err := executeQueryContext(ctx, `stats topk(3, @message) as t`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("stats over empty input should emit one row, got %d", len(rows))
	}
}

// --- error offset accuracy ---

// compileErrorOffset extracts the startCharOffset of a MalformedQueryException.
func compileErrorOffset(t *testing.T, err error) int {
	t.Helper()
	ae, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	qce := ae.RawFields["queryCompileError"].(map[string]interface{})
	loc := qce["location"].(map[string]interface{})
	return loc["startCharOffset"].(int)
}

func TestLikeWithoutPatternOffset(t *testing.T) {
	q := "fields @message | filter ip like"
	err := validateQueryPipeline(q)
	if err == nil {
		t.Fatal("expected error for like without a pattern")
	}
	// The offset must fall inside the query string instead of the previous
	// (-1, 0) zero-value token position.
	if got := compileErrorOffset(t, err); got < 0 || got > len(q) {
		t.Fatalf("startCharOffset = %d, want within [0,%d]", got, len(q))
	}
}

func TestStatsLimitOffsetPointsAtEleventh(t *testing.T) {
	q := "fields @message"
	for i := 0; i < 11; i++ {
		q += " | stats count(*) as c" + strconv.Itoa(i)
	}
	err := validateQueryPipeline(q)
	if err == nil {
		t.Fatal("expected error for more than 10 stats commands")
	}
	want := strings.Index(q, "stats count(*) as c10")
	if got := compileErrorOffset(t, err); got != want {
		t.Fatalf("startCharOffset = %d, want %d (11th stats)", got, want)
	}
}

func TestJoinLimitOffsetPointsAtSecond(t *testing.T) {
	q := "fields @message" +
		" | join type=inner where l.a=r.a (SOURCE logGroups(namePrefix: ['x']))" +
		" | join type=inner where l.b=r.b (SOURCE logGroups(namePrefix: ['y']))"
	err := validateQueryPipeline(q)
	if err == nil {
		t.Fatal("expected error for a second join")
	}
	want := strings.LastIndex(q, "join")
	if got := compileErrorOffset(t, err); got != want {
		t.Fatalf("startCharOffset = %d, want %d (second join)", got, want)
	}
}

// --- output fidelity: @log, @ptr, timestamp formatting ---

func TestBuildRowsSystemFields(t *testing.T) {
	events := testEvents(`{"level": "INFO", "@weird": 1}`)
	rows := buildRows(events, "123456789012")
	if len(rows) != 1 {
		t.Fatalf("rows = %d", len(rows))
	}
	cols := rows[0].ordered()
	want := []string{"@timestamp", "@message", "@logStream", "@log", "@ingestionTime", "@ptr", "@@weird", "level"}
	if len(cols) != len(want) {
		t.Fatalf("columns = %v, want %v", cols, want)
	}
	for i := range want {
		if cols[i] != want[i] {
			t.Fatalf("columns = %v, want %v", cols, want)
		}
	}
	if got := rows[0].fields["@log"]; got != "123456789012:test-group" {
		t.Fatalf("@log = %q, want account:group", got)
	}
	if got := rows[0].fields["@@weird"]; got != "1" {
		t.Fatalf("@-prefixed discovered key should render with @@ prefix, got %q", got)
	}
}

func TestDiscoveredFieldCap(t *testing.T) {
	var b strings.Builder
	b.WriteString("{")
	for i := 0; i < maxDiscoveredJSONFields+10; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(`"f` + strconv.Itoa(i) + `":1`)
	}
	b.WriteString("}")
	rows := buildRows(testEvents(b.String()), "")
	discovered := 0
	for _, k := range rows[0].ordered() {
		if !strings.HasPrefix(k, "@") {
			discovered++
		}
	}
	if discovered != maxDiscoveredJSONFields {
		t.Fatalf("discovered = %d, want the documented cap %d", discovered, maxDiscoveredJSONFields)
	}
}

func TestFormatResultTimestamp(t *testing.T) {
	got := formatResultTimestamp("1700000000000")
	if got != "2023-11-14 22:13:20.000" {
		t.Fatalf("formatResultTimestamp = %q", got)
	}
	if got := formatResultTimestamp("not-a-number"); got != "not-a-number" {
		t.Fatalf("non-numeric values must pass through, got %q", got)
	}
}

func TestRowPointerRoundTrip(t *testing.T) {
	// The message deliberately contains the pointer delimiter; the decoder
	// must treat everything after the third delimiter as the message.
	rows := buildRows(testEvents("left|middle|right"), "123456789012")
	ptr := rows[0].fields["@ptr"]
	if ptr == "" {
		t.Fatal("event row should carry a pointer")
	}
	decoded, err := base64.StdEncoding.DecodeString(ptr)
	if err != nil {
		t.Fatal(err)
	}
	parts := splitPointer(string(decoded))
	if len(parts) != 4 || parts[0] != "test-group" || parts[1] != "stream-1" || parts[3] != "left|middle|right" {
		t.Fatalf("pointer = %q", string(decoded))
	}
	// The pointer survives projection, which is how results stay
	// addressable with GetLogRecord after a fields command.
	ctx := newTestCtx(testEvents("hello"))
	projected, err := executeQueryContext(ctx, `fields @message`)
	if err != nil {
		t.Fatal(err)
	}
	if projected[0].fields["@ptr"] == "" {
		t.Fatal("projected event row should keep @ptr")
	}
	// Aggregate rows built by stats carry no pointer.
	statsRows, err := executeQueryContext(newTestCtx(testEvents("hello")), `stats count(*) as c`)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := statsRows[0].fields["@ptr"]; ok {
		t.Fatal("aggregate rows carry no pointer")
	}
}

// --- hashing and time-series functions ---

func TestHashFunctions(t *testing.T) {
	events := testEvents(`{"msg": "hello"}`)
	ctx := newTestCtx(events)
	rows, err := executeQueryContext(ctx, `fields md5(msg) as m, sha256(msg) as s`)
	if err != nil {
		t.Fatal(err)
	}
	// Reference digests of the ASCII string "hello".
	const md5Hello = "5d41402abc4b2a76b9719d911017c592"
	const sha256Hello = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if rows[0].fields["m"] != md5Hello || rows[0].fields["s"] != sha256Hello {
		t.Fatalf("md5=%q sha256=%q", rows[0].fields["m"], rows[0].fields["s"])
	}
}

func TestTimeSeriesFunctions(t *testing.T) {
	events := testEvents(
		`{"n": 1}`,
		`{"n": 3}`,
	)
	ctx := newTestCtx(events)
	// Ungrouped rate over the whole window: sum 4 over a window of 1ms at
	// 1m intervals gives 4/(1/60000).
	rows, err := executeQueryContext(ctx, `stats rate(n, 1m) as r, countOverTime(n) as c, sumOverTime(n) as s`)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].fields["c"] != "2" || rows[0].fields["s"] != "4" {
		t.Fatalf("countOverTime=%q sumOverTime=%q", rows[0].fields["c"], rows[0].fields["s"])
	}
	r, ok := asNumber(rows[0].fields["r"])
	if !ok || r <= 0 {
		t.Fatalf("rate = %q, want a positive per-interval value", rows[0].fields["r"])
	}

	// histogram buckets the values into equal-width ranges; 1 and 3 with
	// two buckets land one per bucket.
	rows, err = executeQueryContext(newTestCtx(events), `stats histogram(n, 2) as h`)
	if err != nil {
		t.Fatal(err)
	}
	var hist map[string]interface{}
	if err := json.Unmarshal([]byte(rows[0].fields["h"]), &hist); err != nil {
		t.Fatalf("histogram is not a JSON map: %q", rows[0].fields["h"])
	}
	if len(hist) != 2 {
		t.Fatalf("histogram buckets = %v, want 2", hist)
	}
	total := 0
	for _, v := range hist {
		f, _ := asNumber(v)
		total += int(f)
	}
	if total != 2 {
		t.Fatalf("histogram total = %d, want 2", total)
	}
}

func TestRateWithBinGrouping(t *testing.T) {
	events := testEvents(`{"n": 2}`, `{"n": 2}`, `{"n": 2}`)
	// testEvents spaces events one minute apart; bin(1m) yields one bin
	// per event, each holding one row with n=2. With 1m intervals the
	// per-bin rate is 2 per interval.
	ctx := newTestCtx(events)
	rows, err := executeQueryContext(ctx, `stats rate(n, 1m) as r by bin(1m)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("bins = %d, want 3", len(rows))
	}
	for i := range rows {
		if rows[i].fields["r"] != "2" {
			t.Fatalf("bin %d rate = %q, want 2", i, rows[i].fields["r"])
		}
	}
}

// --- timestamp-typed function results ---

func TestTimestampFunctionResults(t *testing.T) {
	// The documented result type of fromMillis, datefloor and dateceil is
	// Timestamp: values render in the result timestamp form while staying
	// numeric inside further expressions.
	rows, err := executeQueryContext(newTestCtx(testEvents("x")),
		`fields fromMillis(1700000000000) as fm, datefloor(fromMillis(1700000000000), 1h) as fl, dateceil(fromMillis(1700000000000), 1h) as fc`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows, 0, "fm"); got != "2023-11-14 22:13:20.000" {
		t.Fatalf("fromMillis rendering = %q", got)
	}
	if got := rowField(t, rows, 0, "fl"); got != "2023-11-14 22:00:00.000" {
		t.Fatalf("datefloor rendering = %q", got)
	}
	if got := rowField(t, rows, 0, "fc"); got != "2023-11-14 23:00:00.000" {
		t.Fatalf("dateceil rendering = %q", got)
	}

	// A stored timestamp value round-trips: toMillis parses the stored
	// rendering back to epoch milliseconds.
	rows, err = executeQueryContext(newTestCtx(testEvents("x")),
		`fields fromMillis(1700000000000) as t | fields toMillis(t) as ms`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows, 0, "ms"); got != "1700000000000" {
		t.Fatalf("toMillis over stored timestamp = %q", got)
	}

	// Comparisons between timestamp-typed values keep the numeric
	// interpretation.
	rows, err = executeQueryContext(newTestCtx(testEvents("x")),
		`filter fromMillis(1700000000000) > fromMillis(1600000000000) | stats count(*) as n`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows, 0, "n"); got != "1" {
		t.Fatalf("timestamp comparison lost the event: n = %q", got)
	}
}

func TestTimestampAggregationRendering(t *testing.T) {
	events := testEvents("a", "b") // 22:13:20 and 22:14:20
	// The passthrough aggregation functions keep the Timestamp rendering
	// of @timestamp.
	rows, err := executeQueryContext(newTestCtx(events),
		`stats earliest(@timestamp) as e, latest(@timestamp) as l, max(@timestamp) as m`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows, 0, "e"); got != "2023-11-14 22:13:20.000" {
		t.Fatalf("earliest(@timestamp) = %q", got)
	}
	if got := rowField(t, rows, 0, "l"); got != "2023-11-14 22:14:20.000" {
		t.Fatalf("latest(@timestamp) = %q", got)
	}
	if got := rowField(t, rows, 0, "m"); got != "2023-11-14 22:14:20.000" {
		t.Fatalf("max(@timestamp) = %q", got)
	}

	// Bin group keys render as timestamps.
	rows, err = executeQueryContext(newTestCtx(events), `stats count(*) as n by bin(1m)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("bins = %d, want 2", len(rows))
	}
	if got := rowField(t, rows, 0, "@bin"); got != "2023-11-14 22:13:00.000" {
		t.Fatalf("first bin key = %q", got)
	}
	if got := rowField(t, rows, 1, "@bin"); got != "2023-11-14 22:14:00.000" {
		t.Fatalf("second bin key = %q", got)
	}
}

func TestFillmissingBinRendering(t *testing.T) {
	// Two events 10 minutes apart with a 5-minute bin: the synthesised
	// middle bin carries the same timestamp rendering as the real bins.
	base := int64(1700000000000) // 22:13:20
	var binned []logEventWithContext
	for i := 0; i < 2; i++ {
		binned = append(binned, logEventWithContext{timestamp: base + int64(i)*10*60*1000, message: `{"v": 1}`})
	}
	rows, err := executeQueryContext(newTestCtx(binned),
		`stats sum(v) as s by bin(5m) | fillmissing with 0 for s`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("fillmissing rows = %d, want 3", len(rows))
	}
	want := []string{
		"2023-11-14 22:10:00.000",
		"2023-11-14 22:15:00.000",
		"2023-11-14 22:20:00.000",
	}
	for i, w := range want {
		if got := rowField(t, rows, i, "@bin"); got != w {
			t.Fatalf("bin %d = %q, want %q", i, got, w)
		}
	}
	if got := rowField(t, rows, 1, "s"); got != "0" {
		t.Fatalf("synthesised bin fill = %q", got)
	}
}

// --- lookup / cidrlookup ---

// lookupTestCtx returns a context whose lookup tables are backed by CSV
// fixtures.
func lookupTestCtx(events []logEventWithContext, tables map[string]string) *execContext {
	ctx := newTestCtx(events)
	ctx.getLookupTable = func(name string) (*parsedLookupTable, error) {
		body, ok := tables[name]
		if !ok {
			return nil, fmt.Errorf("lookup table %s not found", name)
		}
		columns, records, err := parseLookupCSV(body)
		if err != nil {
			return nil, err
		}
		return newParsedLookupTable(columns, records), nil
	}
	return ctx
}

func TestLookupCommand(t *testing.T) {
	tables := map[string]string{
		"users": "id,name,department\nu1,Alice,Engineering\nu2,Bob,Support",
	}
	events := testEvents(
		`{"user_id": "u1", "hint": "x"}`,
		`{"user_id": "u2", "hint": "y"}`,
		`{"user_id": "u3", "hint": "z"}`,
	)
	ctx := lookupTestCtx(events, tables)
	rows, err := executeQueryContext(ctx, `lookup users id as user_id OUTPUT name, department`)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].fields["name"] != "Alice" || rows[0].fields["department"] != "Engineering" {
		t.Fatalf("row0 enrichment: %+v", rows[0].fields)
	}
	if rows[1].fields["name"] != "Bob" {
		t.Fatalf("row1 enrichment: %+v", rows[1].fields)
	}
	// OUTPUT nulls the output fields when no row matches.
	if rows[2].fields["name"] != "" {
		t.Fatalf("unmatched row should carry null outputs, got %+v", rows[2].fields)
	}

	// OUTPUTNEW keeps existing values.
	events = testEvents(
		`{"user_id": "u1", "name": "Preserved"}`,
		`{"user_id": "u1", "hint": "fillme"}`,
	)
	ctx = lookupTestCtx(events, tables)
	rows, err = executeQueryContext(ctx, `lookup users id as user_id OUTPUTNEW name, department`)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].fields["name"] != "Preserved" {
		t.Fatalf("OUTPUTNEW must keep existing values, got %q", rows[0].fields["name"])
	}
	if rows[1].fields["name"] != "Alice" {
		t.Fatalf("OUTPUTNEW fills missing values, got %q", rows[1].fields["name"])
	}
}

func TestLookupMissingColumnFailsQuery(t *testing.T) {
	tables := map[string]string{"users": "id,name\nu1,Alice"}
	ctx := lookupTestCtx(testEvents(`{"user_id": "u1"}`), tables)
	_, err := executeQueryContext(ctx, `lookup users id as user_id OUTPUT nonexistent`)
	if err == nil {
		t.Fatal("lookup against a missing column must fail the query")
	}
}

func TestLookupMissingTableFailsQuery(t *testing.T) {
	ctx := lookupTestCtx(testEvents(`{"user_id": "u1"}`), map[string]string{})
	_, err := executeQueryContext(ctx, `lookup nope id as user_id OUTPUT name`)
	if err == nil {
		t.Fatal("lookup against a missing table must fail the query")
	}
}

func TestCidrLookupCommand(t *testing.T) {
	tables := map[string]string{
		"nets": "cidr,region,owner\n10.0.0.0/8,corp,it\n192.168.1.0/24,office,netops",
	}
	events := testEvents(
		`{"ip": "10.1.2.3"}`,
		`{"ip": "192.168.1.55"}`,
		`{"ip": "203.0.113.9"}`,
	)
	ctx := lookupTestCtx(events, tables)
	rows, err := executeQueryContext(ctx, `cidrlookup nets ip as cidr OUTPUT region, owner`)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].fields["region"] != "corp" || rows[1].fields["region"] != "office" {
		t.Fatalf("cidr enrichment: %+v / %+v", rows[0].fields, rows[1].fields)
	}
	if rows[2].fields["region"] != "" {
		t.Fatalf("unmatched IP should carry null region, got %q", rows[2].fields["region"])
	}

	// Output aliases write to the event field name.
	ctx = lookupTestCtx(events, tables)
	rows, err = executeQueryContext(ctx, `cidrlookup nets ip as cidr OUTPUT region as net_region`)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].fields["net_region"] != "corp" {
		t.Fatalf("aliased output: %+v", rows[0].fields)
	}
}

func TestRowsToCSVUsesColumnOrder(t *testing.T) {
	rows := buildRows(testEvents(`{"b": 1, "a": 2}`), "acct")
	csvOut := rowsToCSV(rows)
	if !strings.HasPrefix(csvOut, "@timestamp,@message,@logStream,@log,@ingestionTime,@ptr,a,b") {
		t.Fatalf("rowsToCSV header = %q", csvOut)
	}
}

func TestResolveLookupTableBodyFromQueryId(t *testing.T) {
	s := &LogsService{}
	rows := buildRows(testEvents(`{"a": 1}`), "acct")
	s.queries.Store("q-done", &queryState{status: "Complete", results: rows})
	s.queries.Store("q-running", &queryState{status: "Running", results: nil})

	body, err := s.resolveLookupTableBody(&LookupTableInput{QueryId: "q-done"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(body, "@timestamp,@message,@logStream,@log,@ingestionTime,@ptr,a") {
		t.Fatalf("ingested body header = %q", body)
	}

	if _, err := s.resolveLookupTableBody(&LookupTableInput{QueryId: "q-running"}); err == nil {
		t.Fatal("a running query must not populate a lookup table")
	}
	if _, err := s.resolveLookupTableBody(&LookupTableInput{QueryId: "nope"}); err == nil {
		t.Fatal("an unknown query id must be rejected")
	}
	if _, err := s.resolveLookupTableBody(&LookupTableInput{}); err == nil {
		t.Fatal("neither tableBody nor queryId must be rejected")
	}
	if _, err := s.resolveLookupTableBody(&LookupTableInput{TableBody: "a\n1\n", QueryId: "q-done"}); err == nil {
		t.Fatal("tableBody and queryId together must be rejected")
	}
}

func TestOffsetsCountCharactersNotBytes(t *testing.T) {
	// The multibyte string literals shift byte positions ahead of character
	// positions; the reported offsets must stay in character units.
	q := `fields @message, "日本語テキスト", "別の文字列" | frobnicate x`
	err := validateQueryPipeline(q)
	if err == nil {
		t.Fatal("expected MalformedQueryException for unknown command")
	}
	byteIdx := strings.Index(q, "| frobnicate") + 2
	want := utf8.RuneCountInString(q[:byteIdx])
	if got := compileErrorOffset(t, err); got != want {
		t.Fatalf("startCharOffset = %d, want %d (character position)", got, want)
	}
}

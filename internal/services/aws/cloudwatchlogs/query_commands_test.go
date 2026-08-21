package cloudwatchlogs

import (
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

// Tests for the pipeline command implementations.

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

	// Two-minute gaps discriminate the maxspan duration: 5m keeps all
	// three events in one session, 1m splits each into its own.
	base := int64(1700000000000)
	var sessionEvents []logEventWithContext
	for i, m := range []string{
		`{"user": "u1", "n": 1}`,
		`{"user": "u1", "n": 2}`,
		`{"user": "u1", "n": 3}`,
	} {
		sessionEvents = append(sessionEvents, logEventWithContext{
			timestamp: base + int64(i)*2*60*1000, // 2-minute gaps
			message:   m,
		})
	}
	for _, tc := range []struct {
		query string
		want  int
	}{{`sessionize user maxspan 5m as sid`, 1}, {`sessionize user maxspan 1m as sid`, 3}} {
		rows3, err := executeQueryContext(newTestCtx(sessionEvents), tc.query)
		if err != nil {
			t.Fatalf("%s: %v", tc.query, err)
		}
		ids := map[string]bool{}
		for _, r := range rows3 {
			ids[r.fields["sid"]] = true
		}
		if len(ids) != tc.want {
			t.Fatalf("%s: session count = %d, want %d", tc.query, len(ids), tc.want)
		}
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
func TestStatsBinOffset(t *testing.T) {
	// Two events one second either side of the shifted boundary merge into
	// one bucket: without the offset they land in separate 1m buckets.
	events := []logEventWithContext{
		{timestamp: 119000, message: `{"sev": "a"}`},
		{timestamp: 121000, message: `{"sev": "a"}`},
	}
	// Documented spelling bin(1m) offset 30s; a bare-millisecond variant
	// must produce the same bucket.
	for _, q := range []string{
		`stats count(*) as n by bin(1m) offset 30s`,
		`stats count(*) as n by bin(1m) offset 30000`,
	} {
		rows, err := executeQueryContext(newTestCtx(events), q)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		if len(rows) != 1 || rowField(t, rows, 0, "n") != "2" {
			t.Fatalf("%s: rows = %+v, want a single merged bucket of 2", q, rows)
		}
		if rowField(t, rows, 0, "@bin") != "1970-01-01 00:01:30.000" {
			t.Fatalf("%s: @bin = %q, want 00:01:30.000 (boundary shifted by 30s)", q, rowField(t, rows, 0, "@bin"))
		}
	}
}

func TestStatsOffsetWithoutBinRejected(t *testing.T) {
	events := testEvents(`{"sev": "a"}`)
	_, err := executeQueryContext(newTestCtx(events), `stats count(*) as n by sev offset 30s`)
	if err == nil {
		t.Fatal("offset without a bin grouping should be rejected")
	}
	ae, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	if ae.Code != "MalformedQueryException" {
		t.Fatalf("unexpected code %s", ae.Code)
	}
}

func TestParseGlobExtraction(t *testing.T) {
	// Glob mode: each * captures one named field; non-matching rows are
	// kept but leave the captured fields unset.
	events := testEvents("GET /index.html 200", "POST /other", "GET /missing 404")
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message "GET * *" as path, status`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3 (parse keeps unmatched rows)", len(rows))
	}
	if rowField(t, rows, 0, "path") != "/index.html" || rowField(t, rows, 0, "status") != "200" {
		t.Fatalf("captures: %+v", rows[0].fields)
	}
	if rowField(t, rows, 1, "path") != "" || rowField(t, rows, 1, "status") != "" {
		t.Fatalf("unmatched row should have no captures: %+v", rows[1].fields)
	}
	if rowField(t, rows, 2, "path") != "/missing" || rowField(t, rows, 2, "status") != "404" {
		t.Fatalf("captures: %+v", rows[2].fields)
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

func TestLogCompareTimeshift(t *testing.T) {
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
	rows, err := executeQueryContext(ctx, `logcompare timeshift 1h`)
	if err != nil {
		t.Fatal(err)
	}
	// The timeshift window covers the same previous hour the diff test
	// uses, so the expectations match: started-pattern diff +1, crashed
	// pattern is new.
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

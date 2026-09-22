package cloudwatchlogs

import (
	"strconv"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
)

// Tests for the pipeline command implementations.

func TestFilterWhereEquivalence(t *testing.T) {
	events := testEvents("alpha", "beta", "alphabet")
	run := func(query string) []queryResultRow {
		rows, err := executeQueryContext(newTestCtx(events), query)
		if err != nil {
			t.Fatal(err)
		}
		return rows
	}
	filterRows := run(`filter @message like "alp"`)
	whereRows := run(`where @message like "alp"`)
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

// Expression columns carry the expression's written form when the query
// gives no alias: a postfix attribute access names its column with its
// source form (not the empty string), an aggregation's auto-label
// renders its argument, and the bin group key equals the emitted @bin
// column's name — a subsequent stats command can address the bucket.
func TestExpressionColumnNames(t *testing.T) {
	events := testEvents(`{"status": "ok"}`, `{"status": "err"}`, `{"status": "err"}`)

	rows, err := executeQueryContext(newTestCtx(events), `fields jsonParse(@message).status`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows, 0, "jsonParse(@message).status"); got != "ok" {
		t.Fatalf("postfix column value = %q", got)
	}
	if _, ok := rows[0].fields[""]; ok {
		t.Fatalf("columns = %v, want no unnamed column", rows[0].columns)
	}

	rows, err = executeQueryContext(newTestCtx(events), `stats count(*) as n, max(jsonParse(@message).status)`)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := rows[0].fields["max(jsonParse(@message).status)"]; !ok {
		t.Fatalf("columns = %v, want the aggregation's auto-label to render its argument", rows[0].columns)
	}

	// The bin group key carries the emitted column's name, so the second
	// stats command's visibility set holds @bin and the grouping reads it.
	rows, err = executeQueryContext(newTestCtx(events),
		`stats count(*) as n by bin(1h) | stats sum(n) as total by @bin`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rowField(t, rows, 0, "total") != "3" {
		t.Fatalf("binned follow-up stats rows = %+v", rows)
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
	// The empty-string events share one value, so they deduplicate to a
	// single row; the event without the field at all is null and retained.
	events := []logEventWithContext{}
	base := int64(1700000000000)
	for i, m := range []string{
		`{"server": "a", "seq": 1}`,
		`{"server": "a", "seq": 2}`,
		`{"server": "b", "seq": 3}`,
		`{"other": 4}`,
		`{"server": "", "seq": 5}`,
		`{"server": "", "seq": 6}`,
	} {
		events = append(events, logEventWithContext{timestamp: base + int64(i)*1000, message: m})
	}
	rows, err := executeQueryContext(newTestCtx(events), `dedup server`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("rows = %d, want 4 (a, b, null-retained, one empty-string)", len(rows))
	}
	empty := 0
	seqs := map[string]bool{}
	for _, r := range rows {
		if v, ok := r.fields["seq"]; ok {
			seqs[v] = true
		}
		if v, ok := r.fields["server"]; ok && v == "" {
			empty++
		}
	}
	if !seqs["2"] {
		t.Fatalf("latest a row should be retained, got %+v", rows)
	}
	if !seqs["6"] || empty != 1 {
		t.Fatalf("empty-string rows should deduplicate to the latest one, got %d rows, seqs %+v", empty, seqs)
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

func TestParseRegexNamedCaptureGroups(t *testing.T) {
	// Regex mode: parse {{field}} /{{regex}}/ with no as clause — the
	// named capture groups (?<name>...) define the extracted fields, in
	// pattern order. Non-matching rows are kept without the fields.
	events := testEvents(
		"eni-0abc123 attached",
		"user=alice, method:GET, latency := 12",
		"user=bob, method:POST, latency := 34",
		"vol-0789 attached",
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message /user=(?<user>.*?), method:(?<method>.*?), latency := (?<latency>.*)/`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("rows = %d, want 4 (parse keeps unmatched rows)", len(rows))
	}
	for _, row := range []int{0, 3} {
		if rowField(t, rows, row, "user") != "" || rowField(t, rows, row, "method") != "" {
			t.Fatalf("unmatched row %d should have no captures: %+v", row, rows[row].fields)
		}
	}
	if rowField(t, rows, 1, "user") != "alice" || rowField(t, rows, 1, "method") != "GET" || rowField(t, rows, 1, "latency") != "12" {
		t.Fatalf("captures: %+v", rows[1].fields)
	}
	if rowField(t, rows, 2, "user") != "bob" || rowField(t, rows, 2, "method") != "POST" || rowField(t, rows, 2, "latency") != "34" {
		t.Fatalf("captures: %+v", rows[2].fields)
	}
}

func TestParseRegexNumberedGroupsWithAs(t *testing.T) {
	// The numbered-group form keeps its as clause; the regex body reaches
	// the compiler verbatim, so character classes and escapes work.
	events := testEvents("ip 10.0.0.1 src", "ip 192.168.1.44 src", "no address here")
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message /(\d+\.\d+\.\d+\.\d+)/ as ip`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3 (parse keeps unmatched rows)", len(rows))
	}
	if rowField(t, rows, 0, "ip") != "10.0.0.1" {
		t.Fatalf("ip = %q, want 10.0.0.1: %+v", rowField(t, rows, 0, "ip"), rows[0].fields)
	}
	if rowField(t, rows, 1, "ip") != "192.168.1.44" {
		t.Fatalf("ip = %q, want 192.168.1.44", rowField(t, rows, 1, "ip"))
	}
	if rowField(t, rows, 2, "ip") != "" {
		t.Fatalf("unmatched row should have no capture: %+v", rows[2].fields)
	}
}

func TestParseRegexMultiMatch(t *testing.T) {
	// The multi-match mode: "Add the keyword multi after the regex
	// pattern" — every match emits its own row, and an event with no match
	// keeps one field-less row.
	events := testEvents(
		"ip 10.0.0.1 then 10.0.0.2",
		"no address here",
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message /(\d+\.\d+\.\d+\.\d+)/ as ip_addr multi`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3 (two matches of the first event plus the unmatched event): %+v", len(rows), rows)
	}
	if rowField(t, rows, 0, "ip_addr") != "10.0.0.1" || rowField(t, rows, 1, "ip_addr") != "10.0.0.2" {
		t.Fatalf("each match must carry its own row: %+v %+v", rows[0].fields, rows[1].fields)
	}
	if rowField(t, rows, 2, "ip_addr") != "" {
		t.Fatalf("unmatched event must keep one field-less row: %+v", rows[2].fields)
	}

	// The named-group form rides multi the same way, and a follow-on
	// stats aggregation counts the multiplied rows per extracted value —
	// the unmatched event's kept row has no ip and is filtered out
	// before the grouping.
	rows2, err := executeQueryContext(newTestCtx(events),
		`parse @message /(?<ip>\d+\.\d+\.\d+\.\d+)/ multi | filter ispresent(ip) | stats count(*) as n by ip`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 2 || rowField(t, rows2, 0, "n") != "1" || rowField(t, rows2, 1, "n") != "1" {
		t.Fatalf("stats over multiplied rows = %+v, want one row per extracted address", rows2)
	}
}

func TestParsePatternFirstDefaultsToMessage(t *testing.T) {
	// The pattern-first form omits the source field: "If fieldName is
	// omitted, @message is used by default." Both the glob and the regex
	// mode ride the default source.
	events := testEvents("user=alice, latency := 12", "user=bob, latency := 34")
	rows, err := executeQueryContext(newTestCtx(events),
		`parse "user=*, latency := *" as user, latency`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
	if rowField(t, rows, 0, "user") != "alice" || rowField(t, rows, 0, "latency") != "12" {
		t.Fatalf("glob captures: %+v", rows[0].fields)
	}
	if rowField(t, rows, 1, "user") != "bob" || rowField(t, rows, 1, "latency") != "34" {
		t.Fatalf("glob captures: %+v", rows[1].fields)
	}
	rows2, err := executeQueryContext(newTestCtx(events),
		`parse /user=(?<user>.*?), latency := (?<latency>\d+)/`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows2, 0, "user") != "alice" || rowField(t, rows2, 1, "latency") != "34" {
		t.Fatalf("named-group captures on the default source: %+v %+v", rows2[0].fields, rows2[1].fields)
	}
}

func TestParseLogfmtMode(t *testing.T) {
	// Logfmt mode: the line's space-separated key=value pairs become a map
	// under the single alias — "The result is a map that you access with
	// dot notation (for example, lf.level, lf.msg)" — with quoted values
	// unescaped.
	events := testEvents(
		`level=error msg="connection failed" duration=42 host=web-1`,
		`level=info msg=started duration=1 host=web-2`,
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message logfmt as lf | filter lf.level = "error" | display lf.msg, lf.duration`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want the single error-level event", len(rows))
	}
	if rowField(t, rows, 0, "lf.msg") != "connection failed" || rowField(t, rows, 0, "lf.duration") != "42" {
		t.Fatalf("logfmt map fields: %+v", rows[0].fields)
	}

	// The omitted-fieldName form rides the default @message source, and a
	// grouped aggregation addresses map keys through the same dot notation.
	rows2, err := executeQueryContext(newTestCtx(events),
		`parse logfmt as lf | stats count(*) as n by lf.host`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 2 || rowField(t, rows2, 0, "n") != "1" || rowField(t, rows2, 1, "n") != "1" {
		t.Fatalf("grouping by a logfmt map key = %+v, want one row per host", rows2)
	}
	hosts := rowField(t, rows2, 0, "lf.host") + "," + rowField(t, rows2, 1, "lf.host")
	if hosts != "web-1,web-2" && hosts != "web-2,web-1" {
		t.Fatalf("grouped hosts = %q, want web-1 and web-2", hosts)
	}

	// The mode's documented form carries exactly one alias.
	if _, err := executeQueryContext(newTestCtx(events), `parse @message logfmt as a, b`); err == nil {
		t.Fatal("logfmt mode with two field names should be rejected")
	}
	if _, err := executeQueryContext(newTestCtx(events), `parse @message csv`); err == nil {
		t.Fatal("csv mode without an as clause should be rejected")
	}
}

func TestParseCsvMode(t *testing.T) {
	// CSV mode: "Each comma-separated value is assigned to the
	// corresponding alias."
	events := testEvents(
		"web-1,GET,/api/users,200,15",
		"web-2,POST,/api/orders,500,42",
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message csv as host, method, path, status, duration | filter status = "500" | display host, path`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want the single 500 event", len(rows))
	}
	if rowField(t, rows, 0, "host") != "web-2" || rowField(t, rows, 0, "path") != "/api/orders" {
		t.Fatalf("csv fields: %+v", rows[0].fields)
	}

	// A record with fewer columns than aliases leaves the missing aliases
	// unset, keeping the event in the results.
	rows2, err := executeQueryContext(newTestCtx(testEvents("web-1,GET")),
		`parse @message csv as a, b, c`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 1 {
		t.Fatalf("rows = %d, want the short record's kept row", len(rows2))
	}
	if rowField(t, rows2, 0, "a") != "web-1" || rowField(t, rows2, 0, "b") != "GET" || rowField(t, rows2, 0, "c") != "" {
		t.Fatalf("short-record assignment: %+v", rows2[0].fields)
	}
}

func TestJSONCommandChainedExtraction(t *testing.T) {
	// The json command's documented form: json field={{fieldName}}
	// "{{key.subkey}}" as {{alias}} — "explicit chained JSON extraction
	// from a previously parsed object field". A row whose source or key is
	// absent keeps its fields without the alias.
	events := testEvents(
		`{"user": {"name": "alice"}, "code": 200}`,
		`no structure here`,
		`{"user": {"title": "admin"}}`,
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`parse @message /(?<payload>\{.*\})/ | json field=payload "user.name" as username | display username`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3 (absent keys keep their rows)", len(rows))
	}
	if rowField(t, rows, 0, "username") != "alice" {
		t.Fatalf("chained extraction: %+v", rows[0].fields)
	}
	if rowField(t, rows, 1, "username") != "" || rowField(t, rows, 2, "username") != "" {
		t.Fatalf("rows without the key keep no alias: %+v %+v", rows[1].fields, rows[2].fields)
	}

	// The discovered-field form: a structured top-level key feeds the same
	// key path into an aggregation.
	events2 := testEvents(`{"requestContext": {"identity": {"sourceIp": "10.1.2.3"}}}`)
	rows2, err := executeQueryContext(newTestCtx(events2),
		`json field=requestContext "identity.sourceIp" as caller_ip | stats count(*) as n by caller_ip`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 1 || rowField(t, rows2, 0, "n") != "1" || rowField(t, rows2, 0, "caller_ip") != "10.1.2.3" {
		t.Fatalf("discovered-field extraction: %+v", rows2)
	}

	// The documented form's members are required together.
	for _, q := range []string{
		`json field=payload`,
		`json @message "a.b" as x`,
		`json field=payload "a.b"`,
	} {
		if _, err := executeQueryContext(newTestCtx(events), q); err == nil {
			t.Fatalf("query %q should be rejected", q)
		}
	}
}

func TestEstimateCommand(t *testing.T) {
	// The estimate command reports "the estimated bytes that the query
	// would scan over the selected log groups and time range, without
	// running the query" — the source events' byte volume, which a prior
	// filter never shrinks.
	events := testEvents("error one", "fine two", "error three")
	var want int64
	for _, e := range events {
		want += int64(len(e.message))
	}
	rows, err := executeQueryContext(newTestCtx(events),
		`fields @message | filter @message like /error/ | estimate`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want the single estimate row", len(rows))
	}
	if got := rowField(t, rows, 0, "estimate"); got != strconv.FormatInt(want, 10) {
		t.Fatalf("estimate = %q, want %d (every event of the window, pre-filter)", got, want)
	}

	// "The estimate command must be the last command in the query", and it
	// takes no arguments.
	if _, err := executeQueryContext(newTestCtx(events), `estimate | limit 5`); err == nil {
		t.Fatal("a command after estimate should be rejected")
	}
	if _, err := executeQueryContext(newTestCtx(events), `estimate extra`); err == nil {
		t.Fatal("estimate with arguments should be rejected")
	}
}

func TestFieldsUnionAcrossCommands(t *testing.T) {
	// "If your query contains multiple fields commands and doesn't include
	// a display command, the results display all of the fields that are
	// specified in the fields commands."
	events := testEvents(`{"a": "1", "b": "2", "c": "3"}`)
	rows, err := executeQueryContext(newTestCtx(events), `fields a | fields b`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rowField(t, rows, 0, "a") != "1" || rowField(t, rows, 0, "b") != "2" {
		t.Fatalf("union projection: %+v", rows[0].fields)
	}
	if _, present := rows[0].fields["c"]; present {
		t.Fatalf("unlisted field c must stay projected out: %+v", rows[0].fields)
	}
}

func TestChainedDisplayShowsLastFieldsFromFullRecord(t *testing.T) {
	// "If you use display more than once in a query, the query results
	// show the field specified in the last occurrence of display command
	// being used" — resolved against the full record, not the previous
	// display's narrowing.
	events := testEvents(`{"a": "1", "b": "2"}`)
	rows, err := executeQueryContext(newTestCtx(events), `display a | display b`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rowField(t, rows, 0, "b") != "2" {
		t.Fatalf("last display's field: %+v", rows[0].fields)
	}
	if _, present := rows[0].fields["a"]; present {
		t.Fatalf("an earlier display's field must not survive: %+v", rows[0].fields)
	}

	// display overrides an accumulated fields set the same way, resolving
	// its field from the full record.
	rows2, err := executeQueryContext(newTestCtx(events), `fields a | display b`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 1 || rowField(t, rows2, 0, "b") != "2" {
		t.Fatalf("display over fields: %+v", rows2[0].fields)
	}
	if _, present := rows2[0].fields["a"]; present {
		t.Fatalf("display replaces the fields projection: %+v", rows2[0].fields)
	}
}

func TestSortNaturalOrdering(t *testing.T) {
	// The sort command's documented natural ordering: the worked
	// example's full ascending sequence — Unicode-ordered symbols, then
	// alphanumeric values starting with numbers (numeric chunks by length
	// first then value), then letter-led values, then number values with
	// the leading-zero tie broken lexically (0, 01, 1).
	want := []string{
		"!", "#", "*%04", "0#", "5A", "111A", "2345_",
		"@", "@_", "A", "A9876fghj", "a12345hfh",
		"0", "01", "1", "2", "3",
	}
	events := testEvents("3", "A", "01", "111A", "@_", "1", "!", "0#", "2", "5A",
		"@", "0", "a12345hfh", "A9876fghj", "2345_", "#", "*%04")
	rows, err := executeQueryContext(newTestCtx(events), `sort @message asc`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != len(want) {
		t.Fatalf("rows = %d, want %d", len(rows), len(want))
	}
	for i, w := range want {
		if got := rowField(t, rows, i, "@message"); got != w {
			t.Fatalf("ascending position %d = %q, want %q (full order: %v)",
				i, got, w, messageColumn(rows))
		}
	}

	// "If you sort in descending order, the sort results are the
	// reverse."
	rows2, err := executeQueryContext(newTestCtx(events), `sort @message desc`)
	if err != nil {
		t.Fatal(err)
	}
	for i, w := range want {
		if got := rowField(t, rows2, len(want)-1-i, "@message"); got != w {
			t.Fatalf("descending position %d = %q, want %q", i, got, w)
		}
	}
}

func messageColumn(rows []queryResultRow) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.fields["@message"]
	}
	return out
}

func TestInWithArrayValuedField(t *testing.T) {
	// "You can use the keyword in to test for set membership and check
	// for elements in an array. To check for elements in an array, put
	// the array after in."
	events := testEvents(
		`{"action": "read", "tags": ["red", "blue"]}`,
		`{"action": "write", "tags": ["green"]}`,
	)
	rows, err := executeQueryContext(newTestCtx(events), `filter "red" in tags`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rowField(t, rows, 0, "action") != "read" {
		t.Fatalf("array membership = %d rows, want the read event", len(rows))
	}
	rows2, err := executeQueryContext(newTestCtx(events), `filter "red" not in tags`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 1 || rowField(t, rows2, 0, "action") != "write" {
		t.Fatalf("negated array membership = %d rows, want the write event", len(rows2))
	}
	// The literal-list form keeps its semantics alongside the array form.
	rows3, err := executeQueryContext(newTestCtx(events), `filter action in ["read"]`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows3) != 1 || rowField(t, rows3, 0, "action") != "read" {
		t.Fatalf("literal-list membership = %d rows, want the read event", len(rows3))
	}
	// A boolean conjunction after the array operand keeps binding the
	// enclosing expression.
	rows4, err := executeQueryContext(newTestCtx(events), `filter "red" in tags and action = "read"`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows4) != 1 {
		t.Fatalf("conjunction after the array operand = %d rows, want 1", len(rows4))
	}
}

func TestDotNotationTreatsJSONStringFieldsAsOpaqueLeaves(t *testing.T) {
	// "Dot notation traverses only fields that are stored as structurally
	// nested JSON objects at ingest time. If a field's value is a
	// JSON-encoded string (a string whose content happens to be valid
	// JSON), dot notation treats it as an opaque leaf value and does not
	// access sub-fields within it ... To access sub-fields within a
	// JSON-encoded string, use jsonParse."
	events := testEvents(`{"api": {"request": {"data": "{\"startTime\":1234,\"endTime\":5678}"}}}`)
	rows, err := executeQueryContext(newTestCtx(events), `filter ispresent(api.request.data.startTime)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("dot notation must not traverse a JSON-encoded string value: %d rows", len(rows))
	}

	// The string itself is the leaf value.
	rows2, err := executeQueryContext(newTestCtx(events), `fields api.request.data`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows2, 0, "api.request.data"); got != `{"startTime":1234,"endTime":5678}` {
		t.Fatalf("leaf value = %q", got)
	}

	// jsonParse is the documented way into the encoded string.
	rows3, err := executeQueryContext(newTestCtx(events), `fields jsonParse(api.request.data).startTime as st`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows3, 0, "st"); got != "1234" {
		t.Fatalf("jsonParse traversal = %q, want 1234", got)
	}
}

func TestStructuresAreNullForScalarFunctionsAndInequal(t *testing.T) {
	// "Maps and lists are treated as null for string, number, and
	// datetime functions" — the documented example: display
	// toupper(json_message) yields empty results — and "Comparing map and
	// list to any other fields result in false." General and JSON
	// functions keep taking structures by their own contracts.
	events := testEvents(`{"user": {"id": "1"}}`)
	rows, err := executeQueryContext(newTestCtx(events),
		`fields toupper(user) as up, strlen(user) as ln, jsonstringify(user) as js`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rowField(t, rows, 0, "up"); got != "" {
		t.Fatalf("toupper over a map = %q, want empty", got)
	}
	if got := rowField(t, rows, 0, "ln"); got != "0" {
		t.Fatalf("strlen over a map = %q, want 0", got)
	}
	if got := rowField(t, rows, 0, "js"); got != `{"id":"1"}` {
		t.Fatalf("jsonstringify over a map = %q", got)
	}

	// Structures never compare equal — not even to themselves.
	rows2, err := executeQueryContext(newTestCtx(events), `filter user = user`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 0 {
		t.Fatalf("map = map must be false: %d rows", len(rows2))
	}
	// ispresent keeps its own contract: a parsed map is present.
	rows3, err := executeQueryContext(newTestCtx(events), `filter ispresent(user)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows3) != 1 {
		t.Fatalf("ispresent over a map = %d rows, want the event", len(rows3))
	}
}

func TestSystemFieldBacktickForm(t *testing.T) {
	// "For system fields that start with @ and contain special characters
	// in the remaining name, place the @ symbol outside the backticks and
	// enclose only the field name portion." The @ composes with the
	// backticked name into one field reference instead of a nested path.
	events := testEvents(`{"@aws.tag.aws:cloudformation:stack-name": "my-stack", "plain": "x"}`)
	rows, err := executeQueryContext(newTestCtx(events),
		"filter @`@aws.tag.aws:cloudformation:stack-name` = \"my-stack\"")
	if err != nil {
		t.Fatal(err)
	}
	// A log key starting with @ is discovered with an additional @
	// prefix, so the composed reference addresses @@aws.tag....
	if len(rows) != 1 {
		t.Fatalf("system-field backtick form = %d rows, want the event", len(rows))
	}
}

func TestParseRegexAfterDivisionAndGlobAsRequired(t *testing.T) {
	// A division slash in an earlier command still means division when a
	// later parse command carries a regex literal; the glob mode still
	// requires its as clause.
	events := []logEventWithContext{
		{timestamp: 1700000000000, message: `{"dur": {"ms": 8}, "eni": "eni-07"}`},
		{timestamp: 1700000000000, message: `{"dur": {"ms": 1}, "eni": "eni-07"}`},
	}
	rows, err := executeQueryContext(newTestCtx(events),
		`filter dur.ms / 2 > 1 | parse @message /"(?<x>eni-\d+)"/ | display x`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rowField(t, rows, 0, "x") != "eni-07" {
		t.Fatalf("rows = %+v, want the dur.ms > 2 row with x = eni-07", rows)
	}
	if _, err := executeQueryContext(newTestCtx(events), `parse @message "eni-*" `); err == nil {
		t.Fatal("glob mode without an as clause should be rejected")
	}
	if _, err := executeQueryContext(newTestCtx(events), `parse @message /(?<x>a)/ trailing`); err == nil {
		t.Fatal("a trailing token after the regex-mode pattern should be rejected")
	}
}

// assertWindowComparison runs one window-comparing query (diff or logcompare
// timeshift) over the shared two-window fixture — a current window holding
// two started events and one crashed event, a previous window holding one
// started event — and asserts the comparison both commands report over it:
// the started pattern grew by one and the crashed pattern is new.
func assertWindowComparison(t *testing.T, query string) {
	t.Helper()
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
	rows, err := executeQueryContext(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
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

func TestDiffComparesPreviousWindow(t *testing.T) {
	assertWindowComparison(t, `diff`)
}

func TestLogCompareTimeshift(t *testing.T) {
	// The timeshift window covers the same previous hour the diff test
	// uses, so the expectations match: started-pattern diff +1, crashed
	// pattern is new.
	assertWindowComparison(t, `logcompare timeshift 1h`)
}

// --- subquery / join / appendcols ---

// A dotted group key keeps its full path as the output column: the first
// segment alone would merge sibling fields (request.id and
// request.type) under one label.
func TestStatsDottedGroupKeyKeepsFullPath(t *testing.T) {
	events := testEvents(
		`{"request": {"id": "a", "type": "x"}}`,
	)
	rows, err := executeQueryContext(newTestCtx(events),
		`stats count(*) by request.id, request.type`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rowField(t, rows, 0, "request.id") != "a" || rowField(t, rows, 0, "request.type") != "x" {
		t.Fatalf("dotted group columns: %+v", rows[0].fields)
	}
}

// A computed aggregation that is neither a plain function call nor a
// literal still labels its output column with the expression text.
func TestStatsComputedExpressionColumn(t *testing.T) {
	events := testEvents(`{"v": 1}`)
	rows, err := executeQueryContext(newTestCtx(events), `stats count(*) + 1`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows, 0, "count(*) + 1") != "2" {
		t.Fatalf("computed column: %+v", rows[0].fields)
	}
}

// topk's field is the SECOND argument ("topk(k: number, fieldName:
// LogField)"): counting the first argument counted the literal k once
// per row, collapsing every result to the single-element ["2"].
func TestStatsTopkCountsTheFieldNotK(t *testing.T) {
	events := testEvents(
		`{"status": "200"}`, `{"status": "200"}`, `{"status": "200"}`,
		`{"status": "500"}`, `{"status": "500"}`,
		`{"status": "404"}`,
	)
	rows, err := executeQueryContext(newTestCtx(events), `stats topk(2, status) as t`)
	if err != nil {
		t.Fatal(err)
	}
	got := rowField(t, rows, 0, "t")
	if !strings.Contains(got, `"200"`) || !strings.Contains(got, `"500"`) ||
		strings.Contains(got, `"404"`) || strings.Contains(got, `"2"`) {
		t.Fatalf("topk(2, status) = %q, want the two most frequent statuses [200 500]", got)
	}
}

// Renaming a topk aggregation must not smuggle it past the topk+by
// rejection: the gate keys on the parsed function, not the alias.
func TestStatsTopkAliasStillRejectsBy(t *testing.T) {
	events := testEvents(`{"ip": "1.2.3.4"}`)
	if _, err := executeQueryContext(newTestCtx(events), `stats topk(3, ip) as t by bin(5m)`); err == nil {
		t.Fatal("topk renamed with `as` bypassed the topk+by rejection")
	}
}

// The limit command carries the documented ceiling ("You can specify a
// limit value of up to 100,000") — one over the ceiling rejects at
// compile time.
func TestLimitCommandCeiling(t *testing.T) {
	events := testEvents(`{"v": 1}`)
	if _, err := executeQueryContext(newTestCtx(events), "fields @message | limit 100001"); err == nil {
		t.Fatal("limit above the ceiling accepted")
	}
	if _, err := executeQueryContext(newTestCtx(events), "fields @message | limit 100000"); err != nil {
		t.Fatalf("limit at the ceiling rejected: %v", err)
	}
}

// count(fieldName) "counts all records that include the specified field
// name" — an empty-string value is an included field — and coalesce
// "returns the first non-null value", so an empty string wins over a
// later argument.
func TestCountAndCoalesceEmptyStringSemantics(t *testing.T) {
	events := testEvents(
		`{"note": ""}`, `{"note": "x"}`, `{"other": 1}`,
	)
	rows, err := executeQueryContext(newTestCtx(events), `stats count(note) as n, count(coalesce(note, "fallback")) as c`)
	if err != nil {
		t.Fatal(err)
	}
	if rowField(t, rows, 0, "n") != "2" {
		t.Fatalf("count(note) = %q, want 2 (the empty string is a present value)", rowField(t, rows, 0, "n"))
	}
	rows, err = executeQueryContext(newTestCtx(events), `fields coalesce(note, "fallback") as picked`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 || rowField(t, rows, 0, "picked") != "" || rowField(t, rows, 2, "picked") != "fallback" {
		t.Fatalf("coalesce must return the empty string (the first non-null value) and fall back only on null, got %+v", rows)
	}
}

// A multi-stats query aggregates its own aggregates: the subsequent
// stats reads the preceding stats' output columns, and the rejected
// forms (raw-field references, bin without a propagated timestamp) are
// compile-time errors, never silent nulls.
func TestMultiStatsAggregatesAggregates(t *testing.T) {
	events := testEvents(
		`{"op": "read", "ms": 10}`,
		`{"op": "read", "ms": 20}`,
		`{"op": "write", "ms": 100}`,
	)
	rows, err := executeQueryContext(newTestCtx(events), `stats sum(ms) as total by op | stats max(total) as busiest`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rowField(t, rows, 0, "busiest") != "100" {
		t.Fatalf("aggregated-of-aggregated rows = %+v, want one row with busiest=100", rows)
	}

	for _, q := range []string{
		`stats sum(ms) as total by op | stats max(strlen(ms)) as m`,
		`stats sum(ms) as total by op | stats avg(total) by bin(5m)`,
	} {
		if err := validateQueryPipeline(q); err == nil {
			t.Fatalf("documented-invalid multi-stats query accepted: %q", q)
		}
	}
}

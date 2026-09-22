// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"sync"
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// valueMetricCall records one metric emission with the value published.
type valueMetricCall struct {
	namespace string
	metric    string
	value     float64
	unit      string
	ts        time.Time
}

// valueRecordingMetricInvoker captures every emission's value, unlike the
// count-only recording invoker the ingestion tests use. The filter fan-out
// appends from its worker goroutines, so every access takes the mutex and
// tests read through the recorded snapshot.
type valueRecordingMetricInvoker struct {
	mu    sync.Mutex
	calls []valueMetricCall
	dims  []map[string]string
}

func (r *valueRecordingMetricInvoker) PutMetricData(_, namespace, metricName string, value float64, ts time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, valueMetricCall{namespace: namespace, metric: metricName, value: value, ts: ts})
	r.dims = append(r.dims, nil)
	return nil
}

func (r *valueRecordingMetricInvoker) PutMetricDataWithDimensions(_, namespace, metricName string, dimensions map[string]string, value float64, ts time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, valueMetricCall{namespace: namespace, metric: metricName, value: value, ts: ts})
	r.dims = append(r.dims, dimensions)
	return nil
}

func (r *valueRecordingMetricInvoker) PutMetricDataWithDimensionsAndUnit(_, namespace, metricName string, dimensions map[string]string, unit string, value float64, ts time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, valueMetricCall{namespace: namespace, metric: metricName, value: value, unit: unit, ts: ts})
	r.dims = append(r.dims, dimensions)
	return nil
}

// recorded returns a snapshot of the emissions and their dimensions; the
// two slices are index-aligned.
func (r *valueRecordingMetricInvoker) recorded() ([]valueMetricCall, []map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	calls := append([]valueMetricCall(nil), r.calls...)
	dims := append([]map[string]string(nil), r.dims...)
	return calls, dims
}

// newMetricEmissionTestService builds a service with storage and the
// value-recording metric invoker wired, plus the group the tests filter.
func newMetricEmissionTestService(t *testing.T) (*LogsService, *valueRecordingMetricInvoker, string) {
	t.Helper()
	svc, store := newTestService(t)

	metrics := &valueRecordingMetricInvoker{}
	svc.SetCloudWatchMetricInvoker(metrics)

	const group = "emission-group"
	createTestLogGroup(t, store, group)
	return svc, metrics, group
}

func putEmissionFilter(t *testing.T, svc *LogsService, group, name, pattern, metricValue string, defaultValue float64, defaultSet bool) {
	t.Helper()
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	if err := store.PutMetricFilter(&logsstore.MetricFilter{
		Name:          name,
		LogGroupName:  group,
		FilterPattern: pattern,
		MetricTransformations: []logsstore.MetricTransformation{{
			MetricName:      name,
			MetricNamespace: "Emission",
			MetricValue:     metricValue,
			DefaultValue:    defaultValue,
			DefaultValueSet: defaultSet,
		}},
	}); err != nil {
		t.Fatalf("put metric filter %s: %v", name, err)
	}
}

func emissionStore(t *testing.T, svc *LogsService) *logsstore.Store {
	t.Helper()
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	return store
}

// TestMetricEmissionValuesNumbersAndReferences pins the metricValue
// semantics: a static number publishes that number per match, a $.field
// reference publishes the event's numeric field value, and a reference the
// event cannot resolve (missing field or non-numeric value) contributes no
// data point. The patterns use the engine's supported comparison forms —
// the unquoted wildcard right-hand side of the guide's worked example is
// an engine-level gap tracked separately.
func TestMetricEmissionValuesNumbersAndReferences(t *testing.T) {
	svc, metrics, group := newMetricEmissionTestService(t)
	putEmissionFilter(t, svc, group, "static", "ERROR", "3", 0, false)
	putEmissionFilter(t, svc, group, "reference", `{ $.latency > 0 }`, "$.latency", 0, false)
	putEmissionFilter(t, svc, group, "missingref", `{ $.requestType = "GET" }`, "$.latency", 0, false)
	store := emissionStore(t, svc)

	now := time.Now().UnixMilli()
	svc.evaluateMetricFilters(store, "us-east-1", group, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: "ERROR one"},
		{Timestamp: now, Message: "ERROR two"},
		{Timestamp: now, Message: `{"latency": 50, "requestType": "GET"}`},
		{Timestamp: now, Message: `{"latency": "slow", "requestType": "GET"}`},
		{Timestamp: now, Message: `{"requestType": "GET"}`},
	}, nil)

	calls, _ := metrics.recorded()
	var staticCount int
	var refValues []float64
	for _, c := range calls {
		switch c.metric {
		case "static":
			staticCount++
			if c.value != 3 {
				t.Fatalf("static metricValue 3 must publish 3, got %v", c.value)
			}
		case "reference":
			refValues = append(refValues, c.value)
		case "missingref":
			t.Fatalf("a reference the event does not carry must contribute no data point, got %v", c.value)
		}
	}
	if staticCount != 2 {
		t.Fatalf("two matching events must publish the static value twice, got %d", staticCount)
	}
	// Only the numeric latency event resolves: the "slow" event matches
	// the pattern but its referenced value is not numeric, and the
	// field-less event's $.latency is absent from the extraction.
	if len(refValues) != 1 || refValues[0] != 50 {
		t.Fatalf("the $.latency reference must publish the event's 50 alone, got %v", refValues)
	}
}

// TestMetricEmissionDefaultValueAtMinuteClose pins the default value
// cadence: a matchless evaluation publishes nothing by itself (the minute
// may still observe a later match), the default of an entirely matchless
// minute publishes once the minute has closed — the next evaluation
// carries the close in-band, and the minute-close worker carries it when
// the group goes quiet — and a minute that observed a match never reports
// the default.
func TestMetricEmissionDefaultValueAtMinuteClose(t *testing.T) {
	svc, metrics, group := newMetricEmissionTestService(t)
	putEmissionFilter(t, svc, group, "count", "ERROR", "1", 0, true)
	store := emissionStore(t, svc)

	now := time.Now().UnixMilli()
	// Five non-matching events in the running minute: nothing publishes.
	svc.evaluateMetricFilters(store, "us-east-1", group, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: "quiet one"},
		{Timestamp: now, Message: "quiet two"},
		{Timestamp: now, Message: "quiet three"},
		{Timestamp: now, Message: "quiet four"},
		{Timestamp: now, Message: "quiet five"},
	}, nil)
	if calls, _ := metrics.recorded(); len(calls) != 0 {
		t.Fatalf("a matchless batch of the running minute must publish nothing, got %v", calls)
	}

	// A second matchless batch, same minute: still nothing.
	svc.evaluateMetricFilters(store, "us-east-1", group, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: "still quiet"},
	}, nil)
	if calls, _ := metrics.recorded(); len(calls) != 0 {
		t.Fatalf("a second matchless batch of the running minute must publish nothing, got %v", calls)
	}

	// A matching batch in the same minute publishes the match value alone.
	svc.evaluateMetricFilters(store, "us-east-1", group, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: "ERROR match"},
	}, nil)
	if calls, _ := metrics.recorded(); len(calls) != 1 || calls[0].value != 1 {
		t.Fatalf("a match must publish the metric value alone, got %v", calls)
	}

	// A matchless batch after the match, same minute: the minute saw a
	// match, so no default either.
	svc.evaluateMetricFilters(store, "us-east-1", group, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: "quiet again"},
	}, nil)
	if calls, _ := metrics.recorded(); len(calls) != 1 {
		t.Fatalf("a matchless batch after a same-minute match must publish nothing, got %v", calls)
	}

	// A closed matchless minute publishes its default through the next
	// evaluation's in-band close, at the closed minute's own timestamp.
	fresh, freshMetrics, freshGroup := newMetricEmissionTestService(t)
	putEmissionFilter(t, fresh, freshGroup, "count", "ERROR", "1", 0, true)
	elapsedMinute := time.Now().Add(-2 * time.Minute)
	fresh.recordMetricMinutePresence("us-east-1", freshGroup, "count", elapsedMinute)
	fresh.evaluateMetricFilters(emissionStore(t, fresh), "us-east-1", freshGroup, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: "quiet after the close"},
	}, nil)
	calls, _ := freshMetrics.recorded()
	if len(calls) != 1 || calls[0].value != 0 {
		t.Fatalf("a closed matchless minute must publish its default once, got %v", calls)
	}
	wantTS := time.Unix(elapsedMinute.Unix()/60*60, 0)
	if !calls[0].ts.Equal(wantTS) {
		t.Fatalf("the default must carry the closed minute's timestamp: got %v want %v", calls[0].ts, wantTS)
	}

	// A matched closed minute owes nothing: the match minute above never
	// published a default even though a later evaluation ran.
	if calls, _ := metrics.recorded(); len(calls) != 1 {
		t.Fatalf("a closed minute that observed a match must never report the default, got %v", calls)
	}

	// The minute-close worker carries the close when no later batch
	// arrives, and a claimed minute publishes exactly once.
	quiet, quietMetrics, quietGroup := newMetricEmissionTestService(t)
	putEmissionFilter(t, quiet, quietGroup, "count", "ERROR", "1", 0, true)
	quiet.recordMetricMinutePresence("us-east-1", quietGroup, "count", elapsedMinute)
	quiet.closeQuietMetricDefaultMinutes()
	calls, _ = quietMetrics.recorded()
	if len(calls) != 1 || calls[0].value != 0 {
		t.Fatalf("the minute-close worker must publish the owed default, got %v", calls)
	}
	quiet.closeQuietMetricDefaultMinutes()
	if calls, _ := quietMetrics.recorded(); len(calls) != 1 {
		t.Fatalf("a claimed closed minute must not publish twice, got %v", calls)
	}

	// The state-machine seam: the running minute claims nothing, a closed
	// matchless minute claims exactly once, a closed matched minute never
	// claims, and the claim state resets per minute. The instants are
	// minute-aligned so the boundaries do not depend on where in the
	// running minute the test starts.
	seam := &LogsService{}
	base := time.Now().Unix() / 60 * 60
	running := time.Unix(base, 0)
	closed := time.Unix(base+60, 0)
	matched := time.Unix(base+120, 0)
	nextRunning := time.Unix(base+180, 0)
	seam.recordMetricMinutePresence("us-east-1", "g", "f", running)
	if _, owed := seam.claimElapsedMetricDefaultMinute("us-east-1", "g", "f", running); owed {
		t.Fatal("the running minute must not claim the default")
	}
	if _, owed := seam.claimElapsedMetricDefaultMinute("us-east-1", "g", "f", closed); !owed {
		t.Fatal("a closed matchless minute must claim the default")
	}
	if _, owed := seam.claimElapsedMetricDefaultMinute("us-east-1", "g", "f", closed); owed {
		t.Fatal("a claimed minute must not claim again")
	}
	seam.markMetricMatched("us-east-1", "g", "f", matched)
	if _, owed := seam.claimElapsedMetricDefaultMinute("us-east-1", "g", "f", matched.Add(60*time.Second)); owed {
		t.Fatal("a minute that observed a match must not report the default")
	}
	seam.recordMetricMinutePresence("us-east-1", "g", "f", matched.Add(60*time.Second))
	if _, owed := seam.claimElapsedMetricDefaultMinute("us-east-1", "g", "f", nextRunning.Add(60*time.Second)); !owed {
		t.Fatal("the next matchless minute must claim the default again")
	}
}

// TestTestMetricFilterEventNumberIsInputPosition pins that eventNumber is
// the matched event's position in the input list, not a match ordinal.
func TestTestMetricFilterEventNumberIsInputPosition(t *testing.T) {
	svc := &LogsService{}
	matches, err := svc.testMetricFilterCore(testMetricFilterInput{
		FilterPattern:    "ERROR",
		FilterPatternSet: true,
		LogEventMessages: []string{
			"nothing here",
			"ERROR at two",
			"still nothing",
			"warn only",
			"ERROR at five",
		},
	})
	if err != nil {
		t.Fatalf("testMetricFilterCore: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("two messages match, got %d records", len(matches))
	}
	for i, want := range []int64{2, 5} {
		got, _ := matches[i]["eventNumber"].(int64)
		if got != want {
			t.Fatalf("match %d must report input position %d, got %v", i, want, matches[i]["eventNumber"])
		}
	}
}

// TestTestMetricFilterEmptyPatternMatchesAll pins the requiredness
// semantics the Put operation already implements for the same member:
// the pattern's shape allows the empty string, so a present empty
// pattern matches every message while an absent member still rejects.
func TestTestMetricFilterEmptyPatternMatchesAll(t *testing.T) {
	svc := &LogsService{}
	matches, err := svc.testMetricFilterCore(testMetricFilterInput{
		FilterPattern:    "",
		FilterPatternSet: true,
		LogEventMessages: []string{"a", "b"},
	})
	if err != nil {
		t.Fatalf("present empty pattern must be legal: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("empty pattern matches all, got %d matches", len(matches))
	}

	_, err = svc.testMetricFilterCore(testMetricFilterInput{
		FilterPattern:    "",
		FilterPatternSet: false,
		LogEventMessages: []string{"a"},
	})
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("absent pattern: code=%q err=%v", code, err)
	}
}

// TestMetricEmissionDimensionsUnitAndSystemFields pins the transformation's
// dimensions and unit members at the emission seam: dimension values are
// value references into the matched event (the operation documentation's
// worked example runs "dimensions": {"Request": "$request"}), a reference
// the event cannot resolve contributes no data point, literals pass
// through, the emitSystemFieldDimensions members contribute @aws.account
// and @aws.region, and the unit rides the datum ("The unit to assign to
// the metric").
func TestMetricEmissionDimensionsUnitAndSystemFields(t *testing.T) {
	svc, metrics, group := newMetricEmissionTestService(t)
	store := emissionStore(t, svc)

	putDimFilter := func(name, pattern string, dims map[string]string, unit string, emit ...string) {
		t.Helper()
		if err := store.PutMetricFilter(&logsstore.MetricFilter{
			Name: name, LogGroupName: group, FilterPattern: pattern,
			EmitSystemFieldDimensions: emit,
			MetricTransformations: []logsstore.MetricTransformation{{
				MetricName: name, MetricNamespace: "Emission", MetricValue: "1",
				Dimensions: dims, Unit: unit,
			}},
		}); err != nil {
			t.Fatalf("put metric filter %s: %v", name, err)
		}
	}
	putDimFilter("dimref", `{ $.requestType = "GET" }`, map[string]string{"Request": "$.requestType"}, "Count")
	putDimFilter("dimmissing", `{ $.requestType = "GET" }`, map[string]string{"Missing": "$.absent"}, "")
	putDimFilter("sysfields", `{ $.requestType = "GET" }`, nil, "", "@aws.account", "@aws.region")
	putDimFilter("unitonly", `{ $.requestType = "GET" }`, nil, "Seconds")

	now := time.Now().UnixMilli()
	svc.evaluateMetricFilters(store, "us-east-1", group, "stream", []logsstore.LogEntry{
		{Timestamp: now, Message: `{"requestType":"GET","latency":50}`},
		{Timestamp: now, Message: `{"requestType":"POST","latency":10}`},
	}, nil)

	type row struct {
		metric string
		dims   map[string]string
		unit   string
	}
	want := map[string]row{
		"dimref":    {metric: "dimref", dims: map[string]string{"Request": "GET"}, unit: "Count"},
		"sysfields": {metric: "sysfields", dims: map[string]string{"@aws.account": "000000000000", "@aws.region": "us-east-1"}, unit: ""},
		"unitonly":  {metric: "unitonly", dims: nil, unit: "Seconds"},
	}
	metricCalls, metricDims := metrics.recorded()
	seen := map[string]int{}
	for i, call := range metricCalls {
		w, ok := want[call.metric]
		if !ok {
			continue
		}
		seen[call.metric]++
		if call.unit != w.unit {
			t.Fatalf("call %d (%s) unit = %q, want %q", i, call.metric, call.unit, w.unit)
		}
		got := metricDims[i]
		if len(got) != len(w.dims) {
			t.Fatalf("call %d (%s) dims = %v, want %v", i, call.metric, got, w.dims)
		}
		for k, v := range w.dims {
			if got[k] != v {
				t.Fatalf("call %d (%s) dims = %v, want %v", i, call.metric, got, w.dims)
			}
		}
	}
	for _, m := range []string{"dimref", "sysfields", "unitonly"} {
		// The POST event matches none of the JSON patterns' requestType
		// guard, so each filter emits exactly once.
		if seen[m] != 1 {
			t.Fatalf("metric %s emitted %d times, want 1", m, seen[m])
		}
	}
	if seen["dimmissing"] != 0 {
		t.Fatalf("unresolvable dimension reference still emitted %d times", seen["dimmissing"])
	}
}

package cloudwatch

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

func newAlarmTestStore(t *testing.T) *cwstore.AlarmStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return cwstore.NewAlarmStore(st, "000000000000", "us-east-1")
}

// A rule referencing its own alarm must be rejected, as must a transitive
// cycle between composites; acyclic references must pass.
func TestValidateCompositeAcyclic(t *testing.T) {
	store := newAlarmTestStore(t)

	composite := func(name, rule string) *cwstore.Alarm {
		a := cwstore.NewAlarm(name, "", "")
		a.AlarmType = cwstore.AlarmTypeCompositeAlarm
		a.AlarmRule = rule
		return a
	}
	if _, err := store.CreateAlarm(composite("base", `ALARM("base-metric")`)); err != nil {
		t.Fatal(err)
	}

	selfRule, err := parseAlarmRule(`ALARM("self-ref")`)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateCompositeAcyclic(store, "self-ref", selfRule); err == nil {
		t.Fatal("self-referencing rule accepted")
	}

	okRule, err := parseAlarmRule(`ALARM("base")`)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateCompositeAcyclic(store, "child", okRule); err != nil {
		t.Fatalf("acyclic rule rejected: %v", err)
	}
	if _, err := store.CreateAlarm(composite("child", `ALARM("base")`)); err != nil {
		t.Fatal(err)
	}

	// "base" updated to reference "child" closes the cycle base → child → base.
	cycleRule, err := parseAlarmRule(`ALARM("child") AND ALARM("base-metric")`)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateCompositeAcyclic(store, "base", cycleRule); err == nil {
		t.Fatal("transitive cycle accepted")
	}
}

// The creation path must reject a cyclic rule before persisting anything.
func TestPutCompositeAlarmRejectsCycle(t *testing.T) {
	store := newAlarmTestStore(t)
	svc := &CloudWatchService{}
	stores := &cloudwatchStores{alarms: store}

	if _, err := svc.putCompositeAlarmCore(stores, &PutCompositeAlarmInput{
		AlarmName: "comp-a",
		AlarmRule: `ALARM("metric-x")`,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.putCompositeAlarmCore(stores, &PutCompositeAlarmInput{
		AlarmName: "comp-b",
		AlarmRule: `ALARM("comp-a")`,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.putCompositeAlarmCore(stores, &PutCompositeAlarmInput{
		AlarmName: "comp-a",
		AlarmRule: `ALARM("comp-b")`,
	}); err == nil {
		t.Fatal("cyclic update accepted")
	}
}

// Cyclic nodes are reported separately so the rest of the graph still
// evaluates.
func TestTopologicalSortLevelsCyclic(t *testing.T) {
	dependencies := map[string][]string{
		"a": {"b"}, // a depends on b
		"b": {"a"}, // b depends on a — cycle
		"c": {"b"}, // downstream of the cycle
		"d": {},    // independent
	}
	levels, cyclic := topologicalSortLevels(dependencies)

	cyclicSet := map[string]bool{}
	for _, n := range cyclic {
		cyclicSet[n] = true
	}
	if !cyclicSet["a"] || !cyclicSet["b"] || !cyclicSet["c"] {
		t.Fatalf("expected a, b and c as cyclic/downstream, got %v", cyclic)
	}

	// d must still be sorted into a level.
	found := false
	for _, level := range levels {
		for _, n := range level {
			if n == "d" {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("independent node d missing from levels %v", levels)
	}
}

// countBreaches must treat the window as half-open [start, end): the
// bucket at the upper boundary belongs to the next evaluation.
func TestCountBreachesHalfOpenWindow(t *testing.T) {
	alarm := &cwstore.Alarm{
		Threshold:          10,
		ComparisonOperator: "GreaterThanThreshold",
		Statistic:          "Average",
		EvaluationPeriods:  1,
		DatapointsToAlarm:  1,
		TreatMissingData:   "ignore",
		Period:             60,
	}
	start := time.Unix(1000000000, 0).UTC().Truncate(60 * time.Second)
	end := start.Add(60 * time.Second)

	stats := []*cwstore.MetricStatistics{
		{Timestamp: start, Average: 50},                        // in window, breaching
		{Timestamp: end, Average: 50},                          // next window's bucket
		{Timestamp: start.Add(-60 * time.Second), Average: 50}, // previous window
	}
	breaches, inWindow := countBreaches(alarm, stats, start, end)
	if breaches != 1 || inWindow != 1 {
		t.Fatalf("expected 1 breach in a 1-bucket window, got breaches=%d inWindow=%d", breaches, inWindow)
	}
}

// Alarms with a period shorter than the tick advance their watermark to
// the newest completed boundary and never re-evaluate old boundaries.
func TestEvaluateAlarmJobSubMinuteWatermark(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	metricStore := cwstore.NewMetricChunkStore(st, "us-east-1", t.TempDir())

	e := newAlarmEvaluator(60*time.Second, 1, nil)

	alarm := cwstore.NewAlarm("sub-minute", "ns", "metric")
	alarm.Period = 10
	alarm.EvaluationPeriods = 1
	alarm.DatapointsToAlarm = 1
	alarm.TreatMissingData = "ignore"
	alarm.State = "INSUFFICIENT_DATA" // no transition from empty stats

	now := time.Now().UTC()
	end := now.Truncate(10 * time.Second)

	// Untracked alarm: single evaluation, watermark set to the newest
	// completed boundary.
	if result := e.evaluateAlarmJob("us-east-1", alarm, metricStore); result != nil {
		t.Fatalf("unexpected transition: %+v", result)
	}
	if got := e.lastEvaluated["us-east-1/sub-minute"]; !got.Equal(end) {
		t.Fatalf("watermark %v, want %v", got, end)
	}

	// No new boundary since the last tick: nothing to do.
	if result := e.evaluateAlarmJob("us-east-1", alarm, metricStore); result != nil {
		t.Fatalf("unexpected transition without a new boundary: %+v", result)
	}
	if got := e.lastEvaluated["us-east-1/sub-minute"]; !got.Equal(end) {
		t.Fatalf("watermark moved without a new boundary: %v", got)
	}

	// A watermark six boundaries old triggers one evaluation per missed
	// boundary; with empty stats no transition fires and the watermark
	// catches up to the newest boundary.
	old := end.Add(-60 * time.Second)
	e.lastEvaluated["us-east-1/sub-minute"] = old
	if result := e.evaluateAlarmJob("us-east-1", alarm, metricStore); result != nil {
		t.Fatalf("unexpected transition during catch-up: %+v", result)
	}
	if got := e.lastEvaluated["us-east-1/sub-minute"]; !got.Equal(end) {
		t.Fatalf("watermark %v after catch-up, want %v", got, end)
	}
}

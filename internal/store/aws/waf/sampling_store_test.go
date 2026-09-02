package waf

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

func newSamplingStoreForTest(t *testing.T) (*SamplingStore, string) {
	t.Helper()
	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return NewSamplingStore(st), dir
}

func TestSamplingStoreRecordAndQuery(t *testing.T) {
	store, _ := newSamplingStoreForTest(t)
	base := time.Now()
	for i := 0; i < 3; i++ {
		err := store.Record("arn", "ruleMetric", SampledRequest{
			RuleNameWithinRuleGroup: "group#rule",
			MetricName:              "ruleMetric",
			Action:                  "Block",
			Timestamp:               base.Add(time.Duration(i) * time.Second),
			ClientIP:                "203.0.113.9",
			URI:                     "/blocked",
			Method:                  "GET",
			HTTPVersion:             "HTTP/1.1",
			Headers:                 []SampledHTTPHeader{{Name: "Host", Value: "example.test"}},
			ResponseCodeSent:        403,
			Labels:                  []string{"awswaf:111122223333:webacl:test:suspect"},
			RequestHeadersInserted:  []SampledHTTPHeader{{Name: "x-amzn-waf-counted", Value: "1"}},
			OverriddenAction:        "Block",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	all, err := store.Query("arn", "ruleMetric", base.Add(-time.Minute), base.Add(time.Minute), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("records = %d, want 3", len(all))
	}
	// Newest first.
	if !all[0].Timestamp.After(all[1].Timestamp) {
		t.Fatalf("records are not newest-first: %+v", all)
	}
	if all[0].URI != "/blocked" || all[0].Headers[0].Name != "Host" {
		t.Fatalf("record fields lost: %+v", all[0])
	}
	if all[0].RuleNameWithinRuleGroup != "group#rule" || all[0].ResponseCodeSent != 403 ||
		len(all[0].Labels) != 1 || len(all[0].RequestHeadersInserted) != 1 || all[0].OverriddenAction != "Block" {
		t.Fatalf("sampled-request disposition members lost: %+v", all[0])
	}
	limited, err := store.Query("arn", "ruleMetric", base.Add(-time.Minute), base.Add(time.Minute), 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(limited) != 2 || !limited[0].Timestamp.After(limited[1].Timestamp) {
		t.Fatalf("limited records = %+v", limited)
	}
	windowed, err := store.Query("arn", "ruleMetric", base.Add(1500*time.Millisecond), base.Add(time.Minute), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(windowed) != 1 {
		t.Fatalf("windowed records = %d, want 1", len(windowed))
	}
	missing, err := store.Query("arn", "otherMetric", base.Add(-time.Minute), base.Add(time.Minute), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 0 {
		t.Fatalf("unknown metric records = %d, want 0", len(missing))
	}
}

// TestSamplingStoreSurvivesReopen pins the persistence requirement:
// samples written before a restart stay retrievable within the
// retention window afterwards.
func TestSamplingStoreSurvivesReopen(t *testing.T) {
	dir := t.TempDir()
	st, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	first := NewSamplingStore(st)
	now := time.Now()
	if err := first.Record("arn", "metric", SampledRequest{
		RuleNameWithinRuleGroup: "group#rule",
		MetricName:              "metric",
		Action:                  "Count",
		Timestamp:               now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := storage.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	second := NewSamplingStore(reopened)
	records, err := second.Query("arn", "metric", now.Add(-time.Minute), now.Add(time.Minute), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Action != "Count" {
		t.Fatalf("records after reopen = %+v", records)
	}
}

func TestSamplingStorePurgeExpired(t *testing.T) {
	store, _ := newSamplingStoreForTest(t)
	now := time.Now()

	// One expired record, SamplingPopulationDepth+1 fresh records of one
	// rule, and one fresh record of another rule.
	if err := store.Record("arn", "metric", SampledRequest{
		MetricName: "metric",
		Timestamp:  now.Add(-SampleRetention - time.Minute),
	}); err != nil {
		t.Fatal(err)
	}
	for i := 0; i <= SamplingPopulationDepth; i++ {
		if err := store.Record("arn", "metric", SampledRequest{
			MetricName: "metric",
			Timestamp:  now.Add(time.Duration(i) * time.Second),
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.Record("arn", "other", SampledRequest{
		MetricName: "other",
		Timestamp:  now,
	}); err != nil {
		t.Fatal(err)
	}

	if err := store.PurgeExpired(now); err != nil {
		t.Fatal(err)
	}

	trimmed, err := store.Query("arn", "metric", now.Add(-2*time.Hour), now.Add(2*time.Hour), 10000)
	if err != nil {
		t.Fatal(err)
	}
	if len(trimmed) != SamplingPopulationDepth {
		t.Fatalf("trimmed records = %d, want %d", len(trimmed), SamplingPopulationDepth)
	}
	// The trim keeps the newest records: the first second of the fresh
	// run is dropped, the expired record is gone. trimmed is
	// newest-first, so its tail is the oldest retained record.
	newest := now.Add(time.Duration(SamplingPopulationDepth) * time.Second)
	if !trimmed[0].Timestamp.Equal(newest) {
		t.Fatalf("newest retained = %s, want %s", trimmed[0].Timestamp, newest)
	}
	if !trimmed[len(trimmed)-1].Timestamp.Equal(now.Add(time.Second)) {
		t.Fatalf("oldest retained = %s, want %s", trimmed[len(trimmed)-1].Timestamp, now.Add(time.Second))
	}
	kept, err := store.Query("arn", "other", now.Add(-time.Hour), now.Add(time.Hour), 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(kept) != 1 {
		t.Fatalf("other-rule records = %d, want 1", len(kept))
	}
}

// TestSamplingStoreCountPopulation pins the matched-request counters the
// population figure reads: every recorded match counts inside its window,
// out-of-window minutes do not, and the retention sweep drops expired
// counters alongside the records.
func TestSamplingStoreCountPopulation(t *testing.T) {
	store, _ := newSamplingStoreForTest(t)
	now := time.Now()

	for i := 0; i < 6; i++ {
		if err := store.Record("arn", "metric", SampledRequest{
			MetricName: "metric",
			Timestamp:  now.Add(time.Duration(i) * time.Second),
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.Record("arn", "metric", SampledRequest{
		MetricName: "metric",
		Timestamp:  now.Add(-4 * time.Hour),
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Record("arn", "other", SampledRequest{
		MetricName: "other",
		Timestamp:  now,
	}); err != nil {
		t.Fatal(err)
	}

	if got := store.CountPopulation("arn", "metric", now.Add(-time.Minute), now.Add(time.Minute)); got != 6 {
		t.Fatalf("population inside the window = %d, want 6", got)
	}
	if got := store.CountPopulation("arn", "metric", now.Add(-5*time.Hour), now.Add(-3*time.Hour)); got != 1 {
		t.Fatalf("population of the aged minute = %d, want 1", got)
	}
	if got := store.CountPopulation("arn", "other", now.Add(-time.Minute), now.Add(time.Minute)); got != 1 {
		t.Fatalf("other-rule population = %d, want 1", got)
	}

	if err := store.PurgeExpired(now); err != nil {
		t.Fatal(err)
	}
	if got := store.CountPopulation("arn", "metric", now.Add(-5*time.Hour), now.Add(-3*time.Hour)); got != 0 {
		t.Fatalf("aged counters survived the purge: %d", got)
	}
	if got := store.CountPopulation("arn", "metric", now.Add(-time.Minute), now.Add(time.Minute)); got != 6 {
		t.Fatalf("fresh counters must survive the purge, got %d", got)
	}
}

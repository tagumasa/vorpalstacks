package s3

import (
	"math"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// TestMetricsFilterMatches pins the filter evaluation the request-metrics
// aggregator applies: whole-bucket filters match everything, prefix filters
// match by key prefix, tag and conjunction predicates consult the object's
// tag set, an empty conjunction member list never fails a match, and an
// access-point predicate never matches on any surface.
func TestMetricsFilterMatches(t *testing.T) {
	objectTags := []types.Tag{{Key: "class", Value: "blue"}}
	tags := func() []types.Tag { return objectTags }

	if !metricsFilterMatches(nil, "any/key.txt", tags) {
		t.Fatal("a missing filter selects the whole bucket")
	}
	prefix := &s3store.MetricsFilter{Prefix: "docs/"}
	if !metricsFilterMatches(prefix, "docs/a.txt", tags) {
		t.Fatal("a matching prefix must select the object")
	}
	if metricsFilterMatches(prefix, "kept/a.txt", tags) {
		t.Fatal("a non-matching prefix must not select the object")
	}
	tag := &s3store.MetricsFilter{Tag: &types.Tag{Key: "class", Value: "blue"}}
	if !metricsFilterMatches(tag, "docs/a.txt", tags) {
		t.Fatal("a matching tag must select the object")
	}
	if metricsFilterMatches(tag, "docs/a.txt", func() []types.Tag { return nil }) {
		t.Fatal("an untagged object must not satisfy a tag filter")
	}
	and := &s3store.MetricsFilter{And: &s3store.MetricsAndOperator{
		Prefix: "docs/",
		Tags:   []types.Tag{{Key: "class", Value: "blue"}},
	}}
	if !metricsFilterMatches(and, "docs/a.txt", tags) {
		t.Fatal("a conjunction must select objects matching every member")
	}
	if metricsFilterMatches(and, "kept/a.txt", tags) {
		t.Fatal("a conjunction must not select objects outside its prefix")
	}
	andNoTags := &s3store.MetricsFilter{And: &s3store.MetricsAndOperator{Prefix: "docs/"}}
	if !metricsFilterMatches(andNoTags, "docs/a.txt", tags) {
		t.Fatal("a conjunction member list without tags must not consult the tag set")
	}
	// The platform has no access-point substrate, so an AccessPointArn
	// predicate can never be satisfied — regardless of its sibling members.
	accessPoint := &s3store.MetricsFilter{AccessPointArn: "arn:aws:s3:us-east-1:123456789012:accesspoint/ap-1"}
	if metricsFilterMatches(accessPoint, "docs/a.txt", tags) {
		t.Fatal("an access-point filter must never match")
	}
	andWithAccessPoint := &s3store.MetricsFilter{And: &s3store.MetricsAndOperator{
		Prefix:         "docs/",
		AccessPointArn: "arn:aws:s3:us-east-1:123456789012:accesspoint/ap-1",
	}}
	if metricsFilterMatches(andWithAccessPoint, "docs/a.txt", tags) {
		t.Fatal("a conjunction carrying an access-point ARN must never match, prefix notwithstanding")
	}
	andAccessPointTag := &s3store.MetricsFilter{And: &s3store.MetricsAndOperator{
		AccessPointArn: "arn:aws:s3:us-east-1:123456789012:accesspoint/ap-1",
		Tags:           []types.Tag{{Key: "class", Value: "blue"}},
	}}
	if metricsFilterMatches(andAccessPointTag, "docs/a.txt", tags) {
		t.Fatal("a conjunction carrying an access-point ARN must never match, tags notwithstanding")
	}
}

// TestRequestMetricAccumulation pins the per-request accumulation: every
// matched request counts AllRequests plus its operation metric, byte counts
// land on the upload/download metrics, and both error counters receive a
// 0-or-1 from every request so their SampleCount spans all requests and
// Average is the documented error rate.
func TestRequestMetricAccumulation(t *testing.T) {
	a := &requestMetricsAggregator{
		window:   time.Minute,
		counters: map[requestMetricsKey]*requestMetricTotals{},
		configs:  map[string]cachedMetricsConfigs{},
	}
	notFound := s3RequestObservation{
		region: "us-east-1", bucket: "bucket", key: "docs/a.txt",
		method: "GET", op: "GetRequests", bytes: 128, status: 404,
		firstByte: 1500 * time.Microsecond, total: 3 * time.Millisecond,
	}
	uploaded := s3RequestObservation{
		region: "us-east-1", bucket: "bucket", key: "docs/b.txt",
		method: "PUT", op: "PutRequests", bytes: 64, status: 200,
		firstByte: 500 * time.Microsecond, total: time.Millisecond,
	}
	a.accumulate(notFound, "filter-a")
	a.accumulate(notFound, "filter-a")
	a.accumulate(uploaded, "filter-a")
	a.accumulate(notFound, "filter-b")

	totals := a.counters[requestMetricsKey{region: "us-east-1", bucket: "bucket", id: "filter-a"}]
	if totals == nil {
		t.Fatal("matched requests must create a totals entry")
	}
	// metric -> {sum, sampleCount, min, max}
	for metric, want := range map[string][4]float64{
		"AllRequests":         {3, 3, 1, 1},
		"GetRequests":         {2, 2, 1, 1},
		"PutRequests":         {1, 1, 1, 1},
		"BytesDownloaded":     {256, 2, 128, 128},
		"BytesUploaded":       {64, 1, 64, 64},
		"4xxErrors":           {2, 3, 0, 1},
		"5xxErrors":           {0, 3, 0, 0},
		"FirstByteLatency":    {3.5, 3, 0.5, 1.5},
		"TotalRequestLatency": {7, 3, 1, 3},
	} {
		d := totals.dists[metric]
		if d == nil {
			t.Fatalf("%s must accumulate a distribution", metric)
		}
		got := [4]float64{d.sum, d.count, d.min, d.max}
		for i := range want {
			if math.Abs(got[i]-want[i]) > 1e-9 {
				t.Fatalf("%s = %v (sum/count/min/max), want %v", metric, got, want)
			}
		}
	}
	other := a.counters[requestMetricsKey{region: "us-east-1", bucket: "bucket", id: "filter-b"}]
	if other == nil || other.dists["AllRequests"] == nil || other.dists["AllRequests"].sum != 1 {
		t.Fatalf("filter-b totals = %+v, want a single AllRequests", other)
	}

	if unit := requestMetricsUnit("BytesUploaded"); unit != "Bytes" {
		t.Fatalf("byte metric unit = %q", unit)
	}
	if unit := requestMetricsUnit("TotalRequestLatency"); unit != "Milliseconds" {
		t.Fatalf("latency metric unit = %q", unit)
	}
	if unit := requestMetricsUnit("AllRequests"); unit != "Count" {
		t.Fatalf("count metric unit = %q", unit)
	}
}

// TestClassifyRequestMetrics pins the method/query classification against
// the AWS metric definitions: object-plane requests classify by method with
// list-oriented GETs on ListRequests; the bucket plane contributes
// HeadRequests, the list requests, DeleteObjects (as DeleteRequests), and
// nothing else beyond AllRequests; SelectObjectContent classifies as
// SelectRequests, never PostRequests.
func TestClassifyRequestMetrics(t *testing.T) {
	cases := []struct {
		name   string
		method string
		target string
		want   string
	}{
		{"object get", "GET", "/bucket/docs/a.txt", "GetRequests"},
		{"object tagging get", "GET", "/bucket/docs/a.txt?tagging", "GetRequests"},
		{"object list parts", "GET", "/bucket/docs/a.txt?uploadId=abc", "ListRequests"},
		{"object head", "HEAD", "/bucket/docs/a.txt", "HeadRequests"},
		{"object put", "PUT", "/bucket/docs/a.txt", "PutRequests"},
		{"object delete", "DELETE", "/bucket/docs/a.txt", "DeleteRequests"},
		{"object select", "POST", "/bucket/docs/a.txt?select&select-type=SQL", "SelectRequests"},
		{"multipart create", "POST", "/bucket/docs/a.txt?uploads", "PostRequests"},
		{"multipart complete", "POST", "/bucket/docs/a.txt?uploadId=abc", "PostRequests"},
		{"restore", "POST", "/bucket/docs/a.txt?restore", "PostRequests"},
		{"form post", "POST", "/bucket/docs/a.txt", "PostRequests"},
		{"bucket head", "HEAD", "/bucket", "HeadRequests"},
		{"delete objects", "POST", "/bucket?delete", "DeleteRequests"},
		{"list versions", "GET", "/bucket?versions", "ListRequests"},
		{"list uploads", "GET", "/bucket?uploads", "ListRequests"},
		{"list v2", "GET", "/bucket?list-type=2", "ListRequests"},
		{"list v1", "GET", "/bucket", "ListRequests"},
		{"bucket acl", "GET", "/bucket?acl", ""},
		{"bucket versioning put", "PUT", "/bucket?versioning", ""},
		{"create bucket", "PUT", "/bucket", ""},
		{"delete bucket", "DELETE", "/bucket", ""},
	}
	for _, tc := range cases {
		u, err := url.ParseRequestURI(tc.target)
		if err != nil {
			t.Fatalf("%s: bad target %q: %v", tc.name, tc.target, err)
		}
		r := &http.Request{Method: tc.method, URL: u}
		key := ""
		if rest := strings.TrimPrefix(u.Path, "/"); strings.Contains(rest, "/") {
			key = rest[strings.Index(rest, "/")+1:]
		}
		if got := classifyRequestMetrics(r, key); got != tc.want {
			t.Fatalf("%s: classifyRequestMetrics(%s %s) = %q, want %q", tc.name, tc.method, tc.target, got, tc.want)
		}
	}
}

// TestRequestLatencyPhases pins the two latency phases: FirstByteLatency
// accumulates the time to the response starting to be returned, while
// TotalRequestLatency additionally covers the response body, so a request
// with a body phase keeps the two metrics distinct.
func TestRequestLatencyPhases(t *testing.T) {
	a := &requestMetricsAggregator{
		window:   time.Minute,
		counters: map[requestMetricsKey]*requestMetricTotals{},
		configs:  map[string]cachedMetricsConfigs{},
	}
	a.accumulate(s3RequestObservation{
		region: "us-east-1", bucket: "bucket", key: "docs/a.txt",
		method: "GET", bytes: 10, status: 200,
		firstByte: 2 * time.Millisecond, total: 5 * time.Millisecond,
	}, "filter-a")
	a.accumulate(s3RequestObservation{
		region: "us-east-1", bucket: "bucket", key: "docs/a.txt",
		method: "HEAD", status: 200,
		firstByte: time.Millisecond, total: 1500 * time.Microsecond,
	}, "filter-a")

	totals := a.counters[requestMetricsKey{region: "us-east-1", bucket: "bucket", id: "filter-a"}]
	first, total := totals.dists["FirstByteLatency"], totals.dists["TotalRequestLatency"]
	if first == nil || total == nil {
		t.Fatal("both latency metrics must accumulate")
	}
	if first.sum >= total.sum {
		t.Fatalf("first-byte sum %v must be below total sum %v when responses carry a body phase", first.sum, total.sum)
	}
	if first.max >= total.max {
		t.Fatalf("first-byte max %v must be below total max %v", first.max, total.max)
	}
	if first.count != 2 || total.count != 2 {
		t.Fatalf("sample counts = %v/%v, want 2/2", first.count, total.count)
	}
}

// TestCopySourceObservation pins the CopyObject source contribution: the
// copy counts PutRequests (and AllRequests and the rest) for the
// destination, while the source bucket receives GetRequests only.
func TestCopySourceObservation(t *testing.T) {
	a := &requestMetricsAggregator{
		window:   time.Minute,
		counters: map[requestMetricsKey]*requestMetricTotals{},
		configs:  map[string]cachedMetricsConfigs{},
	}
	a.accumulate(s3RequestObservation{
		region: "us-east-1", bucket: "dst", key: "copy.txt",
		method: "PUT", op: "PutRequests", status: 200,
		firstByte: time.Millisecond, total: 2 * time.Millisecond,
	}, "dst-filter")
	a.accumulate(s3RequestObservation{
		region: "us-east-1", bucket: "src", key: "orig.txt",
		metricOnly: "GetRequests",
	}, "src-filter")

	dst := a.counters[requestMetricsKey{region: "us-east-1", bucket: "dst", id: "dst-filter"}]
	if dst == nil || dst.dists["AllRequests"] == nil || dst.dists["AllRequests"].sum != 1 || dst.dists["PutRequests"] == nil {
		t.Fatalf("destination totals = %+v, want AllRequests 1 plus PutRequests", dst)
	}
	src := a.counters[requestMetricsKey{region: "us-east-1", bucket: "src", id: "src-filter"}]
	if src == nil {
		t.Fatal("source totals missing")
	}
	gets := src.dists["GetRequests"]
	if gets == nil || gets.sum != 1 {
		t.Fatalf("source GetRequests = %+v, want sum 1", gets)
	}
	for metric := range src.dists {
		if metric != "GetRequests" {
			t.Fatalf("source received %q; the definition claims GetRequests only", metric)
		}
	}
}

// TestRequestMetricStatisticDatapoints pins the published datapoint shape:
// one statistic set per metric with Sum as the value, so the CloudWatch side
// can store Average-per-request semantics.
func TestRequestMetricStatisticDatapoints(t *testing.T) {
	totals := newMetricTotals()
	add := func(name string, v float64) { totals.dists[name].add(v) }
	totals.dists["4xxErrors"] = &requestMetricDistribution{}
	add("4xxErrors", 1)
	add("4xxErrors", 0)
	totals.dists["AllRequests"] = &requestMetricDistribution{}
	add("AllRequests", 1)
	add("AllRequests", 1)

	key := requestMetricsKey{region: "us-east-1", bucket: "bucket", id: "whole"}
	datapoints := statisticDatapoints(key, totals, time.Unix(1700000000, 0).UTC())
	if len(datapoints) != 2 {
		t.Fatalf("datapoint count = %d, want 2", len(datapoints))
	}
	byName := map[string]eventbus.S3MetricDatapoint{}
	for _, dp := range datapoints {
		byName[dp.MetricName] = dp
	}
	fourxx, ok := byName["4xxErrors"]
	if !ok {
		t.Fatal("4xxErrors datapoint missing")
	}
	if fourxx.Value != 1 || fourxx.SampleCount != 2 || fourxx.Minimum != 0 || fourxx.Maximum != 1 {
		t.Fatalf("4xxErrors = %+v, want sum 1 sample 2 min 0 max 1", fourxx)
	}
	if fourxx.Unit != "Count" {
		t.Fatalf("4xxErrors unit = %q", fourxx.Unit)
	}
	all, ok := byName["AllRequests"]
	if !ok {
		t.Fatal("AllRequests datapoint missing")
	}
	if all.Value != 2 || all.SampleCount != 2 || all.Minimum != 1 || all.Maximum != 1 {
		t.Fatalf("AllRequests = %+v, want sum 2 sample 2 min 1 max 1", all)
	}
	for _, dp := range datapoints {
		if len(dp.Dimensions) != 2 {
			t.Fatalf("%s dimensions = %+v, want BucketName and FilterId", dp.MetricName, dp.Dimensions)
		}
	}
}

package s3

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/request"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	s3store "vorpalstacks/internal/store/aws/s3"
)

// S3 request metrics reproduce the AWS request-metric set for the AWS/S3
// namespace: per-request counters, byte counts, error counters, and
// latencies, aggregated per metrics-configuration filter and published with
// the BucketName and FilterId dimensions CloudWatch documents.
const requestMetricsNamespace = "AWS/S3"

// requestMetricsUnit maps a metric name to its documented CloudWatch unit.
func requestMetricsUnit(name string) string {
	switch name {
	case "BytesDownloaded", "BytesUploaded":
		return "Bytes"
	case "FirstByteLatency", "TotalRequestLatency":
		return "Milliseconds"
	default:
		return "Count"
	}
}

// classifyRequestMetrics maps one served request to its per-operation
// metric following the AWS metric definitions: AllRequests counts every
// HTTP request made to the bucket regardless of type; GetRequests is made
// for objects and excludes list operations ("This metric is incremented
// for the source of each CopyObject request" is added separately);
// HeadRequests covers the HEAD requests made to an Amazon S3 bucket,
// object heads included; DeleteRequests "also includes DeleteObjects
// requests"; PostRequests excludes DeleteObjects and SelectObjectContent,
// which classify as DeleteRequests and SelectRequests; ListRequests covers
// the requests that list the contents of a bucket. A bucket-subresource
// request (GetBucketAcl, PutBucketVersioning, CreateBucket, ...) counts in
// AllRequests only: an empty result never suppresses the AllRequests count.
func classifyRequestMetrics(r *http.Request, key string) string {
	query := r.URL.Query()
	if key != "" {
		switch r.Method {
		case http.MethodPost:
			if query.Has("select") {
				return "SelectRequests"
			}
			return "PostRequests"
		case http.MethodHead:
			return "HeadRequests"
		case http.MethodDelete:
			return "DeleteRequests"
		case http.MethodPut:
			return "PutRequests"
		case http.MethodGet:
			// ListParts paginates an in-progress upload and is one of the
			// list-oriented requests GetRequests excludes.
			if query.Has("uploadId") {
				return "ListRequests"
			}
			return "GetRequests"
		}
		return ""
	}
	switch r.Method {
	case http.MethodHead:
		return "HeadRequests"
	case http.MethodPost:
		if query.Has("delete") {
			return "DeleteRequests"
		}
		return ""
	case http.MethodGet:
		if isBucketListRequest(r) {
			return "ListRequests"
		}
		return ""
	}
	return ""
}

// s3RequestObservation is one observed request on the S3 request plane.
// metricOnly marks a partial contribution: a CopyObject request increments
// GetRequests for the source bucket, and nothing else — the definition
// claims only GetRequests for the source.
type s3RequestObservation struct {
	region     string
	bucket     string
	key        string
	method     string
	op         string
	metricOnly string
	bytes      int64
	status     int
	firstByte  time.Duration
	total      time.Duration
}

// requestMetricsAggregator accumulates per-filter request metrics and
// publishes them to the event bus on a fixed window; the CloudWatch service
// ingests the published datapoints into its metric store.
type requestMetricsAggregator struct {
	svc     *S3Service
	window  time.Duration
	started sync.Once

	mu       sync.Mutex
	counters map[requestMetricsKey]*requestMetricTotals
	configs  map[string]cachedMetricsConfigs
}

type requestMetricsKey struct {
	region string
	bucket string
	id     string
}

// requestMetricDistribution is one metric's running statistic set. Every
// matching request contributes one value (1 for request counters, a 0-or-1
// for the error counters, bytes or milliseconds elsewhere), so after window
// aggregation Average = Sum/SampleCount keeps the documented per-request
// meaning — the error rate for 4xxErrors/5xxErrors, bytes per request for
// the byte counters.
type requestMetricDistribution struct {
	sum   float64
	min   float64
	max   float64
	count float64
}

func (d *requestMetricDistribution) add(value float64) {
	d.sum += value
	d.count++
	if d.count == 1 || value < d.min {
		d.min = value
	}
	if d.count == 1 || value > d.max {
		d.max = value
	}
}

// requestMetricTotals holds one filter's per-metric distributions.
type requestMetricTotals struct {
	dists map[string]*requestMetricDistribution
}

type cachedMetricsConfigs struct {
	configs []*s3store.MetricsConfiguration
	expires time.Time
}

func newRequestMetricsAggregator(svc *S3Service) *requestMetricsAggregator {
	// The compressed test window is the same test-runner-only mechanism the
	// inventory cadence uses.
	window := 60 * time.Second
	if os.Getenv("TEST_MODE") == "true" {
		window = 2 * time.Second
	}
	return &requestMetricsAggregator{
		svc:      svc,
		window:   window,
		counters: map[requestMetricsKey]*requestMetricTotals{},
		configs:  map[string]cachedMetricsConfigs{},
	}
}

func newMetricTotals() *requestMetricTotals {
	return &requestMetricTotals{dists: map[string]*requestMetricDistribution{}}
}

// bucketSubresourceParams are the bucket-subresource query parameters a GET
// can carry; a bucket GET with none of them is a ListObjects (V1) request.
var bucketSubresourceParams = []string{
	"acl", "versioning", "encryption", "policy", "policyStatus", "cors",
	"tagging", "lifecycle", "website", "replication", "object-lock",
	"notification", "logging", "ownershipControls", "requestPayment",
	"accelerate", "inventory", "metrics", "location", "publicAccessBlock",
}

func hasBucketSubresource(query url.Values) bool {
	for _, p := range bucketSubresourceParams {
		if query.Has(p) {
			return true
		}
	}
	return false
}

// isBucketListRequest reports whether a bucket-plane GET is a list request
// (?versions, ?uploads, list-type=2, or the ListObjects V1 fall-through).
func isBucketListRequest(r *http.Request) bool {
	if r.Method != http.MethodGet {
		return false
	}
	query := r.URL.Query()
	if query.Has("versions") || query.Has("uploads") || query.Get("list-type") == "2" {
		return true
	}
	return !hasBucketSubresource(query)
}

// recordServedRequest observes one request from the single ServeHTTP
// dispatch point. firstByte is the time from the complete request being
// received to the response starting to be returned (dispatch completion);
// total additionally covers writing the response body, matching the two
// documented latency definitions.
func (s *S3Service) recordServedRequest(reqCtx *request.RequestContext, r *http.Request, bucket, key string, result interface{}, status int, firstByte, total time.Duration) {
	a := s.requestMetrics
	if a == nil || bucket == "" {
		return
	}
	stores, err := s.store(reqCtx)
	if err != nil || stores == nil {
		return
	}
	bytes := int64(0)
	if cl := r.Header.Get("Content-Length"); cl != "" {
		if n, parseErr := strconv.ParseInt(cl, 10, 64); parseErr == nil && n > 0 {
			bytes = n
		}
	}
	if out, ok := result.(*GetObjectOutput); ok && out.ContentLength > 0 {
		bytes = out.ContentLength
	}
	a.record(s3RequestObservation{
		region:    reqCtx.GetRegion(),
		bucket:    bucket,
		key:       key,
		method:    r.Method,
		op:        classifyRequestMetrics(r, key),
		bytes:     bytes,
		status:    status,
		firstByte: firstByte,
		total:     total,
	}, stores)

	// A CopyObject request also increments GetRequests for the source
	// bucket, as the GetRequests definition documents. UploadPartCopy is
	// outside this rule: the definition names CopyObject only.
	if r.Method == http.MethodPut && key != "" && r.Header.Get("x-amz-copy-source") != "" && !r.URL.Query().Has("uploadId") {
		if srcBucket, srcKey, _, parseErr := parseCopySource(r.Header.Get("x-amz-copy-source")); parseErr == nil && srcBucket != "" {
			a.record(s3RequestObservation{
				region:     reqCtx.GetRegion(),
				bucket:     srcBucket,
				key:        srcKey,
				metricOnly: "GetRequests",
			}, stores)
		}
	}
}

func (a *requestMetricsAggregator) record(obs s3RequestObservation, stores *s3Stores) {
	configs := a.bucketConfigs(obs.bucket, stores)
	if len(configs) == 0 {
		return
	}
	a.startFlusher()
	for _, cfg := range configs {
		if !metricsFilterMatches(cfg.Filter, obs.key, func() []types.Tag {
			return a.objectTags(obs, stores)
		}) {
			continue
		}
		a.accumulate(obs, cfg.ID)
	}
}

// metricsFilterMatches evaluates one metrics filter against a request. The
// tags closure is consulted only for tag predicates, so buckets without tag
// filters never pay a tag lookup. An AccessPointArn predicate never matches:
// the platform has no access-point substrate, so no request can arrive
// through an access point and a filter whose only predicate is an access
// point ARN generates no datapoints.
func metricsFilterMatches(filter *s3store.MetricsFilter, key string, tags func() []types.Tag) bool {
	if filter == nil {
		return true
	}
	if filter.AccessPointArn != "" {
		return false
	}
	if filter.Prefix != "" && !strings.HasPrefix(key, filter.Prefix) {
		return false
	}
	if filter.Tag != nil && !metricsTagMatches(filter.Tag, tags()) {
		return false
	}
	if filter.And != nil {
		if filter.And.AccessPointArn != "" {
			return false
		}
		if filter.And.Prefix != "" && !strings.HasPrefix(key, filter.And.Prefix) {
			return false
		}
		if len(filter.And.Tags) > 0 {
			set := tags()
			for i := range filter.And.Tags {
				if !metricsTagMatches(&filter.And.Tags[i], set) {
					return false
				}
			}
		}
	}
	return true
}

func metricsTagMatches(tag *types.Tag, tags []types.Tag) bool {
	if tag == nil {
		return false
	}
	for _, candidate := range tags {
		if candidate.Key == tag.Key && candidate.Value == tag.Value {
			return true
		}
	}
	return false
}

// objectTags resolves the tag set of the observed object for tag predicates.
// A lookup failure yields no tags, which simply fails tag predicates.
func (a *requestMetricsAggregator) objectTags(obs s3RequestObservation, stores *s3Stores) []types.Tag {
	out, err := a.svc.getObjectTaggingCore(context.Background(), stores, &GetObjectTaggingInput{
		Bucket: obs.bucket,
		Key:    obs.key,
	})
	if err != nil || out == nil {
		return nil
	}
	tags := make([]types.Tag, 0, len(out.TagSet))
	for _, t := range out.TagSet {
		tags = append(tags, t.ToCommon())
	}
	return tags
}

// invalidateConfigs drops the cached metrics configurations of one bucket so
// a configuration write is visible to the aggregation immediately instead of
// after the cache TTL. Without this, the observations that precede and
// follow the write within the same window would keep aggregating against
// the stale pre-write (often empty) configuration set.
func (a *requestMetricsAggregator) invalidateConfigs(bucket string) {
	a.mu.Lock()
	delete(a.configs, bucket)
	a.mu.Unlock()
}

// bucketConfigs returns the bucket's metrics configurations through a TTL
// cache so the record path costs no store read on unconfigured buckets and
// at most one per window on configured ones. Configuration writes
// invalidate the entry (invalidateConfigs from the mutation cores).
func (a *requestMetricsAggregator) bucketConfigs(bucket string, stores *s3Stores) []*s3store.MetricsConfiguration {
	now := time.Now()
	a.mu.Lock()
	if cached, ok := a.configs[bucket]; ok && now.Before(cached.expires) {
		a.mu.Unlock()
		return cached.configs
	}
	a.mu.Unlock()

	configs := []*s3store.MetricsConfiguration{}
	if stored, err := stores.buckets.Get(bucket); err == nil && stored != nil {
		for _, cfg := range stored.MetricsConfigurations {
			configs = append(configs, cfg)
		}
	}
	a.mu.Lock()
	a.configs[bucket] = cachedMetricsConfigs{configs: configs, expires: now.Add(a.window)}
	a.mu.Unlock()
	return configs
}

func (a *requestMetricsAggregator) accumulate(obs s3RequestObservation, filterID string) {
	key := requestMetricsKey{region: obs.region, bucket: obs.bucket, id: filterID}
	a.mu.Lock()
	defer a.mu.Unlock()
	totals := a.counters[key]
	if totals == nil {
		totals = newMetricTotals()
		a.counters[key] = totals
	}
	// Every request contributes a 0-or-1 to both error counters, so their
	// SampleCount covers all requests and Average is the documented error
	// rate rather than the share of errored requests that happened to land.
	add := func(name string, value float64) {
		d := totals.dists[name]
		if d == nil {
			d = &requestMetricDistribution{}
			totals.dists[name] = d
		}
		d.add(value)
	}
	if obs.metricOnly != "" {
		add(obs.metricOnly, 1)
		return
	}
	add("AllRequests", 1)
	if obs.op != "" {
		add(obs.op, 1)
	}
	if obs.bytes > 0 {
		switch obs.method {
		case http.MethodGet:
			add("BytesDownloaded", float64(obs.bytes))
		case http.MethodPut:
			add("BytesUploaded", float64(obs.bytes))
		}
	}
	add("4xxErrors", boolValue(obs.status >= 400 && obs.status < 500))
	add("5xxErrors", boolValue(obs.status >= 500))
	firstByteMs := float64(obs.firstByte) / float64(time.Millisecond)
	totalMs := float64(obs.total) / float64(time.Millisecond)
	add("FirstByteLatency", firstByteMs)
	add("TotalRequestLatency", totalMs)
}

func boolValue(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// startFlusher launches the window publisher once; the goroutine lives for
// the process lifetime and only ever holds in-memory counters, so a hard
// shutdown loses at most one window of best-effort datapoints.
func (a *requestMetricsAggregator) startFlusher() {
	a.started.Do(func() {
		logs.Info("s3: request metrics aggregator started", logs.Any("window", a.window))
		go func() {
			ticker := time.NewTicker(a.window)
			defer ticker.Stop()
			for range ticker.C {
				a.flush()
			}
		}()
	})
}

// statisticDatapoints converts one filter's accumulated distributions into
// the datapoints published on the event bus: one CloudWatch statistic set
// per metric carrying Sum (as Value), SampleCount, Minimum, and Maximum.
func statisticDatapoints(key requestMetricsKey, totals *requestMetricTotals, now time.Time) []eventbus.S3MetricDatapoint {
	datapoints := make([]eventbus.S3MetricDatapoint, 0, len(totals.dists))
	for name, d := range totals.dists {
		datapoints = append(datapoints, eventbus.S3MetricDatapoint{
			MetricName:  name,
			Value:       d.sum,
			SampleCount: d.count,
			Minimum:     d.min,
			Maximum:     d.max,
			Unit:        requestMetricsUnit(name),
			Timestamp:   now,
			Dimensions: []eventbus.S3MetricDimension{
				{Name: "BucketName", Value: key.bucket},
				{Name: "FilterId", Value: key.id},
			},
		})
	}
	return datapoints
}

func (a *requestMetricsAggregator) flush() {
	a.mu.Lock()
	pending := a.counters
	a.counters = map[requestMetricsKey]*requestMetricTotals{}
	a.mu.Unlock()
	if len(pending) == 0 || a.svc == nil || a.svc.bus == nil {
		return
	}
	now := time.Now().UTC()
	byRegion := map[string][]eventbus.S3MetricDatapoint{}
	for key, totals := range pending {
		byRegion[key.region] = append(byRegion[key.region], statisticDatapoints(key, totals, now)...)
	}
	for region, datapoints := range byRegion {
		evt := &eventbus.S3RequestMetricsEvent{
			EventBase: eventbus.EventBase{
				Timestamp: now,
				Source:    "aws:s3",
				Region:    region,
				AccountID: a.svc.accountID,
			},
			Datapoints: datapoints,
		}
		if err := a.svc.bus.Publish(context.Background(), evt); err != nil {
			logs.Warn("s3: request metrics publish failed", logs.Err(err))
		}
	}
}

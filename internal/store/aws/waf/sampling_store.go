package waf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const sampledRequestsBucketName = "waf_sampled_requests"

// sampledCountsBucketName holds the per-minute matched-request
// counters that back the GetSampledRequests population figure.
const sampledCountsBucketName = "waf_sample_counts"

// sampleCountBucketWidth is the granularity of the matched-request
// counters: one counter per rule per minute keeps the key count within
// the three-hour retention bounded (180 keys per rule).
const sampleCountBucketWidth = time.Minute

// sampledRequestKeySuffixLen is the zero-padded width of the Unix-nano
// timestamp embedded in a sample key; fixed width keeps the byte order
// of keys equal to their chronological order.
const sampledRequestKeySuffixLen = 20

// SampledHTTPHeader is one forwarded header of a sampled web request.
type SampledHTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TokenInspectionRecord is the token-check outcome retained with a
// sampled request of a Captcha or Challenge rule: the kind's solve
// timestamp carried by the supplied token, and the reason the check
// failed when it did.
type TokenInspectionRecord struct {
	SolveTimestamp int64  `json:"solveTimestamp,omitempty"`
	FailureReason  string `json:"failureReason,omitempty"`
}

// SampledRequest is one sampled web request retained for the
// GetSampledRequests operation. RuleNameWithinRuleGroup carries the
// sampled rule name for matches from inside a rule group and stays
// empty for rules declared directly in the web ACL, where the sampled
// request carries no such name.
type SampledRequest struct {
	RuleNameWithinRuleGroup string              `json:"ruleNameWithinRuleGroup,omitempty"`
	MetricName              string              `json:"metricName"`
	Action                  string              `json:"action"`
	Timestamp               time.Time           `json:"timestamp"`
	ClientIP                string              `json:"clientIp"`
	URI                     string              `json:"uri"`
	Method                  string              `json:"method"`
	HTTPVersion             string              `json:"httpVersion"`
	Headers                 []SampledHTTPHeader `json:"headers"`
	// ResponseCodeSent is the status the enforcement plane sent for the
	// request the sample was taken from.
	ResponseCodeSent int `json:"responseCodeSent,omitempty"`
	// Labels are the fully qualified labels matching rules applied to
	// the request.
	Labels []string `json:"labels,omitempty"`
	// RequestHeadersInserted are the custom headers the matching rule
	// actions inserted into the forwarded request.
	RequestHeadersInserted []SampledHTTPHeader `json:"requestHeadersInserted,omitempty"`
	// OverriddenAction is the action a rule group rule was configured
	// with when a rule action override replaced it; empty when no
	// override applied.
	OverriddenAction string                 `json:"overriddenAction,omitempty"`
	Captcha          *TokenInspectionRecord `json:"captchaInspection,omitempty"`
	Challenge        *TokenInspectionRecord `json:"challengeInspection,omitempty"`
}

// SamplingStore persists sampled web requests in Pebble so samples
// survive a server restart within the three-hour retention window. One
// record is stored per key, ordered by the timestamp embedded in the
// key; per-rule depth is bounded to SamplingPopulationDepth by the
// retention sweep. A parallel counter bucket holds one matched-request
// count per rule per minute, which serves the GetSampledRequests
// population figure.
type SamplingStore struct {
	*common.BaseStore
	counts storage.Bucket
	mu     sync.Mutex
}

// NewSamplingStore creates a new sampled-request store.
func NewSamplingStore(store storage.BasicStorage) *SamplingStore {
	return &SamplingStore{
		BaseStore: common.NewBaseStore(store.Bucket(sampledRequestsBucketName), "waf"),
		counts:    store.Bucket(sampledCountsBucketName),
	}
}

func sampleRecordKey(webACLARN, metricName string, timestamp time.Time) string {
	return fmt.Sprintf("%s|%s|%0*d", webACLARN, metricName, sampledRequestKeySuffixLen, timestamp.UnixNano())
}

// Record stores one sampled web request. Two samples of the same rule
// landing on the identical nanosecond share a key and the earlier
// record is replaced; sampling is telemetry, so the loss is acceptable.
// The rule's matched-request counter for the current minute increments
// alongside, whatever the retention depth holds.
func (s *SamplingStore) Record(webACLARN, metricName string, record SampledRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := sampleRecordKey(webACLARN, metricName, record.Timestamp)
	if err := s.BaseStore.Put(key, record); err != nil {
		return NewStoreError("record_sampled_request", err)
	}
	if err := s.bumpCount(webACLARN, metricName, record.Timestamp); err != nil {
		return NewStoreError("record_sampled_request", err)
	}
	return nil
}

// sampleCountKey identifies one rule's counter for the minute that
// contains the timestamp.
func sampleCountKey(webACLARN, metricName string, ts time.Time) []byte {
	minute := ts.Unix() / int64(sampleCountBucketWidth/time.Second)
	return []byte(fmt.Sprintf("%s|%s|%020d", webACLARN, metricName, minute))
}

// bumpCount increments the rule's matched-request counter of the
// minute containing the timestamp. Callers hold s.mu. A missing
// counter reads as zero (the bucket interface returns a nil value for
// an absent key).
func (s *SamplingStore) bumpCount(webACLARN, metricName string, ts time.Time) error {
	key := sampleCountKey(webACLARN, metricName, ts)
	current := int64(0)
	raw, err := s.counts.Get(key)
	if err != nil {
		return err
	}
	if raw != nil {
		if parsed, err := strconv.ParseInt(string(raw), 10, 64); err == nil {
			current = parsed
		}
	}
	return s.counts.Put(key, []byte(strconv.FormatInt(current+1, 10)))
}

// CountPopulation returns the number of matched requests recorded for
// the rule inside the half-open window [start, end). GetSampledRequests
// samples from among the first 5,000 of them, so the population the
// callers report caps at SamplingPopulationDepth.
func (s *SamplingStore) CountPopulation(webACLARN, metricName string, start, end time.Time) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	prefix := []byte(webACLARN + "|" + metricName + "|")
	iter := s.counts.ScanPrefix(prefix)
	defer iter.Close()

	var total int64
	for iter.Next() {
		minute, err := strconv.ParseInt(string(iter.Value()), 10, 64)
		if err != nil {
			continue
		}
		bucketStart := time.Unix(parseCountKeyMinute(iter.Key())*int64(sampleCountBucketWidth/time.Second), 0)
		bucketEnd := bucketStart.Add(sampleCountBucketWidth)
		if bucketStart.Before(end) && bucketEnd.After(start) {
			total += minute
		}
	}
	return total
}

// parseCountKeyMinute decodes the minute suffix of a counter key; a
// malformed suffix decodes to 0 and is therefore never counted.
func parseCountKeyMinute(key []byte) int64 {
	sep := bytes.LastIndexByte(key, '|')
	if sep < 0 {
		return 0
	}
	minute, err := strconv.ParseInt(string(key[sep+1:]), 10, 64)
	if err != nil {
		return 0
	}
	return minute
}

// Query returns the sampled requests of the rule recorded inside the
// half-open time window [start, end), newest first, up to max results.
// Callers validate the window against the three-hour retention bound;
// records that outlived the bound may still occupy keys but can never
// be selected by a valid window.
func (s *SamplingStore) Query(webACLARN, metricName string, start, end time.Time, max int) ([]SampledRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prefix := []byte(webACLARN + "|" + metricName + "|")
	iter := s.Bucket().ScanPrefix(prefix)
	defer iter.Close()

	var inWindow []SampledRequest
	for iter.Next() {
		var record SampledRequest
		if err := json.Unmarshal(iter.Value(), &record); err != nil {
			return nil, NewStoreError("query_sampled_requests", err)
		}
		if record.Timestamp.Before(start) || !record.Timestamp.Before(end) {
			continue
		}
		inWindow = append(inWindow, record)
	}
	if err := iter.Error(); err != nil {
		return nil, NewStoreError("query_sampled_requests", err)
	}

	// Iteration is key-ascending, which is oldest-first; keep only the
	// newest max records and return them newest-first.
	if len(inWindow) > max {
		inWindow = inWindow[len(inWindow)-max:]
	}
	out := make([]SampledRequest, 0, len(inWindow))
	for i := len(inWindow) - 1; i >= 0; i-- {
		out = append(out, inWindow[i])
	}
	return out, nil
}

// PurgeExpired drops records older than the three-hour retention window
// and trims every rule's retained set to the SamplingPopulationDepth
// newest records; the matched-request counters outside the retention
// window are dropped with them. The retention sweep calls this
// periodically; between sweeps Query bounds every answer regardless.
func (s *SamplingStore) PurgeExpired(now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ruleKey → ascending record keys of that rule.
	perRule := make(map[string][]string)
	if err := s.Bucket().ForEach(func(k, _ []byte) error {
		key := string(k)
		ruleKey, _, ok := splitSampleRecordKey(key)
		if ok {
			perRule[ruleKey] = append(perRule[ruleKey], key)
		}
		return nil
	}); err != nil {
		return NewStoreError("purge_sampled_requests", err)
	}

	cutoffNano := now.Add(-SampleRetention).UnixNano()
	for _, keys := range perRule {
		// Keys are ascending, so expired and over-depth records are a
		// prefix of the slice.
		keepFrom := 0
		for keepFrom < len(keys) && sampleKeyTimestampNano(keys[keepFrom]) < cutoffNano {
			keepFrom++
		}
		if depth := len(keys) - keepFrom; depth > SamplingPopulationDepth {
			keepFrom += depth - SamplingPopulationDepth
		}
		for _, key := range keys[:keepFrom] {
			if err := s.Bucket().Delete([]byte(key)); err != nil {
				return NewStoreError("purge_sampled_requests", err)
			}
		}
	}

	if err := s.purgeExpiredCounts(now); err != nil {
		return NewStoreError("purge_sampled_requests", err)
	}
	return nil
}

// purgeExpiredCounts drops matched-request counters whose minute lies
// entirely outside the three-hour retention window. Callers hold s.mu.
func (s *SamplingStore) purgeExpiredCounts(now time.Time) error {
	cutoff := now.Add(-SampleRetention)
	iter := s.counts.ScanPrefix(nil)
	defer iter.Close()
	for iter.Next() {
		minute := parseCountKeyMinute(iter.Key())
		bucketEnd := time.Unix(minute*int64(sampleCountBucketWidth/time.Second), 0).Add(sampleCountBucketWidth)
		if !bucketEnd.After(cutoff) {
			if err := s.counts.Delete(iter.Key()); err != nil {
				return err
			}
		}
	}
	return iter.Error()
}

// splitSampleRecordKey splits a record key into its rule part
// (web ACL ARN and metric name) and its timestamp suffix. A key
// without the expected shape is reported as not ok and left untouched
// by the sweep.
func splitSampleRecordKey(key string) (ruleKey string, timestamp string, ok bool) {
	sep := strings.LastIndexByte(key, '|')
	if sep < 0 || len(key)-sep-1 != sampledRequestKeySuffixLen {
		return "", "", false
	}
	return key[:sep], key[sep+1:], true
}

// sampleKeyTimestampNano decodes the timestamp suffix of a record key;
// malformed suffixes decode to 0 and are therefore always purged.
func sampleKeyTimestampNano(key string) int64 {
	_, timestamp, ok := splitSampleRecordKey(key)
	if !ok {
		return 0
	}
	nano, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 0
	}
	return nano
}

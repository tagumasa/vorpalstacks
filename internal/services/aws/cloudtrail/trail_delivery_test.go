package cloudtrail

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below cover the trail delivery engine: the documented file
// naming, the record-content log file, the selector matrix, the digest
// chain with its verifiable signature, and the GetTrailStatus bookkeeping
// written under the trail lock.

func TestTrailDeliveryNaming(t *testing.T) {
	at := time.Date(2026, 9, 17, 14, 52, 7, 0, time.UTC)

	logRe := regexp.MustCompile(`^111122223333_CloudTrail_us-east-2_20260917T1452Z_[0-9A-Za-z]{16}\.json\.gz$`)
	name := logFileName("111122223333", "CloudTrail", "us-east-2", at, "Mu0KsOhtH1ar15ZZ")
	if !logRe.MatchString(name) {
		t.Fatalf("log file name %q does not match the documented format", name)
	}
	suffix := randomLogFileNameSuffix()
	if len(suffix) != logFileUniqueStringLength {
		t.Fatalf("unique suffix %q is %d characters, want %d", suffix, len(suffix), logFileUniqueStringLength)
	}

	digest := digestFileName("111122223333", "us-east-2", "my-trail", "us-east-2", at)
	want := "111122223333_CloudTrail-Digest_us-east-2_my-trail_us-east-2_20260917T145207Z.json.gz"
	if digest != want {
		t.Fatalf("digest file name = %q, want %q", digest, want)
	}

	key := trailObjectKey("ct-prefix", "111122223333", "CloudTrail", "us-east-2", at)
	wantKey := "ct-prefix/AWSLogs/111122223333/CloudTrail/us-east-2/2026/09/17/"
	if key != wantKey {
		t.Fatalf("object key = %q, want %q", key, wantKey)
	}
	bare := trailObjectKey("", "111122223333", "CloudTrail-Digest", "us-east-2", at)
	if !strings.HasPrefix(bare, "AWSLogs/") {
		t.Fatalf("unprefixed key %q must start at AWSLogs/", bare)
	}
}

func TestTrailSelectorMatching(t *testing.T) {
	readEvent := &cloudtrailstore.Event{EventName: "ListTrails", EventSource: "cloudtrail.amazonaws.com", ReadOnly: "true", EventCategory: "Management"}
	writeEvent := &cloudtrailstore.Event{EventName: "CreateTrail", EventSource: "cloudtrail.amazonaws.com", ReadOnly: "false", EventCategory: "Management"}

	// No selectors configured: the default selector matches everything.
	if !trailSelectorMatches(&cloudtrailstore.Trail{}, readEvent) {
		t.Fatal("default selector rejected a management event")
	}

	readOnly := &cloudtrailstore.Trail{EventSelectors: []cloudtrailstore.EventSelector{{ReadWriteType: "ReadOnly", IncludeManagementEvents: true}}}
	if !trailSelectorMatches(readOnly, readEvent) {
		t.Fatal("ReadOnly selector rejected a read event")
	}
	if trailSelectorMatches(readOnly, writeEvent) {
		t.Fatal("ReadOnly selector accepted a write event")
	}

	writeOnly := &cloudtrailstore.Trail{EventSelectors: []cloudtrailstore.EventSelector{{ReadWriteType: "WriteOnly", IncludeManagementEvents: true}}}
	if !trailSelectorMatches(writeOnly, writeEvent) {
		t.Fatal("WriteOnly selector rejected a write event")
	}
	if trailSelectorMatches(writeOnly, readEvent) {
		t.Fatal("WriteOnly selector accepted a read event")
	}

	noManagement := &cloudtrailstore.Trail{EventSelectors: []cloudtrailstore.EventSelector{{ReadWriteType: "All", IncludeManagementEvents: false}}}
	if trailSelectorMatches(noManagement, writeEvent) {
		t.Fatal("selector with IncludeManagementEvents=false accepted a management event")
	}

	excluded := &cloudtrailstore.Trail{EventSelectors: []cloudtrailstore.EventSelector{{
		ReadWriteType: "All", IncludeManagementEvents: true,
		ExcludeManagementEventSources: []string{"cloudtrail.amazonaws.com"},
	}}}
	if trailSelectorMatches(excluded, writeEvent) {
		t.Fatal("excluded event source was accepted")
	}

	equals := &cloudtrailstore.Trail{AdvancedEventSelectors: []cloudtrailstore.AdvancedEventSelector{{
		FieldSelectors: []cloudtrailstore.AdvancedFieldSelector{{Field: "eventCategory", Equals: []string{"Management"}}},
	}}}
	if !trailSelectorMatches(equals, writeEvent) {
		t.Fatal("advanced Management selector rejected a management event")
	}
	startsWith := &cloudtrailstore.Trail{AdvancedEventSelectors: []cloudtrailstore.AdvancedEventSelector{{
		FieldSelectors: []cloudtrailstore.AdvancedFieldSelector{{Field: "eventName", StartsWith: []string{"Create"}}},
	}}}
	if !trailSelectorMatches(startsWith, writeEvent) || trailSelectorMatches(startsWith, readEvent) {
		t.Fatal("startsWith selector matched the wrong events")
	}
	notEquals := &cloudtrailstore.Trail{AdvancedEventSelectors: []cloudtrailstore.AdvancedEventSelector{{
		FieldSelectors: []cloudtrailstore.AdvancedFieldSelector{{Field: "readOnly", NotEquals: []string{"true"}}},
	}}}
	if !trailSelectorMatches(notEquals, writeEvent) || trailSelectorMatches(notEquals, readEvent) {
		t.Fatal("readOnly NotEquals selector matched the wrong events")
	}
	unknownField := &cloudtrailstore.Trail{AdvancedEventSelectors: []cloudtrailstore.AdvancedEventSelector{{
		FieldSelectors: []cloudtrailstore.AdvancedFieldSelector{{Field: "eventVersion", Equals: []string{"1.08"}}},
	}}}
	if trailSelectorMatches(unknownField, writeEvent) {
		t.Fatal("a field the store does not carry must never hold")
	}
}

// newDeliveryTestStore returns a pebble-backed store with a logging trail.
func newDeliveryTestStore(t *testing.T, name, bucket string, validation bool) *cloudtrailstore.CloudTrailStore {
	t.Helper()
	store := newErrorTestServiceStore(t)
	trail := cloudtrailstore.NewTrail(name, bucket, "us-east-1")
	trail.LogFileValidationEnabled = validation
	if _, err := store.CreateTrail(trail); err != nil {
		t.Fatalf("create trail failed: %v", err)
	}
	if err := store.StartLogging(name); err != nil {
		t.Fatalf("start logging failed: %v", err)
	}
	return store
}

func putDeliveryEvent(t *testing.T, store *cloudtrailstore.CloudTrailStore, name string, at time.Time) {
	t.Helper()
	body, err := json.Marshal(map[string]interface{}{"eventID": name, "eventName": name, "eventCategory": "Management"})
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	if err := store.PutEvent(&cloudtrailstore.Event{
		EventID:         name,
		EventName:       name,
		ReadOnly:        "false",
		EventSource:     "cloudtrail.amazonaws.com",
		EventTime:       at,
		EventType:       "AwsApiCall",
		EventVersion:    "1.09",
		AwsRegion:       "us-east-1",
		EventCategory:   "Management",
		CloudTrailEvent: string(body),
	}); err != nil {
		t.Fatalf("put event %s: %v", name, err)
	}
}

// flushTrailForTest refreshes the trail record and flushes the trail's
// home-region window — the single-region form of the engine's per-region
// flush.
func flushTrailForTest(t *testing.T, svc *CloudTrailService, store *cloudtrailstore.CloudTrailStore, name string) {
	t.Helper()
	trail, err := store.GetTrail(name)
	if err != nil {
		t.Fatalf("get trail %s: %v", name, err)
	}
	svc.flushTrail(store, store, "us-east-1", trail, &trailDeliveryOutcome{})
}

func gunzip(t *testing.T, data []byte) []byte {
	t.Helper()
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("gunzip: %v", err)
	}
	return out
}

func TestFlushTrailDeliversLogFile(t *testing.T) {
	store := newDeliveryTestStore(t, "deliver-log", "logs-bucket", false)
	fake := newFakeS3Invoker("logs-bucket")
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake})

	// The second event's time is derived from the first: two consecutive
	// clock reads can land in the same millisecond, and equal times leave
	// the file order to the store's most-recent-first walk.
	firstAt := time.Now().UTC()
	putDeliveryEvent(t, store, "evt-deliver-1", firstAt)
	putDeliveryEvent(t, store, "evt-deliver-2", firstAt.Add(time.Millisecond))
	time.Sleep(25 * time.Millisecond)

	flushTrailForTest(t, svc, store, "deliver-log")

	if len(fake.buckets["logs-bucket"]) != 1 {
		t.Fatalf("expected exactly 1 delivered log file, got %d", len(fake.buckets["logs-bucket"]))
	}
	var key string
	for k := range fake.buckets["logs-bucket"] {
		key = k
	}
	nameRe := regexp.MustCompile(`^AWSLogs/acc123/CloudTrail/us-east-1/\d{4}/\d{2}/\d{2}/acc123_CloudTrail_us-east-1_\d{8}T\d{4}Z_[0-9A-Za-z]{16}\.json\.gz$`)
	if !nameRe.MatchString(key) {
		t.Fatalf("delivered key %q does not follow the documented naming", key)
	}

	var parsed struct {
		Records []map[string]interface{} `json:"Records"`
	}
	if err := json.Unmarshal(gunzip(t, fake.buckets["logs-bucket"][key]), &parsed); err != nil {
		t.Fatalf("log file is not {\"Records\":[...]} JSON: %v", err)
	}
	if len(parsed.Records) != 2 {
		t.Fatalf("log file carries %d records, want the 2 delivered events", len(parsed.Records))
	}
	if parsed.Records[0]["eventName"] != "evt-deliver-1" {
		t.Fatalf("record content mismatch: %v", parsed.Records[0])
	}

	trail, err := store.GetTrail("deliver-log")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if trail.LatestDeliveryTime == nil || trail.LatestDeliveryError != "" {
		t.Fatalf("delivery bookkeeping wrong: %+v", trail.LatestDeliveryTime)
	}
	if !trail.LatestDeliveryAttemptSuccess {
		t.Fatal("LatestDeliveryAttemptSucceeded must be true after delivery")
	}
	if wm := trail.DeliveryWatermarks["us-east-1"]; wm == 0 {
		t.Fatal("watermark not advanced after delivery")
	}

	// A second flush with no new events delivers nothing and stays clean.
	flushTrailForTest(t, svc, store, "deliver-log")
	if len(fake.buckets["logs-bucket"]) != 1 {
		t.Fatalf("empty window delivered another file: %d objects", len(fake.buckets["logs-bucket"]))
	}
}

func TestFlushTrailDeliversDigestWithChain(t *testing.T) {
	store := newDeliveryTestStore(t, "deliver-digest", "digest-bucket", true)
	fake := newFakeS3Invoker("digest-bucket")
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake})

	putDeliveryEvent(t, store, "evt-digest-1", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "deliver-digest")

	var digestKey, logKey string
	for k := range fake.buckets["digest-bucket"] {
		if strings.Contains(k, "CloudTrail-Digest/") {
			digestKey = k
		} else {
			logKey = k
		}
	}
	if digestKey == "" || logKey == "" {
		t.Fatalf("expected one log file and one digest, got keys %v", fake.buckets["digest-bucket"])
	}
	digestRe := regexp.MustCompile(`^AWSLogs/acc123/CloudTrail-Digest/us-east-1/\d{4}/\d{2}/\d{2}/acc123_CloudTrail-Digest_us-east-1_deliver-digest_us-east-1_\d{8}T\d{6}Z\.json\.gz$`)
	if !digestRe.MatchString(digestKey) {
		t.Fatalf("digest key %q does not follow the documented naming", digestKey)
	}

	body := gunzip(t, fake.buckets["digest-bucket"][digestKey])
	var digest map[string]interface{}
	if err := json.Unmarshal(body, &digest); err != nil {
		t.Fatalf("digest is not JSON: %v", err)
	}
	if digest["awsAccountId"] != "acc123" || digest["digestSignatureAlgorithm"] != "SHA256withRSA" {
		t.Fatalf("digest header fields wrong: %v", digest)
	}
	fingerprint, _ := digest["digestPublicKeyFingerprint"].(string)
	if len(fingerprint) != 32 {
		t.Fatalf("digestPublicKeyFingerprint = %q, want 32 hex characters", fingerprint)
	}
	if digest["previousDigestS3Object"] != nil {
		t.Fatal("the first digest is a starting digest: previous* members must be absent")
	}
	logFiles, _ := digest["logFiles"].([]interface{})
	if len(logFiles) != 1 {
		t.Fatalf("digest references %d log files, want the 1 delivered", len(logFiles))
	}
	entry := logFiles[0].(map[string]interface{})
	sum := sha256.Sum256(fake.buckets["digest-bucket"][logKey])
	if entry["hashValue"] != hex.EncodeToString(sum[:]) {
		t.Fatalf("digest hashValue does not match the delivered log file content")
	}

	// The signature verifies with the fingerprint-matched key: the
	// validator's path over the uncompressed digest bytes.
	pub, priv, err := store.LoadTrailSigningKey("deliver-digest")
	if err != nil {
		t.Fatalf("load signing key: %v", err)
	}
	if pub.Fingerprint() != fingerprint {
		t.Fatalf("fingerprint %q does not match the trail's key %q", fingerprint, pub.Fingerprint())
	}
	sig, err := hex.DecodeString(fake.metadata["digest-bucket"][digestKey+"/signature"])
	if err != nil {
		t.Fatalf("metadata signature is not hexadecimal: %v", err)
	}
	if fake.metadata["digest-bucket"][digestKey+"/signature-algorithm"] != "SHA256withRSA" {
		t.Fatal("signature-algorithm metadata missing")
	}
	digestSum := sha256.Sum256(body)
	if err := rsa.VerifyPKCS1v15(&priv.PublicKey, crypto.SHA256, digestSum[:], sig); err != nil {
		t.Fatalf("digest signature does not verify: %v", err)
	}
	if _, err := x509.ParsePKCS1PublicKey(pub.Value); err != nil {
		t.Fatalf("trail key is not served as PKCS#1: %v", err)
	}

	trail, _ := store.GetTrail("deliver-digest")
	if trail.LatestDigestTime == nil || len(trail.PendingDigestFiles) != 0 {
		t.Fatalf("digest bookkeeping wrong: time=%v pending=%d", trail.LatestDigestTime, len(trail.PendingDigestFiles))
	}
	firstDigestObject := digestKey

	// After the cadence elapses the next digest chains the previous one.
	if _, err := store.MutateTrail("deliver-digest", func(tr *cloudtrailstore.Trail) error {
		past := trail.LatestDigestTime.Add(-2 * time.Hour)
		tr.LastDigestEnd = &past
		return nil
	}); err != nil {
		t.Fatalf("rewind digest clock: %v", err)
	}
	putDeliveryEvent(t, store, "evt-digest-2", time.Now().UTC())
	// The digest file name carries seconds; the second digest must land
	// in a later second than the first or it overwrites it.
	time.Sleep(1100 * time.Millisecond)
	flushTrailForTest(t, svc, store, "deliver-digest")

	var secondDigestKey string
	for k := range fake.buckets["digest-bucket"] {
		if strings.Contains(k, "CloudTrail-Digest/") && k != firstDigestObject {
			secondDigestKey = k
		}
	}
	if secondDigestKey == "" {
		t.Fatal("second digest not delivered after the cadence elapsed")
	}
	var second map[string]interface{}
	if err := json.Unmarshal(gunzip(t, fake.buckets["digest-bucket"][secondDigestKey]), &second); err != nil {
		t.Fatalf("second digest is not JSON: %v", err)
	}
	if second["previousDigestS3Object"] != firstDigestObject {
		t.Fatalf("second digest chains %v, want %q", second["previousDigestS3Object"], firstDigestObject)
	}
	if second["previousDigestHashAlgorithm"] != "SHA-256" {
		t.Fatalf("previousDigestHashAlgorithm = %v, want SHA-256", second["previousDigestHashAlgorithm"])
	}
}

func TestFlushTrailRecordsDeliveryError(t *testing.T) {
	// The bucket is absent from the fake S3 surface: the delivery fails,
	// the failure is recorded, and the watermark stays so the window is
	// retried.
	store := newDeliveryTestStore(t, "deliver-fail", "missing-bucket", false)
	fake := newFakeS3Invoker()
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake})

	putDeliveryEvent(t, store, "evt-fail-1", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "deliver-fail")

	trail, err := store.GetTrail("deliver-fail")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if trail.LatestDeliveryError == "" {
		t.Fatal("delivery failure was not recorded on the trail")
	}
	if trail.LatestDeliveryAttemptSuccess {
		t.Fatal("LatestDeliveryAttemptSucceeded must be false after a failure")
	}
	if wm, ok := trail.DeliveryWatermarks["us-east-1"]; ok && wm != 0 {
		t.Fatal("watermark advanced past an undelivered window")
	}
}

func TestFlushTrailDeliversToCWLogsAndSNS(t *testing.T) {
	store := newDeliveryTestStore(t, "deliver-multi", "multi-bucket", false)

	const (
		groupARN = "arn:aws:logs:us-east-1:acc123:log-group:CloudTrail/logs"
		topicARN = "arn:aws:sns:us-east-1:acc123:notify"
	)
	if _, err := store.MutateTrail("deliver-multi", func(tr *cloudtrailstore.Trail) error {
		tr.CloudWatchLogsLogGroupARN = groupARN
		tr.SnsTopicARN = topicARN
		return nil
	}); err != nil {
		t.Fatalf("configure destinations: %v", err)
	}

	logs := newFakeLogsInvoker("CloudTrail/logs")
	sns := &fakeSNSInvoker{topics: map[string]string{topicARN: sufficientTopicPolicy(topicARN)}}
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: newFakeS3Invoker("multi-bucket"), sns: sns, logs: logs})

	firstAt := time.Now().UTC()
	putDeliveryEvent(t, store, "evt-multi-1", firstAt)
	putDeliveryEvent(t, store, "evt-multi-2", firstAt.Add(time.Millisecond))
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "deliver-multi")

	// CloudWatch Logs: the documented stream carries one entry per event,
	// each holding the event's record JSON at the event's own time.
	stream := "acc123_CloudTrail_us-east-1"
	entries := logs.events["CloudTrail/logs/"+stream]
	if len(entries) != 2 {
		t.Fatalf("stream %s carries %d entries, want the 2 delivered events", stream, len(entries))
	}
	if !json.Valid([]byte(entries[0].Message)) || !strings.Contains(entries[0].Message, "evt-multi") {
		t.Fatalf("entry message is not the event record JSON: %s", entries[0].Message)
	}
	if entries[0].Timestamp != firstAt.UnixMilli() {
		t.Fatalf("entry timestamp %d is not the event time %d", entries[0].Timestamp, firstAt.UnixMilli())
	}

	// The notification announces the delivered files as the documented
	// JSON object: the bucket and the delivered object keys.
	if len(sns.published) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(sns.published))
	}
	publish := sns.published[0]
	if publish.topicARN != topicARN {
		t.Fatalf("notification went to %s, want the trail's topic", publish.topicARN)
	}
	var message struct {
		S3Bucket    string   `json:"s3Bucket"`
		S3ObjectKey []string `json:"s3ObjectKey"`
	}
	if err := json.Unmarshal([]byte(publish.message), &message); err != nil {
		t.Fatalf("notification message is not the documented JSON: %v", err)
	}
	if message.S3Bucket != "multi-bucket" || len(message.S3ObjectKey) != 1 {
		t.Fatalf("notification content wrong: %+v", message)
	}
	if !strings.HasSuffix(message.S3ObjectKey[0], ".json.gz") {
		t.Fatalf("notification key is not a delivered log file: %s", message.S3ObjectKey[0])
	}

	trail, err := store.GetTrail("deliver-multi")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if trail.LatestCWLogsDeliveryTime == nil || trail.LatestCWLogsDeliveryError != "" {
		t.Fatal("cloudwatch-logs delivery bookkeeping wrong")
	}
	if trail.LatestNotificationTime == nil || trail.LatestNotificationError != "" || !trail.LatestNotificationAttemptSuccess {
		t.Fatal("notification bookkeeping wrong")
	}
}

func TestFlushTrailRecordsCWLogsError(t *testing.T) {
	// The configured log group does not exist: the log file still delivers,
	// the CloudWatch Logs failure is recorded on its own member, and the
	// watermark advances — the destinations report independently.
	store := newDeliveryTestStore(t, "deliver-cwl-fail", "cwl-bucket", false)
	if _, err := store.MutateTrail("deliver-cwl-fail", func(tr *cloudtrailstore.Trail) error {
		tr.CloudWatchLogsLogGroupARN = "arn:aws:logs:us-east-1:acc123:log-group:missing-group"
		return nil
	}); err != nil {
		t.Fatalf("configure destination: %v", err)
	}

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: newFakeS3Invoker("cwl-bucket"), logs: newFakeLogsInvoker()})

	putDeliveryEvent(t, store, "evt-cwl-fail", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "deliver-cwl-fail")

	trail, err := store.GetTrail("deliver-cwl-fail")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if trail.LatestCWLogsDeliveryTime != nil {
		t.Fatal("no cloudwatch-logs delivery time may be recorded for a failed destination")
	}
	if trail.LatestCWLogsDeliveryError == "" {
		t.Fatal("cloudwatch-logs delivery failure was not recorded")
	}
	if trail.LatestDeliveryTime == nil {
		t.Fatal("the S3 log file must still have delivered")
	}
}

func TestFlushTrailRecordsNotificationError(t *testing.T) {
	// The topic is gone at delivery time: the failure is recorded with its
	// attempt outcome, and the delivery pipeline is otherwise unaffected.
	store := newDeliveryTestStore(t, "deliver-notify-fail", "notify-bucket", false)
	if _, err := store.MutateTrail("deliver-notify-fail", func(tr *cloudtrailstore.Trail) error {
		tr.SnsTopicARN = "arn:aws:sns:us-east-1:acc123:vanished"
		return nil
	}); err != nil {
		t.Fatalf("configure destination: %v", err)
	}

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{
		s3:  newFakeS3Invoker("notify-bucket"),
		sns: &fakeSNSInvoker{topics: map[string]string{}},
	})

	putDeliveryEvent(t, store, "evt-notify-fail", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "deliver-notify-fail")

	trail, err := store.GetTrail("deliver-notify-fail")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if trail.LatestNotificationTime != nil {
		t.Fatal("no notification time may be recorded for a failed publish")
	}
	if trail.LatestNotificationError == "" || trail.LatestNotificationAttemptSuccess {
		t.Fatal("notification failure bookkeeping wrong")
	}
	if trail.LatestDeliveryTime == nil {
		t.Fatal("the S3 log file must still have delivered")
	}
}

// panicLookupStore forwards every store operation to the wrapped store
// except the event page walk, which panics — a read-path fault localised
// to the flush that walks it.
type panicLookupStore struct {
	*cloudtrailstore.CloudTrailStore
}

func (s *panicLookupStore) LookupEvents(query cloudtrailstore.EventQuery) ([]*cloudtrailstore.Event, string, error) {
	panic("injected store fault")
}

// TestFlushTrailPanicContained pins the containment contract: a panic in
// one trail's flush is recovered, the process survives, the healthy
// trail's flush completes, and the panicked trail records a delivery
// failure with its window held for the tick's retry.
func TestFlushTrailPanicContained(t *testing.T) {
	store := newDeliveryTestStore(t, "panic-good", "panic-bucket", false)
	bad := cloudtrailstore.NewTrail("panic-bad", "panic-bucket", "us-east-1")
	if _, err := store.CreateTrail(bad); err != nil {
		t.Fatalf("create trail failed: %v", err)
	}
	if err := store.StartLogging("panic-bad"); err != nil {
		t.Fatalf("start logging failed: %v", err)
	}
	putDeliveryEvent(t, store, "evt-panic-1", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)

	goodTrail, err := store.GetTrail("panic-good")
	if err != nil {
		t.Fatalf("get good trail: %v", err)
	}
	badTrail, err := store.GetTrail("panic-bad")
	if err != nil {
		t.Fatalf("get bad trail: %v", err)
	}
	panicking := &panicLookupStore{CloudTrailStore: store}
	fake := newFakeS3Invoker("panic-bucket")
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		svc.flushTrailGuarded(store, panicking, "us-east-1", badTrail)
	}()
	go func() {
		defer wg.Done()
		svc.flushTrailGuarded(store, store, "us-east-1", goodTrail)
	}()
	wg.Wait()

	if len(fake.buckets["panic-bucket"]) != 1 {
		t.Fatalf("the healthy trail's flush must complete, got %d objects", len(fake.buckets["panic-bucket"]))
	}
	recorded, err := store.GetTrail("panic-bad")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if recorded.LatestDeliveryError == "" || recorded.LatestDeliveryAttemptSuccess {
		t.Fatal("the panic outcome was not recorded as a delivery failure")
	}
	if wm, ok := recorded.DeliveryWatermarks["us-east-1"]; ok && wm != 0 {
		t.Fatal("watermark advanced past a panicked window")
	}
}

// TestFlushTrailPanicPreservesPartialDelivery pins the panic path's
// bookkeeping: a flush that already wrote its log file (and delivered the
// CloudWatch Logs leg) before panicking records those facts — the written
// file joins the pending digest set, the stream's watermark advances, and
// the S3 window is held for the tick's retry.
func TestFlushTrailPanicPreservesPartialDelivery(t *testing.T) {
	store := newDeliveryTestStore(t, "panic-partial", "partial-bucket", true)
	const groupARN = "arn:aws:logs:us-east-1:acc123:log-group:CloudTrail/logs"
	if _, err := store.MutateTrail("panic-partial", func(tr *cloudtrailstore.Trail) error {
		tr.CloudWatchLogsLogGroupARN = groupARN
		return nil
	}); err != nil {
		t.Fatalf("configure destination: %v", err)
	}
	putDeliveryEvent(t, store, "evt-panic-partial-1", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)

	logs := newFakeLogsInvoker("CloudTrail/logs")
	// The plain log-file write succeeds; the digest write (metadata put)
	// panics — the fault lands after the flush accumulated its partial
	// state.
	fake := newFakeS3Invoker("partial-bucket")
	fake.panicOnMetadata = true
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake, logs: logs})

	trail, err := store.GetTrail("panic-partial")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	svc.flushTrailGuarded(store, store, "us-east-1", trail)

	if len(fake.buckets["partial-bucket"]) != 1 {
		t.Fatalf("the pre-panic log file must exist in the bucket, got %d objects", len(fake.buckets["partial-bucket"]))
	}
	var writtenKey string
	for k := range fake.buckets["partial-bucket"] {
		writtenKey = k
	}
	stream := "CloudTrail/logs/acc123_CloudTrail_us-east-1"
	if len(logs.events[stream]) != 1 {
		t.Fatalf("the pre-panic cloudwatch-logs delivery must have run, got %d entries", len(logs.events[stream]))
	}

	recorded, err := store.GetTrail("panic-partial")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if recorded.LatestDeliveryError == "" || recorded.LatestDeliveryAttemptSuccess {
		t.Fatal("the panic outcome was not recorded as a delivery failure")
	}
	if len(recorded.PendingDigestFiles) != 1 || recorded.PendingDigestFiles[0].ObjectKey != writtenKey {
		t.Fatalf("the pre-panic file must join the pending digest set, got %+v", recorded.PendingDigestFiles)
	}
	if wm := recorded.CwlogsWatermarks["us-east-1"]; wm == 0 {
		t.Fatal("the succeeded cloudwatch-logs leg must keep its watermark advance")
	}
	if wm, ok := recorded.DeliveryWatermarks["us-east-1"]; ok && wm != 0 {
		t.Fatal("watermark advanced past a panicked window")
	}
}

// TestFlushTrailCWLogsIndependentWatermark pins the CloudWatch Logs leg's
// watermark independence: when the S3 write fails and the tick retries the
// window verbatim, the stream is not re-sent — only events after the
// stream's own watermark go out, and a later event arrives exactly once
// while the S3 window still carries the accumulated backlog.
func TestFlushTrailCWLogsIndependentWatermark(t *testing.T) {
	store := newDeliveryTestStore(t, "cwl-independent", "cwl-indep-bucket", false)
	const groupARN = "arn:aws:logs:us-east-1:acc123:log-group:CloudTrail/logs"
	if _, err := store.MutateTrail("cwl-independent", func(tr *cloudtrailstore.Trail) error {
		tr.CloudWatchLogsLogGroupARN = groupARN
		return nil
	}); err != nil {
		t.Fatalf("configure destination: %v", err)
	}

	logs := newFakeLogsInvoker("CloudTrail/logs")
	// The S3 surface lacks the trail's bucket: every S3 write fails while
	// the stream delivery succeeds.
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: newFakeS3Invoker(), logs: logs})

	putDeliveryEvent(t, store, "evt-cwl-indep-1", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "cwl-independent")

	stream := "CloudTrail/logs/acc123_CloudTrail_us-east-1"
	if len(logs.events[stream]) != 1 {
		t.Fatalf("stream carries %d entries, want the 1 window event", len(logs.events[stream]))
	}
	trail, err := store.GetTrail("cwl-independent")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if wm := trail.CwlogsWatermarks["us-east-1"]; wm == 0 {
		t.Fatal("cloudwatch-logs watermark not advanced after its delivery")
	}
	if trail.LatestDeliveryError == "" {
		t.Fatal("the S3 failure must still be recorded")
	}

	// The tick retries the failed S3 window verbatim; the stream must not
	// receive the retried events again.
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "cwl-independent")
	if len(logs.events[stream]) != 1 {
		t.Fatalf("S3 retry re-sent the stream: %d entries", len(logs.events[stream]))
	}

	// A later event reaches the stream exactly once even though the S3
	// window (still undelivered as a whole) covers both events.
	putDeliveryEvent(t, store, "evt-cwl-indep-2", time.Now().UTC())
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "cwl-independent")
	if entries := logs.events[stream]; len(entries) != 2 {
		t.Fatalf("stream carries %d entries after the new event, want 2", len(entries))
	}
}

// TestFlushTrailCWLogsFailureRedeliversWindow pins the stream leg's own
// recovery: when a window's CloudWatch Logs delivery fails after S3
// succeeded, the events are not lost — the stream watermark stays put and
// a later flush re-reads the span from it, delivering the missed events
// exactly once more.
func TestFlushTrailCWLogsFailureRedeliversWindow(t *testing.T) {
	store := newDeliveryTestStore(t, "cwl-recover", "cwl-recover-bucket", false)
	const groupARN = "arn:aws:logs:us-east-1:acc123:log-group:CloudTrail/logs"
	if _, err := store.MutateTrail("cwl-recover", func(tr *cloudtrailstore.Trail) error {
		tr.CloudWatchLogsLogGroupARN = groupARN
		return nil
	}); err != nil {
		t.Fatalf("configure destination: %v", err)
	}

	logs := newFakeLogsInvoker("CloudTrail/logs")
	logs.failPuts = 1
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: newFakeS3Invoker("cwl-recover-bucket"), logs: logs})

	stream := "CloudTrail/logs/acc123_CloudTrail_us-east-1"
	evtT := time.Now().UTC()
	putDeliveryEvent(t, store, "evt-cwl-recover-1", evtT)
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "cwl-recover")

	// S3 succeeded and the stream failed: the failure is reported and the
	// stream watermark pins to the failed span's start, never past the
	// events it still owes.
	trail, err := store.GetTrail("cwl-recover")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if trail.LatestCWLogsDeliveryError == "" {
		t.Fatal("the stream failure must be reported")
	}
	if wm, ok := trail.CwlogsWatermarks["us-east-1"]; ok && wm > evtT.UnixMilli() {
		t.Fatalf("the stream watermark (%d) advanced past the event it owes (%d)", wm, evtT.UnixMilli())
	}
	if trail.LatestDeliveryError != "" {
		t.Fatal("the S3 leg must not be held by the stream failure")
	}
	if len(logs.events[stream]) != 0 {
		t.Fatalf("stream carries %d entries after the injected failure, want 0", len(logs.events[stream]))
	}

	// The next flush has nothing new for S3 but re-reads the stream's own
	// undelivered span and delivers it.
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "cwl-recover")
	if entries := logs.events[stream]; len(entries) != 1 || !strings.Contains(entries[0].Message, "evt-cwl-recover-1") {
		t.Fatalf("stream carries %+v after recovery, want the missed event exactly once", logs.events[stream])
	}
	trail, err = store.GetTrail("cwl-recover")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if wm := trail.CwlogsWatermarks["us-east-1"]; wm == 0 {
		t.Fatal("the stream watermark must advance after the recovered delivery")
	}
	if trail.LatestCWLogsDeliveryError != "" {
		t.Fatal("the recovered delivery must clear the stream error")
	}
}

// TestDeliverTrailsAllRegionsMultiRegion pins the multi-region contract:
// a trail with IsMultiRegionTrail delivers every active region's events —
// the events read from each region's own store, the file paths and names
// carrying that region — while a single-region trail delivers only its
// home region. All bookkeeping stays on the trail record in its home
// store, keyed per region. The CloudWatch Logs destination receives every
// region's events in the one home-named stream of the group its ARN
// names.
func TestDeliverTrailsAllRegionsMultiRegion(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-multi-region-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: tmpDir})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })

	const home = "us-east-1"
	const away = "eu-west-1"
	homeStorage, err := sm.GetStorage(home)
	if err != nil {
		t.Fatalf("home storage: %v", err)
	}
	homeStore := cloudtrailstore.NewCloudTrailStore(homeStorage, "acc123", home)
	awayStorage, err := sm.GetStorage(away)
	if err != nil {
		t.Fatalf("away storage: %v", err)
	}
	awayStore := cloudtrailstore.NewCloudTrailStore(awayStorage, "acc123", away)

	const logGroup = "CloudTrail/logs"
	logs := newFakeLogsInvoker(logGroup)
	multi := cloudtrailstore.NewTrail("multi-region", "multi-bucket", home)
	multi.IsMultiRegionTrail = true
	multi.CloudWatchLogsLogGroupARN = fmt.Sprintf("arn:aws:logs:%s:acc123:log-group:%s", home, logGroup)
	if _, err := homeStore.CreateTrail(multi); err != nil {
		t.Fatalf("create multi-region trail: %v", err)
	}
	single := cloudtrailstore.NewTrail("single-region", "multi-bucket", home)
	if _, err := homeStore.CreateTrail(single); err != nil {
		t.Fatalf("create single-region trail: %v", err)
	}
	for _, name := range []string{"multi-region", "single-region"} {
		if err := homeStore.StartLogging(name); err != nil {
			t.Fatalf("start logging %s: %v", name, err)
		}
	}

	base := time.Now().UTC()
	putDeliveryEvent(t, homeStore, "evt-home-1", base)
	putDeliveryEvent(t, awayStore, "evt-away-1", base.Add(time.Millisecond))
	time.Sleep(25 * time.Millisecond)

	fake := newFakeS3Invoker("multi-bucket")
	svc := NewCloudTrailService("acc123", home)
	svc.SetStorageManager(sm)
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake, logs: logs})

	svc.deliverTrailsAllRegions()

	perRegion := map[string]int{}
	for key := range fake.buckets["multi-bucket"] {
		parts := strings.Split(key, "/")
		if len(parts) < 4 || parts[0] != "AWSLogs" {
			t.Fatalf("unexpected object key %q", key)
		}
		perRegion[parts[3]]++
	}
	if perRegion[away] != 1 {
		t.Fatalf("multi-region trail must deliver the away region's events, got %d eu-west-1 objects", perRegion[away])
	}
	if perRegion[home] != 2 {
		t.Fatalf("both trails must deliver the home region, got %d us-east-1 objects", perRegion[home])
	}

	// The one home-named stream carries both regions' events: the away
	// event reached the group alongside the home event, and no
	// away-named stream exists.
	homeStream := logGroup + "/acc123_CloudTrail_" + home
	entries := logs.events[homeStream]
	if len(entries) != 2 {
		t.Fatalf("home stream carries %d entries, want both regions' events", len(entries))
	}
	sawAway := false
	for _, e := range entries {
		if strings.Contains(e.Message, "evt-away-1") {
			sawAway = true
		}
	}
	if !sawAway {
		t.Fatal("the away region's event did not reach the home-named stream")
	}
	if _, ok := logs.events[logGroup+"/acc123_CloudTrail_"+away]; ok {
		t.Fatal("an away-named stream must not exist: the trail's one stream is home-named")
	}

	multiTrail, err := homeStore.GetTrail("multi-region")
	if err != nil {
		t.Fatalf("get multi-region trail: %v", err)
	}
	if wm := multiTrail.DeliveryWatermarks[away]; wm == 0 {
		t.Fatal("the away region's watermark must be recorded on the home trail record")
	}
	if wm := multiTrail.DeliveryWatermarks[home]; wm == 0 {
		t.Fatal("the home region's watermark must be recorded on the home trail record")
	}
	singleTrail, err := homeStore.GetTrail("single-region")
	if err != nil {
		t.Fatalf("get single-region trail: %v", err)
	}
	if wm, ok := singleTrail.DeliveryWatermarks[away]; ok && wm != 0 {
		t.Fatal("a single-region trail must not carry another region's watermark")
	}
}

// TestDeliverTrailsAllRegionsDigestChain pins the digest-chain invariant
// across a multi-region walk: the first flush digests exactly once, the
// region flushed after the digest leaves its log file pending, and the
// next digest chains the delivered one while covering the pending file —
// two digests sharing one predecessor (a forked chain) must never be
// written.
func TestDeliverTrailsAllRegionsDigestChain(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-digest-chain-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: tmpDir})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })

	const home = "us-east-1"
	const away = "eu-west-1"
	homeStorage, err := sm.GetStorage(home)
	if err != nil {
		t.Fatalf("home storage: %v", err)
	}
	homeStore := cloudtrailstore.NewCloudTrailStore(homeStorage, "acc123", home)
	awayStorage, err := sm.GetStorage(away)
	if err != nil {
		t.Fatalf("away storage: %v", err)
	}
	awayStore := cloudtrailstore.NewCloudTrailStore(awayStorage, "acc123", away)

	multi := cloudtrailstore.NewTrail("chain-region", "chain-bucket", home)
	multi.IsMultiRegionTrail = true
	multi.LogFileValidationEnabled = true
	if _, err := homeStore.CreateTrail(multi); err != nil {
		t.Fatalf("create multi-region trail: %v", err)
	}
	if err := homeStore.StartLogging("chain-region"); err != nil {
		t.Fatalf("start logging: %v", err)
	}

	base := time.Now().UTC()
	putDeliveryEvent(t, homeStore, "evt-chain-home-1", base)
	putDeliveryEvent(t, awayStore, "evt-chain-away-1", base.Add(time.Millisecond))
	time.Sleep(25 * time.Millisecond)

	fake := newFakeS3Invoker("chain-bucket")
	svc := NewCloudTrailService("acc123", home)
	svc.SetStorageManager(sm)
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake})

	svc.deliverTrailsAllRegions()

	splitDelivered := func() (logKeys, digestKeys []string) {
		for key := range fake.buckets["chain-bucket"] {
			if strings.Contains(key, "CloudTrail-Digest/") {
				digestKeys = append(digestKeys, key)
			} else {
				logKeys = append(logKeys, key)
			}
		}
		return logKeys, digestKeys
	}
	logKeys, digestKeys := splitDelivered()
	if len(digestKeys) != 1 {
		t.Fatalf("the first flush must deliver exactly one digest across every region, got %d: %v", len(digestKeys), digestKeys)
	}
	if len(logKeys) != 2 {
		t.Fatalf("both regions' windows must land as log files, got %d: %v", len(logKeys), logKeys)
	}
	var firstDigest map[string]interface{}
	if err := json.Unmarshal(gunzip(t, fake.buckets["chain-bucket"][digestKeys[0]]), &firstDigest); err != nil {
		t.Fatalf("first digest is not JSON: %v", err)
	}
	if covered, _ := firstDigest["logFiles"].([]interface{}); len(covered) != 1 {
		t.Fatalf("the first digest covers %d log files, want the one flushed before it", len(covered))
	}
	trail, err := homeStore.GetTrail("chain-region")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if len(trail.PendingDigestFiles) != 1 {
		t.Fatalf("the region flushed after the digest must leave its file pending, got %d", len(trail.PendingDigestFiles))
	}
	if trail.LatestDigestTime == nil || trail.PreviousDigestObject != digestKeys[0] {
		t.Fatalf("digest bookkeeping wrong: prev=%q", trail.PreviousDigestObject)
	}

	// After the cadence elapses the next digest chains the first and
	// covers the pending file alongside the new window's.
	if _, err := homeStore.MutateTrail("chain-region", func(tr *cloudtrailstore.Trail) error {
		past := time.Now().UTC().Add(-2 * time.Hour)
		tr.LastDigestEnd = &past
		return nil
	}); err != nil {
		t.Fatalf("rewind digest clock: %v", err)
	}
	putDeliveryEvent(t, homeStore, "evt-chain-home-2", time.Now().UTC())
	// The digest file name carries seconds; the next digest must land in
	// a later second than the first or it overwrites it.
	time.Sleep(1100 * time.Millisecond)
	svc.deliverTrailsAllRegions()

	_, digestKeys = splitDelivered()
	if len(digestKeys) != 2 {
		t.Fatalf("the second flush must deliver exactly one more digest, got %d: %v", len(digestKeys), digestKeys)
	}
	firstKey := trail.PreviousDigestObject
	secondKey := ""
	for _, k := range digestKeys {
		if k != firstKey {
			secondKey = k
		}
	}
	if secondKey == "" {
		t.Fatal("the second digest did not land as a new object")
	}
	var secondDigest map[string]interface{}
	if err := json.Unmarshal(gunzip(t, fake.buckets["chain-bucket"][secondKey]), &secondDigest); err != nil {
		t.Fatalf("second digest is not JSON: %v", err)
	}
	if secondDigest["previousDigestS3Object"] != trail.PreviousDigestObject {
		t.Fatalf("second digest chains %v, want the first digest %q",
			secondDigest["previousDigestS3Object"], trail.PreviousDigestObject)
	}
	if covered, _ := secondDigest["logFiles"].([]interface{}); len(covered) != 2 {
		t.Fatalf("the second digest covers %d log files, want the pending one and the new window's", len(covered))
	}
	final, err := homeStore.GetTrail("chain-region")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if len(final.PendingDigestFiles) != 0 {
		t.Fatalf("pending files must be consumed by the chaining digest, got %d", len(final.PendingDigestFiles))
	}
}

// TestFlushTrailSplitsOversizedWindow pins the window batching: a window
// carrying more than MaxTrailLogFileEvents matched events lands as several
// log files, each bounded by MaxTrailLogFileEvents records, with the
// watermark advancing over the whole window.
func TestFlushTrailSplitsOversizedWindow(t *testing.T) {
	store := newDeliveryTestStore(t, "split-window", "split-bucket", false)
	fake := newFakeS3Invoker("split-bucket")
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: fake})

	const total = MaxTrailLogFileEvents + 1
	// The event times stay inside a few milliseconds: the window's upper
	// bound is the flush's own start time, so a spread wider than the
	// insertion takes would leave the tail events outside the window.
	base := time.Now().UTC()
	for i := 0; i < total; i++ {
		putDeliveryEvent(t, store, fmt.Sprintf("evt-split-%d", i), base.Add(time.Duration(i)*time.Microsecond))
	}
	time.Sleep(25 * time.Millisecond)
	flushTrailForTest(t, svc, store, "split-window")

	objects := fake.buckets["split-bucket"]
	if len(objects) != 2 {
		t.Fatalf("a %d-event window must split into 2 log files, got %d", total, len(objects))
	}
	records := 0
	for _, data := range objects {
		var parsed struct {
			Records []map[string]interface{} `json:"Records"`
		}
		if err := json.Unmarshal(gunzip(t, data), &parsed); err != nil {
			t.Fatalf("log file is not JSON: %v", err)
		}
		if len(parsed.Records) > MaxTrailLogFileEvents {
			t.Fatalf("log file carries %d records, over the %d bound", len(parsed.Records), MaxTrailLogFileEvents)
		}
		records += len(parsed.Records)
	}
	if records != total {
		t.Fatalf("the files together carry %d records, want all %d", records, total)
	}
	trail, err := store.GetTrail("split-window")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}
	if wm := trail.DeliveryWatermarks["us-east-1"]; wm == 0 {
		t.Fatal("watermark not advanced over the split window")
	}
}

func TestTakePutLogEventsBatchBounds(t *testing.T) {
	entries := make([]invokers.LogsLogEntry, invokers.PutLogEventsBatchMaxEvents+5)
	for i := range entries {
		entries[i] = invokers.LogsLogEntry{Timestamp: int64(i), Message: "x"}
	}
	batch, consumed := takePutLogEventsBatch(entries)
	if len(batch) != invokers.PutLogEventsBatchMaxEvents || consumed != invokers.PutLogEventsBatchMaxEvents {
		t.Fatalf("cardinality bound: batch %d consumed %d, want %d", len(batch), consumed, invokers.PutLogEventsBatchMaxEvents)
	}

	// Two messages that each fit alone but not together: the second opens
	// the next batch rather than overflowing this one.
	pair := []invokers.LogsLogEntry{
		{Timestamp: 1, Message: strings.Repeat("a", 600000)},
		{Timestamp: 2, Message: strings.Repeat("b", 600000)},
	}
	batch, consumed = takePutLogEventsBatch(pair)
	if len(batch) != 1 || consumed != 1 {
		t.Fatalf("byte bound: batch %d consumed %d, want 1/1", len(batch), consumed)
	}

	// A message over the per-event bound can never fit any batch: it is
	// skipped (consumed) without entering the batch.
	over := []invokers.LogsLogEntry{
		{Timestamp: 1, Message: strings.Repeat("a", invokers.PutLogEventsBatchMaxBytes)},
	}
	batch, consumed = takePutLogEventsBatch(over)
	if len(batch) != 0 || consumed != 1 {
		t.Fatalf("oversized: batch %d consumed %d, want 0/1", len(batch), consumed)
	}
}

func TestDeliverToCloudWatchLogsChunksAndSkipsOversized(t *testing.T) {
	store := newDeliveryTestStore(t, "deliver-chunk", "chunk-bucket", false)
	const groupARN = "arn:aws:logs:us-east-1:acc123:log-group:CloudTrail/logs"
	if _, err := store.MutateTrail("deliver-chunk", func(tr *cloudtrailstore.Trail) error {
		tr.CloudWatchLogsLogGroupARN = groupARN
		return nil
	}); err != nil {
		t.Fatalf("configure destination: %v", err)
	}
	trail, err := store.GetTrail("deliver-chunk")
	if err != nil {
		t.Fatalf("get trail: %v", err)
	}

	logs := newFakeLogsInvoker("CloudTrail/logs")
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{logs: logs})

	// The window walk serves events most-recent-first; one oversized
	// record sits among them (a recorder-side defect the delivery cannot
	// fix — it must not block the deliverable events).
	base := time.Now().UTC()
	oversized := &cloudtrailstore.Event{
		EventID:         "evt-oversized",
		EventTime:       base.Add(13 * time.Second),
		CloudTrailEvent: strings.Repeat("x", invokers.PutLogEventsBatchMaxBytes+100),
	}
	events := make([]*cloudtrailstore.Event, 0, 26)
	for i := 25; i >= 0; i-- {
		if i == 13 {
			events = append(events, oversized)
			continue
		}
		events = append(events, &cloudtrailstore.Event{
			EventID:         fmt.Sprintf("evt-chunk-%d", i),
			EventTime:       base.Add(time.Duration(i) * time.Second),
			CloudTrailEvent: fmt.Sprintf(`{"eventID":"evt-chunk-%d"}`, i),
		})
	}

	if err := svc.deliverToCloudWatchLogs(store, trail, "us-east-1", events); err != nil {
		t.Fatalf("deliverToCloudWatchLogs: %v", err)
	}

	stream := "CloudTrail/logs/acc123_CloudTrail_us-east-1"
	if len(logs.batchSizes[stream]) < 2 {
		t.Fatalf("expected the delivery split into batches, got call sizes %v", logs.batchSizes[stream])
	}
	for _, size := range logs.batchSizes[stream] {
		if size > invokers.PutLogEventsBatchMaxEvents {
			t.Fatalf("batch of %d exceeds the documented cardinality", size)
		}
	}
	entries := logs.events[stream]
	if len(entries) != 25 {
		t.Fatalf("stream carries %d entries, want the 25 deliverable events", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp < entries[i-1].Timestamp {
			t.Fatalf("entries not chronological at %d: %d after %d", i, entries[i].Timestamp, entries[i-1].Timestamp)
		}
	}
	for _, e := range entries {
		if e.Message == oversized.CloudTrailEvent {
			t.Fatal("oversized record was delivered despite exceeding the per-event bound")
		}
	}
}

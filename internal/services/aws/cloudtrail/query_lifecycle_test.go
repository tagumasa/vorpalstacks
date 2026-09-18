package cloudtrail

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/invokers"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below cover the Lake query lifecycle beyond the engine itself:
// admission under the concurrent-query bound, the request forms StartQuery
// accepts, the query deadline producing TIMED_OUT with partial rows, the
// ListQueries time filters, and DeliveryS3Uri delivery of result and sign
// files whose hashes and RSA signature verify.

// fakeS3Invoker is an in-memory S3Invoker: buckets map to keys to object
// content, which is all the delivery path reads and writes. Object
// metadata (the digest signature) is kept alongside for assertions. The
// permissive mode answers every bucket as existing with a sufficient
// CloudTrail policy — for Core tests whose subject is not the destination.
type fakeS3Invoker struct {
	buckets         map[string]map[string][]byte
	metadata        map[string]map[string]string
	policies        map[string]string
	permissive      bool
	panicOnPut      bool
	panicOnMetadata bool
}

func newFakeS3Invoker(buckets ...string) *fakeS3Invoker {
	f := &fakeS3Invoker{
		buckets:  make(map[string]map[string][]byte),
		metadata: make(map[string]map[string]string),
		policies: make(map[string]string),
	}
	for _, b := range buckets {
		f.buckets[b] = make(map[string][]byte)
		f.policies[b] = sufficientAnyBucketPolicy(b)
	}
	return f
}

// newPermissiveS3Invoker answers every bucket as existing with a
// sufficient policy.
func newPermissiveS3Invoker() *fakeS3Invoker {
	f := newFakeS3Invoker()
	f.permissive = true
	return f
}

// sufficientAnyBucketPolicy is a bucket policy granting the CloudTrail
// service principal the documented write access over the whole bucket —
// sufficient for every account and prefix combination.
func sufficientAnyBucketPolicy(bucket string) string {
	return fmt.Sprintf(`{"Version":"2012-10-17","Statement":[`+
		`{"Sid":"AWSCloudTrailAclCheck20150319","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::%s"},`+
		`{"Sid":"AWSCloudTrailWrite20150319","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::%s/*"}]}`,
		bucket, bucket)
}

func (f *fakeS3Invoker) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	if objects, ok := f.buckets[bucket]; ok {
		if data, ok := objects[key]; ok {
			return data, nil
		}
	}
	return nil, awserrors.NewAWSError("NoSuchKey", "The specified key does not exist.", 404)
}

func (f *fakeS3Invoker) GetObjectVersion(ctx context.Context, region, bucket, key, versionID string, maxBytes int64) ([]byte, error) {
	return f.GetObject(ctx, region, bucket, key, maxBytes)
}

func (f *fakeS3Invoker) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	if f.panicOnPut {
		panic("injected put failure")
	}
	objects, ok := f.buckets[bucket]
	if !ok {
		return awserrors.NewAWSError("NoSuchBucket", "The specified bucket does not exist.", 404)
	}
	objects[key] = data
	return nil
}

func (f *fakeS3Invoker) PutObjectWithMetadata(ctx context.Context, region, bucket, key string, data []byte, contentType string, metadata map[string]string) error {
	if f.panicOnMetadata {
		panic("injected metadata-put failure")
	}
	if err := f.PutObject(ctx, region, bucket, key, data, contentType); err != nil {
		return err
	}
	if _, ok := f.metadata[bucket]; !ok {
		f.metadata[bucket] = make(map[string]string)
	}
	for k, v := range metadata {
		f.metadata[bucket][key+"/"+k] = v
	}
	return nil
}

func (f *fakeS3Invoker) ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error) {
	keys := []string{}
	for k := range f.buckets[bucket] {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (f *fakeS3Invoker) ListObjectEntries(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]invokers.S3ObjectEntry, error) {
	return nil, nil
}

func (f *fakeS3Invoker) BucketExists(ctx context.Context, region, bucket string) (bool, error) {
	if f.permissive {
		return true, nil
	}
	_, ok := f.buckets[bucket]
	return ok, nil
}

func (f *fakeS3Invoker) GetBucketPolicy(ctx context.Context, region, bucket string) (string, error) {
	if f.permissive {
		return sufficientAnyBucketPolicy(bucket), nil
	}
	return f.policies[bucket], nil
}

func (f *fakeS3Invoker) EnsureBucket(ctx context.Context, region, bucket string) error {
	f.buckets[bucket] = make(map[string][]byte)
	return nil
}

func (f *fakeS3Invoker) DeleteObject(ctx context.Context, region, bucket, key string) error {
	delete(f.buckets[bucket], key)
	return nil
}

// seedLifecycleEDS creates an ENABLED event data store the query paths can
// target.
func seedLifecycleEDS(t *testing.T, store *cloudtrailstore.CloudTrailStore, id string) {
	t.Helper()
	_, err := store.CreateEventDataStore(&cloudtrailstore.EventDataStore{
		EventDataStoreID:  id,
		EventDataStoreARN: "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/" + id,
		Name:              "eds-" + id,
		Status:            "ENABLED",
		RetentionPeriod:   90,
	})
	require.NoError(t, err)
}

func requireAWSCode(t *testing.T, err error, code string, httpStatus int) {
	t.Helper()
	require.Error(t, err)
	var apiErr *awserrors.AWSError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, code, apiErr.GetCode())
	assert.Equal(t, httpStatus, apiErr.HTTPStatus)
}

// fakeInvokerRegistry serves the fake invokers through the registry
// surface the service resolves at call time; the embedded nil interface
// covers the registry's unused members.
type fakeInvokerRegistry struct {
	invokers.Registry
	s3   invokers.S3Invoker
	sns  invokers.SNSInvoker
	logs invokers.LogsInvoker
}

func (f fakeInvokerRegistry) S3Invoker() invokers.S3Invoker   { return f.s3 }
func (f fakeInvokerRegistry) SNSInvoker() invokers.SNSInvoker { return f.sns }
func (f fakeInvokerRegistry) LogsInvoker() invokers.LogsInvoker {
	return f.logs
}

func TestStartQueryRequestForms(t *testing.T) {
	store := newQueryTestStore(t)
	seedLifecycleEDS(t, store, "eds-forms")
	svc := NewCloudTrailService("acc123", "us-east-1")

	t.Run("alias form is rejected until dashboard templates exist", func(t *testing.T) {
		_, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
			QueryAlias:      "dashboard-query",
			QueryParameters: []string{"eds-forms"},
		})
		requireAWSCode(t, err, "UnsupportedOperationException", 400)
	})

	t.Run("statement and alias do not mix", func(t *testing.T) {
		_, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
			QueryStatement: "SELECT eventID FROM eds-forms",
			QueryAlias:     "dashboard-query",
		})
		requireAWSCode(t, err, "InvalidParameterException", 400)
	})

	t.Run("parameters without an alias are meaningless", func(t *testing.T) {
		_, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
			QueryParameters: []string{"eds-forms"},
		})
		requireAWSCode(t, err, "InvalidParameterException", 400)
	})

	t.Run("owner account mismatch is not found", func(t *testing.T) {
		_, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
			QueryStatement: "SELECT eventID FROM eds-forms",
			OwnerAccountID: "999999999999",
		})
		requireAWSCode(t, err, "EventDataStoreNotFoundException", 404)
	})

	t.Run("delivery target must be a valid S3 URI", func(t *testing.T) {
		svcWithInvoker := NewCloudTrailService("acc123", "us-east-1")
		svcWithInvoker.SetInvokerRegistry(fakeInvokerRegistry{s3: newFakeS3Invoker("results")})
		_, err := svcWithInvoker.startQueryCore(t.Context(), store, StartQueryInput{
			QueryStatement: "SELECT eventID FROM eds-forms",
			DeliveryS3URI:  "not-an-s3-uri",
		})
		requireAWSCode(t, err, "InvalidS3BucketNameException", 400)
	})

	t.Run("delivery target bucket must exist", func(t *testing.T) {
		svcWithInvoker := NewCloudTrailService("acc123", "us-east-1")
		svcWithInvoker.SetInvokerRegistry(fakeInvokerRegistry{s3: newFakeS3Invoker("results")})
		_, err := svcWithInvoker.startQueryCore(t.Context(), store, StartQueryInput{
			QueryStatement: "SELECT eventID FROM eds-forms",
			DeliveryS3URI:  "s3://missing-bucket",
		})
		requireAWSCode(t, err, "S3BucketDoesNotExistException", 404)
	})

	t.Run("response echoes the event data store owner", func(t *testing.T) {
		resp, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
			QueryStatement: "SELECT eventID FROM eds-forms",
		})
		require.NoError(t, err)
		assert.Equal(t, "acc123", resp["EventDataStoreOwnerAccountId"])
		assert.NotEmpty(t, resp["QueryId"])
	})
}

// Admission follows the documented concurrent-query bound: with the store
// already holding the maximum number of non-terminal queries, a further
// StartQuery is rejected and never recorded.
func TestStartQueryAdmissionBound(t *testing.T) {
	store := newQueryTestStore(t)
	seedLifecycleEDS(t, store, "eds-admit")
	svc := NewCloudTrailService("acc123", "us-east-1")

	seed := func(id string) {
		require.NoError(t, store.SaveQuery(&cloudtrailstore.QueryRecord{
			QueryID:        id,
			EventDataStore: "eds-admit",
			QueryStatus:    "RUNNING",
			StartTime:      time.Now().UTC(),
		}))
	}
	for i := 0; i < cloudtrailstore.MaxConcurrentQueries-1; i++ {
		seed(fmt.Sprintf("q-running-%d", i))
	}

	resp, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
		QueryStatement: "SELECT eventID FROM eds-admit",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp["QueryId"])

	// The admitted query's async executor finalises it on the empty store
	// almost immediately, which would free its admission slot before the
	// rejection below is attempted. Wait for the finalisation, then pin the
	// record back to RUNNING — the executor never touches a terminal record
	// again, so the bound is exercised deterministically.
	queryID := resp["QueryId"].(string)
	require.Eventually(t, func() bool {
		qr, getErr := store.GetQuery(queryID)
		return getErr == nil && cloudtrailstore.QueryTerminalStatus(qr.QueryStatus)
	}, 5*time.Second, 2*time.Millisecond)
	_, err = store.MutateQuery(queryID, func(qr *cloudtrailstore.QueryRecord) error {
		qr.QueryStatus = "RUNNING"
		return nil
	})
	require.NoError(t, err)

	_, err = svc.startQueryCore(t.Context(), store, StartQueryInput{
		QueryStatement: "SELECT eventID FROM eds-admit",
	})
	requireAWSCode(t, err, "MaxConcurrentQueriesException", 429)

	// A terminal query does not occupy an admission slot.
	_, err = store.MutateQuery(resp["QueryId"].(string), func(qr *cloudtrailstore.QueryRecord) error {
		qr.QueryStatus = "FINISHED"
		return nil
	})
	require.NoError(t, err)
	_, err = svc.startQueryCore(t.Context(), store, StartQueryInput{
		QueryStatement: "SELECT eventID FROM eds-admit",
	})
	require.NoError(t, err)
}

// A query whose deadline has passed transitions to TIMED_OUT, keeps the
// rows produced before the deadline, and records its end time.
func TestStartQueryDeadlineTimesOut(t *testing.T) {
	store := newQueryTestStore(t)
	seedLifecycleEDS(t, store, "eds-deadline")
	seedLakeEvent(t, store, "alice", "CreateTrail", false, nil)
	svc := NewCloudTrailService("acc123", "us-east-1")
	// A negative override places the deadline in the past, so the first
	// page check fires deterministically.
	svc.queryDeadline = -time.Hour

	resp, err := svc.startQueryCore(t.Context(), store, StartQueryInput{
		QueryStatement: "SELECT eventID FROM eds-deadline",
	})
	require.NoError(t, err)
	queryID := resp["QueryId"].(string)

	var qr *cloudtrailstore.QueryRecord
	for i := 0; i < 50; i++ {
		qr, err = store.GetQuery(queryID)
		require.NoError(t, err)
		if cloudtrailstore.QueryTerminalStatus(qr.QueryStatus) {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	assert.Equal(t, "TIMED_OUT", qr.QueryStatus)
	require.NotNil(t, qr.EndTime)

	// The timed-out finalisation keeps partial rows retrievable.
	require.NoError(t, store.SaveQuery(&cloudtrailstore.QueryRecord{
		QueryID:        "q-partial",
		EventDataStore: "eds",
		QueryStatus:    "RUNNING",
		StartTime:      time.Now().UTC(),
	}))
	partial := [][]map[string]string{{{"eventID": "e1"}}}
	final := svc.finaliseQueryExecution(store, "q-partial", partial, &lakeExecutionStats{eventsScanned: 3, eventsMatched: 1}, true, nil)
	require.NotNil(t, final)
	assert.Equal(t, "TIMED_OUT", final.QueryStatus)
	assert.Equal(t, partial, final.QueryResultRows)
	assert.Equal(t, int32(1), final.ResultsCount)
}

// ListQueries bounds its results to queries run within the requested
// period; the bounds accept both wire time forms.
func TestListQueriesTimeFilters(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")

	now := time.Now().UTC()
	seed := func(id string, start time.Time) {
		require.NoError(t, store.SaveQuery(&cloudtrailstore.QueryRecord{
			QueryID:        id,
			EventDataStore: "eds-filter",
			QueryStatus:    "FINISHED",
			StartTime:      start,
		}))
	}
	seed("q-old", now.Add(-2*time.Hour))
	seed("q-now", now)
	seed("q-future", now.Add(2*time.Hour))

	resp, err := svc.listQueriesCore(store, ListQueriesInput{
		EventDataStore: "eds-filter",
		StartTimeRaw:   float64(now.Add(-time.Hour).UnixMilli()) / 1000,
		EndTimeRaw:     float64(now.Add(time.Hour).UnixMilli()) / 1000,
	})
	require.NoError(t, err)
	queries := resp["Queries"].([]map[string]interface{})
	require.Len(t, queries, 1)
	assert.Equal(t, "q-now", queries[0]["QueryId"])

	resp, err = svc.listQueriesCore(store, ListQueriesInput{
		EventDataStore: "eds-filter",
		StartTimeStr:   now.Add(time.Hour).Format(time.RFC3339),
	})
	require.NoError(t, err)
	queries = resp["Queries"].([]map[string]interface{})
	require.Len(t, queries, 1)
	assert.Equal(t, "q-future", queries[0]["QueryId"])

	_, err = svc.listQueriesCore(store, ListQueriesInput{
		EventDataStore: "eds-filter",
		StartTimeStr:   now.Format(time.RFC3339),
		EndTimeStr:     now.Add(-time.Hour).Format(time.RFC3339),
	})
	// ListQueries declares InvalidDateRangeException — not the LookupEvents
	// time-range shape — for an out-of-order period.
	requireAWSCode(t, err, "InvalidDateRangeException", 400)
}

// A FINISHED query with a delivery target writes its gzip CSV result file
// and sign file under the documented S3 path; the sign file's hashes match
// the compressed content and its RSA signature verifies against the stored
// public key with the matching fingerprint.
func TestDeliverQueryResults(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")
	invoker := newFakeS3Invoker("results-bucket")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: invoker})

	completed := time.Now().UTC()
	rows := [][]map[string]string{
		{{"eventID": "e1"}, {"eventName": "CreateTrail"}},
		{{"eventID": "e2"}, {"eventName": "DeleteTrail"}},
	}
	qr := &cloudtrailstore.QueryRecord{
		QueryID:         "q-deliver",
		EventDataStore:  "eds-deliver",
		QueryStatement:  "SELECT eventID, eventName FROM eds-deliver",
		QueryStatus:     "FINISHED",
		QueryResultRows: rows,
		ResultsCount:    2,
		StartTime:       completed.Add(-time.Minute),
		EndTime:         &completed,
		DeliveryS3URI:   "s3://results-bucket/ct-prefix",
		DeliveryStatus:  "PENDING",
	}
	require.NoError(t, store.SaveQuery(qr))

	svc.deliverQueryResults(store, qr)

	after, err := store.GetQuery("q-deliver")
	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", after.DeliveryStatus)

	base := fmt.Sprintf("ct-prefix/AWSLogs/acc123/CloudTrail-Lake/Query/%04d/%02d/%02d/q-deliver",
		completed.Year(), completed.Month(), completed.Day())
	gz := invoker.buckets["results-bucket"][base+"/result_1.csv.gz"]
	require.NotEmpty(t, gz, "result file must be delivered at %s", base)

	zr, err := gzip.NewReader(bytes.NewReader(gz))
	require.NoError(t, err)
	csvText, err := io.ReadAll(zr)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(csvText), "eventID,eventName"), "CSV header row: %s", csvText)
	assert.Contains(t, string(csvText), "CreateTrail")
	assert.Contains(t, string(csvText), "DeleteTrail")

	signRaw := invoker.buckets["results-bucket"][base+"/result_sign.json"]
	require.NotEmpty(t, signRaw, "sign file must be delivered")
	var sign queryResultSignFile
	require.NoError(t, json.Unmarshal(signRaw, &sign))
	assert.Equal(t, "1.0", sign.Version)
	assert.Equal(t, "us-east-1", sign.Region)
	assert.Equal(t, "SHA-256", sign.HashAlgorithm)
	assert.Equal(t, "SHA256withRSA", sign.SignatureAlgorithm)
	require.Len(t, sign.Files, 1)
	assert.Equal(t, "result_1.csv.gz", sign.Files[0].FileName)

	sum := sha256.Sum256(gz)
	assert.Equal(t, hex.EncodeToString(sum[:]), sign.Files[0].FileHashValue)

	// The signature verifies against the stored public key the sign file's
	// fingerprint identifies — the validator's path. The key is served in
	// the documented PKCS#1 DER form.
	keys, err := store.ListPublicKeys(nil, nil)
	require.NoError(t, err)
	var verified bool
	for _, pk := range keys {
		if pk.Fingerprint() != sign.PublicKeyFingerprint {
			continue
		}
		rsaPub, err := x509.ParsePKCS1PublicKey(pk.Value)
		require.NoError(t, err)
		digest := sha256.Sum256([]byte(sign.Files[0].FileHashValue))
		sig, err := hex.DecodeString(sign.HashSignature)
		require.NoError(t, err)
		verified = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest[:], sig) == nil
	}
	assert.True(t, verified, "hashSignature must verify against the fingerprint-matched public key")
}

// A delivery interrupted mid-body — here by an injected panic in the
// object write — must record FAILED, not the success it never reached:
// the deferred status record starts from FAILED and only the fully
// delivered path upgrades it.
func TestDeliverQueryResultsPanicRecordsFailed(t *testing.T) {
	store := newQueryTestStore(t)
	svc := NewCloudTrailService("acc123", "us-east-1")
	invoker := newFakeS3Invoker("results-bucket")
	invoker.panicOnPut = true
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: invoker})

	qr := &cloudtrailstore.QueryRecord{
		QueryID:         "q-panic",
		EventDataStore:  "eds-deliver",
		QueryStatement:  "SELECT eventID FROM eds-deliver",
		QueryStatus:     "FINISHED",
		QueryResultRows: [][]map[string]string{{{"eventID": "e1"}}},
		ResultsCount:    1,
		StartTime:       time.Now().UTC().Add(-time.Minute),
		EndTime:         func() *time.Time { t := time.Now().UTC(); return &t }(),
		DeliveryS3URI:   "s3://results-bucket",
		DeliveryStatus:  "PENDING",
	}
	require.NoError(t, store.SaveQuery(qr))

	require.Panics(t, func() { svc.deliverQueryResults(store, qr) })

	after, err := store.GetQuery("q-panic")
	require.NoError(t, err)
	assert.Equal(t, "FAILED", after.DeliveryStatus)
	assert.Empty(t, invoker.buckets["results-bucket"], "no file may be recorded as delivered")
}

// The delivery splits into numbered files at the per-file size bound, each
// carrying the header row, the sign file names every delivered file, and
// every completed file reaches the sink — hash verified — as it is
// finalised, so the delivery path holds one file's content at a time.
func TestBuildQueryResultFilesSplits(t *testing.T) {
	pq, err := parseQueryStatement("SELECT eventID FROM eds")
	require.NoError(t, err)
	rows := make([][]map[string]string, 20)
	for i := range rows {
		rows[i] = []map[string]string{{"eventID": fmt.Sprintf("event-id-%02d", i)}}
	}

	var delivered [][]byte
	entries, err := buildQueryResultFiles(pq, rows, 60, func(name string, content []byte, hashHex string) error {
		sum := sha256.Sum256(content)
		if hex.EncodeToString(sum[:]) != hashHex {
			t.Errorf("hash mismatch on %s", name)
		}
		delivered = append(delivered, content)
		return nil
	})
	require.NoError(t, err)
	require.True(t, len(entries) >= 2, "expected a split, got %d files", len(entries))
	require.Len(t, delivered, len(entries))

	totalRows := 0
	for i, content := range delivered {
		zr, err := gzip.NewReader(bytes.NewReader(content))
		require.NoError(t, err)
		text, err := io.ReadAll(zr)
		require.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(string(text)), "\n")
		assert.Equal(t, "eventID", lines[0], "every chunk carries the header")
		totalRows += len(lines) - 1
		assert.Equal(t, fmt.Sprintf("result_%d.csv.gz", i+1), entries[i].FileName)
	}
	assert.Equal(t, 20, totalRows)

	pub, priv, err := cloudtrailstore.GenerateKeyPair()
	require.NoError(t, err)
	sign, err := buildQueryResultSignFile("us-east-1", entries, pub, priv)
	require.NoError(t, err)
	require.Len(t, sign.Files, len(entries))
	for i, e := range entries {
		assert.Equal(t, e.FileName, sign.Files[i].FileName)
	}
}

// TestListQueriesSevenDayWindow pins the listing's own retention: "Returns
// a list of queries and query statuses for the past seven days"
// (ListQueries) — a record older than the window is excluded from the
// listing even without caller bounds, and the store's retention sweep
// drops it while keeping the recent one.
func TestListQueriesSevenDayWindow(t *testing.T) {
	store := newQueryTestStore(t)
	svc := &CloudTrailService{}
	seedLakeEDS(t, store)
	edsARN := "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/eds"

	oldStart := time.Now().UTC().Add(-cloudtrailstore.ListQueriesWindow - time.Hour)
	require.NoError(t, store.SaveQuery(&cloudtrailstore.QueryRecord{
		QueryID: "q-old", EventDataStore: edsARN, QueryStatus: "FINISHED", StartTime: oldStart,
	}))
	require.NoError(t, store.SaveQuery(&cloudtrailstore.QueryRecord{
		QueryID: "q-recent", EventDataStore: edsARN, QueryStatus: "FINISHED",
		StartTime: time.Now().UTC().Add(-time.Hour),
	}))

	resp, err := svc.listQueriesCore(store, ListQueriesInput{EventDataStore: edsARN})
	require.NoError(t, err)
	ids := map[string]bool{}
	for _, item := range resp["Queries"].([]map[string]interface{}) {
		ids[item["QueryId"].(string)] = true
	}
	assert.True(t, ids["q-recent"], "recent query is listed")
	assert.False(t, ids["q-old"], "a query older than the seven-day window is not listed")

	purged, err := store.PurgeQueriesBefore(time.Now().UTC().Add(-cloudtrailstore.ListQueriesWindow))
	require.NoError(t, err)
	assert.Equal(t, 1, purged, "the sweep drops exactly the aged record")
	_, err = store.GetQuery("q-old")
	assert.ErrorIs(t, err, cloudtrailstore.ErrQueryNotFound)
	_, err = store.GetQuery("q-recent")
	require.NoError(t, err)
}

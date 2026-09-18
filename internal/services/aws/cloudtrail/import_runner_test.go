package cloudtrail

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vorpalstacks/internal/common/invokers"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below cover the import runner: trail records in gzip log files
// land in the destination event data store with their own identifiers and
// category, the file-name date bounds and the destination retention bound
// the copy, unreadable files become ListImportFailures entries and flip
// the import to FAILED, and a stop that lands mid-walk holds against the
// runner's remaining transitions.

// gzipRecords renders one delivered log file: the record JSON objects
// wrapped in the {"Records":[...]} envelope and gzip-compressed. The
// envelope is built by concatenation so the records ride as JSON objects,
// not as escaped strings.
func gzipRecords(t *testing.T, records ...string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write([]byte(`{"Records":[` + strings.Join(records, ",") + `]}`))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

// importRecord renders one trail record with the record-format spellings
// the importer parses.
func importRecord(id, eventTime, eventName string) string {
	return fmt.Sprintf(`{"eventVersion":"1.09","eventID":%q,"eventTime":%q,`+
		`"eventSource":"cloudtrail.amazonaws.com","eventName":%q,`+
		`"awsRegion":"us-east-1","readOnly":false,"eventType":"AwsApiCall","managementEvent":true}`,
		id, eventTime, eventName)
}

// newImportRunnerStore provisions a store with a 7-day destination event
// data store, so the retention bound pins with records days apart.
func newImportRunnerStore(t *testing.T) (*cloudtrailstore.CloudTrailStore, *cloudtrailstore.EventDataStore) {
	t.Helper()
	store := newLockTestStore(t)
	eds, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("import-dst", "acc123", "us-east-1"))
	require.NoError(t, err)
	_, err = store.MutateEventDataStore(eds.EventDataStoreID, func(eds *cloudtrailstore.EventDataStore) error {
		eds.RetentionPeriod = 7
		return nil
	})
	require.NoError(t, err)
	return store, eds
}

// importSourceRaw builds the ImportSource wire value for one S3 location.
func importSourceRaw(uri string) map[string]interface{} {
	return map[string]interface{}{"S3": map[string]interface{}{
		"S3LocationUri":         uri,
		"S3BucketRegion":        "us-east-1",
		"S3BucketAccessRoleArn": "arn:aws:iam::acc123:role/import-role",
	}}
}

// startTestImport runs the new-import core path synchronously and waits
// for the runner goroutine to settle the status.
func startTestImport(t *testing.T, svc *CloudTrailService, store *cloudtrailstore.CloudTrailStore, edsARN, uri string, start, end *time.Time) *cloudtrailstore.Import {
	t.Helper()
	in := StartImportInput{
		Destinations:         []string{edsARN},
		ImportSourceRaw:      importSourceRaw(uri),
		ImportSourceProvided: true,
	}
	if start != nil {
		in.StartEventTimeRaw = float64(start.Unix())
	}
	if end != nil {
		in.EndEventTimeRaw = float64(end.Unix())
	}
	resp, err := svc.startImportCore(t.Context(), store, in)
	require.NoError(t, err)
	return waitImportSettled(t, store, resp["ImportId"].(string))
}

// waitImportSettled polls the import until it reaches a terminal status.
func waitImportSettled(t *testing.T, store *cloudtrailstore.CloudTrailStore, importID string) *cloudtrailstore.Import {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		imp, err := store.GetImport(importID)
		require.NoError(t, err)
		if importTerminalStatus(imp.ImportStatus) {
			return imp
		}
		if time.Now().After(deadline) {
			t.Fatalf("import %s did not settle; status %s", importID, imp.ImportStatus)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestImportRunnerImportsRecords(t *testing.T) {
	store, eds := newImportRunnerStore(t)
	now := time.Now().UTC()
	recent := now.Add(-time.Hour).Format(time.RFC3339)
	older := now.Add(-time.Hour * 2).Format(time.RFC3339)
	outsideRetention := now.AddDate(0, 0, -10).Format(time.RFC3339)

	s3 := newFakeS3Invoker("src-bucket")
	s3.buckets["src-bucket"]["logs/a.json.gz"] = gzipRecords(t,
		importRecord("imp-a1", recent, "CreateTrail"),
		importRecord("imp-a2", older, "DeleteTrail"))
	s3.buckets["src-bucket"]["logs/b.json.gz"] = gzipRecords(t,
		importRecord("imp-b1", recent, "ListTrails"),
		importRecord("imp-b2", outsideRetention, "GetTrail"))
	s3.buckets["src-bucket"]["logs/plain.json"] = []byte(`{"Records":[]}`) // uncompressed: never a candidate
	s3.buckets["src-bucket"]["elsewhere/c.json.gz"] = gzipRecords(t, importRecord("imp-c1", recent, "StartLogging"))

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3})

	imp := startTestImport(t, svc, store, eds.EventDataStoreARN, "s3://src-bucket/logs/", nil, nil)

	assert.Equal(t, "COMPLETED", imp.ImportStatus)
	assert.Equal(t, cloudtrailstore.ImportStatistics{
		PrefixesFound:     1,
		PrefixesCompleted: 1,
		FilesCompleted:    2,
		EventsCompleted:   3,
	}, imp.ImportStatistics, "two gzip files import their three in-retention records; the uncompressed file and the other prefix are out of the walk")

	imported, _, err := store.LookupEDSEvents(eds.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	require.Len(t, imported, 3)
	assert.Equal(t, "imp-a2", imported[0].EventID) // ascending event time: the oldest first
	assert.Equal(t, "imp-a1", imported[1].EventID)
	assert.Equal(t, "imp-b1", imported[2].EventID)
	for _, e := range imported {
		assert.Equal(t, "Management", e.EventCategory, "records without an eventCategory import as management events")
		assert.Contains(t, e.CloudTrailEvent, e.EventID, "the record JSON rides verbatim")
	}
}

func TestImportRunnerFileFailures(t *testing.T) {
	store, eds := newImportRunnerStore(t)
	recent := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)

	s3 := newFakeS3Invoker("src-bucket")
	s3.buckets["src-bucket"]["logs/good.json.gz"] = gzipRecords(t, importRecord("fl-1", recent, "CreateTrail"))
	s3.buckets["src-bucket"]["logs/not-gzip.json.gz"] = []byte("plain text, not a gzip stream")
	var rawJson bytes.Buffer
	zw := gzip.NewWriter(&rawJson)
	_, _ = zw.Write([]byte(`{"notRecords":[]}`))
	require.NoError(t, zw.Close())
	s3.buckets["src-bucket"]["logs/not-records.json.gz"] = rawJson.Bytes()
	s3.buckets["src-bucket"]["logs/bad-record.json.gz"] = gzipRecords(t,
		importRecord("fl-2", recent, "ListTrails"), `"oops"`)

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3})

	imp := startTestImport(t, svc, store, eds.EventDataStoreARN, "s3://src-bucket/logs/", nil, nil)

	assert.Equal(t, "FAILED", imp.ImportStatus, "an import with any failure finishes FAILED")
	assert.Equal(t, int64(2), imp.ImportStatistics.FilesCompleted, "the decoded files completed; the unreadable ones did not")
	assert.Equal(t, int64(2), imp.ImportStatistics.EventsCompleted)
	assert.Equal(t, int64(3), imp.ImportStatistics.FailedEntries)
	require.Len(t, imp.Failures, 3)
	byLocation := map[string]cloudtrailstore.ImportFailure{}
	for _, f := range imp.Failures {
		byLocation[f.Location] = f
	}
	require.Contains(t, byLocation, "s3://src-bucket/logs/not-gzip.json.gz")
	assert.Equal(t, "InvalidFileFormat", byLocation["s3://src-bucket/logs/not-gzip.json.gz"].ErrorType)
	require.Contains(t, byLocation, "s3://src-bucket/logs/not-records.json.gz")
	assert.Equal(t, "InvalidFileFormat", byLocation["s3://src-bucket/logs/not-records.json.gz"].ErrorType)
	require.Contains(t, byLocation, "s3://src-bucket/logs/bad-record.json.gz")
	assert.Equal(t, "InvalidEventRecord", byLocation["s3://src-bucket/logs/bad-record.json.gz"].ErrorType)
	assert.Equal(t, "FAILED", byLocation["s3://src-bucket/logs/not-gzip.json.gz"].Status)
}

func TestImportRunnerTimeBoundsAndDefaultPrefix(t *testing.T) {
	store, eds := newImportRunnerStore(t)
	recent := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)

	s3 := newFakeS3Invoker("src-bucket")
	// The default walk (no prefix in the URI) admits only keys under a
	// CloudTrail path segment; the Config service folder stays out.
	log := func(day int, stamp, unique string) string {
		return fmt.Sprintf("AWSLogs/acc123/CloudTrail/us-east-1/2026/09/%02d/acc123_CloudTrail_us-east-1_%s_%s.json.gz", day, stamp, unique)
	}
	s3.buckets["src-bucket"][log(10, "20260910T1200Z", "aaaaaaaaaaaaaaaa")] = gzipRecords(t, importRecord("tb-old", recent, "CreateTrail"))
	s3.buckets["src-bucket"][log(12, "20260912T1200Z", "bbbbbbbbbbbbbbbb")] = gzipRecords(t, importRecord("tb-in", recent, "CreateTrail"))
	s3.buckets["src-bucket"]["AWSLogs/acc123/Config/us-east-1/2026/09/11/config_20260911T1200Z.json.gz"] = gzipRecords(t, importRecord("tb-config", recent, "PutConfig"))
	s3.buckets["src-bucket"]["AWSLogs/acc123/CloudTrail/undated.json.gz"] = gzipRecords(t, importRecord("tb-undated", recent, "CreateTrail"))

	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3})

	start := time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC)
	imp := startTestImport(t, svc, store, eds.EventDataStoreARN, "s3://src-bucket", &start, &end)

	assert.Equal(t, "COMPLETED", imp.ImportStatus)
	assert.Equal(t, int64(1), imp.ImportStatistics.FilesCompleted,
		"the out-of-range, undated, and non-CloudTrail files are never attempted")

	imported, _, err := store.LookupEDSEvents(eds.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	require.Len(t, imported, 1)
	assert.Equal(t, "tb-in", imported[0].EventID)
}

func TestStartImportGuards(t *testing.T) {
	store, eds := newImportRunnerStore(t)
	recent := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)

	s3 := newFakeS3Invoker("src-bucket")
	s3.buckets["src-bucket"]["logs/a.json.gz"] = gzipRecords(t, importRecord("g-1", recent, "CreateTrail"))
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3})

	t.Run("malformed source URI", func(t *testing.T) {
		_, err := svc.startImportCore(t.Context(), store, StartImportInput{
			Destinations:         []string{eds.EventDataStoreARN},
			ImportSourceRaw:      importSourceRaw("http://src-bucket/logs/"),
			ImportSourceProvided: true,
		})
		requireAWSCode(t, err, "InvalidImportSourceException", 400)
	})

	t.Run("missing source bucket", func(t *testing.T) {
		_, err := svc.startImportCore(t.Context(), store, StartImportInput{
			Destinations:         []string{eds.EventDataStoreARN},
			ImportSourceRaw:      importSourceRaw("s3://ghost-bucket/logs/"),
			ImportSourceProvided: true,
		})
		requireAWSCode(t, err, "InvalidImportSourceException", 400)
	})

	t.Run("unavailable S3 service", func(t *testing.T) {
		bare := NewCloudTrailService("acc123", "us-east-1")
		_, err := bare.startImportCore(t.Context(), store, StartImportInput{
			Destinations:         []string{eds.EventDataStoreARN},
			ImportSourceRaw:      importSourceRaw("s3://src-bucket/logs/"),
			ImportSourceProvided: true,
		})
		requireAWSCode(t, err, "InvalidImportSourceException", 400)
	})

	t.Run("missing access role", func(t *testing.T) {
		_, err := svc.startImportCore(t.Context(), store, StartImportInput{
			Destinations: []string{eds.EventDataStoreARN},
			ImportSourceRaw: map[string]interface{}{"S3": map[string]interface{}{
				"S3LocationUri": "s3://src-bucket/logs/", "S3BucketRegion": "us-east-1"}},
			ImportSourceProvided: true,
		})
		requireAWSCode(t, err, "InvalidImportSourceException", 400)
	})

	t.Run("ongoing import holds a new one", func(t *testing.T) {
		ongoing, err := store.CreateImport(cloudtrailstore.NewImport(
			[]string{eds.EventDataStoreARN}, cloudtrailstore.ImportSource{
				S3LocationURI: "s3://src-bucket/logs/", S3BucketRegion: "us-east-1"}))
		require.NoError(t, err)
		_, err = store.MutateImport(ongoing.ImportID, func(imp *cloudtrailstore.Import) error {
			imp.ImportStatus = "IN_PROGRESS"
			return nil
		})
		require.NoError(t, err)

		_, err = svc.startImportCore(t.Context(), store, StartImportInput{
			Destinations:         []string{eds.EventDataStoreARN},
			ImportSourceRaw:      importSourceRaw("s3://src-bucket/logs/"),
			ImportSourceProvided: true,
		})
		requireAWSCode(t, err, "AccountHasOngoingImportException", 400)

		// A terminal import frees the account again.
		_, err = svc.stopImportCore(store, ImportIDInput{ImportID: ongoing.ImportID})
		require.NoError(t, err)
		imp := startTestImport(t, svc, store, eds.EventDataStoreARN, "s3://src-bucket/logs/", nil, nil)
		assert.Equal(t, "COMPLETED", imp.ImportStatus)
	})
}

// gatedS3Invoker wraps the fake invoker with a hook before each GetObject,
// so a stop landing mid-file is deterministic.
type gatedS3Invoker struct {
	*fakeS3Invoker
	beforeGetObject func(key string)
}

func (g *gatedS3Invoker) GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error) {
	if g.beforeGetObject != nil {
		g.beforeGetObject(key)
	}
	return g.fakeS3Invoker.GetObject(ctx, region, bucket, key, maxBytes)
}

var _ invokers.S3Invoker = (*gatedS3Invoker)(nil)

// TestImportRunnerStopMidWalkHolds pins the between-files stop: a
// StopImport that lands while a file is being read wins against the
// runner's progress write, the walk abandons, and STOPPED stands with no
// completion statistics. The events already copied into the destination
// stay copied — the stop bounds the walk, not the past.
func TestImportRunnerStopMidWalkHolds(t *testing.T) {
	store, eds := newImportRunnerStore(t)
	recent := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)

	s3 := newFakeS3Invoker("src-bucket")
	s3.buckets["src-bucket"]["logs/a.json.gz"] = gzipRecords(t, importRecord("stop-1", recent, "CreateTrail"))
	s3.buckets["src-bucket"]["logs/b.json.gz"] = gzipRecords(t, importRecord("stop-2", recent, "DeleteTrail"))

	started := make(chan string)
	stopped := make(chan struct{})
	gated := &gatedS3Invoker{fakeS3Invoker: s3}
	gated.beforeGetObject = func(key string) {
		if strings.HasSuffix(key, "a.json.gz") {
			close(started)
			<-stopped
		}
	}
	svc := NewCloudTrailService("acc123", "us-east-1")
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: gated})

	created, err := store.CreateImport(cloudtrailstore.NewImport(
		[]string{eds.EventDataStoreARN}, cloudtrailstore.ImportSource{
			S3LocationURI: "s3://src-bucket/logs/", S3BucketRegion: "us-east-1"}))
	require.NoError(t, err)
	go svc.runImport(store, created.ImportID)

	<-started
	_, err = svc.stopImportCore(store, ImportIDInput{ImportID: created.ImportID})
	require.NoError(t, err)
	close(stopped)

	imp := waitImportSettled(t, store, created.ImportID)
	assert.Equal(t, "STOPPED", imp.ImportStatus)
	assert.Equal(t, cloudtrailstore.ImportStatistics{}, imp.ImportStatistics,
		"the progress write refuses over STOPPED, so no completion statistics land")

	// The first file's records were already being copied when the stop
	// landed; the second file is never read. The runner finishes its
	// in-flight file after the gate opens, so the copy is awaited, not
	// assumed.
	require.Eventually(t, func() bool {
		imported, _, err := store.LookupEDSEvents(eds.EventDataStoreID, cloudtrailstore.EDSQuery{})
		return err == nil && len(imported) == 1
	}, 5*time.Second, 5*time.Millisecond, "the first file's records were copied before the walk abandoned")
	imported, _, err := store.LookupEDSEvents(eds.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	require.Len(t, imported, 1)
	assert.Equal(t, "stop-1", imported[0].EventID)
}

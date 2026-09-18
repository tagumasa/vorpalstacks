package testutil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *TestRunner) runCloudTrailImportTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// Create a dedicated EDS for import tests.
	var edsARN string
	var edsID string
	results = append(results, r.RunTest("cloudtrail", "Import_Setup_EDS", func() error {
		resp, err := tc.createEventDataStore("import-test", nil)
		if err != nil {
			return fmt.Errorf("create EDS for import: %w", err)
		}
		if resp.EventDataStoreArn == nil {
			return fmt.Errorf("EDS ARN is nil")
		}
		edsARN = *resp.EventDataStoreArn
		edsID = tc.edsIDFromARN(edsARN)
		return nil
	}))
	defer func() {
		if edsARN != "" {
			_ = tc.deleteEventDataStore(edsID)
		}
	}()

	// Provision the import source: a bucket whose CloudTrail prefix holds
	// two delivered-form gzip log files (three records) plus an
	// uncompressed file the import must ignore, and a broken prefix whose
	// one corrupt file drives the failure path.
	bucket := tc.uniqueName("ct-import")
	results = append(results, r.RunTest("cloudtrail", "Import_Setup_Bucket", func() error {
		if err := tc.createDeliveryBucket(bucket); err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		now := time.Now().UTC()
		recent := now.Add(-time.Hour).Format(time.RFC3339)
		record := func(id, name string) string {
			return fmt.Sprintf(`{"eventVersion":"1.09","eventID":%q,"eventTime":%q,`+
				`"eventSource":"cloudtrail.amazonaws.com","eventName":%q,`+
				`"awsRegion":"us-east-1","readOnly":false,"eventType":"AwsApiCall","managementEvent":true}`,
				id, recent, name)
		}
		day := now.Format("2006/01/02")
		logKey := func(stamp, unique string) string {
			return fmt.Sprintf("CloudTrail/AWSLogs/%s/CloudTrail/%s/%s/%s_CloudTrail_%s_%s_%s.json.gz",
				tc.accountID, tc.region, day, tc.accountID, tc.region, stamp, unique)
		}
		putGzip := func(key string, records ...string) error {
			var buf bytes.Buffer
			zw := gzip.NewWriter(&buf)
			if _, err := zw.Write([]byte(`{"Records":[` + strings.Join(records, ",") + `]}`)); err != nil {
				return err
			}
			if err := zw.Close(); err != nil {
				return err
			}
			_, err := tc.s3Client.PutObject(tc.ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucket), Key: aws.String(key), Body: bytes.NewReader(buf.Bytes()),
			})
			return err
		}
		stamp := now.Format("20060102T1504Z")
		if err := putGzip(logKey(stamp, "aaaaaaaaaaaaaaaa"),
			record("sdk-imp-1", "CreateTrail"), record("sdk-imp-2", "DeleteTrail")); err != nil {
			return fmt.Errorf("put first log file: %w", err)
		}
		if err := putGzip(logKey(stamp, "bbbbbbbbbbbbbbbb"), record("sdk-imp-3", "ListTrails")); err != nil {
			return fmt.Errorf("put second log file: %w", err)
		}
		_, err := tc.s3Client.PutObject(tc.ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket), Key: aws.String("CloudTrail/readme.txt"),
			Body: strings.NewReader("not a log file"),
		})
		if err != nil {
			return fmt.Errorf("put uncompressed object: %w", err)
		}
		_, err = tc.s3Client.PutObject(tc.ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket), Key: aws.String("broken/corrupt.json.gz"),
			Body: strings.NewReader("plain bytes under a gzip name"),
		})
		if err != nil {
			return fmt.Errorf("put corrupt object: %w", err)
		}
		return nil
	}))
	defer tc.deleteDeliveryBucket(bucket)

	waitTerminal := func(importID string) (*cloudtrail.GetImportOutput, error) {
		deadline := time.Now().Add(15 * time.Second)
		for {
			resp, err := tc.client.GetImport(tc.ctx, &cloudtrail.GetImportInput{ImportId: aws.String(importID)})
			if err != nil {
				return nil, fmt.Errorf("GetImport failed: %w", err)
			}
			switch resp.ImportStatus {
			case types.ImportStatusCompleted, types.ImportStatusFailed, types.ImportStatusStopped:
				return resp, nil
			}
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("import %s did not settle; status %s", importID, resp.ImportStatus)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	source := func(uri string) *types.ImportSource {
		return &types.ImportSource{
			S3: &types.S3ImportSource{
				S3LocationUri:         aws.String(uri),
				S3BucketRegion:        aws.String(tc.region),
				S3BucketAccessRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/CloudTrailImport", tc.accountID)),
			},
		}
	}

	// StartImport runs the real copy: two gzip files under the CloudTrail
	// prefix land their three records in the destination, the uncompressed
	// object is never a candidate, and the measured statistics report it.
	var importID string
	results = append(results, r.RunTest("cloudtrail", "StartImport_Success", func() error {
		resp, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			Destinations: []string{edsARN},
			ImportSource: source(fmt.Sprintf("s3://%s/CloudTrail/", bucket)),
		})
		if err != nil {
			return fmt.Errorf("StartImport failed: %w", err)
		}
		if resp.ImportId == nil {
			return fmt.Errorf("ImportId is nil")
		}
		importID = *resp.ImportId
		done, err := waitTerminal(importID)
		if err != nil {
			return err
		}
		if done.ImportStatus != types.ImportStatusCompleted {
			return fmt.Errorf("ImportStatus = %s, want COMPLETED", done.ImportStatus)
		}
		st := done.ImportStatistics
		if st == nil {
			return fmt.Errorf("ImportStatistics is nil")
		}
		if aws.ToInt64(st.PrefixesFound) != 1 || aws.ToInt64(st.PrefixesCompleted) != 1 {
			return fmt.Errorf("prefix statistics = %d/%d, want 1/1", aws.ToInt64(st.PrefixesFound), aws.ToInt64(st.PrefixesCompleted))
		}
		if aws.ToInt64(st.FilesCompleted) != 2 {
			return fmt.Errorf("FilesCompleted = %d, want 2 (the uncompressed object is not a candidate)", aws.ToInt64(st.FilesCompleted))
		}
		if aws.ToInt64(st.EventsCompleted) != 3 {
			return fmt.Errorf("EventsCompleted = %d, want 3", aws.ToInt64(st.EventsCompleted))
		}
		if aws.ToInt64(st.FailedEntries) != 0 {
			return fmt.Errorf("FailedEntries = %d, want 0", aws.ToInt64(st.FailedEntries))
		}
		return nil
	}))

	// The imported trail events are queryable in the destination EDS.
	results = append(results, r.RunTest("cloudtrail", "StartImport_QueriesImportedEvents", func() error {
		rows, err := tc.runLakeQueryAndWait(edsID, fmt.Sprintf(
			"SELECT eventID, eventName FROM %s WHERE eventID = 'sdk-imp-1'", edsID))
		if err != nil {
			return err
		}
		if len(rows) != 1 || rows[0]["eventName"] != "CreateTrail" {
			return fmt.Errorf("query rows = %v, want the imported sdk-imp-1 CreateTrail event", rows)
		}
		return nil
	}))

	// The time bounds arrive as Unix epochs on the JSON wire and must
	// bound the import — never silently drop to an unbounded one. Bounds
	// no delivered file name falls inside skip every candidate file.
	results = append(results, r.RunTest("cloudtrail", "StartImport_TimeBounds", func() error {
		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		resp, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			Destinations:   []string{edsARN},
			ImportSource:   source(fmt.Sprintf("s3://%s/CloudTrail/", bucket)),
			StartEventTime: aws.Time(start),
			EndEventTime:   aws.Time(end),
		})
		if err != nil {
			return fmt.Errorf("StartImport with time bounds failed: %w", err)
		}
		if resp.StartEventTime == nil || !resp.StartEventTime.Equal(start) {
			return fmt.Errorf("StartEventTime echo = %v, want %v", resp.StartEventTime, start)
		}
		if resp.EndEventTime == nil || !resp.EndEventTime.Equal(end) {
			return fmt.Errorf("EndEventTime echo = %v, want %v", resp.EndEventTime, end)
		}

		// The created import carries the bounds on its record too, and
		// completes with no file inside the range.
		done, err := waitTerminal(*resp.ImportId)
		if err != nil {
			return err
		}
		if done.StartEventTime == nil || !done.StartEventTime.Equal(start) {
			return fmt.Errorf("stored StartEventTime = %v, want %v", done.StartEventTime, start)
		}
		if done.ImportStatus != types.ImportStatusCompleted {
			return fmt.Errorf("ImportStatus = %s, want COMPLETED", done.ImportStatus)
		}
		if st := done.ImportStatistics; st != nil && aws.ToInt64(st.FilesCompleted) != 0 {
			return fmt.Errorf("FilesCompleted = %d, want 0 for a range no file name falls inside", aws.ToInt64(st.FilesCompleted))
		}
		return nil
	}))

	// A source bucket that does not exist is rejected at start.
	results = append(results, r.RunTest("cloudtrail", "StartImport_MissingSourceBucket", func() error {
		_, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			Destinations: []string{edsARN},
			ImportSource: source("s3://no-such-import-bucket-xyz/CloudTrail/"),
		})
		return AssertErrorContains(err, "InvalidImportSourceException")
	}))

	// A corrupt log file fails its entry, the failure is listable, and the
	// import finishes FAILED.
	var failedImportID string
	results = append(results, r.RunTest("cloudtrail", "StartImport_Failures", func() error {
		resp, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			Destinations: []string{edsARN},
			ImportSource: source(fmt.Sprintf("s3://%s/broken/", bucket)),
		})
		if err != nil {
			return fmt.Errorf("StartImport failed: %w", err)
		}
		done, err := waitTerminal(*resp.ImportId)
		if err != nil {
			return err
		}
		if done.ImportStatus != types.ImportStatusFailed {
			return fmt.Errorf("ImportStatus = %s, want FAILED", done.ImportStatus)
		}
		failedImportID = *resp.ImportId
		return nil
	}))
	results = append(results, r.RunTest("cloudtrail", "ListImportFailures_Success", func() error {
		if failedImportID == "" {
			return fmt.Errorf("no failed import to list")
		}
		resp, err := tc.client.ListImportFailures(tc.ctx, &cloudtrail.ListImportFailuresInput{
			ImportId: aws.String(failedImportID),
		})
		if err != nil {
			return fmt.Errorf("ListImportFailures failed: %w", err)
		}
		if len(resp.Failures) == 0 {
			return fmt.Errorf("the corrupt file produced no failure entry")
		}
		f := resp.Failures[0]
		if !strings.Contains(aws.ToString(f.Location), "corrupt.json.gz") ||
			aws.ToString(f.ErrorType) == "" || aws.ToString(f.ErrorMessage) == "" {
			return fmt.Errorf("failure entry incomplete: %+v", f)
		}
		return nil
	}))

	// GetImport.
	results = append(results, r.RunTest("cloudtrail", "GetImport_Success", func() error {
		resp, err := tc.client.GetImport(tc.ctx, &cloudtrail.GetImportInput{
			ImportId: aws.String(importID),
		})
		if err != nil {
			return fmt.Errorf("GetImport failed: %w", err)
		}
		if resp.ImportId == nil || *resp.ImportId != importID {
			return fmt.Errorf("ImportId mismatch")
		}
		if resp.ImportStatus == "" {
			return fmt.Errorf("ImportStatus is empty")
		}
		return nil
	}))

	// GetImport_NotFound.
	results = append(results, r.RunTest("cloudtrail", "GetImport_NotFound", func() error {
		_, err := tc.client.GetImport(tc.ctx, &cloudtrail.GetImportInput{
			ImportId: aws.String("nonexistent-import-id"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent import")
		}
		return nil
	}))

	// ListImports.
	results = append(results, r.RunTest("cloudtrail", "ListImports_Success", func() error {
		resp, err := tc.client.ListImports(tc.ctx, &cloudtrail.ListImportsInput{
			Destination: aws.String(edsARN),
		})
		if err != nil {
			return fmt.Errorf("ListImports failed: %w", err)
		}
		found := false
		for _, imp := range resp.Imports {
			if imp.ImportId != nil && *imp.ImportId == importID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("import %s not found in list", importID)
		}
		return nil
	}))

	// StartImport_MissingDestinations.
	results = append(results, r.RunTest("cloudtrail", "StartImport_MissingDestinations", func() error {
		_, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			ImportSource: source(fmt.Sprintf("s3://%s/CloudTrail/", bucket)),
		})
		if err == nil {
			return fmt.Errorf("expected error for missing Destinations")
		}
		return nil
	}))

	return results
}

package testutil

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"vorpalstacks-sdk-tests/config"
)

// dynamoDBExportImportTests pins the export/import contract: exports and
// imports start in the IN_PROGRESS state and complete asynchronously, an
// export snapshots the table as of its ExportTime, and out-of-window
// export times are rejected.
func (r *TestRunner) dynamoDBExportImportTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "ExportImport_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("load config: %v", err),
		})
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })
	bucket := fmt.Sprintf("ddb-expimp-%d", suffix)
	if _, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "ExportImport_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create bucket: %v", err),
		})
	}
	defer cleanupS3Bucket(ctx, s3Client, bucket)

	// Build a table whose state differs before and after a chosen point:
	// two items exist at the point, a third appears afterwards.
	tableName := fmt.Sprintf("ExpImp-%d", suffix)
	tableArn, exportAt, seedErr := seedExportTable(ctx, client, tableName)
	if seedErr != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "ExportImport_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("seed table: %v", seedErr),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

	results = append(results, r.RunTest("dynamodb", "ExportTableToPointInTime_InitialStateInProgress", func() error {
		out, err := client.ExportTableToPointInTime(ctx, &dynamodb.ExportTableToPointInTimeInput{
			TableArn:   aws.String(tableArn),
			S3Bucket:   aws.String(bucket),
			ExportTime: aws.Time(exportAt),
		})
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}
		desc := out.ExportDescription
		if desc.ExportStatus != types.ExportStatusInProgress {
			return fmt.Errorf("initial ExportStatus = %q, want IN_PROGRESS (exports complete asynchronously)", desc.ExportStatus)
		}
		if desc.EndTime != nil {
			return fmt.Errorf("in-progress export must not report EndTime, got %v", desc.EndTime)
		}
		if desc.ExportTime == nil || diffSeconds(*desc.ExportTime, exportAt) > 2 {
			return fmt.Errorf("ExportTime = %v, want ~%v", desc.ExportTime, exportAt)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExportTableToPointInTime_SnapshotsAsOfExportTime", func() error {
		out, err := client.ExportTableToPointInTime(ctx, &dynamodb.ExportTableToPointInTimeInput{
			TableArn:   aws.String(tableArn),
			S3Bucket:   aws.String(bucket),
			ExportTime: aws.Time(exportAt),
		})
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}
		exportArn := aws.ToString(out.ExportDescription.ExportArn)

		desc, err := pollExportStatus(ctx, client, exportArn)
		if err != nil {
			return err
		}
		if desc.ItemCount == nil || *desc.ItemCount != 2 {
			return fmt.Errorf("as-of export ItemCount = %d, want 2 (the state at the export time)", aws.ToInt64(desc.ItemCount))
		}
		if desc.BilledSizeBytes == nil || *desc.BilledSizeBytes == 0 {
			return fmt.Errorf("completed export must report BilledSizeBytes")
		}

		// The exported data must hold the as-of values. The completed
		// export's ExportManifest names manifest-files.json, whose JSON
		// lines point at the gzip-compressed DYNAMODB_JSON data files —
		// the documented consumer path through the export layout.
		manifestURL := aws.ToString(desc.ExportManifest)
		if manifestURL == "" {
			return fmt.Errorf("completed export must report ExportManifest")
		}
		manifestKey := strings.TrimPrefix(manifestURL, "s3://"+bucket+"/")
		manifestObj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(manifestKey)})
		if err != nil {
			return fmt.Errorf("get files manifest %s: %w", manifestKey, err)
		}
		type manifestEntry struct {
			ItemCount     int64  `json:"itemCount"`
			DataFileS3Key string `json:"dataFileS3Key"`
		}
		// manifest-files.json is JSON lines: one entry object per line.
		var entries []manifestEntry
		manifestDec := json.NewDecoder(manifestObj.Body)
		for {
			var entry manifestEntry
			if err := manifestDec.Decode(&entry); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				manifestObj.Body.Close()
				return fmt.Errorf("decode files manifest: %w", err)
			}
			entries = append(entries, entry)
		}
		manifestObj.Body.Close()
		if len(entries) == 0 {
			return fmt.Errorf("files manifest lists no data files")
		}
		if entries[0].ItemCount != 2 {
			return fmt.Errorf("manifest itemCount = %d, want 2", entries[0].ItemCount)
		}

		values := map[string]string{}
		for _, entry := range entries {
			obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(entry.DataFileS3Key)})
			if err != nil {
				return fmt.Errorf("get exported data %s: %w", entry.DataFileS3Key, err)
			}
			gz, gzErr := gzip.NewReader(obj.Body)
			if gzErr != nil {
				obj.Body.Close()
				return fmt.Errorf("data file %s is not gzip: %w", entry.DataFileS3Key, gzErr)
			}
			dec := json.NewDecoder(gz)
			for {
				var line struct {
					Item map[string]map[string]string
				}
				if err := dec.Decode(&line); err != nil {
					break
				}
				if id, ok := line.Item["id"]; ok {
					if v, ok := line.Item["v"]; ok {
						values[id["S"]] = v["S"]
					}
				}
			}
			gz.Close()
			obj.Body.Close()
		}
		if len(values) != 2 || values["k1"] != "v1" || values["k2"] != "keep" {
			return fmt.Errorf("exported items = %v, want the as-of state {k1:v1, k2:keep}", values)
		}

		list, err := client.ListExports(ctx, &dynamodb.ListExportsInput{TableArn: aws.String(tableArn)})
		if err != nil {
			return fmt.Errorf("list exports: %w", err)
		}
		found := false
		for _, summary := range list.ExportSummaries {
			if aws.ToString(summary.ExportArn) == exportArn {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("ListExports must contain the export %s", exportArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExportTableToPointInTime_CurrentStateWhenTimeOmitted", func() error {
		out, err := client.ExportTableToPointInTime(ctx, &dynamodb.ExportTableToPointInTimeInput{
			TableArn: aws.String(tableArn),
			S3Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}
		desc, err := pollExportStatus(ctx, client, aws.ToString(out.ExportDescription.ExportArn))
		if err != nil {
			return err
		}
		if desc.ItemCount == nil || *desc.ItemCount != 3 {
			return fmt.Errorf("current-state export ItemCount = %d, want 3 (the present state)", aws.ToInt64(desc.ItemCount))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExportTableToPointInTime_OutOfWindow_Rejected", func() error {
		_, err := client.ExportTableToPointInTime(ctx, &dynamodb.ExportTableToPointInTimeInput{
			TableArn:   aws.String(tableArn),
			S3Bucket:   aws.String(bucket),
			ExportTime: aws.Time(time.Now().Add(-time.Hour)),
		})
		var invalidTime *types.InvalidExportTimeException
		if err == nil {
			return fmt.Errorf("export before the restorable window must fail")
		}
		if !errors.As(err, &invalidTime) {
			return fmt.Errorf("error = %v, want InvalidExportTimeException", err)
		}
		return nil
	}))

	importPrefix := fmt.Sprintf("import-%d/", suffix)
	if err := writeImportSource(ctx, s3Client, bucket, importPrefix); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "ImportTable_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("write import source: %v", err),
		})
	}
	importTable := fmt.Sprintf("ExpImpImport-%d", suffix)
	results = append(results, r.RunTest("dynamodb", "ImportTable_CompletesAsynchronously", func() error {
		out, err := client.ImportTable(ctx, &dynamodb.ImportTableInput{
			InputFormat: types.InputFormatDynamodbJson,
			S3BucketSource: &types.S3BucketSource{
				S3Bucket:    aws.String(bucket),
				S3KeyPrefix: aws.String(importPrefix),
			},
			TableCreationParameters: &types.TableCreationParameters{
				TableName: aws.String(importTable),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				},
				BillingMode: types.BillingModePayPerRequest,
			},
		})
		if err != nil {
			return fmt.Errorf("import: %w", err)
		}
		desc := out.ImportTableDescription
		if desc.ImportStatus != types.ImportStatusInProgress {
			return fmt.Errorf("initial ImportStatus = %q, want IN_PROGRESS (imports complete asynchronously)", desc.ImportStatus)
		}
		if desc.EndTime != nil {
			return fmt.Errorf("in-progress import must not report EndTime, got %v", desc.EndTime)
		}

		importArn := aws.ToString(desc.ImportArn)
		final, err := pollImportStatus(ctx, client, importArn)
		if err != nil {
			return err
		}
		if final.ImportStatus != types.ImportStatusCompleted {
			return fmt.Errorf("import status = %q (failure %q: %s), want COMPLETED",
				final.ImportStatus, aws.ToString(final.FailureCode), aws.ToString(final.FailureMessage))
		}
		if final.ProcessedItemCount != 2 {
			return fmt.Errorf("ProcessedItemCount = %d, want 2", final.ProcessedItemCount)
		}

		item, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName:      aws.String(importTable),
			Key:            map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "i1"}},
			ConsistentRead: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("get imported item: %w", err)
		}
		if attr, ok := item.Item["v"].(*types.AttributeValueMemberS); !ok || attr.Value != "one" {
			return fmt.Errorf("imported item value = %+v, want one", item.Item["v"])
		}

		importedTableArn := aws.ToString(final.TableArn)
		list, err := client.ListImports(ctx, &dynamodb.ListImportsInput{TableArn: aws.String(importedTableArn)})
		if err != nil {
			return fmt.Errorf("list imports: %w", err)
		}
		found := false
		for _, summary := range list.ImportSummaryList {
			if aws.ToString(summary.ImportArn) == importArn {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("ListImports must contain the import %s", importArn)
		}
		return nil
	}))
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(importTable)})

	results = append(results, r.RunTest("dynamodb", "ImportTable_ItemValidationFailures", func() error {
		// Source data with a duplicate key (a documented overwrite, not an
		// error), a wrong key-type line, and an unparsable line.
		badPrefix := fmt.Sprintf("import-bad-%d/", suffix)
		body := strings.Join([]string{
			`{"Item":{"id":{"N":"1"},"v":{"S":"one"}}}`,
			`{"Item":{"id":{"N":"1"},"v":{"S":"two"}}}`,
			`{"Item":{"id":{"S":"x"},"v":{"S":"bad"}}}`,
			`garbage-not-json`,
		}, "\n") + "\n"
		if _, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(badPrefix + "data.json"),
			Body:   strings.NewReader(body),
		}); err != nil {
			return fmt.Errorf("write import source: %w", err)
		}

		badTable := fmt.Sprintf("ExpImpBadItems-%d", suffix)
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(badTable)})
		out, err := client.ImportTable(ctx, &dynamodb.ImportTableInput{
			InputFormat: types.InputFormatDynamodbJson,
			S3BucketSource: &types.S3BucketSource{
				S3Bucket:    aws.String(bucket),
				S3KeyPrefix: aws.String(badPrefix),
			},
			TableCreationParameters: &types.TableCreationParameters{
				TableName: aws.String(badTable),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeN},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				},
				BillingMode: types.BillingModePayPerRequest,
			},
		})
		if err != nil {
			return fmt.Errorf("import: %w", err)
		}
		final, err := pollImportStatus(ctx, client, aws.ToString(out.ImportTableDescription.ImportArn))
		if err != nil {
			return err
		}
		// Items that fail validation are skipped and counted; the job ends
		// FAILED with the item-validation failure code.
		if final.ImportStatus != types.ImportStatusFailed {
			return fmt.Errorf("import status = %q, want FAILED when items fail validation", final.ImportStatus)
		}
		if aws.ToString(final.FailureCode) != "ItemValidationError" {
			return fmt.Errorf("failure code = %q, want ItemValidationError", aws.ToString(final.FailureCode))
		}
		// Three lines parsed into items (both N-keyed lines plus the
		// mistyped one); two loaded (the duplicate overwrites), two errors
		// (mistyped key, unparsable line).
		if final.ProcessedItemCount != 3 {
			return fmt.Errorf("ProcessedItemCount = %d, want 3", final.ProcessedItemCount)
		}
		if final.ImportedItemCount != 2 {
			return fmt.Errorf("ImportedItemCount = %d, want 2", final.ImportedItemCount)
		}
		if final.ErrorCount != 2 {
			return fmt.Errorf("ErrorCount = %d, want 2", final.ErrorCount)
		}

		// The table holds the single distinct item with the last write's
		// value; the mistyped line stored nothing.
		desc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(badTable)})
		if err != nil {
			return fmt.Errorf("describe table: %w", err)
		}
		if aws.ToInt64(desc.Table.ItemCount) != 1 {
			return fmt.Errorf("table ItemCount = %d, want 1 distinct item (overwrites must not double-count)", aws.ToInt64(desc.Table.ItemCount))
		}
		item, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName:      aws.String(badTable),
			Key:            map[string]types.AttributeValue{"id": &types.AttributeValueMemberN{Value: "1"}},
			ConsistentRead: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("get imported item: %w", err)
		}
		if attr, ok := item.Item["v"].(*types.AttributeValueMemberS); !ok || attr.Value != "two" {
			return fmt.Errorf("imported item value = %+v, want two (last write wins)", item.Item["v"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ImportTable_UnknownBucket_Fails", func() error {
		out, err := client.ImportTable(ctx, &dynamodb.ImportTableInput{
			InputFormat: types.InputFormatDynamodbJson,
			S3BucketSource: &types.S3BucketSource{
				S3Bucket: aws.String(fmt.Sprintf("ddb-no-such-bucket-%d", suffix)),
			},
			TableCreationParameters: &types.TableCreationParameters{
				TableName: aws.String(fmt.Sprintf("ExpImpBad-%d", suffix)),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				},
				BillingMode: types.BillingModePayPerRequest,
			},
		})
		if err != nil {
			return fmt.Errorf("import: %w", err)
		}
		final, err := pollImportStatus(ctx, client, aws.ToString(out.ImportTableDescription.ImportArn))
		if err != nil {
			return err
		}
		if final.ImportStatus != types.ImportStatusFailed {
			return fmt.Errorf("import from a missing bucket must fail, got %q", final.ImportStatus)
		}
		if final.FailureCode == nil {
			return fmt.Errorf("failed import must report a failure code")
		}
		return nil
	}))
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(fmt.Sprintf("ExpImpBad-%d", suffix))})

	return results
}

// seedExportTable creates a table with recovery enabled, writes two items,
// and after a margin writes a third item. The returned export point lies
// between the two writes; the caller deletes the table.
func seedExportTable(ctx context.Context, client *dynamodb.Client, tableName string) (string, time.Time, error) {
	out, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("create table: %w", err)
	}
	tableArn := aws.ToString(out.TableDescription.TableArn)
	if _, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
		TableName: aws.String(tableName),
		PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: aws.Bool(true),
		},
	}); err != nil {
		return "", time.Time{}, fmt.Errorf("enable recovery: %w", err)
	}

	put := func(id, v string) error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: id},
				"v":  &types.AttributeValueMemberS{Value: v},
			},
		})
		return err
	}
	if err := put("k1", "v1"); err != nil {
		return "", time.Time{}, fmt.Errorf("put k1: %w", err)
	}
	if err := put("k2", "keep"); err != nil {
		return "", time.Time{}, fmt.Errorf("put k2: %w", err)
	}
	// The export point must separate the two "at the point" writes from
	// the later one; DynamoDB timestamps carry second-level granularity,
	// so the margin exceeds one second on both sides.
	time.Sleep(1100 * time.Millisecond)
	exportAt := time.Now()
	time.Sleep(1100 * time.Millisecond)
	if err := put("k3", "later"); err != nil {
		return "", time.Time{}, fmt.Errorf("put k3: %w", err)
	}
	return tableArn, exportAt, nil
}

// pollExportStatus polls DescribeExport until the export reaches a final
// state and returns the final description.
func pollExportStatus(ctx context.Context, client *dynamodb.Client, exportArn string) (*types.ExportDescription, error) {
	deadline := time.Now().Add(20 * time.Second)
	for {
		out, err := client.DescribeExport(ctx, &dynamodb.DescribeExportInput{ExportArn: aws.String(exportArn)})
		if err != nil {
			return nil, fmt.Errorf("describe export: %w", err)
		}
		if out.ExportDescription.ExportStatus != types.ExportStatusInProgress {
			return out.ExportDescription, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("export %s did not finish in time", exportArn)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

// pollImportStatus polls DescribeImport until the import reaches a final
// state and returns the final description.
func pollImportStatus(ctx context.Context, client *dynamodb.Client, importArn string) (*types.ImportTableDescription, error) {
	deadline := time.Now().Add(20 * time.Second)
	for {
		out, err := client.DescribeImport(ctx, &dynamodb.DescribeImportInput{ImportArn: aws.String(importArn)})
		if err != nil {
			return nil, fmt.Errorf("describe import: %w", err)
		}
		if out.ImportTableDescription.ImportStatus != types.ImportStatusInProgress {
			return out.ImportTableDescription, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("import %s did not finish in time", importArn)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

// writeImportSource uploads a two-item DynamoDB JSON data file for import.
func writeImportSource(ctx context.Context, s3Client *s3.Client, bucket, prefix string) error {
	body := strings.Join([]string{
		`{"Item":{"id":{"S":"i1"},"v":{"S":"one"}}}`,
		`{"Item":{"id":{"S":"i2"},"v":{"S":"two"}}}`,
	}, "\n") + "\n"
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix + "data.json"),
		Body:   strings.NewReader(body),
	})
	return err
}

// cleanupS3Bucket removes every object in the bucket and then the bucket.
func cleanupS3Bucket(ctx context.Context, s3Client *s3.Client, bucket string) {
	var markers []string
	out, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		return
	}
	for _, obj := range out.Contents {
		markers = append(markers, aws.ToString(obj.Key))
	}
	for _, key := range markers {
		s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	}
	s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
}

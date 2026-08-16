package testutil

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwltypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *TestRunner) runCWLogsToLambda(ic *integClients, ts string) TestResult {
	fnName := fmt.Sprintf("integ-cwl-lambda-%s", ts)
	roleName := fmt.Sprintf("integ-cwl-role-%s", ts)
	logGroupName := fmt.Sprintf("/integ/cwl-lambda/%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	ic.createLambda(fnName, roleName)
	defer ic.deleteLambda(fnName)

	fnARN := fmt.Sprintf("arn:aws:lambda:%s:000000000000:function:%s", ic.region, fnName)

	ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)})
	defer func() {
		ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
	}()

	_, err := ic.cwl.PutSubscriptionFilter(ic.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(logGroupName),
		FilterName:     aws.String("integ-lambda-sub"),
		FilterPattern:  aws.String("[...]"),
		DestinationArn: aws.String(fnARN),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Lambda", func() error { return fmt.Errorf("put subscription filter: %w", err) })
	}

	ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
	})

	ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String("integration test log message"), Timestamp: aws.Int64(time.Now().UnixMilli())},
		},
	})

	return r.pollVerify("CWLogs_Lambda", defaultPollTimeout, func() error {
		return ic.verifyLambdaInvoked(fnName)
	})
}

func (r *TestRunner) runCWLogsToKinesis(ic *integClients, ts string) TestResult {
	streamName := fmt.Sprintf("integ-cwl-kin-%s", ts)
	logGroupName := fmt.Sprintf("/integ/cwl-kinesis/%s", ts)

	err := ic.createKinesisStream(streamName)
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("create stream: %w", err) })
	}
	defer ic.deleteStream(streamName)

	if err := ic.pollStreamActive(streamName, defaultPollTimeout); err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("stream not active: %w", err) })
	}

	streamARN := fmt.Sprintf("arn:aws:kinesis:%s:000000000000:stream/%s", ic.region, streamName)

	ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)})
	defer func() {
		ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
	}()

	ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
	})

	_, err = ic.cwl.PutSubscriptionFilter(ic.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(logGroupName),
		FilterName:     aws.String("integ-kinesis-sub"),
		FilterPattern:  aws.String("[...]"),
		DestinationArn: aws.String(streamARN),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("put subscription filter: %w", err) })
	}

	ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String("kinesis subscription test"), Timestamp: aws.Int64(time.Now().UnixMilli())},
		},
	})

	return r.pollVerify("CWLogs_Kinesis", defaultPollTimeout, func() error {
		streamDesc, err := ic.describeStream(streamName)
		if err != nil {
			return fmt.Errorf("describe stream: %w", err)
		}
		if len(streamDesc.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards in stream")
		}
		shardID := streamDesc.StreamDescription.Shards[0].ShardId
		iter, err := ic.createIteratorFromHorizon(streamName, *shardID)
		if err != nil {
			return fmt.Errorf("create iterator: %w", err)
		}
		records, err := ic.getRecords(iter)
		if err != nil {
			return fmt.Errorf("get records: %w", err)
		}
		if len(records.Records) == 0 {
			return fmt.Errorf("expected records from CW Logs subscription, got 0")
		}
		if !strings.Contains(string(records.Records[0].Data), "awslogs") {
			return fmt.Errorf("expected awslogs envelope, got: %s", string(records.Records[0].Data))
		}
		return nil
	})
}

// runCWLogsLookupTableKMS verifies that a lookup table created with a
// customer-managed KMS key round-trips its content and still serves lookup
// commands in queries: the body is encrypted at rest and decrypted at the
// API and query boundaries.
func (r *TestRunner) runCWLogsLookupTableKMS(ic *integClients, ts string) TestResult {
	testName := "CWLogs_LookupTable_KMS"

	keyResp, err := ic.kms.CreateKey(ic.ctx, &kms.CreateKeyInput{
		Description: aws.String("integ lookup table key " + ts),
	})
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create key: %w", err) })
	}
	keyID := aws.ToString(keyResp.KeyMetadata.KeyId)
	defer func() {
		_, _ = ic.kms.ScheduleKeyDeletion(ic.ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               aws.String(keyID),
			PendingWindowInDays: aws.Int32(7),
		})
	}()

	tableName := fmt.Sprintf("integ_kms_table_%s", ts)
	body := "id,name\nk1,Encrypted\n"
	createResp, err := ic.cwl.CreateLookupTable(ic.ctx, &cloudwatchlogs.CreateLookupTableInput{
		LookupTableName: aws.String(tableName),
		TableBody:       aws.String(body),
		KmsKeyId:        aws.String(keyID),
	})
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create lookup table: %w", err) })
	}
	tableArn := aws.ToString(createResp.LookupTableArn)
	defer func() {
		_, _ = ic.cwl.DeleteLookupTable(ic.ctx, &cloudwatchlogs.DeleteLookupTableInput{
			LookupTableArn: aws.String(tableArn),
		})
	}()

	err = func() error {
		getResp, err := ic.cwl.GetLookupTable(ic.ctx, &cloudwatchlogs.GetLookupTableInput{
			LookupTableArn: aws.String(tableArn),
		})
		if err != nil {
			return fmt.Errorf("get lookup table: %w", err)
		}
		if aws.ToString(getResp.TableBody) != body {
			return fmt.Errorf("encrypted table body round trip: %q", aws.ToString(getResp.TableBody))
		}
		if aws.ToString(getResp.KmsKeyId) != keyID {
			return fmt.Errorf("kmsKeyId = %q, want %q", aws.ToString(getResp.KmsKeyId), keyID)
		}

		// A lookup command against the encrypted table must enrich events.
		logGroupName := fmt.Sprintf("/integ/cwl-kms-lookup/%s", ts)
		_, err = ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)})
		if err != nil {
			return fmt.Errorf("create log group: %w", err)
		}
		defer ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
		_, _ = ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String("s1"),
		})
		now := time.Now().UnixMilli()
		_, _ = ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String("s1"),
			LogEvents:     []cwltypes.InputLogEvent{{Message: aws.String(`{"uid":"k1"}`), Timestamp: aws.Int64(now)}},
		})

		startResp, err := ic.cwl.StartQuery(ic.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 60000),
			EndTime:       aws.Int64(now + 60000),
			LogGroupNames: []string{logGroupName},
			QueryString:   aws.String(fmt.Sprintf(`lookup %s id as uid OUTPUT name | fields uid, name`, tableName)),
		})
		if err != nil {
			return fmt.Errorf("start query: %w", err)
		}
		var rows [][]cwltypes.ResultField
		for i := 0; i < 20; i++ {
			resResp, err := ic.cwl.GetQueryResults(ic.ctx, &cloudwatchlogs.GetQueryResultsInput{QueryId: startResp.QueryId})
			if err != nil {
				return fmt.Errorf("get query results: %w", err)
			}
			if resResp.Status == cwltypes.QueryStatusComplete {
				rows = resResp.Results
				break
			}
			if resResp.Status == cwltypes.QueryStatusFailed || resResp.Status == cwltypes.QueryStatusCancelled {
				return fmt.Errorf("query status %s", resResp.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if len(rows) != 1 {
			return fmt.Errorf("rows = %d, want 1", len(rows))
		}
		for _, f := range rows[0] {
			if aws.ToString(f.Field) == "name" && aws.ToString(f.Value) != "Encrypted" {
				return fmt.Errorf("encrypted table enrichment: name=%q", aws.ToString(f.Value))
			}
		}
		return nil
	}()
	return integResult(testName, err)
}

// runCWLogsScheduledQueryS3 verifies that a scheduled query with an S3
// destination writes its gzipped CSV results under the configured S3 URI and
// reports the delivery through GetScheduledQueryHistory.
func (r *TestRunner) runCWLogsScheduledQueryS3(ic *integClients, ts string) TestResult {
	testName := "CWLogs_ScheduledQuery_S3"

	bucket := fmt.Sprintf("integ-cwl-sched-%s", ts)
	if _, err := ic.s3.CreateBucket(ic.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create bucket: %w", err) })
	}
	defer ic.s3.DeleteBucket(ic.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})

	logGroupName := fmt.Sprintf("/integ/cwl-sched-s3/%s", ts)
	if _, err := ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)}); err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create log group: %w", err) })
	}
	defer ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
	if _, err := ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("s1"),
	}); err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create log stream: %w", err) })
	}
	now := time.Now().UnixMilli()
	if _, err := ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("s1"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String(`{"code":500}`), Timestamp: aws.Int64(now)},
		},
	}); err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("put events: %w", err) })
	}

	createResp, err := ic.cwl.CreateScheduledQuery(ic.ctx, &cloudwatchlogs.CreateScheduledQueryInput{
		Name:                aws.String(fmt.Sprintf("integ-sched-s3-%s", ts)),
		QueryString:         aws.String("fields code"),
		QueryLanguage:       cwltypes.QueryLanguageCwli,
		LogGroupIdentifiers: []string{logGroupName},
		ExecutionRoleArn:    aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
		ScheduleExpression:  aws.String("rate(1 minute)"),
		State:               cwltypes.ScheduledQueryStateEnabled,
		DestinationConfiguration: &cwltypes.DestinationConfiguration{
			S3Configuration: &cwltypes.S3Configuration{
				DestinationIdentifier: aws.String(fmt.Sprintf("s3://%s/results/%s", bucket, ts)),
				RoleArn:               aws.String("arn:aws:iam::123456789012:role/deliver"),
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create scheduled query: %w", err) })
	}
	queryArn := aws.ToString(createResp.ScheduledQueryArn)
	defer ic.cwl.DeleteScheduledQuery(ic.ctx, &cloudwatchlogs.DeleteScheduledQueryInput{Identifier: aws.String(queryArn)})

	// The scheduled query worker fires on its one-minute tick; wait for the
	// delivery to appear under the configured prefix.
	var objectKey string
	res := r.pollVerify(testName, 100*time.Second, func() error {
		listResp, err := ic.s3.ListObjectsV2(ic.ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String("results/" + ts + "/"),
		})
		if err != nil {
			return fmt.Errorf("list objects: %w", err)
		}
		if len(listResp.Contents) == 0 {
			return fmt.Errorf("no delivered objects yet")
		}
		objectKey = aws.ToString(listResp.Contents[0].Key)
		return nil
	})
	if res.Status != "PASS" {
		return res
	}

	verifyErr := func() error {
		getObj, err := ic.s3.GetObject(ic.ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		})
		if err != nil {
			return fmt.Errorf("get object: %w", err)
		}
		defer getObj.Body.Close()
		zr, err := gzip.NewReader(getObj.Body)
		if err != nil {
			return fmt.Errorf("gzip reader: %w", err)
		}
		csv, err := io.ReadAll(zr)
		if err != nil {
			return fmt.Errorf("read gzip: %w", err)
		}
		if !strings.Contains(string(csv), "500") {
			return fmt.Errorf("delivered CSV missing query results: %q", string(csv))
		}

		history, err := ic.cwl.GetScheduledQueryHistory(ic.ctx, &cloudwatchlogs.GetScheduledQueryHistoryInput{
			Identifier: aws.String(queryArn),
			StartTime:  aws.Int64(now - 60000),
			EndTime:    aws.Int64(time.Now().Add(time.Hour).UnixMilli()),
		})
		if err != nil {
			return fmt.Errorf("get history: %w", err)
		}
		for _, trig := range history.TriggerHistory {
			for _, dest := range trig.Destinations {
				if dest.DestinationType == cwltypes.ScheduledQueryDestinationTypeS3 &&
					dest.Status == cwltypes.ActionStatusComplete {
					return nil
				}
			}
		}
		return fmt.Errorf("history does not report a completed S3 delivery: %+v", history.TriggerHistory)
	}()
	return integResult(testName, verifyErr)
}

// integResult builds the delivery test result without a RunTest wrapper:
// the integration harness already prints the Running and result lines.
func integResult(testName string, err error) TestResult {
	res := TestResult{Service: integSvc, TestName: testName, Status: "PASS"}
	if err != nil {
		res.Status = "FAIL"
		res.Error = err.Error()
	}
	return res
}

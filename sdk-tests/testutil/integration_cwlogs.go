package testutil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
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

	fnARN := fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", ic.region, ic.accountID, fnName)

	if _, err := ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)}); err != nil {
		return r.RunTest(integSvc, "CWLogs_Lambda", func() error { return fmt.Errorf("create log group: %w", err) })
	}
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

	if _, err := ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
	}); err != nil {
		return r.RunTest(integSvc, "CWLogs_Lambda", func() error { return fmt.Errorf("create log stream: %w", err) })
	}
	if _, err := ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String("integration test log message"), Timestamp: aws.Int64(time.Now().UnixMilli())},
		},
	}); err != nil {
		return r.RunTest(integSvc, "CWLogs_Lambda", func() error { return fmt.Errorf("put log events: %w", err) })
	}

	return r.pollVerify("CWLogs_Lambda", defaultPollTimeout, func() error {
		if err := ic.verifyLambdaInvoked(fnName); err != nil {
			return err
		}
		return ic.verifyLambdaAwslogsPayload(fnName, "integration test log message")
	})
}

// verifyLambdaAwslogsPayload decodes the subscription payload the function
// received: the handler echoes its invocation event, whose awslogs.data
// member is the base64-wrapped gzip of the DATA_MESSAGE record — the
// documented Lambda invocation format for log-group-level subscription
// filters. Asserting on the decoded record (not just the invocation)
// verifies the delivered content. The scan reads every stream of the
// function's log group so it holds under either stream layout: invocations
// of one execution environment share its stream, and invocations across
// separate environments land in separate streams.
func (ic *integClients) verifyLambdaAwslogsPayload(fnName, message string) error {
	found, err := ic.scanLambdaLogStreams(fnName, func(msg string) (bool, error) {
		var envelope struct {
			Awslogs struct {
				Data string `json:"data"`
			} `json:"awslogs"`
		}
		if err := json.Unmarshal([]byte(msg), &envelope); err != nil || envelope.Awslogs.Data == "" {
			return false, nil
		}
		compressed, err := base64.StdEncoding.DecodeString(envelope.Awslogs.Data)
		if err != nil {
			return false, fmt.Errorf("awslogs.data is not base64: %w", err)
		}
		zr, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			return false, fmt.Errorf("awslogs.data does not carry gzip: %w", err)
		}
		payload, err := io.ReadAll(zr)
		if err != nil {
			return false, fmt.Errorf("read gzip: %w", err)
		}
		// The reachability probe CloudWatch Logs may send when the
		// subscription is created echoes here too, ahead of any data
		// record; it is a documented CONTROL_MESSAGE with no events, not
		// the delivery under test — keep scanning for the DATA record.
		if strings.Contains(string(payload), `"messageType":"CONTROL_MESSAGE"`) {
			return false, nil
		}
		if !strings.Contains(string(payload), `"messageType":"DATA_MESSAGE"`) {
			return false, nil
		}
		if !strings.Contains(string(payload), message) {
			return false, fmt.Errorf("DATA_MESSAGE record does not carry the logged event: %s", string(payload))
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("lambda %s: no invocation log carried a DATA_MESSAGE awslogs envelope", fnName)
	}
	return nil
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

	streamARN := fmt.Sprintf("arn:aws:kinesis:%s:%s:stream/%s", ic.region, ic.accountID, streamName)

	if _, err := ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroupName)}); err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("create log group: %w", err) })
	}
	defer func() {
		ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroupName)})
	}()

	if _, err := ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
	}); err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("create log stream: %w", err) })
	}

	_, err = ic.cwl.PutSubscriptionFilter(ic.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(logGroupName),
		FilterName:     aws.String("integ-kinesis-sub"),
		FilterPattern:  aws.String("[...]"),
		DestinationArn: aws.String(streamARN),
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("put subscription filter: %w", err) })
	}

	if _, err := ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String("test-stream"),
		LogEvents: []cwltypes.InputLogEvent{
			{Message: aws.String("kinesis subscription test"), Timestamp: aws.Int64(time.Now().UnixMilli())},
		},
	}); err != nil {
		return r.RunTest(integSvc, "CWLogs_Kinesis", func() error { return fmt.Errorf("put log events: %w", err) })
	}

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
		// The stream may carry the CONTROL_MESSAGE reachability record the
		// subscription's creation emits ("Sometimes CloudWatch Logs may
		// emit Amazon Kinesis Data Streams records with a 'CONTROL_MESSAGE'
		// type, mainly for checking if the destination is reachable"),
		// ahead of the data batch — scan for both shapes. The record's data
		// is the gzipped payload as-is — one gzip layer, no JSON envelope
		// and no extra base64 wrap (the {"awslogs":{"data":...}} envelope
		// is the Lambda invocation format; the CloudWatch Logs User Guide's
		// Kinesis example reads the record with one base64-decode of the
		// API presentation followed by zcat).
		sawControl := false
		sawData := false
		for _, rec := range records.Records {
			zr, err := gzip.NewReader(bytes.NewReader(rec.Data))
			if err != nil {
				return fmt.Errorf("record data must be a single gzip layer: %w", err)
			}
			payload, err := io.ReadAll(zr)
			if err != nil {
				return fmt.Errorf("read gzip: %w", err)
			}
			switch {
			case strings.Contains(string(payload), `"CONTROL_MESSAGE"`):
				sawControl = true
				if !strings.Contains(string(payload), `"logEvents":[]`) {
					return fmt.Errorf("control record must carry no log events, got: %s", string(payload))
				}
			case strings.Contains(string(payload), `"DATA_MESSAGE"`) && strings.Contains(string(payload), "kinesis subscription test"):
				sawData = true
			}
		}
		if !sawData {
			return fmt.Errorf("records must include the DATA_MESSAGE payload carrying the logged event, got %d records", len(records.Records))
		}
		if !sawControl {
			return fmt.Errorf("records must include the CONTROL_MESSAGE reachability probe, got %d records", len(records.Records))
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
			StartTime:     aws.Int64(now/1000 - 60),
			EndTime:       aws.Int64(now/1000 + 60),
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
		ExecutionRoleArn:    aws.String(fmt.Sprintf("arn:aws:iam::%s:role/scheduled-query-role", ic.accountID)),
		ScheduleExpression:  aws.String("rate(1 minute)"),
		State:               cwltypes.ScheduledQueryStateEnabled,
		DestinationConfiguration: &cwltypes.DestinationConfiguration{
			S3Configuration: &cwltypes.S3Configuration{
				DestinationIdentifier: aws.String(fmt.Sprintf("s3://%s/results/%s", bucket, ts)),
				RoleArn:               aws.String(fmt.Sprintf("arn:aws:iam::%s:role/deliver", ic.accountID)),
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, testName, func() error { return fmt.Errorf("create scheduled query: %w", err) })
	}
	queryArn := aws.ToString(createResp.ScheduledQueryArn)
	defer ic.cwl.DeleteScheduledQuery(ic.ctx, &cloudwatchlogs.DeleteScheduledQueryInput{Identifier: aws.String(queryArn)})

	// The AWS rate() contract runs the first query one full interval
	// after creation; under the regression server's TEST_MODE the period
	// compresses to seconds, so the delivery appears shortly after
	// creation. Wait for it under the configured prefix.
	var objectKey string
	res := r.pollVerify(testName, 30*time.Second, func() error {
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

// runCWLogsVendedDeliveryS3 pins the vended-logs delivery engine's S3
// leg: a delivery source over a log group paired with an S3 delivery
// destination delivers the group's events as one gzipped JSON-lines
// object under the documented AWSLogs/<account-id>/ prefix.
func (r *TestRunner) runCWLogsVendedDeliveryS3(ic *integClients, ts string) TestResult {
	srcName := fmt.Sprintf("vended-s3-src-%s", ts)
	bucket := fmt.Sprintf("vended-delivery-%s", ts)
	srcGroup := fmt.Sprintf("/integ/vended-s3/%s", ts)
	const stream = "s3-delivery/1"

	_, err := ic.s3.CreateBucket(ic.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_VendedDelivery_S3", func() error {
			return fmt.Errorf("create bucket: %w", err)
		})
	}
	defer func() {
		list, _ := ic.s3.ListObjectsV2(ic.ctx, &s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
		for _, o := range list.Contents {
			ic.s3.DeleteObject(ic.ctx, &s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: o.Key})
		}
		ic.s3.DeleteBucket(ic.ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucket)})
	}()

	ic.cwl.CreateLogGroup(ic.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(srcGroup)})
	defer ic.cwl.DeleteLogGroup(ic.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(srcGroup)})
	groups, err := ic.cwl.DescribeLogGroups(ic.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(srcGroup),
	})
	if err != nil || len(groups.LogGroups) == 0 || groups.LogGroups[0].Arn == nil {
		return r.RunTest(integSvc, "CWLogs_VendedDelivery_S3", func() error {
			return fmt.Errorf("resolve source group ARN: %v", err)
		})
	}
	srcArn := groups.LogGroups[0].Arn

	ic.cwl.CreateLogStream(ic.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName: aws.String(srcGroup), LogStreamName: aws.String(stream),
	})
	now := time.Now().UnixMilli()
	for i, msg := range []string{"vended-first", "vended-second"} {
		_, err := ic.cwl.PutLogEvents(ic.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName: aws.String(srcGroup), LogStreamName: aws.String(stream),
			LogEvents: []cwltypes.InputLogEvent{
				{Message: aws.String(msg), Timestamp: aws.Int64(now - int64(2000-i*1000))},
			},
		})
		if err != nil {
			return r.RunTest(integSvc, "CWLogs_VendedDelivery_S3", func() error {
				return fmt.Errorf("put source event %d: %w", i, err)
			})
		}
	}

	if _, err := ic.cwl.PutDeliverySource(ic.ctx, &cloudwatchlogs.PutDeliverySourceInput{
		Name: aws.String(srcName), ResourceArn: srcArn, LogType: aws.String("APPLICATION_LOGS"),
	}); err != nil {
		return r.RunTest(integSvc, "CWLogs_VendedDelivery_S3", func() error {
			return fmt.Errorf("put delivery source: %w", err)
		})
	}
	defer ic.cwl.DeleteDeliverySource(ic.ctx, &cloudwatchlogs.DeleteDeliverySourceInput{Name: aws.String(srcName)})

	dest, err := ic.cwl.PutDeliveryDestination(ic.ctx, &cloudwatchlogs.PutDeliveryDestinationInput{
		Name: aws.String(srcName + "-dest"),
		DeliveryDestinationConfiguration: &cwltypes.DeliveryDestinationConfiguration{
			DestinationResourceArn: aws.String("arn:aws:s3:::" + bucket),
		},
		OutputFormat: cwltypes.OutputFormatJson,
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_VendedDelivery_S3", func() error {
			return fmt.Errorf("put delivery destination: %w", err)
		})
	}
	defer ic.cwl.DeleteDeliveryDestination(ic.ctx, &cloudwatchlogs.DeleteDeliveryDestinationInput{
		Name: aws.String(srcName + "-dest"),
	})

	created, err := ic.cwl.CreateDelivery(ic.ctx, &cloudwatchlogs.CreateDeliveryInput{
		DeliverySourceName:     aws.String(srcName),
		DeliveryDestinationArn: dest.DeliveryDestination.Arn,
		RecordFields:           []string{"time", "message"},
	})
	if err != nil {
		return r.RunTest(integSvc, "CWLogs_VendedDelivery_S3", func() error {
			return fmt.Errorf("create delivery: %w", err)
		})
	}
	defer ic.cwl.DeleteDelivery(ic.ctx, &cloudwatchlogs.DeleteDeliveryInput{Id: created.Delivery.Id})

	prefix := fmt.Sprintf("AWSLogs/%s/%s/", ic.accountID, srcName)
	return r.pollVerify("CWLogs_VendedDelivery_S3", defaultPollTimeout, func() error {
		listed, err := ic.s3.ListObjectsV2(ic.ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket), Prefix: aws.String(prefix),
		})
		if err != nil {
			return fmt.Errorf("list delivered objects: %w", err)
		}
		if len(listed.Contents) == 0 {
			return fmt.Errorf("no delivered object under %s yet", prefix)
		}
		obj, err := ic.s3.GetObject(ic.ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket), Key: listed.Contents[0].Key,
		})
		if err != nil {
			return fmt.Errorf("get delivered object: %w", err)
		}
		defer obj.Body.Close()
		zr, err := gzip.NewReader(obj.Body)
		if err != nil {
			return fmt.Errorf("delivered object is not gzip: %w", err)
		}
		defer zr.Close()
		body, err := io.ReadAll(zr)
		if err != nil {
			return fmt.Errorf("read delivered object: %w", err)
		}
		lines := strings.Split(strings.TrimRight(string(body), "\n"), "\n")
		if len(lines) < 2 {
			return fmt.Errorf("delivered lines: %q", string(body))
		}
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
			return fmt.Errorf("delivered line is not JSON: %w", err)
		}
		if record["message"] != "vended-first" {
			return fmt.Errorf("first delivered record: %v", record)
		}
		if _, ok := record["time"]; !ok {
			return fmt.Errorf("delivered record missing time: %v", record)
		}
		return nil
	})
}

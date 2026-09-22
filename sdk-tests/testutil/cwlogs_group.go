package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (tc *cwlogsTestCtx) basicTests() []TestResult {
	var results []TestResult

	// Every test builds its own group and stream: the family shares no
	// fixture, so one test's failure cannot cascade into its siblings. The
	// newLogGroupFixture / newGroupStreamFixture constructors state that
	// scaffold once; the tests below still each own an isolated group.

	results = append(results, tc.runner.RunTest("logs", "CreateLogGroup_VerifyFields", func() error {
		logGroupName := tc.uniquePrefix("TestLogGroup")
		resp, err := tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		defer tc.deleteLogGroup(logGroupName)
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		descResp, err := tc.collectAllLogGroups(logGroupName)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp) != 1 {
			return fmt.Errorf("expected 1 log group, got %d", len(descResp))
		}
		lg := descResp[0]
		if lg.LogGroupName == nil || *lg.LogGroupName != logGroupName {
			return fmt.Errorf("logGroupName mismatch: got %q", aws.ToString(lg.LogGroupName))
		}
		if lg.Arn == nil || !strings.Contains(*lg.Arn, logGroupName) {
			return fmt.Errorf("arn missing or does not contain group name: %q", aws.ToString(lg.Arn))
		}
		if lg.CreationTime == nil || *lg.CreationTime == 0 {
			return fmt.Errorf("creationTime is zero or nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeLogGroups_ContainsCreated", func() error {
		logGroupName, cleanupGroup, err := tc.newLogGroupFixture("TestLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		resp, err := tc.collectAllLogGroups(logGroupName)
		if err != nil {
			return err
		}
		if len(resp) != 1 {
			return fmt.Errorf("expected 1 log group, got %d", len(resp))
		}
		lg := resp[0]
		if lg.LogGroupName == nil || *lg.LogGroupName != logGroupName {
			return fmt.Errorf("logGroupName mismatch: got %q, want %q", aws.ToString(lg.LogGroupName), logGroupName)
		}
		if lg.Arn == nil || *lg.Arn == "" {
			return fmt.Errorf("arn is nil or empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "CreateLogStream_VerifyFields", func() error {
		logGroupName, cleanupGroup, err := tc.newLogGroupFixture("TestLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		logStreamName := tc.uniquePrefix("TestLogStream")
		if _, err := tc.client.CreateLogStream(tc.ctx, &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		}); err != nil {
			return err
		}

		descResp, err := tc.collectAllLogStreams(logGroupName)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp) != 1 {
			return fmt.Errorf("expected 1 log stream, got %d", len(descResp))
		}
		ls := descResp[0]
		if ls.LogStreamName == nil || *ls.LogStreamName != logStreamName {
			return fmt.Errorf("logStreamName mismatch: got %q", aws.ToString(ls.LogStreamName))
		}
		if ls.Arn == nil || !strings.Contains(*ls.Arn, logStreamName) {
			return fmt.Errorf("arn missing or does not contain stream name: %q", aws.ToString(ls.Arn))
		}
		if ls.CreationTime == nil || *ls.CreationTime == 0 {
			return fmt.Errorf("creationTime is zero or nil")
		}
		return nil
	}))

	// The event-bearing members of the DescribeLogStreams entry: after a
	// known three-event ingestion the entry reports the ingestion bounds,
	// the zero storedBytes the LogStream shape documents for streams, and
	// the upload token of the last batch.
	results = append(results, tc.runner.RunTest("logs", "DescribeLogStreams_EntryMembers", func() error {
		logGroupName, cleanupGroup, err := tc.newGroupStreamFixture("EntryMembersGroup", "entry-stream")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		base := time.Now().UnixMilli() - 60000
		for i, msg := range []string{"entry-one", "entry-two", "entry-three"} {
			if err := tc.putLogEvent(logGroupName, "entry-stream", msg, base+int64(i)*1000); err != nil {
				return fmt.Errorf("put event %d: %v", i, err)
			}
		}

		streams, err := tc.collectAllLogStreams(logGroupName)
		if err != nil {
			return err
		}
		if len(streams) != 1 {
			return fmt.Errorf("expected 1 log stream, got %d", len(streams))
		}
		ls := streams[0]
		if aws.ToInt64(ls.FirstEventTimestamp) != base {
			return fmt.Errorf("firstEventTimestamp = %d, want %d", aws.ToInt64(ls.FirstEventTimestamp), base)
		}
		if aws.ToInt64(ls.LastEventTimestamp) != base+2000 {
			return fmt.Errorf("lastEventTimestamp = %d, want %d", aws.ToInt64(ls.LastEventTimestamp), base+2000)
		}
		if ls.LastIngestionTime == nil || *ls.LastIngestionTime < base {
			return fmt.Errorf("lastIngestionTime = %v, want an ingestion-time stamp", ls.LastIngestionTime)
		}
		// "As of June 17, 2019, this parameter is no longer supported for
		// log streams, and is always reported as zero" (LogStream shape,
		// storedBytes member documentation).
		if ls.StoredBytes == nil || *ls.StoredBytes != 0 {
			return fmt.Errorf("storedBytes = %v, want the documented zero for streams", ls.StoredBytes)
		}
		if aws.ToString(ls.UploadSequenceToken) == "" {
			return fmt.Errorf("uploadSequenceToken missing after ingestion")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_VerifyNextToken", func() error {
		logGroupName, cleanupGroup, err := tc.newLogGroupFixture("TestLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		logStreamName := tc.uniquePrefix("TestLogStream")
		if err := tc.createLogStream(logGroupName, logStreamName); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}

		put := func(msg string) (*cloudwatchlogs.PutLogEventsOutput, error) {
			return tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
				LogGroupName:  aws.String(logGroupName),
				LogStreamName: aws.String(logStreamName),
				LogEvents: []types.InputLogEvent{
					{
						Message:   aws.String(msg),
						Timestamp: aws.Int64(time.Now().UnixMilli()),
					},
				},
			})
		}
		first, err := put("Test log message")
		if err != nil {
			return err
		}
		if first == nil {
			return fmt.Errorf("response is nil")
		}
		// A fully valid batch reports no rejections and hands back the
		// token the next upload addresses.
		if first.RejectedLogEventsInfo != nil {
			return fmt.Errorf("rejectedLogEventsInfo on a valid batch: %+v", first.RejectedLogEventsInfo)
		}
		if aws.ToString(first.NextSequenceToken) == "" {
			return fmt.Errorf("nextSequenceToken missing")
		}
		second, err := put("Test log message two")
		if err != nil {
			return err
		}
		if aws.ToString(second.NextSequenceToken) == "" || aws.ToString(second.NextSequenceToken) == aws.ToString(first.NextSequenceToken) {
			return fmt.Errorf("nextSequenceToken must advance per batch: %q then %q",
				aws.ToString(first.NextSequenceToken), aws.ToString(second.NextSequenceToken))
		}
		if second.RejectedLogEventsInfo != nil {
			return fmt.Errorf("rejectedLogEventsInfo on the second valid batch: %+v", second.RejectedLogEventsInfo)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "FilterLogEvents_VerifyMessage", func() error {
		logGroupName := tc.uniquePrefix("TestLogGroup")
		logStreamName := tc.uniquePrefix("TestLogStream")
		if err := tc.createLogGroup(logGroupName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer tc.deleteLogGroup(logGroupName)
		if err := tc.createLogStream(logGroupName, logStreamName); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		if err := tc.putLogEvent(logGroupName, logStreamName, "Test log message", time.Now().UnixMilli()); err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := tc.client.FilterLogEvents(tc.ctx, &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}
		if len(resp.Events) == 0 {
			return fmt.Errorf("no events returned")
		}
		if resp.Events[0].Message == nil || *resp.Events[0].Message != "Test log message" {
			return fmt.Errorf("message mismatch: got %q, want %q", aws.ToString(resp.Events[0].Message), "Test log message")
		}
		if resp.Events[0].LogStreamName == nil || *resp.Events[0].LogStreamName != logStreamName {
			return fmt.Errorf("logStreamName mismatch: got %q", aws.ToString(resp.Events[0].LogStreamName))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "PutRetentionPolicy_VerifyRetention", func() error {
		logGroupName, cleanupGroup, err := tc.newLogGroupFixture("TestLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		_, err = tc.client.PutRetentionPolicy(tc.ctx, &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(logGroupName),
			RetentionInDays: aws.Int32(7),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.collectAllLogGroups(logGroupName)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp) != 1 {
			return fmt.Errorf("log group not found uniquely: %d matches", len(descResp))
		}
		if descResp[0].RetentionInDays == nil || *descResp[0].RetentionInDays != 7 {
			return fmt.Errorf("retentionInDays mismatch: got %d, want 7", aws.ToInt32(descResp[0].RetentionInDays))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogStream_VerifyGone", func() error {
		logGroupName, cleanupGroup, err := tc.newLogGroupFixture("TestLogGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		logStreamName := tc.uniquePrefix("TestLogStream")
		if err := tc.createLogStream(logGroupName, logStreamName); err != nil {
			return fmt.Errorf("create stream: %v", err)
		}
		_, err = tc.client.DeleteLogStream(tc.ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.collectAllLogStreams(logGroupName)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		for _, ls := range descResp {
			if ls.LogStreamName != nil && *ls.LogStreamName == logStreamName {
				return fmt.Errorf("log stream still exists after delete")
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogGroup_VerifyGone", func() error {
		logGroupName := tc.uniquePrefix("TestLogGroup")
		if err := tc.createLogGroup(logGroupName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		_, err := tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			return err
		}

		descResp, err := tc.collectAllLogGroups(logGroupName)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		for _, lg := range descResp {
			if lg.LogGroupName != nil && *lg.LogGroupName == logGroupName {
				return fmt.Errorf("log group still exists after delete")
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "CreateLogGroup_Duplicate", func() error {
		dupGroupName, cleanupGroup, err := tc.newLogGroupFixture("DupLogGroup")
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer cleanupGroup()

		_, err = tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(dupGroupName),
		})
		if err := AssertErrorContains(err, "ResourceAlreadyExistsException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteLogGroup_NonExistent", func() error {
		_, err := tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String("nonexistent-log-group-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// The model documents the identifier members of these operations as
	// name-or-ARN; each pin below drives its operation through the ARN the
	// object plane reports (trailing ":*" included).

	results = append(results, tc.runner.RunTest("logs", "GetLogGroupFields_ByIdentifierARN", func() error {
		fieldsName, cleanupGroup, err := tc.newGroupStreamFixture("FieldsGroup", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		if err := tc.putLogEvent(fieldsName, "s1", `{"user_id":"u7","code":200}`, time.Now().UnixMilli()); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		fieldsArn, err := tc.findLogGroupARN(fieldsName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		resp, err := tc.client.GetLogGroupFields(tc.ctx, &cloudwatchlogs.GetLogGroupFieldsInput{
			LogGroupIdentifier: fieldsArn,
		})
		if err != nil {
			return fmt.Errorf("get fields by identifier ARN: %v", err)
		}
		var foundUserID bool
		for _, f := range resp.LogGroupFields {
			if aws.ToString(f.Name) == "user_id" {
				foundUserID = true
			}
		}
		if !foundUserID {
			return fmt.Errorf("identifier-ARN field scan missed the user_id field (%d fields)", len(resp.LogGroupFields))
		}

		// The time member is epoch seconds centring a ±8-minute window:
		// centred on now the fresh event is inside, centred an hour ago
		// it is outside.
		centred, err := tc.client.GetLogGroupFields(tc.ctx, &cloudwatchlogs.GetLogGroupFieldsInput{
			LogGroupIdentifier: fieldsArn,
			Time:               aws.Int64(time.Now().Unix()),
		})
		if err != nil {
			return fmt.Errorf("get fields with time: %v", err)
		}
		var centredFound bool
		for _, f := range centred.LogGroupFields {
			if aws.ToString(f.Name) == "user_id" {
				centredFound = true
			}
		}
		if !centredFound {
			return fmt.Errorf("seconds-centred ±8-minute window missed the fresh field (%d fields)", len(centred.LogGroupFields))
		}
		stale, err := tc.client.GetLogGroupFields(tc.ctx, &cloudwatchlogs.GetLogGroupFieldsInput{
			LogGroupIdentifier: fieldsArn,
			Time:               aws.Int64(time.Now().Add(-time.Hour).Unix()),
		})
		if err != nil {
			return fmt.Errorf("get fields with stale time: %v", err)
		}
		for _, f := range stale.LogGroupFields {
			if aws.ToString(f.Name) == "user_id" {
				return fmt.Errorf("±8-minute window around an hour ago included the fresh field")
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "KmsKey_AssociateDisassociateByIdentifierARN", func() error {
		kmsName, cleanupGroup, err := tc.newLogGroupFixture("KmsGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		// A real symmetric key: associating a key that does not exist or is
		// not usable rejects with InvalidParameterException, so the test key
		// must come from the platform's own KMS.
		keyResp, err := tc.kmsClient.CreateKey(tc.ctx, &kms.CreateKeyInput{
			Description: aws.String("cwlogs associate test key"),
		})
		if err != nil {
			return fmt.Errorf("create key: %v", err)
		}
		keyID := aws.ToString(keyResp.KeyMetadata.KeyId)
		defer func() {
			_, _ = tc.kmsClient.ScheduleKeyDeletion(tc.ctx, &kms.ScheduleKeyDeletionInput{
				KeyId:               aws.String(keyID),
				PendingWindowInDays: aws.Int32(7),
			})
		}()
		keyArn := aws.ToString(keyResp.KeyMetadata.Arn)

		kmsArn, err := tc.findLogGroupARN(kmsName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		if _, err := tc.client.AssociateKmsKey(tc.ctx, &cloudwatchlogs.AssociateKmsKeyInput{
			ResourceIdentifier: kmsArn,
			KmsKeyId:           aws.String(keyArn),
		}); err != nil {
			return fmt.Errorf("associate by identifier ARN: %v", err)
		}
		groups, err := tc.collectAllLogGroups(kmsName)
		if err != nil {
			return fmt.Errorf("describe after associate: %v", err)
		}
		if len(groups) != 1 {
			return fmt.Errorf("expected 1 log group, got %d", len(groups))
		}
		if aws.ToString(groups[0].KmsKeyId) != keyArn {
			return fmt.Errorf("kmsKeyId after associate = %q, want %q", aws.ToString(groups[0].KmsKeyId), keyArn)
		}

		if _, err := tc.client.DisassociateKmsKey(tc.ctx, &cloudwatchlogs.DisassociateKmsKeyInput{
			ResourceIdentifier: kmsArn,
		}); err != nil {
			return fmt.Errorf("disassociate by identifier ARN: %v", err)
		}
		groups, err = tc.collectAllLogGroups(kmsName)
		if err != nil {
			return fmt.Errorf("describe after disassociate: %v", err)
		}
		if aws.ToString(groups[0].KmsKeyId) != "" {
			return fmt.Errorf("kmsKeyId survived disassociation: %q", aws.ToString(groups[0].KmsKeyId))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeletionProtection_ByIdentifierARN_BlocksDelete", func() error {
		dpName, cleanupGroup, err := tc.newLogGroupFixture("DelProtGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		dpArn, err := tc.findLogGroupARN(dpName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		if _, err := tc.client.PutLogGroupDeletionProtection(tc.ctx, &cloudwatchlogs.PutLogGroupDeletionProtectionInput{
			LogGroupIdentifier:        dpArn,
			DeletionProtectionEnabled: aws.Bool(true),
		}); err != nil {
			return fmt.Errorf("put deletion protection by identifier ARN: %v", err)
		}
		if _, err := tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dpName),
		}); err == nil {
			return fmt.Errorf("delete succeeded while deletion protection was enabled through the ARN")
		}
		if _, err := tc.client.PutLogGroupDeletionProtection(tc.ctx, &cloudwatchlogs.PutLogGroupDeletionProtectionInput{
			LogGroupIdentifier:        dpArn,
			DeletionProtectionEnabled: aws.Bool(false),
		}); err != nil {
			return fmt.Errorf("clear deletion protection by identifier ARN: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DataProtectionPolicy_ByIdentifierARN_Roundtrip", func() error {
		dppName, cleanupGroup, err := tc.newLogGroupFixture("DppGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		dppArn, err := tc.findLogGroupARN(dppName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		// The two-block structure the policyDocument member documentation
		// mandates: the Audit block with its FindingsDestination and the
		// Deidentify block with its empty MaskConfig, over matching
		// DataIdentifer arrays.
		policyDocument := `{"Name":"data-protection-policy","Version":"2021-06","Statement":[` +
			`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Audit":{"FindingsDestination":{"CloudWatchLogs":{"LogGroup":"dpp-findings"}}}}},` +
			`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`
		putResp, err := tc.client.PutDataProtectionPolicy(tc.ctx, &cloudwatchlogs.PutDataProtectionPolicyInput{
			LogGroupIdentifier: dppArn,
			PolicyDocument:     aws.String(policyDocument),
		})
		if err != nil {
			return fmt.Errorf("put policy by identifier ARN: %v", err)
		}
		if aws.ToString(putResp.LogGroupIdentifier) != aws.ToString(dppArn) {
			return fmt.Errorf("put echo = %q, want the request identifier %q", aws.ToString(putResp.LogGroupIdentifier), aws.ToString(dppArn))
		}
		// A malformed structure rejects with the operation's declared
		// parameter error.
		if _, err := tc.client.PutDataProtectionPolicy(tc.ctx, &cloudwatchlogs.PutDataProtectionPolicyInput{
			LogGroupIdentifier: dppArn,
			PolicyDocument:     aws.String(`{"Name":"malformed","Configuration":[]}`),
		}); err == nil {
			return fmt.Errorf("a document without the two mandated blocks must reject")
		}

		getResp, err := tc.client.GetDataProtectionPolicy(tc.ctx, &cloudwatchlogs.GetDataProtectionPolicyInput{
			LogGroupIdentifier: dppArn,
		})
		if err != nil {
			return fmt.Errorf("get policy by identifier ARN: %v", err)
		}
		if aws.ToString(getResp.PolicyDocument) != policyDocument {
			return fmt.Errorf("policy document mismatch: got %q", aws.ToString(getResp.PolicyDocument))
		}
		if aws.ToString(getResp.LogGroupIdentifier) != aws.ToString(dppArn) {
			return fmt.Errorf("get echo = %q, want the request identifier %q", aws.ToString(getResp.LogGroupIdentifier), aws.ToString(dppArn))
		}

		if _, err := tc.client.DeleteDataProtectionPolicy(tc.ctx, &cloudwatchlogs.DeleteDataProtectionPolicyInput{
			LogGroupIdentifier: dppArn,
		}); err != nil {
			return fmt.Errorf("delete policy by identifier ARN: %v", err)
		}
		if _, err := tc.client.GetDataProtectionPolicy(tc.ctx, &cloudwatchlogs.GetDataProtectionPolicyInput{
			LogGroupIdentifier: dppArn,
		}); err == nil {
			return fmt.Errorf("policy still readable after ARN-addressed delete")
		}
		return nil
	}))

	return results
}

package testutil

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
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"vorpalstacks-sdk-tests/config"
)

// The delivery pins exercise the trail delivery machinery end to end: a
// logging trail's events reach the S3 bucket as gzip {"Records":[...]}
// log files under the documented naming, a validation-enabled trail
// additionally delivers a signed digest file whose signature verifies
// against the fingerprint-matched ListPublicKeys key, and GetTrailStatus
// reports the recorded delivery outcomes.

// createDeliveryBucket creates the trail's S3 bucket with the CloudTrail
// delivery policy attached.
func (tc *cloudTrailTestContext) createDeliveryBucket(name string) error {
	if _, err := tc.s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(name)}); err != nil {
		return err
	}
	return putCloudTrailBucketPolicy(tc.ctx, tc.s3Client, tc.accountID, name)
}

func (tc *cloudTrailTestContext) deleteDeliveryBucket(name string) {
	objs, err := tc.s3Client.ListObjectsV2(tc.ctx, &s3.ListObjectsV2Input{Bucket: aws.String(name)})
	if err != nil {
		return
	}
	for _, o := range objs.Contents {
		_, _ = tc.s3Client.DeleteObject(tc.ctx, &s3.DeleteObjectInput{Bucket: aws.String(name), Key: o.Key})
	}
	_, _ = tc.s3Client.DeleteBucket(tc.ctx, &s3.DeleteBucketInput{Bucket: aws.String(name)})
}

// waitForDeliveryObject polls the bucket until an object under the prefix
// appears (the TEST_MODE delivery tick is one second) or the bound elapses.
func (tc *cloudTrailTestContext) waitForDeliveryObject(bucket, prefix string, bound time.Duration) (string, error) {
	deadline := time.Now().Add(bound)
	for time.Now().Before(deadline) {
		resp, err := tc.s3Client.ListObjectsV2(tc.ctx, &s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: aws.String(prefix)})
		if err != nil {
			return "", fmt.Errorf("list delivery objects: %v", err)
		}
		for _, o := range resp.Contents {
			if strings.HasSuffix(aws.ToString(o.Key), ".json.gz") {
				return aws.ToString(o.Key), nil
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	return "", fmt.Errorf("no object delivered under %s/%s within %s", bucket, prefix, bound)
}

func (r *TestRunner) runCloudTrailDeliveryTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "TrailDelivery_LogFile", func() error {
		bucket := tc.uniqueName("ct-deliver-bucket")
		if err := tc.createDeliveryBucket(bucket); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer tc.deleteDeliveryBucket(bucket)

		name := tc.uniqueName("ct-deliver")
		defer tc.deleteTrail(name)
		if _, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String(bucket),
		}); err != nil {
			return err
		}
		if _, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("start logging: %v", err)
		}

		// The trail's own API traffic is what gets delivered.
		for i := 0; i < 3; i++ {
			if _, err := tc.client.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
				return fmt.Errorf("generate traffic: %v", err)
			}
		}

		prefix := fmt.Sprintf("AWSLogs/%s/CloudTrail/%s/", tc.accountID, tc.region)
		key, err := tc.waitForDeliveryObject(bucket, prefix, 10*time.Second)
		if err != nil {
			return err
		}

		nameRe := regexp.MustCompile(fmt.Sprintf(`^AWSLogs/%s/CloudTrail/%s/\d{4}/\d{2}/\d{2}/%s_CloudTrail_%s_\d{8}T\d{4}Z_[0-9A-Za-z]{16}\.json\.gz$`,
			regexp.QuoteMeta(tc.accountID), regexp.QuoteMeta(tc.region), regexp.QuoteMeta(tc.accountID), regexp.QuoteMeta(tc.region)))
		if !nameRe.MatchString(key) {
			return fmt.Errorf("delivered key %q does not follow the documented naming", key)
		}

		raw, err := tc.getObjectBytes(bucket, key)
		if err != nil {
			return err
		}
		zr, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return fmt.Errorf("delivered object is not gzip: %v", err)
		}
		body, err := io.ReadAll(zr)
		if err != nil {
			return err
		}
		var parsed struct {
			Records []map[string]interface{} `json:"Records"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			return fmt.Errorf("log file is not {\"Records\":[...]} JSON: %v", err)
		}
		sawListTrails := false
		for _, rec := range parsed.Records {
			if rec["eventName"] == "ListTrails" && rec["eventCategory"] == "Management" {
				sawListTrails = true
			}
		}
		if !sawListTrails {
			return fmt.Errorf("log file carries %d records but none is this trail's own ListTrails event", len(parsed.Records))
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail status: %v", err)
		}
		if !aws.ToBool(status.IsLogging) {
			return fmt.Errorf("IsLogging = %v, want true", status.IsLogging)
		}
		if status.LatestDeliveryTime == nil {
			return fmt.Errorf("LatestDeliveryTime not reported after a real delivery")
		}
		// The model types LatestDeliveryAttemptSucceeded as a string
		// ("true"/"false"), not a boolean.
		if aws.ToString(status.LatestDeliveryAttemptSucceeded) != "true" {
			return fmt.Errorf("LatestDeliveryAttemptSucceeded = %v, want \"true\"", status.LatestDeliveryAttemptSucceeded)
		}
		if status.LatestDeliveryError != nil && *status.LatestDeliveryError != "" {
			return fmt.Errorf("LatestDeliveryError = %q after a successful delivery", *status.LatestDeliveryError)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "TrailDelivery_Digest", func() error {
		bucket := tc.uniqueName("ct-digest-bucket")
		if err := tc.createDeliveryBucket(bucket); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer tc.deleteDeliveryBucket(bucket)

		name := tc.uniqueName("ct-digest")
		defer tc.deleteTrail(name)
		if _, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                    aws.String(name),
			S3BucketName:            aws.String(bucket),
			EnableLogFileValidation: aws.Bool(true),
		}); err != nil {
			return err
		}
		if _, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("start logging: %v", err)
		}
		for i := 0; i < 3; i++ {
			if _, err := tc.client.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
				return fmt.Errorf("generate traffic: %v", err)
			}
		}

		digestPrefix := fmt.Sprintf("AWSLogs/%s/CloudTrail-Digest/%s/", tc.accountID, tc.region)
		key, err := tc.waitForDeliveryObject(bucket, digestPrefix, 15*time.Second)
		if err != nil {
			return err
		}
		digestRe := regexp.MustCompile(fmt.Sprintf(`^AWSLogs/%s/CloudTrail-Digest/%s/\d{4}/\d{2}/\d{2}/%s_CloudTrail-Digest_%s_%s_%s_\d{8}T\d{6}Z\.json\.gz$`,
			regexp.QuoteMeta(tc.accountID), regexp.QuoteMeta(tc.region), regexp.QuoteMeta(tc.accountID), regexp.QuoteMeta(tc.region), regexp.QuoteMeta(name), regexp.QuoteMeta(tc.region)))
		if !digestRe.MatchString(key) {
			return fmt.Errorf("digest key %q does not follow the documented naming", key)
		}

		raw, err := tc.getObjectBytes(bucket, key)
		if err != nil {
			return err
		}
		zr, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return fmt.Errorf("digest object is not gzip: %v", err)
		}
		body, err := io.ReadAll(zr)
		if err != nil {
			return err
		}
		var digest struct {
			AWSAccountId               string  `json:"awsAccountId"`
			DigestPublicKeyFingerprint string  `json:"digestPublicKeyFingerprint"`
			DigestSignatureAlgorithm   string  `json:"digestSignatureAlgorithm"`
			PreviousDigestS3Object     *string `json:"previousDigestS3Object"`
			LogFiles                   []struct {
				HashValue string `json:"hashValue"`
			} `json:"logFiles"`
		}
		if err := json.Unmarshal(body, &digest); err != nil {
			return fmt.Errorf("digest is not JSON: %v", err)
		}
		if digest.AWSAccountId != tc.accountID || digest.DigestSignatureAlgorithm != "SHA256withRSA" {
			return fmt.Errorf("digest header fields wrong: %+v", digest)
		}
		if digest.PreviousDigestS3Object != nil {
			return fmt.Errorf("first digest must be a starting digest, got previousDigestS3Object %q", *digest.PreviousDigestS3Object)
		}
		if len(digest.LogFiles) == 0 {
			return fmt.Errorf("digest references no log files though the trail delivered events")
		}

		keys, err := tc.client.ListPublicKeys(tc.ctx, &cloudtrail.ListPublicKeysInput{})
		if err != nil {
			return fmt.Errorf("list public keys: %v", err)
		}
		var rsaPub *rsa.PublicKey
		for _, pk := range keys.PublicKeyList {
			if aws.ToString(pk.Fingerprint) == digest.DigestPublicKeyFingerprint {
				rsaPub, err = x509.ParsePKCS1PublicKey(pk.Value)
				if err != nil {
					return fmt.Errorf("fingerprint-matched key is not PKCS#1 DER: %v", err)
				}
			}
		}
		if rsaPub == nil {
			return fmt.Errorf("no ListPublicKeys key matches the digest fingerprint %q", digest.DigestPublicKeyFingerprint)
		}

		head, err := tc.s3Client.HeadObject(tc.ctx, &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
		if err != nil {
			return fmt.Errorf("head digest object: %v", err)
		}
		// The S3 response surfaces object metadata under its bare keys
		// (the x-amz-meta- prefix is a wire-header concern).
		sigHex := head.Metadata["signature"]
		if sigHex == "" {
			return fmt.Errorf("digest object carries no signature metadata: %v", head.Metadata)
		}
		if head.Metadata["signature-algorithm"] != "SHA256withRSA" {
			return fmt.Errorf("signature-algorithm metadata = %q", head.Metadata["signature-algorithm"])
		}
		sig, err := hex.DecodeString(sigHex)
		if err != nil {
			return fmt.Errorf("signature metadata is not hexadecimal: %v", err)
		}
		digestSum := sha256.Sum256(body)
		if err := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digestSum[:], sig); err != nil {
			return fmt.Errorf("digest signature does not verify: %v", err)
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail status: %v", err)
		}
		if status.LatestDigestDeliveryTime == nil {
			return fmt.Errorf("LatestDigestDeliveryTime not reported after a digest delivery")
		}
		return nil
	}))

	// CloudWatch Logs delivery: the trail's events reach the configured log
	// group on the documented <account>_CloudTrail_<region> stream, and
	// GetTrailStatus reports the outcome.
	results = append(results, r.RunTest("cloudtrail", "TrailDelivery_CWLogs", func() error {
		bucket := tc.uniqueName("ct-cwl-bucket")
		if err := tc.createDeliveryBucket(bucket); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer tc.deleteDeliveryBucket(bucket)

		logGroup := tc.uniqueName("ct-cwl-group")
		if _, err := tc.logsClient.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{LogGroupName: aws.String(logGroup)}); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.logsClient.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{LogGroupName: aws.String(logGroup)})

		name := tc.uniqueName("ct-cwl")
		defer tc.deleteTrail(name)
		logGroupARN := fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", tc.region, tc.accountID, logGroup)
		if _, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:                      aws.String(name),
			S3BucketName:              aws.String(bucket),
			CloudWatchLogsLogGroupArn: aws.String(logGroupARN),
		}); err != nil {
			return err
		}
		if _, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("start logging: %v", err)
		}
		for i := 0; i < 3; i++ {
			if _, err := tc.client.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
				return fmt.Errorf("generate traffic: %v", err)
			}
		}

		// The delivery worker creates the stream named after the account
		// and the trail's region. The poll scans the whole stream from
		// the head following the forward token: under full-regression
		// load the first flush carries a whole flood window at once, so
		// the fixed newest page a bare poll reads can scroll past the
		// ListTrails records within seconds of their delivery (the
		// production cadence is the documented five minutes; this remains
		// a bounded TEST_MODE wait).
		stream := fmt.Sprintf("%s_CloudTrail_%s", tc.accountID, tc.region)
		deadline := time.Now().Add(45 * time.Second)
		sawListTrails := false
		for time.Now().Before(deadline) && !sawListTrails {
			var token *string
			for !sawListTrails {
				events, err := tc.logsClient.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
					LogGroupName: aws.String(logGroup), LogStreamName: aws.String(stream),
					Limit: aws.Int32(200), StartFromHead: aws.Bool(true), NextToken: token,
				})
				if err != nil {
					break
				}
				for _, ev := range events.Events {
					var record map[string]interface{}
					if json.Unmarshal([]byte(aws.ToString(ev.Message)), &record) == nil && record["eventName"] == "ListTrails" {
						sawListTrails = true
					}
				}
				forward := aws.ToString(events.NextForwardToken)
				if sawListTrails || forward == "" || forward == aws.ToString(token) {
					break
				}
				token = events.NextForwardToken
			}
			if !sawListTrails {
				time.Sleep(300 * time.Millisecond)
			}
		}
		if !sawListTrails {
			return fmt.Errorf("no ListTrails record reached log group %s stream %s", logGroup, stream)
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail status: %v", err)
		}
		if status.LatestCloudWatchLogsDeliveryTime == nil {
			return fmt.Errorf("LatestCloudWatchLogsDeliveryTime not reported after a delivery")
		}
		if status.LatestCloudWatchLogsDeliveryError != nil && *status.LatestCloudWatchLogsDeliveryError != "" {
			return fmt.Errorf("LatestCloudWatchLogsDeliveryError = %q after a successful delivery", *status.LatestCloudWatchLogsDeliveryError)
		}
		return nil
	}))

	// SNS notification delivery: a subscribed SQS queue observes the
	// notification JSON that announces the delivered log files.
	results = append(results, r.RunTest("cloudtrail", "TrailDelivery_Notification", func() error {
		bucket := tc.uniqueName("ct-notify-bucket")
		if err := tc.createDeliveryBucket(bucket); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer tc.deleteDeliveryBucket(bucket)

		topic, err := tc.snsClient.CreateTopic(tc.ctx, &sns.CreateTopicInput{Name: aws.String(tc.uniqueName("ct-notify-topic"))})
		if err != nil {
			return fmt.Errorf("create topic: %v", err)
		}
		defer tc.snsClient.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{TopicArn: topic.TopicArn})
		topicPolicy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"AWSCloudTrailSNSPolicy20131101","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"%s"}]}`,
			aws.ToString(topic.TopicArn))
		if _, err := tc.snsClient.SetTopicAttributes(tc.ctx, &sns.SetTopicAttributesInput{
			TopicArn: topic.TopicArn, AttributeName: aws.String("Policy"), AttributeValue: aws.String(topicPolicy),
		}); err != nil {
			return fmt.Errorf("set topic policy: %v", err)
		}

		queueName := tc.uniqueName("ct-notify-queue")
		queue, err := tc.sqsClient.CreateQueue(tc.ctx, &sqs.CreateQueueInput{QueueName: aws.String(queueName)})
		if err != nil {
			return fmt.Errorf("create queue: %v", err)
		}
		queueURL := aws.ToString(queue.QueueUrl)
		defer func() {
			_, _ = tc.sqsClient.DeleteQueue(tc.ctx, &sqs.DeleteQueueInput{QueueUrl: aws.String(queueURL)})
		}()
		if _, err := tc.snsClient.Subscribe(tc.ctx, &sns.SubscribeInput{
			TopicArn: topic.TopicArn,
			Protocol: aws.String("sqs"),
			Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:%s", tc.region, tc.accountID, queueName)),
		}); err != nil {
			return fmt.Errorf("subscribe queue: %v", err)
		}

		name := tc.uniqueName("ct-notify")
		defer tc.deleteTrail(name)
		if _, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String(bucket),
			SnsTopicName: topic.TopicArn,
		}); err != nil {
			return err
		}
		if _, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("start logging: %v", err)
		}
		for i := 0; i < 3; i++ {
			if _, err := tc.client.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
				return fmt.Errorf("generate traffic: %v", err)
			}
		}

		// The notification is the documented JSON object: the bucket and
		// the delivered object keys.
		deadline := time.Now().Add(15 * time.Second)
		notified := false
		for time.Now().Before(deadline) && !notified {
			msgs, err := tc.sqsClient.ReceiveMessage(tc.ctx, &sqs.ReceiveMessageInput{
				QueueUrl: aws.String(queueURL), MaxNumberOfMessages: 10,
			})
			if err == nil {
				for _, m := range msgs.Messages {
					var envelope struct {
						Message string `json:"Message"`
					}
					if json.Unmarshal([]byte(aws.ToString(m.Body)), &envelope) != nil {
						continue
					}
					var message struct {
						S3Bucket    string   `json:"s3Bucket"`
						S3ObjectKey []string `json:"s3ObjectKey"`
					}
					if json.Unmarshal([]byte(envelope.Message), &message) != nil {
						continue
					}
					if message.S3Bucket == bucket && len(message.S3ObjectKey) > 0 && strings.HasSuffix(message.S3ObjectKey[0], ".json.gz") {
						notified = true
					}
				}
			}
			if !notified {
				time.Sleep(300 * time.Millisecond)
			}
		}
		if !notified {
			return fmt.Errorf("no log-file delivery notification reached queue %s", queueName)
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail status: %v", err)
		}
		if status.LatestNotificationTime == nil {
			return fmt.Errorf("LatestNotificationTime not reported after a notification")
		}
		if aws.ToString(status.LatestNotificationAttemptSucceeded) != "true" {
			return fmt.Errorf("LatestNotificationAttemptSucceeded = %v, want \"true\"", status.LatestNotificationAttemptSucceeded)
		}
		return nil
	}))

	// Multi-region delivery: a trail created with IsMultiRegionTrail
	// delivers the events recorded in a second region under that region's
	// own path — the one bucket receives every enabled region's files.
	results = append(results, r.RunTest("cloudtrail", "TrailDelivery_MultiRegion", func() error {
		bucket := tc.uniqueName("ct-multi-bucket")
		if err := tc.createDeliveryBucket(bucket); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		defer tc.deleteDeliveryBucket(bucket)

		// API traffic in the second region is what its window delivers.
		awayRegion := "us-west-2"
		if tc.region == awayRegion {
			awayRegion = "eu-west-1"
		}
		awayCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   awayRegion,
		})
		if err != nil {
			return fmt.Errorf("load second-region config: %v", err)
		}
		awayClient := cloudtrail.NewFromConfig(awayCfg)

		name := tc.uniqueName("ct-multi")
		defer tc.deleteTrail(name)
		if _, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:               aws.String(name),
			S3BucketName:       aws.String(bucket),
			IsMultiRegionTrail: aws.Bool(true),
		}); err != nil {
			return err
		}
		if _, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)}); err != nil {
			return fmt.Errorf("start logging: %v", err)
		}
		for i := 0; i < 3; i++ {
			if _, err := awayClient.ListTrails(tc.ctx, &cloudtrail.ListTrailsInput{}); err != nil {
				return fmt.Errorf("generate second-region traffic: %v", err)
			}
		}

		// The bucket lives in the trail's region; the delivered file's
		// path and name carry the second region.
		prefix := fmt.Sprintf("AWSLogs/%s/CloudTrail/%s/", tc.accountID, awayRegion)
		key, err := tc.waitForDeliveryObject(bucket, prefix, 15*time.Second)
		if err != nil {
			return err
		}
		nameRe := regexp.MustCompile(fmt.Sprintf(`^AWSLogs/%s/CloudTrail/%s/\d{4}/\d{2}/\d{2}/%s_CloudTrail_%s_\d{8}T\d{4}Z_[0-9A-Za-z]{16}\.json\.gz$`,
			regexp.QuoteMeta(tc.accountID), regexp.QuoteMeta(awayRegion), regexp.QuoteMeta(tc.accountID), regexp.QuoteMeta(awayRegion)))
		if !nameRe.MatchString(key) {
			return fmt.Errorf("delivered key %q does not carry the second region's naming", key)
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get trail status: %v", err)
		}
		if status.LatestDeliveryTime == nil {
			return fmt.Errorf("LatestDeliveryTime not reported after a multi-region delivery")
		}
		return nil
	}))

	return results
}

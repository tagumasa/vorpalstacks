package testutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	cloudtraildata "github.com/aws/aws-sdk-go-v2/service/cloudtraildata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go"
	"vorpalstacks-sdk-tests/config"
)

type cloudTrailTestContext struct {
	client     *cloudtrail.Client
	dataClient *cloudtraildata.Client
	s3Client   *s3.Client
	snsClient  *sns.Client
	sqsClient  *sqs.Client
	logsClient *cloudwatchlogs.Client
	ctx        context.Context
	region     string
	accountID  string
}

// uniqueName returns a collision-free resource name for one test run; Unix
// nanoseconds keep concurrent suites from claiming the same name.
func (tc *cloudTrailTestContext) uniqueName(tag string) string {
	return fmt.Sprintf("%s-%d", tag, time.Now().UnixNano())
}

// createTrail creates a trail with only the two required members set; tests
// that need further members build the input inline so their field coverage
// stays visible. The bucket is provisioned with the CloudTrail delivery
// policy first — CreateTrail rejects a bucket without it.
func (tc *cloudTrailTestContext) createTrail(name, bucket string) (*cloudtrail.CreateTrailOutput, error) {
	if err := tc.ensureTrailBucket(bucket); err != nil {
		return nil, err
	}
	return tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
		Name:         aws.String(name),
		S3BucketName: aws.String(bucket),
	})
}

// ensureTrailBucket creates the bucket when missing and attaches the
// AWS-documented CloudTrail delivery policy (s3:GetBucketAcl on the bucket,
// s3:PutObject on the AWSLogs prefix for the CloudTrail service principal).
func (tc *cloudTrailTestContext) ensureTrailBucket(bucket string) error {
	if _, err := tc.s3Client.CreateBucket(tc.ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) || apiErr.ErrorCode() != "BucketAlreadyOwnedByYou" {
			return fmt.Errorf("CreateBucket %s failed: %w", bucket, err)
		}
	}
	return putCloudTrailBucketPolicy(tc.ctx, tc.s3Client, tc.accountID, bucket)
}

// putCloudTrailBucketPolicy attaches the AWS-documented CloudTrail delivery
// policy to an existing bucket; shared by the trail suite and the audit
// integration suite.
func putCloudTrailBucketPolicy(ctx context.Context, client *s3.Client, accountID, bucket string) error {
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[`+
		`{"Sid":"AWSCloudTrailAclCheck20150319","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::%s"},`+
		`{"Sid":"AWSCloudTrailWrite20150319","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::%s/AWSLogs/%s/*"}]}`,
		bucket, bucket, accountID)
	if _, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket), Policy: aws.String(policy),
	}); err != nil {
		return fmt.Errorf("PutBucketPolicy %s failed: %w", bucket, err)
	}
	return nil
}

// deleteTrail deletes a trail by name; intended for deferred cleanup.
func (tc *cloudTrailTestContext) deleteTrail(name string) {
	_, _ = tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})
}

// trailARN builds the ARN of a trail in the test account and region; only
// negative tests need this, since real creates return the ARN themselves.
func (tc *cloudTrailTestContext) trailARN(name string) string {
	return fmt.Sprintf("arn:aws:cloudtrail:%s:%s:trail/%s", tc.region, tc.accountID, name)
}

// createEventDataStore creates an event data store with termination
// protection off so the test can delete it again; retentionDays is nil when
// the test exercises the service default.
func (tc *cloudTrailTestContext) createEventDataStore(tag string, retentionDays *int32) (*cloudtrail.CreateEventDataStoreOutput, error) {
	return tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
		Name:                         aws.String(tc.uniqueName(tag)),
		TerminationProtectionEnabled: aws.Bool(false),
		RetentionPeriod:              retentionDays,
	})
}

// deleteEventDataStore soft-deletes an event data store given its name or
// ARN; deferred callers discard the error, registered cleanup tests return it.
// Termination protection is lifted first — the create default leaves it on,
// and a protected store would survive the deletion and hold the per-region
// event-data-store quota (which counts every lifecycle stage). The call then
// waits for the record to leave the list: the deletion window (a
// sub-second test window in TEST_MODE, seven days in production) keeps
// counting against the quota until the sweep removes the record, and a
// deferred cleanup that returns early would let fast consecutive fixtures
// pile past the quota.
func (tc *cloudTrailTestContext) deleteEventDataStore(id string) error {
	_, err := tc.client.UpdateEventDataStore(tc.ctx, &cloudtrail.UpdateEventDataStoreInput{
		EventDataStore:               aws.String(id),
		TerminationProtectionEnabled: aws.Bool(false),
	})
	if err != nil {
		var notFound *types.EventDataStoreNotFoundException
		if !errors.As(err, &notFound) {
			return err
		}
		return nil
	}
	if _, err = tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
		EventDataStore: aws.String(id),
	}); err != nil {
		return err
	}
	target := tc.edsIDFromARN(id)
	if target == "" {
		target = id
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, listErr := tc.client.ListEventDataStores(tc.ctx, &cloudtrail.ListEventDataStoresInput{})
		if listErr != nil {
			return listErr
		}
		gone := true
		for _, eds := range resp.EventDataStores {
			if edsIDFromListARN(aws.ToString(eds.EventDataStoreArn)) == target {
				gone = false
				break
			}
		}
		if gone {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("event data store %s still listed after deletion window", id)
}

// edsIDFromListARN extracts the event data store ID from a listed ARN.
func edsIDFromListARN(arn string) string {
	if idx := strings.LastIndex(arn, "/"); idx >= 0 {
		return arn[idx+1:]
	}
	return arn
}

// getObjectBytes fetches an object's full content from the test bucket.
func (tc *cloudTrailTestContext) getObjectBytes(bucket, key string) ([]byte, error) {
	resp, err := tc.s3Client.GetObject(tc.ctx, &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	if err != nil {
		return nil, fmt.Errorf("GetObject %s/%s failed: %w", bucket, key, err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// edsIDFromARN returns the ID segment after the final slash of an event data
// store ARN, which every event-data-store operation accepts as its target.
func (tc *cloudTrailTestContext) edsIDFromARN(arn string) string {
	if idx := strings.LastIndex(arn, "/"); idx >= 0 {
		return arn[idx+1:]
	}
	return arn
}

// tagListToMap converts a tag list to a key/value map for membership
// assertions.
func tagListToMap(tags []types.Tag) map[string]string {
	m := make(map[string]string)
	for _, t := range tags {
		if t.Key != nil && t.Value != nil {
			m[*t.Key] = *t.Value
		}
	}
	return m
}

func (r *TestRunner) RunCloudTrailTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudtrail",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cloudtrail.NewFromConfig(cfg)
	ctx := context.Background()
	tc := &cloudTrailTestContext{
		client: client,
		// The cloudtrail-data ingestion service shares the AWS config; its
		// requests are classified by the cloudtrail-data signing name.
		dataClient: cloudtraildata.NewFromConfig(cfg),
		// Path-style addressing, the same as the audit integration's S3
		// client: the platform's S3 plane routes path-style requests.
		s3Client:   s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true }),
		snsClient:  sns.NewFromConfig(cfg),
		sqsClient:  sqs.NewFromConfig(cfg),
		logsClient: cloudwatchlogs.NewFromConfig(cfg),
		ctx:        ctx,
		region:     r.region,
		accountID:  r.accountID,
	}

	results = append(results, r.runCloudTrailTrailTests(tc)...)
	results = append(results, r.runCloudTrailLoggingTests(tc)...)
	results = append(results, r.runCloudTrailSelectorTests(tc)...)
	results = append(results, r.runCloudTrailKeysTests(tc)...)
	results = append(results, r.runCloudTrailDeliveryTests(tc)...)
	results = append(results, r.runCloudTrailPolicyTests(tc)...)
	results = append(results, r.runCloudTrailTagTests(tc)...)
	results = append(results, r.runCloudTrailEventTests(tc)...)
	results = append(results, r.runCloudTrailEdgeTests(tc)...)
	results = append(results, r.runCloudTrailEDSTests(tc)...)
	results = append(results, r.runCloudTrailChannelTests(tc)...)
	results = append(results, r.runCloudTrailQueryTests(tc)...)
	results = append(results, r.runCloudTrailConfigTests(tc)...)
	results = append(results, r.runCloudTrailValidationTests(tc)...)
	results = append(results, r.runCloudTrailImportTests(tc)...)
	results = append(results, r.runCloudTrailQueryAITests(tc)...)
	results = append(results, r.runCloudTrailDataTests(tc)...)

	return results
}

package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"vorpalstacks-sdk-tests/config"
)

type cloudTrailTestContext struct {
	client    *cloudtrail.Client
	ctx       context.Context
	region    string
	accountID string
}

// uniqueName returns a collision-free resource name for one test run; Unix
// nanoseconds keep concurrent suites from claiming the same name.
func (tc *cloudTrailTestContext) uniqueName(tag string) string {
	return fmt.Sprintf("%s-%d", tag, time.Now().UnixNano())
}

// createTrail creates a trail with only the two required members set; tests
// that need further members build the input inline so their field coverage
// stays visible.
func (tc *cloudTrailTestContext) createTrail(name, bucket string) (*cloudtrail.CreateTrailOutput, error) {
	return tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
		Name:         aws.String(name),
		S3BucketName: aws.String(bucket),
	})
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
func (tc *cloudTrailTestContext) deleteEventDataStore(id string) error {
	_, err := tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
		EventDataStore: aws.String(id),
	})
	return err
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
	tc := &cloudTrailTestContext{client: client, ctx: ctx, region: r.region, accountID: r.accountID}

	results = append(results, r.runCloudTrailTrailTests(tc)...)
	results = append(results, r.runCloudTrailLoggingTests(tc)...)
	results = append(results, r.runCloudTrailSelectorTests(tc)...)
	results = append(results, r.runCloudTrailKeysTests(tc)...)
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

	return results
}

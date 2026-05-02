package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"vorpalstacks-sdk-tests/config"
)

type cloudTrailTestContext struct {
	client *cloudtrail.Client
	ctx    context.Context
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
	tc := &cloudTrailTestContext{client: client, ctx: ctx}

	results = append(results, r.runCloudTrailTrailTests(tc)...)
	results = append(results, r.runCloudTrailLoggingTests(tc)...)
	results = append(results, r.runCloudTrailSelectorTests(tc)...)
	results = append(results, r.runCloudTrailKeysTests(tc)...)
	results = append(results, r.runCloudTrailPolicyTests(tc)...)
	results = append(results, r.runCloudTrailTagTests(tc)...)
	results = append(results, r.runCloudTrailEventTests(tc)...)
	results = append(results, r.runCloudTrailEdgeTests(tc)...)

	return results
}

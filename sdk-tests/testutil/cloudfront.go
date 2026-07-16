package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"vorpalstacks-sdk-tests/config"
)

type cfTestContext struct {
	runner *TestRunner
	client *cloudfront.Client
	ctx    context.Context
}

func (tc *cfTestContext) uniquePrefix(tag string) string {
	return fmt.Sprintf("%s-%d", tag, time.Now().UnixNano())
}

func (r *TestRunner) RunCloudFrontTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudfront",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	tc := &cfTestContext{
		runner: r,
		client: cloudfront.NewFromConfig(cfg),
		ctx:    context.Background(),
	}

	results = append(results, cfDistributionTests(tc)...)
	results = append(results, cfTagTests(tc)...)
	results = append(results, cfInvalidationTests(tc)...)
	results = append(results, cfOACTests(tc)...)
	results = append(results, cfCachePolicyTests(tc)...)
	results = append(results, cfOriginRequestPolicyTests(tc)...)
	results = append(results, cfResponseHeadersPolicyTests(tc)...)
	results = append(results, cfEdgeTests(tc)...)

	return results
}

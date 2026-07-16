package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"vorpalstacks-sdk-tests/config"
)

type snsTestContext struct {
	client *sns.Client
	ctx    context.Context
	region string
}

func (r *TestRunner) initSNS() (*snsTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &snsTestContext{
		client: sns.NewFromConfig(cfg),
		ctx:    context.Background(),
		region: r.region,
	}, nil
}

func (tc *snsTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *snsTestContext) createTopic(name string) (string, error) {
	resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
		Name: ptrString(name),
	})
	if err != nil {
		return "", fmt.Errorf("create topic %q: %w", name, err)
	}
	return *resp.TopicArn, nil
}

func (tc *snsTestContext) deleteTopic(arn string) {
	tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
		TopicArn: ptrString(arn),
	})
}

func (r *TestRunner) RunSNSTests() []TestResult {
	tc, err := r.initSNS()
	if err != nil {
		return []TestResult{{
			Service:  "sns",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}

	var results []TestResult

	results = append(results, r.runSNSTopicTests(tc)...)
	results = append(results, r.runSNSPublishTests(tc)...)
	results = append(results, r.runSNSSubscriptionTests(tc)...)
	results = append(results, r.runSNSPlatformTests(tc)...)
	results = append(results, r.runSNSTagTests(tc)...)
	results = append(results, r.runSNSPolicyTests(tc)...)
	results = append(results, r.runSNSEdgeTests(tc)...)

	return results
}

func ptrString(s string) *string { return &s }

func nanoTime() int64 { return time.Now().UnixNano() }

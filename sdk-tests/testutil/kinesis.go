package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunKinesisTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "kinesis",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := kinesis.NewFromConfig(cfg)
	ctx := context.Background()
	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	results = append(results, r.kinesisStreamTests(ctx, client, ts)...)
	results = append(results, r.kinesisRecordTests(ctx, client, ts)...)
	results = append(results, r.kinesisShardTests(ctx, client, ts)...)
	results = append(results, r.kinesisConsumerTests(ctx, client, ts)...)
	results = append(results, r.kinesisConfigTests(ctx, client, ts)...)
	results = append(results, r.kinesisTagTests(ctx, client, ts)...)
	results = append(results, r.kinesisEdgeTests(ctx, client, ts)...)

	return results
}

func kinesisStream(ts, name string) string {
	return fmt.Sprintf("kinesis-%s-%s", name, ts)
}

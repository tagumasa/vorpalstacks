package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
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

// kinesisCreateStream creates a provisioned stream and lets it settle for
// the caller-provided duration before returning; the returned closure
// deletes the stream.
func kinesisCreateStream(ctx context.Context, client *kinesis.Client, streamName string, shardCount int32, settle time.Duration) (func(), error) {
	_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(streamName),
		ShardCount: aws.Int32(shardCount),
	})
	if err != nil {
		return nil, fmt.Errorf("create stream %q failed: %w", streamName, err)
	}
	cleanup := func() {
		client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})
	}
	time.Sleep(settle)
	return cleanup, nil
}

func kinesisTagMap(tags []types.Tag) map[string]string {
	m := make(map[string]string)
	for _, t := range tags {
		m[aws.ToString(t.Key)] = aws.ToString(t.Value)
	}
	return m
}

// kinesisOpenShards keeps shards whose ending sequence number is unset,
// i.e. shards still accepting records.
func kinesisOpenShards(shards []types.Shard) []types.Shard {
	var open []types.Shard
	for _, s := range shards {
		if s.SequenceNumberRange.EndingSequenceNumber == nil || *s.SequenceNumberRange.EndingSequenceNumber == "" {
			open = append(open, s)
		}
	}
	return open
}

func (r *TestRunner) kinesisKMSKeyARN(keyID string) string {
	return fmt.Sprintf("arn:aws:kms:%s:%s:key/%s", r.region, r.accountID, keyID)
}

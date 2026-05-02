package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func (r *TestRunner) kinesisEdgeTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kinesis", "DescribeStream_NonExistent", func() error {
		_, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DeleteStream_NonExistent", func() error {
		_, err := client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "CreateStream_DuplicateName", func() error {
		dupStreamName := fmt.Sprintf("dup-stream-%s", ts)
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(dupStreamName),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(dupStreamName)})

		_, err = client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(dupStreamName),
			ShardCount: aws.Int32(1),
		})
		if err := AssertErrorContains(err, "ResourceInUseException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListShards_NonExistentStream", func() error {
		_, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeLimits", func() error {
		resp, err := client.DescribeLimits(ctx, &kinesis.DescribeLimitsInput{})
		if err != nil {
			return err
		}
		if resp.ShardLimit == nil {
			return fmt.Errorf("ShardLimit is nil")
		}
		if aws.ToInt32(resp.ShardLimit) <= 0 {
			return fmt.Errorf("ShardLimit: expected > 0, got %d", aws.ToInt32(resp.ShardLimit))
		}
		if resp.OpenShardCount == nil {
			return fmt.Errorf("OpenShardCount is nil")
		}
		if aws.ToInt32(resp.OpenShardCount) < 0 {
			return fmt.Errorf("OpenShardCount: expected >= 0, got %d", aws.ToInt32(resp.OpenShardCount))
		}
		return nil
	}))

	return results
}

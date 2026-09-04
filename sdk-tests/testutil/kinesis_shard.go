package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

func (r *TestRunner) kinesisShardTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kinesis", "ListShards", func() error {
		sn := kinesisStream(ts, "shards")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		if len(resp.Shards) != 1 {
			return fmt.Errorf("expected 1 shard, got %d", len(resp.Shards))
		}
		shard := resp.Shards[0]
		if shard.ShardId == nil || *shard.ShardId == "" {
			return fmt.Errorf("ShardId is nil or empty")
		}
		if shard.HashKeyRange == nil {
			return fmt.Errorf("HashKeyRange is nil")
		}
		if shard.SequenceNumberRange == nil {
			return fmt.Errorf("SequenceNumberRange is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListShards_MultiShard", func() error {
		sn := kinesisStream(ts, "multi")
		const shardCount = 3
		cleanup, err := kinesisCreateStream(ctx, client, sn, shardCount, 1*time.Second)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		if len(resp.Shards) != shardCount {
			return fmt.Errorf("expected %d shards, got %d", shardCount, len(resp.Shards))
		}
		shardIDs := make(map[string]bool)
		for _, s := range resp.Shards {
			shardIDs[aws.ToString(s.ShardId)] = true
		}
		if len(shardIDs) != shardCount {
			return fmt.Errorf("expected %d unique shard IDs, got %d", shardCount, len(shardIDs))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListShards_MaxResultsWindow", func() error {
		sn := kinesisStream(ts, "mrw")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 2, 1*time.Second)
		if err != nil {
			return err
		}
		defer cleanup()

		// The MaxResults wire member must bound the page, and the accepted
		// input window is 1-10000.
		one, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(sn),
			MaxResults: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		if len(one.Shards) != 1 {
			return fmt.Errorf("MaxResults 1: got %d shards, want 1", len(one.Shards))
		}
		if one.NextToken == nil {
			return fmt.Errorf("MaxResults 1: NextToken missing")
		}

		for _, maxResults := range []int32{0, 10001} {
			_, err := client.ListShards(ctx, &kinesis.ListShardsInput{
				StreamName: aws.String(sn),
				MaxResults: aws.Int32(maxResults),
			})
			if err == nil {
				return fmt.Errorf("MaxResults %d: expected InvalidArgumentException, got success", maxResults)
			}
			if err := AssertErrorContains(err, "InvalidArgumentException"); err != nil {
				return fmt.Errorf("MaxResults %d: %v", maxResults, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "UpdateShardCount", func() error {
		sn := kinesisStream(ts, "usc")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.UpdateShardCount(ctx, &kinesis.UpdateShardCountInput{
			StreamName:       aws.String(sn),
			TargetShardCount: aws.Int32(2),
			ScalingType:      types.ScalingTypeUniformScaling,
		})
		if err != nil {
			return err
		}
		if resp.CurrentShardCount == nil {
			return fmt.Errorf("CurrentShardCount is nil")
		}
		if aws.ToInt32(resp.TargetShardCount) != 2 {
			return fmt.Errorf("TargetShardCount: got %d, want 2", aws.ToInt32(resp.TargetShardCount))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "MergeShards", func() error {
		sn := kinesisStream(ts, "merge")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 2, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		openShards := kinesisOpenShards(resp.Shards)
		if len(openShards) < 2 {
			return fmt.Errorf("need at least 2 open shards for merge, got %d", len(openShards))
		}
		_, err = client.MergeShards(ctx, &kinesis.MergeShardsInput{
			StreamName:           aws.String(sn),
			ShardToMerge:         openShards[0].ShardId,
			AdjacentShardToMerge: openShards[1].ShardId,
		})
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)

		postResp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("post-merge list shards: %v", err)
		}
		postOpenShards := kinesisOpenShards(postResp.Shards)
		if len(postOpenShards) >= len(openShards) {
			return fmt.Errorf("expected fewer open shards after merge, before=%d after=%d", len(openShards), len(postOpenShards))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "SplitShard", func() error {
		sn := kinesisStream(ts, "split")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		open := kinesisOpenShards(resp.Shards)
		if len(open) == 0 {
			return fmt.Errorf("no open shard found for split")
		}
		_, err = client.SplitShard(ctx, &kinesis.SplitShardInput{
			StreamName:         aws.String(sn),
			ShardToSplit:       open[0].ShardId,
			NewStartingHashKey: aws.String("9223372036854775808"),
		})
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)

		postResp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("post-split list shards: %v", err)
		}
		postOpenShards := kinesisOpenShards(postResp.Shards)
		if len(postOpenShards) < 2 {
			return fmt.Errorf("expected >= 2 open shards after split, got %d", len(postOpenShards))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListShardsWithExclusiveStart", func() error {
		sn := kinesisStream(ts, "lsex")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 1*time.Second)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		if len(resp.Shards) == 0 {
			return fmt.Errorf("expected at least 1 shard, got 0")
		}
		if resp.Shards[0].ShardId == nil {
			return fmt.Errorf("ShardId is nil")
		}
		if resp.NextToken != nil && *resp.NextToken != "" {
			nextResp, err := client.ListShards(ctx, &kinesis.ListShardsInput{
				StreamName:  aws.String(sn),
				NextToken:   resp.NextToken,
				ShardFilter: &types.ShardFilter{Type: types.ShardFilterTypeAtLatest},
			})
			if err != nil {
				return fmt.Errorf("list shards with next token: %v", err)
			}
			if nextResp.Shards == nil {
				return fmt.Errorf("next page shards is nil")
			}
		}
		return nil
	}))

	return results
}

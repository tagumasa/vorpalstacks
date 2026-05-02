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
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)

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
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(shardCount),
		})
		if err != nil {
			return fmt.Errorf("failed to create stream: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(1 * time.Second)

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

	results = append(results, r.RunTest("kinesis", "UpdateShardCount", func() error {
		sn := kinesisStream(ts, "usc")
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)

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
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		var openShards []types.Shard
		for _, s := range resp.Shards {
			if s.SequenceNumberRange.EndingSequenceNumber == nil || *s.SequenceNumberRange.EndingSequenceNumber == "" {
				openShards = append(openShards, s)
			}
		}
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
		var postOpenShards []types.Shard
		for _, s := range postResp.Shards {
			if s.SequenceNumberRange.EndingSequenceNumber == nil || *s.SequenceNumberRange.EndingSequenceNumber == "" {
				postOpenShards = append(postOpenShards, s)
			}
		}
		if len(postOpenShards) >= len(openShards) {
			return fmt.Errorf("expected fewer open shards after merge, before=%d after=%d", len(openShards), len(postOpenShards))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "SplitShard", func() error {
		sn := kinesisStream(ts, "split")
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)

		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}
		var openShard *types.Shard
		for _, s := range resp.Shards {
			if s.SequenceNumberRange.EndingSequenceNumber == nil || *s.SequenceNumberRange.EndingSequenceNumber == "" {
				openShard = &s
				break
			}
		}
		if openShard == nil {
			return fmt.Errorf("no open shard found for split")
		}
		_, err = client.SplitShard(ctx, &kinesis.SplitShardInput{
			StreamName:         aws.String(sn),
			ShardToSplit:       openShard.ShardId,
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
		var postOpenShards []types.Shard
		for _, s := range postResp.Shards {
			if s.SequenceNumberRange.EndingSequenceNumber == nil || *s.SequenceNumberRange.EndingSequenceNumber == "" {
				postOpenShards = append(postOpenShards, s)
			}
		}
		if len(postOpenShards) < 2 {
			return fmt.Errorf("expected >= 2 open shards after split, got %d", len(postOpenShards))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListShardsWithExclusiveStart", func() error {
		sn := kinesisStream(ts, "lsex")
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(1 * time.Second)

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

package testutil

import (
	"context"
	"fmt"
	"math/big"
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

		// The SDKs serialise Timestamp members as epoch-second JSON
		// numbers, so the shard filter and the stream-generation member
		// must both read the numeric wire form end to end. A future
		// moment keeps the open set: the shard started before it and
		// remains open.
		filtered, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(sn),
			ShardFilter: &types.ShardFilter{
				Type:      types.ShardFilterTypeAtTimestamp,
				Timestamp: aws.Time(time.Now().Add(time.Minute)),
			},
		})
		if err != nil {
			return fmt.Errorf("AT_TIMESTAMP shard filter: %v", err)
		}
		if len(filtered.Shards) != 1 {
			return fmt.Errorf("filtered shards: expected 1, got %d", len(filtered.Shards))
		}
		desc, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("describe summary: %v", err)
		}
		if _, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName:              aws.String(sn),
			StreamCreationTimestamp: desc.StreamDescriptionSummary.StreamCreationTimestamp,
		}); err != nil {
			return fmt.Errorf("StreamCreationTimestamp round trip: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeStream_LimitPages", func() error {
		sn := kinesisStream(ts, "dsplim")
		const shardCount = 3
		cleanup, err := kinesisCreateStream(ctx, client, sn, shardCount, 1*time.Second)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(sn),
			Limit:      aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if len(resp.StreamDescription.Shards) != 2 || !aws.ToBool(resp.StreamDescription.HasMoreShards) {
			return fmt.Errorf("page 1: %d shards, HasMoreShards %v", len(resp.StreamDescription.Shards), aws.ToBool(resp.StreamDescription.HasMoreShards))
		}

		resp, err = client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName:            aws.String(sn),
			Limit:                 aws.Int32(2),
			ExclusiveStartShardId: resp.StreamDescription.Shards[1].ShardId,
		})
		if err != nil {
			return err
		}
		if len(resp.StreamDescription.Shards) != 1 || aws.ToBool(resp.StreamDescription.HasMoreShards) {
			return fmt.Errorf("page 2: %d shards, HasMoreShards %v", len(resp.StreamDescription.Shards), aws.ToBool(resp.StreamDescription.HasMoreShards))
		}

		// A Limit above the documented hundred caps to a full page rather
		// than rejecting — the member's range trait accepts the value.
		resp, err = client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(sn),
			Limit:      aws.Int32(150),
		})
		if err != nil {
			return err
		}
		if len(resp.StreamDescription.Shards) != shardCount || aws.ToBool(resp.StreamDescription.HasMoreShards) {
			return fmt.Errorf("limit 150: %d shards, HasMoreShards %v", len(resp.StreamDescription.Shards), aws.ToBool(resp.StreamDescription.HasMoreShards))
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
		// CurrentShardCount reports the count at request time, before the
		// reshaping runs — the documented example answers the pre-update
		// count alongside the target (a three-shard stream scaled to six
		// returns CurrentShardCount 3).
		if aws.ToInt32(resp.CurrentShardCount) != 1 {
			return fmt.Errorf("CurrentShardCount: got %d, want the pre-update count 1", aws.ToInt32(resp.CurrentShardCount))
		}
		if aws.ToInt32(resp.TargetShardCount) != 2 {
			return fmt.Errorf("TargetShardCount: got %d, want 2", aws.ToInt32(resp.TargetShardCount))
		}
		post, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("post-update list shards: %v", err)
		}
		if n := len(kinesisOpenShards(post.Shards)); n != 2 {
			return fmt.Errorf("open shards after 1 to 2: got %d, want 2", n)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "UpdateShardCount_ScaleDownConverges", func() error {
		sn := kinesisStream(ts, "uscsd")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 4, 1*time.Second)
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
		if aws.ToInt32(resp.CurrentShardCount) != 4 {
			return fmt.Errorf("CurrentShardCount: got %d, want the pre-update count 4", aws.ToInt32(resp.CurrentShardCount))
		}
		post, err := client.ListShards(ctx, &kinesis.ListShardsInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("post-scale-down list shards: %v", err)
		}
		open := kinesisOpenShards(post.Shards)
		if len(open) != 2 {
			return fmt.Errorf("open shards after 4 to 2: got %d, want 2 — the scaling must converge, not no-op", len(open))
		}
		// The survivors tile the hash key space contiguously: the lowest
		// starts at 0, the highest ends at the maximum key, and each
		// boundary is adjacent.
		ordered := make([]types.Shard, len(open))
		copy(ordered, open)
		for i := 1; i < len(ordered); i++ {
			for j := i; j > 0 && aws.ToString(ordered[j].HashKeyRange.StartingHashKey) < aws.ToString(ordered[j-1].HashKeyRange.StartingHashKey); j-- {
				ordered[j], ordered[j-1] = ordered[j-1], ordered[j]
			}
		}
		if aws.ToString(ordered[0].HashKeyRange.StartingHashKey) != "0" {
			return fmt.Errorf("key space tiling: lowest open shard starts at %s, want 0", aws.ToString(ordered[0].HashKeyRange.StartingHashKey))
		}
		for i := 1; i < len(ordered); i++ {
			prevEnd := new(big.Int)
			if _, ok := prevEnd.SetString(aws.ToString(ordered[i-1].HashKeyRange.EndingHashKey), 10); !ok {
				return fmt.Errorf("unparseable ending hash key %s", aws.ToString(ordered[i-1].HashKeyRange.EndingHashKey))
			}
			start := new(big.Int)
			if _, ok := start.SetString(aws.ToString(ordered[i].HashKeyRange.StartingHashKey), 10); !ok {
				return fmt.Errorf("unparseable starting hash key %s", aws.ToString(ordered[i].HashKeyRange.StartingHashKey))
			}
			if new(big.Int).Sub(start, prevEnd).Cmp(big.NewInt(1)) != 0 {
				return fmt.Errorf("key space tiling: gap between %s and %s", aws.ToString(ordered[i-1].ShardId), aws.ToString(ordered[i].ShardId))
			}
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

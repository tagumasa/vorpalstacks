package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

func (r *TestRunner) kinesisStreamTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	streamName := kinesisStream(ts, "main")

	results = append(results, r.RunTest("kinesis", "CreateStream", func() error {
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("describe after create: %v", err)
		}
		if aws.ToString(descResp.StreamDescription.StreamName) != streamName {
			return fmt.Errorf("StreamName mismatch: got %q, want %q", aws.ToString(descResp.StreamDescription.StreamName), streamName)
		}
		if descResp.StreamDescription.StreamStatus != types.StreamStatusActive {
			return fmt.Errorf("StreamStatus expected ACTIVE, got %s", descResp.StreamDescription.StreamStatus)
		}
		if descResp.StreamDescription.StreamARN == nil {
			return fmt.Errorf("StreamARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListStreams", func() error {
		var nextToken *string
		found := false
		for {
			resp, err := client.ListStreams(ctx, &kinesis.ListStreamsInput{
				Limit:     aws.Int32(10),
				NextToken: nextToken,
			})
			if err != nil {
				return err
			}
			if resp.StreamNames == nil {
				return fmt.Errorf("StreamNames is nil")
			}
			for _, sn := range resp.StreamNames {
				if sn == streamName {
					found = true
					break
				}
			}
			if found {
				break
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if !found {
			return fmt.Errorf("created stream %q not found in ListStreams response", streamName)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeStream", func() error {
		resp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return err
		}
		if resp.StreamDescription == nil {
			return fmt.Errorf("StreamDescription is nil")
		}
		if aws.ToString(resp.StreamDescription.StreamName) != streamName {
			return fmt.Errorf("StreamName: got %q, want %q", aws.ToString(resp.StreamDescription.StreamName), streamName)
		}
		if resp.StreamDescription.StreamARN == nil {
			return fmt.Errorf("StreamARN is nil")
		}
		if resp.StreamDescription.StreamStatus != types.StreamStatusActive {
			return fmt.Errorf("StreamStatus: expected ACTIVE, got %s", resp.StreamDescription.StreamStatus)
		}
		if len(resp.StreamDescription.Shards) == 0 {
			return fmt.Errorf("expected at least 1 shard, got 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeStreamSummary", func() error {
		resp, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return err
		}
		sd := resp.StreamDescriptionSummary
		if sd == nil {
			return fmt.Errorf("StreamDescriptionSummary is nil")
		}
		if aws.ToString(sd.StreamName) != streamName {
			return fmt.Errorf("StreamName: got %q, want %q", aws.ToString(sd.StreamName), streamName)
		}
		if sd.StreamStatus != types.StreamStatusActive {
			return fmt.Errorf("StreamStatus: expected ACTIVE, got %s", sd.StreamStatus)
		}
		if sd.OpenShardCount == nil || *sd.OpenShardCount < 1 {
			return fmt.Errorf("OpenShardCount: expected >= 1, got %d", aws.ToInt32(sd.OpenShardCount))
		}
		if sd.StreamCreationTimestamp == nil || sd.StreamCreationTimestamp.IsZero() {
			return fmt.Errorf("StreamCreationTimestamp is nil or zero")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "CreateStreamWithTags", func() error {
		sn := kinesisStream(ts, "tags")
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		if err != nil {
			return err
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)

		tagResp, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
			StreamName: aws.String(sn),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		tagMap := make(map[string]string)
		for _, t := range tagResp.Tags {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("tag Environment: got %q, want %q", tagMap["Environment"], "test")
		}
		if tagMap["Owner"] != "test-user" {
			return fmt.Errorf("tag Owner: got %q, want %q", tagMap["Owner"], "test-user")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DeleteStream", func() error {
		delSn := kinesisStream(ts, "del")
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(delSn),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		time.Sleep(500 * time.Millisecond)

		_, err = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(delSn),
		})
		if err != nil {
			return err
		}

		_, err = client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(delSn),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("stream should not exist after delete: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DescribeStream_VerifyTimestamp", func() error {
		sn := kinesisStream(ts, "ts")
		beforeCreate := time.Now()
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)
		resp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(sn),
		})
		if err != nil {
			return err
		}
		ts := resp.StreamDescription.StreamCreationTimestamp
		if ts == nil {
			return fmt.Errorf("StreamCreationTimestamp is nil")
		}
		if ts.IsZero() {
			return fmt.Errorf("StreamCreationTimestamp is zero")
		}
		if ts.Before(beforeCreate.Add(-5*time.Minute)) || ts.After(time.Now().Add(5*time.Minute)) {
			return fmt.Errorf("StreamCreationTimestamp %v is outside expected range", ts)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListStreams_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgStreams []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagStr-%s-%d", pgTs, i)
			_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
				StreamName: aws.String(name),
				ShardCount: aws.Int32(1),
			})
			if err != nil {
				for _, sn := range pgStreams {
					client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
						StreamName:              aws.String(sn),
						EnforceConsumerDeletion: aws.Bool(true),
					})
				}
				return fmt.Errorf("create stream %s: %v", name, err)
			}
			pgStreams = append(pgStreams, name)
		}

		var allStreams []string
		var nextToken *string
		pages := 0
		for {
			resp, err := client.ListStreams(ctx, &kinesis.ListStreamsInput{
				Limit:     aws.Int32(2),
				NextToken: nextToken,
			})
			if err != nil {
				for _, sn := range pgStreams {
					client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
						StreamName:              aws.String(sn),
						EnforceConsumerDeletion: aws.Bool(true),
					})
				}
				return fmt.Errorf("list streams page: %v", err)
			}
			for _, sn := range resp.StreamNames {
				if strings.HasPrefix(sn, "PagStr-"+pgTs) {
					allStreams = append(allStreams, sn)
				}
			}
			pages++
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, sn := range pgStreams {
			client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
				StreamName:              aws.String(sn),
				EnforceConsumerDeletion: aws.Bool(true),
			})
		}
		if len(allStreams) != 5 {
			return fmt.Errorf("expected 5 paginated streams, got %d", len(allStreams))
		}
		if pages < 3 {
			return fmt.Errorf("expected >= 3 pages with limit=2 for 5 streams, got %d", pages)
		}
		return nil
	}))

	return results
}

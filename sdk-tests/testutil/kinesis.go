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

	streamName := fmt.Sprintf("test-stream-%d", time.Now().UnixNano())

	// Test 1: Create Stream
	results = append(results, r.RunTest("kinesis", "CreateStream", func() error {
		resp, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	// Test 2: List Streams
	results = append(results, r.RunTest("kinesis", "ListStreams", func() error {
		resp, err := client.ListStreams(ctx, &kinesis.ListStreamsInput{
			Limit: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.StreamNames == nil {
			return fmt.Errorf("StreamNames is nil")
		}
		return nil
	}))

	// Test 3: Describe Stream
	var streamARN string
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
		if resp.StreamDescription.StreamARN == nil {
			return fmt.Errorf("StreamDescription.StreamARN is nil")
		}
		streamARN = aws.ToString(resp.StreamDescription.StreamARN)
		return nil
	}))

	// Test 4: Describe Stream Summary
	results = append(results, r.RunTest("kinesis", "DescribeStreamSummary", func() error {
		resp, err := client.DescribeStreamSummary(ctx, &kinesis.DescribeStreamSummaryInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return err
		}
		if resp.StreamDescriptionSummary == nil {
			return fmt.Errorf("StreamDescriptionSummary is nil")
		}
		return nil
	}))

	// Test 5: Put Record
	results = append(results, r.RunTest("kinesis", "PutRecord", func() error {
		resp, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),
			Data:         []byte("test-data"),
			PartitionKey: aws.String("partition-key-1"),
		})
		if err != nil {
			return err
		}
		if resp.SequenceNumber == nil {
			return fmt.Errorf("SequenceNumber is nil")
		}
		return nil
	}))

	// Test 6: Put Records
	results = append(results, r.RunTest("kinesis", "PutRecords", func() error {
		resp, err := client.PutRecords(ctx, &kinesis.PutRecordsInput{
			StreamName: aws.String(streamName),
			Records: []types.PutRecordsRequestEntry{
				{
					Data:         []byte("test-data-1"),
					PartitionKey: aws.String("partition-key-1"),
				},
				{
					Data:         []byte("test-data-2"),
					PartitionKey: aws.String("partition-key-2"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Records == nil {
			return fmt.Errorf("Records is nil")
		}
		for i, rec := range resp.Records {
			if rec.SequenceNumber == nil {
				return fmt.Errorf("record[%d].SequenceNumber is nil", i)
			}
		}
		return nil
	}))

	// Test 7: Get Shard Iterator
	var shardID string
	describeResp, _ := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	})
	if len(describeResp.StreamDescription.Shards) > 0 {
		shardID = aws.ToString(describeResp.StreamDescription.Shards[0].ShardId)
	}

	var shardIterator string
	if shardID != "" {
		results = append(results, r.RunTest("kinesis", "GetShardIterator", func() error {
			resp, err := client.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
				StreamName:        aws.String(streamName),
				ShardId:           aws.String(shardID),
				ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
			})
			if err != nil {
				return err
			}
			shardIterator = aws.ToString(resp.ShardIterator)
			return nil
		}))
	}

	// Test 8: List Shards
	results = append(results, r.RunTest("kinesis", "ListShards", func() error {
		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return err
		}
		if resp.Shards == nil {
			return fmt.Errorf("Shards is nil")
		}
		return nil
	}))

	// Test 8b: ListShards_MultiShard - verify all shards returned for a multi-shard stream
	streamNameMulti := fmt.Sprintf("test-stream-multi-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("kinesis", "ListShards_MultiShard", func() error {
		const shardCount = 3
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamNameMulti),
			ShardCount: aws.Int32(shardCount),
		})
		if err != nil {
			return fmt.Errorf("failed to create stream: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(streamNameMulti),
		})
		time.Sleep(1 * time.Second)
		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(streamNameMulti),
		})
		if err != nil {
			return err
		}
		if len(resp.Shards) != shardCount {
			return fmt.Errorf("expected %d shards, got %d", shardCount, len(resp.Shards))
		}
		return nil
	}))

	// Test 9: Get Records (using real iterator)
	if shardIterator != "" {
		results = append(results, r.RunTest("kinesis", "GetRecords", func() error {
			resp, err := client.GetRecords(ctx, &kinesis.GetRecordsInput{
				ShardIterator: aws.String(shardIterator),
			})
			if err != nil {
				return err
			}
			if resp.Records == nil {
				return fmt.Errorf("Records is nil")
			}
			return nil
		}))
	}

	// Test 10: Enable Enhanced Monitoring
	results = append(results, r.RunTest("kinesis", "EnableEnhancedMonitoring", func() error {
		resp, err := client.EnableEnhancedMonitoring(ctx, &kinesis.EnableEnhancedMonitoringInput{
			StreamName: aws.String(streamName),
			ShardLevelMetrics: []types.MetricsName{
				types.MetricsNameIncomingBytes,
				types.MetricsNameOutgoingBytes,
			},
		})
		if err != nil {
			return err
		}
		if resp.CurrentShardLevelMetrics == nil {
			return fmt.Errorf("CurrentShardLevelMetrics is nil")
		}
		return nil
	}))

	// Test 11: Disable Enhanced Monitoring
	results = append(results, r.RunTest("kinesis", "DisableEnhancedMonitoring", func() error {
		resp, err := client.DisableEnhancedMonitoring(ctx, &kinesis.DisableEnhancedMonitoringInput{
			StreamName:        aws.String(streamName),
			ShardLevelMetrics: []types.MetricsName{},
		})
		if err != nil {
			return err
		}
		if resp.CurrentShardLevelMetrics == nil {
			return fmt.Errorf("CurrentShardLevelMetrics is nil")
		}
		return nil
	}))

	// Test 12: Add Tags to Stream
	results = append(results, r.RunTest("kinesis", "AddTagsToStream", func() error {
		_, err := client.AddTagsToStream(ctx, &kinesis.AddTagsToStreamInput{
			StreamName: aws.String(streamName),
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		return err
	}))

	// Test 13: List Tags for Stream
	results = append(results, r.RunTest("kinesis", "ListTagsForStream", func() error {
		resp, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("Tags is nil")
		}
		return nil
	}))

	// Test 14: Remove Tags from Stream
	results = append(results, r.RunTest("kinesis", "RemoveTagsFromStream", func() error {
		_, err := client.RemoveTagsFromStream(ctx, &kinesis.RemoveTagsFromStreamInput{
			StreamName: aws.String(streamName),
			TagKeys:    []string{"Environment"},
		})
		return err
	}))

	// Test 15: Start Stream Encryption
	results = append(results, r.RunTest("kinesis", "StartStreamEncryption", func() error {
		_, err := client.StartStreamEncryption(ctx, &kinesis.StartStreamEncryptionInput{
			StreamName:     aws.String(streamName),
			EncryptionType: types.EncryptionTypeKms,
			KeyId:          aws.String("arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"),
		})
		return err
	}))

	// Test 16: Stop Stream Encryption
	results = append(results, r.RunTest("kinesis", "StopStreamEncryption", func() error {
		_, err := client.StopStreamEncryption(ctx, &kinesis.StopStreamEncryptionInput{
			StreamName:     aws.String(streamName),
			EncryptionType: types.EncryptionTypeKms,
			KeyId:          aws.String("arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"),
		})
		return err
	}))

	// Test 17: Create Stream with Tags
	streamName2 := fmt.Sprintf("test-stream-tags-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("kinesis", "CreateStreamWithTags", func() error {
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName2),
			ShardCount: aws.Int32(1),
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(streamName2),
		})
		return err
	}))

	// Test 18: Increase Stream Retention Period
	results = append(results, r.RunTest("kinesis", "IncreaseStreamRetentionPeriod", func() error {
		_, err := client.IncreaseStreamRetentionPeriod(ctx, &kinesis.IncreaseStreamRetentionPeriodInput{
			StreamName:           aws.String(streamName),
			RetentionPeriodHours: aws.Int32(48),
		})
		return err
	}))

	// Test 19: Decrease Stream Retention Period
	results = append(results, r.RunTest("kinesis", "DecreaseStreamRetentionPeriod", func() error {
		_, err := client.DecreaseStreamRetentionPeriod(ctx, &kinesis.DecreaseStreamRetentionPeriodInput{
			StreamName:           aws.String(streamName),
			RetentionPeriodHours: aws.Int32(24),
		})
		return err
	}))

	// Test 20: Register Stream Consumer
	consumerName := fmt.Sprintf("test-consumer-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("kinesis", "RegisterStreamConsumer", func() error {
		resp, err := client.RegisterStreamConsumer(ctx, &kinesis.RegisterStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		if err != nil {
			return err
		}
		if resp.Consumer == nil {
			return fmt.Errorf("Consumer is nil")
		}
		return nil
	}))

	// Test 21: Describe Stream Consumer
	results = append(results, r.RunTest("kinesis", "DescribeStreamConsumer", func() error {
		resp, err := client.DescribeStreamConsumer(ctx, &kinesis.DescribeStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		if err != nil {
			return err
		}
		if resp.ConsumerDescription == nil {
			return fmt.Errorf("ConsumerDescription is nil")
		}
		if resp.ConsumerDescription.ConsumerARN == nil {
			return fmt.Errorf("ConsumerDescription.ConsumerARN is nil")
		}
		return nil
	}))

	// Test 22: List Stream Consumers
	results = append(results, r.RunTest("kinesis", "ListStreamConsumers", func() error {
		resp, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
			StreamARN: aws.String(streamARN),
		})
		if err != nil {
			return err
		}
		if resp.Consumers == nil {
			return fmt.Errorf("Consumers is nil")
		}
		return nil
	}))

	// Test 23: SubscribeToShard (event stream API)
	var consumerARN string
	if streamARN != "" && consumerName != "" {
		describeConsumerResp, describeConsumerErr := client.DescribeStreamConsumer(ctx, &kinesis.DescribeStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		if describeConsumerErr == nil {
			consumerARN = aws.ToString(describeConsumerResp.ConsumerDescription.ConsumerARN)
		}
	}
	if consumerARN != "" && shardID != "" {
		results = append(results, r.RunTest("kinesis", "SubscribeToShard", func() error {
			type result struct {
				err error
			}
			resultCh := make(chan result, 1)

			go func() {
				subCtx, cancel := context.WithCancel(ctx)
				defer cancel()

				resp, err := client.SubscribeToShard(subCtx, &kinesis.SubscribeToShardInput{
					ConsumerARN: aws.String(consumerARN),
					ShardId:     aws.String(shardID),
					StartingPosition: &types.StartingPosition{
						Type: types.ShardIteratorTypeTrimHorizon,
					},
				})
				if err != nil {
					resultCh <- result{err: fmt.Errorf("SubscribeToShard failed: %v", err)}
					return
				}
				defer resp.GetStream().Close()

				eventCh := resp.GetStream().Events()
				timeout := time.After(10 * time.Second)
				for {
					select {
					case <-timeout:
						resultCh <- result{err: fmt.Errorf("timed out waiting for SubscribeToShard event")}
						return
					case event, ok := <-eventCh:
						if !ok {
							resultCh <- result{err: fmt.Errorf("stream closed without receiving any event")}
							return
						}
						switch v := event.(type) {
						case *types.SubscribeToShardEventStreamMemberSubscribeToShardEvent:
							if len(v.Value.Records) == 0 {
								resultCh <- result{err: fmt.Errorf("SubscribeToShardEvent contained zero records")}
								return
							}
							resultCh <- result{err: nil}
							return
						case *types.UnknownUnionMember:
							resultCh <- result{err: fmt.Errorf("SubscribeToShard returned unknown event type: %s", v.Tag)}
							return
						}
					}
				}
			}()

			select {
			case r := <-resultCh:
				return r.err
			case <-time.After(15 * time.Second):
				return fmt.Errorf("SubscribeToShard test timed out (server event stream may not be closing connection)")
			}
		}))
	} else {
		results = append(results, TestResult{
			Service:  "kinesis",
			TestName: "SubscribeToShard",
			Status:   "SKIP",
			Error:    "consumerARN or shardID not available from prior tests",
		})
	}

	// Test 24: Deregister Stream Consumer
	results = append(results, r.RunTest("kinesis", "DeregisterStreamConsumer", func() error {
		_, err := client.DeregisterStreamConsumer(ctx, &kinesis.DeregisterStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		return err
	}))

	// Test 25: Update Stream Mode (create new stream for this test)
	streamName4 := fmt.Sprintf("test-stream-mode-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("kinesis", "UpdateStreamMode", func() error {
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName4),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(streamName4),
		})
		if err != nil {
			return err
		}
		_, err = client.UpdateStreamMode(ctx, &kinesis.UpdateStreamModeInput{
			StreamARN: descResp.StreamDescription.StreamARN,
			StreamModeDetails: &types.StreamModeDetails{
				StreamMode: types.StreamModeOnDemand,
			},
		})
		if err != nil {
			return err
		}
		_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(streamName4),
		})
		return nil
	}))

	// Test 26: Update Shard Count
	results = append(results, r.RunTest("kinesis", "UpdateShardCount", func() error {
		resp, err := client.UpdateShardCount(ctx, &kinesis.UpdateShardCountInput{
			StreamName:       aws.String(streamName),
			TargetShardCount: aws.Int32(2),
			ScalingType:      types.ScalingTypeUniformScaling,
		})
		if err != nil {
			return err
		}
		if resp.CurrentShardCount == nil {
			return fmt.Errorf("CurrentShardCount is nil")
		}
		return nil
	}))

	// Test 27: Merge Shards - create new stream with 2 shards for merge test
	streamName5 := fmt.Sprintf("test-stream-merge-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("kinesis", "MergeShards", func() error {
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName5),
			ShardCount: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(streamName5),
		})
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
			StreamName:           aws.String(streamName5),
			ShardToMerge:         openShards[0].ShardId,
			AdjacentShardToMerge: openShards[1].ShardId,
		})
		if err != nil {
			return err
		}
		_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(streamName5),
		})
		return nil
	}))

	// Test 28: Split Shard
	results = append(results, r.RunTest("kinesis", "SplitShard", func() error {
		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(streamName),
		})
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
			StreamName:         aws.String(streamName),
			ShardToSplit:       openShard.ShardId,
			NewStartingHashKey: aws.String("9223372036854775808"),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	// Test 29: Describe Limits
	results = append(results, r.RunTest("kinesis", "DescribeLimits", func() error {
		resp, err := client.DescribeLimits(ctx, &kinesis.DescribeLimitsInput{})
		if err != nil {
			return err
		}
		if resp.ShardLimit == nil {
			return fmt.Errorf("ShardLimit is nil")
		}
		return nil
	}))

	// Test 30: Delete Stream
	results = append(results, r.RunTest("kinesis", "DeleteStream", func() error {
		_, err := client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(streamName),
		})
		return err
	}))

	// Test 31: Create Stream for ListShards verification
	streamName3 := fmt.Sprintf("test-stream-list-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("kinesis", "ListShardsWithExclusiveStart", func() error {
		client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName3),
			ShardCount: aws.Int32(1),
		})
		time.Sleep(1 * time.Second)
		resp, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String(streamName3),
		})
		if err != nil {
			return err
		}
		if resp.Shards == nil {
			return fmt.Errorf("Shards is nil")
		}
		client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(streamName3),
		})
		return nil
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("kinesis", "DescribeStream_NonExistent", func() error {
		_, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent stream")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "DeleteStream_NonExistent", func() error {
		_, err := client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent stream")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "CreateStream_DuplicateName", func() error {
		dupStreamName := fmt.Sprintf("dup-stream-%d", time.Now().UnixNano())
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(dupStreamName),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(dupStreamName),
		})

		_, err = client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(dupStreamName),
			ShardCount: aws.Int32(1),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate stream name")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "PutRecord_GetRecords_Roundtrip", func() error {
		rtStreamName := fmt.Sprintf("rt-stream-%d", time.Now().UnixNano())
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(rtStreamName),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{
			StreamName: aws.String(rtStreamName),
		})
		time.Sleep(1 * time.Second)

		testData := []byte("roundtrip-kinesis-data-verify")
		putResp, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(rtStreamName),
			Data:         testData,
			PartitionKey: aws.String("partition-1"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		if putResp.SequenceNumber == nil {
			return fmt.Errorf("sequence number is nil")
		}

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
			StreamName: aws.String(rtStreamName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards")
		}

		iterResp, err := client.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
			StreamName:        aws.String(rtStreamName),
			ShardId:           descResp.StreamDescription.Shards[0].ShardId,
			ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
		})
		if err != nil {
			return fmt.Errorf("get iterator: %v", err)
		}

		getResp, err := client.GetRecords(ctx, &kinesis.GetRecordsInput{
			ShardIterator: iterResp.ShardIterator,
		})
		if err != nil {
			return fmt.Errorf("get records: %v", err)
		}
		if len(getResp.Records) == 0 {
			return fmt.Errorf("no records returned")
		}
		if string(getResp.Records[0].Data) != string(testData) {
			return fmt.Errorf("data mismatch: got %q, want %q", string(getResp.Records[0].Data), string(testData))
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListShards_NonExistentStream", func() error {
		_, err := client.ListShards(ctx, &kinesis.ListShardsInput{
			StreamName: aws.String("nonexistent-stream-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent stream")
		}
		return nil
	}))

	return results
}

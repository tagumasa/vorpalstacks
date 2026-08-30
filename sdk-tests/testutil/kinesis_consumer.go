package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

func (r *TestRunner) kinesisConsumerTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	streamName := kinesisStream(ts, "consumer")
	consumerName := fmt.Sprintf("consumer-%s", ts)

	var streamARN string
	var shardID string

	results = append(results, r.RunTest("kinesis", "RegisterStreamConsumer", func() error {
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		time.Sleep(500 * time.Millisecond)

		descResp, descErr := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
		if descErr != nil {
			return fmt.Errorf("describe: %v", descErr)
		}
		streamARN = aws.ToString(descResp.StreamDescription.StreamARN)
		if len(descResp.StreamDescription.Shards) > 0 {
			shardID = aws.ToString(descResp.StreamDescription.Shards[0].ShardId)
		}

		resp, err := client.RegisterStreamConsumer(ctx, &kinesis.RegisterStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		if err != nil {
			return err
		}
		if resp.Consumer == nil {
			return fmt.Errorf("consumer is nil")
		}
		if aws.ToString(resp.Consumer.ConsumerName) != consumerName {
			return fmt.Errorf("ConsumerName: got %q, want %q", aws.ToString(resp.Consumer.ConsumerName), consumerName)
		}
		if aws.ToString(resp.Consumer.ConsumerARN) == "" {
			return fmt.Errorf("ConsumerARN is empty")
		}
		if resp.Consumer.ConsumerStatus != types.ConsumerStatusActive {
			return fmt.Errorf("ConsumerStatus: expected ACTIVE, got %s", resp.Consumer.ConsumerStatus)
		}
		return nil
	}))

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
		cd := resp.ConsumerDescription
		if aws.ToString(cd.ConsumerARN) == "" {
			return fmt.Errorf("ConsumerARN is empty")
		}
		if aws.ToString(cd.ConsumerName) != consumerName {
			return fmt.Errorf("ConsumerName: got %q, want %q", aws.ToString(cd.ConsumerName), consumerName)
		}
		if cd.ConsumerStatus != types.ConsumerStatusActive {
			return fmt.Errorf("ConsumerStatus: expected ACTIVE, got %s", cd.ConsumerStatus)
		}
		if cd.StreamARN == nil {
			return fmt.Errorf("StreamARN is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListStreamConsumers", func() error {
		resp, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
			StreamARN: aws.String(streamARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Consumers) == 0 {
			return fmt.Errorf("expected at least 1 consumer, got 0")
		}
		found := false
		for _, c := range resp.Consumers {
			if aws.ToString(c.ConsumerName) == consumerName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("consumer %q not found in ListStreamConsumers", consumerName)
		}
		return nil
	}))

	var consumerARN string
	if streamARN != "" && consumerName != "" {
		descConsumerResp, descConsumerErr := client.DescribeStreamConsumer(ctx, &kinesis.DescribeStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		if descConsumerErr == nil {
			consumerARN = aws.ToString(descConsumerResp.ConsumerDescription.ConsumerARN)
		}
	}

	if consumerARN != "" && shardID != "" {
		if _, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),
			Data:         []byte("subscribe-test-data"),
			PartitionKey: aws.String("subscribe-pk"),
		}); err != nil {
			results = append(results, SetupFailResult("kinesis", "PutRecord setup failed: %v", err))
			return results
		}

		results = append(results, r.RunTest("kinesis", "SubscribeToShard", func() error {
			type subResult struct {
				err error
			}
			resultCh := make(chan subResult, 1)

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
					resultCh <- subResult{err: fmt.Errorf("SubscribeToShard failed: %v", err)}
					return
				}
				defer resp.GetStream().Close()

				eventCh := resp.GetStream().Events()
				timeout := time.After(10 * time.Second)
				for {
					select {
					case <-timeout:
						resultCh <- subResult{err: fmt.Errorf("timed out waiting for SubscribeToShard event")}
						return
					case event, ok := <-eventCh:
						if !ok {
							resultCh <- subResult{err: fmt.Errorf("stream closed without receiving any event")}
							return
						}
						switch v := event.(type) {
						case *types.SubscribeToShardEventStreamMemberSubscribeToShardEvent:
							if len(v.Value.Records) == 0 {
								resultCh <- subResult{err: fmt.Errorf("SubscribeToShardEvent contained zero records")}
								return
							}
							resultCh <- subResult{err: nil}
							return
						case *types.UnknownUnionMember:
							resultCh <- subResult{err: fmt.Errorf("SubscribeToShard returned unknown event type: %s", v.Tag)}
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

	results = append(results, r.RunTest("kinesis", "ListStreamConsumers_MaxResultsWindow", func() error {
		// Fill the stream to its per-stream consumer quota so the paging
		// rules are observable: the accepted input window is 1-10000 and
		// every provided value bounds the page.
		for i := 1; i <= 19; i++ {
			if _, err := client.RegisterStreamConsumer(ctx, &kinesis.RegisterStreamConsumerInput{
				StreamARN:    aws.String(streamARN),
				ConsumerName: aws.String(fmt.Sprintf("%s-reader%03d", consumerName, i)),
			}); err != nil {
				return fmt.Errorf("register reader %d: %v", i, err)
			}
		}

		// An omitted MaxResults returns every registered consumer — well
		// below the documented default page of 100 — without a token.
		all, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
			StreamARN: aws.String(streamARN),
		})
		if err != nil {
			return err
		}
		if len(all.Consumers) != 20 {
			return fmt.Errorf("omitted MaxResults: got %d consumers, want 20", len(all.Consumers))
		}
		if all.NextToken != nil {
			return fmt.Errorf("omitted MaxResults: unexpected NextToken for a complete page")
		}

		// A provided value bounds the page and reports the remainder.
		five, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
			StreamARN:  aws.String(streamARN),
			MaxResults: aws.Int32(5),
		})
		if err != nil {
			return err
		}
		if len(five.Consumers) != 5 || five.NextToken == nil {
			return fmt.Errorf("MaxResults 5: got %d consumers, want 5 with a token", len(five.Consumers))
		}

		tail, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
			StreamARN: aws.String(streamARN),
			NextToken: five.NextToken,
		})
		if err != nil {
			return err
		}
		if len(tail.Consumers) != 15 || tail.NextToken != nil {
			return fmt.Errorf("token follow-up: got %d consumers, want 15 without more", len(tail.Consumers))
		}

		for _, maxResults := range []int32{0, 10001} {
			_, err := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
				StreamARN:  aws.String(streamARN),
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

	results = append(results, r.RunTest("kinesis", "RegisterStreamConsumer_PerStreamQuota", func() error {
		// The stream already carries its twenty registered consumers from
		// the preceding rows; the documented per-stream quota is twenty
		// for Provisioned and On-demand Standard streams.
		_, err := client.RegisterStreamConsumer(ctx, &kinesis.RegisterStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(fmt.Sprintf("%s-over", consumerName)),
		})
		if err == nil {
			return fmt.Errorf("expected LimitExceededException at the quota boundary, got success")
		}
		return AssertErrorContains(err, "LimitExceededException")
	}))

	results = append(results, r.RunTest("kinesis", "DeregisterStreamConsumer", func() error {
		_, err := client.DeregisterStreamConsumer(ctx, &kinesis.DeregisterStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
		})
		if err != nil {
			return err
		}

		listResp, listErr := client.ListStreamConsumers(ctx, &kinesis.ListStreamConsumersInput{
			StreamARN: aws.String(streamARN),
		})
		if listErr != nil {
			return fmt.Errorf("list after deregister: %v", listErr)
		}
		for _, c := range listResp.Consumers {
			if aws.ToString(c.ConsumerName) == consumerName {
				return fmt.Errorf("consumer %q should not appear after deregister", consumerName)
			}
		}
		_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})
		return nil
	}))

	return results
}

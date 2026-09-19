package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func (r *TestRunner) kinesisTagTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	streamName := kinesisStream(ts, "tag")
	if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(streamName),
		ShardCount: aws.Int32(1),
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream setup failed: %v", err)}
	}
	if _, err := kinesisDescribeWhenReady(ctx, client, streamName, 10*time.Second); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream setup failed: stream not ready: %v", err)}
	}

	results = append(results, r.RunTest("kinesis", "AddTagsToStream", func() error {
		_, err := client.AddTagsToStream(ctx, &kinesis.AddTagsToStreamInput{
			StreamName: aws.String(streamName),
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		if err != nil {
			return err
		}

		tagResp, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("list tags after add: %v", err)
		}
		tagMap := kinesisTagMap(tagResp.Tags)
		if tagMap["Environment"] != "test" {
			return fmt.Errorf("tag Environment: got %q, want %q", tagMap["Environment"], "test")
		}
		if tagMap["Owner"] != "test-user" {
			return fmt.Errorf("tag Owner: got %q, want %q", tagMap["Owner"], "test-user")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "ListTagsForStream", func() error {
		resp, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 2 {
			return fmt.Errorf("expected >= 2 tags, got %d", len(resp.Tags))
		}
		tagMap := kinesisTagMap(resp.Tags)
		if _, ok := tagMap["Environment"]; !ok {
			return fmt.Errorf("tag Environment not found")
		}
		if _, ok := tagMap["Owner"]; !ok {
			return fmt.Errorf("tag Owner not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "RemoveTagsFromStream", func() error {
		_, err := client.RemoveTagsFromStream(ctx, &kinesis.RemoveTagsFromStreamInput{
			StreamName: aws.String(streamName),
			TagKeys:    []string{"Environment"},
		})
		if err != nil {
			return err
		}

		tagResp, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
			StreamName: aws.String(streamName),
		})
		if err != nil {
			return fmt.Errorf("list tags after remove: %v", err)
		}
		for _, t := range tagResp.Tags {
			if aws.ToString(t.Key) == "Environment" {
				return fmt.Errorf("tag Environment should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.kinesisARNTagTests(ctx, client, ts)...)
	results = append(results, r.kinesisTagVerifyTests(ctx, client, ts)...)
	results = append(results, r.kinesisConsumerTagTests(ctx, client, ts)...)

	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})

	return results
}

func (r *TestRunner) kinesisARNTagTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	tagStreamName := kinesisStream(ts, "tagarn")
	if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
		StreamName: aws.String(tagStreamName),
		ShardCount: aws.Int32(1),
	}); err != nil {
		return []TestResult{SetupFailResult("kinesis", "CreateStream (tagarn) setup failed: %v", err)}
	}
	tagStreamDesc, err := kinesisDescribeWhenReady(ctx, client, tagStreamName, 10*time.Second)
	if err != nil {
		return []TestResult{SetupFailResult("kinesis", "DescribeStream (tagarn) setup failed: %v", err)}
	}
	var tagStreamARN string
	if tagStreamDesc != nil && tagStreamDesc.StreamDescription != nil {
		tagStreamARN = aws.ToString(tagStreamDesc.StreamDescription.StreamARN)
	}

	if tagStreamARN != "" {
		results = append(results, r.RunTest("kinesis", "TagResource", func() error {
			_, err := client.TagResource(ctx, &kinesis.TagResourceInput{
				ResourceARN: aws.String(tagStreamARN),
				Tags: map[string]string{
					"TagTest":  "value1",
					"TagTest2": "value2",
				},
			})
			if err != nil {
				return err
			}

			tagResp, err := client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
				ResourceARN: aws.String(tagStreamARN),
			})
			if err != nil {
				return fmt.Errorf("list tags after TagResource: %v", err)
			}
			tagMap := kinesisTagMap(tagResp.Tags)
			if tagMap["TagTest"] != "value1" {
				return fmt.Errorf("tag TagTest: got %q, want %q", tagMap["TagTest"], "value1")
			}
			if tagMap["TagTest2"] != "value2" {
				return fmt.Errorf("tag TagTest2: got %q, want %q", tagMap["TagTest2"], "value2")
			}
			return nil
		}))

		results = append(results, r.RunTest("kinesis", "ListTagsForResource", func() error {
			resp, err := client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
				ResourceARN: aws.String(tagStreamARN),
			})
			if err != nil {
				return err
			}
			if len(resp.Tags) < 2 {
				return fmt.Errorf("expected >= 2 tags, got %d", len(resp.Tags))
			}
			return nil
		}))

		results = append(results, r.RunTest("kinesis", "UntagResource", func() error {
			_, err := client.UntagResource(ctx, &kinesis.UntagResourceInput{
				ResourceARN: aws.String(tagStreamARN),
				TagKeys:     []string{"TagTest"},
			})
			if err != nil {
				return err
			}
			resp, err := client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
				ResourceARN: aws.String(tagStreamARN),
			})
			if err != nil {
				return err
			}
			for _, t := range resp.Tags {
				if aws.ToString(t.Key) == "TagTest" {
					return fmt.Errorf("TagTest should have been removed")
				}
			}
			tagMap := kinesisTagMap(resp.Tags)
			if _, ok := tagMap["TagTest2"]; !ok {
				return fmt.Errorf("TagTest2 should still be present after untagging TagTest")
			}
			return nil
		}))
	} else {
		results = append(results, TestResult{Service: "kinesis", TestName: "TagResource", Status: "SKIP", Error: "tagStreamARN not available"})
		results = append(results, TestResult{Service: "kinesis", TestName: "ListTagsForResource", Status: "SKIP", Error: "tagStreamARN not available"})
		results = append(results, TestResult{Service: "kinesis", TestName: "UntagResource", Status: "SKIP", Error: "tagStreamARN not available"})
	}

	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(tagStreamName)})
	return results
}

func (r *TestRunner) kinesisTagVerifyTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kinesis", "ListTagsForResource_StreamCreated", func() error {
		sn := kinesisStream(ts, "tlcr")
		_, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(sn),
			ShardCount: aws.Int32(1),
			Tags: map[string]string{
				"CreatedBy": "sdk-test",
				"Project":   "vorpalstacks",
			},
		})
		if err != nil {
			return err
		}
		defer client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(sn)})
		time.Sleep(500 * time.Millisecond)

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
		if err != nil {
			return err
		}

		resp, err := client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
			ResourceARN: descResp.StreamDescription.StreamARN,
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) < 2 {
			return fmt.Errorf("expected >= 2 tags from stream creation, got %d", len(resp.Tags))
		}
		tagMap := kinesisTagMap(resp.Tags)
		if tagMap["CreatedBy"] != "sdk-test" {
			return fmt.Errorf("tag CreatedBy: got %q, want %q", tagMap["CreatedBy"], "sdk-test")
		}
		if tagMap["Project"] != "vorpalstacks" {
			return fmt.Errorf("tag Project: got %q, want %q", tagMap["Project"], "vorpalstacks")
		}
		return nil
	}))

	return results
}

// kinesisConsumerTagTests pins the consumer half of the tag resource
// space: create-time tags reach ListTagsForResource, TagResource and
// UntagResource address the consumer by ARN, and the request-level tag
// bounds (TagMap size on AddTagsToStream, TagKeyList size on
// RemoveTagsFromStream, the Limit range on ListTagsForStream) answer
// InvalidArgumentException.
func (r *TestRunner) kinesisConsumerTagTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	streamName := kinesisStream(ts, "ctagc")
	consumerName := fmt.Sprintf("ctagc-reader-%s", ts)
	var consumerARN string

	results = append(results, r.RunTest("kinesis", "RegisterStreamConsumer_CreateTimeTags", func() error {
		if _, err := client.CreateStream(ctx, &kinesis.CreateStreamInput{
			StreamName: aws.String(streamName),
			ShardCount: aws.Int32(1),
		}); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		descResp, err := kinesisDescribeWhenReady(ctx, client, streamName, 10*time.Second)
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		streamARN := aws.ToString(descResp.StreamDescription.StreamARN)

		resp, err := client.RegisterStreamConsumer(ctx, &kinesis.RegisterStreamConsumerInput{
			StreamARN:    aws.String(streamARN),
			ConsumerName: aws.String(consumerName),
			Tags:         map[string]string{"Origin": "register"},
		})
		if err != nil {
			return err
		}
		if resp.Consumer == nil {
			return fmt.Errorf("consumer is nil")
		}
		consumerARN = aws.ToString(resp.Consumer.ConsumerARN)

		tagResp, err := client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
			ResourceARN: aws.String(consumerARN),
		})
		if err != nil {
			return fmt.Errorf("list tags after register: %v", err)
		}
		if got := kinesisTagMap(tagResp.Tags)["Origin"]; got != "register" {
			return fmt.Errorf("create-time tag Origin: got %q, want %q", got, "register")
		}
		return nil
	}))

	if consumerARN != "" {
		results = append(results, r.RunTest("kinesis", "TagResource_ConsumerARN", func() error {
			if _, err := client.TagResource(ctx, &kinesis.TagResourceInput{
				ResourceARN: aws.String(consumerARN),
				Tags:        map[string]string{"Stage": "api"},
			}); err != nil {
				return err
			}
			tagResp, err := client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
				ResourceARN: aws.String(consumerARN),
			})
			if err != nil {
				return err
			}
			tagMap := kinesisTagMap(tagResp.Tags)
			if len(tagMap) != 2 || tagMap["Origin"] != "register" || tagMap["Stage"] != "api" {
				return fmt.Errorf("consumer tags after TagResource: %v", tagMap)
			}

			if _, err := client.UntagResource(ctx, &kinesis.UntagResourceInput{
				ResourceARN: aws.String(consumerARN),
				TagKeys:     []string{"Origin"},
			}); err != nil {
				return err
			}
			tagResp, err = client.ListTagsForResource(ctx, &kinesis.ListTagsForResourceInput{
				ResourceARN: aws.String(consumerARN),
			})
			if err != nil {
				return err
			}
			tagMap = kinesisTagMap(tagResp.Tags)
			if len(tagMap) != 1 || tagMap["Stage"] != "api" {
				return fmt.Errorf("consumer tags after UntagResource: %v", tagMap)
			}
			return nil
		}))
	} else {
		results = append(results, TestResult{Service: "kinesis", TestName: "TagResource_ConsumerARN", Status: "SKIP", Error: "consumerARN not available"})
	}

	results = append(results, r.RunTest("kinesis", "AddTagsToStream_TagCountLimit", func() error {
		tooMany := make(map[string]string, 51)
		for i := 0; i < 51; i++ {
			tooMany[fmt.Sprintf("k%d", i)] = "v"
		}
		_, err := client.AddTagsToStream(ctx, &kinesis.AddTagsToStreamInput{
			StreamName: aws.String(streamName),
			Tags:       tooMany,
		})
		return expectAWSErrorCode(err, "InvalidArgumentException")
	}))

	results = append(results, r.RunTest("kinesis", "RemoveTagsFromStream_TagKeyCountLimit", func() error {
		tagKeys := make([]string, 51)
		for i := range tagKeys {
			tagKeys[i] = fmt.Sprintf("k%d", i)
		}
		_, err := client.RemoveTagsFromStream(ctx, &kinesis.RemoveTagsFromStreamInput{
			StreamName: aws.String(streamName),
			TagKeys:    tagKeys,
		})
		return expectAWSErrorCode(err, "InvalidArgumentException")
	}))

	results = append(results, r.RunTest("kinesis", "ListTagsForStream_LimitRange", func() error {
		_, err := client.ListTagsForStream(ctx, &kinesis.ListTagsForStreamInput{
			StreamName: aws.String(streamName),
			Limit:      aws.Int32(51),
		})
		return expectAWSErrorCode(err, "InvalidArgumentException")
	}))

	if consumerARN != "" {
		_, _ = client.DeregisterStreamConsumer(ctx, &kinesis.DeregisterStreamConsumerInput{
			ConsumerARN: aws.String(consumerARN),
		})
	}
	_, _ = client.DeleteStream(ctx, &kinesis.DeleteStreamInput{StreamName: aws.String(streamName)})
	return results
}

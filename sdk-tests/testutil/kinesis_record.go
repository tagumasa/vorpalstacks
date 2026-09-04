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

func (r *TestRunner) kinesisRecordTests(ctx context.Context, client *kinesis.Client, ts string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("kinesis", "PutRecord", func() error {
		sn := kinesisStream(ts, "rec")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(sn),
			Data:         []byte("test-data"),
			PartitionKey: aws.String("partition-key-1"),
		})
		if err != nil {
			return err
		}
		if resp.SequenceNumber == nil || *resp.SequenceNumber == "" {
			return fmt.Errorf("SequenceNumber is nil or empty")
		}
		if resp.ShardId == nil || *resp.ShardId == "" {
			return fmt.Errorf("ShardId is nil or empty")
		}
		return nil
	}))

	// PartitionKey @length(1,256) counts Unicode characters (the shape
	// carries no pattern), so a 200-character CJK key (600 bytes) is legal.
	results = append(results, r.RunTest("kinesis", "PutRecord_PartitionKeyMultibyteAccepted", func() error {
		sn := kinesisStream(ts, "recmb")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		key := strings.Repeat("\u65e5", 200)
		resp, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(sn),
			Data:         []byte("test-data"),
			PartitionKey: aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("PutRecord with 200-character CJK partition key: %v", err)
		}
		if resp.SequenceNumber == nil || *resp.SequenceNumber == "" {
			return fmt.Errorf("SequenceNumber is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "PutRecords", func() error {
		sn := kinesisStream(ts, "recs")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		resp, err := client.PutRecords(ctx, &kinesis.PutRecordsInput{
			StreamName: aws.String(sn),
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
			return fmt.Errorf("records is nil")
		}
		if resp.FailedRecordCount != nil && *resp.FailedRecordCount > 0 {
			return fmt.Errorf("FailedRecordCount: %d", *resp.FailedRecordCount)
		}
		for i, rec := range resp.Records {
			if rec.SequenceNumber == nil || *rec.SequenceNumber == "" {
				return fmt.Errorf("record[%d].SequenceNumber is nil or empty", i)
			}
			if rec.ShardId == nil || *rec.ShardId == "" {
				return fmt.Errorf("record[%d].ShardId is nil or empty", i)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "GetShardIterator", func() error {
		sn := kinesisStream(ts, "iter")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards")
		}
		shardID := descResp.StreamDescription.Shards[0].ShardId

		resp, err := client.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
			StreamName:        aws.String(sn),
			ShardId:           shardID,
			ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
		})
		if err != nil {
			return err
		}
		if resp.ShardIterator == nil || *resp.ShardIterator == "" {
			return fmt.Errorf("ShardIterator is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "GetRecords", func() error {
		sn := kinesisStream(ts, "getrec")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 500*time.Millisecond)
		if err != nil {
			return err
		}
		defer cleanup()

		testData := []byte("getrecords-test-data")
		putResp, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(sn),
			Data:         testData,
			PartitionKey: aws.String("pk"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		expectedSeq := aws.ToString(putResp.SequenceNumber)

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		iterResp, err := client.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
			StreamName:        aws.String(sn),
			ShardId:           descResp.StreamDescription.Shards[0].ShardId,
			ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
		})
		if err != nil {
			return fmt.Errorf("get iterator: %v", err)
		}

		resp, err := client.GetRecords(ctx, &kinesis.GetRecordsInput{
			ShardIterator: iterResp.ShardIterator,
		})
		if err != nil {
			return err
		}
		if len(resp.Records) == 0 {
			return fmt.Errorf("expected at least 1 record, got 0")
		}
		rec := resp.Records[0]
		if string(rec.Data) != string(testData) {
			return fmt.Errorf("Data mismatch: got %q, want %q", string(rec.Data), string(testData))
		}
		if aws.ToString(rec.PartitionKey) != "pk" {
			return fmt.Errorf("PartitionKey: got %q, want %q", aws.ToString(rec.PartitionKey), "pk")
		}
		if aws.ToString(rec.SequenceNumber) != expectedSeq {
			return fmt.Errorf("SequenceNumber: got %q, want %q", aws.ToString(rec.SequenceNumber), expectedSeq)
		}
		return nil
	}))

	results = append(results, r.RunTest("kinesis", "PutRecord_GetRecords_Roundtrip", func() error {
		sn := kinesisStream(ts, "rt")
		cleanup, err := kinesisCreateStream(ctx, client, sn, 1, 1*time.Second)
		if err != nil {
			return err
		}
		defer cleanup()

		testData := []byte("roundtrip-kinesis-data-verify")
		putResp, err := client.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(sn),
			Data:         testData,
			PartitionKey: aws.String("partition-1"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		if putResp.SequenceNumber == nil {
			return fmt.Errorf("sequence number is nil")
		}

		descResp, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{StreamName: aws.String(sn)})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(descResp.StreamDescription.Shards) == 0 {
			return fmt.Errorf("no shards")
		}

		iterResp, err := client.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
			StreamName:        aws.String(sn),
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

	return results
}

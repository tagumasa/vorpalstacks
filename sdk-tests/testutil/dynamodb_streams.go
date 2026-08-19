package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodbstreams "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	streamtypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"

	"vorpalstacks-sdk-tests/config"
)

// dynamoDBStreamsTests pins the DynamoDB Streams read path over the AWS
// Streams SDK client: describe a stream, obtain a TRIM_HORIZON shard
// iterator, drain the captured records, poll for new ones, and list the
// account's streams.
func (r *TestRunner) dynamoDBStreamsTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()

	cfg, cfgErr := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if cfgErr != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "Streams_DescribeStream_ReturnsShards",
			Status:   "FAIL",
			Error:    fmt.Sprintf("load config: %v", cfgErr),
		})
	}
	sc := dynamodbstreams.NewFromConfig(cfg)

	tableName := fmt.Sprintf("streams-baseline-%d", suffix)
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
		StreamSpecification: &dynamodbtypes.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: dynamodbtypes.StreamViewTypeNewAndOldImages,
		},
	}); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "Streams_DescribeStream_ReturnsShards",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create table %s: %v", tableName, err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
	if err := waitKinesisDestTableActive(ctx, client, tableName); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "Streams_DescribeStream_ReturnsShards",
			Status:   "FAIL",
			Error:    fmt.Sprintf("wait active: %v", err),
		})
	}

	for _, key := range []string{"item-1", "item-2"} {
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: key},
			},
		}); err != nil {
			return append(results, TestResult{
				Service:  "dynamodb",
				TestName: "Streams_DescribeStream_ReturnsShards",
				Status:   "FAIL",
				Error:    fmt.Sprintf("put %s: %v", key, err),
			})
		}
	}

	tableDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
	if err != nil || tableDesc.Table.LatestStreamArn == nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "Streams_DescribeStream_ReturnsShards",
			Status:   "FAIL",
			Error:    fmt.Sprintf("describe table for stream ARN: %v", err),
		})
	}
	streamArn := *tableDesc.Table.LatestStreamArn

	// Every stream exposes a single shard whose id is derived from the
	// stream ARN, so the id is read from DescribeStream.
	descResp, err := sc.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{StreamArn: aws.String(streamArn)})
	if err != nil || len(descResp.StreamDescription.Shards) == 0 {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "Streams_DescribeStream_ReturnsShards",
			Status:   "FAIL",
			Error:    fmt.Sprintf("describe stream for shard id: %v", err),
		})
	}
	shardID := *descResp.StreamDescription.Shards[0].ShardId

	results = append(results, r.RunTest("dynamodb", "Streams_DescribeStream_ReturnsShards", func() error {
		resp, err := sc.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{StreamArn: aws.String(streamArn)})
		if err != nil {
			return err
		}
		if resp.StreamDescription == nil {
			return fmt.Errorf("no StreamDescription")
		}
		sd := resp.StreamDescription
		if sd.StreamStatus != streamtypes.StreamStatusEnabled {
			return fmt.Errorf("expected StreamStatus=ENABLED, got %v", sd.StreamStatus)
		}
		if len(sd.Shards) == 0 {
			return fmt.Errorf("expected at least one shard")
		}
		if sd.Shards[0].ShardId == nil || *sd.Shards[0].ShardId == "" {
			return fmt.Errorf("shard has no id")
		}
		if sd.StreamViewType != streamtypes.StreamViewTypeNewAndOldImages {
			return fmt.Errorf("expected NEW_AND_OLD_IMAGES, got %v", sd.StreamViewType)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Streams_GetShardIterator_ReadsCapturedRecords", func() error {
		itResp, err := sc.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
			StreamArn:         aws.String(streamArn),
			ShardId:           aws.String(shardID),
			ShardIteratorType: streamtypes.ShardIteratorTypeTrimHorizon,
		})
		if err != nil {
			return err
		}
		if itResp.ShardIterator == nil {
			return fmt.Errorf("no shard iterator returned")
		}

		recResp, err := sc.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: itResp.ShardIterator,
		})
		if err != nil {
			return err
		}
		if len(recResp.Records) != 2 {
			return fmt.Errorf("expected 2 captured records, got %d", len(recResp.Records))
		}
		if recResp.Records[0].EventName != "INSERT" || recResp.Records[1].EventName != "INSERT" {
			return fmt.Errorf("expected INSERT events, got %v/%v", recResp.Records[0].EventName, recResp.Records[1].EventName)
		}
		if recResp.Records[0].Dynamodb.Keys["pk"] == nil {
			return fmt.Errorf("record carries no key")
		}
		if recResp.NextShardIterator == nil {
			return fmt.Errorf("expected NextShardIterator for further polling")
		}

		// The next poll starts after the drained records and stays empty
		// until a new change arrives.
		next, err := sc.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: recResp.NextShardIterator,
		})
		if err != nil {
			return err
		}
		if len(next.Records) != 0 {
			return fmt.Errorf("expected an empty follow-up poll, got %d", len(next.Records))
		}
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: "item-3"},
			},
		}); err != nil {
			return err
		}
		fresh, err := sc.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{
			ShardIterator: next.NextShardIterator,
		})
		if err != nil {
			return err
		}
		if len(fresh.Records) != 1 || fresh.Records[0].EventName != "INSERT" {
			return fmt.Errorf("expected the new INSERT record, got %d records", len(fresh.Records))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Streams_ListStreams_ContainsTableStream", func() error {
		resp, err := sc.ListStreams(ctx, &dynamodbstreams.ListStreamsInput{})
		if err != nil {
			return err
		}
		for _, s := range resp.Streams {
			if s.StreamArn != nil && *s.StreamArn == streamArn {
				return nil
			}
		}
		return fmt.Errorf("stream %s not listed", streamArn)
	}))

	return results
}

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

		// Disabling and re-enabling the stream starts a new generation: a
		// fresh TRIM_HORIZON iterator on the new ARN serves none of the
		// superseded generation's records, and the generation's own write
		// reads back under the new ARN.
		for _, enabled := range []bool{false, true} {
			spec := &dynamodbtypes.StreamSpecification{StreamEnabled: aws.Bool(enabled)}
			if enabled {
				spec.StreamViewType = dynamodbtypes.StreamViewTypeNewAndOldImages
			}
			if _, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
				TableName:           aws.String(tableName),
				StreamSpecification: spec,
			}); err != nil {
				return fmt.Errorf("update table stream enabled=%v: %v", enabled, err)
			}
		}
		reDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err != nil || reDesc.Table.LatestStreamArn == nil {
			return fmt.Errorf("describe table after re-enable: %v", err)
		}
		newArn := *reDesc.Table.LatestStreamArn
		if newArn == streamArn {
			return fmt.Errorf("re-enable must mint a fresh stream ARN, got %q twice", newArn)
		}
		reStream, err := sc.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{StreamArn: aws.String(newArn)})
		if err != nil || len(reStream.StreamDescription.Shards) == 0 {
			return fmt.Errorf("describe re-enabled stream: %v", err)
		}
		reIt, err := sc.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
			StreamArn:         aws.String(newArn),
			ShardId:           reStream.StreamDescription.Shards[0].ShardId,
			ShardIteratorType: streamtypes.ShardIteratorTypeTrimHorizon,
		})
		if err != nil {
			return err
		}
		reRec, err := sc.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{ShardIterator: reIt.ShardIterator})
		if err != nil {
			return err
		}
		if len(reRec.Records) != 0 {
			return fmt.Errorf("re-enabled stream must start empty, got %d records of the superseded generation", len(reRec.Records))
		}
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: "item-4"},
			},
		}); err != nil {
			return err
		}
		reFresh, err := sc.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{ShardIterator: reRec.NextShardIterator})
		if err != nil {
			return err
		}
		if len(reFresh.Records) != 1 || reFresh.Records[0].EventName != "INSERT" {
			return fmt.Errorf("expected the new generation's single INSERT record, got %+v", reFresh.Records)
		}
		return nil
	}))

	// The read-path test above re-enabled the table's stream, so the listed
	// stream is the current generation: the ARN asserted is read from the
	// table's LatestStreamArn at list time.
	results = append(results, r.RunTest("dynamodb", "Streams_ListStreams_ContainsTableStream", func() error {
		desc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err != nil || desc.Table.LatestStreamArn == nil {
			return fmt.Errorf("describe table for stream ARN: %v", err)
		}
		currentArn := *desc.Table.LatestStreamArn
		resp, err := sc.ListStreams(ctx, &dynamodbstreams.ListStreamsInput{})
		if err != nil {
			return err
		}
		for _, s := range resp.Streams {
			if s.StreamArn != nil && *s.StreamArn == currentArn {
				return nil
			}
		}
		return fmt.Errorf("stream %s not listed", currentArn)
	}))

	// Writes committed through ExecuteTransaction capture stream records in
	// the same storage transaction: INSERT, UPDATE and DELETE statements each
	// produce their event.
	results = append(results, r.RunTest("dynamodb", "Streams_ExecuteTransaction_CapturesWriteRecords", func() error {
		txnTable := fmt.Sprintf("streams-executetx-%d", suffix)
		cleanupTable, err := createDynamoTestTable(ctx, client, txnTable,
			withDynamoStream(dynamodbtypes.StreamViewTypeNewAndOldImages))
		if err != nil {
			return err
		}
		defer cleanupTable()
		if err := waitKinesisDestTableActive(ctx, client, txnTable); err != nil {
			return err
		}

		for _, key := range []string{"tx-u", "tx-d"} {
			if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(txnTable),
				Item: map[string]dynamodbtypes.AttributeValue{
					"id":  &dynamodbtypes.AttributeValueMemberS{Value: key},
					"val": &dynamodbtypes.AttributeValueMemberS{Value: "seed"},
				},
			}); err != nil {
				return err
			}
		}

		if _, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []dynamodbtypes.ParameterizedStatement{
				{Statement: aws.String("INSERT INTO \"" + txnTable + "\" VALUE {'id': 'tx-i', 'val': 'inserted'}")},
				{Statement: aws.String("UPDATE \"" + txnTable + "\" SET val = 'updated' WHERE id = 'tx-u'")},
				{Statement: aws.String("DELETE FROM \"" + txnTable + "\" WHERE id = 'tx-d'")},
			},
		}); err != nil {
			return err
		}

		tblDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(txnTable)})
		if err != nil || tblDesc.Table.LatestStreamArn == nil {
			return fmt.Errorf("describe table for stream ARN: %v", err)
		}
		txDesc, err := sc.DescribeStream(ctx, &dynamodbstreams.DescribeStreamInput{StreamArn: tblDesc.Table.LatestStreamArn})
		if err != nil || len(txDesc.StreamDescription.Shards) == 0 {
			return fmt.Errorf("describe stream: %v", err)
		}
		itResp, err := sc.GetShardIterator(ctx, &dynamodbstreams.GetShardIteratorInput{
			StreamArn:         tblDesc.Table.LatestStreamArn,
			ShardId:           txDesc.StreamDescription.Shards[0].ShardId,
			ShardIteratorType: streamtypes.ShardIteratorTypeTrimHorizon,
		})
		if err != nil {
			return err
		}
		recResp, err := sc.GetRecords(ctx, &dynamodbstreams.GetRecordsInput{ShardIterator: itResp.ShardIterator})
		if err != nil {
			return err
		}
		// Two seeded puts then the transaction's three statements.
		if len(recResp.Records) != 5 {
			return fmt.Errorf("expected 5 captured records, got %d", len(recResp.Records))
		}
		wantEvents := []string{"INSERT", "INSERT", "INSERT", "MODIFY", "REMOVE"}
		for i, want := range wantEvents {
			if string(recResp.Records[i].EventName) != want {
				return fmt.Errorf("record %d event = %s, want %s", i, recResp.Records[i].EventName, want)
			}
		}
		modified := recResp.Records[3].Dynamodb
		if modified == nil || modified.NewImage == nil {
			return fmt.Errorf("MODIFY record carries no new image: %+v", modified)
		}
		if v, ok := modified.NewImage["val"].(*streamtypes.AttributeValueMemberS); !ok || v.Value != "updated" {
			return fmt.Errorf("MODIFY new image val mismatch: %v", modified.NewImage["val"])
		}
		return nil
	}))

	return results
}

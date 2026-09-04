package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// dynamoDBItemCollectionMetricsTests exercises the ReturnItemCollectionMetrics
// contract on its own fixtures: item collections exist only on tables with a
// local secondary index, and the single-item writes report one entry object
// while the batched writes report one entry per table.
func (r *TestRunner) dynamoDBItemCollectionMetricsTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	lsiTable := fmt.Sprintf("ICM-LSI-%d", time.Now().UnixNano())
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(lsiTable),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("lsi_sk"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
		},
		BillingMode: types.BillingModePayPerRequest,
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: aws.String("lsi1"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("lsi_sk"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
		},
	}); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "ItemCollectionMetrics_SetupLSITable",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create LSI table: %v", err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(lsiTable)})

	plainTable := fmt.Sprintf("ICM-Plain-%d", time.Now().UnixNano())
	cleanupTable, err := createDynamoTestTable(ctx, client, plainTable)
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "ItemCollectionMetrics_SetupPlainTable",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create plain table: %v", err),
		})
	}
	defer cleanupTable()

	results = append(results, r.RunTest("dynamodb", "PutItem_ItemCollectionMetrics", func() error {
		resp, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:                   aws.String(lsiTable),
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
			Item: map[string]types.AttributeValue{
				"pk":     &types.AttributeValueMemberS{Value: "icm-p1"},
				"sk":     &types.AttributeValueMemberS{Value: "a"},
				"lsi_sk": &types.AttributeValueMemberS{Value: "x"},
				"val":    &types.AttributeValueMemberS{Value: "one"},
			},
		})
		if err != nil {
			return err
		}
		return assertSingleItemCollectionMetrics(resp.ItemCollectionMetrics, "icm-p1")
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_ItemCollectionMetrics", func() error {
		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                   aws.String(lsiTable),
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "icm-p1"},
				"sk": &types.AttributeValueMemberS{Value: "a"},
			},
			UpdateExpression: aws.String("SET val = :v"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberS{Value: "updated"},
			},
		})
		if err != nil {
			return err
		}
		return assertSingleItemCollectionMetrics(resp.ItemCollectionMetrics, "icm-p1")
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteItem_ItemCollectionMetrics", func() error {
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(lsiTable),
			Item: map[string]types.AttributeValue{
				"pk":     &types.AttributeValueMemberS{Value: "icm-p2"},
				"sk":     &types.AttributeValueMemberS{Value: "b"},
				"lsi_sk": &types.AttributeValueMemberS{Value: "y"},
			},
		}); err != nil {
			return fmt.Errorf("seed item for delete: %w", err)
		}
		resp, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName:                   aws.String(lsiTable),
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "icm-p2"},
				"sk": &types.AttributeValueMemberS{Value: "b"},
			},
		})
		if err != nil {
			return err
		}
		return assertSingleItemCollectionMetrics(resp.ItemCollectionMetrics, "icm-p2")
	}))

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_ItemCollectionMetrics", func() error {
		resp, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
			RequestItems: map[string][]types.WriteRequest{
				lsiTable: {
					{
						PutRequest: &types.PutRequest{
							Item: map[string]types.AttributeValue{
								"pk":     &types.AttributeValueMemberS{Value: "icm-batch"},
								"sk":     &types.AttributeValueMemberS{Value: "c"},
								"lsi_sk": &types.AttributeValueMemberS{Value: "z"},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.ItemCollectionMetrics == nil {
			return fmt.Errorf("expected ItemCollectionMetrics for ReturnItemCollectionMetrics=SIZE")
		}
		entries, ok := resp.ItemCollectionMetrics[lsiTable]
		if !ok || len(entries) == 0 {
			return fmt.Errorf("expected an ItemCollectionMetrics entry for table %s, got %v", lsiTable, resp.ItemCollectionMetrics)
		}
		return assertSingleItemCollectionMetrics(&entries[0], "icm-batch")
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_ItemCollectionMetrics", func() error {
		resp, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
			TransactItems: []types.TransactWriteItem{
				{
					Put: &types.Put{
						TableName: aws.String(lsiTable),
						Item: map[string]types.AttributeValue{
							"pk":     &types.AttributeValueMemberS{Value: "icm-txn"},
							"sk":     &types.AttributeValueMemberS{Value: "d"},
							"lsi_sk": &types.AttributeValueMemberS{Value: "w"},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.ItemCollectionMetrics == nil {
			return fmt.Errorf("expected ItemCollectionMetrics for ReturnItemCollectionMetrics=SIZE")
		}
		entries, ok := resp.ItemCollectionMetrics[lsiTable]
		if !ok || len(entries) == 0 {
			return fmt.Errorf("expected an ItemCollectionMetrics entry for table %s, got %v", lsiTable, resp.ItemCollectionMetrics)
		}
		return assertSingleItemCollectionMetrics(&entries[0], "icm-txn")
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_ItemCollectionMetrics_NoLSIOmitted", func() error {
		resp, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:                   aws.String(plainTable),
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "icm-plain"},
				"val": &types.AttributeValueMemberS{Value: "no index"},
			},
		})
		if err != nil {
			return err
		}
		if resp.ItemCollectionMetrics != nil {
			return fmt.Errorf("expected no ItemCollectionMetrics for a table without local secondary indexes, got %v", resp.ItemCollectionMetrics)
		}
		return nil
	}))

	return results
}

// assertSingleItemCollectionMetrics verifies one ItemCollectionMetrics
// entry: the written item's partition key as the ItemCollectionKey and a
// two-element size estimate range.
func assertSingleItemCollectionMetrics(metrics *types.ItemCollectionMetrics, wantPK string) error {
	if metrics == nil {
		return fmt.Errorf("expected ItemCollectionMetrics for ReturnItemCollectionMetrics=SIZE")
	}
	if metrics.ItemCollectionKey == nil {
		return fmt.Errorf("ItemCollectionKey missing from the metrics entry")
	}
	pk, ok := metrics.ItemCollectionKey["pk"].(*types.AttributeValueMemberS)
	if !ok || pk.Value != wantPK {
		return fmt.Errorf("ItemCollectionKey missing partition key %q: %v", wantPK, metrics.ItemCollectionKey)
	}
	if len(metrics.SizeEstimateRangeGB) != 2 {
		return fmt.Errorf("expected a two-element SizeEstimateRangeGB, got %v", metrics.SizeEstimateRangeGB)
	}
	return nil
}

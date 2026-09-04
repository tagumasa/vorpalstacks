package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBBasicBatchTests(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem", func() error {
		resp, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: {
					{
						PutRequest: &types.PutRequest{
							Item: map[string]types.AttributeValue{
								"id":   &types.AttributeValueMemberS{Value: "batch1"},
								"data": &types.AttributeValueMemberS{Value: "batch item 1"},
							},
						},
					},
					{
						PutRequest: &types.PutRequest{
							Item: map[string]types.AttributeValue{
								"id":   &types.AttributeValueMemberS{Value: "batch2"},
								"data": &types.AttributeValueMemberS{Value: "batch item 2"},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.UnprocessedItems == nil {
			return fmt.Errorf("UnprocessedItems is nil")
		}
		verifyResp, verifyErr := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				tableName: {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "batch1"}},
						{"id": &types.AttributeValueMemberS{Value: "batch2"}},
					},
				},
			},
		})
		if verifyErr != nil {
			return fmt.Errorf("BatchWriteItem verification BatchGetItem failed: %w", verifyErr)
		}
		if items, ok := verifyResp.Responses[tableName]; !ok || len(items) != 2 {
			return fmt.Errorf("BatchWriteItem verification: expected 2 items, got %d", len(items))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchGetItem", func() error {
		resp, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				tableName: {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "batch1"}},
						{"id": &types.AttributeValueMemberS{Value: "batch2"}},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if len(resp.Responses) == 0 {
			return fmt.Errorf("no responses")
		}
		items, ok := resp.Responses[tableName]
		if !ok || len(items) != 2 {
			return fmt.Errorf("expected 2 items from BatchGetItem, got %d", len(items))
		}
		for _, item := range items {
			if id, ok := item["id"].(*types.AttributeValueMemberS); ok {
				if id.Value != "batch1" && id.Value != "batch2" {
					return fmt.Errorf("unexpected item id: %s", id.Value)
				}
			} else {
				return fmt.Errorf("expected 'id' attribute in batch item")
			}
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBBatchNonExistentTest(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "BatchGetItem_NonExistentTable", func() error {
		_, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				"NonExistentTable_xyz": {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "k"}},
					},
				},
			},
		})
		return expectResourceNotFound(err)
	}))

	return results
}

func (r *TestRunner) dynamoDBBatchEdgeCaseTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	// === BATCH OPERATION EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_DeleteRequest", func() error {
		bwTable := fmt.Sprintf("BWDel-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, bwTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(bwTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "del1"},
				"name": &types.AttributeValueMemberS{Value: "ToDelete"},
			},
		})

		_, err = client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				bwTable: {
					{
						DeleteRequest: &types.DeleteRequest{
							Key: map[string]types.AttributeValue{
								"id": &types.AttributeValueMemberS{Value: "del1"},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("BatchWriteItem DeleteRequest: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(bwTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "del1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.Item) != 0 {
			return fmt.Errorf("item should be deleted after BatchWriteItem DeleteRequest")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchGetItem_Projection", func() error {
		bgTable := fmt.Sprintf("BGProj-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, bgTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(bgTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "bp1"},
				"name":  &types.AttributeValueMemberS{Value: "Alice"},
				"email": &types.AttributeValueMemberS{Value: "alice@test.com"},
			},
		})

		resp, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				bgTable: {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "bp1"}},
					},
					ProjectionExpression: aws.String("id, name"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("batch get: %v", err)
		}
		items := resp.Responses[bgTable]
		if len(items) != 1 {
			return fmt.Errorf("expected 1 item, got %d", len(items))
		}
		if len(items[0]) != 2 {
			return fmt.Errorf("expected 2 projected attributes, got %d", len(items[0]))
		}
		if _, ok := items[0]["id"]; !ok {
			return fmt.Errorf("expected 'id' in projected item")
		}
		if _, ok := items[0]["name"]; !ok {
			return fmt.Errorf("expected 'name' in projected item")
		}
		if _, ok := items[0]["email"]; ok {
			return fmt.Errorf("did not expect 'email' in projected item")
		}
		nameVal, ok := items[0]["name"].(*types.AttributeValueMemberS)
		if !ok || nameVal.Value != "Alice" {
			return fmt.Errorf("expected name=Alice, got %v", items[0]["name"])
		}
		return nil
	}))

	return results
}

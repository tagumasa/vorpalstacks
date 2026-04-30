package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBBasicTransactionTests(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems", func() error {
		resp, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Put: &types.Put{
						TableName: aws.String(tableName),
						Item: map[string]types.AttributeValue{
							"id":   &types.AttributeValueMemberS{Value: "transact1"},
							"name": &types.AttributeValueMemberS{Value: "Transact Item"},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("TransactWriteItems response is nil")
		}
		verifyResp, verifyErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "transact1"},
			},
		})
		if verifyErr != nil {
			return fmt.Errorf("TransactWriteItems verification failed: %w", verifyErr)
		}
		if verifyResp.Item == nil {
			return fmt.Errorf("item not found after TransactWriteItems Put")
		}
		name, ok := verifyResp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Transact Item" {
			return fmt.Errorf("expected name='Transact Item', got %v", verifyResp.Item["name"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactGetItems", func() error {
		resp, err := client.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
			TransactItems: []types.TransactGetItem{
				{
					Get: &types.Get{
						TableName: aws.String(tableName),
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "transact1"},
						},
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
		item := resp.Responses[0].Item
		if item == nil {
			return fmt.Errorf("expected non-nil Item in TransactGetItems response")
		}
		name, ok := item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Transact Item" {
			return fmt.Errorf("expected name='Transact Item', got %v", item["name"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteTransaction", func() error {
		resp, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String("SELECT * FROM \"" + tableName + "\" WHERE id = 'transact1'"),
				},
			},
		})
		if err != nil {
			return err
		}
		if len(resp.Responses) == 0 {
			return fmt.Errorf("no responses")
		}
		if resp.Responses[0].Item == nil {
			return fmt.Errorf("expected non-nil Item in ExecuteTransaction response")
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBTransactionEdgeCaseTests(ctx context.Context, client *dynamodb.Client, compTableName string) []TestResult {
	var results []TestResult

	// === TRANSACTION EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_MultipleOps", func() error {
		twTable := fmt.Sprintf("TW-Multi-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(twTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(twTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(twTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "t1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						TableName: aws.String(twTable),
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "t1"},
						},
						UpdateExpression:          aws.String("SET #v = #v + :inc"),
						ExpressionAttributeNames:  map[string]string{"#v": "val"},
						ExpressionAttributeValues: map[string]types.AttributeValue{":inc": &types.AttributeValueMemberN{Value: "5"}},
					},
				},
				{
					Put: &types.Put{
						TableName: aws.String(twTable),
						Item: map[string]types.AttributeValue{
							"id":  &types.AttributeValueMemberS{Value: "t2"},
							"val": &types.AttributeValueMemberN{Value: "42"},
						},
					},
				},
				{
					ConditionCheck: &types.ConditionCheck{
						TableName: aws.String(twTable),
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "t1"},
						},
						ConditionExpression: aws.String("attribute_exists(id)"),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("transact write: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(twTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "t1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get t1: %v", err)
		}
		val, ok := resp.Item["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "15" {
			return fmt.Errorf("expected val=15 after transact update, got %v", val)
		}

		resp2, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(twTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "t2"},
			},
		})
		if err != nil {
			return fmt.Errorf("get t2: %v", err)
		}
		if resp2.Item == nil {
			return fmt.Errorf("expected t2 to be created")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_ConditionFail", func() error {
		twTable := fmt.Sprintf("TW-CondF-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(twTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(twTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(twTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "nofail"},
				"val": &types.AttributeValueMemberN{Value: "1"},
			},
		})

		_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						TableName: aws.String(twTable),
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "nofail"},
						},
						UpdateExpression:         aws.String("SET #v = :x"),
						ConditionExpression:      aws.String("#v = :expect"),
						ExpressionAttributeNames: map[string]string{"#v": "val"},
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":x":      &types.AttributeValueMemberN{Value: "2"},
							":expect": &types.AttributeValueMemberN{Value: "999"},
						},
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected TransactionCanceledException")
		}
		errMsg := err.Error()
		if !strings.Contains(errMsg, "TransactionCanceledException") {
			return fmt.Errorf("expected TransactionCanceledException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_Delete", func() error {
		twTable := fmt.Sprintf("TW-Del-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(twTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(twTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(twTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "td1"},
				"name": &types.AttributeValueMemberS{Value: "ToDelete"},
			},
		})

		_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Delete: &types.Delete{
						TableName: aws.String(twTable),
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "td1"},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("transact delete: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(twTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "td1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(resp.Item) != 0 {
			return fmt.Errorf("item should be deleted after TransactWriteItems Delete")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactGetItems_NonExistentKey", func() error {
		tgTable := fmt.Sprintf("TG-NE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tgTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tgTable)})

		resp, err := client.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
			TransactItems: []types.TransactGetItem{
				{
					Get: &types.Get{
						TableName: aws.String(tgTable),
						Key: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberS{Value: "nonexistent"},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("transact get: %v", err)
		}
		if len(resp.Responses) != 1 {
			return fmt.Errorf("expected 1 response, got %d", len(resp.Responses))
		}
		if resp.Responses[0].Item != nil {
			return fmt.Errorf("expected nil Item for non-existent key")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteTransaction_MixedReadWrite", func() error {
		_, err := client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String("SELECT * FROM \"" + compTableName + "\" WHERE pk = 'user1'"),
				},
				{
					Statement: aws.String("INSERT INTO \"" + compTableName + "\" VALUE {'pk': 'user3', 'sk': 'meta', 'name': 'Charlie'}"),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected TransactionConflictException for mixed read/write")
		}
		if !strings.Contains(err.Error(), "TransactionConflictException") {
			return fmt.Errorf("expected TransactionConflictException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteTransaction_WriteOnly", func() error {
		etTable := fmt.Sprintf("ETWrite-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(etTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(etTable)})

		_, err = client.ExecuteTransaction(ctx, &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String("INSERT INTO \"" + etTable + "\" VALUE {'id': 'et1', 'val': 'hello'}"),
				},
				{
					Statement: aws.String("INSERT INTO \"" + etTable + "\" VALUE {'id': 'et2', 'val': 'world'}"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("execute transaction: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(etTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "et1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Item == nil {
			return fmt.Errorf("expected et1 to be created")
		}
		resp2, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(etTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "et2"},
			},
		})
		if err != nil {
			return fmt.Errorf("get et2: %v", err)
		}
		if resp2.Item == nil {
			return fmt.Errorf("expected et2 to be created")
		}
		val, ok := resp2.Item["val"].(*types.AttributeValueMemberS)
		if !ok || val.Value != "world" {
			return fmt.Errorf("expected et2 val=world, got %v", resp2.Item["val"])
		}
		return nil
	}))

	return results
}

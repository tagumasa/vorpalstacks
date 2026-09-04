package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBBasicCRUD(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "PutItem", func() error {
		resp, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "test1"},
				"name":  &types.AttributeValueMemberS{Value: "Test Item"},
				"count": &types.AttributeValueMemberN{Value: "42"},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("PutItem response is nil")
		}
		verifyResp, verifyErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "test1"}},
		})
		if verifyErr != nil {
			return fmt.Errorf("PutItem verification GetItem failed: %w", verifyErr)
		}
		if verifyResp.Item == nil {
			return fmt.Errorf("PutItem verification: item not found after PutItem")
		}
		if count, ok := verifyResp.Item["count"].(*types.AttributeValueMemberN); !ok || count.Value != "42" {
			return fmt.Errorf("PutItem verification: count mismatch, expected 42")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "GetItem", func() error {
		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "test1"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Item == nil {
			return fmt.Errorf("item not found")
		}
		name, ok := resp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Test Item" {
			return fmt.Errorf("item name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem", func() error {
		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "test1"},
			},
			UpdateExpression: aws.String("SET #n = :name"),
			ExpressionAttributeNames: map[string]string{
				"#n": "name",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name": &types.AttributeValueMemberS{Value: "Updated"},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			return err
		}
		if resp.Attributes == nil {
			return fmt.Errorf("attributes not found")
		}
		updatedName, ok := resp.Attributes["name"].(*types.AttributeValueMemberS)
		if !ok || updatedName.Value != "Updated" {
			return fmt.Errorf("name mismatch after update: got %v", updatedName)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBReturnValueTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "PutItem_ReturnValues", func() error {
		rvTable := fmt.Sprintf("RVTable-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, rvTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		resp, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(rvTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "rv1"},
				"name":  &types.AttributeValueMemberS{Value: "Alice"},
				"count": &types.AttributeValueMemberN{Value: "10"},
			},
			ReturnValues: types.ReturnValueAllOld,
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}
		if resp.Attributes != nil {
			return fmt.Errorf("first PutItem with ReturnValues=ALL_OLD should have nil Attributes for new item, got %v", resp.Attributes)
		}

		resp2, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(rvTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "rv1"},
				"name":  &types.AttributeValueMemberS{Value: "Bob"},
				"count": &types.AttributeValueMemberN{Value: "20"},
			},
			ReturnValues: types.ReturnValueAllOld,
		})
		if err != nil {
			return fmt.Errorf("put item 2: %v", err)
		}
		if resp2.Attributes == nil {
			return fmt.Errorf("second PutItem with ReturnValues=ALL_OLD should return old attributes")
		}
		oldName, ok := resp2.Attributes["name"].(*types.AttributeValueMemberS)
		if !ok || oldName.Value != "Alice" {
			return fmt.Errorf("old name should be 'Alice', got %v", oldName)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_ReturnUpdatedAttributes", func() error {
		uaTable := fmt.Sprintf("UATable-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, uaTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(uaTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "ua1"},
				"val": &types.AttributeValueMemberN{Value: "0"},
				"tags": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "a"},
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(uaTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ua1"},
			},
			UpdateExpression: aws.String("ADD #v :inc, SET #t = list_append(#t, :newTag)"),
			ExpressionAttributeNames: map[string]string{
				"#v": "val",
				"#t": "tags",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":inc":    &types.AttributeValueMemberN{Value: "5"},
				":newTag": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "b"}}},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected updated attributes")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "5" {
			return fmt.Errorf("expected val=5, got %v", val)
		}
		tags, ok := resp.Attributes["tags"].(*types.AttributeValueMemberL)
		if !ok {
			return fmt.Errorf("expected list for tags, got %T", resp.Attributes["tags"])
		}
		if len(tags.Value) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tags.Value))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_ReturnConsumedCapacity", func() error {
		qTable := fmt.Sprintf("QCapTable-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, qTable, withDynamoHashKey("pk"))
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(qTable),
			Item: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "key1"},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(qTable),
			KeyConditionExpression: aws.String("pk = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "key1"},
			},
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}
		if resp.ConsumedCapacity == nil {
			return fmt.Errorf("expected ConsumedCapacity in response")
		}
		if resp.ConsumedCapacity.TableName == nil || *resp.ConsumedCapacity.TableName != qTable {
			return fmt.Errorf("ConsumedCapacity.TableName mismatch, got %v", resp.ConsumedCapacity.TableName)
		}
		return nil
	}))

	return results
}

// === PUT ITEM EDGE CASES ===

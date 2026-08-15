package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBProjectionTests(ctx context.Context, client *dynamodb.Client, compTableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "GetItem_ProjectionExpression", func() error {
		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(compTableName),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "user1"},
				"sk": &types.AttributeValueMemberS{Value: "meta"},
			},
			ProjectionExpression: aws.String("name, age"),
		})
		if err != nil {
			return err
		}
		if len(resp.Item) != 2 {
			return fmt.Errorf("expected 2 projected attributes, got %d", len(resp.Item))
		}
		if _, ok := resp.Item["name"]; !ok {
			return fmt.Errorf("expected 'name' in projection")
		}
		if _, ok := resp.Item["age"]; !ok {
			return fmt.Errorf("expected 'age' in projection")
		}
		if _, ok := resp.Item["pk"]; ok {
			return fmt.Errorf("did not expect 'pk' in projection")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "GetItem_ProjectionWithAttrNames", func() error {
		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(compTableName),
			Key: map[string]types.AttributeValue{
				"pk": &types.AttributeValueMemberS{Value: "user1"},
				"sk": &types.AttributeValueMemberS{Value: "meta"},
			},
			ProjectionExpression: aws.String("#n, #a"),
			ExpressionAttributeNames: map[string]string{
				"#n": "name",
				"#a": "age",
			},
		})
		if err != nil {
			return err
		}
		if len(resp.Item) != 2 {
			return fmt.Errorf("expected 2 projected attributes, got %d", len(resp.Item))
		}
		name, ok := resp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Alice" {
			return fmt.Errorf("expected name=Alice, got %v", resp.Item["name"])
		}
		age, ok := resp.Item["age"].(*types.AttributeValueMemberN)
		if !ok || age.Value != "30" {
			return fmt.Errorf("expected age=30, got %v", resp.Item["age"])
		}
		return nil
	}))

	return results
}

// === DELETE ITEM EDGE CASES ===

func (r *TestRunner) dynamoDBDeleteEdgeCaseTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "DeleteItem_NonExistentKey_NoCondition", func() error {
		delTable := fmt.Sprintf("DelNE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(delTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(delTable)})

		_, err = client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(delTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "nonexistent"},
			},
		})
		if err != nil {
			return fmt.Errorf("DeleteItem non-existent key without condition should succeed, got: %v", err)
		}
		scanResp, scanErr := client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(delTable)})
		if scanErr != nil {
			return fmt.Errorf("scan after delete: %v", scanErr)
		}
		if scanResp.Count != 0 {
			return fmt.Errorf("expected empty table after deleting non-existent key, got %d items", scanResp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteItem_ReturnValuesAllOld", func() error {
		rvDelTable := fmt.Sprintf("RVDel-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(rvDelTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(rvDelTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(rvDelTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "rvdel1"},
				"name":  &types.AttributeValueMemberS{Value: "ToDelete"},
				"count": &types.AttributeValueMemberN{Value: "99"},
			},
		})

		resp, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(rvDelTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "rvdel1"},
			},
			ReturnValues: types.ReturnValueAllOld,
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected old attributes in response")
		}
		oldName, ok := resp.Attributes["name"].(*types.AttributeValueMemberS)
		if !ok || oldName.Value != "ToDelete" {
			return fmt.Errorf("expected old name 'ToDelete', got %v", oldName)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteItem_ReturnValuesAllOld_NonExistent", func() error {
		rvDelTable := fmt.Sprintf("RVDelNE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(rvDelTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(rvDelTable)})

		resp, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(rvDelTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "nonexistent"},
			},
			ReturnValues: types.ReturnValueAllOld,
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		if resp.Attributes != nil {
			return fmt.Errorf("expected nil Attributes for non-existent key, got %v", resp.Attributes)
		}
		return nil
	}))

	return results
}

// === UPDATE ITEM EDGE CASES ===

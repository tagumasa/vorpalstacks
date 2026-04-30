package testutil

import (
	"bytes"
	"context"
	"errors"
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
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(rvTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(rvTable)})

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
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(uaTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(uaTable)})

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
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(qTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(qTable)})

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

func (r *TestRunner) dynamoDBPutConditionTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "PutItem_ConditionPass", func() error {
		condTable := fmt.Sprintf("CondPut-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(condTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(condTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(condTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "cp1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
			ConditionExpression: aws.String("attribute_not_exists(id)"),
		})
		if err != nil {
			return fmt.Errorf("PutItem with attribute_not_exists should succeed: %v", err)
		}
		verifyResp, verifyErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(condTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "cp1"},
			},
		})
		if verifyErr != nil {
			return fmt.Errorf("verify get: %v", verifyErr)
		}
		if verifyResp.Item == nil {
			return fmt.Errorf("item not found after conditional PutItem")
		}
		val, ok := verifyResp.Item["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "10" {
			return fmt.Errorf("expected val=10, got %v", verifyResp.Item["val"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_ConditionFail", func() error {
		condTable := fmt.Sprintf("CondPutF-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(condTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(condTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(condTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "cpf1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(condTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "cpf1"},
				"val": &types.AttributeValueMemberN{Value: "20"},
			},
			ConditionExpression: aws.String("attribute_not_exists(id)"),
		})
		if err == nil {
			return fmt.Errorf("expected ConditionalCheckFailedException")
		}
		var ccf *types.ConditionalCheckFailedException
		if !errors.As(err, &ccf) {
			return fmt.Errorf("expected ConditionalCheckFailedException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_ReturnConsumedCapacity", func() error {
		rcTable := fmt.Sprintf("RCapP-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(rcTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(rcTable)})

		resp, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(rcTable),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "rc1"},
			},
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		if resp.ConsumedCapacity == nil {
			return fmt.Errorf("expected ConsumedCapacity in PutItem response")
		}
		if resp.ConsumedCapacity.TableName == nil || *resp.ConsumedCapacity.TableName != rcTable {
			return fmt.Errorf("ConsumedCapacity.TableName mismatch, got %v", resp.ConsumedCapacity.TableName)
		}
		if resp.ConsumedCapacity.CapacityUnits == nil || *resp.ConsumedCapacity.CapacityUnits <= 0 {
			return fmt.Errorf("expected positive CapacityUnits, got %v", resp.ConsumedCapacity.CapacityUnits)
		}
		return nil
	}))

	return results
}

// === GET ITEM EDGE CASES ===

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

func (r *TestRunner) dynamoDBUpdateExpressionTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "UpdateItem_CreateNonExistent", func() error {
		uaTable := fmt.Sprintf("UACreate-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(uaTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(uaTable)})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(uaTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "new1"},
			},
			UpdateExpression: aws.String("SET #n = :name"),
			ExpressionAttributeNames: map[string]string{
				"#n": "name",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":name": &types.AttributeValueMemberS{Value: "CreatedViaUpdate"},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected attributes")
		}
		idVal, ok := resp.Attributes["id"].(*types.AttributeValueMemberS)
		if !ok || idVal.Value != "new1" {
			return fmt.Errorf("expected id=new1, got %v", idVal)
		}
		nameVal, ok := resp.Attributes["name"].(*types.AttributeValueMemberS)
		if !ok || nameVal.Value != "CreatedViaUpdate" {
			return fmt.Errorf("expected name=CreatedViaUpdate, got %v", nameVal)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_IfNotExists", func() error {
		ineTable := fmt.Sprintf("INETable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ineTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ineTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ineTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "ine1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(ineTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ine1"},
			},
			UpdateExpression:          aws.String("SET #v = if_not_exists(#v, :zero) + :inc"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":zero": &types.AttributeValueMemberN{Value: "0"}, ":inc": &types.AttributeValueMemberN{Value: "5"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected attributes")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "15" {
			return fmt.Errorf("expected val=15 (10+5), got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_IfNotExists_NoExisting", func() error {
		ineTable := fmt.Sprintf("INENE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ineTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ineTable)})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(ineTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "inene1"},
			},
			UpdateExpression:          aws.String("SET #v = if_not_exists(#v, :zero) + :inc"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":zero": &types.AttributeValueMemberN{Value: "0"}, ":inc": &types.AttributeValueMemberN{Value: "5"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "5" {
			return fmt.Errorf("expected val=5 (0+5), got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_Arithmetic", func() error {
		arithTable := fmt.Sprintf("Arith-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(arithTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(arithTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(arithTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "a1"},
				"val": &types.AttributeValueMemberN{Value: "100"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(arithTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "a1"},
			},
			UpdateExpression:          aws.String("SET #v = #v - :dec"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":dec": &types.AttributeValueMemberN{Value: "30"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "70" {
			return fmt.Errorf("expected val=70 (100-30), got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_Remove", func() error {
		rmTable := fmt.Sprintf("RmTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(rmTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(rmTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(rmTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "rm1"},
				"name":  &types.AttributeValueMemberS{Value: "Alice"},
				"email": &types.AttributeValueMemberS{Value: "alice@test.com"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(rmTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "rm1"},
			},
			UpdateExpression: aws.String("REMOVE email"),
			ReturnValues:     types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if _, ok := resp.Attributes["email"]; ok {
			return fmt.Errorf("expected 'email' to be removed")
		}
		if _, ok := resp.Attributes["name"]; !ok {
			return fmt.Errorf("expected 'name' to remain")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_AddNumber", func() error {
		addTable := fmt.Sprintf("AddN-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(addTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(addTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(addTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "an1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(addTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "an1"},
			},
			UpdateExpression:          aws.String("ADD #v :inc"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":inc": &types.AttributeValueMemberN{Value: "5"}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "15" {
			return fmt.Errorf("expected val=15, got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_AddStringSet", func() error {
		ssTable := fmt.Sprintf("AddSS-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ssTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ssTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ssTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ss1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"a", "b"}},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(ssTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ss1"},
			},
			UpdateExpression:          aws.String("ADD #t :newTags"),
			ExpressionAttributeNames:  map[string]string{"#t": "tags"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":newTags": &types.AttributeValueMemberSS{Value: []string{"b", "c"}}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		tags, ok := resp.Attributes["tags"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("expected SS type for tags")
		}
		if len(tags.Value) != 3 {
			return fmt.Errorf("expected 3 tags (a,b,c), got %d", len(tags.Value))
		}
		tagMap := make(map[string]bool)
		for _, t := range tags.Value {
			tagMap[t] = true
		}
		for _, exp := range []string{"a", "b", "c"} {
			if !tagMap[exp] {
				return fmt.Errorf("expected tag %q in set", exp)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_DeleteStringSet", func() error {
		dsTable := fmt.Sprintf("DelSS-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(dsTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(dsTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(dsTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ds1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(dsTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ds1"},
			},
			UpdateExpression:          aws.String("DELETE #t :remove"),
			ExpressionAttributeNames:  map[string]string{"#t": "tags"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":remove": &types.AttributeValueMemberSS{Value: []string{"a", "c"}}},
			ReturnValues:              types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		tags, ok := resp.Attributes["tags"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("expected SS type for tags")
		}
		if len(tags.Value) != 1 || tags.Value[0] != "b" {
			return fmt.Errorf("expected tags=[b], got %v", tags.Value)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_UpdatedOld", func() error {
		uoTable := fmt.Sprintf("UO-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(uoTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(uoTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(uoTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "uo1"},
				"val":  &types.AttributeValueMemberN{Value: "10"},
				"name": &types.AttributeValueMemberS{Value: "Old"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(uoTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "uo1"},
			},
			UpdateExpression:          aws.String("SET #v = :new"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":new": &types.AttributeValueMemberN{Value: "20"}},
			ReturnValues:              types.ReturnValueUpdatedOld,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected updated old attributes")
		}
		if _, ok := resp.Attributes["val"]; !ok {
			return fmt.Errorf("expected 'val' in UPDATED_OLD response")
		}
		if _, ok := resp.Attributes["name"]; ok {
			return fmt.Errorf("did not expect unchanged 'name' in UPDATED_OLD response")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "10" {
			return fmt.Errorf("expected val=10 in UPDATED_OLD, got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_UpdatedNew", func() error {
		unTable := fmt.Sprintf("UN-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(unTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(unTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(unTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "un1"},
				"val":  &types.AttributeValueMemberN{Value: "10"},
				"name": &types.AttributeValueMemberS{Value: "Old"},
			},
		})

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(unTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "un1"},
			},
			UpdateExpression:          aws.String("SET #v = :new"),
			ExpressionAttributeNames:  map[string]string{"#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":new": &types.AttributeValueMemberN{Value: "20"}},
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp.Attributes == nil {
			return fmt.Errorf("expected updated new attributes")
		}
		val, ok := resp.Attributes["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "20" {
			return fmt.Errorf("expected val=20 in UPDATED_NEW, got %v", val)
		}
		if _, ok := resp.Attributes["name"]; ok {
			return fmt.Errorf("did not expect unchanged 'name' in UPDATED_NEW response")
		}
		return nil
	}))

	// UpdateItem list_append appends elements to an existing list attribute.
	results = append(results, r.RunTest("dynamodb", "UpdateItem_ListAppendNested", func() error {
		listTable := fmt.Sprintf("ListAppend-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(listTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(listTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(listTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "la1"},
				"items": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "a"}}},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(listTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "la1"},
			},
			UpdateExpression: aws.String("SET items = list_append(items, :new)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":new": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "b"}}},
			},
			ReturnValues: types.ReturnValueAllNew,
		})
		if err != nil {
			return fmt.Errorf("list_append update failed: %v", err)
		}
		items, ok := resp.Attributes["items"].(*types.AttributeValueMemberL)
		if !ok {
			return fmt.Errorf("expected list for items, got %T", resp.Attributes["items"])
		}
		if len(items.Value) != 2 {
			return fmt.Errorf("expected 2 items after list_append, got %d", len(items.Value))
		}
		first, ok := items.Value[0].(*types.AttributeValueMemberS)
		if !ok || first.Value != "a" {
			return fmt.Errorf("expected first element 'a', got %v", items.Value[0])
		}
		second, ok := items.Value[1].(*types.AttributeValueMemberS)
		if !ok || second.Value != "b" {
			return fmt.Errorf("expected second element 'b', got %v", items.Value[1])
		}
		return nil
	}))

	// list_append on a non-List attribute (SS, NS, BS, etc.) must return ValidationException.
	results = append(results, r.RunTest("dynamodb", "UpdateItem_ListAppend_TypeMismatch", func() error {
		laTable := fmt.Sprintf("ListAppendErr-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(laTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(laTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(laTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ss1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"a", "b"}},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(laTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "ss1"},
			},
			UpdateExpression: aws.String("SET tags = list_append(tags, :extra)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":extra": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberS{Value: "c"}}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for list_append on SS, got nil")
		}
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

// === CONDITION EXPRESSION TESTS ===

func (r *TestRunner) dynamoDBConditionExpressionTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "Condition_AttributeExists_True", func() error {
		ceTable := fmt.Sprintf("CE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ceTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ceTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ceTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ce1"},
				"name": &types.AttributeValueMemberS{Value: "Test"},
			},
		})

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(ceTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ce1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("attribute_exists(name)"),
			ExpressionAttributeNames:  map[string]string{"#s": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "active"}},
		})
		if err != nil {
			return fmt.Errorf("attribute_exists should pass: %v", err)
		}
		getResp, getErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(ceTable),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ce1"}},
		})
		if getErr != nil {
			return fmt.Errorf("get after update: %v", getErr)
		}
		status, ok := getResp.Item["status"].(*types.AttributeValueMemberS)
		if !ok || status.Value != "active" {
			return fmt.Errorf("expected status=active after conditional update, got %v", getResp.Item["status"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Condition_AttributeNotExists_False", func() error {
		ceTable := fmt.Sprintf("CENE-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ceTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ceTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ceTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "cene1"},
				"name": &types.AttributeValueMemberS{Value: "Test"},
			},
		})

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(ceTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "cene1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("attribute_not_exists(name)"),
			ExpressionAttributeNames:  map[string]string{"#s": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "active"}},
		})
		if err == nil {
			return fmt.Errorf("expected ConditionalCheckFailedException")
		}
		var ccf *types.ConditionalCheckFailedException
		if !errors.As(err, &ccf) {
			return fmt.Errorf("expected ConditionalCheckFailedException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Condition_BeginsWith", func() error {
		bwTable := fmt.Sprintf("BW-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(bwTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(bwTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(bwTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "bw1"},
				"name": &types.AttributeValueMemberS{Value: "HelloWorld"},
			},
		})

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(bwTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "bw1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("begins_with(name, :prefix)"),
			ExpressionAttributeNames:  map[string]string{"#s": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "matched"}, ":prefix": &types.AttributeValueMemberS{Value: "Hello"}},
		})
		if err != nil {
			return fmt.Errorf("begins_with should pass: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(bwTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "bw1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("begins_with(name, :prefix)"),
			ExpressionAttributeNames:  map[string]string{"#s": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "nope"}, ":prefix": &types.AttributeValueMemberS{Value: "XYZ"}},
		})
		if err == nil {
			return fmt.Errorf("expected ConditionalCheckFailedException for non-matching begins_with")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Condition_Contains", func() error {
		ctTable := fmt.Sprintf("CT-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ctTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ctTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(ctTable),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: "ct1"},
				"tags": &types.AttributeValueMemberSS{Value: []string{"go", "java", "python"}},
			},
		})

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(ctTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ct1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("contains(tags, :tag)"),
			ExpressionAttributeNames:  map[string]string{"#s": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "matched"}, ":tag": &types.AttributeValueMemberS{Value: "java"}},
		})
		if err != nil {
			return fmt.Errorf("contains on StringSet should pass: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(ctTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ct1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("contains(tags, :tag)"),
			ExpressionAttributeNames:  map[string]string{"#s": "status"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "nope"}, ":tag": &types.AttributeValueMemberS{Value: "rust"}},
		})
		if err == nil {
			return fmt.Errorf("expected ConditionalCheckFailedException for non-matching contains")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Condition_ComparisonOperators", func() error {
		coTable := fmt.Sprintf("CO-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(coTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(coTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(coTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "co1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		tests := []struct {
			cond string
			val  string
			pass bool
		}{
			{"#v = :x", "10", true},
			{"#v <> :x", "20", true},
			{"#v < :x", "20", true},
			{"#v <= :x", "10", true},
			{"#v > :x", "5", true},
			{"#v >= :x", "10", true},
			{"#v < :x", "5", false},
			{"#v > :x", "20", false},
		}
		for _, tc := range tests {
			_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				TableName:                 aws.String(coTable),
				Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "co1"}},
				UpdateExpression:          aws.String("SET #s = :s"),
				ConditionExpression:       aws.String(tc.cond),
				ExpressionAttributeNames:  map[string]string{"#v": "val", "#s": "status"},
				ExpressionAttributeValues: map[string]types.AttributeValue{":s": &types.AttributeValueMemberS{Value: "ok"}, ":x": &types.AttributeValueMemberN{Value: tc.val}},
			})
			if tc.pass && err != nil {
				return fmt.Errorf("condition '%s' with val '%s' should pass: %v", tc.cond, tc.val, err)
			}
			if !tc.pass && err == nil {
				return fmt.Errorf("condition '%s' with val '%s' should fail", tc.cond, tc.val)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Condition_AND_OR", func() error {
		aoTable := fmt.Sprintf("AO-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(aoTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(aoTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(aoTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "ao1"},
				"val":    &types.AttributeValueMemberN{Value: "10"},
				"active": &types.AttributeValueMemberBOOL{Value: true},
			},
		})

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(aoTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ao1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("active = :t AND #v > :x"),
			ExpressionAttributeNames:  map[string]string{"#s": "status", "#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "and-pass"}, ":t": &types.AttributeValueMemberBOOL{Value: true}, ":x": &types.AttributeValueMemberN{Value: "5"}},
		})
		if err != nil {
			return fmt.Errorf("AND condition should pass: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(aoTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ao1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("active = :f AND #v > :x"),
			ExpressionAttributeNames:  map[string]string{"#s": "status", "#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "and-fail"}, ":f": &types.AttributeValueMemberBOOL{Value: false}, ":x": &types.AttributeValueMemberN{Value: "5"}},
		})
		if err == nil {
			return fmt.Errorf("AND condition (one false) should fail")
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(aoTable),
			Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "ao1"}},
			UpdateExpression:          aws.String("SET #s = :v"),
			ConditionExpression:       aws.String("active = :f OR #v > :x"),
			ExpressionAttributeNames:  map[string]string{"#s": "status", "#v": "val"},
			ExpressionAttributeValues: map[string]types.AttributeValue{":v": &types.AttributeValueMemberS{Value: "or-pass"}, ":f": &types.AttributeValueMemberBOOL{Value: false}, ":x": &types.AttributeValueMemberN{Value: "5"}},
		})
		if err != nil {
			return fmt.Errorf("OR condition (one true) should pass: %v", err)
		}

		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBTypeTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "PutItem_GetItem_AllScalarTypes", func() error {
		allTypesTable := fmt.Sprintf("AllTypes-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(allTypesTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(allTypesTable)})

		binaryData := []byte("\x00\x01\x02\xff\xfe")
		putItem := map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: "alltypes1"},
			"str_val":  &types.AttributeValueMemberS{Value: "hello"},
			"num_val":  &types.AttributeValueMemberN{Value: "3.14"},
			"bin_val":  &types.AttributeValueMemberB{Value: binaryData},
			"bool_val": &types.AttributeValueMemberBOOL{Value: true},
			"null_val": &types.AttributeValueMemberNULL{Value: true},
		}
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(allTypesTable),
			Item:      putItem,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(allTypesTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "alltypes1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Item == nil {
			return fmt.Errorf("item is nil")
		}

		s, ok := resp.Item["str_val"].(*types.AttributeValueMemberS)
		if !ok || s.Value != "hello" {
			return fmt.Errorf("str_val mismatch: got %v", resp.Item["str_val"])
		}
		n, ok := resp.Item["num_val"].(*types.AttributeValueMemberN)
		if !ok || n.Value != "3.14" {
			return fmt.Errorf("num_val mismatch: got %v", resp.Item["num_val"])
		}
		b, ok := resp.Item["bin_val"].(*types.AttributeValueMemberB)
		if !ok || !bytes.Equal(b.Value, binaryData) {
			return fmt.Errorf("bin_val mismatch: got %v", resp.Item["bin_val"])
		}
		bo, ok := resp.Item["bool_val"].(*types.AttributeValueMemberBOOL)
		if !ok || bo.Value != true {
			return fmt.Errorf("bool_val mismatch: got %v", resp.Item["bool_val"])
		}
		nu, ok := resp.Item["null_val"].(*types.AttributeValueMemberNULL)
		if !ok || nu.Value != true {
			return fmt.Errorf("null_val mismatch: got %v", resp.Item["null_val"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_GetItem_SetTypes", func() error {
		setTable := fmt.Sprintf("SetTypes-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(setTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(setTable)})

		nsItem := map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: "setitem1"},
			"nums": &types.AttributeValueMemberNS{Value: []string{"1", "2.5", "-10"}},
			"strs": &types.AttributeValueMemberSS{Value: []string{"alpha", "beta"}},
			"bins": &types.AttributeValueMemberBS{Value: [][]byte{{0xCA, 0xFE}, {0xDE, 0xAD}}},
			"map_v": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"inner_str": &types.AttributeValueMemberS{Value: "deep"},
				"inner_num": &types.AttributeValueMemberN{Value: "42"},
			}},
			"list_v": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberN{Value: "1"},
				&types.AttributeValueMemberS{Value: "two"},
				&types.AttributeValueMemberBOOL{Value: false},
			}},
		}
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(setTable),
			Item:      nsItem,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(setTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "setitem1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Item == nil {
			return fmt.Errorf("item is nil")
		}

		ns, ok := resp.Item["nums"].(*types.AttributeValueMemberNS)
		if !ok {
			return fmt.Errorf("nums: expected NS, got %T", resp.Item["nums"])
		}
		nsSet := make(map[string]bool)
		for _, v := range ns.Value {
			nsSet[v] = true
		}
		for _, expected := range []string{"1", "2.5", "-10"} {
			if !nsSet[expected] {
				return fmt.Errorf("nums missing %q: got %v", expected, ns.Value)
			}
		}

		ss, ok := resp.Item["strs"].(*types.AttributeValueMemberSS)
		if !ok {
			return fmt.Errorf("strs: expected SS, got %T", resp.Item["strs"])
		}
		ssSet := make(map[string]bool)
		for _, v := range ss.Value {
			ssSet[v] = true
		}
		for _, expected := range []string{"alpha", "beta"} {
			if !ssSet[expected] {
				return fmt.Errorf("strs missing %q: got %v", expected, ss.Value)
			}
		}

		bs, ok := resp.Item["bins"].(*types.AttributeValueMemberBS)
		if !ok {
			return fmt.Errorf("bins: expected BS, got %T", resp.Item["bins"])
		}
		if len(bs.Value) != 2 {
			return fmt.Errorf("bins: expected 2 elements, got %d", len(bs.Value))
		}

		m, ok := resp.Item["map_v"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("map_v: expected M, got %T", resp.Item["map_v"])
		}
		innerS, ok := m.Value["inner_str"].(*types.AttributeValueMemberS)
		if !ok || innerS.Value != "deep" {
			return fmt.Errorf("map_v.inner_str mismatch: got %v", m.Value["inner_str"])
		}
		innerN, ok := m.Value["inner_num"].(*types.AttributeValueMemberN)
		if !ok || innerN.Value != "42" {
			return fmt.Errorf("map_v.inner_num mismatch: got %v", m.Value["inner_num"])
		}

		l, ok := resp.Item["list_v"].(*types.AttributeValueMemberL)
		if !ok || len(l.Value) != 3 {
			return fmt.Errorf("list_v: expected 3 elements, got %d", len(l.Value))
		}
		ln, ok := l.Value[0].(*types.AttributeValueMemberN)
		if !ok || ln.Value != "1" {
			return fmt.Errorf("list_v[0] mismatch: got %v", l.Value[0])
		}
		ls, ok := l.Value[1].(*types.AttributeValueMemberS)
		if !ok || ls.Value != "two" {
			return fmt.Errorf("list_v[1] mismatch: got %v", l.Value[1])
		}
		lb, ok := l.Value[2].(*types.AttributeValueMemberBOOL)
		if !ok || lb.Value != false {
			return fmt.Errorf("list_v[2] mismatch: got %v", l.Value[2])
		}
		return nil
	}))

	return results
}

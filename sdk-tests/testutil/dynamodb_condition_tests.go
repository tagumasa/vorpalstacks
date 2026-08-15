package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

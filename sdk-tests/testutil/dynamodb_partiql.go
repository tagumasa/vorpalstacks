package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBBasicPartiQLTests(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement (PartiQL)", func() error {
		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + tableName + "\" VALUE {'id': 'partiql1', 'name': 'PartiQL Item'}"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("ExecuteStatement response is nil")
		}
		verifyResp, verifyErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "partiql1"},
			},
		})
		if verifyErr != nil {
			return fmt.Errorf("PartiQL INSERT verification failed: %w", verifyErr)
		}
		if verifyResp.Item == nil {
			return fmt.Errorf("item not found after PartiQL INSERT")
		}
		name, ok := verifyResp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "PartiQL Item" {
			return fmt.Errorf("expected name='PartiQL Item', got %v", verifyResp.Item["name"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement (SELECT)", func() error {
		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("SELECT * FROM \"" + tableName + "\" WHERE id = 'partiql1'"),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("no items found")
		}
		name, ok := resp.Items[0]["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "PartiQL Item" {
			return fmt.Errorf("expected name='PartiQL Item', got %v", resp.Items[0]["name"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchExecuteStatement", func() error {
		resp, err := client.BatchExecuteStatement(ctx, &dynamodb.BatchExecuteStatementInput{
			Statements: []types.BatchStatementRequest{
				{
					Statement: aws.String("UPDATE \"" + tableName + "\" SET #n = :name WHERE id = 'batch1'"),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "Updated via Batch"},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Responses == nil {
			return fmt.Errorf("BatchExecuteStatement Responses is nil")
		}
		if len(resp.Responses) == 0 {
			return fmt.Errorf("expected at least one response in BatchExecuteStatement")
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBPartiQLInsertParamsTest(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	// PartiQL INSERT with parameterised values (?) substitutes placeholders correctly.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_InsertWithParams", func() error {
		paramTable := fmt.Sprintf("ParamInsert-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(paramTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(paramTable)})

		_, err = client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + paramTable + "\" VALUE {'id': ?, 'name': ?}"),
			Parameters: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "param1"},
				&types.AttributeValueMemberS{Value: "Parametrised Item"},
			},
		})
		if err != nil {
			return fmt.Errorf("parametrised INSERT failed: %v", err)
		}
		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(paramTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "param1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get after param insert: %v", err)
		}
		if getResp.Item == nil {
			return fmt.Errorf("item not found after parametrised INSERT")
		}
		name, ok := getResp.Item["name"].(*types.AttributeValueMemberS)
		if !ok || name.Value != "Parametrised Item" {
			return fmt.Errorf("expected name='Parametrised Item', got %v", getResp.Item["name"])
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBPartiQLEdgeCaseTests(ctx context.Context, client *dynamodb.Client, compTableName string) []TestResult {
	var results []TestResult

	// === PARTIQL EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_SelectWhere", func() error {
		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("SELECT * FROM \"" + compTableName + "\" WHERE pk = 'user1' AND sk = 'order2'"),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) != 1 {
			return fmt.Errorf("expected 1 item, got %d", len(resp.Items))
		}
		amt, ok := resp.Items[0]["amount"].(*types.AttributeValueMemberN)
		if !ok || amt.Value != "200" {
			return fmt.Errorf("expected amount=200, got %v", resp.Items[0]["amount"])
		}
		status, ok := resp.Items[0]["status"].(*types.AttributeValueMemberS)
		if !ok || status.Value != "pending" {
			return fmt.Errorf("expected status=pending, got %v", resp.Items[0]["status"])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_Update", func() error {
		puTable := fmt.Sprintf("PQUpd-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(puTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(puTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(puTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "pu1"},
				"val": &types.AttributeValueMemberN{Value: "10"},
			},
		})

		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("UPDATE \"" + puTable + "\" SET val = 20 WHERE id = 'pu1'"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(puTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "pu1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		val, ok := getResp.Item["val"].(*types.AttributeValueMemberN)
		if !ok || val.Value != "20" {
			return fmt.Errorf("expected val=20 after PartiQL UPDATE, got %v", val)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_Delete", func() error {
		pdTable := fmt.Sprintf("PQDel-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(pdTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pdTable)})

		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(pdTable),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "pd1"},
				"val": &types.AttributeValueMemberN{Value: "99"},
			},
		})

		_, err = client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("DELETE FROM \"" + pdTable + "\" WHERE id = 'pd1'"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}

		getResp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(pdTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "pd1"},
			},
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if len(getResp.Item) != 0 {
			return fmt.Errorf("item should be deleted after PartiQL DELETE")
		}
		return nil
	}))

	// SELECT with column projection returns only the specified columns, not all attributes.
	results = append(results, r.RunTest("dynamodb", "ExecuteStatement_SelectProjection", func() error {
		projTable := fmt.Sprintf("PQProj-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(projTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(projTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(projTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "proj1"},
				"name":   &types.AttributeValueMemberS{Value: "Alice"},
				"age":    &types.AttributeValueMemberN{Value: "30"},
				"active": &types.AttributeValueMemberBOOL{Value: true},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("SELECT name, age FROM \"" + projTable + "\" WHERE id = 'proj1'"),
		})
		if err != nil {
			return fmt.Errorf("select projection: %v", err)
		}
		if len(resp.Items) != 1 {
			return fmt.Errorf("expected 1 item, got %d", len(resp.Items))
		}
		item := resp.Items[0]
		if _, ok := item["name"]; !ok {
			return fmt.Errorf("missing projected column 'name'")
		}
		if _, ok := item["age"]; !ok {
			return fmt.Errorf("missing projected column 'age'")
		}
		if _, ok := item["active"]; ok {
			return fmt.Errorf("unprojected column 'active' should not be present")
		}
		if _, ok := item["id"]; ok {
			return fmt.Errorf("unprojected column 'id' should not be present")
		}
		return nil
	}))

	return results
}

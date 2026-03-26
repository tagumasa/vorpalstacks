package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunDynamoDBTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := dynamodb.NewFromConfig(cfg)
	ctx := context.Background()

	tableName := fmt.Sprintf("TestTable-%d", time.Now().UnixNano())
	tableARN := fmt.Sprintf("arn:aws:dynamodb:%s:000000000000:table/%s", r.region, tableName)

	results = append(results, r.RunTest("dynamodb", "CreateTable", func() error {
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeTable", func() error {
		resp, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.Table == nil {
			return fmt.Errorf("table not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ListTables", func() error {
		resp, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
		if err != nil {
			return err
		}
		if resp.TableNames == nil {
			return fmt.Errorf("table names is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "test1"},
				"name":  &types.AttributeValueMemberS{Value: "Test Item"},
				"count": &types.AttributeValueMemberN{Value: "42"},
			},
		})
		return err
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
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			KeyConditionExpression: aws.String("id = :id"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":id": &types.AttributeValueMemberS{Value: "test1"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count == 0 {
			return fmt.Errorf("no items found")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.Count == 0 {
			return fmt.Errorf("no items found")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem", func() error {
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
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
		return err
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
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteItem", func() error {
		_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "test1"},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "TagResource", func() error {
		_, err := client.TagResource(ctx, &dynamodb.TagResourceInput{
			ResourceArn: aws.String(tableARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("Test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "ListTagsOfResource", func() error {
		resp, err := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
			ResourceArn: aws.String(tableARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) == 0 {
			return fmt.Errorf("no tags found")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &dynamodb.UntagResourceInput{
			ResourceArn: aws.String(tableARN),
			TagKeys:     []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive", func() error {
		_, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableName),
			TimeToLiveSpecification: &types.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeTimeToLive", func() error {
		resp, err := client.DescribeTimeToLive(ctx, &dynamodb.DescribeTimeToLiveInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.TimeToLiveDescription == nil {
			return fmt.Errorf("TTL description not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateBackup", func() error {
		backupName := fmt.Sprintf("TestBackup-%d", time.Now().UnixNano())
		_, err := client.CreateBackup(ctx, &dynamodb.CreateBackupInput{
			TableName:  aws.String(tableName),
			BackupName: aws.String(backupName),
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "ListBackups", func() error {
		resp, err := client.ListBackups(ctx, &dynamodb.ListBackupsInput{})
		if err != nil {
			return err
		}
		if resp.BackupSummaries == nil {
			return fmt.Errorf("backup summaries is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeContinuousBackups", func() error {
		resp, err := client.DescribeContinuousBackups(ctx, &dynamodb.DescribeContinuousBackupsInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.ContinuousBackupsDescription == nil {
			return fmt.Errorf("continuous backups description not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateContinuousBackups", func() error {
		_, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
			TableName: aws.String(tableName),
			PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
				PointInTimeRecoveryEnabled: aws.Bool(true),
			},
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "ExecuteStatement (PartiQL)", func() error {
		_, err := client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
			Statement: aws.String("INSERT INTO \"" + tableName + "\" VALUE {'id': 'partiql1', 'name': 'PartiQL Item'}"),
		})
		return err
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
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems", func() error {
		_, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
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
		return err
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
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchExecuteStatement", func() error {
		_, err := client.BatchExecuteStatement(ctx, &dynamodb.BatchExecuteStatementInput{
			Statements: []types.BatchStatementRequest{
				{
					Statement: aws.String("UPDATE \"" + tableName + "\" SET #n = :name WHERE id = 'batch1'"),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "Updated via Batch"},
					},
				},
			},
		})
		return err
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
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTable", func() error {
		_, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName: aws.String(tableName),
		})
		return err
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteTable", func() error {
		_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("dynamodb", "GetItem_NonExistentTable", func() error {
		_, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String("NoSuchTable_xyz"),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_NonExistentTable", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("NoSuchTable_xyz"),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_NonExistentTable", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String("NoSuchTable_xyz"),
			KeyConditionExpression: aws.String("id = :id"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_NonExistentTable", func() error {
		_, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeTable_NonExistentTable", func() error {
		_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteTable_NonExistentTable", func() error {
		_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_ConditionalCheckFail", func() error {
		errTable := fmt.Sprintf("CondTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(errTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(errTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(errTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "cond1"},
				"status": &types.AttributeValueMemberS{Value: "active"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(errTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "cond1"},
			},
			UpdateExpression:    aws.String("SET #s = :val"),
			ConditionExpression: aws.String("#s = :expected"),
			ExpressionAttributeNames: map[string]string{
				"#s": "status",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":val":      &types.AttributeValueMemberS{Value: "inactive"},
				":expected": &types.AttributeValueMemberS{Value: "deleted"},
			},
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

	results = append(results, r.RunTest("dynamodb", "GetItem_NonExistentKey", func() error {
		errTable := fmt.Sprintf("GetItemErr-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(errTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(errTable)})

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(errTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "nonexistent"},
			},
		})
		if err != nil {
			return fmt.Errorf("GetItem on non-existent key should not error, got: %v", err)
		}
		if len(resp.Item) != 0 {
			return fmt.Errorf("expected empty item for non-existent key, got %d attributes", len(resp.Item))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteItem_ConditionalCheckFail", func() error {
		errTable := fmt.Sprintf("DelCondTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(errTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(errTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(errTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "del1"},
				"status": &types.AttributeValueMemberS{Value: "active"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		_, err = client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(errTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "del1"},
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
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateTable_DuplicateName", func() error {
		dupTable := fmt.Sprintf("DupTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(dupTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(dupTable)})

		_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(dupTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate table name")
		}
		var ri *types.ResourceInUseException
		if !errors.As(err, &ri) {
			return fmt.Errorf("expected ResourceInUseException, got: %T: %v", err, err)
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

	return results
}

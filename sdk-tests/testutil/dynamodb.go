package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		resp, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return err
		}
		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		return nil
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
		resp, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableName),
			TimeToLiveSpecification: &types.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		if resp.TimeToLiveSpecification == nil {
			return fmt.Errorf("TimeToLiveSpecification is nil")
		}
		return nil
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
		resp, err := client.CreateBackup(ctx, &dynamodb.CreateBackupInput{
			TableName:  aws.String(tableName),
			BackupName: aws.String(backupName),
		})
		if err != nil {
			return err
		}
		if resp.BackupDetails == nil {
			return fmt.Errorf("BackupDetails is nil")
		}
		client.DeleteBackup(ctx, &dynamodb.DeleteBackupInput{
			BackupArn: resp.BackupDetails.BackupArn,
		})
		return nil
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
		resp, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
			TableName: aws.String(tableName),
			PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
				PointInTimeRecoveryEnabled: aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		if resp.ContinuousBackupsDescription == nil {
			return fmt.Errorf("ContinuousBackupsDescription is nil")
		}
		return nil
	}))

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
		return nil
	}))

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
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTable", func() error {
		resp, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		return nil
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

	// === COMPOSITE KEY TABLE SETUP ===

	compTableName := fmt.Sprintf("CompTable-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("dynamodb", "CreateTable_CompositeKey", func() error {
		resp, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(compTableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return err
		}
		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.KeySchema == nil || len(resp.TableDescription.KeySchema) != 2 {
			return fmt.Errorf("expected 2 key schema elements")
		}
		return nil
	}))

	// Seed composite key data
	for _, item := range []map[string]types.AttributeValue{
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "meta"}, "name": &types.AttributeValueMemberS{Value: "Alice"}, "age": &types.AttributeValueMemberN{Value: "30"}, "active": &types.AttributeValueMemberBOOL{Value: true}},
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "order1"}, "amount": &types.AttributeValueMemberN{Value: "100"}, "status": &types.AttributeValueMemberS{Value: "shipped"}},
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "order2"}, "amount": &types.AttributeValueMemberN{Value: "200"}, "status": &types.AttributeValueMemberS{Value: "pending"}},
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "order3"}, "amount": &types.AttributeValueMemberN{Value: "50"}, "status": &types.AttributeValueMemberS{Value: "delivered"}},
		{"pk": &types.AttributeValueMemberS{Value: "user2"}, "sk": &types.AttributeValueMemberS{Value: "meta"}, "name": &types.AttributeValueMemberS{Value: "Bob"}, "age": &types.AttributeValueMemberN{Value: "25"}, "active": &types.AttributeValueMemberBOOL{Value: false}},
		{"pk": &types.AttributeValueMemberS{Value: "user2"}, "sk": &types.AttributeValueMemberS{Value: "order1"}, "amount": &types.AttributeValueMemberN{Value: "300"}, "status": &types.AttributeValueMemberS{Value: "shipped"}},
	} {
		client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(compTableName),
			Item:      item,
		})
	}

	// === PUT ITEM EDGE CASES ===

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
		return nil
	}))

	// === GET ITEM EDGE CASES ===

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
		return nil
	}))

	// === DELETE ITEM EDGE CASES ===

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

	// === UPDATE ITEM EDGE CASES ===

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

	// === CONDITION EXPRESSION TESTS ===

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

	// === QUERY EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "Query_CompositeKey", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "user1"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 4 {
			return fmt.Errorf("expected 4 items for user1, got %d", resp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_SortKeyCondition", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk AND sk = :sk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "user1"},
				":sk": &types.AttributeValueMemberS{Value: "order2"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 1 {
			return fmt.Errorf("expected 1 item, got %d", resp.Count)
		}
		amt, ok := resp.Items[0]["amount"].(*types.AttributeValueMemberN)
		if !ok || amt.Value != "200" {
			return fmt.Errorf("expected amount=200, got %v", amt)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_SortKeyBeginsWith", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":     &types.AttributeValueMemberS{Value: "user1"},
				":prefix": &types.AttributeValueMemberS{Value: "order"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 3 {
			return fmt.Errorf("expected 3 order items, got %d", resp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_SortKeyBetween", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk AND sk BETWEEN :low AND :high"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":   &types.AttributeValueMemberS{Value: "user1"},
				":low":  &types.AttributeValueMemberS{Value: "order1"},
				":high": &types.AttributeValueMemberS{Value: "order3"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 3 {
			return fmt.Errorf("expected 3 items in BETWEEN, got %d", resp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_ScanIndexForward", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "user1"},
			},
			ScanIndexForward: aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if resp.Count != 4 {
			return fmt.Errorf("expected 4 items, got %d", resp.Count)
		}
		firstSK, ok := resp.Items[0]["sk"].(*types.AttributeValueMemberS)
		if !ok || firstSK.Value != "order3" {
			return fmt.Errorf("expected first item sk=order3 in descending order, got %v", firstSK)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_FilterExpression", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
			FilterExpression:       aws.String("#s = :status"),
			ExpressionAttributeNames: map[string]string{
				"#s": "status",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":     &types.AttributeValueMemberS{Value: "user1"},
				":prefix": &types.AttributeValueMemberS{Value: "order"},
				":status": &types.AttributeValueMemberS{Value: "shipped"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 1 {
			return fmt.Errorf("expected 1 shipped order, got %d", resp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_Limit", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "user1"},
			},
			Limit: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if resp.Count > 2 {
			return fmt.Errorf("expected at most 2 items with Limit=2, got %d", resp.Count)
		}
		if resp.LastEvaluatedKey == nil {
			return fmt.Errorf("expected LastEvaluatedKey when Limit < total items")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_ProjectionExpression", func() error {
		resp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(compTableName),
			KeyConditionExpression: aws.String("pk = :pk AND sk = :sk"),
			ProjectionExpression:   aws.String("pk, sk, amount"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "user1"},
				":sk": &types.AttributeValueMemberS{Value: "order1"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 1 {
			return fmt.Errorf("expected 1 item, got %d", resp.Count)
		}
		if len(resp.Items[0]) != 3 {
			return fmt.Errorf("expected 3 projected attributes, got %d", len(resp.Items[0]))
		}
		return nil
	}))

	// === SCAN EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "Scan_FilterExpression", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        aws.String(compTableName),
			FilterExpression: aws.String("#s = :status"),
			ExpressionAttributeNames: map[string]string{
				"#s": "status",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":status": &types.AttributeValueMemberS{Value: "shipped"},
			},
		})
		if err != nil {
			return err
		}
		if resp.Count != 2 {
			return fmt.Errorf("expected 2 shipped items, got %d", resp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_ProjectionExpression", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:            aws.String(compTableName),
			ProjectionExpression: aws.String("pk, name"),
		})
		if err != nil {
			return err
		}
		for _, item := range resp.Items {
			if _, ok := item["pk"]; !ok {
				return fmt.Errorf("expected 'pk' in projected item")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_Limit", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(compTableName),
			Limit:     aws.Int32(3),
		})
		if err != nil {
			return err
		}
		if resp.Count > 3 {
			return fmt.Errorf("expected at most 3 items with Limit=3, got %d", resp.Count)
		}
		return nil
	}))

	// === BATCH OPERATION EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_DeleteRequest", func() error {
		bwTable := fmt.Sprintf("BWDel-%d", time.Now().UnixNano())
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
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(bgTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(bgTable)})

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
		return nil
	}))

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
		return nil
	}))

	// === TABLE OPERATION EDGE CASES ===

	results = append(results, r.RunTest("dynamodb", "CreateTable_Validation_NoKeySchema", func() error {
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String("BadTable"),
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for missing KeySchema")
		}
		if err == nil {
			return fmt.Errorf("expected error for missing KeySchema")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteTable_Protected", func() error {
		dpTable := fmt.Sprintf("DP-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(dpTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode:               types.BillingModePayPerRequest,
			DeletionProtectionEnabled: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer func() {
			client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
				TableName:                 aws.String(dpTable),
				DeletionProtectionEnabled: aws.Bool(false),
			})
			client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(dpTable)})
		}()

		_, err = client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(dpTable),
		})
		if err == nil {
			return fmt.Errorf("expected ResourceInUseException for protected table")
		}
		var ri *types.ResourceInUseException
		if !errors.As(err, &ri) {
			return fmt.Errorf("expected ResourceInUseException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateTable_GSI", func() error {
		gsiTable := fmt.Sprintf("GSITable-%d", time.Now().UnixNano())
		resp, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(gsiTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("gsi_pk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
			},
			BillingMode: types.BillingModePayPerRequest,
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("gsi1"),
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String("gsi_pk"), KeyType: types.KeyTypeHash},
						{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create with GSI: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(gsiTable)})

		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.GlobalSecondaryIndexes == nil || len(resp.TableDescription.GlobalSecondaryIndexes) != 1 {
			return fmt.Errorf("expected 1 GSI in description")
		}
		if resp.TableDescription.GlobalSecondaryIndexes[0].IndexName == nil || *resp.TableDescription.GlobalSecondaryIndexes[0].IndexName != "gsi1" {
			return fmt.Errorf("expected GSI name 'gsi1'")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateTable_StreamSpec", func() error {
		stTable := fmt.Sprintf("StreamTable-%d", time.Now().UnixNano())
		resp, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(stTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
			StreamSpecification: &types.StreamSpecification{
				StreamEnabled:  aws.Bool(true),
				StreamViewType: types.StreamViewTypeNewImage,
			},
		})
		if err != nil {
			return fmt.Errorf("create with stream: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(stTable)})

		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.StreamSpecification == nil || resp.TableDescription.StreamSpecification.StreamEnabled == nil || !*resp.TableDescription.StreamSpecification.StreamEnabled {
			return fmt.Errorf("expected StreamEnabled=true")
		}
		if resp.TableDescription.LatestStreamArn == nil {
			return fmt.Errorf("expected LatestStreamArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTable_EnableSSE", func() error {
		sseTable := fmt.Sprintf("SSETable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(sseTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(sseTable)})

		resp, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName: aws.String(sseTable),
			SSESpecification: &types.SSESpecification{
				Enabled: aws.Bool(true),
				SSEType: types.SSETypeAes256,
			},
		})
		if err != nil {
			return fmt.Errorf("update SSE: %v", err)
		}
		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.SSEDescription == nil {
			return fmt.Errorf("expected SSEDescription")
		}
		if resp.TableDescription.SSEDescription.Status != types.SSEStatusEnabled {
			return fmt.Errorf("expected SSEStatus=ENABLED, got %v", resp.TableDescription.SSEDescription.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTable_AddGSI", func() error {
		agTable := fmt.Sprintf("AddGSI-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(agTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("gsi_pk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(agTable)})

		resp, err := client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName: aws.String(agTable),
			GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
				{
					Create: &types.CreateGlobalSecondaryIndexAction{
						IndexName: aws.String("new_gsi"),
						KeySchema: []types.KeySchemaElement{
							{AttributeName: aws.String("gsi_pk"), KeyType: types.KeyTypeHash},
						},
						Projection: &types.Projection{
							ProjectionType: types.ProjectionTypeAll,
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("update add GSI: %v", err)
		}
		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.GlobalSecondaryIndexes == nil || len(resp.TableDescription.GlobalSecondaryIndexes) != 1 {
			return fmt.Errorf("expected 1 GSI after update")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive_Disable", func() error {
		ttlTable := fmt.Sprintf("TTLDis-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(ttlTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(ttlTable)})

		client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(ttlTable),
			TimeToLiveSpecification: &types.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		})

		resp, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(ttlTable),
			TimeToLiveSpecification: &types.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(false),
			},
		})
		if err != nil {
			return fmt.Errorf("disable TTL: %v", err)
		}
		if resp.TimeToLiveSpecification == nil {
			return fmt.Errorf("TimeToLiveSpecification is nil")
		}
		return nil
	}))

	// Clean up composite key table
	client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(compTableName)})

	// === PAGINATION TESTS ===

	results = append(results, r.RunTest("dynamodb", "Query_Limit_Pagination", func() error {
		pagTableName := fmt.Sprintf("PagTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(pagTableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}

		for i := 0; i < 5; i++ {
			_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(pagTableName),
				Item: map[string]types.AttributeValue{
					"pk":   &types.AttributeValueMemberS{Value: fmt.Sprintf("item-%d", i)},
					"data": &types.AttributeValueMemberS{Value: "pagination"},
				},
			})
			if err != nil {
				client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagTableName)})
				return fmt.Errorf("put item %d: %v", i, err)
			}
		}

		var allItems []string
		var exclusiveStartKey map[string]types.AttributeValue
		pageCount := 0
		for {
			resp, err := client.Query(ctx, &dynamodb.QueryInput{
				TableName:              aws.String(pagTableName),
				KeyConditionExpression: aws.String("pk = :pk"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":pk": &types.AttributeValueMemberS{Value: "item-0"},
				},
				Limit:             aws.Int32(2),
				ExclusiveStartKey: exclusiveStartKey,
			})
			if err != nil {
				client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagTableName)})
				return fmt.Errorf("query page: %v", err)
			}
			pageCount++
			for _, item := range resp.Items {
				if pk, ok := item["pk"].(*types.AttributeValueMemberS); ok {
					allItems = append(allItems, pk.Value)
				}
			}
			if resp.LastEvaluatedKey != nil {
				exclusiveStartKey = resp.LastEvaluatedKey
			} else {
				break
			}
		}

		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagTableName)})
		if len(allItems) != 1 {
			return fmt.Errorf("expected 1 item for pk=item-0, got %d", len(allItems))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ListGlobalTables_Pagination", func() error {
		pagGT1 := fmt.Sprintf("PagGT-%d-1", time.Now().UnixNano())
		pagGT2 := fmt.Sprintf("PagGT-%d-2", time.Now().UnixNano())
		_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(pagGT1),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String("us-east-1")},
			},
		})
		if err != nil {
			return fmt.Errorf("create global table %s: %v", pagGT1, err)
		}
		_, err = client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
			GlobalTableName: aws.String(pagGT2),
			ReplicationGroup: []types.Replica{
				{RegionName: aws.String("us-east-1")},
			},
		})
		if err != nil {
			client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagGT1)})
			return fmt.Errorf("create global table %s: %v", pagGT2, err)
		}

		found := map[string]bool{pagGT1: false, pagGT2: false}
		var exclusiveStartName *string
		pageCount := 0
		for {
			resp, err := client.ListGlobalTables(ctx, &dynamodb.ListGlobalTablesInput{
				Limit:                         aws.Int32(1),
				ExclusiveStartGlobalTableName: exclusiveStartName,
			})
			if err != nil {
				client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagGT1)})
				client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagGT2)})
				return fmt.Errorf("list global tables page: %v", err)
			}
			pageCount++
			for _, gt := range resp.GlobalTables {
				if gt.GlobalTableName != nil {
					if _, ok := found[*gt.GlobalTableName]; ok {
						found[*gt.GlobalTableName] = true
					}
				}
			}
			if resp.LastEvaluatedGlobalTableName != nil && *resp.LastEvaluatedGlobalTableName != "" {
				exclusiveStartName = resp.LastEvaluatedGlobalTableName
			} else {
				break
			}
		}

		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagGT1)})
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pagGT2)})
		for name, f := range found {
			if !f {
				return fmt.Errorf("created global table %q not found in ListGlobalTables", name)
			}
		}
		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages, got %d", pageCount)
		}
		return nil
	}))

	return results
}

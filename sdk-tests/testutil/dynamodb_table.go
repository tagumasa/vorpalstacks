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

func (r *TestRunner) dynamoDBCreateTableTests(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

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
		if resp.TableDescription.TableName == nil || *resp.TableDescription.TableName != tableName {
			return fmt.Errorf("TableName mismatch: got %v", resp.TableDescription.TableName)
		}
		if len(resp.TableDescription.KeySchema) != 1 || *resp.TableDescription.KeySchema[0].AttributeName != "id" {
			return fmt.Errorf("KeySchema mismatch")
		}
		if len(resp.TableDescription.AttributeDefinitions) != 1 || *resp.TableDescription.AttributeDefinitions[0].AttributeName != "id" {
			return fmt.Errorf("AttributeDefinitions mismatch")
		}
		if resp.TableDescription.TableArn == nil {
			return fmt.Errorf("TableArn is nil")
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
		if *resp.Table.TableName != tableName {
			return fmt.Errorf("TableName mismatch: got %v, want %s", resp.Table.TableName, tableName)
		}
		if resp.Table.TableStatus != types.TableStatusActive {
			return fmt.Errorf("expected Active, got %s", resp.Table.TableStatus)
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
		found := false
		for _, name := range resp.TableNames {
			if name == tableName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created table %s not found in ListTables response", tableName)
		}
		return nil
	}))

	// ListTables with Limit paginates via ExclusiveStartTableName/LastEvaluatedTableName.
	results = append(results, r.RunTest("dynamodb", "ListTables_Pagination", func() error {
		var tableNames []string
		for i := 0; i < 5; i++ {
			pn := fmt.Sprintf("PagList-%d-%d", time.Now().UnixNano(), i)
			_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
				TableName: aws.String(pn),
				AttributeDefinitions: []types.AttributeDefinition{
					{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
				},
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				},
				BillingMode: types.BillingModePayPerRequest,
			})
			if err != nil {
				return fmt.Errorf("create %s: %v", pn, err)
			}
			tableNames = append(tableNames, pn)
		}
		for _, n := range tableNames {
			defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(n)})
		}

		var collected []string
		var exclusiveStart *string
		for {
			resp, err := client.ListTables(ctx, &dynamodb.ListTablesInput{
				Limit:                   aws.Int32(2),
				ExclusiveStartTableName: exclusiveStart,
			})
			if err != nil {
				return fmt.Errorf("ListTables page: %v", err)
			}
			collected = append(collected, resp.TableNames...)
			if resp.LastEvaluatedTableName == nil {
				break
			}
			exclusiveStart = resp.LastEvaluatedTableName
		}
		for _, tn := range tableNames {
			found := false
			for _, c := range collected {
				if c == tn {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("table %s not found in paginated ListTables results (got %d tables)", tn, len(collected))
			}
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBDeleteItemTest(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "DeleteItem", func() error {
		_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "test1"},
			},
		})
		if err != nil {
			return err
		}
		verifyResp, verifyErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "test1"}},
		})
		if verifyErr != nil {
			return fmt.Errorf("DeleteItem verification GetItem failed: %w", verifyErr)
		}
		if verifyResp.Item != nil && len(verifyResp.Item) > 0 {
			return fmt.Errorf("DeleteItem verification: item still exists after deletion")
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBUpdateTableTests(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

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
		if resp.TableDescription.TableName == nil || *resp.TableDescription.TableName != tableName {
			return fmt.Errorf("TableName mismatch: expected %s, got %v", tableName, resp.TableDescription.TableName)
		}
		return nil
	}))

	// UpdateTable must not mutate the original table description (deep copy required).
	results = append(results, r.RunTest("dynamodb", "UpdateTable_DoesNotCorruptOriginal", func() error {
		descResp, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		originalName := *descResp.Table.TableName
		_, err = client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		descResp2, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if descResp2.Table.TableName == nil || *descResp2.Table.TableName != originalName {
			return fmt.Errorf("UpdateTable corrupted table name: expected %s, got %v", originalName, descResp2.Table.TableName)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBDeleteTableTest(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "DeleteTable", func() error {
		resp, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.TableName == nil || *resp.TableDescription.TableName != tableName {
			return fmt.Errorf("TableName mismatch: expected %s, got %v", tableName, resp.TableDescription.TableName)
		}
		if resp.TableDescription.TableStatus != types.TableStatusArchived && resp.TableDescription.TableStatus != types.TableStatusArchiving && resp.TableDescription.TableStatus != types.TableStatusDeleting {
			return fmt.Errorf("expected ARCHIVED, ARCHIVING or DELETING status, got %v", resp.TableDescription.TableStatus)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBCompositeKeySetup(ctx context.Context, client *dynamodb.Client, compTableName string) []TestResult {
	var results []TestResult

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

	for _, item := range []map[string]types.AttributeValue{
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "meta"}, "name": &types.AttributeValueMemberS{Value: "Alice"}, "age": &types.AttributeValueMemberN{Value: "30"}, "active": &types.AttributeValueMemberBOOL{Value: true}},
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "order1"}, "amount": &types.AttributeValueMemberN{Value: "100"}, "status": &types.AttributeValueMemberS{Value: "shipped"}},
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "order2"}, "amount": &types.AttributeValueMemberN{Value: "200"}, "status": &types.AttributeValueMemberS{Value: "pending"}},
		{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "order3"}, "amount": &types.AttributeValueMemberN{Value: "50"}, "status": &types.AttributeValueMemberS{Value: "delivered"}},
		{"pk": &types.AttributeValueMemberS{Value: "user2"}, "sk": &types.AttributeValueMemberS{Value: "meta"}, "name": &types.AttributeValueMemberS{Value: "Bob"}, "age": &types.AttributeValueMemberN{Value: "25"}, "active": &types.AttributeValueMemberBOOL{Value: false}},
		{"pk": &types.AttributeValueMemberS{Value: "user2"}, "sk": &types.AttributeValueMemberS{Value: "order1"}, "amount": &types.AttributeValueMemberN{Value: "300"}, "status": &types.AttributeValueMemberS{Value: "shipped"}},
	} {
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(compTableName),
			Item:      item,
		}); err != nil {
			return append(results, TestResult{Service: "dynamodb", TestName: "CompositeKeySetup", Status: "FAIL", Error: fmt.Sprintf("put composite item: %v", err)})
		}
	}

	return results
}

func (r *TestRunner) dynamoDBDuplicateTableTest(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

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

	return results
}

func (r *TestRunner) dynamoDBTableEdgeCaseTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "CreateTable_Validation_NoKeySchema", func() error {
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String("BadTable"),
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for missing KeySchema")
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
		if resp.TimeToLiveSpecification.Enabled == nil || *resp.TimeToLiveSpecification.Enabled {
			return fmt.Errorf("expected Enabled=false after disabling TTL")
		}
		if resp.TimeToLiveSpecification.AttributeName == nil || *resp.TimeToLiveSpecification.AttributeName != "ttl" {
			return fmt.Errorf("expected AttributeName=ttl, got %v", resp.TimeToLiveSpecification.AttributeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateTable_SSESpecification", func() error {
		sseCreateTable := fmt.Sprintf("SSECreate-%d", time.Now().UnixNano())
		resp, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(sseCreateTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
			SSESpecification: &types.SSESpecification{
				Enabled: aws.Bool(true),
				SSEType: types.SSETypeAes256,
			},
		})
		if err != nil {
			return fmt.Errorf("create with SSE: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(sseCreateTable)})

		if resp.TableDescription == nil {
			return fmt.Errorf("TableDescription is nil")
		}
		if resp.TableDescription.SSEDescription == nil {
			return fmt.Errorf("expected SSEDescription in CreateTable response with SSESpecification")
		}
		if resp.TableDescription.SSEDescription.Status != types.SSEStatusEnabled {
			return fmt.Errorf("expected SSEStatus=ENABLED, got %v", resp.TableDescription.SSEDescription.Status)
		}
		if resp.TableDescription.SSEDescription.SSEType != types.SSETypeAes256 {
			return fmt.Errorf("expected SSEType=AES256, got %v", resp.TableDescription.SSEDescription.SSEType)
		}
		return nil
	}))

	return results
}

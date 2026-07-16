package testutil

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBBasicQueryScan(ctx context.Context, client *dynamodb.Client, tableName string) []TestResult {
	var results []TestResult

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
		if resp.Count != 1 {
			return fmt.Errorf("expected 1 item, got %d", resp.Count)
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
		if resp.Count < 1 {
			return fmt.Errorf("expected at least 1 item, got %d", resp.Count)
		}
		if resp.ScannedCount < resp.Count {
			return fmt.Errorf("ScannedCount (%d) should be >= Count (%d)", resp.ScannedCount, resp.Count)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBConditionNotCaseTest(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	// NOT operator in condition expressions is case-insensitive (not, NOT, Not all valid).
	results = append(results, r.RunTest("dynamodb", "Query_ConditionNot_CaseInsensitive", func() error {
		notTable := fmt.Sprintf("NotCase-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(notTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(notTable)})

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(notTable),
			Item: map[string]types.AttributeValue{
				"id":    &types.AttributeValueMemberS{Value: "not1"},
				"val":   &types.AttributeValueMemberS{Value: "active"},
				"count": &types.AttributeValueMemberN{Value: "5"},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		scanResp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        aws.String(notTable),
			FilterExpression: aws.String("Not(val = :inactive)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":inactive": &types.AttributeValueMemberS{Value: "inactive"},
			},
		})
		if err != nil {
			return fmt.Errorf("scan with lowercase 'Not' failed: %v", err)
		}
		if scanResp.Count != 1 {
			return fmt.Errorf("expected 1 item with NOT filter, got %d", scanResp.Count)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBCompositeQueryScanTests(ctx context.Context, client *dynamodb.Client, compTableName string) []TestResult {
	var results []TestResult

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
		if resp.ScannedCount > 3 {
			return fmt.Errorf("expected at most 3 scanned items with Limit=3, got %d", resp.ScannedCount)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBNestedTypeTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "Query_Scan_NestedTypes", func() error {
		nestedTable := fmt.Sprintf("NestedTypes-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(nestedTable),
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
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(nestedTable)})

		mapItem := map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "nested1"},
			"config": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"ttl": &types.AttributeValueMemberN{Value: "3600"},
				"flags": &types.AttributeValueMemberL{Value: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "enabled"},
					&types.AttributeValueMemberBOOL{Value: true},
				}},
				"metadata": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"version": &types.AttributeValueMemberS{Value: "1.0"},
					"hash":    &types.AttributeValueMemberB{Value: []byte{0xAB, 0xCD}},
				}},
			}},
		}
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(nestedTable),
			Item:      mapItem,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		queryResp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(nestedTable),
			KeyConditionExpression: aws.String("pk = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "nested1"},
			},
		})
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}
		if queryResp.Count != 1 {
			return fmt.Errorf("query count: expected 1, got %d", queryResp.Count)
		}
		item := queryResp.Items[0]
		config, ok := item["config"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("config: expected M, got %T", item["config"])
		}
		ttl, ok := config.Value["ttl"].(*types.AttributeValueMemberN)
		if !ok || ttl.Value != "3600" {
			return fmt.Errorf("config.ttl mismatch: got %v", config.Value["ttl"])
		}
		meta, ok := config.Value["metadata"].(*types.AttributeValueMemberM)
		if !ok {
			return fmt.Errorf("config.metadata: expected M, got %T", config.Value["metadata"])
		}
		hashVal, ok := meta.Value["hash"].(*types.AttributeValueMemberB)
		if !ok || !bytes.Equal(hashVal.Value, []byte{0xAB, 0xCD}) {
			return fmt.Errorf("config.metadata.hash mismatch: got %v", meta.Value["hash"])
		}

		scanResp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(nestedTable),
		})
		if err != nil {
			return fmt.Errorf("scan: %v", err)
		}
		if scanResp.Count < 1 {
			return fmt.Errorf("scan count: expected >= 1, got %d", scanResp.Count)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBPaginationTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

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

func (r *TestRunner) dynamoDBIndexQueryTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "GSI_Query_Scan", func() error {
		gsiTestTable := fmt.Sprintf("GSIQuery-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(gsiTestTable),
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
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(gsiTestTable)})

		for _, item := range []map[string]types.AttributeValue{
			{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "a"}, "gsi_pk": &types.AttributeValueMemberS{Value: "cat_a"}, "data": &types.AttributeValueMemberS{Value: "first"}},
			{"pk": &types.AttributeValueMemberS{Value: "user1"}, "sk": &types.AttributeValueMemberS{Value: "b"}, "gsi_pk": &types.AttributeValueMemberS{Value: "cat_b"}, "data": &types.AttributeValueMemberS{Value: "second"}},
			{"pk": &types.AttributeValueMemberS{Value: "user2"}, "sk": &types.AttributeValueMemberS{Value: "c"}, "gsi_pk": &types.AttributeValueMemberS{Value: "cat_a"}, "data": &types.AttributeValueMemberS{Value: "third"}},
		} {
			if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(gsiTestTable),
				Item:      item,
			}); err != nil {
				return fmt.Errorf("put item: %v", err)
			}
		}

		gsiQueryResp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(gsiTestTable),
			IndexName:              aws.String("gsi1"),
			KeyConditionExpression: aws.String("gsi_pk = :gpk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":gpk": &types.AttributeValueMemberS{Value: "cat_a"},
			},
		})
		if err != nil {
			return fmt.Errorf("GSI query: %v", err)
		}
		if gsiQueryResp.Count != 2 {
			return fmt.Errorf("GSI query count: expected 2, got %d", gsiQueryResp.Count)
		}
		dataSet := make(map[string]bool)
		for _, item := range gsiQueryResp.Items {
			d, ok := item["data"].(*types.AttributeValueMemberS)
			if !ok {
				return fmt.Errorf("GSI query item: data is not S")
			}
			dataSet[d.Value] = true
		}
		if !dataSet["first"] || !dataSet["third"] {
			return fmt.Errorf("GSI query items: expected 'first' and 'third', got %v", dataSet)
		}

		gsiScanResp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(gsiTestTable),
			IndexName: aws.String("gsi1"),
		})
		if err != nil {
			return fmt.Errorf("GSI scan: %v", err)
		}
		if gsiScanResp.Count != 3 {
			return fmt.Errorf("GSI scan count: expected 3, got %d", gsiScanResp.Count)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "LSI_Query", func() error {
		lsiTable := fmt.Sprintf("LSIQuery-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(lsiTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("lsi_sk"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
			},
			BillingMode: types.BillingModePayPerRequest,
			LocalSecondaryIndexes: []types.LocalSecondaryIndex{
				{
					IndexName: aws.String("lsi1"),
					KeySchema: []types.KeySchemaElement{
						{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
						{AttributeName: aws.String("lsi_sk"), KeyType: types.KeyTypeRange},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(lsiTable)})

		for _, item := range []map[string]types.AttributeValue{
			{"pk": &types.AttributeValueMemberS{Value: "p1"}, "sk": &types.AttributeValueMemberS{Value: "main_a"}, "lsi_sk": &types.AttributeValueMemberS{Value: "lsi_x"}, "val": &types.AttributeValueMemberS{Value: "item_x"}},
			{"pk": &types.AttributeValueMemberS{Value: "p1"}, "sk": &types.AttributeValueMemberS{Value: "main_b"}, "lsi_sk": &types.AttributeValueMemberS{Value: "lsi_x"}, "val": &types.AttributeValueMemberS{Value: "item_y"}},
			{"pk": &types.AttributeValueMemberS{Value: "p1"}, "sk": &types.AttributeValueMemberS{Value: "main_c"}, "lsi_sk": &types.AttributeValueMemberS{Value: "lsi_z"}, "val": &types.AttributeValueMemberS{Value: "item_z"}},
		} {
			if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(lsiTable),
				Item:      item,
			}); err != nil {
				return fmt.Errorf("put item: %v", err)
			}
		}

		lsiQueryResp, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(lsiTable),
			IndexName:              aws.String("lsi1"),
			KeyConditionExpression: aws.String("pk = :pk AND lsi_sk = :lsk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":  &types.AttributeValueMemberS{Value: "p1"},
				":lsk": &types.AttributeValueMemberS{Value: "lsi_x"},
			},
		})
		if err != nil {
			return fmt.Errorf("LSI query: %v", err)
		}
		if lsiQueryResp.Count != 2 {
			return fmt.Errorf("LSI query count: expected 2, got %d", lsiQueryResp.Count)
		}
		valSet := make(map[string]bool)
		for _, item := range lsiQueryResp.Items {
			v, ok := item["val"].(*types.AttributeValueMemberS)
			if !ok {
				return fmt.Errorf("LSI query item: val is not S")
			}
			valSet[v.Value] = true
		}
		if !valSet["item_x"] || !valSet["item_y"] {
			return fmt.Errorf("LSI query items: expected 'item_x' and 'item_y', got %v", valSet)
		}
		return nil
	}))

	return results
}

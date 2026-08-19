package testutil

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"vorpalstacks-sdk-tests/config"
)

// dynamoDBQueryScanConformanceTests pins the query/scan data-plane contract:
// the partition key equality requirement of Query, and pagination that must
// reach every secondary index member even when the page limit is smaller than
// the member count.
func (r *TestRunner) dynamoDBQueryScanConformanceTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	tableName := fmt.Sprintf("QsConf-%d", time.Now().UnixNano())
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("gsik"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("gsi-idx"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("gsik"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "QueryScanConformance_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create table: %v", err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

	// Non-members are written first so the base-table storage order starts
	// with items the index does not contain; a paginator that bounds the raw
	// scan before filtering would lose the trailing members.
	for i := 0; i < 7; i++ {
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: fmt.Sprintf("plain-%02d", i)},
				"val": &types.AttributeValueMemberS{Value: "v"},
			},
		})
		if err != nil {
			return append(results, TestResult{
				Service:  "dynamodb",
				TestName: "QueryScanConformance_Setup",
				Status:   "FAIL",
				Error:    fmt.Sprintf("put plain item: %v", err),
			})
		}
	}
	memberKeys := []string{"member-a", "member-b", "member-c"}
	for _, key := range memberKeys {
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":   &types.AttributeValueMemberS{Value: key},
				"gsik": &types.AttributeValueMemberS{Value: "g-" + key},
			},
		})
		if err != nil {
			return append(results, TestResult{
				Service:  "dynamodb",
				TestName: "QueryScanConformance_Setup",
				Status:   "FAIL",
				Error:    fmt.Sprintf("put member item: %v", err),
			})
		}
	}

	// A keys-only GSI table for projection conformance checks.
	keysTableName := fmt.Sprintf("QsKeys-%d", time.Now().UnixNano())
	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(keysTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("cat"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("cat-idx"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("cat"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeKeysOnly},
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "QueryScanConformance_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create keys-only table: %v", err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(keysTableName)})
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(keysTableName),
		Item: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberS{Value: "k1"},
			"cat":     &types.AttributeValueMemberS{Value: "c1"},
			"payload": &types.AttributeValueMemberS{Value: "not projected"},
		},
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "QueryScanConformance_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("put keys-only item: %v", err),
		})
	}

	results = append(results, r.RunTest("dynamodb", "Scan_GSI_Limit_ReachesAllItems", func() error {
		var collected []string
		var lastEvaluated map[string]types.AttributeValue
		pages := 0
		for {
			input := &dynamodb.ScanInput{
				TableName: aws.String(tableName),
				IndexName: aws.String("gsi-idx"),
				Limit:     aws.Int32(2),
			}
			if lastEvaluated != nil {
				input.ExclusiveStartKey = lastEvaluated
			}
			resp, err := client.Scan(ctx, input)
			if err != nil {
				return fmt.Errorf("scan page: %v", err)
			}
			for _, item := range resp.Items {
				idAttr, ok := item["id"].(*types.AttributeValueMemberS)
				if !ok {
					return fmt.Errorf("page item missing id attribute")
				}
				collected = append(collected, idAttr.Value)
			}
			pages++
			if resp.LastEvaluatedKey == nil || len(resp.LastEvaluatedKey) == 0 {
				break
			}
			lastEvaluated = resp.LastEvaluatedKey
			if pages > 10 {
				return fmt.Errorf("pagination did not terminate after %d pages", pages)
			}
		}
		if len(collected) != len(memberKeys) {
			return fmt.Errorf("expected %d index members across pages, got %d: %v", len(memberKeys), len(collected), collected)
		}
		seen := map[string]bool{}
		for _, key := range collected {
			seen[key] = true
		}
		for _, expected := range memberKeys {
			if !seen[expected] {
				return fmt.Errorf("index member %s missing from paginated results", expected)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_MissingPartitionKeyCondition", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			KeyConditionExpression: aws.String("val = :v"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberS{Value: "v"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException when the key condition omits the partition key")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "Query_GSI_MissingIndexKeyCondition", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			IndexName:              aws.String("gsi-idx"),
			KeyConditionExpression: aws.String("id = :v"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberS{Value: "member-a"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException when the key condition omits the index partition key")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_Count_Select", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(tableName),
			Select:    types.SelectCount,
		})
		if err != nil {
			return err
		}
		if resp.Count != 10 {
			return fmt.Errorf("expected count of 10 items, got %d", resp.Count)
		}
		if len(resp.Items) != 0 {
			return fmt.Errorf("COUNT select must not return items, got %d", len(resp.Items))
		}
		if resp.ScannedCount != 10 {
			return fmt.Errorf("expected scanned count 10, got %d", resp.ScannedCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_MissingTotalSegments", func() error {
		_, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(tableName),
			Segment:   aws.Int32(0),
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException when Segment is given without TotalSegments")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "Query_GSI_ConsistentRead_Rejected", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			IndexName:              aws.String("gsi-idx"),
			ConsistentRead:         aws.Bool(true),
			KeyConditionExpression: aws.String("gsik = :g"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":g": &types.AttributeValueMemberS{Value: "g-member-a"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for strongly consistent read on a global secondary index")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_StrongConsistent_RCU", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName:              aws.String(tableName),
			ConsistentRead:         aws.Bool(true),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		if err != nil {
			return err
		}
		if resp.ConsumedCapacity == nil || resp.ConsumedCapacity.CapacityUnits == nil {
			return fmt.Errorf("expected ConsumedCapacity with CapacityUnits")
		}
		expected := float64(resp.ScannedCount) * 1.0
		if *resp.ConsumedCapacity.CapacityUnits != expected {
			return fmt.Errorf("expected %.1f capacity units for a strongly consistent scan of %d items, got %.1f",
				expected, resp.ScannedCount, *resp.ConsumedCapacity.CapacityUnits)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_GSI_KeysOnly_ReturnsProjectedAttributes", func() error {
		resp, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String(keysTableName),
			IndexName: aws.String("cat-idx"),
		})
		if err != nil {
			return err
		}
		if resp.Count != 1 {
			return fmt.Errorf("expected 1 index member, got %d", resp.Count)
		}
		for name := range resp.Items[0] {
			if name != "id" && name != "cat" {
				return fmt.Errorf("keys-only index scan must not return attribute %s", name)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_GSI_AllAttributesNotProjected_Rejected", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(keysTableName),
			IndexName:              aws.String("cat-idx"),
			Select:                 types.SelectAllAttributes,
			KeyConditionExpression: aws.String("cat = :c"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":c": &types.AttributeValueMemberS{Value: "c1"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for ALL_ATTRIBUTES on a keys-only global secondary index")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "Query_Descending_ContinuesAfterDeletedCursor", func() error {
		// When the cursor item is deleted between pages of a descending
		// query, the next page must still return the remaining smaller
		// sort keys.
		descTable := fmt.Sprintf("QsDesc-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(descTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeN},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(descTable)})

		for i := 1; i <= 3; i++ {
			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(descTable),
				Item: map[string]types.AttributeValue{
					"id": &types.AttributeValueMemberS{Value: "pk"},
					"sk": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i)},
				},
			})
			if err != nil {
				return fmt.Errorf("put item %d: %v", i, err)
			}
		}

		page1, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(descTable),
			KeyConditionExpression: aws.String("id = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "pk"},
			},
			ScanIndexForward: aws.Bool(false),
			Limit:            aws.Int32(1),
		})
		if err != nil {
			return fmt.Errorf("first page: %v", err)
		}
		if len(page1.Items) != 1 || page1.LastEvaluatedKey == nil {
			return fmt.Errorf("first page must return one item with a cursor: %+v", page1)
		}
		if top, ok := page1.Items[0]["sk"].(*types.AttributeValueMemberN); !ok || top.Value != "3" {
			return fmt.Errorf("a descending query must return the highest sort key first: %+v", page1.Items)
		}
		if _, err = client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(descTable),
			Key:       page1.LastEvaluatedKey,
		}); err != nil {
			return fmt.Errorf("delete cursor item: %v", err)
		}

		page2, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(descTable),
			KeyConditionExpression: aws.String("id = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "pk"},
			},
			ScanIndexForward:  aws.Bool(false),
			ExclusiveStartKey: page1.LastEvaluatedKey,
		})
		if err != nil {
			return fmt.Errorf("second page: %v", err)
		}
		gotKeys := make([]string, 0, len(page2.Items))
		for _, item := range page2.Items {
			if n, ok := item["sk"].(*types.AttributeValueMemberN); ok {
				gotKeys = append(gotKeys, n.Value)
			} else {
				gotKeys = append(gotKeys, fmt.Sprintf("%T", item["sk"]))
			}
		}
		if !reflect.DeepEqual(gotKeys, []string{"2", "1"}) {
			return fmt.Errorf("remaining sort keys after the deleted cursor must be [2 1], got %v", gotKeys)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_KeyTypeMismatch_Rejected", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberN{Value: "12345"},
				"val": &types.AttributeValueMemberS{Value: "v"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for a number-valued string key attribute")
		}
		if err := AssertErrorContains(err, "Type mismatch for key"); err != nil {
			return err
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "GetItem_KeyTypeMismatch_Rejected", func() error {
		_, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberN{Value: "12345"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for a wrong-typed GetItem key")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_KeyTypeMismatch_Rejected", func() error {
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: {
					{
						PutRequest: &types.PutRequest{
							Item: map[string]types.AttributeValue{
								"id":  &types.AttributeValueMemberN{Value: "12345"},
								"val": &types.AttributeValueMemberS{Value: "v"},
							},
						},
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for a wrong-typed key in BatchWriteItem")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_ReturnValues_AllNew_Rejected", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:    aws.String(tableName),
			ReturnValues: types.ReturnValueAllNew,
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "rv-invalid"},
				"val": &types.AttributeValueMemberS{Value: "v"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: PutItem accepts only NONE and ALL_OLD")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_NumberSetDuplicates_Rejected", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "dup-set"},
				"val": &types.AttributeValueMemberS{Value: "v"},
				"ns":  &types.AttributeValueMemberNS{Value: []string{"1", "2", "1.0"}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: set elements must be unique")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_MultiplyOperator_Rejected", func() error {
		_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:        aws.String(tableName),
			Key:              map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "member-a"}},
			UpdateExpression: aws.String("SET gsik = :v * :w"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v": &types.AttributeValueMemberN{Value: "2"},
				":w": &types.AttributeValueMemberN{Value: "3"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: SET arithmetic supports only plus and minus")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_UndefinedValuePlaceholder_Rejected", func() error {
		_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:        aws.String(tableName),
			Key:              map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "member-a"}},
			UpdateExpression: aws.String("SET gsik = :missing"),
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: undefined value placeholder")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_RemoveUndefinedName_Rejected", func() error {
		_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:        aws.String(tableName),
			Key:              map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "member-a"}},
			UpdateExpression: aws.String("REMOVE #missing"),
			ExpressionAttributeNames: map[string]string{
				"#defined": "gsik",
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: undefined attribute name in REMOVE")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateItem_RemoveDefinedName_Succeeds", func() error {
		_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:        aws.String(tableName),
			Key:              map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "member-b"}},
			UpdateExpression: aws.String("REMOVE #defined"),
			ExpressionAttributeNames: map[string]string{
				"#defined": "gsik",
			},
		})
		if err != nil {
			return err
		}
		resp, getErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "member-b"}},
		})
		if getErr != nil {
			return getErr
		}
		if _, still := resp.Item["gsik"]; still {
			return fmt.Errorf("REMOVE with a defined name must remove the attribute: %v", resp.Item)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_DuplicatePut_Rejected", func() error {
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: {
					{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "dup-key"},
					}}},
					{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "dup-key"},
					}}},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: two puts with identical keys in one request")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_PutAndDeleteSameItem_Rejected", func() error {
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: {
					{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "dup-key"},
					}}},
					{DeleteRequest: &types.DeleteRequest{Key: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "dup-key"},
					}}},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: put and delete of the same item in one request")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "BatchGetItem_DuplicateKeys_Rejected", func() error {
		_, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				tableName: {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "member-a"}},
						{"id": &types.AttributeValueMemberS{Value: "member-a"}},
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException: duplicate keys in BatchGetItem request")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "BatchGetItem_NonexistentKeys_ConsumeCapacity", func() error {
		// Requests for nonexistent items consume the minimum read capacity
		// units, so two present keys plus one absent key are charged as
		// three reads; a table with only an absent key still reports the
		// charge.
		resp, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				tableName: {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "member-a"}},
						{"id": &types.AttributeValueMemberS{Value: "member-b"}},
						{"id": &types.AttributeValueMemberS{Value: "absent-key"}},
					},
				},
				keysTableName: {
					Keys: []map[string]types.AttributeValue{
						{"id": &types.AttributeValueMemberS{Value: "absent-in-keys"}},
					},
				},
			},
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
		if err != nil {
			return err
		}
		capacity := map[string]float64{}
		for _, cc := range resp.ConsumedCapacity {
			if cc.TableName == nil || cc.CapacityUnits == nil {
				return fmt.Errorf("ConsumedCapacity entry missing TableName or CapacityUnits: %+v", cc)
			}
			capacity[*cc.TableName] = *cc.CapacityUnits
		}
		if got := capacity[tableName]; got != 1.5 {
			return fmt.Errorf("expected 1.5 capacity units for %s (three reads at 0.5 each), got %v", tableName, got)
		}
		if got := capacity[keysTableName]; got != 0.5 {
			return fmt.Errorf("expected 0.5 capacity units for %s (one read of an absent key), got %v", keysTableName, got)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_NonexistentTable", func() error {
		_, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String("NoSuchTxnTable"),
					Item: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "k"},
					},
				}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ResourceNotFoundException for a nonexistent transaction table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_InvalidTableName_Rejected", func() error {
		_, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String("ab"),
					Item: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "k"},
					},
				}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for an invalid transaction table name")
		}
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_DuplicateTarget_StayedCanceled", func() error {
		_, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String(tableName),
					Item: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "txn-dup"},
					},
				}},
				{Delete: &types.Delete{
					TableName: aws.String(tableName),
					Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "txn-dup"}},
				}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected TransactionCanceledException: the same item is targeted twice")
		}
		var canceled *types.TransactionCanceledException
		if !errors.As(err, &canceled) {
			return fmt.Errorf("expected TransactionCanceledException, got: %T: %v", err, err)
		}
		if len(canceled.CancellationReasons) != 2 {
			return fmt.Errorf("expected 2 cancellation reasons, got %d", len(canceled.CancellationReasons))
		}
		if canceled.CancellationReasons[1].Code == nil || *canceled.CancellationReasons[1].Code != "ValidationError" {
			return fmt.Errorf("expected ValidationError reason on the duplicate item, got %+v", canceled.CancellationReasons[1])
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_ConditionFailure_ReasonMessage", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":  &types.AttributeValueMemberS{Value: "txn-cond"},
				"val": &types.AttributeValueMemberS{Value: "original"},
			},
		})
		if err != nil {
			return err
		}
		_, err = client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName:           aws.String(tableName),
					ConditionExpression: aws.String("val = :expected"),
					Item: map[string]types.AttributeValue{
						"id":  &types.AttributeValueMemberS{Value: "txn-cond"},
						"val": &types.AttributeValueMemberS{Value: "replaced"},
					},
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":expected": &types.AttributeValueMemberS{Value: "different"},
					},
				}},
			},
		})
		if err == nil {
			return fmt.Errorf("expected TransactionCanceledException for the unmet condition")
		}
		var canceled *types.TransactionCanceledException
		if !errors.As(err, &canceled) {
			return fmt.Errorf("expected TransactionCanceledException, got: %T: %v", err, err)
		}
		if len(canceled.CancellationReasons) != 1 {
			return fmt.Errorf("expected 1 cancellation reason, got %d", len(canceled.CancellationReasons))
		}
		reason := canceled.CancellationReasons[0]
		if reason.Code == nil || *reason.Code != "ConditionalCheckFailed" {
			return fmt.Errorf("expected ConditionalCheckFailed code, got %+v", reason)
		}
		if reason.Message == nil || *reason.Message != "The conditional request failed" {
			return fmt.Errorf("expected the conditional-failure message, got %+v", reason)
		}
		// The transaction must have rolled back: the original value stays.
		resp, getErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "txn-cond"}},
		})
		if getErr != nil {
			return getErr
		}
		valAttr, ok := resp.Item["val"].(*types.AttributeValueMemberS)
		if !ok || valAttr.Value != "original" {
			return fmt.Errorf("cancelled transaction must not modify the item: %v", resp.Item)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_IdempotentRetry", func() error {
		token := fmt.Sprintf("idem-%d", time.Now().UnixNano())
		input := &dynamodb.TransactWriteItemsInput{
			ClientRequestToken: aws.String(token),
			TransactItems: []types.TransactWriteItem{
				{Update: &types.Update{
					TableName:        aws.String(tableName),
					Key:              map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "idem-key"}},
					UpdateExpression: aws.String("ADD cnt :one"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":one": &types.AttributeValueMemberN{Value: "1"},
					},
				}},
			},
		}
		for i := 0; i < 2; i++ {
			_, err := client.TransactWriteItems(ctx, input)
			if err != nil {
				return fmt.Errorf("attempt %d: %v", i+1, err)
			}
		}
		// An ADD re-applied by the retry would leave 2 behind; the
		// idempotent replay must leave the effect of a single execution.
		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "idem-key"}},
		})
		if err != nil {
			return err
		}
		valAttr, ok := resp.Item["cnt"].(*types.AttributeValueMemberN)
		if !ok || valAttr.Value != "1" {
			return fmt.Errorf("idempotent retry must not re-apply the write: %v", resp.Item)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_IdempotentMismatch", func() error {
		token := fmt.Sprintf("idem-mm-%d", time.Now().UnixNano())
		first := &dynamodb.TransactWriteItemsInput{
			ClientRequestToken: aws.String(token),
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String(tableName),
					Item: map[string]types.AttributeValue{
						"id":  &types.AttributeValueMemberS{Value: "idem-mm-key"},
						"val": &types.AttributeValueMemberS{Value: "first"},
					},
				}},
			},
		}
		if _, err := client.TransactWriteItems(ctx, first); err != nil {
			return fmt.Errorf("first attempt: %v", err)
		}
		second := &dynamodb.TransactWriteItemsInput{
			ClientRequestToken: aws.String(token),
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String(tableName),
					Item: map[string]types.AttributeValue{
						"id":  &types.AttributeValueMemberS{Value: "idem-mm-key"},
						"val": &types.AttributeValueMemberS{Value: "different-payload"},
					},
				}},
			},
		}
		_, err := client.TransactWriteItems(ctx, second)
		if err == nil {
			return fmt.Errorf("expected IdempotentParameterMismatchException for a reused token with a different payload")
		}
		var mismatch *types.IdempotentParameterMismatchException
		if !errors.As(err, &mismatch) {
			return fmt.Errorf("expected IdempotentParameterMismatchException, got: %T: %v", err, err)
		}
		resp, getErr := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "idem-mm-key"}},
		})
		if getErr != nil {
			return getErr
		}
		valAttr, ok := resp.Item["val"].(*types.AttributeValueMemberS)
		if !ok || valAttr.Value != "first" {
			return fmt.Errorf("mismatched replay must not execute: %v", resp.Item)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_ReplayCapacity_IsReadUnits", func() error {
		// The initial transactional write reports write capacity units; a
		// replay with the same client token reports read capacity units.
		token := fmt.Sprintf("idem-cap-%d", time.Now().UnixNano())
		input := &dynamodb.TransactWriteItemsInput{
			ClientRequestToken:     aws.String(token),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String(tableName),
					Item: map[string]types.AttributeValue{
						"id":  &types.AttributeValueMemberS{Value: "idem-cap-key"},
						"val": &types.AttributeValueMemberS{Value: "v"},
					},
				}},
			},
		}
		first, err := client.TransactWriteItems(ctx, input)
		if err != nil {
			return fmt.Errorf("first attempt: %v", err)
		}
		if len(first.ConsumedCapacity) == 0 {
			return fmt.Errorf("initial call must report ConsumedCapacity")
		}
		cc := first.ConsumedCapacity[0]
		if cc.WriteCapacityUnits == nil || *cc.WriteCapacityUnits != 2.0 {
			return fmt.Errorf("initial call must report write capacity units of 2.0, got %+v", cc)
		}
		if cc.ReadCapacityUnits != nil {
			return fmt.Errorf("initial call must not report read capacity units, got %+v", cc)
		}
		replay, err := client.TransactWriteItems(ctx, input)
		if err != nil {
			return fmt.Errorf("replay attempt: %v", err)
		}
		if len(replay.ConsumedCapacity) == 0 {
			return fmt.Errorf("replay must report ConsumedCapacity")
		}
		rcc := replay.ConsumedCapacity[0]
		if rcc.ReadCapacityUnits == nil || *rcc.ReadCapacityUnits != 2.0 {
			return fmt.Errorf("replay must report read capacity units of 2.0, got %+v", rcc)
		}
		if rcc.WriteCapacityUnits != nil {
			return fmt.Errorf("replay must not report write capacity units, got %+v", rcc)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_ConcurrentSameToken_ExecutedOnce", func() error {
		// Concurrent retries of one client request token must not execute
		// the transaction twice: a caller that loses the race observes the
		// in-progress claim and fails with TransactionInProgressException.
		token := fmt.Sprintf("idem-race-%d", time.Now().UnixNano())
		newInput := func() *dynamodb.TransactWriteItemsInput {
			return &dynamodb.TransactWriteItemsInput{
				ClientRequestToken: aws.String(token),
				TransactItems: []types.TransactWriteItem{
					{Update: &types.Update{
						TableName:        aws.String(tableName),
						Key:              map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "idem-race-key"}},
						UpdateExpression: aws.String("ADD raceCnt :one"),
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":one": &types.AttributeValueMemberN{Value: "1"},
						},
					}},
				},
			}
		}
		const callers = 2
		var wg sync.WaitGroup
		outcomes := make([]error, callers)
		for i := 0; i < callers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, err := client.TransactWriteItems(ctx, newInput())
				outcomes[i] = err
			}(i)
		}
		wg.Wait()
		for i, err := range outcomes {
			if err == nil {
				continue
			}
			var inProgress *types.TransactionInProgressException
			if !errors.As(err, &inProgress) {
				return fmt.Errorf("caller %d: unexpected error: %v", i, err)
			}
		}
		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "idem-race-key"}},
		})
		if err != nil {
			return err
		}
		cnt, ok := resp.Item["raceCnt"].(*types.AttributeValueMemberN)
		if !ok || cnt.Value != "1" {
			return fmt.Errorf("the transaction must take effect exactly once, raceCnt was %v", resp.Item["raceCnt"])
		}
		return nil
	}))

	return results
}

// dynamoDBGlobalTableReplicationTests verifies that batch and transaction
// writes replicate to global table replica regions, not only single-item
// writes.
func (r *TestRunner) dynamoDBGlobalTableReplicationTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	replicaRegion := "us-west-2"
	if r.region == replicaRegion {
		replicaRegion = "eu-west-1"
	}
	replicaCfg, cfgErr := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   replicaRegion,
	})
	if cfgErr != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "GlobalTableReplication_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("load replica config: %v", cfgErr),
		})
	}
	replicaClient := dynamodb.NewFromConfig(replicaCfg)

	tableName := fmt.Sprintf("GtRepl-%d", time.Now().UnixNano())
	for _, c := range []*dynamodb.Client{client, replicaClient} {
		_, err := c.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
			StreamSpecification: &types.StreamSpecification{
				StreamEnabled:  aws.Bool(true),
				StreamViewType: types.StreamViewTypeNewAndOldImages,
			},
		})
		if err != nil {
			return append(results, TestResult{
				Service:  "dynamodb",
				TestName: "GlobalTableReplication_Setup",
				Status:   "FAIL",
				Error:    fmt.Sprintf("create table in region: %v", err),
			})
		}
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
	defer replicaClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

	_, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
		GlobalTableName: aws.String(tableName),
		ReplicationGroup: []types.Replica{
			{RegionName: aws.String(r.region)},
			{RegionName: aws.String(replicaRegion)},
		},
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "GlobalTableReplication_Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create global table: %v", err),
		})
	}

	waitForReplica := func(key string, expectFound bool) error {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			resp, getErr := replicaClient.GetItem(ctx, &dynamodb.GetItemInput{
				TableName: aws.String(tableName),
				Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: key}},
			})
			if getErr != nil {
				return fmt.Errorf("replica get %s: %v", key, getErr)
			}
			found := resp.Item != nil && len(resp.Item) > 0
			if found == expectFound {
				return nil
			}
			time.Sleep(200 * time.Millisecond)
		}
		return fmt.Errorf("replica region %s did not observe key %s (expected present=%v) within 5s", replicaRegion, key, expectFound)
	}

	results = append(results, r.RunTest("dynamodb", "BatchWriteItem_ReplicatesToGlobalTableReplica", func() error {
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: {
					{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"id":  &types.AttributeValueMemberS{Value: "batch-key"},
						"val": &types.AttributeValueMemberS{Value: "from-batch"},
					}}},
				},
			},
		})
		if err != nil {
			return err
		}
		if waitErr := waitForReplica("batch-key", true); waitErr != nil {
			return waitErr
		}
		resp, getErr := replicaClient.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "batch-key"}},
		})
		if getErr != nil {
			return getErr
		}
		valAttr, ok := resp.Item["val"].(*types.AttributeValueMemberS)
		if !ok || valAttr.Value != "from-batch" {
			return fmt.Errorf("replica item attribute mismatch: %v", resp.Item)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactWriteItems_ReplicatesToGlobalTableReplica", func() error {
		_, err := client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{Put: &types.Put{
					TableName: aws.String(tableName),
					Item: map[string]types.AttributeValue{
						"id":  &types.AttributeValueMemberS{Value: "transact-key"},
						"val": &types.AttributeValueMemberS{Value: "from-transact"},
					},
				}},
			},
		})
		if err != nil {
			return err
		}
		if waitErr := waitForReplica("transact-key", true); waitErr != nil {
			return waitErr
		}
		resp, getErr := replicaClient.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "transact-key"}},
		})
		if getErr != nil {
			return getErr
		}
		valAttr, ok := resp.Item["val"].(*types.AttributeValueMemberS)
		if !ok || valAttr.Value != "from-transact" {
			return fmt.Errorf("replica item attribute mismatch: %v", resp.Item)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "BatchWriteDelete_ReplicatesToGlobalTableReplica", func() error {
		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: {
					{DeleteRequest: &types.DeleteRequest{Key: map[string]types.AttributeValue{
						"id": &types.AttributeValueMemberS{Value: "batch-key"},
					}}},
				},
			},
		})
		if err != nil {
			return err
		}
		return waitForReplica("batch-key", false)
	}))

	return results
}

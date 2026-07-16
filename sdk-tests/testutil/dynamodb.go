package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	// Phase 1: Main table lifecycle (shared state)
	results = append(results, r.dynamoDBCreateTableTests(ctx, client, tableName)...)
	results = append(results, r.dynamoDBBasicCRUD(ctx, client, tableName)...)
	results = append(results, r.dynamoDBBasicQueryScan(ctx, client, tableName)...)
	results = append(results, r.dynamoDBBasicBatchTests(ctx, client, tableName)...)
	results = append(results, r.dynamoDBDeleteItemTest(ctx, client, tableName)...)
	results = append(results, r.dynamoDBMainAdvancedTests(ctx, client, tableName, tableARN)...)
	results = append(results, r.dynamoDBBasicPartiQLTests(ctx, client, tableName)...)
	results = append(results, r.dynamoDBBasicTransactionTests(ctx, client, tableName)...)
	results = append(results, r.dynamoDBUpdateTableTests(ctx, client, tableName)...)

	// Self-contained edge cases (create own tables)
	results = append(results, r.dynamoDBPartiQLInsertParamsTest(ctx, client)...)
	results = append(results, r.dynamoDBConditionNotCaseTest(ctx, client)...)

	// Delete main table
	results = append(results, r.dynamoDBDeleteTableTest(ctx, client, tableName)...)

	// Phase 2: Error cases (self-contained)
	results = append(results, r.dynamoDBNonExistentTableTests(ctx, client)...)
	results = append(results, r.dynamoDBConditionalCheckTests(ctx, client)...)
	results = append(results, r.dynamoDBBatchNonExistentTest(ctx, client)...)
	results = append(results, r.dynamoDBDuplicateTableTest(ctx, client)...)

	// Phase 3: Return value and consumed capacity (self-contained)
	results = append(results, r.dynamoDBReturnValueTests(ctx, client)...)

	// Phase 4: Composite key table (shared compTableName)
	compTableName := fmt.Sprintf("CompTable-%d", time.Now().UnixNano())
	results = append(results, r.dynamoDBCompositeKeySetup(ctx, client, compTableName)...)
	results = append(results, r.dynamoDBPutConditionTests(ctx, client)...)
	results = append(results, r.dynamoDBProjectionTests(ctx, client, compTableName)...)
	results = append(results, r.dynamoDBDeleteEdgeCaseTests(ctx, client)...)
	results = append(results, r.dynamoDBUpdateExpressionTests(ctx, client)...)
	results = append(results, r.dynamoDBConditionExpressionTests(ctx, client)...)
	results = append(results, r.dynamoDBCompositeQueryScanTests(ctx, client, compTableName)...)
	results = append(results, r.dynamoDBPartiQLEdgeCaseTests(ctx, client, compTableName)...)
	results = append(results, r.dynamoDBTransactionEdgeCaseTests(ctx, client, compTableName)...)
	results = append(results, r.dynamoDBBatchEdgeCaseTests(ctx, client)...)
	client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &compTableName})

	// Phase 5: Table edge cases (self-contained)
	results = append(results, r.dynamoDBTableEdgeCaseTests(ctx, client)...)

	// Phase 6: Pagination (self-contained)
	results = append(results, r.dynamoDBPaginationTests(ctx, client)...)

	// Phase 7: Types and index queries (self-contained)
	results = append(results, r.dynamoDBTypeTests(ctx, client)...)
	results = append(results, r.dynamoDBNestedTypeTests(ctx, client)...)
	results = append(results, r.dynamoDBIndexQueryTests(ctx, client)...)

	return results
}

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
	tableARN := fmt.Sprintf("arn:aws:dynamodb:%s:%s:table/%s", r.region, r.accountID, tableName)

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

	// Item collection metrics on self-contained LSI and plain fixtures
	results = append(results, r.dynamoDBItemCollectionMetricsTests(ctx, client)...)

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

	// Phase 4b: Query/Scan conformance (self-contained)
	results = append(results, r.dynamoDBQueryScanConformanceTests(ctx, client)...)
	results = append(results, r.dynamoDBGlobalTableReplicationTests(ctx, client)...)

	// Phase 4c: Point-in-time recovery (self-contained)
	results = append(results, r.dynamoDBPITRTests(ctx, client)...)

	// Phase 4d: Export/import lifecycle (self-contained)
	results = append(results, r.dynamoDBExportImportTests(ctx, client)...)

	// Phase 4e: Kinesis streaming destinations (self-contained)
	results = append(results, r.dynamoDBKinesisDestinationTests(ctx, client)...)

	// Phase 4f: Global table creation preconditions (self-contained)
	results = append(results, r.dynamoDBGlobalTableValidationTests(ctx, client)...)

	// Phase 4g: DescribeLimits/DescribeEndpoints shapes (self-contained)
	results = append(results, r.dynamoDBUtilityTests(ctx, client)...)

	// Phase 4h: TTL update contract (self-contained)
	results = append(results, r.dynamoDBTTLValidationTests(ctx, client)...)

	// Phase 4i: Input validation contracts (self-contained)
	results = append(results, r.dynamoDBInputValidationTests(ctx, client)...)

	// Phase 4j: DynamoDB Streams read path (self-contained)
	results = append(results, r.dynamoDBStreamsTests(ctx, client)...)

	// Phase 4k: Contributor insights rules and access aggregation (self-contained)
	results = append(results, r.dynamoDBContributorInsightsTests(ctx, client)...)
	results = append(results, r.dynamoDBContributorReadPathsTests(ctx, client)...)

	// Phase 4l: Baseline coverage for previously untested operations (self-contained)
	results = append(results, r.dynamoDBBaselineCoverageTests(ctx, client)...)

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

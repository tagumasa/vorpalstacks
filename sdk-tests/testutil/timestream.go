package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunTimestreamTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "timestream",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := timestreamwrite.NewFromConfig(cfg)
	ctx := context.Background()

	databaseName := fmt.Sprintf("test-db-%d", time.Now().UnixNano())
	tableName := fmt.Sprintf("test-table-%d", time.Now().UnixNano())

	// Test 1: Create Database
	results = append(results, r.RunTest("timestream", "CreateDatabase", func() error {
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(databaseName),
		})
		return err
	}))

	// Test 2: List Databases
	results = append(results, r.RunTest("timestream", "ListDatabases", func() error {
		_, err := client.ListDatabases(ctx, &timestreamwrite.ListDatabasesInput{
			MaxResults: aws.Int32(10),
		})
		return err
	}))

	// Test 3: Describe Database
	results = append(results, r.RunTest("timestream", "DescribeDatabase", func() error {
		_, err := client.DescribeDatabase(ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String(databaseName),
		})
		return err
	}))

	// Test 4: Create Table
	results = append(results, r.RunTest("timestream", "CreateTable", func() error {
		_, err := client.CreateTable(ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	// Test 5: List Tables
	results = append(results, r.RunTest("timestream", "ListTables", func() error {
		_, err := client.ListTables(ctx, &timestreamwrite.ListTablesInput{
			DatabaseName: aws.String(databaseName),
			MaxResults:   aws.Int32(10),
		})
		return err
	}))

	// Test 6: Describe Table
	results = append(results, r.RunTest("timestream", "DescribeTable", func() error {
		_, err := client.DescribeTable(ctx, &timestreamwrite.DescribeTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	// Test 7: Write Records
	results = append(results, r.RunTest("timestream", "WriteRecords", func() error {
		_, err := client.WriteRecords(ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
			Records: []types.Record{
				{
					Dimensions: []types.Dimension{
						{
							Name:  aws.String("dim1"),
							Value: aws.String("value1"),
						},
					},
					MeasureName:      aws.String("metric1"),
					MeasureValue:     aws.String("100"),
					MeasureValueType: types.MeasureValueTypeDouble,
					Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
					TimeUnit:         types.TimeUnitMilliseconds,
				},
			},
		})
		return err
	}))

	// Test 8: Update Table
	results = append(results, r.RunTest("timestream", "UpdateTable", func() error {
		_, err := client.UpdateTable(ctx, &timestreamwrite.UpdateTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
			RetentionProperties: &types.RetentionProperties{
				MemoryStoreRetentionPeriodInHours:  aws.Int64(24),
				MagneticStoreRetentionPeriodInDays: aws.Int64(7),
			},
		})
		return err
	}))

	// Test 9: Describe Endpoints
	results = append(results, r.RunTest("timestream", "DescribeEndpoints", func() error {
		_, err := client.DescribeEndpoints(ctx, &timestreamwrite.DescribeEndpointsInput{})
		return err
	}))

	// Test 10: Delete Table
	results = append(results, r.RunTest("timestream", "DeleteTable", func() error {
		_, err := client.DeleteTable(ctx, &timestreamwrite.DeleteTableInput{
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	// Test 11: Update Database
	results = append(results, r.RunTest("timestream", "UpdateDatabase", func() error {
		_, err := client.UpdateDatabase(ctx, &timestreamwrite.UpdateDatabaseInput{
			DatabaseName: aws.String(databaseName),
			KmsKeyId:     aws.String("arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"),
		})
		return err
	}))

	// Test 12: Delete Database
	results = append(results, r.RunTest("timestream", "DeleteDatabase", func() error {
		_, err := client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(databaseName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("timestream", "DescribeDatabase_NonExistent", func() error {
		_, err := client.DescribeDatabase(ctx, &timestreamwrite.DescribeDatabaseInput{
			DatabaseName: aws.String("nonexistent-db-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent database")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeTable_NonExistent", func() error {
		_, err := client.DescribeTable(ctx, &timestreamwrite.DescribeTableInput{
			DatabaseName: aws.String("nonexistent-db-xyz"),
			TableName:    aws.String("nonexistent-table-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CreateDatabase_Duplicate", func() error {
		dupDBName := fmt.Sprintf("dup-db-%d", time.Now().UnixNano())
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})

		_, err = client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(dupDBName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate database")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "WriteRecords_GetRecords_Roundtrip", func() error {
		rtDBName := fmt.Sprintf("rt-db-%d", time.Now().UnixNano())
		rtTableName := fmt.Sprintf("rt-table-%d", time.Now().UnixNano())
		_, err := client.CreateDatabase(ctx, &timestreamwrite.CreateDatabaseInput{
			DatabaseName: aws.String(rtDBName),
		})
		if err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		defer client.DeleteDatabase(ctx, &timestreamwrite.DeleteDatabaseInput{
			DatabaseName: aws.String(rtDBName),
		})

		_, err = client.CreateTable(ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(rtDBName),
			TableName:    aws.String(rtTableName),
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}

		measureValue := fmt.Sprintf("verify-%d", time.Now().UnixNano())
		_, err = client.WriteRecords(ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(rtDBName),
			TableName:    aws.String(rtTableName),
			Records: []types.Record{
				{
					MeasureName:      aws.String("cpu_utilization"),
					MeasureValue:     aws.String(measureValue),
					MeasureValueType: types.MeasureValueTypeDouble,
					Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
					TimeUnit:         types.TimeUnitMilliseconds,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("write: %v", err)
		}
		return nil
	}))

	return results
}

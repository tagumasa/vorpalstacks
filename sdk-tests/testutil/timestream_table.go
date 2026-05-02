package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
)

func (r *TestRunner) runTimestreamTableTests(tc *tsTestContext) []TestResult {
	var results []TestResult
	dbName := tc.uniqueName("tbl-db")
	tableName := tc.uniqueName("tbl")

	results = append(results, r.RunTest("timestream", "TableLifeycle_Setup", func() error {
		if err := tc.createDatabase(dbName); err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CreateTable_ContentVerify", func() error {
		resp, err := tc.writeClient.CreateTable(tc.ctx, &timestreamwrite.CreateTableInput{
			DatabaseName: aws.String(dbName),
			TableName:    aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.Table == nil {
			return fmt.Errorf("Table is nil")
		}
		if resp.Table.TableName == nil || *resp.Table.TableName != tableName {
			return fmt.Errorf("expected TableName=%s, got %v", tableName, resp.Table.TableName)
		}
		if resp.Table.DatabaseName == nil || *resp.Table.DatabaseName != dbName {
			return fmt.Errorf("expected DatabaseName=%s, got %v", dbName, resp.Table.DatabaseName)
		}
		if resp.Table.Arn == nil || *resp.Table.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		if resp.Table.CreationTime == nil {
			return fmt.Errorf("CreationTime is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeTable_ContentVerify", func() error {
		resp, err := tc.writeClient.DescribeTable(tc.ctx, &timestreamwrite.DescribeTableInput{
			DatabaseName: aws.String(dbName),
			TableName:    aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.Table == nil {
			return fmt.Errorf("Table is nil")
		}
		if resp.Table.TableName == nil || *resp.Table.TableName != tableName {
			return fmt.Errorf("expected TableName=%s, got %v", tableName, resp.Table.TableName)
		}
		if resp.Table.DatabaseName == nil || *resp.Table.DatabaseName != dbName {
			return fmt.Errorf("expected DatabaseName=%s, got %v", dbName, resp.Table.DatabaseName)
		}
		if resp.Table.Arn == nil || *resp.Table.Arn == "" {
			return fmt.Errorf("Arn is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListTables_ContainsCreated", func() error {
		resp, err := tc.writeClient.ListTables(tc.ctx, &timestreamwrite.ListTablesInput{
			DatabaseName: aws.String(dbName),
			MaxResults:   aws.Int32(20),
		})
		if err != nil {
			return err
		}
		found := false
		for _, t := range resp.Tables {
			if t.TableName != nil && *t.TableName == tableName {
				found = true
				if t.Arn == nil || *t.Arn == "" {
					return fmt.Errorf("listed table has nil/empty Arn")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("table %s not found in list", tableName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "WriteRecords_RecordsIngested", func() error {
		resp, err := tc.writeClient.WriteRecords(tc.ctx, &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String(dbName),
			TableName:    aws.String(tableName),
			Records: []types.Record{
				{
					Dimensions: []types.Dimension{
						{Name: aws.String("host"), Value: aws.String("server1")},
					},
					MeasureName:      aws.String("cpu_utilization"),
					MeasureValue:     aws.String("85.5"),
					MeasureValueType: types.MeasureValueTypeDouble,
					Time:             aws.String(fmt.Sprintf("%d", time.Now().UnixMilli())),
					TimeUnit:         types.TimeUnitMilliseconds,
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.RecordsIngested == nil {
			return fmt.Errorf("RecordsIngested is nil")
		}
		if resp.RecordsIngested.Total != 1 {
			return fmt.Errorf("expected RecordsIngested.Total=1, got %d", resp.RecordsIngested.Total)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UpdateTable_Retention", func() error {
		resp, err := tc.writeClient.UpdateTable(tc.ctx, &timestreamwrite.UpdateTableInput{
			DatabaseName: aws.String(dbName),
			TableName:    aws.String(tableName),
			RetentionProperties: &types.RetentionProperties{
				MemoryStoreRetentionPeriodInHours:  aws.Int64(48),
				MagneticStoreRetentionPeriodInDays: aws.Int64(30),
			},
		})
		if err != nil {
			return err
		}
		if resp.Table == nil {
			return fmt.Errorf("Table is nil")
		}
		if resp.Table.RetentionProperties == nil {
			return fmt.Errorf("RetentionProperties is nil")
		}
		if *resp.Table.RetentionProperties.MemoryStoreRetentionPeriodInHours != 48 {
			return fmt.Errorf("expected MemoryStoreRetentionPeriodInHours=48, got %d", *resp.Table.RetentionProperties.MemoryStoreRetentionPeriodInHours)
		}
		if *resp.Table.RetentionProperties.MagneticStoreRetentionPeriodInDays != 30 {
			return fmt.Errorf("expected MagneticStoreRetentionPeriodInDays=30, got %d", *resp.Table.RetentionProperties.MagneticStoreRetentionPeriodInDays)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeEndpoints_Write", func() error {
		resp, err := tc.writeClient.DescribeEndpoints(tc.ctx, &timestreamwrite.DescribeEndpointsInput{})
		if err != nil {
			return err
		}
		if len(resp.Endpoints) == 0 {
			return fmt.Errorf("expected at least 1 endpoint")
		}
		if resp.Endpoints[0].Address == nil || *resp.Endpoints[0].Address == "" {
			return fmt.Errorf("endpoint Address is nil or empty")
		}
		if resp.Endpoints[0].CachePeriodInMinutes == 0 {
			return fmt.Errorf("endpoint CachePeriodInMinutes is 0")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DeleteTable_Basic", func() error {
		_, err := tc.writeClient.DeleteTable(tc.ctx, &timestreamwrite.DeleteTableInput{
			DatabaseName: aws.String(dbName),
			TableName:    aws.String(tableName),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeTable_NonExistent", func() error {
		_, err := tc.writeClient.DescribeTable(tc.ctx, &timestreamwrite.DescribeTableInput{
			DatabaseName: aws.String("nonexistent-db-xyz"),
			TableName:    aws.String("nonexistent-table-xyz"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("timestream", "TableLifecycle_Cleanup", func() error {
		tc.deleteDatabase(dbName)
		return nil
	}))

	return results
}

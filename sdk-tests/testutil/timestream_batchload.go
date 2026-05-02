package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
)

func (r *TestRunner) runTimestreamBatchLoadTests(tc *tsTestContext) []TestResult {
	var results []TestResult
	blDBName := tc.uniqueName("bl-db")
	blTableName := tc.uniqueName("bl-table")
	blTaskID := ""

	results = append(results, r.RunTest("timestream", "BatchLoad_Setup", func() error {
		if err := tc.createDatabase(blDBName); err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		if err := tc.createTable(blDBName, blTableName); err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CreateBatchLoadTask_Verify", func() error {
		resp, err := tc.writeClient.CreateBatchLoadTask(tc.ctx, &timestreamwrite.CreateBatchLoadTaskInput{
			TargetDatabaseName: aws.String(blDBName),
			TargetTableName:    aws.String(blTableName),
			DataSourceConfiguration: &types.DataSourceConfiguration{
				DataFormat: types.BatchLoadDataFormatCsv,
				DataSourceS3Configuration: &types.DataSourceS3Configuration{
					BucketName: aws.String("test-bucket"),
				},
			},
			ReportConfiguration: &types.ReportConfiguration{
				ReportS3Configuration: &types.ReportS3Configuration{
					BucketName: aws.String("report-bucket"),
				},
			},
			DataModelConfiguration: &types.DataModelConfiguration{
				DataModel: &types.DataModel{
					DimensionMappings: []types.DimensionMapping{
						{
							SourceColumn:      aws.String("col1"),
							DestinationColumn: aws.String("dim1"),
						},
					},
					MeasureNameColumn: aws.String("measure"),
					TimeColumn:        aws.String("ts"),
					TimeUnit:          types.TimeUnitMilliseconds,
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId == "" {
			return fmt.Errorf("TaskId is nil or empty")
		}
		blTaskID = *resp.TaskId
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeBatchLoadTask_Verify", func() error {
		resp, err := tc.writeClient.DescribeBatchLoadTask(tc.ctx, &timestreamwrite.DescribeBatchLoadTaskInput{
			TaskId: aws.String(blTaskID),
		})
		if err != nil {
			return err
		}
		if resp.BatchLoadTaskDescription == nil {
			return fmt.Errorf("expected BatchLoadTaskDescription")
		}
		desc := resp.BatchLoadTaskDescription
		if desc.TaskId == nil || *desc.TaskId != blTaskID {
			return fmt.Errorf("expected TaskId=%s, got %v", blTaskID, desc.TaskId)
		}
		if desc.CreationTime == nil {
			return fmt.Errorf("expected CreationTime to be set")
		}
		if desc.TargetDatabaseName == nil || *desc.TargetDatabaseName != blDBName {
			return fmt.Errorf("expected TargetDatabaseName=%s, got %v", blDBName, desc.TargetDatabaseName)
		}
		if desc.TargetTableName == nil || *desc.TargetTableName != blTableName {
			return fmt.Errorf("expected TargetTableName=%s, got %v", blTableName, desc.TargetTableName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListBatchLoadTasks_Verify", func() error {
		resp, err := tc.writeClient.ListBatchLoadTasks(tc.ctx, &timestreamwrite.ListBatchLoadTasksInput{})
		if err != nil {
			return err
		}
		if len(resp.BatchLoadTasks) == 0 {
			return fmt.Errorf("expected at least 1 task, got 0")
		}
		found := false
		for _, t := range resp.BatchLoadTasks {
			if t.TaskId != nil && *t.TaskId == blTaskID {
				found = true
				if t.TaskStatus == "" {
					return fmt.Errorf("batch load task has empty TaskStatus")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("batch load task %s not found in list", blTaskID)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ResumeBatchLoadTask_CompletedTask", func() error {
		_, err := tc.writeClient.ResumeBatchLoadTask(tc.ctx, &timestreamwrite.ResumeBatchLoadTaskInput{
			TaskId: aws.String(blTaskID),
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	results = append(results, r.RunTest("timestream", "DescribeBatchLoadTask_NonExistent", func() error {
		_, err := tc.writeClient.DescribeBatchLoadTask(tc.ctx, &timestreamwrite.DescribeBatchLoadTaskInput{
			TaskId: aws.String("nonexistent-task-id"),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("timestream", "BatchLoad_Cleanup", func() error {
		tc.deleteTable(blDBName, blTableName)
		tc.deleteDatabase(blDBName)
		return nil
	}))

	return results
}

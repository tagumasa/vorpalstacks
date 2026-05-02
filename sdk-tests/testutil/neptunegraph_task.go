package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func (r *TestRunner) runNeptunegraphImportTaskTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	var importTaskID string

	results = append(results, r.RunTest("neptunegraph", "StartImportTask", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.StartImportTask(tc.ctx, &neptunegraph.StartImportTaskInput{
			GraphIdentifier: aws.String(tc.graphID),
			Source:          aws.String("s3://test-bucket/import-data/"),
			RoleArn:         aws.String("arn:aws:iam::000000000000:role/NeptuneImportRole"),
			Format:          types.FormatCsv,
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId == "" {
			return fmt.Errorf("expected non-empty import task ID")
		}
		if resp.Status != types.ImportTaskStatusInitializing && resp.Status != types.ImportTaskStatusImporting {
			return fmt.Errorf("expected initialising/inProgress status, got %s", resp.Status)
		}
		importTaskID = *resp.TaskId
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetImportTask", func() error {
		if importTaskID == "" {
			return fmt.Errorf("no import task ID")
		}
		resp, err := tc.client.GetImportTask(tc.ctx, &neptunegraph.GetImportTaskInput{
			TaskIdentifier: aws.String(importTaskID),
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId != importTaskID {
			return fmt.Errorf("expected taskId=%s, got %v", importTaskID, resp.TaskId)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListImportTasks", func() error {
		resp, err := tc.client.ListImportTasks(tc.ctx, &neptunegraph.ListImportTasksInput{})
		if err != nil {
			return err
		}
		if resp.Tasks == nil {
			return fmt.Errorf("expected non-nil Tasks list")
		}
		found := false
		for _, t := range resp.Tasks {
			if t.TaskId != nil && *t.TaskId == importTaskID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("import task not found in ListImportTasks")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "CancelImportTask", func() error {
		if importTaskID == "" {
			return fmt.Errorf("no import task ID")
		}
		resp, err := tc.client.CancelImportTask(tc.ctx, &neptunegraph.CancelImportTaskInput{
			TaskIdentifier: aws.String(importTaskID),
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId != importTaskID {
			return fmt.Errorf("expected taskId=%s, got %v", importTaskID, resp.TaskId)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "CreateGraphUsingImportTask", func() error {
		importGraphName := fmt.Sprintf("sdk-impgraph-%s", tc.tsNano[len(tc.tsNano)-6:])
		resp, err := tc.client.CreateGraphUsingImportTask(tc.ctx, &neptunegraph.CreateGraphUsingImportTaskInput{
			GraphName: aws.String(importGraphName),
			Source:    aws.String("s3://test-bucket/import-data/"),
			RoleArn:   aws.String("arn:aws:iam::000000000000:role/NeptuneImportRole"),
			Format:    types.FormatCsv,
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId == "" {
			return fmt.Errorf("expected non-empty task ID from CreateGraphUsingImportTask")
		}
		if resp.Status == "" {
			return fmt.Errorf("expected non-empty status")
		}
		_, _ = tc.client.CancelImportTask(tc.ctx, &neptunegraph.CancelImportTaskInput{
			TaskIdentifier: resp.TaskId,
		})
		if resp.GraphId != nil && *resp.GraphId != "" {
			_, _ = tc.client.DeleteGraph(tc.ctx, &neptunegraph.DeleteGraphInput{
				GraphIdentifier: resp.GraphId,
				SkipSnapshot:    aws.Bool(true),
			})
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunegraphExportTaskTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	var exportTaskID string

	results = append(results, r.RunTest("neptunegraph", "StartExportTask", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.StartExportTask(tc.ctx, &neptunegraph.StartExportTaskInput{
			GraphIdentifier:  aws.String(tc.graphID),
			Destination:      aws.String("s3://test-bucket/export-data/"),
			KmsKeyIdentifier: aws.String("arn:aws:kms:us-east-1:000000000000:key/12345678-1234-1234-1234-123456789012"),
			RoleArn:          aws.String("arn:aws:iam::000000000000:role/NeptuneExportRole"),
			Format:           types.ExportFormatCsv,
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId == "" {
			return fmt.Errorf("expected non-empty export task ID")
		}
		exportTaskID = *resp.TaskId
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetExportTask", func() error {
		if exportTaskID == "" {
			return fmt.Errorf("no export task ID")
		}
		resp, err := tc.client.GetExportTask(tc.ctx, &neptunegraph.GetExportTaskInput{
			TaskIdentifier: aws.String(exportTaskID),
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId != exportTaskID {
			return fmt.Errorf("expected taskId=%s, got %v", exportTaskID, resp.TaskId)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListExportTasks", func() error {
		resp, err := tc.client.ListExportTasks(tc.ctx, &neptunegraph.ListExportTasksInput{})
		if err != nil {
			return err
		}
		if resp.Tasks == nil {
			return fmt.Errorf("expected non-nil Tasks list")
		}
		found := false
		for _, t := range resp.Tasks {
			if t.TaskId != nil && *t.TaskId == exportTaskID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("export task not found in ListExportTasks")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListExportTasks_FilterByGraph", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.ListExportTasks(tc.ctx, &neptunegraph.ListExportTasksInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
		if err != nil {
			return err
		}
		if resp.Tasks == nil {
			return fmt.Errorf("expected non-nil Tasks list")
		}
		for _, t := range resp.Tasks {
			if t.TaskId != nil && *t.TaskId == exportTaskID {
				return nil
			}
		}
		return fmt.Errorf("export task not found when filtering by graph")
	}))

	results = append(results, r.RunTest("neptunegraph", "CancelExportTask", func() error {
		if exportTaskID == "" {
			return fmt.Errorf("no export task ID")
		}
		resp, err := tc.client.CancelExportTask(tc.ctx, &neptunegraph.CancelExportTaskInput{
			TaskIdentifier: aws.String(exportTaskID),
		})
		if err != nil {
			return err
		}
		if resp.TaskId == nil || *resp.TaskId != exportTaskID {
			return fmt.Errorf("expected taskId=%s, got %v", exportTaskID, resp.TaskId)
		}
		return nil
	}))

	return results
}

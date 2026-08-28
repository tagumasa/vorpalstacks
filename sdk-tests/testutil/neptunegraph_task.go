package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func (r *TestRunner) runNeptunegraphImportTaskTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	acct := tc.accountID
	var importTaskID string

	results = append(results, r.RunTest("neptunegraph", "StartImportTask", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.client.StartImportTask(tc.ctx, &neptunegraph.StartImportTaskInput{
			GraphIdentifier: aws.String(tc.graphID),
			Source:          aws.String("s3://test-bucket/import-data/"),
			RoleArn:         aws.String(fmt.Sprintf("arn:aws:iam::%s:role/NeptuneImportRole", acct)),
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
		if err := requireID(importTaskID, "import task ID"); err != nil {
			return err
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
		if err := requireID(importTaskID, "import task ID"); err != nil {
			return err
		}
		tasks, err := paginate(func(next *string) ([]types.ImportTaskSummary, *string, error) {
			resp, err := tc.client.ListImportTasks(tc.ctx, &neptunegraph.ListImportTasksInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return resp.Tasks, resp.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, t := range tasks {
			if t.TaskId != nil && *t.TaskId == importTaskID {
				return nil
			}
		}
		return fmt.Errorf("import task not found in ListImportTasks")
	}))

	results = append(results, r.RunTest("neptunegraph", "CancelImportTask", func() error {
		if err := requireID(importTaskID, "import task ID"); err != nil {
			return err
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
		importGraphName := tc.unique("sdk-impgraph")
		createTags := map[string]string{"env": "import-created", "suite": "neptunegraph"}
		resp, err := tc.client.CreateGraphUsingImportTask(tc.ctx, &neptunegraph.CreateGraphUsingImportTaskInput{
			GraphName: aws.String(importGraphName),
			Source:    aws.String("s3://test-bucket/import-data/"),
			RoleArn:   aws.String(fmt.Sprintf("arn:aws:iam::%s:role/NeptuneImportRole", acct)),
			Format:    types.FormatCsv,
			Tags:      createTags,
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
		if resp.GraphId == nil || *resp.GraphId == "" {
			_, _ = tc.client.CancelImportTask(tc.ctx, &neptunegraph.CancelImportTaskInput{
				TaskIdentifier: resp.TaskId,
			})
			return fmt.Errorf("expected non-empty graph ID from CreateGraphUsingImportTask")
		}
		graph, err := tc.getGraph(*resp.GraphId)
		if err != nil {
			return err
		}
		if graph.Arn == nil || *graph.Arn == "" {
			return fmt.Errorf("expected non-empty ARN for import-created graph")
		}
		tagResp, err := tc.client.ListTagsForResource(tc.ctx, &neptunegraph.ListTagsForResourceInput{
			ResourceArn: graph.Arn,
		})
		if err != nil {
			return err
		}
		for k, v := range createTags {
			if got, ok := tagResp.Tags[k]; !ok || got != v {
				return fmt.Errorf("import-created graph missing create-time tag %s=%s (got %q)", k, v, got)
			}
		}
		_, _ = tc.client.CancelImportTask(tc.ctx, &neptunegraph.CancelImportTaskInput{
			TaskIdentifier: resp.TaskId,
		})
		_, _ = tc.client.DeleteGraph(tc.ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: resp.GraphId,
			SkipSnapshot:    aws.Bool(true),
		})
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunegraphExportTaskTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	var exportTaskID string
	reg := tc.region
	acct := tc.accountID

	results = append(results, r.RunTest("neptunegraph", "StartExportTask", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.client.StartExportTask(tc.ctx, &neptunegraph.StartExportTaskInput{
			GraphIdentifier:  aws.String(tc.graphID),
			Destination:      aws.String("s3://test-bucket/export-data/"),
			KmsKeyIdentifier: aws.String(fmt.Sprintf("arn:aws:kms:%s:%s:key/12345678-1234-1234-1234-123456789012", reg, acct)),
			RoleArn:          aws.String(fmt.Sprintf("arn:aws:iam::%s:role/NeptuneExportRole", acct)),
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
		if err := requireID(exportTaskID, "export task ID"); err != nil {
			return err
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
		if err := requireID(exportTaskID, "export task ID"); err != nil {
			return err
		}
		tasks, err := paginate(func(next *string) ([]types.ExportTaskSummary, *string, error) {
			resp, err := tc.client.ListExportTasks(tc.ctx, &neptunegraph.ListExportTasksInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return resp.Tasks, resp.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, t := range tasks {
			if t.TaskId != nil && *t.TaskId == exportTaskID {
				return nil
			}
		}
		return fmt.Errorf("export task not found in ListExportTasks")
	}))

	results = append(results, r.RunTest("neptunegraph", "ListExportTasks_FilterByGraph", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
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
		if err := requireID(exportTaskID, "export task ID"); err != nil {
			return err
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

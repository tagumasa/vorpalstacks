package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"vorpalstacks-sdk-tests/config"
)

type disableHostPrefixFinalize struct{}

func (*disableHostPrefixFinalize) ID() string { return "DisableHostPrefixFinalize" }

func (*disableHostPrefixFinalize) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
	ctx = smithyhttp.SetHostnameImmutable(ctx, true)
	return next.HandleFinalize(ctx, in)
}

func newNeptuneGraphClient(r *TestRunner) (*neptunegraph.Client, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, err
	}
	return neptunegraph.NewFromConfig(cfg, func(o *neptunegraph.Options) {
		o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
			return stack.Finalize.Add(&disableHostPrefixFinalize{}, middleware.Before)
		})
	}), nil
}

func (r *TestRunner) RunNeptunegraphTests() []TestResult {
	var results []TestResult

	client, err := newNeptuneGraphClient(r)
	if err != nil {
		return append(results, SetupFailResult("neptunegraph", "Failed to create client: %v", err))
	}

	ctx := context.Background()
	tsNano := fmt.Sprintf("%d", time.Now().UnixNano())
	graphName := fmt.Sprintf("sdk-graph-%s", tsNano[len(tsNano)-8:])
	snapshotName := fmt.Sprintf("sdk-snap-%s", tsNano[len(tsNano)-8:])

	var graphID string
	var graphARN string

	// === Graph Lifecycle ===

	results = append(results, r.RunTest("neptunegraph", "CreateGraph", func() error {
		resp, err := client.CreateGraph(ctx, &neptunegraph.CreateGraphInput{
			GraphName:          aws.String(graphName),
			ProvisionedMemory:  aws.Int32(128),
			DeletionProtection: aws.Bool(false),
			PublicConnectivity: aws.Bool(false),
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "sdk-test",
			},
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id == "" {
			return fmt.Errorf("expected non-empty graph ID")
		}
		if resp.Name == nil || *resp.Name != graphName {
			return fmt.Errorf("expected graphName=%s, got %v", graphName, resp.Name)
		}
		if resp.Status != types.GraphStatusAvailable {
			return fmt.Errorf("expected status AVAILABLE, got %s", resp.Status)
		}
		graphID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID from CreateGraph")
		}
		resp, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != graphID {
			return fmt.Errorf("expected graphId=%s, got %v", graphID, resp.Id)
		}
		if resp.Name == nil || *resp.Name != graphName {
			return fmt.Errorf("expected name=%s, got %v", graphName, resp.Name)
		}
		if resp.ProvisionedMemory == nil || *resp.ProvisionedMemory != 128 {
			return fmt.Errorf("expected provisionedMemory=128, got %v", resp.ProvisionedMemory)
		}
		if resp.Arn == nil {
			return fmt.Errorf("expected non-empty ARN")
		}
		graphARN = *resp.Arn
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListGraphs", func() error {
		resp, err := client.ListGraphs(ctx, &neptunegraph.ListGraphsInput{})
		if err != nil {
			return err
		}
		if resp.Graphs == nil {
			return fmt.Errorf("expected non-nil Graphs list")
		}
		found := false
		for _, g := range resp.Graphs {
			if g.Id != nil && *g.Id == graphID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created graph not found in ListGraphs")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "UpdateGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.UpdateGraph(ctx, &neptunegraph.UpdateGraphInput{
			GraphIdentifier:    aws.String(graphID),
			ProvisionedMemory:  aws.Int32(256),
			DeletionProtection: aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "UpdateGraph_Verify", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err != nil {
			return err
		}
		if resp.ProvisionedMemory == nil || *resp.ProvisionedMemory != 256 {
			return fmt.Errorf("expected provisionedMemory=256 after update, got %v", resp.ProvisionedMemory)
		}
		if resp.DeletionProtection == nil || !*resp.DeletionProtection {
			return fmt.Errorf("expected deletionProtection=true after update, got %v", resp.DeletionProtection)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "StopGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.StopGraph(ctx, &neptunegraph.StopGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "StopGraph_Verify", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.GraphStatusStopped {
			return fmt.Errorf("expected status STOPPED, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "StartGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.StartGraph(ctx, &neptunegraph.StartGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "StartGraph_Verify", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.GraphStatusAvailable {
			return fmt.Errorf("expected status AVAILABLE after start, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ResetGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.ResetGraph(ctx, &neptunegraph.ResetGraphInput{
			GraphIdentifier: aws.String(graphID),
			SkipSnapshot:    aws.Bool(true),
		})
		return err
	}))

	// === Snapshots ===

	var snapshotID string

	results = append(results, r.RunTest("neptunegraph", "CreateGraphSnapshot", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.CreateGraphSnapshot(ctx, &neptunegraph.CreateGraphSnapshotInput{
			GraphIdentifier: aws.String(graphID),
			SnapshotName:    aws.String(snapshotName),
			Tags: map[string]string{
				"Type": "sdk-test",
			},
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id == "" {
			return fmt.Errorf("expected non-empty snapshot ID")
		}
		if resp.Name == nil || *resp.Name != snapshotName {
			return fmt.Errorf("expected snapshotName=%s, got %v", snapshotName, resp.Name)
		}
		snapshotID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetGraphSnapshot", func() error {
		if snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		resp, err := client.GetGraphSnapshot(ctx, &neptunegraph.GetGraphSnapshotInput{
			SnapshotIdentifier: aws.String(snapshotID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != snapshotID {
			return fmt.Errorf("expected snapshotId=%s, got %v", snapshotID, resp.Id)
		}
		if resp.Status == "" {
			return fmt.Errorf("expected non-nil snapshot status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListGraphSnapshots", func() error {
		resp, err := client.ListGraphSnapshots(ctx, &neptunegraph.ListGraphSnapshotsInput{})
		if err != nil {
			return err
		}
		if resp.GraphSnapshots == nil {
			return fmt.Errorf("expected non-nil GraphSnapshots list")
		}
		found := false
		for _, s := range resp.GraphSnapshots {
			if s.Id != nil && *s.Id == snapshotID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created snapshot not found in ListGraphSnapshots")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListGraphSnapshots_FilterByGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.ListGraphSnapshots(ctx, &neptunegraph.ListGraphSnapshotsInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err != nil {
			return err
		}
		if resp.GraphSnapshots == nil {
			return fmt.Errorf("expected non-nil GraphSnapshots list")
		}
		for _, s := range resp.GraphSnapshots {
			if s.Id != nil && *s.Id == snapshotID {
				return nil
			}
		}
		return fmt.Errorf("snapshot not found when filtering by graph")
	}))

	// === Private Graph Endpoints ===

	results = append(results, r.RunTest("neptunegraph", "CreatePrivateGraphEndpoint", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.CreatePrivateGraphEndpoint(ctx, &neptunegraph.CreatePrivateGraphEndpointInput{
			GraphIdentifier: aws.String(graphID),
			VpcId:           aws.String("vpc-test123"),
			SubnetIds:       []string{"subnet-aaa111", "subnet-bbb222"},
		})
		if err != nil {
			return err
		}
		if resp.VpcId == nil || *resp.VpcId != "vpc-test123" {
			return fmt.Errorf("expected vpcId=vpc-test123, got %v", resp.VpcId)
		}
		if resp.Status == "" {
			return fmt.Errorf("expected non-nil endpoint status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetPrivateGraphEndpoint", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.GetPrivateGraphEndpoint(ctx, &neptunegraph.GetPrivateGraphEndpointInput{
			GraphIdentifier: aws.String(graphID),
			VpcId:           aws.String("vpc-test123"),
		})
		if err != nil {
			return err
		}
		if resp.VpcId == nil || *resp.VpcId != "vpc-test123" {
			return fmt.Errorf("expected vpcId=vpc-test123, got %v", resp.VpcId)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListPrivateGraphEndpoints", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.ListPrivateGraphEndpoints(ctx, &neptunegraph.ListPrivateGraphEndpointsInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err != nil {
			return err
		}
		if resp.PrivateGraphEndpoints == nil {
			return fmt.Errorf("expected non-nil PrivateGraphEndpoints list")
		}
		if len(resp.PrivateGraphEndpoints) == 0 {
			return fmt.Errorf("expected at least one private endpoint")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "DeletePrivateGraphEndpoint", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.DeletePrivateGraphEndpoint(ctx, &neptunegraph.DeletePrivateGraphEndpointInput{
			GraphIdentifier: aws.String(graphID),
			VpcId:           aws.String("vpc-test123"),
		})
		return err
	}))

	// === Import Tasks ===

	var importTaskID string

	results = append(results, r.RunTest("neptunegraph", "StartImportTask", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.StartImportTask(ctx, &neptunegraph.StartImportTaskInput{
			GraphIdentifier: aws.String(graphID),
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
		resp, err := client.GetImportTask(ctx, &neptunegraph.GetImportTaskInput{
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
		resp, err := client.ListImportTasks(ctx, &neptunegraph.ListImportTasksInput{})
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
		_, err := client.CancelImportTask(ctx, &neptunegraph.CancelImportTaskInput{
			TaskIdentifier: aws.String(importTaskID),
		})
		return err
	}))

	// === Export Tasks ===

	var exportTaskID string

	results = append(results, r.RunTest("neptunegraph", "StartExportTask", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.StartExportTask(ctx, &neptunegraph.StartExportTaskInput{
			GraphIdentifier:  aws.String(graphID),
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
		resp, err := client.GetExportTask(ctx, &neptunegraph.GetExportTaskInput{
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
		resp, err := client.ListExportTasks(ctx, &neptunegraph.ListExportTasksInput{})
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
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.ListExportTasks(ctx, &neptunegraph.ListExportTasksInput{
			GraphIdentifier: aws.String(graphID),
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
		_, err := client.CancelExportTask(ctx, &neptunegraph.CancelExportTaskInput{
			TaskIdentifier: aws.String(exportTaskID),
		})
		return err
	}))

	// === Tags ===

	results = append(results, r.RunTest("neptunegraph", "TagResource", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if graphARN == "" {
			return fmt.Errorf("no graph ARN")
		}
		_, err := client.TagResource(ctx, &neptunegraph.TagResourceInput{
			ResourceArn: aws.String(graphARN),
			Tags: map[string]string{
				"ExtraTag":  "extra-value",
				"CreatedBy": "sdk-tests",
			},
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "ListTagsForResource", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if graphARN == "" {
			return fmt.Errorf("no graph ARN")
		}
		resp, err := client.ListTagsForResource(ctx, &neptunegraph.ListTagsForResourceInput{
			ResourceArn: aws.String(graphARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("expected non-nil Tags map")
		}
		if resp.Tags["Environment"] != "test" {
			return fmt.Errorf("expected tag Environment=test, got %v", resp.Tags["Environment"])
		}
		if resp.Tags["Owner"] != "sdk-test" {
			return fmt.Errorf("expected tag Owner=sdk-test, got %v", resp.Tags["Owner"])
		}
		if resp.Tags["ExtraTag"] != "extra-value" {
			return fmt.Errorf("expected tag ExtraTag=extra-value, got %v", resp.Tags["ExtraTag"])
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "UntagResource", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if graphARN == "" {
			return fmt.Errorf("no graph ARN")
		}
		_, err := client.UntagResource(ctx, &neptunegraph.UntagResourceInput{
			ResourceArn: aws.String(graphARN),
			TagKeys:     []string{"ExtraTag"},
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "ListTagsForResource_AfterUntag", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		if graphARN == "" {
			return fmt.Errorf("no graph ARN")
		}
		resp, err := client.ListTagsForResource(ctx, &neptunegraph.ListTagsForResourceInput{
			ResourceArn: aws.String(graphARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("expected non-nil Tags map")
		}
		if _, ok := resp.Tags["ExtraTag"]; ok {
			return fmt.Errorf("expected ExtraTag to be removed")
		}
		if resp.Tags["Environment"] != "test" {
			return fmt.Errorf("expected tag Environment=test still present, got %v", resp.Tags["Environment"])
		}
		return nil
	}))

	// === Data Plane ===

	results = append(results, r.RunTest("neptunegraph", "ExecuteQuery_BasicMatch", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.ExecuteQuery(ctx, &neptunegraph.ExecuteQueryInput{
			GraphIdentifier: aws.String(graphID),
			Language:        types.QueryLanguageOpenCypher,
			QueryString:     aws.String("MATCH (n) RETURN n LIMIT 1"),
		})
		if err != nil {
			return fmt.Errorf("ExecuteQuery failed: %v", err)
		}
		if resp == nil || resp.Payload == nil {
			return fmt.Errorf("expected payload from ExecuteQuery")
		}
		resp.Payload.Close()
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "CancelQuery_NotImplemented", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.CancelQuery(ctx, &neptunegraph.CancelQueryInput{
			GraphIdentifier: aws.String(graphID),
			QueryId:         aws.String("q-fake123456"),
		})
		if err := AssertErrorContains(err, "NotImplementedException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetGraphSummary", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.GetGraphSummary(ctx, &neptunegraph.GetGraphSummaryInput{
			GraphIdentifier: aws.String(graphID),
			Mode:            types.GraphSummaryModeBasic,
		})
		if err != nil {
			return err
		}
		if resp.GraphSummary == nil {
			return fmt.Errorf("expected non-nil GraphSummary")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListQueries", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := client.ListQueries(ctx, &neptunegraph.ListQueriesInput{
			GraphIdentifier: aws.String(graphID),
			MaxResults:      aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.Queries == nil {
			return fmt.Errorf("expected non-nil Queries list")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetQuery_NotFound", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.GetQuery(ctx, &neptunegraph.GetQueryInput{
			GraphIdentifier: aws.String(graphID),
			QueryId:         aws.String("q-nonexist00"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === CreateGraphUsingImportTask ===

	results = append(results, r.RunTest("neptunegraph", "CreateGraphUsingImportTask", func() error {
		importGraphName := fmt.Sprintf("sdk-impgraph-%s", tsNano[len(tsNano)-6:])
		resp, err := client.CreateGraphUsingImportTask(ctx, &neptunegraph.CreateGraphUsingImportTaskInput{
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
		_, _ = client.CancelImportTask(ctx, &neptunegraph.CancelImportTaskInput{
			TaskIdentifier: resp.TaskId,
		})
		if resp.GraphId != nil && *resp.GraphId != "" {
			_, _ = client.DeleteGraph(ctx, &neptunegraph.DeleteGraphInput{
				GraphIdentifier: resp.GraphId,
				SkipSnapshot:    aws.Bool(true),
			})
		}
		return nil
	}))

	// === RestoreGraphFromSnapshot ===

	var restoredGraphID string

	results = append(results, r.RunTest("neptunegraph", "RestoreGraphFromSnapshot", func() error {
		if snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		restoreName := fmt.Sprintf("sdk-restored-%s", tsNano[len(tsNano)-6:])
		resp, err := client.RestoreGraphFromSnapshot(ctx, &neptunegraph.RestoreGraphFromSnapshotInput{
			SnapshotIdentifier: aws.String(snapshotID),
			GraphName:          aws.String(restoreName),
			ProvisionedMemory:  aws.Int32(128),
			DeletionProtection: aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id == "" {
			return fmt.Errorf("expected non-empty restored graph ID")
		}
		if resp.Name == nil || *resp.Name != restoreName {
			return fmt.Errorf("expected name=%s, got %v", restoreName, resp.Name)
		}
		if resp.SourceSnapshotId == nil || *resp.SourceSnapshotId != snapshotID {
			return fmt.Errorf("expected sourceSnapshotId=%s, got %v", snapshotID, resp.SourceSnapshotId)
		}
		restoredGraphID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "RestoreGraphFromSnapshot_Verify", func() error {
		if restoredGraphID == "" {
			return fmt.Errorf("no restored graph ID")
		}
		resp, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(restoredGraphID),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.GraphStatusAvailable {
			return fmt.Errorf("expected status AVAILABLE for restored graph, got %s", resp.Status)
		}
		if resp.SourceSnapshotId == nil || *resp.SourceSnapshotId != snapshotID {
			return fmt.Errorf("expected sourceSnapshotId=%s, got %v", snapshotID, resp.SourceSnapshotId)
		}
		return nil
	}))

	if restoredGraphID != "" {
		_, _ = client.DeleteGraph(ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: aws.String(restoredGraphID),
			SkipSnapshot:    aws.Bool(true),
		})
	}

	// === Delete Snapshot ===

	results = append(results, r.RunTest("neptunegraph", "DeleteGraphSnapshot", func() error {
		if snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		_, err := client.DeleteGraphSnapshot(ctx, &neptunegraph.DeleteGraphSnapshotInput{
			SnapshotIdentifier: aws.String(snapshotID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraphSnapshot_Verify", func() error {
		if snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		_, err := client.GetGraphSnapshot(ctx, &neptunegraph.GetGraphSnapshotInput{
			SnapshotIdentifier: aws.String(snapshotID),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === Error Cases ===

	results = append(results, r.RunTest("neptunegraph", "GetGraph_NotFound", func() error {
		_, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String("g-nonexist00"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraph_NotFound", func() error {
		_, err := client.DeleteGraph(ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: aws.String("g-nonexist00"),
			SkipSnapshot:    aws.Bool(true),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// === Cleanup ===

	results = append(results, r.RunTest("neptunegraph", "UpdateGraph_DisableDeletionProtection", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.UpdateGraph(ctx, &neptunegraph.UpdateGraphInput{
			GraphIdentifier:    aws.String(graphID),
			DeletionProtection: aws.Bool(false),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraph", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.DeleteGraph(ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: aws.String(graphID),
			SkipSnapshot:    aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraph_Verify", func() error {
		if graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := client.GetGraph(ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(graphID),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

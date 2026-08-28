package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func (r *TestRunner) runNeptunegraphGraphTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	graphName := tc.unique("sdk-graph")

	results = append(results, r.RunTest("neptunegraph", "CreateGraph", func() error {
		resp, err := tc.client.CreateGraph(tc.ctx, &neptunegraph.CreateGraphInput{
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
		if resp.ProvisionedMemory == nil || *resp.ProvisionedMemory != 128 {
			return fmt.Errorf("expected provisionedMemory=128, got %v", resp.ProvisionedMemory)
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("expected non-empty ARN")
		}
		tc.graphID = *resp.Id
		tc.graphARN = *resp.Arn
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetGraph", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.getGraph(tc.graphID)
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected graphId=%s, got %v", tc.graphID, resp.Id)
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
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListGraphs", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		graphs, err := paginate(func(next *string) ([]types.GraphSummary, *string, error) {
			resp, err := tc.client.ListGraphs(tc.ctx, &neptunegraph.ListGraphsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return resp.Graphs, resp.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, g := range graphs {
			if g.Id != nil && *g.Id == tc.graphID {
				return nil
			}
		}
		return fmt.Errorf("created graph not found in ListGraphs")
	}))

	results = append(results, r.RunTest("neptunegraph", "UpdateGraph", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.client.UpdateGraph(tc.ctx, &neptunegraph.UpdateGraphInput{
			GraphIdentifier:    aws.String(tc.graphID),
			ProvisionedMemory:  aws.Int32(256),
			DeletionProtection: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected id=%s, got %v", tc.graphID, resp.Id)
		}
		if resp.Status == "" {
			return fmt.Errorf("expected non-empty status")
		}
		if resp.ProvisionedMemory == nil || *resp.ProvisionedMemory != 256 {
			return fmt.Errorf("expected provisionedMemory=256, got %v", resp.ProvisionedMemory)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "UpdateGraph_Verify", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.getGraph(tc.graphID)
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
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.client.StopGraph(tc.ctx, &neptunegraph.StopGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected id=%s, got %v", tc.graphID, resp.Id)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "StopGraph_Verify", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.getGraph(tc.graphID)
		if err != nil {
			return err
		}
		if resp.Status != types.GraphStatusStopped {
			return fmt.Errorf("expected status STOPPED, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "StartGraph", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.client.StartGraph(tc.ctx, &neptunegraph.StartGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected id=%s, got %v", tc.graphID, resp.Id)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "StartGraph_Verify", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.getGraph(tc.graphID)
		if err != nil {
			return err
		}
		if resp.Status != types.GraphStatusAvailable {
			return fmt.Errorf("expected status AVAILABLE after start, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ResetGraph", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		resp, err := tc.client.ResetGraph(tc.ctx, &neptunegraph.ResetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
			SkipSnapshot:    aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected id=%s, got %v", tc.graphID, resp.Id)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ResetGraph_FinalSnapshotCreated", func() error {
		if err := tc.requireGraph(); err != nil {
			return err
		}
		snapshotsFor := func() ([]types.GraphSnapshotSummary, error) {
			return paginate(func(next *string) ([]types.GraphSnapshotSummary, *string, error) {
				resp, err := tc.client.ListGraphSnapshots(tc.ctx, &neptunegraph.ListGraphSnapshotsInput{
					GraphIdentifier: aws.String(tc.graphID),
					NextToken:       next,
				})
				if err != nil {
					return nil, nil, err
				}
				return resp.GraphSnapshots, resp.NextToken, nil
			})
		}
		before, err := snapshotsFor()
		if err != nil {
			return err
		}
		resp, err := tc.client.ResetGraph(tc.ctx, &neptunegraph.ResetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
			SkipSnapshot:    aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected id=%s, got %v", tc.graphID, resp.Id)
		}
		after, err := snapshotsFor()
		if err != nil {
			return err
		}
		if len(after) <= len(before) {
			return fmt.Errorf("expected a final snapshot after ResetGraph(skipSnapshot=false), before=%d after=%d", len(before), len(after))
		}
		for _, sn := range after {
			if sn.SourceGraphId != nil && *sn.SourceGraphId == tc.graphID {
				if sn.Status != types.SnapshotStatusAvailable && sn.Status != types.SnapshotStatusCreating {
					return fmt.Errorf("unexpected final snapshot status %s", sn.Status)
				}
				return nil
			}
		}
		return fmt.Errorf("final snapshot for graph %s not found in list", tc.graphID)
	}))

	return results
}

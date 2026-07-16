package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func (r *TestRunner) runNeptunegraphGraphTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	graphName := fmt.Sprintf("sdk-graph-%s", tc.tsNano[len(tc.tsNano)-8:])

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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID from CreateGraph")
		}
		resp, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
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
		resp, err := tc.client.ListGraphs(tc.ctx, &neptunegraph.ListGraphsInput{})
		if err != nil {
			return err
		}
		if resp.Graphs == nil {
			return fmt.Errorf("expected non-nil Graphs list")
		}
		found := false
		for _, g := range resp.Graphs {
			if g.Id != nil && *g.Id == tc.graphID {
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
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

	return results
}

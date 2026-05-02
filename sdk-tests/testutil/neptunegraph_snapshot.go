package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func (r *TestRunner) runNeptunegraphSnapshotTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult
	snapshotName := fmt.Sprintf("sdk-snap-%s", tc.tsNano[len(tc.tsNano)-8:])

	results = append(results, r.RunTest("neptunegraph", "CreateGraphSnapshot", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.CreateGraphSnapshot(tc.ctx, &neptunegraph.CreateGraphSnapshotInput{
			GraphIdentifier: aws.String(tc.graphID),
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
		tc.snapshotID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetGraphSnapshot", func() error {
		if tc.snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		resp, err := tc.client.GetGraphSnapshot(tc.ctx, &neptunegraph.GetGraphSnapshotInput{
			SnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.snapshotID {
			return fmt.Errorf("expected snapshotId=%s, got %v", tc.snapshotID, resp.Id)
		}
		if resp.Status == "" {
			return fmt.Errorf("expected non-nil snapshot status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListGraphSnapshots", func() error {
		resp, err := tc.client.ListGraphSnapshots(tc.ctx, &neptunegraph.ListGraphSnapshotsInput{})
		if err != nil {
			return err
		}
		if resp.GraphSnapshots == nil {
			return fmt.Errorf("expected non-nil GraphSnapshots list")
		}
		found := false
		for _, s := range resp.GraphSnapshots {
			if s.Id != nil && *s.Id == tc.snapshotID {
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.ListGraphSnapshots(tc.ctx, &neptunegraph.ListGraphSnapshotsInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
		if err != nil {
			return err
		}
		if resp.GraphSnapshots == nil {
			return fmt.Errorf("expected non-nil GraphSnapshots list")
		}
		for _, s := range resp.GraphSnapshots {
			if s.Id != nil && *s.Id == tc.snapshotID {
				return nil
			}
		}
		return fmt.Errorf("snapshot not found when filtering by graph")
	}))

	results = append(results, r.runNeptunegraphRestoreTests(tc)...)

	results = append(results, r.RunTest("neptunegraph", "DeleteGraphSnapshot", func() error {
		if tc.snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		resp, err := tc.client.DeleteGraphSnapshot(tc.ctx, &neptunegraph.DeleteGraphSnapshotInput{
			SnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.snapshotID {
			return fmt.Errorf("expected snapshotId=%s, got %v", tc.snapshotID, resp.Id)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraphSnapshot_Verify", func() error {
		if tc.snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		_, err := tc.client.GetGraphSnapshot(tc.ctx, &neptunegraph.GetGraphSnapshotInput{
			SnapshotIdentifier: aws.String(tc.snapshotID),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunegraphRestoreTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult

	var restoredGraphID string
	results = append(results, r.RunTest("neptunegraph", "RestoreGraphFromSnapshot", func() error {
		if tc.snapshotID == "" {
			return fmt.Errorf("no snapshot ID")
		}
		restoreName := fmt.Sprintf("sdk-restored-%s", tc.tsNano[len(tc.tsNano)-6:])
		resp, err := tc.client.RestoreGraphFromSnapshot(tc.ctx, &neptunegraph.RestoreGraphFromSnapshotInput{
			SnapshotIdentifier: aws.String(tc.snapshotID),
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
		if resp.SourceSnapshotId == nil || *resp.SourceSnapshotId != tc.snapshotID {
			return fmt.Errorf("expected sourceSnapshotId=%s, got %v", tc.snapshotID, resp.SourceSnapshotId)
		}
		restoredGraphID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "RestoreGraphFromSnapshot_Verify", func() error {
		if restoredGraphID == "" {
			return fmt.Errorf("no restored graph ID")
		}
		resp, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(restoredGraphID),
		})
		if err != nil {
			return err
		}
		if resp.Status != types.GraphStatusAvailable {
			return fmt.Errorf("expected status AVAILABLE for restored graph, got %s", resp.Status)
		}
		if resp.SourceSnapshotId == nil || *resp.SourceSnapshotId != tc.snapshotID {
			return fmt.Errorf("expected sourceSnapshotId=%s, got %v", tc.snapshotID, resp.SourceSnapshotId)
		}
		return nil
	}))

	if restoredGraphID != "" {
		_, _ = tc.client.DeleteGraph(tc.ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: aws.String(restoredGraphID),
			SkipSnapshot:    aws.Bool(true),
		})
	}

	return results
}

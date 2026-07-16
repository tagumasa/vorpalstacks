package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
)

func (r *TestRunner) runNeptunegraphEdgeTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunegraph", "GetGraph_NotFound", func() error {
		_, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String("g-nonexist00"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraph_NotFound", func() error {
		_, err := tc.client.DeleteGraph(tc.ctx, &neptunegraph.DeleteGraphInput{
			GraphIdentifier: aws.String("g-nonexist00"),
			SkipSnapshot:    aws.Bool(true),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "UpdateGraph_DisableDeletionProtection", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.UpdateGraph(tc.ctx, &neptunegraph.UpdateGraphInput{
			GraphIdentifier:    aws.String(tc.graphID),
			DeletionProtection: aws.Bool(false),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != tc.graphID {
			return fmt.Errorf("expected id=%s, got %v", tc.graphID, resp.Id)
		}
		if resp.DeletionProtection == nil || *resp.DeletionProtection {
			return fmt.Errorf("expected deletionProtection=false, got %v", resp.DeletionProtection)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "DeleteGraph", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.DeleteGraph(tc.ctx, &neptunegraph.DeleteGraphInput{
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

	results = append(results, r.RunTest("neptunegraph", "DeleteGraph_Verify", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := tc.client.GetGraph(tc.ctx, &neptunegraph.GetGraphInput{
			GraphIdentifier: aws.String(tc.graphID),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

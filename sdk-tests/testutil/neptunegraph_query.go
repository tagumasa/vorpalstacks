package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func (r *TestRunner) runNeptunegraphQueryTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunegraph", "ExecuteQuery_BasicMatch", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.ExecuteQuery(tc.ctx, &neptunegraph.ExecuteQueryInput{
			GraphIdentifier: aws.String(tc.graphID),
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

	results = append(results, r.RunTest("neptunegraph", "CancelQuery_NonExistentQuery", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := tc.client.CancelQuery(tc.ctx, &neptunegraph.CancelQueryInput{
			GraphIdentifier: aws.String(tc.graphID),
			QueryId:         aws.String("q-fake123456"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "GetGraphSummary", func() error {
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.GetGraphSummary(tc.ctx, &neptunegraph.GetGraphSummaryInput{
			GraphIdentifier: aws.String(tc.graphID),
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		resp, err := tc.client.ListQueries(tc.ctx, &neptunegraph.ListQueriesInput{
			GraphIdentifier: aws.String(tc.graphID),
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
		if tc.graphID == "" {
			return fmt.Errorf("no graph ID")
		}
		_, err := tc.client.GetQuery(tc.ctx, &neptunegraph.GetQueryInput{
			GraphIdentifier: aws.String(tc.graphID),
			QueryId:         aws.String("q-nonexist00"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

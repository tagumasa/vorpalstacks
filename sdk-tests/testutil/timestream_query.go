package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
)

func (r *TestRunner) runTimestreamQueryTests(tc *tsTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("timestream", "DescribeEndpoints_Query", func() error {
		resp, err := tc.queryClient.DescribeEndpoints(tc.ctx, &timestreamquery.DescribeEndpointsInput{})
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

	results = append(results, r.RunTest("timestream", "PrepareQuery_Verify", func() error {
		resp, err := tc.queryClient.PrepareQuery(tc.ctx, &timestreamquery.PrepareQueryInput{
			QueryString: aws.String("SELECT 1"),
		})
		if err != nil {
			return err
		}
		if resp.QueryString == nil || *resp.QueryString != "SELECT 1" {
			return fmt.Errorf("expected QueryString='SELECT 1', got %v", resp.QueryString)
		}
		if len(resp.Columns) == 0 {
			return fmt.Errorf("expected non-empty Columns")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "CancelQuery", func() error {
		queryResp, err := tc.queryClient.Query(tc.ctx, &timestreamquery.QueryInput{
			QueryString: aws.String("SELECT 1"),
		})
		if err != nil {
			return fmt.Errorf("query: %v", err)
		}
		if queryResp.QueryId == nil {
			return fmt.Errorf("QueryId is nil")
		}
		cancelResp, err := tc.queryClient.CancelQuery(tc.ctx, &timestreamquery.CancelQueryInput{
			QueryId: queryResp.QueryId,
		})
		if err != nil {
			return err
		}
		if cancelResp.CancellationMessage == nil {
			return fmt.Errorf("CancellationMessage is nil")
		}
		return nil
	}))

	return results
}

package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestContext) testQueryResults() []TestResult {
	var results []TestResult

	// Shared query fixture: start once, wait for completion, then both
	// result scenarios below read from it. A failure surfaces as one FAIL
	// row named after the setup step it replaced.
	resultsQueryId, err := tc.startAndWaitForQuery("SHOW DATABASES")
	if err != nil {
		return append(results, TestResult{
			Service:  "athena",
			TestName: "GetQueryResults_StartQuery",
			Status:   "FAIL",
			Error:    fmt.Sprintf("query setup failed: %v", err),
		})
	}

	results = append(results, tc.runner.RunTest("athena", "GetQueryResults", func() error {
		resp, err := tc.client.GetQueryResults(tc.ctx, &athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(resultsQueryId),
		})
		if err != nil {
			return err
		}
		rs := resp.ResultSet
		if rs == nil {
			return fmt.Errorf("result set is nil")
		}
		if rs.ResultSetMetadata == nil {
			return fmt.Errorf("result set metadata is nil")
		}
		if len(rs.Rows) == 0 {
			return fmt.Errorf("result set has no rows")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetQueryRuntimeStatistics", func() error {
		resp, err := tc.client.GetQueryRuntimeStatistics(tc.ctx, &athena.GetQueryRuntimeStatisticsInput{
			QueryExecutionId: aws.String(resultsQueryId),
		})
		if err != nil {
			return err
		}
		if resp.QueryRuntimeStatistics == nil {
			return fmt.Errorf("query runtime statistics is nil")
		}
		return nil
	}))

	return results
}

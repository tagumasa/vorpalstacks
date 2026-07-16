package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestCtx) testQueryResults() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	var resultsQueryId string
	results = append(results, tc.runner.RunTest("athena", "GetQueryResults_StartQuery", func() error {
		resp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SHOW DATABASES"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
		})
		if err != nil {
			return err
		}
		resultsQueryId = aws.ToString(resp.QueryExecutionId)
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetQueryResults_WaitForCompletion", func() error {
		for i := 0; i < 30; i++ {
			resp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
				QueryExecutionId: aws.String(resultsQueryId),
			})
			if err != nil {
				return err
			}
			state := resp.QueryExecution.Status.State
			if state == types.QueryExecutionStateSucceeded {
				return nil
			}
			if state == types.QueryExecutionStateFailed || state == types.QueryExecutionStateCancelled {
				return fmt.Errorf("query ended in state %s", state)
			}
			time.Sleep(500 * time.Millisecond)
		}
		return fmt.Errorf("query did not complete within timeout")
	}))

	results = append(results, tc.runner.RunTest("athena", "GetQueryResults", func() error {
		resp, err := client.GetQueryResults(ctx, &athena.GetQueryResultsInput{
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
		resp, err := client.GetQueryRuntimeStatistics(ctx, &athena.GetQueryRuntimeStatisticsInput{
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

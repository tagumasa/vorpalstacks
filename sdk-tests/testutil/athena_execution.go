package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestCtx) testQueryExecution() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	var queryExecutionId string
	results = append(results, tc.runner.RunTest("athena", "StartQueryExecution", func() error {
		resp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SELECT 1"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
			ResultConfiguration: &types.ResultConfiguration{
				OutputLocation: aws.String("s3://test-bucket/athena/"),
			},
		})
		if err != nil {
			return err
		}
		queryExecutionId = aws.ToString(resp.QueryExecutionId)
		if queryExecutionId == "" {
			return fmt.Errorf("QueryExecutionId is empty")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetQueryExecution", func() error {
		resp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: aws.String(queryExecutionId),
		})
		if err != nil {
			return err
		}
		qe := resp.QueryExecution
		if qe == nil {
			return fmt.Errorf("query execution is nil")
		}
		if aws.ToString(qe.Query) != "SELECT 1" {
			return fmt.Errorf("expected query 'SELECT 1', got %q", aws.ToString(qe.Query))
		}
		if qe.Status == nil {
			return fmt.Errorf("status is nil")
		}
		if qe.QueryExecutionContext == nil || aws.ToString(qe.QueryExecutionContext.Database) != "default" {
			return fmt.Errorf("expected query context database 'default'")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListQueryExecutions", func() error {
		resp, err := client.ListQueryExecutions(ctx, &athena.ListQueryExecutionsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.QueryExecutionIds == nil {
			return fmt.Errorf("query execution IDs list is nil")
		}
		var found bool
		for _, id := range resp.QueryExecutionIds {
			if id == queryExecutionId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("started query execution ID %q not found in list", queryExecutionId)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "StopQueryExecution", func() error {
		resp, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("/* SLOW */ SELECT 1"),
		})
		if err != nil {
			return err
		}
		stopQueryId := resp.QueryExecutionId

		_, err = client.StopQueryExecution(ctx, &athena.StopQueryExecutionInput{
			QueryExecutionId: stopQueryId,
		})
		if err != nil {
			return fmt.Errorf("StopQueryExecution failed: %v", err)
		}

		getResp, err := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: stopQueryId,
		})
		if err != nil {
			return err
		}
		if getResp.QueryExecution.Status.State != types.QueryExecutionStateCancelled {
			return fmt.Errorf("expected CANCELLED, got %s", getResp.QueryExecution.Status.State)
		}
		return nil
	}))

	var batchQEId1, batchQEId2 string
	results = append(results, tc.runner.RunTest("athena", "BatchGetQueryExecution_Setup", func() error {
		resp1, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SELECT 1"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
		})
		if err != nil {
			return err
		}
		batchQEId1 = aws.ToString(resp1.QueryExecutionId)

		resp2, err := client.StartQueryExecution(ctx, &athena.StartQueryExecutionInput{
			QueryString: aws.String("SELECT 2"),
			QueryExecutionContext: &types.QueryExecutionContext{
				Database: aws.String("default"),
			},
		})
		if err != nil {
			return err
		}
		batchQEId2 = aws.ToString(resp2.QueryExecutionId)

		for i := 0; i < 30; i++ {
			r1, _ := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{QueryExecutionId: aws.String(batchQEId1)})
			r2, _ := client.GetQueryExecution(ctx, &athena.GetQueryExecutionInput{QueryExecutionId: aws.String(batchQEId2)})
			if r1.QueryExecution.Status.State == types.QueryExecutionStateSucceeded &&
				r2.QueryExecution.Status.State == types.QueryExecutionStateSucceeded {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
		return fmt.Errorf("queries did not complete within timeout")
	}))

	results = append(results, tc.runner.RunTest("athena", "BatchGetQueryExecution", func() error {
		resp, err := client.BatchGetQueryExecution(ctx, &athena.BatchGetQueryExecutionInput{
			QueryExecutionIds: []string{batchQEId1, batchQEId2, "nonexistent-qe-id"},
		})
		if err != nil {
			return err
		}
		if len(resp.QueryExecutions) != 2 {
			return fmt.Errorf("expected 2 query executions, got %d", len(resp.QueryExecutions))
		}
		if len(resp.UnprocessedQueryExecutionIds) != 1 {
			return fmt.Errorf("expected 1 unprocessed, got %d", len(resp.UnprocessedQueryExecutionIds))
		}
		return nil
	}))

	return results
}

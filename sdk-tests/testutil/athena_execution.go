package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestContext) testQueryExecution() []TestResult {
	var results []TestResult

	var queryExecutionId string
	results = append(results, tc.runner.RunTest("athena", "StartQueryExecution", func() error {
		id, err := tc.startQuery("SELECT 1")
		if err != nil {
			return err
		}
		queryExecutionId = id
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetQueryExecution", func() error {
		resp, err := tc.client.GetQueryExecution(tc.ctx, &athena.GetQueryExecutionInput{
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
		allIDs, err := tc.allQueryExecutions()
		if err != nil {
			return err
		}
		for _, id := range allIDs {
			if id == queryExecutionId {
				return nil
			}
		}
		return fmt.Errorf("started query execution ID %q not found in list", queryExecutionId)
	}))

	results = append(results, tc.runner.RunTest("athena", "StopQueryExecution", func() error {
		stopQueryId, err := tc.startQuery("/* SLOW */ SELECT 1")
		if err != nil {
			return err
		}

		_, err = tc.client.StopQueryExecution(tc.ctx, &athena.StopQueryExecutionInput{
			QueryExecutionId: aws.String(stopQueryId),
		})
		if err != nil {
			return fmt.Errorf("StopQueryExecution failed: %v", err)
		}

		getResp, err := tc.client.GetQueryExecution(tc.ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: aws.String(stopQueryId),
		})
		if err != nil {
			return err
		}
		if getResp.QueryExecution.Status.State != types.QueryExecutionStateCancelled {
			return fmt.Errorf("expected CANCELLED, got %s", getResp.QueryExecution.Status.State)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "BatchGetQueryExecution", func() error {
		batchQEId1, err := tc.startQuery("SELECT 1")
		if err != nil {
			return fmt.Errorf("setup start 1 failed: %w", err)
		}
		batchQEId2, err := tc.startQuery("SELECT 2")
		if err != nil {
			return fmt.Errorf("setup start 2 failed: %w", err)
		}

		if _, err := tc.waitForQuery(batchQEId1); err != nil {
			return fmt.Errorf("setup query 1 did not complete: %w", err)
		}
		if _, err := tc.waitForQuery(batchQEId2); err != nil {
			return fmt.Errorf("setup query 2 did not complete: %w", err)
		}

		resp, err := tc.client.BatchGetQueryExecution(tc.ctx, &athena.BatchGetQueryExecutionInput{
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

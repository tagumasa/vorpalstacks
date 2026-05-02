package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
)

func (r *TestRunner) runNeptunedataServerAPITests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "ListGremlinQueries", func() error {
		resp, err := tc.client.ListGremlinQueries(tc.ctx, &neptunedata.ListGremlinQueriesInput{})
		if err != nil {
			return err
		}
		if resp.AcceptedQueryCount == nil {
			return fmt.Errorf("expected non-nil AcceptedQueryCount")
		}
		if resp.RunningQueryCount == nil {
			return fmt.Errorf("expected non-nil RunningQueryCount")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ListOpenCypherQueries", func() error {
		resp, err := tc.client.ListOpenCypherQueries(tc.ctx, &neptunedata.ListOpenCypherQueriesInput{})
		if err != nil {
			return err
		}
		if resp.AcceptedQueryCount == nil {
			return fmt.Errorf("expected non-nil AcceptedQueryCount")
		}
		if resp.RunningQueryCount == nil {
			return fmt.Errorf("expected non-nil RunningQueryCount")
		}
		return nil
	}))

	// Cancel tests use nonexistent query IDs because vorpalstacks executes
	// queries synchronously — by the time a cancel request reaches the server
	// the query has already completed.  Only error-path coverage is feasible
	// without an artificially slow query execution environment.
	results = append(results, r.RunTest("neptunedata", "CancelGremlinQuery", func() error {
		_, err := tc.client.CancelGremlinQuery(tc.ctx, &neptunedata.CancelGremlinQueryInput{
			QueryId: aws.String("nonexistent-query-id"),
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return err
		}
		return nil
	}))

	// See note above regarding synchronous execution and cancel timing.
	results = append(results, r.RunTest("neptunedata", "CancelOpenCypherQuery", func() error {
		_, err := tc.client.CancelOpenCypherQuery(tc.ctx, &neptunedata.CancelOpenCypherQueryInput{
			QueryId: aws.String("nonexistent-query-id"),
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return err
		}
		return nil
	}))

	// GetOpenCypherQueryStatus cannot be tested with a real query ID because
	// vorpalstacks executes synchronously and ListOpenCypherQueries only
	// returns running/waiting queries per AWS spec.  The query completes
	// before the list call returns, so no queryId can be obtained.
	// Verify that nonexistent query returns proper error.
	results = append(results, r.RunTest("neptunedata", "GetOpenCypherQueryStatus_NotFound", func() error {
		_, err := tc.client.GetOpenCypherQueryStatus(tc.ctx, &neptunedata.GetOpenCypherQueryStatusInput{
			QueryId: aws.String("nonexistent-query-id"),
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunedataUnsupportedTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "GetSparqlStatistics_Unsupported", func() error {
		_, err := tc.client.GetSparqlStatistics(tc.ctx, &neptunedata.GetSparqlStatisticsInput{})
		if err := AssertErrorContains(err, "UnsupportedOperationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetRDFGraphSummary_Unsupported", func() error {
		_, err := tc.client.GetRDFGraphSummary(tc.ctx, &neptunedata.GetRDFGraphSummaryInput{})
		if err := AssertErrorContains(err, "UnsupportedOperationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "StartMLDataProcessingJob_Unsupported", func() error {
		_, err := tc.client.StartMLDataProcessingJob(tc.ctx, &neptunedata.StartMLDataProcessingJobInput{
			InputDataS3Location:     aws.String("s3://test/ml-input"),
			ProcessedDataS3Location: aws.String("s3://test/ml-output"),
		})
		if err := AssertErrorContains(err, "UnsupportedOperationException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunedataEdgeTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "ErrorCase_InvalidCypherSyntax", func() error {
		_, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("INVALID CYPHER QUERY"),
		})
		if err := AssertErrorContains(err, "MalformedQueryException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_InvalidGremlinSyntax", func() error {
		_, err := tc.client.ExecuteGremlinQuery(tc.ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.INVALID_STEP()"),
		})
		if err := AssertErrorContains(err, "MalformedQueryException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_FastResetInvalidToken", func() error {
		_, err := tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  aws.String("invalid-token-12345"),
		})
		if err := AssertErrorContains(err, "PreconditionsFailedException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_NonExistentLoaderJob", func() error {
		_, err := tc.client.GetLoaderJobStatus(tc.ctx, &neptunedata.GetLoaderJobStatusInput{
			LoadId: aws.String("nonexistent-load-id"),
		})
		if err := AssertErrorContains(err, "BulkLoadIdNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_CancelNonExistentLoaderJob", func() error {
		_, err := tc.client.CancelLoaderJob(tc.ctx, &neptunedata.CancelLoaderJobInput{
			LoadId: aws.String("nonexistent-load-id"),
		})
		if err := AssertErrorContains(err, "BulkLoadIdNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_EmptyCypherQuery", func() error {
		_, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String(""),
		})
		if err := AssertErrorContains(err, "MissingParameterException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_EmptyGremlinQuery", func() error {
		_, err := tc.client.ExecuteGremlinQuery(tc.ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String(""),
		})
		if err := AssertErrorContains(err, "MissingParameterException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

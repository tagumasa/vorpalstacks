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

	// Cancel and status tests use nonexistent query IDs because vorpalstacks
	// executes queries synchronously — by the time a cancel or status request
	// reaches the server the query has already completed, and
	// ListOpenCypherQueries only returns running/waiting queries per AWS
	// spec, so no real queryId can be obtained.  Only error-path coverage is
	// feasible without an artificially slow query execution environment.
	results = append(results, r.RunTest("neptunedata", "QueryCancelStatusErrorPaths", func() error {
		for _, c := range []struct {
			name string
			call func() error
		}{
			{"CancelGremlinQuery", func() error {
				_, err := tc.client.CancelGremlinQuery(tc.ctx, &neptunedata.CancelGremlinQueryInput{
					QueryId: aws.String("nonexistent-query-id"),
				})
				return err
			}},
			{"CancelOpenCypherQuery", func() error {
				_, err := tc.client.CancelOpenCypherQuery(tc.ctx, &neptunedata.CancelOpenCypherQueryInput{
					QueryId: aws.String("nonexistent-query-id"),
				})
				return err
			}},
			{"GetOpenCypherQueryStatus", func() error {
				_, err := tc.client.GetOpenCypherQueryStatus(tc.ctx, &neptunedata.GetOpenCypherQueryStatusInput{
					QueryId: aws.String("nonexistent-query-id"),
				})
				return err
			}},
		} {
			if err := expectAWSErrorCode(c.call(), "BadRequestException"); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunedataUnsupportedTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	// Each row pins the UnsupportedOperationException contract of one RDF or
	// ML surface that has no substrate on this platform.
	results = append(results, r.RunTest("neptunedata", "UnsupportedSurfaces", func() error {
		for _, c := range []struct {
			name string
			call func() error
		}{
			{"GetSparqlStatistics", func() error {
				_, err := tc.client.GetSparqlStatistics(tc.ctx, &neptunedata.GetSparqlStatisticsInput{})
				return err
			}},
			{"GetRDFGraphSummary", func() error {
				_, err := tc.client.GetRDFGraphSummary(tc.ctx, &neptunedata.GetRDFGraphSummaryInput{})
				return err
			}},
			{"StartMLDataProcessingJob", func() error {
				_, err := tc.client.StartMLDataProcessingJob(tc.ctx, &neptunedata.StartMLDataProcessingJobInput{
					InputDataS3Location:     aws.String("s3://test/ml-input"),
					ProcessedDataS3Location: aws.String("s3://test/ml-output"),
				})
				return err
			}},
		} {
			if err := expectAWSErrorCode(c.call(), "UnsupportedOperationException"); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runNeptunedataEdgeTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	// Each row pins the error contract of one query/loader API against a
	// malformed, empty, or nonexistent request.
	results = append(results, r.RunTest("neptunedata", "ErrorCases", func() error {
		for _, c := range []struct {
			name string
			call func() error
			code string
		}{
			{"InvalidCypherSyntax", func() error {
				_, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
					OpenCypherQuery: aws.String("INVALID CYPHER QUERY"),
				})
				return err
			}, "MalformedQueryException"},
			{"InvalidGremlinSyntax", func() error {
				_, err := tc.client.ExecuteGremlinQuery(tc.ctx, &neptunedata.ExecuteGremlinQueryInput{
					GremlinQuery: aws.String("g.INVALID_STEP()"),
				})
				return err
			}, "MalformedQueryException"},
			{"FastResetInvalidToken", func() error {
				_, err := tc.client.ExecuteFastReset(tc.ctx, &neptunedata.ExecuteFastResetInput{
					Action: types.ActionPerformReset,
					Token:  aws.String("invalid-token-12345"),
				})
				return err
			}, "PreconditionsFailedException"},
			{"NonExistentLoaderJob", func() error {
				_, err := tc.client.GetLoaderJobStatus(tc.ctx, &neptunedata.GetLoaderJobStatusInput{
					LoadId: aws.String("nonexistent-load-id"),
				})
				return err
			}, "BulkLoadIdNotFoundException"},
			{"CancelNonExistentLoaderJob", func() error {
				_, err := tc.client.CancelLoaderJob(tc.ctx, &neptunedata.CancelLoaderJobInput{
					LoadId: aws.String("nonexistent-load-id"),
				})
				return err
			}, "BulkLoadIdNotFoundException"},
			{"EmptyCypherQuery", func() error {
				_, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
					OpenCypherQuery: aws.String(""),
				})
				return err
			}, "MissingParameterException"},
			{"EmptyGremlinQuery", func() error {
				_, err := tc.client.ExecuteGremlinQuery(tc.ctx, &neptunedata.ExecuteGremlinQueryInput{
					GremlinQuery: aws.String(""),
				})
				return err
			}, "MissingParameterException"},
		} {
			if err := expectAWSErrorCode(c.call(), c.code); err != nil {
				return fmt.Errorf("%s: %w", c.name, err)
			}
		}
		return nil
	}))

	return results
}

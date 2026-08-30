package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"

	"vorpalstacks-sdk-tests/config"
)

type athenaTestContext struct {
	client *athena.Client
	ctx    context.Context
	runner *TestRunner
}

func (r *TestRunner) RunAthenaTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{
			Service:  "athena",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}

	tc := &athenaTestContext{
		client: athena.NewFromConfig(cfg),
		ctx:    context.Background(),
		runner: r,
	}

	results = append(results, tc.testWorkGroups()...)
	results = append(results, tc.testDataCatalogs()...)
	results = append(results, tc.testDatabases()...)
	results = append(results, tc.testNamedQueries()...)
	results = append(results, tc.testQueryExecution()...)
	results = append(results, tc.testQueryResults()...)
	results = append(results, tc.testPreparedStatements()...)
	results = append(results, tc.testCapacityReservations()...)
	results = append(results, tc.testTagging()...)
	results = append(results, tc.testSQL()...)
	results = append(results, tc.testValidation()...)
	results = append(results, tc.testEdgeCases()...)

	return results
}

// uniqueName builds a per-run unique resource name from the given prefix.
func (tc *athenaTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()%1000000000)
}

// workgroupARN builds the Athena work group ARN for this runner's region
// and account.
func (tc *athenaTestContext) workgroupARN(name string) string {
	return fmt.Sprintf("arn:aws:athena:%s:%s:workgroup/%s", tc.runner.region, tc.runner.AccountID(), name)
}

// datacatalogARN builds the Athena data catalog ARN for this runner's
// region and account.
func (tc *athenaTestContext) datacatalogARN(name string) string {
	return fmt.Sprintf("arn:aws:athena:%s:%s:datacatalog/%s", tc.runner.region, tc.runner.AccountID(), name)
}

// createWorkGroup creates a work group and returns its name. Callers
// defer deleteWorkGroup for cleanup.
func (tc *athenaTestContext) createWorkGroup(name string, configuration *types.WorkGroupConfiguration) error {
	_, err := tc.client.CreateWorkGroup(tc.ctx, &athena.CreateWorkGroupInput{
		Name:          aws.String(name),
		Configuration: configuration,
	})
	return err
}

// deleteWorkGroup deletes a work group, ignoring errors (best-effort
// cleanup).
func (tc *athenaTestContext) deleteWorkGroup(name string) {
	_, _ = tc.client.DeleteWorkGroup(tc.ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(name)})
}

// createDataCatalog creates a GLUE data catalog and returns its name.
func (tc *athenaTestContext) createDataCatalog(name, description string) error {
	_, err := tc.client.CreateDataCatalog(tc.ctx, &athena.CreateDataCatalogInput{
		Name:        aws.String(name),
		Type:        types.DataCatalogTypeGlue,
		Description: aws.String(description),
	})
	return err
}

// deleteDataCatalog deletes a data catalog, ignoring errors (best-effort
// cleanup).
func (tc *athenaTestContext) deleteDataCatalog(name string) {
	_, _ = tc.client.DeleteDataCatalog(tc.ctx, &athena.DeleteDataCatalogInput{Name: aws.String(name)})
}

// createNamedQuery creates a named query in the default database and
// returns its id. Callers defer deleteNamedQuery for cleanup.
func (tc *athenaTestContext) createNamedQuery(name, queryString, description string) (string, error) {
	resp, err := tc.client.CreateNamedQuery(tc.ctx, &athena.CreateNamedQueryInput{
		Name:        aws.String(name),
		Database:    aws.String("default"),
		QueryString: aws.String(queryString),
		Description: aws.String(description),
	})
	if err != nil {
		return "", err
	}
	id := aws.ToString(resp.NamedQueryId)
	if id == "" {
		return "", fmt.Errorf("NamedQueryId is empty")
	}
	return id, nil
}

// deleteNamedQuery deletes a named query, ignoring errors (best-effort
// cleanup).
func (tc *athenaTestContext) deleteNamedQuery(id string) {
	_, _ = tc.client.DeleteNamedQuery(tc.ctx, &athena.DeleteNamedQueryInput{NamedQueryId: aws.String(id)})
}

// startQuery starts a query execution against the default database with
// the shared test output location and returns its id.
func (tc *athenaTestContext) startQuery(queryString string) (string, error) {
	resp, err := tc.client.StartQueryExecution(tc.ctx, &athena.StartQueryExecutionInput{
		QueryString: aws.String(queryString),
		QueryExecutionContext: &types.QueryExecutionContext{
			Database: aws.String("default"),
		},
		ResultConfiguration: &types.ResultConfiguration{
			OutputLocation: aws.String("s3://test-bucket/athena/"),
		},
	})
	if err != nil {
		return "", err
	}
	id := aws.ToString(resp.QueryExecutionId)
	if id == "" {
		return "", fmt.Errorf("QueryExecutionId is empty")
	}
	return id, nil
}

// waitForQuery polls GetQueryExecution until the query reaches a terminal
// state (SUCCEEDED, FAILED, or CANCELLED) or the max poll count is reached.
func (tc *athenaTestContext) waitForQuery(queryExecutionId string) (types.QueryExecutionState, error) {
	for i := 0; i < 60; i++ {
		resp, err := tc.client.GetQueryExecution(tc.ctx, &athena.GetQueryExecutionInput{
			QueryExecutionId: aws.String(queryExecutionId),
		})
		if err != nil {
			return "", err
		}
		state := resp.QueryExecution.Status.State
		if state == types.QueryExecutionStateSucceeded ||
			state == types.QueryExecutionStateFailed ||
			state == types.QueryExecutionStateCancelled {
			if state != types.QueryExecutionStateSucceeded {
				reason := aws.ToString(resp.QueryExecution.Status.StateChangeReason)
				return state, fmt.Errorf("query %s ended in %s: %s", queryExecutionId, state, reason)
			}
			return state, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return "", fmt.Errorf("query %s did not complete within 30s", queryExecutionId)
}

// startAndWaitForQuery starts a query and waits until it reaches a
// terminal state.
func (tc *athenaTestContext) startAndWaitForQuery(queryString string) (string, error) {
	id, err := tc.startQuery(queryString)
	if err != nil {
		return "", err
	}
	if _, err := tc.waitForQuery(id); err != nil {
		return "", err
	}
	return id, nil
}

// allWorkGroups walks ListWorkGroups to completion across all pages.
func (tc *athenaTestContext) allWorkGroups() ([]types.WorkGroupSummary, error) {
	return paginate(func(next *string) ([]types.WorkGroupSummary, *string, error) {
		resp, err := tc.client.ListWorkGroups(tc.ctx, &athena.ListWorkGroupsInput{
			MaxResults: aws.Int32(50),
			NextToken:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.WorkGroups, resp.NextToken, nil
	})
}

// allNamedQueryIDs walks ListNamedQueries to completion across all pages.
func (tc *athenaTestContext) allNamedQueryIDs() ([]string, error) {
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListNamedQueries(tc.ctx, &athena.ListNamedQueriesInput{
			MaxResults: aws.Int32(50),
			NextToken:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.NamedQueryIds, resp.NextToken, nil
	})
}

// allQueryExecutions walks ListQueryExecutions to completion across all
// pages.
func (tc *athenaTestContext) allQueryExecutions() ([]string, error) {
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListQueryExecutions(tc.ctx, &athena.ListQueryExecutionsInput{
			MaxResults: aws.Int32(50),
			NextToken:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.QueryExecutionIds, resp.NextToken, nil
	})
}

// createPreparedStatement creates a prepared statement and returns its
// name.
func (tc *athenaTestContext) createPreparedStatement(workGroup, statementName, queryStatement string) error {
	_, err := tc.client.CreatePreparedStatement(tc.ctx, &athena.CreatePreparedStatementInput{
		StatementName:  aws.String(statementName),
		WorkGroup:      aws.String(workGroup),
		QueryStatement: aws.String(queryStatement),
	})
	return err
}

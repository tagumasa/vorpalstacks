package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/athena"
	"vorpalstacks-sdk-tests/config"
)

type athenaTestCtx struct {
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

	tc := &athenaTestCtx{
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
	results = append(results, tc.testTagging()...)
	results = append(results, tc.testEdgeCases()...)

	return results
}

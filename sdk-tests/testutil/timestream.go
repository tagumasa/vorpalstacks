package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"vorpalstacks-sdk-tests/config"
)

type tsTestContext struct {
	writeClient *timestreamwrite.Client
	queryClient *timestreamquery.Client
	iamClient   *iam.Client
	ctx         context.Context
}

func newTSTestContext(r *TestRunner) (*tsTestContext, []TestResult) {
	var results []TestResult
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, []TestResult{{
			Service:  "timestream",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}
	return &tsTestContext{
		writeClient: timestreamwrite.NewFromConfig(cfg),
		queryClient: timestreamquery.NewFromConfig(cfg),
		iamClient:   iam.NewFromConfig(cfg),
		ctx:         context.Background(),
	}, results
}

func (tc *tsTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *tsTestContext) createDatabase(name string) error {
	_, err := tc.writeClient.CreateDatabase(tc.ctx, &timestreamwrite.CreateDatabaseInput{
		DatabaseName: aws.String(name),
	})
	return err
}

func (tc *tsTestContext) createTable(dbName, tableName string) error {
	_, err := tc.writeClient.CreateTable(tc.ctx, &timestreamwrite.CreateTableInput{
		DatabaseName: aws.String(dbName),
		TableName:    aws.String(tableName),
	})
	return err
}

func (tc *tsTestContext) deleteDatabase(name string) {
	_, _ = tc.writeClient.DeleteDatabase(tc.ctx, &timestreamwrite.DeleteDatabaseInput{
		DatabaseName: aws.String(name),
	})
}

func (tc *tsTestContext) deleteTable(dbName, tableName string) {
	_, _ = tc.writeClient.DeleteTable(tc.ctx, &timestreamwrite.DeleteTableInput{
		DatabaseName: aws.String(dbName),
		TableName:    aws.String(tableName),
	})
}

func (r *TestRunner) RunTimestreamTests() []TestResult {
	tc, setupResults := newTSTestContext(r)
	if tc == nil {
		return setupResults
	}

	var results []TestResult
	results = append(results, setupResults...)
	results = append(results, r.runTimestreamDatabaseTests(tc)...)
	results = append(results, r.runTimestreamTableTests(tc)...)
	results = append(results, r.runTimestreamBatchLoadTests(tc)...)
	results = append(results, r.runTimestreamQueryTests(tc)...)
	results = append(results, r.runTimestreamScheduledTests(tc)...)
	results = append(results, r.runTimestreamAccountTests(tc)...)
	results = append(results, r.runTimestreamEdgeTests(tc)...)
	return results
}

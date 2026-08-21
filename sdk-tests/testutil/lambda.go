package testutil

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"vorpalstacks-sdk-tests/config"
)

const lambdaFunctionCode = "exports.handler = async (event) => { return { statusCode: 200, body: 'Hello' }; };"

// zipLambdaCode wraps raw JavaScript source code into a zip archive
// with entry name "index.js", matching the handler "index.handler".
// AWS Lambda requires ZipFile to be a base64-encoded zip archive.
func zipLambdaCode(src string) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, err := w.Create("index.js")
	if err != nil {
		return nil, fmt.Errorf("zip create index.js: %w", err)
	}
	if _, err := f.Write([]byte(src)); err != nil {
		return nil, fmt.Errorf("zip write source: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("zip close: %w", err)
	}
	return buf.Bytes(), nil
}

func (r *TestRunner) RunLambdaTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "lambda",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := lambda.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	cwlClient := cloudwatchlogs.NewFromConfig(cfg)
	ctx := context.Background()

	createIAMRole := func(roleName string) error {
		return IAMCreateRole(iamClient, roleName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`)
	}

	deleteIAMRole := func(roleName string) {
		IAMDeleteRole(iamClient, roleName)
	}

	results = append(results, runLambdaFunctionTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaAliasTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaLayerTests(r, ctx, client)...)
	results = append(results, runLambdaESMTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaESMEngineTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaConfigTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaPermissionTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole, r.region)...)
	results = append(results, runLambdaReferenceTests(r, ctx, client, cwlClient, createIAMRole, deleteIAMRole)...)

	return results
}

func deleteLambdaLogGroup(cwlClient *cloudwatchlogs.Client, ctx context.Context, functionName string) {
	_, err := cwlClient.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String("/aws/lambda/" + functionName),
	})
	if err != nil {
		log.Printf("lambda: failed to delete log group for %s: %v", functionName, err)
	}
}

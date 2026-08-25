package testutil

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"vorpalstacks-sdk-tests/config"
)

const lambdaFunctionCode = "exports.handler = async (event) => { return { statusCode: 200, body: 'Hello' }; };"

// lambdaTestContext carries the shared clients and naming seed for the
// Lambda suite. The per-area run* functions receive it instead of the
// former six-parameter chain.
type lambdaTestContext struct {
	r      *TestRunner
	ctx    context.Context
	client *lambda.Client
	cwl    *cloudwatchlogs.Client
	iam    *iam.Client
	ts     string
}

func (tc *lambdaTestContext) unique(prefix string) string {
	return prefix + "-" + tc.ts
}

// createRole creates an IAM role assumable by lambda.amazonaws.com and
// returns its ARN together with a cleanup function that deletes it.
func (tc *lambdaTestContext) createRole(name string) (string, func(), error) {
	trustPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
	if err := IAMCreateRole(tc.iam, name, trustPolicy); err != nil {
		return "", nil, err
	}
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.r.accountID, name)
	return roleARN, func() { IAMDeleteRole(tc.iam, name) }, nil
}

// createFunction zips the given JavaScript handler source into a zip
// archive and creates the function with the platform-default shape
// (nodejs22x runtime, "index.handler" entry point), returning the function
// ARN together with a cleanup function that deletes the function and its
// log group. Options mutate the request for non-default shapes such as
// environment variables, memory limits, or alternate runtimes.
func (tc *lambdaTestContext) createFunction(name, roleARN, handlerSource string, opts ...func(*lambda.CreateFunctionInput)) (string, func(), error) {
	zipCode, err := zipLambdaCode(handlerSource)
	if err != nil {
		return "", nil, fmt.Errorf("zip lambda code: %w", err)
	}
	input := &lambda.CreateFunctionInput{
		FunctionName: aws.String(name),
		Runtime:      types.RuntimeNodejs22x,
		Role:         aws.String(roleARN),
		Handler:      aws.String("index.handler"),
		Code:         &types.FunctionCode{ZipFile: zipCode},
	}
	for _, opt := range opts {
		opt(input)
	}
	resp, err := tc.client.CreateFunction(tc.ctx, input)
	if err != nil {
		return "", nil, err
	}
	cleanup := func() {
		tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(name)})
		deleteLambdaLogGroup(tc.cwl, tc.ctx, name)
	}
	return aws.ToString(resp.FunctionArn), cleanup, nil
}

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

	tc := &lambdaTestContext{
		r:      r,
		ctx:    context.Background(),
		client: lambda.NewFromConfig(cfg),
		cwl:    cloudwatchlogs.NewFromConfig(cfg),
		iam:    iam.NewFromConfig(cfg),
		ts:     fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	results = append(results, runLambdaFunctionTests(tc)...)
	results = append(results, runLambdaAliasTests(tc)...)
	results = append(results, runLambdaLayerTests(tc)...)
	results = append(results, runLambdaESMTests(tc)...)
	results = append(results, runLambdaESMEngineTests(tc)...)
	results = append(results, runLambdaConfigTests(tc)...)
	results = append(results, runLambdaPermissionTests(tc)...)
	results = append(results, runLambdaReferenceTests(tc)...)

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

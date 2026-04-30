package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"vorpalstacks-sdk-tests/config"
)

const lambdaFunctionCode = "exports.handler = async (event) => { return { statusCode: 200, body: 'Hello' }; };"

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
	ctx := context.Background()

	createIAMRole := func(roleName string) error {
		return IAMCreateRole(iamClient, roleName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}`)
	}

	deleteIAMRole := func(roleName string) {
		IAMDeleteRole(iamClient, roleName)
	}

	results = append(results, runLambdaFunctionTests(r, ctx, client, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaAliasTests(r, ctx, client, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaLayerTests(r, ctx, client)...)
	results = append(results, runLambdaESMTests(r, ctx, client, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaConfigTests(r, ctx, client, createIAMRole, deleteIAMRole)...)
	results = append(results, runLambdaPermissionTests(r, ctx, client, createIAMRole, deleteIAMRole, r.region)...)

	return results
}

package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunAPIGatewayTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "apigateway",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := apigateway.NewFromConfig(cfg)
	ctx := context.Background()

	apiName := fmt.Sprintf("TestAPI-%d", time.Now().UnixNano())
	var apiID string

	results = append(results, r.RunTest("apigateway", "CreateRestApi", func() error {
		resp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:        aws.String(apiName),
			Description: aws.String("Test API"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("API ID is nil")
		}
		apiID = *resp.Id
		return nil
	}))

	results = append(results, r.runAPIGatewayRestApiTests(ctx, client, apiID, apiName)...)
	results = append(results, r.runAPIGatewayResourceTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayMethodLifecycleTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayDeploymentTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayValidatorTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayModelTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayAuthorizerTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayApiKeyTests(ctx, client)...)
	results = append(results, r.runAPIGatewayUsagePlanTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayDomainTests(ctx, client, apiID)...)
	results = append(results, r.runAPIGatewayEdgeTests(ctx, client, apiID)...)

	results = append(results, r.RunTest("apigateway", "DeleteRestApi", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	return results
}

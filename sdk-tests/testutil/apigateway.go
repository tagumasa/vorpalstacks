package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"vorpalstacks-sdk-tests/config"
)

// apigwTestContext carries the suite-wide state shared by every API Gateway
// runner function: the runner, the request context, the SDK client, and the
// REST API created once at section head that most tests build upon.
type apigwTestContext struct {
	r              *TestRunner
	ctx            context.Context
	client         *apigateway.Client
	apiID          string
	apiName        string
	rootResourceID string
}

// uniqueName returns a name unique across suite runs by joining the prefix
// with the current Unix-nanosecond timestamp.
func (tc *apigwTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// createAPI creates a REST API and returns its id and root resource id.
// The caller owns the API lifetime and pairs the call with deleteAPI when
// the API is private to one scenario.
func (tc *apigwTestContext) createAPI(name string) (string, string, error) {
	resp, err := tc.client.CreateRestApi(tc.ctx, &apigateway.CreateRestApiInput{
		Name: aws.String(name),
	})
	if err != nil {
		return "", "", err
	}
	if resp.Id == nil {
		return "", "", fmt.Errorf("API ID is nil")
	}
	if resp.RootResourceId == nil {
		return "", "", fmt.Errorf("root resource ID is nil")
	}
	return *resp.Id, *resp.RootResourceId, nil
}

// deleteAPI deletes a REST API by id, ignoring an empty id so deferred
// cleanup stays safe when creation failed before an id was produced.
func (tc *apigwTestContext) deleteAPI(id string) {
	if id == "" {
		return
	}
	_, _ = tc.client.DeleteRestApi(tc.ctx, &apigateway.DeleteRestApiInput{
		RestApiId: aws.String(id),
	})
}

// createDeployment creates a deployment on the given API and returns its id.
func (tc *apigwTestContext) createDeployment(apiID, description string) (string, error) {
	input := &apigateway.CreateDeploymentInput{RestApiId: aws.String(apiID)}
	if description != "" {
		input.Description = aws.String(description)
	}
	resp, err := tc.client.CreateDeployment(tc.ctx, input)
	if err != nil {
		return "", err
	}
	if resp.Id == nil {
		return "", fmt.Errorf("deployment ID is nil")
	}
	return *resp.Id, nil
}

// findRootResource returns the id of the "/" resource of the given API.
func (tc *apigwTestContext) findRootResource(apiID string) (string, error) {
	resources, err := tc.client.GetResources(tc.ctx, &apigateway.GetResourcesInput{
		RestApiId: aws.String(apiID),
	})
	if err != nil {
		return "", fmt.Errorf("get resources: %v", err)
	}
	for _, res := range resources.Items {
		if res.Path != nil && *res.Path == "/" {
			return *res.Id, nil
		}
	}
	return "", fmt.Errorf("root resource not found")
}

// allRestApis walks every GetRestApis page.
func (tc *apigwTestContext) allRestApis() ([]types.RestApi, error) {
	return paginate(func(next *string) ([]types.RestApi, *string, error) {
		resp, err := tc.client.GetRestApis(tc.ctx, &apigateway.GetRestApisInput{
			Limit:    aws.Int32(500),
			Position: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Items, resp.Position, nil
	})
}

// allApiKeys walks every GetApiKeys page.
func (tc *apigwTestContext) allApiKeys() ([]types.ApiKey, error) {
	return paginate(func(next *string) ([]types.ApiKey, *string, error) {
		resp, err := tc.client.GetApiKeys(tc.ctx, &apigateway.GetApiKeysInput{
			Limit:    aws.Int32(500),
			Position: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Items, resp.Position, nil
	})
}

// allUsagePlans walks every GetUsagePlans page.
func (tc *apigwTestContext) allUsagePlans() ([]types.UsagePlan, error) {
	return paginate(func(next *string) ([]types.UsagePlan, *string, error) {
		resp, err := tc.client.GetUsagePlans(tc.ctx, &apigateway.GetUsagePlansInput{
			Limit:    aws.Int32(500),
			Position: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Items, resp.Position, nil
	})
}

// allModels walks every GetModels page of the given API.
func (tc *apigwTestContext) allModels(apiID string) ([]types.Model, error) {
	return paginate(func(next *string) ([]types.Model, *string, error) {
		resp, err := tc.client.GetModels(tc.ctx, &apigateway.GetModelsInput{
			RestApiId: aws.String(apiID),
			Limit:     aws.Int32(500),
			Position:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Items, resp.Position, nil
	})
}

// allRequestValidators walks every GetRequestValidators page of the given API.
func (tc *apigwTestContext) allRequestValidators(apiID string) ([]types.RequestValidator, error) {
	return paginate(func(next *string) ([]types.RequestValidator, *string, error) {
		resp, err := tc.client.GetRequestValidators(tc.ctx, &apigateway.GetRequestValidatorsInput{
			RestApiId: aws.String(apiID),
			Limit:     aws.Int32(500),
			Position:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Items, resp.Position, nil
	})
}

// allAuthorizers walks every GetAuthorizers page of the given API.
func (tc *apigwTestContext) allAuthorizers(apiID string) ([]types.Authorizer, error) {
	return paginate(func(next *string) ([]types.Authorizer, *string, error) {
		resp, err := tc.client.GetAuthorizers(tc.ctx, &apigateway.GetAuthorizersInput{
			RestApiId: aws.String(apiID),
			Limit:     aws.Int32(500),
			Position:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Items, resp.Position, nil
	})
}

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
	tc := &apigwTestContext{r: r, ctx: context.Background(), client: client}

	// Section head: create the REST API every per-file runner builds upon.
	// On failure this surfaces as one FAIL row named after the step.
	apiName := tc.uniqueName("TestAPI")
	apiID, rootResourceID, err := tc.createAPI(apiName)
	if err != nil {
		return append(results, TestResult{
			Service:  "apigateway",
			TestName: "CreateRestApi",
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}
	tc.apiID = apiID
	tc.apiName = apiName
	tc.rootResourceID = rootResourceID

	results = append(results, r.runAPIGatewayRestApiTests(tc)...)
	results = append(results, r.runAPIGatewayResourceTests(tc)...)
	results = append(results, r.runAPIGatewayMethodLifecycleTests(tc)...)
	results = append(results, r.runAPIGatewayDeploymentTests(tc)...)
	results = append(results, r.runAPIGatewayValidatorTests(tc)...)
	results = append(results, r.runAPIGatewayModelTests(tc)...)
	results = append(results, r.runAPIGatewayAuthorizerTests(tc)...)
	results = append(results, r.runAPIGatewayApiKeyTests(tc)...)
	results = append(results, r.runAPIGatewayUsagePlanTests(tc)...)
	results = append(results, r.runAPIGatewayDomainTests(tc)...)
	results = append(results, r.runAPIGatewayEdgeTests(tc)...)
	results = append(results, r.runAPIGatewayDeepAuditTests(tc)...)
	results = append(results, r.runAPIGatewayValidationTests(tc)...)

	// Teardown: delete the shared REST API. On failure this surfaces as one
	// FAIL row named after the step.
	resp, err := tc.client.DeleteRestApi(tc.ctx, &apigateway.DeleteRestApiInput{
		RestApiId: aws.String(apiID),
	})
	if err != nil || resp == nil {
		msg := "response is nil"
		if err != nil {
			msg = err.Error()
		}
		results = append(results, TestResult{
			Service:  "apigateway",
			TestName: "DeleteRestApi",
			Status:   "FAIL",
			Error:    msg,
		})
	}

	return results
}

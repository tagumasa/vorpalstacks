package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayMethodLifecycleTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "PutMethod_AuthorizationTypes", func() error {
		pmAPI := fmt.Sprintf("PmAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(pmAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("secure"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		for _, authType := range []string{"NONE", "AWS_IAM", "CUSTOM"} {
			_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
				RestApiId:         createResp.Id,
				ResourceId:        resResp.Id,
				HttpMethod:        aws.String("GET"),
				AuthorizationType: aws.String(authType),
			})
			if err != nil {
				return fmt.Errorf("put method with auth %s: %v", authType, err)
			}
			getResp, err := client.GetMethod(ctx, &apigateway.GetMethodInput{
				RestApiId:  createResp.Id,
				ResourceId: resResp.Id,
				HttpMethod: aws.String("GET"),
			})
			if err != nil {
				return fmt.Errorf("get method with auth %s: %v", authType, err)
			}
			if getResp.AuthorizationType == nil || *getResp.AuthorizationType != authType {
				return fmt.Errorf("auth type mismatch for %s, got %v", authType, getResp.AuthorizationType)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethod_InvalidAuthType", func() error {
		iaAPI := fmt.Sprintf("IaAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(iaAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("invalid-auth"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         createResp.Id,
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("INVALID_TYPE"),
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for invalid auth type, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegration_Types", func() error {
		itAPI := fmt.Sprintf("ItAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(itAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("inttest"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         createResp.Id,
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("POST"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		for _, intType := range []types.IntegrationType{
			types.IntegrationTypeMock,
			types.IntegrationTypeHttp,
			types.IntegrationTypeHttpProxy,
			types.IntegrationTypeAwsProxy,
		} {
			_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
				RestApiId:  createResp.Id,
				ResourceId: resResp.Id,
				HttpMethod: aws.String("POST"),
				Type:       intType,
			})
			if err != nil {
				return fmt.Errorf("put integration type %s: %v", intType, err)
			}
			getResp, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
				RestApiId:  createResp.Id,
				ResourceId: resResp.Id,
				HttpMethod: aws.String("POST"),
			})
			if err != nil {
				return fmt.Errorf("get integration type %s: %v", intType, err)
			}
			if getResp.Type != intType {
				return fmt.Errorf("type mismatch, expected %s got %s", intType, getResp.Type)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "TestInvokeMethod", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(apiID),
			ParentId:  aws.String(apiID),
			PathPart:  aws.String("mock"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}
		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(apiID),
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("POST"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}
		_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("POST"),
			Type:       types.IntegrationTypeMock,
			RequestTemplates: map[string]string{
				"application/json": "{\"statusCode\": 200}",
			},
		})
		if err != nil {
			return fmt.Errorf("put integration: %v", err)
		}
		resp, err := client.TestInvokeMethod(ctx, &apigateway.TestInvokeMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("POST"),
			Body:       aws.String(`{"test": "data"}`),
		})
		if err != nil {
			return err
		}
		if resp.Status != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.Status)
		}
		if resp.Log == nil {
			return fmt.Errorf("log is nil")
		}
		return nil
	}))

	results = append(results, r.vtlTests(ctx, client, apiID)...)

	results = append(results, r.RunTest("apigateway", "MethodWithIntegration_FullLifecycle", func() error {
		lcAPI := fmt.Sprintf("LcAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(lcAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("lifecycle"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         createResp.Id,
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
			OperationName:     aws.String("GetLifecycle"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:             createResp.Id,
			ResourceId:            resResp.Id,
			HttpMethod:            aws.String("GET"),
			Type:                  types.IntegrationTypeMock,
			IntegrationHttpMethod: aws.String("POST"),
			Uri:                   aws.String("https://httpbin.org/post"),
			RequestParameters:     map[string]string{"integration.request.header.X-Custom": "'static'"},
			RequestTemplates:      map[string]string{"application/json": "{\"statusCode\":200}"},
			PassthroughBehavior:   aws.String("WHEN_NO_MATCH"),
			TimeoutInMillis:       aws.Int32(3000),
			CacheNamespace:        aws.String("lifecycle"),
			CacheKeyParameters:    []string{"header.Authorization"},
		})
		if err != nil {
			return fmt.Errorf("put integration: %v", err)
		}

		getIntResp, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get integration: %v", err)
		}
		if getIntResp.Uri == nil || *getIntResp.Uri != "https://httpbin.org/post" {
			return fmt.Errorf("uri mismatch, got %v", getIntResp.Uri)
		}
		if getIntResp.TimeoutInMillis != 3000 {
			return fmt.Errorf("timeoutInMillis mismatch, got %d", getIntResp.TimeoutInMillis)
		}

		_, err = client.PutIntegrationResponse(ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:          createResp.Id,
			ResourceId:         resResp.Id,
			HttpMethod:         aws.String("GET"),
			StatusCode:         aws.String("200"),
			ResponseParameters: map[string]string{"method.response.header.Content-Type": "integration.response.header.Content-Type"},
			ResponseTemplates:  map[string]string{"application/json": "$input.json('$')"},
			SelectionPattern:   aws.String("2\\d{2}"),
		})
		if err != nil {
			return fmt.Errorf("put integration response: %v", err)
		}

		_, err = client.PutMethodResponse(ctx, &apigateway.PutMethodResponseInput{
			RestApiId:          createResp.Id,
			ResourceId:         resResp.Id,
			HttpMethod:         aws.String("GET"),
			StatusCode:         aws.String("200"),
			ResponseParameters: map[string]bool{"method.response.header.Content-Type": true},
			ResponseModels:     map[string]string{"application/json": "Empty"},
		})
		if err != nil {
			return fmt.Errorf("put method response: %v", err)
		}

		_, err = client.DeleteMethodResponse(ctx, &apigateway.DeleteMethodResponseInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete method response: %v", err)
		}

		_, err = client.DeleteIntegrationResponse(ctx, &apigateway.DeleteIntegrationResponseInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete integration response: %v", err)
		}

		_, err = client.DeleteIntegration(ctx, &apigateway.DeleteIntegrationInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete integration: %v", err)
		}

		_, err = client.DeleteMethod(ctx, &apigateway.DeleteMethodInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete method: %v", err)
		}

		_, err = client.DeleteResource(ctx, &apigateway.DeleteResourceInput{
			RestApiId:  createResp.Id,
			ResourceId: resResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete resource: %v", err)
		}
		return nil
	}))

	return results
}

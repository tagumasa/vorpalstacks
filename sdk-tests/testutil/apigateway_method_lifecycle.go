package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayMethodLifecycleTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "PutMethod_AuthorizationTypes", func() error {
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("PmAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		resResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("secure"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		for _, authType := range []string{"NONE", "AWS_IAM", "CUSTOM"} {
			_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
				RestApiId:         aws.String(ownAPI),
				ResourceId:        resResp.Id,
				HttpMethod:        aws.String("GET"),
				AuthorizationType: aws.String(authType),
			})
			if err != nil {
				return fmt.Errorf("put method with auth %s: %v", authType, err)
			}
			getResp, err := tc.client.GetMethod(tc.ctx, &apigateway.GetMethodInput{
				RestApiId:  aws.String(ownAPI),
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
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("IaAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		resResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("invalid-auth"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(ownAPI),
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
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("ItAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		resResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("inttest"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(ownAPI),
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
			input := &apigateway.PutIntegrationInput{
				RestApiId:  aws.String(ownAPI),
				ResourceId: resResp.Id,
				HttpMethod: aws.String("POST"),
				Type:       intType,
			}
			// URI is required for all non-MOCK integration types.
			if intType != types.IntegrationTypeMock {
				input.Uri = aws.String("http://example.com/test")
			}
			if intType == types.IntegrationTypeAws || intType == types.IntegrationTypeAwsProxy {
				input.IntegrationHttpMethod = aws.String("POST")
			}
			_, err = tc.client.PutIntegration(tc.ctx, input)
			if err != nil {
				return fmt.Errorf("put integration type %s: %v", intType, err)
			}
			getResp, err := tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
				RestApiId:  aws.String(ownAPI),
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
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(tc.apiID),
			ParentId:  aws.String(tc.rootResourceID),
			PathPart:  aws.String("mock"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}
		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(tc.apiID),
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("POST"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}
		_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
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
		resp, err := tc.client.TestInvokeMethod(tc.ctx, &apigateway.TestInvokeMethodInput{
			RestApiId:  aws.String(tc.apiID),
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

	results = append(results, r.vtlTests(tc)...)

	results = append(results, r.RunTest("apigateway", "MethodWithIntegration_FullLifecycle", func() error {
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("LcAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		resResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("lifecycle"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(ownAPI),
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
			OperationName:     aws.String("GetLifecycle"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:             aws.String(ownAPI),
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

		getIntResp, err := tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(ownAPI),
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

		_, err = tc.client.PutIntegrationResponse(tc.ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:          aws.String(ownAPI),
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

		_, err = tc.client.PutMethodResponse(tc.ctx, &apigateway.PutMethodResponseInput{
			RestApiId:          aws.String(ownAPI),
			ResourceId:         resResp.Id,
			HttpMethod:         aws.String("GET"),
			StatusCode:         aws.String("200"),
			ResponseParameters: map[string]bool{"method.response.header.Content-Type": true},
			ResponseModels:     map[string]string{"application/json": "Empty"},
		})
		if err != nil {
			return fmt.Errorf("put method response: %v", err)
		}

		_, err = tc.client.DeleteMethodResponse(tc.ctx, &apigateway.DeleteMethodResponseInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete method response: %v", err)
		}

		_, err = tc.client.DeleteIntegrationResponse(tc.ctx, &apigateway.DeleteIntegrationResponseInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete integration response: %v", err)
		}

		_, err = tc.client.DeleteIntegration(tc.ctx, &apigateway.DeleteIntegrationInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete integration: %v", err)
		}

		_, err = tc.client.DeleteMethod(tc.ctx, &apigateway.DeleteMethodInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete method: %v", err)
		}

		_, err = tc.client.DeleteResource(tc.ctx, &apigateway.DeleteResourceInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete resource: %v", err)
		}
		return nil
	}))

	return results
}

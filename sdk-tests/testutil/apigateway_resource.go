package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayResourceTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	var resourceID string
	results = append(results, r.RunTest("apigateway", "CreateResource", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(apiID),
			ParentId:  aws.String(apiID),
			PathPart:  aws.String("test"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("resource ID is nil")
		}
		if resp.PathPart == nil || *resp.PathPart != "test" {
			return fmt.Errorf("pathPart mismatch, got %v", resp.PathPart)
		}
		if resp.Path == nil || *resp.Path != "/test" {
			return fmt.Errorf("path mismatch, got %v", resp.Path)
		}
		if resp.ParentId == nil || *resp.ParentId != apiID {
			return fmt.Errorf("parentId mismatch, got %v", resp.ParentId)
		}
		resourceID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetResource", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetResource(ctx, &apigateway.GetResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != resourceID {
			return fmt.Errorf("resource ID mismatch, got %v", resp.Id)
		}
		if resp.Path == nil || *resp.Path != "/test" {
			return fmt.Errorf("path mismatch, got %v", resp.Path)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetResources", func() error {
		if apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := client.GetResources(ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(apiID),
		})
		if err != nil {
			return err
		}
		if resp.Items == nil {
			return fmt.Errorf("items list is nil")
		}
		if len(resp.Items) < 2 {
			return fmt.Errorf("expected at least 2 resources, got %d", len(resp.Items))
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateResource", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateResource(ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/pathPart"),
					Value: aws.String("items"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Path == nil || *resp.Path != "/items" {
			return fmt.Errorf("path not updated, got %v", resp.Path)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(apiID),
			ResourceId:        aws.String(resourceID),
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
			ApiKeyRequired:    false,
		})
		if err != nil {
			return err
		}
		if resp.HttpMethod == nil || *resp.HttpMethod != "GET" {
			return fmt.Errorf("httpMethod mismatch, got %v", resp.HttpMethod)
		}
		if resp.AuthorizationType == nil || *resp.AuthorizationType != "NONE" {
			return fmt.Errorf("authorizationType mismatch, got %v", resp.AuthorizationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetMethod(ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return err
		}
		if resp.HttpMethod == nil || *resp.HttpMethod != "GET" {
			return fmt.Errorf("httpMethod mismatch, got %v", resp.HttpMethod)
		}
		if resp.AuthorizationType == nil || *resp.AuthorizationType != "NONE" {
			return fmt.Errorf("authorizationType mismatch, got %v", resp.AuthorizationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateMethod(ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/authorizationType"),
					Value: aws.String("AWS_IAM"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.AuthorizationType == nil || *resp.AuthorizationType != "AWS_IAM" {
			return fmt.Errorf("authorizationType not updated, got %v", resp.AuthorizationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			Type:       types.IntegrationTypeMock,
			RequestTemplates: map[string]string{
				"application/json": "{\"statusCode\": 200}",
			},
		})
		if err != nil {
			return err
		}
		if resp.Type != types.IntegrationTypeMock {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		if resp.RequestTemplates == nil || resp.RequestTemplates["application/json"] != "{\"statusCode\": 200}" {
			return fmt.Errorf("requestTemplates mismatch, got %v", resp.RequestTemplates)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return err
		}
		if resp.Type != types.IntegrationTypeMock {
			return fmt.Errorf("type mismatch, got %v", resp.Type)
		}
		if resp.RequestTemplates == nil || resp.RequestTemplates["application/json"] != "{\"statusCode\": 200}" {
			return fmt.Errorf("requestTemplates mismatch, got %v", resp.RequestTemplates)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateIntegration(ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/timeoutInMillis"),
					Value: aws.String("5000"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.TimeoutInMillis != 5000 {
			return fmt.Errorf("timeoutInMillis not updated, got %d", resp.TimeoutInMillis)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutIntegrationResponse(ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			ResponseTemplates: map[string]string{
				"application/json": "{\"message\": \"ok\"}",
			},
			SelectionPattern: aws.String("2\\d{2}"),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		if resp.ResponseTemplates == nil || resp.ResponseTemplates["application/json"] != "{\"message\": \"ok\"}" {
			return fmt.Errorf("responseTemplates mismatch, got %v", resp.ResponseTemplates)
		}
		if resp.SelectionPattern == nil || *resp.SelectionPattern != "2\\d{2}" {
			return fmt.Errorf("selectionPattern mismatch, got %v", resp.SelectionPattern)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetIntegrationResponse(ctx, &apigateway.GetIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		if resp.ResponseTemplates == nil || resp.ResponseTemplates["application/json"] != "{\"message\": \"ok\"}" {
			return fmt.Errorf("responseTemplates mismatch, got %v", resp.ResponseTemplates)
		}
		if resp.SelectionPattern == nil || *resp.SelectionPattern != "2\\d{2}" {
			return fmt.Errorf("selectionPattern mismatch, got %v", resp.SelectionPattern)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.UpdateIntegrationResponse(ctx, &apigateway.UpdateIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/selectionPattern"),
					Value: aws.String("ok"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.SelectionPattern == nil || *resp.SelectionPattern != "ok" {
			return fmt.Errorf("selectionPattern not updated, got %v", resp.SelectionPattern)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethodResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.PutMethodResponse(ctx, &apigateway.PutMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			ResponseModels: map[string]string{
				"application/json": "Empty",
			},
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		if resp.ResponseModels == nil || resp.ResponseModels["application/json"] != "Empty" {
			return fmt.Errorf("responseModels mismatch, got %v", resp.ResponseModels)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetMethodResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := client.GetMethodResponse(ctx, &apigateway.GetMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil || *resp.StatusCode != "200" {
			return fmt.Errorf("statusCode mismatch, got %v", resp.StatusCode)
		}
		if resp.ResponseModels == nil || resp.ResponseModels["application/json"] != "Empty" {
			return fmt.Errorf("responseModels mismatch, got %v", resp.ResponseModels)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteMethodResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteMethodResponse(ctx, &apigateway.DeleteMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetMethodResponse(ctx, &apigateway.GetMethodResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err == nil {
			return fmt.Errorf("GetMethodResponse should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteIntegrationResponse", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteIntegrationResponse(ctx, &apigateway.DeleteIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetIntegrationResponse(ctx, &apigateway.GetIntegrationResponseInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err == nil {
			return fmt.Errorf("GetIntegrationResponse should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteIntegration", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteIntegration(ctx, &apigateway.DeleteIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err == nil {
			return fmt.Errorf("GetIntegration should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteMethod", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteMethod(ctx, &apigateway.DeleteMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetMethod(ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err == nil {
			return fmt.Errorf("GetMethod should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteResource", func() error {
		if apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := client.DeleteResource(ctx, &apigateway.DeleteResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetResource(ctx, &apigateway.GetResourceInput{
			RestApiId:  aws.String(apiID),
			ResourceId: aws.String(resourceID),
		})
		if err == nil {
			return fmt.Errorf("GetResource should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateResource_NestedPath", func() error {
		crAPI := fmt.Sprintf("CrAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(crAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		usersResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("users"),
		})
		if err != nil {
			return fmt.Errorf("create users resource: %v", err)
		}

		userIdResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  usersResp.Id,
			PathPart:  aws.String("{userId}"),
		})
		if err != nil {
			return fmt.Errorf("create userId resource: %v", err)
		}
		if userIdResp.Path == nil || *userIdResp.Path != "/users/{userId}" {
			return fmt.Errorf("nested path mismatch, got %v", userIdResp.Path)
		}

		resResp, err := client.GetResources(ctx, &apigateway.GetResourcesInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		if len(resResp.Items) < 3 {
			return fmt.Errorf("expected at least 3 resources (root, users, {userId}), got %d", len(resResp.Items))
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateResource_CascadeChildPaths", func() error {
		ccAPI := fmt.Sprintf("CcAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(ccAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		parentResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("items"),
		})
		if err != nil {
			return fmt.Errorf("create parent: %v", err)
		}

		childResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  parentResp.Id,
			PathPart:  aws.String("{id}"),
		})
		if err != nil {
			return fmt.Errorf("create child: %v", err)
		}
		if childResp.Path == nil || *childResp.Path != "/items/{id}" {
			return fmt.Errorf("child path mismatch before rename, got %v", childResp.Path)
		}

		_, err = client.UpdateResource(ctx, &apigateway.UpdateResourceInput{
			RestApiId:  createResp.Id,
			ResourceId: parentResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/pathPart"), Value: aws.String("products")},
			},
		})
		if err != nil {
			return fmt.Errorf("update parent pathPart: %v", err)
		}

		childAfter, err := client.GetResource(ctx, &apigateway.GetResourceInput{
			RestApiId:  createResp.Id,
			ResourceId: childResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get child after rename: %v", err)
		}
		if childAfter.Path == nil || *childAfter.Path != "/products/{id}" {
			return fmt.Errorf("child path not cascaded, got %v", childAfter.Path)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateResource_DuplicateConflict", func() error {
		dcAPI := fmt.Sprintf("DcAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(dcAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("dup"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}

		_, err = client.CreateResource(ctx, &apigateway.CreateResourceInput{
			RestApiId: createResp.Id,
			ParentId:  createResp.Id,
			PathPart:  aws.String("dup"),
		})
		if err := AssertErrorContains(err, "ConflictException"); err != nil {
			return fmt.Errorf("expected ConflictException for duplicate resource, got: %v", err)
		}
		return nil
	}))

	return results
}

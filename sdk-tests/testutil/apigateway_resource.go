package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayResourceTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var resourceID string
	results = append(results, r.RunTest("apigateway", "CreateResource", func() error {
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(tc.apiID),
			ParentId:  aws.String(tc.rootResourceID),
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
		if resp.ParentId == nil || *resp.ParentId != tc.rootResourceID {
			return fmt.Errorf("parentId mismatch, got %v", resp.ParentId)
		}
		resourceID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetResource", func() error {
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.GetResource(tc.ctx, &apigateway.GetResourceInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" {
			return fmt.Errorf("API ID not available")
		}
		resp, err := tc.client.GetResources(tc.ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.UpdateResource(tc.ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.GetMethod(tc.ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.UpdateMethod(tc.ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.PutIntegrationResponse(tc.ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.GetIntegrationResponse(tc.ctx, &apigateway.GetIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.UpdateIntegrationResponse(tc.ctx, &apigateway.UpdateIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.PutMethodResponse(tc.ctx, &apigateway.PutMethodResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		resp, err := tc.client.GetMethodResponse(tc.ctx, &apigateway.GetMethodResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := tc.client.DeleteMethodResponse(tc.ctx, &apigateway.DeleteMethodResponseInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetMethodResponse(tc.ctx, &apigateway.GetMethodResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := tc.client.DeleteIntegrationResponse(tc.ctx, &apigateway.DeleteIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetIntegrationResponse(tc.ctx, &apigateway.GetIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := tc.client.DeleteIntegration(tc.ctx, &apigateway.DeleteIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := tc.client.DeleteMethod(tc.ctx, &apigateway.DeleteMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetMethod(tc.ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(tc.apiID),
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
		if tc.apiID == "" || resourceID == "" {
			return fmt.Errorf("API ID or resource ID not available")
		}
		_, err := tc.client.DeleteResource(tc.ctx, &apigateway.DeleteResourceInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetResource(tc.ctx, &apigateway.GetResourceInput{
			RestApiId:  aws.String(tc.apiID),
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
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("CrAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		usersResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("users"),
		})
		if err != nil {
			return fmt.Errorf("create users resource: %v", err)
		}

		userIdResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  usersResp.Id,
			PathPart:  aws.String("{userId}"),
		})
		if err != nil {
			return fmt.Errorf("create userId resource: %v", err)
		}
		if userIdResp.Path == nil || *userIdResp.Path != "/users/{userId}" {
			return fmt.Errorf("nested path mismatch, got %v", userIdResp.Path)
		}

		resResp, err := tc.client.GetResources(tc.ctx, &apigateway.GetResourcesInput{
			RestApiId: aws.String(ownAPI),
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
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("CcAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		parentResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("items"),
		})
		if err != nil {
			return fmt.Errorf("create parent: %v", err)
		}

		childResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  parentResp.Id,
			PathPart:  aws.String("{id}"),
		})
		if err != nil {
			return fmt.Errorf("create child: %v", err)
		}
		if childResp.Path == nil || *childResp.Path != "/items/{id}" {
			return fmt.Errorf("child path mismatch before rename, got %v", childResp.Path)
		}

		_, err = tc.client.UpdateResource(tc.ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: parentResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/pathPart"), Value: aws.String("products")},
			},
		})
		if err != nil {
			return fmt.Errorf("update parent pathPart: %v", err)
		}

		childAfter, err := tc.client.GetResource(tc.ctx, &apigateway.GetResourceInput{
			RestApiId:  aws.String(ownAPI),
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
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("DcAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		_, err = tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("dup"),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}

		_, err = tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("dup"),
		})
		if err := AssertErrorContains(err, "ConflictException"); err != nil {
			return fmt.Errorf("expected ConflictException for duplicate resource, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateResource_RootPathPart_Rejected", func() error {
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("RpAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		_, err = tc.client.UpdateResource(tc.ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: aws.String(ownRoot),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/pathPart"), Value: aws.String("root")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for root pathPart update, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateResource_UnsupportedPatchOp_Rejected", func() error {
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("UpAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		resResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("items"),
		})
		if err != nil {
			return fmt.Errorf("create resource: %v", err)
		}

		_, err = tc.client.UpdateResource(tc.ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpMove, From: aws.String("/pathPart"), Path: aws.String("/pathPart")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for unsupported patch op, got: %v", err)
		}
		return nil
	}))

	return results
}

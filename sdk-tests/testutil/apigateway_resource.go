package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayResourceTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var resourceID string
	results = append(results, r.RunTest("apigateway", "CreateResource", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "UpdateMethod_UnknownPatchPath_Rejected", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		_, err := tc.client.UpdateMethod(tc.ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:   types.OpReplace,
					Path: aws.String("/bogusSetting"),
				},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for unknown patch path, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegration", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "UpdateMethod_AuthorizationScopesIndexForms_Rejected", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		// Numeric index addressing appears nowhere in the official patch
		// tables: the documented list operations address the whole
		// /authorizationScopes member only.
		for _, po := range []types.PatchOperation{
			{Op: types.OpAdd, Path: aws.String("/authorizationScopes/0"), Value: aws.String("read:dogs")},
			{Op: types.OpRemove, Path: aws.String("/authorizationScopes/0")},
			{Op: types.OpAdd, Path: aws.String("/authorizationScopes/-"), Value: aws.String("read:dogs")},
		} {
			_, err := tc.client.UpdateMethod(tc.ctx, &apigateway.UpdateMethodInput{
				RestApiId:       aws.String(tc.apiID),
				ResourceId:      aws.String(resourceID),
				HttpMethod:      aws.String("GET"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s, got: %v", po.Op, *po.Path, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetIntegration", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "UpdateMethod_WholeMemberMapPatches", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		_, err := tc.client.UpdateMethod(tc.ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpAdd,
					Path:  aws.String("/requestModels"),
					Value: aws.String(`{"application/json":"Empty"}`),
				},
				{
					Op:    types.OpAdd,
					Path:  aws.String("/requestParameters"),
					Value: aws.String(`{"method.request.header.Authorization":true}`),
				},
				{
					Op:    types.OpAdd,
					Path:  aws.String("/authorizationScopes"),
					Value: aws.String("read:cats"),
				},
			},
		})
		if err != nil {
			return err
		}
		methodResp, err := tc.client.GetMethod(tc.ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get method: %v", err)
		}
		if len(methodResp.RequestModels) != 1 || methodResp.RequestModels["application/json"] != "Empty" {
			return fmt.Errorf("whole-member requestModels add not applied, got %v", methodResp.RequestModels)
		}
		if !methodResp.RequestParameters["method.request.header.Authorization"] {
			return fmt.Errorf("whole-member requestParameters add not applied, got %v", methodResp.RequestParameters)
		}
		if len(methodResp.AuthorizationScopes) != 1 || methodResp.AuthorizationScopes[0] != "read:cats" {
			return fmt.Errorf("whole-member authorizationScopes add not applied, got %v", methodResp.AuthorizationScopes)
		}

		// The patch table allows /requestParameters replace except for
		// MOCK integrations; this method's integration is MOCK.
		_, err = tc.client.UpdateMethod(tc.ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/requestParameters"),
					Value: aws.String(`{"method.request.querystring.page":true}`),
				},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for requestParameters replace on a MOCK integration, got: %v", err)
		}

		_, err = tc.client.UpdateMethod(tc.ctx, &apigateway.UpdateMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/requestModels")},
				{Op: types.OpRemove, Path: aws.String("/authorizationScopes")},
			},
		})
		if err != nil {
			return err
		}
		methodResp, err = tc.client.GetMethod(tc.ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get method: %v", err)
		}
		if len(methodResp.RequestModels) != 0 {
			return fmt.Errorf("whole-member requestModels remove did not clear the map, got %v", methodResp.RequestModels)
		}
		if len(methodResp.AuthorizationScopes) != 0 {
			return fmt.Errorf("whole-member authorizationScopes remove did not clear the list, got %v", methodResp.AuthorizationScopes)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegration", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "UpdateIntegration_WholeMemberMapPatches", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		_, err := tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/requestTemplates"),
					Value: aws.String(`{"application/xml":"<status>ok</status>"}`),
				},
				{Op: types.OpRemove, Path: aws.String("/cacheKeyParameters")},
			},
		})
		if err != nil {
			return err
		}
		integResp, err := tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get integration: %v", err)
		}
		if len(integResp.RequestTemplates) != 1 || integResp.RequestTemplates["application/xml"] != "<status>ok</status>" {
			return fmt.Errorf("whole-member requestTemplates replace not applied, got %v", integResp.RequestTemplates)
		}

		_, err = tc.client.UpdateIntegrationResponse(tc.ctx, &apigateway.UpdateIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/responseTemplates"),
					Value: aws.String(`{"text/plain":"ok"}`),
				},
			},
		})
		if err != nil {
			return err
		}
		intRespResp, err := tc.client.GetIntegrationResponse(tc.ctx, &apigateway.GetIntegrationResponseInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("get integration response: %v", err)
		}
		if len(intRespResp.ResponseTemplates) != 1 || intRespResp.ResponseTemplates["text/plain"] != "ok" {
			return fmt.Errorf("whole-member responseTemplates replace not applied, got %v", intRespResp.ResponseTemplates)
		}

		_, err = tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/requestTemplates")},
			},
		})
		if err != nil {
			return err
		}
		integResp, err = tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get integration: %v", err)
		}
		if len(integResp.RequestTemplates) != 0 {
			return fmt.Errorf("whole-member requestTemplates remove did not clear the map, got %v", integResp.RequestTemplates)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegration_CacheKeyParametersWholeReplace", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		// The /cacheKeyParameters row documents add, replace and remove:
		// replace sets the whole list from a JSON array of strings.
		resp, err := tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/cacheKeyParameters"),
					Value: aws.String(`["method.request.header.Authorization","method.request.querystring.q"]`),
				},
			},
		})
		if err != nil {
			return err
		}
		want := []string{"method.request.header.Authorization", "method.request.querystring.q"}
		if len(resp.CacheKeyParameters) != 2 || resp.CacheKeyParameters[0] != want[0] || resp.CacheKeyParameters[1] != want[1] {
			return fmt.Errorf("whole replace not applied, got %v", resp.CacheKeyParameters)
		}

		// add appends a single entry.
		resp, err = tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/cacheKeyParameters"), Value: aws.String("method.request.path.op")},
			},
		})
		if err != nil {
			return err
		}
		if len(resp.CacheKeyParameters) != 3 || resp.CacheKeyParameters[2] != "method.request.path.op" {
			return fmt.Errorf("whole add did not append, got %v", resp.CacheKeyParameters)
		}

		// Numeric index addressing appears nowhere in the official patch
		// tables.
		_, err = tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/cacheKeyParameters/0")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for the indexed form, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegration_MockIntegrationPatches_Rejected", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		// The uri and httpMethod rows document replace "except for MOCK
		// integrations", and the type row marks every operation Not
		// supported — this integration is MOCK.
		for _, po := range []types.PatchOperation{
			{Op: types.OpReplace, Path: aws.String("/uri"), Value: aws.String("http://example.test/")},
			{Op: types.OpReplace, Path: aws.String("/httpMethod"), Value: aws.String("POST")},
			{Op: types.OpReplace, Path: aws.String("/type"), Value: aws.String("AWS")},
		} {
			_, err := tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
				RestApiId:       aws.String(tc.apiID),
				ResourceId:      aws.String(resourceID),
				HttpMethod:      aws.String("GET"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s of a MOCK integration, got: %v", po.Op, *po.Path, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateIntegration_TargetAndTransferMode", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
		}
		// The integrationTarget and responseTransferMode rows document
		// replace only.
		target := "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-alb/abc/def"
		resp, err := tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/integrationTarget"), Value: aws.String(target)},
				{Op: types.OpReplace, Path: aws.String("/responseTransferMode"), Value: aws.String("STREAM")},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.IntegrationTarget) != target {
			return fmt.Errorf("integrationTarget not applied, got %v", resp.IntegrationTarget)
		}
		if string(resp.ResponseTransferMode) != "STREAM" {
			return fmt.Errorf("responseTransferMode not applied, got %v", resp.ResponseTransferMode)
		}
		getResp, err := tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(tc.apiID),
			ResourceId: aws.String(resourceID),
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return err
		}
		if string(getResp.ResponseTransferMode) != "STREAM" {
			return fmt.Errorf("responseTransferMode not persisted, got %v", getResp.ResponseTransferMode)
		}

		// An invalid enum value and a non-replace operation reject.
		for _, po := range []types.PatchOperation{
			{Op: types.OpReplace, Path: aws.String("/responseTransferMode"), Value: aws.String("STREAMED")},
			{Op: types.OpAdd, Path: aws.String("/integrationTarget"), Value: aws.String("x")},
		} {
			_, err := tc.client.UpdateIntegration(tc.ctx, &apigateway.UpdateIntegrationInput{
				RestApiId:       aws.String(tc.apiID),
				ResourceId:      aws.String(resourceID),
				HttpMethod:      aws.String("GET"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s, got: %v", po.Op, *po.Path, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutMethodResponse", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetMethodResponse should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteIntegrationResponse", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetIntegrationResponse should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteIntegration", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetIntegration should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteMethod", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetMethod should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteResource", func() error {
		if err := tc.require(tc.apiID, resourceID); err != nil {
			return err
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
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetResource should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateResource_NestedPath", func() error {
		ownAPI, ownRoot, err := tc.createOwnAPI("CrAPI")
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
		ownAPI, ownRoot, err := tc.createOwnAPI("CcAPI")
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
		ownAPI, ownRoot, err := tc.createOwnAPI("DcAPI")
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
		ownAPI, ownRoot, err := tc.createOwnAPI("RpAPI")
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
		ownAPI, ownRoot, err := tc.createOwnAPI("UpAPI")
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

	results = append(results, r.RunTest("apigateway", "UpdateResource_ParentId", func() error {
		ownAPI, ownRoot, err := tc.createOwnAPI("PrAPI")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		aResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("group"),
		})
		if err != nil {
			return fmt.Errorf("create resource a: %v", err)
		}
		bResp, err := tc.client.CreateResource(tc.ctx, &apigateway.CreateResourceInput{
			RestApiId: aws.String(ownAPI),
			ParentId:  aws.String(ownRoot),
			PathPart:  aws.String("item"),
		})
		if err != nil {
			return fmt.Errorf("create resource b: %v", err)
		}

		// The /parentId row documents replace: moving b under a updates
		// its full path.
		resp, err := tc.client.UpdateResource(tc.ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: bResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/parentId"), Value: aResp.Id},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.ParentId) != aws.ToString(aResp.Id) || aws.ToString(resp.Path) != "/group/item" {
			return fmt.Errorf("parentId replace not applied, got parent=%v path=%v", resp.ParentId, resp.Path)
		}

		// Moving a resource under its own subtree would create a cycle.
		_, err = tc.client.UpdateResource(tc.ctx, &apigateway.UpdateResourceInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: aResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/parentId"), Value: bResp.Id},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for a cyclic parentId, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateMethodResponse_NilMapAddParameter", func() error {
		ownAPI, ownRoot, err := tc.createOwnAPI("MrAPI")
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
		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(ownAPI),
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}
		// The response is created without ResponseParameters/ResponseModels,
		// so the stored nested maps are nil until the update applier
		// initialises them — the first add-patch must not panic.
		_, err = tc.client.PutMethodResponse(tc.ctx, &apigateway.PutMethodResponseInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("put method response: %v", err)
		}
		updated, err := tc.client.UpdateMethodResponse(tc.ctx, &apigateway.UpdateMethodResponseInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/responseParameters/method.response.header.X-Custom"), Value: aws.String("true")},
				{Op: types.OpAdd, Path: aws.String("/responseModels/text"), Value: aws.String("Empty")},
			},
		})
		if err != nil {
			return fmt.Errorf("update method response over nil maps: %v", err)
		}
		if v, ok := updated.ResponseParameters["method.response.header.X-Custom"]; !ok || !v {
			return fmt.Errorf("responseParameters not applied, got %v", updated.ResponseParameters)
		}
		if m := updated.ResponseModels["text"]; m != "Empty" {
			return fmt.Errorf("responseModels not applied, got %v", updated.ResponseModels)
		}
		getResp, err := tc.client.GetMethodResponse(tc.ctx, &apigateway.GetMethodResponseInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: resResp.Id,
			HttpMethod: aws.String("GET"),
			StatusCode: aws.String("200"),
		})
		if err != nil {
			return fmt.Errorf("get method response: %v", err)
		}
		if v, ok := getResp.ResponseParameters["method.response.header.X-Custom"]; !ok || !v {
			return fmt.Errorf("responseParameters not persisted, got %v", getResp.ResponseParameters)
		}
		return nil
	}))

	return results
}

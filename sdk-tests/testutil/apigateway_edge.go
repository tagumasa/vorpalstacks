package testutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayEdgeTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "TagResource_UntagResource", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{name: "restapi-plane", probe: func() error {
				ownAPI, _, err := tc.createOwnAPI("TagAPI")
				if err != nil {
					return fmt.Errorf("create: %v", err)
				}
				defer tc.deleteAPI(ownAPI)

				arn := tc.resourceARN("restapis", ownAPI)

				_, err = tc.client.TagResource(tc.ctx, &apigateway.TagResourceInput{
					ResourceArn: aws.String(arn),
					Tags: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				})
				if err != nil {
					return fmt.Errorf("tag: %v", err)
				}

				tagResp, err := tc.client.GetTags(tc.ctx, &apigateway.GetTagsInput{
					ResourceArn: aws.String(arn),
				})
				if err != nil {
					return fmt.Errorf("get tags: %v", err)
				}
				if tagResp.Tags == nil || tagResp.Tags["key1"] != "value1" {
					return fmt.Errorf("tags mismatch, got %v", tagResp.Tags)
				}

				_, err = tc.client.UntagResource(tc.ctx, &apigateway.UntagResourceInput{
					ResourceArn: aws.String(arn),
					TagKeys:     []string{"key2"},
				})
				if err != nil {
					return fmt.Errorf("untag: %v", err)
				}

				tagResp2, err := tc.client.GetTags(tc.ctx, &apigateway.GetTagsInput{
					ResourceArn: aws.String(arn),
				})
				if err != nil {
					return fmt.Errorf("get tags after untag: %v", err)
				}
				if _, exists := tagResp2.Tags["key2"]; exists {
					return fmt.Errorf("key2 should have been removed")
				}
				if tagResp2.Tags["key1"] != "value1" {
					return fmt.Errorf("key1 should still exist")
				}
				return nil
			}},
			{name: "usageplan-plane", probe: func() error {
				planID, err := tc.createOwnUsagePlan("TagPlan")
				if err != nil {
					return fmt.Errorf("create usage plan: %v", err)
				}
				defer tc.deleteUsagePlan(planID)

				arn := tc.resourceARN("usageplans", planID)

				_, err = tc.client.TagResource(tc.ctx, &apigateway.TagResourceInput{
					ResourceArn: aws.String(arn),
					Tags: map[string]string{
						"env":  "test",
						"team": "qa",
					},
				})
				if err != nil {
					return fmt.Errorf("tag usage plan: %v", err)
				}

				tagResp, err := tc.client.GetTags(tc.ctx, &apigateway.GetTagsInput{
					ResourceArn: aws.String(arn),
				})
				if err != nil {
					return fmt.Errorf("get tags: %v", err)
				}
				if tagResp.Tags["env"] != "test" || tagResp.Tags["team"] != "qa" {
					return fmt.Errorf("tags mismatch: %v", tagResp.Tags)
				}

				_, err = tc.client.UntagResource(tc.ctx, &apigateway.UntagResourceInput{
					ResourceArn: aws.String(arn),
					TagKeys:     []string{"team"},
				})
				if err != nil {
					return fmt.Errorf("untag: %v", err)
				}

				tagResp2, err := tc.client.GetTags(tc.ctx, &apigateway.GetTagsInput{
					ResourceArn: aws.String(arn),
				})
				if err != nil {
					return fmt.Errorf("get tags after untag: %v", err)
				}
				if _, exists := tagResp2.Tags["team"]; exists {
					return fmt.Errorf("team should have been removed")
				}
				if tagResp2.Tags["env"] != "test" {
					return fmt.Errorf("env should still exist")
				}
				return nil
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	// Operations against resources that do not exist fail with the
	// modelled NotFoundException, as the service model specifies —
	// including tag operations against a stage that does not exist.
	results = append(results, r.RunTest("apigateway", "NonExistentResources", func() error {
		rows := []struct {
			name  string
			probe func() error
		}{
			{name: "GetRestApi", probe: func() error {
				_, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
					RestApiId: aws.String("nonexistent_xyz"),
				})
				return AssertErrorContains(err, "NotFoundException")
			}},
			{name: "DeleteRestApi", probe: func() error {
				_, err := tc.client.DeleteRestApi(tc.ctx, &apigateway.DeleteRestApiInput{
					RestApiId: aws.String("nonexistent_xyz"),
				})
				return AssertErrorContains(err, "NotFoundException")
			}},
			{name: "GetStage", probe: func() error {
				ownAPI, _, err := tc.createOwnAPI("TmpAPI")
				if err != nil {
					return fmt.Errorf("create: %v", err)
				}
				defer tc.deleteAPI(ownAPI)
				_, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
					RestApiId: aws.String(ownAPI),
					StageName: aws.String("nonexistent_stage"),
				})
				return AssertErrorContains(err, "NotFoundException")
			}},
			{name: "TagResource_NonExistentStage", probe: func() error {
				arn := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/no-such-api-%d/stages/prod",
					tc.r.region, time.Now().UnixNano())
				_, err := tc.client.TagResource(tc.ctx, &apigateway.TagResourceInput{
					ResourceArn: aws.String(arn),
					Tags:        map[string]string{"Environment": "test"},
				})
				if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
					return fmt.Errorf("TagResource: %v", aerr)
				}
				_, err = tc.client.UntagResource(tc.ctx, &apigateway.UntagResourceInput{
					ResourceArn: aws.String(arn),
					TagKeys:     []string{"Environment"},
				})
				if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
					return fmt.Errorf("UntagResource: %v", aerr)
				}
				_, err = tc.client.GetTags(tc.ctx, &apigateway.GetTagsInput{
					ResourceArn: aws.String(arn),
				})
				if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
					return fmt.Errorf("GetTags: %v", aerr)
				}
				return nil
			}},
		}
		for _, row := range rows {
			if err := row.probe(); err != nil {
				return fmt.Errorf("%s: %v", row.name, err)
			}
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runAPIGatewayDeepAuditTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "PutMethod_AuthorizationScopes_RoundTrip", func() error {
		ownAPI, rootId, err := tc.createOwnAPI("AuthScopeAPI")
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer tc.deleteAPI(ownAPI)

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:           aws.String(ownAPI),
			ResourceId:          &rootId,
			HttpMethod:          aws.String("GET"),
			AuthorizationType:   aws.String("NONE"),
			AuthorizationScopes: []string{"read", "write"},
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		method, err := tc.client.GetMethod(tc.ctx, &apigateway.GetMethodInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: &rootId,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get method: %v", err)
		}

		if len(method.AuthorizationScopes) != 2 {
			return fmt.Errorf("expected 2 authorizationScopes, got %d", len(method.AuthorizationScopes))
		}
		if method.AuthorizationScopes[0] != "read" || method.AuthorizationScopes[1] != "write" {
			return fmt.Errorf("authorizationScopes mismatch: %v", method.AuthorizationScopes)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "PutIntegration_TlsConfig_Timeout_RoundTrip", func() error {
		ownAPI, rootId, err := tc.createOwnAPI("TlsAPI")
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer tc.deleteAPI(ownAPI)

		_, err = tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId:         aws.String(ownAPI),
			ResourceId:        &rootId,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		_, err = tc.client.PutIntegration(tc.ctx, &apigateway.PutIntegrationInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: &rootId,
			HttpMethod: aws.String("GET"),
			Type:       "HTTP",
			Uri:        aws.String("https://example.com"),
			TlsConfig: &types.TlsConfig{
				InsecureSkipVerification: true,
			},
		})
		if err != nil {
			return fmt.Errorf("put integration: %v", err)
		}

		integ, err := tc.client.GetIntegration(tc.ctx, &apigateway.GetIntegrationInput{
			RestApiId:  aws.String(ownAPI),
			ResourceId: &rootId,
			HttpMethod: aws.String("GET"),
		})
		if err != nil {
			return fmt.Errorf("get integration: %v", err)
		}

		if integ.TlsConfig == nil {
			return fmt.Errorf("tlsConfig is nil")
		}
		if !integ.TlsConfig.InsecureSkipVerification {
			return fmt.Errorf("insecureSkipVerification should be true")
		}
		if integ.TimeoutInMillis != 29000 {
			return fmt.Errorf("expected default timeoutInMillis 29000, got %d", integ.TimeoutInMillis)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateUsagePlan_QuotaOffset", func() error {
		createResp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(tc.uniqueName("OffsetPlan")),
			Quota: &types.QuotaSettings{
				Limit:  1000,
				Offset: 10,
				Period: types.QuotaPeriodTypeMonth,
			},
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(aws.ToString(createResp.Id))

		plan, err := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}

		if plan.Quota == nil {
			return fmt.Errorf("quota is nil")
		}
		if plan.Quota.Offset != 10 {
			return fmt.Errorf("expected offset 10, got %d", plan.Quota.Offset)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateRestApi_DisableExecuteApiEndpoint", func() error {
		createResp, err := tc.client.CreateRestApi(tc.ctx, &apigateway.CreateRestApiInput{
			Name:                      aws.String(tc.uniqueName("DisableExec")),
			DisableExecuteApiEndpoint: true,
			MinimumCompressionSize:    aws.Int32(0),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(*createResp.Id)

		api, err := tc.client.GetRestApi(tc.ctx, &apigateway.GetRestApiInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}

		if !api.DisableExecuteApiEndpoint {
			return fmt.Errorf("disableExecuteApiEndpoint should be true")
		}
		if api.MinimumCompressionSize == nil {
			return fmt.Errorf("minimumCompressionSize should be present (0)")
		}
		if *api.MinimumCompressionSize != 0 {
			return fmt.Errorf("expected minimumCompressionSize 0, got %d", *api.MinimumCompressionSize)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateRestApi_EndpointConfiguration", func() error {
		// The create-time endpointConfiguration carries ipAddressType and
		// vpcEndpointIds beside types, and both round-trip.
		createResp, err := tc.client.CreateRestApi(tc.ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tc.uniqueName("EndpConf")),
			EndpointConfiguration: &types.EndpointConfiguration{
				Types:          []types.EndpointType{types.EndpointTypePrivate},
				IpAddressType:  types.IpAddressTypeDualstack,
				VpcEndpointIds: []string{"vpce-0abc123"},
			},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(*createResp.Id)
		if createResp.EndpointConfiguration == nil ||
			len(createResp.EndpointConfiguration.VpcEndpointIds) != 1 ||
			createResp.EndpointConfiguration.VpcEndpointIds[0] != "vpce-0abc123" ||
			createResp.EndpointConfiguration.IpAddressType != types.IpAddressTypeDualstack {
			return fmt.Errorf("endpointConfiguration not echoed on create, got %+v", createResp.EndpointConfiguration)
		}

		_, err = tc.client.CreateRestApi(tc.ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tc.uniqueName("EndpConfBad")),
			EndpointConfiguration: &types.EndpointConfiguration{
				Types:         []types.EndpointType{types.EndpointTypeRegional},
				IpAddressType: "bogus",
			},
		})
		if aerr := AssertErrorContains(err, "BadRequestException"); aerr != nil {
			return fmt.Errorf("expected BadRequestException for ipAddressType bogus, got: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GatewayResponse_CRUD", func() error {
		ownAPI, _, err := tc.createAPI(tc.uniqueName("GwResp"))
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer tc.deleteAPI(ownAPI)

		put, err := tc.client.PutGatewayResponse(tc.ctx, &apigateway.PutGatewayResponseInput{
			RestApiId:    aws.String(ownAPI),
			ResponseType: types.GatewayResponseTypeThrottled,
			StatusCode:   aws.String("429"),
			ResponseParameters: map[string]string{
				"gatewayresponse.header.X-RateLimit": "lim",
			},
			ResponseTemplates: map[string]string{"application/json": `{"message":"slow down"}`},
		})
		if err != nil {
			return err
		}
		if put.DefaultResponse {
			return fmt.Errorf("a customised gateway response must not be the default one")
		}
		if aws.ToString(put.StatusCode) != "429" || put.ResponseParameters["gatewayresponse.header.X-RateLimit"] != "lim" {
			return fmt.Errorf("put echo mismatch, got %+v", put)
		}

		get, err := tc.client.GetGatewayResponse(tc.ctx, &apigateway.GetGatewayResponseInput{
			RestApiId:    aws.String(ownAPI),
			ResponseType: types.GatewayResponseTypeThrottled,
		})
		if err != nil {
			return err
		}
		if get.ResponseType != types.GatewayResponseTypeThrottled || get.ResponseTemplates["application/json"] != `{"message":"slow down"}` {
			return fmt.Errorf("get mismatch, got %+v", get)
		}

		list, err := tc.client.GetGatewayResponses(tc.ctx, &apigateway.GetGatewayResponsesInput{
			RestApiId: aws.String(ownAPI),
		})
		if err != nil {
			return err
		}
		if len(list.Items) != 1 || list.Items[0].ResponseType != types.GatewayResponseTypeThrottled {
			return fmt.Errorf("expected exactly the THROTTLED response in listing, got %+v", list.Items)
		}

		// The documented patch surface: replace of /statusCode and
		// add/replace/remove on the parameter and template maps.
		upd, err := tc.client.UpdateGatewayResponse(tc.ctx, &apigateway.UpdateGatewayResponseInput{
			RestApiId:    aws.String(ownAPI),
			ResponseType: types.GatewayResponseTypeThrottled,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/statusCode"), Value: aws.String("503")},
				{Op: types.OpAdd, Path: aws.String("/responseParameters/gatewayresponse.header.X-Reason"), Value: aws.String("overload")},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(upd.StatusCode) != "503" || upd.ResponseParameters["gatewayresponse.header.X-Reason"] != "overload" {
			return fmt.Errorf("patch not applied, got %+v", upd)
		}

		if _, err := tc.client.DeleteGatewayResponse(tc.ctx, &apigateway.DeleteGatewayResponseInput{
			RestApiId:    aws.String(ownAPI),
			ResponseType: types.GatewayResponseTypeThrottled,
		}); err != nil {
			return err
		}
		_, err = tc.client.GetGatewayResponse(tc.ctx, &apigateway.GetGatewayResponseInput{
			RestApiId:    aws.String(ownAPI),
			ResponseType: types.GatewayResponseTypeThrottled,
		})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", aerr)
		}

		_, err = tc.client.PutGatewayResponse(tc.ctx, &apigateway.PutGatewayResponseInput{
			RestApiId:    aws.String(ownAPI),
			ResponseType: types.GatewayResponseType("NOT_A_MODELLED_TYPE"),
		})
		if aerr := AssertErrorContains(err, "BadRequestException"); aerr != nil {
			return fmt.Errorf("expected BadRequestException for unmodelled responseType, got: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "FlushStageCache_Operations", func() error {
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("FlushStage"))
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		if _, err := tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId: aws.String(ownAPI), ResourceId: aws.String(ownRoot),
			HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
		}); err != nil {
			return fmt.Errorf("put method: %v", err)
		}
		if _, err := tc.client.CreateDeployment(tc.ctx, &apigateway.CreateDeploymentInput{
			RestApiId: aws.String(ownAPI), StageName: aws.String("flushstage"),
		}); err != nil {
			return fmt.Errorf("create deployment: %v", err)
		}

		// Both flush operations reply 202 with an empty body; no
		// cache-cluster precondition is documented.
		if _, err := tc.client.FlushStageCache(tc.ctx, &apigateway.FlushStageCacheInput{
			RestApiId: aws.String(ownAPI), StageName: aws.String("flushstage"),
		}); err != nil {
			return fmt.Errorf("flush stage cache: %v", err)
		}
		if _, err := tc.client.FlushStageAuthorizersCache(tc.ctx, &apigateway.FlushStageAuthorizersCacheInput{
			RestApiId: aws.String(ownAPI), StageName: aws.String("flushstage"),
		}); err != nil {
			return fmt.Errorf("flush stage authorizers cache: %v", err)
		}

		_, err = tc.client.FlushStageCache(tc.ctx, &apigateway.FlushStageCacheInput{
			RestApiId: aws.String(ownAPI), StageName: aws.String("missing-stage"),
		})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("expected NotFoundException for missing stage, got: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "SdkTypes_StaticTable", func() error {
		typesResp, err := tc.client.GetSdkTypes(tc.ctx, &apigateway.GetSdkTypesInput{})
		if err != nil {
			return err
		}
		var android *types.SdkType
		for i := range typesResp.Items {
			if typesResp.Items[i].Id != nil && *typesResp.Items[i].Id == "android" {
				android = &typesResp.Items[i]
			}
		}
		if android == nil {
			return fmt.Errorf("android SDK type missing, got %+v", typesResp.Items)
		}
		required := map[string]bool{}
		for _, prop := range android.ConfigurationProperties {
			if prop.Required {
				required[aws.ToString(prop.Name)] = true
			}
		}
		for _, name := range []string{"groupId", "artifactId", "artifactVersion", "invokerPackage"} {
			if !required[name] {
				return fmt.Errorf("android SDK type missing required property %s, got %+v", name, android.ConfigurationProperties)
			}
		}

		java, err := tc.client.GetSdkType(tc.ctx, &apigateway.GetSdkTypeInput{Id: aws.String("java")})
		if err != nil {
			return err
		}
		if len(java.ConfigurationProperties) != 2 || !java.ConfigurationProperties[0].Required {
			return fmt.Errorf("java SDK type properties mismatch, got %+v", java.ConfigurationProperties)
		}

		_, err = tc.client.GetSdkType(tc.ctx, &apigateway.GetSdkTypeInput{Id: aws.String("no-such-sdk")})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("expected NotFoundException for unknown SDK type, got: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetExport_Document", func() error {
		ownAPI, ownRoot, err := tc.createAPI(tc.uniqueName("Export"))
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		if _, err := tc.client.PutMethod(tc.ctx, &apigateway.PutMethodInput{
			RestApiId: aws.String(ownAPI), ResourceId: aws.String(ownRoot),
			HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
			OperationName: aws.String("ListExport"),
		}); err != nil {
			return fmt.Errorf("put method: %v", err)
		}
		if _, err := tc.client.CreateModel(tc.ctx, &apigateway.CreateModelInput{
			RestApiId: aws.String(ownAPI), Name: aws.String("ExportModel"),
			ContentType: aws.String("application/json"),
			Schema:      aws.String(`{"type":"object","properties":{"id":{"type":"string"}}}`),
		}); err != nil {
			return fmt.Errorf("create model: %v", err)
		}
		if _, err := tc.client.CreateDeployment(tc.ctx, &apigateway.CreateDeploymentInput{
			RestApiId: aws.String(ownAPI), StageName: aws.String("exportstage"),
		}); err != nil {
			return fmt.Errorf("create deployment: %v", err)
		}

		exported, err := tc.client.GetExport(tc.ctx, &apigateway.GetExportInput{
			RestApiId:  aws.String(ownAPI),
			StageName:  aws.String("exportstage"),
			ExportType: aws.String("swagger"),
			Accepts:    aws.String("application/json"),
		})
		if err != nil {
			return err
		}
		if aws.ToString(exported.ContentType) != "application/json" {
			return fmt.Errorf("contentType mismatch, got %v", exported.ContentType)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(exported.Body, &doc); err != nil {
			return fmt.Errorf("export body is not JSON: %v", err)
		}
		if doc["swagger"] != "2.0" {
			return fmt.Errorf("expected swagger 2.0 document, got %v", doc["swagger"])
		}
		paths, ok := doc["paths"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("paths missing from export, got %v", doc["paths"])
		}
		rootOp, ok := paths["/"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("root path missing from export, got %v", paths)
		}
		get, ok := rootOp["get"].(map[string]interface{})
		if !ok || get["operationId"] != "ListExport" {
			return fmt.Errorf("operation not exported, got %v", rootOp["get"])
		}
		defs, ok := doc["definitions"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("definitions missing from export, got %v", doc["definitions"])
		}
		if _, ok := defs["ExportModel"]; !ok {
			return fmt.Errorf("model definition missing from export, got %v", defs)
		}

		// The Go SDK client pins Accept to application/json on every API
		// Gateway call (its build-stack AcceptHeader middleware overwrites
		// the operation serializer's Accept from Accepts), so the YAML
		// branch is unreachable through this client; a server-side unit
		// test pins it instead. Here the oas30 family is pinned via JSON.
		oasExport, err := tc.client.GetExport(tc.ctx, &apigateway.GetExportInput{
			RestApiId:  aws.String(ownAPI),
			StageName:  aws.String("exportstage"),
			ExportType: aws.String("oas30"),
		})
		if err != nil {
			return err
		}
		if aws.ToString(oasExport.ContentType) != "application/json" || !strings.Contains(string(oasExport.Body), `"openapi": "3.0.1"`) {
			return fmt.Errorf("oas30 export mismatch, type %v body %s", oasExport.ContentType, oasExport.Body)
		}

		_, err = tc.client.GetExport(tc.ctx, &apigateway.GetExportInput{
			RestApiId: aws.String(ownAPI), StageName: aws.String("exportstage"),
			ExportType: aws.String("bogus"),
		})
		if aerr := AssertErrorContains(err, "BadRequestException"); aerr != nil {
			return fmt.Errorf("expected BadRequestException for bogus exportType, got: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetModelTemplate_Skeleton", func() error {
		ownAPI, _, err := tc.createAPI(tc.uniqueName("Tpl"))
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		if _, err := tc.client.CreateModel(tc.ctx, &apigateway.CreateModelInput{
			RestApiId: aws.String(ownAPI), Name: aws.String("TplModel"),
			ContentType: aws.String("application/json"),
			Schema:      aws.String(`{"type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"}}}`),
		}); err != nil {
			return fmt.Errorf("create model: %v", err)
		}

		template, err := tc.client.GetModelTemplate(tc.ctx, &apigateway.GetModelTemplateInput{
			RestApiId: aws.String(ownAPI), ModelName: aws.String("TplModel"),
		})
		if err != nil {
			return err
		}
		if template.Value == nil ||
			!strings.Contains(*template.Value, "$inputRoot.id") ||
			!strings.Contains(*template.Value, "$inputRoot.name") {
			return fmt.Errorf("template skeleton missing property echoes, got %v", template.Value)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "Account_RoundTrip", func() error {
		account, err := tc.client.GetAccount(tc.ctx, &apigateway.GetAccountInput{})
		if err != nil {
			return err
		}
		usagePlans := false
		for _, f := range account.Features {
			if f == "UsagePlans" {
				usagePlans = true
			}
		}
		if !usagePlans {
			return fmt.Errorf("features = %v, want the UsagePlans entry", account.Features)
		}
		if account.ThrottleSettings == nil {
			return fmt.Errorf("throttleSettings missing from the default account")
		}
		if account.ThrottleSettings.RateLimit != 10000 ||
			account.ThrottleSettings.BurstLimit != 5000 {
			return fmt.Errorf("default throttleSettings = %v/%v, want 10000/5000",
				account.ThrottleSettings.RateLimit, account.ThrottleSettings.BurstLimit)
		}

		updated, err := tc.client.UpdateAccount(tc.ctx, &apigateway.UpdateAccountInput{
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/cloudwatchRoleArn"), Value: aws.String("arn:aws:iam::123456789012:role/apigateway-metrics")},
				{Op: types.OpAdd, Path: aws.String("/features"), Value: aws.String("ProbeFeature")},
			},
		})
		if err != nil {
			return err
		}
		defer tc.client.UpdateAccount(tc.ctx, &apigateway.UpdateAccountInput{
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/cloudwatchRoleArn"), Value: aws.String("")},
				{Op: types.OpRemove, Path: aws.String("/features"), Value: aws.String("ProbeFeature")},
			},
		})
		if aws.ToString(updated.CloudwatchRoleArn) != "arn:aws:iam::123456789012:role/apigateway-metrics" {
			return fmt.Errorf("cloudwatchRoleArn not stored: %v", updated.CloudwatchRoleArn)
		}
		probeFeature := false
		for _, f := range updated.Features {
			if f == "ProbeFeature" {
				probeFeature = true
			}
		}
		if !probeFeature {
			return fmt.Errorf("added feature missing: %v", updated.Features)
		}

		reread, err := tc.client.GetAccount(tc.ctx, &apigateway.GetAccountInput{})
		if err != nil {
			return err
		}
		if aws.ToString(reread.CloudwatchRoleArn) != "arn:aws:iam::123456789012:role/apigateway-metrics" {
			return fmt.Errorf("cloudwatchRoleArn did not persist: %v", reread.CloudwatchRoleArn)
		}

		if _, err := tc.client.UpdateAccount(tc.ctx, &apigateway.UpdateAccountInput{
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/features"), Value: aws.String("UsagePlans")},
			},
		}); err == nil {
			return fmt.Errorf("removing the UsagePlans feature was accepted")
		} else if aerr := AssertErrorContains(err, "BadRequestException"); aerr != nil {
			return fmt.Errorf("expected BadRequestException for UsagePlans removal, got: %v", aerr)
		}

		if _, err := tc.client.UpdateAccount(tc.ctx, &apigateway.UpdateAccountInput{
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/apiKeyVersion"), Value: aws.String("2")},
			},
		}); err == nil {
			return fmt.Errorf("patching an unsupported path was accepted")
		} else if aerr := AssertErrorContains(err, "BadRequestException"); aerr != nil {
			return fmt.Errorf("expected BadRequestException for the unsupported path, got: %v", aerr)
		}
		return nil
	}))

	return results
}

package testutil

import (
	"fmt"
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

	return results
}

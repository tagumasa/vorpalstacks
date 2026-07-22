package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayEdgeTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "TagResource_UntagResource_ListTags", func() error {
		tagAPI := fmt.Sprintf("TagAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tagAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		arn := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s", r.region, *createResp.Id)

		_, err = client.TagResource(ctx, &apigateway.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		tagResp, err := client.GetTags(ctx, &apigateway.GetTagsInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get tags: %v", err)
		}
		if tagResp.Tags == nil || tagResp.Tags["key1"] != "value1" {
			return fmt.Errorf("tags mismatch, got %v", tagResp.Tags)
		}

		_, err = client.UntagResource(ctx, &apigateway.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"key2"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tagResp2, err := client.GetTags(ctx, &apigateway.GetTagsInput{
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
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApi_NonExistent", func() error {
		_, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRestApi_NonExistent", func() error {
		_, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStage_NonExistent", func() error {
		tmpAPI := fmt.Sprintf("TmpAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tmpAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("nonexistent_stage"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runAPIGatewayDeepAuditTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "PutMethod_AuthorizationScopes_RoundTrip", func() error {
		tmpAPI := fmt.Sprintf("AuthScopeAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tmpAPI),
		})
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		var rootId string
		resources, err := client.GetResources(ctx, &apigateway.GetResourcesInput{RestApiId: createResp.Id})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		for _, res := range resources.Items {
			if res.Path != nil && *res.Path == "/" {
				rootId = *res.Id
				break
			}
		}
		if rootId == "" {
			return fmt.Errorf("root resource not found")
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:           createResp.Id,
			ResourceId:          &rootId,
			HttpMethod:          aws.String("GET"),
			AuthorizationType:   aws.String("NONE"),
			AuthorizationScopes: []string{"read", "write"},
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		method, err := client.GetMethod(ctx, &apigateway.GetMethodInput{
			RestApiId:  createResp.Id,
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
		tmpAPI := fmt.Sprintf("TlsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tmpAPI),
		})
		if err != nil {
			return fmt.Errorf("create api: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		var rootId string
		resources, err := client.GetResources(ctx, &apigateway.GetResourcesInput{RestApiId: createResp.Id})
		if err != nil {
			return fmt.Errorf("get resources: %v", err)
		}
		for _, res := range resources.Items {
			if res.Path != nil && *res.Path == "/" {
				rootId = *res.Id
				break
			}
		}
		if rootId == "" {
			return fmt.Errorf("root resource not found")
		}

		_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
			RestApiId:         createResp.Id,
			ResourceId:        &rootId,
			HttpMethod:        aws.String("GET"),
			AuthorizationType: aws.String("NONE"),
		})
		if err != nil {
			return fmt.Errorf("put method: %v", err)
		}

		_, err = client.PutIntegration(ctx, &apigateway.PutIntegrationInput{
			RestApiId:  createResp.Id,
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

		integ, err := client.GetIntegration(ctx, &apigateway.GetIntegrationInput{
			RestApiId:  createResp.Id,
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

	results = append(results, r.RunTest("apigateway", "TagResource_UsagePlan", func() error {
		planName := fmt.Sprintf("TagPlan-%d", time.Now().UnixNano())
		createResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(planName),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: createResp.Id})

		arn := fmt.Sprintf("arn:aws:apigateway:%s::/usageplans/%s", r.region, *createResp.Id)

		_, err = client.TagResource(ctx, &apigateway.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags: map[string]string{
				"env":  "test",
				"team": "qa",
			},
		})
		if err != nil {
			return fmt.Errorf("tag usage plan: %v", err)
		}

		tagResp, err := client.GetTags(ctx, &apigateway.GetTagsInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get tags: %v", err)
		}
		if tagResp.Tags["env"] != "test" || tagResp.Tags["team"] != "qa" {
			return fmt.Errorf("tags mismatch: %v", tagResp.Tags)
		}

		_, err = client.UntagResource(ctx, &apigateway.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"team"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tagResp2, err := client.GetTags(ctx, &apigateway.GetTagsInput{
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
	}))

	results = append(results, r.RunTest("apigateway", "CreateUsagePlan_QuotaOffset", func() error {
		planName := fmt.Sprintf("OffsetPlan-%d", time.Now().UnixNano())
		createResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(planName),
			Quota: &types.QuotaSettings{
				Limit:  1000,
				Offset: 10,
				Period: types.QuotaPeriodTypeMonth,
			},
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: createResp.Id})

		plan, err := client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
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
		apiName := fmt.Sprintf("DisableExec-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name:                      aws.String(apiName),
			DisableExecuteApiEndpoint: true,
			MinimumCompressionSize:    aws.Int32(0),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		api, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
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

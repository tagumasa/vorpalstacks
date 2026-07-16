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

func (r *TestRunner) runAPIGatewayUsagePlanTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "CreateUsagePlan", func() error {
		resp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name:        aws.String("test-usage-plan"),
			Description: aws.String("Test usage plan"),
			Throttle: &types.ThrottleSettings{
				BurstLimit: 10,
				RateLimit:  5.0,
			},
			Quota: &types.QuotaSettings{
				Limit:  1000,
				Period: types.QuotaPeriodTypeMonth,
			},
			Tags: map[string]string{
				"team": "backend",
			},
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("usage plan ID is nil")
		}
		if resp.Name == nil || *resp.Name != "test-usage-plan" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Description == nil || *resp.Description != "Test usage plan" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		if resp.Throttle == nil || resp.Throttle.BurstLimit != 10 {
			return fmt.Errorf("throttle burstLimit mismatch")
		}
		if resp.Throttle.RateLimit != 5.0 {
			return fmt.Errorf("throttle rateLimit mismatch, got %v", resp.Throttle.RateLimit)
		}
		if resp.Quota == nil || resp.Quota.Period != types.QuotaPeriodTypeMonth {
			return fmt.Errorf("quota period mismatch")
		}
		if resp.Quota.Limit != 1000 {
			return fmt.Errorf("quota limit mismatch, got %v", resp.Quota.Limit)
		}
		if resp.Tags == nil || resp.Tags["team"] != "backend" {
			return fmt.Errorf("tags mismatch, got %v", resp.Tags)
		}
		return nil
	}))

	var usagePlanID string
	results = append(results, r.RunTest("apigateway", "GetUsagePlans", func() error {
		resp, err := client.GetUsagePlans(ctx, &apigateway.GetUsagePlansInput{
			Limit: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan")
		}
		for _, item := range resp.Items {
			if item.Name != nil && *item.Name == "test-usage-plan" {
				usagePlanID = *item.Id
				break
			}
		}
		if usagePlanID == "" {
			return fmt.Errorf("test-usage-plan not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		resp, err := client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "test-usage-plan" {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.Throttle == nil || resp.Throttle.BurstLimit != 10 {
			return fmt.Errorf("throttle burstLimit mismatch, got %v", resp.Throttle)
		}
		if resp.Quota == nil || resp.Quota.Limit != 1000 {
			return fmt.Errorf("quota limit mismatch, got %v", resp.Quota)
		}
		if resp.Description == nil || *resp.Description != "Test usage plan" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		resp, err := client.UpdateUsagePlan(ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/name"),
					Value: aws.String("updated-usage-plan"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "updated-usage-plan" {
			return fmt.Errorf("name not updated, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		_, err := client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
		})
		if err == nil {
			return fmt.Errorf("GetUsagePlan should fail after delete")
		}
		if !strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("expected NotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateUsagePlanKey_Lifecycle", func() error {
		keyResp, err := client.CreateApiKey(ctx, &apigateway.CreateApiKeyInput{
			Name:    aws.String("upk-test-key"),
			Enabled: true,
		})
		if err != nil {
			return fmt.Errorf("create api key: %v", err)
		}
		defer client.DeleteApiKey(ctx, &apigateway.DeleteApiKeyInput{ApiKey: keyResp.Id})

		upResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String("upk-test-plan"),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		upkResp, err := client.CreateUsagePlanKey(ctx, &apigateway.CreateUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
			KeyType:     aws.String("API_KEY"),
		})
		if err != nil {
			return fmt.Errorf("create usage plan key: %v", err)
		}
		if upkResp.Id == nil {
			return fmt.Errorf("usage plan key ID is nil")
		}

		getResp, err := client.GetUsagePlanKey(ctx, &apigateway.GetUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan key: %v", err)
		}
		if getResp.Type == nil || *getResp.Type != "API_KEY" {
			return fmt.Errorf("type mismatch, got %v", getResp.Type)
		}

		keysResp, err := client.GetUsagePlanKeys(ctx, &apigateway.GetUsagePlanKeysInput{
			UsagePlanId: upResp.Id,
			Limit:       aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("get usage plan keys: %v", err)
		}
		if len(keysResp.Items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan key")
		}

		_, err = client.DeleteUsagePlanKey(ctx, &apigateway.DeleteUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete usage plan key: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsage", func() error {
		upResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("usage-plan-%d", time.Now().UnixNano())),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		now := time.Now().UTC()
		startDate := now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate := now.Format("2006-01-02")
		resp, err := client.GetUsage(ctx, &apigateway.GetUsageInput{
			UsagePlanId: upResp.Id,
			StartDate:   aws.String(startDate),
			EndDate:     aws.String(endDate),
		})
		if err != nil {
			return err
		}
		if resp.UsagePlanId == nil || *resp.UsagePlanId != *upResp.Id {
			return fmt.Errorf("usagePlanId mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UsagePlan_WithApiStages", func() error {
		usAPI := fmt.Sprintf("UsAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(usAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		depResp, err := client.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
			RestApiId: createResp.Id,
		})
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		_, err = client.CreateStage(ctx, &apigateway.CreateStageInput{
			RestApiId:    createResp.Id,
			StageName:    aws.String("api-stage"),
			DeploymentId: depResp.Id,
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		upResp, err := client.CreateUsagePlan(ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("us-plan-%d", time.Now().UnixNano())),
			ApiStages: []types.ApiStage{
				{
					ApiId: createResp.Id,
					Stage: aws.String("api-stage"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}

		getResp, err := client.GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: upResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		if len(getResp.ApiStages) == 0 {
			return fmt.Errorf("expected apiStages to be set")
		}

		client.DeleteUsagePlan(ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})
		return nil
	}))

	return results
}

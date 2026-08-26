package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayUsagePlanTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "CreateUsagePlan", func() error {
		resp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
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
		items, err := tc.allUsagePlans()
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan")
		}
		found := containsID(items, func(item *types.UsagePlan) bool {
			return item.Name != nil && *item.Name == "test-usage-plan"
		})
		if found == nil {
			return fmt.Errorf("test-usage-plan not found")
		}
		usagePlanID = *found.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsagePlan", func() error {
		if usagePlanID == "" {
			return fmt.Errorf("usage plan ID not available")
		}
		resp, err := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
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
		resp, err := tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
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
		_, err := tc.client.DeleteUsagePlan(tc.ctx, &apigateway.DeleteUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
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
		keyResp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
			Name:    aws.String("upk-test-key"),
			Enabled: true,
		})
		if err != nil {
			return fmt.Errorf("create api key: %v", err)
		}
		defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: keyResp.Id})

		upResp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String("upk-test-plan"),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.client.DeleteUsagePlan(tc.ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		upkResp, err := tc.client.CreateUsagePlanKey(tc.ctx, &apigateway.CreateUsagePlanKeyInput{
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

		getResp, err := tc.client.GetUsagePlanKey(tc.ctx, &apigateway.GetUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan key: %v", err)
		}
		if getResp.Type == nil || *getResp.Type != "API_KEY" {
			return fmt.Errorf("type mismatch, got %v", getResp.Type)
		}

		keysResp, err := tc.client.GetUsagePlanKeys(tc.ctx, &apigateway.GetUsagePlanKeysInput{
			UsagePlanId: upResp.Id,
			Limit:       aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("get usage plan keys: %v", err)
		}
		if len(keysResp.Items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan key")
		}

		_, err = tc.client.DeleteUsagePlanKey(tc.ctx, &apigateway.DeleteUsagePlanKeyInput{
			UsagePlanId: upResp.Id,
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete usage plan key: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsage", func() error {
		upResp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("usage-plan-%d", time.Now().UnixNano())),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.client.DeleteUsagePlan(tc.ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		now := time.Now().UTC()
		startDate := now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate := now.Format("2006-01-02")
		resp, err := tc.client.GetUsage(tc.ctx, &apigateway.GetUsageInput{
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
		ownAPI, _, err := tc.createAPI(tc.uniqueName("UsAPI"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(ownAPI)
		depID, err := tc.createDeployment(ownAPI, "")
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		_, err = tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(ownAPI),
			StageName:    aws.String("api-stage"),
			DeploymentId: aws.String(depID),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		upResp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("us-plan-%d", time.Now().UnixNano())),
			ApiStages: []types.ApiStage{
				{
					ApiId: aws.String(ownAPI),
					Stage: aws.String("api-stage"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}

		getResp, err := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: upResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		if len(getResp.ApiStages) == 0 {
			return fmt.Errorf("expected apiStages to be set")
		}

		tc.client.DeleteUsagePlan(tc.ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateUsagePlan_ApiStageIndexOutOfRange_Rejected", func() error {
		upResp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(fmt.Sprintf("oob-plan-%d", time.Now().UnixNano())),
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.client.DeleteUsagePlan(tc.ctx, &apigateway.DeleteUsagePlanInput{UsagePlanId: upResp.Id})

		// The plan has no api stages; index 0 is out of range and must be
		// rejected instead of silently appending an empty stage entry.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: upResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/apiStages/0/apiId"), Value: aws.String("abc123")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for out-of-range apiStages index, got: %v", err)
		}

		getResp, getErr := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: upResp.Id,
		})
		if getErr != nil {
			return fmt.Errorf("get usage plan: %v", getErr)
		}
		if len(getResp.ApiStages) != 0 {
			return fmt.Errorf("expected no api stages after rejected update, got %d", len(getResp.ApiStages))
		}
		return nil
	}))

	return results
}

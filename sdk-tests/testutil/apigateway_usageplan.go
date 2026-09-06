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
		if err := tc.require(usagePlanID); err != nil {
			return err
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
		if err := tc.require(usagePlanID); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "UpdateUsagePlan_WholeMemberRemoves", func() error {
		if err := tc.require(usagePlanID); err != nil {
			return err
		}
		resp, err := tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/quota")},
				{Op: types.OpRemove, Path: aws.String("/throttle")},
			},
		})
		if err != nil {
			return err
		}
		if resp.Quota != nil {
			return fmt.Errorf("whole-member quota remove did not clear the settings, got %v", resp.Quota)
		}
		if resp.Throttle != nil {
			return fmt.Errorf("whole-member throttle remove did not clear the settings, got %v", resp.Throttle)
		}

		// The whole-member rows document remove only.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/quota"), Value: aws.String(`{}`)},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for whole-member quota replace, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteUsagePlan", func() error {
		if err := tc.require(usagePlanID); err != nil {
			return err
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

	results = append(results, r.RunTest("apigateway", "CreateUsagePlan_InvalidPerStageThrottle_Rejected", func() error {
		_, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name: aws.String(tc.uniqueName("bad-throttle")),
			ApiStages: []types.ApiStage{
				{
					ApiId: aws.String(tc.apiID),
					Stage: aws.String("test"),
					Throttle: map[string]types.ThrottleSettings{
						"GET": {BurstLimit: 20000},
					},
				},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for per-stage throttle over the limit, got: %v", err)
		}
		if !strings.Contains(fmt.Sprintf("%v", err), "per-stage throttle burstLimit") {
			return fmt.Errorf("expected the per-stage throttle message, got: %v", err)
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

		planID, err := tc.createOwnUsagePlan("upk-test-plan")
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(planID)

		upkResp, err := tc.client.CreateUsagePlanKey(tc.ctx, &apigateway.CreateUsagePlanKeyInput{
			UsagePlanId: aws.String(planID),
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
			UsagePlanId: aws.String(planID),
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan key: %v", err)
		}
		if getResp.Type == nil || *getResp.Type != "API_KEY" {
			return fmt.Errorf("type mismatch, got %v", getResp.Type)
		}

		keysResp, err := tc.client.GetUsagePlanKeys(tc.ctx, &apigateway.GetUsagePlanKeysInput{
			UsagePlanId: aws.String(planID),
			Limit:       aws.Int32(100),
		})
		if err != nil {
			return fmt.Errorf("get usage plan keys: %v", err)
		}
		if len(keysResp.Items) == 0 {
			return fmt.Errorf("expected at least 1 usage plan key")
		}

		_, err = tc.client.DeleteUsagePlanKey(tc.ctx, &apigateway.DeleteUsagePlanKeyInput{
			UsagePlanId: aws.String(planID),
			KeyId:       keyResp.Id,
		})
		if err != nil {
			return fmt.Errorf("delete usage plan key: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetUsage", func() error {
		planID, err := tc.createOwnUsagePlan("usage-plan")
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(planID)

		now := time.Now().UTC()
		startDate := now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate := now.Format("2006-01-02")
		resp, err := tc.client.GetUsage(tc.ctx, &apigateway.GetUsageInput{
			UsagePlanId: aws.String(planID),
			StartDate:   aws.String(startDate),
			EndDate:     aws.String(endDate),
		})
		if err != nil {
			return err
		}
		if resp.UsagePlanId == nil || *resp.UsagePlanId != planID {
			return fmt.Errorf("usagePlanId mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UsagePlan_WithApiStages", func() error {
		ownAPI, _, err := tc.createOwnAPI("UsAPI")
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
			Name: aws.String(tc.uniqueName("us-plan")),
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

		tc.deleteUsagePlan(aws.ToString(upResp.Id))
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateUsagePlan_ApiStagesPatches", func() error {
		planID, err := tc.createOwnUsagePlan("apsp-plan")
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(planID)

		// The whole-member /apiStages row documents add and remove; the
		// developer guide example carries the value as apiId:stageName.
		upd, err := tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/apiStages"), Value: aws.String(tc.apiID + ":test")},
			},
		})
		if err != nil {
			return err
		}
		if len(upd.ApiStages) != 1 || aws.ToString(upd.ApiStages[0].ApiId) != tc.apiID || aws.ToString(upd.ApiStages[0].Stage) != "test" {
			return fmt.Errorf("whole-member apiStages add not applied, got %+v", upd.ApiStages)
		}

		addr := "/apiStages/" + tc.apiID + ":test"

		// Whole-throttle replace with the documented JSON object form,
		// keyed by {resourcePath}/{httpMethod} — "//GET" is the API root.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String(addr + "/throttle"), Value: aws.String(`{"//GET":{"rateLimit":1,"burstLimit":2}}`)},
			},
		})
		if err != nil {
			return err
		}
		assertThrottle := func(rate, burst float64) error {
			getResp, err := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
				UsagePlanId: aws.String(planID),
			})
			if err != nil {
				return fmt.Errorf("get usage plan: %v", err)
			}
			if len(getResp.ApiStages) != 1 {
				return fmt.Errorf("api stage lost, got %+v", getResp.ApiStages)
			}
			ts, ok := getResp.ApiStages[0].Throttle["//GET"]
			if !ok {
				return fmt.Errorf("throttle key //GET missing, got %+v", getResp.ApiStages[0].Throttle)
			}
			if ts.RateLimit != rate || float64(ts.BurstLimit) != burst {
				return fmt.Errorf("throttle entry mismatch, got %+v", ts)
			}
			return nil
		}
		if err := assertThrottle(1, 2); err != nil {
			return err
		}

		// A single-method rateLimit patch updates the entry.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String(addr + "/throttle//GET/rateLimit"), Value: aws.String("0.5")},
			},
		})
		if err != nil {
			return err
		}
		if err := assertThrottle(0.5, 2); err != nil {
			return err
		}

		// Removing the method throttle key clears just that key.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String(addr + "/throttle//GET")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(planID),
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		if len(getResp.ApiStages) != 1 || len(getResp.ApiStages[0].Throttle) != 0 {
			return fmt.Errorf("method throttle remove did not clear the key, got %+v", getResp.ApiStages)
		}

		// The escaped resource path token stays as addressed in the stored
		// key ("~1items~1{id}/GET"), matching the method-keyed map
		// convention the official CLI update-stage example output shows.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String(addr + "/throttle/~1items~1{id}/GET/rateLimit"), Value: aws.String("2")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err = tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(planID),
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		escaped, ok := getResp.ApiStages[0].Throttle["~1items~1{id}/GET"]
		if !ok {
			return fmt.Errorf("escaped-key entry missing, got %+v", getResp.ApiStages[0].Throttle)
		}
		if escaped.RateLimit != 2 {
			return fmt.Errorf("escaped-key rateLimit mismatch, got %+v", escaped)
		}
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String(addr + "/throttle/~1items~1{id}/GET")},
			},
		})
		if err != nil {
			return err
		}

		// Addressing a stage the plan does not carry rejects.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/apiStages/" + tc.apiID + ":nope/throttle"), Value: aws.String(`{}`)},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for an unknown api stage address, got: %v", err)
		}

		// The whole-member remove with the value drops the entry.
		_, err = tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(planID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/apiStages"), Value: aws.String(tc.apiID + ":test")},
			},
		})
		if err != nil {
			return err
		}
		getResp, err = tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(planID),
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		if len(getResp.ApiStages) != 0 {
			return fmt.Errorf("whole-member apiStages remove did not drop the entry, got %+v", getResp.ApiStages)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateUsagePlan_ApiStageIndexForms_Rejected", func() error {
		planID, err := tc.createOwnUsagePlan("idx-plan")
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(planID)

		// Numeric index addressing appears nowhere in the official patch
		// tables: the documented element addressing is
		// /apiStages/{apiId}:{stageName}/....
		for _, po := range []types.PatchOperation{
			{Op: types.OpReplace, Path: aws.String("/apiStages/0/apiId"), Value: aws.String("abc123")},
			{Op: types.OpAdd, Path: aws.String("/apiStages/0"), Value: aws.String("abc123:test")},
		} {
			_, err := tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
				UsagePlanId:     aws.String(planID),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s, got: %v", po.Op, *po.Path, err)
			}
		}

		getResp, getErr := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: aws.String(planID),
		})
		if getErr != nil {
			return fmt.Errorf("get usage plan: %v", getErr)
		}
		if len(getResp.ApiStages) != 0 {
			return fmt.Errorf("expected no api stages after rejected updates, got %d", len(getResp.ApiStages))
		}
		return nil
	}))

	return results
}

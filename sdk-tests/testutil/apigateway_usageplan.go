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

		// keyId restricts the listing to the plans associated with the key.
		keyResp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
			Name: aws.String(tc.uniqueName("keyid-filter-key")),
		})
		if err != nil {
			return fmt.Errorf("create api key: %v", err)
		}
		defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: keyResp.Id})
		ownPlan, err := tc.createOwnUsagePlan("keyid-filter-plan")
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		defer tc.deleteUsagePlan(ownPlan)
		if _, err := tc.client.CreateUsagePlanKey(tc.ctx, &apigateway.CreateUsagePlanKeyInput{
			UsagePlanId: aws.String(ownPlan),
			KeyId:       keyResp.Id,
			KeyType:     aws.String("API_KEY"),
		}); err != nil {
			return fmt.Errorf("create usage plan key: %v", err)
		}
		byKey, err := tc.client.GetUsagePlans(tc.ctx, &apigateway.GetUsagePlansInput{
			KeyId: keyResp.Id,
		})
		if err != nil {
			return err
		}
		if len(byKey.Items) != 1 || aws.ToString(byKey.Items[0].Id) != ownPlan {
			return fmt.Errorf("keyId filter did not match exactly the associated plan, got %+v", byKey.Items)
		}
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

		// The scalar rows of the official patch table: description, the
		// plan-level throttle members, and the quota members.
		upd, err := tc.client.UpdateUsagePlan(tc.ctx, &apigateway.UpdateUsagePlanInput{
			UsagePlanId: aws.String(usagePlanID),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/description"), Value: aws.String("scalar rows")},
				{Op: types.OpReplace, Path: aws.String("/throttle/burstLimit"), Value: aws.String("700")},
				{Op: types.OpReplace, Path: aws.String("/throttle/rateLimit"), Value: aws.String("800.5")},
				{Op: types.OpReplace, Path: aws.String("/quota/limit"), Value: aws.String("500")},
				{Op: types.OpReplace, Path: aws.String("/quota/offset"), Value: aws.String("1")},
				{Op: types.OpReplace, Path: aws.String("/quota/period"), Value: aws.String("WEEK")},
			},
		})
		if err != nil {
			return fmt.Errorf("scalar row patch: %v", err)
		}
		if aws.ToString(upd.Description) != "scalar rows" {
			return fmt.Errorf("description row not applied, got %v", upd.Description)
		}
		if upd.Throttle == nil || upd.Throttle.BurstLimit != 700 || upd.Throttle.RateLimit != 800.5 {
			return fmt.Errorf("throttle rows not applied, got %+v", upd.Throttle)
		}
		if upd.Quota == nil || upd.Quota.Limit != 500 || upd.Quota.Offset != 1 ||
			string(upd.Quota.Period) != "WEEK" {
			return fmt.Errorf("quota rows not applied, got %+v", upd.Quota)
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

		// A forced multi-page walk: two more keys bring the plan past a
		// two-key page limit, and every page's position feeds back.
		extra := make([]string, 0, 2)
		for i := 0; i < 2; i++ {
			extraKey, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
				Name: aws.String(fmt.Sprintf("upk-page-key-%d", i)),
			})
			if err != nil {
				return fmt.Errorf("create extra key %d: %v", i, err)
			}
			defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: extraKey.Id})
			if _, err := tc.client.CreateUsagePlanKey(tc.ctx, &apigateway.CreateUsagePlanKeyInput{
				UsagePlanId: aws.String(planID),
				KeyId:       extraKey.Id,
				KeyType:     aws.String("API_KEY"),
			}); err != nil {
				return fmt.Errorf("attach extra key %d: %v", i, err)
			}
			extra = append(extra, aws.ToString(extraKey.Id))
		}
		want := append([]string{aws.ToString(keyResp.Id)}, extra...)
		var got []string
		var position *string
		pages := 0
		for {
			page, err := tc.client.GetUsagePlanKeys(tc.ctx, &apigateway.GetUsagePlanKeysInput{
				UsagePlanId: aws.String(planID),
				Limit:       aws.Int32(2),
				Position:    position,
			})
			if err != nil {
				return fmt.Errorf("page %d: %v", pages, err)
			}
			pages++
			for _, k := range page.Items {
				got = append(got, aws.ToString(k.Id))
			}
			position = page.Position
			if position == nil {
				break
			}
		}
		if pages < 2 {
			return fmt.Errorf("expected the page limit to force multiple pages, got %d", pages)
		}
		if len(got) != len(want) {
			return fmt.Errorf("walked %d keys %v, want %d %v", len(got), got, len(want), want)
		}
		gotSet := make(map[string]bool, len(got))
		for _, id := range got {
			gotSet[id] = true
		}
		for _, id := range want {
			if !gotSet[id] {
				return fmt.Errorf("created key %q missing from the walk, got %v", id, got)
			}
		}
		for _, id := range extra {
			if _, err := tc.client.DeleteUsagePlanKey(tc.ctx, &apigateway.DeleteUsagePlanKeyInput{
				UsagePlanId: aws.String(planID),
				KeyId:       aws.String(id),
			}); err != nil {
				return fmt.Errorf("delete extra key association %s: %v", id, err)
			}
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
		keyResp, err := tc.client.CreateApiKey(tc.ctx, &apigateway.CreateApiKeyInput{
			Name: aws.String(tc.uniqueName("usage-shape-key")),
		})
		if err != nil {
			return fmt.Errorf("create api key: %v", err)
		}
		defer tc.client.DeleteApiKey(tc.ctx, &apigateway.DeleteApiKeyInput{ApiKey: keyResp.Id})

		// The Usage shape reports a quota draw-down, so the plan carries a
		// daily quota to draw from.
		planResp, err := tc.client.CreateUsagePlan(tc.ctx, &apigateway.CreateUsagePlanInput{
			Name:  aws.String(tc.uniqueName("usage-shape-plan")),
			Quota: &types.QuotaSettings{Limit: 100, Period: types.QuotaPeriodTypeDay},
		})
		if err != nil {
			return fmt.Errorf("create usage plan: %v", err)
		}
		planID := *planResp.Id
		defer tc.deleteUsagePlan(planID)

		if _, err := tc.client.CreateUsagePlanKey(tc.ctx, &apigateway.CreateUsagePlanKeyInput{
			UsagePlanId: aws.String(planID),
			KeyId:       keyResp.Id,
			KeyType:     aws.String("API_KEY"),
		}); err != nil {
			return fmt.Errorf("create usage plan key: %v", err)
		}

		today := time.Now().Format("2006-01-02")
		resp, err := tc.client.GetUsage(tc.ctx, &apigateway.GetUsageInput{
			UsagePlanId: aws.String(planID),
			StartDate:   aws.String(today),
			EndDate:     aws.String(today),
		})
		if err != nil {
			return err
		}
		if resp.UsagePlanId == nil || *resp.UsagePlanId != planID {
			return fmt.Errorf("usagePlanId mismatch")
		}
		usage, ok := resp.Items[*keyResp.Id]
		if !ok {
			return fmt.Errorf("values not indexed by API key id %s: %v", *keyResp.Id, resp.Items)
		}
		if len(usage) != 1 || len(usage[0]) != 2 || usage[0][0] != 0 || usage[0][1] != 100 {
			return fmt.Errorf("expected daily log [0 100] for an unused key, got %v", usage)
		}

		// A replace of /remaining with an integer is the only documented
		// patch operation for usage.
		upd, err := tc.client.UpdateUsage(tc.ctx, &apigateway.UpdateUsageInput{
			UsagePlanId: aws.String(planID),
			KeyId:       keyResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/remaining"), Value: aws.String("10")},
			},
		})
		if err != nil {
			return err
		}
		if upd.StartDate == nil || *upd.StartDate != today || upd.EndDate == nil || *upd.EndDate != today {
			return fmt.Errorf("update usage response not scoped to the grant date, got %v..%v", upd.StartDate, upd.EndDate)
		}
		updUsage, ok := upd.Items[*keyResp.Id]
		if !ok || len(updUsage) != 1 || updUsage[0][0] != 0 || updUsage[0][1] != 10 {
			return fmt.Errorf("granted remaining 10 not reported, got %v", upd.Items)
		}

		after, err := tc.client.GetUsage(tc.ctx, &apigateway.GetUsageInput{
			UsagePlanId: aws.String(planID),
			KeyId:       keyResp.Id,
			StartDate:   aws.String(today),
			EndDate:     aws.String(today),
		})
		if err != nil {
			return err
		}
		afterUsage, ok := after.Items[*keyResp.Id]
		if !ok || len(afterUsage) != 1 || afterUsage[0][0] != 0 || afterUsage[0][1] != 10 {
			return fmt.Errorf("granted override not reflected in GetUsage, got %v", after.Items)
		}

		_, err = tc.client.UpdateUsage(tc.ctx, &apigateway.UpdateUsageInput{
			UsagePlanId: aws.String(planID),
			KeyId:       keyResp.Id,
			PatchOperations: []types.PatchOperation{
				{Op: types.OpAdd, Path: aws.String("/remaining"), Value: aws.String("5")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for add on /remaining, got: %v", err)
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
		defer tc.deleteUsagePlan(aws.ToString(upResp.Id))

		getResp, err := tc.client.GetUsagePlan(tc.ctx, &apigateway.GetUsagePlanInput{
			UsagePlanId: upResp.Id,
		})
		if err != nil {
			return fmt.Errorf("get usage plan: %v", err)
		}
		if len(getResp.ApiStages) == 0 {
			return fmt.Errorf("expected apiStages to be set")
		}
		if aws.ToString(getResp.ApiStages[0].ApiId) != ownAPI || aws.ToString(getResp.ApiStages[0].Stage) != "api-stage" {
			return fmt.Errorf("apiStage mismatch, got %+v", getResp.ApiStages)
		}
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

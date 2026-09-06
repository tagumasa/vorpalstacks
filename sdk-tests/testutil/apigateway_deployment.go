package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func (r *TestRunner) runAPIGatewayDeploymentTests(tc *apigwTestContext) []TestResult {
	var results []TestResult

	var deploymentID string
	results = append(results, r.RunTest("apigateway", "CreateDeployment", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.CreateDeployment(tc.ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   aws.String(tc.apiID),
			Description: aws.String("test deployment"),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil {
			return fmt.Errorf("deployment ID is nil")
		}
		if resp.Description == nil || *resp.Description != "test deployment" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		deploymentID = *resp.Id
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployment", func() error {
		if err := tc.require(tc.apiID, deploymentID); err != nil {
			return err
		}
		resp, err := tc.client.GetDeployment(tc.ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return err
		}
		if resp.Id == nil || *resp.Id != deploymentID {
			return fmt.Errorf("deployment ID mismatch, got %v", resp.Id)
		}
		if resp.CreatedDate == nil {
			return fmt.Errorf("createdDate is nil")
		}
		if resp.Description == nil || *resp.Description != "test deployment" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateDeployment", func() error {
		if err := tc.require(tc.apiID, deploymentID); err != nil {
			return err
		}
		resp, err := tc.client.UpdateDeployment(tc.ctx, &apigateway.UpdateDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated deployment"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated deployment" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}

		// Verify the description change persists via a fresh read.
		getResp, err := tc.client.GetDeployment(tc.ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return fmt.Errorf("get deployment: %v", err)
		}
		if getResp.Description == nil || *getResp.Description != "updated deployment" {
			return fmt.Errorf("deployment description mismatch, got %v", getResp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetDeployments", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.GetDeployments(tc.ctx, &apigateway.GetDeploymentsInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return fmt.Errorf("expected at least 1 deployment")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "CreateStage", func() error {
		if err := tc.require(tc.apiID, deploymentID); err != nil {
			return err
		}
		resp, err := tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(tc.apiID),
			StageName:    aws.String("test"),
			DeploymentId: aws.String(deploymentID),
			Description:  aws.String("test stage"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.StageName == nil || *resp.StageName != "test" {
			return fmt.Errorf("stageName mismatch, got %v", resp.StageName)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != deploymentID {
			return fmt.Errorf("deploymentId mismatch, got %v", resp.DeploymentId)
		}
		if resp.Description == nil || *resp.Description != "test stage" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}

		// The documented stage variable constraints hold at create
		// ingress: names are alphanumeric/underscore, values match the
		// documented charset and cannot be empty, and the canary's stage
		// variable overrides are stage variables too.
		for _, input := range []apigateway.CreateStageInput{
			{
				RestApiId:    aws.String(tc.apiID),
				StageName:    aws.String("badvarname"),
				DeploymentId: aws.String(deploymentID),
				Variables:    map[string]string{"has space": "a b"},
			},
			{
				RestApiId:    aws.String(tc.apiID),
				StageName:    aws.String("badvarvalue"),
				DeploymentId: aws.String(deploymentID),
				Variables:    map[string]string{"env": "has space"},
			},
			{
				RestApiId:    aws.String(tc.apiID),
				StageName:    aws.String("emptyvarvalue"),
				DeploymentId: aws.String(deploymentID),
				Variables:    map[string]string{"env": ""},
			},
			{
				RestApiId:    aws.String(tc.apiID),
				StageName:    aws.String("badcanaryvars"),
				DeploymentId: aws.String(deploymentID),
				CanarySettings: &types.CanarySettings{
					StageVariableOverrides: map[string]string{"a/b": "v"},
				},
			},
		} {
			_, err := tc.client.CreateStage(tc.ctx, &input)
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for stage variable validation on stage %s, got: %v", *input.StageName, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStage", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return err
		}
		if resp.StageName == nil || *resp.StageName != "test" {
			return fmt.Errorf("stage name mismatch, got %v", resp.StageName)
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != deploymentID {
			return fmt.Errorf("deploymentId mismatch, got %v", resp.DeploymentId)
		}
		if resp.Description == nil || *resp.Description != "test stage" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStages", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.GetStages(tc.ctx, &apigateway.GetStagesInput{
			RestApiId: aws.String(tc.apiID),
		})
		if err != nil {
			return err
		}
		if len(resp.Item) == 0 {
			return fmt.Errorf("expected at least 1 stage")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		resp, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/description"),
					Value: aws.String("updated stage"),
				},
				{
					Op:    types.OpReplace,
					Path:  aws.String("/variables/env"),
					Value: aws.String("prod"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "updated stage" {
			return fmt.Errorf("description not updated, got %v", resp.Description)
		}

		// Verify the stage variable round-trips via a fresh read.
		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if stageResp.Variables == nil || stageResp.Variables["env"] != "prod" {
			return fmt.Errorf("stage variables not set, got %v", stageResp.Variables)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_VariablePatchConstraints", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		// The /variables/* row of the official patch table documents
		// replace only: per-key add and remove reject — on the stage
		// variables and on the canary overrides, which are stage
		// variables too.
		for _, po := range []types.PatchOperation{
			{Op: types.OpAdd, Path: aws.String("/variables/lambda_alias"), Value: aws.String("v1")},
			{Op: types.OpRemove, Path: aws.String("/variables/env")},
			{Op: types.OpAdd, Path: aws.String("/canarySettings/stageVariableOverrides/extra"), Value: aws.String("v1")},
		} {
			_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
				RestApiId:       aws.String(tc.apiID),
				StageName:       aws.String("test"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s, got: %v", po.Op, *po.Path, err)
			}
		}

		// A per-key replace with a legal name round-trips, on the stage
		// variables and on the canary overrides.
		_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/variables/feature_x"), Value: aws.String("on")},
				{Op: types.OpReplace, Path: aws.String("/canarySettings/stageVariableOverrides/ov1"), Value: aws.String("v1")},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if stageResp.Variables["feature_x"] != "on" || stageResp.Variables["env"] != "prod" {
			return fmt.Errorf("per-key variable replace not applied, got %v", stageResp.Variables)
		}
		if stageResp.CanarySettings == nil ||
			stageResp.CanarySettings.StageVariableOverrides["ov1"] != "v1" {
			return fmt.Errorf("canary override replace not applied, got %+v", stageResp.CanarySettings)
		}

		// a~1b unescapes to a/b, which the documented stage variable name
		// charset rejects even for replace.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/variables/a~1b"), Value: aws.String("v")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for charset-invalid variable name, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_MethodSettingsPatches", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		// The documented setting family addresses members as
		// /{resourcePath}/{httpMethod}/{group}/{member}; the returned map is
		// keyed by the as-addressed pointer token joined to the method (the
		// official CLI example output shows "~1resourceName/GET" and the
		// wildcard as "*/*").
		_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/~1pets~1{petId}/GET/logging/loglevel"),
					Value: aws.String("INFO"),
				},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		setting, ok := stageResp.MethodSettings["~1pets~1{petId}/GET"]
		if !ok {
			return fmt.Errorf("as-addressed method settings key missing, got %v", stageResp.MethodSettings)
		}
		if setting.LoggingLevel == nil || *setting.LoggingLevel != "INFO" {
			return fmt.Errorf("loggingLevel not set, got %v", setting.LoggingLevel)
		}

		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/*/*/metrics/enabled"),
					Value: aws.String("true"),
				},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		wildcard, ok := stageResp.MethodSettings["*/*"]
		if !ok {
			return fmt.Errorf("wildcard method settings key missing, got %v", stageResp.MethodSettings)
		}
		if !wildcard.MetricsEnabled {
			return fmt.Errorf("wildcard metricsEnabled not set, got %v", wildcard.MetricsEnabled)
		}

		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/~1pets/GET/caching/unauthorizedCacheControlHeaderStrategy"),
					Value: aws.String("SUCCEED_WITHOUT_RESPONSE_HEADER"),
				},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		strategy, ok := stageResp.MethodSettings["~1pets/GET"]
		if !ok {
			return fmt.Errorf("~1pets key missing, got %v", stageResp.MethodSettings)
		}
		if string(strategy.UnauthorizedCacheControlHeaderStrategy) != "SUCCEED_WITHOUT_RESPONSE_HEADER" {
			return fmt.Errorf("unauthorizedCacheControlHeaderStrategy not set, got %v", strategy.UnauthorizedCacheControlHeaderStrategy)
		}

		// The undocumented per-key form rejects, as does an invalid strategy
		// value and a non-replace operation on the family.
		for _, po := range []types.PatchOperation{
			{Op: types.OpReplace, Path: aws.String("/methodSettings/~1pets/GET/loggingLevel"), Value: aws.String("INFO")},
			{Op: types.OpReplace, Path: aws.String("/~1pets/GET/caching/unauthorizedCacheControlHeaderStrategy"), Value: aws.String("FAIL_WITH_40X")},
			{Op: types.OpAdd, Path: aws.String("/~1pets/GET/logging/dataTrace"), Value: aws.String("true")},
		} {
			_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
				RestApiId:       aws.String(tc.apiID),
				StageName:       aws.String("test"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s, got: %v", po.Op, *po.Path, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_MalformedPatchPathTokens_Rejected", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		// A trailing slash carries an empty key token: not a documented
		// patch form, so every operation kind must fail instead of
		// writing an empty-string key.
		for _, po := range []types.PatchOperation{
			{Op: types.OpAdd, Path: aws.String("/variables/"), Value: aws.String("v")},
			{Op: types.OpAdd, Path: aws.String("/canarySettings/stageVariableOverrides/"), Value: aws.String("v")},
			// Whole-member paths support only the operations documented
			// in the official patch table.
			{Op: types.OpAdd, Path: aws.String("/variables"), Value: aws.String(`{"a":"b"}`)},
			{Op: types.OpReplace, Path: aws.String("/canarySettings")},
			{Op: types.OpRemove, Path: aws.String("/methodSettings")},
			// The whole-member variables value must respect the
			// documented stage variable name charset.
			{Op: types.OpReplace, Path: aws.String("/variables"), Value: aws.String(`{"bad name":"v"}`)},
		} {
			_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
				RestApiId:       aws.String(tc.apiID),
				StageName:       aws.String("test"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for path %q, got: %v", *po.Path, err)
			}
		}
		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if _, ok := stageResp.Variables[""]; ok {
			return fmt.Errorf("empty-string variable key stored, got %v", stageResp.Variables)
		}
		if stageResp.Variables["env"] != "prod" {
			return fmt.Errorf("rejected patches changed the variables, got %v", stageResp.Variables)
		}
		if stageResp.CanarySettings != nil {
			if _, ok := stageResp.CanarySettings.StageVariableOverrides[""]; ok {
				return fmt.Errorf("empty-string override key stored, got %v", stageResp.CanarySettings.StageVariableOverrides)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_WholeMemberPatchPaths", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		// Whole-member replace sets the map from the JSON object value.
		_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/variables"),
					Value: aws.String(`{"env2":"staging","feature_x":"on"}`),
				},
				{
					Op:    types.OpReplace,
					Path:  aws.String("/canarySettings/stageVariableOverrides"),
					Value: aws.String(`{"ov1":"v1"}`),
				},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if len(stageResp.Variables) != 2 || stageResp.Variables["env2"] != "staging" || stageResp.Variables["feature_x"] != "on" {
			return fmt.Errorf("whole-member variables replace not applied, got %v", stageResp.Variables)
		}
		if stageResp.CanarySettings == nil || len(stageResp.CanarySettings.StageVariableOverrides) != 1 ||
			stageResp.CanarySettings.StageVariableOverrides["ov1"] != "v1" {
			return fmt.Errorf("whole-member override replace not applied, got %v", stageResp.CanarySettings)
		}

		// Whole-member replace of the method settings map: keys are
		// {resource_path}/{http_method} with real slashes and entries
		// carry the AWS member names.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:    types.OpReplace,
					Path:  aws.String("/methodSettings"),
					Value: aws.String(`{"*/*":{"metricsEnabled":true,"loggingLevel":"ERROR"}}`),
				},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if len(stageResp.MethodSettings) != 1 {
			return fmt.Errorf("whole-member method settings replace did not replace the map, got %v", stageResp.MethodSettings)
		}
		setting, ok := stageResp.MethodSettings["*/*"]
		if !ok {
			return fmt.Errorf("wildcard method settings key missing, got %v", stageResp.MethodSettings)
		}
		if !setting.MetricsEnabled ||
			setting.LoggingLevel == nil || *setting.LoggingLevel != "ERROR" {
			return fmt.Errorf("method settings entry not decoded, got %+v", setting)
		}

		// Whole-member remove clears the members the patch table allows
		// remove on.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpRemove, Path: aws.String("/variables")},
				{Op: types.OpRemove, Path: aws.String("/canarySettings")},
			},
		})
		if err != nil {
			return err
		}
		stageResp, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if len(stageResp.Variables) != 0 {
			return fmt.Errorf("whole-member variables remove did not clear the map, got %v", stageResp.Variables)
		}
		if stageResp.CanarySettings != nil {
			return fmt.Errorf("whole-member canary remove did not clear the settings, got %v", stageResp.CanarySettings)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_UnknownPatchPath_Rejected", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
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

	results = append(results, r.RunTest("apigateway", "UpdateStage_UnsupportedPatchOps_Rejected", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		// The official patch table marks these (path, op) cells Not
		// supported: scalar rows are replace-only, per-key method
		// settings are replace-only, and canary scalar removes are not a
		// supported form.
		for _, po := range []types.PatchOperation{
			{Op: types.OpRemove, Path: aws.String("/description")},
			{Op: types.OpAdd, Path: aws.String("/cacheClusterEnabled"), Value: aws.String("true")},
			{Op: types.OpAdd, Path: aws.String("/deploymentId"), Value: aws.String("d")},
			{Op: types.OpRemove, Path: aws.String("/canarySettings/percentTraffic")},
			{Op: types.OpRemove, Path: aws.String("/methodSettings/~1pets~1GET")},
		} {
			_, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
				RestApiId:       aws.String(tc.apiID),
				StageName:       aws.String("test"),
				PatchOperations: []types.PatchOperation{po},
			})
			if err := AssertErrorContains(err, "BadRequestException"); err != nil {
				return fmt.Errorf("expected BadRequestException for op %s on %s, got: %v", po.Op, *po.Path, err)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_CopyCanaryDeployment", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		// Create a fresh deployment, point the canary at it, then promote
		// it with the documented copy operation (from
		// /canarySettings/deploymentId to /deploymentId).
		dep, err := tc.client.CreateDeployment(tc.ctx, &apigateway.CreateDeploymentInput{
			RestApiId:   aws.String(tc.apiID),
			Description: aws.String("canary candidate"),
		})
		if err != nil {
			return fmt.Errorf("create deployment: %v", err)
		}
		if dep.Id == nil || *dep.Id == "" {
			return fmt.Errorf("deployment id missing")
		}
		canaryDeployment := *dep.Id

		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/canarySettings/deploymentId"), Value: aws.String(canaryDeployment)},
				{Op: types.OpReplace, Path: aws.String("/canarySettings/percentTraffic"), Value: aws.String("100.0")},
			},
		})
		if err != nil {
			return err
		}

		stageResp, err := tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("get stage: %v", err)
		}
		if stageResp.DeploymentId == nil || *stageResp.DeploymentId == canaryDeployment {
			return fmt.Errorf("precondition failed: the stage already serves the canary deployment, got %v", stageResp.DeploymentId)
		}

		resp, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
			PatchOperations: []types.PatchOperation{
				{
					Op:   types.OpCopy,
					Path: aws.String("/deploymentId"),
					From: aws.String("/canarySettings/deploymentId"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.DeploymentId == nil || *resp.DeploymentId != canaryDeployment {
			return fmt.Errorf("canary deployment not promoted, got %v want %s", resp.DeploymentId, canaryDeployment)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteStage", func() error {
		if err := tc.require(tc.apiID); err != nil {
			return err
		}
		_, err := tc.client.DeleteStage(tc.ctx, &apigateway.DeleteStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetStage(tc.ctx, &apigateway.GetStageInput{
			RestApiId: aws.String(tc.apiID),
			StageName: aws.String("test"),
		})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetStage should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteDeployment", func() error {
		if err := tc.require(tc.apiID, deploymentID); err != nil {
			return err
		}
		_, err := tc.client.DeleteDeployment(tc.ctx, &apigateway.DeleteDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		_, err = tc.client.GetDeployment(tc.ctx, &apigateway.GetDeploymentInput{
			RestApiId:    aws.String(tc.apiID),
			DeploymentId: aws.String(deploymentID),
		})
		if aerr := AssertErrorContains(err, "NotFoundException"); aerr != nil {
			return fmt.Errorf("GetDeployment should fail with NotFoundException after delete: %v", aerr)
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "UpdateStage_AccessLogDestinationArn_Validated", func() error {
		apiID, _, err := tc.createOwnAPI("AlAPI")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteAPI(apiID)

		depID, err := tc.createDeployment(apiID, "")
		if err != nil {
			return fmt.Errorf("deploy: %v", err)
		}

		_, err = tc.client.CreateStage(tc.ctx, &apigateway.CreateStageInput{
			RestApiId:    aws.String(apiID),
			StageName:    aws.String("access-log"),
			DeploymentId: aws.String(depID),
		})
		if err != nil {
			return fmt.Errorf("create stage: %v", err)
		}

		// A Lambda function ARN is not a valid access log destination.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("access-log"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/accessLogSettings/destinationArn"),
					Value: aws.String("arn:aws:lambda:us-east-1:123456789012:function:logs")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for Lambda destination ARN, got: %v", err)
		}

		// Firehose delivery streams must begin with amazon-apigateway-.
		_, err = tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("access-log"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/accessLogSettings/destinationArn"),
					Value: aws.String("arn:aws:firehose:us-east-1:123456789012:deliverystream/my-stream")},
			},
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return fmt.Errorf("expected BadRequestException for non-prefixed Firehose ARN, got: %v", err)
		}

		// A CloudWatch Logs log group ARN is accepted.
		stageResp, err := tc.client.UpdateStage(tc.ctx, &apigateway.UpdateStageInput{
			RestApiId: aws.String(apiID),
			StageName: aws.String("access-log"),
			PatchOperations: []types.PatchOperation{
				{Op: types.OpReplace, Path: aws.String("/accessLogSettings/destinationArn"),
					Value: aws.String("arn:aws:logs:us-east-1:123456789012:log-group:my-api-logs")},
			},
		})
		if err != nil {
			return fmt.Errorf("update stage with valid logs ARN: %v", err)
		}
		if stageResp.AccessLogSettings == nil || stageResp.AccessLogSettings.DestinationArn == nil ||
			*stageResp.AccessLogSettings.DestinationArn != "arn:aws:logs:us-east-1:123456789012:log-group:my-api-logs" {
			return fmt.Errorf("accessLogSettings destinationArn not stored, got %v", stageResp.AccessLogSettings)
		}
		return nil
	}))

	return results
}

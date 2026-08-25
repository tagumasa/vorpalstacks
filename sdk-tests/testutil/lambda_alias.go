package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaAliasTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	funcName := tc.unique("AliasFunc")
	roleARN, cleanupRole, err := tc.createRole(tc.unique("AliasRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Alias_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupRole()

	_, cleanupFn, err := tc.createFunction(funcName, roleARN, lambdaFunctionCode)
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Alias_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupFn()

	results = append(results, tc.r.RunTest("lambda", "PublishVersion", func() error {
		resp, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return err
		}
		if resp.Version == nil {
			return fmt.Errorf("version is nil")
		}
		if *resp.Version != "1" {
			return fmt.Errorf("first published version should be 1, got %v", resp.Version)
		}
		if resp.FunctionName == nil || *resp.FunctionName != funcName {
			return fmt.Errorf("FunctionName mismatch, got %v", resp.FunctionName)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListVersionsByFunction", func() error {
		resp, err := tc.client.ListVersionsByFunction(tc.ctx, &lambda.ListVersionsByFunctionInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("versions list is nil")
		}
		hasLatest := false
		for _, v := range resp.Versions {
			if v.Version != nil && *v.Version == "$LATEST" {
				hasLatest = true
				break
			}
		}
		if !hasLatest {
			return fmt.Errorf("$LATEST version not found in ListVersionsByFunction")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateAlias", func() error {
		resp, err := tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(funcName),
			Name:            aws.String("live"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "live" {
			return fmt.Errorf("alias name mismatch, got %v", resp.Name)
		}
		if resp.FunctionVersion == nil || *resp.FunctionVersion != "$LATEST" {
			return fmt.Errorf("FunctionVersion mismatch, got %v", resp.FunctionVersion)
		}
		if resp.AliasArn == nil || *resp.AliasArn == "" {
			return fmt.Errorf("AliasArn is nil or empty")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetAlias", func() error {
		resp, err := tc.client.GetAlias(tc.ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("live"),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != "live" {
			return fmt.Errorf("alias name mismatch, got %v", resp.Name)
		}
		if resp.FunctionVersion == nil || *resp.FunctionVersion != "$LATEST" {
			return fmt.Errorf("FunctionVersion mismatch, got %v", resp.FunctionVersion)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateAlias", func() error {
		resp, err := tc.client.UpdateAlias(tc.ctx, &lambda.UpdateAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("live"),
			Description:  aws.String("Production alias"),
		})
		if err != nil {
			return err
		}
		if resp.Description == nil || *resp.Description != "Production alias" {
			return fmt.Errorf("Description not updated, got %v", resp.Description)
		}
		if resp.Name == nil || *resp.Name != "live" {
			return fmt.Errorf("alias name mismatch, got %v", resp.Name)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListAliases", func() error {
		resp, err := tc.client.ListAliases(tc.ctx, &lambda.ListAliasesInput{
			FunctionName: aws.String(funcName),
		})
		if err != nil {
			return err
		}
		if resp.Aliases == nil {
			return fmt.Errorf("aliases list is nil")
		}
		found := false
		for _, a := range resp.Aliases {
			if a.Name != nil && *a.Name == "live" {
				found = true
				if a.FunctionVersion == nil || *a.FunctionVersion != "$LATEST" {
					return fmt.Errorf("alias FunctionVersion mismatch, got %v", a.FunctionVersion)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("alias 'live' not found in ListAliases")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteAlias", func() error {
		_, err := tc.client.DeleteAlias(tc.ctx, &lambda.DeleteAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("live"),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetAlias(tc.ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("live"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "PublishVersion_VerifyVersion", func() error {
		pvFunc := tc.unique("PvFunc")
		pvRole, cleanupPvRole, err := tc.createRole(tc.unique("PvRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupPvRole()
		_, cleanupPvFn, err := tc.createFunction(pvFunc, pvRole, "exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanupPvFn()

		resp, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(pvFunc),
		})
		if err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if resp.Version == nil || *resp.Version == "$LATEST" {
			return fmt.Errorf("published version should not be $LATEST, got %v", resp.Version)
		}
		if resp.Version != nil && *resp.Version != "1" {
			return fmt.Errorf("first published version should be 1, got %v", resp.Version)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateAlias_DuplicateName", func() error {
		caFunc := tc.unique("CaFunc")
		caRole, cleanupCaRole, err := tc.createRole(tc.unique("CaRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupCaRole()
		_, cleanupCaFn, err := tc.createFunction(caFunc, caRole, "exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer cleanupCaFn()

		_, err = tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(caFunc),
			Name:            aws.String("prod"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err != nil {
			return fmt.Errorf("first alias: %v", err)
		}

		_, err = tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(caFunc),
			Name:            aws.String("prod"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err := AssertErrorContains(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetAlias_NonExistent", func() error {
		_, err := tc.client.GetAlias(tc.ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("nonexistent-alias-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "AliasQualifier_ReturnsPublishedVersionConfig", func() error {
		aqFunc := tc.unique("AqFunc")
		aqRole, cleanupAqRole, err := tc.createRole(tc.unique("AqRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupAqRole()
		_, cleanupAqFn, err := tc.createFunction(aqFunc, aqRole, "exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer cleanupAqFn()

		if _, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(aqFunc),
			Description:  aws.String("version one"),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if _, err := tc.client.UpdateFunctionConfiguration(tc.ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(aqFunc),
			MemorySize:   aws.Int32(256),
		}); err != nil {
			return fmt.Errorf("update configuration: %v", err)
		}
		if _, err := tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(aqFunc),
			Name:            aws.String("stable"),
			FunctionVersion: aws.String("1"),
		}); err != nil {
			return fmt.Errorf("create alias: %v", err)
		}

		// An alias qualifier reports the configuration of the published
		// version the alias points to, not the mutable $LATEST state.
		cfg, err := tc.client.GetFunctionConfiguration(tc.ctx, &lambda.GetFunctionConfigurationInput{
			FunctionName: aws.String(aqFunc),
			Qualifier:    aws.String("stable"),
		})
		if err != nil {
			return err
		}
		if cfg.Version == nil || *cfg.Version != "1" {
			return fmt.Errorf("alias-qualified Version should be 1, got %v", cfg.Version)
		}
		if cfg.Description == nil || *cfg.Description != "version one" {
			return fmt.Errorf("alias-qualified Description should be the published one, got %v", cfg.Description)
		}
		if cfg.MemorySize == nil || *cfg.MemorySize != 128 {
			return fmt.Errorf("alias-qualified MemorySize should be the published 128, got %v", cfg.MemorySize)
		}

		full, err := tc.client.GetFunction(tc.ctx, &lambda.GetFunctionInput{
			FunctionName: aws.String(aqFunc),
			Qualifier:    aws.String("stable"),
		})
		if err != nil {
			return err
		}
		if full.Configuration == nil || full.Configuration.Version == nil || *full.Configuration.Version != "1" {
			return fmt.Errorf("GetFunction with alias should carry version 1 configuration")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateAlias_RoutingWeightValidation", func() error {
		if _, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(funcName),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}

		// A weight outside [0, 1] is rejected.
		_, err := tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(funcName),
			Name:            aws.String("weighted"),
			FunctionVersion: aws.String("1"),
			RoutingConfig: &types.AliasRoutingConfiguration{
				AdditionalVersionWeights: map[string]float64{"2": 1.5},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
			return err
		}

		// Routing may only target published versions.
		_, err = tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(funcName),
			Name:            aws.String("weighted"),
			FunctionVersion: aws.String("1"),
			RoutingConfig: &types.AliasRoutingConfiguration{
				AdditionalVersionWeights: map[string]float64{"9": 0.5},
			},
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}

		// A valid routing config is accepted and echoed.
		resp, err := tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(funcName),
			Name:            aws.String("weighted"),
			FunctionVersion: aws.String("1"),
			RoutingConfig: &types.AliasRoutingConfiguration{
				AdditionalVersionWeights: map[string]float64{"2": 0.5},
			},
		})
		if err != nil {
			return err
		}
		if resp.RoutingConfig == nil || resp.RoutingConfig.AdditionalVersionWeights["2"] != 0.5 {
			return fmt.Errorf("RoutingConfig not echoed, got %v", resp.RoutingConfig)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateAlias_ClearDescription", func() error {
		if _, err := tc.client.UpdateAlias(tc.ctx, &lambda.UpdateAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("weighted"),
			Description:  aws.String("to be cleared"),
		}); err != nil {
			return fmt.Errorf("set description: %v", err)
		}
		resp, err := tc.client.UpdateAlias(tc.ctx, &lambda.UpdateAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("weighted"),
			Description:  aws.String(""),
		})
		if err != nil {
			return err
		}
		if resp.Description != nil && *resp.Description != "" {
			return fmt.Errorf("empty Description should clear the value, got %q", *resp.Description)
		}
		return nil
	}))

	return results
}

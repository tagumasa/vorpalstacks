package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaAliasTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	cwlClient *cloudwatchlogs.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
) []TestResult {
	var results []TestResult

	funcName := fmt.Sprintf("AliasFunc-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("AliasRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

	if err := createIAMRole(roleName); err != nil {
		return []TestResult{{Service: "lambda", TestName: "Alias_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer deleteIAMRole(roleName)

	zipCode, err := zipLambdaCode(lambdaFunctionCode)
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Alias_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to zip lambda code: %v", err)}}
	}

	_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(funcName),
		Runtime:      types.RuntimeNodejs22x,
		Role:         aws.String(roleARN),
		Handler:      aws.String("index.handler"),
		Code:         &types.FunctionCode{ZipFile: zipCode},
	})
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Alias_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(funcName)})
	defer deleteLambdaLogGroup(cwlClient, ctx, funcName)

	results = append(results, r.RunTest("lambda", "PublishVersion", func() error {
		resp, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
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

	results = append(results, r.RunTest("lambda", "ListVersionsByFunction", func() error {
		resp, err := client.ListVersionsByFunction(ctx, &lambda.ListVersionsByFunctionInput{
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

	results = append(results, r.RunTest("lambda", "CreateAlias", func() error {
		resp, err := client.CreateAlias(ctx, &lambda.CreateAliasInput{
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

	results = append(results, r.RunTest("lambda", "GetAlias", func() error {
		resp, err := client.GetAlias(ctx, &lambda.GetAliasInput{
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

	results = append(results, r.RunTest("lambda", "UpdateAlias", func() error {
		resp, err := client.UpdateAlias(ctx, &lambda.UpdateAliasInput{
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

	results = append(results, r.RunTest("lambda", "ListAliases", func() error {
		resp, err := client.ListAliases(ctx, &lambda.ListAliasesInput{
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

	results = append(results, r.RunTest("lambda", "DeleteAlias", func() error {
		_, err := client.DeleteAlias(ctx, &lambda.DeleteAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("live"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetAlias(ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("live"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "PublishVersion_VerifyVersion", func() error {
		pvFunc := fmt.Sprintf("PvFunc-%d", time.Now().UnixNano())
		pvRoleName := fmt.Sprintf("PvRole-%d", time.Now().UnixNano())
		pvRole := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, pvRoleName)
		pvCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("zip lambda code: %v", err)
		}
		if err := createIAMRole(pvRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(pvRoleName)
		_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(pvFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(pvRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: pvCode},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(pvFunc)})
		defer deleteLambdaLogGroup(cwlClient, ctx, pvFunc)

		resp, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
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

	results = append(results, r.RunTest("lambda", "CreateAlias_DuplicateName", func() error {
		caFunc := fmt.Sprintf("CaFunc-%d", time.Now().UnixNano())
		caRoleName := fmt.Sprintf("CaRole-%d", time.Now().UnixNano())
		caRole := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, caRoleName)
		caCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("zip lambda code: %v", err)
		}
		if err := createIAMRole(caRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(caRoleName)
		_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(caFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(caRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: caCode},
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(caFunc)})
		defer deleteLambdaLogGroup(cwlClient, ctx, caFunc)

		_, err = client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(caFunc),
			Name:            aws.String("prod"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err != nil {
			return fmt.Errorf("first alias: %v", err)
		}

		_, err = client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(caFunc),
			Name:            aws.String("prod"),
			FunctionVersion: aws.String("$LATEST"),
		})
		if err := AssertErrorContains(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetAlias_NonExistent", func() error {
		_, err := client.GetAlias(ctx, &lambda.GetAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("nonexistent-alias-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "AliasQualifier_ReturnsPublishedVersionConfig", func() error {
		aqFunc := fmt.Sprintf("AqFunc-%d", time.Now().UnixNano())
		aqRoleName := fmt.Sprintf("AqRole-%d", time.Now().UnixNano())
		aqRole := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, aqRoleName)
		aqCode, err := zipLambdaCode("exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("zip lambda code: %v", err)
		}
		if err := createIAMRole(aqRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(aqRoleName)
		_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(aqFunc),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(aqRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: aqCode},
		})
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(aqFunc)})
		defer deleteLambdaLogGroup(cwlClient, ctx, aqFunc)

		if _, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(aqFunc),
			Description:  aws.String("version one"),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}
		if _, err := client.UpdateFunctionConfiguration(ctx, &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(aqFunc),
			MemorySize:   aws.Int32(256),
		}); err != nil {
			return fmt.Errorf("update configuration: %v", err)
		}
		if _, err := client.CreateAlias(ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(aqFunc),
			Name:            aws.String("stable"),
			FunctionVersion: aws.String("1"),
		}); err != nil {
			return fmt.Errorf("create alias: %v", err)
		}

		// An alias qualifier reports the configuration of the published
		// version the alias points to, not the mutable $LATEST state.
		cfg, err := client.GetFunctionConfiguration(ctx, &lambda.GetFunctionConfigurationInput{
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

		full, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
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

	results = append(results, r.RunTest("lambda", "CreateAlias_RoutingWeightValidation", func() error {
		if _, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
			FunctionName: aws.String(funcName),
		}); err != nil {
			return fmt.Errorf("publish: %v", err)
		}

		// A weight outside [0, 1] is rejected.
		_, err := client.CreateAlias(ctx, &lambda.CreateAliasInput{
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
		_, err = client.CreateAlias(ctx, &lambda.CreateAliasInput{
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
		resp, err := client.CreateAlias(ctx, &lambda.CreateAliasInput{
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

	results = append(results, r.RunTest("lambda", "UpdateAlias_ClearDescription", func() error {
		if _, err := client.UpdateAlias(ctx, &lambda.UpdateAliasInput{
			FunctionName: aws.String(funcName),
			Name:         aws.String("weighted"),
			Description:  aws.String("to be cleared"),
		}); err != nil {
			return fmt.Errorf("set description: %v", err)
		}
		resp, err := client.UpdateAlias(ctx, &lambda.UpdateAliasInput{
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

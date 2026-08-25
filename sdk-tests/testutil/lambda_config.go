package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaConfigTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	results = append(results, runLambdaProvisionedConcurrencyTests(tc)...)
	results = append(results, runLambdaEventInvokeConfigTests(tc)...)
	results = append(results, runLambdaFunctionUrlConfigTests(tc)...)
	results = append(results, runLambdaResponseStreamTests(tc)...)
	results = append(results, runLambdaConfigErrorTests(tc)...)

	return results
}

func runLambdaProvisionedConcurrencyTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	pcFuncName := tc.unique("PcFunc")
	pcRole, cleanupPcRole, err := tc.createRole(tc.unique("PcRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupPcRole()

	_, cleanupPcFn, err := tc.createFunction(pcFuncName, pcRole, "exports.handler = async () => { return 1; };")
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupPcFn()

	publishResp, err := tc.client.PublishVersion(tc.ctx, &lambda.PublishVersionInput{
		FunctionName: aws.String(pcFuncName),
	})
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to publish version: %v", err)}}
	}
	pcVersion := *publishResp.Version

	results = append(results, tc.r.RunTest("lambda", "PutProvisionedConcurrencyConfig", func() error {
		resp, err := tc.client.PutProvisionedConcurrencyConfig(tc.ctx, &lambda.PutProvisionedConcurrencyConfigInput{
			FunctionName:                    aws.String(pcFuncName),
			Qualifier:                       aws.String(pcVersion),
			ProvisionedConcurrentExecutions: aws.Int32(5),
		})
		if err != nil {
			return err
		}
		if resp.AllocatedProvisionedConcurrentExecutions == nil {
			return fmt.Errorf("AllocatedProvisionedConcurrentExecutions is nil")
		}
		if resp.RequestedProvisionedConcurrentExecutions == nil || *resp.RequestedProvisionedConcurrentExecutions != 5 {
			return fmt.Errorf("RequestedProvisionedConcurrentExecutions mismatch, expected 5, got %v", resp.RequestedProvisionedConcurrentExecutions)
		}
		if resp.Status == "" {
			return fmt.Errorf("Status is empty")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "PutProvisionedConcurrencyConfig_RejectsLatest", func() error {
		_, err := tc.client.PutProvisionedConcurrencyConfig(tc.ctx, &lambda.PutProvisionedConcurrencyConfigInput{
			FunctionName:                    aws.String(pcFuncName),
			Qualifier:                       aws.String("$LATEST"),
			ProvisionedConcurrentExecutions: aws.Int32(5),
		})
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "PutProvisionedConcurrencyConfig_UnknownQualifier", func() error {
		_, err := tc.client.PutProvisionedConcurrencyConfig(tc.ctx, &lambda.PutProvisionedConcurrencyConfigInput{
			FunctionName:                    aws.String(pcFuncName),
			Qualifier:                       aws.String("999"),
			ProvisionedConcurrentExecutions: aws.Int32(5),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetProvisionedConcurrencyConfig", func() error {
		resp, err := tc.client.GetProvisionedConcurrencyConfig(tc.ctx, &lambda.GetProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(pcFuncName),
			Qualifier:    aws.String(pcVersion),
		})
		if err != nil {
			return err
		}
		if resp.Status == "" {
			return fmt.Errorf("status is empty")
		}
		if resp.RequestedProvisionedConcurrentExecutions == nil || *resp.RequestedProvisionedConcurrentExecutions != 5 {
			return fmt.Errorf("RequestedProvisionedConcurrentExecutions mismatch, expected 5, got %v", resp.RequestedProvisionedConcurrentExecutions)
		}
		if resp.AllocatedProvisionedConcurrentExecutions == nil {
			return fmt.Errorf("AllocatedProvisionedConcurrentExecutions is nil")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListProvisionedConcurrencyConfigs", func() error {
		resp, err := tc.client.ListProvisionedConcurrencyConfigs(tc.ctx, &lambda.ListProvisionedConcurrencyConfigsInput{
			FunctionName: aws.String(pcFuncName),
		})
		if err != nil {
			return err
		}
		if resp.ProvisionedConcurrencyConfigs == nil {
			return fmt.Errorf("configs list is nil")
		}
		if len(resp.ProvisionedConcurrencyConfigs) == 0 {
			return fmt.Errorf("expected at least 1 config")
		}
		found := false
		for _, c := range resp.ProvisionedConcurrencyConfigs {
			if c.FunctionArn != nil && *c.FunctionArn != "" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no valid config found in ListProvisionedConcurrencyConfigs")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteProvisionedConcurrencyConfig", func() error {
		_, err := tc.client.DeleteProvisionedConcurrencyConfig(tc.ctx, &lambda.DeleteProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(pcFuncName),
			Qualifier:    aws.String(pcVersion),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetProvisionedConcurrencyConfig(tc.ctx, &lambda.GetProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(pcFuncName),
			Qualifier:    aws.String(pcVersion),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetProvisionedConcurrencyConfig_NonExistent", func() error {
		_, err := tc.client.GetProvisionedConcurrencyConfig(tc.ctx, &lambda.GetProvisionedConcurrencyConfigInput{
			FunctionName: aws.String(pcFuncName),
			Qualifier:    aws.String(pcVersion),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

func runLambdaEventInvokeConfigTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	eicFuncName := tc.unique("EicFunc")
	eicRole, cleanupEicRole, err := tc.createRole(tc.unique("EicRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "PutFunctionEventInvokeConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupEicRole()

	_, cleanupEicFn, err := tc.createFunction(eicFuncName, eicRole, "exports.handler = async () => { return 1; };")
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "PutFunctionEventInvokeConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupEicFn()

	results = append(results, tc.r.RunTest("lambda", "PutFunctionEventInvokeConfig", func() error {
		maxAge := int32(3600)
		maxRetries := int32(2)
		resp, err := tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:             aws.String(eicFuncName),
			MaximumEventAgeInSeconds: aws.Int32(maxAge),
			MaximumRetryAttempts:     aws.Int32(maxRetries),
		})
		if err != nil {
			return err
		}
		if resp.LastModified == nil {
			return fmt.Errorf("LastModified is nil")
		}
		if resp.MaximumEventAgeInSeconds == nil || *resp.MaximumEventAgeInSeconds != 3600 {
			return fmt.Errorf("MaximumEventAgeInSeconds mismatch, got %v", resp.MaximumEventAgeInSeconds)
		}
		if resp.MaximumRetryAttempts == nil || *resp.MaximumRetryAttempts != 2 {
			return fmt.Errorf("MaximumRetryAttempts mismatch, got %v", resp.MaximumRetryAttempts)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunctionEventInvokeConfig", func() error {
		resp, err := tc.client.GetFunctionEventInvokeConfig(tc.ctx, &lambda.GetFunctionEventInvokeConfigInput{
			FunctionName: aws.String(eicFuncName),
		})
		if err != nil {
			return err
		}
		if resp.MaximumEventAgeInSeconds == nil || *resp.MaximumEventAgeInSeconds != 3600 {
			return fmt.Errorf("MaximumEventAgeInSeconds mismatch, got %v", resp.MaximumEventAgeInSeconds)
		}
		if resp.MaximumRetryAttempts == nil || *resp.MaximumRetryAttempts != 2 {
			return fmt.Errorf("MaximumRetryAttempts mismatch, got %v", resp.MaximumRetryAttempts)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "PutFunctionEventInvokeConfig_Defaults", func() error {
		// A Put that omits the retry and age members applies the AWS
		// defaults (2 retries, 6 hours) instead of zeros.
		resp, err := tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName: aws.String(eicFuncName),
		})
		if err != nil {
			return err
		}
		if resp.MaximumRetryAttempts == nil || *resp.MaximumRetryAttempts != 2 {
			return fmt.Errorf("default MaximumRetryAttempts should be 2, got %v", resp.MaximumRetryAttempts)
		}
		if resp.MaximumEventAgeInSeconds == nil || *resp.MaximumEventAgeInSeconds != 21600 {
			return fmt.Errorf("default MaximumEventAgeInSeconds should be 21600, got %v", resp.MaximumEventAgeInSeconds)
		}
		// Restore the explicit values the subsequent list test
		// expects to observe.
		_, err = tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:             aws.String(eicFuncName),
			MaximumEventAgeInSeconds: aws.Int32(3600),
			MaximumRetryAttempts:     aws.Int32(2),
		})
		return err
	}))

	results = append(results, tc.r.RunTest("lambda", "ListFunctionEventInvokeConfigs", func() error {
		resp, err := tc.client.ListFunctionEventInvokeConfigs(tc.ctx, &lambda.ListFunctionEventInvokeConfigsInput{
			FunctionName: aws.String(eicFuncName),
		})
		if err != nil {
			return err
		}
		if resp.FunctionEventInvokeConfigs == nil {
			return fmt.Errorf("configs list is nil")
		}
		if len(resp.FunctionEventInvokeConfigs) == 0 {
			return fmt.Errorf("expected at least 1 config")
		}
		cfg := resp.FunctionEventInvokeConfigs[0]
		if cfg.MaximumEventAgeInSeconds == nil || *cfg.MaximumEventAgeInSeconds != 3600 {
			return fmt.Errorf("MaximumEventAgeInSeconds mismatch in list, got %v", cfg.MaximumEventAgeInSeconds)
		}
		if cfg.MaximumRetryAttempts == nil || *cfg.MaximumRetryAttempts != 2 {
			return fmt.Errorf("MaximumRetryAttempts mismatch in list, got %v", cfg.MaximumRetryAttempts)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteFunctionEventInvokeConfig", func() error {
		_, err := tc.client.DeleteFunctionEventInvokeConfig(tc.ctx, &lambda.DeleteFunctionEventInvokeConfigInput{
			FunctionName: aws.String(eicFuncName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetFunctionEventInvokeConfig(tc.ctx, &lambda.GetFunctionEventInvokeConfigInput{
			FunctionName: aws.String(eicFuncName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionEventInvokeConfig_PartialUpdate", func() error {
		resp, err := tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:             aws.String(eicFuncName),
			MaximumEventAgeInSeconds: aws.Int32(3600),
			MaximumRetryAttempts:     aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if resp.MaximumRetryAttempts == nil || *resp.MaximumRetryAttempts != 2 {
			return fmt.Errorf("setup: MaximumRetryAttempts mismatch, got %v", resp.MaximumRetryAttempts)
		}

		updateResp, err := tc.client.UpdateFunctionEventInvokeConfig(tc.ctx, &lambda.UpdateFunctionEventInvokeConfigInput{
			FunctionName:         aws.String(eicFuncName),
			MaximumRetryAttempts: aws.Int32(0),
		})
		if err != nil {
			return err
		}
		if updateResp.MaximumRetryAttempts == nil || *updateResp.MaximumRetryAttempts != 0 {
			return fmt.Errorf("MaximumRetryAttempts not updated, got %v", updateResp.MaximumRetryAttempts)
		}
		if updateResp.MaximumEventAgeInSeconds == nil || *updateResp.MaximumEventAgeInSeconds != 3600 {
			return fmt.Errorf("MaximumEventAgeInSeconds should be preserved, got %v", updateResp.MaximumEventAgeInSeconds)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionEventInvokeConfig_CreateIfAbsent", func() error {
		updateResp, err := tc.client.UpdateFunctionEventInvokeConfig(tc.ctx, &lambda.UpdateFunctionEventInvokeConfigInput{
			FunctionName:         aws.String(eicFuncName),
			MaximumRetryAttempts: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		if updateResp.MaximumRetryAttempts == nil || *updateResp.MaximumRetryAttempts != 1 {
			return fmt.Errorf("MaximumRetryAttempts mismatch, got %v", updateResp.MaximumRetryAttempts)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionEventInvokeConfig_DestinationAtomic", func() error {
		destName := tc.unique("dest")
		onSuccessArn := fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", tc.r.region, tc.r.accountID, destName)
		onFailureArn := fmt.Sprintf("arn:aws:sqs:%s:%s:%s", tc.r.region, tc.r.accountID, destName)
		_, err := tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName: aws.String(eicFuncName),
			DestinationConfig: &types.DestinationConfig{
				OnSuccess: &types.OnSuccess{
					Destination: aws.String(onSuccessArn),
				},
				OnFailure: &types.OnFailure{
					Destination: aws.String(onFailureArn),
				},
			},
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.UpdateFunctionEventInvokeConfig(tc.ctx, &lambda.UpdateFunctionEventInvokeConfigInput{
			FunctionName: aws.String(eicFuncName),
			DestinationConfig: &types.DestinationConfig{
				OnSuccess: &types.OnSuccess{
					Destination: aws.String(onSuccessArn),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.DestinationConfig == nil {
			return fmt.Errorf("DestinationConfig is nil")
		}
		if resp.DestinationConfig.OnSuccess == nil || *resp.DestinationConfig.OnSuccess.Destination != onSuccessArn {
			return fmt.Errorf("OnSuccess mismatch, got %v", resp.DestinationConfig.OnSuccess)
		}
		if resp.DestinationConfig.OnFailure != nil {
			return fmt.Errorf("OnFailure should be nil (atomic replacement), got %v", resp.DestinationConfig.OnFailure)
		}
		return nil
	}))

	return results
}

func runLambdaFunctionUrlConfigTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	furlFuncName := tc.unique("FurlFunc")
	furlRole, cleanupFurlRole, err := tc.createRole(tc.unique("FurlRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateFunctionUrlConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupFurlRole()

	_, cleanupFurlFn, err := tc.createFunction(furlFuncName, furlRole, "exports.handler = async () => { return 1; };")
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "CreateFunctionUrlConfig_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupFurlFn()

	results = append(results, tc.r.RunTest("lambda", "CreateFunctionUrlConfig", func() error {
		resp, err := tc.client.CreateFunctionUrlConfig(tc.ctx, &lambda.CreateFunctionUrlConfigInput{
			FunctionName: aws.String(furlFuncName),
			AuthType:     types.FunctionUrlAuthTypeNone,
		})
		if err != nil {
			return err
		}
		if resp.FunctionUrl == nil || *resp.FunctionUrl == "" {
			return fmt.Errorf("FunctionUrl is nil or empty")
		}
		if resp.AuthType != types.FunctionUrlAuthTypeNone {
			return fmt.Errorf("AuthType mismatch, got %v", resp.AuthType)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetFunctionUrlConfig", func() error {
		resp, err := tc.client.GetFunctionUrlConfig(tc.ctx, &lambda.GetFunctionUrlConfigInput{
			FunctionName: aws.String(furlFuncName),
		})
		if err != nil {
			return err
		}
		if resp.FunctionUrl == nil || *resp.FunctionUrl == "" {
			return fmt.Errorf("FunctionUrl is nil or empty")
		}
		if resp.AuthType != types.FunctionUrlAuthTypeNone {
			return fmt.Errorf("AuthType mismatch, expected NONE, got %v", resp.AuthType)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "UpdateFunctionUrlConfig", func() error {
		resp, err := tc.client.UpdateFunctionUrlConfig(tc.ctx, &lambda.UpdateFunctionUrlConfigInput{
			FunctionName: aws.String(furlFuncName),
			AuthType:     types.FunctionUrlAuthTypeAwsIam,
		})
		if err != nil {
			return err
		}
		if resp.AuthType != types.FunctionUrlAuthTypeAwsIam {
			return fmt.Errorf("AuthType not updated, got %v", resp.AuthType)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListFunctionUrlConfigs", func() error {
		resp, err := tc.client.ListFunctionUrlConfigs(tc.ctx, &lambda.ListFunctionUrlConfigsInput{
			FunctionName: aws.String(furlFuncName),
		})
		if err != nil {
			return err
		}
		if resp.FunctionUrlConfigs == nil {
			return fmt.Errorf("url configs list is nil")
		}
		if len(resp.FunctionUrlConfigs) == 0 {
			return fmt.Errorf("expected at least 1 url config")
		}
		cfg := resp.FunctionUrlConfigs[0]
		if cfg.AuthType != types.FunctionUrlAuthTypeAwsIam {
			return fmt.Errorf("AuthType mismatch in list, expected AWS_IAM, got %v", cfg.AuthType)
		}
		if cfg.FunctionUrl == nil || *cfg.FunctionUrl == "" {
			return fmt.Errorf("FunctionUrl is nil or empty in list")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteFunctionUrlConfig", func() error {
		_, err := tc.client.DeleteFunctionUrlConfig(tc.ctx, &lambda.DeleteFunctionUrlConfigInput{
			FunctionName: aws.String(furlFuncName),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetFunctionUrlConfig(tc.ctx, &lambda.GetFunctionUrlConfigInput{
			FunctionName: aws.String(furlFuncName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "CreateFunctionUrlConfig_QualifierValidation", func() error {
		uqFunc := tc.unique("UqFunc")
		uqRole, cleanupUqRole, err := tc.createRole(tc.unique("UqRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupUqRole()
		_, cleanupUqFn, err := tc.createFunction(uqFunc, uqRole, "exports.handler = async () => { return 1; };")
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer cleanupUqFn()
		if _, err := tc.client.CreateAlias(tc.ctx, &lambda.CreateAliasInput{
			FunctionName:    aws.String(uqFunc),
			Name:            aws.String("live"),
			FunctionVersion: aws.String("$LATEST"),
		}); err != nil {
			return fmt.Errorf("create alias: %v", err)
		}

		// The URL qualifier names an alias; a numeric version is
		// not a valid URL qualifier.
		_, err = tc.client.CreateFunctionUrlConfig(tc.ctx, &lambda.CreateFunctionUrlConfigInput{
			FunctionName: aws.String(uqFunc),
			AuthType:     types.FunctionUrlAuthTypeNone,
			Qualifier:    aws.String("1"),
		})
		if err := AssertErrorContains(err, "InvalidParameterValueException"); err != nil {
			return err
		}

		resp, err := tc.client.CreateFunctionUrlConfig(tc.ctx, &lambda.CreateFunctionUrlConfigInput{
			FunctionName: aws.String(uqFunc),
			AuthType:     types.FunctionUrlAuthTypeNone,
			Qualifier:    aws.String("live"),
		})
		if err != nil {
			return err
		}
		if resp.FunctionUrl == nil || *resp.FunctionUrl == "" {
			return fmt.Errorf("FunctionUrl is nil or empty")
		}
		return nil
	}))

	return results
}

func runLambdaResponseStreamTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	iaFuncName := tc.unique("IaFunc")
	iaRole, cleanupIaRole, err := tc.createRole(tc.unique("IaRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "ResponseStream_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupIaRole()

	_, cleanupIaFn, err := tc.createFunction(iaFuncName, iaRole, "exports.handler = async () => { return { statusCode: 200 }; };")
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "ResponseStream_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupIaFn()

	results = append(results, tc.r.RunTest("lambda", "InvokeWithResponseStream", func() error {
		resp, err := tc.client.InvokeWithResponseStream(tc.ctx, &lambda.InvokeWithResponseStreamInput{
			FunctionName: aws.String(iaFuncName),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		if resp.ResponseStreamContentType == nil {
			return fmt.Errorf("ResponseStreamContentType is nil")
		}
		return nil
	}))

	return results
}

func runLambdaConfigErrorTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	results = append(results, tc.r.RunTest("lambda", "GetFunctionUrlConfig_NoConfig", func() error {
		nofcFuncName := tc.unique("NofcFunc")
		nofcRole, cleanupNofcRole, err := tc.createRole(tc.unique("NofcRole"))
		if err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer cleanupNofcRole()
		_, cleanupNofcFn, err := tc.createFunction(nofcFuncName, nofcRole, "",
			func(input *lambda.CreateFunctionInput) {
				input.Code = &types.FunctionCode{ZipFile: []byte("code")}
			})
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer cleanupNofcFn()

		_, err = tc.client.GetFunctionUrlConfig(tc.ctx, &lambda.GetFunctionUrlConfigInput{
			FunctionName: aws.String(nofcFuncName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "PutFunctionEventInvokeConfig_NonExistent", func() error {
		_, err := tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:             aws.String("nonexistent-func-xyz-123"),
			MaximumEventAgeInSeconds: aws.Int32(3600),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

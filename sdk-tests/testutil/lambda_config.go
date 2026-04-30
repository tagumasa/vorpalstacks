package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaConfigTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
) []TestResult {
	var results []TestResult

	// --- Provisioned Concurrency ---

	pcFuncName := fmt.Sprintf("PcFunc-%d", time.Now().UnixNano())
	pcRoleName := fmt.Sprintf("PcRole-%d", time.Now().UnixNano())
	pcRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", pcRoleName)
	pcCode := []byte("exports.handler = async () => { return 1; };")

	if err := createIAMRole(pcRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(pcRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(pcFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(pcRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: pcCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(pcFuncName)})

			publishResp, err := client.PublishVersion(ctx, &lambda.PublishVersionInput{
				FunctionName: aws.String(pcFuncName),
			})
			if err != nil {
				results = append(results, TestResult{Service: "lambda", TestName: "PutProvisionedConcurrencyConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to publish version: %v", err)})
			} else {
				pcVersion := *publishResp.Version

				results = append(results, r.RunTest("lambda", "PutProvisionedConcurrencyConfig", func() error {
					resp, err := client.PutProvisionedConcurrencyConfig(ctx, &lambda.PutProvisionedConcurrencyConfigInput{
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

				results = append(results, r.RunTest("lambda", "GetProvisionedConcurrencyConfig", func() error {
					resp, err := client.GetProvisionedConcurrencyConfig(ctx, &lambda.GetProvisionedConcurrencyConfigInput{
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

				results = append(results, r.RunTest("lambda", "ListProvisionedConcurrencyConfigs", func() error {
					resp, err := client.ListProvisionedConcurrencyConfigs(ctx, &lambda.ListProvisionedConcurrencyConfigsInput{
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

				results = append(results, r.RunTest("lambda", "DeleteProvisionedConcurrencyConfig", func() error {
					_, err := client.DeleteProvisionedConcurrencyConfig(ctx, &lambda.DeleteProvisionedConcurrencyConfigInput{
						FunctionName: aws.String(pcFuncName),
						Qualifier:    aws.String(pcVersion),
					})
					if err != nil {
						return err
					}
					_, err = client.GetProvisionedConcurrencyConfig(ctx, &lambda.GetProvisionedConcurrencyConfigInput{
						FunctionName: aws.String(pcFuncName),
						Qualifier:    aws.String(pcVersion),
					})
					if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
						return err
					}
					return nil
				}))

				results = append(results, r.RunTest("lambda", "GetProvisionedConcurrencyConfig_NonExistent", func() error {
					_, err := client.GetProvisionedConcurrencyConfig(ctx, &lambda.GetProvisionedConcurrencyConfigInput{
						FunctionName: aws.String(pcFuncName),
						Qualifier:    aws.String(pcVersion),
					})
					if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
						return err
					}
					return nil
				}))
			}
		}
	}

	// --- Event Invoke Config ---

	eicFuncName := fmt.Sprintf("EicFunc-%d", time.Now().UnixNano())
	eicRoleName := fmt.Sprintf("EicRole-%d", time.Now().UnixNano())
	eicRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", eicRoleName)
	eicCode := []byte("exports.handler = async () => { return 1; };")

	if err := createIAMRole(eicRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "PutFunctionEventInvokeConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(eicRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(eicFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(eicRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: eicCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "PutFunctionEventInvokeConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(eicFuncName)})

			results = append(results, r.RunTest("lambda", "PutFunctionEventInvokeConfig", func() error {
				maxAge := int32(3600)
				maxRetries := int32(2)
				resp, err := client.PutFunctionEventInvokeConfig(ctx, &lambda.PutFunctionEventInvokeConfigInput{
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

			results = append(results, r.RunTest("lambda", "GetFunctionEventInvokeConfig", func() error {
				resp, err := client.GetFunctionEventInvokeConfig(ctx, &lambda.GetFunctionEventInvokeConfigInput{
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

			results = append(results, r.RunTest("lambda", "ListFunctionEventInvokeConfigs", func() error {
				resp, err := client.ListFunctionEventInvokeConfigs(ctx, &lambda.ListFunctionEventInvokeConfigsInput{
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

			results = append(results, r.RunTest("lambda", "DeleteFunctionEventInvokeConfig", func() error {
				_, err := client.DeleteFunctionEventInvokeConfig(ctx, &lambda.DeleteFunctionEventInvokeConfigInput{
					FunctionName: aws.String(eicFuncName),
				})
				if err != nil {
					return err
				}
				_, err = client.GetFunctionEventInvokeConfig(ctx, &lambda.GetFunctionEventInvokeConfigInput{
					FunctionName: aws.String(eicFuncName),
				})
				if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
					return err
				}
				return nil
			}))
		}
	}

	// --- Function URL Config ---

	furlFuncName := fmt.Sprintf("FurlFunc-%d", time.Now().UnixNano())
	furlRoleName := fmt.Sprintf("FurlRole-%d", time.Now().UnixNano())
	furlRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", furlRoleName)
	furlCode := []byte("exports.handler = async () => { return 1; };")

	if err := createIAMRole(furlRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "CreateFunctionUrlConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(furlRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(furlFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(furlRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: furlCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "CreateFunctionUrlConfig_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(furlFuncName)})

			results = append(results, r.RunTest("lambda", "CreateFunctionUrlConfig", func() error {
				resp, err := client.CreateFunctionUrlConfig(ctx, &lambda.CreateFunctionUrlConfigInput{
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

			results = append(results, r.RunTest("lambda", "GetFunctionUrlConfig", func() error {
				resp, err := client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
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

			results = append(results, r.RunTest("lambda", "UpdateFunctionUrlConfig", func() error {
				resp, err := client.UpdateFunctionUrlConfig(ctx, &lambda.UpdateFunctionUrlConfigInput{
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

			results = append(results, r.RunTest("lambda", "ListFunctionUrlConfigs", func() error {
				resp, err := client.ListFunctionUrlConfigs(ctx, &lambda.ListFunctionUrlConfigsInput{
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

			results = append(results, r.RunTest("lambda", "DeleteFunctionUrlConfig", func() error {
				_, err := client.DeleteFunctionUrlConfig(ctx, &lambda.DeleteFunctionUrlConfigInput{
					FunctionName: aws.String(furlFuncName),
				})
				if err != nil {
					return err
				}
				_, err = client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
					FunctionName: aws.String(furlFuncName),
				})
				if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
					return err
				}
				return nil
			}))
		}
	}

	// --- Response Stream ---

	iaFuncName := fmt.Sprintf("IaFunc-%d", time.Now().UnixNano())
	iaRoleName := fmt.Sprintf("IaRole-%d", time.Now().UnixNano())
	iaRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", iaRoleName)
	iaCode := []byte("exports.handler = async () => { return { statusCode: 200 }; };")

	if err := createIAMRole(iaRoleName); err != nil {
		results = append(results, TestResult{Service: "lambda", TestName: "ResponseStream_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create IAM role: %v", err)})
	} else {
		defer deleteIAMRole(iaRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(iaFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(iaRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: iaCode},
		})
		if err != nil {
			results = append(results, TestResult{Service: "lambda", TestName: "ResponseStream_Setup", Status: "FAIL", Error: fmt.Sprintf("Failed to create function: %v", err)})
		} else {
			defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(iaFuncName)})

			results = append(results, r.RunTest("lambda", "InvokeWithResponseStream", func() error {
				resp, err := client.InvokeWithResponseStream(ctx, &lambda.InvokeWithResponseStreamInput{
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
		}
	}

	// --- Config Error Cases ---

	results = append(results, r.RunTest("lambda", "GetFunctionUrlConfig_NoConfig", func() error {
		nofcFuncName := fmt.Sprintf("NofcFunc-%d", time.Now().UnixNano())
		nofcRoleName := fmt.Sprintf("NofcRole-%d", time.Now().UnixNano())
		nofcRole := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", nofcRoleName)
		if err := createIAMRole(nofcRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		defer deleteIAMRole(nofcRoleName)
		_, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
			FunctionName: aws.String(nofcFuncName),
			Runtime:      types.RuntimeNodejs22x,
			Role:         aws.String(nofcRole),
			Handler:      aws.String("index.handler"),
			Code:         &types.FunctionCode{ZipFile: []byte("code")},
		})
		if err != nil {
			return fmt.Errorf("create function: %v", err)
		}
		defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(nofcFuncName)})

		_, err = client.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
			FunctionName: aws.String(nofcFuncName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "PutFunctionEventInvokeConfig_NonExistent", func() error {
		_, err := client.PutFunctionEventInvokeConfig(ctx, &lambda.PutFunctionEventInvokeConfigInput{
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

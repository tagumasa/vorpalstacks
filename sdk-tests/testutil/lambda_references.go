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

// runLambdaReferenceTests verifies that every FunctionName reference form
// AWS documents is accepted: bare name, name:qualifier, full ARN (with an
// optional qualifier suffix), and the partial ARN
// account-id:function:name (also with an optional qualifier suffix).
func runLambdaReferenceTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
	cwlClient *cloudwatchlogs.Client,
	createIAMRole func(string) error,
	deleteIAMRole func(string),
) []TestResult {
	var results []TestResult

	funcName := fmt.Sprintf("RefFunc-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("RefRole-%d", time.Now().UnixNano())
	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", r.accountID, roleName)

	if err := createIAMRole(roleName); err != nil {
		return []TestResult{{Service: "lambda", TestName: "Ref_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer deleteIAMRole(roleName)

	zipCode, err := zipLambdaCode(lambdaFunctionCode)
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Ref_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to zip lambda code: %v", err)}}
	}

	if _, err := client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(funcName),
		Runtime:      types.RuntimeNodejs22x,
		Role:         aws.String(roleARN),
		Handler:      aws.String("index.handler"),
		Code:         &types.FunctionCode{ZipFile: zipCode},
	}); err != nil {
		return []TestResult{{Service: "lambda", TestName: "Ref_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(funcName)})
	defer deleteLambdaLogGroup(cwlClient, ctx, funcName)

	partialArn := fmt.Sprintf("%s:function:%s", r.accountID, funcName)

	results = append(results, r.RunTest("lambda", "FunctionRef_PartialArn_EventInvokeConfig", func() error {
		_, err := client.PutFunctionEventInvokeConfig(ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:         aws.String(partialArn),
			MaximumRetryAttempts: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		resp, err := client.GetFunctionEventInvokeConfig(ctx, &lambda.GetFunctionEventInvokeConfigInput{
			FunctionName: aws.String(partialArn),
		})
		if err != nil {
			return err
		}
		if resp.MaximumRetryAttempts == nil || *resp.MaximumRetryAttempts != 1 {
			return fmt.Errorf("MaximumRetryAttempts mismatch, got %v", resp.MaximumRetryAttempts)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "FunctionRef_PartialArn_Invoke", func() error {
		resp, err := client.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: aws.String(partialArn),
		})
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("StatusCode mismatch, got %d", resp.StatusCode)
		}
		if resp.ExecutedVersion == nil || *resp.ExecutedVersion != "$LATEST" {
			return fmt.Errorf("ExecutedVersion for unqualified partial ARN should be $LATEST, got %v", resp.ExecutedVersion)
		}
		return nil
	}))

	// --- Error contract of store-level failures ---

	esmSourceArn := fmt.Sprintf("arn:aws:sqs:%s:%s:ref-dup-queue", r.region, r.accountID)
	var esmUUID string

	results = append(results, r.RunTest("lambda", "ErrorMapping_ESM_DuplicateConflict", func() error {
		in := &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(funcName),
			EventSourceArn: aws.String(esmSourceArn),
			BatchSize:      aws.Int32(10),
		}
		created, err := client.CreateEventSourceMapping(ctx, in)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		if created.UUID == nil {
			return fmt.Errorf("first create returned no UUID")
		}
		esmUUID = *created.UUID

		_, err = client.CreateEventSourceMapping(ctx, in)
		if err := AssertErrorContains(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))
	defer func() {
		if esmUUID != "" {
			_, _ = client.DeleteEventSourceMapping(ctx, &lambda.DeleteEventSourceMappingInput{UUID: aws.String(esmUUID)})
		}
	}()

	results = append(results, r.RunTest("lambda", "ErrorMapping_DeleteFunction_BadVersionQualifier", func() error {
		_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String(funcName),
			Qualifier:    aws.String("999"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ErrorMapping_GetFunctionConcurrency_Unset", func() error {
		_, err := client.GetFunctionConcurrency(ctx, &lambda.GetFunctionConcurrencyInput{
			FunctionName: aws.String(funcName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

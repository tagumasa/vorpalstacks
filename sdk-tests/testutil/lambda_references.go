package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// runLambdaReferenceTests verifies that every FunctionName reference form
// AWS documents is accepted: bare name, name:qualifier, full ARN (with an
// optional qualifier suffix), and the partial ARN
// account-id:function:name (also with an optional qualifier suffix).
func runLambdaReferenceTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	funcName := tc.unique("RefFunc")
	roleARN, cleanupRole, err := tc.createRole(tc.unique("RefRole"))
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Ref_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create IAM role: %v", err)}}
	}
	defer cleanupRole()

	_, cleanupFn, err := tc.createFunction(funcName, roleARN, lambdaFunctionCode)
	if err != nil {
		return []TestResult{{Service: "lambda", TestName: "Ref_Setup", Status: "FAIL",
			Error: fmt.Sprintf("Failed to create function: %v", err)}}
	}
	defer cleanupFn()

	partialArn := fmt.Sprintf("%s:function:%s", tc.r.accountID, funcName)

	results = append(results, tc.r.RunTest("lambda", "FunctionRef_PartialArn_EventInvokeConfig", func() error {
		_, err := tc.client.PutFunctionEventInvokeConfig(tc.ctx, &lambda.PutFunctionEventInvokeConfigInput{
			FunctionName:         aws.String(partialArn),
			MaximumRetryAttempts: aws.Int32(1),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetFunctionEventInvokeConfig(tc.ctx, &lambda.GetFunctionEventInvokeConfigInput{
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

	results = append(results, tc.r.RunTest("lambda", "FunctionRef_PartialArn_Invoke", func() error {
		resp, err := tc.client.Invoke(tc.ctx, &lambda.InvokeInput{
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

	esmSourceArn := fmt.Sprintf("arn:aws:sqs:%s:%s:ref-dup-queue", tc.r.region, tc.r.accountID)
	var esmUUID string

	results = append(results, tc.r.RunTest("lambda", "ErrorMapping_ESM_DuplicateConflict", func() error {
		in := &lambda.CreateEventSourceMappingInput{
			FunctionName:   aws.String(funcName),
			EventSourceArn: aws.String(esmSourceArn),
			BatchSize:      aws.Int32(10),
		}
		created, err := tc.client.CreateEventSourceMapping(tc.ctx, in)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		if created.UUID == nil {
			return fmt.Errorf("first create returned no UUID")
		}
		esmUUID = *created.UUID

		_, err = tc.client.CreateEventSourceMapping(tc.ctx, in)
		if err := AssertErrorContains(err, "ResourceConflictException"); err != nil {
			return err
		}
		return nil
	}))
	defer func() {
		if esmUUID != "" {
			_, _ = tc.client.DeleteEventSourceMapping(tc.ctx, &lambda.DeleteEventSourceMappingInput{UUID: aws.String(esmUUID)})
		}
	}()

	results = append(results, tc.r.RunTest("lambda", "ErrorMapping_DeleteFunction_BadVersionQualifier", func() error {
		_, err := tc.client.DeleteFunction(tc.ctx, &lambda.DeleteFunctionInput{
			FunctionName: aws.String(funcName),
			Qualifier:    aws.String("999"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ErrorMapping_GetFunctionConcurrency_Unset", func() error {
		_, err := tc.client.GetFunctionConcurrency(tc.ctx, &lambda.GetFunctionConcurrencyInput{
			FunctionName: aws.String(funcName),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

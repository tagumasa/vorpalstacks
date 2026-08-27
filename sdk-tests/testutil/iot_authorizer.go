package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// runIoTAuthorizerTests covers the custom-authorizer lifecycle: Create/
// Describe/Update/List/Delete plus the default-authorizer management
// (SetDefault/DescribeDefault/ClearDefault) and TestInvokeAuthorizer, with a
// NotFound negative path. All names use uniqueName.
func (r *TestRunner) runIoTAuthorizerTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	authName := uniqueName("authorizer")
	fnARN := "arn:aws:lambda:" + tc.region + ":" + tc.accountID + ":function:auth-fn"

	defer func() {
		// Clearing the default first avoids delete being blocked while default.
		tc.client.ClearDefaultAuthorizer(tc.ctx, &iot.ClearDefaultAuthorizerInput{})
		tc.client.DeleteAuthorizer(tc.ctx, &iot.DeleteAuthorizerInput{AuthorizerName: aws.String(authName)})
	}()

	results = append(results, r.RunTest("iot", "Authorizer_CreateAuthorizer", func() error {
		out, err := tc.client.CreateAuthorizer(tc.ctx, &iot.CreateAuthorizerInput{
			AuthorizerName:        aws.String(authName),
			AuthorizerFunctionArn: aws.String(fnARN),
			TokenKeyName:          aws.String("token"),
			Status:                "ACTIVE",
		})
		if err != nil {
			return fmt.Errorf("CreateAuthorizer failed: %w", err)
		}
		if out.AuthorizerName == nil || *out.AuthorizerName != authName {
			return fmt.Errorf("expected authorizerName=%s, got %v", authName, out.AuthorizerName)
		}
		if out.AuthorizerArn == nil || *out.AuthorizerArn == "" {
			return fmt.Errorf("expected non-empty authorizerArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_DescribeAuthorizer", func() error {
		out, err := tc.client.DescribeAuthorizer(tc.ctx, &iot.DescribeAuthorizerInput{AuthorizerName: aws.String(authName)})
		if err != nil {
			return fmt.Errorf("DescribeAuthorizer failed: %w", err)
		}
		if out.AuthorizerDescription == nil {
			return fmt.Errorf("expected non-nil authorizerDescription")
		}
		d := out.AuthorizerDescription
		if d.AuthorizerName == nil || *d.AuthorizerName != authName {
			return fmt.Errorf("expected authorizerName=%s, got %v", authName, d.AuthorizerName)
		}
		if d.AuthorizerFunctionArn == nil || *d.AuthorizerFunctionArn != fnARN {
			return fmt.Errorf("expected function ARN, got %v", d.AuthorizerFunctionArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_UpdateAuthorizer", func() error {
		fnARN2 := "arn:aws:lambda:" + tc.region + ":" + tc.accountID + ":function:auth-fn-v2"
		if _, err := tc.client.UpdateAuthorizer(tc.ctx, &iot.UpdateAuthorizerInput{
			AuthorizerName:        aws.String(authName),
			AuthorizerFunctionArn: aws.String(fnARN2),
		}); err != nil {
			return fmt.Errorf("UpdateAuthorizer failed: %w", err)
		}
		out, err := tc.client.DescribeAuthorizer(tc.ctx, &iot.DescribeAuthorizerInput{AuthorizerName: aws.String(authName)})
		if err != nil {
			return fmt.Errorf("DescribeAuthorizer after update failed: %w", err)
		}
		if aws.ToString(out.AuthorizerDescription.AuthorizerFunctionArn) != fnARN2 {
			return fmt.Errorf("expected updated function ARN")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_SetDefaultAuthorizer", func() error {
		out, err := tc.client.SetDefaultAuthorizer(tc.ctx, &iot.SetDefaultAuthorizerInput{AuthorizerName: aws.String(authName)})
		if err != nil {
			return fmt.Errorf("SetDefaultAuthorizer failed: %w", err)
		}
		if aws.ToString(out.AuthorizerArn) == "" {
			return fmt.Errorf("expected non-empty authorizerArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_SetDefaultAuthorizer_NotFoundRejected", func() error {
		_, err := tc.client.SetDefaultAuthorizer(tc.ctx, &iot.SetDefaultAuthorizerInput{AuthorizerName: aws.String(uniqueName("no-such-authorizer"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Authorizer_DescribeDefaultAuthorizer", func() error {
		// The response shape is the authorizerDescription structure, the
		// same as DescribeAuthorizer's output.
		out, err := tc.client.DescribeDefaultAuthorizer(tc.ctx, &iot.DescribeDefaultAuthorizerInput{})
		if err != nil {
			return fmt.Errorf("DescribeDefaultAuthorizer failed: %w", err)
		}
		if out.AuthorizerDescription == nil {
			return fmt.Errorf("expected non-nil authorizerDescription")
		}
		if aws.ToString(out.AuthorizerDescription.AuthorizerName) != authName {
			return fmt.Errorf("expected authorizerName=%s, got %v", authName, out.AuthorizerDescription.AuthorizerName)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Authorizer_TestInvokeAuthorizer", func() error {
		_, err := tc.client.TestInvokeAuthorizer(tc.ctx, &iot.TestInvokeAuthorizerInput{
			AuthorizerName: aws.String(authName),
			Token:          aws.String("test-token"),
			TokenSignature: aws.String("test-signature"),
		})
		return err
	}))

	results = append(results, r.RunTest("iot", "Authorizer_ClearDefaultAuthorizer", func() error {
		_, err := tc.client.ClearDefaultAuthorizer(tc.ctx, &iot.ClearDefaultAuthorizerInput{})
		return err
	}))

	results = append(results, r.RunTest("iot", "Authorizer_ListAuthorizers_IncludesCreated", func() error {
		out, err := tc.client.ListAuthorizers(tc.ctx, &iot.ListAuthorizersInput{})
		if err != nil {
			return fmt.Errorf("ListAuthorizers failed: %w", err)
		}
		for _, a := range out.Authorizers {
			if a.AuthorizerName != nil && *a.AuthorizerName == authName {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d authorizers", authName, len(out.Authorizers))
	}))

	results = append(results, r.RunTest("iot", "Authorizer_DescribeAuthorizer_NotFound", func() error {
		_, err := tc.client.DescribeAuthorizer(tc.ctx, &iot.DescribeAuthorizerInput{AuthorizerName: aws.String(uniqueName("nope-auth"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Authorizer_DeleteAuthorizer", func() error {
		_, err := tc.client.DeleteAuthorizer(tc.ctx, &iot.DeleteAuthorizerInput{AuthorizerName: aws.String(authName)})
		if err != nil {
			return fmt.Errorf("DeleteAuthorizer failed: %w", err)
		}
		return nil
	}))

	return results
}

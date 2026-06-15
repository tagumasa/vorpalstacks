package testutil

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

const testIoTPolicyDocument = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "iot:Connect",
      "Resource": "*"
    }
  ]
}`

func (r *TestRunner) runIoTPolicyTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "Policy_CreatePolicy", func() error {
		out, err := tc.client.CreatePolicy(tc.ctx, &iot.CreatePolicyInput{
			PolicyName:     aws.String("test-policy-1"),
			PolicyDocument: aws.String(testIoTPolicyDocument),
		})
		if err != nil {
			var httpErr interface{ HTTPStatusCode() int }
			if errors.As(err, &httpErr) && httpErr.HTTPStatusCode() == 409 {
				return nil
			}
			return fmt.Errorf("CreatePolicy failed: %w", err)
		}
		if out.PolicyName == nil || *out.PolicyName != "test-policy-1" {
			return fmt.Errorf("expected policyName=test-policy-1, got %v", out.PolicyName)
		}
		if out.PolicyArn == nil || *out.PolicyArn == "" {
			return fmt.Errorf("expected non-empty policyArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_GetPolicy", func() error {
		out, err := tc.client.GetPolicy(tc.ctx, &iot.GetPolicyInput{
			PolicyName: aws.String("test-policy-1"),
		})
		if err != nil {
			return fmt.Errorf("GetPolicy failed: %w", err)
		}
		if out.PolicyName == nil || *out.PolicyName != "test-policy-1" {
			return fmt.Errorf("expected policyName=test-policy-1, got %v", out.PolicyName)
		}
		if out.PolicyDocument == nil || *out.PolicyDocument == "" {
			return fmt.Errorf("expected non-empty policyDocument")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_ListPolicies", func() error {
		out, err := tc.client.ListPolicies(tc.ctx, &iot.ListPoliciesInput{})
		if err != nil {
			return fmt.Errorf("ListPolicies failed: %w", err)
		}
		if out.Policies == nil {
			return fmt.Errorf("expected non-nil policies list")
		}
		found := false
		for _, p := range out.Policies {
			if p.PolicyName != nil && *p.PolicyName == "test-policy-1" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("test-policy-1 not found in list of %d policies", len(out.Policies))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_AttachDetach", func() error {
		certOut, err := tc.client.CreateKeysAndCertificate(tc.ctx, &iot.CreateKeysAndCertificateInput{
			SetAsActive: true,
		})
		if err != nil {
			return fmt.Errorf("CreateKeysAndCertificate failed: %w", err)
		}

		_, err = tc.client.AttachPolicy(tc.ctx, &iot.AttachPolicyInput{
			PolicyName: aws.String("test-policy-1"),
			Target:     certOut.CertificateArn,
		})
		if err != nil {
			return fmt.Errorf("AttachPolicy failed: %w", err)
		}

		policies, err := tc.client.ListAttachedPolicies(tc.ctx, &iot.ListAttachedPoliciesInput{
			Target: certOut.CertificateArn,
		})
		if err != nil {
			return fmt.Errorf("ListAttachedPolicies failed: %w", err)
		}
		if len(policies.Policies) == 0 {
			return fmt.Errorf("expected at least one policy after AttachPolicy")
		}
		found := false
		for _, p := range policies.Policies {
			if p.PolicyName != nil && *p.PolicyName == "test-policy-1" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("test-policy-1 not found in attached policies")
		}

		_, err = tc.client.UpdateCertificate(tc.ctx, &iot.UpdateCertificateInput{
			CertificateId: certOut.CertificateId,
			NewStatus:     types.CertificateStatusInactive,
		})
		if err != nil {
			return fmt.Errorf("UpdateCertificate to INACTIVE failed: %w", err)
		}

		_, err = tc.client.DetachPolicy(tc.ctx, &iot.DetachPolicyInput{
			PolicyName: aws.String("test-policy-1"),
			Target:     certOut.CertificateArn,
		})
		if err != nil {
			return fmt.Errorf("DetachPolicy failed: %w", err)
		}

		_, err = tc.client.DeleteCertificate(tc.ctx, &iot.DeleteCertificateInput{
			CertificateId: certOut.CertificateId,
		})
		if err != nil {
			return fmt.Errorf("DeleteCertificate failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_DeletePolicy", func() error {
		_, _ = tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{
			PolicyName: aws.String("test-policy-1"),
		})
		return nil
	}))

	return results
}
